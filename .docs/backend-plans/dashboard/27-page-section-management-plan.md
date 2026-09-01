# KẾ HOẠCH BACKEND: QUẢN LÝ SECTION CỦA TỪNG PAGE

**Dự án:** TechBite

**Tài liệu nguồn:** `.docs/ideas/dashboard/26-page-section-management-idea.md`

**Giải pháp:** Semi-Structured CMS với `page_sections` và `page_section_items`

**Stack bắt buộc:** Go, Gin, GORM, Google Wire, MySQL, Redis

**Phạm vi:** Admin CRUD/reorder Section, liên kết Category/Post/Media, Media Library picker và mở rộng read model Public Page.

## 1. Hiện trạng và quyết định kiến trúc

- Module `Page`, `PageMedia`, `Media`, `Category`, `Post`, JWT middleware, Redis client và Google Wire đã tồn tại; phải tái sử dụng, không tạo auth/cache client riêng.
- Backend hiện chưa có `PageSection`, `PageSectionItem`, API Section, Media Library list/search hoặc implementation Public Page API. `.docs/backend-plans/dashboard/26-page-public-api-plan.md` là prerequisite cho phần response public.
- `CATEGORY` trong phạm vi này ánh xạ tới bảng sản phẩm `categories`, đúng theo IDEA. Không ánh xạ tới `post_categories`; nếu sau này cần danh mục bài viết phải thêm `POST_CATEGORY`, không overload `CATEGORY`.
- `item_id` là polymorphic reference nên không có foreign key chung. Use Case phải batch-validate loại, tồn tại, soft delete, quyền Media và duplicate trước khi ghi.
- Dùng CQRS thực dụng:
  - Command service xử lý create/update/delete/reorder/sync trong transaction.
  - Query service dựng Admin read model và Public projection bằng số query cố định, không N+1.
- Mọi mutation cùng Page phải lock row `pages` bằng `SELECT ... FOR UPDATE` để tuần tự hóa create/reorder/delete/sync và tránh snapshot stale.
- Section/Item dùng soft delete. Xóa vật lý Page vẫn cascade toàn bộ Section/Item; soft-delete Page giữ cấu hình nhưng mọi Admin/Public query scoped không trả Page đó.
- Admin Section không cache Redis trong v1 vì dữ liệu nhỏ, thay đổi thường xuyên và collection key động. Redis chỉ cache Public Page aggregate; mutation xóa exact key sau commit.
- Category/Post/Media là dữ liệu nguồn. Xóa Item hoặc Section chỉ xóa liên kết, tuyệt đối không xóa dữ liệu nguồn.

---

## 2. Trụ cột 1 — Thiết kế dữ liệu

### 2.1. GORM entity `PageSection`

Tạo `apps/backend/internal/domain/page/entity/page_section.go`:

```go
package entity

import (
    "time"

    mediaEntity "go_refine_dashboard_be/internal/domain/media/entity"
    "gorm.io/gorm"
)

type PageSectionStatus string

const (
    PageSectionStatusActive   PageSectionStatus = "ACTIVE"
    PageSectionStatusInactive PageSectionStatus = "INACTIVE"
)

type PageSection struct {
    ID     uint   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
    PageID uint   `gorm:"column:page_id;not null;uniqueIndex:uq_page_sections_page_key,priority:1;index:idx_page_sections_page_status_sort,priority:1" json:"pageId"`
    Key    string `gorm:"column:key;type:varchar(100);not null;uniqueIndex:uq_page_sections_page_key,priority:2" json:"key"`
    Name   string `gorm:"column:name;type:varchar(255);not null" json:"name"`

    Title       string `gorm:"column:title;type:varchar(255);not null;default:''" json:"title"`
    Description string `gorm:"column:description;type:text;not null" json:"description"`

    BackgroundColor   *string `gorm:"column:background_color;type:varchar(9)" json:"backgroundColor"`
    BackgroundMediaID *uint   `gorm:"column:background_media_id;index:idx_page_sections_background_media" json:"backgroundMediaId"`
    FeatureMediaID    *uint   `gorm:"column:feature_media_id;index:idx_page_sections_feature_media" json:"featureMediaId"`

    SortOrder int               `gorm:"column:sort_order;not null;default:0;index:idx_page_sections_page_status_sort,priority:3" json:"sortOrder"`
    Status    PageSectionStatus `gorm:"column:status;type:varchar(20);not null;default:'ACTIVE';index:idx_page_sections_page_status_sort,priority:2" json:"status"`

    Page            *Page                `gorm:"foreignKey:PageID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
    BackgroundMedia *mediaEntity.Media   `gorm:"foreignKey:BackgroundMediaID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"backgroundMedia,omitempty"`
    FeatureMedia    *mediaEntity.Media   `gorm:"foreignKey:FeatureMediaID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"featureMedia,omitempty"`
    Items           []PageSectionItem    `gorm:"foreignKey:SectionID;references:ID" json:"-"`
    CreatedAt       time.Time            `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
    UpdatedAt       time.Time            `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
    DeletedAt       gorm.DeletedAt       `gorm:"column:deleted_at;index:idx_page_sections_page_status_sort,priority:4" json:"-"`
}

func (PageSection) TableName() string { return "page_sections" }
```

- `uq_page_sections_page_key(page_id, key)` giữ key ổn định và duy nhất trong một Page.
- Nếu tạo lại một key đã soft-delete, repository restore cùng row bằng `Unscoped`, cập nhật metadata và đặt `deleted_at = NULL`; không insert row trùng.
- `idx_page_sections_page_status_sort(page_id, status, sort_order, deleted_at)` phục vụ Admin/Public list có thứ tự.
- Background/Feature Media nullable, `ON DELETE SET NULL`; Service vẫn chặn hard-delete Media đang attached để tránh mất tài sản ngoài ý muốn.
- Có thể bổ sung `Sections []PageSection` với `json:"-"` vào `Page` nếu cần navigation GORM, nhưng query chính không dùng auto-preload.

### 2.2. GORM entity `PageSectionItem`

Tạo `apps/backend/internal/domain/page/entity/page_section_item.go`:

```go
package entity

import (
    "time"

    "gorm.io/gorm"
)

type PageSectionItemType string

const (
    PageSectionItemTypeCategory PageSectionItemType = "CATEGORY"
    PageSectionItemTypePost     PageSectionItemType = "POST"
    PageSectionItemTypeMedia    PageSectionItemType = "MEDIA"
)

type PageSectionItem struct {
    ID         uint                `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
    SectionID  uint                `gorm:"column:section_id;not null;uniqueIndex:uq_page_section_item_link,priority:1;index:idx_page_section_items_section_collection_sort,priority:1" json:"sectionId"`
    ItemType   PageSectionItemType `gorm:"column:item_type;type:varchar(20);not null;uniqueIndex:uq_page_section_item_link,priority:3;index:idx_page_section_items_source,priority:1" json:"itemType"`
    ItemID     uint                `gorm:"column:item_id;not null;uniqueIndex:uq_page_section_item_link,priority:4;index:idx_page_section_items_source,priority:2" json:"itemId"`
    Collection string              `gorm:"column:collection;type:varchar(100);not null;uniqueIndex:uq_page_section_item_link,priority:2;index:idx_page_section_items_section_collection_sort,priority:2" json:"collection"`
    SortOrder  int                 `gorm:"column:sort_order;not null;default:0;index:idx_page_section_items_section_collection_sort,priority:3" json:"sortOrder"`

    Section   *PageSection   `gorm:"foreignKey:SectionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
    CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
    UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
    DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index:idx_page_section_items_section_collection_sort,priority:4;index:idx_page_section_items_source,priority:3" json:"-"`
}

func (PageSectionItem) TableName() string { return "page_section_items" }
```

- Không khai báo GORM relation từ `ItemID` tới ba bảng nguồn.
- `uq_page_section_item_link(section_id, collection, item_type, item_id)` cấm duplicate lịch sử. Re-add dùng `CreateOrRestoreItem`, không insert duplicate.
- Một Item được phép xuất hiện ở collection khác. Một collection active chỉ được có một `item_type`; Service enforce vì MySQL không thể biểu diễn constraint này trên schema hai bảng hiện tại.
- `idx_page_section_items_source(item_type, item_id, deleted_at)` phục vụ reverse lookup để invalidate các Public Page đang tham chiếu nguồn.

### 2.3. Migration MySQL

Tạo:

```text
apps/backend/database/migrations/<timestamp>_create_page_sections_tables.up.sql
apps/backend/database/migrations/<timestamp>_create_page_sections_tables.down.sql
```

Migration up:

```sql
CREATE TABLE `page_sections` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `page_id` INT UNSIGNED NOT NULL,
  `key` VARCHAR(100) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `title` VARCHAR(255) NOT NULL DEFAULT '',
  `description` TEXT NOT NULL,
  `background_color` VARCHAR(9) NULL,
  `background_media_id` INT UNSIGNED NULL,
  `feature_media_id` INT UNSIGNED NULL,
  `sort_order` INT NOT NULL DEFAULT 0,
  `status` VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_page_sections_page_key` (`page_id`, `key`),
  KEY `idx_page_sections_page_status_sort` (`page_id`, `status`, `sort_order`, `deleted_at`),
  KEY `idx_page_sections_background_media` (`background_media_id`),
  KEY `idx_page_sections_feature_media` (`feature_media_id`),
  CONSTRAINT `chk_page_sections_status` CHECK (`status` IN ('ACTIVE', 'INACTIVE')),
  CONSTRAINT `chk_page_sections_sort_order` CHECK (`sort_order` >= 0),
  CONSTRAINT `fk_page_sections_page` FOREIGN KEY (`page_id`) REFERENCES `pages` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_page_sections_background_media` FOREIGN KEY (`background_media_id`) REFERENCES `media` (`id`)
    ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `fk_page_sections_feature_media` FOREIGN KEY (`feature_media_id`) REFERENCES `media` (`id`)
    ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `page_section_items` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `section_id` INT UNSIGNED NOT NULL,
  `item_type` VARCHAR(20) NOT NULL,
  `item_id` INT UNSIGNED NOT NULL,
  `collection` VARCHAR(100) NOT NULL,
  `sort_order` INT NOT NULL DEFAULT 0,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_page_section_item_link` (`section_id`, `collection`, `item_type`, `item_id`),
  KEY `idx_page_section_items_section_collection_sort` (`section_id`, `collection`, `sort_order`, `deleted_at`),
  KEY `idx_page_section_items_source` (`item_type`, `item_id`, `deleted_at`),
  CONSTRAINT `chk_page_section_items_type` CHECK (`item_type` IN ('CATEGORY', 'POST', 'MEDIA')),
  CONSTRAINT `chk_page_section_items_sort_order` CHECK (`sort_order` >= 0),
  CONSTRAINT `fk_page_section_items_section` FOREIGN KEY (`section_id`) REFERENCES `page_sections` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

Migration down phải drop `page_section_items` trước, sau đó `page_sections`. Không sửa migration Page/Page Media đã chạy.

---

## 3. Trụ cột 2 — DDD/CQRS, repository và transaction

### 3.1. Cấu trúc module

```text
apps/backend/internal/
├── controller/page_section_controller.go
├── domain/page/entity/
│   ├── page_section.go
│   └── page_section_item.go
├── infrastructure/
│   ├── repository/page_section_repo.go
│   └── queryservice/page_section_query_service.go
└── usecases/pagesection/
    ├── dto/page_section_dto.go
    ├── queryservice/page_section_query_service.go
    └── service/page_section_service.go
```

- Controller chỉ bind URI/query/body, lấy identity đã verify từ Gin Context, gọi command/query service và map domain error.
- Command service giữ toàn bộ validation/business rule, điều phối transaction và cache invalidation.
- Repository chỉ persistence/locking, không map HTTP hoặc quyết định quyền.
- Query service dùng projection typed để aggregate Media/collection/source data; không trả GORM entity trực tiếp.

### 3.2. Command repository

```go
type PageSectionRepository interface {
    FindPageForUpdate(ctx context.Context, tx *gorm.DB, pageID uint) (*pageEntity.Page, error)
    FindSectionForUpdate(ctx context.Context, tx *gorm.DB, pageID, sectionID uint) (*pageEntity.PageSection, error)
    FindActiveSectionsForUpdate(ctx context.Context, tx *gorm.DB, pageID uint) ([]pageEntity.PageSection, error)
    FindMediaByIDsForUpdate(ctx context.Context, tx *gorm.DB, ids []uint) ([]mediaEntity.Media, error)
    FindSourceIDsForUpdate(ctx context.Context, tx *gorm.DB, itemType pageEntity.PageSectionItemType, ids []uint) ([]uint, error)
    CreateOrRestoreSection(ctx context.Context, tx *gorm.DB, section *pageEntity.PageSection) error
    UpdateSection(ctx context.Context, tx *gorm.DB, section *pageEntity.PageSection) error
    SoftDeleteSectionAndItems(ctx context.Context, tx *gorm.DB, sectionID uint) error
    BulkUpdateSectionOrder(ctx context.Context, tx *gorm.DB, pageID uint, items []SectionOrderUpdate) error
    FindCollectionItemsForUpdate(ctx context.Context, tx *gorm.DB, sectionID uint, collection string) ([]pageEntity.PageSectionItem, error)
    CreateOrRestoreItem(ctx context.Context, tx *gorm.DB, item *pageEntity.PageSectionItem) error
    SoftDeleteItemsNotIn(ctx context.Context, tx *gorm.DB, sectionID uint, collection string, itemType pageEntity.PageSectionItemType, ids []uint) error
    SoftDeleteCollection(ctx context.Context, tx *gorm.DB, sectionID uint, collection string) error
    MarkMediaAttached(ctx context.Context, tx *gorm.DB, ids []uint) error
}
```

- `FindSourceIDsForUpdate` dùng đúng bảng theo enum, GORM default scope loại soft-deleted rows và trả batch ID; tuyệt đối không nối tên bảng từ raw client string.
- `BulkUpdateSectionOrder` dùng một `CASE id WHEN ...` hoặc prepared updates trong một transaction; không ghép sort field/value không tin cậy vào SQL.
- `CreateOrRestoreSection/Item` query `Unscoped` theo unique key. Chỉ row có `DeletedAt.Valid=true` mới được restore `deleted_at = NULL`; row đang active trả conflict, không bị ghi đè như một lệnh create.
- Lock order thống nhất: Page → Section(s) → source rows → links. Không đảo thứ tự để tránh deadlock.

### 3.3. Query service

```go
type PageSectionQueryService interface {
    GetSections(ctx context.Context, pageID uint, query dto.GetPageSectionsQuery) (*dto.PageSectionListResponse, error)
    GetSection(ctx context.Context, pageID, sectionID uint) (*dto.PageSectionResponse, error)
    GetItems(ctx context.Context, pageID, sectionID uint, query dto.GetSectionItemsQuery) (*dto.PageSectionItemListResponse, error)
    GetPublicSections(ctx context.Context, pageID uint) ([]dto.PublicPageSectionResponse, error)
    FindAffectedPublicPageSlugs(ctx context.Context, itemType pageEntity.PageSectionItemType, itemIDs []uint) ([]string, error)
    FindMediaAffectedPublicPageSlugs(ctx context.Context, mediaIDs []uint) ([]string, error)
}
```

Query strategy:

- Admin Section list: query sections theo Page + preload/batch Background/Feature Media + một grouped query đếm `(section_id, collection, item_type)`. Số query cố định, không query collection theo từng card.
- Admin Items: query pivot theo Section/collection, xác minh một Item Type, sau đó một batch query đúng bảng nguồn; giữ thứ tự `sort_order ASC, id ASC`.
- Public: sections ACTIVE + Media trực tiếp, tất cả pivot active, rồi tối đa ba batch query Category/Post/Media. Tổng query cố định theo loại, không tăng theo số Section/Item.
- Reverse lookup chỉ join `page_section_items → page_sections → pages`, filter Page PUBLISHED/non-deleted và trả distinct slug. Media lookup thêm `background_media_id`/`feature_media_id`.

### 3.4. Command flows

#### Create Section

1. Validate/normalize key, name, title, description, color, status, optional sort order.
2. Begin transaction, lock Page active; Page missing/soft-deleted → `ErrPageNotFound`.
3. Enforce tối đa 100 active Section/Page.
4. Lock và validate Background/Feature Media mới: tồn tại, `image/*`; Admin dùng mọi Media, Staff chỉ Media của mình.
5. Nếu `sortOrder` nil, append `MAX(sort_order)+1`; nếu insert giữa danh sách, shift Section phía sau trong cùng transaction.
6. Create/restore Section; mark Media mới `attached`; commit.
7. Sau commit, xóa exact Public Page cache và query lại response.

#### Update Section

1. Lock Page rồi Section thuộc đúng Page.
2. `key` chỉ đổi khi request gửi giá trị mới; vẫn validate regex/unique. Frontend chịu trách nhiệm confirmation UX.
3. Staff chỉ cần ownership check đối với Media ID mới/thay đổi; Media đang giữ nguyên không chặn việc sửa metadata.
4. Update full metadata, không update `sort_order` qua endpoint này; reorder có endpoint riêng.
5. Commit → invalidate public cache → trả response.

#### Delete Section

- Lock Page/Section, soft-delete toàn bộ Item active rồi soft-delete Section trong một transaction.
- Chuẩn hóa lại `sort_order` các Section còn lại về `0..n-1`.
- Không xóa Category/Post/Media và không chuyển Media về `temporary`.
- Success `204`; invalidate public cache sau commit.

#### Reorder Section

- Lock Page và toàn bộ active Section.
- Payload phải là hoán vị chính xác của tất cả Section active: đủ ID, không ID lạ/trùng, `sortOrder` liên tục `0..n-1`.
- Sai/stale snapshot trả `409 SECTION_ORDER_CONFLICT`, không update một phần.
- Bulk update atomic rồi invalidate public cache.

#### Sync Collection

1. Validate collection regex, Item Type, tối đa 100 Item, ID unique và order liên tục `0..n-1`.
2. Lock Page → Section → collection links → source rows.
3. Nếu collection active đã có Item Type khác, trả `409 COLLECTION_TYPE_CONFLICT`. Muốn đổi type phải sync `[]` bằng type hiện tại trước.
4. Batch-validate toàn bộ nguồn trước destructive mutation:
   - `CATEGORY`: row `categories` tồn tại, chưa soft-delete.
   - `POST`: row `posts` tồn tại, chưa soft-delete.
   - `MEDIA`: row `media` tồn tại, chưa soft-delete, MIME `image/*`; Staff chỉ bị ownership check với Media mới thêm, không chặn reorder/remove Media đã link từ trước.
5. Soft-delete Item không thuộc snapshot, restore/create Item còn lại, mark Media mới `attached`.
6. Commit → invalidate public cache → query lại collection đã hydrate.

Sync `[]` là hợp lệ. Nếu collection chưa tồn tại thì là idempotent no-op; nếu tồn tại, `itemType` phải khớp type hiện tại.

---

## 4. API contract và Auth Context

### 4.1. DTO binding chung

```go
type PageSectionsURI struct {
    PageID uint `uri:"page_id" binding:"required,min=1"`
}

type PageSectionURI struct {
    PageID    uint `uri:"page_id" binding:"required,min=1"`
    SectionID uint `uri:"section_id" binding:"required,min=1"`
}

type SectionCollectionURI struct {
    PageID     uint   `uri:"page_id" binding:"required,min=1"`
    SectionID  uint   `uri:"section_id" binding:"required,min=1"`
    Collection string `uri:"collection" binding:"required,max=100"`
}

type GetPageSectionsQuery struct {
    Current  int    `form:"current" binding:"omitempty,min=1"`
    PageSize int    `form:"pageSize" binding:"omitempty,min=1,max=100"`
    Status   string `form:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
}

type GetSectionItemsQuery struct {
    Collection string `form:"collection" binding:"required,max=100"`
    Current    int    `form:"current" binding:"omitempty,min=1"`
    PageSize   int    `form:"pageSize" binding:"omitempty,min=1,max=100"`
}
```

Collection/key regex `^[a-z0-9]+(?:_[a-z0-9]+)*$` và color regex `^#[0-9A-Fa-f]{6}([0-9A-Fa-f]{2})?$` được validate tại service sau trim.

### 4.2. Mutation DTO

```go
type CreatePageSectionRequest struct {
    Key                 string `json:"key" binding:"required,max=100"`
    Name                string `json:"name" binding:"required,max=255"`
    Title               string `json:"title" binding:"omitempty,max=255"`
    Description         string `json:"description" binding:"omitempty,max=5000"`
    BackgroundColor     *string `json:"backgroundColor" binding:"omitempty,max=9"`
    BackgroundMediaID   *uint   `json:"backgroundMediaId" binding:"omitempty,min=1"`
    FeatureMediaID      *uint   `json:"featureMediaId" binding:"omitempty,min=1"`
    SortOrder           *int    `json:"sortOrder" binding:"omitempty,gte=0"`
    Status              string  `json:"status" binding:"required,oneof=ACTIVE INACTIVE"`
}

type UpdatePageSectionRequest struct {
    Key                 string  `json:"key" binding:"required,max=100"`
    Name                string  `json:"name" binding:"required,max=255"`
    Title               string  `json:"title" binding:"omitempty,max=255"`
    Description         string  `json:"description" binding:"omitempty,max=5000"`
    BackgroundColor     *string `json:"backgroundColor" binding:"omitempty,max=9"`
    BackgroundMediaID   *uint   `json:"backgroundMediaId" binding:"omitempty,min=1"`
    FeatureMediaID      *uint   `json:"featureMediaId" binding:"omitempty,min=1"`
    Status              string  `json:"status" binding:"required,oneof=ACTIVE INACTIVE"`
}

type ReorderSectionItem struct {
    ID        uint `json:"id" binding:"required,min=1"`
    SortOrder int  `json:"sortOrder" binding:"gte=0"`
}

type ReorderPageSectionsRequest struct {
    Sections []ReorderSectionItem `json:"sections" binding:"required,min=1,max=100,dive"`
}

type SyncSectionItem struct {
    ItemID    uint `json:"itemId" binding:"required,min=1"`
    SortOrder int  `json:"sortOrder" binding:"gte=0"`
}

type SyncSectionCollectionRequest struct {
    ItemType string            `json:"itemType" binding:"required,oneof=CATEGORY POST MEDIA"`
    Items    []SyncSectionItem `json:"items" binding:"max=100,dive"`
}
```

PUT Section là full replacement cho metadata. JSON `null` ở Media/color nghĩa là clear; `sortOrder` chỉ đi qua reorder endpoint.

### 4.3. Response DTO

```go
type SectionMediaResponse struct {
    ID           uint   `json:"id"`
    FileName     string `json:"fileName"`
    OriginalURL  string `json:"originalUrl"`
    ThumbnailURL string `json:"thumbnailUrl,omitempty"`
    MediumURL    string `json:"mediumUrl,omitempty"`
    MimeType     string `json:"mimeType"`
}

type SectionCollectionSummary struct {
    Collection string `json:"collection"`
    ItemType   string `json:"itemType"`
    Total      int64  `json:"total"`
}

type PageSectionResponse struct {
    ID                uint                       `json:"id"`
    PageID            uint                       `json:"pageId"`
    Key               string                     `json:"key"`
    Name              string                     `json:"name"`
    Title             string                     `json:"title"`
    Description       string                     `json:"description"`
    BackgroundColor   *string                    `json:"backgroundColor"`
    BackgroundMediaID *uint                      `json:"backgroundMediaId"`
    BackgroundMedia   *SectionMediaResponse      `json:"backgroundMedia"`
    FeatureMediaID    *uint                      `json:"featureMediaId"`
    FeatureMedia      *SectionMediaResponse      `json:"featureMedia"`
    SortOrder         int                        `json:"sortOrder"`
    Status            string                     `json:"status"`
    Collections       []SectionCollectionSummary `json:"collections"`
    CreatedAt         time.Time                  `json:"createdAt"`
    UpdatedAt         time.Time                  `json:"updatedAt"`
}

type PageSectionListResponse struct {
    Data  []PageSectionResponse `json:"data"`
    Total int64                 `json:"total"`
}
```

Hydrated Item dùng một DTO typed, không dùng `any`/`interface{}`:

```go
type SectionItemDataResponse struct {
    ID           uint   `json:"id"`
    Name         string `json:"name,omitempty"`
    Title        string `json:"title,omitempty"`
    Slug         string `json:"slug,omitempty"`
    TypeCode     string `json:"typeCode,omitempty"`
    ImageURL     string `json:"imageUrl,omitempty"`
    Status       string `json:"status,omitempty"`
    FileName     string `json:"fileName,omitempty"`
    OriginalURL  string `json:"originalUrl,omitempty"`
    ThumbnailURL string `json:"thumbnailUrl,omitempty"`
    MediumURL    string `json:"mediumUrl,omitempty"`
    MimeType     string `json:"mimeType,omitempty"`
}

type PageSectionItemResponse struct {
    ID         uint                    `json:"id"`
    SectionID  uint                    `json:"sectionId"`
    ItemType   string                  `json:"itemType"`
    ItemID     uint                    `json:"itemId"`
    Collection string                  `json:"collection"`
    SortOrder  int                     `json:"sortOrder"`
    Data       SectionItemDataResponse `json:"data"`
}

type PageSectionItemListResponse struct {
    Data  []PageSectionItemResponse `json:"data"`
    Total int64                     `json:"total"`
}

type ErrorResponse struct {
    Error   string            `json:"error"`
    Code    string            `json:"code"`
    Message string            `json:"message,omitempty"`
    Fields  map[string]string `json:"fields,omitempty"`
}
```

- Collections/Data luôn là `[]` khi rỗng, không trả `null`.
- Không trả `author_id`, `owner_id`, filesystem path, token/JTI, password, SQL error hoặc GORM entity nội bộ.

### 4.4. Endpoints

#### GET `/api/v1/admin/pages/:page_id/sections`

- Query: `current`, `pageSize`, optional `status`.
- Mặc định `current=1`, `pageSize=100`; order cố định `sort_order ASC, id ASC`.
- `200 PageSectionListResponse`; Page không tồn tại/soft-delete → `404 PAGE_NOT_FOUND`.

#### POST `/api/v1/admin/pages/:page_id/sections`

- Body `CreatePageSectionRequest`.
- `201 { "data": PageSectionResponse }`.
- Duplicate active/key restore race → `409 SECTION_KEY_CONFLICT`.

#### GET `/api/v1/admin/pages/:page_id/sections/:section_id`

- `200 { "data": PageSectionResponse }`.
- Section không thuộc Page hoặc soft-deleted → `404 SECTION_NOT_FOUND`; không tiết lộ Section thuộc Page khác.

#### PUT `/api/v1/admin/pages/:page_id/sections/:section_id`

- Body `UpdatePageSectionRequest`.
- `200 { "data": PageSectionResponse, "message": "Cập nhật Section thành công" }`.

#### DELETE `/api/v1/admin/pages/:page_id/sections/:section_id`

- Soft-delete Section + Item, normalize order, không xóa source.
- Success `204 No Content`.

#### PUT `/api/v1/admin/pages/:page_id/section-order`

- Body `ReorderPageSectionsRequest` chứa toàn bộ active Section.
- Success `200 { "message": "Sắp xếp Section thành công", "data": [] }` theo thứ tự mới.

#### GET `/api/v1/admin/pages/:page_id/sections/:section_id/items`

- Query bắt buộc `collection`; optional `current`, `pageSize`.
- `200 PageSectionItemListResponse`, order `sort_order ASC, id ASC`.
- Collection chưa có → `{ "data": [], "total": 0 }`.

#### PUT `/api/v1/admin/pages/:page_id/sections/:section_id/items/:collection`

- Body `SyncSectionCollectionRequest`; đây là snapshot cuối cùng.
- Success `200 { "message": "Đồng bộ collection thành công", "data": [] }` đã hydrate/sắp xếp.

### 4.5. Route và JWT Context

```go
pageSections := admin.Group("/pages/:page_id")
pageSections.Use(middleware.RoleMiddleware("ADMIN", "STAFF"))
{
    pageSections.GET("/sections", pageSectionController.GetSections)
    pageSections.POST("/sections", pageSectionController.CreateSection)
    pageSections.GET("/sections/:section_id", pageSectionController.GetSection)
    pageSections.PUT("/sections/:section_id", pageSectionController.UpdateSection)
    pageSections.DELETE("/sections/:section_id", pageSectionController.DeleteSection)
    pageSections.PUT("/section-order", pageSectionController.ReorderSections)
    pageSections.GET("/sections/:section_id/items", pageSectionController.GetItems)
    pageSections.PUT("/sections/:section_id/items/:collection", pageSectionController.SyncItems)
}
```

- `admin` đã chạy `AuthMiddleware(redisClient)`: verify chữ ký/expiry trước Redis blacklist lookup, rồi set `userId` và `role`.
- Controller type-assert `userId uint`, `role string`; không nhận identity/owner từ body/query.
- ADMIN quản lý mọi Section/Media; STAFF theo quyền Page hiện hữu và chỉ gắn Media mới do mình sở hữu. CUSTOMER bị `403` trước controller.

### 4.6. Error mapping

| Điều kiện | HTTP | Code |
| --- | ---: | --- |
| URI/query/body/key/color/collection/order sai | 400 | `VALIDATION_ERROR` |
| JWT thiếu/sai/hết hạn/blacklisted | 401 | `UNAUTHORIZED` |
| Role sai hoặc Staff gắn Media người khác | 403 | `FORBIDDEN` |
| Page không tồn tại | 404 | `PAGE_NOT_FOUND` |
| Section không thuộc Page/đã xóa | 404 | `SECTION_NOT_FOUND` |
| Category/Post/Media nguồn không tồn tại | 404 | `SOURCE_NOT_FOUND` |
| Key trùng | 409 | `SECTION_KEY_CONFLICT` |
| Snapshot reorder stale | 409 | `SECTION_ORDER_CONFLICT` |
| Collection đang thuộc Item Type khác | 409 | `COLLECTION_TYPE_CONFLICT` |
| Duplicate/race link | 409 | `SECTION_ITEM_CONFLICT` |
| MySQL/Redis/config unexpected | 500 | `INTERNAL_ERROR` |

Redis invalidation lỗi sau commit chỉ log, không biến mutation DB thành `500` nếu dữ liệu đã commit thành công.

---

## 5. Media Library cho picker

Frontend cần chọn Background/Feature/Media có sẵn, nhưng Backend hiện chưa có list Media. Mở rộng module Media hiện hữu:

```text
GET /api/v1/admin/media?current=1&pageSize=20&q={search}&mimeType=image&status={optional}
```

```go
type GetMediaListQuery struct {
    Current  int    `form:"current" binding:"omitempty,min=1"`
    PageSize int    `form:"pageSize" binding:"omitempty,min=1,max=100"`
    Q        string `form:"q" binding:"omitempty,max=255"`
    MimeType string `form:"mimeType" binding:"omitempty,oneof=image"`
    Status   string `form:"status" binding:"omitempty,oneof=temporary attached"`
}

type MediaListItemResponse struct {
    ID           uint   `json:"id"`
    FileName     string `json:"fileName"`
    OriginalURL  string `json:"originalUrl"`
    ThumbnailURL string `json:"thumbnailUrl,omitempty"`
    MediumURL    string `json:"mediumUrl,omitempty"`
    MimeType     string `json:"mimeType"`
    Size         int64  `json:"size"`
    Status       string `json:"status"`
    CreatedAt    time.Time `json:"createdAt"`
}
```

- Mở rộng `MediaRepository.FindAndCount`, `MediaService.ListMedia`, `MediaController.ListMedia`; không tạo module Media mới.
- Route dùng `AuthMiddleware` + `RoleMiddleware("ADMIN", "STAFF")`.
- ADMIN thấy mọi Media hợp lệ; STAFF chỉ thấy `owner_id = userId`. Response không trả `ownerId` hoặc filesystem path.
- Search escape `%`, `_`, `\`; order cố định `created_at DESC, id DESC`; `mimeType=image` map server-side thành `mime_type LIKE 'image/%'`.
- Không cache search/list Redis. Frontend chịu trách nhiệm debounce 300 ms.
- Category/Post picker tái sử dụng `GET /admin/categories` và `GET /admin/posts`; không tạo endpoint bản sao.

---

## 6. Mở rộng Public Page read model

Sau khi plan `26-page-public-api-plan.md` được thực thi, bổ sung `sections` vào một response `GET /api/v1/pages/:slug`; không tạo public endpoint riêng theo Section.

```go
type PublicSectionItemDataResponse struct {
    ID           uint   `json:"id"`
    Name         string `json:"name,omitempty"`
    Title        string `json:"title,omitempty"`
    Slug         string `json:"slug,omitempty"`
    TypeCode     string `json:"typeCode,omitempty"`
    ImageURL     string `json:"imageUrl,omitempty"`
    OriginalURL  string `json:"originalUrl,omitempty"`
    ThumbnailURL string `json:"thumbnailUrl,omitempty"`
    MediumURL    string `json:"mediumUrl,omitempty"`
}

type PublicPageSectionItemResponse struct {
    ItemType string                        `json:"itemType"`
    ItemID   uint                          `json:"itemId"`
    SortOrder int                          `json:"sortOrder"`
    Data     PublicSectionItemDataResponse `json:"data"`
}

type PublicPageSectionResponse struct {
    ID              uint                                           `json:"id"`
    Key             string                                         `json:"key"`
    Title           string                                         `json:"title"`
    Description     string                                         `json:"description"`
    BackgroundColor *string                                        `json:"backgroundColor"`
    BackgroundMedia *SectionMediaResponse                          `json:"backgroundMedia"`
    FeatureMedia    *SectionMediaResponse                          `json:"featureMedia"`
    SortOrder       int                                            `json:"sortOrder"`
    Collections     map[string][]PublicPageSectionItemResponse     `json:"collections"`
}
```

- Không trả `name`, Section status, timestamps hoặc collection metadata Admin.
- Chỉ Section `ACTIVE`, order `sort_order ASC, id ASC`.
- Category public chỉ hydrate row `ACTIVE`; Post phải chưa soft-delete; Media phải chưa soft-delete, `attached`, `image/*`.
- Source không còn public-valid bị bỏ khỏi response nhưng pivot vẫn giữ cho Dashboard xử lý/khôi phục; không tự trỏ sang Item khác.
- `collections` luôn `{}` khi không có Item; mỗi mảng giữ order pivot.
- Public Page query dùng số batch cố định như mục 3.3, không N+1.

---

## 7. Trụ cột 3 — Redis cache và token state

### 7.1. Public cache

```text
public:pages:slug:{slug}
```

- Cache-aside toàn bộ Public Page response đã gồm Media, SEO và Section; TTL 15 phút bằng `SetEx` theo plan Public Page.
- Không cache GORM entity, raw ownership, Admin metadata hoặc token.
- Section GET Admin và Media Library list đọc trực tiếp MySQL; index/pagination đủ cho v1.
- Mỗi mutation Section/Item lấy Page slug từ row đã lock; chỉ `DEL public:pages:slug:{slug}` sau transaction commit.
- Redis read lỗi ở Public API fallback MySQL; Redis write/delete lỗi log và không làm hỏng response/mutation MySQL.
- Không dùng Redis `KEYS`.

### 7.2. Invalidation khi dữ liệu nguồn thay đổi

Public Page cache chứa snapshot Category/Post/Media đã hydrate nên thay đổi nguồn cũng phải invalidate:

- Category update/status/delete/bulk-status: reverse lookup `item_type=CATEGORY`, lấy distinct Page slug rồi pipeline `DEL` exact keys sau commit.
- Post update/delete: reverse lookup `item_type=POST` và invalidate các Page liên quan.
- Media delete/metadata mutation: lookup `item_type=MEDIA` cộng `background_media_id`/`feature_media_id`; thu thập slug trước hard delete nếu FK sẽ `SET NULL`.
- Upload Media chưa liên kết không invalidate Page.
- Reverse invalidation failure chỉ log; TTL 15 phút là safety net, không rollback source mutation đã commit.

Tạo helper dùng chung, ví dụ `PublicPageCacheInvalidator`, nhận `*redis.Client` và query service; Category/Post/Media services gọi helper thay vì tự viết Redis key logic.

### 7.3. JWT blacklist, refresh rotation và logout

Page Section không tạo token key mới. Tất cả Admin route kế thừa các pattern Auth:

```text
auth:bl:at:{jti}
auth:rf:{user_id}:{jti}
auth:rf:index:{user_id}
```

- Auth middleware phải lấy `JWT_SECRET` từ environment/config, verify algorithm/signature/expiry trước khi gọi Redis, dùng request context, rồi check `auth:bl:at:{jti}`.
- Logout chỉ chạy sau AuthMiddleware; xóa refresh session/JTI, `SREM` index và `SetEx auth:bl:at:{jti}` với TTL chính xác `exp - now`. TTL <= 0 thì không ghi key rác.
- Refresh rotation phải atomic bằng Redis Lua hoặc transaction tương đương:
  1. Kiểm tra old refresh key tồn tại.
  2. Xóa old key + `SREM` old JTI.
  3. Tạo key new JTI + `SADD` index + đặt TTL.
  4. Chỉ trả token mới khi script thành công.
- Old JTI không tồn tại được coi là replay: dùng index theo user để thu hồi toàn bộ refresh session liên quan và buộc đăng nhập lại; không lưu full token trong Redis.
- Auth implementation hiện tại rotate theo nhiều lệnh rời và hard-code secret; đây là shared security debt phải được harden trước production, không nhân bản logic trong Page Section service.

---

## 8. Wire, route, OpenAPI và seeder

### 8.1. Google Wire

```go
var PageSectionSet = wire.NewSet(
    repository.NewPageSectionRepository,
    infrastructureQuery.NewPageSectionQueryService,
    pageSectionService.NewPageSectionService,
    controller.NewPageSectionController,
)

func InitializePageSectionController(
    db *gorm.DB,
    redisClient *redis.Client,
) *controller.PageSectionController
```

- Controller có command/query service qua constructor; `*gorm.DB`/Redis không đi trực tiếp vào Handler.
- Regenerate `internal/di/wire_gen.go`; không sửa generated code thủ công.
- Khởi tạo controller trong `main.go`, đăng ký route sau `pages` base và trước start server.
- Mở rộng `MediaSet` chỉ bằng methods trên repository/service/controller hiện tại; constructor signature không cần đổi nếu list không dùng Redis.

### 8.2. OpenAPI

Cập nhật `.docs/api-endpoints.yaml`:

- Tám endpoint Page Section/Item.
- Một endpoint Admin Media Library.
- Public Page response có `sections`.
- BearerAuth cho toàn bộ Admin routes; public route không BearerAuth.
- Schemas request/response/error và status `200/201/204/400/401/403/404/409/500`.

### 8.3. Seeder

- Bổ sung seeder idempotent `page_sections.go` sau seeder Page/Category/Post/Media.
- Lookup Page bằng slug và nguồn bằng stable slug/file name; không hard-code auto-increment ID.
- Upsert Section theo `(page_id, key)`, restore nếu soft-delete; sync Item theo collection và order.
- Seeder không xóa Section/Item ngoài fixture của nó và phải chạy lại không tạo duplicate.

---

## 9. Kiểm thử và acceptance criteria

### 9.1. Database/repository

- Migration up/down đúng thứ tự, FK/index/check constraints khớp kiểu `INT UNSIGNED` của bảng nguồn.
- Unique Page/key và Section/collection/type/item hoạt động; re-create/re-add restore row soft-deleted.
- Hard-delete Page cascade Section/Item; hard-delete Media hợp lệ set null Background/Feature.
- Polymorphic `item_id` không có FK nhưng Use Case từ chối source thiếu/sai loại.
- Query Section summary và hydrate Item có số query cố định, order ổn định, không N+1.

### 9.2. Service/transaction

- Create/Update từ chối key/color/status/MIME/ownership sai mà không ghi một phần.
- Concurrent create/reorder/sync cùng Page được tuần tự hóa bởi Page lock; không duplicate/order gap.
- Reorder chỉ chấp nhận exact snapshot, stale payload không update một phần.
- Sync validate toàn bộ trước soft-delete; duplicate, gap, missing source, type conflict hoặc Media forbidden giữ nguyên collection cũ.
- Sync `[]`, delete Section và remove Item không xóa dữ liệu nguồn/physical Media.
- Staff giữ được Media cũ do Admin gắn khi sửa metadata/reorder, nhưng không thêm Media mới của owner khác.

### 9.3. API/Auth

- List responses luôn `{ data: [], total: 0 }` khi rỗng; nested Section ID của Page khác trả `404`.
- Controller mỏng, bind đủ URI/query/body và không nhận userId/role/owner từ client.
- ADMIN/STAFF đúng quyền; CUSTOMER `403`; JWT invalid/expired/revoked `401` trước service.
- Error response không lộ SQL, stack, path file, owner, claims/token hay dữ liệu nhạy cảm.
- Media Library phân trang/search/filter ảnh đúng; Staff chỉ thấy Media của mình.

### 9.4. Public/cache

- Public Page chỉ trả Section ACTIVE, source public-valid, collection/order đúng và không trả `name` Admin.
- Cache hit không query MySQL; Section/Item mutation xóa exact Page key sau commit.
- Category/Post/Media source mutation invalidate mọi Public Page tham chiếu qua reverse lookup.
- Redis unavailable không làm mutation MySQL rollback hoặc Public GET thất bại nếu DB còn hoạt động.

### 9.5. Token security

- Middleware verify token trước Redis blacklist lookup và dùng JTI key.
- Logout blacklist access JTI đúng TTL còn lại.
- Refresh rotation chỉ một request dùng được old JTI; concurrent/replay thu hồi session theo policy, không tồn tại khoảng trống old/new key.

---

## 10. Trình tự triển khai

1. Tạo migration + entity `PageSection`, `PageSectionItem`.
2. Viết command repository, query projections/service và test batch query/locking.
3. Viết DTO, command service, domain errors và transaction flows.
4. Viết thin controller, `PageSectionSet`, Wire injector và routes.
5. Mở rộng Media Library list/search với ownership/RBAC.
6. Cập nhật Public Page read model/response và exact cache invalidation, gồm reverse invalidation từ Category/Post/Media.
7. Cập nhật OpenAPI và seeder idempotent.
8. Chạy migration, unit/integration/API tests, regenerate Wire và kiểm tra race/concurrent snapshot.
9. Tích hợp Frontend theo `.docs/frontend-plans/dashboard/27-page-section-management-plan.md`.
10. Harden shared JWT secret/atomic refresh rotation trước production nếu chưa được xử lý ở module Auth.
