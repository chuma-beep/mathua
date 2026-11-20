'use client'

import { useMemo, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Slider } from '../Controls'
import { rootsSchema } from '../schemas'
import {
  ROOT_FUNCTIONS,
  bisection,
  newton,
  newtonTangent,
  rootFunction,
  type RootMethod,
} from '../math/roots'
import { clipLineToRect } from '../math/conic'
import { niceTicks, scaleLinear } from '../scale'
import { fmt } from '../format'
import type { TemplateProps } from '../types'

const W = 360
const H = 320
const PAD = 30
const X_MIN = -3
const X_MAX = 3
const Y_MIN = -4
const Y_MAX = 6

export default function RootIteration({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => rootsSchema.parse(params), [params])
  const [fn, setFn] = useState(defaults.fn)
  const [method, setMethod] = useState<RootMethod>(defaults.method)
  const [a, setA] = useState(defaults.a)
  const [b, setB] = useState(defaults.b)
  const [start, setStart] = useState(defaults.start)
  const [iterations, setIterations] = useState(defaults.iterations)

  const spec = rootFunction(fn)
  const X = scaleLinear(X_MIN, X_MAX, PAD, W - PAD)
  const Y = scaleLinear(Y_MAX, Y_MIN, PAD, H - PAD)
  const ticksX = niceTicks(X_MIN, X_MAX, 4).filter((value) => Math.abs(value) > 1e-9)
  const ticksY = niceTicks(Y_MIN, Y_MAX, 4).filter((value) => Math.abs(value) > 1e-9)

  const curve = useMemo(() => {
    const points: Array<[number, number]> = []
    for (let index = 0; index <= 240; index++) {
      const x = X_MIN + ((X_MAX - X_MIN) * index) / 240
      const y = spec.f(x)
      if (Number.isFinite(y)) points.push([X(x), Y(Math.max(Y_MIN - 2, Math.min(Y_MAX + 2, y)))])
    }
    return points.map((point, index) => `${index === 0 ? 'M' : 'L'}${point[0].toFixed(2)},${point[1].toFixed(2)}`).join(' ')
  }, [spec, X, Y])

  const bisectionResult = useMemo(
    () => (method === 'bisection' ? bisection(spec.f, a, b, iterations) : null),
    [method, spec, a, b, iterations],
  )
  const newtonResult = useMemo(
    () => (method === 'newton' ? newton(spec.f, spec.df, start, iterations) : null),
    [method, spec, start, iterations],
  )
  const result = bisectionResult ?? newtonResult
  const root = result?.root ?? Number.NaN
  const error = Number.isFinite(root) ? Math.abs(root - spec.exactValue) : Number.NaN

  const reset = () => {
    setFn(defaults.fn)
    setMethod(defaults.method)
    setA(defaults.a)
    setB(defaults.b)
    setStart(defaults.start)
    setIterations(defaults.iterations)
  }

  const plotLine = (point: [number, number], direction: [number, number], key: string, dashed = false) => {
    const clipped = clipLineToRect(point, direction, { x0: X_MIN, x1: X_MAX, y0: Y_MIN, y1: Y_MAX })
    if (!clipped) return null
    return (
      <line
        key={key}
        x1={X(clipped[0][0])}
        y1={Y(clipped[0][1])}
        x2={X(clipped[1][0])}
        y2={Y(clipped[1][1])}
        stroke="var(--accent-blue)"
        strokeWidth="1.2"
        strokeDasharray={dashed ? '4 3' : undefined}
      />
    )
  }

  const currentBracket =
    bisectionResult && bisectionResult.steps.length > 0
      ? bisectionResult.steps[bisectionResult.steps.length - 1]
      : null
  const lastNewton = newtonResult && newtonResult.steps.length > 0 ? newtonResult.steps[newtonResult.steps.length - 1] : null

  return (
    <FigureFrame
      title={title ?? 'Iteration converges to the root'}
      hint="Change the method, the bracket or the starting guess. Bisection keeps the sign change; Newton follows the tangent to the axis."
      caption={
        caption ??
        'A numerical root is the limit of an iteration, not a guess read off the graph. Bisection halves a bracket with a sign change. Newton follows the tangent at the current guess to its x-intercept. The plotted curve is the function; the markers and lines are the actual iterations.'
      }
      onReset={reset}
      readout={[
        { label: 'function', value: spec.expression },
        { label: 'method', value: method === 'bisection' ? 'Bisection' : 'Newton' },
        { label: 'iterations', value: String(iterations) },
        { label: 'root estimate', value: Number.isFinite(root) ? fmt(root, 8) : '—' },
        { label: 'exact root', value: spec.exact },
        { label: 'error', value: Number.isFinite(error) ? fmt(error, 3) : '—' },
        { label: 'status', value: result?.status ?? 'invalid' },
        ...(method === 'bisection'
          ? [
              { label: 'initial bracket', value: `[${fmt(a, 4)}, ${fmt(b, 4)}]` },
              { label: 'bracket', value: currentBracket ? `[${fmt(currentBracket.a, 4)}, ${fmt(currentBracket.b, 4)}]` : '—' },
              { label: 'f(a), f(b)', value: currentBracket ? `${fmt(spec.f(currentBracket.a), 3)}, ${fmt(spec.f(currentBracket.b), 3)}` : '—' },
            ]
          : [
              { label: 'last guess', value: lastNewton ? fmt(lastNewton.x, 5) : '—' },
              { label: 'slope', value: lastNewton ? fmt(lastNewton.slope, 3) : '—' },
            ]),
        ...(result?.note ? [{ label: 'note', value: result.note }] : []),
      ]}
    >
      <svg
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        role="img"
        aria-label={`${method} iteration for ${spec.expression}`}
      >
        <defs>
          <clipPath id="root-plot">
            <rect x={PAD} y={PAD} width={W - PAD * 2} height={H - PAD * 2} />
          </clipPath>
        </defs>

        <line x1={PAD} y1={Y(0)} x2={W - PAD} y2={Y(0)} stroke="var(--text-muted)" strokeWidth="0.8" />
        <line x1={X(0)} y1={PAD} x2={X(0)} y2={H - PAD} stroke="var(--text-muted)" strokeWidth="0.8" />

        <g clipPath="url(#root-plot)">
          <path d={curve} fill="none" stroke="var(--text-primary)" strokeWidth="1.6" />

          {bisectionResult && currentBracket ? (
            <>
              <rect
                x={X(Math.min(currentBracket.a, currentBracket.b))}
                y={PAD}
                width={Math.abs(X(currentBracket.b) - X(currentBracket.a))}
                height={H - PAD * 2}
                fill="var(--accent-blue)"
                opacity="0.08"
              />
              {bisectionResult.steps.map((step) => (
                <g key={step.index}>
                  <circle cx={X(step.midpoint)} cy={Y(0)} r="2.4" fill="var(--accent-teal)" />
                  <line
                    x1={X(step.midpoint)}
                    y1={Y(0)}
                    x2={X(step.midpoint)}
                    y2={Y(step.fMidpoint)}
                    stroke="var(--accent-teal)"
                    strokeWidth="0.8"
                    strokeDasharray="2 2"
                  />
                </g>
              ))}
              <circle cx={X(a)} cy={Y(0)} r="3.5" fill="var(--accent-blue)" />
              <circle cx={X(b)} cy={Y(0)} r="3.5" fill="var(--accent-blue)" />
              <text x={X(a)} y={Y(0) + 16} textAnchor="middle" fontSize="9" fill="var(--accent-blue)">
                a
              </text>
              <text x={X(b)} y={Y(0) + 16} textAnchor="middle" fontSize="9" fill="var(--accent-blue)">
                b
              </text>
            </>
          ) : null}

          {newtonResult
            ? newtonResult.steps.map((step) => {
                const tangent = newtonTangent(spec.f, spec.df, step.x)
                return (
                  <g key={step.index}>
                    {Number.isFinite(tangent.slope)
                      ? plotLine([step.x, tangent.y], [1, tangent.slope], `tangent-${step.index}`, true)
                      : null}
                    <circle cx={X(step.x)} cy={Y(tangent.y)} r="3" fill="var(--accent-blue)" />
                    <line
                      x1={X(step.x)}
                      y1={Y(tangent.y)}
                      x2={X(step.next)}
                      y2={Y(0)}
                      stroke="var(--accent-teal)"
                      strokeWidth="0.9"
                      strokeDasharray="2 2"
                    />
                    <circle cx={X(step.next)} cy={Y(0)} r="3" fill="var(--accent-teal)" />
                  </g>
                )
              })
            : null}

          {Number.isFinite(root) ? <circle cx={X(root)} cy={Y(0)} r="4" fill="var(--accent-green)" /> : null}
        </g>

        <g aria-hidden="true">
          {ticksX.map((value) => (
            <g key={value}>
              <line x1={X(value)} y1={Y(0) - 3} x2={X(value)} y2={Y(0) + 3} stroke="var(--text-muted)" strokeWidth="0.8" />
              <text x={X(value)} y={Y(0) + 15} textAnchor="middle" fontSize="8" fill="var(--text-muted)">
                {value}
              </text>
            </g>
          ))}
          {ticksY.map((value) => (
            <g key={value}>
              <line x1={X(0) - 3} y1={Y(value)} x2={X(0) + 3} y2={Y(value)} stroke="var(--text-muted)" strokeWidth="0.8" />
              <text x={X(0) - 7} y={Y(value) + 3} textAnchor="end" fontSize="8" fill="var(--text-muted)">
                {value}
              </text>
            </g>
          ))}
        </g>

        <text x={PAD + 2} y={PAD + 10} fontSize="9" fill="var(--text-muted)">
          {spec.expression}
        </text>
      </svg>

      <div className="mt-3 space-y-2">
        <div className="flex flex-wrap items-center gap-2">
          <label className="font-mono text-[10px] text-mathua-muted" htmlFor="root-function">
            Function
          </label>
          <select
            id="root-function"
            aria-label="Function"
            value={fn}
            onChange={(event) => setFn(event.target.value as typeof fn)}
            className="h-7 border border-mathua-border bg-transparent px-1 font-mono text-[10px] text-mathua-secondary"
          >
            {ROOT_FUNCTIONS.map((item) => (
              <option key={item.id} value={item.id}>
                {item.expression}
              </option>
            ))}
          </select>
          <label className="font-mono text-[10px] text-mathua-muted" htmlFor="root-method">
            Method
          </label>
          <select
            id="root-method"
            aria-label="Method"
            value={method}
            onChange={(event) => setMethod(event.target.value as RootMethod)}
            className="h-7 border border-mathua-border bg-transparent px-1 font-mono text-[10px] text-mathua-secondary"
          >
            <option value="bisection">bisection</option>
            <option value="newton">Newton</option>
          </select>
        </div>

        <Slider label="Iterations" value={iterations} min={1} max={12} step={1} onChange={setIterations} />
        {method === 'bisection' ? (
          <>
            <Slider label="a" value={a} min={-3} max={3} step={0.1} onChange={setA} />
            <Slider label="b" value={b} min={-3} max={3} step={0.1} onChange={setB} />
          </>
        ) : (
          <Slider label="start" value={start} min={-3} max={3} step={0.1} onChange={setStart} />
        )}
      </div>
    </FigureFrame>
  )
}
