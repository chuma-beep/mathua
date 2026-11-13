/**
 * Mathua's Cloudflare Worker.
 *
 * The static export is served from Workers static assets. Everything under
 * /api is proxied to the Fly origin, which keeps the database, the Go server
 * and the in-memory diagnostic and quiz sessions where they already run.
 *
 * The proxy adds X-Mathua-Edge-Secret. The origin refuses /api requests
 * without it once MATHUA_EDGE_SECRET is configured, so a direct request to
 * mathua.fly.dev cannot choose its own rate-limit bucket by forging
 * CF-Connecting-IP.
 */
const ORIGIN = 'https://mathua.fly.dev'
const APEX = 'mathua.com'
const EDGE_SECRET_HEADER = 'X-Mathua-Edge-Secret'

function isAPI(pathname) {
  return pathname === '/api' || pathname.startsWith('/api/')
}

async function proxyAPI(request, env) {
  const target = new URL(request.url)
  target.protocol = 'https:'
  target.hostname = 'mathua.fly.dev'

  const headers = new Headers(request.headers)
  headers.set(EDGE_SECRET_HEADER, env.EDGE_SECRET)
  // Fly's edge overwrites Fly-Client-IP with the Worker's address. The Go
  // server reads CF-Connecting-IP instead, which Cloudflare sets for the real
  // client and which the guard makes safe to trust.
  headers.delete('Fly-Client-IP')

  const init = {
    method: request.method,
    headers,
    redirect: 'manual',
  }
  if (request.method !== 'GET' && request.method !== 'HEAD') {
    init.body = request.body
  }

  try {
    return await fetch(target.toString(), init)
  } catch {
    return new Response(JSON.stringify({ error: 'origin unavailable' }), {
      status: 502,
      headers: { 'Content-Type': 'application/json' },
    })
  }
}

export default {
  async fetch(request, env) {
    const url = new URL(request.url)

    if (url.hostname === `www.${APEX}`) {
      url.hostname = APEX
      return Response.redirect(url.toString(), 301)
    }

    if (isAPI(url.pathname)) {
      return proxyAPI(request, env)
    }

    return env.ASSETS.fetch(request)
  },
}
