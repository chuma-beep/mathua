'use client'

import { useReducer } from 'react'
import { useTheme } from '../../hooks/useTheme'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import Header from '../../components/Header'
import Footer from '../../components/Footer'
import SectionHeader from '../../components/SectionHeader'
import { signup, login } from '../../lib/api'
import { setToken, setUserInfo } from '../../lib/auth'

type LoginState = {
  tab: 'login' | 'signup'
  name: string
  username: string
  password: string
  error: string
  loading: boolean
}

type LoginAction =
  | { type: 'SET_TAB'; tab: 'login' | 'signup' }
  | { type: 'SET_NAME'; name: string }
  | { type: 'SET_USERNAME'; username: string }
  | { type: 'SET_PASSWORD'; password: string }
  | { type: 'SET_ERROR'; error: string }
  | { type: 'SET_LOADING'; loading: boolean }

const initialState: LoginState = {
  tab: 'login',
  name: '',
  username: '',
  password: '',
  error: '',
  loading: false,
}

function loginReducer(state: LoginState, action: LoginAction): LoginState {
  switch (action.type) {
    case 'SET_TAB':
      return { ...state, tab: action.tab, error: '' }
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
  }
}

export default function LoginPage() {
  const { mounted } = useTheme()
  const { push } = useRouter()
  const [state, dispatch] = useReducer(loginReducer, initialState)

  if (!mounted) return <div style={{ background: 'var(--bg)', minHeight: '100vh' }} />

  const handleSubmit = async () => {
    if (state.tab === 'signup' && !state.name.trim()) { dispatch({ type: 'SET_ERROR', error: 'Name is required' }); return }
    if (!state.username.trim() || !state.password) { dispatch({ type: 'SET_ERROR', error: 'Username and password are required' }); return }
    dispatch({ type: 'SET_ERROR', error: '' })
    dispatch({ type: 'SET_LOADING', loading: true })
    try {
      const res = state.tab === 'signup'
        ? await signup(state.name.trim(), state.username.trim(), state.password)
        : await login(state.username.trim(), state.password)
      setToken(res.token)
      setUserInfo({ student_id: res.student_id, name: res.name, username: state.username.trim(), concepts_mastered: 0, current_streak: 0, level: 'Novice', diagnostic_completed: res.diagnostic_completed })
      push(res.diagnostic_completed ? '/session' : '/onboard')
    } catch (e: any) {
      dispatch({ type: 'SET_ERROR', error: e.message || 'Authentication failed' })
    } finally { dispatch({ type: 'SET_LOADING', loading: false }) }
  }

  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-6 max-sm:px-4">
        <section className="pt-8 max-w-md mx-auto mt-12">
          <span className="flex mb-4">
            <Link href="/" className="text-mathua-secondary text-sm hover:text-mathua-primary">Back</Link>
          </span>
          <SectionHeader label="Account" title={state.tab === 'login' ? 'Welcome back' : 'Create account'} />
          <div className="flex gap-2 mt-6 mb-4">
            <button onClick={() => dispatch({ type: 'SET_TAB', tab: 'login' })} className={`flex-1 rounded-none h-12 text-sm font-medium ${state.tab === 'login' ? 'border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white' : 'bg-mathua-surface-elevated border border-mathua-border text-mathua-secondary'}`}>Login</button>
            <button onClick={() => dispatch({ type: 'SET_TAB', tab: 'signup' })} className={`flex-1 rounded-none h-12 text-sm font-medium ${state.tab === 'signup' ? 'border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white' : 'bg-mathua-surface-elevated border border-mathua-border text-mathua-secondary'}`}>Sign Up</button>
          </div>
          <div className="bg-mathua-surface border border-mathua-border rounded-none p-6 space-y-4">
            {state.tab === 'signup' && (
              <div>
                <label htmlFor="name" className="font-mono text-[10px] uppercase text-mathua-muted">Name</label>
                <input id="name" type="text" value={state.name} onChange={(e) => dispatch({ type: 'SET_NAME', name: e.target.value })} onKeyDown={(e) => e.key === 'Enter' && handleSubmit()} placeholder="Your name" className="w-full mt-1 bg-mathua-code border border-mathua-border rounded-none h-12 px-4 font-mono text-base text-mathua-primary placeholder:text-mathua-muted focus:outline-none focus:border-mathua-blue" />
              </div>
            )}
            <div>
              <label htmlFor="username" className="font-mono text-[10px] uppercase text-mathua-muted">Username</label>
              <input id="username" type="text" value={state.username} onChange={(e) => dispatch({ type: 'SET_USERNAME', username: e.target.value })} onKeyDown={(e) => e.key === 'Enter' && handleSubmit()} placeholder="username" className="w-full mt-1 bg-mathua-code border border-mathua-border rounded-none h-12 px-4 font-mono text-base text-mathua-primary placeholder:text-mathua-muted focus:outline-none focus:border-mathua-blue" />
            </div>
            <div>
              <label htmlFor="password" className="font-mono text-[10px] uppercase text-mathua-muted">Password</label>
              <input id="password" type="password" value={state.password} onChange={(e) => dispatch({ type: 'SET_PASSWORD', password: e.target.value })} onKeyDown={(e) => e.key === 'Enter' && handleSubmit()} placeholder="password" className="w-full mt-1 bg-mathua-code border border-mathua-border rounded-none h-12 px-4 font-mono text-base text-mathua-primary placeholder:text-mathua-muted focus:outline-none focus:border-mathua-blue" />
            </div>
            {state.error && <p className="text-mathua-red text-xs">{state.error}</p>}
            <button onClick={handleSubmit} disabled={state.loading} className="w-full border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none h-12 font-medium text-sm disabled:opacity-50">
              {state.loading ? 'Loading\u2026' : state.tab === 'signup' ? 'Create Account' : 'Login'}
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
    </>
  )
}
