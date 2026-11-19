'use client'

import dynamic from 'next/dynamic'
import { useEffect, useRef, useState } from 'react'
import { crossSectionArea, lateralArea, radius, volume } from '../../prototypes/solid'

/**
 * Prototype surface for ADR-059. Untracked, dev-only, not linked from anywhere.
 * Each renderer is a separate dynamic import so the build attributes its own
 * chunk and an ordinary lesson never downloads any of them.
 */
const SolidThree = dynamic(() => import('../../prototypes/SolidThree'), { ssr: false })
const SolidCanvas = dynamic(() => import('../../prototypes/SolidCanvas'), { ssr: false })
const SolidSvg = dynamic(() => import('../../prototypes/SolidSvg'), { ssr: false })

type Mode = 'three' | 'canvas' | 'svg'

function webglAvailable(): boolean {
  if (typeof document === 'undefined') return false
  try {
    const canvas = document.createElement('canvas')
    return Boolean(canvas.getContext('webgl2') || canvas.getContext('webgl'))
  } catch {
    return false
  }
}

export default function Proto3d() {
  const [mode, setMode] = useState<Mode>('three')
  const [timing, setTiming] = useState<number | null>(null)
  const started = useRef(0)
  const [nonce, setNonce] = useState(0)
  const [webgl, setWebgl] = useState(true)
  useEffect(() => setWebgl(webglAvailable()), [])
  const [fallback, setFallback] = useState(false)
  const effective: Mode = mode === 'three' && (!webgl || fallback) ? 'svg' : mode

  // One clock for all three: the parent effect runs after the child's, so this
  // measures click to mounted-and-first-drawn without each prototype owning a
  // different definition of "ready".
  useEffect(() => {
    if (started.current === 0) return
    const id = requestAnimationFrame(() => {
      setTiming(performance.now() - started.current)
      started.current = 0
    })
    return () => cancelAnimationFrame(id)
  }, [effective, nonce])

  return (
    <main className="mx-auto max-w-3xl px-4 py-8">
      {/* Every route needs a way off it (test/navCoverage.test.ts), including
          a prototype that is not linked from anywhere. */}
      <a href="/learn" className="font-mono text-[11px] text-mathua-secondary hover:text-mathua-blue">
        ← Learn
      </a>
      <p className="mt-4 font-mono text-[11px] uppercase tracking-wider text-mathua-muted">Prototype, local only</p>
      <h1 className="mt-1 text-xl text-mathua-primary">3D solid of revolution</h1>
      <p className="mt-2 font-mono text-[11px] leading-relaxed text-mathua-secondary">
        r(x) = 1 + 0.35 sin(1.6x) on [0, 4]. Volume {volume().toFixed(3)}, lateral area{' '}
        {lateralArea().toFixed(3)}, cross-section at x = 2 is {crossSectionArea(2).toFixed(3)} (radius{' '}
        {radius(2).toFixed(3)}).
      </p>

      <div className="mt-4 flex flex-wrap gap-2">
        {(['three', 'canvas', 'svg'] as Mode[]).map((option) => (
          <button
            key={option}
            type="button"
            onClick={() => {
              setTiming(null)
              started.current = performance.now()
              setNonce((n) => n + 1)
              setMode(option)
            }}
            className={`border px-3 py-1 font-mono text-[11px] ${
              mode === option
                ? 'border-mathua-blue text-mathua-blue'
                : 'border-mathua-border text-mathua-secondary hover:border-mathua-blue'
            }`}
          >
            {option === 'three' ? 'Three.js' : option === 'canvas' ? 'Canvas 2D' : 'SVG'}
          </button>
        ))}
        <button
          type="button"
          onClick={() => {
            started.current = performance.now()
            setFallback((value) => !value)
          }}
          className="border border-mathua-border px-3 py-1 font-mono text-[11px] text-mathua-secondary hover:border-mathua-blue"
        >
          {fallback ? 'fallback on' : 'force fallback'}
        </button>
      </div>

      <p className="mt-2 font-mono text-[10px] text-mathua-muted">
        webgl {String(webgl)} · rendering {effective} · init {timing === null ? '…' : `${timing.toFixed(1)} ms`}
      </p>

      <div className="mt-4">
        {effective === 'three' && <SolidThree onReady={setTiming} />}
        {effective === 'canvas' && <SolidCanvas onReady={setTiming} />}
        {effective === 'svg' && <SolidSvg onReady={setTiming} />}
      </div>

      <p className="mt-4 font-mono text-[10px] text-mathua-muted">
        Drag or use the arrow keys to turn the solid. The slider moves the cutting plane.
      </p>
    </main>
  )
}
