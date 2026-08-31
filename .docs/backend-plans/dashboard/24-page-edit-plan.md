# KẾ HOẠCH BACKEND: CHỈNH SỬA TRANG TĨNH (PAGE EDIT)

**Dự án:** TechBite  
**Tài liệu nguồn:** `.docs/ideas/dashboard/23-page-edit-idea.md`  
**Phụ thuộc:** `22-page-list-plan.md`, `23-page-create-plan.md`

## 1. Mục tiêu và ranh giới

- Bổ sung Detail và Update cho aggregate `Page`: `GET /api/v1/admin/pages/:page_id` và `PUT /api/v1/admin/pages/:page_id`.
- Cập nhật title, slug, content và status trên cùng một record Page; không cho đổi sang Post, gắn post type/category hoặc sửa `author_id`.
- Kế thừa Page Entity/migration, DDD layers, JWT/RBAC, Redis client và Wire set từ các Page plans trước.
- Page Media và SEO là sub-resource độc lập, không nhận nested Media/SEO vào body PUT này. Các plan tiếp theo sẽ dùng Page ID đã tồn tại.

## 2. Trụ cột 1 — Thiết kế dữ liệu

### 2.1. Tái sử dụng `Page` Entity

Không tạo migration mới. Dùng cùng bảng `pages` và entity:

```go
type Page struct {
    ID        uint           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
    Title     string         `gorm:"column:title;type:varchar(255);not null" json:"title"`
    Slug      string         `gorm:"column:slug;type:varchar(255);not null;uniqueIndex:uq_pages_slug" json:"slug"`
    Content   string         `gorm:"column:content;type:longtext;not null" json:"content"`
    Status    PageStatus     `gorm:"column:status;type:varchar(20);not null;default:'DRAFT';index:idx_pages_status_deleted" json:"status"`
    AuthorID  uint           `gorm:"column:author_id;not null;index:idx_pages_author_id" json:"authorId"`
    CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
    UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
    DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index:idx_pages_status_deleted" json:"-"`
}
```

- `FindByID` luôn dùng GORM scoped query, vì vậy Page soft-deleted được xem là không tồn tại.
- `author_id` là immutable audit field trong hạng mục này; `Update` chỉ cập nhật `title`, `slug`, `content`, `status`, `updated_at`.
- Slug vẫn unique trong namespace `pages`; không tạo global slug registry hay redirect từ slug cũ.

### 2.2. Repository cần bổ sung

```go
FindByID(ctx context.Context, id uint) (*entity.Page, error)
CheckSlugExistsExceptID(ctx context.Context, slug string, id uint) (bool, error)
Update(ctx context.Context, page *entity.Page) error
```

- `CheckSlugExistsExceptID` thêm điều kiện `slug = ? AND id <> ?`, lọc soft-delete theo scope GORM.
- `Update` chỉ dùng column whitelist, không gọi `Save` với object từ client:

```go
Updates(map[string]interface{}{
    "title":   page.Title,
    "slug":    page.Slug,
    "content": page.Content,
    "status":  page.Status,
})
```

- Sau update, query lại Page scoped để trả timestamps mới. Duplicate-key database do race condition được Service map sang lỗi slug conflict.

## 3. Trụ cột 2 — API Contract và Auth Context

### 3.1. Routes và RBAC

Mở rộng route group Page đã được bảo vệ bởi AuthMiddleware + `RoleMiddleware("ADMIN", "STAFF")`:

```go
pages := admin.Group("/pages")
pages.Use(middleware.RoleMiddleware("ADMIN", "STAFF"))
{
    pages.GET("", pageController.GetPages)
    pages.POST("", pageController.CreatePage)
    pages.GET("/:page_id", pageController.GetPageByID)
    pages.PUT("/:page_id", pageController.UpdatePage)
    pages.DELETE("/:page_id", pageController.DeletePage)
}
```

- Auth middleware xác thực JWT, kiểm tra blacklist JTI ở Redis, đặt `userId`/`role` vào Gin Context trước Controller.
- GET/PUT chỉ cho `ADMIN`, `STAFF`; `CUSTOMER` nhận `403`.
- `authorId`, `userId`, `role`, Media và SEO từ body bị bỏ qua vì DTO không định nghĩa các field đó.

### 3.2. Lấy chi tiết Page

**Method & Route:** `GET /api/v1/admin/pages/:page_id`

```go
type GetPageByIDURI struct {
    PageID uint `uri:"page_id" binding:"required,min=1"`
}

type PageResponse struct {
    ID        uint      `json:"id"`
    Title     string    `json:"title"`
    Slug      string    `json:"slug"`
    Content   string    `json:"content"`
    Status    string    `json:"status"`
    AuthorID  uint      `json:"authorId"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}
```

Success `200`:

```json
{
  "data": {
    "id": 1,
    "title": "Giới thiệu",
    "slug": "gioi-thieu",
    "content": "<p>Nội dung trang</p>",
    "status": "PUBLISHED",
    "authorId": 2,
    "createdAt": "2026-08-31T10:00:00Z",
    "updatedAt": "2026-08-31T10:30:00Z"
  }
}
```

- `404 page_not_found` khi ID không tồn tại hoặc đã soft-delete.
- Không preload Page Media hoặc SEO vào response này; frontend tải bằng sub-resource/API SEO riêng để tránh N+1 và giữ aggregate contract rõ ràng.

### 3.3. Cập nhật Page

**Method & Route:** `PUT /api/v1/admin/pages/:page_id`

```go
type UpdatePageURI struct {
    PageID uint `uri:"page_id" binding:"required,min=1"`
}

type UpdatePageRequest struct {
    Title   string `json:"title" binding:"required,max=255"`
    Slug    string `json:"slug" binding:"required,max=255"`
    Content string `json:"content" binding:"required"`
    Status  string `json:"status" binding:"required,oneof=DRAFT PUBLISHED"`
}
```

Request:

```json
{
  "title": "Giới thiệu về TechBite",
  "slug": "gioi-thieu",
  "content": "<p>Nội dung đã cập nhật</p>",
  "status": "PUBLISHED"
}
```

Success `200`:

```json
{
  "data": {
    "id": 1,
    "title": "Giới thiệu về TechBite",
    "slug": "gioi-thieu",
    "content": "<p>Nội dung đã cập nhật</p>",
    "status": "PUBLISHED",
    "authorId": 2,
    "createdAt": "2026-08-31T10:00:00Z",
    "updatedAt": "2026-08-31T11:00:00Z"
  },
  "message": "Cập nhật trang thành công"
}
```

### 3.4. Service rules và error mapping

Mở rộng `PageService`:

```go
GetPageByID(ctx context.Context, id uint) (*dto.PageResponse, error)
UpdatePage(ctx context.Context, id uint, req *dto.UpdatePageRequest) (*dto.PageResponse, error)
```

`UpdatePage` thực hiện theo thứ tự:

1. Tìm Page scoped theo ID; không có trả `ErrPageNotFound`.
2. Trim title, slug, content; reject title/content rỗng sau trim.
3. Validate slug regex `^[a-z0-9]+(?:-[a-z0-9]+)*$` và status enum.
4. Check slug conflict, loại trừ Page hiện tại.
5. Lưu `oldSlug`, cập nhật đúng bốn field được phép rồi repository update.
6. Sau khi DB commit, invalidate cache bằng cả `oldSlug` và `newSlug`; nếu title/content/slug đổi thì xóa SEO resolved cache cho Page.
7. Map entity updated sang `PageResponse` và trả về.

| Lỗi | HTTP | Response |
| --- | --- | --- |
| URI/body/validation sai | 400 | `validation_error` |
| Page không tồn tại/đã xóa | 404 | `page_not_found` |
| Slug trùng | 409 | `slug_conflict` cùng `fields.slug` |
| JWT invalid/revoked | 401 | do AuthMiddleware |
| Role không hợp lệ | 403 | do RoleMiddleware |
| Hạ tầng | 500 | `internal_error`, không lộ SQL/Redis |

Controller chỉ bind URI/body, gọi Service và map domain error; không thực hiện GORM, Redis scan hoặc check slug.

## 4. Trụ cột 3 — Redis cache và token state

### 4.1. Invalidation cache

Page Edit chưa tạo public response cache, nhưng phải chuẩn bị invalidation để không trả dữ liệu stale khi Public Page/SEO được kích hoạt:

```text
admin:pages:list:*
public:pages:slug:{oldSlug}
public:pages:slug:{newSlug}
seo:page:{pageID}
```

- Xóa cả old/new slug để trường hợp đổi slug không giữ URL cũ trong cache.
- Xóa `seo:page:{pageID}` khi title/content/slug thay đổi vì SEO resolver dùng các field này làm fallback.
- Cache invalidation chạy sau update thành công; Redis lỗi chỉ log, không rollback transaction Page.
- Không viết cache ở API admin detail/update và không xóa SEO DB record trong Update.

### 4.2. JWT và session state

- Không bổ sung token/cache auth mới.
- Access token revoked bị chặn bởi `auth:bl:at:{jti}` trước mọi thao tác Page.
- Refresh token rotation và Logout tiếp tục dùng Redis/`SetEx` TTL hiện hữu; API Edit không được tự quản lý refresh token.

## 5. Media, SEO và API documentation

- Page Media (`/admin/pages/:page_id/media`) và Page SEO (`/admin/seo-meta/page/:page_id`) là dependency UI nhưng không thuộc API core Edit. Chúng chỉ chạy khi Page ID hợp lệ và được quy hoạch trong tài liệu riêng.
- Trước khi bật SEO Page, `EntityRegistry` phải resolve `EntityPage` từ bảng `pages` với `deleted_at IS NULL`; hiện trạng `ErrEntityTypeNotConfigured` không đáp ứng UI Page Edit.
- Bổ sung GET/PUT `/api/v1/admin/pages/{page_id}` vào `.docs/api-endpoints.yaml`, cùng BearerAuth, DTO schemas và responses `200`, `400`, `401`, `403`, `404`, `409`.
- Bổ sung PageController methods vào `PageSet` hiện có và regenerate `wire_gen.go`; không tạo DI set mới.

## 6. Kiểm thử và tiêu chí hoàn thành

### Service/Repository

- Get detail trả đúng content/status; Page soft-deleted không được trả.
- Update title/content/status hợp lệ; `authorID` không đổi và only-whitelist columns được ghi.
- Slug sai, title/content whitespace, status sai và slug trùng đều trả domain error đúng.
- Đổi slug invalidates old/new public key; đổi title/content invalidates SEO cache key; Redis lỗi không làm mất update DB.

### API/Middleware

- `ADMIN`/`STAFF` GET/PUT thành công; `CUSTOMER` nhận `403`; thiếu/sai/revoked JWT nhận `401`.
- GET/PUT Page không tồn tại hoặc soft-deleted trả `404`.
- PUT `409` có `fields.slug`; response không chứa deleted fields, JWT claims hoặc SQL details.
- PUT không thể thay `authorId` dù client gửi thêm field; request Media/SEO fields bị bỏ qua, không tạo relation ngoài ý muốn.
- Sau PUT, GET List trả title/slug/status mới và cache keys được invalidated sau commit.
