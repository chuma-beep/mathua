'use client'

import { useEffect, useRef, useId } from 'react'
import mermaid from 'mermaid'

interface MermaidDiagramProps {
  chart: string
  theme?: 'dark' | 'light'
  className?: string
}

const CSS_FALLBACKS: Record<string, string> = {
  '--bg': '#0b0f1a',
  '--surface': '#0f172a',
  '--surface-elevated': '#1e293b',
  '--surface-highlight': '#334155',
  '--text-primary': '#e8e2d5',
  '--text-secondary': '#9fa8b4',
  '--accent-gold': '#c8a96e',
  '--border': '#1e2d45',
  '--border-strong': '#2a3f5f',
}

function resolveVar(name: string): string {
  const raw = getComputedStyle(document.documentElement).getPropertyValue(name).trim()
  if (raw) return raw
  return CSS_FALLBACKS[name] || '#333'
}

function buildThemeVars(isDark: boolean): Record<string, string> {
  return {
    primaryColor: resolveVar('--surface-elevated'),
    primaryTextColor: resolveVar('--text-primary'),
    primaryBorderColor: resolveVar('--accent-gold'),
    lineColor: resolveVar('--accent-gold'),
    secondaryColor: resolveVar('--surface'),
    tertiaryColor: resolveVar('--surface-highlight'),
    textColor: resolveVar('--text-secondary'),
    edgeLabelBackground: resolveVar('--bg'),
    nodeBorder: resolveVar('--border-strong'),
    clusterBkg: resolveVar('--bg'),
    clusterBorder: resolveVar('--border'),
    titleColor: resolveVar('--text-primary'),
    fontSize: '13px',
    mainBkg: resolveVar('--bg'),
  }
}

export default function MermaidDiagram({ chart, theme = 'dark', className = '' }: MermaidDiagramProps) {
  const containerRef = useRef<HTMLDivElement>(null)
  const id = useId().replace(/[^a-zA-Z0-9]/g, '')
  const previousTheme = useRef<string | null>(null)

  const isDark = theme === 'dark'

  useEffect(() => {
    const themeVariables = buildThemeVars(isDark)

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
        containerRef.current.innerHTML = `<pre style="color:var(--accent-red);font-family:JetBrains Mono,Fira Code,monospace;font-size:12px;padding:1rem;white-space:pre-wrap">Mermaid parse error:\n${err}</pre>`
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
