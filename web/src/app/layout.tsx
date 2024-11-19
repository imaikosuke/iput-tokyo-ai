// src/app/layout.tsx
import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "AIPUT TOKYO",
  description: "大学に関する質問に AI がお答えします",
  icons: {
    icon: [
      {
        url: "/iput-logo.png",
        type: "image/png",
      }
    ]
  }
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
