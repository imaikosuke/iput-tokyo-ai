// src/features/question/QuestionForm.tsx
"use client";

import { useState } from "react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import * as z from "zod";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Form, FormControl, FormField, FormItem, FormLabel } from "@/components/ui/form";
import { Textarea } from "@/components/ui/textarea";
import ReactMarkdown from "react-markdown";

const formSchema = z.object({
  question: z.string().min(1, "質問を入力してください"),
});

export function QuestionForm() {
  const [answer, setAnswer] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      question: "",
    },
  });

  const onSubmit = async (values: z.infer<typeof formSchema>) => {
    setIsLoading(true);
    setError(null);
    setAnswer("");

    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/query/`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ content: values.question }),
      });

      if (!response.ok) {
        throw new Error(`APIエラー: ${response.status}`);
      }

      const data = await response.text();
      setAnswer(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "予期せぬエラーが発生しました");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
          <FormField
            control={form.control}
            name="question"
            render={({ field }) => (
              <FormItem>
                <FormLabel>質問を入力してください</FormLabel>
                <FormControl>
                  <Textarea
                    placeholder="例: 情報工学科について教えてください"
                    className="min-h-[100px]"
                    {...field}
                  />
                </FormControl>
              </FormItem>
            )}
          />
          <Button type="submit" className="w-full" disabled={isLoading}>
            {isLoading ? "送信中..." : "送信"}
          </Button>
        </form>
      </Form>

      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      {answer && (
        <div className="border rounded-md p-4 bg-gray-50">
          <h2 className="font-medium mb-2">回答:</h2>
          <div className="whitespace-pre-wrap">
            <ReactMarkdown>{answer}</ReactMarkdown>
          </div>
        </div>
      )}

      {isLoading && (
        <div className="flex justify-center items-center space-x-2">
          <div className="animate-spin h-5 w-5 border-2 border-blue-500 rounded-full border-t-transparent" />
          <span>回答を生成中...</span>
        </div>
      )}
    </div>
  );
}
