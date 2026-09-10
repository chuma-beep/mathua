'use client'

import { useCallback, useEffect, useState } from 'react'
import Header from '../../../components/Header'
import Footer from '../../../components/Footer'
import KatexContent from '../../../components/KatexContent'
import {
  adminLogin,
  adminLogout,
  listReports,
  updateReportStatus,
  type QuestionReport,
  type ReportStatus,
} from '../../../lib/api'

const STATUSES = ['open', 'confirmed', 'fixed', 'dismissed', 'all']
const SOURCES = ['all', 'study', 'diagnostic', 'quiz', 'lesson', 'concept', 'review']
const SESSION_KEY = 'mathua_admin_session'

interface AdminSession {
  token: string
  expires_at: string
}

function loadStoredSession(): AdminSession | null {
  if (typeof window === 'undefined') return null
  try {
    const raw = sessionStorage.getItem(SESSION_KEY)
    if (!raw) return null
    const s = JSON.parse(raw) as AdminSession
    if (!s.token) return null
    if (s.expires_at && new Date(s.expires_at).getTime() < Date.now()) {
      sessionStorage.removeItem(SESSION_KEY)
      return null
    }
    return s
  } catch {
    return null
  }
}

export default function AdminReportsPage() {
  const [session, setSession] = useState<AdminSession | null>(null)
  const [password, setPassword] = useState('')
  const [loggingIn, setLoggingIn] = useState(false)
  const [status, setStatus] = useState('open')
  const [source, setSource] = useState('all')
  const [reports, setReports] = useState<QuestionReport[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const dropSession = useCallback((message: string) => {
    sessionStorage.removeItem(SESSION_KEY)
    setSession(null)
    setReports([])
    setError(message)
  }, [])

  const load = useCallback(
    async (s: AdminSession, filter: string) => {
      setLoading(true)
      setError('')
      try {
        setReports(await listReports(s.token, filter))
      } catch {
        dropSession('Session expired or password changed — log in again.')
      } finally {
        setLoading(false)
      }
    },
    [dropSession],
  )

  useEffect(() => {
    const stored = loadStoredSession()
    if (stored) {
      setSession(stored)
      load(stored, 'open')
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  async function handleLogin() {
    if (!password || loggingIn) return
    setLoggingIn(true)
    setError('')
    try {
      const res = await adminLogin(password)
      const s = { token: res.token, expires_at: res.expires_at }
      sessionStorage.setItem(SESSION_KEY, JSON.stringify(s))
      setSession(s)
      setPassword('')
      await load(s, status)
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Login failed')
    } finally {
      setLoggingIn(false)
    }
  }

  async function handleLogout() {
    if (session) await adminLogout(session.token)
    sessionStorage.removeItem(SESSION_KEY)
    setSession(null)
    setReports([])
    setError('')
  }

  async function resolve(id: number, next: ReportStatus) {
    if (!session) return
    try {
      await updateReportStatus(session.token, id, next)
      setReports(prev => prev.map(r => (r.id === id ? { ...r, status: next } : r)))
    } catch {
      dropSession('Session expired or password changed — log in again.')
    }
  }

  return (
    <>
      <Header />
      <main className="max-w-container mx-auto px-4 sm:px-6 pt-20 pb-16">
        <h1 className="font-serif text-2xl text-mathua-primary">Content reports</h1>
        <p className="mt-1 font-mono text-xs text-mathua-muted">
          User complaints about questions, explanations, lessons, and diagrams.
        </p>

        {!session ? (
          <div className="mt-6 max-w-md">
            <p className="font-mono text-[11px] text-mathua-secondary uppercase tracking-wider mb-2">
              Admin login
            </p>
            <div className="flex flex-col sm:flex-row gap-2">
              <form
                onSubmit={e => { e.preventDefault(); handleLogin() }}
                className="flex flex-col sm:flex-row gap-2 flex-1"
              >
                <label htmlFor="admin-password" className="sr-only">Admin password</label>
                <input
                  id="admin-password"
                  aria-label="Admin password"
                  type="password"
                  value={password}
                  onChange={e => setPassword(e.target.value)}
                  placeholder="Admin password"
                  autoComplete="current-password"
                  className="flex-1 bg-mathua-surface border border-mathua-border px-3 h-11 text-sm font-mono text-mathua-primary placeholder:text-mathua-muted outline-none"
                />
                <button
                  type="submit"
                  disabled={loggingIn || !password}
                  className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white transition-colors px-4 h-11 text-sm font-mono disabled:opacity-50"
                >
                  {loggingIn ? 'Checking…' : 'Log in'}
                </button>
              </form>
            </div>
            {error && <p className="mt-3 font-mono text-xs text-red-400">{error}</p>}
          </div>
        ) : (
          <>
            <div className="mt-6 flex items-center gap-3 flex-wrap">
              <span className="font-mono text-[11px] text-mathua-green uppercase tracking-wider">
                ● Admin session
              </span>
              <span className="font-mono text-[11px] text-mathua-muted">
                expires {new Date(session.expires_at).toLocaleString()}
              </span>
              <button
                type="button"
                onClick={handleLogout}
                className="font-mono text-[11px] text-mathua-muted hover:text-mathua-primary uppercase tracking-wider"
              >
                Log out
              </button>
            </div>

            <div className="mt-4 flex gap-2 flex-wrap">
              {STATUSES.map(s => (
                <button
                  type="button"
                  key={s}
                  onClick={() => { setStatus(s); load(session, s) }}
                  className={`px-3 h-8 font-mono text-[11px] uppercase tracking-wider border transition-colors ${
                    status === s
                      ? 'border-mathua-blue text-mathua-blue'
                      : 'border-mathua-border text-mathua-muted hover:text-mathua-primary'
                  }`}
                >
                  {s}
                </button>
              ))}
            </div>

            <div className="mt-2 flex gap-2 flex-wrap items-center">
              <span className="font-mono text-[10px] text-mathua-muted uppercase tracking-wider">Source:</span>
              {SOURCES.map(s => (
                <button
                  type="button"
                  key={s}
                  onClick={() => setSource(s)}
                  className={`px-3 h-8 font-mono text-[11px] uppercase tracking-wider border transition-colors ${
                    source === s
                      ? 'border-mathua-blue text-mathua-blue'
                      : 'border-mathua-border text-mathua-muted hover:text-mathua-primary'
                  }`}
                >
                  {s}
                </button>
              ))}
            </div>

            {error && <p className="mt-4 font-mono text-xs text-red-400">{error}</p>}
            {loading && <p className="mt-4 font-mono text-xs text-mathua-muted">Loading…</p>}

            <div className="mt-6 space-y-3">
              {reports.filter(r => source === 'all' || (r.source || 'study') === source).map(r => (
                <div key={r.id} className="border border-mathua-border bg-mathua-surface p-4">
                  <div className="flex items-center gap-2 flex-wrap font-mono text-[10px] uppercase tracking-wider">
                    <span className="text-mathua-blue">#{r.id}</span>
                    <span className="text-mathua-secondary">{r.kind}</span>
                    <span className="text-mathua-secondary">{r.reason}</span>
                    <span className="text-mathua-muted">{r.concept_id}</span>
                    <span className="text-mathua-muted">{r.status}</span>
                    <span className="text-mathua-muted ml-auto">{r.created_at}</span>
                  </div>
                  {r.question && (
                    <KatexContent className="mt-2 text-xs font-mono text-mathua-primary whitespace-pre-wrap">
                      {r.question}
                    </KatexContent>
                  )}
                  {r.expected && (
                    <p className="mt-1 font-mono text-[11px] text-mathua-secondary">Expected: {r.expected}</p>
                  )}
                  {r.explanation && (
                    <KatexContent className="mt-1 text-[11px] font-mono text-mathua-secondary">
                      {r.explanation}
                    </KatexContent>
                  )}
                  {r.detail && (
                    <p className="mt-1 font-mono text-[11px] text-mathua-primary">“{r.detail}”</p>
                  )}
                  <div className="mt-1 font-mono text-[10px] text-mathua-muted">
                    reporter: {r.reporter_id || 'anon'}
                    {r.lesson_id ? ` · lesson: ${r.lesson_id}` : ''}
                    {r.source ? ` · source: ${r.source}` : ''}
                  </div>
                  <div className="mt-2 flex gap-2">
                    {(['confirmed', 'fixed', 'dismissed'] as ReportStatus[]).map(next => (
                      <button
                        type="button"
                        key={next}
                        onClick={() => resolve(r.id, next)}
                        disabled={r.status === next}
                        className="font-mono text-[10px] uppercase tracking-wider border border-mathua-border text-mathua-muted hover:text-mathua-blue hover:border-mathua-blue px-2 h-7 disabled:opacity-40"
                      >
                        {next}
                      </button>
                    ))}
                  </div>
                </div>
              ))}
              {!loading && !error && reports.length === 0 && (
                <p className="font-mono text-xs text-mathua-muted">No reports with status “{status}”.</p>
              )}
            </div>
          </>
        )}
      </main>
      <Footer />
    </>
  )
}
