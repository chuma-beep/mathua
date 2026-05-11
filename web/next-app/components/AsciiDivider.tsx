interface AsciiDividerProps {
  pattern?: 'double' | 'dash' | 'wave'
  className?: string
}

const patterns: Record<string, string> = {
  double: '══════════════════════════════════════════════════════════',
  dash: '─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─',
  wave: '╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌',
}

export default function AsciiDivider({ pattern = 'double', className = '' }: AsciiDividerProps) {
  return (
    <div className={`ascii-divider py-4 animate-ascii-pulse ${className}`} aria-hidden="true">
      {patterns[pattern]}
    </div>
  )
}
