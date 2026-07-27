# QUY HOẠCH KIẾN TRÚC BACK-END: THÊM MỚI SẢN PHẨM

**Dự án:** TechBite
**Module:** Product Management
**Tài liệu tham chiếu:** `.docs/ideas/dashboard/08-product-create-idea.md`

*(Ghi chú: Bản thiết kế này tập trung vào luồng tạo mới dữ liệu (Create) cho sản phẩm, đồng bộ với kiến trúc phân tầng DDD/Clean Architecture đang có).*

---

## TRỤ CỘT 1: THIẾT KẾ DỮ LIỆU (DATABASE SCHEMA STRUCT)

Thiết kế này sử dụng GORM Entity `Product` để lưu trữ thông tin sản phẩm, giá cả, và tồn kho. Nó có mối quan hệ BelongsTo với `Category` và hỗ trợ Soft Delete.

```go
package entity

import (
	"time"
	"gorm.io/gorm"
)

// Product đại diện cho bảng sản phẩm trong cơ sở dữ liệu
type Product struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	CategoryID  uint           `gorm:"index;not null" json:"categoryId"`
	Price       float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	SalePrice   *float64       `gorm:"type:decimal(10,2)" json:"salePrice"` // Nullable nếu không có giá khuyến mãi
	SKU         string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"sku"`
	Stock       int            `gorm:"type:int;default:0;not null" json:"stock"`
	SoldCount   int            `gorm:"type:int;default:0;not null" json:"soldCount"`
	ImageURL    string         `gorm:"type:varchar(500)" json:"image"` // Trả về dạng `image` để mapping dễ dàng với UI
	Status      string         `gorm:"type:varchar(20);default:'ACTIVE';index" json:"status"` // ACTIVE, HIDDEN
	
	// Quan hệ ngoại (Foreign Key)
	Category    Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`

	// Audit logs & Soft Delete
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Product) TableName() string {
	return "products"
}
```

### Ràng buộc & Tối ưu (Constraints & Optimization):
- **Khóa chính & Khóa ngoại:** `ID` tự tăng làm khóa chính. `CategoryID` tham chiếu đến `categories(id)` và được đánh index để tối ưu lọc sản phẩm theo danh mục.
- **Index:** `uniqueIndex` ở `SKU` để đảm bảo mã sản phẩm là duy nhất trên toàn hệ thống. `index` cho `Status` để tối ưu các truy vấn chỉ lấy sản phẩm "ACTIVE".
- **Kiểu dữ liệu Tiền tệ:** Sử dụng `decimal(10,2)` ở GORM tag cho cột Price và SalePrice để tránh sai số thập phân trong DB.
- **Soft Delete:** Tích hợp `gorm.DeletedAt` giúp bảo toàn toàn vẹn dữ liệu cho các Order lịch sử.

---

## TRỤ CỘT 2: GIAO KÈO API (API CONTRACT & CONTEXT AUTH)

### 1. Tạo mới sản phẩm
- **Method & Route:** `POST /api/v1/admin/products`
- **Auth:** Yêu cầu đi qua Middleware Auth (Kiểm tra token hợp lệ, lấy `userId`, phân quyền Admin/Operator).

**Request Binding Struct:**
```go
package dto

type CreateProductRequest struct {
	Name        string   `json:"name" binding:"required,max=255"`
	Description string   `json:"description"`
	CategoryID  uint     `json:"categoryId" binding:"required"`
	Price       float64  `json:"price" binding:"required,gt=0"`
	SalePrice   *float64 `json:"salePrice" binding:"omitempty,gte=0"`
	SKU         string   `json:"sku"` // Backend tự động sinh nếu rỗng
	Stock       int      `json:"stock" binding:"omitempty,gte=0"`
	ImageURL    string   `json:"image"`
	Status      string   `json:"status" binding:"omitempty,oneof=ACTIVE HIDDEN"`
}
```

**Response Payload (Success - 201 Created):**
```go
type ProductResponse struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	CategoryID  uint     `json:"categoryId"`
	Category    *CategoryResponse `json:"category,omitempty"`
	Price       float64  `json:"price"`
	SalePrice   *float64 `json:"salePrice"`
	SKU         string   `json:"sku"`
	Stock       int      `json:"stock"`
	SoldCount   int      `json:"soldCount"`
	ImageURL    string   `json:"image"`
	Status      string   `json:"status"`
	CreatedAt   string   `json:"createdAt"`
}
```

### 2. Các Quy tắc Nghiệp vụ (Business Rules) tại tầng UseCase / Service

Khi tạo mới một product, Service layer (UseCase) cần xử lý các luồng logic (Business validations) sau:
1. **Auto-Generate SKU:** 
   Nếu client không gửi `req.SKU`, hệ thống tự sinh mã SKU ngẫu nhiên kết hợp prefix (VD: `PRD-<UUID_Short>` hoặc logic tuỳ biến).
2. **Kiểm tra trùng lặp SKU:** 
   Thực thi câu truy vấn `Count` để đảm bảo `SKU` (người dùng nhập hoặc tự sinh) là duy nhất. Báo lỗi 400 (Bad Request) nếu trùng lặp.
3. **Kiểm tra Validate Category:** 
   Gọi Repository kiểm tra xem `CategoryID` gửi lên có thật sự tồn tại trong DB không. Báo lỗi 404/400 nếu danh mục không hợp lệ.
4. **Kiểm tra logic Giá khuyến mãi:** 
   Nếu client có truyền lên `SalePrice`, bắt buộc ở tầng Service phải check: `SalePrice < Price`. Mặc dù FE đã check nhưng BE không bao giờ được phép tin tưởng tuyệt đối vào FE.
5. **Gán giá trị mặc định:** 
   Nếu `Status` rỗng, tự động set về `"ACTIVE"`. `SoldCount` khởi tạo bằng `0`.

---

## TRỤ CỘT 3: XỬ LÝ CACHE & QUẢN LÝ TRẠNG THÁI (REDIS INTEGRATION)

### Xóa Cache (Cache Invalidation)
Mỗi khi có một sản phẩm mới được tạo và publish ở trạng thái `"ACTIVE"`, bộ nhớ đệm danh sách sản phẩm ở trang chủ (hoặc theo danh mục) của người dùng cuối có thể bị lỗi thời. Hệ thống cần xóa (Invalidate) các Redis key bị ảnh hưởng.

- **Key Pattern:** `techbite:products:home_list`, `techbite:categories:products:*`.
- **Quy trình xử lý tại UseCase:**
  ```go
  // 1. Lưu sản phẩm vào DB
  err := u.productRepo.Create(ctx, productEntity)
  if err != nil {
      return err
  }
  
  // 2. Xóa cache nếu sản phẩm Đang bán
  if productEntity.Status == "ACTIVE" {
      // Xóa cache danh sách tổng
      u.redisClient.Del(ctx, "techbite:products:home_list")
      // Xóa cache danh sách theo CategoryID
      cacheKey := fmt.Sprintf("techbite:categories:products:%d", productEntity.CategoryID)
      u.redisClient.Del(ctx, cacheKey)
  }
  ```

*(Lưu ý: Luồng Upload Ảnh cho Sản phẩm sẽ sử dụng chung với API Upload có sẵn `POST /api/v1/media/upload`, Client nhận URL và truyền vào field `image` trong JSON Payload).*
