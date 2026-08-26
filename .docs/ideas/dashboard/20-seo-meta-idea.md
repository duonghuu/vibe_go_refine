# Ý TƯỞNG: Bảng quản lý SEO dùng chung (SEO Meta)

## 1. Thông tin chung

* **Dự án:** TechBite.
* **Tính năng:** Quản lý SEO Meta dùng chung cho nhiều loại entity.
* **Mục đích:** Tách dữ liệu SEO khỏi các bảng nghiệp vụ như `posts`, `products`, `post_categories`, `product_categories` và `pages`.
* **Phạm vi hiện tại:** Mỗi entity chỉ có một bộ SEO duy nhất. Hệ thống chưa hỗ trợ đa ngôn ngữ.

### 1.1. Phân biệt `seo_meta` và `post_meta`

Hai bảng có mục đích khác nhau và không được dùng thay thế cho nhau:

* `seo_meta`: nguồn dữ liệu duy nhất cho các trường SEO.
* `post_meta`: tiếp tục lưu metadata mở rộng dạng key-value, ví dụ `source_url`, `view_count`, `author_notes` hoặc các cấu hình đặc thù khác.

Các key SEO sau không được tạo mới trong `post_meta`:

```text
seo_title
seo_description
canonical_url
og_title
og_description
og_image
twitter_title
twitter_description
twitter_image
robots
schema_json
```

## 2. Định hướng kiến trúc

`seo_meta` sử dụng polymorphic relation để dùng chung cho nhiều entity:

```text
entity_type = post
entity_id   = 10
```

Các entity được hỗ trợ trong giai đoạn đầu:

```text
post
product
post_category
product_category
page
```

`entity_type` không được nhận tự do từ Frontend. Backend phải duy trì whitelist/registry và kiểm tra entity tương ứng tồn tại trước khi tạo hoặc cập nhật SEO.

Vì polymorphic relation không thể tạo foreign key trực tiếp tới nhiều bảng, Service phải chịu trách nhiệm:

* Kiểm tra `entity_type` hợp lệ.
* Kiểm tra `entity_id` tồn tại và chưa bị soft-delete.
* Không cho phép tạo bản ghi SEO mồ côi.
* Xử lý xóa SEO khi entity bị xóa theo lifecycle của entity.

## 3. Thiết kế dữ liệu

```sql
CREATE TABLE seo_meta (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    entity_type VARCHAR(50) NOT NULL,
    entity_id BIGINT UNSIGNED NOT NULL,

    meta_title VARCHAR(255) NULL,
    meta_description VARCHAR(500) NULL,
    meta_keywords VARCHAR(500) NULL,

    canonical_url VARCHAR(500) NULL,

    og_title VARCHAR(255) NULL,
    og_description VARCHAR(500) NULL,
    og_image VARCHAR(500) NULL,
    og_type VARCHAR(100) NULL,

    twitter_title VARCHAR(255) NULL,
    twitter_description VARCHAR(500) NULL,
    twitter_image VARCHAR(500) NULL,
    twitter_card VARCHAR(100) NULL,

    robots VARCHAR(100) NULL,
    schema_json JSON NULL,

    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,

    UNIQUE KEY uq_seo_meta_entity (entity_type, entity_id),
    INDEX idx_seo_meta_entity (entity_type, entity_id)
);
```

### 3.1. Quy tắc dữ liệu

* `entity_type`: chỉ nhận giá trị trong registry của Backend.
* `entity_id`: số nguyên dương, phải tồn tại trong entity tương ứng.
* Mỗi entity chỉ có tối đa một bản ghi `seo_meta`.
* Không cần `locale` ở giai đoạn hiện tại.
* Không nên soft-delete `seo_meta`; khi entity bị xóa, SEO nên được xóa cùng lifecycle hoặc được dọn bởi job cleanup.
* `schema_json` phải là JSON hợp lệ, có giới hạn kích thước và chỉ cho phép các schema type đã được hệ thống chấp nhận.
* `og_image` và `twitter_image` hiện lưu URL để giữ thiết kế đơn giản. Khi Media được chuẩn hóa cho SEO, có thể chuyển sang lưu `media_id` và resolve URL ở Service.

### 3.2. Khả năng mở rộng đa ngôn ngữ

Thiết kế hiện tại áp dụng một bộ SEO cho mỗi entity. Khi triển khai đa ngôn ngữ trong tương lai:

1. Bổ sung cột `locale VARCHAR(10)` với giá trị mặc định phù hợp.
2. Thay unique key bằng `(entity_type, entity_id, locale)`.
3. Bổ sung fallback theo ngôn ngữ mặc định.
4. Không thay đổi API nghiệp vụ hiện tại nếu `locale` được để optional.

## 4. Quy tắc nghiệp vụ

### 4.1. Validation

`meta_title`:

* Không bắt buộc.
* Tối đa 255 ký tự.
* Chuỗi rỗng sau khi trim được xem như không có giá trị.

`meta_description`:

* Không bắt buộc.
* Tối đa 500 ký tự.
* Khuyến nghị khoảng 150–160 ký tự.

`canonical_url`:

* Không bắt buộc.
* Nếu có giá trị thì phải là URL hợp lệ.
* Nếu rỗng, Service build URL từ `base_url` theo environment và slug của entity.
* Không tự động ghi URL fallback trở lại database.

`robots`:

* Không bắt buộc.
* Mặc định là `index,follow`.
* Giai đoạn đầu hỗ trợ:

```text
index,follow
noindex,follow
noindex,nofollow
```

`schema_json`:

* Không bắt buộc.
* Phải parse được như JSON.
* Khi render ra HTML phải serialize an toàn, không cho phép chèn script tùy ý ngoài JSON-LD hợp lệ.

### 4.2. Fallback khi đọc SEO

Fallback được resolve tại Service/SEO Resolver, không lưu ngược vào database:

```text
meta_title          → entity.title
meta_description    → excerpt hoặc content đã strip HTML và cắt ngắn
canonical_url       → base_url + entity.slug
og_title            → resolved meta_title
og_description      → resolved meta_description
og_image            → thumbnail của entity nếu có
twitter_title       → resolved og_title
twitter_description → resolved og_description
twitter_image       → resolved og_image
robots              → index,follow
```

Resolver phải phân biệt giá trị `NULL`, chuỗi rỗng và chuỗi chỉ chứa whitespace. Việc cắt nội dung phải an toàn với Unicode và không làm hỏng HTML/UTF-8.

## 5. 

## 6. API dự kiến

Base route quản trị:

```text
/api/v1/admin/seo-meta
```

### 6.1. Lấy SEO của entity

```http
GET /api/v1/admin/seo-meta/:entityType/:entityId
```

Yêu cầu JWT hợp lệ và role `ADMIN` hoặc `STAFF`.

Response nên phân biệt giá trị lưu trong DB và giá trị sau fallback:

```json
{
  "data": {
    "entity_type": "post",
    "entity_id": 10,
    "meta_title": null,
    "resolved_title": "Tiêu đề bài viết",
    "meta_description": null,
    "resolved_description": "Mô tả được sinh từ nội dung bài viết",
    "canonical_url": null,
    "resolved_canonical_url": "https://example.com/posts/example-post",
    "og_title": null,
    "og_description": null,
    "og_image": null,
    "robots": "index,follow",
    "schema_json": null
  }
}
```

### 6.2. Upsert SEO của entity

```http
PUT /api/v1/admin/seo-meta/:entityType/:entityId
```

Request body:

```json
{
  "meta_title": "Tiêu đề SEO",
  "meta_description": "Mô tả SEO",
  "meta_keywords": "golang, backend",
  "canonical_url": "https://example.com/bai-viet",
  "og_title": "Tiêu đề chia sẻ",
  "og_description": "Mô tả chia sẻ",
  "og_image": "https://example.com/uploads/og-image.jpg",
  "og_type": "article",
  "twitter_title": "Tiêu đề Twitter",
  "twitter_description": "Mô tả Twitter",
  "twitter_image": "https://example.com/uploads/twitter-image.jpg",
  "twitter_card": "summary_large_image",
  "robots": "index,follow",
  "schema_json": {
    "@context": "https://schema.org",
    "@type": "Article"
  }
}
```

API phải tự kiểm tra `entityType` và `entityId`, không nhận `entity_type`/`entity_id` từ request body để tránh sai lệch giữa URI và payload.

### 6.3. Chuẩn lỗi

```json
{
  "error": "validation_error",
  "message": "Dữ liệu SEO không hợp lệ",
  "fields": {
    "canonical_url": "URL không hợp lệ"
  }
}
```

Các status chính:

* `400 Bad Request`: URI hoặc payload không hợp lệ.
* `401 Unauthorized`: JWT thiếu, sai, hết hạn hoặc bị blacklist.
* `403 Forbidden`: không đủ quyền quản trị SEO.
* `404 Not Found`: entity không tồn tại hoặc đã bị soft-delete.
* `422 Unprocessable Entity`: dữ liệu SEO không đạt validation.
* `500 Internal Server Error`: lỗi hệ thống.

## 7. Luồng xử lý

### 7.1. Tạo entity có SEO

1. Admin nhập dữ liệu entity.
2. Backend validate và lưu entity.
3. Upload media nếu có.
4. Gọi upsert `seo_meta` với các trường SEO người dùng nhập.
5. Khi public site đọc dữ liệu, SEO Resolver áp dụng fallback.

SEO Meta là module riêng, không cập nhật ngầm trong transaction Post/Product nếu không có yêu cầu nghiệp vụ rõ ràng.

### 7.2. Chỉnh sửa entity

1. Query entity.
2. Query `seo_meta` theo `entity_type` và `entity_id`.
3. Trả cả dữ liệu gốc và dữ liệu đã resolve cho Frontend.
4. Admin cập nhật SEO qua API SEO Meta.
5. Invalidate cache sau khi upsert thành công.

### 7.3. Xóa entity

Khi entity bị xóa, Service phải xóa bản ghi `seo_meta` tương ứng hoặc đưa vào quy trình cleanup rõ ràng. Không để tồn tại SEO record mồ côi.

## 8. Cache

Cache key thống nhất:

```text
seo:{entity_type}:{entity_id}
```

* Public read sử dụng cache-aside nếu cần tối ưu.
* Admin read có thể đọc trực tiếp database hoặc dùng chung cache sau khi invalidate.
* Sau `PUT`, bắt buộc xóa cache SEO tương ứng.
* Lỗi Redis không được làm mất dữ liệu đã ghi thành công vào MySQL; phải ghi log để theo dõi lỗi invalidate.

## 9. Phân quyền và bảo mật

* GET/PUT admin đều đi qua JWT Auth Middleware.
* Quyền thao tác tối thiểu: `ADMIN`, `STAFF`.
* Không tin tưởng `userId` hoặc `role` từ body/query; lấy từ JWT Context.
* Không trả token, thông tin nhạy cảm hoặc stack trace trong response.
* `schema_json`, image URL và canonical URL phải được validate trước khi lưu/render.

## 10. Tiêu chí nghiệm thu

* `post_meta` không còn là nguồn đọc SEO.
* Một entity chỉ có tối đa một `seo_meta` record.
* `entity_type` ngoài whitelist bị từ chối.
* Entity không tồn tại không thể tạo SEO.
* Upsert có thể chạy lặp lại mà không tạo dữ liệu trùng.
* Fallback hoạt động độc lập cho từng trường SEO.
* Giá trị fallback không bị ghi đè vào dữ liệu gốc.
* Cache được invalidate sau khi cập nhật SEO.
* Xóa entity không để lại SEO record mồ côi.
* Có unit test cho SEO Resolver và integration test cho API/unique constraint.

## 11. Khả năng mở rộng sau này

Các hướng mở rộng không thuộc phạm vi hiện tại:

* SEO theo `locale`.
* Lưu `media_id` thay cho image URL.
* Versioning và lịch sử chỉnh sửa SEO.
* Preview SERP/Open Graph trong Dashboard.
* Sitemap và kiểm tra broken canonical URL.