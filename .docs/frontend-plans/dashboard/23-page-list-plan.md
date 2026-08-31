# FRONTEND PLAN: Quản trị danh sách Trang (Page List)

## 1. Mục tiêu

- Xây dựng màn hình CMS danh sách Page dựa trên `.docs/ideas/dashboard/21-page-list-idea.md`.
- Page là resource độc lập (`pages`), không truyền `typeCode` và không tái sử dụng resource `posts`.
- Cho phép Admin/Staff xem danh sách, tìm kiếm, lọc trạng thái, phân trang, đi tới tạo/sửa và xóa mềm Page.
- Tuân thủ React + Refine + MUI, TypeScript strict, UI component không gọi API và không dùng `any` trong code mới.

## 2. Phạm vi

### Trong phạm vi

- Thêm Page List resource và route `/pages` trong Refine.
- DataGrid server-side với các cột ID, tiêu đề, slug, trạng thái, cập nhật lúc và thao tác.
- Tìm kiếm title/slug bằng debounce 300ms; lọc `DRAFT`/`PUBLISHED`; phân trang và sắp xếp đồng bộ URL.
- Loading, error, empty, xác nhận xóa và thông báo kết quả.

### Ngoài phạm vi

- Form tạo/sửa Page, Page Media Card, SEO Card và public Webview route.
- Quản lý post type, category, hierarchy, menu hoặc bulk actions.
- Thay đổi generic data provider hiện tại ngoài việc dùng resource `pages` qua contract REST đã có.

## 3. Cấu trúc và điều hướng

### Files dự kiến

```text
apps/frontend/src/pages/pages/
├── index.ts
├── list.tsx
├── page-list-toolbar.tsx
├── page-status-chip.tsx
└── page-types.ts
```

### Đăng ký Refine

- Bổ sung resource `pages` trong `apps/frontend/src/App.tsx`:
  - `list: "/pages"`
  - `create: "/pages/create"`
  - `edit: "/pages/edit/:id"`
  - `meta.canDelete: true`
- Trong hạng mục này chỉ render route `/pages`; các route Create/Edit là dependency của các hạng mục frontend kế tiếp. Nút tạo/sửa điều hướng đến các URL này mà không tự triển khai form ở đây.

## 4. UI và component

### 4.1. Smart component: `PageList`

- Dùng `useDataGrid<IPageListItem, HttpError>` cho resource `pages`, `syncWithLocation: true` và phân trang/sắp xếp phía server.
- Dùng state cục bộ `searchInput`; debounce bằng `lodash/debounce` với delay 300ms rồi gọi `setFilters`.
- Khi search trống, gỡ filter; không gọi API theo từng ký tự.
- Dùng `useDelete` hoặc `DeleteButton` của Refine cho soft delete và invalidate danh sách sau thành công.
- Hiển thị title, `List` container và DataGrid; điều phối Loading/Error/Empty nhưng không đặt raw fetch vào component.

### 4.2. Dumb components

- `PageListToolbar`: nhận `searchValue`, `status`, callbacks, trạng thái loading; gồm search input, Status Select và nút Create.
- `PageStatusChip`: chỉ nhận `DRAFT | PUBLISHED`; dùng màu semantic từ MUI Theme (`warning` cho draft, `success` cho published), không hard-code HEX.
- Không tách DataGrid cell renderer thành component nếu không cần tái sử dụng; mọi renderer vẫn phải nhận type rõ ràng.

### 4.3. DataGrid

| Cột | Cấu hình |
| --- | --- |
| `id` | Rộng cố định, sortable. |
| `title` | Cột chính, nhấn để điều hướng Edit. |
| `slug` | Ellipsis/tooltip khi dài, sortable. |
| `status` | `PageStatusChip`, lọc server-side. |
| `updatedAt` | Định dạng `vi-VN`, sortable. |
| actions | EditButton và DeleteButton có tooltip/xác nhận. |

- Dùng density `compact`, tắt column menu nếu không cần và giữ chiều cao hàng phù hợp text.
- Loading dùng state của DataGrid; Error dùng Alert kèm nút thử lại; Empty dùng Empty State và nút tạo Page.

## 5. Data contracts và Refine mapping

```ts
export type PageStatus = "DRAFT" | "PUBLISHED";

export interface IPageListItem {
  id: number;
  title: string;
  slug: string;
  status: PageStatus;
  authorId: number;
  createdAt: string;
  updatedAt: string;
}

export interface IPageListResponse {
  data: IPageListItem[];
  total: number;
}
```

Backend contract cần tương thích với data provider hiện tại:

```text
GET /api/v1/admin/pages
  ?current=<number>
  &pageSize=<number>
  &title_like=<string>
  &status=DRAFT|PUBLISHED
  &sortBy=id|title|slug|updated_at
  &order=asc|desc

DELETE /api/v1/admin/pages/:id
```

- `title_like` biểu thị tìm theo title hoặc slug theo đặc tả IDEA; frontend không lọc client-side.
- `status` không chọn nghĩa là không truyền filter.
- API list phải trả `{ data, total }`; API delete trả success body mà Refine có thể unwrap qua data provider hiện có.

## 6. Luồng tương tác

1. Mở `/pages`: Refine gọi danh sách trang đầu tiên, đồng thời phản ánh query params trên URL.
2. Nhập ô search: cập nhật UI ngay; sau 300ms, filter `title_like` được set và DataGrid quay về trang đầu.
3. Chọn trạng thái: cập nhật filter `status` tức thời, reset trang đầu và đồng bộ URL.
4. Chuyển trang/sắp xếp: DataGrid gửi paging/sort xuống API thông qua Refine.
5. Nhấn tạo/sửa: điều hướng lần lượt tới `/pages/create` và `/pages/edit/:id`.
6. Nhấn xóa: mở confirmation; sau success hiển thị Snackbar và refresh grid; lỗi giữ dữ liệu đang xem và hiển thị thông báo backend.

## 7. Acceptance test

- Hiển thị đúng dữ liệu API, `total`, pagination và các state loading/error/empty.
- Search chỉ phát request sau tối thiểu 300ms kể từ lần gõ cuối; xóa search trả về danh sách không filter.
- Filter `DRAFT` và `PUBLISHED` gửi đúng query, có thể kết hợp search, và giữ lại khi reload URL.
- Sort các cột được hỗ trợ gửi `sortBy`/`order` chính xác; không phát sinh sort client-side.
- Delete có confirmation, success refresh bảng, API lỗi hiển thị Snackbar và không điều hướng sai.
- TypeScript compile không có `any` trong module Page mới; giao diện dùng MUI Theme tokens.

## 8. Phụ thuộc và giả định

- Backend Page List/Delete, AuthMiddleware và response `{ data, total }` đã sẵn sàng theo các IDEA/Page plan.
- Quyền truy cập đã được Backend kiểm soát bằng JWT/RBAC; frontend chỉ dựa vào AuthProvider hiện có để xử lý `401`.
- Route Create/Edit sẽ được hoàn thiện bởi frontend plans riêng; chưa thêm shortcut/menu đặc biệt ngoài resource Page.
