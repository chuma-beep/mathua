'use client'

import * as React from 'react'
import { cn } from '@/lib/utils'

// Mathua skin: square corners, hairline borders, mono text, no ring/shadow.
const Input = React.forwardRef<HTMLInputElement, React.ComponentProps<'input'>>(
  ({ className, type, ...props }, ref) => (
    <input
      type={type}
      ref={ref}
      className={cn(
        'w-full min-w-0 h-12 rounded-none border border-mathua-border bg-mathua-code px-4 font-mono text-base text-mathua-primary placeholder:text-mathua-muted outline-none transition-colors focus:border-mathua-blue disabled:cursor-not-allowed disabled:opacity-50',
        className
      )}
      {...props}
    />
  )
)
Input.displayName = 'Input'

export { Input }
