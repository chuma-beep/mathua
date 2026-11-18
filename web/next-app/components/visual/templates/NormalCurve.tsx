'use client'

import { useCallback, useMemo, useRef, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Slider, Toggle } from '../Controls'
import { normalSchema } from '../schemas'
import { normalPdf, normalProbability } from '../math/normal'
import { clamp, invertLinear, pathFromPoints, scaleLinear } from '../scale'
import { fmt } from '../format'
import { useDrag } from '../useDrag'
import type { TemplateProps } from '../types'

const W = 340
const H = 240
const PAD = 26

export default function NormalCurve({ params, title, caption }: TemplateProps) {
  const initial = useMemo(() => normalSchema.parse(params), [params])
  const [mu, setMu] = useState(initial.mu)
  const [sigma, setSigma] = useState(initial.sigma)
  const [x, setX] = useState(initial.x)
  const [bands, setBands] = useState(initial.bands)
  const svgRef = useRef<SVGSVGElement | null>(null)

  const lo = mu - 4 * sigma
  const hi = mu + 4 * sigma
  const peak = normalPdf(mu, mu, sigma)
  const X = scaleLinear(lo, hi, PAD, W - PAD)
  const Y = scaleLinear(0, peak * 1.16, H - PAD, PAD)
  const invertX = invertLinear(lo, hi, PAD, W - PAD)

  const samples = useMemo(() => {
    const out: Array<[number, number]> = []
    for (let i = 0; i <= 240; i++) {
      const xv = lo + ((hi - lo) * i) / 240
      out.push([xv, normalPdf(xv, mu, sigma)])
    }
    return out
  }, [lo, hi, mu, sigma])

  const tailEnd = clamp(x, lo, hi)
  const tail = samples.filter(([xv]) => xv <= tailEnd)
  const areaPath =
    tail.length > 1
      ? `${pathFromPoints(tail.map(([xv, yv]) => [X(xv), Y(yv)]))} L ${X(tailEnd)} ${Y(0)} L ${X(lo)} ${Y(0)} Z`
      : ''

  const onDrag = useCallback(
    (point: { x: number; y: number }) => {
      setX(clamp(invertX(point.x), -12, 12))
    },
    [invertX],
  )
  useDrag(svgRef, onDrag)

  const reset = () => {
    setMu(initial.mu)
    setSigma(initial.sigma)
    setX(initial.x)
    setBands(initial.bands)
  }

  const z = (x - mu) / sigma
  const p = normalProbability(x, mu, sigma)

  return (
    <FigureFrame
      title={title ?? 'The normal curve'}
      caption={
        caption ??
        'The shaded tail is P(X ≤ x). Moving μ slides the curve, σ widens it, and x moves the cut. The band probability stays the same after a change of μ or σ.'
      }
      onReset={reset}
      readout={[
        { label: 'z', value: fmt(z) },
        { label: 'P(X ≤ x)', value: fmt(p, 4) },
        { label: 'P(X > x)', value: fmt(1 - p, 4) },
        { label: 'μ ± σ', value: '0.6827' },
        { label: 'μ', value: fmt(mu) },
        { label: 'σ', value: fmt(sigma) },
      ]}
    >
      <svg
        ref={svgRef}
        viewBox={`0 0 ${W} ${H}`}
        className="h-auto w-full touch-none font-mono"
        style={{ touchAction: 'none' }}
        role="img"
        aria-label="A normal density with a shaded left tail"
      >
        {bands &&
          [-3, -2, -1, 1, 2, 3].map((k) => (
            <g key={k}>
              <line
                x1={X(mu + k * sigma)}
                y1={PAD}
                x2={X(mu + k * sigma)}
                y2={H - PAD}
                stroke="var(--border-strong)"
                strokeWidth="0.8"
                strokeDasharray="3 3"
              />
              <text x={X(mu + k * sigma)} y={H - PAD + 12} textAnchor="middle" fontSize="9" fill="var(--text-muted)">
                {k > 0 ? `${k}σ` : `−${Math.abs(k)}σ`}
              </text>
            </g>
          ))}

        <line x1={PAD} y1={Y(0)} x2={W - PAD} y2={Y(0)} stroke="var(--text-muted)" strokeWidth="0.8" />
        {areaPath && <path d={areaPath} fill="var(--accent-blue)" fillOpacity="0.2" />}
        <path d={pathFromPoints(samples.map(([xv, yv]) => [X(xv), Y(yv)]))} fill="none" stroke="var(--text-primary)" strokeWidth="1.8" />
        <line x1={X(clamp(x, lo, hi))} y1={PAD} x2={X(clamp(x, lo, hi))} y2={H - PAD} stroke="var(--accent-blue)" strokeWidth="1.2" />
        <circle cx={X(clamp(x, lo, hi))} cy={Y(normalPdf(clamp(x, lo, hi), mu, sigma))} r="3.6" fill="var(--accent-blue)" />

        <text x={X(mu)} y={PAD + 2} textAnchor="middle" fontSize="9" fill="var(--text-muted)">μ</text>
        <text x={X(clamp(x, lo, hi)) + 6} y={PAD + 16} fontSize="10" fill="var(--accent-blue)">x</text>
        <text x={W - PAD} y={H - PAD + 12} textAnchor="end" fontSize="9" fill="var(--text-muted)">drag to move x</text>
      </svg>

      <div className="mt-3 space-y-2">
        <Slider label="μ" value={mu} min={-4} max={4} step={0.05} onChange={setMu} />
        <Slider label="σ" value={sigma} min={0.1} max={4} step={0.05} onChange={setSigma} />
        <Slider label="x" value={x} min={-10} max={10} step={0.05} onChange={setX} />
        <Toggle label="Show σ bands" checked={bands} onChange={setBands} />
      </div>
    </FigureFrame>
  )
}
