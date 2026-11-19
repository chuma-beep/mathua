'use client'

import { useMemo, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Select, Slider } from '../Controls'
import { riemannSchema } from '../schemas'
import { evalPoly, definiteIntegral } from '../math/poly'
import { riemannSum, type RiemannRule } from '../math/integrate'
import { clamp, niceTicks, pathFromPoints, scaleLinear } from '../scale'
import { XAxisTicks, YAxisTicks } from '../Axis'
import { fmt } from '../format'
import type { TemplateProps } from '../types'

const W = 320
const H = 250
const PAD = 28

export default function RiemannSum({ params, title, caption }: TemplateProps) {
  const initial = useMemo(() => riemannSchema.parse(params), [params])
  const [n, setN] = useState(() => clamp(initial.n, 1, 256))
  const [rule, setRule] = useState<RiemannRule>(initial.rule)

  const { coeffs } = initial
  // An interval is a set, so an authored a > b describes the same interval as
  // b > a. Without this, a reversed pair drew a mirrored curve with zero-width
  // rectangles while the readout still showed numbers: the picture and the
  // arithmetic disagreed about which way round the interval was.
  const a = Math.min(initial.a, initial.b)
  const b = Math.max(initial.a, initial.b)
  const exact = definiteIntegral(coeffs, a, b)
  const sum = riemannSum(coeffs, a, b, n, rule)
  const dx = (b - a) / n

  const samples = useMemo(() => {
    const out: Array<[number, number]> = []
    for (let i = 0; i <= 200; i++) {
      const xv = a + ((b - a) * i) / 200
      out.push([xv, evalPoly(coeffs, xv)])
    }
    return out
  }, [coeffs, a, b])

  const yValues = samples.map((p) => p[1]).concat([0])
  const yMin = Math.min(...yValues) - 0.4
  const yMax = Math.max(...yValues) + 0.4
  const X = scaleLinear(a, b, PAD, W - PAD)
  const Y = scaleLinear(yMin, yMax, H - PAD, PAD)
  const xTicks = niceTicks(a, b, 4)
  const yTicks = niceTicks(yMin, yMax, 4)

  const reset = () => {
    setN(clamp(initial.n, 1, 256))
    setRule(initial.rule)
  }

  const rects = Array.from({ length: n }, (_, i) => {
    const x0 = a + i * dx
    const x1 = x0 + dx
    const sampleX = rule === 'left' ? x0 : rule === 'right' ? x1 : (x0 + x1) / 2
    const y = evalPoly(coeffs, sampleX)
    const top = Math.min(Y(0), Y(y))
    const height = Math.abs(Y(y) - Y(0))
    return { key: i, x: X(x0), width: Math.max(0, X(x1) - X(x0)), top, height }
  })

  return (
    <FigureFrame
      title={title ?? 'Riemann sum'}
      caption={
        caption ??
        'Each rectangle samples the function once. Increasing n narrows every rectangle, and the sum approaches the exact area.'
      }
      onReset={reset}
      readout={[
        { label: 'sum', value: fmt(sum, 4) },
        { label: 'exact', value: fmt(exact, 4) },
        { label: 'error', value: fmt(Math.abs(sum - exact), 4) },
        { label: 'n', value: String(n) },
        { label: 'Δx', value: fmt(dx, 3) },
      ]}
    >
      <svg viewBox={`0 0 ${W} ${H}`} className="mx-auto block h-auto w-full max-w-[420px] font-mono" role="img" aria-label="Rectangles approximating the area under a curve">
        {yMin < 0 && yMax > 0 && (
          <line x1={PAD} y1={Y(0)} x2={W - PAD} y2={Y(0)} stroke="var(--text-muted)" strokeWidth="0.8" />
        )}
        {a < 0 && b > 0 && <line x1={X(0)} y1={PAD} x2={X(0)} y2={H - PAD} stroke="var(--text-muted)" strokeWidth="0.8" />}
        <XAxisTicks ticks={xTicks} scale={X} at={H - PAD} format={(v) => fmt(v, 1)} />
        <YAxisTicks ticks={yTicks} scale={Y} at={PAD} format={(v) => fmt(v, 1)} />

        {rects.map((r) => (
          <rect
            key={r.key}
            x={r.x}
            y={r.top}
            width={r.width}
            height={r.height}
            fill="var(--accent-blue)"
            fillOpacity="0.14"
            stroke="var(--accent-blue)"
            strokeWidth="0.8"
          />
        ))}
        <path d={pathFromPoints(samples.map(([xv, yv]) => [X(xv), Y(yv)]))} fill="none" stroke="var(--text-primary)" strokeWidth="1.8" />
        <text x={PAD + 12} y={PAD + 10} fontSize="10" fill="var(--accent-blue)">n = {n}</text>
        <text x={W - PAD} y={PAD + 10} textAnchor="end" fontSize="10" fill="var(--accent-teal)">
          {rule} sum {fmt(sum, 3)}
        </text>
      </svg>

      <div className="mt-3 space-y-2">
        <Slider
          label="n"
          value={n}
          min={1}
          max={256}
          step={1}
          onChange={(value) => setN(clamp(Math.round(value), 1, 256))}
        />
        <Select
          label="sample"
          value={rule}
          options={[
            { value: 'left', label: 'left' },
            { value: 'right', label: 'right' },
            { value: 'mid', label: 'midpoint' },
          ]}
          onChange={(v) => setRule(v as RiemannRule)}
        />
      </div>
    </FigureFrame>
  )
}
