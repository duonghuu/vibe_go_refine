# KẾ HOẠCH BACKEND: TẠO TRANG TĨNH (PAGE CREATE)

**Dự án:** TechBite  
**Tài liệu nguồn:** `.docs/ideas/dashboard/22-page-create-idea.md`  
**Phụ thuộc:** `.docs/backend-plans/dashboard/22-page-list-plan.md`

## 1. Mục tiêu và phạm vi

- Cung cấp API tạo Page độc lập tại `POST /api/v1/admin/pages` cho Admin/Staff.
- Tái sử dụng aggregate `Page`, bảng `pages`, `PageRepository`, Auth/RBAC, Redis client và Wire set được quy hoạch ở Page List.
- Tạo Page trước; Media gallery và SEO chỉ được liên kết qua API riêng sau khi response trả `pageId`. Không nhận payload Media/SEO trong Create để tránh record liên kết mồ côi.
- Không triển khai Page hierarchy, menu, schedule, public API hoặc Page Edit trong hạng mục này.

## 2. Trụ cột 1 — Thiết kế dữ liệu

### 2.1. Tái sử dụng Entity và migration Page

Không tạo bảng hoặc migration thứ hai. Dùng `entity.Page` và bảng `pages` từ Page List:

```go
type PageStatus string

const (
    PageStatusDraft     PageStatus = "DRAFT"
    PageStatusPublished PageStatus = "PUBLISHED"
)

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

- `author_id` được gán từ JWT Context; không nhận từ JSON body.
- `status` chỉ nhận `DRAFT` hoặc `PUBLISHED`; DB `CHECK` và Service cùng bảo vệ invariant.
- `uq_pages_slug` đảm bảo slug duy nhất trong namespace Page, kể cả bản ghi đã soft-delete. V1 chưa có global slug registry xuyên `pages`/`posts`, do đó không query hay thay đổi Post để kiểm tra slug chéo.

### 2.2. Repository bổ sung

Mở rộng `PageRepository`:

```go
Create(ctx context.Context, page *entity.Page) error
CheckSlugExists(ctx context.Context, slug string) (bool, error)
```

- `Create` dùng `db.WithContext(ctx).Create(page)` để GORM trả `ID`, timestamps sau khi insert.
- `CheckSlugExists` query `Model(&entity.Page{}).Where("slug = ?", slug).Count(&count)`. Mặc định GORM chỉ kiểm tra row chưa bị soft-delete; unique index vẫn là tuyến phòng thủ cuối cùng với row soft-deleted/race condition.
- Service phải map lỗi duplicate-key ở tầng insert thành `ErrPageSlugConflict`, không trả lỗi MySQL thô.

## 3. Trụ cột 2 — API Contract, validation và Auth Context

### 3.1. Route và phân quyền

Mở rộng group Page đã đăng ký trong `main.go`:

```go
pages := admin.Group("/pages")
pages.Use(middleware.RoleMiddleware("ADMIN", "STAFF"))
{
    pages.GET("", pageController.GetPages)
    pages.POST("", pageController.CreatePage)
    pages.DELETE("/:page_id", pageController.DeletePage)
}
```

- `admin` kế thừa `AuthMiddleware(redisClient)`: verify JWT signature, kiểm tra blacklist `auth:bl:at:{jti}`, rồi đặt `userId` và `role` vào Gin Context.
- `CreatePage` lấy `authorID` duy nhất từ `ctx.Get("userId")` sau khi type assert `uint`; body không có `authorId`, `userId` hoặc `role`.
- `CUSTOMER` bị `RoleMiddleware` từ chối với `403` trước khi Controller chạy.

### 3.2. Request và response DTO

Vị trí: `internal/usecases/page/dto/page_dto.go`.

```go
type CreatePageRequest struct {
    Title   string `json:"title" binding:"required,max=255"`
    Slug    string `json:"slug" binding:"required,max=255"`
    Content string `json:"content" binding:"required"`
    Status  string `json:"status" binding:"required,oneof=DRAFT PUBLISHED"`
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

**Endpoint:** `POST /api/v1/admin/pages`

Request:

```json
{
  "title": "Giới thiệu",
  "slug": "gioi-thieu",
  "content": "<p>Nội dung trang</p>",
  "status": "DRAFT"
}
```

Success `201 Created`:

```json
{
  "data": {
    "id": 1,
    "title": "Giới thiệu",
    "slug": "gioi-thieu",
    "content": "<p>Nội dung trang</p>",
    "status": "DRAFT",
    "authorId": 2,
    "createdAt": "2026-08-31T10:00:00Z",
    "updatedAt": "2026-08-31T10:00:00Z"
  }
}
```

### 3.3. Validation và business rules tại Service

`PageService` mở rộng interface:

```go
CreatePage(ctx context.Context, authorID uint, req *dto.CreatePageRequest) (*dto.PageResponse, error)
```

Trình tự:

1. Trim `title`, `slug`, `content`; title/content sau trim không rỗng.
2. Validate slug bằng regex `^[a-z0-9]+(?:-[a-z0-9]+)*$`; không tự sinh slug ở Backend.
3. Parse/validate `status` là `DRAFT` hoặc `PUBLISHED`.
4. Gọi `CheckSlugExists`; nếu tồn tại trả `ErrPageSlugConflict`.
5. Tạo `entity.Page` với `AuthorID=authorID`, gọi repository `Create`.
6. Nếu insert gặp duplicate-key do race condition, map `ErrPageSlugConflict`; lỗi khác wrap có ngữ cảnh và trả về Controller.
7. Chỉ sau khi insert commit thành công mới invalidate cache list; trả `PageResponse` chứa content.

Domain errors:

```go
var (
    ErrPageInvalidInput  = errors.New("dữ liệu trang không hợp lệ")
    ErrPageSlugConflict  = errors.New("slug trang đã tồn tại")
)
```

### 3.4. Error mapping Controller

Controller chỉ `ShouldBindJSON`, đọc identity từ Context, gọi Service và map lỗi:

| Tình huống | HTTP | Body |
| --- | --- | --- |
| JSON/body/field không hợp lệ | 400 | `{ "error": "validation_error", "message": "..." }` |
| JWT thiếu/sai/revoked | 401 | do AuthMiddleware trả trước Controller |
| Không có role phù hợp | 403 | do RoleMiddleware trả trước Controller |
| Slug trùng | 409 | `{ "error": "slug_conflict", "message": "...", "fields": { "slug": "..." } }` |
| Lỗi hệ thống | 500 | `{ "error": "internal_error", "message": "Không thể tạo trang" }` |

Không trả `DeletedAt`, password, claims JWT, SQL error hoặc chi tiết Redis trong response.

### 3.5. DI và tài liệu API

- Giữ `PageSet` từ Page List: `NewPageRepository -> NewPageService -> NewPageController`; chỉ thêm method/constructor dependencies nếu chưa có Redis client.
- `InitializePageController(db, redisClient)` tiếp tục là entrypoint Wire duy nhất; regenerate `wire_gen.go`.
- Bổ sung `POST /api/v1/admin/pages` vào `.docs/api-endpoints.yaml`, gồm BearerAuth, schema request, response `201`, `400`, `401`, `403`, `409`.

## 4. Trụ cột 3 — Redis cache và token state

### 4.1. Cache invalidation

- Admin Page List hiện không cache vì có filter, search và pagination; Create vẫn gọi `invalidateCache(ctx)` cho pattern tương lai `admin:pages:list:*` sau insert thành công.
- Không tạo Page public cache, Page Media cache hoặc SEO cache trong hạng mục này; chúng chỉ được thêm khi API/aggregate tương ứng tồn tại.
- Nếu Redis nil hoặc lỗi scan/delete cache, ghi log và không rollback bản ghi Page đã tạo; database là source of truth.

### 4.2. Token/Redis hiện hữu

- Không tạo token/session key mới cho Page.
- Page Create kế thừa blacklist key `auth:bl:at:{jti}`, refresh session key `auth:rf:{userId}:{jti}`, Refresh Token Rotation và Logout `SetEx` TTL hiện tại.
- Khi access token đã bị revoke, middleware chặn request trước khi Service tạo Page; không có side effect database/cache.

## 5. Kiểm thử và tiêu chí hoàn thành

### Test Service/Repository

- Tạo `DRAFT` và `PUBLISHED` gán đúng `authorID`, timestamps, ID và response.
- Reject title/content chỉ có whitespace, status ngoài enum, slug sai regex và slug vượt giới hạn.
- Slug trùng trả domain error `ErrPageSlugConflict`; mô phỏng duplicate-key tại insert cũng map cùng error.
- Không có Media/SEO row hay request nào được tạo từ Create Page.
- Invalidation chỉ chạy sau Create thành công; Redis lỗi không làm Create thất bại.

### Test API/Middleware

- `ADMIN`/`STAFF` tạo được Page và nhận `201 {data}`.
- `CUSTOMER` nhận `403`; thiếu/sai/revoked token nhận `401`; `authorId` giả trong body bị Gin bỏ qua và không thay đổi author từ JWT.
- Payload không hợp lệ nhận `400`; slug trùng nhận `409` kèm `fields.slug`; lỗi nội bộ không lộ SQL.
- Page vừa tạo xuất hiện qua `GET /api/v1/admin/pages` với đúng status.

## 6. Ghi chú triển khai sau Create

- Frontend chỉ chuyển tới `/pages/edit/:id` sau success; không upload trước khi Page tồn tại.
- Page Media và SEO sẽ có plan/API riêng; khi triển khai SEO cần thay `EntityPage` từ `ErrEntityTypeNotConfigured` sang Page resolver.
- Page Edit sẽ dùng lại `PageResponse` và không được đổi `authorId` do client gửi lên.
