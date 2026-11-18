'use client'

import { createContext, useContext, useId, useState, type ReactNode } from 'react'
import type { PredictSpec } from '../../lib/api'

// Predict before you reveal.
//
// The proposal's teaching sequence has six steps and the fifth is "ask learners
// to predict the outcome before revealing it". A figure that responds to a
// control teaches a learner what the control does; a figure that asks first
// teaches them whether their model of the idea was right, which is the part that
// changes what they believe.
//
// The prediction lives on the visual, not on a template, so a figure added later
// inherits it without the dispatcher, the frame or eleven templates each growing
// a prop. The context is provided by `VisualFigure` and consumed by
// `FigureFrame`, which is the only place that renders chrome.
//
// Answering is local and is never sent anywhere. A prediction is not evidence of
// mastery, so it cannot be graded by the server, cannot grant XP and cannot move
// a review date.

const PredictContext = createContext<PredictSpec | null>(null)

export function PredictProvider({ predict, children }: { predict?: PredictSpec | null; children: ReactNode }) {
  return <PredictContext.Provider value={predict ?? null}>{children}</PredictContext.Provider>
}

export function usePredict(): PredictSpec | null {
  return useContext(PredictContext)
}

/**
 * The question and its options. Rendered inside `FigureFrame` above the drawing,
 * with native radios so arrow keys move between the answers.
 *
 * Feedback is revealed as soon as an option is chosen and stays available while
 * the learner drags the figure, which is the point: the prediction is a claim
 * they can now test. Wrong answers are never punished, hidden or reset.
 */
export function PredictBlock({ predict }: { predict: PredictSpec }) {
  const group = useId()
  const [choice, setChoice] = useState<number | null>(null)
  const answered = choice !== null
  const correct = choice === predict.answer

  return (
    <fieldset className="mb-3 border border-mathua-border px-3 py-2">
      <legend className="px-1 font-mono text-[10px] tracking-wide text-mathua-muted">Predict first</legend>
      <p className="font-mono text-[11px] leading-relaxed text-mathua-primary">{predict.question}</p>
      <div className="mt-2 space-y-1">
        {predict.options.map((option, index) => {
          const id = `${group}-${index}`
          return (
            <label key={id} htmlFor={id} className="flex cursor-pointer items-start gap-2 font-mono text-[11px] text-mathua-secondary">
              <input
                id={id}
                type="radio"
                name={group}
                checked={choice === index}
                onChange={() => setChoice(index)}
                className="mt-[2px]"
                style={{ accentColor: 'var(--accent-blue)' }}
              />
              <span>{option}</span>
            </label>
          )
        })}
      </div>
      <div role="status" aria-live="polite" className="mt-2 min-h-[1.25rem]">
        {answered && (
          <p className={`font-mono text-[11px] ${correct ? 'text-mathua-green' : 'text-mathua-red'}`}>
            {correct ? 'Yes. ' : 'Not quite. '}
            <span className="text-mathua-secondary">{predict.explain}</span>
          </p>
        )}
      </div>
    </fieldset>
  )
}
