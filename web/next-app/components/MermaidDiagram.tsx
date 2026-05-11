'use client'

import { useEffect, useRef, useId } from 'react'
import mermaid from 'mermaid'

let mermaidInitialized = false

function initMermaid() {
  if (mermaidInitialized) return
  mermaid.initialize({
    startOnLoad: false,
    theme: 'base',
    securityLevel: 'sandbox',
    fontFamily: "'JetBrains Mono', 'Fira Code', monospace",
  })
  mermaidInitialized = true
}

interface MermaidDiagramProps {
  chart: string
  theme?: 'dark' | 'light'
  className?: string
}

export default function MermaidDiagram({ chart, theme = 'dark', className = '' }: MermaidDiagramProps) {
  const containerRef = useRef<HTMLDivElement>(null)
  const id = useId().replace(/[^a-zA-Z0-9]/g, '')
  const previousTheme = useRef<string | null>(null)

  const isDark = theme === 'dark'

  useEffect(() => {
    initMermaid()

    const themeVariables: Record<string, string> = isDark
      ? {
          primaryColor: 'var(--surface-elevated, #1e293b)',
          primaryTextColor: 'var(--text-primary, #e8e2d5)',
          primaryBorderColor: 'var(--accent-gold, #c8a96e)',
          lineColor: 'var(--accent-gold, #c8a96e)',
          secondaryColor: 'var(--surface, #0f172a)',
          tertiaryColor: 'var(--surface-highlight, #334155)',
          textColor: 'var(--text-secondary, #9fa8b4)',
          edgeLabelBackground: 'var(--bg, #0b0f1a)',
          nodeBorder: 'var(--border-strong, #2a3f5f)',
          clusterBkg: 'var(--bg, #0b0f1a)',
          clusterBorder: 'var(--border, #1e2d45)',
          titleColor: 'var(--text-primary, #e8e2d5)',
          fontSize: '13px',
        }
      : {
          primaryColor: 'var(--surface-elevated, #f5f1e6)',
          primaryTextColor: 'var(--text-primary, #1e293b)',
          primaryBorderColor: 'var(--accent-gold, #b8933e)',
          lineColor: 'var(--accent-gold, #b8933e)',
          secondaryColor: 'var(--surface, #faf8f0)',
          tertiaryColor: 'var(--surface-highlight, #e8e4d6)',
          textColor: 'var(--text-secondary, #475569)',
          edgeLabelBackground: 'var(--bg, #fefcf4)',
          nodeBorder: 'var(--border-strong, #b8b5a8)',
          clusterBkg: 'var(--bg, #fefcf4)',
          clusterBorder: 'var(--border, #d6d3c8)',
          titleColor: 'var(--text-primary, #1e293b)',
          fontSize: '13px',
        }

    mermaid.initialize({
      startOnLoad: false,
      theme: 'base',
      securityLevel: 'sandbox',
      themeVariables,
      fontFamily: "'JetBrains Mono', 'Fira Code', monospace",
    })

    if (previousTheme.current === theme && containerRef.current?.innerHTML) return

    const renderDiagram = async () => {
      if (!containerRef.current) return
      try {
        const { svg } = await mermaid.render(`mermaid-${id}`, chart.trim())
        containerRef.current.innerHTML = svg
        previousTheme.current = theme
      } catch (err) {
        containerRef.current.innerHTML = `<pre style="color:var(--accent-red);font-family:var(--mono);font-size:12px;padding:1rem;white-space:pre-wrap">Mermaid parse error:\n${err}</pre>`
      }
    }

    renderDiagram()
  }, [chart, theme, id, isDark])

  return (
    <div
      ref={containerRef}
      className={`flex justify-center overflow-x-auto py-4 ${className}`}
      style={{ background: 'var(--bg)' }}
    />
  )
}
