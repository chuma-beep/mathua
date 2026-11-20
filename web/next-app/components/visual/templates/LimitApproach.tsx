'use client'

import { useMemo, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Select, Slider } from '../Controls'
import { XAxisTicks } from '../Axis'
import { limitSchema } from '../schemas'
import { continuousAt, limitExists, limitValue, oneSidedLimits, valueAtPoint, type LimitKind } from '../math/limit'
import { pathFromPoints, scaleLinear } from '../scale'
import { fmt } from '../format'
import type { TemplateProps } from '../types'

const W = 340
const H = 250
const PAD = 28
const CAP = 8

const KIND_OPTIONS: Array<{ value: LimitKind; label: string }> = [
  { value: 'removable', label: 'Hole' },
  { value: 'jump', label: 'Jump' },
  { value: 'vertical_asymptote', label: 'Asymptote' },
]

function limitText(value: number): string {
  if (value === Number.POSITIVE_INFINITY) return '+∞'
  if (value === Number.NEGATIVE_INFINITY) return '−∞'
  if (!Number.isFinite(value)) return 'undefined'
  return fmt(value)
}

export default function LimitApproach({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => limitSchema.parse(params), [params])
  const [kind, setKind] = useState<LimitKind>(defaults.kind)
  const [a, setA] = useState(defaults.a)
  const [delta, setDelta] = useState(defaults.delta)

  const domain: [number, number] = [a - 1.6, a + 1.6]
  const sides = oneSidedLimits(kind, a)
  const exists = limitExists(kind, a)
  const continuous = continuousAt(kind, a)

  // Two paths, one either side of a, because the point at a is exactly where the
  // question is and a single path would draw straight over it.
  const paths = useMemo(() => {
    const left: Array<[number, number]> = []
    const right: Array<[number, number]> = []
    const gap = 0.012
    const lo = a - 1.6
    const hi = a + 1.6
    for (let i = 0; i <= 320; i++) {
      const x = lo + ((hi - lo) * i) / 320
      if (Math.abs(x - a) < gap) continue
      const y = limitValue(kind, a, x)
      if (!Number.isFinite(y) || Math.abs(y) > CAP) continue
      ;(x < a ? left : right).push([x, y])
    }
    return { left, right }
  }, [kind, a])

  const X = scaleLinear(domain[0], domain[1], PAD, W - PAD)
  // The visible range follows the values that were actually drawn, so a figure
  // about a hole near y = 4 does not spend half its height on empty axes.
  const drawn = [...paths.left, ...paths.right].map(([, y]) => y).concat([0])
  const lo = Math.min(...drawn)
  const hi = Math.max(...drawn)
  const padding = Math.max(0.6, (hi - lo) * 0.18)
  const Y = scaleLinear(Math.max(-CAP - 1, lo - padding), Math.min(CAP + 1, hi + padding), H - PAD, PAD)

  const reset = () => {
    setKind(defaults.kind)
    setA(defaults.a)
    setDelta(defaults.delta)
  }

  return (
    <FigureFrame
      title={title ?? 'Two sides, one limit'}
      hint="Narrow δ and watch the two sides. Whether they meet is the question."
      caption={
        caption ??
        'The shaded band is |x - a| < δ. A limit exists when both sides approach the same value as the band narrows, which a hole allows and a jump prevents.'
      }
      onReset={reset}
      readout={[
        { label: 'left limit', value: limitText(sides.left) },
        { label: 'right limit', value: limitText(sides.right) },
        { label: 'f(a)', value: limitText(valueAtPoint(kind, a)) },
        { label: 'two-sided', value: exists ? 'exists' : 'does not exist' },
        { label: 'continuous', value: continuous ? 'yes' : 'no' },
        { label: 'δ', value: fmt(delta) },
      ]}
    >
      <svg
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        role="img"
        aria-label="A function near a point, with a moving band around that point"
      >
        <rect x={X(a - delta)} y={PAD} width={Math.max(0, X(a + delta) - X(a - delta))} height={H - PAD * 2} fill="var(--accent-blue)" fillOpacity="0.09" />
        <line x1={PAD} y1={Y(0)} x2={W - PAD} y2={Y(0)} stroke="var(--text-muted)" strokeWidth="0.8" />

        <XAxisTicks ticks={[a - 1, a, a + 1]} scale={X} at={Y(0)} format={(v) => fmt(v)} />
        {exists && Number.isFinite(sides.left) && (
          <line x1={PAD} y1={Y(sides.left)} x2={W - PAD} y2={Y(sides.left)} stroke="var(--accent-green)" strokeWidth="1" strokeDasharray="4 4" />
        )}
        {Number.isFinite(sides.left) && Number.isFinite(sides.right) && Math.abs(sides.left - sides.right) > 1e-12 && (
          <>
            <line x1={PAD} y1={Y(sides.left)} x2={W - PAD} y2={Y(sides.left)} stroke="var(--accent-teal)" strokeWidth="1" strokeDasharray="4 4" />
            <line x1={PAD} y1={Y(sides.right)} x2={W - PAD} y2={Y(sides.right)} stroke="var(--accent-teal)" strokeWidth="1" strokeDasharray="4 4" />
          </>
        )}

        <line x1={X(a)} y1={PAD} x2={X(a)} y2={H - PAD} stroke="var(--text-muted)" strokeWidth="0.9" strokeDasharray="2 3" />

        <path d={pathFromPoints(paths.left.map(([x, y]) => [X(x), Y(y)]))} fill="none" stroke="var(--text-primary)" strokeWidth="1.8" />
        <path d={pathFromPoints(paths.right.map(([x, y]) => [X(x), Y(y)]))} fill="none" stroke="var(--text-primary)" strokeWidth="1.8" />

        {kind === 'removable' && <circle cx={X(a)} cy={Y(sides.left)} r="4" fill="var(--bg)" stroke="var(--accent-blue)" strokeWidth="1.6" />}
        {kind === 'jump' && (
          <>
            <circle cx={X(a)} cy={Y(sides.left)} r="3.6" fill="var(--bg)" stroke="var(--accent-teal)" strokeWidth="1.6" />
            <circle cx={X(a)} cy={Y(sides.right)} r="3.6" fill="var(--accent-teal)" />
          </>
        )}
        {kind === 'vertical_asymptote' && (
          <line x1={X(a)} y1={PAD} x2={X(a)} y2={H - PAD} stroke="var(--accent-blue)" strokeWidth="1.4" />
        )}

        <text x={PAD + 2} y={PAD + 10} fontSize="10" fill={exists ? 'var(--accent-green)' : 'var(--accent-teal)'} stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
          {exists ? 'both sides converge' : kind === 'jump' ? 'jump: no limit' : 'unbounded'}
        </text>
      </svg>

      <div className="mt-3 space-y-2">
        <Select label="Shape" value={kind} options={KIND_OPTIONS} onChange={(v) => setKind(v as LimitKind)} />
        <Slider label="a" value={a} min={-2} max={2} step={0.1} onChange={setA} />
        <Slider label="δ" value={delta} min={0.05} max={1.5} step={0.01} onChange={setDelta} />
      </div>
    </FigureFrame>
  )
}
