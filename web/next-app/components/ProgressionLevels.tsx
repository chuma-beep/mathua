interface Level {
  num: string
  name: string
  range: string
  elite?: boolean
}

interface ProgressionLevelsProps {
  levels: Level[]
  className?: string
}

export default function ProgressionLevels({ levels, className = '' }: ProgressionLevelsProps) {
  return (
    <div className={className}>
      <h3
        style={{
          fontFamily: "'DM Serif Display', serif",
          fontWeight: 400,
          fontSize: '1.1rem',
          color: 'var(--text-primary)',
          marginBottom: '0.75rem',
          letterSpacing: '-0.01em',
        }}
      >
        Levels of mastery
      </h3>
      <table style={{ width: '100%', borderCollapse: 'collapse', border: 'none', background: 'transparent' }}>
        <tbody>
          {levels.map((level) => (
            <tr
              key={level.num}
              style={{
                borderBottom: '0.5px solid var(--border)',
                background: 'transparent',
                ...(level.elite ? { borderTop: '1px solid var(--accent-gold)' } : {}),
              }}
            >
              <td
                style={{
                  fontFamily: "'DM Serif Display', serif",
                  fontSize: '13px',
                  color: 'var(--accent-gold)',
                  padding: '8px 12px 8px 0',
                  whiteSpace: 'nowrap',
                  verticalAlign: 'baseline',
                }}
              >
                {parseInt(level.num)}.
              </td>
              <td
                style={{
                  fontFamily: "'DM Serif Display', serif",
                  fontSize: '1rem',
                  color: 'var(--text-primary)',
                  padding: '8px 0',
                  verticalAlign: 'baseline',
                }}
              >
                {level.name}
              </td>
              <td
                style={{
                  fontFamily: "'DM Serif Display', serif",
                  fontSize: '12px',
                  color: 'var(--text-muted)',
                  textAlign: 'right',
                  padding: '8px 0 8px 12px',
                  whiteSpace: 'nowrap',
                  verticalAlign: 'baseline',
                }}
              >
                {level.range} concepts
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
