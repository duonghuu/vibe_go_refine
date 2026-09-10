import type { HeaderContent } from "@/types/home";

interface HeaderTopBarProps {
  socialLinks: HeaderContent["socialLinks"];
  phoneLabel: HeaderContent["phoneLabel"];
  phone: HeaderContent["phone"];
  email: HeaderContent["email"];
}

export default function HeaderTopBar({ socialLinks, phoneLabel, phone, email }: HeaderTopBarProps) {
  return (
    <div className="border-b border-white/5 bg-[#222328]">
      <div className="mx-auto flex max-w-7xl flex-col items-center justify-between gap-2 px-5 py-3 text-sm text-[#919194] md:flex-row lg:px-8">
        <div className="flex gap-5">
          {socialLinks.map((socialLink) => (
            <a key={socialLink.href} href={socialLink.href} target="_blank" rel="noreferrer" aria-label={socialLink.ariaLabel}>
              <i className={socialLink.icon} />
            </a>
          ))}
        </div>
        <div className="flex flex-col items-center gap-2 sm:flex-row sm:gap-8">
          <a href={phone.href}>{phoneLabel} <span className="text-white">{phone.label}</span></a>
          <a href={email.href}><i className="fa fa-envelope mr-2" /><span className="text-white">{email.label}</span></a>
        </div>
      </div>
    </div>
  );
}
