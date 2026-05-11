interface DomainRow {
  domain: string
  count: number
  comingSoon?: boolean
}

interface DomainTableProps {
  rows: DomainRow[]
  className?: string
}

export default function DomainTable({ rows, className = '' }: DomainTableProps) {
  return (
    <div className={className}>
      <table
        style={{
          width: '100%',
          borderCollapse: 'collapse',
          border: 'none',
          background: 'transparent',
        }}
      >
        <thead>
          <tr>
            <th
              style={{
                fontFamily: "'JetBrains Mono', 'Fira Code', monospace",
                fontWeight: 400,
                fontSize: '0.8rem',
                textTransform: 'uppercase',
                letterSpacing: '0.1em',
                color: 'var(--accent-gold)',
                textAlign: 'left',
                padding: '10px 20px 10px 0',
                borderBottom: '1px solid var(--accent-gold)',
                background: 'transparent',
              }}
            >
              Domain
            </th>
            <th
              style={{
                fontFamily: "'JetBrains Mono', 'Fira Code', monospace",
                fontWeight: 400,
                fontSize: '0.8rem',
                textTransform: 'uppercase',
                letterSpacing: '0.1em',
                color: 'var(--accent-gold)',
                textAlign: 'right',
                padding: '10px 0 10px 20px',
                borderBottom: '1px solid var(--accent-gold)',
                background: 'transparent',
              }}
            >
              Concepts
            </th>
          </tr>
        </thead>
        <tbody>
          {rows.map((row) => (
            <tr key={row.domain} style={{ borderBottom: '0.5px solid var(--border)', background: 'transparent' }}>
              <td
                style={{
                  fontFamily: "'JetBrains Mono', 'Fira Code', monospace",
                  fontSize: '1rem',
                  color: 'var(--text-secondary)',
                  padding: '14px 20px 14px 0',
                  background: 'transparent',
                }}
              >
                {row.domain}
                {row.comingSoon && (
                  <span
                    style={{
                      fontFamily: "'JetBrains Mono', 'Fira Code', monospace",
                      fontStyle: 'italic',
                      fontSize: '0.9rem',
                      color: 'var(--text-muted)',
                      marginLeft: '0.5rem',
                    }}
                  >
                    — coming in v1.1
                  </span>
                )}
              </td>
              <td
                style={{
                  fontFamily: "'JetBrains Mono', 'Fira Code', monospace",
                  fontSize: '0.9rem',
                  color: 'var(--accent-gold)',
                  textAlign: 'right',
                  padding: '14px 0 14px 20px',
                  background: 'transparent',
                }}
              >
                {row.count}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
