'use client'

import { useEffect, useState } from 'react'
import { adminMe, type AdminMe } from '../lib/api'

// The caller's admin capability, fetched once per page from /api/admin/me.
//
// This is the client's copy of the permission list, and it is a convenience: it decides which
// sections render, never what the server allows. A section hidden here is still refused by the
// server if reached another way, and a section shown here still fails with a 403 if the account
// lacks the permission — which is what keeps the two from being the same decision.
export function useAdminMe() {
  const [me, setMe] = useState<AdminMe | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    let cancelled = false
    adminMe()
      .then((m) => { if (!cancelled) setMe(m) })
      .catch(() => { if (!cancelled) setMe(null) })
      .finally(() => { if (!cancelled) setLoading(false) })
    return () => { cancelled = true }
  }, [])

  const can = (permission: string) => me?.permissions.includes(permission) === true
  return { me, loading, can, isStaff: me !== null }
}
