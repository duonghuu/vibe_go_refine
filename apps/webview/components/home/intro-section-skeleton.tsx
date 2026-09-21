export default function IntroSectionSkeleton() {
  return (
    <section className="animate-pulse py-24" aria-busy="true" aria-label="Đang tải nội dung giới thiệu">
      <div className="mx-auto max-w-7xl px-5 lg:px-8">
        <div className="mb-3 h-5 w-64 rounded bg-black/10" />
        <div className="h-12 max-w-3xl rounded bg-black/10 sm:h-14" />
        <div className="mt-16 grid gap-10 md:grid-cols-3">
          {["intro-skeleton-one", "intro-skeleton-two", "intro-skeleton-three"].map((key) => (
            <div key={key}>
              <div className="h-16 w-16 rounded bg-black/10" />
              <div className="mt-6 mb-3 h-7 w-3/4 rounded bg-black/10" />
              <div className="h-20 rounded bg-black/10" />
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
