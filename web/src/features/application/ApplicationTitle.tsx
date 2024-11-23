// src/features/application/ApplicationTitle.tsx
"use client";
import { motion } from "framer-motion";

export default function ApplicationTitle() {
  return (
    <motion.h1
      initial={{ opacity: 0, y: -20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5 }}
      className="text-4xl font-bold text-center text-blue-800 mb-8"
    >
      IPUT TOKYO AI
    </motion.h1>
  );
}
