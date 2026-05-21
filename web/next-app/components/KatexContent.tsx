'use client'

import React, { ReactNode } from 'react'
import ReactMarkdown from 'react-markdown'
import remarkMath from 'remark-math'
import rehypeKatex from 'rehype-katex'
import remarkGfm from 'remark-gfm'
import 'katex/dist/katex.min.css'

function headingId(text: string): string {
  return text
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-|-$/g, '')
}

function extractText(children: ReactNode): string {
  let text = ''
  React.Children.forEach(children, (child) => {
    if (typeof child === 'string' || typeof child === 'number') {
      text += child
    } else if (child && typeof child === 'object' && 'props' in child) {
      text += extractText((child as any).props.children)
    }
  })
  return text
}

export default function KatexContent({ children }: { children: string }) {
  let content = children

  // Convert lesson LaTeX delimiters to $...$ / $$...$$
  // remark-math expects $...$ for inline and $$...$$ on their own lines for display.
  // Two styles: Algebrica (double backslash: \\(...\\) / \\[...\\]) and Wikipedia (single backslash: \(...\) / \[...\]).

  // --- Inline math ---
  // Algebrica style: \\(...\\) — unescape \\ back to \ inside
  content = content.replace(/\\\\\(([\s\S]*?)\\\\\)/g, (_, inner) => {
    const clean = inner.replace(/\\\\/g, '\\')
    if (clean.includes('\n')) return '\n$$\n' + clean.trim() + '\n$$\n'
    return '$' + clean + '$'
  })
  // Wikipedia style: \(...\) — backslashes inside are already single
  content = content.replace(/\\\(([\s\S]*?)\\\)/g, (_, inner) => {
    if (inner.includes('\n')) return '\n$$\n' + inner.trim() + '\n$$\n'
    return '$' + inner + '$'
  })

  // --- Display math ---
  // Use block $$...$$ (own lines) when multi-line,
  // and inline $\displaystyle ...$ when single-line (e.g. inside table cells).
  // Algebrica style: \\[...\\]
  content = content.replace(/\\\\\[([\s\S]*?)\\\\\]/g, (_, inner) => {
    const clean = inner.replace(/\\\\/g, '\\')
    if (inner.includes('\n')) {
      return '\n$$\n' + clean.trim() + '\n$$\n'
    }
    return '$\\displaystyle ' + clean.trim() + '$'
  })
  // Wikipedia style: \[...\]
  content = content.replace(/\\\[([\s\S]*?)\\\]/g, (_, inner) => {
    if (inner.includes('\n')) {
      return '\n$$\n' + inner.trim() + '\n$$\n'
    }
    return '$\\displaystyle ' + inner.trim() + '$'
  })

  return (
    <div className="katex-content text-sm leading-relaxed">
      <ReactMarkdown
        remarkPlugins={[remarkMath, remarkGfm]}
        rehypePlugins={[rehypeKatex]}
        components={{
          a: ({ children }) => <>{children}</>,
          code: ({ children }) => (
            <code className="bg-mathua-code px-1 rounded-none text-mathua-secondary">
              {children}
            </code>
          ),
          img: ({ src, alt }) => {
            let url = src
            if (url && !url.startsWith('http') && !url.startsWith('/')) {
              url = '/diagrams/algebrica/' + url.split('/').pop()
            }
            return <img src={url} alt={alt || ''} className="max-w-full h-auto my-4 mx-auto" />
          },
          h1: ({ children, ...props }) => {
            const text = extractText(children)
            return <h1 id={headingId(text)} {...props}>{children}</h1>
          },
          h2: ({ children, ...props }) => {
            const text = extractText(children)
            return <h2 id={headingId(text)} {...props}>{children}</h2>
          },
          h3: ({ children, ...props }) => {
            const text = extractText(children)
            return <h3 id={headingId(text)} {...props}>{children}</h3>
          },
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
        .katex-content h1:hover .anchor-link,
        .katex-content h2:hover .anchor-link,
        .katex-content h3:hover .anchor-link { opacity: 1; }
      `}</style>
    </div>
  )
}
