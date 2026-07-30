# FRONTEND PLAN: Thêm người dùng mới (Create User)

## 1. Mục tiêu
- Xây dựng màn hình Thêm người dùng mới (Admin Create User) dựa trên tài liệu ý tưởng `11-user-create-idea.md`.
- Sử dụng các hook của Refine.js (như `useForm`) và các component Material UI (MUI) kết hợp `react-hook-form` để quản lý state và validation.
- Thiết kế layout 2 cột nhất quán với các màn hình trước đó (Product Create/Edit), chia các nhóm thông tin logic (Thông tin cơ bản, Mật khẩu, Phân quyền, Trạng thái) giúp Admin thao tác dễ dàng và tránh sai sót.

## 2. Phạm vi chức năng
- **Thông tin cơ bản:** Họ và tên, Email đăng nhập.
- **Thông tin đăng nhập:** Mật khẩu, Xác nhận mật khẩu (Hỗ trợ ẩn/hiện mật khẩu, kiểm tra độ mạnh mật khẩu).
- **Phân quyền (Role):** Chọn quyền cho người dùng (ADMIN, STAFF, CUSTOMER) kèm mô tả/cảnh báo chi tiết.
- **Trạng thái tài khoản:** Kích hoạt (ACTIVE) hoặc Vô hiệu hóa (INACTIVE) kèm cảnh báo.
- **Thông tin hệ thống:** Hiển thị dưới dạng Placeholder các trường sẽ được tự động tạo (User ID, Created At, Updated At, Last Login).
- **Hành động form:** Hủy bỏ (quay về danh sách), Tạo người dùng mới.

## 3. Luồng màn hình (Main Content)

### 3.1. Header & Page Title
- **Tiêu đề:** `Thêm người dùng mới`
- **Breadcrumb:** `Trang chủ > Quản trị hệ thống > Người dùng > Thêm mới`

### 3.2. Form Layout (Sử dụng CSS Grid / Flexbox)
Bố cục 2 cột (Main Column chiếm khoảng 2/3, Side Column chiếm 1/3) trên Desktop, 1 cột trên Tablet/Mobile:

- **Cột chính (Main Column - 8/12):**
  - **Card "Thông tin người dùng":**
    - **Họ và tên (*):** Bắt buộc, tối đa 100 ký tự. Tự động Focus khi mở trang.
    - **Email (*):** Bắt buộc, chuẩn định dạng Email.
  - **Card "Thông tin đăng nhập":**
    - **Mật khẩu (*):** Bắt buộc, tối thiểu 8 ký tự. Tích hợp nút Toggle (Hiện/Ẩn) và thanh hiển thị độ mạnh mật khẩu.
    - **Xác nhận mật khẩu (*):** Bắt buộc, phải khớp với trường Mật khẩu.
  - **Card "Phân quyền":**
    - **Role (*):** Radio hoặc Select (ADMIN, STAFF, CUSTOMER). Hiển thị Alert cảnh báo nếu chọn ADMIN.

- **Cột phụ (Side Column - 4/12):**
  - **Card "Trạng thái tài khoản":**
    - **Trạng thái:** MUI `Select` hoặc `RadioGroup` (ACTIVE, INACTIVE). Mặc định là ACTIVE. Hiển thị cảnh báo nếu chọn INACTIVE.
  - **Card "Thông tin hệ thống":**
    - Hiển thị Text dạng Placeholder/Disabled cho: User ID, Created At, Updated At, Last Login (Label: "Tự động tạo").

### 3.3. Sticky Actions (Thanh hành động)
- Nút **Hủy bỏ:** Hủy thao tác, hiển thị hộp thoại xác nhận nếu đã nhập liệu (Unsaved Changes Warning), điều hướng về `/users`.
- Nút **Tạo người dùng (Primary):** Trigger submit form. Thành công hiển thị Toast và quay về danh sách User.

## 4. Component Breakdown
- `UserCreatePage`: Component chính, sử dụng thẻ `<Create>` của `@refinedev/mui` và hook `useForm`.
- `UserBasicInfoCard`: Chứa input Họ và tên, Email.
- `UserAuthCard`: Chứa input Mật khẩu, Xác nhận mật khẩu, có logic Toggle Password Visibility và tính toán Password Strength.
- `UserPermissionCard`: Chứa lựa chọn Role và Alert hiển thị cảnh báo tương ứng.
- `UserStatusCard`: Chứa lựa chọn trạng thái (ACTIVE/INACTIVE).
- `UserSystemInfoCard`: Chứa thông tin placeholder hệ thống.

## 5. Data & State Management

### 5.1. Refine Hooks & Form
- Sử dụng `@refinedev/react-hook-form` với hook `useForm` kết nối resource `users`.
- `mutationMode` mặc định (pessimistic) để xử lý tạo mới.

### 5.2. Form Validation (Quy tắc kiểm tra)
- **Họ và tên:** `required`, `maxLength: 100`.
- **Email:** `required`, `pattern` email hợp lệ. (Validation trùng email sẽ do backend trả về và form hiển thị lỗi).
- **Mật khẩu:** `required`, `minLength: 8`.
- **Xác nhận mật khẩu:** Custom validation kiểm tra khớp với trường `mật khẩu`.
- **Role:** `required` (ADMIN, STAFF, CUSTOMER).
- **Status:** `required` (ACTIVE, INACTIVE).

## 6. Các tính năng nâng cao (Nghiệp vụ cốt lõi)
- **Real-time Validation:** Sử dụng `mode: "onBlur"` hoặc `"onChange"` để feedback tức thời cho người nhập.
- **Auto Focus:** Tự động focus vào ô "Họ và tên" khi render trang lần đầu (`autoFocus` prop).
- **Password Visibility Toggle:** Sử dụng `InputAdornment` của MUI để thêm icon mắt Hiện/Ẩn password.
- **Password Strength Indicator:** Viết hàm (hoặc dùng thư viện nhẹ) để đánh giá độ mạnh password (chữ hoa, chữ thường, số, ký tự đặc biệt) và hiển thị thanh tiến trình (LinearProgress) màu tương ứng (Yếu - Đỏ, Vừa - Vàng, Mạnh - Xanh).
- **Unsaved Changes Warning:** Tích hợp `warnWhenUnsavedChanges: true` trong cấu hình hook của Refine.
- **Toast Notification:** Refine mặc định hỗ trợ, nhưng cần đảm bảo hook bắt được error message từ Backend (ví dụ: "Email đã tồn tại") để hiển thị toast hoặc gán lỗi vào input Email.

## 7. Checklist Triển khai
- [ ] Tạo file `apps/frontend/src/pages/users/create.tsx` (hoặc cấu trúc thư mục tương đương).
- [ ] Khai báo page `UserCreate` vào Refine resources trong `App.tsx` (nếu chưa có).
- [ ] Khởi tạo `useForm` với các default values (`status: 'ACTIVE'`, `role: 'CUSTOMER'`).
- [ ] Xây dựng layout 2 cột Grid, đảm bảo responsive.
- [ ] Xây dựng `UserBasicInfoCard` (Name, Email).
- [ ] Xây dựng `UserAuthCard` (Password, Confirm Password + Nút Toggle Mắt).
- [ ] Thêm logic tính toán và hiển thị Password Strength Indicator.
- [ ] Xây dựng `UserPermissionCard` (Role Selection + Alert).
- [ ] Xây dựng `UserStatusCard` và `UserSystemInfoCard`.
- [ ] Định nghĩa Validation Rules cho các trường sử dụng react-hook-form.
- [ ] Tích hợp tính năng Cảnh báo thoát trang (Unsaved changes).
- [ ] Kiểm tra API `POST /users`, đảm bảo xử lý lỗi backend trả về và điều hướng thành công.
