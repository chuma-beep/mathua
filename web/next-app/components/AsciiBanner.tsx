'use client'

import { useEffect, useState } from 'react'

interface AsciiBannerProps {
  text: string
  className?: string
  delay?: number
}

export default function AsciiBanner({ text, className = '', delay = 200 }: AsciiBannerProps) {
  const [revealed, setRevealed] = useState(false)

  useEffect(() => {
    const t = setTimeout(() => setRevealed(true), delay)
    return () => clearTimeout(t)
  }, [delay])

  return (
    <pre
      className={`ascii-text select-none ${className}`}
      style={{
        clipPath: revealed ? 'inset(0 0 0 0)' : 'inset(0 100% 0 0)',
        transition: 'clip-path 0.8s steps(40)',
        fontSize: '12px',
        lineHeight: 1.3,
      }}
      aria-hidden="true"
    >
      {text}
    </pre>
  )
}
