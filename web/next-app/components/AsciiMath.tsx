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
    let intervalId: ReturnType<typeof setInterval> | null = null

    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          let i = 0
          intervalId = setInterval(() => {
            i++
            setVisibleLines(i)
            if (i >= lines.length && intervalId) {
              clearInterval(intervalId)
              intervalId = null
            }
          }, 120)
          observer.unobserve(el)
        }
      },
      { threshold: 0.3 }
    )
    observer.observe(el)
    return () => {
      observer.disconnect()
      if (intervalId) clearInterval(intervalId)
    }
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
