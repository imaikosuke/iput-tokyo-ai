import Link from "next/link";
import * as React from "react";

// import {
//   ContextMenu,
//   ContextMenuContent,
//   ContextMenuItem,
//   ContextMenuTrigger,
// } from "@openstatus/ui/src/components/context-menu";
import Image from "next/image";

export function AppLogo() {
  return (
    // <ContextMenu>
    //   <ContextMenuTrigger>
    <div>
      <Link href="/" className="flex items-center gap-2 font-cal">
        <Image
          src="/iput-tokyo.png"
          alt="IPUT TOKYO LOGO"
          height={30}
          width={30}
          className="rounded-full border border-border bg-transparent"
        />
        AIPUT TOKYO
      </Link>
      {/* </ContextMenuTrigger> */}
      {/* <ContextMenuContent> */}
      {/* <ContextMenuItem asChild> */}
      <a href="/iput-logo.png" download="iput-logo.png">
        Download SVG
      </a>
      {/* </ContextMenuItem> */}
      {/* </ContextMenuContent> */}
    </div>
    // </ContextMenu>
  );
}
