'use client'

import { useCallback, useMemo, useRef, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Select, Slider, Toggle } from '../Controls'
import { transform2dSchema } from '../schemas'
import { apply, compose, dilation, identity, reflection, rotation, type Mat2, type Vec2 } from '../math/matrix'
import { maxSideDelta, polygonPath } from '../math/geometry'
import { clamp, pathFromPoints, scaleLinear } from '../scale'
import { fmt, fmtVec } from '../format'
import { useDrag } from '../useDrag'
import type { TemplateProps } from '../types'

const W = 320
const H = 320
const PAD = 24
const DEG = Math.PI / 180

type Kind = 'rotation' | 'reflection' | 'translation' | 'dilation' | 'composition'

const KIND_OPTIONS: Array<{ value: Kind; label: string }> = [
  { value: 'rotation', label: 'Rotation' },
  { value: 'reflection', label: 'Reflection in the y-axis' },
  { value: 'translation', label: 'Translation' },
  { value: 'dilation', label: 'Dilation' },
  { value: 'composition', label: 'Rotation then reflection' },
]

export default function Transform2D({ params, title, caption }: TemplateProps) {
  const initial = useMemo(() => transform2dSchema.parse(params), [params])
  const [kind, setKind] = useState<Kind>(initial.kind)
  const [angle, setAngle] = useState(initial.angle)
  const [scale, setScale] = useState(initial.scale)
  const [dx, setDx] = useState(initial.dx)
  const [dy, setDy] = useState(initial.dy)
  const [showOrbit, setShowOrbit] = useState(initial.showOrbit)
  const svgRef = useRef<SVGSVGElement | null>(null)

  const shape = useMemo<Vec2[]>(
    () => initial.shape.map((pair) => [Number(pair[0]), Number(pair[1])] as Vec2),
    [initial.shape],
  )

  const matrix: Mat2 = useMemo(() => {
    switch (kind) {
      case 'rotation':
        return rotation(angle * DEG)
      case 'reflection':
        return reflection('y')
      case 'translation':
        return identity()
      case 'dilation':
        return dilation(scale)
      case 'composition':
        return compose(reflection('y'), rotation(angle * DEG))
    }
  }, [kind, angle, scale])

  const offset: Vec2 = kind === 'translation' ? [dx, dy] : [0, 0]
  const image = shape.map((p) => {
    const q = apply(matrix, p)
    return [q[0] + offset[0], q[1] + offset[1]] as Vec2
  })

  const extent = Math.max(
    2.4,
    ...shape.flatMap((p) => [Math.abs(p[0]), Math.abs(p[1])]),
    ...image.flatMap((p) => [Math.abs(p[0]), Math.abs(p[1])]),
  ) * 1.15

  const X = scaleLinear(-extent, extent, PAD, W - PAD)
  const Y = scaleLinear(extent, -extent, PAD, H - PAD)
  const cx = X(0)
  const cy = Y(0)

  const onDrag = useCallback(
    (point: { x: number; y: number }) => {
      if (kind === 'rotation' || kind === 'composition') {
        const a = Math.atan2(-(point.y - cy), point.x - cx) / DEG
        setAngle((a + 360) % 360)
      } else if (kind === 'translation') {
        setDx(clamp((point.x - cx) * (2 * extent) / (W - 2 * PAD), -4, 4))
        setDy(clamp((point.y - cy) * (-2 * extent) / (H - 2 * PAD), -4, 4))
      }
    },
    [kind, cy, cx, extent],
  )
  const draggable = kind === 'rotation' || kind === 'composition' || kind === 'translation'
  useDrag(svgRef, onDrag, draggable)

  const reset = () => {
    setKind(initial.kind)
    setAngle(initial.angle)
    setScale(initial.scale)
    setDx(initial.dx)
    setDy(initial.dy)
    setShowOrbit(initial.showOrbit)
  }

  const gridLines = useMemo(() => {
    const n = Math.ceil(extent)
    const out: Array<{ x1: number; y1: number; x2: number; y2: number }> = []
    for (let k = -n; k <= n; k++) {
      out.push({ x1: X(k), y1: Y(-n), x2: X(k), y2: Y(n) })
      out.push({ x1: X(-n), y1: Y(k), x2: X(n), y2: Y(k) })
    }
    return out
  }, [extent, X, Y])

  const delta = maxSideDelta(shape, image)
  const matrixText = `[[${fmt(matrix[0][0])}, ${fmt(matrix[0][1])}], [${fmt(matrix[1][0])}, ${fmt(matrix[1][1])}]]`

  return (
    <FigureFrame
      title={title ?? 'A geometric transformation'}
      caption={
        caption ??
        'The dashed outline is the original, the filled outline is its image. An isometry preserves every side length; a dilation does not. Drag the figure to turn or move it.'
      }
      onReset={reset}
      readout={[
        { label: 'matrix', value: matrixText },
        ...(kind === 'translation' ? [{ label: 'offset', value: fmtVec(offset) }] : []),
        { label: 'image', value: image.slice(0, 3).map((point) => fmtVec(point)).join('  ') },
        { label: 'max side change', value: fmt(delta, 3) },
        { label: 'area scale', value: fmt(Math.abs(matrix[0][0] * matrix[1][1] - matrix[0][1] * matrix[1][0])) },
      ]}
    >
      <svg
        ref={svgRef}
        viewBox={`0 0 ${W} ${H}`}
        className="h-auto w-full font-mono"
        style={{ touchAction: draggable ? 'none' : 'auto' }}
        role="img"
        aria-label="A polygon and its transformed image"
      >
        {gridLines.map((l, i) => (
          <line key={i} x1={l.x1} y1={l.y1} x2={l.x2} y2={l.y2} stroke="var(--border)" strokeWidth="0.5" />
        ))}
        <line x1={PAD} y1={cy} x2={W - PAD} y2={cy} stroke="var(--text-muted)" strokeWidth="0.9" />
        <line x1={cx} y1={PAD} x2={cx} y2={H - PAD} stroke="var(--text-muted)" strokeWidth="0.9" />

        {showOrbit &&
          shape.map((p, i) => (
            <circle
              key={i}
              cx={cx}
              cy={cy}
              r={Math.hypot(X(p[0]) - cx, Y(p[1]) - cy)}
              fill="none"
              stroke="var(--border-strong)"
              strokeWidth="0.8"
              strokeDasharray="3 4"
            />
          ))}

        <path d={polygonPath(shape, X, Y)} fill="none" stroke="var(--accent-green)" strokeWidth="1.4" strokeDasharray="5 4" />
        <path d={polygonPath(image, X, Y)} fill="var(--accent-blue)" fillOpacity="0.12" stroke="var(--accent-blue)" strokeWidth="1.8" />

        {shape.map((p, i) => (
          <circle key={`o${i}`} cx={X(p[0])} cy={Y(p[1])} r="2.6" fill="var(--accent-green)" />
        ))}
        {image.map((p, i) => (
          <circle key={`i${i}`} cx={X(p[0])} cy={Y(p[1])} r="2.6" fill="var(--accent-blue)" />
        ))}

        {draggable && (
          <path
            d={pathFromPoints([
              [cx, cy],
              kind === 'translation' ? [X(offset[0]), Y(offset[1])] : [X(image[0][0]), Y(image[0][1])],
            ])}
            stroke="var(--accent-blue)"
            strokeWidth="1"
            strokeDasharray="3 3"
            fill="none"
            markerEnd="url(#transform-arrow)"
          />
        )}
        <defs>
          <marker id="transform-arrow" viewBox="0 0 8 8" refX="7" refY="4" markerWidth="6" markerHeight="6" orient="auto">
            <path d="M0,0 L8,4 L0,8 z" fill="var(--accent-blue)" />
          </marker>
        </defs>

        <text x={W - PAD} y={H - PAD + 12} textAnchor="end" fontSize="9" fill="var(--text-muted)">
          {draggable ? 'drag to change the transformation' : 'original and image'}
        </text>
      </svg>

      <div className="mt-3 space-y-2">
        <Select label="Transformation" value={kind} options={KIND_OPTIONS} onChange={(v) => setKind(v as Kind)} />
        {(kind === 'rotation' || kind === 'composition') && (
          <Slider label="angle" value={angle} min={0} max={360} step={1} onChange={setAngle} format={(v) => `${fmt(v, 0)}°`} />
        )}
        {kind === 'dilation' && (
          <Slider label="scale" value={scale} min={0.25} max={2.5} step={0.05} onChange={setScale} />
        )}
        {kind === 'translation' && (
          <>
            <Slider label="dx" value={dx} min={-4} max={4} step={0.1} onChange={setDx} />
            <Slider label="dy" value={dy} min={-4} max={4} step={0.1} onChange={setDy} />
          </>
        )}
        <Toggle label="Show orbits" checked={showOrbit} onChange={setShowOrbit} />
      </div>
    </FigureFrame>
  )
}
