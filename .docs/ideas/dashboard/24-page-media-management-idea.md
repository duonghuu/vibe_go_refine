# Ý TƯỞNG: Quản lý hình ảnh Trang (Page Media)

## 1. Thông tin chung

* **Dự án:** TechBite.
* **Tính năng:** Liên kết Page và Media.
* **Mục đích:** Cho phép mỗi trang có một ảnh đại diện và một gallery ảnh sắp xếp được.
* **Phạm vi UI:** Không có màn hình độc lập; Media Card được nhúng trong màn hình tạo/sửa Page.

---

## 2. Dữ liệu & collections

* Bảng liên kết riêng: `page_media`; không dùng `post_media` vì bảng hiện tại có foreign key tới `posts`.
* `page_id`: Page sở hữu liên kết.
* `media_id`: File Media đã upload.
* `collection`: chỉ nhận `thumbnail` hoặc `gallery`.
* `sort_order`: xác định thứ tự trong gallery.

---

## 3. Quy tắc nghiệp vụ

* Một Page có tối đa một `thumbnail`; upload/chọn thumbnail mới thay thế liên kết cũ.
* Gallery nhận nhiều ảnh, hỗ trợ kéo thả để thay đổi `sort_order` và đồng bộ toàn bộ collection.
* Chỉ nhận Media là ảnh; Admin có thể dùng mọi Media hợp lệ, Staff chỉ gắn Media do mình sở hữu theo quy tắc Media hiện hữu.
* Gỡ ảnh chỉ gỡ liên kết với Page, không xóa file Media gốc.
* Thay đổi Media phải làm mới cache public của Page.

---

## 4. API dự kiến

```text
GET    /api/v1/admin/pages/:page_id/media?collection=thumbnail|gallery
POST   /api/v1/admin/pages/:page_id/media
PUT    /api/v1/admin/pages/:page_id/media/:collection
DELETE /api/v1/admin/pages/:page_id/media/:media_id?collection=thumbnail|gallery
POST   /api/v1/media/upload
```

Đồng bộ gallery:

```json
{
  "media": [
    { "id": 102, "sortOrder": 0 },
    { "id": 103, "sortOrder": 1 }
  ]
}
```

---

## 5. Trải nghiệm UI

* Thumbnail hiển thị preview, thay thế và gỡ bỏ.
* Gallery hiển thị preview, trạng thái upload, lỗi upload, gỡ từng ảnh và kéo thả để sắp xếp.
* Trong lúc Page chưa được tạo, card vô hiệu hóa thao tác liên kết và hướng dẫn lưu Page trước.

