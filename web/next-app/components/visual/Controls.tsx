'use client'

import { useId } from 'react'

interface SliderProps {
  label: string
  value: number
  min: number
  max: number
  step?: number
  onChange: (value: number) => void
  format?: (value: number) => string
}

/**
 * A labelled range input with a visible value. Native range inputs already
 * answer arrow keys, Home and End, which is the keyboard path every drag has.
 */
export function Slider({ label, value, min, max, step = 0.1, onChange, format }: SliderProps) {
  const id = useId()
  return (
    <div className="flex items-center gap-2">
      <label htmlFor={id} className="w-14 shrink-0 font-mono text-[10px] uppercase tracking-wide text-mathua-muted">
        {label}
      </label>
      <input
        id={id}
        type="range"
        min={min}
        max={max}
        step={step}
        value={value}
        onChange={(e) => onChange(Number(e.target.value))}
        className="min-w-0 flex-1"
        style={{ accentColor: 'var(--accent-blue)' }}
      />
      <output htmlFor={id} className="w-14 shrink-0 text-right font-mono text-[11px] tabular-nums text-mathua-primary">
        {format ? format(value) : String(value)}
      </output>
    </div>
  )
}

interface NumberFieldProps {
  label: string
  value: number
  min?: number
  max?: number
  step?: number
  onChange: (value: number) => void
}

/** A compact numeric entry, used where a slider would hide the exact value. */
export function NumberField({ label, value, min = -9, max = 9, step = 0.1, onChange }: NumberFieldProps) {
  const id = useId()
  return (
    <div className="flex items-center gap-1.5">
      <label htmlFor={id} className="font-mono text-[10px] uppercase tracking-wide text-mathua-muted">
        {label}
      </label>
      <input
        id={id}
        type="number"
        value={value}
        min={min}
        max={max}
        step={step}
        onChange={(e) => {
          const next = Number(e.target.value)
          if (Number.isFinite(next)) onChange(next)
        }}
        className="w-16 border border-mathua-border bg-mathua-bg px-1.5 py-0.5 text-right font-mono text-[11px] tabular-nums text-mathua-primary"
      />
    </div>
  )
}

interface SelectProps<T extends string> {
  label: string
  value: T
  options: Array<{ value: T; label: string }>
  onChange: (value: T) => void
}

export function Select<T extends string>({ label, value, options, onChange }: SelectProps<T>) {
  const id = useId()
  return (
    <div className="flex items-center gap-2">
      <label htmlFor={id} className="font-mono text-[10px] uppercase tracking-wide text-mathua-muted">
        {label}
      </label>
      <select
        id={id}
        value={value}
        onChange={(e) => onChange(e.target.value as T)}
        className="border border-mathua-border bg-mathua-bg px-1.5 py-0.5 font-mono text-[11px] text-mathua-primary"
      >
        {options.map((option) => (
          <option key={option.value} value={option.value}>
            {option.label}
          </option>
        ))}
      </select>
    </div>
  )
}

interface ToggleProps {
  label: string
  checked: boolean
  onChange: (checked: boolean) => void
}

export function Toggle({ label, checked, onChange }: ToggleProps) {
  const id = useId()
  return (
    <label htmlFor={id} className="flex items-center gap-2 font-mono text-[10px] text-mathua-secondary">
      <input
        id={id}
        type="checkbox"
        checked={checked}
        onChange={(e) => onChange(e.target.checked)}
        style={{ accentColor: 'var(--accent-blue)' }}
      />
      {label}
    </label>
  )
}
