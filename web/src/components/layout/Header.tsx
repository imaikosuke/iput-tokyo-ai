// src/components/layout/Header.tsx
"use client";

import Image from "next/image";
import Link from "next/link";
import { Button } from "@/components/ui/button";

export default function Header() {
  const scrollToSection = (sectionId: string) => {
    const section = document.getElementById(sectionId);
    const header = document.querySelector("header");

    if (section && header) {
      const headerHeight = header.offsetHeight;
      const sectionTop = section.offsetTop - headerHeight;

      window.scrollTo({
        top: sectionTop,
        behavior: "smooth",
      });
    }
  };

  return (
    <header className="fixed top-0 left-0 right-0 z-50">
      <div className="mx-auto px-4 sm:px-6 lg:px-8">
        <div className="relative backdrop-blur-md bg-white/40 rounded-full m-4 shadow-lg border border-white/50">
          <div className="flex items-center justify-between h-16 px-4">
            {/* Logo and Title */}
            <div className="flex items-center space-x-3">
              <Image
                src="/iput-logo.png"
                alt="IPUT TOKYO LOGO"
                width={32}
                height={32}
                className="rounded-full"
              />
              <span className="text-xl font-semibold">AIPUT TOKYO</span>
            </div>

            {/* Navigation Links */}
            <nav className="hidden md:flex space-x-16">
              <button
                onClick={() => scrollToSection("about")}
                className="text-lg font-medium text-gray-700 hover:text-gray-900"
              >
                About
              </button>
              <button
                onClick={() => scrollToSection("problem")}
                className="text-lg font-medium text-gray-700 hover:text-gray-900"
              >
                Problem
              </button>
              <button
                onClick={() => scrollToSection("development")}
                className="text-lg font-medium text-gray-700 hover:text-gray-900"
              >
                Development
              </button>
            </nav>

            {/* Let's Start Button */}
            <div>
              <Button
                variant="default"
                className="rounded-full shadow-md hover:shadow-lg transition-shadow duration-300"
              >
                <Link href="/application">Let&apos;s Start</Link>
              </Button>
            </div>
          </div>
        </div>
      </div>
    </header>
  );
}
