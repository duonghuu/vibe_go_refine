export default function Loading() {
  return (
    <section
      className="relative animate-pulse bg-[#242424] py-48"
      aria-busy="true"
      aria-label="Đang tải trang chủ"
    >
      <div className="mx-auto max-w-7xl px-5 lg:px-8">
        <div className="mb-3 h-5 w-48 rounded bg-white/20" />
        <div className="mb-3 h-14 max-w-3xl rounded bg-white/20 sm:h-20" />
        <div className="mb-10 h-14 max-w-3xl rounded bg-white/20 sm:h-20" />
        <div className="h-12 w-40 rounded-full bg-white/20" />
      </div>
    </section>
  );
}
