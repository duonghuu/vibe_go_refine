# KẾ HOẠCH BACKEND: CHỈNH SỬA BÀI VIẾT (ADMIN EDIT POST)

**Dự án:** TechBite  
**Module:** Post Management  
**Tài liệu tham chiếu:** `.docs/ideas/dashboard/15-post-create-idea.md`, `.docs/backend-plans/dashboard/15-post-create-plan.md`, `.docs/backend-plans/dashboard/16-post-list-plan.md`, `.docs/backend-plans/dashboard/13-post-media-plan.md`, `.docs/backend-plans/dashboard/14-post-meta-plan.md`

## 1. Bối cảnh và mục tiêu

- Bổ sung khả năng lấy chi tiết và cập nhật bài viết từ Dashboard Admin.
- Cho phép cập nhật `title`, `slug`, `content` và `categoryId`.
- `typeCode` được xác định từ bản ghi hiện tại, không cho client chuyển bài viết sang Post Type khác trong API Update.
- Giữ nguyên `authorId` và các audit field; không nhận các trường này từ request body.
- Tái sử dụng Entity `Post`, Repository, Service, Controller, Auth Middleware và Redis client hiện có.
- Post Media và Post Meta là các aggregate/module riêng. Không cập nhật ngầm chúng trong transaction Post Update; nếu cần đồng bộ, gọi các API media/meta riêng sau khi Post Update thành công.

### 1.1. Hiện trạng cần bổ sung

Module Post hiện đã có `POST /api/v1/admin/posts`, `GET /api/v1/admin/posts` và Delete, nhưng chưa có:

- `GET /api/v1/admin/posts/:id`.
- `PUT /api/v1/admin/posts/:id`.
- Repository methods tìm theo ID, kiểm tra slug ngoại trừ chính bản ghi đang sửa và cập nhật entity.
- DTO `UpdatePostRequest` và mapping lỗi nghiệp vụ sang HTTP status phù hợp.

## 2. Kiến trúc triển khai DDD/Clean Architecture

Luồng request:

`Gin Route -> AuthMiddleware -> RoleMiddleware -> PostController -> PostService -> PostRepository -> MySQL`

- **Controller:** Parse `id`, bind JSON, kiểm tra lỗi binding, lấy `userId`/`role` từ Gin Context nếu cần audit/authorization, gọi Service và trả response.
- **Service/Use Case:** Validate dữ liệu, kiểm tra Post tồn tại, bảo toàn `typeCode`/`authorId`, kiểm tra slug unique, kiểm tra category cùng `typeCode`, thực hiện update và invalidate cache.
- **Repository:** Chỉ thực hiện truy vấn GORM, không chứa business rule.
- **DI:** Bổ sung method vào interface hiện tại; không tạo graph mới. `wire.go` tiếp tục cung cấp `PostRepository -> PostService -> PostController`, sau đó regenerate `wire_gen.go`.

### 2.1. Phân quyền

- Route phải kế thừa `AuthMiddleware(redisClient)` từ group `/api/v1/admin`.
- Bổ sung `RoleMiddleware("ADMIN", "STAFF")` cho route Update/Get detail nếu chính sách Admin hiện tại giới hạn Content Management cho hai role này.
- Nếu hệ thống chính thức thêm role `EDITOR`, thêm role đó vào policy, không lấy role từ body/query.
- `userId` và `role` phải được lấy từ JWT claims đã được middleware xác thực và set vào context. Không tin tưởng `authorId`, `userId` hoặc `role` từ client.

## 3. Trụ cột 1: Thiết kế dữ liệu

### 3.1. GORM Entity tái sử dụng

Không tạo bảng mới. Cập nhật/duy trì Entity tại `apps/backend/internal/domain/post/entity/post.go`:

```go
package entity

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID         uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	TypeCode   string         `gorm:"type:varchar(50);not null;index" json:"typeCode"`
	Title      string         `gorm:"type:varchar(255);not null" json:"title"`
	Slug       string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Content    string         `gorm:"type:longtext;not null" json:"content"`
	AuthorID   uint           `gorm:"not null;index" json:"authorId"`
	CategoryID *uint          `gorm:"index" json:"categoryId"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Post) TableName() string { return "posts" }
```

### 3.2. Constraints và tính toàn vẹn

- `ID`: khóa chính tự tăng.
- `TypeCode`: bắt buộc, foreign key logic đến `post_types.code`, có index; Update không được thay đổi.
- `Title`: bắt buộc, tối đa 255 ký tự.
- `Slug`: bắt buộc, unique trên toàn bảng theo migration hiện tại (`idx_posts_slug`). Khi kiểm tra trùng phải loại trừ `id` đang chỉnh sửa.
- `Content`: `longtext`, không được rỗng sau khi chuẩn hóa whitespace/HTML rỗng ở tầng Service.
- `AuthorID`: bắt buộc, không thay đổi trong Edit.
- `CategoryID`: nullable, có index và phải tham chiếu Post Category có cùng `TypeCode`; nếu null thì bỏ phân loại.
- `DeletedAt`: soft delete; bản ghi đã soft-delete không được trả về ở Get detail và không được Update.
- Không cần migration cho chức năng Edit nếu schema hiện tại đã khớp `posts` và `category_id`.

### 3.3. Transaction và cạnh tranh dữ liệu

- Update Post chính thực hiện trong transaction DB khi cần kết hợp kiểm tra tồn tại, category và update.
- Unique index MySQL vẫn là lớp bảo vệ cuối cùng cho slug. Nếu xảy ra duplicate key do race condition, map về `409 Conflict`.
- Không dùng `Save` với toàn bộ request body. Chỉ update whitelist fields để tránh ghi đè `typeCode`, `authorId`, `createdAt`, `deletedAt`.
- Media/Meta không nằm trong transaction này; dùng endpoint sync riêng theo đúng module đã thiết kế.

## 4. Trụ cột 2: API Contract

Base route: `/api/v1/admin/posts`  
Auth: `AuthMiddleware` bắt buộc; tất cả request phải kiểm tra JWT signature trước, sau đó kiểm tra blacklist JTI trong Redis.

### 4.1. Lấy chi tiết bài viết

**Method & Route:** `GET /api/v1/admin/posts/:id`

**URI binding:**

```go
type GetPostByIDURI struct {
	ID uint `uri:"id" binding:"required,min=1"`
}
```

**Response 200:**

```go
type PostResponse struct {
	ID         uint      `json:"id"`
	TypeCode   string    `json:"typeCode"`
	Title      string    `json:"title"`
	Slug       string    `json:"slug"`
	Content    string    `json:"content"`
	AuthorID   uint      `json:"authorId"`
	CategoryID *uint     `json:"categoryId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type GetPostResponse struct {
	Data PostResponse `json:"data"`
}
```

**Business rules:**

- Chỉ tìm bản ghi chưa soft-delete.
- Có thể preload/resolve thông tin category nếu cần, nhưng response tối thiểu phải có `categoryId` để Refine đổ form.
- Không trả `DeletedAt` hoặc dữ liệu nội bộ không cần thiết.

**HTTP errors:**

- `400 Bad Request`: ID không hợp lệ.
- `401 Unauthorized`: thiếu/sai/hết hạn JWT hoặc JTI bị blacklist.
- `403 Forbidden`: role không được phép truy cập Content Management.
- `404 Not Found`: không tồn tại hoặc đã soft-delete.
- `500 Internal Server Error`: lỗi truy vấn không dự đoán được.

### 4.2. Cập nhật bài viết

**Method & Route:** `PUT /api/v1/admin/posts/:id`

**URI binding:**

```go
type UpdatePostURI struct {
	ID uint `uri:"id" binding:"required,min=1"`
}
```

**Request binding:**

```go
type UpdatePostRequest struct {
	Title      string `json:"title" binding:"required,max=255"`
	Slug       string `json:"slug" binding:"required,max=255"`
	Content    string `json:"content" binding:"required"`
	CategoryID *uint  `json:"categoryId"`
}
```

**Validation bổ sung tại Service:**

- Trim/chuẩn hóa các field string trước khi kiểm tra; không chấp nhận title/content chỉ gồm whitespace.
- Slug phải khớp `^[a-z0-9-]+$` và không chứa dấu cách.
- `CategoryID == nil` là hợp lệ; nếu có giá trị thì category phải tồn tại, chưa bị xóa mềm và `category.type_code == post.type_code`.
- Không nhận `TypeCode`, `AuthorID`, `CreatedAt`, `UpdatedAt`, `DeletedAt` từ request.

**Response 200:**

```go
type UpdatePostResponse struct {
	Data    PostResponse `json:"data"`
	Message string       `json:"message"`
}
```

`Message` mặc định: `Cập nhật bài viết thành công`.

**HTTP errors:**

- `400 Bad Request`: JSON/URI không hợp lệ, field không đạt validation, category không cùng loại bài viết.
- `401 Unauthorized`: JWT không hợp lệ, hết hạn hoặc bị revoke.
- `403 Forbidden`: role không có quyền Edit Post.
- `404 Not Found`: Post không tồn tại/đã bị soft-delete hoặc category được tham chiếu không tồn tại.
- `409 Conflict`: slug đã được dùng bởi Post khác.
- `500 Internal Server Error`: lỗi DB/Redis hoặc lỗi hệ thống.

### 4.3. Chuẩn lỗi đề xuất

Các endpoint mới nên trả cấu trúc nhất quán để Data Provider hiển thị được lỗi:

```go
type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message,omitempty"`
	Fields  map[string]string `json:"fields,omitempty"`
}
```

- Lỗi slug trùng: `fields.slug = "slug đã tồn tại, vui lòng chọn slug khác"`.
- Không trả stack trace, SQL query, JWT, refresh token hoặc thông tin nhạy cảm.
- Dùng sentinel/domain errors như `ErrPostNotFound`, `ErrPostSlugConflict`, `ErrPostCategoryTypeMismatch` để Controller map status, không dò chuỗi lỗi.

## 5. Service và Repository contract

### 5.1. Repository interface

Mở rộng `apps/backend/internal/infrastructure/repository/post_repo.go`:

```go
type PostRepository interface {
	Create(ctx context.Context, post *entity.Post) error
	FindByID(ctx context.Context, id uint) (*entity.Post, error)
	CheckSlugExistsExceptID(ctx context.Context, slug string, id uint) (bool, error)
	FindAndCount(ctx context.Context, title, typeCode string, offset, limit int, sort string) ([]entity.Post, int64, error)
	Update(ctx context.Context, post *entity.Post) error
	Delete(ctx context.Context, id uint) error
}
```

Repository implementation:

- `FindByID`: `Where("id = ?", id).First(&post)`; để GORM tự loại soft-deleted record.
- `CheckSlugExistsExceptID`: `Where("slug = ? AND id <> ?", slug, id).Count(...)`.
- `Update`: chỉ nhận entity đã được Service whitelist/validate; có thể dùng `Updates(map[string]interface{}{...})` với các field được phép.
- Không ghép trực tiếp `sort` từ client vào SQL nếu chưa whitelist cột/order ở Service; chỉ cho phép `id`, `title`, `created_at`, `updated_at` và `asc/desc`.

### 5.2. Service interface

Mở rộng `apps/backend/internal/usecases/post/service/post_service.go`:

```go
type PostService interface {
	CreatePost(ctx context.Context, authorID uint, req *dto.CreatePostRequest) (*dto.PostResponse, error)
	GetPostByID(ctx context.Context, id uint) (*dto.PostResponse, error)
	GetPosts(ctx context.Context, query dto.GetPostsQuery) (*dto.PaginatedPostResponse, error)
	UpdatePost(ctx context.Context, id uint, req *dto.UpdatePostRequest) (*dto.PostResponse, error)
	DeletePost(ctx context.Context, id uint) error
}
```

`UpdatePost` flow:

1. Load Post hiện tại; trả `ErrPostNotFound` nếu không có.
2. Validate request và slug format.
3. Kiểm tra slug trùng, loại trừ chính `id`.
4. Nếu có `CategoryID`, kiểm tra category tồn tại và cùng `TypeCode` với Post.
5. Chỉ gán `Title`, `Slug`, `Content`, `CategoryID`; giữ nguyên `TypeCode`, `AuthorID`.
6. Update trong transaction.
7. Invalidate cache Post detail, list và các cache liên quan đến type cũ.
8. Map entity đã cập nhật thành `PostResponse`.

## 6. Controller và Router

### 6.1. Controller methods

Bổ sung:

- `GetPostByID(ctx *gin.Context)`: bind URI, gọi `GetPostByID`, map `gorm.ErrRecordNotFound`/domain error thành `404`.
- `UpdatePost(ctx *gin.Context)`: bind URI + JSON, gọi `UpdatePost`, không đọc `typeCode`/`authorId` từ body.

Controller không được kiểm tra slug, category type hoặc tự viết transaction.

### 6.2. Router

Trong nhóm `admin := api.Group("/admin")` hiện đã có `AuthMiddleware`:

```go
posts := admin.Group("/posts")
posts.Use(middleware.RoleMiddleware("ADMIN", "STAFF"))
{
	posts.GET("", postController.GetPosts)
	posts.GET("/:id", postController.GetPostByID)
	posts.POST("", postController.CreatePost)
	posts.PUT("/:id", postController.UpdatePost)
	posts.DELETE("/:post_id", postController.DeletePost)
}
```

Nếu chưa muốn thay đổi policy cho List/Create/Delete, có thể gắn `RoleMiddleware` riêng cho `GET /:id` và `PUT /:id`, nhưng phải có policy rõ ràng và nhất quán với Content Management.

## 7. Trụ cột 3: Redis, cache và token state

### 7.1. Cache Post

Dashboard ưu tiên realtime, vì vậy cache detail/list là tùy chọn. Nếu bật cache, dùng namespace rõ ràng:

- Detail: `admin:posts:detail:{id}`.
- List: `admin:posts:list:type:{typeCode}:page:{page}:size:{pageSize}:search:{hash}`.
- Category tree: tái sử dụng cache của Post Category nếu module đó đã bật.

Khuyến nghị TTL:

- Detail/list admin: 1-5 phút.
- Không lưu request body hoặc nội dung nhạy cảm ngoài response Post cần thiết.

Sau `PUT` thành công, bắt buộc invalidate:

- `admin:posts:detail:{id}`.
- Toàn bộ list cache `admin:posts:list:*` vì slug/title/category/type filter có thể thay đổi.
- Nếu tương lai cho phép đổi `typeCode`, invalidate cả type cũ và type mới; trong plan này `typeCode` immutable nên chỉ cần type hiện tại.
- Cache public/storefront liên quan đến slug hoặc post content nếu tồn tại.

Dùng `SCAN` theo batch hoặc quản lý Set key; không dùng `KEYS` trên production vì có thể block Redis. Lỗi invalidate cần được log/monitor; quyết định trả lỗi mutation hay chỉ retry phải được thống nhất ở infrastructure policy, nhưng DB update không được rollback chỉ vì cache miss.

### 7.2. Auth Middleware và Blacklist

Không tạo auth flow mới cho Edit. Tái sử dụng `AuthMiddleware` hiện tại:

1. Parse và verify chữ ký JWT.
2. Lấy `jti`, kiểm tra `auth:bl:at:{jti}` trong Redis.
3. Nếu tồn tại, abort `401 Unauthorized` trước khi vào Controller.
4. Set `userId` và `role` vào Gin Context.

Khi `/logout` được gọi, access token/JTI phải được blacklist bằng `SetEx` với TTL còn lại của JWT. Không lưu toàn bộ token lâu dài trong Redis; refresh session tiếp tục dùng key `auth:rf:{userId}:{jti}`.

### 7.3. Refresh Token Rotation

Chức năng Edit không thay đổi `/refresh-token`, nhưng plan triển khai phải bảo toàn quy tắc hệ thống:

- Kiểm tra refresh JTI tồn tại.
- Thu hồi JTI cũ và tạo cặp token/JTI mới trong luồng atomic/transaction phù hợp.
- Khi phát hiện replay, xóa toàn bộ refresh session của user và yêu cầu đăng nhập lại.
- Refresh token tiếp tục nằm trong HttpOnly/Secure/SameSite cookie theo auth contract hiện tại.

## 8. DDD/CQRS và tính nhất quán dữ liệu

- Command `UpdatePost` chịu trách nhiệm thay đổi Post và phát invalidation event nội bộ qua service/infrastructure hiện có; không để Controller gọi Redis trực tiếp.
- Query `GetPostByID` chỉ đọc Post chưa soft-delete và trả DTO response riêng, không expose trực tiếp GORM entity.
- Repository không gọi Redis và không chứa authorization.
- Nếu sau này cần audit lịch sử thay đổi, bổ sung module audit riêng; không nhồi audit log vào `PostMeta`.

## 9. Kiểm thử và tiêu chí nghiệm thu

### Unit/Service tests

- Update title/content/category thành công.
- Đổi slug thành slug mới thành công.
- Giữ nguyên slug hiện tại không bị báo conflict.
- Slug trùng Post khác trả `ErrPostSlugConflict`.
- Post không tồn tại/đã soft-delete trả `ErrPostNotFound`.
- Category null hợp lệ.
- Category khác `TypeCode` bị từ chối.
- Request không thể thay đổi `TypeCode` hoặc `AuthorID`.
- Cache invalidation được gọi sau update thành công.

### Integration/API tests

- `GET /api/v1/admin/posts/:id` trả đầy đủ Content và CategoryID.
- `PUT /api/v1/admin/posts/:id` trả `200` và payload `data` mới.
- Thiếu JWT, JWT hết hạn hoặc JTI blacklist trả `401`.
- Role không hợp lệ trả `403`.
- ID sai trả `400`; ID không tồn tại trả `404`.
- Payload sai validation trả `400` với `fields` nếu xác định được field.
- Duplicate slug trả `409`.
- Update xong, list/detail không trả cache cũ.

### Checklist triển khai

- [ ] Thêm `UpdatePostRequest`, `GetPostByIDURI`, `UpdatePostURI` và response/error DTO.
- [ ] Thêm `FindByID`, `CheckSlugExistsExceptID`, `Update` vào Post Repository.
- [ ] Thêm `GetPostByID` và `UpdatePost` vào Post Service.
- [ ] Thêm domain errors và mapping HTTP status tại Controller.
- [ ] Thêm `GET /:id` và `PUT /:id` vào router Admin Posts.
- [ ] Áp dụng RoleMiddleware theo policy Content Management.
- [ ] Invalidate detail/list cache bằng SCAN hoặc Set key sau mutation.
- [ ] Cập nhật `.docs/api-endpoints.yaml` cho GET/PUT `/api/v1/admin/posts/{id}`.
- [ ] Regenerate Google Wire nếu interface/constructor thay đổi.
- [ ] Viết unit test service, repository và integration test cho các case trên.
- [ ] Chạy migration kiểm tra schema hiện tại; không tạo migration mới nếu `posts` đã đủ cột/index.

