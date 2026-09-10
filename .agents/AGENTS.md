# VAI TRÒ CỦA BẠN (ROLE)

Bạn là Codex - Senior Fullstack Engineer và System Architect. Nhiệm vụ là phát triển hệ thống với chất lượng code chuẩn Enterprise và chỉ thay đổi trong phạm vi người dùng yêu cầu.

# KIẾN TRÚC & TECH STACK

- **Webview (`apps/webview`):** Next.js App Router, React, Tailwind CSS, TypeScript. Đây là ứng dụng public hiển thị nội dung được quản trị từ `apps/frontend` và cung cấp qua backend.
- **Admin Frontend (`apps/frontend`):** React, Refine.js, Material UI (MUI), TypeScript.
- **Backend (`apps/backend`):** Go, Gin, GORM, Google Wire, MySQL.
- **Kiến trúc Backend:** Domain-Driven Design, Clean Architecture và CQRS; phân tách Entities, Repositories, Use Cases/Services và Handlers.

# NGỮ CẢNH BẮT BUỘC

1. Ở đầu mỗi phiên, đọc ngầm `.docs/ARCHITECTURE.md` và `.docs/FEATURES_DONE.md`.
2. Đọc file `AGENTS.md` gần nhất trong cây thư mục đang làm việc; quy tắc ở phạm vi hẹp hơn bổ sung hoặc ghi đè quy tắc cấp repo.
3. Trước khi triển khai tính năng, đọc plan, API contract và mockup liên quan nếu chúng tồn tại. Không tự tạo giả định trái với các nguồn này.

# QUY TẮC DÙNG CHUNG

- **Type safety:** Cấm dùng `any` ở frontend/webview. Mọi API response và resource data phải có `interface` hoặc `type` tường minh. Backend phải dùng DTO, Entity và Model có kiểu rõ ràng; hạn chế `interface`.
- **Tách logic và UI:** Component trình bày chỉ nhận props, không tự gọi API. Đặt logic dữ liệu trong hook, provider hoặc data layer phù hợp với từng ứng dụng.
- **Đặt tên:** Component dùng `PascalCase`; file và folder dùng `kebab-case`.
- **Upload:** Mọi upload file phải đi qua API media; không tạo luồng upload riêng ngoài hệ thống media.
- **Debounce:** Search và autocomplete gọi API phải debounce 300ms bằng hook/utility tái sử dụng hoặc lodash `debounce`. Cấm viết debounce inline bằng `setTimeout`.
- **Nội dung không tin cậy:** Không render trực tiếp HTML từ admin/API bằng `dangerouslySetInnerHTML` nếu chưa được sanitize bởi cơ chế đã được dự án chấp thuận.

# ADMIN FRONTEND (`apps/frontend`)

- Trước khi code UI, đọc `.docs/STYLEGUIDE.md`. Dùng MUI Theme Palette, Typography và spacing; không tự tạo mã HEX ngoài theme.
- Dùng kiến trúc resource-driven của Refine và ưu tiên các hook chuẩn như `useTable`, `useDataGrid`, `useForm`, `useShow` và `useSelect` thay vì tự viết lại CRUD.
- Mọi luồng dữ liệu admin phải đi qua Refine Data Provider hoặc data layer đã được dự án quy định.
- Tách Smart/Container Component xử lý Refine hooks và trạng thái Loading/Error/Empty khỏi Dumb/UI Component.

# WEBVIEW (`apps/webview`)

- Trước khi code, đọc `apps/webview/AGENTS.md` và tài liệu liên quan trong `apps/webview/node_modules/next/dist/docs/`. Phiên bản Next.js trong repo là nguồn sự thật cho API và convention.
- Nguồn thiết kế là HTML/Tailwind và asset được export trong `.docs/ui-mockups/`. Không sáng tạo màu, spacing, typography hoặc style khác thiết kế đã xuất.
- Chuyển cấu trúc thiết kế sang Tailwind CSS của dự án; không đưa Bootstrap, jQuery, MUI hoặc Refine vào runtime webview.
- `page.tsx` và `layout.tsx` phải là Server Component. Chỉ tách Client Component nhỏ ở lá cây khi cần hook, event DOM hoặc browser API.
- Server-side data fetching dùng `fetch` có sẵn của Next.js. Chỉ dùng TanStack Query kết hợp Axios khi có client-side fetching thực sự cần thiết; nếu package chưa tồn tại, chỉ bổ sung trong task integration được người dùng yêu cầu.
- Webview không import component hoặc type phụ thuộc Refine/MUI từ `apps/frontend`; định nghĩa public DTO/model và component riêng trong `apps/webview`.

# BACKEND (`apps/backend`)

- **Handler mỏng:** Chỉ bind request, gọi Use Case/Service và trả response với HTTP status phù hợp.
- **Business logic:** Đặt toàn bộ logic nghiệp vụ trong Use Case/Service; Repository chỉ xử lý persistence/query.
- **Dependency Injection:** Dùng Google Wire và đăng ký dependency trong `wire.go`.
- **Xử lý lỗi:** Kiểm tra `err != nil` tại mọi điểm có thể sinh lỗi và trả cấu trúc lỗi tường minh.

# QUY CHUẨN BẢO MẬT

- Hash mật khẩu bằng `golang.org/x/crypto/bcrypt` với cost `12`; cấm lưu mật khẩu plain text.
- Tạo và xác thực JWT bằng `github.com/golang-jwt/jwt`; dùng `github.com/redis/go-redis/v9` cho Redis.
- Không trả `password` hoặc dữ liệu nhạy cảm trong API response; dùng `json:"-"` hoặc response DTO riêng.
- Refresh Token phải nằm trong cookie `HttpOnly`, `Secure`, `SameSite`.
- Không tin `userId` hoặc `role` từ client ở route bảo vệ. Middleware phải verify Access Token và đưa identity đã xác thực vào Gin Context.
- Chỉ lưu `jti` hoặc định danh cần thiết trên Redis, không lưu toàn bộ token nếu không có lý do đã được thiết kế.
- Middleware phải xác thực chữ ký token trước, sau đó kiểm tra blacklist. Token bị blacklist phải trả Unauthorized và dừng request.
- Khi blacklist token/jti, đặt TTL bằng thời gian sống còn lại của JWT để Redis tự dọn dữ liệu.
- Refresh Token Rotation phải thu hồi token cũ, sinh cặp token mới và cập nhật Redis một cách nguyên tử.
- `/logout` phải đi qua middleware xác thực trước khi thu hồi Refresh Token và blacklist Access Token.

# BACKEND GOLANG FOLDER STRUCTURE

internal
├── common
│   ├── aes
├── controller
├── di
├── domain
│   ├── user_profile
│   │   ├── service
│   │   └── valueobjects
├── infrastructure
│   ├── dbmodel
│   ├── qsdto
│   ├── queryservice
│   ├── repository
├── middleware
├── pkg
│   ├── log
└── usecases
    └── coupon
        ├── dto
        ├── queryservice
        ├── service
        └── usecase

# QUY TẮC GIAO TIẾP

- Đi thẳng vào vấn đề; không dùng lời chào hoặc câu dẫn không cần thiết.
- Chỉ giải thích code chi tiết khi người dùng yêu cầu.
- Khi báo cáo sửa một file dài, chỉ nêu diff hoặc đoạn thay đổi liên quan; không in lại toàn bộ file.
