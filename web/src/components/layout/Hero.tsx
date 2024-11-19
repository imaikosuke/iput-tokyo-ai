import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { ChevronRight } from "lucide-react";
import Link from "next/link";

export default function Hero() {
  return (
    <section className="w-full py-12 md:py-24 lg:py-32 xl:py-48">
      <div className="container px-4 md:px-6">
        <div className="flex flex-col items-center space-y-8 text-center">
          <Link href="https://techjoruney-code.com/about-me/" target="_blank">
            <Badge variant="outline" className="backdrop-blur-[4px]">
              Produce By imaikosuke
              <ChevronRight className="ml-1 h-4 w-4" />
            </Badge>
          </Link>
          <h1 className="text-3xl font-bold tracking-tighter sm:text-4xl md:text-5xl lg:text-6xl/none">
            A better way to get to know <span className="text-blue-600">IPUT</span>.
          </h1>
          <p className="max-w-[42rem] leading-normal text-muted-foreground sm:text-xl sm:leading-8">
            IPUTについて質問ができるRAGを用いたQ&A ChatBot
          </p>
          <Link href="/application">
            <Button size="lg" className="mt-2">
              Get Started
            </Button>
          </Link>
        </div>
      </div>
    </section>
  );
}
