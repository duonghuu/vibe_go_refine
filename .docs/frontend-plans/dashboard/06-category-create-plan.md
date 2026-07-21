# FRONTEND PLAN: Thêm danh mục sản phẩm (Create Category)

## 1. Mục tiêu
- Xây dựng màn hình (hoặc Drawer) Thêm mới danh mục sản phẩm (Create Category) dựa trên tài liệu ý tưởng `06-category-create-idea.md`.
- Sử dụng các hook của Refine.js (như `useForm` hoặc `useDrawerForm`) và các component Material UI (MUI) kết hợp `react-hook-form` để quản lý state và validation.
- Thiết kế giao diện rõ ràng để nhập thông tin cơ bản, chọn danh mục cha theo cấu trúc phân cấp, tải lên ảnh đại diện và thiết lập các thông số khác.

## 2. Phạm vi chức năng
- **Thông tin cơ bản:** Tên danh mục, Slug (tự động sinh từ Tên), Mô tả.
- **Danh mục cha:** Chọn từ danh sách phân cấp (Tree-view Dropdown). Cho phép chọn "Không có" để tạo danh mục gốc.
- **Hình ảnh:** Khu vực upload ảnh với tính năng Preview (xem trước), Thay đổi và Xóa ảnh.
- **Thiết lập:** Trạng thái (Hoạt động / Ẩn) và Thứ tự hiển thị.
- **Hành động form:** Hủy, Lưu (lưu và quay lại danh sách), Lưu & Thêm mới (lưu và reset form để tạo tiếp).

## 3. Luồng màn hình (Main Content)

### 3.1. Header & Page Title
- **Tiêu đề:** `Thêm danh mục` (sử dụng style tiêu đề chuẩn của dự án).
- **Breadcrumb:** `Trang chủ > Cấu hình hệ thống > Danh mục > Thêm mới`.

### 3.2. Form Layout (Sử dụng CSS Grid / Flexbox)
Chia layout thành 2 cột (ví dụ cột trái chiếm 2/3, cột phải chiếm 1/3) hoặc xếp dọc nếu dùng dạng Drawer:
- **Cột chính (Main Column):**
  - Card "Thông tin cơ bản":
    - **Tên danh mục (*):** Bắt buộc nhập, MUI `TextField`.
    - **Slug:** MUI `TextField`. Mặc định tự sinh ra từ trường Tên. Có thể chỉnh sửa.
    - **Danh mục cha:** Custom `Select` hoặc Tree Select hiển thị dạng phân cấp.
    - **Mô tả:** MUI `TextField` multiline (khoảng 4-5 rows).
- **Cột phụ (Sidebar Column):**
  - Card "Hình ảnh":
    - Khu vực kéo thả / chọn file. 
    - Validate dung lượng (≤ 5MB) và định dạng (jpg, jpeg, png, webp).
    - Hiển thị ảnh xem trước khi chọn.
  - Card "Thiết lập":
    - **Trạng thái:** MUI `RadioGroup` hoặc `Select` (Mặc định: Hoạt động).
    - **Thứ tự hiển thị:** MUI `TextField` dạng number (Mặc định: 0, tối thiểu: 0).

### 3.3. Footer Actions (Dưới cùng form)
- Nút **Hủy:** Xóa dữ liệu chưa lưu và đóng / điều hướng về trang danh sách.
- Nút **Lưu & Thêm mới:** `Button` (MUI) gọi API tạo, sau đó báo thành công và gọi `reset()` để clear form.
- Nút **Lưu:** `Button` (MUI) ưu tiên cao nhất, lưu thành công thì điều hướng về list.

## 4. Component Breakdown
- `CategoryCreatePage` (hoặc `CategoryCreateDrawer`): Thành phần chính bọc `Create` của Refine và sử dụng `useForm`.
- `CategoryBasicInfoCard`: Component bọc các trường Tên, Slug, Danh mục cha, Mô tả.
- `CategoryMediaCard`: Component tích hợp trình Upload ảnh (tái sử dụng component `ImageUploader` nếu có).
- `CategorySettingsCard`: Component cho cài đặt trạng thái, thứ tự.
- `ParentCategoryTreeSelect`: Dropdown chuyên dụng dùng `useSelect` gọi API danh mục và hiển thị dạng cây.

## 5. Data & State Management
### 5.1. Refine Hooks & Form
- Sử dụng `@refinedev/react-hook-form` với hook `useForm` (hoặc `useDrawerForm`) trên resource `categories`.
- Lấy danh sách category cho Dropdown bằng hook `useSelect` hoặc custom query.
- Xử lý upload ảnh: Gọi API POST `/api/v1/media/upload` (thông thường qua custom upload logic) trước khi submit form, hoặc dùng base64 nếu API backend hỗ trợ.

### 5.2. Form Validation (Yup / Zod / React Hook Form rules)
- **Tên danh mục:** `required` (Bắt buộc).
- **Slug:** Bắt buộc (nếu có regex: chỉ chứa chữ thường, số, dấu `-`).
- **Thứ tự:** `min: 0`.
- **Hình ảnh:** Tự custom logic báo lỗi file size (> 5MB) và mimetype.

## 6. Các tính năng nâng cao (Nghiệp vụ cốt lõi)
- **Auto-Generate Slug:** Theo dõi sự thay đổi của field "Tên danh mục" bằng `watch("name")`. Dùng `useEffect` kết hợp hàm `slugify` (có xử lý tiếng Việt) để tự động điền vào field "slug" nếu người dùng chưa can thiệp sửa tay slug.
- **Xử lý Save & New:** Ghi đè hàm submit mặc định của `useForm`. Nếu người dùng click "Lưu & Thêm mới", gắn flag. Trong callback `onMutationSuccess`, nếu flag này là true, ta gọi hàm `reset()` thay vì Refine tự động redirect về danh sách.

## 7. Checklist Triển khai
- [ ] Khởi tạo UI `CategoryCreatePage` chia cột theo Card.
- [ ] Cấu hình `useForm` và kết nối với MUI form.
- [ ] Cấu hình schema validation (Tên, Slug, Thứ tự...).
- [ ] Làm logic Auto-Slug (debounce).
- [ ] Tích hợp API lấy danh sách danh mục cha vào select/dropdown.
- [ ] Tạo component Upload ảnh đạt chuẩn (preview, error text dung lượng).
- [ ] Cấu hình nút Submit xử lý 2 trường hợp: Lưu thường và Lưu & Thêm mới.
- [ ] Kiểm tra tích hợp tạo mới toàn luồng, toast thông báo thành công từ Refine.
