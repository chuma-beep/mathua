'use client'

import { useId, useMemo, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Slider, Toggle } from '../Controls'
import { XAxisTicks, YAxisTicks } from '../Axis'
import { vectorFieldSchema } from '../schemas'
import {
  VF_LABELS,
  VF_NAMES,
  arrowLength,
  circulation,
  flux,
  loopArea,
  loopPath,
  streamlines,
  vfCurl,
  vfDivergence,
  vfValue,
  type VectorFieldPreset,
} from '../math/vectorfield'
import { scaleLinear } from '../scale'
import { fmt } from '../format'
import type { TemplateProps } from '../types'

const W = 360
const H = 360
const PAD = 34
const EXTENT = 2.4
const GRID = 9

export default function VectorFieldLab({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => vectorFieldSchema.parse(params), [params])
  const [preset, setPreset] = useState<VectorFieldPreset>(defaults.preset)
  const [radius, setRadius] = useState(defaults.radius)
  const [showLines, setShowLines] = useState(defaults.showStreamlines)
  const arrowId = useId()

  const X = scaleLinear(-EXTENT, EXTENT, PAD, W - PAD)
  const Y = scaleLinear(-EXTENT, EXTENT, H - PAD, PAD)
  const ticks = [-2, -1, 0, 1, 2]
  const path = useMemo(() => loopPath(radius, 80), [radius])
  const lines = useMemo(() => (showLines ? streamlines(preset) : []), [preset, showLines])

  const arrows = useMemo(() => {
    const out: Array<{ x1: number; y1: number; x2: number; y2: number }> = []
    for (let row = 0; row < GRID; row++) {
      for (let column = 0; column < GRID; column++) {
        const x = -EXTENT + 0.4 + (column * (2 * EXTENT - 0.8)) / (GRID - 1)
        const y = -EXTENT + 0.4 + (row * (2 * EXTENT - 0.8)) / (GRID - 1)
        const [ux, uy] = vfValue(preset, x, y)
        const length = Math.hypot(ux, uy)
        if (length < 1e-9) continue
        const reach = 0.1 + 0.2 * arrowLength(preset, x, y)
        out.push({
          x1: X(x),
          y1: Y(y),
          x2: X(x + (ux / length) * reach),
          y2: Y(y + (uy / length) * reach),
        })
      }
    }
    return out
  }, [preset, X, Y])

  const div = vfDivergence(preset)
  const curl = vfCurl(preset)
  const circ = circulation(preset, radius)
  const out = flux(preset, radius)
  const area = loopArea(radius)
  const green = Math.abs(circ - curl * area) < 1e-9

  const reset = () => {
    setPreset(defaults.preset)
    setRadius(defaults.radius)
    setShowLines(defaults.showStreamlines)
  }

  return (
    <FigureFrame
      title={title ?? 'Divergence and curl are two readings of one field'}
      hint="Change the field and the loop. Arrows are the field at each point, the circle is the loop the integrals run around, and the streamlines follow the same arrows."
      caption={
        caption ??
        'Each arrow is the vector at its base. The circulation is the line integral around the red loop and the flux is the flow through it. Green\'s theorem says the circulation equals the enclosed area times the curl, and the divergence theorem says the flux equals the area times the divergence.'
      }
      onReset={reset}
      controls={
        <div className="grid grid-cols-1 gap-2">
          <div className="flex items-center gap-2">
            <label htmlFor="vector-field" className="w-20 shrink-0 font-mono text-xs text-mathua-muted">
              Field
            </label>
            <select
              id="vector-field"
              aria-label="Field"
              value={preset}
              onChange={(event) => setPreset(event.target.value as VectorFieldPreset)}
              className="h-10 min-w-0 flex-1 border border-mathua-border bg-transparent px-2 font-mono text-xs text-mathua-secondary"
            >
              {Object.entries(VF_NAMES).map(([value, label]) => (
                <option key={value} value={value}>
                  {label} — {VF_LABELS[value as VectorFieldPreset]}
                </option>
              ))}
            </select>
          </div>
          <Slider label="Loop radius" value={radius} min={0.4} max={2} step={0.1} onChange={setRadius} />
          <Toggle label="Show streamlines" checked={showLines} onChange={setShowLines} />
        </div>
      }
      readout={[
        { label: 'field', value: `${VF_NAMES[preset]} · ${VF_LABELS[preset]}` },
        { label: 'divergence', value: fmt(div, 2) },
        { label: 'curl', value: fmt(curl, 2) },
        { label: 'loop radius', value: fmt(radius, 2) },
        { label: 'enclosed area', value: fmt(area, 3) },
        { label: 'circulation', value: fmt(circ, 3) },
        { label: 'curl × area', value: fmt(curl * area, 3) },
        { label: 'flux', value: fmt(out, 3) },
        { label: 'divergence × area', value: fmt(div * area, 3) },
        { label: 'Green check', value: green ? 'holds' : 'fails' },
      ]}
    >
      <svg
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        role="img"
        aria-label={`Vector field ${VF_LABELS[preset]} with a loop of radius ${fmt(radius, 2)}`}
      >
        <defs>
          <marker
            id={arrowId}
            viewBox="0 0 8 8"
            refX="7"
            refY="4"
            markerWidth="4.5"
            markerHeight="4.5"
            markerUnits="userSpaceOnUse"
            orient="auto"
          >
            <path d="M0,0 L8,4 L0,8 z" fill="var(--text-muted)" />
          </marker>
        </defs>
        <line x1={PAD} x2={W - PAD} y1={Y(0)} y2={Y(0)} stroke="var(--text-muted)" strokeWidth="0.8" />
        <line x1={X(0)} x2={X(0)} y1={PAD} y2={H - PAD} stroke="var(--text-muted)" strokeWidth="0.8" />
        <XAxisTicks ticks={ticks} scale={X} at={Y(0)} />
        <YAxisTicks ticks={ticks} scale={Y} at={X(0)} />

        {arrows.map((arrow, index) => (
          <line
            key={index}
            x1={arrow.x1}
            y1={arrow.y1}
            x2={arrow.x2}
            y2={arrow.y2}
            stroke="var(--text-muted)"
            strokeWidth="0.9"
            opacity="0.85"
            markerEnd={`url(#${arrowId})`}
          />
        ))}

        {lines.map((line, index) => (
          <polyline
            key={index}
            points={line.map(([x, y]) => `${X(x)},${Y(y)}`).join(' ')}
            fill="none"
            stroke="var(--accent-blue)"
            strokeWidth="1"
            opacity="0.55"
          />
        ))}

        <polyline
          points={path.map(([x, y]) => `${X(x)},${Y(y)}`).join(' ')}
          fill="none"
          stroke="var(--accent-red)"
          strokeWidth="1.8"
        />
        <circle cx={X(radius)} cy={Y(0)} r="3" fill="var(--accent-red)" />
        <text x={X(radius) + 5} y={Y(0) - 5} fontSize="8" fill="var(--accent-red)">
          r = {fmt(radius, 2)}
        </text>

        <text x={PAD} y={PAD - 10} fontSize="9" fill="var(--text-muted)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
          div {fmt(div, 2)} · curl {fmt(curl, 2)}
        </text>
      </svg>
    </FigureFrame>
  )
}
