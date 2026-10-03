'use client'

// The only file in Mathua that imports `mathlive`. Everything else goes through
// MathInput.tsx, so the dependency has exactly one edge to reason about.
//
// Why this module is separate from MathInput.tsx: MathLive is a custom element
// that needs a DOM and real text measurement (importing it registers
// `<math-field>` as a side effect). That makes it unsafe to evaluate during
// server rendering, so it must sit behind its own module boundary and be pulled
// in lazily. This is that boundary.
//
// Everything here is public MathLive API: own the element, choose the keyboard
// layout on focus, report the value upward, drive Enter to submit. Nothing here
// decides anything about grading or holds answer state.

import { useCallback, useEffect, useRef, type MutableRefObject } from 'react'
// Side-effect import: this is what registers the <math-field> custom element.
// The type import below is erased at runtime, so it cannot do this job.
import 'mathlive'
import type { MathfieldElement } from 'mathlive'
import { layoutForMode } from './layouts'
import type { MathKeyboardMode } from './keyboards'
import { toPlainAnswer } from './plainAnswer'
import type { MathInputValue } from './types'

/*
 * `<math-field>` is a custom element MathLive registers at import time. React 18
 * has no JSX type for it and MathLive's types declare the class but not the
 * intrinsic element, so this is the one place that bridges the two. It lives here
 * rather than in a standalone .d.ts because a declaration file has no runtime
 * presence and the repo's orphan-module guard treats that as dead code.
 */
declare global {
  namespace JSX {
    interface IntrinsicElements {
      'math-field': React.DetailedHTMLProps<React.HTMLAttributes<MathfieldElement>, MathfieldElement> & {
        /** MathLive attribute: LaTeX rendered while the field is empty. */
        placeholder?: string
        /** MathLive attribute: 'auto' | 'manual' | 'sandboxed'. */
        'math-virtual-keyboard-policy'?: string
        /**
         * `class`, not `className`. React 18 forwards props to custom elements
         * verbatim instead of mapping className to class, so className would be
         * written as the attribute `classname` and the theme CSS would never match.
         */
        class?: string
      }
    }
  }
}

export interface MathLiveFieldHandle {
  focus: () => void
}

interface Props {
  value: string
  mode: MathKeyboardMode
  disabled: boolean
  placeholder?: string
  /** Accessible name. MathLive gives the field `role="group"`, so this is what a screen reader announces. */
  ariaLabel: string
  onChange: (value: MathInputValue) => void
  onSubmit: () => void
  handleRef?: MutableRefObject<MathLiveFieldHandle | null>
  /**
   * Incremented by the wrapper each time focus is requested. A counter rather than
   * a boolean because the host asks for focus when a question is revealed, which on
   * the first question is *before* this chunk has resolved — a boolean that was
   * already true would not re-run the effect once the element finally exists.
   */
  focusRequest?: number
}

export default function MathLiveField({
  value,
  mode,
  disabled,
  placeholder,
  ariaLabel,
  onChange,
  onSubmit,
  handleRef,
  focusRequest = 0,
}: Props) {
  const ref = useRef<MathfieldElement | null>(null)
  // MathLive owns the element's value, so a value arriving from outside has to
  // be pushed in. Tracking what we last sent keeps us from resetting the
  // learner's cursor on every keystroke.
  const lastPushed = useRef<string>('')
  // The virtual keyboard is a page-wide singleton and the Learn feed can hold
  // several fields at once. Remember which layout we installed so a re-focus
  // doesn't stomp on one a sibling field just set.
  const installedMode = useRef<MathKeyboardMode | null>(null)
  // So the listeners, registered once, always call the current handler.
  const onSubmitRef = useRef(onSubmit)
  onSubmitRef.current = onSubmit
  // installKeyboard is declared below; the pointerdown listener needs it now.
  const installKeyboardRef = useRef<() => void>(() => {})

  const emit = useCallback(() => {
    const el = ref.current
    if (!el) return
    const latex = el.value ?? ''
    lastPushed.current = latex
    onChange({ latex, plainAnswer: toPlainAnswer(el.getValue('plain-text') ?? '') })
  }, [onChange])

  // Push an externally-owned value in — ChoiceOptions filling the field, or a
  // reused card. Only when it genuinely differs from what we last emitted.
  useEffect(() => {
    const el = ref.current
    if (!el || value === lastPushed.current) return
    lastPushed.current = value
    if (el.value !== value) {
      el.value = value
      emit()
    }
  }, [value, emit])

  useEffect(() => {
    if (focusRequest > 0) handleRef?.current?.focus()
  }, [focusRequest, handleRef])

  // `disabled` is a property on the element, not a reflected attribute we can
  // rely on across React versions.
  useEffect(() => {
    const el = ref.current
    if (el) el.readOnly = disabled
  }, [disabled])

  const installKeyboard = useCallback(() => {
    const kb = typeof window === 'undefined' ? undefined : window.mathVirtualKeyboard
    if (!kb) return
    // The layout is a page-wide singleton, so every field installs its own on focus
    // or they overwrite each other in the Learn feed.
    if (installedMode.current !== mode) {
      kb.layouts = [layoutForMode(mode)]
      installedMode.current = mode
    }
    // Show it only where there is no physical keyboard to type on. Opening a
    // 4-row keypad over the answer on every focus of a desktop browser is the
    // "giant calculator" the design is supposed to avoid; on desktop the toggle
    // stays available for the symbolic keys.
    if (window.matchMedia?.('(pointer: coarse)').matches) kb.show()
  }, [mode])

  installKeyboardRef.current = installKeyboard

  // Enter is bound with a native listener rather than a React onKeyDown prop.
  // Keystrokes originate in the contenteditable sink inside MathLive's shadow root;
  // they are composed and do reach the host, but React's delegated root listener does
  // not reliably see them, and this one must not be flaky.
  useEffect(() => {
    const el = ref.current
    if (!el) return
    const onKeyDown = (e: KeyboardEvent) => {
      if (e.key !== 'Enter') return
      // Deliberately no preventDefault(). MathLive's own Enter handling commits the
      // expression; suppressing it makes the next keystroke replay the field's
      // contents, so "2+2" + Enter produced "2+22+2". Recorded, not assumed.
      onSubmitRef.current?.()
    }
    el.addEventListener('keydown', onKeyDown)
    el.addEventListener('keydown', onKeyDown)
    return () => el.removeEventListener('keydown', onKeyDown)
  }, [])

  // Raise the keypad when the field is touched.
  //
  // Bound natively on the element rather than through React's `onFocus`: MathLive
  // focuses a contenteditable sink inside its own shadow root, and a tap on a touch
  // profile does not reliably surface as a focus event on the host. `pointerdown`
  // always does.
  useEffect(() => {
    const el = ref.current
    if (!el) return
    const onPointerDown = () => installKeyboardRef.current()
    el.addEventListener('pointerdown', onPointerDown)
    return () => el.removeEventListener('pointerdown', onPointerDown)
  }, [])

  // Dismiss the keypad on a touch outside it.
  //
  // Known limitation, measured rather than assumed: in MathLive 0.111.0 under
  // Chromium, `mathVirtualKeyboard.hide()` had no effect — calling it directly from
  // the page left both `visible === true` and the panel in the DOM, under `auto` and
  // under `manual` alike. So this listener expresses the intent and does not yet
  // achieve it, and the mobile suite does not assert dismissal. Tracked as the first
  // thing to revisit if the keypad is found covering an exercise.
  useEffect(() => {
    if (typeof document === 'undefined') return
    const onPointerDown = (e: PointerEvent) => {
      const target = e.target as Element | null
      if (!target) return
      if (target.closest?.('math-field')) return
      // The keypad's own classes are all MLK__-prefixed; there is no single
      // wrapper class to test for, and hiding it mid-tap cancels the keypress.
      if (target.closest?.('[class*="MLK__"]')) return
      window.mathVirtualKeyboard?.hide()
    }
    document.addEventListener('pointerdown', onPointerDown, true)
    return () => document.removeEventListener('pointerdown', onPointerDown, true)
  }, [])

  // A Learn card can unmount while focused (a 3-miss halt, a concept swap).
  // Leaving the singleton keyboard up would cover the next exercise.
  useEffect(() => {
    const el = ref.current
    return () => {
      if (typeof window === 'undefined') return
      if (document.activeElement === el) window.mathVirtualKeyboard?.hide()
    }
  }, [])

  if (handleRef) {
    handleRef.current = { focus: () => ref.current?.focus() }
  }

  return (
    <math-field
      ref={ref}
      class="mathua-math-field"
      placeholder={placeholder}
      aria-label={ariaLabel}
      onInput={emit}
      onFocus={installKeyboard}
    />
  )
}