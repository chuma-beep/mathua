'use client'

import { useMemo, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Slider } from '../Controls'
import { XAxisTicks } from '../Axis'
import { inferenceSchema } from '../schemas'
import {
  densityCurve,
  distributionPdf,
  inferenceModel,
  type InferenceDistribution,
} from '../math/inference'
import { scaleLinear } from '../scale'
import { fmt } from '../format'
import type { TemplateProps } from '../types'

const W = 360
const H = 330
const PAD = 32
const XMIN = -4
const XMAX = 6
const PLOT_TOP = 26
const PLOT_BOTTOM = 214

type Point = [number, number]

function areaPath(points: Point[], baseline: number): string {
  if (points.length === 0) return ''
  const top = points.map(([x, y], index) => `${index === 0 ? 'M' : 'L'}${x.toFixed(2)},${y.toFixed(2)}`).join(' ')
  const last = points[points.length - 1]
  const first = points[0]
  return `${top} L${last[0].toFixed(2)},${baseline.toFixed(2)} L${first[0].toFixed(2)},${baseline.toFixed(2)} Z`
}

export default function InferenceLab({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => inferenceSchema.parse(params), [params])
  const [dist, setDist] = useState<InferenceDistribution>(defaults.dist)
  const [df, setDf] = useState(defaults.df)
  const [alpha, setAlpha] = useState(defaults.alpha)
  const [effect, setEffect] = useState(defaults.effect)
  const [observed, setObserved] = useState(defaults.observed)

  const model = useMemo(
    () => inferenceModel({ dist, df, alpha, effect, observed }),
    [dist, df, alpha, effect, observed],
  )
  const X = scaleLinear(XMIN, XMAX, PAD, W - PAD)
  const highest = Math.max(
    distributionPdf(dist, df, 0),
    distributionPdf(dist, df, effect),
    distributionPdf(dist, df, observed),
  )
  const Y = scaleLinear(0, highest * 1.08, PLOT_BOTTOM, PLOT_TOP)
  const ticks = [-4, -2, 0, 2, 4, 6]

  const curve = (center: number): Point[] =>
    densityCurve(dist, df, XMIN - center, XMAX - center, 160).map(([x, density]) => [
      X(x + center),
      Y(density),
    ])

  const nullCurve = curve(0)
  const altCurve = curve(effect)
  const critical = model.critical
  const leftTail = Number.isFinite(critical)
    ? densityCurve(dist, df, XMIN, -critical, 90).map(([x, density]) => [X(x), Y(density)] as Point)
    : []
  const rightTail = Number.isFinite(critical)
    ? densityCurve(dist, df, critical, XMAX, 130).map(([x, density]) => [X(x), Y(density)] as Point)
    : []

  const reset = () => {
    setDist(defaults.dist)
    setDf(defaults.df)
    setAlpha(defaults.alpha)
    setEffect(defaults.effect)
    setObserved(defaults.observed)
  }

  const decision = model.reject ? 'reject the null' : 'fail to reject'
  const intervalText = Number.isFinite(model.interval[0])
    ? `[${fmt(model.interval[0], 2)}, ${fmt(model.interval[1], 2)}]`
    : '—'

  return (
    <FigureFrame
      title={title ?? 'A test is one distribution and two tails'}
      hint="Move the observed statistic and the true effect. The shaded tails are alpha, the vertical line is the observation, and the interval below uses the same critical value."
      caption={
        caption ??
        'Under the null the statistic follows the blue distribution. The rejection region is the two tails whose total area is alpha. The observed statistic fixes the p-value, the decision and the interval. The orange curve is the alternative, so its mass inside the rejection region is the power.'
      }
      onReset={reset}
      readout={[
        { label: 'null distribution', value: dist === 'z' ? 'standard normal' : `t with ${df} df` },
        { label: 'alpha', value: fmt(alpha, 3) },
        { label: 'critical value', value: fmt(critical, 3) },
        { label: 'true effect', value: fmt(effect, 2) },
        { label: 'observed statistic', value: fmt(model.statistic, 3) },
        { label: 'p-value', value: fmt(model.pValue, 4) },
        { label: 'decision', value: decision },
        { label: 'power', value: fmt(model.power, 3) },
        { label: 'Type II error', value: fmt(model.typeTwo, 3) },
        { label: 'interval', value: intervalText },
        ...(model.note ? [{ label: 'note', value: model.note }] : []),
      ]}
    >
      <svg
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        role="img"
        aria-label={`Two-sided ${dist} test with alpha ${fmt(alpha, 2)} and observed statistic ${fmt(observed, 2)}`}
      >
        <line x1={PAD} x2={W - PAD} y1={PLOT_BOTTOM} y2={PLOT_BOTTOM} stroke="var(--text-muted)" strokeWidth="0.8" />

        <path d={areaPath(leftTail, PLOT_BOTTOM)} fill="var(--accent-red)" fillOpacity="0.22" />
        <path d={areaPath(rightTail, PLOT_BOTTOM)} fill="var(--accent-red)" fillOpacity="0.22" />

        <path
          d={nullCurve.map(([x, y], index) => `${index === 0 ? 'M' : 'L'}${x.toFixed(2)},${y.toFixed(2)}`).join(' ')}
          fill="none"
          stroke="var(--accent-blue)"
          strokeWidth="2"
        />
        <path
          d={altCurve.map(([x, y], index) => `${index === 0 ? 'M' : 'L'}${x.toFixed(2)},${y.toFixed(2)}`).join(' ')}
          fill="none"
          stroke="var(--chart-3)"
          strokeWidth="1.6"
          strokeDasharray="5 3"
        />

        <line x1={X(-critical)} x2={X(-critical)} y1={PLOT_TOP} y2={PLOT_BOTTOM} stroke="var(--accent-red)" strokeWidth="1" strokeDasharray="3 3" />
        <line x1={X(critical)} x2={X(critical)} y1={PLOT_TOP} y2={PLOT_BOTTOM} stroke="var(--accent-red)" strokeWidth="1" strokeDasharray="3 3" />
        <line x1={X(observed)} x2={X(observed)} y1={PLOT_TOP - 2} y2={PLOT_BOTTOM + 4} stroke="var(--text-primary)" strokeWidth="1.8" />

        <text x={X(critical) + 4} y={PLOT_TOP + 10} fontSize="8" fill="var(--accent-red)">
          {fmt(critical, 2)}
        </text>
        <text x={X(observed) + 4} y={PLOT_BOTTOM - 6} fontSize="8" fill="var(--text-primary)">
          {dist} = {fmt(observed, 2)}
        </text>
        <text x={PAD} y={PLOT_TOP + 10} fontSize="9" fill="var(--text-muted)">
          null
        </text>
        <text x={PAD + 38} y={PLOT_TOP + 10} fontSize="9" fill="var(--chart-3)">
          alternative
        </text>

        <text x={PAD} y={250} fontSize="9" fill="var(--text-muted)">
          interval from the same critical value
        </text>
        <line x1={X(model.interval[0])} x2={X(model.interval[1])} y1={274} y2={274} stroke="var(--accent-teal)" strokeWidth="4" strokeLinecap="round" />
        <line x1={X(model.interval[0])} x2={X(model.interval[0])} y1={268} y2={280} stroke="var(--accent-teal)" strokeWidth="1.5" />
        <line x1={X(model.interval[1])} x2={X(model.interval[1])} y1={268} y2={280} stroke="var(--accent-teal)" strokeWidth="1.5" />
        <circle cx={X(observed)} cy={274} r="3.5" fill="var(--text-primary)" />
        <text x={X(observed)} y={292} textAnchor="middle" fontSize="8" fill="var(--text-muted)">
          {fmt(observed, 2)}
        </text>

        <XAxisTicks ticks={ticks} scale={X} at={PLOT_BOTTOM} />
      </svg>

      <div className="mt-3 grid grid-cols-1 gap-2">
        <div className="flex items-center gap-2">
          <label htmlFor="inference-distribution" className="w-24 shrink-0 font-mono text-[10px] text-mathua-muted">
            Distribution
          </label>
          <select
            id="inference-distribution"
            aria-label="Distribution"
            value={dist}
            onChange={(event) => setDist(event.target.value as InferenceDistribution)}
            className="h-7 min-w-0 flex-1 border border-mathua-border bg-transparent px-1 font-mono text-[10px] text-mathua-secondary"
          >
            <option value="z">standard normal</option>
            <option value="t">Student t</option>
          </select>
        </div>
        {dist === 't' ? <Slider label="Degrees of freedom" value={df} min={1} max={60} step={1} onChange={setDf} /> : null}
        <Slider label="Alpha" value={alpha} min={0.01} max={0.2} step={0.01} onChange={setAlpha} />
        <Slider label="True effect" value={effect} min={-1} max={3} step={0.1} onChange={setEffect} />
        <Slider label="Observed statistic" value={observed} min={-4} max={6} step={0.1} onChange={setObserved} />
      </div>
    </FigureFrame>
  )
}
