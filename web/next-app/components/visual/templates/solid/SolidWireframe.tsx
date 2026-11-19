'use client'

import { useRef } from 'react'
import { radius, profile, X_MIN, X_MAX, CENTER_X } from '../../math/solid'
import { linePath, project, ringPath, type Camera } from './project'
import type { SolidViewProps } from './props'

const SIZE = 320
const SEGMENTS = 28

/**
 * The non-WebGL fallback: a wireframe drawn with the same geometry and the same
 * camera as the Three.js scene. It has no fill and no lighting, so it shows the
 * shape rather than the surface, which is the honest thing to show when the
 * browser cannot render the surface at all.
 */
export default function SolidWireframe({ theta, phi, cut, showCut, onOrbit }: SolidViewProps) {
  const dragging = useRef(false)
  const camera: Camera = { theta, phi, distance: 9, scale: SIZE * 0.3, cx: SIZE / 2, cy: SIZE / 2 }

  const stations = profile(18)
  const rings = stations.map(([x, r]) => {
    const points = Array.from({ length: SEGMENTS }, (_, j) => {
      const angle = (2 * Math.PI * j) / SEGMENTS
      return project({ x: x - CENTER_X, y: r * Math.cos(angle), z: r * Math.sin(angle) }, camera)
    })
    return ringPath(points)
  })
  const meridians = Array.from({ length: 12 }, (_, m) => {
    const angle = (2 * Math.PI * m) / 12
    return linePath(
      stations.map(([x, r]) =>
        project({ x: x - CENTER_X, y: r * Math.cos(angle), z: r * Math.sin(angle) }, camera),
      ),
    )
  })

  const cutRing = showCut
    ? (() => {
        const r = radius(cut)
        return ringPath(
          Array.from({ length: SEGMENTS }, (_, j) => {
            const angle = (2 * Math.PI * j) / SEGMENTS
            return project({ x: cut - CENTER_X, y: r * Math.cos(angle), z: r * Math.sin(angle) }, camera)
          }),
        )
      })()
    : ''

  return (
    <svg
      viewBox={`0 0 ${SIZE} ${SIZE}`}
      className="mx-auto block h-auto w-full max-w-[420px] cursor-grab font-mono active:cursor-grabbing"
      role="img"
      aria-label="A solid of revolution drawn as a wireframe"
      onPointerDown={(event) => {
        dragging.current = true
        event.currentTarget.setPointerCapture(event.pointerId)
      }}
      onPointerMove={(event) => {
        if (!dragging.current) return
        onOrbit(event.movementX * 0.01, event.movementY * 0.01)
      }}
      onPointerUp={() => {
        dragging.current = false
      }}
      onPointerCancel={() => {
        dragging.current = false
      }}
    >
      <g fill="none" stroke="var(--accent-blue)" strokeWidth="0.7" opacity="0.75">
        {rings.map((d, i) => (
          <path key={`r${i}`} d={d} />
        ))}
        {meridians.map((d, i) => (
          <path key={`m${i}`} d={d} />
        ))}
      </g>
      <path d={linePath([project({ x: X_MIN - CENTER_X, y: 0, z: 0 }, camera), project({ x: X_MAX - CENTER_X, y: 0, z: 0 }, camera)])} stroke="var(--text-muted)" strokeWidth="0.8" strokeDasharray="4 4" />
      {cutRing && <path d={cutRing} fill="var(--accent-teal)" fillOpacity="0.16" stroke="var(--accent-teal)" strokeWidth="1.6" />}
      <text x={SIZE - 8} y={14} textAnchor="end" fontSize="9" fill="var(--text-muted)">
        wireframe fallback
      </text>
    </svg>
  )
}
