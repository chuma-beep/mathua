'use client'

import SectionHeader from '../../../components/SectionHeader'
import AsciiDivider from '../../../components/AsciiDivider'

const headingFont = "'Lora', serif"
const bodyFont = "'Lora', serif"
const monoFont = "'JetBrains Mono', 'Fira Code', monospace"

const h2Style: React.CSSProperties = {
  fontFamily: headingFont,
  fontWeight: 400,
  fontSize: '1.3rem',
  color: 'var(--text-primary)',
  borderBottom: '0.5px solid var(--border)',
  paddingBottom: '0.5rem',
  marginBottom: '1rem',
}

const h3Style: React.CSSProperties = {
  fontFamily: headingFont,
  fontWeight: 400,
  fontSize: '1rem',
  color: 'var(--text-primary)',
  marginTop: '1.5rem',
  marginBottom: '0.5rem',
}

const bodyStyle: React.CSSProperties = {
  fontFamily: bodyFont,
  fontSize: '0.95rem',
  color: 'var(--text-secondary)',
  lineHeight: 1.85,
  marginBottom: '1rem',
}

const codeBlockStyle: React.CSSProperties = {
  background: 'var(--surface)',
  border: 'none',
  padding: '20px 24px',
  fontFamily: monoFont,
  fontSize: '13px',
  color: 'var(--text-secondary)',
  whiteSpace: 'pre',
  overflowX: 'auto',
  lineHeight: 1.6,
  margin: '1rem 0',
}

const tableHeaderStyle: React.CSSProperties = {
  fontFamily: monoFont,
  fontWeight: 400,
  fontSize: '0.85rem',
  textTransform: 'uppercase',
  letterSpacing: '0.08em',
  color: 'var(--accent-gold)',
  padding: '10px 14px 10px 0',
  borderBottom: '1px solid var(--accent-gold)',
  textAlign: 'left',
  background: 'transparent',
}

const tableCellStyle: React.CSSProperties = {
  fontFamily: bodyFont,
  fontSize: '0.9rem',
  color: 'var(--text-secondary)',
  padding: '10px 14px 10px 0',
  borderBottom: '0.5px solid var(--border)',
  background: 'transparent',
}

export default function ContributingPage() {
  return (
    <div className="max-w-container mx-auto px-6 max-sm:px-4">
      <section className="pt-8">
        <SectionHeader label="Community" title="Contributing to Mathua" />
        <p style={{ ...bodyStyle, textAlign: 'center', maxWidth: '640px', margin: '0 auto 2rem' }}>
          Mathua is community-built. Every concept, every generator, every line of the
          concept graph was added by someone who wanted to help others learn math better.
          The most useful thing you can contribute is a new concept — and it takes exactly
          three pieces: a JSON entry, a Go generator, and a fuzz test.
        </p>
      </section>

      <AsciiDivider pattern="wave" />

      {/* Getting started */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>Getting started</h2>
        <p style={bodyStyle}>
          You'll need Go 1.21+. Clone the repo and make sure the TUI runs:
        </p>
        <pre style={codeBlockStyle}>
{`git clone https://github.com/chuma-beep/mathua.git
cd mathua
go build ./cmd/mathua
./mathua`}
        </pre>
        <p style={bodyStyle}>
          The concept graph lives in{' '}
          <code style={{ fontFamily: monoFont, fontSize: '0.9em', color: 'var(--accent-gold)' }}>data/concepts.json</code>.
          Generators live in{' '}
          <code style={{ fontFamily: monoFont, fontSize: '0.9em', color: 'var(--accent-gold)' }}>internal/generator/</code>.
          Here are the kinds of contributions that move the needle:
        </p>
        {[
          'New concepts — the single most impactful thing you can add.',
          'New generators — make existing or new concepts produce better problems.',
          'Bug fixes in the scheduling engine or graders.',
          'Documentation and diagram improvements.',
        ].map((item, i) => (
          <div key={i} style={{ ...bodyStyle, marginBottom: '0.4rem' }}>
            <span style={{ color: 'var(--border-strong)', marginRight: '0.25rem', fontFamily: monoFont }}>·</span>
            {item}
          </div>
        ))}
      </section>

      <AsciiDivider pattern="dash" />

      {/* Step 1 */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>Step 1: Define a new concept</h2>
        <p style={bodyStyle}>
          Every concept lives in{' '}
          <code style={{ fontFamily: monoFont, fontSize: '0.9em', color: 'var(--accent-gold)' }}>data/concepts.json</code>.
          Add a JSON object with an ID, label, domain, prerequisite list, and mastery thresholds.
        </p>
        <pre style={codeBlockStyle}>
{`{
  "id":                "arith.mult.tables",
  "label":             "Multiplication tables 1-12",
  "domain":            "arithmetic",
  "subdomain":         "arithmetic.multiplication",
  "grading_type":      "numeric",
  "prerequisites":     ["arith.add.multi"],
  "mastery_threshold": {
    "streak":           7,
    "avg_time_seconds": 6.0
  }
}`}
        </pre>

        <div style={{ overflowX: 'auto', margin: '1.5rem 0' }}>
          <table style={{
            width: '100%',
            borderCollapse: 'collapse',
            border: 'none',
            background: 'transparent',
          }}>
            <thead>
              <tr style={{ background: 'transparent' }}>
                <th style={tableHeaderStyle}>Field</th>
                <th style={tableHeaderStyle}>Type</th>
                <th style={tableHeaderStyle}>Required</th>
                <th style={tableHeaderStyle}>Description</th>
              </tr>
            </thead>
            <tbody>
              {[
                ['id', 'string', 'Yes', 'Unique dot-separated identifier (e.g., arith.add.single)'],
                ['label', 'string', 'Yes', 'Human-readable name displayed in the UI'],
                ['domain', 'string', 'Yes', 'One of the 16 domain categories'],
                ['subdomain', 'string', 'Optional', 'Nested grouping within a domain'],
                ['grading_type', 'enum', 'Yes', 'numeric for arithmetic, polynomial for algebra'],
                ['prerequisites', 'string[]', 'Yes', 'Concept IDs that must be mastered first'],
                ['streak', 'number', 'Yes', 'Consecutive correct answers required for mastery'],
                ['avg_time_seconds', 'number', 'Yes', 'Maximum acceptable average response time'],
              ].map(([field, type, req, desc], i) => (
                <tr key={i} style={{ background: 'transparent' }}>
                  <td style={{ ...tableCellStyle, fontFamily: monoFont, fontSize: '0.8rem' }}>{field}</td>
                  <td style={tableCellStyle}>{type}</td>
                  <td style={tableCellStyle}>
                    <span style={{
                      fontFamily: monoFont,
                      fontSize: '10px',
                      color: req === 'Yes' ? 'var(--accent-green)' : 'var(--text-muted)',
                      border: '0.5px solid',
                      borderColor: req === 'Yes' ? 'var(--accent-green)' : 'var(--border)',
                      padding: '1px 6px',
                    }}>
                      {req}
                    </span>
                  </td>
                  <td style={tableCellStyle}>{desc}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        <div style={{
          borderLeft: '2px solid var(--accent-gold)',
          paddingLeft: '1.5rem',
          fontFamily: bodyFont,
          fontStyle: 'italic',
          fontSize: '0.9rem',
          color: 'var(--text-muted)',
          lineHeight: 1.7,
          marginTop: '1rem',
        }}>
          The prerequisites list is the most important field. What must a student absolutely
          know before attempting this? If in doubt, add the prerequisite — the graph validator
          will catch cycles.
        </div>
      </section>

      <AsciiDivider pattern="dash" />

      {/* Step 2 */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>Step 2: Teach Mathua to ask questions</h2>
        <p style={bodyStyle}>
          A generator is a Go function that produces a unique problem every time it's called.
          There is no static question bank — every problem is built on demand. Generators live in{' '}
          <code style={{ fontFamily: monoFont, fontSize: '0.9em', color: 'var(--accent-gold)' }}>internal/generator/[domain]/</code> and
          implement the Generator interface.
        </p>
        <pre style={codeBlockStyle}>
{`package arithmetic

import "mathua/internal/generator"

type AddSingleGen struct{}

func (g *AddSingleGen) Generate(difficulty float64) generator.Problem {
    max := int(5 + difficulty*4)
    a   := rand.Intn(max) + 2
    b   := rand.Intn(max) + 2
    ans := a + b
    return generator.Problem{
        Question:    fmt.Sprintf("%d + %d = ?", a, b),
        Answer:      strconv.Itoa(ans),
        Explanation: fmt.Sprintf("%d + %d = %d", a, b, ans),
    }
}`}
        </pre>
        <h3 style={h3Style}>A few guidelines</h3>
        {[
          'Use the difficulty parameter to scale operand sizes. At 0.0, trivial. At 1.0, challenging for someone at that level.',
          'Always return an Explanation — it is shown when a student asks to see the solution.',
          'Use crypto/rand or math/rand with a seeded source. No hardcoded problems.',
          'Stay deterministic with respect to difficulty. A student should not get a meaningfully harder problem at the same difficulty.',
        ].map((rule, i) => (
          <div key={i} style={{ ...bodyStyle, marginBottom: '0.4rem' }}>
            <span style={{ color: 'var(--border-strong)', marginRight: '0.25rem', fontFamily: monoFont }}>·</span>
            {rule}
          </div>
        ))}
      </section>

      <AsciiDivider pattern="wave" />

      {/* Step 3 */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>Step 3: Prove it works</h2>
        <p style={bodyStyle}>
          Every generator needs a fuzz test that asserts 1 000 valid samples. This catches
          edge cases — division by zero, negative operand ranges, malformed output — before
          a student ever sees them.
        </p>
        <pre style={codeBlockStyle}>
{`func TestAddSingleGen(t *testing.T) {
    g := &AddSingleGen{}
    for i := 0; i < 1000; i++ {
        p := g.Generate(rand.Float64())
        assert.NotEmpty(t, p.Question)
        assert.NotEmpty(t, p.Answer)
        assert.NotEmpty(t, p.Explanation)
        assert.True(t, evaluateNumeric(p))
    }
}`}
        </pre>
        <p style={bodyStyle}>
          Run your tests before opening a PR:
        </p>
        <pre style={codeBlockStyle}>
          go test ./internal/generator/... -count 1000
        </pre>
      </section>

      <AsciiDivider pattern="dash" />

      {/* Validator */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>The graph validator</h2>
        <p style={bodyStyle}>
          A validator runs on every pull request. It checks two invariants before any
          merge can happen:
        </p>
        {[
          'No cycles — concept A cannot require B while B requires A.',
          'No orphans — every prerequisite must exist in the graph.',
        ].map((rule, i) => (
          <div key={i} style={{ ...bodyStyle, marginBottom: '0.4rem' }}>
            <span style={{ color: 'var(--accent-gold)', fontFamily: monoFont, fontSize: '13px' }}>
              {i + 1}.
            </span>{' '}
            {rule}
          </div>
        ))}
        <pre style={codeBlockStyle}>
          go run scripts/validate_graph.go
        </pre>
      </section>

      <AsciiDivider pattern="wave" />

      {/* Submitting a PR */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>Submitting a pull request</h2>
        <p style={bodyStyle}>
          Once everything passes locally, here's the full checklist:
        </p>
        {[
          'Add the concept to concepts.json with correct prerequisites.',
          'Write the generator in the appropriate domain subdirectory.',
          'Write the fuzz test with 1 000 samples.',
          'Run go test ./... and go run scripts/validate_graph.go locally.',
          'Open a PR. The CI pipeline runs the validator and all tests automatically.',
          'A maintainer reviews the concept ordering, thresholds, and generator quality.',
        ].map((step, i) => (
          <div key={i} style={{ ...bodyStyle, marginBottom: '0.4rem' }}>
            <span style={{ color: 'var(--accent-gold)', fontFamily: monoFont, fontSize: '13px' }}>
              {i + 1}.
            </span>{' '}
            {step}
          </div>
        ))}
        <p style={{ ...bodyStyle, marginTop: '1rem' }}>
          Reviews usually happen within a few days. If a week passes with no response, feel
          free to ping the thread. We read every PR.
        </p>
      </section>

      <AsciiDivider pattern="dash" />

      {/* Design conventions */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>Design conventions</h2>
        {[
          'Concept IDs follow domain.subdomain.descriptor — lower case, no spaces.',
          'Mastery thresholds are pragmatic. Single-digit addition should require faster response (6–8 s) than multi-digit multiplication (15–20 s).',
          'Subdomains group related concepts. If a domain grows past 15 concepts, consider introducing subdomains.',
          'Difficulty scaling should be linear where sensible. The jump from 0.0 to 1.0 should feel meaningful, not extreme.',
        ].map((rule, i) => (
          <div key={i} style={{ ...bodyStyle, marginBottom: '0.4rem' }}>
            <span style={{ color: 'var(--border-strong)', marginRight: '0.25rem', fontFamily: monoFont }}>·</span>
            {rule}
          </div>
        ))}
        <div style={{
          borderLeft: '2px solid var(--accent-gold)',
          paddingLeft: '1.5rem',
          fontFamily: bodyFont,
          fontStyle: 'italic',
          fontSize: '0.9rem',
          color: 'var(--text-muted)',
          lineHeight: 1.7,
          marginTop: '1.5rem',
        }}>
          Read the full design system in{' '}
          <a href="https://github.com/chuma-beep/mathua/blob/main/DESIGN.md" style={{ color: 'var(--accent-gold)', textDecoration: 'underline' }}>
            DESIGN.md
          </a>{' '}
          for color palette, typography, spacing, and component patterns.
        </div>
      </section>
    </div>
  )
}
