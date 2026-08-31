# FRONTEND PLAN: Chỉnh sửa trang tĩnh (Page Edit)

## 1. Mục tiêu

- Xây dựng màn hình `/pages/edit/:id` từ `.docs/ideas/dashboard/23-page-edit-idea.md`.
- Cho phép Admin/Staff tải và cập nhật title, slug, rich content, `DRAFT`/`PUBLISHED`, thumbnail/gallery và SEO của một Page đã tồn tại.
- Giữ Page là resource độc lập: không hiển thị hoặc gửi `post_type`, category, hierarchy, menu hay `authorId`.
- Điều phối mutation theo thứ tự Page → Media → SEO, với error/retry riêng để không mất dữ liệu đã lưu một phần.

## 2. Phạm vi

### Trong phạm vi

- Fetch detail, Edit form, URL preview, update Page và trạng thái loading/error/not-found.
- Tích hợp Page Media Card (thumbnail/gallery) và SEO Meta Card cho `entityType: "page"`.
- Cảnh báo thay đổi chưa lưu, validation, map field error slug, success/partial-success feedback và retry Media/SEO.
- Refactor tối thiểu các component/hook Media/SEO của Post thành phần dùng chung, giữ nguyên hành vi Post đang có.

### Ngoài phạm vi

- Redirect từ slug cũ, Page public route/rendering, menu, hierarchy, scheduled publishing, Page Builder và bulk update.
- Thay đổi nội dung/bố cục của các form Post ngoài việc chuyển chúng sang shared adapter tương thích.
- Tự động xuất bản hoặc gọi public API ngay khi status đổi.

## 3. Cấu trúc và routing

### Files Page

```text
apps/frontend/src/pages/pages/
├── index.ts
├── edit.tsx
├── page-form-types.ts
├── page-basic-info-card.tsx
├── page-public-url-preview.tsx
└── page-edit-actions.tsx
```

`PageEdit` là smart component; các card/form/action/preview là dumb component chỉ nhận props và callback typed.

### Shared Media và SEO

Tách phần đang gắn cứng với Post thành adapter dùng chung, sau đó giữ wrapper Post tương thích:

```text
apps/frontend/src/pages/content-shared/
├── entity-media-api.ts
├── entity-media-types.ts
├── use-entity-media.ts
├── entity-media-card.tsx
├── entity-seo-api.ts
├── entity-seo-types.ts
├── use-entity-seo.ts
└── entity-seo-card.tsx
```

- `posts/use-post-media.ts` và `posts/use-post-seo.ts` trở thành adapter truyền `entityPath: "posts"`, `entityType: "post"`; không thay đổi public props/hành vi của Post UI.
- `PageEdit` dùng adapter `entityPath: "pages"`, `entityType: "page"` và copywriting “trang”.
- Shared UI không gọi API; API module/hook là nơi duy nhất gọi `customRequest`/Refine mutation.

### Route Refine

Hoàn thiện resource và route đã được khai báo trước đó:

```tsx
<Route path="/pages">
  <Route index element={<PageList />} />
  <Route path="create" element={<PageCreate />} />
  <Route path="edit/:id" element={<PageEdit />} />
</Route>
```

## 4. UI và trải nghiệm

### 4.1. Header và action bar

- Tiêu đề: **Chỉnh sửa trang**.
- Breadcrumb: `Trang chủ > Quản trị nội dung > Trang > Chỉnh sửa`.
- Nút **Hủy** quay `/pages`; nếu Page/Media/SEO đang dirty, dùng unsaved-change warning/confirmation.
- Nút **Lưu thay đổi** giữ `status` hiện tại trong form; khi đang load hoặc saving thì disabled.
- Không tự thay đổi status khi sửa một field khác. Status dùng MUI control với đúng hai giá trị `DRAFT`, `PUBLISHED`.

### 4.2. Layout hai cột

- Cột chính:
  - `PageBasicInfoCard` tái sử dụng từ Create: title, slug và Rich Text Editor.
  - `PagePublicUrlPreview` hiển thị đường dẫn tương đối `/${slug}`; đây chỉ là preview, không gọi public API và không tạo redirect slug cũ.
- Cột phụ:
  - Status card.
  - `EntityMediaCard` cho thumbnail và gallery.
  - `EntitySeoCard` với `entityType="page"`.
- Desktop dùng grid 2/3–1/3; màn hình nhỏ chuyển về một cột. Dùng MUI Theme tokens, không hard-code màu.

### 4.3. Media card

- Thumbnail: tối đa một ảnh, preview/thay thế/gỡ.
- Gallery: upload nhiều ảnh, preview, gỡ từng ảnh, drag-and-drop reorder và `sortOrder` theo vị trí hiện tại.
- Lúc mở Edit, Media tải riêng với skeleton. Lỗi Media không chặn Page form; hiển thị Alert và Retry.
- Xóa file temporary gọi Media API ngay; gỡ ảnh attached chỉ đổi local state và sync khi Save.

### 4.4. SEO card

- Reuse các trường hiện có: meta title/description/keywords, canonical, Open Graph, Twitter, robots và Schema JSON-LD.
- GET trả giá trị saved + resolved fallback; preview phải phân biệt input tùy chỉnh với fallback và không gửi `resolved_*` về server.
- SEO load/error/retry độc lập; nếu chưa có record, khởi tạo form rỗng với resolved preview từ Page.

## 5. TypeScript contracts và state

```ts
export type PageStatus = "DRAFT" | "PUBLISHED";

export interface IPageResponse {
  id: number;
  title: string;
  slug: string;
  content: string;
  status: PageStatus;
  authorId: number;
  createdAt: string;
  updatedAt: string;
}

export interface IPageUpdatePayload {
  title: string;
  slug: string;
  content: string;
  status: PageStatus;
}

export type EntityMediaCollection = "thumbnail" | "gallery";

export interface IEntityMediaItem {
  id: number;
  entityId: number;
  mediaId: number;
  collection: EntityMediaCollection;
  sortOrder: number;
  media?: IUploadedMedia;
}
```

- Page form không có `id`, `authorId` hay `entityType` để submit.
- `useEntityMedia` quản lý `thumbnail`, `gallery`, `isLoading`, `isUploading`, `isSyncing`, `error`, `isDirty`.
- `useEntitySeo` quản lý SEO form, `isLoading`, `isSaving`, `error`, `isDirty`, resolved preview và retry.
- `isPageDirty` tổng hợp `formState.isDirty || media.isDirty || seo.isDirty` để dùng cho warning/cancel. Không dùng `any` trong contracts, adapters hay components mới.

## 6. Refine và data flow

### 6.1. Tải dữ liệu

1. Validate `id` từ route trước khi gọi API; ID không hợp lệ hiển thị error state và link về `/pages`.
2. Dùng `useForm<IPageResponse, HttpError, IPageUpdatePayload>` với `action: "edit"`, `resource: "pages"`, `redirect: false`, `warnWhenUnsavedChanges: true` để Refine tải `GET /admin/pages/:id`.
3. Khi Page detail thành công, khởi tạo Media và SEO query song song với `enabled: Boolean(pageId)`.
4. Chỉ Page `404` là lỗi chặn form. Media/SEO lỗi chỉ làm card tương ứng lỗi; dữ liệu Page chính vẫn giữ được.

### 6.2. Slug và validation

- Edit không auto-generate slug theo title; slug hiện tại luôn được giữ đến khi người dùng tự sửa.
- Dùng cùng validation với Page Create: title/slug tối đa 255, slug regex `^[a-z0-9]+(?:-[a-z0-9]+)*$`, content không rỗng sau khi loại HTML trống.
- Không có API kiểm tra slug theo mỗi lượt gõ. Lỗi `409` khi lưu map vào `setError("slug", ...)`.

### 6.3. Thứ tự Save

1. Trigger validation cho Page form và SEO form; nếu không hợp lệ, dừng và focus lỗi.
2. Gọi `PUT /admin/pages/:id` với Page payload khi Page form dirty. Nếu Page không dirty, giữ response Page hiện tại làm context.
3. Nếu Page update thành công (hoặc không có thay đổi), sync `thumbnail` và `gallery` chỉ khi `media.isDirty`.
4. Upsert SEO chỉ khi `seo.isDirty`.
5. Khi mọi bước thành công: invalidate `pages`, Page detail, Page Media và SEO query; reset dirty state; hiển thị toast thành công và quay `/pages`.
6. Nếu Page update lỗi: không gọi Media/SEO và giữ tất cả local state.
7. Nếu Media hoặc SEO lỗi sau Page update: không redirect, hiển thị **Page đã lưu, Media/SEO chưa đồng bộ**, giữ phần lỗi để Retry mà không gửi lại Page update.

## 7. API contract prerequisites

### Page core

```text
GET /api/v1/admin/pages/:id
PUT /api/v1/admin/pages/:id
```

`PUT` body:

```json
{
  "title": "Giới thiệu về TechBite",
  "slug": "gioi-thieu",
  "content": "<p>Nội dung đã cập nhật</p>",
  "status": "PUBLISHED"
}
```

`GET`/`PUT` trả `{ "data": IPageResponse }`; `404` cho Page soft-deleted/không tồn tại, `409` khi slug trùng, `400` validation và `401/403` Auth/RBAC.

### Page Media

```text
POST   /api/v1/media/upload
DELETE /api/v1/media/:id
GET    /api/v1/admin/pages/:page_id/media?collection=thumbnail|gallery
PUT    /api/v1/admin/pages/:page_id/media/:collection
```

Sync request dùng contract chung:

```json
{
  "media": [
    { "id": 101, "sortOrder": 0 },
    { "id": 102, "sortOrder": 1 }
  ]
}
```

### SEO

```text
GET /api/v1/admin/seo-meta/page/:page_id
PUT /api/v1/admin/seo-meta/page/:page_id
```

- Backend phải kích hoạt `EntityPage` resolver; hiện registry của dự án vẫn trả `ErrEntityTypeNotConfigured` cho Page.
- API Media và SEO là prerequisite: không dùng mock/fallback request nếu chúng chưa tồn tại.

## 8. Loading, error và success

- **Page loading:** skeleton toàn form; không render form rỗng trước khi GET Page xong.
- **Page 404/error:** Error state rõ ràng với action quay danh sách; không thử tải Media/SEO tiếp.
- **Media/SEO loading:** skeleton riêng ở mỗi card.
- **Media/SEO error:** Alert + Retry, không reset Page form hoặc local changes ở card còn lại.
- **Update success:** Snackbar và refresh preview URL/status/list.
- **Partial success:** thông báo bước đã lưu và lỗi còn lại; Retry chỉ gọi mutation thất bại.
- **401:** giữ luồng refresh/redirect tập trung ở AuthProvider/data provider, không tự xử lý token tại PageEdit.

## 9. Acceptance test

- Mở `/pages/edit/:id` tải đúng Page; Page không tồn tại/soft-deleted hiển thị 404 state.
- Form giữ slug khi title đổi; chỉ payload user-edit mới thay slug, URL preview phản ánh slug local.
- `PUT` chỉ gửi title/slug/content/status, không gửi authorId, post type, category, Media hoặc SEO fields.
- Media tải/sync đúng thumbnail/gallery, chặn >1 thumbnail, giữ sort order, và không gọi theo từng ảnh ở list.
- SEO Page gửi `entityType=page`, không gửi `resolved_*`, xử lý empty/error/retry.
- Save Page → Media → SEO đúng thứ tự; Page lỗi chặn steps sau; Media/SEO lỗi tạo partial-success và retry không PUT Page lần nữa.
- Cancel/route change cảnh báo khi bất kỳ Page, Media hoặc SEO state dirty.
- Module mới type-safe, không dùng `any`; Post Media/SEO regression vẫn hoạt động với adapter `post`.

## 10. Phụ thuộc và giả định

- Backend Page List/Create đã tạo `pages` resource; cần bổ sung Page Detail/Update, Page Media và Page SEO trước khi tích hợp UI thật.
- Slug chỉ unique trong namespace Page theo plan hiện có; không triển khai redirect hoặc registry slug xuyên Post/Page.
- Public URL preview dùng relative path; base URL/canonical cuối cùng do SEO service resolve.
- `PageEdit` giữ status hiện tại khi lưu; không thêm nút Lưu nháp/Xuất bản riêng ở màn hình này.
