// Canonical public origin for metadata, robots and sitemap.
// Override per deployment with NEXT_PUBLIC_SITE_URL (e.g. a self-hosted
// instance); defaults to the public origin.
export const SITE_URL = process.env.NEXT_PUBLIC_SITE_URL ?? 'https://mathua.com'

// Public, indexable routes only. App views (profile, session, settings,
// goals, diagnose, onboard, admin) and query-driven pages stay out.
//
// /study is gone from this list because the route is closed: it forwards to /learn, and a
// search result that redirects is worse than no result. /domains takes its place as the
// curriculum-orientation surface, which is what the library was being indexed for.
export const INDEXABLE_ROUTES = [
  '/',
  '/how-it-works',
  '/domains',
  '/graph',
  '/leaderboard',
  '/glossary',
  '/docs',
  '/docs/architecture',
  '/docs/administration',
  '/docs/contributing',
  '/docs/efficacy',
  '/docs/math-input',
  '/docs/system-design',
  '/licenses',
  '/privacy',
  '/note',
  '/login',
] as const
