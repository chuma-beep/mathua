'use client'

import { useCallback, useMemo, useRef, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { NumberField, Select, Toggle } from '../Controls'
import { XAxisTicks } from '../Axis'
import { spreadSchema } from '../schemas'
import { deviations, mean, squaredDeviations, stdDev, sumSquares, variance } from '../math/spread'
import { clamp, scaleLinear } from '../scale'
import { fmt } from '../format'
import { useDrag } from '../useDrag'
import type { TemplateProps } from '../types'

const W = 340
const H = 250
const PAD = 26
const AXIS_Y = 176
const MIN_X = 0
const MAX_X = 10

export default function SpreadMeasure({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => spreadSchema.parse(params), [params])
  const [values, setValues] = useState<number[]>(defaults.values)
  const [squared, setSquared] = useState(defaults.squared)
  const [sample, setSample] = useState(defaults.sample)
  const [selected, setSelected] = useState(0)
  const svgRef = useRef<SVGSVGElement | null>(null)

  const X = scaleLinear(MIN_X, MAX_X, PAD, W - PAD)
  const invX = scaleLinear(PAD, W - PAD, MIN_X, MAX_X)

  const m = mean(values)
  const gaps = deviations(values)
  const squares = squaredDeviations(values)
  const ss = sumSquares(values)
  const v = variance(values, sample)
  const sd = stdDev(values, sample)

  const maxSquare = Math.max(0, ...squares)
  // A zero gap must draw as no bar. The old floor of one pixel gave every point a
  // visible bar even when every value equalled the mean, so a sample with variance
  // zero looked like five tiny squares rather than none.
  const hasSpread = maxSquare > 1e-12
  const barTop = 34
  const barHeight = (square: number) => (hasSpread ? (square / maxSquare) * (AXIS_Y - barTop - 16) : 0)

  const onDrag = useCallback(
    (point: { x: number; y: number }) => {
      const next = clamp(Math.round(invX(point.x) * 10) / 10, MIN_X, MAX_X)
      setValues((prev) => {
        let best = 0
        let bestDistance = Number.POSITIVE_INFINITY
        prev.forEach((value, i) => {
          const d = Math.abs(value - next)
          if (d < bestDistance) {
            bestDistance = d
            best = i
          }
        })
        const copy = [...prev]
        copy[best] = next
        return copy
      })
    },
    [invX],
  )
  useDrag(svgRef, onDrag)

  const setPoint = (index: number, value: number) => {
    setValues((prev) => {
      const copy = [...prev]
      copy[index] = clamp(value, MIN_X, MAX_X)
      return copy
    })
  }

  const addPoint = () => {
    if (values.length >= 8) return
    const far = clamp(m + (m >= 5 ? -1.8 : 1.8), MIN_X, MAX_X)
    setValues([...values, Math.round(far * 10) / 10])
    setSelected(values.length)
  }

  const removePoint = () => {
    if (values.length <= 2) return
    setValues(values.slice(0, -1))
    setSelected((prev) => Math.min(prev, values.length - 2))
  }

  const reset = () => {
    setValues(defaults.values)
    setSquared(defaults.squared)
    setSample(defaults.sample)
    setSelected(0)
  }

  const safeIndex = Math.min(selected, values.length - 1)

  return (
    <FigureFrame
      title={title ?? 'Standard deviation is a typical gap'}
      hint="Drag a point, or pick one and type its value. The bars are the squared gaps."
      caption={
        caption ??
        'The mean is the balance point, so the signed gaps add to zero. Variance averages their squares, which is why one distant point moves it more than several near ones, and the standard deviation puts the answer back in the units of the data.'
      }
      onReset={reset}
      readout={[
        { label: 'n', value: String(values.length) },
        { label: 'mean', value: fmt(m, 3) },
        { label: 'sum of squares', value: fmt(ss, 3) },
        { label: 'variance', value: fmt(v, 3) },
        { label: 'std dev', value: fmt(sd, 3) },
        { label: 'divides by', value: sample ? `n − 1 = ${values.length - 1}` : `n = ${values.length}` },
      ]}
    >
      <svg
        ref={svgRef}
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        style={{ touchAction: 'none' }}
        role="img"
        aria-label="Points on a number line, their gaps from the mean, and the squared gaps as bars"
      >
        <line x1={PAD} y1={AXIS_Y} x2={W - PAD} y2={AXIS_Y} stroke="var(--text-muted)" strokeWidth="0.9" />
        <XAxisTicks ticks={[0, 2, 4, 6, 8, 10]} scale={X} at={AXIS_Y} />

        {/* The mean, and the signed gaps that cancel against it. */}
        <line x1={X(m)} y1={barTop - 6} x2={X(m)} y2={AXIS_Y + 12} stroke="var(--accent-green)" strokeWidth="1.2" strokeDasharray="4 3" />
        <text x={X(m) + 5} y={barTop - 8} fontSize="9" fill="var(--accent-green)">
          mean {fmt(m)}
        </text>

        {values.map((value, i) => {
          const top = AXIS_Y - barHeight(squares[i])
          return squared ? (
            <rect
              key={`b${i}`}
              x={X(value) - 7}
              y={top}
              width="14"
              height={AXIS_Y - top}
              fill="var(--accent-blue)"
              fillOpacity="0.16"
              stroke="var(--accent-blue)"
              strokeWidth="0.8"
            />
          ) : (
            <line key={`g${i}`} x1={X(m)} y1={AXIS_Y + 8} x2={X(value)} y2={AXIS_Y + 8} stroke="var(--accent-teal)" strokeWidth="2" />
          )
        })}

        {values.map((value, i) => (
          <circle
            key={`p${i}`}
            cx={X(value)}
            cy={AXIS_Y}
            r={i === safeIndex ? 5 : 4}
            fill={i === safeIndex ? 'var(--accent-teal)' : 'var(--accent-blue)'}
          />
        ))}

        {squared && (
          <text x={PAD + 2} y={barTop - 2} fontSize="9" fill="var(--text-muted)">
            squared gap
          </text>
        )}
        {squared && !hasSpread && (
          <text x={PAD + 2} y={barTop + 12} fontSize="9" fill="var(--text-muted)">
            no spread: every value equals the mean
          </text>
        )}
      </svg>

      <div className="mt-3 space-y-2">
        <div className="flex flex-wrap items-center gap-2">
          <button
            type="button"
            onClick={addPoint}
            disabled={values.length >= 8}
            className="border border-mathua-border px-2 py-0.5 font-mono text-[10px] text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue disabled:opacity-40"
          >
            Add point
          </button>
          <button
            type="button"
            onClick={removePoint}
            disabled={values.length <= 2}
            className="border border-mathua-border px-2 py-0.5 font-mono text-[10px] text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue disabled:opacity-40"
          >
            Remove point
          </button>
          <span className="font-mono text-[10px] text-mathua-muted">{values.length} points</span>
        </div>

        <Select
          label="Point"
          value={String(safeIndex)}
          options={values.map((_, i) => ({ value: String(i), label: `P${i + 1}` }))}
          onChange={(value) => setSelected(Number(value))}
        />
        <NumberField label="x" value={values[safeIndex]} min={MIN_X} max={MAX_X} onChange={(value) => setPoint(safeIndex, value)} />
        <div className="flex flex-wrap gap-4">
          <Toggle label="Squared gaps" checked={squared} onChange={setSquared} />
          <Toggle label="Sample (n − 1)" checked={sample} onChange={setSample} />
        </div>
      </div>
    </FigureFrame>
  )
}
