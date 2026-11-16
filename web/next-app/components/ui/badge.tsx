import * as React from "react"
import { cva, type VariantProps } from "class-variance-authority"

import { cn } from "@/lib/utils"

const badgeVariants = cva(
  "inline-flex items-center rounded-none border px-1.5 py-0 font-mono text-[10px] uppercase tracking-wider transition-colors focus:outline-none",
  {
    variants: {
      variant: {
        default: "border-mathua-blue bg-mathua-blue-faint text-mathua-blue",
        secondary: "border-mathua-border bg-mathua-surface-elevated text-mathua-secondary",
        destructive: "border-mathua-red bg-mathua-red-faint text-mathua-red",
        outline: "border-mathua-border text-mathua-secondary",
        up: "border-mathua-green bg-mathua-green-faint text-mathua-green",
        down: "border-mathua-red bg-mathua-red-faint text-mathua-red",
        stay: "border-mathua-border text-mathua-muted",
      },
    },
    defaultVariants: {
      variant: "default",
    },
  }
)

export interface BadgeProps
  extends React.HTMLAttributes<HTMLDivElement>,
    VariantProps<typeof badgeVariants> {}

function Badge({ className, variant, ...props }: BadgeProps) {
  return (
    <div className={cn(badgeVariants({ variant }), className)} {...props} />
  )
}

export { Badge, badgeVariants }
