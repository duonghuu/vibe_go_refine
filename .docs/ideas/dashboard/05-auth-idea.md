# Ý TƯỞNG: Xác thực người dùng (JWT Authentication)

## 1. Thông tin chung (Meta Info)

* **Dự án:** TechBite.
* **Tính năng:** Xác thực người dùng (JWT Authentication).
* **Mục đích:** Cung cấp cơ chế xác thực và quản lý phiên đăng nhập cho toàn bộ hệ thống thông qua JWT Access Token và Refresh Token. Module chỉ chịu trách nhiệm xác thực, không quản lý thông tin người dùng (đã được tách thành User Module).

---

# 2. Đối tượng & Trải nghiệm (Target & UX)

* **Người dùng chính:**

  * ADMIN
  * STAFF
  * CUSTOMER

* **Hành động chính:**

  * Đăng nhập bằng Email và Password.
  * Tự động duy trì phiên làm việc bằng Refresh Token.
  * Đăng xuất.
  * Lấy thông tin người dùng đang đăng nhập.
  * Tự động gia hạn Access Token khi hết hạn.
  * Tự động chuyển về màn hình Login khi phiên đăng nhập hết hiệu lực.

* **Cảm xúc mang lại:**

  * Đăng nhập nhanh.
  * Bảo mật.
  * Phiên làm việc ổn định.
  * Người dùng gần như không nhận thấy quá trình Refresh Token.

---

# 3. Đặc tả Thiết kế (Design Specs)

## Phong cách UI

* Giao diện Login đơn giản.
* Responsive.
* Form validation rõ ràng.
* Thông báo lỗi trực quan.
* Loading khi đang xác thực.

---

## Màu sắc chủ đạo (Brand Colors)

* **Cam (Primary):** `bg-[#ff8c42]`
* **Xanh lá (Success):** `bg-green-500`
* **Đỏ (Error):** `bg-red-500`
* **Xanh dương (Info):** `bg-blue-500`

---

## Cấu trúc màn hình (Top to Bottom)

### Login Page

Hiển thị

* Logo hệ thống
* Email
* Password
* Remember Me (tùy chọn)
* Nút Đăng nhập

---

### Loading

Trong quá trình Login

* Disable Button
* Hiển thị Loading Spinner

---

### Thông báo lỗi

Các trường hợp

* Sai Email hoặc Password
* Tài khoản bị khóa (INACTIVE)
* Tài khoản không tồn tại
* Lỗi hệ thống

---

# 4. Dữ liệu cốt lõi (Mock Data)

## Login Request

```json
{
    "email": "admin@techbite.com",
    "password": "12345678"
}
```

---

## Login Response

```json
{
    "accessToken": "eyJhbGciOi...",
    "refreshToken": "4ab7db18...",
    "expiresIn": 1800,
    "user": {
        "id": 1,
        "email": "admin@techbite.com",
        "name": "System Admin",
        "role": "ADMIN"
    }
}
```

---

## JWT Payload

```json
{
    "sub": "1",
    "role": "ADMIN",
    "exp": 1789999999,
    "iat": 1789990000
}
```

Payload chỉ chứa thông tin cần thiết:

* User ID
* Role
* Issued At
* Expired At

---

# 5. Chức năng nổi bật cho Authentication

## Đăng nhập (Login)

Người dùng đăng nhập bằng

* Email
* Password

Quy trình

* Kiểm tra dữ liệu đầu vào.
* Tìm User theo Email.
* Kiểm tra trạng thái tài khoản.
* So sánh Password bằng BCrypt.
* Sinh Access Token.
* Sinh Refresh Token.
* Lưu Refresh Token vào Database.
* Cập nhật `last_login_at`.
* Trả Access Token và Refresh Token.

---

## Đăng xuất (Logout)

Khi Logout

* Thu hồi (Revoke) Refresh Token hiện tại.
* Xóa Token phía Client.
* Kết thúc phiên làm việc.

---

## Refresh Token

Khi Access Token hết hạn

* Client gửi Refresh Token.
* Kiểm tra Refresh Token trong Database.
* Kiểm tra thời hạn sử dụng.
* Thu hồi Refresh Token cũ (Token Rotation).
* Sinh Access Token mới.
* Sinh Refresh Token mới.
* Lưu Refresh Token mới.
* Trả về bộ Token mới.

---

## Lấy thông tin người dùng hiện tại

API

```
GET /auth/me
```

Trả về

* ID
* Email
* Name
* Role

Thông tin được lấy từ User Module dựa trên User ID trong JWT.

---

## Middleware xác thực JWT

Tất cả API yêu cầu đăng nhập đều đi qua Middleware.

Thực hiện

* Đọc Authorization Header.
* Kiểm tra định dạng Bearer Token.
* Verify JWT.
* Kiểm tra Expired.
* Lấy User ID và Role.
* Lưu thông tin vào Gin Context.

Ví dụ

```go
ctx.Set("userID", 1)
ctx.Set("role", "ADMIN")
```

---

## Phân quyền (Authorization)

Module Authentication chỉ xác định danh tính và vai trò.

Các Middleware phân quyền

* RequireAuthenticated()
* RequireRole("ADMIN")
* RequireRole("STAFF")
* RequireRole("CUSTOMER")

Có thể mở rộng sang Permission hoặc RBAC trong tương lai.

---

## Session Management

Refresh Token được lưu trong bảng riêng.

Thông tin

* User ID
* Token Hash
* Expired At
* Revoked At
* Device (khuyến nghị)
* IP Address (khuyến nghị)

Không lưu Refresh Token dạng Plain Text.

---

## Tự động Refresh Token (Frontend)

RefineJS sử dụng Axios Interceptor.

Luồng xử lý

* Gửi Access Token.
* Nếu API trả về `401 Unauthorized` do Access Token hết hạn:

  * Gọi `/auth/refresh`.
  * Lưu Access Token mới.
  * Gửi lại Request ban đầu.
* Nếu Refresh Token cũng hết hạn:

  * Logout.
  * Điều hướng về Login.

---

## Kiểm tra dữ liệu đầu vào

### Login

* Email bắt buộc.
* Email đúng định dạng.
* Password bắt buộc.

---

### Refresh Token

* Refresh Token bắt buộc.
* Refresh Token hợp lệ.
* Chưa hết hạn.
* Chưa bị thu hồi.

---

## Quy tắc nghiệp vụ (Business Rules)

* Chỉ User có trạng thái `ACTIVE` mới được đăng nhập.
* Password luôn được lưu dưới dạng `password_hash` (BCrypt).
* Access Token có thời gian sống ngắn (ví dụ: 30 phút).
* Refresh Token có thời gian sống dài hơn (ví dụ: 7 ngày).
* Refresh Token chỉ được lưu dưới dạng Hash trong Database.
* Mỗi lần Refresh phải sinh Refresh Token mới (Token Rotation).
* Refresh Token cũ phải bị thu hồi sau khi Refresh thành công.
* `last_login_at` chỉ được cập nhật sau khi đăng nhập thành công.
* JWT chỉ chứa thông tin tối thiểu (User ID, Role, Expired).
* Mọi API yêu cầu xác thực phải đi qua JWT Middleware.

---

# 6. Định hướng UI Component

* **Login Form:** Form Email, Password và Remember Me.
* **Login Button:** Hiển thị Loading khi xác thực.
* **Validation Message:** Hiển thị lỗi nhập liệu và lỗi xác thực.
* **Notification:** Thông báo Login thành công hoặc thất bại.
* **Auth Provider (RefineJS):** Quản lý Login, Logout, Check Auth, Get Identity và Get Permissions.
* **Axios Interceptor:** Tự động gắn Access Token, Refresh Token và Retry Request.
* **Protected Route:** Bảo vệ các trang yêu cầu đăng nhập.
* **Unauthorized Page (401):** Hiển thị khi chưa xác thực.
* **Forbidden Page (403):** Hiển thị khi không đủ quyền truy cập.

---

# 7. API

| Method | Endpoint        | Mô tả                             |
| ------ | --------------- | --------------------------------- |
| POST   | `/auth/login`   | Đăng nhập                         |
| POST   | `/auth/refresh` | Làm mới Access Token              |
| POST   | `/auth/logout`  | Đăng xuất                         |
| GET    | `/auth/me`      | Lấy thông tin người dùng hiện tại |

---

# 8. Kiến trúc Module

Module Authentication chỉ xử lý xác thực và quản lý phiên đăng nhập.

```text
Auth Module
│
├── Login
├── Logout
├── Refresh Token
├── JWT Middleware
├── Authorization Middleware
└── Session Management
```

Các thông tin người dùng như Email, Name, Role, Status... được truy xuất từ **User Module** thông qua `UserRepository`. Điều này giúp tách biệt rõ trách nhiệm giữa **Authentication** (xác thực) và **User Management** (quản lý người dùng), phù hợp với kiến trúc DDD/Clean Architecture và thuận tiện mở rộng sang OAuth2, SSO hoặc RBAC trong tương lai.
