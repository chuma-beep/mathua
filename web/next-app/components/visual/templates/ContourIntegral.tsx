'use client'

import { useMemo, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Select, Slider, Toggle } from '../Controls'
import { contourIntegralSchema } from '../schemas'
import {
  FnLabels,
  fnValue,
  integrate,
  pathAt,
  poleInfo,
  type Complex,
  type ContourFn,
  type Path,
} from '../math/contourIntegral'
import { scaleLinear } from '../scale'
import { fmt } from '../format'
import type { TemplateProps } from '../types'

const PANEL = 300
const PANEL_H = 220
const PAD = 24

const PATH_OPTIONS = [
  { value: 'circle' as const, label: 'Circle' },
  { value: 'ellipse' as const, label: 'Ellipse' },
  { value: 'wavy' as const, label: 'Wavy closed curve' },
  { value: 'line' as const, label: 'Line segment' },
]

const FN_OPTIONS = (Object.keys(FnLabels) as ContourFn[]).map((key) => ({ value: key, label: FnLabels[key] }))

type PathKind = 'circle' | 'ellipse' | 'wavy' | 'line'

function formatComplex(value: Complex | null): string {
  if (!value || !Number.isFinite(value[0]) || !Number.isFinite(value[1])) return 'undefined'
  return `${fmt(value[0])} ${value[1] < 0 ? '−' : '+'} ${fmt(Math.abs(value[1]))}i`
}

export default function ContourIntegral({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => contourIntegralSchema.parse(params), [params])
  const [kind, setKind] = useState<PathKind>(defaults.path)
  const [fn, setFn] = useState<ContourFn>(defaults.fn)
  const [reverse, setReverse] = useState(defaults.reverse)
  const [radius, setRadius] = useState(defaults.radius)
  const [centerX, setCenterX] = useState(defaults.centerX)
  const [t, setT] = useState(defaults.t)

  const direction: 1 | -1 = reverse ? -1 : 1

  const path = useMemo<Path>(() => {
    const center: Complex = [centerX, 0]
    switch (kind) {
      case 'circle':
        return { kind: 'circle', center, radius }
      case 'ellipse':
        return { kind: 'ellipse', center, rx: radius * 1.5, ry: radius }
      case 'wavy':
        return { kind: 'wavy', center, radius, amplitude: radius * 0.22, lobes: 3 }
      case 'line':
        return { kind: 'line', a: [centerX - radius, -radius], b: [centerX + radius, radius] }
    }
  }, [kind, centerX, radius])

  // The integral does not depend on t, so it is computed once per path and
  // function; t only chooses how much of the accumulated sequence to show.
  const { partial, total } = useMemo(() => integrate(path, fn, direction, 2000), [path, fn, direction])
  const pole = useMemo(() => poleInfo(path), [path])
  const divergent = !Number.isFinite(total[0]) || !Number.isFinite(total[1])

  const progress = Math.max(0, Math.min(1, t))
  const index = Math.max(0, Math.ceil(progress * Math.max(1, partial.length)) - 1)
  const running = useMemo<Complex | null>(
    () => (divergent ? null : progress >= 1 ? total : (partial[index] ?? ([0, 0] as Complex))),
    [divergent, progress, total, partial, index],
  )

  const pathPoints = useMemo(() => {
    const out: Complex[] = []
    for (let i = 0; i <= 180; i++) out.push(pathAt(path, 1, i / 180).z)
    return out
  }, [path])

  const imagePoints = useMemo(() => {
    const out: Array<Complex | null> = []
    for (let i = 0; i <= 240; i++) out.push(fnValue(fn, pathAt(path, direction, i / 240).z))
    return out
  }, [path, fn, direction])

  const zExtent = useMemo(() => {
    const max = Math.max(1.2, ...pathPoints.map(([x, y]) => Math.max(Math.abs(x), Math.abs(y))))
    return max * 1.2
  }, [pathPoints])

  const wExtent = useMemo(() => {
    const values = imagePoints.filter((v): v is Complex => v !== null)
    if (running) values.push(running)
    if (!divergent) values.push(total)
    const max = Math.max(1, ...values.map(([x, y]) => Math.max(Math.abs(x), Math.abs(y))))
    return Math.min(9, max * 1.2)
  }, [imagePoints, running, total, divergent])

  const zX = scaleLinear(-zExtent, zExtent, PAD, PANEL - PAD)
  const zY = scaleLinear(zExtent, -zExtent, PAD, PANEL_H - PAD)
  const wX = scaleLinear(-wExtent, wExtent, PAD, PANEL - PAD)
  const wY = scaleLinear(wExtent, -wExtent, PAD, PANEL_H - PAD)

  const sample = pathAt(path, direction, progress).z
  const tip = { x: zX(sample[0]), y: zY(sample[1]) }
  const next = pathAt(path, direction, Math.min(1, progress + 0.02)).z
  const arrowAngle = (Math.atan2(zY(next[1]) - tip.y, zX(next[0]) - tip.x) * 180) / Math.PI

  const reset = () => {
    setKind(defaults.path)
    setFn(defaults.fn)
    setReverse(defaults.reverse)
    setRadius(defaults.radius)
    setCenterX(defaults.centerX)
    setT(defaults.t)
  }

  const trace = divergent ? [] : partial.slice(0, index + 1)

  return (
    <FigureFrame
      title={title ?? 'A contour integral accumulates in the plane'}
      hint="Sweep t. The vector in the right-hand plane is the integral so far, not an area."
      caption={
        caption ??
        'The left plane holds the path and its direction; the right holds the image of the path and the running value of the integral. A pole on the path makes the integral diverge. A segment and a circle are solved exactly; an ellipse and a wavy curve are minimised numerically, and the detection row says which.'
      }
      onReset={reset}
      readout={[
        { label: 'I(t)', value: formatComplex(running) },
        { label: '|I(t)|', value: running ? fmt(Math.hypot(running[0], running[1])) : 'undefined' },
        { label: 'total', value: formatComplex(total) },
        { label: 't', value: fmt(progress, 2) },
        { label: 'pole', value: pole.crosses ? 'on the path' : pole.certainty === 'uncertain' ? 'not found (unverified)' : 'clear' },
        // A near miss of 1e-4 must not print as 0 beside a "clear" verdict, which
        // is the same mistake the basis figure made with a rounded determinant.
        { label: 'closest', value: pole.closest >= 1e-3 ? fmt(pole.closest, 3) : pole.closest.toExponential(1) },
        { label: 'detection', value: pole.certainty },
      ]}
    >
      <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
        <svg viewBox={`0 0 ${PANEL} ${PANEL_H}`} className="block h-auto w-full font-mono" role="img" aria-label="The path and its direction in the z plane">
          <line x1={PAD} y1={zY(0)} x2={PANEL - PAD} y2={zY(0)} stroke="var(--text-muted)" strokeWidth="0.8" />
          <line x1={zX(0)} y1={PAD} x2={zX(0)} y2={PANEL_H - PAD} stroke="var(--text-muted)" strokeWidth="0.8" />

          <polyline
            points={pathPoints.map(([x, y]) => `${zX(x)},${zY(y)}`).join(' ')}
            fill="none"
            stroke="var(--accent-blue)"
            strokeWidth="1.8"
          />

          <circle cx={zX(0)} cy={zY(0)} r="4" fill="none" stroke="var(--accent-red)" strokeWidth="1.4" />
          <text x={zX(0) + 6} y={zY(0) - 6} fontSize="9" fill="var(--accent-red)">0</text>

          <circle cx={tip.x} cy={tip.y} r="4.5" fill="var(--accent-green)" />
          <path
            d="M -4 2.4 L 4 0 L -4 -2.4 Z"
            transform={`translate(${tip.x} ${tip.y}) rotate(${arrowAngle})`}
            fill="var(--accent-green)"
          />
          <text x={tip.x + 7} y={tip.y + 12} fontSize="9" fill="var(--accent-green)">z(t)</text>
          <text x={PAD + 2} y={PAD + 10} fontSize="9" fill="var(--text-muted)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
            path, direction {direction === 1 ? '→' : '←'}
          </text>
        </svg>

        <svg viewBox={`0 0 ${PANEL} ${PANEL_H}`} className="block h-auto w-full font-mono" role="img" aria-label="The image of the path and the accumulated integral in the w plane">
          <line x1={PAD} y1={wY(0)} x2={PANEL - PAD} y2={wY(0)} stroke="var(--text-muted)" strokeWidth="0.8" />
          <line x1={wX(0)} y1={PAD} x2={wX(0)} y2={PANEL_H - PAD} stroke="var(--text-muted)" strokeWidth="0.8" />

          <polyline
            points={imagePoints
              .filter((v): v is Complex => v !== null)
              .map(([x, y]) => `${wX(x)},${wY(y)}`)
              .join(' ')}
            fill="none"
            stroke="var(--text-muted)"
            strokeWidth="1.2"
            strokeDasharray="3 3"
            opacity={divergent ? 0.25 : 0.8}
          />

          {!divergent && trace.length > 1 && (
            <polyline
              points={trace.map(([x, y]) => `${wX(x)},${wY(y)}`).join(' ')}
              fill="none"
              stroke="var(--accent-teal)"
              strokeWidth="1"
              opacity="0.55"
            />
          )}

          {running && (
            <>
              <line
                x1={wX(0)}
                y1={wY(0)}
                x2={wX(running[0])}
                y2={wY(running[1])}
                stroke="var(--accent-blue)"
                strokeWidth="2"
                markerEnd="url(#contour-arrow)"
              />
              <circle cx={wX(running[0])} cy={wY(running[1])} r="3.6" fill="var(--accent-blue)" />
              <text x={wX(running[0]) + 6} y={wY(running[1]) - 6} fontSize="9" fill="var(--accent-blue)">
                I(t)
              </text>
            </>
          )}

          {divergent && (
            <text x={PANEL / 2} y={PANEL_H / 2} textAnchor="middle" fontSize="10" fill="var(--accent-red)">
              path passes through 0: no value
            </text>
          )}

          <defs>
            <marker id="contour-arrow" markerWidth="6" markerHeight="6" refX="5" refY="3" orient="auto">
              <path d="M0,0 L6,3 L0,6 Z" fill="var(--accent-blue)" />
            </marker>
          </defs>
        </svg>
      </div>

      <div className="mt-3 space-y-2">
        <div className="grid grid-cols-1 gap-2 sm:grid-cols-2">
          <Select label="Path" value={kind} options={PATH_OPTIONS} onChange={(value) => setKind(value as PathKind)} />
          <Select label="Function" value={fn} options={FN_OPTIONS} onChange={(value) => setFn(value as ContourFn)} />
        </div>
        <Slider label="t" value={progress} min={0} max={1} step={0.005} onChange={setT} />
        <Slider label="size" value={radius} min={0.3} max={2.5} step={0.05} onChange={setRadius} />
        <Slider label="centre" value={centerX} min={-2} max={2} step={0.05} onChange={setCenterX} />
        <Toggle label="Traverse backwards" checked={reverse} onChange={setReverse} />
      </div>
    </FigureFrame>
  )
}
