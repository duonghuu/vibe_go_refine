import type { ServicesContent } from "@/types/home";
import SectionHeading from "./section-heading";

interface ServicesSectionProps {
  content: ServicesContent;
}

export default function ServicesSection({ content }: ServicesSectionProps) {
  return (
    <section id="services" className="border-t border-black/5 py-24" aria-labelledby="services-title">
      <div className="mx-auto max-w-7xl px-5 lg:px-8">
        <SectionHeading id="services-title" heading={content.heading} align="center" tone="default" className="mx-auto mb-16 max-w-2xl" />
        <div className="grid gap-x-12 gap-y-14 md:grid-cols-2 lg:grid-cols-3">
          {content.items.map((service) => <article key={service.title} className="relative pl-20"><i className={`${service.icon} absolute left-0 top-1 text-5xl text-[#242424]/40`} /><h3 className="mb-3 text-xl font-semibold text-[#242424]">{service.title}</h3><p className="leading-8 text-black/65">{service.description}</p></article>)}
        </div>
      </div>
    </section>
  );
}
