'use client'

import { useMemo, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Slider } from '../Controls'
import { solidRevSchema } from '../schemas'
import {
  SOLID_FORMULAS,
  SOLID_LABELS,
  SOLID_X0,
  SOLID_X1,
  innerRadius,
  outerRadius,
  shellHeight,
  sliceArea,
  sliceVolume,
  solidVolume,
  type SolidMode,
} from '../math/solidrev'
import { SHADES, project, type Vec3 } from '../math/project3d'
import { fmt } from '../format'
import type { TemplateProps } from '../types'

const W = 360
const H = 340
const PAD = 26
const VIEW = { theta: 0.7, phi: 0.62 }
const CENTER = (SOLID_X0 + SOLID_X1) / 2
const DX = 0.05

interface Ring {
  points: Array<[number, number]>
  depth: number
  shade: 0 | 1 | 2 | 3
}

export default function SolidRevolution({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => solidRevSchema.parse(params), [params])
  const [mode, setMode] = useState<SolidMode>(defaults.mode)
  const [slice, setSlice] = useState(defaults.slice)

  const geometry = useMemo(() => {
    const rings: Vec3[][] = []
    if (mode === 'shell') {
      for (let i = 0; i <= 14; i++) {
        const x = SOLID_X0 + ((SOLID_X1 - SOLID_X0) * i) / 14
        const ring: Vec3[] = []
        for (let j = 0; j <= 40; j++) {
          const angle = (2 * Math.PI * j) / 40
          ring.push({ x: x * Math.cos(angle), y: x * Math.sin(angle), z: shellHeight(x) })
        }
        rings.push(ring)
      }
    } else {
      for (let i = 0; i <= 18; i++) {
        const x = SOLID_X0 + ((SOLID_X1 - SOLID_X0) * i) / 18
        const r = outerRadius(x)
        const ring: Vec3[] = []
        for (let j = 0; j <= 40; j++) {
          const angle = (2 * Math.PI * j) / 40
          ring.push({ x: x - CENTER, y: r * Math.cos(angle), z: r * Math.sin(angle) })
        }
        rings.push(ring)
      }
    }

    const projected = rings.map((ring) => ring.map((point) => project(point, VIEW)))
    const all = projected.flat()
    const minX = Math.min(...all.map((p) => p.x))
    const maxX = Math.max(...all.map((p) => p.x))
    const minY = Math.min(...all.map((p) => p.y))
    const maxY = Math.max(...all.map((p) => p.y))
    const unit = Math.min(
      (W - 2 * PAD) / (maxX - minX || 1),
      (H - 2 * PAD) / (maxY - minY || 1),
    )
    const padX = (W - 2 * PAD - (maxX - minX) * unit) / 2
    const padY = (H - 2 * PAD - (maxY - minY) * unit) / 2
    const SX = (value: number) => PAD + padX + (value - minX) * unit
    const SY = (value: number) => H - PAD - padY - (value - minY) * unit

    const quads: Ring[] = []
    for (let i = 0; i < projected.length - 1; i++) {
      for (let j = 0; j < projected[i].length - 1; j++) {
        const corners = [projected[i][j], projected[i + 1][j], projected[i + 1][j + 1], projected[i][j + 1]]
        quads.push({
          points: corners.map((p) => [SX(p.x), SY(p.y)] as [number, number]),
          depth: corners.reduce((sum, p) => sum + p.depth, 0) / 4,
          shade: (j % 8 < 4 ? 1 : 2) as 0 | 1 | 2 | 3,
        })
      }
    }
    quads.sort((a, b) => a.depth - b.depth)

    const ringPath = (points: Vec3[]) =>
      points.map((point) => {
        const p = project(point, VIEW)
        return `${SX(p.x)},${SY(p.y)}`
      }).join(' ')

    let slicePath = ''
    let holePath = ''
    if (mode === 'shell') {
      const top: Vec3[] = []
      const bottom: Vec3[] = []
      for (let j = 0; j <= 48; j++) {
        const angle = (2 * Math.PI * j) / 48
        top.push({ x: slice * Math.cos(angle), y: slice * Math.sin(angle), z: shellHeight(slice) })
        bottom.push({ x: slice * Math.cos(angle), y: slice * Math.sin(angle), z: 0 })
      }
      slicePath = ringPath(top)
      holePath = ringPath(bottom)
    } else {
      const outer: Vec3[] = []
      const inner: Vec3[] = []
      for (let j = 0; j <= 48; j++) {
        const angle = (2 * Math.PI * j) / 48
        outer.push({ x: slice - CENTER, y: outerRadius(slice) * Math.cos(angle), z: outerRadius(slice) * Math.sin(angle) })
        inner.push({ x: slice - CENTER, y: innerRadius(slice) * Math.cos(angle), z: innerRadius(slice) * Math.sin(angle) })
      }
      slicePath = ringPath(outer)
      holePath = ringPath(inner)
    }

    return { quads, slicePath, holePath, SX, SY }
  }, [mode, slice])

  const outer = mode === 'shell' ? slice : outerRadius(slice)
  const volume = solidVolume(mode)
  const area = sliceArea(mode, slice)

  const reset = () => {
    setMode(defaults.mode)
    setSlice(defaults.slice)
  }

  return (
    <FigureFrame
      title={title ?? 'One body, three ways to slice it'}
      hint="Move the slice and change the method. The volume is the sum of the slices, so the formula and the picture are the same statement."
      caption={
        caption ??
        'The body is swept from the drawn profile. A disk slice is a thin cylinder, a washer slice is a thin annulus and a shell slice is a thin tube. Each readout is the integral of that slice over the interval. Disks and washers share one body about the horizontal axis; the shell body is a second solid about the vertical axis, so its volume is not meant to match theirs.'
      }
      onReset={reset}
      controls={
        <>
          <div className="flex items-center gap-2">
            <label htmlFor="solid-method" className="w-20 shrink-0 font-mono text-xs text-mathua-muted">
              Method
            </label>
            <select
              id="solid-method"
              aria-label="Method"
              value={mode}
              onChange={(event) => setMode(event.target.value as SolidMode)}
              className="h-10 min-w-0 flex-1 border border-mathua-border bg-transparent px-2 font-mono text-xs text-mathua-secondary"
            >
              {Object.entries(SOLID_LABELS).map(([value, label]) => (
                <option key={value} value={value}>
                  {label}
                </option>
              ))}
            </select>
          </div>
          <Slider label="Slice position" value={slice} min={0.1} max={2} step={0.05} onChange={setSlice} />
        </>
      }
      readout={[
        { label: 'method', value: SOLID_LABELS[mode] },
        { label: 'formula', value: SOLID_FORMULAS[mode] },
        { label: 'slice at x', value: fmt(slice, 2) },
        { label: 'radius there', value: fmt(outer, 3) },
        { label: 'slice area', value: fmt(area, 4) },
        { label: 'slice volume', value: fmt(sliceVolume(mode, slice, DX), 4) },
        { label: 'thickness used', value: fmt(DX, 2) },
        { label: 'total volume', value: fmt(volume, 4) },
        { label: 'interval', value: `${fmt(SOLID_X0)} ≤ x ≤ ${fmt(SOLID_X1)}` },
      ]}
    >
      <svg
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        role="img"
        aria-label={`${SOLID_LABELS[mode]} body with the slice at x = ${fmt(slice, 2)} highlighted`}
      >
        {geometry.quads.map((quad, index) => (
          <polygon
            key={index}
            points={quad.points.map(([x, y]) => `${x},${y}`).join(' ')}
            fill={SHADES[quad.shade]}
            fillOpacity="0.35"
            stroke="var(--bg)"
            strokeWidth="0.3"
          />
        ))}

        {geometry.slicePath ? (
          <polygon points={geometry.slicePath} fill="var(--accent-red)" fillOpacity="0.35" stroke="var(--accent-red)" strokeWidth="1.4" />
        ) : null}
        {geometry.holePath ? (
          <polygon points={geometry.holePath} fill="var(--bg)" fillOpacity="0.85" stroke="var(--accent-red)" strokeWidth="1" />
        ) : null}

        <text x={PAD} y={16} fontSize="8" fill="var(--text-muted)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
          {SOLID_FORMULAS[mode]}
        </text>
        <text x={PAD} y={30} fontSize="8" fill="var(--accent-red)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
          slice at x = {fmt(slice, 2)}
        </text>
      </svg>
    </FigureFrame>
  )
}
