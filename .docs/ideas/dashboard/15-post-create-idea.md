# Ý TƯỞNG: Quản trị Thêm mới bài viết (Admin Create Post)

## 1. Thông tin chung (Meta Info)

* **Dự án:** TechBite.
* **Tính năng:** Màn hình Thêm mới bài viết (Admin Create Post).
* **Mục đích:** Cho phép Admin/Editor nhập thông tin chi tiết để tạo mới một bài viết (post) vào hệ thống.

---

## 2. Đối tượng & Trải nghiệm (Target & UX)

* **Người dùng chính:** Quản trị viên (Admin), Biên tập viên (Editor), Người viết bài (Author).
* **Hành động chính:**
  * Nhập thông tin cơ bản: Tiêu đề (title), đường dẫn (slug), nội dung (content).
  * `type_id` được hệ thống tự động nhận diện thông qua URL (Path Params hoặc Query Params), không yêu cầu người dùng chọn.
  * Lưu bài viết hoặc Lưu nháp (nếu hệ thống hỗ trợ trạng thái).
* **Cảm xúc mang lại:** Form nhập liệu tối giản, sạch sẽ, cung cấp không gian rộng rãi cho trình soạn thảo văn bản (Editor). Thao tác lưu mượt mà, phản hồi ngay lập tức.

---

## 3. Đặc tả Thiết kế (Design Specs)

### Phong cách UI

* Bố cục dạng Grid/Card, chia thành 2 cột (Main Column cho nhập nội dung và Side Column cho thiết lập phụ) để tận dụng không gian Desktop.
* Thiết kế nhất quán: bo góc nhẹ (`rounded-xl` hoặc `rounded-2xl`), khoảng trắng hợp lý, Typography rõ ràng dễ đọc.
* Sticky Action Bar (Thanh công cụ cố định trên cùng hoặc dưới cùng) để người dùng luôn có thể nhấn "Lưu" bài viết kể cả khi cuộn chuột qua bài viết dài.
* Responsive: Layout 2 cột sẽ chuyển thành 1 cột trên màn hình Tablet/Mobile.

### Màu sắc chủ đạo (Brand Colors)

* **Primary:** Nút "Lưu", "Tạo mới": Dùng màu chủ đạo của hệ thống (VD: `bg-[#ff8c42]`).
* **Đỏ:** Báo lỗi nhập liệu (Validation Error): `text-red-500`.
* **Xám nhạt/Trắng:** Màu nền cho các Card chứa form.

### Cấu trúc màn hình (Top to Bottom)

#### Header & Sticky Actions
* Tiêu đề: "Thêm bài viết mới"
* Breadcrumb: Bảng điều khiển > Bài viết > Thêm mới
* Nút hành động:
  * Nút **Hủy bỏ** (Quay về danh sách)
  * Nút **Lưu bài viết** (Primary Action)

#### Bố cục 2 Cột (2-Column Layout)

**Cột chính (Main Column - Chiếm khoảng 2/3 không gian):**

1. **Thông tin cơ bản (Basic Info Card)**
   * **Tiêu đề (title):** (Input text - Bắt buộc). Tiêu đề của bài viết.
   * **Đường dẫn (slug):** (Input text - Bắt buộc). Đường dẫn thân thiện SEO (tự động sinh từ tiêu đề nếu để trống).
   * **Nội dung (content):** (Rich Text Editor). Trình soạn thảo văn bản chính của bài viết. Không gian hiển thị rộng rãi, tích hợp công cụ định dạng văn bản.

**Cột phụ (Side Column - Chiếm khoảng 1/3 không gian):**

*(Ghi chú: Nếu không có thiết lập phụ nào khác, form có thể dùng bố cục 1 cột)*

1. **Thông tin phân loại (Categorization)**
   * **Danh mục bài viết (category_id):** (Select / Dropdown - Không bắt buộc). Lấy danh sách từ bảng `post_categories` dựa theo `typeCode`. Hỗ trợ chọn danh mục cha/con. Có thể dùng Autocomplete nếu danh mục dài.

2. **Thông tin ngầm (Hidden Data)**
   * **Loại bài viết (type_id):** Không hiển thị chọn trên UI. Frontend sẽ tự động lấy giá trị này từ URL (ví dụ: `?type_code=NEWS` hoặc route parameter) để submit ngầm xuống API.

---

## 4. Kiểm tra dữ liệu (Validation & Error Handling)

* **Tiêu đề (title):** Bắt buộc nhập. Giới hạn độ dài tối đa phù hợp (ví dụ 255 ký tự).
* **Đường dẫn (slug):** Bắt buộc nhập. Phải duy nhất, định dạng chuỗi thân thiện với URL (không dấu, các từ cách nhau bởi dấu gạch ngang). Nếu người dùng không nhập, Frontend tự động sinh từ `title`.
* **Nội dung (content):** Bắt buộc nhập. Báo lỗi nếu để trống.
* **Loại bài viết (type_id):** Bắt buộc. Hệ thống kiểm tra giá trị từ URL (typeCode); nếu không có hoặc không hợp lệ, chặn submit và báo lỗi hoặc điều hướng về trang lỗi/danh sách.

---

## 5. Chức năng nổi bật

* **Tự động tạo Slug (Auto-slug):** Khi người dùng gõ nội dung vào ô "Tiêu đề", hệ thống hỗ trợ tự động điền giá trị chuyển đổi vào ô "Đường dẫn" (Slug).
* **Real-time Validation:** Báo lỗi ngay lập tức khi người dùng nhập sai định dạng hoặc bỏ trống trường bắt buộc và rời khỏi ô nhập đó (onBlur/onChange).
* **Unsaved Changes Warning:** Hiển thị cảnh báo (Confirm Dialog) nếu người dùng điều hướng sang trang khác khi đang soạn thảo bài viết mà chưa lưu.
* **Toast Notification:** Hiển thị thông báo (Snackbar) góc màn hình báo thao tác thành công hoặc chi tiết lỗi kết nối API.

---

## 6. Định hướng UI Component

* **Page Layout / Header:** Chứa Breadcrumb và nhóm nút thao tác ở vị trí cố định.
* **Card Component:** Khối giao diện chứa từng nhóm field (Basic Info, Organization).
* **Form Elements:**
  * `TextField` / `Input` (dành cho title, slug).
* **Rich Text Editor:** Sử dụng thư viện Editor (như Mui-Tiptap, Quill, hoặc TinyMCE) cho trường `content`.
* **Toast / Snackbar:** Component phản hồi trạng thái từ Material UI.
* **Dialog:** Hộp thoại cảnh báo.
