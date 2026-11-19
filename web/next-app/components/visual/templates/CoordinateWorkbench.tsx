'use client'

import { useCallback, useId, useMemo, useRef, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { NumberField, Select, Slider } from '../Controls'
import { XAxisTicks, YAxisTicks } from '../Axis'
import { coordinateSchema } from '../schemas'
import { distance, intersect, lineThrough, lineWithSlope, midpoint, type Line, type Point } from '../math/coords'
import { clamp, scaleLinear } from '../scale'
import { fmt } from '../format'
import { useDrag } from '../useDrag'
import type { TemplateProps } from '../types'

const W = 340
const H = 320
const PAD = 30
const EXTENT = 5

type Mode = 'measure' | 'intersect'

const MODE_OPTIONS: Array<{ value: Mode; label: string }> = [
  { value: 'measure', label: 'Distance and midpoint' },
  { value: 'intersect', label: 'Two lines meet' },
]

/** A minus sign for display, because the ASCII hyphen reads as a dash in prose. */
function num(value: number): string {
  return fmt(value, 2).replace('-', '−')
}

/** Wraps a negative number so a subtraction reads as `4 − (−3)` rather than `4 − -3`. */
function wrap(value: number): string {
  const text = num(value)
  return value < 0 ? `(${text})` : text
}

function lineText(line: Line): string {
  if (line.kind === 'vertical') return `x = ${num(line.x)}`
  return `y = ${num(line.m)}x ${line.b < 0 ? '−' : '+'} ${num(Math.abs(line.b))}`
}

export default function CoordinateWorkbench({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => {
    const p = coordinateSchema.parse(params)
    return {
      mode: p.mode as Mode,
      a: [Number(p.a[0]), Number(p.a[1])] as Point,
      b: [Number(p.b[0]), Number(p.b[1])] as Point,
      c: [Number(p.c[0]), Number(p.c[1])] as Point,
      slope: p.slope,
    }
  }, [params])

  const [mode, setMode] = useState<Mode>(defaults.mode)
  const [a, setA] = useState<Point>(defaults.a)
  const [b, setB] = useState<Point>(defaults.b)
  const [c, setC] = useState<Point>(defaults.c)
  const [slope, setSlope] = useState(defaults.slope)
  const [selected, setSelected] = useState(0)
  const svgRef = useRef<SVGSVGElement | null>(null)
  const uid = useId().replace(/:/g, '')

  const points: Point[] = mode === 'measure' ? [a, b, c] : [a, b, c]
  const setters = [setA, setB, setC] as const
  const labels = ['A', 'B', 'C'] as const

  const X = scaleLinear(-EXTENT, EXTENT, PAD, W - PAD)
  const Y = scaleLinear(EXTENT, -EXTENT, H - PAD, PAD)
  const unit = (W - PAD * 2) / (EXTENT * 2)
  const cx = X(0)
  const cy = Y(0)

  const m = midpoint(a, b)
  const d = distance(a, b)
  // One point does not determine a line. When A and B coincide the line through
  // them is undefined, and reporting the vertical line that `lineThrough` falls
  // back to would invent a line from a single point.
  const samePoints = distance(a, b) < 1e-9
  const line1: Line | null = samePoints ? null : lineThrough(a, b)
  const line2 = lineWithSlope(c, slope)
  const meet = line1 ? intersect(line1, line2) : null
  // The lines keep going outside the frame, so a shallow pair can meet well beyond
  // it. The coordinate is still true, but a green dot at (14.5, 11.5) is clipped
  // and the reader sees "meets" with no crossing anywhere on screen.
  const meetVisible =
    meet?.kind === 'point' && Math.abs(meet.x) <= EXTENT && Math.abs(meet.y) <= EXTENT

  const corner: Point = [b[0], a[1]]

  const onDrag = useCallback(
    (point: { x: number; y: number }) => {
      const px = clamp((point.x - cx) / unit, -4.5, 4.5)
      const py = clamp((cy - point.y) / unit, -4.5, 4.5)
      const target = Math.round(px * 10) / 10
      const targetY = Math.round(py * 10) / 10
      let best = 0
      let bestDistance = Number.POSITIVE_INFINITY
      points.forEach(([x, y], i) => {
        const gap = Math.hypot(x - target, y - targetY)
        if (gap < bestDistance) {
          bestDistance = gap
          best = i
        }
      })
      setSelected(best)
      setters[best]([target, targetY])
    },
    // `points` is rebuilt each render from a, b and c, so depend on those.
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [cx, cy, unit, a, b, c],
  )
  useDrag(svgRef, onDrag)

  const editCoordinate = (axis: 0 | 1, value: number) => {
    const current = points[selected]
    const next: Point = axis === 0 ? [clamp(value, -4.5, 4.5), current[1]] : [current[0], clamp(value, -4.5, 4.5)]
    setters[selected](next)
  }

  const reset = () => {
    setMode(defaults.mode)
    setA(defaults.a)
    setB(defaults.b)
    setC(defaults.c)
    setSlope(defaults.slope)
    setSelected(0)
  }

  const working = useMemo(() => {
    if (mode === 'measure') {
      return [
        `AB = √((x₂ − x₁)² + (y₂ − y₁)²)`,
        `= √((${num(b[0])} − ${wrap(a[0])})² + (${num(b[1])} − ${wrap(a[1])})²)`,
        `= √(${num((b[0] - a[0]) ** 2)} + ${num((b[1] - a[1]) ** 2)}) = √${num((b[0] - a[0]) ** 2 + (b[1] - a[1]) ** 2)} ≈ ${num(d)}`,
        `M = ((x₁ + x₂)/2, (y₁ + y₂)/2) = (${num(m[0])}, ${num(m[1])})`,
      ]
    }
    if (!line1 || !meet) {
      return [
        'A and B are the same point, so they do not determine a line.',
        `line 2: ${lineText(line2)} (through C with slope ${num(slope)})`,
        'Move A or B apart to get a line through two points.',
      ]
    }
    const second = `line 2: ${lineText(line2)} (through C with slope ${num(slope)})`
    const first = `line 1: ${lineText(line1)} (through A and B)`
    if (meet.kind === 'parallel') {
      return [first, second, 'Same slope, different intercept: no point is shared.']
    }
    if (meet.kind === 'coincident') {
      return [first, second, 'The equations describe one line, so every point on it is shared.']
    }
    return [first, second, `Shared point: (${num(meet.x)}, ${num(meet.y)})`]
  }, [mode, a, b, slope, d, m, line1, line2, meet])

  const readout =
    mode === 'measure'
      ? [
          { label: 'A', value: `(${fmt(a[0])}, ${fmt(a[1])})` },
          { label: 'B', value: `(${fmt(b[0])}, ${fmt(b[1])})` },
          { label: 'Δx', value: fmt(b[0] - a[0]) },
          { label: 'Δy', value: fmt(b[1] - a[1]) },
          { label: 'distance', value: fmt(d, 3) },
          { label: 'midpoint', value: `(${fmt(m[0])}, ${fmt(m[1])})` },
        ]
      : [
          { label: 'line 1', value: line1 ? lineText(line1) : 'undefined' },
          { label: 'line 2', value: lineText(line2) },
          {
            label: 'intersection',
            value: !meet
              ? 'undefined'
              : meet.kind === 'point'
                ? `(${fmt(meet.x)}, ${fmt(meet.y)})`
                : meet.kind === 'parallel'
                  ? 'none'
                  : 'all shared',
          },
          {
            label: 'status',
            value: !meet ? 'undefined' : meet.kind === 'point' ? 'meets' : meet.kind === 'parallel' ? 'parallel' : 'coincident',
          },
        ]

  const lineEnds = (line: Line): [Point, Point] => {
    if (line.kind === 'vertical') {
      return [
        [line.x, -EXTENT],
        [line.x, EXTENT],
      ]
    }
    return [
      [-EXTENT, line.m * -EXTENT + line.b],
      [EXTENT, line.m * EXTENT + line.b],
    ]
  }

  return (
    <FigureFrame
      title={title ?? 'Points, distances and lines on the plane'}
      hint="Drag a point, or choose one and type its coordinates. The working below follows the drawing."
      caption={
        caption ??
        'The distance is the hypotenuse of the horizontal and vertical gaps, and the midpoint is their average. Two lines meet unless their slopes agree, in which case they are parallel or the same line.'
      }
      onReset={reset}
      readout={readout}
    >
      <svg
        ref={svgRef}
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] cursor-grab font-mono active:cursor-grabbing"
        style={{ touchAction: 'none' }}
        role="img"
        aria-label="Points, segments and lines on a coordinate plane"
      >
        <defs>
          <clipPath id={`${uid}-plot`}>
            <rect x={PAD} y={PAD} width={W - PAD * 2} height={H - PAD * 2} />
          </clipPath>
        </defs>

        <g clipPath={`url(#${uid}-plot)`}>
          <XAxisTicks ticks={[-4, -2, 2, 4]} scale={X} at={cy} />
          <YAxisTicks ticks={[-4, -2, 2, 4]} scale={Y} at={cx} />
          <line x1={PAD} y1={cy} x2={W - PAD} y2={cy} stroke="var(--text-muted)" strokeWidth="0.8" />
          <line x1={cx} y1={PAD} x2={cx} y2={H - PAD} stroke="var(--text-muted)" strokeWidth="0.8" />

          {mode === 'measure' ? (
            <>
              <line x1={X(a[0])} y1={Y(a[1])} x2={X(corner[0])} y2={Y(corner[1])} stroke="var(--accent-teal)" strokeWidth="1" strokeDasharray="4 3" />
              <line x1={X(corner[0])} y1={Y(corner[1])} x2={X(b[0])} y2={Y(b[1])} stroke="var(--accent-teal)" strokeWidth="1" strokeDasharray="4 3" />
              <path
                d={`M ${X(corner[0])} ${Y(corner[1] - 0.35 * Math.sign(b[1] - a[1] || 1))} L ${X(corner[0] + 0.35 * Math.sign(b[0] - a[0] || 1))} ${Y(corner[1] - 0.35 * Math.sign(b[1] - a[1] || 1))} L ${X(corner[0] + 0.35 * Math.sign(b[0] - a[0] || 1))} ${Y(corner[1])}`}
                fill="none"
                stroke="var(--text-muted)"
                strokeWidth="0.8"
              />
              <line x1={X(a[0])} y1={Y(a[1])} x2={X(b[0])} y2={Y(b[1])} stroke="var(--accent-blue)" strokeWidth="1.8" />
              <text x={X((a[0] + corner[0]) / 2)} y={Y(a[1]) - 6} textAnchor="middle" fontSize="9" fill="var(--accent-teal)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
                Δx {num(b[0] - a[0])}
              </text>
              <text x={X(corner[0]) + 6} y={Y((a[1] + b[1]) / 2)} fontSize="9" fill="var(--accent-teal)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
                Δy {num(b[1] - a[1])}
              </text>
              <circle cx={X(m[0])} cy={Y(m[1])} r="4" fill="none" stroke="var(--accent-green)" strokeWidth="1.6" />
              <text x={X(m[0]) + 8} y={Y(m[1]) + 16} fontSize="10" fill="var(--accent-green)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
                M
              </text>
            </>
          ) : (
            <>
              {[line1, line2].filter((line): line is Line => line !== null).map((line, i) => {
                const [from, to] = lineEnds(line)
                return (
                  <line
                    key={i}
                    x1={X(from[0])}
                    y1={Y(from[1])}
                    x2={X(to[0])}
                    y2={Y(to[1])}
                    stroke={i === 0 ? 'var(--accent-blue)' : 'var(--accent-teal)'}
                    strokeWidth="1.6"
                  />
                )
              })}
              {meet?.kind === 'point' && meetVisible && <circle cx={X(meet.x)} cy={Y(meet.y)} r="5" fill="var(--accent-green)" />}
              {meet?.kind === 'point' && !meetVisible && (
                <text x={cx} y={PAD + 12} textAnchor="middle" fontSize="10" fill="var(--accent-teal)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
                  meets outside this window
                </text>
              )}
              {!meet && (
                <text x={cx} y={PAD + 12} textAnchor="middle" fontSize="10" fill="var(--accent-red)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
                  A and B coincide: no line
                </text>
              )}
              {meet && meet.kind !== 'point' && (
                <text x={cx} y={PAD + 12} textAnchor="middle" fontSize="10" fill="var(--accent-red)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
                  {meet.kind === 'parallel' ? 'parallel: no shared point' : 'coincident: every point shared'}
                </text>
              )}
            </>
          )}

          {(mode === 'measure' ? [a, b] : [a, b, c]).map((point, i) => (
            <g key={i}>
              <circle cx={X(point[0])} cy={Y(point[1])} r={selected === i ? 5.5 : 4} fill={selected === i ? 'var(--accent-teal)' : 'var(--accent-blue)'} />
              <text x={X(point[0]) + 8} y={Y(point[1]) - 8} fontSize="10" fill="var(--text-primary)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
                {labels[i]}
              </text>
            </g>
          ))}
        </g>
      </svg>

      <div className="mt-2 space-y-0.5 border border-mathua-border bg-mathua-code p-2 font-mono text-[10px] leading-relaxed text-mathua-secondary" data-testid="visual-working">
        {working.map((line, i) => (
          <p key={i} className="whitespace-pre-wrap break-words">
            {line}
          </p>
        ))}
      </div>

      <div className="mt-3 space-y-2">
        <Select label="Mode" value={mode} options={MODE_OPTIONS} onChange={(v) => setMode(v as Mode)} />
        <Select
          label="Point"
          value={String(selected)}
          options={labels.map((label, i) => ({ value: String(i), label }))}
          onChange={(v) => setSelected(Number(v))}
        />
        <div className="grid grid-cols-2 gap-2">
          <NumberField label="x" value={points[selected][0]} min={-4.5} max={4.5} onChange={(v) => editCoordinate(0, v)} />
          <NumberField label="y" value={points[selected][1]} min={-4.5} max={4.5} onChange={(v) => editCoordinate(1, v)} />
        </div>
        {mode === 'intersect' && <Slider label="slope" value={slope} min={-3} max={3} step={0.1} onChange={setSlope} />}
      </div>
    </FigureFrame>
  )
}
