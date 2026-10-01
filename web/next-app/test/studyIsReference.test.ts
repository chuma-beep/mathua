import { describe, it, expect } from 'vitest'
import { readFileSync } from 'node:fs'
import { join } from 'node:path'

const APP = join(__dirname, '..')
const components = readFileSync(join(APP, 'app', 'study', 'components.tsx'), 'utf-8')

// "Study explains. Learn teaches." A reference page that renders a learner's
// mastery as a percentage and a progress bar invites them to do something on it,
// which is the confusion the split exists to prevent. ADR-021 is specific:
// percentages leave Study, while the per-concept status words and badges inside
// LessonDetail stay, because those are what tell a reader where to begin.
describe('Study is reference, not a progress surface', () => {
  const domainOverview = components.slice(
    components.indexOf('export function DomainOverview'),
    components.indexOf('// ── Domain Drill-Down'),
  )
  const drillDown = components.slice(
    components.indexOf('export function DomainDrillDown'),
    components.indexOf('export function LessonDetail'),
  )
  const lessonDetail = components.slice(components.indexOf('export function LessonDetail'))

  // Asserting on specific expressions was not enough: reintroducing the
  // percentage as domainAgg[domain]?.pct matched none of them and the suite
  // stayed green. These views must not mention the aggregate or any mastery
  // quantity at all, however it is spelled.
  const progression = /domainAgg|agg\b|mastered|\bpct\b|allMastered|allDone|h-\[3px\]/g

  it('renders nothing progression-derived on the domain list', () => {
    // The aggregate is still computed in page.tsx, but only to decide whether
    // to offer the "Start here" orientation.
    expect(domainOverview.match(progression)).toBeNull()
  })

  it('still shows how much is in a domain — that is corpus, not learner', () => {
    expect(domainOverview).toMatch(/lessons\.length} lesson/)
  })

  it('renders nothing progression-derived in the drill-down', () => {
    expect(drillDown).not.toMatch(/Domain progress/)
    expect(drillDown.match(progression)).toBeNull()
  })

  it('shows concept count in the drill-down but not how many are mastered', () => {
    expect(drillDown).toMatch(/concept\{total !== 1 \? 's' : ''\}/)
  })

  it('keeps the per-concept status words and badges in the lesson detail', () => {
    // These are the part that actually helps: "before you start" says which
    // prerequisites you already hold, in words, with no percentage.
    expect(lessonDetail).toMatch(/not started/)
    expect(lessonDetail).toMatch(/in progress/)
    expect(lessonDetail).toMatch(/mastered/)
    expect(lessonDetail).toMatch(/MasteryBadge/)
  })
})
