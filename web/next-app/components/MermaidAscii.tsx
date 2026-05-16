'use client'

import { useMemo } from 'react'
import { renderMermaidASCII } from 'beautiful-mermaid'

interface MermaidAsciiProps {
  code: string
  useAscii?: boolean
  className?: string
}

export default function MermaidAscii({ code, useAscii = false, className = '' }: MermaidAsciiProps) {
  const ascii = useMemo(() => {
    try {
      return renderMermaidASCII(code, { useAscii })
    } catch {
      return null
    }
  }, [code, useAscii])

  if (!ascii) return null

  return (
    <pre
      className={`ascii-text ${className}`}
      style={{
        fontSize: '11px',
        lineHeight: 1.3,
        margin: '16px 0',
      }}
      aria-hidden="true"
    >
      {ascii}
    </pre>
  )
}
