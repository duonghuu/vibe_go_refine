# Ý TƯỞNG: Quản trị Danh sách trang tĩnh (Admin Page List)

## 1. Thông tin chung

* **Dự án:** TechBite.
* **Tính năng:** Quản lý danh sách Trang (Page).
* **Mục đích:** Cho phép Admin và Staff quản lý các trang nội dung tĩnh như Giới thiệu, Liên hệ và Chính sách.
* **Ranh giới kiến trúc:** `Page` là thực thể độc lập, không phải một `post_type` và không thuộc `posts`.

---

## 2. Đối tượng & trải nghiệm

* **Người dùng chính:** Admin, Staff.
* **Hành động chính:** Xem, tìm kiếm, lọc theo trạng thái, tạo mới, chỉnh sửa và xóa mềm trang.
* **Trải nghiệm mong muốn:** Bảng quản trị gọn, phản hồi nhanh, nhận biết rõ trang nào đang là nháp hoặc đã công khai.

---

## 3. Đặc tả màn hình

### Header

* Tiêu đề: **Quản lý trang**.
* Breadcrumb: `Trang chủ > Quản trị nội dung > Trang`.
* Nút chính: **Tạo trang mới**.

### Thanh công cụ

* Ô tìm kiếm theo `title` hoặc `slug`, debounce 300ms.
* Bộ lọc trạng thái: `DRAFT`, `PUBLISHED`.
* Trạng thái lọc, phân trang và sắp xếp đồng bộ với URL.

### DataGrid

| Cột | Nội dung |
| --- | --- |
| ID | Định danh trang. |
| Tiêu đề | Nhấn để mở màn hình chỉnh sửa. |
| Slug | Đường dẫn public, hiển thị rút gọn nếu dài. |
| Trạng thái | Badge `DRAFT` hoặc `PUBLISHED`. |
| Cập nhật lúc | Thời điểm cập nhật gần nhất. |
| Thao tác | Sửa và Xóa. |

* DataGrid dùng mật độ `dense` hoặc `medium`.
* Khi danh sách rỗng, hiển thị empty state cùng nút tạo trang mới.
* Xóa luôn yêu cầu Confirm Dialog trước khi gọi API.

---

## 4. Quy tắc nghiệp vụ

* `DRAFT` chỉ xuất hiện trong CMS; `PUBLISHED` có thể được public API trả về.
* Xóa là soft delete, không hiển thị ở danh sách mặc định hay public API.
* Không có phân loại, danh mục, post type, trang cha-con hoặc cấu hình menu trong phiên bản đầu.
* Tất cả thao tác CMS yêu cầu Access Token hợp lệ; quyền ghi dành cho `ADMIN` và `STAFF`.

---

## 5. API dự kiến

```text
GET    /api/v1/admin/pages?current=1&pageSize=10&title_like=&status=
DELETE /api/v1/admin/pages/:id
```

Response danh sách tương thích Refine:

```json
{
  "data": [
    {
      "id": 1,
      "title": "Giới thiệu",
      "slug": "gioi-thieu",
      "status": "PUBLISHED",
      "authorId": 2,
      "createdAt": "2026-08-31T10:00:00Z",
      "updatedAt": "2026-08-31T10:30:00Z"
    }
  ],
  "total": 1
}
```

---

## 6. Trạng thái phản hồi

* Loading: hiển thị loading state của DataGrid.
* Error: hiển thị thông báo lỗi và cho phép tải lại.
* Success: toast sau khi xóa thành công và refresh danh sách.

