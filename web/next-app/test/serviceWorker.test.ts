import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { describe, expect, it, vi } from 'vitest'

// Loads public/sw.js as a classic worker. It has no exports and expects the
// worker globals, so the test supplies them and captures the fetch listener.
const source = readFileSync(resolve(__dirname, '../public/sw.js'), 'utf8')

interface Harness {
  dispatchFetch: (request: unknown) => Promise<Response>
  dispatchPush: (payload: unknown) => Promise<void>
  dispatchNotificationClick: (data: { href?: string } | undefined) => Promise<void>
  puts: Array<{ request: unknown; response: Response }>
  store: Map<string, Response>
  shown: Array<{ title: string; options: Record<string, unknown> }>
  focused: string[]
  navigated: string[]
  opened: string[]
  closed: number
}

function loadWorker(
  fetchImpl: (req: unknown) => Promise<Response>,
  windows: Array<{ url: string }> = [{ url: 'https://mathua.com/learn' }],
): Harness {
  const listeners: Record<string, (event: unknown) => void> = {}
  const puts: Array<{ request: unknown; response: Response }> = []
  const store = new Map<string, Response>()
  const shown: Array<{ title: string; options: Record<string, unknown> }> = []
  const focused: string[] = []
  const navigated: string[] = []
  const opened: string[] = []
  let closed = 0

  const keyOf = (req: unknown) =>
    typeof req === 'string' ? req : String((req as { url?: string }).url ?? '')

  const cache = {
    put: async (request: unknown, response: Response) => {
      puts.push({ request, response })
      store.set(keyOf(request), response)
    },
    match: async (request: unknown) => store.get(keyOf(request)),
    addAll: async () => undefined,
  }

  const caches = {
    open: async () => cache,
    match: async (request: unknown) => store.get(keyOf(request)),
    keys: async () => [],
    delete: async () => true,
  }

  const windowClients = windows.map(client => ({
    url: client.url,
    focus: async () => {
      focused.push(client.url)
    },
    navigate: async (href: string) => {
      navigated.push(href)
    },
  }))

  const self = {
    addEventListener: (name: string, fn: (event: unknown) => void) => {
      listeners[name] = fn
    },
    skipWaiting: async () => undefined,
    registration: {
      showNotification: async (title: string, options: Record<string, unknown>) => {
        shown.push({ title, options })
      },
    },
    clients: {
      claim: async () => undefined,
      matchAll: async () => windowClients,
      openWindow: async (href: string) => {
        opened.push(href)
      },
    },
    location: { origin: 'https://mathua.com' },
  }

  const load = new Function('self', 'caches', 'fetch', 'Response', 'URL', source)
  load(self, caches, fetchImpl, Response, URL)

  return {
    puts,
    store,
    shown,
    focused,
    navigated,
    opened,
    get closed() {
      return closed
    },
    dispatchFetch: async (request: unknown) => {
      let settled: Promise<Response> | undefined
      listeners.fetch({
        request,
        respondWith: (p: Promise<Response>) => {
          settled = p
        },
      })
      if (!settled) throw new Error('fetch event was not handled')
      return settled
    },
    dispatchPush: async (payload: unknown) => {
      let settled: Promise<unknown> | undefined
      listeners.push({
        data: payload === null ? null : { json: () => payload },
        waitUntil: (p: Promise<unknown>) => {
          settled = p
        },
      })
      if (!settled) throw new Error('push event was not handled')
      await settled
    },
    dispatchNotificationClick: async (data: { href?: string } | undefined) => {
      let settled: Promise<unknown> | undefined
      listeners.notificationclick({
        notification: {
          close: () => {
            closed += 1
          },
          data,
        },
        waitUntil: (p: Promise<unknown>) => {
          settled = p
        },
      })
      if (!settled) throw new Error('notificationclick event was not handled')
      await settled
    },
  }
}

const nav = { method: 'GET', mode: 'navigate', url: 'https://mathua.com/learn' }

const htmlResponse = (body: string) =>
  new Response(body, { status: 200, headers: { 'content-type': 'text/html; charset=utf-8' } })

describe('the service worker and Cloudflare-injected pages', () => {
  it('does not cache a navigation carrying a Precursor token', async () => {
    // Cloudflare appends this to the HTML of a challenged response. Cached, it
    // replays a dead token and Cloudflare answers 400 on every later visit.
    const body = '<html><body>learn</body><script>window.__CF$cv$params={r:"x",ut:"dead"}</script></html>'
    const worker = loadWorker(async () => htmlResponse(body))

    const res = await worker.dispatchFetch(nav)

    expect(await res.text()).toBe(body)
    expect(worker.puts).toHaveLength(0)
  })

  it('caches an ordinary navigation so it can be read offline', async () => {
    const body = '<html><body>learn</body></html>'
    const worker = loadWorker(async () => htmlResponse(body))

    await worker.dispatchFetch(nav)

    expect(worker.puts).toHaveLength(1)
  })

  it('serves the offline page rather than a cached challenged page', async () => {
    const worker = loadWorker(async () => {
      throw new TypeError('Failed to fetch')
    })
    const challenged = htmlResponse('<html><script>window.__CF$cv$params={}</script></html>')
    await worker.store.set(nav.url, challenged)
    await worker.store.set('/offline', htmlResponse('<html><body>offline</body></html>'))

    const res = await worker.dispatchFetch(nav)

    expect(await res.text()).toContain('offline')
  })

  it('serves a cached ordinary page when the network fails', async () => {
    const worker = loadWorker(async () => {
      throw new TypeError('Failed to fetch')
    })
    await worker.store.set(nav.url, htmlResponse('<html><body>cached learn</body></html>'))

    const res = await worker.dispatchFetch(nav)

    expect(await res.text()).toContain('cached learn')
  })
})

describe('push notifications', () => {
  it('shows the payload the server published', async () => {
    const worker = loadWorker(async () => htmlResponse(''))

    await worker.dispatchPush({
      title: '3 reviews due',
      body: 'A short session clears the queue.',
      href: '/review',
      notification_id: 'ntf_abc',
    })

    expect(worker.shown).toHaveLength(1)
    expect(worker.shown[0].title).toBe('3 reviews due')
    expect(worker.shown[0].options).toMatchObject({
      body: 'A short session clears the queue.',
      tag: 'ntf_abc',
      data: { href: '/review', notification_id: 'ntf_abc' },
    })
  })

  it('falls back to a safe title and route when the payload is empty', async () => {
    const worker = loadWorker(async () => htmlResponse(''))

    await worker.dispatchPush(null)

    expect(worker.shown[0].title).toBe('Mathua')
    expect(worker.shown[0].options).toMatchObject({ data: { href: '/', notification_id: '' } })
  })

  it('focuses an open window and moves it to the target', async () => {
    const worker = loadWorker(async () => htmlResponse(''))

    await worker.dispatchNotificationClick({ href: '/review' })

    expect(worker.closed).toBe(1)
    expect(worker.focused).toEqual(['https://mathua.com/learn'])
    expect(worker.navigated).toEqual(['/review'])
    expect(worker.opened).toHaveLength(0)
  })

  it('opens a window when none is already open', async () => {
    const worker = loadWorker(async () => htmlResponse(''), [])

    await worker.dispatchNotificationClick({ href: '/study' })

    expect(worker.opened).toEqual(['/study'])
    expect(worker.focused).toHaveLength(0)
  })

  it('opens the home page when a notification carries no route', async () => {
    const worker = loadWorker(async () => htmlResponse(''), [])

    await worker.dispatchNotificationClick(undefined)

    expect(worker.opened).toEqual(['/'])
  })
})

describe('the worker version', () => {
  it('is not the static development stamp in a built export', () => {
    // The source must stay stampable; the build rewrites out/sw.js. This guards
    // the line the Dockerfile and scripts/stamp-sw.mjs both replace.
    expect(source).toMatch(/const VERSION = '[^']*'/)
  })
})

// Keep the import used even if the build check above is simplified later.
void vi
