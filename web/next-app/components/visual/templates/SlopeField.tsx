'use client'

import { useMemo, useRef, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Slider } from '../Controls'
import { XAxisTicks, YAxisTicks } from '../Axis'
import { odeSlopeSchema } from '../schemas'
import { ODE_FIELDS, integrateSlopeField, odeField, slopeAt, type OdeFieldId, type OdeMethod } from '../math/ode'
import { clamp, invertLinear, pathFromPoints, scaleLinear } from '../scale'
import { fmt } from '../format'
import { useDrag } from '../useDrag'
import type { TemplateProps } from '../types'

const W = 360
const H = 360
const PAD = 28
const EXTENT = 3
const BOUNDS = { x0: -EXTENT, x1: EXTENT, y0: -EXTENT, y1: EXTENT }

export default function SlopeField({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => odeSlopeSchema.parse(params), [params])
  const [field, setField] = useState<OdeFieldId>(defaults.field)
  const [x0, setX0] = useState(defaults.x0)
  const [y0, setY0] = useState(defaults.y0)
  const [step, setStep] = useState(defaults.step)
  const [method, setMethod] = useState<OdeMethod>(defaults.method)
  const svgRef = useRef<SVGSVGElement | null>(null)

  const info = odeField(field)
  const trajectory = useMemo(
    () => integrateSlopeField(field, x0, y0, {
      xMin: BOUNDS.x0,
      xMax: BOUNDS.x1,
      yMin: BOUNDS.y0,
      yMax: BOUNDS.y1,
      step,
      method,
    }),
    [field, x0, y0, step, method],
  )

  const X = scaleLinear(-EXTENT, EXTENT, PAD, W - PAD)
  const Y = scaleLinear(EXTENT, -EXTENT, PAD, H - PAD)
  const toX = invertLinear(-EXTENT, EXTENT, PAD, W - PAD)
  const toY = invertLinear(EXTENT, -EXTENT, PAD, H - PAD)
  // +-3 sit on the frame, so their labels were half outside it.
  const ticks = [-2, -1, 0, 1, 2]

  useDrag(svgRef, (pointer) => {
    setX0(clamp(toX(pointer.x), -3, 3))
    setY0(clamp(toY(pointer.y), -3, 3))
  })

  const reset = () => {
    setField(defaults.field)
    setX0(defaults.x0)
    setY0(defaults.y0)
    setStep(defaults.step)
    setMethod(defaults.method)
  }

  const segments = useMemo(() => {
    const result: Array<{ x1: number; y1: number; x2: number; y2: number; slope: number }> = []
    const count = 15
    for (let i = 0; i <= count; i++) {
      for (let j = 0; j <= count; j++) {
        const x = -EXTENT + (2 * EXTENT * i) / count
        const y = -EXTENT + (2 * EXTENT * j) / count
        const slope = slopeAt(field, x, y)
        if (!Number.isFinite(slope)) continue
        const length = Math.hypot(1, slope)
        const ux = 1 / length
        const uy = slope / length
        // Half-length in data units, not pixels. It was 5, which in a window of
        // +-3 drew every mark across the whole figure and turned the field into
        // a web of crossing lines. One grid cell is 0.4, so 0.16 is a mark.
        const half = 0.16
        result.push({
          x1: X(x - ux * half),
          y1: Y(y - uy * half),
          x2: X(x + ux * half),
          y2: Y(y + uy * half),
          slope,
        })
      }
    }
    return result
  }, [field, X, Y])

  const curve = trajectory.points.map((point) => [X(point.x), Y(point.y)] as [number, number])
  const start = { x: x0, y: y0 }
  const finish = trajectory.points[trajectory.points.length - 1] ?? start
  const slope = slopeAt(field, start.x, start.y)

  return (
    <FigureFrame
      title={title ?? 'The slope field is the differential equation'}
      hint="Drag the starting point. Every short segment is the slope dy/dx at that point, and the curve follows those slopes numerically."
      caption={
        caption ??
        'Each segment is drawn with slope f(x, y), so the field is the differential equation itself. The curve is integrated from the starting point in both directions. Change the method or step size to see numerical error, not a different equation.'
      }
      onReset={reset}
      readout={[
        { label: 'equation', value: info.equation },
        { label: 'start', value: `(${fmt(start.x, 2)}, ${fmt(start.y, 2)})` },
        { label: 'slope at start', value: Number.isFinite(slope) ? fmt(slope, 3) : 'undefined' },
        { label: 'method', value: method === 'rk4' ? 'RK4' : 'Euler' },
        { label: 'step h', value: fmt(step, 2) },
        { label: 'integrated steps', value: String(trajectory.steps) },
        { label: 'end', value: `(${fmt(finish.x, 2)}, ${fmt(finish.y, 2)})` },
        { label: 'status', value: trajectory.status },
        ...(trajectory.note ? [{ label: 'note', value: trajectory.note }] : []),
      ]}
    >
      <svg
        ref={svgRef}
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        style={{ touchAction: 'none' }}
        role="img"
        aria-label={`Slope field for ${info.equation} with the solution through (${fmt(x0, 2)}, ${fmt(y0, 2)})`}
      >
        <line x1={PAD} y1={Y(0)} x2={W - PAD} y2={Y(0)} stroke="var(--text-muted)" strokeWidth="0.8" />
        <line x1={X(0)} y1={PAD} x2={X(0)} y2={H - PAD} stroke="var(--text-muted)" strokeWidth="0.8" />
        <XAxisTicks ticks={ticks} scale={X} at={Y(0)} />
        <YAxisTicks ticks={ticks} scale={Y} at={X(0)} />

        {segments.map((segment, index) => (
          <line
            key={index}
            x1={segment.x1}
            y1={segment.y1}
            x2={segment.x2}
            y2={segment.y2}
            stroke="var(--text-muted)"
            strokeWidth="1"
            opacity="0.7"
          />
        ))}

        {curve.length > 1 ? (
          <path d={pathFromPoints(curve)} fill="none" stroke="var(--accent-blue)" strokeWidth="2" />
        ) : null}

        <circle cx={X(x0)} cy={Y(y0)} r="4.5" fill="var(--accent-green)" />
        <text
          x={X(x0) + 8}
          y={Y(y0) - 7}
          fontSize="9"
          fill="var(--text-primary)"
          stroke="var(--bg)"
          strokeWidth="3"
          paintOrder="stroke"
        >
          start
        </text>
      </svg>

      <div className="mt-3 space-y-2">
        <div className="flex items-center gap-2">
          <label htmlFor="slope-field" className="w-10 shrink-0 font-mono text-[10px] text-mathua-muted">
            Field
          </label>
          <select
            id="slope-field"
            aria-label="Field"
            value={field}
            onChange={(event) => setField(event.target.value as OdeFieldId)}
            className="h-7 min-w-0 border border-mathua-border bg-transparent px-1 font-mono text-[10px] text-mathua-secondary"
          >
            {ODE_FIELDS.map((item) => (
              <option key={item.id} value={item.id}>
                {item.equation}
              </option>
            ))}
          </select>
        </div>
        <div className="flex items-center gap-2">
          <label htmlFor="slope-method" className="w-10 shrink-0 font-mono text-[10px] text-mathua-muted">
            Method
          </label>
          <select
            id="slope-method"
            aria-label="Method"
            value={method}
            onChange={(event) => setMethod(event.target.value as OdeMethod)}
            className="h-7 min-w-0 border border-mathua-border bg-transparent px-1 font-mono text-[10px] text-mathua-secondary"
          >
            <option value="rk4">RK4</option>
            <option value="euler">Euler</option>
          </select>
        </div>
        <div className="grid grid-cols-1 gap-2">
          <Slider label="x₀" value={x0} min={-3} max={3} step={0.05} onChange={setX0} />
          <Slider label="y₀" value={y0} min={-3} max={3} step={0.05} onChange={setY0} />
          <Slider label="step h" value={step} min={0.02} max={0.5} step={0.01} onChange={setStep} />
        </div>
      </div>
    </FigureFrame>
  )
}
