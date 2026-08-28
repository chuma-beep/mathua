'use client'

import dynamic from 'next/dynamic'
import SectionHeader from '../../../components/SectionHeader'
import AsciiDivider from '../../../components/AsciiDivider'
import Loading from '../../../components/Loading'

const ArchitectureFlow = dynamic(() => import('../../../components/ArchitectureFlow'), {
  ssr: false,
  loading: () => (
    <div style={{
      height: 440, border: '0.5px solid var(--border)',
      display: 'flex', alignItems: 'center', justifyContent: 'center',
      color: 'var(--text-muted)', fontFamily: "'IBM Plex Mono', monospace", fontSize: '13px',
    }}>
      <Loading label="LOADING ARCHITECTURE DIAGRAM" />
    </div>
  ),
})

const DiagnosticFlow = dynamic(() => import('../../../components/DiagnosticFlow'), {
  ssr: false,
  loading: () => (
    <div style={{
      height: 420, border: '0.5px solid var(--border)',
      display: 'flex', alignItems: 'center', justifyContent: 'center',
      color: 'var(--text-muted)', fontFamily: "'IBM Plex Mono', monospace", fontSize: '13px',
    }}>
      <Loading label="LOADING DIAGNOSTIC DIAGRAM" />
    </div>
  ),
})

const DataModelFlow = dynamic(() => import('../../../components/SystemDesignFlow/DataModelFlow'), {
  ssr: false,
  loading: () => (
    <div style={{
      height: 320, border: '0.5px solid var(--border)',
      display: 'flex', alignItems: 'center', justifyContent: 'center',
      color: 'var(--text-muted)', fontFamily: "'IBM Plex Mono', monospace", fontSize: '13px',
    }}>
      <Loading label="LOADING DATA MODEL DIAGRAM" />
    </div>
  ),
})

const RequestFlow = dynamic(() => import('../../../components/SystemDesignFlow/RequestFlow'), {
  ssr: false,
  loading: () => (
    <div style={{
      height: 600, border: '0.5px solid var(--border)',
      display: 'flex', alignItems: 'center', justifyContent: 'center',
      color: 'var(--text-muted)', fontFamily: "'IBM Plex Mono', monospace", fontSize: '13px',
    }}>
      <Loading label="LOADING REQUEST FLOW DIAGRAM" />
    </div>
  ),
})

const headingFont = "'IBM Plex Serif', serif"
const bodyFont = "'IBM Plex Serif', serif"
const monoFont = "'IBM Plex Mono', monospace"

const h2Style: React.CSSProperties = {
  fontFamily: headingFont,
  fontWeight: 400,
  fontSize: '1.3rem',
  color: 'var(--text-primary)',
  borderBottom: '0.5px solid var(--border)',
  paddingBottom: '0.5rem',
  marginBottom: '1rem',
}

const bodyStyle: React.CSSProperties = {
  fontFamily: bodyFont,
  fontSize: '0.95rem',
  color: 'var(--text-secondary)',
  lineHeight: 1.85,
  marginBottom: '1rem',
}

const calloutStyle: React.CSSProperties = {
  borderLeft: '2px solid var(--accent-blue)',
  paddingLeft: '1.5rem',
  fontFamily: bodyFont,
  fontStyle: 'italic',
  fontSize: '0.9rem',
  color: 'var(--text-muted)',
  lineHeight: 1.7,
  marginTop: '1.5rem',
}

export default function ArchitecturePage() {
  return (
    <div className="max-w-container mx-auto px-4 sm:px-6 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0">
      <section className="pt-8 min-w-0 overflow-hidden">
        <SectionHeader label="Documentation" title="Architecture & System Design" />
        <p style={{ ...bodyStyle, textAlign: 'center', maxWidth: '640px', margin: '0 auto 2rem' }}>
          Mathua is a single Go binary with a web delivery mode. Click nodes in the diagrams below for details.
        </p>
        <div style={{ textAlign: 'center', marginBottom: '2rem' }}>
          <a
            href="https://github.com/chuma-beep/mathua/blob/main/docs/system-design.md"
            className="link-underline"
            style={{ fontFamily: monoFont, fontSize: '12px', color: 'var(--text-muted)' }}
          >
            View full system design doc on GitHub
          </a>
        </div>
      </section>

      {/* I. Five-Layer Architecture */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>I. Five-Layer Architecture</h2>
        <p style={bodyStyle}>
          The system is organized into five layers. The UI has a single implementation -- the web frontend
          (React/Next.js) -- calling into the Go engine through the REST API.
          The API exposes 25+ REST endpoints through net/http. The Core Engine handles DAG loading,
          SM-2 scheduling, problem generation, mastery tracking, scoring, and the CAT diagnostic.
          The Grading layer dispatches to 6+ grader strategies with a SymPy subprocess for symbolic
          math. Storage is abstracted behind a Repository interface.
        </p>

        <ArchitectureFlow />
      </section>

      <AsciiDivider pattern="wave" />

      {/* II. Data Model */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>II. Data Model</h2>
        <p style={bodyStyle}>
          Five core tables store all application state. The `concept_progress` table embeds SM-2
          spaced repetition fields (repetitions, interval, efactor) alongside mastery state and
          weakness scores. All DDL uses CREATE TABLE IF NOT EXISTS for idempotent bootstrapping.
          Incremental migrations add columns with ALTER TABLE guarded by error-checking.
        </p>

        <DataModelFlow />

        <p style={bodyStyle}>
          The SM-2 upsert uses SQLite&apos;s ON CONFLICT ... DO UPDATE to atomically save all progress
          fields in one statement. XP tracking is date-aware: xp_today resets when xp_date differs
          from the current date, preserving xp_total as a lifetime accumulator.
        </p>

        <div style={calloutStyle}>
          The PostgresStore exists as a stub with all methods returning &quot;not implemented.&quot;
          The same Repository interface works for both databases -- the schema is identical.
        </div>
      </section>

      <AsciiDivider pattern="wave" />

      {/* III. Request Flow */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>III. End-to-End Request Flow</h2>
        <p style={bodyStyle}>
          When a student submits an answer, the request passes through 12 distinct stages before
          the next question is served. The entire pipeline runs synchronously in a single Go
          goroutine, completing in under 100ms for numeric grading and under 500ms for SymPy-based
          grading (including subprocess round-trip).
        </p>

        <RequestFlow />

        <div style={{ marginTop: '2rem' }}>
          {[
            'Client sends POST /api/answer with {session_id, answer, elapsed}.',
            'Auth middleware validates JWT and injects student ID into request context.',
            'Engine.SubmitAnswer routes to the correct grader by grading_type.',
            'The Mastery Machine evaluates whether a state transition is earned.',
            'SM-2 Compute recalculates repetition count, interval, and easiness factor.',
            'repo.UpsertProgress atomically saves all fields via ON CONFLICT DO UPDATE.',
            'If weakness > 0.3, it propagates to dependent concepts at w * 0.3.',
            'repo.RecordAttempt stores the raw answer for analytics.',
            'XP is computed (base * timeMult * streakMult) and added to the student.',
            'Engine.NextQuestion asks the scheduler for the next concept and generates a problem.',
          ].map((step, i) => (
            <div key={step} style={{ ...bodyStyle, marginBottom: '0.5rem' }}>
              <span style={{ color: 'var(--accent-blue)', fontFamily: monoFont, fontSize: '13px', fontWeight: 500 }}>
                {i + 1}.
              </span>
              {' '}{step}
            </div>
          ))}
        </div>
      </section>

      <AsciiDivider pattern="wave" />

      {/* IV. Grading System */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>IV. Grading System</h2>
        <p style={bodyStyle}>
          The grading system uses a strategy pattern dispatched by grading_type. Numeric grading
          (int, float, fraction, mixed, scientific notation) is pure Go with big.Rat for exact
          rational arithmetic and 1e-9 float tolerance. Multiple choice is case-insensitive with
          single-letter matching.           Comparison handles operators: {'>'} {'<'} {'='} {'>='} {'<='} {'!='}. Ordering and
          tuple graders do positional exact matching.
        </p>
        <p style={bodyStyle}>
          For symbolic math -- polynomial and expression grading -- Mathua spawns a long-lived
          Python subprocess running SymPy. The Go client sends a JSON pair over stdin; SymPy parses
          both into expression trees and checks equivalence via simplify(expected - answer) == 0.
          If Python/SymPy are not installed, grading falls back gracefully to a pure-Go symbolic
          string normalizer. The subprocess has a 10-second timeout and 500-character input limit.
        </p>
        <p style={bodyStyle}>
          The complex grader preprocesses polar form, handles plus-minus notation, and delegates to
          SymPy. In total there are 8 grading strategies: numeric, multiple choice, comparison,
          ordering, tuple, complex, symbolic (fallback), and SymPy (polynomial/expression).
        </p>

        <div style={calloutStyle}>
          The D2 diagram for the grading pipeline is available in the <a href="/docs" className="link-underline" style={{ color: 'var(--accent-blue)' }}>system design documentation</a>.
        </div>
      </section>

      <AsciiDivider pattern="wave" />

      {/* V. CAT Diagnostic */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>V. Computerised Adaptive Testing</h2>
        <p style={bodyStyle}>
          The diagnostic engine locates a student&apos;s knowledge frontier using binary search on
          the topologically sorted concept graph. This reduces the assessment from 284 questions
          (one per concept) to approximately 20-35.
        </p>

        <DiagnosticFlow />

        <div style={{ marginTop: '2rem' }}>
          {[
            'The concept graph is sorted topologically. The diagnostic starts at the midpoint.',
            'Correct answers within the time limit move the probe forward toward harder concepts.',
            'Incorrect or slow answers move backward toward foundational material.',
            'After 3 consecutive correct answers in a region, the frontier is considered located.',
            'The diagnostic records a mastery estimate for every concept passed through.',
          ].map((step, i) => (
            <div key={step} style={{ ...bodyStyle, marginBottom: '0.5rem' }}>
              <span style={{ color: 'var(--accent-blue)', fontFamily: monoFont, fontSize: '13px', fontWeight: 500 }}>
                {i + 1}.
              </span>
              {' '}{step}
            </div>
          ))}
        </div>
        <div style={calloutStyle}>
          The diagnostic can be retaken at any time. Retaking does not delete progress: it creates
          a new estimate that is merged with existing data, always preferring the more optimistic
          estimate.
        </div>
      </section>

      <AsciiDivider pattern="wave" />

      {/* VI. Key Trade-offs */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>VI. Key Design Decisions</h2>

        <div style={{ ...bodyStyle, marginBottom: '1.5rem' }}>
          <strong style={{ color: 'var(--text-primary)', fontFamily: monoFont, fontSize: '13px' }}>
            Go over Python or Node.js?
          </strong>
          <p style={{ margin: '4px 0 0' }}>
            Single ~15MB binary that is both API server and static file server. No runtime dependencies.
            ~5ms startup vs ~500ms for Python. Goroutines for concurrency. SymPy bridge via os/exec.
          </p>
        </div>

        <div style={{ ...bodyStyle, marginBottom: '1.5rem' }}>
          <strong style={{ color: 'var(--text-primary)', fontFamily: monoFont, fontSize: '13px' }}>
            Generated questions over static bank?
          </strong>
          <p style={{ margin: '4px 0 0' }}>
            Every problem is procedurally generated by a parameterized Go function. Infinite variety,
            no memorization, adaptive difficulty. Trade-off: generation latency (ns for numeric, ~200ms
            for SymPy). Worth it.
          </p>
        </div>

        <div style={{ ...bodyStyle, marginBottom: '1.5rem' }}>
          <strong style={{ color: 'var(--text-primary)', fontFamily: monoFont, fontSize: '13px' }}>
            SQLite + PostgreSQL instead of one database?
          </strong>
          <p style={{ margin: '4px 0 0' }}>
            SQLite for zero-config local development, PostgreSQL for production web. Same schema,
            same Repository interface. Transparent via DATABASE_URL. The PostgresStore is currently
            a stub awaiting implementation.
          </p>
        </div>

        <div style={{ ...bodyStyle, marginBottom: '1.5rem' }}>
          <strong style={{ color: 'var(--text-primary)', fontFamily: monoFont, fontSize: '13px' }}>
            Raw SQL over an ORM?
          </strong>
          <p style={{ margin: '4px 0 0' }}>
            Five tables, straightforward relationships. Raw SQL gives full control over the SM-2
            upsert query and leaderboard computation. Transparent debugging. Easy to port between
            SQLite and PostgreSQL.
          </p>
        </div>

        <div style={calloutStyle}>
          For the full system design document covering all 11 sections in detail (including SM-2
          algorithm internals, scoring formulas, level system, and complete API reference),
          see the <a href="https://github.com/chuma-beep/mathua/blob/main/docs/system-design.md" className="link-underline" style={{ color: 'var(--accent-blue)' }}>system design markdown document</a>.
        </div>
      </section>
    </div>
  )
}
