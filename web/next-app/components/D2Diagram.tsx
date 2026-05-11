'use client'

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
  'student-model': {
    light: '/diagrams/student-model-light.svg',
    dark: '/diagrams/student-model-dark.svg',
  },
}

interface D2DiagramProps {
  name: keyof typeof diagramMap
  theme?: 'dark' | 'light'
  className?: string
}

export default function D2Diagram({ name, theme = 'dark', className = '' }: D2DiagramProps) {
  const pair = diagramMap[name]
  if (!pair) return <div style={{ color: 'var(--accent-red)', fontFamily: 'JetBrains Mono, Fira Code, monospace', fontSize: '12px', padding: '1rem' }}>Diagram "{name}" not found</div>

  return (
    <div className={`flex justify-center py-4 ${className}`}>
      <img
        src={theme === 'dark' ? pair.dark : pair.light}
        alt={`${name} diagram`}
        style={{ maxWidth: '100%', height: 'auto' }}
      />
    </div>
  )
}
