'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { getAttempts, deleteAccount, getMe, type AttemptRecord } from '../lib/api'
import { getUserInfo, signOut, clearGuest } from '../lib/auth'
import { attemptsToCSV } from './DangerZone'

const DELETE_PHRASE = 'delete my account'

// DeleteAccount: destructive, explicit, and honest about scope. Deletion
// wipes the account row and every owned row (progress, attempts, sessions,
// quizzes, identities, settings, plans, avatars, filed reports) — nothing
// is kept. Password accounts additionally confirm the current password;
// OAuth-only and guest rows need the phrase alone.
export default function DeleteAccount() {
  const { push } = useRouter()
  const [phrase, setPhrase] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  // Tri-state, not a boolean: `null` means "we could not find out" (getMe
  // failed or is still in flight). Hiding the field on `null` used to strand a
  // password account — the delete went out with no password, the server
  // answered 401 "current password is required", and there was no field on
  // screen to fill. So when unknown we SHOW the field but do not require it:
  // the server is the real gate, and an OAuth-only learner can ignore it.
  const [needsPassword, setNeedsPassword] = useState<boolean | null>(null)
  const [busy, setBusy] = useState(false)
  const [exporting, setExporting] = useState(false)
  const [exported, setExported] = useState(false)
  const armed = phrase.trim() === DELETE_PHRASE && (needsPassword !== true || password !== '')

  useEffect(() => {
    const cached = getUserInfo()?.has_password
    if (cached !== undefined) setNeedsPassword(!!cached)
    getMe()
      .then(me => setNeedsPassword(!!me.has_password))
      .catch(() => setNeedsPassword(cached === undefined ? null : !!cached))
  }, [])

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

  async function handleDelete(e: React.FormEvent) {
    e.preventDefault()
    if (!armed || busy) return
    setBusy(true)
    setError('')
    try {
      if (needsPassword === true) {
        await deleteAccount({ phrase: phrase.trim(), password })
      } else {
        // When we don't know whether a password is required, send it if the
        // learner typed one; the server decides, and a 401 can now be retried
        // because the field is on screen.
        await deleteAccount(password ? { phrase: phrase.trim(), password } : { phrase: phrase.trim() })
      }
      signOut()
      clearGuest()
      push('/')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Delete failed')
      setBusy(false)
    }
  }

  return (
    <div className="border border-mathua-red bg-mathua-surface p-4 sm:p-6 min-w-0">
      <span className="font-mono text-sm text-mathua-red">Delete account</span>
      <p className="text-mathua-muted text-xs mt-1">
        Deletes your account and everything in it. This cannot be undone.
      </p>
      <div className="mt-3 grid grid-cols-1 sm:grid-cols-2 gap-3">
        <div>
          <p className="font-mono text-[11px] uppercase tracking-wider text-mathua-muted mb-1">Wiped</p>
          <ul className="font-mono text-[11px] text-mathua-secondary space-y-0.5">
            <li>Account, avatar, settings, daily goal</li>
            <li>Concept mastery, streaks, review schedule</li>
            <li>Attempts log, sessions and quizzes</li>
            <li>Plans, filed reports and XP — profile and leaderboard</li>
          </ul>
        </div>
        <div>
          <p className="font-mono text-[11px] uppercase tracking-wider text-mathua-muted mb-1">Kept</p>
          <ul className="font-mono text-[11px] text-mathua-secondary space-y-0.5">
            <li>Nothing — deletion is total</li>
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
          {exporting ? 'Exporting…' : exported ? 'Export again' : 'Export attempts log (CSV) first'}
        </button>
      </div>
      <form onSubmit={handleDelete} className="mt-4">
        <label htmlFor="delete-phrase" className="font-mono text-[11px] uppercase text-mathua-muted">
          Type <span className="text-mathua-primary">delete my account</span> to confirm
        </label>
        <div className="flex flex-col sm:flex-row gap-2 mt-2">
          <input
            id="delete-phrase"
            type="text"
            value={phrase}
            onChange={e => setPhrase(e.target.value)}
            placeholder="delete my account"
            autoComplete="off"
            className="flex-1 min-w-0 bg-mathua-bg border border-mathua-border px-3 py-2 font-mono text-xs text-mathua-primary"
          />
          {needsPassword !== false && (
            <input
              id="delete-password"
              type="password"
              value={password}
              onChange={e => setPassword(e.target.value)}
              placeholder="Current password"
              autoComplete="current-password"
              // Distinct from the change-password field higher up this page:
              // two controls sharing an accessible name made the Settings
              // page ambiguous for screen readers and for tests.
              aria-label="Current password, to confirm deletion"
              className="flex-1 min-w-0 bg-mathua-bg border border-mathua-border px-3 py-2 font-mono text-xs text-mathua-primary"
            />
          )}
          <button
            type="submit"
            disabled={!armed || busy}
            className="border border-mathua-red text-mathua-red hover:bg-mathua-red hover:text-white rounded-none px-4 h-10 text-xs font-mono disabled:opacity-40 min-h-[40px]"
          >
            {busy ? 'Deleting…' : 'Delete my account'}
          </button>
        </div>
      </form>
      {error && <p className="mt-2 font-mono text-[11px] text-mathua-red">{error}</p>}
    </div>
  )
}
