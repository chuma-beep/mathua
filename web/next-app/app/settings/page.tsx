'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import AsciiDivider from '../../components/AsciiDivider'
import { getSettings, updateSettings, updateProfileName, changePassword, uploadAvatar, deleteAvatar, avatarImageUrl, getIdentities, deleteIdentity, createLinkToken, requestEmailVerification, startOAuthLogin, getConfig, OAUTH_LABELS, type OAuthProvider, enableShare, disableShare, type UserSettings } from '../../lib/api'
import { isLoggedIn, getUserInfo, setUserInfo } from '../../lib/auth'
import { DICEBEAR_STYLES, dicebearUrl, randomDicebear, type DicebearPick } from '../../lib/dicebear'
import { Switch } from '../../components/ui/switch'
import { Dices } from 'lucide-react'
import Loading from '../../components/Loading'
import Avatar from '../../components/Avatar'

export default function SettingsPage() {
  const { push } = useRouter()
  const [settings, setSettings] = useState<UserSettings>({})
  const [saved, setSaved] = useState(false)
  const [loading, setLoading] = useState(true)
  const [shareToken, setShareToken] = useState('')
  const [shareUrl, setShareUrl] = useState('')
  const [shareBusy, setShareBusy] = useState(false)
  const [restricted, setRestricted] = useState(false)
  const [copied, setCopied] = useState(false)
  const [saveError, setSaveError] = useState('')
  const [displayName, setDisplayName] = useState('')
  const [nameBusy, setNameBusy] = useState(false)
  const [nameMsg, setNameMsg] = useState('')
  const [photoBusy, setPhotoBusy] = useState(false)
  const [photoMsg, setPhotoMsg] = useState('')
  const [photoVersion, setPhotoVersion] = useState<number | null>(null)
  const [pendingDice, setPendingDice] = useState<{ style: DicebearPick['style']; seed: string } | null>(null)
  const [imageBusy, setImageBusy] = useState(false)
  const [imageMsg, setImageMsg] = useState('')
  const [recoveryEmail, setRecoveryEmail] = useState('')
  const [cpCurrent, setCpCurrent] = useState('')
  const [cpNew, setCpNew] = useState('')
  const [cpBusy, setCpBusy] = useState(false)
  const [cpMsg, setCpMsg] = useState('')
  const [identities, setIdentities] = useState<{ provider: string; email: string }[]>([])
  const [idBusy, setIdBusy] = useState(false)
  const [idMsg, setIdMsg] = useState('')
  const [providers, setProviders] = useState<OAuthProvider[]>([])
  const [verifyMsg, setVerifyMsg] = useState('')
  const [verifyBusy, setVerifyBusy] = useState(false)

  useEffect(() => {
    if (!isLoggedIn()) {
      setRestricted(true)
      setLoading(false)
      return
    }
    getSettings().then(s => {
      setSettings(s)
      setLoading(false)
    }).catch((e) => {
      console.error('getSettings failed:', e)
      setLoading(false)
    })
    getIdentities().then(setIdentities).catch(() => {})
    getConfig().then(cfg => {
      if (cfg && Array.isArray(cfg.providers)) {
        setProviders(cfg.providers.filter((p): p is OAuthProvider => p in OAUTH_LABELS))
      }
    }).catch(() => {})
    setDisplayName(getUserInfo()?.name ?? '')
    setRecoveryEmail(getUserInfo()?.email ?? '')
    try {
      const savedShare = localStorage.getItem('mathua_share_url')
      if (savedShare) setShareUrl(savedShare)
    } catch { /* private mode — share link just won't persist */ }
  }, [])

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
      setSaveError('Couldn\u2019t save — check your connection and try again.')
      setTimeout(() => setSaveError(''), 3000)
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
  const pausedLabel = (() => {
    if (!settings.pause_until) return ''
    const d = new Date(settings.pause_until + 'T00:00:00Z')
    return isNaN(d.getTime()) ? settings.pause_until : d.toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric' })
  })()

  const handleShare = async (enable: boolean) => {
    setShareBusy(true)
    try {
      if (enable) {
        const res = await enableShare()
        setShareToken(res.token)
        const url = `${window.location.origin}/share?token=${res.token}`
        setShareUrl(url)
        try { localStorage.setItem('mathua_share_url', url) } catch { /* ignore */ }
      } else {
        await disableShare()
        setShareToken('')
        setShareUrl('')
        try { localStorage.removeItem('mathua_share_url') } catch { /* ignore */ }
      }
    } catch {
      console.error('share toggle failed')
    } finally {
      setShareBusy(false)
    }
  }

  // Profile helpers: every save strips the legacy avatar_preset key (quiet
  // migration to the DiceBear model — presets are no longer offered).
  const stripPreset = (s: UserSettings): UserSettings => {
    const next = { ...s }
    delete (next as unknown as Record<string, unknown>).avatar_preset
    return next
  }

  const gallerySeed = (typeof window !== 'undefined' ? getUserInfo()?.student_id : undefined) ?? 'guest'

  // One effective seed everywhere: a staged surprise re-seeds the whole
  // grid (new character, different outfits); otherwise the saved pick's
  // seed; otherwise the stable per-user seed. After Save, pending clears
  // and the grid falls back to the just-saved seed — no visual jump.
  const gridSeed = pendingDice?.seed ?? settings.avatar_dicebear?.seed ?? gallerySeed

  const handleSaveName = async () => {
    const name = displayName.trim()
    if (!name) {
      setNameMsg('Name can’t be empty.')
      return
    }
    setNameBusy(true)
    setNameMsg('')
    try {
      const res = await updateProfileName(name, recoveryEmail.trim())
      const info = getUserInfo()
      if (info) setUserInfo({ ...info, name: res.name, email: res.email || info.email })
      setDisplayName(res.name)
      setNameMsg('Profile saved.')
    } catch (e) {
      // Surface friendly server messages (e.g. "email already in use");
      // fall back to the generic notice for network failures.
      const serverMessage = (e as { serverMessage?: string })?.serverMessage
      setNameMsg(serverMessage || 'Couldn’t save — check your connection and try again.')
    } finally {
      setNameBusy(false)
    }
  }

  // Gallery picks and Surprise only stage a preview — Save image persists.
  const handlePickStyle = (style: DicebearPick['style']) => {
    const seed = pendingDice?.seed ?? settings.avatar_dicebear?.seed ?? gallerySeed
    setPendingDice({ style, seed })
    setImageMsg('')
  }

  const handleSurprise = () => {
    const pick = randomDicebear()
    setPendingDice({ style: pick.style, seed: pick.seed })
    setImageMsg('')
  }

  const pendingDiffers =
    pendingDice !== null &&
    (pendingDice.style !== settings.avatar_dicebear?.style ||
      pendingDice.seed !== settings.avatar_dicebear?.seed ||
      settings.avatar_custom === true)

  const handleSaveImage = async () => {
    if (!pendingDice) return
    setImageBusy(true)
    setImageMsg('')
    try {
      const next = stripPreset({ ...settings, avatar_custom: false, avatar_dicebear: { style: pendingDice.style, seed: pendingDice.seed } })
      setSettings(next)
      await updateSettings(next)
      if (settings.avatar_custom) await deleteAvatar().catch(() => {})
      setPhotoVersion(null)
      setPendingDice(null)
      setImageMsg('Image saved.')
      // Broadcast so the Header picture refreshes in-tab (same object —
      // setUserInfo always dispatches auth-changed).
      const savedInfo = getUserInfo()
      if (savedInfo) setUserInfo({ ...savedInfo })
    } catch {
      setImageMsg('Couldn’t save — check your connection and try again.')
    } finally {
      setImageBusy(false)
    }
  }

  const handleUpload = async (file: File | undefined) => {
    if (!file) return
    if (file.size > 512 * 1024) {
      setPhotoMsg('Photo must be 512KB or smaller.')
      return
    }
    setPhotoBusy(true)
    setPhotoMsg('')
    try {
      await uploadAvatar(file)
      const next = stripPreset({ ...settings, avatar_custom: true, avatar_version: Date.now() })
      setSettings(next)
      setPendingDice(null)
      setPhotoVersion(Date.now())
      setPhotoMsg('Photo updated.')
      const uploadedInfo = getUserInfo()
      if (uploadedInfo) setUserInfo({ ...uploadedInfo })
    } catch {
      setPhotoMsg('Upload failed — use PNG, JPEG, GIF, or WebP up to 512KB.')
    } finally {
      setPhotoBusy(false)
    }
  }

  const handleRemovePhoto = async () => {
    setPhotoBusy(true)
    try {
      await deleteAvatar()
      const next = stripPreset({ ...settings, avatar_custom: false, avatar_version: Date.now() })
      setSettings(next)
      setPendingDice(null)
      setPhotoVersion(null)
      const removedInfo = getUserInfo()
      if (removedInfo) setUserInfo({ ...removedInfo })
    } catch { console.error('deleteAvatar failed') }
    finally { setPhotoBusy(false) }
  }

  const handleChangePassword = async () => {
    if (!cpCurrent || !cpNew) {
      setCpMsg('Enter your current and a new password.')
      return
    }
    setCpBusy(true)
    setCpMsg('')
    try {
      await changePassword(cpCurrent, cpNew)
      setCpCurrent('')
      setCpNew('')
      setCpMsg('Password changed.')
    } catch (e: unknown) {
      setCpMsg(e instanceof Error ? e.message : 'Change failed — try again.')
    } finally {
      setCpBusy(false)
    }
  }

  const currentPhotoUrl = settings.avatar_custom ? avatarImageUrl(photoVersion ?? undefined) : undefined
  const currentDice = settings.avatar_dicebear

  const refreshIdentities = () => {
    getIdentities().then(setIdentities).catch(() => {})
  }

  const handleVerifyEmail = async () => {
    setVerifyBusy(true)
    setVerifyMsg('')
    try {
      await requestEmailVerification()
      setVerifyMsg('Verification link sent — check your inbox.')
      const info = getUserInfo()
      if (info) setUserInfo({ ...info })
    } catch (e: unknown) {
      setVerifyMsg(e instanceof Error ? e.message : 'Request failed — try again.')
    } finally {
      setVerifyBusy(false)
    }
  }

  const handleConnect = async (provider: string) => {
    setIdBusy(true)
    setIdMsg('')
    try {
      const linkToken = await createLinkToken()
      startOAuthLogin(provider, { intent: 'link', linkToken })
    } catch {
      setIdMsg('Couldn’t start connecting — try again.')
      setIdBusy(false)
    }
  }

  const handleDisconnect = async (provider: string) => {
    setIdBusy(true)
    setIdMsg('')
    try {
      await deleteIdentity(provider)
      refreshIdentities()
    } catch (e: unknown) {
      setIdMsg(e instanceof Error ? e.message : 'Disconnect failed — try again.')
    } finally {
      setIdBusy(false)
    }
  }

  // Single-pass lookup for the connect-button list below.
  const connectedSet = new Set(identities.map(id => id.provider))

  const userEmail = typeof window !== 'undefined' ? getUserInfo()?.email ?? '' : ''
  const userEmailVerified = typeof window !== 'undefined' ? getUserInfo()?.email_verified ?? false : false

  if (restricted) {
    return (
      <>
        <Header />
        <div className="max-w-container mx-auto px-4 sm:px-6 py-20 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-20 text-center overflow-x-hidden min-w-0">
          <p className="font-mono text-xs text-mathua-secondary mb-4">Settings needs an account: your preferences are stored per account.</p>
          <Link href="/login" className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white px-6 py-2 font-mono text-xs min-h-[36px] inline-flex items-center justify-center">Sign in</Link>
        </div>
        <BottomTabs />
        <Footer />
      </>
    )
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
            <button type="button" onClick={() => { if (window.history.length > 1) window.history.back(); else push('/profile') }} className="text-mathua-secondary text-sm hover:text-mathua-primary">
              ← Back
            </button>
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
                    ? `Paused until ${pausedLabel}. Due reviews are hidden while paused. Pick another duration to change it, or resume.`
                    : 'Take a break for 30, 60, or 90 days. Due reviews are hidden while paused.'}
                </p>
                <div className="flex flex-wrap gap-2 mt-3">
                  {[30, 60, 90].map(d => (
                    <button
                      type="button"
                      key={d}
                      onClick={() => handlePause(d)}
                      className="border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue rounded-none px-4 h-9 text-xs font-mono disabled:opacity-40"
                    >
                      Pause {d} days
                    </button>
                  ))}
                  {paused && (
                    <button
                      type="button"
                      onClick={() => handlePause(null)}
                      className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none px-4 h-9 text-xs font-mono"
                    >
                      Resume now
                    </button>
                  )}
                </div>
              </div>

              <div className="border-t border-mathua-border pt-6 min-w-0">
                <span className="font-mono text-sm text-mathua-primary">Profile</span>
                <p className="text-mathua-muted text-xs mt-1">
                  Your display name and profile picture. Username stays the same.
                </p>
                <label htmlFor="display-name" className="font-mono text-[11px] uppercase text-mathua-muted mt-4 block">Display name</label>
                <form onSubmit={e => { e.preventDefault(); handleSaveName() }} className="flex flex-col sm:flex-row gap-2 mt-2">
                  <input
                    id="display-name"
                    type="text"
                    value={displayName}
                    maxLength={50}
                    onChange={(e) => setDisplayName(e.target.value)}
                    placeholder="Your name"
                    enterKeyHint="done"
                    className="flex-1 min-w-0 bg-mathua-code border border-mathua-border rounded-none h-10 px-3 font-mono text-sm text-mathua-primary placeholder:text-mathua-muted focus:outline-none focus:border-mathua-blue"
                  />
                  <button
                    type="submit"
                    disabled={nameBusy}
                    className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none px-4 h-10 text-xs font-mono disabled:opacity-50 shrink-0 inline-flex items-center gap-2"
                  >
                    {nameBusy ? (<><Loading inline size={11} /> Saving…</>) : 'Save profile'}
                  </button>
                </form>
                {nameMsg && <p className="font-mono text-[11px] text-mathua-secondary mt-2">{nameMsg}</p>}

                <label htmlFor="recovery-email" className="font-mono text-[11px] uppercase text-mathua-muted mt-4 block">Recovery email (optional)</label>
                <p className="text-mathua-muted text-xs mt-1">Used only for password reset. Empty clears it.</p>
                <div className="flex flex-col sm:flex-row gap-2 mt-2">
                  <input
                    id="recovery-email"
                    type="email"
                    value={recoveryEmail}
                    onChange={(e) => setRecoveryEmail(e.target.value)}
                    placeholder="you@example.com"
                    autoComplete="email"
                    className="flex-1 min-w-0 bg-mathua-code border border-mathua-border rounded-none h-10 px-3 font-mono text-sm text-mathua-primary placeholder:text-mathua-muted focus:outline-none focus:border-mathua-blue"
                  />
                </div>
                {userEmail !== '' && !userEmailVerified && (
                  <div className="mt-2 flex flex-wrap items-center gap-2">
                    <span className="font-mono text-[11px] text-mathua-muted">Email unverified: logins won’t merge until you verify.</span>
                    <button
                      type="button"
                      onClick={handleVerifyEmail}
                      disabled={verifyBusy}
                      className="border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue rounded-none px-3 h-8 text-[11px] font-mono disabled:opacity-50"
                    >
                      {verifyBusy ? 'Sending…' : 'Verify email'}
                    </button>
                  </div>
                )}
                {verifyMsg && <p className="font-mono text-[11px] text-mathua-secondary mt-2">{verifyMsg}</p>}

                <div className="flex items-center gap-4 mt-6">
                  <Avatar
                    seed={gridSeed}
                    name={displayName || 'You'}
                    size={56}
                    url={pendingDice ? dicebearUrl(pendingDice.style, pendingDice.seed) : currentPhotoUrl ?? (currentDice ? dicebearUrl(currentDice.style as DicebearPick['style'], currentDice.seed) : undefined)}
                  />
                  <div className="min-w-0">
                    <p className="font-mono text-[11px] uppercase text-mathua-muted">Current picture</p>
                    <p className="text-mathua-muted text-xs mt-1 break-words">
                      {pendingDice
                        ? 'New preview — not saved yet.'
                        : settings.avatar_custom ? 'Your uploaded photo.' : currentDice ? 'Your DiceBear character.' : 'Automatic — pick a character or upload a photo.'}
                    </p>
                  </div>
                </div>

                <p className="font-mono text-[11px] uppercase text-mathua-muted mt-6 mb-2">Pick a character</p>
                <div className="grid grid-cols-4 sm:grid-cols-6 gap-2">
                  {DICEBEAR_STYLES.map((style) => {
                    const shown = pendingDice ?? currentDice
                    const selected = !settings.avatar_custom && shown?.style === style
                    return (
                      <button
                        type="button"
                        key={style}
                        onClick={() => handlePickStyle(style)}
                        title={style}
                        aria-label={`Pick ${style} character`}
                        className={`p-1 border rounded-none transition-colors bg-mathua-code ${selected ? 'border-mathua-blue' : 'border-transparent hover:border-mathua-border'}`}
                      >
                        {/* eslint-disable-next-line @next/next/no-img-element */}
                        <img src={dicebearUrl(style, gridSeed)} alt={`${style} avatar character`} width={64} height={64} loading="lazy" className="w-full h-auto" />
                      </button>
                    )
                  })}
                </div>
                <div className="flex flex-wrap gap-2 mt-3">
                  <button
                    type="button"
                    onClick={handleSurprise}
                    className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none px-4 h-9 text-xs font-mono inline-flex items-center gap-2"
                  >
                    <Dices className="size-4" aria-hidden />
                    Surprise me
                  </button>
                  <button
                    type="button"
                    onClick={handleSaveImage}
                    disabled={!pendingDiffers || imageBusy}
                    className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none px-4 h-9 text-xs font-mono disabled:opacity-50 inline-flex items-center gap-2"
                  >
                    {imageBusy ? (<><Loading inline size={11} /> Saving…</>) : 'Save image'}
                  </button>
                  <label className="border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue rounded-none px-4 h-9 text-xs font-mono inline-flex items-center gap-2 cursor-pointer focus-within:border-mathua-blue focus-within:text-mathua-blue">
                    {photoBusy ? (<><Loading inline size={11} /> Uploading…</>) : 'Upload photo'}
                    <input
                      type="file"
                      accept="image/png,image/jpeg,image/gif,image/webp"
                      aria-label="Upload profile photo"
                      aria-describedby="photo-limits"
                      className="hidden"
                      disabled={photoBusy}
                      onChange={(e) => { handleUpload(e.target.files?.[0]); e.target.value = '' }}
                    />
                  </label>
                  {settings.avatar_custom && (
                    <button
                      type="button"
                      onClick={handleRemovePhoto}
                      disabled={photoBusy}
                      className="border border-mathua-border text-mathua-secondary hover:border-mathua-red hover:text-mathua-red rounded-none px-4 h-9 text-xs font-mono disabled:opacity-50"
                    >
                      Remove photo
                    </button>
                  )}
                </div>
                {photoMsg && <p className="font-mono text-[11px] text-mathua-secondary mt-2">{photoMsg}</p>}
                {imageMsg && <p className="font-mono text-[11px] text-mathua-secondary mt-2">{imageMsg}</p>}
                <p id="photo-limits" className="text-mathua-muted text-[11px] mt-3">
                  Characters by <a href="https://www.dicebear.com" target="_blank" rel="noreferrer" className="text-mathua-blue hover:text-mathua-blue-hover">DiceBear</a>. Photos up to 512KB.
                </p>
              </div>

              <div className="border-t border-mathua-border pt-6 min-w-0">
                <span className="font-mono text-sm text-mathua-primary">Change password</span>
                <p className="text-mathua-muted text-xs mt-1">
                  Needs your current password. Forgot it? Sign out and use “Forgot password?” on the login page.
                </p>
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 mt-3">
                  <input
                    type="password"
                    value={cpCurrent}
                    onChange={(e) => setCpCurrent(e.target.value)}
                    placeholder="Current password"
                    autoComplete="current-password"
                    aria-label="Current password"
                    className="min-w-0 bg-mathua-code border border-mathua-border rounded-none h-10 px-3 font-mono text-sm text-mathua-primary placeholder:text-mathua-muted focus:outline-none focus:border-mathua-blue"
                  />
                  <input
                    type="password"
                    value={cpNew}
                    onChange={(e) => setCpNew(e.target.value)}
                    placeholder="New password"
                    autoComplete="new-password"
                    aria-label="New password"
                    className="min-w-0 bg-mathua-code border border-mathua-border rounded-none h-10 px-3 font-mono text-sm text-mathua-primary placeholder:text-mathua-muted focus:outline-none focus:border-mathua-blue"
                  />
                </div>
                <button
                  type="button"
                  onClick={handleChangePassword}
                  disabled={cpBusy}
                  className="mt-3 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none px-4 h-10 text-xs font-mono disabled:opacity-50 inline-flex items-center gap-2"
                >
                  {cpBusy ? (<><Loading inline size={11} /> Saving…</>) : 'Change password'}
                </button>
                {cpMsg && <p className="font-mono text-[11px] text-mathua-secondary mt-2">{cpMsg}</p>}
              </div>

              <div className="border-t border-mathua-border pt-6 min-w-0">
                <span className="font-mono text-sm text-mathua-primary">Connected accounts</span>
                <p className="text-mathua-muted text-xs mt-1">
                  Sign in with any of these: verified emails merge into this account, never a new one.
                </p>
                <div className="mt-3 space-y-2">
                  {identities.map(id => (
                    <div key={id.provider} className="flex items-center justify-between gap-3 border border-mathua-border px-3 h-10 min-w-0">
                      <span className="font-mono text-xs text-mathua-primary truncate">
                        {OAUTH_LABELS[id.provider as OAuthProvider] ?? id.provider}
                        {id.email && <span className="text-mathua-muted"> · {id.email}</span>}
                      </span>
                      <button
                        type="button"
                        onClick={() => handleDisconnect(id.provider)}
                        disabled={idBusy}
                        className="shrink-0 font-mono text-[11px] text-mathua-muted hover:text-mathua-red disabled:opacity-50"
                      >
                        Disconnect
                      </button>
                    </div>
                  ))}
                  {providers.flatMap(p => connectedSet.has(p) ? [] : [
                    <button
                      type="button"
                      key={p}
                      onClick={() => handleConnect(p)}
                      disabled={idBusy}
                      className="w-full border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue rounded-none px-4 h-10 text-xs font-mono disabled:opacity-50"
                    >
                      Connect {OAUTH_LABELS[p]}
                    </button>
                  ])}
                  {providers.length === 0 && (
                    <p className="font-mono text-[11px] text-mathua-muted">No external login services configured on this server.</p>
                  )}
                </div>
                {idMsg && <p className="font-mono text-[11px] text-mathua-secondary mt-2">{idMsg}</p>}
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
                      aria-label="Share link"
                      onFocus={e => e.currentTarget.select()}
                      className="w-full bg-mathua-code border border-mathua-border rounded-none h-10 px-3 font-mono text-xs text-mathua-primary focus:outline-none focus:border-mathua-blue min-w-0"
                    />
                    <div className="flex gap-2 mt-2">
                      <button
                        type="button"
                        onClick={() => {
                          const done = () => { setCopied(true); setTimeout(() => setCopied(false), 2000) }
                          if (navigator.clipboard?.writeText) {
                            navigator.clipboard.writeText(shareUrl).then(done).catch(() => done())
                          } else { done() }
                        }}
                        className="border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue rounded-none px-4 h-9 text-xs font-mono"
                      >
                        {copied ? 'Copied!' : 'Copy link'}
                      </button>
                      <button
                        type="button"
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
                    type="button"
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
            {saveError && (
              <p className="text-mathua-red text-xs text-center mt-4 font-mono">
                {saveError}
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
