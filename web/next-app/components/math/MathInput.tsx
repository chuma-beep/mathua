'use client'

// Mathua's mathematical answer editor.
//
// This is the only component the rest of the app imports for math input. It owns
// the visual states, the dynamic-import boundary and the fallback; it knows nothing
// about MathLive's element or about grading.
//
//   MathInput.tsx            ← every answer host imports this
//     └─ next/dynamic(ssr:false)
//        └─ MathLiveField.tsx   ← the only file importing 'mathlive'
//
// Why a fallback at all: MathLive is ~800 kB and registers a custom element that
// needs real text measurement. If that chunk fails to load — old browser, blocked
// asset, flaky connection — the learner must still be able to answer. So a failure
// degrades to the plain `<Input>` Mathua already had, which grades identically for
// every concept. That is the whole point of the fallback: it is not a lesser
// editor, it is the editor that always works.

import dynamic from 'next/dynamic'
import { Component, lazy, Suspense, useEffect, useRef, useState, type MutableRefObject, type ReactNode } from 'react'
import { Input } from '@/components/ui/input'
import SymbolPalette from '@/components/SymbolPalette'
import { presentationForQuestion } from './keyboards'
import type { MathInputStatus, MathInputValue } from './types'
import type { MathLiveFieldHandle } from './MathLiveField'

const SelfCheck = dynamic(() => import('./SelfCheck'), { ssr: false })

const EMPTY: MathInputValue = { latex: '', plainAnswer: '' }

// `React.lazy` rather than `next/dynamic({ ssr: false })`.
//
// Both give the same production result — MathLive lands in its own chunk that is
// never fetched during server rendering or on a route without an answer field.
// But `next/dynamic` resolves through Next's build-time chunk loader, which does
// not exist under vitest, so the loading state there never resolves and the
// component cannot be tested at all. `lazy` is plain React: it works in both, and
// the failure path below can be reached and asserted.
const MathLiveField = lazy(() => import('./MathLiveField'))

export interface MathFocusHandle {
  focus: () => void
}

export type { MathInputValue, MathInputStatus } from './types'
export type { MathKeyboardMode } from './keyboards'

interface Props {
  /** LaTeX. The wrapper's source of truth for the math field's content. */
  value?: string
  onChange?: (value: MathInputValue) => void
  onSubmit?: () => void
  disabled?: boolean
  placeholder?: string
  mode?: import('./keyboards').MathKeyboardMode
  status?: MathInputStatus
  /** Accessible name. Hosts already pass "Your answer" as a label. */
  ariaLabel?: string
  id?: string
  /** Focus handle, so a host can focus the editor without knowing which control rendered. */
  focusRef?: MutableRefObject<MathFocusHandle | null>
  className?: string
}

/**
 * Catches a MathLive chunk that fails to load or throws while mounting, and renders
 * the plain input instead.
 *
 * This is the only mechanism that can catch it. A chunk that fails to *load* never
 * mounts, so no ref check would ever fire, and `next/dynamic` surfaces the failure as
 * a React render error rather than a rejected promise. A boundary is the one thing
 * that sees both.
 */
class MathLiveBoundary extends Component<
  { fallback: ReactNode; onFallback: () => void; children: ReactNode },
  { failed: boolean }
> {
  state = { failed: false }

  static getDerivedStateFromError() {
    return { failed: true }
  }

  componentDidCatch(error: Error) {
    // Surfaced deliberately: a silent 800 kB failure looks exactly like a broken
    // deploy, and whoever sees this next needs to know which one it was.
    console.error('mathlive failed to load; falling back to the plain answer input:', error)
  }

  render() {
    return this.state.failed ? this.props.fallback : this.props.children
  }
}

/**
 * The plain answer input, used in two places.
 *
 * Once when the MathLive chunk fails and Mathua has no choice, and once while that
 * chunk is still *loading*. The second case used to render
 * `<div className="h-12" aria-hidden />` — a div of exactly the right size in exactly
 * the right place, which is the worst kind of loading state, because it looks like an
 * answer field and silently swallows the tap.
 *
 * Measured (`e2e/answer-diagnosis.spec.ts`): between the answer UI appearing and the
 * editor being usable there is a window of 322ms with the chunk cached and 3116ms
 * with a 3s chunk delay — which on a real phone connection is longer still. During
 * that window a tap on the answer area landed on that div, the keystrokes went
 * nowhere, and Check Answer was correctly but unhelpfully disabled because no answer
 * existed. Value propagation was never at fault: once the editor was up, input
 * reached React state in 2-7ms and the first tap submitted in 12-14ms.
 *
 * So this is not a cosmetic placeholder. It is the editor that always works, shown
 * while the better one loads — and whatever is typed here carries across, because the
 * typed text is reported as the LaTeX too and `MathLiveField` pushes it in on mount.
 */
function PlainAnswerInput({
  value,
  onChange,
  onSubmit,
  disabled,
  placeholder,
  ariaLabel,
  id,
  className,
  inputRef,
  loading,
}: {
  value: string
  onChange: (plainAnswer: string) => void
  onSubmit?: () => void
  disabled: boolean
  placeholder?: string
  ariaLabel: string
  id?: string
  className: string
  inputRef?: MutableRefObject<HTMLInputElement | null>
  /**
   * Marks the loading instance. The attribute is not decoration: the editor and the
   * loading input share a placeholder, so without it a query for "the answer field"
   * can return the transient control and a tap on it goes nowhere — which is the
   * exact failure this component was changed to fix, reappearing in the tests.
   */
  loading?: boolean
}) {
  return (
    <Input
      ref={inputRef}
      id={id}
      value={value}
      data-mathua-loading={loading ? 'true' : undefined}
      aria-busy={loading ? true : undefined}
      onChange={e => onChange(e.target.value)}
      onKeyDown={e => {
        if (e.key === 'Enter') {
          e.preventDefault()
          onSubmit?.()
        }
      }}
      placeholder={placeholder}
      enterKeyHint="go"
      disabled={disabled}
      aria-label={ariaLabel}
      className={className}
    />
  )
}

/**
 * Reports the moment the boundary trips, so `focus()` can point at the plain input
 * instead of at an element that never mounted.
 */
function FallbackSlot({ children, onFallback }: { children: ReactNode; onFallback: () => void }) {
  useEffect(() => {
    onFallback()
  }, [onFallback])
  return <>{children}</>
}

/**
 * Once MathLive has mounted we stay on it: remounting mid-answer would discard what
 * the learner typed. A failed mount is the only way to reach the fallback, and that
 * is a terminal state for this component instance.
 */
export default function MathInput({
  value,
  onChange,
  onSubmit,
  disabled,
  placeholder,
  mode = 'algebra',
  status = 'default',
  ariaLabel = 'Your answer',
  id,
  focusRef,
  className,
}: Props) {
  const fieldRef = useRef<MathLiveFieldHandle | null>(null)
  const fallbackRef = useRef<HTMLInputElement | null>(null)
  const [fallbackActive, setFallbackActive] = useState(false)
  const [focusRequest, setFocusRequest] = useState(0)
  // What the learner actually typed, for the optional self-check. Tracked separately
  // from `value` because a host's value is the *plain* answer it will submit, and it
  // only catches up after a render — the hint must not wait for the round trip.
  const [typedLatex, setTypedLatex] = useState('')
  // What has been typed into the plain input while MathLive loads. Local, because the
  // host's `value` is LaTeX and this control is showing plain text; the two only meet
  // when `MathLiveField` mounts and pushes `value` in.
  const [loadingText, setLoadingText] = useState('')

  const isDisabled = disabled === true || status === 'disabled'
  const layout = layoutClasses()
  // Only the standalone fallback needs a status *class*; see fallbackStatusClasses.
  const fallbackShell = fallbackStatusClasses(status, isDisabled)

  if (focusRef) {
    focusRef.current = {
      focus: () => {
        // Whichever plain input is mounted wins: the Suspense one while the chunk
        // loads, the boundary one if it failed. Testing the ref rather than
        // `fallbackActive` is what makes the loading case focusable at all.
        if (fallbackRef.current) {
          fallbackRef.current.focus()
          return
        }
        // Queue the request rather than dropping it: on the first question the host
        // autofocuses while this chunk is still loading, and a plain call would find
        // no element and silently do nothing.
        setFocusRequest(n => n + 1)
        fieldRef.current?.focus()
      },
    }
  }

  const fallback = (
    <PlainAnswerInput
      inputRef={fallbackRef}
      id={id}
      value={value ?? ''}
      onChange={plain => onChange?.({ latex: plain, plainAnswer: plain })}
      onSubmit={() => onSubmit?.()}
      disabled={isDisabled}
      placeholder={placeholder}
      ariaLabel={ariaLabel}
      className={`w-full min-w-0 min-h-12 bg-mathua-code font-mono text-base ${fallbackShell} ${className ?? ''}`}
    />
  )

  return (
    <MathLiveBoundary fallback={fallback} onFallback={() => setFallbackActive(true)}>
      <div className={`mathua-math-input ${layout} ${className ?? ''}`}>
        {/* The loading state is a working plain input, not a spacer. See
            PlainAnswerInput: a div of the right size in the right place looked like an
            answer field and swallowed every tap and keystroke for the length of the
            chunk load. Typed text is reported as the LaTeX too, so `MathLiveField`
            picks it up when it mounts. */}
        <Suspense
          fallback={
            <PlainAnswerInput
              inputRef={fallbackRef}
              id={id}
              value={loadingText}
              onChange={plain => {
                setLoadingText(plain)
                onChange?.({ latex: plain, plainAnswer: plain })
              }}
              onSubmit={() => onSubmit?.()}
              disabled={isDisabled}
              placeholder={placeholder}
              ariaLabel={ariaLabel}
              className={`w-full ${className ?? ''}`}
              loading
            />
          }
        >
          <MathLiveField
            value={value ?? ''}
            mode={mode}
            disabled={isDisabled}
            status={status}
            placeholder={placeholder}
            ariaLabel={ariaLabel}
            onChange={v => onChange?.(v)}
            onLatex={setTypedLatex}
            onSubmit={() => onSubmit?.()}
            handleRef={fieldRef}
            focusRequest={focusRequest}
          />
        </Suspense>
        {/* Its own boundary, deliberately outside the field's. A lazily-loaded
            sibling suspending inside the field's <Suspense> re-suspends that boundary
            and React discards the already-rendered field — remounting it, which reset
            the answer to the host's last value and lost what the learner had typed. */}
        <Suspense fallback={null}>
          <SelfCheck latex={typedLatex} active={!isDisabled} />
        </Suspense>
      </div>
    </MathLiveBoundary>
  )
}

// Mathua's states, expressed with the same tokens as the rest of the UI: square
// corners, hairline borders, no ring. The colours are CSS variables MathLive's own
// stylesheet also reads (see app/globals.css), so the field follows light/dark
// without a second palette.
//
// Split in two, because the box belongs to the control and not to the wrapper.
//
// These two used to be one string applied to the wrapper `<div>`, and that produced
// two nested bordered rectangles around the answer field: the wrapper drew
// `border border-mathua-border bg-mathua-code px-3 py-2 focus-within:border-mathua-blue`
// while `<math-field>` drew its own `border: 1px solid var(--border)` and turned blue
// on `:focus-within` (app/globals.css). Same blue — `mathua-blue` *is*
// `var(--accent-blue)`, tailwind.config.js — so it read as one colour in two states
// rather than as a deliberate two-tone treatment. Both transition over 0.15s, so they
// lit up together.
//
// Worse, the wrapper also wraps `SelfCheck`, so once anything was typed the outer box
// grew downward to enclose the hint. And the verdict colours landed on the wrapper,
// *outside* the field's opaque background, so `correct` drew a coloured ring around a
// still-grey box instead of colouring the box the learner was looking at.
function layoutClasses(): string {
  // `min-h-12` and the padding moved to the control. `app/globals.css` sets the
  // field's own `min-height: 3rem` and padding to match `Input`, because
  // `e2e/mobile.spec.ts` asserts the field is at least 48px and it only reached that
  // height via this wrapper.
  return 'w-full min-w-0'
}

// The plain `<input>` fallback draws its own box, so it still needs the status colour
// as a class. The math editor cannot use one: globals.css's
// `math-field.mathua-math-field` selector is more specific than any Tailwind utility
// class, so a `border-mathua-green` passed to `<math-field>` would lose to the
// `border-color` declared there. `MathLiveField` therefore takes `status` and sets
// `data-status`, which globals.css styles.
function fallbackStatusClasses(status: MathInputStatus, disabled: boolean): string {
  if (disabled) return 'border-mathua-border opacity-50 cursor-not-allowed'
  return legacyStatusClass(status)
}

// ---------------------------------------------------------------------------
// MathAnswerInput — what the answer hosts actually render.
//
// One component, two controls behind it, chosen from the question's own metadata.
// The important property is what the host does *not* have to know: it keeps a plain
// string in state, calls `onChange` with a plain string, and submits that string to
// the same endpoint it always did. LaTeX, MathJSON and MathLive all stay on this
// side of the line, which is why none of the five answer hosts needed to change
// anything but this one import.
// ---------------------------------------------------------------------------

export interface MathAnswerInputProps {
  /** The answer, as the plain string the graders have always received. */
  value: string
  onChange: (plainAnswer: string) => void
  onSubmit?: () => void
  disabled?: boolean
  gradingType?: string
  conceptId?: string
  status?: MathInputStatus
  id?: string
  ariaLabel?: string
  placeholder?: string
  className?: string
  /** Native mobile keyboard for text mode. Ignored by the math editor, which has its own keypad. */
  inputMode?: 'numeric' | 'text'
  focusRef?: MutableRefObject<MathFocusHandle | null>
  /** Render the legacy SymbolPalette under a text-mode input. Off in math mode, where the keypad replaces it. */
  showPalette?: boolean
}

export function MathAnswerInput({
  value,
  onChange,
  onSubmit,
  disabled,
  gradingType,
  conceptId,
  status = 'default',
  id,
  ariaLabel = 'Your answer',
  placeholder,
  className,
  inputMode,
  focusRef,
  showPalette = true,
}: MathAnswerInputProps) {
  const { editor, keyboard } = presentationForQuestion(gradingType, conceptId)

  // The plain input and the SymbolPalette share one element ref, owned here rather
  // than by the host: in math mode there is no input element to point at, and a
  // host holding a null element ref cannot focus what it does not know about.
  const inputRef = useRef<HTMLInputElement | null>(null)

  // MathLive's source of truth is LaTeX, the host's is a plain string. Hold the
  // LaTeX here so typing does not round-trip the answer through a parser on every
  // keystroke, and let the field push plain text upward.
  const latexRef = useRef(value)

  if (editor === 'text') {
    if (focusRef) focusRef.current = { focus: () => inputRef.current?.focus() }
    return (
      <div className={className}>
        <Input
          ref={inputRef}
          id={id}
          value={value}
          onChange={e => onChange(e.target.value)}
          onKeyDown={e => {
            if (e.key === 'Enter') {
              e.preventDefault()
              onSubmit?.()
            }
          }}
          placeholder={placeholder}
          inputMode={inputMode}
          enterKeyHint="go"
          disabled={disabled}
          aria-label={ariaLabel}
          className={status === 'default' ? '' : legacyStatusClass(status)}
        />
        {showPalette ? <SymbolPalette targetRef={inputRef} onInsert={onChange} /> : null}
      </div>
    )
  }

  return (
    <MathInput
      value={latexRef.current}
      mode={keyboard}
      status={status}
      disabled={disabled}
      placeholder={placeholder}
      ariaLabel={ariaLabel}
      id={id}
      className={className}
      focusRef={focusRef}
      onSubmit={() => onSubmit?.()}
      onChange={v => {
        latexRef.current = v.latex
        onChange(v.plainAnswer)
      }}
    />
  )
}

function legacyStatusClass(status: MathInputStatus): string {
  if (status === 'correct') return 'border-mathua-green-faint'
  if (status === 'incorrect') return 'border-mathua-red-faint'
  return ''
}

export { EMPTY as EMPTY_MATH_INPUT_VALUE }