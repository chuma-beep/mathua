'use client'

import { useEffect, useState, useRef } from 'react'

interface Opts {
  hideThreshold?: number
  topOffset?: number
  idleMs?: number
  bottomOffset?: number
}

export function useScrollDirection({
  hideThreshold = 12,
  topOffset = 40,
  idleMs = 300,
  bottomOffset = 24,
}: Opts = {}) {
  const [hidden, setHidden] = useState(false)
  const lastY = useRef(0)
  const idleTimer = useRef<number | null>(null)

  useEffect(() => {
    let ticking = false
    const onScroll = () => {
      if (ticking) return
      ticking = true
      requestAnimationFrame(() => {
        const y = window.scrollY
        const delta = y - lastY.current
        const atBottom = window.innerHeight + y >= document.documentElement.scrollHeight - bottomOffset

        if (idleTimer.current !== null) window.clearTimeout(idleTimer.current)
        idleTimer.current = window.setTimeout(() => setHidden(false), idleMs)

        if (atBottom) setHidden(false)
        else if (y <= topOffset) setHidden(false)
        else if (delta > hideThreshold) setHidden(true)
        else if (delta < -hideThreshold) setHidden(false)

        lastY.current = y
        ticking = false
      })
    }
    window.addEventListener('scroll', onScroll, { passive: true })
    return () => {
      window.removeEventListener('scroll', onScroll)
      if (idleTimer.current !== null) window.clearTimeout(idleTimer.current)
    }
  }, [hideThreshold, topOffset, idleMs, bottomOffset])

  return hidden
}
