# QUY HOẠCH KIẾN TRÚC BACK-END: QUẢN LÝ SEO CHO BÀI VIẾT (POST SEO META)

**Dự án:** TechBite  
**Module:** SEO Meta dùng chung  
**Tài liệu nguồn:** `.docs/ideas/dashboard/20-seo-meta-idea.md`  
**Phạm vi triển khai đầu tiên:** `post`, đồng thời thiết kế registry để hỗ trợ `product`, `post_category`, `product_category`, `page`.

## 1. Mục tiêu và nguyên tắc

- Tạo bảng `seo_meta` làm nguồn dữ liệu duy nhất cho metadata SEO.
- `post_meta` tiếp tục chỉ lưu metadata mở rộng dạng key-value; tuyệt đối không đọc hoặc ghi các key SEO trong bảng này.
- Mỗi entity chỉ có một SEO record, chưa có `locale`.
- `entity_type` và `entity_id` chỉ lấy từ URI. Backend dùng registry whitelist để xác minh loại entity và kiểm tra bản ghi tồn tại, chưa soft-delete.
- Vì quan hệ polymorphic không thể tạo foreign key tới nhiều bảng, tính toàn vẹn liên kết được bảo đảm ở Service và lifecycle hook/use case xóa entity.
- Không tự động ghi giá trị fallback vào database.

## 2. Kiến trúc DDD/Clean Architecture

Luồng chính:

`Gin Route -> AuthMiddleware -> RoleMiddleware -> SEO Meta Controller -> SEO Meta Service/Resolver -> Repository/Entity Registry -> MySQL/Redis`

### 2.1. Thành phần cần tạo

```text
apps/backend/internal/domain/seometa/entity/seo_meta.go
apps/backend/internal/infrastructure/repository/seo_meta_repo.go
apps/backend/internal/infrastructure/repository/entity_registry.go
apps/backend/internal/usecases/seometa/dto/seo_meta_dto.go
apps/backend/internal/usecases/seometa/service/seo_meta_service.go
apps/backend/internal/controller/seo_meta_controller.go
apps/backend/database/migrations/<timestamp>_create_seo_meta_table.up.sql
apps/backend/database/migrations/<timestamp>_create_seo_meta_table.down.sql
```

Đăng ký dependency trong `internal/di/wire.go` và regenerate `wire_gen.go` bằng Google Wire. Handler chỉ bind URI/JSON và map lỗi; business rule, resolver, registry và cache nằm ở Service; Repository chỉ truy vấn GORM.

### 2.2. Registry entity

Tạo `EntityRegistry` với các implementation/checker tái sử dụng Repository hiện có:

```go
type EntityType string

const (
	EntityPost            EntityType = "post"
	EntityProduct         EntityType = "product"
	EntityPostCategory    EntityType = "post_category"
	EntityProductCategory EntityType = "product_category"
	EntityPage            EntityType = "page"
)

type EntityReference struct {
	Type         EntityType
	ID           uint
	Title        string
	Slug         string
	Description  string
	Excerpt      string
	Content      string
	ThumbnailURL string
}

type EntityRegistry interface {
	Resolve(ctx context.Context, entityType EntityType, entityID uint) (*EntityReference, error)
}
```

- Registry chỉ chấp nhận 5 giá trị whitelist; giá trị khác trả lỗi `unsupported_entity_type`.
- `Resolve` phải query với điều kiện `deleted_at IS NULL`; không trả SEO cho entity không tồn tại hoặc đã soft-delete.
- `product_category` cần ánh xạ đúng entity/bảng danh mục sản phẩm hiện tại (`categories`) theo kiến trúc hiện có; không suy diễn từ `post_categories`.
- `page` phải có checker/repository khi entity Page được tạo. Trong thời gian chưa có bảng Page, loại này trả lỗi cấu hình rõ ràng, không cho tạo SEO mồ côi.
- Không nhận `userId`, `role`, `authorId` hoặc entity identity từ body/query.

## 3. Trụ cột 1: Thiết kế dữ liệu

### 3.1. GORM Entity

```go
package entity

import (
	"time"

	"gorm.io/datatypes"
)

type SEOMeta struct {
	ID                 uint           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	EntityType         string         `gorm:"column:entity_type;type:varchar(50);not null" json:"entity_type"`
	EntityID           uint           `gorm:"column:entity_id;not null" json:"entity_id"`
	MetaTitle          *string        `gorm:"column:meta_title;type:varchar(255)" json:"meta_title"`
	MetaDescription    *string        `gorm:"column:meta_description;type:varchar(500)" json:"meta_description"`
	MetaKeywords       *string        `gorm:"column:meta_keywords;type:varchar(500)" json:"meta_keywords"`
	CanonicalURL       *string        `gorm:"column:canonical_url;type:varchar(500)" json:"canonical_url"`
	OGTitle            *string        `gorm:"column:og_title;type:varchar(255)" json:"og_title"`
	OGDescription      *string        `gorm:"column:og_description;type:varchar(500)" json:"og_description"`
	OGImage            *string        `gorm:"column:og_image;type:varchar(500)" json:"og_image"`
	OGType             *string        `gorm:"column:og_type;type:varchar(100)" json:"og_type"`
	TwitterTitle       *string        `gorm:"column:twitter_title;type:varchar(255)" json:"twitter_title"`
	TwitterDescription *string        `gorm:"column:twitter_description;type:varchar(500)" json:"twitter_description"`
	TwitterImage       *string        `gorm:"column:twitter_image;type:varchar(500)" json:"twitter_image"`
	TwitterCard        *string        `gorm:"column:twitter_card;type:varchar(100)" json:"twitter_card"`
	Robots             *string        `gorm:"column:robots;type:varchar(100)" json:"robots"`
	SchemaJSON         datatypes.JSON `gorm:"column:schema_json;type:json" json:"schema_json"`
	CreatedAt          time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (SEOMeta) TableName() string { return "seo_meta" }
```

`*string` bảo toàn sự khác biệt giữa `NULL`, chuỗi rỗng và whitespace để Resolver xử lý đúng. `schema_json` không trả stack trace hay dữ liệu nội bộ; DTO response phải dùng field đã validate/serialize.

### 3.2. Migration và ràng buộc MySQL

Migration `create_seo_meta_table` cần tạo:

```sql
CREATE TABLE `seo_meta` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `entity_type` varchar(50) NOT NULL,
  `entity_id` bigint unsigned NOT NULL,
  `meta_title` varchar(255) DEFAULT NULL,
  `meta_description` varchar(500) DEFAULT NULL,
  `meta_keywords` varchar(500) DEFAULT NULL,
  `canonical_url` varchar(500) DEFAULT NULL,
  `og_title` varchar(255) DEFAULT NULL,
  `og_description` varchar(500) DEFAULT NULL,
  `og_image` varchar(500) DEFAULT NULL,
  `og_type` varchar(100) DEFAULT NULL,
  `twitter_title` varchar(255) DEFAULT NULL,
  `twitter_description` varchar(500) DEFAULT NULL,
  `twitter_image` varchar(500) DEFAULT NULL,
  `twitter_card` varchar(100) DEFAULT NULL,
  `robots` varchar(100) DEFAULT NULL,
  `schema_json` json DEFAULT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_seo_meta_entity` (`entity_type`, `entity_id`),
  KEY `idx_seo_meta_entity` (`entity_type`, `entity_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

- Không thêm `DeletedAt`: SEO được xóa theo lifecycle entity hoặc cleanup use case; không dùng soft-delete cho bảng này.
- `entity_id` là số nguyên dương; không có foreign key vật lý do polymorphic relation.
- Unique key là lớp bảo vệ cuối cùng cho upsert đồng thời. Repository phải map duplicate-key race thành lỗi xung đột phù hợp.
- Down migration chỉ được drop bảng `seo_meta` khi chạy migration rollback có chủ đích.
- Chưa thêm `locale`; khi mở rộng đa ngôn ngữ, đổi unique key thành `(entity_type, entity_id, locale)` và bổ sung fallback locale.

### 3.3. Repository contract

```go
type SEOMetaRepository interface {
	FindByEntity(ctx context.Context, entityType string, entityID uint) (*entity.SEOMeta, error)
	Upsert(ctx context.Context, seo *entity.SEOMeta) (*entity.SEOMeta, error)
	DeleteByEntity(ctx context.Context, entityType string, entityID uint) error
}
```

`Upsert` phải dùng transaction hoặc MySQL upsert theo unique key, chỉ update whitelist SEO fields. `DeleteByEntity` được gọi trong transaction/lifecycle xóa entity nếu cùng DB transaction có thể bảo đảm; nếu entity module tách transaction, ghi log và đưa cleanup vào service retryable.

## 4. Trụ cột 2: API Contract và Auth

Base route: `/api/v1/admin/seo-meta`.

Tất cả endpoint kế thừa `AuthMiddleware(redisClient)` và `RoleMiddleware("ADMIN", "STAFF")`. Auth Middleware hiện có đã verify JWT trước, kiểm tra `auth:bl:at:{jti}` trong Redis và set `userId`/`role` vào Gin Context. Service/Handler chỉ dùng các context value này cho audit/authorization, không nhận identity từ payload.

### 4.1. GET SEO của entity

**Method:** `GET`  
**Route:** `/api/v1/admin/seo-meta/:entityType/:entityId`  
**Auth:** JWT hợp lệ + role `ADMIN` hoặc `STAFF`.

```go
type GetSEOMetaURI struct {
	EntityType string `uri:"entityType" binding:"required,max=50"`
	EntityID   uint   `uri:"entityId" binding:"required,min=1"`
}

type SEOMetaResponse struct {
	EntityType         string          `json:"entity_type"`
	EntityID           uint            `json:"entity_id"`
	MetaTitle          *string         `json:"meta_title"`
	ResolvedTitle      string          `json:"resolved_title"`
	MetaDescription    *string         `json:"meta_description"`
	ResolvedDescription string         `json:"resolved_description"`
	MetaKeywords       *string         `json:"meta_keywords"`
	CanonicalURL       *string         `json:"canonical_url"`
	ResolvedCanonicalURL string        `json:"resolved_canonical_url"`
	OGTitle            *string         `json:"og_title"`
	OGDescription      *string         `json:"og_description"`
	OGImage            *string         `json:"og_image"`
	OGType             *string         `json:"og_type"`
	TwitterTitle       *string         `json:"twitter_title"`
	TwitterDescription *string         `json:"twitter_description"`
	TwitterImage       *string         `json:"twitter_image"`
	TwitterCard        *string         `json:"twitter_card"`
	Robots             string          `json:"robots"`
	SchemaJSON         datatypes.JSON  `json:"schema_json"`
}

type GetSEOMetaResponse struct {
	Data SEOMetaResponse `json:"data"`
}
```

Nếu entity hợp lệ nhưng chưa có SEO record, trả `200` với các field gốc là `null` và field resolved dùng fallback. Nếu entity không tồn tại/đã soft-delete, trả `404`.

### 4.2. PUT upsert SEO

**Method:** `PUT`  
**Route:** `/api/v1/admin/seo-meta/:entityType/:entityId`  
**Auth:** JWT hợp lệ + role `ADMIN` hoặc `STAFF`.

```go
type UpsertSEOMetaURI struct {
	EntityType string `uri:"entityType" binding:"required,max=50"`
	EntityID   uint   `uri:"entityId" binding:"required,min=1"`
}

type UpsertSEOMetaRequest struct {
	MetaTitle          *string         `json:"meta_title"`
	MetaDescription    *string         `json:"meta_description"`
	MetaKeywords       *string         `json:"meta_keywords"`
	CanonicalURL       *string         `json:"canonical_url"`
	OGTitle            *string         `json:"og_title"`
	OGDescription      *string         `json:"og_description"`
	OGImage            *string         `json:"og_image"`
	OGType             *string         `json:"og_type"`
	TwitterTitle       *string         `json:"twitter_title"`
	TwitterDescription *string         `json:"twitter_description"`
	TwitterImage       *string         `json:"twitter_image"`
	TwitterCard        *string         `json:"twitter_card"`
	Robots             *string         `json:"robots"`
	SchemaJSON         datatypes.JSON  `json:"schema_json"`
}

type UpsertSEOMetaResponse struct {
	Data    SEOMetaResponse `json:"data"`
	Message string          `json:"message"`
}
```

`entity_type`/`entity_id` không được xuất hiện trong request body. Service phải resolve entity trước khi validate và ghi. Upsert lặp lại không tạo record trùng.

### 4.3. Xóa SEO thủ công

**Method:** `DELETE`  
**Route:** `/api/v1/admin/seo-meta/:entityType/:entityId`  
**Auth:** JWT hợp lệ + role `ADMIN` hoặc `STAFF`.

```go
type DeleteSEOMetaResponse struct {
	Message string `json:"message"`
}
```

Xóa SEO không xóa entity. Có thể trả `204 No Content` hoặc `200` theo convention response hiện tại; ưu tiên `200` với message để tương thích dashboard. Sau thành công bắt buộc invalidate cache.

### 4.4. Chuẩn lỗi

```go
type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}
```

- `400`: URI/body không bind được, entity ID không hợp lệ.
- `401`: thiếu, sai, hết hạn hoặc blacklist JWT.
- `403`: role không phải `ADMIN`/`STAFF`.
- `404`: entity không tồn tại/đã soft-delete hoặc SEO record được yêu cầu không tồn tại theo policy.
- `409`: cạnh tranh unique key hoặc xung đột upsert.
- `422`: validation SEO thất bại (`canonical_url`, `robots`, JSON schema...).
- `500`: lỗi MySQL/Redis hoặc lỗi hệ thống; không trả stack trace.

## 5. Business Validation và SEO Resolver

### 5.1. Validation tại Service

- Trim toàn bộ string; chuỗi rỗng sau trim chuyển thành `nil`.
- `meta_title` tối đa 255 ký tự; `meta_description` tối đa 500; keywords tối đa 500.
- URL canonical/image phải parse được và là URL hợp lệ; không chấp nhận scheme nguy hiểm như `javascript:`. Image hiện lưu URL, chưa đổi sang `media_id`.
- `robots` mặc định `index,follow`; chỉ nhận `index,follow`, `noindex,follow`, `noindex,nofollow`.
- `og_type`/`twitter_card` dùng whitelist cấu hình của module; không nhận giá trị tùy ý nếu chưa được hệ thống chấp nhận.
- `schema_json` phải là JSON hợp lệ, giới hạn kích thước payload; chỉ chấp nhận schema JSON-LD type được hệ thống cho phép. Khi render phải serialize an toàn, không nối chuỗi HTML/script tùy ý.

### 5.2. Resolver

Resolver nhận `SEOMeta` nullable và `EntityReference`, trả DTO đã lưu + resolved:

```text
resolved_title           = meta_title hoặc entity.title/name
resolved_description     = meta_description hoặc excerpt hoặc content strip HTML rồi cắt Unicode an toàn
resolved_canonical_url   = canonical_url hoặc base_url + entity.slug
resolved_og_title        = og_title hoặc resolved_title
resolved_og_description  = og_description hoặc resolved_description
resolved_og_image        = og_image hoặc entity.thumbnail
resolved_twitter_title   = twitter_title hoặc resolved_og_title
resolved_twitter_desc    = twitter_description hoặc resolved_og_description
resolved_twitter_image   = twitter_image hoặc resolved_og_image
robots                   = robots hoặc index,follow
```

Fallback phụ thuộc `entity_type`: Post dùng `Title`, `Content`, `Slug`, thumbnail; Product dùng `Name`, `Description`, `Slug`, `ImageURL`; category dùng `Name`, `Description`, `Slug`, `ImageURL`; Page dùng adapter tương ứng. Không ghi fallback ngược vào `seo_meta`.

## 6. Redis Cache và token state

### 6.1. Cache SEO

- Key chuẩn: `seo:{entity_type}:{entity_id}`.
- Dữ liệu cache là JSON của response resolved hoặc domain snapshot; dùng cache-aside với TTL cấu hình, ví dụ 1 giờ.
- GET: đọc Redis trước, cache miss đọc MySQL + registry, sau đó set cache.
- PUT/DELETE: ghi/xóa MySQL thành công rồi `DEL seo:{entity_type}:{entity_id}`.
- Redis lỗi không được làm fail request đã ghi MySQL thành công; phải log lỗi invalidate. Tránh trả dữ liệu cache cũ sau PUT bằng cách chỉ set cache mới sau khi DB commit.
- Khi entity update title/slug/content/thumbnail, lifecycle entity phải invalidate cùng key vì fallback đã thay đổi.

### 6.2. JWT blacklist/refresh/logout

SEO API không tạo token mới nhưng bắt buộc tái sử dụng cơ chế hiện có:

- `AuthMiddleware` verify chữ ký trước rồi mới kiểm tra Redis key `auth:bl:at:{jti}`; nếu tồn tại abort `401`.
- `userId` và `role` chỉ lấy từ claims đã set trong Gin Context.
- Khi `/logout`, blacklist JTI/access token với `SetEx`/`Set` theo TTL còn lại của JWT; xóa refresh token state tương ứng.
- `/refresh-token` phải rotate refresh JTI trong một luồng atomic: revoke token cũ, tạo cặp token mới, lưu state mới; replay phải thu hồi toàn bộ refresh token của user theo chính sách Auth.
- Refresh token chỉ nằm trong HttpOnly, Secure, SameSite cookie; không đưa token hoặc stack trace vào response SEO.

## 7. Tích hợp lifecycle và DI

- Khi `Post`, `Product`, `PostCategory`, `Category` (product category) hoặc `Page` bị xóa, use case xóa entity phải gọi `DeleteByEntity` tương ứng hoặc phát hiện cleanup trong cùng transaction/service boundary.
- Không gộp SEO vào transaction Create/Update Post nếu không có yêu cầu; Dashboard gọi PUT SEO riêng sau khi entity tạo/cập nhật thành công.
- Bổ sung `SEOMetaSet` trong `wire.go`: `SEO Meta Repository -> Entity Registry -> SEO Meta Service -> SEO Meta Controller`, truyền `*gorm.DB` và `*redis.Client` khi cần.
- Thêm route group `admin.Group("/seo-meta")`, kế thừa Auth Middleware và gắn `RoleMiddleware("ADMIN", "STAFF")` cho GET/PUT/DELETE.
- Cập nhật OpenAPI/API endpoint registry nếu dự án đang quản lý contract tại file riêng.

## 8. Tiêu chí nghiệm thu và kiểm thử

- `post_meta` không đọc/ghi bất kỳ key SEO nào.
- Entity type ngoài whitelist bị từ chối; entity không tồn tại/soft-delete không tạo được SEO.
- Unique `(entity_type, entity_id)` hoạt động; upsert lặp lại không tạo duplicate.
- GET trả đồng thời dữ liệu gốc và resolved; fallback độc lập từng field và không ghi ngược DB.
- Validation đúng giới hạn, URL, robots và JSON schema; trim phân biệt null/rỗng/whitespace.
- PUT/DELETE invalidate đúng key; Redis lỗi không làm mất dữ liệu MySQL.
- Xóa entity không để lại SEO record mồ côi.
- Unit test cho registry, Unicode-safe description truncation và SEO Resolver.
- Repository/integration test cho migration, unique constraint, upsert và cleanup lifecycle.
- API test cho `200`, `400`, `401`, `403`, `404`, `409`, `422`, `500`; xác nhận không nhận identity từ body/query.

## 9. Ngoài phạm vi hiện tại

- SEO theo `locale`/đa ngôn ngữ.
- Đổi image URL sang `media_id`.
- Versioning/audit history SEO.
- SERP/Open Graph preview, sitemap và kiểm tra broken canonical URL.
