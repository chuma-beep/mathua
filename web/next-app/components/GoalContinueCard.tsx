'use client'

import Link from 'next/link'
import { useEffect, useState } from 'react'
import { getGoal, getGraph, getProgress } from '../lib/api'
import { getUserInfo } from '../lib/auth'
import { goalProgress, nextGoalStep, type GoalConcept } from '../lib/goalScope'

interface GoalCardState {
  goal: GoalConcept
  learned: number
  total: number
  pct: number
  step: GoalConcept | null
}

/**
 * A small bridge between the graph's learning goal and the daily study loop.
 *
 * The goal already trims the graph, but without this card the Learn shelf still
 * looks like the goal does not exist. It is presentation only: progress and the
 * next step come from the same goalScope helpers the graph uses.
 */
export default function GoalContinueCard() {
  const [state, setState] = useState<GoalCardState | null>(null)

  useEffect(() => {
    let cancelled = false
    const user = getUserInfo()
    if (!user) return

    Promise.all([getGoal(), getGraph(), getProgress(user.student_id)])
      .then(([goalId, graph, progress]) => {
        if (cancelled || !goalId) return
        const concepts: GoalConcept[] = graph.nodes.map(node => ({
          id: node.id,
          label: node.label,
          domain: node.domain,
          prerequisites: node.prerequisites,
        }))
        const goal = concepts.find(concept => concept.id === goalId)
        if (!goal) return

        const summary = goalProgress(goalId, concepts, progress)
        const stepId = nextGoalStep(goalId, concepts, progress)
        const step = stepId ? concepts.find(concept => concept.id === stepId) ?? null : null
        setState({ goal, ...summary, step })
      })
      .catch(() => {
        if (!cancelled) setState(null)
      })

    return () => {
      cancelled = true
    }
  }, [])

  if (!state) return null

  return (
    <section aria-label="Your goal" className="mb-4 border border-mathua-border bg-mathua-surface p-4">
      <div className="flex flex-wrap items-start justify-between gap-3">
        <div className="min-w-0">
          <p className="font-mono text-[10px] uppercase tracking-wider text-mathua-muted">Your goal</p>
          <h2 className="mt-1 font-mono text-sm text-mathua-primary break-words">{state.goal.label}</h2>
        </div>
        <Link
          href={`/graph?goal=${encodeURIComponent(state.goal.id)}`}
          className="shrink-0 font-mono text-[10px] uppercase tracking-wider text-mathua-blue hover:underline"
        >
          View path
        </Link>
      </div>

      <div className="mt-3 h-1.5 w-full bg-mathua-border" aria-hidden="true">
        <div className="h-full bg-mathua-blue" style={{ width: `${state.pct}%` }} />
      </div>
      <p className="mt-2 font-mono text-[10px] uppercase tracking-wider text-mathua-muted">
        {state.learned} of {state.total} concepts learned
      </p>

      {state.step ? (
        <Link
          href={`/learn?concept=${encodeURIComponent(state.step.id)}`}
          className="mt-3 inline-flex items-center gap-1 font-mono text-xs uppercase tracking-wider text-mathua-blue hover:text-mathua-blue-hover"
        >
          Continue toward {state.step.label} &#8594;
        </Link>
      ) : (
        <p className="mt-3 text-sm text-mathua-secondary">
          No ready step yet. Review a prerequisite to unlock the next one.
        </p>
      )}
    </section>
  )
}
