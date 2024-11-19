// src/app/application/page.tsx
import { QuestionForm } from "@/features/question/QuestionForm";

export default function Application() {
  return (
    <main className="container mx-auto p-4 max-w-3xl">
      <h1 className="text-2xl font-bold mb-6">東京国際工科専門職大学 Q&A</h1>
      <QuestionForm />
    </main>
  );
}
