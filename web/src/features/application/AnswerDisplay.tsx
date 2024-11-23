// web/src/features/application/AnswerDisplay.tsx
"use client";

import { motion } from "framer-motion";
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import ReactMarkdown from "react-markdown";

interface AnswerDisplayProps {
  answer: string | null;
  isLoading: boolean;
}

export function AnswerDisplay({ answer, isLoading }: AnswerDisplayProps) {
  return (
    <motion.div
      initial={{ opacity: 0, x: 20 }}
      animate={{ opacity: 1, x: 0 }}
      transition={{ duration: 0.5 }}
      className="flex-1 md:w-1/2"
    >
      <Card className="h-full shadow-lg hover:shadow-xl transition-shadow duration-300">
        <CardContent className="p-6">
          <h2 className="text-2xl font-semibold mb-4 text-blue-700">回答</h2>
          {isLoading ? (
            <AnswerSkeleton />
          ) : answer ? (
            <p className="text-lg">
              <ReactMarkdown>{answer}</ReactMarkdown>
            </p>
          ) : (
            <p className="text-lg text-gray-500">質問を入力してください。回答がここに表示されます。</p>
          )}
        </CardContent>
      </Card>
    </motion.div>
  );
}

function AnswerSkeleton() {
  return (
    <div className="space-y-2">
      <Skeleton className="h-4 w-full" />
      <Skeleton className="h-4 w-full" />
      <Skeleton className="h-4 w-full" />
      <Skeleton className="h-4 w-3/4" />
    </div>
  );
}
