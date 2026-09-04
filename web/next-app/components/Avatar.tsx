'use client'

interface AvatarProps {
  seed: string
  name: string
  size?: number
  url?: string
  preset?: number | null
  className?: string
}

// Discord-like palettes: flat bg + initial
const PALETTES: { bg: string; fg: string }[] = [
  { bg: '#5865F2', fg: '#ffffff' }, // blurple
  { bg: '#57F287', fg: '#23272A' }, // green
  { bg: '#FEE75C', fg: '#23272A' }, // yellow
  { bg: '#EB459E', fg: '#ffffff' }, // pink
  { bg: '#ED4245', fg: '#ffffff' }, // red
  { bg: '#2CC9A7', fg: '#23272A' }, // teal
  { bg: '#FAA61A', fg: '#23272A' }, // orange
  { bg: '#80848E', fg: '#ffffff' }, // gray
]

function hashSeed(s: string): number {
  let h = 0
  for (let i = 0; i < s.length; i++) {
    h = (h * 31 + s.charCodeAt(i)) >>> 0
  }
  return h
}

export function avatarPalette(seed: string): { bg: string; fg: string } {
  if (!seed) return PALETTES[0]
  return PALETTES[hashSeed(seed) % PALETTES.length]
}

export function initialFor(name: string): string {
  const trimmed = name.trim()
  if (!trimmed) return '?'
  // use first grapheme, uppercased
  return Array.from(trimmed)[0]?.toUpperCase() ?? '?'
}

export default function Avatar({ seed, name, size = 64, url, preset, className }: AvatarProps) {
  const palette = preset != null ? PALETTES[preset % PALETTES.length] : avatarPalette(seed)
  const initial = initialFor(name)
  const fontSize = Math.round(size * 0.42)

  if (url) {
    return (
      <img
        src={url}
        alt={`${name} avatar`}
        width={size}
        height={size}
        className={`rounded-full object-cover shrink-0 ${className ?? ''}`}
        style={{ width: size, height: size }}
        loading="lazy"
        referrerPolicy="no-referrer"
      />
    )
  }

  return (
    <div
      aria-label={`${name} avatar`}
      className={`rounded-full flex items-center justify-center shrink-0 select-none ${className ?? ''}`}
      style={{
        width: size,
        height: size,
        background: palette.bg,
        color: palette.fg,
        fontFamily: "'IBM Plex Mono', monospace",
        fontSize,
        fontWeight: 600,
        lineHeight: 1,
      }}
    >
      {initial}
    </div>
  )
}
