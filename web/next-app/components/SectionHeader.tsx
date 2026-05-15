interface SectionHeaderProps {
  label?: string
  title: string
  center?: boolean
  className?: string
}

export default function SectionHeader({ label, title, center = false, className = '' }: SectionHeaderProps) {
  return (
    <div className={`${center ? 'text-center' : ''} ${className}`}>
      {label && (
        <span className="font-mono text-[11px] uppercase tracking-[0.15em] text-mathua-muted block mb-2">
          {label}
        </span>
      )}
      <h2
        className="font-serif font-normal text-[1.9rem] max-sm:text-[1.5rem] leading-tight text-mathua-primary"
        style={{
          fontFamily: "'IBM Plex Serif', serif",
          fontWeight: 400,
          color: 'var(--text-primary)',
          borderLeft: '2px solid var(--accent-blue)',
          paddingLeft: '1rem',
          marginBottom: '1.5rem',
        }}
      >
        {title}
      </h2>
    </div>
  )
}
