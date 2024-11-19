import Image from "next/image"

export default function About() {
  return (
    <section id="about" className="w-full py-12 md:py-24 lg:py-32">
      <div className="container px-4 md:px-6">
        <h2 className="text-4xl md:text-5xl font-bold text-center mb-12 bg-gradient-to-br from-blue-800 to-blue-300 text-transparent bg-clip-text">
          About
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-8 items-center">
          <div className="relative">
            <Image
              src="/about-section.png"
              alt="AIPUT TOKYOについて"
              width={600}
              height={400}
              className="rounded-lg shadow-lg transform hover:scale-105 transition-transform duration-300"
            />
          </div>
          <div className="space-y-6">
            <h3 className="text-2xl md:text-3xl font-semibold relative inline-block">
              AIPUT TOKYOとは
              <span className="absolute -bottom-1 left-0 w-full h-3 bg-blue-200 opacity-50"></span>
            </h3>
            <p className="text-lg text-gray-700 leading-relaxed">
              IPUT 東京校について質問できるRAGを用いたQ&A ChatBot
            </p>
          </div>
        </div>
      </div>
    </section>
  )
}