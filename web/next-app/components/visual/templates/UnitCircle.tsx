'use client'

import { useMemo, useRef, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Slider, Toggle } from '../Controls'
import { unitCircleSchema } from '../schemas'
import {
  TAU,
  degreesToRadians,
  projectToUnitCircle,
  radiansToDegrees,
  unitPoint,
} from '../math/trig'
import { arcPath, clamp, invertLinear, pathFromPoints, scaleLinear } from '../scale'
import { fmt } from '../format'
import { useDrag } from '../useDrag'
import type { TemplateProps } from '../types'

const W = 340
const H = 446

const CX = 170
const CY = 138
const R = 82

const PLOT_X0 = 30
const PLOT_X1 = 310
const PLOT_Y0 = 30
const PLOT_Y1 = 246

const GRAPH_X0 = 30
const GRAPH_X1 = 310
const GRAPH_Y0 = 292
const GRAPH_Y1 = 414

const COMMON_ANGLES = [0, 30, 45, 60, 90, 120, 135, 150, 180, 210, 225, 240, 270, 300, 315, 330]

const angleLabel = (degrees: number) => `${fmt(degrees, 1)}°`

export default function UnitCircle({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => unitCircleSchema.parse(params), [params])
  const [angle, setAngle] = useState(defaults.angle)
  const [showSin, setShowSin] = useState(defaults.showSin)
  const [showCos, setShowCos] = useState(defaults.showCos)
  const svgRef = useRef<SVGSVGElement | null>(null)

  const point = useMemo(() => unitPoint(degreesToRadians(angle)), [angle])
  const theta = point.theta
  const normalizedDegrees = radiansToDegrees(theta)
  const selectedAngle = COMMON_ANGLES.find((value) => Math.abs(value - normalizedDegrees) < 0.5)
  const tanValue = point.tan === null ? 'undefined' : fmt(point.tan, 3)

  const toX = invertLinear(-1.25, 1.25, PLOT_X0, PLOT_X1)
  const toY = invertLinear(1.25, -1.25, PLOT_Y0, PLOT_Y1)
  const circleX = scaleLinear(-1.25, 1.25, PLOT_X0, PLOT_X1)
  const circleY = scaleLinear(1.25, -1.25, PLOT_Y0, PLOT_Y1)
  const graphX = scaleLinear(0, TAU, GRAPH_X0, GRAPH_X1)
  const graphY = scaleLinear(1.1, -1.1, GRAPH_Y0, GRAPH_Y1)

  const reset = () => {
    setAngle(defaults.angle)
    setShowSin(defaults.showSin)
    setShowCos(defaults.showCos)
  }

  const onDrag = (pointer: { x: number; y: number }) => {
    const target: [number, number] = [toX(pointer.x), toY(pointer.y)]
    const next = projectToUnitCircle(target)
    setAngle(clamp(radiansToDegrees(next), -360, 360))
  }
  useDrag(svgRef, onDrag)

  const pathFor = (fn: (value: number) => number) => {
    const points: Array<[number, number]> = []
    for (let index = 0; index <= 180; index++) {
      const value = (index / 180) * TAU
      points.push([graphX(value), graphY(fn(value))])
    }
    return pathFromPoints(points)
  }

  const sinePath = pathFor(Math.sin)
  const cosinePath = pathFor(Math.cos)

  const pointX = CX + point.cos * R
  const pointY = CY - point.sin * R
  const projectionX = CX + point.cos * R
  const projectionY = CY - point.sin * R

  return (
    <FigureFrame
      title={title ?? 'One angle, one point, two graphs'}
      hint="Drag the point around the circle, or change the angle. Sine is the height, cosine is the horizontal coordinate, and both reappear on the graph below."
      caption={
        caption ??
        'A point on the unit circle has coordinates (cos θ, sin θ). The angle turns counterclockwise from the positive x-axis, and the same angle marks the sine and cosine values on the graph below. Adding 360° returns to the same point.'
      }
      onReset={reset}
      controls={
        <div className="space-y-2">
          <Slider label="Angle" value={angle} min={-360} max={360} step={1} onChange={setAngle} suffix="°" />
          <div className="flex items-center gap-2">
            <label htmlFor="unit-circle-common-angle" className="shrink-0 font-mono text-xs text-mathua-muted">
              Common angle
            </label>
            <select
              id="unit-circle-common-angle"
              aria-label="Common angle"
              value={selectedAngle === undefined ? '' : String(selectedAngle)}
              onChange={(event) => {
                if (event.target.value !== '') setAngle(Number(event.target.value))
              }}
              className="h-10 min-w-0 border border-mathua-border bg-transparent px-2 font-mono text-xs text-mathua-secondary"
            >
              {selectedAngle === undefined && <option value="">Custom</option>}
              {COMMON_ANGLES.map((value) => (
                <option key={value} value={value}>
                  {value}°
                </option>
              ))}
            </select>
          </div>
          <div className="grid grid-cols-2 gap-2">
            <Toggle label="Show sine" checked={showSin} onChange={setShowSin} />
            <Toggle label="Show cosine" checked={showCos} onChange={setShowCos} />
          </div>
        </div>
      }
      readout={[
        { label: 'angle', value: angleLabel(angle) },
        { label: 'angle on circle', value: angleLabel(normalizedDegrees) },
        { label: 'radians', value: `${fmt(theta, 3)} rad` },
        { label: 'point', value: `(${fmt(point.cos, 3)}, ${fmt(point.sin, 3)})` },
        { label: 'cos θ', value: fmt(point.cos, 4) },
        { label: 'sin θ', value: fmt(point.sin, 4) },
        { label: 'tan θ', value: tanValue },
        { label: 'quadrant', value: point.quadrant === 'axis' ? 'on an axis' : point.quadrant },
        { label: 'reference angle', value: angleLabel(radiansToDegrees(point.referenceAngle)) },
        { label: 'sin² + cos²', value: fmt(point.sin * point.sin + point.cos * point.cos, 6) },
      ]}
    >
      <svg
        ref={svgRef}
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        style={{ touchAction: 'none' }}
        role="img"
        aria-label={`Unit circle at ${angleLabel(normalizedDegrees)} with its sine and cosine graph`}
      >
        <line x1={PLOT_X0} y1={CY} x2={PLOT_X1} y2={CY} stroke="var(--text-muted)" strokeWidth="0.8" />
        <line x1={CX} y1={PLOT_Y0} x2={CX} y2={PLOT_Y1} stroke="var(--text-muted)" strokeWidth="0.8" />
        <circle cx={CX} cy={CY} r={R} fill="none" stroke="var(--text-muted)" strokeWidth="1" />

        {[-1, 1].map((value) => (
          <g key={value}>
            <line
              x1={circleX(value)}
              y1={CY - 3}
              x2={circleX(value)}
              y2={CY + 3}
              stroke="var(--text-muted)"
              strokeWidth="0.8"
            />
            <text x={circleX(value)} y={CY + 14} textAnchor="middle" fontSize="8" fill="var(--text-muted)">
              {value}
            </text>
          </g>
        ))}
        {[-1, 1].map((value) => (
          <g key={value}>
            <line
              x1={CX - 3}
              y1={circleY(value)}
              x2={CX + 3}
              y2={circleY(value)}
              stroke="var(--text-muted)"
              strokeWidth="0.8"
            />
            <text x={CX - 10} y={circleY(value) + 3} textAnchor="end" fontSize="8" fill="var(--text-muted)">
              {value}
            </text>
          </g>
        ))}

        {theta > 1e-9 ? (
          <path d={arcPath(CX, CY, 28, theta)} fill="none" stroke="var(--accent-green)" strokeWidth="1.5" data-part="angle-arc" />
        ) : null}

        <line
          x1={CX}
          y1={CY}
          x2={pointX}
          y2={pointY}
          stroke="var(--text-primary)"
          strokeWidth="1.4"
        />

        {showCos ? (
          <line
            x1={pointX}
            y1={pointY}
            x2={projectionX}
            y2={CY}
            stroke="var(--accent-teal)"
            strokeWidth="1.2"
            strokeDasharray="4 3"
            data-part="cosine-projection"
          />
        ) : null}
        {showSin ? (
          <line
            x1={pointX}
            y1={pointY}
            x2={CX}
            y2={projectionY}
            stroke="var(--accent-blue)"
            strokeWidth="1.2"
            strokeDasharray="4 3"
            data-part="sine-projection"
          />
        ) : null}

        <circle cx={pointX} cy={pointY} r="4.5" fill="var(--accent-green)" data-part="circle-point" />
        <text x={pointX + 8} y={pointY - 7} fontSize="9" fill="var(--text-primary)">
          P
        </text>

        <text x={CX + 20} y={CY - 10} fontSize="9" fill="var(--accent-green)">
          θ
        </text>
        <text x={PLOT_X0 + 2} y={PLOT_Y0 + 10} fontSize="9" fill="var(--text-muted)">
          unit circle
        </text>

        <line x1={GRAPH_X0} y1={graphY(0)} x2={GRAPH_X1} y2={graphY(0)} stroke="var(--text-muted)" strokeWidth="0.8" />
        <line x1={GRAPH_X0} y1={GRAPH_Y0} x2={GRAPH_X0} y2={GRAPH_Y1} stroke="var(--text-muted)" strokeWidth="0.8" />

        {showSin ? (
          <path d={sinePath} fill="none" stroke="var(--accent-blue)" strokeWidth="1.5" data-part="sine-curve" />
        ) : null}
        {showCos ? (
          <path d={cosinePath} fill="none" stroke="var(--accent-teal)" strokeWidth="1.5" data-part="cosine-curve" />
        ) : null}

        <line
          x1={graphX(theta)}
          y1={GRAPH_Y0}
          x2={graphX(theta)}
          y2={GRAPH_Y1}
          stroke="var(--text-muted)"
          strokeWidth="0.8"
          strokeDasharray="3 3"
        />
        {showSin ? (
          <circle
            cx={graphX(theta)}
            cy={graphY(point.sin)}
            r="3.5"
            fill="var(--accent-blue)"
            data-part="sine-marker"
          />
        ) : null}
        {showCos ? (
          <circle
            cx={graphX(theta)}
            cy={graphY(point.cos)}
            r="3.5"
            fill="var(--accent-teal)"
            data-part="cosine-marker"
          />
        ) : null}

        {[0, Math.PI / 2, Math.PI, (3 * Math.PI) / 2, TAU].map((value, index) => (
          <g key={value}>
            <line
              x1={graphX(value)}
              y1={GRAPH_Y1}
              x2={graphX(value)}
              y2={GRAPH_Y1 + 4}
              stroke="var(--text-muted)"
              strokeWidth="0.8"
            />
            <text x={graphX(value)} y={GRAPH_Y1 + 15} textAnchor="middle" fontSize="8" fill="var(--text-muted)">
              {['0', 'π/2', 'π', '3π/2', '2π'][index]}
            </text>
          </g>
        ))}

        <text x={GRAPH_X0 + 2} y={GRAPH_Y0 + 8} fontSize="9" fill="var(--accent-blue)">
          sin
        </text>
        <text x={GRAPH_X0 + 26} y={GRAPH_Y0 + 8} fontSize="9" fill="var(--accent-teal)">
          cos
        </text>
      </svg>

    </FigureFrame>
  )
}
