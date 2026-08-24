# FRONTEND PLAN: Chỉnh sửa danh mục bài viết (Edit Post Category)

## 1. Mục tiêu
- Xây dựng màn hình Chỉnh sửa danh mục bài viết (Edit Post Category) dựa trên tài liệu ý tưởng `18-post-category-edit-idea.md`.
- Giao diện đồng nhất với màn hình thêm danh mục, sử dụng các hook của Refine.js (`useForm`) kết hợp Material UI (MUI).
- Đảm bảo hiển thị đầy đủ thông tin cũ, loại trừ chính nó và các danh mục con trong danh sách "Danh mục cha" để tránh lỗi vòng lặp (circular reference).

## 2. Phạm vi chức năng
- **Loại bài viết (Post Type):** Hiển thị dạng Text/Read-only, không cho phép chỉnh sửa.
- **Thông tin cơ bản:** Tên danh mục, Slug, Mô tả, Danh mục cha (ngoại trừ nó và danh mục con của nó).
- **Hình ảnh:** Hiển thị ảnh cũ. Khu vực upload cho phép Đổi/Xóa ảnh.
- **Thiết lập:** Trạng thái (Hoạt động / Ẩn) và Thứ tự hiển thị.
- **Hành động form:** Hủy, Lưu (cập nhật và quay lại danh sách).

## 3. Luồng màn hình (Main Content)

### 3.1. Header & Page Title
- **Tiêu đề:** `Chỉnh sửa danh mục bài viết`
- **Breadcrumb:** `Trang chủ > Quản trị nội dung > Danh mục > Chỉnh sửa`

### 3.2. Form Layout (Cấu trúc Grid 2 cột)
Đồng nhất 100% về kiến trúc UI với trang tạo danh mục bài viết:
- **Cột chính (Main Column - chiếm 2/3 độ rộng):**
  - Card "Thông tin cơ bản":
    - **Loại bài viết (Post Type):** Field Read-only hoặc Text, thể hiện Loại bài viết đang thuộc về (lấy từ dữ liệu fetch được).
    - **Tên danh mục (*):** Bắt buộc, TextField.
    - **Slug (*):** TextField. Không tự động generate real-time để tránh hỏng SEO của URL cũ. Chỉ update khi người dùng cố tình thay đổi.
    - **Danh mục cha:** Tree Select / Dropdown phân cấp. Chỉ fetch và hiển thị dữ liệu theo `type_code` của danh mục. Cần lọc bỏ node hiện tại (bằng `id`) và các node con để không tự chọn làm cha của chính mình.
    - **Mô tả:** TextField đa dòng (multiline).
- **Cột phụ (Sidebar Column - chiếm 1/3 độ rộng):**
  - Card "Hình ảnh":
    - Kế thừa Component Upload đang có, truyền vào ảnh hiện tại từ dữ liệu cũ (nếu có).
  - Card "Thiết lập hiển thị":
    - **Thứ tự (Sort Order):** Number TextField.
    - **Trạng thái:** Switch hoặc Radio Group (ACTIVE / INACTIVE).

### 3.3. Footer Actions (Dưới cùng form)
- Nút **Hủy:** Xóa/Discard mọi thay đổi và chuyển hướng quay lại danh sách.
- Nút **Lưu:** Gọi API PUT `/api/v1/admin/post-categories/{id}`, hiển thị Toast thành công, điều hướng về danh sách. Bỏ nút "Lưu & Thêm mới" so với form Create.

## 4. Component Breakdown
- `PostCategoryEditPage`: Wrapper component chứa logic Refine `useForm` (action: "edit").
- `PostCategoryBasicInfoCard`: Chứa Loại bài viết, Tên, Slug, Danh mục cha, Mô tả. Tái sử dụng/Mở rộng từ form Create nếu có thể, thêm logic read-only cho Post Type.
- `PostCategoryMediaCard`: Quản lý hiển thị và thay đổi ảnh. Tái sử dụng component upload chung.
- `PostCategorySettingsCard`: Quản lý Trạng thái và Thứ tự.
- `ParentPostCategoryTreeSelect`: Select component custom dùng `useSelect` gọi API `/api/v1/admin/post-categories/tree?type_code={code}`. Thêm props `excludeId` để lọc bỏ trên UI nhánh danh mục hiện tại.

## 5. Data & State Management
### 5.1. Refine Hooks & API
- Khai báo form bằng `@refinedev/react-hook-form` với resource `post-categories` và action `edit`. Refine sẽ tự động gọi API `GET /{id}` để lấy dữ liệu đổ vào form.
- Dựa vào data GET được, xác định `type_code` để gọi API lấy cây danh mục cha.
- Khi submit, payload gửi PUT request cần bao gồm các thông tin thay đổi. Giao tiếp Upload qua Endpoint chung `/api/v1/media/upload` nếu có thay ảnh mới trước khi submit data form.

### 5.2. Form Validation (Rules)
- **Tên danh mục:** Bắt buộc, <= 255 ký tự.
- **Slug:** Bắt buộc, đúng Regex chuẩn slug `^[a-z0-9-]+$`.
- **Danh mục cha:** Có thể null/rỗng.
- **Thứ tự:** >= 0.
- **Hình ảnh:** Tương thích jpg/png/webp, size <= 5$MB.

## 6. Các tính năng nâng cao (Nghiệp vụ cốt lõi)
- **Chặn vòng lặp danh mục (Circular Reference):** Component Tree Select phải lọc dữ liệu trả về từ API, xóa node có ID trùng với danh mục đang edit (kéo theo xóa toàn bộ node con của nó trên tree) để đảm bảo không gán nhầm.
- **Auto-Generate Slug Disabled:** Ở màn hình Edit, tính năng gõ Tên tự động sinh Slug bị tắt đi hoặc yêu cầu user chủ động click nút "Tạo Slug" (nếu có thiết kế thêm) để tránh sửa nhầm URL đang SEO.

## 7. Checklist Triển khai
- [ ] Khởi tạo `PostCategoryEditPage` với Refine `useForm(action: 'edit')`.
- [ ] Dựng Layout chính chia cột Card giống form Create.
- [ ] Thiết lập Form Data mặc định (defaultValues) khi fetch xong GET API.
- [ ] Field "Loại bài viết" hiển thị Read-only.
- [ ] Xử lý lọc bỏ node hiện tại và con của nó trong component `ParentPostCategoryTreeSelect`.
- [ ] Tắt tính năng tự động sinh Slug khi Edit.
- [ ] Gắn Component Upload ảnh chung (xử lý hiển thị URL ảnh cũ).
- [ ] Code logic nút "Lưu" gửi PUT request và custom Toast error handling.
