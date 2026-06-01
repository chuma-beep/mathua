'use client'

import { Toaster as Sonner, toast } from 'sonner'

export { toast }

type ToasterProps = React.ComponentProps<typeof Sonner>

export function Toaster({ ...props }: ToasterProps) {
  return (
    <Sonner
      theme="system"
      position="bottom-right"
      closeButton
      toastOptions={{
        unstyled: true,
        classNames: {
          toast:
            'font-mono text-sm px-3 py-2 bg-[var(--surface)] text-[var(--text-primary)] border-[0.5px] border-[var(--border)] flex items-center gap-2 w-fit max-w-[360px]',
          title: 'text-[var(--text-primary)] font-mono text-sm',
          description: 'text-[var(--text-muted)] font-mono text-xs',
          actionButton:
            'bg-[var(--surface-highlight)] border-[0.5px] border-[var(--border)] text-[var(--text-primary)] px-2 py-0.5 text-xs font-mono',
          cancelButton:
            'bg-[var(--surface)] border-[0.5px] border-[var(--border)] text-[var(--text-muted)] px-2 py-0.5 text-xs font-mono',
          closeButton:
            'text-[var(--text-muted)] hover:text-[var(--text-primary)] border-none bg-transparent',
          error:
            'border-[var(--accent-red)]',
          success:
            'border-[var(--accent-green)]',
          info:
            'border-[var(--accent-blue)]',
          warning:
            'border-[var(--accent-yellow)]',
          icon: 'text-[var(--text-muted)]',
        },
      }}
      {...props}
    />
  )
}
