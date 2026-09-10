import Link from "next/link";
import type { HeroContent } from "@/types/home";

interface HeroSectionProps {
  content: HeroContent;
}

export default function HeroSection({ content }: HeroSectionProps) {
  return (
    <section style={{ backgroundImage: `url(${content.backgroundImage})` }} className="relative bg-cover bg-[10%_0%] py-48 before:absolute before:inset-0 before:bg-black/50" aria-labelledby="hero-title">
      <div className="relative mx-auto max-w-7xl px-5 lg:px-8">
        <span className="mb-3 block text-white">{content.eyebrow}</span>
        <h1 id="hero-title" className="mb-10 max-w-3xl text-4xl font-semibold leading-tight text-white sm:text-6xl lg:text-7xl">
          {content.titleLines.map((line, index) => <span key={line}>{index > 0 && <br />}{line}</span>)}
        </h1>
        <Link href={content.cta.href} className="inline-flex items-center rounded-full bg-[#f75757] px-10 py-4 text-xs uppercase text-white transition-colors hover:bg-[#dd0b0b]">{content.cta.label} <i className="fa fa-angle-right ml-3 text-base" /></Link>
      </div>
    </section>
  );
}
