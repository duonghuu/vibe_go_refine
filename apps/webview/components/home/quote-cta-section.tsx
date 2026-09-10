import Link from "next/link";
import type { QuoteCtaContent } from "@/types/home";

interface QuoteCtaSectionProps {
  content: QuoteCtaContent;
}

export default function QuoteCtaSection({ content }: QuoteCtaSectionProps) {
  return (
    <section id="quote" className="relative py-20">
      <div className="mx-auto max-w-7xl px-5 lg:px-8">
        <div className="flex flex-col items-start justify-between gap-8 rounded border border-black/5 bg-[#f5f8f9] p-10 md:flex-row md:items-center">
          <div><span className="text-sm uppercase tracking-wide text-[#f75757]">{content.heading.eyebrow}</span><h2 className="mt-2 text-2xl font-semibold text-[#242424]">{content.heading.title}</h2></div>
          <Link href={content.cta.href} className="shrink-0 rounded-full bg-[#f75757] px-8 py-4 text-xs uppercase text-white transition-colors hover:bg-[#dd0b0b]">{content.cta.label}</Link>
        </div>
      </div>
    </section>
  );
}
