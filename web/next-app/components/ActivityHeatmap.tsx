'use client'

import { useState, useEffect } from 'react'
import type { DailyActivity } from '../lib/api'
import YearlyGrid from './YearlyGrid'
import MonthlyCards from './MonthlyCards'

const STORAGE_KEY = 'mathua-activity-view'

interface Props {
  data: DailyActivity[]
}

export default function ActivityHeatmap({ data }: Props) {
  const [view, setView] = useState<'yearly' | 'monthly'>('monthly')

  useEffect(() => {
    const stored = localStorage.getItem(STORAGE_KEY)
    if (stored === 'yearly' || stored === 'monthly') {
      setView(stored)
    }
  }, [])

  function switchView(v: 'yearly' | 'monthly') {
    setView(v)
    localStorage.setItem(STORAGE_KEY, v)
  }

  const totalQuestions = data.reduce((sum, d) => sum + d.questions, 0)
  const totalCorrect = data.reduce((sum, d) => sum + d.correct, 0)
  const totalDays = data.filter((d) => d.questions > 0).length

  return (
    <div className="w-full max-w-full min-w-0 overflow-hidden">
      <div className="flex flex-col gap-2 mb-4 sm:flex-row sm:items-center sm:justify-between sm:gap-3">
        <div className="font-mono text-[11px] text-mathua-muted min-w-0 truncate">
          {totalDays} active days · {totalQuestions} questions ·{' '}
          {totalQuestions > 0 ? Math.round((totalCorrect / totalQuestions) * 100) : 0}% correct
        </div>

        <div className="inline-flex shrink-0 self-start sm:self-auto border-[0.5px] border-mathua-border-strong rounded-none overflow-hidden">
          <button
            type="button"
            onClick={() => switchView('monthly')}
            aria-pressed={view === 'monthly'}
            className={`font-mono text-[11px] px-3 min-h-[36px] inline-flex items-center justify-center border-none cursor-pointer transition-colors ${
              view === 'monthly' ? 'bg-mathua-blue text-white' : 'bg-transparent text-mathua-muted hover:text-mathua-primary'
            }`}
          >
            Monthly
          </button>
          <button
            type="button"
            onClick={() => switchView('yearly')}
            aria-pressed={view === 'yearly'}
            className={`font-mono text-[11px] px-3 min-h-[36px] inline-flex items-center justify-center border-none cursor-pointer transition-colors ${
              view === 'yearly' ? 'bg-mathua-blue text-white' : 'bg-transparent text-mathua-muted hover:text-mathua-primary'
            }`}
          >
            Yearly
          </button>
        </div>
      </div>

      <div className="w-full max-w-full min-w-0 overflow-hidden">
        {view === 'yearly' ? (
          <YearlyGrid data={data} />
        ) : (
          <MonthlyCards data={data} />
        )}
      </div>
    </div>
  )
}
