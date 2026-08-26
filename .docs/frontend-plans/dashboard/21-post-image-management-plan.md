# FRONTEND PLAN: Quản lý hình ảnh bài viết (Post Image Management)

## 1. Mục tiêu

- Bổ sung quản lý ảnh đại diện (`thumbnail`) và thư viện ảnh (`gallery`) vào màn hình Tạo mới/Chỉnh sửa bài viết.
- Giữ `media.id` trong Frontend state để tạo quan hệ qua `post_media`; không chỉ lưu URL như luồng Product hiện tại.
- Tái sử dụng trải nghiệm upload, preview và xóa ảnh của Product nhưng chuẩn hóa thành component dùng chung cho Post.
- Hiển thị thumbnail trong Post List bằng dữ liệu trả về từ một API list duy nhất, không gọi API Media riêng cho từng dòng.
- Tuân thủ Refine.js, React Hook Form, MUI và TypeScript strict; không sử dụng `any`.

## 2. Phạm vi chức năng

### 2.1. Trong phạm vi

- Upload một ảnh thumbnail.
- Upload nhiều ảnh gallery.
- Preview, thay thế, xóa ảnh.
- Kéo-thả hoặc thao tác sắp xếp lại gallery và cập nhật `sortOrder`.
- Tải danh sách Media hiện tại khi mở trang Edit.
- Sync thumbnail/gallery sau khi lưu Post.
- Hiển thị thumbnail ở Post List.
- Loading, success, empty và error state cho upload, fetch Media và sync.
- Cảnh báo rời trang khi có ảnh hoặc form chưa lưu.

### 2.2. Ngoài phạm vi giai đoạn đầu

- Collection `content` trong Rich Text Editor.
- Collection `attachment` cho PDF/DOC.
- Xây dựng Media Library độc lập.
- Tải gallery trong Post List hoặc preload toàn bộ Media cho mỗi dòng.

## 3. Luồng màn hình

### 3.1. Create Post

Giữ layout 2 cột hiện tại của `PostCreate`:

- Cột chính: Card thông tin cơ bản gồm title, slug, content.
- Cột phụ: Card phân loại và Card `Hình ảnh bài viết`.
- Footer: Hủy và Lưu bài viết.

Card hình ảnh gồm hai khu vực:

1. **Ảnh đại diện**: một preview duy nhất, nút chọn/thay thế/xóa.
2. **Thư viện ảnh**: vùng chọn nhiều file, lưới preview, trạng thái upload, nút xóa và sắp xếp.

Khi người dùng chọn file, Frontend gọi `POST /api/v1/media/upload`, lưu lại object Media trả về và giữ trạng thái `temporary` cho tới khi Post được lưu.

### 3.2. Edit Post

- Sau khi Post detail load thành công, gọi `GET /api/v1/admin/posts/{postId}/media`.
- Tách response theo `collection` thành `thumbnail` và `gallery`.
- Hiển thị loading riêng trong Card hình ảnh; không làm form Post bị blank khi Media API đang tải.
- Nếu Media API lỗi, giữ form Post hoạt động nhưng hiển thị cảnh báo và cho phép Retry.
- Khi lưu, cập nhật Post và sync từng collection theo contract Backend.

### 3.3. Post List

- Thêm cột `Thumbnail` ở đầu hoặc gần cột Title.
- Dùng `thumbnailUrl` trong response list; không gọi `GET /posts/{id}/media` cho từng row.
- Dùng `thumbnailUrl` với `mediumUrl` nếu Backend cung cấp; fallback về placeholder khi không có ảnh.
- Giữ server-side pagination, sorting và debounce tìm kiếm hiện tại.

## 4. Component breakdown

### 4.1. Component dùng chung

- `PostMediaCard`: Card tổng hợp thumbnail và gallery; chỉ nhận props/state handlers, không tự gọi API.
- `MediaUploader`: Input file/dropzone, kiểm tra định dạng/dung lượng và gọi callback upload.
- `MediaPreview`: Preview một file, tên file, trạng thái upload và nút xóa.
- `ThumbnailMediaField`: Dumb component cho collection `thumbnail`.
- `GalleryMediaField`: Dumb component cho danh sách gallery và reorder.
- `PostMediaGrid`: Lưới preview gallery responsive.
- `PostMediaErrorState`: Hiển thị lỗi Media và nút Retry.

### 4.2. Logic/container

- `usePostMedia`: hook quản lý state thumbnail/gallery, upload, remove, reorder, fetch và sync.
- `post-media-api.ts`: các hàm gọi upload, list Media, sync collection và delete temporary Media.
- `post-media-types.ts`: TypeScript interfaces cho Upload response, Post Media response và request payload.
- `PostCreate` và `PostEdit`: Smart components điều phối `useForm`, `usePostMedia`, submit và redirect.

UI component không được gọi API trực tiếp; mọi request đi qua hook/API module.

## 5. TypeScript contracts

```ts
export type MediaStatus = "temporary" | "attached";

export interface IUploadedMedia {
  id: number;
  fileName: string;
  originalUrl: string;
  thumbnailUrl?: string;
  mediumUrl?: string;
  status: MediaStatus;
}

export interface IPostMediaItem {
  id: number;
  postId: number;
  mediaId: number;
  collection: "thumbnail" | "gallery";
  sortOrder: number;
  media?: IUploadedMedia;
}

export interface IPostMediaResponse {
  data: IPostMediaItem[];
  total: number;
}

export interface IMediaSyncItem {
  id: number;
  sortOrder: number;
}

export interface ISyncPostMediaRequest {
  media: IMediaSyncItem[];
}

export interface IPostMediaFormValue {
  thumbnail: IUploadedMedia | null;
  gallery: IUploadedMedia[];
}
```

Post request type cần mở rộng theo Backend contract:

```ts
export interface IPostMediaPayload {
  thumbnailId: number | null;
  gallery: IMediaSyncItem[];
}
```

Nếu Backend chưa hỗ trợ nested `media` trong `POST/PUT /posts`, Frontend tạm thời tạo/cập nhật Post trước, sau đó gọi hai API sync collection. Khi đó phải giữ Post ID và hiển thị thông báo rõ nếu một bước sync thất bại.

## 6. Data và state management

### 6.1. Upload state

`usePostMedia` quản lý tối thiểu:

```ts
interface PostMediaState {
  thumbnail: IUploadedMedia | null;
  gallery: IUploadedMedia[];
  isUploading: boolean;
  uploadingCount: number;
  isLoading: boolean;
  isSyncing: boolean;
  error: string | null;
  isDirty: boolean;
}
```

- File đang upload không được đưa vào payload sync cho tới khi có `id`.
- `uploadingCount > 0` hoặc `isSyncing` thì khóa nút Lưu.
- Xóa Media temporary chưa liên kết có thể gọi DELETE ngay; xóa Media attached chỉ thay đổi danh sách local và gỡ liên kết khi sync.
- Reorder chỉ cập nhật local state; gửi `sortOrder` khi người dùng submit.

### 6.2. Refine và React Hook Form

- Tiếp tục dùng `useForm` của `@refinedev/react-hook-form` cho Post.
- Không nhét `File` object vào request Post; chỉ gửi Media ID và metadata cần thiết.
- `warnWhenUnsavedChanges: true` phải phản ánh cả `formState.isDirty` và `postMedia.isDirty`.
- Khi `onMutationSuccess` của Post chạy, chỉ redirect sau khi Media sync thành công.
- Sau mutation thành công, invalidate `posts` để Post List lấy thumbnail mới.

### 6.3. API client

API module dùng `API_URL` và cơ chế Auth hiện có. Không tạo fetch wrapper riêng có logic refresh token khác với provider.

Các hàm cần có:

```ts
uploadMedia(file: File): Promise<IUploadedMedia>;
getPostMedia(postId: number, collection?: "thumbnail" | "gallery"): Promise<IPostMediaResponse>;
syncPostMedia(postId: number, collection: "thumbnail" | "gallery", items: IMediaSyncItem[]): Promise<void>;
deleteTemporaryMedia(mediaId: number): Promise<void>;
```

Nếu dùng `fetch`, phải giữ `credentials`/Authorization tương thích với AuthProvider hiện tại và parse lỗi về `HttpError` rõ ràng.

## 7. UI validation và interaction

### 7.1. File validation

- Chỉ nhận `image/jpeg`, `image/png`, `image/webp`, phù hợp với Media service hiện tại.
- Giới hạn dung lượng theo cấu hình Backend; UI nên chặn sớm bằng giới hạn 5MB nếu đó là quy ước hiện tại của Dashboard.
- Không tin phần mở rộng file; Backend vẫn là nơi kiểm tra MIME/Magic Bytes cuối cùng.
- Gallery không cho chọn file trùng với cùng Media ID.

### 7.2. Thumbnail

- Chỉ cho phép một item.
- Chọn ảnh mới thay thế preview local; ảnh cũ vẫn giữ liên kết ở Backend cho tới khi sync thành công.
- Xóa thumbnail gửi danh sách rỗng trong sync để gỡ liên kết.

### 7.3. Gallery

- `sortOrder` được đánh lại liên tục từ `1` sau mỗi lần reorder.
- Hiển thị empty state `Chưa có ảnh trong thư viện` khi danh sách rỗng.
- Có thể dùng HTML5 drag events hoặc thư viện drag-and-drop đã được phê duyệt; không tự thêm dependency nếu chưa cần thiết.
- Hỗ trợ keyboard/focus state cơ bản cho nút xóa và reorder trên màn hình nhỏ.

## 8. API contract prerequisite

Frontend phụ thuộc các endpoint:

```text
POST   /api/v1/media/upload
DELETE /api/v1/media/:id
GET    /api/v1/admin/posts/:post_id/media?collection=thumbnail|gallery
PUT    /api/v1/admin/posts/:post_id/media/:collection
```

Request sync thống nhất dùng camelCase ở Frontend:

```json
{
  "media": [
    { "id": 101, "sortOrder": 0 }
  ]
}
```

Backend hiện có DTO dùng `sort_order`; cần thống nhất contract hoặc để data provider/API module map `sortOrder` sang `sort_order` tại một nơi duy nhất. Không để từng component tự map field.

Post List cần Backend trả thêm `thumbnailUrl` trong response list bằng một JOIN thumbnail. Frontend không được triển khai vòng lặp gọi Media API theo từng row.

## 9. Loading, Error, Empty và Success

- **Upload loading:** preview skeleton/progress, khóa item đang upload và hiển thị số lượng đang xử lý.
- **Media loading:** skeleton trong Card hình ảnh ở Edit.
- **Upload error:** giữ các ảnh upload thành công, hiển thị lỗi cho file thất bại và cho phép thử lại.
- **Media GET error:** thông báo lỗi và nút Retry; không xóa dữ liệu form Post.
- **Post save error:** giữ nguyên form và danh sách ảnh local.
- **Media sync error:** không redirect; hiển thị collection lỗi, cho phép Retry sync hoặc đi tới Edit để khôi phục.
- **Empty thumbnail:** hiển thị placeholder nhẹ, không coi là lỗi nếu thumbnail là tùy chọn.
- **Success:** Snackbar `Lưu bài viết và hình ảnh thành công`, invalidate resource `posts`, rồi redirect đúng `type_code`.

## 10. Tối ưu Post List

- Mở rộng `IPostResponse` với `thumbnailUrl?: string | null`.
- Cột Thumbnail chỉ render URL đã có trong row data.
- Dùng kích thước ảnh `thumbnailUrl`/`mediumUrl`, không tải `originalUrl` kích thước lớn ở DataGrid.
- Dùng placeholder khi ảnh lỗi hoặc không có thumbnail.
- Không preload gallery, attachment hoặc toàn bộ `media` object trong list.
- Kiểm tra Network tab để bảo đảm một lần load list chỉ có request Post List, không phát sinh request Media theo số lượng row.

## 11. Checklist triển khai

- [ ] Tạo `post-media-types.ts` với các interface strict, không dùng `any`.
- [ ] Tạo `post-media-api.ts` dùng Auth/API config hiện có.
- [ ] Tạo `usePostMedia` xử lý upload, fetch, remove, reorder và sync.
- [ ] Tạo `PostMediaCard`, `ThumbnailMediaField`, `GalleryMediaField` và trạng thái UI tương ứng.
- [ ] Tích hợp Card hình ảnh vào `posts/create.tsx`.
- [ ] Tích hợp fetch Media, Card hình ảnh và sync vào `posts/edit.tsx`.
- [ ] Mở rộng request type của Post với payload Media nếu Backend hỗ trợ transaction.
- [ ] Nếu chưa có nested payload, triển khai fallback create/update Post rồi sync hai collection.
- [ ] Bổ sung `thumbnailUrl` vào type và cột của `posts/list.tsx`.
- [ ] Bảo đảm Post List không tạo N+1 request Media.
- [ ] Kiểm thử upload thành công/thất bại, file sai định dạng, file quá lớn, xóa temporary, xóa attached, reorder gallery và retry sync.
- [ ] Kiểm thử refresh trang Edit, unsaved warning, redirect theo `type_code` và double-submit.
- [ ] Cập nhật `.docs/api-endpoints.yaml`/contract nếu Backend đổi `sort_order` sang `sortOrder` hoặc thêm nested `media`.

