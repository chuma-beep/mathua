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
    <div>
      {/* Summary + toggle row */}
      <div
        style={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          flexWrap: 'wrap',
          gap: 12,
          marginBottom: 16,
        }}
      >
        <div
          style={{
            fontFamily: "'IBM Plex Mono', monospace",
            fontSize: 11,
            color: 'var(--text-muted)',
          }}
        >
          {totalDays} active days · {totalQuestions} questions ·{' '}
          {totalQuestions > 0 ? Math.round((totalCorrect / totalQuestions) * 100) : 0}% correct
        </div>

        <div style={{ display: 'flex', border: '0.5px solid var(--border-strong)' }}>
          <button
            onClick={() => switchView('monthly')}
            style={{
              fontFamily: "'IBM Plex Mono', monospace",
              fontSize: 11,
              padding: '4px 12px',
              background: view === 'monthly' ? 'var(--accent-blue)' : 'transparent',
              color: view === 'monthly' ? 'var(--bg)' : 'var(--text-muted)',
              border: 'none',
              cursor: 'pointer',
              transition: 'background 0.15s',
            }}
          >
            Monthly
          </button>
          <button
            onClick={() => switchView('yearly')}
            style={{
              fontFamily: "'IBM Plex Mono', monospace",
              fontSize: 11,
              padding: '4px 12px',
              background: view === 'yearly' ? 'var(--accent-blue)' : 'transparent',
              color: view === 'yearly' ? 'var(--bg)' : 'var(--text-muted)',
              border: 'none',
              cursor: 'pointer',
              transition: 'background 0.15s',
            }}
          >
            Yearly
          </button>
        </div>
      </div>

      {/* Active view */}
      {view === 'yearly' ? (
        <YearlyGrid data={data} />
      ) : (
        <MonthlyCards data={data} />
      )}
    </div>
  )
}
