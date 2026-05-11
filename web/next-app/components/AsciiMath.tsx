'use client'

import { useEffect, useState, useRef } from 'react'

interface AsciiMathProps {
  lines: string[]
  className?: string
}

export default function AsciiMath({ lines, className = '' }: AsciiMathProps) {
  const [visibleLines, setVisibleLines] = useState(0)
  const ref = useRef<HTMLPreElement>(null)

  useEffect(() => {
    const el = ref.current
    if (!el) return

    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          let i = 0
          const interval = setInterval(() => {
            i++
            setVisibleLines(i)
            if (i >= lines.length) clearInterval(interval)
          }, 120)
          observer.unobserve(el)
        }
      },
      { threshold: 0.3 }
    )
    observer.observe(el)
    return () => observer.disconnect()
  }, [lines.length])

  return (
    <pre
      ref={ref}
      className={`ascii-text text-center ${className}`}
      style={{ fontSize: '12px' }}
      aria-hidden="true"
    >
      {lines.slice(0, visibleLines).join('\n')}
    </pre>
  )
}
