import { getIntroContent } from "@/lib/home/get-intro-content";
import IntroSection from "./intro-section";

export default async function IntroSectionContainer() {
  const result = await getIntroContent();
  if (result.state === "success" && result.content) {
    return <IntroSection content={result.content} />;
  }

  if (result.state === "empty") {
    return (
      <section className="py-24" aria-labelledby="intro-empty-title" role="status">
        <div className="mx-auto max-w-7xl px-5 lg:px-8">
          <p id="intro-empty-title" className="text-black/65">
            Chưa có nội dung giới thiệu.
          </p>
        </div>
      </section>
    );
  }

  return (
    <section className="py-24" aria-labelledby="intro-error-title" role="alert">
      <div className="mx-auto max-w-7xl px-5 lg:px-8">
        <p id="intro-error-title" className="text-black/65">
          Không thể tải nội dung giới thiệu lúc này.
        </p>
      </div>
    </section>
  );
}
