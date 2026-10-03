'use client'

// The optional hint under the answer field.
//
// Deliberately narrow. It shows the learner what *their own* expression simplifies to.
// It never compares against the expected answer, because the client does not have it —
// the server holds that, and shipping it here to power a hint would quietly turn the
// hint into an answer oracle and undo ADR-005. There is therefore no way for this
// component to influence a verdict, and `onChange` below feeds nothing but the field's
// own display.
//
// Compute Engine is 788 kB gzipped — nearly four times MathLive. It is never part of any
// route's initial JS; it is fetched the first time a learner pauses mid-answer with
// arithmetic on screen, and never fetched at all otherwise.
//
// That is still a large thing to pull down because someone typed a sum, so on a
// metered or slow connection the hint is skipped entirely. A one-line simplification is
// not worth someone's data allowance, and a feature that silently costs money is worse
// than one that is visibly absent.

import { useEffect, useRef, useState } from 'react'

/** Below this the engine is not worth loading or showing anything for. */
const MIN_LENGTH = 3
/** Long enough that the hint has settled before a keystroke interrupts it. */
const DEBOUNCE_MS = 700

/**
 * Whether pulling the engine is acceptable on this connection.
 *
 * `saveData` is the user telling the browser they are on a metered plan. The
 * effective-type check covers the slow-2G/3G cases, where an 788 kB fetch is the
 * difference between a usable page and a stalled one.
 */
export function isFetchReasonable(): boolean {
  if (typeof navigator === 'undefined') return false
  const c = (navigator as Navigator & { connection?: { saveData?: boolean; effectiveType?: string } }).connection
  if (!c) return true
  if (c.saveData) return false
  return c.effectiveType !== 'slow-2g' && c.effectiveType !== '2g'
}

interface Props {
  /** The learner's LaTeX, straight from the field. */
  latex: string
  /** Rendered only when the learner has an answer worth checking. */
  active: boolean
}

export default function SelfCheck({ latex, active }: Props) {
  const [hint, setHint] = useState<string | null>(null)
  const timer = useRef<ReturnType<typeof setTimeout> | null>(null)
  // The engine import is shared module-wide, so a second field reuses the first.
  const pending = useRef(0)

  useEffect(() => {
    if (!active || latex.trim().length < MIN_LENGTH) {
      setHint(null)
      return
    }
    if (timer.current) clearTimeout(timer.current)

    // Debounced, and only for arithmetic: CE will read prose as a product of
    // letter-variables and call it valid, so anything without an operator, a digit,
    // a power or a radical is not sent to it at all.
    if (!/[\d^{}]|[+\-*/=]|\\frac|\\sqrt/u.test(latex)) {
      setHint(null)
      return
    }
    // …and not at all where the download would be the wrong trade.
    if (!isFetchReasonable()) {
      setHint(null)
      return
    }

    const token = ++pending.current
    timer.current = setTimeout(async () => {
      const { loadEngine, simplificationHint } = await import('./equivalence')
      const engine = await loadEngine()
      if (!engine || token !== pending.current) return
      setHint(simplificationHint(engine, latex)?.text ?? null)
    }, DEBOUNCE_MS)

    return () => {
      if (timer.current) clearTimeout(timer.current)
    }
  }, [latex, active])

  if (!hint) return null
  return (
    // Not a verdict, and not styled like one: no green, no cross, nothing that
    // could be mistaken for the grader having spoken.
    <p data-testid="math-self-check" className="mt-1.5 font-mono text-[11px] text-mathua-muted">
      {hint}
    </p>
  )
}
