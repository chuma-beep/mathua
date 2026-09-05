// Curated DiceBear avatar styles (https://www.dicebear.com, HTTP API v10.x).
// Styles carry individual licenses (CC0/CC-BY mix) — credited in README/footer.
// URLs are seed-stable: the same style+seed always renders the same picture
// and caches in the browser indefinitely. SVG is used (higher rate limit).

export const DICEBEAR_VERSION = '10.x'
const DICEBEAR_BASE = `https://api.dicebear.com/${DICEBEAR_VERSION}`

// Peculiar, friendly picks — robots, pixel art, doodles, adventurers.
export const DICEBEAR_STYLES = [
  'bottts',
  'pixel-art',
  'adventurer',
  'adventurer-neutral',
  'lorelei',
  'croodles',
  'fun-emoji',
  'dylan',
  'micah',
  'notionists',
  'open-peeps',
  'personas',
] as const

export type DicebearStyle = (typeof DICEBEAR_STYLES)[number]

export interface DicebearPick {
  style: DicebearStyle
  seed: string
}

export function isDicebearStyle(s: string): s is DicebearStyle {
  return (DICEBEAR_STYLES as readonly string[]).includes(s)
}

export function dicebearUrl(style: DicebearStyle, seed: string): string {
  return `${DICEBEAR_BASE}/${style}/svg?seed=${encodeURIComponent(seed)}`
}

function randomSeed(): string {
  const rand = Math.random().toString(36).slice(2, 10)
  return `mathua-${rand}-${Date.now().toString(36)}`
}

// randomDicebear rolls a surprise combo: random style + random seed.
export function randomDicebear(): DicebearPick & { url: string } {
  const style = DICEBEAR_STYLES[Math.floor(Math.random() * DICEBEAR_STYLES.length)]
  const seed = randomSeed()
  return { style, seed, url: dicebearUrl(style, seed) }
}

export function pickUrl(pick: DicebearPick): string {
  return dicebearUrl(pick.style, pick.seed)
}

import { getSettings, updateSettings, type UserSettings } from './api'
import { getUserInfo } from './auth'

export interface EnsureDeps {
  getUserInfo: () => { avatar_url?: string } | null
  getSettings: () => Promise<UserSettings>
  updateSettings: (s: UserSettings) => Promise<void>
}

const defaultDeps: EnsureDeps = { getUserInfo, getSettings, updateSettings }

// ensureDicebearAvatar assigns a surprise DiceBear combo to photo-less users
// (no custom upload, no Google photo, no existing pick — including a legacy
// preset, which is grandfathered, never overwritten). Call after auth
// completes (Profile mount covers all flows). Returns the new pick or null.
export async function ensureDicebearAvatar(
  deps: EnsureDeps = defaultDeps,
): Promise<(DicebearPick & { url: string }) | null> {
  const user = deps.getUserInfo()
  if (!user) return null
  if (user.avatar_url) return null
  let settings: UserSettings
  try {
    settings = await deps.getSettings()
  } catch {
    return null
  }
  if (settings.avatar_custom || settings.avatar_dicebear || settings.avatar_preset != null) return null
  const pick = randomDicebear()
  try {
    await deps.updateSettings({ ...settings, avatar_dicebear: { style: pick.style, seed: pick.seed } })
  } catch {
    return null
  }
  return pick
}
