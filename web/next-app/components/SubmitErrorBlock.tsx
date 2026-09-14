'use client'

import { getErrorMessage, getErrorStatus } from '../lib/api'

export interface SubmitError {
  message: string
  status?: number
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
  const sessionGone = error.status === 404
  const skipsExhausted = onSkip != null && (skipsLeft ?? MAX_SKIPS) <= 0
  return (
    <div
      role="alert"
      className="border border-mathua-red/60 bg-mathua-surface p-4 text-left"
    >
      <p className="font-mono text-xs text-mathua-red">
        {sessionGone
          ? 'Session expired — the server restarted and forgot this run.'
          : "Couldn't submit your answer — it's still in the box."}
      </p>
      <p className="mt-1 font-mono text-[11px] text-mathua-muted break-words [overflow-wrap:anywhere]">
        {error.status != null ? `(${error.status}) ` : ''}{error.message}
      </p>
      <div className="mt-3 flex flex-wrap items-center gap-2">
        {!sessionGone && (
          <button
            type="button"
            onClick={onRetry}
            disabled={retrying}
            className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white px-4 h-10 font-mono text-xs disabled:opacity-50"
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
          No skips left — restart for a fresh run.
        </p>
      )}
    </div>
  )
}
