'use client'

import { useMemo, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Slider } from '../Controls'
import { XAxisTicks } from '../Axis'
import { accumulationSchema } from '../schemas'
import {
  ACC_LABELS,
  accumulationModel,
  areaCurve,
  integrand,
  type AccumulationFn,
} from '../math/accumulation'
import { scaleLinear } from '../scale'
import { fmt } from '../format'
import type { TemplateProps } from '../types'

const W = 360
const H = 350
const PAD = 34
const F_TOP = 44
const F_BOTTOM = 168
const A_TOP = 210
const A_BOTTOM = 320
const LO = 0
const HI = 3

export default function AccumulationLab({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => accumulationSchema.parse(params), [params])
  const [fn, setFn] = useState<AccumulationFn>(defaults.fn)
  const [a, setA] = useState(defaults.a)
  const [x, setX] = useState(defaults.x)

  const model = accumulationModel(fn, a, x)
  const left = Math.min(a, x)
  const right = Math.max(a, x)

  const X = scaleLinear(LO, HI, PAD, W - PAD)
  const curve = useMemo(() => {
    const out: Array<[number, number]> = []
    for (let index = 0; index <= 120; index++) {
      const value = LO + ((HI - LO) * index) / 120
      out.push([value, integrand(fn, value)])
    }
    return out
  }, [fn])
  const fMax = Math.max(0.4, ...curve.map(([, value]) => value))
  const FY = scaleLinear(0, fMax * 1.15, F_BOTTOM, F_TOP)

  const areaPoints = useMemo(
    () => curve.filter(([value]) => value >= left - 1e-9 && value <= right + 1e-9),
    [curve, left, right],
  )
  const areaPath = areaPoints.length
    ? `M ${X(areaPoints[0][0])} ${FY(0)} ` +
      areaPoints.map(([value, height]) => `L ${X(value)} ${FY(height)}`).join(' ') +
      ` L ${X(areaPoints[areaPoints.length - 1][0])} ${FY(0)} Z`
    : ''

  const A = useMemo(() => areaCurve(fn, a, LO, HI), [fn, a])
  const aMin = Math.min(0, ...A.map(([, value]) => value)) - 0.2
  const aMax = Math.max(0.2, ...A.map(([, value]) => value)) + 0.2
  const AY = scaleLinear(aMin, aMax, A_BOTTOM, A_TOP)

  const reset = () => {
    setFn(defaults.fn)
    setA(defaults.a)
    setX(defaults.x)
  }

  return (
    <FigureFrame
      title={title ?? 'The area function, and its derivative underneath'}
      hint="Move x. The lower curve is the area accumulated so far, and its slope at x is the height of the upper curve there. That equality is the fundamental theorem."
      caption={
        caption ??
        'The shaded area is A(x), the integral of f from a to x. The lower panel plots A against x. A tangent at x has slope f(x), so differentiating the area returns the integrand.'
      }
      onReset={reset}
      readout={[
        { label: 'function', value: ACC_LABELS[fn] },
        { label: 'left end a', value: fmt(a, 2) },
        { label: 'right end x', value: fmt(x, 2) },
        { label: 'area A(x)', value: fmt(model.area, 4) },
        { label: 'f(x)', value: fmt(model.value, 4) },
        { label: 'measured A′(x)', value: fmt(model.measured, 4) },
        { label: 'difference', value: fmt(model.error, 8) },
        { label: 'status', value: model.status === 'ok' ? 'theorem holds' : 'invalid input' },
      ]}
    >
      <svg
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        role="img"
        aria-label={`Area under ${ACC_LABELS[fn]} from ${fmt(a, 2)} to ${fmt(x, 2)} and the area function below it`}
      >
        <text x={PAD} y={F_TOP - 10} fontSize="9" fill="var(--text-muted)">
          f(x)
        </text>
        <line x1={PAD} x2={W - PAD} y1={FY(0)} y2={FY(0)} stroke="var(--text-muted)" strokeWidth="0.8" />
        <XAxisTicks ticks={[0, 1, 2, 3]} scale={X} at={A_BOTTOM} />
        <path d={areaPath} fill="var(--accent-blue)" fillOpacity="0.2" stroke="none" />
        <polyline
          points={curve.map(([value, height]) => `${X(value)},${FY(height)}`).join(' ')}
          fill="none"
          stroke="var(--accent-blue)"
          strokeWidth="1.8"
        />
        <line x1={X(a)} x2={X(a)} y1={F_TOP} y2={F_BOTTOM} stroke="var(--accent-teal)" strokeWidth="1" strokeDasharray="3 3" />
        <line x1={X(x)} x2={X(x)} y1={F_TOP} y2={F_BOTTOM} stroke="var(--text-primary)" strokeWidth="1.4" />
        <text x={X(a) + 3} y={F_TOP + 8} fontSize="8" fill="var(--accent-teal)">
          a
        </text>
        <circle cx={X(x)} cy={FY(model.value)} r="3" fill="var(--text-primary)" />

        <text x={PAD} y={A_TOP - 10} fontSize="9" fill="var(--text-muted)">
          A(x) and its tangent
        </text>
        <line x1={PAD} x2={W - PAD} y1={AY(0)} y2={AY(0)} stroke="var(--text-muted)" strokeWidth="0.8" />
        <polyline
          points={A.map(([value, height]) => `${X(value)},${AY(height)}`).join(' ')}
          fill="none"
          stroke="var(--accent-teal)"
          strokeWidth="1.8"
        />
        {model.status === 'ok' ? (
          <line
            x1={X(Math.max(LO, x - 0.45))}
            y1={AY(model.area - model.measured * (x - Math.max(LO, x - 0.45)))}
            x2={X(Math.min(HI, x + 0.45))}
            y2={AY(model.area + model.measured * (Math.min(HI, x + 0.45) - x))}
            stroke="var(--accent-red)"
            strokeWidth="1.2"
            strokeDasharray="4 3"
          />
        ) : null}
        <circle cx={X(x)} cy={AY(model.area)} r="3.4" fill="var(--accent-teal)" />
      </svg>

      <div className="mt-3 grid grid-cols-1 gap-2">
        <div className="flex items-center gap-2">
          <label htmlFor="accumulation-fn" className="w-20 shrink-0 font-mono text-[10px] text-mathua-muted">
            Function
          </label>
          <select
            id="accumulation-fn"
            aria-label="Function"
            value={fn}
            onChange={(event) => setFn(event.target.value as AccumulationFn)}
            className="h-7 min-w-0 flex-1 border border-mathua-border bg-transparent px-1 font-mono text-[10px] text-mathua-secondary"
          >
            {Object.entries(ACC_LABELS).map(([value, label]) => (
              <option key={value} value={value}>
                {label}
              </option>
            ))}
          </select>
        </div>
        <Slider label="Left end a" value={a} min={0} max={2} step={0.1} onChange={setA} />
        <Slider label="Right end x" value={x} min={0} max={3} step={0.1} onChange={setX} />
      </div>
    </FigureFrame>
  )
}
