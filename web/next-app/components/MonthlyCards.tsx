'use client'

import { useMemo, useState } from 'react'
import type { DailyActivity } from '../lib/api'
import MonthCard from './MonthCard'

interface Props {
  data: DailyActivity[]
}

export default function MonthlyCards({ data }: Props) {
  const currentMonth = new Date().getMonth()
  const [expanded, setExpanded] = useState<number>(currentMonth)

  const months = useMemo(() => {
    const result: { year: number; month: number }[] = []
    const now = new Date()
    for (let i = 11; i >= 0; i--) {
      const d = new Date(now.getFullYear(), now.getMonth() - i, 1)
      result.push({ year: d.getFullYear(), month: d.getMonth() })
    }
    return result
  }, [])

  return (
    <div className="w-full max-w-full min-w-0 mx-auto">
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
