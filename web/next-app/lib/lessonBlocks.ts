import type { KpInfo, LessonKpsRes } from './api'
import { stripMathDelimiters } from './lessonMath'

// Worked-example blocks for one concept, deduplicated.
//
// This is the content half of what used to be `app/study/components.tsx`. It moved here with
// the reference layer rather than with the route: the route is closed, the material is not, and
// a module that lives under a deleted page directory is a module nobody can find.
//
// The dedupe rule is the interesting part and is preserved exactly. A KP shard may name the same
// lesson section twice under byte-different headings, so identical bodies recur. Rendering one
// of them and pointing at the other is not an option — there is no second place to point, and
// the learner would read it twice.

/** Normalize a section name the way Go's `normSectionKey` does. Mirrored by `norm_section` in `scripts/audit_lessons.py`. */
export function normSection(s: string): string {
  return s.replace(/\\[()[\]]|\$\$?/g, '').split(/\s+/).join(' ').trim()
}

/** Short stable hash for body text (djb2), paired with the normalized section to form the dedupe key. */
export function hashBody(s: string): string {
  let h = 5381
  for (let i = 0; i < s.length; i++) h = ((h << 5) + h + s.charCodeAt(i)) >>> 0
  return h.toString(36)
}

export interface ExampleBlock {
  key: string
  label: string
  /** Other labels that resolved to this same body, so a rename does not hide content. */
  also: string[]
  subgoals: string[]
  body: string
  conceptId: string
  index: number
  diagram: string | null
}

/**
 * Turn one concept's shard into its blocks.
 *
 * The key is (normalized section + body hash) rather than either half alone: two distinct
 * sections can legitimately share a body, and one section can legitimately hold two bodies, so
 * a key built from only one of them hides content in one direction or the other.
 */
export function buildExampleBlocks(conceptId: string, res: LessonKpsRes): ExampleBlock[] {
  const kps: KpInfo[] = res.kps ?? []
  const conceptDiagram = res.diagram ?? null
  const blocks: ExampleBlock[] = []
  const seen = new Map<string, ExampleBlock>()
  const seenDiagrams = new Set<string>()

  kps.forEach((kp, k) => {
    // The shard is an external payload, so every field is defaulted here rather than at the
    // point of use. A KP with no label reaches `stripMathDelimiters` on the way to the
    // summary text, and that function calls `.replace` on what it is given — so a missing
    // label took down the whole page rather than rendering an unlabelled example.
    const body = kp.worked_example ?? ''
    const label = kp.label ?? ''
    const key = `${normSection(kp.section ?? '')}::${hashBody(body)}`

    if (body) {
      const first = seen.get(key)
      if (first) {
        const norm = (l: string) => stripMathDelimiters(l).toLowerCase()
        if (label && ![first.label, ...first.also].some(l => norm(l) === norm(label))) {
          first.also.push(label)
        }
        return
      }
    }

    const block: ExampleBlock = {
      key: `lesson-kp-${conceptId}-${k}`,
      label,
      also: [],
      subgoals: (kp.subgoals ?? []).filter((s): s is string => typeof s === 'string'),
      body,
      conceptId,
      index: k,
      diagram: null,
    }
    if (body) seen.set(key, block)
    // The concept's own diagram rides on its first block only. Attaching it to every block
    // repeats the same figure three times on one screen.
    if (conceptDiagram && !seenDiagrams.has(conceptDiagram) && k === 0) {
      seenDiagrams.add(conceptDiagram)
      block.diagram = conceptDiagram
    }
    blocks.push(block)
  })

  return blocks
}