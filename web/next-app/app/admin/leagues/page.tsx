'use client'

import AdminShell from '@/components/admin/AdminShell'
import { useCallback, useEffect, useMemo, useState } from 'react'
import Link from 'next/link'
import SectionHeader from '../../../components/SectionHeader'
import Loading from '../../../components/Loading'
import Avatar from '../../../components/Avatar'
import { dicebearUrl, isDicebearStyle } from '../../../lib/dicebear'
import { API_BASE, getAdminLeagues, getErrorMessage, type AdminLeague, type AdminLeagueMember } from '../../../lib/api'

// The league pyramid, week by week.
//
// The public league page answers one question: where am I. That leaves an owner unable to see a
// cohort at all, which is how a report saying "my league is broken" becomes unanswerable: there
// was no screen on which the league existed. This page is that screen.
//
// Two reads, deliberately. The pyramid is one grouped query, so opening the page does not load
// every roster; a roster is fetched when a league is opened. A week can hold dozens of leagues
// and each holds up to thirty learners, so eager loading would be a page that gets slower as the
// product succeeds.
//
// The week is held as null until the client knows today's date. The site is a static export, so
// rendering "this week" during the build would bake whichever Monday the build ran on into the
// HTML and then disagree with the browser on hydration.

/** The Monday of the league week containing `d`, as the server keys it (UTC). */
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

/** `2026-10-19` as `19 October 2026`. Fixed to UTC so the label matches the key. */
function weekLabel(week: string): string {
  const t = new Date(`${week}T00:00:00Z`)
  return t.toLocaleDateString('en-GB', { day: 'numeric', month: 'long', year: 'numeric', timeZone: 'UTC' })
}

const TIER_ORDER = ['bronze', 'silver', 'gold', 'platinum', 'diamond']

function tierLabel(tier: string): string {
  return tier.charAt(0).toUpperCase() + tier.slice(1)
}

function memberLabel(m: AdminLeagueMember): string {
  return m.name || m.username || m.student_id
}

function MemberAvatar({ member }: { member: AdminLeagueMember }) {
  // A custom upload arrives as a public by-id URL, so it wins over the remote
  // Google photo and the generated pick in the same order the rest of the app
  // uses. A legacy preset has no URL and becomes the Avatar's initial palette.
  const pick = member.avatar_dicebear
  // A custom upload comes back as a server-relative path because the API and
  // the web export are separate origins. A Google photo or any other absolute
  // URL is already complete.
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

export default function AdminLeaguesPage() {
  const [week, setWeek] = useState<string | null>(null)
  const [leagues, setLeagues] = useState<AdminLeague[] | null>(null)
  const [seats, setSeats] = useState(0)
  const [promote, setPromote] = useState(0)
  const [error, setError] = useState('')
  const [rosters, setRosters] = useState<Record<string, AdminLeagueMember[]>>({})
  const [open, setOpen] = useState('')
  const [rosterError, setRosterError] = useState('')

  // Today's Monday, read once in the browser rather than during the build.
  useEffect(() => { setWeek(weekKeyOf(new Date())) }, [])

  const load = useCallback((key: string) => {
    setError('')
    getAdminLeagues({ week: key })
      .then((res) => { setLeagues(res.leagues); setSeats(res.seats); setPromote(res.promote); setWeek(res.week) })
      .catch((e) => { setError(getErrorMessage(e)); setLeagues([]) })
  }, [])

  useEffect(() => { if (week) load(week) }, [week, load])

  async function toggle(league: AdminLeague) {
    const key = `${league.tier}/${league.league_no}`
    setRosterError('')
    if (open === key) { setOpen(''); return }
    setOpen(key)
    if (rosters[key]) return
    try {
      const res = await getAdminLeagues({ week: week ?? undefined, tier: league.tier, league: league.league_no })
      const roster = res.leagues[0]?.members ?? []
      setRosters((prev) => ({ ...prev, [key]: roster }))
    } catch (e) {
      // Collapse as well as report. Leaving the league open would show its loading state for a
      // request that has already failed, which reads as "still working" rather than "cannot read".
      setRosterError(getErrorMessage(e))
      setOpen('')
    }
  }

  // Grouped for rendering, in pyramid order. A tier with no league this week is omitted rather
  // than shown empty: "no gold league exists yet" is a fact about the week, not about the tier.
  const tiers = useMemo(() => {
    const byTier = new Map<string, AdminLeague[]>()
    for (const l of leagues ?? []) {
      const bucket = byTier.get(l.tier)
      if (bucket) bucket.push(l)
      else byTier.set(l.tier, [l])
    }
    return TIER_ORDER.filter((t) => byTier.has(t)).map((t) => [t, byTier.get(t) ?? []] as const)
  }, [leagues])

  const thisWeek = week ? week === weekKeyOf(new Date()) : false

  return (
    <AdminShell>
      <main className="mx-auto w-full max-w-[1200px] min-w-0 px-4 sm:px-6 py-8 sm:py-12 overflow-x-hidden">
        <SectionHeader label="Administration" title="Leagues" />
  
        <div className="mb-4 flex flex-wrap items-center gap-3">
          <button
            type="button"
            disabled={!week}
            onClick={() => week && load(shiftWeek(week, -7))}
            className="border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue px-3 h-9 font-mono text-xs disabled:opacity-50"
          >
            Earlier week
          </button>
          <p className="font-mono text-xs text-mathua-primary">
            {week ? `Week of ${weekLabel(week)}` : 'Loading the week'}
          </p>
          <button
            type="button"
            disabled={!week || thisWeek}
            onClick={() => week && load(shiftWeek(week, 7))}
            className="border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue px-3 h-9 font-mono text-xs disabled:opacity-50"
          >
            Later week
          </button>
        </div>

        {error && (
          <p role="alert" className="mb-4 border border-mathua-red bg-mathua-surface p-4 font-mono text-xs text-mathua-red">{error}</p>
        )}

        {!leagues && !error && <Loading label="LOADING LEAGUES" />}

        {leagues && leagues.length === 0 && !error && (
          <p role="status" className="border border-mathua-border bg-mathua-surface p-4 font-mono text-xs text-mathua-secondary">
            No league was formed in this week. Cohorts appear once a learner has earned XP.
          </p>
        )}

        {leagues && leagues.length > 0 && (
          <p className="mb-6 font-mono text-[11px] text-mathua-muted">
            {leagues.length} league{leagues.length === 1 ? '' : 's'}, {seats} learner{seats === 1 ? '' : 's'} seated.
            {' '}A full league moves {promote} up and {promote} down; short leagues move a third of themselves and cohorts under six do not move.
          </p>
        )}

        {rosterError && (
          <p role="alert" className="mb-4 border border-mathua-red bg-mathua-surface p-4 font-mono text-xs text-mathua-red">{rosterError}</p>
        )}

        <div className="space-y-8">
          {tiers.map(([tier, items]) => (
            <section key={tier} className="space-y-3">
              <div className="flex flex-wrap items-baseline justify-between gap-2 border-b border-mathua-border pb-2">
                <h2 className="font-mono text-xs text-mathua-primary">{tierLabel(tier)}</h2>
                <span className="font-mono text-[10px] text-mathua-muted">
                  {items.length} league{items.length === 1 ? '' : 's'}
                </span>
              </div>
              <ul className="space-y-2">
                {items.map((l) => {
                  const key = `${l.tier}/${l.league_no}`
                  const expanded = open === key
                  return (
                    <li key={key} className="border border-mathua-border bg-mathua-surface">
                      <button
                        type="button"
                        aria-expanded={expanded}
                        aria-label={`${tierLabel(l.tier)} league ${l.league_no}, ${l.size} learner${l.size === 1 ? '' : 's'}`}
                        onClick={() => { void toggle(l) }}
                        className="flex w-full flex-wrap items-baseline justify-between gap-2 px-3 py-3 text-left"
                      >
                        <span className="font-mono text-xs text-mathua-primary">League {l.league_no}</span>
                        <span className="font-mono text-[10px] text-mathua-muted">
                          {l.size} learner{l.size === 1 ? '' : 's'}
                          {l.movers > 0 ? `, ${l.movers} up and ${l.movers} down` : ', nobody moves'}
                        </span>
                      </button>
                      {expanded && (
                        <div className="border-t border-mathua-border px-3 py-3">
                          {!rosters[key] && <Loading label="LOADING STANDINGS" />}
                          {rosters[key] && rosters[key].length === 0 && (
                            <p className="font-mono text-[11px] text-mathua-secondary">Nobody holds a seat in this league.</p>
                          )}
                          {rosters[key] && rosters[key].length > 0 && (
                            <table className="w-full border-collapse font-mono text-[11px]">
                              <thead>
                                <tr className="text-mathua-muted">
                                  <th scope="col" className="py-1 pr-2 text-left font-normal">Rank</th>
                                  <th scope="col" className="py-1 pr-2 text-left font-normal">Learner</th>
                                  <th scope="col" className="py-1 text-right font-normal">XP this week</th>
                                </tr>
                              </thead>
                              <tbody>
                                {rosters[key].map((m) => (
                                  <tr key={m.student_id} className="border-t border-mathua-border">
                                    <td className="py-1 pr-2 text-mathua-secondary">{m.rank}</td>
                                    <td className="py-1 pr-2">
                                      <span className="flex min-w-0 items-center gap-2">
                                        <MemberAvatar member={m} />
                                        <span className="min-w-0">
                                          <Link
                                            href={`/admin/users/detail?id=${encodeURIComponent(m.student_id)}`}
                                            className="text-mathua-blue hover:underline"
                                          >
                                            {memberLabel(m)}
                                          </Link>
                                          {m.username && m.username !== memberLabel(m) && (
                                            <span className="ml-2 text-mathua-muted">@{m.username}</span>
                                          )}
                                        </span>
                                      </span>
                                    </td>
                                    <td className="py-1 text-right text-mathua-primary">{m.xp}</td>
                                  </tr>
                                ))}
                              </tbody>
                            </table>
                          )}
                        </div>
                      )}
                    </li>
                  )
                })}
              </ul>
            </section>
          ))}
        </div>
      </main>
    </AdminShell>
  )
}
