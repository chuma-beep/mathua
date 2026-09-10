'use client'

import { useMemo } from 'react'
import type { DailyActivity } from '../lib/api'

const GAP = 2

const DAY_HEADERS = ['Su', 'Mo', 'Tu', 'We', 'Th', 'Fr', 'Sa']
const MONTH_NAMES = ['January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December']

function cellColor(level: number): string {
  if (level === 0) return 'transparent'
  const pcts = { 1: '18', 2: '38', 3: '60', 4: '85' }
  return `color-mix(in oklch, var(--accent-blue) ${pcts[level as keyof typeof pcts]}%, var(--bg))`
}

interface Props {
  year: number
  month: number
  data: DailyActivity[]
  expanded: boolean
  onToggle: () => void
}

export default function MonthCard({ year, month, data, expanded, onToggle }: Props) {
  const monthData = useMemo(() => {
    const prefix = `${year}-${String(month + 1).padStart(2, '0')}`
    return data.filter((d) => d.date.startsWith(prefix))
  }, [data, year, month])

  const { calendar, stats } = useMemo(() => {
    const firstDay = new Date(year, month, 1)
    const lastDay = new Date(year, month + 1, 0)
    const daysInMonth = lastDay.getDate()
    const startDow = firstDay.getDay()

    const activityMap = new Map<string, { questions: number; correct: number; concepts: string[] }>()
    for (const d of data) {
      activityMap.set(d.date, { questions: d.questions, correct: d.correct, concepts: d.concepts })
    }

    const calendar: ({ date: number; questions: number; correct: number; level: number } | null)[][] = []
    let day = 1
    for (let week = 0; week < 6; week++) {
      const row: ({ date: number; questions: number; correct: number; level: number } | null)[] = []
      for (let dow = 0; dow < 7; dow++) {
        if ((week === 0 && dow < startDow) || day > daysInMonth) {
          row.push(null)
        } else {
          const dateStr = `${year}-${String(month + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`
          const act = activityMap.get(dateStr)
          const questions = act?.questions ?? 0
          const correct = act?.correct ?? 0
          row.push({ date: day, questions, correct, level: 0 })
          day++
        }
      }
      calendar.push(row)
    }

    const totalQ = monthData.reduce((sum, d) => sum + d.questions, 0)
    const totalC = monthData.reduce((sum, d) => sum + d.correct, 0)
    const allConcepts = new Set<string>()
    for (const d of monthData) {
      for (const cid of d.concepts) allConcepts.add(cid)
    }
    const accuracy = totalQ > 0 ? Math.round((totalC / totalQ) * 100) : 0

    let maxQ = 0
    for (const d of monthData) {
      if (d.questions > maxQ) maxQ = d.questions
    }

    for (const week of calendar) {
      for (const cell of week) {
        if (cell) {
          if (cell.questions === 0 || maxQ === 0) {
            cell.level = 0
          } else {
            const ratio = cell.questions / maxQ
            if (ratio <= 0.25) cell.level = 1
            else if (ratio <= 0.50) cell.level = 2
            else if (ratio <= 0.75) cell.level = 3
            else cell.level = 4
          }
        }
      }
    }

    return {
      calendar,
      stats: { questions: totalQ, correct: totalC, accuracy, concepts: allConcepts.size },
    }
  }, [data, year, month, monthData])

  const activeDays = monthData.filter((d) => d.questions > 0).length

  return (
    <div className="w-full max-w-full min-w-0 border-[0.5px] border-mathua-border mb-1 overflow-hidden bg-transparent">
      <button
        type="button"
        onClick={onToggle}
        aria-expanded={expanded}
        aria-label={`${MONTH_NAMES[month]} ${year} activity`}
        className="flex w-full max-w-full min-w-0 items-center justify-between gap-2 px-3 py-2 min-h-[36px] bg-transparent border-none cursor-pointer text-mathua-primary font-mono"
      >
        <span className="min-w-0 flex-1 truncate text-left text-[13px] font-normal">
          {MONTH_NAMES[month]} {year}
          <span className="ml-2 text-[11px] text-mathua-muted">
            {activeDays} day{activeDays !== 1 ? 's' : ''}
          </span>
        </span>
        <span
          className="shrink-0 text-[14px] text-mathua-muted transition-transform duration-200"
          style={{ transform: expanded ? 'rotate(90deg)' : 'rotate(0deg)' }}
        >
          ▸
        </span>
      </button>

      {expanded && (
        <div className="px-3 pb-3 pt-0 border-t-[0.5px] border-mathua-border w-full max-w-full min-w-0 overflow-hidden">
          <div className="flex flex-wrap gap-x-4 gap-y-1 py-2 font-mono text-[11px] text-mathua-secondary">
            <span>{stats.questions} questions</span>
            <span className="text-mathua-muted">{stats.accuracy}% correct</span>
            <span className="text-mathua-muted">{stats.concepts} topics</span>
          </div>

          <div className="flex justify-center w-full max-w-full min-w-0 overflow-hidden">
            <div className="inline-block max-w-full">
              <div className="flex gap-[2px] mb-[2px]">
                {DAY_HEADERS.map((dh) => (
                  <div
                    key={dh}
                    className="w-[11px] h-[11px] flex items-center justify-center font-mono text-[8px] text-mathua-muted shrink-0"
                  >
                    {dh}
                  </div>
                ))}
              </div>

              {calendar.map((week, wi) => (
                <div key={wi} className="flex gap-[2px] mb-[2px]">
                  {week.map((cell, di) => (
                    <div
                      key={di}
                      title={
                        cell
                          ? `${MONTH_NAMES[month].slice(0, 3)} ${cell.date}: ${cell.questions}q, ${cell.questions > 0 ? Math.round((cell.correct / cell.questions) * 100) : 0}%`
                          : ''
                      }
                      className="w-[11px] h-[11px] flex items-center justify-center font-mono text-[7px] shrink-0 border-[0.5px]"
                      style={{
                        background: cell ? cellColor(cell.level) : 'transparent',
                        borderColor: cell ? 'var(--border)' : 'transparent',
                        color: cell && cell.questions > 0 ? 'var(--accent-blue)' : 'var(--text-muted)',
                      }}
                    >
                      {cell?.date}
                    </div>
                  ))}
                </div>
              ))}
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
