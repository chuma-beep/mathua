'use client'

import { useMemo, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Slider } from '../Controls'
import { XAxisTicks } from '../Axis'
import { numberLineSchema } from '../schemas'
import { balanceState, compareSigned, percentChange, ratioForm, tenthsValue, type ArithmeticMode } from '../math/arith'
import { niceTicks, scaleLinear } from '../scale'
import { fmt } from '../format'
import type { TemplateProps } from '../types'

const W = 360
const H = 300
const PAD = 30

const MODE_LABELS: Record<ArithmeticMode, string> = {
  signed: 'Signed number line',
  decimal: 'Tenths on a line',
  percent: 'Percentage change',
  ratio: 'Ratio tape',
  balance: 'Balance equation',
}

const MODE_TITLES: Record<ArithmeticMode, string> = {
  signed: 'Left of zero is smaller',
  decimal: 'Ten equal parts make one',
  percent: 'A percent is a part of the base',
  ratio: 'The same comparison, simplified',
  balance: 'Do the same to both sides',
}

export default function NumberLineLab({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => numberLineSchema.parse(params), [params])
  const [mode, setMode] = useState<ArithmeticMode>(defaults.mode)
  const [a, setA] = useState(defaults.a)
  const [b, setB] = useState(defaults.b)
  const [tenths, setTenths] = useState(defaults.tenths)
  const [percent, setPercent] = useState(defaults.percent)
  const [base, setBase] = useState(defaults.base)
  const [ratioA, setRatioA] = useState(defaults.ratioA)
  const [ratioB, setRatioB] = useState(defaults.ratioB)
  const [guess, setGuess] = useState(defaults.guess)
  const [constant, setConstant] = useState(defaults.constant)
  const [total, setTotal] = useState(defaults.total)

  const reset = () => {
    setMode(defaults.mode)
    setA(defaults.a)
    setB(defaults.b)
    setTenths(defaults.tenths)
    setPercent(defaults.percent)
    setBase(defaults.base)
    setRatioA(defaults.ratioA)
    setRatioB(defaults.ratioB)
    setGuess(defaults.guess)
    setConstant(defaults.constant)
    setTotal(defaults.total)
  }

  const signed = compareSigned(a, b)
  const decimal = tenthsValue(tenths)
  const change = percentChange(percent, base)
  const ratio = ratioForm(ratioA, ratioB)
  const balance = balanceState(guess, constant, total)

  const readout =
    mode === 'signed'
      ? [
          { label: 'A', value: fmt(a) },
          { label: 'B', value: fmt(b) },
          { label: 'comparison', value: signed.sentence },
          { label: 'leftmost', value: fmt(signed.left) },
          { label: 'rightmost', value: fmt(signed.right) },
        ]
      : mode === 'decimal'
        ? [
            { label: 'tenths', value: `${tenths}/10` },
            { label: 'decimal', value: fmt(decimal, 1) },
            { label: 'fraction', value: `${tenths}/10` },
          ]
        : mode === 'percent'
          ? [
              { label: 'percent', value: `${fmt(percent)}%` },
              { label: 'base', value: fmt(base) },
              { label: 'change', value: fmt(change.amount) },
              { label: 'result', value: fmt(change.result) },
              { label: 'direction', value: change.direction },
            ]
          : mode === 'ratio'
            ? [
                { label: 'ratio', value: `${ratio.a} : ${ratio.b}` },
                { label: 'simplified', value: `${ratio.simplifiedA} : ${ratio.simplifiedB}` },
                { label: 'as a fraction', value: `${ratio.a}/${ratio.b}` },
              ]
            : [
                { label: 'equation', value: `x + ${fmt(constant)} = ${fmt(total)}` },
                { label: 'guess x', value: fmt(guess) },
                { label: 'left side', value: fmt(balance.left) },
                { label: 'right side', value: fmt(balance.right) },
                { label: 'balanced', value: balance.balanced ? 'yes' : 'no' },
                { label: 'solution', value: fmt(balance.solution) },
              ]

  const signedScale = scaleLinear(-10, 10, PAD, W - PAD)
  const signedY = 155
  const signedTicks = niceTicks(-10, 10, 5)

  return (
    <FigureFrame
      title={title ?? MODE_TITLES[mode]}
      hint="Choose a model, then drag a slider. Every marker, bar and readout comes from the same value."
      caption={
        caption ??
        'A number line, a bar, a tape and a balance are four ways to see one quantity. The position or length is the value, and changing it changes every label at once.'
      }
      onReset={reset}
      readout={readout}
    >
      <svg
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        role="img"
        aria-label={`${MODE_LABELS[mode]} with ${mode === 'signed' ? `${fmt(a)} and ${fmt(b)}` : mode === 'decimal' ? `${tenths} tenths` : mode === 'percent' ? `${fmt(percent)} percent of ${fmt(base)}` : mode === 'ratio' ? `ratio ${ratioA} to ${ratioB}` : `x plus ${fmt(constant)} equals ${fmt(total)}`}`}
      >
        <defs>
          <clipPath id="number-line-plot">
            <rect x={PAD} y={PAD} width={W - PAD * 2} height={H - PAD * 2} />
          </clipPath>
        </defs>

        {mode === 'signed' ? (
          <>
            <line x1={PAD} y1={signedY} x2={W - PAD} y2={signedY} stroke="var(--text-muted)" strokeWidth="1" />
            <XAxisTicks ticks={signedTicks} scale={signedScale} at={signedY} />
            <circle cx={signedScale(a)} cy={signedY} r="5" fill="var(--accent-blue)" />
            <text x={signedScale(a)} y={signedY - 12} textAnchor="middle" fontSize="11" fill="var(--accent-blue)">
              A {fmt(a)}
            </text>
            <circle cx={signedScale(b)} cy={signedY} r="5" fill="var(--chart-2)" />
            <text x={signedScale(b)} y={signedY + 22} textAnchor="middle" fontSize="11" fill="var(--chart-2)">
              B {fmt(b)}
            </text>
            <text x={W / 2} y={70} textAnchor="middle" fontSize="11" fill="var(--text-primary)">
              {signed.sentence}
            </text>
          </>
        ) : null}

        {mode === 'decimal' ? (
          <>
            <line x1={PAD} y1={170} x2={W - PAD} y2={170} stroke="var(--text-muted)" strokeWidth="1" />
            {Array.from({ length: 11 }, (_, index) => {
              const x = PAD + (index / 10) * (W - PAD * 2)
              return (
                <g key={index}>
                  <line x1={x} y1={164} x2={x} y2={176} stroke="var(--text-muted)" strokeWidth="1" />
                  {index % 5 === 0 ? (
                    <text x={x} y={192} textAnchor="middle" fontSize="9" fill="var(--text-muted)">
                      {index === 10 ? '1' : `0.${index}`}
                    </text>
                  ) : null}
                </g>
              )
            })}
            <circle cx={PAD + (tenths / 10) * (W - PAD * 2)} cy={170} r="6" fill="var(--accent-blue)" />
            <text x={W / 2} y={90} textAnchor="middle" fontSize="22" fill="var(--text-primary)">
              {fmt(decimal, 1)}
            </text>
            <text x={W / 2} y={116} textAnchor="middle" fontSize="11" fill="var(--text-muted)">
              {tenths} of 10 equal parts
            </text>
          </>
        ) : null}

        {mode === 'percent' ? (
          <>
            <rect x={PAD} y={120} width={W - PAD * 2} height={40} fill="var(--bg-subtle)" stroke="var(--text-muted)" />
            <rect
              x={PAD}
              y={120}
              width={((W - PAD * 2) * Math.max(0, Math.min(100, percent))) / 100}
              height={40}
              fill="var(--accent-blue)"
              opacity={0.35}
            />
            <line
              x1={PAD + ((W - PAD * 2) * Math.max(0, Math.min(100, percent))) / 100}
              y1={114}
              x2={PAD + ((W - PAD * 2) * Math.max(0, Math.min(100, percent))) / 100}
              y2={166}
              stroke="var(--accent-blue)"
              strokeWidth="2"
            />
            <text x={W / 2} y={90} textAnchor="middle" fontSize="13" fill="var(--text-primary)">
              {fmt(percent)}% of {fmt(base)}
            </text>
            <text x={W / 2} y={200} textAnchor="middle" fontSize="12" fill="var(--text-primary)">
              change {fmt(change.amount)} → {fmt(change.result)}
            </text>
            <text x={PAD} y={220} fontSize="9" fill="var(--text-muted)">
              0%
            </text>
            <text x={W - PAD} y={220} textAnchor="end" fontSize="9" fill="var(--text-muted)">
              100%
            </text>
          </>
        ) : null}

        {mode === 'ratio' ? (
          <>
            {[0, 1].map((row) => {
              const count = row === 0 ? ratio.a : ratio.b
              const width = (W - PAD * 2) / Math.max(ratio.a, ratio.b)
              const y = 110 + row * 70
              return (
                <g key={row}>
                  {Array.from({ length: count }, (_, index) => (
                    <rect
                      key={index}
                      x={PAD + index * width}
                      y={y}
                      width={Math.max(1, width - 2)}
                      height={40}
                      fill={row === 0 ? 'var(--accent-blue)' : 'var(--chart-2)'}
                      opacity={0.3}
                      stroke={row === 0 ? 'var(--accent-blue)' : 'var(--chart-2)'}
                    />
                  ))}
                  <text x={PAD} y={y - 8} fontSize="10" fill="var(--text-primary)">
                    {row === 0 ? `A = ${ratio.a}` : `B = ${ratio.b}`}
                  </text>
                </g>
              )
            })}
            <text x={W / 2} y={230} textAnchor="middle" fontSize="12" fill="var(--text-primary)">
              {ratio.sentence}
            </text>
          </>
        ) : null}

        {mode === 'balance' ? (
          <>
            <line
              x1={W / 2}
              y1={120}
              x2={W / 2}
              y2={190}
              stroke="var(--text-muted)"
              strokeWidth="2"
            />
            <line
              x1={70}
              y1={120}
              x2={W - 70}
              y2={120}
              stroke="var(--text-primary)"
              strokeWidth="3"
              transform={`rotate(${balance.balanced ? 0 : balance.left > balance.right ? 4 : -4} ${W / 2} 120)`}
            />
            <g transform={`translate(${70} 145)`}>
              <rect x={-32} y={0} width={64} height={34} fill="var(--bg-subtle)" stroke="var(--accent-blue)" />
              <text x={0} y={21} textAnchor="middle" fontSize="11" fill="var(--accent-blue)">
                x + {fmt(constant)}
              </text>
              <text x={0} y={52} textAnchor="middle" fontSize="10" fill="var(--text-primary)">
                {fmt(balance.left)}
              </text>
            </g>
            <g transform={`translate(${W - 70} 145)`}>
              <rect x={-32} y={0} width={64} height={34} fill="var(--bg-subtle)" stroke="var(--chart-2)" />
              <text x={0} y={21} textAnchor="middle" fontSize="11" fill="var(--chart-2)">
                {fmt(total)}
              </text>
              <text x={0} y={52} textAnchor="middle" fontSize="10" fill="var(--text-primary)">
                {fmt(balance.right)}
              </text>
            </g>
            <text x={W / 2} y={240} textAnchor="middle" fontSize="11" fill={balance.balanced ? 'var(--accent-green)' : 'var(--accent-red)'}>
              {balance.balanced ? 'balanced' : `x must be ${fmt(balance.solution)}`}
            </text>
          </>
        ) : null}
      </svg>

      <div className="mt-3 flex items-center gap-2">
        <label htmlFor="number-line-mode" className="font-mono text-[10px] text-mathua-muted">
          Model
        </label>
        <select
          id="number-line-mode"
          aria-label="Model"
          value={mode}
          onChange={(event) => setMode(event.target.value as ArithmeticMode)}
          className="h-7 border border-mathua-border bg-transparent px-1 font-mono text-[10px] text-mathua-secondary"
        >
          {(['signed', 'decimal', 'percent', 'ratio', 'balance'] as const).map((item) => (
            <option key={item} value={item}>
              {MODE_LABELS[item]}
            </option>
          ))}
        </select>
      </div>

      <div className="mt-2 grid grid-cols-1 gap-2">
        {mode === 'signed' ? (
          <>
            <Slider label="Number A" value={a} min={-10} max={10} step={1} onChange={setA} />
            <Slider label="Number B" value={b} min={-10} max={10} step={1} onChange={setB} />
          </>
        ) : null}
        {mode === 'decimal' ? (
          <Slider label="Tenths" value={tenths} min={0} max={10} step={1} onChange={setTenths} />
        ) : null}
        {mode === 'percent' ? (
          <>
            <Slider label="Percent" value={percent} min={-100} max={100} step={1} onChange={setPercent} />
            <Slider label="Base" value={base} min={0} max={200} step={1} onChange={setBase} />
          </>
        ) : null}
        {mode === 'ratio' ? (
          <>
            <Slider label="First term" value={ratioA} min={1} max={12} step={1} onChange={setRatioA} />
            <Slider label="Second term" value={ratioB} min={1} max={12} step={1} onChange={setRatioB} />
          </>
        ) : null}
        {mode === 'balance' ? (
          <>
            <Slider label="Guess x" value={guess} min={-10} max={10} step={1} onChange={setGuess} />
            <Slider label="Constant b" value={constant} min={-10} max={10} step={1} onChange={setConstant} />
            <Slider label="Right side c" value={total} min={-10} max={10} step={1} onChange={setTotal} />
          </>
        ) : null}
      </div>
    </FigureFrame>
  )
}
