export default function ContributingPage() {
  return (
    <div className="kami-container">
      <span className="kami-overline">Community</span>
      <h1>Contributing to Mathua</h1>
      <p>
        The most impactful contributions are new concepts and improved generators. Adding
        a concept requires exactly three things: a JSON entry in the concept graph, a Go
        generator function, and a fuzz test. This guide walks through each step.
      </p>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 'var(--kami-space-lg)', marginBottom: 'var(--kami-space-2xl)' }}>
        <div className="kami-card">
          <div className="kami-metric">
            <div className="kami-metric-number">3</div>
            <div className="kami-metric-label">Steps to contribute</div>
          </div>
        </div>
        <div className="kami-card">
          <div className="kami-metric">
            <div className="kami-metric-number">1000</div>
            <div className="kami-metric-label">Fuzz samples required</div>
          </div>
        </div>
        <div className="kami-card">
          <div className="kami-metric">
            <div className="kami-metric-number">0</div>
            <div className="kami-metric-label">Cycles tolerated</div>
          </div>
        </div>
      </div>

      <hr />

      {/* Step 1 */}
      <section>
        <h2>
          <span style={{ color: 'var(--kami-accent)', fontWeight: 500 }}>I.</span>{' '}
          Add a node to the concept graph
        </h2>
        <p>
          Every concept in Mathua lives in <code>data/concepts.json</code>. Adding a concept
          means adding a JSON object with an ID, label, domain, prerequisite list, and
          mastery thresholds.
        </p>
        <pre>{`{
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
}`}</pre>

        <table>
          <thead>
            <tr>
              <th>Field</th>
              <th>Type</th>
              <th>Required</th>
              <th>Description</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td><code>id</code></td>
              <td>string</td>
              <td>
                <span className="kami-tag kami-tag--strong" style={{ fontSize: '0.625rem', padding: '1px 6px' }}>Yes</span>
              </td>
              <td>Unique dot-separated identifier (e.g., <code>arith.add.single</code>)</td>
            </tr>
            <tr>
              <td><code>label</code></td>
              <td>string</td>
              <td>
                <span className="kami-tag kami-tag--strong" style={{ fontSize: '0.625rem', padding: '1px 6px' }}>Yes</span>
              </td>
              <td>Human-readable name displayed in the UI</td>
            </tr>
            <tr>
              <td><code>domain</code></td>
              <td>string</td>
              <td>
                <span className="kami-tag kami-tag--strong" style={{ fontSize: '0.625rem', padding: '1px 6px' }}>Yes</span>
              </td>
              <td>One of the 16 domain categories (e.g., <code>arithmetic</code>)</td>
            </tr>
            <tr>
              <td><code>subdomain</code></td>
              <td>string</td>
              <td>
                <span className="kami-tag kami-tag--medium" style={{ fontSize: '0.625rem', padding: '1px 6px' }}>Optional</span>
              </td>
              <td>Nested grouping within a domain</td>
            </tr>
            <tr>
              <td><code>grading_type</code></td>
              <td>enum</td>
              <td>
                <span className="kami-tag kami-tag--strong" style={{ fontSize: '0.625rem', padding: '1px 6px' }}>Yes</span>
              </td>
              <td><code>numeric</code> for arithmetic, <code>polynomial</code> for algebra</td>
            </tr>
            <tr>
              <td><code>prerequisites</code></td>
              <td>string[]</td>
              <td>
                <span className="kami-tag kami-tag--strong" style={{ fontSize: '0.625rem', padding: '1px 6px' }}>Yes</span>
              </td>
              <td>List of concept IDs that must be mastered first</td>
            </tr>
            <tr>
              <td><code>mastery_threshold.streak</code></td>
              <td>number</td>
              <td>
                <span className="kami-tag kami-tag--strong" style={{ fontSize: '0.625rem', padding: '1px 6px' }}>Yes</span>
              </td>
              <td>Consecutive correct answers required</td>
            </tr>
            <tr>
              <td><code>mastery_threshold.avg_time_seconds</code></td>
              <td>number</td>
              <td>
                <span className="kami-tag kami-tag--strong" style={{ fontSize: '0.625rem', padding: '1px 6px' }}>Yes</span>
              </td>
              <td>Maximum acceptable average response time</td>
            </tr>
          </tbody>
        </table>

        <blockquote>
          The prerequisite list is the most important field. Think carefully — what must a student
          absolutely know before attempting this concept? If in doubt, add the prerequisite. The
          graph validator will catch cycles.
        </blockquote>
      </section>

      <hr />

      {/* Step 2 */}
      <section>
        <h2>
          <span style={{ color: 'var(--kami-accent)', fontWeight: 500 }}>II.</span>{' '}
          Write the generator function
        </h2>
        <p>
          Each concept needs a generator — a Go function that produces a unique problem every
          time it is called. Generators live in <code>internal/generator/[domain]/</code> and
          implement the <code>Generator</code> interface.
        </p>
        <pre>{`package arithmetic

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
}`}</pre>
        <h3>Guidelines for good generators</h3>
        <ul>
          <li>Use the <code>difficulty</code> parameter to scale operand sizes. At 0.0 the problem
            should be trivial; at 1.0 it should be challenging for a student at that level.</li>
          <li>Always provide an <code>Explanation</code> field. It is shown when a student asks
            to see the solution.</li>
          <li>Use <code>crypto/rand</code> or <code>math/rand</code> with a seeded source. Do not
            hardcode problems — the generator must produce unique output each call.</li>
          <li>Keep the generator deterministic with respect to difficulty. A student should not
            get a meaningfully harder problem at the same difficulty.</li>
        </ul>
      </section>

      <hr />

      {/* Step 3 */}
      <section>
        <h2>
          <span style={{ color: 'var(--kami-accent)', fontWeight: 500 }}>III.</span>{' '}
          Write a fuzz test
        </h2>
        <p>
          Every generator must have a fuzz test that asserts 1000 valid samples. This catches
          edge cases — division by zero, negative operand ranges, malformed output — before
          they reach a student.
        </p>
        <pre>{`func TestAddSingleGen(t *testing.T) {
    g := &AddSingleGen{}
    for i := 0; i < 1000; i++ {
        p := g.Generate(rand.Float64())
        assert.NotEmpty(t, p.Question)
        assert.NotEmpty(t, p.Answer)
        assert.NotEmpty(t, p.Explanation)
        // Optionally: validate the answer programmatically
        assert.True(t, evaluateNumeric(p))
    }
}`}</pre>
        <blockquote>
          Run generator fuzz tests with <code>go test ./internal/generator/... -run TestFuzz -count 1000</code>.
          All tests must pass before the PR is merged.
        </blockquote>
      </section>

      <hr />

      {/* Validator */}
      <section>
        <span className="kami-overline">Validation</span>
        <h2>Graph validator — automatic on every PR</h2>
        <p>
          A graph validator runs on every pull request. It checks two invariants:
        </p>
        <ol>
          <li><strong>No cycles.</strong> The DAG must remain acyclic. A cycle means concept A
            requires B, which requires A — creating an impossible prerequisite chain.</li>
          <li><strong>No orphans.</strong> Every prerequisite referenced must exist in the graph.
            A reference to a non-existent concept ID is rejected.</li>
        </ol>
        <p>
          Run the validator locally before opening a PR:
        </p>
        <pre>{`go run scripts/validate_graph.go`}</pre>
      </section>

      <hr />

      {/* Pull Request Flow */}
      <section>
        <span className="kami-overline">Workflow</span>
        <h2>Pull request checklist</h2>
        <ol>
          <li>Add the concept to <code>data/concepts.json</code> with correct prerequisites.</li>
          <li>Write the generator in the appropriate domain subdirectory.</li>
          <li>Write the fuzz test with 1000 samples.</li>
          <li>Run <code>go test ./...</code> and <code>go run scripts/validate_graph.go</code> locally.</li>
          <li>Open a PR. The CI pipeline runs the validator and all tests.</li>
          <li>A maintainer reviews the concept ordering, threshold values, and generator quality.</li>
        </ol>
      </section>

      <hr />

      {/* Design Philosophy */}
      <section>
        <h2>Design conventions</h2>
        <ul>
          <li><strong>Concept IDs</strong> follow the pattern <code>domain.subdomain.descriptor</code>.
            Use snake_case, no spaces, no special characters.</li>
          <li><strong>Mastery thresholds</strong> are pragmatic. A single-digit addition should require
            faster response (6&ndash;8 seconds) than multi-digit multiplication (15&ndash;20 seconds).</li>
          <li><strong>Subdomains</strong> group related concepts. If a domain grows beyond 15 concepts,
            consider introducing subdomains.</li>
          <li><strong>Difficulty scaling</strong> should be linear where possible. The difference
            between difficulty 0.0 and 1.0 should feel meaningful but not extreme.</li>
        </ul>
        <blockquote>
          Read the full design system in <a href="https://github.com/chuma-beep/mathua/blob/main/DESIGN.md">DESIGN.md</a> for
          color palette, typography, spacing, and component patterns used across the project.
        </blockquote>
      </section>
    </div>
  )
}
