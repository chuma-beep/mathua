'use client'

import { useId, useMemo, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Select, Slider } from '../Controls'
import { XAxisTicks, YAxisTicks } from '../Axis'
import { taylorSchema } from '../schemas'
import { fnValue, maxError, taylorCoeffs, taylorValue, type SeriesFn } from '../math/series'
import { niceTicks, pathFromPoints, scaleLinear } from '../scale'
import { fmt } from '../format'
import type { TemplateProps } from '../types'

const W = 340
const H = 240
const PAD = 28

const FN_OPTIONS: Array<{ value: SeriesFn; label: string }> = [
  { value: 'sin', label: 'sin x' },
  { value: 'cos', label: 'cos x' },
  { value: 'exp', label: 'e^x' },
  { value: 'geometric', label: '1 / (1 - x)' },
]

export default function TaylorSeries({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => taylorSchema.parse(params), [params])
  const [fn, setFn] = useState<SeriesFn>(defaults.fn)
  const [center, setCenter] = useState(defaults.center)
  const [degree, setDegree] = useState(defaults.degree)
  const uid = useId().replace(/:/g, '')
  const window = defaults.window

  const coeffs = useMemo(() => taylorCoeffs(fn, center, degree), [fn, center, degree])
  const samples = useMemo(() => {
    // One path per continuous piece. Joining across an undefined sample would draw
    // a straight line through the pole of 1/(1-x), which is the opposite of what a
    // pole looks like.
    const curve: Array<Array<[number, number]>> = []
    const poly: Array<Array<[number, number]>> = []
    let currentCurve: Array<[number, number]> = []
    let currentPoly: Array<[number, number]> = []
    for (let i = 0; i <= 240; i++) {
      const x = window[0] + ((window[1] - window[0]) * i) / 240
      const f = fnValue(fn, x)
      const p = taylorValue(coeffs, center, x)
      if (Number.isFinite(f) && Math.abs(f) < 40) currentCurve.push([x, f])
      else if (currentCurve.length > 0) {
        curve.push(currentCurve)
        currentCurve = []
      }
      if (Number.isFinite(p) && Math.abs(p) < 40) currentPoly.push([x, p])
      else if (currentPoly.length > 0) {
        poly.push(currentPoly)
        currentPoly = []
      }
    }
    if (currentCurve.length > 0) curve.push(currentCurve)
    if (currentPoly.length > 0) poly.push(currentPoly)
    return { curve, poly }
  }, [fn, coeffs, center, window])

  const yValues = [...samples.curve.flat(), ...samples.poly.flat()].map(([, y]) => y).concat([0])
  const loose = Math.max(3, Math.min(12, Math.max(...yValues.map(Math.abs))))
  const yMin = -loose
  const yMax = loose
  const X = scaleLinear(window[0], window[1], PAD, W - PAD)
  const Y = scaleLinear(yMin, yMax, H - PAD, PAD)

  const error = maxError(fn, center, degree, window[0], window[1])
  const leading = coeffs[coeffs.length - 1]
  // The centre can be dragged onto the pole at x = 1, where there is no Taylor
  // polynomial at all. Say so instead of printing a number that came from an
  // undefined coefficient.
  const singular = !Number.isFinite(fnValue(fn, center))

  const reset = () => {
    setFn(defaults.fn)
    setCenter(defaults.center)
    setDegree(defaults.degree)
  }

  return (
    <FigureFrame
      title={title ?? 'Taylor polynomials of increasing degree'}
      hint="Raise the degree and move the centre: the polynomial matches the curve at the centre."
      caption={
        caption ??
        'Each extra term adds one derivative of agreement at the centre. The polynomial is exact there and drifts away from the curve as x moves, and the drift is what the error readout measures.'
      }
      onReset={reset}
      readout={[
        { label: 'degree', value: String(degree) },
        { label: 'centre', value: fmt(center) },
        { label: 'max error', value: singular || Number.isNaN(error) ? 'undefined' : fmt(error, 4) },
        { label: 'error at right', value: fmt(Math.abs(fnValue(fn, window[1]) - taylorValue(coeffs, center, window[1])), 4) },
        { label: 'leading coeff', value: fmt(leading, 4) },
      ]}
    >
      <svg
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        role="img"
        aria-label="A curve and its Taylor polynomial at increasing degree"
      >
        <line x1={PAD} y1={Y(0)} x2={W - PAD} y2={Y(0)} stroke="var(--text-muted)" strokeWidth="0.8" />
        {window[0] < 0 && window[1] > 0 && (
          <line x1={X(0)} y1={PAD} x2={X(0)} y2={H - PAD} stroke="var(--text-muted)" strokeWidth="0.8" />
        )}

        {/* The polynomial runs far outside the plot once the degree is low and the
            window is wide, so it is clipped rather than allowed to cross the hint. */}
        <defs>
          <clipPath id={`${uid}-plot`}>
            <rect x={PAD} y={PAD} width={W - PAD * 2} height={H - PAD * 2} />
          </clipPath>
        </defs>
        <g clipPath={`url(#${uid}-plot)`}>
          {samples.curve.map((segment, i) => (
            <path
              key={`f${i}`}
              d={pathFromPoints(segment.map(([x, y]) => [X(x), Y(y)]))}
              fill="none"
              stroke="var(--text-primary)"
              strokeWidth="1.8"
            />
          ))}
          {samples.poly.map((segment, i) => (
            <path
              key={`p${i}`}
              d={pathFromPoints(segment.map(([x, y]) => [X(x), Y(y)]))}
              fill="none"
              stroke="var(--accent-blue)"
              strokeWidth="1.6"
              strokeDasharray="6 3"
            />
          ))}
        </g>

        <line x1={X(center)} y1={PAD} x2={X(center)} y2={H - PAD} stroke="var(--accent-green)" strokeWidth="1" strokeDasharray="3 3" />
        {!singular && <circle cx={X(center)} cy={Y(fnValue(fn, center))} r="4" fill="var(--accent-green)" />}
        {singular && (
          <text x={PAD + 34} y={PAD + 10} fontSize="9" fill="var(--accent-red)">
            pole at the centre
          </text>
        )}

        <XAxisTicks ticks={niceTicks(window[0], window[1], 5)} scale={X} at={Y(0)} />
        <YAxisTicks ticks={niceTicks(yMin, yMax, 4).filter((v) => Math.abs(v) > 1e-9)} scale={Y} at={X(0)} />

        <text x={PAD + 2} y={PAD + 10} fontSize="10" fill="var(--text-primary)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
          f
        </text>
        <text x={PAD + 2} y={PAD + 22} fontSize="10" fill="var(--accent-blue)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
          P{degree}
        </text>
      </svg>

      <div className="mt-3 space-y-2">
        <Select label="Function" value={fn} options={FN_OPTIONS} onChange={(v) => setFn(v as SeriesFn)} />
        <Slider label="degree" value={degree} min={0} max={14} step={1} onChange={(v) => setDegree(Math.round(v))} />
        <Slider label="centre" value={center} min={-2} max={2} step={0.1} onChange={setCenter} />
      </div>
    </FigureFrame>
  )
}
