'use client'

import { useState } from 'react'
import { setDailyXPGoal } from '../lib/api'

const MIN_GOAL = 5
const MAX_GOAL = 200
const STEP = 5
const PRESETS = [15, 30, 50]

function clamp(g: number): number {
  return Math.min(MAX_GOAL, Math.max(MIN_GOAL, Math.round(g)))
}

export default function GoalStepper({
  goal,
  onGoalChange,
}: {
  goal: number
  onGoalChange: (goal: number) => void
}) {
  const [pending, setPending] = useState(false)
  const [error, setError] = useState('')

  async function save(next: number) {
    const value = clamp(next)
    if (value === goal || pending) return
    const prev = goal
    setError('')
    onGoalChange(value) // optimistic — ProfileStats bar updates instantly
    setPending(true)
    try {
      await setDailyXPGoal(value)
    } catch {
      onGoalChange(prev) // rollback
      setError('Couldn’t save goal — check connection and retry.')
    } finally {
      setPending(false)
    }
  }

  return (
    <div className="border border-mathua-border bg-mathua-surface p-3 min-w-0">
      <div className="font-mono text-[10px] uppercase text-mathua-muted tracking-wider">Daily goal</div>
      <div className="flex items-center gap-2 mt-2">
        <button
          onClick={() => save(goal - STEP)}
          disabled={pending}
          aria-label="Decrease daily XP goal"
          className="font-mono text-sm text-mathua-muted border border-mathua-border-strong w-11 min-h-[44px] inline-flex items-center justify-center hover:text-mathua-blue hover:border-mathua-blue transition-colors disabled:opacity-40"
        >
          −
        </button>
        <div aria-live="polite" className="flex-1 text-center font-mono text-[16px] text-mathua-primary min-w-0">
          {goal} <span className="text-[11px] text-mathua-muted">XP/day</span>
        </div>
        <button
          onClick={() => save(goal + STEP)}
          disabled={pending}
          aria-label="Increase daily XP goal"
          className="font-mono text-sm text-mathua-muted border border-mathua-border-strong w-11 min-h-[44px] inline-flex items-center justify-center hover:text-mathua-blue hover:border-mathua-blue transition-colors disabled:opacity-40"
        >
          +
        </button>
      </div>
      <div className="flex gap-2 mt-2">
        {PRESETS.map((p) => (
          <button
            key={p}
            onClick={() => save(p)}
            disabled={pending}
            aria-label={`Set daily goal to ${p} XP`}
            aria-pressed={goal === p}
            className={`flex-1 font-mono text-[11px] min-h-[44px] inline-flex items-center justify-center border transition-colors disabled:opacity-40 ${
              goal === p
                ? 'border-mathua-blue text-mathua-blue bg-mathua-surface-elevated'
                : 'border-mathua-border text-mathua-muted hover:text-mathua-primary'
            }`}
          >
            {p}
          </button>
        ))}
      </div>
      {pending && (
        <div role="status" className="font-mono text-[10px] text-mathua-muted mt-2">
          Saving…
        </div>
      )}
      {error && (
        <div role="alert" className="font-mono text-[10px] text-mathua-red mt-2 break-words">
          {error}
        </div>
      )}
    </div>
  )
}
