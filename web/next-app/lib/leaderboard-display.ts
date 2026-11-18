import type { LeaderboardAvatar } from './api'
import { resolveAvatar } from './dicebear'

type Named = { id: string; name: string }
type Avatared = Named & { avatar?: LeaderboardAvatar }

// The server decides the avatar shape; the client only turns it into a URL. A
// generated avatar is seeded from the public id, never from the account id.
export function avatarURL(row: Avatared): string | undefined {
  const av = row.avatar
  if (!av) return undefined
  if (av.kind === 'custom' && av.url) {
    return av.version ? `${av.url}?v=${av.version}` : av.url
  }
  return resolveAvatar(
    { avatar_url: '' },
    { avatar_dicebear: av.style ? { style: av.style, seed: row.id } : null },
  ).url
}

// The server resolves a display name, but a blank cell is a silent failure.
// Fall back to the public id, which is never empty.
export function nameOf(row: Named): string {
  return row.name.trim() || row.id
}

export function shortDay(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return ''
  return d.toLocaleDateString('en-GB', { day: 'numeric', month: 'short', timeZone: 'UTC' })
}
