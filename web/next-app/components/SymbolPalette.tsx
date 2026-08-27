'use client'

const SYMBOLS = [
  'π', '±', '∞', '√', '∛', '²', '³', '°',
  '≠', '≈', '≤', '≥', '×', '÷', '·',
  '∫', '∑', '∂', 'θ', 'α', 'β',
  '(', ')', '[', ']', '|', '/',
]

interface Props {
  targetRef: React.RefObject<HTMLInputElement>
  onInsert: (val: string) => void
  compact?: boolean
}

export default function SymbolPalette({ targetRef, onInsert, compact }: Props) {
  function insert(sym: string) {
    const el = targetRef.current
    if (!el) {
      onInsert(sym)
      return
    }
    const start = el.selectionStart ?? el.value.length
    const end = el.selectionEnd ?? el.value.length
    const before = el.value.slice(0, start)
    const after = el.value.slice(end)
    const next = before + sym + after
    onInsert(next)
    requestAnimationFrame(() => {
      el.focus()
      const pos = start + sym.length
      el.setSelectionRange(pos, pos)
    })
  }

  return (
    <div className={`flex flex-wrap gap-1.5 p-2 bg-mathua-surface border border-mathua-border ${compact ? 'mt-1' : 'mt-2'}`}>
      {SYMBOLS.map((s) => (
        <button
          key={s}
          type="button"
          onClick={() => insert(s)}
          className="w-7 h-7 font-mono text-xs border border-mathua-border-strong bg-mathua-code text-mathua-primary hover:border-mathua-blue hover:text-mathua-blue"
          aria-label={`Insert ${s}`}
        >
          {s}
        </button>
      ))}
    </div>
  )
}
