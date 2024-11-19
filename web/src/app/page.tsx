// src/app/page.tsx

import { HomePageLayout } from "@/components/layout/HomePageLayout";
import Hero from "@/components/layout/Hero";
import About from "@/components/layout/About";
import Problem from "@/components/layout/Problem";
import Development from "@/components/layout/Development";

export default function Home() {
  return (
    <HomePageLayout>
      <div className="grid gap-8">
        <Hero />
        <About />
        <Problem />
        <Development />
      </div>
    </HomePageLayout>
  );
}
