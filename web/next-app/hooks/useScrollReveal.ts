import { useEffect } from 'react'

// Reveal below-the-fold sections as they scroll into view. Elements opt in
// with `data-reveal`; CSS supplies the hidden state and the transition, and
// this hook flips `.is-visible` once. No IntersectionObserver (or reduced
// motion) means the content is shown immediately.
export function useScrollReveal(enabled = true) {
  useEffect(() => {
    if (!enabled) return
    const els = Array.from(document.querySelectorAll<HTMLElement>('[data-reveal]'))
    if (els.length === 0) return
    if (typeof IntersectionObserver === 'undefined') {
      for (const el of els) el.classList.add('is-visible')
      return
    }
    const observer = new IntersectionObserver(
      entries => {
        for (const entry of entries) {
          if (!entry.isIntersecting) continue
          entry.target.classList.add('is-visible')
          observer.unobserve(entry.target)
        }
      },
      { threshold: 0.15 }
    )
    for (const el of els) observer.observe(el)
    return () => observer.disconnect()
  }, [enabled])
}
