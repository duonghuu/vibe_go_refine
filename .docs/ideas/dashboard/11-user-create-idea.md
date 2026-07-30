# Ý TƯỞNG: Quản trị Thêm người dùng (Admin Create User)

## 1. Thông tin chung (Meta Info)

* **Dự án:** TechBite.
* **Tính năng:** Màn hình Thêm người dùng (Admin Create User).
* **Mục đích:** Cho phép Quản trị viên tạo mới tài khoản người dùng, thiết lập thông tin cơ bản, mật khẩu ban đầu, phân quyền và trạng thái tài khoản.

---

# 2. Đối tượng & Trải nghiệm (Target & UX)

* **Người dùng chính:** Quản trị viên (Admin).

### Hành động chính

* Nhập thông tin người dùng.
* Thiết lập Email đăng nhập.
* Thiết lập mật khẩu ban đầu.
* Phân quyền người dùng.
* Thiết lập trạng thái tài khoản.
* Lưu người dùng mới hoặc hủy thao tác.

### Cảm xúc mang lại

Form được chia thành các nhóm thông tin rõ ràng, thống nhất với màn hình Edit User và Edit Product. Các trường quan trọng được sắp xếp hợp lý giúp Admin tạo tài khoản nhanh chóng và hạn chế sai sót.

---

# 3. Đặc tả Thiết kế (Design Specs)

## Phong cách UI

* Bố cục dạng Grid/Card gồm **2 cột** (Main Column và Side Column).
* Thiết kế thống nhất với Product Create/Edit và User Edit.
* Card bo góc (`rounded-xl` hoặc `rounded-2xl`).
* Sticky Action Bar luôn hiển thị nút lưu.
* Responsive:

  * Desktop: 2 cột.
  * Tablet/Mobile: 1 cột.

---

## Màu sắc chủ đạo (Brand Colors)

* **Cam thương hiệu:** Nút Tạo mới `bg-[#ff8c42]`
* **Đỏ:** Validation/Error `text-red-500`
* **Xanh lá:** ACTIVE
* **Đỏ:** INACTIVE
* **Xanh dương:** ADMIN
* **Tím:** STAFF
* **Xám:** CUSTOMER

---

## Cấu trúc màn hình (Top to Bottom)

### Header & Sticky Actions

* Tiêu đề

```
Thêm người dùng mới
```

* Breadcrumb

```
Trang chủ
>
Quản trị hệ thống
>
Người dùng
>
Thêm mới
```

### Nhóm nút hành động

* Hủy bỏ
* Tạo người dùng

---

# Bố cục 2 cột (2-Column Layout)

## Cột chính (Main Column - khoảng 2/3)

### 1. Thông tin người dùng (Basic Information Card)

Hiển thị các trường

* Họ và tên (*)
* Email (*)

---

### 2. Thông tin đăng nhập (Authentication Card)

* Mật khẩu (*)
* Xác nhận mật khẩu (*)

Mật khẩu hỗ trợ

* Hiện/Ẩn mật khẩu.
* Hiển thị yêu cầu về độ mạnh mật khẩu.

---

### 3. Phân quyền (Permission Card)

Role

```
ADMIN
STAFF
CUSTOMER
```

Hiển thị mô tả từng quyền.

Nếu chọn ADMIN

Hiển thị Alert

> Người dùng này sẽ có toàn quyền quản trị hệ thống.

---

## Cột phụ (Side Column - khoảng 1/3)

### 1. Trạng thái tài khoản (Status Card)

Trạng thái

```
ACTIVE
INACTIVE
```

Giá trị mặc định

```
ACTIVE
```

Nếu chọn INACTIVE

Hiển thị cảnh báo

> Người dùng sẽ chưa thể đăng nhập cho đến khi được kích hoạt.

---

### 2. Thông tin hệ thống (System Information Card)

Hiển thị thông tin sẽ được hệ thống tự động tạo sau khi lưu.

* User ID
* Created At
* Updated At
* Last Login

Tất cả hiển thị dạng Placeholder hoặc "Tự động tạo".

---

# 4. Kiểm tra dữ liệu (Validation & Error Handling)

### Họ và tên

* Bắt buộc.
* Tối đa 100 ký tự.

---

### Email

* Bắt buộc.
* Đúng định dạng Email.
* Không được trùng trong hệ thống.

---

### Mật khẩu

* Bắt buộc.
* Tối thiểu 8 ký tự.
* Khuyến nghị gồm chữ hoa, chữ thường, số và ký tự đặc biệt.

---

### Xác nhận mật khẩu

* Bắt buộc.
* Phải giống mật khẩu.

---

### Role

Chỉ nhận

* ADMIN
* STAFF
* CUSTOMER

---

### Status

Chỉ nhận

* ACTIVE
* INACTIVE

---

# 5. Chức năng nổi bật

* **Real-time Validation** cho Email, Password và Confirm Password.
* **Password Visibility Toggle** cho phép hiện/ẩn mật khẩu.
* **Password Strength Indicator** hiển thị mức độ mạnh của mật khẩu.
* **Unsaved Changes Warning** khi rời trang mà chưa lưu.
* **Toast Notification** khi tạo thành công hoặc thất bại.
* **Tự động Focus** vào trường Họ và tên khi mở trang.

---

# 6. Quy tắc nghiệp vụ (Business Rules)

* Email phải là duy nhất trong toàn hệ thống.
* Password được mã hóa bằng BCrypt trước khi lưu.
* Chỉ ADMIN mới được tạo tài khoản có Role là ADMIN.
* User mới mặc định có trạng thái **ACTIVE** (có thể thay đổi trước khi lưu).
* `created_at`, `updated_at`, `last_login_at` được hệ thống tự động sinh.
* `last_login_at` ban đầu bằng `NULL`.

---

# 7. API

## Tạo người dùng

```
POST /users
```

Request Body

```json
{
  "name": "Nguyễn Văn A",
  "email": "staff@techbite.com",
  "password": "Password@123",
  "confirm_password": "Password@123",
  "role": "STAFF",
  "status": "ACTIVE"
}
```

---

# 8. Chức năng sau khi tạo

## Thành công

* Hiển thị Toast

```
Tạo người dùng thành công.
```

* Chuyển về màn hình danh sách User.
* Danh sách tự Refresh và giữ nguyên bộ lọc, phân trang hiện tại.

---

## Thất bại

Hiển thị lỗi từ Backend.

Ví dụ

```
Email đã tồn tại trong hệ thống.
```

hoặc

```
Bạn không có quyền tạo tài khoản ADMIN.
```

---

# 9. Định hướng UI Component

### Page Layout

* Header
* Breadcrumb
* Sticky Action Bar

### Form Layout

* Grid
* Card

### Form Components

* TextField
* PasswordField
* Select
* Switch hoặc Radio
* Alert
* Divider

### Action Components

* Create Button
* Cancel Button

### Feedback Components

* Password Strength Indicator
* Toast / Snackbar
* Confirm Dialog (Unsaved Changes)
