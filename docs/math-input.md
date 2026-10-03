# Math input

How a learner's answer travels from their fingers to a verdict, and what has to be
true at each step.

The short version: **the answer is a plain string on the wire, and it always has
been.** The math editor changed where that string comes from and how good it is. It
did not change the wire format, the endpoints, or the graders. This document records
the contract, the parts of MathJSON that are reserved for later, and the places where
the third-party libraries did not do what their names suggest.

- The editor: `web/next-app/components/math/`
- The compatibility shim: `web/next-app/components/math/plainAnswer.ts`
- The proof it survives the graders: `internal/grader/mathlive_compat_test.go`

---

## 1. The wire format today

```
learner types  →  MathLive (LaTeX internally)
               →  MathAnswerInput holds the LaTeX, publishes plain text
               →  POST /api/study/answer   { concept_id, answer, elapsed, question }
               →  engine.GradeAnswer(conceptID, expected, answer)
               →  grader.Router / GradedGenerator / SymPy
               →  correct, streak, mastery, XP
```

`answer` is a string. It was a string before the editor existed and it is a string
now. No endpoint field was added, no grading type was added, and no mastery rule was
touched. That was deliberate: ADR-005 makes `Engine.GradeAnswer` the single grading
entry point, and keeping it untouched is what makes the editor safe to ship.

A host sees exactly this:

```tsx
<MathAnswerInput value={answer} onChange={setAnswer} gradingType={…} conceptId={…} />
```

It stores a string and submits a string. `MathInputValue` is internal to the wrapper:

```ts
type MathInputValue = { latex: string; plainAnswer: string; mathJson?: unknown }
```

`mathJson` is reserved and always `undefined` today — see §3.

### Which concepts get the editor

`components/math/keyboards.ts` is the only place this is decided.

| `grading_type` | control | concepts |
|---|---|---|
| `numeric`, `symbolic`, `expression`, `polynomial`, `complex` | MathLive editor | 415 |
| `multiple_choice`, `tuple`, `ordering`, `comparison` | plain `<input>` + `SymbolPalette` | 242 |

The split is not cosmetic. A third of the corpus answers with a tap target or a short
list; a math editor there is a downgrade, not an upgrade. The server already sends
`grading_type` on every question (ADR-009), so this needed no new schema and no
per-question special cases.

---

## 2. Why a shim exists

MathLive's own plain-text output is not what Mathua's graders read. Recorded from
0.111.0 in Chromium (`web/next-app/test/fixtures/mathlive-plain-text.json`,
regenerate with `node scripts/record-mathlive-corpus.mjs`):

| Learner types | MathLive emits | Grader verdict |
|---|---|---|
| `\frac{1}{2}` | `(1)/(2)` | **rejected** — `numericGrader` parses only `1/2` |
| `\frac{2}{4}` | `(2)/(4)` | rejected |
| `4\frac{1}{10}` | `4(1)/(10)` | rejected — reads as `41/10`, not `4 1/10` |
| `7 \div 2` | `7 -: 2` | **false miss** — SymPy reads `-:` as subtraction |
| `\sin(x)+\cos(x)` | `sin (x)+cos (x)` | unreliable — SymPy's implicit-multiplication pass needs `sin(` |

The fifth row is the one that already bit: the retired `SymbolPalette` inserted the
`×` glyph, which `grading/sympy_service.py` rejects outright. Not a new class of bug —
an existing one.

`toPlainAnswer()` removes four artifacts, in this order, and nothing else:

1. `-:` → `/` (division)
2. `4(1)/(10)` → `4 1/10` (mixed number — **must** precede rule 3)
3. `(1)/(2)` → `1/2`
4. drop a space between a name and its `(…` or `^…`

It is a compatibility shim over a third-party serializer, not a parser. It never
reinterprets an expression. Each rule has a test in `test/plainAnswer.test.ts`, and
`internal/grader/mathlive_compat_test.go` re-implements it independently and asserts
the two agree by running the real graders over the recording. That test is the
evidence that the bridge is sound; it is not decoration.

---

## 3. MathJSON: reserved, not used

`MathInputValue.mathJson` is `undefined` and no call site reads it. It exists so the
wire format can change later without touching the five answer hosts.

**MathLive 0.111.0 cannot produce MathJSON.** Verified against the shipped types and
bundle:

- `OutputFormat` is `"ascii-math" | "latex" | "latex-expanded" | "latex-unstyled" | "latex-without-placeholders" | "typst" | "math-ml" | "plain-text" | "spoken" | …`. There is no `mathjson`.
- `getValue('mathjson')` no longer exists.
- `convertLatexToMathJson()` is no longer exported. The only surviving MathJSON code in the 1.4 MB bundle is a private `mathJsonToLatex` used for the clipboard.

MathJSON is now Compute Engine's exclusive domain (`@cortex-js/compute-engine`, also
MIT). So if that layer is ever added, the path is **MathLive → LaTeX → `ce.parse()` →
`.json`**, not MathLive → MathJSON.

MathLive still *consumes* MathJSON in two places — clipboard paste with
`format: "math-json"`, and as its default `serializeToLatex` — and it does **not**
bundle the library that does it. It looks the engine up on a global symbol:

```js
globalThis[Symbol.for("io.cortexjs.compute-engine")].ComputeEngine
```

and, finding nothing, logs *"The CortexJS Compute Engine library is not available"*
and returns an empty string. So if `@cortex-js/compute-engine` is ever added, simply
importing it is enough — MathLive picks it up through that global with no wiring. That
is the cheapest possible first step, and it fixes MathLive's MathJSON paste path
rather than changing any of ours.

### The subset a Go grader would actually receive

**Read from the Compute Engine 0.147 documentation, not executed locally** — Compute
Engine is not a dependency yet, so this table is what the library documents, not
something this repository has observed. The first version of it that *is* observed
should come from a corpus recorded the way
`test/fixtures/mathlive-plain-text.json` was.

The brief this work came from asked for "numbers, symbols, Add, Subtract, Multiply,
Divide, Power, Negate, Sqrt, Rational, plus function applications". Two of those do
not survive canonicalisation, and three it omits do. Against
`ce.parse()`'s canonical form:

| Expression | MathJSON | Note |
|---|---|---|
| `\frac{1}{2}` | `["Rational", 1, 2]` | **not** `["Divide", …]` |
| `\frac{x}{y}` | `["Rational", "x", "y"]` | symbolic denominator, still `Rational` |
| `\sqrt{9}` | `["Square", 3]` | evaluates to the number |
| `\sqrt{x}` | `["Power", "x", "Half"]` | **`"Half"` is a symbol, not `1/2`** |
| `\sqrt[3]{x}` | `["Power", "x", ["Rational", 1, 3]]` | |
| `x^2` | `["Square", "x"]` | **not** `["Power", "x", 2]` |
| `a - b` | `["Add", ["Negate", "b"], "a"]` | **`Subtract` is eliminated** |
| `\div 2` | `["Divide", …]` | only in `{form: 'raw'}`; canonical form folds it |
| `3.14` | `3.14`, or `{"num": "3.14…"}` | high precision needs the object form |
| `x` | `"x"`, or `{"sym": "x"}` | metadata needs the object form |

So a Go-side subset would need `Add`, `Negate`, `Multiply`, `Power`, `Square`,
`Divide` (raw only), `Rational`, `Sqrt` (raw only), `Apply`, plus the `"Half"`
symbol — and it would have to accept that **`Subtract` and `Sqrt` essentially never
arrive**.

Recommended stance if this is ever built: parse the raw form (`{form: 'raw'}`) rather
than canonical, so the tree mirrors what the learner typed, and keep `Sqrt`/`Divide`
in the subset anyway for safety.

### Is it worth building?

Not yet, on the evidence. `getValue('plain-text')` — MathLive's own public API, no
Compute Engine — already round-trips fractions, exponents, roots, mixed numbers,
implicit multiplication, subtraction, powers, ratios, remainder formats and symbolic
variables through the existing graders, proven in `mathlive_compat_test.go`. The
corpus test currently covers 12 of 34 recorded cases through the pure-Go path; the
remaining 22 are the SymPy-backed types, which SymPy already handles with `trigsimp`,
`radsimp`, `expand`, `factor` and random-substitution numeric checks.

A MathJSON layer would earn its place if it bought *better verdicts* — evaluating a
constant expression, or reasoning about equivalent forms without a subprocess. It
should not be added merely to have a structured representation.

---

## 4. The fallback grading path

Two different fallbacks, often confused:

**The grading fallback.** `grader.Result.Unavailable` marks "the grader could not
run" — a missing or broken SymPy runtime, a timeout, a busy semaphore. ADR-006 is
explicit: this is never a student miss. No attempt, no streak reset, no weakness bump,
no XP, and the client is told to retry rather than penalise.

**The editor fallback.** If MathLive's chunk fails to load — blocked asset, old
browser, flaky connection — `MathInput` falls back to the plain `<input>`. An error
boundary catches it, because a chunk that fails to *load* never mounts and no ref
check would ever fire. This is not a lesser editor: the plain input grades correctly
for every concept in the corpus, which is exactly what a learner with a blocked asset
needs. It is the editor that always works.

Both are tested (`test/mathInputFallback.test.tsx`, and the existing grader tests).

### Known limitation, not yet fixed

`mathVirtualKeyboard.hide()` does not work in MathLive 0.111.0 under Chromium.
Called directly from the page it left `visible === true` and the panel in the DOM,
under both the `auto` and `manual` policies, and neither an outside tap nor a blur
freed it. `MathLiveField` calls `hide()` on a pointerdown outside the field, but the
mobile suite does **not** assert dismissal, because the behaviour is unproven. Until
it is, the working exits on a phone are submitting (which blurs the field) or the
keypad's own hide keycap.

---

## 5. API surprises

Everything here cost time to find and is worth writing down.

**`getValue('plain-text')` returns `""` in jsdom**, always — so the component suite
stubs the library out (`test/stubs/mathlive.ts`, aliased in `vitest.config.mts`) and
the real behaviour is verified in Playwright. MathLive's `node` export condition also
resolves to the SSR build under vitest, which never registers `<math-field>` at all.

**`next/dynamic({ ssr: false })` hangs under vitest.** It resolves through Next's
build-time chunk loader, which does not exist outside the Next build, so the loading
state never resolves. `React.lazy` + `Suspense` gives the same production result and
is testable.

**React 18 does not map `className` on custom elements.** It forwards props verbatim,
so `<math-field className="…">` writes the attribute `classname`, which no CSS selector
matches. It must be `class`.

**A bare-string Vite alias is a prefix match.** Aliasing `mathlive` also rewrote
`mathlive/static.css`, and the lazy import rejected. Use an exact RegExp — and put the
alias under `test.alias`, because vitest externalises `node_modules` deps *before*
`resolve.alias` applies.

**`mathlive/static.css` is the import path, not the filename.** The `exports` map
exposes `"./static.css"`, not `"./mathlive-static.css"`.

**MathLive loads its own webfonts from `<origin>/fonts` at runtime.** There is no
server to serve them under `output: 'export'`, so `scripts/copy-mathlive-fonts.mjs`
copies them from `node_modules` on `predev`/`prebuild`. Without them every fraction
renders in a fallback face.

**Keycap `width` accepts only `0.5 | 1 | 1.5 | 2 | 5`.** Anything else is a type error
at best and silently ignored at worst.

**Keycap width is derived from `10cqw`** — 10% of the keypad panel — which lands
around 37px on a 412px phone. `--keycap-width` overrides it. Height alone was never
the binding constraint for the 44px touch-target floor; this was.

**Typing the `√` glyph produces `\surd`, not `\sqrt{9}`.** `\surd` is an ordinary
symbol with the radicand outside it. This is why the keypad's `√` key inserts
`\sqrt{#0}` — and why the retired `SymbolPalette`'s `√` could never have worked.

**`preventDefault()` on Enter is unnecessary.** MathLive treats Enter as "commit the
expression". Leave it alone.

**A backspace immediately after a typed fraction does nothing** in 0.111.0 — the caret
sits after the fraction's trailing placeholder rather than adjacent to a removable atom.
A second backspace removes it. Asserted as observed, not as expected.

**Playwright's atomic `locator.click()` does not focus MathLive**; `mouse.down()` and
`mouse.up()` separated by a delay does. And typing needs ~40 ms per key, or MathLive's
asynchronous parser drops characters — a harness artefact, not a product bug.

---

## 6. Where the keyboard choice lives

`keyboards.ts` maps `grading_type` + `domain` → `MathKeyboardMode`; `layouts.ts` holds
the five keycap grids as plain data. No React component branches on a domain.

```
grading_type ──► editor   (math | text)          keyboards.ts
concept_id ──► domain ────► MathKeyboardMode ──► layout   layouts.ts
```

Unmapped domains fall back to `algebra`, the most general of the five: falling back to
`arithmetic` would hide the letters a linear-algebra or number-theory question needs,
which is the worse failure.

Structural keys insert structures. The fraction key inserts `\frac{#@}{#0}` — with a
navigable placeholder — not `/`. A keycap that inserted the character would look
identical in a screenshot and behave nothing like a fraction. `test/keyboards.test.ts`
asserts this for every mode, and `e2e/math-input.spec.ts` asserts it against the real
library.