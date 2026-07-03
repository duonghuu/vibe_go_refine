---
name: code-ui
description: Chuyển đổi bản vẽ từ Stitch thành Code React/Refine.js, tuân thủ chặt chẽ Frontend Plan, ép buộc dùng Framework chuẩn và ưu tiên tái sử dụng Component.
triggers:
  - "/code-ui"
  - "code giao diện"
---

# NHIỆM VỤ: LẬP TRÌNH GIAO DIỆN (UI EXECUTION)

Khi nhận lệnh `/code-ui [Tên dự án] [Tên Màn Hình / Component]`, bạn đang đóng vai trò là một Frontend Developer thi công. Hãy BẮT BUỘC thực hiện tuần tự 4 bước sau một cách im lặng, chỉ báo cáo kết quả cuối cùng:

## BƯỚC 1: NẠP NGỮ CẢNH VÀ QUY HOẠCH
Trước khi làm bất cứ điều gì, bạn BẮT BUỘC phải đọc ngầm 2 file:
1. `.docs/STYLEGUIDE.md` (Để lấy cấu hình Material UI Palette, Theme và các token thiết kế cho Admin).
2. Tương ứng file kế hoạch trong `.docs/frontend-plans/` (Để biết cấu trúc trang/component và Interface Props).

## BƯỚC 2: QUÉT THƯ VIỆN & TÁI SỬ DỤNG (QUAN TRỌNG TỐI THƯỢNG)
- Quét các thư mục components dùng chung của dự án.
- Nếu bản vẽ của Stitch có chứa các phần tử UI cơ bản... BẠN PHẢI TÌM XEM component hoặc hệ thống Material UI tương ứng đã được thiết lập chưa.
- **Luật thép:** Nếu ĐÃ CÓ hoặc Material UI đã hỗ trợ sẵn (ví dụ: `DataGrid`, `Card`, `Button`, `TextField`), tuyệt đối không code lại bừa bãi. BẮT BUỘC phải import và cấu hình theo đúng hệ thống design system đã thống nhất. Chỉ tạo file mới cho những layout nghiệp vụ đặc thù chưa từng xuất hiện.

## BƯỚC 3: ĐỌC BẢN VẼ TỪ STITCH
- Kết nối với bản vẽ mà người dùng vừa cung cấp.
- Ánh xạ (Map) các thuộc tính đồ họa sang các component và thuộc tính (sx props) của **Material UI**. Không dùng mã HEX tự chế, chỉ dùng token màu từ Theme Palette.

## BƯỚC 4: SINH CODE VÀ ÉP KHUÔN FRAMEWORK
Tiến hành gõ code vào thư mục dự án theo các quy tắc của bản Kế hoạch, ĐỒNG THỜI tuân thủ tuyệt đối Đạo luật Framework sau:

[RÀNG BUỘC FRAMEWORK: REFINE.JS & MATERIAL UI]
1. Kiến trúc CRUD: Tận dụng triệt để kiến trúc resource-driven của Refine.js. Sử dụng các component bọc chuẩn như `<List>`, `<Create>`, `<Edit>`, `<Show>` từ gói `@refinedev/mui`.
2. Cú pháp: Sử dụng Functional Component (Arrow Function). Khai báo kiểu dữ liệu tường minh bằng `interface` hoặc `type` (TypeScript), nghiêm cấm sử dụng `any`.
3. Quản lý luồng dữ liệu: Phân tách rõ ràng Logic và UI. Sử dụng các React Hooks đặc trưng của Refine (`useTable`, `useForm`, `useShow`, `useSelect`) để gắn kết trạng thái, không tự viết lại logic tương tác CRUD.
4. UI Densities: Đảm bảo các bảng dữ liệu (`<DataGrid>`) sử dụng cấu hình mật độ hiển thị cao (`dense` hoặc `medium`) để tối ưu hóa không gian làm việc của admin.

## DATA FETCHING
- Bạn bị cấm tự gọi trực tiếp API trong bước này.
- Thay vào đó, bạn phải cấu hình mock data thông qua cơ chế `dataProvider` của Refine.js hoặc định nghĩa dữ liệu giả (mock data) có cấu trúc chuẩn.
- Dữ liệu giả phải được định nghĩa bằng interface và type rõ ràng.
- Hình ảnh dữ liệu giả phải được lấy từ các nguồn CDN ổn định hoặc unsplash.com.
- Dữ liệu giả phải được định nghĩa ở file mock data (.json) trong thư mục `.docs/mock-data/`.

## BÁO CÁO KẾT QUẢ
Sau khi code xong, in ra danh sách các file `.tsx` vừa tạo hoặc chỉnh sửa. Hỏi người dùng xem có cần điều chỉnh bố cục UI, khoảng cách padding/margin nào không trước khi chạy lệnh `/save`.