import { test, expect } from '@playwright/test'

const ENTRY = {
  rank: 1,
  student_id: 's1',
  name: 'Ada',
  username: 'ada',
  avatar_dicebear: { style: 'bottts', seed: 's1' },
  avatar_custom: false,
  mastered: 5,
  streak: 2,
  level: 'Novice',
  score: 500,
}

const LEAGUES = {
  week: '2026-W36',
  leagues: [
    {
      tier: 'bronze',
      members: [
        {
          student_id: 's1',
          name: 'Ada',
          username: 'ada',
          avatar_dicebear: { style: 'bottts', seed: 's1' },
          avatar_custom: false,
          tier: 'bronze',
          total_mastered: 5,
          weekly_mastered: 5,
          moved: 0,
        },
      ],
    },
  ],
}

test('leaderboard shows names with avatars in table and leagues', async ({ page }) => {
  await page.route('**/api/leaderboard', route => route.fulfill({ json: [ENTRY] }))
  await page.route('**/api/leagues', route => route.fulfill({ json: LEAGUES }))
  await page.goto('/leaderboard')
  // Table row: avatar image + name.
  const row = page.locator('tbody tr', { hasText: 'Ada' })
  await expect(row).toBeVisible({ timeout: 30_000 })
  await expect(row.getByRole('img', { name: 'Ada avatar' })).toBeVisible()
  // League member: same name + avatar.
  const leagues = page.locator('section', { hasText: 'promotion' })
  await expect(leagues.getByText('Ada').first()).toBeVisible()
  await expect(page.getByRole('img', { name: 'Ada avatar' }).nth(1)).toBeVisible()
})

test('blank names fall back to username, never a blank cell', async ({ page }) => {
  await page.route('**/api/leaderboard', route =>
    route.fulfill({ json: [{ ...ENTRY, name: '   ', username: 'mystery_fox42' }] })
  )
  await page.route('**/api/leagues', route => route.fulfill({ json: { week: '2026-W36', leagues: [] } }))
  await page.goto('/leaderboard')
  await expect(page.locator('tbody tr', { hasText: 'mystery_fox42' })).toBeVisible({ timeout: 30_000 })
})
