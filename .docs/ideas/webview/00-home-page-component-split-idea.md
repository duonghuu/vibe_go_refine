# Ý TƯỞNG: Chia nhỏ trang chủ Webview thành component

> Phạm vi tài liệu: chỉ đề xuất ý tưởng và định hướng; không bao gồm triển khai code.

## 1. Mục tiêu

- Chia nhỏ `apps/webview/app/page.tsx` thành các component có trách nhiệm rõ ràng, dễ đọc, dễ kiểm thử và dễ thay dữ liệu fixture bằng Public Page API sau này.
- Giữ `app/page.tsx` là Server Component và chỉ dùng Client Component tại những nhánh thật sự cần state, effect hoặc browser API.
- Giữ nguyên giao diện, thứ tự section, anchor, responsive behavior và asset hiện tại trong bước refactor đầu tiên.
- Chuẩn bị ranh giới dữ liệu theo hướng `API DTO -> normalized model -> UI props`, nhưng không nối API thật trong phạm vi của việc chia component.

## 2. Hiện trạng

`apps/webview/app/page.tsx` hiện trực tiếp render:

1. Header.
2. Hero.
3. Intro/benefits.
4. About.
5. Statistics.
6. Services.
7. Contact CTA.
8. Testimonials.
9. Latest articles.
10. Quote CTA.
11. Footer.

Ba phần tương tác đã được tách:

- `InteractiveHeader`: menu mobile và dropdown.
- `AnimatedCounter`: `IntersectionObserver` và animation frame.
- `TestimonialCarousel`: resize listener, timer và điều hướng carousel.

Các vấn đề chính:

- Một file route đang giữ quá nhiều markup và chi tiết trình bày.
- Dữ liệu, type và việc đọc JSON cùng nằm trong `components/home-data.ts`.
- Header/footer chưa có ranh giới layout rõ ràng để tái sử dụng cho route khác.
- Nhiều mẫu UI lặp lại như eyebrow + heading, container, CTA và màu thương hiệu.
- Public Page/Section contract hiện chưa biểu diễn đủ toàn bộ nội dung của mockup. Các trường như CTA label/link, số điện thoại, icon dịch vụ, testimonial quote/name/role chưa có mapping trực tiếp trong `PublicPageSectionResponse`.

## 3. Nguyên tắc chia tách

### 3.1. Tách theo section nghiệp vụ

Mỗi section cấp cao của landing page là một component riêng. Không tách thành component chỉ vì một đoạn JSX dài vài dòng; chỉ tách khi phần đó có ít nhất một trong các đặc điểm:

- Có ý nghĩa nội dung độc lập.
- Có props/data contract riêng.
- Có thể được bật/tắt hoặc sắp xếp bằng CMS trong tương lai.
- Có logic tương tác hoặc cách xử lý media riêng.
- Có khả năng tái sử dụng trên route khác.

### 3.2. Server mặc định, Client ở lá

- `app/page.tsx`, footer và các section tĩnh tiếp tục là Server Component.
- `AnimatedCounter` chỉ hydrate con số, không biến toàn bộ `StatisticsSection` thành Client Component.
- `TestimonialCarousel` là client leaf nằm bên trong `TestimonialsSection` server.
- Header nên tách phần top bar tĩnh khỏi navigation tương tác. Chỉ `HeaderNavigation` cần `"use client"`.
- Form subscribe ở footer vẫn là markup server trong giai đoạn hiện tại. Khi có submit thật, tách riêng `NewsletterForm` dùng Server Action hoặc Client Component theo contract được duyệt.

### 3.3. Component trình bày không đọc dữ liệu

- Section chỉ nhận props typed và render UI.
- Không import JSON, gọi API hoặc normalize response bên trong section.
- Route hoặc server data layer chịu trách nhiệm chọn nguồn dữ liệu và truyền model đã chuẩn hóa xuống component.
- Không dùng `any`; dữ liệu từ API sau này phải đi qua DTO/type guard/normalizer.

### 3.4. Không trừu tượng hóa quá sớm

- Chỉ tạo `SectionHeading` cho mẫu eyebrow + heading đang thực sự lặp lại.
- Chưa tạo một component chung kiểu `GenericSection` với nhiều `className`/boolean để bao phủ mọi layout.
- Chưa tạo registry CMS ở bước refactor thuần cấu trúc. Registry chỉ nên xuất hiện khi route bắt đầu render danh sách `sections` từ Public API.
- Giữ Tailwind class sát component sở hữu layout; token màu lặp lại tiếp tục dùng CSS variables hiện có khi triển khai.

## 4. Cấu trúc file đề xuất

```text
apps/webview/
├── app/
│   └── page.tsx
├── components/
│   ├── site/
│   │   ├── site-header.tsx
│   │   ├── header-top-bar.tsx
│   │   ├── header-navigation.tsx       # Client Component
│   │   └── site-footer.tsx
│   └── home/
│       ├── section-heading.tsx
│       ├── hero-section.tsx
│       ├── intro-section.tsx
│       ├── about-section.tsx
│       ├── statistics-section.tsx
│       ├── animated-counter.tsx        # Client Component
│       ├── services-section.tsx
│       ├── contact-cta-section.tsx
│       ├── testimonials-section.tsx
│       ├── testimonial-carousel.tsx    # Client Component
│       ├── latest-articles-section.tsx
│       └── quote-cta-section.tsx
├── data/
│   └── webview-home.json
├── lib/
│   └── home/
│       └── get-home-content.ts
└── types/
    └── home.ts
```

Không cần tạo `home-page.tsx`: chính `app/page.tsx` đã là composition root của trang chủ. Thêm một wrapper chỉ chuyển toàn bộ JSX sang file khác mà không tạo ranh giới trách nhiệm mới.

## 5. Trách nhiệm từng thành phần

| Thành phần | Trách nhiệm | Boundary |
| --- | --- | --- |
| `app/page.tsx` | Lấy/chuẩn hóa dữ liệu, compose header, các section và footer theo thứ tự | Server |
| `site-header.tsx` | Ghép top bar tĩnh và navigation | Server |
| `header-top-bar.tsx` | Social, phone, email | Server |
| `header-navigation.tsx` | Menu mobile, dropdown, đóng menu khi điều hướng | Client |
| `site-footer.tsx` | Link groups, thông tin thương hiệu, subscribe markup, copyright | Server |
| `section-heading.tsx` | Render eyebrow + heading với biến thể `left/center` và `default/inverse` | Server |
| `hero-section.tsx` | Hero background, overlay, title và primary CTA | Server |
| `intro-section.tsx` | Heading và danh sách lợi ích | Server |
| `about-section.tsx` | Feature image, nội dung giới thiệu và CTA | Server |
| `statistics-section.tsx` | Grid statistics; truyền từng giá trị cho counter | Server + client leaf |
| `services-section.tsx` | Grid dịch vụ | Server |
| `contact-cta-section.tsx` | Banner nền và contact card | Server |
| `testimonials-section.tsx` | Heading và chuẩn bị danh sách testimonial | Server + client leaf |
| `latest-articles-section.tsx` | Danh sách bài mới, `next/image`, metadata hiển thị | Server |
| `quote-cta-section.tsx` | CTA cuối trang | Server |
| `get-home-content.ts` | Đọc fixture hiện tại; sau này là điểm chuyển sang Public Page API và normalizer | Server-only data layer |
| `types/home.ts` | Normalized model và props dùng trong webview | Không phụ thuộc Refine/MUI |

## 6. Model dữ liệu đề xuất

Type UI chỉ chứa đúng dữ liệu component cần, không dùng trực tiếp DTO backend:

```ts
export interface LinkContent {
  label: string;
  href: string;
}

export interface SectionHeadingContent {
  eyebrow: string;
  title: string;
}

export interface StatItem {
  value: number;
  suffix: string;
  label: string;
}

export interface ServiceItem {
  icon: string;
  title: string;
  description: string;
}

export interface TestimonialItem {
  quote: string;
  name: string;
  role: string;
}

export interface ArticleItem {
  image: string;
  title: string;
  categoryLabels: string[];
  authorName: string;
  href: string;
}
```

Khi mở rộng fixture, có thể tạo `HomeContent` chứa dữ liệu của toàn bộ section. Mỗi section vẫn chỉ nhận phần dữ liệu của chính nó, ví dụ:

```ts
interface ServicesSectionProps {
  heading: SectionHeadingContent;
  items: ServiceItem[];
}

interface StatisticsSectionProps {
  items: StatItem[];
}
```

Không truyền toàn bộ `HomeContent` vào mọi section vì sẽ làm props phụ thuộc rộng và khó tái sử dụng.

## 7. Hình dạng `app/page.tsx` sau refactor

```tsx
export default async function HomePage() {
  const content = await getHomeContent();

  return (
    <>
      <SiteHeader />
      <main>
        <HeroSection content={content.hero} />
        <IntroSection content={content.intro} />
        <AboutSection content={content.about} />
        <StatisticsSection items={content.stats} />
        <ServicesSection content={content.services} />
        <ContactCtaSection content={content.contact} />
        <TestimonialsSection content={content.testimonials} />
        <LatestArticlesSection content={content.latestArticles} />
        <QuoteCtaSection content={content.quote} />
      </main>
      <SiteFooter content={content.footer} />
    </>
  );
}
```

`getHomeContent()` có thể là hàm đồng bộ trong thời gian dùng JSON; chữ ký async được cân nhắc chỉ khi chuyển sang server `fetch`. Không cần cố tình khai báo async trước khi có I/O thật.

## 8. Chiến lược chuyển đổi sang CMS sau này

### Giai đoạn 1 — Refactor không đổi hành vi

- Extract từng section từ JSX hiện tại.
- Di chuyển type khỏi `components/home-data.ts` sang `types/home.ts`.
- Đưa việc đọc `webview-home.json` vào data layer.
- Tách header static/client boundary và footer.
- Giữ nguyên text, asset, anchor ID, DOM semantics và Tailwind classes.

Đây là giai đoạn nên thực hiện trước vì có thể review bằng visual regression và không phụ thuộc thay đổi backend.

### Giai đoạn 2 — Hoàn thiện fixture theo model mới

- Bổ sung dữ liệu hero, intro, about, contact, quote và footer vào fixture runtime nếu các nội dung này cần cấu hình.
- Thêm normalizer kiểm tra/fallback dữ liệu thay vì dùng ép kiểu `as SomeType[]` tại module data.
- Giữ `.docs/mock-data/webview-home.json` làm tài liệu tham khảo; runtime chỉ import `apps/webview/data/webview-home.json` để tránh hai nguồn sự thật.

### Giai đoạn 3 — Nối Public Page API

- Khai báo `PublicPageDto` và `PublicPageSectionDto` riêng với normalized `HomeContent`.
- Fetch tại `app/page.tsx` hoặc server-only data layer.
- Chỉ nhận Page `PUBLISHED` và Section `ACTIVE` đã được backend lọc; webview vẫn kiểm tra phòng thủ với key không hỗ trợ/payload lỗi.
- Dùng registry tường minh `section.key -> renderer`, không suy luận component từ title:

```text
hero                -> HeroSection
intro               -> IntroSection
about               -> AboutSection
statistics          -> StatisticsSection
services            -> ServicesSection
contact_cta         -> ContactCtaSection
testimonials        -> TestimonialsSection
latest_articles     -> LatestArticlesSection
quote_cta           -> QuoteCtaSection
```

- Render theo `sortOrder`; key không hỗ trợ được bỏ qua có chủ đích hoặc dùng fallback đã thiết kế, không làm sập toàn trang.
- Không render `content` HTML bằng `dangerouslySetInnerHTML` nếu chưa có sanitizer được dự án phê duyệt.

### Khoảng trống contract cần quyết định trước Giai đoạn 3

`PublicPageSectionResponse` hiện chỉ có metadata chung, media và collection `CATEGORY/POST/MEDIA`. Nó chưa đủ cho mọi section trong mockup. Cần chọn một trong hai hướng ở một plan backend riêng:

1. Mở rộng section schema bằng payload có type/schema được validate cho các trường CTA, icon và testimonial.
2. Quy ước các nội dung đó thành Post/Category/Media và bổ sung mapping rõ ràng, chấp nhận giới hạn của từng loại nguồn.

Không nên nhét dữ liệu JSON tùy ý vào `description`, suy luận icon từ title hoặc hard-code một nửa nội dung trong component và lấy nửa còn lại từ API.

## 9. Thứ tự triển khai khuyến nghị

1. Chụp baseline giao diện desktop/mobile và ghi nhận anchor hiện tại.
2. Tạo `types/home.ts` và data loader; giữ dữ liệu runtime không đổi.
3. Extract các section tĩnh, bắt đầu từ hero đến quote CTA.
4. Di chuyển hai client leaf vào `components/home` và cập nhật import.
5. Tách header thành server shell + client navigation.
6. Extract footer nhưng chưa đưa vào `app/layout.tsx` cho đến khi có thêm route xác nhận cùng dùng một site shell.
7. Rút gọn `app/page.tsx` thành composition root.
8. Chạy lint/build và so sánh giao diện với mockup ở các breakpoint chính.

Mỗi bước nên là thay đổi cơ học, không đồng thời sửa nội dung, style hoặc nối API để diff dễ review và rollback.

## 10. Tiêu chí hoàn thành

- `apps/webview/app/page.tsx` chỉ còn trách nhiệm lấy dữ liệu và compose page; không chứa markup chi tiết của section.
- Mỗi section cấp cao có file riêng, tên file kebab-case và component PascalCase.
- Component trình bày không import JSON và không gọi API.
- `page.tsx` cùng mọi section tĩnh vẫn là Server Component.
- Chỉ navigation, counter và carousel tạo client bundle.
- Không dùng `any`; props và model đều typed tường minh.
- Giao diện, nội dung, thứ tự section, anchor, link, asset và responsive behavior không thay đổi so với trước refactor.
- `next/image` tiếp tục được dùng cho ảnh bài viết; media có kích thước/`sizes` phù hợp khi triển khai.
- `npm run lint` và `npm run build` trong `apps/webview` đều pass.
- Có kiểm tra thủ công menu mobile, dropdown, counter, carousel và điều hướng anchor.

## 11. Ngoài phạm vi

- Không nối Public Page API trong task chia component.
- Không thay đổi backend, database hoặc API contract.
- Không thiết kế lại màu, typography, spacing hay nội dung.
- Không chuyển header/footer vào root layout trước khi xác định hành vi của các route khác, loading, error và not-found.
- Không đưa Bootstrap, jQuery, MUI hoặc Refine vào runtime webview.
- Không xóa asset/plugin legacy trong cùng đợt refactor; việc audit và dọn asset nên là task riêng.

## 12. Kết luận

Hướng phù hợp nhất là **feature-first theo từng section**, với `app/page.tsx` làm composition root, section tĩnh là Server Component và ba vùng tương tác được giữ ở các client leaf nhỏ. Refactor nên hoàn tất độc lập trước khi thêm registry CMS. Cách này giảm rủi ro thay đổi giao diện, làm rõ data boundary và tạo đường chuyển sang Public Page API mà không buộc component UI phụ thuộc vào contract backend hiện chưa đủ biểu đạt landing page.
