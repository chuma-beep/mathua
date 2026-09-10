import type { CSSProperties } from 'react'

export const headingFont = "'IBM Plex Serif', serif"
export const bodyFont = "'IBM Plex Serif', serif"
export const monoFont = "'IBM Plex Mono', monospace"

export const sections = [
  { id: 'concept-graph', label: 'Concept graph' },
  { id: 'student-model', label: 'Student model' },
  { id: 'spaced-repetition', label: 'Spaced repetition' },
  { id: 'diagnostic', label: 'Diagnostic' },
  { id: 'scheduler', label: 'Scheduler' },
  { id: 'scoring', label: 'Scoring' },
  { id: 'generators', label: 'Generators' },
  { id: 'symbolic-grading', label: 'Expression grading' },
]

export const inlineCodeStyle: CSSProperties = {
  fontFamily: monoFont,
  fontSize: '0.9em',
  color: 'var(--accent-blue)',
  overflowWrap: 'anywhere',
  wordBreak: 'break-all',
}

export const mutedCodeStyle: CSSProperties = {
  fontFamily: monoFont,
  fontSize: '0.9em',
  color: 'var(--text-secondary)',
  overflowWrap: 'anywhere',
  wordBreak: 'break-all',
}

export const h2Style: CSSProperties = {
  fontFamily: headingFont,
  fontWeight: 400,
  fontSize: 'clamp(1.4rem, 4vw, 1.7rem)',
  color: 'var(--text-primary)',
  borderBottom: '0.5px solid var(--border)',
  paddingBottom: '0.5rem',
  marginBottom: '1.5rem',
  letterSpacing: '-0.01em',
}

export const bodyStyle: CSSProperties = {
  fontFamily: bodyFont,
  fontSize: 'clamp(0.9rem, 0.85rem + 0.5vw, 1rem)',
  color: 'var(--text-secondary)',
  lineHeight: 1.85,
  marginBottom: '1rem',
}

export const codeBlockStyle: CSSProperties = {
  background: 'transparent',
  border: 'none',
  borderLeft: '2px solid var(--accent-blue)',
  borderRadius: 0,
  padding: '0.5rem 0 0.5rem 1.5rem',
  fontFamily: monoFont,
  fontSize: '13px',
  color: 'var(--text-secondary)',
  whiteSpace: 'pre',
  overflowX: 'auto',
  lineHeight: 1.6,
  margin: '1rem 0',
  width: '100%',
  maxWidth: '100%',
  minWidth: 0,
}

export const tableHeaderStyle: CSSProperties = {
  fontFamily: headingFont,
  fontWeight: 400,
  fontSize: 'clamp(0.8rem, 0.75rem + 0.3vw, 0.9rem)',
  textTransform: 'uppercase',
  letterSpacing: '0.1em',
  color: 'var(--accent-blue)',
  padding: '10px 14px 10px 12px',
  borderBottom: '1px solid var(--accent-blue)',
  textAlign: 'left',
  background: 'transparent',
}

export const tableCellStyle: CSSProperties = {
  fontFamily: bodyFont,
  fontSize: 'clamp(0.85rem, 0.8rem + 0.4vw, 1rem)',
  color: 'var(--text-secondary)',
  padding: '10px 14px 10px 12px',
  borderBottom: '0.5px solid var(--border)',
  background: 'transparent',
}

export const h1Style: CSSProperties = {
  fontFamily: headingFont,
  fontWeight: 400,
  fontSize: '1.9rem',
  color: 'var(--text-primary)',
  marginBottom: '1rem',
  letterSpacing: '-0.01em',
  borderBottom: '0.5px solid var(--border)',
  paddingBottom: '0.5rem',
}

export const stepLiStyle: CSSProperties = {
  fontFamily: bodyFont,
  fontSize: '0.95rem',
  color: 'var(--text-secondary)',
  lineHeight: 1.7,
  marginBottom: '0.75rem',
  paddingLeft: '2.25rem',
  position: 'relative',
}

export function roman(n: number): string {
  const r = ['I', 'II', 'III', 'IV', 'V', 'VI', 'VII', 'VIII', 'IX', 'X']
  return r[n - 1] ?? String(n)
}
