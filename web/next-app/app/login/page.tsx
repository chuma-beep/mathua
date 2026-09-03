'use client'

import Loading from '../../components/Loading'

import { Suspense, useEffect, useReducer, useState } from 'react'
import { useTheme } from '../../hooks/useTheme'
import { useRouter, useSearchParams } from 'next/navigation'
import Link from 'next/link'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import Footer from '../../components/Footer'
import SectionHeader from '../../components/SectionHeader'
import { signup, login, validateToken, API_BASE, getConfig } from '../../lib/api'
import { setToken, setUserInfo, clearToken } from '../../lib/auth'

type LoginState = {
  tab: 'login' | 'signup'
  name: string
  username: string
  password: string
  error: string
  loading: boolean
  showPassword: boolean
  attempted: boolean
  fieldErrors: { name?: string; username?: string; password?: string }
}

type LoginAction =
  | { type: 'SET_TAB'; tab: 'login' | 'signup' }
  | { type: 'SET_NAME'; name: string }
  | { type: 'SET_USERNAME'; username: string }
  | { type: 'SET_PASSWORD'; password: string }
  | { type: 'SET_ERROR'; error: string }
  | { type: 'SET_LOADING'; loading: boolean }
  | { type: 'TOGGLE_SHOW_PASSWORD' }
  | { type: 'SET_ATTEMPTED'; attempted: boolean }
  | { type: 'SET_FIELD_ERRORS'; fieldErrors: LoginState['fieldErrors'] }

const initialState: LoginState = {
  tab: 'login',
  name: '',
  username: '',
  password: '',
  error: '',
  loading: false,
  showPassword: false,
  attempted: false,
  fieldErrors: {},
}

function loginReducer(state: LoginState, action: LoginAction): LoginState {
  switch (action.type) {
    case 'SET_TAB':
      return { ...state, tab: action.tab, error: '', fieldErrors: {} }
    case 'SET_NAME':
      return { ...state, name: action.name }
    case 'SET_USERNAME':
      return { ...state, username: action.username }
    case 'SET_PASSWORD':
      return { ...state, password: action.password }
    case 'SET_ERROR':
      return { ...state, error: action.error }
    case 'SET_LOADING':
      return { ...state, loading: action.loading }
    case 'TOGGLE_SHOW_PASSWORD':
      return { ...state, showPassword: !state.showPassword }
    case 'SET_ATTEMPTED':
      return { ...state, attempted: action.attempted }
    case 'SET_FIELD_ERRORS':
      return { ...state, fieldErrors: action.fieldErrors }
  }
}

// Mirrors internal/auth.ValidatePassword — keep both in sync.
function passwordIssues(pw: string): string[] {
  const issues: string[] = []
  if (pw.length < 8) issues.push('at least 8 characters')
  if (!/\d/.test(pw)) issues.push('a number')
  if (!/[^A-Za-z0-9]/.test(pw)) issues.push('a special character')
  return issues
}

function validateFields(
  tab: 'login' | 'signup',
  values: { name: string; username: string; password: string }
): LoginState['fieldErrors'] {
  const errors: LoginState['fieldErrors'] = {}
  if (tab === 'signup' && !values.name.trim()) errors.name = 'Name is required'
  if (!values.username.trim()) errors.username = 'Username is required'
  if (!values.password) {
    errors.password = 'Password is required'
  } else if (tab === 'signup') {
    const issues = passwordIssues(values.password)
    if (issues.length > 0) errors.password = 'Password needs ' + issues.join(', ')
  }
  return errors
}

const inputClassName = (hasError: boolean) =>
  `w-full mt-1 bg-mathua-code border rounded-none h-12 px-4 pr-12 font-mono text-base text-mathua-primary placeholder:text-mathua-muted focus:outline-none ${
    hasError ? 'border-mathua-red' : 'border-mathua-border focus:border-mathua-blue'
  }`

function EyeIcon({ off }: { off?: boolean }) {
  return (
    <svg
      width="20"
      height="20"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="1.8"
      strokeLinecap="round"
      strokeLinejoin="round"
      aria-hidden="true"
    >
      <path d="M2 12s3.5-7 10-7 10 7 10 7-3.5 7-10 7-10-7-10-7Z" />
      <circle cx="12" cy="12" r="3" />
      {off && <path d="m4 4 16 16" />}
    </svg>
  )
}

function LoginInner() {
  const { mounted } = useTheme()
  const { push } = useRouter()
  const searchParams = useSearchParams()
  const [state, dispatch] = useReducer(loginReducer, initialState)
  const [googleReady, setGoogleReady] = useState(false)
  const [googleError, setGoogleError] = useState('')
  const [googleLoading, setGoogleLoading] = useState(false)

  // Google redirect callback: ?token=...&id=...&name=...
  useEffect(() => {
    const token = searchParams.get('token')
    const id = searchParams.get('id')
    const name = searchParams.get('name')
    const err = searchParams.get('error')
    if (err) setGoogleError(err === 'google_denied' ? 'Google sign-in was cancelled' : 'Google sign-in failed — try again')
    if (token && id) {
      setToken(token)
      setUserInfo({ student_id: id, name: name || 'Google user', username: '', concepts_mastered: 0, current_streak: 0, level: 'Novice', diagnostic_completed: false })
      validateToken().then(v => { if (!v.valid) clearToken() })
      push('/profile')
    }
  }, [searchParams, push])

  // GIS One-Tap init
  useEffect(() => {
    let cancelled = false
    getConfig().then(cfg => {
      const cid = (cfg as unknown as { google_client_id?: string }).google_client_id
      if (!cid || cancelled) return
      const src = 'https://accounts.google.com/gsi/client'
      const existing = document.querySelector(`script[src="${src}"]`) as HTMLScriptElement | null
      const init = () => {
        const g = (window as unknown as { google?: { accounts: { id: { initialize: (o: unknown) => void; renderButton: (a: HTMLElement, b: unknown) => void; prompt: () => void } } } }).google
        if (!g) return
        try {
          g.accounts.id.initialize({
            client_id: cid,
            callback: async (resp: { credential?: string }) => {
              if (!resp?.credential) return
              setGoogleLoading(true)
              setGoogleError('')
              try {
                const r = await fetch(`${API_BASE}/api/auth/google`, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ id_token: resp.credential }) })
                if (!r.ok) throw new Error(await r.text())
                const data = await r.json() as { token: string; student_id: string; name: string; diagnostic_completed: boolean }
                setToken(data.token)
                setUserInfo({ student_id: data.student_id, name: data.name, username: '', concepts_mastered: 0, current_streak: 0, level: 'Novice', diagnostic_completed: data.diagnostic_completed })
                push('/profile')
              } catch (e) { setGoogleError((e as Error).message || 'Google sign-in failed') } finally { setGoogleLoading(false) }
            },
            auto_select: false,
            cancel_on_tap_outside: true,
          })
          const el = document.getElementById('g_id_onload')
          if (el) g.accounts.id.renderButton(el, { theme: 'outline', size: 'large', width: 320, text: 'continue_with', shape: 'square' })
          g.accounts.id.prompt()
          setGoogleReady(true)
        } catch {}
      }
      if (existing) { init(); return }
      const s = document.createElement('script'); s.src = src; s.async = true; s.defer = true; s.onload = init; document.head.appendChild(s)
    }).catch(()=>{})
    return () => { cancelled = true }
  }, [push])

  const handleGoogleRedirect = () => {
    setGoogleError('')
    window.location.href = `${API_BASE}/api/auth/google/login?return=${encodeURIComponent('/profile')}`
  }

  if (!mounted) return <div style={{ background: 'var(--bg)', minHeight: '100vh' }} />

  const updateField = (action: LoginAction) => {
    dispatch(action)
    if (state.attempted) {
      const next = {
        name: action.type === 'SET_NAME' ? (action as { name: string }).name : state.name,
        username:
          action.type === 'SET_USERNAME' ? (action as { username: string }).username : state.username,
        password:
          action.type === 'SET_PASSWORD' ? (action as { password: string }).password : state.password,
      }
      dispatch({ type: 'SET_FIELD_ERRORS', fieldErrors: validateFields(state.tab, next) })
    }
  }

  const handleSubmit = async () => {
    dispatch({ type: 'SET_ATTEMPTED', attempted: true })
    const fieldErrors = validateFields(state.tab, state)
    dispatch({ type: 'SET_FIELD_ERRORS', fieldErrors })
    if (Object.keys(fieldErrors).length > 0) return
    dispatch({ type: 'SET_ERROR', error: '' })
    dispatch({ type: 'SET_LOADING', loading: true })
    try {
      const res = state.tab === 'signup'
        ? await signup(state.name.trim(), state.username.trim(), state.password)
        : await login(state.username.trim(), state.password)
      setToken(res.token)
      setUserInfo({ student_id: res.student_id, name: res.name, username: state.username.trim(), concepts_mastered: 0, current_streak: 0, level: 'Novice', diagnostic_completed: res.diagnostic_completed })
      const verified = await validateToken()
      if (!verified.valid) {
        clearToken()
        dispatch({ type: 'SET_ERROR', error: 'Something went wrong, but we\'re working on it.' })
        dispatch({ type: 'SET_LOADING', loading: false })
        return
      }
      push('/profile')
    } catch (e: any) {
      dispatch({ type: 'SET_ERROR', error: e.message || 'Authentication failed' })
    } finally { dispatch({ type: 'SET_LOADING', loading: false }) }
  }

  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-4 sm:px-6 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0">
        <section className="pt-8 max-w-md mx-auto mt-8 sm:mt-12 min-w-0 overflow-hidden">
          <span className="flex mb-4">
            <Link href="/" className="text-mathua-secondary text-sm hover:text-mathua-primary">Back</Link>
          </span>
          <SectionHeader label="Account" title={state.tab === 'login' ? 'Welcome back' : 'Create account'} />
          <div className="flex gap-2 mt-6 mb-4 min-w-0">
            <button onClick={() => dispatch({ type: 'SET_TAB', tab: 'login' })} className={`flex-1 min-w-0 min-h-[44px] rounded-none h-12 text-sm font-medium px-2 ${state.tab === 'login' ? 'bg-mathua-blue text-white' : 'bg-mathua-surface-elevated border border-mathua-border text-mathua-secondary'}`}>Login</button>
            <button onClick={() => dispatch({ type: 'SET_TAB', tab: 'signup' })} className={`flex-1 min-w-0 min-h-[44px] rounded-none h-12 text-sm font-medium px-2 ${state.tab === 'signup' ? 'bg-mathua-blue text-white' : 'bg-mathua-surface-elevated border border-mathua-border text-mathua-secondary'}`}>Sign Up</button>
          </div>
          <div className="bg-mathua-surface border border-mathua-border rounded-none p-4 sm:p-6 space-y-4 w-full max-w-full min-w-0 overflow-hidden">
            {state.tab === 'signup' && (
              <div>
                <label htmlFor="name" className="font-mono text-[10px] uppercase text-mathua-muted">Name</label>
                <input
                  id="name"
                  type="text"
                  value={state.name}
                  onChange={(e) => updateField({ type: 'SET_NAME', name: e.target.value })}
                  onKeyDown={(e) => e.key === 'Enter' && handleSubmit()}
                  placeholder="Your name"
                  autoComplete="name"
                  aria-invalid={!!state.fieldErrors.name}
                  aria-describedby={state.fieldErrors.name ? 'name-error' : undefined}
                  className={inputClassName(!!state.fieldErrors.name)}
                />
                {state.fieldErrors.name && <p id="name-error" className="text-mathua-red text-xs mt-1">{state.fieldErrors.name}</p>}
              </div>
            )}
            <div>
              <label htmlFor="username" className="font-mono text-[10px] uppercase text-mathua-muted">Username</label>
              <input
                id="username"
                type="text"
                value={state.username}
                onChange={(e) => updateField({ type: 'SET_USERNAME', username: e.target.value })}
                onKeyDown={(e) => e.key === 'Enter' && handleSubmit()}
                placeholder="username"
                autoComplete="username"
                aria-invalid={!!state.fieldErrors.username}
                aria-describedby={state.fieldErrors.username ? 'username-error' : undefined}
                className={inputClassName(!!state.fieldErrors.username)}
              />
              {state.fieldErrors.username && <p id="username-error" className="text-mathua-red text-xs mt-1">{state.fieldErrors.username}</p>}
            </div>
            <div>
              <label htmlFor="password" className="font-mono text-[10px] uppercase text-mathua-muted">Password</label>
              <div className="relative">
                <input
                  id="password"
                  type={state.showPassword ? 'text' : 'password'}
                  value={state.password}
                  onChange={(e) => updateField({ type: 'SET_PASSWORD', password: e.target.value })}
                  onKeyDown={(e) => e.key === 'Enter' && handleSubmit()}
                  placeholder="password"
                  autoComplete={state.tab === 'signup' ? 'new-password' : 'current-password'}
                  aria-invalid={!!state.fieldErrors.password}
                  aria-describedby={state.fieldErrors.password ? 'password-error' : undefined}
                  className={inputClassName(!!state.fieldErrors.password)}
                />
                <button
                  type="button"
                  onClick={() => dispatch({ type: 'TOGGLE_SHOW_PASSWORD' })}
                  aria-label={state.showPassword ? 'Hide password' : 'Show password'}
                  aria-pressed={state.showPassword}
                  title={state.showPassword ? 'Hide password' : 'Show password'}
                  className="absolute right-3 top-1/2 -translate-y-1/2 mt-0.5 text-mathua-muted hover:text-mathua-primary"
                >
                  <EyeIcon off={state.showPassword} />
                </button>
              </div>
              {state.fieldErrors.password && <p id="password-error" className="text-mathua-red text-xs mt-1">{state.fieldErrors.password}</p>}
            </div>
            {state.error && <p className="text-mathua-red text-xs">{state.error}</p>}
            <button onClick={handleSubmit} disabled={state.loading} data-testid="auth-submit" className="w-full border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none h-12 font-medium text-sm disabled:opacity-50">
              {state.loading ? (<><Loading inline size={13} /> Loading…</>) : state.tab === 'signup' ? 'Create Account' : 'Login'}
            </button>
            <div className="flex items-center gap-3 my-2">
              <div className="h-px flex-1 bg-mathua-border" />
              <span className="font-mono text-[10px] uppercase text-mathua-muted">or</span>
              <div className="h-px flex-1 bg-mathua-border" />
            </div>
            <div className="space-y-2">
              <div id="g_id_onload" className="flex justify-center min-h-[44px] items-center" />
              {!googleReady && (
                <button onClick={handleGoogleRedirect} disabled={googleLoading} className="w-full border border-mathua-border bg-white text-[#3c4043] hover:bg-gray-50 rounded-none h-12 font-medium text-sm flex items-center justify-center gap-2 disabled:opacity-50">
                  <svg width="18" height="18" viewBox="0 0 48 48"><path fill="#EA4335" d="M24 9.5c3.54 0 6.71 1.22 9.21 3.6l6.85-6.85C35.9 2.38 30.47 0 24 0 14.62 0 6.51 5.38 2.56 13.22l7.98 6.19C12.43 13.72 17.74 9.5 24 9.5z"/><path fill="#4285F4" d="M46.98 24.55c0-1.57-.15-3.09-.38-4.55H24v9.02h12.94c-.58 2.96-2.26 5.48-4.78 7.18l7.73 6c4.51-4.18 7.09-10.36 7.09-17.65z"/><path fill="#FBBC05" d="M10.53 28.59c-.48-1.45-.76-2.99-.76-5.09s.27-3.64.76-5.09l-7.98-6.19C.92 15.77 0 19.69 0 24s.92 8.23 2.56 11.78l7.97-6.19z"/><path fill="#34A853" d="M24 48c6.48 0 11.93-2.13 15.89-5.81l-7.73-6c-2.15 1.45-4.92 2.3-8.16 2.3-6.26 0-11.57-4.22-13.47-9.91l-7.98 6.19C6.51 42.62 14.62 48 24 48z"/><path fill="none" d="M0 0h48v48H0z"/></svg>
                  {googleLoading ? 'Connecting…' : 'Continue with Google'}
                </button>
              )}
              {googleReady && (
                <button onClick={handleGoogleRedirect} disabled={googleLoading} className="w-full border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue rounded-none h-11 font-mono text-xs">
                  {googleLoading ? 'Connecting…' : 'Continue with Google (redirect)'}
                </button>
              )}
              {googleError && <p className="text-mathua-red text-xs text-center">{googleError}</p>}
              <p className="font-mono text-[10px] text-mathua-muted text-center">Google will link to existing account by email if username exists (username stays unique)</p>
            </div>
            <div className="mt-3 text-center">
              <Link href="/profile" className="text-mathua-muted text-xs hover:text-mathua-secondary">
                Skip for now: try without account
              </Link>
            </div>
          </div>
        </section>
      </div>
      <Footer />
      <BottomTabs />
    </>
  )
}

export default function LoginPage() {
  return (
    <Suspense fallback={<div style={{ background: 'var(--bg)', minHeight: '100vh' }} />}>
      <LoginInner />
    </Suspense>
  )
}
