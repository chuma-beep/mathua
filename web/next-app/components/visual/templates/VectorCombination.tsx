'use client'

import { useCallback, useMemo, useRef, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { NumberField, Slider, Toggle } from '../Controls'
import { XAxisTicks, YAxisTicks } from '../Axis'
import { vectorCombinationSchema } from '../schemas'
import { addVec, angleBetween, length, scaleVec, subVec, type Vec2 } from '../math/matrix'
import { clamp, pathFromPoints, scaleLinear } from '../scale'
import { fmt } from '../format'
import { useDrag } from '../useDrag'
import type { TemplateProps } from '../types'

const SIZE = 320
const PAD = 28
const LIMIT = 4

export default function VectorCombination({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => {
    const p = vectorCombinationSchema.parse(params)
    return {
      a: [Number(p.a[0]), Number(p.a[1])] as Vec2,
      b: [Number(p.b[0]), Number(p.b[1])] as Vec2,
      k: p.k,
      showSum: p.showSum,
    }
  }, [params])

  const [a, setA] = useState<Vec2>(defaults.a)
  const [b, setB] = useState<Vec2>(defaults.b)
  const [k, setK] = useState(defaults.k)
  const [showSum, setShowSum] = useState(defaults.showSum)
  const svgRef = useRef<SVGSVGElement | null>(null)

  const sum = addVec(a, b)
  const difference = subVec(a, b)
  const scaled = scaleVec(b, k)
  // A zero vector has no direction, so there is no angle to report. The helper
  // returns 0 as a guard, which would otherwise read as "they are parallel".
  const angleText =
    length(a) < 1e-9 || length(b) < 1e-9 ? 'undefined' : `${fmt((angleBetween(a, b) * 180) / Math.PI, 1)}°`
  const scaledVisible = length(scaled) > 1e-9

  const reach = Math.max(
    ...a.map(Math.abs),
    ...b.map(Math.abs),
    ...sum.map(Math.abs),
    ...scaled.map(Math.abs),
  )
  // At least 3.2 so the whole ±3 tick row sits inside the plot, plus headroom for
  // a label at an arrowhead rather than past the edge of the frame.
  const extent = Math.max(3.2, reach * 1.12)
  const X = scaleLinear(-extent, extent, PAD, SIZE - PAD)
  const Y = scaleLinear(extent, -extent, PAD, SIZE - PAD)
  const cx = X(0)
  const cy = Y(0)
  // A label beside an arrowhead flips inward near the right edge, because the
  // alternative is a clipped label, which is what "a + b" was.
  const anchorAt = (v: number) => (v > extent * 0.55 ? 'end' : 'start') as 'start' | 'end'
  const dxAt = (v: number) => (v > extent * 0.55 ? -8 : 8)

  const onDrag = useCallback(
    (point: { x: number; y: number }) => {
      const unit = (SIZE - PAD * 2) / (extent * 2)
      const px = clamp((point.x - cx) / unit, -LIMIT, LIMIT)
      const py = clamp((cy - point.y) / unit, -LIMIT, LIMIT)
      const da = Math.hypot(px - a[0], py - a[1])
      const db = Math.hypot(px - b[0], py - b[1])
      if (da <= db) setA([px, py])
      else setB([px, py])
    },
    [cx, cy, extent, a, b],
  )
  useDrag(svgRef, onDrag)

  const reset = () => {
    setA(defaults.a)
    setB(defaults.b)
    setK(defaults.k)
    setShowSum(defaults.showSum)
  }

  const arrow = (to: Vec2, color: string, marker: string) => (
    <line x1={cx} y1={cy} x2={X(to[0])} y2={Y(to[1])} stroke={color} strokeWidth="1.8" markerEnd={`url(#vec-${marker})`} />
  )

  return (
    <FigureFrame
      title={title ?? 'Adding vectors, and scaling one'}
      hint="Drag either arrowhead, or type the components. k scales b."
      caption={
        caption ??
        'The sum is the diagonal of the parallelogram spanned by a and b, and a − b is the other diagonal. Multiplying by k scales the length by |k| and reverses the direction when k is negative.'
      }
      onReset={reset}
      readout={[
        { label: 'a', value: `${fmt(a[0])}, ${fmt(a[1])}` },
        { label: 'b', value: `${fmt(b[0])}, ${fmt(b[1])}` },
        { label: 'a + b', value: `${fmt(sum[0])}, ${fmt(sum[1])}` },
        { label: 'a − b', value: `${fmt(difference[0])}, ${fmt(difference[1])}` },
        { label: 'k b', value: `${fmt(scaled[0])}, ${fmt(scaled[1])}` },
        { label: '|a|', value: fmt(length(a)) },
        // A zero vector has no direction, so there is no angle to report. The
        // helper returns 0 as a guard, which would read as "they are parallel".
        { label: 'angle a,b', value: angleText },
      ]}
    >
      <svg
        ref={svgRef}
        viewBox={`0 0 ${SIZE} ${SIZE}`}
        className="mx-auto block h-auto w-full max-w-[420px] cursor-grab font-mono active:cursor-grabbing"
        style={{ touchAction: 'none' }}
        role="img"
        aria-label="Two vectors, their sum and a scaled copy"
      >
        <defs>
          {[
            ['blue', 'var(--accent-blue)'],
            ['teal', 'var(--accent-teal)'],
            ['green', 'var(--accent-green)'],
          ].map(([name, color]) => (
            <marker key={name} id={`vec-${name}`} markerWidth="6" markerHeight="6" refX="5" refY="3" orient="auto">
              <path d="M0,0 L6,3 L0,6 Z" fill={color} />
            </marker>
          ))}
        </defs>

        <line x1={PAD} y1={cy} x2={SIZE - PAD} y2={cy} stroke="var(--text-muted)" strokeWidth="0.8" />
        <line x1={cx} y1={PAD} x2={cx} y2={SIZE - PAD} stroke="var(--text-muted)" strokeWidth="0.8" />

        {showSum && (
          <>
            <path
              d={pathFromPoints([
                [X(a[0]), Y(a[1])],
                [X(sum[0]), Y(sum[1])],
                [X(b[0]), Y(b[1])],
              ])}
              fill="none"
              stroke="var(--accent-blue)"
              strokeWidth="0.9"
              strokeDasharray="3 3"
              opacity="0.55"
            />
            {arrow(sum, 'var(--accent-green)', 'green')}
          </>
        )}

        {arrow(a, 'var(--accent-blue)', 'blue')}
        {arrow(b, 'var(--accent-teal)', 'teal')}
        {scaledVisible && (
          <line x1={cx} y1={cy} x2={X(scaled[0])} y2={Y(scaled[1])} stroke="var(--accent-green)" strokeWidth="1.2" strokeDasharray="5 3" opacity="0.8" />
        )}
        {/* a − b is the step from b's tip to a's tip, so it is drawn where the
            learner can read subtraction off the picture and not only off the list. */}
        <line
          x1={X(b[0])}
          y1={Y(b[1])}
          x2={X(a[0])}
          y2={Y(a[1])}
          stroke="var(--text-muted)"
          strokeWidth="1.1"
          strokeDasharray="4 3"
        />

        <circle cx={X(a[0])} cy={Y(a[1])} r="4" fill="var(--accent-blue)" />
        <circle cx={X(b[0])} cy={Y(b[1])} r="4" fill="var(--accent-teal)" />

        <text x={X(a[0]) + dxAt(a[0])} y={Y(a[1]) - 7} textAnchor={anchorAt(a[0])} fontSize="10" fill="var(--accent-blue)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
          a
        </text>
        <text x={X(b[0]) + dxAt(b[0])} y={Y(b[1]) - 7} textAnchor={anchorAt(b[0])} fontSize="10" fill="var(--accent-teal)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
          b
        </text>
        <text x={X(difference[0] / 2) + 8} y={Y(difference[1] / 2) - 6} fontSize="9" fill="var(--text-muted)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
          a − b
        </text>
        {scaledVisible && (
          <text x={X(scaled[0]) + dxAt(scaled[0])} y={Y(scaled[1]) + 14} textAnchor={anchorAt(scaled[0])} fontSize="9" fill="var(--accent-green)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
            k b
          </text>
        )}
        {showSum && (
          <text x={X(sum[0]) + dxAt(sum[0])} y={Y(sum[1]) + 14} textAnchor={anchorAt(sum[0])} fontSize="10" fill="var(--accent-green)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
            a + b
          </text>
        )}

        <XAxisTicks ticks={[-3, -2, -1, 1, 2, 3]} scale={X} at={cy} />
        <YAxisTicks ticks={[-3, -2, -1, 1, 2, 3]} scale={Y} at={cx} />
      </svg>

      <div className="mt-3 space-y-2">
        <div className="grid grid-cols-2 gap-2">
          <NumberField label="a x" value={a[0]} min={-4} max={4} onChange={(v) => setA([v, a[1]])} />
          <NumberField label="a y" value={a[1]} min={-4} max={4} onChange={(v) => setA([a[0], v])} />
          <NumberField label="b x" value={b[0]} min={-4} max={4} onChange={(v) => setB([v, b[1]])} />
          <NumberField label="b y" value={b[1]} min={-4} max={4} onChange={(v) => setB([b[0], v])} />
        </div>
        <Slider label="k" value={k} min={-3} max={3} step={0.1} onChange={setK} />
        <Toggle label="Show the parallelogram" checked={showSum} onChange={setShowSum} />
      </div>
    </FigureFrame>
  )
}
