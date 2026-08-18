# QUY HOẠCH KIẾN TRÚC BACK-END: THÊM DANH MỤC BÀI VIẾT

**Dự án:** TechBite
**Module:** Post Category Management
**Tài liệu tham chiếu:** `.docs/ideas/dashboard/17-post-category-create-idea.md`

*(Ghi chú: Bản thiết kế này tập trung vào luồng tạo mới (Create) và lấy danh sách dạng cây (Tree) cho danh mục bài viết theo loại bài viết - Post Type. Đồng bộ với kiến trúc DDD/Clean Architecture đang có).*

---

## TRỤ CỘT 1: THIẾT KẾ DỮ LIỆU (DATABASE SCHEMA STRUCT)

Sử dụng GORM Entity `PostCategory`. Khác với danh mục sản phẩm, danh mục bài viết được gộp nhóm dựa trên `TypeCode` (Loại bài viết), đồng thời vẫn hỗ trợ cấu trúc cây (Parent-Child) và Soft Delete.

```go
package entity

import (
	"time"
	"gorm.io/gorm"
)

// PostCategory đại diện cho bảng danh mục bài viết trong cơ sở dữ liệu
type PostCategory struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	TypeCode    string         `gorm:"type:varchar(50);not null;index" json:"typeCode"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Slug        string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	ParentID    *uint          `gorm:"index" json:"parentId"` // Nullable để hỗ trợ danh mục gốc
	Description string         `gorm:"type:text" json:"description"`
	ImageURL    string         `gorm:"type:varchar(500)" json:"imageUrl"`
	SortOrder   int            `gorm:"type:int;default:0;index" json:"sortOrder"`
	Status      string         `gorm:"type:varchar(20);default:'ACTIVE';index" json:"status"` // ACTIVE, INACTIVE
	
	// Quan hệ nội bộ (Self-referencing) để tạo dạng cây
	Parent      *PostCategory  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children    []PostCategory `gorm:"foreignKey:ParentID" json:"children,omitempty"`

	// Audit logs & Soft Delete
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (PostCategory) TableName() string {
	return "post_categories"
}
```

### Ràng buộc & Tối ưu (Constraints & Optimization):
- **Khóa chính & Ngoại:** `ID` tự tăng. `ParentID` tham chiếu cùng bảng.
- **Index:** `uniqueIndex` ở `Slug` đảm bảo tính duy nhất toàn cục. `index` ở `TypeCode`, `ParentID`, `SortOrder`, và `Status` giúp filter nhanh chóng theo loại bài viết.
- **Quan hệ PostType (nếu có Entity PostType độc lập):** `TypeCode` có thể được hiểu là Khóa ngoại lỏng tham chiếu tới bảng `post_types(code)` (Có thể khai báo Foreign Key nếu Database chặt chẽ).

---

## TRỤ CỘT 2: GIAO KÈO API (API CONTRACT & CONTEXT AUTH)

### 1. Lấy danh sách danh mục cha theo Loại bài viết
- **Method & Route:** `GET /api/v1/admin/post-categories/tree?type_code={code}`
- **Auth:** Yêu cầu đi qua Middleware Auth.
- **Mục đích:** Fetch danh sách cây danh mục chỉ thuộc về một `type_code` chỉ định, để cung cấp dữ liệu cho dropdown Parent ID ở Frontend.

**Response Payload (Success - 200 OK):**
```go
// Trả về một mảng phẳng nhưng đã được sắp xếp hoặc mảng lồng nhau tùy chuẩn Data Provider Refine (thường Refine cần mảng phẳng).
// Ở đây trả về mảng phẳng cho options:
type PostCategoryTreeResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"` // Có thể được format thêm prefix "--" theo level
	TypeCode string `json:"typeCode"`
}
// Payload thực tế sẽ được wrap trong cấu trúc chuẩn của Gin: { "data": [...] }
```

### 2. Tạo mới danh mục bài viết
- **Method & Route:** `POST /api/v1/admin/post-categories`
- **Auth:** Yêu cầu Middleware Auth.

**Request Binding Struct:**
```go
package dto

type CreatePostCategoryRequest struct {
	TypeCode    string `json:"typeCode" binding:"required,max=50"`
	Name        string `json:"name" binding:"required,max=255"`
	Slug        string `json:"slug" binding:"required,max=255"`
	ParentID    *uint  `json:"parentId"` // Nullable
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
	SortOrder   int    `json:"sortOrder" binding:"omitempty,min=0"`
	Status      string `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
}
```

**Response Payload (Success - 201 Created):**
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

### 3. Các Quy tắc Nghiệp vụ (Business Rules) tại UseCase
1. **Kiểm tra hợp lệ TypeCode:** Nếu có bảng `post_types`, kiểm tra xem `typeCode` có tồn tại không.
2. **Kiểm tra trùng lặp Slug:** `Count` kiểm tra tính duy nhất của `req.Slug`. Nếu đã tồn tại, ném `409 Conflict`.
3. **Ràng buộc Danh mục cha (Cross-Type Validation):** 
   Nếu `ParentID` được gửi lên, truy vấn lấy thông tin Category cha và kiểm tra: `ParentCategory.TypeCode == req.TypeCode`. CẤM GẮN danh mục cha thuộc Loại A cho danh mục con thuộc Loại B.
4. **Trạng thái mặc định:** Nếu `Status` rỗng, gán `"ACTIVE"`.

---

## TRỤ CỘT 3: XỬ LÝ CACHE & QUẢN LÝ TRẠNG THÁI (REDIS INTEGRATION)

### Xóa Cache theo Phân mảnh Loại bài viết (Cache Invalidation)
Khác với sản phẩm, dữ liệu danh mục bài viết bị chia mảnh dựa vào `typeCode`. Việc lưu cache cũng phải dựa theo Namespace của từng loại bài viết để tối ưu bộ nhớ.

- **Key Pattern:** `techbite:post_categories:tree:{typeCode}`
- **Quy trình xử lý tại UseCase khi Create thành công:**
  ```go
  // Lưu entity vào MySQL
  err := u.postCategoryRepo.Create(ctx, categoryEntity)
  if err != nil {
      return err
  }
  
  // Xóa cache của ĐÚNG loại bài viết bị ảnh hưởng
  cacheKey := fmt.Sprintf("techbite:post_categories:tree:%s", categoryEntity.TypeCode)
  u.redisClient.Del(ctx, cacheKey)
  ```
  *(Điều này đảm bảo khi tạo danh mục Tin tức (NEWS), cache của danh mục Dịch vụ (SERVICES) không bị xóa nhầm, tăng tỉ lệ Cache Hit).*
