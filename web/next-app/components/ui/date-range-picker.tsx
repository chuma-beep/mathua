'use client'

import * as React from 'react'
import type { DateRange } from 'react-day-picker'
import { Button } from '@/components/ui/button'
import { Calendar } from '@/components/ui/calendar'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'

function label(value?: Date) {
  return value ? value.toLocaleDateString('en-GB', { day: 'numeric', month: 'short', year: 'numeric', timeZone: 'UTC' }) : ''
}

// The picker's contract is a UTC midnight for the day the operator chose. The
// calendar works in local days, so both directions are converted here rather
// than at each call site. A picked day therefore never drifts when the browser
// is east or west of UTC.
export function dateKey(value: Date): string {
  return value.toISOString().slice(0, 10)
}

function fromCalendarDay(value: Date): Date {
  return new Date(Date.UTC(value.getFullYear(), value.getMonth(), value.getDate()))
}

function toCalendarDay(value: Date): Date {
  return new Date(value.getUTCFullYear(), value.getUTCMonth(), value.getUTCDate())
}

export function DateRangePicker({
  value,
  onChange,
  min,
  max,
  align = 'start',
}: {
  value: DateRange | undefined
  onChange: (range: DateRange | undefined) => void
  min?: Date
  max?: Date
  align?: 'start' | 'center' | 'end'
}) {
  const [open, setOpen] = React.useState(false)
  const selected = React.useMemo(
    () => ({
      from: value?.from ? toCalendarDay(value.from) : undefined,
      to: value?.to ? toCalendarDay(value.to) : undefined,
    }),
    [value],
  )
  const select = React.useCallback(
    (next: DateRange | undefined) => {
      onChange(next?.from ? { from: fromCalendarDay(next.from), to: next.to ? fromCalendarDay(next.to) : undefined } : undefined)
    },
    [onChange],
  )
  const text = value?.from
    ? value.to
      ? `${label(value.from)} – ${label(value.to)}`
      : label(value.from)
    : 'Select a date range'
  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button variant="outline" className="h-9 min-w-[15rem] justify-start font-mono text-xs">
          {text}
        </Button>
      </PopoverTrigger>
      <PopoverContent align={align} className="w-auto p-0">
        <Calendar
          mode="range"
          selected={selected}
          onSelect={select}
          numberOfMonths={2}
          disabled={[{ before: min ?? new Date('2026-10-19T00:00:00Z') }, ...(max ? [{ after: max }] : [])]}
        />
        <div className="flex items-center justify-between border-t border-mathua-border px-3 py-2">
          <button type="button" className="font-mono text-[10px] text-mathua-muted hover:text-mathua-blue" onClick={() => onChange(undefined)}>
            Clear
          </button>
          <button type="button" className="font-mono text-[10px] text-mathua-blue hover:underline" onClick={() => setOpen(false)}>
            Done
          </button>
        </div>
      </PopoverContent>
    </Popover>
  )
}
