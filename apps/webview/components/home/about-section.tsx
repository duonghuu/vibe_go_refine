import Link from "next/link";
import type { AboutContent } from "@/types/home";

interface AboutSectionProps {
  content: AboutContent;
}

export default function AboutSection({ content }: AboutSectionProps) {
  return (
    <section id="about" className="relative py-24" aria-labelledby="about-title">
      <div style={{ backgroundImage: `url(${content.image.src})` }} className="absolute inset-y-0 left-0 hidden w-[45%] bg-cover bg-center lg:block" />
      <div className="mx-auto max-w-7xl px-5 lg:px-8">
        <div className="ml-auto max-w-2xl lg:pl-20">
          <span className="text-sm uppercase tracking-wide text-[#f75757]">{content.heading.eyebrow}</span>
          <h2 id="about-title" className="relative mt-3 mb-8 text-3xl font-semibold leading-tight text-[#242424] md:text-4xl">{content.heading.title}</h2>
          <div className="relative pl-0 md:pl-20">
            <h3 className="mb-3 text-xl font-semibold text-[#242424]">{content.featureTitle}</h3>
            <p className="mb-10 leading-8 text-black/65">{content.description}</p>
            <Link href={content.cta.href} className="inline-flex rounded-full bg-[#f75757] px-8 py-4 text-xs uppercase text-white transition-colors hover:bg-[#dd0b0b]">{content.cta.label}</Link>
          </div>
        </div>
      </div>
    </section>
  );
}
