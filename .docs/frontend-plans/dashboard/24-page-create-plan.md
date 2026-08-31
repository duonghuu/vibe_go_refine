# FRONTEND PLAN: Tạo trang tĩnh (Page Create)

## 1. Mục tiêu

- Xây dựng màn hình `/pages/create` từ `.docs/ideas/dashboard/22-page-create-idea.md`.
- Tạo Page độc lập qua resource `pages`, không truyền `typeCode`, category hay `authorId`.
- Cho phép lưu cùng một form thành `DRAFT` hoặc `PUBLISHED`, sau đó chuyển tới màn hình Edit để gắn Media và SEO.
- Tuân thủ Refine, React Hook Form, MUI, TypeScript strict và `warnWhenUnsavedChanges` hiện có.

## 2. Phạm vi

### Trong phạm vi

- Route/Refine resource Create, form title/slug/content, action Lưu nháp/Xuất bản/Hủy.
- Auto-slug đến khi người dùng tự chỉnh slug; validation client-side và map lỗi field từ Backend.
- Bố cục desktop hai cột, responsive một cột; card Media và SEO ở trạng thái chờ Page được tạo.
- Loading, submit success/error, unsaved-change warning và redirect có chủ đích.

### Ngoài phạm vi

- Upload/sync thumbnail/gallery, gọi API SEO và màn hình Edit Page; chúng thuộc luồng sau khi Page đã có `id`.
- Public Page, page hierarchy, menu, schedule xuất bản, Post Type hoặc Page Builder.
- Thay đổi hành vi Post Create đang ổn định.

## 3. Cấu trúc và routing

### Files dự kiến

```text
apps/frontend/src/pages/pages/
├── index.ts
├── create.tsx
├── page-form-types.ts
├── page-basic-info-card.tsx
└── page-post-create-notice.tsx

apps/frontend/src/utils/
└── generate-slug.ts
```

- `create.tsx` là smart component, điều phối Refine mutation, React Hook Form, action status và điều hướng.
- `PageBasicInfoCard` và `PagePostCreateNotice` là dumb component; chỉ nhận value/error/callback, không gọi API.
- Trích xuất helper slug thuần đang được lặp lại trong các form hiện hữu thành `generate-slug.ts`; giữ kết quả tương thích slug lowercase/kebab-case của Post/Product/Category.

### Refine registration

Trong `apps/frontend/src/App.tsx`, resource `pages` có `create: "/pages/create"`; bổ sung route:

```tsx
<Route path="/pages">
  <Route index element={<PageList />} />
  <Route path="create" element={<PageCreate />} />
</Route>
```

`PageEdit` chưa là dependency runtime của form, nhưng là đích redirect phải được triển khai cùng hoặc trước khi phát hành Create.

## 4. UI và trải nghiệm

### 4.1. Header và action bar

- Tiêu đề: **Tạo trang mới**.
- Breadcrumb: `Trang chủ > Quản trị nội dung > Trang > Tạo mới`.
- Nút **Hủy**: nếu form clean thì chuyển `/pages`; nếu dirty thì dùng luồng warning của Refine/confirmation hiện có.
- Nút **Lưu nháp** submit cùng form với `status: "DRAFT"`.
- Nút **Xuất bản** submit cùng form với `status: "PUBLISHED"`.
- Khi mutation đang pending, khóa cả hai nút lưu và Hủy để tránh gửi trùng hoặc điều hướng giữa chừng.

### 4.2. Layout form

- Grid desktop: cột chính chiếm khoảng hai phần ba, cột phụ một phần ba; tại breakpoint nhỏ chuyển thành một cột.
- Cột chính dùng `PageBasicInfoCard`:
  - `title`: MUI TextField bắt buộc.
  - `slug`: MUI TextField bắt buộc, có helper text về URL public.
  - `content`: Rich Text Editor đang được dự án dùng cho Post, điều khiển bằng `Controller` từ React Hook Form.
- Cột phụ:
  - Card trạng thái hiển thị action đang chọn và giải thích Draft/Published.
  - `PagePostCreateNotice` cho **Hình ảnh** và **SEO**: disabled, không chứa upload/form SEO, hiển thị hướng dẫn “Lưu trang trước để quản lý hình ảnh và SEO”.
- Dùng MUI Theme tokens, spacing/radius/border từ Design System; không hard-code màu HEX.

## 5. TypeScript contracts và form state

```ts
export type PageStatus = "DRAFT" | "PUBLISHED";

export interface IPageCreateFormValues {
  title: string;
  slug: string;
  content: string;
}

export interface IPageCreatePayload extends IPageCreateFormValues {
  status: PageStatus;
}

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
```

- `authorId` chỉ có trong response, không có trong form hoặc payload.
- Status không cần field input thường trực: `PageCreate` giữ `submitStatus: PageStatus | null` và chỉ ghép vào payload ngay trước `onFinish`.
- `IPageResponse` phải được dùng làm generic của Refine mutation, không dùng `any` hay cast không kiểm soát.

## 6. Refine, validation và auto-slug

### 6.1. Form integration

- Dùng `useForm<IPageCreateFormValues, HttpError, IPageCreatePayload>` từ `@refinedev/react-hook-form` với `resource: "pages"`.
- Dùng `handleSubmit` cho cả hai action; action handler truyền status vào một hàm `submitPage(status)` duy nhất.
- `onMutationSuccess` lấy `data.id` từ response và điều hướng tới `/pages/edit/:id`; không gọi Media/SEO API trong Create.
- `onMutationError` giữ nguyên title, slug, content và status đang chọn; dùng notification provider để hiển thị lỗi.

### 6.2. Validation client-side

| Field | Quy tắc |
| --- | --- |
| `title` | Trim không rỗng, tối đa 255 ký tự. |
| `slug` | Trim không rỗng, tối đa 255 ký tự, regex `^[a-z0-9]+(?:-[a-z0-9]+)*$`. |
| `content` | Sau khi loại HTML rỗng/whitespace phải còn nội dung. |

- Validation form chạy `onBlur`; submit luôn trigger lại toàn bộ validation.
- Backend là nguồn xác thực cuối cùng. Khi `409` trả fields cho `slug`, adapter/mutation map lỗi đó vào `setError("slug", ...)`; lỗi chung hiển thị Snackbar/Alert.
- Không gọi API kiểm tra slug theo từng lần gõ.

### 6.3. Auto-slug

- Theo dõi `title` và cờ `isSlugManuallyEdited`.
- Chỉ gọi `generateSlug(title)` để set slug khi cờ này là `false`.
- Bất kỳ thay đổi trực tiếp nào vào slug đặt cờ thành `true`; sau đó title không ghi đè slug.
- Nếu người dùng xóa slug hoàn toàn, đặt cờ về `false` để title tiếp tục sinh slug. Không tự chuyển đổi lại slug khi submit.

## 7. API contract prerequisite

Frontend phụ thuộc endpoint bảo vệ:

```text
POST /api/v1/admin/pages
```

Request:

```json
{
  "title": "Giới thiệu",
  "slug": "gioi-thieu",
  "content": "<p>Nội dung trang</p>",
  "status": "DRAFT"
}
```

Response thành công `201`:

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

- API phải nằm dưới AuthMiddleware và RoleMiddleware `ADMIN`/`STAFF`.
- `400` cho payload không hợp lệ; `401/403` để AuthProvider xử lý theo luồng hiện có; `409` cho slug trùng với lỗi field; `500` cho lỗi hệ thống.
- Backend Create Page chưa có trong backend plan `22-page-list-plan.md`; cần có backend plan/implementation Create trước khi UI này được tích hợp thật.

## 8. Luồng submit và trạng thái

1. Admin nhập form; title sinh slug nếu slug chưa được chỉnh tay.
2. Nhấn **Lưu nháp** hoặc **Xuất bản**; frontend validate rồi gọi `POST /admin/pages` với status tương ứng.
3. Khi request pending, khóa action và hiển thị feedback đang lưu.
4. Thành công: hiển thị toast theo status, reset dirty state, invalidate resource `pages`, chuyển `/pages/edit/:id`.
5. Lỗi: giữ form/dirty state; field error hiển thị cạnh input; API/general error hiển thị thông báo. Không tạo hoặc sync Media/SEO.

## 9. Acceptance test

- Route `/pages/create` chỉ render khi authenticated; resource `pages` dùng đúng endpoint admin.
- Form có đủ title/slug/content, hai action status và layout responsive.
- Title auto-generate slug; slug chỉnh tay không bị title ghi đè; slug rỗng quay về auto-generation.
- Validation chặn submit với title/slug/content không hợp lệ; không có request slug-check theo mỗi keystroke.
- Lưu nháp gửi `DRAFT`, xuất bản gửi `PUBLISHED`; không gửi `authorId`, Media hoặc SEO payload.
- `409` hiển thị lỗi tại slug; lỗi khác giữ dữ liệu người dùng; `201` toast, invalidate list và redirect đúng ID.
- Media/SEO cards trước khi tạo chỉ là trạng thái hướng dẫn, không phát sinh request.
- TypeScript build không có `any` ở module Page Create và UI components không gọi API trực tiếp.

## 10. Phụ thuộc và giả định

- `PageList`/resource plan đã khai báo URL `/pages`, `/pages/create` và `/pages/edit/:id`.
- Backend sẽ cung cấp Create Page contract, Entity/Page status và role guard đúng như tài liệu IDEA.
- Rich Text Editor từ Post được tái sử dụng nguyên component/cấu hình; không bổ sung editor library mới.
- Page Media, Page SEO và Edit Page được triển khai sau, dùng ID từ response Create thay vì cố lưu dữ liệu phụ trước khi Page tồn tại.
