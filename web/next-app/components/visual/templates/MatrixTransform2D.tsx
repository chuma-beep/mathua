'use client'

import { useCallback, useEffect, useId, useMemo, useRef, useState } from 'react'
import FigureFrame from '../FigureFrame'
import MatrixBox from '../MatrixBox'
import { NumberField, Slider } from '../Controls'
import { matrix2dSchema } from '../schemas'
import {
  apply,
  collapseState,
  columns,
  det,
  fromFlat,
  identity,
  inverse,
  lerpMatrix,
  mul,
  type CollapseState,
  type Mat2,
  type Vec2,
} from '../math/matrix'
import { describeCollapse } from '../math/matrix'
import { fmt, fmtVec } from '../format'
import { placeLabels, type LabelItem } from '../labels'
import { useReducedMotion } from '../useReducedMotion'
import type { TemplateProps } from '../types'

// A matrix acting on the plane, explored by hand and then, if the learner asks,
// watched as a move.
//
// The primary interaction is direct: the columns of A are drawn as arrows and a
// learner drags their heads, so editing the matrix and seeing the plane are one
// gesture. The animation is secondary and never starts on its own, because a
// tween the learner did not ask for is a demonstration rather than an
// experiment. With it running the plane is drawn at M(t) = (1 - t)I + tA while
// A's own columns stay put as the destination, so the moving grid can be compared
// against the matrix it is heading for.

const SIZE = 320
const PAD = 26
// 3.2 rather than the prototype's 2.4: the authored figure starts at A v = (3, 2.2),
// and a window that clips the answer would show the lesson's own example leaving
// the frame. Drags stay inside 2, so the window is not slack the learner notices.
const EXTENT = 3.2
const UNIT = (SIZE - PAD * 2) / (EXTENT * 2)
const CX = SIZE / 2
const CY = SIZE / 2
const LIMIT = 2
const EDGE = EXTENT + 1.5
const GRID = Array.from({ length: Math.ceil(EXTENT) * 2 + 1 }, (_, i) => i - Math.ceil(EXTENT))
const DURATION_MS = 900

// In data units, so the tolerance shrinks and grows with the drawing. Big enough
// for a fingertip, small enough that two knobs 35px apart stay distinguishable:
// 0.42 units is roughly 23px here.
const PICK = 0.42

// --accent-teal is oklch(0.68 0.15 250), five degrees of hue from --accent-blue,
// so two columns drawn in those two tokens were nearly the same colour and the
// second arrow read as grey. --chart-2 is the theme's only orange and is defined
// in both schemes; the dash pattern on the second column stays as a second
// channel, because colour alone fails in grayscale.
const COL1 = 'var(--accent-blue)'
const COL2 = 'var(--chart-2)'

const X = (x: number) => CX + x * UNIT
const Y = (y: number) => CY - y * UNIT
const INV_X = (px: number) => (px - CX) / UNIT
const INV_Y = (py: number) => (CY - py) / UNIT
const clamp = (value: number, low: number, high: number) => Math.min(high, Math.max(low, value))

type Flat2 = [number, number, number, number]
type Target = 'c1' | 'c2' | 'v'

interface Challenge {
  id: string
  label: string
  passed: (matrix: Mat2) => boolean
}

const CHALLENGES: Challenge[] = [
  { id: 'area', label: 'Double the area', passed: (m) => det(m) > 1.9 },
  { id: 'flip', label: 'Flip orientation', passed: (m) => det(m) < -0.05 },
  { id: 'collapse', label: 'Collapse to a line', passed: (m) => Math.abs(det(m)) < 0.02 },
  {
    id: 'square',
    label: 'Columns at 90°',
    passed: (m) => {
      // A zero column has no direction, so the dot product is zero for the wrong
      // reason. The same defect the VectorCombination review caught.
      const [first, second] = columns(m)
      const a = Math.hypot(first[0], first[1])
      const b = Math.hypot(second[0], second[1])
      return a > 0.3 && b > 0.3 && Math.abs(first[0] * second[0] + first[1] * second[1]) < 0.05
    },
  },
]

const PRESETS: Array<{ label: string; matrix: Flat2 }> = [
  { label: 'Identity', matrix: [1, 0, 0, 1] },
  { label: 'Quarter turn', matrix: [0, -1, 1, 0] },
  { label: 'Stretch', matrix: [1.8, 0, 0, 0.6] },
  { label: 'Reflect', matrix: [1, 0, 0, -1] },
  { label: 'Shear', matrix: [1, 1, 0, 1] },
  { label: 'Project', matrix: [1, 0, 0, 0] },
]

/** Shared chip shape, so a display toggle, a preset and the step control look
 *  like the same class of thing rather than three separate inventions. */
const chip = (on: boolean) =>
  `border px-2 py-1 font-mono text-[10px] ${
    on
      ? 'border-mathua-blue text-mathua-blue'
      : 'border-mathua-border text-mathua-muted hover:border-mathua-blue hover:text-mathua-blue'
  }`

const sameVec = (a: Vec2, b: Vec2) => Math.abs(a[0] - b[0]) < 1e-9 && Math.abs(a[1] - b[1]) < 1e-9

/**
 * A determinant printed at 2 decimals rounds a genuinely thin matrix to `0`,
 * which reads as singular beside a finite inverse. Values that small keep four
 * decimals so the number agrees with the classification.
 */
const detText = (value: number) => (value !== 0 && Math.abs(value) < 0.01 ? fmt(value, 4) : fmt(value))

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
  const [showWarp, setShowWarp] = useState(true)
  // 3 is the finished picture; pressing Step wraps to 0 and uncovers the
  // combination one term at a time, so Av is something the learner builds.
  const [step, setStep] = useState(3)
  const [t, setT] = useState(1)
  const [playing, setPlaying] = useState(false)
  const [hovered, setHovered] = useState<Target | null>(null)
  const drag = useRef<Target | null>(null)
  const tRef = useRef(1)
  const reduced = useReducedMotion()
  const uid = useId().replace(/:/g, '')

  const A = fromFlat(matrix)
  const B: Mat2 | null = then ? fromFlat(then) : null
  const composed = B ? mul(B, A) : null
  const M = lerpMatrix(identity(), A, t)

  const mid = apply(A, vector)
  const image = composed ? apply(composed, vector) : mid
  const moving = apply(M, vector)

  const determinant = det(A)
  const area = Math.abs(determinant)
  const collapsed = area < 1e-9
  const orientation = collapsed ? 'collapsed' : determinant > 0 ? 'preserved' : 'reversed'
  const shape = collapseState(A, UNIT)
  const movingShape = collapseState(M, UNIT)
  const inv = inverse(A)
  const [a1, a2] = columns(A)
  const [m1, m2] = columns(M)

  // The moving legs are the same construction as A's, evaluated at M(t), so the
  // picture explains why M v is where it is at every frame rather than only at
  // the end.
  const firstLeg = apply(M, [vector[0], 0])
  const secondLeg = apply(M, [0, vector[1]])

  // Every arrow either fits the fixed window or is clipped at its edge; the
  // readout still carries the true value, and this flag says the drawing is not
  // showing all of it.
  const offscreen = [matrix[0], matrix[1], matrix[2], matrix[3], vector[0], vector[1], image[0], image[1]].some(
    (value) => Math.abs(value) > EXTENT,
  )

  const setTime = useCallback((next: number) => {
    tRef.current = next
    setT(next)
  }, [])

  // Editing A is a statement about the destination, so it stops the journey and
  // puts the figure back on A. A vector edit is not: the learner is moving the
  // thing being transformed, which is worth watching mid-flight.
  const editMatrix = useCallback(
    (next: Flat2) => {
      setMatrix(next)
      setPlaying(false)
      setTime(1)
    },
    [setTime],
  )

  // One loop, owned by this effect. Pausing, finishing, editing A or unmounting
  // all tear it down through this cleanup; nothing outside here schedules a frame.
  useEffect(() => {
    if (!playing) return
    if (reduced) {
      // No frames at all: the move is shown as its endpoint.
      setTime(1)
      setPlaying(false)
      return
    }
    let frame = 0
    let last = performance.now()
    const tick = (now: number) => {
      const next = Math.min(1, tRef.current + (now - last) / DURATION_MS)
      last = now
      setTime(next)
      if (next >= 1) {
        setPlaying(false)
        return
      }
      frame = requestAnimationFrame(tick)
    }
    frame = requestAnimationFrame(tick)
    return () => cancelAnimationFrame(frame)
  }, [playing, reduced, setTime])

  const reset = () => {
    setMatrix(defaults.matrix)
    setThen(defaults.then)
    setVector(defaults.vector)
    setShowGrid(defaults.showGrid)
    setShowBasis(defaults.showBasis)
    setShowWarp(true)
    setStep(3)
    setPlaying(false)
    setTime(1)
    setHovered(null)
  }

  const setEntry = (which: 'a' | 'b', index: number, value: number) => {
    const target = which === 'a' ? matrix : then
    if (!target) return
    const next = [...target] as Flat2
    next[index] = value
    if (which === 'a') editMatrix(next)
    else setThen(next)
  }

  const setColumn = (index: 0 | 1, next: Vec2) => {
    const flat = [...matrix] as Flat2
    if (index === 0) {
      flat[0] = next[0]
      flat[2] = next[1]
    } else {
      flat[1] = next[0]
      flat[3] = next[1]
    }
    editMatrix(flat)
  }

  const pick = (px: number, py: number): Target | null => {
    const x = INV_X(px)
    const y = INV_Y(py)
    const near = (p: Vec2) => Math.hypot(p[0] - x, p[1] - y)
    const options: Array<[Target, number]> = [
      ['c1', near(a1)],
      ['c2', near(a2)],
      ['v', near(vector)],
    ]
    options.sort((a, b) => a[1] - b[1])
    // Empty space is not a grab: pressing the background used to throw v at the
    // pointer, which made every stray click an edit.
    return options[0][1] < PICK ? options[0][0] : null
  }

  const moveTo = (px: number, py: number) => {
    const target = drag.current
    if (!target) return
    const next: Vec2 = [
      Math.round(clamp(INV_X(px), -LIMIT, LIMIT) * 10) / 10,
      Math.round(clamp(INV_Y(py), -LIMIT, LIMIT) * 10) / 10,
    ]
    if (target === 'c1') setColumn(0, next)
    else if (target === 'c2') setColumn(1, next)
    else setVector(next)
  }

  const local = (event: React.PointerEvent<SVGSVGElement>): [number, number] => {
    const rect = event.currentTarget.getBoundingClientRect()
    return [((event.clientX - rect.left) / rect.width) * SIZE, ((event.clientY - rect.top) / rect.height) * SIZE]
  }

  // What the combination is building: at t = 1 it is Av, before that it is the
  // moving M(t)v, which is why the label changes rather than the arithmetic.
  const imageVisible = t < 1 || step >= 3
  const imageTip: Vec2 = t < 1 ? moving : mid
  const imageLabel = t < 1 ? 'M(t)v' : 'Av'

  const labelItems: LabelItem[] = []
  if (composed) {
    if (showBasis && !sameVec(a1, [1, 0])) labelItems.push({ id: 'c1', text: 'A e₁', tip: a1, side: -1 })
    if (showBasis && !sameVec(a2, [0, 1])) labelItems.push({ id: 'c2', text: 'A e₂', tip: a2, side: 1 })
    labelItems.push({ id: 'v', text: 'v', tip: vector, side: -1 })
    labelItems.push({ id: 'image', text: 'Av', tip: mid, side: 1 })
    labelItems.push({ id: 'composed', text: 'BAv', tip: image, side: 1 })
  } else {
    if (imageVisible) labelItems.push({ id: 'image', text: imageLabel, tip: imageTip, side: 1 })
    // Identity suppression: a column that already is the basis vector gets no
    // label of its own, or the drawing shows `A e₁` on top of `e₁` at the start
    // of every animation.
    if (showBasis && !sameVec(a1, [1, 0])) labelItems.push({ id: 'c1', text: 'A e₁', tip: a1, side: -1 })
    if (showBasis && !sameVec(a2, [0, 1])) labelItems.push({ id: 'c2', text: 'A e₂', tip: a2, side: 1 })
    labelItems.push({ id: 'v', text: 'v', tip: vector, side: -1 })
  }
  const labels = placeLabels(labelItems, { width: SIZE, height: SIZE, pad: PAD, toX: X, toY: Y })

  const arrow = ({
    to,
    from = [0, 0],
    colour,
    marker,
    dash,
    handle = false,
    id,
  }: {
    to: Vec2
    from?: Vec2
    colour: string
    marker: string
    dash?: string
    handle?: boolean
    id?: Target
  }) => {
    const active = id !== undefined && hovered === id
    return (
      <g>
        <line
          x1={X(from[0])}
          y1={Y(from[1])}
          x2={X(to[0])}
          y2={Y(to[1])}
          stroke={colour}
          strokeWidth={active ? 2.8 : 2.2}
          strokeDasharray={dash}
          markerEnd={`url(#${marker})`}
        />
        {handle && (
          <g>
            {/* The halo appears only under the pointer. It is what tells the
                learner which of two nearby knobs is about to move. */}
            {active && <circle cx={X(to[0])} cy={Y(to[1])} r="13" fill={colour} fillOpacity="0.14" />}
            <circle cx={X(to[0])} cy={Y(to[1])} r={active ? 8.5 : 7} fill="var(--bg)" stroke={colour} strokeWidth={active ? 2.2 : 1.6} />
            <circle cx={X(to[0])} cy={Y(to[1])} r="2.4" fill={colour} />
          </g>
        )}
      </g>
    )
  }

  /** The image of the unit square, drawn according to what the matrix does to it. */
  const squareImage = (m: Mat2, state: CollapseState) => {
    const [c1, c2] = columns(m)
    const d = det(m)
    const colour = d < 0 ? 'var(--accent-red)' : COL1
    if (state === 'singular') {
      // Rank one: every point of the square lands on the segment from the origin
      // to A(e₁ + e₂), so there is no interior left to shade.
      const diagonal = [c1[0] + c2[0], c1[1] + c2[1]] as Vec2
      return (
        <line x1={X(0)} y1={Y(0)} x2={X(diagonal[0])} y2={Y(diagonal[1])} stroke="var(--accent-red)" strokeWidth="3" />
      )
    }
    const points = [[0, 0] as Vec2, c1, [c1[0] + c2[0], c1[1] + c2[1]] as Vec2, c2]
      .map(([x, y]) => `${X(x)},${Y(y)}`)
      .join(' ')
    if (state === 'pixel') {
      // Thinner than a pixel: an outline that stays visible at one pixel, and no
      // fill, because a fill would claim an area the drawing cannot show.
      return <polygon points={points} fill="none" stroke={colour} strokeWidth={1} opacity={0.9} />
    }
    return (
      <polygon
        points={points}
        fill={colour}
        fillOpacity={state === 'near' ? 0.08 : 0.16}
        stroke={colour}
        strokeOpacity={state === 'near' ? 0.35 : 0.55}
        strokeWidth={1}
      />
    )
  }

  const note = (() => {
    if (t >= 1 && shape === 'singular') return 'det A = 0: the plane collapses to a line'
    if (t < 1 && movingShape === 'singular') return 'M(t) is singular: the plane is a line here'
    if (t >= 1 && shape === 'pixel') return `det A = ${fmt(determinant, 4)}: too thin to fill, still invertible`
    if (offscreen) return 'an arrow leaves the window'
    return ''
  })()

  return (
    <FigureFrame
      title={title ?? 'Matrix as a transformation'}
      hint="Drag a knob to rewrite a column, or type A below. The green knob moves v."
      caption={
        caption ??
          'The columns are the arrows, so the matrix and the picture are the same fact. The area of the shaded square is |det A|. A⁻¹ exists exactly when det A ≠ 0, even when the image is too thin to see.'
      }
      onReset={reset}
      readout={[
        { label: 'A', value: `[[${fmt(matrix[0], 1)}, ${fmt(matrix[1], 1)}], [${fmt(matrix[2], 1)}, ${fmt(matrix[3], 1)}]]` },
        { label: 'det A', value: detText(determinant) },
        {
          label: 'A⁻¹',
          value: inv
            ? `[[${fmt(inv[0][0])},${fmt(inv[0][1])}],[${fmt(inv[1][0])},${fmt(inv[1][1])}]]`
            : 'no inverse (singular)',
        },
        { label: 'area', value: detText(area) },
        { label: 'orientation', value: orientation },
        { label: 'unit square', value: describeCollapse(shape, area) },
        { label: 'v', value: fmtVec(vector) },
        { label: 'Av', value: fmtVec(mid) },
        // The same product written as a combination of A's columns. The arrows
        // show the two scaled columns and their sum; this states the equation, so
        // a reader does not have to infer it from the picture.
        { label: 'columns', value: `Av = ${fmt(vector[0])}·(${fmt(a1[0])}, ${fmt(a1[1])}) + ${fmt(vector[1])}·(${fmt(a2[0])}, ${fmt(a2[1])})` },
        ...(composed
          ? [
              { label: 'BA', value: `[[${fmt(composed[0][0])}, ${fmt(composed[0][1])}], [${fmt(composed[1][0])}, ${fmt(composed[1][1])}]]` },
              { label: 'BAv', value: fmtVec(image) },
              { label: 'det BA', value: detText(det(composed)) },
            ]
          : []),
        // |Av| is the length of A v, which is `mid` whether or not a second matrix
        // follows. Using `image` here labelled the composed length as |Av| and hid
        // the intermediate one entirely.
        { label: '|Av|', value: fmt(Math.hypot(mid[0], mid[1])) },
        ...(composed ? [{ label: '|BAv|', value: fmt(Math.hypot(image[0], image[1])) }] : []),
      ]}
    >
      <MatrixBox matrix={A} symbol="A" />

      <div className="mt-2 grid grid-cols-2 gap-2">
        <NumberField label="a" value={matrix[0]} onChange={(v) => setEntry('a', 0, v)} />
        <NumberField label="b" value={matrix[1]} onChange={(v) => setEntry('a', 1, v)} />
        <NumberField label="c" value={matrix[2]} onChange={(v) => setEntry('a', 2, v)} />
        <NumberField label="d" value={matrix[3]} onChange={(v) => setEntry('a', 3, v)} />
      </div>

      {then && (
        <div className="mt-2 space-y-1">
          <p className="font-mono text-[10px] text-mathua-muted">then B, applied after A</p>
          <div className="grid grid-cols-2 gap-2">
            <NumberField label="a₂" value={then[0]} onChange={(v) => setEntry('b', 0, v)} />
            <NumberField label="b₂" value={then[1]} onChange={(v) => setEntry('b', 1, v)} />
            <NumberField label="c₂" value={then[2]} onChange={(v) => setEntry('b', 2, v)} />
            <NumberField label="d₂" value={then[3]} onChange={(v) => setEntry('b', 3, v)} />
          </div>
        </div>
      )}

      <svg
        viewBox={`0 0 ${SIZE} ${SIZE}`}
        className="mx-auto mt-2 block h-auto w-full max-w-[420px] touch-none font-mono"
        style={{ touchAction: 'none', cursor: hovered ? 'grab' : 'default' }}
        role="img"
        aria-label="A matrix acting on the plane, its two columns draggable"
        onPointerDown={(event) => {
          const [px, py] = local(event)
          const target = pick(px, py)
          if (!target) {
            setHovered(null)
            return
          }
          drag.current = target
          setHovered(target)
          event.currentTarget.setPointerCapture(event.pointerId)
        }}
        onPointerMove={(event) => {
          const [px, py] = local(event)
          if (drag.current) {
            moveTo(px, py)
            return
          }
          setHovered(pick(px, py))
        }}
        onPointerUp={() => {
          drag.current = null
        }}
        onPointerCancel={() => {
          drag.current = null
        }}
        onPointerLeave={() => {
          if (!drag.current) setHovered(null)
        }}
      >
        <defs>
          {[
            ['lab-1', COL1],
            ['lab-2', COL2],
            ['lab-green', 'var(--accent-green)'],
            ['lab-muted', 'var(--text-muted)'],
            ['lab-ink', 'var(--text-primary)'],
          ].map(([id, colour]) => (
            // userSpaceOnUse, not the default strokeWidth: scaled markers were
            // 17px triangles on a 2.4px stroke and read as a second arrowhead.
            <marker key={id} id={id} markerUnits="userSpaceOnUse" markerWidth="9" markerHeight="9" refX="7.5" refY="4.5" orient="auto">
              <path d="M0,0 L9,4.5 L0,9 Z" fill={colour} />
            </marker>
          ))}
          <clipPath id={`${uid}-plot`}>
            <rect x={PAD} y={PAD} width={SIZE - PAD * 2} height={SIZE - PAD * 2} />
          </clipPath>
        </defs>

        {showGrid &&
          GRID.map((k) => (
            <g key={k}>
              <line x1={X(k)} y1={PAD} x2={X(k)} y2={SIZE - PAD} stroke="var(--border-strong)" strokeWidth="0.5" opacity={showWarp ? 0.3 : 0.45} />
              <line x1={PAD} y1={Y(k)} x2={SIZE - PAD} y2={Y(k)} stroke="var(--border-strong)" strokeWidth="0.5" opacity={showWarp ? 0.3 : 0.45} />
            </g>
          ))}
        <line x1={PAD} y1={CY} x2={SIZE - PAD} y2={CY} stroke="var(--text-muted)" strokeWidth="0.8" />
        <line x1={CX} y1={PAD} x2={CX} y2={SIZE - PAD} stroke="var(--text-muted)" strokeWidth="0.8" />

        {/* The plane itself, moved by the same matrix the arrows encode. Without
            this the figure shows two arrows changing while the plane they describe
            never does, which is the part a learner cannot see. */}
        {showWarp && (
          <g clipPath={`url(#${uid}-plot)`}>
            {/* The unit square before the move, dashed, so the shaded image has
                something to be compared against. */}
            <polygon
              points={[[0, 0], [1, 0], [1, 1], [0, 1]].map(([x, y]) => `${X(x)},${Y(y)}`).join(' ')}
              fill="none"
              stroke="var(--text-muted)"
              strokeWidth="0.8"
              strokeDasharray="3 3"
              opacity="0.7"
            />
            {squareImage(M, movingShape)}
            {GRID.map((k) => {
              const down = apply(M, [k, -EDGE])
              const up = apply(M, [k, EDGE])
              const left = apply(M, [-EDGE, k])
              const right = apply(M, [EDGE, k])
              return (
                <g key={`warp-${k}`}>
                  <line x1={X(down[0])} y1={Y(down[1])} x2={X(up[0])} y2={Y(up[1])} stroke="var(--text-muted)" strokeWidth="0.7" opacity="0.24" />
                  <line x1={X(left[0])} y1={Y(left[1])} x2={X(right[0])} y2={Y(right[1])} stroke="var(--text-muted)" strokeWidth="0.7" opacity="0.24" />
                </g>
              )
            })}
          </g>
        )}

        {/* The column combination, evaluated at M(t) so it explains the moving
            image as well as the final one. */}
        {(t < 1 || step >= 1) && (
          <line
            x1={X(0)}
            y1={Y(0)}
            x2={X(firstLeg[0])}
            y2={Y(firstLeg[1])}
            stroke={COL1}
            strokeWidth="1.6"
            strokeDasharray="4 3"
            opacity={t < 1 || step >= 3 ? 0.45 : 0.9}
          />
        )}
        {(t < 1 || step >= 2) && (
          <line
            x1={X(firstLeg[0])}
            y1={Y(firstLeg[1])}
            x2={X(moving[0])}
            y2={Y(moving[1])}
            stroke={COL2}
            strokeWidth="1.6"
            strokeDasharray="4 3"
            opacity={t < 1 || step >= 3 ? 0.45 : 0.9}
          />
        )}

        {arrow({ to: [1, 0], colour: 'var(--text-muted)', marker: 'lab-muted', dash: '2 3' })}
        {arrow({ to: [0, 1], colour: 'var(--text-muted)', marker: 'lab-muted', dash: '2 3' })}

        {showBasis && (
          <>
            {arrow({ to: a1, colour: COL1, marker: 'lab-1', handle: true, id: 'c1' })}
            {arrow({ to: a2, colour: COL2, marker: 'lab-2', dash: '5 3', handle: true, id: 'c2' })}
          </>
        )}

        {arrow({ to: vector, colour: 'var(--accent-green)', marker: 'lab-green', handle: true, id: 'v' })}

        {composed ? (
          <>
            <line x1={X(vector[0])} y1={Y(vector[1])} x2={X(mid[0])} y2={Y(mid[1])} stroke={COL2} strokeWidth="1.4" strokeDasharray="4 3" markerEnd="url(#lab-2)" />
            <line x1={X(mid[0])} y1={Y(mid[1])} x2={X(image[0])} y2={Y(image[1])} stroke={COL1} strokeWidth="1.8" markerEnd="url(#lab-1)" />
          </>
        ) : (
          imageVisible && arrow({ to: imageTip, colour: 'var(--text-primary)', marker: 'lab-ink', dash: t < 1 ? '4 3' : undefined })
        )}

        {labels.map((label) => {
          const active = hovered !== null && label.id === hovered
          return (
            <text
              key={label.id}
              x={label.x}
              y={label.y}
              textAnchor="middle"
              fontSize="11"
              fontWeight={active ? 600 : 400}
              fill="var(--text-primary)"
              stroke="var(--bg)"
              strokeWidth="3"
              paintOrder="stroke"
            >
              {label.text}
            </text>
          )
        })}

        {note && (
          <text x={PAD} y={SIZE - PAD - 4} fontSize="9" fill={shape === 'singular' || movingShape === 'singular' ? 'var(--accent-red)' : 'var(--text-muted)'}>
            {note}
          </text>
        )}
      </svg>

      <div className="mt-2 flex flex-wrap items-center gap-1.5">
        <button type="button" aria-pressed={showGrid} aria-label="Show the reference grid" onClick={() => setShowGrid(!showGrid)} className={chip(showGrid)}>
          Grid
        </button>
        <button type="button" aria-pressed={showWarp} aria-label="Show the plane that A produces" onClick={() => setShowWarp(!showWarp)} className={chip(showWarp)}>
          Warp
        </button>
        <button
          type="button"
          aria-label={playing ? 'Pause' : 'Play'}
          onClick={() => {
            if (playing) {
              setPlaying(false)
              return
            }
            setTime(0)
            setPlaying(true)
          }}
          className={chip(playing)}
        >
          {playing ? 'Pause' : 'Play'}
        </button>
        <button type="button" aria-label="Show A" onClick={() => { setPlaying(false); setTime(1) }} className={chip(t >= 1)}>
          Show A
        </button>
        <button
          type="button"
          onClick={() => setStep((current) => (current + 1) % 4)}
          className={chip(step < 3)}
          aria-label="Step through the combination that builds Av"
        >
          {step === 3 ? 'Rebuild Av' : 'Next term'}
        </button>
        <span className="font-mono text-[10px] text-mathua-muted">
          {t < 1 ? `M(t), t = ${fmt(t)}` : step === 0 ? 'v only' : step === 1 ? 'x·A e₁' : step === 2 ? '+ y·A e₂' : 'Av'}
        </span>
      </div>

      {/* Scrubbing is the keyboard-and-touch path through the animation, and it is
          also the only way to inspect M(t) when reduced motion forbids playback. */}
      <div className="mt-2">
        <Slider
          label="t"
          value={t}
          min={0}
          max={1}
          step={0.01}
          onChange={(next) => {
            setPlaying(false)
            setTime(next)
          }}
        />
      </div>

      {/* v on one line instead of two stacked fields: it is one vector, and the
          plot already shows both of its numbers as a position. */}
      <div className="mt-2 flex items-center gap-1.5">
        <span className="font-mono text-[11px] text-mathua-secondary">v =</span>
        <NumberField label="vx" value={vector[0]} min={-LIMIT} max={LIMIT} onChange={(next) => setVector((previous) => [clamp(next, -LIMIT, LIMIT), previous[1]])} />
        <NumberField label="vy" value={vector[1]} min={-LIMIT} max={LIMIT} onChange={(next) => setVector((previous) => [previous[0], clamp(next, -LIMIT, LIMIT)])} />
      </div>

      {/* One row instead of six chips wrapping to two: the presets are a shortcut,
          not a control a learner uses mid-drag, and the select states which preset,
          if any, is current. */}
      <div className="mt-2 flex items-center gap-1.5">
        <span className="font-mono text-[10px] text-mathua-muted">Preset</span>
        <select
          aria-label="Preset matrix"
          value={PRESETS.find((preset) => preset.matrix.every((value, index) => Math.abs(value - matrix[index]) < 1e-9))?.label ?? ''}
          onChange={(event) => {
            const preset = PRESETS.find((item) => item.label === event.target.value)
            if (!preset) return
            editMatrix(preset.matrix)
          }}
          className="h-7 border border-mathua-border bg-transparent px-1 font-mono text-[10px] text-mathua-secondary"
        >
          {!PRESETS.some((preset) => preset.matrix.every((value, index) => Math.abs(value - matrix[index]) < 1e-9)) && <option value="">Custom</option>}
          {PRESETS.map((preset) => (
            <option key={preset.label} value={preset.label}>
              {preset.label}
            </option>
          ))}
        </select>
      </div>

      <ul className="mt-2 grid grid-cols-1 gap-x-2 gap-y-0.5 sm:grid-cols-2">
        {CHALLENGES.map((challenge) => {
          const done = challenge.passed(A)
          return (
            <li key={challenge.id} className={`font-mono text-[10px] ${done ? 'text-mathua-green' : 'text-mathua-muted'}`}>
              {done ? '✓' : '·'} {challenge.label}
            </li>
          )
        })}
      </ul>
    </FigureFrame>
  )
}
