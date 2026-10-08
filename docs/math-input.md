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
MIT), and it **is** a dependency now — for the self-check, not for MathLive. So if a
MathJSON wire format is ever added, the path is **MathLive → LaTeX → `ce.parse()` →
`.json`**, not MathLive → MathJSON.

Importing the library is also all it takes to switch on MathLive's own MathJSON support:
it resolves the engine through `globalThis[Symbol.for("io.cortexjs.compute-engine")]`,
which CE sets as a side effect of being imported, and logs an error when it is absent.
`test/equivalence.test.ts` asserts that global exists.

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

**Observed, not read from documentation.** `@cortex-js/compute-engine@0.147.0` is
installed, and this table is `JSON.stringify(ce.parse(latex).json)` run against it. That
matters, because the library's *documented* canonical form does not match what it
emits: the docs say `Sqrt` becomes `["Power", x, "Half"]` and `["Power", x, 2]`
becomes `["Square", x]`. In 0.147.0 neither happens — `Sqrt` and `Root` survive intact
and `x^2` stays a `Power`. A Go grader written against the documentation would fail on
almost every input.

| Expression | MathJSON from 0.147.0 | Note |
|---|---|---|
| `\frac{1}{2}` | `["Rational", 1, 2]` | literal fractions become `Rational`, not `Divide` |
| `\frac{x}{y}` | `["Divide", "x", "y"]` | a symbolic denominator stays a `Divide` |
| `\frac{1}{2}+\frac{1}{3}` | `["Rational", 5, 6]` | parsed *and* evaluated on the way in |
| `x^2` | `["Power", "x", 2]` | **not** `["Square", "x"]` |
| `x^2+2x+1` | `["Add", ["Power","x",2], ["Multiply",2,"x"], 1]` | |
| `\sqrt{9}` | `3` | a root of a number is just the number |
| `\sqrt{x}` | `["Sqrt", "x"]` | **`Sqrt` survives** canonicalisation |
| `\sqrt[3]{x}` | `["Root", "x", 3]` | |
| `x-2` | `["Add", "x", -2]` | **`Subtract` is eliminated**, as documented |
| `-x` | `["Negate", "x"]` | |
| `2x+3` | `["Add", ["Multiply", 2, "x"], 3]` | implicit multiplication is explicit here |
| `\sin(x)` | `["Sin", "x"]` | functions appear bare, not wrapped in `Apply` |
| `\pi` | `"Pi"` | a named constant, not a symbol |
| `3.14` | `3.14` | high precision needs the `{"num": "…"}` object form |

So the subset a Go grader would need is `Add`, `Multiply`, `Power`, `Divide`, `Negate`,
`Rational`, `Sqrt`, `Root`, `Sin`/`Cos`/… , `Pi`, and bare symbols — and it must
tolerate that **`Subtract` never arrives**, `Divide` only survives with a symbolic
denominator, and constant subexpressions arrive already evaluated.

### The feedback contract

`equivalence()` is three-valued, and the naming is deliberate — a string union rather
than `true | false | undefined`, because `'unknown'` is an *absence* and `undefined`
would be indistinguishable from a function that failed to return.

| state | meaning |
|---|---|
| `'equivalent'` | the two expressions are the same |
| `'different'` | they are provably not |
| `'unknown'` | Compute Engine could not decide, and we decline to guess |

The actual results, so the states are not aspirational:

| comparison | result |
|---|---|
| `x²+2x+1` vs `(x+1)²` | `equivalent` |
| `2^10` vs `1024` | `different` |
| `\sqrt{x²}` vs `x` | `unknown` — correctly refuses; it is `|x|` |
| `\sqrt{x²}` vs `\lvert x\rvert` | `equivalent` |
| `2x+3` vs `x²+2x+1` | **`unknown`**, not `different` |

That last row is the one that matters. CE will not call two plainly unequal symbolic
expressions *different* — it declines rather than guess — which is why `'unknown'` is a
first-class outcome and not an error path.

**No state renders a verdict.** Not `'equivalent'`, not `'different'`, and above all not
`'unknown'`: nothing may render as "correct", "incorrect", "wrong", "try again" or "not
yet". `test/selfCheck.test.tsx` asserts that for all three states, using a real pair
that reaches each one.

This is not a missing feature, and the reason is the answer oracle. The client does not
hold the expected answer — grading is server-side and answers are anchored there
(ADR-005) — so there is nothing to compare against without shipping the answer to the
browser. `simplificationHint` is what the UI uses instead: it reports what the
learner's own expression simplifies to and cannot produce a verdict even in principle,
because it never sees a second expression.

### What Compute Engine does when asked to judge equality

Measured, because the answer decides whether an equivalence feature can exist at all
(`test/equivalence.test.ts`):

| Comparison | `isEqual` | `isIdenticallyEqual` | `simplify` fallback |
|---|---|---|---|
| `2x+3` vs `2*x+3` | `true` | `true` | — |
| `(x+1)^2` vs `x^2+2x+1` | **`undefined`** | `true` | — |
| `sin^2(x)+cos^2(x)` vs `1` | **`undefined`** | **`undefined`** | `1` — resolves it |
| `\sqrt{x^2}` vs `x` | `undefined` | **`undefined`** (correctly refuses) | stays unequal |
| `\sqrt{x^2}` vs `\|x\|` | `undefined` | `true` | — |
| `2x+3` vs `2x+4` | `undefined` | **`undefined`** | stays unequal |
| `2^10` vs `1024` | `true` | `false` | — |

Three things follow.

**`isEqual` is unusable for grading.** It returns `undefined` for the symbolic
identities that matter — expanded versus factored, and every trig identity. It answers
"do these have the same *value*", and two expressions with a free variable do not.

**`isIdenticallyEqual` is the only one that can refuse**, which is exactly what is
wanted for `sqrt(x^2)` vs `x`: it declines rather than confirming a false identity, and
it knows the right answer when asked about `|x|`. It also declines to call `2x+3` and
`2x+4` *different*, so a hint built on it would have to say "I don't know" for plainly
wrong input.

**`simplify` is the tie-breaker**, and it is what makes trig identities work — and it is
also what produces `|x|` from `sqrt(x^2)`, which is the single best argument for using a
CAS here: the naive reading of `sqrt(x^2)` is `x`, and that is wrong for every negative
x.

Which is why the shipped feature is a **self-check** — "that simplifies to X" about the
learner's own expression — and not a comparison against an answer. The client does not
have the expected answer, by design; see §4.

### Is it worth building?

Not yet, on the evidence. `getValue('plain-text')` — MathLive's own public API, no
Compute Engine — already round-trips fractions, exponents, roots, mixed numbers,
implicit multiplication, subtraction, powers, ratios, remainder formats and symbolic
variables through the existing graders, proven in `mathlive_compat_test.go`. The
corpus test currently covers 12 of 34 recorded cases through the pure-Go path; the
remaining 22 are the SymPy-backed types, which SymPy already handles with `trigsimp`,
`radsimp`, `expand`, `factor` and random-substitution numeric checks.

A MathJSON layer would earn its place if it bought *better verdicts*. The one gap
`plain-text` leaves is a constant expression the numeric grader cannot evaluate
(`10^6`, `12 × 4`); a MathJSON-aware numeric path would close it. It should not be added
merely to have a structured representation.

---

### The self-check, and what it costs

`components/math/SelfCheck.tsx` shows the learner what their own expression simplifies
to: `1/2 + 1/3` → "that simplifies to 5/6", `sqrt(x^2)` → "that simplifies to |x|".

It is deliberately not an equivalence check against the expected answer. The client does
not hold the expected answer, and shipping it to power a hint would turn the hint into
an answer oracle and undo ADR-005. No code path lets the hint affect a submitted value,
and `test/selfCheck.test.tsx` asserts that: the hint appears in the DOM while the value
the host receives is untouched.

The cost is the part worth arguing about. Compute Engine is **788 kB gzipped** — nearly
four times MathLive. It is never in any route's initial JS, and never fetched unless a
learner pauses mid-answer with arithmetic on screen. Measured:

- Routes without math input: **+0 kB**.
- Routes with the editor: **+1 kB gzip** (the hint's shell).
- One **788 kB gzip** download, on first use, expression-valued concepts only.

It is skipped entirely when `navigator.connection.saveData` is set, or the effective
type is `2g`/`slow-2g`. A one-line simplification is not worth someone's data allowance,
and a feature that quietly spends it is worse than one that is visibly absent.

If that trade is wrong, `SelfCheck.tsx` is the only file to delete: nothing else in
`components/math/` imports it, and removing it costs no other behavior.

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

### The keycap `latex`/`label` trap

Two of Mathua's keypad keys were broken in **all five layouts** and no test noticed:
backspace typed the literal text `[backspace]` into the answer field, and the zero key
typed `[0]`. Both were declared like this:

```ts
const BACKSPACE: Keycap = { latex: '[backspace]', width: 2 }
const ZERO_WIDE: Keycap = { latex: '[0]', width: 2 }
```

    Assumed: a keycap whose `latex` names a bracketed command performs that command.
    Actual 0.111.0: `latex` is *inserted*. `executeKeycapCommand` resolves in order —
    `command`, then `insert`, then `key`, then `latex` (run as an insert), and only
    then `typedText(label)`. Separately, `normalizeKeycap` merges MathLive's own
    `KEYCAP_SHORTCUTS` definition only when the keycap is a bare string, or carries a
    `label` or `key` naming a known shortcut. `latex` is not one of those.
    Difference: these keycaps matched no shortcut *and* fell through to the insert
    branch, so they typed their own source text. `KEYCAP_SHORTCUTS["[backspace]"]`
    carries `command: "performWithFeedback(deleteBackward)"` and
    `KEYCAP_SHORTCUTS["[0]"]` carries `latex: "0"`; neither was consulted.
    Fix: name the shortcut in `label`, which is what MathLive's own layouts do
    (`{ label: "[backspace]", width: 1 }`). `test/keyboards.test.ts` now fails on any
    bracketed token in a keycap's `latex`, which is the whole bug class.
    Source: `node_modules/mathlive/mathlive.mjs` — `executeKeycapCommand`,
    `normalizeKeycap`, `renderKeycap`, `KEYCAP_SHORTCUTS`.

A second trap sat behind it, and it is why the bug survived: **`renderKeycap` only
adds the `MLK__keycap` class when the keycap's class does not already contain
`separator`, `action`, `shift`, `fnbutton` or `bigfnbutton`.** So backspace, the caret
arrows, the dismiss key and every separator render *without* it. Every selector in the
suite used `.MLK__keycap`, including the 44px touch-target audit — so that audit
measured only the keys that were already working and reported a clean pass over a
keypad whose delete key typed its own name. `e2e/helpers/answer.ts` now exposes
`keycaps()` (the row's visible direct children) and `glyphKeycap()`, and the audit
asserts that it reached the delete and caret keys, so "no undersized keycaps" cannot
quietly go back to meaning "no keycaps measured".

### The loading window, measured

Between the answer UI appearing and MathLive being usable there is a window. It used to
be filled with `<div className="h-12" aria-hidden />`: the right size, in the right
place, and completely inert. A tap on it looked like it landed on an answer field and
then did nothing.

Measured (`make diagnose-answer` equivalent: `e2e/answer-diagnosis.spec.ts`, everything
on one `performance.now()` clock):

| MathLive chunk latency | window with no editor |
|---|---|
| 0ms (cached) | 322ms |
| 500ms | 611ms |
| 1500ms | 1611ms |
| 3000ms | 3116ms |

And once the editor is up, nothing else is slow: input reaches React state in 2-7ms and
the first tap on Check reaches the API and returns in 12-14ms, on desktop and both
mobile profiles.

The loading state is now `PlainAnswerInput` — the same working plain input used when
the chunk fails outright. So the learner can answer during the window rather than
repeating themselves, and whatever they type is adopted when MathLive arrives.

Two things to be careful about if you re-measure this:

- **CDP network throttling does not reach the MathLive chunk here.** It reported
  `transferSize: 0` and completed in 39ms under a 4 Mbps profile while the main chunks
  correctly took 1100-1400ms — it is served from Chromium's memory cache. A first pass
  took that as "the network does not matter". `Network.setCacheDisabled` did not change
  it either.
- **Delay the chunk instead, and prove the delay applied.** `delayMathLiveChunk()`
  finds its target by scanning the build output for MathLive's own markers, routes at
  context level, and counts its interceptions; every scenario reports `delayApplied`.

The 310ms floor is a floor for this harness, not a phone: desktop CPU, cached chunk.
Network contribution is unmeasured and the cost of parsing a 794 kB decoded bundle on a
mid-range phone is unknown. Both are larger.

### Dismissal, and why Mathua no longer tries

Mathua used to hide the virtual keypad from a capture-phase `pointerdown` listener on
`document`: any tap outside the field, and the keypad went away. It had to go.

MathLive raises the keyboard from exactly one place — a document `focusin` listener
that calls `show()` after 300ms, gated on `policy === 'auto' && hasEditableContent`.
There is no re-show path. So hiding the panel while the field *still holds focus*
strands it: no focus transition means no `focusin`, so tapping the field again does
nothing, and the learner has to keep tapping. On iOS Safari a tap on non-focusable
content frequently does not blur the input, so this reproduced reliably on a phone
while being invisible in Chromium.

MathLive already hides the keypad itself, from a `focusout` listener, when focus
genuinely leaves the field. That is the behavior worth having, so the outside-tap
listener was redundant *and* harmful. Mathua now hides the keypad in exactly one
place: when a focused field unmounts, since the Learn feed swaps cards under you.

`MathLiveField`'s `pointerdown` handler is the recovery path — it calls `show()`, and
because `show()` only dispatches a message, calling it when the panel is already
visible is harmless.

### Not yet fixed

`mathVirtualKeyboard.hide()` had no measurable effect under Chromium when called
directly from the page: `visible` stayed `true` and the panel stayed in the DOM, under
both `auto` and `manual`. That is why the strand could not be reproduced in Playwright
and why `e2e/mobile-math.spec.ts` asserts the *call* Mathua makes rather than the
resulting visibility — spying on `mathVirtualKeyboard.hide` fails on the old listener
and passes now. The recovery behavior itself still wants a real-device check.

### Every question gets the editor

The editor was once scoped to the five expression-valued `grading_type`s. That left 223
of 657 concepts — 209 `multiple_choice`, 11 `tuple`, 1 `ordering`, 2 `comparison` —
typing mathematics into a plain `<input>`, because the answer hosts render a text field
even when the question also offers tappable options.

Making it universal was **measured, not assumed**. The multiple-choice grader is
`strings.EqualFold(expected, answer)`, and MathLive reads bare adjacent letters as
implicit multiplication:

| expected | MathLive emitted | grades? |
|---|---|---|
| `yes` | `"y e s"` | ✗ |
| `no` | `"n o"` | ✗ |
| `odd` | `"o d d"` | ✗ |
| `addition` | `"a d d i t i o n"` | ✗ |
| `pi/2` | `"p i/2"` | ✗ |
| `1/2` | `"(1)/(2)"` | ✓ |

`yes` is the single most common expected answer in the corpus (673) and `no` is second
(171), so universal MathLive would have broken roughly 850 questions. Two rules close
it, both covered by `internal/grader/mathlive_compat_test.go` against a *recording* of
what the library emits:

- `fixImplicitLetterSpacing` rejoins a word the keyboard split into letters. The guard is
  the hard part: `2pi r` also contains a letter-space-letter, and collapsing it gives
  `2pir`. So the run is extended left over its whole adjacent letter sequence and only
  collapsed when what precedes *that* is a letter or nothing.
- `fixFractionParens` accepts a symbolic numerator, `(pi)/(2)` → `pi/2`, while staying
  narrow enough to leave `(x + 1)/(y)` alone.

Multiple choice still renders `ChoiceOptions` above the field, so yes/no remains a tap.

### The keypad answers for the corpus

Two things were missing and both blocked real questions:

- **No comma in any layout.** Every `tuple` and `ordering` answer is a list, so those 12
  concepts had no way to enter their answer at all.
- **No `< > ≤ ≥`.** `dec.basics.compare` and `frac.ops.compare` grade a comparison.

And a structural one, which was the expensive part: **MathLive's toolbar navigates
*layouts*, not *layers*.** It emits `data-layer="${layout.layers[0].id}"` for each entry
in `kb.layouts`, so a layout with three `layers` renders only the first and the other two
are unreachable. A keycap carrying a `layer` property only gains a `layer-switch` *class*
and does nothing when pressed. Layers 2+ had been unreachable from the start — the old
operators and scientific-notation layer had no way in, which is the same class of bug as
the keycap that typed its own name (ADR-029).

So each page is its own single-layer `Layout`, `layoutsForMode` returns them together, and
`kb.layouts = [...pages, 'alphabetic']`. Every layer also needs a `label`, because the
toolbar renders it and shows "untitled" without one. Letters come from MathLive's own
`alphabetic` layout rather than a 26-key grid that could not hold the 44px floor at
320px.

## 5. Mobile setup

Real-device verification lives in [`mobile-device-checklist.md`](./mobile-device-checklist.md).
Nothing in this document is a substitute for it: the automated suite uses emulated
viewports, which cannot reproduce OS-keyboard suppression, iOS viewport units, safe-area
insets, iOS focus zoom, or real parse cost. The numbers below are floors.

What is actually configured, all verified against the pinned MathLive 0.111.0 rather
than remembered from an earlier version.

**Virtual keyboard policy.** Not set. 0.111.0's default is `'auto'`, confirmed in the
bundle's own defaults object, and `auto` is the behavior Mathua wants: raise the keypad
on touch, leave it alone where there is a physical keyboard. The attribute exists as both
`mathVirtualKeyboardPolicy` and `math-virtual-keyboard-policy`. It is left at the default
rather than set redundantly, because setting it would assert something already true and
would need re-asserting if the default ever changed.

**Keyboard layouts.** Five, chosen from the concept's domain by `keyboards.ts`:
`arithmetic`, `fractions`, `algebra`, `geometry`, `calculus`. Installed on focus
(`kb.layouts = [layoutForMode(mode)]`), because the keyboard is a page-wide singleton and
every field would otherwise overwrite the last.

**Keyboard container.** MathLive's own. It positions the panel against the field, so it
does not cover the editor. What it cannot do is shrink the *page*, so
`geometrychange` is republished as one CSS custom property:

| what | where |
|---|---|
| keypad height | `--mathua-keyboard-height` on `:root` |
| keypad state | `data-mathua-keyboard="open"｜"closed"` on `:root` |

Stylesheets adapt to that instead of a hardcoded pixel height, which would be wrong on
every device. Safe-area is handled where it already was — `ChromeToggle`,
`ConceptGraphFlow`, `/history` — and deliberately *not* added to the answer area, because
MathLive's panel already sits above the home indicator and double-padding it would push
the field away from the keypad.

**Fonts.** `MathfieldElement._fontsDirectory` defaults to `'./fonts/'`, and Mathua copies
MathLive's webfonts there at build time (`scripts/copy-mathlive-fonts.mjs`). 300 kB of
woff2. The property is left at its default because the default is already correct.

**Sounds.** Disabled, and this was a real bug:

```
MathfieldElement.soundsDirectory = null   // static; the instance accessor throws
MathfieldElement.keypressVibration = false
```

0.111.0 defaults `_soundsDirectory` to `'./sounds'` and `keypressSound` to
`keypress-standard.wav`, loaded with a `fetch()` **on every keystroke**. Mathua ships
`public/fonts` and never shipped `public/sounds`, so each key requested a file that does
not exist. Measured: two taps produced two requests for
`/…/chunks/sounds/keypress-standard.wav`. `soundsDirectory = null` is the documented way
to load no sounds; `test/mobile-math.spec.ts` asserts none are requested.

**Compute Engine.** Dynamically imported by `SelfCheck`, never in initial JS. 787 kB
gzip, fetched once, on first use, and skipped entirely when
`navigator.connection.saveData` is set or the effective type is `2g`/`slow-3g`.

**Client-only loading.** `MathLiveField.tsx` is the only file that imports `mathlive`,
and it sits behind `React.lazy` inside a boundary. The boundary catches a chunk that
fails to *load*, which no ref check can: a chunk that never mounts has no ref to check.

### First Load JS, before and after

Measured from the built `out/`, gzip, summing the `<script src>` tags each prerendered
route actually ships.

| Route | Before | After | Δ |
|---|---|---|---|
| `/learn` | 436 kB | 437 kB | +1 |
| `/review` | 433 kB | 434 kB | +1 |
| `/goals` | 440 kB | 441 kB | +1 |
| `/onboard` | 438 kB | 438 kB | 0 |
| `/study` (later closed — a 460 B forwarder; see ADR-046) | 462 kB | 462 kB | 0 |
| `/profile` | 248 kB | 248 kB | 0 |
| `/history` | 411 kB | 411 kB | 0 |
| `/graph` | 252 kB | 252 kB | 0 |
| `/settings` | 247 kB | 246 kB | −1 |
| `/` | 236 kB | 236 kB | 0 |
| **total** | **3603 kB** | **3605 kB** | **+2** |

Lazily fetched, never in the above:

| chunk | gzip | what |
|---|---|---|
| `c48e1edc…` | 787 kB | Compute Engine — first simplification only |
| `428ea25f…` | 215 kB | MathLive — on any expression-valued question |

The MathLive chunk is 794 kB decoded, which is why the loading window is ~310ms even
with the bytes already cached: it is parse and evaluate, not download.

## 6. API surprises

Everything here cost time to find and is worth writing down.

**`math-mode-space` is ignored; set the property.** This is the spacebar. MathLive's own
`alphabetic` layout — reachable from the keypad toolbar on every question — ships

```js
{ label: " ", width: 1.5 }
```

A label and a width, and no `command`, `insert`, `latex` or `key`. `executeKeycapCommand`
resolves a keycap as `command` → `insert` → `key` → `latex` → `typedText(label)`, so
with none of the first four set it falls through to `typedText(" ")`. That reaches the
`[Space]` keystroke handler, which does:

```js
if (mathfield.options.mathModeSpace) { insert(mathModeSpace) }
else                               { moveAfterParent() }
```

`mathModeSpace` defaults to `""`. **So the spacebar moved the caret and inserted nothing** —
an upstream defect, not a Mathua regression. Three generator answer formats need a real
space: `Q R` (remainder), `a sqrt(b)` and `P=[[...]] D=[[...]]`.

The payload is `\,`, not a literal space, because TeX discards literal spaces in math mode.
`\,` is the proven round-trip value: the recorded corpus holds
`expression 5 R 3 <- 5\ R\ 3`, which comes back as plain text `5 R 3` and grades.

Set it as a **property**, not the documented attribute. Measured on 0.111.0 with React 18:
the attribute is present in the DOM and `placeholder` and `class` are both honoured from
it, but `mathModeSpace` stays `""` and the spacebar stays inert. `el.mathModeSpace = '\,'`
and `el.setOptions({ mathModeSpace: '\,' })` both work. An attribute that is silently
ignored is worse than none — it reads as configuration and is not.

Two more traps on the way, both the ADR-029 class again. A keycap of `{ label: ' ' }`
*does* work (it maps to `[Space]`) but renders as a blank face, which is half of why the
key read as broken. The obvious fix — a visible `label: '␣'` — **types the glyph `␣` into
the field**, because with no payload beside the label it falls through to
`typedText('␣')`. The payload has to be stated alongside: `{ label: '␣', latex: '\,' }`.

**The keypad's width budget is in keycap-widths, and it depends on the viewport.**
MathLive does not stretch keys to fill its panel: measured on a Pixel 7 (400px panel),
every keycap renders 45px and a `width: 2` keycap renders 93px, whatever the row holds.
So the binding constraint is the row's total *units*, and the ceiling is the point at
which `units × 45` stops fitting:

| panel | 5u | 6u | 7u | 8u | 9u | 10u |
|---|---|---|---|---|---|---|
| 400px (412px phone) | 45 | 45 | 45 | 45 | **42** | 37 |
| 308px (320px phone) | 45 | 45 | **41** | **36** | — | — |

Eight units is the design target and is what `test/keyboards.test.ts` asserts for all five
modes. **At 320px the ceiling is six units, and thirteen rows across four modes are already
over it** — so the 44px floor holds at 412px and does not hold at 320px. The e2e
touch-target audit passes at 320px only because it measures whichever layout the default
question selects, `ARITHMETIC`, the one mode whose rows all fit. Closing that means
splitting the eight-keycap rows across modes, which is a redesign rather than a fix.

**`geometrychange` carries no `detail`.** The 0.111.0 docs say
`evt.detail.boundingRect`. The library dispatches `new Event('geometrychange')` — a
plain `Event`, no detail at all — from a `ResizeObserver` and from `stateChanged()`. A
handler written to the documented shape reads `undefined`, computes a height of zero, and
reports the keypad as closed for the rest of the session: the exact failure the listener
exists to prevent. Read `kb.boundingRect` instead.

**`soundsDirectory` is static, and the instance accessor throws.** Setting
`field.soundsDirectory` throws `Error("Use MathfieldElement.soundsDirectory instead")`,
so it has to be the class. The default is `'./sounds'`, and sounds are fetched per
keystroke.

**The keycap `latex` is an insert, not a command.** See "The keycap `latex`/`label` trap"
above — backspace and the zero key typed their own source text in all five layouts.

**`getValue('plain-text')` returns `""` in jsdom**, always — so the component suite
stubs the library out (`test/stubs/mathlive.ts`, aliased in `vitest.config.mts`) and
the real behavior is verified in Playwright. MathLive's `node` export condition also
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

## 7. Where the keyboard choice lives

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

## 8. The LaTeX corpus check

`web/next-app/scripts/normalize-latex.mjs` runs every math span in the lesson corpus
through Compute Engine and reports the malformed ones (`make latex-normalize`).

The corpus gate (`test/latexCorpus.test.ts`) already proves KaTeX can *render* every
lesson, but that is a weak signal: KaTeX is forgiving by design and will render a
malformed fraction by guessing. This is the strict counterpart.

Four things about it are not obvious, and each was learned the hard way:

**The spans come from `cmd/latexdump -spans`, not from a regex.** `latexnorm.Scan` is
the only thing in the tree that knows a dollar can be a delimiter or a currency amount,
and it runs over the canonicalized body — the text the API actually serves. A regex
would report every "$5 and up" in the corpus as a syntax error.

**`ce.parse` does not throw.** It is total: a LaTeX syntax error comes back as a tree
containing an `["Error", reason, …]` node. A try/catch reports a corpus with thousands
of broken spans as perfectly clean, which is worse than no report because it is a
confident lie. Failure is read off the returned box.

**A parse failure is not automatically a defect.** CE is a numeric/symbolic CAS; this
corpus is largely abstract algebra, group theory and topology. CE reports `U \times F`
as `incompatible-type: number, function` — having parsed the syntax perfectly and then
declined to type a product of two sets — and does the same for `\mathbb{R}^2 \setminus
\{0\}`. Filing those as "unparseable" buries the real defects under ~1,600 spans of
noise and gets the report switched off. Reasons are bucketed by what CE objected to:
`syntax` (a real defect: unclosed delimiter, stray operator, stray token),
`not-numeric` (understood, refused to type it), `ce-gap` (a command CE has no rule for,
which cannot be told from valid-but-unsupported without a human).

**Every span is accounted for, and the buckets reconcile.** The summary asserts
`clean + malformed + unchecked == spans` and exits 2 if not, because that check is what
caught the first two versions of this script reporting "20,549 of 46,872" while
happily reporting success. It also states `checked_pct`, because a report that quietly
covers 44% of its input is not a clean bill of health. `make validate` does not include
it and CI runs it with `continue-on-error: true` — the corpus has pre-existing defects,
so a hard gate would block main on content rather than on a regression.

Current state of the corpus, after `cmd/latexfix` repaired what was mechanical:

| | spans |
|---|---|
| checked (parsed cleanly) | 37,657 |
| **malformed** | **228** (137 unique, in 60 files) |
| unchecked: not an expression (environments, prose-in-math, arrows) | 6,480 |
| unchecked: valid but outside CE's numeric type system | 773 |
| unchecked: command or operand CE does not implement | 1,650 |
| engine crashes (CE overflowed its own recursion) | 0 |
| **total** | **46,821** across 577 files |

`cmd/latexfix` took it from 517. Three classes were mechanical and are now fixed:

- **Unicode apostrophes in derivatives** (162 → 1). `f’(x)` is a typo. `‘a \text{is
  less than}\ b’` is a *quotation mark*, and rewriting it to `'` makes LaTeX render
  primes — a silent change to what the lesson teaches. The rule only fires on an
  attached apostrophe, and never on an opening curly quote.
- **Scrape artifacts**: `42✓` correctness ticks (68) and HTML entities leaked into
  math (`&gt;`, `&#39;`) (32).
- **Doubled delimiters** (14,936 `\(`, plus the display-math `\[`). The Algebrica
  scrape wrote `\(`, which `latexnorm` reads as an escaped backslash followed by a
  paren — so every Algebrica formula sat in a *prose* region and no region-based
  tool could see it. The canonicalizer is what turned it into `$ … $`, which is why
  the checker found defects the fixer could not reach. Un-escaping is provably
  free: `canonicalize` emits byte-identical output for `\(x\)` and `\(x\)`, and
  that invariant is a test.

Two findings worth keeping:

- **`\[` needs a discriminator.** A row break in an `align` environment is
  `\[2pt]` — two backslashes then a bracket holding a length, closed by a *single*
  `]`. So `\]` is unambiguously a display closer, while `\[` is a display opener
  only when not followed by a length. The corpus agrees: 4,496 `\]` against ~4,500
  non-row-break `\[`, and 1,157 genuine row breaks left alone.
- **Removing a `✓` can create a defect.** `$42✓$` → `$42$` is a *bare numeral set as
  math*, which the canonicalizer escapes to `\$42\$` — the learner sees literal
  dollar signs (ADR-016). So a repaired span that becomes a bare numeral loses its
  delimiters too: `$42✓$` → `42`. The old form was itself broken; the learner was
  seeing `\$1✓\$`.

Of the 517, the largest single cause is a Unicode apostrophe where LaTeX wants ASCII:
`f’(x) = 18x - 4` is 162 of them, from the scrape. KaTeX renders it anyway, which is
exactly why no existing gate caught it. The rest are unbalanced delimiters, stray `&`
outside an `aligned` environment, and one prose paragraph that `latexnorm.Scan`
classified as math.
