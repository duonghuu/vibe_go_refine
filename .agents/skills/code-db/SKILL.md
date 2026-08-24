---
name: code-db
description: Chuyển đổi Thiết kế Dữ liệu từ BACKEND_PLAN và TRỰC TIẾP TẠO FILE Struct Models (GORM), định nghĩa Migration SQL, cấu hình ánh xạ quan hệ thực thể tuân thủ chuẩn Enterprise.
triggers:
  - "/code-db"
  - "viết schema"
---

# NHIỆM VỤ: DATABASE ADMINISTRATOR (DBA) & THỰC THI FILE

Khi nhận lệnh `/code-db [Đường dẫn file Backend Plan]`, bạn đóng vai trò là một Database Administrator cấp cao. Nhiệm vụ của bạn là đọc "Trụ cột 1: Thiết kế Dữ liệu" trong bản kế hoạch và TRỰC TIẾP GHI RA FILE vật lý trong dự án.

BẮT BUỘC thực hiện tuần tự và tuân thủ các quy tắc sắt đá sau:

## 1. NẠP NGỮ CẢNH & XÁC ĐỊNH VỊ TRÍ FILE

* Xác định cấu trúc dự án dựa trên Clean Architecture / Domain-Driven Design để định vị nơi lưu trữ.
* **[QUAN TRỌNG]:** Tự động xác định đường dẫn lưu file chuẩn của Golang và GORM. Toàn bộ các Struct Models định nghĩa Schema và cấu hình ánh xạ bảng phải được đặt trong thư mục thuộc tầng Domain/Entity (Ví dụ: `internal/domain/entities/` hoặc `pkg/models/`). File Migration SQL (`.sql`) phải được lưu vào đúng thư mục quản lý migration dữ liệu (Ví dụ: `deployments/migrations/` hoặc `database/migrations/`).

## 2. QUY TẮC ĐẶT TÊN (NAMING CONVENTION)

* **Tên Struct Model:** PascalCase, số ít (VD: `User`, `Product`).
* **Tên Trường (Fields trong Struct):** PascalCase cho Struct Field và bắt buộc định nghĩa tag `gorm` kèm tag `json` dạng snake_case/camelCase phù hợp với cấu hình frontend (VD: `CreatedAt time.Time `gorm:"column:created_at" json:"created_at"``).
* **Trường ID:** Bắt buộc sử dụng khóa chính dạng tự động tăng (`autoIncrement`) hoặc chuỗi định danh duy nhất (UUID/ULID) tùy thuộc vào thiết kế thực thể.
* **Tên Table thực tế trong DB:** Định nghĩa tường minh thông qua hàm `TableName() string` cho từng struct theo quy chuẩn snake_case, số nhiều (VD: `users`, `products`, `point_transactions`).

## 3. RÀNG BUỘC & QUAN HỆ (RELATIONS & CONSTRAINTS)

* Mọi quan hệ 1-N hoặc N-N BẮT BUỘC phải cấu hình tường minh bằng các thẻ GORM tags (`foreignKey`, `references`, `many2many`).
* BẮT BUỘC xử lý hành vi xóa và lưu vết: Tích hợp GORM Soft Delete bằng cách nhúng trường `DeletedAt gorm.DeletedAt `gorm:"index"`` vào các thực thể quan trọng để tránh mất mát dữ liệu của doanh nghiệp. Sử dụng ràng buộc `OnDelete:CASCADE` thông qua GORM tags chỉ đối với các dữ liệu phụ thuộc hoàn toàn vào thực thể cha.

## 4. TỐI ƯU HIỆU SUẤT (INDEXING)

* Tự động đánh chỉ mục bằng cách khai báo thuộc tính `index` hoặc `uniqueIndex` trong tag `gorm` cho các trường thường xuyên xuất hiện trong mệnh đề `WHERE` (ví dụ: các trường tìm kiếm, các trường lọc trạng thái, mã coupon, token identifiers).
* Đánh chỉ mục hỗn hợp (Composite Index) cho các cụm dữ liệu hay truy vấn đồng thời để tối ưu hiệu năng cho hệ thống Enterprise.

## 5. THỰC THI GHI FILE (FILE EXECUTION - LỆNH BẮT BUỘC)

* BẠN BỊ CẤM chỉ in code ra màn hình chat.
* BẮT BUỘC phải thực hiện hành động **Tạo mới** hoặc **Ghi đè/Cập nhật (Patch)** đoạn code Struct Models và file Migration vừa sinh ra vào đúng các file vật lý trên hệ thống.
* Tuyệt đối không làm hỏng, xóa mất hoặc ghi đè sai các Struct/Bảng dữ liệu đã tồn tại trước đó trong file nguồn.

## 6. BÁO CÁO KẾT QUẢ

* Sau khi ghi file thành công, in ra một thông báo ngắn gọn: *"✅ Đã cập nhật thành công thiết kế Struct Models và Migration vào file [Tên các đường dẫn file vật lý]."*
* Liệt kê ngắn gọn các Struct mới, các mối quan hệ (Relations) và các trường được đánh chỉ mục (Index) vừa thêm vào để System Architect xem xét.