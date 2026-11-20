'use client'

import { useMemo, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Slider } from '../Controls'
import { XAxisTicks, YAxisTicks } from '../Axis'
import { analysisSchema } from '../schemas'
import {
  SEQUENCES,
  metricBall,
  partialSums,
  sequenceView,
  type MetricId,
  type SequenceId,
} from '../math/analysis'
import { niceTicks, scaleLinear } from '../scale'
import { fmt } from '../format'
import type { TemplateProps } from '../types'

const W = 360
const H = 320
const PAD = 34

export default function AnalysisLab({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => analysisSchema.parse(params), [params])
  const [mode, setMode] = useState(defaults.mode)
  const [sequence, setSequence] = useState<SequenceId>(defaults.sequence)
  const [epsilon, setEpsilon] = useState(defaults.epsilon)
  const [metric, setMetric] = useState<MetricId>(defaults.metric)
  const [radius, setRadius] = useState(defaults.radius)

  const view = useMemo(() => sequenceView(sequence, epsilon, 30), [sequence, epsilon])
  const sums = useMemo(() => partialSums(sequence, 30), [sequence])
  const ball = useMemo(() => metricBall(metric, radius), [metric, radius])
  const spec = SEQUENCES[sequence]

  const X = scaleLinear(1, 30, PAD, W - PAD)
  const Y = scaleLinear(-1.2, 2.6, H - PAD, PAD)
  const MX = scaleLinear(-2.5, 2.5, PAD, W - PAD)
  const MY = scaleLinear(-2.5, 2.5, H - PAD, PAD)
  const ticks = niceTicks(-2, 2, 4).filter((value) => Math.abs(value) > 1e-9)
  const series = sums.map((value, index) => `${X(index + 1)},${Y(Math.max(-1.2, Math.min(2.6, value)))} `).join(' ')
  const terms = view.terms.map((value, index) => `${X(index + 1)},${Y(value)}`).join(' ')

  const reset = () => {
    setMode(defaults.mode)
    setSequence(defaults.sequence)
    setEpsilon(defaults.epsilon)
    setMetric(defaults.metric)
    setRadius(defaults.radius)
  }

  return (
    <FigureFrame
      title={title ?? 'Make the epsilon explicit'}
      hint="The band is the actual epsilon-neighbourhood of the limit. The threshold is the first index after which every later term stays inside it."
      caption={
        caption ??
        'A sequence converges when every epsilon band eventually contains the whole tail. The metric mode changes the distance formula itself, so the ball is an L1, L2 or L∞ ball rather than the same picture relabelled.'
      }
      onReset={reset}
      readout={[
        { label: 'mode', value: mode },
        { label: 'object', value: mode === 'metric' ? `${metric} ball` : `${spec.formula}` },
        { label: 'limit', value: mode === 'metric' ? '—' : fmt(spec.limit, 4) },
        { label: 'epsilon', value: mode === 'metric' ? fmt(radius, 2) : fmt(epsilon, 2) },
        {
          label: 'threshold',
          value: mode === 'metric' ? String(ball.pointsInside) : view.threshold === null ? 'not reached in 30 terms' : String(view.threshold),
        },
        {
          label: 'tail distance',
          value: mode === 'metric' ? '—' : view.tailDistance === null ? '—' : fmt(view.tailDistance, 4),
        },
        {
          label: 'inside',
          value: mode === 'metric' ? `${ball.pointsInside} lattice points` : view.within ? 'the last term is inside' : 'the last term is outside',
        },
        { label: 'series tail', value: mode === 'sequence' ? fmt(sums[sums.length - 1], 4) : '—' },
      ]}
    >
      <svg
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        role="img"
        aria-label={mode === 'metric' ? `${metric} ball of radius ${fmt(radius, 2)}` : `${spec.label} with an epsilon band`}
      >
        {mode === 'metric' ? (
          <>
            <line x1={PAD} y1={MY(0)} x2={W - PAD} y2={MY(0)} stroke="var(--text-muted)" strokeWidth="0.8" />
            <line x1={MX(0)} y1={PAD} x2={MX(0)} y2={H - PAD} stroke="var(--text-muted)" strokeWidth="0.8" />
            {ball.points.map((point) => (
              <circle
                key={`${point.x},${point.y}`}
                cx={MX(point.x)}
                cy={MY(point.y)}
                r={point.inside ? 5 : 3}
                fill={point.inside ? 'var(--accent-blue)' : 'var(--text-muted)'}
                fillOpacity={point.inside ? 0.85 : 0.35}
              />
            ))}
            <XAxisTicks ticks={ticks} scale={MX} at={MY(0)} />
            <YAxisTicks ticks={ticks} scale={MY} at={MX(0)} />
          </>
        ) : (
          <>
            <rect
              x={PAD}
              y={Y(spec.limit + epsilon)}
              width={W - PAD * 2}
              height={Math.max(1, Y(spec.limit - epsilon) - Y(spec.limit + epsilon))}
              fill="var(--accent-blue)"
              fillOpacity={0.12}
            />
            <line x1={PAD} y1={Y(spec.limit)} x2={W - PAD} y2={Y(spec.limit)} stroke="var(--accent-blue)" strokeWidth={1.2} />
            <polyline points={terms} fill="none" stroke="var(--chart-3)" strokeWidth={1.4} />
            {view.terms.map((value, index) => (
              <circle key={index} cx={X(index + 1)} cy={Y(value)} r={2.5} fill="var(--accent-blue)" />
            ))}
            {mode === 'series' ? <polyline points={series} fill="none" stroke="var(--accent-teal)" strokeWidth={1.6} /> : null}
            {view.threshold !== null ? (
              <line x1={X(view.threshold)} y1={PAD} x2={X(view.threshold)} y2={H - PAD} stroke="var(--accent-red)" strokeDasharray="4 3" />
            ) : null}
            <XAxisTicks ticks={[1, 10, 20, 30]} scale={X} at={H - PAD} />
            <YAxisTicks ticks={[-1, 0, 1, 2]} scale={Y} at={PAD} />
          </>
        )}
      </svg>

      <div className="mt-3 space-y-2">
        <select
          aria-label="Model"
          value={mode}
          onChange={(event) => setMode(event.target.value as typeof mode)}
          className="h-8 w-full border border-mathua-border bg-transparent px-2 font-mono text-[11px]"
        >
          <option value="sequence">Sequence</option>
          <option value="series">Partial sums</option>
          <option value="metric">Metric ball</option>
        </select>
        {mode === 'metric' ? (
          <div className="grid grid-cols-1 gap-2">
            <select
              aria-label="Metric"
              value={metric}
              onChange={(event) => setMetric(event.target.value as MetricId)}
              className="h-8 w-full border border-mathua-border bg-transparent px-2 font-mono text-[11px]"
            >
              <option value="euclidean">Euclidean L2</option>
              <option value="manhattan">Manhattan L1</option>
              <option value="chebyshev">Chebyshev L∞</option>
            </select>
            <Slider label="Radius" value={radius} min={0.5} max={3} step={0.1} onChange={setRadius} />
          </div>
        ) : (
          <div className="grid grid-cols-1 gap-2">
            <select
              aria-label="Sequence"
              value={sequence}
              onChange={(event) => setSequence(event.target.value as SequenceId)}
              className="h-8 w-full border border-mathua-border bg-transparent px-2 font-mono text-[11px]"
            >
              {Object.entries(SEQUENCES).map(([id, item]) => (
                <option key={id} value={id}>
                  {item.label}
                </option>
              ))}
            </select>
            <Slider label="Epsilon" value={epsilon} min={0.01} max={1} step={0.01} onChange={setEpsilon} />
          </div>
        )}
      </div>
    </FigureFrame>
  )
}
