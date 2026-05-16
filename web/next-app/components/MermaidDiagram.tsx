'use client'

import { useMemo } from 'react'
import { renderMermaidSVG } from 'beautiful-mermaid'

interface MermaidDiagramProps {
  code: string
  className?: string
}

export default function MermaidDiagram({ code, className = '' }: MermaidDiagramProps) {
  const { svg, error } = useMemo(() => {
    try {
      return {
        svg: renderMermaidSVG(code, {
          bg: 'var(--bg)',
          fg: 'var(--text-primary)',
          accent: 'var(--accent-blue)',
          muted: 'var(--text-muted)',
          surface: 'var(--surface)',
          border: 'var(--border-strong)',
          line: 'var(--border-strong)',
          font: "'IBM Plex Mono', 'IBM Plex Serif', monospace",
          transparent: true,
        }),
        error: null,
      }
    } catch (err) {
      return { svg: null, error: err instanceof Error ? err.message : String(err) }
    }
  }, [code])

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

  return (
    <div
      className={className}
      style={{
        display: 'flex',
        justifyContent: 'center',
        margin: '24px 0',
      }}
      dangerouslySetInnerHTML={{ __html: svg! }}
    />
  )
}
