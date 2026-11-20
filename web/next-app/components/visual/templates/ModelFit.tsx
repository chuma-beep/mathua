'use client'

import { useMemo, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Slider } from '../Controls'
import { XAxisTicks, YAxisTicks } from '../Axis'
import { modelFitSchema } from '../schemas'
import { REGRESSION_DATA, logisticProbability, predict, trainRegression } from '../math/ml'
import { niceTicks, scaleLinear } from '../scale'
import { fmt } from '../format'
import type { TemplateProps } from '../types'

const W = 360
const H = 320
const PAD = 34

export default function ModelFit({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => modelFitSchema.parse(params), [params])
  const [mode, setMode] = useState(defaults.mode)
  const [rate, setRate] = useState(defaults.rate)
  const [steps, setSteps] = useState(defaults.steps)
  const [threshold, setThreshold] = useState(defaults.threshold)

  const fit = useMemo(() => trainRegression(rate, steps), [rate, steps])
  const last = fit.steps[fit.steps.length - 1]
  const weight = 2
  const probability = logisticProbability(weight, threshold, 1)
  const X = scaleLinear(0, 5, PAD, W - PAD)
  const Y = scaleLinear(0, 7, H - PAD, PAD)
  const LX = scaleLinear(-3, 3, PAD, W - PAD)
  const LY = scaleLinear(0, 1, H - PAD, PAD)
  const ticks = niceTicks(0, 6, 4).filter((value) => Math.abs(value) > 1e-9)

  const reset = () => {
    setMode(defaults.mode)
    setRate(defaults.rate)
    setSteps(defaults.steps)
    setThreshold(defaults.threshold)
  }

  return (
    <FigureFrame
      title={title ?? 'Fit a line by decreasing its loss'}
      hint="Each step uses the exact gradient of the mean squared error. The line and the loss are the values after the chosen number of steps."
      caption={
        caption ??
        'The regression mode starts at slope 0 and intercept 0 and records every gradient step. The logistic mode applies the sigmoid to w x + b, so its output is a probability rather than a class label.'
      }
      onReset={reset}
      readout={[
        { label: 'model', value: mode === 'regression' ? 'linear regression' : 'logistic regression' },
        {
          label: mode === 'regression' ? 'line' : 'score',
          value: mode === 'regression' && last ? `y = ${fmt(last.slope, 3)}x + ${fmt(last.intercept, 3)}` : `2x + ${fmt(threshold, 2)}`,
        },
        {
          label: mode === 'regression' ? 'MSE' : 'probability at x = 1',
          value: mode === 'regression' && last ? fmt(last.loss, 4) : fmt(probability, 4),
        },
        { label: 'steps', value: mode === 'regression' ? String(fit.steps.length) : '—' },
        { label: 'gradient', value: mode === 'regression' && last ? `(${fmt(last.gradSlope, 3)}, ${fmt(last.gradIntercept, 3)})` : '—' },
        { label: 'status', value: mode === 'regression' ? fit.status : probability > 0.5 ? 'positive class' : 'negative class' },
      ]}
    >
      <svg
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        role="img"
        aria-label={mode === 'regression' ? 'Least squares regression line and data' : 'Logistic probability curve'}
      >
        {mode === 'regression' ? (
          <>
            <line x1={PAD} y1={Y(0)} x2={W - PAD} y2={Y(0)} stroke="var(--text-muted)" strokeWidth="0.8" />
            <line x1={X(0)} y1={PAD} x2={X(0)} y2={H - PAD} stroke="var(--text-muted)" strokeWidth="0.8" />
            {REGRESSION_DATA.map((point) => (
              <circle key={point.x} cx={X(point.x)} cy={Y(point.y)} r={4} fill="var(--chart-3)" />
            ))}
            {last ? (
              <line
                x1={X(0)}
                y1={Y(predict(last.slope, last.intercept, 0))}
                x2={X(5)}
                y2={Y(predict(last.slope, last.intercept, 5))}
                stroke="var(--accent-blue)"
                strokeWidth={1.8}
              />
            ) : null}
            <XAxisTicks ticks={ticks} scale={X} at={Y(0)} />
            <YAxisTicks ticks={ticks} scale={Y} at={X(0)} />
          </>
        ) : (
          <>
            <line x1={PAD} y1={LY(0)} x2={W - PAD} y2={LY(0)} stroke="var(--text-muted)" strokeWidth="0.8" />
            <line x1={LX(0)} y1={PAD} x2={LX(0)} y2={H - PAD} stroke="var(--text-muted)" strokeWidth="0.8" />
            <polyline
              points={Array.from({ length: 121 }, (_v, index) => {
                const x = -3 + (6 * index) / 120
                return `${LX(x)},${LY(logisticProbability(weight, threshold, x))}`
              }).join(' ')}
              fill="none"
              stroke="var(--accent-blue)"
              strokeWidth={1.8}
            />
            <line x1={LX(-3)} y1={LY(0.5)} x2={LX(3)} y2={LY(0.5)} stroke="var(--text-muted)" strokeDasharray="4 3" />
            <circle cx={LX(1)} cy={LY(probability)} r={4} fill="var(--accent-teal)" />
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
          <option value="regression">Linear regression</option>
          <option value="logistic">Logistic regression</option>
        </select>
        {mode === 'regression' ? (
          <div className="grid grid-cols-1 gap-2">
            <Slider label="Learning rate" value={rate} min={0.01} max={0.5} step={0.01} onChange={setRate} />
            <Slider label="Steps" value={steps} min={1} max={20} step={1} onChange={setSteps} />
          </div>
        ) : (
          <Slider label="Bias b" value={threshold} min={-2} max={2} step={0.1} onChange={setThreshold} />
        )}
      </div>
    </FigureFrame>
  )
}
