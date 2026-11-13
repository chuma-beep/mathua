import { afterEach, describe, expect, it, vi } from 'vitest'
import { resolveAPIBase } from '../lib/apiBase'

afterEach(() => { vi.unstubAllEnvs() })

// A static export has no API routes. If the base URL is missing, every call
// 404s on the page's own origin and the deployment looks broken; the build must
// refuse instead.
describe('resolveAPIBase', () => {
  it('uses the configured origin and drops a trailing slash', () => {
    expect(resolveAPIBase('https://mathua.fly.dev/')).toBe('https://mathua.fly.dev')
    expect(resolveAPIBase('  https://api.example.com  ')).toBe('https://api.example.com')
  })

  it('fails a production build that has no API origin', () => {
    vi.stubEnv('NODE_ENV', 'production')
    expect(() => resolveAPIBase('')).toThrow(/NEXT_PUBLIC_API_URL is required/)
    expect(() => resolveAPIBase(undefined)).toThrow(/NEXT_PUBLIC_API_URL is required/)
  })

  it('keeps same-origin relative calls in development', () => {
    vi.stubEnv('NODE_ENV', 'development')
    expect(resolveAPIBase(undefined)).toBe('')
  })
})

// The Go image serves the export and the API together: the same empty base that
// is fatal on a static host is correct there, and only an explicit flag may say so.
it('accepts an explicit same-origin build', () => {
  vi.stubEnv('NODE_ENV', 'production')
  expect(resolveAPIBase('', true)).toBe('')
})

// A developer .env that names localhost must not survive into a Worker or
// Docker build. The same-origin flag wins over any configured origin.
it('prefers same-origin over a configured origin', () => {
  vi.stubEnv('NODE_ENV', 'production')
  expect(resolveAPIBase('http://localhost:8080', true)).toBe('')
})
