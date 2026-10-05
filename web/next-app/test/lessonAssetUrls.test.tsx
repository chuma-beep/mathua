import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import KatexContent, { resolveLessonImageUrl } from '../components/KatexContent'

/**
 * `/study` is the reference view: the whole lesson, figures in place.
 *
 * The lesson body is canonical and `/learn` selects figures out of it (see LessonAssets).
 * That split only holds if rendering the body still renders every figure in document order —
 * so these assert the reference surface directly rather than inferring it from `/learn`.
 */
describe('reference lesson figures', () => {
  it('interleaves figures with the prose that surrounds them', () => {
    const body = [
      'The angle is positive in the first quadrant.',
      '',
      '![Unit circle with notable angles](svg/unit-circle-labeled.svg)',
      '',
      'The length of $PR$ is the sine of the angle.',
    ].join('\n')

    const { container } = render(<KatexContent>{body}</KatexContent>)

    const imgs = Array.from(container.querySelectorAll('img'))
    expect(imgs).toHaveLength(1)
    expect(imgs[0].getAttribute('src')).toBe('/diagrams/algebrica/unit-circle-labeled.svg')
    expect(imgs[0].getAttribute('alt')).toBe('Unit circle with notable angles')

    // Document order: prose, then the figure, then more prose. A figure appended after the
    // lesson instead of placed inside it reads as an afterthought.
    const text = Array.from(container.querySelectorAll('p, img')).map(
      (el) => el.tagName === 'IMG' ? `[figure:${el.getAttribute('src')}]` : (el.textContent ?? ''),
    )
    const firstProse = text.findIndex((t) => t.includes('positive in the first quadrant'))
    const figure = text.findIndex((t) => t.startsWith('[figure:'))
    const lastProse = text.findIndex((t) => t.includes('is the sine of the angle'))
    expect(firstProse).toBeGreaterThanOrEqual(0)
    expect(figure).toBeGreaterThan(firstProse)
    expect(lastProse).toBeGreaterThan(figure)
  })

  it('keeps every figure in a lesson that has many', () => {
    // 254 figures across the corpus, none of which a learner should have to notice missing.
    const parts: string[] = []
    for (let i = 1; i <= 12; i++) {
      parts.push(`Step ${i}.`, '', `![Figure ${i}](svg/fig-${i}.svg)`, '')
    }
    const { container } = render(<KatexContent>{parts.join('\n')}</KatexContent>)
    const srcs = Array.from(container.querySelectorAll('img')).map((i) => i.getAttribute('src'))
    expect(srcs).toEqual(Array.from({ length: 12 }, (_, i) => `/diagrams/algebrica/fig-${i + 1}.svg`))
  })

  it('falls back to the file name when alt text is absent', () => {
    // 133 of the corpus's figures are `![]`. An unlabelled graphic is worse than a clumsy
    // one, so the name is used rather than nothing.
    const { container } = render(<KatexContent>{'![](svg/floor-and-ceiling-functions-1.svg)'}</KatexContent>)
    expect(container.querySelector('img')?.getAttribute('alt')).toBe('floor-and-ceiling-functions-1.svg')
  })
})

describe('resolveLessonImageUrl', () => {
  // Mirrored by `internal/lessons/assets_test.go`, which gates the corpus against a body
  // referencing a figure that is not shipped. Both copies must agree or the reference view
  // and the Learn view resolve the same figure differently.
  it.each([
    ['svg/unit-circle.svg', '/diagrams/algebrica/unit-circle.svg'],
    ['/diagrams/algebrica/x.png', '/diagrams/algebrica/x.png'],
    ['https://example.org/x.svg', 'https://example.org/x.svg'],
    ['http://example.org/x.svg', 'http://example.org/x.svg'],
    ['x.svg', '/diagrams/algebrica/x.svg'],
    ['deep/nested/dir/x.svg', '/diagrams/algebrica/x.svg'],
    [undefined, ''],
    ['', ''],
  ])('%s -> %s', (input, expected) => {
    expect(resolveLessonImageUrl(input)).toBe(expected)
  })
})