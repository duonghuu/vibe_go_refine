# Ý TƯỞNG: Quản lý hình ảnh bài viết (Post Image Management)

## 1. Thông tin chung (Meta Info)

* **Dự án:** TechBite.
* **Tính năng:** Quản lý ảnh đại diện và thư viện ảnh của bài viết.
* **Mục đích:** Tích hợp thao tác upload, chọn, sắp xếp và gỡ ảnh trực tiếp vào màn hình Tạo/Sửa bài viết; dùng bảng trung gian `post_media` làm nguồn dữ liệu liên kết chính giữa `posts` và `media`.
* **Phạm vi giai đoạn đầu:** Hỗ trợ `thumbnail` và `gallery`. `content` và `attachment` được giữ trong mô hình dữ liệu để phát triển cùng Rich Text Editor và upload tài liệu ở giai đoạn sau.

---

## 2. Định hướng kiến trúc

### 2.1. Nguồn dữ liệu chuẩn

Không lưu URL ảnh đại diện trực tiếp vào bảng `posts`. Mọi ảnh của bài viết được quản lý qua quan hệ:

```text
Upload ảnh
  -> media (status = temporary)
  -> lưu bài viết
  -> post_media (post_id, media_id, collection, sort_order)
  -> media (status = attached)
```

`post_media` là nguồn dữ liệu chuẩn để truy vấn ảnh của một bài viết. Bảng `media` chỉ lưu metadata file, URL và trạng thái vòng đời; không thêm cột `post_id` vào `media`.

### 2.2. Tham khảo luồng Product

Màn hình Product hiện có luồng chọn file, upload tức thời, preview và xóa preview. Luồng UI này có thể tái sử dụng cho Post.

Tuy nhiên, Post không chỉ lưu URL như trường `products.image_url`. Frontend phải giữ đầy đủ metadata upload trả về, tối thiểu là:

```ts
type UploadedMedia = {
  id: number;
  originalUrl: string;
  thumbnailUrl?: string;
  mediumUrl?: string;
  fileName: string;
  status: "temporary" | "attached";
};
```

Giữ `media.id` là điều kiện bắt buộc để tạo liên kết chính xác trong `post_media`, đồng bộ thứ tự ảnh và dọn file tạm khi người dùng hủy thao tác.

---

## 3. Collection và quy tắc nghiệp vụ

| Collection | Mục đích | Quy tắc |
| --- | --- | --- |
| `thumbnail` | Ảnh cover, hiển thị tại danh sách bài viết và đầu trang chi tiết | Mỗi Post có tối đa một ảnh; có thể để trống nếu nghiệp vụ cho phép. |
| `gallery` | Album/slider ảnh tại trang chi tiết | Có nhiều ảnh; thứ tự hiển thị theo `sort_order`. |
| `content` | Ảnh được chèn từ Rich Text Editor | Không hiển thị trong card upload thông thường; sẽ liên kết khi editor hỗ trợ chèn ảnh theo `media.id`. |
| `attachment` | File người đọc tải về | Chờ media service hỗ trợ PDF/DOC và UI đính kèm tài liệu. |

Các quy tắc bắt buộc:

* Một `media_id` không được liên kết trùng trong cùng `post_id` và `collection`.
* `thumbnail` phải được thay thế theo cơ chế sync, không thêm dồn nhiều record.
* `gallery` được đồng bộ toàn bộ danh sách, gồm cả thứ tự mới sau khi kéo-thả.
* Chỉ cho phép liên kết Media còn tồn tại, chưa bị xóa mềm, và thuộc quyền quản lý của người dùng hiện tại hoặc người dùng có vai trò Admin.
* Khi một Media không còn bất kỳ liên kết nào với Post hay đối tượng khác, nó có thể trở lại trạng thái có thể dọn dẹp; cron sẽ dọn những file `temporary` quá hạn.

---

## 4. Trải nghiệm giao diện (UI/UX)

### 4.1. Vị trí trên form Post

Thêm Card **Hình ảnh bài viết** ở cột phụ, đặt sau Card **Phân loại** trong màn hình Tạo mới và Chỉnh sửa bài viết.

### 4.2. Ảnh đại diện

* Chọn hoặc kéo-thả một file ảnh.
* Upload xảy ra ngay khi người dùng chọn file.
* Sau upload, hiển thị preview, tên file và nút thay thế/xóa.
* Khi chọn ảnh mới, state chỉ giữ ảnh mới làm `thumbnail`; ảnh cũ chỉ bị gỡ liên kết sau khi lưu thành công.
* Nếu xóa ảnh vừa upload nhưng chưa lưu bài viết, Frontend gọi `DELETE /api/v1/media/:id` để dọn ngay file tạm.

### 4.3. Thư viện ảnh

* Cho phép chọn nhiều file ảnh hoặc kéo-thả nhiều file vào vùng upload.
* Hiển thị dạng lưới thumbnail; mỗi item có preview, trạng thái upload và nút xóa.
* Hỗ trợ kéo-thả để sắp xếp lại; vị trí trong mảng tương ứng với `sort_order`.
* Các ảnh chỉ được sync sau khi người dùng bấm Lưu; khi upload/lưu lỗi phải hiển thị Snackbar có thông điệp cụ thể.

### 4.4. Trạng thái và an toàn dữ liệu

* Trong lúc upload hoặc lưu, disable nút Lưu để tránh gửi danh sách Media chưa hoàn thiện.
* Cảnh báo rời trang khi form hoặc danh sách ảnh có thay đổi chưa lưu.
* Trang Edit tải ảnh bằng API Post Media trước khi render preview đầy đủ.
* Nếu tạo Post thành công nhưng sync ảnh thất bại, phải báo lỗi rõ ràng và điều hướng tới Edit Post để người dùng hoàn tất; phương án ưu tiên vẫn là xử lý transaction ở Backend để tránh trạng thái này.

---

## 5. Luồng xử lý

### 5.1. Tạo bài viết

1. Người dùng upload ảnh; mỗi ảnh tạo một record `media` có `status = temporary`.
2. Frontend giữ danh sách `UploadedMedia` trong form state.
3. Người dùng bấm Lưu.
4. Backend tạo `posts`, xác thực các Media, tạo record `post_media` và chuyển Media đang được dùng thành `attached` trong cùng transaction.
5. Trả về Post cùng dữ liệu Media đã liên kết; Frontend điều hướng về danh sách bài viết.

Payload đề xuất:

```json
{
  "typeCode": "NEWS",
  "title": "Xu hướng công nghệ mới",
  "slug": "xu-huong-cong-nghe-moi",
  "content": "...",
  "categoryId": 3,
  "media": {
    "thumbnailId": 101,
    "gallery": [
      { "id": 102, "sortOrder": 1 },
      { "id": 103, "sortOrder": 2 }
    ]
  }
}
```

### 5.2. Chỉnh sửa bài viết

1. Frontend tải Post và `GET /api/v1/admin/posts/:post_id/media`.
2. Render Media theo `collection` và `sort_order`.
3. Người dùng thay ảnh đại diện, thêm/xóa/sắp xếp gallery.
4. Khi lưu, Backend cập nhật thông tin Post và đồng bộ Media theo collection trong một transaction.
5. Cache chi tiết/list Post và cache Media của Post được invalidated sau khi commit thành công.

---

## 6. API đề xuất

Module `post_media` hiện đã có các sub-resource phù hợp cho trang Edit:

```text
GET    /api/v1/admin/posts/:post_id/media?collection=gallery
POST   /api/v1/admin/posts/:post_id/media
PUT    /api/v1/admin/posts/:post_id/media/:collection
DELETE /api/v1/admin/posts/:post_id/media/:media_id?collection=gallery
```

API sync là API ưu tiên cho thao tác từ UI vì idempotent và đồng bộ đúng trạng thái cuối cùng của giao diện:

```json
PUT /api/v1/admin/posts/42/media/thumbnail
{
  "media": [{ "id": 101, "sortOrder": 0 }]
}
```

```json
PUT /api/v1/admin/posts/42/media/gallery
{
  "media": [
    { "id": 102, "sortOrder": 1 },
    { "id": 103, "sortOrder": 2 }
  ]
}
```

Đối với Create Post, ưu tiên nhận object `media` ngay trong `POST /posts` để toàn bộ tạo Post và liên kết ảnh diễn ra trong một transaction. Các endpoint sub-resource vẫn hữu ích cho Edit, reorder hoặc trường hợp lưu từng collection độc lập.

---

## 7. Ràng buộc Database và vòng đời dữ liệu

### 7.1. Ràng buộc cần có cho `post_media`

* Bật Foreign Key `post_id -> posts(id)` và `media_id -> media(id)`; khi hard-delete Media cần cascade các liên kết tương ứng.
* Thêm unique index `(post_id, media_id, collection)` để tránh liên kết trùng.
* Giữ index `(post_id, collection, sort_order)` để tối ưu truy vấn hiển thị gallery.
* `thumbnail` tối đa một record/Post được kiểm soát trong Service/transaction. MySQL không hỗ trợ partial unique index trực tiếp cho điều kiện `collection = 'thumbnail'`.

### 7.2. Soft delete và cascade

`posts` hiện dùng soft delete. Vì vậy FK `ON DELETE CASCADE` không tự chạy khi GORM chỉ cập nhật `deleted_at`.

Khi xóa mềm một Post, Post Service phải xóa mềm toàn bộ record `post_media` của Post trong cùng transaction. Khi hard-delete một Media, các liên kết `post_media` có thể được cascade ở database.

### 7.3. Trạng thái Media

Không gắn `post_id` trực tiếp vào `media`, vì một Media có thể được tái sử dụng ở nhiều Post hoặc collection. `status = attached` có nghĩa Media đang có ít nhất một liên kết hợp lệ; nếu hệ thống về sau cho tái sử dụng giữa Product và Post, việc xác định trạng thái phải xét toàn bộ các liên kết thay vì chỉ dựa vào `media.product_id`.

---

## 8. Validation, bảo mật và lỗi

* Chỉ chấp nhận các collection hợp lệ: `thumbnail`, `gallery`, `attachment`, `content`.
* Sync thumbnail chỉ nhận tối đa một phần tử; gallery không nhận ID trùng.
* Kiểm tra Post tồn tại trước mọi thao tác Media.
* Kiểm tra Media tồn tại, không bị xóa, đúng định dạng cho collection và thuộc quyền của người dùng.
* Lỗi validation trả `400`, thiếu Post/Media trả `404`, không đủ quyền trả `403`; không trả `500` cho lỗi nghiệp vụ.
* API upload và API Post Media đều phải dùng Auth Middleware; user ID lấy từ JWT context, không lấy từ request body.
* Chuẩn hóa payload JSON về `camelCase` để đồng nhất với Post và Product; tránh việc UI mới phải dùng lẫn `sortOrder` và `sort_order`.

---

## 9. Tiêu chí hoàn thành (Acceptance Criteria)

* Admin/Staff tạo và sửa Post có thể upload, preview, thay thế, xóa và sắp xếp ảnh.
* Một Post có tối đa một thumbnail và nhiều gallery ảnh có thứ tự ổn định.
* Reload trang Edit hiển thị đúng danh sách Media đã lưu.
* Xóa ảnh khỏi gallery chỉ gỡ liên kết, không xóa file đang được dùng ở Post/đối tượng khác.
* Media mới upload nhưng không được dùng được dọn ngay khi người dùng xóa hoặc được cron dọn sau thời hạn.
* Không có record `post_media` mồ côi, record trùng, hoặc ảnh thuộc người dùng khác bị gắn trái phép.

---

## 10. Hiệu năng khi hiển thị danh sách Post

### 10.1. Đánh đổi khi không lưu URL trực tiếp trong `posts`

Không lưu `thumbnail_url` trực tiếp trong bảng `posts` giúp quan hệ dữ liệu nhất quán và hỗ trợ nhiều collection, nhưng API danh sách phải JOIN thêm `post_media` và `media` để lấy ảnh đại diện.

Chi phí JOIN này thường không đáng kể ở quy mô hiện tại nếu truy vấn được tối ưu. Vấn đề cần tránh là gọi API Media riêng cho từng Post, tạo ra mô hình N+1 query.

### 10.2. Truy vấn đề xuất cho Post List

Repository lấy danh sách Post chỉ JOIN collection `thumbnail`, không tải toàn bộ gallery:

```sql
SELECT
  p.id,
  p.title,
  p.slug,
  p.type_code,
  p.created_at,
  m.original_url AS thumbnail_url
FROM posts p
LEFT JOIN post_media pm
  ON pm.post_id = p.id
  AND pm.collection = 'thumbnail'
  AND pm.deleted_at IS NULL
LEFT JOIN media m
  ON m.id = pm.media_id
  AND m.deleted_at IS NULL
WHERE p.deleted_at IS NULL
ORDER BY p.created_at DESC
LIMIT 20 OFFSET 0;
```

API list chỉ trả các trường cần thiết cho bảng:

```json
{
  "id": 42,
  "title": "Xu hướng công nghệ mới",
  "slug": "xu-huong-cong-nghe-moi",
  "thumbnailUrl": "/uploads/cover.jpg",
  "createdAt": "2026-08-25T10:00:00Z"
}
```

### 10.3. Index và nguyên tắc tải dữ liệu

* Bổ sung index `(post_id, collection, sort_order)` cho `post_media`; nếu dữ liệu soft-delete lớn, cân nhắc index bao gồm `deleted_at` theo chiến lược MySQL đang sử dụng.
* Không gọi `GET /posts/:post_id/media` riêng cho từng dòng trong list.
* Không `Preload` toàn bộ gallery hoặc attachment trong API list.
* Chỉ tải gallery ở trang chi tiết, trang Edit hoặc khi người dùng mở khu vực ảnh.
* Có thể cache kết quả Post List và thumbnail bằng Redis khi lưu lượng đọc tăng; mọi thao tác sync Media phải invalidate cache liên quan.

### 10.4. Quyết định thiết kế

Giai đoạn đầu **không thêm `thumbnail_url` vào `posts`**. Dùng `JOIN` thumbnail trong một truy vấn list duy nhất là đủ hiệu năng và tránh phải đồng bộ URL ở hai nơi.

Chỉ xem xét denormalize `thumbnail_url` vào `posts` khi profiling thực tế chứng minh JOIN là bottleneck. Nếu áp dụng sau này, Service bắt buộc cập nhật URL đồng thời khi thumbnail được thay đổi để tránh dữ liệu không nhất quán.
