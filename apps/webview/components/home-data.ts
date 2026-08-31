import homeData from "@/data/webview-home.json";

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
}

export const stats = homeData.stats as StatItem[];
export const services = homeData.services as ServiceItem[];
export const testimonials = homeData.testimonials as TestimonialItem[];
export const articles = homeData.articles as ArticleItem[];
