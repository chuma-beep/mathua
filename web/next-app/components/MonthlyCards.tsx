'use client'

import { useEffect, useMemo, useRef, useState } from 'react'
import type { DailyActivity } from '../lib/api'
import MonthCard, { MONTH_NAMES } from './MonthCard'

interface Props {
  data: DailyActivity[]
}

function monthKey(year: number, month: number): string {
  return `${year}-${String(month + 1).padStart(2, '0')}`
}

export default function MonthlyCards({ data }: Props) {
  // Horizontal month strip, oldest → newest: only months with activity plus
  // the current month as an anchor — never a wall of empty accordions. Every
  // visible month shows its grid (no expand/collapse); the strip opens
  // scrolled to the newest end.
  const months = useMemo(() => {
    const now = new Date()
    const current = { year: now.getFullYear(), month: now.getMonth() }
    const active = new Map<string, { year: number; month: number }>()
    for (const d of data) {
      if (d.questions <= 0) continue
      const parts = d.date.split('-').map(Number)
      if (parts.length < 2 || !Number.isFinite(parts[0]) || !Number.isFinite(parts[1])) continue
      const key = monthKey(parts[0], parts[1] - 1)
      if (!active.has(key)) active.set(key, { year: parts[0], month: parts[1] - 1 })
    }
    const currentKey = monthKey(current.year, current.month)
    if (!active.has(currentKey)) active.set(currentKey, current)
    return [...active.values()].sort((a, b) => a.year - b.year || a.month - b.month)
  }, [data])

  const stripRef = useRef<HTMLDivElement>(null)
  useEffect(() => {
    const el = stripRef.current
    if (el) el.scrollLeft = el.scrollWidth
  }, [months.length])

  // Mobile pager: one month at a time with prev/next buttons (no scroll
  // bar). Desktop keeps the sideways strip below. Both render; CSS picks.
  const [selected, setSelected] = useState<number | null>(null)
  const current = selected ?? months.length - 1
  const clamped = months.length === 0 ? 0 : Math.min(Math.max(current, 0), months.length - 1)
  const shown = months[clamped]

  return (
    <>
      <div className="sm:hidden">
        <div className="flex items-center justify-between gap-2 mb-3">
          <button
            type="button"
            onClick={() => setSelected(clamped - 1)}
            disabled={clamped <= 0}
            aria-label="Previous month"
            className="border border-mathua-border px-4 py-2 font-mono text-xs text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue disabled:opacity-40 min-h-[44px] inline-flex items-center"
          >
            ‹ Prev
          </button>
          <span className="font-mono text-xs text-mathua-primary" aria-live="polite">
            {months.length > 0 ? `${MONTH_NAMES[shown.month]} ${shown.year}` : ''}
          </span>
          <button
            type="button"
            onClick={() => setSelected(clamped + 1)}
            disabled={clamped >= months.length - 1}
            aria-label="Next month"
            className="border border-mathua-border px-4 py-2 font-mono text-xs text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue disabled:opacity-40 min-h-[44px] inline-flex items-center"
          >
            Next ›
          </button>
        </div>
        {shown && (
          <div className="flex justify-center">
            <MonthCard
              key={`${shown.year}-${shown.month}`}
              year={shown.year}
              month={shown.month}
              data={data}
            />
          </div>
        )}
      </div>
      <div
        ref={stripRef}
        data-testid="month-strip"
        className="w-full max-w-full min-w-0 overflow-x-auto hidden sm:block"
      >
        <div className="flex gap-3 w-max min-w-full px-0.5 py-0.5">
          {months.map((m) => (
            <MonthCard
              key={monthKey(m.year, m.month)}
              year={m.year}
              month={m.month}
              data={data}
            />
          ))}
        </div>
      </div>
    </>
  )
}
