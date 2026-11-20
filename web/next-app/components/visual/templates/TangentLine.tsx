'use client'

import { useCallback, useMemo, useRef, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Slider } from '../Controls'
import { tangentSchema } from '../schemas'
import { derivative, evalPoly } from '../math/poly'
import { clamp, invertLinear, niceTicks, pathFromPoints, scaleLinear } from '../scale'
import { XAxisTicks, YAxisTicks } from '../Axis'
import { fmt } from '../format'
import { useDrag } from '../useDrag'
import type { TemplateProps } from '../types'

const W = 320
const H = 260
const PAD = 26

export default function TangentLine({ params, title, caption }: TemplateProps) {
  const initial = useMemo(() => tangentSchema.parse(params), [params])
  const [x, setX] = useState(() => clamp(initial.x, initial.window[0], initial.window[1]))
  const [h, setH] = useState(initial.h)
  const svgRef = useRef<SVGSVGElement | null>(null)

  const { coeffs, window } = initial
  const slopeCoeffs = derivative(coeffs)
  const f = evalPoly(coeffs, x)
  const fh = evalPoly(coeffs, x + h)
  const tangent = evalPoly(slopeCoeffs, x)
  const secant = (fh - f) / h

  // The right edge grows to hold x + h. Without this, dragging x to the end of the
  // window with a large h put the secant's second point outside the viewBox: the
  // curve, the dot and the label were clipped while the readout still reported a
  // slope, so the figure stopped showing the thing its caption claims.
  const domain = useMemo<[number, number]>(() => [window[0], Math.max(window[1], x + h)], [window, x, h])

  const samples = useMemo(() => {
    const out: Array<[number, number]> = []
    for (let i = 0; i <= 240; i++) {
      const xv = domain[0] + ((domain[1] - domain[0]) * i) / 240
      out.push([xv, evalPoly(coeffs, xv)])
    }
    return out
  }, [coeffs, domain])

  const yValues = useMemo(() => {
    const values = samples.map((p) => p[1])
    const lineAt = (m: number, at: number) => f + m * (at - x)
    values.push(lineAt(secant, domain[0]), lineAt(secant, domain[1]))
    values.push(lineAt(tangent, domain[0]), lineAt(tangent, domain[1]))
    values.push(f, fh)
    return values.filter(Number.isFinite)
  }, [samples, secant, tangent, f, fh, domain, x])

  const yMin = Math.min(...yValues) - 0.6
  const yMax = Math.max(...yValues) + 0.6
  const X = scaleLinear(domain[0], domain[1], PAD, W - PAD)
  const Y = scaleLinear(yMin, yMax, H - PAD, PAD)
  const invertX = invertLinear(domain[0], domain[1], PAD, W - PAD)
  const xTicks = niceTicks(domain[0], domain[1], 5)
  const yTicks = niceTicks(yMin, yMax, 4)

  const onDrag = useCallback(
    (point: { x: number; y: number }) => {
      setX(clamp(invertX(point.x), window[0], window[1]))
    },
    [invertX, window],
  )
  useDrag(svgRef, onDrag)

  const reset = () => {
    setX(clamp(initial.x, window[0], window[1]))
    setH(initial.h)
  }

  const curvePath = pathFromPoints(samples.map(([xv, yv]) => [X(xv), Y(yv)]))
  const line = (m: number) =>
    `M ${X(domain[0])} ${Y(f + m * (domain[0] - x))} L ${X(domain[1])} ${Y(f + m * (domain[1] - x))}`

  return (
    <FigureFrame
      title={title ?? 'Tangent and secant through a curve'}
      caption={
        caption ??
        'Drag the point or move the x slider. As h shrinks, the secant slope approaches the derivative at the point.'
      }
      onReset={reset}
      hint="Drag the plot, or move x and h with the sliders."
      readout={[
        { label: 'f(x)', value: fmt(f) },
        { label: "f'(x)", value: fmt(tangent) },
        { label: 'secant', value: fmt(secant) },
        { label: 'difference', value: fmt(Math.abs(secant - tangent), 3) },
        { label: 'h', value: fmt(h) },
      ]}
    >
      <svg
        ref={svgRef}
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] cursor-grab font-mono active:cursor-grabbing"
        style={{ touchAction: 'pan-y' }}
        role="img"
        aria-label="A curve with a movable point, its tangent and a secant"
      >
        <defs>
          <marker id="tangent-dot" markerWidth="6" markerHeight="6" refX="3" refY="3">
            <circle cx="3" cy="3" r="2.4" fill="var(--accent-blue)" />
          </marker>
        </defs>

        {yMin < 0 && yMax > 0 && (
          <line x1={PAD} y1={Y(0)} x2={W - PAD} y2={Y(0)} stroke="var(--text-muted)" strokeWidth="0.8" />
        )}
        {window[0] < 0 && window[1] > 0 && (
          <line x1={X(0)} y1={PAD} x2={X(0)} y2={H - PAD} stroke="var(--text-muted)" strokeWidth="0.8" />
        )}
        <XAxisTicks ticks={xTicks} scale={X} at={H - PAD} format={(v) => fmt(v, 1)} />
        <YAxisTicks ticks={yTicks} scale={Y} at={PAD} format={(v) => fmt(v, 1)} />

        <path d={curvePath} fill="none" stroke="var(--text-primary)" strokeWidth="1.8" />
        <path d={line(secant)} fill="none" stroke="var(--accent-teal)" strokeWidth="1.4" strokeDasharray="5 4" />
        <path d={line(tangent)} fill="none" stroke="var(--accent-blue)" strokeWidth="1.5" />

        <line x1={X(x)} y1={Y(f)} x2={X(x + h)} y2={Y(fh)} stroke="var(--accent-teal)" strokeWidth="1" opacity="0.7" />
        <circle cx={X(x)} cy={Y(f)} r="4" fill="var(--accent-blue)" />
        <circle cx={X(x + h)} cy={Y(fh)} r="3.2" fill="var(--accent-teal)" />

        <text x={X(x) - 8} y={Y(f) - 10} textAnchor="end" fontSize="10" fill="var(--accent-blue)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">x</text>
        <text
          x={X(x + h) + (X(x + h) > W - PAD - 30 ? -8 : 8)}
          y={Y(fh) + 16}
          textAnchor={X(x + h) > W - PAD - 30 ? 'end' : 'start'}
          fontSize="10"
          fill="var(--accent-teal)"
          stroke="var(--bg)"
          strokeWidth="3"
          paintOrder="stroke"
        >
          x + h
        </text>
        <text x={PAD + 2} y={PAD + 10} fontSize="10" fill="var(--accent-blue)">tangent</text>
        <text x={PAD + 2} y={PAD + 22} fontSize="10" fill="var(--accent-teal)">secant</text>
      </svg>

      <div className="mt-3 space-y-2">
        <Slider label="x" value={x} min={window[0]} max={window[1]} step={0.01} onChange={setX} />
        <Slider label="h" value={h} min={0.05} max={3} step={0.01} onChange={setH} />
      </div>
    </FigureFrame>
  )
}
