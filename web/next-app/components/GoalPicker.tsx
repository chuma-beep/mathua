'use client'

import { useEffect, useMemo, useRef, useState } from 'react'

export interface GoalOption {
  id: string
  label: string
  domain: string
}

interface GoalPickerProps {
  concepts: GoalOption[]
  currentGoalId?: string | null
  onSelect: (id: string) => void
  onClose: () => void
}

/**
 * A searchable concept picker for the graph's goal control.
 *
 * It exists because "select a node, then look below the graph" was the only way
 * to set a goal, and nobody found it. The picker is deliberately small: the
 * graph page already owns the concept list and the goal setter.
 */
export default function GoalPicker({ concepts, currentGoalId, onSelect, onClose }: GoalPickerProps) {
  const [query, setQuery] = useState('')
  const inputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    inputRef.current?.focus()
  }, [])

  useEffect(() => {
    const onKey = (event: KeyboardEvent) => {
      if (event.key === 'Escape') onClose()
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [onClose])

  const results = useMemo(() => {
    const q = query.trim().toLowerCase()
    const matching = q
      ? concepts.filter(
          concept =>
            concept.label.toLowerCase().includes(q) ||
            concept.id.toLowerCase().includes(q) ||
            concept.domain.toLowerCase().includes(q),
        )
      : concepts
    return matching.slice(0, 12)
  }, [concepts, query])

  return (
    <div
      role="presentation"
      onClick={onClose}
      className="fixed inset-0 z-50 flex items-start justify-center bg-black/60 px-4 pt-[12vh]"
    >
      <div
        role="dialog"
        aria-modal="true"
        aria-label="Choose a goal"
        onClick={event => event.stopPropagation()}
        className="w-full max-w-lg border border-mathua-border bg-mathua-surface p-4 font-mono"
      >
        <div className="flex items-start justify-between gap-4">
          <div className="min-w-0">
            <h2 className="text-sm text-mathua-primary">Choose a goal</h2>
            <p className="mt-1 text-[11px] text-mathua-secondary">
              The graph will keep only this concept and the path that leads to it.
            </p>
          </div>
          <button
            type="button"
            onClick={onClose}
            className="shrink-0 text-[10px] uppercase tracking-wider text-mathua-muted hover:text-mathua-primary"
          >
            Close
          </button>
        </div>

        <input
          ref={inputRef}
          value={query}
          onChange={event => setQuery(event.target.value)}
          placeholder="Search concepts..."
          aria-label="Search concepts for a goal"
          className="mt-4 w-full min-h-[44px] border-[0.5px] border-mathua-border-strong bg-mathua-surface-elevated px-3 text-[12px] text-mathua-primary outline-none"
        />

        <div className="mt-2 max-h-[min(55dvh,360px)] overflow-y-auto">
          {results.length === 0 ? (
            <p className="px-1 py-4 text-[11px] text-mathua-secondary">No matching concepts.</p>
          ) : (
            <ul className="m-0 list-none p-0">
              {results.map(concept => (
                <li key={concept.id} className="border-b-[0.5px] border-mathua-border last:border-b-0">
                  <button
                    type="button"
                    onClick={() => onSelect(concept.id)}
                    className="flex w-full items-center gap-2 px-1 py-2.5 text-left text-[11px] text-mathua-secondary hover:bg-mathua-code"
                  >
                    <span className="min-w-0 flex-1 truncate">{concept.label}</span>
                    <span className="shrink-0 text-[9px] uppercase tracking-wider text-mathua-muted">
                      {currentGoalId === concept.id ? 'Current goal' : concept.domain.replace(/_/g, ' ')}
                    </span>
                  </button>
                </li>
              ))}
            </ul>
          )}
        </div>
      </div>
    </div>
  )
}
