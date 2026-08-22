import { chromium } from 'playwright'
const browser = await chromium.launch()
const page = await browser.newPage()
await page.goto('http://localhost:4177/graph', { waitUntil: 'networkidle' })
const data = await page.evaluate(() => {
  const nav = performance.getEntriesByType('navigation')[0]
  const res = performance.getEntriesByType('resource')
    .filter(r => r.name.includes('_next'))
    .map(r => ({ n: r.name.split('/').pop().slice(0, 40), ms: Math.round(r.responseEnd - r.startTime), kb: Math.round(r.transferSize / 1024), start: Math.round(r.startTime) }))
    .sort((a, b) => b.ms - a.ms)
    .slice(0, 8)
  const paint = performance.getEntriesByType('paint').map(p => `${p.name}=${Math.round(p.startTime)}ms`)
  return { domInteractive: Math.round(nav.domInteractive), loadEvent: Math.round(nav.loadEventEnd), paint, top: res }
})
console.log(JSON.stringify(data, null, 1))
await browser.close()
