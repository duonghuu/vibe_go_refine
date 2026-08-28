---
name: code-ui-webview
description: Chuyển đổi bản vẽ thành Code React/Vue, tuân thủ chặt chẽ Frontend Plan, ép buộc dùng Framework chuẩn và ưu tiên tái sử dụng Component.
triggers:
  - "/code-ui-webview"
  - "code giao diện webview"
---

# NHIỆM VỤ: LẬP TRÌNH GIAO DIỆN (UI EXECUTION)

Khi nhận lệnh `/code-ui-webview [Tên dự án] [Tên Màn Hình / Component]`, bạn đang đóng vai trò là một Frontend Developer thi công. Hãy BẮT BUỘC thực hiện tuần tự 4 bước sau một cách im lặng, chỉ báo cáo kết quả cuối cùng:

## BƯỚC 1: NẠP NGỮ CẢNH VÀ QUY HOẠCH
Trước khi làm bất cứ điều gì, bạn BẮT BUỘC phải đọc ngầm 2 file:
1. `.docs/STYLEGUIDE.md` (Để lấy biến màu Tailwind, font chữ, quy tắc bo góc).
2. Tương ứng file kế hoạch trong `.docs/frontend-plans/` (Để biết cấu trúc Smart/Dumb Component và Interface Props).

## BƯỚC 2: QUÉT THƯ VIỆN & TÁI SỬ DỤNG (QUAN TRỌNG TỐI THƯỢNG)
- Quét thư mục `apps/frontend/components/ui/` (Thư viện UI dùng chung).
- Nếu bản vẽ có chứa Nút bấm (Button), Input, Thẻ (Card)... BẠN PHẢI TÌM XEM component đó đã tồn tại chưa.
- **Luật thép:** Nếu ĐÃ CÓ, tuyệt đối không code lại, BẮT BUỘC phải import component đó vào để dùng. Chỉ tạo file mới cho những cấu trúc đặc thù chưa từng xuất hiện.

## BƯỚC 3: ĐỌC BẢN VẼ
- Kết nối với bản vẽ mà người dùng vừa cung cấp.
- Ánh xạ (Map) các thuộc tính đồ họa sang các class của **Tailwind CSS**. Không dùng CSS thuần hay mã HEX lạ.

## BƯỚC 4: SINH CODE VÀ ÉP KHUÔN FRAMEWORK
Tiến hành gõ code vào thư mục dự án theo các quy tắc của bản Kế hoạch, ĐỒNG THỜI tuân thủ tuyệt đối Đạo luật Framework sau:

[RÀNG BUỘC FRAMEWORK: NEXT.JS APP ROUTER]
1. Kiến trúc Client/Server: Mặc định là Server Component. CHỈ thêm `'use client'` nếu Component đó thực sự cần dùng Hook (`useState`, `useEffect`) hoặc xử lý sự kiện DOM.
2. Cú pháp: Sử dụng Arrow Function. Khai báo Props bằng `interface` (TypeScript).
3. Tối ưu: BẮT BUỘC dùng `<Image />` từ `next/image` và `<Link />` từ `next/link`.
4. State: Ưu tiên Zustand cho Global State và Server Actions cho mutate dữ liệu.

## DATA FETCHING
- Bạn bị cấm gọi API trong bước này
- Thay vào đó, bạn phải dùng dữ liệu giả (mock data)
- Dữ liệu giả phải được định nghĩa bằng interface và type
- Hình ảnh dữ liệu giả phải được lấy từ unsplash.com
- Dữ liệu giả phải được định nghĩa ở file mock data (.json) trong thư mục `.docs/mock-data/`

## BÁO CÁO KẾT QUẢ
Sau khi code xong, in ra danh sách các file `.tsx` vừa tạo hoặc chỉnh sửa. Hỏi người dùng xem có cần điều chỉnh padding/margin nào không trước khi chạy lệnh `/save`.