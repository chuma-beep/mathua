'use client'

import { useMemo } from 'react'
import type { DailyActivity } from '../lib/api'

function last7(data: DailyActivity[]): { date: string; questions: number; correct: number }[] {
  const byDate = new Map(data.map((d) => [d.date, d]))
  const days: { date: string; questions: number; correct: number }[] = []
  const today = new Date()
  for (let i = 6; i >= 0; i--) {
    const d = new Date(today)
    d.setDate(today.getDate() - i)
    const key = d.toISOString().slice(0, 10)
    const found = byDate.get(key)
    days.push({ date: key, questions: found?.questions ?? 0, correct: found?.correct ?? 0 })
  }
  return days
}

export default function WeeklyChart({ data }: { data: DailyActivity[] }) {
  const days = useMemo(() => last7(data), [data])
  const max = Math.max(1, ...days.map((d) => d.questions))
  const total = days.reduce((s, d) => s + d.questions, 0)
  const correct = days.reduce((s, d) => s + d.correct, 0)
  const pct = total > 0 ? Math.round((correct / total) * 100) : 0
  const summary = `${total} questions in the last 7 days · ${pct}% correct`

  return (
    <div className="border border-mathua-border bg-mathua-surface p-3 min-w-0">
      <div className="font-mono text-[10px] uppercase text-mathua-muted tracking-wider">This week</div>
      {total === 0 ? (
        <div className="mt-2">
          <p className="font-mono text-[11px] text-mathua-secondary leading-relaxed">
            No activity this week yet — answer a Lesson to light up the chart.
          </p>
          <a
            href="/study"
            className="mt-2 w-full border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white px-4 py-2 font-mono text-xs min-h-[44px] inline-flex items-center justify-center transition-colors"
          >
            Open Study →
          </a>
        </div>
      ) : (
        <>
          <div
            role="img"
            aria-label={`Weekly activity chart. ${summary}`}
            className="flex items-end gap-1.5 h-20 mt-3"
          >
            {days.map((d, i) => {
              const h = d.questions > 0 ? Math.max(8, Math.round((d.questions / max) * 100)) : 4
              const acc = d.questions > 0 ? d.correct / d.questions : 0
              const isToday = i === days.length - 1
              return (
                <div key={d.date} className="flex-1 min-w-0 flex flex-col items-center justify-end h-full" title={`${d.date}: ${d.questions} questions`}>
                  <div
                    className={`w-full ${isToday ? 'bg-mathua-blue' : acc >= 0.7 ? 'bg-mathua-teal' : 'bg-mathua-border-strong'}`}
                    style={{ height: `${h}%` }}
                  />
                  <div className={`font-mono text-[8px] mt-1 ${isToday ? 'text-mathua-blue' : 'text-mathua-muted'}`}>
                    {['S', 'M', 'T', 'W', 'T', 'F', 'S'][new Date(d.date + 'T12:00:00').getDay()]}
                  </div>
                </div>
              )
            })}
          </div>
          <div className="font-mono text-[10px] text-mathua-muted mt-2 truncate" title={summary}>
            {total} questions · {pct}% correct
          </div>
          <details className="mt-1">
            <summary className="font-mono text-[10px] text-mathua-blue cursor-pointer min-h-[44px] inline-flex items-center">
              Daily breakdown
            </summary>
            <table className="w-full font-mono text-[10px] text-mathua-secondary mt-1">
              <tbody>
                {days.map((d) => (
                  <tr key={d.date}>
                    <td className="py-0.5 pr-2">{d.date}</td>
                    <td className="py-0.5 text-right">{d.questions} q · {d.correct} ✓</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </details>
        </>
      )}
    </div>
  )
}
