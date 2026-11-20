'use client'

import { useId, useMemo, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Slider, Toggle } from '../Controls'
import { XAxisTicks, YAxisTicks } from '../Axis'
import { sinusoidSchema } from '../schemas'
import {
  formatPi,
  sinusoid,
  sinusoidFormula,
  sinusoidValue,
  type SinusoidForm,
} from '../math/sinusoid'
import { TAU } from '../math/trig'
import { niceTicks, scaleLinear } from '../scale'
import { fmt } from '../format'
import type { TemplateProps } from '../types'

const W = 360
const H = 320
const PAD = 32
const X_MIN = -TAU
const X_MAX = TAU

const FORM_LABELS: Record<SinusoidForm, string> = {
  sin: 'sine',
  cos: 'cosine',
}

export default function SinusoidExplorer({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => sinusoidSchema.parse(params), [params])
  const [form, setForm] = useState<SinusoidForm>(defaults.form)
  const [amplitude, setAmplitude] = useState(defaults.amplitude)
  const [frequency, setFrequency] = useState(defaults.frequency)
  const [phase, setPhase] = useState(defaults.phase)
  const [vertical, setVertical] = useState(defaults.vertical)
  const [showReference, setShowReference] = useState(defaults.showReference)
  const clipId = useId()

  const wave = useMemo(
    () => sinusoid({ form, amplitude, frequency, phase, vertical }),
    [form, amplitude, frequency, phase, vertical],
  )
  const equation = sinusoidFormula(wave)

  const yTop = Math.max(2, Math.ceil(wave.max + 1))
  const yBottom = Math.min(-2, Math.floor(wave.min - 1))
  const X = scaleLinear(X_MIN, X_MAX, PAD, W - PAD)
  const Y = scaleLinear(yTop, yBottom, PAD, H - PAD)
  const majorX = [-TAU, -Math.PI, 0, Math.PI, TAU]
  const minorX = [-(3 * Math.PI) / 2, -Math.PI / 2, Math.PI / 2, (3 * Math.PI) / 2]
  const ticksY = niceTicks(yBottom, yTop, 4)

  const curve = useMemo(() => {
    let d = ''
    let drawing = false
    for (let index = 0; index <= 360; index++) {
      const x = X_MIN + ((X_MAX - X_MIN) * index) / 360
      const y = sinusoidValue(wave, x)
      if (!Number.isFinite(y) || y < yBottom || y > yTop) {
        drawing = false
        continue
      }
      d += `${drawing ? 'L' : 'M'}${X(x).toFixed(2)},${Y(y).toFixed(2)} `
      drawing = true
    }
    return d.trim()
  }, [wave, X, Y, yBottom, yTop])

  const reference = useMemo(() => {
    let d = ''
    for (let index = 0; index <= 360; index++) {
      const x = X_MIN + ((X_MAX - X_MIN) * index) / 360
      d += `${index === 0 ? 'M' : 'L'}${X(x).toFixed(2)},${Y(Math.sin(x)).toFixed(2)} `
    }
    return d.trim()
  }, [X, Y])

  const reset = () => {
    setForm(defaults.form)
    setAmplitude(defaults.amplitude)
    setFrequency(defaults.frequency)
    setPhase(defaults.phase)
    setVertical(defaults.vertical)
    setShowReference(defaults.showReference)
  }

  const phaseDirection = phase > 1e-10 ? 'right' : phase < -1e-10 ? 'left' : 'none'
  const phaseValue = phaseDirection === 'none' ? '0' : `${formatPi(Math.abs(phase))} ${phaseDirection}`
  const bandTop = Y(Math.min(yTop, wave.max))
  const bandBottom = Y(Math.max(yBottom, wave.min))

  return (
    <FigureFrame
      title={title ?? 'One wave, four parameters'}
      hint="Change A, B, C and D. A sets the height, B the period, C the horizontal shift and D the midline."
      caption={
        caption ??
        'For y = A f(B(x − C)) + D, the amplitude is |A|, the period is 2π/|B|, C is the phase shift and D is the midline. The dashed curve is y = sin x for comparison.'
      }
      onReset={reset}
      readout={[
        { label: 'form', value: FORM_LABELS[form] },
        { label: 'equation', value: equation },
        { label: 'amplitude', value: fmt(wave.peak) },
        { label: 'period', value: formatPi(wave.period) },
        { label: 'frequency', value: `${fmt(Math.abs(wave.frequency) / TAU, 3)} cycles/unit` },
        { label: 'phase shift', value: phaseValue },
        { label: 'vertical shift', value: fmt(wave.vertical) },
        { label: 'maximum', value: fmt(wave.max) },
        { label: 'minimum', value: fmt(wave.min) },
        { label: 'midline', value: fmt(wave.midline) },
      ]}
    >
      <svg
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        role="img"
        aria-label={`Graph of ${equation}`}
      >
        <defs>
          <clipPath id={clipId}>
            <rect x={PAD} y={PAD} width={W - PAD * 2} height={H - PAD * 2} />
          </clipPath>
        </defs>

        <line x1={PAD} y1={Y(0)} x2={W - PAD} y2={Y(0)} stroke="var(--text-muted)" strokeWidth="0.8" />
        <line x1={X(0)} y1={PAD} x2={X(0)} y2={H - PAD} stroke="var(--text-muted)" strokeWidth="0.8" />

        <g clipPath={`url(#${clipId})`}>
          <rect x={PAD} y={bandTop} width={W - PAD * 2} height={Math.max(0, bandBottom - bandTop)} fill="var(--accent-blue)" opacity="0.06" />

          <line
            x1={PAD}
            y1={Y(wave.midline)}
            x2={W - PAD}
            y2={Y(wave.midline)}
            stroke="var(--text-muted)"
            strokeWidth="0.9"
            strokeDasharray="4 4"
          />

          {phase >= X_MIN && phase <= X_MAX ? (
            <line
              x1={X(phase)}
              y1={PAD + 26}
              x2={X(phase)}
              y2={H - PAD}
              stroke="var(--chart-3)"
              strokeWidth="1"
              strokeDasharray="3 4"
            />
          ) : null}

          {showReference ? (
            <path d={reference} fill="none" stroke="var(--text-muted)" strokeWidth="1" strokeDasharray="5 4" />
          ) : null}

          <path d={curve} fill="none" stroke="var(--accent-blue)" strokeWidth="2.2" strokeLinecap="round" />
        </g>

        {minorX.map((value) => (
          <line key={value} x1={X(value)} y1={Y(0)} x2={X(value)} y2={Y(0) + 3} stroke="var(--text-muted)" strokeWidth="0.7" />
        ))}
        <XAxisTicks ticks={majorX} scale={X} at={Y(0)} format={(value) => formatPi(value)} />
        <YAxisTicks ticks={ticksY} scale={Y} at={X(0)} />

        <text x={PAD + 2} y={PAD + 10} fontSize="10" fill="var(--accent-blue)">
          {equation}
        </text>
        {showReference ? (
          <text x={W - PAD - 2} y={H - PAD - 4} fontSize="9" fill="var(--text-muted)" textAnchor="end">
            dashed: sin x
          </text>
        ) : null}
        <text x={X(phase)} y={PAD + 20} fontSize="8" fill="var(--chart-3)" textAnchor="middle">
          {phaseDirection === 'none' ? '' : 'C'}
        </text>
      </svg>

      <div className="mt-3 space-y-2">
        <div className="flex items-center gap-2">
          <label htmlFor={`${clipId}-form`} className="w-14 shrink-0 font-mono text-[10px] text-mathua-muted">
            Form
          </label>
          <select
            id={`${clipId}-form`}
            aria-label="Form"
            value={form}
            onChange={(event) => setForm(event.target.value as SinusoidForm)}
            className="h-7 min-w-0 border border-mathua-border bg-transparent px-1 font-mono text-[10px] text-mathua-secondary"
          >
            <option value="sin">Sine</option>
            <option value="cos">Cosine</option>
          </select>
        </div>
        <Slider label="Amplitude A" value={amplitude} min={-3} max={3} step={0.1} onChange={setAmplitude} />
        <Slider label="Coefficient B" value={frequency} min={0.25} max={4} step={0.05} onChange={setFrequency} />
        <Slider label="Phase C" value={phase} min={-3.14} max={3.14} step={0.05} onChange={setPhase} />
        <Slider label="Vertical D" value={vertical} min={-8} max={8} step={0.1} onChange={setVertical} />
        <Toggle label="Show reference" checked={showReference} onChange={setShowReference} />
      </div>
    </FigureFrame>
  )
}
