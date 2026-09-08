'use client'

// Root error boundary: replaces the whole layout when the root itself
// throws, so it must render its own <html>/<body> and cannot rely on the
// root layout's fonts, theme script, or Toaster.
export default function GlobalError({
  error,
  reset,
}: {
  error: Error & { digest?: string }
  reset: () => void
}) {
  return (
    <html lang="en" className="dark">
      <body style={{ background: '#0b0f1a', color: '#e8e2d5', fontFamily: 'monospace' }}>
        <main style={{ display: 'flex', minHeight: '100vh', alignItems: 'center', justifyContent: 'center', padding: '24px' }}>
          <section style={{ width: '100%', maxWidth: '640px' }}>
            <div style={{ fontSize: '10px', textTransform: 'uppercase', letterSpacing: '0.18em', opacity: 0.6 }}>
              Mathua / error § 500
            </div>
            <h1 style={{ fontSize: '2rem', marginTop: '24px' }}>Something broke.</h1>
            <p role="alert" style={{ marginTop: '12px', fontSize: '14px', opacity: 0.7 }}>
              {error?.message || 'An unexpected error interrupted rendering.'}
              {error?.digest ? ` (ref ${error.digest})` : ''}
            </p>
            <div style={{ marginTop: '32px', display: 'flex', gap: '12px' }}>
              <button
                type="button"
                onClick={() => reset()}
                style={{ border: '1px solid #c8a96e', color: '#c8a96e', background: 'transparent', padding: '10px 20px', fontSize: '12px', cursor: 'pointer' }}
              >
                ↻ Try again
              </button>
              <a href="/" style={{ border: '1px solid #444', color: '#e8e2d5', padding: '10px 20px', fontSize: '12px', textDecoration: 'none' }}>
                ← Return home
              </a>
            </div>
          </section>
        </main>
      </body>
    </html>
  )
}
