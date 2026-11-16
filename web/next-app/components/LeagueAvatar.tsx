'use client'

import Avatar from './Avatar'

// A league frame is a border, never a skin. The avatar image keeps its exact
// size; every SVG element lives in the gutter around it, so nothing
// draws over the picture. A learner with no tier or no XP this week gets the
// plain avatar back unchanged rather than a default Bronze ring.

export const LEAGUE_TIER_ORDER = ['bronze', 'silver', 'gold', 'platinum', 'diamond'] as const
export type LeagueTier = (typeof LEAGUE_TIER_ORDER)[number]

const MIN_GEOMETRY_SIZE = 20

const TIER_COLOUR: Record<LeagueTier, { main: string; soft: string }> = {
  bronze: { main: 'var(--league-bronze)', soft: 'var(--league-bronze-soft)' },
  silver: { main: 'var(--league-silver)', soft: 'var(--league-silver-soft)' },
  gold: { main: 'var(--league-gold)', soft: 'var(--league-gold-soft)' },
  platinum: { main: 'var(--league-platinum)', soft: 'var(--league-platinum-soft)' },
  diamond: { main: 'var(--league-diamond)', soft: 'var(--league-diamond-soft)' },
}

/** Normalises a payload tier. Anything outside the five tiers is no frame. */
export function normalizeLeagueTier(tier?: string | null): LeagueTier | null {
  const value = (tier ?? '').trim().toLowerCase()
  return (LEAGUE_TIER_ORDER as readonly string[]).includes(value) ? (value as LeagueTier) : null
}

export function leagueFrameGeometry(size: number) {
  const collapsed = size < MIN_GEOMETRY_SIZE
  const gutter = collapsed ? Math.max(2, Math.round(size * 0.08)) : Math.max(6, Math.round(size * 0.24))
  const outer = size + gutter * 2
  const padPercent = (gutter / outer) * 100
  const imageRadiusPercent = 50 - padPercent
  return {
    collapsed,
    gutter,
    outer,
    imageRadius: size / 2,
    // The closest any frame geometry comes to the image. Exported so the test
    // can assert the gutter is real rather than trusting the viewBox.
    frameInnerRadius: size / 2 + gutter * 0.2,
    padPercent,
    imageRadiusPercent,
  }
}

type Props = {
  tier?: string | null
  /** Earned XP. A learner with none keeps the plain avatar. */
  weeklyXp?: number | null
  seed: string
  name: string
  size?: number
  url?: string
  preset?: number | null
  className?: string
  avatarClassName?: string
}

function polar(radius: number, degrees: number): { x: number; y: number } {
  const radians = ((degrees - 90) * Math.PI) / 180
  return { x: 50 + radius * Math.cos(radians), y: 50 + radius * Math.sin(radians) }
}

function segment(radius: number, start: number, end: number): string {
  const a = polar(radius, start)
  const b = polar(radius, end)
  return `M ${a.x.toFixed(3)} ${a.y.toFixed(3)} A ${radius.toFixed(3)} ${radius.toFixed(3)} 0 0 1 ${b.x.toFixed(3)} ${b.y.toFixed(3)}`
}

function octagon(radius: number): string {
  return Array.from({ length: 8 }, (_, i) => {
    const p = polar(radius, i * 45)
    return `${i === 0 ? 'M' : 'L'} ${p.x.toFixed(3)} ${p.y.toFixed(3)}`
  }).join(' ') + ' Z'
}

function octagonPath(radius: number, rotation = 22.5): string {
  return Array.from({ length: 8 }, (_, i) => {
    const p = polar(radius, rotation + i * 45)
    return `${i === 0 ? 'M' : 'L'} ${p.x.toFixed(3)} ${p.y.toFixed(3)}`
  }).join(' ') + ' Z'
}

function Frame({ tier, size }: { tier: LeagueTier; size: number }) {
  const { collapsed, padPercent, imageRadiusPercent } = leagueFrameGeometry(size)
  const colour = TIER_COLOUR[tier]
  const pad = padPercent
  const inner = imageRadiusPercent
  const common = { fill: 'none', strokeLinecap: 'round' as const, strokeLinejoin: 'round' as const }
  const detailed = size >= 40
  const glow = detailed ? (
    <circle cx="50" cy="50" r={inner + pad * 0.86} fill="none" stroke={colour.soft} strokeWidth={pad * 0.24} opacity="0.16" />
  ) : null
  const dot = (angle: number, radius: number, radiusPx: number, fill = colour.soft) => {
    const p = polar(radius, angle)
    return <circle cx={p.x} cy={p.y} r={radiusPx} fill={fill} stroke={colour.main} strokeWidth={radiusPx * 0.4} />
  }

  if (collapsed) {
    return (
      <circle
        cx="50"
        cy="50"
        r={inner + pad * 0.5}
        stroke={colour.main}
        strokeWidth={pad * 0.55}
        {...common}
      />
    )
  }


  if (tier === 'bronze') {
    const r = inner + pad * 0.52
    return (
      <>
        {glow}
        <path d={octagonPath(r)} stroke={colour.main} strokeWidth={pad * 0.3} {...common} />
        <path d={octagonPath(inner + pad * 0.3)} stroke={colour.soft} strokeWidth={pad * 0.11} opacity="0.95" {...common} />
        {detailed && [22.5, 67.5, 112.5, 157.5, 202.5, 247.5, 292.5, 337.5].map((angle) => (
          <g key={angle}>{dot(angle, r, pad * 0.09)}</g>
        ))}
        <polygon
          points={`${50 - pad * 0.7},${50 + inner} ${50 - pad * 0.45},${50 + inner + pad * 0.28} ${50 - pad * 0.24},${50 + inner + pad * 0.16} ${50},${50 + inner + pad * 0.75} ${50 + pad * 0.24},${50 + inner + pad * 0.16} ${50 + pad * 0.45},${50 + inner + pad * 0.28} ${50 + pad * 0.7},${50 + inner}`}
          fill={colour.main}
          opacity="0.5"
          stroke={colour.soft}
          strokeWidth={pad * 0.08}
        />
        {detailed && <path d={`M ${50 - pad * 0.35} ${50 + inner + pad * 0.4} L 50 ${50 + inner + pad * 0.62} L ${50 + pad * 0.35} ${50 + inner + pad * 0.4}`} stroke={colour.soft} strokeWidth={pad * 0.06} fill="none" opacity="0.8" />}
      </>
    )
  }

  if (tier === 'silver') {
    const edge = 50 - inner
    return (
      <>
        {glow}
        <circle cx="50" cy="50" r={inner + pad * 0.68} stroke={colour.soft} strokeWidth={pad * 0.05} opacity="0.7" {...common} />
        <circle cx="50" cy="50" r={inner + pad * 0.3} stroke={colour.soft} strokeWidth={pad * 0.16} {...common} />
        <circle cx="50" cy="50" r={inner + pad * 0.55} stroke={colour.main} strokeWidth={pad * 0.18} {...common} />
        <polygon points={`${edge},42 ${pad * 0.08},36 ${pad * 0.3},50 ${pad * 0.08},64 ${edge},58`} fill={colour.soft} opacity="0.5" stroke={colour.main} strokeWidth={pad * 0.1} />
        <polygon points={`${100 - edge},42 ${100 - pad * 0.08},36 ${100 - pad * 0.3},50 ${100 - pad * 0.08},64 ${100 - edge},58`} fill={colour.soft} opacity="0.5" stroke={colour.main} strokeWidth={pad * 0.1} />
        {detailed && (
          <>
            <path d={`M ${pad * 0.16} 44 L ${edge - pad * 0.04} 48`} stroke={colour.main} strokeWidth={pad * 0.05} opacity="0.8" />
            <path d={`M ${pad * 0.16} 56 L ${edge - pad * 0.04} 52`} stroke={colour.main} strokeWidth={pad * 0.05} opacity="0.8" />
            <path d={`M ${100 - pad * 0.16} 44 L ${100 - edge + pad * 0.04} 48`} stroke={colour.main} strokeWidth={pad * 0.05} opacity="0.8" />
            <path d={`M ${100 - pad * 0.16} 56 L ${100 - edge + pad * 0.04} 52`} stroke={colour.main} strokeWidth={pad * 0.05} opacity="0.8" />
            {dot(0, inner + pad * 0.55, pad * 0.09)}
            {dot(90, inner + pad * 0.55, pad * 0.09)}
            {dot(180, inner + pad * 0.55, pad * 0.09)}
            {dot(270, inner + pad * 0.55, pad * 0.09)}
          </>
        )}
      </>
    )
  }

  if (tier === 'gold') {
    const base = 50 - inner
    return (
      <>
        {glow}
        <circle cx="50" cy="50" r={inner + pad * 0.32} stroke={colour.main} strokeWidth={pad * 0.22} {...common} />
        <circle cx="50" cy="50" r={inner + pad * 0.5} stroke={colour.soft} strokeWidth={pad * 0.06} opacity="0.9" {...common} />
        {detailed && [0, 45, 90, 135, 180, 225, 270, 315].map((angle) => (
          <g key={angle}>{dot(angle, inner + pad * 0.45, pad * 0.075)}</g>
        ))}
        <polygon
          points={`${50 - pad * 0.9},${base} ${50 - pad * 0.7},${base - pad * 0.5} ${50 - pad * 0.35},${base - pad * 0.28} ${50 - pad * 0.2},${pad * 0.12} ${50},${base - pad * 0.62} ${50 + pad * 0.2},${pad * 0.12} ${50 + pad * 0.35},${base - pad * 0.28} ${50 + pad * 0.7},${base - pad * 0.5} ${50 + pad * 0.9},${base}`}
          fill={colour.main}
          opacity="0.82"
          stroke={colour.soft}
          strokeWidth={pad * 0.08}
        />
        {detailed && (
          <>
            <circle cx={50} cy={pad * 0.18} r={pad * 0.055} fill={colour.soft} stroke={colour.main} strokeWidth={pad * 0.03} />
            <circle cx={50 - pad * 0.35} cy={base - pad * 0.28} r={pad * 0.05} fill={colour.soft} stroke={colour.main} strokeWidth={pad * 0.03} />
            <circle cx={50 + pad * 0.35} cy={base - pad * 0.28} r={pad * 0.05} fill={colour.soft} stroke={colour.main} strokeWidth={pad * 0.03} />
          </>
        )}
        <polygon points={`50,${50 + inner} ${50 + pad * 0.3},${50 + inner + pad * 0.28} 50,${50 + inner + pad * 0.62} ${50 - pad * 0.3},${50 + inner + pad * 0.28}`} fill={colour.soft} stroke={colour.main} strokeWidth={pad * 0.07} />
      </>
    )
  }

  if (tier === 'platinum') {
    const edge = 50 - inner
    const r = inner + pad * 0.34
    return (
      <>
        {glow}
        <circle cx="50" cy="50" r={inner + pad * 0.72} stroke={colour.soft} strokeWidth={pad * 0.07} opacity="0.55" {...common} />
        <path d={octagonPath(inner + pad * 0.55, 0)} stroke={colour.main} strokeWidth={pad * 0.2} {...common} />
        {Array.from({ length: 16 }, (_, i) => {
          const centre = i * 22.5
          return (
            <path
              key={i}
              d={segment(r, centre - 8, centre + 8)}
              stroke={colour.main}
              strokeWidth={pad * 0.3}
              {...common}
            />
          )
        })}
        <polygon points={`${edge},44 ${pad * 0.06},38 ${pad * 0.24},50 ${pad * 0.06},62 ${edge},56`} fill={colour.main} opacity="0.35" stroke={colour.soft} strokeWidth={pad * 0.07} />
        <polygon points={`${100 - edge},44 ${100 - pad * 0.06},38 ${100 - pad * 0.24},50 ${100 - pad * 0.06},62 ${100 - edge},56`} fill={colour.main} opacity="0.35" stroke={colour.soft} strokeWidth={pad * 0.07} />
        {detailed && [45, 135, 225, 315].map((angle) => {
          const p = polar(inner + pad * 0.78, angle)
          const a = polar(inner + pad * 0.52, angle - 9)
          const b = polar(inner + pad * 0.52, angle + 9)
          return <polygon key={angle} points={`${a.x},${a.y} ${p.x},${p.y} ${b.x},${b.y}`} fill={colour.main} opacity="0.45" stroke={colour.soft} strokeWidth={pad * 0.05} />
        })}
      </>
    )
  }

  const radius = inner + pad * 0.58
  const bevelRadius = inner + pad * 0.34
  return (
    <>
      {glow}
      <path d={octagonPath(radius)} stroke={colour.main} strokeWidth={pad * 0.28} {...common} />
      <path d={octagonPath(bevelRadius)} stroke={colour.soft} strokeWidth={pad * 0.12} {...common} />
      {[0, 45, 90, 135, 180, 225, 270, 315].map((angle) => {
        const a = polar(bevelRadius, angle)
        const b = polar(radius - pad * 0.1, angle)
        return <path key={angle} d={`M ${a.x.toFixed(3)} ${a.y.toFixed(3)} L ${b.x.toFixed(3)} ${b.y.toFixed(3)}`} stroke={angle % 90 === 0 ? colour.main : angle % 45 === 0 ? 'var(--league-diamond-violet)' : colour.soft} strokeWidth={pad * 0.07} opacity="0.85" />
      })}
      <path d="M 50 1.4 L 54.5 5.9 L 50 10.4 L 45.5 5.9 Z" fill={colour.soft} stroke={colour.main} strokeWidth={pad * 0.06} />
      <path d="M 50 2.4 L 50 9.4" stroke="var(--league-diamond-violet)" strokeWidth={pad * 0.05} opacity="0.8" />
      <path d={`M 50 ${100 - 10.4} L 54.5 ${100 - 5.9} L 50 ${100 - 1.4} L 45.5 ${100 - 5.9} Z`} fill={colour.soft} stroke={colour.main} strokeWidth={pad * 0.06} />
      {detailed && [0, 90, 180, 270].map((angle) => {
        const a = polar(radius - pad * 0.02, angle)
        const b = polar(radius + pad * 0.09, angle - 6)
        const c = polar(radius + pad * 0.09, angle + 6)
        const d = polar(radius - pad * 0.02, angle + 12)
        return <polygon key={angle} points={`${a.x},${a.y} ${b.x},${b.y} ${c.x},${c.y} ${d.x},${d.y}`} fill={colour.soft} opacity="0.35" stroke={colour.main} strokeWidth={pad * 0.05} />
      })}
    </>
  )
}

export default function LeagueAvatar({
  tier,
  weeklyXp,
  seed,
  name,
  size = 64,
  url,
  preset,
  className,
  avatarClassName,
}: Props) {
  const league = normalizeLeagueTier(tier)
  const earned = typeof weeklyXp === 'number' && Number.isFinite(weeklyXp) && weeklyXp > 0

  if (!league || !earned) {
    return <Avatar seed={seed} name={name} size={size} url={url} preset={preset} className={`${avatarClassName ?? ''} ${className ?? ''}`.trim()} />
  }

  const { outer, imageRadius, frameInnerRadius, collapsed } = leagueFrameGeometry(size)
  return (
    <span
      className={`relative inline-flex shrink-0 items-center justify-center ${className ?? ''}`}
      style={{ width: outer, height: outer }}
      data-league-frame={league}
      data-collapsed={collapsed ? 'true' : 'false'}
      data-image-radius={imageRadius}
      data-frame-inner-radius={frameInnerRadius}
    >
      <Avatar seed={seed} name={name} size={size} url={url} preset={preset} className={avatarClassName} />
      <svg
        aria-hidden="true"
        className="pointer-events-none absolute inset-0 h-full w-full motion-reduce:transition-none"
        viewBox="0 0 100 100"
        width={outer}
        height={outer}
      >
        <Frame tier={league} size={size} />
      </svg>
      <span className="sr-only">{league.charAt(0).toUpperCase() + league.slice(1)} league</span>
    </span>
  )
}
