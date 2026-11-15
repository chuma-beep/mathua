'use client'

import katex from 'katex'
import 'katex/dist/katex.min.css'

// One piece of mathematics from a reply, rendered the way the lesson surfaces
// render it.
//
// This uses `dangerouslySetInnerHTML` on purpose, and the reason is narrow:
// KaTeX does not take HTML, it takes TeX and produces its own output tree. With
// `trust: false` the commands that would emit an attribute or a URL, \href,
// \url, \htmlClass, \includegraphics, are refused and rendered as literal text,
// so no string the model wrote can reach an anchor or an event handler. That is
// the same configuration rehype-katex uses for lesson content, which is why this
// does not import the Markdown pipeline: full Markdown would turn
// `[x](https://elsewhere)` into a link, and links in a reply are supposed to
// come only from the server's [[id]] table.
//
// `throwOnError: false` turns a malformed expression into a red inline error for
// that one span rather than an exception that takes the panel down.
export default function AssistantMath({ tex, display }: { tex: string; display: boolean }) {
	const html = katex.renderToString(tex, {
		displayMode: display,
		throwOnError: false,
		strict: 'ignore',
		trust: false,
	})

	return (
		<span
			className={display ? 'my-1 block overflow-x-auto' : 'inline-block align-middle'}
			// aislop-ignore-line security/dangerously-set-innerhtml -- the string is KaTeX's output tree built from TeX with trust:false, never model HTML
			dangerouslySetInnerHTML={{ __html: html }}
		/>
	)
}
