"use client";

import Link from "next/link";
import { useState } from "react";

const InteractiveHeader = () => {
  const [isOpen, setIsOpen] = useState(false);
  const [openDropdown, setOpenDropdown] = useState<string | null>(null);

  const toggleDropdown = (name: string) => {
    setOpenDropdown((current) => (current === name ? null : name));
  };

  return (
    <header className="relative z-20 bg-[#242424] text-white">
      <div className="border-b border-white/5 bg-[#222328]">
        <div className="mx-auto flex max-w-7xl flex-col items-center justify-between gap-2 px-5 py-3 text-sm text-[#919194] md:flex-row lg:px-8">
          <div className="flex gap-5">
            <a href="https://www.facebook.com/themefisher" target="_blank" rel="noreferrer" aria-label="Facebook"><i className="ti-facebook" /></a>
            <a href="https://twitter.com/themefisher" target="_blank" rel="noreferrer" aria-label="Twitter"><i className="ti-twitter" /></a>
            <a href="https://github.com/themefisher/" target="_blank" rel="noreferrer" aria-label="GitHub"><i className="ti-github" /></a>
          </div>
          <div className="flex flex-col items-center gap-2 sm:flex-row sm:gap-8">
            <a href="tel:+23-345-67890">Call Us : <span className="text-white">+23-345-67890</span></a>
            <a href="mailto:support@gmail.com"><i className="fa fa-envelope mr-2" /><span className="text-white">support@gmail.com</span></a>
          </div>
        </div>
      </div>
      <nav className="mx-auto flex max-w-7xl items-center justify-between px-5 py-5 lg:px-8" aria-label="Main navigation">
        <Link href="/" className="text-xl font-semibold tracking-wide">Mega<span className="text-[#f75757]">kit.</span></Link>
        <button type="button" className="text-xl lg:hidden" onClick={() => setIsOpen((current) => !current)} aria-expanded={isOpen} aria-controls="main-menu" aria-label="Toggle navigation">
          <i className="fa fa-bars" />
        </button>
        <div id="main-menu" className={`${isOpen ? "block" : "hidden"} absolute left-0 top-full w-full bg-[#242424] px-5 pb-6 lg:static lg:block lg:w-auto lg:bg-transparent lg:p-0`}>
          <ul className="flex flex-col items-center gap-5 text-sm uppercase lg:flex-row">
            <li><Link href="/" className="text-white transition-colors hover:text-[#f75757]" onClick={() => setIsOpen(false)}>Home</Link></li>
            <li className="relative text-center">
              <button type="button" className="text-white transition-colors hover:text-[#f75757]" onClick={() => toggleDropdown("about")} aria-expanded={openDropdown === "about"}>About <i className="fa fa-angle-down ml-1" /></button>
              {openDropdown === "about" && <ul className="mt-3 space-y-2 text-xs lg:absolute lg:right-0 lg:w-48 lg:bg-white lg:p-4 lg:text-left lg:shadow-lg"><li><a href="#about" className="text-white hover:text-[#f75757] lg:text-[#242424]">Our company</a></li><li><a href="#quote" className="text-white hover:text-[#f75757] lg:text-[#242424]">Pricing</a></li></ul>}
            </li>
            <li><a href="#services" className="text-white transition-colors hover:text-[#f75757]" onClick={() => setIsOpen(false)}>Services</a></li>
            <li><a href="#portfolio" className="text-white transition-colors hover:text-[#f75757]" onClick={() => setIsOpen(false)}>Portfolio</a></li>
            <li className="relative text-center">
              <button type="button" className="text-white transition-colors hover:text-[#f75757]" onClick={() => toggleDropdown("blog")} aria-expanded={openDropdown === "blog"}>Blog <i className="fa fa-angle-down ml-1" /></button>
              {openDropdown === "blog" && <ul className="mt-3 space-y-2 text-xs lg:absolute lg:right-0 lg:w-48 lg:bg-white lg:p-4 lg:text-left lg:shadow-lg"><li><a href="#latest-news" className="text-white hover:text-[#f75757] lg:text-[#242424]">Blog Grid</a></li><li><a href="#latest-news" className="text-white hover:text-[#f75757] lg:text-[#242424]">Blog with Sidebar</a></li><li><a href="#latest-news" className="text-white hover:text-[#f75757] lg:text-[#242424]">Blog Single</a></li></ul>}
            </li>
            <li><a href="#contact" className="text-white transition-colors hover:text-[#f75757]" onClick={() => setIsOpen(false)}>Contact</a></li>
            <li><a href="#contact" className="rounded-full border-2 border-[#f75757] px-6 py-3 text-white transition-colors hover:bg-[#f75757]" onClick={() => setIsOpen(false)}>Get a Quote</a></li>
          </ul>
        </div>
      </nav>
    </header>
  );
};

export default InteractiveHeader;
