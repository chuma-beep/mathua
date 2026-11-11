import { render } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { answerCard, answerField } from './helpers/answerInput'

// The math editor's chunk resolves after the first render and React replaces the
// loading input with the editor in place. A test that typed into the loading
// input and then asks for its card is holding a node that has left the document,
// and `closest` on a detached node cannot see the card above it. That threw "no
// question card found" in CI, on a machine slow enough for the replacement to
// land between typing and the click.
describe('answerCard', () => {
  it('finds the card when the field it was given has been replaced', () => {
    const { container } = render(
      <div>
        <div className="border bg-mathua-surface p-6" data-testid="card">
          <input placeholder="Loading" readOnly />
        </div>
      </div>,
    )
    const card = container.querySelector('[data-testid="card"]') as HTMLElement
    const stale = card.querySelector('input') as HTMLInputElement

    stale.remove()
    const live = document.createElement('input')
    live.placeholder = 'Your answer'
    card.appendChild(live)

    expect(document.contains(stale)).toBe(false)
    expect(answerCard(stale)).toBe(card)
    expect(answerField()).toBe(live)
  })
})
