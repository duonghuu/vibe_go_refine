# FRONTEND PLAN: Quản lý loại bài viết (Post Types)

## 1. Mục tiêu
- Xây dựng cụm tính năng Quản lý loại bài viết (Post Types) bao gồm trang danh sách, trang thêm mới và trang chỉnh sửa dựa trên tài liệu ý tưởng `12-post-types-create-idea.md`.
- Đảm bảo tính nhất quán của giao diện người dùng, sử dụng các trang riêng biệt (Pages) thay vì Modal/Drawer cho thao tác Thêm/Sửa.
- Tích hợp các hook của Refine.js (`useTable`, `useForm`, `useShow`, `@refinedev/mui`) kết hợp với Material UI (MUI) để xây dựng UI và quản lý trạng thái, dữ liệu.

## 2. Phạm vi chức năng
- **Trang Danh sách (List):**
  - Hiển thị danh sách các loại bài viết dưới dạng bảng (DataGrid).
  - Tìm kiếm theo `Code` hoặc `Tên loại bài viết` (có debounce).
  - Lọc dữ liệu theo Trạng thái (ACTIVE/INACTIVE).
  - Hỗ trợ sắp xếp trên DataGrid và tính năng kéo thả (Drag-and-drop) để thay đổi `sort_order`.
  - Hỗ trợ thao tác xóa hoặc xóa hàng loạt (Bulk Actions) với cơ chế kiểm tra ràng buộc (không cho phép xóa nếu có bài viết đang sử dụng).
- **Trang Thêm mới (Create):** Form nhập liệu riêng biệt để tạo loại bài viết mới với các validation rules chặt chẽ.
- **Trang Chỉnh sửa (Edit):** Form chỉnh sửa loại bài viết tương tự trang Thêm mới nhưng điền sẵn dữ liệu.

## 3. Luồng màn hình (Main Content)

### 3.1. Trang Danh sách (List Page)
- **Header:** 
  - Tiêu đề: `Quản lý loại bài viết`.
  - Breadcrumb: `Trang chủ > Quản trị nội dung > Loại bài viết`.
  - Nút hành động: `Thêm loại bài viết mới` (Điều hướng sang trang Create).
- **Khu vực Tìm kiếm và Lọc:** TextField tìm kiếm, Select/ToggleButton cho bộ lọc trạng thái.
- **Bảng dữ liệu (Table):**
  - Cột: ID, Code, Tên loại, Trạng thái (hiển thị dưới dạng Badge), Thứ tự hiển thị, Thao tác (Sửa/Xóa).
  - Màu sắc Trạng thái: Xanh lá (ACTIVE), Đỏ (INACTIVE).
  - Các cột dữ liệu sử dụng hiển thị mật độ cao (dense/medium) tối ưu hiển thị như chuẩn Architecture.

### 3.2. Trang Thêm mới & Chỉnh sửa (Create/Edit Page)
- **Header:** 
  - Tiêu đề: `Thêm loại bài viết mới` / `Chỉnh sửa loại bài viết`.
  - Breadcrumb tương ứng.
- **Layout Form (CSS Grid / Flexbox):** Cấu trúc 1 cột hoặc 2 cột.
  - **Mã loại (Code) *:** TextField bắt buộc. Ràng buộc: Chữ hoa, số, dấu gạch dưới, không khoảng trắng, tối đa 50 ký tự. Chỉ khả dụng khi thêm mới, disable/read-only khi chỉnh sửa.
  - **Tên loại (Name) *:** TextField bắt buộc. Tối đa 100 ký tự.
  - **Thứ tự (Sort Order) *:** Number input (Tự động tính toán tăng dần khi tạo mới).
  - **Trạng thái (Status):** Switch hoặc Radio (ACTIVE/INACTIVE).
- **Sticky Actions (Thanh hành động):** Nút Lưu (Primary) và Nút Hủy (quay lại trang List, có cảnh báo Unsaved Changes).

## 4. Component Breakdown
- `PostTypeListPage`: Trang hiển thị danh sách, chứa bảng DataGrid và các bộ lọc.
- `PostTypeCreatePage`: Trang form thêm mới loại bài viết.
- `PostTypeEditPage`: Trang form chỉnh sửa loại bài viết.
- `PostTypeForm`: Component form dùng chung cho cả Create và Edit, tái sử dụng logic nhập liệu và validation.
- `PostTypeStatusBadge`: Component hiển thị trạng thái bằng màu sắc tương ứng.

## 5. Data & State Management

### 5.1. Refine Resources
- Khai báo resource `post_types` trong `App.tsx` với các đường dẫn list, create, edit tương ứng.

### 5.2. Refine Hooks & Form
- List: Sử dụng hook `useDataGrid` (nếu dùng `@refinedev/mui`) kết hợp với DataGrid của MUI.
- Create/Edit: Sử dụng `@refinedev/react-hook-form` với hook `useForm`.

### 5.3. Form Validation (Quy tắc kiểm tra)
- **Code:** `required`, `pattern: /^[A-Z0-9_]+$/` (chỉ chữ hoa, số và gạch dưới), `maxLength: 50`.
- **Name:** `required`, `maxLength: 100`.
- **Sort Order:** `min: 1`, pattern số nguyên dương.

## 6. Các tính năng nâng cao (Nghiệp vụ cốt lõi)
- **Debounce Search:** Tích hợp cơ chế hoãn hàm (debounce 300ms) trên Client cho ô tìm kiếm trước khi trigger filter lên Server (Sử dụng lodash debounce hoặc useDebounce hook).
- **Drag-and-drop (Reorder):** Sử dụng thư viện bên ngoài hoặc tính năng Row Reorder của DataGrid để cập nhật `sort_order` và gửi API.
- **Bulk Delete / Soft Delete Constraint:** Khi người dùng ấn nút Xóa, hiển thị Dialog xác nhận. Xử lý lỗi trả về từ Backend (qua interceptor của DataProvider) nếu `post_type` đang được bài viết sử dụng để hiển thị thông báo lỗi chính xác trên Toast.

## 7. Checklist Triển khai
- [ ] Khai báo resource `post_types` trong Refine config.
- [ ] Dựng UI trang `PostTypeListPage` (DataGrid, Search, Filter).
- [ ] Dựng component dùng chung `PostTypeForm`.
- [ ] Dựng UI trang `PostTypeCreatePage` và `PostTypeEditPage`.
- [ ] Tích hợp React Hook Form và cấu hình Validation theo Spec.
- [ ] Tích hợp tính năng Cảnh báo thoát trang (Unsaved changes).
- [ ] Test luồng tạo mới, báo lỗi trùng Code từ BE.
- [ ] Test luồng chỉnh sửa dữ liệu.
- [ ] Xử lý logic xóa 1 bản ghi và xóa nhiều (Bulk Action) kèm xác nhận.
- [ ] Tích hợp tính năng kéo thả đổi vị trí (Reorder).
