'use client'

import { useEffect, useMemo, useRef, useState } from 'react'
import { project, type Camera } from './project'
import { profile, radius, X_MAX, X_MIN } from './solid'

/**
 * Option B: hand-rolled SVG projection.
 *
 * SVG cannot shade three thousand quads without putting three thousand nodes in
 * the document, so this draws the surface as a wireframe: rings around the axis
 * and meridians along it. It is a legitimate fallback, and it is not a shaded
 * solid, which is the finding this prototype exists to record.
 */
export default function SolidSvg({ onReady }: { onReady?: (ms: number) => void }) {
  const [theta, setTheta] = useState(0.7)
  const [phi, setPhi] = useState(1.15)
  const [cut, setCut] = useState(true)
  const [cutX, setCutX] = useState(2)
  const size = 320

  const camera: Camera = { theta, phi, distance: 9, // 0.26 rather than 0.34: perspective magnifies the near side by about 1.3,
      // and at 0.34 the silhouette overflowed the canvas on that side.
      scale: size * 0.26, cx: size / 2, cy: size / 2 }
  const rings = 18
  const meridians = 12

  const ringPaths = useMemo(() => {
    const out: Array<{ front: string; back: string }> = []
    for (let i = 0; i <= rings; i++) {
      const x = X_MIN + ((X_MAX - X_MIN) * i) / rings
      const r = radius(x)
      let front = ''
      let back = ''
      // A ring crosses from the near half to the far half twice, so each half is
      // its own subpath and each subpath needs an M. Joining them with L drew a
      // chord across the solid, and the browser rejected the d attribute.
      let lastSide: boolean | null = null
      for (let j = 0; j <= 72; j++) {
        const a = (2 * Math.PI * j) / 72
        const p = project({ x: x - (X_MIN + X_MAX) / 2, y: r * Math.cos(a), z: r * Math.sin(a) }, camera)
        const near = p.depth < camera.distance
        const command = `${near === lastSide ? 'L' : 'M'} ${p.x.toFixed(1)} ${p.y.toFixed(1)}`
        if (near) front += command
        else back += command
        lastSide = near
      }
      out.push({ front, back })
    }
    return out
  }, [camera.theta, camera.phi])

  const meridianPaths = useMemo(() => {
    const out: string[] = []
    for (let m = 0; m < meridians; m++) {
      const a = (2 * Math.PI * m) / meridians
      let path = ''
      for (const [x, r] of profile(48)) {
        const p = project({ x, y: r * Math.cos(a), z: r * Math.sin(a) }, camera)
        path += `${path ? 'L' : 'M'} ${p.x.toFixed(1)} ${p.y.toFixed(1)}`
      }
      out.push(path)
    }
    return out
  }, [camera.theta, camera.phi])

  const section = useMemo(() => {
    const r = radius(cutX)
    let path = ''
    for (let i = 0; i <= 64; i++) {
      const a = (2 * Math.PI * i) / 64
      const p = project({ x: cutX - (X_MIN + X_MAX) / 2, y: r * Math.cos(a), z: r * Math.sin(a) }, camera)
      path += `${path ? 'L' : 'M'} ${p.x.toFixed(1)} ${p.y.toFixed(1)}`
    }
    return path
  }, [cutX, camera.theta, camera.phi])

  const started = useRef(performance.now())
  useEffect(() => {
    onReady?.(performance.now() - started.current)
    // Report once, on mount; a per-turn report would be a render side effect.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const turn = (dTheta: number, dPhi: number) => {
    setTheta((t) => t + dTheta)
    setPhi((p) => Math.min(Math.PI - 0.15, Math.max(0.15, p + dPhi)))
  }

  return (
    <div>
      <svg
        viewBox={`0 0 ${size} ${size}`}
        data-proto="svg"
        className="mx-auto block h-auto w-full max-w-[420px] cursor-grab touch-none focus-visible:outline focus-visible:outline-2 focus-visible:outline-mathua-blue"
        role="img"
        aria-label="A wireframe solid of revolution with a circular cross-section"
        tabIndex={0}
        onKeyDown={(event) => {
          if (event.key === 'ArrowLeft') turn(-0.12, 0)
          else if (event.key === 'ArrowRight') turn(0.12, 0)
          else if (event.key === 'ArrowUp') turn(0, -0.1)
          else if (event.key === 'ArrowDown') turn(0, 0.1)
          else return
          event.preventDefault()
        }}
        onPointerDown={(event) => {
          let lastX = event.clientX
          let lastY = event.clientY
          const move = (moveEvent: PointerEvent) => {
            turn((moveEvent.clientX - lastX) * 0.01, -(moveEvent.clientY - lastY) * 0.01)
            lastX = moveEvent.clientX
            lastY = moveEvent.clientY
          }
          const up = () => {
            window.removeEventListener('pointermove', move)
            window.removeEventListener('pointerup', up)
          }
          window.addEventListener('pointermove', move)
          window.addEventListener('pointerup', up)
        }}
      >
        {ringPaths.map((ring, i) => (
          <g key={i}>
            <path d={ring.back} fill="none" stroke="var(--accent-blue)" strokeWidth="0.6" opacity="0.18" />
            <path d={ring.front} fill="none" stroke="var(--accent-blue)" strokeWidth="0.9" opacity="0.55" />
          </g>
        ))}
        {meridianPaths.map((path, i) => (
          <path key={`m${i}`} d={path} fill="none" stroke="var(--accent-blue)" strokeWidth="1" opacity="0.75" />
        ))}
        {cut && <path d={section} fill="none" stroke="var(--accent-teal)" strokeWidth="1.8" />}
      </svg>
      <div className="mt-2 flex flex-wrap items-center gap-2 font-mono text-[10px] text-mathua-secondary">
        <label className="flex items-center gap-1">
          <input type="checkbox" checked={cut} onChange={(e) => setCut(e.target.checked)} />
          cross-section
        </label>
        <input
          type="range"
          aria-label="cross-section position"
          min={X_MIN}
          max={X_MAX}
          step={0.02}
          value={cutX}
          onChange={(e) => setCutX(Number(e.target.value))}
          className="min-w-[110px] flex-1"
        />
        <output className="tabular-nums">{cutX.toFixed(2)}</output>
      </div>
    </div>
  )
}
