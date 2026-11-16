'use client'

import AdminShell from '@/components/admin/AdminShell'
import { useCallback, useEffect, useMemo, useState } from 'react'
import Link from 'next/link'
import SectionHeader from '../../../components/SectionHeader'
import Avatar from '../../../components/Avatar'
import { Badge } from '../../../components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '../../../components/ui/card'
import { Skeleton } from '../../../components/ui/skeleton'
import { Table, TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRow } from '../../../components/ui/table'
import { dicebearUrl, isDicebearStyle } from '../../../lib/dicebear'
import { API_BASE, getAdminLeagues, getErrorMessage, type AdminLeague, type AdminLeagueMember } from '../../../lib/api'

// The league console is a reading surface, not a queue. A week has five tiers,
// each with zero or more leagues, and a league has up to thirty seats. The page
// therefore separates the week summary, the tier pyramid and one roster, and
// reads only the selected roster. Nothing here changes a league: XP decides the
// seats and the server applies movement when the week closes.

const LEAGUE_SIZE = 30
const TIER_ORDER = ['bronze', 'silver', 'gold', 'platinum', 'diamond']
const TIER_CLASS: Record<string, string> = {
  bronze: 'bg-[var(--chart-2)]',
  silver: 'bg-[var(--chart-5)]',
  gold: 'bg-[var(--chart-3)]',
  platinum: 'bg-[var(--chart-1)]',
  diamond: 'bg-[var(--chart-4)]',
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

function MemberAvatar({ member }: { member: AdminLeagueMember }) {
  const pick = member.avatar_dicebear
  const storedURL = member.avatar_url
  const url = storedURL
    ? (storedURL.startsWith('/') ? `${API_BASE}${storedURL}` : storedURL)
    : (pick && isDicebearStyle(pick.style) ? dicebearUrl(pick.style, pick.seed) : undefined)
  return (
    <Avatar
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
            <TableRow key={m.student_id} className={up ? 'bg-mathua-green-faint' : down ? 'bg-mathua-red-faint' : undefined}>
              <TableCell className="text-mathua-secondary">{m.rank}</TableCell>
              <TableCell>
                <span className="flex min-w-0 items-center gap-2">
                  <MemberAvatar member={m} />
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
                {up ? <Badge variant="up">Up</Badge> : down ? <Badge variant="down">Down</Badge> : <Badge variant="stay">Stays</Badge>}
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

  // The static export cannot know today's Monday at build time, so the week is
  // read from the URL in the browser and only then fetched.
  useEffect(() => {
    const params = new URLSearchParams(window.location.search)
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
    params.set('week', week)
    if (selected) {
      params.set('tier', selected.tier)
      params.set('league', String(selected.league))
    }
    window.history.replaceState(null, '', `?${params.toString()}`)
  }, [week, selected])

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
              <div className="flex flex-wrap items-baseline justify-between gap-2 border-b border-mathua-border pb-2">
                <h2 id="pyramid-heading" className="font-mono text-xs uppercase tracking-wider text-mathua-primary">Tier pyramid</h2>
                <span className="font-mono text-[10px] text-mathua-muted">Seats and fill by tier</span>
              </div>
              <div className="mt-3 grid grid-cols-2 gap-2 sm:grid-cols-3 lg:grid-cols-5">
                {TIER_ORDER.map(tier => {
                  const items = byTier.get(tier) ?? []
                  const tierSeats = items.reduce((n, l) => n + l.size, 0)
                  const tierCapacity = items.length * LEAGUE_SIZE
                  const tierFill = tierCapacity > 0 ? Math.round((tierSeats / tierCapacity) * 100) : 0
                  const first = items[0]
                  return (
                    <button
                      key={tier}
                      type="button"
                      disabled={!first}
                      onClick={() => first && choose(first)}
                      className={`border p-3 text-left transition-colors ${first ? 'border-mathua-border bg-mathua-surface hover:border-mathua-blue' : 'border-mathua-border bg-mathua-surface opacity-55'}`}
                    >
                      <div className="flex items-center justify-between gap-2">
                        <span className="font-mono text-[10px] uppercase tracking-wider text-mathua-secondary">{tierLabel(tier)}</span>
                        <span className={`h-2 w-2 ${TIER_CLASS[tier]}`} aria-hidden="true" />
                      </div>
                      <p className="mt-2 font-mono text-xl tabular-nums text-mathua-primary">{items.length}</p>
                      <p className="font-mono text-[10px] text-mathua-muted">league{items.length === 1 ? '' : 's'}</p>
                      <p className="mt-2 font-mono text-[10px] text-mathua-secondary">{tierSeats} / {tierCapacity} seats</p>
                      <p className="font-mono text-[10px] text-mathua-muted">{tierFill}% full</p>
                    </button>
                  )
                })}
              </div>
              <div className="mt-3 flex h-2 w-full overflow-hidden border border-mathua-border" aria-hidden="true">
                {TIER_ORDER.map(tier => {
                  const tierSeats = (byTier.get(tier) ?? []).reduce((n, l) => n + l.size, 0)
                  const width = seats > 0 ? (tierSeats / seats) * 100 : 0
                  return <span key={tier} className={TIER_CLASS[tier]} style={{ width: `${width}%` }} />
                })}
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
      </main>
    </AdminShell>
  )
}
