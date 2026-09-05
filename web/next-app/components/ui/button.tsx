'use client'

import * as React from 'react'
import { Slot } from '@radix-ui/react-slot'
import { cva, type VariantProps } from 'class-variance-authority'
import { cn } from '@/lib/utils'

// Mathua skin: square corners, mono micro-labels, hairline borders.
// No rounded-*, no slate — all color flows through mathua/sidebar tokens.
const buttonVariants = cva(
  'inline-flex items-center justify-center gap-2 whitespace-nowrap font-mono text-xs transition-colors focus-visible:outline-none focus-visible:border-mathua-blue disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0',
  {
    variants: {
      variant: {
        default: 'bg-mathua-blue text-white hover:bg-mathua-blue-hover',
        destructive: 'bg-mathua-red text-white hover:opacity-90',
        outline:
          'border border-mathua-border-strong bg-transparent text-mathua-muted hover:text-mathua-blue hover:border-mathua-blue',
        secondary: 'bg-mathua-surface-elevated text-mathua-primary hover:bg-mathua-surface-highlight',
        ghost: 'text-mathua-muted hover:text-mathua-blue hover:bg-mathua-surface-elevated',
        link: 'text-mathua-blue underline-offset-4 hover:underline',
      },
      size: {
        default: 'h-10 px-4 py-2 min-h-[44px]',
        sm: 'h-9 px-3 min-h-[36px]',
        lg: 'h-11 px-8 min-h-[44px]',
        icon: 'h-10 w-10 min-h-[44px] min-w-[44px]',
      },
    },
    defaultVariants: {
      variant: 'default',
      size: 'default',
    },
  }
)

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {
  asChild?: boolean
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, asChild = false, ...props }, ref) => {
    const Comp = asChild ? Slot : 'button'
    return <Comp className={cn(buttonVariants({ variant, size, className }))} ref={ref} {...props} />
  }
)
Button.displayName = 'Button'

export { Button, buttonVariants }
