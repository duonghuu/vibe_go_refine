# WEBVIEW PLAN: Chia nhỏ trang chủ thành component

## 1. Mục tiêu và phạm vi

Refactor `apps/webview/app/page.tsx` thành composition root cho trang chủ, chia markup theo site shell và từng section nghiệp vụ. Kết quả phải giữ nguyên giao diện, nội dung, asset, thứ tự section, anchor và responsive behavior hiện tại.

Phạm vi gồm:

- Chuẩn hóa toàn bộ nội dung đang hard-code và dữ liệu JSON hiện có thành fixture typed cho trang chủ.
- Tạo các component section độc lập; component UI chỉ nhận props đã typed.
- Giữ Server Component làm mặc định; chỉ hydration phần navigation, counter và carousel.
- Xác nhận bằng lint, build và kiểm tra thủ công các tương tác đang có.

Ngoài phạm vi:

- Không gọi Public Page/Section API, không tạo CMS registry và không sửa backend/API contract.
- Không đổi nội dung, thiết kế, Tailwind class, asset, metadata, cache hay cấu trúc route.
- Không chuyển header/footer vào `app/layout.tsx`; chưa có route thứ hai xác nhận dùng chung site shell.
- Không thêm dependency hoặc dọn plugin/asset legacy trong cùng thay đổi.

## 2. Hiện trạng và quyết định kiến trúc

`page.tsx` hiện render trực tiếp Header, 9 section nội dung và Footer. Các client leaf đã tồn tại là menu/header, animated counter và testimonial carousel; fixture hiện chỉ chứa stats, services, testimonials và articles.

Áp dụng các quyết định sau:

| Khu vực | Quyết định |
| --- | --- |
| Route `/` | `app/page.tsx` tiếp tục là Server Component, chỉ lấy `HomeContent` và compose UI. |
| Site shell | Tách header thành Server shell + top bar server + navigation client; footer là Server Component. |
| Home sections | Mỗi section cấp cao là một Server Component riêng, giữ nguyên semantic HTML và class Tailwind hiện tại. |
| Client boundary | Chỉ `HeaderNavigation`, `AnimatedCounter`, `TestimonialCarousel` dùng `"use client"`; props truyền xuống phải serializable. |
| Nguồn runtime | Chỉ `apps/webview/data/webview-home.json` được import lúc chạy; `.docs/mock-data/webview-home.json` chỉ là tài liệu tham khảo. |
| Data layer | `getHomeContent()` đọc fixture và trả `HomeContent` typed đồng bộ; component không tự đọc JSON hoặc normalize data. |
| Tái sử dụng | Chỉ tạo `SectionHeading` cho mẫu eyebrow + title lặp lại; không tạo generic section hoặc CMS registry sớm. |

## 3. Cấu trúc file sau refactor

```text
apps/webview/
├── app/
│   └── page.tsx
├── components/
│   ├── site/
│   │   ├── header-navigation.tsx       # Client Component
│   │   ├── header-top-bar.tsx
│   │   ├── site-footer.tsx
│   │   └── site-header.tsx
│   └── home/
│       ├── about-section.tsx
│       ├── animated-counter.tsx        # Client Component
│       ├── contact-cta-section.tsx
│       ├── hero-section.tsx
│       ├── intro-section.tsx
│       ├── latest-articles-section.tsx
│       ├── quote-cta-section.tsx
│       ├── section-heading.tsx
│       ├── services-section.tsx
│       ├── statistics-section.tsx
│       ├── testimonial-carousel.tsx    # Client Component
│       └── testimonials-section.tsx
├── data/
│   └── webview-home.json
├── lib/home/
│   └── get-home-content.ts
└── types/
    └── home.ts
```

Xóa sau khi migrate import:

- `components/interactive-header.tsx`
- `components/animated-counter.tsx`
- `components/testimonial-carousel.tsx`
- `components/home-data.ts`

Không tạo `home-page.tsx`; file này chỉ làm trung gian mà không tạo ranh giới trách nhiệm mới.

## 4. Data contract

### 4.1. Type UI

Tạo `types/home.ts` với các type tường minh sau:

```ts
export interface LinkContent {
  label: string;
  href: string;
}

export interface SocialLink extends LinkContent {
  icon: string;
  ariaLabel: string;
}

export interface SectionHeadingContent {
  eyebrow: string;
  title: string;
}

export interface ImageContent {
  src: string;
  alt: string;
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
  image: ImageContent;
  title: string;
  categories: string[];
  authorName: string;
  href: string;
}
```

`HomeContent` phải chứa đủ dữ liệu cho `header`, `hero`, `intro`, `about`, `stats`, `services`, `contactCta`, `testimonials`, `latestArticles`, `quoteCta` và `footer`.

Quy ước dữ liệu bắt buộc:

- Hero lưu `titleLines: string[]`, sau đó component tự chèn line break; tuyệt đối không lưu JSX hoặc HTML trong JSON.
- Navigation dùng discriminated union `link`/`dropdown`; dropdown có danh sách `LinkContent` con.
- Mọi CTA, phone, mail, social URL, footer link và background image chuyển vào fixture.
- Ảnh article dùng `ImageContent`; `alt` được thiết lập có chủ đích. Các ảnh decorative có thể dùng chuỗi rỗng.
- Stable technical ID như `hero-title`, `about`, `services`, `contact`, `latest-news`, `quote` giữ trong component, không lấy từ data.

### 4.2. Fixture và loader

- Mở rộng `data/webview-home.json` theo `HomeContent`; giữ nguyên mọi giá trị text/URL/asset đang hiển thị.
- Tạo `lib/home/get-home-content.ts`, import fixture và kiểm tra cấu trúc qua gán typed `const homeContent: HomeContent = homeData`.
- `getHomeContent(): HomeContent` trả fixture đồng bộ. Không dùng `as SomeType[]`, `any`, network request hoặc API response DTO trong task này.
- Khi background URL đến từ fixture, section dùng `style={{ backgroundImage: ... }}`; Tailwind giữ các class `bg-cover`, position và overlay hiện tại để không làm thay đổi cách hiển thị.

## 5. Component responsibility và cách compose

| Component | Props/dữ liệu nhận | Trách nhiệm |
| --- | --- | --- |
| `SiteHeader` | `HomeContent["header"]` | Ghép `HeaderTopBar` và `HeaderNavigation`. |
| `HeaderTopBar` | social links, phone, email | Render social/contact tĩnh. |
| `HeaderNavigation` | brand, navigation, quote CTA | State menu mobile/dropdown và đóng menu khi chọn link. |
| `HeroSection` | hero data | Background, overlay, eyebrow, h1 và CTA. |
| `IntroSection` | heading, benefit items | Heading và grid 3 benefits. |
| `AboutSection` | about data | Ảnh feature desktop, copy và CTA. |
| `StatisticsSection` | `StatItem[]` | Grid 4 counter; gửi từng item vào client leaf. |
| `ServicesSection` | heading, `ServiceItem[]` | Grid service. |
| `ContactCtaSection` | contact CTA data | Banner nền và contact card. |
| `TestimonialsSection` | heading, `TestimonialItem[]` | Heading và truyền list vào carousel. |
| `LatestArticlesSection` | heading, `ArticleItem[]` | Article grid; tiếp tục dùng `next/image`. |
| `QuoteCtaSection` | quote CTA data | Banner CTA cuối nội dung chính. |
| `SiteFooter` | footer data | Link groups, subscribe markup, brand/contact/copyright/social. |
| `SectionHeading` | `id`, heading, `align`, `tone` | Mẫu eyebrow + title với hai biến thể `left/center`, `default/inverse`. |

`app/page.tsx` có hình dạng sau:

```tsx
export default function HomePage() {
  const content = getHomeContent();

  return (
    <>
      <SiteHeader content={content.header} />
      <main>
        <HeroSection content={content.hero} />
        <IntroSection content={content.intro} />
        <AboutSection content={content.about} />
        <StatisticsSection items={content.stats} />
        <ServicesSection content={content.services} />
        <ContactCtaSection content={content.contactCta} />
        <TestimonialsSection content={content.testimonials} />
        <LatestArticlesSection content={content.latestArticles} />
        <QuoteCtaSection content={content.quoteCta} />
      </main>
      <SiteFooter content={content.footer} />
    </>
  );
}
```

## 6. Trình tự triển khai

1. Ghi nhận baseline UI tại mobile, tablet và desktop; kiểm tra các anchor và hành vi header/carousel/counter hiện tại.
2. Tạo `HomeContent` và các type item; mở rộng fixture bằng chính dữ liệu đang render, không đổi giá trị hiển thị.
3. Thêm `getHomeContent()` và chuyển route sang đọc data ở đúng một điểm.
4. Extract lần lượt các section tĩnh, sao chép nguyên markup/Tailwind/ARIA hiện có rồi thay hard-code bằng props typed.
5. Tạo `SectionHeading` và chỉ thay thế ở những section có cấu trúc eyebrow + heading tương đương.
6. Tách header: phần top bar là Server Component, phần state/menu là `HeaderNavigation` client. Di chuyển counter/carousel hiện có vào `components/home` mà không đổi logic tương tác.
7. Extract footer, chuyển dữ liệu liên quan vào fixture và giữ subscribe ở dạng presentational form không submit.
8. Rút gọn `page.tsx`, kiểm tra không còn import JSON trực tiếp trong component presentation và xóa các file component cũ.

## 7. Kiểm thử và tiêu chí nghiệm thu

Chạy trong `apps/webview`:

```bash
npm run lint
npm run build
```

Kiểm tra thủ công:

- Giao diện desktop, tablet và mobile không đổi về text, ảnh, màu, spacing, grid, overlay và thứ tự section.
- Menu hamburger, dropdown About/Blog và thao tác đóng menu sau điều hướng hoạt động như cũ.
- Counter chỉ bắt đầu khi vào viewport; carousel tự chuyển, đổi theo breakpoint và các dot hoạt động.
- Toàn bộ anchor/link, mailto/tel và external link vẫn đúng; external link giữ `rel="noreferrer"`.
- Chỉ có một `h1`; tất cả `aria-labelledby`, label và trạng thái `aria-expanded` được giữ.
- `next/image` vẫn được dùng cho article ảnh; component chỉ nhận `src` và `alt` typed.
- `page.tsx`, site shell và section tĩnh không có `"use client"`; chỉ ba client leaf tạo client bundle.
- Không có `any`, type assertion mảng từ JSON hoặc import Refine/MUI vào webview.

## 8. Hướng tích hợp CMS sau refactor

Task sau mới định nghĩa `PublicPageDto`/`PublicPageSectionDto`, normalizer và registry tường minh `section.key -> component`. Registry sẽ render theo `sortOrder`, bỏ qua key/payload không hỗ trợ có chủ đích và không dùng title để suy luận component.

Public Section contract hiện chưa thể hiện CTA, icon dịch vụ, phone hoặc testimonial. Trước khi tích hợp API, cần có plan backend riêng để mở rộng schema có validation hoặc quy ước nguồn dữ liệu rõ ràng. Không dùng `description` để nhét JSON tùy ý và không render HTML chưa sanitize.
