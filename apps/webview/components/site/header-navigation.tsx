"use client";

import Link from "next/link";
import { useState } from "react";
import type { BrandContent, LinkContent, NavigationItem } from "@/types/home";

interface HeaderNavigationProps {
  brand: BrandContent;
  navigation: NavigationItem[];
  quoteCta: LinkContent;
}

export default function HeaderNavigation({ brand, navigation, quoteCta }: HeaderNavigationProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [openDropdown, setOpenDropdown] = useState<string | null>(null);

  const toggleDropdown = (name: string) => {
    setOpenDropdown((current) => (current === name ? null : name));
  };

  const closeMenu = () => {
    setIsOpen(false);
    setOpenDropdown(null);
  };

  return (
    <nav className="mx-auto flex max-w-7xl items-center justify-between px-5 py-5 lg:px-8" aria-label="Main navigation">
      <Link href={brand.href} className="text-xl font-semibold tracking-wide">
        {brand.name}<span className="text-[#f75757]">{brand.accent}</span>
      </Link>
      <button type="button" className="text-xl lg:hidden" onClick={() => setIsOpen((current) => !current)} aria-expanded={isOpen} aria-controls="main-menu" aria-label="Toggle navigation">
        <i className="fa fa-bars" />
      </button>
      <div id="main-menu" className={`${isOpen ? "block" : "hidden"} absolute left-0 top-full w-full bg-[#242424] px-5 pb-6 lg:static lg:block lg:w-auto lg:bg-transparent lg:p-0`}>
        <ul className="flex flex-col items-center gap-5 text-sm uppercase lg:flex-row">
          {navigation.map((item) => {
            if (item.kind === "link") {
              return (
                <li key={item.label}>
                  <Link href={item.href} className="text-white transition-colors hover:text-[#f75757]" onClick={closeMenu}>{item.label}</Link>
                </li>
              );
            }

            return (
              <li key={item.label} className="relative text-center">
                <button type="button" className="text-white transition-colors hover:text-[#f75757]" onClick={() => toggleDropdown(item.label)} aria-expanded={openDropdown === item.label}>
                  {item.label} <i className="fa fa-angle-down ml-1" />
                </button>
                {openDropdown === item.label && (
                  <ul className="mt-3 space-y-2 text-xs lg:absolute lg:right-0 lg:w-48 lg:bg-white lg:p-4 lg:text-left lg:shadow-lg">
                    {item.children.map((child) => (
                      <li key={child.label}>
                        <Link href={child.href} className="text-white hover:text-[#f75757] lg:text-[#242424]" onClick={closeMenu}>{child.label}</Link>
                      </li>
                    ))}
                  </ul>
                )}
              </li>
            );
          })}
          <li>
            <Link href={quoteCta.href} className="rounded-full border-2 border-[#f75757] px-6 py-3 text-white transition-colors hover:bg-[#f75757]" onClick={closeMenu}>{quoteCta.label}</Link>
          </li>
        </ul>
      </div>
    </nav>
  );
}
