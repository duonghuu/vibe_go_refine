import type { HeaderContent } from "@/types/home";
import HeaderNavigation from "./header-navigation";
import HeaderTopBar from "./header-top-bar";

interface SiteHeaderProps {
  content: HeaderContent;
}

export default function SiteHeader({ content }: SiteHeaderProps) {
  return (
    <header className="relative z-20 bg-[#242424] text-white">
      <HeaderTopBar socialLinks={content.socialLinks} phoneLabel={content.phoneLabel} phone={content.phone} email={content.email} />
      <HeaderNavigation brand={content.brand} navigation={content.navigation} quoteCta={content.quoteCta} />
    </header>
  );
}
