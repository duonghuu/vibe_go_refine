# Ý TƯỞNG: Chỉnh sửa trang tĩnh (Admin Edit Page)

## 1. Thông tin chung

* **Dự án:** TechBite.
* **Tính năng:** Chỉnh sửa Trang (Page).
* **Mục đích:** Cho phép cập nhật nội dung, đường dẫn và trạng thái công khai của một trang đã tạo.

---

## 2. Trải nghiệm

* Form hiển thị dữ liệu hiện có ngay khi tải xong.
* Người dùng có thể chuyển đổi giữa nháp và công khai mà không tạo bản ghi Page khác.
* Slug chỉ thay đổi khi người dùng chủ động sửa, tránh làm đổi URL ngoài ý muốn.

---

## 3. Đặc tả màn hình

* Tiêu đề: **Chỉnh sửa trang**.
* Breadcrumb: `Trang chủ > Quản trị nội dung > Trang > Chỉnh sửa`.
* Bố cục và field giống màn hình Tạo trang: title, slug, rich content, status, thumbnail/gallery và SEO.
* Hiển thị preview URL public theo slug, nhưng không tự dựng hoặc thay đổi menu website.
* Nút **Lưu thay đổi** giữ trạng thái đang chọn; **Hủy** quay về danh sách.

---

## 4. Quy tắc nghiệp vụ

* Page không được đổi sang Post hoặc gắn `post_type`.
* Page bị soft-delete hoặc không tồn tại trả trạng thái không tìm thấy.
* Cập nhật title/content/slug/status làm mới cache trang public.
* Đổi slug làm URL cũ không còn hợp lệ; phiên bản đầu không tạo redirect URL cũ.
* SEO và Media là các thao tác con của Page hiện tại, chỉ thực hiện khi `id` hợp lệ.

---

## 5. API dự kiến

```text
GET /api/v1/admin/pages/:id
PUT /api/v1/admin/pages/:id
```

```json
{
  "title": "Giới thiệu về TechBite",
  "slug": "gioi-thieu",
  "content": "<p>Nội dung đã cập nhật</p>",
  "status": "PUBLISHED"
}
```

---

## 6. Trạng thái phản hồi

* Loading: skeleton cho form và các card liên quan.
* Error: thông báo lỗi, có thao tác thử lại khi không tải được Page/Media/SEO.
* Success: toast thành công, refresh dữ liệu Page và URL preview.

