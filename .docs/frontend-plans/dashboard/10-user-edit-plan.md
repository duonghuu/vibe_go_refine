# KẾ HOẠCH FRONTEND: Chỉnh sửa người dùng (Admin Edit User)

## 1. Thông tin chung
- **Tính năng:** Màn hình Chỉnh sửa người dùng (Admin Edit User).
- **Mục đích:** Cho phép Quản trị viên cập nhật thông tin (Họ tên), phân quyền (Role) và trạng thái (Status) của tài khoản.
- **Vị trí/Route:** `apps/frontend/src/pages/users/edit.tsx` (Route tương ứng trong App.tsx: `/users/edit/:id`).

## 2. Kiến trúc & Framework
- **Framework:** Refine.js (`@refinedev/mui`, `@refinedev/core`) kết hợp với Material UI (MUI).
- **Layout:** Tận dụng component `<Edit>` từ `@refinedev/mui` làm container chính để hiển thị Header, Breadcrumb và các action (Save, Delete - tùy chọn). 
- **Data Hook:** Sử dụng `useForm` của `@refinedev/mui` (tích hợp react-hook-form) để tự động fetch data user (`useOne`) và xử lý submit update (`useUpdate`).
- **Notification:** Dùng toast/snackbar mặc định của Refine khi cập nhật thành công/thất bại.

## 3. Cấu trúc UI Component (Grid 2 cột)

Màn hình áp dụng bố cục `<Grid container spacing={3}>`.
Trên Desktop, cột chính chiếm `xs={12} md={8}`, cột phụ chiếm `xs={12} md={4}`.

### 3.1. Cột chính (Main Column)

**Card 1: Thông tin cơ bản (Basic Information)**
- `TextField`: Họ và tên (Bắt buộc, max 100 char). Bind với `name`.
- `TextField`: Email. Bind với `email`. Cấu hình `disabled` hoặc `InputProps={{ readOnly: true }}`.
- *Thiết kế: Dùng `<Card>` và `<CardContent>` của MUI.*

**Card 2: Phân quyền (Permission)**
- `Select` hoặc `RadioGroup`: Vai trò (Role). Lựa chọn: ADMIN, STAFF, CUSTOMER.
- *Lưu ý UI:* Có thể hiển thị thêm helper text giải thích các quyền. Nếu chọn ADMIN, hiển thị một `<Alert severity="warning">` cảnh báo cấp quyền quản trị.

**Card 3: Trạng thái tài khoản (Status)**
- `Switch` hoặc `RadioGroup`: Trạng thái (ACTIVE / INACTIVE).
- *Lưu ý UI:* Nếu chuyển sang INACTIVE, hiển thị `<Alert severity="error">` cảnh báo người dùng sẽ không thể đăng nhập.

### 3.2. Cột phụ (Side Column)

**Card 4: Trạng thái hiện tại (Account Status Badge)**
- Hiển thị nhanh các `Chip` hoặc `Badge` thể hiện Role và Status hiện tại của bản ghi (chưa phải giá trị đang edit).
- Ví dụ: `<Chip label="ADMIN" color="primary" />` `<Chip label="ACTIVE" color="success" />`.

**Card 5: Thông tin hệ thống (System Information)**
- Hiển thị thông tin chỉ đọc (Readonly): User ID, Ngày tạo (Created At), Ngày cập nhật (Updated At), Đăng nhập lần cuối (nếu có).
- Sử dụng `<Typography>` để hiển thị cấu trúc `Label: Value`.

**Card 6: Thông tin đăng nhập (Account / Security)**
- `TextField`: Password (Readonly, hiển thị dạng `********`).
- Nút Action: `Button` Reset Password (sẽ trigger một modal hoặc gọi API riêng, trong phạm vi form edit user thì nút này có thể là action phụ).

## 4. Xử lý Form & Validation (react-hook-form)

Sử dụng `useForm` và truyền validation rule vào `Controller` (hoặc register):

- **Họ tên (`name`):** 
  - `required: "Vui lòng nhập họ và tên"`
  - `maxLength: { value: 100, message: "Tối đa 100 ký tự" }`
- **Role (`role`):** 
  - Chỉ nhận các giá trị hợp lệ.
- **Status (`status`):** 
  - Chỉ nhận ACTIVE / INACTIVE.

## 5. Các bước triển khai (Workflow)

1. **Bước 1: Tạo Mock Data**
   - Cập nhật `.docs/mock-data/users/edit.json` (nếu cần) hoặc trả dữ liệu tĩnh từ mock data provider để test giao diện.

2. **Bước 2: Cập nhật Router/App.tsx**
   - Đảm bảo resource `users` được định nghĩa `edit: "/users/edit/:id"`.

3. **Bước 3: Code giao diện `pages/users/edit.tsx`**
   - Bọc trang trong component `<Edit>`.
   - Setup `useForm` của `@refinedev/mui`.
   - Xây dựng layout Grid 2 cột.
   - Thêm các Controller của react-hook-form cho từng input field (Name, Email readonly, Role, Status).
   - Thêm cảnh báo bằng component `<Alert>`.
   - Render UI cột phụ với thông tin hệ thống lấy từ dữ liệu trả về `queryResult.data`.

4. **Bước 4: Loại bỏ Drawer ở Danh sách (`pages/users/list.tsx`)**
   - Cập nhật lại `list.tsx` để nút Edit (trên mỗi dòng của DataGrid) khi click sẽ chuyển hướng sang trang `/users/edit/:id` thay vì mở Drawer.
   - (Sử dụng `<EditButton hideText recordItemId={row.id} />` của `@refinedev/mui`).
