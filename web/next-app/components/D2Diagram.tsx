'use client'

import Image from 'next/image'

const diagramMap: Record<string, { light: string; dark: string }> = {
  architecture: {
    light: '/diagrams/architecture-light.svg',
    dark: '/diagrams/architecture-dark.svg',
  },
  'cat-diagnostic': {
    light: '/diagrams/cat-diagnostic-light.svg',
    dark: '/diagrams/cat-diagnostic-dark.svg',
  },
  contributing: {
    light: '/diagrams/contributing-light.svg',
    dark: '/diagrams/contributing-dark.svg',
  },
  platforms: {
    light: '/diagrams/platforms-light.svg',
    dark: '/diagrams/platforms-dark.svg',
  },
  'system-design-architecture': {
    light: '/diagrams/system-design-architecture-light.svg',
    dark: '/diagrams/system-design-architecture-dark.svg',
  },
  'system-design-data-model': {
    light: '/diagrams/system-design-data-model-light.svg',
    dark: '/diagrams/system-design-data-model-dark.svg',
  },
  'system-design-request-flow': {
    light: '/diagrams/system-design-request-flow-light.svg',
    dark: '/diagrams/system-design-request-flow-dark.svg',
  },
  'system-design-sm2-mastery': {
    light: '/diagrams/system-design-sm2-mastery-light.svg',
    dark: '/diagrams/system-design-sm2-mastery-dark.svg',
  },
  'system-design-grading-pipeline': {
    light: '/diagrams/system-design-grading-pipeline-light.svg',
    dark: '/diagrams/system-design-grading-pipeline-dark.svg',
  },
  'system-design-cat-diagnostic': {
    light: '/diagrams/system-design-cat-diagnostic-light.svg',
    dark: '/diagrams/system-design-cat-diagnostic-dark.svg',
  },
}

interface D2DiagramProps {
  name: keyof typeof diagramMap
  theme?: 'dark' | 'light'
  className?: string
}

export default function D2Diagram({ name, theme = 'dark', className = '' }: D2DiagramProps) {
  const pair = diagramMap[name]
  if (!pair) return <div style={{ color: 'var(--accent-red)', fontFamily: 'IBM Plex Mono, Fira Code, monospace', fontSize: '12px', padding: '1rem' }}>Diagram "{name}" not found</div>

  return (
    <div className={`flex justify-center py-4 ${className}`}>
      <Image
        src={theme === 'dark' ? pair.dark : pair.light}
        alt={`${name} diagram`}
        width={800}
        height={400}
        style={{ maxWidth: '100%', height: 'auto' }}
      />
    </div>
  )
}
