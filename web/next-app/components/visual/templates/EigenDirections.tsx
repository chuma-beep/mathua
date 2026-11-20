'use client'

import { useCallback, useMemo, useRef, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { NumberField, Slider, Toggle } from '../Controls'
import { eigenSchema } from '../schemas'
import { angleBetween, apply, det, dot, eigen, fromFlat, normalize, trace, type Vec2 } from '../math/matrix'
import { clamp, scaleLinear } from '../scale'
import { fmt } from '../format'
import { useDrag } from '../useDrag'
import type { TemplateProps } from '../types'

const SIZE = 320
const PAD = 26
const DEG = 180 / Math.PI

export default function EigenDirections({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => {
    const p = eigenSchema.parse(params)
    return {
      matrix: [
        Number(p.matrix[0]),
        Number(p.matrix[1]),
        Number(p.matrix[2]),
        Number(p.matrix[3]),
      ] as [number, number, number, number],
      vector: [Number(p.vector[0]), Number(p.vector[1])] as Vec2,
      showCircle: p.showCircle,
    }
  }, [params])

  const [matrix, setMatrix] = useState<[number, number, number, number]>(defaults.matrix)
  const [vector, setVector] = useState<Vec2>(() => normalize(defaults.vector))
  const [showCircle, setShowCircle] = useState(defaults.showCircle)
  const svgRef = useRef<SVGSVGElement | null>(null)

  const m = fromFlat(matrix)
  const image = apply(m, vector)
  // The Rayleigh quotient. On an eigenvector it is the eigenvalue, which is why
  // it is printed beside the angle: the learner sees both settle at once.
  const quotient = dot(vector, image)
  const deviation = angleBetween(vector, image) * DEG
  // Collinear, not "within two degrees". An eigenvector whose eigenvalue is
  // negative maps to the opposite ray, which is the same line, and the old angle
  // test called that 180 degrees and refused to name it. With |v| = 1 the cross
  // product is |A v| sin(theta), so it vanishes exactly on that line and nowhere
  // else. A zero image is also an eigenvector, with eigenvalue zero.
  const imageLength = Math.hypot(image[0], image[1])
  const cross = vector[0] * image[1] - vector[1] * image[0]
  const isEigen = imageLength < 1e-9 || Math.abs(cross) < 1e-6 * imageLength

  // The frame follows the matrix, not the current image: an extent that tracked
  // |A v| rescaled the drawing on every pointer move, so the point slid away from
  // the cursor while it was being dragged.
  const rowSum = Math.max(
    Math.abs(m[0][0]) + Math.abs(m[0][1]),
    Math.abs(m[1][0]) + Math.abs(m[1][1]),
  )
  const spectrum = eigen(m)
  const extent = Math.max(1.6, rowSum * 1.15)
  const X = scaleLinear(-extent, extent, PAD, SIZE - PAD)
  const Y = scaleLinear(extent, -extent, PAD, SIZE - PAD)
  const cx = X(0)
  const cy = Y(0)

  const setEntry = (index: number, value: number) => {
    setMatrix((prev) => {
      const next = [...prev] as [number, number, number, number]
      next[index] = clamp(value, -4, 4)
      return next
    })
  }

  const onDrag = useCallback(
    (point: { x: number; y: number }) => {
      const dx = (point.x - cx) / ((SIZE - PAD * 2) / (2 * extent))
      const dy = (cy - point.y) / ((SIZE - PAD * 2) / (2 * extent))
      const unit = normalize([dx, dy])
      if (unit[0] !== 0 || unit[1] !== 0) setVector(unit)
    },
    [cx, cy, extent],
  )
  useDrag(svgRef, onDrag)

  const angle = ((Math.atan2(vector[1], vector[0]) * 180) / Math.PI + 360) % 360
  const setAngle = (degrees: number) => {
    const radians = (degrees * Math.PI) / 180
    setVector([Math.cos(radians), Math.sin(radians)])
  }

  const reset = () => {
    setMatrix(defaults.matrix)
    setVector(normalize(defaults.vector))
    setShowCircle(defaults.showCircle)
  }

  const arrow = (to: Vec2, color: string) => (
    <line
      x1={cx}
      y1={cy}
      x2={X(to[0])}
      y2={Y(to[1])}
      stroke={color}
      strokeWidth="1.8"
      markerEnd={`url(#eigen-arrow-${color.includes('blue') ? 'blue' : 'teal'})`}
    />
  )

  return (
    <FigureFrame
      title={title ?? 'Eigenvectors are directions that do not turn'}
      hint="Drag the point around the circle, or type the matrix entries."
      caption={
        caption ??
        'A v is where the matrix sends the unit vector v. When A v stays on the same line as v, the direction is an eigenvector and the length scales by the eigenvalue. The angle between them is the test.'
      }
      onReset={reset}
      readout={[
        { label: 'A', value: `[[${fmt(matrix[0], 1)}, ${fmt(matrix[1], 1)}], [${fmt(matrix[2], 1)}, ${fmt(matrix[3], 1)}]]` },
        { label: 'Av', value: `(${fmt(image[0])}, ${fmt(image[1])})` },
        { label: 'v·Av', value: fmt(quotient) },
        { label: 'turn', value: `${fmt(deviation, 1)}°` },
        { label: 'state', value: isEigen ? 'eigenvector' : 'not an eigenvector' },
        {
          label: 'eigenvalues',
          value:
            spectrum.kind === 'complex'
              ? `${fmt(spectrum.re)} ± ${fmt(spectrum.im)}i`
              : spectrum.repeated
                ? `${fmt(spectrum.values[0])} (twice)`
                : `${fmt(spectrum.values[0])}, ${fmt(spectrum.values[1])}`,
        },
        { label: 'det A', value: fmt(det(m)) },
        { label: 'trace', value: fmt(trace(m)) },
      ]}
    >
      <svg
        ref={svgRef}
        viewBox={`0 0 ${SIZE} ${SIZE}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        style={{ touchAction: 'none' }}
        role="img"
        aria-label="A unit vector, its image under a matrix, and the angle between them"
      >
        <defs>
          {[
            ['blue', 'var(--accent-blue)'],
            ['teal', 'var(--accent-teal)'],
          ].map(([name, color]) => (
            <marker key={name} id={`eigen-arrow-${name}`} markerWidth="6" markerHeight="6" refX="5" refY="3" orient="auto">
              <path d="M0,0 L6,3 L0,6 Z" fill={color} />
            </marker>
          ))}
        </defs>

        <line x1={PAD} y1={cy} x2={SIZE - PAD} y2={cy} stroke="var(--text-muted)" strokeWidth="0.8" />
        <line x1={cx} y1={PAD} x2={cx} y2={SIZE - PAD} stroke="var(--text-muted)" strokeWidth="0.8" />
        {showCircle && (
          <circle cx={cx} cy={cy} r={X(1) - cx} fill="none" stroke="var(--border-strong)" strokeDasharray="3 4" />
        )}

        {/* The line through v, so "same line" is a visible thing rather than a claim. */}
        <line
          x1={X(-vector[0] * 6)}
          y1={Y(-vector[1] * 6)}
          x2={X(vector[0] * 6)}
          y2={Y(vector[1] * 6)}
          stroke={isEigen ? 'var(--accent-green)' : 'var(--border-strong)'}
          strokeWidth="0.8"
          strokeDasharray="2 4"
        />

        {arrow(vector, 'var(--accent-blue)')}
        {arrow(image, 'var(--accent-teal)')}

        <circle cx={X(vector[0])} cy={Y(vector[1])} r="4" fill="var(--accent-blue)" />
        <circle cx={X(image[0])} cy={Y(image[1])} r="4" fill="var(--accent-teal)" />

        <text x={X(vector[0]) + 7} y={Y(vector[1]) - 7} fontSize="10" fill="var(--accent-blue)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
          v
        </text>
        <text x={X(image[0]) + 7} y={Y(image[1]) - 7} fontSize="10" fill="var(--accent-teal)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
          Av
        </text>

        {spectrum.kind === 'complex' && (
          <text x={PAD} y={SIZE - PAD - 18} fontSize="9" fill="var(--text-muted)">
            no real eigen-direction: λ = {fmt(spectrum.re)} ± {fmt(spectrum.im)}i
          </text>
        )}
        {spectrum.kind === 'real' && spectrum.repeated && (
          <text x={PAD} y={SIZE - PAD - 18} fontSize="9" fill="var(--text-muted)">
            repeated eigenvalue: one direction
          </text>
        )}
        {isEigen && (
          <text x={PAD} y={PAD + 10} fontSize="10" fill="var(--accent-green)">
            eigenvector, λ ≈ {fmt(quotient)}
          </text>
        )}
      </svg>

      <div className="mt-3 space-y-2">
        <div className="grid grid-cols-2 gap-2">
          <NumberField label="a" value={matrix[0]} onChange={(v) => setEntry(0, v)} />
          <NumberField label="b" value={matrix[1]} onChange={(v) => setEntry(1, v)} />
          <NumberField label="c" value={matrix[2]} onChange={(v) => setEntry(2, v)} />
          <NumberField label="d" value={matrix[3]} onChange={(v) => setEntry(3, v)} />
        </div>
        <Slider label="angle" value={angle} min={0} max={360} step={1} onChange={setAngle} suffix="°" />
        <Toggle label="Show the unit circle" checked={showCircle} onChange={setShowCircle} />
      </div>
    </FigureFrame>
  )
}
