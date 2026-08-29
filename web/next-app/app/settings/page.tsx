'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import AsciiDivider from '../../components/AsciiDivider'
import { getSettings, updateSettings, enableShare, disableShare, type UserSettings } from '../../lib/api'
import { isLoggedIn } from '../../lib/auth'
import { Switch } from '../../components/ui/switch'
import Loading from '../../components/Loading'

export default function SettingsPage() {
  const { push } = useRouter()
  const [settings, setSettings] = useState<UserSettings>({})
  const [saved, setSaved] = useState(false)
  const [loading, setLoading] = useState(true)
  const [shareToken, setShareToken] = useState('')
  const [shareUrl, setShareUrl] = useState('')
  const [shareBusy, setShareBusy] = useState(false)

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

  const handlePause = async (days: number | null) => {
    const next: UserSettings = { ...settings }
    if (days === null) {
      next.pause_until = null
    } else {
      const d = new Date()
      d.setUTCDate(d.getUTCDate() + days)
      next.pause_until = d.toISOString().slice(0, 10)
    }
    setSettings(next)
    setSaved(false)
    try {
      await updateSettings(next)
      setSaved(true)
      setTimeout(() => setSaved(false), 2000)
    } catch {
      console.error('updateSettings failed')
    }
  }

  const paused = !!settings.pause_until

  const handleShare = async (enable: boolean) => {
    setShareBusy(true)
    try {
      if (enable) {
        const res = await enableShare()
        setShareToken(res.token)
        setShareUrl(`${window.location.origin}/share?token=${res.token}`)
      } else {
        await disableShare()
        setShareToken('')
        setShareUrl('')
      }
    } catch {
      console.error('share toggle failed')
    } finally {
      setShareBusy(false)
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

              <div className="border-t border-mathua-border pt-6 min-w-0">
                <span className="font-mono text-sm text-mathua-primary">Pause learning</span>
                <p className="text-mathua-muted text-xs mt-1">
                  {paused
                    ? `Paused until ${settings.pause_until}. Due reviews are hidden while paused.`
                    : 'Take a break for 30, 60, or 90 days. Due reviews are hidden while paused.'}
                </p>
                <div className="flex flex-wrap gap-2 mt-3">
                  {[30, 60, 90].map(d => (
                    <button
                      key={d}
                      onClick={() => handlePause(d)}
                      disabled={paused && settings.pause_until !== null}
                      className="border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue rounded-none px-4 h-9 text-xs font-mono disabled:opacity-40"
                    >
                      Pause {d} days
                    </button>
                  ))}
                  {paused && (
                    <button
                      onClick={() => handlePause(null)}
                      className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none px-4 h-9 text-xs font-mono"
                    >
                      Resume now
                    </button>
                  )}
                </div>
              </div>

              <div className="border-t border-mathua-border pt-6 min-w-0">
                <span className="font-mono text-sm text-mathua-primary">Share with parent / teacher</span>
                <p className="text-mathua-muted text-xs mt-1">
                  Generate a read-only link to this student&apos;s progress, activity, and weak spots.
                </p>
                {shareUrl ? (
                  <div className="mt-3 min-w-0">
                    <input
                      readOnly
                      value={shareUrl}
                      onFocus={e => e.currentTarget.select()}
                      className="w-full bg-mathua-code border border-mathua-border rounded-none h-10 px-3 font-mono text-xs text-mathua-primary focus:outline-none focus:border-mathua-blue min-w-0"
                    />
                    <div className="flex gap-2 mt-2">
                      <button
                        onClick={() => { navigator.clipboard?.writeText(shareUrl) }}
                        className="border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue rounded-none px-4 h-9 text-xs font-mono"
                      >
                        Copy link
                      </button>
                      <button
                        onClick={() => handleShare(false)}
                        disabled={shareBusy}
                        className="border border-mathua-red text-mathua-red hover:bg-mathua-red hover:text-white rounded-none px-4 h-9 text-xs font-mono disabled:opacity-50"
                      >
                        Disable share
                      </button>
                    </div>
                  </div>
                ) : (
                  <button
                    onClick={() => handleShare(true)}
                    disabled={shareBusy}
                    className="mt-3 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none px-4 h-9 text-xs font-mono disabled:opacity-50"
                  >
                    {shareBusy ? 'Generating…' : 'Enable share link'}
                  </button>
                )}
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
