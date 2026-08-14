# Ý TƯỞNG: Liên kết Bài viết và Media (Post Media)

## 1. Thông tin chung (Meta Info)

* **Dự án:** TechBite.
* **Tính năng:** Liên kết bài viết và hình ảnh/tài liệu (Post Media).
* **Mục đích:** Quản lý mối quan hệ n-n (hoặc 1-n) giữa bài viết (Post) và các tập tin đa phương tiện (Media). Cho phép gắn hình ảnh/tài liệu vào bài viết theo các bộ sưu tập (collection) khác nhau để phục vụ hiển thị ở UI.
* **Lưu ý quan trọng:** Tính năng này chỉ xử lý ở Backend API, **không cần giao diện (UI) độc lập**. Các thao tác liên kết media sẽ được tích hợp trực tiếp vào màn hình Thêm/Sửa Bài viết.

---

# 2. Đặc tả dữ liệu (Data Specs)

Bảng trung gian `post_media`:
* **`post_id`**: ID của bài viết (Khóa ngoại liên kết tới bảng `posts`).
* **`media_id`**: ID của tập tin đa phương tiện (Khóa ngoại liên kết tới bảng `media`).
* **`collection`**: Phân loại tập tin theo mục đích sử dụng.
* **`sort_order`**: Số thứ tự để sắp xếp các media trong cùng một collection.

Các giá trị hợp lệ của `collection`:
* `thumbnail`: Hình thu nhỏ đại diện cho bài viết.
* `gallery`: Hình ảnh trong bộ sưu tập (slider, album) của bài viết.
* `attachment`: Tài liệu đính kèm (pdf, doc, v.v.) để người dùng tải về.
* `content`: Hình ảnh được chèn trực tiếp trong nội dung bài viết (thông qua WYSIWYG editor).

---

# 3. Kiểm tra dữ liệu (Validation & Error Handling)

### post_id
* Bắt buộc (Required).
* Phải tồn tại trong hệ thống (tồn tại trong bảng `posts`).

### media_id
* Bắt buộc (Required).
* Phải tồn tại trong hệ thống (tồn tại trong bảng `media`).

### collection
* Bắt buộc (Required).
* Chỉ nhận các giá trị: `thumbnail`, `gallery`, `attachment`, `content`.

### sort_order
* Phải là số nguyên (Integer).
* Mặc định là 0 hoặc tự động tăng theo thứ tự thêm vào nếu không truyền.

---

# 4. Quy tắc nghiệp vụ (Business Rules)

* **Xóa bài viết (Cascade Delete):** Khi một bài viết (`post_id`) bị xóa, toàn bộ liên kết tương ứng trong bảng `post_media` cũng phải bị xóa.
* **Xóa media (Cascade Delete):** Khi một tập tin media (`media_id`) bị xóa khỏi hệ thống, các liên kết của nó trong `post_media` cũng phải bị xóa theo.
* **Quy tắc Thumbnail:** Thông thường mỗi bài viết chỉ có duy nhất 1 `thumbnail`. Nếu người dùng cập nhật `thumbnail`, hệ thống cần gỡ bỏ (hoặc thay thế) liên kết `thumbnail` cũ của bài viết đó.
* **Quy tắc Gallery & Attachment:** Một bài viết có thể có nhiều `gallery` và `attachment`. Cần chú ý xử lý `sort_order` để khi hiển thị ra ngoài Frontend, danh sách media được sắp xếp đúng như người dùng đã thao tác.
* **Cơ chế Sync (Đồng bộ):** Để tiện cho thao tác từ Frontend, Backend nên hỗ trợ cơ chế Sync (gửi toàn bộ danh sách ID cần gắn, Backend tự động tính toán để thêm mới liên kết chưa có và xóa liên kết không còn trong danh sách).

---

# 5. Đặc tả API (API Specs)

*API này có thể đứng độc lập (sub-resource của Post) hoặc được gộp chung vào payload khi tạo/cập nhật Post.*

## 5.1. Lấy danh sách Media của một bài viết

```
GET /posts/{post_id}/media
```
*(Query param tuỳ chọn: `?collection=gallery` để chỉ lấy media thuộc collection cụ thể)*

**Response (Thành công 200):**

```json
{
  "data": [
    {
      "post_id": 1,
      "media_id": 101,
      "collection": "thumbnail",
      "sort_order": 1,
      "media": {
        "id": 101,
        "url": "https://example.com/thumbnail.jpg",
        "file_name": "thumbnail.jpg",
        "mime_type": "image/jpeg"
      }
    },
    {
      "post_id": 1,
      "media_id": 102,
      "collection": "gallery",
      "sort_order": 1,
      "media": {
        "id": 102,
        "url": "https://example.com/gallery1.jpg",
        "file_name": "gallery1.jpg",
        "mime_type": "image/jpeg"
      }
    }
  ]
}
```

## 5.2. Liên kết Media vào bài viết

```
POST /posts/{post_id}/media
```

**Request Body:**

```json
{
  "media_id": 101,
  "collection": "thumbnail",
  "sort_order": 1
}
```

## 5.3. Đồng bộ (Sync) Media cho một bộ sưu tập (Collection)

*API này thường dùng nhất để cập nhật danh sách hình ảnh (ví dụ: cập nhật lại toàn bộ Gallery).*

```
PUT /posts/{post_id}/media/{collection}
```

**Request Body:**
*(Danh sách các media thuộc collection này. Những media nào trước đó thuộc collection này nhưng không có trong mảng dưới đây sẽ bị xóa liên kết khỏi bài viết).*

```json
{
  "media": [
    { "id": 102, "sort_order": 1 },
    { "id": 103, "sort_order": 2 },
    { "id": 105, "sort_order": 3 }
  ]
}
```

## 5.4. Gỡ liên kết Media khỏi bài viết

```
DELETE /posts/{post_id}/media/{media_id}
```
*(Hoặc có thể truyền thêm query param `?collection=gallery` nếu một media_id có thể thuộc nhiều collection khác nhau trong cùng 1 bài viết).*

---

# 6. Dữ liệu mẫu (Sample Data)

```sql
-- Bảng trung gian: post_media
-- (Giả định post có id = 10, media có id từ 100 đến 105)

INSERT INTO post_media (post_id, media_id, collection, sort_order) VALUES
(10, 100, 'thumbnail', 1),
(10, 101, 'gallery', 1),
(10, 102, 'gallery', 2),
(10, 103, 'gallery', 3),
(10, 104, 'attachment', 1),
(10, 105, 'content', 1);
```
