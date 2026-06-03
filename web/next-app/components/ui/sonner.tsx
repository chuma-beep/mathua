'use client'

import { Toaster as Sonner, toast } from 'sonner'
import 'sonner/dist/styles.css'

export { toast }

type ToasterProps = React.ComponentProps<typeof Sonner>

export function Toaster({ ...props }: ToasterProps) {
  return (
    <Sonner
      theme="system"
      position="bottom-right"
      closeButton
      {...props}
    />
  )
}
