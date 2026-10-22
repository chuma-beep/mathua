// Focus the control that advances to the next question, clear of the floating tab bar.
//
// The answer hosts used `autoFocus` on Next. Focus scrolls an element into view only when
// its border box is outside the viewport, and it ignores `scroll-margin` when deciding
// that. On a 320x568 phone Next was fully on screen but under the tab bar's chevron
// toggle, so nothing scrolled and the tap meant for Next landed on "Hide navigation".
// `scrollIntoView` measures the box expanded by `scroll-margin`, which `.clear-of-tabs`
// sets below lg (app/globals.css), so it moves Next up when it is too close to the bar.
//
// Pass as a callback ref. It is a stable module function, so React calls it once on
// mount and once with null on unmount.
export function focusClearOfTabs(el: HTMLElement | null): void {
  if (!el) return
  el.focus({ preventScroll: true })
  // jsdom does not implement scrollIntoView.
  el.scrollIntoView?.({ block: 'nearest' })
}
