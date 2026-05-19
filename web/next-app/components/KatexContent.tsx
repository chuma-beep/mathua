'use client'

import ReactMarkdown from 'react-markdown'
import remarkMath from 'remark-math'
import rehypeKatex from 'rehype-katex'
import remarkGfm from 'remark-gfm'
import 'katex/dist/katex.min.css'

export default function KatexContent({ children }: { children: string }) {
  let content = children

  // Convert lesson LaTeX delimiters from Markdown-escaped form to $...$ / $$...$$
  // The lesson files use \\(...\\) (inline) and \\[...\\] (display),
  // which remark-math expects as $...$ and $$...$$.
  // Inside math content, also unescape \\ → \ so that \\{ becomes \{, \\sin becomes \sin, etc.
  content = content.replace(/\\\\\(([\s\S]*?)\\\\\)/g, (_, inner) => '$' + inner.replace(/\\\\/g, '\\') + '$')
  content = content.replace(/\\\\\[([\s\S]*?)\\\\\]/g, (_, inner) => {
    const clean = inner.replace(/\\\\/g, '\\')
    return '\n$$\n' + clean + '\n$$\n'
  })

  return (
    <div className="katex-content text-sm leading-relaxed">
      <ReactMarkdown
        remarkPlugins={[remarkMath, remarkGfm]}
        rehypePlugins={[rehypeKatex]}
        components={{
          a: ({ href, children }) => (
            <a href={href} target="_blank" rel="noopener noreferrer"
               className="text-mathua-blue hover:underline">
              {children}
            </a>
          ),
          code: ({ children }) => (
            <code className="bg-mathua-code px-1 rounded-none text-mathua-secondary">
              {children}
            </code>
          ),
        }}
      >
        {content}
      </ReactMarkdown>
      <style jsx global>{`
        .katex-content h1 { font-size: 1.5rem; font-weight: 700; margin: 1.5rem 0 0.75rem; font-family: var(--font-serif); }
        .katex-content h2 { font-size: 1.25rem; font-weight: 700; margin: 1.25rem 0 0.5rem; font-family: var(--font-serif); }
        .katex-content h3 { font-size: 1.1rem; font-weight: 600; margin: 1rem 0 0.5rem; }
        .katex-content p { margin: 0.75rem 0; }
        .katex-content ul, .katex-content ol { margin: 0.5rem 0; padding-left: 1.5rem; }
        .katex-content li { margin: 0.25rem 0; }
        .katex-content hr { border: 0; border-top: 1px solid; margin: 1.5rem 0; opacity: 0.3; }
        .katex-content strong { font-weight: 700; }
      `}</style>
    </div>
  )
}
