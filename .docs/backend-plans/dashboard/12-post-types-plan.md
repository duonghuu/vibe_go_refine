# QUY HOẠCH KIẾN TRÚC BACK-END: QUẢN LÝ LOẠI BÀI VIẾT (POST TYPES)

**Dự án:** TechBite
**Module:** Post Type Management
**Tài liệu tham chiếu:** `.docs/ideas/dashboard/12-post-types-create-idea.md`

*(Ghi chú: Bản thiết kế này định nghĩa kiến trúc cho tính năng Quản lý loại bài viết (Post Types) bao gồm các API Create, Read, Update, Delete, đồng bộ với kiến trúc phân tầng DDD/Clean Architecture).*

---

## TRỤ CỘT 1: THIẾT KẾ DỮ LIỆU (DATABASE SCHEMA STRUCT)

Thiết kế GORM Entity `PostType` với cấu trúc đơn giản để quản lý danh mục bài viết. Sử dụng Soft Delete (`DeletedAt`) để ngăn chặn xóa mất dữ liệu ngoài ý muốn.

```go
package entity

import (
	"time"
	"gorm.io/gorm"
)

// PostType đại diện cho bảng loại bài viết trong cơ sở dữ liệu
type PostType struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Code      string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Status    string         `gorm:"type:varchar(20);default:'ACTIVE';index" json:"status"` // ACTIVE, INACTIVE
	SortOrder int            `gorm:"type:int;default:0;index" json:"sort_order"`

	// Audit logs & Soft Delete
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (PostType) TableName() string {
	return "post_types"
}
```

### Ràng buộc & Tối ưu (Constraints & Optimization):
- **Khóa chính:** `ID` tự tăng làm khóa chính.
- **Index:** `uniqueIndex` ở `Code` để đảm bảo mã loại không trùng lặp và tìm kiếm nhanh. `index` cho `Status` và `SortOrder` để tối ưu câu lệnh WHERE và ORDER BY.
- **Soft Delete:** Hỗ trợ `DeletedAt` tự động của GORM.

---

## TRỤ CỘT 2: GIAO KÈO API (API CONTRACT & CONTEXT AUTH)

Tất cả các API dưới đây **BẮT BUỘC** phải đi qua Middleware Auth của Gin để xác thực Access Token từ Redis (đảm bảo chỉ Admin/Staff được truy cập).

### 1. Lấy danh sách loại bài viết
- **Method & Route:** `GET /api/v1/admin/post-types`
- **Chức năng:** Trả về danh sách kèm phân trang, hỗ trợ lọc theo status, tìm kiếm theo code/name.

### 2. Tạo mới loại bài viết
- **Method & Route:** `POST /api/v1/admin/post-types`

**Request Binding Struct:**
```go
package dto

type CreatePostTypeRequest struct {
	Code      string `json:"code" binding:"required,max=50,alphanum"` // Cần custom validator cho format A-Z0-9_ nếu cần
	Name      string `json:"name" binding:"required,max=100"`
	Status    string `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
	SortOrder int    `json:"sort_order" binding:"omitempty,min=1"`
}
```

### 3. Cập nhật loại bài viết
- **Method & Route:** `PUT /api/v1/admin/post-types/:id`

**Request Binding Struct:**
```go
package dto

type UpdatePostTypeRequest struct {
	Name      string `json:"name" binding:"required,max=100"`
	Status    string `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
	SortOrder int    `json:"sort_order" binding:"omitempty,min=1"`
}
// Lưu ý: Trường Code thường không cho phép cập nhật để tránh ảnh hưởng logic phân loại (hoặc nếu có, phải validate unique).
```

### 4. Xóa loại bài viết
- **Method & Route:** `DELETE /api/v1/admin/post-types/:id`

**Response Payload (Chung cho các API thành công - ngoại trừ Delete/List):**
```go
type PostTypeResponse struct {
	ID        uint   `json:"id"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	SortOrder int    `json:"sort_order"`
}
```

### Quy tắc Nghiệp vụ (Business Rules) tại tầng UseCase
1. **Kiểm tra trùng lặp Code:** Khi POST, thực thi truy vấn để đảm bảo `req.Code` chưa tồn tại trong cơ sở dữ liệu. Ném lỗi 409 Conflict nếu đã tồn tại.
2. **Khởi tạo trạng thái mặc định:** Nếu không truyền lên `Status`, gán mặc định `"ACTIVE"`.
3. **Ràng buộc Xóa (Delete Constraint):** Không được xóa `post_type` nếu đang có Bài viết (Post) nào sử dụng (Foreign Key Constraint). Bắt buộc phải query đếm số lượng bài viết liên kết trước khi thực hiện xóa. Nếu có, ném lỗi 400 Bad Request kèm thông báo "Không thể xóa loại bài viết đang được sử dụng".

---

## TRỤ CỘT 3: XỬ LÝ CACHE & QUẢN LÝ TRẠNG THÁI (REDIS INTEGRATION)

### Xóa Cache (Cache Invalidation)
Mỗi khi có hành động `POST`, `PUT`, `DELETE` đối với PostType, hệ thống cần xóa cache danh sách loại bài viết (nếu có được lưu trong Redis để hiển thị menu hoặc bộ lọc cho ứng dụng Frontend/Khách hàng).

- **Key Pattern:** `techbite:post_types:active_list`
- **Quy trình xử lý tại UseCase:**
  ```go
  // Thực hiện lưu/cập nhật/xóa post_type trong MySQL
  // ...
  
  // Xóa cache
  u.redisClient.Del(ctx, "techbite:post_types:active_list")
  ```

*(Lưu ý: Xác thực tài khoản Admin/Staff được kế thừa hoàn toàn từ cơ chế AuthMiddleware có sẵn kiểm tra Access Token qua Redis, không cần cấu hình mới).*
