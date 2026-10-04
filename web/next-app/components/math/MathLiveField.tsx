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
// A value import, not just the type: the two settings below are statics on the class,
// and the instance accessors deliberately throw ("Use MathfieldElement.soundsDirectory
// instead"), so there is no way to reach them from an element.
import { MathfieldElement } from 'mathlive'
import { layoutsForMode } from './layouts'
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
  /** Called with the LaTeX as typed, for optional decoration. Never for grading. */
  onLatex?: (latex: string) => void
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
  onLatex,
  onSubmit,
  handleRef,
  focusRequest = 0,
}: Props) {
  // Per-keystroke feedback Mathua does not want, configured off once per document.
  //
  // Verified against 0.111.0 rather than assumed: `MathfieldElement._soundsDirectory`
  // defaults to `'./sounds'` and `keypressSound` to `keypress-standard.wav` and friends,
  // loaded with `fetch(`${dir}/${file}`)` on every key. Mathua copies MathLive's
  // *fonts* into `public/fonts` but never shipped `public/sounds`, so every keystroke
  // was issuing a request for a file that does not exist. Setting the static to `null`
  // is the documented way to prevent any sound from being loaded, and `keypressVibration`
  // is the same idea for the vibration API.
  useEffect(() => {
    MathfieldElement.soundsDirectory = null
    MathfieldElement.keypressVibration = false
  }, [])

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

  // Stable on purpose: hosts pass a fresh arrow to `onChange` on every render, and an
  // `emit` that changed identity re-ran the sync effect below, which then wrote the
  // host's (stale) value back over what the learner had just typed.
  const onChangeRef = useRef(onChange)
  onChangeRef.current = onChange
  const onLatexRef = useRef(onLatex)
  onLatexRef.current = onLatex
  const emit = useCallback(() => {
    const el = ref.current
    if (!el) return
    const latex = el.value ?? ''
    lastPushed.current = latex
    onChangeRef.current({ latex, plainAnswer: toPlainAnswer(el.getValue('plain-text') ?? '') })
    onLatexRef.current?.(latex)
  }, [])

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
    // `value` only. See the note on `emit`.
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

  // What a spacebar press inserts. `mathModeSpace` defaults to `""`, and with it empty
  // the `[Space]` keystroke handler only calls `moveAfterParent()` -- the caret moves
  // and nothing is inserted. That is why the spacebar looked dead: MathLive's own
  // `alphabetic` layout ships `{ label: " ", width: 1.5 }`, which declares no command,
  // no insert and no latex, so `executeKeycapCommand` falls through to
  // `typedText(" ")` and lands on that empty branch. It is an upstream defect, not a
  // Mathua regression.
  //
  // `\,` and not a literal space: TeX discards literal spaces in math mode, so a spacebar
  // inserting `" "` would render and serialise as nothing. `\,` is the payload proven to
  // survive the round trip -- the recorded corpus holds `expression 5 R 3 <- 5\ R\ 3`,
  // which comes back as the plain text `5 R 3` and grades. Three generator answer formats
  // need it: `Q R`, `a sqrt(b)` and `P=[[...]] D=[[...]]`.
  //
  // Set as a property, not as the documented `math-mode-space` attribute. Measured on
  // 0.111.0 with React 18: the attribute is present in the DOM and `placeholder` and
  // `class` are both honoured from it, but `mathModeSpace` stays `""` and the spacebar
  // stays inert. The property and `setOptions({ mathModeSpace })` both work. An attribute
  // that is silently ignored is worse than none -- it reads as configuration and is not.
  useEffect(() => {
    const el = ref.current
    if (el) el.mathModeSpace = '\\,'
  }, [])

  const installKeyboard = useCallback(() => {
    const kb = typeof window === 'undefined' ? undefined : window.mathVirtualKeyboard
    if (!kb) return
    // The layout is a page-wide singleton, so every field installs its own on focus
    // or they overwrite each other in the Learn feed.
    if (installedMode.current !== mode) {
      // One MathLive *layout* per Mathua page, plus MathLive's own `alphabetic`, which
      // gives every letter
      // without a 26-key grid Mathua would have to design for a 320px screen and fail
      // the 44px touch-target floor on. With more than one entry in `kb.layouts`,
      // MathLive renders its own layout switcher in the keypad toolbar -- verified, not
      // assumed: the toolbar read "Arithmeticabc" before this was wired up, because a
      // probe set `layouts` by hand and MathLive offered the switcher unprompted.
      //
      // This is what makes a yes/no question typeable. Multiple choice still renders
      // `ChoiceOptions` above the field, so `yes` is a tap; the alphabet is there for the
      // learner who would rather type, and `addition`, `odd` and `Z` have no other route
      // in.
      kb.layouts = [...layoutsForMode(mode), 'alphabetic']
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
    return () => el.removeEventListener('keydown', onKeyDown)
  }, [])

  // Raise the keypad when the field is touched.
  //
  // Bound natively on the element rather than through React's `onFocus`: MathLive
  // focuses a contenteditable sink inside its own shadow root, and a tap on a touch
  // profile does not reliably surface as a focus event on the host. `pointerdown`
  // always does.
  //
  // This is also the recovery path. MathLive raises the keyboard from exactly one
  // place — a document `focusin` listener — so if anything hides the panel while the
  // field still holds focus, tapping the field again is the only thing that can bring
  // it back, and that is this handler. MathLive's own `focusout` handler hides the
  // panel when focus genuinely leaves, which is the behaviour we want; nothing else
  // needs to.
  useEffect(() => {
    const el = ref.current
    if (!el) return
    const onPointerDown = () => installKeyboardRef.current()
    el.addEventListener('pointerdown', onPointerDown)
    return () => el.removeEventListener('pointerdown', onPointerDown)
  }, [])

  // Publish the keypad's height so page layout can keep clear of it.
  //
  // MathLive positions its own panel relative to the field, so it does not cover the
  // editor -- that part is its job. What it cannot do is shrink the *page*: a sticky
  // header, or a submit button anchored to the bottom of the viewport, still sits
  // underneath a panel whose height Mathua has no idea about. `geometrychange` carries
  // the new bounding rectangle, and this turns it into one CSS custom property so
  // stylesheets can adapt without a hardcoded pixel value that is wrong on every
  // device.
  useEffect(() => {
    if (typeof window === 'undefined') return
    const kb = window.mathVirtualKeyboard
    if (!kb) return
    // The rect is read from `kb.boundingRect`, *not* from the event.
    //
    // 0.111.0's own type documentation says `evt.detail.boundingRect`, and that is
    // wrong: the library dispatches `new Event('geometrychange')` — a plain Event with
    // no `detail` at all — from a ResizeObserver and from `stateChanged()`. A handler
    // written to the documented shape reads `undefined`, computes a height of zero, and
    // reports the keypad as closed forever, which is precisely the failure this exists to
    // prevent. The measured symptom of getting it wrong is in git history.
    const publish = () => {
      const height = Math.round(kb.boundingRect?.height ?? 0)
      document.documentElement.style.setProperty('--mathua-keyboard-height', `${height}px`)
      document.documentElement.dataset.mathuaKeyboard = height > 0 ? 'open' : 'closed'
    }
    const onGeometryChange = () => publish()
    kb.addEventListener('geometrychange', onGeometryChange)
    // The panel can already be up when this mounts (a later card in the Learn feed),
    // so publish the current geometry rather than waiting for the first change.
    publish()
    return () => {
      kb.removeEventListener('geometrychange', onGeometryChange)
      document.documentElement.style.removeProperty('--mathua-keyboard-height')
      delete document.documentElement.dataset.mathuaKeyboard
    }
  }, [])

  // A Learn card can unmount while focused (a 3-miss halt, a concept swap).
  // Leaving the singleton keyboard up would cover the next exercise.
  //
  // This is the only place Mathua hides the keypad. There used to be a
  // capture-phase `pointerdown` listener on `document` that hid it on any tap
  // outside the field, and it had to go: hiding without requiring focus to leave
  // left the field focused with no keyboard and no way back, because MathLive only
  // shows on a *focus transition*. Tapping a non-focusable area — the question
  // text, the page background — did exactly that, and the learner had to tap the
  // field several times before it came back. MathLive's own `focusout` handler
  // covers the case that should actually dismiss.
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