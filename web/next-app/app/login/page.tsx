'use client'

import { useState } from 'react'
import { useTheme } from '../../hooks/useTheme'
import { useRouter } from 'next/navigation'
import Header from '../../components/Header'
import Footer from '../../components/Footer'
import SectionHeader from '../../components/SectionHeader'
import { signup, login } from '../../lib/api'
import { setToken, setUserInfo } from '../../lib/auth'

export default function LoginPage() {
  const { mounted } = useTheme()
  const router = useRouter()
  const [tab, setTab] = useState<'login' | 'signup'>('login')
  const [name, setName] = useState('')
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  if (!mounted) return <div style={{ background: 'var(--bg)', minHeight: '100vh' }} />

  const handleSubmit = async () => {
    if (tab === 'signup' && !name.trim()) { setError('Name is required'); return }
    if (!username.trim() || !password) { setError('Username and password are required'); return }
    setError(''); setLoading(true)
    try {
      const res = tab === 'signup'
        ? await signup(name.trim(), username.trim(), password)
        : await login(username.trim(), password)
      setToken(res.token)
      setUserInfo({ student_id: res.student_id, name: res.name, username: username.trim(), concepts_mastered: 0, current_streak: 0, level: 'Novice' })
      router.push('/goals')
    } catch (e: any) {
      setError(e.message || 'Authentication failed')
    } finally { setLoading(false) }
  }

  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-6 max-sm:px-4">
        <section className="pt-8 max-w-md mx-auto mt-12">
          <span className="flex mb-4">
            <a href="/" className="text-mathua-secondary text-sm hover:text-mathua-primary">Back</a>
          </span>
          <SectionHeader label="Account" title={tab === 'login' ? 'Welcome back' : 'Create account'} />
          <div className="flex gap-2 mt-6 mb-4">
            <button onClick={() => { setTab('login'); setError('') }} className={`flex-1 rounded-md h-10 text-sm font-medium ${tab === 'login' ? 'bg-mathua-blue text-white' : 'bg-mathua-surface-elevated border border-mathua-border text-mathua-secondary'}`}>Login</button>
            <button onClick={() => { setTab('signup'); setError('') }} className={`flex-1 rounded-md h-10 text-sm font-medium ${tab === 'signup' ? 'bg-mathua-blue text-white' : 'bg-mathua-surface-elevated border border-mathua-border text-mathua-secondary'}`}>Sign Up</button>
          </div>
          <div className="bg-mathua-surface border border-mathua-border rounded-lg p-6 space-y-4">
            {tab === 'signup' && (
              <div>
                <label className="font-mono text-[10px] uppercase text-mathua-muted">Name</label>
                <input type="text" value={name} onChange={(e) => setName(e.target.value)} onKeyDown={(e) => e.key === 'Enter' && handleSubmit()} placeholder="Your name" autoFocus className="w-full mt-1 bg-mathua-code border border-mathua-border rounded-md h-10 px-3 font-mono text-sm text-mathua-primary placeholder:text-mathua-muted focus:outline-none focus:border-mathua-blue" />
              </div>
            )}
            <div>
              <label className="font-mono text-[10px] uppercase text-mathua-muted">Username</label>
              <input type="text" value={username} onChange={(e) => setUsername(e.target.value)} onKeyDown={(e) => e.key === 'Enter' && handleSubmit()} placeholder="username" className="w-full mt-1 bg-mathua-code border border-mathua-border rounded-md h-10 px-3 font-mono text-sm text-mathua-primary placeholder:text-mathua-muted focus:outline-none focus:border-mathua-blue" />
            </div>
            <div>
              <label className="font-mono text-[10px] uppercase text-mathua-muted">Password</label>
              <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} onKeyDown={(e) => e.key === 'Enter' && handleSubmit()} placeholder="password" className="w-full mt-1 bg-mathua-code border border-mathua-border rounded-md h-10 px-3 font-mono text-sm text-mathua-primary placeholder:text-mathua-muted focus:outline-none focus:border-mathua-blue" />
            </div>
            {error && <p className="text-mathua-red text-xs">{error}</p>}
            <button onClick={handleSubmit} disabled={loading} className="w-full bg-mathua-blue text-white hover:bg-mathua-blue-hover rounded-md h-12 font-medium text-sm disabled:opacity-50">
              {loading ? 'Loading...' : tab === 'signup' ? 'Create Account' : 'Login'}
            </button>
          </div>
        </section>
      </div>
      <Footer />
    </>
  )
}
