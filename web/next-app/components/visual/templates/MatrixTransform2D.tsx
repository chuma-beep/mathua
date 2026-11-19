'use client'

import { useCallback, useId, useMemo, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { NumberField, Slider, Toggle } from '../Controls'
import { XAxisTicks, YAxisTicks } from '../Axis'
import { matrix2dSchema } from '../schemas'
import { apply, det, fromFlat, mul, type Mat2 } from '../math/matrix'
import { fmt, fmtVec } from '../format'
import type { TemplateProps } from '../types'

const SIZE = 320
const PAD = 26
const EXTENT = 3.6

type Flat2 = [number, number, number, number]

export default function MatrixTransform2D({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => {
    const p = matrix2dSchema.parse(params)
    return {
      matrix: [Number(p.matrix[0]), Number(p.matrix[1]), Number(p.matrix[2]), Number(p.matrix[3])] as Flat2,
      then: p.then
        ? ([Number(p.then[0]), Number(p.then[1]), Number(p.then[2]), Number(p.then[3])] as Flat2)
        : null,
      vector: [Number(p.vector[0]), Number(p.vector[1])] as [number, number],
      showGrid: p.showGrid,
      showBasis: p.showBasis,
    }
  }, [params])

  const [matrix, setMatrix] = useState<Flat2>(defaults.matrix)
  const [then, setThen] = useState<Flat2 | null>(defaults.then)
  const [vector, setVector] = useState<[number, number]>(defaults.vector)
  const [showGrid, setShowGrid] = useState(defaults.showGrid)
  const [showBasis, setShowBasis] = useState(defaults.showBasis)
  const uid = useId().replace(/:/g, '')

  const A = fromFlat(matrix)
  const B: Mat2 | null = then ? fromFlat(then) : null
  const composed = B ? mul(B, A) : null
  const mid = apply(A, vector)
  const image = composed ? apply(composed, vector) : mid
  // The image of the unit square is the area claim, so the determinant is read three
  // ways: its value, the area it makes, and the orientation it keeps or reverses. A
  // zero determinant is not a small area, it is a collapsed one, and the drawing has to
  // show that rather than a thin quadrilateral.
  const determinant = det(A)
  const area = Math.abs(determinant)
  const collapsed = area < 1e-9
  const orientation = collapsed ? 'collapsed' : determinant > 0 ? 'preserved' : 'reversed'
  const imageColor = determinant < 0 ? 'var(--accent-red)' : 'var(--accent-blue)'
  const diagonal = apply(A, [1, 1])

  const unit = (SIZE - PAD * 2) / (EXTENT * 2)
  const cx = SIZE / 2
  const cy = SIZE / 2
  const X = useCallback((x: number) => cx + x * unit, [cx, unit])
  const Y = useCallback((y: number) => cy - y * unit, [cy, unit])

  const setEntry = (which: 'a' | 'b', index: number, value: number) => {
    const target = which === 'a' ? matrix : then
    if (!target) return
    const next = [...target] as Flat2
    next[index] = value
    if (which === 'a') setMatrix(next)
    else setThen(next)
  }

  const reset = () => {
    setMatrix(defaults.matrix)
    setThen(defaults.then)
    setVector(defaults.vector)
    setShowGrid(defaults.showGrid)
    setShowBasis(defaults.showBasis)
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
  const squarePath = (pts: Array<[number, number]>) =>
    pts.map((p, i) => `${i === 0 ? 'M' : 'L'} ${X(p[0])} ${Y(p[1])}`).join(' ') + ' Z'

  return (
    <FigureFrame
      title={title ?? 'Matrix as a transformation'}
      caption={
        caption ??
        'The grid is the plane itself. Each basis arrow is where the matrix sends e₁ and e₂, and the shaded square is the image of the unit square. Its area is |det A|, and the arrows around its edge run the same way as the original when det A is positive and the opposite way when it is negative. At det A = 0 the square collapses to a segment. Volume would need a third dimension, so this figure stops at area.'
      }
      onReset={reset}
      readout={[
        { label: 'A', value: `[[${fmt(matrix[0], 1)}, ${fmt(matrix[1], 1)}], [${fmt(matrix[2], 1)}, ${fmt(matrix[3], 1)}]]` },
        { label: 'det A', value: fmt(determinant) },
        { label: 'area', value: fmt(area) },
        { label: 'orientation', value: orientation },
        { label: 'v', value: fmtVec(vector) },
        { label: 'Av', value: fmtVec(mid) },
        ...(composed
          ? [
              { label: 'BA', value: `[[${fmt(composed[0][0])}, ${fmt(composed[0][1])}], [${fmt(composed[1][0])}, ${fmt(composed[1][1])}]]` },
              { label: 'BAv', value: fmtVec(image) },
              { label: 'det BA', value: fmt(det(composed)) },
            ]
          : []),
        // |Av| is the length of A v, which is `mid` whether or not a second matrix
        // follows. Using `image` here labelled the composed length as |Av| and hid
        // the intermediate one entirely.
        { label: '|Av|', value: fmt(Math.hypot(mid[0], mid[1])) },
        ...(composed ? [{ label: '|BAv|', value: fmt(Math.hypot(image[0], image[1])) }] : []),
      ]}
    >
      <svg
        viewBox={`0 0 ${SIZE} ${SIZE}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        role="img"
        aria-label="A matrix acting on the plane"
      >
        <defs>
          {[
            ['blue', 'var(--accent-blue)'],
            ['teal', 'var(--accent-teal)'],
            ['green', 'var(--accent-green)'],
          ].map(([name, color]) => (
            <marker key={name} id={`${uid}-${name}`} markerWidth="6" markerHeight="6" refX="5" refY="3" orient="auto">
              <path d="M0,0 L6,3 L0,6 Z" fill={color} />
            </marker>
          ))}
          <clipPath id={`${uid}-plot`}>
            <rect x={PAD} y={PAD} width={SIZE - PAD * 2} height={SIZE - PAD * 2} />
          </clipPath>
        </defs>

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
        <line x1={PAD} y1={cy} x2={SIZE - PAD} y2={cy} stroke="var(--text-muted)" strokeWidth="0.8" />
        <line x1={cx} y1={PAD} x2={cx} y2={SIZE - PAD} stroke="var(--text-muted)" strokeWidth="0.8" />

        <g clipPath={`url(#${uid}-plot)`}>
          {gridLines.map((line, i) => {
            const a = apply(A, line.from)
            const b = apply(A, line.to)
            return (
              <line
                key={`t${i}`}
                x1={X(a[0])}
                y1={Y(a[1])}
                x2={X(b[0])}
                y2={Y(b[1])}
                stroke="var(--accent-blue)"
                strokeWidth="0.6"
                opacity="0.22"
              />
            )
          })}
        </g>

        <path d={squarePath(unitSquare)} fill="none" stroke="var(--accent-green)" strokeWidth="1" strokeDasharray="4 3" opacity="0.9" />

        {collapsed ? (
          <>
            {/* A rank-one image is a segment from the origin to A(e1 + e2): every point of
                the square lands on it, so there is no interior left to shade. */}
            <line
              x1={X(0)}
              y1={Y(0)}
              x2={X(diagonal[0])}
              y2={Y(diagonal[1])}
              stroke="var(--accent-red)"
              strokeWidth="3"
            />
            <text x={PAD} y={SIZE - PAD - 4} fontSize="10" fill="var(--accent-red)">
              area 0: the plan is squashed onto a line
            </text>
          </>
        ) : (
          <>
            <path
              d={squarePath(unitSquare.map((p) => apply(A, p)))}
              fill={imageColor}
              fillOpacity="0.12"
              stroke={imageColor}
              strokeWidth="1.4"
            />
            {/* Traversal arrows, so orientation is visible rather than only printed: the
                same vertex order goes round the original and round the image, which means
                the same direction when det A > 0 and the opposite when it is negative. */}
            {unitSquare.map((point, i) => {
              const next = unitSquare[(i + 1) % unitSquare.length]
              const a = apply(A, point)
              const b = apply(A, next)
              const mx = (X(a[0]) + X(b[0])) / 2
              const my = (Y(a[1]) + Y(b[1])) / 2
              const angle = (Math.atan2(Y(b[1]) - Y(a[1]), X(b[0]) - X(a[0])) * 180) / Math.PI
              return (
                <path
                  key={`trav${i}`}
                  d="M -3.5 2.2 L 3.5 0 L -3.5 -2.2 Z"
                  transform={`translate(${mx} ${my}) rotate(${angle})`}
                  fill={imageColor}
                />
              )
            })}
          </>
        )}

        {showBasis && (
          <>
            <line x1={X(0)} y1={Y(0)} x2={X(A[0][0])} y2={Y(A[1][0])} stroke="var(--accent-blue)" strokeWidth="2" markerEnd={`url(#${uid}-blue)`} />
            <line x1={X(0)} y1={Y(0)} x2={X(A[0][1])} y2={Y(A[1][1])} stroke="var(--accent-teal)" strokeWidth="2" markerEnd={`url(#${uid}-teal)`} />
            <text x={X(A[0][0]) + 6} y={Y(A[1][0]) + 14} fontSize="10" fill="var(--accent-blue)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
              Ae₁
            </text>
            <text x={X(A[0][1]) - 26} y={Y(A[1][1]) - 6} fontSize="10" fill="var(--accent-teal)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
              Ae₂
            </text>
          </>
        )}

        <line x1={X(0)} y1={Y(0)} x2={X(vector[0])} y2={Y(vector[1])} stroke="var(--accent-green)" strokeWidth="1.4" markerEnd={`url(#${uid}-green)`} />
        <text x={X(vector[0]) + 6} y={Y(vector[1]) + 14} fontSize="10" fill="var(--accent-green)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
          v
        </text>

        {composed ? (
          <>
            <line x1={X(vector[0])} y1={Y(vector[1])} x2={X(mid[0])} y2={Y(mid[1])} stroke="var(--accent-teal)" strokeWidth="1.4" strokeDasharray="4 3" markerEnd={`url(#${uid}-teal)`} />
            <line x1={X(mid[0])} y1={Y(mid[1])} x2={X(image[0])} y2={Y(image[1])} stroke="var(--accent-blue)" strokeWidth="1.8" markerEnd={`url(#${uid}-blue)`} />
            <text x={X(mid[0]) + 6} y={Y(mid[1]) - 6} fontSize="10" fill="var(--accent-teal)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
              Av
            </text>
            <text x={X(image[0]) + 6} y={Y(image[1]) + 14} fontSize="10" fill="var(--accent-blue)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
              BAv
            </text>
          </>
        ) : (
          <>
            <line x1={X(vector[0])} y1={Y(vector[1])} x2={X(image[0])} y2={Y(image[1])} stroke="var(--accent-green)" strokeWidth="0.8" strokeDasharray="3 3" opacity="0.6" />
            <line x1={X(0)} y1={Y(0)} x2={X(image[0])} y2={Y(image[1])} stroke="var(--accent-blue)" strokeWidth="2.4" markerEnd={`url(#${uid}-blue)`} />
            <text x={X(image[0]) + 7} y={Y(image[1]) + 14} fontSize="10" fill="var(--accent-blue)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
              Av
            </text>
          </>
        )}

        <XAxisTicks ticks={[-3, -2, -1, 1, 2, 3]} scale={X} at={cy} />
        <YAxisTicks ticks={[-3, -2, -1, 1, 2, 3]} scale={Y} at={cx} />
      </svg>

      <div className="mt-3 space-y-2">
        <div className="grid grid-cols-2 gap-2">
          <NumberField label="a" value={matrix[0]} onChange={(v) => setEntry('a', 0, v)} />
          <NumberField label="b" value={matrix[1]} onChange={(v) => setEntry('a', 1, v)} />
          <NumberField label="c" value={matrix[2]} onChange={(v) => setEntry('a', 2, v)} />
          <NumberField label="d" value={matrix[3]} onChange={(v) => setEntry('a', 3, v)} />
        </div>

        {then && (
          <div className="space-y-1">
            <p className="font-mono text-[10px] text-mathua-muted">then B, applied after A</p>
            <div className="grid grid-cols-2 gap-2">
              <NumberField label="a₂" value={then[0]} onChange={(v) => setEntry('b', 0, v)} />
              <NumberField label="b₂" value={then[1]} onChange={(v) => setEntry('b', 1, v)} />
              <NumberField label="c₂" value={then[2]} onChange={(v) => setEntry('b', 2, v)} />
              <NumberField label="d₂" value={then[3]} onChange={(v) => setEntry('b', 3, v)} />
            </div>
          </div>
        )}

        <Slider label="vx" value={vector[0]} min={-3} max={3} step={0.1} onChange={(v) => setVector((p) => [v, p[1]])} />
        <Slider label="vy" value={vector[1]} min={-3} max={3} step={0.1} onChange={(v) => setVector((p) => [p[0], v])} />
        <div className="flex flex-wrap gap-4">
          <Toggle label="Reference grid" checked={showGrid} onChange={setShowGrid} />
          <Toggle label="Basis images" checked={showBasis} onChange={setShowBasis} />
        </div>
      </div>
    </FigureFrame>
  )
}
