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
      text += extractText((child as React.ReactElement<{ children?: ReactNode }>).props.children)
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

// rehype-katex renders failures two ways: <span class="katex-error"> when
// even the lenient retry throws, and red text (errorColor) when KaTeX itself
// renders the ParseError. Replace both with a labeled <span> fallback so
// users never see red TeX soup, and log once per span so new corpus breakage
// is visible in telemetry instead of silent.
const KATEX_ERROR_COLOR_RE = /#cc0000|204,\s*0,\s*0/i // must match errorColor below
let fontCheckDone = false
function warnFontOnce(): void {
  // KaTeX renders negation slashes (e.g. \neq) with a private-use glyph that
  // exists only in KaTeX_Main: without the webfont users see tofu boxes even
  // though parsing succeeded. Diagnose once instead of failing silently.
  if (fontCheckDone || typeof window === 'undefined') return
  fontCheckDone = true
  try {
    const fonts = (window as unknown as { fonts?: { check(spec: string): boolean } }).fonts
    if (fonts && !fonts.check('16px KaTeX_Main')) {
      // eslint-disable-next-line no-console
      console.warn('[KatexContent] KaTeX_Main webfont not loaded — math symbols may show as boxes')
    }
  } catch {
    // Font check API unavailable; nothing to diagnose.
  }
}
function rehypeKatexFallback() {
  return (tree: unknown) => {
    const walk = (node: any, parent: any, index: number, mathKind: string): void => {
      if (!node || typeof node !== 'object') return
      const classes: unknown = node.properties?.className
      const classList: string[] = Array.isArray(classes) ? (classes as unknown[]).map(String) : []
      // KaTeX replaces math-display/math-inline scopes with its own output:
      // .katex-display marks display mode, plain .katex marks inline.
      const kind = classList.includes('katex-display')
        ? 'display'
        : classList.includes('katex') && !classList.includes('katex-error')
          ? 'inline'
          : mathKind
      const style: unknown = node.properties?.style
      const isErrorNode =
        node.type === 'element' && classList.includes('katex-error')
      const isErrorColor =
        node.type === 'element' &&
        node.tagName === 'span' &&
        typeof style === 'string' &&
        KATEX_ERROR_COLOR_RE.test(style)
      if (
        (isErrorNode || isErrorColor) &&
        parent &&
        Array.isArray(parent.children) &&
        index >= 0
      ) {
        const texts: string[] = []
        const collect = (n: any): void => {
          if (!n || typeof n !== 'object') return
          if (typeof n.value === 'string') texts.push(n.value)
          if (Array.isArray(n.children)) n.children.forEach(collect)
        }
        collect(node)
        const source: string = texts.join('')
        if (typeof window !== 'undefined') {
          warnFontOnce()
          // Provenance for triage: route + span kind + full-ish source. The
          // `\[`-storm investigation showed a bare excerpt is unidentifiable.
          const route = window.location.pathname + window.location.search
          // eslint-disable-next-line no-console
          console.error(
            `[KatexContent] KaTeX failed to render (${kind} @ ${route}); showing source fallback:`,
            source.slice(0, 500)
          )
        }
        const detail: string =
          typeof node.properties?.title === 'string' && node.properties.title
            ? node.properties.title
            : source
        parent.children[index] = {
          type: 'element',
          // NOTE: plain span, not code — the components.code override below
          // would discard a code node's className/title.
          tagName: 'span',
          properties: {
            className: ['math-fallback'],
            title: `Math failed to render (${detail.slice(0, 200)})`,
          },
          children: [{ type: 'text', value: source }],
        }
        return
      }
      if (Array.isArray(node.children)) {
        node.children.forEach((child: unknown, i: number) => walk(child, node, i, kind))
      }
    }
    walk(tree, null, -1, 'math')
  }
}

export default function KatexContent({ children, className = '' }: { children: string; className?: string }) {
  const content = prepareLessonMath(children)
  // In production KaTeX errors are still non-throwing (labeled fallback via
  // rehypeKatexFallback below) but we emit a console warning for telemetry —
  // vitest corpus gate uses throwOnError:true. Note spans are entity-encoded
  // by prepareLessonMath (`&#92;`), so the suspect check below only ever
  // matches prose leaks, never span contents.
  if (typeof window !== 'undefined' && process.env.NODE_ENV === 'development') {
    const suspect = content.match(/<span class="math-(display|inline)">[^<]*\\[^<]*<\/span>/)
    if (suspect) {
      console.warn('[KatexContent] math span with raw backslash may fail KaTeX:', suspect[0].slice(0, 160))
    }
    const proseLeak = content
      .replace(/<span class="math-(display|inline)">[\s\S]*?<\/span>/g, ' ')
      .match(/\\(frac|begin|end|text|cdot|pi)\b/)
    if (proseLeak) {
      console.warn('[KatexContent] possible raw LaTeX in prose:', proseLeak[0])
    }
  }

  return (
    <div className={`mathua-lesson katex-content text-sm leading-relaxed ${className}`}>
      <ReactMarkdown
        remarkPlugins={[remarkGfm]}
        rehypePlugins={[rehypeRaw, [rehypeSanitize, sanitizeSchema], [rehypeKatex, {
          throwOnError: false,
          trust: false,
          strict: 'warn',
          output: 'html',
          errorColor: '#cc0000',
          macros: lessonMacros,
        }], rehypeKatexFallback]}
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
            // A missing alt falls back to the file name (never empty on content
            // images); loading is lazy to avoid layout/first-paint cost.
            const file = (url || '').split('/').pop()?.split('?')[0] || 'diagram'
            // eslint-disable-next-line @next/next/no-img-element
            return <img src={url} alt={alt || file} loading="lazy" className="max-w-full h-auto my-4 mx-auto" />
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
        .katex-content .math-fallback { font-family: var(--font-mono, monospace); font-size: 0.85em; background: var(--code-bg); border: 1px dashed currentColor; padding: 0.1em 0.35em; opacity: 0.85; overflow-wrap: anywhere; }
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
