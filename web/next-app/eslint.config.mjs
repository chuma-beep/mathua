import { defineConfig, globalIgnores } from 'eslint/config'
import nextVitals from 'eslint-config-next/core-web-vitals'

// Next 16 removed `next lint`, so the CLI runs the same Core Web Vitals rules
// directly. The lint surface stays at app, components and lib, which is what
// `next lint` covered before. The TypeScript rule set is deliberately not added
// here; that is a separate change with its own findings.
const eslintConfig = defineConfig([
  ...nextVitals,
  {
    rules: {
      // Next 16 ships the React Compiler lint rules as errors. Adopting them
      // means a separate refactor across the app, so keep the pre-upgrade lint
      // surface until that change.
      'react-hooks/set-state-in-effect': 'off',
      'react-hooks/refs': 'off',
      'react-hooks/preserve-manual-memoization': 'off',
      'react-hooks/immutability': 'off',
      'react-hooks/purity': 'off',
    },
  },
  globalIgnores(['.next/**', 'out/**', 'build/**', 'next-env.d.ts']),
])

export default eslintConfig
