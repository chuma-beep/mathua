'use client'

import { useMemo, useState } from 'react'
import type { DailyActivity } from '../lib/api'
import MonthCard from './MonthCard'

interface Props {
  data: DailyActivity[]
}

export default function MonthlyCards({ data }: Props) {
  const currentMonth = new Date().getMonth()
  const currentYear = new Date().getFullYear()
  const [expanded, setExpanded] = useState<number>(currentMonth)

  const months = useMemo(() => {
    // Generate last 12 months from today
    const result: { year: number; month: number }[] = []
    const now = new Date()
    for (let i = 11; i >= 0; i--) {
      const d = new Date(now.getFullYear(), now.getMonth() - i, 1)
      result.push({ year: d.getFullYear(), month: d.getMonth() })
    }
    return result
  }, [])

  return (
    <div style={{ maxWidth: 820, margin: '0 auto', width: '100%' }}>
      {months.map((m) => (
        <MonthCard
          key={`${m.year}-${m.month}`}
          year={m.year}
          month={m.month}
          data={data}
          expanded={expanded === m.month}
          onToggle={() => setExpanded(expanded === m.month ? -1 : m.month)}
        />
      ))}
    </div>
  )
}
