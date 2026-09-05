'use client'

interface ProgressBarProps {
  answered: number
  estimatedTotal: number
  className?: string
}

// Shared progress bar for the Diagnostic (adaptive CAT, 15-45 questions).
// Determinate when the backend reports an estimate; indeterminate shimmer
// fallback while the estimate is unknown (estimatedTotal <= 0).
export default function ProgressBar({ answered, estimatedTotal, className = '' }: ProgressBarProps) {
  const indeterminate = !estimatedTotal || estimatedTotal <= 0
  const safeAnswered = Math.max(0, answered)
  const pct = indeterminate ? 0 : Math.min((safeAnswered / estimatedTotal) * 100, 100)
  const displayAnswered = indeterminate ? safeAnswered : Math.min(safeAnswered, estimatedTotal)

  return (
    <div className={`mb-4 ${className}`}>
      <div className="flex justify-between text-[10px] font-mono text-mathua-muted mb-1">
        <span>Progress bar</span>
        <span aria-live="polite">
          {indeterminate ? `Question ${Math.max(safeAnswered, 1)} · finding your frontier…` : `${displayAnswered} / ${estimatedTotal}`}
        </span>
      </div>
      {indeterminate ? (
        <div
          role="progressbar"
          aria-label="Progress bar"
          aria-busy="true"
          className="h-1.5 bg-mathua-code rounded-full overflow-hidden relative"
        >
          <div className="progress-bar-indeterminate h-full bg-mathua-blue rounded-full" />
          <style>{`@keyframes progress-bar-slide { 0% { transform: translateX(-100%); } 100% { transform: translateX(300%); } } .progress-bar-indeterminate { width: 33%; animation: progress-bar-slide 1.4s ease-in-out infinite; } @media (prefers-reduced-motion: reduce) { .progress-bar-indeterminate { animation: none; width: 33%; } }`}</style>
        </div>
      ) : (
        <div
          role="progressbar"
          aria-label="Progress bar"
          aria-valuemin={0}
          aria-valuemax={estimatedTotal}
          aria-valuenow={displayAnswered}
          aria-valuetext={`Question ${displayAnswered} of about ${estimatedTotal}`}
          className="h-1.5 bg-mathua-code rounded-full overflow-hidden"
        >
          <div
            className="h-full bg-mathua-blue rounded-full transition-all duration-500"
            style={{ width: `${pct}%` }}
          />
        </div>
      )}
    </div>
  )
}
