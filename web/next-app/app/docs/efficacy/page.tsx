import type { Metadata } from 'next'
import Link from 'next/link'
import Image from 'next/image'

export const metadata: Metadata = {
  title: 'Efficacy — Mathua',
  description: 'How Mathua measures learning: first/second-pass instrumentation, retention, and the evidence figure.',
}

const headingFont = "var(--font-space-grotesk), 'Space Grotesk', serif"
const bodyFont = "var(--font-space-grotesk), 'Space Grotesk', serif"
const monoFont = "var(--font-jetbrains-mono), 'JetBrains Mono', monospace"

const h2Style: React.CSSProperties = {
  fontFamily: headingFont,
  fontWeight: 400,
  fontSize: 'clamp(1.15rem, 1rem + 0.8vw, 1.3rem)',
  color: 'var(--text-primary)',
  borderBottom: '0.5px solid var(--border)',
  paddingBottom: '0.5rem',
  marginBottom: '1rem',
  marginTop: '2.5rem',
}

const bodyStyle: React.CSSProperties = {
  fontFamily: bodyFont,
  fontSize: 'clamp(0.9rem, 0.85rem + 0.5vw, 1rem)',
  color: 'var(--text-secondary)',
  lineHeight: 1.85,
  marginBottom: '1rem',
}

const codeBlockStyle: React.CSSProperties = {
  background: 'var(--surface)',
  border: 'none',
  padding: '16px 20px',
  fontFamily: monoFont,
  fontSize: '13px',
  color: 'var(--text-secondary)',
  whiteSpace: 'pre',
  overflowX: 'auto',
  lineHeight: 1.6,
  margin: '0 0 1rem 0',
  width: '100%',
  maxWidth: '100%',
  minWidth: 0,
}

const captionStyle: React.CSSProperties = {
  fontFamily: monoFont,
  fontSize: '12px',
  color: 'var(--text-muted)',
  lineHeight: 1.7,
  marginTop: '0.75rem',
}

function Ol({ items }: { items: React.ReactNode[] }) {
  return (
    <ol style={{ ...bodyStyle, paddingLeft: '1.5rem', marginTop: 0 }}>
      {items.map((item, i) => (
        <li key={i} style={{ marginBottom: '0.5rem' }}>{item}</li>
      ))}
    </ol>
  )
}

function Ul({ items }: { items: React.ReactNode[] }) {
  return (
    <ul style={{ ...bodyStyle, paddingLeft: '1.5rem', marginTop: 0 }}>
      {items.map((item, i) => (
        <li key={i} style={{ marginBottom: '0.5rem' }}>{item}</li>
      ))}
    </ul>
  )
}

export default function EfficacyPage() {
  return (
    <div className="max-w-[720px] mx-auto px-4 sm:px-6 py-12 sm:py-20 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-20 overflow-x-hidden min-w-0">
      <Link href="/docs" style={{ fontFamily: monoFont, fontSize: '13px', color: 'var(--text-muted)' }}>
        ← Docs
      </Link>
      <h1 style={{
        fontFamily: headingFont,
        fontWeight: 400,
        fontSize: '1.9rem',
        color: 'var(--text-primary)',
        marginBottom: '0.5rem',
        marginTop: '1rem',
        borderBottom: '0.5px solid var(--border)',
        paddingBottom: '0.5rem',
      }}>
        Efficacy
      </h1>
      <p style={{ ...bodyStyle, color: 'var(--text-muted)' }}>
        How Mathua measures learning — and the honest limits of what the numbers can claim.
      </p>

      <h2 style={{ ...h2Style, marginTop: '2rem' }}>What is recorded</h2>
      <p style={bodyStyle}>
        Every answer writes one <span style={{ fontFamily: monoFont, fontSize: '0.9em' }}>attempts</span> row:
        student, concept, correctness, elapsed seconds, timestamp. Nothing is sampled or
        aggregated at write time — every metric below is computed from the raw log.
      </p>

      <h2 style={h2Style}>Point-in-time figure</h2>
      <div style={{
        border: '0.5px solid var(--border)',
        padding: '1rem',
        background: 'var(--surface)',
      }}>
        <Image
          src="/figures/efficacy-overview.svg"
          alt="First-pass and second-pass rates from the dev attempt log, with Math Academy parity markers"
          width={1198}
          height={778}
          style={{ maxWidth: '100%', height: 'auto', display: 'block' }}
        />
      </div>
      <p style={captionStyle}>
        Generated from the dev attempt log by <span style={{ color: 'var(--text-secondary)' }}>python3 figures/efficacy.py</span> (style
        adapter: <span style={{ color: 'var(--text-secondary)' }}>figures/mathua_style.py</span>). The MA 93%/98% lines are
        external reference markers, not Mathua data. The weekly-trend panel is deliberately omitted: with
        70 + 1 attempts across two week-buckets a trend line would mislead. Regenerate after real usage accumulates:
      </p>
      <pre style={codeBlockStyle}>{`python3 figures/fetch_fonts.py   # once per machine (gitignored display fonts)
python3 figures/efficacy.py      # dump mathua.db -> figures/efficacy_data.json + render
# then copy figures/output/efficacy-overview.svg to web/next-app/public/figures/`}</pre>

      <h2 style={h2Style}>Longitudinal metrics</h2>
      <p style={bodyStyle}>
        <span style={{ fontFamily: monoFont, fontSize: '0.9em' }}>GET /api/efficacy/trend</span> returns weekly
        Monday–Sunday buckets with per-week pass rates, plus a coarse retention signal (students active in
        two or more distinct weeks ÷ all students who ever attempted) and the first-pass trend (last week
        minus first week). The Profile page renders the last 12 weeks as a first-pass sparkline plus the
        retention figure.
      </p>

      <h2 style={h2Style}>How to read it</h2>
      <Ol items={[
        <span key="1"><strong style={{ color: 'var(--text-primary)', fontWeight: 400 }}>Compare within a student, not across.</strong> A learner&apos;s own weekly first-pass rate is a more honest signal than the product-wide number, which mixes learners at very different frontiers.</span>,
        <span key="2"><strong style={{ color: 'var(--text-primary)', fontWeight: 400 }}>Watch active_students before the rate.</strong> A week with three attempts can show 0% or 100% by accident.</span>,
        <span key="3"><strong style={{ color: 'var(--text-primary)', fontWeight: 400 }}>Second-pass is the recovery metric.</strong> A gap between first- and second-pass rates means material is learnable with feedback but not yet automatic.</span>,
      ]} />

      <h2 style={h2Style}>Limitations</h2>
      <Ul items={[
        'No control group and no randomisation: these are usage metrics, not an efficacy claim. Mathua does not assert the 4× / 93% / 98% figures from Math Academy\u2019s marketing.',
        'Retention is binary (returned / did not) over the whole window; it does not model decay or churn timing.',
        'Buckets are UTC weeks; a learner near a timezone boundary can appear in two weeks for one study session.',
        'The attempt log is the source of truth. If it is cleared, all metrics reset.',
      ]} />

      <h2 style={h2Style}>Reproducing an analysis</h2>
      <pre style={codeBlockStyle}>{`# product-wide point-in-time
curl -s localhost:8080/api/efficacy/all | jq
# weekly trend + retention
curl -s localhost:8080/api/efficacy/trend | jq`}</pre>
      <p style={{ ...bodyStyle, fontStyle: 'italic', color: 'var(--text-muted)' }}>
        Both read only the local database; no external service or API key is involved. Full definitions live in <Link href="https://github.com/chuma-beep/mathua/blob/main/docs/efficacy.md" className="link-underline" style={{ color: 'var(--accent-blue)' }}>docs/efficacy.md</Link>.
      </p>
    </div>
  )
}
