import type { DetailedHTMLProps, HTMLAttributes } from 'react'
import type { MathfieldElement } from 'mathlive'

// MathLive registers <math-field> when the module is imported, but its types
// declare the class only. React 19 moved the JSX namespace into the react
// module, so the intrinsic is declared here rather than in a global namespace.
declare module 'react' {
  namespace JSX {
    interface IntrinsicElements {
      'math-field': DetailedHTMLProps<HTMLAttributes<MathfieldElement>, MathfieldElement> & {
        /** MathLive attribute: LaTeX rendered while the field is empty. */
        placeholder?: string
        /** MathLive attribute: 'auto' | 'manual' | 'sandboxed'. */
        'math-virtual-keyboard-policy'?: string
        /**
         * `class`, not `className`. React forwards props to custom elements
         * verbatim instead of mapping className to class, so className would be
         * written as the attribute `classname` and the theme CSS would never match.
         */
        class?: string
      }
    }
  }
}
