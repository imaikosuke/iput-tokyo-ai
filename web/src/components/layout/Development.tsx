import Image from "next/image";

export default function Development() {
  return (
    <section id="development" className="w-full py-12 md:py-24 lg:py-32">
      <div className="container px-4 md:px-6">
        <h2 className="text-4xl md:text-5xl font-bold text-center mb-12 bg-gradient-to-br from-blue-800 to-blue-300 text-transparent bg-clip-text">
          Development
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-8 items-center">
          <div className="relative">
            <Image
              src="/development-section.png"
              alt="AIPUT TOKYOの開発者のアイコン"
              width={600}
              height={400}
              className="rounded-lg shadow-lg transform hover:scale-105 transition-transform duration-300"
            />
          </div>
          <div className="space-y-6">
            <h3 className="text-2xl md:text-3xl font-semibold relative inline-block">
              開発元
              <span className="absolute -bottom-1 left-0 w-full h-3 bg-blue-200 opacity-50"></span>
            </h3>
            <p className="text-lg text-gray-700 leading-relaxed">
              IPUTの学生１人で開発・運営をしています。私の高校生の時の実体験から役に立てればと思ったのがきっかけであり、モチベーションにもなっています。
            </p>
          </div>
        </div>
      </div>
    </section>
  );
}
