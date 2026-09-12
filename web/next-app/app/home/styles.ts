import type { CSSProperties } from 'react'

export const loadingGraphStyle: CSSProperties = {
  height: '280px',
  background: 'transparent',
  borderRadius: 0,
  border: '0.5px solid var(--border)',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  color: 'var(--text-muted)',
  fontFamily: "var(--font-jetbrains-mono), 'JetBrains Mono', monospace",
  fontSize: '13px',
  width: '100%',
  maxWidth: '100%',
}
