'use client'

import Link from 'next/link'
import type { MyLeague } from '@/lib/api'
import { avatarURL, nameOf, shortDay } from '@/lib/leaderboard-display'
import { shiftWeekKey, weekLabel, weekStartKey } from '@/lib/leaderboard-period'
import LeagueAvatar from '@/components/LeagueAvatar'
import { LeaderboardTableSkeleton } from '@/components/skeletons/LeaderboardSkeleton'
import { Table, TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { PeriodArrows, WeekPicker } from './period-picker'

function outcomeText(league: MyLeague): string {
  switch (league.movement) {
    case 'up':
      return league.to_tier ? `Promoted to ${league.to_tier}` : 'Promoted'
    case 'down':
      return league.to_tier ? `Relegated to ${league.to_tier}` : 'Relegated'
    case 'stay':
      return league.to_tier ? `Stayed in ${league.to_tier}` : 'Stayed'
    case 'new':
      return 'New in this league'
    default:
      return ''
  }
}

export default function LeagueBoard({
  league,
  week,
  earliest,
  loading,
  onWeek,
}: {
  league: MyLeague | null
  week: string
  earliest?: string
  loading: boolean
  onWeek: (week: string) => void
}) {
  const current = weekStartKey(new Date())
  const canGoBack = !earliest || week > earliest
  const canGoForward = week < current
  const outcome = league && league.final ? outcomeText(league) : ''

  if (!loading && league === null) {
    return (
      <p className="px-2 py-8 text-center text-sm text-mathua-muted">
        <Link href="/login" className="text-mathua-blue hover:underline">
          Sign in
        </Link>{' '}
        to join a league.
      </p>
    )
  }

  if (!loading && league && league.tier === '') {
    return (
      <>
        <div className="mb-4 flex flex-wrap items-center justify-center gap-2">
          <PeriodArrows
            onPrev={() => onWeek(shiftWeekKey(week, -1))}
            onNext={() => onWeek(shiftWeekKey(week, 1))}
            prevDisabled={!canGoBack}
            nextDisabled={!canGoForward}
            label="week"
          />
          <WeekPicker value={week} onChange={onWeek} min={earliest} max={current} />
        </div>
        <p className="px-2 py-8 text-center text-sm text-mathua-muted">
          {league.final
            ? 'You were not in a league that week.'
            : 'You are not in a league yet. Answer one question this week and you join a Bronze league with room.'}
        </p>
      </>
    )
  }

  return (
    <div className="min-w-0">
      <div className="mb-4 flex flex-wrap items-center justify-between gap-3">
        <div className="flex flex-wrap items-center gap-2">
          <PeriodArrows
            onPrev={() => onWeek(shiftWeekKey(week, -1))}
            onNext={() => onWeek(shiftWeekKey(week, 1))}
            prevDisabled={!canGoBack}
            nextDisabled={!canGoForward}
            label="week"
          />
          <WeekPicker value={week} onChange={onWeek} min={earliest} max={current} />
        </div>
        {league && (
          <p className="font-mono text-[11px] text-mathua-muted">
            {shortDay(league.week_start)} to {shortDay(league.week_end)} · {league.size} learners
          </p>
        )}
      </div>

      {league && (
        <>
          <div className="mb-4 flex flex-wrap items-baseline justify-between gap-2">
            <h2 className="font-mono text-xs uppercase tracking-wider text-mathua-primary">
              {league.tier} league {league.league_no}
            </h2>
            <div className="flex items-center gap-2">
              {league.final ? (
                <span className="font-mono text-[10px] uppercase tracking-wider text-mathua-muted">Final</span>
              ) : (
                <span className="font-mono text-[10px] uppercase tracking-wider text-mathua-blue">Live</span>
              )}
              {outcome && <span className="font-mono text-[11px] text-mathua-secondary">{outcome}</span>}
            </div>
          </div>
          <p className="mb-4 text-sm text-mathua-secondary">
            {league.final
              ? 'This week is closed. The table is the result learners saw when it ended.'
              : `Top ${league.promote} move up a tier next Monday, bottom ${league.demote} move down. Ties on XP are ordered by name so positions are stable.`}
          </p>
        </>
      )}

      <Table className="min-w-[320px] border border-mathua-border bg-mathua-surface">
        <TableCaption className="sr-only">
          {weekLabel(week)} league standings
        </TableCaption>
        <TableHeader>
          <TableRow className="hover:bg-transparent">
            <TableHead className="w-10 bg-mathua-surface-elevated p-2.5 text-left sm:w-12 sm:p-[14px_20px]">#</TableHead>
            <TableHead className="bg-mathua-surface-elevated p-2.5 sm:p-[14px_20px]">Learner</TableHead>
            <TableHead className="bg-mathua-surface-elevated p-2.5 text-right sm:p-[14px_20px]">XP</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {loading && <LeaderboardTableSkeleton rows={8} />}
          {!loading &&
            league?.members.map((member) => (
              <TableRow key={member.id} className={member.is_me ? 'bg-mathua-blue-faint hover:bg-mathua-blue-faint' : undefined}>
                <TableCell className="p-2.5 text-sm text-mathua-muted sm:p-[14px_20px]">{member.rank}</TableCell>
                <TableCell className="max-w-[160px] p-2.5 text-sm font-medium text-mathua-primary sm:max-w-none sm:p-[14px_20px]">
                  <span className="flex min-w-0 items-center gap-2">
                    <LeagueAvatar
                      tier={league.tier}
                      weeklyXp={member.xp}
                      seed={member.id}
                      name={nameOf(member)}
                      size={26}
                      url={avatarURL(member)}
                      className="shrink-0"
                    />
                    <span className="block truncate" title={nameOf(member)}>
                      {nameOf(member)}
                    </span>
                    {member.is_me && (
                      <span className="shrink-0 font-mono text-[10px] uppercase text-mathua-blue">you</span>
                    )}
                  </span>
                </TableCell>
                <TableCell className="whitespace-nowrap p-2.5 text-right text-sm text-mathua-green sm:p-[14px_20px]">
                  {member.xp}
                </TableCell>
              </TableRow>
            ))}
        </TableBody>
      </Table>
    </div>
  )
}
