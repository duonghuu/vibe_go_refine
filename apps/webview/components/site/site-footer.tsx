import Link from "next/link";
import type { FooterContent, LinkContent } from "@/types/home";

interface SiteFooterProps {
  content: FooterContent;
}

function FooterLink({ link }: { link: LinkContent }) {
  return <Link href={link.href}>{link.label}</Link>;
}

export default function SiteFooter({ content }: SiteFooterProps) {
  return (
    <footer className="py-20" aria-label="Footer">
      <div className="mx-auto max-w-7xl px-5 lg:px-8">
        <div className="grid gap-10 sm:grid-cols-2 lg:grid-cols-4">
          {content.linkGroups.map((group) => (
            <div key={group.title}>
              <h3 className="mb-5 text-lg font-semibold text-[#242424]">{group.title}</h3>
              <ul className="space-y-3 leading-7 text-black/65">
                {group.links.map((link) => <li key={link.label}><FooterLink link={link} /></li>)}
              </ul>
            </div>
          ))}
          <div>
            <h3 className="mb-5 text-lg font-semibold text-[#242424]">{content.subscribe.title}</h3>
            <p className="mb-4 leading-8 text-black/65">{content.subscribe.description}</p>
            <form className="max-w-xs" action="#">
              <input type="text" className="mb-3 h-12 w-full border border-black/5 bg-[#f5f8f9] px-4 outline-none focus:border-[#f75757]" placeholder={content.subscribe.placeholder} aria-label="Subscribe email" />
              <button type="button" className="rounded bg-[#f75757] px-6 py-3 text-xs uppercase text-white hover:bg-[#dd0b0b]">{content.subscribe.buttonLabel}</button>
            </form>
          </div>
          <div>
            <Link href={content.brand.href} className="mb-5 block text-2xl font-semibold tracking-wide">
              {content.brand.name}<span className="text-[#f75757]">{content.brand.accent}</span>
            </Link>
            <h6 className="mb-2 text-sm"><a href={content.email.href}>{content.email.label}</a></h6>
            <a href={content.phone.href} className="text-xl text-[#f75757]">{content.phone.label}</a>
          </div>
        </div>
        <div className="mt-12 grid gap-5 border-t border-black/5 pt-6 text-sm text-black/65 md:grid-cols-3">
          <p>{content.copyright.prefix} <span className="text-[#f75757]">{content.copyright.brand}</span> {content.copyright.by} <a href={content.copyright.link.href} target="_blank" rel="noreferrer">{content.copyright.link.label}</a></p>
          <p>{content.distributedBy.prefix} <a href={content.distributedBy.link.href} target="_blank" rel="noreferrer">{content.distributedBy.link.label}</a></p>
          <div className="flex gap-4 md:justify-end">
            {content.socialLinks.map((socialLink) => <a key={socialLink.href} href={socialLink.href} target="_blank" rel="noreferrer" aria-label={socialLink.ariaLabel}>{socialLink.label}</a>)}
          </div>
        </div>
      </div>
    </footer>
  );
}
