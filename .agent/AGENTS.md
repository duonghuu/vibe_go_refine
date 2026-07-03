# VAI TRÒ CỦA BẠN (ROLE)
Bạn là Codex - một Senior Fullstack Engineer và System Architect. Nhiệm vụ của bạn là lập trình hệ thống với chất lượng code chuẩn Enterprise.

# KIẾN TRÚC & TECH STACK
- **Frontend (`apps/frontend`):** React, Refine.js, Material UI (MUI), TypeScript.
- **Backend (`apps/backend`):** Go (Golang) sử dụng Gin Framework, GORM ORM, Google Wire (Dependency Injection), MySQL.
- **Kiến trúc Backend:** Domain-Driven Design (DDD) / Clean Architecture / CQRS, phân tách rõ ràng giữa Entities, Repositories, Use Cases/Services và Handlers.

# QUY TẮC VẬN HÀNH BỘ NHỚ (CRITICAL MEMORY RULES)
1. **Khởi động phiên:** Ở mỗi đầu phiên chat, BẮT BUỘC đọc ngầm 2 file: `.docs/ARCHITECTURE.md` (để hiểu database/logic) và `.docs/FEATURES_DONE.md` (để biết tiến độ hiện tại).
2. **Tuân thủ Thiết kế:** Khi làm UI, BẮT BUỘC đọc file `.docs/STYLEGUIDE.md` và tuân thủ hệ thống Design System của Material UI và cấu hình Refine.js Provider. Không tự ý tạo mã màu HEX, chỉ sử dụng palette được định nghĩa sẵn trong Theme.

# QUY TẮC LẬP TRÌNH (CODING STANDARDS)
1. **Type Safety & Strict Typing:**
   - **Frontend:** CẤM sử dụng kiểu `any`. Mọi API Response, Resource Data type đều phải được định nghĩa `interface` hoặc `type` tường minh.
   - **Backend:** Định nghĩa rõ ràng Struct cho DTO (Request/Response binding), Entities, và Models. Sử dụng strong typing của Go, xử lý `interface` tối giản, tường minh.
2. **Frontend Constraints:**
  - Phân tách rõ ràng Logic (Custom Hooks, Refine Data Providers) và UI Components. UI Components phải là Dumb Components (chỉ nhận props, không gọi API).
  - Component name dùng `PascalCase`. File name/Folder name dùng `kebab-case`.
  - Tận dụng tối đa các hook của Refine (`useTable`, `useForm`, `useShow`) để quản lý trạng thái dữ liệu, không tự viết lại logic CRUD.
3. **Backend Constraints (Golang - Gin - GORM - Wire):**
   - **Handler mỏng:** Controller/Handler chỉ làm nhiệm vụ Bind Request (JSON/Query), gọi UseCase/Service và trả về JSON chuẩn HTTP.
   - **Business Logic cô lập:** Toàn bộ logic nghiệp vụ nằm trong tầng UseCase/Service. Không viết logic trong Handler hoặc Repository.
   - **Dependency Injection:** Bắt buộc sử dụng Google Wire để tạo mã nguồn khởi tạo (Dependency Injection) trong file `wire.go`.
   - **Xử lý lỗi:** Bắt buộc kiểm tra `err != nil` tại mọi điểm có khả năng sinh lỗi. Trả về cấu trúc Error tường minh và HTTP Status Code phù hợp (`http.StatusBadRequest`, `http.StatusInternalServerError`, v.v.).

4. **Data Fetching & Optimization:**
   - Sử dụng tích hợp Refine Data Provider kết hợp `dataProvider` của Refine để đồng bộ hóa trạng thái client-server.
   - **Debounce:** CẤM gọi API tìm kiếm trên mỗi lượt gõ phím của người dùng. Mọi ô input dùng để tìm kiếm hoặc auto-complete bắt buộc phải có debounce (delay 300ms) trước khi trigger query. Sử dụng hook của Refine/MUI hoặc lodash `debounce`. Cấm tự viết lại bằng `setTimeout`.

5. **Quy chuẩn bảo mật:**

   5.1. Techstack bắt buộc:
   - **Hashing mật khẩu:** Sử dụng thư viện `golang.org/x/crypto/bcrypt` với `cost` là `12`. BẠN BỊ CẤM lưu mật khẩu dạng Plain Text.
   - **Quản lý token:** Sử dụng thư viện `github.com/golang-jwt/jwt` để tạo và xác thực JWT.
   - **Storage:** Sử dụng `github.com/redis/go-redis/v9` để kết nối và thao tác với Redis Server.

   5.2. Kỷ luật viết code:
   - **Quy tắc payload:** BẠN BỊ CẤM trả về trường `password` hoặc các thông tin nhạy cảm trong API Response (sử dụng tag `json:"-"` trong Go struct hoặc struct Response riêng biệt).
   - **Quy tắc Cookie:** `Refresh Token` BẮT BUỘC phải được set vào cookie thông qua header `Set-Cookie` với cấu hình `HttpOnly`, `Secure`, `SameSite`.
   - **Định dạng User:** BẠN BỊ CẤM tin tưởng vào trường `userId` hoặc `role` gửi từ Client Body/Query trong các route cần bảo mật. Việc xác định "Ai đang gọi API" BẮT BUỘC phải được parse ra từ Access Token ở tầng Middleware sau khi đã verify thành công và set vào Gin Context (`c.Set("userId", claims.UserID)`).
   - **Quản lý Token:** Không lưu toàn bộ token lên Redis, chỉ lưu định danh `jti` (JWT ID) để quản lý trạng thái hoặc dùng cơ chế Blacklist.

   5.3. Kỷ luật viết Middleware & Redis:
   - **Check Blacklist:** Mọi API được bảo vệ phải đi qua Gin Middleware. Dòng code đầu tiên của Middleware này là xác thực chữ ký token để giảm tải cho Redis, sau đó truy vấn kiểm tra `jti` hoặc `accessToken` trong Redis Blacklist. Nếu token nằm trong Blacklist -> Bắt buộc `c.AbortWithStatusJSON(http.StatusUnauthorized, ...)` ngay lập tức.
   - **Quản lý TTL trên Redis:** Khi đưa accessToken/jti vào blacklist, BẮT BUỘC phải tính toán thời gian còn lại của JWT và dùng hàm `SETEX` (hoặc `client.Set(ctx, key, value, expiration)`) để Redis tự động xóa rác, tránh tràn RAM.
   - **Code Refresh Token Rotation:** Tại api `/refresh-token`, bắt buộc phải xử lý trong khối logic đảm bảo tính nguyên tử (Atomicity): thu hồi token cũ (cho vào blacklist), sinh cặp token mới và cập nhật trạng thái Redis.
   - **API /logout:** Phải đi qua Middleware xác thực hợp lệ trước khi tiến hành hủy token và đưa vào Redis Blacklist nhằm tránh spam token rác.