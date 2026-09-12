// Generate the 1200x630 Open Graph card from the precomputed graph.
// Run manually after `npm run graph:build`; the PNG is committed so the
// Docker/CI build never needs a browser. Usage: npm run og:build
import { readFileSync, writeFileSync } from 'node:fs'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'
import { chromium } from 'playwright'

const root = join(dirname(fileURLToPath(import.meta.url)), '..')
const concepts = JSON.parse(readFileSync(join(root, 'data', 'concepts.json'), 'utf8'))
const positions = JSON.parse(readFileSync(join(root, 'public', 'data', 'concepts.positions.json'), 'utf8'))

const DOMAIN_COLORS = {
  arithmetic: '#4db8a0', fractions: '#a8a0f0', prealgebra: '#7dd3fc', algebra: '#e8a849',
  geometry: '#86efac', trigonometry: '#fda4af', complex_numbers: '#e879f9', precalculus: '#fda4af',
  calculus: '#e879f9', linear_algebra: '#86efac', statistics: '#67e8f9', machine_learning: '#818cf8',
  discrete_math: '#f0abab', number_theory: '#a8e6cf', differential_equations: '#fdba74',
  abstract_algebra: '#c4b5fd', topology: '#f9a8d4',
}

// Fixed tilt so the 3D cloud reads as depth in a still image.
const ANGLE = 0.42
const cos = Math.cos(ANGLE)
const sin = Math.sin(ANGLE)
const points = concepts
  .filter(c => positions[c.id])
  .map(c => {
    const [x, y, z] = positions[c.id]
    return { x: x * cos - z * sin, y, domain: c.domain }
  })

const xs = points.map(p => p.x)
const ys = points.map(p => p.y)
const minX = Math.min(...xs)
const maxX = Math.max(...xs)
const minY = Math.min(...ys)
const maxY = Math.max(...ys)
const W = 560
const H = 560
const PAD = 30
const px = p => PAD + ((p.x - minX) / (maxX - minX)) * (W - PAD * 2)
const py = p => PAD + (1 - (p.y - minY) / (maxY - minY)) * (H - PAD * 2)

const dots = points
  .map(p => {
    const color = DOMAIN_COLORS[p.domain] ?? '#5a6577'
    return `<circle cx="${px(p).toFixed(1)}" cy="${py(p).toFixed(1)}" r="4.2" fill="${color}" opacity="0.9"/>`
  })
  .join('')

const domainCount = new Set(concepts.map(c => c.domain)).size
const html = `<!doctype html>
<html><head><meta charset="utf-8"/>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link href="https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:wght@400&family=IBM+Plex+Serif:wght@400&display=block" rel="stylesheet">
<style>
  * { margin: 0; padding: 0; box-sizing: border-box; }
  body {
    width: 1200px; height: 630px; overflow: hidden;
    background: #0b0f1a; color: #e8e2d5;
    font-family: 'IBM Plex Serif', Georgia, serif;
    background-image:
      linear-gradient(rgba(42,63,95,0.35) 1px, transparent 1px),
      linear-gradient(90deg, rgba(42,63,95,0.35) 1px, transparent 1px);
    background-size: 48px 48px;
  }
  .frame { position: absolute; inset: 18px; border: 1px solid #1e2d45; }
  .copy { position: absolute; left: 72px; top: 172px; width: 560px; }
  .brand { font-family: 'IBM Plex Mono', monospace; font-size: 22px; color: #c8a96e; letter-spacing: 0.08em; }
  h1 { font-size: 62px; font-weight: 400; line-height: 1.12; margin-top: 18px; }
  p { font-size: 21px; color: #9fa8b4; margin-top: 22px; line-height: 1.6; }
  .meta { font-family: 'IBM Plex Mono', monospace; font-size: 15px; color: #5a6577; margin-top: 34px; letter-spacing: 0.04em; }
  .graph { position: absolute; right: 44px; top: 35px; width: 560px; height: 560px; opacity: 0.92; }
</style></head>
<body>
  <div class="frame"></div>
  <div class="copy">
    <div class="brand">λ Mathua</div>
    <h1>Adaptive math<br/>learning platform</h1>
    <p>Mastery-gated. Every concept unlocks only when its prerequisites are proven — fast and correct.</p>
    <div class="meta">${concepts.length} concepts · ${domainCount} domains · MIT · open source</div>
  </div>
  <svg class="graph" viewBox="0 0 ${W} ${H}">${dots}</svg>
</body></html>`

const browser = await chromium.launch()
const page = await browser.newPage({ viewport: { width: 1200, height: 630 }, deviceScaleFactor: 1 })
await page.setContent(html, { waitUntil: 'networkidle' })
await page.waitForTimeout(600)
const png = await page.screenshot({ type: 'png' })
await browser.close()

const dest = join(root, 'public', 'og.png')
writeFileSync(dest, png)
console.log(`og:build ${concepts.length} nodes / ${domainCount} domains -> ${dest} (${Math.round(png.length / 1024)}KB)`)
