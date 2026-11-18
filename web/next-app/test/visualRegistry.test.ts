import { describe, it, expect } from 'vitest'
import fs from 'node:fs'
import path from 'node:path'
import { REPO_ROOT } from './helpers/roots'
import { VISUAL_TEMPLATES, VISUAL_TYPES, isVisualType } from '../components/visual/registry'

// The registry is the contract between three files that are edited by hand:
// data/visualizations/types.json (what the Go loader permits), registry.ts (what
// the browser can render) and the shards (what is actually authored). None of
// them imports another, so this test is the only place a drift is visible.

const allowlist = JSON.parse(
  fs.readFileSync(path.join(REPO_ROOT, 'data', 'visualizations', 'types.json'), 'utf8'),
) as { id: string; status: string }[]

function shardFiles(): string[] {
  const dir = path.join(REPO_ROOT, 'data', 'lessons', 'kp')
  return fs.readdirSync(dir).filter((name) => name.endsWith('.json')).map((name) => path.join(dir, name))
}

describe('visual registry', () => {
  it('matches the Go allowlist exactly', () => {
    expect(VISUAL_TYPES).toEqual(allowlist.map((entry) => entry.id))
  })

  it('every allowlisted id has a schema and a loader', () => {
    for (const entry of allowlist) {
      const template = VISUAL_TEMPLATES[entry.id]
      expect(template, `${entry.id} has no registry entry`).toBeDefined()
      expect(template.schema).toBeDefined()
      expect(typeof template.load).toBe('function')
    }
  })

  it('rejects a type that is not on the allowlist', () => {
    expect(isVisualType('linalg.matrix2d')).toBe(true)
    expect(isVisualType('linalg.not-a-template')).toBe(false)
  })

  it('every authored curriculum visual parses against its template schema', () => {
    let found = 0
    const problems: string[] = []
    for (const file of shardFiles()) {
      const steps = JSON.parse(fs.readFileSync(file, 'utf8')) as { label?: string; visual?: { type: string; params?: unknown } }[]
      steps.forEach((step, index) => {
        if (!step.visual) return
        found++
        const entry = VISUAL_TEMPLATES[step.visual.type]
        if (!entry) {
          problems.push(`${path.basename(file)} step ${index + 1}: unknown type ${step.visual.type}`)
          return
        }
        const parsed = entry.schema.safeParse(step.visual.params ?? {})
        if (!parsed.success) {
          problems.push(
            `${path.basename(file)} step ${index + 1} (${step.label ?? 'unlabelled'}): ${parsed.error.issues
              .map((issue) => issue.path.join('.') || issue.message)
              .join(', ')}`,
          )
        }
      })
    }
    expect(problems).toEqual([])
    expect(found).toBeGreaterThanOrEqual(6)
  })
})
