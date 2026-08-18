# Ý TƯỞNG: Quản trị Danh sách bài viết (Admin List Post)

## 1. Thông tin chung (Meta Info)

* **Dự án:** TechBite.
* **Tính năng:** Màn hình Danh sách bài viết (Admin List Post).
* **Mục đích:** Cho phép Admin/Editor quản lý, xem, tìm kiếm, lọc và phân trang các bài viết đã được tạo trong hệ thống. Hỗ trợ truy cập nhanh vào các chức năng Thêm, Sửa, Xóa.

---

## 2. Đối tượng & Trải nghiệm (Target & UX)

* **Người dùng chính:** Quản trị viên (Admin), Biên tập viên (Editor), Người viết bài (Author).
* **Hành động chính:**
  * Xem danh sách bài viết dưới dạng bảng (DataGrid).
  * Phân loại, hiển thị bài viết tự động theo `type_code` (lấy từ URL path hoặc query params).
  * Tìm kiếm bài viết theo tiêu đề (có tích hợp debounce).
  * Phân trang (Pagination) và sắp xếp (Sorting).
  * Thực hiện các thao tác: Xem chi tiết, Chỉnh sửa, Xóa (có xác nhận).
* **Cảm xúc mang lại:** Giao diện quản lý trực quan, gọn gàng, dữ liệu hiển thị với mật độ tối ưu. Các thao tác lọc, tìm kiếm phản hồi nhanh, bảng dữ liệu tối ưu không gian dễ dàng bao quát thông tin.

---

## 3. Đặc tả Thiết kế (Design Specs)

### Phong cách UI

* Bảng dữ liệu (DataGrid) chiếm trọng tâm không gian, sử dụng thiết lập mật độ cao (dense/medium) phù hợp quản trị.
* Thiết kế phẳng, các nút thao tác (Action) trên từng hàng gọn gàng (ưu tiên dùng Icon thay vì Text dài).
* Bố cục rõ ràng, nhất quán: Header (Tiêu đề, Breadcrumb, Nút Tạo mới) -> Thanh công cụ (Tìm kiếm/Lọc) -> Bảng dữ liệu (DataGrid) -> Phân trang (Pagination).

### Màu sắc chủ đạo (Brand Colors)

* **Primary:** Nút "Tạo bài viết mới": Dùng màu chủ đạo của hệ thống (VD: `bg-[#ff8c42]`).
* **Đỏ:** Báo lỗi hoặc hành động nguy hiểm như Xóa bài viết (`text-red-500`).
* **Trắng/Xám nhạt:** Nền của các Card chứa bảng dữ liệu.

### Cấu trúc màn hình (Top to Bottom)

#### Header & Actions
* Tiêu đề: Tự động đổi theo loại bài viết dựa vào URL (ví dụ: "Danh sách Tin tức", "Danh sách Trang").
* Breadcrumb: Bảng điều khiển > Bài viết > Danh sách
* Nút hành động chính:
  * Nút **Tạo bài viết mới** (Primary Action - Chuyển hướng sang trang Tạo mới của loại bài viết tương ứng).

#### Thanh công cụ (Toolbar)
* **Ô tìm kiếm (Search Input):** Tìm kiếm theo tiêu đề bài viết. Bắt buộc có debounce 300ms.
* **Bộ lọc (Filter):** (Tùy chọn mở rộng sau này) Lọc theo trạng thái, khoảng thời gian tạo.

#### Bảng dữ liệu (DataGrid)
* **Cột hiển thị (Columns):**
  * **ID / STT:** Định danh hoặc số thứ tự.
  * **Tiêu đề (Title):** Tên bài viết (có thể nhấn vào để sang trang chỉnh sửa).
  * **Đường dẫn (Slug):** Hiển thị text rút gọn (ellipsis).
  * **Ngày tạo (Created At):** Định dạng ngày tháng năm dễ đọc.
  * **Hành động (Actions):** Chứa các nút (Icons): Chỉnh sửa (Edit), Xóa (Delete), Xem chi tiết (Show).

---

## 4. Kiểm tra dữ liệu (Validation & Error Handling)

* **Xóa bài viết:** Bắt buộc hiển thị hộp thoại xác nhận (Confirm Dialog) trước khi gọi API xóa thực sự để tránh thao tác nhầm.
* **Xử lý lỗi tải dữ liệu:** Hiển thị trạng thái Loading (Skeleton hoặc Spinner) khi đang fetch API. Nếu lỗi, báo lỗi (Toast/Snackbar).
* **Trạng thái rỗng (Empty State):** Khi không có bài viết nào, hiển thị giao diện báo trống kèm nút "Tạo bài viết mới" ngay giữa màn hình.

---

## 5. Chức năng nổi bật

* **Tự động lọc theo Type:** Tương tự trang Create, trang List tự động trích xuất mã loại bài viết từ URL (`typeCode`) để truyền vào bộ lọc ngầm xuống API, đảm bảo người dùng đang xem đúng danh sách của loại bài viết đó.
* **Tối ưu hóa hiệu năng (Debounce):** Tìm kiếm không gọi API liên tục mà đợi người dùng ngừng gõ 300ms, giảm tải cho backend (tuân thủ quy định tối ưu hóa kiến trúc).
* **Đồng bộ URL:** Mọi thông tin về bộ lọc, phân trang, sắp xếp tự động đồng bộ lên URL, giúp Admin dễ dàng tải lại trang hoặc chia sẻ liên kết mà không mất trạng thái hiện tại.
* **Phân trang phía Server (Server-side Pagination):** Tích hợp hoàn hảo với Refine.js Data Provider.

---

## 6. Định hướng UI Component

* **Page Layout / Header:** Component tiêu chuẩn cho trang List của Refine kết hợp Breadcrumb và Action button.
* **DataGrid Component:** Dùng `<DataGrid>` của thư viện MUI X để quản lý bảng dữ liệu với tính năng sorting/pagination tích hợp.
* **Search / Form Elements:** Dùng `TextField` của MUI cho thanh tìm kiếm.
* **Action Buttons:** `IconButton` và `Tooltip` của MUI cho các tác vụ trong cột Hành động.
* **Dialog:** Hộp thoại cảnh báo (Dialog của MUI) khi thực hiện tác vụ xóa.
