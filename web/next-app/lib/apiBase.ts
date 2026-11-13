// One place that decides where the API lives.
//
// The static export is served by a host that has no API routes, so a build
// without NEXT_PUBLIC_API_URL does not degrade, every request becomes a 404
// on the page's own origin and the site looks broken while loading perfectly.
// Failing the build is the only failure mode a deploy cannot ship.
export function resolveAPIBase(
  configured: string | undefined,
  sameOrigin = process.env.NEXT_PUBLIC_SAME_ORIGIN_API === '1',
): string {
  // A same-origin build serves the export and the API together, so the base is
  // intentionally empty. This wins over a configured value because a developer
  // .env with localhost must never be baked into a Worker or Docker build.
  if (sameOrigin) return ''
  const base = (configured ?? '').trim().replace(/\/+$/, '')
  if (base) return base
  // The Go image serves the static export and the API on one origin, so an
  // empty base is correct there. It has to be declared, because the same empty
  // value on a static host is the 404-everything failure this guard exists for.
  if (process.env.NODE_ENV === 'production') {
    throw new Error(
      'NEXT_PUBLIC_API_URL is required: without it the browser calls /api/* on the static host, where every request is a 404.',
    )
  }
  return ''
}

export const API_BASE = resolveAPIBase(process.env.NEXT_PUBLIC_API_URL)
