'use client'

import { useEffect, type RefObject } from 'react'

/**
 * Take an element out of the tab order and the accessibility tree while it is
 * hidden.
 *
 * The hidden chrome is only transformed and faded, so without this its links
 * stay focusable and screen-reader visible -- a keyboard user tabs into a nav
 * they cannot see. React 18 does not render `inert` (it only learned the
 * attribute in 19), so it is set on the node directly.
 *
 * `aria-hidden` is set alongside it: `inert` implies it where supported, but
 * assistive tech that predates `inert` still needs the explicit signal.
 */
export function useInertWhen<T extends HTMLElement>(
  ref: RefObject<T | null>,
  inert: boolean,
): void {
  useEffect(() => {
    const el = ref.current
    if (!el) return
    if (inert) {
      el.setAttribute('inert', '')
      el.setAttribute('aria-hidden', 'true')
    } else {
      el.removeAttribute('inert')
      el.removeAttribute('aria-hidden')
    }
  }, [ref, inert])
}
