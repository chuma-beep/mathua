'use client'

import { useState, useMemo, useEffect, useRef } from 'react'
import Fuse from 'fuse.js'

interface SearchItem {
  title: string
  body: string
  domain: string
  concepts: string[]
}

interface SearchBarProps {
  items: SearchItem[]
  onSelect: (item: SearchItem) => void
  placeholder?: string
}

export default function SearchBar({ items, onSelect, placeholder = 'Search lessons…' }: SearchBarProps) {
  const [query, setQuery] = useState('')
  const [showResults, setShowResults] = useState(false)
  const ref = useRef<HTMLDivElement>(null)

  const fuse = useMemo(() => {
    return new Fuse(items, {
      keys: [
        { name: 'title', weight: 2 },
        { name: 'concepts', weight: 1.5 },
        { name: 'body', weight: 1 },
        { name: 'domain', weight: 1 },
      ],
      threshold: 0.4,
      minMatchCharLength: 2,
    })
  }, [items])

  const results = useMemo(() => {
    if (!query.trim()) return []
    return fuse.search(query.trim()).slice(0, 10)
  }, [query, fuse])

  useEffect(() => {
    function handleClick(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) {
        setShowResults(false)
      }
    }
    document.addEventListener('mousedown', handleClick)
    return () => document.removeEventListener('mousedown', handleClick)
  }, [])

  return (
    <div ref={ref} className="relative w-full max-w-xl mx-auto mb-6">
      <div className="relative">
        <svg
          className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-mathua-muted"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
        <input
          type="text"
          value={query}
          onChange={e => { setQuery(e.target.value); setShowResults(true) }}
          onFocus={() => setShowResults(true)}
          placeholder={placeholder}
          className="w-full h-11 pl-10 pr-4 border border-mathua-border bg-mathua-surface text-sm font-mono text-mathua-primary placeholder:text-mathua-muted focus:outline-none focus:border-mathua-blue rounded-none"
        />
        {query && (
          <button
            onClick={() => { setQuery(''); setShowResults(false) }}
            className="absolute right-3 top-1/2 -translate-y-1/2 text-mathua-muted hover:text-mathua-primary text-xs"
          >
            ✕
          </button>
        )}
      </div>
      {showResults && query.trim() && (
        <div className="absolute top-full left-0 right-0 mt-1 border border-mathua-border bg-mathua-surface z-50 max-h-80 overflow-y-auto">
          {results.length === 0 ? (
            <div className="p-3 text-xs font-mono text-mathua-muted text-center">
              No results found
            </div>
          ) : (
            results.map((r, i) => (
              <button
                key={i}
                onClick={() => { onSelect(r.item); setShowResults(false); setQuery('') }}
                className="w-full text-left p-3 hover:bg-mathua-code border-b border-mathua-border last:border-b-0 transition-colors"
              >
                <div className="font-mono text-sm text-mathua-primary">{r.item.title}</div>
                <div className="font-mono text-[10px] text-mathua-muted mt-0.5">
                  {r.item.domain} &middot; {r.item.concepts.length} concept{r.item.concepts.length !== 1 ? 's' : ''}
                </div>
                <div className="font-mono text-[10px] text-mathua-secondary mt-0.5 line-clamp-1">
                  {r.item.body.replace(/[#*\[\]()\\$]/g, '').slice(0, 120)}
                </div>
              </button>
            ))
          )}
        </div>
      )}
    </div>
  )
}
