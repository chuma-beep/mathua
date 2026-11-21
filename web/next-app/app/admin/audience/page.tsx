'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import dynamic from 'next/dynamic'
import AdminShell from '@/components/admin/AdminShell'
import SectionHeader from '@/components/SectionHeader'
import { Skeleton } from '@/components/ui/skeleton'
import { Table, TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { useAdminMe } from '@/hooks/useAdminMe'
import {
  fetchAdminLocations,
  getErrorMessage,
  type AdminLocationCountry,
  type AdminLocationDay,
  type AdminLocationDeviceDay,
  type AdminLocations,
} from '../../../lib/api'

// three-globe builds a WebGL scene when its module is evaluated, and that
// touches `window`. The page is still prerendered, so the globe loads only in
// the browser and the server render keeps the card's shape.
const CountryGlobe = dynamic(() => import('@/components/admin/CountryGlobe'), {
  ssr: false,
  loading: () => <div className="aspect-square w-full max-w-[560px] animate-pulse bg-mathua-surface" aria-hidden />,
})

// The device gallery builds its scene the same way, and is only meaningful in a
// browser with WebGL, so it follows the globe behind a dynamic import.
const DeviceGallery = dynamic(() => import('@/components/admin/DeviceGallery'), {
  ssr: false,
  loading: () => (
    <div className="mx-auto h-[218px] w-full max-w-[620px] animate-pulse bg-mathua-surface sm:h-[236px]" aria-hidden />
  ),
})

// Where learners are, at country granularity.
//
// The country is resolved by Cloudflare at the edge and stored once a day for a
// signed-in learner. There is no city, no coordinate and no IP address in the
// answer, so this page cannot become a tracker even for an administrator. The
// honest caveat is in the copy: a VPN or a carrier can put a learner somewhere
// they are not, and a country is a coarse answer to "where" by construction.
//
// Growth is shown as two separate facts, because they are different questions:
// activity growth (distinct learners in the selected window against the window
// before it) and new accounts (signups in the window against the one before it). A country can grow in
// practice without a single new account, and a country can add accounts that
// have not returned yet.

const regionNames = (() => {
  try {
    return new Intl.DisplayNames(['en'], { type: 'region' })
  } catch {
    return null
  }
})()

function countryName(code: string): string {
  try {
    return regionNames?.of(code) ?? code
  } catch {
    return code
  }
}

type Delta = { label: string; tone: 'up' | 'down' | 'flat' }

function delta(current: number, previous: number): Delta {
  if (previous === 0) {
    return current === 0 ? { label: 'no activity yet', tone: 'flat' } : { label: 'first recorded window', tone: 'up' }
  }
  if (current === previous) return { label: 'no change', tone: 'flat' }
  const change = current - previous
  const percent = Math.round((Math.abs(change) / previous) * 100)
  return { label: `${change > 0 ? '+' : '-'}${percent}%`, tone: change > 0 ? 'up' : 'down' }
}

function deltaTone(tone: Delta['tone']): string {
  if (tone === 'up') return 'text-mathua-green'
  if (tone === 'down') return 'text-mathua-muted'
  return 'text-mathua-muted'
}

function growthNote(d: Delta, days: number): string {
  if (d.label === 'first recorded window') return `No previous ${days}-day window to compare`
  if (d.label === 'no activity yet') return 'No activity recorded yet'
  return `${d.label} against the previous ${days} days`
}

function shortDay(day: string): string {
  const date = new Date(`${day}T00:00:00Z`)
  if (Number.isNaN(date.getTime())) return day
  return date.toLocaleDateString(undefined, { month: 'short', day: 'numeric', timeZone: 'UTC' })
}

function Sparkline({ daily }: { daily: AdminLocationDay[] }) {
  const width = 640
  const height = 132
  const pad = 10
  const max = Math.max(1, ...daily.map((d) => d.active))
  const step = daily.length > 1 ? (width - pad * 2) / (daily.length - 1) : 0
  const points = daily.map((d, i) => {
    const x = pad + i * step
    const y = height - pad - (d.active / max) * (height - pad * 2)
    return [x, y] as const
  })
  const line = points.map((p, i) => `${i === 0 ? 'M' : 'L'}${p[0].toFixed(1)} ${p[1].toFixed(1)}`).join(' ')
  const area = `${line} L${width - pad} ${height - pad} L${pad} ${height - pad} Z`
  const peak = daily.length > 0 ? Math.max(...daily.map((d) => d.active)) : 0

  return (
    <div className="flex flex-1 flex-col text-mathua-blue">
      <svg viewBox={`0 0 ${width} ${height}`} className="min-h-32 max-h-48 w-full flex-1 sm:max-h-64" aria-hidden="true" preserveAspectRatio="none">
        <defs>
          <linearGradient id="locations-fill" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stopColor="currentColor" stopOpacity="0.28" />
            <stop offset="100%" stopColor="currentColor" stopOpacity="0" />
          </linearGradient>
        </defs>
        <line x1={pad} y1={height - pad} x2={width - pad} y2={height - pad} stroke="var(--border)" strokeWidth="1" />
        <path d={area} fill="url(#locations-fill)" />
        <path d={line} fill="none" stroke="currentColor" strokeWidth="2" vectorEffect="non-scaling-stroke" />
      </svg>
      <div className="mt-2 flex justify-between font-mono text-[10px] text-mathua-muted">
        <span>{daily.length > 0 ? shortDay(daily[0].day) : ''}</span>
        <span>
          Peak {peak} {peak === 1 ? 'learner' : 'learners'} in a day
        </span>
        <span>{daily.length > 0 ? shortDay(daily[daily.length - 1].day) : ''}</span>
      </div>
      <p className="sr-only">
        Distinct active learners per day over the selected window, ending at {peak} at the peak.
      </p>
    </div>
  )
}

// Three lines, one per device class, over the same window as the activity
// sparkline. The classes overlap: a learner who signs in on a phone and a
// laptop counts on both lines that day, so the lines add up above the daily
// total. The caption says so rather than letting the reader assume a split.
function DeviceLines({ daily }: { daily: AdminLocationDeviceDay[] }) {
  const width = 640
  const height = 132
  const pad = 10
  const max = Math.max(
    1,
    ...daily.flatMap((d) => [d.mobile, d.tablet, d.desktop]),
  )
  const step = daily.length > 1 ? (width - pad * 2) / (daily.length - 1) : 0
  const series = [
    { key: 'mobile', label: 'Mobile', tone: 'text-mathua-blue' },
    { key: 'desktop', label: 'Desktop', tone: 'text-mathua-teal' },
    { key: 'tablet', label: 'Tablet', tone: 'text-mathua-green' },
  ] as const
  const path = (key: 'mobile' | 'desktop' | 'tablet') =>
    daily
      .map((d, i) => {
        const x = pad + i * step
        const y = height - pad - (d[key] / max) * (height - pad * 2)
        return `${i === 0 ? 'M' : 'L'}${x.toFixed(1)} ${y.toFixed(1)}`
      })
      .join(' ')

  return (
    <div>
      <ul className="mb-2 flex flex-wrap items-center gap-x-4 gap-y-1 font-mono text-[10px] text-mathua-muted">
        {series.map((s) => (
          <li key={s.key} className="flex items-center gap-1.5">
            <span aria-hidden="true" className={`h-0.5 w-4 ${s.tone} bg-current`} />
            <span>{s.label}</span>
          </li>
        ))}
      </ul>
      <svg viewBox={`0 0 ${width} ${height}`} className="h-32 w-full" aria-hidden="true" preserveAspectRatio="none">
        <line x1={pad} y1={height - pad} x2={width - pad} y2={height - pad} stroke="var(--border)" strokeWidth="1" />
        {series.map((s) => (
          <path
            key={s.key}
            d={path(s.key)}
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            vectorEffect="non-scaling-stroke"
            className={s.tone}
          />
        ))}
      </svg>
      <div className="mt-2 flex justify-between font-mono text-[10px] text-mathua-muted">
        <span>{daily.length > 0 ? shortDay(daily[0].day) : ''}</span>
        <span>One class per learner per day</span>
        <span>{daily.length > 0 ? shortDay(daily[daily.length - 1].day) : ''}</span>
      </div>
      <p className="sr-only">
        Distinct signed-in learners per day for mobile, desktop and tablet over the selected
        window. A learner who uses two classes on one day appears on both lines.
      </p>
    </div>
  )
}

function StatCard({ label, value, note, tone }: { label: string; value: string; note: string; tone?: Delta['tone'] }) {
  return (
    <div className="border border-mathua-border bg-mathua-surface p-4">
      <p className="font-mono text-[10px] uppercase tracking-[0.15em] text-mathua-muted">{label}</p>
      <p className="mt-2 font-serif text-3xl leading-none text-mathua-primary">{value}</p>
      <p className={`mt-2 font-mono text-[10px] ${tone ? deltaTone(tone) : 'text-mathua-muted'}`}>{note}</p>
    </div>
  )
}

function activityBar(country: AdminLocationCountry, max: number): number {
  if (max <= 0 || country.active <= 0) return 0
  return Math.max(3, Math.round((country.active / max) * 100))
}

function fmt(value: number): string {
  return value.toLocaleString()
}

// The page is one window, not three overlapping ones: the stat cards, the
// country table, the globe, the sparkline and the device series all describe the
// same number of days, and the comparison is always the equally long window
// before it. Ninety days is the widest offer because the daily series is drawn
// in full at every width, and a year of points would collapse to noise.
const RANGES = [7, 30, 90] as const

export default function AdminAudiencePage() {
  const { isStaff } = useAdminMe()
  const [days, setDays] = useState(30)
  const [data, setData] = useState<AdminLocations | null>(null)
  const [error, setError] = useState('')

  useEffect(() => {
    if (!isStaff) return
    let live = true
    fetchAdminLocations(days)
      .then((res) => { if (live) { setData(res); setError('') } })
      .catch((e) => { if (live) setError(getErrorMessage(e)) })
    return () => { live = false }
  }, [isStaff, days])

  // The visible numbers and the labels must agree even while a new window is
  // loading, so labels follow the data that is on screen, not the pending pick.
  const windowDays = data?.days ?? days
  const activeDelta = data ? delta(data.active, data.active_prev) : null
  const newDelta = data ? delta(data.new, data.new_prev) : null
  const visitsDelta = data ? delta(data.visits, data.visits_prev) : null
  const maxActive = data ? Math.max(0, ...data.countries.map((c) => c.active)) : 0
  const movers = data
    ? [...data.countries]
        .filter((c) => c.active > c.active_prev && c.active > 0)
        .sort((a, b) => (b.active - b.active_prev) - (a.active - a.active_prev))
        .slice(0, 3)
    : []
  const newcomers = data
    ? [...data.countries]
        .filter((c) => c.new > 0)
        .sort((a, b) => b.new - a.new || a.code.localeCompare(b.code))
        .slice(0, 3)
    : []

  return (
    <AdminShell>
      <main className="mx-auto w-full max-w-[1200px] min-w-0 px-4 py-8 sm:px-6 sm:py-12 overflow-x-hidden">
        <SectionHeader label="Administration" title="Audience" />
        <p className="-mt-3 mb-6 font-mono text-[11px] text-mathua-muted">
          Who is using Mathua, on what and where activity is growing. Country comes from the
          Cloudflare edge and the device class from the sign-in request, both refreshed at most
          once a day per signed-in learner. A VPN or mobile carrier can place someone elsewhere;
          no IP address, city or user agent is stored.
        </p>

        <div className="mb-6 flex flex-wrap items-center gap-2">
          <span className="font-mono text-[10px] uppercase tracking-[0.15em] text-mathua-muted">Window</span>
          <div role="group" aria-label="Audience window" className="flex border border-mathua-border">
            {RANGES.map((value) => (
              <button
                key={value}
                type="button"
                onClick={() => setDays(value)}
                aria-pressed={days === value}
                className={`h-8 border-l border-mathua-border px-3 font-mono text-[11px] first:border-l-0 ${
                  days === value
                    ? 'bg-mathua-blue text-mathua-bg'
                    : 'text-mathua-secondary hover:text-mathua-primary'
                }`}
              >
                {value} days
              </button>
            ))}
          </div>
          <span className="font-mono text-[10px] text-mathua-muted">
            Compared with the {days} days before
          </span>
        </div>

        {!isStaff && (
          <p role="status" className="mb-6 border border-mathua-border bg-mathua-surface p-4 font-mono text-xs text-mathua-secondary">
            This area needs a staff role. If you hold one, sign in with that account.
          </p>
        )}

        {error && (
          <p role="alert" className="mb-6 border border-mathua-red bg-mathua-surface p-4 font-mono text-xs text-mathua-red">
            {error}
          </p>
        )}

        {!data && !error && isStaff && (
          <div className="space-y-3" aria-busy="true">
            <p role="status" className="sr-only">Loading audience…</p>
            <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
              {[0, 1, 2, 3].map((i) => <Skeleton key={i} className="h-28 w-full" />)}
            </div>
            <div className="grid gap-3 lg:grid-cols-3">
              <Skeleton className="h-80 w-full lg:col-span-2" />
              <Skeleton className="h-80 w-full" />
            </div>
          </div>
        )}

        {data && (
          <>
            <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
              <StatCard
                label="Known accounts"
                value={`${fmt(data.located)} of ${fmt(data.total)}`}
                note={data.total > 0 ? `${Math.round((data.located / data.total) * 100)}% have a recorded country` : 'No accounts yet'}
              />
              <StatCard
                label={`Active, last ${windowDays} days`}
                value={fmt(data.active)}
                note={growthNote(activeDelta ?? { label: 'no activity yet', tone: 'flat' }, windowDays)}
                tone={activeDelta?.tone}
              />
              <StatCard
                label={`Visits, last ${windowDays} days`}
                value={fmt(data.visits)}
                note={growthNote(visitsDelta ?? { label: 'no activity yet', tone: 'flat' }, windowDays)}
                tone={visitsDelta?.tone}
              />
              <StatCard
                label={`New accounts, last ${windowDays} days`}
                value={fmt(data.new)}
                note={growthNote(newDelta ?? { label: 'no activity yet', tone: 'flat' }, windowDays)}
                tone={newDelta?.tone}
              />
            </div>

            {data.located === 0 ? (
              <div className="mt-3 border border-mathua-border bg-mathua-surface p-6">
                <p className="font-mono text-xs text-mathua-secondary">No countries recorded yet.</p>
                <p className="mt-2 max-w-[60ch] font-mono text-[11px] text-mathua-muted">
                  The first country is stored the next time a signed-in learner makes a request
                  through the site. History builds from that day, so the globe fills in as people
                  return.
                </p>
              </div>
            ) : (
              <>
                <div className="mt-3 grid gap-3 lg:grid-cols-3">
                  <div className="border border-mathua-border bg-mathua-surface p-4 lg:col-span-2">
                    <div className="mb-1 flex flex-wrap items-baseline justify-between gap-2">
                      <h3 className="font-mono text-[10px] uppercase tracking-[0.15em] text-mathua-muted">Where learners are</h3>
                      <span
                        className="hidden font-mono text-[10px] text-mathua-muted sm:inline"
                        title="Arcs run from the busiest country to the next six by activity; they rank activity, not travel."
                      >
                        Marker size: active in the window · arcs: busiest six
                      </span>
                    </div>
                    <CountryGlobe countries={data.countries} />
                  </div>

                  <div className="flex flex-col border border-mathua-border bg-mathua-surface p-4">
                    <h3 className="font-mono text-[10px] uppercase tracking-[0.15em] text-mathua-muted">Fastest growing</h3>
                    {movers.length === 0 ? (
                      <p className="mt-3 font-mono text-[11px] text-mathua-muted">
                        No country has grown in this window yet. Growth compares the last
                        {windowDays} days with the {windowDays} before them.
                      </p>
                    ) : (
                      <ol className="mt-3 space-y-3">
                        {movers.map((c) => {
                          const d = delta(c.active, c.active_prev)
                          return (
                            <li key={c.code} className="flex items-baseline justify-between gap-3">
                              <span className="min-w-0">
                                <span className="block truncate font-mono text-xs text-mathua-primary">{countryName(c.code)}</span>
                                <span className="block font-mono text-[10px] text-mathua-muted">
                                  {c.active_prev} to {c.active} active
                                </span>
                              </span>
                              <span className={`shrink-0 font-mono text-sm ${deltaTone(d.tone)}`}>{d.label}</span>
                            </li>
                          )
                        })}
                      </ol>
                    )}
                    <p className="mt-4 border-t border-mathua-border pt-3 font-mono text-[10px] text-mathua-muted">
                      Growth is activity, not signups.
                    </p>
                    {newcomers.length > 0 && (
                      <div className="mt-4 border-t border-mathua-border pt-3">
                        <p className="font-mono text-[10px] uppercase tracking-[0.12em] text-mathua-muted">
                          New accounts, last {windowDays} days
                        </p>
                        <ol className="mt-2 space-y-2">
                          {newcomers.map((c) => (
                            <li key={c.code} className="flex items-baseline justify-between gap-3">
                              <span className="truncate font-mono text-[11px] text-mathua-secondary">{countryName(c.code)}</span>
                              <span className="shrink-0 font-mono text-[11px] text-mathua-primary">+{fmt(c.new)}</span>
                            </li>
                          ))}
                        </ol>
                      </div>
                    )}
                    <div className="mt-auto border-t border-mathua-border pt-3">
                      <div className="flex items-baseline justify-between gap-2 font-mono text-[10px] uppercase tracking-[0.12em] text-mathua-muted">
                        <span>Coverage</span>
                        <span>{Math.round((data.located / data.total) * 100)}%</span>
                      </div>
                      <div className="mt-2 h-1 w-full bg-mathua-border">
                        <div
                          className="h-full bg-mathua-blue"
                          style={{ width: `${Math.max(3, Math.round((data.located / data.total) * 100))}%` }}
                        />
                      </div>
                      <p className="mt-2 font-mono text-[10px] text-mathua-muted">
                        {fmt(data.total - data.located)} accounts have no country yet. The next
                        signed-in visit records one.
                      </p>
                    </div>
                  </div>
                </div>

                <div className="mt-3 grid gap-3 lg:grid-cols-5">
                  <div className="flex flex-col border border-mathua-border bg-mathua-surface p-4 lg:col-span-2">
                    <div className="mb-3 flex flex-wrap items-baseline justify-between gap-2">
                      <h3 className="font-mono text-[10px] uppercase tracking-[0.15em] text-mathua-muted">Activity, last {windowDays} days</h3>
                      <span className="font-mono text-[10px] text-mathua-muted" title="A visit is one signed-in learner-day.">
                        {fmt(data.visits)} visits
                      </span>
                    </div>
                    <Sparkline daily={data.daily} />
                    <div className="mt-auto grid grid-cols-2 gap-3 border-t border-mathua-border pt-3">
                      <div>
                        <p className="font-mono text-[10px] uppercase tracking-[0.12em] text-mathua-muted">Visits per active learner</p>
                        <p className="mt-1 font-mono text-lg text-mathua-primary">
                          {data.active > 0 ? (data.visits / data.active).toFixed(1) : '—'}
                        </p>
                        <p className="font-mono text-[10px] text-mathua-muted">One learner-day each</p>
                      </div>
                      <div>
                        <p className="font-mono text-[10px] uppercase tracking-[0.12em] text-mathua-muted">Busiest day</p>
                        <p className="mt-1 font-mono text-lg text-mathua-primary">
                          {Math.max(0, ...data.daily.map((d) => d.active))}
                        </p>
                        <p className="font-mono text-[10px] text-mathua-muted">learners in one day</p>
                      </div>
                    </div>
                  </div>

                  <div className="border border-mathua-border bg-mathua-surface p-4 lg:col-span-3">
                    <div className="mb-1 flex flex-wrap items-baseline justify-between gap-2">
                      <h3 className="font-mono text-[10px] uppercase tracking-[0.15em] text-mathua-muted">Devices, last {windowDays} days</h3>
                      <span className="font-mono text-[10px] text-mathua-muted">Drag the gallery to turn it</span>
                    </div>
                    {data.devices.length === 0 ? (
                      <p className="mt-3 font-mono text-[11px] text-mathua-muted">
                        Devices are recorded from the next signed-in visit.
                      </p>
                    ) : (
                      <>
                        <DeviceGallery devices={data.devices} />
                        <div className="mt-4 border-t border-mathua-border pt-3">
                          <div className="mb-2 flex flex-wrap items-baseline justify-between gap-2">
                            <h4 className="font-mono text-[10px] uppercase tracking-[0.12em] text-mathua-muted">Daily trend</h4>
                            <span className="font-mono text-[10px] text-mathua-muted">Distinct learners per day</span>
                          </div>
                          <DeviceLines daily={data.device_daily} />
                        </div>
                      </>
                    )}
                  </div>
                </div>

                <div className="mt-3 border border-mathua-border bg-mathua-surface">
                  <Table className="block md:table" role="table">
                    <TableCaption className="sr-only">
                      Countries with known accounts, active learners in the selected window and the window before it, and new accounts
                    </TableCaption>
                    <TableHeader className="hidden md:table-header-group" role="rowgroup">
                      <TableRow className="hover:bg-transparent" role="row">
                        <TableHead className="px-4 py-2">Country</TableHead>
                        <TableHead className="px-4 py-2 text-right">Known</TableHead>
                        <TableHead className="px-4 py-2 text-right">Active</TableHead>
                        <TableHead className="px-4 py-2 text-right">Previous</TableHead>
                        <TableHead className="px-4 py-2 text-right">Change</TableHead>
                        <TableHead className="px-4 py-2 text-right">New</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody className="block md:table-row-group" role="rowgroup">
                      {data.countries.map((c) => {
                        const d = delta(c.active, c.active_prev)
                        return (
                          <TableRow
                            key={c.code}
                            role="row"
                            className="grid grid-cols-3 gap-x-4 gap-y-2 border-b border-mathua-border px-4 py-3 last:border-0 hover:bg-transparent md:table-row md:px-0 md:py-0"
                          >
                            <TableCell role="cell" className="col-span-3 px-0 py-0 md:table-cell md:px-4 md:py-3">
                              <span className="font-mono text-xs text-mathua-primary">{countryName(c.code)}</span>
                              <span className="ml-2 font-mono text-[10px] text-mathua-muted">{c.code}</span>
                            </TableCell>
                            <TableCell
                              role="cell"
                              className="flex flex-col gap-1 px-0 py-0 md:table-cell md:px-4 md:py-3 md:text-right"
                            >
                              <span className="font-mono text-[10px] uppercase tracking-[0.1em] text-mathua-muted md:hidden">Known</span>
                              <span className="font-mono text-xs text-mathua-secondary">{fmt(c.known)}</span>
                            </TableCell>
                            <TableCell role="cell" className="flex flex-col gap-1 px-0 py-0 md:table-cell md:px-4 md:py-3">
                              <span className="font-mono text-[10px] uppercase tracking-[0.1em] text-mathua-muted md:hidden">Active</span>
                              <div className="flex items-baseline justify-between gap-2 md:items-center md:justify-end">
                                <span className="flex items-center gap-2">
                                  <span className="font-mono text-xs text-mathua-primary">{fmt(c.active)}</span>
                                  <span aria-hidden="true" className="hidden h-1 w-16 bg-mathua-border sm:block">
                                    <span className="block h-full bg-mathua-blue" style={{ width: `${activityBar(c, maxActive)}%` }} />
                                  </span>
                                </span>
                              </div>
                            </TableCell>
                            <TableCell
                              role="cell"
                              className="hidden px-0 py-0 md:table-cell md:px-4 md:py-3 md:text-right"
                            >
                              <span className="font-mono text-xs text-mathua-secondary">{fmt(c.active_prev)}</span>
                            </TableCell>
                            <TableCell
                              role="cell"
                              className="hidden px-0 py-0 md:table-cell md:px-4 md:py-3 md:text-right"
                            >
                              <span className="font-mono text-[10px] uppercase tracking-[0.1em] text-mathua-muted md:hidden">Change</span>
                              <span className={`font-mono text-xs ${deltaTone(d.tone)}`}>{d.label}</span>
                            </TableCell>
                            <TableCell
                              role="cell"
                              className="flex flex-col gap-1 px-0 py-0 md:table-cell md:px-4 md:py-3 md:text-right"
                            >
                              <span className="font-mono text-[10px] uppercase tracking-[0.1em] text-mathua-muted md:hidden">New</span>
                              <span className="font-mono text-xs text-mathua-secondary">
                                {fmt(c.new)}
                                {c.new_prev > 0 && (
                                  <span
                                    className="ml-1 text-mathua-muted"
                                    title={`${c.new_prev} in the previous ${windowDays} days`}
                                  >
                                    ({c.new - c.new_prev >= 0 ? '+' : ''}{c.new - c.new_prev})
                                  </span>
                                )}
                              </span>
                            </TableCell>
                          </TableRow>
                        )
                      })}
                    </TableBody>
                  </Table>
                </div>
              </>
            )}

            <p className="mt-4 font-mono text-[10px] text-mathua-muted">
              Generated <time dateTime={data.generated_at}>{new Date(data.generated_at).toLocaleString()}</time>.
              Country data is refreshed at most once a day per signed-in learner.
            </p>
          </>
        )}

        <p className="mt-8 font-mono text-[11px] text-mathua-muted">
          Account-level detail lives in{' '}
          <Link href="/admin/users" className="text-mathua-secondary hover:text-mathua-blue hover:underline">Users</Link>.
          This page never shows a city, a coordinate or an IP address.
        </p>
      </main>
    </AdminShell>
  )
}
