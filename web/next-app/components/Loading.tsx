'use client'

import { useEffect, useState } from 'react'

const FRAMES = ['⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏']

interface LoadingProps {
  label?: string
  full?: boolean
  size?: number
  inline?: boolean
}

export default function Loading({ label, full = false, size = 18, inline = false }: LoadingProps) {
  const [frame, setFrame] = useState(0)
  const [reduced, setReduced] = useState(true)

  useEffect(() => {
    const mq = window.matchMedia('(prefers-reduced-motion: reduce)')
    setReduced(mq.matches)
    if (mq.matches) return
    const id = setInterval(() => setFrame(f => (f + 1) % FRAMES.length), 80)
    return () => clearInterval(id)
  }, [])

  if (inline) {
    return (
      <span
        aria-hidden
        style={{
          display: 'inline-block',
          fontFamily: "'IBM Plex Mono', monospace",
          fontSize: size,
          lineHeight: 1,
          width: size,
          textAlign: 'center',
          verticalAlign: '-0.125em',
        }}
      >
        {FRAMES[reduced ? 0 : frame]}
      </span>
    )
  }

  return (
    <div
      role="status"
      aria-live="polite"
      style={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        gap: 10,
        width: '100%',
        minHeight: full ? 'clamp(320px, 50vh, 520px)' : undefined,
        padding: full ? undefined : '1rem 0',
      }}
    >
      <span
        aria-hidden
        style={{
          fontFamily: "'IBM Plex Mono', monospace",
          fontSize: size,
          lineHeight: 1,
          color: 'var(--accent-blue)',
          display: 'inline-block',
          width: size,
          textAlign: 'center',
        }}
      >
        {FRAMES[reduced ? 0 : frame]}
      </span>
      {label && (
        <span
          style={{
            fontFamily: "'IBM Plex Mono', monospace",
            fontSize: 12,
            color: 'var(--text-muted)',
            letterSpacing: '0.04em',
          }}
        >
          {label}
        </span>
      )}
    </div>
  )
}
