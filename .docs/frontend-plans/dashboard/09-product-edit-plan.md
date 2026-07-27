# FRONTEND PLAN: Chỉnh sửa sản phẩm (Edit Product)

## 1. Mục tiêu
- Xây dựng màn hình Chỉnh sửa sản phẩm dựa trên tài liệu ý tưởng `09-product-edit-idea.md`.
- Kế thừa và tái sử dụng cấu trúc component từ màn hình Tạo mới (Product Create).
- Tích hợp hook `useForm` của Refine.js để lấy dữ liệu (fetch data) hiện tại của sản phẩm và cập nhật (mutate) dữ liệu lên Server.

## 2. Phạm vi chức năng
- Hiển thị đầy đủ thông tin sản phẩm đã có vào các trường của form (Prefill data).
- Cho phép chỉnh sửa: Tên sản phẩm, Mô tả chi tiết, Giá bán gốc, Giá khuyến mãi, SKU, Số lượng, Danh mục, Hình ảnh, Trạng thái.
- Ràng buộc Validation như màn hình Tạo mới.
- Nút hành động: Hủy bỏ, Lưu & Ẩn (Lưu trạng thái HIDDEN), Lưu & Xuất bản (Lưu trạng thái ACTIVE).

## 3. Luồng màn hình (Main Content)

### 3.1. Header & Page Title
- **Tiêu đề:** `Chỉnh sửa sản phẩm`
- **Breadcrumb:** `Bảng điều khiển > Sản phẩm > Chỉnh sửa`
- Component bọc ngoài: Dùng `<Edit>` từ `@refinedev/mui`.

### 3.2. Form Layout (Giữ nguyên cấu trúc Grid 2 cột như Tạo mới)
- **Cột chính (Main Column):**
  - Card "Thông tin cơ bản" (Tên, Mô tả - Rich Text Editor)
  - Card "Giá cả" (Giá gốc, Giá khuyến mãi)
  - Card "Quản lý kho" (SKU, Tồn kho)
- **Cột phụ (Side Column):**
  - Card "Hình ảnh" (Dropzone Upload)
  - Card "Phân loại" (Dropdown Category ID)
  - Card "Trạng thái" (Switch / Radio)

### 3.3. Sticky Footer Actions
- Kế thừa logic UI từ màn hình tạo mới, ghi đè nút "Lưu" (SaveButton) mặc định của Refine.
- Tùy biến hàm submit của `useForm` để mutate trường `status` tùy theo nút click.

## 4. Component Breakdown
- `ProductEditPage`: Component bọc màn hình chính, sử dụng thẻ `<Edit>` và `useForm`.
- Tái sử dụng các thẻ Card Component từ màn hình Create (nếu đã tách Component) hoặc sao chép layout cấu trúc.

## 5. Data & State Management

### 5.1. Refine Hooks & Form
- Sử dụng `@refinedev/react-hook-form` với hook `useForm`. Hook này sẽ tự động gọi API `GET /products/:id` dựa trên URL params (VD: `/products/edit/1`) để lấy dữ liệu prefill vào form.
- Khi submit, form tự động gọi API `PATCH/PUT /products/:id`.

### 5.2. Xử lý Image/File trong chế độ Edit
- Đảm bảo hiển thị trước (preview) của ảnh cũ (lấy URL từ `queryResult` của `useForm`).
- Khi user chọn file mới, upload file lên server, lấy URL mới và cập nhật field `image` trong form. Nếu user không chọn file mới, giữ nguyên URL cũ.

## 6. Checklist Triển khai
- [ ] Khởi tạo page component `ProductEditPage` bằng `<Edit>`.
- [ ] Tích hợp `useForm` của `react-hook-form`.
- [ ] Thiết lập dữ liệu mặc định (prefill) cho các field từ API, đặc biệt lưu ý hiển thị hình ảnh có sẵn.
- [ ] Render 2 cột layout và các form input (Autocomplete cho Category, Dropzone cho ảnh).
- [ ] Thêm custom validation check Giá khuyến mãi < Giá bán.
- [ ] Xử lý logic Unsaved Changes (`warnWhenUnsavedChanges: true`).
- [ ] Test luồng edit: Cập nhật thành công, quay về màn List hoặc giữ ở màn Edit.
