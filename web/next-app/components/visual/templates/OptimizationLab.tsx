'use client'

import { useMemo, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Slider } from '../Controls'
import { XAxisTicks, YAxisTicks } from '../Axis'
import { optimizationSchema } from '../schemas'
import { gradientDescent, linearProgram, objective } from '../math/optimization'
import { contourSegments } from '../math/contour'
import { niceTicks, scaleLinear } from '../scale'
import { fmt } from '../format'
import type { TemplateProps } from '../types'

const W = 360
const H = 320
const PAD = 32
const BOUNDS = { x0: -3, x1: 3, y0: -3, y1: 3 }
const LEVELS = [0.5, 1, 2, 4, 8, 16]

export default function OptimizationLab({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => optimizationSchema.parse(params), [params])
  const [mode, setMode] = useState(defaults.mode)
  const [startX, setStartX] = useState(defaults.startX)
  const [startY, setStartY] = useState(defaults.startY)
  const [rate, setRate] = useState(defaults.rate)
  const [steps, setSteps] = useState(defaults.steps)
  const [candidateX, setCandidateX] = useState(defaults.candidateX)
  const [candidateY, setCandidateY] = useState(defaults.candidateY)

  const run = useMemo(() => gradientDescent(startX, startY, rate, steps), [startX, startY, rate, steps])
  const lp = useMemo(() => linearProgram(candidateX, candidateY), [candidateX, candidateY])
  const X = scaleLinear(BOUNDS.x0, BOUNDS.x1, PAD, W - PAD)
  const Y = scaleLinear(BOUNDS.y0, BOUNDS.y1, H - PAD, PAD)
  const ticks = niceTicks(-3, 3, 4).filter((value) => Math.abs(value) > 1e-9)

  const reset = () => {
    setMode(defaults.mode)
    setStartX(defaults.startX)
    setStartY(defaults.startY)
    setRate(defaults.rate)
    setSteps(defaults.steps)
    setCandidateX(defaults.candidateX)
    setCandidateY(defaults.candidateY)
  }

  const contours = useMemo(
    () =>
      mode === 'descent'
        ? LEVELS.flatMap((level) => contourSegments((x, y) => objective(x, y), BOUNDS, level, 60, 60))
        : [],
    [mode],
  )

  const path = run.points.map((point) => `${X(point.x)},${Y(point.y)}`).join(' ')
  const last = run.points[run.points.length - 1]

  return (
    <FigureFrame
      title={title ?? 'Follow the gradient, or test a constraint'}
      hint="Change the step size and count. The path is the sequence of points the iteration computes."
      caption={
        caption ??
        'The contour plot is the actual objective x² + 2y². Each dot is one gradient step. In the feasible-region mode every constraint is evaluated at the candidate point.'
      }
      onReset={reset}
      readout={[
        {
          label: 'objective',
          value: mode === 'descent' ? 'x² + 2y²' : 'maximise 3x + 2y',
        },
        {
          label: 'start',
          value: mode === 'descent' ? `(${fmt(startX)}, ${fmt(startY)})` : `(${fmt(candidateX)}, ${fmt(candidateY)})`,
        },
        {
          label: 'gradient',
          value: mode === 'descent' && last ? `(${fmt(last.gradX)}, ${fmt(last.gradY)})` : '—',
        },
        {
          label: 'value',
          value: mode === 'descent' && last ? fmt(last.value, 3) : fmt(lp.objectiveValue, 3),
        },
        { label: 'steps', value: mode === 'descent' ? String(run.points.length) : String(lp.vertices.length) },
        { label: 'status', value: mode === 'descent' ? run.status : lp.feasible ? 'feasible' : 'violates a constraint' },
        {
          label: 'optimum',
          value: `(${fmt(lp.optimum.x)}, ${fmt(lp.optimum.y)}) = ${fmt(lp.optimum.value)}`,
        },
        {
          label: 'constraints',
          value: lp.candidate.map((item) => `${item.label}: ${fmt(item.value, 2)}`).join('; '),
        },
      ]}
    >
      <svg
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        role="img"
        aria-label={mode === 'descent' ? 'Gradient descent path on contour lines' : 'Feasible region and candidate point'}
      >
        <line x1={PAD} y1={Y(0)} x2={W - PAD} y2={Y(0)} stroke="var(--text-muted)" strokeWidth="0.8" />
        <line x1={X(0)} y1={PAD} x2={X(0)} y2={H - PAD} stroke="var(--text-muted)" strokeWidth="0.8" />

        {mode === 'descent' ? (
          <>
            {contours.map((segment, index) => (
              <line
                key={index}
                x1={X(segment[0][0])}
                y1={Y(segment[0][1])}
                x2={X(segment[1][0])}
                y2={Y(segment[1][1])}
                stroke="var(--chart-3)"
                strokeWidth={1}
              />
            ))}
            <polyline points={path} fill="none" stroke="var(--accent-blue)" strokeWidth={1.8} />
            {run.points.map((point) => (
              <circle key={point.index} cx={X(point.x)} cy={Y(point.y)} r={2.6} fill="var(--accent-blue)" />
            ))}
            <circle cx={X(startX)} cy={Y(startY)} r={5} fill="none" stroke="var(--accent-red)" strokeWidth={1.4} />
          </>
        ) : (
          <>
            <polygon
              points={lp.vertices.map((vertex) => `${X(vertex.x)},${Y(vertex.y)}`).join(' ')}
              fill="var(--accent-blue)"
              fillOpacity={0.12}
              stroke="var(--accent-blue)"
              strokeWidth={1.5}
            />
            <line x1={X(0)} y1={Y(0)} x2={X(2.5)} y2={Y(0)} stroke="var(--chart-3)" strokeDasharray="4 3" />
            <circle cx={X(candidateX)} cy={Y(candidateY)} r={5} fill={lp.feasible ? 'var(--accent-teal)' : 'var(--accent-red)'} />
            <circle cx={X(lp.optimum.x)} cy={Y(lp.optimum.y)} r={3.5} fill="var(--accent-blue)" />
          </>
        )}

        <XAxisTicks ticks={ticks} scale={X} at={Y(0)} />
        <YAxisTicks ticks={ticks} scale={Y} at={X(0)} />
      </svg>

      <div className="mt-3 space-y-2">
        <select
          aria-label="Model"
          value={mode}
          onChange={(event) => setMode(event.target.value as typeof mode)}
          className="h-8 w-full border border-mathua-border bg-transparent px-2 font-mono text-[11px]"
        >
          <option value="descent">Gradient descent</option>
          <option value="feasible">Feasible region</option>
        </select>
        {mode === 'descent' ? (
          <div className="grid grid-cols-1 gap-2">
            <Slider label="Start x" value={startX} min={-3} max={3} step={0.1} onChange={setStartX} />
            <Slider label="Start y" value={startY} min={-3} max={3} step={0.1} onChange={setStartY} />
            <Slider label="Learning rate" value={rate} min={0.01} max={0.5} step={0.01} onChange={setRate} />
            <Slider label="Steps" value={steps} min={1} max={20} step={1} onChange={setSteps} />
          </div>
        ) : (
          <div className="grid grid-cols-1 gap-2">
            <Slider label="Candidate x" value={candidateX} min={0} max={3} step={0.1} onChange={setCandidateX} />
            <Slider label="Candidate y" value={candidateY} min={0} max={3} step={0.1} onChange={setCandidateY} />
          </div>
        )}
      </div>
    </FigureFrame>
  )
}
