'use client'

import { useCallback, useMemo, useRef, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { NumberField, Select, Slider } from '../Controls'
import { complexMapSchema } from '../schemas'
import { angleBetweenDirections, angleFactorAt, MAP_LABELS, mapDerivative, mapIsAnalytic, mapIsConformalAt, mapValue, mappedDirectionAngle, type ComplexMap } from '../math/complexMap'
import { clamp, scaleLinear } from '../scale'
import { fmt } from '../format'
import { useDrag } from '../useDrag'
import type { TemplateProps } from '../types'

const PAD = 12
const PANEL = 150
const GAP = 28
const LEFT = PAD
const RIGHT = LEFT + PANEL + GAP
const W = RIGHT + PANEL + PAD
const H = PANEL + 46
const TOP = 28
const DEG = 180 / Math.PI

const MAP_OPTIONS: Array<{ value: ComplexMap; label: string }> = [
  { value: 'square', label: 'z squared' },
  { value: 'exp', label: 'e to the z' },
  { value: 'inverse', label: '1 / z' },
  { value: 'conjugate', label: 'the conjugate' },
]

export default function ComplexMapping({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => {
    const p = complexMapSchema.parse(params)
    return { map: p.map, z: [Number(p.z[0]), Number(p.z[1])] as [number, number], turn: p.turn }
  }, [params])

  const [map, setMap] = useState<ComplexMap>(defaults.map)
  const [z, setZ] = useState<[number, number]>(defaults.z)
  const [turn, setTurn] = useState(defaults.turn)
  const svgRef = useRef<SVGSVGElement | null>(null)

  const image = mapValue(map, z)
  const derivative = mapDerivative(map, z)
  const analytic = mapIsAnalytic(map)

  // Two directions through z, symmetric about the positive x axis, so the
  // opening angle is exactly `turn` and the comparison is between like and like.
  const half = (turn / DEG) / 2
  const before = [-half, half]
  const after = before.map((angle) => mappedDirectionAngle(map, z, angle))
  const angleBefore = angleBetweenDirections(before[0], before[1]) * DEG
  const factor = angleFactorAt(map, z)
  // Null means the map has no value here, so there is no image angle to report.
  const angleAfter = factor === null ? null : angleBefore * factor
  const imageFinite = Number.isFinite(image[0]) && Number.isFinite(image[1])
  const afterFinite = after.every((angle) => Number.isFinite(angle))
  const conformal = mapIsConformalAt(map, z)

  const extentZ = 2.3
  // A pole has no finite image, and letting NaN into the scale turned every
  // coordinate in the right panel into NaN, which drew nothing at all.
  const reach = imageFinite ? Math.max(Math.abs(image[0]), Math.abs(image[1])) : 0
  const extentW = Math.max(1.8, reach * 1.3)

  const zX = scaleLinear(-extentZ, extentZ, LEFT, LEFT + PANEL)
  const zY = scaleLinear(extentZ, -extentZ, TOP + PANEL, TOP)
  const wX = scaleLinear(-extentW, extentW, RIGHT, RIGHT + PANEL)
  const wY = scaleLinear(extentW, -extentW, TOP + PANEL, TOP)

  const onDrag = useCallback(
    (point: { x: number; y: number }) => {
      const x = clamp(((point.x - LEFT) / PANEL) * 2 * extentZ - extentZ, -3, 3)
      const y = clamp(extentZ - ((point.y - TOP) / PANEL) * 2 * extentZ, -3, 3)
      setZ([Number(x.toFixed(2)), Number(y.toFixed(2))])
    },
    [extentZ],
  )
  useDrag(svgRef, onDrag)

  const reset = () => {
    setMap(defaults.map)
    setZ(defaults.z)
    setTurn(defaults.turn)
  }

  const ray = (x: (v: number) => number, y: (v: number) => number, from: [number, number], angle: number, color: string) => {
    const length = 0.95
    return (
      <line
        x1={x(from[0])}
        y1={y(from[1])}
        x2={x(from[0] + length * Math.cos(angle))}
        y2={y(from[1] + length * Math.sin(angle))}
        stroke={color}
        strokeWidth="1.6"
      />
    )
  }

  const arrow = (x: (v: number) => number, y: (v: number) => number, from: [number, number], to: [number, number]) => (
    <line x1={x(from[0])} y1={y(from[1])} x2={x(to[0])} y2={y(to[1])} stroke="var(--accent-blue)" strokeWidth="1.6" markerEnd="url(#map-arrow)" />
  )

  return (
    <FigureFrame
      title={title ?? 'Angles survive an analytic map'}
      hint="Drag z and change the opening angle. Watch what happens to that angle on the right."
      caption={
        caption ??
        'An analytic map rotates every direction at a point by the same amount, so the angle between two curves survives. Conjugation reflects instead, so the angle keeps its size and loses its sense.'
      }
      onReset={reset}
      readout={[
        { label: 'z', value: `(${fmt(z[0])}, ${fmt(z[1])})` },
        { label: 'f(z)', value: imageFinite ? `(${fmt(image[0])}, ${fmt(image[1])})` : 'undefined' },
        { label: '|f′(z)|', value: derivative ? fmt(Math.hypot(derivative[0], derivative[1])) : '—' },
        { label: 'angle in', value: `${fmt(angleBefore, 1)}°` },
        { label: 'angle out', value: angleAfter === null ? 'undefined' : `${fmt(angleAfter, 1)}°` },
        { label: 'conformal', value: conformal ? 'yes' : 'no' },
      ]}
    >
      <svg
        ref={svgRef}
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        style={{ touchAction: 'none' }}
        role="img"
        aria-label="A point with two directions, and their images under a complex map"
      >
        <defs>
          <marker id="map-arrow" markerWidth="6" markerHeight="6" refX="5" refY="3" orient="auto">
            <path d="M0,0 L6,3 L0,6 Z" fill="var(--accent-blue)" />
          </marker>
        </defs>

        {[
          { x: zX, y: zY, left: LEFT, title: 'z-plane' },
          { x: wX, y: wY, left: RIGHT, title: 'w = f(z)' },
        ].map((panel) => (
          <g key={panel.title}>
            <rect x={panel.left} y={TOP} width={PANEL} height={PANEL} fill="none" stroke="var(--border)" strokeWidth="0.8" />
            <line x1={panel.left} y1={panel.y(0)} x2={panel.left + PANEL} y2={panel.y(0)} stroke="var(--text-muted)" strokeWidth="0.7" />
            <line x1={panel.x(0)} y1={TOP} x2={panel.x(0)} y2={TOP + PANEL} stroke="var(--text-muted)" strokeWidth="0.7" />
            <text x={panel.left} y={TOP - 8} fontSize="9" fill="var(--text-muted)">
              {panel.title}
            </text>
          </g>
        ))}

        {before.map((angle, i) => (
          <g key={`b${i}`}>
            {ray(zX, zY, z, angle, i === 0 ? 'var(--accent-teal)' : 'var(--accent-green)')}
          </g>
        ))}
        <circle cx={zX(z[0])} cy={zY(z[1])} r="4" fill="var(--accent-blue)" />
        <text x={zX(z[0]) + 7} y={zY(z[1]) - 6} fontSize="10" fill="var(--accent-blue)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
          z
        </text>

        {imageFinite && afterFinite && after.map((angle, i) => (
          <g key={`a${i}`}>
            {ray(wX, wY, image, angle, i === 0 ? 'var(--accent-teal)' : 'var(--accent-green)')}
          </g>
        ))}
        {imageFinite && (
          <>
            <circle cx={wX(image[0])} cy={wY(image[1])} r="4" fill="var(--accent-blue)" />
            <text x={wX(image[0]) + 7} y={wY(image[1]) - 6} fontSize="10" fill="var(--accent-blue)" stroke="var(--bg)" strokeWidth="3" paintOrder="stroke">
              f(z)
            </text>
            {arrow(wX, wY, [0, 0], image)}
          </>
        )}
        {!imageFinite && (
          <text x={RIGHT + PANEL / 2} y={TOP + PANEL / 2} textAnchor="middle" fontSize="9" fill="var(--accent-teal)">
            no value here
          </text>
        )}
        {factor === 2 && (
          <text x={RIGHT} y={TOP - 8} fontSize="9" fill="var(--accent-teal)">
            critical point: the angle doubles
          </text>
        )}

        <text x={LEFT} y={TOP + PANEL + 18} fontSize="9" fill="var(--text-muted)">
          {`${fmt(angleBefore, 0)}° between the rays`}
        </text>
        <text x={RIGHT} y={TOP + PANEL + 18} fontSize="9" fill={conformal ? 'var(--accent-blue)' : 'var(--accent-teal)'}>
          {angleAfter === null ? 'no image angle' : `${fmt(angleAfter, 0)}° after the map`}
        </text>
      </svg>

      <div className="mt-3 space-y-2">
        <Select label="Map" value={map} options={MAP_OPTIONS} onChange={(v) => setMap(v as ComplexMap)} />
        <Slider label="opening" value={turn} min={10} max={170} step={1} onChange={setTurn} suffix="°" />
        <div className="grid grid-cols-2 gap-2">
          <NumberField label="x" value={z[0]} min={-3} max={3} onChange={(v) => setZ([v, z[1]])} />
          <NumberField label="y" value={z[1]} min={-3} max={3} onChange={(v) => setZ([z[0], v])} />
        </div>
      </div>
    </FigureFrame>
  )
}
