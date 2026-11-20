'use client'

import { useMemo, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Slider } from '../Controls'
import { conditionalAreaSchema } from '../schemas'
import { conditionalArea } from '../math/probability'
import { fmt } from '../format'
import type { TemplateProps } from '../types'

const SIZE = 116
const LEFT_X = 40
const RIGHT_X = 204
const TOP = 72
const BOTTOM = TOP + SIZE

export default function ConditionalArea({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => conditionalAreaSchema.parse(params), [params])
  const [pB, setPB] = useState(defaults.pB)
  const [pAgivenB, setPAgivenB] = useState(defaults.pAgivenB)
  const [pAgivenBc, setPAgivenBc] = useState(defaults.pAgivenBc)

  const model = conditionalArea(pB, pAgivenB, pAgivenBc)

  const reset = () => {
    setPB(defaults.pB)
    setPAgivenB(defaults.pAgivenB)
    setPAgivenBc(defaults.pAgivenBc)
  }

  if (!model) {
    return (
      <FigureFrame title={title ?? 'Conditioning shrinks the sample space'} readout={[]}>
        <p className="font-mono text-[10px] text-mathua-muted">Enter finite probabilities.</p>
      </FigureFrame>
    )
  }

  const bWidth = SIZE * model.pB
  const bcWidth = SIZE - bWidth
  const bBand = SIZE * model.pAgivenB
  const bcBand = SIZE * model.pAgivenBc

  return (
    <FigureFrame
      title={title ?? 'Conditioning shrinks the sample space'}
      hint="Move the three sliders. The left square is the whole sample space split at P(B); the right square is the column for B stretched to full width, which is what conditioning on B does."
      caption={
        caption ??
        'Area is probability. In the left square the B column has width P(B) and the event A is a band inside each column at height P(A|B) or P(A|B^c). The band inside the B column has area P(A ∩ B). Conditioning on B rescales that column to the full square, so the band height becomes P(A|B). Equal heights in the two columns means independence.'
      }
      onReset={reset}
      readout={[
        { label: 'P(B)', value: fmt(model.pB, 3) },
        { label: 'P(Bᶜ)', value: fmt(model.pBc, 3) },
        { label: 'P(A | B)', value: fmt(model.pAgivenB, 3) },
        { label: 'P(A | Bᶜ)', value: fmt(model.pAgivenBc, 3) },
        { label: 'P(A ∩ B)', value: fmt(model.pAandB, 4) },
        { label: 'P(A ∩ Bᶜ)', value: fmt(model.pAandBc, 4) },
        { label: 'P(A)', value: fmt(model.pA, 4) },
        {
          label: 'independence',
          value: model.independent ? 'holds: the bands are level' : 'fails: the bands differ',
        },
      ]}
    >
      <svg
        viewBox={`0 0 360 300`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        role="img"
        aria-label={`Unit square split at P(B) = ${fmt(model.pB, 2)} with the event A as a band of height P(A|B) = ${fmt(model.pAgivenB, 2)}`}
      >
        <text x={LEFT_X} y={TOP - 30} fontSize="9" fill="var(--text-muted)">
          sample space
        </text>
        <text x={RIGHT_X} y={TOP - 30} fontSize="9" fill="var(--text-muted)">
          conditioned on B
        </text>

        <rect x={LEFT_X} y={TOP} width={SIZE} height={SIZE} fill="var(--bg)" stroke="var(--text-muted)" strokeWidth="1" />
        <rect x={LEFT_X} y={TOP} width={bWidth} height={SIZE} fill="var(--accent-teal)" fillOpacity="0.08" />
        <rect x={LEFT_X} y={BOTTOM - bBand} width={bWidth} height={bBand} fill="var(--accent-blue)" fillOpacity="0.35" />
        <rect
          x={LEFT_X + bWidth}
          y={BOTTOM - bcBand}
          width={bcWidth}
          height={bcBand}
          fill="var(--accent-red)"
          fillOpacity="0.28"
        />
        <line x1={LEFT_X + bWidth} x2={LEFT_X + bWidth} y1={TOP} y2={BOTTOM} stroke="var(--text-muted)" strokeWidth="1.2" />
        <text x={LEFT_X + bWidth / 2} y={TOP - 8} textAnchor="middle" fontSize="8" fill="var(--accent-teal)">
          B
        </text>
        <text x={LEFT_X + bWidth + bcWidth / 2} y={TOP - 8} textAnchor="middle" fontSize="8" fill="var(--text-muted)">
          Bᶜ
        </text>
        <text x={LEFT_X + bWidth / 2} y={BOTTOM - bBand / 2 + 3} textAnchor="middle" fontSize="7" fill="var(--bg)">
          {fmt(model.pAandB, 3)}
        </text>
        <text x={LEFT_X + bWidth + bcWidth / 2} y={BOTTOM - bcBand / 2 + 3} textAnchor="middle" fontSize="7" fill="var(--text-primary)">
          {fmt(model.pAandBc, 3)}
        </text>

        <path
          d={`M ${LEFT_X + SIZE + 8} ${TOP + SIZE / 2} L ${RIGHT_X - 10} ${TOP + SIZE / 2}`}
          stroke="var(--text-muted)"
          strokeWidth="1.2"
          markerEnd="url(#cond-arrow)"
        />
        <defs>
          <marker id="cond-arrow" viewBox="0 0 8 8" refX="7" refY="4" markerWidth="6" markerHeight="6" orient="auto">
            <path d="M0,0 L8,4 L0,8 z" fill="var(--text-muted)" />
          </marker>
        </defs>

        <rect x={RIGHT_X} y={TOP} width={SIZE} height={SIZE} fill="var(--bg)" stroke="var(--text-muted)" strokeWidth="1" />
        <rect x={RIGHT_X} y={BOTTOM - SIZE * model.pAgivenB} width={SIZE} height={SIZE * model.pAgivenB} fill="var(--accent-blue)" fillOpacity="0.35" />
        <line
          x1={RIGHT_X}
          x2={RIGHT_X + SIZE}
          y1={BOTTOM - SIZE * model.pAgivenB}
          y2={BOTTOM - SIZE * model.pAgivenB}
          stroke="var(--accent-blue)"
          strokeWidth="1.4"
        />
        <text x={RIGHT_X + SIZE / 2} y={BOTTOM - (SIZE * model.pAgivenB) / 2 + 3} textAnchor="middle" fontSize="8" fill="var(--bg)">
          P(A | B) = {fmt(model.pAgivenB, 2)}
        </text>
        <text x={RIGHT_X + SIZE / 2} y={BOTTOM + 14} textAnchor="middle" fontSize="8" fill="var(--text-muted)">
          area P(B) × P(A | B)
        </text>
        <text x={LEFT_X + SIZE / 2} y={BOTTOM + 14} textAnchor="middle" fontSize="8" fill="var(--text-muted)">
          area 1
        </text>
      </svg>

      <div className="mt-3 grid grid-cols-1 gap-2">
        <Slider label="P(B)" value={pB} min={0.1} max={0.9} step={0.01} onChange={setPB} />
        <Slider label="P(A | B)" value={pAgivenB} min={0} max={1} step={0.01} onChange={setPAgivenB} />
        <Slider label="P(A | Bᶜ)" value={pAgivenBc} min={0} max={1} step={0.01} onChange={setPAgivenBc} />
      </div>
    </FigureFrame>
  )
}
