---
name: integrate-api-webview
description: Biến UI tĩnh (Mock Data) thành UI động bằng cách gọi API thực tế. Bắt buộc xử lý đủ 3 trạng thái (Loading, Success, Error).
triggers:
  - "/integrate-webview"
  - "nối api webview"
---

# NHIỆM VỤ: INTEGRATION ENGINEER (KỸ SƯ TÍCH HỢP)

Khi nhận lệnh `/integrate-webview [Tên file UI] với [Tên API/File Backend Plan]`, bạn đóng vai trò là người nối "đường ống dữ liệu" giữa giao diện tĩnh và máy chủ.

BẮT BUỘC thực hiện tuần tự các bước sau:

## 1. NẠP NGỮ CẢNH CÔNG NGHỆ (DATA FETCHING RULES)
- Đọc `AGENTS.md` để biết dự án đang dùng thư viện Data Fetching nào (VD: `fetch` thuần của Next.js, `Axios`, `React Query` hay `RTK Query`). BẠN BỊ CẤM dùng sai thư viện.
- Rà soát file UI hiện tại để xem dữ liệu Mock Data đang nằm ở đâu.
- Đọc file Backend Plan (hoặc API Contract) được cung cấp để biết chính xác Endpoint, Method, và Payload cấu trúc ra sao.
- Trong NextJS không được chuyển page về client (`'use client'`), bắt buộc phải xử lý data ở server.

## 2. QUY TẮC "KIỀNG 3 CHÂN" KHI GỌI API (BẮT BUỘC)
Khi xóa Mock Data và thay bằng hàm gọi API, BẮT BUỘC phải code đủ 3 trạng thái của giao diện:
1. **[Trạng thái Loading]:** Phải hiển thị Skeleton loading hoặc Spinner khi đang chờ data. Không được để màn hình trắng.
2. **[Trạng thái Error]:** Bọc API trong `try-catch`. Nếu API lỗi, phải có UI báo lỗi rõ ràng (VD: Toast notification hoặc thông báo "Không tải được dữ liệu").
3. **[Trạng thái Empty/Success]:** Nếu data rỗng, hiển thị component "Chưa có dữ liệu". Nếu có data, map (đổ) data vào UI chuẩn xác.

## 3. THỰC THI GHI FILE
- Tiến hành cập nhật trực tiếp vào file UI (`.tsx` hoặc `.vue`).
- Xóa bỏ hoàn toàn Mock Data cũ.
- Tự động import các custom hooks hoặc thư viện call API cần thiết lên đầu file.
- Không dùng hook client: Tuyệt đối không import useRouter từ next/router hay next/navigation vì chúng yêu cầu 'use client'.
- Sử dụng redirect: Import trực tiếp redirect từ next/navigation. Hàm này có thể gọi ngay trong Server Component (khi fetch data bị lỗi/hoặc cần redirect theo điều kiện dữ liệu) hoặc trong Server Actions (khi user submit form).

## 4. BÁO CÁO KẾT QUẢ
- In ra thông báo: *"✅ Đã nối thành công API [Tên API] vào màn hình [Tên Màn Hình]."*
- Ghi chú lại cho người dùng biết bạn đã xử lý Loading và Error như thế nào.