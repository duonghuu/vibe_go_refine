# Ý TƯỞNG: Chỉnh sửa danh mục sản phẩm (Edit Category)

## 1. Thông tin chung

* **Dự án:** TechBite
* **Tính năng:** Chỉnh sửa danh mục sản phẩm
* **Mục đích:** Cập nhật thông tin danh mục sản phẩm hiện có.

---

## 2. Đối tượng

* Quản trị viên (Admin)
* Nhân viên quản lý thực đơn (Menu Planner)

---

## 3. Bố cục màn hình

### Header

* Tiêu đề: **Chỉnh sửa danh mục**
* Breadcrumb:

  ```
  Trang chủ > Cấu hình hệ thống > Danh mục > Chỉnh sửa
  ```

### Form thông tin

#### Thông tin cơ bản

* Tên danh mục (*)
* Slug
* Danh mục cha (Dropdown Tree)
* Mô tả

#### Hình ảnh

* Hiển thị ảnh hiện tại
* Thay ảnh
* Xóa ảnh

#### Thiết lập

* Thứ tự hiển thị
* Trạng thái

  * Hoạt động
  * Ẩn

#### Thông tin hệ thống (Chỉ đọc)

* ID
* Ngày tạo
* Ngày cập nhật
* Người tạo
* Người cập nhật

### Footer

* Hủy
* Lưu thay đổi

---

## 4. Validate

| Trường       | Quy tắc                                    |
| ------------ | ------------------------------------------ |
| Tên danh mục | Bắt buộc                                   |
| Slug         | Không trùng, chỉ gồm chữ thường, số và "-" |
| Danh mục cha | Không được chọn chính nó hoặc tạo vòng lặp |
| Thứ tự       | ≥ 0                                        |
| Ảnh          | jpg, jpeg, png, webp; ≤ 5MB                |

---

## 5. Business Rules

* Load dữ liệu hiện tại khi mở màn hình.
* Cho phép thay đổi tất cả thông tin ngoại trừ **ID**.
* Nếu thay đổi **Tên**, hệ thống vẫn giữ **Slug** hiện tại; chỉ thay đổi khi người dùng chỉnh sửa Slug.
* Không cho phép chọn chính danh mục đang sửa hoặc danh mục con của nó làm danh mục cha.
* Khi chuyển trạng thái từ **Hoạt động** sang **Ẩn**, hiển thị cảnh báo nếu danh mục đang chứa sản phẩm hoặc có danh mục con.
* Cảnh báo khi rời trang nếu có thay đổi chưa lưu (Unsaved Changes).
* Sau khi lưu thành công, quay về màn hình danh sách danh mục.

---

## 6. API

* `GET /api/v1/admin/categories/{id}`
* `GET /api/v1/admin/categories/tree` (Danh sách danh mục cha)
* `PUT /api/v1/admin/categories/{id}`
* `POST /api/v1/media/upload` (Upload ảnh)

---

## 7. Request Body

```json
{
  "name": "Burger",
  "slug": "burger",
  "parentId": null,
  "description": "Các món Burger",
  "imageUrl": "https://example.com/images/categories/mobile.png",
  "sortOrder": 1,
  "status": "ACTIVE"
}
```