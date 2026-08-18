# Ý TƯỞNG: Thêm danh mục bài viết (Create Post Category)

## 1. Thông tin chung (Meta Info)

* **Dự án:** TechBite
* **Tính năng:** Thêm danh mục bài viết (Post Category Create)
* **Mục đích:** Tạo mới danh mục chuyên dành cho các loại bài viết (News, Services, Guides,...), tách biệt hoàn toàn với danh mục sản phẩm. Việc phân loại sẽ phụ thuộc vào từng loại bài viết (Post Type).

---

## 2. Đối tượng & Trải nghiệm (Target & UX)

* **Người dùng chính:** Quản trị viên (Admin), Biên tập viên (Editor).

### Trải nghiệm
Giao diện form nhập liệu rõ ràng, tiện lợi. Hỗ trợ tự động sinh slug từ tên danh mục. Khi chọn Loại bài viết (Post Type), hệ thống tự động tải danh sách "Danh mục cha" tương ứng của loại bài viết đó để đảm bảo tính phân cấp logic.

---

## 3. Bố cục màn hình (Design Specs)

### Header & Actions

* Tiêu đề: **Thêm danh mục bài viết**
* Breadcrumb: `Trang chủ > Quản trị nội dung > Danh mục > Thêm mới`

### Form thông tin

#### 3.1. Phân loại & Thông tin cơ bản
* **Loại bài viết (Post Type) (*):** được hệ thống tự động nhận diện thông qua URL (Path Params hoặc Query Params), không yêu cầu người dùng chọn. 
  * *Hành vi:* Bắt buộc. Loại bài viết khác nhau, danh sách "Danh mục cha" sẽ được lấy lại tương ứng.
* **Tên danh mục (*):** Text input.
* **Slug:** Text input (tự động sinh từ tên danh mục, cho phép tùy chỉnh).
* **Danh mục cha:** Dropdown Tree (Mặc định: Không có/Root). Chỉ hiển thị các danh mục thuộc `Post Type` đã chọn.
* **Mô tả:** Textarea (tùy chọn).

#### 3.2. Hình ảnh (Tùy chọn)
* **Ảnh đại diện:** Upload/Select từ Media Library.
* Hỗ trợ Preview, Thay ảnh, Xóa ảnh.

#### 3.3. Thiết lập hiển thị
* **Thứ tự hiển thị (Sort Order):** Number input, tự động điền (ví dụ: tự động lấy max + 1 hoặc mặc định là 0).
* **Trạng thái (Status):** Switch/Radio.
  * `ACTIVE` (Hoạt động)
  * `INACTIVE` (Ẩn)

### Footer
* Nút **Hủy** (quay về danh sách).
* Nút **Lưu** (lưu và quay về danh sách).
* Nút **Lưu & Thêm mới** (lưu và reset form để tạo tiếp).

---

## 4. Kiểm tra dữ liệu (Validation & Error Handling)

| Trường | Quy tắc |
| :--- | :--- |
| **Loại bài viết** | Bắt buộc chọn. |
| **Tên danh mục** | Bắt buộc, tối đa 255 ký tự. |
| **Slug** | Bắt buộc, không trùng lặp trong cùng một hệ thống (hoặc cùng Post Type), chỉ chứa chữ thường, số, dấu `-`. |
| **Danh mục cha** | Không được tạo vòng lặp (Circular reference). Chỉ được chọn danh mục cùng Post Type. |
| **Thứ tự** | Là số nguyên $\ge 0$. |
| **Ảnh đại diện** | Định dạng jpg, jpeg, png, webp. Dung lượng $\le 5$MB. |

---

## 5. Quy tắc nghiệp vụ (Business Rules)

* **Ràng buộc với Post Type:** Danh mục bài viết được nhóm theo loại bài viết. Một danh mục "Tin trong nước" phải thuộc loại "Tin tức (NEWS)". 
* **Auto Slug:** Khi người dùng gõ "Tên danh mục", slug sẽ được tạo tự động real-time nếu người dùng chưa can thiệp chỉnh sửa thủ công vào trường slug.
* **Cây phân cấp (Hierarchy):** Khi tạo một danh mục con, nó kế thừa Post Type của danh mục cha. Form cần khóa chặt logic này bằng cách load đúng Cây danh mục cha theo Post Type được chọn.

---

## 6. API Endpoints

* `GET /api/v1/admin/post-types` (Lấy danh sách Post Types để điền vào Dropdown).
* `GET /api/v1/admin/post-categories/tree?type_code={code}` (Lấy danh sách danh mục cha theo loại bài viết).
* `POST /api/v1/admin/post-categories` (Submit tạo mới danh mục).
* `POST /api/v1/media/upload` (API upload ảnh dùng chung nếu có).

---

## 7. Cấu trúc Request Body

```json
{
  "type_code": "NEWS",
  "name": "Tin trong nước",
  "slug": "tin-trong-nuoc",
  "parent_id": null,
  "description": "Các tin tức thời sự trong nước",
  "image_url": "https://example.com/images/post-categories/tin-trong-nuoc.jpg",
  "sort_order": 1,
  "status": "ACTIVE"
}
```

---

## 8. Chức năng sau khi thao tác

### Thành công
* Hiển thị thông báo (Toast): `Thêm danh mục bài viết thành công.`
* Chuyển hướng hoặc reset form dựa theo nút người dùng bấm (Lưu vs Lưu & Thêm mới).

### Thất bại
* Hiển thị Toast lỗi trả về từ Backend (ví dụ: `Slug đã tồn tại`, `Loại bài viết không hợp lệ`).
