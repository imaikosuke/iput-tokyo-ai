// web/src/features/application/ApplicationLayout.tsx
"use client";

import { useState } from "react";
import { motion } from "framer-motion";
import { QuestionForm } from "./QuestionForm";
import { FAQTemplates } from "./FAQTemplates";
import { AnswerDisplay } from "./AnswerDisplay";

export default function ApplicationLayout() {
  const [answer, setAnswer] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(false);

  const handleSubmit = async (question: string) => {
    setIsLoading(true);
    setAnswer(null);
    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/query/`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ content: question }),
      });

      if (!response.ok) {
        throw new Error(`APIエラー: ${response.status}`);
      }

      const data = await response.text();
      setAnswer(data);
    } catch (err) {
      console.error(err);
      setAnswer("エラーが発生しました。もう一度お試しください。");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5 }}
      className="flex flex-col md:flex-row gap-8"
    >
      <div className="flex-1 space-y-8">
        <QuestionForm onSubmit={handleSubmit} />
        <FAQTemplates onQuestionSelect={handleSubmit} isAnswerDisplayed={true} />
      </div>
      <AnswerDisplay answer={answer} isLoading={isLoading} />
    </motion.div>
  );
}
