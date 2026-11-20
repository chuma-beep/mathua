'use client'

import { useMemo, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { NumberField } from '../Controls'
import { gaussianSchema } from '../schemas'
import { eliminate, outcomeText, type Row } from '../math/rowReduce'
import { fmt } from '../format'
import type { TemplateProps } from '../types'

const COLUMNS = ['x', 'y', '=']

/**
 * Gaussian elimination, one operation at a time.
 *
 * The point is not the answer but the permission to transform: each button
 * performs one elementary row operation and says what it did, because
 * `R2 ← R2 − 3R1` keeps the solution set only for a reason a learner can state.
 */
export default function GaussianElimination({ params, title, caption }: TemplateProps) {
  const defaults = useMemo<Row[]>(() => {
    const p = gaussianSchema.parse(params)
    return p.rows.map((row) => [Number(row[0]), Number(row[1]), Number(row[2])] as Row)
  }, [params])

  const [rows, setRows] = useState<Row[]>(defaults)
  const [step, setStep] = useState(0)

  const outcome = useMemo(() => eliminate(rows), [rows])
  const total = outcome.steps.length
  const shown = step === 0 ? rows : outcome.steps[step - 1].matrix
  const description =
    step === 0
      ? 'The system as it was written: two equations in x and y.'
      : outcome.steps[step - 1].description

  const setCell = (r: number, c: number, value: number) => {
    setRows((prev) =>
      prev.map((row, i) => (i === r ? (row.map((v, j) => (j === c ? value : v)) as Row) : row)),
    )
    setStep(0)
  }

  const reset = () => {
    setRows(defaults)
    setStep(0)
  }

  return (
    <FigureFrame
      title={title ?? 'Elimination, one row operation at a time'}
      hint="Step forward and read the sentence: each operation rewrites an equation without changing the solution."
      caption={
        caption ??
        'A row operation replaces one equation by a combination of the two, so any pair (x, y) that satisfied the old system still satisfies the new one. Elimination stops at a form the answer can be read from, and the last column is that answer.'
      }
      onReset={reset}
      readout={[
        { label: 'step', value: `${step} of ${total}` },
        { label: 'operation', value: step === 0 ? 'as written' : `row ${step}` },
        { label: 'outcome', value: outcomeText(outcome) },
        {
          label: 'matrix',
          value: shown.map((row) => `[${row.map((v) => fmt(v)).join('  ')}]`).join('  '),
        },
      ]}
    >
      <div className="space-y-3">
        <div className="space-y-1">
          {shown.map((row, r) => (
            <div key={r} className="flex items-center gap-2 border-x border-mathua-border px-3 py-1 font-mono text-[12px] text-mathua-primary">
              <span className="w-6 shrink-0 text-mathua-muted">{`R${r + 1}`}</span>
              {row.map((value, c) => (
                <span
                  key={c}
                  className={`w-14 text-right tabular-nums ${c === 2 ? 'border-l border-mathua-border pl-3 text-mathua-secondary' : ''}`}
                >
                  {fmt(value)}
                </span>
              ))}
            </div>
          ))}
        </div>

        <p className="min-h-[2.5rem] border border-mathua-border bg-mathua-code px-2 py-1.5 font-mono text-[10px] leading-relaxed text-mathua-secondary" data-testid="elimination-step">
          {description}
        </p>

        <div className="flex flex-wrap items-center gap-2">
          <button
            type="button"
            onClick={() => setStep((s) => Math.max(0, s - 1))}
            disabled={step === 0}
            className="border border-mathua-border px-2 py-0.5 font-mono text-[10px] text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue disabled:opacity-40"
          >
            Back
          </button>
          <button
            type="button"
            onClick={() => setStep((s) => Math.min(total, s + 1))}
            disabled={step >= total}
            className="border border-mathua-border px-2 py-0.5 font-mono text-[10px] text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue disabled:opacity-40"
          >
            Step
          </button>
          <span className="font-mono text-[10px] text-mathua-muted">{`${step} of ${total} operations`}</span>
        </div>

        <div className="grid grid-cols-3 gap-2 border-t border-mathua-border pt-3">
          {COLUMNS.map((column) => (
            <p key={column} className="font-mono text-[10px] text-mathua-muted">
              {column}
            </p>
          ))}
          {rows.map((row, r) =>
            row.map((value, c) => (
              <NumberField
                key={`${r}-${c}`}
                label={`R${r + 1} ${COLUMNS[c]}`}
                value={value}
                min={-9}
                max={9}
                step={1}
                onChange={(v) => setCell(r, c, v)}
              />
            )),
          )}
        </div>
      </div>
    </FigureFrame>
  )
}
