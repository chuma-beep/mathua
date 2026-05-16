'use client'

import { useMemo } from 'react'
import { renderMermaidSVG } from 'beautiful-mermaid'
import { useTheme } from '../hooks/useTheme'

interface MermaidDiagramProps {
  code: string
  className?: string
}

const LIGHT_COLORS = {
  bg: '#ffffff',
  fg: '#18181b',
  accent: '#2563eb',
  line: '#2563eb',
  muted: '#71717a',
  surface: '#e4e4e7',
  border: '#a1a1aa',
}

const DARK_COLORS = {
  bg: '#18181b',
  fg: '#fafafa',
  accent: '#60a5fa',
  line: '#60a5fa',
  muted: '#a1a1aa',
  surface: '#27272a',
  border: '#52525b',
}

export default function MermaidDiagram({ code, className = '' }: MermaidDiagramProps) {
  const { theme, mounted } = useTheme()

  const { svg, error } = useMemo(() => {
    if (!mounted) return { svg: null, error: null }
    try {
      const colors = theme === 'light' ? LIGHT_COLORS : DARK_COLORS
      return {
        svg: renderMermaidSVG(code, {
          ...colors,
          font: "'IBM Plex Mono', 'IBM Plex Serif', monospace",
          transparent: true,
        }),
        error: null,
      }
    } catch (err) {
      return { svg: null, error: err instanceof Error ? err.message : String(err) }
    }
  }, [code, theme, mounted])

  if (error) {
    return (
      <pre style={{
        fontFamily: "'IBM Plex Mono', monospace",
        fontSize: '12px',
        color: 'var(--accent-red)',
        borderLeft: '2px solid var(--accent-red)',
        padding: '8px 16px',
        margin: '16px 0',
        whiteSpace: 'pre-wrap',
      }}>
        {error}
      </pre>
    )
  }

  if (!mounted || !svg) {
    return (
      <div
        className={className}
        style={{
          display: 'flex',
          justifyContent: 'center',
          margin: '24px 0',
          minHeight: '120px',
        }}
      />
    )
  }

  return (
    <div
      className={className}
      style={{
        display: 'flex',
        justifyContent: 'center',
        margin: '24px 0',
      }}
      dangerouslySetInnerHTML={{ __html: svg }}
    />
  )
}
