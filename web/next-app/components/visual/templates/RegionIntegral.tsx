'use client'

import { useMemo, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Slider } from '../Controls'
import { XAxisTicks, YAxisTicks } from '../Axis'
import { regionSchema } from '../schemas'
import {
  DENSITY_LABELS,
  REGION_EQUATIONS,
  REGION_LABELS,
  areaElement,
  cellStats,
  densityValue,
  exactIntegral,
  regionBounds,
  regionCells,
  regionSum,
  reversedBounds,
  type Density,
  type RegionKind,
} from '../math/region'

import { fmt } from '../format'
import type { TemplateProps } from '../types'

const W = 360
const H = 330
const PAD = 34

export default function RegionIntegral({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => regionSchema.parse(params), [params])
  const [kind, setKind] = useState<RegionKind>(defaults.region)
  const [density, setDensity] = useState<Density>(defaults.density)
  const [n, setN] = useState(defaults.n)

  const bounds = regionBounds(kind)
  const cells = useMemo(() => regionCells(kind, n), [kind, n])
  const estimate = useMemo(() => regionSum(kind, density, n), [kind, density, n])
  const exact = exactIntegral(kind, density)
  const stats = cellStats(kind, n)
  const error = Math.abs(estimate - exact)

  // Equal cells must be drawn as squares, so both axes share one scale and the
  // shorter dimension is padded rather than stretched.
  const spanX = bounds.x1 - bounds.x0 + 0.2
  const spanY = bounds.y1 - bounds.y0 + 0.2
  const unit = Math.min((W - 2 * PAD) / spanX, (H - 2 * PAD) / spanY)
  const originX = PAD + ((W - 2 * PAD) - spanX * unit) / 2
  const originY = PAD + ((H - 2 * PAD) - spanY * unit) / 2
  const X = (value: number) => originX + (value - (bounds.x0 - 0.1)) * unit
  const Y = (value: number) => originY + (bounds.y1 + 0.1 - value) * unit
  const dx = (bounds.x1 - bounds.x0) / n
  const dy = (bounds.y1 - bounds.y0) / n
  const ticksX = kind === 'rect' ? [0, 1, 2] : [0, 0.5, 1]
  const ticksY = [0, 0.5, 1]

  const boundary = useMemo(() => {
    if (kind !== 'between') return null
    const lower: Array<[number, number]> = []
    const upper: Array<[number, number]> = []
    for (let index = 0; index <= 80; index++) {
      const x = index / 80
      lower.push([x, x * x])
      upper.push([x, x])
    }
    return { lower, upper }
  }, [kind])

  const reset = () => {
    setKind(defaults.region)
    setDensity(defaults.density)
    setN(defaults.n)
  }

  return (
    <FigureFrame
      title={title ?? 'The integral is a sum of small cells'}
      hint="Increase the grid and watch the estimate approach the exact value. Only centres inside the region contribute, so a coarse grid misses the boundary cells."
      caption={
        caption ??
        'The region is partitioned into equal cells. A cell counts when its centre lies inside, and its contribution is the density there times the cell area. The exact value comes from the bounds below, not from the grid, so the error column is a real measurement.'
      }
      onReset={reset}
      readout={[
        { label: 'region', value: REGION_LABELS[kind] },
        { label: 'bounds', value: REGION_EQUATIONS[kind] },
        { label: 'density', value: DENSITY_LABELS[density] },
        { label: 'cells inside', value: String(stats.count) },
        { label: 'cell area', value: fmt(stats.area, 5) },
        { label: 'estimate', value: fmt(estimate, 4) },
        { label: 'exact', value: fmt(exact, 4) },
        { label: 'error', value: fmt(error, 4) },
        { label: 'reversed', value: reversedBounds(kind) },
        { label: 'area element', value: areaElement(kind) },
      ]}
    >
      <svg
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        role="img"
        aria-label={`${REGION_LABELS[kind]} partitioned into ${n} by ${n} cells, shaded where the centre is inside`}
      >
        <line x1={PAD} x2={W - PAD} y1={Y(bounds.y0)} y2={Y(bounds.y0)} stroke="var(--text-muted)" strokeWidth="0.8" />
        <line x1={X(bounds.x0)} x2={X(bounds.x0)} y1={PAD} y2={H - PAD} stroke="var(--text-muted)" strokeWidth="0.8" />
        <XAxisTicks ticks={ticksX} scale={X} at={Y(bounds.y0)} />
        <YAxisTicks ticks={ticksY} scale={Y} at={X(bounds.x0)} />

        {cells.map((cell, index) => (
          <rect
            key={index}
            x={X(cell.x - dx / 2)}
            y={Y(cell.y + dy / 2)}
            width={Math.max(0, X(cell.x + dx / 2) - X(cell.x - dx / 2))}
            height={Math.max(0, Y(cell.y - dy / 2) - Y(cell.y + dy / 2))}
            fill={cell.inside ? 'var(--accent-blue)' : 'transparent'}
            fillOpacity={cell.inside ? 0.12 + 0.5 * Math.min(1, densityValue(density, cell.x, cell.y)) : 0}
            stroke="var(--text-muted)"
            strokeWidth="0.3"
            strokeOpacity={cell.inside ? 0.5 : 0.18}
          />
        ))}

        {boundary ? (
          <g fill="none" stroke="var(--accent-red)" strokeWidth="1.6">
            <polyline points={boundary.lower.map(([x, y]) => `${X(x)},${Y(y)}`).join(' ')} />
            <polyline points={boundary.upper.map(([x, y]) => `${X(x)},${Y(y)}`).join(' ')} />
          </g>
        ) : null}

        {kind === 'polar' ? (
          <path
            d={`M ${X(1)} ${Y(0)} A ${X(1) - X(0)} ${Y(0) - Y(1)} 0 0 0 ${X(0)} ${Y(1)}`}
            fill="none"
            stroke="var(--accent-red)"
            strokeWidth="1.6"
          />
        ) : null}

        <text x={PAD} y={PAD - 8} fontSize="9" fill="var(--text-muted)">
          {n} × {n} cells · {stats.count} inside
        </text>
      </svg>

      <div className="mt-3 grid grid-cols-1 gap-2">
        <div className="flex items-center gap-2">
          <label htmlFor="region-kind" className="w-20 shrink-0 font-mono text-[10px] text-mathua-muted">
            Region
          </label>
          <select
            id="region-kind"
            aria-label="Region"
            value={kind}
            onChange={(event) => setKind(event.target.value as RegionKind)}
            className="h-7 min-w-0 flex-1 border border-mathua-border bg-transparent px-1 font-mono text-[10px] text-mathua-secondary"
          >
            {Object.entries(REGION_LABELS).map(([value, label]) => (
              <option key={value} value={value}>
                {label}
              </option>
            ))}
          </select>
        </div>
        <div className="flex items-center gap-2">
          <label htmlFor="region-density" className="w-20 shrink-0 font-mono text-[10px] text-mathua-muted">
            Density
          </label>
          <select
            id="region-density"
            aria-label="Density"
            value={density}
            onChange={(event) => setDensity(event.target.value as Density)}
            className="h-7 min-w-0 flex-1 border border-mathua-border bg-transparent px-1 font-mono text-[10px] text-mathua-secondary"
          >
            {Object.entries(DENSITY_LABELS).map(([value, label]) => (
              <option key={value} value={value}>
                {label}
              </option>
            ))}
          </select>
        </div>
        <Slider label="Cells per side" value={n} min={4} max={40} step={2} onChange={setN} />
      </div>
    </FigureFrame>
  )
}
