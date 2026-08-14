# FRONTEND PLAN: Thêm mới bài viết (Create Post)

## 1. Mục tiêu
- Xây dựng màn hình Thêm mới bài viết (Post) dựa trên tài liệu ý tưởng `15-post-create-idea.md`.
- Sử dụng các hook của Refine.js (như `useForm`) và các component Material UI (MUI) kết hợp `react-hook-form` để quản lý state và validation.
- Xây dựng layout gọn gàng, tập trung không gian lớn cho phần soạn thảo nội dung (Rich Text Editor). 
- Chú ý: `type_id` sẽ được tự động trích xuất từ URL (không hiển thị form để chọn) và dùng để phân loại bài viết.

## 2. Phạm vi chức năng
- **Thông tin cơ bản:** Tiêu đề bài viết (title), Đường dẫn (slug), Nội dung bài viết (content).
- **Thông tin ngầm:** `type_id` lấy từ URL (Query Param `?type_id=...` hoặc Path Param) để đính kèm vào payload khi tạo bài viết.
- **Hành động form:** Hủy bỏ (quay về danh sách), Lưu bài viết.
- **Lưu thành công:** Redirect về trang danh sách tin tương ứng với type đang tương tác

## 3. Luồng màn hình (Main Content)

### 3.1. Header & Page Title
- **Tiêu đề:** `Thêm bài viết mới`
- **Breadcrumb:** `Bảng điều khiển > Bài viết > Thêm mới`

### 3.2. Form Layout (Sử dụng CSS Grid / Flexbox)
Bố cục có thể hiển thị dưới dạng 1 cột hoặc 2 cột tùy theo không gian, ưu tiên 1 cột tập trung (do không có thiết lập phụ phức tạp):

- **Card "Thông tin cơ bản":**
  - **Tiêu đề (title) (*):** Bắt buộc nhập, MUI `TextField`. Khi gõ sẽ tự động tạo chuỗi cho Slug.
  - **Đường dẫn (slug) (*):** Bắt buộc nhập, sinh tự động từ Tiêu đề, MUI `TextField`. Có thể cho phép người dùng tùy chỉnh lại.
  - **Nội dung (content) (*):** Bắt buộc nhập. Tích hợp Rich Text Editor (như Mui-Tiptap, Quill, hoặc TinyMCE) chiếm không gian rộng để dễ soạn thảo.

### 3.3. Sticky Footer Actions (Cố định ở dưới hoặc trên cùng)
- Nút **Hủy bỏ:** Hủy bỏ việc tạo, hiển thị dialog confirm nếu form đã bị thay đổi, điều hướng về trang danh sách bài viết.
- Nút **Lưu bài viết (Primary):** Trigger form submit.
- Các button có style đồng bộ với các trang đã triển khai.

## 4. Component Breakdown
- `PostCreatePage`: Component bọc màn hình chính, sử dụng thẻ `<Create>` từ `@refinedev/mui` và dùng hook `useForm`. Xử lý logic trích xuất `type_id` từ URL tại đây.
- `PostBasicInfoCard`: Component chứa các field Tiêu đề, Đường dẫn, và Trình soạn thảo văn bản.

## 5. Data & State Management

### 5.1. Refine Hooks & Form
- Sử dụng `@refinedev/react-hook-form` với hook `useForm` tương tác với resource `posts`.
- Sử dụng Hook của React Router (ví dụ: `useLocation`, `useSearchParams` hoặc `useParams`) để parse URL lấy giá trị `type_id`.
- Trong sự kiện `onFinish` hoặc bằng cách override lại hàm `submit` của react-hook-form, tiến hành nhồi thêm trường `type_id` vào payload chuẩn bị gửi đi.

### 5.2. Form Validation (Quy tắc kiểm tra)
- **Tiêu đề (title):** `required` (Bắt buộc). Độ dài tối đa 255.
- **Đường dẫn (slug):** `required`. Phải đúng định dạng URL thân thiện (không dấu, cách nhau bởi gạch ngang).
- **Nội dung (content):** `required` (Bắt buộc).
- **Loại bài viết (type_id):** Kiểm tra ngay lúc khởi tạo component, nếu `type_id` không tồn tại hoặc không hợp lệ -> Hiển thị cảnh báo và điều hướng (redirect) người dùng về lại trang danh sách, chặn hiển thị form tạo.

## 6. Các tính năng nâng cao (Nghiệp vụ cốt lõi)
- **Tự động sinh Slug:** Sử dụng event `onChange` trên field Tiêu đề, kết hợp với thư viện xử lý chuỗi (như `slugify`) để tự động điền giá trị vào field Slug.
- **Real-time Validation:** Khai báo chế độ `mode: "onBlur"` hoặc `"onChange"` trong `react-hook-form` để hiển thị lỗi ngay khi user nhập sai.
- **Unsaved Changes Warning:** Bật tính năng `warnWhenUnsavedChanges` của Refine để kích hoạt một popup cảnh báo nếu người dùng đang nhập dở nội dung bài viết nhưng bấm link chuyển sang màn hình khác.
- **Bảo mật URL Param:** Cần ép kiểu (Type Casting) an toàn cho `type_id` khi đọc từ URL trước khi đưa vào payload (tránh lỗi định dạng dữ liệu).

## 7. Checklist Triển khai
- [ ] Khởi tạo page component `PostCreatePage`.
- [ ] Xử lý logic lấy `type_id` từ URL (Query Param hoặc Path Param), validate sớm và báo lỗi/redirect nếu thiếu.
- [ ] Setup `useForm` (react-hook-form + Refine).
- [ ] Bố cục giao diện Card "Thông tin cơ bản".
- [ ] Cài đặt hoặc tái sử dụng Rich Text Editor Component cho trường "Nội dung".
- [ ] Viết logic tự động sinh Slug từ Tiêu đề.
- [ ] Định nghĩa Validation Schema (Ràng buộc tiêu đề, slug, nội dung).
- [ ] Ghi đè hàm submit để chèn `type_id` vào payload JSON trước khi gọi API.
- [ ] Tích hợp popup cảnh báo thoát form khi đang soạn dở (`warnWhenUnsavedChanges`).
- [ ] Kiểm tra submit luồng tạo và hiển thị thông báo toast thành công/thất bại.
