'use client'

import { useMemo, useState } from 'react'
import type { DailyActivity } from '../lib/api'
import MonthCard from './MonthCard'

interface Props {
  data: DailyActivity[]
}

function monthKey(year: number, month: number): string {
  return `${year}-${String(month + 1).padStart(2, '0')}`
}

export default function MonthlyCards({ data }: Props) {
  // Only months with activity plus the current month as an anchor — never a
  // wall of twelve mostly-empty accordions. Newest first. Expansion is keyed
  // by year-month (a bare month index collides across years).
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
    return [...active.values()].sort((a, b) => b.year - a.year || b.month - a.month)
  }, [data])

  const [expanded, setExpanded] = useState<string | null>(null)
  const openKey = expanded ?? (months.length > 0 ? monthKey(months[0].year, months[0].month) : null)

  return (
    <div className="w-full max-w-full min-w-0 mx-auto">
      {months.map((m) => {
        const key = monthKey(m.year, m.month)
        return (
          <MonthCard
            key={key}
            year={m.year}
            month={m.month}
            data={data}
            expanded={openKey === key}
            onToggle={() => setExpanded(openKey === key ? '' : key)}
          />
        )
      })}
    </div>
  )
}
