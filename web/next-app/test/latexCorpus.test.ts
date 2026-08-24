import { describe, it } from 'vitest'
import { execSync } from 'node:child_process'
import fs from 'node:fs'
import path from 'node:path'
import katex from 'katex'
import { lessonMacros, prepareLessonMath } from '../lib/lessonMath'

// Corpus gate: every lesson body, exactly as the Go API serves it
// (canonicalized via cmd/latexdump), is pushed through the frontend math
// pipeline (prepareLessonMath) and every resulting math span is rendered by
// KaTeX with throwOnError.
//
// Errors are compared against a committed baseline
// (test/latexCorpusBaseline.json): a file may not GAIN a new error kind, and
// a file that becomes clean must be dropped from the baseline. Regenerate
// with UPDATE_BASELINE=1 npx vitest run test/latexCorpus.test.ts.
//
// Sibling gates: internal/latex/corpus_test.go (canonicalization idempotence
// + validator warnings) runs on the Go side.

const CORPUS_ROOT = path.join(__dirname, '../../../data/lessons')
const REPO_ROOT = path.join(CORPUS_ROOT, '../..')
const BASELINE_PATH = path.join(__dirname, 'latexCorpusBaseline.json')

const UPDATE_BASELINE = process.env.UPDATE_BASELINE === '1'

// Normalize a KaTeX error to its stable kind (positions and inputs churn).
function errorKind(msg: string): string {
  return msg.replace(/ at position \d+:[\s\S]*$/, '').trim()
}

function unescapeHtml(s: string): string {
  return s
    .replace(/&amp;/g, '&')
    .replace(/&lt;/g, '<')
    .replace(/&gt;/g, '>')
    .replace(/&quot;/g, '"')
    .replace(/&#39;/g, "'")
}

function canonicalBodies(): Record<string, string> {
  const stdout = execSync('go run ./cmd/latexdump', {
    cwd: REPO_ROOT,
    maxBuffer: 256 * 1024 * 1024,
  })
  return JSON.parse(stdout.toString())
}

describe('lesson corpus renders through KaTeX', () => {
  it(
    'renders every served math span without new KaTeX errors',
    { timeout: 300_000 },
    () => {
      const bodies = canonicalBodies()
      const files = Object.keys(bodies)
      if (files.length < 400) throw new Error(`corpus shrank unexpectedly: ${files.length} files`)

      // file -> set of error kinds
      const fresh: Record<string, Set<string>> = {}
      let spans = 0
      let errorCount = 0

      for (const rel of files) {
        const prepared = prepareLessonMath(bodies[rel])

        const spanRe = /<span class="math-(display|inline)">([\s\S]*?)<\/span>/g
        let m: RegExpExecArray | null
        while ((m = spanRe.exec(prepared)) !== null) {
          const [, kind, escapedInner] = m
          spans++
          const inner = unescapeHtml(escapedInner)
          try {
            katex.renderToString(inner, {
              displayMode: kind === 'display',
              throwOnError: true,
              trust: false,
              macros: lessonMacros,
            })
          } catch (err) {
            errorCount++
            const msg = err instanceof Error ? err.message : String(err)
            let set = fresh[rel]
          if (!set) { set = new Set(); fresh[rel] = set }
            set.add(errorKind(msg))
          }
        }
      }

      if (UPDATE_BASELINE) {
        const out: Record<string, string[]> = {}
        for (const rel of Object.keys(fresh).sort()) {
          out[rel] = Array.from(fresh[rel]).sort()
        }
        fs.writeFileSync(BASELINE_PATH, JSON.stringify(out, null, 2) + '\n')
        console.log(
          `baseline rewritten: ${Object.keys(out).length} files with errors, ${errorCount} broken spans`
        )
        return
      }

      let baseline: Record<string, string[]> = {}
      if (fs.existsSync(BASELINE_PATH)) {
        baseline = JSON.parse(fs.readFileSync(BASELINE_PATH, 'utf8'))
      } else {
        throw new Error('latexCorpusBaseline.json missing; run with UPDATE_BASELINE=1 first')
      }

      const newErrors: string[] = []
      for (const rel of Object.keys(fresh).sort()) {
        const allowed = new Set(baseline[rel] ?? [])
        if (allowed.size === 0) {
          for (const k of Array.from(fresh[rel])) newErrors.push(`${rel}: ${k} (file was clean)`)
        } else {
          for (const k of Array.from(fresh[rel])) {
            if (!allowed.has(k)) newErrors.push(`${rel}: ${k} (new kind)`)
          }
        }
      }
      for (const rel of Object.keys(baseline).sort()) {
        if (!fresh[rel]) {
          throw new Error(
            `${rel} no longer has KaTeX errors — remove it from latexCorpusBaseline.json (UPDATE_BASELINE=1)`
          )
        }
      }

      if (newErrors.length > 0) {
        throw new Error(
          `${errorCount} broken spans; ${newErrors.length} NEW error kinds:\n` + newErrors.slice(0, 20).join('\n')
        )
      }
    }
  )
})
