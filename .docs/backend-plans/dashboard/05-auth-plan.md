# QUY HOẠCH KIẾN TRÚC BACK-END: XÁC THỰC NGƯỜI DÙNG (JWT AUTHENTICATION)

**Dự án:** TechBite
**Module:** Authentication (Auth)
**Tài liệu tham chiếu:** `.docs/ideas/dashboard/05-auth-idea.md` & `.docs/ARCHITECTURE.md`

---

## TRỤ CỘT 1: THIẾT KẾ DỮ LIỆU (DATABASE SCHEMA STRUCT)

Tuân thủ nguyên tắc từ file Kiến trúc (`ARCHITECTURE.md`) và quy định của dự án, Module Authentication **KHÔNG TẠO BẢNG MySQL MỚI** để lưu Refresh Token (Session). Việc quản lý phiên người dùng được đẩy hoàn toàn lên **Redis** để đảm bảo khả năng mở rộng và thao tác thu hồi quyền truy cập (Revoke) theo thời gian thực một cách nguyên tử.

Thực thể `User` để truy vấn thông tin sẽ sử dụng trực tiếp từ User Module. Dưới đây là các Struct/DTO binding ở tầng Application/Domain của Auth Module:

```go
package dto

// LoginRequest binding JSON từ client gửi lên
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// UserProfile chứa thông tin công khai trả về cho User
type UserProfile struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

// TokenResponse payload trả về cho client
type TokenResponse struct {
	AccessToken string      `json:"accessToken"`
	ExpiresIn   int         `json:"expiresIn"`
	User        UserProfile `json:"user"` // Có thể nil nếu dùng api refresh token
}

// JwtCustomClaims định nghĩa Payload trong JWT Token
type JwtCustomClaims struct {
	UserID uint   `json:"userId"`
	Role   string `json:"role"`
	Jti    string `json:"jti"` // JWT ID - Dùng để whitelist/blacklist trên Redis
	// Kế thừa RegisteredClaims từ thư viện golang-jwt/jwt
}
```

### Ràng buộc & Bảo mật:
- **BCrypt:** Mật khẩu truyền lên trong `LoginRequest` sẽ được đối chiếu với `password` dạng hash trong DB thông qua `golang.org/x/crypto/bcrypt`.
- **JWT Payload:** Chứa tối thiểu `userId`, `role` và `jti` (JWT ID). Cấm tuyệt đối đính kèm các thông tin nhạy cảm khác.

---

## TRỤ CỘT 2: GIAO KÈO API (API CONTRACT & CONTEXT AUTH)

Tất cả các API dưới đây nằm trong nhóm `/api/v1/auth`. Việc trả về `RefreshToken` BẮT BUỘC được set thông qua **Cookie HTTP-Only**.

### 1. API Đăng nhập
- **Method & Route:** `POST /api/v1/auth/login`
- **Request Binding:** JSON binding vào struct `LoginRequest`.
- **Response Payload:** Trả về struct `TokenResponse` (gồm `accessToken`, `expiresIn`, thông tin `user` rút gọn).
- **Security Action:** 
  - Tạo `Access Token` (Ví dụ: 15-30 phút).
  - Tạo `Refresh Token` chứa `jti` ngẫu nhiên (Ví dụ: 7 ngày).
  - Gắn `Refresh Token` vào Response Header: `Set-Cookie: refreshToken=...; HttpOnly; Secure; SameSite=Strict`.

### 2. API Làm mới Token (Refresh Token Rotation)
- **Method & Route:** `POST /api/v1/auth/refresh-token`
- **Request Binding:** Không nhận body. Backend tự động trích xuất `refreshToken` từ request Cookie.
- **Response Payload:** Trả về struct `TokenResponse` (gồm `accessToken`, `expiresIn`).
- **Security Action:** 
  - Validate chữ ký `refreshToken`.
  - Thực hiện nguyên tử (Atomic) trên Redis: Check `jti` cũ, hủy `jti` cũ, sinh `jti` mới và lưu vào Redis.
  - Gắn `Refresh Token` mới vào Response Header `Set-Cookie`.

### 3. API Lấy thông tin tài khoản hiện tại
- **Method & Route:** `GET /api/v1/auth/me`
- **Request Binding:** Không có. Yêu cầu truyền `Authorization: Bearer <accessToken>`.
- **Middleware:** Bắt buộc đi qua Gin `AuthMiddleware`. Lấy `userId` từ Gin Context (`c.Get("userId")`).
- **Response Payload:** Trả về struct `UserProfile` được map từ DB.

### 4. API Đăng xuất
- **Method & Route:** `POST /api/v1/auth/logout`
- **Request Binding:** Trích xuất `accessToken` từ Header và `refreshToken` từ Cookie.
- **Middleware:** Bắt buộc đi qua Gin `AuthMiddleware`.
- **Response Payload:** Thành công (200 OK) dạng JSON `{ "message": "success" }`.
- **Security Action:** Xóa Cookie `refreshToken` (Set-Cookie với max-age=0), đưa `accessToken` vào Blacklist và xóa `jti` của Refresh Token trên Redis.

---

## TRỤ CỘT 3: XỬ LÝ CACHE & QUẢN LÝ TRẠNG THÁI TOKEN (REDIS INTEGRATION)

Việc quản lý token sử dụng **Redis** để đạt độ trễ thấp và ngăn chặn tức thì (Real-time Revocation) các token không mong muốn.

### 1. Quản lý Refresh Token (Whitelist JTI)
- **Cơ chế:** Lưu định danh `jti` (JWT ID) của Refresh Token. Nếu `jti` không có trong Redis, token đó bị xem là vô hiệu lực (đã bị thu hồi hoặc đã đăng xuất).
- **Key Pattern:** `auth:rf:{userId}:{jti}` -> Trạng thái: `valid`.
- **Lưu trữ (SET):** Sử dụng hàm `client.SetEx(ctx, key, "valid", refreshTokenExpiration)` để Redis tự động xóa rác khi token hết hạn.
- **Refresh Rotation:** Tại endpoint `/refresh-token`, sử dụng kịch bản xử lý tuần tự (hoặc transaction Redis `MULTI/EXEC` nếu cần):
  1. Kiểm tra tồn tại `auth:rf:{userId}:{old_jti}`. Nếu không tồn tại -> Báo lỗi 401 (hoặc cảnh báo Replay Attack, tiến hành xóa tất cả phiên).
  2. Xóa key `old_jti`.
  3. Set key `new_jti` với thời gian sống (TTL) của Refresh Token mới.

### 2. Quản lý Access Token (Blacklist)
- **Cơ chế:** Lưu danh sách các Access Token hoặc `jti` của Access Token ĐÃ BỊ HỦY (khi gọi API `/logout` hoặc chủ động Revoke).
- **Key Pattern:** `auth:bl:at:{accessToken}` (hoặc lấy hash của accessToken để làm key tiết kiệm không gian) -> Giá trị: `revoked`.
- **Lưu trữ (SET):** Tính toán thời gian sống còn lại (TTL = `exp` - `time.Now()`). Lưu vào Redis bằng `client.SetEx(ctx, key, "revoked", ttl)`. Khi hết hạn thực tế, token tự động biến mất khỏi Redis và chữ ký jwt sẽ lo phần `expired` phía sau.
- **Middleware Authorization:** Tại dòng đầu tiên, sau khi parse signature thành công, `AuthMiddleware` thực hiện truy vấn `GET auth:bl:at:{accessToken}`. Nếu có dữ liệu trả về -> `c.AbortWithStatusJSON(http.StatusUnauthorized, ...)` ngay lập tức.

### 3. Tính Nguyên Tử & Tối Ưu
- Khi tài khoản bị Block hoặc thay đổi quyền, Admin có thể gọi logic `KEYS auth:rf:{userId}:*` để lấy danh sách tất cả các thiết bị đăng nhập và xóa nó ngay lập tức (Force Logout all devices).
- Việc quản lý Token như thế này thỏa mãn hoàn toàn bài toán Replay Attack và Memory Leak trên hạ tầng Redis Server.
