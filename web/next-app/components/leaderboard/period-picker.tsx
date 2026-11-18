'use client'

import * as React from 'react'
import { Button } from '@/components/ui/button'
import { Calendar } from '@/components/ui/calendar'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuLabel,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuSeparator,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { ChevronDown, ChevronLeft, ChevronRight } from '@/components/icons'
import {
  monthLabel,
  monthStartKey,
  parsePeriodKey,
  toCalendarDay,
  weekLabel,
  weekStartKey,
} from '@/lib/leaderboard-period'

const triggerClass =
  'h-9 justify-start gap-2 rounded-none border-mathua-border bg-mathua-surface px-3 font-mono text-xs text-mathua-primary hover:border-mathua-blue'

export function WeekPicker({
  value,
  onChange,
  min,
  max,
}: {
  value: string
  onChange: (key: string) => void
  min?: string
  max?: string
}) {
  const [open, setOpen] = React.useState(false)
  const selected = parsePeriodKey(value)
  const minDate = min ? parsePeriodKey(min) : undefined
  const maxDate = max ? parsePeriodKey(max) : undefined
  const disabled = [
    ...(minDate ? [{ before: toCalendarDay(minDate) }] : []),
    ...(maxDate ? [{ after: toCalendarDay(maxDate) }] : []),
  ]
  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button variant="outline" className={triggerClass} aria-label="Choose a week">
          {weekLabel(value)}
          <ChevronDown size={14} aria-hidden="true" />
        </Button>
      </PopoverTrigger>
      <PopoverContent align="center" className="w-auto rounded-none border-mathua-border p-0">
        <Calendar
          mode="single"
          selected={selected ? toCalendarDay(selected) : undefined}
          defaultMonth={selected ? toCalendarDay(selected) : maxDate ? toCalendarDay(maxDate) : undefined}
          onSelect={(day) => {
            if (!day) return
            onChange(weekStartKey(day))
            setOpen(false)
          }}
          disabled={disabled.length > 0 ? disabled : undefined}
          className="rounded-none"
        />
        <p className="border-t border-mathua-border px-3 py-2 font-mono text-[10px] text-mathua-muted">
          Any day opens that Monday to Sunday week, UTC.
        </p>
      </PopoverContent>
    </Popover>
  )
}

function monthsBetween(min: string, max: string) {
  const first = parsePeriodKey(monthStartKey(parsePeriodKey(min) ?? new Date()))
  const last = parsePeriodKey(monthStartKey(parsePeriodKey(max) ?? new Date()))
  if (!first || !last) return []
  const out: { key: string; label: string; year: string }[] = []
  const cursor = new Date(last)
  while (cursor >= first) {
    const key = monthStartKey(cursor)
    out.push({ key, label: monthLabel(key), year: String(cursor.getUTCFullYear()) })
    cursor.setUTCMonth(cursor.getUTCMonth() - 1)
  }
  return out
}

export function MonthPicker({
  value,
  onChange,
  min,
  max,
}: {
  value: string
  onChange: (key: string) => void
  min?: string
  max?: string
}) {
  const months = React.useMemo(() => monthsBetween(min ?? value, max ?? value), [min, max, value])
  const currentYear = new Date().getUTCFullYear()
  const recent = months.filter((m) => m.year === String(currentYear))
  const older = months.filter((m) => m.year !== String(currentYear))
  const years = Array.from(new Set(older.map((m) => m.year)))

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" className={triggerClass} aria-label="Choose a month">
          {monthLabel(value)}
          <ChevronDown size={14} aria-hidden="true" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="center" className="max-h-[320px] w-48 overflow-y-auto rounded-none border-mathua-border">
        <DropdownMenuLabel className="font-mono text-[10px] uppercase tracking-wider text-mathua-muted">
          Month
        </DropdownMenuLabel>
        <DropdownMenuSeparator className="bg-mathua-border" />
        <DropdownMenuRadioGroup value={value} onValueChange={onChange}>
          {recent.map((month) => (
            <DropdownMenuRadioItem key={month.key} value={month.key} className="rounded-none font-mono text-xs">
              {month.label}
            </DropdownMenuRadioItem>
          ))}
        </DropdownMenuRadioGroup>
        {years.map((year) => (
          <DropdownMenuSub key={year}>
            <DropdownMenuSubTrigger className="rounded-none font-mono text-xs">{year}</DropdownMenuSubTrigger>
            <DropdownMenuSubContent className="max-h-[280px] w-44 overflow-y-auto rounded-none border-mathua-border">
              <DropdownMenuRadioGroup value={value} onValueChange={onChange}>
                {older
                  .filter((month) => month.year === year)
                  .map((month) => (
                    <DropdownMenuRadioItem key={month.key} value={month.key} className="rounded-none font-mono text-xs">
                      {month.label}
                    </DropdownMenuRadioItem>
                  ))}
              </DropdownMenuRadioGroup>
            </DropdownMenuSubContent>
          </DropdownMenuSub>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}

export function PeriodArrows({
  onPrev,
  onNext,
  prevDisabled,
  nextDisabled,
  label,
}: {
  onPrev: () => void
  onNext: () => void
  prevDisabled?: boolean
  nextDisabled?: boolean
  label: string
}) {
  return (
    <div className="flex items-center gap-1">
      <Button
        type="button"
        variant="outline"
        size="icon"
        onClick={onPrev}
        disabled={prevDisabled}
        aria-label={`Previous ${label}`}
        className="h-9 w-9 rounded-none border-mathua-border bg-mathua-surface text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue"
      >
        <ChevronLeft size={15} aria-hidden="true" />
      </Button>
      <Button
        type="button"
        variant="outline"
        size="icon"
        onClick={onNext}
        disabled={nextDisabled}
        aria-label={`Next ${label}`}
        className="h-9 w-9 rounded-none border-mathua-border bg-mathua-surface text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue"
      >
        <ChevronRight size={15} aria-hidden="true" />
      </Button>
    </div>
  )
}
