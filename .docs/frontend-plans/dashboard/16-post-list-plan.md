# FRONTEND PLAN: Quản trị danh sách bài viết (List Post)

## 1. Mục tiêu
- Xây dựng màn hình Quản lý danh sách bài viết (Post) dựa trên tài liệu ý tưởng `16-post-list-idea.md`.
- Sử dụng các hook của Refine.js (`useDataGrid`, `useDelete`) và component `<DataGrid>` của Material UI (MUI X) để hiển thị dữ liệu bảng tối ưu nhất.
- Đảm bảo luồng trải nghiệm đồng bộ với các màn hình danh sách đã có (như Product List, Category List), trong đó đặc biệt chú trọng tích hợp tự động lọc theo `typeCode` từ URL.

## 2. Phạm vi chức năng
- **Hiển thị danh sách:** Hiển thị bài viết trong DataGrid (bảng dữ liệu) có mật độ (density) cao/vừa để dễ quan sát.
- **Tìm kiếm (Search):** Tìm kiếm bài viết theo tiêu đề với cơ chế debounce 300ms bắt buộc.
- **Lọc ngầm (Hidden Filter):** Chỉ hiển thị các bài viết thuộc về `typeCode` hiện tại (lấy từ URL path hoặc query params).
- **Thao tác nhanh:** Xem chi tiết (Show), Chỉnh sửa (Edit) bài viết (chuyển hướng sang trang tương ứng), và Xóa (Delete) bài viết (có xác nhận Dialog).
- **Phân trang (Pagination):** Quản lý phân trang Server-side đồng bộ tự động với DataGrid của Refine.

## 3. Luồng màn hình (Main Content)

### 3.1. Header & Page Title
- Sử dụng thẻ `<List>` của `@refinedev/mui`.
- **Tiêu đề & Breadcrumb:** Thay đổi linh hoạt dựa trên `typeCode` đang xem (ví dụ: "Danh sách Tin tức", "Danh sách Trang").
- **Nút hành động:** Nút "Tạo mới" (Create) điều hướng sang trang thêm bài viết ứng với `typeCode` hiện tại.

### 3.2. Thanh công cụ tìm kiếm và lọc (Toolbar)
- Nằm phía trên DataGrid.
- Chứa một thẻ `TextField` (MUI) để tìm kiếm theo tiêu đề. Cần gắn sự kiện onChange và bọc bằng hàm `debounce` để tránh spam request.
- Các giá trị tìm kiếm và lọc sẽ được đồng bộ lên URL tự động nhờ Refine.js.

### 3.3. Bảng dữ liệu (DataGrid)
- Các cột dữ liệu chính:
  - **ID:** Định danh.
  - **Tiêu đề (Title):** Cột chính hiển thị tên bài viết.
  - **Đường dẫn (Slug):** Hiển thị rút gọn hoặc ellipsis.
  - **Tác giả (Author):** ID hoặc tên người viết.
  - **Ngày tạo (Created At):** Định dạng hiển thị ngày/tháng/năm dễ đọc.
  - **Hành động (Actions):** Sử dụng các component chuẩn của Refine như `<EditButton>`, `<ShowButton>`, `<DeleteButton>` cho từng hàng.

## 4. Component Breakdown
- `PostListPage`: Component chính bọc toàn trang, khởi tạo `useDataGrid`.
- `PostListToolbar`: Component chứa thanh tìm kiếm, filter.
- (Tái sử dụng) DataGrid của MUI X cùng hệ thống cột cấu hình linh hoạt.

## 5. Data & State Management

### 5.1. Refine Hooks & API
- Dùng `@refinedev/mui` với hook `useDataGrid` trỏ vào resource `posts`.
- Trích xuất tham số `typeCode` từ URL (thông qua React Router hook như `useParams` hoặc Refine routing).
- Truyền tham số `typeCode` vào mảng `filters` của `useDataGrid` ở chế độ ngầm định (permanent hoặc initial filter) để gọi API giới hạn kết quả trả về đúng loại bài viết đó.

### 5.2. Search State (Debounce)
- Quản lý state của input tìm kiếm cục bộ, kết hợp sử dụng lodash `debounce` hoặc hook tuỳ chỉnh trì hoãn 300ms.
- Khi giá trị debounced thay đổi, kích hoạt hàm `setFilters` của `useDataGrid` để cập nhật query lên Server.

## 6. Các tính năng nâng cao (Nghiệp vụ cốt lõi)
- **Cơ chế lọc tự động theo Type:** Tính năng này buộc Front-end phải chặn việc load toàn bộ `posts`. API chỉ được gọi khi đã xác định được `typeCode`. Nếu không có `typeCode`, tự động báo lỗi hoặc redirect.
- **Server-side Pagination & Sorting:** Tất cả các thao tác chuyển trang, bấm tiêu đề cột để sắp xếp trên DataGrid phải được truyền thẳng xuống Server thông qua query params của URL. Refine.js quản lý tự động việc này.
- **Bảo vệ hành động Xóa:** Tận dụng `<DeleteButton>` của Refine đã có sẵn cơ chế Confirmation Modal hoặc dùng Dialog tùy chỉnh của MUI. Sau khi xóa thành công hiển thị Toast notification.

## 7. Checklist Triển khai
- [ ] Khởi tạo page component `PostListPage`.
- [ ] Trích xuất `typeCode` từ URL và validate tồn tại.
- [ ] Thiết lập `useDataGrid` với tham số `filters` ngầm định áp dụng cho `typeCode`.
- [ ] Xây dựng UI thanh công cụ tìm kiếm và tích hợp hook debounce (300ms) để gọi API.
- [ ] Khai báo cấu trúc cột (Columns) cho bảng dữ liệu (ID, Tiêu đề, Slug, Ngày tạo).
- [ ] Tích hợp cột Hành động (Action) sử dụng `<EditButton>`, `<ShowButton>`, `<DeleteButton>`.
- [ ] Tinh chỉnh cài đặt của DataGrid (pagination, density).
- [ ] Kiểm tra đồng bộ URL (bộ lọc, phân trang, từ khóa tìm kiếm) khi load lại trang.
- [ ] Xử lý Loading state và Empty state (dùng các tuỳ chọn mặc định của MUI Grid hoặc custom component).
