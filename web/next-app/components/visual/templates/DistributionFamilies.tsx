'use client'

import { useCallback, useMemo, useRef, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Select, Slider } from '../Controls'
import { XAxisTicks, YAxisTicks } from '../Axis'
import { familySchema } from '../schemas'
import { FAMILIES, FAMILY_LABELS, familyCdf, familyMean, familyPdf, familySd, familyTail, type Family } from '../math/families'
import { clamp, invertLinear, pathFromPoints, scaleLinear } from '../scale'
import { fmt } from '../format'
import { useDrag } from '../useDrag'
import type { TemplateProps } from '../types'

const W = 340
const H = 250
const PAD = 28
const X_MIN = -4.5
const X_MAX = 5.5

const COLORS: Record<Family, string> = {
  normal: 'var(--accent-blue)',
  uniform: 'var(--accent-teal)',
  exponential: 'var(--accent-green)',
}

const OPTIONS = FAMILIES.map((value) => ({ value, label: FAMILY_LABELS[value] }))

export default function DistributionFamilies({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => familySchema.parse(params), [params])
  const [emphasis, setEmphasis] = useState<Family>(defaults.emphasis)
  const [x, setX] = useState(defaults.x)
  const svgRef = useRef<SVGSVGElement | null>(null)

  const X = scaleLinear(X_MIN, X_MAX, PAD, W - PAD)
  const Y = scaleLinear(0, 1.15, H - PAD, PAD)
  const invertX = invertLinear(X_MIN, X_MAX, PAD, W - PAD)

  // The emphasised family draws its tail; the others stay visible but dimmed, so
  // the comparison is on one axis rather than across three separate plots.
  const curves = useMemo(
    () =>
      FAMILIES.map((family) => {
        const points: Array<[number, number]> = []
        for (let i = 0; i <= 260; i++) {
          const xv = X_MIN + ((X_MAX - X_MIN) * i) / 260
          points.push([xv, familyPdf(family, xv)])
        }
        return { family, points }
      }),
    [],
  )

  const tail = useMemo(() => {
    const points: Array<[number, number]> = []
    for (let i = 0; i <= 260; i++) {
      const xv = X_MIN + ((clamp(x, X_MIN, X_MAX) - X_MIN) * i) / 260
      points.push([xv, familyPdf(emphasis, xv)])
    }
    return points
  }, [emphasis, x])

  const onDrag = useCallback(
    (point: { x: number; y: number }) => {
      setX(clamp(invertX(point.x), X_MIN, X_MAX))
    },
    [invertX],
  )
  useDrag(svgRef, onDrag)

  const reset = () => {
    setEmphasis(defaults.emphasis)
    setX(defaults.x)
  }

  return (
    <FigureFrame
      title={title ?? 'One spread, three shapes'}
      hint="Drag the line or move x. Every family here has standard deviation 1."
      caption={
        caption ??
        'Three densities with the same spread and different shapes. The shaded tail is the probability the emphasised family puts below x, and it is not the same number for the others: a mean and a standard deviation do not name a distribution. The exponential starts at zero, so its mean is 1 rather than 0, which the mean row reports for whichever family is emphasised.'
      }
      onReset={reset}
      readout={[
        { label: 'x', value: fmt(x) },
        { label: 'f(x)', value: fmt(familyPdf(emphasis, x), 3) },
        { label: 'P(X ≤ x)', value: fmt(familyCdf(emphasis, x), 4) },
        { label: 'P normal', value: fmt(familyCdf('normal', x), 4) },
        { label: 'P uniform', value: fmt(familyCdf('uniform', x), 4) },
        { label: 'P exp', value: fmt(familyCdf('exponential', x), 4) },
        { label: 'mean', value: fmt(familyMean(emphasis), 2) },
        { label: 'sd', value: fmt(familySd(emphasis), 2) },
        // Named, not "tail": the shaded region is the left tail, so an unqualified
        // "tail" beside it read as the shaded area while printing its complement.
        { label: 'P(X > x)', value: fmt(familyTail(emphasis, x), 4) },
      ]}
    >
      <svg
        ref={svgRef}
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] cursor-grab font-mono active:cursor-grabbing"
        style={{ touchAction: 'pan-y' }}
        role="img"
        aria-label="Three continuous densities drawn on one axis with a shaded tail"
      >
        <line x1={PAD} y1={Y(0)} x2={W - PAD} y2={Y(0)} stroke="var(--text-muted)" strokeWidth="0.8" />

        {tail.length > 1 && (
          <path
            d={`${pathFromPoints(tail.map(([xv, yv]) => [X(xv), Y(yv)]))} L ${X(clamp(x, X_MIN, X_MAX))} ${Y(0)} L ${X(X_MIN)} ${Y(0)} Z`}
            fill={COLORS[emphasis]}
            fillOpacity="0.16"
          />
        )}

        {curves.map(({ family, points }) => (
          <path
            key={family}
            d={pathFromPoints(points.map(([xv, yv]) => [X(xv), Y(yv)]))}
            fill="none"
            stroke={COLORS[family]}
            strokeWidth={family === emphasis ? 1.9 : 1.1}
            opacity={family === emphasis ? 1 : 0.45}
          />
        ))}

        <line x1={X(x)} y1={PAD} x2={X(x)} y2={H - PAD} stroke="var(--text-muted)" strokeWidth="1" strokeDasharray="3 3" />
        {FAMILIES.map((family) => (
          <circle key={`d-${family}`} cx={X(x)} cy={Y(familyPdf(family, x))} r={family === emphasis ? 4 : 2.6} fill={COLORS[family]} opacity={family === emphasis ? 1 : 0.5} />
        ))}

        <XAxisTicks ticks={[-4, -2, 0, 2, 4]} scale={X} at={Y(0)} />
        {/* Density is a height, so the axis needs a number: without one the
            exponential's peak at zero looks like decoration rather than 1. */}
        <YAxisTicks ticks={[0.5, 1]} scale={Y} at={PAD} format={(v) => fmt(v, 1)} />
        {/* A legend row above the plot, because floating labels beside each curve
            landed on the curves and the vertical line. */}
        {FAMILIES.map((family, i) => (
          <g key={`l-${family}`} opacity={family === emphasis ? 1 : 0.6}>
            <line x1={PAD + i * 82} y1={PAD - 12} x2={PAD + 10 + i * 82} y2={PAD - 12} stroke={COLORS[family]} strokeWidth="2" />
            <text x={PAD + 14 + i * 82} y={PAD - 9} fontSize="9" fill={COLORS[family]}>
              {family === 'exponential' ? 'exp' : family}
            </text>
          </g>
        ))}
      </svg>

      <div className="mt-3 space-y-2">
        <Select label="Emphasise" value={emphasis} options={OPTIONS} onChange={(v) => setEmphasis(v as Family)} />
        <Slider label="x" value={x} min={X_MIN} max={X_MAX} step={0.01} onChange={setX} />
      </div>
    </FigureFrame>
  )
}
