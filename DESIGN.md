# Mathua Design System

A mathematics learning platform design system built on a single serif family (IBM Plex Serif) paired with its companion monospace (IBM Plex Mono). Dark-first, warm-neutral, gold-accented. Terminal toolkit for the desktop TUI; warm parchment for the web. Precision above decoration.

---

## Design Philosophy

**One family, two voices.** IBM Plex Serif carries all text — headings and body use weight alone for hierarchy. IBM Plex Mono carries code, labels, nav, and structure. The two were designed together by Bold Monday and share the same construction principles. Nothing clashes.

The palette is warm: dark surfaces are `#0b0f1a` (a hint of blue, not pure black) and light surfaces are `#fefcf4` (warmest parchment, never pure white). A single gold accent (`#c8a96e` light / `#b8933e` dark) marks mastery, theorems, and structure. Secondary accents serve semantic roles — green for correct, red for errors, teal for decay.

---

## Color Palette

### Dark Mode (default)

| Token | Hex | Usage |
|---|---|---|
| `--bg` | `#0b0f1a` | Page background |
| `--surface` | `#0f172a` | Cards, raised panels |
| `--surface-elevated` | `#1e293b` | Hover states, code blocks |
| `--surface-highlight` | `#334155` | Active elements |
| `--text-primary` | `#e8e2d5` | Headings, body text |
| `--text-secondary` | `#9fa8b4` | Secondary text, descriptions |
| `--text-muted` | `#5a6577` | Labels, metadata, captions |
| `--accent-gold` | `#c8a96e` | **Theorem gold** — mastery, links, structural accents |
| `--accent-gold-hover` | `#d4b87a` | Gold hover state |
| `--accent-blue` | `#3b82f6` | **Integral blue** — links (legacy), interactive states |
| `--accent-blue-hover` | `#60a5fa` | Blue hover state |
| `--accent-green` | `#10b981` | **Proof green** — correct answers, mastered |
| `--accent-red` | `#ef4444` | **Error red** — incorrect answers |
| `--accent-teal` | `#4db8a0` | **Decay teal** — concepts needing review |
| `--border` | `#1e2d45` | Card borders, dividers |
| `--border-strong` | `#2a3f5f` | Strong separators |
| `--code-bg` | `#1e293b` | Code block backgrounds |
| `--division` | `#1e2d45` | Section dividers |

### Light Mode

| Token | Hex | Usage |
|---|---|---|
| `--bg` | `#fefcf4` | Page background |
| `--surface` | `#faf8f0` | Cards |
| `--surface-elevated` | `#f5f1e6` | Hover states, code blocks |
| `--surface-highlight` | `#e8e4d6` | Active elements |
| `--text-primary` | `#1e293b` | Headings, body |
| `--text-secondary` | `#475569` | Descriptions |
| `--text-muted` | `#94a3b8` | Labels, captions |
| `--accent-gold` | `#b8933e` | Theorem gold |
| `--accent-gold-hover` | `#a07a2e` | Gold hover |
| `--accent-teal` | `#3a9a8a` | Decay teal |
| `--accent-blue` | `#2563eb` | Integral blue |
| `--accent-blue-hover` | `#1d4ed8` | Blue hover |
| `--accent-green` | `#059669` | Proof green |
| `--accent-red` | `#dc2626` | Errors |
| `--border` | `#d6d3c8` | Card borders |
| `--border-strong` | `#b8b5a8` | Strong separators |
| `--code-bg` | `#f0ede0` | Code block backgrounds |
| `--division` | `#d6d3c8` | Section dividers |

---

## Typography

| Role | Font Family | Weight | Size / Line | Responsive |
|---|---|---|---|---|
| Page Title | `'IBM Plex Serif', serif` | 400 | `clamp(1.5rem, 5vw, 1.9rem)` / 1.2 | `max-sm:text-[1.5rem]` (SectionHeader) |
| Section Heading | `'IBM Plex Serif', serif` | 400 | `1.9rem` / 1.2 | `max-sm:text-[1.5rem]` |
| Card Heading | `'IBM Plex Serif', serif` | 400 | `1.1rem` / 1.3 | |
| Body | `'IBM Plex Serif', serif` | 400 | 15px / 1.75 | `max-width: 640px` on mobile, bumps to 16px |
| Body Small | `'IBM Plex Serif', serif` | 400 | 14px / 1.6 | |
| Inline Code | `'IBM Plex Mono', monospace` | 400 | 0.9em | |
| Nav / Labels | `'IBM Plex Mono', monospace` | 400 | 11–12px / 1.5 | |
| Code Blocks | `'IBM Plex Mono', monospace` | 400 | 13px / 1.6 | |
| Metric Numbers | `'IBM Plex Mono', monospace` | 400 | 2rem / 1 | `font-variant-numeric: tabular-nums` |
| ASCII Art | `'IBM Plex Mono', monospace` | 400 | 12px / 1.6 | `overflow-x: auto` |

### Fluid typography pattern

Where font sizes need to adapt to screen width, use `clamp()`:

```css
font-size: clamp(min, preferred, max);
```

Examples:
- Hero heading: `clamp(2.2rem, 5vw, 3.8rem)`
- Body subtext: `clamp(0.95rem, 2.5vw, 1.1rem)`
- Small labels: `clamp(11px, 3vw, 12px)`

---

## Spacing & Layout

| Token | Value | Note |
|---|---|---|
| Container max-width | `1100px` (`max-w-container`) | Applied via `mx-auto px-6 max-sm:px-4` on all pages |
| Section padding (y) | `py-20 max-sm:py-12` | 80px desktop, 48px mobile |
| Code block padding | `16px 20px` | Unified across docs pages |
| Flex gap (default) | `24px` | Between section columns |
| Flex gap (tight) | `12px` | Between inline elements |
| Column gap (stacked mobile) | `16px` (`max-sm:gap-4`) | After stacking to `flex-col` |

---

## Borders

| Element | Style |
|---|---|
| Cards, panels | `0.5px solid var(--border)` |
| Strong separators | `0.5px solid var(--border-strong)` |
| Section dividers | `0.5px solid var(--border)` |
| Gold accents | `2px solid var(--accent-gold)` (left-border on formula blocks) |
| Buttons | No border-radius (`borderRadius: 0`) |
| Code blocks | `border none` with `background: var(--surface)` |

Prefer 0.5px borders throughout — they render as 1px on most displays and feel lighter than full 1px lines. No border-radius on buttons. Code blocks get surface backgrounds instead of borders.

---

## Link Hover Effect

All clickable text links use the `.link-underline` class — a gold line that animates from left to right on hover via a `::after` pseudo-element:

```css
.link-underline {
  position: relative;
  text-decoration: none;
}

.link-underline::after {
  content: '';
  position: absolute;
  left: 0;
  bottom: -2px;
  width: 100%;
  height: 1px;
  background: var(--accent-gold);
  transform: scaleX(0);
  transform-origin: left;
  transition: transform 0.2s ease;
}

.link-underline:hover::after {
  transform: scaleX(1);
}
```

Applied to: Header nav links, Footer links, docs page links, in-page reference links.

---

## Components

### Header

Sticky blurred nav bar at the top of every page. Single component in `components/Header.tsx`, shared across all 8 pages.

- Position: `sticky top-0 z-100`, backdrop-filter: `blur(8px)`
- Background: semitransparent (`rgba(11, 15, 26, 0.95)` dark / `rgba(254, 252, 244, 0.95)` light)
- Font: IBM Plex Mono, 12px links, 13px brand name in `var(--accent-gold)`
- Mobile: horizontal scroll (`overflow-x: auto`) with `white-space: nowrap` — no hamburger, no JS state
- Theme toggle: inline button in bar (not floating), hidden until `mounted`

### Footer

In `components/Footer.tsx`. Links use `display: flex; flex-wrap: wrap` with `gap: 4px 12px` — wraps gracefully on any screen width. Tagline in IBM Plex Serif.

### FormulaBlock

In `components/FormulaBlock.tsx`. Renders formulas with a gold left-border accent. Uses `white-space: pre-wrap` + `max-width: 100%` to prevent horizontal overflow on mobile.

### D2 Diagram

In `components/D2Diagram.tsx`. Maps diagram names to pre-rendered light/dark SVG pairs served from `public/diagrams/`. D2 source files live in `diagrams/*.d2` at the project root. Generated with:

```bash
d2 -t 0  diagrams/name.d2 diagrams/generated/name-light.svg
d2 -t 200 diagrams/name.d2 diagrams/generated/name-dark.svg
```

SVGs are committed to the repo. The component picks the right variant based on the `theme` prop from `useTheme()`.

### Section Header

In `components/SectionHeader.tsx`. Renders a centered page-section heading with a gold left-border block. Title font: IBM Plex Serif, `1.9rem max-sm:1.5rem`.

### Pipeline, DomainTable, ProgressionLevels

All use `width: 100%` and IBM Plex Mono for labels/structure. Pipeline wraps with `flex-wrap`. DomainTable uses `overflow-x: auto` for safe table scrolling.

---

## Diagrams (D2)

All architecture and workflow diagrams are pre-rendered D2 SVGs. No browser-side rendering. No Mermaid runtime.

| Diagram | File | Used in |
|---|---|---|
| Architecture (5-layer stack) | `architecture.d2` | `/docs/architecture` |
| CAT diagnostic (binary search) | `cat-diagnostic.d2` | `/docs/architecture`, `/diagnose` |
| Contributing (PR workflow) | `contributing.d2` | Pre-rendered but not displayed (numbered list is clearer) |
| Student model (state pipeline) | `student-model.d2` | Available for future use |
| Platforms (Web vs Desktop) | `platforms.d2` | Home page "One engine. Two ways to run it." |

Generate with `npm run diagrams`.

---

## Responsive Patterns

| Pattern | Implementation |
|---|---|
| Container padding | `px-6 max-sm:px-4` — 24px → 16px |
| Vertical spacing | `py-20 max-sm:py-12` — 80px → 48px |
| Column stacking | `flex flex-wrap gap-8 max-sm:flex-col max-sm:gap-4` |
| Nav overflow | Horizontal scroll on mobile (`overflow-x: auto`) |
| Footer wrapping | `display: flex; flex-wrap: wrap; gap: 4px 12px` |
| SVG scaling | `max-width: 100%; height: auto` |
| Code blocks | `overflow-x: auto` + `white-space: pre-wrap` for formulas |
| Fluid typography | `clamp()` on font sizes that need to scale |
| Hidden on mobile | `max-sm:hidden` for dividers/elements that break stacking |
| Vertical divider fix | `max-sm:hidden` on `borderLeft` dividers between stacked columns |

---

## Tailwind Configuration

Key customizations in `tailwind.config.js`:

```js
maxWidth: {
  container: '1100px',
},
fontFamily: {
  serif: ['IBM Plex Serif', 'serif'],
  sans: ['IBM Plex Serif', 'serif'],
  mono: ['IBM Plex Mono', 'monospace'],
}
```
