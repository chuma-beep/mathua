'use client'

import { useMemo, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Slider } from '../Controls'
import { fractionLineSchema } from '../schemas'
import { addFractions, compareFractions, fractionValue, simplify } from '../math/fraction'
import { fmt } from '../format'
import type { TemplateProps } from '../types'

const W = 360
const H = 300
const LEFT = 42
const RIGHT = 318

export default function FractionLine({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => fractionLineSchema.parse(params), [params])
  const [numerator, setNumerator] = useState(defaults.numerator)
  const [denominator, setDenominator] = useState(defaults.denominator)
  const [otherNumerator, setOtherNumerator] = useState(defaults.otherNumerator)
  const [otherDenominator, setOtherDenominator] = useState(defaults.otherDenominator)

  const first = { numerator, denominator }
  const second = { numerator: otherNumerator, denominator: otherDenominator }
  const firstValue = fractionValue(first)
  const secondValue = fractionValue(second)
  const domain = Math.max(1, Math.ceil(firstValue), Math.ceil(secondValue))
  const comparison = compareFractions(first, second)
  const sum = addFractions(first, second)
  const firstSimplified = simplify(first)
  const secondSimplified = simplify(second)
  const same = comparison === 0

  const reset = () => {
    setNumerator(defaults.numerator)
    setDenominator(defaults.denominator)
    setOtherNumerator(defaults.otherNumerator)
    setOtherDenominator(defaults.otherDenominator)
  }

  const x = (value: number) => LEFT + (Math.min(domain, Math.max(0, value)) / domain) * (RIGHT - LEFT)

  const numberLine = (value: number, denominatorValue: number, y: number, color: string, key: string) => {
    const ticks: Array<{ x: number; whole: boolean }> = []
    const count = Math.min(96, denominatorValue * domain)
    for (let index = 0; index <= count; index++) {
      const at = index / denominatorValue
      ticks.push({ x: x(at), whole: Math.abs(at - Math.round(at)) < 1e-9 })
    }
    return (
      <g key={key}>
        <line x1={LEFT} y1={y} x2={RIGHT} y2={y} stroke="var(--text-muted)" strokeWidth="1" />
        {ticks.map((tick, index) => (
          <line
            key={index}
            x1={tick.x}
            y1={y - (tick.whole ? 4 : 2)}
            x2={tick.x}
            y2={y + (tick.whole ? 4 : 2)}
            stroke="var(--text-muted)"
            strokeWidth={tick.whole ? 1 : 0.6}
          />
        ))}
        <circle cx={x(value)} cy={y} r="4" fill={color} />
        <text x={x(value)} y={y - 9} textAnchor="middle" fontSize="9" fill={color}>
          {value.toFixed(3).replace(/0+$/, '').replace(/\.$/, '')}
        </text>
      </g>
    )
  }

  const bar = (
    value: number,
    denominatorValue: number,
    y: number,
    color: string,
    key: string,
  ) => {
    const width = (RIGHT - LEFT) / domain
    const segment = width / denominatorValue
    const count = Math.min(denominatorValue, Math.ceil(value * denominatorValue))
    return (
      <g key={key}>
        {Array.from({ length: count }, (_, index) => (
          <rect
            key={index}
            x={LEFT + index * segment}
            y={y}
            width={Math.max(0, segment - 1)}
            height="22"
            fill={color}
            opacity={index < value * denominatorValue - 1e-9 ? 0.85 : 0.25}
            stroke="var(--bg)"
            strokeWidth="1"
          />
        ))}
        <rect
          x={LEFT}
          y={y}
          width={Math.min(RIGHT - LEFT, value * width)}
          height="22"
          fill="none"
          stroke={color}
          strokeWidth="1"
        />
      </g>
    )
  }

  return (
    <FigureFrame
      title={title ?? 'The same amount, two ways'}
      hint="Move the numerators and denominators. The number line places each fraction, and the bars split the same whole into its denominator's parts."
      caption={
        caption ??
        'A fraction is a point on the number line and an area of a whole. Dividing every part into more equal pieces changes the numerator and denominator but not the point, which is why two different fraction symbols can name the same amount.'
      }
      onReset={reset}
      readout={[
        { label: 'first fraction', value: `${numerator}/${denominator}` },
        { label: 'first decimal', value: fmt(firstValue, 4) },
        { label: 'first simplified', value: `${firstSimplified.numerator}/${firstSimplified.denominator}` },
        { label: 'second fraction', value: `${otherNumerator}/${otherDenominator}` },
        { label: 'second decimal', value: fmt(secondValue, 4) },
        { label: 'second simplified', value: `${secondSimplified.numerator}/${secondSimplified.denominator}` },
        { label: 'comparison', value: comparison < 0 ? '<' : comparison > 0 ? '>' : '=' },
        { label: 'equivalent', value: same ? 'yes' : 'no' },
        { label: 'sum', value: `${sum.numerator}/${sum.denominator}` },
        { label: 'common denominator', value: String(sum.denominator) },
      ]}
    >
      <svg
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        role="img"
        aria-label={`Number line and fraction bars comparing ${numerator}/${denominator} with ${otherNumerator}/${otherDenominator}`}
      >
        <text x={LEFT} y={22} fontSize="8" fill="var(--text-primary)">
          {`${numerator}/${denominator}`}
        </text>
        {numberLine(firstValue, denominator, 58, 'var(--accent-blue)', 'first-line')}
        <text x={LEFT} y={108} fontSize="8" fill="var(--text-primary)">
          {`${otherNumerator}/${otherDenominator}`}
        </text>
        {numberLine(secondValue, otherDenominator, 144, 'var(--accent-green)', 'second-line')}

        {bar(firstValue, denominator, 188, 'var(--accent-blue)', 'first-bar')}
        {bar(secondValue, otherDenominator, 228, 'var(--accent-green)', 'second-bar')}

        <text x={LEFT} y={176} fontSize="8" fill="var(--accent-blue)">
          {`${numerator} of ${denominator} equal parts`}
        </text>
        <text x={LEFT} y={216} fontSize="8" fill="var(--accent-green)">
          {`${otherNumerator} of ${otherDenominator} equal parts`}
        </text>
        <text x={LEFT} y={278} fontSize="9" fill="var(--text-primary)">
          {`${numerator}/${denominator} ${comparison < 0 ? '<' : comparison > 0 ? '>' : '='} ${otherNumerator}/${otherDenominator}`}
        </text>
      </svg>

      <div className="mt-3 grid grid-cols-1 gap-1.5">
        <Slider label="Numerator" value={numerator} min={0} max={8} step={1} onChange={setNumerator} />
        <Slider label="Denominator" value={denominator} min={1} max={8} step={1} onChange={setDenominator} />
        <Slider
          label="Other numerator"
          value={otherNumerator}
          min={0}
          max={8}
          step={1}
          onChange={setOtherNumerator}
        />
        <Slider
          label="Other denominator"
          value={otherDenominator}
          min={1}
          max={8}
          step={1}
          onChange={setOtherDenominator}
        />
      </div>
    </FigureFrame>
  )
}
