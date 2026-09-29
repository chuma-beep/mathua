'use client'

import { useState, useEffect, useRef } from 'react'
import ReportButton from './ReportButton'
import type { ReportKind } from '../lib/api'

interface ReportMenuProps {
  conceptId?: string
  lessonId?: string
  kind: ReportKind
  question?: string
  blockId: string
}

// Icon-only per-block report menu. The block id is prefixed into the report
// payload so moderators can locate the exact example with no extra schema.
export default function ReportMenu({ conceptId, lessonId, kind, question, blockId }: ReportMenuProps) {
  const [open, setOpen] = useState(false)
  const ref = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (!open) return
    function onDown(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false)
    }
    function onKey(e: KeyboardEvent) {
      if (e.key === 'Escape') setOpen(false)
    }
    document.addEventListener('mousedown', onDown)
    document.addEventListener('keydown', onKey)
    return () => {
      document.removeEventListener('mousedown', onDown)
      document.removeEventListener('keydown', onKey)
    }
  }, [open ])

  return (
    <div ref={ref} className="relative shrink-0">
      <button
        type="button"
        onClick={() => setOpen(o => !o)}
        aria-label={`Report options, example ${blockId}`}
        aria-expanded={open}
        className="w-9 h-9 min-w-[36px] min-h-[36px] inline-flex items-center justify-center font-mono text-sm text-mathua-muted hover:text-mathua-blue hover:border-mathua-blue border border-transparent transition-colors"
      >
        ⋯
      </button>
      {open && (
        <div role="menu" className="absolute right-0 top-full mt-1 z-30 border border-mathua-border bg-mathua-surface p-3 min-w-[220px]">
          <ReportButton
            conceptId={conceptId}
            lessonId={lessonId}
            kind={kind}
            question={question ? `[block ${blockId}] ${question}` : `[block ${blockId}]`}
            label="Report a problem"
          />
        </div>
      )}
    </div>
  )
}
