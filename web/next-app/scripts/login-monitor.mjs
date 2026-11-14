import { chromium } from '@playwright/test'

// Production login monitor. The Fly health check stayed green through an
// eight-hour login outage because it never exercised the browser bundle. This
// script checks what a learner actually does: render the form, submit a wrong
// password, and click the Google button. When Turnstile is enabled the form
// stops at the verification prompt, so the script checks that prompt and then
// proves the API still rejects a tokenless login.
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
    const username = `login_monitor_${stamp}`
    const password = `wrong-password-${stamp}`
    await page.locator('#username').fill(username)
    await page.locator('#password').fill(password)
    await page.locator('[data-testid="auth-submit"]').click()

    const invalid = page.getByText(/Invalid credentials/i).first()
    const turnstilePrompt = page.getByText(/Complete the verification/i).first()
    const outcome = await Promise.race([
      invalid.waitFor({ state: 'visible', timeout: 20000 }).then(() => 'invalid'),
      turnstilePrompt.waitFor({ state: 'visible', timeout: 20000 }).then(() => 'turnstile'),
    ]).catch(() => fail('neither Invalid credentials nor the verification prompt appeared'))

    let apiStatus = 0
    if (outcome === 'turnstile') {
      const result = await page.evaluate(async ({ username, password }) => {
        const res = await fetch('/api/auth/login', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ username, password }),
        })
        return { status: res.status, body: await res.text() }
      }, { username, password })
      apiStatus = result.status
      if (result.status !== 403 || !result.body.includes('captcha verification failed')) {
        fail(`tokenless login returned ${result.status}: ${result.body}`)
      }
    }

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
      invalidCredentials: outcome === 'invalid',
      turnstilePrompt: outcome === 'turnstile',
      apiStatus: apiStatus || undefined,
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
