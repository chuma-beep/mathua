'use client'

import Link from 'next/link'

export default function NotFound() {
  return (
    <main className="relative flex min-h-screen items-center justify-center overflow-hidden bg-mathua-bg px-6 py-16">
      <div
        aria-hidden="true"
        className="pointer-events-none absolute inset-0 opacity-[0.08]"
        style={{
          backgroundImage:
            "linear-gradient(var(--border) 1px, transparent 1px), linear-gradient(90deg, var(--border) 1px, transparent 1px)",
          backgroundSize: "48px 48px",
          maskImage: "radial-gradient(ellipse at center, black 40%, transparent 78%)",
        }}
      />

      <section className="relative z-10 w-full max-w-2xl">
        <div className="flex items-center justify-between border-b border-mathua-border pb-3 font-mono text-[10px] uppercase tracking-[0.18em] text-mathua-secondary">
          <span>Mathua / error</span>
          <span aria-hidden="true">§ 404</span>
        </div>

        <div className="mt-10 flex items-baseline gap-6">
          <h1 className="font-serif text-[7.5rem] leading-none font-bold text-mathua-blue tabular-nums">
            404
          </h1>
          <div className="hidden h-24 w-px bg-mathua-border sm:block" aria-hidden="true" />
          <div className="hidden font-mono text-xs leading-relaxed text-mathua-secondary sm:block">
            <div>lim<sub>x→/</sub> page(x)</div>
            <div className="mt-1 text-mathua-primary">= ∅</div>
          </div>
        </div>

        <h2 className="mt-8 font-serif text-2xl font-semibold text-mathua-primary">
          This route is undefined.
        </h2>
        <p className="mt-3 max-w-md font-mono text-sm leading-relaxed text-mathua-muted">
          The concept you requested isn&apos;t in the graph — it may have been moved,
          renamed, or never existed in this domain.
        </p>

        <div className="mt-10 flex flex-wrap items-center gap-3">
          <Link
            href="/"
            className="inline-flex items-center justify-center border border-mathua-blue bg-mathua-blue px-5 py-2.5 font-mono text-xs uppercase tracking-[0.18em] text-mathua-bg transition-colors hover:bg-mathua-blue-hover focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-mathua-primary"
          >
            ← Return home
          </Link>
          <button
            type="button"
            onClick={() => history.back()}
            className="inline-flex items-center justify-center border border-mathua-border bg-transparent px-5 py-2.5 font-mono text-xs uppercase tracking-[0.18em] text-mathua-primary transition-colors hover:bg-mathua-surface focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-mathua-primary"
          >
            Go back
          </button>
        </div>

        <div className="mt-16 flex items-center justify-between border-t border-mathua-border pt-3 font-mono text-[10px] uppercase tracking-[0.18em] text-mathua-secondary">
          <span>status: not_found</span>
          <span>code: 0x194</span>
        </div>
      </section>
    </main>
  )
}
