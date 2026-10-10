import { describe, it, expect } from 'vitest'
import fs from 'node:fs'
import path from 'node:path'
import { LOGIN_HERO } from '../lib/loginHero'

import { APP_ROOT } from './helpers/roots'

// The hero is an <img> at `inset-0 size-full object-cover`, so the panel sizes it
// rather than the other way round. That makes a missing file a 404 on the login
// page -- the one page a new learner is guaranteed to reach -- and it makes a
// wrong `width`/`height` a wrong box reserved before the file arrives. Neither
// throws: the image still lays out, the panel still looks plausible, and nothing
// else in the suite looks at the hero.

type Dimensions = { width: number; height: number }

/** Real pixel dimensions, read from the file's own header. */
function dimensionsOf(abs: string): Dimensions {
  const head = fs.readFileSync(abs).subarray(0, 64 * 1024)

  // JPEG: walk the marker segments to the first SOF, which carries the frame size.
  if (head[0] === 0xff && head[1] === 0xd8) {
    for (let i = 2; i < head.length - 9; i++) {
      if (head[i] !== 0xff) continue
      const marker = head[i + 1]
      const len = head.readUInt16BE(i + 2)
      if (marker >= 0xc0 && marker <= 0xcf && ![0xc4, 0xc8, 0xcc].includes(marker)) {
        return { height: head.readUInt16BE(i + 5), width: head.readUInt16BE(i + 7) }
      }
      i += len
    }
  }

  // PNG: IHDR is always the first chunk, at a fixed offset.
  if (head.readUInt32BE(0) === 0x89504e47) {
    return { width: head.readUInt32BE(16), height: head.readUInt32BE(20) }
  }

  throw new Error(`unrecognised image header in ${abs}`)
}

const abs = path.join(APP_ROOT, 'public', LOGIN_HERO.file.replace(/^\//, ''))

describe('LOGIN_HERO', () => {
  it('names a file that exists', () => {
    expect(fs.existsSync(abs), `${LOGIN_HERO.file} is declared but not on disk`).toBe(true)
  })

  it('declares the real dimensions of the file', () => {
    const { width, height } = dimensionsOf(abs)
    expect(LOGIN_HERO.width, `declares width ${LOGIN_HERO.width}, file is ${width}`).toBe(width)
    expect(LOGIN_HERO.height, `declares height ${LOGIN_HERO.height}, file is ${height}`).toBe(height)
  })
})