// Canonical public origin for metadata, robots and sitemap.
// Override per deployment with NEXT_PUBLIC_SITE_URL (e.g. a self-hosted
// Fly instance); defaults to the Vercel front-end.
export const SITE_URL = process.env.NEXT_PUBLIC_SITE_URL ?? 'https://mathua.vercel.app'

// Public, indexable routes only. App views (profile, session, settings,
// goals, diagnose, onboard, admin) and query-driven pages stay out.
export const INDEXABLE_ROUTES = [
  '/',
  '/how-it-works',
  '/study',
  '/graph',
  '/leaderboard',
  '/docs',
  '/docs/architecture',
  '/docs/contributing',
  '/docs/efficacy',
  '/docs/system-design',
  '/note',
  '/login',
] as const
