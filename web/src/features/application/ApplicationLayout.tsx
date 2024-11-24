// web/src/features/application/ApplicationLayout.tsx
"use client";

import { useState } from "react";
import { motion } from "framer-motion";
import { QuestionForm } from "./QuestionForm";
import { FAQTemplates } from "./FAQTemplates";
import { AnswerDisplay } from "./AnswerDisplay";
import { PastQuestions } from "./PastQuestions";

interface QnA {
  question: string;
  answer: string;
}

interface APIResponse {
  answer: string;
}

export default function ApplicationLayout() {
  const [answer, setAnswer] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [pastQnAs, setPastQnAs] = useState<QnA[]>([]);

  const handleSubmit = async (question: string, predefinedAnswer?: string) => {
    setIsLoading(true);
    setAnswer(null);

    // FAQの場合は事前に用意した回答を表示
    if (predefinedAnswer) {
      setAnswer(predefinedAnswer);
      setPastQnAs((prev) => [...prev, { question, answer: predefinedAnswer }]);
      setIsLoading(false);
      return;
    }

    try {
      // 通常の質問の場合はAPI呼び出し
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

      const data = (await response.json()) as APIResponse;
      setAnswer(data.answer);
      setPastQnAs((prev) => [...prev, { question, answer: data.answer }]);
    } catch (err) {
      console.error(err);
      setAnswer("エラーが発生しました。もう一度お試しください。");
    } finally {
      setIsLoading(false);
    }
  };

  const handlePastQuestionClick = (qna: QnA) => {
    setAnswer(qna.answer);
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
        <PastQuestions pastQnAs={pastQnAs} onQuestionClick={handlePastQuestionClick} />
      </div>
      <AnswerDisplay answer={answer} isLoading={isLoading} />
    </motion.div>
  );
}
