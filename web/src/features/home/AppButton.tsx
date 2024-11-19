import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import Link from "next/link";

export default function AppButton({ className }: { className?: React.ReactNode }) {
  return (
    <Button asChild className={cn("rounded-full", className)}>
      <Link href="/application">Let&apos;s Start</Link>
    </Button>
  );
}
