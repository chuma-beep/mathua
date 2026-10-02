'use client'

import { ChevronDown, ChevronUp } from 'lucide-react'
import { useChrome } from '../hooks/useChrome'

// The chrome's tap control. It has to outlive the bars it controls: when they
// are off screen there is nothing left to tap, so both of these stay put and one
// of them is always the way back in.
//
// Deliberately not rendered inside Header. The header's backdrop-blur creates a
// containing block that would trap this position:fixed child, and the header
// goes inert while hidden, which would disable the control exactly when it is
// needed. The root layout mounts it instead.

// The tab bar sits at this offset from the bottom; the chevron goes just above
// it so it is never under the bar's links.
const BAR_OFFSET = 'max(12px, env(safe-area-inset-bottom))'
const BAR_HEIGHT = 56

export default function ChromeToggle() {
  const { visible, toggle, hasChrome, hasTabs } = useChrome()

  // Nothing to hide on a page without chrome (mobile /profile has neither a
  // header nor a tab bar), so a control here would do nothing.
  if (!hasChrome) return null

  // Above the tab bar when there is one and it is showing; at the bottom edge
  // otherwise, which is both the no-bar case and the hidden case.
  const chevronBottom = hasTabs && visible
    ? `calc(${BAR_OFFSET} + ${BAR_HEIGHT + 2}px)`
    : 'max(6px, env(safe-area-inset-bottom))'

  return (
    <>
      {/* Full-width bottom strip: the "tap the screen to bring it back" target.
          Live only while the bars are hidden. When they are showing this would
          sit under the tab bar and in front of whatever content is at the bottom
          of the page, so it is inert then and the chevron does the hiding. */}
      <div
        aria-hidden="true"
        data-testid="chrome-tap-strip"
        onClick={visible ? undefined : toggle}
        className={`lg:hidden fixed inset-x-0 z-30 ${visible ? 'pointer-events-none' : 'pointer-events-auto'}`}
        style={{ bottom: 0, height: 'calc(28px + env(safe-area-inset-bottom))' }}
      />

      <button
        type="button"
        onClick={toggle}
        aria-expanded={visible}
        aria-label={visible ? 'Hide navigation' : 'Show navigation'}
        data-testid="chrome-toggle"
        className="lg:hidden fixed left-1/2 z-40 flex size-11 -translate-x-1/2 items-center justify-center rounded-full border border-mathua-border bg-mathua-surface text-mathua-muted transition-[bottom] duration-300 ease-out active:text-mathua-blue motion-reduce:transition-none"
        style={{ bottom: chevronBottom }}
      >
        {visible ? (
          <ChevronDown size={18} strokeWidth={2} aria-hidden="true" />
        ) : (
          <ChevronUp size={18} strokeWidth={2} aria-hidden="true" />
        )}
      </button>
    </>
  )
}
