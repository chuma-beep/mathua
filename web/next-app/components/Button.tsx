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
    'bg-mathua-blue text-white hover:bg-mathua-blue-hover',
  outline:
    'border border-mathua-border-strong text-mathua-primary hover:border-mathua-blue hover:text-mathua-blue',
  ghost:
    'text-mathua-secondary hover:bg-mathua-surface-elevated',
}

export default function Button({
  children,
  variant = 'primary',
  href,
  onClick,
  className = '',
}: ButtonProps) {
  const classes = `inline-flex items-center justify-center h-10 px-6 rounded-md text-[13px] font-medium transition-all duration-200 cursor-pointer ${variants[variant]} ${className}`

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
