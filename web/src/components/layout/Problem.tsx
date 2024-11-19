import Image from "next/image";

export default function Problem() {
  return (
    <section id="problem" className="w-full py-12 md:py-24 lg:py-32">
      <div className="container px-4 md:px-6">
        <h2 className="text-4xl md:text-5xl font-bold text-center mb-12 bg-gradient-to-br from-blue-800 to-blue-300 text-transparent bg-clip-text">
          Problem
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-8 items-center">
          <div className="space-y-6 order-2 md:order-1">
            <h3 className="text-2xl md:text-3xl font-semibold relative inline-block">
              解決したい課題
              <span className="absolute -bottom-1 left-0 w-full h-3 bg-blue-200 opacity-50"></span>
            </h3>
            <p className="text-lg text-gray-700 leading-relaxed">
              IPUTに興味があっても情報源が少なくて調べることができずに進学しても大丈夫か不安になる
            </p>
          </div>
          <div className="relative order-1 md:order-2">
            <Image
              src="/problem-section.png"
              alt="AIPUT TOKYOによって解決したい課題"
              width={600}
              height={400}
              className="rounded-lg shadow-lg transform hover:scale-105 transition-transform duration-300"
            />
          </div>
        </div>
      </div>
    </section>
  );
}
