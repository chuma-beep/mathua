'use client'

import { InlineMath, BlockMath } from 'react-katex'
import 'katex/dist/katex.min.css'

interface Segment {
  type: 'text' | 'inline' | 'display'
  content: string
}

function parseLaTeX(text: string): Segment[] {
  const segments: Segment[] = []

  // Pattern matches \\(...\\) for inline math and \\[...\\] for display math.
  // The lesson files use double-backslash delimiters because they were written
  // for a Markdown pipeline that strips one level of escaping.
  const re = /\\\\\(([\s\S]*?)\\\\\)|\\\\\[([\s\S]*?)\\\\\]/g

  let lastIdx = 0
  let m: RegExpExecArray | null

  while ((m = re.exec(text)) !== null) {
    if (m.index > lastIdx) {
      segments.push({ type: 'text', content: text.slice(lastIdx, m.index) })
    }
    if (m[1] !== undefined) {
      segments.push({ type: 'inline', content: m[1] })
    } else if (m[2] !== undefined) {
      segments.push({ type: 'display', content: m[2] })
    }
    lastIdx = re.lastIndex
  }

  if (lastIdx < text.length) {
    segments.push({ type: 'text', content: text.slice(lastIdx) })
  }

  return segments
}

export default function KatexContent({ children }: { children: string }) {
  const segments = parseLaTeX(children)

  return (
    <span className="katex-content">
      {segments.map((seg, i) => {
        switch (seg.type) {
          case 'inline':
            return <InlineMath key={i} math={seg.content} />
          case 'display':
            return (
              <span key={i} className="block my-4 overflow-x-auto">
                <BlockMath math={seg.content} />
              </span>
            )
          default:
            return (
              <span key={i} className="whitespace-pre-wrap">
                {seg.content}
              </span>
            )
        }
      })}
    </span>
  )
}
