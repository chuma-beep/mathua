// The star count is public data, read from GitHub in the browser. It is cached
// for six hours because the unauthenticated API allows sixty requests an hour
// per address, and a header that fetched on every navigation would spend that
// budget on one reader. A failure simply leaves the label as "Star", never an
// error in the chrome.

const REPO = 'chuma-beep/mathua'
const CACHE_KEY = 'mathua:github-stars:v1'
const MAX_AGE_MS = 6 * 60 * 60 * 1000

interface CachedStars {
	count: number
	at: number
}

export function formatStarCount(count: number): string {
	if (!Number.isFinite(count) || count < 0) return ''
	if (count < 1000) return String(Math.round(count))
	const thousands = count / 1000
	// GitHub's own rule: one decimal below a hundred thousand, whole numbers
	// above it, and never a trailing .0.
	if (thousands >= 100) return `${Math.round(thousands)}k`
	return `${thousands.toFixed(1).replace(/\.0$/, '')}k`
}

function readCached(now: number): number | null {
	if (typeof window === 'undefined') return null
	try {
		const raw = window.localStorage.getItem(CACHE_KEY)
		if (!raw) return null
		const parsed = JSON.parse(raw) as CachedStars
		if (typeof parsed.count !== 'number' || typeof parsed.at !== 'number') return null
		if (now - parsed.at > MAX_AGE_MS) return null
		return parsed.count
	} catch {
		return null
	}
}

function writeCached(count: number, now: number): void {
	if (typeof window === 'undefined') return
	try {
		window.localStorage.setItem(CACHE_KEY, JSON.stringify({ count, at: now }))
	} catch {
		// Storage can be unavailable in private modes. The count is not worth
		// failing a render over.
	}
}

export async function loadGitHubStars(): Promise<number | null> {
	const now = Date.now()
	const cached = readCached(now)
	if (cached !== null) return cached
	try {
		const res = await fetch(`https://api.github.com/repos/${REPO}`, {
			headers: { Accept: 'application/vnd.github+json' },
		})
		if (!res.ok) return null
		const body = (await res.json()) as { stargazers_count?: unknown }
		const count = body.stargazers_count
		if (typeof count !== 'number' || !Number.isFinite(count)) return null
		writeCached(count, now)
		return count
	} catch {
		return null
	}
}
