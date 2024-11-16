import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "standalone",
  env: {
    NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL,
  },
  generateBuildId: async () => {
    return process.env.BUILD_ID || new Date().getTime().toString();
  },
  // async rewrites() {
  //   return [
  //     {
  //       source: "/api/:path*",
  //       destination: "http://server:9020/:path*"
  //     }
  //   ]
  // }
};

export default nextConfig;
