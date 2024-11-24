// web/src/types/faq.ts
import { Building, DiffIcon, Lightbulb, Users } from "lucide-react";
import type { LucideIcon } from "lucide-react";

export interface FAQItem {
  question: string;
  answer: string;
  icon: LucideIcon;
}

export const faqData: FAQItem[] = [
  {
    question: "東京国際工科専門職大学とは？",
    answer:
      "東京国際工科専門職大学は、2020年に設立された専門職大学です。情報工学とデジタルエンタテインメントの分野で、実践的な技術と理論を学ぶことができる教育機関です。「Designer in Society（社会とともにあるデザイナー）」という教育理念のもと、技術力と創造力を兼ね備えた人材の育成を目指しています。",
    icon: Building,
  },
  {
    question: "一般的な大学や専門学校と何が違いますか？",
    answer:
      "専門職大学である本学の特徴は、理論と実践を組み合わせた教育にあります。一般の大学と比べて実習や企業との連携が多く、専門学校と比べて理論的な学びも充実しています。また、卒業時に「学士（専門職）」の学位が授与され、大学院進学も可能です。",
    icon: DiffIcon,
  },
  {
    question: "特徴はなんですか？",
    answer:
      "本学の主な特徴は以下の3つです:\n\n1. 第一線で活躍する専門家による実践的な教育\n2. 企業や地域社会との密接な連携による実践的なプロジェクト学習\n3. 最新の設備と少人数制による手厚い指導\n\nまた、1年次から実習・演習を多く取り入れ、実践力を段階的に身につけていく教育課程も特徴です。",
    icon: Lightbulb,
  },
  {
    question: "どんな人に向いていますか？",
    answer:
      "以下のような方に特に向いています：\n\n1. 技術力と創造力を活かしたい方\n2. 実践的な学びを通じて即戦力となりたい方\n3. デジタル技術を使って社会に貢献したい方\n4. 理論と実践の両方を学びたい方\n\n特に、技術だけでなく、その活用方法や社会での役割についても深く考えたい方に適しています。",
    icon: Users,
  },
];
