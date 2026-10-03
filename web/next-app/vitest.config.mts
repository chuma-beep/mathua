import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'
import path from 'path'

export default defineConfig({
  plugins: [react()],
  test: {
    environment: 'jsdom',
    environmentOptions: {
      jsdom: {
        url: 'http://localhost/',
      },
    },
    setupFiles: ['./test/setup.ts'],
    alias: [
      { find: /^mathlive$/, replacement: path.resolve(__dirname, 'test/stubs/mathlive.ts') },
    ],
    server: {
      deps: {
        // Resolve `mathlive` through the alias above. Without this, vitest
        // externalises it as a node_modules dependency before resolution runs and
        // the alias never applies.
        inline: ['mathlive'],
      },
    },
    include: ['test/**/*.test.{ts,tsx}'],
    css: false,
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, '.'),
      // MathLive cannot mount in jsdom (real text measurement, and its `node`
      // export condition resolves to the SSR build that never registers the
      // element). The stub keeps the component suite about Mathua's behaviour;
      // the real library is exercised by the Playwright math-input suites.
      // Exact match only. A bare-string alias is a prefix match in Vite, so
      // `mathlive/mathlive-static.css` would rewrite to `<stub>/mathlive-static.css`
      // and the dynamic import would reject instead of resolving.
      mathlive: { find: /^mathlive$/, replacement: path.resolve(__dirname, 'test/stubs/mathlive.ts') },
    },
  },
})
