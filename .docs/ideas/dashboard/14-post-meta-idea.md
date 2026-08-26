# Ý TƯỞNG: Thông tin mở rộng Bài viết (Post Meta)

## 1. Thông tin chung (Meta Info)

* **Dự án:** TechBite.
* **Tính năng:** Thông tin mở rộng bài viết (Post Meta).
* **Mục đích:** Quản lý thông tin metadata linh hoạt cho bài viết dưới dạng Key-Value (ví dụ: cấu hình SEO title/description, custom fields, rating, luợt xem, các cấu hình tuỳ biến khác không nằm trong cấu trúc bảng `posts` cố định).

---

# 2. Đặc tả dữ liệu (Data Specs)

Bảng `post_meta`:
* **`id`**: Khóa chính tự tăng.
* **`post_id`**: ID của bài viết (Khóa ngoại liên kết tới bảng `posts`).
* **`key`**: Tên của trường thông tin mở rộng (Key - Dạng chuỗi).
* **`value`**: Giá trị của trường thông tin mở rộng (Value - Dạng text/chuỗi, có thể lưu trữ JSON nếu cần).

Các ví dụ về `key` phổ biến:
* `seo_title`: Tiêu đề SEO của bài viết.
* `seo_description`: Mô tả SEO.
* `source_url`: Nguồn bài viết (nếu đi copy).
* `view_count`: Lượt xem bài viết (có thể lưu riêng để tránh lock bảng bài viết khi cập nhật).

---

# 3. Kiểm tra dữ liệu (Validation & Error Handling)

### post_id
* Bắt buộc (Required).
* Phải tồn tại trong hệ thống (tồn tại trong bảng `posts`).

### key
* Bắt buộc (Required).
* Phải là chuỗi (String).
* Không được để trống.
* Khuyến nghị độ dài không vượt quá 255 ký tự.
* Có thể kết hợp `post_id` và `key` thành một `Unique Index` (Mỗi bài viết chỉ có một giá trị duy nhất cho một khóa cụ thể).

### value
* Không bắt buộc, nhưng tuỳ vào nghiệp vụ của từng `key`.
* Dạng chuỗi (Text).

---

# 4. Quy tắc nghiệp vụ (Business Rules)

* **Xóa bài viết (Cascade Delete):** Khi một bài viết (`post_id`) bị xóa, toàn bộ meta tương ứng trong bảng `post_meta` cũng phải bị xóa.
* **Upsert (Cập nhật hoặc Thêm mới):** Thường khi lưu bài viết, Frontend sẽ gửi lên một danh sách Key-Value. Backend cần xử lý cơ chế Upsert: Nếu `key` chưa tồn tại đối với `post_id` đó thì Insert, nếu đã tồn tại thì Update `value`.
* **Cơ chế Sync (Đồng bộ):** Nếu Frontend gửi danh sách meta ít hơn danh sách meta hiện có, các `key` không có mặt trong danh sách gửi lên (nhưng có trong DB) có thể bị xóa (tuỳ thuộc vào thiết kế nghiệp vụ xem đây là Update toàn bộ hay Update một phần).
* **Định dạng dữ liệu Value:** Có thể quy ước sử dụng JSON cho cột `value` nếu cần lưu trữ phức tạp (mảng, object).

---

# 5. Đặc tả API (API Specs)

*API này có thể đứng độc lập hoặc được gộp chung vào payload khi tạo/cập nhật Post.*

## 5.1. Lấy danh sách Meta của một bài viết

```
GET /posts/{post_id}/meta
```

**Response (Thành công 200):**

```json
{
  "data": [
    {
      "id": 1,
      "post_id": 1,
      "key": "seo_title",
      "value": "Tiêu đề SEO siêu chuẩn"
    },
    {
      "id": 2,
      "post_id": 1,
      "key": "seo_description",
      "value": "Mô tả hấp dẫn về bài viết để hiển thị trên kết quả tìm kiếm Google."
    }
  ]
}
```

## 5.2. Cập nhật/Thêm mới danh sách Meta cho bài viết (Sync/Upsert)

```
PUT /posts/{post_id}/meta
```

**Request Body:**

```json
{
  "meta": [
    {
      "key": "seo_title",
      "value": "Cập nhật tiêu đề SEO"
    },
    {
      "key": "source_url",
      "value": "https://example.com/source"
    }
  ]
}
```

## 5.3. Xóa một Meta cụ thể khỏi bài viết

```
DELETE /posts/{post_id}/meta/{key}
```
*(Xóa theo `key` sẽ dễ dàng hơn so với xóa theo `id` vì Frontend thường kiểm soát thông qua tên `key`).*

---

# 6. Dữ liệu mẫu (Sample Data)

```sql
-- Bảng: post_meta
-- (Giả định post có id = 10)

INSERT INTO post_meta (post_id, key, value) VALUES
(10, 'seo_title', 'Hướng dẫn học Golang cho người mới bắt đầu'),
(10, 'seo_description', 'Bài viết cung cấp những kiến thức cơ bản và lộ trình học Golang hiệu quả nhất năm nay.'),
(10, 'view_count', '1250'),
(10, 'author_notes', 'Cần cập nhật lại phần Interface vào tuần sau');
```
