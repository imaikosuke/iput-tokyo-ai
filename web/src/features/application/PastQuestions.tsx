// web/src/features/application/PastQuestions.tsx
"use client";

import { motion } from "framer-motion";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { History } from "lucide-react";

interface QnA {
  question: string;
  answer: string;
}

interface PastQuestionsProps {
  pastQnAs: QnA[];
  onQuestionClick: (qna: QnA) => void;
}

const truncateText = (text: string, maxLength: number) => {
  if (text.length <= maxLength) return text;
  return text.slice(0, maxLength) + "・・・";
};

export function PastQuestions({ pastQnAs, onQuestionClick }: PastQuestionsProps) {
  return (
    <Card className="shadow-lg hover:shadow-xl transition-shadow duration-300">
      <CardContent className="p-6">
        <motion.div
          className="space-y-4"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.5, delay: 0.4 }}
        >
          <h2 className="text-xl font-semibold text-center text-gray-700 flex items-center justify-center">
            <History className="mr-2 h-6 w-6" />
            回答履歴
          </h2>
          {pastQnAs.length === 0 ? (
            <div className="text-center text-gray-500">
              <p>過去にしたQ&Aを確認できます</p>
              <p className="text-sm mt-2">※このページを閉じるとリセットされます</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 gap-2">
              {pastQnAs.map((qna, index) => (
                <motion.div
                  key={index}
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ duration: 0.3, delay: index * 0.1 }}
                >
                  <Button
                    variant="outline"
                    className="w-full text-left justify-start h-auto py-2 px-4 shadow-sm hover:shadow-md transition-shadow duration-300"
                    onClick={() => onQuestionClick(qna)}
                  >
                    <span className="text-sm">{truncateText(qna.question, 30)}</span>
                  </Button>
                </motion.div>
              ))}
            </div>
          )}
        </motion.div>
      </CardContent>
    </Card>
  );
}
