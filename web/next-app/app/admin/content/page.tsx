'use client'

import { useEffect, useState } from 'react'
import Header from '../../../components/Header'
import BottomTabs from '../../../components/BottomTabs'
import Footer from '../../../components/Footer'
import SectionHeader from '../../../components/SectionHeader'
import AdminNav from '../../../components/AdminNav'
import Loading from '../../../components/Loading'
import { listAdminContent, getErrorMessage, type AdminContentConcept } from '../../../lib/api'

// Content inspection.
//
// Read-only, deliberately. The corpus is code-adjacent data — concepts and generated questions —
// with no safe write path in the data model, and inventing one would be the fragile editor the
// brief warns against. What an operator actually needs from here is to answer "what is this
// concept, what does it depend on, and has anyone complained about it", which is a read.
export default function AdminContentPage() {
  const [q, setQ] = useState('')
  const [items, setItems] = useState<AdminContentConcept[] | null>(null)
  const [total, setTotal] = useState(0)
  const [error, setError] = useState('')

  useEffect(() => {
    let cancelled = false
    const t = setTimeout(() => {
      listAdminContent(q)
        .then((r) => { if (!cancelled) { setItems(r.concepts); setTotal(r.total) } })
        .catch((e) => { if (!cancelled) { setError(getErrorMessage(e)); setItems([]) } })
    }, 200)
    return () => { cancelled = true; clearTimeout(t) }
  }, [q])

  return (
    <>
      <Header />
      <div className="pt-[var(--chrome-top)] lg:pt-0">
        <main className="mx-auto w-full max-w-[900px] min-w-0 px-4 sm:px-6 py-8 sm:py-12 overflow-x-hidden">
          <SectionHeader label="Administration" title="Content" />
          <AdminNav />

          <input
            type="search"
            value={q}
            onChange={(e) => setQ(e.target.value)}
            placeholder="Search concepts by id or label"
            className="mb-4 w-full border border-mathua-border bg-mathua-code px-3 h-11 font-mono text-xs text-mathua-primary"
          />

          {error && <p role="alert" className="mb-4 border border-mathua-red bg-mathua-surface p-4 font-mono text-xs text-mathua-red">{error}</p>}
          {!items && !error && <Loading label="LOADING CONCEPTS" />}

          {items && (
            <>
              <p className="mb-3 font-mono text-[11px] text-mathua-muted">{total} concept{total === 1 ? '' : 's'}</p>
              <ul className="space-y-1">
                {items.map((c) => (
                  <li key={c.id} className="border border-mathua-border bg-mathua-surface px-4 py-3">
                    <div className="flex items-baseline justify-between gap-3">
                      <span className="font-mono text-xs text-mathua-primary">{c.label}</span>
                      <span className="font-mono text-[10px] text-mathua-muted">{c.id}</span>
                    </div>
                    <div className="mt-1 font-mono text-[10px] text-mathua-secondary">
                      {c.domain}{c.subdomain ? ` / ${c.subdomain}` : ''} · {c.grading_type}
                      {c.prerequisites.length > 0 && ` · needs ${c.prerequisites.length}`}
                    </div>
                  </li>
                ))}
              </ul>
            </>
          )}
        </main>
      </div>
      <BottomTabs />
      <Footer />
    </>
  )
}
