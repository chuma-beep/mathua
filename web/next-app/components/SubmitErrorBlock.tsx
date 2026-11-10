'use client'

import { useEffect, useRef } from 'react'
import { getErrorMessage, getErrorStatus } from '../lib/api'
import { classifyFailure, failureHeadline } from '../lib/failure'

export interface SubmitError {
  message: string
  status?: number
  /** A caller-specific headline when the generic failure bucket would mislead. */
  headline?: string
}

// oxlint-disable-next-line anti-slop/no-unknown-parameters -- catch value: unknown is the correct error contract
export function toSubmitError(err: unknown, fallback: string): SubmitError {  const message = getErrorMessage(err)
  return { message: message || fallback, status: getErrorStatus(err) }
}

export const MAX_SKIPS = 3

// Persistent, in-flow submit failure. Replaces native alert(): the message
// carries the server's own words (status + detail) so the next occurrence
// is diagnosable on sight, and the actions offer a way forward.
export default function SubmitErrorBlock({
  error,
  onRetry,
  onSkip,
  skipsLeft,
  onRestart,
  restartLabel,
  retrying,
}: {
  error: SubmitError
  onRetry: () => void
  onSkip?: () => void
  skipsLeft?: number
  onRestart: () => void
  restartLabel: string
  retrying?: boolean
}) {
  const failure = classifyFailure(error)
  const sessionGone = failure.kind === 'not_found'
  const throttled = failure.kind === 'throttled'
  const skipsExhausted = onSkip != null && (skipsLeft ?? MAX_SKIPS) <= 0

  // A 429 means the answer was refused for pace, not for being wrong, so retry
  // it once automatically rather than making the learner press a button for a
  // condition the server clears in a couple of seconds. The ref stops a loop if
  // the retry is refused again.
  const autos = useRef<string | null>(null)
  useEffect(() => {
    if (!throttled || retrying) return
    const key = `${error.status}:${error.message}`
    if (autos.current === key) return
    autos.current = key
    const timer = setTimeout(onRetry, 2000)
    return () => clearTimeout(timer)
  }, [throttled, retrying, error.status, error.message, onRetry])
  return (
    <div
      role="alert"
      className="border border-mathua-red bg-mathua-surface p-4 text-left"
    >
      <p className="font-mono text-xs text-mathua-red">
        {sessionGone
          ? 'Session expired, the server restarted and forgot this run.'
          : error.headline ?? failureHeadline(failure)}
      </p>
      <p className="mt-1 font-mono text-[11px] text-mathua-muted break-words [overflow-wrap:anywhere]">
        {error.status != null ? `(${error.status}) ` : ''}{error.message}
      </p>
      <div className="mt-3 flex flex-wrap items-center gap-2">
        {!sessionGone && failure.retryable && (
          <button
            type="button"
            onClick={onRetry}
            disabled={retrying}
            className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-4 h-10 font-mono text-xs disabled:opacity-50"
          >
            {retrying ? 'Retrying…' : 'Retry'}
          </button>
        )}
        {!sessionGone && onSkip && !skipsExhausted && (
          <button
            type="button"
            onClick={onSkip}
            className="border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue px-4 h-10 font-mono text-xs"
          >
            Skip →
          </button>
        )}
        <button
          type="button"
          onClick={onRestart}
          className="font-mono text-[11px] text-mathua-muted hover:text-mathua-primary underline underline-offset-2"
        >
          {restartLabel}
        </button>
      </div>
      {!sessionGone && skipsExhausted && (
        <p className="mt-2 font-mono text-[10px] text-mathua-muted">
          No skips left, restart for a fresh run.
        </p>
      )}
    </div>
  )
}
