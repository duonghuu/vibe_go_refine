# QUY HOẠCH KIẾN TRÚC BACK-END: QUẢN TRỊ POST MEDIA

## 1. Thiết kế Dữ liệu (Database Schema Struct)

### GORM Model: PostMedia

```go
package entity

import (
	"time"
	"gorm.io/gorm"
)

// PostMedia đại diện cho bảng liên kết giữa Post và Media
type PostMedia struct {
	ID         uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID     uint           `gorm:"not null;index" json:"post_id"`
	MediaID    uint           `gorm:"not null;index" json:"media_id"`
	Collection string         `gorm:"type:varchar(50);not null;index" json:"collection"` // thumbnail, gallery, attachment, content
	SortOrder  int            `gorm:"type:int;default:0" json:"sort_order"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	
	// Relationships (để load eager loading với Preload)
	// Post  *Post  `gorm:"foreignKey:PostID" json:"post,omitempty"`
	// Media *Media `gorm:"foreignKey:MediaID" json:"media,omitempty"`
}
```

*   **Ràng buộc & Indexing:**
    *   Sử dụng Composite Index hoặc Single Index trên `PostID` và `Collection` để tối ưu hóa truy vấn lấy danh sách media theo collection của một bài viết cụ thể.
    *   Trường `DeletedAt` đảm bảo xóa mềm (Soft Delete) theo chuẩn của GORM.
    *   Liên kết cascade delete khi xóa bài viết sẽ được xử lý ở tầng Service/UseCase hoặc thông qua khai báo constraint trực tiếp tại MySQL.

---

## 2. Giao kèo API (API Contract & Context Auth)

Dưới đây là các định nghĩa endpoint phục vụ giao tiếp với Refine.js Dashboard. Hầu hết các Endpoint ghi đều yêu cầu bảo mật.

### 2.1. Lấy danh sách Media của bài viết

*   **Method & Route:** `GET /api/v1/posts/:post_id/media`
*   **Query Params:** `collection` (tuỳ chọn - để lọc theo loại)
*   **Auth:** Tùy nghiệp vụ bài viết (Công khai thì không cần Auth, xem dạng Dashboard Admin thì cần).

**Response Payload (Thành công 200 OK):**
```go
// Tương thích với Refine.js Data Provider mong đợi một mảng data và total
type GetPostMediaResponse struct {
	Data  []PostMediaResponse `json:"data"`
	Total int64               `json:"total"`
}

type PostMediaResponse struct {
	ID         uint        `json:"id"`
	PostID     uint        `json:"post_id"`
	MediaID    uint        `json:"media_id"`
	Collection string      `json:"collection"`
	SortOrder  int         `json:"sort_order"`
	Media      interface{} `json:"media,omitempty"` // Thông tin thực tế của Media (URL, Name, MimeType)
}
```

### 2.2. Liên kết Media vào bài viết

*   **Method & Route:** `POST /api/v1/posts/:post_id/media`
*   **Auth:** `RequireAuthMiddleware` (Bắt buộc) - Chỉ Admin/Editor có quyền.

**Request Binding Struct:**
```go
type CreatePostMediaRequest struct {
	MediaID    uint   `json:"media_id" binding:"required"`
	Collection string `json:"collection" binding:"required,oneof=thumbnail gallery attachment content"`
	SortOrder  *int   `json:"sort_order"` // Dùng pointer để phân biệt có truyền hay không, nếu null thì gán mặc định 0
}
```

### 2.3. Đồng bộ (Sync) Media cho một Collection

Thao tác gỡ và thêm mới hàng loạt được tích hợp vào một API duy nhất để dễ dàng sử dụng bên Refine.js.

*   **Method & Route:** `PUT /api/v1/posts/:post_id/media/:collection`
*   **Auth:** `RequireAuthMiddleware` (Bắt buộc).

**Request Binding Struct:**
```go
type SyncPostMediaRequest struct {
	// Danh sách các ID của media sẽ thuộc về bài viết ở collection này
	// Backend sẽ xoá những media không nằm trong list và thêm/cập nhật media trong list
	Media []MediaSyncItem `json:"media" binding:"required,dive"`
}

type MediaSyncItem struct {
	ID        uint `json:"id" binding:"required"`
	SortOrder int  `json:"sort_order"`
}
```

### 2.4. Gỡ liên kết Media khỏi bài viết

*   **Method & Route:** `DELETE /api/v1/posts/:post_id/media/:media_id`
*   **Query Params:** `?collection=gallery` (Bắt buộc nếu cùng 1 media_id gắn cho nhiều collection trong 1 bài viết).
*   **Auth:** `RequireAuthMiddleware` (Bắt buộc).

**Response Payload:** Trả về `204 No Content` nếu gỡ thành công.

---

## 3. Xử lý Cache & Quản lý Trạng thái Token (Redis Integration)

### 3.1. Quản lý trạng thái Token và Phân quyền (Auth & Context)

*   **Xác thực tập trung:** Tất cả các endpoint mang tính chất Write/Update (như `POST`, `PUT`, `DELETE`) của module Media này bắt buộc phải đi qua `AuthMiddleware` của Gin.
*   **Trích xuất danh tính:** Middleware tiến hành parse Access Token (JWT), lấy ra thuộc tính định danh và đưa vào context (`c.Set("userId", claims.UserID)`, `c.Set("role", claims.Role)`). Tầng Handler / UseCase bên dưới tuyệt đối không lấy userId trực tiếp từ Body JSON để tránh bị mạo danh.
*   **Bảo vệ thời gian thực (Blacklist):** Middleware cần thực thi dòng code kiểm tra JTI / raw accessToken hiện tại có nằm trong Blacklist của Redis hay không. Nếu có `c.AbortWithStatusJSON(http.StatusUnauthorized)` ngay lập tức, ngăn ngừa lộ lọt dữ liệu khi account bị thu hồi quyền khẩn cấp.

### 3.2. Caching dữ liệu đọc bằng Redis

Trong các hệ thống content, thông tin Media của bài viết thay đổi rất ít nhưng được gọi (read) vô cùng thường xuyên bởi Client bên ngoài.

*   **Key Pattern:** `cache:post:{post_id}:media:{collection}` (Ví dụ: `cache:post:123:media:gallery`).
*   **Chiến lược lưu trữ (Caching Strategy):**
    *   Sử dụng **Cache-Aside Pattern**. API `GET /api/v1/posts/:post_id/media` đầu tiên sẽ gọi vào Redis (ví dụ: hàm `HGET` hoặc `GET`).
    *   Nếu rỗng (Miss), Backend tiến hành truy vấn MySQL, serialize Struct ra byte JSON và đẩy vào Redis với một TTL hợp lý thông qua hàm `SetEx` (ví dụ 24 tiếng).
*   **Chiến lược xóa Cache (Invalidation):**
    *   Khi quản trị viên thao tác `POST`, `PUT (Sync)`, hoặc `DELETE` tại Refine Dashboard, tầng UseCase/Service xử lý lưu Database thành công BẮT BUỘC gọi lệnh xóa (Invalidate) Redis Key tương ứng với `post_id` đó (`client.Del(ctx, cacheKey)`).
    *   Hành động này đảm bảo khi người dùng tải lại trang web, nội dung gallery hoặc thumbnail sẽ luôn là phiên bản mới nhất ngay sau khi Admin bấm lưu.
