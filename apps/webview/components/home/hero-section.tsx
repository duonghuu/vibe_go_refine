import Link from "next/link";
import type { HeroSectionContent } from "@/types/home";

interface HeroSectionProps {
  content: HeroSectionContent;
  backgroundImage: string | null;
  imageState: "success" | "empty" | "error";
}

export default function HeroSection({ content, backgroundImage, imageState }: HeroSectionProps) {
  return (
    <section
      style={backgroundImage ? { backgroundImage: `url(${backgroundImage})` } : undefined}
      className="relative bg-[#242424] bg-cover bg-[10%_0%] py-48 before:absolute before:inset-0 before:bg-black/50"
      aria-labelledby={content.title ? "hero-title" : undefined}
      aria-label={content.title ? undefined : "Hero"}
    >
      <div className="relative mx-auto max-w-7xl px-5 lg:px-8">
        {content.eyebrow && <span className="mb-3 block text-white">{content.eyebrow}</span>}
        {content.title && (
          <h1 id="hero-title" className={`max-w-3xl text-4xl font-semibold leading-tight text-white sm:text-6xl lg:text-7xl ${content.description ? "mb-4" : "mb-10"}`}>
            {content.title}
          </h1>
        )}
        {content.description && <p className="mb-10 max-w-2xl text-lg text-white/80">{content.description}</p>}
        {content.cta && <Link href={content.cta.href} className="inline-flex items-center rounded-full bg-[#f75757] px-10 py-4 text-base uppercase !text-white transition-colors hover:bg-[#dd0b0b]">{content.cta.label} <i className="fa fa-angle-right ml-3 text-base" /></Link>}
        {imageState === "empty" && (
          <p className="mt-6 text-sm text-white/80" role="status">Chưa có hình ảnh hero được cấu hình.</p>
        )}
        {imageState === "error" && (
          <p className="mt-6 text-sm text-white/80" role="alert">Không thể tải hình ảnh hero lúc này.</p>
        )}
      </div>
    </section>
  );
}
