# KẾ HOẠCH BACKEND: QUẢN LÝ HÌNH ẢNH TRANG (PAGE MEDIA)

**Dự án:** TechBite  
**Tài liệu nguồn:** `.docs/ideas/dashboard/24-page-media-management-idea.md`  
**Phạm vi:** Liên kết `thumbnail` và `gallery` với Page đã tồn tại qua bảng `page_media` độc lập.

## 1. Hiện trạng và nguyên tắc

- Không dùng `post_media`: bảng này có foreign key tới `posts`. Page Media có entity, repository, service, controller, migration và Redis namespace riêng.
- Các kế hoạch Page List/Create/Edit đã quy hoạch `pages` và module `Page`, nhưng code/migration Page hiện chưa có trong `apps/backend`; đây là prerequisite bắt buộc.
- Tái sử dụng Gin, GORM, Google Wire, `AuthMiddleware`, `RoleMiddleware`, Media entity và Redis client. Controller chỉ bind URI/query/body, đọc identity từ Context, gọi service và map lỗi.
- `POST /api/v1/media/upload` hiện đã được bảo vệ bởi AuthMiddleware nhưng `MediaController` dùng `mockUserId`. Phải sửa shared Media API lấy `userId`/`role` từ JWT Context trước khi áp rule Staff ownership.
- Không xóa file Media khi gỡ Page Media. Mọi mutation link dùng transaction MySQL; Redis lỗi chỉ log sau commit, không rollback DB.

## 2. Trụ cột 1 — Thiết kế dữ liệu

### 2.1. GORM entity

Tạo `apps/backend/internal/domain/page/entity/page_media.go`:

```go
package entity

import (
    "time"

    mediaEntity "go_refine_dashboard_be/internal/domain/media/entity"
    "gorm.io/gorm"
)

type PageMediaCollection string

const (
    PageMediaCollectionThumbnail PageMediaCollection = "thumbnail"
    PageMediaCollectionGallery   PageMediaCollection = "gallery"
)

type PageMedia struct {
    ID         uint                `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
    PageID     uint                `gorm:"column:page_id;not null;index:idx_page_media_page_collection_sort,priority:1;uniqueIndex:uq_page_media_link,priority:1" json:"pageId"`
    MediaID    uint                `gorm:"column:media_id;not null;index:idx_page_media_media_id;uniqueIndex:uq_page_media_link,priority:2" json:"mediaId"`
    Collection PageMediaCollection `gorm:"column:collection;type:varchar(20);not null;index:idx_page_media_page_collection_sort,priority:2;uniqueIndex:uq_page_media_link,priority:3" json:"collection"`
    SortOrder  int                 `gorm:"column:sort_order;not null;default:0;index:idx_page_media_page_collection_sort,priority:3" json:"sortOrder"`

    Media     *mediaEntity.Media `gorm:"foreignKey:MediaID;references:ID" json:"media,omitempty"`
    CreatedAt time.Time          `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
    UpdatedAt time.Time          `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
    DeletedAt gorm.DeletedAt     `gorm:"column:deleted_at;index:idx_page_media_page_collection_sort,priority:4" json:"-"`
}

func (PageMedia) TableName() string { return "page_media" }
```

- Pivot soft-delete chỉ đánh dấu liên kết, không ảnh hưởng Media/file gốc.
- `uq_page_media_link(page_id, media_id, collection)` cấm link trùng. Repository phải restore row soft-deleted bằng `Unscoped` thay vì insert duplicate.
- MySQL không có partial unique index thuận tiện cho “một thumbnail active” cùng soft delete. Mọi write phải `SELECT ... FOR UPDATE` row Page trong transaction, rồi service enforce tối đa một thumbnail; Page lock tuần tự hóa các thao tác cùng Page.

### 2.2. Migration MySQL

Tạo migration mới `database/migrations/<timestamp>_create_page_media_table.up.sql`, không sửa migration `post_media`:

```sql
CREATE TABLE page_media (
  id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  page_id INT UNSIGNED NOT NULL,
  media_id INT UNSIGNED NOT NULL,
  collection VARCHAR(20) NOT NULL,
  sort_order INT NOT NULL DEFAULT 0,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  deleted_at DATETIME(3) NULL DEFAULT NULL,
  CONSTRAINT chk_page_media_collection CHECK (collection IN ('thumbnail', 'gallery')),
  CONSTRAINT chk_page_media_sort_order CHECK (sort_order >= 0),
  CONSTRAINT fk_page_media_page FOREIGN KEY (page_id) REFERENCES pages(id)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT fk_page_media_media FOREIGN KEY (media_id) REFERENCES media(id)
    ON DELETE CASCADE ON UPDATE CASCADE,
  UNIQUE KEY uq_page_media_link (page_id, media_id, collection),
  KEY idx_page_media_page_collection_sort (page_id, collection, sort_order, deleted_at),
  KEY idx_page_media_media_id (media_id),
  KEY idx_page_media_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

Migration down chỉ `DROP TABLE page_media`. Nếu Page bị soft-delete, Page service phải soft-delete pivot trong transaction; FK cascade chỉ là tuyến phòng thủ khi hard-delete/migration.

### 2.3. Repository

Tạo `internal/infrastructure/repository/page_media_repo.go`:

```go
type PageMediaRepository interface {
    FindPageForUpdate(ctx context.Context, tx *gorm.DB, pageID uint) (*pageEntity.Page, error)
    FindByPageID(ctx context.Context, pageID uint, collection string) ([]pageEntity.PageMedia, int64, error)
    FindMediaByIDsForUpdate(ctx context.Context, tx *gorm.DB, ids []uint) ([]mediaEntity.Media, error)
    CreateOrRestore(ctx context.Context, tx *gorm.DB, item *pageEntity.PageMedia) error
    SoftDeleteNotIn(ctx context.Context, tx *gorm.DB, pageID uint, collection string, mediaIDs []uint) error
    SoftDeleteLink(ctx context.Context, tx *gorm.DB, pageID, mediaID uint, collection string) (bool, error)
    MarkAttached(ctx context.Context, tx *gorm.DB, mediaIDs []uint) error
}
```

- GET là một query có `Preload("Media")`, optional collection, `ORDER BY sort_order ASC, id ASC`; không có N+1 lookup Media.
- Media được load theo batch `WHERE id IN ?` và lock trước khi mutation; GORM scoped query tự loại row soft-deleted.
- `CreateOrRestore` khôi phục link đã xóa hoặc update `sort_order`; chỉ insert nếu link chưa từng tồn tại.
- `SoftDeleteNotIn` xóa mềm toàn bộ link không thuộc snapshot, nhận `[]` để clear collection. Không `Unscoped().Delete` pivot trong flow nghiệp vụ bình thường.

## 3. Trụ cột 2 — API contract và DDD/CQRS

### 3.1. Module và Wire

```text
apps/backend/internal/
├── controller/page_media_controller.go
├── domain/page/entity/page_media.go
├── infrastructure/repository/page_media_repo.go
└── usecases/pagemedia/
    ├── dto/page_media_dto.go
    └── service/page_media_service.go
```

Thêm vào `internal/di/wire.go`, rồi regenerate `wire_gen.go`:

```go
var PageMediaSet = wire.NewSet(
    repository.NewPageMediaRepository,
    pagemediaService.NewPageMediaService,
    controller.NewPageMediaController,
)

func InitializePageMediaController(db *gorm.DB, redisClient *redis.Client) *controller.PageMediaController
```

Service nhận `PageMediaRepository`, `*gorm.DB`, `*redis.Client`. Không đưa GORM/Redis vào Handler.

### 3.2. DTO

Tạo `internal/usecases/pagemedia/dto/page_media_dto.go`:

```go
type MediaResponse struct {
    ID           uint   `json:"id"`
    FileName     string `json:"fileName"`
    OriginalURL  string `json:"originalUrl"`
    ThumbnailURL string `json:"thumbnailUrl,omitempty"`
    MediumURL    string `json:"mediumUrl,omitempty"`
    MimeType     string `json:"mimeType"`
    Size         int64  `json:"size"`
    Status       string `json:"status"`
}

type PageMediaResponse struct {
    ID         uint           `json:"id"`
    PageID     uint           `json:"pageId"`
    MediaID    uint           `json:"mediaId"`
    Collection string         `json:"collection"`
    SortOrder  int            `json:"sortOrder"`
    Media      *MediaResponse `json:"media,omitempty"`
}

type CreatePageMediaRequest struct {
    MediaID    uint   `json:"media_id" binding:"required,min=1"`
    Collection string `json:"collection" binding:"required,oneof=thumbnail gallery"`
    SortOrder  *int   `json:"sort_order" binding:"omitempty,gte=0"`
}

type MediaSyncItem struct {
    ID        uint `json:"id" binding:"required,min=1"`
    SortOrder int  `json:"sort_order" binding:"gte=0"`
}

type SyncPageMediaRequest struct {
    Media []MediaSyncItem `json:"media" binding:"max=100,dive"`
}
```

- Response chuẩn camelCase; request giữ `media_id`/`sort_order` tương thích API Post Media hiện có. Adapter Frontend map `mediaId`/`sortOrder` tại một boundary duy nhất.
- `media` không `required`: mảng rỗng là valid để gỡ thumbnail/xóa gallery. Response tái dùng Media DTO an toàn, không trả filesystem path, owner metadata, token hay password.

### 3.3. Route và auth context

```go
pageMedia := admin.Group("/pages/:page_id/media")
pageMedia.Use(middleware.RoleMiddleware("ADMIN", "STAFF"))
{
    pageMedia.GET("", pageMediaController.GetPageMedia)
    pageMedia.POST("", pageMediaController.CreatePageMedia)
    pageMedia.PUT("/:collection", pageMediaController.SyncPageMedia)
    pageMedia.DELETE("/:media_id", pageMediaController.DeletePageMedia)
}
```

`admin` đã có `AuthMiddleware(redisClient)`: middleware xác minh signature JWT trước, kiểm tra `auth:bl:at:{jti}`, sau đó set `userId` (`uint`) và `role` (`string`). Controller type-assert identity này, không đọc role/owner/user từ input client.

### 3.4. Endpoints

#### GET `/api/v1/admin/pages/:page_id/media?collection=thumbnail|gallery`

```go
type GetPageMediaURI struct { PageID uint `uri:"page_id" binding:"required,min=1"` }
type GetPageMediaQuery struct { Collection string `form:"collection" binding:"omitempty,oneof=thumbnail gallery"` }
type PaginatedPageMediaResponse struct {
    Data []PageMediaResponse `json:"data"`
    Total int64 `json:"total"`
}
```

Success `200`: `{ "data": [], "total": 0 }`; empty là hợp lệ. Page không tồn tại/soft-deleted trả `404`.

#### POST `/api/v1/admin/pages/:page_id/media`

Body `CreatePageMediaRequest`, success `201`: `{ "data": PageMediaResponse }`.

1. Lock/check Page active.
2. Validate Media tồn tại, không soft-delete, `mime_type` là `image/*` và quyền ownership.
3. `ADMIN` được gắn mọi Media hợp lệ; `STAFF` chỉ khi `media.owner_id == userId`.
4. Thumbnail soft-delete thumbnail active cũ, attach/restore ảnh mới tại `sort_order=0` trong một transaction.
5. Gallery dùng `sort_order` request hoặc `MAX(sort_order)+1` trong Page lock. Media mới được đánh dấu `attached`; Media cũ không bị xóa/đổi temporary.

#### PUT `/api/v1/admin/pages/:page_id/media/:collection`

```go
type SyncPageMediaURI struct {
    PageID uint `uri:"page_id" binding:"required,min=1"`
    Collection string `uri:"collection" binding:"required,oneof=thumbnail gallery"`
}
```

```json
{
  "media": [
    { "id": 102, "sort_order": 0 },
    { "id": 103, "sort_order": 1 }
  ]
}
```

- Thumbnail tối đa một item và, nếu có, `sort_order=0`.
- Gallery tối đa 100 item; ID/sort order unique, order liên tục `0..n-1`.
- Service validate đầy đủ Page/Media/MIME/ownership/duplicates trước destructive change. Trong một transaction: lock Page → soft-delete links không thuộc snapshot → restore/create snapshot → mark Media mới `attached` → commit.
- `200`: `{ "message": "Đồng bộ hình ảnh trang thành công", "data": [] }`, `data` được query lại sau commit theo thứ tự mới.

#### DELETE `/api/v1/admin/pages/:page_id/media/:media_id?collection=thumbnail|gallery`

```go
type DeletePageMediaURI struct {
    PageID uint `uri:"page_id" binding:"required,min=1"`
    MediaID uint `uri:"media_id" binding:"required,min=1"`
}
type DeletePageMediaQuery struct { Collection string `form:"collection" binding:"omitempty,oneof=thumbnail gallery"` }
```

- Nếu `collection` rỗng và Media có link ở nhiều collection, trả `400 collection_required`; nếu xác định một link thì được gỡ.
- Service lock Page, kiểm tra quyền Media owner/role, soft-delete pivot, commit và invalidate cache. Không có active link trả `404`; success `204 No Content`.
- Endpoint không gọi `MediaRepository.Delete`, không xóa DB/file Media và không đổi trạng thái Media.

### 3.5. Service, transaction và errors

```go
type PageMediaService interface {
    GetPageMedia(ctx context.Context, pageID uint, collection string) ([]dto.PageMediaResponse, int64, error)
    CreatePageMedia(ctx context.Context, pageID uint, req *dto.CreatePageMediaRequest, userID uint, role string) (*dto.PageMediaResponse, error)
    SyncPageMedia(ctx context.Context, pageID uint, collection string, req *dto.SyncPageMediaRequest, userID uint, role string) ([]dto.PageMediaResponse, error)
    DeletePageMedia(ctx context.Context, pageID, mediaID uint, collection string, userID uint, role string) error
}
```

Domain errors: `ErrPageNotFound`, `ErrMediaNotFound`, `ErrMediaForbidden`, `ErrCollectionInvalid`, `ErrThumbnailLimit`, `ErrMediaDuplicate`, `ErrSortOrderInvalid`, `ErrMediaTypeInvalid`, `ErrPageMediaNotFound`, `ErrCollectionRequired`.

| Điều kiện | HTTP | Code |
| --- | ---: | --- |
| URI/query/body/collection/sort sai | 400 | `VALIDATION_ERROR` |
| Page, Media hoặc link không tồn tại | 404 | `NOT_FOUND` |
| Không đủ role hoặc Staff không sở hữu Media | 403 | `FORBIDDEN` |
| Duplicate/race khi attach | 409 | `MEDIA_CONFLICT` |
| JWT missing/invalid/revoked | 401 | `UNAUTHORIZED` |
| MySQL/Redis bất ngờ | 500 | `INTERNAL_ERROR` |

Controller map domain error sang `ErrorResponse{Error, Code, Details?}`, không lộ SQL, JWT claims hoặc owner ID của người khác.

## 4. Củng cố shared Media API

- `MediaController.UploadFile`/`DeleteFile` thay `mockUserId` bằng identity từ Gin Context, sau AuthMiddleware.
- `MediaService.DeleteFile` chỉ xóa Media `temporary`: ADMIN có thể dọn mọi file temporary, STAFF chỉ dọn file do mình sở hữu; Media `attached` luôn trả `409 MEDIA_ATTACHED`. Gỡ Page Media luôn gọi endpoint Page Media.
- Giữ Magic Bytes JPEG/PNG/WEBP và 5 MB. Map MIME/size sai → `400`, ownership → `403`, missing Media → `404`, không dùng `500` cho business error.

## 5. Trụ cột 3 — Redis cache và token state

### 5.1. Cache

```text
cache:page:{page_id}:media:all
cache:page:{page_id}:media:thumbnail
cache:page:{page_id}:media:gallery
public:pages:slug:{slug}
seo:page:{page_id}
```

- GET Media dùng cache-aside: cache miss query/preload MySQL, serialize DTO, `SetEx` TTL 15 phút. Cache response rỗng là hợp lệ.
- Sau POST/PUT/DELETE commit, xóa trực tiếp ba Page Media keys, public slug cache và SEO resolved cache (thumbnail là SEO fallback). Không dùng Redis `KEYS`.
- Public Page API là hạng mục sau sẽ ghi `public:pages:slug:*`; Page Media chỉ invalidates key đã biết. Redis failure log/fallback DB, không rollback mutation đã commit.

### 5.2. Token state

- Không tạo token/session key mới. Protected routes kế thừa blacklist `auth:bl:at:{jti}` và Role middleware.
- Logout giữ `SetEx` theo TTL access token còn lại. Refresh rotation là luồng dùng chung: revoke refresh JTI cũ, tạo token/JTI mới và lưu session Redis phải atomic (Redis transaction/optimistic locking), tránh replay; Page Media không tự refresh token.

## 6. Trình tự triển khai

1. Thực thi Page List/Create/Edit baseline (`pages`, Page module, Wire, routes).
2. Harden shared Media identity/delete guard.
3. Thêm migration/entity `page_media`.
4. Viết repository lock/preload/restore/soft-delete methods.
5. Viết DTO, PageMediaService rules, transactions và Redis invalidation.
6. Thêm thin controller, routes, `PageMediaSet`; regenerate `wire_gen.go`.
7. Cập nhật `.docs/api-endpoints.yaml` với bốn endpoints, schemas, BearerAuth và error responses.
8. Tích hợp Frontend sau khi API contracts được kiểm thử.

## 7. Acceptance criteria

- Migration tạo FK/index/check constraints; không có collection sai hoặc `sort_order < 0`.
- GET không N+1, gallery trả `sort_order ASC, id ASC`; Page soft-deleted trả `404`.
- POST thumbnail replace atomic; concurrent requests không để lại hai thumbnail active.
- PUT validates toàn bộ trước mutation; ID trùng, gap order, MIME sai, Media missing hay Staff dùng Media người khác không thay đổi collection cũ.
- Sync `[]`, gỡ và reorder chỉ thay link, không xóa Media/file hay chuyển Media dùng chung về temporary.
- Admin dùng Media mọi owner; Staff chỉ own Media; Customer `403`; JWT invalid/revoked `401` trước Service.
- Upload ghi đúng owner JWT, không còn `mockUserId`; DELETE shared Media attached bị chặn.
- Cache chỉ invalidates sau DB commit; public Page/SEO cache không stale sau thumbnail/gallery change.
- Controller không chứa business/GORM/Redis logic; Wire generated file được regenerate; response không lộ dữ liệu nhạy cảm.
