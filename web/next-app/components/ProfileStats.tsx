'use client'

import type { Scores } from '../lib/api'

interface Props {
  name: string
  scores: Scores
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

const mono: React.CSSProperties = { fontFamily: "'IBM Plex Mono', monospace" }
const serif: React.CSSProperties = { fontFamily: "'IBM Plex Serif', serif" }
const muted: React.CSSProperties = { ...mono, fontSize: 11, color: 'var(--text-muted)' }

export default function ProfileStats({ name, scores }: Props) {
  const xpPct = scores.daily_xp_goal > 0
    ? Math.min(100, Math.round((scores.xp_today / scores.daily_xp_goal) * 100))
    : 0
  const levelNum = LEVEL_NAMES[scores.level] ?? '·'

  return (
    <div
      style={{
        display: 'flex',
        flexWrap: 'wrap',
        gap: 20,
        alignItems: 'flex-start',
        border: '0.5px solid var(--border)',
        padding: '20px 24px',
      }}
    >
      {/* Level badge */}
      <div
        style={{
          border: '0.5px solid var(--border-strong)',
          padding: '14px 16px',
          textAlign: 'center',
          minWidth: 80,
        }}
      >
        <div style={{ ...mono, fontSize: 9, color: 'var(--text-muted)', marginBottom: 4 }}>
          LV
        </div>
        <div
          style={{
            ...mono,
            fontSize: 36,
            fontWeight: 400,
            color: 'var(--accent-blue)',
            lineHeight: 1,
          }}
        >
          {levelNum}
        </div>
        <div style={{ ...mono, fontSize: 10, color: 'var(--text-primary)', marginTop: 4 }}>
          {scores.level}
        </div>
      </div>

      {/* Stats */}
      <div style={{ flex: 1, minWidth: 240 }}>
        <h2 style={{ ...serif, fontSize: '1.2rem', fontWeight: 400, color: 'var(--text-primary)', marginBottom: 2 }}>
          {name}
        </h2>

        <div style={{ display: 'flex', flexWrap: 'wrap', gap: '12px 24px', marginTop: 10 }}>
          <div>
            <div style={muted}>XP TODAY</div>
            <div style={{ ...mono, fontSize: 16, color: 'var(--text-primary)' }}>
              {scores.xp_today}
              <span style={{ ...muted, fontSize: 11 }}> / {scores.daily_xp_goal}</span>
            </div>
            {/* XP bar */}
            <div
              style={{
                marginTop: 4,
                height: 4,
                width: '100%',
                maxWidth: 200,
                background: 'var(--border)',
              }}
            >
              <div
                style={{
                  height: '100%',
                  width: `${xpPct}%`,
                  background: 'var(--accent-blue)',
                  transition: 'width 0.3s',
                }}
              />
            </div>
          </div>

          <div>
            <div style={muted}>TOTAL XP</div>
            <div style={{ ...mono, fontSize: 16, color: 'var(--text-primary)' }}>
              {scores.lifetime_points.toLocaleString()}
            </div>
          </div>

          <div>
            <div style={muted}>STREAK</div>
            <div style={{ ...mono, fontSize: 16, color: scores.current_streak > 0 ? 'var(--accent-blue)' : 'var(--text-muted)' }}>
              {scores.current_streak} day{scores.current_streak !== 1 ? 's' : ''}
            </div>
          </div>

          <div>
            <div style={muted}>MASTERED</div>
            <div style={{ ...mono, fontSize: 16, color: 'var(--text-primary)' }}>
              {scores.concepts_mastered}
            </div>
          </div>

          <div>
            <div style={muted}>WEEKLY</div>
            <div style={{ ...mono, fontSize: 16, color: 'var(--text-primary)' }}>
              {scores.weekly_score.toLocaleString()}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
