'use client'

import { useState, useEffect } from 'react'
import { isLoggedIn, getUserInfo, fetchMe, type UserInfo } from '../lib/auth'

export function useAuthState() {
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

  // The role is re-read from the server, once per page load.
  //
  // `mathua_user` is a cache of `/api/auth/me`, and a cache of a privilege is only as fresh as
  // its last write: a phone that signed in before Admin V2 carries no role, and one that signed
  // in after it had `moderator`/`owner` dropped by a comparison that listed only two of the four
  // roles. Either way `isStaffRole(user.role)` answered false for good, so the Admin entry never
  // appeared on that device while a re-authenticated desktop showed it.
  //
  // `fetchMe` repairs the stored session and broadcasts `auth-changed`, so this resolves through
  // the same listener the sign-in path uses rather than being a parallel source of truth. It runs
  // only when a role is missing — a learner is fetched once and then left alone, because
  // "the server said not staff" is worth caching and "we could not ask" is not.
  useEffect(() => {
    if (typeof window === 'undefined' || !isLoggedIn()) return
    if (getUserInfo()?.role) return
    let cancelled = false
    fetchMe().then((me) => {
      if (!cancelled && me?.role) setUser(me)
    })
    return () => { cancelled = true }
  }, [])

  return { loggedIn, user }
}