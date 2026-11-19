'use client'

import { useEffect, useRef, useState } from 'react'
import { project, quads, shade, type Camera } from './project'
import { radius, slope, X_MAX, X_MIN } from './solid'

/**
 * Option C: Canvas 2D with a painter's sort.
 *
 * No engine, no WebGL, no lights: the quads are projected here and filled from
 * the back forwards, so a nearer quad covers a farther one. That is the whole
 * depth model, and it is enough for a convex surface of revolution seen from
 * outside, which is the only thing this figure draws.
 */
export default function SolidCanvas({ onReady }: { onReady?: (ms: number) => void }) {
  const canvas = useRef<HTMLCanvasElement | null>(null)
  const [theta, setTheta] = useState(0.7)
  const [phi, setPhi] = useState(1.15)
  const [cut, setCut] = useState(true)
  const [cutX, setCutX] = useState(2)

  const turn = (dTheta: number, dPhi: number) => {
    setTheta((t) => t + dTheta)
    setPhi((p) => Math.min(Math.PI - 0.15, Math.max(0.15, p + dPhi)))
  }

  const reported = useRef(false)
  useEffect(() => {
    const element = canvas.current
    if (!element) return
    const started = performance.now()
    const context = element.getContext('2d')
    if (!context) return

    const size = element.clientWidth || 320
    const dpr = Math.min(window.devicePixelRatio || 1, 2)
    element.width = Math.round(size * dpr)
    element.height = Math.round(size * dpr)
    context.setTransform(dpr, 0, 0, dpr, 0, 0)

    const camera: Camera = {
      theta,
      phi,
      distance: 9,
      // 0.26 rather than 0.34: perspective magnifies the near side by about 1.3,
      // and at 0.34 the silhouette overflowed the canvas on that side.
      scale: size * 0.26,
      cx: size / 2,
      cy: size / 2,
    }

    const faces = quads(28, 48, radius, slope, X_MIN, X_MAX)
      .map((quad) => {
        const projected = quad.points.map((point) => project(point, camera))
        const depth = projected.reduce((sum, p) => sum + p.depth, 0) / projected.length
        return { projected, depth, brightness: shade(quad.normal) }
      })
      .sort((a, b) => b.depth - a.depth)

    // The two end caps, as polygons in the same painter's order.
    for (const [x, normalX] of [[X_MIN, -1] as const, [X_MAX, 1] as const]) {
      const r = radius(x)
      const projected = Array.from({ length: 48 }, (_, i) => {
        const a = (2 * Math.PI * i) / 48
        return project({ x: x - (X_MIN + X_MAX) / 2, y: r * Math.cos(a), z: r * Math.sin(a) }, camera)
      })
      faces.push({
        projected,
        depth: projected.reduce((sum, p) => sum + p.depth, 0) / projected.length,
        brightness: shade({ x: normalX, y: 0, z: 0 }),
      })
    }
    faces.sort((a, b) => b.depth - a.depth)

    context.clearRect(0, 0, size, size)
    for (const face of faces) {
      const tone = Math.round(30 + 170 * face.brightness)
      context.fillStyle = `rgb(${Math.round(tone * 0.45)}, ${Math.round(tone * 0.72)}, ${tone})`
      context.strokeStyle = context.fillStyle
      context.beginPath()
      face.projected.forEach((p, i) => (i === 0 ? context.moveTo(p.x, p.y) : context.lineTo(p.x, p.y)))
      context.closePath()
      context.fill()
      context.stroke()
    }

    if (cut) {
      const r = radius(cutX)
      context.strokeStyle = '#2f8f6f'
      context.lineWidth = 1.6
      context.beginPath()
      for (let i = 0; i <= 64; i++) {
        const a = (2 * Math.PI * i) / 64
        const p = project({ x: cutX - (X_MIN + X_MAX) / 2, y: r * Math.cos(a), z: r * Math.sin(a) }, camera)
        if (i === 0) context.moveTo(p.x, p.y)
        else context.lineTo(p.x, p.y)
      }
      context.stroke()
    }

    if (!reported.current) {
      reported.current = true
      onReady?.(performance.now() - started)
    }
  }, [theta, phi, cut, cutX, onReady])

  return (
    <div>
      <canvas
        ref={canvas}
        data-proto="canvas"
        role="img"
        aria-label="A solid of revolution with a circular cross-section"
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
        className="mx-auto block w-full max-w-[420px] cursor-grab touch-none focus-visible:outline focus-visible:outline-2 focus-visible:outline-mathua-blue active:cursor-grabbing"
      />
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
