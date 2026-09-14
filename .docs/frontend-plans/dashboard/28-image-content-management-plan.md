# FRONTEND PLAN: Quản lý nội dung hình ảnh theo Type

**Tài liệu nguồn:** `.docs/ideas/img-management-with-type.md`

**Backend contract:** `.docs/backend-plans/img-management-with-type-plan.md`, `.docs/api-endpoints.yaml`

**Ứng dụng:** `apps/frontend`

**Stack:** React 19, TypeScript, Refine.js, React Hook Form, Material UI, MUI DataGrid

---

## 1. Mục tiêu

- Xây dựng hai khu vực quản trị độc lập cho domain `ImageContent`:
  - **Loại nội dung hình ảnh**: ADMIN cấu hình Type, giới hạn item và các metadata field được bật/bắt buộc.
  - **Nội dung hình ảnh**: ADMIN/STAFF quản lý ảnh và metadata theo Type.
- Form Item phải sinh động từ `fieldConfig`; không hard-code riêng LOGO, SLIDER hoặc PARTNER.
- Tái sử dụng Media Library và API upload hiện có; frontend chỉ lưu `mediaId`, không tạo luồng upload mới.
- Hỗ trợ list/search/filter/sort, create/edit/delete, reorder snapshot, cảnh báo dữ liệu chưa lưu và phản hồi lỗi nghiệp vụ từ Backend.
- Tuân thủ Refine resource-driven, Smart/Dumb Component, TypeScript strict, không dùng `any`, đủ Loading/Error/Empty/Success.

---

## 2. Hiện trạng và phạm vi

### 2.1. Hiện trạng liên quan

- `App.tsx` đã đăng ký các resource chuẩn và hỗ trợ nhiều Data Provider (`default`, `pageSections`).
- `customRequest` đã xử lý Access Token, Refresh Token và request JSON/FormData tập trung.
- `useDebounce` dùng lodash với mặc định `300ms` đã tồn tại tại `apps/frontend/src/hooks/use-debounce.ts`.
- Media Library có endpoint list `/api/v1/admin/media`; upload dùng `/api/v1/media/upload`.
- Page Section đã có picker ảnh server-side và có thể tái sử dụng pattern chọn Media, nhưng cần tách component dùng chung nếu việc tái sử dụng không làm phụ thuộc ngược vào module Page.
- Backend ImageContent đã có CRUD Type, CRUD Item, reorder và Public API. Admin UI chưa có resource, route hoặc component tương ứng.
- Chưa có mockup riêng cho tính năng này; UI lấy `.docs/STYLEGUIDE.md`, MUI Theme và các màn hình Post Type/Page làm nguồn thiết kế.

### 2.2. Trong phạm vi

- Resource, route và menu cho `image-content-types` và `image-contents`.
- Type List/Create/Edit và confirmation delete.
- Item List/Create/Edit, Media picker/upload, confirmation delete và reorder theo Type.
- Data layer typed cho API ImageContent, normalizer response và typed error mapping.
- Role-aware UI: Type chỉ ADMIN được quản trị; Item dành cho ADMIN/STAFF.
- Loading, skeleton, retry, empty state, snackbar, unsaved changes, responsive và accessibility.

### 2.3. Ngoài phạm vi

- Không thay đổi `apps/webview` hoặc render LOGO/SLIDER/PARTNER ngoài public site.
- Không xây editor field schema tùy ý ngoài bốn field `name`, `description`, `secondaryDescription`, `url`.
- Không tạo/xóa file Media trực tiếp khi xóa ImageContent Item.
- Không seed dữ liệu, sửa migration hoặc thay đổi hành vi Post/Page/Page Section.
- Không hỗ trợ bulk delete, revision, scheduling, locale, A/B testing hoặc crop ảnh trong task này.

---

## 3. Backend prerequisites cần chốt trước bước tích hợp

### 3.1. STAFF cần Type catalog dạng read-only

Item form cần `code`, `name`, `fieldConfig`, `maxItems` và `itemCount` của mọi Type, kể cả Type chưa có Item. Public API không trả `fieldConfig`, không trả Type INACTIVE và không phù hợp cho Admin form.

Backend hiện đặt `RoleMiddleware("ADMIN")` trên toàn bộ group `/admin/image-content-types`, vì vậy STAFF không thể tạo Item động theo Type. Trước khi tích hợp cần chọn một trong hai phương án, ưu tiên phương án đầu:

1. Cho `ADMIN/STAFF` gọi `GET /admin/image-content-types` và `GET /admin/image-content-types/:id`; giữ `POST/PUT/DELETE` chỉ ADMIN.
2. Tạo endpoint read-only riêng `/admin/image-content-type-options` cho ADMIN/STAFF với đủ `fieldConfig`, `maxItems`, `itemCount` và status.

Frontend tuyệt đối không hard-code cấu hình LOGO/SLIDER/PARTNER để né thiếu contract này.

### 3.2. Bảo toàn metadata của field đang disabled

- Admin detail trả raw metadata để có thể khôi phục khi Type bật lại field.
- Form Edit chỉ cho sửa field đang enabled; field disabled phải được hiển thị read-only trong khu vực **Dữ liệu đang tạm ẩn** nếu có giá trị.
- Update API phải giữ nguyên giá trị field disabled khi request không gửi field đó. Backend cần phân biệt omitted với explicit `null`, hoặc chỉ update các field enabled.
- Frontend không gửi raw giá trị disabled vào payload vì contract hiện từ chối `FIELD_NOT_ENABLED`; cũng không gửi `null` nếu điều đó làm mất dữ liệu đang lưu.

### 3.3. Error code và response ổn định

- Backend nên trả đúng code chi tiết: `TYPE_NOT_FOUND`, `IMAGE_CONTENT_NOT_FOUND`, `MEDIA_NOT_FOUND`, `TYPE_CODE_CONFLICT`, `TYPE_CONFIG_CONFLICT`, `MAX_ITEMS_EXCEEDED`, `FIELD_REQUIRED`, `FIELD_NOT_ENABLED`, `MEDIA_TYPE_INVALID`, `FORBIDDEN`, `ORDER_CONFLICT`.
- Data layer không suy luận nghiệp vụ từ chuỗi `message`; `code` là nguồn quyết định UI, `message` dùng để hiển thị.
- DELETE `204` không có body; provider phải xử lý trực tiếp, không gọi `response.json()`.

---

## 4. Information architecture, routing và RBAC

### 4.1. Refine resources

```ts
{
  name: "image-content-types",
  list: "/image-content-types",
  create: "/image-content-types/create",
  edit: "/image-content-types/edit/:id",
  meta: {
    label: "Loại nội dung hình ảnh",
    canDelete: true,
    requiredRoles: ["ADMIN"],
  },
},
{
  name: "image-contents",
  list: "/image-contents",
  create: "/image-contents/create",
  edit: "/image-contents/edit/:id",
  meta: {
    label: "Nội dung hình ảnh",
    canDelete: true,
    requiredRoles: ["ADMIN", "STAFF"],
  },
}
```

- Dùng tên resource có dấu gạch ngang để ánh xạ trực tiếp tới `/api/v1/admin/image-content-types` và `/api/v1/admin/image-contents`.
- Đăng ký route tương ứng trong vùng `Authenticated` của `App.tsx`.
- Type menu và route chỉ hiển thị/cho vào khi identity có role `ADMIN`.
- Item menu và route cho `ADMIN`, `STAFF`; `CUSTOMER` không hiển thị.
- Route guard phía frontend chỉ phục vụ UX. Backend RBAC vẫn là lớp bảo vệ quyết định.

### 4.2. Điều hướng giữa Type và Item

- Từ Type List có action **Quản lý nội dung** dẫn tới:

```text
/image-contents?typeCode=SLIDER
```

- Từ Item List, nút **Thêm nội dung** giữ Type đang lọc:

```text
/image-contents/create?typeCode=SLIDER
```

- Edit Item không cho đổi Type vì Backend không hỗ trợ cập nhật `typeCode`; muốn đổi Type phải xóa và tạo Item mới.

---

## 5. Cấu trúc file dự kiến

```text
apps/frontend/src/
├── App.tsx
├── components/
│   └── media/
│       ├── image-media-picker-dialog.tsx
│       ├── image-media-upload-field.tsx
│       └── media-preview.tsx
├── providers/
│   └── image-content-data-provider.ts
└── pages/
    ├── image-content-types/
    │   ├── index.ts
    │   ├── list.tsx
    │   ├── create.tsx
    │   ├── edit.tsx
    │   ├── image-content-type-types.ts
    │   ├── image-content-type-normalizers.ts
    │   ├── image-content-type-validation.ts
    │   └── components/
    │       ├── image-content-type-form.tsx
    │       ├── field-config-matrix.tsx
    │       ├── type-status-chip.tsx
    │       ├── type-capacity-chip.tsx
    │       └── type-delete-dialog.tsx
    └── image-contents/
        ├── index.ts
        ├── list.tsx
        ├── create.tsx
        ├── edit.tsx
        ├── image-content-types.ts
        ├── image-content-api.ts
        ├── image-content-normalizers.ts
        ├── image-content-validation.ts
        ├── use-image-content-form.ts
        ├── use-image-content-reorder.ts
        └── components/
            ├── image-content-form.tsx
            ├── dynamic-metadata-fields.tsx
            ├── disabled-metadata-summary.tsx
            ├── image-content-list-toolbar.tsx
            ├── image-content-status-chip.tsx
            ├── image-content-delete-dialog.tsx
            ├── image-content-reorder-dialog.tsx
            └── image-content-empty-state.tsx
```

Phân trách nhiệm:

- Page `list/create/edit` là Smart Component, sở hữu Refine hooks, route/query state và mutation state.
- Form, matrix, chip, preview, dialog là Dumb Component nhận typed props/callback.
- `image-content-data-provider.ts` là gateway duy nhất cho CRUD chuẩn và typed API error.
- `image-content-api.ts` chỉ chứa custom action chưa thuộc CRUD chuẩn: reorder, tải toàn bộ snapshot theo Type và upload Media.
- Normalizer nhận `unknown`, dùng type guard và trả model chuẩn; component không ép kiểu response trực tiếp.
- Nếu Media picker được tách dùng chung, Page Section/Post Media chuyển sang dùng chung ở task refactor riêng; không làm tăng scope triển khai ImageContent.

---

## 6. TypeScript contracts

### 6.1. Type và field configuration

```ts
export type ImageContentStatus = "ACTIVE" | "INACTIVE";
export type ImageContentFieldKey =
  | "name"
  | "description"
  | "secondaryDescription"
  | "url";

export interface IImageContentFieldRule {
  enabled: boolean;
  required: boolean;
}

export interface IImageContentFieldConfig {
  name: IImageContentFieldRule;
  description: IImageContentFieldRule;
  secondaryDescription: IImageContentFieldRule;
  url: IImageContentFieldRule;
}

export interface IImageContentType {
  id: number;
  code: string;
  name: string;
  status: ImageContentStatus;
  sortOrder: number;
  maxItems: number | null;
  fieldConfig: IImageContentFieldConfig;
  itemCount: number;
  createdAt: string;
  updatedAt: string;
}

export interface IImageContentTypeFormValues {
  code: string;
  name: string;
  status: ImageContentStatus;
  sortOrder: number;
  hasItemLimit: boolean;
  maxItems: number | null;
  fieldConfig: IImageContentFieldConfig;
}

export interface ICreateImageContentTypePayload {
  code: string;
  name: string;
  status: ImageContentStatus;
  sortOrder: number;
  maxItems: number | null;
  fieldConfig: IImageContentFieldConfig;
}

export type IUpdateImageContentTypePayload = Omit<
  ICreateImageContentTypePayload,
  "code"
>;
```

- `hasItemLimit` chỉ tồn tại trong form; payload map `false` thành `maxItems: null`.
- `code` được trim + uppercase khi tạo và không xuất hiện trong Update payload.
- Dùng `ImageContentFieldKey[]` cố định cho matrix; không dùng `Record<string, ...>` nhận key tùy ý từ UI.

### 6.2. Item và Media

```ts
export interface IImageContentMedia {
  id: number;
  fileName: string;
  originalUrl: string;
  thumbnailUrl?: string;
  mediumUrl?: string;
  mimeType: string;
  status: "temporary" | "attached";
}

export interface IImageContent {
  id: number;
  typeCode: string;
  typeName: string;
  fieldConfig: IImageContentFieldConfig;
  mediaId: number;
  media: IImageContentMedia;
  name: string | null;
  description: string | null;
  secondaryDescription: string | null;
  url: string | null;
  sortOrder: number;
  status: ImageContentStatus;
  createdBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface IImageContentFormValues {
  type: IImageContentType | null;
  media: IImageContentMedia | null;
  name: string;
  description: string;
  secondaryDescription: string;
  url: string;
  status: ImageContentStatus;
}

export interface ICreateImageContentPayload {
  typeCode: string;
  mediaId: number;
  name: string | null;
  description: string | null;
  secondaryDescription: string | null;
  url: string | null;
  status: ImageContentStatus;
}

export interface IUpdateImageContentPayload {
  mediaId: number;
  name?: string | null;
  description?: string | null;
  secondaryDescription?: string | null;
  url?: string | null;
  status: ImageContentStatus;
}

export interface IImageContentOrderPayload {
  typeCode: string;
  items: Array<{ id: number; sortOrder: number }>;
}
```

- Form giữ object `type/media` để render label và preview; payload chỉ gửi identifier.
- Chuỗi enabled được trim; chuỗi rỗng map thành `null`.
- Field disabled bị loại khỏi Update payload sau khi backend đáp ứng quy tắc bảo toàn dữ liệu tại mục 3.2.
- Preview URL dùng `thumbnailUrl ?? mediumUrl ?? originalUrl` qua helper resolve URL dùng `BACKEND_URL`.

### 6.3. Typed API error

```ts
export type ImageContentErrorCode =
  | "VALIDATION_ERROR"
  | "TYPE_NOT_FOUND"
  | "IMAGE_CONTENT_NOT_FOUND"
  | "MEDIA_NOT_FOUND"
  | "TYPE_CODE_CONFLICT"
  | "TYPE_CONFIG_CONFLICT"
  | "MAX_ITEMS_EXCEEDED"
  | "FIELD_REQUIRED"
  | "FIELD_NOT_ENABLED"
  | "MEDIA_TYPE_INVALID"
  | "FORBIDDEN"
  | "ORDER_CONFLICT"
  | "INTERNAL_ERROR";

export interface IImageContentApiError {
  error: string;
  code: ImageContentErrorCode;
  message: string;
  fields?: Partial<Record<ImageContentFieldKey | "code" | "maxItems", string>>;
}
```

- `readImageContentResponse` kiểm tra status và normalize body `unknown` thành `HttpError` mở rộng có `code/fields`.
- UI ưu tiên field error từ `fields`; nếu không có thì hiển thị snackbar theo `code`.

---

## 7. Data Provider và data flow

### 7.1. Named Data Provider

- Đăng ký `imageContentDataProvider` trong `App.tsx` bên cạnh `default` và `pageSections`.
- Các page/hook của module truyền `dataProviderName: "imageContent"`.
- Provider triển khai `getList`, `getOne`, `create`, `update`, `deleteOne` cho đúng hai resource.
- Query list map chuẩn Refine sang `_start`, `_end`, `_sort`, `_order`, `q`, `status`, `typeCode`.
- List trả `{ data, total }`; detail/create/update unwrap `{ data }` đúng một lần.
- Không dùng `any`; generic mặc định là `BaseRecord`, dữ liệu raw là `unknown`.
- DELETE xử lý `204` như success và trả record tối thiểu `{ id }` cho Refine.

### 7.2. Search, filter và URL state

- Search dùng `useDebounce(search, 300)` trước khi gọi `setFilters`.
- Filter Type/Status và pagination đồng bộ URL bằng `syncWithLocation`.
- Khi đổi `typeCode`, reset trang về 1.
- Type List search theo code/name; Item List search theo name/description/secondaryDescription/url.
- DataGrid dùng server pagination và server sorting; sort field chỉ dùng whitelist backend hỗ trợ.

### 7.3. Reorder snapshot

- Reorder API yêu cầu toàn bộ Item chưa xóa của một Type, gồm ACTIVE và INACTIVE.
- Không dùng riêng các row của trang DataGrid hiện tại.
- Khi mở Reorder Dialog, helper tải tuần tự các page kích thước 100 theo `typeCode`, `_sort=sortOrder`, `_order=ASC` cho tới khi đủ `total` hoặc chạm 1000.
- Merge phải kiểm tra ID trùng, tổng row và TypeCode; response thay đổi trong lúc tải thì refetch từ đầu một lần trước khi báo lỗi.
- Chỉ gửi một request `PUT /admin/image-contents/order` sau khi người dùng bấm **Lưu thứ tự**.

---

## 8. Màn hình Loại nội dung hình ảnh

### 8.1. Type List

Header:

- Tiêu đề **Loại nội dung hình ảnh**.
- Mô tả ngắn về việc Type điều khiển field hiển thị của Item.
- Nút **Thêm loại nội dung** chỉ dành cho ADMIN.

Toolbar:

- Search code/tên có debounce 300ms.
- Filter status: Tất cả, ACTIVE, INACTIVE.
- Nút reset filter.

DataGrid `density="medium"` gồm:

- `code`: monospace Chip hoặc Typography, không chỉnh inline.
- `name`.
- `fieldConfig`: các Chip nhỏ như `Tên · bắt buộc`, `URL · tùy chọn`; field disabled không hiển thị hoặc gom trong tooltip.
- `itemCount/maxItems`: `3 / 10`, `4 / Không giới hạn`; cảnh báo khi đạt giới hạn.
- `status`: semantic Chip.
- `sortOrder`.
- `updatedAt`.
- actions: **Quản lý nội dung**, **Sửa**, **Xóa**.

States:

- Loading: skeleton/DataGrid loading overlay.
- Empty do chưa có Type: CTA tạo Type đầu tiên.
- Empty do filter: thông báo không có kết quả và nút reset filter.
- Error: Alert + **Thử lại**; không hiển thị empty state giả.

### 8.2. Type Create/Edit

Layout desktop hai cột, mobile một cột:

- Card **Thông tin chung**: code, name, sortOrder.
- Card **Trạng thái và giới hạn**: status, toggle giới hạn, maxItems.
- Card full-width **Cấu hình metadata**: matrix Enabled/Required.
- Sticky action footer: **Hủy**, **Lưu**.

Field Config Matrix:

| Field | Label UI | Enabled | Required |
| --- | --- | --- | --- |
| `name` | Tên | Switch | Checkbox/Switch |
| `description` | Mô tả chính | Switch | Checkbox/Switch |
| `secondaryDescription` | Mô tả phụ | Switch | Checkbox/Switch |
| `url` | Liên kết | Switch | Checkbox/Switch |

Quy tắc tương tác:

- Required disabled khi Enabled=false.
- Tắt Enabled tự đặt Required=false trong form state.
- Edit khóa `code` và không đưa code vào payload.
- `maxItems`: null khi không giới hạn; nếu bật phải là số nguyên `1..1000`.
- Khi Edit giảm max hoặc bật Required, hiển thị warning rằng dữ liệu hiện có có thể gây conflict.
- `TYPE_CONFIG_CONFLICT` giữ form dirty, không redirect, hiển thị message/field details.
- `TYPE_CODE_CONFLICT` focus field code.
- Delete mở confirmation; `TYPE_IN_USE` giải thích cần xóa Item trước và cung cấp link sang Item List đã lọc Type.

Validation client:

- Code sau normalize phải khớp `^[A-Z][A-Z0-9_]{0,49}$`.
- Name trim, required, tối đa 100 ký tự.
- Sort order là số nguyên `>= 0`.
- Required chỉ hợp lệ khi Enabled.
- Client validation hỗ trợ UX nhưng không thay thế Backend validation.

---

## 9. Màn hình Nội dung hình ảnh

### 9.1. Item List

Header:

- Tiêu đề **Nội dung hình ảnh**.
- Nút **Thêm nội dung**.
- Nút **Sắp xếp** chỉ enable khi đã chọn đúng một Type và Type có từ hai Item trở lên.

Toolbar:

- Type Select/Autocomplete lấy từ Type catalog.
- Status filter: Tất cả, ACTIVE, INACTIVE.
- Search metadata có debounce 300ms.
- Capacity summary của Type đang chọn: `itemCount/maxItems`.

DataGrid:

- Preview ảnh 56×40 hoặc 64×44, `objectFit="cover"`, rounded theo Theme.
- Tên; nếu Name disabled/rỗng dùng `typeName #id` làm display fallback.
- Type name + code.
- Mô tả ngắn, ellipsis và tooltip.
- URL dạng link an toàn; external mở tab mới với `rel="noopener noreferrer"`, internal chỉ hiển thị text/icon nếu Dashboard không có route tương ứng.
- Status Chip.
- Sort order.
- Updated time.
- Actions: Edit/Delete.

Create capacity:

- Chỉ disable CTA khi frontend biết chắc `maxItems !== null && itemCount >= maxItems`.
- Nếu stale và Backend trả `MAX_ITEMS_EXCEEDED`, refetch Type catalog, giữ form và hiển thị snackbar.
- Item INACTIVE vẫn được tính vào capacity.

### 9.2. Item Create/Edit

Luồng Create:

1. Chọn Type; nếu URL có `typeCode`, preselect sau khi Type catalog tải xong.
2. Hiển thị summary Type: status, capacity và field đang bật.
3. Chọn/upload ảnh.
4. Render metadata field theo `fieldConfig`.
5. Chọn status và submit.

Luồng Edit:

- Tải Item detail, dùng `fieldConfig` trong response làm snapshot ban đầu.
- Đồng thời lấy Type catalog/detail mới nhất; nếu config khác snapshot, hiển thị info alert và dùng config mới nhất cho validation.
- Type hiển thị read-only.
- Giữ `sortOrder` read-only; đổi vị trí qua Reorder Dialog.
- Raw metadata của field disabled hiển thị trong `DisabledMetadataSummary`, không cho sửa và không gửi vào payload.

Dynamic fields:

- `name`: TextField, maxLength 255.
- `description`: multiline TextField, maxLength 5000, có bộ đếm ký tự.
- `secondaryDescription`: multiline TextField, maxLength 5000.
- `url`: TextField, maxLength 500; nhận internal path bắt đầu bằng một `/` hoặc absolute `http/https`; từ chối `//`, `javascript:`, `data:` và control character.
- Field `enabled=false` không render trong vùng editable.
- Field `required=true` có dấu `*`, React Hook Form rule required và helper text.
- Khi đổi Type trong Create, unregister/reset metadata field không còn enabled để tránh gửi dữ liệu stale.

### 9.3. Media picker và upload

- `ImageMediaPickerDialog` single-select, server pagination và search debounce 300ms.
- Chỉ lấy `mimeType=image`; row hiển thị preview, fileName, MIME, kích thước nếu response có.
- Selected Media được giữ khi đổi page/search.
- Upload gọi API Media hiện tại bằng FormData; sau success chọn ngay Media vừa upload.
- Preview dùng URL variant ưu tiên thumbnail → medium → original.
- Nút **Đổi ảnh** không xóa Media cũ.
- Nếu `MEDIA_NOT_FOUND` hoặc `MEDIA_TYPE_INVALID`, giữ form và yêu cầu chọn lại.
- STAFF không tự lọc quyền sở hữu dựa vào dữ liệu không đầy đủ; Backend quyết định. Khi `FORBIDDEN`, hiển thị thông báo Media không thuộc quyền sử dụng và mở lại picker.

### 9.4. Delete Item

- Confirmation nêu rõ chỉ xóa liên kết/nội dung quản trị, không xóa file Media.
- Success đóng dialog, refetch Item List và Type catalog để cập nhật `itemCount`.
- Error giữ row và dialog state nhất quán, không optimistic remove nếu chưa có rollback rõ ràng.

---

## 10. Reorder UX

- Reorder chỉ hoạt động trong một Type; TypeCode hiển thị read-only trong dialog.
- Dialog `maxWidth="md"`, danh sách ảnh + tên fallback + status + vị trí.
- Hỗ trợ drag/drop nếu dùng API native hoặc thư viện đã có; luôn có nút lên/xuống cho keyboard/mobile.
- Drag/drop chỉ cập nhật local state, không request theo từng thao tác.
- Chuẩn hóa `sortOrder` liên tục `0..n-1` trước submit.
- Footer có **Hoàn tác**, **Hủy**, **Lưu thứ tự**.
- Trong lúc save, khóa thao tác và chống double-submit.
- `ORDER_CONFLICT`: giữ dialog mở, thông báo dữ liệu đã thay đổi và cung cấp nút **Tải lại thứ tự mới**; không tự ghi đè snapshot stale.
- Success refetch Item List + Type catalog, reset dirty state và snackbar.
- Đóng dialog khi order dirty phải confirmation.

---

## 11. Error, loading và notification strategy

Mapping tối thiểu:

| Backend code | Hành vi UI |
| --- | --- |
| `VALIDATION_ERROR` | Giữ form, map `fields`, focus field đầu tiên |
| `TYPE_CODE_CONFLICT` | Error tại code, không redirect |
| `TYPE_CONFIG_CONFLICT` | Alert trong Type form, giữ dữ liệu đã nhập |
| `TYPE_IN_USE` | Dialog hướng dẫn mở danh sách Item theo Type |
| `MAX_ITEMS_EXCEEDED` | Refetch capacity, disable create nếu đã đầy |
| `FIELD_REQUIRED` | Error đúng dynamic field |
| `FIELD_NOT_ENABLED` | Refetch Type config và rebuild form |
| `MEDIA_NOT_FOUND` | Clear Media selection nếu Media đã bị xóa |
| `MEDIA_TYPE_INVALID` | Báo chỉ chấp nhận ảnh |
| `FORBIDDEN` | Snackbar quyền hạn; không retry tự động mutation |
| `ORDER_CONFLICT` | Giữ reorder dialog và yêu cầu reload snapshot |
| `401` | Dùng auth refresh/logout flow hiện có |
| `500` | Snackbar lỗi hệ thống và cho retry an toàn |

- Không hiển thị raw SQL, stack trace hoặc toàn bộ response object.
- Query retry có giới hạn; mutation không tự retry để tránh ghi lặp.
- Success snackbar riêng cho create/update/delete/reorder.

---

## 12. Unsaved changes, accessibility và responsive

### 12.1. Unsaved changes

- Create/Edit tích hợp `warnWhenUnsavedChanges` hiện có.
- Hủy hoặc back khi dirty mở confirmation.
- Media upload thành công nhưng chưa submit Item vẫn được coi là form dirty; không tự xóa file khi người dùng rời trang trong scope này.
- Reorder dialog có dirty state độc lập.

### 12.2. Accessibility

- Mọi icon action có tooltip và `aria-label` cụ thể.
- Ảnh preview có alt từ `name`, fallback `typeName #id` hoặc fileName.
- Switch Enabled/Required có label liên kết; không truyền nghĩa chỉ bằng màu.
- Focus trở lại nút mở dialog sau khi đóng.
- Reorder có nút keyboard, thông báo vị trí và live region khi di chuyển nếu triển khai drag/drop.

### 12.3. Responsive và Theme

- Desktop form hai cột; tablet/mobile một cột.
- DataGrid giữ cột ảnh/name/type/status/action; cột mô tả/time có thể ẩn ở breakpoint nhỏ.
- Dialog full-screen trên mobile khi form/picker dài.
- Chỉ dùng palette, typography, spacing, border và shadow từ MUI Theme/STYLEGUIDE; không thêm mã HEX mới.
- DataGrid `density="medium"`, toolbar wrap trên mobile.

---

## 13. Trình tự triển khai

1. Chốt ba backend prerequisites tại mục 3.
2. Tạo TypeScript contracts, normalizer và typed error parser.
3. Tạo/register `imageContentDataProvider` và Refine resources/routes.
4. Triển khai role-aware menu/route guard.
5. Triển khai Type List/Create/Edit/Delete.
6. Tách hoặc tạo Media picker dùng cho ImageContent.
7. Triển khai Item List/Create/Edit/Delete với dynamic fields.
8. Triển khai fetch-all snapshot và Reorder Dialog.
9. Hoàn thiện loading/error/empty/unsaved/a11y/responsive.
10. Chạy lint/typecheck/build và kiểm thử các role/API error.

---

## 14. Checklist kiểm thử

### 14.1. Type

- [ ] ADMIN xem/tạo/sửa/xóa Type; STAFF/CUSTOMER không thấy menu và bị route guard chặn.
- [ ] Search debounce, status filter, pagination và sorting đồng bộ URL.
- [ ] Code normalize uppercase, regex đúng và khóa khi edit.
- [ ] Tắt Enabled tự tắt Required.
- [ ] maxItems null/1/1000 map payload đúng; 0 và >1000 bị chặn.
- [ ] `TYPE_CODE_CONFLICT`, `TYPE_CONFIG_CONFLICT`, `TYPE_IN_USE` hiển thị đúng ngữ cảnh.
- [ ] List cập nhật `itemCount` sau mutation Item.

### 14.2. Item

- [ ] ADMIN/STAFF đọc Type catalog và quản lý Item; CUSTOMER bị chặn.
- [ ] Create preselect Type từ query string; Type INACTIVE vẫn quản trị được.
- [ ] LOGO đạt một Item thì CTA bị disable; Backend conflict vẫn được xử lý nếu state stale.
- [ ] Dynamic form chỉ render field enabled và bắt required đúng config.
- [ ] Field disabled có raw dữ liệu không bị mất sau update field khác.
- [ ] URL internal/http/https hợp lệ; protocol-relative/javascript/data/control character bị chặn.
- [ ] Picker chỉ hiển thị ảnh, upload dùng Media API hiện có.
- [ ] STAFF chọn Media không sở hữu nhận FORBIDDEN rõ ràng.
- [ ] Xóa Item không gọi delete Media và capacity được cập nhật.

### 14.3. Reorder và trạng thái chung

- [ ] Reorder tải toàn bộ Item qua nhiều page, gồm ACTIVE và INACTIVE.
- [ ] Snapshot gửi đủ ID, không trùng, sortOrder liên tục.
- [ ] `ORDER_CONFLICT` không làm mất local order và có reload.
- [ ] Loading/Error/Empty/Success không nhấp nháy hoặc chồng trạng thái.
- [ ] Unsaved warning hoạt động cho Type form, Item form và reorder.
- [ ] Keyboard, screen reader label và mobile layout dùng được.

### 14.4. Lệnh xác minh

```text
cd apps/frontend
npm run build
```

- TypeScript compile không lỗi.
- Không phát sinh `any` trong các file mới.
- Không có raw `fetch` ngoài provider/API layer được chỉ định.
- Không có mã màu hard-code mới ngoài dữ liệu nội dung do người dùng nhập.

---

## 15. Tiêu chí nghiệm thu

- ADMIN cấu hình được Type mới từ bốn field chuẩn mà không cần sửa frontend.
- ADMIN/STAFF tạo và quản lý Item theo `fieldConfig`, dùng Media hiện có và nhận đúng validation/error.
- Capacity tính cả Item INACTIVE, LOGO không vượt quá một Item.
- Reorder gửi full snapshot và xử lý concurrent conflict an toàn.
- Field disabled không bị lộ trong vùng editable và không làm mất raw metadata khi cập nhật Item.
- List/detail/form đáp ứng Refine contract, đủ loading/error/empty và không dùng `any`.
- Không thay đổi Webview; tích hợp Public API vào giao diện public được thực hiện ở task riêng.
