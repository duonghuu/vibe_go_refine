# KẾ HOẠCH BACKEND: THÊM MỚI BÀI VIẾT (POST CREATE)

## 1. Thiết kế Dữ liệu (Database Schema Struct)

### GORM Model: `Post`
Vị trí đề xuất: `apps/backend/internal/domain/post/entity/post.go`

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
	AuthorID  uint           `gorm:"not null;index" json:"authorId"` // Parse từ JWT Context
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```

**Ràng buộc & Tối ưu hoá:**
- `TypeCode`: Có thể hiểu là Foreign key tham chiếu logic đến `code` của bảng `post_types`. Đánh index để hỗ trợ truy vấn danh sách theo loại nhanh chóng.
- `Slug`: Đánh `uniqueIndex` bắt buộc vì slug phải là duy nhất trên toàn hệ thống (hoặc ít nhất duy nhất theo `TypeCode`).
- `Content`: Kiểu `longtext` để lưu trữ nội dung bài viết định dạng phong phú (Rich Text/HTML).

---

## 2. Giao kèo API (API Contract & Context Auth)

### Endpoint Tạo mới
- **Method & Route:** `POST /api/v1/posts`
- **Xác thực (Auth & Permission):**
  - Phải đi qua `AuthMiddleware` để verify token hợp lệ và chặn token nằm trong Redis Blacklist.
  - Phân quyền (RBAC): Dành cho Role `ADMIN`, `STAFF` hoặc `EDITOR` tuỳ thiết lập.
  - **Critical:** Lấy `userId` từ Gin Context (`c.Get("userId")`) và gán làm `AuthorID` thay vì lấy từ payload client gửi lên, đảm bảo bảo mật.

### DTO (Request/Response)
Vị trí đề xuất: `apps/backend/internal/usecases/post/dto/post_dto.go`

**Request Payload:**
```go
package dto

type CreatePostRequest struct {
	TypeCode  string `json:"typeCode" binding:"required,max=50"`
	Title     string `json:"title" binding:"required,max=255"`
	Slug      string `json:"slug" binding:"required,max=255"`
	Content   string `json:"content" binding:"required"`
}
```

**Response Payload:**
```go
import "time"

type PostResponse struct {
	ID        uint      `json:"id"`
	TypeCode  string    `json:"typeCode"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Content   string    `json:"content"`
	AuthorID  uint      `json:"authorId"`
	CreatedAt time.Time `json:"createdAt"`
}
```

---

## 3. Xử lý Cache & Quản lý Trạng thái Token (Redis Integration)

- **Redis Cache Invalidation:** Khi tạo bài viết mới thành công, UseCase cần tiến hành invalidate (xoá) bộ nhớ đệm của các API liên quan như `GET /api/v1/posts` (cụ thể là xoá key theo pattern `posts:list:*` hoặc `posts:type:*`) để dữ liệu danh sách bài viết trên CMS (Refine) hoặc Storefront được làm mới ngay lập tức.
- **Quản lý Token (Chuẩn hệ thống):**
  - Mọi request tạo bài viết bắt buộc phải bị từ chối `401 Unauthorized` nếu `accessToken` (hoặc `jti`) của người dùng đang nằm trong Blacklist của Redis (VD: vừa bị logout hoặc thu hồi quyền).
- **Chống Spam (Rate Limiting):** Cân nhắc thiết lập cơ chế giới hạn tần suất tạo bài viết trên Redis (VD: 10 bài / phút / 1 user) để bảo vệ hệ thống khỏi các đợt tấn công nhồi dữ liệu (Flood).
