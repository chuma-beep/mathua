'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { useAuthState } from './useAuthState'

// Landing rule: logged-in visits to / bounce to /profile (tokens persist
// across tabs via localStorage; cross-tab logins arrive through the
// auth-changed/storage sync). The header logo points logged-in users
// straight at /profile; marketing content stays reachable via
// /docs, /how-it-works, and /note, which never bounce.
// Returns true once the decision has been made (safe to render).
export function useHomeRedirect(mounted: boolean): boolean {
  const { loggedIn } = useAuthState()
  const { replace } = useRouter()
  const [checked, setChecked] = useState(false)

  useEffect(() => {
    if (!mounted) return
    if (loggedIn) {
      replace('/profile')
      return
    }
    setChecked(true)
  }, [mounted, loggedIn, replace])

  return checked
}
