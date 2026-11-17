'use client'

import AdminShell from '@/components/admin/AdminShell'
import { useCallback, useEffect, useMemo, useState } from 'react'
import Link from 'next/link'
import SectionHeader from '../../../components/SectionHeader'
import Avatar from '../../../components/Avatar'
import { ArrowDown, ArrowUp } from '../../../components/icons'
import LeagueAvatar from '../../../components/LeagueAvatar'
import { Card, CardContent, CardHeader, CardTitle } from '../../../components/ui/card'
import { Skeleton } from '../../../components/ui/skeleton'
import { Table, TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRow } from '../../../components/ui/table'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '../../../components/ui/tabs'
import { Bar, BarChart, CartesianGrid, Line, LineChart, XAxis, YAxis } from 'recharts'
import { ChartContainer, ChartLegend, ChartLegendContent, ChartTooltip, ChartTooltipContent, type ChartConfig } from '../../../components/ui/chart'
import { dicebearUrl, isDicebearStyle } from '../../../lib/dicebear'
import { API_BASE, getAdminLeagueRange, getAdminLeagues, getErrorMessage, type AdminLeague, type AdminLeagueMember, type AdminLeagueRange } from '../../../lib/api'
import { DateRangePicker } from '../../../components/ui/date-range-picker'
import type { DateRange } from 'react-day-picker'

// The league console is a reading surface, not a queue. A week has five tiers,
// each with zero or more leagues, and a league has up to thirty seats. The page
// therefore separates the week summary, the tier pyramid and one roster, and
// reads only the selected roster. Nothing here changes a league: XP decides the
// seats and the server applies movement when the week closes.

type Scope = 'week' | 'month' | 'year' | 'custom'

const LEAGUE_SIZE = 30
const TIER_ORDER = ['bronze', 'silver', 'gold', 'platinum', 'diamond']
const TIER_COLOUR: Record<string, string> = {
  bronze: 'var(--league-bronze)',
  silver: 'var(--league-silver)',
  gold: 'var(--league-gold)',
  platinum: 'var(--league-platinum)',
  diamond: 'var(--league-diamond)',
}

function weekKeyOf(d: Date): string {
  const day = new Date(Date.UTC(d.getUTCFullYear(), d.getUTCMonth(), d.getUTCDate()))
  const sinceMonday = (day.getUTCDay() + 6) % 7
  day.setUTCDate(day.getUTCDate() - sinceMonday)
  return day.toISOString().slice(0, 10)
}

function shiftWeek(week: string, days: number): string {
  const t = new Date(`${week}T00:00:00Z`)
  t.setUTCDate(t.getUTCDate() + days)
  return t.toISOString().slice(0, 10)
}

function weekLabel(week: string): string {
  const t = new Date(`${week}T00:00:00Z`)
  return t.toLocaleDateString('en-GB', { day: 'numeric', month: 'short', year: 'numeric', timeZone: 'UTC' })
}

function tierLabel(tier: string): string {
  return tier.charAt(0).toUpperCase() + tier.slice(1)
}

function memberLabel(m: AdminLeagueMember): string {
  return m.name || m.username || m.student_id
}

function recentWeeks(week: string, count: number): string[] {
  return Array.from({ length: count }, (_, i) => shiftWeek(week, -7 * i))
}

type Selected = { tier: string; league: number }

function leagueKey(l: Pick<AdminLeague, 'tier' | 'league_no'> | Selected): string {
  return 'league_no' in l ? `${l.tier}/${l.league_no}` : `${l.tier}/${l.league}`
}

function MemberAvatar({ member, tier }: { member: AdminLeagueMember; tier: string }) {
  const pick = member.avatar_dicebear
  const storedURL = member.avatar_url
  const url = storedURL
    ? (storedURL.startsWith('/') ? `${API_BASE}${storedURL}` : storedURL)
    : (pick && isDicebearStyle(pick.style) ? dicebearUrl(pick.style, pick.seed) : undefined)
  return (
    <LeagueAvatar
      tier={tier}
      weeklyXp={member.xp}
      seed={member.student_id}
      name={memberLabel(member)}
      size={24}
      url={url}
      preset={member.avatar_preset ?? null}
    />
  )
}

function RosterTable({ league, members }: { league: AdminLeague; members: AdminLeagueMember[] }) {
  if (members.length === 0) {
    return <p className="py-8 text-center font-mono text-[11px] text-mathua-muted">Nobody holds a seat in this league.</p>
  }
  const promote = league.movers
  const relegate = league.movers
  return (
    <Table>
      <TableCaption>
        Top {promote} promote and bottom {relegate} relegate when the week closes. Rank and XP are the current week only.
      </TableCaption>
      <TableHeader>
        <TableRow>
          <TableHead className="w-12">Rank</TableHead>
          <TableHead>Learner</TableHead>
          <TableHead className="w-24 text-right">XP</TableHead>
          <TableHead className="w-24 text-right">Movement</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {members.map((m) => {
          const up = promote > 0 && m.rank <= promote
          const down = relegate > 0 && m.rank > league.size - relegate
          return (
            <TableRow key={m.student_id}>
              <TableCell className="text-mathua-secondary">{m.rank}</TableCell>
              <TableCell>
                <span className="flex min-w-0 items-center gap-2">
                  <MemberAvatar member={m} tier={league.tier} />
                  <span className="min-w-0">
                    <Link
                      href={`/admin/users/detail?id=${encodeURIComponent(m.student_id)}`}
                      className="block truncate text-mathua-blue hover:underline"
                      title={memberLabel(m)}
                    >
                      {memberLabel(m)}
                    </Link>
                    {m.username && m.username !== memberLabel(m) && (
                      <span className="block truncate text-[10px] text-mathua-muted">@{m.username}</span>
                    )}
                  </span>
                </span>
              </TableCell>
              <TableCell className="text-right tabular-nums text-mathua-primary">{m.xp.toLocaleString()}</TableCell>
              <TableCell className="text-right">
                {up ? (
                  <span className="inline-flex items-center gap-1 text-mathua-green">
                    <ArrowUp size={13} aria-hidden="true" />
                    Up
                  </span>
                ) : down ? (
                  <span className="inline-flex items-center gap-1 text-mathua-red">
                    <ArrowDown size={13} aria-hidden="true" />
                    Down
                  </span>
                ) : (
                  <span className="inline-flex items-center gap-1 text-mathua-muted">
                    <span aria-hidden="true">·</span>
                    Stays
                  </span>
                )}
              </TableCell>
            </TableRow>
          )
        })}
      </TableBody>
    </Table>
  )
}

function RosterPanel({
  league,
  members,
  loading,
  error,
  onRetry,
}: {
  league: AdminLeague | null
  members: AdminLeagueMember[]
  loading: boolean
  error: string
  onRetry: () => void
}) {
  if (!league) {
    return <p className="py-10 text-center font-mono text-[11px] text-mathua-muted">Choose a league to read its roster.</p>
  }
  return (
    <div>
      <div className="mb-3 flex flex-wrap items-baseline justify-between gap-2">
        <div>
          <p className="font-mono text-[10px] uppercase tracking-wider text-mathua-muted">{tierLabel(league.tier)}</p>
          <h3 className="font-serif text-lg text-mathua-primary">League {league.league_no}</h3>
        </div>
        <p className="font-mono text-[10px] text-mathua-muted">
          {league.size} of {LEAGUE_SIZE} seats
        </p>
      </div>
      {error ? (
        <div className="border border-mathua-red bg-mathua-surface p-3">
          <p role="alert" className="font-mono text-[11px] text-mathua-red">{error}</p>
          <button type="button" onClick={onRetry} className="mt-2 h-8 border border-mathua-border px-2.5 font-mono text-[11px] text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue">
            Retry
          </button>
        </div>
      ) : loading ? (
        <div className="space-y-2">
          {Array.from({ length: 5 }, (_, i) => <Skeleton key={i} className="h-9 w-full" />)}
        </div>
      ) : (
        <RosterTable league={league} members={members} />
      )}
    </div>
  )
}

type PyramidRow = {
  tier: string
  leagues: number
  seats: number
  capacity: number
  unique: number
  onSelect?: () => void
}

function LeaguePyramid({ rows }: { rows: PyramidRow[] }) {
  const totalSeats = rows.reduce((n, r) => n + r.seats, 0)
  return (
    <>
      <div className="mt-3 grid grid-cols-2 gap-2 sm:grid-cols-3 lg:grid-cols-5">
        {rows.map(row => {
          const fill = row.capacity > 0 ? Math.round((row.seats / row.capacity) * 100) : 0
          const content = (
            <>
              <div className="flex items-center justify-between gap-2">
                <span className="font-mono text-[10px] uppercase tracking-wider text-mathua-secondary">{tierLabel(row.tier)}</span>
                <span className="h-2 w-2" style={{ backgroundColor: TIER_COLOUR[row.tier] }} aria-hidden="true" />
              </div>
              <p className="mt-2 font-mono text-xl tabular-nums text-mathua-primary">{row.leagues}</p>
              <p className="font-mono text-[10px] text-mathua-muted">league{row.leagues === 1 ? '' : 's'}</p>
              <p className="mt-2 font-mono text-[10px] text-mathua-secondary">{row.seats.toLocaleString()} / {row.capacity.toLocaleString()} seats</p>
              <p className="font-mono text-[10px] text-mathua-muted">{fill}% full{row.unique > 0 ? ` · ${row.unique.toLocaleString()} learners` : ''}</p>
            </>
          )
          return row.onSelect ? (
            <button key={row.tier} type="button" onClick={row.onSelect} className="border border-mathua-border bg-mathua-surface p-3 text-left hover:border-mathua-blue">
              {content}
            </button>
          ) : (
            <div key={row.tier} className={`border border-mathua-border bg-mathua-surface p-3 ${row.leagues === 0 ? 'opacity-55' : ''}`}>
              {content}
            </div>
          )
        })}
      </div>
      <div className="mt-3 flex h-2 w-full overflow-hidden border border-mathua-border" aria-hidden="true">
        {rows.map(row => (
          <span key={row.tier} style={{ width: `${totalSeats > 0 ? (row.seats / totalSeats) * 100 : 0}%`, backgroundColor: TIER_COLOUR[row.tier] }} />
        ))}
      </div>
    </>
  )
}

function TierMix({ rows }: { rows: PyramidRow[] }) {
  const total = rows.reduce((n, r) => n + r.seats, 0)
  const max = Math.max(1, ...rows.map(r => r.seats))
  return (
    <div className="mt-3 space-y-3">
      {rows.map(row => {
        const share = total > 0 ? (row.seats / total) * 100 : 0
        return (
          <div key={row.tier} className="grid grid-cols-[5.5rem_1fr_8rem] items-center gap-3">
            <span className="font-mono text-[10px] uppercase tracking-wider text-mathua-secondary">{tierLabel(row.tier)}</span>
            <span className="h-3 w-full bg-mathua-surface-elevated">
              <span className="block h-3" style={{ width: `${(row.seats / max) * 100}%`, backgroundColor: TIER_COLOUR[row.tier] }} />
            </span>
            <span className="text-right font-mono text-[10px] tabular-nums text-mathua-muted">
              {row.seats.toLocaleString()} seats · {share.toFixed(1)}%
            </span>
          </div>
        )
      })}
    </div>
  )
}

function TierViews({ rows, label }: { rows: PyramidRow[]; label: string }) {
  return (
    <Tabs defaultValue="pyramid">
      <div className="flex flex-wrap items-end justify-between gap-2">
        <TabsList className="h-10 justify-start gap-4 rounded-none border-b border-mathua-border bg-transparent p-0 text-mathua-secondary">
          <TabsTrigger
            value="pyramid"
            className="h-10 rounded-none border-b-2 border-transparent bg-transparent px-3 py-0 font-mono text-[11px] text-mathua-secondary shadow-none data-[state=active]:border-mathua-blue data-[state=active]:bg-transparent data-[state=active]:text-mathua-blue data-[state=active]:shadow-none focus-visible:ring-mathua-blue"
          >
            Tier pyramid
          </TabsTrigger>
          <TabsTrigger
            value="mix"
            className="h-10 rounded-none border-b-2 border-transparent bg-transparent px-3 py-0 font-mono text-[11px] text-mathua-secondary shadow-none data-[state=active]:border-mathua-blue data-[state=active]:bg-transparent data-[state=active]:text-mathua-blue data-[state=active]:shadow-none focus-visible:ring-mathua-blue"
          >
            Tier mix
          </TabsTrigger>
        </TabsList>
        <span className="pb-2 font-mono text-[10px] text-mathua-muted">{label}</span>
      </div>
      <TabsContent value="pyramid"><LeaguePyramid rows={rows} /></TabsContent>
      <TabsContent value="mix"><TierMix rows={rows} /></TabsContent>
    </Tabs>
  )
}

const NETFLOW_CONFIG = {
  net: { label: 'Net flow', color: 'var(--chart-1)' },
} satisfies ChartConfig

const RATE_CONFIG = {
  promoted: { label: 'Promoted', color: 'var(--chart-3)' },
  demoted: { label: 'Demoted', color: 'var(--chart-4)' },
} satisfies ChartConfig

function MovementSection({ data }: { data: AdminLeagueRange }) {
  const comparable = Math.max(1, data.seats - data.new)
  const rows = [
    { label: 'Promoted', count: data.promoted, unique: data.unique_promoted, rate: (data.promoted / comparable) * 100, tone: 'text-mathua-green', bg: 'bg-mathua-green' },
    { label: 'Demoted', count: data.demoted, unique: data.unique_demoted, rate: (data.demoted / comparable) * 100, tone: 'text-mathua-red', bg: 'bg-mathua-red' },
    { label: 'Stayed', count: data.stayed, unique: 0, rate: (data.stayed / comparable) * 100, tone: 'text-mathua-blue', bg: 'bg-mathua-blue' },
    { label: 'New seats', count: data.new, unique: 0, rate: (data.new / Math.max(1, data.seats)) * 100, tone: 'text-mathua-secondary', bg: 'bg-mathua-muted' },
    { label: 'Left', count: data.left, unique: 0, rate: (data.left / Math.max(1, data.seats)) * 100, tone: 'text-mathua-secondary', bg: 'bg-mathua-border-strong' },
  ]
  const max = Math.max(1, ...rows.map(r => r.count))
  const flow = data.buckets.map(b => ({ label: b.label, net: b.promoted - b.demoted }))
  const rates = data.buckets.map(b => {
    const denominator = Math.max(1, b.seats - b.new)
    return { label: b.label, promoted: Number(((b.promoted / denominator) * 100).toFixed(1)), demoted: Number(((b.demoted / denominator) * 100).toFixed(1)) }
  })
  const matrix = new Map(data.matrix.map(cell => [`${cell.from}/${cell.to}`, cell.count]))
  return (
    <section className="mt-8" aria-labelledby="movement-heading">
      <div className="flex flex-wrap items-baseline justify-between gap-2 border-b border-mathua-border pb-2">
        <div>
          <h2 id="movement-heading" className="font-mono text-xs uppercase tracking-wider text-mathua-primary">Movement</h2>
          <p className="mt-1 font-mono text-[10px] text-mathua-muted">How seats changed between adjacent weeks.</p>
        </div>
        <span className="font-mono text-[10px] text-mathua-muted">{data.weeks} weeks · {data.seats} seats</span>
      </div>

      <div className="mt-4 grid gap-3 lg:grid-cols-[2fr_1fr]">
        <div className="border border-mathua-border bg-mathua-surface p-4">
          <div className="space-y-3">
            {rows.map(row => (
              <div key={row.label} className="grid grid-cols-[6rem_1fr_7rem] items-center gap-3">
                <span className="font-mono text-[10px] uppercase tracking-wider text-mathua-secondary">{row.label}</span>
                <span className="h-2 w-full bg-mathua-surface-elevated">
                  <span className={`block h-2 ${row.bg}`} style={{ width: `${(row.count / max) * 100}%` }} />
                </span>
                <span className="text-right font-mono text-[11px] tabular-nums text-mathua-primary">
                  {row.count.toLocaleString()}
                  <span className="ml-1 text-[10px] text-mathua-muted">{row.rate.toFixed(1)}%</span>
                </span>
              </div>
            ))}
          </div>
          <p className="mt-3 font-mono text-[10px] text-mathua-muted">Unique movers: {data.unique_promoted} promoted · {data.unique_demoted} demoted. Events can exceed learners when someone moves more than once.</p>
        </div>

        <div className="border border-mathua-border bg-mathua-surface p-4">
          <p className="font-mono text-[10px] uppercase tracking-wider text-mathua-muted">Net flow</p>
          <p className={`mt-2 font-mono text-3xl tabular-nums ${data.promoted - data.demoted >= 0 ? 'text-mathua-green' : 'text-mathua-red'}`}>
            {data.promoted - data.demoted >= 0 ? '+' : ''}{(data.promoted - data.demoted).toLocaleString()}
          </p>
          <p className="mt-1 font-mono text-[10px] text-mathua-muted">Promotions minus demotions</p>
          <p className="mt-4 font-mono text-[11px] text-mathua-secondary">Retention after next week</p>
          <p className="font-mono text-lg tabular-nums text-mathua-primary">{data.seats > 0 ? (((data.seats - data.left) / data.seats) * 100).toFixed(1) : '—'}%</p>
        </div>
      </div>

      <div className="mt-4 grid gap-3 lg:grid-cols-2">
        <div className="border border-mathua-border bg-mathua-surface p-4">
          <p className="font-mono text-[10px] uppercase tracking-wider text-mathua-secondary">Net flow by week</p>
          <ChartContainer config={NETFLOW_CONFIG} className="mt-3 h-56 w-full">
            <BarChart data={flow} margin={{ top: 4, right: 4, left: -18, bottom: 0 }}>
              <CartesianGrid vertical={false} strokeDasharray="3 3" />
              <XAxis dataKey="label" tickLine={false} axisLine={false} tickMargin={8} minTickGap={18} />
              <YAxis tickLine={false} axisLine={false} />
              <ChartTooltip
                cursor={{ fill: 'var(--surface-elevated)' }}
                content={<ChartTooltipContent indicator="line" />}
              />
              <Bar dataKey="net" fill="var(--color-net)" />
            </BarChart>
          </ChartContainer>
        </div>
        <div className="border border-mathua-border bg-mathua-surface p-4">
          <p className="font-mono text-[10px] uppercase tracking-wider text-mathua-secondary">Movement rates</p>
          <ChartContainer config={RATE_CONFIG} className="mt-3 h-56 w-full">
            <LineChart data={rates} margin={{ top: 4, right: 4, left: -18, bottom: 0 }}>
              <CartesianGrid vertical={false} strokeDasharray="3 3" />
              <XAxis dataKey="label" tickLine={false} axisLine={false} tickMargin={8} minTickGap={18} />
              <YAxis tickLine={false} axisLine={false} tickFormatter={v => `${v}%`} />
              <ChartTooltip content={<ChartTooltipContent indicator="line" />} />
              <ChartLegend content={<ChartLegendContent />} />
              <Line type="monotone" dataKey="promoted" stroke="var(--color-promoted)" strokeWidth={2} dot={false} />
              <Line type="monotone" dataKey="demoted" stroke="var(--color-demoted)" strokeWidth={2} dot={false} />
            </LineChart>
          </ChartContainer>
        </div>
      </div>

      <div className="mt-4 overflow-x-auto border border-mathua-border bg-mathua-surface">
        <table className="w-full min-w-[40rem] text-left">
          <caption className="px-3 pt-3 text-left font-mono text-[10px] text-mathua-muted">Previous tier by current tier. First appearances and returns are excluded.</caption>
          <thead className="border-b border-mathua-border">
            <tr>
              <th scope="col" className="px-3 py-2 font-mono text-[10px] uppercase tracking-wider text-mathua-muted">From</th>
              {TIER_ORDER.map(tier => <th key={tier} scope="col" className="px-3 py-2 text-right font-mono text-[10px] uppercase tracking-wider text-mathua-muted">{tierLabel(tier)}</th>)}
            </tr>
          </thead>
          <tbody>
            {TIER_ORDER.map(from => (
              <tr key={from} className="border-b border-mathua-border last:border-0">
                <th scope="row" className="px-3 py-2 font-mono text-[10px] uppercase tracking-wider text-mathua-secondary">{tierLabel(from)}</th>
                {TIER_ORDER.map(to => {
                  const count = matrix.get(`${from}/${to}`) ?? 0
                  const intensity = Math.min(0.45, count / Math.max(1, data.seats))
                  return (
                    <td key={to} className="px-3 py-2 text-right font-mono text-[11px] tabular-nums text-mathua-primary" style={{ backgroundColor: count > 0 ? `color-mix(in oklab, var(--chart-1) ${Math.round(intensity * 100)}%, transparent)` : undefined }}>
                      {count || '·'}
                    </td>
                  )
                })}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </section>
  )
}

function RangePanel({
  data,
  loading,
  error,
  onRetry,
  rankBy,
  onRankBy,
  onOpenWeek,
}: {
  data: AdminLeagueRange | null
  loading: boolean
  error: string
  onRetry: () => void
  rankBy: 'range' | 'latest'
  onRankBy: (value: 'range' | 'latest') => void
  onOpenWeek: (week: string, tier: string, league: number) => void
}) {
  const [jumpWeek, setJumpWeek] = useState('')
  const [jumpLeagues, setJumpLeagues] = useState<AdminLeague[]>([])
  useEffect(() => {
    if (!data || data.weeks_available.length === 0) { setJumpWeek(''); setJumpLeagues([]); return }
    if (!jumpWeek || !data.weeks_available.includes(jumpWeek)) setJumpWeek(data.weeks_available[data.weeks_available.length - 1])
  }, [data, jumpWeek])
  useEffect(() => {
    if (!jumpWeek) { setJumpLeagues([]); return }
    let live = true
    getAdminLeagues({ week: jumpWeek }).then(res => { if (live) setJumpLeagues(res.leagues) }).catch(() => { if (live) setJumpLeagues([]) })
    return () => { live = false }
  }, [jumpWeek])
  if (loading && !data) {
    return <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">{Array.from({ length: 4 }, (_, i) => <Skeleton key={i} className="h-24 w-full" />)}</div>
  }
  if (error) {
    return (
      <div className="border border-mathua-red bg-mathua-surface p-3">
        <p role="alert" className="font-mono text-[11px] text-mathua-red">{error}</p>
        <button type="button" onClick={onRetry} className="mt-2 h-8 border border-mathua-border px-2.5 font-mono text-[11px] text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue">Retry</button>
      </div>
    )
  }
  if (!data) return null
  const learners = rankBy === 'range' ? data.top_range : data.top_latest
  const maxTier = Math.max(1, ...data.tiers.map(t => t.seats))
  return (
    <>
      <div className="grid gap-2 sm:grid-cols-2 lg:grid-cols-4">
        <Card><CardHeader><CardTitle>League weeks</CardTitle></CardHeader><CardContent><p className="font-mono text-2xl tabular-nums text-mathua-primary">{data.weeks}</p><p className="mt-1 font-mono text-[10px] text-mathua-muted">{data.seats.toLocaleString()} seats</p></CardContent></Card>
        <Card><CardHeader><CardTitle>Unique learners</CardTitle></CardHeader><CardContent><p className="font-mono text-2xl tabular-nums text-mathua-primary">{data.unique_learners.toLocaleString()}</p><p className="mt-1 font-mono text-[10px] text-mathua-muted">Across the range</p></CardContent></Card>
        <Card><CardHeader><CardTitle>Promoted</CardTitle></CardHeader><CardContent><p className="font-mono text-2xl tabular-nums text-mathua-green">{data.promoted.toLocaleString()}</p><p className="mt-1 font-mono text-[10px] text-mathua-muted">{data.demoted.toLocaleString()} demoted · {data.stayed.toLocaleString()} stayed</p></CardContent></Card>
        <Card><CardHeader><CardTitle>New seats</CardTitle></CardHeader><CardContent><p className="font-mono text-2xl tabular-nums text-mathua-primary">{data.new.toLocaleString()}</p><p className="mt-1 font-mono text-[10px] text-mathua-muted">First appearance in a league</p></CardContent></Card>
      </div>

      <section className="mt-8" aria-labelledby="range-tiers">
        <div id="range-tiers">
          <TierViews
            label="Seats, leagues and unique learners across the range"
            rows={TIER_ORDER.map(tier => {
              const row = data.tiers.find(t => t.tier === tier)
              const leagues = row?.leagues ?? 0
              const seats = row?.seats ?? 0
              return { tier, leagues, seats, capacity: leagues * LEAGUE_SIZE, unique: row?.unique_learners ?? 0 }
            })}
          />
        </div>
      </section>

      <MovementSection data={data} />

      <section className="mt-8" aria-labelledby="range-weeks">
        <div className="flex items-baseline justify-between border-b border-mathua-border pb-2">
          <h2 id="range-weeks" className="font-mono text-xs uppercase tracking-wider text-mathua-primary">Weekly breakdown</h2>
          <span className="font-mono text-[10px] text-mathua-muted">Select a row to open that week</span>
        </div>
        <div className="mt-3 overflow-x-auto">
          <Table>
            <TableHeader><TableRow><TableHead>Period</TableHead><TableHead className="text-right">Weeks</TableHead><TableHead className="text-right">Seats</TableHead><TableHead className="text-right">Unique</TableHead><TableHead className="text-right">Up</TableHead><TableHead className="text-right">Down</TableHead></TableRow></TableHeader>
            <TableBody>
              {data.buckets.map(b => (
                <TableRow key={b.key}>
                  <TableCell>
                    {/^\d{4}-\d{2}-\d{2}$/.test(b.key) ? (
                      <Link href={`/admin/leagues?week=${b.key}`} className="text-mathua-blue hover:underline">{weekLabel(b.key)}</Link>
                    ) : b.key}
                  </TableCell>
                  <TableCell className="text-right tabular-nums">{b.weeks}</TableCell>
                  <TableCell className="text-right tabular-nums">{b.seats}</TableCell>
                  <TableCell className="text-right tabular-nums">{b.unique_learners}</TableCell>
                  <TableCell className="text-right tabular-nums text-mathua-green">{b.promoted}</TableCell>
                  <TableCell className="text-right tabular-nums text-mathua-red">{b.demoted}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      </section>

      <section className="mt-8" aria-labelledby="range-top">
        <div className="flex flex-wrap items-center justify-between gap-2 border-b border-mathua-border pb-2">
          <h2 id="range-top" className="font-mono text-xs uppercase tracking-wider text-mathua-primary">Top learners</h2>
          <div className="flex gap-1" role="group" aria-label="Ranking basis">
            <button type="button" aria-pressed={rankBy === 'range'} onClick={() => onRankBy('range')} className={`h-8 border px-2.5 font-mono text-[11px] ${rankBy === 'range' ? 'border-mathua-blue text-mathua-blue' : 'border-mathua-border text-mathua-secondary'}`}>XP in range</button>
            <button type="button" aria-pressed={rankBy === 'latest'} onClick={() => onRankBy('latest')} className={`h-8 border px-2.5 font-mono text-[11px] ${rankBy === 'latest' ? 'border-mathua-blue text-mathua-blue' : 'border-mathua-border text-mathua-secondary'}`}>Latest week</button>
          </div>
        </div>
        {learners.length === 0 ? <p className="py-8 text-center font-mono text-[11px] text-mathua-muted">No ranked learners in this range.</p> : (
          <ol className="mt-3 divide-y divide-mathua-border border border-mathua-border bg-mathua-surface">
            {learners.map((l, i) => (
              <li key={l.student_id} className="flex items-center gap-3 px-3 py-2">
                <span className="w-6 font-mono text-[11px] text-mathua-muted">{i + 1}</span>
                <span className="min-w-0 flex-1 truncate font-mono text-[11px] text-mathua-primary">{l.name || l.username || l.student_id}</span>
                <span className="font-mono text-[10px] uppercase tracking-wider text-mathua-muted">{l.tier || '—'}</span>
                <span className="w-20 text-right font-mono text-[11px] tabular-nums text-mathua-primary">{l.xp.toLocaleString()}</span>
              </li>
            ))}
          </ol>
        )}
      </section>

      <section className="mt-8" aria-labelledby="range-open">
        <div className="flex flex-wrap items-end justify-between gap-3 border-b border-mathua-border pb-2">
          <h2 id="range-open" className="font-mono text-xs uppercase tracking-wider text-mathua-primary">Open a league roster</h2>
          <span className="font-mono text-[10px] text-mathua-muted">Choose a week, then a league</span>
        </div>
        <div className="mt-3 flex flex-wrap items-center gap-2">
          <label className="flex items-center gap-2 font-mono text-[10px] text-mathua-muted">
            Week
            <select aria-label="Roster week" value={jumpWeek} onChange={e => setJumpWeek(e.target.value)} className="h-9 border border-mathua-border bg-mathua-surface px-2 font-mono text-xs text-mathua-primary">
              {data.weeks_available.map(w => <option key={w} value={w}>{weekLabel(w)}</option>)}
            </select>
          </label>
          <label className="flex items-center gap-2 font-mono text-[10px] text-mathua-muted">
            League
            <select
              aria-label="Roster league"
              defaultValue=""
              disabled={jumpLeagues.length === 0}
              onChange={e => {
                const [tier, league] = e.target.value.split('/')
                if (tier && league) onOpenWeek(jumpWeek, tier, Number(league))
              }}
              className="h-9 border border-mathua-border bg-mathua-surface px-2 font-mono text-xs text-mathua-primary"
            >
              <option value="" disabled>Choose a league</option>
              {TIER_ORDER.filter(tier => jumpLeagues.some(l => l.tier === tier)).map(tier => (
                <optgroup key={tier} label={tierLabel(tier)}>
                  {jumpLeagues.filter(l => l.tier === tier).map(l => (
                    <option key={`${l.tier}/${l.league_no}`} value={`${l.tier}/${l.league_no}`}>
                      {tierLabel(l.tier)} league {l.league_no} · {l.size} of {LEAGUE_SIZE}
                    </option>
                  ))}
                </optgroup>
              ))}
            </select>
          </label>
        </div>
      </section>
    </>
  )
}

export default function AdminLeaguesPage() {
  const [week, setWeek] = useState<string | null>(null)
  const [leagues, setLeagues] = useState<AdminLeague[] | null>(null)
  const [seats, setSeats] = useState(0)
  const [promote, setPromote] = useState(0)
  const [demote, setDemote] = useState(0)
  const [error, setError] = useState('')
  const [selected, setSelected] = useState<Selected | null>(null)
  const [roster, setRoster] = useState<AdminLeagueMember[]>([])
  const [rosterLoading, setRosterLoading] = useState(false)
  const [rosterError, setRosterError] = useState('')
  const [retryToken, setRetryToken] = useState(0)
  const [scope, setScope] = useState<Scope>('week')
  const [month, setMonth] = useState(() => new Date().toISOString().slice(0, 7))
  const [year, setYear] = useState(() => String(new Date().getUTCFullYear()))
  const [custom, setCustom] = useState<DateRange | undefined>(undefined)
  const [rangeData, setRangeData] = useState<AdminLeagueRange | null>(null)
  const [rangeLoading, setRangeLoading] = useState(false)
  const [rangeError, setRangeError] = useState('')
  const [rangeRetry, setRangeRetry] = useState(0)
  const [rankBy, setRankBy] = useState<'range' | 'latest'>('range')

  // The static export cannot know today's Monday at build time, so the week is
  // read from the URL in the browser and only then fetched.
  useEffect(() => {
    const params = new URLSearchParams(window.location.search)
    const rawScope = params.get('scope')
    if (rawScope === 'week' || rawScope === 'month' || rawScope === 'year' || rawScope === 'custom') setScope(rawScope)
    const rawMonth = params.get('month')
    if (rawMonth && /^\d{4}-\d{2}$/.test(rawMonth)) setMonth(rawMonth)
    const rawYear = params.get('year')
    if (rawYear && /^\d{4}$/.test(rawYear)) setYear(rawYear)
    const rawFrom = params.get('from')
    const rawTo = params.get('to')
    if (rawFrom && rawTo) setCustom({ from: new Date(`${rawFrom}T00:00:00Z`), to: new Date(`${rawTo}T00:00:00Z`) })
    const rawWeek = params.get('week')
    const parsedWeek = rawWeek && /^\d{4}-\d{2}-\d{2}$/.test(rawWeek) ? rawWeek : weekKeyOf(new Date())
    const rawTier = params.get('tier')
    const rawLeague = Number(params.get('league'))
    setWeek(parsedWeek)
    if (rawTier && TIER_ORDER.includes(rawTier) && Number.isInteger(rawLeague) && rawLeague > 0) {
      setSelected({ tier: rawTier, league: rawLeague })
    }
  }, [])

  const load = useCallback((key: string) => {
    setError('')
    getAdminLeagues({ week: key })
      .then((res) => {
        setLeagues(res.leagues)
        setSeats(res.seats)
        setPromote(res.promote)
        setDemote(res.demote)
        setWeek(res.week)
      })
      .catch((e) => {
        setError(getErrorMessage(e))
        setLeagues([])
      })
  }, [])

  useEffect(() => { if (week) load(week) }, [week, load])

  // Keep the URL as the source of truth for a refresh or a shared link.
  useEffect(() => {
    if (!week) return
    const params = new URLSearchParams()
    params.set('scope', scope)
    if (scope === 'week') {
      params.set('week', week)
      if (selected) {
        params.set('tier', selected.tier)
        params.set('league', String(selected.league))
      }
    } else if (scope === 'month') {
      params.set('month', month)
    } else if (scope === 'year') {
      params.set('year', year)
    } else if (custom?.from && custom?.to) {
      params.set('from', custom.from.toISOString().slice(0, 10))
      params.set('to', custom.to.toISOString().slice(0, 10))
    }
    window.history.replaceState(null, '', `?${params.toString()}`)
  }, [week, selected, scope, month, year, custom])

  // A week with leagues should open on something rather than an empty roster.
  useEffect(() => {
    if (!leagues || leagues.length === 0) { setSelected(null); return }
    const exists = selected && leagues.some(l => l.tier === selected.tier && l.league_no === selected.league)
    if (!exists) setSelected({ tier: leagues[0].tier, league: leagues[0].league_no })
  }, [leagues, selected])

  useEffect(() => {
    if (!week || !selected) { setRoster([]); return }
    let live = true
    setRosterLoading(true)
    setRosterError('')
    getAdminLeagues({ week, tier: selected.tier, league: selected.league })
      .then((res) => {
        if (live) setRoster(res.leagues[0]?.members ?? [])
      })
      .catch((e) => {
        if (live) { setRoster([]); setRosterError(getErrorMessage(e)) }
      })
      .finally(() => { if (live) setRosterLoading(false) })
    return () => { live = false }
  }, [week, selected, retryToken])

  const rangeBounds = useMemo(() => {
    if (scope === 'month') {
      const [y, m] = month.split('-').map(Number)
      const first = new Date(Date.UTC(y, m - 1, 1))
      const last = new Date(Date.UTC(y, m, 0))
      return { from: weekKeyOf(first), to: weekKeyOf(last), group: 'week' as const }
    }
    if (scope === 'year') {
      return { from: weekKeyOf(new Date(Date.UTC(Number(year), 0, 1))), to: weekKeyOf(new Date(Date.UTC(Number(year), 11, 31))), group: 'month' as const }
    }
    if (scope === 'custom' && custom?.from && custom?.to) {
      const days = Math.round((custom.to.getTime() - custom.from.getTime()) / 86400000)
      return { from: weekKeyOf(custom.from), to: weekKeyOf(custom.to), group: days > 90 ? 'month' as const : 'week' as const }
    }
    return null
  }, [scope, month, year, custom])

  useEffect(() => {
    if (scope === 'week' || !rangeBounds) { setRangeData(null); return }
    let live = true
    setRangeLoading(true)
    setRangeError('')
    getAdminLeagueRange(rangeBounds.from, rangeBounds.to, rangeBounds.group)
      .then(data => { if (live) setRangeData(data) })
      .catch(e => { if (live) { setRangeData(null); setRangeError(getErrorMessage(e)) } })
      .finally(() => { if (live) setRangeLoading(false) })
    return () => { live = false }
  }, [scope, rangeBounds, rangeRetry])

  const byTier = useMemo(() => {
    const map = new Map<string, AdminLeague[]>()
    for (const l of leagues ?? []) {
      const bucket = map.get(l.tier)
      if (bucket) bucket.push(l)
      else map.set(l.tier, [l])
    }
    return map
  }, [leagues])

  const selectedLeague = useMemo(
    () => leagues?.find(l => selected && l.tier === selected.tier && l.league_no === selected.league) ?? null,
    [leagues, selected],
  )

  const capacity = (leagues?.length ?? 0) * LEAGUE_SIZE
  const fill = capacity > 0 ? Math.round((seats / capacity) * 100) : 0
  const largest = leagues?.reduce<AdminLeague | null>((best, l) => (!best || l.size > best.size ? l : best), null)
  const thisWeek = week === weekKeyOf(new Date())

  function choose(l: AdminLeague) {
    setSelected({ tier: l.tier, league: l.league_no })
  }

  return (
    <AdminShell>
      <main className="mx-auto w-full max-w-[1200px] min-w-0 px-4 sm:px-6 py-8 sm:py-12 overflow-x-hidden">
        <SectionHeader label="Administration" title="Leagues" />

        <div className="mb-4 flex flex-wrap items-center gap-1" role="group" aria-label="League time scope">
          {(['week', 'month', 'year', 'custom'] as Scope[]).map(item => (
            <button
              key={item}
              type="button"
              aria-pressed={scope === item}
              onClick={() => setScope(item)}
              className={`h-9 border px-3 font-mono text-xs capitalize ${scope === item ? 'border-mathua-blue bg-mathua-blue-faint text-mathua-blue' : 'border-mathua-border text-mathua-secondary hover:border-mathua-blue'}`}
            >
              {item}
            </button>
          ))}
          {scope === 'month' && (
            <input type="month" aria-label="League month" value={month} onChange={e => setMonth(e.target.value)} className="ml-2 h-9 border border-mathua-border bg-mathua-surface px-2 font-mono text-xs text-mathua-primary" />
          )}
          {scope === 'year' && (
            <select aria-label="League year" value={year} onChange={e => setYear(e.target.value)} className="ml-2 h-9 border border-mathua-border bg-mathua-surface px-2 font-mono text-xs text-mathua-primary">
              {Array.from({ length: 5 }, (_, i) => String(new Date().getUTCFullYear() - i)).map(y => <option key={y} value={y}>{y}</option>)}
            </select>
          )}
          {scope === 'custom' && (
            <span className="ml-2"><DateRangePicker value={custom} onChange={setCustom} max={new Date()} /></span>
          )}
        </div>

        {scope !== 'week' && (
          <RangePanel
            data={rangeData}
            loading={rangeLoading}
            error={rangeError}
            onRetry={() => setRangeRetry(n => n + 1)}
            rankBy={rankBy}
            onRankBy={setRankBy}
            onOpenWeek={(nextWeek, tier, league) => {
              setScope('week')
              setWeek(nextWeek)
              setSelected({ tier, league })
            }}
          />
        )}

        {scope === 'week' && (<>
        <div className="mb-5 flex flex-wrap items-center gap-2">
          <button
            type="button"
            disabled={!week}
            onClick={() => week && load(shiftWeek(week, -7))}
            className="h-9 border border-mathua-border px-3 font-mono text-xs text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue disabled:opacity-50"
          >
            Previous
          </button>
          <button
            type="button"
            disabled={!week || thisWeek}
            onClick={() => week && load(shiftWeek(week, 7))}
            className="h-9 border border-mathua-border px-3 font-mono text-xs text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue disabled:opacity-50"
          >
            Next
          </button>
          <select
            aria-label="League week"
            value={week ?? ''}
            disabled={!week}
            onChange={e => load(e.target.value)}
            className="h-9 border border-mathua-border bg-mathua-surface px-2 font-mono text-xs text-mathua-primary"
          >
            {week && recentWeeks(weekKeyOf(new Date()), 12).concat(week).filter((w, i, all) => all.indexOf(w) === i).sort().reverse().map(w => (
              <option key={w} value={w}>{weekLabel(w)}</option>
            ))}
          </select>
          {!thisWeek && week && (
            <button
              type="button"
              onClick={() => load(weekKeyOf(new Date()))}
              className="h-9 px-2 font-mono text-[11px] text-mathua-blue hover:underline"
            >
              This week
            </button>
          )}
          <span className="ml-auto font-mono text-[10px] text-mathua-muted">
            {week ? `Week of ${weekLabel(week)} (UTC)` : 'Loading the week'}
          </span>
        </div>

        {error && (
          <div className="mb-4 border border-mathua-red bg-mathua-surface p-3">
            <p role="alert" className="font-mono text-[11px] text-mathua-red">{error}</p>
            <button type="button" onClick={() => week && load(week)} className="mt-2 h-8 border border-mathua-border px-2.5 font-mono text-[11px] text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue">
              Retry
            </button>
          </div>
        )}

        {!leagues && !error && (
          <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
            {Array.from({ length: 4 }, (_, i) => <Skeleton key={i} className="h-24 w-full" />)}
          </div>
        )}

        {leagues && leagues.length === 0 && !error && (
          <p role="status" className="border border-mathua-border bg-mathua-surface p-4 font-mono text-xs text-mathua-secondary">
            No league was formed in this week. Cohorts appear once a learner has earned XP.
          </p>
        )}

        {leagues && leagues.length > 0 && (
          <>
            <div className="grid gap-2 sm:grid-cols-2 lg:grid-cols-4">
              <Card><CardHeader><CardTitle>Leagues formed</CardTitle></CardHeader><CardContent><p className="font-mono text-2xl tabular-nums text-mathua-primary">{leagues.length}</p><p className="mt-1 font-mono text-[10px] text-mathua-muted">Across {TIER_ORDER.filter(t => byTier.has(t)).length} tiers</p></CardContent></Card>
              <Card><CardHeader><CardTitle>Learners seated</CardTitle></CardHeader><CardContent><p className="font-mono text-2xl tabular-nums text-mathua-primary">{seats.toLocaleString()}</p><p className="mt-1 font-mono text-[10px] text-mathua-muted">of {capacity.toLocaleString()} seats · {fill}% full</p></CardContent></Card>
              <Card><CardHeader><CardTitle>Movement rule</CardTitle></CardHeader><CardContent><p className="font-mono text-2xl tabular-nums text-mathua-primary">{promote} / {demote}</p><p className="mt-1 font-mono text-[10px] text-mathua-muted">Promoted and relegated per full league</p></CardContent></Card>
              <Card><CardHeader><CardTitle>Largest league</CardTitle></CardHeader><CardContent><p className="font-mono text-2xl tabular-nums text-mathua-primary">{largest?.size ?? 0}</p><p className="mt-1 font-mono text-[10px] text-mathua-muted">{largest ? `${tierLabel(largest.tier)} league ${largest.league_no}` : 'No league'}</p></CardContent></Card>
            </div>

            <section className="mt-8" aria-labelledby="pyramid-heading">
              <div id="pyramid-heading">
                <TierViews
                  label="Seats and fill by tier"
                  rows={TIER_ORDER.map(tier => {
                    const items = byTier.get(tier) ?? []
                    const tierSeats = items.reduce((n, l) => n + l.size, 0)
                    const first = items[0]
                    return {
                      tier,
                      leagues: items.length,
                      seats: tierSeats,
                      capacity: items.length * LEAGUE_SIZE,
                      unique: 0,
                      onSelect: first ? () => choose(first) : undefined,
                    }
                  })}
                />
              </div>
            </section>

            <section className="mt-8" aria-labelledby="rosters-heading">
              <div className="flex flex-wrap items-end justify-between gap-3 border-b border-mathua-border pb-2">
                <h2 id="rosters-heading" className="font-mono text-xs uppercase tracking-wider text-mathua-primary">League roster</h2>
                <label className="flex items-center gap-2 font-mono text-[10px] text-mathua-muted">
                  League
                  <select
                    aria-label="League roster"
                    value={selected ? `${selected.tier}/${selected.league}` : ''}
                    onChange={e => {
                      const [tier, league] = e.target.value.split('/')
                      setSelected({ tier, league: Number(league) })
                    }}
                    className="h-9 border border-mathua-border bg-mathua-surface px-2 font-mono text-xs text-mathua-primary"
                  >
                    {TIER_ORDER.filter(t => byTier.has(t)).map(tier => (
                      <optgroup key={tier} label={tierLabel(tier)}>
                        {(byTier.get(tier) ?? []).map(l => (
                          <option key={leagueKey(l)} value={`${l.tier}/${l.league_no}`}>
                            {tierLabel(l.tier)} league {l.league_no} · {l.size} of {LEAGUE_SIZE}
                          </option>
                        ))}
                      </optgroup>
                    ))}
                  </select>
                </label>
              </div>
              <Card className="mt-3">
                <CardContent className="pt-4">
                  <RosterPanel
                    league={selectedLeague}
                    members={roster}
                    loading={rosterLoading}
                    error={rosterError}
                    onRetry={() => setRetryToken(t => t + 1)}
                  />
                </CardContent>
              </Card>
            </section>
          </>
        )}
        </>)}
      </main>
    </AdminShell>
  )
}
