'use client'

interface MasteryBadgeProps {
  status?: string
  size?: 'sm' | 'md'
}

const statusColors = {
  MASTERED: 'bg-green-500',
  PRACTICING: 'bg-yellow-400',
  LEARNING: 'bg-yellow-600',
  DECAYING: 'bg-orange-400',
  UNSEEN: 'bg-mathua-border',
} satisfies Record<string, string>

const statusLabels = {
  MASTERED: 'M',
  PRACTICING: 'P',
  LEARNING: 'L',
  DECAYING: 'D',
  UNSEEN: '',
} satisfies Record<string, string>

export default function MasteryBadge({ status, size = 'sm' }: MasteryBadgeProps) {
  const color = statusColors[status || 'UNSEEN'] || 'bg-mathua-border'
  const label = statusLabels[status || 'UNSEEN'] || ''
  const dim = size === 'sm' ? 'w-2.5 h-2.5' : 'w-3.5 h-3.5'
  const fontSize = size === 'sm' ? 'text-[6px]' : 'text-[8px]'

  return (
    <span
      className={`inline-flex items-center justify-center rounded-full ${dim} ${color} ${fontSize} text-white font-mono`}
      title={status || 'UNSEEN'}
    >
      {label}
    </span>
  )
}
