interface FormulaBlockProps {
  code: string
  label?: string
  className?: string
}

export default function FormulaBlock({ code, label, className = '' }: FormulaBlockProps) {
  return (
    <div className={`flex flex-col items-center ${className}`}>
      <pre className="bg-mathua-code border border-mathua-border border-l-[2px] border-l-mathua-blue rounded-r-md p-5 font-mono text-xs text-mathua-blue whitespace-pre overflow-x-auto leading-relaxed text-center inline-block">
        {code}
      </pre>
      {label && (
        <span className="font-mono text-[11px] text-mathua-muted mt-1">{label}</span>
      )}
    </div>
  )
}
