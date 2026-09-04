/* Mathua service worker — hand-rolled for the Next.js static export.
 *
 * Strategy:
 *   - navigations (HTML): network-first, cache fallback, then /offline.html
 *   - hashed static assets + images + fonts: cache-first (immutable)
 *   - GET /api/*: network-first with cache fallback (read-only endpoints;
 *     non-GET requests are never intercepted so answer submissions fail
 *     loudly instead of silently succeeding from a stale cache)
 */
const VERSION = 'v2'
const SHELL_CACHE = `mathua-shell-${VERSION}`
const ASSET_CACHE = `mathua-assets-${VERSION}`
const DATA_CACHE = `mathua-data-${VERSION}`

// '/offline' resolves to offline.html on both the Go server (clean-URL
// rewrite) and static hosts; '.html' would 301 on some hosts, and a
// redirected cached response is invalid for navigations.
const PRECACHE_URLS = ['/', '/offline']

const ASSET_PATTERN = /\/_next\/static\/|\/diagrams\/|\/icons\/|\.(?:png|jpe?g|gif|svg|ico|webp|woff2?|ttf|css|js|webmanifest)$/i

self.addEventListener('install', event => {
  event.waitUntil(
    caches
      .open(SHELL_CACHE)
      .then(cache => cache.addAll(PRECACHE_URLS))
      .then(() => self.skipWaiting())
  )
})

self.addEventListener('activate', event => {
  event.waitUntil(
    caches
      .keys()
      .then(keys =>
        Promise.all(
          keys
            .filter(key => ![SHELL_CACHE, ASSET_CACHE, DATA_CACHE].includes(key))
            .map(key => caches.delete(key))
        )
      )
      .then(() => self.clients.claim())
  )
})

async function networkFirst(request, cacheName, fallbackUrls) {
  const cache = await caches.open(cacheName)
  try {
    const response = await fetch(request)
    if (response.ok) cache.put(request, response.clone())
    return response
  } catch (err) {
    const cached = await cache.match(request)
    if (cached) return cached
    for (const url of fallbackUrls ?? []) {
      const fallback = await caches.match(url)
      if (fallback) return fallback
    }
    throw err
  }
}

async function cacheFirst(request, cacheName) {
  const cache = await caches.open(cacheName)
  const cached = await cache.match(request)
  if (cached) return cached
  const response = await fetch(request)
  if (response.ok) cache.put(request, response.clone())
  return response
}

self.addEventListener('fetch', event => {
  const { request } = event
  if (request.method !== 'GET') return

  const url = new URL(request.url)
  if (url.origin !== self.location.origin) return

  if (request.mode === 'navigate') {
    event.respondWith(networkFirst(request, SHELL_CACHE, ['/offline', '/offline.html']))
    return
  }

  if (ASSET_PATTERN.test(url.pathname)) {
    event.respondWith(cacheFirst(request, ASSET_CACHE))
    return
  }

  if (url.pathname.startsWith('/api/')) {
    event.respondWith(networkFirst(request, DATA_CACHE))
  }
})
