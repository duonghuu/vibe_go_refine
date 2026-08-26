# BACKEND PLAN: Quản lý hình ảnh bài viết (Post Image Management)

## 1. Mục tiêu và phạm vi

Tổ chức việc upload, liên kết, sắp xếp và gỡ ảnh của bài viết thông qua bảng trung gian `post_media`. Giai đoạn đầu hỗ trợ `thumbnail` và `gallery`; `content` và `attachment` được giữ để mở rộng sau.

`thumbnail` tối đa một ảnh cho mỗi Post. `gallery` có nhiều ảnh và hiển thị theo `sort_order`. Không thêm `post_id` vào bảng `media`; một Media có thể được dùng lại ở nhiều Post hoặc collection.

Stack bắt buộc: Go, Gin, GORM, MySQL, Redis và Google Wire; triển khai theo DDD/Clean Architecture hiện có.

## 2. Đánh giá hiện trạng và phần cần bổ sung

Các module hiện có cần tái sử dụng:

- `domain/media`, `usecases/media`, `MediaController`: upload ảnh, kiểm tra Magic Bytes và lưu `temporary`.
- `domain/post/entity/Post` và `PostMedia`.
- `postmedia` Service/Repository/Controller: đã có GET/POST/PUT-sync/DELETE.
- `AuthMiddleware`/`RoleMiddleware`: xác minh JWT, kiểm tra blacklist JTI, set `userId`/`role` vào Gin Context.
- Redis client và Google Wire trong `internal/di`.

Các điểm cần sửa:

1. Foreign key của `post_media` đang bị comment trong migration.
2. Service chưa xác nhận Post tồn tại, ownership Media, collection thumbnail duy nhất hoặc MIME type.
3. GET Media đang tải Media trong vòng lặp, gây N+1 query.
4. Sync cần transaction và chiến lược unique/soft-delete rõ ràng.
5. Soft delete Post không kích hoạt MySQL cascade.
6. Post List chưa trả thumbnail bằng một JOIN duy nhất.
7. DTO Post Media đang dùng `snake_case`, trong khi DTO Post dùng `camelCase`; cần map tập trung.

## 3. Kiến trúc module và trách nhiệm

```text
Gin Controller
    -> PostMedia/PostImage Use Case
        -> Post Repository + PostMedia Repository + Media Repository
            -> MySQL transaction
        -> Redis cache invalidation
```

### 3.1. Controller

Chỉ bind URI/query/body, lấy identity đã xác thực từ Gin Context, gọi Use Case và trả response/error. Không kiểm tra ownership hoặc thao tác GORM.

### 3.2. Service/Use Case

Mở rộng `PostMediaService` hoặc tạo `PostImageService` để:

- Validate Post, Media, collection, MIME và ownership.
- Enforce một thumbnail.
- Sync một collection trong transaction.
- Đổi Media `temporary` thành `attached` sau liên kết thành công.
- Kiểm tra Media còn được dùng ở Product/Post trước khi dọn.
- Invalidate cache sau commit.
- Orchestrate tạo/cập nhật Post cùng Media khi request Post có nested `media`.

### 3.3. Repository

- `PostRepository`: truy vấn Post và Post List có thumbnail bằng `LEFT JOIN`.
- `PostMediaRepository`: query JOIN Media, create/delete/sync và soft-delete theo Post.
- `MediaRepository`: bulk lookup theo ID, kiểm tra owner/status/MIME và cập nhật status.
- Mọi repository nhận `context.Context` và dùng chung `*gorm.DB` transaction khi cần.

## 4. Thiết kế dữ liệu

### 4.1. Entity `PostMedia`

```go
type PostMediaCollection string

const (
	PostMediaCollectionThumbnail  PostMediaCollection = "thumbnail"
	PostMediaCollectionGallery    PostMediaCollection = "gallery"
	PostMediaCollectionAttachment PostMediaCollection = "attachment"
	PostMediaCollectionContent    PostMediaCollection = "content"
)

type PostMedia struct {
	ID         uint                `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	PostID     uint                `gorm:"column:post_id;not null;index:idx_post_media_post_collection,priority:1" json:"postId"`
	MediaID    uint                `gorm:"column:media_id;not null;index:idx_post_media_media_id" json:"mediaId"`
	Collection PostMediaCollection `gorm:"column:collection;type:varchar(50);not null;index:idx_post_media_post_collection,priority:2" json:"collection"`
	SortOrder  int                 `gorm:"column:sort_order;not null;default:0;index:idx_post_media_post_collection,priority:3" json:"sortOrder"`
	CreatedAt  time.Time           `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time           `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	DeletedAt  gorm.DeletedAt      `gorm:"column:deleted_at;index" json:"-"`
	Media      *mediaEntity.Media  `gorm:"foreignKey:MediaID;references:ID" json:"media,omitempty"`
}

func (PostMedia) TableName() string { return "post_media" }
```

`Post` không bắt buộc có field ngược `[]PostMedia`; repository dùng JOIN để tránh preload ngoài ý muốn.

### 4.2. Entity `Media`

Giữ các field hiện có: `ID`, `OriginalName`, `FileName`, `MimeType`, `Size`, `OriginalUrl`, `ThumbnailUrl`, `MediumUrl`, `Status`, `OwnerID`, `ProductID`, timestamps và `gorm.DeletedAt`, với các tag column/json hiện tại.

Không sử dụng `ProductID` cho Post. `status = attached` nghĩa Media đang có ít nhất một liên kết hợp lệ; khi gỡ liên kết phải kiểm tra cả `post_media` và liên kết Product trước khi đánh dấu lại trạng thái/dọn file.

### 4.3. Ràng buộc MySQL và migration

Tạo migration mới, không sửa migration đã chạy:

```sql
ALTER TABLE post_media
  ADD CONSTRAINT fk_post_media_post
    FOREIGN KEY (post_id) REFERENCES posts(id)
    ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT fk_post_media_media
    FOREIGN KEY (media_id) REFERENCES media(id)
    ON DELETE CASCADE ON UPDATE CASCADE,
  ADD UNIQUE KEY uq_post_media_link (post_id, media_id, collection),
  ADD KEY idx_post_media_post_collection_sort
    (post_id, collection, sort_order, deleted_at);
```

Trước khi thêm unique index phải xử lý record trùng. Khi sync collection, hard-delete liên kết cũ trong transaction hoặc dùng upsert/restore có kiểm soát; không để soft-deleted record ngăn việc gắn lại cùng Media.

Vì `posts` dùng soft delete, Post Service phải soft-delete các liên kết `post_media` trong cùng transaction khi xóa Post. FK cascade chỉ áp dụng khi hard-delete thực sự.

## 5. DTO và API contract

### 5.1. DTO dùng chung

```go
type MediaResponse struct {
	ID           uint   `json:"id"`
	FileName     string `json:"fileName"`
	OriginalUrl  string `json:"originalUrl"`
	ThumbnailUrl string `json:"thumbnailUrl,omitempty"`
	MediumUrl    string `json:"mediumUrl,omitempty"`
	MimeType     string `json:"mimeType"`
	Size         int64  `json:"size"`
	Status       string `json:"status"`
}

type PostMediaResponse struct {
	ID         uint           `json:"id"`
	PostID     uint           `json:"postId"`
	MediaID    uint           `json:"mediaId"`
	Collection string         `json:"collection"`
	SortOrder  int            `json:"sortOrder"`
	Media      *MediaResponse `json:"media,omitempty"`
}

type MediaSyncItem struct {
	ID        uint `json:"id" binding:"required,min=1"`
	SortOrder int  `json:"sortOrder" binding:"gte=0"`
}

type SyncPostMediaRequest struct {
	Media []MediaSyncItem `json:"media" binding:"required,dive"`
}

type PostMediaPayload struct {
	ThumbnailID *uint           `json:"thumbnailId"`
	Gallery     []MediaSyncItem `json:"gallery" binding:"max=100,dive"`
}
```

Nếu cần duy trì client cũ đang gửi `post_id`, `media_id`, `sort_order`, chỉ map alias ở transport layer; domain và response mới dùng `camelCase` thống nhất với Post/Product.

### 5.2. GET danh sách Media của Post

```text
GET /api/v1/admin/posts/:post_id/media?collection=thumbnail|gallery
```

Middleware: `AuthMiddleware` và role `ADMIN`/`STAFF`.

```go
type GetPostMediaURI struct {
	PostID uint `uri:"post_id" binding:"required,min=1"`
}

type GetPostMediaQuery struct {
	Collection string `form:"collection" binding:"omitempty,oneof=thumbnail gallery attachment content"`
}

type PaginatedPostMediaResponse struct {
	Data  []PostMediaResponse `json:"data"`
	Total int64               `json:"total"`
}
```

Repository phải trả Media bằng một JOIN, sắp xếp `sort_order ASC, id ASC`, không gọi `GetMediaByID` trong vòng lặp. Lỗi: `400` URI/query, `404` Post, `401/403` auth, `500` hệ thống.

### 5.3. POST liên kết một Media

```text
POST /api/v1/admin/posts/:post_id/media
```

Middleware: `AuthMiddleware` + `RoleMiddleware("ADMIN", "STAFF")`.

```go
type CreatePostMediaRequest struct {
	MediaID    uint                `json:"mediaId" binding:"required,min=1"`
	Collection PostMediaCollection `json:"collection" binding:"required,oneof=thumbnail gallery attachment content"`
	SortOrder  int                 `json:"sortOrder" binding:"gte=0"`
}
```

Service kiểm tra Post, Media, owner/role, MIME và duplicate. Với `thumbnail`, nếu đã có ảnh thì replace trong transaction. Response `201`: `{ "data": PostMediaResponse }`.

### 5.4. PUT sync một collection

```text
PUT /api/v1/admin/posts/:post_id/media/:collection
```

`collection` chỉ nhận `thumbnail`, `gallery` trong giai đoạn đầu. Middleware: `AuthMiddleware` + `RoleMiddleware("ADMIN", "STAFF")`.

Quy tắc:

- `thumbnail`: mảng tối đa một item.
- `gallery`: tối đa 100 item, ID không trùng, `sortOrder >= 0`.
- Tất cả Media phải tồn tại, không bị xóa, đúng MIME và thuộc owner; ADMIN có thể thao tác Media của người khác.
- Transaction xác thực toàn bộ danh sách trước, gỡ liên kết collection cũ, tạo danh sách mới, cập nhật Media status và commit.
- Cache chỉ invalidate sau commit.

Request:

```json
{
  "media": [
    { "id": 101, "sortOrder": 0 },
    { "id": 102, "sortOrder": 1 }
  ]
}
```

Response `200`:

```go
type SyncPostMediaResponse struct {
	Message string              `json:"message"`
	Data    []PostMediaResponse `json:"data"`
}
```

Lỗi: `400` validation, `403` ownership/role, `404` Post/Media, `409` duplicate/race condition.

### 5.5. DELETE gỡ liên kết

```text
DELETE /api/v1/admin/posts/:post_id/media/:media_id?collection=gallery
```

Middleware: `AuthMiddleware` + `RoleMiddleware("ADMIN", "STAFF")`.

```go
type DeletePostMediaURI struct {
	PostID  uint `uri:"post_id" binding:"required,min=1"`
	MediaID uint `uri:"media_id" binding:"required,min=1"`
}

type DeletePostMediaQuery struct {
	Collection string `form:"collection" binding:"omitempty,oneof=thumbnail gallery attachment content"`
}
```

Service xóa liên kết trong transaction, không xóa file nếu Media vẫn được dùng ở nơi khác. Response `204 No Content`; lỗi nghiệp vụ trả `404` hoặc `403` phù hợp.

### 5.6. Tạo/Cập nhật Post kèm Media

Để tránh trạng thái Post đã tạo nhưng sync ảnh thất bại, mở rộng DTO Post:

```go
type CreatePostRequest struct {
	TypeCode   string            `json:"typeCode" binding:"required,max=50"`
	Title      string            `json:"title" binding:"required,max=255"`
	Slug       string            `json:"slug" binding:"required,max=255"`
	Content    string            `json:"content" binding:"required"`
	CategoryID *uint             `json:"categoryId"`
	Media      *PostMediaPayload `json:"media"`
}

type UpdatePostRequest struct {
	Title      string            `json:"title" binding:"required,max=255"`
	Slug       string            `json:"slug" binding:"required,max=255"`
	Content    string            `json:"content" binding:"required"`
	CategoryID *uint             `json:"categoryId"`
	Media      *PostMediaPayload `json:"media"`
}
```

`Media == nil` ở Update nghĩa là không thay đổi liên kết ảnh; `Media != nil` nghĩa là sync thumbnail/gallery theo payload. Create tạo Post và liên kết Media trong một DB transaction. Nếu chưa thể mở rộng Post Service, Frontend được phép fallback POST/PUT Post rồi gọi hai API sync, nhưng phải có retry và thông báo lỗi rõ ràng.

Post response/list bổ sung thumbnail:

```go
type PostResponse struct {
	ID           uint      `json:"id"`
	TypeCode     string    `json:"typeCode"`
	Title        string    `json:"title"`
	Slug         string    `json:"slug"`
	Content      string    `json:"content,omitempty"`
	AuthorID     uint      `json:"authorId"`
	CategoryID   *uint     `json:"categoryId"`
	ThumbnailURL *string   `json:"thumbnailUrl,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
```

Post List chỉ JOIN collection `thumbnail` và lấy URL cần thiết; không gọi API Media theo từng dòng.

## 6. Transaction và nghiệp vụ trạng thái Media

### 6.1. Luồng Create

1. Controller lấy `authorID` từ JWT context, không nhận từ body.
2. Service kiểm tra slug/category/type và Media payload.
3. Mở `db.Transaction`.
4. Tạo Post.
5. Lock/read Media IDs, xác thực owner/role/status/MIME.
6. Tạo PostMedia thumbnail/gallery, enforce unique và sort order.
7. Cập nhật Media được dùng thành `attached`.
8. Commit; invalidate Post/Media cache.

Nếu bất kỳ bước nào thất bại, rollback toàn bộ Post và liên kết; file temporary vẫn do cleanup job xử lý.

### 6.2. Luồng Update và Sync

Update Post và Sync Media trong cùng transaction khi request có `media`. Không gỡ liên kết cũ trước khi toàn bộ Media mới được validate để tránh mất dữ liệu khi một ID sai.

### 6.3. Xóa Post

Post Service soft-delete Post và toàn bộ `post_media` trong một transaction. Không xóa Media vật lý vì Media có thể được chia sẻ.

## 7. Cache Redis

### 7.1. Key pattern

```text
cache:post:{post_id}:media:all
cache:post:{post_id}:media:{collection}
cache:admin:posts:list:{hash_of_query}
cache:admin:posts:detail:{post_id}
```

TTL đề xuất: 24 giờ cho Media Post, giảm nếu nội dung thay đổi thường xuyên.

### 7.2. Cache-aside

GET Media/List Post đọc Redis trước. Cache miss query JOIN MySQL, serialize DTO và `SetEx`. Không dùng `KEYS` trong production; dùng `SCAN` theo pattern hoặc lưu danh sách key liên quan.

Sau commit phải invalidate:

- `cache:post:{id}:media:*`
- `cache:admin:posts:detail:{id}`
- các `cache:admin:posts:list:*` và list theo `typeCode`.

### 7.3. Thumbnail trong Post List

```sql
SELECT p.id, p.title, p.slug, p.type_code, p.created_at,
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
LIMIT ? OFFSET ?;
```

Index bắt buộc: `(post_id, collection, sort_order, deleted_at)` trên `post_media`, cùng các index hiện có trên `posts` và `media`.

## 8. Auth, ownership và Redis token state

### 8.1. Middleware

Tái sử dụng `AuthMiddleware(redisClient)` hiện có: xác minh chữ ký JWT trước, lấy JTI kiểm tra blacklist Redis, sau đó set `userId` và `role` vào Context. Tất cả API write của Media/Post Media và nhóm Admin đều phải đi qua middleware.

Không tin `ownerId`, `userId` hoặc `role` từ body/query. Service lấy identity đã được middleware xác thực và kiểm tra kiểu trước khi xử lý.

### 8.2. Phân quyền

- `ADMIN`: xem và quản lý mọi Post Media.
- `STAFF`: quản lý Post Media trong quyền Dashboard, chỉ gắn Media được cấp quyền.
- `CUSTOMER` không truy cập route quản trị.
- DELETE temporary Media kiểm tra `media.owner_id == userID`, trừ ADMIN.

### 8.3. Refresh/Logout

Không thay đổi contract Auth hiện có. Refresh Token rotation phải thu hồi JTI cũ và tạo JTI mới trong atomic operation; logout đưa JTI/access token vào blacklist bằng `SetEx` với TTL còn lại. Auth Middleware dùng key hiện tại `auth:bl:at:{jti}` để chặn request sau logout.

## 9. Error contract

```go
type ErrorResponse struct {
	Error   string            `json:"error"`
	Code    string            `json:"code,omitempty"`
	Details map[string]string `json:"details,omitempty"`
}
```

| Tình huống | HTTP | Code |
| --- | ---: | --- |
| URI/query/body sai | 400 | `VALIDATION_ERROR` |
| Post/Media không tồn tại | 404 | `NOT_FOUND` |
| User không sở hữu Media | 403 | `FORBIDDEN` |
| Thumbnail/Media trùng hoặc race | 409 | `MEDIA_CONFLICT` |
| Token hết hạn/blacklist | 401 | `UNAUTHORIZED` |
| DB/Redis/storage failure | 500 | `INTERNAL_ERROR` |

Controller không trả `500` cho lỗi validation hoặc lỗi nghiệp vụ đã biết.

## 10. Google Wire và thay đổi module

Cập nhật providers trong `internal/di/wire.go` và regenerate `wire_gen.go`:

- `NewPostMediaRepository` với query JOIN/transaction mới.
- `NewMediaRepository` bổ sung bulk validation/status methods.
- `NewPostMediaService` hoặc `NewPostImageService` nhận `*gorm.DB`, repositories và Redis client.
- `NewPostController` nhận use case orchestrate Post + Media nếu bật nested payload.

Không tạo dependency ngược từ Controller xuống Repository; không đưa business rule vào `main.go` hoặc Handler.

## 11. Trình tự triển khai

1. Migration FK/index/unique và xử lý record trùng hiện có.
2. Entity enum/associations và DTO typed.
3. Refactor PostMedia Repository thành query JOIN, bulk Media lookup và transaction sync.
4. Bổ sung validation/ownership/status trong Service.
5. Bổ sung nested Media transaction cho Create/Update Post.
6. Bổ sung Post List thumbnail JOIN và response `thumbnailUrl`.
7. Chuẩn hóa Controller status code/error contract.
8. Cập nhật Wire, route role middleware và OpenAPI/API contract.
9. Bổ sung Redis cache/invalidation và kiểm thử cache hit/miss.
10. Viết test cho rollback, thumbnail uniqueness, ownership, soft delete và N+1 prevention.

## 12. Acceptance Criteria

- Tạo Post với thumbnail/gallery thành công trong một transaction.
- Một Post không thể có hơn một thumbnail.
- Gallery giữ đúng thứ tự sau sync và reload.
- Media không thuộc user khác không thể bị liên kết trái phép.
- Xóa Post không xóa Media đang được dùng ở nơi khác.
- GET Post Media không phát sinh N+1 query.
- Post List lấy thumbnail bằng một JOIN và không gọi API Media theo từng dòng.
- Cache được invalidate sau commit, không phục vụ thumbnail cũ.
- Tất cả API write trả status/error đúng contract và được bảo vệ bởi JWT/JTI blacklist.
