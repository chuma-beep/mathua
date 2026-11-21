'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import AdminShell from '@/components/admin/AdminShell'
import CountryGlobe from '@/components/admin/CountryGlobe'
import SectionHeader from '@/components/SectionHeader'
import { Skeleton } from '@/components/ui/skeleton'
import { Table, TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { useAdminMe } from '@/hooks/useAdminMe'
import {
  fetchAdminLocations,
  getErrorMessage,
  type AdminLocationCountry,
  type AdminLocationDay,
  type AdminLocations,
} from '../../../lib/api'

// Where learners are, at country granularity.
//
// The country is resolved by Cloudflare at the edge and stored once a day for a
// signed-in learner. There is no city, no coordinate and no IP address in the
// answer, so this page cannot become a tracker even for an administrator. The
// honest caveat is in the copy: a VPN or a carrier can put a learner somewhere
// they are not, and a country is a coarse answer to "where" by construction.
//
// Growth is shown as two separate facts, because they are different questions:
// activity growth (distinct learners practising this week against last week) and
// new accounts (signups this week against last week). A country can grow in
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
    return current === 0 ? { label: 'no activity yet', tone: 'flat' } : { label: 'first recorded week', tone: 'up' }
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

function growthNote(d: Delta): string {
  if (d.label === 'first recorded week') return 'First recorded week, no prior week to compare'
  if (d.label === 'no activity yet') return 'No activity recorded yet'
  return `${d.label} against the previous 7 days`
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
    <div className="text-mathua-blue">
      <svg viewBox={`0 0 ${width} ${height}`} className="h-32 w-full" aria-hidden="true" preserveAspectRatio="none">
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
        Distinct active learners per day over the last thirty days, ending at {peak} at the peak.
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
  if (max <= 0 || country.active_7d <= 0) return 0
  return Math.max(3, Math.round((country.active_7d / max) * 100))
}

function fmt(value: number): string {
  return value.toLocaleString()
}

export default function AdminLocationsPage() {
  const { isStaff } = useAdminMe()
  const [data, setData] = useState<AdminLocations | null>(null)
  const [error, setError] = useState('')

  useEffect(() => {
    if (!isStaff) return
    let live = true
    fetchAdminLocations()
      .then((res) => { if (live) { setData(res); setError('') } })
      .catch((e) => { if (live) setError(getErrorMessage(e)) })
    return () => { live = false }
  }, [isStaff])

  const activeDelta = data ? delta(data.active_7d, data.active_prev_7d) : null
  const newDelta = data ? delta(data.new_7d, data.new_prev_7d) : null
  const visitsDelta = data ? delta(data.visits_7d, data.visits_prev_7d) : null
  const maxDevice = data ? Math.max(0, ...data.devices.map((d) => d.active_7d)) : 0
  const maxActive = data ? Math.max(0, ...data.countries.map((c) => c.active_7d)) : 0
  const movers = data
    ? [...data.countries]
        .filter((c) => c.active_7d > c.active_prev_7d && c.active_7d > 0)
        .sort((a, b) => (b.active_7d - b.active_prev_7d) - (a.active_7d - a.active_prev_7d))
        .slice(0, 3)
    : []
  const newcomers = data
    ? [...data.countries]
        .filter((c) => c.new_7d > 0)
        .sort((a, b) => b.new_7d - a.new_7d || a.code.localeCompare(b.code))
        .slice(0, 3)
    : []

  return (
    <AdminShell>
      <main className="mx-auto w-full max-w-[1200px] min-w-0 px-4 py-8 sm:px-6 sm:py-12 overflow-x-hidden">
        <SectionHeader label="Administration" title="Locations" />
        <p className="-mt-3 mb-6 font-mono text-[11px] text-mathua-muted">
          Countries where accounts are and where activity is growing. The country comes from
          Cloudflare&apos;s edge, is refreshed for a signed-in learner at most once a day, and is
          approximate: a VPN or a mobile carrier can place a learner elsewhere. No IP address or city
          is stored.
        </p>

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
            <p role="status" className="sr-only">Loading locations…</p>
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
                label="Active this week"
                value={fmt(data.active_7d)}
                note={growthNote(activeDelta ?? { label: 'no activity yet', tone: 'flat' })}
                tone={activeDelta?.tone}
              />
              <StatCard
                label="Visits this week"
                value={fmt(data.visits_7d)}
                note={growthNote(visitsDelta ?? { label: 'no activity yet', tone: 'flat' })}
                tone={visitsDelta?.tone}
              />
              <StatCard
                label="New accounts this week"
                value={fmt(data.new_7d)}
                note={growthNote(newDelta ?? { label: 'no activity yet', tone: 'flat' })}
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
                      <span className="hidden font-mono text-[10px] text-mathua-muted sm:inline">Marker size: active this week</span>
                    </div>
                    <CountryGlobe countries={data.countries} />
                  </div>

                  <div className="border border-mathua-border bg-mathua-surface p-4">
                    <h3 className="font-mono text-[10px] uppercase tracking-[0.15em] text-mathua-muted">Fastest growing</h3>
                    {movers.length === 0 ? (
                      <p className="mt-3 font-mono text-[11px] text-mathua-muted">
                        No country has grown against last week yet. Growth compares the last seven
                        days with the seven before them.
                      </p>
                    ) : (
                      <ol className="mt-3 space-y-3">
                        {movers.map((c) => {
                          const d = delta(c.active_7d, c.active_prev_7d)
                          return (
                            <li key={c.code} className="flex items-baseline justify-between gap-3">
                              <span className="min-w-0">
                                <span className="block truncate font-mono text-xs text-mathua-primary">{countryName(c.code)}</span>
                                <span className="block font-mono text-[10px] text-mathua-muted">
                                  {c.active_prev_7d} to {c.active_7d} active
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
                          New accounts this week
                        </p>
                        <ol className="mt-2 space-y-2">
                          {newcomers.map((c) => (
                            <li key={c.code} className="flex items-baseline justify-between gap-3">
                              <span className="truncate font-mono text-[11px] text-mathua-secondary">{countryName(c.code)}</span>
                              <span className="shrink-0 font-mono text-[11px] text-mathua-primary">+{fmt(c.new_7d)}</span>
                            </li>
                          ))}
                        </ol>
                      </div>
                    )}
                  </div>
                </div>

                <div className="mt-3 border border-mathua-border bg-mathua-surface p-4">
                  <div className="mb-3 flex flex-wrap items-baseline justify-between gap-2">
                    <h3 className="font-mono text-[10px] uppercase tracking-[0.15em] text-mathua-muted">Activity, last 30 days</h3>
                    <span className="font-mono text-[10px] text-mathua-muted" title="A visit is one signed-in learner-day.">
                      {fmt(data.visits_30d)} visits
                    </span>
                  </div>
                  <Sparkline daily={data.daily} />
                  <div className="mt-4 border-t border-mathua-border pt-3">
                    <div className="flex flex-wrap items-baseline justify-between gap-2">
                      <h4 className="font-mono text-[10px] uppercase tracking-[0.12em] text-mathua-muted">Devices this week</h4>
                      <span className="font-mono text-[10px] text-mathua-muted">One class per learner</span>
                    </div>
                    {data.devices.length === 0 ? (
                      <p className="mt-2 font-mono text-[11px] text-mathua-muted">
                        Devices are recorded from the next signed-in visit.
                      </p>
                    ) : (
                      <ul className="mt-2 grid gap-2 sm:grid-cols-3">
                        {data.devices.map((device) => (
                          <li key={device.device}>
                            <div className="flex items-baseline justify-between gap-2 font-mono text-[11px]">
                              <span className="capitalize text-mathua-secondary">{device.device}</span>
                              <span className="text-mathua-primary">{fmt(device.active_7d)}</span>
                            </div>
                            <div className="mt-1 h-1 w-full bg-mathua-border">
                              <div
                                className="h-full bg-mathua-blue"
                                style={{ width: `${maxDevice > 0 ? Math.max(3, Math.round((device.active_7d / maxDevice) * 100)) : 0}%` }}
                              />
                            </div>
                          </li>
                        ))}
                      </ul>
                    )}
                  </div>
                </div>

                <div className="mt-3 border border-mathua-border bg-mathua-surface">
                  <Table className="block md:table" role="table">
                    <TableCaption className="sr-only">
                      Countries with known accounts, active learners this week and last week, and new accounts
                    </TableCaption>
                    <TableHeader className="hidden md:table-header-group" role="rowgroup">
                      <TableRow className="hover:bg-transparent" role="row">
                        <TableHead className="px-4 py-2">Country</TableHead>
                        <TableHead className="px-4 py-2 text-right">Known</TableHead>
                        <TableHead className="px-4 py-2 text-right">Active 7d</TableHead>
                        <TableHead className="px-4 py-2 text-right">Previous 7d</TableHead>
                        <TableHead className="px-4 py-2 text-right">Change</TableHead>
                        <TableHead className="px-4 py-2 text-right">New 7d</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody className="block md:table-row-group" role="rowgroup">
                      {data.countries.map((c) => {
                        const d = delta(c.active_7d, c.active_prev_7d)
                        return (
                          <TableRow
                            key={c.code}
                            role="row"
                            className="grid grid-cols-2 gap-x-4 gap-y-2 border-b border-mathua-border px-4 py-3 last:border-0 hover:bg-transparent md:table-row md:px-0 md:py-0"
                          >
                            <TableCell role="cell" className="col-span-2 px-0 py-0 md:table-cell md:px-4 md:py-3">
                              <span className="font-mono text-xs text-mathua-primary">{countryName(c.code)}</span>
                              <span className="ml-2 font-mono text-[10px] text-mathua-muted">{c.code}</span>
                            </TableCell>
                            <TableCell
                              role="cell"
                              className="flex items-baseline justify-between gap-2 px-0 py-0 md:table-cell md:px-4 md:py-3 md:text-right"
                            >
                              <span className="font-mono text-[10px] uppercase tracking-[0.1em] text-mathua-muted md:hidden">Known</span>
                              <span className="font-mono text-xs text-mathua-secondary">{fmt(c.known)}</span>
                            </TableCell>
                            <TableCell role="cell" className="px-0 py-0 md:table-cell md:px-4 md:py-3">
                              <div className="flex items-baseline justify-between gap-2 md:items-center md:justify-end">
                                <span className="font-mono text-[10px] uppercase tracking-[0.1em] text-mathua-muted md:hidden">Active 7d</span>
                                <span className="flex items-center gap-2">
                                  <span className="font-mono text-xs text-mathua-primary">{fmt(c.active_7d)}</span>
                                  <span aria-hidden="true" className="h-1 w-12 bg-mathua-border md:w-16">
                                    <span className="block h-full bg-mathua-blue" style={{ width: `${activityBar(c, maxActive)}%` }} />
                                  </span>
                                </span>
                              </div>
                            </TableCell>
                            <TableCell
                              role="cell"
                              className="hidden px-0 py-0 md:table-cell md:px-4 md:py-3 md:text-right"
                            >
                              <span className="font-mono text-xs text-mathua-secondary">{fmt(c.active_prev_7d)}</span>
                            </TableCell>
                            <TableCell
                              role="cell"
                              className="flex items-baseline justify-between gap-2 px-0 py-0 md:table-cell md:px-4 md:py-3 md:text-right"
                            >
                              <span className="font-mono text-[10px] uppercase tracking-[0.1em] text-mathua-muted md:hidden">Change</span>
                              <span className={`font-mono text-xs ${deltaTone(d.tone)}`}>{d.label}</span>
                            </TableCell>
                            <TableCell
                              role="cell"
                              className="flex items-baseline justify-between gap-2 px-0 py-0 md:table-cell md:px-4 md:py-3 md:text-right"
                            >
                              <span className="font-mono text-[10px] uppercase tracking-[0.1em] text-mathua-muted md:hidden">New 7d</span>
                              <span className="font-mono text-xs text-mathua-secondary">
                                {fmt(c.new_7d)}
                                {c.new_prev_7d > 0 && (
                                  <span
                                    className="ml-1 text-mathua-muted"
                                    title={`${c.new_prev_7d} in the previous seven days`}
                                  >
                                    ({c.new_7d - c.new_prev_7d >= 0 ? '+' : ''}{c.new_7d - c.new_prev_7d})
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
