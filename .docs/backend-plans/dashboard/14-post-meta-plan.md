# QUY HOẠCH KIẾN TRÚC BACK-END: POST META

## TRỤ CỘT 1: THIẾT KẾ DỮ LIỆU (DATABASE SCHEMA STRUCT)

**Entity: `PostMeta`**

Đây là bảng dùng để lưu trữ metadata mở rộng cho thực thể `Post`. Do tính chất dữ liệu linh hoạt (Key-Value), bảng này sử dụng kiểu lưu trữ chuỗi cho giá trị.

```go
package entity

import (
	"time"
)

type PostMeta struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID    uint      `gorm:"not null;index:idx_post_meta_post_key,unique" json:"post_id"` // Khóa ngoại và unique index kết hợp với key
	Key       string    `gorm:"type:varchar(255);not null;index:idx_post_meta_post_key,unique" json:"key"`
	Value     string    `gorm:"type:text" json:"value"` // Text để có thể chứa nội dung dài hoặc JSON string
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// Có thể bỏ qua Soft Delete (gorm.DeletedAt) đối với bảng meta để tối ưu hiệu năng vì meta liên kết chặt với lifecycle của Post. Nếu Post bị xóa cứng/mềm thì tùy thuộc ứng dụng xử lý cascade.
}
```

**Ràng buộc (Constraints):**
- **Khóa ngoại:** `PostID` tham chiếu đến `id` của bảng `posts` với tùy chọn `ON DELETE CASCADE` (khi cấu hình Migration).
- **Indexing:** `Composite Unique Index` (`idx_post_meta_post_key`) trên 2 cột `post_id` và `key` đảm bảo mỗi bài viết chỉ có duy nhất một giá trị cho một khóa cụ thể, thuận tiện cho việc Upsert.

---

## TRỤ CỘT 2: GIAO KÈO API (API CONTRACT & CONTEXT AUTH)

Các API này phục vụ Dashboard để lấy thông tin meta và đồng bộ meta.

### 1. Lấy danh sách Meta của một Bài viết

- **Method:** `GET`
- **Route:** `/api/v1/posts/:postId/meta`
- **Middleware:** `RequireAuth` (Yêu cầu JWT Token hợp lệ)

**Request Binding:**
```go
type GetPostMetaRequest struct {
	PostID uint `uri:"postId" binding:"required"`
}
```

**Response Payload:**
```go
type PostMetaResponse struct {
	ID        uint      `json:"id"`
	PostID    uint      `json:"post_id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetPostMetaResponse struct {
	Data []PostMetaResponse `json:"data"`
}
```

### 2. Đồng bộ (Sync/Upsert) Meta cho Bài viết

- **Method:** `PUT`
- **Route:** `/api/v1/posts/:postId/meta`
- **Middleware:** `RequireAuth`, `RequireRole(ADMIN, STAFF)` (Chỉ quyền quản trị hoặc nhân viên mới được phép thao tác)

**Request Binding:**
```go
type MetaItemRequest struct {
	Key   string `json:"key" binding:"required,max=255"`
	Value string `json:"value"`
}

type SyncPostMetaRequest struct {
	PostID uint              `uri:"postId" binding:"required"`
	Meta   []MetaItemRequest `json:"meta" binding:"required,dive"`
}
```

**Response Payload:**
```go
type SyncPostMetaResponse struct {
	Message string `json:"message"` // "Đồng bộ meta thành công"
}
```

### 3. Xóa một Meta cụ thể của Bài viết

- **Method:** `DELETE`
- **Route:** `/api/v1/posts/:postId/meta/:key`
- **Middleware:** `RequireAuth`, `RequireRole(ADMIN, STAFF)`

**Request Binding:**
```go
type DeletePostMetaRequest struct {
	PostID uint   `uri:"postId" binding:"required"`
	Key    string `uri:"key" binding:"required"`
}
```

**Response Payload:**
```go
type DeletePostMetaResponse struct {
	Message string `json:"message"` // "Xóa meta thành công"
}
```

---

## TRỤ CỘT 3: XỬ LÝ CACHE & QUẢN LÝ TRẠNG THÁI TOKEN (REDIS INTEGRATION)

### Xử lý Cache (Caching)
- **Key Pattern:** `post:{postId}:meta`
- **Cơ chế lưu trữ:** Sử dụng kiểu dữ liệu `Hash` hoặc lưu toàn bộ mảng JSON String trong Redis thông qua hàm `SETEX` với thời gian sống (TTL) phù hợp (ví dụ: 1-2 giờ) đối với Client App. Với API dành cho Dashboard, có thể bỏ qua Cache để luôn lấy dữ liệu Real-time, hoặc buộc phải Invalidate Cache mỗi khi gọi API `PUT` (Sync) và `DELETE` phía trên.
- **Luồng vô hiệu hóa (Invalidation):** 
  - Tại API Sync (`PUT`) hoặc API Delete (`DELETE`), sau khi cập nhật thành công xuống cơ sở dữ liệu MySQL, Backend BẮT BUỘC thực thi lệnh `client.Del(ctx, "post:{postId}:meta")` để xóa rác khỏi Redis.

### Quản lý Token Auth
- Vì các Endpoints trên đều là Protected Route (được bảo vệ bởi JWT):
  - Middleware sẽ trích xuất token, xác minh JWT Signature.
  - Sau đó kiểm tra chuỗi `jti` hoặc `accessToken` trên **Redis Blacklist**. Nếu tồn tại => Chặn đứng ngay tại Middleware bằng `http.StatusUnauthorized`.
  - Parse các claim (`userId`, `role`) đẩy vào Gin Context (`c.Set("userId", claims.UserID)`).
  - Tầng Handler/UseCase BẮT BUỘC phải lấy `userId` và `role` từ Context này, không được phép đọc từ body request để tránh giả mạo dữ liệu quyền hạn.
