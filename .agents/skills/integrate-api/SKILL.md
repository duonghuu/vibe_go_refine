---
name: integrate-api
description: Biến UI tĩnh (Mock Data) thành UI động bằng cách tích hợp API thực tế vào Refine.js Resource Providers. Bắt buộc xử lý đủ 3 trạng thái (Loading, Success, Error).
triggers:
  - "/integrate"
  - "nối api"
---

# NHIỆM VỤ: INTEGRATION ENGINEER (KỸ SƯ TÍCH HỢP)

Khi nhận lệnh `/integrate [Tên file UI] với [Tên API/File Backend Plan]`, bạn đóng vai trò là người nối "đường ống dữ liệu" giữa giao diện quản trị tĩnh và hệ thống Backend Golang (Gin - GORM).

BẮT BUỘC thực hiện tuần tự các bước sau:

## 1. NẠP NGỮ CẢNH CÔNG NGHỆ (DATA FETCHING RULES)
- Rà soát file UI hiện tại để xem cấu trúc Data Types đang được định nghĩa như thế nào.
- Đọc file Backend Plan (hoặc API Contract) để đối chiếu kiểu dữ liệu (Struct mapping) giữa Go Backend và TypeScript Frontend.
- **Luật thép Framework:** Sử dụng cơ chế `dataProvider` tích hợp sẵn của Refine.js kết hợp với các hooks dữ liệu chuẩn (`useTable`, `useForm`, `useShow`, `useSelect`) để tự động hóa việc fetch dữ liệu, đồng bộ hóa mutation và quản lý cache. BẠN BỊ CẤM gọi trực tiếp các hàm `fetch` hay `axios` phân tán bên ngoài các hooks này khi tương tác với các resource chuẩn của hệ thống.

## 2. QUY TẮC "KIỀNG 3 CHÂN" KHI GỌI API (BẮT BUỘC)
Khi loại bỏ Mock Data và gắn kết dữ liệu thực qua Refine.js và Material UI, BẮT BUỘC phải đảm bảo giao diện xử lý mượt mà đủ 3 trạng thái thông qua các thuộc tính trả về từ hooks:
1. **[Trạng thái Loading]:** Tận dụng thuộc tính `tableQueryResult.isLoading` hoặc `queryResult.isLoading`. Trong lúc đợi, bắt buộc hiển thị Material UI `<Skeleton />` hoặc thành phần loading bọc ngoài của Refine để tránh giật lag màn hình.
2. **[Trạng thái Error]:** Hệ thống phải tự động kích hoạt `useNotification` của Refine để bắn ra `Snackbar (Toast)` cảnh báo lỗi từ Backend (dựa theo mã HTTP Status Code trả về như 400, 401, 500) và điều hướng an toàn nếu token hết hạn.
3. **[Trạng thái Empty/Success]:** Nếu dữ liệu mảng trả về trống, sử dụng component trống chuẩn của hệ thống để thông báo trực quan. Nếu có dữ liệu, đổ trực tiếp vào thuộc tính `rows` của MUI `<DataGrid>` thông qua data mapped từ Refine.

## 3. THỰC THI GHI FILE
- Tiến hành cập nhật trực tiếp vào file UI (`.tsx`).
- Xóa bỏ hoàn toàn Mock Data cũ. Thay thế bằng việc truyền các tham số `resource`, `filters`, `sorters` vào hook dữ liệu của Refine.
- Định nghĩa chặt chẽ Generic Types cho các hook (Ví dụ: `useTable<IProduct, IHttpError>`) để đảm bảo tính an toàn dữ liệu cuối (Type Safety), nghiêm cấm lạm dụng kiểu `any`.
- Kiểm tra tích hợp cơ chế `debounce: 300ms` cho toàn bộ các tác vụ lọc dữ liệu và ô tìm kiếm tự động trên grid.

## 4. BÁO CÁO KẾT QUẢ
- In ra thông báo ngắn gọn: *"✅ Đã nối thành công API [Tên API] vào màn hình quản trị [Tên Màn Hình]."*
- Liệt kê các Hook Refine đã sử dụng kèm theo giải trình ngắn về cách cấu hình xử lý trạng thái dữ liệu (Loading/Error/Empty) để kiểm tra.