# QUY HOẠCH KIẾN TRÚC BACK-END: CHỈNH SỬA DANH MỤC BÀI VIẾT

**Dự án:** TechBite
**Module:** Post Category Management
**Tài liệu tham chiếu:** `.docs/ideas/dashboard/18-post-category-edit-idea.md`

*(Ghi chú: Bản thiết kế này tập trung vào luồng lấy chi tiết (Get By ID) và cập nhật (Update) cho danh mục bài viết. Đồng bộ với kiến trúc DDD/Clean Architecture đang có).*

---

## TRỤ CỘT 1: THIẾT KẾ DỮ LIỆU (DATABASE SCHEMA STRUCT)

Sử dụng lại GORM Entity `PostCategory` đã định nghĩa ở phần tạo mới. Không có thay đổi về cấu trúc bảng cho tác vụ Edit, tuy nhiên cần chú ý đến quy tắc không cập nhật trường `TypeCode`.

```go
package entity

import (
	"time"
	"gorm.io/gorm"
)

// (Entity tái sử dụng từ quá trình Create)
type PostCategory struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	TypeCode    string         `gorm:"type:varchar(50);not null;index" json:"typeCode"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Slug        string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	ParentID    *uint          `gorm:"index" json:"parentId"`
	Description string         `gorm:"type:text" json:"description"`
	ImageURL    string         `gorm:"type:varchar(500)" json:"imageUrl"`
	SortOrder   int            `gorm:"type:int;default:0;index" json:"sortOrder"`
	Status      string         `gorm:"type:varchar(20);default:'ACTIVE';index" json:"status"`
	
	Parent      *PostCategory  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children    []PostCategory `gorm:"foreignKey:ParentID" json:"children,omitempty"`

	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
```

---

## TRỤ CỘT 2: GIAO KÈO API (API CONTRACT & CONTEXT AUTH)

### 1. Lấy chi tiết danh mục bài viết (Get By ID)
- **Method & Route:** `GET /api/v1/admin/post-categories/:id`
- **Auth:** Yêu cầu đi qua Middleware Auth.
- **Mục đích:** Fetch dữ liệu chi tiết của danh mục để fill vào form chỉnh sửa ở Frontend.

**Response Payload (Success - 200 OK):**
```go
type PostCategoryResponse struct {
	ID          uint   `json:"id"`
	TypeCode    string `json:"typeCode"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	ParentID    *uint  `json:"parentId"`
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
	SortOrder   int    `json:"sortOrder"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
}
```

### 2. Cập nhật danh mục bài viết
- **Method & Route:** `PUT /api/v1/admin/post-categories/:id`
- **Auth:** Yêu cầu Middleware Auth.

**Request Binding Struct:**
Lưu ý: Không cho phép cập nhật `TypeCode` qua API này để đảm bảo tính nhất quán dữ liệu của các bài viết đã gắn với danh mục.

```go
package dto

type UpdatePostCategoryRequest struct {
	Name        string `json:"name" binding:"required,max=255"`
	Slug        string `json:"slug" binding:"required,max=255"`
	ParentID    *uint  `json:"parentId"` // Nullable
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
	SortOrder   int    `json:"sortOrder" binding:"omitempty,min=0"`
	Status      string `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
}
```

**Response Payload (Success - 200 OK):** Trả về `PostCategoryResponse` đã được cập nhật.

### 3. Các Quy tắc Nghiệp vụ (Business Rules) tại UseCase
1. **Kiểm tra tồn tại:** Truy vấn lấy entity cũ trong DB dựa trên `:id`. Nếu không thấy, trả về `404 Not Found`.
2. **Kiểm tra trùng lặp Slug:** `Count` kiểm tra xem `req.Slug` có bị trùng với một danh mục nào khác (ngoại trừ chính ID hiện tại) không. Nếu trùng, trả về `409 Conflict`.
3. **Ràng buộc Danh mục cha (Circular Reference & Cross-Type):**
   - Không được gán `ParentID` bằng chính ID đang cập nhật.
   - Không được gán `ParentID` là ID của một trong những node con/cháu của danh mục đang cập nhật (Circular Reference Validation).
   - Truy vấn thông tin `ParentCategory`, đảm bảo `ParentCategory.TypeCode` khớp với `TypeCode` gốc của danh mục hiện tại.

---

## TRỤ CỘT 3: XỬ LÝ CACHE & QUẢN LÝ TRẠNG THÁI (REDIS INTEGRATION)

### Xóa Cache Danh mục
Do một danh mục bài viết thay đổi thông tin (tên, parent, thứ tự), nó sẽ ảnh hưởng đến kết quả của API lấy danh sách cây (Tree). Do đó cần invalidate cache.

- **Key Pattern:** `techbite:post_categories:tree:{typeCode}`
- **Quy trình xử lý tại UseCase khi Update thành công:**
  ```go
  // Lấy entity cũ để biết TypeCode
  oldEntity, err := u.postCategoryRepo.FindByID(ctx, id)
  
  // Lưu thay đổi vào MySQL...
  err = u.postCategoryRepo.Update(ctx, updatedEntity)
  
  // Xóa cache của ĐÚNG loại bài viết bị ảnh hưởng
  cacheKey := fmt.Sprintf("techbite:post_categories:tree:%s", oldEntity.TypeCode)
  u.redisClient.Del(ctx, cacheKey)
  ```
