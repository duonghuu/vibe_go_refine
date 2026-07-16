# BACKEND PLAN: Quản trị danh sách sản phẩm (Admin Product List)

Dựa trên yêu cầu nghiệp vụ từ `01-product-list-idea.md` và kiến trúc chuẩn của hệ thống, dưới đây là bản quy hoạch thiết kế Backend chi tiết.

## Trụ cột 1: Thiết kế Dữ liệu (Database Schema Struct)

Tạo các file models trong `internal/domain/product/entity/` và `internal/domain/category/entity/`.

```go
package entity

import (
	"time"
	"gorm.io/gorm"
)

// Category đại diện cho danh mục sản phẩm (Hỗ trợ cấu trúc cây)
type Category struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null;index" json:"name"`
	ParentID  *uint          `gorm:"index" json:"parentId"` // Nullable cho root category
	Children  []Category     `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Product đại diện cho thực thể Sản phẩm
type Product struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	SKU         string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"sku"`
	Name        string         `gorm:"type:varchar(255);not null;index:idx_product_name" json:"name"`
	CategoryID  uint           `gorm:"not null;index" json:"categoryId"`
	Category    Category       `gorm:"foreignKey:CategoryID" json:"category"`
	Price       float64        `gorm:"type:decimal(10,2);not null" json:"price"` // Giá gốc
	SalePrice   *float64       `gorm:"type:decimal(10,2)" json:"salePrice"`      // Giá khuyến mãi
	Stock       int            `gorm:"type:int;not null;default:0" json:"stock"`
	SoldCount   int            `gorm:"type:int;not null;default:0" json:"soldCount"`
	Status      string         `gorm:"type:varchar(20);not null;default:'ACTIVE';index" json:"status"` // ACTIVE, HIDDEN, OUT_OF_STOCK
	ImageURL    string         `gorm:"type:text" json:"image"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
```

*Ràng buộc & Tối ưu:*
- `SKU` có `uniqueIndex` để truy vấn nhanh và tránh trùng lặp.
- `Name` và `Status` được đánh Index (`index`) để phục vụ tìm kiếm và filter.
- Áp dụng `gorm.DeletedAt` cho cơ chế Soft Delete, đảm bảo tính toàn vẹn dữ liệu cho các Order liên quan trong quá khứ.

## Trụ cột 2: Giao kèo API (API Contract & Context Auth)

Các API quản trị được bọc trong Router Group `/api/v1/admin/products` và **BẮT BUỘC** đi qua `AuthMiddleware` (Validate Access Token qua Redis Blacklist) và `RoleMiddleware` (Chỉ cấp quyền cho `ADMIN`, `STAFF`).

### 1. Lấy danh sách sản phẩm (Phân trang & Lọc)
- **Method & Route:** `GET /api/v1/admin/products`
- **Request Query Binding:**
```go
type ListProductReq struct {
	Page       int     `form:"_start" binding:"min=0"`
	PageSize   int     `form:"_end" binding:"min=1"`
	Sort       string  `form:"_sort"`  // field name
	Order      string  `form:"_order"` // ASC, DESC
	Search     string  `form:"q"`      // Search theo name, SKU
	CategoryID *uint   `form:"categoryId"`
	Status     string  `form:"status"`
	MinPrice   *float64 `form:"minPrice"`
	MaxPrice   *float64 `form:"maxPrice"`
	MinStock   *int    `form:"minStock"`
}
```
- **Response Payload:**
```go
type ListProductRes struct {
	Data  []Product `json:"data"`
	Total int64     `json:"total"` // Phục vụ Refine.js DataGrid
}
```

### 2. Thêm mới Sản phẩm
- **Method & Route:** `POST /api/v1/admin/products`
- **Request Body Binding:**
```go
type CreateProductReq struct {
	Name       string   `json:"name" binding:"required,max=255"`
	SKU        string   `json:"sku" binding:"required,max=50"`
	CategoryID uint     `json:"categoryId" binding:"required"`
	Price      float64  `json:"price" binding:"required,gt=0"`
	SalePrice  *float64 `json:"salePrice" binding:"omitempty,gte=0"`
	Stock      int      `json:"stock" binding:"required,gte=0"`
	Status     string   `json:"status" binding:"required,oneof=ACTIVE HIDDEN OUT_OF_STOCK"`
	ImageURL   string   `json:"image" binding:"omitempty,url"`
}
```
- **Response Payload:** `201 Created` kèm theo dữ liệu `Product` vừa tạo.

### 3. Cập nhật & Thao tác hàng loạt
- **Cập nhật:** `PUT /api/v1/admin/products/:id` (Full update qua `UpdateProductReq`).
- **Đổi trạng thái (Soft):** `PATCH /api/v1/admin/products/:id/status` (Body: `{"status": "HIDDEN"}`).
- **Xóa (Soft Delete):** `DELETE /api/v1/admin/products/:id`.
- **Thao tác hàng loạt (Bulk):**
  - `POST /api/v1/admin/products/bulk-delete` (Body: `{"ids": [1, 2, 3]}`)
  - `PATCH /api/v1/admin/products/bulk-status` (Body: `{"ids": [1, 2], "status": "HIDDEN"}`)

## Trụ cột 3: Xử lý Cache & Quản lý Trạng thái Token (Redis Integration)

### 1. Quản lý Auth & Token (Security Critical)
- **Xác thực Request (Middleware):** Khi gọi các API `/admin/products/*`, `AuthMiddleware` sẽ lấy Access Token từ Header, xác thực chữ ký JWT, parse `userId`/`role`. Dòng code tiếp theo sẽ `GET` từ Redis Blacklist. Nếu tồn tại -> chặn `401 Unauthorized`.
- **Logout & Blacklist:** 
  Khi Admin/Staff logout, lấy Access Token (từ Header) tính TTL còn lại và đưa vào Redis: `SETEX blacklist:token:<accessToken> <TTL> "1"`. Đồng thời `DEL refresh_token:jti:<userId>`.

### 2. Caching Dữ liệu Sản phẩm / Danh mục
Do đặc thù trang Quản trị (Admin Dashboard) yêu cầu dữ liệu theo thời gian thực sát nhất để quản lý Tồn kho (`Stock`), ta áp dụng chiến lược Cache như sau:
- **Products List / Detail:** KHÔNG CACHE toàn bộ danh sách ở tầng Admin (tránh sai lệch số liệu tồn kho). Mọi query được tối ưu thông qua Index của MySQL.
- **Categories:** Dữ liệu danh mục rất ít khi thay đổi.
  - Cấu trúc Key: `cache:categories:list` (Type: String - JSON).
  - TTL: 24 giờ.
  - Invalidation (Xóa cache): Mỗi khi có API tạo/sửa/xóa Category được gọi thành công, tự động thực thi `DEL cache:categories:list` trên Redis để ép hệ thống fetch lại từ DB trong lần request tiếp theo.
- **Thống kê tổng quan (Summary Cards):**
  - Card (Tổng số, Đang bán, Ngừng bán, Sắp hết hàng) yêu cầu Query khá nặng nếu dữ liệu lớn.
  - Cấu trúc Key: `cache:products:summary_metrics`
  - TTL: 5 phút.
  - Background Job (Cron): Có thể cân nhắc viết cron job nhỏ update Key này mỗi 5 phút hoặc xóa cache khi có sự kiện thay đổi dữ liệu lớn (Bulk update).
