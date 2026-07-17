# Ý TƯỞNG: Quản trị người dùng (Admin User Management)

## 1. Thông tin chung (Meta Info)

* **Dự án:** TechBite.
* **Tính năng:** Màn hình Quản trị người dùng (Admin User Management).
* **Mục đích:** Quản lý toàn bộ tài khoản người dùng trong hệ thống, bao gồm tạo mới, cập nhật thông tin, phân quyền, kích hoạt/vô hiệu hóa tài khoản, đặt lại mật khẩu và theo dõi lịch sử đăng nhập gần nhất.

---

# 2. Đối tượng & Trải nghiệm (Target & UX)

* **Người dùng chính:** Quản trị viên (Admin).

* **Hành động chính:**

  * Tìm kiếm nhanh người dùng theo Email hoặc Tên.
  * Lọc theo Vai trò (Role) và Trạng thái (Status).
  * Sắp xếp theo ID, Tên, Email, Ngày tạo hoặc Lần đăng nhập cuối.
  * Thêm người dùng mới.
  * Chỉnh sửa thông tin người dùng.
  * Reset mật khẩu.
  * Khóa hoặc mở khóa tài khoản.
  * Xem thông tin chi tiết người dùng.

* **Cảm xúc mang lại:** Giao diện quản trị chuyên nghiệp, dễ quan sát, thao tác nhanh, giúp quản trị viên kiểm soát tài khoản và phân quyền hiệu quả.

---

# 3. Đặc tả Thiết kế (Design Specs)

## Phong cách UI

* Dashboard hiện đại, tối giản.
* Bảng dữ liệu (Data Table) hỗ trợ phân trang, tìm kiếm, lọc và sắp xếp.
* Card thống kê phía trên.
* Các Badge thể hiện Role và Status rõ ràng.
* Modal hoặc Drawer cho Create/Edit User.

---

## Màu sắc chủ đạo (Brand Colors)

* **Cam (Thêm mới):** `bg-[#ff8c42]`
* **Xanh lá (ACTIVE):** `bg-green-500`
* **Đỏ (INACTIVE):** `bg-red-500`
* **Xanh dương (ADMIN):** `bg-blue-500`
* **Tím (STAFF):** `bg-purple-500`
* **Xám (CUSTOMER):** `bg-gray-500`

---

## Cấu trúc màn hình (Top to Bottom)

### Header

* Tiêu đề: **Quản lý người dùng**
* Breadcrumb

```
Trang chủ > Quản trị hệ thống > Người dùng
```

Nhóm nút hành động

* Nút **Thêm người dùng**
* Nút **Làm mới dữ liệu**

---

### Thống kê nhanh (Summary Cards)

Hiển thị 4 Card

* Tổng số User
* Tổng Admin
* Tổng Staff
* User đang hoạt động

(Có thể mở rộng thêm Tổng Customer hoặc User bị khóa.)

---

### Thanh tìm kiếm & Bộ lọc

#### Search

* Email
* Tên

#### Filter

* Role

  * ADMIN
  * STAFF
  * CUSTOMER

* Status

  * ACTIVE
  * INACTIVE

#### Sort

* ID
* Name
* Email
* Created At
* Last Login

#### Nút

* Search
* Reset Filter

---

### Danh sách người dùng (User Table)

Các cột hiển thị

| Cột        | Mô tả                  |
| ---------- | ---------------------- |
| ID         | ID người dùng          |
| Email      | Email đăng nhập        |
| Tên        | Họ tên                 |
| Role       | Badge quyền            |
| Status     | Badge trạng thái       |
| Last Login | Lần đăng nhập gần nhất |
| Created At | Ngày tạo               |
| Action     | Các thao tác           |

---

### Action Menu

Mỗi dòng có các chức năng

* Xem chi tiết
* Chỉnh sửa
* Reset Password
* Disable
* Enable

---

### Bulk Action (Khuyến nghị)

Khi chọn nhiều User

* Enable hàng loạt
* Disable hàng loạt

(Lưu ý: Không cho phép Disable chính tài khoản đang đăng nhập hoặc Admin cuối cùng của hệ thống.)

---

# 4. Dữ liệu cốt lõi (Mock Data)

| ID | Email                                             | Tên          | Role     | Status   | Last Login       | Created At |
| -- | ------------------------------------------------- | ------------ | -------- | -------- | ---------------- | ---------- |
| 1  | [admin@techbite.com](mailto:admin@techbite.com)   | System Admin | ADMIN    | ACTIVE   | 2026-07-17 08:30 | 2026-01-01 |
| 2  | [staff1@techbite.com](mailto:staff1@techbite.com) | Nguyễn Văn A | STAFF    | ACTIVE   | 2026-07-16 17:40 | 2026-03-12 |
| 3  | [staff2@techbite.com](mailto:staff2@techbite.com) | Trần Thị B   | STAFF    | INACTIVE | 2026-07-10 09:20 | 2026-03-20 |
| 4  | [customer1@gmail.com](mailto:customer1@gmail.com) | Lê Văn C     | CUSTOMER | ACTIVE   | 2026-07-17 09:05 | 2026-05-08 |
| 5  | [customer2@gmail.com](mailto:customer2@gmail.com) | Phạm Thị D   | CUSTOMER | INACTIVE | 2026-06-29 15:15 | 2026-05-18 |

---

# 5. Chức năng nổi bật cho Quản lý User

## Quản lý phân quyền (Role Management)

Hệ thống hỗ trợ 3 vai trò:

* ADMIN
* STAFF
* CUSTOMER

Badge hiển thị màu khác nhau giúp dễ nhận biết.

---

## Quản lý trạng thái tài khoản

Hai trạng thái

* ACTIVE
* INACTIVE

Khi Disable

* Không thể đăng nhập.
* JWT Token hiện tại bị vô hiệu (nếu hệ thống hỗ trợ Blacklist/Session Management).

---

## Reset Password

Admin có thể

* Nhập mật khẩu mới.
* Xác nhận mật khẩu.
* Hệ thống mã hóa bằng BCrypt trước khi lưu.

---

## Theo dõi đăng nhập

Hiển thị

* Last Login

Giúp Admin dễ dàng phát hiện

* User lâu không sử dụng
* User mới đăng nhập
* Tài khoản chưa từng đăng nhập

---

## Tìm kiếm & Bộ lọc

Cho phép kết hợp nhiều điều kiện

* Search Email
* Search Name
* Filter Role
* Filter Status

Giúp tìm kiếm nhanh trong số lượng lớn tài khoản.

---

## Sắp xếp dữ liệu

Hỗ trợ Sort

* ID
* Email
* Name
* Created At
* Last Login

Ascending / Descending.

---

## Phân trang

Hỗ trợ

* Page
* Limit
* Total Record

Phù hợp với lượng dữ liệu lớn.

---

## Kiểm tra dữ liệu đầu vào

### Create User

* Email bắt buộc.
* Email không được trùng.
* Password tối thiểu 8 ký tự.
* Name bắt buộc.
* Role hợp lệ.
* Status hợp lệ.

### Update User

* Không cho phép thay đổi Email (khuyến nghị để đảm bảo tính ổn định của định danh đăng nhập).
* Chỉ cập nhật:

  * Name
  * Role
  * Status

---

## Quy tắc nghiệp vụ (Business Rules)

* Không được xóa tài khoản Admin cuối cùng.
* Không được Disable chính tài khoản đang đăng nhập.
* Chỉ ADMIN mới được tạo hoặc phân quyền ADMIN.
* Email là duy nhất trong toàn hệ thống.
* Password luôn được lưu dưới dạng `password_hash` (BCrypt).
* `last_login_at` chỉ được cập nhật sau khi đăng nhập thành công.

---

# 6. Định hướng UI Component

* **Summary Cards:** Hiển thị thống kê nhanh số lượng User theo Role và Status.
* **User Data Table:** Bảng dữ liệu hỗ trợ Search, Filter, Sort và Pagination.
* **Create/Edit User Drawer (hoặc Modal):** Form thêm/sửa người dùng.
* **User Detail Drawer:** Hiển thị thông tin chi tiết người dùng.
* **Reset Password Dialog:** Hộp thoại nhập mật khẩu mới và xác nhận.
* **Enable/Disable Confirmation Dialog:** Xác nhận trước khi thay đổi trạng thái tài khoản.
* **Role Badge:** Badge màu theo từng vai trò (ADMIN, STAFF, CUSTOMER).
* **Status Badge:** Badge màu theo trạng thái (ACTIVE, INACTIVE).
* **Filter Panel:** Bộ lọc Role, Status kết hợp ô tìm kiếm.
* **Pagination Component:** Điều hướng và thay đổi số bản ghi mỗi trang.
