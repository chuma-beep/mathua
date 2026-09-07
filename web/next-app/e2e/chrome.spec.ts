import { test, expect } from '@playwright/test'

// Regression: layout + page both rendered chrome on these routes
// (double headers on /docs, double bottom tab bars on /diagnose + /goals).
// Pages own their chrome; layouts are metadata-only (docs/layout keeps its
// section header by design).
test.beforeEach(async ({ context }) => {
  await context.route('**/api/**', route => route.fulfill({ status: 404, body: '' }))
})

for (const route of ['/docs', '/diagnose', '/goals', '/how-it-works']) {
  test(`${route} renders exactly one header, footer and bottom tab bar`, async ({ page }) => {
    await page.goto(route)
    await expect(page.locator('header')).toHaveCount(1, { timeout: 30_000 })
    await expect(page.getByLabel('Profile', { exact: true })).toHaveCount(1)
    if (route === '/docs') {
      await expect(page.locator('footer')).toHaveCount(1)
      // Docs section header survives; the default app header must not.
      await expect(page.getByText('Architecture and contributing guide')).toBeVisible()
    }
  })
}
