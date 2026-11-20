'use client'

import { useEffect, useState } from 'react'
import { Toaster as Sonner, toast } from 'sonner'
import 'sonner/dist/styles.css'

export { toast }

type ToasterProps = React.ComponentProps<typeof Sonner>

/**
 * Toasts come up from the bottom on every screen size. On a phone the bottom
 * edge holds the tab bar and the support bubble, so the toast is centred and
 * lifted above both rather than moved to the top, where it covered the page
 * heading.
 */
export function Toaster({ ...props }: ToasterProps) {
  const [layout, setLayout] = useState<{
    position: 'bottom-right' | 'bottom-center'
    offset: number
  }>({ position: 'bottom-right', offset: 132 })

  useEffect(() => {
    const query = window.matchMedia('(max-width: 640px)')
    const apply = () =>
      setLayout(
        query.matches
          ? { position: 'bottom-center', offset: 24 }
          : { position: 'bottom-right', offset: 132 },
      )
    apply()
    query.addEventListener('change', apply)
    return () => query.removeEventListener('change', apply)
  }, [])

  return (
    <Sonner
      theme="system"
      position={layout.position}
      offset={layout.offset}
      closeButton
      {...props}
    />
  )
}
