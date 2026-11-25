'use client'

import { useId } from 'react'
import { Input } from '@/components/ui/input'
import { Switch } from '@/components/ui/switch'
import { Select as UiSelect, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'

// Form controls for the visual figures, built on the shadcn primitives in
// components/ui so a figure's fields are the same fields the rest of Mathua
// uses. Nothing here is a hand-rolled control.
//
// Two details are deliberate. The native number spinners are hidden and
// replaced by real buttons: the browser's own stepper is a different size,
// colour and hit target in every browser, it cannot be labelled for a screen
// reader, and on a dark theme it is often invisible. And every numeric entry
// fills its cell, because a fixed 64px box clipped a value like `-1.25`.

const numeral =
  'h-10 rounded-none border-mathua-border bg-mathua-code px-2 text-center font-mono text-xs tabular-nums text-mathua-primary focus:border-mathua-blue'

/** Hides the platform stepper in every engine; the buttons take over. */
const noSpinner =
  '[appearance:textfield] [&::-webkit-inner-spin-button]:appearance-none [&::-webkit-outer-spin-button]:appearance-none'

function StepButton({
  label,
  onClick,
  disabled,
  children,
}: {
  label: string
  onClick: () => void
  disabled: boolean
  children: React.ReactNode
}) {
  return (
    <button
      type="button"
      aria-label={label}
      onClick={onClick}
      disabled={disabled}
      className="flex h-10 w-10 shrink-0 items-center justify-center border border-mathua-border bg-mathua-surface font-mono text-[13px] leading-none text-mathua-secondary transition-colors hover:border-mathua-blue hover:text-mathua-blue disabled:cursor-not-allowed disabled:opacity-40"
    >
      {children}
    </button>
  )
}

interface StepperProps {
  id?: string
  label: string
  value: number
  min: number
  max: number
  step: number
  onChange: (value: number) => void
  className?: string
}

/** A number box with explicit decrement and increment buttons. */
function Stepper({ id, label, value, min, max, step, onChange, className }: StepperProps) {
  const commit = (next: number) => {
    if (!Number.isFinite(next)) return
    onChange(Math.min(max, Math.max(min, next)))
  }
  return (
    <div className={`flex items-stretch ${className ?? ''}`}>
      <StepButton label={`Decrease ${label}`} onClick={() => commit(value - step)} disabled={value <= min}>
        −
      </StepButton>
      <Input
        id={id}
        type="number"
        value={Number(value.toFixed(4))}
        min={min}
        max={max}
        step={step}
        aria-label={id ? undefined : `${label} value`}
        onChange={(e) => commit(Number(e.target.value))}
        className={`${numeral} ${noSpinner} w-14 rounded-none border-x-0`}
      />
      <StepButton label={`Increase ${label}`} onClick={() => commit(value + step)} disabled={value >= max}>
        +
      </StepButton>
    </div>
  )
}

/** A stacked label above its control. */
function Field({ label, htmlFor, children }: { label: string; htmlFor?: string; children: React.ReactNode }) {
  return (
    <div className="flex min-w-0 flex-col gap-1">
      <label htmlFor={htmlFor} className="font-mono text-[10px] text-mathua-muted">
        {label}
      </label>
      {children}
    </div>
  )
}

interface SliderProps {
  label: string
  value: number
  min: number
  max: number
  step?: number
  onChange: (value: number) => void
  /** A unit shown after the number box, e.g. a degree sign. */
  suffix?: string
}

/**
 * A range with a stepped number box beside it. The range is the fast path and
 * answers arrow keys, Home and End; the box is the precise one, because a slope
 * of 1.37 is not reachable by dragging. Both write the same clamped value, so
 * they cannot disagree.
 */
export function Slider({ label, value, min, max, step = 0.1, onChange, suffix }: SliderProps) {
  const id = useId()
  const commit = (next: number) => {
    if (!Number.isFinite(next)) return
    onChange(Math.min(max, Math.max(min, next)))
  }
  return (
    // The label always owns its line, and the range wraps beside the stepper.
    // A narrow figure-controls strip cannot hold all three on one row; doing so
    // pushed the stepper's + button past the card edge.
    <div className="space-y-1">
      <label htmlFor={id} className="block font-mono text-[10px] text-mathua-muted">
        {label}
      </label>
      <div className="flex flex-wrap items-center gap-2">
        <input
          id={id}
          type="range"
          min={min}
          max={max}
          step={step}
          value={value}
          onChange={(e) => commit(Number(e.target.value))}
          className="min-w-[8rem] flex-1"
          style={{ accentColor: 'var(--accent-blue)' }}
        />
        <span className="flex shrink-0 items-center gap-2">
          <Stepper label={label} value={value} min={min} max={max} step={step} onChange={onChange} />
          {suffix && <span className="font-mono text-[10px] text-mathua-muted">{suffix}</span>}
        </span>
      </div>
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

/** A numeric entry that fills its cell, for values a slider would only approximate. */
export function NumberField({ label, value, min = -9, max = 9, step = 0.1, onChange }: NumberFieldProps) {
  const id = useId()
  return (
    <Field label={label} htmlFor={id}>
      <Stepper id={id} label={label} value={value} min={min} max={max} step={step} onChange={onChange} className="w-full [&>input]:w-full" />
    </Field>
  )
}

interface SelectProps<T extends string> {
  label: string
  value: T
  options: Array<{ value: T; label: string }>
  onChange: (value: T) => void
}

/** shadcn Select, full width, so a choice reads as a control rather than a label. */
export function Select<T extends string>({ label, value, options, onChange }: SelectProps<T>) {
  const id = useId()
  return (
    <Field label={label} htmlFor={id}>
      <UiSelect value={value} onValueChange={(next) => onChange(next as T)}>
        <SelectTrigger id={id} aria-label={label}>
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          {options.map((option) => (
            <SelectItem key={option.value} value={option.value}>
              {option.label}
            </SelectItem>
          ))}
        </SelectContent>
      </UiSelect>
    </Field>
  )
}

interface ToggleProps {
  label: string
  checked: boolean
  onChange: (checked: boolean) => void
}

/** A labelled switch, the same control the settings pages use for a boolean. */
export function Toggle({ label, checked, onChange }: ToggleProps) {
  const id = useId()
  return (
    <div className="flex min-h-11 items-center justify-between gap-3">
      <label htmlFor={id} className="font-mono text-xs text-mathua-secondary">
        {label}
      </label>
      <Switch id={id} checked={checked} onCheckedChange={onChange} className="min-h-0" />
    </div>
  )
}
