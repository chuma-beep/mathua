import { ReactNode } from 'react'

interface InfoCardProps {
  title: string
  children: ReactNode
  className?: string
}

export default function InfoCard({ title, children, className = '' }: InfoCardProps) {
  return (
    <div
      className={`bg-mathua-surface border border-mathua-border border-l-[3px] border-l-mathua-blue rounded-lg p-6 flex-1 min-w-[250px] ${className}`}
    >
      <h4 className="font-serif font-normal text-xl text-mathua-blue text-center mb-2.5">
        {title}
      </h4>
      <p className="text-mathua-secondary text-sm leading-relaxed text-center">
        {children}
      </p>
    </div>
  )
}
