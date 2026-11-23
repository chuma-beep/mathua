import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { formatStarCount, loadGitHubStars } from '../lib/githubStars'

beforeEach(() => {
	window.localStorage.clear()
})

afterEach(() => {
	vi.unstubAllGlobals()
})

describe('formatStarCount', () => {
	it('leaves counts below a thousand exact', () => {
		expect(formatStarCount(0)).toBe('0')
		expect(formatStarCount(999)).toBe('999')
	})

	it('shortens thousands the way GitHub does', () => {
		expect(formatStarCount(1000)).toBe('1k')
		expect(formatStarCount(1234)).toBe('1.2k')
		expect(formatStarCount(12345)).toBe('12.3k')
		expect(formatStarCount(123456)).toBe('123k')
	})

	it('refuses nonsense rather than printing NaN', () => {
		expect(formatStarCount(Number.NaN)).toBe('')
		expect(formatStarCount(-1)).toBe('')
	})
})

describe('loadGitHubStars', () => {
	it('reads once and serves the cached count afterwards', async () => {
		const fetchMock = vi.fn().mockResolvedValue({
			ok: true,
			json: async () => ({ count: 1234 }),
		})
		vi.stubGlobal('fetch', fetchMock)

		expect(await loadGitHubStars()).toBe(1234)
		expect(await loadGitHubStars()).toBe(1234)
		expect(fetchMock).toHaveBeenCalledTimes(1)
	})

	it('returns null on a failed response and caches nothing', async () => {
		const fetchMock = vi.fn().mockResolvedValue({ ok: false, json: async () => ({}) })
		vi.stubGlobal('fetch', fetchMock)

		expect(await loadGitHubStars()).toBeNull()
		expect(window.localStorage.getItem('mathua:github-stars:v1')).toBeNull()
	})
})
