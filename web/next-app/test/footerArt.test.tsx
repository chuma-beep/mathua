import { describe, it, expect } from 'vitest'
import fs from 'node:fs'
import path from 'node:path'
import { FOOTER_ART } from '../components/Footer'

import { APP_ROOT } from './helpers/roots'

// The footer draws one of FOOTER_ART on every load, and each entry names two
// pictures: a landscape for desktop and a portrait for phones. The image is an
// <img> at `height: auto`, so the picture sets the footer's height, and the
// `width`/`height` on the element give the browser the intrinsic ratio before the
// file arrives.
//
// That makes a wrong `width`/`height` the failure worth guarding. It is not an
// exception: the image still renders, the footer still looks plausible, and the
// only symptom is a layout jump on a slow connection while the browser reflows
// once the real dimensions arrive. Nothing else in the suite would notice.
//
// The second failure is a file that is listed but absent, which 404s on exactly
// the visits that pick it and leaves the footer empty -- intermittent by
// construction, because the pick is random. With art direction there are twice
// as many files, so twice as many ways to get this wrong.

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

const variants = FOOTER_ART.flatMap(art => [
  { ...art.desktop, slot: 'desktop' as const },
  { ...art.mobile, slot: 'mobile' as const },
])

describe('FOOTER_ART', () => {
  it('is not empty', () => {
    // The sentinel. Every assertion below passes vacuously against an empty pool,
    // which is the shape of a guard that reads nothing.
    expect(FOOTER_ART.length).toBeGreaterThan(0)
    expect(variants.length).toBe(FOOTER_ART.length * 2)
  })

  it('names only files that exist', () => {
    for (const art of variants) {
      const abs = path.join(APP_ROOT, 'public', art.file.replace(/^\//, ''))
      expect(fs.existsSync(abs), `${art.slot} ${art.file} is listed but not on disk`).toBe(true)
    }
  })

  it('declares the real dimensions of every file', () => {
    for (const art of variants) {
      const abs = path.join(APP_ROOT, 'public', art.file.replace(/^\//, ''))
      const { width, height } = dimensionsOf(abs)
      expect(art.width, `${art.file} declares width ${art.width}, file is ${width}`).toBe(width)
      expect(art.height, `${art.file} declares height ${art.height}, file is ${height}`).toBe(height)
    }
  })

  it('pairs a landscape desktop with a portrait mobile', () => {
    // The whole premise of the art direction. A landscape shown across a phone is
    // a letterbox slot; a portrait shown across a desktop is a thin strip with
    // dead space either side. Swapping them would still render, and would still
    // look like a deliberate choice -- which is what makes it worth asserting.
    for (const art of FOOTER_ART) {
      const desktop = art.desktop.width / art.desktop.height
      const mobile = art.mobile.width / art.mobile.height
      expect(desktop, `${art.desktop.file} is not a landscape`).toBeGreaterThanOrEqual(1.6)
      expect(mobile, `${art.mobile.file} is not a portrait`).toBeLessThanOrEqual(0.8)
    }
  })

  it('keeps one desktop aspect ratio across the pool', () => {
    // The pick is random on every load. Desktop artwork of mixed ratios would
    // change the footer's height from visit to visit -- the page reflows under the
    // reader, and nothing about it looks like a bug at the moment it happens.
    // Mobile is deliberately exempt: the portrait set ranges from 9:16 to 2:3.
    const ratios = FOOTER_ART.map(a => a.desktop.width / a.desktop.height)
    const spread = Math.max(...ratios) - Math.min(...ratios)
    expect(spread, `ratios range ${Math.min(...ratios).toFixed(3)}..${Math.max(...ratios).toFixed(3)}`).toBeLessThan(0.05)
  })

  it('states a tone for every picture', () => {
    for (const art of variants) {
      expect(['light', 'dark'], `${art.slot} ${art.file}`).toContain(art.tone)
    }
  })

  it('gives every pool entry a tone for both breakpoints', () => {
    // The ink is read per breakpoint, so an entry missing one leaves the footer
    // on the other variant's tone -- and the failure is invisible until that
    // breakpoint is loaded.
    for (const art of FOOTER_ART) {
      expect(['light', 'dark']).toContain(art.desktop.tone)
      expect(['light', 'dark']).toContain(art.mobile.tone)
    }
  })
})