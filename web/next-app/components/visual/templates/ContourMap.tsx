'use client'

import { useCallback, useMemo, useRef, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { NumberField, Select, Slider } from '../Controls'
import { XAxisTicks, YAxisTicks } from '../Axis'
import { contourSchema } from '../schemas'
import { FIELD_LABELS, fieldGradient, fieldLevels, fieldValue, type Field } from '../math/field'
import { contourSegments } from '../math/contour'
import { clamp, scaleLinear } from '../scale'
import { fmt } from '../format'
import { useDrag } from '../useDrag'
import type { TemplateProps } from '../types'

const SIZE = 320
const PAD = 30
const EXTENT = 3
const DEG = 180 / Math.PI

const FIELD_OPTIONS: Array<{ value: Field; label: string }> = [
  { value: 'bowl', label: FIELD_LABELS.bowl },
  { value: 'saddle', label: FIELD_LABELS.saddle },
  { value: 'product', label: FIELD_LABELS.product },
  { value: 'ripple', label: FIELD_LABELS.ripple },
]

export default function ContourMap({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => {
    const p = contourSchema.parse(params)
    return {
      field: p.field as Field,
      levels: p.levels,
      probe: [Number(p.probe[0]), Number(p.probe[1])] as [number, number],
    }
  }, [params])

  const [field, setField] = useState<Field>(defaults.field)
  const [levels, setLevels] = useState(defaults.levels)
  const [probe, setProbe] = useState<[number, number]>(defaults.probe)
  const svgRef = useRef<SVGSVGElement | null>(null)

  const X = scaleLinear(-EXTENT, EXTENT, PAD, SIZE - PAD)
  const Y = scaleLinear(EXTENT, -EXTENT, PAD, SIZE - PAD)
  const cx = X(0)
  const cy = Y(0)
  const unit = (SIZE - PAD * 2) / (2 * EXTENT)

  const levelValues = useMemo(() => fieldLevels(field, EXTENT, levels), [field, levels])
  const curves = useMemo(() => {
    // Built inside the memo so the dependency array stays static: an object literal
    // in the outer scope is a new reference every render and would defeat the memo.
    const bounds = { x0: -EXTENT, x1: EXTENT, y0: -EXTENT, y1: EXTENT }
    return levelValues.map((level) => ({
      level,
      segments: contourSegments((x, y) => fieldValue(field, x, y), bounds, level, 90, 90),
    }))
  }, [field, levelValues])

  const value = fieldValue(field, probe[0], probe[1])
  const gradient = fieldGradient(field, probe[0], probe[1])
  const magnitude = Math.hypot(gradient[0], gradient[1])
  // At a critical point the gradient is the zero vector, and atan2(0, 0) would
  // report 0 degrees — a direction the function does not have. There is no arrow
  // to be perpendicular to and no uphill to point at, so both are undefined.
  const atCritical = magnitude < 1e-9
  const angle = atCritical ? null : (Math.atan2(gradient[1], gradient[0]) * DEG + 360) % 360

  const onDrag = useCallback(
    (point: { x: number; y: number }) => {
      const px = clamp((point.x - cx) / unit, -EXTENT, EXTENT)
      const py = clamp((cy - point.y) / unit, -EXTENT, EXTENT)
      setProbe([Number(px.toFixed(3)), Number(py.toFixed(3))])
    },
    [cx, cy, unit],
  )
  useDrag(svgRef, onDrag)

  const reset = () => {
    setField(defaults.field)
    setLevels(defaults.levels)
    setProbe(defaults.probe)
  }

  // The arrow shows direction at a readable length: a gradient of 8 and one of 0.5
  // both fit the figure, so the length is capped and the number carries the size.
  const arrowLength = magnitude === 0 ? 0 : Math.min(1.15, 0.35 + Math.min(magnitude, 6) * 0.13)
  const dirX = magnitude === 0 ? 0 : gradient[0] / magnitude
  const dirY = magnitude === 0 ? 0 : gradient[1] / magnitude
  const tipX = probe[0] + dirX * arrowLength
  const tipY = probe[1] + dirY * arrowLength
  const tangentSpan = 1.3

  return (
    <FigureFrame
      title={title ?? 'Level curves and the gradient'}
      hint="Drag the probe. The gradient points across the level curve, never along it."
      caption={
        caption ??
        'Each curve is one value of the function. The arrow is the gradient there: it is perpendicular to the curve through the probe, points uphill and its length grows where the curves crowd together.'
      }
      onReset={reset}
      readout={[
        { label: 'f(x, y)', value: fmt(value, 3) },
        { label: '∂f/∂x', value: fmt(gradient[0], 3) },
        { label: '∂f/∂y', value: fmt(gradient[1], 3) },
        { label: '|∇f|', value: fmt(magnitude, 3) },
        { label: 'gradient angle', value: angle === null ? 'undefined' : `${fmt(angle, 0)}°` },
        { label: 'levels', value: String(levels) },
      ]}
    >
      <svg
        ref={svgRef}
        viewBox={`0 0 ${SIZE} ${SIZE}`}
        className="mx-auto block h-auto w-full max-w-[420px] cursor-grab font-mono active:cursor-grabbing"
        style={{ touchAction: 'none' }}
        role="img"
        aria-label="Contour lines of a two-variable function with the gradient at a probe point"
      >
        <defs>
          <marker id="contour-arrow" markerWidth="6" markerHeight="6" refX="5" refY="3" orient="auto">
            <path d="M0,0 L6,3 L0,6 Z" fill="var(--accent-blue)" />
          </marker>
        </defs>

        <line x1={PAD} y1={cy} x2={SIZE - PAD} y2={cy} stroke="var(--text-muted)" strokeWidth="0.8" />
        <line x1={cx} y1={PAD} x2={cx} y2={SIZE - PAD} stroke="var(--text-muted)" strokeWidth="0.8" />
        <XAxisTicks ticks={[-2, -1, 1, 2]} scale={X} at={cy} />
        <YAxisTicks ticks={[-2, -1, 1, 2]} scale={Y} at={cx} />
        {curves.map(({ level, segments }, index) => (
          <g key={level}>
            {segments.map((segment, i) => (
              <line
                key={i}
                x1={X(segment[0][0])}
                y1={Y(segment[0][1])}
                x2={X(segment[1][0])}
                y2={Y(segment[1][1])}
                stroke="var(--accent-teal)"
                strokeWidth="1.2"
                strokeLinecap="round"
                opacity={0.28 + (0.5 * (index + 1)) / Math.max(1, curves.length)}
              />
            ))}
          </g>
        ))}


        {/* The level direction at the probe: perpendicular to the gradient, so the
            perpendicularity is visible rather than asserted in the caption. Hidden at a
            critical point, where the gradient is zero and every direction is tangent. */}
        {!atCritical && (
          <line
            x1={X(probe[0] - dirY * tangentSpan)}
            y1={Y(probe[1] + dirX * tangentSpan)}
            x2={X(probe[0] + dirY * tangentSpan)}
            y2={Y(probe[1] - dirX * tangentSpan)}
            stroke="var(--accent-green)"
            strokeWidth="0.9"
            strokeDasharray="4 3"
            opacity="0.8"
          />
        )}

        {magnitude > 0 && (
          <line x1={X(probe[0])} y1={Y(probe[1])} x2={X(tipX)} y2={Y(tipY)} stroke="var(--accent-blue)" strokeWidth="1.8" markerEnd="url(#contour-arrow)" />
        )}

        {atCritical && (
          <text x={SIZE - PAD} y={PAD + 10} textAnchor="end" fontSize="9" fill="var(--text-muted)">
            critical point: no gradient
          </text>
        )}

        <circle cx={X(probe[0])} cy={Y(probe[1])} r="4" fill="var(--accent-blue)" />
        <text x={X(probe[0]) + 7} y={Y(probe[1]) - 7} fontSize="10" fill="var(--accent-blue)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
          p
        </text>
      </svg>

      <div className="mt-3 space-y-2">
        <Select label="Field" value={field} options={FIELD_OPTIONS} onChange={(v) => setField(v as Field)} />
        <Slider label="levels" value={levels} min={3} max={12} step={1} onChange={(v) => setLevels(Math.round(v))} />
        <div className="grid grid-cols-2 gap-2">
          <NumberField label="x" value={probe[0]} min={-EXTENT} max={EXTENT} onChange={(v) => setProbe([v, probe[1]])} />
          <NumberField label="y" value={probe[1]} min={-EXTENT} max={EXTENT} onChange={(v) => setProbe([probe[0], v])} />
        </div>
      </div>
    </FigureFrame>
  )
}
