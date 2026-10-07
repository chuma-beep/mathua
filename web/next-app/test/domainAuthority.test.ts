import { describe, it, expect } from 'vitest'
import fs from 'node:fs'
import path from 'node:path'

// One authority for the domain list: lib/graphDomains.ts.
//
// Six modules used to carry their own copy. `DomainProgress` listed 15 domains and omitted
// precalculus and machine_learning, so 50 concepts appeared nowhere in the product — with no
// empty row and no error to say so. Three more modules omitted machine_learning. The dashboard's
// per-domain table was the worst of them because it was the one place a learner would look for
// "which subjects have I not touched".
//
// The guard is coverage, not prohibition. A module is allowed to declare its own domain list —
// an onboarding picker or a goals view may legitimately order domains differently, and
// countByDomain's own comment says presentation order is a view concern. What is not allowed is a
// list that silently omits a domain, because that is indistinguishable from a domain that does
// not exist.

// test/ -> next-app -> repo root. Getting this wrong makes every assertion below pass vacuously,
// because the scan then visits nothing and finds nothing — which is what the first version of
// test/noEmoji.test.ts did.
const REPO = path.resolve(__dirname, '..', '..', '..')
const ROOT = path.join(REPO, 'web', 'next-app')
const AUTHORED_DIRS = ['app', 'components', 'lib']
const CANONICAL = path.join(ROOT, 'lib', 'graphDomains.ts')

function corpusDomains(): string[] {
  const dir = path.join(REPO, 'data', 'concepts')
  const out = new Set<string>()
  for (const file of fs.readdirSync(dir)) {
    if (file === 'enrichment.json' || !file.endsWith('.json')) continue
    const data = JSON.parse(fs.readFileSync(path.join(dir, file), 'utf8'))
    for (const c of Array.isArray(data) ? data : Object.values(data)) {
      out.add((c as { domain: string }).domain)
    }
  }
  return [...out].sort()
}

/** Every source file we author, excluding the canonical module itself. */
function authoredFiles(): { rel: string; text: string }[] {
  const out: { rel: string; text: string }[] = []
  const walk = (dir: string) => {
    for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
      if (entry.name === 'node_modules' || entry.name.startsWith('.')) continue
      const abs = path.join(dir, entry.name)
      if (entry.isDirectory()) {
        walk(abs)
      } else if (/\.(ts|tsx)$/.test(entry.name) && abs !== CANONICAL) {
        out.push({ rel: path.relative(ROOT, abs), text: fs.readFileSync(abs, 'utf8') })
      }
    }
  }
  for (const d of AUTHORED_DIRS) walk(path.join(ROOT, d))
  return out
}

/**
 * Find declarations that enumerate the domain set.
 *
 * The discriminator is the declaration's own name or shape, not how many domains it mentions.
 * A count is not enough: `components/math/layouts.ts` names five domains in its geometry and
 * algebra layout tables, and those are subject-matter references — they exist to type equations,
 * not to enumerate the curriculum. Flagging them would bury the real offence under noise.
 *
 * So: a declaration counts when its name says "domain", or when it is an array literal of quoted
 * domain ids. Both are things only an enumeration does.
 */
function domainDeclarations(text: string, domains: string[]): { name: string; declared: string[] }[] {
  const found: { name: string; declared: string[] }[] = []
  for (const m of text.matchAll(/=\s*(\[|\{)/g)) {
    const window = text.slice(m.index, m.index + 2500)
    const declared = domains.filter(
      (d) => new RegExp(`['"]${d}['"]`).test(window) || new RegExp(`\\b${d}\\s*:`).test(window),
    )
    if (declared.length < 3) continue
    const line = text.slice(0, m.index).split('\n').pop() ?? ''
    const name = (line.match(/(?:const|export const)\s+(\w+)/) ?? [])[1] ?? ''
    const isArrayOfIds = m[1] === '[' && domains.some((d) => new RegExp(`['"]${d}['"]`).test(window))
    if (!/domain/i.test(name) && !isArrayOfIds) continue
    found.push({ name: name || 'array of domain ids', declared })
  }
  return found
}

describe('domain list authority', () => {
  const domains = corpusDomains()

  it('finds a corpus to check against, and it is the whole corpus', () => {
    // A guard with nothing to compare against reports nothing. The list is read from the corpus
    // files rather than hardcoded, so it cannot disagree with what the product ships.
    expect(domains.length).toBeGreaterThanOrEqual(17)
  })

  it('the canonical module covers every domain in the corpus', () => {
    const canon = fs.readFileSync(CANONICAL, 'utf8')
    const missing = domains.filter((d) => !new RegExp(`\\b${d}\\b`).test(canon))
    expect(missing, `lib/graphDomains.ts omits ${missing.join(', ')}`).toEqual([])
  })

  it('no other module declares a domain list that omits a domain', () => {
    const offenders: string[] = []
    for (const { rel, text } of authoredFiles()) {
      for (const { name, declared } of domainDeclarations(text, domains)) {
        const missing = domains.filter((d) => !declared.includes(d))
        if (missing.length > 0) {
          offenders.push(`${rel} ${name}: omits ${missing.join(', ')}`)
        }
      }
    }
    expect(
      offenders.length === 0 ? 'no partial domain list' : offenders.join('\n'),
    ).toBe('no partial domain list')
  })

  it('does not confuse a subject-matter reference with an enumeration', () => {
    // Geometry and algebra layout tables name five domains each. They exist to type equations,
    // not to enumerate the curriculum, and their names say nothing about domains — which is
    // exactly the distinction the detector uses. Without it this test file flags six files and
    // the real offence is invisible.
    const layouts = fs.readFileSync(path.join(ROOT, 'components', 'math', 'layouts.ts'), 'utf8')
    expect(domainDeclarations(layouts, domains)).toEqual([])
  })

  it('does still catch a named-but-partial domain map', () => {
    // The name test must not become a loophole: a partial list called DOMAIN_LABELS is the
    // original bug exactly.
    const domainsHere = corpusDomains()
    const partial = `const DOMAIN_LABELS = { ${domainsHere
      .slice(0, 5)
      .map((d) => `${d}: 1`)
      .join(', ')} }`
    const found = domainDeclarations(partial, domainsHere)
    expect(found.length).toBe(1)
    expect(found[0].name).toBe('DOMAIN_LABELS')
  })

  it('finds the real declarations in the codebase, so the guard is not inert', () => {
    // A guard that finds nothing anywhere passes vacuously. Require that the detector actually
    // sees at least the canonical module's own consumers.
    const seen = new Set<string>()
    for (const { text } of authoredFiles()) {
      for (const { name } of domainDeclarations(text, domains)) seen.add(name)
    }
    expect(seen.size, 'the detector found no domain enumeration anywhere').toBeGreaterThan(0)
  })
})
