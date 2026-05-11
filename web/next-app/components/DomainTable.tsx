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
      <table className="w-full border-collapse bg-mathua-surface rounded-lg overflow-hidden border border-mathua-border">
        <thead>
          <tr>
            <th className="font-mono text-[11px] uppercase tracking-[0.1em] text-mathua-muted text-left p-[14px_20px] bg-mathua-surface-elevated border-b border-mathua-border">
              Domain
            </th>
            <th className="font-mono text-[11px] uppercase tracking-[0.1em] text-mathua-muted text-right p-[14px_20px] bg-mathua-surface-elevated border-b border-mathua-border">
              Concepts
            </th>
          </tr>
        </thead>
        <tbody>
          {rows.map((row) => (
            <tr key={row.domain} className="border-b border-mathua-border last:border-b-0">
              <td className="p-[14px_20px] text-mathua-primary text-sm">
                {row.domain}
                {row.comingSoon && (
                  <span className="text-[10px] text-mathua-blue bg-mathua-blue/10 px-2 py-0.5 rounded ml-2">
                    Coming soon
                  </span>
                )}
              </td>
              <td className="p-[14px_20px] font-mono text-mathua-muted text-right">
                {row.count}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
