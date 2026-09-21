import Image from "next/image";
import type { IntroContent } from "@/types/home";
import SectionHeading from "./section-heading";

interface IntroSectionProps {
  content: IntroContent;
}

export default function IntroSection({ content }: IntroSectionProps) {
  return (
    <section className="py-24" aria-labelledby="intro-title">
      <div className="mx-auto max-w-7xl px-5 lg:px-8">
        <SectionHeading id="intro-title" heading={content.heading} align="left" tone="default" className="mb-16 max-w-3xl" />
        <div className="grid gap-10 md:grid-cols-3">
          {content.benefits.map((benefit) => (
            <article key={benefit.title}>
              <Image src={benefit.image.src} alt={benefit.image.alt} width={64} height={64} unoptimized className="h-16 w-16 object-contain" />
              <h3 className="mt-6 mb-3 text-xl font-semibold text-[#242424]">{benefit.title}</h3>
              <p className="leading-8 text-black/65">{benefit.description}</p>
            </article>
          ))}
        </div>
      </div>
    </section>
  );
}
