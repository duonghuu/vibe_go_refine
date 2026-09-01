# FRONTEND PLAN: Quản lý Section của từng Page

## 1. Mục tiêu

- Triển khai khu vực **Quản lý Section** trong `/pages/edit/:id` theo `.docs/ideas/dashboard/26-page-section-management-idea.md`.
- Cho phép Admin/Staff tạo, sửa, xóa, bật/tắt và sắp xếp các Section của một Page đã tồn tại.
- Cho phép Section cấu hình metadata chung, Background Media, Feature Media và các collection liên kết `CATEGORY`, `POST`, `MEDIA`.
- Giữ mô hình Semi-Structured CMS: Frontend chỉ gửi ID và thứ tự liên kết; không sao chép dữ liệu Category/Post/Media vào payload Section.
- Tuân thủ Refine.js, React Hook Form, MUI, TypeScript strict, Smart/Dumb Component và đủ Loading/Error/Empty/Success.

---

## 2. Hiện trạng và phạm vi

### 2.1. Hiện trạng liên quan

- `PageEdit` đã tồn tại tại `apps/frontend/src/pages/pages/edit.tsx`, tải Page bằng `useForm` và điều phối Page core, Page Media, SEO.
- Page Media/SEO hiện đang tái sử dụng các module mang tên Post (`PostMediaCard`, `usePostMedia`, `PostSeoCard`, `usePostSeo`). Tính năng Section không refactor lại các module này nếu không cần thiết.
- Frontend chưa có module Page Section, chưa có Media Library list/search và Backend chưa có các endpoint Section trong Idea.
- API Category và Post đã có list phân trang/tìm kiếm. API Media hiện chỉ có upload/delete và API liên kết Post/Page Media, chưa đủ cho popup chọn Media có sẵn.

### 2.2. Trong phạm vi

- Nhúng Section Manager vào Page Edit sau khu vực Page core/Media/SEO hiện tại.
- Danh sách Section dạng accordion/card, hiển thị metadata, status, Media preview và số Item theo collection.
- Form tạo/sửa Section bằng dialog có Stepper: thông tin chung → giao diện/Media → xác nhận.
- Reorder Section bằng pointer và nút lên/xuống; chỉ gửi một snapshot khi người dùng bấm lưu thứ tự.
- Confirmation Dialog khi xóa Section; optimistic status toggle có rollback khi lỗi.
- Collection Manager: xem Item, mở picker, chọn nhiều, bỏ chọn, reorder và đồng bộ snapshot nguyên tử.
- Item Picker dùng server pagination, debounce tìm kiếm `300ms` và giữ selection khi đổi trang/tìm kiếm.
- Media Picker dùng chung cho Background Media, Feature Media và collection `MEDIA`; hỗ trợ single-select hoặc multi-select theo ngữ cảnh.
- Tích hợp dirty state của Section vào cảnh báo rời Page Edit.

### 2.3. Ngoài phạm vi

- Không xây Page Builder kéo thả tự do, preview Webview trực tiếp hoặc render component theo `section.key` trong Dashboard.
- Không tạo Section trong `/pages/create`; màn hình Create hiển thị notice **Sections — Lưu trang trước để quản lý các Section**.
- Không hỗ trợ Section lồng nhau, revision, scheduled publishing, đa ngôn ngữ, A/B testing hoặc JSON layout tự do.
- Không sửa Page public/Webview; Frontend Admin chỉ quản lý dữ liệu cho Public API sử dụng sau.
- Không tự tạo, sửa hoặc xóa Category/Post/Media nguồn từ trong Item Picker.
- Không thêm loại Item ngoài `CATEGORY`, `POST`, `MEDIA`.

---

## 3. Điểm tích hợp, routing và bố cục

### 3.1. Routing

- Không tạo Refine route/resource hiển thị độc lập cho Section.
- Giữ route hiện tại:

```text
/pages/edit/:id
```

- Section là nested sub-resource của Page và chỉ được tải sau khi `pageId` hợp lệ, Page detail tải thành công.
- Page `404` hoặc soft-delete vẫn là lỗi chặn toàn màn hình; không gọi Section API trong trường hợp này.
- `PageCreate` mở rộng `PagePostCreateNotice` với icon/label Section và giữ hành vi redirect sang `/pages/edit/:id` sau khi tạo Page thành công.

### 3.2. Vị trí trong Page Edit

- Đặt `PageSectionsManager` thành card toàn chiều rộng bên dưới grid Page core/Media/SEO và phía trên action footer.
- Không đặt Section Manager ở cột phụ vì danh sách collection, picker và reorder cần không gian ngang.
- Trên desktop, danh sách Section dùng accordion toàn chiều rộng; trên mobile vẫn là một cột và actions được wrap thành nhiều dòng.
- Dùng MUI Theme tokens cho màu, border, spacing và typography. `backgroundColor` trong form là dữ liệu nội dung do biên tập viên cấu hình, không phải màu chrome của Dashboard.

### 3.3. Ranh giới Save

- Page core/Media/SEO tiếp tục dùng nút **Lưu thay đổi** của Page Edit.
- Section dùng mutation và nút lưu riêng trong từng ngữ cảnh: **Lưu Section**, **Lưu thứ tự**, **Lưu collection**.
- Không âm thầm gộp mutation Section vào chuỗi Page → Media → SEO vì lỗi Section không được làm mất hoặc gửi lại Page core.
- Nút lưu Page bị disable nếu dialog/picker Section đang submit, nhưng không bắt buộc lưu lại các Section đã hoàn tất mutation.

---

## 4. Cấu trúc file dự kiến

```text
apps/frontend/src/pages/pages/
├── edit.tsx
├── create.tsx
└── sections/
    ├── page-section-types.ts
    ├── page-section-api.ts
    ├── page-section-normalizers.ts
    ├── page-section-validation.ts
    ├── use-page-sections.ts
    ├── use-section-collection.ts
    ├── page-sections-manager.tsx
    ├── page-section-card.tsx
    ├── page-section-form-dialog.tsx
    ├── page-section-form-steps.tsx
    ├── section-media-field.tsx
    ├── section-collections.tsx
    ├── section-collection-panel.tsx
    ├── section-item-row.tsx
    ├── section-item-picker-dialog.tsx
    ├── media-picker-dialog.tsx
    ├── section-delete-dialog.tsx
    └── section-order-actions.tsx

apps/frontend/src/hooks/
└── use-debounce.ts             # tạo nếu dự án chưa có utility dùng chung
```

Phân trách nhiệm:

- `PageEdit`: Smart coordinator cấp Page; truyền `pageId`, trạng thái disabled và nhận `isDirty`/`isMutating` từ Section Manager.
- `usePageSections`: tải list, create/update/delete/toggle/reorder, rollback và invalidate state Section.
- `useSectionCollection`: lazy load, quản lý snapshot Item, reorder, sync và retry một collection.
- `page-section-api.ts`: nơi duy nhất ghép nested URL và gọi `customRequest`; component không gọi API trực tiếp.
- `page-section-normalizers.ts`: parse response không tin cậy, chuẩn hóa snake_case/camelCase và URL Media.
- Các Card/Dialog/Field/Row là dumb components, chỉ nhận typed props và callbacks.
- Không dùng `any`. Dữ liệu API thô dùng `unknown` và type guard/normalizer trước khi đưa vào UI.

---

## 5. TypeScript contracts

### 5.1. Section

```ts
export type PageSectionStatus = "ACTIVE" | "INACTIVE";
export type PageSectionItemType = "CATEGORY" | "POST" | "MEDIA";

export interface ISectionMediaSummary {
  id: number;
  fileName: string;
  originalUrl: string;
  thumbnailUrl?: string;
  mediumUrl?: string;
  mimeType: string;
}

export interface IPageSectionCollectionSummary {
  collection: string;
  itemType: PageSectionItemType;
  total: number;
}

export interface IPageSection {
  id: number;
  pageId: number;
  key: string;
  name: string;
  title: string | null;
  description: string | null;
  backgroundColor: string | null;
  backgroundMediaId: number | null;
  backgroundMedia: ISectionMediaSummary | null;
  featureMediaId: number | null;
  featureMedia: ISectionMediaSummary | null;
  sortOrder: number;
  status: PageSectionStatus;
  collections: IPageSectionCollectionSummary[];
  createdAt: string;
  updatedAt: string;
}

export interface IPageSectionFormValues {
  key: string;
  name: string;
  title: string;
  description: string;
  backgroundColor: string;
  backgroundMedia: ISectionMediaSummary | null;
  featureMedia: ISectionMediaSummary | null;
  status: PageSectionStatus;
}

export interface IUpsertPageSectionPayload {
  key: string;
  name: string;
  title: string | null;
  description: string | null;
  backgroundColor: string | null;
  backgroundMediaId: number | null;
  featureMediaId: number | null;
  sortOrder: number;
  status: PageSectionStatus;
}

export interface IReorderPageSectionsPayload {
  sections: Array<{ id: number; sortOrder: number }>;
}
```

- Form giữ object Media để preview nhưng payload chỉ gửi `backgroundMediaId` và `featureMediaId`.
- `sortOrder` khi create lấy từ `sections.length`; không cho người dùng nhập tay.
- Response list phải có `collections` summary để UI không gọi N+1 chỉ nhằm lấy số lượng.

### 5.2. Collection Item dạng discriminated union

```ts
export interface ICategorySectionData {
  id: number;
  name: string;
  slug: string;
  imageUrl?: string;
  status: "ACTIVE" | "HIDDEN";
}

export interface IPostSectionData {
  id: number;
  title: string;
  slug: string;
  typeCode: string;
  thumbnailUrl?: string;
}

export interface IMediaSectionData extends ISectionMediaSummary {
  status: "temporary" | "attached";
}

interface IPageSectionItemBase {
  id: number;
  sectionId: number;
  itemId: number;
  collection: string;
  sortOrder: number;
}

export type IPageSectionItem =
  | (IPageSectionItemBase & { itemType: "CATEGORY"; data: ICategorySectionData })
  | (IPageSectionItemBase & { itemType: "POST"; data: IPostSectionData })
  | (IPageSectionItemBase & { itemType: "MEDIA"; data: IMediaSectionData });

export interface ISyncSectionCollectionPayload {
  itemType: PageSectionItemType;
  items: Array<{ itemId: number; sortOrder: number }>;
}

export interface ISectionCollectionState {
  collection: string;
  itemType: PageSectionItemType;
  items: IPageSectionItem[];
  isLoading: boolean;
  isSaving: boolean;
  isDirty: boolean;
  error: string | null;
}
```

- Không dùng một interface `data` lỏng hoặc ép kiểu trong component; `itemType` quyết định component preview và field được đọc.
- Cùng một collection chỉ có một `itemType`. Normalizer phải từ chối response trộn loại thay vì render sai dữ liệu.

---

## 6. UI và tương tác

### 6.1. Header và trạng thái tổng

- Header card: **Các Section của trang**, mô tả ngắn, tổng số Section và nút **Thêm Section**.
- Khi không có Section, hiển thị Empty State với CTA tạo Section đầu tiên; không coi là lỗi.
- Khi GET lỗi, hiển thị `Alert` và nút **Thử lại**; Page form/Media/SEO vẫn hoạt động.
- Khi load lần đầu, render skeleton card/accordion; không render danh sách rỗng rồi nhấp nháy.

### 6.2. Section card/accordion

Summary hiển thị:

- Drag handle, `name`, technical `key`, `title` rút gọn.
- Status Chip `ACTIVE`/`INACTIVE` dùng palette semantics của Theme.
- Chip tổng Item theo collection, ví dụ `featured_posts · 6`.
- Actions: lên, xuống, sửa, bật/tắt và xóa; action có tooltip và `aria-label`.

Khi expand:

- Hiển thị description, Background/Feature Media preview và metadata không chỉnh trực tiếp.
- Hiển thị `SectionCollections` theo từng collection summary.
- Lazy load Item của collection khi panel được mở lần đầu; không request toàn bộ mọi collection ngay lúc mở Page Edit.

### 6.3. Tạo/sửa Section

- Dùng `Dialog` responsive (`maxWidth="md"`, full-screen ở breakpoint nhỏ) và MUI Stepper vì form có nhiều nhóm cấu hình.
- Bước 1 — **Thông tin chung**: key, name, title, description, status.
- Bước 2 — **Giao diện**: background color, Background Media, Feature Media.
- Bước 3 — **Xác nhận**: preview giá trị và nút lưu.
- Tạo mới gọi `POST`; sau success đóng dialog, append/refetch Section và tự expand Section mới để người dùng thêm collection tùy chọn.
- Chỉnh sửa gọi `PUT`; không chờ nút lưu Page core.
- Trong Edit, `key` mặc định read-only. Action **Đổi key** mở confirmation giải thích Webview có thể phụ thuộc key; chỉ sau xác nhận mới unlock field.
- Đóng dialog khi dirty phải có confirmation; submit lỗi giữ nguyên step/form và map field error nếu Backend cung cấp.

### 6.4. Background color và Media

- Background color là optional. Dùng MUI color input kết hợp text field để nhập giá trị hợp lệ; có action xóa về `null`.
- Dashboard không dùng màu Section làm màu chrome; chỉ hiển thị swatch preview có border từ Theme.
- Background/Feature dùng `MediaPickerDialog` ở chế độ single-select, chỉ hiện Media ảnh.
- Gỡ Media chỉ đặt field local về `null`; không gọi xóa file Media.
- Preview URL dùng `thumbnailUrl ?? mediumUrl ?? originalUrl`, qua helper resolve URL chung.

### 6.5. Reorder Section

- Pointer drag/drop chỉ đổi state local; không gọi API trong mỗi `dragover`/drop.
- Có nút **Di chuyển lên/xuống** cho keyboard/mobile; disabled ở phần tử đầu/cuối.
- Sau thay đổi, hiển thị action bar **Lưu thứ tự** và **Hoàn tác**.
- `Lưu thứ tự` gửi toàn bộ `{ id, sortOrder }` liên tục từ `0` qua một request.
- Nếu mutation lỗi, giữ thứ tự local, hiển thị lỗi và cho Retry; **Hoàn tác** khôi phục snapshot server gần nhất.
- Trong khi reorder đang dirty, create/delete/toggle bị disable để tránh snapshot chứa ID stale.

### 6.6. Xóa và bật/tắt

- Xóa bắt buộc có Confirmation Dialog, nêu rõ Item liên kết sẽ bị xóa nhưng dữ liệu nguồn không bị xóa.
- Delete success remove card/invalidate list và toast; lỗi giữ card, đóng progress và hiển thị message phù hợp.
- Toggle status có thể optimistic update, nhưng phải giữ bản copy trước mutation để rollback nếu `PUT` lỗi.
- Không cho double-submit hoặc thao tác lại cùng Section trong khi mutation đang chạy.

---

## 7. Collection Manager và Item Picker

### 7.1. Collection panel

- Mỗi collection hiển thị technical key, Item Type, tổng Item, trạng thái load và nút **Chọn dữ liệu**.
- Các collection thông dụng có label thân thiện:
  - `categories` → Danh mục.
  - `featured_posts` → Bài viết nổi bật.
  - `gallery` → Thư viện ảnh.
- Collection key khác vẫn hiển thị được bằng text humanized; không hard-code UI chỉ cho ba key mẫu.
- Thêm collection mới: chọn Item Type, nhập/chọn key gợi ý và bắt buộc chọn ít nhất một Item. Vì schema không có bảng định nghĩa collection riêng, collection rỗng không được lưu.
- Sync mảng rỗng có nghĩa gỡ toàn bộ Item và collection summary sẽ biến mất sau refetch.

### 7.2. Item Picker

- Một dialog dùng chung, nhận `itemType`, `collection`, `selectedIds`, `multiple`, callback confirm.
- List server-side gồm checkbox, thumbnail/avatar phù hợp, title/name/filename và metadata ngắn.
- Tìm kiếm dùng `useDebounce(value, 300)`; không dùng `setTimeout` tự viết trong component.
- Selection được lưu bằng `Set<number>` ở container và không mất khi đổi trang, filter hoặc search.
- Item đã selected vẫn checked khi xuất hiện lại; không thêm duplicate.
- Nút xác nhận hiển thị số lượng đã chọn và trả danh sách ID theo thứ tự hiện tại.
- `CATEGORY` gọi resource `categories`, filter `q`; có thể lọc status nếu Backend hỗ trợ.
- `POST` gọi resource `posts`, filter `title_like`; có filter tùy chọn `typeCode` để thu hẹp loại bài viết.
- `MEDIA` gọi resource `media`, filter search và `mimeType=image`; endpoint list Media là prerequisite Backend mới.

### 7.3. Đồng bộ và reorder Item

- Khi mở picker, dùng items đã lưu làm initial selection.
- Sau confirm, tạo local snapshot: giữ thứ tự Item cũ còn được chọn, append Item mới theo thứ tự chọn, chuẩn hóa `sortOrder` từ `0`.
- Người dùng có thể reorder Item bằng drag/drop và nút lên/xuống trước khi bấm **Lưu collection**.
- Chỉ `PUT` một snapshot đầy đủ; không gửi request trên từng add/remove/reorder.
- Sync lỗi giữ local state và mở action Retry. Sync success refetch collection + summary Section, reset dirty và toast.
- Không xóa Category/Post/Media nguồn khi gỡ khỏi collection.

---

## 8. Data flow và state management

### 8.1. Tải Section

1. `PageEdit` validate route ID và tải Page detail như hiện tại.
2. Sau khi Page detail success, enable `usePageSections(pageId)`.
3. Hook gọi một request list, normalize Media URL, sort Section theo `sortOrder` và lưu `serverOrderSnapshot`.
4. Collection Item chỉ được tải khi expand panel; cache theo key `{ pageId, sectionId, collection }`.
5. Error Section/collection chỉ ảnh hưởng khu vực tương ứng, không reset Page form, Page Media hoặc SEO.

### 8.2. Mutation Section

- `createSection`: POST, toast, invalidate/refetch list, expand Section mới.
- `updateSection`: PUT full payload, update/refetch card và reset dialog dirty state.
- `toggleSection`: PUT từ snapshot đầy đủ với status mới; rollback nếu lỗi.
- `deleteSection`: DELETE, remove/invalidate list và xóa cache collection của Section.
- `reorderSections`: PUT snapshot toàn danh sách; refetch để lấy source of truth Backend.
- `syncCollection`: PUT snapshot collection; refetch collection và list summary.

### 8.3. Dirty state và Page Edit

`usePageSections` expose tối thiểu:

```ts
export interface UsePageSectionsResult {
  sections: IPageSection[];
  isLoading: boolean;
  isMutating: boolean;
  isOrderDirty: boolean;
  hasDraftChanges: boolean;
  error: string | null;
  reload(): Promise<void>;
}
```

- `hasDraftChanges` là `isOrderDirty || sectionFormDirty || anyCollectionDirty`.
- `PageEdit` mở cancel confirmation khi `formState.isDirty || pageMedia.isDirty || pageSeo.isDirty || pageSections.hasDraftChanges`.
- Mutation Section đã success không còn dirty và không cần gửi lại khi Save Page.
- Khi route change trong lúc dialog/picker dirty hoặc reorder chưa lưu, dùng warning chung của Page Edit và confirmation nội bộ khi đóng dialog.
- Không lưu `File` object hoặc raw API DTO trong Section form/state.

### 8.4. Cache và invalidation Frontend

- Sau mutation Section, invalidate query list Page Section và đúng query collection; không invalidate toàn bộ ứng dụng.
- Sau Section/Item mutation success, invalidate Page detail/public-related key nếu Frontend đang cache response aggregate, nhưng không tự gọi Public API.
- Backend chịu trách nhiệm xóa Redis public Page cache; Frontend chỉ refresh dữ liệu Dashboard.
- Không request Item theo từng row. Mỗi collection dùng một request list phân trang hoặc snapshot endpoint.

---

## 9. API contract prerequisites

### 9.1. Page Section

```text
GET    /api/v1/admin/pages/:page_id/sections
POST   /api/v1/admin/pages/:page_id/sections
GET    /api/v1/admin/pages/:page_id/sections/:section_id
PUT    /api/v1/admin/pages/:page_id/sections/:section_id
DELETE /api/v1/admin/pages/:page_id/sections/:section_id
PUT    /api/v1/admin/pages/:page_id/section-order
```

- List response dùng wrapper `{ data, total }` và mỗi Section có `collections` summary.
- Create/Update trả `{ data: IPageSection }` để UI không tự dựng ID hoặc timestamps.
- API phải xác minh `section_id` thuộc `page_id`; Frontend không coi ID nested là đủ quyền.

### 9.2. Section Item

```text
GET /api/v1/admin/pages/:page_id/sections/:section_id/items?collection={collection}
PUT /api/v1/admin/pages/:page_id/sections/:section_id/items/:collection
```

GET response:

```json
{
  "data": [
    {
      "id": 901,
      "sectionId": 10,
      "itemType": "POST",
      "itemId": 501,
      "collection": "featured_posts",
      "sortOrder": 0,
      "data": {
        "id": 501,
        "title": "Bài viết nổi bật",
        "slug": "bai-viet-noi-bat",
        "typeCode": "NEWS"
      }
    }
  ],
  "total": 1
}
```

PUT body:

```json
{
  "itemType": "POST",
  "items": [
    { "itemId": 501, "sortOrder": 0 },
    { "itemId": 498, "sortOrder": 1 }
  ]
}
```

- Thống nhất camelCase cho Admin API. Nếu Backend trả snake_case, chỉ normalizer/API module được map; component luôn dùng camelCase.
- Backend trả dữ liệu hydrate tối thiểu đúng `itemType`; không buộc Frontend gọi detail cho từng Item.

### 9.3. Picker data sources

```text
GET /api/v1/admin/categories?current=1&pageSize=20&q={search}
GET /api/v1/admin/posts?current=1&pageSize=20&title_like={search}&typeCode={optional}
GET /api/v1/admin/media?current=1&pageSize=20&q={search}&mimeType=image
```

- Endpoint Media list/search là prerequisite bắt buộc cho Background/Feature/Media picker; không dùng Page Media relation API thay cho Media Library.
- Media list chỉ trả Media người dùng được phép dùng theo RBAC/ownership; UI không lọc quyền bằng dữ liệu client.
- Cả ba endpoint cần phân trang phẳng `{ data, total }` tương thích Refine.

---

## 10. Validation, error và feedback

### 10.1. Validation client

- `key`: required, tối đa 100, regex `^[a-z0-9]+(?:_[a-z0-9]+)*$`.
- `name`: required, tối đa 255.
- `title`: optional, tối đa 255.
- `description`: optional; giới hạn UI phải thống nhất với Backend plan trước khi code.
- `backgroundColor`: empty hoặc định dạng màu được Backend chấp nhận; không gửi chuỗi rỗng, map về `null`.
- Background/Feature Media chỉ nhận item có MIME type `image/*`.
- `collection`: required, tối đa 100, regex `^[a-z0-9]+(?:_[a-z0-9]+)*$`.
- `itemType`: union cố định; `itemId > 0`; collection mới phải có ít nhất một Item.
- `sortOrder`: không nhập tay, luôn normalize từ `0` trước khi request.

### 10.2. Mapping lỗi HTTP

- `400`: map field validation nếu có, nếu không hiển thị Alert/toast trong dialog/panel.
- `401`: để AuthProvider/customRequest refresh hoặc redirect; không tự đọc/ghi token trong module Section.
- `403`: thông báo không có quyền quản lý Page/Media, giữ state local để người dùng không mất thay đổi.
- `404`: Page lỗi toàn màn hình; Section/Item nguồn lỗi tại card/picker và có reload.
- `409`: key trùng, Item trùng hoặc race reorder; map vào key/collection, refetch source of truth trước khi retry khi cần.
- `500`: thông báo chung, giữ dialog/local order/selection để retry.

### 10.3. Feedback

- Mutation success hiển thị Snackbar góc phải và tự đóng theo cấu hình NotificationProvider hiện tại.
- Mutation lỗi không đóng dialog/picker và không xóa selection.
- Nút mutation có loading, khóa double-submit; không khóa toàn Page Edit nếu chỉ một collection đang retry.
- Partial state phải rõ ràng: Section metadata đã lưu nhưng collection lỗi là hai mutation độc lập, UI chỉ retry collection.

---

## 11. Hiệu năng, accessibility và bảo mật

- Search/autocomplete debounce `300ms` bằng `useDebounce`; cancel/stale query không được ghi đè kết quả mới.
- Picker dùng server pagination, page size mặc định 20; không tải toàn bộ Category/Post/Media vào client.
- Lazy load collection khi expand; list Section chỉ nhận summary để tránh hydrate toàn bộ dữ liệu khi không cần.
- Không tạo N+1 request theo Section card hoặc Item row.
- Reorder local và sync snapshot một lần; không gọi API theo mỗi thao tác kéo.
- Dialog có focus trap, title/description liên kết bằng ARIA; IconButton có `aria-label`; reorder có keyboard fallback.
- Preview ảnh có `alt` từ name/title/fileName và fallback khi URL lỗi.
- UI không gửi `userId`, `role`, `ownerId`, `pageId` trong body nếu ID đã nằm trong URL; Backend lấy identity từ Access Token.
- Không render HTML từ title/description/data nguồn bằng `dangerouslySetInnerHTML`.
- Không log payload/token hoặc chi tiết response nhạy cảm ra console.

---

## 12. Acceptance criteria

- Mở Page Edit hợp lệ tải một request Section list sau khi Page detail success; Page không tồn tại không gọi Section API.
- Empty/Error/Loading/Success của Section độc lập với Page core, Page Media và SEO.
- Tạo Section validate đúng, gửi Media ID thay vì Media object, thêm card đúng thứ tự và tự expand sau success.
- Edit Section giữ `key` read-only cho tới khi người dùng xác nhận đổi; lỗi `409` map vào field key.
- Toggle status rollback khi API lỗi; delete luôn có confirmation và không xóa dữ liệu Category/Post/Media nguồn.
- Reorder pointer/keyboard chỉ đổi local state; một lần lưu gửi snapshot liên tục từ `0`, retry không mất thứ tự local.
- Collection lazy load đúng loại và thứ tự; picker giữ selection qua pagination/search, không thêm duplicate.
- Search Category/Post/Media debounce `300ms`; kiểm tra Network không có request theo từng ký tự hoặc từng Item row.
- Đồng bộ collection gửi một Item Type và snapshot đầy đủ; sync mảng rỗng gỡ collection khỏi summary sau refetch.
- Background/Feature dùng single Media picker; collection Media dùng multi-select; chỉ ảnh hợp lệ được chọn.
- Unsaved warning bao gồm Section dialog, collection và reorder dirty; mutation Section đã thành công không bị gửi lại khi lưu Page.
- Toàn bộ module mới strict typed, không dùng `any`, component UI không gọi API và không hard-code màu Dashboard.
- `npm run build` tại `apps/frontend` thành công sau khi triển khai.

---

## 13. Checklist triển khai

- [ ] Hoàn thiện Backend Plan/API cho `page_sections`, `page_section_items`, aggregate collection summary và public cache invalidation.
- [ ] Bổ sung API Media Library list/search có phân trang, MIME filter và RBAC/ownership.
- [ ] Chốt camelCase response/payload và shape hydrated `data` cho ba Item Type.
- [ ] Tạo types, validation, normalizers và API module Section; không dùng `any`.
- [ ] Tạo `usePageSections`, nối list/create/update/delete/toggle/reorder và retry.
- [ ] Tạo Section Manager/Card/Empty/Error/Skeleton/Delete Dialog/Order Actions.
- [ ] Tạo React Hook Form + Stepper dialog cho create/edit và Media fields.
- [ ] Tạo `useSectionCollection`, Collection Panel, Item Row và snapshot sync.
- [ ] Tạo Item Picker dùng Refine list hooks, server pagination, persistent selection và `useDebounce(300)`.
- [ ] Tạo Media Picker single/multi dùng Media Library endpoint; không dùng Page Media link API.
- [ ] Tích hợp `PageSectionsManager` và aggregate dirty/mutating state vào `PageEdit`.
- [ ] Thêm notice Section disabled vào Page Create; không gọi API khi chưa có Page ID.
- [ ] Test unit normalizer/validation/reorder nếu test infrastructure có sẵn; không thêm test dependency ngoài phạm vi nếu dự án chưa cấu hình.
- [ ] Test thủ công Loading/Error/Empty/Success, 400/403/404/409/500, retry, double-submit và unsaved navigation.
- [ ] Kiểm tra Network không N+1, debounce đúng 300ms, một request snapshot cho reorder/sync.
- [ ] Chạy TypeScript build và regression Page Edit core/Media/SEO.
