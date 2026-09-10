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
  cta: LinkContent;
}

export interface BrandContent {
  name: string;
  accent: string;
  href: string;
}

export interface NavigationLink extends LinkContent {
  kind: "link";
}

export interface NavigationDropdown {
  kind: "dropdown";
  label: string;
  children: LinkContent[];
}

export type NavigationItem = NavigationLink | NavigationDropdown;

export interface HeaderContent {
  brand: BrandContent;
  socialLinks: SocialLink[];
  phoneLabel: string;
  phone: LinkContent;
  email: LinkContent;
  navigation: NavigationItem[];
  quoteCta: LinkContent;
}

export interface HeroContent {
  backgroundImage: string;
  eyebrow: string;
  titleLines: string[];
  cta: LinkContent;
}

export interface BenefitItem {
  icon: string;
  title: string;
  description: string;
}

export interface IntroContent {
  heading: SectionHeadingContent;
  benefits: BenefitItem[];
}

export interface AboutContent {
  image: ImageContent;
  heading: SectionHeadingContent;
  featureTitle: string;
  description: string;
  cta: LinkContent;
}

export interface ServicesContent {
  heading: SectionHeadingContent;
  items: ServiceItem[];
}

export interface ContactCtaContent {
  backgroundImage: string;
  heading: SectionHeadingContent;
  description: string;
  phoneIcon: string;
  phone: LinkContent;
}

export interface TestimonialsContent {
  heading: SectionHeadingContent;
  items: TestimonialItem[];
}

export interface LatestArticlesContent {
  backgroundImage: string;
  heading: SectionHeadingContent;
  items: ArticleItem[];
}

export interface QuoteCtaContent {
  heading: SectionHeadingContent;
  cta: LinkContent;
}

export interface FooterLinkGroup {
  title: string;
  links: LinkContent[];
}

export interface FooterContent {
  brand: BrandContent;
  linkGroups: FooterLinkGroup[];
  subscribe: {
    title: string;
    description: string;
    placeholder: string;
    buttonLabel: string;
  };
  email: LinkContent;
  phone: LinkContent;
  copyright: {
    prefix: string;
    brand: string;
    by: string;
    link: LinkContent;
  };
  distributedBy: {
    prefix: string;
    link: LinkContent;
  };
  socialLinks: SocialLink[];
}

export interface HomeContent {
  header: HeaderContent;
  hero: HeroContent;
  intro: IntroContent;
  about: AboutContent;
  stats: StatItem[];
  services: ServicesContent;
  contactCta: ContactCtaContent;
  testimonials: TestimonialsContent;
  latestArticles: LatestArticlesContent;
  quoteCta: QuoteCtaContent;
  footer: FooterContent;
}
