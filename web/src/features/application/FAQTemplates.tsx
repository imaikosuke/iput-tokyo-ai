'use client'

import { useState, useEffect, useRef } from 'react'
import { motion } from 'framer-motion'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { MessageCircleQuestionIcon as QuestionMarkCircle, Building, DiffIcon, Lightbulb, Users } from 'lucide-react'

interface FAQTemplatesProps {
  onQuestionSelect: (question: string) => void
  isAnswerDisplayed: boolean
}

const faqQuestions = [
  { question: "東京国際工科専門職大学とは？", icon: Building },
  { question: "一般的な大学や専門学校と何が違いますか？", icon: DiffIcon },
  { question: "特徴はなんですか？", icon: Lightbulb },
  { question: "どんな人に向いていますか？", icon: Users }
]

export function FAQTemplates({ onQuestionSelect, isAnswerDisplayed }: FAQTemplatesProps) {
  const [isSingleColumn, setIsSingleColumn] = useState(false)
  const containerRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const checkOverflow = () => {
      if (containerRef.current) {
        const buttons = containerRef.current.querySelectorAll('button')
        let hasOverflow = false

        buttons.forEach((button) => {
          if (button.scrollWidth > button.clientWidth) {
            hasOverflow = true
          }
        })

        setIsSingleColumn(hasOverflow || isAnswerDisplayed)
      }
    }

    checkOverflow()
    window.addEventListener('resize', checkOverflow)

    return () => {
      window.removeEventListener('resize', checkOverflow)
    }
  }, [isAnswerDisplayed])

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
            <QuestionMarkCircle className="mr-2 h-6 w-6" />
            よくある質問
          </h2>
          <div
            ref={containerRef}
            className={`grid gap-2 ${
              isSingleColumn ? 'grid-cols-1' : 'grid-cols-1 md:grid-cols-2'
            }`}
          >
            {faqQuestions.map(({ question, icon: Icon }, index) => (
              <motion.div
                key={index}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.3, delay: index * 0.1 }}
              >
                <Button
                  variant="outline"
                  className="w-full text-left justify-start h-auto py-2 px-4 shadow-sm hover:shadow-md transition-shadow duration-300"
                  onClick={() => onQuestionSelect(question)}
                >
                  <Icon className="mr-2 h-4 w-4 flex-shrink-0" />
                  <span className="text-sm">{question}</span>
                </Button>
              </motion.div>
            ))}
          </div>
        </motion.div>
      </CardContent>
    </Card>
  )
}

