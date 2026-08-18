# KẾ HOẠCH BACKEND: QUẢN TRỊ DANH SÁCH BÀI VIẾT (POST LIST)

## 1. Thiết kế Dữ liệu (Database Schema Struct)

### Tái sử dụng GORM Model: `Post`
Vị trí: `apps/backend/internal/domain/post/entity/post.go` (Đã được tạo từ luồng Post Create).

```go
package entity

import (
	"time"
	"gorm.io/gorm"
)

type Post struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	TypeCode  string         `gorm:"type:varchar(50);not null;index" json:"typeCode"`
	Title     string         `gorm:"type:varchar(255);not null" json:"title"`
	Slug      string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Content   string         `gorm:"type:longtext;not null" json:"content"`
	AuthorID  uint           `gorm:"not null;index" json:"authorId"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```

**Ràng buộc & Tối ưu hoá truy vấn:**
- **Index `TypeCode`:** Được sử dụng để tăng tốc độ phân loại bài viết theo `NEWS`, `PAGE`... khi Frontend gửi yêu cầu lấy danh sách theo loại.
- **Index `DeletedAt`:** Hỗ trợ tính năng Soft Delete mặc định của GORM, giúp ẩn bài viết đã xóa khỏi danh sách thay vì xóa vĩnh viễn khỏi CSDL, đảm bảo toàn vẹn dữ liệu cho Dashboard.
- **Tìm kiếm `Title`:** Việc tìm kiếm (LIKE) trên trường `title` phụ thuộc vào tốc độ quét, nếu dữ liệu rất lớn có thể cân nhắc thêm Full-Text Index trong tương lai, hiện tại dùng thao tác quét tiêu chuẩn kết hợp phân trang là đủ cho admin.

---

## 2. Giao kèo API (API Contract & Context Auth)

### 2.1. Lấy danh sách (List API)
- **Method & Route:** `GET /api/v1/posts`
- **Xác thực:** Cần đi qua `AuthMiddleware` để verify Token, chặn Blacklist. Admin hoặc tài khoản có quyền truy cập Content Management.

**Request Binding Struct (Query Parameters):**
Vị trí đề xuất: `apps/backend/internal/usecases/post/dto/post_query_dto.go`
```go
package dto

type GetPostsQuery struct {
	Page     int    `form:"current" binding:"omitempty,min=1"`
	PageSize int    `form:"pageSize" binding:"omitempty,min=1"`
	Title    string `form:"title_like" binding:"omitempty"` // Tìm kiếm title chứa chuỗi này
	TypeCode string `form:"typeCode" binding:"omitempty,max=50"` // Lọc theo loại bài viết
    SortBy   string `form:"sortBy" binding:"omitempty"` // Ví dụ: "id", "createdAt"
	Order    string `form:"order" binding:"omitempty,oneof=asc desc"`
}
```
*(Ghi chú: Cấu trúc parameters phải tương thích với chuẩn gọi API ngầm của Refine.js DataProvider).*

**Response Payload:**
```go
package dto

import "time"

type PostResponse struct {
	ID        uint      `json:"id"`
	TypeCode  string    `json:"typeCode"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
    // KHÔNG trả về Content trong List để giảm payload JSON.
	AuthorID  uint      `json:"authorId"`
	CreatedAt time.Time `json:"createdAt"`
}

type PaginatedPostResponse struct {
	Data  []PostResponse `json:"data"`
	Total int64          `json:"total"`
}
```

### 2.2. Xóa bài viết (Delete API)
- **Method & Route:** `DELETE /api/v1/posts/:id`
- **Xác thực:** Đi qua `AuthMiddleware`. Yêu cầu quyền quản trị viên, đọc `userId` từ JWT context để ghi log nếu cần.

**Response Payload:** Trả về chuẩn JSON cấu trúc rỗng `{}`, HTTP Status 200 OK. Hệ thống thực hiện Soft Delete (GORM `Delete`).

---

## 3. Xử lý Cache & Quản lý Trạng thái Token (Redis Integration)

- **Quản lý phiên qua Token:** Tại dòng code đầu tiên của Middleware bảo vệ 2 API trên, hệ thống sử dụng khoá `jti` từ Access Token để truy vấn Redis Server. Nếu token nằm trong Blacklist (do đã nhấn `/logout` hoặc bị ban), request bị từ chối ngay tức khắc với HTTP Status `401 Unauthorized`.
- **Cơ chế Caching dữ liệu (Tùy chọn cho Admin):**
  - Mặc dù trang Admin luôn yêu cầu dữ liệu realtime nhất, nếu lưu lượng truy cập bảng danh sách cao, có thể cache query list trên Redis theo key format: `admin:posts:list:type:{typeCode}:page:{page}:size:{pageSize}:search:{searchKeyword}`.
  - TTL (thời gian sống) của cache này nên ngắn (ví dụ: 1-5 phút).
  - **Cache Invalidation (Nguyên tắc xoá):** Mỗi khi API `POST` (Tạo mới), `PUT` (Cập nhật), hoặc `DELETE` (Xoá) bài viết được gọi thành công, Service liên quan BẮT BUỘC phải thực hiện xoá hàng loạt các key bắt đầu bằng `admin:posts:list:*` thông qua lệnh `SCAN` hoặc lưu quản lý Set Key trên Redis để người dùng thấy kết quả cập nhật mới nhất ở màn hình List.
