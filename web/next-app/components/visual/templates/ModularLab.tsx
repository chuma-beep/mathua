'use client'

import { useMemo, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Slider } from '../Controls'
import { modularSchema } from '../schemas'
import { congruence, euclid, factorize, modPow } from '../math/modular'
import { fmt } from '../format'
import type { TemplateProps } from '../types'

const W = 360
const H = 300
const PAD = 30

export default function ModularLab({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => modularSchema.parse(params), [params])
  const [mode, setMode] = useState(defaults.mode)
  const [a, setA] = useState(defaults.a)
  const [n, setN] = useState(defaults.n)
  const [base, setBase] = useState(defaults.base)
  const [exponent, setExponent] = useState(defaults.exponent)
  const [b, setB] = useState(defaults.b)

  const reset = () => {
    setMode(defaults.mode)
    setA(defaults.a)
    setN(defaults.n)
    setBase(defaults.base)
    setExponent(defaults.exponent)
    setB(defaults.b)
  }

  const cong = congruence(a, n)
  const power = modPow(base, exponent, n)
  const gcd = euclid(a, b)
  const fact = factorize(a)

  const circle = Array.from({ length: Math.max(1, n) }, (_v, index) => {
    const angle = (index / Math.max(1, n)) * Math.PI * 2 - Math.PI / 2
    return { index, x: 180 + 78 * Math.cos(angle), y: 128 + 78 * Math.sin(angle) }
  })
  const selected = cong?.residue ?? -1

  return (
    <FigureFrame
      title={title ?? 'Exact integer arithmetic'}
      hint="Choose a model, then change the integers. Every residue, quotient and remainder is computed, not sketched."
      caption={caption ?? 'Modular arithmetic keeps the remainder. Repeated squaring reduces the base at every step. Euclid replaces the pair by the divisor and the remainder, and prime factorization stops when the quotient is prime.'}
      onReset={reset}
      readout={[
        { label: 'model', value: mode },
        { label: 'a', value: String(a) },
        { label: 'modulus n', value: String(n) },
        {
          label: 'a = qn + r',
          value: cong ? `${a} = ${cong.quotient}·${n} + ${cong.residue}` : '—',
        },
        {
          label: 'remainder',
          value: cong ? String(cong.residue) : '—',
        },
        {
          label: 'modular power',
          value: power ? `${base}^${exponent} ≡ ${power.value} (mod ${n})` : '—',
        },
        {
          label: 'gcd',
          value: gcd ? `gcd(${Math.abs(a)}, ${Math.abs(b)}) = ${gcd.gcd}` : '—',
        },
        {
          label: 'factorization',
          value: fact ? `${a} = ${fact.factors.join(' × ')}` : '—',
        },
      ]}
    >
      <svg
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        role="img"
        aria-label={`${mode} model for ${a} modulo ${n}`}
      >
        {(mode === 'residue' || mode === 'power') && (
          <>
            <circle cx="180" cy="128" r="78" fill="none" stroke="var(--text-muted)" strokeWidth="1" />
            {circle.map((point) => (
              <g key={point.index}>
                <circle
                  cx={point.x}
                  cy={point.y}
                  r={point.index === selected ? 5 : 2.5}
                  fill={point.index === selected ? 'var(--accent-blue)' : 'var(--text-muted)'}
                />
                <text x={point.x} y={point.y - 7} textAnchor="middle" fontSize="9" fill="var(--text-secondary)">
                  {point.index}
                </text>
              </g>
            ))}
            <text x="90" y="238" fontSize="10" fill="var(--text-secondary)">
              residues mod {n}
            </text>
            <text x="90" y="254" fontSize="10" fill="var(--accent-blue)">
              {a} lands on {selected}
            </text>
          </>
        )}

        {(mode === 'euclid' || mode === 'factor') && (
          <>
            <text x={PAD} y={PAD + 8} fontSize="10" fill="var(--accent-blue)">
              {mode === 'euclid' ? `gcd(${Math.abs(a)}, ${Math.abs(b)})` : `factor ${a}`}
            </text>
            {mode === 'euclid'
              ? (gcd?.steps ?? []).slice(0, 11).map((step, index) => (
                  <text key={index} x={PAD} y={PAD + 28 + index * 18} fontSize="10" fill="var(--text-secondary)">
                    {step.a} = {step.quotient}·{step.b} + {step.remainder}
                  </text>
                ))
              : (fact?.tree ?? []).slice(0, 11).map((line, index) => (
                  <text key={index} x={PAD} y={PAD + 28 + index * 18} fontSize="10" fill="var(--text-secondary)">
                    {line}
                  </text>
                ))}
            {mode === 'euclid' && gcd ? (
              <text x={PAD} y={PAD + 28 + gcd.steps.length * 18} fontSize="10" fill="var(--accent-green)">
                gcd = {gcd.gcd}
              </text>
            ) : null}
            {mode === 'factor' && fact ? (
              <text x={PAD} y={PAD + 28 + fact.tree.length * 18} fontSize="10" fill="var(--accent-green)">
                {fact.factors.join(' × ')}
              </text>
            ) : null}
          </>
        )}

        {mode === 'power' && power && (
          <>
            {power.steps.slice(0, 8).map((step, index) => (
              <text key={index} x={PAD} y={20 + index * 18} fontSize="10" fill="var(--text-secondary)">
                e {step.exponent}: base {step.base}, result {step.result}
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
            <option value="residue">Residue</option>
            <option value="power">Repeated squaring</option>
            <option value="euclid">Euclid</option>
            <option value="factor">Factorization</option>
          </select>
        </label>
        <div className="grid grid-cols-1 gap-2">
          <Slider label="a" value={a} min={-50} max={100} step={1} onChange={setA} />
          <Slider label="modulus n" value={n} min={2} max={12} step={1} onChange={setN} />
          <Slider label="base" value={base} min={1} max={12} step={1} onChange={setBase} />
          <Slider label="exponent" value={exponent} min={1} max={30} step={1} onChange={setExponent} />
          <Slider label="b" value={b} min={1} max={100} step={1} onChange={setB} />
        </div>
      </div>
    </FigureFrame>
  )
}
