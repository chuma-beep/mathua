'use client'

import { useCallback, useMemo, useRef, useState, useId } from 'react'
import FigureFrame from '../FigureFrame'
import { NumberField, Toggle } from '../Controls'
import { complexMultiplySchema } from '../schemas'
import { carg, cmul, cmod, type Complex } from '../math/complex'
import { arcPath, clamp, scaleLinear } from '../scale'
import { fmt } from '../format'
import { useDrag } from '../useDrag'
import type { TemplateProps } from '../types'

const SIZE = 320
const PAD = 30
const LIMIT = 3

export default function ComplexMultiply({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => {
    const p = complexMultiplySchema.parse(params)
    return {
      z: [Number(p.z[0]), Number(p.z[1])] as Complex,
      w: [Number(p.w[0]), Number(p.w[1])] as Complex,
      showUnitCircle: p.showUnitCircle,
    }
  }, [params])
  const [z, setZ] = useState<Complex>(defaults.z)
  const [w, setW] = useState<Complex>(defaults.w)
  const [showUnitCircle, setShowUnitCircle] = useState(defaults.showUnitCircle)
  const svgRef = useRef<SVGSVGElement | null>(null)
  const uid = useId().replace(/:/g, '')

  const product = cmul(z, w)
  // Fit the view to the four points rather than to a radius: a fixed origin-centred square
  // spends most of its area on quadrants the figure never draws in, which is why the vectors
  // used to sit in a small cluster at the bottom of an empty canvas. One span for both axes
  // keeps the unit circle a circle.
  const lo = Math.min(0, z[0], w[0], product[0], z[1], w[1], product[1]) - 0.7
  const hi = Math.max(0, z[0], w[0], product[0], z[1], w[1], product[1]) + 0.7
  const unit = (SIZE - PAD * 2) / Math.max(2, hi - lo)
  const X = scaleLinear(lo, hi, PAD, SIZE - PAD)
  const Y = scaleLinear(hi, lo, PAD, SIZE - PAD)
  const cx = X(0)
  const cy = Y(0)

  const onDrag = useCallback(
    (point: { x: number; y: number }) => {
      const wx = clamp((point.x - cx) / unit, -LIMIT, LIMIT)
      const wy = clamp((cy - point.y) / unit, -LIMIT, LIMIT)
      const dz = Math.hypot(wx - z[0], wy - z[1])
      const dw = Math.hypot(wx - w[0], wy - w[1])
      if (dz <= dw) setZ([wx, wy])
      else setW([wx, wy])
    },
    [cx, cy, unit, z, w],
  )
  useDrag(svgRef, onDrag)

  const reset = () => {
    setZ(defaults.z)
    setW(defaults.w)
    setShowUnitCircle(defaults.showUnitCircle)
  }

  const radius = 42
  const arrow = (from: Complex, to: Complex, color: string, marker: string) => (
    <line
      x1={X(from[0])}
      y1={Y(from[1])}
      x2={X(to[0])}
      y2={Y(to[1])}
      stroke={color}
      strokeWidth="1.8"
      markerEnd={`url(#${uid}-${marker})`}
    />
  )

  return (
    <FigureFrame
      title={title ?? 'Complex multiplication'}
      caption={
        caption ??
        'Drag either point, or edit its parts. Multiplying adds the arguments and multiplies the moduli, so the product is a rotation and a scaling of the plane.'
      }
      onReset={reset}
      hint="Drag the nearer of z and w, or type its parts."
      readout={[
        { label: '|z|', value: fmt(cmod(z)) },
        { label: 'arg z', value: fmt(carg(z)) },
        { label: '|w|', value: fmt(cmod(w)) },
        { label: 'arg w', value: fmt(carg(w)) },
        { label: '|zw|', value: fmt(cmod(product)) },
        { label: 'arg zw', value: fmt(carg(product)) },
        { label: 'arg sum', value: fmt(carg(z) + carg(w)) },
      ]}
    >
      <svg
        ref={svgRef}
        viewBox={`0 0 ${SIZE} ${SIZE}`}
        className="mx-auto block h-auto w-full max-w-[420px] cursor-grab touch-none font-mono active:cursor-grabbing"
        style={{ touchAction: 'none' }}
        role="img"
        aria-label="Two complex numbers, their arguments and their product"
      >
        <defs>
          {[
            ['blue', 'var(--accent-blue)'],
            ['green', 'var(--accent-green)'],
            ['teal', 'var(--accent-teal)'],
          ].map(([name, color]) => (
            <marker key={name} id={`${uid}-${name}`} markerWidth="6" markerHeight="6" refX="5" refY="3" orient="auto">
              <path d="M0,0 L6,3 L0,6 Z" fill={color} />
            </marker>
          ))}
        </defs>

        <line x1={PAD} y1={cy} x2={SIZE - PAD} y2={cy} stroke="var(--text-muted)" strokeWidth="0.8" />
        <line x1={cx} y1={PAD} x2={cx} y2={SIZE - PAD} stroke="var(--text-muted)" strokeWidth="0.8" />
        {showUnitCircle && <circle cx={cx} cy={cy} r={unit} fill="none" stroke="var(--border-strong)" strokeDasharray="3 4" />}

        <path d={arcPath(cx, cy, radius, carg(z))} fill="none" stroke="var(--accent-blue)" strokeWidth="1" opacity="0.6" />
        <path d={arcPath(cx, cy, radius + 10, carg(w))} fill="none" stroke="var(--accent-green)" strokeWidth="1" opacity="0.6" />
        <path d={arcPath(cx, cy, radius + 20, carg(product))} fill="none" stroke="var(--accent-teal)" strokeWidth="1" opacity="0.6" />

        {arrow([0, 0], z, 'var(--accent-blue)', 'blue')}
        {arrow([0, 0], w, 'var(--accent-green)', 'green')}
        {arrow([0, 0], product, 'var(--accent-teal)', 'teal')}

        <circle cx={X(z[0])} cy={Y(z[1])} r="4" fill="var(--accent-blue)" />
        <circle cx={X(w[0])} cy={Y(w[1])} r="4" fill="var(--accent-green)" />
        <circle cx={X(product[0])} cy={Y(product[1])} r="3.2" fill="var(--accent-teal)" />

        <text x={X(z[0]) + 7} y={Y(z[1]) - 7} fontSize="10" fill="var(--accent-blue)">z</text>
        <text x={X(w[0]) + 7} y={Y(w[1]) - 7} fontSize="10" fill="var(--accent-green)">w</text>
        <text x={X(product[0]) + 7} y={Y(product[1]) + 12} fontSize="10" fill="var(--accent-teal)">zw</text>
      </svg>

      <div className="mt-3 grid grid-cols-2 gap-2">
        <NumberField label="z re" value={z[0]} onChange={(v) => setZ([v, z[1]])} />
        <NumberField label="z im" value={z[1]} onChange={(v) => setZ([z[0], v])} />
        <NumberField label="w re" value={w[0]} onChange={(v) => setW([v, w[1]])} />
        <NumberField label="w im" value={w[1]} onChange={(v) => setW([w[0], v])} />
      </div>
      <div className="mt-2">
        <Toggle label="Show unit circle" checked={showUnitCircle} onChange={setShowUnitCircle} />
      </div>
    </FigureFrame>
  )
}
