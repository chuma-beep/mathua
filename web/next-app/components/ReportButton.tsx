'use client'

import { useState, useEffect } from 'react'
import { submitReport, type ReportKind, type ReportReason } from '../lib/api'

interface ReportButtonProps {
  conceptId?: string
  lessonId?: string
  kind: ReportKind
  question?: string
  expected?: string
  explanation?: string
  source?: string
  sessionId?: string
  attemptId?: string
  /** What the learner has typed, and the format the field gave them, so a
      moderator can rebuild the attempt instead of guessing at it. */
  answerGiven?: string
  gradingType?: string
  answerHint?: string
  answerInputMode?: string
  /** Button caption. Defaults to 'Report a problem'; pass a specific one
      whenever several report buttons share a row or card. */
  label?: string
}

const REASONS: { value: ReportReason; label: string }[] = [
  { value: 'wrong_answer', label: 'Wrong answer' },
  { value: 'bad_explanation', label: 'Bad explanation' },
  { value: 'unclear', label: 'Unclear / confusing' },
  { value: 'formatting', label: 'Formatting / math rendering' },
  { value: 'other', label: 'Other' },
]

export default function ReportButton(props: ReportButtonProps) {
  const [open, setOpen] = useState(false)
  const [reason, setReason] = useState<ReportReason>('wrong_answer')
  const [detail, setDetail] = useState('')
  const [sending, setSending] = useState(false)
  const [sent, setSent] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    if (sent) {
      setSent(false)
    }
  }, [props.attemptId, props.question])

  useEffect(() => {
    if (!sent) return
    const timer = setTimeout(() => {
      setSent(false)
    }, 4000)
    return () => clearTimeout(timer)
  }, [sent])

  async function handleSubmit() {
    if (sending || sent) return
    setSending(true)
    setError('')
    try {
      await submitReport({
        concept_id: props.conceptId,
        lesson_id: props.lessonId,
        kind: props.kind,
        question: props.question,
        expected: props.expected,
        explanation: props.explanation,
        source: props.source,
        session_id: props.sessionId,
        attempt_id: props.attemptId,
        reason,
        detail: detail.trim(),
        answer_given: props.answerGiven,
        grading_type: props.gradingType,
        answer_hint: props.answerHint,
        answer_input_mode: props.answerInputMode,
      })
      setSent(true)
      setOpen(false)
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Could not send report')
    } finally {
      setSending(false)
    }
  }

  if (sent) {
    return (
      <span className="font-mono text-[10px] text-mathua-green uppercase tracking-wider">
        ✓ Thanks — reported
      </span>
    )
  }

  if (!open) {
    return (
      <button
        type="button"
        onClick={() => setOpen(true)}
        title="Complain about this question or explanation"
        className="font-mono text-[10px] text-mathua-muted hover:text-mathua-blue transition-colors uppercase tracking-wider"
      >
        {props.label ?? 'Report a problem'}
      </button>
    )
  }

  // A fixed sheet, not an inline block. Inline, the panel moved whenever the
  // verdict above it appeared or the question changed, which dismissed the
  // native picker on Android about half a second after it opened. The sheet
  // sits above the keyboard and its position cannot be shifted by the form.
  return (
    <div
      role="dialog"
      aria-modal="true"
      aria-label="Report a problem"
      className="fixed inset-0 z-[3000] flex items-end sm:items-center justify-center bg-black/60 p-0 sm:p-4"
    >
    <div className="w-full sm:max-w-md max-h-[85vh] overflow-y-auto border border-mathua-border bg-mathua-bg p-3 pb-[calc(0.75rem+env(safe-area-inset-bottom))]">
      <div className="flex items-center justify-between mb-2">
        <span className="font-mono text-[10px] text-mathua-secondary uppercase tracking-wider">
          What&apos;s wrong?
        </span>
        <button
          type="button"
          onClick={() => setOpen(false)}
          className="font-mono text-[10px] text-mathua-muted hover:text-mathua-primary"
        >
          Cancel
        </button>
      </div>
      <select
        value={reason}
        onChange={e => setReason(e.target.value as ReportReason)}
        className="w-full bg-mathua-surface border border-mathua-border px-2 h-9 text-xs font-mono text-mathua-primary outline-none"
      >
        {REASONS.map(r => (
          <option key={r.value} value={r.value}>{r.label}</option>
        ))}
      </select>
      <textarea
        autoFocus
        value={detail}
        onChange={e => setDetail(e.target.value)}
        placeholder="Optional details (what did you expect?)"
        rows={2}
        maxLength={2000}
        className="mt-2 w-full bg-mathua-surface border border-mathua-border px-2 py-2 text-xs font-mono text-mathua-primary placeholder:text-mathua-muted outline-none resize-y"
      />
      {error && (
        <p className="mt-1 font-mono text-[10px] text-red-400">{error}</p>
      )}
      <button
        type="button"
        onClick={handleSubmit}
        disabled={sending}
        className="mt-2 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint transition-colors px-3 h-9 text-xs font-mono disabled:opacity-50"
      >
        {sending ? 'Sending…' : 'Send report'}
      </button>
    </div>
    </div>
  )
}
