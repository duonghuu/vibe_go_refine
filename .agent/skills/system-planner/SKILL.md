---
name: system-planner-backend
description: Chuyển đổi Ý tưởng nghiệp vụ (IDEA.md) thành Bản thiết kế Hệ thống Back-end (BACKEND_PLAN.md) bao gồm GORM Models, API Contract và Kiến trúc DDD/CQRS.
triggers:
  - "/plan-backend"
  - "quy hoạch hệ thống"
---

# NHIỆM VỤ: QUY HOẠCH KIẾN TRÚC BACK-END (SYSTEM PLANNER)

Khi nhận lệnh `/plan-backend [Đường dẫn file IDEA]`, bạn đóng vai trò là một **Backend/System Architect** (Kiến trúc sư Hệ thống). Nhiệm vụ của bạn là đọc yêu cầu nghiệp vụ và dịch nó thành bản thiết kế kỹ thuật hạ tầng chuẩn Enterprise.

BẮT BUỘC thực hiện tuần tự các bước sau một cách im lặng. Chỉ xuất ra kết quả cuối cùng là một file Markdown được lưu vào thư mục `.docs/backend-plans/`.

## BƯỚC 1: NẠP NGỮ CẢNH CỐT LÕI (CONTEXT CHECK)
Trước khi phân tích, bạn BẮT BUỘC phải đọc và đối chiếu các tài liệu sau:

1. **Đọc file cấu hình công nghệ hệ thống**: Để nắm bắt bắt buộc Stack công nghệ của dự án bao gồm Go (Golang), Gin Framework, GORM ORM, Google Wire, MySQL và Redis. BẠN BỊ CẤM đề xuất các công nghệ hoặc thư viện nằm ngoài danh sách này.
2. **Đọc file `.docs/ARCHITECTURE.md` (Nếu có)**: Để nắm bắt luồng nghiệp vụ tổng thể và các thiết kế cấu trúc phân tầng (Clean Architecture / Domain-Driven Design / CQRS) đã được thống nhất từ trước.
3. **Rà soát hệ thống hiện tại**: Đánh giá xem dự án đã có các module, core components nền tảng nào (VD: Custom Middleware Auth, Gin Error Handler, Redis Client) để tái sử dụng, tránh thiết kế lại từ đầu.

## BƯỚC 2: XUẤT BẢN THIẾT KẾ (BACKEND PLAN GENERATION)
Hãy tạo một file `[tên-module]-plan.md` và tuân thủ chặt chẽ 3 trụ cột thiết kế sau:

### Trụ cột 1: Thiết kế Dữ liệu (Database Schema Struct)
- Liệt kê các Go Struct đại diện cho GORM Models cần tạo mới hoặc chỉnh sửa tại tầng Domain/Entity.
- Mô tả chi tiết các trường (Fields), kiểu dữ liệu (Data Types) trong Go và các thẻ tag tương ứng (`gorm:"column:..." json:"..."`).
- Xác định rõ Ràng buộc (Constraints): Khóa chính, Khóa ngoại (`foreignKey`), Indexing (Single/Composite Index) để tối ưu hóa truy vấn MySQL, và tích hợp Soft Delete thông qua `gorm.DeletedAt`.

### Trụ cột 2: Giao kèo API (API Contract & Context Auth)
- Định nghĩa rõ các Endpoints sẽ cung cấp cho Frontend (Refine.js).
- BẮT BUỘC có đủ:
  + `Method` (GET, POST, PUT, DELETE) & `Route` (VD: `/api/v1/loyalty/scenarios/:id`).
  + `Request Binding Struct` (Sử dụng thẻ binding của Gin như `binding:"required"` để validate payload).
  + `Response Payload Struct` (Thành công & Thất bại, tuyệt đối không lộ trường nhạy cảm như `password` thông qua tag `json:"-"` hoặc Struct chuyên biệt).
- Ghi chú rõ API nào cần đi qua Gin Middleware để xác thực Access Token và parse `userId`/`role` trực tiếp từ JWT.

### Trụ cột 3: Xử lý Cache & Quản lý Trạng thái Token (Redis Integration)
- Đánh giá luồng nghiệp vụ: Quy hoạch cấu trúc lưu trữ và Key Pattern trên Redis đối với các tác vụ Caching dữ liệu ít biến động.
- Cấu hình bảo mật nâng cao: Thiết kế luồng Blacklist Token trên Redis bằng hàm `SetEx` có tính toán TTL còn lại của JWT, và cơ chế xử lý nguyên tử cho luồng Refresh Token Rotation tại endpoint `/refresh-token` cùng luồng giải phóng rác tại endpoint `/logout`.

## BÁO CÁO KẾT QUẢ
Sau khi tạo file xong, in ra một thông báo ngắn gọn: 
*"✅ Đã hoàn tất bản quy hoạch Back-end tại file [Tên file]. Sẵn sàng để thực thi mã nguồn API."*