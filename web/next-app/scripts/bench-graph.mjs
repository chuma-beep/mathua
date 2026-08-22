// One-off production load benchmark for /graph (not committed to test suite).
// Serves nothing itself: expects `npm run start:static` on :3001 already running.
import { chromium } from 'playwright'

const url = process.argv[2] || 'http://localhost:3001/graph'

const browser = await chromium.launch()
const page = await browser.newPage()
const t0 = Date.now()
await page.goto(url, { waitUntil: 'commit' })
await page.waitForSelector('.react-flow__node', { timeout: 30_000 })
const nodesVisible = Date.now() - t0
const count = await page.locator('.react-flow__node').count()
console.log(`graph interactive: ${nodesVisible}ms, nodes in DOM: ${count}`)
await browser.close()
