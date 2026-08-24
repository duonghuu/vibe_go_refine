# Ý TƯỞNG: Chỉnh sửa danh mục bài viết (Edit Post Category)

## 1. Thông tin chung (Meta Info)

* **Dự án:** TechBite
* **Tính năng:** Chỉnh sửa danh mục bài viết (Post Category Edit)
* **Mục đích:** Chỉnh sửa thông tin chi tiết của một danh mục bài viết hiện có.

---

## 2. Đối tượng & Trải nghiệm (Target & UX)

* **Người dùng chính:** Quản trị viên (Admin), Biên tập viên (Editor).

### Trải nghiệm
Giao diện form nhập liệu rõ ràng, hiển thị sẵn thông tin cũ. Người dùng có thể chỉnh sửa các thông tin như tên, mô tả, danh mục cha, trạng thái. Danh sách "Danh mục cha" loại trừ chính danh mục đang sửa và các danh mục con của nó để tránh vòng lặp.

---

## 3. Bố cục màn hình (Design Specs)

### Header & Actions

* Tiêu đề: **Chỉnh sửa danh mục bài viết**
* Breadcrumb: `Trang chủ > Quản trị nội dung > Danh mục > Chỉnh sửa`

### Form thông tin

#### 3.1. Phân loại & Thông tin cơ bản
* **Loại bài viết (Post Type):** Hiển thị dạng Text hoặc Read-only. Không cho phép đổi loại bài viết của một danh mục đã tạo.
* **Tên danh mục (*):** Text input.
* **Slug:** Text input.
* **Danh mục cha:** Dropdown Tree (Mặc định: Không có/Root). Chỉ hiển thị các danh mục thuộc cùng `Post Type`, KHÔNG hiển thị danh mục đang sửa và các danh mục con của nó.
* **Mô tả:** Textarea (tùy chọn).

#### 3.2. Hình ảnh (Tùy chọn)
* **Ảnh đại diện:** Upload/Select từ Media Library.
* Hỗ trợ Preview, Thay ảnh, Xóa ảnh. Hiển thị ảnh cũ nếu có.

#### 3.3. Thiết lập hiển thị
* **Thứ tự hiển thị (Sort Order):** Number input.
* **Trạng thái (Status):** Switch/Radio.
  * `ACTIVE` (Hoạt động)
  * `INACTIVE` (Ẩn)

### Footer
* Nút **Hủy** (quay về danh sách).
* Nút **Lưu** (lưu và quay về danh sách).

---

## 4. Kiểm tra dữ liệu (Validation & Error Handling)

| Trường | Quy tắc |
| :--- | :--- |
| **Tên danh mục** | Bắt buộc, tối đa 255 ký tự. |
| **Slug** | Bắt buộc, không trùng lặp trong cùng một hệ thống (ngoại trừ chính nó), chỉ chứa chữ thường, số, dấu `-`. |
| **Danh mục cha** | Không được tạo vòng lặp (Circular reference). Không được chọn chính nó hoặc các danh mục con của nó làm cha. |
| **Thứ tự** | Là số nguyên >= 0. |
| **Ảnh đại diện** | Định dạng jpg, jpeg, png, webp. Dung lượng <= 5MB. |

---

## 5. Quy tắc nghiệp vụ (Business Rules)

* **Không đổi Post Type:** Không cho phép thay đổi Loại bài viết của một danh mục đã tồn tại để đảm bảo tính toàn vẹn dữ liệu của các bài viết thuộc danh mục đó.
* **Chống vòng lặp danh mục:** Danh sách Dropdown chọn danh mục cha phải loại trừ danh mục hiện tại và toàn bộ nhánh con của nó.
* **Slug:** Không auto-generate như lúc tạo mới để tránh tự ý làm hỏng URL đã SEO (chỉ
thay đổi khi user chủ động nhập).
---

## 6. API Endpoints

* `GET /api/v1/admin/post-categories/{id}` (Lấy chi tiết danh mục để điền vào form).
* `GET /api/v1/admin/post-categories/tree?type_code={code}` (Lấy danh sách danh mục cha theo loại bài viết, FE có thể tự filter bỏ các node không hợp lệ).
* `PUT /api/v1/admin/post-categories/{id}` (Submit cập nhật danh mục).
* `POST /api/v1/media/upload` (API upload ảnh dùng chung nếu có).

---

## 7. Cấu trúc Request Body (Cập nhật)

```json
{
  "name": "Tin trong nước (Update)",
  "slug": "tin-trong-nuoc",
  "parent_id": null,
  "description": "Các tin tức thời sự trong nước",
  "image_url": "https://example.com/images/post-categories/tin-trong-nuoc.jpg",
  "sort_order": 2,
  "status": "ACTIVE"
}
```

---

## 8. Chức năng sau khi thao tác

### Thành công
* Hiển thị thông báo (Toast): `Cập nhật danh mục bài viết thành công.`
* Chuyển hướng quay về danh sách.

### Thất bại
* Hiển thị Toast lỗi trả về từ Backend (ví dụ: `Slug đã tồn tại`, `Danh mục cha không hợp lệ`).
