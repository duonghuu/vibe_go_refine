# FRONTEND PLAN: Quản lý hình ảnh Trang (Page Media)

## 1. Mục tiêu

- Triển khai Media Card trong `/pages/edit/:id` theo `.docs/ideas/dashboard/24-page-media-management-idea.md`.
- Cho phép Page đã tồn tại có tối đa một ảnh đại diện (`thumbnail`) và nhiều ảnh thư viện (`gallery`) có thứ tự.
- Tái sử dụng luồng upload/preview/reorder của Post Media qua adapter dùng chung, nhưng gọi Page Media API và dùng bảng `page_media` ở Backend.
- Đồng bộ liên kết Media an toàn, phân biệt xóa file tạm với gỡ liên kết của Media đã gắn, đồng thời phản ánh trạng thái loading/error/retry riêng cho Card.

## 2. Phạm vi

### Trong phạm vi

- Tải, upload ảnh, preview, thay thế/gỡ thumbnail, thêm/gỡ nhiều ảnh gallery và sắp xếp lại gallery.
- Đồng bộ snapshot của từng collection khi người dùng lưu Page; retry riêng bước Media nếu Page core đã lưu thành công.
- Kiểm tra sớm định dạng/kích thước ảnh, loading/error/empty/success state, accessibility cơ bản và cảnh báo rời trang khi Media có thay đổi chưa lưu.
- Refactor tối thiểu Post Media sang các module dùng chung để giữ UI/logic Media thống nhất cho Post và Page.
- Invalidate data Page/Media sau khi sync để public Page cache được Backend làm mới và dashboard không hiển thị dữ liệu cũ.

### Ngoài phạm vi

- Media Library độc lập, tìm kiếm/chọn một Media đã có, bulk upload, crop/resize phía client, caption/alt text hay collection ngoài `thumbnail`/`gallery`.
- Hiển thị thumbnail/gallery trên Page List hoặc Page public Webview; các response đó cần được quy hoạch riêng để tránh N+1 request.
- Thay đổi API Post Media, thay đổi Page core/SEO form hoặc tự xóa file Media đã được gắn với Page.

> Lưu ý: Idea cho phép Admin gắn mọi Media hợp lệ và Staff chỉ gắn Media của mình. UI hiện chỉ có `POST /media/upload`; chưa có endpoint list/search Media để xây bộ chọn Media có sẵn. Vì vậy release này hỗ trợ **upload ảnh mới**. Khi có API Media Library, picker sẽ gọi endpoint đó và dùng cùng `attachExisting`/sync contract, không thay đổi UI Card hay Page form.

## 3. Điểm tích hợp và cấu trúc file

Page Media là sub-resource, chỉ được khởi tạo khi `pageId` hợp lệ. `PageCreate` tiếp tục hiển thị card disabled với hướng dẫn lưu Page trước; không upload hoặc tạo liên kết khi chưa có Page.

```text
apps/frontend/src/pages/content-shared/
├── entity-media-types.ts
├── entity-media-api.ts
├── use-entity-media.ts
└── entity-media-card.tsx

apps/frontend/src/pages/posts/
├── post-media-types.ts          # re-export/type adapter tương thích nếu cần
├── post-media-api.ts            # wrapper entityPath="posts"
└── use-post-media.ts            # wrapper hook giữ public API hiện tại

apps/frontend/src/pages/pages/
├── edit.tsx                     # điều phối Page → Media → SEO
└── page-media-card.tsx          # wrapper copywriting “hình ảnh trang”, nếu shared Card không nhận label
```

- `EntityMediaCard` và các preview/uploader bên trong là dumb components: chỉ nhận data, disabled state và callbacks typed; không gọi API.
- `entity-media-api.ts` là nơi duy nhất ghép URL, gửi `FormData`, chuẩn hóa camelCase/snake_case và parse response không tin cậy.
- `useEntityMedia` là container state cho upload, load, remove, reorder, sync và retry. `PageEdit` chỉ truyền `entityPath: "pages"`, `entityId: pageId` và labels Page.
- Không thêm dependency drag-and-drop mới. Tiếp tục dùng HTML5 drag events đang có, cộng với nút di chuyển lên/xuống cho thao tác bàn phím/màn hình nhỏ.
- Không sửa public props/hành vi của `PostMediaCard` trong cùng thay đổi; Post adapters cần được kiểm thử regression sau refactor.

## 4. TypeScript contracts

```ts
export type EntityMediaCollection = "thumbnail" | "gallery";
export type MediaStatus = "temporary" | "attached";
export type EntityMediaPath = "posts" | "pages";

export interface IUploadedMedia {
  id: number;
  fileName: string;
  originalUrl: string;
  thumbnailUrl?: string;
  mediumUrl?: string;
  status: MediaStatus;
}

export interface IEntityMediaItem {
  id: number;
  entityId: number;
  mediaId: number;
  collection: EntityMediaCollection;
  sortOrder: number;
  media?: IUploadedMedia;
}

export interface IEntityMediaListResponse {
  data: IEntityMediaItem[];
  total: number;
}

export interface IMediaSyncItem {
  id: number;
  sortOrder: number;
}

export interface ISyncEntityMediaRequest {
  media: IMediaSyncItem[];
}

export interface EntityMediaState {
  thumbnail: IUploadedMedia | null;
  gallery: IUploadedMedia[];
  isLoading: boolean;
  isUploading: boolean;
  uploadingCount: number;
  isSyncing: boolean;
  error: string | null;
  isDirty: boolean;
}
```

- Raw API DTO có thể chấp nhận cả `file_name`/`fileName`, `sort_order`/`sortOrder`, `page_id`/`pageId` tại đúng lớp normalizer; phần UI chỉ dùng contracts camelCase ở trên.
- `IEntityMediaItem.entityId` biểu diễn Page ID trong adapter Page; không tạo `postId` trong module dùng chung.
- Không dùng `any`, `File` không xuất hiện trong request sync, và `status` trả về từ Media API phải được khai báo union rõ ràng.

## 5. API contract và quy ước gọi API

### 5.1. Endpoint phụ thuộc

```text
POST   /api/v1/media/upload
DELETE /api/v1/media/:id
GET    /api/v1/admin/pages/:page_id/media?collection=thumbnail|gallery
POST   /api/v1/admin/pages/:page_id/media
PUT    /api/v1/admin/pages/:page_id/media/:collection
DELETE /api/v1/admin/pages/:page_id/media/:media_id?collection=thumbnail|gallery
```

- Tất cả request custom dùng `customRequest`/`authenticatedFetch` hiện hữu để kế thừa access-token refresh và cách xử lý `401`; không tạo fetch client riêng.
- `GET` không truyền `collection` để tải cả thumbnail và gallery bằng một request. Nếu Backend chỉ hỗ trợ filter bắt buộc, hook gọi hai request song song và hợp nhất kết quả ở API module.
- `POST /media/upload` gửi `multipart/form-data` với field `file`, lấy `IUploadedMedia` có `id` trước khi thêm vào local state.
- `PUT` là luồng UI chuẩn khi Save: body luôn là snapshot đầy đủ của collection, kể cả `[]` để gỡ thumbnail hoặc xóa toàn bộ gallery. API module map camelCase state sang Backend contract đã thống nhất:

```json
{
  "media": [
    { "id": 102, "sort_order": 0 },
    { "id": 103, "sort_order": 1 }
  ]
}
```

- `POST /pages/:page_id/media` và `DELETE /pages/:page_id/media/:media_id` là endpoint đầy đủ cho thao tác liên kết đơn lẻ/API consumer khác. Media Card không gọi chúng trong normal Save flow, vì `PUT` snapshot tránh race giữa replace, remove và reorder. Khi Media Library được bổ sung, adapter có thể dùng `POST` để attach ngay; state vẫn cần được sync/refresh theo collection.
- `DELETE /media/:id` chỉ được gọi khi người dùng gỡ một `temporary` Media vừa upload mà chưa từng được attach. Với Media `attached`, thao tác gỡ chỉ bỏ item khỏi local collection; `PUT` mới gỡ record `page_media`, tuyệt đối không gọi xóa file gốc.

### 5.2. Phản hồi và lỗi cần xử lý

- `200/201`: normalise response, update state và invalidate query đã stale khi cần.
- `400`: hiển thị lỗi validation theo file/collection; không xoá ảnh đã upload thành công khác.
- `403`: thông báo không đủ quyền gắn Media. UI không tự suy diễn ownership từ client.
- `404`: Page hoặc Media không tồn tại; Page `404` do `PageEdit` xử lý toàn màn hình, Media `404` hiện Alert trong Card và có Retry.
- `409`: thumbnail bị trùng/race hoặc Media đã gắn không hợp lệ; reload collection trước khi cho người dùng thử lại.
- `401`: để AuthProvider xử lý refresh/redirect tập trung.

## 6. UI và tương tác

### 6.1. Vị trí trong Page Edit

- Card **Hình ảnh trang** nằm ở cột phụ của `/pages/edit/:id`, dưới card Status và trước SEO Card như `.docs/frontend-plans/dashboard/25-page-edit-plan.md`.
- Dùng `Card`, `Stack`, `Alert`, `Skeleton`, `LinearProgress`, `IconButton` và palette/spacing/radius của MUI Theme; không hard-code HEX hay tạo style hệ riêng.
- Trong `/pages/create`, card hiển thị disabled notice “Lưu trang trước để quản lý hình ảnh”, không render input file đang hoạt động.

### 6.2. Thumbnail

- Empty state hiển thị vùng chọn/kéo thả một ảnh cùng accepted formats và size limit.
- Khi đã có ảnh, hiển thị preview dùng `thumbnailUrl ?? mediumUrl ?? originalUrl`, file name và actions **Thay thế** / **Gỡ**.
- Upload thumbnail mới thay local preview ngay khi API upload thành công; không gỡ thumbnail cũ trên Backend cho tới khi `PUT .../thumbnail` thành công.
- Thumbnail tối đa một item; sync gửi item tại `sort_order: 0`, hoặc mảng rỗng khi người dùng gỡ.

### 6.3. Gallery

- Nhận nhiều file qua file picker hoặc drag/drop, render grid responsive gồm preview, tên rút gọn, nút gỡ và handle kéo thả.
- Item đang upload có progress/loading state và không thể reorder/gỡ trước khi có Media ID. Nếu nhiều file, hiển thị `uploadingCount`; lỗi một file không làm chặn file còn lại.
- Sau mỗi add/remove/reorder, thứ tự local được đánh lại liên tục từ `0`. `sortOrder` không được nhập tay.
- Drag-and-drop cập nhật local state, không gửi request trên từng thao tác kéo. Cung cấp controls “Di chuyển lên/xuống” có `aria-label` để người dùng keyboard vẫn sắp xếp được.
- Empty gallery là trạng thái hợp lệ: “Chưa có ảnh trong thư viện”.

### 6.4. Validation file

- Chỉ nhận `image/jpeg`, `image/png`, `image/webp`; giới hạn 5 MB/file theo quy ước Media hiện tại. Backend vẫn xác minh MIME/Magic Bytes cuối cùng.
- Không chấp nhận item duplicate cùng `media.id` trong gallery; nếu API picker được thêm, duplicate phải bị chặn trước khi local state đổi.
- Không có ô search/autocomplete trong release này. Nếu thêm Media picker sau này, truy vấn tìm kiếm bắt buộc debounce 300 ms bằng utility được dự án phê duyệt.

## 7. Data flow và trạng thái lưu

### 7.1. Tải dữ liệu

1. `PageEdit` validate route `id` trước; chỉ sau khi Page detail thành công mới gọi `useEntityMedia({ entityPath: "pages", entityId })`.
2. Hook tải Media riêng với Page form. Trong lúc tải, chỉ Card render Skeleton, không làm Page form rỗng hay khóa việc sửa Page core.
3. Hook tách `thumbnail` và `gallery`, sort gallery tăng dần theo `sortOrder`, sau đó reset `isDirty: false`.
4. Media load lỗi giữ nguyên Page/SEO form và hiển thị Alert + **Thử lại** gọi lại `load`.

### 7.2. Lưu trong Page Edit

`PageEdit` giữ flow đã quy hoạch: **Page core → Page Media → SEO**.

1. Validate Page/SEO form; nếu lỗi, không gọi Media.
2. Nếu form Page dirty, gọi `PUT /admin/pages/:id`. Nếu lỗi, dừng; local Media vẫn còn nguyên.
3. Nếu `media.isDirty`, lần lượt gọi `PUT .../thumbnail`, rồi `PUT .../gallery`. Disable action save trong lúc `isUploading || isSyncing` để không gửi snapshot thiếu Media ID hoặc double-submit.
4. Nếu cần sync hai collection, một collection thất bại sẽ dừng bước SEO, giữ state Media dirty và lưu lại collection lỗi. Nút **Thử đồng bộ Media** chỉ retry Media, không gửi lại Page core.
5. Khi cả ba bước thành công, invalidate `pages`, Page detail, Page Media và SEO cache/query; reset dirty state, toast thành công, điều hướng `/pages` theo hành vi Page Edit.
6. Nếu Page đã lưu nhưng Media lỗi, không redirect. Hiển thị toast/Alert “Trang đã lưu, hình ảnh chưa đồng bộ” và để người dùng retry; không làm mất gallery/thumbnail local.

### 7.3. Unsaved changes và cancel

- `isPageDirty = formState.isDirty || media.isDirty || seo.isDirty`.
- `warnWhenUnsavedChanges`/confirmation hiện có phải được kích hoạt khi Media local thay đổi, kể cả Page form không dirty.
- Hủy/rời route trong lúc upload hoặc sync phải khóa thao tác hoặc yêu cầu xác nhận theo luồng chung; không chủ động gọi DELETE cho các Media pending chỉ vì user mở dialog Cancel.

## 8. Cache và tính nhất quán

- Không tự gọi API public để làm mới cache. Mỗi mutation attach/sync/delete liên kết phải để Backend invalidate public Page cache theo rule nghiệp vụ.
- Sau sync thành công, Frontend invalidate Refine resource `pages` và query key Page Media; load lại collection để lấy thứ tự/source of truth do Backend trả về.
- Không làm N+1: Page List/Public Page cần thumbnail phải nhận `thumbnailUrl` trong response aggregate của endpoint tương ứng, không gọi Page Media API cho từng Page.
- Trường hợp sync response không trả collection mới, hook phải `GET` lại một lần sau success; không đoán sort order từ response cũ.

## 9. Acceptance criteria

- `/pages/create` không cho upload/gắn Media khi chưa có Page ID; sau khi tạo redirect tới `/pages/edit/:id` thì Card hoạt động.
- Mở Page Edit tải đúng thumbnail/gallery, gallery luôn đúng thứ tự `sortOrder`, và lỗi Media không chặn Page form/SEO Card.
- Chỉ có một thumbnail local/saved; thay thế rồi Save gỡ liên kết thumbnail cũ và giữ Media file cũ.
- Gallery hỗ trợ upload nhiều ảnh, remove/reorder bằng pointer và keyboard, sau Save Backend nhận snapshot `sort_order` liên tục từ `0`.
- Gỡ Media attached chỉ gọi sync/gỡ liên kết, không gọi `DELETE /media/:id`; gỡ Media temporary gọi endpoint cleanup phù hợp.
- Upload/sync lỗi giữ tất cả local state liên quan, thông báo rõ collection/file lỗi, cho retry mà không gửi lại Page core đã thành công.
- Thao tác Save bị khóa trong upload/sync; cảnh báo rời trang xuất hiện khi Media dirty.
- Admin/Staff nhận đúng `403` từ Backend theo ownership/RBAC; client không gửi hay quyết định role/owner.
- API calls sử dụng auth client hiện có; module Page/Post sau refactor type-safe và không dùng `any` mới.

## 10. Checklist triển khai

- [ ] Backend tạo `page_media` và đầy đủ Page Media endpoints, xác thực ảnh/ownership, unique thumbnail, sort order và invalidation public Page cache.
- [ ] Thống nhất response GET, wrapper `{ data, total }` và body sync dùng `sort_order` tại API boundary.
- [ ] Tạo `entity-media-types.ts`, `entity-media-api.ts`, `use-entity-media.ts`, UI Card/preview/uploader shared.
- [ ] Chuyển Post Media qua adapter/wrapper và chạy regression Create/Edit Post trước khi dùng cho Page.
- [ ] Tích hợp Page Media Card vào `PageEdit`, disabled notice vào `PageCreate`, và nối aggregate dirty/save/retry flow.
- [ ] Test upload file hợp lệ/không hợp lệ/quá dung lượng, partial upload failure, replace/remove thumbnail, empty gallery, reorder, refresh Edit, retry và double submit.
- [ ] Test authorization Admin/Staff, Page/Media 404, race/409 và xác nhận Media attached không bị xóa vật lý.
- [ ] Kiểm tra Network: một lần mở Page Edit chỉ tải Page, Media và SEO theo từng subresource; không có request Media theo số lượng preview/Page List rows.
- [ ] Chỉ bổ sung picker Media Library khi API list/search và permission contract được phê duyệt; dùng debounce 300 ms cho search.
