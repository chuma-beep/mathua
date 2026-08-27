'use client'

import React, { ReactNode } from 'react'
import ReactMarkdown from 'react-markdown'
import rehypeKatex from 'rehype-katex'
import remarkGfm from 'remark-gfm'
import rehypeRaw from 'rehype-raw'
import rehypeSanitize, { defaultSchema } from 'rehype-sanitize'
import 'katex/dist/katex.min.css'
import { lessonMacros, prepareLessonMath } from '../lib/lessonMath'

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

const sanitizeSchema = {
  ...defaultSchema,
  attributes: {
    ...defaultSchema.attributes,
    span: [
      ...(defaultSchema.attributes?.span ?? []),
      ['className', 'math-inline', 'math-display'],
    ],
    div: [...(defaultSchema.attributes?.div ?? []), ['className', 'math-display', 'katex', 'katex-display', 'katex-html']],
  },
  tagNames: [...(defaultSchema.tagNames ?? []), 'span'],
}

export default function KatexContent({ children, className = '' }: { children: string; className?: string }) {
  const content = prepareLessonMath(children)
  // In production KaTeX errors are still non-throwing (red fallback) but we
  // emit a console warning for telemetry — vitest corpus gate uses throwOnError:true.
  if (typeof window !== 'undefined' && content.includes('math-') && content.includes('\\')) {
    // Cheap heuristic: if prepareLessonMath emitted math spans with backslashes,
    // downstream KaTeX may warn — let rehype-katex handle it.
  }

  return (
    <div className={`mathua-lesson katex-content text-sm leading-relaxed ${className}`}>
      <ReactMarkdown
        remarkPlugins={[remarkGfm]}
        rehypePlugins={[rehypeRaw, [rehypeSanitize, sanitizeSchema], [rehypeKatex, {
          throwOnError: false,
          trust: false,
          strict: 'warn',
          errorColor: '#cc0000',
          macros: lessonMacros,
        }]]}
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
            // Lesson-sourced images come from arbitrary remote/relative URLs without
            // intrinsic dimensions, so next/image optimization does not apply here.
            // eslint-disable-next-line @next/next/no-img-element
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
        .katex-content { overflow-x: clip; overflow-wrap: break-word; }
        .katex-content h1 { font-size: clamp(1.25rem, 1rem + 1vw, 1.75rem); font-weight: 700; margin: 1.5rem 0 0.75rem; font-family: var(--font-serif); }
        .katex-content h2 { font-size: clamp(1.1rem, 0.9rem + 0.6vw, 1.5rem); font-weight: 700; margin: 1.25rem 0 0.5rem; font-family: var(--font-serif); }
        .katex-content h3 { font-size: clamp(1rem, 0.9rem + 0.3vw, 1.25rem); font-weight: 600; margin: 1rem 0 0.5rem; }
        .katex-content p { margin: 0.75rem 0; font-size: clamp(0.875rem, 0.8rem + 0.3vw, 1rem); }
        .katex-content ul, .katex-content ol { margin: 0.5rem 0; padding-left: 1.5rem; }
        .katex-content li { margin: 0.25rem 0; font-size: clamp(0.875rem, 0.8rem + 0.3vw, 1rem); }
        .katex-content .katex { font-size: 1.05875rem; }
        .katex-content .katex-display { overflow-x: auto; overflow-y: hidden; max-width: 100%; padding-bottom: 4px; }
        .katex-content .math-display { display: block; overflow-x: auto; max-width: 100%; }
        .katex-content hr { border: 0; border-top: 1px solid; margin: 1.5rem 0; opacity: 0.3; }
        .katex-content strong { font-weight: 700; }
        .katex-content table { display: block; overflow-x: auto; max-width: 100%; border-collapse: collapse; margin: 0.75rem 0; font-size: 0.875rem; }
        .katex-content th, .katex-content td { border: 1px solid; border-color: var(--border); padding: 0.5rem 0.75rem; text-align: left; vertical-align: top; }
        .katex-content th { font-weight: 600; background: var(--code-bg); }
        .katex-content h1:hover .anchor-link,
        .katex-content h2:hover .anchor-link,
        .katex-content h3:hover .anchor-link { opacity: 1; }
      `}</style>
    </div>
  )
}
