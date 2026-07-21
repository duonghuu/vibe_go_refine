# QUY HOẠCH KIẾN TRÚC BACK-END: THÊM DANH MỤC SẢN PHẨM

**Dự án:** TechBite
**Module:** Category Management
**Tài liệu tham chiếu:** `.docs/ideas/dashboard/06-category-create-idea.md`

*(Ghi chú: Bản thiết kế này tập trung vào luồng tạo mới dữ liệu (Create) cho danh mục sản phẩm, đồng bộ với kiến trúc phân tầng DDD/Clean Architecture đang có).*

---

## TRỤ CỘT 1: THIẾT KẾ DỮ LIỆU (DATABASE SCHEMA STRUCT)

Thiết kế này sử dụng GORM Entity `Category` hỗ trợ cấu trúc cây (Parent-Child) cho phép danh mục chứa danh mục con, và hỗ trợ Soft Delete để tránh rủi ro mất dữ liệu.

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
	
	// Quan hệ nội bộ (Self-referencing) để tạo dạng cây
	Parent      *Category      `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children    []Category     `gorm:"foreignKey:ParentID" json:"children,omitempty"`

	// Audit logs & Soft Delete
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Category) TableName() string {
	return "categories"
}
```

### Ràng buộc & Tối ưu (Constraints & Optimization):
- **Khóa chính & Khóa ngoại:** `ID` tự tăng làm khóa chính. `ParentID` tham chiếu đến `ID` cùng bảng.
- **Index:** `uniqueIndex` ở `Slug` để đảm bảo không bị trùng lặp URL và query nhanh. `index` cho `ParentID`, `SortOrder`, và `Status`.
- **Soft Delete:** Có `DeletedAt`.

---

## TRỤ CỘT 2: GIAO KÈO API (API CONTRACT & CONTEXT AUTH)

### 1. Tạo mới danh mục
- **Method & Route:** `POST /api/v1/admin/categories`
- **Auth:** Yêu cầu đi qua Middleware Auth (Kiểm tra token, phân quyền Admin/Menu Planner).

**Request Binding Struct:**
```go
package dto

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,max=255"`
	Slug        string `json:"slug" binding:"required,max=255"`
	ParentID    *uint  `json:"parentId"` // Nullable nếu là danh mục gốc
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
	SortOrder   int    `json:"sortOrder" binding:"omitempty,min=0"`
	Status      string `json:"status" binding:"omitempty,oneof=ACTIVE HIDDEN"`
}
```

**Response Payload (Success - 201 Created):**
```go
type CategoryResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	ParentID     *uint  `json:"parentId"`
	Description  string `json:"description"`
	ImageURL     string `json:"imageUrl"`
	SortOrder    int    `json:"sortOrder"`
	Status       string `json:"status"`
	ProductCount int64  `json:"productCount"` // Mặc định là 0 khi vừa tạo mới
	CreatedAt    string `json:"createdAt"`
}
```

### 2. Các Quy tắc Nghiệp vụ (Business Rules) tại tầng UseCase / Service

Khi tạo mới một category, Service layer (UseCase) cần xử lý các logic phòng thủ (Defensive logic) sau:
1. **Kiểm tra trùng lặp Slug:**
   Thực thi câu truy vấn `Count` để đảm bảo `req.Slug` chưa tồn tại trong cơ sở dữ liệu. Nếu tồn tại, ném lỗi 400/409 để chặn quá trình tạo mới.
2. **Khởi tạo trạng thái mặc định:**
   Nếu Client không gửi lên `Status`, hệ thống tự động gán giá trị mặc định là `"ACTIVE"`.
3. **Cấu trúc vòng lặp (Circular Dependency):**
   Mặc dù rủi ro tạo vòng lặp khi "Tạo mới" là cực thấp (vì bản ghi chưa có ID), nhưng vẫn cần kiểm tra xem `ParentID` truyền lên có hợp lệ và tồn tại hay không.

---

## TRỤ CỘT 3: XỬ LÝ CACHE & QUẢN LÝ TRẠNG THÁI (REDIS INTEGRATION)

### Xóa Cache (Cache Invalidation)
Mỗi khi có một danh mục mới được tạo ra thành công, bộ nhớ đệm lưu trữ danh sách cây danh mục (thường phục vụ cho Frontend Store/Menu khách hàng) sẽ bị sai lệch. 
Hệ thống **BẮT BUỘC** phải xóa (invalidate) cache này.

- **Key Pattern:** `techbite:categories:active_tree` (và các keys danh sách khác nếu có).
- **Quy trình xử lý tại UseCase:**
  ```go
  // Thực hiện lưu category vào MySQL
  err := u.categoryRepo.Create(ctx, categoryEntity)
  if err != nil {
      return err
  }
  
  // Xóa cache để ép buộc truy vấn lại từ DB ở lượt gọi tiếp theo
  u.redisClient.Del(ctx, "techbite:categories:active_tree")
  ```

*(Lưu ý: Luồng Upload Ảnh sẽ được xử lý riêng rẽ thông qua endpoint `POST /api/v1/media/upload`, Frontend tự động gọi API đó trước để lấy URL rồi mới truyền vào payload của API Create Category).*
