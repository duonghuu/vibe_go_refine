# QUY HOẠCH KIẾN TRÚC BACK-END: QUẢN TRỊ NGƯỜI DÙNG (USER MANAGEMENT)

**Dự án:** TechBite
**Module:** Admin User Management
**Tài liệu tham chiếu:** `.docs/ideas/dashboard/04-user-list-idea.md`

---

## TRỤ CỘT 1: THIẾT KẾ DỮ LIỆU (DATABASE SCHEMA STRUCT)

Sử dụng GORM để định nghĩa Entity `User`.

```go
package entity

import (
	"time"
	"gorm.io/gorm"
)

// User đại diện cho tài khoản người dùng trong hệ thống
type User struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"` // Không trả về qua JSON
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Role      string         `gorm:"type:varchar(20);default:'CUSTOMER';index" json:"role"` // ADMIN, STAFF, CUSTOMER
	Status    string         `gorm:"type:varchar(20);default:'ACTIVE';index" json:"status"` // ACTIVE, INACTIVE
	LastLogin *time.Time     `gorm:"index" json:"lastLogin"`

	// Audit logs & Soft Delete
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```

### Ràng buộc & Tối ưu (Constraints & Optimization):
- **Index:** `Email` được gán `uniqueIndex` đảm bảo tính duy nhất. Các trường dùng để filter như `Role`, `Status` được đánh `index` để tối ưu truy vấn. Tương tự với `LastLogin` dùng để sort.
- **Bảo mật:** Trường `Password` được gắn thẻ `json:"-"` để hệ thống tự động loại bỏ trường này khi serialize/marshal ra HTTP Response, tránh lộ mật khẩu đã băm (hash).
- **Soft Delete:** Tích hợp `gorm.DeletedAt` đảm bảo an toàn dữ liệu lịch sử.

---

## TRỤ CỘT 2: GIAO KÈO API (API CONTRACT & CONTEXT AUTH)

Tất cả các API dưới đây đều nằm trong nhóm `/api/v1/admin/users` và BẮT BUỘC phải đi qua:
1. `AuthMiddleware`: Xác thực JWT Access Token hợp lệ.
2. `RoleMiddleware("ADMIN", "STAFF")`: Giới hạn quyền truy cập. Một số thao tác nhạy cảm (Tạo, Sửa, Đổi mật khẩu) yêu cầu `RoleMiddleware("ADMIN")`.

### 1. Lấy danh sách người dùng (Phân trang, Tìm kiếm, Lọc)
- **Method & Route:** `GET /api/v1/admin/users`
- **Middleware:** `AuthMiddleware`, `RoleMiddleware("ADMIN", "STAFF")`
- **Query Parameters:** `_start`, `_end`, `_sort`, `_order`, `q` (tìm theo Name, Email), `role`, `status`.
- **Response Payload (200 OK):** Trả về `{ "data": [User...], "total": X }`.

### 2. Xem chi tiết người dùng
- **Method & Route:** `GET /api/v1/admin/users/:id`
- **Middleware:** `AuthMiddleware`, `RoleMiddleware("ADMIN", "STAFF")`
- **Response Payload:** Object `User`.

### 3. Tạo mới người dùng
- **Method & Route:** `POST /api/v1/admin/users`
- **Middleware:** `AuthMiddleware`, `RoleMiddleware("ADMIN")`
- **Request Binding Struct:**
```go
type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=8,max=255"`
	Name     string `json:"name" binding:"required,max=255"`
	Role     string `json:"role" binding:"required,oneof=ADMIN STAFF CUSTOMER"`
	Status   string `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
}
```

### 4. Cập nhật thông tin cơ bản người dùng (Không đổi Email/Password)
- **Method & Route:** `PUT /api/v1/admin/users/:id`
- **Middleware:** `AuthMiddleware`, `RoleMiddleware("ADMIN")`
- **Business Rule:** Không cho cập nhật `Email` để bảo toàn định danh.
- **Request Binding Struct:**
```go
type UpdateUserRequest struct {
	Name   *string `json:"name" binding:"omitempty,max=255"`
	Role   *string `json:"role" binding:"omitempty,oneof=ADMIN STAFF CUSTOMER"`
	Status *string `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
}
```

### 5. Cập nhật trạng thái (Enable/Disable)
- **Method & Route:** `PATCH /api/v1/admin/users/:id/status`
- **Middleware:** `AuthMiddleware`, `RoleMiddleware("ADMIN")`
- **Business Rule:** Từ chối thao tác nếu `id` trùng với `userId` trong Context (tự khóa chính mình).
- **Request Binding Struct:**
```go
type UpdateUserStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=ACTIVE INACTIVE"`
}
```

### 6. Cập nhật trạng thái hàng loạt (Bulk Status)
- **Method & Route:** `PATCH /api/v1/admin/users/bulk-status`
- **Middleware:** `AuthMiddleware`, `RoleMiddleware("ADMIN")`
- **Request Binding Struct:**
```go
type BulkUpdateUserStatusRequest struct {
	IDs    []uint `json:"ids" binding:"required,min=1"`
	Status string `json:"status" binding:"required,oneof=ACTIVE INACTIVE"`
}
```

### 7. Reset Password
- **Method & Route:** `PATCH /api/v1/admin/users/:id/reset-password`
- **Middleware:** `AuthMiddleware`, `RoleMiddleware("ADMIN")`
- **Request Binding Struct:**
```go
type ResetPasswordRequest struct {
	NewPassword string `json:"newPassword" binding:"required,min=8,max=255"`
}
```

---

## TRỤ CỘT 3: XỬ LÝ CACHE & QUẢN LÝ TRẠNG THÁI TOKEN (REDIS INTEGRATION)

Module quản lý User có tính nhạy cảm cao về bảo mật. Cần áp dụng nghiêm ngặt các quy trình xử lý Token với Redis:

### 1. Luồng xử lý khi Disable tài khoản (INACTIVE)
Khi Admin gọi API `PATCH .../status` và set user thành `INACTIVE` (hoặc Reset Password):
- **Bắt buộc Invalidates Session:** Hệ thống phải truy cập Redis và **XÓA TOÀN BỘ** `jti` (JWT ID của Refresh Token) thuộc về `userID` này đang được lưu trên Redis. (Pattern: `auth:refresh_token:<userID>:*`).
- Điều này ép user trên tất cả các thiết bị không thể lấy Access Token mới và sẽ bị đẩy ra ngoài sau khi Access Token hiện tại hết hạn (hoặc đưa luôn Access Token vào Blacklist nếu hệ thống có lưu track Access Token tập trung).

### 2. Luồng bảo vệ Access Token (Redis Blacklist)
- **Endpoint Logout (`POST /logout`):** Đọc JWT Access Token, tính toán khoảng thời gian còn lại (TTL) và lưu JWT Signature vào Redis bằng `client.SetEx(ctx, "blacklist:"+tokenSig, "revoked", ttl)`.
- **Auth Middleware:** Trước khi xác thực Role, Middleware phải verify JWT Signature, sau đó kiểm tra xem token này có nằm trong Blacklist của Redis hay không. Nếu có `-> 401 Unauthorized` ngay lập tức.
- **Nguyên lý thu gom rác:** Do sử dụng `SetEx` với TTL khớp đúng với thời gian hết hạn của JWT, Redis sẽ tự động dọn dẹp các token bị Blacklist khi chúng hết hạn thật sự, tránh gây phình bộ nhớ RAM.

### 3. Business Logic Level
- Quá trình thực hiện `Reset Password` cần chạy nguyên tử trong 1 Unit of Work. Băm mật khẩu (Bcrypt) -> Update DB -> Invalidate các session hiện tại của User trong Redis. 
- Mật khẩu sử dụng cost mặc định (e.g. 10 hoặc 12) của gói `golang.org/x/crypto/bcrypt`.
