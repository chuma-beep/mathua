'use client'

import { useMemo, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Slider } from '../Controls'
import { surfaceSchema } from '../schemas'
import {
  FIELD_LABELS,
  directional,
  fieldCritical,
  fieldGradient,
  fieldHessian,
  fieldValue,
  saddleTest,
  tangentPlane,
  type Field,
} from '../math/field'
import { project } from '../math/project3d'
import { fmt } from '../format'
import type { TemplateProps } from '../types'

const W = 360
const H = 340
const PAD = 30
const EXTENT = 1.8
const N = 16
const Z_SCALE = 0.38
const VIEW = { theta: 0.62, phi: 0.48 }

export default function SurfaceLab({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => surfaceSchema.parse(params), [params])
  const [field, setField] = useState<Field>(defaults.field)
  const [x0, setX0] = useState(defaults.x0)
  const [y0, setY0] = useState(defaults.y0)

  const mesh = useMemo(() => {
    const points: Array<Array<{ x: number; y: number; z: number; depth: number; sx: number; sy: number }>> = []
    for (let i = 0; i <= N; i++) {
      const row = []
      for (let j = 0; j <= N; j++) {
        const x = -EXTENT + (2 * EXTENT * i) / N
        const y = -EXTENT + (2 * EXTENT * j) / N
        const z = fieldValue(field, x, y)
        const p = project({ x, y, z: z * Z_SCALE }, VIEW)
        row.push({ x, y, z, depth: p.depth, sx: p.x, sy: p.y })
      }
      points.push(row)
    }
    const all = points.flat()
    const minX = Math.min(...all.map((p) => p.sx))
    const maxX = Math.max(...all.map((p) => p.sx))
    const minY = Math.min(...all.map((p) => p.sy))
    const maxY = Math.max(...all.map((p) => p.sy))
    const minDepth = Math.min(...all.map((p) => p.depth))
    const maxDepth = Math.max(...all.map((p) => p.depth))
    const unit = Math.min(
      (W - 2 * PAD) / (maxX - minX || 1),
      (H - 2 * PAD) / (maxY - minY || 1),
    )
    const padX = (W - 2 * PAD - (maxX - minX) * unit) / 2
    const padY = (H - 2 * PAD - (maxY - minY) * unit) / 2
    const SX = (value: number) => PAD + padX + (value - minX) * unit
    const SY = (value: number) => H - PAD - padY - (value - minY) * unit

    // A wireframe rather than filled quads: a height field seen in axonometric
    // projection has cells that fold over near a saddle, and a filled polygon
    // that folds draws a shape the surface does not have. Lines that cross are
    // read as a grid, which is the honest reading.
    const lines: Array<{ points: string; depth: number; weight: number }> = []
    for (let i = 0; i <= N; i++) {
      const row = points[i]
      lines.push({
        points: row.map((p) => `${SX(p.sx)},${SY(p.sy)}`).join(' '),
        depth: row.reduce((sum, p) => sum + p.depth, 0) / row.length,
        weight: i % 3 === 0 ? 1.1 : 0.6,
      })
    }
    for (let j = 0; j <= N; j++) {
      const column = points.map((row) => row[j])
      lines.push({
        points: column.map((p) => `${SX(p.sx)},${SY(p.sy)}`).join(' '),
        depth: column.reduce((sum, p) => sum + p.depth, 0) / column.length,
        weight: j % 3 === 0 ? 1.1 : 0.6,
      })
    }
    lines.sort((a, b) => a.depth - b.depth)
    return { SX, SY, lines, minDepth, maxDepth }
  }, [field])

  const plane = tangentPlane(field, x0, y0)
  const test = saddleTest(field, x0, y0)
  const hessian = fieldHessian(field, x0, y0)
  const [fx, fy] = fieldGradient(field, x0, y0)
  const probe = project({ x: x0, y: y0, z: plane.z0 }, VIEW)
  const gradientLength = Math.hypot(fx, fy)
  const dir = gradientLength > 1e-9 ? [fx / gradientLength, fy / gradientLength] : [0, 0]
  const tip = project({ x: x0 + dir[0] * 0.7, y: y0 + dir[1] * 0.7, z: plane.z0 }, VIEW)
  const patch = 0.75
  const patchCorners: Array<[number, number]> = [
    [x0 - patch, y0 - patch],
    [x0 + patch, y0 - patch],
    [x0 + patch, y0 + patch],
    [x0 - patch, y0 + patch],
  ]
  const patchPoints = patchCorners
    .map(([x, y]) => project({ x, y, z: plane.at(x, y) }, VIEW))
    .map((p) => `${mesh.SX(p.x)},${mesh.SY(p.y)}`)
    .join(' ')

  const reset = () => {
    setField(defaults.field)
    setX0(defaults.x0)
    setY0(defaults.y0)
  }

  return (
    <FigureFrame
      title={title ?? 'The tangent plane matches the value and both slopes'}
      hint="Move the probe. The plane touches the surface at the probe and tilts with the two partial derivatives. At a critical point the tilt vanishes and the Hessian decides the shape."
      caption={
        caption ??
        'The surface is z = f(x, y). The shaded patch is the tangent plane, whose coefficients are the two partial derivatives at the probe. The arrow is the gradient, which points uphill in the base plane. The discriminant D = f_xx f_yy − f_xy² classifies a critical point.'
      }
      onReset={reset}
      readout={[
        { label: 'surface', value: FIELD_LABELS[field] },
        { label: 'probe', value: `(${fmt(x0, 2)}, ${fmt(y0, 2)})` },
        { label: 'z at probe', value: fmt(plane.z0, 3) },
        { label: 'gradient', value: `(${fmt(fx, 2)}, ${fmt(fy, 2)})` },
        { label: 'tangent plane', value: plane.equation },
        { label: 'D = f_xx f_yy − f_xy²', value: fmt(test.discriminant, 3) },
        { label: 'f_xx', value: fmt(hessian.xx, 2) },
        { label: 'classification', value: test.classification },
        { label: 'directional at 45°', value: fmt(directional(field, x0, y0, Math.PI / 4), 3) },
        { label: 'stationary point', value: fieldCritical(field).map(([x, y]) => `(${fmt(x)}, ${fmt(y)})`).join(', ') },
      ]}
    >
      <svg
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        role="img"
        aria-label={`Surface ${FIELD_LABELS[field]} with the tangent plane at (${fmt(x0, 2)}, ${fmt(y0, 2)})`}
      >
        {mesh.lines.map((line, index) => {
          const t = mesh.maxDepth > mesh.minDepth ? (line.depth - mesh.minDepth) / (mesh.maxDepth - mesh.minDepth) : 1
          return (
            <polyline
              key={index}
              points={line.points}
              fill="none"
              stroke="var(--accent-blue)"
              strokeWidth={line.weight}
              opacity={0.18 + 0.62 * t}
            />
          )
        })}

        <polygon points={patchPoints} fill="var(--accent-teal)" fillOpacity="0.28" stroke="var(--accent-teal)" strokeWidth="1.4" />
        <line
          x1={mesh.SX(probe.x)}
          y1={mesh.SY(probe.y)}
          x2={mesh.SX(tip.x)}
          y2={mesh.SY(tip.y)}
          stroke="var(--accent-red)"
          strokeWidth="1.8"
          markerEnd="url(#surface-arrow)"
        />
        <defs>
          <marker id="surface-arrow" viewBox="0 0 8 8" refX="7" refY="4" markerWidth="6" markerHeight="6" orient="auto">
            <path d="M0,0 L8,4 L0,8 z" fill="var(--accent-red)" />
          </marker>
        </defs>
        <circle cx={mesh.SX(probe.x)} cy={mesh.SY(probe.y)} r="4" fill="var(--text-primary)" />
        <text x={PAD} y={16} fontSize="8" fill="var(--text-muted)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
          gradient ({fmt(fx, 2)}, {fmt(fy, 2)})
        </text>
      </svg>

      <div className="mt-3 grid grid-cols-1 gap-2">
        <div className="flex items-center gap-2">
          <label htmlFor="surface-field" className="w-20 shrink-0 font-mono text-[10px] text-mathua-muted">
            Surface
          </label>
          <select
            id="surface-field"
            aria-label="Surface"
            value={field}
            onChange={(event) => setField(event.target.value as Field)}
            className="h-7 min-w-0 flex-1 border border-mathua-border bg-transparent px-1 font-mono text-[10px] text-mathua-secondary"
          >
            {Object.entries(FIELD_LABELS).map(([value, label]) => (
              <option key={value} value={value}>
                z = {label}
              </option>
            ))}
          </select>
        </div>
        <Slider label="Probe x" value={x0} min={-2} max={2} step={0.1} onChange={setX0} />
        <Slider label="Probe y" value={y0} min={-2} max={2} step={0.1} onChange={setY0} />
      </div>
    </FigureFrame>
  )
}
