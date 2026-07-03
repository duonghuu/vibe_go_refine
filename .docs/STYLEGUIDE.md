# STYLEGUIDE DESIGN SYSTEM: B2B ADMIN PLATFORM

## 1. Bảng màu & Phân cấp trực quan (Color & Hierarchy)
- **Primary Action (Màu chốt hạ/Tiến trình):** `palette.primary.main` (Màu xanh chủ đạo hệ thống. Sử dụng cho Sidebar Active, các nút hành động chính như "Thêm mới", "Tiếp tục", "Xem chi tiết").
- **Success Action (Màu đột biến dữ liệu thành công):** `palette.success.main` (Màu xanh lá. Chỉ dùng DUY NHẤT cho các nút xác nhận lưu thay đổi mang tính chất chốt hạ như: "Lưu cấu hình", "Kích hoạt Scenario", "Phát hành Coupon" để định hướng hành vi an toàn).
- **Secondary Action:** `palette.secondary.main` hoặc `palette.action.disabled` (Màu xám/slate. Sử dụng cho nút "Hủy", "Quay lại", các bộ lọc danh sách chưa active).
- **Background Hệ thống:** Nền layout tổng thể dùng màu xám siêu nhạt `palette.background.default` (`#f8fafc`). Toàn bộ các thẻ thống kê (Metrics Card) và Bảng dữ liệu (Data Grid) bắt buộc dùng nền trắng tinh (`#ffffff`) để tạo độ tương phản phẳng, giúp giảm mỏi mắt khi vận hành lâu.
- **Hiển thị Số liệu (Metrics):** Số liệu tổng quan dùng font lớn bold màu sẫm (`#0f172a`). Các chỉ số tăng trưởng (Growth Rate) nếu Dương dùng `text-green-600`, nếu Âm dùng `text-red-600`.

## 2. Thành phần đặc thù Admin (Admin UI Components)
- **Data/Metrics Card (Thẻ thống kê tổng quan):**
  - Khung viền mỏng phẳng (Flat Border) hoặc đổ bóng rất nhẹ (Elevation 1), không dùng hiệu ứng hovering lòe loẹt.
  - Phía trên bên phải luôn kèm theo Icon đại diện cho dữ liệu (User Icon, Coupon Icon, Scenario Diagram Icon).
- **Trạng thái Thực thể (Entity Status Chips):**
  - `Active / Đang chạy`: Dạng Chip nền xanh lá nhạt, chữ xanh đậm.
  - `Draft / Nháp`: Dạng Chip nền xám nhạt, chữ xám đậm.
  - `Paused / Tạm dừng`: Dạng Chip nền cam nhạt, chữ cam đậm.
- **Quy chuẩn Bảng dữ liệu (Data Grid / Table Layout):**
  - Cột cuối cùng bên phải luôn là cột "Hành động" (Actions) cố định (Pinned/Sticky Column) chứa cụm nút thao tác nhanh: Xem (Eye), Sửa (Edit), Xóa (Delete).
  - Độ cao dòng (Row Density): Sử dụng chế độ hiển thị `dense` mặc định của MUI để tối ưu mật độ thông tin trên một màn hình hiển thị lớn.

## 3. Trải nghiệm người dùng & Ràng buộc Form (UX & Form Constraints)
- **Luồng xử lý Biểu mẫu (Form Mutation):**
  - **Form ngắn (< 5 trường nhập liệu):** Sử dụng `Drawer` (Ngăn kéo trượt từ phải sang) thông qua hook `useDrawerForm` của Refine để tạo mới hoặc chỉnh sửa trực tiếp tại trang danh sách. KHÔNG chuyển trang để tránh làm mất dấu bộ lọc dữ liệu hiện tại của Admin.
  - **Form dài / Cấu hình kịch bản phức tạp:** Chuyển hướng trang riêng, bắt buộc áp dụng UI `Stepper` chia nhỏ tiến trình (Bước 1: Thông tin cơ bản -> Bước 2: Thiết lập điều kiện/Target nhóm user -> Bước 3: Nội dung Broadcast/Coupon -> Bước 4: Review & Kích hoạt).
- **Phản hồi hệ thống (Feedback & Alerts):**
  - Khi biểu mẫu Mutation thành công (Tạo/Sửa), hệ thống bắt buộc kích hoạt `Snackbar (Toast)` góc trên cùng bên phải báo thành công và tự đóng sau 3 giây.
  - Hành động bấm nút "Xóa" hoặc "Hủy kích hoạt" bắt buộc phải mở `Confirmation Dialog` (Hộp thoại xác nhận) trước khi gọi API thực tế.
- **Tối ưu hóa Tìm kiếm (Search & Filters):**
  - Mọi ô nhập liệu tìm kiếm tự động hoặc bộ lọc trên Header của Table bắt buộc phải tích hợp cơ chế `debounce 300ms` thông qua hook của Refine dữ liệu trước khi trigger fetch API xuống backend, cấm bắn request liên tục theo mỗi lượt bấm phím.