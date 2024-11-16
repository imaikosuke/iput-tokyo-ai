// src/app/layout.tsx
import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "東京国際工科専門職大学 Q&A",
  description: "大学に関する質問に AI がお答えします",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ja">
      <body className="min-h-screen bg-background font-sans antialiased">{children}</body>
    </html>
  );
}
