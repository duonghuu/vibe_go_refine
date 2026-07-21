# FRONTEND PLAN: Chỉnh sửa danh mục sản phẩm (Edit Category)

## 1. Mục tiêu
- Xây dựng màn hình Chỉnh sửa danh mục sản phẩm (Edit Category) dựa trên tài liệu ý tưởng `07-category-edit-idea.md`.
- Kế thừa và tái sử dụng tối đa các components từ form Create (như UI layout, Tree Select).
- Đảm bảo tải chính xác dữ liệu hiện tại (Initial Data) của danh mục bằng các hook của Refine.js (`useForm`, `useShow`, v.v.).
- Xử lý các quy tắc nghiệp vụ đặc thù cho việc chỉnh sửa như giới hạn danh mục cha, cảnh báo chuyển trạng thái và quản lý trạng thái form thay đổi.

## 2. Phạm vi chức năng
- **Thông tin cơ bản:** Cập nhật Tên danh mục, Slug, Mô tả.
- **Danh mục cha:** Thay đổi danh mục cha (ngăn chọn chính nó hoặc vòng lặp cây).
- **Hình ảnh:** Hiển thị ảnh hiện tại, cho phép tải lên ảnh mới để thay thế hoặc xóa ảnh.
- **Thiết lập:** Cập nhật trạng thái (Hoạt động / Ẩn) và Thứ tự hiển thị.
- **Thông tin hệ thống:** Vùng chỉ đọc hiển thị ID, Ngày tạo/cập nhật, Người tạo/cập nhật.
- **Hành động form:** Hủy (với cảnh báo Unsaved Changes), Lưu thay đổi.

## 3. Luồng màn hình (Main Content)

### 3.1. Header & Page Title
- **Tiêu đề:** `Chỉnh sửa danh mục`
- **Breadcrumb:** `Trang chủ > Cấu hình hệ thống > Danh mục > Chỉnh sửa`

### 3.2. Form Layout (CSS Grid / Flexbox - Kế thừa từ Create Layout)
Sử dụng cấu trúc 2 cột tương tự màn hình Create:
- **Cột chính (Main Column):**
  - Card "Thông tin cơ bản":
    - **Tên danh mục (*):** Bắt buộc, MUI `TextField`.
    - **Slug:** MUI `TextField`. Có thể chỉnh sửa. Không tự động override nếu Tên thay đổi.
    - **Danh mục cha:** Custom `Select` hoặc Tree Select hiển thị dạng phân cấp (disable chính nó và danh mục con của nó).
    - **Mô tả:** MUI `TextField` multiline.
- **Cột phụ (Sidebar Column):**
  - Card "Hình ảnh":
    - Hiển thị ảnh hiện hành của danh mục (nếu có).
    - Nút / Vùng Upload thay thế, Nút Xóa ảnh. Validate dung lượng và định dạng giống Create.
  - Card "Thiết lập":
    - **Trạng thái:** MUI `Select` / `RadioGroup`. Cảnh báo nếu chuyển sang "Ẩn".
    - **Thứ tự hiển thị:** MUI `TextField` dạng number.
  - Card "Thông tin hệ thống" (Chỉ đọc, hiển thị dạng text/list mờ):
    - ID, Ngày tạo, Ngày cập nhật, Người tạo, Người cập nhật.

### 3.3. Footer Actions
- Nút **Hủy:** Xóa dữ liệu chưa lưu và điều hướng về trang danh sách. Nếu có sửa đổi (dirty), bật dialog xác nhận.
- Nút **Lưu thay đổi:** Xác thực dữ liệu và gọi API PUT, thành công thì điều hướng về trang danh sách.

## 4. Component Breakdown
- `CategoryEditPage`: Component gốc sử dụng tính năng Edit từ `@refinedev/react-hook-form` (hoặc bọc Refine `Edit`).
- `CategoryBasicInfoCard`: Tái sử dụng/Chia sẻ từ Create, bổ sung logic truyền vào `disabled` cho parent category hợp lý.
- `CategoryMediaCard`: Tái sử dụng từ Create, xử lý truyền data ảnh hiện hành (`defaultImage`).
- `CategorySettingsCard`: Tái sử dụng từ Create.
- `SystemInfoCard`: Component mới chỉ dùng trong màn Edit hiển thị Audit logs của bản ghi.

## 5. Data & State Management
### 5.1. Refine Hooks & Form
- Sử dụng hook `useForm` với action là `"edit"` từ `@refinedev/react-hook-form` gắn với resource `categories`. Refine sẽ tự động fetch data record tương ứng và fill vào react-hook-form.
- Danh sách category cha: Sử dụng custom API để truyền tree nhưng phải filter (hoặc disable) bản ghi hiện tại và tất cả descendants (con cháu) của bản ghi đó.
- Upload/Thay ảnh: Có thể gọi API upload riêng hoặc truyền Base64 tùy thiết kế API backend.

### 5.2. Form Validation
- Tương tự Create (Tên bắt buộc, Slug định dạng chữ thường/số/dấu `-`, Ảnh chuẩn format/size).

## 6. Các tính năng nâng cao (Nghiệp vụ cốt lõi)
- **Cảnh báo Unsaved Changes:** Khai thác thuộc tính `isDirty` của `react-hook-form`. Gắn event báo khi người dùng nhấn Hủy, Back hoặc đóng tab/page mà chưa lưu (dùng hook `useWarnAboutChange` của Refine nếu có).
- **Logic Slug:** Khác với Create, màn hình Edit **không được** tự động cập nhật Slug theo Tên (để bảo toàn SEO/URL cũ). Chỉ thay đổi nếu người dùng cố tình gõ vào field Slug.
- **Vòng lặp danh mục (Circular Dependency):** Filter dữ liệu dropdown Parent ID: ID của danh mục được chọn làm Parent không được trùng với ID bản ghi đang sửa, hoặc nằm trong list descendant IDs.
- **Cảnh báo khi Ẩn (Inactive):** Lắng nghe sự thay đổi của field "Trạng thái". Nếu giá trị cũ = ACTIVE, giá trị mới = INACTIVE, gọi dialog (hoặc hiển thị hint text) cảnh báo: "Việc ẩn danh mục này có thể ảnh hưởng đến sản phẩm và danh mục con bên trong".

## 7. Checklist Triển khai
- [ ] Khởi tạo `CategoryEditPage` bọc `Edit` component.
- [ ] Tích hợp `useForm(action: "edit")`, đảm bảo dữ liệu map vào input chuẩn.
- [ ] Thiết lập component `SystemInfoCard` cho thông tin readonly.
- [ ] Điều chỉnh dropdown Danh mục cha (loại trừ chính nó và các con/cháu).
- [ ] Xử lý logic Slug cố định, chỉ cập nhật khi gõ thủ công.
- [ ] Cấu hình cảnh báo "Unsaved Changes" qua `isDirty`.
- [ ] Xử lý confirm dialog khi đổi trạng thái sang "Ẩn".
- [ ] Kiểm tra kết nối API `PUT /api/v1/admin/categories/{id}`, toast thành công và redirect về danh sách.
