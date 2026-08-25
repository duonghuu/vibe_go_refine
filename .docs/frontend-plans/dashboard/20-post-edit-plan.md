# FRONTEND PLAN: Chỉnh sửa bài viết (Admin Edit Post)

## 1. Mục tiêu

- Xây dựng màn hình Chỉnh sửa bài viết dựa trên `15-post-create-idea.md` và UI `PostCreate` hiện có.
- Cho phép Admin/Editor tải dữ liệu bài viết theo `id`, chỉnh sửa Tiêu đề, Slug, Nội dung và Danh mục, sau đó cập nhật qua Refine.js.
- Giữ nguyên `typeCode` của bài viết; không cho người dùng chuyển bài viết sang loại bài viết khác trong màn hình Edit.
- Tái sử dụng layout, card, autocomplete danh mục, notification và cơ chế cảnh báo thay đổi chưa lưu của luồng Create.

## 2. Phạm vi chức năng

### 2.1. Chức năng chính

- **Thông tin cơ bản:** Tiêu đề, Slug, Nội dung.
- **Phân loại:** Danh mục bài viết, chỉ hiển thị các danh mục thuộc `typeCode` của bài viết.
- **Loại bài viết:** Hiển thị read-only dưới dạng mã/tên loại; không bind thành input có thể sửa.
- **Điều hướng:** Hủy và quay về danh sách đúng `type_code`; Lưu và quay về danh sách sau khi cập nhật thành công.
- **Trạng thái dữ liệu:** Loading khi tải chi tiết, Error khi API lỗi/không tìm thấy, Success khi form đã được đổ dữ liệu.
- **Unsaved changes:** Cảnh báo khi rời trang trong lúc form có thay đổi chưa lưu.

### 2.2. Ngoài phạm vi màn hình Edit cơ bản

- Post Media và Post Meta đã có API riêng. Không gộp các mutation này vào payload cập nhật Post chính.
- Nếu UI yêu cầu quản lý thumbnail/attachment hoặc SEO meta trong cùng màn hình, triển khai thành các card độc lập và gọi lần lượt:
  - `GET /api/v1/admin/posts/{post_id}/media` và `PUT /api/v1/admin/posts/{post_id}/media/{collection}`.
  - `GET /api/v1/admin/posts/{postId}/meta` và `PUT /api/v1/admin/posts/{postId}/meta`.
- Các card mở rộng chỉ được submit sau khi Post chính cập nhật thành công, có xử lý lỗi riêng và không làm mất dữ liệu form chính.

## 3. Luồng màn hình

### 3.1. Header

- Tiêu đề: `Chỉnh sửa bài viết`.
- Breadcrumb: `Bảng điều khiển > Bài viết > Chỉnh sửa`.
- Có thể hiển thị `typeCode` và tên loại bài viết ở vùng header hoặc trong card thông tin; giá trị này là read-only.

### 3.2. Form layout

Giữ layout 2 cột và style của `PostCreate`:

- **Cột chính (`md={8}`):** Card `Thông tin cơ bản`.
  - `Tiêu đề (*)`: TextField, lấy giá trị cũ từ API.
  - `Đường dẫn (Slug) (*)`: TextField, giữ nguyên slug cũ khi người dùng đổi tiêu đề; chỉ thay đổi nếu người dùng chủ động sửa.
  - `Nội dung (*)`: dùng component editor hiện có của Create. Nếu editor Rich Text chưa được chuẩn hóa, trước mắt dùng cùng component/textarea hiện tại để bảo toàn HTML/text trả về từ API.
- **Cột phụ (`md={4}`):** Card `Phân loại`.
  - `Loại bài viết`: Text/disabled field, hiển thị `typeCode` hoặc tên Post Type.
  - `Danh mục bài viết`: Autocomplete phân cấp, fetch theo `typeCode`, hiển thị category hiện tại sau khi dữ liệu và options đã sẵn sàng.
- **Footer actions:**
  - `Hủy`: nếu form dirty thì hiển thị Confirmation Dialog; sau đó về `/posts?type_code={typeCode}`.
  - `Lưu thay đổi`: disabled trong lúc mutation hoặc khi đang load; gọi `PUT`, hiện Snackbar thành công và quay về list.

## 4. Component breakdown

- `PostEditPage`: Smart/container component, đọc `id`, tải dữ liệu, khởi tạo `useForm`, điều phối submit và redirect.
- `PostBasicInfoCard`: Dumb component nhận `register`, `control`, `errors`, `watch` và render title/slug/content.
- `PostClassificationCard`: Dumb component render `typeCode` read-only và category Autocomplete.
- `PostCategoryAutocomplete`: Tái sử dụng logic phân cấp từ Create; nhận `typeCode`, `value`, `onChange`, `loading`, `error`.
- `PostEditFooterActions`: Nhận `isLoading`, `onCancel`, `onSubmit` và render các nút hành động.
- Có thể trích xuất `PostForm` dùng chung cho Create/Edit sau khi hai màn hình có cùng contract; logic auto-slug phải bật riêng cho Create và tắt mặc định ở Edit.

## 5. Data & state management

### 5.1. Refine resource và route

- Bổ sung resource `posts.edit: "/posts/edit/:id"` trong `apps/frontend/src/App.tsx`.
- Bổ sung route `<Route path="edit/:id" element={<PostEdit />} />` dưới `/posts`.
- Đổi import page từ `{ PostCreate, PostList }` thành có thêm `PostEdit`.
- Giữ resource name là `posts` để Data Provider tự ánh xạ sang `/api/v1/admin/posts`.

### 5.2. Hook và fetch dữ liệu

- Dùng `useForm<IPostResponse, HttpError, IUpdatePostRequest>` từ `@refinedev/react-hook-form` với:
  - `action: "edit"`.
  - `resource: "posts"`.
  - `redirect: false` để kiểm soát redirect về list theo `typeCode`.
  - `warnWhenUnsavedChanges: true`.
- Refine tự gọi `GET /api/v1/admin/posts/{id}` để lấy bản ghi chi tiết; không dùng dữ liệu thiếu Content từ list làm default value.
- Đọc `id` bằng `useParams`/`useParsed`, ép kiểu và chặn render form nếu id không hợp lệ.
- Lấy `typeCode` ưu tiên từ response chi tiết. Query `type_code` chỉ dùng làm fallback cho redirect hoặc trạng thái loading; không được dùng query client để ghi đè `typeCode` của record.
- Chỉ bật query category sau khi `typeCode` đã có (`enabled: Boolean(typeCode)`).

### 5.3. TypeScript contracts

```ts
interface IPostResponse {
  id: number;
  typeCode: string;
  title: string;
  slug: string;
  content: string;
  authorId: number;
  categoryId?: number | null;
  createdAt: string;
}

interface IUpdatePostRequest {
  title: string;
  slug: string;
  content: string;
  categoryId?: number | null;
}
```

- Không gửi `id`, `authorId` hoặc `typeCode` do người dùng nhập trong update payload.
- `typeCode` được dùng làm context/đối chiếu; backend quyết định giá trị hợp lệ theo record hiện tại.
- Không sử dụng `any`; API error phải dùng `HttpError` hoặc type cụ thể từ data provider.

## 6. API contract prerequisite

Frontend Edit phụ thuộc các endpoint sau tại nhóm `/api/v1/admin/posts`:

- `GET /api/v1/admin/posts/{id}`: trả về `data` chứa đầy đủ `id`, `typeCode`, `title`, `slug`, `content`, `authorId`, `categoryId`, `createdAt`.
- `PUT /api/v1/admin/posts/{id}`: nhận `title`, `slug`, `content`, `categoryId` và trả về record đã cập nhật trong `data`.
- API cần trả lỗi có thể hiển thị được cho các trường hợp `404`, slug trùng `409`, dữ liệu không hợp lệ `400`, hết phiên `401`.
- Sau mutation thành công backend cần invalidate cache danh sách Post. Frontend chỉ cần `invalidate`/refetch resource `posts` thông qua Refine sau khi redirect.
- Cập nhật `.docs/api-endpoints.yaml` và route/controller backend trước khi tích hợp UI nếu hai endpoint trên chưa được tạo.

## 7. Validation và hành vi Edit

- `title`: bắt buộc, tối đa 255 ký tự.
- `slug`: bắt buộc, tối đa 255 ký tự, regex `^[a-z0-9-]+$`.
- `content`: bắt buộc; cần kiểm tra nội dung rỗng sau khi loại bỏ whitespace/HTML rỗng nếu dùng Rich Text Editor.
- `categoryId`: nullable; chỉ nhận option thuộc `typeCode` hiện tại.
- `typeCode`: bắt buộc từ response; nếu thiếu hoặc không hợp lệ thì hiển thị lỗi và quay về list.
- Không auto-generate slug theo `title` trong Edit để tránh thay đổi URL/SEO ngoài ý muốn. Nếu cần hỗ trợ, dùng hành động rõ ràng `Tạo lại Slug`, không chạy tự động khi title thay đổi.
- Khi API trả lỗi slug trùng, gắn lỗi vào field `slug` và hiển thị Snackbar mô tả lỗi backend.

## 8. Loading, Error, Empty và mutation feedback

- **Loading:** hiển thị loading state của `<Edit>`/form; không render input với default value rỗng trước khi GET hoàn tất.
- **404/GET error:** hiển thị thông báo `Không tìm thấy bài viết` hoặc lỗi backend, có nút quay về danh sách theo `typeCode` fallback.
- **Category error:** vẫn cho phép hiển thị form chính nhưng báo lỗi tại card Phân loại; không âm thầm reset `categoryId` cũ.
- **Submit loading:** khóa các nút và ngăn submit lặp.
- **Success:** Snackbar `Cập nhật bài viết thành công`, sau đó điều hướng `/posts?type_code={resolvedTypeCode}`.
- **Failure:** giữ nguyên dữ liệu người dùng nhập, hiển thị error notification và lỗi field nếu có.

## 9. Checklist triển khai

- [ ] Tạo `apps/frontend/src/pages/posts/edit.tsx` với type rõ ràng, không dùng `any`.
- [ ] Bổ sung resource edit và route edit cho `posts` trong `App.tsx`.
- [ ] Bổ sung/kiểm tra API `GET /api/v1/admin/posts/{id}` và `PUT /api/v1/admin/posts/{id}`.
- [ ] Dùng `useForm(action: "edit")`, load detail và map default values đúng kiểu `categoryId` nullable.
- [ ] Tái sử dụng layout/card/category hierarchy của `PostCreate`.
- [ ] Hiển thị `typeCode` read-only và khóa việc đổi loại bài viết.
- [ ] Tắt auto-slug trong Edit; validate slug khi người dùng chỉnh thủ công.
- [ ] Tích hợp `warnWhenUnsavedChanges` và Confirmation Dialog cho nút Hủy.
- [ ] Xử lý đầy đủ Loading, Success, Error, 404 và lỗi slug trùng.
- [ ] Kiểm thử redirect giữ đúng `type_code`, refresh trang Edit, submit hai lần và category không thuộc type.
- [ ] Nếu bật Media/Meta card, kiểm thử các mutation liên quan độc lập sau khi Post chính cập nhật thành công.

