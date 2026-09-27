'use client'

import { useMemo } from 'react'
import KatexContent from './KatexContent'
import { parseChoices } from '../lib/choices'

interface Props {
  question: string
  value: string
  onPick: (optionText: string) => void
  disabled?: boolean
}

// Tappable options for `multiple_choice` questions that carry an "A) …" list.
// Renders nothing (caller keeps the free-text input) when there are no options.
export default function ChoiceOptions({ question, value, onPick, disabled }: Props) {
  const options = useMemo(() => parseChoices(question), [question])
  if (!options) return null
  return (
    <div
      role="group"
      aria-label="Answer choices"
      className="mt-2 grid grid-cols-1 sm:grid-cols-2 gap-2 min-w-0"
    >
      {options.map(o => {
        const selected = value.trim().toLowerCase() === o.text.trim().toLowerCase()
        return (
          <button
            key={o.letter}
            type="button"
            onClick={() => onPick(o.text)}
            disabled={disabled}
            aria-pressed={selected}
            className={`flex items-start gap-2 text-left px-3 py-2 min-h-[44px] font-mono text-[12px] border rounded-none disabled:opacity-50 ${
              selected
                ? 'border-mathua-blue text-mathua-blue'
                : 'border-mathua-border-strong text-mathua-primary hover:border-mathua-blue hover:text-mathua-blue'
            }`}
          >
            <span className="shrink-0 font-semibold">{o.letter}</span>
            <span className="min-w-0 break-words">
              <KatexContent>{o.text}</KatexContent>
            </span>
          </button>
        )
      })}
    </div>
  )
}
