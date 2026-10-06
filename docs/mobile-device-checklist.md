# Mobile device checklist

Manual verification for the MathLive answer editor. Everything here is **not** covered by
Playwright, and the reason is worth stating plainly before the list rather than after it.

## Why this file exists

`e2e/mobile-math.spec.ts` runs under `mobile-chromium` and `mobile-320`, which are
emulated viewports with synthetic touch. They are genuinely useful. They caught the
backspace key typing its own name. They caught the zero key typing `[0]`. They caught a
loading state that swallowed every tap. They are not a substitute for a device.

They cannot reproduce:

| behavior | why emulation misses it |
|---|---|
| native / OS keyboard suppression | Chromium has no OS keyboard to suppress |
| iOS Safari viewport units | `100vh`, `dvh` and the URL bar behave differently, and only on Safari |
| safe-area insets | a desktop browser reports zero insets |
| keyboard open/close animation | instantaneous in emulation |
| iOS focus zoom | a 16px threshold Safari enforces and Chromium does not |
| real CPU and parse cost | a desktop core, and a cached MathLive chunk |

Two consequences recorded in `docs/math-input.md`, both of which make the automated
numbers a **floor** rather than a measurement:

- The no-input window measured **322ms** with the MathLive chunk cached and **3116ms**
  with a 3s chunk delay. The 310ms floor is parse-and-evaluate of a 794 kB decoded bundle
  on a desktop CPU. On a mid-range phone both terms are larger.
- Network contribution is **unmeasured**. CDP throttling cannot reach the chunk in this
  harness. Chromium serves the chunk from its memory cache, and it reports
  `transferSize: 0` regardless of the applied profile.

So: a green suite means the wiring is right. It does not mean the experience is.

## How to use this file

Run each block on a real device, in a real browser, not a desktop device-emulation mode.
Record the device, OS, browser and build, then tick or note. **A failure here outranks a
passing suite.**

If something fails that this file does not cover, that is a gap in this file. Add it.

---

## Device 1 — iPhone, current supported iOS Safari

Recommended: any iPhone still receiving security updates. Record the exact model.

| # | check | expected | result |
|---|---|---|---|
| 1.1 | Open a `/learn` task on a **cold load** (clear the tab, then reload) | | |
| 1.2 | Time from the answer UI appearing to the keypad being usable | record | |
| 1.3 | Tap the answer field **immediately**, before the keypad appears | the tap is not swallowed; a plain input is usable while the editor loads | |
| 1.4 | Type during that window | the keystrokes are accepted, and survive the editor arriving | |
| 1.5 | Type `2+2`, tap Check Answer **once** | submits on the first tap, one request | |
| 1.6 | Does the OS keyboard appear at all? | it must not | |
| 1.7 | Is the question still visible with the keypad up? | yes | |
| 1.8 | Is the Check Answer button reachable with the keypad up? | yes | |
| 1.9 | Tap Check Answer with the keypad up | lands on the first tap | |
| 1.10 | Page zoom on focus | **no zoom** — a 16px threshold Safari enforces | |
| 1.11 | Long-press in the answer field | context menu behaves; selection is not destroyed | |
| 1.12 | Dismiss the keypad (its own keycap) | it closes, and tapping the field brings it back | |
| 1.13 | Tap a non-focusable area (the question text), then tap the field | the keypad comes back on the next tap | |
| 1.14 | Rotate to landscape and back | the field and Check stay reachable | |
| 1.15 | Bottom spacing above the home indicator | not doubled, not clipped | |
| 1.16 | Horizontal scroll | none required to answer | |
| 1.17 | Backspace key | deletes; never inserts the text `[backspace]` | |
| 1.18 | Zero key | inserts `0`; never the text `[0]` | |
| 1.19 | Fraction key | builds a fraction with two navigable slots | |
| 1.20 | Full flow, touch only: open → focus → answer → modify → submit → feedback → continue | completes | |
| 1.21 | Type `1/2+1/3`, pause | the self-check offers "that simplifies to 5/6" | |
| 1.22 | Switch the keypad to `abc` and tap the wide bottom bar | inserts a space; does **not** just move the caret (ADR-035) | |
| 1.23 | Tap the `␣` key on the arithmetic grid | inserts a space; never types the glyph `␣` | |
| 1.24 | Tap every key in the `Algebra` and `Geometry` grids | all at least 44px on a 412px phone (ADR-035) | |
| 1.22 | On mobile data, observe first use | note whether the CE download is acceptable to you | |
| 1.23 | DevTools → Network, filter `.wav` | **no requests** | |
| 1.24 | DevTools → Elements, `:root` | `data-mathua-keyboard="open"` and a non-zero `--mathua-keyboard-height` while the keypad is up | |

## Device 2 — mid-range Android, current supported Chrome

"Mid-range" matters more than "current": the questions are about CPU, and a flagship
hides parse cost. A device from roughly the last 3–4 years.

| # | check | expected | result |
|---|---|---|---|
| 2.1 | Cold load, time to the keypad being usable | record | |
| 2.2 | Tap the field immediately | not swallowed | |
| 2.3 | Type `2+2`, tap Check once | submits on the first tap | |
| 2.4 | Viewport movement when the keypad opens | no jump; the field is not hidden | |
| 2.5 | Does the layout oscillate as the keypad animates? | no | |
| 2.6 | Bottom safe-area / gesture bar | not overlapped | |
| 2.7 | Backspace, zero, fraction keys | correct | |
| 2.8 | Touch-only full flow | completes | |
| 2.9 | **First-focus latency** | record; flag anything that visibly freezes the UI | |
| 2.10 | **First-use MathLive load** | record; this is the parse cost of 794 kB decoded | |
| 2.11 | **First-use Compute Engine load** (after typing `1/2+1/3` and pausing) | record; 787 kB gzip | |
| 2.12 | Does typing stutter while the editor is up? | no | |
| 2.13 | Network → filter `wav` | **no requests** | |
| 2.14 | 320px-equivalent (or a narrow device) | no horizontal scroll; keypad ≤ half the screen | |
| 2.15 | Chrome remote debugging → `Emulation.setDeviceMetricsOverride` is **not** a substitute for 2.1–2.13 | | |

---

## Known open items, for whoever runs this

These are unresolved and belong in front of a device, not a test:

1. **Keypad dismissal after an outside tap.** Playwright cannot reproduce it: under
   Chromium `mathVirtualKeyboard.hide()` is a no-op, so the failure mode is invisible in
   CI. The outside-tap `hide()` was removed (ADR-030) and the `pointerdown` handler is
   the recovery path, but **1.13** is the check that matters.
2. **320px: Check enabled but the tap did not land** within 1500ms in two diagnosis
   scenarios, while the identical flow passed on `mobile-chromium`. Whether the keypad
   covers the button is unresolved. Check **1.9** and **2.8**, and **2.14**.
3. **Is 787 kB of Compute Engine acceptable on mobile data?** It is skipped on
   `saveData` and on 2g, but *not* on 3g or 4g. If the answer is no, deleting
   `components/math/SelfCheck.tsx` reverts that feature outright — nothing else in
   `components/math/` imports it.
4. **The real cold-start window.** Every number in `docs/math-input.md` is a floor. **1.2**
   and **2.1** replace them with measurements.

## Recording

One row per device per run, so regressions are visible rather than remembered.

| date | device | OS | browser | MathLive | notes |
|---|---|---|---|---|---|
| | | | | 0.111.0 | |
