'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { getAttempts, resetAccount, type AttemptRecord } from '../lib/api'

const RESET_PHRASE = 'reset my progress'

/**
 * Whether the typed phrase confirms the action.
 *
 * Case-insensitive, and deliberately so. This input is `type="text"` with no
 * `autoCapitalize`, so Safari capitalises the first letter as you type and Android
 * keyboards frequently autocorrect: "reset my progress" arrives as "Reset my progress".
 * Both the comparison here and the server's were exact, so the button stayed disabled
 * forever on a phone and the server would have answered the capitalised phrase with
 * "confirmation phrase does not match" — telling someone who had typed exactly what
 * the label asked for that they had mistyped it. Chromium does not auto-capitalise,
 * which is why the suite was green throughout.
 *
 * The phrase confirms intent; it is not a secret. Casing and surrounding whitespace say
 * nothing about how deliberate the action was, and refusing to arm over them locks a
 * learner out of a reset they explicitly asked for. Anything else still fails.
 */
export function phraseConfirms(typed: string, expected: string): boolean {
  return typed.trim().toLowerCase() === expected
}

export function attemptsToCSV(rows: AttemptRecord[]): string {
  const head = 'timestamp,concept_id,question,answer,expected,correct,elapsed_seconds,source'
  const esc = (v: string | number | boolean) => `"${String(v ?? '').replace(/"/g, '""')}"`
  return [head, ...rows.map(a => [
    a.timestamp, a.concept_id, a.question, a.answer, a.expected,
    a.correct ? '1' : '0', a.elapsed_seconds, a.source,
  ].map(esc).join(','))].join('\n')
}

// DangerZone: destructive, explicit, and honest about scope. Reset wipes
// learning evidence; choices (account, settings, goal, course,
// destination/deadline/pace) survive.
export default function DangerZone() {
  const { push } = useRouter()
  const [phrase, setPhrase] = useState('')
  const [busy, setBusy] = useState(false)
  const [exporting, setExporting] = useState(false)
  const [exported, setExported] = useState(false)
  const [error, setError] = useState('')
  const armed = phraseConfirms(phrase, RESET_PHRASE)

  async function handleExport() {
    setExporting(true)
    setError('')
    try {
      const rows: AttemptRecord[] = []
      let offset = 0
      for (;;) {
        const res = await getAttempts({ limit: 500, offset })
        rows.push(...res.attempts)
        if (rows.length >= res.total || res.attempts.length === 0) break
        offset += res.attempts.length
      }
      const blob = new Blob([attemptsToCSV(rows)], { type: 'text/csv' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = 'mathua-attempts.csv'
      a.click()
      URL.revokeObjectURL(url)
      setExported(true)
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Export failed')
    } finally {
      setExporting(false)
    }
  }

  async function handleReset(e: React.FormEvent) {
    e.preventDefault()
    if (!armed || busy) return
    setBusy(true)
    setError('')
    try {
      // Send the canonical phrase, so the server's own gate compares like with like.
      await resetAccount(RESET_PHRASE)
      push('/profile')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Reset failed')
      setBusy(false)
    }
  }

  return (
    <div className="border border-mathua-red bg-mathua-surface p-4 sm:p-6 min-w-0">
      <span className="font-mono text-sm text-mathua-red">Reset progress</span>
      <p className="text-mathua-muted text-xs mt-1">
        Wipes your learning record. This cannot be undone.
      </p>
      <div className="mt-3 grid grid-cols-1 sm:grid-cols-2 gap-3">
        <div>
          <p className="font-mono text-[11px] uppercase tracking-wider text-mathua-muted mb-1">Wiped</p>
          <ul className="font-mono text-[11px] text-mathua-secondary space-y-0.5">
            <li>Concept mastery, streaks, review schedule</li>
            <li>Attempts log and sessions</li>
            <li>XP — profile and leaderboard</li>
            <li>Diagnostic flag and plan evidence</li>
          </ul>
        </div>
        <div>
          <p className="font-mono text-[11px] uppercase tracking-wider text-mathua-muted mb-1">Kept</p>
          <ul className="font-mono text-[11px] text-mathua-secondary space-y-0.5">
            <li>Account, avatar, settings, daily goal</li>
            <li>Course enrollment</li>
            <li>Destination, deadline and pace</li>
          </ul>
        </div>
      </div>
      <div className="mt-4">
        <button
          type="button"
          onClick={handleExport}
          disabled={exporting}
          className="border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue rounded-none px-4 h-9 text-xs font-mono disabled:opacity-40"
        >
          {exporting ? 'Exporting…' : exported ? 'Export again' : 'Export attempts log (CSV)'}
        </button>
      </div>
      <form onSubmit={handleReset} className="mt-4">
        <label htmlFor="reset-phrase" className="font-mono text-[11px] uppercase text-mathua-muted">
          Type <span className="text-mathua-primary">reset my progress</span> to confirm
        </label>
        <div className="flex flex-col sm:flex-row gap-2 mt-2">
          <input
            id="reset-phrase"
            type="text"
            value={phrase}
            onChange={e => setPhrase(e.target.value)}
            placeholder="reset my progress"
            autoComplete="off"
            autoCapitalize="none"
            autoCorrect="off"
            spellCheck={false}
            enterKeyHint="go"
            aria-describedby={phrase && !armed ? 'reset-phrase-hint' : undefined}
            className="flex-1 min-w-0 bg-mathua-bg border border-mathua-border px-3 py-2 font-mono text-xs text-mathua-primary"
          />
          <button
            type="submit"
            disabled={!armed || busy}
            className="border border-mathua-red text-mathua-red hover:bg-mathua-red hover:text-white rounded-none px-4 h-10 text-xs font-mono disabled:opacity-40 min-h-[40px]"
          >
            {busy ? 'Resetting…' : 'Reset everything above'}
          </button>
        </div>
      </form>
      {/* A greyed-out button with no explanation is indistinguishable from a broken
          one, which is how this hid for as long as it did. */}
      {phrase && !armed && (
        <p id="reset-phrase-hint" role="status" className="mt-2 font-mono text-[11px] text-mathua-secondary">
          That does not match. Type <span className="text-mathua-primary">{RESET_PHRASE}</span> exactly.
        </p>
      )}
      {error && <p className="mt-2 font-mono text-[11px] text-mathua-red">{error}</p>}
    </div>
  )
}
