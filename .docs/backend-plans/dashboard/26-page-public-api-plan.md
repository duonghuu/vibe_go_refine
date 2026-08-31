# KẾ HOẠCH BACKEND: PUBLIC PAGE API

**Dự án:** TechBite  
**Tài liệu nguồn:** `.docs/ideas/dashboard/25-page-public-api-idea.md`  
**Phạm vi:** API read-only công khai trả Page đã xuất bản theo slug, gồm content, Page Media đã sắp xếp và SEO đã resolve.

## 1. Hiện trạng, phụ thuộc và nguyên tắc

- API này là query-side read model, không tạo aggregate/bảng mới. Nó tái sử dụng `pages`, `page_media`, `media`, `seo_meta` từ các kế hoạch Page List/Create/Edit, Page Media và SEO Meta.
- Các module/migration Page và Page Media hiện mới được quy hoạch, chưa tồn tại trong code. Public API chỉ được bật sau khi các dependency đó hoạt động và `EntityRegistry` hỗ trợ `EntityPage`.
- Không tái dùng `Post`/`post_media`, không gọi API admin nội bộ, và không để Webview ghép Page, Media, SEO qua nhiều HTTP request.
- `GET /api/v1/pages/:slug` không cần Access Token, chỉ trả Page `PUBLISHED` có `deleted_at IS NULL`. DRAFT, soft-deleted và slug không tồn tại cùng trả `404 PAGE_NOT_FOUND` để không tiết lộ trạng thái nội dung.
- Handler mỏng; query, status visibility, map Media, SEO fallback, cache-aside và invalidation nằm trong Public Page service/repository. Không thêm công nghệ ngoài Go, Gin, GORM, Wire, MySQL, Redis.

## 2. Trụ cột 1 — Dữ liệu và read model

### 2.1. Entity/migration

Không tạo table hoặc migration mới. Reuse đúng các persistent entity đã quy hoạch:

```go
// internal/domain/page/entity/page.go
type Page struct {
    ID        uint       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
    Title     string     `gorm:"column:title;type:varchar(255);not null" json:"title"`
    Slug      string     `gorm:"column:slug;type:varchar(255);not null;uniqueIndex:uq_pages_slug" json:"slug"`
    Content   string     `gorm:"column:content;type:longtext;not null" json:"content"`
    Status    PageStatus `gorm:"column:status;type:varchar(20);not null" json:"status"`
    AuthorID  uint       `gorm:"column:author_id;not null" json:"-"`
    DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

// internal/domain/page/entity/page_media.go
type PageMedia struct {
    PageID     uint                `gorm:"column:page_id;not null;index:idx_page_media_page_collection_sort,priority:1" json:"-"`
    MediaID    uint                `gorm:"column:media_id;not null" json:"-"`
    Collection PageMediaCollection `gorm:"column:collection;type:varchar(20);not null;index:idx_page_media_page_collection_sort,priority:2" json:"-"`
    SortOrder  int                 `gorm:"column:sort_order;not null;index:idx_page_media_page_collection_sort,priority:3" json:"-"`
    DeletedAt  gorm.DeletedAt      `gorm:"column:deleted_at;index:idx_page_media_page_collection_sort,priority:4" json:"-"`
}
```

- Public query luôn lọc `pages.status = 'PUBLISHED'`, Page/Pivot/Media `deleted_at IS NULL`, `collection IN ('thumbnail', 'gallery')`.
- Index `uq_pages_slug` bảo đảm lookup slug một bản ghi; `idx_page_media_page_collection_sort(page_id, collection, sort_order, deleted_at)` hỗ trợ media query đã sắp xếp. Không thêm index chỉ cho một GET cho đến khi `EXPLAIN` chứng minh cần thiết.
- Không trả `author_id`, Media `owner_id`, `status`, `mime_type`, size, original file name, internal database timestamps, refresh token hay dữ liệu admin khác.

### 2.2. Public projection và repository

Tạo projection typed không map sang bảng mới:

```go
type PublicPageProjection struct {
    ID        uint      `gorm:"column:id" json:"-"`
    Title     string    `gorm:"column:title" json:"-"`
    Slug      string    `gorm:"column:slug" json:"-"`
    Content   string    `gorm:"column:content" json:"-"`
    UpdatedAt time.Time `gorm:"column:updated_at" json:"-"`
    SEO       *seoEntity.SEOMeta
    Thumbnail *PublicPageMediaProjection
    Gallery   []PublicPageMediaProjection
}

type PublicPageMediaProjection struct {
    ID           uint   `gorm:"column:id" json:"-"`
    OriginalURL  string `gorm:"column:original_url" json:"-"`
    ThumbnailURL string `gorm:"column:thumbnail_url" json:"-"`
    MediumURL    string `gorm:"column:medium_url" json:"-"`
    SortOrder    int    `gorm:"column:sort_order" json:"-"`
}
```

Tạo `internal/infrastructure/repository/public_page_repo.go`:

```go
type PublicPageRepository interface {
    FindPublishedBySlug(ctx context.Context, slug string) (*PublicPageProjection, error)
}
```

`FindPublishedBySlug` dùng số query cố định, không N+1:

1. Query Page published theo slug, `LEFT JOIN seo_meta` với `entity_type = 'page'`.
2. Nếu Page tồn tại, query `page_media JOIN media` một lần, filter pivot/media active, sort `collection ASC, sort_order ASC, page_media.id ASC`; repository tách thumbnail/gallery trong projection.

Không đưa SEO record hay Media links của Page khác vào response. `gorm.ErrRecordNotFound` được map thành `ErrPublicPageNotFound` ở service.

### 2.3. SEO resolver dùng chung

Refactor pure resolver từ `usecases/seometa/service/seo_meta_service.go` sang file dùng chung, ví dụ `internal/usecases/seometa/service/resolver.go`:

```go
type ResolvedSEO struct {
    ResolvedTitle        string
    ResolvedDescription  string
    ResolvedCanonicalURL string
    OGTitle              string
    OGDescription        string
    OGImage              string
    OGType               string
    TwitterTitle         string
    TwitterDescription   string
    TwitterImage         string
    TwitterCard          string
    Robots               string
    SchemaJSON           json.RawMessage
}

func ResolvePublicPageSEO(meta *seoEntity.SEOMeta, ref PageSEOReference, publicSiteURL string) ResolvedSEO
```

- `PageSEOReference` lấy title, slug, content và thumbnail từ projection; không query lại Page trong resolver.
- Fallback title lấy Page title; description strip HTML, trim và cắt an toàn theo rune; image dùng thumbnail Page; robots mặc định `index,follow`.
- Canonical dùng `seo_meta.canonical_url` nếu có, ngược lại tạo từ `PUBLIC_SITE_URL` + slug. `PUBLIC_SITE_URL` là config bắt buộc, phải parse/validate `http` hoặc `https` lúc khởi động; không hard-code domain hoặc trả slug tương đối như resolver hiện hữu.
- Public DTO chỉ trả giá trị **resolved**. Raw custom SEO fields không cần lộ ra ngoài API public.
- Mở rộng `EntityRegistry.Resolve` cho `EntityPage`: query `pages` active và thumbnail Page bằng left join. Đây là prerequisite để Admin SEO Page và public fallback dùng cùng semantics; registry vẫn không trả DRAFT qua public route vì service public tự lọc PUBLISHED.

## 3. Trụ cột 2 — API contract và DDD/CQRS

### 3.1. Module, DI và route

```text
apps/backend/internal/
├── controller/public_page_controller.go
├── infrastructure/repository/public_page_repo.go
└── usecases/publicpage/
    ├── dto/public_page_dto.go
    └── service/public_page_service.go
```

Thêm `PublicPageSet` vào `internal/di/wire.go`, inject `*gorm.DB`, `*redis.Client` và `PublicSiteURLConfig`, rồi regenerate `wire_gen.go`:

```go
var PublicPageSet = wire.NewSet(
    repository.NewPublicPageRepository,
    publicpageService.NewPublicPageService,
    controller.NewPublicPageController,
)
```

Đăng ký public route **ngoài** `admin` group, không gắn `AuthMiddleware` hoặc `RoleMiddleware`:

```go
pagesPublic := api.Group("/pages")
{
    pagesPublic.GET("/:slug", publicPageController.GetBySlug)
}
```

Không có body, userId hoặc role trong API này. Các admin write Page/Page Media/SEO vẫn bắt buộc kế thừa JWT middleware như các kế hoạch trước.

### 3.2. Request binding và response DTO

```go
type GetPublicPageURI struct {
    Slug string `uri:"slug" binding:"required,max=255"`
}

type PublicPageMediaResponse struct {
    ID           uint   `json:"id"`
    OriginalURL  string `json:"originalUrl"`
    ThumbnailURL string `json:"thumbnailUrl,omitempty"`
    MediumURL    string `json:"mediumUrl,omitempty"`
}

type PublicPageSEOResponse struct {
    ResolvedTitle        string          `json:"resolvedTitle"`
    ResolvedDescription  string          `json:"resolvedDescription"`
    ResolvedCanonicalURL string          `json:"resolvedCanonicalUrl"`
    OGTitle              string          `json:"ogTitle,omitempty"`
    OGDescription        string          `json:"ogDescription,omitempty"`
    OGImage              string          `json:"ogImage,omitempty"`
    OGType               string          `json:"ogType,omitempty"`
    TwitterTitle         string          `json:"twitterTitle,omitempty"`
    TwitterDescription   string          `json:"twitterDescription,omitempty"`
    TwitterImage         string          `json:"twitterImage,omitempty"`
    TwitterCard          string          `json:"twitterCard,omitempty"`
    Robots               string          `json:"robots"`
    SchemaJSON           json.RawMessage `json:"schemaJson,omitempty"`
}

type PublicPageResponse struct {
    ID        uint                      `json:"id"`
    Title     string                    `json:"title"`
    Slug      string                    `json:"slug"`
    Content   string                    `json:"content"`
    UpdatedAt time.Time                 `json:"updatedAt"`
    Thumbnail *PublicPageMediaResponse  `json:"thumbnail"`
    Gallery   []PublicPageMediaResponse `json:"gallery"`
    SEO       PublicPageSEOResponse     `json:"seo"`
}

type GetPublicPageResponse struct {
    Data PublicPageResponse `json:"data"`
}
```

`gallery` luôn là `[]`, không phải `null`; `thumbnail` là `null` nếu Page chưa gắn thumbnail. `schemaJson` chỉ được trả khi JSON-LD hợp lệ theo validation của SEO service.

### 3.3. GET Page theo slug

```text
GET /api/v1/pages/:slug
```

Service flow:

1. Trim slug rồi validate regex lowercase `^[a-z0-9]+(?:-[a-z0-9]+)*$`; không tự lowercase hoặc redirect ở API. Input sai trả `400 INVALID_SLUG` trước cache/DB.
2. Read `public:pages:slug:{slug}`; cache hit unmarshal `GetPublicPageResponse` và trả ngay.
3. Cache miss gọi repository, chỉ nhận Page published active. Không có record trả `ErrPublicPageNotFound`.
4. Map thumbnail/gallery theo `sort_order`, gọi resolver chung để build SEO resolved, serialize DTO rồi `SetEx` cache.
5. Trả `200 { "data": ... }`.

Ví dụ response:

```json
{
  "data": {
    "id": 1,
    "title": "Giới thiệu",
    "slug": "gioi-thieu",
    "content": "<p>Nội dung trang</p>",
    "updatedAt": "2026-08-31T10:30:00Z",
    "thumbnail": null,
    "gallery": [],
    "seo": {
      "resolvedTitle": "Giới thiệu",
      "resolvedDescription": "Nội dung trang",
      "resolvedCanonicalUrl": "https://example.com/gioi-thieu",
      "robots": "index,follow"
    }
  }
}
```

### 3.4. Error contract

```go
type PublicErrorResponse struct {
    Error string `json:"error"`
    Code  string `json:"code"`
}
```

| Điều kiện | HTTP | Code | Response |
| --- | ---: | --- | --- |
| Slug sai format | 400 | `INVALID_SLUG` | lỗi validation chung |
| DRAFT, soft-deleted hoặc không có slug | 404 | `PAGE_NOT_FOUND` | cùng message không tiết lộ trạng thái |
| DB/Redis/config unexpected | 500 | `INTERNAL_ERROR` | không trả SQL/cache detail |

Không trả `401`/`403` vì route là public; không cache `400`, `404`, `500` ở v1. Nếu cần rate limit public sau này, thêm Gin middleware ở API gateway/route layer, không đưa logic vào service.

## 4. Trụ cột 3 — Redis cache và token state

### 4.1. Cache-aside và invalidation

```text
public:pages:slug:{normalized_slug}
```

- Cache lưu toàn bộ `GetPublicPageResponse` đã resolved, TTL 15 phút qua `SetEx`; không cache GORM entity, raw SEO config, token hoặc DTO admin.
- Redis read lỗi: log và query MySQL; Redis write lỗi: vẫn trả `200` nếu DB/read model thành công.
- Invalidate chỉ sau DB transaction commit:
  - Page Create/Update/Delete: xóa slug cũ và slug mới; status DRAFT ↔ PUBLISHED luôn xóa key.
  - Page Media POST/PUT/DELETE: lấy slug từ Page row đã lock rồi xóa key.
  - SEO Meta Page PUT/DELETE: registry resolve Page slug rồi xóa key.
- Đồng thời xóa `seo:page:{pageID}` khi Page title/content/slug/thumbnail đổi, vì SEO fallback phụ thuộc các field đó. Không dùng `KEYS`; các key cần xóa đều xác định trực tiếp.
- Cache miss concurrent có thể tạo truy vấn lặp nhưng không ảnh hưởng correctness; không cache 404 để Page vừa publish xuất hiện ngay sau invalidation.

### 4.2. Token/Redis security

- Public GET không parse/lưu JWT và không dùng blacklist key. Cache key chỉ dùng slug đã validate, không lấy raw header/query để tránh cache-key poisoning.
- Admin mutations giữ `AuthMiddleware`: verify signature trước rồi kiểm tra `auth:bl:at:{jti}`. `userId`/`role` chỉ lấy từ verified claims.
- Logout tiếp tục `SetEx` blacklist với TTL còn lại. Refresh token rotation phải atomically revoke JTI cũ, tạo/lưu JTI mới; replay theo policy Auth phải thu hồi phiên liên quan. Public Page service không tham gia token flow.

## 5. Kiểm thử và tiêu chí nghiệm thu

- Page PUBLISHED trả đúng Page/content, thumbnail nullable, gallery `[]` và order ổn định; không có N+1 query Media.
- DRAFT, Page soft-deleted và slug không tồn tại cùng trả `404 PAGE_NOT_FOUND` với body không phân biệt.
- Slug rỗng, uppercase, ký tự nguy hiểm hoặc gap hyphen bị từ chối `400`; input hợp lệ dùng đúng format slug Page admin.
- SEO custom/fallback trả đúng resolved title, HTML-stripped description, canonical absolute URL, thumbnail fallback, robots và JSON-LD hợp lệ; raw SEO fields không bị trả.
- Cache hit không gọi MySQL; cache miss write đúng full response; Redis unavailable vẫn đọc API từ MySQL.
- Update title/content/status/slug, sync thumbnail/gallery hoặc upsert/delete SEO đều xóa đúng public key **sau commit**; slug cũ không còn stale sau rename.
- Route không yêu cầu Authorization; API public không trả author/owner/status Media/MIME/size/tokens hoặc lỗi SQL.
- Wire generated file được regenerate, OpenAPI `.docs/api-endpoints.yaml` có endpoint/schemas `200/400/404/500`, và Webview có thể render từ đúng một request API.

## 6. Trình tự triển khai

1. Hoàn thành thực thi `pages`, Page Media, SEO Meta Page và `EntityRegistry.EntityPage` theo các plan trước.
2. Bổ sung `PUBLIC_SITE_URL` validated config và tách SEO resolver thuần dùng chung.
3. Tạo public repository projection/query, DTO và service cache-aside.
4. Thêm controller, `PublicPageSet`, Wire injector và route public.
5. Bổ sung invalidation hooks cho Page/Page Media/SEO sau commit.
6. Cập nhật OpenAPI, unit/integration/API tests và kiểm tra query/cache bằng log hoặc metrics hiện có.
