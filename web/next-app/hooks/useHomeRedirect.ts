'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { consumeHomeView } from '../lib/auth'
import { useAuthState } from './useAuthState'

// Landing rule: logged-in visits to / bounce to /profile (tokens persist
// across tabs via localStorage; cross-tab logins arrive through the
// auth-changed/storage sync). An explicit Home click sets a one-shot
// per-tab flag and behaves normally. Guests are unaffected.
// Returns true once the decision has been made (safe to render).
export function useHomeRedirect(mounted: boolean): boolean {
  const { loggedIn } = useAuthState()
  const { replace } = useRouter()
  const [checked, setChecked] = useState(false)

  useEffect(() => {
    if (!mounted) return
    if (consumeHomeView()) {
      setChecked(true)
      return
    }
    if (loggedIn) {
      replace('/profile')
      return
    }
    setChecked(true)
  }, [mounted, loggedIn, replace])

  return checked
}
