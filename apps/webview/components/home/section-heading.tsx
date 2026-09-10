import type { SectionHeadingContent } from "@/types/home";

interface SectionHeadingProps {
  id: string;
  heading: SectionHeadingContent;
  align: "left" | "center";
  tone: "default" | "inverse";
  className?: string;
}

export default function SectionHeading({ id, heading, align, tone, className = "" }: SectionHeadingProps) {
  const alignmentClass = align === "center" ? "text-center" : "";
  const titleClass = tone === "inverse" ? "text-white" : "text-[#242424]";

  return (
    <div className={`${alignmentClass} ${className}`}>
      <span className="text-sm uppercase tracking-wide text-[#f75757]">{heading.eyebrow}</span>
      <h2 id={id} className={`mt-3 text-3xl font-semibold leading-tight md:text-4xl ${titleClass}`}>{heading.title}</h2>
    </div>
  );
}
