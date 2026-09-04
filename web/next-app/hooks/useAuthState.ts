'use client'

import { useState, useEffect } from 'react'
import { isLoggedIn, getUserInfo, type UserInfo } from '../lib/auth'

export function useAuthState(): { loggedIn: boolean; user: UserInfo | null } {
  const [loggedIn, setLoggedIn] = useState(() => typeof window !== 'undefined' && isLoggedIn())
  const [user, setUser] = useState<UserInfo | null>(() => (typeof window !== 'undefined' ? getUserInfo() : null))

  useEffect(() => {
    const sync = () => {
      setLoggedIn(isLoggedIn())
      setUser(getUserInfo())
    }
    sync()
    window.addEventListener('auth-changed', sync)
    window.addEventListener('storage', sync)
    return () => {
      window.removeEventListener('auth-changed', sync)
      window.removeEventListener('storage', sync)
    }
  }, [])

  return { loggedIn, user }
}
