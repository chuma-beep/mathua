import { getPushKey, subscribePush, unsubscribePush } from './api'

// The subscribe flow must stay usable before React hydrates and after a service
// worker update, so it reads browser state at call time rather than caching it.
// A browser without PushManager cannot receive anything, and returning a named
// result keeps the settings copy specific instead of a generic failure.

export type PushSubscribeResult = 'subscribed' | 'unsupported' | 'unavailable' | 'denied'

export function browserPushSupported(): boolean {
  return (
    typeof window !== 'undefined' &&
    'serviceWorker' in navigator &&
    'PushManager' in window &&
    typeof Notification !== 'undefined'
  )
}

// decodeBase64URL turns a VAPID key into the byte array pushManager expects.
// The URL-safe alphabet and the missing padding are both part of the format.
function decodeBase64URL(raw: string): Uint8Array<ArrayBuffer> {
  const padded = raw.replace(/-/g, '+').replace(/_/g, '/')
  const binary = atob(padded + '='.repeat((4 - (padded.length % 4)) % 4))
  const bytes = new Uint8Array(new ArrayBuffer(binary.length))
  for (let i = 0; i < binary.length; i += 1) {
    bytes[i] = binary.charCodeAt(i)
  }
  return bytes
}

async function activeRegistration(): Promise<ServiceWorkerRegistration | null> {
  if (!('serviceWorker' in navigator)) return null
  const existing = await navigator.serviceWorker.getRegistration('/')
  if (existing) return existing
  try {
    return await navigator.serviceWorker.ready
  } catch {
    return null
  }
}

export async function subscribeBrowserPush(): Promise<PushSubscribeResult> {
  if (!browserPushSupported()) return 'unsupported'

  const key = await getPushKey()
  if (!key.configured || !key.public_key) return 'unavailable'

  // Ask for permission before touching the subscription, so a denial leaves no
  // half-created endpoint behind.
  const permission = await Notification.requestPermission()
  if (permission !== 'granted') return 'denied'

  const registration = await activeRegistration()
  if (!registration) return 'unsupported'

  const existing = await registration.pushManager.getSubscription()
  const subscription =
    existing ??
    (await registration.pushManager.subscribe({
      userVisibleOnly: true,
      applicationServerKey: decodeBase64URL(key.public_key),
    }))

  const json = subscription.toJSON()
  const p256dh = json.keys?.p256dh ?? ''
  const auth = json.keys?.auth ?? ''
  if (!json.endpoint || !p256dh || !auth) return 'unsupported'

  await subscribePush({ endpoint: json.endpoint, keys: { p256dh, auth } })
  return 'subscribed'
}

export async function unsubscribeBrowserPush(): Promise<void> {
  if (!browserPushSupported()) return
  const registration = await activeRegistration()
  if (!registration) return
  const subscription = await registration.pushManager.getSubscription()
  if (!subscription) return
  // Tell the server first: if the browser drops the subscription before the
  // request lands, the endpoint stays live in the table with no device to
  // deliver to.
  await unsubscribePush(subscription.endpoint)
  await subscription.unsubscribe()
}
