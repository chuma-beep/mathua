'use client'

import { useMemo } from 'react'
import type { DailyActivity } from '../lib/api'

const CELL_SIZE = 11
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

    // Build calendar grid: 7 columns × up to 6 rows
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

    // Compute stats
    const totalQ = monthData.reduce((sum, d) => sum + d.questions, 0)
    const totalC = monthData.reduce((sum, d) => sum + d.correct, 0)
    const allConcepts = new Set<string>()
    for (const d of monthData) {
      for (const cid of d.concepts) allConcepts.add(cid)
    }
    const accuracy = totalQ > 0 ? Math.round((totalC / totalQ) * 100) : 0

    // Compute max questions for intensity
    let maxQ = 0
    for (const d of monthData) {
      if (d.questions > maxQ) maxQ = d.questions
    }

    // Assign levels to calendar cells
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
    <div
      style={{
        border: '0.5px solid var(--border)',
        marginBottom: 4,
      }}
    >
      <button
        onClick={onToggle}
        style={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          width: '100%',
          padding: '8px 10px',
          background: 'transparent',
          border: 'none',
          cursor: 'pointer',
          color: 'var(--text-primary)',
          fontFamily: "'IBM Plex Mono', monospace",
        }}
      >
        <span style={{ fontSize: 13, fontWeight: 400 }}>
          {MONTH_NAMES[month]} {year}
          <span style={{ color: 'var(--text-muted)', marginLeft: 8, fontSize: 11 }}>
            {activeDays} day{activeDays !== 1 ? 's' : ''}
          </span>
        </span>
        <span style={{ fontSize: 14, color: 'var(--text-muted)', transition: 'transform 0.2s', transform: expanded ? 'rotate(90deg)' : 'rotate(0deg)' }}>
          ▸
        </span>
      </button>

      {expanded && (
        <div style={{ padding: '0 10px 10px', borderTop: '0.5px solid var(--border)' }}>
          {/* Stats row */}
          <div
            style={{
              display: 'flex',
              gap: 16,
              padding: '8px 0',
              fontFamily: "'IBM Plex Mono', monospace",
              fontSize: 11,
              color: 'var(--text-secondary)',
            }}
          >
            <span>{stats.questions} questions</span>
            <span style={{ color: 'var(--text-muted)' }}>{stats.accuracy}% correct</span>
            <span style={{ color: 'var(--text-muted)' }}>{stats.concepts} topics</span>
          </div>

          {/* Mini calendar */}
          <div style={{ display: 'inline-block' }}>
            {/* Day headers */}
            <div style={{ display: 'flex', gap: GAP, marginBottom: 2 }}>
              {DAY_HEADERS.map((dh) => (
                <div
                  key={dh}
                  style={{
                    width: CELL_SIZE,
                    height: CELL_SIZE,
                    fontFamily: "'IBM Plex Mono', monospace",
                    fontSize: 8,
                    color: 'var(--text-muted)',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                  }}
                >
                  {dh}
                </div>
              ))}
            </div>

            {/* Calendar grid */}
            {calendar.map((week, wi) => (
              <div key={wi} style={{ display: 'flex', gap: GAP, marginBottom: GAP }}>
                {week.map((cell, di) => (
                  <div
                    key={di}
                    title={
                      cell
                        ? `${MONTH_NAMES[month].slice(0, 3)} ${cell.date}: ${cell.questions}q, ${cell.questions > 0 ? Math.round((cell.correct / cell.questions) * 100) : 0}%`
                        : ''
                    }
                    style={{
                      width: CELL_SIZE,
                      height: CELL_SIZE,
                      background: cell ? cellColor(cell.level) : 'transparent',
                      border: cell ? '0.5px solid var(--border)' : 'none',
                      fontFamily: "'IBM Plex Mono', monospace",
                      fontSize: 7,
                      color: cell && cell.questions > 0 ? 'var(--accent-blue)' : 'var(--text-muted)',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                    }}
                  >
                    {cell?.date}
                  </div>
                ))}
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}
