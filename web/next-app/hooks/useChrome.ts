'use client'

import { useSyncExternalStore } from 'react'
import {
  subscribeChrome,
  getChromeVisible,
  getChromeServerVisible,
  getChromeHasChrome,
  getChromeHasTabs,
  getChromeServerHasChrome,
  getChromeServerHasTabs,
  toggleChrome,
} from '../lib/chrome'

/**
 * Whether the app chrome (sticky header + mobile tab bar) is showing, a toggle
 * for it, and which parts are mounted.
 *
 * Both navs call this so they never disagree about whether they are on screen
 * -- a header that hid while the tab bar stayed would look like a bug rather
 * than a gesture. The tap control uses `hasChrome` to stay out of the way on
 * pages that have no chrome to hide, and `hasTabs` to know whether to sit above
 * a tab bar.
 */
export function useChrome(): {
  visible: boolean
  toggle: () => void
  hasChrome: boolean
  hasTabs: boolean
} {
  const visible = useSyncExternalStore(
    subscribeChrome,
    getChromeVisible,
    getChromeServerVisible,
  )
  const hasChrome = useSyncExternalStore(
    subscribeChrome,
    getChromeHasChrome,
    getChromeServerHasChrome,
  )
  const hasTabs = useSyncExternalStore(
    subscribeChrome,
    getChromeHasTabs,
    getChromeServerHasTabs,
  )
  return { visible, toggle: toggleChrome, hasChrome, hasTabs }
}
