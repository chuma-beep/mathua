interface DomainRow {
  domain: string
  count: number
  comingSoon?: boolean
}

interface DomainTableProps {
  rows: DomainRow[]
  className?: string
}

const thStyle: React.CSSProperties = {
  fontFamily: "'IBM Plex Mono', monospace",
  fontWeight: 400,
  fontSize: '0.8rem',
  textTransform: 'uppercase',
  letterSpacing: '0.1em',
  color: 'var(--accent-blue)',
  padding: '10px 20px 10px 0',
  borderBottom: '1px solid var(--accent-blue)',
  background: 'transparent',
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
                ...thStyle,
                textAlign: 'left',
              }}
            >
              Domain
            </th>
            <th
              style={{
                ...thStyle,
                textAlign: 'right',
                padding: '10px 0 10px 20px',
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
                  fontFamily: "'IBM Plex Mono', monospace",
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
                      fontFamily: "'IBM Plex Mono', monospace",
                      fontStyle: 'italic',
                      fontSize: '0.9rem',
                      color: 'var(--text-muted)',
                      marginLeft: '0.5rem',
                    }}
                  >
                    : coming in v1.1
                  </span>
                )}
              </td>
              <td
                style={{
                  fontFamily: "'IBM Plex Mono', monospace",
                  fontSize: '0.9rem',
                  color: 'var(--accent-blue)',
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
