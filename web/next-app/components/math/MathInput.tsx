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

import { Component, lazy, Suspense, useEffect, useRef, useState, type MutableRefObject, type ReactNode } from 'react'
import { Input } from '@/components/ui/input'
import SymbolPalette from '@/components/SymbolPalette'
import { presentationForQuestion } from './keyboards'
import type { MathInputStatus, MathInputValue } from './types'
import type { MathLiveFieldHandle } from './MathLiveField'

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

  const isDisabled = disabled === true || status === 'disabled'
  const shell = statusClasses(status, isDisabled)

  if (focusRef) {
    focusRef.current = {
      focus: () => {
        if (fallbackActive) {
          fallbackRef.current?.focus()
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
    <Input
      ref={fallbackRef}
      id={id}
      value={value ?? ''}
      onChange={e => onChange?.({ latex: '', plainAnswer: e.target.value })}
      onKeyDown={e => {
        if (e.key === 'Enter') {
          e.preventDefault()
          onSubmit?.()
        }
      }}
      placeholder={placeholder}
      disabled={isDisabled}
      aria-label={ariaLabel}
      className={`${shell} ${className ?? ''}`}
    />
  )

  return (
    <MathLiveBoundary fallback={fallback} onFallback={() => setFallbackActive(true)}>
      <div className={`mathua-math-input ${shell} ${className ?? ''}`}>
        <Suspense fallback={<div className="h-12" aria-hidden />}>
          <MathLiveField
            value={value ?? ''}
            mode={mode}
            disabled={isDisabled}
            placeholder={placeholder}
            ariaLabel={ariaLabel}
            onChange={v => onChange?.(v)}
            onSubmit={() => onSubmit?.()}
            handleRef={fieldRef}
            focusRequest={focusRequest}
          />
        </Suspense>
      </div>
    </MathLiveBoundary>
  )
}

// Mathua's states, expressed with the same tokens as the rest of the UI: square
// corners, hairline borders, no ring. The colours are CSS variables MathLive's own
// stylesheet also reads (see app/globals.css), so the field follows light/dark
// without a second palette.
function statusClasses(status: MathInputStatus, disabled: boolean): string {
  const base =
    'w-full min-w-0 min-h-12 border bg-mathua-code px-3 py-2 font-mono text-base text-mathua-primary transition-colors'
  if (disabled) return `${base} border-mathua-border opacity-50 cursor-not-allowed`
  switch (status) {
    case 'correct':
      return `${base} border-green-600/50`
    case 'incorrect':
      return `${base} border-red-600/50`
    default:
      return `${base} border-mathua-border focus-within:border-mathua-blue`
  }
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
  if (status === 'correct') return 'border-green-600/50'
  if (status === 'incorrect') return 'border-red-600/50'
  return ''
}

export { EMPTY as EMPTY_MATH_INPUT_VALUE }