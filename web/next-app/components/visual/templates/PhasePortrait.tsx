'use client'

import { useMemo, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Toggle } from '../Controls'
import { XAxisTicks, YAxisTicks } from '../Axis'
import { phaseSchema } from '../schemas'
import {
  phaseLine,
  phaseLineSlope,
  phaseNullclines,
  phaseSystem,
  phaseTrajectory,
  phaseVector,
  type PhaseLineModel,
  type PhaseSystemId,
} from '../math/phase'
import { scaleLinear } from '../scale'
import { fmt } from '../format'
import type { TemplateProps } from '../types'

const W = 360
const H = 320
const PAD = 32
const LINE_MIN = -1.4
const LINE_MAX = 1.9
const PLANE = 3

const LINE_LABELS: Record<PhaseLineModel, string> = {
  decay: 'Exponential decay',
  logistic: 'Logistic growth',
  bistable: 'Bistable switch',
  threshold: 'Threshold',
}

const SYSTEM_LABELS: Record<PhaseSystemId, string> = {
  saddle: 'Saddle',
  'stable-node': 'Stable node',
  spiral: 'Stable spiral',
  center: 'Centre',
}

const SEEDS: Array<[number, number]> = [
  [0.9, 0.35],
  [-0.9, 0.4],
  [0.4, 0.9],
  [-0.4, -0.9],
  [1.8, -1.1],
  [-1.8, 1.2],
]

function arrow(x: number, y: number, direction: number, color: string) {
  const size = 5
  return (
    <polygon
      points={`${x},${y} ${x - size * direction},${y - size * 0.7} ${x - size * direction},${y + size * 0.7}`}
      fill={color}
    />
  )
}

export default function PhasePortrait({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => phaseSchema.parse(params), [params])
  const [mode, setMode] = useState(defaults.mode)
  const [lineModel, setLineModel] = useState<PhaseLineModel>(defaults.lineModel)
  const [systemId, setSystemId] = useState<PhaseSystemId>(defaults.system)
  const [showTrajectories, setShowTrajectories] = useState(defaults.showTrajectories)

  const line = useMemo(() => phaseLine(lineModel), [lineModel])
  const system = useMemo(() => phaseSystem(systemId), [systemId])
  const nullclines = useMemo(() => phaseNullclines(system), [system])

  const X = scaleLinear(-PLANE, PLANE, PAD, W - PAD)
  const Y = scaleLinear(-PLANE, PLANE, H - PAD, PAD)
  const lineX = scaleLinear(LINE_MIN, LINE_MAX, PAD, W - PAD)
  const slopeMax = Math.max(1, ...[-1.2, -0.6, 0, 0.6, 1.2, 1.8].map((x) => Math.abs(phaseLineSlope(lineModel, x))))
  const lineY = scaleLinear(-slopeMax * 1.1, slopeMax * 1.1, 214, 70)
  const ticks = [-3, -2, -1, 0, 1, 2, 3]

  const curve = useMemo(() => {
    const points: Array<[number, number]> = []
    for (let index = 0; index <= 180; index++) {
      const x = LINE_MIN + ((LINE_MAX - LINE_MIN) * index) / 180
      points.push([x, phaseLineSlope(lineModel, x)])
    }
    return points
  }, [lineModel])

  const trajectories = useMemo(
    () =>
      showTrajectories
        ? SEEDS.map((seed) => phaseTrajectory(system, seed, 3.4, 130)).filter(
            (trajectory) => trajectory.points.length > 1,
          )
        : [],
    [showTrajectories, system],
  )

  const reset = () => {
    setMode(defaults.mode)
    setLineModel(defaults.lineModel)
    setSystemId(defaults.system)
    setShowTrajectories(defaults.showTrajectories)
  }

  const lineReadout = line.equilibria.map((item) => `${fmt(item.x)} ${item.stable ? 'stable' : 'unstable'}`).join(', ')
  const eigenText = Array.isArray(system.eigenvalues)
    ? system.eigenvalues.map((value) => fmt(value)).join(', ')
    : `${fmt(system.eigenvalues.real)} ± ${fmt(system.eigenvalues.imaginary)}i`

  return (
    <FigureFrame
      title={title ?? 'Stability is read from the field, not from one trajectory'}
      hint="Switch between a phase line and a phase plane. Arrows follow the sign of the derivative or the vector field, and the trajectories are integrated from that same field."
      caption={
        caption ??
        'On a phase line an equilibrium is stable when nearby arrows point toward it. In the plane the vector field is the right-hand side of the system. The eigenvalues decide whether the origin is a saddle, a node, a spiral or a centre.'
      }
      onReset={reset}
      readout={
        mode === 'line'
          ? [
              { label: 'model', value: line.label },
              { label: 'equation', value: line.equation },
              { label: 'equilibria', value: lineReadout },
              { label: 'flow', value: 'left where x′ < 0, right where x′ > 0' },
              { label: 'status', value: line.equilibria.length > 0 ? 'equilibria found' : 'no equilibrium' },
            ]
          : [
              { label: 'system', value: SYSTEM_LABELS[systemId] },
              { label: 'matrix', value: `[[${system.matrix.join(', ')}]]` },
              { label: 'trace', value: fmt(system.trace) },
              { label: 'determinant', value: fmt(system.determinant) },
              { label: 'eigenvalues', value: eigenText },
              { label: 'classification', value: system.classification },
              { label: 'stability', value: system.stability },
              { label: 'trajectories', value: showTrajectories ? String(trajectories.length) : 'hidden' },
            ]
      }
    >
      <svg
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        role="img"
        aria-label={
          mode === 'line'
            ? `Phase line for ${line.equation} with its equilibria`
            : `Phase plane for the ${SYSTEM_LABELS[systemId]} with vector field and trajectories`
        }
      >
        {mode === 'line' ? (
          <>
            <line x1={PAD} x2={W - PAD} y1={214} y2={214} stroke="var(--text-muted)" strokeWidth="0.8" />
            <line x1={lineX(0)} x2={lineX(0)} y1={62} y2={220} stroke="var(--text-muted)" strokeWidth="0.6" strokeDasharray="3 3" />
            <path
              d={curve.map(([x, y], index) => `${index === 0 ? 'M' : 'L'}${lineX(x).toFixed(2)},${lineY(y).toFixed(2)}`).join(' ')}
              fill="none"
              stroke="var(--chart-3)"
              strokeWidth="1.4"
            />
            <text x={PAD} y={40} fontSize="9" fill="var(--text-muted)">
              x′ as a function of x
            </text>
            {line.equilibria.map((equilibrium) => {
              const x = lineX(equilibrium.x)
              return (
                <g key={equilibrium.x}>
                  <circle cx={x} cy={214} r="4.5" fill={equilibrium.stable ? 'var(--accent-green)' : 'var(--accent-red)'} />
                  <text x={x} y={238} textAnchor="middle" fontSize="8" fill="var(--text-muted)">
                    {fmt(equilibrium.x)}
                  </text>
                  <text x={x} y={250} textAnchor="middle" fontSize="8" fill={equilibrium.stable ? 'var(--accent-green)' : 'var(--accent-red)'}>
                    {equilibrium.stable ? 'stable' : 'unstable'}
                  </text>
                  <circle cx={x} cy={lineY(0)} r="3" fill="var(--text-primary)" />
                </g>
              )
            })}
            {[-1.2, -0.8, -0.4, 0.2, 0.6, 1.2, 1.6].map((x) => {
              const slope = phaseLineSlope(lineModel, x)
              return arrow(lineX(x), 272, Math.sign(slope) || 1, 'var(--accent-blue)')
            })}
            <text x={PAD} y={290} fontSize="8" fill="var(--text-muted)">
              flow direction on the line
            </text>
          </>
        ) : (
          <>
            <line x1={PAD} x2={W - PAD} y1={Y(0)} y2={Y(0)} stroke="var(--text-muted)" strokeWidth="0.8" />
            <line x1={X(0)} x2={X(0)} y1={PAD} y2={H - PAD} stroke="var(--text-muted)" strokeWidth="0.8" />

            {nullclines.map((item, index) => {
              const [px, py] = item.point
              const [dx, dy] = item.direction
              const scale = 10
              return (
                <line
                  key={index}
                  x1={X(px - dx * scale)}
                  y1={Y(py - dy * scale)}
                  x2={X(px + dx * scale)}
                  y2={Y(py + dy * scale)}
                  stroke={index === 0 ? 'var(--accent-teal)' : 'var(--chart-3)'}
                  strokeWidth="1"
                  strokeDasharray="4 3"
                />
              )
            })}

            {Array.from({ length: 11 }, (_, row) => (
              <g key={row}>
                {Array.from({ length: 11 }, (_, column) => {
                  const x = -3 + column * 0.6
                  const y = -3 + row * 0.6
                  const [vx, vy] = phaseVector(system, x, y)
                  const length = Math.hypot(vx, vy)
                  if (length < 1e-9) return null
                  const scale = 0.16
                  return (
                    <line
                      key={`${row}-${column}`}
                      x1={X(x)}
                      y1={Y(y)}
                      x2={X(x + (vx / length) * scale)}
                      y2={Y(y + (vy / length) * scale)}
                      stroke="var(--text-muted)"
                      strokeWidth="0.7"
                      opacity="0.7"
                    />
                  )
                })}
              </g>
            ))}

            {trajectories.map((trajectory, index) => (
              <polyline
                key={index}
                points={trajectory.points
                  .filter(([x, y]) => Math.abs(x) <= PLANE && Math.abs(y) <= PLANE)
                  .map(([x, y]) => `${X(x)},${Y(y)}`)
                  .join(' ')}
                fill="none"
                stroke="var(--accent-blue)"
                strokeWidth="1.3"
                opacity="0.85"
              />
            ))}

            <circle cx={X(0)} cy={Y(0)} r="4" fill="var(--text-primary)" />
            <text x={PAD} y={PAD + 8} fontSize="8" fill="var(--accent-teal)">
              x′ = 0
            </text>
            <text x={PAD + 44} y={PAD + 8} fontSize="8" fill="var(--chart-3)">
              y′ = 0
            </text>
            <XAxisTicks ticks={ticks} scale={X} at={Y(0)} />
            <YAxisTicks ticks={ticks} scale={Y} at={X(0)} />
          </>
        )}
      </svg>

      <div className="mt-3 grid grid-cols-1 gap-2">
        <div className="flex items-center gap-2">
          <label htmlFor="phase-mode" className="w-24 shrink-0 font-mono text-[10px] text-mathua-muted">
            View
          </label>
          <select
            id="phase-mode"
            aria-label="View"
            value={mode}
            onChange={(event) => setMode(event.target.value as typeof mode)}
            className="h-7 min-w-0 flex-1 border border-mathua-border bg-transparent px-1 font-mono text-[10px] text-mathua-secondary"
          >
            <option value="line">Phase line</option>
            <option value="plane">Phase plane</option>
          </select>
        </div>
        {mode === 'line' ? (
          <div className="flex items-center gap-2">
            <label htmlFor="phase-line-model" className="w-24 shrink-0 font-mono text-[10px] text-mathua-muted">
              Equation
            </label>
            <select
              id="phase-line-model"
              aria-label="Equation"
              value={lineModel}
              onChange={(event) => setLineModel(event.target.value as PhaseLineModel)}
              className="h-7 min-w-0 flex-1 border border-mathua-border bg-transparent px-1 font-mono text-[10px] text-mathua-secondary"
            >
              {Object.entries(LINE_LABELS).map(([value, label]) => (
                <option key={value} value={value}>
                  {label}
                </option>
              ))}
            </select>
          </div>
        ) : (
          <>
            <div className="flex items-center gap-2">
              <label htmlFor="phase-system" className="w-24 shrink-0 font-mono text-[10px] text-mathua-muted">
                System
              </label>
              <select
                id="phase-system"
                aria-label="System"
                value={systemId}
                onChange={(event) => setSystemId(event.target.value as PhaseSystemId)}
                className="h-7 min-w-0 flex-1 border border-mathua-border bg-transparent px-1 font-mono text-[10px] text-mathua-secondary"
              >
                {Object.entries(SYSTEM_LABELS).map(([value, label]) => (
                  <option key={value} value={value}>
                    {label}
                  </option>
                ))}
              </select>
            </div>
            <Toggle label="Show trajectories" checked={showTrajectories} onChange={setShowTrajectories} />
          </>
        )}
      </div>
    </FigureFrame>
  )
}
