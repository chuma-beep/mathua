import { describe, it, expect, vi } from 'vitest'
import {
  DICEBEAR_STYLES,
  dicebearUrl,
  isDicebearStyle,
  randomDicebear,
  ensureDicebearAvatar,
  resolveAvatar,
} from '../lib/dicebear'

describe('dicebear', () => {
  it('builds pinned, encoded URLs', () => {
    expect(dicebearUrl('bottts', 'abc')).toBe('https://api.dicebear.com/10.x/bottts/svg?seed=abc')
    expect(dicebearUrl('pixel-art', 'a b')).toContain('seed=a%20b')
  })

  it('curates a non-empty gallery of known styles', () => {
    expect(DICEBEAR_STYLES.length).toBeGreaterThanOrEqual(8)
    for (const s of DICEBEAR_STYLES) expect(isDicebearStyle(s)).toBe(true)
    expect(isDicebearStyle('nope')).toBe(false)
  })

  it('randomDicebear rolls valid combos', () => {
    for (let i = 0; i < 20; i++) {
      const pick = randomDicebear()
      expect(isDicebearStyle(pick.style)).toBe(true)
      expect(pick.seed.length).toBeGreaterThan(0)
      expect(pick.url).toBe(dicebearUrl(pick.style, pick.seed))
    }
  })

  it('ensureDicebearAvatar assigns a surprise to photo-less users', async () => {
    const updateSettings = vi.fn().mockResolvedValue(undefined)
    const deps = {
      getUserInfo: () => ({}),
      getSettings: async () => ({}),
      updateSettings,
    }
    const pick = await ensureDicebearAvatar(deps)
    expect(pick).not.toBeNull()
    expect(updateSettings).toHaveBeenCalledOnce()
    expect(updateSettings.mock.calls[0][0].avatar_dicebear).toEqual({ style: pick!.style, seed: pick!.seed })
  })

  it('ensureDicebearAvatar respects existing choices (photo, pick, legacy preset)', async () => {
    const updateSettings = vi.fn()
    const base = { getSettings: async () => ({}), updateSettings }
    expect(await ensureDicebearAvatar({ ...base, getUserInfo: () => ({ avatar_url: 'https://x/y.png' }) })).toBeNull()
    expect(await ensureDicebearAvatar({ ...base, getUserInfo: () => ({}), getSettings: async () => ({ avatar_dicebear: { style: 'bottts', seed: 's' } }) })).toBeNull()
    expect(await ensureDicebearAvatar({ ...base, getUserInfo: () => ({}), getSettings: async () => ({ avatar_preset: 3 }) })).toBeNull()
    expect(await ensureDicebearAvatar({ ...base, getUserInfo: () => null })).toBeNull()
    expect(updateSettings).not.toHaveBeenCalled()
  })

  it('ensureDicebearAvatar swallows backend failures', async () => {
    const deps = {
      getUserInfo: () => ({}),
      getSettings: async () => { throw new Error('offline') },
      updateSettings: vi.fn(),
    }
    expect(await ensureDicebearAvatar(deps)).toBeNull()
  })
})

describe('resolveAvatar', () => {
  const user = { avatar_url: 'https://google/photo.jpg' }

  it('prefers custom upload with roaming version', () => {
    const r = resolveAvatar(user, { avatar_custom: true, avatar_version: 123 })
    expect(r.url).toContain('/api/avatar/me?v=123')
    expect(r.preset).toBeNull()
  })

  it('prefers Google photo over dicebear pick', () => {
    const r = resolveAvatar(user, { avatar_dicebear: { style: 'bottts', seed: 's' } })
    expect(r.url).toBe('https://google/photo.jpg')
  })

  it('uses the dicebear pick without a photo', () => {
    const r = resolveAvatar(null, { avatar_dicebear: { style: 'bottts', seed: 's' } })
    expect(r.url).toBe('https://api.dicebear.com/10.x/bottts/svg?seed=s')
  })

  it('grandfathers legacy presets, ignores invalid styles', () => {
    expect(resolveAvatar(null, { avatar_preset: 3 }).preset).toBe(3)
    expect(resolveAvatar(null, { avatar_dicebear: { style: 'nope', seed: 's' } })).toEqual({
      url: undefined,
      preset: null,
    })
  })

  it('falls back to automatic initial with nothing set', () => {
    expect(resolveAvatar(null, {})).toEqual({ url: undefined, preset: null })
  })
})
