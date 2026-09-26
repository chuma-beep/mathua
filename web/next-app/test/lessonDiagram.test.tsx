import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'

vi.mock('next/image', () => ({
  default: ({ alt }: { alt?: string }) => <img alt={alt ?? ''} />,
}))

const svgDoc = '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><circle cx="50" cy="50" r="40"/></svg>'
const evilSvg = '<svg xmlns="http://www.w3.org/2000/svg"><script>alert(1)</script><circle cx="50" cy="50" r="40" onclick="evil()"/></svg>'

const fetchMock = vi.fn<
  (url: string) => Promise<{ ok: boolean; status: number; text: () => Promise<string>; json?: () => Promise<Record<string, { title: string }>> }>
>()
const metaTitles: Record<string, { title: string }> = {}

function mockFetch(svgDoc: string, ok = true) {
  fetchMock.mockImplementation(async (url: string) => {
    if (url === '/diagrams/meta.json') {
      return { ok: true, status: 200, text: async () => '{}', json: async () => ({ ...metaTitles }) }
    }
    return { ok, status: ok ? 200 : 404, text: async () => svgDoc }
  })
}

type LessonDiagramProps = { src: string; alt: string; maxHeight?: number }

describe('LessonDiagram', () => {
  let LessonDiagram: (props: LessonDiagramProps) => React.JSX.Element

  beforeEach(async () => {
    fetchMock.mockReset()
    vi.stubGlobal('fetch', fetchMock)
    for (const k of Object.keys(metaTitles)) delete metaTitles[k]
    mockFetch(svgDoc)
    // Fresh module per test: the metadata cache is module-scoped.
    vi.resetModules()
    LessonDiagram = (await import('../components/LessonDiagram')).default
  })

  it('inlines sanitized SVG content', async () => {
    const { container } = render(<LessonDiagram src="/diagrams/algebrica/example.svg" alt="Example diagram" />)
    await waitFor(() => expect(container.querySelector('.lesson-diagram svg')).toBeTruthy())
    expect(container.querySelector('.lesson-diagram svg circle')).toBeTruthy()
  })

  it('strips scripts and event handlers from vendored SVGs', async () => {
    mockFetch(evilSvg)
    const { container } = render(<LessonDiagram src="/diagrams/algebrica/evil.svg" alt="Evil diagram" />)
    await waitFor(() => expect(container.querySelector('.lesson-diagram svg')).toBeTruthy())
    expect(container.querySelector('.lesson-diagram script')).toBeNull()
    expect(container.querySelector('.lesson-diagram [onclick]')).toBeNull()
    expect(container.querySelector('.lesson-diagram svg circle')).toBeTruthy()
  })

  it('upgrades alt text from diagram metadata when available', async () => {
    metaTitles['/diagrams/algebrica/example.svg'] = { title: 'Parabola opening upward' }
    render(<LessonDiagram src="/diagrams/algebrica/example.svg" alt="Example diagram" />)
    await waitFor(() => expect(screen.getByLabelText('Parabola opening upward')).toBeTruthy())
  })

  it('renders PNG through the plain image path', () => {
    render(<LessonDiagram src="/diagrams/algebrica/example.png" alt="Example diagram" />)
    const img = screen.getByAltText('Example diagram') as HTMLImageElement
    expect(img.tagName).toBe('IMG')
  })

  it('falls back to the image when the SVG fetch fails', async () => {
    mockFetch('', false)
    render(<LessonDiagram src="/diagrams/algebrica/missing.svg" alt="Missing diagram" />)
    await waitFor(() => expect(screen.getByAltText('Missing diagram')).toBeTruthy())
    expect(screen.getByAltText('Missing diagram').tagName).toBe('IMG')
  })
})
