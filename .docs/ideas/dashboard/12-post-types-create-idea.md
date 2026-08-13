# Ý TƯỞNG: Quản trị Loại bài viết (Post Type Management)

## 1. Thông tin chung (Meta Info)

* **Dự án:** TechBite.
* **Tính năng:** Quản lý loại bài viết (Post Type Management).
* **Mục đích:** Xác định và phân loại các loại bài viết khác nhau trong hệ thống, giúp tổ chức nội dung và áp dụng các quy tắc hiển thị phù hợp cho từng loại.

---

# 2. Đối tượng & Trải nghiệm (Target & UX)

* **Người dùng chính:** Quản trị viên (Admin), Biên tập viên (Editor).

### Hành động chính

* Xem danh sách các loại bài viết.
* Thêm mới loại bài viết.
* Chỉnh sửa thông tin loại bài viết.
* Sắp xếp thứ tự hiển thị.
* Thay đổi trạng thái kích hoạt.
* Xóa loại bài viết (nếu chưa được sử dụng).

### Cảm xúc mang lại

Giao diện quản lý trực quan, dễ dàng phân loại và tìm kiếm. Hỗ trợ thao tác nhanh với các loại bài viết phổ biến đã được định nghĩa sẵn.

---

# 3. Đặc tả Thiết kế (Design Specs)

## Phong cách UI

* Bảng danh sách (Table) với các cột: ID, Code, Tên loại, Trạng thái, Thứ tự, Thao tác.
* Sử dụng Badge để hiển thị trạng thái.
* Trang riêng biệt (Page) để thêm/sửa loại bài viết.
* Drag-and-drop để sắp xếp thứ tự (sort_order).
* Responsive trên mọi thiết bị.

---

## Màu sắc chủ đạo (Brand Colors)

* **Cam thương hiệu:** Nút tạo mới `bg-[#ff8c42]`
* **Xanh lá:** ACTIVE
* **Đỏ:** INACTIVE
* **Xanh dương:** ID, Code
* **Xám:** Text thông thường

---

## Cấu trúc màn hình (Top to Bottom)

### Header & Actions

* Tiêu đề: `Quản lý loại bài viết`
* Nút: `Thêm loại bài viết mới`
* Breadcrumb: `Trang chủ > Quản trị nội dung > Loại bài viết`

### Bảng danh sách (Table)

| ID | Code | Tên loại | Trạng thái | Thứ tự | Thao tác |
|----|------|----------|------------|---------|----------|
| 1  | NEWS | Tin tức  | ACTIVE  | 1       | Sửa - Xóa   |
| 2  | SERVICE | Dịch vụ | ACTIVE | 2 | Sửa - Xóa |
| 3  | POLICY | Chính sách | ACTIVE | 3 | Sửa - Xóa |
| 4  | GUIDE | Hướng dẫn | ACTIVE | 4 | Sửa - Xóa |
| 5  | ABOUT | Giới thiệu | ACTIVE | 5 | Sửa - Xóa |

### Form thêm/sửa (Page riêng biệt)

* **Mã loại (Code)***: Text input, chỉ cho phép chữ hoa, số và dấu gạch dưới.
* **Tên loại (Name)***: Text input.
* **Trạng thái (Status)**: Switch/Radio (ACTIVE/INACTIVE).
* **Thứ tự (Sort Order)**: Number input, tự động tăng.

---

# 4. Kiểm tra dữ liệu (Validation & Error Handling)

### Code

* Bắt buộc.
* Chỉ chứa chữ hoa (A-Z), số (0-9) và dấu gạch dưới (_).
* Không được trùng trong hệ thống.
* Tối đa 50 ký tự.
* Không chứa khoảng trắng.

### Name

* Bắt buộc.
* Tối đa 100 ký tự.

### Status

* Chỉ nhận: `ACTIVE` hoặc `INACTIVE`.

### Sort Order

* Phải là số nguyên dương.
* Không được trùng trong cùng hệ thống (tự động điều chỉnh).

---

# 5. Chức năng nổi bật

* **Tìm kiếm:** Tìm kiếm theo Code hoặc Name.
* **Lọc:** Lọc theo trạng thái (ACTIVE/INACTIVE).
* **Sắp xếp:** Click vào tiêu đề cột để sắp xếp.
* **Drag-and-drop:** Kéo thả để thay đổi thứ tự.
* **Xóa mềm:** Kiểm tra ràng buộc trước khi xóa (không xóa nếu đã có bài viết sử dụng).
* **Bulk Actions:** Xóa nhiều, thay đổi trạng thái hàng loạt.

---

# 6. Quy tắc nghiệp vụ (Business Rules)

* Mỗi bài viết (Post) bắt buộc phải thuộc một `post_type`.
* Không thể xóa `post_type` đang được sử dụng bởi bất kỳ bài viết nào.
* `sort_order` mặc định tăng dần theo thứ tự tạo mới.
* Khi thay đổi `sort_order`, hệ thống tự động cập nhật lại thứ tự cho toàn bộ danh sách.
* `code` là duy nhất và được sử dụng trong URL hoặc API để phân loại.

---

# 7. API

## Lấy danh sách loại bài viết

```
GET /post-types
```

Response:

```json
{
  "data": [
    {
      "id": 1,
      "code": "NEWS",
      "name": "Tin tức",
      "status": "ACTIVE",
      "sort_order": 1
    },
    {
      "id": 2,
      "code": "SERVICE",
      "name": "Dịch vụ",
      "status": "ACTIVE",
      "sort_order": 2
    }
  ]
}
```

## Tạo mới loại bài viết

```
POST /post-types
```

Request Body:

```json
{
  "code": "RECRUITMENT",
  "name": "Tuyển dụng",
  "status": "ACTIVE",
  "sort_order": 6
}
```

## Cập nhật loại bài viết

```
PUT /post-types/{id}
```

## Xóa loại bài viết

```
DELETE /post-types/{id}
```

---

# 8. Chức năng sau khi thao tác

## Thành công

* Hiển thị Toast: `Thêm loại bài viết thành công.`
* Tự động refresh danh sách.

## Thất bại

Hiển thị lỗi từ Backend:

```
Mã loại đã tồn tại trong hệ thống.
```
hoặc
```
Không thể xóa loại bài viết đang được sử dụng.
```

---

# 9. Định hướng UI Component

### Page Components
* Page Header
* Search & Filter Bar
* Action Button (Create New)

### Table Components
* Data Table with Pagination
* Status Badge
* Action Buttons (Edit, Delete)

### Form Components
* Text Field (Code, Name)
* Number Input (Sort Order)
* Switch/Radio (Status)

### Feedback Components
* Toast/Snackbar
* Confirm Dialog (Delete)
* Loading Spinner

---

**Dữ liệu mẫu cho post_types:**

```sql
INSERT INTO post_types (id, code, name, status, sort_order) VALUES
(1, 'NEWS', 'Tin tức', 'ACTIVE', 1),
(2, 'SERVICE', 'Dịch vụ', 'ACTIVE', 2),
(3, 'POLICY', 'Chính sách', 'ACTIVE', 3),
(4, 'GUIDE', 'Hướng dẫn', 'ACTIVE', 4),
(5, 'ABOUT', 'Giới thiệu', 'ACTIVE', 5),
(6, 'RECRUITMENT', 'Tuyển dụng', 'INACTIVE', 6),
(7, 'EVENT', 'Sự kiện', 'ACTIVE', 7),
(8, 'PARTNER', 'Đối tác', 'ACTIVE', 8);
```