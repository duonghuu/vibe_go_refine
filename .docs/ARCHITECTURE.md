# KIẾN TRÚC HỆ THỐNG: E-COMMERCE PLATFORM

## 1. Sơ đồ dữ liệu cốt lõi (Core Database Entities)
- **Product (Sản phẩm):** Lưu thông tin cơ bản, `price` (giá gốc), `salePrice` (giá khuyến mãi), `stock` (tồn kho).
- **Category (Danh mục):** Phân loại sản phẩm (Quần áo, Điện tử...). Hỗ trợ cấu trúc cây (Parent - Child).
- **Cart & CartItem (Giỏ hàng):** Lưu trạng thái giỏ hàng tạm thời. Nếu user chưa login, dùng Session/Local Storage. Nếu đã login, lưu vào DB.
- **Order & OrderItem (Đơn hàng):** Lưu trữ lịch sử mua hàng bất biến (Snapshot). Giá sản phẩm trong OrderItem BẮT BUỘC lưu cứng tại thời điểm mua, không tham chiếu lại giá bảng Product.

## 2. Luồng nghiệp vụ tối quan trọng (Critical Business Logic)
- **Tính toán tiền (Pricing):** - Mọi phép tính tiền tệ (Tổng đơn, Giảm giá, Thuế) BẮT BUỘC thực hiện ở Backend (Go).
  - TUYỆT ĐỐI không tin tưởng dữ liệu giá tiền gửi lên từ Frontend. Frontend chỉ gửi `productId` và `quantity`.
- **Quản lý Tồn kho (Inventory):**
  - Trừ `stock` ngay khi thanh toán thành công.
  - Nếu `stock === 0`, API không cho phép thêm vào giỏ hàng.
- **Thanh toán (Payment Checkout):**
  - Trạng thái Order mặc định là `PENDING`.
  - Chỉ chuyển sang `PAID` qua Webhook trả về từ Cổng thanh toán (Stripe/VNPay).

## 3. Quy chuẩn API & Tối ưu hóa Dashboard (API & Datagrid Standards)
- **Phân trang (Pagination):** Mọi API trả về danh sách (Ví dụ: Lấy danh sách sản phẩm) bắt buộc phải có phân trang theo cấu trúc phẳng chuẩn tương thích với Refine.js Data Provider bao gồm mảng dữ liệu (`data`) và tổng số bản ghi (`total`).
- **Tối ưu hóa Truy vấn & Hiệu năng UI:** - Mọi ô tìm kiếm, auto-complete hoặc bộ lọc trên giao diện Dashboard bắt buộc phải tích hợp cơ chế hoãn hàm (`debounce 300ms`) từ phía Client trước khi gọi xuống Server nhằm tránh overload cơ sở dữ liệu.
  - Bảng dữ liệu Grid (`<DataGrid>` của Material UI) trong trang quản trị bắt buộc bật cấu hình mật độ hiển thị cao (`dense` hoặc `medium`) để tối ưu hóa không gian hiển thị số liệu.
- **Bảo mật:** API tạo Đơn hàng và API xuất báo cáo tài chính/vận hành (Export Excel/CSV) nặng bắt buộc phải có Rate Limit để chống Spam tài nguyên hệ thống.

## 4. Kiến trúc Trang Quản trị (Dashboard System Architecture)
- **Kiến trúc FE (Refine.js & Material UI):** - Tận dụng kiến trúc Resource-driven của Refine.js, ánh xạ trực tiếp các thực thể dữ liệu (`Product`, `Category`, `Order`, `User`) thành các Resource quản trị tương ứng thông qua các React hooks chuẩn (`useTable`, `useForm`, `useShow`).
  - Phân tách tuyệt đối giữa Smart/Container Component (Xử lý hook dữ liệu Refine, trạng thái Loading/Error/Empty) và Dumb/UI Component (Chỉ nhận props hiển thị và render UI theo Material UI Design System).
- **Kiến trúc BE (Golang - Gin - GORM - Wire):**
  - **Tầng Handler:** Chỉ bind dữ liệu Request JSON/Query từ Dashboard gửi lên, gọi UseCase và trả về cấu trúc JSON chuẩn HTTP.
  - **Tầng UseCase/Service:** Xử lý toàn bộ logic tính toán số liệu tổng quan (Metrics Aggregation), lọc dữ liệu nâng cao (Advanced Filtering) và xử lý phân quyền.
  - **Tầng Repository:** Thực hiện truy vấn dữ liệu tối ưu qua GORM. Toàn bộ các dependency được đăng ký tập trung tại file `wire.go`.

## 5. Module Auth (Xác thực, phân quyền & bảo mật)

### 5.1. Tổng quan nghiệp vụ
- Hệ thống sử dụng cơ chế **JWT** kết hợp với **Redis** để quản lý phiên đăng nhập theo thời gian thực
- Việc quản lý Token qua Redis nhằm mục đích: Thu hồi quyền truy cập ngay lập tức khi phát hiện rủi ro và chặn các phiên đăng nhập khi người dùng đã Logout

### 5.2. Định nghĩa vai trò

- `ADMIN`: Toàn quyền quản trị hệ thống. Có quyền truy cập mọi Dashboard cấu hình cấp cao, xem báo cáo tổng thể, duyệt User mới, quản lý cấu hình hệ thống và thu hồi quyền truy cập của bất kỳ tài khoản nào.
- `STAFF`: Nhân viên quản lý vận hành. Chỉ cho phép truy cập và thực thi một số Resource nhất định được quy hoạch trong Refine.js Layout (Ví dụ: Tạo order, duyệt trạng thái đơn hàng, cập nhật stock sản phẩm, xem lịch sử order, không có quyền xóa dữ liệu cốt lõi hoặc cấu hình hệ thống).
- `CUSTOMER`: Khách hàng nội bộ. Chỉ được xem order của chính mình, đặt hàng, quản lý sản phẩm trong giỏ hàng của chính mình.

### 5.3. Chiến lược Refresh Token Rotation & Quản lý JTI

- **Luồng chuẩn**: Khi Access Token (15 phút) hết hạn, Refine AuthProvider tự động gọi API `/refresh-token` gửi kèm Refresh Token nằm trong HttpOnly Cookie. Golang Backend kiểm tra tính hợp lệ, hủy định danh `jti` cũ trên Redis, sinh cặp Token mới (`Refresh Token B` + `Access Token New`), đồng thời lưu `jti` mới vào Redis Server.
- **Kịch bản tấn công Replay Attack**:
+ Nếu hacker trộm được `REFRESH TOKEN A` và cố gắng sử dụng lại nó (Trong khi Client thực đã đổi sang `REFRESH TOKEN B`)
+ Hệ thống kiểm tra Redis thấy `REFRESH TOKEN A` không tồn tại hoặc đã bị đánh dấu sử dụng hoặc thu hồi
+ **Hành động bắt buộc**: Phát tín hiệu CẢNH BÁO BẢO MẬT. Ngay lập tức xóa TOÀN BỘ Refresh Token của `userId` trong Redis để ép tất cả thiết bị phải đăng nhập lại từ đầu

### 5.4. Chiến lược đăng xuất

- Khi Admin/Staff bấm nút đăng xuất trên giao diện Refine.js, request `/logout` bắt buộc đi qua Middleware xác thực token.
- Hệ thống Backend tiến hành xóa mã `jti` của Refresh Token tương ứng khỏi Redis Server.
- Đồng thời, lấy Access Token hiện tại, tính toán thời gian tồn tại còn lại (TTL) và đưa token này vào danh sách **Blacklist** trên Redis thông qua hàm `client.Set(ctx, accessToken, "blacklisted", expiration)` để tự động giải phóng bộ nhớ, chặn đứng request ngay tại dòng code đầu tiên của Middleware bảo vệ.