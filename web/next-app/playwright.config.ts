import { defineConfig, devices } from '@playwright/test'

const PORT = 4173
const baseURL = `http://localhost:${PORT}`

// The diagnosis suite is a measurement harness, not a regression test: a scenario that throws is
// printed and swallowed, so the test reports green either way. It cost 273.6s of the 648.5s the
// e2e suite spent executing on 2026-10-23, 205.5s of that on shard 2, across 27 of 258
// executions, and one of those runs printed four FAILURE lines without failing CI. It stays in
// the repository for local investigation and runs only when PW_DIAGNOSIS=1.
const diagnosisOptIn = !!process.env.PW_DIAGNOSIS

export default defineConfig({
  testDir: './e2e',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  // PW_WORKERS overrides (CI shards set 1 to halve peak browser memory).
  workers: process.env.PW_WORKERS ? Number(process.env.PW_WORKERS) : process.env.CI ? 2 : undefined,
  reporter: process.env.CI ? [['list'], ['html', { open: 'never' }]] : 'list',
  // Generous per-test budget. These specs boot MathLive and KaTeX and stub a dozen endpoints, so
  // the page needs real CPU time before it settles and CI runs two workers across three
  // projects, so a spec can be competing with another browser for most of its wall clock. At
  // 30s this suite failed intermittently on `toBeEnabled` / `toBeVisible` waits that had not
  // been given enough room, which reads as a product fault and is not one. A test that needs
  // longer than this is genuinely stuck and should say so in a minute rather than hide.
  timeout: 60_000,
  use: {
    baseURL,
    trace: 'retain-on-failure',
    // The export registers `public/sw.js`, whose fetch handler calls `fetch` itself, so its
    // requests never pass through `page.route`. A service worker that answered from cache would
    // also let a test pass against the export of an earlier build. Blocking it keeps the fixture's
    // network guard the only thing that can answer a request.
    serviceWorkers: 'block',
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
      // `testIgnore` was the mobile spec only, so this project also collected the diagnosis
      // harness. It is excluded here unless the harness is asked for by name.
      testIgnore: diagnosisOptIn
        ? /mobile(-math)?\.spec\.ts/
        : [/mobile(-math)?\.spec\.ts/, /-diagnosis\.spec\.ts$/],
    },
    {
      name: 'mobile-chromium',
      use: { ...devices['Pixel 7'] },
      // With PW_DIAGNOSIS=1 the diagnosis spec runs on the real touch profiles as well as
      // desktop: the behaviour under investigation is touch-specific, and a desktop-only run
      // would measure a code path a phone never takes.
      testMatch: diagnosisOptIn
        ? [/mobile(-math)?\.spec\.ts/, /visualFigure\.spec\.ts/, /-diagnosis\.spec\.ts$/]
        : [/mobile(-math)?\.spec\.ts/, /visualFigure\.spec\.ts/],
    },
    {
      name: 'mobile-320',
      use: { ...devices['Pixel 7'], viewport: { width: 320, height: 568 } },
      testMatch: diagnosisOptIn
        ? [/mobile(-math)?\.spec\.ts/, /visualFigure\.spec\.ts/, /-diagnosis\.spec\.ts$/]
        : [/mobile(-math)?\.spec\.ts/, /visualFigure\.spec\.ts/],
    },
  ],
  webServer: {
    command: `npx serve out -l ${PORT} --no-clipboard`,
    url: baseURL,
    reuseExistingServer: !process.env.CI,
    timeout: 60_000,
  },
})
