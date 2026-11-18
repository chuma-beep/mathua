'use client'

import type { LeaderboardEntry } from '@/lib/api'
import { avatarURL, nameOf } from '@/lib/leaderboard-display'
import LeagueAvatar from '@/components/LeagueAvatar'

export default function Podium({ entries, score }: { entries: LeaderboardEntry[]; score: string }) {
  if (entries.length < 3) return null
  return (
    <ol className="mb-6 grid list-none grid-cols-1 gap-3 p-0 sm:grid-cols-3">
      {entries.slice(0, 3).map((row, i) => (
        <li
          key={`podium-${row.id}`}
          className={`flex min-w-0 items-center gap-3 border border-mathua-border bg-mathua-surface px-4 py-3 ${
            i === 0 ? 'sm:order-2' : i === 1 ? 'sm:order-1' : 'sm:order-3'
          }`}
        >
          <span className="w-5 shrink-0 font-mono text-lg text-mathua-blue">{row.rank}</span>
          <LeagueAvatar
            tier={row.league}
            weeklyXp={row.league_xp}
            seed={row.id}
            name={nameOf(row)}
            size={30}
            url={avatarURL(row)}
            className="shrink-0"
          />
          <span className="min-w-0 flex-1">
            <span className="block truncate text-sm text-mathua-primary" title={nameOf(row)}>
              {nameOf(row)}
            </span>
            <span className="block font-mono text-[10px] uppercase tracking-wider text-mathua-muted">
              {row.xp} XP {score}
            </span>
          </span>
        </li>
      ))}
    </ol>
  )
}
