# Ý TƯỞNG: Thêm danh mục sản phẩm (Create Category)

## 1. Thông tin chung

* **Dự án:** TechBite
* **Tính năng:** Thêm danh mục sản phẩm
* **Mục đích:** Tạo mới danh mục sản phẩm để tổ chức menu theo cấu trúc phân cấp.

---

## 2. Đối tượng

* Quản trị viên (Admin)
* Nhân viên quản lý thực đơn (Menu Planner)

---

## 3. Bố cục màn hình

### Header

* Tiêu đề: **Thêm danh mục**
* Breadcrumb:

  ```
  Trang chủ > Cấu hình hệ thống > Danh mục > Thêm mới
  ```

### Form thông tin

#### Thông tin cơ bản

* Tên danh mục (*)
* Slug (tự động sinh từ tên, cho phép chỉnh sửa)
* Danh mục cha (Dropdown Tree, mặc định: Không có)
* Mô tả

#### Hình ảnh

* Upload ảnh đại diện
* Hỗ trợ Preview, Thay ảnh, Xóa ảnh

#### Thiết lập

* Thứ tự hiển thị
* Trạng thái

  * Hoạt động
  * Ẩn

### Footer

* Hủy
* Lưu
* Lưu & Thêm mới

* Lưu ý: sử dụng lại Header và Sidebar chung đã có của trang
---

## 4. Validate

| Trường       | Quy tắc                                    |
| ------------ | ------------------------------------------ |
| Tên danh mục | Bắt buộc                                   |
| Slug         | Không trùng, chỉ gồm chữ thường, số và "-" |
| Danh mục cha | Không được tạo vòng lặp                    |
| Thứ tự       | ≥ 0                                        |
| Ảnh          | jpg, jpeg, png, webp; ≤ 5MB                |

---

## 5. Business Rules

* Tự động sinh Slug khi nhập tên.
* Có thể chỉnh sửa Slug trước khi lưu.
* Không cho phép chọn danh mục cha tạo cấu trúc vòng lặp.
* Danh mục mặc định là **Hoạt động**.
* Sau khi tạo thành công:

  * **Lưu:** Quay về danh sách.
  * **Lưu & Thêm mới:** Reset form để tiếp tục tạo danh mục khác.

---

## 6. API

* `GET /api/v1/admin/categories/tree` (Danh sách danh mục cha)
* `POST /api/v1/admin/categories`
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
