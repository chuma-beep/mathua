# AGENTS.md — Mathua

Hard rules for anyone or anything editing this repository. Each one states the rule, why it
exists, and how it is checked, because a rule nobody can verify is an aspiration.

---

## 1. No emoji anywhere we author

**Emoji are not permitted in any file we author** — application source, components, tests,
scripts, Go, or documentation. This is not a style preference; see the reasons below.

### Why

- **An emoji is not a theme surface.** Its glyph, color and baseline come from the *font*, so it
  cannot inherit `currentColor`, cannot follow light and dark mode, and cannot be sized
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
  library directly — the wrapper is what pins the size to 16px and the color to `currentColor`,
  and the animated icons render a `<div>` whose default color is not `currentColor`, so an icon
  added outside the wrapper does not follow light and dark mode and does not keep the collapsed
  44px sidebar rail even.
- **Text.** A label beats a symbol most of the time. In tables, a word ("Yes", "Not stored")
  beats a mark.

Icons animate on mount and again on hover. Two are driven imperatively instead — the compass
spins forward on menu open and reverse on close, the sun-moon on theme toggle — because hover
rarely fires on touch screens. Those pass `isAnimated={false}` and call `startAnimation` through
a ref, guarded on `prefers-reduced-motion`.

`test/iconAuthority.test.ts` enforces the surface, and `lucide-react`, `motion/react` and
`framer-motion` are no longer dependencies. The hand-written `compass.tsx` and `sun-moon.tsx`
were deleted: they were a third icon system built on `motion`, and removing that dependency made
**every route's bundle smaller** while adding animation to nineteen icons.

### Typography that is not emoji

These are deliberate and allowed: arrows (`→ ← ↑ ↓ ↔ ⇒ ⇔ ⟶ ⟹`), check and ballot marks
(`✓ ✗ ✔ ✕ ✖`), decorative stars (`✦ ✧`), bullets (`· •`), floor/ceiling brackets
(`⌈ ⌉ ⌊ ⌋`), mathematical symbols (`∈ ∑ ∏ √ ∞ ≈ × ÷ ≠ ≥ ≤ ° ² ′`), and punctuation
(`… — – § †`). They are typography with Unicode codepoints Mathua chose, not pictographs, and
some of them — `⌈`/`⌉` especially — appear inside generated mathematics where "clearing an
emoji" would corrupt a set.

### How it is checked

`web/next-app/test/noEmoji.test.ts` scans every authored file and fails on any emoji outside
the list above. It also asserts that it is scanning the right files at all, and that its own
fixtures are written as escape sequences — because a guard that reads nothing reports nothing,
and one that flags its own fixtures gets deleted rather than fixed.
---

## 6. Two tiers of prose, and the split is a decision

Mathua's documentation is written in **two registers**, and both are deliberate.

**Procedural and onboarding prose follows Simplified Technical English.** It covers `CONTRIBUTING.md`,
`docs/mobile-device-checklist.md`, `docs/architecture.md`, and the instructional sections of
`README.md`. STE is a controlled language from the aerospace industry: short sentences, active
voice, one instruction each, no semicolons, no contractions, no "should" or "may". It exists
because a newcomer following steps under time pressure must not have to parse a long sentence.

**Analytical prose does not.** `CONTEXT.md` and the ADRs, `docs/mastery-*.md`,
`docs/system-design.md`, and `docs/lesson-media.md` keep their current register. Their value is the
evidence in them: measured ratios such as "853 of 875 edges", the rejection of two wrong
measurement methods, and a published retraction of one of my own wrong numbers. The STE rules
strip exactly that. A retraction that reads "the measurement was wrong" instead of naming the
Kahn bug that caused it is not more readable, it is less useful.

`scripts/check_ste.py` enforces the first tier. It is a **derived structural subset, not
ASD-STE100 conformance**, and must never be described as conformance — the upstream skill says of
itself that it "does not certify compliance with ASD-STE100". Two reasons it is a
reimplementation rather than a vendored call:

- **Licensing.** MIT covers that skill's own text. The ASD-STE100 *dictionary* is ASD's property
  and is explicitly not MIT. Copying the word list into this repo would put ASD's dictionary in
  our tree.
- **The dictionary is the wrong tool here.** A mathematics and pedagogy codebase is mostly
  technical names — concept, learner, mastery, scheduler, prerequisite — which STE permits
  anyway. A vocabulary gate would either nag constantly or need an allowlist wide enough to be
  meaningless. The structural rules are the part that improves readability, and they need no
  dictionary.

Moving a document between tiers is a decision to record here, not one to make silently.

---

## 7. A documented number and the constant it describes are one fact

If the code owns a number, the docs may quote it but never own it. Same for corpus counts, which
come from the corpus.

**Why:** README stated a 150 XP mastery-check gate twice while `internal/xp/policy.go` said 50. The
same figure was wrong on the landing page until it was fixed separately. Nothing was broken and no
test failed, because a string in Markdown cannot break anything — which is exactly why it needs a
check. The same class of bug put "1-3 knowledge points" in CONTRIBUTING against "3 KPs per concept"
in README, with all 657 shards holding exactly 3 and nothing enforcing either.

**How:** `scripts/check_docs.py` reads each fact from its owner — the Go constant, the corpus, the
shards — and fails on any doc that disagrees. **It must never hardcode a fact it can read**, because
a checker that pins its own copy of the number goes stale in exactly the way it was written to
prevent.

One judgement call it makes: a number next to "concepts" is only checked when the sentence also
scopes it to the whole corpus. "calculus is 128 concepts" is a legitimate subset claim that no
checker can verify without knowing which subset is meant.
