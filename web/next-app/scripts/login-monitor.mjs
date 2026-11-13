import { chromium } from '@playwright/test'

// Production login monitor. The Fly health check stayed green through an
// eight-hour login outage because it never exercised the browser bundle. This
// script checks what a learner actually does: render the form, submit a wrong
// password, and click the Google button.
const BASE = process.env.MATHUA_BASE_URL || 'https://mathua.com'
const EXPECTED_CALLBACK = `${BASE}/api/auth/google/callback`

function fail(message) {
  throw new Error(message)
}

async function main() {
  const configResponse = await fetch(`${BASE}/api/config`)
  if (!configResponse.ok) fail(`config request failed: ${configResponse.status}`)
  const config = await configResponse.json()
  if (!config.auth_enabled) fail('auth is disabled in production')
  if (!config.google_client_id) fail('google_client_id is missing from /api/config')

  const browser = await chromium.launch()
  const page = await browser.newPage({ viewport: { width: 1280, height: 800 } })

  const consoleErrors = []
  const failedRequests = []
  page.on('console', message => {
    if (message.type() === 'error') consoleErrors.push(message.text())
  })
  page.on('requestfailed', request => {
    failedRequests.push(`${request.url()} (${request.failure()?.errorText ?? 'failed'})`)
  })

  try {
    await page.goto(`${BASE}/login`, { waitUntil: 'domcontentloaded', timeout: 45000 })
    await page.locator('#username').waitFor({ state: 'visible', timeout: 20000 })
    await page.locator('#password').waitFor({ state: 'visible', timeout: 10000 })
    await page.locator('[data-testid="auth-submit"]').waitFor({ state: 'visible', timeout: 10000 })
    const google = page.getByRole('button', { name: /Continue with Google/i }).first()
    await google.waitFor({ state: 'visible', timeout: 20000 })

    const stamp = Date.now()
    await page.locator('#username').fill(`login_monitor_${stamp}`)
    await page.locator('#password').fill(`wrong-password-${stamp}`)
    await page.locator('[data-testid="auth-submit"]').click()
    await page.getByText(/Invalid credentials/i).first().waitFor({ state: 'visible', timeout: 20000 })

    await google.click()
    await page.waitForURL(/accounts\.google\.com/, { timeout: 25000 })
    const url = new URL(page.url())
    const clientId = url.searchParams.get('client_id')
    const redirectUri = url.searchParams.get('redirect_uri')
    if (clientId !== config.google_client_id) {
      fail(`google client_id mismatch: ${clientId ?? 'missing'}`)
    }
    if (redirectUri !== EXPECTED_CALLBACK) {
      fail(`redirect_uri mismatch: ${redirectUri ?? 'missing'}`)
    }

    console.log(JSON.stringify({
      ok: true,
      base: BASE,
      invalidCredentials: true,
      googleReached: true,
      clientId,
      redirectUri,
      consoleErrors,
      failedRequests,
    }))
  } finally {
    await browser.close()
  }
}

main().catch(error => {
  console.error(JSON.stringify({ ok: false, base: BASE, error: String(error) }))
  process.exitCode = 1
})
