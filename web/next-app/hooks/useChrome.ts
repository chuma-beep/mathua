'use client'

import { useSyncExternalStore } from 'react'
import {
  subscribeChrome,
  getChromeVisible,
  getChromeServerVisible,
  toggleChrome,
} from '../lib/chrome'

/**
 * Whether the app chrome (sticky header + mobile tab bar) is showing, and a
 * toggle for it. Both navs call this so they never disagree about whether they
 * are on screen -- a header that hid while the tab bar stayed would look like
 * a bug rather than a gesture.
 */
export function useChrome(): { visible: boolean; toggle: () => void } {
  const visible = useSyncExternalStore(
    subscribeChrome,
    getChromeVisible,
    getChromeServerVisible,
  )
  return { visible, toggle: toggleChrome }
}
