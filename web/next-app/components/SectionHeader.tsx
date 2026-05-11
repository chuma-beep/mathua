interface SectionHeaderProps {
  label: string
  title: string
  center?: boolean
  className?: string
}

export default function SectionHeader({ label, title, center = true, className = '' }: SectionHeaderProps) {
  return (
    <div className={`${center ? 'text-center' : ''} ${className}`}>
      <span className="section-label">{label}</span>
      <h2 className="font-serif font-medium text-[28px] leading-tight text-mathua-primary">
        {title}
      </h2>
    </div>
  )
}
