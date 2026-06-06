'use client'

import dynamic from 'next/dynamic'
import SectionHeader from '../../../components/SectionHeader'
import AsciiDivider from '../../../components/AsciiDivider'

const ArchitectureFlow = dynamic(() => import('../../../components/SystemDesignFlow/ArchitectureFlow'), {
  ssr: false,
  loading: () => (
    <div style={{
      height: 520, border: '0.5px solid var(--border)',
      display: 'flex', alignItems: 'center', justifyContent: 'center',
      color: 'var(--text-muted)', fontFamily: "'IBM Plex Mono', monospace", fontSize: '13px',
    }}>
      Loading architecture diagram&hellip;
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
      Loading data model diagram&hellip;
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
      Loading request flow diagram&hellip;
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

export default function SystemDesignPage() {
  return (
    <div className="max-w-container mx-auto px-6 max-sm:px-4">
      <section className="pt-8">
        <SectionHeader label="Documentation" title="System Design" />
        <p style={{ ...bodyStyle, textAlign: 'center', maxWidth: '640px', margin: '0 auto 2rem' }}>
          Mathua is a single Go binary with two delivery modes. The engine core is identical --
          only the presentation layer differs. Click nodes in the diagrams below for details.
        </p>
        <div style={{ textAlign: 'center', marginBottom: '2rem' }}>
          <a
            href="https://github.com/chuma-beep/mathua/blob/main/docs/system-design.md"
            className="link-underline"
            style={{ fontFamily: monoFont, fontSize: '12px', color: 'var(--text-muted)' }}
          >
            View source on GitHub
          </a>
        </div>
      </section>

      {/* Section I: Architecture */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>I. Architecture: Five layers, one binary</h2>
        <p style={bodyStyle}>
          The system is organized into five layers stacked vertically. The UI layer has two implementations --
          a web frontend (React/Next.js) and a desktop TUI (Bubble Tea) -- both calling into the same Go engine
          through different entry points. The API layer exposes 25+ REST endpoints through net/http.
          The Core Engine handles DAG loading, SM-2 scheduling, problem generation, mastery tracking,
          scoring, and the CAT diagnostic. The Grading layer dispatches to 6+ grader strategies with a
          SymPy subprocess for symbolic math. Storage is abstracted behind a Repository interface.
        </p>

        <ArchitectureFlow />

        <div style={{ marginTop: '2rem' }}>
          {[
            'The DAG Loader validates concept prerequisites and produces a topological sort via Kahn\'s algorithm.',
            'The SM-2 Scheduler selects the next concept using priority scoring (recency, mastery, weakness) with a 70/30 new/review balance.',
            'All problems are generated on demand by parameterized Go functions -- no static question bank.',
            'The Mastery Machine tracks state transitions: UNSEEN -> LEARNING -> PRACTICING -> MASTERED -> DECAYING.',
            'Weakness propagates from concepts with score > 0.3 to their dependents at w * 0.3.',
          ].map((item, i) => (
            <div key={item} style={{ ...bodyStyle, marginBottom: '0.5rem' }}>
              <span style={{ color: 'var(--accent-blue)', fontFamily: monoFont, fontSize: '13px', fontWeight: 500 }}>
                {i + 1}.
              </span>
              {' '}{item}
            </div>
          ))}
        </div>
      </section>

      <AsciiDivider pattern="wave" />

      {/* Section II: Data Model */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>II. Data Model: Five tables, SM-2 integrated</h2>
        <p style={bodyStyle}>
          Five core tables store all application state. The `concept_progress` table is the most complex --
          it embeds SM-2 spaced repetition fields (repetitions, interval, efactor) alongside mastery state
          and weakness scores. All DDL uses CREATE TABLE IF NOT EXISTS for idempotent bootstrapping.
          Incremental migrations add columns with ALTER TABLE guarded by error-checking.
        </p>

        <DataModelFlow />

        <p style={bodyStyle}>
          The SM-2 upsert query uses SQLite's ON CONFLICT ... DO UPDATE pattern to atomically save
          all progress fields in a single statement. XP tracking is date-aware: xp_today resets when
          xp_date differs from the current date, preserving xp_total as a lifetime accumulator.
        </p>

        <div style={calloutStyle}>
          The PostgresStore exists as a stub with all methods returning "not implemented."
          The same Repository interface works for both databases -- the schema is identical.
        </div>
      </section>

      <AsciiDivider pattern="wave" />

      {/* Section III: Request Flow */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>III. Request Flow: Answer submission end-to-end</h2>
        <p style={bodyStyle}>
          When a student submits an answer, the request passes through 12 distinct stages before
          the next question is served. The entire pipeline runs synchronously within a single
          Go goroutine, completing in under 100ms for numeric grading and under 500ms for
          SymPy-based grading (including subprocess round-trip).
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

        <div style={calloutStyle}>
          The diagnostic engine locates a student's knowledge frontier using binary search on
          the topologically sorted concept graph. This reduces the assessment from 284 questions
          to approximately 20-35. Retaking the diagnostic uses optimistic merge: new estimates
          are merged with existing data, always preferring the more optimistic estimate.
        </div>
      </section>
    </div>
  )
}
