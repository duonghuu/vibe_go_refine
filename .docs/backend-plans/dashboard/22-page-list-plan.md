# KẾ HOẠCH BACKEND: QUẢN TRỊ DANH SÁCH TRANG TĨNH (PAGE LIST)

**Dự án:** TechBite  
**Tài liệu nguồn:** `.docs/ideas/dashboard/21-page-list-idea.md`  
**Phạm vi:** Tạo nền tảng `Page` độc lập và API Admin List/Delete phục vụ Resource `pages` của Refine.

## 1. Hiện trạng và nguyên tắc thiết kế

- `posts` phụ thuộc `post_types`, có category và luồng nghiệp vụ bài viết; không dùng cho Page.
- Page là aggregate riêng trong `domain/page`, không có `type_code`, category, hierarchy, menu, Media, SEO hay public API trong hạng mục này.
- Tái sử dụng Gin, GORM, Google Wire, `AuthMiddleware`, `RoleMiddleware` và Redis client hiện có. Handler chỉ bind/map response; validation, phân trang, whitelist sort và soft delete nằm trong Service.
- Cấu trúc cần bổ sung:

```text
internal/domain/page/entity/page.go
internal/infrastructure/repository/page_repo.go
internal/usecases/page/dto/page_dto.go
internal/usecases/page/service/page_service.go
internal/controller/page_controller.go
database/migrations/<timestamp>_create_pages_table.{up,down}.sql
```

- Đăng ký `PageSet`, `InitializePageController` trong `internal/di/wire.go`, regenerate `wire_gen.go`, và khai báo routes trong `main.go`.

---

## 2. Trụ cột 1 — Thiết kế dữ liệu

### 2.1. GORM Entity

Vị trí: `apps/backend/internal/domain/page/entity/page.go`.

```go
package entity

import (
    "time"

    "gorm.io/gorm"
)

type PageStatus string

const (
    PageStatusDraft     PageStatus = "DRAFT"
    PageStatusPublished PageStatus = "PUBLISHED"
)

type Page struct {
    ID       uint       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
    Title    string     `gorm:"column:title;type:varchar(255);not null" json:"title"`
    Slug     string     `gorm:"column:slug;type:varchar(255);not null;uniqueIndex:uq_pages_slug" json:"slug"`
    Content  string     `gorm:"column:content;type:longtext;not null" json:"content"`
    Status   PageStatus `gorm:"column:status;type:varchar(20);not null;default:'DRAFT';index:idx_pages_status_deleted" json:"status"`
    AuthorID uint       `gorm:"column:author_id;not null;index:idx_pages_author_id" json:"authorId"`

    CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
    UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
    DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index:idx_pages_status_deleted" json:"-"`
}

func (Page) TableName() string { return "pages" }
```

`AuthorID` là audit field; giá trị này chỉ được nhận từ `userId` do JWT middleware đặt vào Gin Context. Hạng mục List/Delete không tạo DTO ghi Page, nhưng Entity đã bao gồm các trường bắt buộc để các luồng Create/Edit kế tiếp dùng cùng bảng.

### 2.2. Migration MySQL

Tạo migration mới, không sửa migration đã chạy:

```sql
CREATE TABLE pages (
  id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  slug VARCHAR(255) NOT NULL,
  content LONGTEXT NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
  author_id INT NOT NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  deleted_at DATETIME(3) NULL DEFAULT NULL,
  CONSTRAINT chk_pages_status CHECK (status IN ('DRAFT', 'PUBLISHED')),
  CONSTRAINT fk_pages_author
    FOREIGN KEY (author_id) REFERENCES users(id)
    ON DELETE RESTRICT ON UPDATE CASCADE,
  UNIQUE KEY uq_pages_slug (slug),
  KEY idx_pages_status_deleted (status, deleted_at),
  KEY idx_pages_updated_at (updated_at),
  KEY idx_pages_author_id (author_id),
  KEY idx_pages_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

Migration down xóa bảng `pages` sau khi đã xóa foreign key/index theo cơ chế `DROP TABLE pages`; không ảnh hưởng `posts`, `post_types`, `seo_meta` hoặc `media`.

### 2.3. Truy vấn Repository

`PageRepository` cần các hàm:

```go
FindAndCount(ctx context.Context, filter PageListFilter) ([]entity.Page, int64, error)
FindByID(ctx context.Context, id uint) (*entity.Page, error)
Delete(ctx context.Context, id uint) error
```

- `FindAndCount` dùng một base query `Model(&entity.Page{})`; GORM tự áp điều kiện `deleted_at IS NULL`.
- Nếu có `title_like`, lọc có escape wildcard trên `(title LIKE ? OR slug LIKE ?)`.
- Nếu có `status`, lọc đúng enum; không tự suy diễn giá trị khác.
- Sort chỉ nhận whitelist: `id`, `title`, `slug`, `updated_at`. Mọi giá trị khác dùng mặc định `updated_at DESC`; không đưa chuỗi sort từ query thẳng vào `Order`.
- Áp dụng `LIMIT/OFFSET` sau khi đếm tổng; default `current=1`, `pageSize=10`, giới hạn tối đa `100`.
- `Delete` dùng `db.Delete(&entity.Page{}, id)`, do đó soft delete. Service phải kiểm tra Page tồn tại trước để không trả success giả cho ID không tồn tại/đã xóa.

---

## 3. Trụ cột 2 — API Contract và phân quyền

### 3.1. Route và Middleware

Đăng ký dưới `admin := api.Group("/admin")`, vốn đã có `AuthMiddleware(redisClient)`:

```go
pages := admin.Group("/pages")
pages.Use(middleware.RoleMiddleware("ADMIN", "STAFF"))
{
    pages.GET("", pageController.GetPages)
    pages.DELETE("/:page_id", pageController.DeletePage)
}
```

- `AuthMiddleware` xác thực chữ ký JWT trước, kiểm tra blacklist JTI ở Redis, rồi đặt `userId`/`role` vào Context.
- `RoleMiddleware` chặn `CUSTOMER` với `403`; List và Delete đều chỉ dành cho `ADMIN`, `STAFF`.
- Không lấy `userId`, `role`, status hoặc sort authority từ request body/client ngoài contract query đã validate.

### 3.2. List Page

**Method & Route:** `GET /api/v1/admin/pages`

**Query binding:**

```go
type GetPagesQuery struct {
    Current  int    `form:"current" binding:"omitempty,min=1"`
    PageSize int    `form:"pageSize" binding:"omitempty,min=1,max=100"`
    Search   string `form:"title_like" binding:"omitempty,max=255"`
    Status   string `form:"status" binding:"omitempty,oneof=DRAFT PUBLISHED"`
    SortBy   string `form:"sortBy" binding:"omitempty,oneof=id title slug updated_at"`
    Order    string `form:"order" binding:"omitempty,oneof=asc desc"`
}
```

**Response DTO:**

```go
type PageListItemResponse struct {
    ID        uint       `json:"id"`
    Title     string     `json:"title"`
    Slug      string     `json:"slug"`
    Status    PageStatus `json:"status"`
    AuthorID  uint       `json:"authorId"`
    CreatedAt time.Time  `json:"createdAt"`
    UpdatedAt time.Time  `json:"updatedAt"`
}

type PaginatedPageResponse struct {
    Data  []PageListItemResponse `json:"data"`
    Total int64                  `json:"total"`
}
```

- Không trả `content` trong list để giảm payload.
- Success `200`: `{ "data": [...], "total": 42 }`; danh sách rỗng luôn là `[]`, không phải `null`.
- `400`: URI/query không hợp lệ; `401`: thiếu/invalid/revoked JWT; `403`: không đúng role; `500`: lỗi truy vấn không lộ chi tiết DB.

### 3.3. Delete Page

**Method & Route:** `DELETE /api/v1/admin/pages/:page_id`

**URI binding:**

```go
type DeletePageURI struct {
    PageID uint `uri:"page_id" binding:"required,min=1"`
}

type DeletePageResponse struct {
    Message string `json:"message"`
}
```

**Business rules:**

1. Validate URI, tìm Page với `deleted_at IS NULL`.
2. Không tồn tại hoặc đã soft-delete trả `404 page_not_found`.
3. Gọi repository soft delete; không hard delete nội dung.
4. Chỉ invalid cache sau khi delete thành công.
5. Success `200`: `{ "message": "Xóa trang thành công" }`.

**Chuẩn lỗi:**

```json
{
  "error": "page_not_found",
  "message": "Không tìm thấy trang"
}
```

Controller map domain errors: `ErrPageNotFound -> 404`, validation/binding -> `400`, mọi lỗi không xác định -> `500`. Không xử lý GORM, phân quyền hay cache trong Controller.

### 3.4. Service và DI

```go
type PageService interface {
    GetPages(ctx context.Context, query dto.GetPagesQuery) (*dto.PaginatedPageResponse, error)
    DeletePage(ctx context.Context, id uint) error
}
```

- `GetPages` chuẩn hóa pagination, xây filter typed, gọi repository, map entity sang DTO list không content.
- `DeletePage` kiểm tra tồn tại trước rồi soft delete; không triển khai xóa Media/SEO vì các quan hệ đó chưa thuộc phạm vi List Page.
- `PageSet` gồm `NewPageRepository`, `NewPageService`, `NewPageController`; `InitializePageController(db, redisClient)` nhận Redis để chuẩn bị invalidation. Regenerate `wire_gen.go` thay vì sửa tay.

---

## 4. Trụ cột 3 — Redis cache và trạng thái token

### 4.1. Cache Page List

- Không cache `GET /admin/pages` ở v1: danh sách có pagination, search, filter/sort và chỉ dành cho CMS; truy vấn có index và phân trang là đủ.
- Service vẫn có một `invalidateCache` có chủ đích cho pattern tương lai `admin:pages:list:*`; Delete gọi hàm này an toàn khi Redis nil hoặc chưa tồn tại key.
- Không tạo key public, Media hay SEO trong hạng mục này. Khi các module đó được xây, public Page cache sẽ được quy hoạch riêng để tránh cache không đầy đủ.

### 4.2. JWT/Redis hiện có

- Không tạo cơ chế token riêng cho Page.
- Auth middleware sử dụng key blacklist hiện tại `auth:bl:at:{jti}` và chỉ kiểm tra Redis sau khi JWT signature hợp lệ.
- Logout tiếp tục đặt blacklist bằng `SetEx` với TTL còn lại của access token; refresh token rotation tiếp tục revoke refresh JTI cũ và lưu JTI mới theo luồng hiện tại.
- API Page kế thừa hoàn toàn các quy tắc trên; không được ghi token hoặc nội dung Page vào key auth.

---

## 5. Kiểm thử và cập nhật tài liệu

### Test Service/Repository

- List mặc định: pagination, `total`, mảng rỗng và sort `updated_at DESC`.
- Search khớp title hoặc slug; kết hợp search với status.
- `DRAFT`/`PUBLISHED` filter; query status/sort/pageSize không hợp lệ bị controller trả `400`.
- Soft-deleted Page không xuất hiện trong list; delete lần hai và delete ID không tồn tại trả `404`.
- Sort whitelist không thể tạo SQL injection.

### Test API/Middleware

- `GET`/`DELETE` không token hoặc token revoked trả `401`.
- Role `CUSTOMER` trả `403`; `ADMIN` và `STAFF` thực hiện được.
- List response đúng `{data,total}` và không chứa `content`.
- Delete success chỉ soft-delete record và list sau đó không còn record.

### Tài liệu

- Bổ sung `GET /api/v1/admin/pages` và `DELETE /api/v1/admin/pages/{page_id}` vào `.docs/api-endpoints.yaml`.
- Chưa kích hoạt `EntityPage` của SEO registry: nó tiếp tục trả `ErrEntityTypeNotConfigured` đến khi có backend plan/implementation SEO Page riêng.
