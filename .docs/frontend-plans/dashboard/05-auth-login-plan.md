# FRONTEND PLAN: Xác thực và Đăng nhập (Auth & Login)

## 1. Mục tiêu
- Xây dựng màn hình Đăng nhập (Login) và cơ chế xác thực toàn cục cho hệ thống TechBite.
- Cung cấp trải nghiệm đăng nhập mượt mà, an toàn với cơ chế JWT (Access Token & Refresh Token).
- Đảm bảo người dùng luôn được duy trì phiên đăng nhập tự động (Refresh Token Rotation) mà không bị gián đoạn.
- Phân tách rõ ràng luồng UI Component (MUI) và Data/Auth Hook (Refine.js).

## 2. Phạm vi chức năng
### 2.1. Giao diện Đăng nhập
- Form đăng nhập cơ bản với Email và Password.
- Hỗ trợ tuỳ chọn "Remember Me" (nếu cần).
- Xử lý các trạng thái Validation (Email không đúng định dạng, Password trống, v.v.).
- Xử lý trạng thái Loading (disable nút submit và hiển thị spinner).

### 2.2. Xử lý Trạng thái & API Auth
- Tích hợp RefineJS `AuthProvider`: `login`, `logout`, `check`, `getIdentity`, `onError`.
- Triển khai cơ chế Interceptor (Axios) để tự động gắn Access Token vào Header `Authorization`.
- Triển khai luồng Refresh Token Rotation tự động khi Access Token hết hạn (bắt lỗi 401).

### 2.3. Điều hướng & Phân quyền
- Bảo vệ các routes nội bộ (Dashboard) thông qua cơ chế `Authenticated` của Refine.
- Tự động chuyển hướng về trang `/login` nếu phiên bản hết hạn hoàn toàn (Refresh Token thất bại).
- Quản lý Identity (thông tin User) để hiển thị trên Header (Avatar, Name, Role).

## 3. Luồng màn hình (Screen Flow)
### 3.1. Trang Login (`/login`)
- **UI Layout**: Giao diện tập trung, thẻ (Card) đăng nhập nằm giữa màn hình. Background đơn giản hoặc theo nhận diện thương hiệu TechBite.
- **Thành phần**:
  - Logo hệ thống.
  - Ô nhập liệu: Email, Password.
  - Checkbox: Ghi nhớ đăng nhập (Remember me).
  - Nút: Đăng nhập (kèm hiệu ứng loading).
  - Thông báo lỗi (Alert/Snackbar) khi đăng nhập thất bại.

### 3.2. Quá trình xử lý
1. Người dùng nhập Email + Password và nhấn Đăng nhập.
2. Refine kích hoạt hook `useLogin`.
3. Trạng thái nút chuyển sang `Loading`.
4. Gọi API `/auth/login`.
5. Thành công:
   - Lưu Token (Access Token lưu bộ nhớ/localStorage, Refresh Token lưu HTTPOnly Cookie).
   - Fetch Identity (`/auth/me`).
   - Chuyển hướng vào trang chính (`/`).
6. Thất bại:
   - Hiển thị Toast Error rõ ràng (Sai tài khoản, Tài khoản bị khoá...).
   - Reset Loading state.

## 4. Component Breakdown
### 4.1. Pages
- `LoginPage` (`apps/frontend/src/pages/auth/login.tsx`): Trang đăng nhập chính.

### 4.2. UI Components (Dumb Components)
- `AuthLayout`: Bố cục chung cho các trang xác thực.
- `LoginForm`: Form chứa các input fields (có sử dụng MUI `TextField`).
- `SubmitButton`: Nút đăng nhập với trạng thái loading tích hợp.

### 4.3. Providers & Hooks (Smart Logic)
- `authProvider` (`apps/frontend/src/providers/authProvider.ts`): Cốt lõi xử lý logic JWT.
- `axiosInstance` (`apps/frontend/src/providers/axiosInstance.ts`): Cấu hình Interceptors xử lý Header và Refresh Token.

## 5. Data & State
### 5.1. Dữ liệu API (Mock & Thực tế)
- **Request Login**: `{ "email": "admin@techbite.com", "password": "..." }`
- **Response Login**: `{ "accessToken": "...", "expiresIn": 1800, "user": { ... } }` (Refresh token trả về qua Set-Cookie HTTPOnly).

### 5.2. State Frontend
- Trạng thái các trường input (Email, Password).
- Validation state (Error messages).
- Trạng thái `isLoading` khi submit.
- Cấu trúc Identity lưu trữ: `User ID`, `Email`, `Name`, `Role`.

## 6. Quy ước hiển thị & Design Specs
- **Màu sắc**:
  - Nút Submit: Dùng Primary Action `palette.primary.main`.
  - Nút Loading/Disabled: `palette.action.disabled`.
  - Toast lỗi: `bg-red-500` hoặc Error Snackbar của MUI.
- **Typo**: Font chữ lớn rõ ràng.
- **Card**: Sử dụng viền mỏng phẳng, đổ bóng nhẹ (Elevation 1) theo chuẩn STYLEGUIDE.

## 7. Tương tác chính & API Integration
### 7.1. API Endpoints
- `POST /auth/login`: Xác thực credentials.
- `POST /auth/logout`: Đăng xuất, xoá token.
- `POST /auth/refresh-token`: Làm mới token tự động.
- `GET /auth/me`: Lấy thông tin user đăng nhập.

### 7.2. Axios Interceptor Logic
- **Request**: Thêm `Authorization: Bearer <accessToken>`.
- **Response Error (401)**:
  - Nếu API lỗi 401: Tạm dừng các request queue.
  - Gọi API `/auth/refresh-token`.
  - Nếu thành công: Cập nhật Access Token mới, chạy lại request queue.
  - Nếu thất bại: Gọi AuthProvider `logout` để đẩy về màn `/login`.

## 8. Trạng thái giao diện
- `Initial`: Form trống sẵn sàng nhập.
- `Typing`: Đang nhập dữ liệu, tắt bớt validate hiển thị trước đó.
- `Submitting (Loading)`: Nút Đăng nhập disable, có vòng quay loading.
- `Error`: Snackbar hoặc Alert đỏ hiện lên báo lỗi.
- `Success`: Redirect mượt mà vào Dashboard.

## 9. Gợi ý triển khai kỹ thuật
- **Form Management**: Sử dụng `useForm` từ `@refinedev/react-hook-form` hoặc MUI kết hợp `react-hook-form` để quản lý validation tối ưu.
- **Validation Schema**: Sử dụng `zod` hoặc `yup` để chuẩn hoá rule: Email phải hợp lệ, Password không được để trống.
- **Axios**: Triển khai `axios-retry` hoặc custom queue để xử lý đồng bộ nhiều request bị 401 cùng lúc.
- **Bảo mật**: Tuyệt đối không lưu Refresh Token trong LocalStorage (Nên được Set-Cookie `HttpOnly` từ backend, FE không cần quản lý Refresh Token thủ công).
- **Tuân thủ DDD**: `authProvider.ts` hoàn toàn tách biệt khỏi UI Component.

## 10. Định tuyến đề xuất
- `/login`: Route public dành cho đăng nhập.
- Các route khác (ngoại trừ login/register) đều bọc bởi `<Authenticated>`.

## 11. Tiêu chí hoàn thành
- Giao diện Đăng nhập đúng Styleguide, responsive.
- Đăng nhập thành công với API thực tế (hoặc Mock chuẩn).
- Bắt lỗi form đầy đủ (Email sai định dạng, thiếu password).
- Hiển thị Toast lỗi khi sai mật khẩu.
- Xử lý được Interceptor Refresh Token (401 Retry logic).
- Đăng xuất thành công, xoá trắng dữ liệu phiên và Cookie.

## 12. Checklist triển khai
- [ ] Tạo file `axiosInstance.ts` cấu hình baseURL và Interceptors (cho JWT, 401, Retry).
- [ ] Viết logic `authProvider.ts` cài đặt 5 phương thức cốt lõi.
- [ ] Tích hợp `authProvider` vào component `<Refine>` ở `App.tsx`.
- [ ] Tạo UI `LoginPage` với MUI Layout.
- [ ] Code Form (Email, Password) với validation.
- [ ] Kết nối hook `useLogin` vào form submit.
- [ ] Test luồng Đăng nhập (Thành công/Thất bại).
- [ ] Test luồng tự động Refresh Token khi Access Token hết hạn (mock timeout).
- [ ] Test luồng Đăng xuất.
