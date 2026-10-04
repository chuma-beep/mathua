import { defineConfig, devices } from '@playwright/test'

const PORT = 4173
const baseURL = `http://localhost:${PORT}`

export default defineConfig({
  testDir: './e2e',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  // PW_WORKERS overrides (CI shards set 1 to halve peak browser memory).
  workers: process.env.PW_WORKERS ? Number(process.env.PW_WORKERS) : process.env.CI ? 2 : undefined,
  reporter: process.env.CI ? [['list'], ['html', { open: 'never' }]] : 'list',
  timeout: 30_000,
  use: {
    baseURL,
    trace: 'retain-on-failure',
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
      testIgnore: /mobile(-math)?\.spec\.ts/,
    },
    {
      name: 'mobile-chromium',
      use: { ...devices['Pixel 7'] },
      // The diagnosis spec runs on the real touch profiles as well as desktop: the
      // behaviour under investigation is touch-specific, and a desktop-only run
      // would measure a code path a phone never takes.
      testMatch: [/mobile(-math)?\.spec\.ts/, /-diagnosis\.spec\.ts$/],
    },
    {
      name: 'mobile-320',
      use: { ...devices['Pixel 7'], viewport: { width: 320, height: 568 } },
      testMatch: [/mobile(-math)?\.spec\.ts/, /-diagnosis\.spec\.ts$/],
    },
  ],
  webServer: {
    command: `npx serve out -l ${PORT} --no-clipboard`,
    url: baseURL,
    reuseExistingServer: !process.env.CI,
    timeout: 60_000,
  },
})
