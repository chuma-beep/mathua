'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useTheme } from '../../hooks/useTheme'
import Header from '../../components/Header'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import { isLoggedIn } from '../../lib/auth'

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

interface Course {
  id: string; name: string; grade: string; description: string; targets: string[]
}

export default function GoalsPage() {
  const { mounted } = useTheme()
  const router = useRouter()
  const [courses, setCourses] = useState<Course[]>([])
  const [loading, setLoading] = useState(true)
  const [setting, setSetting] = useState<string | null>(null)
  const [error, setError] = useState('')

  useEffect(() => {
    if (!isLoggedIn()) { router.push('/login'); return }
    fetch(`${API_BASE}/api/courses`, {
      headers: { Authorization: `Bearer ${localStorage.getItem('mathua_token')}` }
    })
      .then(r => r.json())
      .then(d => { setCourses(d.courses || []); setLoading(false) })
      .catch(() => setLoading(false))
  }, [])

  const startCourse = async (courseID: string) => {
    setSetting(courseID)
    setError('')
    try {
      const res = await fetch(`${API_BASE}/api/courses/${courseID}/diagnostic`, {
        method: 'POST',
        headers: { Authorization: `Bearer ${localStorage.getItem('mathua_token')}` }
      })
      if (!res.ok) throw new Error('Failed')
      router.push('/session')
    } catch {
      setError('Could not start course')
    } finally {
      setSetting(null)
    }
  }

  if (!mounted) return <div style={{ background: 'var(--bg)', minHeight: '100vh' }} />

  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-6 max-sm:px-4">
        <section className="pt-8">
          <span className="flex mb-4">
            <a href="/" className="text-mathua-secondary text-sm hover:text-mathua-primary">Back</a>
          </span>
          <SectionHeader label="Choose your course" title="What do you want to learn?" />
          <p className="text-mathua-secondary text-sm text-center mt-2 mb-8">
            Pick a course below. The diagnostic will assess your current level and create a personalised learning path.
          </p>
          {loading && <p className="text-center text-mathua-muted text-sm">Loading courses...</p>}
          {error && <p className="text-center text-mathua-red text-sm mb-4">{error}</p>}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 max-w-3xl mx-auto">
            {courses.map(c => (
              <div key={c.id} className="bg-mathua-surface border border-mathua-border rounded-lg p-6 flex flex-col">
                <div className="flex justify-between items-start mb-2">
                  <h3 className="text-mathua-primary font-medium">{c.name}</h3>
                  <span className="font-mono text-[11px] text-mathua-muted">{c.grade}</span>
                </div>
                <p className="text-mathua-secondary text-sm flex-1 mb-4">{c.description}</p>
                <button
                  onClick={() => startCourse(c.id)}
                  disabled={setting !== null}
                  className="w-full bg-mathua-blue text-white hover:bg-mathua-blue-hover rounded-md h-10 font-medium text-sm disabled:opacity-50"
                >
                  {setting === c.id ? 'Starting...' : 'Start Diagnostic'}
                </button>
              </div>
            ))}
          </div>
        </section>
      </div>
      <Footer />
    </>
  )
}
