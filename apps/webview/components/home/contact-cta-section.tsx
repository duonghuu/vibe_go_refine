import type { ContactCtaContent } from "@/types/home";

interface ContactCtaSectionProps {
  content: ContactCtaContent;
}

export default function ContactCtaSection({ content }: ContactCtaSectionProps) {
  return (
    <section id="contact" style={{ backgroundImage: `url(${content.backgroundImage})` }} className="relative bg-cover bg-center py-28 before:absolute before:inset-0 before:bg-black/50" aria-labelledby="cta-title">
      <div className="relative mx-auto max-w-7xl px-5 lg:px-8">
        <div className="max-w-lg rounded bg-white p-10">
          <span className="text-sm uppercase tracking-wide text-[#f75757]">{content.heading.eyebrow}</span>
          <h2 id="cta-title" className="mt-2 mb-5 text-3xl font-semibold leading-tight text-[#242424]">{content.heading.title}</h2>
          <p className="mb-4 text-xl text-black/65">{content.description}</p>
          <h3 className="text-2xl font-semibold text-[#242424]"><i className={`${content.phoneIcon} mr-3 text-[#f75757]`} />{content.phone.label}</h3>
        </div>
      </div>
    </section>
  );
}
