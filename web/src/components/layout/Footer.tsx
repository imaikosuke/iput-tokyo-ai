export default function Footer() {
  return (
    <footer className="w-full py-6 bg-white">
      <div className="container px-4 md:px-6 mx-auto">
        <div className="border-t border-gray-200 pt-6">
          <p className="text-center text-sm text-gray-500">
            © {new Date().getFullYear()} AIPUT TOKYO. All rights reserved.
          </p>
        </div>
      </div>
    </footer>
  );
}
