'use client'

import { useId, useMemo, useRef, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Slider } from '../Controls'
import { XAxisTicks, YAxisTicks } from '../Axis'
import { conicSchema } from '../schemas'
import {
  classifyGeneralConic,
  clipLineToRect,
  conicParameterRange,
  conicPointAt,
  discriminant,
  focalMeasurements,
  generalConic,
  projectToConic,
  type ConicBranch,
  type ConicKind,
  type Intercepts,
  type Vec2,
} from '../math/conic'
import { contourSegments } from '../math/contour'
import { invertLinear, niceTicks, scaleLinear } from '../scale'
import { fmt } from '../format'
import { useDrag } from '../useDrag'
import type { TemplateProps } from '../types'

const TYPE_LABELS: Record<ConicKind, string> = {
  circle: 'circle',
  ellipse: 'ellipse',
  hyperbola: 'hyperbola',
  parabola: 'parabola',
  point: 'a single point',
  intersecting: 'two intersecting lines',
  parallel: 'two parallel lines',
  line: 'one line',
  empty: 'no real points',
  'all-plane': 'every point',
  invalid: 'invalid coefficients',
}

const W = 340
const H = 300
const PAD = 30
const EXTENT = 2.8
const BOUNDS = { x0: -EXTENT, x1: EXTENT, y0: -EXTENT, y1: EXTENT }

type Coefficients = [number, number, number, number, number, number]

const PRESETS: Array<{ label: string; coefficients: Coefficients }> = [
  { label: 'Circle', coefficients: [1, 0, 1, 0, 0, -1] },
  { label: 'Ellipse', coefficients: [1, 0, 0.5, 0, 0, -1] },
  { label: 'Parabola', coefficients: [1, 0, 0, 0, -1, 0] },
  { label: 'Hyperbola', coefficients: [1, 0, -1, 0, 0, -1] },
  { label: 'Translated ellipse', coefficients: [1, 0, 4, -2, 16, 13] },
  { label: 'Degenerate', coefficients: [1, 0, -1, 0, 0, 0] },
]

const sameCoefficients = (a: Coefficients, b: Coefficients) =>
  a.every((value, index) => Math.abs(value - b[index]) < 1e-9)

function describeIntercepts(intercepts: Intercepts, axis: 'x' | 'y'): string {
  switch (intercepts.status) {
    case 'two':
      return intercepts.values.map((value) => `${axis} = ${fmt(value)}`).join(', ')
    case 'one':
      return `${axis} = ${fmt(intercepts.values[0])}`
    case 'double':
      return `${axis} = ${fmt(intercepts.values[0])} (double)`
    case 'coincident':
      return `every ${axis}`
    case 'none':
      return 'none'
  }
}

function pointText(point: Vec2): string {
  return `(${fmt(point[0])}, ${fmt(point[1])})`
}

export default function ConicSection({ params, title, caption }: TemplateProps) {
  const clipId = useId()
  const svgRef = useRef<SVGSVGElement | null>(null)
  const defaults = useMemo(() => conicSchema.parse(params), [params])
  const initial = useMemo<Coefficients>(
    () => [defaults.A, defaults.B, defaults.C, defaults.D, defaults.E, defaults.F],
    [defaults],
  )
  const [coefficients, setCoefficients] = useState<Coefficients>(initial)
  const [rawProbe, setRawProbe] = useState(0.7)
  const [branch, setBranch] = useState<ConicBranch>(1)
  const [A, B, C, D, E, F] = coefficients

  const geometry = useMemo(() => classifyGeneralConic(A, B, C, D, E, F), [A, B, C, D, E, F])
  const curve =
    geometry.kind === 'circle' ||
    geometry.kind === 'ellipse' ||
    geometry.kind === 'hyperbola' ||
    geometry.kind === 'parabola'
  const probeRange = conicParameterRange(geometry)
  const probeParameter = probeRange
    ? Math.min(probeRange.max, Math.max(probeRange.min, rawProbe))
    : 0
  const probePoint = curve && probeRange ? conicPointAt(geometry, probeParameter, branch) : null
  const measurements = probePoint ? focalMeasurements(geometry, probePoint) : null
  const probeOutside =
    probePoint !== null &&
    (probePoint[0] < BOUNDS.x0 ||
      probePoint[0] > BOUNDS.x1 ||
      probePoint[1] < BOUNDS.y0 ||
      probePoint[1] > BOUNDS.y1)

  const segments = useMemo(
    () =>
      curve
        ? contourSegments((x, y) => generalConic(A, B, C, D, E, F, x, y), BOUNDS, 0, 96, 96)
        : [],
    [A, B, C, D, E, F, curve],
  )

  const X = scaleLinear(-EXTENT, EXTENT, PAD, W - PAD)
  const Y = scaleLinear(EXTENT, -EXTENT, PAD, H - PAD)
  const toX = invertLinear(-EXTENT, EXTENT, PAD, W - PAD)
  const toY = invertLinear(EXTENT, -EXTENT, PAD, H - PAD)
  const ticks = niceTicks(-EXTENT, EXTENT, 4).filter((value) => Math.abs(value) > 1e-9)

  const drag = (point: { x: number; y: number }) => {
    if (!curve) return
    const projection = projectToConic(geometry, [toX(point.x), toY(point.y)], branch)
    if (projection) setRawProbe(projection.parameter)
  }
  useDrag(svgRef, drag, curve)

  const reset = () => {
    setCoefficients(initial)
    setRawProbe(0.7)
    setBranch(1)
  }

  const selectPreset = (label: string) => {
    const preset = PRESETS.find((item) => item.label === label)
    if (preset) {
      setCoefficients(preset.coefficients)
      setRawProbe(0.7)
      setBranch(1)
    }
  }

  const selected = PRESETS.find((preset) => sameCoefficients(preset.coefficients, coefficients))

  const setCoefficient = (index: number, value: number) => {
    setCoefficients((previous) => {
      const next = [...previous] as Coefficients
      next[index] = value
      return next
    })
  }

  const line = (point: Vec2, direction: Vec2, key: string, dashed = false, color = 'var(--accent-blue)') => {
    const clipped = clipLineToRect(point, direction, BOUNDS)
    if (!clipped) return null
    return (
      <line
        key={key}
        x1={X(clipped[0][0])}
        y1={Y(clipped[0][1])}
        x2={X(clipped[1][0])}
        y2={Y(clipped[1][1])}
        stroke={color}
        strokeWidth={1.6}
        strokeDasharray={dashed ? '5 4' : undefined}
      />
    )
  }

  const axis = (center: Vec2, direction: Vec2, half: number, key: string, dashed = false) => {
    const start: Vec2 = [center[0] - direction[0] * half, center[1] - direction[1] * half]
    const end: Vec2 = [center[0] + direction[0] * half, center[1] + direction[1] * half]
    return (
      <line
        key={key}
        x1={X(start[0])}
        y1={Y(start[1])}
        x2={X(end[0])}
        y2={Y(end[1])}
        stroke="var(--chart-2)"
        strokeWidth={1.2}
        strokeDasharray={dashed ? '4 3' : undefined}
      />
    )
  }

  const centreValue = geometry.center ? pointText(geometry.center) : '—'
  const vertexValue = geometry.vertex ? pointText(geometry.vertex) : '—'
  const openingValue = geometry.opening
    ? `(${fmt(geometry.opening[0], 2)}, ${fmt(geometry.opening[1], 2)})`
    : '—'
  const axesValue =
    geometry.kind === 'parabola' && geometry.p !== undefined
      ? `p = ${fmt(geometry.p, 3)}`
      : geometry.a !== undefined && geometry.b !== undefined
        ? `${fmt(geometry.a, 3)} × ${fmt(geometry.b, 3)}`
        : '—'
  const latusRectumValue =
    geometry.kind === 'parabola' && geometry.p !== undefined ? fmt(4 * geometry.p, 4) : '—'
  const angleLabel =
    geometry.kind === 'ellipse'
      ? 'minor-axis angle'
      : geometry.kind === 'hyperbola'
        ? 'transverse-axis angle'
        : geometry.kind === 'parabola'
          ? 'squared-axis angle'
          : 'axis angle'
  const angleValue =
    geometry.theta !== undefined ? `${fmt((geometry.theta * 180) / Math.PI, 1)}°` : '—'
  const interceptsValue = `${describeIntercepts(geometry.intercepts.x, 'x')}; ${describeIntercepts(geometry.intercepts.y, 'y')}`
  const fociValue =
    geometry.foci === undefined || geometry.foci.length === 0
      ? '—'
      : geometry.kind === 'circle'
        ? `both at ${pointText(geometry.foci[0])}`
        : geometry.foci.map((focus, index) => `F${index + 1} ${pointText(focus)}`).join(', ')
  const probeValue = probePoint
    ? `${pointText(probePoint)}${probeOutside ? ' outside plot' : ''}`
    : '—'
  const focalValue =
    measurements === null
      ? '—'
      : measurements.focalSum !== null
        ? geometry.kind === 'circle'
          ? `PF = ${fmt(measurements.toFoci[0], 4)} = r`
          : `PF1 + PF2 = ${fmt(measurements.focalSum, 4)} = 2a = ${fmt(geometry.focalConstant ?? 2 * (geometry.a ?? 0), 4)}`
        : measurements.focalDifference !== null
          ? `|PF1 − PF2| = ${fmt(measurements.focalDifference, 4)} = 2a = ${fmt(geometry.focalConstant ?? 2 * (geometry.a ?? 0), 4)}`
          : measurements.focusDistance !== null
            ? `PF = ${fmt(measurements.focusDistance, 4)}`
            : '—'
  const directrixValue = geometry.directrix
    ? `through ${pointText(geometry.directrix.point)}, direction (${fmt(geometry.directrix.direction[0], 2)}, ${fmt(geometry.directrix.direction[1], 2)})`
    : '—'
  const latusRectum =
    geometry.kind === 'parabola' &&
    geometry.foci?.[0] &&
    geometry.squared &&
    geometry.p !== undefined
      ? (() => {
          const focus = geometry.foci[0]
          const half = 2 * geometry.p
          return [
            [focus[0] - geometry.squared[0] * half, focus[1] - geometry.squared[1] * half] as Vec2,
            [focus[0] + geometry.squared[0] * half, focus[1] + geometry.squared[1] * half] as Vec2,
          ]
        })()
      : null
  const directrixDistanceValue =
    measurements?.directrixDistance !== null && measurements?.directrixDistance !== undefined
      ? fmt(measurements.directrixDistance, 4)
      : '—'

  return (
    <FigureFrame
      title={title ?? 'The general equation covers every conic'}
      hint="Change A to F. B turns the axes, D and E move the centre, and a linear term is what makes a parabola. Drag the probe or use its slider to move a point along the curve."
      caption={
        caption ??
        'The general equation covers every conic. A and C set the quadratic shape, B turns the axes, and D and E move the centre. Past the balance point |B| = 2√(AC) an ellipse becomes a hyperbola; a parabola needs a linear term.'
      }
      onReset={reset}
      readout={[
        { label: 'A, B, C', value: `${fmt(A)}, ${fmt(B)}, ${fmt(C)}` },
        { label: 'D, E, F', value: `${fmt(D)}, ${fmt(E)}, ${fmt(F)}` },
        { label: 'discriminant', value: fmt(discriminant(A, B, C), 3) },
        {
          label: 'squared terms',
          value: `${(Math.abs(A) > 1e-9 ? 1 : 0) + (Math.abs(C) > 1e-9 ? 1 : 0)} written`,
        },
        { label: 'type', value: TYPE_LABELS[geometry.kind] },
        { label: 'centre', value: centreValue },
        { label: 'vertex', value: vertexValue },
        { label: 'opening', value: openingValue },
        { label: 'axes', value: axesValue },
        { label: 'latus rectum', value: latusRectumValue },
        { label: angleLabel, value: angleValue },
        { label: 'foci', value: fociValue },
        { label: 'probe', value: probeValue },
        { label: 'focal property', value: focalValue },
        { label: 'directrix', value: directrixValue },
        { label: 'distance to directrix', value: directrixDistanceValue },
        {
          label: 'eccentricity',
          value: geometry.eccentricity === null ? '—' : `${fmt(geometry.eccentricity, 3)} (constant)`,
        },
        { label: 'intercepts', value: interceptsValue },
        { label: 'certainty', value: geometry.certainty },
        ...(geometry.note ? [{ label: 'note', value: geometry.note }] : []),
      ]}
    >
      <svg
        ref={svgRef}
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        role="img"
        aria-label={`${TYPE_LABELS[geometry.kind]} from its general equation`}
      >
        <defs>
          <clipPath id={clipId}>
            <rect x={PAD} y={PAD} width={W - PAD * 2} height={H - PAD * 2} />
          </clipPath>
        </defs>

        <line x1={PAD} y1={Y(0)} x2={W - PAD} y2={Y(0)} stroke="var(--text-muted)" strokeWidth="0.8" />
        <line x1={X(0)} y1={PAD} x2={X(0)} y2={H - PAD} stroke="var(--text-muted)" strokeWidth="0.8" />

        <g clipPath={`url(#${clipId})`}>
          {segments.map((segment, index) => (
            <line
              key={index}
              x1={X(segment[0][0])}
              y1={Y(segment[0][1])}
              x2={X(segment[1][0])}
              y2={Y(segment[1][1])}
              stroke="var(--accent-blue)"
              strokeWidth={1.7}
              strokeLinecap="round"
            />
          ))}

          {geometry.kind === 'ellipse' || geometry.kind === 'circle'
            ? geometry.center &&
              geometry.major &&
              geometry.minor &&
              geometry.a !== undefined &&
              geometry.b !== undefined && (
                <>
                  {axis(geometry.center, geometry.major, geometry.a, 'major')}
                  {axis(geometry.center, geometry.minor, geometry.b, 'minor', true)}
                  <circle cx={X(geometry.center[0])} cy={Y(geometry.center[1])} r="2.5" fill="var(--chart-2)" />
                </>
              )
            : null}

          {geometry.kind === 'hyperbola' &&
          geometry.center &&
          geometry.transverse &&
          geometry.conjugate &&
          geometry.a !== undefined &&
          geometry.b !== undefined ? (
            <>
              {axis(geometry.center, geometry.transverse, geometry.a, 'transverse')}
              {axis(geometry.center, geometry.conjugate, geometry.b, 'conjugate', true)}
              {line(
                geometry.center,
                [
                  geometry.transverse[0] * geometry.a + geometry.conjugate[0] * geometry.b,
                  geometry.transverse[1] * geometry.a + geometry.conjugate[1] * geometry.b,
                ],
                'asymptote-1',
                true,
              )}
              {line(
                geometry.center,
                [
                  geometry.transverse[0] * geometry.a - geometry.conjugate[0] * geometry.b,
                  geometry.transverse[1] * geometry.a - geometry.conjugate[1] * geometry.b,
                ],
                'asymptote-2',
                true,
              )}
              <circle cx={X(geometry.center[0])} cy={Y(geometry.center[1])} r="2.5" fill="var(--chart-2)" />
            </>
          ) : null}

          {geometry.kind === 'parabola' && geometry.vertex && geometry.opening ? (
            <>
              {line(geometry.vertex, geometry.opening, 'parabola-axis', true)}
              <circle cx={X(geometry.vertex[0])} cy={Y(geometry.vertex[1])} r="2.5" fill="var(--chart-2)" />
            </>
          ) : null}

          {geometry.directrix
            ? line(geometry.directrix.point, geometry.directrix.direction, 'directrix', true, 'var(--chart-3)')
            : null}

          {latusRectum ? (
            <line
              data-part="latus-rectum"
              x1={X(latusRectum[0][0])}
              y1={Y(latusRectum[0][1])}
              x2={X(latusRectum[1][0])}
              y2={Y(latusRectum[1][1])}
              stroke="var(--chart-3)"
              strokeWidth={1.2}
            />
          ) : null}

          {geometry.lines.map((item, index) => line(item.point, item.direction, `line-${index}`))}

          {probePoint ? (
            <>
              {geometry.foci?.map((focus, index) => (
                <line
                  key={`probe-focus-${index}`}
                  x1={X(probePoint[0])}
                  y1={Y(probePoint[1])}
                  x2={X(focus[0])}
                  y2={Y(focus[1])}
                  stroke="var(--chart-2)"
                  strokeWidth={1.1}
                  strokeDasharray="3 3"
                />
              ))}
              <circle
                cx={X(probePoint[0])}
                cy={Y(probePoint[1])}
                r="4"
                fill="var(--chart-2)"
                stroke="var(--bg)"
                strokeWidth={1.5}
              />
            </>
          ) : null}

          {geometry.foci?.map((focus, index) => (
            <g key={`focus-${index}`}>
              <circle cx={X(focus[0])} cy={Y(focus[1])} r="3" fill="var(--accent-blue)" />
              <text
                x={X(focus[0]) + 5}
                y={Y(focus[1]) - 4}
                fontSize="9"
                fill="var(--accent-blue)"
                stroke="var(--bg)"
                strokeWidth="3"
                paintOrder="stroke"
              >
                {geometry.kind === 'circle' ? 'F' : `F${index + 1}`}
              </text>
            </g>
          ))}

          {geometry.kind === 'point' && geometry.center ? (
            <circle cx={X(geometry.center[0])} cy={Y(geometry.center[1])} r="4" fill="var(--accent-blue)" />
          ) : null}
        </g>

        <XAxisTicks ticks={ticks} scale={X} at={Y(0)} />
        <YAxisTicks ticks={ticks} scale={Y} at={X(0)} />

        <text x={PAD + 2} y={PAD + 10} fontSize="10" fill="var(--accent-blue)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
          {TYPE_LABELS[geometry.kind]}
        </text>
        {geometry.kind === 'empty' || geometry.kind === 'all-plane' || geometry.kind === 'invalid' ? (
          <text x={W / 2} y={H / 2} textAnchor="middle" fontSize="10" fill="var(--text-muted)">
            {geometry.kind === 'all-plane' ? 'every point satisfies the equation' : geometry.kind === 'invalid' ? 'coefficients must be finite' : 'no real points'}
          </text>
        ) : null}
      </svg>

      {curve && probeRange ? (
        <div className="mt-3 space-y-2">
          <div className="flex items-center justify-between gap-2">
            <label htmlFor={`${clipId}-probe`} className="font-mono text-[10px] text-mathua-muted">
              Probe position
            </label>
            {geometry.kind === 'hyperbola' ? (
              <span className="font-mono text-[10px] text-mathua-muted">branch {branch === 1 ? '1' : '2'}</span>
            ) : null}
          </div>
          <input
            id={`${clipId}-probe`}
            type="range"
            aria-label="Probe position"
            min={probeRange.min}
            max={probeRange.max}
            step={probeRange.step}
            value={probeParameter}
            onChange={(event) => setRawProbe(Number(event.target.value))}
            className="w-full"
            style={{ accentColor: 'var(--chart-2)' }}
          />
          {geometry.kind === 'hyperbola' ? (
            <div role="group" aria-label="Hyperbola branch" className="flex items-center gap-2">
              <span className="font-mono text-[10px] text-mathua-muted">branch</span>
              <button
                type="button"
                aria-pressed={branch === 1}
                onClick={() => setBranch(1)}
                className="border border-mathua-border px-2 py-1 font-mono text-[10px] text-mathua-secondary aria-pressed:border-mathua-blue aria-pressed:text-mathua-blue"
              >
                First branch
              </button>
              <button
                type="button"
                aria-pressed={branch === -1}
                onClick={() => setBranch(-1)}
                className="border border-mathua-border px-2 py-1 font-mono text-[10px] text-mathua-secondary aria-pressed:border-mathua-blue aria-pressed:text-mathua-blue"
              >
                Second branch
              </button>
            </div>
          ) : null}
        </div>
      ) : (
        <p className="mt-3 font-mono text-[10px] text-mathua-muted">
          Focal probe unavailable for this equation.
        </p>
      )}

      <div className="mt-3 grid grid-cols-1 gap-2">
        {(['A', 'B', 'C', 'D', 'E', 'F'] as const).map((label, index) => (
          <Slider
            key={label}
            label={label}
            value={coefficients[index]}
            min={index < 3 ? -2 : -4}
            max={index < 3 ? 2 : 4}
            step={0.05}
            onChange={(value) => setCoefficient(index, value)}
          />
        ))}
      </div>

      <div className="mt-2 flex items-center gap-1.5">
        <span className="font-mono text-[10px] text-mathua-muted">Preset</span>
        <select
          aria-label="Preset equation"
          value={selected?.label ?? ''}
          onChange={(event) => selectPreset(event.target.value)}
          className="h-7 border border-mathua-border bg-transparent px-1 font-mono text-[10px] text-mathua-secondary"
        >
          {!selected && <option value="">Custom</option>}
          {PRESETS.map((preset) => (
            <option key={preset.label} value={preset.label}>
              {preset.label}
            </option>
          ))}
        </select>
      </div>
    </FigureFrame>
  )
}
