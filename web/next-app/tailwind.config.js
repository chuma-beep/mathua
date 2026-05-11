/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './app/**/*.{js,ts,jsx,tsx,mdx}',
    './components/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        mathua: {
          bg: 'var(--bg)',
          surface: 'var(--surface)',
          'surface-elevated': 'var(--surface-elevated)',
          'surface-highlight': 'var(--surface-highlight)',
          primary: 'var(--text-primary)',
          secondary: 'var(--text-secondary)',
          muted: 'var(--text-muted)',
          gold: 'var(--accent-gold)',
          'gold-hover': 'var(--accent-gold-hover)',
          blue: 'var(--accent-blue)',
          'blue-hover': 'var(--accent-blue-hover)',
          green: 'var(--accent-green)',
          red: 'var(--accent-red)',
          border: 'var(--border)',
          'border-strong': 'var(--border-strong)',
          code: 'var(--code-bg)',
          division: 'var(--division)',
        },
      },
      fontFamily: {
        serif: ['Playfair Display', 'Georgia', 'serif'],
        sans: ['Inter', '-apple-system', 'system-ui', 'sans-serif'],
        mono: ['JetBrains Mono', 'monospace'],
      },
      maxWidth: {
        container: '1100px',
      },
      borderRadius: {
        card: '8px',
        pill: '9999px',
      },
      keyframes: {
        'ascii-reveal': {
          '0%': { clipPath: 'inset(0 100% 0 0)' },
          '100%': { clipPath: 'inset(0 0 0 0)' },
        },
        'ascii-pulse': {
          '0%, 100%': { opacity: '0.4' },
          '50%': { opacity: '1' },
        },
        'progress-fill': {
          '0%': { width: '0%' },
          '100%': { width: 'var(--progress-width)' },
        },
      },
      animation: {
        'ascii-reveal': 'ascii-reveal 1.5s steps(30) forwards',
        'ascii-pulse': 'ascii-pulse 2s ease-in-out infinite',
        'progress-fill': 'progress-fill 0.8s ease-out forwards',
      },
      backgroundImage: {
        'noise': "url(\"data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.65' numOctaves='3' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)' opacity='0.03'/%3E%3C/svg%3E\")",
      },
    },
  },
  plugins: [],
}
