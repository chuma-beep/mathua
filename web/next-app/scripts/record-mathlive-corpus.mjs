// Regenerate test/fixtures/mathlive-plain-text.json.
//
// The fixture is a *recording*, not a specification: it is what
// `MathfieldElement.getValue('plain-text')` actually returns in a real browser for
// the answer shapes Mathua's graders have to survive. This script is the only way
// that recording is produced, so the file can always be re-derived rather than
// hand-edited into agreement with itself.
//
// It runs the real library in Chromium. It cannot run in jsdom — MathLive measures
// text for real, and its `node` export condition resolves to the SSR build that
// never registers `<math-field>`, which is why plain-text comes back empty there.
//
//   node scripts/record-mathlive-corpus.mjs
//
// Run it after any MathLive version bump and read the diff before committing it:
// the corpus IS the contract, so a change here is a change to what the server will
// be asked to grade.

import { chromium } from '@playwright/test'
import { createServer } from 'node:http'
import { readFile, writeFile } from 'node:fs/promises'
import { extname, join, dirname } from 'node:path'
import { fileURLToPath } from 'node:url'
import { createRequire } from 'node:module'

const here = dirname(fileURLToPath(import.meta.url))
const appRoot = join(here, '..')
const mathliveDir = join(appRoot, 'node_modules', 'mathlive')
const require = createRequire(import.meta.url)
const { version } = require(join(mathliveDir, 'package.json'))

const PAGE = `<!doctype html><html><head><meta charset="utf-8"></head><body>
<script type="module">
import './mathlive.min.mjs'
window.__CASES = __CASES__
window.__probe = async () => {
  const mf = document.createElement('math-field')
  document.body.appendChild(mf)
  await new Promise(r => setTimeout(r, 60))
  const out = []
  for (const [gradingType, expected, latex] of window.__CASES) {
    mf.value = latex
    await new Promise(r => setTimeout(r, 5))
    const g = f => { try { return mf.getValue(f) } catch (e) { return 'ERR ' + e.message } }
    out.push({ gradingType, expected, latex, plain: g('plain-text'), back: mf.getValue('latex') })
  }
  return out
}
window.__ready = true
</script></body></html>`

// Every answer shape Mathua's existing graders must survive, grouped by the
// grading_type that will receive it. `numeric` rows are the answer forms a learner
// is expected to submit; `numeric-expression` rows are the ones they produce when
// they show their working, kept separate because the numeric grader parses literals
// only and cannot evaluate them (see internal/grader/mathlive_compat_test.go).
const CASES = [
  ['numeric', '-5', '-5'],
  ['numeric', '0.75', '0.75'],
  ['numeric', '7/2', '7\\div2'],
  ['numeric', '1/2', '\\frac{1}{2}'],
  ['numeric', '1/2', '\\frac{2}{4}'],
  ['numeric', '4 1/10', '4\\frac{1}{10}'],
  ['numeric', '-1 1/2', '-1\\frac{1}{2}'],
  ['numeric', '0.0314', '3.14e-2'],
  ['numeric-expression', '1000000', '10^{6}'],
  ['numeric-expression', '0.001', '10^{-3}'],
  ['numeric-expression', '0.0314', '3.14\\times 10^{-2}'],
  ['numeric-expression', '48', '12\\times4'],
  ['expression', '2x+3', '2x+3'],
  ['expression', 'x+4=10', 'x+4=10'],
  ['expression', '6x^2', '6x^{2}'],
  ['expression', 'x^2+2x+1', 'x^{2}+2x+1'],
  ['expression', '(x+1)^2', '(x+1)^{2}'],
  ['expression', '2x', '2\\pi r'],
  ['expression', '-3x-5', '-3x-5'],
  ['expression', 'sqrt(x^2)', '\\sqrt{x^{2}}'],
  ['expression', 'x^(n+1)', 'x^{n+1}'],
  ['expression', 'x^2', 'x^{2}'],
  ['symbolic', 'sqrt(2)', '\\sqrt{2}'],
  ['symbolic', '1/sqrt(3)', '\\frac{1}{\\sqrt{3}}'],
  ['symbolic', '(x+y)/(x-y)', '\\frac{x+y}{x-y}'],
  ['expression', 'sin(x)+cos(x)', '\\sin(x)+\\cos(x)'],
  ['expression', 'sin^2(x)+cos^2(x)', '\\sin^{2}(x)+\\cos^{2}(x)'],
  ['expression', '5 R 3', '5\\ R\\ 3'],
  ['expression', '3sqrt(2)', '3\\sqrt{2}'],
  ['expression', '2:3', '2:3'],
  ['tuple', '2,3', '2,3'],
  ['tuple', '(2,3)', '(2,3)'],
  ['ordering', '2,3,5', '2,3,5'],
  ['comparison', '<', '<'],
  ['comparison', '>', '>'],
  ['comparison', '<=', '\\le'],
  ['comparison', '>=', '\\ge'],
  // Multiple choice, which used to be the one grading type with no recorded cases at
  // all -- so `GradingMultipleChoice` sat in the compat test's pure-Go set and was
  // never actually exercised.
  //
  // These are the answer shapes the corpus really contains, taken from the most common
  // `expected` values across internal/generator: 673 "yes", 171 "no", then small
  // integers, then words. `choiceGrader` is `strings.EqualFold(expected, answer)`,
  // so every one of these is an exact string comparison and the risk is specific:
  // MathLive reads bare letters as implicit multiplication, so "yes" could come back
  // as "y e s" and "Z" as a single symbol, and neither would match. Measured rather
  // than assumed, which is the whole point of recording them.
  ['multiple_choice', 'yes', 'yes'],
  ['multiple_choice', 'no', 'no'],
  ['multiple_choice', '1/2', '\\frac{1}{2}'],
  ['multiple_choice', 'pi/2', 'pi/2'],
  ['multiple_choice', 'pi/2', '\\frac{\\pi}{2}'],
  ['multiple_choice', 'odd', 'odd'],
  ['multiple_choice', 'Z', 'Z'],
  ['multiple_choice', 'addition', 'addition'],
  ['multiple_choice', '2', '2'],
  ['multiple_choice', '-1', '-1'],
]

async function main() {
  const srv = createServer(async (req, res) => {
    const name = (req.url === '/' ? 'index.html' : req.url.split('?')[0]).replace(/^\//, '')
    if (name === 'mathlive.min.mjs') {
      res.writeHead(200, { 'Content-Type': 'text/javascript' })
      res.end(await readFile(join(mathliveDir, 'mathlive.min.mjs')))
      return
    }
    const types = { '.html': 'text/html', '.mjs': 'text/javascript' }
    res.writeHead(200, { 'Content-Type': types[extname(name)] ?? 'application/octet-stream' })
    res.end(PAGE.replace('__CASES__', JSON.stringify(CASES)))
  })
  if (process.env.DEBUG_REQUESTS) srv.on('request', r => console.log('  req', r.url))
  const port = 4400 + (process.pid % 200)
  await new Promise(r => srv.listen(port, r))

  const browser = await chromium.launch()
  const page = await browser.newPage()
  page.on('pageerror', e => console.error('  page error:', e.message))
  page.on('console', m => {
    // MathLive warns when it cannot fetch its webfonts from /fonts; that is expected
    // here and only pollutes the output.
    if (m.type() === 'error' && !m.text().includes('fonts')) console.error('  console:', m.text())
  })
  await page.goto(`http://localhost:${port}/`)
  await page.waitForFunction(() => window.__ready === true, { timeout: 30000 })
  const rows = await page.evaluate(() => window.__probe())
  await browser.close()
  srv.close()

  const out = {
    _comment:
      'Recorded by scripts/record-mathlive-corpus.mjs from mathlive in Chromium, via ' +
      "MathfieldElement.getValue('plain-text'). Not hand-written: do not edit, re-record.",
    mathliveVersion: version,
    cases: rows.map(r => ({ gradingType: r.gradingType, expected: r.expected, latex: r.latex, plain: r.plain })),
  }
  const dest = join(appRoot, 'test', 'fixtures', 'mathlive-plain-text.json')
  await writeFile(dest, JSON.stringify(out, null, 2) + '\n')
  console.log(`recorded ${out.cases.length} cases from mathlive ${version} → test/fixtures/mathlive-plain-text.json`)
  for (const r of out.cases) console.log(`  ${r.gradingType.padEnd(17)} ${r.latex.padEnd(24)} → ${JSON.stringify(r.plain)}`)
}

main().catch(err => {
  console.error('recording failed:', err)
  process.exit(1)
})