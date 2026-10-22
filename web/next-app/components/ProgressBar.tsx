'use client'

interface ProgressBarProps {
  answered: number
  budget: number
  className?: string
}

// Shared progress bar for the Diagnostic. The budget is the probe count of this
// run, not the cover size, so the bar tracks questions answered.
export default function ProgressBar({ answered, budget, className = '' }: ProgressBarProps) {
  const safeAnswered = Math.max(0, answered)
  const total = Math.max(0, budget)
  const done = Math.min(safeAnswered, total)
  const pct = total > 0 ? (done / total) * 100 : 0
  const indeterminate = total <= 0
  const questionLabel = `Question ${Math.max(safeAnswered + 1, 1)}`

  return (
    <div className={`mb-4 ${className}`}>
      <div className="flex justify-between text-[10px] font-mono text-mathua-muted mb-1">
        <span>Progress bar</span>
        <span aria-live="polite">
          {questionLabel}
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
          aria-valuemax={total}
          aria-valuenow={done}
          aria-valuetext={questionLabel}
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
