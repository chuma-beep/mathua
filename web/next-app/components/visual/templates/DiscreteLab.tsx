'use client'

import { useMemo, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Slider } from '../Controls'
import { discreteSchema } from '../schemas'
import { binomial, combinations, setOperation, truthTable, type SetOp } from '../math/discrete'
import type { TemplateProps } from '../types'

const W = 360
const H = 300

const SET_LABELS = ['1', '2', '3', '4', '5', '6']

export default function DiscreteLab({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => discreteSchema.parse(params), [params])
  const [mode, setMode] = useState(defaults.mode)
  const [op, setOp] = useState(defaults.op)
  const [setOpValue, setSetOpValue] = useState<SetOp>('union')
  const [leftMask, setLeftMask] = useState(defaults.leftMask)
  const [rightMask, setRightMask] = useState(defaults.rightMask)
  const [n, setN] = useState(defaults.n)
  const [k, setK] = useState(defaults.k)

  const reset = () => {
    setMode(defaults.mode)
    setOp(defaults.op)
    setLeftMask(defaults.leftMask)
    setRightMask(defaults.rightMask)
    setN(defaults.n)
    setK(defaults.k)
  }

  const table = truthTable(op)
  const sets = setOperation(setOpValue, leftMask, rightMask)
  const combos = combinations(n, k)

  return (
    <FigureFrame
      title={title ?? 'A finite structure you can enumerate'}
      hint="Change the connective, the two sets or the subset size. Every row and element is listed exhaustively."
      caption={caption ?? 'A truth table enumerates every assignment. A set operation is computed from membership. A combination lists the actual k-subsets, so the binomial count can be checked by counting them.'}
      onReset={reset}
      readout={[
        { label: 'model', value: mode },
        { label: 'truth rows', value: mode === 'truth' ? String(table.rows.length) : '—' },
        { label: 'true rows', value: mode === 'truth' ? String(table.trueCount) : '—' },
        { label: 'left set', value: sets.left.join(', ') || 'empty' },
        { label: 'right set', value: sets.right.join(', ') || 'empty' },
        { label: 'result set', value: sets.result.join(', ') || 'empty' },
        { label: 'cardinality', value: String(sets.count) },
        { label: 'C(n, k)', value: mode === 'combinations' ? `${binomial(n, k)} = ${combos.length} listed` : '—' },
      ]}
    >
      <svg
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        role="img"
        aria-label={`${mode} discrete structure`}
      >
        {mode === 'truth' && (
          <>
            <text x={30} y={28} fontSize="10" fill="var(--accent-blue)">
              {table.expression}
            </text>
            {['P', 'Q', 'result'].map((label, col) => (
              <text key={label} x={60 + col * 100} y={54} textAnchor="middle" fontSize="10" fill="var(--text-muted)">
                {label}
              </text>
            ))}
            {table.rows.map((row, index) => (
              <g key={index}>
                <text x={60} y={78 + index * 28} textAnchor="middle" fontSize="11" fill="var(--text-secondary)">
                  {row.p ? 'T' : 'F'}
                </text>
                <text x={160} y={78 + index * 28} textAnchor="middle" fontSize="11" fill="var(--text-secondary)">
                  {row.q ? 'T' : 'F'}
                </text>
                <text
                  x={260}
                  y={78 + index * 28}
                  textAnchor="middle"
                  fontSize="11"
                  fill={row.value ? 'var(--accent-green)' : 'var(--accent-red)'}
                >
                  {row.value ? 'T' : 'F'}
                </text>
              </g>
            ))}
          </>
        )}

        {mode === 'sets' && (
          <>
            <circle cx="140" cy="135" r="72" fill="var(--accent-blue)" fillOpacity="0.12" stroke="var(--accent-blue)" />
            <circle cx="220" cy="135" r="72" fill="var(--accent-green)" fillOpacity="0.12" stroke="var(--accent-green)" />
            <text x="94" y="62" fontSize="10" fill="var(--accent-blue)">A</text>
            <text x="256" y="62" fontSize="10" fill="var(--accent-green)">B</text>
            {SET_LABELS.map((label, index) => {
              const inLeft = sets.left.includes(index + 1)
              const inRight = sets.right.includes(index + 1)
              const x = inLeft && inRight ? 180 : inLeft ? 118 : inRight ? 242 : 180
              const y = 112 + (index % 3) * 22
              return (
                <text key={label} x={x} y={y} textAnchor="middle" fontSize="11" fill="var(--text-secondary)">
                  {label}
                </text>
              )
            })}
            <text x={30} y={252} fontSize="10" fill="var(--text-secondary)">
              {sets.op}: A {sets.left.join(', ') || 'empty'} | B {sets.right.join(', ') || 'empty'}
            </text>
            <text x={30} y={270} fontSize="10" fill="var(--accent-blue)">
              result {sets.result.join(', ') || 'empty'}
            </text>
          </>
        )}

        {mode === 'combinations' && (
          <>
            <text x={30} y={30} fontSize="10" fill="var(--accent-blue)">
              C({n}, {k}) = {binomial(n, k)}
            </text>
            {combos.slice(0, 18).map((subset, index) => (
              <text key={index} x={30 + (index % 2) * 165} y={56 + Math.floor(index / 2) * 22} fontSize="10" fill="var(--text-secondary)">
                {'{' + subset.join(', ') + '}'}
              </text>
            ))}
          </>
        )}
      </svg>

      <div className="mt-3 space-y-2">
        <label className="flex items-center justify-between gap-2 font-mono text-[10px] text-mathua-muted">
          <span>Model</span>
          <select
            aria-label="Model"
            value={mode}
            onChange={(event) => setMode(event.target.value as typeof mode)}
            className="h-7 border border-mathua-border bg-transparent px-1 text-[10px] text-mathua-secondary"
          >
            <option value="truth">Truth table</option>
            <option value="sets">Set operation</option>
            <option value="combinations">Combinations</option>
          </select>
        </label>
        {mode === 'truth' && (
          <label className="flex items-center justify-between gap-2 font-mono text-[10px] text-mathua-muted">
            <span>Connective</span>
            <select
              aria-label="Connective"
              value={op}
              onChange={(event) => setOp(event.target.value as typeof op)}
              className="h-7 border border-mathua-border bg-transparent px-1 text-[10px] text-mathua-secondary"
            >
              <option value="and">P and Q</option>
              <option value="or">P or Q</option>
              <option value="implies">P implies Q</option>
              <option value="xor">P xor Q</option>
            </select>
          </label>
        )}
        {mode === 'sets' && (
          <label className="flex items-center justify-between gap-2 font-mono text-[10px] text-mathua-muted">
            <span>Operation</span>
            <select
              aria-label="Set operation"
              value={sets.op}
              onChange={(event) => setSetOpValue(event.target.value as SetOp)}
              className="h-7 border border-mathua-border bg-transparent px-1 text-[10px] text-mathua-secondary"
            >
              <option value="union">Union</option>
              <option value="intersection">Intersection</option>
              <option value="difference">Difference</option>
              <option value="symmetric">Symmetric difference</option>
            </select>
          </label>
        )}
        <div className="grid grid-cols-1 gap-2">
          {mode === 'sets' && (
            <>
              <Slider label="Left mask" value={leftMask} min={0} max={63} step={1} onChange={setLeftMask} />
              <Slider label="Right mask" value={rightMask} min={0} max={63} step={1} onChange={setRightMask} />
            </>
          )}
          {mode === 'combinations' && (
            <>
              <Slider label="n" value={n} min={1} max={8} step={1} onChange={setN} />
              <Slider label="k" value={k} min={0} max={8} step={1} onChange={setK} />
            </>
          )}
        </div>
      </div>
    </FigureFrame>
  )
}
