# QUY HOẠCH KIẾN TRÚC BACK-END: QUẢN TRỊ DANH MỤC SẢN PHẨM

**Dự án:** TechBite
**Module:** Category Management
**Tài liệu tham chiếu:** `.docs/ideas/dashboard/02-category-list-idea.md`

*(Ghi chú: Theo yêu cầu, tạm thời bỏ qua xác thực người dùng trong bản thiết kế này, tập trung vào thiết kế dữ liệu, API cho nghiệp vụ quản lý danh mục và caching).*

---

## TRỤ CỘT 1: THIẾT KẾ DỮ LIỆU (DATABASE SCHEMA STRUCT)

Sử dụng GORM để định nghĩa Entity `Category`. Cấu trúc này hỗ trợ dạng cây (Parent-Child) thông qua trường `ParentID` tham chiếu đến chính bảng `Category`.

```go
package entity

import (
	"time"
	"gorm.io/gorm"
)

// Category đại diện cho bảng danh mục sản phẩm trong cơ sở dữ liệu
type Category struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Slug        string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	ParentID    *uint          `gorm:"index" json:"parentId"` // Nullable để hỗ trợ danh mục gốc
	Description string         `gorm:"type:text" json:"description"`
	ImageURL    string         `gorm:"type:varchar(500)" json:"imageUrl"`
	SortOrder   int            `gorm:"type:int;default:0;index" json:"sortOrder"`
	Status      string         `gorm:"type:varchar(20);default:'ACTIVE';index" json:"status"` // ACTIVE, HIDDEN
	
	// Quan hệ nội bộ (Self-referencing)
	Parent      *Category      `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children    []Category     `gorm:"foreignKey:ParentID" json:"children,omitempty"`

	// Audit logs & Soft Delete
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
```

### Ràng buộc & Tối ưu (Constraints & Optimization):
- **Index:** Đánh index ở `Slug` (dùng để tìm kiếm nhanh), `ParentID` (phục vụ truy vấn cây danh mục), `SortOrder` (khi cần sắp xếp hiển thị) và `Status`.
- **Soft Delete:** Tích hợp `gorm.DeletedAt` để tránh mất dữ liệu nhầm lẫn và dễ dàng rollback.
- **Ràng buộc Xóa (Cascade/Protect):** Ở tầng Service sẽ chặn xóa nếu danh mục có sản phẩm liên kết (Dựa theo Idea). Không cấu hình Cascade Delete ở Database để an toàn dữ liệu.

---

## TRỤ CỘT 2: GIAO KÈO API (API CONTRACT & CONTEXT AUTH)

Quy định cấu trúc trả về chuẩn theo Refine.js Data Provider: Danh sách phải có `data` và `total`.

### 1. Lấy danh sách danh mục (Phân trang, Tìm kiếm, Lọc)
- **Method & Route:** `GET /api/v1/categories`
- **Query Parameters:**
  - `_start`, `_end` (cho phân trang)
  - `_sort`, `_order` (sắp xếp)
  - `q` (tìm kiếm theo Name, Slug)
  - `status`, `parentId` (bộ lọc)
- **Response Payload (Success - 200 OK):**
```json
{
  "data": [
    {
      "id": 1,
      "name": "Đồ uống",
      "slug": "do-uong",
      "parentId": null,
      "imageUrl": "...",
      "sortOrder": 1,
      "status": "ACTIVE",
      "productCount": 18, 
      "createdAt": "..."
    }
  ],
  "total": 1
}
```
*(Lưu ý: `productCount` được tính toán thông qua Aggregation Query ở tầng Repository hoặc gộp kết quả ở tầng UseCase).*

### 2. Lấy chi tiết danh mục
- **Method & Route:** `GET /api/v1/categories/:id`
- **Response Payload (Success - 200 OK):** Trả về object Category (giống object trong mảng `data` ở trên).

### 3. Tạo mới danh mục
- **Method & Route:** `POST /api/v1/categories`
- **Request Binding Struct:**
```go
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,max=255"`
	Slug        string `json:"slug" binding:"required,max=255"` // Có thể auto-generate từ Name nếu frontend không gửi
	ParentID    *uint  `json:"parentId"`
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
	SortOrder   int    `json:"sortOrder"`
	Status      string `json:"status" binding:"omitempty,oneof=ACTIVE HIDDEN"`
}
```

### 4. Cập nhật danh mục
- **Method & Route:** `PUT /api/v1/categories/:id` (hoặc `PATCH`)
- **Request Binding Struct:**
```go
type UpdateCategoryRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=255"`
	Slug        *string `json:"slug" binding:"omitempty,max=255"`
	ParentID    *uint   `json:"parentId"`
	Description *string `json:"description"`
	ImageURL    *string `json:"imageUrl"`
	SortOrder   *int    `json:"sortOrder"`
	Status      *string `json:"status" binding:"omitempty,oneof=ACTIVE HIDDEN"`
}
```

### 5. Xóa danh mục
- **Method & Route:** `DELETE /api/v1/categories/:id`
- **Logic xử lý:** Kiểm tra xem danh mục có đang chứa sản phẩm hay không (Count product where category_id = id). Nếu > 0, trả về lỗi 400 Bad Request yêu cầu dọn dẹp sản phẩm trước. Nếu là danh mục cha và có danh mục con, cảnh báo hoặc đưa danh mục con lên làm gốc (Tùy logic cụ thể sẽ handle ở Service).

### 6. Cập nhật trạng thái hàng loạt (Bulk Status)
- **Method & Route:** `PATCH /api/v1/categories/bulk-status`
- **Request Binding Struct:**
```go
type BulkUpdateStatusRequest struct {
	IDs    []uint `json:"ids" binding:"required,min=1"`
	Status string `json:"status" binding:"required,oneof=ACTIVE HIDDEN"`
}
```

### 7. Kéo thả đổi vị trí (Reorder)
- **Method & Route:** `PATCH /api/v1/categories/reorder`
- **Request Binding Struct:**
```go
type ReorderItem struct {
	ID        uint `json:"id" binding:"required"`
	SortOrder int  `json:"sortOrder" binding:"required"`
}

type ReorderCategoriesRequest struct {
	Items []ReorderItem `json:"items" binding:"required,min=1"`
}
```

---

## TRỤ CỘT 3: XỬ LÝ CACHE (REDIS INTEGRATION)

Dữ liệu danh mục thay đổi rất ít (chỉ thay đổi khi Admin cấu hình) nhưng lại được truy vấn cực kỳ nhiều (hiển thị menu cho tất cả khách hàng). Do đó, cần cấu hình Cache với Redis để giảm tải cho DB.

### 1. Chiến lược Cache cho Menu khách hàng
- **Key Pattern:** `techbite:categories:active_tree`
- **Data:** Lưu trữ cây danh mục (chỉ chứa các danh mục có `status='ACTIVE'`) dưới dạng JSON String.
- **TTL (Time to Live):** 24 giờ (hoặc không giới hạn).

### 2. Chiến lược Cache Invalidation (Xóa cache)
Bất cứ khi nào có hành động thay đổi dữ liệu bảng Category từ Admin, hệ thống BẮT BUỘC phải thực hiện xóa cache:
- Create (`POST /api/v1/categories`)
- Update (`PUT /api/v1/categories/:id`)
- Delete (`DELETE /api/v1/categories/:id`)
- Bulk Status (`PATCH /api/v1/categories/bulk-status`)
- Reorder (`PATCH /api/v1/categories/reorder`)

Hành động ở tầng Service (Golang) sau khi commit Transaction thành công với GORM:
```go
redisClient.Del(ctx, "techbite:categories:active_tree")
```
Khi request tiếp theo từ khách hàng gọi tới, API sẽ kiểm tra Redis (Cache Miss), truy vấn lại từ MySQL và ghi đè lại vào Redis (Cache Hit cho các lượt sau).

---
*Ghi chú: Vấn đề bảo mật API cho Admin (Token, RBAC, Redis Blacklist) sẽ được tái áp dụng dựa theo các module nền tảng hiện có của dự án.*
