import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Megakit | Creative Agency",
  description: "Megakit creative agency and business solutions.",
};

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return <html lang="en" suppressHydrationWarning><head><link rel="stylesheet" href="/plugins/themify/css/themify-icons.css" /><link rel="stylesheet" href="/plugins/fontawesome/css/all.css" /><link rel="preconnect" href="https://fonts.googleapis.com" /><link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Hind:300,400,500,600,700|Poppins:300,400,600,700" /></head><body>{children}</body></html>;
}
