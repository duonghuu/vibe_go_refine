# Ý TƯỞNG: Quản lý Section của từng Page (Semi-Structured CMS)

## 1. Thông tin chung

* **Dự án:** TechBite.
* **Tính năng:** Quản lý các Section thuộc từng Page.
* **Giải pháp:** Semi-Structured CMS (CMS bán cấu trúc).
* **Mục đích:** Cho phép một Page tổ chức nội dung thành nhiều Section có thứ tự và metadata chung; các dữ liệu liên kết như Category, Post và Media được quản lý độc lập, không nhúng thành JSON trong Page hoặc Section.
* **Phạm vi UI:** Khu vực quản lý Section được nhúng trong màn hình chỉnh sửa Page. Page phải được tạo trước khi có thể thêm Section.

---

## 2. Bài toán & nguyên tắc thiết kế

* Hệ thống có nhiều Page; mỗi Page có thể không có Section hoặc có nhiều Section với cấu trúc phức tạp.
* Page biết chính xác các Section của mình thông qua quan hệ `pages 1-N page_sections`.
* Section chỉ giữ metadata và các Media đại diện dùng chung cho phần hiển thị.
* Category, Post và Media trong từng collection không được sao chép vào Section mà được tham chiếu qua `page_section_items`.
* Cấu trúc bán cố định giúp Backend vẫn validate được dữ liệu, đồng thời cho phép một Section có nhiều collection như `categories`, `featured_posts`, `gallery` mà không phải thêm cột mới.
* Trường `content` hiện tại của Page vẫn được giữ. Page có thể dùng content, Section hoặc kết hợp cả hai; việc render do Webview quyết định theo contract public.

```text
Page
  └── PageSection (hero, feature, product, ...)
        ├── backgroundMedia -> Media
        ├── featureMedia    -> Media
        └── PageSectionItem[]
              ├── CATEGORY -> Category
              ├── POST     -> Post
              └── MEDIA    -> Media
```

---

## 3. Mô hình dữ liệu

### 3.1. Bảng `page_sections`

```text
page_sections
-------------
id
page_id

key                     // hero, feature, product...
name                    // Tên quản trị, ví dụ: Hero Banner

title
description

background_color
background_media_id
feature_media_id

sort_order
status

created_at
updated_at
```

Quan hệ và ràng buộc:

* `page_id` là foreign key tới `pages.id`; một Page có nhiều Section.
* `background_media_id` và `feature_media_id` là foreign key nullable tới `media.id`.
* Khi xóa vật lý Page, Section và Item con bị xóa cascade. Page soft-delete không làm mất cấu hình Section.
* Khi Media nền hoặc Media feature bị xóa hợp lệ, foreign key tương ứng được đặt về `NULL`.
* `key` là duy nhất trong phạm vi một Page bằng unique key `(page_id, key)`.
* Có index `(page_id, status, sort_order)` để tải và sắp xếp Section của Page.
* `status` chỉ nhận `ACTIVE` hoặc `INACTIVE`; mặc định `ACTIVE`.

### 3.2. Bảng `page_section_items`

```text
page_section_items
------------------
id

section_id

item_type               // CATEGORY | POST | MEDIA
item_id

collection              // categories | featured_posts | gallery

sort_order
```

Quan hệ và ràng buộc:

* `section_id` là foreign key tới `page_sections.id`; xóa Section sẽ xóa cascade toàn bộ Item thuộc Section.
* `item_type` xác định loại thực thể mà `item_id` tham chiếu: `CATEGORY` tới `categories`, `POST` tới `posts`, `MEDIA` tới `media`.
* Do `item_id` là liên kết đa hình, database không tạo được một foreign key chung. Use Case bắt buộc kiểm tra thực thể tồn tại, chưa bị soft-delete và người dùng có quyền liên kết trước khi lưu.
* `collection` là semantic key mở rộng, viết dạng `snake_case`; không đóng cứng thành enum. Ví dụ: `categories`, `featured_posts`, `latest_posts`, `gallery`.
* Không cho phép trùng một Item trong cùng collection, bằng unique key `(section_id, collection, item_type, item_id)`.
* Một Item vẫn có thể xuất hiện ở hai collection khác nhau trong cùng Section.
* Có index `(section_id, collection, sort_order)` để tải đúng thứ tự từng collection.

---

## 4. Metadata của Section

| Trường | Ý nghĩa |
| :--- | :--- |
| `key` | Khóa kỹ thuật ổn định để Webview chọn component render, ví dụ `hero`, `featured_news`. |
| `name` | Tên hiển thị trong Dashboard để biên tập viên nhận biết Section. |
| `title` | Tiêu đề nội dung của Section, có thể để trống. |
| `description` | Mô tả nội dung, có thể để trống. |
| `background_color` | Màu nền cấu hình cho Section, có thể để trống để Webview dùng theme mặc định. |
| `background_media_id` | Ảnh nền tùy chọn, tham chiếu Media. |
| `feature_media_id` | Ảnh nổi bật tùy chọn, tham chiếu Media. |
| `sort_order` | Vị trí Section trong Page, bắt đầu từ `0`. |
| `status` | `ACTIVE` được trả ra public; `INACTIVE` chỉ tồn tại trong Dashboard. |

---

## 5. Trải nghiệm quản trị

### 5.1. Danh sách Section trong Page

* Hiển thị các Section dạng card/accordion theo `sort_order`.
* Mỗi Section hiển thị tên quản trị, key, trạng thái, tiêu đề và số lượng Item theo collection.
* Hỗ trợ thêm, chỉnh sửa, xóa, bật/tắt và kéo thả để sắp xếp Section.
* Xóa Section phải có hộp thoại xác nhận vì toàn bộ liên kết Item của Section cũng bị xóa.
* Nếu Page chưa được tạo, khu vực Section bị vô hiệu hóa và hướng dẫn người dùng lưu Page trước.

### 5.2. Form Section

* Thông tin chung: `key`, `name`, `title`, `description`, `status`.
* Giao diện: `background_color`, chọn/xóa Background Media, chọn/xóa Feature Media.
* Collection: hiển thị từng nhóm Item, cho phép thêm nhiều Item, gỡ Item và kéo thả thay đổi thứ tự.
* Việc gỡ Media khỏi Section chỉ xóa liên kết, không xóa file Media gốc.

### 5.3. Luồng chọn dữ liệu liên kết

```text
+ Add Category / Post / Media
        ↓
Mở popup danh sách tương ứng
        ↓
Tìm kiếm, lọc và chọn nhiều bản ghi
        ↓
Xác nhận lựa chọn
        ↓
Lưu item_id vào page_section_items
        ↓
Refresh collection và giữ đúng sort_order
```

* Popup phải giữ lại các Item đã chọn khi phân trang hoặc tìm kiếm.
* Item đã có trong collection được đánh dấu selected và không được thêm trùng.
* Ô tìm kiếm/autocomplete bắt buộc debounce `300ms` bằng hook/thư viện chuẩn của dự án.
* UI phải xử lý đủ Loading, Error, Empty và Success cho danh sách Section và từng collection.

---

## 6. Quy tắc nghiệp vụ & validation

* Chỉ được thao tác Section khi Page tồn tại và chưa bị soft-delete.
* `key` bắt buộc, tối đa 100 ký tự, chỉ gồm chữ thường, số và dấu gạch dưới; không cho đổi `key` ngoài ý muốn sau khi Webview đã phụ thuộc vào key đó.
* `name` bắt buộc, tối đa 255 ký tự. `title` tối đa 255 ký tự.
* `sort_order` là số nguyên lớn hơn hoặc bằng `0`.
* Khi sắp xếp lại, Backend đồng bộ thứ tự trong transaction và chuẩn hóa thành dãy liên tục từ `0`.
* Background Media và Feature Media phải tồn tại, chưa bị xóa và có MIME type ảnh.
* `item_type` chỉ nhận `CATEGORY`, `POST`, `MEDIA`; `item_id` phải lớn hơn `0` và đúng loại thực thể.
* `collection` bắt buộc, tối đa 100 ký tự, định dạng `snake_case`.
* Một request đồng bộ collection phải dùng cùng `item_type`; không trộn Category, Post và Media trong một collection.
* Thao tác thêm/xóa/sắp xếp Item được thực hiện trong transaction để tránh collection lưu dở dang.
* Admin được quản lý mọi Section. Staff chỉ được thao tác theo quyền Page và quy tắc sở hữu Media hiện có; không tin tưởng `userId` gửi từ client.
* Mọi thay đổi Section, Item, thứ tự hoặc trạng thái phải xóa cache public của Page tương ứng.

---

## 7. API quản trị dự kiến

```text
GET    /api/v1/admin/pages/:page_id/sections
POST   /api/v1/admin/pages/:page_id/sections
GET    /api/v1/admin/pages/:page_id/sections/:section_id
PUT    /api/v1/admin/pages/:page_id/sections/:section_id
DELETE /api/v1/admin/pages/:page_id/sections/:section_id
PUT    /api/v1/admin/pages/:page_id/section-order

GET    /api/v1/admin/pages/:page_id/sections/:section_id/items?collection={collection}
PUT    /api/v1/admin/pages/:page_id/sections/:section_id/items/:collection
```

Các popup lựa chọn tái sử dụng API danh sách hiện có của Category, Post và Media. Không tạo bản sao dữ liệu chỉ để phục vụ Section.

### 7.1. Tạo Section

```json
{
  "key": "featured_news",
  "name": "Tin tức nổi bật",
  "title": "Nổi bật hôm nay",
  "description": "Các bài viết được biên tập viên lựa chọn",
  "backgroundColor": null,
  "backgroundMediaId": null,
  "featureMediaId": 120,
  "sortOrder": 1,
  "status": "ACTIVE"
}
```

### 7.2. Đồng bộ thứ tự Section

```json
{
  "sections": [
    { "id": 10, "sortOrder": 0 },
    { "id": 12, "sortOrder": 1 },
    { "id": 11, "sortOrder": 2 }
  ]
}
```

### 7.3. Đồng bộ một collection

```json
{
  "itemType": "POST",
  "items": [
    { "itemId": 501, "sortOrder": 0 },
    { "itemId": 498, "sortOrder": 1 },
    { "itemId": 477, "sortOrder": 2 }
  ]
}
```

Request đồng bộ là trạng thái cuối cùng của collection: Item không còn trong payload sẽ bị gỡ liên kết, Item mới sẽ được thêm và tất cả `sort_order` được cập nhật nguyên tử.

---

## 8. Public Page API

API public hiện có được mở rộng để trả thêm `sections`:

```text
GET /api/v1/pages/:slug
```

```json
{
  "data": {
    "id": 1,
    "title": "Trang chủ",
    "slug": "trang-chu",
    "content": "",
    "sections": [
      {
        "id": 10,
        "key": "featured_news",
        "title": "Nổi bật hôm nay",
        "description": "Các bài viết được lựa chọn",
        "backgroundColor": null,
        "backgroundMedia": null,
        "featureMedia": { "id": 120, "originalUrl": "/uploads/feature.webp" },
        "sortOrder": 0,
        "collections": {
          "featured_posts": [
            { "itemType": "POST", "itemId": 501, "sortOrder": 0, "data": {} }
          ]
        }
      }
    ]
  }
}
```

* Chỉ trả Section `ACTIVE`, sắp xếp theo `sort_order` tăng dần.
* Item trong mỗi collection được hydrate thành dữ liệu public tối thiểu của đúng thực thể và sắp xếp theo `sort_order`.
* Item đã bị xóa hoặc không còn hợp lệ không được trả public nhưng liên kết không tự động chuyển sang một thực thể khác.
* Response không trả metadata nội bộ hoặc dữ liệu nhạy cảm.

---

## 9. Xử lý lỗi

* `400 Bad Request`: payload, key, collection, item type hoặc thứ tự không hợp lệ.
* `401 Unauthorized`: chưa xác thực.
* `403 Forbidden`: không có quyền quản lý Page hoặc Media liên quan.
* `404 Not Found`: Page, Section hoặc Item nguồn không tồn tại/đã bị xóa.
* `409 Conflict`: key Section hoặc Item trong collection bị trùng.
* `500 Internal Server Error`: lỗi database/hạ tầng; không trả chi tiết nội bộ cho client.

---

## 10. Ngoài phạm vi phiên bản đầu

* Không xây dựng page builder tự do hoặc lưu layout tùy ý bằng JSON.
* Không hỗ trợ Section lồng Section.
* Không hỗ trợ lịch sử phiên bản, duyệt nội dung, lịch xuất bản, đa ngôn ngữ hoặc A/B testing.
* Không tự động xóa Category, Post hay Media gốc khi xóa Item hoặc Section.
* Không cho phép plugin tự định nghĩa `item_type` ngoài ba loại đã thống nhất.

---

## 11. Tiêu chí nghiệm thu

* Một Page có thể tạo, sửa, xóa, bật/tắt và sắp xếp nhiều Section.
* Mỗi Section lưu được metadata chung, Background Media và Feature Media.
* Mỗi Section quản lý được nhiều collection Category, Post hoặc Media, chọn nhiều và sắp xếp được (optional).
* Không thể liên kết bản ghi không tồn tại, sai loại, không có quyền hoặc thêm trùng trong cùng collection.
* Xóa Section xóa toàn bộ liên kết Item nhưng không xóa dữ liệu Category, Post, Media gốc.
* API public chỉ trả Page công khai, Section đang hoạt động và Item hợp lệ theo đúng thứ tự.
* Mọi thay đổi Section/Item làm mới cache public của Page.
* Dashboard xử lý đầy đủ Loading, Error, Empty, Success và debounce tìm kiếm `300ms`.
