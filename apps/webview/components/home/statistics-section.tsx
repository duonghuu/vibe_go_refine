import type { StatItem } from "@/types/home";
import AnimatedCounter from "./animated-counter";

interface StatisticsSectionProps {
  items: StatItem[];
}

export default function StatisticsSection({ items }: StatisticsSectionProps) {
  return (
    <section className="py-24" aria-label="Company statistics">
      <div className="mx-auto grid max-w-7xl grid-cols-2 gap-10 px-5 md:grid-cols-4 lg:px-8">
        {items.map((stat) => <div key={stat.label} className="text-center"><h3 className="mb-0 text-3xl font-semibold text-[#242424] md:text-5xl"><AnimatedCounter value={stat.value} suffix={stat.suffix} /></h3><p className="text-[#808080]">{stat.label}</p></div>)}
      </div>
    </section>
  );
}
