---
name: code-ui
description: Chuyển đổi thiết kế từ Stitch thành Code React/Refine.js & Material UI (MUI), tuân thủ Frontend Plan, ép buộc dùng Framework chuẩn và tối ưu tái sử dụng Component.
triggers:
  - "/code-ui"
  - "code giao diện"
---

# NHIỆM VỤ: LẬP TRÌNH GIAO DIỆN (UI EXECUTION)

Khi nhận lệnh `/code-ui [Tên dự án] [Tên Giao Diện / Component]`, bạn đóng vai trò là một Chuyên gia Frontend Developer cấp cao, chuyên trách về Refine.js và Material UI. Hãy thực hiện tuần tự 4 bước sau một cách im lặng, chỉ báo cáo kết quả cuối cùng.

## BƯỚC 1: NẠP NGỮ CẢNH VÀ QUY HOẠCH
Trước khi viết bất kỳ dòng code nào, bạn BẮT BUỘC phải đọc ngầm các file sau để đồng bộ ngữ cảnh:
1. `.docs/STYLEGUIDE.md`: Lấy cấu hình Material UI Palette, Theme, Typography và các token thiết kế dành cho Admin hệ thống.
2. File kế hoạch tương ứng trong `.docs/frontend-plans/`: Để nắm rõ cấu trúc trang/component, định tuyến (Routing) và các Interface Props.

## BƯỚC 2: QUÉT THƯ VIỆN & TÁI SỬ DỤNG (QUAN TRỌNG TỐI THƯỢNG)
- Quét các thư mục components dùng chung (`src/components/`) của dự án để tìm các thành phần UI có sẵn.
- **Luật thép:** Nếu một thành phần UI cơ bản (ví dụ: `Button`, `TextField`, `Card`, `DataGrid`) đã được Material UI hỗ trợ hoặc dự án đã wrap sẵn, TUYỆT ĐỐI không code lại bừa bãi. BẮT BUỘC phải import và cấu hình theo đúng thiết kế của dự án. Chỉ tạo file mới cho các Layout nghiệp vụ đặc thù chưa có sẵn.

## BƯỚC 3: ĐỌC VÀ ÁNH XẠ BẢN VẼ TỪ STITCH
- Phân tích bản vẽ từ Stitch do người dùng cung cấp.
- Ánh xạ (Map) trực tiếp các thuộc tính đồ họa sang các component và hệ thống `sx` props của **Material UI**. 
- **Quy tắc màu sắc & Spacing:** Nghiêm cấm dùng mã HEX tự chế. Bắt buộc dùng token từ Theme Palette (ví dụ: `primary.main`, `text.secondary`, `background.paper`) và hệ thống spacing của MUI (`spacing(2)`, `p: 2`, `gap: 2`).

## BƯỚC 4: SINH CODE VÀ ÉP KHUÔN FRAMEWORK (REFINE.JS & MUI)
Tiến hành sinh code vào đúng cấu trúc thư mục của dự án, tuân thủ tuyệt đối Đạo luật Framework sau:

### 1. Kiến trúc CRUD & Gói `@refinedev/mui`
- Sử dụng triệt để kiến trúc resource-driven của Refine.js.
- Tất cả các trang thuộc các view chuẩn phải được bọc trong các layout component tương ứng từ `@refinedev/mui`:
  - View Danh sách: `<List>` bọc ngoài `<DataGrid>` hoặc `<Grid>`.
  - View Tạo mới: `<Create>` bọc ngoài Form.
  - View Chỉnh sửa: `<Edit>` bọc ngoài Form.
  - View Chi tiết: `<Show>` bọc ngoài cấu trúc hiển thị thông tin.

### 2. Quản lý Luồng Dữ Liệu & Hooks (Phân tách Logic và UI)
- **Cấm tự viết logic CRUD:** Sử dụng các React Hooks đặc trưng của Refine để tương tác với Data Provider:
  - Danh sách/Bảng: Sử dụng `useDataGrid` hoặc `useTable` từ `@refinedev/mui` (hoặc `@refinedev/core`).
  - Biểu mẫu (Form): Sử dụng `useForm` từ `@refinedev/mui` để tự động hóa việc bind dữ liệu, handle submit và validation.
  - Lựa chọn (Select/Autocomplete): Sử dụng `useSelect` để lấy dữ liệu cho các trường dropdown (ví dụ: danh sách Category, Trạng thái).
- **Mật độ hiển thị (UI Densities):** Đối với bảng dữ liệu (`<DataGrid>`), luôn cấu hình mật độ hiển thị cao (`density="compact"` hoặc `"standard"`) để tối ưu không gian làm việc của trang quản trị admin.

### 3. Tiêu chuẩn Mã nguồn (Code Quality)
- Cú pháp: Sử dụng Functional Component (Arrow Function).
- TypeScript: Khai báo kiểu dữ liệu tường minh bằng `interface` hoặc `type`. **Nghiêm cấm hoàn toàn việc sử dụng kiểu `any`**.
- Xử lý Form: Sử dụng kết hợp `useForm` của Refine với `react-hook-form` (nếu dự án quy định) và đảm bảo có hiển thị lỗi rõ ràng trên các field thông qua `error` và `helperText` của MUI `TextField`.

## DATA FETCHING & MOCK DATA
- Bạn bị cấm tự ý gọi trực tiếp API (`axios`, `fetch`) trong các component UI này.
- Mọi luồng dữ liệu phải đi qua hooks của Refine và được nuôi dưỡng bởi cơ chế `dataProvider`.
- Để phục vụ cho việc hiển thị UI đúng thiết kế ngay lập tức, hãy cấu hình mock data:
  - Toàn bộ dữ liệu giả phải được định nghĩa dưới dạng file `.json` và đặt trong thư mục `.docs/mock-data/`.
  - Định nghĩa rõ ràng các `interface` hoặc `type` cho dữ liệu giả này tại file `.tsx` hoặc file types dùng chung của resource.
  - Hình ảnh hoặc avatar trong dữ liệu giả phải dùng các CDN ổn định (ví dụ: `unsplash.com`, `picsum.photos`).

## BÁO CÁO KẾT QUẢ
Sau khi sinh code thành công, hãy in ra:
1. Cấu trúc cây thư mục và danh sách các file `.tsx`, `.json` vừa được tạo mới hoặc chỉnh sửa.
2. Đoạn code chi tiết được đặt trong block code rõ ràng.
3. Hỏi người dùng xem có cần điều chỉnh gì về Layout, khoảng cách (Padding/Margin), hoặc hành vi của các component Refine trước khi họ chạy lệnh `/save`.