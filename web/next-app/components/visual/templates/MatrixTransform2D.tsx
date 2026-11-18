'use client'

import { useCallback, useMemo, useState, useId } from 'react'
import FigureFrame from '../FigureFrame'
import { NumberField, Slider, Toggle } from '../Controls'
import { matrix2dSchema } from '../schemas'
import { apply, det, fromFlat } from '../math/matrix'
import { fmt, fmtVec } from '../format'
import type { TemplateProps } from '../types'

const SIZE = 320
const EXTENT = 3.6

export default function MatrixTransform2D({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => {
    const p = matrix2dSchema.parse(params)
    return {
      matrix: [Number(p.matrix[0]), Number(p.matrix[1]), Number(p.matrix[2]), Number(p.matrix[3])] as [number, number, number, number],
      vector: [Number(p.vector[0]), Number(p.vector[1])] as [number, number],
      showGrid: p.showGrid,
    }
  }, [params])
  const [matrix, setMatrix] = useState<[number, number, number, number]>(defaults.matrix)
  const [vector, setVector] = useState<[number, number]>(defaults.vector)
  const [showGrid, setShowGrid] = useState(defaults.showGrid)
  const uid = useId().replace(/:/g, '')

  const m = fromFlat(matrix)
  const image = apply(m, vector)
  const area = det(m)

  const pad = 24
  const unit = (SIZE - pad * 2) / (EXTENT * 2)
  const cx = SIZE / 2
  const cy = SIZE / 2
  const X = useCallback((x: number) => cx + x * unit, [cx, unit])
  const Y = useCallback((y: number) => cy - y * unit, [cy, unit])

  const setEntry = (index: number, value: number) => {
    setMatrix((prev) => {
      const next = [...prev] as [number, number, number, number]
      next[index] = value
      return next
    })
  }

  const reset = () => {
    setMatrix(defaults.matrix)
    setVector(defaults.vector)
    setShowGrid(defaults.showGrid)
  }

  const gridLines = useMemo(() => {
    if (!showGrid) return []
    const lines: Array<{ from: [number, number]; to: [number, number] }> = []
    for (let k = -4; k <= 4; k++) {
      if (k === 0) continue
      lines.push({ from: [k, -EXTENT], to: [k, EXTENT] })
      lines.push({ from: [-EXTENT, k], to: [EXTENT, k] })
    }
    return lines
  }, [showGrid])

  const unitSquare: Array<[number, number]> = [
    [0, 0],
    [1, 0],
    [1, 1],
    [0, 1],
  ]

  return (
    <FigureFrame
      title={title ?? 'Matrix as a transformation'}
      caption={
        caption ??
        'The grid is the plane itself. Each basis arrow is where the matrix sends e₁ and e₂, and the shaded square is the image of the unit square. Its area is |det A|.'
      }
      onReset={reset}
      readout={[
        { label: 'A', value: `[[${fmt(matrix[0], 1)}, ${fmt(matrix[1], 1)}], [${fmt(matrix[2], 1)}, ${fmt(matrix[3], 1)}]]` },
        { label: 'det A', value: fmt(area) },
        { label: 'v', value: fmtVec(vector) },
        { label: 'Av', value: fmtVec(image) },
        { label: '|Av|', value: fmt(Math.hypot(image[0], image[1])) },
      ]}
    >
      <svg viewBox={`0 0 ${SIZE} ${SIZE}`} className="h-auto w-full font-mono" role="img" aria-label="A matrix acting on the plane">
        <defs>
          <marker id={`${uid}-blue`} markerWidth="6" markerHeight="6" refX="5" refY="3" orient="auto">
            <path d="M0,0 L6,3 L0,6 Z" fill="var(--accent-blue)" />
          </marker>
          <marker id={`${uid}-teal`} markerWidth="6" markerHeight="6" refX="5" refY="3" orient="auto">
            <path d="M0,0 L6,3 L0,6 Z" fill="var(--accent-teal)" />
          </marker>
          <marker id={`${uid}-green`} markerWidth="6" markerHeight="6" refX="5" refY="3" orient="auto">
            <path d="M0,0 L6,3 L0,6 Z" fill="var(--accent-green)" />
          </marker>
        </defs>

        {/* The untransformed grid, kept faint so the deformation is visible against it. */}
        {gridLines.map((line, i) => (
          <line
            key={`g${i}`}
            x1={X(line.from[0])}
            y1={Y(line.from[1])}
            x2={X(line.to[0])}
            y2={Y(line.to[1])}
            stroke="var(--border-strong)"
            strokeWidth="0.6"
            strokeDasharray="2 4"
            opacity="0.7"
          />
        ))}
        <line x1={pad} y1={cy} x2={SIZE - pad} y2={cy} stroke="var(--text-muted)" strokeWidth="0.8" />
        <line x1={cx} y1={pad} x2={cx} y2={SIZE - pad} stroke="var(--text-muted)" strokeWidth="0.8" />

        {/* The image of the grid: each line is mapped through A, so the deformation is
            the thing being shown rather than an assertion about it. */}
        {gridLines.map((line, i) => {
          const a = apply(m, line.from)
          const b = apply(m, line.to)
          return (
            <line
              key={`t${i}`}
              x1={X(a[0])}
              y1={Y(a[1])}
              x2={X(b[0])}
              y2={Y(b[1])}
              stroke="var(--accent-blue)"
              strokeWidth="0.7"
              opacity="0.35"
            />
          )
        })}

        <path
          d={`M ${X(0)} ${Y(0)} L ${X(unitSquare[1][0])} ${Y(unitSquare[1][1])} L ${X(unitSquare[2][0])} ${Y(unitSquare[2][1])} L ${X(unitSquare[3][0])} ${Y(unitSquare[3][1])} Z`}
          fill="none"
          stroke="var(--accent-green)"
          strokeWidth="1"
          strokeDasharray="4 3"
          opacity="0.7"
        />
        <path
          d={`M ${X(0)} ${Y(0)} L ${X(apply(m, unitSquare[1])[0])} ${Y(apply(m, unitSquare[1])[1])} L ${X(apply(m, unitSquare[2])[0])} ${Y(apply(m, unitSquare[2])[1])} L ${X(apply(m, unitSquare[3])[0])} ${Y(apply(m, unitSquare[3])[1])} Z`}
          fill="var(--accent-blue)"
          fillOpacity="0.12"
          stroke="var(--accent-blue)"
          strokeWidth="1.4"
        />

        <line x1={X(0)} y1={Y(0)} x2={X(m[0][0])} y2={Y(m[1][0])} stroke="var(--accent-blue)" strokeWidth="2" markerEnd={`url(#${uid}-blue)`} />
        <line x1={X(0)} y1={Y(0)} x2={X(m[0][1])} y2={Y(m[1][1])} stroke="var(--accent-teal)" strokeWidth="2" markerEnd={`url(#${uid}-teal)`} />
        <line x1={X(0)} y1={Y(0)} x2={X(vector[0])} y2={Y(vector[1])} stroke="var(--accent-green)" strokeWidth="1.4" markerEnd={`url(#${uid}-green)`} />
        <line x1={X(vector[0])} y1={Y(vector[1])} x2={X(image[0])} y2={Y(image[1])} stroke="var(--accent-green)" strokeWidth="0.8" strokeDasharray="3 3" opacity="0.6" />
        <line x1={X(0)} y1={Y(0)} x2={X(image[0])} y2={Y(image[1])} stroke="var(--accent-blue)" strokeWidth="2.4" markerEnd={`url(#${uid}-blue)`} />

        <text x={X(m[0][0]) + 6} y={Y(m[1][0]) - 6} fontSize="10" fill="var(--accent-blue)">Ae₁</text>
        <text x={X(m[0][1]) + 6} y={Y(m[1][1]) - 6} fontSize="10" fill="var(--accent-teal)">Ae₂</text>
        <text x={X(image[0]) + 7} y={Y(image[1]) + 12} fontSize="10" fill="var(--accent-blue)">Av</text>
      </svg>

      <div className="mt-3 space-y-2">
        <div className="grid grid-cols-2 gap-2">
          <NumberField label="a" value={matrix[0]} onChange={(v) => setEntry(0, v)} />
          <NumberField label="b" value={matrix[1]} onChange={(v) => setEntry(1, v)} />
          <NumberField label="c" value={matrix[2]} onChange={(v) => setEntry(2, v)} />
          <NumberField label="d" value={matrix[3]} onChange={(v) => setEntry(3, v)} />
        </div>
        <Slider label="vx" value={vector[0]} min={-3} max={3} step={0.1} onChange={(v) => setVector((p) => [v, p[1]])} />
        <Slider label="vy" value={vector[1]} min={-3} max={3} step={0.1} onChange={(v) => setVector((p) => [p[0], v])} />
        <Toggle label="Show reference grid" checked={showGrid} onChange={setShowGrid} />
      </div>
    </FigureFrame>
  )
}
