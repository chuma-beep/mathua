'use client'

import { useMemo, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Slider } from '../Controls'
import { XAxisTicks, YAxisTicks } from '../Axis'
import { conicSchema } from '../schemas'
import type { ConicType } from '../math/conic'
import { classifyQuadratic, discriminant, quadraticForm, xIntercepts, yIntercepts } from '../math/conic'
import { contourSegments } from '../math/contour'
import { niceTicks, scaleLinear } from '../scale'
import { fmt } from '../format'
import type { TemplateProps } from '../types'

const TYPE_LABELS: Record<ConicType, string> = {
  circle: 'circle',
  ellipse: 'ellipse',
  hyperbola: 'hyperbola',
  parallel: 'two parallel lines',
  empty: 'no real points',
}

const W = 340
const H = 300
const PAD = 30
const EXTENT = 2.8

export default function ConicSection({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => conicSchema.parse(params), [params])
  const [A, setA] = useState(defaults.A)
  const [B, setB] = useState(defaults.B)
  const [C, setC] = useState(defaults.C)

  const shape = classifyQuadratic(A, B, C)
  const segments = useMemo(
    () =>
      contourSegments(
        (x, y) => quadraticForm(A, B, C, x, y) - 1,
        { x0: -EXTENT, x1: EXTENT, y0: -EXTENT, y1: EXTENT },
        0,
        96,
        96,
      ),
    [A, B, C],
  )

  const X = scaleLinear(-EXTENT, EXTENT, PAD, W - PAD)
  const Y = scaleLinear(EXTENT, -EXTENT, PAD, H - PAD)
  const ticks = niceTicks(-EXTENT, EXTENT, 4).filter((v) => Math.abs(v) > 1e-9)

  const reset = () => {
    setA(defaults.A)
    setB(defaults.B)
    setC(defaults.C)
  }

  const signs =
    shape.squaredTerms < 2 ? 'one squared term' : shape.sameSign ? 'same sign' : 'opposite signs'
  const intercepts = [...xIntercepts(A, C).map((v) => `x ${fmt(v)}`), ...yIntercepts(A, C).map((v) => `y ${fmt(v)}`)]

  return (
    <FigureFrame
      title={title ?? 'The discriminant decides the conic'}
      hint="Change A, B and C. B only rotates the figure; A and C decide its kind."
      caption={
        caption ??
        'The curve is Ax² + Bxy + Cy² = 1. B² − 4AC is negative for an ellipse and positive for a hyperbola, and the cross term can only turn the figure, never change which of the two it is. When the discriminant is zero the equation loses a direction and the curve becomes a pair of parallel lines, because a parabola would need a linear term this equation does not have.'
      }
      onReset={reset}
      readout={[
        { label: 'A, B, C', value: `${fmt(A)}, ${fmt(B)}, ${fmt(C)}` },
        { label: 'discriminant', value: fmt(discriminant(A, B, C), 3) },
        { label: 'squared terms', value: String(shape.squaredTerms) },
        { label: 'signs', value: signs },
        { label: 'type', value: TYPE_LABELS[shape.type] },
        {
          label: 'eccentricity',
          value: Number.isFinite(shape.eccentricity) ? fmt(shape.eccentricity, 3) : '—',
        },
        { label: 'intercepts', value: intercepts.length > 0 ? intercepts.join(', ') : 'none' },
      ]}
    >
      <svg
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        role="img"
        aria-label="A conic section with its axes and discriminant"
      >
        <line x1={PAD} y1={Y(0)} x2={W - PAD} y2={Y(0)} stroke="var(--text-muted)" strokeWidth="0.8" />
        <line x1={X(0)} y1={PAD} x2={X(0)} y2={H - PAD} stroke="var(--text-muted)" strokeWidth="0.8" />

        {segments.map((seg, i) => (
          <line
            key={i}
            x1={X(seg[0][0])}
            y1={Y(seg[0][1])}
            x2={X(seg[1][0])}
            y2={Y(seg[1][1])}
            stroke="var(--accent-blue)"
            strokeWidth="1.7"
            strokeLinecap="round"
          />
        ))}

        <XAxisTicks ticks={ticks} scale={X} at={Y(0)} />
        <YAxisTicks ticks={ticks} scale={Y} at={X(0)} />

        <text x={PAD + 2} y={PAD + 10} fontSize="10" fill="var(--accent-blue)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
          {TYPE_LABELS[shape.type]}
        </text>
        {segments.length === 0 && (
          <text x={W / 2} y={H / 2} textAnchor="middle" fontSize="10" fill="var(--text-muted)">
            no real points
          </text>
        )}
      </svg>

      <div className="mt-3 space-y-2">
        <Slider label="A" value={A} min={-2} max={2} step={0.05} onChange={setA} />
        <Slider label="B" value={B} min={-2} max={2} step={0.05} onChange={setB} />
        <Slider label="C" value={C} min={-2} max={2} step={0.05} onChange={setC} />
      </div>
    </FigureFrame>
  )
}
