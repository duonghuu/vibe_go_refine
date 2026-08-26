# FRONTEND PLAN: Quản lý SEO cho bài viết (Post SEO Meta)

## 1. Mục tiêu

- Bổ sung khu vực quản lý SEO vào luồng Tạo mới và Chỉnh sửa bài viết.
- Sử dụng seo_meta làm nguồn dữ liệu SEO duy nhất của Post.
- Giữ post_meta cho metadata mở rộng dạng key-value; Frontend không gửi các key SEO vào API Post Meta.
- Cho phép Admin/Staff nhập SEO tùy chỉnh và xem trước giá trị sau fallback.
- Tách mutation SEO khỏi mutation Post chính, phù hợp với kiến trúc module riêng của Backend.
- Tuân thủ Refine.js, React Hook Form, MUI và TypeScript strict; không dùng any.
- Chưa triển khai locale; mỗi Post hiện chỉ có một bộ SEO.

## 2. Phạm vi chức năng

### 2.1. Trong phạm vi

- Hiển thị SEO Meta card trong posts/create và posts/edit/:id.
- Quản lý:
  - SEO Title, SEO Description, Meta Keywords.
  - Canonical URL.
  - Open Graph Title/Description/Image/Type.
  - Twitter Title/Description/Image/Card.
  - Robots.
  - Schema JSON-LD.
- Tải SEO hiện tại khi mở trang Edit.
- Upsert SEO sau khi Post được tạo hoặc cập nhật.
- Hiển thị fallback preview cho Title, Description, Canonical URL và ảnh.
- Xử lý Loading, Success, Empty, Error và Retry riêng cho SEO.
- Cảnh báo rời trang khi SEO có thay đổi chưa lưu.
- Invalidate/refetch dữ liệu SEO sau mutation thành công.

### 2.2. Ngoài phạm vi

- Quản lý SEO cho Product, Category hoặc Page.
- Quản lý đa ngôn ngữ/locale.
- Chỉnh sửa dữ liệu post_meta.
- Xây dựng Media Library mới.
- Lưu media_id cho Open Graph/Twitter image; giai đoạn đầu dùng URL.
- Preview SERP/Open Graph hoàn chỉnh.
- Tự động sinh hoặc ghi fallback ngược vào database.

## 3. Vị trí và luồng màn hình

### 3.1. Post Create

Giữ layout hiện tại của PostCreate:

- Cột chính: Thông tin cơ bản gồm title, slug và content.
- Cột phụ: Phân loại, hình ảnh bài viết và card SEO Meta.
- Footer: Hủy và Lưu bài viết.

Thứ tự xử lý khi tạo mới:

1. Admin nhập thông tin Post và SEO.
2. Frontend tạo Post trước để nhận postId.
3. Frontend sync Media nếu có.
4. Frontend gọi PUT /api/v1/admin/seo-meta/post/{postId}.
5. Chỉ redirect khi các bước bắt buộc đã hoàn tất.
6. Nếu Post tạo thành công nhưng SEO thất bại, thông báo rõ và giữ khả năng Retry trong trang Edit.

SEO là dữ liệu tùy chọn. Nếu Admin không nhập trường SEO nào, Frontend vẫn cho phép tạo Post.

### 3.2. Post Edit

Khi mở posts/edit/:id:

1. Tải Post detail bằng useForm hiện có.
2. Tải SEO bằng GET /api/v1/admin/seo-meta/post/{postId}.
3. Hiển thị giá trị gốc trong input và giá trị resolved trong preview/helper text.
4. Cập nhật Post trước.
5. Sync Media theo luồng hiện có.
6. Upsert SEO.
7. Chỉ redirect sau khi các mutation được xử lý theo đúng thứ tự.
8. Nếu SEO lỗi, không làm mất dữ liệu Post đã lưu; hiển thị Retry SEO.

Có thể lưu Post mà không có SEO, nhưng phải hiển thị rõ trạng thái SEO chưa đồng bộ.

### 3.3. Không tạo route SEO độc lập

SEO của Post là card thuộc Post Create/Edit, không tạo menu/resource list riêng cho seo_meta ở giai đoạn đầu. API vẫn dùng route generic để tái sử dụng cho các entity khác sau này.

## 4. UI layout và UX

### 4.1. SEO Meta Card

Dùng MUI Card/CardHeader/CardContent cùng border, radius và spacing với các card hiện tại.

Chia thành các nhóm Accordion:

1. Cơ bản:
   - SEO Title.
   - SEO Description.
   - Meta Keywords.
   - Robots.
2. Canonical:
   - Canonical URL.
   - Helper text: để trống sẽ dùng URL từ slug.
3. Open Graph:
   - OG Title.
   - OG Description.
   - OG Image URL.
   - OG Type.
4. Twitter Card:
   - Twitter Title.
   - Twitter Description.
   - Twitter Image URL.
   - Twitter Card.
5. Structured Data:
   - Schema JSON-LD bằng multiline TextField hoặc editor đã có.
   - Hiển thị lỗi JSON trước khi submit.

### 4.2. Fallback preview

Hiển thị khu vực SEO Preview:

- Title preview: giá trị nhập hoặc resolved title.
- Description preview: giá trị nhập hoặc resolved description.
- Canonical preview: canonical nhập hoặc URL fallback.
- OG image preview: URL nhập hoặc thumbnail Post nếu Backend trả resolved URL.
- Gắn nhãn rõ Đã tùy chỉnh hoặc Fallback từ bài viết.

Preview không được gửi các trường resolved_* lên Backend.

### 4.3. Nút hành động

- Lưu thay đổi dùng palette.success.main theo Styleguide.
- Hủy dùng màu secondary/action disabled.
- Khi đang lưu Post/Media/SEO, khóa nút submit để tránh mutation trùng.
- Khi SEO lỗi độc lập, hiển thị nút Thử lại SEO.
- Không dùng màu HEX hard-code.

## 5. Component breakdown

### 5.1. Smart/container

- PostCreate: điều phối useForm, usePostMedia, usePostSeo và thứ tự mutation.
- PostEdit: điều phối Post query, Media query, usePostSeo, dirty state và redirect.
- usePostSeo(postId): quản lý fetch, form state, fallback preview, validate JSON và upsert.
- post-seo-api.ts: module gọi API SEO, không để component gọi fetch trực tiếp.
- post-seo-types.ts: toàn bộ TypeScript contract.
- seo-preview-utils.ts: hàm thuần để resolve preview, không có side effect.

### 5.2. Dumb/UI components

- PostSeoCard: layout tổng thể, chỉ nhận props/callback.
- SeoBasicFields: title, description, keywords, robots.
- SeoCanonicalFields: canonical URL.
- SeoOpenGraphFields: OG fields.
- SeoTwitterFields: Twitter fields.
- SeoSchemaField: textarea/editor và lỗi JSON.
- SeoResolvedPreview: hiển thị giá trị gốc/resolved.
- SeoSyncErrorState: lỗi API và nút Retry.

UI component không gọi API, không điều phối Refine mutation và không chứa business rule.

## 6. Data và state management

### 6.1. Refine integration

- Giữ useForm hiện tại cho Post.
- Không tạo Refine resource list cho seo_meta.
- Dùng useCustom/useCustomMutation hoặc API module có type rõ ràng để GET/PUT SEO.
- Dùng queryOptions.enabled: Boolean(postId) cho SEO GET ở Edit.
- Không dùng useOne với resource polymorphic nếu Data Provider không ánh xạ đúng route.
- Dùng cơ chế AuthProvider/dataProvider hiện tại, không viết refresh token riêng.

### 6.2. Hook state

Hook nên cung cấp state tương đương:

    interface PostSeoState {
      data: IPostSeoResponse | null;
      isLoading: boolean;
      isSaving: boolean;
      isError: boolean;
      error: HttpError | null;
      isDirty: boolean;
    }

Và các hành động:

    load(): Promise<void>
    save(values: IPostSeoFormValues): Promise<void>
    retry(): Promise<void>
    reset(): void

Tên type thực tế có thể điều chỉnh theo convention hiện có nhưng không được dùng any.

### 6.3. Dirty state

Dirty state của trang phải tổng hợp:

- Post form isDirty.
- Post Media isDirty.
- Post SEO isDirty.

Khi một phần dirty:

- warnWhenUnsavedChanges phải bật.
- Nút Hủy mở Confirmation Dialog.
- Không reset SEO local state khi Post query refetch ngoài chủ ý của người dùng.

## 7. TypeScript contracts

    type SeoRobots =
      | "index,follow"
      | "noindex,follow"
      | "noindex,nofollow";

    type SeoOgType = "article" | "website" | "product";
    type SeoTwitterCard = "summary" | "summary_large_image";

    type JsonPrimitive = string | number | boolean | null;
    type JsonValue =
      | JsonPrimitive
      | JsonValue[]
      | { [key: string]: JsonValue };

    interface IPostSeoFormValues {
      metaTitle: string;
      metaDescription: string;
      metaKeywords: string;
      canonicalUrl: string;
      ogTitle: string;
      ogDescription: string;
      ogImage: string;
      ogType: SeoOgType | "";
      twitterTitle: string;
      twitterDescription: string;
      twitterImage: string;
      twitterCard: SeoTwitterCard | "";
      robots: SeoRobots;
      schemaJsonText: string;
    }

    interface IPostSeoResponse {
      entityType: "post";
      entityId: number;
      resolvedTitle: string;
      resolvedDescription: string;
      resolvedCanonicalUrl: string;
      resolvedOgTitle: string;
      resolvedOgDescription: string;
      resolvedOgImage: string | null;
      resolvedTwitterTitle: string;
      resolvedTwitterDescription: string;
      resolvedTwitterImage: string | null;
      schemaJson: JsonValue | null;
    }

Form dùng schemaJsonText là chuỗi để TextField/editor có thể chỉnh sửa. API adapter phải parse chuỗi này thành schema_json trước khi PUT và serialize schema_json thành chuỗi có format khi đổ dữ liệu vào form.

Frontend API adapter phải map rõ camelCase sang snake_case nếu dataProvider không tự xử lý. Không gửi entityType, entityId hoặc các trường resolved_* trong body.

## 8. Validation

- metaTitle: optional, tối đa 255 ký tự.
- metaDescription: optional, tối đa 500 ký tự; hiển thị counter và gợi ý 150–160 ký tự.
- metaKeywords: optional, tối đa 500 ký tự.
- canonicalUrl: optional; nếu nhập phải là absolute URL hợp lệ.
- ogTitle: optional, tối đa 255 ký tự.
- ogDescription: optional, tối đa 500 ký tự.
- ogImage: optional; nếu nhập phải là URL hợp lệ.
- twitterTitle: optional, tối đa 255 ký tự.
- twitterDescription: optional, tối đa 500 ký tự.
- twitterImage: optional; nếu nhập phải là URL hợp lệ.
- robots: chỉ cho phép options Backend hỗ trợ.
- schemaJson: optional; phải parse được JSON và hiển thị lỗi theo field.

Frontend validation chỉ cải thiện UX; Backend vẫn là nơi validate cuối cùng.

Chuẩn hóa input:

- Trim các field text trước khi gửi.
- Chuyển chuỗi rỗng thành null để Backend áp dụng fallback.
- Không tự sinh canonical URL và không ghi fallback vào form như giá trị tùy chỉnh.
- Không debounce input SEO vì đây là form local.

## 9. API contract prerequisite

Frontend phụ thuộc:

    GET /api/v1/admin/seo-meta/post/{postId}
    PUT /api/v1/admin/seo-meta/post/{postId}

GET cần trả:

- Giá trị gốc đã lưu.
- Giá trị resolved để hiển thị preview.
- entityType = post.
- entityId.
- 404 nếu chưa có SEO record; Frontend xử lý như Empty state.

PUT cần:

- Upsert theo post và postId.
- Trả record SEO đã lưu.
- Trả lỗi field-level cho URL, độ dài, robots và schema JSON.
- Trả 401 để AuthProvider xử lý refresh/redirect.
- Trả 403 nếu role không có quyền.

Cập nhật .docs/api-endpoints.yaml khi Backend contract được triển khai.

## 10. Loading, Error, Empty và Success

- Loading Post: SEO card hiển thị skeleton riêng.
- Loading SEO: field disabled hoặc skeleton tới khi GET hoàn tất.
- Empty SEO/404: khởi tạo form rỗng và hiển thị resolved fallback từ Post.
- GET error: Alert trong SEO card, có Retry; giữ Post form hoạt động.
- Post mutation error: giữ toàn bộ dữ liệu Post/Media/SEO local.
- SEO mutation error: không redirect; giữ dữ liệu local, hiển thị field errors và Retry.
- Success: Snackbar Đã lưu SEO cho bài viết; invalidate query SEO và resource Post.
- Partial success: thông báo Bài viết đã lưu, SEO chưa đồng bộ.
- 401: để dataProvider/AuthProvider xử lý thống nhất.

## 11. Trình tự mutation

### 11.1. Create

    Validate Post
      → Create Post
      → Sync Post Media
      → Upsert Post SEO
      → Invalidate posts/SEO
      → Redirect

Nếu Media không có thay đổi thì bỏ qua bước Sync Media.

### 11.2. Edit

    Validate Post
      → Update Post
      → Sync Post Media nếu dirty
      → Upsert Post SEO nếu dirty hoặc cần tạo record
      → Invalidate posts/SEO
      → Redirect

- Không gửi SEO nếu card chưa dirty và đã có dữ liệu hợp lệ.
- Nếu người dùng xóa toàn bộ SEO field, vẫn gửi payload null để Backend áp dụng fallback.
- Không chạy song song mutation nếu thứ tự lỗi có thể tạo trạng thái khó khôi phục.

## 12. Checklist triển khai

- [ ] Tạo apps/frontend/src/pages/posts/post-seo-types.ts.
- [ ] Tạo apps/frontend/src/pages/posts/post-seo-api.ts.
- [ ] Tạo apps/frontend/src/pages/posts/use-post-seo.ts.
- [ ] Tạo dumb components cho SEO card và preview.
- [ ] Tích hợp SEO card vào posts/create.tsx.
- [ ] Tích hợp SEO query/card/mutation vào posts/edit.tsx.
- [ ] Tổng hợp dirty state giữa Post, Media và SEO.
- [ ] Xử lý SEO GET 404 như Empty state.
- [ ] Xử lý Loading, Error, Retry, Partial Success và 401.
- [ ] Map camelCase/snake_case tại một API adapter duy nhất.
- [ ] Không gửi post_meta, resolved_*, entityType hoặc entityId trong body PUT.
- [ ] Không sử dụng any trong type, hook hoặc component.
- [ ] Kiểm thử fallback preview, URL validation, JSON validation, reset field về null và robots options.
- [ ] Kiểm thử Create/Edit, refresh trang, dirty warning, double-submit và Post save thành công nhưng SEO thất bại.
- [ ] Cập nhật API contract/OpenAPI trước khi tích hợp thật.
