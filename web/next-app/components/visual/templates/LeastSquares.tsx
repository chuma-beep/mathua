'use client'

import { useCallback, useMemo, useRef, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { NumberField, Select, Toggle } from '../Controls'
import { XAxisTicks } from '../Axis'
import { regressionSchema } from '../schemas'
import { correlation, hasXSpread, intercept, means, predict, slope, sse, type Point } from '../math/regression'
import { clamp, scaleLinear } from '../scale'
import { fmt } from '../format'
import { useDrag } from '../useDrag'
import type { TemplateProps } from '../types'

const W = 340
const H = 280
const PAD = 30
const X_MIN = -1
const X_MAX = 7
const Y_MIN = -1
const Y_MAX = 8

export default function LeastSquares({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => {
    const p = regressionSchema.parse(params)
    return {
      points: p.points.map((pair) => [Number(pair[0]), Number(pair[1])] as Point),
      showResiduals: p.showResiduals,
    }
  }, [params])

  const [points, setPoints] = useState<Point[]>(defaults.points)
  const [showResiduals, setShowResiduals] = useState(defaults.showResiduals)
  const [selected, setSelected] = useState(0)
  const svgRef = useRef<SVGSVGElement | null>(null)

  const X = scaleLinear(X_MIN, X_MAX, PAD, W - PAD)
  const Y = scaleLinear(Y_MIN, Y_MAX, H - PAD, PAD)
  const invX = scaleLinear(PAD, W - PAD, X_MIN, X_MAX)
  const invY = scaleLinear(H - PAD, PAD, Y_MIN, Y_MAX)

  const m = slope(points)
  const b = intercept(points)
  const r = correlation(points)
  const [mx, my] = means(points)
  // A cloud whose x values are all equal has no unique least-squares line, so the
  // figure reports the fit as undefined instead of drawing the zero slope the
  // guard in `slope` returns.
  const defined = hasXSpread(points)

  // The nearest point takes the drag, so overlapping points stay separable.
  const onDrag = useCallback(
    (point: { x: number; y: number }) => {
      const px = clamp(invX(point.x), X_MIN + 0.1, X_MAX - 0.1)
      const py = clamp(invY(point.y), Y_MIN + 0.1, Y_MAX - 0.1)
      setPoints((prev) => {
        let best = 0
        let bestDistance = Number.POSITIVE_INFINITY
        prev.forEach(([x, y], i) => {
          const d = Math.hypot(x - px, y - py)
          if (d < bestDistance) {
            bestDistance = d
            best = i
          }
        })
        const next = [...prev]
        next[best] = [px, py]
        return next
      })
    },
    [invX, invY],
  )
  useDrag(svgRef, onDrag)

  const editPoint = (axis: 0 | 1, value: number) => {
    setPoints((prev) => {
      const next = [...prev]
      // Removing a point can leave the selection past the end; clamping here stops
      // an edit from appending a stray point at a stale index.
      const index = Math.min(selected, next.length - 1)
      const current = next[index] ?? [0, 0]
      next[index] = axis === 0 ? [clamp(value, X_MIN + 0.1, X_MAX - 0.1), current[1]] : [current[0], clamp(value, Y_MIN + 0.1, Y_MAX - 0.1)]
      return next
    })
  }

  const reset = () => {
    setPoints(defaults.points)
    setSelected(0)
    setShowResiduals(defaults.showResiduals)
  }

  const addPoint = () => {
    if (points.length >= 12) return
    // A deterministic new point: the free spot above the current line, so it
    // visibly moves the fit without random noise.
    const x = clamp((points.length * 1.3) % 6 + 0.6, X_MIN + 0.5, X_MAX - 0.5)
    const y = clamp(b + m * x + 1.2, Y_MIN + 0.5, Y_MAX - 0.5)
    setPoints([...points, [x, y]])
  }

  const removePoint = () => {
    if (points.length <= 2) return
    setPoints(points.slice(0, -1))
  }

  return (
    <FigureFrame
      title={title ?? 'The least-squares line'}
      hint="Drag a point. The line follows, and the sum of squares falls to the minimum for these points."
      caption={
        caption ??
        'The line minimizes the sum of the squared vertical gaps, so a single far point pulls harder than a near one. The line always passes through the mean point of the data.'
      }
      onReset={reset}
      readout={[
        { label: 'slope', value: defined ? fmt(m, 3) : 'undefined' },
        { label: 'intercept', value: defined ? fmt(b, 3) : 'undefined' },
        { label: 'r', value: defined ? fmt(r, 3) : 'undefined' },
        { label: 'sum of squares', value: defined ? fmt(sse(points), 3) : 'undefined' },
        { label: 'mean point', value: `(${fmt(mx)}, ${fmt(my)})` },
        { label: 'at x = 6', value: defined ? fmt(predict(points, 6)) : 'undefined' },
      ]}
    >
      <svg
        ref={svgRef}
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        style={{ touchAction: 'none' }}
        role="img"
        aria-label="A scatter plot with its least-squares line and residuals"
      >
        <line x1={PAD} y1={Y(0)} x2={W - PAD} y2={Y(0)} stroke="var(--text-muted)" strokeWidth="0.8" />
        <line x1={X(0)} y1={PAD} x2={X(0)} y2={H - PAD} stroke="var(--text-muted)" strokeWidth="0.8" />

        {defined && showResiduals &&
          points.map(([x, y], i) => (
            <line
              key={`r${i}`}
              x1={X(x)}
              y1={Y(y)}
              x2={X(x)}
              y2={Y(b + m * x)}
              stroke="var(--accent-teal)"
              strokeWidth="1"
              strokeDasharray="2 3"
              opacity="0.75"
            />
          ))}

        {defined && (
          <line
            x1={X(X_MIN)}
            y1={Y(b + m * X_MIN)}
            x2={X(X_MAX)}
            y2={Y(b + m * X_MAX)}
            stroke="var(--accent-blue)"
            strokeWidth="1.8"
          />
        )}
        {!defined && (
          <text x={PAD + 2} y={PAD + 10} fontSize="9" fill="var(--accent-red)">
            x does not vary: no unique line
          </text>
        )}

        <circle cx={X(mx)} cy={Y(my)} r="4.5" fill="none" stroke="var(--accent-green)" strokeWidth="1.4" />
        {points.map(([x, y], i) => (
          <circle key={i} cx={X(x)} cy={Y(y)} r="4" fill="var(--accent-blue)" />
        ))}

        <XAxisTicks ticks={[0, 2, 4, 6]} scale={X} at={Y(0)} />
        <text x={PAD + 2} y={PAD + 10} fontSize="9" fill="var(--text-muted)">
          mean point
        </text>
      </svg>

      <div className="mt-3 space-y-2">
        <div className="flex flex-wrap items-center gap-2">
          <button
            type="button"
            onClick={addPoint}
            disabled={points.length >= 12}
            className="border border-mathua-border px-2 py-0.5 font-mono text-[10px] text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue disabled:opacity-40"
          >
            Add point
          </button>
          <button
            type="button"
            onClick={removePoint}
            disabled={points.length <= 2}
            className="border border-mathua-border px-2 py-0.5 font-mono text-[10px] text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue disabled:opacity-40"
          >
            Remove point
          </button>
          <span className="font-mono text-[10px] text-mathua-muted">{points.length} points</span>
        </div>
        {/* The keyboard path for the drag: choose a point, then type its two
            coordinates. The select sits on its own row so neither field has to
            shrink to a sliver at 320px. */}
        <Select
          label="Point"
          value={String(Math.min(selected, points.length - 1))}
          options={points.map((_, i) => ({ value: String(i), label: `P${i + 1}` }))}
          onChange={(v) => setSelected(Number(v))}
        />
        <div className="grid grid-cols-2 gap-2">
          <NumberField label="x" value={points[Math.min(selected, points.length - 1)][0]} min={X_MIN} max={X_MAX} onChange={(v) => editPoint(0, v)} />
          <NumberField label="y" value={points[Math.min(selected, points.length - 1)][1]} min={Y_MIN} max={Y_MAX} onChange={(v) => editPoint(1, v)} />
        </div>
        <Toggle label="Show residuals" checked={showResiduals} onChange={setShowResiduals} />
      </div>
    </FigureFrame>
  )
}
