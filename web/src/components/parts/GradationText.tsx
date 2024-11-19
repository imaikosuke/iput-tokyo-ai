interface GradationTextProps {
  children: React.ReactNode;
  className?: string;
}

export default function GradationText({ children, className = "" }: GradationTextProps) {
  return (
    <span
      className={`
        inline-block 
        bg-gradient-to-br 
        from-gray-900 
        to-gray-400 
        bg-clip-text 
        text-transparent 
        font-bold 
        ${className}
      `}
    >
      {children}
    </span>
  );
}
