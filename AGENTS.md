# AGENTS.md — Mathua

Hard rules for anyone or anything editing this repository. Each one states the rule, why it
exists, and how it is checked, because a rule nobody can verify is an aspiration.

---

## 1. No emoji anywhere we author

**Emoji are not permitted in any file we author** — application source, components, tests,
scripts, Go, or documentation. This is not a style preference; see the reasons below.

### Why

- **An emoji is not a theme surface.** Its glyph, colour and baseline come from the *font*, so it
  cannot inherit `currentColor`, cannot follow light and dark mode, and cannot be sized
  consistently against the mono type Mathua's UI uses. Everything else in the app is a theme
  token; an emoji is the one thing that is not.
- **They render differently per platform.** The same codepoint is a different picture on macOS,
  iOS, Windows and Android, so a UI that looks right on a developer's machine is not the UI a
  learner sees. That is how the `⏵` on the graph's flow toggle ended up at an unpredictable
  baseline next to the word "Flow" inside a button already set to `fontSize: 10`.
- **They are not accessible by default.** They are announced by name ("hourglass"), which on a
  line that already says "3 concepts due for review" is noise, and they vary in how they
  degrade. Icons get `aria-hidden` when decorative.

### What to use instead

- **An icon component.** `@animateicons/react` (lucide set) for animated icons,
  `lucide-react` for static ones, and `components/icons/` for the hand-written animated ones
  already there (`compass.tsx`, `sun-moon.tsx`). Decorate decorative icons with `aria-hidden`.
- **Text.** A label beats a symbol most of the time. In tables, a word ("Yes", "Not stored")
  beats a mark.

If an icon is needed and neither library has it, add it to `components/icons/` following
`compass.tsx` rather than reaching for a glyph.

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