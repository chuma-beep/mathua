/* Mathua service worker, hand-rolled for the Next.js static export.
 *
 * Strategy:
 *   - navigations (HTML): network-first, cache fallback, then /offline.html
 *   - hashed static assets + images + fonts: cache-first (immutable)
 *   - the API is never intercepted. Its responses are per-learner and must
 *     not be served from a cache shared by whoever uses this browser next;
 *     a stale GET here would show another account's progress offline.
 *
 * VERSION is stamped at image build time (Dockerfile rewrites the line below
 * with a UTC timestamp) so every deploy ships byte-different sw.js: browsers
 * install the new worker, and activate purges all caches from older builds.
 * A static version would let clients run stale hashed JS indefinitely, old
 * chunks 404 on fresh images while this worker keeps serving its cached
 * copies (cache-first), pinning the client on pre-fix rendering forever.
 */
const VERSION = 'v2-local'
const SHELL_CACHE = `mathua-shell-${VERSION}`
const ASSET_CACHE = `mathua-assets-${VERSION}`

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
            .filter(key => ![SHELL_CACHE, ASSET_CACHE].includes(key))
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

// Cloudflare injects a per-request Precursor token into the HTML:
// `window.__CF$cv$params = {r, t, u, ut, i}` followed by a script that appends
// /cdn-cgi/challenge-platform/scripts/precursor/main.js. The token is only valid
// for the response it arrived with. Storing that page and replaying it later
// makes the injected script post a dead token, and Cloudflare answers 400 to a
// challenge the page can no longer complete, which is the `h/undefined/...`
// request in the console, repeated on every visit because the cache never
// refreshes. Store navigation HTML only when it carries no token, and never
// serve a cached page that does.
const CHALLENGE_MARKER = '__CF$cv$params'

async function hasChallengeToken(response) {
  try {
    const type = response.headers.get('content-type') || ''
    if (!type.includes('text/html')) return false
    const text = await response.clone().text()
    return text.includes(CHALLENGE_MARKER)
  } catch {
    // If the body cannot be read, treat the page as uncacheable: a missed
    // offline page is a smaller harm than replaying a stale challenge.
    return true
  }
}

async function navigationNetworkFirst(request) {
  const cache = await caches.open(SHELL_CACHE)
  try {
    const response = await fetch(request)
    if (response.ok && !(await hasChallengeToken(response))) {
      cache.put(request, response.clone())
    }
    return response
  } catch (err) {
    const cached = await cache.match(request)
    if (cached && !(await hasChallengeToken(cached))) return cached
    const offline = await caches.match('/offline')
    if (offline) return offline
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
    event.respondWith(navigationNetworkFirst(request))
    return
  }

  if (ASSET_PATTERN.test(url.pathname)) {
    event.respondWith(cacheFirst(request, ASSET_CACHE))
    return
  }

})

/* Push notifications.
 *
 * The payload is the JSON the server publishes for one delivery: title, body,
 * href and notification_id. Nothing here fetches anything, because a push can
 * arrive when the learner is offline and a fetch that hangs keeps the event
 * alive until the browser times it out.
 *
 * Both handlers look their globals up inside the callback. The worker source is
 * also executed by the unit test harness, which supplies only self, caches,
 * fetch, Response and URL; a module-scope self.registration would fail there
 * before any event is dispatched.
 */
self.addEventListener('push', event => {
  let data = {}
  try {
    data = event.data ? event.data.json() : {}
  } catch (err) {
    data = {}
  }
  const title = data.title || 'Mathua'
  const options = {
    body: data.body || '',
    icon: '/icons/icon-192.png',
    badge: '/icons/icon-192.png',
    tag: data.notification_id || undefined,
    data: {
      href: data.href || '/',
      notification_id: data.notification_id || '',
    },
  }
  event.waitUntil(self.registration.showNotification(title, options))
})

/* A click focuses an existing same-origin window and moves it to the target,
 * or opens one when none is open. Reusing the window matters on mobile, where
 * openWindow would stack a second copy of the app behind the first. */
self.addEventListener('notificationclick', event => {
  event.notification.close()
  const data = event.notification.data || {}
  const base = data.href || '/'
  // The click is reported by the page, not the worker, so the learner's own
  // session authenticates it. Tagging the URL keeps the report idempotent: the
  // app strips ?ntf after sending it once.
  const href = data.notification_id
    ? base + (base.includes('?') ? '&' : '?') + 'ntf=' + encodeURIComponent(data.notification_id)
    : base

  event.waitUntil((async () => {
    const windows = await self.clients.matchAll({ type: 'window', includeUncontrolled: true })
    for (const client of windows) {
      if (new URL(client.url).origin !== self.location.origin) continue
      await client.focus()
      if (typeof client.navigate === 'function') {
        await client.navigate(href)
      }
      return
    }
    await self.clients.openWindow(href)
  })())
})
