# FRONTEND PLAN: Quản trị người dùng (Admin User Management)

## 1. Mục tiêu
- Xây dựng màn hình quản lý người dùng toàn diện cho hệ thống bằng Refine.js và Material UI.
- Cung cấp công cụ tìm kiếm, lọc, phân trang và thao tác dữ liệu (Tạo, Sửa, Khóa/Mở khoá tài khoản, Reset Password).
- Tuân thủ nghiêm ngặt hệ thống Design System (`STYLEGUIDE.md`) và luồng nghiệp vụ bảo mật (`ARCHITECTURE.md`).

## 2. Phạm vi chức năng
### 2.1. Danh sách người dùng
- Hiển thị dữ liệu dạng bảng DataGrid (`dense` row density).
- Phân trang chuẩn tương thích với data provider (page, limit, total).
- Search: Tìm kiếm theo Tên, Email (bắt buộc dùng `debounce 300ms`).
- Filter: Lọc theo Role (ADMIN, STAFF, CUSTOMER) và Status (ACTIVE, INACTIVE).
- Sort: ID, Name, Email, Created At, Last Login.

### 2.2. Thao tác trên từng dòng (Row Actions)
- Xem chi tiết người dùng.
- Chỉnh sửa thông tin (Tên, Role, Status - không đổi Email).
- Đặt lại mật khẩu (Reset Password).
- Vô hiệu hóa (Disable) / Kích hoạt (Enable) tài khoản.

### 2.3. Thao tác hàng loạt (Bulk Actions)
- Chọn nhiều user thông qua checkbox trên DataGrid.
- Enable/Disable hàng loạt.

### 2.4. Khối Thống kê nhanh (Metrics)
- Tổng số lượng người dùng.
- Tổng số lượng Admin.
- Tổng số lượng Staff.
- Số lượng User đang hoạt động (ACTIVE).

## 3. Luồng màn hình
### 3.1. Header & Breadcrumb
- Tiêu đề: **Quản lý người dùng**.
- Breadcrumb: `Trang chủ > Quản trị hệ thống > Người dùng`.
- Nút hành động: `Thêm người dùng` (Primary Action - màu `palette.primary.main`), `Làm mới`.

### 3.2. Form thêm/sửa người dùng (Mutation)
- Vì form có ít trường (< 5 trường nhập liệu), sử dụng **Drawer** (trượt từ phải sang) thông qua `useDrawerForm` để không làm mất dấu bộ lọc dữ liệu hiện tại của admin.
- Validate: Bắt buộc điền Name, Email, Password (min 8 kí tự khi tạo), Role, Status hợp lệ.

### 3.3. Dialog & Phản hồi hệ thống
- Mọi thao tác đổi trạng thái (Disable/Enable/Reset Password) bắt buộc dùng `Confirmation Dialog` trước khi gọi API.
- Sau khi hành động thành công, kích hoạt Snackbar/Toast tự đóng sau 3 giây.

## 4. Component Breakdown
### 4.1. Page Container
- `UserList`: Component chính bọc toàn bộ trang, kết nối hook data Refine.

### 4.2. UI Blocks
- `UserSummaryCards`: Component hiển thị 4 thẻ thống kê trên cùng.
- `UserFilterBar`: Khu vực tìm kiếm (có debounce) và lọc.
- `UserTable`: Bọc `DataGrid` của Material UI.
- `UserRowActions`: Cụm nút thao tác (Eye, Edit, Key - Reset password, Lock/Unlock) ở cột Pinned Actions cuối.
- `UserFormDrawer`: Drawer chứa form tạo mới hoặc cập nhật người dùng.
- `ResetPasswordDialog`: Hộp thoại yêu cầu nhập mật khẩu mới và xác nhận.
- `ConfirmStatusDialog`: Hộp thoại xác nhận thay đổi trạng thái tài khoản.

### 4.3. Shared Primitives
- `RoleBadge`: Chip màu đại diện cho quyền (`ADMIN` - Xanh dương, `STAFF` - Tím, `CUSTOMER` - Xám).
- `StatusBadge`: Chip màu trạng thái (`ACTIVE` - Xanh lá, `INACTIVE` - Đỏ).

## 5. Data & State
### 5.1. Resource
- Định nghĩa trong Refine app: `users`

### 5.2. Các trường dữ liệu cốt lõi
- `id`: Định danh duy nhất.
- `email`: Khóa đăng nhập chính (readonly khi edit).
- `name`: Tên đầy đủ.
- `role`: Phân quyền (ADMIN, STAFF, CUSTOMER).
- `status`: Trạng thái (ACTIVE, INACTIVE).
- `lastLoginAt`: Đăng nhập cuối.
- `createdAt`: Ngày tạo.

### 5.3. Frontend State
- Từ khoá tìm kiếm & filter state (liên kết với useDataGrid).
- Trạng thái Drawer (open/close, id đang edit).
- Trạng thái các Dialog (Reset Password, Confirm Status).

## 6. Quy tắc Nghiệp vụ (Business Rules) cần lưu ý ở giao diện
- Không hiển thị nút Disable cho chính user đang đăng nhập (lấy từ auth state).
- Nút Action (Sửa, Đổi trạng thái) nên ẩn hoặc disable nếu user không có quyền tương ứng (vd: Staff không được sửa user Admin).
- Màu sắc nút "Lưu thay đổi" mang tính chốt hạ dùng `palette.success.main`.
- Nút xoá (nếu có) hoặc khoá dùng màu nguy hiểm, các nút huỷ dùng `palette.secondary.main`.

## 7. Checklist Triển khai
- [ ] Cấu hình resource `users` trong `App.tsx` (list, create, edit, show).
- [ ] Tạo file `pages/users/list.tsx`.
- [ ] Xây dựng block thống kê `UserSummaryCards`.
- [ ] Xây dựng thanh công cụ `UserFilterBar` với `debounce`.
- [ ] Render MUI `DataGrid` cấu hình `dense`, kèm `UserRowActions`.
- [ ] Viết các helper components `RoleBadge`, `StatusBadge`.
- [ ] Xây dựng `UserFormDrawer` dùng cho Create/Edit.
- [ ] Xây dựng `ResetPasswordDialog` và `ConfirmStatusDialog`.
- [ ] Kiểm tra responsive, error handling và loading skeletons.
