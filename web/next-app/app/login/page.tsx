'use client'

import Loading from '../../components/Loading'

import { useReducer } from 'react'
import { useTheme } from '../../hooks/useTheme'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import Footer from '../../components/Footer'
import SectionHeader from '../../components/SectionHeader'
import { signup, login, validateToken } from '../../lib/api'
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

export default function LoginPage() {
  const { mounted } = useTheme()
  const { push } = useRouter()
  const [state, dispatch] = useReducer(loginReducer, initialState)

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
      push(res.diagnostic_completed ? '/session' : '/onboard')
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
            <div className="mt-3 text-center">
              <Link href="/session" className="text-mathua-muted text-xs hover:text-mathua-secondary">
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
