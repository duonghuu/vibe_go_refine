# FRONTEND PLAN: Thêm danh mục bài viết (Create Post Category)

## 1. Mục tiêu
- Xây dựng màn hình Thêm mới danh mục bài viết (Create Post Category) dựa trên tài liệu ý tưởng `17-post-category-create-idea.md`.
- Giao diện đồng nhất với màn hình thêm danh mục sản phẩm đang có, tuân thủ `STYLEGUIDE.md` và sử dụng các hook của Refine.js (`useForm`) kết hợp Material UI (MUI).
- Xử lý mượt mà nghiệp vụ: phân loại danh mục theo Loại bài viết (Post Type) được tự động lấy từ URL, đảm bảo cây danh mục cha luôn chính xác với loại bài viết đó.

## 2. Phạm vi chức năng
- **Loại bài viết (Post Type):** Hệ thống tự động trích xuất từ URL (Path hoặc Query Params). Đóng vai trò làm context bắt buộc.
- **Thông tin cơ bản:** Tên danh mục, Slug (tự động sinh từ Tên), Mô tả, Danh mục cha (chỉ hiển thị danh mục thuộc cùng Post Type).
- **Hình ảnh:** Khu vực upload ảnh đại diện có Preview, Đổi/Xóa ảnh.
- **Thiết lập:** Trạng thái (Hoạt động / Ẩn) và Thứ tự hiển thị.
- **Hành động form:** Hủy, Lưu (quay lại danh sách), Lưu & Thêm mới (reset form để tiếp tục tạo).

## 3. Luồng màn hình (Main Content)

### 3.1. Header & Page Title
- **Tiêu đề:** `Thêm danh mục bài viết`
- **Breadcrumb:** `Trang chủ > Quản trị nội dung > Danh mục > Thêm mới`

### 3.2. Form Layout (Cấu trúc Grid 2 cột)
Đồng nhất 100% về kiến trúc UI với trang thêm danh mục sản phẩm:
- **Cột chính (Main Column - chiếm 2/3 độ rộng):**
  - Card "Thông tin cơ bản":
    - **Tên danh mục (*):** Bắt buộc, TextField.
    - **Slug (*):** TextField. Tự động sinh từ Tên, cho phép chỉnh sửa thủ công.
    - **Danh mục cha:** Tree Select / Dropdown phân cấp. Chỉ fetch và hiển thị dữ liệu theo `type_code` lấy từ URL.
    - **Mô tả:** TextField đa dòng (multiline).
- **Cột phụ (Sidebar Column - chiếm 1/3 độ rộng):**
  - Card "Hình ảnh":
    - Kế thừa Component Upload đang có. Hỗ trợ validate kích thước $\le$ 5MB, chuẩn hình ảnh.
  - Card "Thiết lập hiển thị":
    - **Thứ tự (Sort Order):** Number TextField (mặc định là 0).
    - **Trạng thái:** Switch hoặc Radio Group (ACTIVE / INACTIVE).

### 3.3. Footer Actions (Dưới cùng form)
- Dùng màu chuẩn theo Styleguide (`palette.primary.main` và `palette.action.disabled`).
- Nút **Hủy:** Xóa dữ liệu form và chuyển hướng quay lại danh sách.
- Nút **Lưu & Thêm mới:** Gọi API POST, hiển thị Toast thành công, và gọi hàm reset form.
- Nút **Lưu:** Gọi API POST, hiển thị Toast thành công, điều hướng về danh sách.

## 4. Component Breakdown
- `PostCategoryCreatePage`: Wrapper component chứa logic Refine `useForm`.
- `PostCategoryBasicInfoCard`: Chứa Tên, Slug, Danh mục cha, Mô tả.
- `PostCategoryMediaCard`: Quản lý hiển thị và upload ảnh.
- `PostCategorySettingsCard`: Quản lý Trạng thái và Thứ tự.
- `ParentPostCategoryTreeSelect`: Select component custom dùng `useSelect` gọi API `/api/v1/admin/post-categories/tree?type_code={code}`.

## 5. Data & State Management
### 5.1. Refine Hooks & API
- Lấy `type_code` từ Refine routing qua `useParsed` hoặc custom hook đọc Query/Path params.
- Khởi tạo form bằng `@refinedev/react-hook-form` với resource `post-categories`. Khi submit, payload phải đính kèm `type_code`.
- Giao tiếp Upload qua Endpoint chung `/api/v1/media/upload` trước khi submit data.

### 5.2. Form Validation (Rules)
- **type_code:** Bắt buộc (kiểm tra ngầm).
- **Tên danh mục:** Bắt buộc, $\le$ 255 ký tự.
- **Slug:** Bắt buộc, đúng Regex chuẩn slug `^[a-z0-9-]+$`.
- **Thứ tự:** $\ge 0$.
- **Hình ảnh:** Tương thích jpg/png/webp, size $\le 5$MB.

## 6. Các tính năng nâng cao (Nghiệp vụ cốt lõi)
- **Ràng buộc danh mục cha theo URL Context:** Form chỉ load được cây danh mục cha khi `type_code` sẵn sàng. Tránh tình trạng gán nhầm parent khác Post Type.
- **Auto-Generate Slug:** Bắt sự kiện thay đổi của `name` qua `watch()`, debounce tạo slug tự động nếu trường `slug` chưa bị focus chỉnh sửa (isDirty = false).
- **Save & New Flow:** Can thiệp hàm `onMutationSuccess` của form để xác định cờ (flag) từ button submit, nếu là "Lưu & Thêm mới" thì ngăn chặn hành động redirect mặc định, thay vào đó gọi `reset()`.

## 7. Checklist Triển khai
- [ ] Lấy và kiểm tra `type_code` từ URL vào giao diện `Create`.
- [ ] Dựng Layout chính chia cột Card giống với CategoryCreatePage.
- [ ] Cấu hình hook `useForm` với bộ rule validation nghiêm ngặt.
- [ ] Xử lý Real-time Auto-Slug.
- [ ] Code component `ParentPostCategoryTreeSelect` fetch theo `type_code`.
- [ ] Gắn Component Upload ảnh chung.
- [ ] Code logic hai nút Submit ("Lưu", "Lưu & Thêm mới") và custom Toast error handling.
