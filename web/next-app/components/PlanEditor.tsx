'use client'

import { useCallback, useEffect, useMemo, useState } from 'react'
import { getDestinations, getEstimate, getScores, savePlan, getCurrentPlan, type DestinationStatus, type EstimateRes, type StudyPlan } from '../lib/api'
import { getUserInfo, getGuestId } from '../lib/auth'
import { clampGoal, daysFor, monthLabel } from '../lib/plan'

// Learning planner: destination + (daily effort XOR deadline) → estimate.
// Dates render as months (estimates, never promises); quiz eligibility and
// mastery are independent of every number here. Shared by /plan (full) and
// Settings (compact): compact hides the workload detail list, the what-if
// slider, and the plan-delta line, keeping destination/mode/inputs/save.
export default function PlanEditor({ compact = false }: { compact?: boolean }) {
  const [dests, setDests] = useState<DestinationStatus[]>([])
  const [destId, setDestId] = useState('')
  const [mode, setMode] = useState<'effort' | 'deadline'>('effort')
  const [dailyGoal, setDailyGoal] = useState(10)
  const [deadlineDays, setDeadlineDays] = useState(90)
  const [restDays, setRestDays] = useState(1)
  const [whatIf, setWhatIf] = useState(10)
  const [est, setEst] = useState<EstimateRes | null>(null)
  const [saved, setSaved] = useState<StudyPlan | null>(null)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    const info = getUserInfo()
    const sid = info?.student_id || getGuestId() || ''
    Promise.all([getDestinations().catch(() => [] as DestinationStatus[]), (sid ? getScores(sid) : Promise.resolve(null)).catch(() => null), getCurrentPlan().catch(() => ({ plan: null as StudyPlan | null }))]).then(([d, s, p]) => {
      setDests(d)
      if (d.length > 0) setDestId(d[0].id)
      const guestGoal = (() => {
        try {
          const v = Number(window.localStorage.getItem('mathua_daily_goal'))
          return v > 0 ? v : null
        } catch {
          return null
        }
      })()
      const g = guestGoal ?? s?.daily_xp_goal ?? 30
      setDailyGoal(g)
      setWhatIf(g)
      if (p.plan) setSaved(p.plan)
      setLoading(false)
    })
  }, [])

  const fetchEst = useCallback(async (id: string, goal: number, deadline: number, rest: number) => {
    if (!id) return
    setError('')
    try {
      const r = await getEstimate(id, {
        daily_goal: clampGoal(goal),
        deadline_days: mode === 'deadline' ? deadline : undefined,
        rest_days: rest,
        diagnostic_min: 10,
      })
      setEst(r)
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Estimate failed')
    }
  }, [mode])

  useEffect(() => {
    if (destId) fetchEst(destId, dailyGoal, deadlineDays, restDays)
  }, [destId, dailyGoal, deadlineDays, restDays, fetchEst])

  const whatIfDays = useMemo(
    () => (est ? daysFor(est.estimate.xp_remaining, whatIf, restDays) : -1),
    [est, whatIf, restDays],
  )

  async function handleSave() {
    if (!destId) return
    setSaving(true)
    try {
      const p = await savePlan({ destination: destId, daily_goal: clampGoal(dailyGoal), deadline_days: mode === 'deadline' ? deadlineDays : 0, rest_days: restDays })
      setSaved(p)
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Save failed')
    } finally {
      setSaving(false)
    }
  }

  if (loading) return <div className="p-6 font-mono text-xs text-mathua-muted">Loading planner…</div>

  const e = est?.estimate
  return (
    <div>
      <label htmlFor="plan-dest" className="mt-5 block font-mono text-[11px] uppercase tracking-wider text-mathua-muted">What do you want to learn?</label>
      <select id="plan-dest" value={destId} onChange={ev => setDestId(ev.target.value)} className="mt-1 w-full bg-mathua-surface border border-mathua-border px-3 py-2.5 font-mono text-sm text-mathua-primary min-h-[44px]">
        {dests.map(d => <option key={d.id} value={d.id}>{d.name} ({Math.round(d.pct * 100)}% done)</option>)}
      </select>

      <div className="mt-4 flex gap-2" role="group" aria-label="Planning direction">
        {(['effort', 'deadline'] as const).map(m => (
          <button key={m} type="button" onClick={() => setMode(m)} aria-pressed={mode === m} className={`flex-1 border px-4 py-2.5 font-mono text-xs min-h-[44px] transition-colors ${mode === m ? 'border-mathua-blue text-mathua-blue' : 'border-mathua-border text-mathua-secondary hover:border-mathua-blue'}`}>
            {m === 'effort' ? 'Plan by daily effort' : 'Plan by deadline'}
          </button>
        ))}
      </div>

      {mode === 'effort' ? (
        <div className="mt-3">
          <label htmlFor="plan-goal" className="font-mono text-[11px] uppercase tracking-wider text-mathua-muted">Daily XP target</label>
          <div className="mt-1 flex gap-2">
            {[5, 10, 20, 50].map(p => (
              <button key={p} type="button" onClick={() => setDailyGoal(p)} aria-pressed={dailyGoal === p} className={`flex-1 border px-2 py-2 font-mono text-xs min-h-[40px] ${dailyGoal === p ? 'border-mathua-blue text-mathua-blue' : 'border-mathua-border text-mathua-secondary'}`}>{p}</button>
            ))}
            <input id="plan-goal" value={dailyGoal} onChange={ev => setDailyGoal(clampGoal(Number(ev.target.value) || 0))} inputMode="numeric" className="w-20 bg-mathua-surface border border-mathua-border px-2 py-2 font-mono text-xs text-mathua-primary text-center" aria-label="Custom daily XP target" />
          </div>
        </div>
      ) : (
        <div className="mt-3">
          <label htmlFor="plan-deadline" className="font-mono text-[11px] uppercase tracking-wider text-mathua-muted">Target completion (days from now)</label>
          <div className="mt-1 flex gap-2">
            {[30, 90, 180, 365].map(p => (
              <button key={p} type="button" onClick={() => setDeadlineDays(p)} aria-pressed={deadlineDays === p} className={`flex-1 border px-2 py-2 font-mono text-xs min-h-[40px] ${deadlineDays === p ? 'border-mathua-blue text-mathua-blue' : 'border-mathua-border text-mathua-secondary'}`}>{p >= 365 ? '1 year' : p >= 30 ? `${Math.round(p / 30)} mo` : `${p}d`}</button>
            ))}
            <input id="plan-deadline" value={deadlineDays} onChange={ev => setDeadlineDays(Math.max(1, Number(ev.target.value) || 1))} inputMode="numeric" className="w-20 bg-mathua-surface border border-mathua-border px-2 py-2 font-mono text-xs text-mathua-primary text-center" aria-label="Custom deadline in days" />
          </div>
        </div>
      )}

      <div className="mt-3">
        <label htmlFor="plan-rest" className="font-mono text-[11px] uppercase tracking-wider text-mathua-muted">Rest days per week: {restDays}</label>
        <input id="plan-rest" type="range" min={0} max={3} value={restDays} onChange={ev => setRestDays(Number(ev.target.value))} className="mt-1 w-full" />
      </div>

      {error && <p className="mt-3 font-mono text-xs text-red-400">{error}</p>}

      {e && (
        <div className="mt-5 border border-mathua-border bg-mathua-surface p-5">
          <div className="flex items-baseline justify-between gap-2 flex-wrap">
            <p className="font-mono text-[11px] uppercase tracking-wider text-mathua-muted">Remaining workload</p>
            <p className="font-mono text-sm text-mathua-primary">≈ {Math.round(e.xp_remaining).toLocaleString()} XP</p>
          </div>
          {mode === 'effort' ? (
            <p className="mt-3 font-serif text-2xl text-mathua-primary">
              ≈ {e.days > 0 ? Math.ceil(e.days) : '–'} days
              {e.finish_date && <span className="text-mathua-secondary text-lg"> · {monthLabel(e.finish_date)}</span>}
            </p>
          ) : (
            <div className="mt-3">
              <p className="font-serif text-2xl text-mathua-primary">≈ {Math.ceil(e.required_per_day)} XP/day needed</p>
              {!e.feasible && <p className="mt-1 font-mono text-xs text-red-400">Above 3× your recent pace — extend the deadline or raise daily effort gradually.</p>}
            </div>
          )}
          <p className="mt-1 font-mono text-[11px] text-mathua-muted">
            {est?.pace.source === 'measured'
              ? `At your recent pace (${est.pace.rate.toFixed(0)} XP/day, ${est.pace.active_days}/${est.pace.trailing_days} active days).`
              : 'No history yet — assuming you hit your target every day.'}
          </p>
          {!compact && (
            <dl className="mt-4 space-y-1 font-mono text-xs text-mathua-secondary">
              <div className="flex justify-between"><dt>Topics remaining</dt><dd className="text-mathua-primary">{e.remaining.length} of {e.total}</dd></div>
              <div className="flex justify-between"><dt>Focused time remaining</dt><dd className="text-mathua-primary">≈ {Math.round(e.time_min_remaining / 60)}h {(Math.round(e.time_min_remaining % 60))}m</dd></div>
              <div className="flex justify-between"><dt>Reviews due</dt><dd className="text-mathua-primary">{e.reviews_due}</dd></div>
              <div className="flex justify-between"><dt>Mastery checks ahead</dt><dd className="text-mathua-primary">{e.quizzes_ahead}</dd></div>
              <div className="flex justify-between"><dt>Diagnostic + planning time</dt><dd className="text-mathua-primary">≈ {e.diagnostic_min} min</dd></div>
              <div className="flex justify-between"><dt>Initial estimate coverage</dt><dd className="text-mathua-primary">{est?.probes.length} probe topics</dd></div>
            </dl>
          )}
          {est?.fresh_start ? (
            <p className="mt-3 font-mono text-xs text-mathua-secondary">
              Fresh start — complete a few questions and your estimate will sharpen. No behind-schedule figure until a new baseline lands.
            </p>
          ) : (
            !compact && est?.plan_delta_days !== undefined && est.plan_delta_days !== null && (
              <p className="mt-3 font-mono text-xs text-mathua-blue">
                {(est.plan_delta_days >= 0 ? '≈' + Math.abs(est.plan_delta_days).toFixed(0) + ' days ahead of' : '≈' + Math.abs(est.plan_delta_days).toFixed(0) + ' days behind') + ' your saved plan'}
              </p>
            )
          )}
          {!compact && (
            <div className="mt-4">
              <label htmlFor="plan-whatif" className="font-mono text-[11px] uppercase tracking-wider text-mathua-muted">What if {whatIf} XP/day? {whatIfDays > 0 && <span className="text-mathua-primary">→ ≈ {Math.ceil(whatIfDays)} days</span>}</label>
              <input id="plan-whatif" type="range" min={2} max={60} step={1} value={whatIf} onChange={ev => setWhatIf(Number(ev.target.value))} className="mt-1 w-full" />
            </div>
          )}
          <div className="mt-4">
            <button type="button" onClick={handleSave} disabled={saving} className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-6 py-2.5 font-mono text-xs min-h-[44px] disabled:opacity-50">
              {saving ? 'Saving…' : saved ? 'Update my plan' : 'Plan my learning'}
            </button>
          </div>
          <p className="mt-3 font-mono text-[11px] text-mathua-muted">Illustrative estimate from Mathua&apos;s measured curriculum workload — actual completion depends on pace, reviews, and checks. It never changes quiz eligibility or mastery rules.</p>
        </div>
      )}
    </div>
  )
}
