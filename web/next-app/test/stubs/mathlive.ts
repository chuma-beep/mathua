// A stand-in for `mathlive` under vitest.
//
// MathLive cannot run in jsdom: it is a custom element that measures real text
// (ResizeObserver, font metrics, MathML layout), and its `node` export condition
// resolves to the SSR build, which never registers the element at all. Every
// attempt to load it here ends in a blank field or a thrown selector error.
//
// Stubbing is therefore the honest option for the component suite: these tests are
// about Mathua's behaviour — which control is chosen, what gets submitted, what the
// learner sees — and none of that is MathLive's. What the real library does inside
// the field is verified by the Playwright suites (`e2e/math-input.spec.ts`,
// `e2e/mobile-math.spec.ts`), which run Chromium for real.
//
// The stub is deliberately thin but not a no-op: it round-trips `value`, fires
// `input` the way the real element does, and reports a plain-text form, so the
// submission path is exercised end to end.
//
// It declares nothing globally. `tsc` does not read vitest aliases, so it still
// sees MathLive's own `Window.mathVirtualKeyboard` augmentation; redeclaring it
// here would be a conflict rather than a stub.

/** Structural stand-in for the parts of MathLive's element the wrapper touches. */
export interface MathfieldElement extends HTMLElement {
  value: string
  readOnly: boolean
  mathVirtualKeyboardPolicy: string
  getValue(format?: string): string
  focus(): void
}

class StubMathField extends HTMLElement implements MathfieldElement {
  #value = ''
  #listeners = new Map<string, Set<(e: Event) => void>>()

  connectedCallback() {
    // The real element is focusable: MathLive attaches a shadow root with
    // `delegatesFocus` and an internal contenteditable sink. jsdom honours
    // neither, so without a tabindex `focus()` is a no-op and every autofocus
    // assertion in the suite would be vacuous.
    if (!this.hasAttribute('tabindex')) this.setAttribute('tabindex', '0')
  }

  get value(): string {
    return this.#value
  }

  /**
   * The real element fires `input` on every keystroke. Assigning `.value` here
   * does too, so a test that sets a value exercises the same React onChange path
   * a real keystroke would.
   */
  set value(next: string) {
    this.#value = next
    this.dispatchEvent(new Event('input', { bubbles: true }))
  }

  get readOnly(): boolean {
    return this.hasAttribute('readonly')
  }

  set readOnly(v: boolean) {
    if (v) this.setAttribute('readonly', '')
    else this.removeAttribute('readonly')
  }

  mathVirtualKeyboardPolicy = 'auto'

  /**
   * Emulates just enough of MathLive's plain-text serializer for the wrapper's
   * normalization to be exercised for real: a fraction key has to come back as
   * `(1)/(2)`, not as LaTeX, or `toPlainAnswer` would be a no-op under test and the
   * component suite would pass while the production path did nothing.
   *
   * Deliberately crude. The authoritative recording of MathLive's serializer is
   * test/fixtures/mathlive-plain-text.json, and the exact-output assertions live in
   * test/plainAnswer.test.ts against that fixture.
   */
  getValue(format?: string): string {
    if (format !== 'plain-text') return this.#value
    return this.#value
      .replace(/\\frac\{([^{}]*)\}\{([^{}]*)\}/g, '($1)/($2)')
      .replace(/\\sqrt\{([^{}]*)\}/g, 'sqrt($1)')
      .replace(/\^\{([^{}]*)\}/g, '^($1)')
      .replace(/\\times/g, ' * ')
      .replace(/\\div/g, ' -: ')
      .replace(/\\pi/g, 'pi')
      .replace(/\\left|\\right/g, '')
  }

  focus() {
    // super.focus() so jsdom updates document.activeElement; the real element
    // is focusable and the hosts' autofocus assertions depend on that.
    super.focus()
    this.dispatchEvent(new Event('focus'))
  }

  addEventListener(type: string, handler: (e: Event) => void) {
    let set = this.#listeners.get(type)
    if (!set) {
      set = new Set()
      this.#listeners.set(type, set)
    }
    set.add(handler)
    super.addEventListener(type, handler as EventListener)
  }

  removeEventListener(type: string, handler: (e: Event) => void) {
    this.#listeners.get(type)?.delete(handler)
    super.removeEventListener(type, handler as EventListener)
  }
}

if (typeof customElements !== 'undefined' && !customElements.get('math-field')) {
  customElements.define('math-field', StubMathField)
}

// The wrapper talks to `window.mathVirtualKeyboard`, which MathLive installs on
// import. Assigned through a cast rather than a global declaration; see the note
// at the top of this file.
if (typeof window !== 'undefined' && !(window as { mathVirtualKeyboard?: unknown }).mathVirtualKeyboard) {
  ;(window as unknown as { mathVirtualKeyboard: unknown }).mathVirtualKeyboard = {
    visible: false,
    layouts: [] as unknown,
    show(this: { visible: boolean }) {
      this.visible = true
    },
    hide(this: { visible: boolean }) {
      this.visible = false
    },
    addEventListener() {},
    // Present on the real object, and required: the geometry listener's cleanup calls it,
    // and a stub without it fails on unmount rather than on the assertion under test.
    removeEventListener() {},
    boundingRect: { height: 0, width: 0, top: 0, left: 0, right: 0, bottom: 0, x: 0, y: 0 },
  }
}

/**
 * The real module exports the class as a value, not only as a type, because
 * `soundsDirectory` and `keypressVibration` are *statics* and the instance accessors
 * deliberately throw. `MathLiveField` imports it as a value to set them, so the stub has
 * to carry it or every component test dies with
 * "Cannot set properties of undefined (setting 'soundsDirectory')".
 */
export const MathfieldElement = Object.assign(StubMathField, {
  soundsDirectory: './sounds' as string | null,
  keypressVibration: true,
})

export const version = { mathlive: '0.111.0-stub' }
export default {}