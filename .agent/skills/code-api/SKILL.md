---
name: code-api
description: Viết các API Endpoints dựa trên BACKEND_PLAN, kết nối trực tiếp với Database và TRỰC TIẾP TẠO FILE vật lý.
triggers:
  - "/code-api"
  - "viết api"
---

# NHIỆM VỤ: BACKEND DEVELOPER (API INTEGRATOR)

Khi nhận lệnh `/code-api [Đường dẫn file Backend Plan]`, bạn đóng vai trò là một Backend Developer thực thi. Nhiệm vụ của bạn là hiện thực hóa "Trụ cột 2: Giao kèo API" thành code thực tế có thể chạy được.

BẮT BUỘC thực hiện tuần tự các bước sau một cách im lặng:

## 1. NẠP NGỮ CẢNH "KIỀNG 3 CHÂN"
Để viết API không bị lỗi, BẠN BẮT BUỘC phải đọc 3 nguồn dữ liệu sau:
1. **`AGENTS.md`** hoặc cấu trúc thư mục hiện tại: Xác định kiến trúc phân tầng (Clean Architecture / Domain-Driven Design) để đặt mã nguồn đúng lớp (Handlers, UseCases/Services, Repositories).
2. **File Backend Plan được truyền vào**: Để biết cần viết những API nào (GET, POST, PUT, DELETE), Payload ra sao.
3. **Mô hình Database thực tế** (các file Go struct định nghĩa GORM Models tại tầng Domain/Entity): Để đảm bảo code truy xuất đúng tên trường, đúng kiểu dữ liệu và đúng quan hệ bảng đã được `/code-db` tạo ra trước đó.

## 2. KỶ LUẬT VIẾT API (DEFENSIVE PROGRAMMING)
Khi sinh code Golang, BẮT BUỘC áp dụng các lớp phòng thủ nghiêm ngặt:
- **Xác thực đầu vào (Validation):** Tuyệt đối không tin tưởng dữ liệu từ Client. Bắt buộc sử dụng các thẻ binding của Gin (ví dụ: `binding:"required"`) để validate struct Request Payload ngay tại tầng Handler.
- **Xử lý lỗi (Error Handling):** Bắt buộc kiểm tra `err != nil` ở mọi bước gọi hàm (từ Bind JSON, gọi UseCase, đến tầng DB). Tuyệt đối không nuốt lỗi (swallow error).
- **Mã trạng thái (Status Code):** Trả về HTTP Status Code chuẩn mực (`http.StatusOK`, `http.StatusCreated`, `http.StatusBadRequest`, `http.StatusUnauthorized`, `http.StatusInternalServerError`) dưới dạng JSON object rõ ràng.
- **Bảo mật**: Tuyệt đối không trả về trường `password` hoặc thông tin nhạy cảm (Sử dụng struct Response riêng biệt hoặc tag `json:"-"`). Parse thông tin định danh (`userId`, `role`) từ Gin Context đã được xử lý qua Auth Middleware, cấm tin tưởng dữ liệu body gửi từ client.

## 3. THỰC THI GHI FILE (FILE EXECUTION)
- Phân tách code rõ ràng theo mô hình kiến trúc:
  - **Tầng Handler:** Chỉ bind request, gọi UseCase, trả về HTTP JSON.
  - **Tầng UseCase/Service:** Xử lý toàn bộ logic nghiệp vụ (Business Logic).
  - **Tầng Repository:** Thực hiện các câu lệnh truy vấn GORM (`db.WithContext(ctx).Find(...)`).
- TRỰC TIẾP TẠO MỚI HOẶC GHI ĐÈ file code vật lý (`.go`) vào các thư mục tương ứng trong dự án.
- Bắt buộc cập nhật file cấu hình Dependency Injection `wire.go` và chạy lệnh sinh code tự động để đăng ký các Provider mới (Handler, UseCase, Repository) vừa tạo.

## 4. XUẤT BÁO CÁO VÀ TÀI LIỆU API (API SPECS)
- Sau khi ghi file code thành công, BẮT BUỘC tự động tạo hoặc cập nhật file `.docs/api-endpoints.yaml` để lưu trữ tài liệu danh sách các API vừa được tạo.
- Định dạng file YAML phải tuân thủ cấu trúc OpenAPI/Swagger rõ ràng, bao gồm: Path, Method, Request Body, Response (Success/Error), và Đường dẫn tới các file code Go vật lý tương ứng.
- **Tại giao diện chat:** Chỉ cần in ra một thông báo ngắn gọn kèm danh sách các file vật lý đã tạo/sửa đổi và đường dẫn tới file YAML để người dùng theo dõi.