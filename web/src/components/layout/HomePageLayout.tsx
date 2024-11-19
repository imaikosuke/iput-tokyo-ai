// src/components/layout/HomePageLayout.tsx

import Footer from "./Footer";
import Header from "./Header";

export function HomePageLayout({ children }: { children: React.ReactNode }) {
  return (
    <main className="flex min-h-screen w-full flex-col items-center justify-center gap-8 p-4 md:p-8">
      <Header />
      <div className="mx-auto flex w-full max-w-4xl flex-1 flex-col items-start justify-center pt-36">
        {children}
      </div>
      <Footer />
    </main>
  );
}
