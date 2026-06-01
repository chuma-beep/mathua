'use client'

import { useMemo, useState } from 'react'
import type { DailyActivity } from '../lib/api'

const CELL_W = 12
const CELL_H = 12
const GAP = 2
const LABEL_W = 28
const COLS = 53

const DAY_LABELS = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']
const MONTH_LABELS = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']

function intensityLevel(count: number, max: number): number {
  if (count === 0 || max === 0) return 0
  const ratio = count / max
  if (ratio <= 0.25) return 1
  if (ratio <= 0.50) return 2
  if (ratio <= 0.75) return 3
  return 4
}

function cellColor(level: number): string {
  if (level === 0) return 'transparent'
  const pcts = { 1: '18', 2: '38', 3: '60', 4: '85' }
  return `color-mix(in oklch, var(--accent-blue) ${pcts[level as keyof typeof pcts]}%, var(--bg))`
}

function formatDate(dateStr: string): string {
  const d = new Date(dateStr + 'T00:00:00')
  return d.toLocaleDateString('en-US', { weekday: 'short', month: 'short', day: 'numeric', year: 'numeric' })
}

interface Props {
  data: DailyActivity[]
}

export default function YearlyGrid({ data }: Props) {
  const [tooltip, setTooltip] = useState<{ date: string; questions: number; correct: number; x: number; y: number } | null>(null)

  const { grid, monthLabels, maxQuestions } = useMemo(() => {
    const activityMap = new Map<string, { questions: number; correct: number }>()
    let maxQ = 0
    for (const d of data) {
      const q = d.questions
      if (q > maxQ) maxQ = q
      activityMap.set(d.date, { questions: q, correct: d.correct })
    }

    const today = new Date()
    const todayStr = today.toISOString().slice(0, 10)
    const todayDate = new Date(todayStr + 'T00:00:00')

    // Start 364 days back, aligned to Sunday
    const startDate = new Date(todayDate)
    startDate.setDate(startDate.getDate() - 364)
    while (startDate.getDay() !== 0) {
      startDate.setDate(startDate.getDate() - 1)
    }

    const cells: { date: string; questions: number; correct: number; level: number; row: number; col: number }[] = []
    const ml: { col: number; label: string }[] = []

    let lastMonth = -1
    for (let c = 0; c < COLS; c++) {
      for (let r = 0; r < 7; r++) {
        const date = new Date(startDate)
        date.setDate(date.getDate() + c * 7 + r)
        const dateStr = date.toISOString().slice(0, 10)

        if (date > todayDate) continue

        if (date.getMonth() !== lastMonth && r === 0) {
          lastMonth = date.getMonth()
          ml.push({ col: c, label: MONTH_LABELS[date.getMonth()] })
        }

        const act = activityMap.get(dateStr)
        cells.push({
          date: dateStr,
          questions: act?.questions ?? 0,
          correct: act?.correct ?? 0,
          level: 0,
          row: r,
          col: c,
        })
      }
    }

    // Assign intensity levels based on max
    for (const cell of cells) {
      cell.level = intensityLevel(cell.questions, maxQ)
    }

    return { grid: cells, monthLabels: ml, maxQuestions: maxQ }
  }, [data])

  const svgW = LABEL_W + COLS * (CELL_W + GAP)
  const svgH = 16 + (7 * (CELL_H + GAP)) + GAP + 24

  return (
    <div style={{ overflowX: 'auto', position: 'relative' }}>
      <svg
        width={svgW}
        height={svgH}
        style={{ fontFamily: "'IBM Plex Mono', monospace", display: 'block', minWidth: svgW }}
      >
        {/* Month labels */}
        {monthLabels.map((m) => (
          <text
            key={m.col}
            x={LABEL_W + m.col * (CELL_W + GAP)}
            y={10}
            style={{ fontSize: 9, fill: 'var(--text-muted)' }}
          >
            {m.label}
          </text>
        ))}

        {/* Day-of-week labels */}
        {[1, 3, 5].map((r) => (
          <text
            key={r}
            x={0}
            y={20 + r * (CELL_H + GAP) + CELL_H / 2 + 4}
            style={{ fontSize: 9, fill: 'var(--text-muted)' }}
          >
            {DAY_LABELS[r].slice(0, 2)}
          </text>
        ))}

        {/* Cells */}
        {grid.map((cell) => (
          <rect
            key={`${cell.col}-${cell.row}`}
            x={LABEL_W + cell.col * (CELL_W + GAP)}
            y={16 + cell.row * (CELL_H + GAP) + GAP}
            width={CELL_W}
            height={CELL_H}
            style={{
              fill: cellColor(cell.level),
              stroke: 'var(--border)',
              strokeWidth: 0.5,
              cursor: 'pointer',
              shapeRendering: 'crispEdges',
            }}
            onMouseEnter={(e) => {
              const rect = (e.target as SVGRectElement).getBoundingClientRect()
              setTooltip({
                date: cell.date,
                questions: cell.questions,
                correct: cell.correct,
                x: rect.left + rect.width / 2,
                y: rect.top - 8,
              })
            }}
            onMouseLeave={() => setTooltip(null)}
          />
        ))}

        {/* Legend */}
        <text x={svgW - 180} y={svgH - 6} style={{ fontSize: 9, fill: 'var(--text-muted)' }}>
          Less
        </text>
        {[0, 1, 2, 3, 4].map((level) => (
          <rect
            key={level}
            x={svgW - 140 + level * (CELL_W + GAP)}
            y={svgH - 18}
            width={CELL_W}
            height={CELL_H}
            style={{
              fill: cellColor(level),
              stroke: 'var(--border)',
              strokeWidth: 0.5,
              shapeRendering: 'crispEdges',
            }}
          />
        ))}
        <text x={svgW - 64} y={svgH - 6} style={{ fontSize: 9, fill: 'var(--text-muted)' }}>
          More
        </text>
      </svg>

      {/* Tooltip */}
      {tooltip && (
        <div
          style={{
            position: 'fixed',
            left: tooltip.x,
            top: tooltip.y,
            transform: 'translate(-50%, -100%)',
            background: 'var(--bg)',
            border: '0.5px solid var(--border-strong)',
            padding: '4px 8px',
            fontFamily: "'IBM Plex Mono', monospace",
            fontSize: 11,
            color: 'var(--text-primary)',
            whiteSpace: 'nowrap',
            zIndex: 100,
            pointerEvents: 'none',
          }}
        >
          <div>{formatDate(tooltip.date)}</div>
          <div>
            {tooltip.questions} questions ·{' '}
            {tooltip.questions > 0
              ? Math.round((tooltip.correct / tooltip.questions) * 100)
              : 0}
            % correct
          </div>
        </div>
      )}
    </div>
  )
}
