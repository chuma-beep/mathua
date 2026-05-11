# Mathua Design System — Modern Academic

A mathematics learning platform design system bridging scholarly typographic tradition with contemporary dark-mode-first web aesthetics. Monochrome-precise surfaces with gold theorem accents, serif structural headings, and woven ASCII art.

---

## Design Philosophy

**Modern Academic** — Like a well-typeset mathematics journal redesigned for the digital age. Scholarly rigor meets contemporary web design. Dark-first, gold-accented, serif-headed. Precision above decoration.

---

## Color Palette

### Dark Mode (default)

| Token | Hex | Usage |
|---|---|---|
| `--bg` | `#0a0f1a` | Page background |
| `--surface` | `#0f172a` | Cards, raised panels |
| `--surface-elevated` | `#1e293b` | Hover states, code blocks |
| `--surface-highlight` | `#334155` | Active elements, strong borders |
| `--text-primary` | `#f8fafc` | Headings, body text |
| `--text-secondary` | `#94a3b8` | Secondary text, descriptions |
| `--text-muted` | `#64748b` | Labels, metadata, captions |
| `--accent-gold` | `#c8a96e` | **Theorem gold** — mastered states, hero accents |
| `--accent-gold-hover` | `#d4b87a` | Gold hover state |
| `--accent-blue` | `#3b82f6` | **Integral blue** — links, active nav, CTA buttons |
| `--accent-blue-hover` | `#60a5fa` | Blue hover state |
| `--accent-green` | `#10b981` | **Proof green** — success, correct answers, mastered |
| `--accent-red` | `#ef4444` | Errors, incorrect |
| `--border` | `#334155` | Card borders |
| `--border-strong` | `#475569` | Strong separators |
| `--code-bg` | `#1e293b` | Code/formula block backgrounds |
| `--division` | `#1e293b` | Section dividers |

### Light Mode

| Token | Hex | Usage |
|---|---|---|
| `--bg` | `#fefcf4` | Page background — warmest paper |
| `--surface` | `#faf8f0` | Cards |
| `--surface-elevated` | `#f5f1e6` | Hover states, code blocks |
| `--surface-highlight` | `#e8e4d6` | Active elements |
| `--text-primary` | `#1e293b` | Headings, body |
| `--text-secondary` | `#475569` | Descriptions |
| `--text-muted` | `#94a3b8` | Labels, captions |
| `--accent-gold` | `#b8933e` | Theorem gold |
| `--accent-gold-hover` | `#a07a2e` | Gold hover |
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

| Role | Font Family | Weight | Size / Line |
|---|---|---|---|
| Hero Title | `'Cormorant Garamond', serif` | 600 | 56px / 1.1 |
| Page Title | `'Cormorant Garamond', serif` | 500 | 36px / 1.2 |
| Section Heading | `'Cormorant Garamond', serif` | 500 | 28px / 1.3 |
| Card Heading | `'Cormorant Garamond', serif` | 400 | 20px / 1.4 |
| Body | `'Inter', -apple-system, sans-serif` | 400 | 16px / 1.7 |
| Body Small | `'Inter', -apple-system, sans-serif` | 400 | 14px / 1.6 |
| Strong | `'Inter', -apple-system, sans-serif` | 600 | — |
| Table Headers | `'JetBrains Mono', monospace` | 400 | 11px / 1.5 |
| Metric Numbers | `'JetBrains Mono', monospace` | 500 | 32px / 1 |
| Mono Labels | `'JetBrains Mono', monospace` | 400 | 11px / 1.5 |
| Mono Code | `'JetBrains Mono', monospace` | 400 | 13px / 1.8 |
| Stat Badges | `'JetBrains Mono', monospace` | 500 | 11px / 1.5 |

---

## Spacing & Sizing

| Token | Value |
|---|---|
| Container width | `1100px` |
| Section padding (y) | `80px` desktop, `48px` mobile |
| Card padding | `24px` |
| Grid gap | `24px` |
| Tight gap | `12px` |
| Button height | `40px` |
| Button padding (x) | `24px` |

---

## Border Radius

| Element | Value |
|---|---|
| Cards | `8px` |
| Buttons | `6px` |
| Pills / Badges | `9999px` |
| Code blocks | `6px` |
| Inputs | `6px` |

---

## Shadows

| Token | Value |
|---|---|
| Card hover | `0 4px 12px rgba(0, 0, 0, 0.15)` |
| Elevated | `0 8px 24px rgba(0, 0, 0, 0.2)` |

---

## ASCII Art Animations

| Animation | Trigger | Effect |
|---|---|---|
| `ascii-reveal` | On scroll / load | Character-by-character typewriter reveal via `steps()` |
| `ascii-pulse` | Continuous | Opacity pulse on divider lines |
| `progress-fill` | On data change | Width fill from 0 to target |

---

## Tailwind Mapping

DESIGN.md tokens are mapped to Tailwind via `theme.extend.colors.mathua.*`:

| Tailwind Class | DESIGN.md Token |
|---|---|
| `bg-mathua-bg` | `--bg` |
| `bg-mathua-surface` | `--surface` |
| `bg-mathua-surface-elevated` | `--surface-elevated` |
| `text-mathua-primary` | `--text-primary` |
| `text-mathua-secondary` | `--text-secondary` |
| `text-mathua-muted` | `--text-muted` |
| `text-mathua-gold` / `bg-mathua-gold` | `--accent-gold` |
| `text-mathua-blue` / `bg-mathua-blue` | `--accent-blue` |
| `text-mathua-green` / `bg-mathua-green` | `--accent-green` |
| `border-mathua-border` | `--border` |
| `bg-mathua-code` | `--code-bg` |

Keyframe animations:
- `animate-ascii-pulse` — opacity pulse
- `animate-ascii-reveal` — typewriter reveal
- `animate-progress-fill` — width fill

---

## Component Patterns

### Button
- `.btn-primary` → `bg-mathua-blue text-white rounded-md h-10 px-6 font-medium`
- `.btn-outline` → `border border-mathua-border-strong text-mathua-primary rounded-md h-10 px-6 font-medium hover:border-mathua-blue hover:text-mathua-blue`
- `.btn-ghost` → `text-mathua-secondary rounded-md h-10 px-6 font-medium hover:bg-mathua-surface-elevated`

### Card
- `bg-mathua-surface border border-mathua-border rounded-lg p-6`

### InfoCard
- Card base + `border-l-[3px] border-l-mathua-blue`

### Code/Formula Block
- `bg-mathua-code border border-mathua-border rounded-md p-5 font-mono text-mathua-blue text-xs whitespace-pre overflow-x-auto`

### Pipeline State Pill
- UNSEEN / LEARNING / PRACTICING: `bg-mathua-surface border border-mathua-border text-mathua-secondary`
- MASTERED: `bg-mathua-green text-white border-mathua-green`
- DECAYING: `bg-mathua-gold text-mathua-bg`

### ASCII Banner
- `font-mono text-mathua-blue` with `animate-ascii-reveal` on load
