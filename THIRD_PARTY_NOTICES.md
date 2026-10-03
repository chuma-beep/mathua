# Third-Party Notices

Mathua is MIT licensed (see `LICENSE`). It also depends on, ships, or redistributes
content from the projects below. This file records those obligations in one place so
they are not scattered across the README and forgotten when a dependency changes.

## Software dependencies

### MathLive — `mathlive@0.111.0`

Used for the answer editor in `web/next-app/components/math/`. Provides the
`<math-field>` custom element: mathematical input, structure navigation, and the
virtual keyboard. Licensed under the MIT License.

MathLive is **used, not vendored**. Nothing in it is copied into this repository; the
integration goes through its public API, configuration, custom keyboard layouts, CSS
variables and `::part()`. MathLive brings its own webfonts (the KaTeX family), which
`web/next-app/scripts/copy-mathlive-fonts.mjs` copies from `node_modules` into
`public/fonts` at build time. Those files are build output and are not committed.

```
Copyright (c) 2017 - present Arno Gourdol. All rights reserved.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

Upstream: <https://github.com/arnog/mathlive> · Documentation: <https://mathlive.io/mathfield/>

### Compute Engine — `@cortex-js/compute-engine@0.147.0`

Used for one thing: `web/next-app/components/math/SelfCheck.tsx`, which tells a learner
what their own expression simplifies to. Licensed under the MIT License.

Compute Engine is **used, not vendored**, and is a browser dependency only — it is
dynamically imported so it never enters any route's initial JS, and it is skipped on
metered connections. It is *not* used for grading: the Go engine and SymPy remain the
only graders (ADR-005), and the self-check cannot reach a submitted answer.

Note for anyone extending this: importing the library also switches on MathLive's
MathJSON paste support, because Compute Engine registers itself under
`globalThis[Symbol.for("io.cortexjs.compute-engine")]` as a side effect of being
imported. That is an implicit coupling and is asserted in
`web/next-app/test/equivalence.test.ts`.

```
MIT License

Copyright (c) 2019 CortexJS

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

Upstream: <https://cortexjs.io/compute-engine/>

### Other npm dependencies

KaTeX (MIT) renders question text, answers and lesson content; three.js and
`@react-three/fiber` (MIT) render the landing concept graph; dagre, lucide-react,
remark-math and rehype-katex (MIT) support layout, icons and Markdown. Their licence
texts ship inside each package under `node_modules/`, and are reproduced here on
upgrade by the package manager.

## Go dependencies

The Go module's dependency licences are recorded in `go.sum` and in each module's own
repository. SymPy (MIT) is the grading runtime; it is pinned in
`grading/requirements.txt` and installed into the image.

## Content

### Algebrica — CC BY-NC 4.0

Lesson content and diagrams under `data/lessons/` and `figures/` are sourced from
[Algebrica](https://algebrica.org) by Antonio Lupetti, used under
[CC BY-NC 4.0](https://creativecommons.org/licenses/by-nc/4.0/).

**This is a non-commercial licence.** Mathua's own MIT licence does not extend to this
content: the lesson material and diagrams remain CC BY-NC 4.0 and carry that
restriction with them.

### DiceBear — per-style licences

Profile avatars are served by [DiceBear](https://www.dicebear.com). Individual avatar
styles carry their own licences; see their
[license overview](https://www.dicebear.com/licenses/).

### OpenStax, MIT OpenCourseWare, Art of Problem Solving

Referenced for content design and curriculum structure. No code or content is taken
from any of them. See the Acknowledgements section of `README.md`.