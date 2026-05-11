import { ReactNode } from 'react'

interface ButtonProps {
  children: ReactNode
  variant?: 'primary' | 'outline' | 'ghost'
  href?: string
  onClick?: () => void
  className?: string
}

const variants = {
  primary:
    'text-mathua-bg font-mono',
  outline:
    'border border-mathua-border-strong text-mathua-secondary font-mono',
  ghost:
    'text-mathua-muted font-sans',
}

export default function Button({
  children,
  variant = 'primary',
  href,
  onClick,
  className = '',
}: ButtonProps) {
  const classes = `inline-flex items-center justify-center h-10 px-6 font-mono text-[13px] tracking-[0.04em] transition-all duration-200 cursor-pointer ${variant === 'primary' ? 'bg-[var(--accent-gold)] text-[var(--bg)]' : variant === 'outline' ? 'border text-[var(--text-secondary)]' : ''} ${className}`

  if (href) {
    return (
      <a href={href} className={classes}>
        {children}
      </a>
    )
  }

  return (
    <button onClick={onClick} className={classes}>
      {children}
    </button>
  )
}
