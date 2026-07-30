# Ý TƯỞNG: Quản trị Chỉnh sửa người dùng (Admin Edit User)

## 1. Thông tin chung (Meta Info)

* **Dự án:** TechBite.
* **Tính năng:** Màn hình Chỉnh sửa người dùng (Admin Edit User).
* **Mục đích:** Cho phép Quản trị viên cập nhật thông tin, phân quyền và trạng thái của tài khoản người dùng trong hệ thống.

---

# 2. Đối tượng & Trải nghiệm (Target & UX)

* **Người dùng chính:** Quản trị viên (Admin).

* **Hành động chính**

  * Cập nhật thông tin cá nhân.
  * Thay đổi vai trò (Role).
  * Thay đổi trạng thái tài khoản.
  * Theo dõi thông tin đăng nhập.
  * Lưu thay đổi hoặc hủy thao tác.

* **Cảm xúc mang lại**

  Form được nhóm thành từng khu vực rõ ràng, giúp Admin dễ quan sát và cập nhật thông tin mà không bị rối khi số lượng trường dữ liệu tăng lên trong tương lai.

---

# 3. Đặc tả Thiết kế (Design Specs)

## Phong cách UI

* Bố cục dạng Grid/Card gồm **2 cột** (Main Column và Side Column).
* Thiết kế đồng nhất với Product Edit.
* Card bo góc (`rounded-xl` hoặc `rounded-2xl`).
* Sticky Action Bar luôn hiển thị nút lưu.
* Responsive:

  * Desktop: 2 cột.
  * Tablet/Mobile: 1 cột.

---

## Màu sắc chủ đạo (Brand Colors)

* **Cam thương hiệu:** Nút Lưu `bg-[#ff8c42]`
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
Chỉnh sửa người dùng: Nguyễn Văn A
```

* Breadcrumb

```
Trang chủ
>
Quản trị hệ thống
>
Người dùng
>
Chỉnh sửa
```

### Nhóm nút hành động

* Hủy bỏ
* Lưu thay đổi

---

# Bố cục 2 cột (2-Column Layout)

## Cột chính (Main Column - khoảng 2/3)

### 1. Thông tin người dùng (Basic Information Card)

Hiển thị các trường

* Họ và tên (*)
* Email (Readonly)
* Ghi chú (Tùy chọn, nếu hệ thống hỗ trợ)

---

### 2. Phân quyền (Permission Card)

* Role

```
ADMIN
STAFF
CUSTOMER
```

Hiển thị mô tả quyền tương ứng.

Nếu chuyển sang ADMIN

Hiển thị Alert cảnh báo.

---

### 3. Trạng thái tài khoản (Status Card)

* Status

```
ACTIVE
INACTIVE
```

Hiển thị bằng Switch hoặc Radio.

Nếu chuyển sang INACTIVE

Hiển thị cảnh báo

> Người dùng sẽ không thể đăng nhập sau khi lưu.

---

## Cột phụ (Side Column - khoảng 1/3)

### 1. Thông tin hệ thống (System Information Card)

Hiển thị

* User ID
* Created At
* Updated At
* Last Login

Toàn bộ chỉ đọc.

---

### 2. Thông tin đăng nhập (Account Card)

Hiển thị

* Email
* Password

Password

```
************
```

Không cho chỉnh sửa.

Có nút

```
Reset Password
```

---

### 3. Trạng thái tài khoản (Account Status Card)

Hiển thị Badge

```
Role Badge

Status Badge
```

Ví dụ

```
ADMIN

ACTIVE
```

Giúp Admin nhìn nhanh trạng thái hiện tại.

---

# 4. Kiểm tra dữ liệu (Validation & Error Handling)

### Họ tên

* Bắt buộc.
* Tối đa 100 ký tự.

---

### Email

* Chỉ đọc.
* Không cho thay đổi.

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

### Business Validation

* Không được Disable chính tài khoản đang đăng nhập.
* Không được hạ quyền Admin cuối cùng.
* Chỉ ADMIN được cấp quyền ADMIN.

---

# 5. Chức năng nổi bật

* **Real-time Validation** khi nhập Name hoặc thay đổi Role/Status.
* **Unsaved Changes Warning** khi rời trang mà chưa lưu.
* **Toast Notification** khi cập nhật thành công hoặc thất bại.
* **Confirmation Dialog** khi:

  * Disable User.
  * Thay đổi Role thành ADMIN.
* **Refresh Detail** sau khi lưu thành công.

---

# 6. Định hướng UI Component

### Page Layout

* Header
* Breadcrumb
* Sticky Action Bar

### Form Layout

* Grid
* Card

### Form Components

* TextField
* Select
* Switch
* Alert
* Badge
* Divider

### Action Components

* Save Button
* Cancel Button
* Reset Password Button

### Feedback Components

* Toast / Snackbar
* Confirm Dialog

