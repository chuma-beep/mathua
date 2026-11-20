'use client'

import { useCallback, useMemo, useRef, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { NumberField, Select, Toggle } from '../Controls'
import { XAxisTicks, YAxisTicks } from '../Axis'
import { complexFieldSchema } from '../schemas'
import { CF_LABELS, cfCritical, cfDerivative, cfGradients, cfLevelValues, cfPoles, cfValue, cfZeros, type ComplexFn, type Vec2 } from '../math/complexField'
import { contourSegments } from '../math/contour'
import { clamp, scaleLinear } from '../scale'
import { dot, length } from '../math/matrix'
import { fmt } from '../format'
import { useDrag } from '../useDrag'
import type { TemplateProps } from '../types'

const SIZE = 320
const PAD = 30
const EXTENT = 2.4
const DEG = 180 / Math.PI

const FN_OPTIONS: Array<{ value: ComplexFn; label: string }> = (Object.keys(CF_LABELS) as ComplexFn[]).map((fn) => ({
  value: fn,
  label: CF_LABELS[fn],
}))

function fmtComplex(z: Vec2 | null): string {
  if (!z) return 'undefined'
  const sign = z[1] < 0 ? '−' : '+'
  return `${fmt(z[0])} ${sign} ${fmt(Math.abs(z[1]))}i`
}

export default function ComplexLevels({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => {
    const p = complexFieldSchema.parse(params)
    return {
      fn: p.fn,
      // The parsed tuple arrives optional-element typed; build it once so state
      // keeps the Vec2 shape the math module expects.
      probe: [Number(p.probe[0]), Number(p.probe[1])] as Vec2,
      levels: p.levels,
      showImag: p.showImag,
    }
  }, [params])
  const [fn, setFn] = useState<ComplexFn>(defaults.fn)
  const [probe, setProbe] = useState<Vec2>(defaults.probe)
  const [showImag, setShowImag] = useState(defaults.showImag)
  const svgRef = useRef<SVGSVGElement | null>(null)

  const X = scaleLinear(-EXTENT, EXTENT, PAD, SIZE - PAD)
  const Y = scaleLinear(EXTENT, -EXTENT, PAD, SIZE - PAD)
  const invX = scaleLinear(PAD, SIZE - PAD, -EXTENT, EXTENT)
  const invY = scaleLinear(SIZE - PAD, PAD, -EXTENT, EXTENT)

  const realCurves = useMemo(() => {
    // Built inside the memo: an object literal outside it is a new value on every
    // render, so the dependency would never be stable.
    const bounds = { x0: -EXTENT, x1: EXTENT, y0: -EXTENT, y1: EXTENT }
    return cfLevelValues(fn, 'u', EXTENT, defaults.levels).map((level) => ({
      level,
      segments: contourSegments((x, y) => cfValue(fn, [x, y])?.[0] ?? Number.NaN, bounds, level, 80, 80),
    }))
  }, [fn, defaults.levels])

  const imagCurves = useMemo(() => {
    const bounds = { x0: -EXTENT, x1: EXTENT, y0: -EXTENT, y1: EXTENT }
    return cfLevelValues(fn, 'v', EXTENT, defaults.levels).map((level) => ({
      level,
      segments: contourSegments((x, y) => cfValue(fn, [x, y])?.[1] ?? Number.NaN, bounds, level, 80, 80),
    }))
  }, [fn, defaults.levels])

  const onDrag = useCallback(
    (point: { x: number; y: number }) => {
      setProbe([clamp(invX(point.x), -EXTENT, EXTENT), clamp(invY(point.y), -EXTENT, EXTENT)])
    },
    [invX, invY],
  )
  useDrag(svgRef, onDrag)

  const reset = () => {
    setFn(defaults.fn)
    setProbe(defaults.probe)
    setShowImag(defaults.showImag)
  }

  const value = cfValue(fn, probe)
  const derivative = cfDerivative(fn, probe)
  const gradients = cfGradients(fn, probe)
  // Three different reasons a number may be missing, kept apart because they
  // mean different things: no value at all (a pole), a value that is zero (the
  // argument is undefined), and a vanishing derivative (no tangent).
  const atPole = value === null
  const atZero = value !== null && Math.hypot(value[0], value[1]) < 1e-9
  const atCritical = derivative !== null && Math.hypot(derivative[0], derivative[1]) < 1e-9

  // Orthogonality is a statement about two tangents. Where either gradient
  // vanishes there is no tangent, so the honest answer is "undefined".
  const angle =
    gradients && !atCritical
      ? Math.acos(clamp(dot(gradients.du, gradients.dv) / (length(gradients.du) * length(gradients.dv)), -1, 1)) * DEG
      : null

  // "Fails" is a claim about a function that has values. At a pole there is
  // nothing to test, so the honest answer is undefined; every function in this
  // family is analytic wherever it is defined.
  const cauchyRiemann = gradients
    ? Math.abs(gradients.du[0] - gradients.dv[1]) < 1e-9 && Math.abs(gradients.du[1] + gradients.dv[0]) < 1e-9
      ? 'satisfied'
      : 'fails'
    : 'undefined'

  const drawCurves = (curves: typeof realCurves, color: string) =>
    curves.map(({ level, segments }) =>
      segments.map(([[x1, y1], [x2, y2]], i) => (
        <line key={`${level}-${i}`} x1={X(x1)} y1={Y(y1)} x2={X(x2)} y2={Y(y2)} stroke={color} strokeWidth="1.1" opacity="0.55" />
      )),
    )

  return (
    <FigureFrame
      title={title ?? 'Level curves of a complex function'}
      hint="Drag the probe. The blue and teal families cross at a right angle wherever the derivative is not zero."
      caption={
        caption ??
        'The blue curves are where Re f is constant and the teal where Im f is. Cauchy-Riemann makes the two families perpendicular, except at a critical point, where neither has a tangent.'
      }
      onReset={reset}
      readout={[
        { label: 'f(z)', value: fmtComplex(value) },
        { label: '|f|', value: value ? fmt(Math.hypot(value[0], value[1])) : '—' },
        { label: 'arg f', value: value && !atZero ? fmt(Math.atan2(value[1], value[0])) : 'undefined' },
        { label: "|f′|", value: derivative ? fmt(Math.hypot(derivative[0], derivative[1])) : '—' },
        { label: 'Cauchy-Riemann', value: cauchyRiemann },
        { label: 'angle', value: angle === null ? 'undefined' : `${fmt(angle, 1)}°` },
      ]}
    >
      <svg
        ref={svgRef}
        viewBox={`0 0 ${SIZE} ${SIZE}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        style={{ touchAction: 'none' }}
        role="img"
        aria-label="Level curves of the real and imaginary parts of a complex function"
      >
        <line x1={PAD} y1={Y(0)} x2={SIZE - PAD} y2={Y(0)} stroke="var(--text-muted)" strokeWidth="0.8" />
        <line x1={X(0)} y1={PAD} x2={X(0)} y2={SIZE - PAD} stroke="var(--text-muted)" strokeWidth="0.8" />

        {drawCurves(realCurves, 'var(--accent-blue)')}
        {showImag && drawCurves(imagCurves, 'var(--accent-teal)')}

        {cfZeros(fn)
          .filter(([x, y]) => Math.abs(x) <= EXTENT && Math.abs(y) <= EXTENT)
          .map(([x, y], i) => (
            <circle key={`z${i}`} cx={X(x)} cy={Y(y)} r="4" fill="var(--bg)" stroke="var(--accent-green)" strokeWidth="1.6" />
          ))}

        {cfPoles(fn)
          .filter(([x, y]) => Math.abs(x) <= EXTENT && Math.abs(y) <= EXTENT)
          .map(([x, y], i) => (
            <circle key={`p${i}`} cx={X(x)} cy={Y(y)} r="4" fill="var(--bg)" stroke="var(--accent-red)" strokeWidth="1.6" />
          ))}

        {cfCritical(fn)
          .filter(([x, y]) => Math.abs(x) <= EXTENT && Math.abs(y) <= EXTENT)
          .map(([x, y], i) => (
            <g key={`c${i}`} stroke="var(--accent-green)" strokeWidth="1.4">
              <line x1={X(x) - 4} y1={Y(y) - 4} x2={X(x) + 4} y2={Y(y) + 4} />
              <line x1={X(x) - 4} y1={Y(y) + 4} x2={X(x) + 4} y2={Y(y) - 4} />
            </g>
          ))}

        {gradients && !atCritical && (
          <>
            <line
              x1={X(probe[0])}
              y1={Y(probe[1])}
              x2={X(probe[0] + gradients.du[0] * 0.18)}
              y2={Y(probe[1] + gradients.du[1] * 0.18)}
              stroke="var(--accent-blue)"
              strokeWidth="1.6"
              markerEnd="url(#levels-arrow-u)"
            />
            <line
              x1={X(probe[0])}
              y1={Y(probe[1])}
              x2={X(probe[0] + gradients.dv[0] * 0.18)}
              y2={Y(probe[1] + gradients.dv[1] * 0.18)}
              stroke="var(--accent-teal)"
              strokeWidth="1.6"
              markerEnd="url(#levels-arrow-v)"
            />
          </>
        )}

        <defs>
          {[['u', 'var(--accent-blue)'], ['v', 'var(--accent-teal)']].map(([key, color]) => (
            <marker key={key} id={`levels-arrow-${key}`} markerWidth="6" markerHeight="6" refX="5" refY="3" orient="auto">
              <path d="M0,0 L6,3 L0,6 Z" fill={color} />
            </marker>
          ))}
        </defs>

        <circle cx={X(probe[0])} cy={Y(probe[1])} r="4.2" fill="var(--text-primary)" />
        <text x={X(probe[0]) + 7} y={Y(probe[1]) - 7} fontSize="10" fill="var(--text-primary)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
          z
        </text>

        <text x={PAD} y={SIZE - 8} fontSize="9" fill="var(--accent-blue)">
          Re f
        </text>
        {showImag && (
          <text x={PAD + 34} y={SIZE - 8} fontSize="9" fill="var(--accent-teal)">
            Im f
          </text>
        )}
        {atPole ? (
          <text x={SIZE - PAD} y={PAD + 2} textAnchor="end" fontSize="9" fill="var(--accent-red)">
            pole: no value here
          </text>
        ) : atCritical ? (
          <text x={SIZE - PAD} y={PAD + 2} textAnchor="end" fontSize="9" fill="var(--text-muted)">
            critical point: no tangent
          </text>
        ) : null}

        <XAxisTicks ticks={[-2, -1, 1, 2]} scale={X} at={Y(0)} />
        <YAxisTicks ticks={[-2, -1, 1, 2]} scale={Y} at={X(0)} />
      </svg>

      <div className="mt-3 space-y-2">
        <Select label="Function" value={fn} options={FN_OPTIONS} onChange={(value) => setFn(value as ComplexFn)} />
        <div className="grid grid-cols-2 gap-2">
          <NumberField label="x" value={probe[0]} min={-EXTENT} max={EXTENT} onChange={(v) => setProbe([v, probe[1]])} />
          <NumberField label="y" value={probe[1]} min={-EXTENT} max={EXTENT} onChange={(v) => setProbe([probe[0], v])} />
        </div>
        <Toggle label="Show Im f" checked={showImag} onChange={setShowImag} />
      </div>
    </FigureFrame>
  )
}
