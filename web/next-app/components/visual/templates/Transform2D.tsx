'use client'

import { useCallback, useMemo, useRef, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { NumberField, Select, Slider, Toggle } from '../Controls'
import { transform2dSchema } from '../schemas'
import { apply, compose, dilation, identity, reflection, rotation, type Mat2, type Vec2 } from '../math/matrix'
import { maxSideDelta, polygonPath } from '../math/geometry'
import { clamp, niceTicks, pathFromPoints, scaleLinear } from '../scale'
import { XAxisTicks, YAxisTicks } from '../Axis'
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
  { value: 'reflection', label: 'Reflection' },
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

  // The polygon is the learner's, not the shard's: the authored shape is the
  // starting point, and vertices are added, removed and dragged from there. A
  // transformation of a shape somebody built is a different lesson from a
  // transformation of a triangle the figure shipped with.
  const defaults = useMemo<Vec2[]>(
    () => initial.shape.map((pair) => [Number(pair[0]), Number(pair[1])] as Vec2),
    [initial.shape],
  )
  const [shape, setShape] = useState<Vec2[]>(defaults)
  const [selected, setSelected] = useState(0)

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

  // Memoised because `extent` depends on it: a fresh `[0, 0]` on every render made
  // the frame recompute and the dependency look unstable to the hook linter.
  const offset = useMemo<Vec2>(() => (kind === 'translation' ? [dx, dy] : [0, 0]), [kind, dx, dy])
  const image = useMemo(
    () =>
      shape.map((p) => {
        const q = apply(matrix, p)
        return [q[0] + offset[0], q[1] + offset[1]] as Vec2
      }),
    [shape, matrix, offset],
  )

  // The frame is set by the authored polygon and its image, never by the live
  // shape. Deriving it from `shape` made the view zoom while a vertex was being
  // dragged: moving one point outward grew the extent, which rescaled every other
  // point under the pointer and made the drag accelerate away from the cursor.
  const extent = useMemo(() => {
    const authored = defaults.flatMap((p) => {
      const q = apply(matrix, p)
      return [Math.abs(p[0]), Math.abs(p[1]), Math.abs(q[0]) + Math.abs(offset[0]), Math.abs(q[1]) + Math.abs(offset[1])]
    })
    return Math.max(2.4, ...authored) * 1.15
  }, [defaults, matrix, offset])

  const X = scaleLinear(-extent, extent, PAD, W - PAD)
  const Y = scaleLinear(extent, -extent, PAD, H - PAD)
  const cx = X(0)
  const cy = Y(0)
  const axisTicks = niceTicks(-extent, extent, 4).filter((value) => value !== 0)

  const toData = useCallback(
    (point: { x: number; y: number }): Vec2 => [
      Math.round((((point.x - PAD) / (W - PAD * 2)) * 2 * extent - extent) * 10) / 10,
      Math.round((((H - PAD - point.y) / (H - PAD * 2)) * 2 * extent - extent) * 10) / 10,
    ],
    [extent],
  )

  const onDrag = useCallback(
    (point: { x: number; y: number }) => {
      // A vertex under the pointer takes the drag, so construction and
      // transformation share one canvas. Sixteen user units is about a finger.
      let nearest = -1
      let nearestDistance = Number.POSITIVE_INFINITY
      shape.forEach((vertex, index) => {
        const d = Math.hypot(X(vertex[0]) - point.x, Y(vertex[1]) - point.y)
        if (d < nearestDistance) {
          nearestDistance = d
          nearest = index
        }
      })
      if (nearest >= 0 && nearestDistance <= 16) {
        const [nx, ny] = toData(point)
        // Inside the frame: a vertex outside it would be drawn clipped, and the
        // learner would be dragging something they cannot see.
        const limit = extent * 0.9
        setShape((prev) => {
          const next = [...prev]
          next[nearest] = [clamp(nx, -limit, limit), clamp(ny, -limit, limit)]
          return next
        })
        setSelected(nearest)
        return
      }
      if (kind === 'rotation' || kind === 'composition') {
        const a = Math.atan2(-(point.y - cy), point.x - cx) / DEG
        setAngle((a + 360) % 360)
      } else if (kind === 'translation') {
        setDx(clamp((point.x - cx) * (2 * extent) / (W - 2 * PAD), -4, 4))
        setDy(clamp((point.y - cy) * (-2 * extent) / (H - 2 * PAD), -4, 4))
      }
    },
    [kind, cy, cx, extent, shape, X, Y, toData],
  )
  const transformsByDrag = kind === 'rotation' || kind === 'composition' || kind === 'translation'
  useDrag(svgRef, onDrag)

  const addVertex = () => {
    if (shape.length >= 8) return
    const a = shape[shape.length - 1]
    const b = shape[0]
    // Midpoint of the closing edge, nudged outward so the new vertex is visible.
    const shifted: Vec2 = [(a[0] + b[0]) / 2 + 0.5, (a[1] + b[1]) / 2 + 0.5]
    setShape([...shape, shifted])
    setSelected(shape.length)
  }

  const removeVertex = () => {
    if (shape.length <= 3) return
    setShape(shape.filter((_, index) => index !== selected))
    setSelected((prev) => Math.max(0, Math.min(prev, shape.length - 2)))
  }

  const editVertex = (axis: 0 | 1, value: number) => {
    setShape((prev) => {
      const next = [...prev]
      const current = next[selected] ?? [0, 0]
      next[selected] = axis === 0 ? [value, current[1]] : [current[0], value]
      return next
    })
  }

  const reset = () => {
    setShape(defaults)
    setSelected(0)
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
        'The dashed outline is the polygon you build, the filled outline is its image. Add or drag a vertex to make a shape, then transform it: an isometry keeps every side length, a dilation does not.'
      }
      onReset={reset}
      hint="Drag a vertex to reshape the polygon, drag empty space to change the transformation, or use the controls."
      readout={[
        { label: 'matrix', value: matrixText },
        ...(kind === 'translation' ? [{ label: 'offset', value: fmtVec(offset) }] : []),
        { label: 'image', value: image.slice(0, 3).map((point) => fmtVec(point)).join('  ') },
        { label: 'vertices', value: String(shape.length) },
        { label: 'max side change', value: fmt(delta, 3) },
        { label: 'area scale', value: fmt(Math.abs(matrix[0][0] * matrix[1][1] - matrix[0][1] * matrix[1][0])) },
      ]}
    >
      <svg
        ref={svgRef}
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] cursor-grab font-mono active:cursor-grabbing"
        style={{ touchAction: 'none' }}
        role="img"
        aria-label="A polygon and its transformed image"
      >
        {gridLines.map((l, i) => (
          <line key={i} x1={l.x1} y1={l.y1} x2={l.x2} y2={l.y2} stroke="var(--border)" strokeWidth="0.5" />
        ))}
        <line x1={PAD} y1={cy} x2={W - PAD} y2={cy} stroke="var(--text-muted)" strokeWidth="0.9" />
        <line x1={cx} y1={PAD} x2={cx} y2={H - PAD} stroke="var(--text-muted)" strokeWidth="0.9" />
        <XAxisTicks ticks={axisTicks} scale={X} at={cy} />
        <YAxisTicks ticks={axisTicks} scale={Y} at={cx} />

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

        {transformsByDrag && (
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

      </svg>

      <div className="mt-3 space-y-2">
        <Select label="Transformation" value={kind} options={KIND_OPTIONS} onChange={(v) => setKind(v as Kind)} />
        {(kind === 'rotation' || kind === 'composition') && (
          <Slider label="angle" value={angle} min={0} max={360} step={1} onChange={setAngle} suffix="°" />
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
        <div className="flex flex-wrap items-center gap-2">
          <button
            type="button"
            onClick={addVertex}
            disabled={shape.length >= 8}
            className="border border-mathua-border px-2 py-0.5 font-mono text-[10px] text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue disabled:opacity-40"
          >
            Add vertex
          </button>
          <button
            type="button"
            onClick={removeVertex}
            disabled={shape.length <= 3}
            className="border border-mathua-border px-2 py-0.5 font-mono text-[10px] text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue disabled:opacity-40"
          >
            Remove vertex
          </button>
          <span className="font-mono text-[10px] text-mathua-muted">{shape.length} vertices</span>
        </div>
        <Select
          label="Vertex"
          value={String(Math.min(selected, shape.length - 1))}
          options={shape.map((_, index) => ({ value: String(index), label: `V${index + 1}` }))}
          onChange={(value) => setSelected(Number(value))}
        />
        <div className="grid grid-cols-2 gap-2">
          <NumberField label="x" value={shape[Math.min(selected, shape.length - 1)][0]} min={-8} max={8} onChange={(value) => editVertex(0, value)} />
          <NumberField label="y" value={shape[Math.min(selected, shape.length - 1)][1]} min={-8} max={8} onChange={(value) => editVertex(1, value)} />
        </div>
      </div>
    </FigureFrame>
  )
}
