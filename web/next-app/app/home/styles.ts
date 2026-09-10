import type { CSSProperties } from 'react'

export const headingFont = "'IBM Plex Serif', serif"
export const bodyFont = "'IBM Plex Serif', serif"
export const monoFont = "'IBM Plex Mono', monospace"

export const loadingGraphStyle: CSSProperties = {
  height: 'clamp(320px, 60dvh, 520px)',
  background: 'transparent',
  borderRadius: 0,
  border: '0.5px solid var(--border)',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  color: 'var(--text-muted)',
  fontFamily: monoFont,
  fontSize: '13px',
  width: '100%',
  maxWidth: '100%',
}

export const ctaPrimaryStyle: CSSProperties = {
  background: 'var(--accent-blue)',
  color: 'var(--bg)',
  fontFamily: monoFont,
  fontSize: '13px',
  padding: '10px 22px',
  border: 'none',
  borderRadius: '2px',
  letterSpacing: '0.04em',
  textTransform: 'none',
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  textDecoration: 'none',
  transition: 'background 0.2s',
}

export const ctaSecondaryStyle: CSSProperties = {
  background: 'transparent',
  color: 'var(--text-secondary)',
  fontFamily: monoFont,
  fontSize: '13px',
  padding: '10px 22px',
  border: '0.5px solid var(--border-strong)',
  borderRadius: '2px',
  letterSpacing: '0.04em',
  textTransform: 'none',
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  textDecoration: 'none',
  transition: 'background 0.2s, color 0.2s',
}

export const statsRowStyle: CSSProperties = {
  fontFamily: monoFont,
  fontSize: 'clamp(11px, 3vw, 12px)',
  letterSpacing: '0.04em',
  color: 'var(--text-muted)',
  marginBottom: '2.5rem',
  display: 'flex',
  flexWrap: 'wrap',
  justifyContent: 'center',
  gap: '2px 8px',
}

export const codeQuoteStyle: CSSProperties = {
  background: 'transparent',
  border: 'none',
  borderLeft: '2px solid var(--accent-blue)',
  borderRadius: 0,
  padding: '0.5rem 0 0.5rem 1rem',
  fontFamily: monoFont,
  fontSize: 'clamp(11px, 2.5vw, 13px)',
  color: 'var(--text-secondary)',
  whiteSpace: 'pre',
  overflowX: 'auto',
  textAlign: 'left',
  lineHeight: 1.6,
  display: 'inline-block',
}
