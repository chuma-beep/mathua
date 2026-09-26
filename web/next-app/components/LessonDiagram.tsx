'use client'

import { useState, useEffect } from 'react'
import Image from 'next/image'
import DOMPurify from 'isomorphic-dompurify'

interface LessonDiagramProps {
  src: string
  alt: string
  maxHeight?: number
}

// SVG allowlist: shapes, text and containers only. Scripts, event-handler
// attributes, foreignObject and external references never survive —
// vendored third-party SVGs go through this before touching the DOM.
const SVG_PROFILE: Parameters<typeof DOMPurify.sanitize>[1] = {
  USE_PROFILES: { svg: true },
  FORBID_TAGS: ['foreignObject', 'script', 'style', 'a', 'image', 'use'],
  FORBID_ATTR: ['on*'],
}

function sanitizeSvg(raw: string, alt: string): string {
  const clean = DOMPurify.sanitize(raw, SVG_PROFILE)
  // Ensure a screen-reader name: keep an existing <title>, otherwise
  // inject one from the alt text so the inline graphic is never anonymous.
  if (/<title[\s>]/i.test(clean)) return clean
  const escaped = alt.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
  return clean.replace(/<svg([^>]*)>/i, `<svg$1 role="img"><title>${escaped}</title>`)
}

interface DiagramMeta {
  title?: string
  description?: string
  source?: string
}

// Per-asset metadata (data/diagrams/meta.json, published copy under
// public/diagrams/): titles upgrade the alt text, source records the
// CC BY-NC attribution for vendored Algebrica assets. Module-cached,
// one fetch per session at most.
let metaCache: Promise<Record<string, DiagramMeta>> | null = null
function loadMeta(): Promise<Record<string, DiagramMeta>> {
  if (!metaCache) {
    metaCache = fetch('/diagrams/meta.json')
      .then((res) => (res.ok ? res.json() : {}))
      .catch(() => ({}))
  }
  return metaCache
}

// Lesson diagram with the recolor policy applied: SVGs inline sanitized
// (themeable, crisp, accessible); PNGs keep the next/image path. A fetch
// failure falls back to the plain image so a missing asset never blanks.
export default function LessonDiagram({ src, alt, maxHeight = 180 }: LessonDiagramProps) {
  const isSvg = src.toLowerCase().split('?')[0].endsWith('.svg')
  const [label, setLabel] = useState(alt)
  const [inline, setInline] = useState<string | null>(null)

  useEffect(() => {
    setLabel(alt)
    setInline(null)
    let cancelled = false
    // Metadata first so the sanitized <title> and labels use the real
    // title when the asset has one; falls back to the alt prop silently.
    loadMeta()
      .catch(() => ({} as Record<string, DiagramMeta>))
      .then((meta) => {
        if (cancelled) return
        const title = meta[src]?.title
        const name = title || alt
        if (title) setLabel(title)
        if (!isSvg) return
        return fetch(src)
          .then((res) => {
            if (!res.ok) throw new Error(`diagram fetch failed: ${res.status}`)
            return res.text()
          })
          .then((raw) => {
            if (!cancelled) setInline(sanitizeSvg(raw, name))
          })
      })
      .catch(() => {})
    return () => {
      cancelled = true
    }
  }, [src, isSvg, alt])

  if (inline === null) {
    // PNG path, or SVG still loading / failed to fetch: plain image so a
    // missing or unloadable asset never blanks the diagram slot.
    return (
      <Image
        src={src}
        alt={label}
        width={220}
        height={180}
        className="max-w-full h-auto"
        style={{ maxHeight: `${maxHeight}px` }}
        unoptimized
      />
    )
  }
  return (
    <div
      className="lesson-diagram"
      role="img"
      aria-label={label}
      style={{ maxHeight: `${maxHeight}px` }}
      dangerouslySetInnerHTML={{ __html: inline }}
    />
  )
}
