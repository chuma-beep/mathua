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
`components/math/` imports it, and removing it costs no other behaviour.

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

## 7. The LaTeX corpus check

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

Current state of the corpus:

| | spans |
|---|---|
| checked (parsed cleanly) | 37,465 |
| **malformed** | **517** (329 unique) |
| unchecked: not an expression (environments, prose-in-math, arrows) | 6,496 |
| unchecked: valid but outside CE's numeric type system | 773 |
| unchecked: command or operand CE does not implement | 1,621 |
| engine crashes (CE overflowed its own recursion) | 0 |
| **total** | **46,872** across 577 files |

Of the 517, the largest single cause is a Unicode apostrophe where LaTeX wants ASCII:
`f’(x) = 18x - 4` is 162 of them, from the scrape. KaTeX renders it anyway, which is
exactly why no existing gate caught it. The rest are unbalanced delimiters, stray `&`
outside an `aligned` environment, and one prose paragraph that `latexnorm.Scan`
classified as math.
