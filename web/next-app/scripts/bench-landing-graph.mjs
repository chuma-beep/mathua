// Landing-page concept graph benchmark (mobile emulation).
// Serves the static export and measures the real user path:
//   canvas ready, transferred graph bytes, idle/interaction frame times,
//   long tasks and JS heap.
//
// Usage:
//   node scripts/bench-landing-graph.mjs [url]
//   node scripts/bench-landing-graph.mjs --stress=5000
//   (default: starts `serve out` on :3999 itself)
import { spawn } from 'node:child_process'
import { chromium, devices } from 'playwright'
import { layoutDAG3D } from '../lib/layoutDAG3D.ts'

const PORT = 3999
const stressArg = process.argv.find(a => a.startsWith('--stress='))
const stress = stressArg ? Number(stressArg.split('=')[1]) : 0
const url = process.argv.find(a => a.startsWith('http')) || `http://localhost:${PORT}/`

function buildStress(n) {
  const domains = Array.from({ length: 17 }, (_, i) => `d${i}`)
  const concepts = []
  for (let i = 0; i < n; i++) {
    const prereqs = []
    if (i > 0) prereqs.push(`c${i - 1}`)
    if (i > 20) prereqs.push(`c${Math.floor((i - 20) * 0.7)}`)
    concepts.push({ id: `c${i}`, prerequisites: prereqs, domain: domains[i % domains.length] })
  }
  const { positions, importance } = layoutDAG3D(concepts, { domainOrder: domains })
  const index = new Map(concepts.map((c, i) => [c.id, i]))
  const nodes = concepts.map(c => [c.id, `Concept ${c.id.slice(1)}`, c.domain, importance[c.id] ?? 0])
  const edges = []
  for (const c of concepts) {
    for (const p of c.prerequisites) edges.push([index.get(p), index.get(c.id)])
  }
  const domainCounts = {}
  for (const c of concepts) domainCounts[c.domain] = (domainCounts[c.domain] ?? 0) + 1
  return {
    graph: { v: 2, nodes, edges, domains: domainCounts },
    positions,
  }
}

let server
if (!process.argv.find(a => a.startsWith('http'))) {
  server = spawn('npx', ['serve', 'out', '-l', String(PORT), '--no-clipboard'], {
    stdio: 'ignore',
    detached: false,
  })
  // npx cold start varies; poll until the static export answers.
  const deadline = Date.now() + 20_000
  for (;;) {
    try {
      const res = await fetch(url)
      if (res.ok) break
    } catch {
      /* not up yet */
    }
    if (Date.now() > deadline) throw new Error(`static server did not start on ${url}`)
    await new Promise(r => setTimeout(r, 250))
  }
}

const browser = await chromium.launch()
const ctx = await browser.newContext({ ...devices['Pixel 7'] })
const page = await ctx.newPage()

// Stress mode: serve a synthetic graph through the real fetch path so the
// renderer/LOD/labels are exercised at 1k-5k nodes without shipping fixtures.
if (stress > 0) {
  const { graph, positions } = buildStress(stress)
  await page.route('**/data/concepts.graph.json', route => route.fulfill({ json: graph }))
  await page.route('**/data/concepts.positions.json', route => route.fulfill({ json: positions }))
}

const resources = []
page.on('response', async res => {
  const u = res.url()
  if (!u.includes('/_next/') && !u.includes('/data/')) return
  try {
    const body = await res.body()
    resources.push({ url: u.replace(/^https?:\/\/[^/]+/, ''), status: res.status(), bytes: body.length })
  } catch {
    /* ignore */
  }
})
const consoleErrors = []
page.on('console', m => {
  if (m.type() === 'error') consoleErrors.push(m.text())
})
page.on('pageerror', e => consoleErrors.push(String(e)))

await page.addInitScript(() => {
  window.__lt = []
  if (typeof PerformanceObserver !== 'undefined') {
    try {
      new PerformanceObserver(list => {
        for (const e of list.getEntries()) window.__lt.push(Math.round(e.duration))
      }).observe({ entryTypes: ['longtask'] })
    } catch {
      /* unsupported */
    }
  }
})

const t0 = Date.now()
await page.goto(url, { waitUntil: 'load' })
await page.waitForSelector('canvas', { timeout: 45_000 })
const canvasReadyMs = Date.now() - t0

// Let the graph settle (positions + first frames).
await page.evaluate(() => document.querySelector('canvas')?.scrollIntoView({ block: 'center' }))
await page.waitForTimeout(3000)

async function frameStats(ms) {
  return page.evaluate(
    ms =>
      new Promise(resolve => {
        const deltas = []
        let last = performance.now()
        const stop = last + ms
        function loop(t) {
          deltas.push(t - last)
          last = t
          if (t < stop) requestAnimationFrame(loop)
          else resolve(deltas)
        }
        requestAnimationFrame(loop)
      }),
    ms
  )
}

function summarize(deltas) {
  const s = deltas.slice(1).sort((a, b) => a - b)
  if (s.length === 0) return { frames: 0 }
  const q = p => s[Math.min(s.length - 1, Math.floor(s.length * p))]
  return {
    frames: s.length,
    avg: +(s.reduce((a, b) => a + b, 0) / s.length).toFixed(2),
    p50: +q(0.5).toFixed(2),
    p95: +q(0.95).toFixed(2),
    max: +s[s.length - 1].toFixed(2),
  }
}

const idle = summarize(await frameStats(1500))

// Pan: touch scroll gesture across the canvas (CDP synthesizes real touch).
const cdp = await ctx.newCDPSession(page)
const box = await page.locator('canvas').boundingBox()
let pan = { frames: 0 }
let pinch = { frames: 0 }
if (box) {
  const cx = box.x + box.width / 2
  const cy = box.y + box.height / 2
  const measuring = frameStats(2500)
  await cdp.send('Input.synthesizeScrollGesture', {
    x: Math.round(cx),
    y: Math.round(cy),
    xDistance: -220,
    yDistance: 120,
    speed: 1200,
    gestureSourceType: 'touch',
  })
  pan = summarize(await measuring)

  const measuring2 = frameStats(2500)
  await cdp.send('Input.synthesizePinchGesture', {
    x: Math.round(cx),
    y: Math.round(cy),
    scaleFactor: 1.8,
    relativeSpeed: 600,
    gestureSourceType: 'touch',
  })
  pinch = summarize(await measuring2)
}

const info = await page.evaluate(() => {
  const nav = performance.getEntriesByType('navigation')[0]
  const paints = performance.getEntriesByType('paint').map(p => `${p.name}=${Math.round(p.startTime)}ms`)
  let renderer = 'unknown'
  try {
    const gl = document.createElement('canvas').getContext('webgl')
    const dbg = gl?.getExtension('WEBGL_debug_renderer_info')
    if (gl && dbg) renderer = String(gl.getParameter(dbg.UNMASKED_RENDERER_WEBGL))
  } catch {
    /* ignore */
  }
  return {
    domInteractive: Math.round(nav.domInteractive),
    load: Math.round(nav.loadEventEnd),
    paints,
    longTasks: (window.__lt || []).slice(0, 20),
    heapMB: performance.memory ? Math.round(performance.memory.usedJSHeapSize / 1048576) : null,
    dpr: window.devicePixelRatio,
    renderer,
  }
})

const byUrl = new Map()
for (const r of resources) {
  const prev = byUrl.get(r.url)
  if (!prev || r.bytes > prev.bytes) byUrl.set(r.url, r)
}
const top = [...byUrl.values()].sort((a, b) => b.bytes - a.bytes)

console.log(
  JSON.stringify(
    {
      url,
      stress,
      canvasReadyMs,
      idle,
      pan,
      pinch,
      info,
      transferredKB: Math.round(top.reduce((s, r) => s + r.bytes, 0) / 1024),
      topResources: top.slice(0, 12).map(r => ({ url: r.url, kb: Math.round(r.bytes / 1024) })),
      consoleErrors,
    },
    null,
    1
  )
)

await browser.close()
server?.kill()
