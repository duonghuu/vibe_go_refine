import AboutSection from "@/components/home/about-section";
import ContactCtaSection from "@/components/home/contact-cta-section";
import HeroSection from "@/components/home/hero-section";
import IntroSection from "@/components/home/intro-section";
import LatestArticlesSection from "@/components/home/latest-articles-section";
import QuoteCtaSection from "@/components/home/quote-cta-section";
import ServicesSection from "@/components/home/services-section";
import StatisticsSection from "@/components/home/statistics-section";
import TestimonialsSection from "@/components/home/testimonials-section";
import SiteFooter from "@/components/site/site-footer";
import SiteHeader from "@/components/site/site-header";
import { getHomeContent } from "@/lib/home/get-home-content";

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
