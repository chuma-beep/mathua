'use client'

import { useMemo, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Slider } from '../Controls'
import { XAxisTicks } from '../Axis'
import { seriesSumSchema } from '../schemas'
import {
  SUM_LABELS,
  SUM_NAMES,
  partialSum,
  ratioTerm,
  seriesLimit,
  sumStatus,
  sumTerm,
  type SumKind,
} from '../math/series'
import { scaleLinear } from '../scale'
import { fmt } from '../format'
import type { TemplateProps } from '../types'

const W = 360
const H = 350
const PAD = 34
const TERM_TOP = 42
const TERM_BOTTOM = 156
const SUM_TOP = 196
const SUM_BOTTOM = 310

export default function SeriesLab({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => seriesSumSchema.parse(params), [params])
  const [kind, setKind] = useState<SumKind>(defaults.kind)
  const [count, setCount] = useState(defaults.terms)

  const terms = useMemo(
    () => Array.from({ length: count }, (_, index) => sumTerm(kind, index + 1)),
    [kind, count],
  )
  const sums = useMemo(
    () => terms.map((_, index) => partialSum(kind, index + 1)),
    [kind, terms],
  )
  const limit = seriesLimit(kind)
  const status = sumStatus(kind)
  const total = partialSum(kind, count)
  const ratio = ratioTerm(kind, count)
  const maxTerm = Math.max(0.2, ...terms.map((term) => Math.abs(term)))

  const X = scaleLinear(1, count + 1, PAD, W - PAD)
  const TY = scaleLinear(-maxTerm * 1.15, maxTerm * 1.15, TERM_BOTTOM, TERM_TOP)
  const upper = Number.isFinite(limit) ? Math.max(limit * 1.1, Math.max(...sums) * 1.05) : Math.max(...sums) * 1.05
  const SY = scaleLinear(0, upper, SUM_BOTTOM, SUM_TOP)
  const barWidth = Math.max(1.5, (W - 2 * PAD) / count - 1)

  const reset = () => {
    setKind(defaults.kind)
    setCount(defaults.terms)
  }

  return (
    <FigureFrame
      title={title ?? 'The terms decide, the partial sums show'}
      hint="Add terms. A convergent series settles near its dashed limit. A divergent one can look flat for a while and then keep climbing, so read the ratio beside the curve."
      caption={
        caption ??
        'The upper bars are the terms, with the sign carried by direction. The lower curve is the running sum after each term. A convergent series heads for the dashed line. The harmonic sums look flat over a short window and still grow without bound, so the test decides, not the eye.'
      }
      onReset={reset}
      readout={[
        { label: 'series', value: SUM_LABELS[kind] },
        { label: 'status', value: status },
        { label: 'limit', value: Number.isFinite(limit) ? fmt(limit, 5) : 'no finite limit' },
        { label: 'terms shown', value: String(count) },
        { label: 'last term', value: fmt(terms[terms.length - 1], 5) },
        { label: 'partial sum', value: fmt(total, 5) },
        {
          label: 'gap to limit',
          value: Number.isFinite(limit) ? fmt(Math.abs(limit - total), 5) : '—',
        },
        { label: 'ratio |aₙ₊₁/aₙ|', value: fmt(ratio, 4) },
      ]}
    >
      <svg
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        role="img"
        aria-label={`Terms and partial sums of the ${SUM_NAMES[kind]} series over ${count} terms`}
      >
        <text x={PAD} y={TERM_TOP - 8} fontSize="9" fill="var(--text-muted)">
          terms aₖ
        </text>
        <line x1={PAD} x2={W - PAD} y1={TY(0)} y2={TY(0)} stroke="var(--text-muted)" strokeWidth="0.8" />
        <XAxisTicks ticks={[1, Math.round(count / 2), count]} scale={X} at={SUM_BOTTOM} />
        {terms.map((term, index) => {
          const top = Math.min(TY(0), TY(term))
          const height = Math.abs(TY(term) - TY(0))
          return (
            <rect
              key={index}
              x={X(index + 1) - barWidth / 2}
              y={top}
              width={barWidth}
              height={Math.max(0.5, height)}
              fill={term >= 0 ? 'var(--accent-blue)' : 'var(--accent-red)'}
              fillOpacity="0.45"
              stroke={term >= 0 ? 'var(--accent-blue)' : 'var(--accent-red)'}
              strokeWidth="0.5"
            />
          )
        })}

        <text x={PAD} y={SUM_TOP - 8} fontSize="9" fill="var(--text-muted)">
          partial sums Sₙ
        </text>
        <line x1={PAD} x2={W - PAD} y1={SY(0)} y2={SY(0)} stroke="var(--text-muted)" strokeWidth="0.8" />
        {Number.isFinite(limit) && limit <= upper ? (
          <>
            <line
              x1={PAD}
              x2={W - PAD}
              y1={SY(limit)}
              y2={SY(limit)}
              stroke="var(--accent-teal)"
              strokeWidth="1"
              strokeDasharray="4 3"
            />
            <text x={W - PAD} y={SY(limit) - 4} textAnchor="end" fontSize="8" fill="var(--accent-teal)">
              limit {fmt(limit, 4)}
            </text>
          </>
        ) : null}
        <polyline
          points={sums.map((value, index) => `${X(index + 1)},${SY(value)}`).join(' ')}
          fill="none"
          stroke="var(--accent-blue)"
          strokeWidth="1.8"
        />
        {sums.map((value, index) => (
          <circle key={index} cx={X(index + 1)} cy={SY(value)} r="1.8" fill="var(--accent-blue)" />
        ))}
      </svg>

      <div className="mt-3 grid grid-cols-1 gap-2">
        <div className="flex items-center gap-2">
          <label htmlFor="series-kind" className="w-20 shrink-0 font-mono text-[10px] text-mathua-muted">
            Series
          </label>
          <select
            id="series-kind"
            aria-label="Series"
            value={kind}
            onChange={(event) => setKind(event.target.value as SumKind)}
            className="h-7 min-w-0 flex-1 border border-mathua-border bg-transparent px-1 font-mono text-[10px] text-mathua-secondary"
          >
            {Object.entries(SUM_NAMES).map(([value, label]) => (
              <option key={value} value={value}>
                {label}
              </option>
            ))}
          </select>
        </div>
        <Slider label="Terms" value={count} min={2} max={40} step={1} onChange={setCount} />
      </div>
    </FigureFrame>
  )
}
