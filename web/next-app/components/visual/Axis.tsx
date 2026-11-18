'use client'

// Numeric scale labels for the plot templates.
//
// A curve with no numbers on its axes is a sketch: the reader cannot tell
// whether the tangent at x = 1 is steep or whether the window is simply zoomed
// out, and "the secant approaches the tangent" becomes a claim about two lines
// rather than a measurable fact. Each plot therefore carries its own ticks,
// drawn faintly so they label the drawing instead of competing with it.

interface AxisProps {
  ticks: number[]
  scale: (value: number) => number
  /** Pixel position of the axis the labels hang from. */
  at: number
  format?: (value: number) => string
}

function label(value: number, format?: (value: number) => string): string {
  return format ? format(value) : String(Number(value.toFixed(4)))
}

export function XAxisTicks({ ticks, scale, at, format }: AxisProps) {
  return (
    <g aria-hidden="true">
      {ticks.map((value) => (
        <g key={value}>
          <line x1={scale(value)} y1={at} x2={scale(value)} y2={at + 3} stroke="var(--text-muted)" strokeWidth="0.7" />
          <text x={scale(value)} y={at + 12} textAnchor="middle" fontSize="8" fill="var(--text-muted)">
            {label(value, format)}
          </text>
        </g>
      ))}
    </g>
  )
}

export function YAxisTicks({ ticks, scale, at, format }: AxisProps) {
  return (
    <g aria-hidden="true">
      {ticks.map((value) => (
        <g key={value}>
          <line x1={at - 3} y1={scale(value)} x2={at} y2={scale(value)} stroke="var(--text-muted)" strokeWidth="0.7" />
          <text x={at - 5} y={scale(value)} dy="0.32em" textAnchor="end" fontSize="8" fill="var(--text-muted)">
            {label(value, format)}
          </text>
        </g>
      ))}
    </g>
  )
}
