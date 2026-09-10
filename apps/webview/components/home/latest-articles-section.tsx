import Image from "next/image";
import Link from "next/link";
import type { LatestArticlesContent } from "@/types/home";
import SectionHeading from "./section-heading";

interface LatestArticlesSectionProps {
  content: LatestArticlesContent;
}

export default function LatestArticlesSection({ content }: LatestArticlesSectionProps) {
  return (
    <section id="latest-news" style={{ backgroundImage: `url(${content.backgroundImage})` }} className="relative bg-cover bg-center py-24 before:absolute before:inset-0 before:bg-black/80" aria-labelledby="news-title">
      <div className="relative mx-auto max-w-7xl px-5 lg:px-8">
        <SectionHeading id="news-title" heading={content.heading} align="center" tone="inverse" className="mx-auto mb-16 max-w-2xl" />
        <div className="grid gap-10 md:grid-cols-3">
          {content.items.map((article) => (
            <article key={article.title}>
              <Image src={article.image.src} alt={article.image.alt} width={1600} height={1088} className="h-56 w-full rounded object-cover" />
              <div className="mt-5">
                <div className="text-sm text-white/50">
                  {article.categories.map((category) => <span key={category}><Link href={article.href}>{category}</Link><span className="mx-2">/</span></span>)}
                  <span><i className="fa fa-user mr-2" />{article.authorName}</span>
                </div>
                <h3 className="mt-4 mb-8 text-2xl font-semibold leading-snug text-white"><Link href={article.href}>{article.title}</Link></h3>
                <Link href={article.cta.href} className="inline-flex rounded-full border-2 border-[#f75757] px-6 py-3 text-xs uppercase text-white transition-colors hover:bg-[#f75757]">{article.cta.label}</Link>
              </div>
            </article>
          ))}
        </div>
      </div>
    </section>
  );
}
