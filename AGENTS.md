# AGENTS.md, Mathua

Hard rules for anyone or anything editing this repository. Each one states the rule, why it
exists and how it is checked, because a rule nobody can verify is an aspiration.

---

## 1. No emoji anywhere we author

**Emoji are not permitted in any file we author**, application source, components, tests,
scripts, Go or documentation. This is not a style preference; see the reasons below.

### Why

- **An emoji is not a theme surface.** Its glyph, color and baseline come from the *font*, so it
  cannot inherit `currentColor`, cannot follow light and dark mode and cannot be sized
  consistently against the mono type Mathua's UI uses. Everything else in the app is a theme
  token; an emoji is the one thing that is not.
- **They render differently per platform.** The same codepoint is a different picture on macOS,
  iOS, Windows and Android, so a UI that looks right on a developer's machine is not the UI a
  learner sees. That is how U+23F5 on the graph's flow toggle ended up at an unpredictable
  baseline next to the word "Flow" inside a button already set to `fontSize: 10`.
- **They are not accessible by default.** They are announced by name ("hourglass"), which on a
  line that already says "3 concepts due for review" is noise, and they vary in how they
  degrade. Icons get `aria-hidden` when decorative.

### What to use instead

- **`components/icons`.** Every icon in the app comes from that one module, which wraps
  `@animateicons/react` (the animated lucide set). Import from there and never from an icon
  library directly, the wrapper is what pins the size to 16px and the color to `currentColor`,
  and the animated icons render a `<div>` whose default color is not `currentColor`, so an icon
  added outside the wrapper does not follow light and dark mode and does not keep the collapsed
  44px sidebar rail even.
- **Text.** A label beats a symbol most of the time. In tables, a word ("Yes", "Not stored")
  beats a mark.

Icons animate on mount and again on hover. Two are driven imperatively instead, the compass
spins forward on menu open and reverse on close, the sun-moon on theme toggle, because hover
rarely fires on touch screens. Those pass `isAnimated={false}` and call `startAnimation` through
a ref, guarded on `prefers-reduced-motion`.

`test/iconAuthority.test.ts` enforces the surface, and `lucide-react`, `motion/react` and
`framer-motion` are no longer dependencies. The hand-written `compass.tsx` and `sun-moon.tsx`
were deleted: they were a third icon system built on `motion`, and removing that dependency made
**every route's bundle smaller** while adding animation to nineteen icons.

### Typography that is not emoji

These are deliberate and allowed: arrows (`→ ← ↑ ↓ ↔ ⇒ ⇔ ⟶ ⟹`), check and ballot marks
(`✓ ✗ ✔ ✕ ✖`), decorative stars (`✦ ✧`), bullets (`· •`), floor/ceiling brackets
(`⌈ ⌉ ⌊ ⌋`), mathematical symbols (`∈ ∑ ∏ √ ∞ ≈ × ÷ ≠ ≥ ≤ ° ² ′`) and punctuation
(`…, – § †`). They are typography with Unicode codepoints Mathua chose, not pictographs, and
some of them, `⌈`/`⌉` especially, appear inside generated mathematics where "clearing an
emoji" would corrupt a set.

**Circled letters are not in that list, and they were removed from the corpus on 2026-10-15.**
OpenStax wrote a two-part example as `ⓐ … ⓑ …` and the ingester copied the codepoint straight
through: 2340 of them across 72 lessons, in example text and inside figure captions. At lesson size
`ⓐ` is an `a` inside a circle, which is what `@` is, so a reader sees a stray email sign in front
of a sentence. No theme font carries the codepoint, so the browser substitutes a face at a
different size and baseline, the same objection as an emoji. And the corpus already spelled the
parts out in prose: `part (a)` appears 22 times next to a marker that still said `ⓐ`. They are now
`(a)` … `(e)`, and rule 17 of `scripts/audit_lessons.py` fails on any enclosed-alphanumeric glyph
in lesson text.

### How it is checked

`web/next-app/test/noEmoji.test.ts` scans every authored file and fails on any emoji outside
the list above. It also asserts that it is scanning the right files at all, and that its own
fixtures are written as escape sequences, because a guard that reads nothing reports nothing,
and one that flags its own fixtures gets deleted rather than fixed.
---

## 2. Sections 2 through 5 were removed and the numbers were not reused

The rules below are numbered 1, 6 and 7. Sections 2 through 5 existed once and were removed.
Their numbers were left alone on purpose, so that a rule cited as "AGENTS.md §6" in a review or
a commit message still points at the section that author read.

Do not renumber, and do not reuse 2 through 5 for new rules. If a rule is added, give it the
next unused number.

---

## 6. Two tiers of prose, and the split is a decision

Mathua's documentation is written in **two registers**, and both are deliberate.

**Procedural and onboarding prose follows Simplified Technical English.** It covers `CONTRIBUTING.md`,
`docs/mobile-device-checklist.md`, `docs/architecture.md` and the instructional sections of
`README.md`. STE is a controlled language from the aerospace industry: short sentences, active
voice, one instruction each, no semicolons, no contractions, no "should" or "may". It exists
because a newcomer following steps under time pressure must not have to parse a long sentence.

**Analytical prose does not.** `CONTEXT.md` and the ADRs, `docs/mastery-*.md`,
`docs/system-design.md`, `docs/administration.md` and `docs/lesson-media.md` keep their current
register. Their value is the evidence in them: measured ratios such as "853 of 875 edges", the
rejection of two wrong measurement methods and a published retraction of one of my own wrong
numbers. The STE rules strip exactly that. A retraction that reads "the measurement was wrong"
instead of naming the Kahn bug that caused it is not more readable, it is less useful.

`scripts/check_ste.py` enforces the first tier. It is a **derived structural subset, not
ASD-STE100 conformance**, and must never be described as conformance, the upstream skill says of
itself that it "does not certify compliance with ASD-STE100". Two reasons it is a
reimplementation rather than a vendored call:

- **Licensing.** MIT covers that skill's own text. The ASD-STE100 *dictionary* is ASD's property
  and is explicitly not MIT. Copying the word list into this repo would put ASD's dictionary in
  our tree.
- **The dictionary is the wrong tool here.** A mathematics and pedagogy codebase is mostly
  technical names, concept, learner, mastery, scheduler, prerequisite, which STE permits
  anyway. A vocabulary gate would either nag constantly or need an allowlist wide enough to be
  meaningless. The structural rules are the part that improves readability, and they need no
  dictionary.

Moving a document between tiers is a decision to record here, not one to make silently.

**Recorded moves.** `docs/administration.md` was split out of `docs/architecture.md` on
2026-10-12 and placed in the analytical tier. The administration section was 68% of
`architecture.md` and it is design rationale, "a count read before a write loses a race", which
is the register STE strips. Keeping it in the STE file meant checking it against rules written
for procedure. The overview that remained, the five-layer stack and the diagnostic, is
descriptive and stays in the STE tier. `check_ste.py --mode descriptive docs/administration.md`
reports a violation, which is the expected result for a document in this tier.

---

## 7. A documented number and the constant it describes are one fact

If the code owns a number, the docs may quote it but never own it. Same for corpus counts, which
come from the corpus.

**Why:** README stated a 150 XP mastery-check gate twice while `internal/xp/policy.go` said 50. The
same figure was wrong on the landing page until it was fixed separately. Nothing was broken and no
test failed, because a string in Markdown cannot break anything, which is exactly why it needs a
check. The same class of bug put "1-3 knowledge points" in CONTRIBUTING against "3 KPs per concept"
in README, with all 2000 shards holding exactly 3 and nothing enforcing either.

**How:** `scripts/check_docs.py` reads each fact from its owner, the Go constant, the corpus, the
shards and fails on any doc that disagrees. **It must never hardcode a fact it can read**, because
a checker that pins its own copy of the number goes stale in exactly the way it was written to
prevent.

One judgement call it makes: a number next to "concepts" is only checked when the sentence also
scopes it to the whole corpus. "calculus is 128 concepts" is a legitimate subset claim that no
checker can verify without knowing which subset is meant.

---

## 8. No Oxford comma in copy we author

**Join the last two items of a list with "and", never ", and".** So: "tokens, passwords and
sessions", not "tokens, passwords, and sessions". The same applies to "or".

### Why

This is a voice decision, not a grammar one. Both forms are correct English; mixing them is what
reads as unedited. Mathua's copy is short, mostly one sentence at a time, and a list that ends in
", and" pulls the eye to the conjunction as though it were a new clause. The owner asked for one
consistent form, and consistency is the whole benefit, a rule applied to half the sentences
delivers none of it.

It is deliberately narrow: a comma before "and" that joins two independent clauses is correct and
stays. "The pool is empty, and the queue says so" is two clauses, not a list. Only the final
separator of a list of three or more changes.

### How it is checked

`scripts/check_oxford_comma.py` scans authored prose, `AGENTS.md`, `CONTEXT.md`, `CONTRIBUTING.md`,
`README.md`, `docs/`, `web/next-app/app`, `web/next-app/components` and `web/next-app/lib` and
reports a candidate whenever a line contains ", and" or ", or" after an earlier comma.

It reports; it does not block. The heuristic cannot tell a serial comma from a clause comma, so a
hard failure would either be wrong or force contorted sentences, and a checker that is wrong gets
disabled rather than obeyed. The rule is enforced in review; the script exists so a reviewer has
the candidates in front of them instead of relying on noticing. It is wired into the advisory
`Documentation (advisory)` CI job for the same reason the documentation checks are advisory: a
prose nit must never hold a release.
