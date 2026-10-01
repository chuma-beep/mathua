'use client'

import { useState } from 'react'
import { setDailyXPGoal } from '../lib/api'
import { isLoggedIn } from '../lib/auth'
import { clampGoal, GOAL_PRESETS } from '../lib/plan'

const GUEST_KEY = 'mathua_daily_goal'

export function getGuestGoal(): number | null {
  if (typeof window === 'undefined') return null
  const v = Number(window.localStorage.getItem(GUEST_KEY))
  return v > 0 ? v : null
}

// DailyGoalControl: presets + custom target. Authed users write through to
// the server; guests persist locally. Calls back so the page can refresh.
export default function DailyGoalControl({ current, onSaved, bare }: { current: number; onSaved: (goal: number) => void; bare?: boolean }) {
  const [custom, setCustom] = useState('')
  const [busy, setBusy] = useState(false)
  const [msg, setMsg] = useState('')

  async function save(goal: number) {
    const g = clampGoal(goal)
    setBusy(true)
    setMsg('')
    try {
      if (isLoggedIn()) {
        await setDailyXPGoal(g)
      } else {
        window.localStorage.setItem(GUEST_KEY, String(g))
      }
      onSaved(g)
      setMsg(`Goal set to ${g} XP/day`)
      setCustom('')
    } catch (e) {
      setMsg(e instanceof Error ? e.message : 'Could not save goal')
    } finally {
      setBusy(false)
    }
  }

  return (
    <div className={bare ? undefined : 'mt-4 border border-mathua-border bg-mathua-surface p-4'}>
      <p className="font-mono text-[11px] uppercase tracking-wider text-mathua-muted">Daily XP target · currently {current}/day</p>
      <div className="mt-2 flex flex-wrap gap-2">
        {GOAL_PRESETS.map(p => (
          <button
            key={p}
            type="button"
            disabled={busy}
            onClick={() => save(p)}
            aria-pressed={current === p}
            className={`px-4 py-2 font-mono text-xs min-h-[40px] border transition-colors disabled:opacity-50 ${
              current === p
                ? 'border-mathua-blue text-mathua-blue'
                : 'border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue'
            }`}
          >
            {p} XP
          </button>
        ))}
      </div>
      <form
        className="mt-2 flex gap-2"
        onSubmit={e => {
          e.preventDefault()
          if (custom.trim()) save(Number(custom))
        }}
      >
        <label htmlFor="daily-goal-custom" className="sr-only">Custom daily XP target</label>
        <input
          id="daily-goal-custom"
          value={custom}
          onChange={e => setCustom(e.target.value)}
          placeholder="Custom (1–10000)…"
          inputMode="numeric"
          className="flex-1 min-w-0 bg-mathua-bg border border-mathua-border px-3 py-2 font-mono text-xs text-mathua-primary"
        />
        <button type="submit" disabled={busy || !custom.trim()} className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-5 font-mono text-xs min-h-[40px] disabled:opacity-50">
          Set
        </button>
      </form>
      <p className="mt-2 font-mono text-[11px] text-mathua-muted">≈ {Math.max(1, Math.round(current / 2))} questions/day at ~2 XP each. Changing your target never moves quiz eligibility or mastery.</p>
      {msg && <p className="mt-1 font-mono text-[11px] text-mathua-secondary">{msg}</p>}
    </div>
  )
}
