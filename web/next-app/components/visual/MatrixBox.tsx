'use client'

import { fmt } from './format'
import type { Mat2 } from './math/matrix'

/**
 * A bracketed 2x2 matrix, display only.
 *
 * The columns are tinted to match the arrows that draw them, and each pair of
 * entries carries the name of its column. The shape matters: the bracket is what
 * says "this is one object", which a row of four number fields does not, and the
 * column footers are what tie each number to the arrow a learner can move.
 *
 * It is deliberately not an editor. Typing happens in the four number fields
 * below it and by dragging the arrowheads; this component only states what the
 * figure is currently drawing, and carries the whole matrix in its aria-label
 * because the two-letter footer glyphs are not a sentence.
 */
export default function MatrixBox({ matrix, symbol = 'A' }: { matrix: Mat2; symbol?: string }) {
  const rows = [
    [matrix[0][0], matrix[0][1]],
    [matrix[1][0], matrix[1][1]],
  ]
  const label = `${symbol} equals the matrix [[${rows[0].map((value) => fmt(value)).join(', ')}], [${rows[1]
    .map((value) => fmt(value))
    .join(', ')}]]`

  return (
    <div className="flex items-center justify-center gap-2">
      <span className="font-mono text-[12px] text-mathua-secondary">{symbol} =</span>
      <svg
        viewBox="0 0 132 82"
        style={{ width: 132, height: 82, flexShrink: 0 }}
        role="img"
        aria-label={label}
      >
        <path d="M4 4 L4 58 M4 4 L12 4 M4 58 L12 58" fill="none" stroke="var(--text-secondary)" strokeWidth="1.4" />
        <path d="M128 4 L128 58 M120 4 L128 4 M120 58 L128 58" fill="none" stroke="var(--text-secondary)" strokeWidth="1.4" />
        {rows.map((row, r) =>
          row.map((value, c) => (
            <text
              key={`${r}-${c}`}
              x={c === 0 ? 44 : 88}
              y={r === 0 ? 26 : 46}
              textAnchor="middle"
              fontSize="14"
              fill={c === 0 ? 'var(--accent-blue)' : 'var(--chart-2)'}
              fontFamily="inherit"
            >
              {fmt(value)}
            </text>
          )),
        )}
        <text x={44} y={72} textAnchor="middle" fontSize="8" fill="var(--accent-blue)">
          {symbol} e₁
        </text>
        <text x={88} y={72} textAnchor="middle" fontSize="8" fill="var(--chart-2)">
          {symbol} e₂
        </text>
      </svg>
    </div>
  )
}
