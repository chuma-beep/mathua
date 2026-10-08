'use client'

import { useEffect, useState } from 'react'
import KatexContent from './KatexContent'
import LessonDiagram from './LessonDiagram'
import ReportMenu from './ReportMenu'
import { getLessonBody, getLessonKPs, getLessons, type LessonKpsRes } from '../lib/api'
import { getUserInfo } from '../lib/auth'
import { buildExampleBlocks, type ExampleBlock } from '../lib/lessonBlocks'
import { stripMathDelimiters } from '../lib/lessonMath'
import { concepts } from '../lib/conceptData'

// The reference layer, scoped to one concept and reachable only from inside the learning flow.
//
// This used to be `/study`: a browsable library of whole lessons, with its own search, its own
// domain index, and a link to it in every header, tab bar, sidebar and footer. That was a second
// way to approach the same mathematics, and the product had decided against it — "Learn teaches,
// Study explains" only holds if the explaining is reachable without leaving the loop. So the
// route is gone and the material stayed.
//
// What is deliberately *not* here: a way to browse. There is no domain list, no lesson list, no
// search box, no sibling-concept chips, and no link that lands a learner on a concept they did
// not come here for. Everything below is about the one concept being learned. That is the whole
// difference between reference material and a second learning path, and it is behavioral rather
// than cosmetic: the material can explain, and it cannot hand the learner a curriculum.
//
// Reading is also not evidence. There is no "mark as read", no completion, no "I studied this"
// — nothing below moves a mastery status, and the note at the foot of the panel says so
// explicitly, because a panel of authoritative-looking prose sitting next to the questions is
// exactly where a learner would expect one.

interface Prereq {
  id: string
  label: string
  status: string
}

interface Props {
  conceptId: string
  /**
   * Prerequisites from the parent's readiness fetch.
   *
   * Passed in rather than re-fetched so the reference panel cannot become a second place that
   * knows what a prerequisite is. `/learn`'s banner offers the first three with links into the
   * right flow; this lists all of them, as names and status, so nothing is quietly dropped.
   */
  prerequisites?: Prereq[]
}

function statusWord(status: string | undefined): string {
  switch ((status ?? '').toUpperCase()) {
    case 'MASTERED':
      return 'learned'
    case 'DECAYING':
      return 'learned, due a check'
    case 'LEARNING':
    case 'PRACTICING':
      return 'in progress'
    default:
      return 'not started'
  }
}

// Only reached for a concept whose shard carries no knowledge points at all. It costs one lesson
// index fetch, so it stays off the common path — and every concept in the corpus has a shard, so
// in practice this never runs. It is here so a concept authored without one still gets an
// explanation rather than an empty panel.
async function fetchLessonTitle(conceptId: string): Promise<string | null> {
  const res = await getLessons(getUserInfo()?.student_id)
  for (const lessons of Object.values(res.lessons)) {
    for (const l of lessons) {
      if (l.concepts.includes(conceptId)) return l.title
    }
  }
  return null
}

export default function ConceptReference({ conceptId, prerequisites = [] }: Props) {
  const [shard, setShard] = useState<LessonKpsRes | null>(null)
  const [lessonTitle, setLessonTitle] = useState<string | null>(null)
  const [body, setBody] = useState<string>('')
  const [failed, setFailed] = useState(false)

  useEffect(() => {
    let cancelled = false
    setShard(null)
    setLessonTitle(null)
    setBody('')
    setFailed(false)
    getLessonKPs(conceptId)
      .then(res => { if (!cancelled) setShard(res) })
      .catch(() => { if (!cancelled) setFailed(true) })
    return () => { cancelled = true }
  }, [conceptId])

  const blocks: ExampleBlock[] = shard ? buildExampleBlocks(conceptId, shard) : []
  const needsBody = shard !== null && blocks.length === 0

  useEffect(() => {
    if (!needsBody) return
    let cancelled = false
    fetchLessonTitle(conceptId)
      .then(async (title) => {
        if (cancelled || !title) return
        const text = await getLessonBody(title)
        if (cancelled) return
        setLessonTitle(title)
        setBody(text)
      })
      .catch(() => { if (!cancelled) setFailed(true) })
    return () => { cancelled = true }
  }, [needsBody, conceptId])

  const label = concepts.find(c => c.id === conceptId)?.label || conceptId

  return (
    <details className="mb-4 border border-mathua-border bg-mathua-surface">
      <summary className="font-mono text-[11px] uppercase tracking-wider text-mathua-muted px-4 py-3 cursor-pointer hover:text-mathua-blue">
        Reference — the full lesson
      </summary>

      <div className="px-4 pb-4 min-w-0 overflow-hidden">
        {prerequisites.length > 0 && (
          <div className="mb-4">
            <h3 className="font-mono text-[11px] text-mathua-muted mb-2 uppercase tracking-wider">
              This depends on ({prerequisites.length})
            </h3>
            <ul className="space-y-1">
              {prerequisites.map(p => (
                <li key={p.id} className="font-mono text-[11px] text-mathua-secondary">
                  {p.label} <span className="text-mathua-muted">· {statusWord(p.status)}</span>
                </li>
              ))}
            </ul>
            <p className="mt-2 font-mono text-[10px] text-mathua-muted">
              A prerequisite that needs work is offered above, in the learning flow.
            </p>
          </div>
        )}

        {failed && (
          <p role="alert" className="font-mono text-xs text-mathua-red">
            Couldn&apos;t load the reference for {label}. The questions below still work.
          </p>
        )}

        {blocks.length > 0 && (
          <div data-testid="reference-blocks">
            <h3 className="font-mono text-[11px] text-mathua-muted mb-2 uppercase tracking-wider">
              Worked examples ({blocks.length})
            </h3>
            <div className="space-y-3">
              {blocks.map((b, i) => (
                <div key={b.key} className="border border-mathua-border p-4 min-w-0 overflow-hidden">
                  <div className="flex items-start gap-2 min-w-0">
                    <div className="flex-1 min-w-0">
                      {b.body ? (
                        <details open={i === 0}>
                          <summary className="font-mono text-xs text-mathua-primary cursor-pointer">
                            {i + 1}. {b.label ? stripMathDelimiters(b.label) : 'Worked example'}
                            {b.also.length > 0 && (
                              <span className="text-mathua-muted"> · also {b.also.map(stripMathDelimiters).join(', ')}</span>
                            )}
                          </summary>
                          {b.diagram && (
                            <div className="mt-3">
                              <LessonDiagram src={b.diagram} alt={`Worked diagram for ${label}`} />
                            </div>
                          )}
                          {b.subgoals.length > 0 && (
                            <ul className="mt-2 space-y-1">
                              {b.subgoals.map((sg, j) => (
                                <li key={j} className="font-mono text-[11px] text-mathua-secondary pl-3 relative before:content-['–'] before:absolute before:left-0">
                                  {stripMathDelimiters(sg)}
                                </li>
                              ))}
                            </ul>
                          )}
                          <div className="mt-2 bg-mathua-code border border-mathua-border p-3 text-sm overflow-hidden">
                            <KatexContent>{b.body}</KatexContent>
                          </div>
                        </details>
                      ) : (
                        <>
                          <p className="font-mono text-xs text-mathua-primary">
                            {i + 1}. {b.label ? stripMathDelimiters(b.label) : 'Worked example'}
                            {b.also.length > 0 && (
                              <span className="text-mathua-muted"> · also {b.also.map(stripMathDelimiters).join(', ')}</span>
                            )}
                          </p>
                          {b.subgoals.length > 0 && (
                            <ul className="mt-2 space-y-1">
                              {b.subgoals.map((sg, j) => (
                                <li key={j} className="font-mono text-[11px] text-mathua-secondary pl-3 relative before:content-['–'] before:absolute before:left-0">
                                  {stripMathDelimiters(sg)}
                                </li>
                              ))}
                            </ul>
                          )}
                        </>
                      )}
                    </div>
                    <ReportMenu
                      conceptId={conceptId}
                      lessonId={lessonTitle ?? undefined}
                      kind={b.diagram ? 'diagram' : 'worked_example'}
                      question={b.diagram ?? `${b.label}: ${b.body}`.slice(0, 2000)}
                      blockId={`${conceptId}/${b.index}`}
                    />
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {needsBody && body && (
          <div>
            <h3 className="font-mono text-[11px] text-mathua-muted mb-2 uppercase tracking-wider">
              Lesson text{lessonTitle ? ` — ${lessonTitle}` : ''}
            </h3>
            <div className="border border-mathua-border bg-mathua-code p-3 text-sm overflow-hidden">
              <KatexContent>{body}</KatexContent>
            </div>
            <div className="mt-2 flex justify-end">
              <ReportMenu
                conceptId={conceptId}
                lessonId={lessonTitle ?? undefined}
                kind="lesson_body"
                question={body.slice(0, 2000)}
                blockId={`${conceptId}/body`}
              />
            </div>
          </div>
        )}

        <p className="mt-4 border-t border-mathua-border pt-3 font-mono text-[10px] text-mathua-muted">
          Reference explains. It does not count as progress — only the questions below move this
          concept.
        </p>
      </div>
    </details>
  )
}