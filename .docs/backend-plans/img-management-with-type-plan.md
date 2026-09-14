# KẾ HOẠCH BACKEND: QUẢN LÝ NỘI DUNG HÌNH ẢNH THEO TYPE

**Tài liệu nguồn:** `.docs/ideas/img-management-with-type.md`

**Module:** `ImageContent`

**Stack bắt buộc:** Go, Gin, GORM, Google Wire, MySQL, Redis

**Phạm vi:** Database, Admin API cho Refine, Public API và cache. Không bao gồm code UI Admin hoặc tích hợp render vào `apps/webview`.

## 1. Hiện trạng và quyết định kiến trúc

- `media` đã quản lý file, URL biến thể, MIME, owner và trạng thái `temporary/attached`; mọi ảnh của module mới bắt buộc chọn hoặc upload qua API Media hiện tại.
- `post_types/posts` không phù hợp để tái sử dụng trực tiếp vì kéo theo `slug`, category, content và SEO. `page_section_items` cũng không phù hợp vì dữ liệu Logo/Slider/Partner là toàn site và pivot hiện không có metadata theo ngữ cảnh.
- Tạo domain riêng gồm `image_content_types` và `image_contents`, theo mô hình Type + Item. Type quyết định các field được bật/bắt buộc và giới hạn số item; Item lưu metadata hiển thị gắn với một Media.
- Phạm vi mở rộng động chỉ áp dụng cho tập field chuẩn đã hỗ trợ: `name`, `description`, `secondaryDescription`, `url`. Thêm một Type dùng tổ hợp các field này không cần đổi schema/backend. Thêm loại field hoàn toàn mới vẫn cần migration và code.
- Image, trạng thái và thứ tự là field hệ thống của mọi Item, không nằm trong `field_config`.
- Dùng CQRS thực dụng:
  - Command service xử lý create/update/delete/reorder trong transaction.
  - Query service dựng Admin read model và Public projection, không trả GORM entity trực tiếp.
- Admin Type chỉ dành cho `ADMIN`; Admin Item dành cho `ADMIN` và `STAFF`. Public endpoint không yêu cầu token.
- Redis chỉ cache Public projection. Admin list/detail luôn đọc MySQL để tránh dữ liệu cấu hình cũ.

---

## 2. Trụ cột 1 — Thiết kế dữ liệu

### 2.1. Value objects dùng chung

Tạo package `apps/backend/internal/domain/imagecontent/entity`:

```go
package entity

type ImageContentStatus string

const (
    ImageContentStatusActive   ImageContentStatus = "ACTIVE"
    ImageContentStatusInactive ImageContentStatus = "INACTIVE"
)

type ImageContentFieldRule struct {
    Enabled  bool `json:"enabled"`
    Required bool `json:"required"`
}

type ImageContentFieldConfig struct {
    Name                 ImageContentFieldRule `json:"name"`
    Description          ImageContentFieldRule `json:"description"`
    SecondaryDescription ImageContentFieldRule `json:"secondaryDescription"`
    URL                  ImageContentFieldRule `json:"url"`
}
```

Quy tắc bất biến:

- `Required=true` chỉ hợp lệ khi `Enabled=true`.
- JSON chỉ nhận bốn key trên; DTO dùng struct typed để không chấp nhận field tùy ý.
- Type code được trim, uppercase và phải khớp `^[A-Z][A-Z0-9_]{0,49}$`.
- `max_items = NULL` nghĩa là không giới hạn; giá trị khác null nằm trong `1..1000`.

### 2.2. GORM entity `ImageContentType`

Tạo `apps/backend/internal/domain/imagecontent/entity/image_content_type.go`:

```go
package entity

import (
    "time"

    "gorm.io/gorm"
)

type ImageContentType struct {
    ID          uint                    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
    Code        string                  `gorm:"column:code;type:varchar(50);not null;uniqueIndex:uq_image_content_types_code" json:"code"`
    Name        string                  `gorm:"column:name;type:varchar(100);not null" json:"name"`
    Status      ImageContentStatus      `gorm:"column:status;type:varchar(20);not null;default:'ACTIVE';index:idx_image_content_types_status_sort,priority:1" json:"status"`
    SortOrder   int                     `gorm:"column:sort_order;not null;default:0;index:idx_image_content_types_status_sort,priority:2" json:"sortOrder"`
    MaxItems    *uint                   `gorm:"column:max_items;type:int unsigned" json:"maxItems"`
    FieldConfig ImageContentFieldConfig `gorm:"column:field_config;type:json;serializer:json;not null" json:"-"`

    Items     []ImageContent `gorm:"foreignKey:TypeCode;references:Code" json:"-"`
    CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
    UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
    DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index:idx_image_content_types_deleted_at" json:"-"`
}

func (ImageContentType) TableName() string { return "image_content_types" }
```

- `code` là public identifier ổn định và không được đổi qua API Update.
- Create với code đang active trả conflict. Nếu code chỉ tồn tại ở row soft-deleted, repository restore row đó và cập nhật cấu hình thay vì insert trùng unique key.
- `Items` chỉ phục vụ khai báo quan hệ; query list/detail dùng projection/batch query rõ ràng.

### 2.3. GORM entity `ImageContent`

Tạo `apps/backend/internal/domain/imagecontent/entity/image_content.go`:

```go
package entity

import (
    "time"

    mediaEntity "go_refine_dashboard_be/internal/domain/media/entity"
    "gorm.io/gorm"
)

type ImageContent struct {
    ID       uint   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
    TypeCode string `gorm:"column:type_code;type:varchar(50);not null;index:idx_image_contents_type_status_sort,priority:1" json:"typeCode"`
    MediaID  uint   `gorm:"column:media_id;not null;index:idx_image_contents_media_id" json:"mediaId"`

    Name                 *string `gorm:"column:name;type:varchar(255)" json:"name"`
    Description          *string `gorm:"column:description;type:text" json:"description"`
    SecondaryDescription *string `gorm:"column:secondary_description;type:text" json:"secondaryDescription"`
    TargetURL            *string `gorm:"column:target_url;type:varchar(500)" json:"url"`

    SortOrder int                `gorm:"column:sort_order;not null;default:0;index:idx_image_contents_type_status_sort,priority:3" json:"sortOrder"`
    Status    ImageContentStatus `gorm:"column:status;type:varchar(20);not null;default:'ACTIVE';index:idx_image_contents_type_status_sort,priority:2" json:"status"`
    CreatedBy uint               `gorm:"column:created_by;type:int;not null;index:idx_image_contents_created_by" json:"createdBy"`

    Type      *ImageContentType  `gorm:"foreignKey:TypeCode;references:Code;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
    Media     *mediaEntity.Media `gorm:"foreignKey:MediaID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
    CreatedAt time.Time          `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
    UpdatedAt time.Time          `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
    DeletedAt gorm.DeletedAt     `gorm:"column:deleted_at;index:idx_image_contents_type_status_sort,priority:4;index:idx_image_contents_deleted_at" json:"-"`
}

func (ImageContent) TableName() string { return "image_contents" }
```

- Metadata nullable vì mỗi Type bật một tập field khác nhau. Service trim chuỗi và chuyển chuỗi rỗng thành `nil`.
- Một Media được phép xuất hiện ở nhiều Item/Type; không tạo unique constraint trên `media_id`.
- `created_by` dùng `INT` signed để khớp migration `users.id`; Go vẫn dùng `uint` vì identity trong Gin Context hiện là `uint`.
- `ON DELETE RESTRICT` ở `media_id` là lớp bảo vệ cuối cùng ngoài `media.status=attached` của Media Service.
- Không hoàn nguyên Media về `temporary` khi Item bị xóa vì Media có thể còn được Product/Post/Page/Section hoặc Item khác sử dụng.

### 2.4. Migration MySQL

Tạo cặp migration timestamp mới, ví dụ:

```sql
CREATE TABLE `image_content_types` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `code` VARCHAR(50) NOT NULL,
  `name` VARCHAR(100) NOT NULL,
  `status` VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
  `sort_order` INT NOT NULL DEFAULT 0,
  `max_items` INT UNSIGNED NULL,
  `field_config` JSON NOT NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_image_content_types_code` (`code`),
  KEY `idx_image_content_types_status_sort` (`status`, `sort_order`, `deleted_at`),
  KEY `idx_image_content_types_deleted_at` (`deleted_at`),
  CONSTRAINT `chk_image_content_types_status`
    CHECK (`status` IN ('ACTIVE', 'INACTIVE')),
  CONSTRAINT `chk_image_content_types_sort_order`
    CHECK (`sort_order` >= 0),
  CONSTRAINT `chk_image_content_types_max_items`
    CHECK (`max_items` IS NULL OR (`max_items` >= 1 AND `max_items` <= 1000))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `image_contents` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `type_code` VARCHAR(50) NOT NULL,
  `media_id` INT UNSIGNED NOT NULL,
  `name` VARCHAR(255) NULL,
  `description` TEXT NULL,
  `secondary_description` TEXT NULL,
  `target_url` VARCHAR(500) NULL,
  `sort_order` INT NOT NULL DEFAULT 0,
  `status` VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
  `created_by` INT NOT NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL,
  PRIMARY KEY (`id`),
  KEY `idx_image_contents_type_status_sort`
    (`type_code`, `status`, `sort_order`, `deleted_at`),
  KEY `idx_image_contents_media_id` (`media_id`),
  KEY `idx_image_contents_created_by` (`created_by`),
  KEY `idx_image_contents_deleted_at` (`deleted_at`),
  CONSTRAINT `chk_image_contents_status`
    CHECK (`status` IN ('ACTIVE', 'INACTIVE')),
  CONSTRAINT `chk_image_contents_sort_order`
    CHECK (`sort_order` >= 0),
  CONSTRAINT `fk_image_contents_type`
    FOREIGN KEY (`type_code`) REFERENCES `image_content_types` (`code`)
    ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT `fk_image_contents_media`
    FOREIGN KEY (`media_id`) REFERENCES `media` (`id`)
    ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT `fk_image_contents_created_by`
    FOREIGN KEY (`created_by`) REFERENCES `users` (`id`)
    ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

Down migration drop `image_contents` trước, sau đó `image_content_types`. Không sửa schema `media`, `posts`, `pages` hoặc `page_sections`.

---

## 3. Tổ chức DDD/CQRS và repository contract

### 3.1. Cấu trúc module

```text
apps/backend/internal
├── controller
│   └── image_content_controller.go
├── domain/imagecontent/entity
│   ├── image_content.go
│   └── image_content_type.go
├── infrastructure
│   ├── queryservice/image_content_query_service.go
│   └── repository/image_content_repository.go
└── usecases/imagecontent
    ├── dto/image_content_dto.go
    └── service/image_content_service.go
```

- Controller chỉ bind URI/query/body, lấy identity từ Gin Context, gọi service/query service và map error.
- Command service giữ validation, transaction, lock order và cache invalidation.
- Repository chỉ persistence/locking; không quyết định HTTP, RBAC hoặc field requirement.
- Query service dùng projection typed và preload/batch Media để tránh N+1.

### 3.2. Repository contract

```go
type ImageContentOrderUpdate struct {
    ID        uint
    SortOrder int
}

type ImageContentRepository interface {
    FindTypeForUpdate(ctx context.Context, tx *gorm.DB, id uint) (*entity.ImageContentType, error)
    FindTypeByCodeForUpdate(ctx context.Context, tx *gorm.DB, code string) (*entity.ImageContentType, error)
    FindContentForUpdate(ctx context.Context, tx *gorm.DB, id uint) (*entity.ImageContent, error)
    FindContentsByTypeForUpdate(ctx context.Context, tx *gorm.DB, typeCode string) ([]entity.ImageContent, error)
    FindMediaForUpdate(ctx context.Context, tx *gorm.DB, mediaID uint) (*mediaEntity.Media, error)

    CreateOrRestoreType(ctx context.Context, tx *gorm.DB, value *entity.ImageContentType) error
    UpdateType(ctx context.Context, tx *gorm.DB, value *entity.ImageContentType) error
    SoftDeleteType(ctx context.Context, tx *gorm.DB, id uint) error
    CreateContent(ctx context.Context, tx *gorm.DB, value *entity.ImageContent) error
    UpdateContent(ctx context.Context, tx *gorm.DB, value *entity.ImageContent) error
    SoftDeleteContent(ctx context.Context, tx *gorm.DB, id uint) error
    BulkUpdateContentOrder(ctx context.Context, tx *gorm.DB, typeCode string, updates []ImageContentOrderUpdate) error
    MarkMediaAttached(ctx context.Context, tx *gorm.DB, mediaIDs []uint) error

    FindTypes(ctx context.Context, filter TypeListFilter) ([]entity.ImageContentType, map[string]int64, int64, error)
    FindTypeByID(ctx context.Context, id uint) (*entity.ImageContentType, int64, error)
    FindContents(ctx context.Context, filter ContentListFilter) ([]entity.ImageContent, int64, error)
    FindContentByID(ctx context.Context, id uint) (*entity.ImageContent, error)
    FindPublicProjection(ctx context.Context) ([]PublicImageContentTypeRow, error)
}
```

Query rules:

- Sort field từ client phải map qua whitelist; không nối trực tiếp `_sort` vào SQL.
- Search escape `%`, `_` và `\\` trước `LIKE`; Type search trên `code/name`, Item search trên `name/description/secondary_description/target_url`.
- Admin Item list preload Media bằng query cố định và giữ `sort_order ASC, id ASC` mặc định.
- Type list lấy `itemCount` bằng một grouped query theo `type_code`, không count từng row.
- Public query chỉ lấy Type/Item `ACTIVE`, Media chưa soft-delete, `status='attached'`, MIME `image/%`; trả type và item theo order ổn định.

### 3.3. Lock order và transaction

Mọi command theo thứ tự lock thống nhất để tránh deadlock:

```text
ImageContentType → ImageContent rows → Media row
```

- Create Item: lock Type → lock toàn bộ Item cùng Type để đếm/append order → lock Media → insert → mark attached.
- Update Item: lock Type của Item → lock Item → lock Media mới nếu đổi ảnh → validate/update.
- Delete Item: lock Type → lock Item và danh sách cùng Type → soft delete → normalize order.
- Reorder: lock Type → lock toàn bộ Item chưa xóa của Type → validate full snapshot → bulk update.
- Update/Delete Type: lock Type → lock toàn bộ Item của Type → kiểm tra tương thích hoặc `TYPE_IN_USE`.
- Cache chỉ invalidate sau khi transaction commit thành công. Redis lỗi không rollback mutation MySQL.

---

## 4. Quy tắc nghiệp vụ Command

### 4.1. Create/restore Type

1. Normalize code uppercase; validate code/name/status/sort/maxItems/fieldConfig.
2. `required=true && enabled=false` trả `TYPE_CONFIG_INVALID`.
3. Trong transaction, tìm code bằng `Unscoped`:
   - Row active tồn tại: `TYPE_CODE_CONFLICT`.
   - Row soft-deleted: restore cùng ID, ghi cấu hình mới, đặt `deleted_at=NULL`.
   - Không tồn tại: insert mới.
4. Commit, invalidate public cache, trả Type response.

### 4.2. Update Type

- Không nhận hoặc cập nhật `code`.
- Lock Type và toàn bộ Item chưa soft-delete.
- Nếu giảm `maxItems` dưới số Item hiện tại: `TYPE_CONFIG_CONFLICT`.
- Nếu một field chuyển thành required, mọi Item hiện tại phải có giá trị khác rỗng; nếu không trả `TYPE_CONFIG_CONFLICT` kèm field/count trong error details.
- Disable field không xóa dữ liệu đã lưu; Admin detail vẫn trả raw metadata để có thể khôi phục cấu hình, nhưng Public projection không xuất field đang disabled.
- Type `INACTIVE` không xuất hiện trong Public API; Item vẫn được phép quản trị để chuẩn bị trước khi bật Type.

### 4.3. Delete Type

- Chỉ soft-delete khi không còn Item chưa soft-delete.
- Type còn Item trả `409 TYPE_IN_USE`; không cascade hoặc xóa Item ngầm.
- Seeded Type không có cơ chế đặc quyền riêng; cùng tuân theo quy tắc trên.

### 4.4. Create Item

1. Validate Type tồn tại; Type có thể `ACTIVE` hoặc `INACTIVE`.
2. `maxItems` đếm tất cả Item chưa soft-delete, kể cả `INACTIVE`.
3. Validate metadata theo `field_config`:
   - Field disabled mà request có giá trị khác rỗng: `FIELD_NOT_ENABLED`.
   - Field required thiếu/rỗng: `FIELD_REQUIRED`.
   - `name` tối đa 255; mỗi description tối đa 5000; URL tối đa 500.
4. URL hợp lệ khi là relative path bắt đầu bằng một `/` hoặc absolute URL scheme `http/https`; từ chối protocol-relative, `javascript:`, `data:` và control character.
5. Media phải tồn tại, chưa soft-delete và MIME bắt đầu `image/`.
6. `ADMIN` được dùng mọi Media; `STAFF` chỉ được gắn Media có `owner_id=userId`.
7. Nếu không truyền `sortOrder`, append vào cuối. Nếu truyền vị trí `0..count`, shift các Item sau trong cùng transaction.
8. Insert Item, đổi Media sang `attached`, commit và invalidate cache.

### 4.5. Update Item

- `typeCode`, `createdBy` và `sortOrder` không cập nhật qua endpoint này; đổi thứ tự dùng endpoint reorder, đổi Type bằng delete/create.
- Validate full replacement metadata theo field config hiện tại.
- Staff chỉ bị ownership check khi chọn Media mới; giữ nguyên Media cũ không chặn sửa metadata.
- Khi Type đang `INACTIVE`, Item vẫn update được nhưng không public.

### 4.6. Delete và reorder Item

- Delete chỉ soft-delete Item, không xóa Media và không đổi Media về `temporary`.
- Sau delete, normalize toàn bộ Item còn lại cùng Type thành order liên tục `0..n-1`.
- Reorder request phải là hoán vị chính xác của mọi Item chưa soft-delete thuộc Type:
  - Đủ ID, không ID lạ/trùng.
  - `sortOrder` liên tục `0..n-1`.
  - Snapshot thiếu/thừa do thao tác đồng thời trả `409 ORDER_CONFLICT`, không update một phần.
- Reorder bao gồm cả Item `ACTIVE` và `INACTIVE` để một Item bật lại vẫn có vị trí xác định.

---

## 5. Trụ cột 2 — API contract và Auth Context

### 5.1. Request binding structs

```go
type ImageContentIDURI struct {
    ID uint `uri:"id" binding:"required,min=1"`
}

type ImageContentFieldRuleRequest struct {
    Enabled  bool `json:"enabled"`
    Required bool `json:"required"`
}

type ImageContentFieldConfigRequest struct {
    Name                 ImageContentFieldRuleRequest `json:"name"`
    Description          ImageContentFieldRuleRequest `json:"description"`
    SecondaryDescription ImageContentFieldRuleRequest `json:"secondaryDescription"`
    URL                  ImageContentFieldRuleRequest `json:"url"`
}

type GetImageContentTypesQuery struct {
    Start  int    `form:"_start" binding:"omitempty,min=0"`
    End    int    `form:"_end" binding:"omitempty,min=1,max=1000"`
    Search string `form:"q" binding:"omitempty,max=100"`
    Status string `form:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
    SortBy string `form:"_sort" binding:"omitempty,oneof=id code name status sortOrder createdAt updatedAt"`
    Order  string `form:"_order" binding:"omitempty,oneof=ASC DESC asc desc"`
}

type CreateImageContentTypeRequest struct {
    Code        string                          `json:"code" binding:"required,max=50"`
    Name        string                          `json:"name" binding:"required,max=100"`
    Status      string                          `json:"status" binding:"required,oneof=ACTIVE INACTIVE"`
    SortOrder   int                             `json:"sortOrder" binding:"gte=0"`
    MaxItems    *uint                           `json:"maxItems" binding:"omitempty,min=1,max=1000"`
    FieldConfig *ImageContentFieldConfigRequest `json:"fieldConfig" binding:"required"`
}

type UpdateImageContentTypeRequest struct {
    Name        string                          `json:"name" binding:"required,max=100"`
    Status      string                          `json:"status" binding:"required,oneof=ACTIVE INACTIVE"`
    SortOrder   int                             `json:"sortOrder" binding:"gte=0"`
    MaxItems    *uint                           `json:"maxItems" binding:"omitempty,min=1,max=1000"`
    FieldConfig *ImageContentFieldConfigRequest `json:"fieldConfig" binding:"required"`
}

type GetImageContentsQuery struct {
    Start    int    `form:"_start" binding:"omitempty,min=0"`
    End      int    `form:"_end" binding:"omitempty,min=1,max=1000"`
    Search   string `form:"q" binding:"omitempty,max=255"`
    TypeCode string `form:"typeCode" binding:"omitempty,max=50"`
    Status   string `form:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
    SortBy   string `form:"_sort" binding:"omitempty,oneof=id typeCode name status sortOrder createdAt updatedAt"`
    Order    string `form:"_order" binding:"omitempty,oneof=ASC DESC asc desc"`
}

type CreateImageContentRequest struct {
    TypeCode            string  `json:"typeCode" binding:"required,max=50"`
    MediaID             uint    `json:"mediaId" binding:"required,min=1"`
    Name                *string `json:"name" binding:"omitempty,max=255"`
    Description         *string `json:"description" binding:"omitempty,max=5000"`
    SecondaryDescription *string `json:"secondaryDescription" binding:"omitempty,max=5000"`
    URL                 *string `json:"url" binding:"omitempty,max=500"`
    SortOrder           *int    `json:"sortOrder" binding:"omitempty,gte=0"`
    Status              string  `json:"status" binding:"required,oneof=ACTIVE INACTIVE"`
}

type UpdateImageContentRequest struct {
    MediaID             uint    `json:"mediaId" binding:"required,min=1"`
    Name                *string `json:"name" binding:"omitempty,max=255"`
    Description         *string `json:"description" binding:"omitempty,max=5000"`
    SecondaryDescription *string `json:"secondaryDescription" binding:"omitempty,max=5000"`
    URL                 *string `json:"url" binding:"omitempty,max=500"`
    Status              string  `json:"status" binding:"required,oneof=ACTIVE INACTIVE"`
}

type ImageContentOrderItem struct {
    ID        uint `json:"id" binding:"required,min=1"`
    SortOrder int  `json:"sortOrder" binding:"gte=0"`
}

type ReorderImageContentsRequest struct {
    TypeCode string                  `json:"typeCode" binding:"required,max=50"`
    Items    []ImageContentOrderItem `json:"items" binding:"required,max=1000,dive"`
}

type GetPublicImageContentsQuery struct {
    TypeCodes string `form:"typeCodes" binding:"omitempty,max=1024"`
}
```

- `_start/_end/_sort/_order` khớp Simple REST provider hiện tại. Service kiểm tra `_end > _start`, giới hạn page size tối đa 100 và dùng default `_start=0`, `_end=20`.
- Public `typeCodes` là danh sách comma-separated, tối đa 20 code unique. Mỗi code normalize uppercase và validate regex; code không tồn tại/inactive được bỏ qua thay vì trả 404.

### 5.2. Response DTOs

```go
type ImageContentFieldRuleResponse struct {
    Enabled  bool `json:"enabled"`
    Required bool `json:"required"`
}

type ImageContentFieldConfigResponse struct {
    Name                 ImageContentFieldRuleResponse `json:"name"`
    Description          ImageContentFieldRuleResponse `json:"description"`
    SecondaryDescription ImageContentFieldRuleResponse `json:"secondaryDescription"`
    URL                  ImageContentFieldRuleResponse `json:"url"`
}

type ImageContentTypeResponse struct {
    ID          uint                            `json:"id"`
    Code        string                          `json:"code"`
    Name        string                          `json:"name"`
    Status      string                          `json:"status"`
    SortOrder   int                             `json:"sortOrder"`
    MaxItems    *uint                           `json:"maxItems"`
    FieldConfig ImageContentFieldConfigResponse `json:"fieldConfig"`
    ItemCount   int64                           `json:"itemCount"`
    CreatedAt   time.Time                       `json:"createdAt"`
    UpdatedAt   time.Time                       `json:"updatedAt"`
}

type AdminImageContentMediaResponse struct {
    ID           uint   `json:"id"`
    FileName     string `json:"fileName"`
    OriginalURL  string `json:"originalUrl"`
    ThumbnailURL string `json:"thumbnailUrl,omitempty"`
    MediumURL    string `json:"mediumUrl,omitempty"`
    MimeType     string `json:"mimeType"`
    Status       string `json:"status"`
}

type ImageContentResponse struct {
    ID                   uint                            `json:"id"`
    TypeCode             string                          `json:"typeCode"`
    TypeName             string                          `json:"typeName"`
    FieldConfig          ImageContentFieldConfigResponse `json:"fieldConfig"`
    MediaID              uint                            `json:"mediaId"`
    Media                AdminImageContentMediaResponse  `json:"media"`
    Name                 *string                         `json:"name"`
    Description          *string                         `json:"description"`
    SecondaryDescription *string                         `json:"secondaryDescription"`
    URL                  *string                         `json:"url"`
    SortOrder            int                             `json:"sortOrder"`
    Status               string                          `json:"status"`
    CreatedBy            uint                            `json:"createdBy"`
    CreatedAt            time.Time                       `json:"createdAt"`
    UpdatedAt            time.Time                       `json:"updatedAt"`
}

type ImageContentTypeListResponse struct {
    Data  []ImageContentTypeResponse `json:"data"`
    Total int64                      `json:"total"`
}

type ImageContentListResponse struct {
    Data  []ImageContentResponse `json:"data"`
    Total int64                  `json:"total"`
}

type ImageContentErrorResponse struct {
    Error   string            `json:"error"`
    Code    string            `json:"code"`
    Message string            `json:"message"`
    Fields  map[string]string `json:"fields,omitempty"`
}
```

Admin response luôn trả raw metadata kể cả field đang disabled để không làm mất dữ liệu khi Type tạm đổi cấu hình. Không trả GORM relation hoặc thông tin nhạy cảm của User/Media owner.

### 5.3. Admin Type endpoints

| Method | Route | Auth/RBAC | Thành công |
| --- | --- | --- | --- |
| `GET` | `/api/v1/admin/image-content-types` | Access Token + `ADMIN` | `200 ImageContentTypeListResponse` |
| `POST` | `/api/v1/admin/image-content-types` | Access Token + `ADMIN` | `201 {data: ImageContentTypeResponse}` |
| `GET` | `/api/v1/admin/image-content-types/:id` | Access Token + `ADMIN` | `200 {data: ImageContentTypeResponse}` |
| `PUT` | `/api/v1/admin/image-content-types/:id` | Access Token + `ADMIN` | `200 {data: ImageContentTypeResponse}` |
| `DELETE` | `/api/v1/admin/image-content-types/:id` | Access Token + `ADMIN` | `204` |

### 5.4. Admin Item endpoints

| Method | Route | Auth/RBAC | Thành công |
| --- | --- | --- | --- |
| `GET` | `/api/v1/admin/image-contents` | Access Token + `ADMIN/STAFF` | `200 ImageContentListResponse` |
| `POST` | `/api/v1/admin/image-contents` | Access Token + `ADMIN/STAFF` | `201 {data: ImageContentResponse}` |
| `GET` | `/api/v1/admin/image-contents/:id` | Access Token + `ADMIN/STAFF` | `200 {data: ImageContentResponse}` |
| `PUT` | `/api/v1/admin/image-contents/:id` | Access Token + `ADMIN/STAFF` | `200 {data: ImageContentResponse}` |
| `DELETE` | `/api/v1/admin/image-contents/:id` | Access Token + `ADMIN/STAFF` | `204` |
| `PUT` | `/api/v1/admin/image-contents/order` | Access Token + `ADMIN/STAFF` | `200 {data: [], message: string}` |

Route `/order` phải đăng ký trước `/:id` để Gin không coi `order` là ID.

Tất cả Admin route:

1. Đi qua `AuthMiddleware(redisClient)` để verify token, expiry, signature và blacklist.
2. Đi qua `RoleMiddleware(...)` theo bảng trên.
3. Controller lấy `userId` và `role` từ Gin Context; không nhận hai giá trị này từ payload/header tùy ý.

### 5.5. Error mapping

| Điều kiện | HTTP | Code |
| --- | ---: | --- |
| Binding, code, URL, pagination hoặc field config sai | 400 | `VALIDATION_ERROR` |
| Chưa đăng nhập/token sai/hết hạn/revoked | 401 | `UNAUTHORIZED` |
| Sai role hoặc Staff dùng Media không sở hữu | 403 | `FORBIDDEN` |
| Type không tồn tại | 404 | `TYPE_NOT_FOUND` |
| Item không tồn tại | 404 | `IMAGE_CONTENT_NOT_FOUND` |
| Media không tồn tại | 404 | `MEDIA_NOT_FOUND` |
| Media không phải ảnh | 400 | `MEDIA_TYPE_INVALID` |
| Field required thiếu | 400 | `FIELD_REQUIRED` |
| Request gửi field bị disable | 400 | `FIELD_NOT_ENABLED` |
| Code Type đã tồn tại | 409 | `TYPE_CODE_CONFLICT` |
| Update cấu hình không tương thích | 409 | `TYPE_CONFIG_CONFLICT` |
| Type đã đạt giới hạn Item | 409 | `MAX_ITEMS_EXCEEDED` |
| Xóa Type đang có Item | 409 | `TYPE_IN_USE` |
| Snapshot reorder stale/sai | 409 | `ORDER_CONFLICT` |
| DB/Redis/auth unexpected | 500 | `INTERNAL_ERROR` |

Controller không trả SQL error, Redis error, stack trace hoặc chi tiết filesystem.

---

## 6. Public API contract

### 6.1. Endpoint

```text
GET /api/v1/public/image-contents?typeCodes=LOGO,SLIDER,PARTNER
```

- Public, không dùng AuthMiddleware.
- Nếu không có `typeCodes`, trả tất cả Type `ACTIVE`.
- Type hợp lệ nhưng không có Item trả `items: []`.
- Type không tồn tại hoặc `INACTIVE` bị bỏ khỏi `data`; toàn bộ kết quả rỗng vẫn trả `200 {"data":[]}`.

### 6.2. Public response DTO

```go
type PublicImageContentMediaResponse struct {
    ID           uint   `json:"id"`
    OriginalURL  string `json:"originalUrl"`
    ThumbnailURL string `json:"thumbnailUrl,omitempty"`
    MediumURL    string `json:"mediumUrl,omitempty"`
    MimeType     string `json:"mimeType"`
}

type PublicImageContentItemResponse struct {
    ID                   uint                             `json:"id"`
    Name                 *string                          `json:"name,omitempty"`
    Description          *string                          `json:"description,omitempty"`
    SecondaryDescription *string                          `json:"secondaryDescription,omitempty"`
    URL                  *string                          `json:"url,omitempty"`
    SortOrder            int                              `json:"sortOrder"`
    Media                PublicImageContentMediaResponse  `json:"media"`
}

type PublicImageContentGroupResponse struct {
    TypeCode string                           `json:"typeCode"`
    Name     string                           `json:"name"`
    Items    []PublicImageContentItemResponse `json:"items"`
}

type PublicImageContentResponse struct {
    Data []PublicImageContentGroupResponse `json:"data"`
}
```

Projection rules:

- Chỉ Type/Item `ACTIVE`, chưa soft-delete.
- Media phải chưa soft-delete, `attached` và MIME `image/%`; Item có Media không còn public-valid bị bỏ khỏi response, không xóa DB tự động.
- Type order: `image_content_types.sort_order ASC, id ASC`.
- Item order: `image_contents.sort_order ASC, id ASC`.
- Field disabled trong `field_config` bị ép `nil` trước serialize, kể cả DB còn raw value.
- Mảng `data/items` luôn là `[]`, không trả `null`.

Ví dụ:

```json
{
  "data": [
    {
      "typeCode": "SLIDER",
      "name": "Slider",
      "items": [
        {
          "id": 12,
          "name": "Giải pháp cho doanh nghiệp",
          "description": "Nội dung mô tả chính",
          "secondaryDescription": "Nội dung bổ sung",
          "url": "/gioi-thieu",
          "sortOrder": 0,
          "media": {
            "id": 42,
            "originalUrl": "/uploads/slider.webp",
            "mediumUrl": "/uploads/slider-medium.webp",
            "mimeType": "image/webp"
          }
        }
      ]
    }
  ]
}
```

---

## 7. Trụ cột 3 — Redis cache và quản lý token

### 7.1. Public cache

Key duy nhất cho aggregate toàn site:

```text
public:image-contents:v1
```

- Cache-aside toàn bộ `PublicImageContentResponse` chưa filter, TTL 15 phút bằng `SetEx`.
- Public request có `typeCodes` đọc aggregate từ cache/DB rồi filter trong memory; tối đa 20 Type nên không tạo biến thể key theo query.
- Không cache GORM entity, Admin response, raw disabled fields, owner hoặc token.
- Mọi create/update/delete Type/Item và reorder Item gọi `DEL public:image-contents:v1` sau commit.
- Redis read/unmarshal lỗi: log, fallback MySQL và vẫn trả response nếu DB thành công.
- Redis `SetEx`/`DEL` lỗi: log, không biến mutation đã commit thành lỗi; TTL là safety net.
- Không dùng Redis `KEYS` hoặc wildcard invalidation.
- Nếu sau này Media có API cập nhật URL/variant, Media Service phải reverse lookup `image_contents.media_id` và invalidate cùng key. Upload Media chưa gắn không invalidate.

### 7.2. JWT blacklist, refresh rotation và logout

Module không tạo token key riêng. Admin routes tái sử dụng các pattern Auth hiện tại:

```text
auth:bl:at:{jti}
auth:rf:{user_id}:{jti}
auth:rf:index:{user_id}
```

Target security state bắt buộc cho shared Auth:

- Middleware verify algorithm, signature và expiry trước, sau đó check `auth:bl:at:{jti}`; identity chỉ được đưa vào Gin Context sau khi toàn bộ kiểm tra thành công.
- Logout phải đi qua AuthMiddleware, xóa refresh session/JTI và ghi blacklist bằng `SetEx` với TTL chính xác `exp-now`; TTL không dương thì không ghi key rác.
- Refresh Token Rotation phải atomic bằng Redis Lua hoặc transaction tương đương: kiểm tra old key, xóa old JTI, tạo new JTI/index và TTL trong một thao tác logic.
- Old refresh JTI không còn tồn tại được coi là replay; dùng `auth:rf:index:{user_id}` để thu hồi toàn bộ session của user và buộc đăng nhập lại.
- Chỉ lưu JTI/session marker trong Redis, không lưu full access/refresh token.
- Auth hiện đã có `auth:bl:at:{jti}` và `auth:rf:{user_id}:{jti}`, nhưng rotation đang gồm nhiều lệnh rời và secret đang hard-code. Đây là security debt dùng chung cần harden trước production; không nhân bản auth logic trong ImageContent.

---

## 8. Wire, routes, OpenAPI và seeder

### 8.1. Google Wire

```go
var ImageContentSet = wire.NewSet(
    repository.NewImageContentRepository,
    imageContentQuery.NewImageContentQueryService,
    imageContentService.NewImageContentService,
    controller.NewImageContentController,
)

func InitializeImageContentController(
    db *gorm.DB,
    redisClient *redis.Client,
) *controller.ImageContentController {
    wire.Build(ImageContentSet)
    return &controller.ImageContentController{}
}
```

- Thêm initializer vào `wire.go`, regenerate `wire_gen.go`; không sửa generated file thủ công.
- Controller nhận Command Service và Query Service qua constructor, không tự query GORM/Redis.

### 8.2. Route registration

```go
imageContentController := di.InitializeImageContentController(db, redisClient)

imageContentTypes := admin.Group("/image-content-types")
imageContentTypes.Use(middleware.RoleMiddleware("ADMIN"))
{
    imageContentTypes.GET("", imageContentController.GetTypes)
    imageContentTypes.POST("", imageContentController.CreateType)
    imageContentTypes.GET("/:id", imageContentController.GetType)
    imageContentTypes.PUT("/:id", imageContentController.UpdateType)
    imageContentTypes.DELETE("/:id", imageContentController.DeleteType)
}

imageContents := admin.Group("/image-contents")
imageContents.Use(middleware.RoleMiddleware("ADMIN", "STAFF"))
{
    imageContents.GET("", imageContentController.GetContents)
    imageContents.POST("", imageContentController.CreateContent)
    imageContents.PUT("/order", imageContentController.ReorderContents)
    imageContents.GET("/:id", imageContentController.GetContent)
    imageContents.PUT("/:id", imageContentController.UpdateContent)
    imageContents.DELETE("/:id", imageContentController.DeleteContent)
}

public := api.Group("/public")
public.GET("/image-contents", imageContentController.GetPublicContents)
```

### 8.3. OpenAPI

Cập nhật `.docs/api-endpoints.yaml`:

- Mười hai endpoint Type/Item/Public ở mục 5 và 6.
- BearerAuth cho Admin routes; Public route không khai báo BearerAuth.
- Đủ schema request/response/error và status `200/201/204/400/401/403/404/409/500`.
- Mô tả rõ query Simple REST `_start/_end/_sort/_order`, filter `typeCode/status/q` và public `typeCodes`.

### 8.4. Seeder

Tạo `apps/backend/cmd/seeder/image_content_types.go` và gọi sau `SeedUsers`, trước mọi seeder Item tương lai.

| Code | Name | Sort | Max | Field config |
| --- | --- | ---: | ---: | --- |
| `LOGO` | Logo | 0 | 1 | Tất cả metadata disabled |
| `SLIDER` | Slider | 1 | null | name/description/url required; secondaryDescription enabled optional |
| `PARTNER` | Đối tác | 2 | null | name/url required; hai description disabled |

- Seeder lookup `Unscoped` theo code; insert nếu thiếu, restore nếu row đã soft-delete.
- Seeder không ghi đè Type đang active để giữ cấu hình do Admin đã thay đổi.
- Không seed Item hoặc Media mẫu trong phạm vi này.

---

## 9. Kiểm thử và acceptance criteria

### 9.1. Database/repository

- Migration up/down đúng thứ tự; kiểu FK khớp `media.id` unsigned và `users.id` signed.
- CHECK constraint chặn status/order/maxItems sai; JSON field lưu/đọc được bằng GORM serializer.
- Type code unique cả khi row soft-delete; create lại restore đúng row.
- Media FK `RESTRICT` chặn hard delete ngoài Service khi đang được tham chiếu.
- List/search/filter/sort dùng whitelist, escape wildcard và trả `{data,total}` đúng Refine.
- Public query có số query cố định, không N+1, order ổn định.

### 9.2. Service/transaction

- Từ chối `required && !enabled`, Type code sai và maxItems ngoài giới hạn.
- Update Type thành config required không hợp lệ với Item hiện hữu trả conflict và không ghi một phần.
- `LOGO` không tạo được Item thứ hai kể cả Item đầu đang `INACTIVE`.
- Slider thiếu name/description/url; Partner thiếu name/url; Logo gửi metadata disabled đều bị từ chối.
- Media thiếu, soft-delete, không phải ảnh hoặc Staff không sở hữu đều bị từ chối trước insert/update.
- Create vị trí giữa, delete và reorder luôn tạo thứ tự `0..n-1`.
- Reorder thiếu/thừa/trùng/ID khác Type trả `ORDER_CONFLICT` và rollback toàn bộ.
- Delete Item không xóa file/Media và không đổi status Media về temporary.
- Redis invalidation chỉ chạy sau commit; Redis lỗi không rollback DB.

### 9.3. Handler/API/Auth

- Admin Type: CUSTOMER/STAFF bị 403; ADMIN CRUD được.
- Admin Item: CUSTOMER bị 403; ADMIN/STAFF được, với ownership Media đúng quy tắc.
- Token thiếu/sai/hết hạn/blacklisted trả 401 trước Handler.
- Binding URI/query/body sai trả error contract thống nhất, không lộ internal error.
- Public không cần token, chỉ trả Type/Item active, bỏ field disabled và luôn trả array rỗng thay vì null.
- Public filter normalize/deduplicate code, từ chối format sai và giới hạn tối đa 20 Type.
- Cache miss/hit cho response tương đương; cache JSON hỏng hoặc Redis unavailable fallback DB.

### 9.4. Lệnh xác minh

```text
go test ./...
go vet ./...
wire ./internal/di
```

Chạy migration up/down trên MySQL 8, chạy seeder hai lần và kiểm tra OpenAPI contract với các endpoint thực tế.

### 9.5. Tiêu chí hoàn tất

- Quản trị được Type và Item cho `LOGO`, `SLIDER`, `PARTNER` qua API tương thích Refine.
- Có thể tạo Type mới sử dụng tổ hợp field chuẩn mà không đổi database hoặc backend.
- `LOGO` luôn tối đa một Item; field requirement của Slider/Partner được backend cưỡng chế.
- Public API trả aggregate đúng Type/order/status, không lộ dữ liệu Admin và được cache/invalidate an toàn.
- Không tạo luồng upload mới, không thay đổi schema hoặc hành vi của Post/Page/Page Section.
- Auth Context, RBAC, Wire, migration, seeder và OpenAPI được quy hoạch đầy đủ để bước `/code-db` và `/code-api` triển khai không cần quyết định thêm.
