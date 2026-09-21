# WEBVIEW PLAN: Tích hợp IntroSection với Page Section API

## 1. Mục tiêu

Thay phần heading mock của `IntroSection` trên route `/` bằng dữ liệu Section được quản trị tại Admin Page `/pages/edit/6`, tương ứng Page đã publish có slug `trang-chu`.

Kết quả cần đạt:

- Webview lấy nội dung public theo slug `trang-chu`, không phụ thuộc ID database `6` ở runtime.
- Section được chọn bằng technical key ổn định `home_creative`, không suy luận từ tên hoặc nội dung hiển thị.
- `IntroSection` tiếp tục là presentation component, không tự gọi API.
- Route `/` vẫn là Next.js Server Component; không thêm `"use client"`, Axios hoặc TanStack Query.
- Có đủ trạng thái loading, success, empty và error ở phạm vi riêng của Intro.
- Loại bỏ heading mock của Intro khỏi runtime fixture sau khi API hoạt động.

## 2. Hiện trạng đã xác minh

### 2.1. Webview

- `apps/webview/app/page.tsx` đang là async Server Component, dùng `dynamic = "force-dynamic"`.
- Hero đã có data layer gọi API bằng `fetch(..., { cache: "no-store" })`; Intro vẫn nhận toàn bộ `content.intro` từ `apps/webview/data/webview-home.json`.
- `IntroSection` cần hai field heading:
  - `heading.eyebrow`
  - `heading.title`
- Ba benefit card được lấy động từ `section.collections.benefits`; mỗi item cần cung cấp title/name, description và image URL trong public item data.
- `app/loading.tsx` chỉ tạo skeleton cấp route. Theo tài liệu Next.js 16 trong repo, loading riêng cho Intro nên dùng `<Suspense>` gần async component để không chặn toàn bộ trang.

### 2.2. Dữ liệu Admin hiện tại

Dữ liệu read-only trong database tại thời điểm lập plan:

| Trường | Giá trị |
| --- | --- |
| Page ID | `6` |
| Page slug | `trang-chu` |
| Page status | `PUBLISHED` |
| Section ID | `3` |
| Section key | `home_creative` |
| Section status | `ACTIVE` |
| Section sort order | `0` |
| Section title | `We are creative & expert people` |
| Section description | `We work with business & provide solution to client with their business problem` |
| Section collections | Chưa có item |

Mapping cho Intro:

| Public Section | Intro UI |
| --- | --- |
| `section.title` | `content.heading.eyebrow` |
| `section.description` | `content.heading.title` |
| `section.collections.benefits[*]` | `content.benefits` |

Không dùng `section.name`: đây là tên quản trị và public contract trong backend plan đã chủ đích không trả field này.

### 2.3. Khoảng trống backend bắt buộc

Backend hiện chỉ đăng ký public route:

```text
GET /api/v1/public/image-contents
```

Các Page Section endpoint hiện có đều nằm dưới `/api/v1/admin`, yêu cầu JWT và role `ADMIN`/`STAFF`; Webview public tuyệt đối không gọi các endpoint này.

Plan backend đã định nghĩa nhưng code runtime chưa triển khai:

```text
GET /api/v1/pages/:slug
```

Theo:

- `.docs/backend-plans/dashboard/26-page-public-api-plan.md`
- Mục “Mở rộng Public Page read model” trong `.docs/backend-plans/dashboard/27-page-section-management-plan.md`

Vì vậy Public Page API có `sections` là prerequisite bắt buộc trước khi bật tích hợp Intro ở Webview.

## 3. Phạm vi

### Trong phạm vi

- Đọc Page public bằng slug `trang-chu`.
- Tìm Section `ACTIVE` có key `home_creative` từ response public.
- Map `title` và `description` thành heading của Intro.
- Đọc collection `benefits` và loop các item public thành benefit card; không fallback sang fixture local.
- Thêm type DTO, runtime validation, normalizer/data loader và UI state cho Intro.
- Dùng `<Suspense>` để stream skeleton riêng của Intro.
- Xóa `intro.heading` khỏi fixture sau khi API contract đã chạy thật.

### Ngoài phạm vi

- Không gọi API admin `/admin/pages/6/sections` từ Webview.
- Không hard-code Page ID `6` trong Webview.
- Không đưa Access Token, cookie admin hoặc role vào public request.
- Không thay đổi giao diện, spacing, màu sắc, grid hoặc ảnh của ba benefit card.
- Không suy diễn item thiếu title hoặc image; item không đủ dữ liệu bị bỏ qua và collection rỗng chuyển sang empty state.
- Không nhét JSON vào `section.description` để vượt qua giới hạn schema.
- Không tích hợp các section Home còn lại trong cùng task.

## 4. API contract mục tiêu

### 4.1. Endpoint

```http
GET /api/v1/pages/trang-chu
Accept: application/json
```

Webview tạo URL bằng:

```ts
getWebviewApiEndpoint("/pages/trang-chu")
```

Không nối chuỗi trực tiếp từ biến môi trường và không dùng `WEBVIEW_PUBLIC_API_URL` làm API base; biến đó chỉ dùng để resolve public media URL.

### 4.2. Response tối thiểu Intro cần dùng

```json
{
  "data": {
    "id": 6,
    "title": "Trang chủ",
    "slug": "trang-chu",
    "content": "...",
    "updatedAt": "2026-09-18T00:00:00Z",
    "thumbnail": null,
    "gallery": [],
    "seo": {},
    "sections": [
      {
        "id": 3,
        "key": "home_creative",
        "title": "We are creative & expert people",
        "description": "We work with business & provide solution to client with their business problem",
        "backgroundColor": null,
        "backgroundMedia": null,
        "featureMedia": null,
        "sortOrder": 0,
        "collections": {}
      }
    ]
  }
}
```

Backend phải bảo đảm:

- Chỉ Page `PUBLISHED`, chưa soft-delete được public.
- Chỉ Section `ACTIVE`, chưa soft-delete được trả về.
- Section order cố định theo `sort_order ASC, id ASC`.
- `collections` là `{}` khi không có item, không trả `null`.
- Không trả `name`, `status`, timestamps admin, owner hoặc dữ liệu xác thực.
- Cache public Page gồm cả Sections bằng key `public:pages:slug:trang-chu`; mutation Section/Item phải invalidate key sau transaction commit.

### 4.3. Error semantics

| Điều kiện | HTTP/API | Webview state |
| --- | --- | --- |
| Page không tồn tại, DRAFT hoặc soft-deleted | `404 PAGE_NOT_FOUND` | `empty` |
| Response 200 nhưng không có `home_creative` | `200` | `empty` |
| Section tồn tại nhưng title hoặc description rỗng sau trim | `200` | `empty` |
| Network timeout, DNS, API 5xx | lỗi request | `error` |
| JSON sai contract | `200` malformed | `error` |
| Section hợp lệ | `200` | `success` |

Không hiển thị raw backend error, URL nội bộ hoặc stack trace trên UI public.

## 5. Thiết kế type và data layer Webview

### 5.1. Public DTO

Tạo `apps/webview/types/public-page.ts` với type độc lập, không import type từ Frontend/Refine:

```ts
export interface PublicPageSectionItemDto {
  itemType: "CATEGORY" | "POST" | "MEDIA";
  itemId: number;
  sortOrder: number;
  data: Record<string, unknown>;
}

export interface PublicPageSectionDto {
  id: number;
  key: string;
  title: string;
  description: string;
  backgroundColor: string | null;
  sortOrder: number;
  collections: Record<string, PublicPageSectionItemDto[]>;
}

export interface PublicPageDto {
  id: number;
  slug: string;
  sections: PublicPageSectionDto[];
}

export interface PublicPageResponse {
  data: PublicPageDto;
}
```

DTO thực tế có thể khai báo thêm Media/SEO dùng chung theo contract backend, nhưng không dùng `any` và không ép kiểu JSON bằng `as PublicPageResponse`.

### 5.2. Runtime validation

Tạo type guard cho các field Intro thực sự phụ thuộc:

- Root là object và có `data` object.
- `data.id` là number, `data.slug` là string, `data.sections` là array.
- Mỗi section có `id`, `key`, `title`, `description`, `sortOrder` đúng kiểu.
- `collections` là object; Intro chưa đọc item nhưng vẫn kiểm tra object/null boundary đúng contract.

Response sai contract trả `error`; không render dữ liệu một phần không đáng tin cậy.

### 5.3. Loader và normalizer

Tạo `apps/webview/lib/home/get-intro-content.ts`:

```ts
export type IntroContentState = "success" | "empty" | "error";

export interface IntroContentResult {
  state: IntroContentState;
  content: IntroContent | null;
}

export async function getIntroContent(): Promise<IntroContentResult>;
```

Luồng xử lý:

1. Tạo endpoint `/pages/trang-chu` bằng `getWebviewApiEndpoint`.
2. Gọi `fetch(endpoint, { cache: "no-store" })` để đồng bộ với route `force-dynamic` và pattern Hero hiện tại.
3. `404` trả `empty`; non-2xx khác trả `error`.
4. Parse JSON vào `unknown`, validate bằng type guard.
5. Tìm đúng `section.key === "home_creative"`.
6. Trim `section.title` và `section.description`.
7. Nếu title hoặc description rỗng, hoặc Section không có trong response, thì trả `empty` để không tạo heading thiếu nội dung/ARIA label rỗng.
8. Đọc `section.collections.benefits`, loop và map item data thành `BenefitItem`; không dùng benefits từ fixture.
9. Catch network/JSON errors và trả `error`.

Đặt technical constants trong data layer:

```ts
const HOME_PAGE_SLUG = "trang-chu";
const INTRO_SECTION_KEY = "home_creative";
```

Không dùng title/name để tìm Section và không fallback sang Section đầu tiên.

## 6. Component và composition

### 6.1. Presentation component

Giữ `apps/webview/components/home/intro-section.tsx` là component chỉ nhận `IntroContent` và render success UI hiện tại. Component không import config, không gọi `fetch` và không biết Page slug/API DTO.

### 6.2. Async Server Component

Tạo `apps/webview/components/home/intro-section-container.tsx`:

- Không nhận dữ liệu benefits từ fixture.
- `await getIntroContent()`.
- `success`: render `<IntroSection content={result.content} />`.
- `empty`: render section-level empty state “Chưa có nội dung giới thiệu.” với `role="status"`.
- `error`: render section-level error state “Không thể tải nội dung giới thiệu lúc này.” với `role="alert"`.
- Không throw lỗi dự kiến lên route-level error boundary vì lỗi một section không được làm sập toàn bộ Home.

### 6.3. Loading state

Tạo `apps/webview/components/home/intro-section-skeleton.tsx` giữ cùng footprint:

- `section.py-24`, container `max-w-7xl`.
- Hai skeleton bar cho eyebrow/title.
- Grid ba cột với placeholder ảnh 64x64, title và description.
- `aria-busy="true"`, label “Đang tải nội dung giới thiệu”.
- Dùng Tailwind hiện có, không thêm dependency skeleton.

Trong `apps/webview/app/page.tsx`:

```tsx
<Suspense fallback={<IntroSectionSkeleton />}>
  <IntroSectionContainer />
</Suspense>
```

`page.tsx` vẫn là Server Component. Không await Intro ở composition root để phần còn lại có thể stream độc lập.

## 7. Thay đổi fixture và UI types

Sau khi API backend sẵn sàng và integration test pass:

- Xóa `intro` khỏi `apps/webview/data/webview-home.json` vì heading và benefits đều lấy từ API.
- Giữ `IntroContent` là type render:

```ts
export interface IntroContent {
  heading: SectionHeadingContent;
  benefits: BenefitItem[];
}
```

- `IntroContent` chỉ được tạo sau normalizer API; không còn heading hoặc benefits mock làm fallback ngầm.

Nếu API chưa được deploy, chưa xóa heading fixture trong commit backend prerequisite để tránh làm hỏng Home trong khoảng chuyển tiếp. Việc bật integration và xóa mock heading nên nằm trong cùng release Webview.

## 8. Cấu trúc file dự kiến

```text
apps/webview/
├── app/
│   └── page.tsx                              # Suspense boundary cho Intro
├── components/home/
│   ├── intro-section.tsx                    # success/presentation UI
│   ├── intro-section-container.tsx          # async Server Component
│   └── intro-section-skeleton.tsx           # loading UI
├── data/
│   └── webview-home.json                    # không còn intro mock
├── lib/home/
│   └── get-intro-content.ts                 # fetch, guard, normalize, state
└── types/
    ├── home.ts                              # IntroContent
    └── public-page.ts                       # public API DTO
```

Backend prerequisite dự kiến theo plan đã có:

```text
apps/backend/internal/
├── controller/public_page_controller.go
├── infrastructure/repository/public_page_repo.go
└── usecases/publicpage/
    ├── dto/public_page_dto.go
    └── service/public_page_service.go
```

Đồng thời cập nhật `main.go`, Wire và `.docs/api-endpoints.yaml` để route public và schema `sections` là contract chính thức.

## 9. Trình tự triển khai

1. Hoàn thiện Public Page API `GET /api/v1/pages/:slug` theo backend plan 26.
2. Mở rộng read model public trả `sections` theo backend plan 27; bảo đảm chỉ `ACTIVE`, order ổn định và cache invalidation hoạt động.
3. Bổ sung OpenAPI example cho Page `trang-chu` và Section `home_creative`.
4. Thêm DTO/type guard Public Page trong Webview.
5. Tạo `getIntroContent`, kiểm thử success/empty/error bằng payload fixture hoặc mocked `fetch`.
6. Tạo async container và skeleton, đặt Suspense boundary trong `app/page.tsx`.
7. Xóa toàn bộ intro mock khỏi `webview-home.json`, lấy benefits từ `collections.benefits`.
8. Chạy lint, TypeScript, production build và kiểm tra thủ công với Admin update.

## 10. Kiểm thử và tiêu chí nghiệm thu

### Contract/backend

- `GET /api/v1/pages/trang-chu` không cần Authorization và trả Page `PUBLISHED`.
- Response chứa Section `home_creative`, không chứa Section `INACTIVE` hoặc soft-deleted.
- Đổi title/description tại `/pages/edit/6` làm cache `public:pages:slug:trang-chu` bị invalidate sau commit.
- Page DRAFT hoặc không tồn tại trả `404 PAGE_NOT_FOUND` mà không lộ trạng thái nội bộ.
- Public payload không trả `name`, Section status, owner hoặc dữ liệu admin.

### Webview states

- API chậm: skeleton Intro xuất hiện trong Suspense boundary; Hero và các section khác không bị thay bằng skeleton Intro.
- API 200 hợp lệ: eyebrow/title lấy từ Section ID 3, collection `benefits` được loop và render thành các benefit card.
- Không có `home_creative` hoặc Page 404: hiển thị empty state rõ ràng.
- API 5xx/network lỗi/malformed JSON: hiển thị error state, không làm sập toàn trang.
- Tắt Section trong Admin rồi reload: Intro chuyển sang empty state.
- Bật lại/chỉnh nội dung rồi reload: Intro hiển thị nội dung mới, không cần build lại Webview.

### Quality gates

Chạy trong `apps/webview`:

```bash
npm run lint
npx tsc --noEmit
npm run build
```

Kiểm tra thêm:

- Không còn `intro` trong runtime fixture.
- Không có request tới `/api/v1/admin/pages/6/sections` từ Webview.
- Không có `any`, type assertion mù hoặc import type từ `apps/frontend`.
- `page.tsx`, Intro container và Intro presentation component không có `"use client"`.
- Không thay đổi layout hoặc Tailwind classes của success state ngoài phần dữ liệu lấy từ API.

## 11. Rủi ro và hướng mở rộng

- Public Page API hiện chưa có trong runtime; Webview integration không thể hoàn tất an toàn trước prerequisite này.
- Item nguồn cần có metadata public đủ để map title, description và image; item thiếu trường bắt buộc sẽ bị bỏ qua thay vì suy diễn từ `fileName` hoặc lưu JSON trong description.
- Khi nhiều Home section cùng dùng Public Page API, nên tách `getPublicPageBySlug` thành loader dùng chung. Các GET giống nhau trong cùng React tree được Next.js memoize, nhưng vẫn cần một normalizer/contract trung tâm để tránh mỗi section diễn giải payload khác nhau.
- Backend Redis TTL 15 phút chỉ là safety net; invalidation sau mutation Section mới là cơ chế bảo đảm nội dung Admin cập nhật được phản ánh ngay.
