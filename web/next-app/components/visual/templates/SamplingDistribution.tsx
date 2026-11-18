'use client'

import { useMemo, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Select, Slider } from '../Controls'
import { XAxisTicks } from '../Axis'
import { samplingSchema } from '../schemas'
import { histogram, mean, populationPdf, sampleMeans, sd, type Population } from '../math/sampling'
import { normalPdf } from '../math/normal'
import { pathFromPoints, scaleLinear } from '../scale'
import { fmt } from '../format'
import type { TemplateProps } from '../types'

const W = 340
const H = 250
const PAD = 28
const LO = -3.4
const HI = 3.4
const BINS = 34

const POP_OPTIONS: Array<{ value: Population; label: string }> = [
  { value: 'skewed', label: 'Skewed' },
  { value: 'uniform', label: 'Uniform' },
  { value: 'normal', label: 'Normal' },
]

export default function SamplingDistribution({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => samplingSchema.parse(params), [params])
  const [population, setPopulation] = useState<Population>(defaults.population)
  const [n, setN] = useState(defaults.n)
  const [draws, setDraws] = useState(defaults.draws)

  const values = useMemo(() => sampleMeans(population, n, draws, defaults.seed), [population, n, draws, defaults.seed])
  const counts = useMemo(() => histogram(values, LO, HI, BINS), [values])
  const maxCount = Math.max(1, ...counts)
  const standardError = 1 / Math.sqrt(n)
  const binWidth = (HI - LO) / BINS

  const X = scaleLinear(LO, HI, PAD, W - PAD)
  const Y = scaleLinear(0, maxCount * 1.12, H - PAD, PAD)

  const densityCurve = useMemo(() => {
    const out: Array<[number, number]> = []
    for (let i = 0; i <= 160; i++) {
      const x = LO + ((HI - LO) * i) / 160
      out.push([x, populationPdf(population, x)])
    }
    return out
  }, [population])
  const maxDensity = Math.max(...densityCurve.map(([, y]) => y))

  const expectedCurve = useMemo(() => {
    const out: Array<[number, number]> = []
    for (let i = 0; i <= 160; i++) {
      const x = LO + ((HI - LO) * i) / 160
      out.push([x, normalPdf(x, 0, standardError) * draws * binWidth])
    }
    return out
  }, [standardError, draws, binWidth])

  const reset = () => {
    setPopulation(defaults.population)
    setN(defaults.n)
    setDraws(defaults.draws)
  }

  return (
    <FigureFrame
      title={title ?? 'The sampling distribution of the mean'}
      hint="Change n: the bars narrow even though the population does not."
      caption={
        caption ??
        'Each bar counts the mean of one sample of size n. The grey curve is the population, drawn to its own scale; the dashed curve is the normal the central limit theorem predicts, with spread 1/√n.'
      }
      onReset={reset}
      readout={[
        { label: 'n', value: String(n) },
        { label: 'standard error', value: fmt(standardError, 3) },
        { label: 'observed sd', value: fmt(sd(values), 3) },
        { label: 'mean of means', value: fmt(mean(values), 3) },
        { label: 'samples', value: String(draws) },
      ]}
    >
      <svg
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        role="img"
        aria-label="A histogram of sample means beside the population density"
      >
        <line x1={PAD} y1={Y(0)} x2={W - PAD} y2={Y(0)} stroke="var(--text-muted)" strokeWidth="0.8" />

        {counts.map((count, i) => {
          const x0 = LO + i * binWidth
          return count === 0 ? null : (
            <rect
              key={i}
              x={X(x0)}
              y={Y(count)}
              width={Math.max(0, X(x0 + binWidth) - X(x0) - 0.6)}
              height={Math.max(0, Y(0) - Y(count))}
              fill="var(--accent-blue)"
              fillOpacity="0.18"
              stroke="var(--accent-blue)"
              strokeWidth="0.7"
            />
          )
        })}

        {/* The population, rescaled to the histogram's height so the two spreads can
            be compared on one axis. It is a shape reference, not a probability scale. */}
        <path
          d={pathFromPoints(densityCurve.map(([x, y]) => [X(x), Y((y / maxDensity) * maxCount * 0.55)]))}
          fill="none"
          stroke="var(--text-muted)"
          strokeWidth="1.2"
          strokeDasharray="2 3"
        />
        <path
          d={pathFromPoints(expectedCurve.map(([x, y]) => [X(x), Y(y)]))}
          fill="none"
          stroke="var(--accent-teal)"
          strokeWidth="1.5"
          strokeDasharray="6 3"
        />

        <XAxisTicks ticks={[-3, -2, -1, 0, 1, 2, 3]} scale={X} at={Y(0)} />
        <text x={PAD + 2} y={PAD + 10} fontSize="9" fill="var(--text-muted)">
          population (shape only)
        </text>
        <text x={W - PAD} y={PAD + 10} textAnchor="end" fontSize="9" fill="var(--accent-teal)">
          1/√n
        </text>
      </svg>

      <div className="mt-3 space-y-2">
        <Select label="Population" value={population} options={POP_OPTIONS} onChange={(v) => setPopulation(v as Population)} />
        <Slider label="n" value={n} min={1} max={64} step={1} onChange={(v) => setN(Math.round(v))} />
        <Slider label="draws" value={draws} min={100} max={2000} step={100} onChange={(v) => setDraws(Math.round(v))} />
      </div>
    </FigureFrame>
  )
}
