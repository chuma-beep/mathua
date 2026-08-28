'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import AsciiDivider from '../../components/AsciiDivider'
import { getSettings, updateSettings, type UserSettings } from '../../lib/api'
import { isLoggedIn } from '../../lib/auth'
import { Switch } from '../../components/ui/switch'
import Loading from '../../components/Loading'

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

  const handleCheckedChange = async (checked: boolean) => {
    const next = { ...settings, show_timer: checked }
    setSettings(next)
    setSaved(false)
    try {
      await updateSettings(next)
      setSaved(true)
      setTimeout(() => setSaved(false), 2000)
    } catch {
      // Keep optimistic state so the toggle still works when the API is
      // unreachable (e.g. static preview). Persist will apply on next save.
      console.error('updateSettings failed')
    }
  }

  if (loading) {
    return (
      <>
        <Header />
        <div className="max-w-container mx-auto px-4 sm:px-6 pt-20 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 text-center overflow-x-hidden min-w-0">
          <Loading label="LOADING SETTINGS" />
        </div>
        <BottomTabs />
        <Footer />
      </>
    )
  }

  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-4 sm:px-6 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0">
        <section className="pt-8 min-w-0 overflow-hidden">
          <span className="flex justify-between items-center mb-4">
            <Link href="/" className="text-mathua-secondary text-sm hover:text-mathua-primary">
              ← Back
            </Link>
          </span>

          <div className="max-w-lg mx-auto mt-8 sm:mt-12 w-full max-w-full min-w-0 px-2 sm:px-0">
            <SectionHeader label="Settings" title="Preferences" />
            <p className="text-mathua-secondary text-sm text-center mt-2 mb-8">
              Customize your learning experience.
            </p>

            <div className="bg-mathua-surface border border-mathua-border rounded-none p-4 sm:p-6 space-y-6 w-full max-w-full min-w-0 overflow-hidden">
              <div className="flex items-center justify-between gap-4 min-w-0">
                <div className="min-w-0 flex-1">
                  <span className="font-mono text-sm text-mathua-primary">Show answer timer</span>
                  <p className="text-mathua-muted text-xs mt-1">
                    Display elapsed time while answering questions.
                  </p>
                </div>
                <Switch
                  checked={!!settings.show_timer}
                  onCheckedChange={handleCheckedChange}
                  aria-label="Toggle answer timer"
                />
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
      <BottomTabs />
    </>
  )
}
