'use client'

import { useCallback, useMemo, useRef, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { NumberField } from '../Controls'
import { XAxisTicks, YAxisTicks } from '../Axis'
import { basisChangeSchema } from '../schemas'
import { basisMatrix, combine, coordinatesOf, isNearlyDependent } from '../math/basis'
import { det, type Vec2 } from '../math/matrix'
import { clamp, pathFromPoints, scaleLinear } from '../scale'
import { fmt } from '../format'
import { useDrag } from '../useDrag'
import type { TemplateProps } from '../types'

const SIZE = 320
const PAD = 26
const EXTENT = 3.4

export default function BasisChange({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => {
    const p = basisChangeSchema.parse(params)
    return {
      u: [Number(p.u[0]), Number(p.u[1])] as Vec2,
      v: [Number(p.v[0]), Number(p.v[1])] as Vec2,
      point: [Number(p.point[0]), Number(p.point[1])] as Vec2,
    }
  }, [params])

  const [u, setU] = useState<Vec2>(defaults.u)
  const [v, setV] = useState<Vec2>(defaults.v)
  const [point, setPoint] = useState<Vec2>(defaults.point)
  const svgRef = useRef<SVGSVGElement | null>(null)

  const basis = basisMatrix(u, v)
  const d = det(basis)
  const coords = coordinatesOf(basis, point)
  const rebuilt = coords ? combine(basis, coords) : null
  const fragile = coords !== null && isNearlyDependent(basis)
  // fmt rounds to three decimals, so a determinant of 1e-4 would print as 0 beside
  // coordinates that exist. Scientific notation keeps the readout from contradicting
  // the figure: the basis is thin, not singular.
  const detText = Math.abs(d) < 1e-3 && d !== 0 ? d.toExponential(2) : fmt(d, 3)
  const coordsVisible =
    coords !== null && coords.every((value) => Number.isFinite(value) && Math.abs(value) <= EXTENT * 4)

  const X = scaleLinear(-EXTENT, EXTENT, PAD, SIZE - PAD)
  const Y = scaleLinear(EXTENT, -EXTENT, PAD, SIZE - PAD)
  const invertX = scaleLinear(PAD, SIZE - PAD, -EXTENT, EXTENT)
  const invertY = scaleLinear(SIZE - PAD, PAD, -EXTENT, EXTENT)
  const cx = X(0)
  const cy = Y(0)

  const onDrag = useCallback(
    (p: { x: number; y: number }) => {
      setPoint([clamp(invertX(p.x), -EXTENT, EXTENT), clamp(invertY(p.y), -EXTENT, EXTENT)])
    },
    [invertX, invertY],
  )
  useDrag(svgRef, onDrag)

  const reset = () => {
    setU(defaults.u)
    setV(defaults.v)
    setPoint(defaults.point)
  }

  // The lattice the basis draws: every line of constant coordinate, so the point's
  // coordinates can be read off by eye rather than only from the readout.
  const lattice = useMemo(() => {
    const lines: Array<[Vec2, Vec2]> = []
    for (let k = -4; k <= 4; k++) {
      lines.push([combine(basis, [k, -4]), combine(basis, [k, 4])])
      lines.push([combine(basis, [-4, k]), combine(basis, [4, k])])
    }
    return lines
  }, [basis])

  const arrow = (to: Vec2, color: string, id: string) => (
    <line x1={cx} y1={cy} x2={X(to[0])} y2={Y(to[1])} stroke={color} strokeWidth="1.8" markerEnd={`url(#${id})`} />
  )

  return (
    <FigureFrame
      title={title ?? 'The same point, different coordinates'}
      hint="Drag the point. Then change u and v and watch its coordinates change while it stays put."
      caption={
        caption ??
        'The grid is the basis. The coordinates (α, β) say how far to go along u and along v, so the pair depends on the arrows even though the point does not.'
      }
      onReset={reset}
      readout={[
        { label: 'det P', value: detText },
        { label: 'coords', value: coords ? `(α, β) = (${fmt(coords[0])}, ${fmt(coords[1])})` : 'undefined' },
        { label: 'αu + βv', value: rebuilt ? `(${fmt(rebuilt[0])}, ${fmt(rebuilt[1])})` : 'no basis' },
        { label: 'point', value: `(${fmt(point[0])}, ${fmt(point[1])})` },
      ]}
    >
      <svg
        ref={svgRef}
        viewBox={`0 0 ${SIZE} ${SIZE}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        style={{ touchAction: 'none' }}
        role="img"
        aria-label="A point and the grid drawn by two basis vectors"
      >
        <defs>
          {[
            ['u', 'var(--accent-teal)'],
            ['v', 'var(--accent-green)'],
            ['p', 'var(--accent-blue)'],
          ].map(([name, color]) => (
            <marker key={name} id={`basis-${name}`} markerWidth="6" markerHeight="6" refX="5" refY="3" orient="auto">
              <path d="M0,0 L6,3 L0,6 Z" fill={color} />
            </marker>
          ))}
          {/* The lattice runs to the extent of the basis, which is past the axes:
              clipped so the grid stops at the plot instead of the figure edge. */}
          <clipPath id="basis-plot">
            <rect x={PAD} y={PAD} width={SIZE - PAD * 2} height={SIZE - PAD * 2} />
          </clipPath>
        </defs>

        <g clipPath="url(#basis-plot)">
        {Math.abs(d) > 1e-9 &&
          lattice.map(([from, to], i) => (
            <line
              key={i}
              x1={X(from[0])}
              y1={Y(from[1])}
              x2={X(to[0])}
              y2={Y(to[1])}
              stroke="var(--border-strong)"
              strokeWidth="0.7"
              strokeDasharray="3 3"
              opacity="0.95"
            />
          ))}
        </g>

        <line x1={PAD} y1={cy} x2={SIZE - PAD} y2={cy} stroke="var(--text-muted)" strokeWidth="0.8" />
        <line x1={cx} y1={PAD} x2={cx} y2={SIZE - PAD} stroke="var(--text-muted)" strokeWidth="0.8" />
        <XAxisTicks ticks={[-3, -2, -1, 1, 2, 3]} scale={X} at={cy} />
        <YAxisTicks ticks={[-3, -2, -1, 1, 2, 3]} scale={Y} at={cx} />

        {arrow(u, 'var(--accent-teal)', 'basis-u')}
        {arrow(v, 'var(--accent-green)', 'basis-v')}

        {coordsVisible && (
          <path
            d={pathFromPoints([
              [X(0), Y(0)],
              [X(coords[0] * u[0]), Y(coords[0] * u[1])],
              [X(point[0]), Y(point[1])],
            ])}
            fill="none"
            stroke="var(--accent-blue)"
            strokeWidth="1.2"
            strokeDasharray="3 3"
            opacity="0.95"
          />
        )}

        <circle cx={X(point[0])} cy={Y(point[1])} r="4.5" fill="var(--accent-blue)" />

        <text x={X(u[0]) + 6} y={Y(u[1]) + 12} fontSize="10" fill="var(--accent-teal)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
          u
        </text>
        <text x={X(v[0]) + 6} y={Y(v[1]) - 6} fontSize="10" fill="var(--accent-green)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
          v
        </text>
        <text x={X(point[0]) + 7} y={Y(point[1]) + 14} fontSize="10" fill="var(--accent-blue)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
          p
        </text>

        {!coords && (
          <text x={PAD + 2} y={PAD + 12} fontSize="10" fill="var(--accent-teal)">
            u and v are dependent: no coordinates
          </text>
        )}
        {fragile && (
          <text x={PAD + 2} y={PAD + 12} fontSize="10" fill="var(--accent-teal)">
            nearly dependent: coordinates are huge
          </text>
        )}
      </svg>

      <div className="mt-3 grid grid-cols-2 gap-2">
        <NumberField label="ux" value={u[0]} min={-3} max={3} onChange={(value) => setU([value, u[1]])} />
        <NumberField label="uy" value={u[1]} min={-3} max={3} onChange={(value) => setU([u[0], value])} />
        <NumberField label="vx" value={v[0]} min={-3} max={3} onChange={(value) => setV([value, v[1]])} />
        <NumberField label="vy" value={v[1]} min={-3} max={3} onChange={(value) => setV([v[0], value])} />
      </div>
    </FigureFrame>
  )
}
