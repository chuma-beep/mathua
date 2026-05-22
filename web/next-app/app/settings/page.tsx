'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import Header from '../../components/Header'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import AsciiDivider from '../../components/AsciiDivider'
import { getSettings, updateSettings, type UserSettings } from '../../lib/api'
import { isLoggedIn } from '../../lib/auth'

export default function SettingsPage() {
  const { push } = useRouter()
  const [settings, setSettings] = useState<UserSettings>({})
  const [saved, setSaved] = useState(false)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!isLoggedIn()) {
      push('/login')
      return
    }
    getSettings().then(s => {
      setSettings(s)
      setLoading(false)
    }).catch((e) => {
      console.error('getSettings failed:', e)
      setLoading(false)
    })
  }, [push])

  const toggleTimer = async () => {
    const next = { ...settings, show_timer: !settings.show_timer }
    setSettings(next)
    setSaved(false)
    try {
      await updateSettings(next)
      setSaved(true)
      setTimeout(() => setSaved(false), 2000)
    } catch {
      // revert
      setSettings(settings)
    }
  }

  if (loading) {
    return (
      <>
        <Header />
        <div className="max-w-container mx-auto px-6 max-sm:px-4 pt-20 text-center">
          <p className="text-mathua-muted text-sm">Loading…</p>
        </div>
        <Footer />
      </>
    )
  }

  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-6 max-sm:px-4">
        <section className="pt-8">
          <span className="flex justify-between items-center mb-4">
            <Link href="/" className="text-mathua-secondary text-sm hover:text-mathua-primary">
              ← Back
            </Link>
          </span>

          <div className="max-w-lg mx-auto mt-12">
            <SectionHeader label="Settings" title="Preferences" />
            <p className="text-mathua-secondary text-sm text-center mt-2 mb-8">
              Customize your learning experience.
            </p>

            <div className="bg-mathua-surface border border-mathua-border rounded-none p-6 space-y-6">
              <div className="flex items-center justify-between">
                <div>
                  <span className="font-mono text-sm text-mathua-primary">Show answer timer</span>
                  <p className="text-mathua-muted text-xs mt-1">
                    Display elapsed time while answering questions.
                  </p>
                </div>
                <button
                  onClick={toggleTimer}
                  className={`relative w-12 h-6 rounded-full transition-colors ${
                    settings.show_timer ? 'bg-mathua-blue' : 'bg-mathua-border-strong'
                  }`}
                >
                  <span
                    className={`absolute top-0.5 left-0.5 w-5 h-5 bg-white rounded-full transition-transform ${
                      settings.show_timer ? 'translate-x-6' : 'translate-x-0'
                    }`}
                  />
                </button>
              </div>
            </div>

            {saved && (
              <p className="text-mathua-green text-xs text-center mt-4 font-mono">
                Settings saved.
              </p>
            )}
          </div>
        </section>

        <AsciiDivider pattern="wave" />
        <Footer />
      </div>
    </>
  )
}
