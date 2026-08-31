"use client";

import { useEffect, useState } from "react";
import type { TestimonialItem } from "./home-data";

interface TestimonialCarouselProps {
  items: TestimonialItem[];
}

const TestimonialCarousel = ({ items }: TestimonialCarouselProps) => {
  const [activeIndex, setActiveIndex] = useState(0);
  const [visibleCount, setVisibleCount] = useState(2);
  const maxIndex = Math.max(items.length - visibleCount, 0);

  useEffect(() => {
    const updateVisibleCount = () => setVisibleCount(window.innerWidth < 600 ? 1 : 2);
    updateVisibleCount();
    window.addEventListener("resize", updateVisibleCount);
    return () => window.removeEventListener("resize", updateVisibleCount);
  }, []);

  useEffect(() => {
    const timer = window.setInterval(() => setActiveIndex((current) => current >= maxIndex ? 0 : current + 1), 6000);
    return () => window.clearInterval(timer);
  }, [maxIndex]);

  return <div className="overflow-hidden">
    <div className="flex transition-transform duration-500" style={{ transform: `translateX(-${activeIndex * (100 / visibleCount)}%)` }}>
      {items.map((item, index) => <article key={`${item.name}-${item.quote}-${index}`} className="shrink-0 px-4 md:w-1/2" style={{ width: `${100 / visibleCount}%` }}>
        <div className="relative py-10 pl-16 pr-5"><i className="ti-quote-left absolute left-5 top-5 text-4xl text-[#f75757]" /><p className="mb-7 text-xl italic leading-[1.9] text-[#242424]">{item.quote}</p><h3 className="mb-0 text-lg font-semibold text-[#242424]">{item.name}</h3><p className="text-[#808080]">{item.role}</p></div>
      </article>)}
    </div>
    <div className="mt-4 flex justify-center gap-2">
      {Array.from({ length: maxIndex + 1 }, (_, index) => <button key={index} type="button" onClick={() => setActiveIndex(index)} className={`h-2 w-2 rounded-full ${activeIndex === index ? "bg-[#f75757]" : "bg-[#dedede]"}`} aria-label={`Show testimonial ${index + 1}`} />)}
    </div>
  </div>;
};

export default TestimonialCarousel;
