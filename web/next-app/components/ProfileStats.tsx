'use client'

import type { Scores } from '../lib/api'
import Avatar from './Avatar'

interface Props {
  name: string
  scores: Scores
  avatarSeed?: string
  avatarUrl?: string
  avatarPreset?: number | null
}

const LEVEL_NAMES: Record<string, string> = {
  'Novice': '01',
  'Apprentice': '02',
  'Student': '03',
  'Scholar': '04',
  'Adept': '05',
  'Expert': '06',
  'Master': '07',
  'Grandmaster': '08',
  'Math Architect': '09',
}

export default function ProfileStats({ name, scores, avatarSeed, avatarUrl, avatarPreset }: Props) {
  const xpPct = scores.daily_xp_goal > 0
    ? Math.min(100, Math.round((scores.xp_today / scores.daily_xp_goal) * 100))
    : 0
  const levelNum = LEVEL_NAMES[scores.level] ?? '·'

  return (
    <div className="flex flex-col sm:flex-row flex-wrap gap-4 sm:gap-5 items-start border-[0.5px] border-mathua-border p-4 sm:p-5 sm:px-6 w-full max-w-full min-w-0 overflow-hidden">
      <div className="flex flex-row gap-4 sm:contents w-full sm:w-auto">
        <div className="shrink-0 self-start">
          <Avatar seed={avatarSeed ?? name} name={name} size={72} url={avatarUrl} preset={avatarPreset ?? undefined} />
        </div>
        <div className="border-[0.5px] border-mathua-border-strong px-4 py-3.5 text-center min-w-[72px] sm:min-w-[80px] shrink-0 self-start">
          <div className="font-mono text-[9px] text-mathua-muted mb-1">
            LV
          </div>
          <div className="font-mono text-[36px] font-normal text-mathua-blue leading-none">
            {levelNum}
          </div>
          <div className="font-mono text-[10px] text-mathua-primary mt-1">
            {scores.level}
          </div>
        </div>
      </div>

      <div className="flex-1 min-w-0 w-full sm:min-w-0">
        <h2 className="font-serif text-[1.2rem] font-normal text-mathua-primary mb-0.5 truncate">
          {name}
        </h2>

        <div className="flex flex-wrap gap-x-6 gap-y-3 mt-2.5">
          <div className="min-w-0">
            <div className="font-mono text-[11px] text-mathua-muted">XP TODAY</div>
            <div className="font-mono text-[16px] text-mathua-primary">
              {scores.xp_today}
              <span className="font-mono text-[11px] text-mathua-muted"> / {scores.daily_xp_goal}</span>
            </div>
            <div className="mt-1 h-1 w-full max-w-[200px] bg-mathua-border">
              <div
                className="h-full bg-mathua-blue transition-all duration-300"
                style={{ width: `${xpPct}%` }}
              />
            </div>
          </div>

          <div className="min-w-0">
            <div className="font-mono text-[11px] text-mathua-muted">TOTAL XP</div>
            <div className="font-mono text-[16px] text-mathua-primary">
              {scores.lifetime_points.toLocaleString()}
            </div>
          </div>

          <div className="min-w-0">
            <div className="font-mono text-[11px] text-mathua-muted">STREAK</div>
            <div className={`font-mono text-[16px] ${scores.current_streak > 0 ? 'text-mathua-blue' : 'text-mathua-muted'}`}>
              {scores.current_streak} day{scores.current_streak !== 1 ? 's' : ''}
            </div>
          </div>

          <div className="min-w-0">
            <div className="font-mono text-[11px] text-mathua-muted">MASTERED</div>
            <div className="font-mono text-[16px] text-mathua-primary">
              {scores.concepts_mastered}
            </div>
          </div>

          <div className="min-w-0">
            <a href="/leaderboard" className="font-mono text-[11px] text-mathua-muted hover:text-mathua-blue no-underline inline-flex items-center min-h-[20px]">
              WEEKLY →
            </a>
            <div className="font-mono text-[16px] text-mathua-primary">
              {scores.weekly_score.toLocaleString()}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
