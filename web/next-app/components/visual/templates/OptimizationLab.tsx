'use client'

import { useId, useMemo, useRef, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Slider } from '../Controls'
import { XAxisTicks, YAxisTicks } from '../Axis'
import { optimizationSchema } from '../schemas'
import {
  gradientDescent,
  linearProgram,
  objective,
  objectiveContour,
  type LpShiftWhich,
} from '../math/optimization'
import { clipLineToRect } from '../math/conic'
import { contourSegments } from '../math/contour'
import { invertLinear, niceTicks, scaleLinear } from '../scale'
import { useDrag } from '../useDrag'
import { fmt } from '../format'
import type { TemplateProps } from '../types'

const W = 360
const H = 320
const PAD = 32
const BOUNDS = { x0: -3, x1: 3, y0: -3, y1: 3 }
const LP_BOUNDS = { x0: 0, x1: 3, y0: 0, y1: 3 }
const LEVELS = [2, 4, 6, 8, 10, 12]

const gradientArrow = (x: number, y: number, scale: number): [number, number] => {
  const length = Math.hypot(x, y)
  if (length <= 1e-12) return [x, y]
  return [(x / length) * scale, (y / length) * scale]
}

export default function OptimizationLab({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => optimizationSchema.parse(params), [params])
  const [mode, setMode] = useState(defaults.mode)
  const [startX, setStartX] = useState(defaults.startX)
  const [startY, setStartY] = useState(defaults.startY)
  const [rate, setRate] = useState(defaults.rate)
  const [steps, setSteps] = useState(defaults.steps)
  const [candidateX, setCandidateX] = useState(defaults.candidateX)
  const [candidateY, setCandidateY] = useState(defaults.candidateY)
  const [rhs, setRhs] = useState<LpShiftWhich>(defaults.rhs)
  const [delta, setDelta] = useState(defaults.delta)
  const svgRef = useRef<SVGSVGElement | null>(null)
  const clipId = useId()
  const arrowId = useId()

  const run = useMemo(() => gradientDescent(startX, startY, rate, steps), [startX, startY, rate, steps])
  const lp = useMemo(
    () => linearProgram(candidateX, candidateY, { which: rhs, delta }),
    [candidateX, candidateY, rhs, delta],
  )
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
    setRhs(defaults.rhs)
    setDelta(defaults.delta)
  }

  const toX = (value: number) => invertLinear(BOUNDS.x0, BOUNDS.x1, PAD, W - PAD)(value)
  const toY = (value: number) => invertLinear(BOUNDS.y0, BOUNDS.y1, H - PAD, PAD)(value)

  useDrag(
    svgRef,
    (point) => {
      const x = Math.min(3, Math.max(0, toX(point.x)))
      const y = Math.min(3, Math.max(0, toY(point.y)))
      setCandidateX(x)
      setCandidateY(y)
    },
    mode === 'feasible',
  )

  const contours = useMemo(
    () =>
      mode === 'descent'
        ? LEVELS.flatMap((level) => contourSegments((x, y) => objective(x, y), BOUNDS, level, 60, 60))
        : [],
    [mode],
  )

  const path = run.points.map((point) => `${X(point.x)},${Y(point.y)}`).join(' ')
  const last = run.points[run.points.length - 1]

  const objectiveLines = useMemo(
    () =>
      mode === 'feasible'
        ? LEVELS.map((level) => {
            const contour = objectiveContour(level)
            const clipped = clipLineToRect(contour.point, contour.direction, LP_BOUNDS)
            return clipped
              ? {
                  level,
                  from: [X(clipped[0][0]), Y(clipped[0][1])] as [number, number],
                  to: [X(clipped[1][0]), Y(clipped[1][1])] as [number, number],
                }
              : null
          }).filter((item): item is NonNullable<typeof item> => item !== null)
        : [],
    [mode, X, Y],
  )

  const constraintLines = useMemo(
    () =>
      mode === 'feasible'
        ? lp.constraints.map((constraint) => {
            const [nx, ny] = constraint.normal
            const point: [number, number] =
              Math.abs(ny) > 1e-9
                ? [0, constraint.bound / ny]
                : [constraint.bound / nx, 0]
            const direction: [number, number] =
              Math.abs(ny) > 1e-9 ? [ny, -nx] : [0, 1]
            const clipped = clipLineToRect(point, direction, LP_BOUNDS)
            return clipped
              ? {
                  label: constraint.label,
                  active: lp.active.includes(constraint.label),
                  from: [X(clipped[0][0]), Y(clipped[0][1])] as [number, number],
                  to: [X(clipped[1][0]), Y(clipped[1][1])] as [number, number],
                }
              : null
          }).filter((item): item is NonNullable<typeof item> => item !== null)
        : [],
    [mode, lp, X, Y],
  )

  const optimum = lp.optimum
  const gradientEnd = gradientArrow(3, 2, 0.62)
  const normalArrows = lp.optimumKkt.multipliers.map((multiplier) => {
    const direction = gradientArrow(multiplier.normal[0], multiplier.normal[1], 0.45)
    return {
      label: multiplier.label,
      to: [X(optimum.x + direction[0]), Y(optimum.y + direction[1])] as [number, number],
    }
  })

  const activeText = lp.kkt.active.length ? lp.kkt.active.join(', ') : 'none'
  const multiplierText = lp.kkt.multipliers.length
    ? lp.kkt.multipliers.map((item) => `${item.label}: ${fmt(item.value, 3)}`).join('; ')
    : '—'
  const optimumMultiplierText = lp.optimumKkt.multipliers.length
    ? lp.optimumKkt.multipliers.map((item) => `${item.label}: ${fmt(item.value, 3)}`).join('; ')
    : '—'
  const shadowText = lp.shadow
    ? `${lp.shadow.label}: ${fmt(lp.shadow.value, 3)} per unit`
    : '—'
  const rangeText = lp.allowable
    ? `${lp.allowable.label}: [${fmt(lp.allowable.min)}, ${fmt(lp.allowable.max)}]`
    : '—'

  return (
    <FigureFrame
      title={title ?? 'Follow the gradient, or test a constraint'}
      hint={
        mode === 'descent'
          ? 'Change the step size and count. The path is the sequence of points the iteration computes.'
          : 'Drag the candidate or use its sliders. The active constraints, KKT multipliers and dual value all come from the same drawn polygon.'
      }
      caption={
        caption ??
        'The contour plot is the actual objective x² + 2y². Each dot is one gradient step. In the feasible-region mode the shaded polygon is the intersection of the half-planes, and the dual value is read from the multipliers of the active constraints.'
      }
      onReset={reset}
      readout={
        mode === 'descent'
          ? [
              { label: 'objective', value: 'x² + 2y²' },
              { label: 'start', value: `(${fmt(startX)}, ${fmt(startY)})` },
              {
                label: 'gradient',
                value: last ? `(${fmt(last.gradX)}, ${fmt(last.gradY)})` : '—',
              },
              { label: 'value', value: last ? fmt(last.value, 3) : '—' },
              { label: 'steps', value: String(run.points.length) },
              { label: 'status', value: run.status },
              ...(run.note ? [{ label: 'note', value: run.note }] : []),
            ]
          : [
              { label: 'objective', value: 'maximise 3x + 2y' },
              { label: 'candidate', value: `(${fmt(candidateX)}, ${fmt(candidateY)})` },
              { label: 'value', value: fmt(lp.objectiveValue, 3) },
              {
                label: 'status',
                value: lp.status === 'infeasible' ? 'no feasible polygon' : lp.feasible ? 'feasible' : 'violates a constraint',
              },
              { label: 'active at candidate', value: activeText },
              { label: 'KKT at candidate', value: lp.kkt.kkt ? 'holds' : 'fails' },
              { label: 'candidate multipliers', value: multiplierText },
              { label: 'KKT residual', value: fmt(lp.kkt.residual, 4) },
              { label: 'optimum', value: `(${fmt(optimum.x)}, ${fmt(optimum.y)}) = ${fmt(optimum.value, 3)}` },
              { label: 'optimum multipliers', value: optimumMultiplierText },
              { label: 'dual value', value: fmt(lp.dualValue, 3) },
              { label: 'primal-dual gap', value: fmt(lp.gap, 3) },
              { label: 'shadow price', value: shadowText },
              { label: 'allowable range', value: rangeText },
              {
                label: 'constraints',
                value: lp.candidate.map((item) => `${item.label}: ${fmt(item.value, 2)}`).join('; '),
              },
            ]
      }
    >
      <svg
        ref={svgRef}
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        style={{ touchAction: 'none' }}
        role="img"
        aria-label={
          mode === 'descent'
            ? 'Gradient descent path on contour lines'
            : `Feasible region with the candidate at (${fmt(candidateX)}, ${fmt(candidateY)})`
        }
      >
        <defs>
          <clipPath id={clipId}>
            <rect x={PAD} y={PAD} width={W - PAD * 2} height={H - PAD * 2} />
          </clipPath>
          <marker id={arrowId} viewBox="0 0 8 8" refX="7" refY="4" markerWidth="6" markerHeight="6" orient="auto">
            <path d="M0,0 L8,4 L0,8 z" fill="var(--accent-red)" />
          </marker>
        </defs>

        <line x1={PAD} y1={Y(0)} x2={W - PAD} y2={Y(0)} stroke="var(--text-muted)" strokeWidth="0.8" />
        <line x1={X(0)} y1={PAD} x2={X(0)} y2={H - PAD} stroke="var(--text-muted)" strokeWidth="0.8" />

        <g clipPath={`url(#${clipId})`}>
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
              {objectiveLines.map((line) => (
                <line
                  key={line.level}
                  x1={line.from[0]}
                  y1={line.from[1]}
                  x2={line.to[0]}
                  y2={line.to[1]}
                  stroke="var(--chart-3)"
                  strokeWidth={0.9}
                  strokeDasharray="3 3"
                />
              ))}
              {constraintLines.map((line) => (
                <line
                  key={line.label}
                  x1={line.from[0]}
                  y1={line.from[1]}
                  x2={line.to[0]}
                  y2={line.to[1]}
                  stroke={line.active ? 'var(--accent-teal)' : 'var(--text-muted)'}
                  strokeWidth={line.active ? 2.2 : 1}
                />
              ))}
              <line
                x1={X(optimum.x)}
                y1={Y(optimum.y)}
                x2={X(optimum.x + gradientEnd[0])}
                y2={Y(optimum.y + gradientEnd[1])}
                stroke="var(--accent-red)"
                strokeWidth={1.6}
                markerEnd={`url(#${arrowId})`}
              />
              {normalArrows.map((arrow) => (
                <line
                  key={arrow.label}
                  x1={X(optimum.x)}
                  y1={Y(optimum.y)}
                  x2={arrow.to[0]}
                  y2={arrow.to[1]}
                  stroke="var(--accent-teal)"
                  strokeWidth={1.2}
                  markerEnd={`url(#${arrowId})`}
                  opacity={0.8}
                />
              ))}
              <circle cx={X(candidateX)} cy={Y(candidateY)} r={5} fill={lp.feasible ? 'var(--accent-teal)' : 'var(--accent-red)'} />
              <circle cx={X(optimum.x)} cy={Y(optimum.y)} r={3.5} fill="var(--accent-blue)" />
            </>
          )}
        </g>

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
            <div className="flex items-center gap-2">
              <label htmlFor={`${clipId}-rhs`} className="w-24 shrink-0 font-mono text-[10px] text-mathua-muted">
                Shift bound
              </label>
              <select
                id={`${clipId}-rhs`}
                aria-label="Shift bound"
                value={rhs}
                onChange={(event) => setRhs(event.target.value as LpShiftWhich)}
                className="h-7 min-w-0 flex-1 border border-mathua-border bg-transparent px-1 font-mono text-[10px] text-mathua-secondary"
              >
                <option value="none">none</option>
                <option value="first">x + 2y ≤ 4</option>
                <option value="second">2x + y ≤ 5</option>
              </select>
            </div>
            <Slider label="Change δ" value={delta} min={-1} max={1} step={0.05} onChange={setDelta} />
          </div>
        )}
      </div>
    </FigureFrame>
  )
}
