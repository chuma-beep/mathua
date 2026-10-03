#!/usr/bin/env node
// Run every math span in the lesson corpus through Compute Engine and report the
// spans that are genuinely malformed.
//
// Why this exists: the corpus gate (test/latexCorpus.test.ts) proves KaTeX can
// *render* the lessons. It says nothing about whether the expressions are
// well-formed, because KaTeX is forgiving by design — it renders a malformed
// fraction by guessing. Something stricter is needed to find the spans that are
// subtly wrong, and Compute Engine is the strict thing already in the tree.
//
// Four decisions, each of which could have gone the other way:
//
// 1. The spans come from cmd/latexdump, not from a regex over the Markdown.
//    `latexnorm.Scan` is the only thing in the tree that knows a dollar can be
//    a delimiter or a currency amount, and it runs over the *canonicalized*
//    body — the text the API serves. A regex here would report every "$5 and
//    up" in the corpus as a syntax error, which is a worse failure than no
//    report at all: it would train everyone to ignore the output.
//
// 2. Compute Engine's `parse` is total: it does not throw on a LaTeX syntax
//    error, it returns a tree containing an `["Error", reason, …]` node. A
//    try/catch therefore reports a corpus with thousands of broken spans as
//    perfectly clean, which is worse than no report because it is a confident
//    lie. Failure is detected from the returned box, not from an exception.
//
// 3. A parse failure is not automatically a defect in the corpus, and treating
//    it as one is the trap here. CE is a numeric/symbolic CAS, and this corpus is
//    largely abstract algebra, group theory and topology. It reports
//    `U \times F` as `incompatible-type: number, function` — having parsed the
//    syntax perfectly and then declined to type a product of two sets. It does
//    the same for `\mathbb{R}^2 \setminus \{0\}` and for `\xrightarrow{d}`, which
//    it simply has no command for. Roughly 1,500 of the 46,872 spans land in
//    that bucket. Filing them as "unparseable" would bury the ~130 real defects
//    under noise and get this report switched off.
//
//    So the reasons are split by *what CE objected to*:
//
//      syntax   closing delimiters, stray operators, stray delimiters — the
//               LaTeX is wrong. These are the report.
//      not-numeric  CE understood it and refused to type it. Valid math, wrong
//               tool. Counted, not reported.
//      ce-gap   a command CE does not implement, or an operand it wanted and
//               could not find. Cannot be distinguished from valid-but-unsupported
//               without a human, so it is counted, not reported.
//
// 4. Every span gets an outcome and the buckets reconcile to the corpus total.
//    Nothing is skipped quietly, and the summary states what fraction of the
//    corpus was actually checkable — a report that quietly covers 44% of its
//    input is not a clean bill of health.
//
// Output is JSON Lines, one object per line, shaped like Go's slog JSON
// handler: time, level, msg, then flat attributes. Readable with `jq`, and
// unmarshalable in Go with no adapter.
//
// Usage:
//   node scripts/normalize-latex.mjs              # report, exit 1 on defects
//   node scripts/normalize-latex.mjs --quiet      # summary + defects only
//   node scripts/normalize-latex.mjs data/lessons/authored
//
// Exit codes: 0 clean, 1 at least one malformed span, 2 the harness itself broke
// (Go build failed, bad JSON, unreconciled counts) — deliberately distinct from
// 1, so a broken report is never mistaken for a clean corpus.

import { ComputeEngine } from '@cortex-js/compute-engine'
import { execFileSync } from 'node:child_process'
import { pathToFileURL } from 'node:url'

const args = process.argv.slice(2)
const VALUED = new Set(['--limit'])
const QUIET = args.includes('--quiet')
/** The first bare word that is not a flag or a flag's value. */
const ROOT =
  args.find((a, i) => !a.startsWith('--') && !VALUED.has(args[i - 1])) ?? 'data/lessons'

function die(msg, code = 2) {
  process.stderr.write(`${msg}\n`)
  process.exit(code)
}

/** One JSONL record, slog-shaped. */
function log(level, msg, attrs = {}) {
  process.stdout.write(JSON.stringify({ time: new Date().toISOString(), level, msg, ...attrs }) + '\n')
}

/**
 * LaTeX that is valid but is not an expression CE could ever be asked to type.
 * Pre-parse, so the common cases cost nothing.
 */
const NON_EXPRESSION = [
  {
    re: /\\begin\{(align|aligned|equation|gather|multline|cases|array|matrix|pmatrix|bmatrix)\b/,
    why: 'environment',
  },
  { re: /\\label\{|\\tag\{|\\ref\{|\\eqref\{|\\cite\{/, why: 'cross-reference' },
  { re: /\\text\{|\\mbox\{|\\emph\{/, why: 'prose inside math' },
  { re: /\\begin\{(itemize|enumerate)\}/, why: 'list' },
  { re: /\\includegraphics|\\begin\{figure\}/, why: 'figure' },
  { re: /\\[xu]rightarrow|\\longrightarrow|\\implies|\\iff\b/, why: 'diagram arrow' },
  { re: /\\setminus|\\cup\b|\\cap\b|\\subseteq|\\in\b/, why: 'set membership' },
]

/** Classify a span without parsing it. Returns {outcome, why}. */
export function classify(latex) {
  for (const { re, why } of NON_EXPRESSION) {
    if (re.test(latex)) return { outcome: 'unsupported', why }
  }
  return { outcome: 'parse' }
}

/**
 * Bucket a CE parse result.
 *
 * Returns null when the span is clean, else `{bucket, reason}` where bucket is
 * one of `syntax` (a real defect), `not-numeric` or `ce-gap`.
 *
 * `missing` is deliberately filed under ce-gap. It is genuinely ambiguous — it
 * fires on both `2+` and on `d: \Omega^k \to \Omega^{k+1}` — and when a signal
 * cannot tell a defect from a valid span, the conservative home is the bucket
 * that does not accuse the corpus.
 */
export function classifyResult(box) {
  const hit = findError(box)
  if (!hit) return null
  const reason = hit.reason
  const short = reason.slice(0, 120)

  // CE understood the syntax and then refused the *semantics*: it wanted a number
  // and got a set, a morphism or a formal symbol. Valid math, wrong tool.
  if (/incompatible-type|non-evaluable|expect a/.test(reason)) {
    return { bucket: 'not-numeric', reason: short }
  }
  // A command, function or symbol CE has no rule for. Cannot be told apart from
  // valid-but-unsupported without a human, so it is counted, not reported.
  if (/^(unexpected-command|missing|invalid-base|expected-function|invalid-symbol|unexpected-argument)/.test(reason)) {
    return { bucket: 'ce-gap', reason: short }
  }
  // CE's own point-algebra objections, which are not defects in the LaTeX.
  if (/^no-(product-between-points|division-by-point)/.test(reason)) {
    return { bucket: 'ce-gap', reason: short }
  }
  // Anything left is a genuine structural fault: an unclosed delimiter, a stray
  // operator, a token that has no place in the expression.
  return { bucket: 'syntax', reason: short }
}

/**
 * The first error node in a parsed expression, or null.
 *
 * `box.type === 'error'` catches a span that failed wholesale. The recursive walk
 * is not redundant: a failed subexpression can leave a parent that still parses,
 * and `2 + \frac{3}{` is exactly that shape.
 *
 * The reason CE gives ('expected-closing-delimiter', 'unexpected-operator', …) is
 * far more useful than "failed", so it is lifted out of the tree.
 */
export function findError(box) {
  if (!box) return { reason: 'parsed to nothing' }
  const walk = (node, depth) => {
    if (depth > 64 || node === null || typeof node !== 'object') return null
    if (Array.isArray(node)) {
      if (node[0] === 'Error') return { reason: normalizeReason(node[1]) }
      for (const child of node) {
        const hit = walk(child, depth + 1)
        if (hit) return hit
      }
      return null
    }
    for (const key of Object.keys(node)) {
      const hit = walk(node[key], depth + 1)
      if (hit) return hit
    }
    return null
  }
  // Walk first and box.type second, because the walk yields the *whole* reason and
  // box.type does not. Reading the reason off the box instead truncates
  // "ErrorCode,'incompatible-type',…" to "ErrorCode", which then fails to match
  // every type-error pattern and files 773 valid spans as malformed.
  const found = walk(box.json, 0)
  if (found) return found
  return box.type === 'error' ? { reason: 'unknown error' } : null
}

/**
 * CE's reasons arrive as `"'unexpected-command'"` or
 * `"ErrorCode,'incompatible-type','number',…"` — quoted, and sometimes prefixed
 * with a serializer tag. Left as-is, every anchored pattern in classifyResult
 * silently fails to match and the bucket lands in `malformed`, which is how 1,563
 * valid spans first showed up as defects. Normalized once, here.
 */
export function normalizeReason(raw) {
  let s = String(raw).trim()
  if (s.startsWith('ErrorCode,')) s = s.slice('ErrorCode,'.length)
  // The apostrophes are CE's own token delimiters inside a reason code. Dropping
  // them all is both simpler and more readable than unwrapping only a whole-string
  // pair, which misses `'a','b'` and leaves a leading quote that defeats every
  // anchored pattern downstream.
  return s.replace(/'/g, '')
}

/**
 * The snippet shown in a report. A failing span is usually long, and the useful
 * part is where the parser gave up, so this keeps the head and the tail rather
 * than the first N characters.
 */
export function snippet(s, max = 120) {
  const flat = String(s).replace(/\s+/g, ' ').trim()
  if (flat.length <= max) return flat
  const head = Math.ceil((max - 1) * 0.7)
  return flat.slice(0, head) + '…' + flat.slice(-(max - 1 - head))
}

/** scripts/ lives in web/next-app/, so the Go module is three levels up. */
const REPO = new URL('../../../', import.meta.url).pathname

function isMain(url) {
  return process.argv[1] !== undefined && url === pathToFileURL(process.argv[1]).href
}

function loadSpans(root) {
  let raw
  try {
    raw = execFileSync('go', ['run', './cmd/latexdump', '-spans', root], {
      cwd: REPO,
      encoding: 'utf8',
      maxBuffer: 512 * 1024 * 1024,
    })
  } catch (e) {
    die(`latexdump failed (is the Go toolchain up?): ${e.stderr || e.message}`)
  }
  const spans = []
  for (const [n, line] of raw.split('\n').entries()) {
    if (!line.trim()) continue
    try {
      spans.push(JSON.parse(line))
    } catch {
      // A malformed line is a harness bug, not a corpus problem, and skipping it
      // would quietly understate the corpus.
      die(`latexdump line ${n + 1} is not JSON: ${snippet(line, 200)}`)
    }
  }
  return spans
}

// ---- run ----
// Guarded so a test can import the classifiers above without running the scan.
if (isMain(import.meta.url)) main()

function main() {
const spans = loadSpans(ROOT)

// Deduplicate by content before parsing. The corpus repeats itself heavily — the
// same few thousand expressions across 585 files — and a defect's real signal is
// which files contain it, not how many times it recurs.
const unique = new Map()
for (const s of spans) {
  let e = unique.get(s.content)
  if (!e) {
    e = { content: s.content, kind: s.kind, count: 0, files: new Set() }
    unique.set(s.content, e)
  }
  e.count++
  e.files.add(s.file)
}

const ce = new ComputeEngine()
const syntax = []
const notNumeric = new Map()
const ceGap = new Map()
const skipped = new Map()
let ok = 0
let crashed = 0

// CE logs a line to stderr for every canonicalization error it recovers from.
// That is the same information the tree already carries, and at ~1,500 errors it
// drowns the report it is supposed to be annotating.
const realError = console.error
console.error = () => {}

for (const entry of unique.values()) {
  const pre = classify(entry.content)
  if (pre.outcome === 'unsupported') {
    skipped.set(pre.why, (skipped.get(pre.why) ?? 0) + entry.count)
    continue
  }
  let box
  try {
    box = ce.parse(entry.content)
  } catch {
    // CE is not perfectly total: a few inputs overflow its own recursion limit
    // during canonicalization. That is a fault in the checker, not a finding
    // about the corpus, so it is counted separately and never reported.
    crashed += entry.count
    continue
  }
  const verdict = classifyResult(box)
  if (!verdict) {
    ok += entry.count
  } else if (verdict.bucket === 'syntax') {
    syntax.push({ entry, error: verdict.reason })
  } else if (verdict.bucket === 'not-numeric') {
    notNumeric.set(verdict.reason, (notNumeric.get(verdict.reason) ?? 0) + entry.count)
  } else {
    ceGap.set(verdict.reason, (ceGap.get(verdict.reason) ?? 0) + entry.count)
  }
}
console.error = realError

const sum = (m) => [...m.values()].reduce((a, b) => a + b, 0)
const skippedSpans = sum(skipped)
const notNumericSpans = sum(notNumeric)
const ceGapSpans = sum(ceGap)
const syntaxSpans = syntax.reduce((n, f) => n + f.entry.count, 0)

// Every count is in *span occurrences*, not unique expressions, so the buckets can
// be reconciled against the corpus. Getting this wrong is not a rounding error:
// the first version of this script mixed the two and reported 20,549 of 46,872,
// which the check below is what caught.
const accounted = ok + skippedSpans + notNumericSpans + ceGapSpans + syntaxSpans + crashed
if (accounted !== spans.length) {
  die(`internal error: accounted for ${accounted} of ${spans.length} spans. The checker is wrong, not the corpus.`)
}

if (!QUIET) {
  for (const { entry, error } of syntax) {
    log('WARN', 'malformed span', {
      kind: entry.kind,
      snippet: snippet(entry.content),
      error,
      occurrences: entry.count,
      file_count: entry.files.size,
      files: [...entry.files].slice(0, 8),
    })
  }
  for (const [why, n] of [...skipped].sort((a, b) => b[1] - a[1])) {
    log('INFO', 'span is not an expression', { why, spans: n })
  }
  for (const [reason, n] of [...notNumeric].sort((a, b) => b[1] - a[1]).slice(0, 15)) {
    log('INFO', 'valid but outside CE numeric type system', { reason, spans: n })
  }
  for (const [reason, n] of [...ceGap].sort((a, b) => b[1] - a[1]).slice(0, 15)) {
    log('INFO', 'command or operand CE does not implement', { reason, spans: n })
  }
}

// The honest summary: what was checked, and what was not.
const checkable = ok + syntaxSpans
log('INFO', 'normalize-latex summary', {
  files: new Set(spans.map(s => s.file)).size,
  spans: spans.length,
  unique_expressions: unique.size,
  checked: checkable,
  checked_pct: spans.length === 0 ? 0 : Math.round((100 * checkable) / spans.length),
  clean: ok,
  malformed: syntaxSpans,
  malformed_unique: syntax.length,
  unchecked_not_expression: skippedSpans,
  unchecked_not_numeric: notNumericSpans,
  unchecked_ce_gap: ceGapSpans,
  engine_crashes: crashed,
})

if (spans.length > 0 && checkable === 0) {
  log('WARN', 'nothing was checkable — the corpus and the checker have diverged')
}

process.exit(syntax.length > 0 ? 1 : 0)
}
