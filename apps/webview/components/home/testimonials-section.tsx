import type { TestimonialsContent } from "@/types/home";
import SectionHeading from "./section-heading";
import TestimonialCarousel from "./testimonial-carousel";

interface TestimonialsSectionProps {
  content: TestimonialsContent;
}

export default function TestimonialsSection({ content }: TestimonialsSectionProps) {
  return (
    <section className="py-24" aria-labelledby="testimonial-title">
      <div className="mx-auto max-w-7xl px-5 lg:px-8">
        <SectionHeading id="testimonial-title" heading={content.heading} align="left" tone="default" className="mb-10 max-w-2xl" />
        <TestimonialCarousel items={content.items} />
      </div>
    </section>
  );
}
