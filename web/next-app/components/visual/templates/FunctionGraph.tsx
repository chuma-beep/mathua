'use client'

import { useId, useMemo, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Slider, Toggle } from '../Controls'
import { XAxisTicks, YAxisTicks } from '../Axis'
import { functionGraphSchema } from '../schemas'
import {
  FUNCTION_FAMILIES,
  domainOf,
  endBehavior,
  functionFormula,
  functionValue,
  inverseFormula,
  inverseIsFunction,
  rangeOf,
  type FunctionFamily,
} from '../math/functionGraph'
import { scaleLinear } from '../scale'
import { fmt } from '../format'
import type { TemplateProps } from '../types'

const W = 360
const H = 320
const PAD = 30
const X_MIN = -4
const X_MAX = 4
const Y_MIN = -4
const Y_MAX = 4

export default function FunctionGraph({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => functionGraphSchema.parse(params), [params])
  const [family, setFamily] = useState<FunctionFamily>(defaults.family)
  const [a, setA] = useState(defaults.a)
  const [h, setH] = useState(defaults.h)
  const [k, setK] = useState(defaults.k)
  const [showInverse, setShowInverse] = useState(defaults.showInverse)
  const clipId = useId()

  const X = scaleLinear(X_MIN, X_MAX, PAD, W - PAD)
  const Y = scaleLinear(Y_MAX, Y_MIN, PAD, H - PAD)
  const ticksX = [-4, -3, -2, -1, 0, 1, 2, 3, 4]
  const ticksY = [-4, -3, -2, -1, 0, 1, 2, 3, 4]

  const formula = functionFormula(family, a, h, k)
  const ends = endBehavior(family, a)
  const inverse = inverseIsFunction(family, a)
  const inverseText = inverseFormula(family, a, h, k) ?? 'not a function on this domain'

  const curve = useMemo(() => {
    let d = ''
    let drawing = false
    for (let index = 0; index <= 240; index++) {
      const x = X_MIN + ((X_MAX - X_MIN) * index) / 240
      const y = functionValue(family, a, h, k, x)
      if (y === null || y < Y_MIN || y > Y_MAX) {
        drawing = false
        continue
      }
      d += `${drawing ? 'L' : 'M'}${X(x).toFixed(2)},${Y(y).toFixed(2)} `
      drawing = true
    }
    return d.trim()
  }, [family, a, h, k, X, Y])

  const reflected = useMemo(() => {
    if (!showInverse) return ''
    let d = ''
    let drawing = false
    for (let index = 0; index <= 240; index++) {
      const x = X_MIN + ((X_MAX - X_MIN) * index) / 240
      const y = functionValue(family, a, h, k, x)
      if (y === null || y < X_MIN || y > X_MAX || x < Y_MIN || x > Y_MAX) {
        drawing = false
        continue
      }
      d += `${drawing ? 'L' : 'M'}${X(y).toFixed(2)},${Y(x).toFixed(2)} `
      drawing = true
    }
    return d.trim()
  }, [family, a, h, k, showInverse, X, Y])

  const reset = () => {
    setFamily(defaults.family)
    setA(defaults.a)
    setH(defaults.h)
    setK(defaults.k)
    setShowInverse(defaults.showInverse)
  }

  const vertex = family === 'quadratic' ? `(${fmt(h, 2)}, ${fmt(k, 2)})` : '—'

  return (
    <FigureFrame
      title={title ?? 'The function and its inverse share a graph'}
      hint="Change the family, the scale a, the shift h and the lift k. The dashed curve is the reflection in y = x."
      caption={
        caption ??
        'A function and its inverse swap x and y, so their graphs are reflections in the line y = x. A quadratic is not one-to-one on all real x, so its reflection is not a function unless the domain is restricted.'
      }
      onReset={reset}
      readout={[
        { label: 'family', value: FUNCTION_FAMILIES.find((item) => item.id === family)?.label ?? family },
        { label: 'formula', value: formula },
        { label: 'end behavior', value: `left ${ends.left}, right ${ends.right}` },
        { label: 'domain', value: domainOf(family) },
        { label: 'range', value: rangeOf(family, a, k) },
        { label: 'vertex', value: vertex },
        { label: 'inverse', value: inverse ? 'a function' : 'not a function on this domain' },
        { label: 'inverse formula', value: inverseText },
      ]}
    >
      <svg
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        role="img"
        aria-label={`Graph of ${formula} with its inverse reflection`}
      >
        <defs>
          <clipPath id={clipId}>
            <rect x={PAD} y={PAD} width={W - PAD * 2} height={H - PAD * 2} />
          </clipPath>
        </defs>

        <line x1={PAD} y1={Y(0)} x2={W - PAD} y2={Y(0)} stroke="var(--text-muted)" strokeWidth="0.8" />
        <line x1={X(0)} y1={PAD} x2={X(0)} y2={H - PAD} stroke="var(--text-muted)" strokeWidth="0.8" />

        <g clipPath={`url(#${clipId})`}>
          <line
            x1={X(X_MIN)}
            y1={Y(X_MIN)}
            x2={X(X_MAX)}
            y2={Y(X_MAX)}
            stroke="var(--text-muted)"
            strokeWidth="0.8"
            strokeDasharray="3 3"
          />
          <path d={curve} fill="none" stroke="var(--accent-blue)" strokeWidth="2" />
          {showInverse ? (
            <path d={reflected} fill="none" stroke="var(--accent-green)" strokeWidth="1.8" strokeDasharray="5 4" />
          ) : null}
          {family === 'quadratic' ? (
            <circle cx={X(h)} cy={Y(k)} r="3.5" fill="var(--accent-blue)" />
          ) : null}
        </g>

        <XAxisTicks ticks={ticksX} scale={X} at={Y(0)} />
        <YAxisTicks ticks={ticksY} scale={Y} at={X(0)} />

        <text x={PAD + 2} y={PAD + 10} fontSize="10" fill="var(--accent-blue)">
          {formula}
        </text>
        {showInverse ? (
          <text x={W - PAD - 2} y={H - PAD - 4} fontSize="9" fill="var(--accent-green)" textAnchor="end">
            inverse
          </text>
        ) : null}
      </svg>

      <div className="mt-3 space-y-2">
        <div className="flex items-center gap-2">
          <label htmlFor="function-family" className="w-14 shrink-0 font-mono text-[10px] text-mathua-muted">
            Family
          </label>
          <select
            id="function-family"
            aria-label="Family"
            value={family}
            onChange={(event) => setFamily(event.target.value as FunctionFamily)}
            className="h-7 min-w-0 border border-mathua-border bg-transparent px-1 font-mono text-[10px] text-mathua-secondary"
          >
            {FUNCTION_FAMILIES.map((item) => (
              <option key={item.id} value={item.id}>
                {item.label}
              </option>
            ))}
          </select>
        </div>
        <Slider label="Scale a" value={a} min={-2} max={2} step={0.1} onChange={setA} />
        <Slider label="Shift h" value={h} min={-2} max={2} step={0.1} onChange={setH} />
        <Slider label="Lift k" value={k} min={-2} max={2} step={0.1} onChange={setK} />
        <Toggle label="Show inverse" checked={showInverse} onChange={setShowInverse} />
      </div>
    </FigureFrame>
  )
}
