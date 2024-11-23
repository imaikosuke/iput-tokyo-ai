// src/app/application/page.tsx
import ApplicationLayout from "@/features/application/ApplicationLayout";
import ApplicationTitle from "@/features/application/ApplicationTitle";

export default function ApplicationPage() {
  return (
    <main className="flex min-h-screen bg-gradient-to-b from-blue-200 to-white p-8">
      <div className="w-full max-w-6xl mx-auto flex flex-col">
        <ApplicationTitle />
        <ApplicationLayout />
      </div>
    </main>
  );
}
