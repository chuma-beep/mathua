export default function ArchitecturePage() {
  return (
    <div className="kami-container">
      <span className="kami-overline">Documentation</span>
      <h1>Architecture</h1>
      <p>
        Mathua is a single Go binary with two delivery modes. The engine core is identical —
        only the presentation layer differs. This document describes each layer and how they compose.
      </p>

      {/* High-Level Diagram */}
      <section style={{ marginTop: 'var(--kami-space-2xl)', marginBottom: 'var(--kami-space-3xl)' }}>
        <pre style={{ lineHeight: 1.4 }}>{`┌───────────────────────────────────────────────────────────┐
│  Web browser (React + Next.js)  │  Desktop TUI (Bubble Tea) │
├───────────────────────────────────────────────────────────┤
│      REST API — net/http  (TUI calls engine directly)     │
├───────────────────────────────────────────────────────────┤
│   DAG Loader   │   Scheduler    │   Mastery   │  Generators│
│   Path Finder  │  CAT Engine    │   Timer     │  Scoring   │
├───────────────────────────────────────────────────────────┤
│  Numeric grader (Go)  │  Symbolic grader (SymPy, future)  │
├───────────────────────────────────────────────────────────┤
│  SQLite (desktop)     │  Postgres (web)                   │
│  concepts.json — community-editable DAG                   │
└───────────────────────────────────────────────────────────┘`}</pre>
      </section>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 'var(--kami-space-lg)', marginBottom: 'var(--kami-space-2xl)' }}>
        <div className="kami-card">
          <div className="kami-metric">
            <div className="kami-metric-number">1</div>
            <div className="kami-metric-label">Binary</div>
          </div>
        </div>
        <div className="kami-card">
          <div className="kami-metric">
            <div className="kami-metric-number">2</div>
            <div className="kami-metric-label">UI Modes</div>
          </div>
        </div>
        <div className="kami-card">
          <div className="kami-metric">
            <div className="kami-metric-number">5</div>
            <div className="kami-metric-label">Engine Modules</div>
          </div>
        </div>
      </div>

      <hr />

      {/* Layer 1: UI */}
      <section>
        <span className="kami-overline">Layer 1</span>
        <h2>UI &mdash; Two presentations, one engine</h2>
        <p>
          Mathua ships with two user interfaces that share the same engine core. The desktop
          interface requires no account, no network, and stores everything in SQLite. The web
          interface uses React with Next.js and connects to a shared Postgres database.
        </p>

        <table>
          <thead>
            <tr>
              <th>Property</th>
              <th>Desktop TUI</th>
              <th>Web App</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td>Framework</td>
              <td>Bubble Tea (Go)</td>
              <td>React + Next.js</td>
            </tr>
            <tr>
              <td>Account</td>
              <td>Not required</td>
              <td>Required</td>
            </tr>
            <tr>
              <td>Network</td>
              <td>Fully offline</td>
              <td>Requires server</td>
            </tr>
            <tr>
              <td>Storage</td>
              <td>SQLite</td>
              <td>PostgreSQL</td>
            </tr>
            <tr>
              <td>Leaderboard</td>
              <td>Local only</td>
              <td>Global weekly</td>
            </tr>
            <tr>
              <td>Graph View</td>
              <td>ASCII tree</td>
              <td>Interactive 3D</td>
            </tr>
          </tbody>
        </table>
      </section>

      <hr />

      {/* Layer 2: API */}
      <section>
        <span className="kami-overline">Layer 2</span>
        <h2>API &mdash; Four services, one router</h2>
        <p>
          The web server exposes a REST API through Go&rsquo;s standard <code>net/http</code> package.
          No external HTTP framework is used. The API is organized into four services:
        </p>

        <ul>
          <li><strong>Session Service</strong> — Starts and manages practice sessions. Returns the next problem based on scheduler output.</li>
          <li><strong>Graph Service</strong> — Exposes the concept DAG for visualization. Returns node metadata, prerequisite edges, and domain groupings.</li>
          <li><strong>Leaderboard Service</strong> — Computes and returns the weekly leaderboard. Scores reset every Monday at 00:00 UTC.</li>
          <li><strong>Diagnostic Service</strong> — Runs the Computerised Adaptive Testing engine to locate a student&rsquo;s knowledge frontier.</li>
        </ul>

        <blockquote>
          The TUI bypasses the API layer entirely — it calls the engine modules directly through Go function interfaces. The API layer exists solely for the web client.
        </blockquote>
      </section>

      <hr />

      {/* Layer 3: Core Engine */}
      <section>
        <span className="kami-overline">Layer 3</span>
        <h2>Core Engine &mdash; Five modules</h2>

        <h3>DAG Loader</h3>
        <p>
          Loads <code>concepts.json</code> and builds an in-memory directed acyclic graph. Validates that
          the graph contains no cycles and no orphaned prerequisites. Runs on every server start
          and on every PR through CI. The in-memory representation supports efficient topological
          sorting and prerequisite lookups.
        </p>

        <h3>Scheduler</h3>
        <p>
          The scheduler selects the next concept after every submitted answer. It computes a priority
          score for every eligible concept — those whose prerequisites are all mastered. The priority
          formula weights recency and mastery, with bonuses for decaying concepts and newly unlocked
          material. Three hard rules apply: prerequisites must be mastered, no repeats, and a 70/30
          balance between new material and review.
        </p>
        <pre>{`priority = (0.7 × days_since_last_seen)
         + (0.3 × (1.0 − mastery_score))
         + 5.0  if status == DECAYING
         + 2.0  if newly unlocked`}</pre>

        <h3>Mastery &amp; Spaced Repetition (SM-2)</h3>
        <p>
          Each concept tracks a mastery score based on streak and response time. When mastered,
          a concept enters a spaced repetition schedule using a simplified SM-2 algorithm. The
          interval grows with each successful review and resets on failure. Reviews are woven into
          normal sessions — there is no separate &ldquo;review mode.&rdquo;
        </p>
        <pre>{`first review:   1 day
second review:  3 days
third review:   interval × 2.5
interval minimum: 1
ease factor:     2.5, min 1.3
decrease on fail: −0.2`}</pre>

        <h3>Generators</h3>
        <p>
          Every problem is generated on demand by a parameterised Go function — a generator.
          There is no static question bank. Each generator accepts a difficulty parameter
          (0.0–1.0) and returns a problem struct with a question, answer, and explanation.
          The scheduler passes difficulty based on the student&rsquo;s current mastery score.
        </p>
        <pre>{`type Generator interface {
    Generate(difficulty float64) Problem
}

type Problem struct {
    Question    string
    Answer      string
    Explanation string
    Difficulty  float64
}`}</pre>

        <h3>Scoring &amp; Leaderboard</h3>
        <p>
          Mathua tracks two scores. The lifetime topic score reflects permanent mastery — it never
          decreases. The weekly score resets every Monday and is used for the global leaderboard.
          Both incorporate concept counts, speed bonuses, and streak multipliers.
        </p>
        <pre>{`topic_score  = (mastered_concepts × 100)
              + Σ speed_bonus
weekly_score = (mastered_this_week × 100)
              + speed_bonus
              + (current_streak × 10)`}</pre>
      </section>

      <hr />

      {/* Layer 4: Grading */}
      <section>
        <span className="kami-overline">Layer 4</span>
        <h2>Grading &mdash; Numeric and symbolic</h2>
        <p>
          Mathua uses two grading strategies, selected by concept type. Arithmetic concepts
          use the numeric grader — simple string or integer comparison. Algebra concepts use
          a pure-Go polynomial grader with no external dependencies.
        </p>
        <table>
          <thead>
            <tr>
              <th>Grader</th>
              <th>Used for</th>
              <th>Strategy</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td>Numeric</td>
              <td>Counting, arithmetic, fractions</td>
              <td>String / integer equality</td>
            </tr>
            <tr>
              <td>Polynomial</td>
              <td>Algebra (future)</td>
              <td>Parse → expand → normalise → compare coefficients</td>
            </tr>
          </tbody>
        </table>
        <blockquote>
          The polynomial grader handles factoring, expansion, and equation-solving in pure Go. No Python, no SymPy, no subprocess dependencies. The same grader runs identically in the web server and desktop TUI.
        </blockquote>
      </section>

      <hr />

      {/* Layer 5: Storage */}
      <section>
        <span className="kami-overline">Layer 5</span>
        <h2>Storage &mdash; SQLite and Postgres</h2>
        <p>
          The desktop TUI uses SQLite for fully offline operation — no installation, no
          configuration, no network. The web server uses PostgreSQL for concurrent access
          and global state. The data schema is identical across both databases, and the
          engine core uses a repository interface that abstracts the storage backend.
        </p>
        <table>
          <thead>
            <tr>
              <th>Entity</th>
              <th>Stored in</th>
              <th>Description</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td>Concept graph</td>
              <td>concepts.json</td>
              <td>Community-editable, validated on PR</td>
            </tr>
            <tr>
              <td>Student progress</td>
              <td>SQLite / Postgres</td>
              <td>Per-concept state, mastery, intervals</td>
            </tr>
            <tr>
              <td>Session history</td>
              <td>SQLite / Postgres</td>
              <td>Every problem, answer, and timing</td>
            </tr>
            <tr>
              <td>Weekly scores</td>
              <td>Postgres only</td>
              <td>Global leaderboard, reset every Monday</td>
            </tr>
          </tbody>
        </table>
      </section>

      <hr />

      {/* CAT Engine */}
      <section>
        <span className="kami-overline">Algorithm</span>
        <h2>Diagnostic — Computerised Adaptive Testing</h2>
        <p>
          The diagnostic engine locates a student&rsquo;s knowledge frontier using binary search on
          the topologically sorted concept graph. This reduces the assessment from 284 questions
          (one per concept) to approximately 20&ndash;35.
        </p>
        <ol>
          <li>The concept graph is sorted topologically. The diagnostic starts at the midpoint.</li>
          <li>Correct answers within the time limit move the probe forward toward harder concepts.</li>
          <li>Incorrect or slow answers move backward toward foundational material.</li>
          <li>After 3 consecutive correct answers in a region, the frontier is considered located.</li>
          <li>The diagnostic records a mastery estimate for every concept passed through.</li>
        </ol>
        <blockquote>
          The diagnostic can be retaken at any time. Retaking does not delete progress — it creates
          a new estimate that is merged with existing data, always preferring the more optimistic
          estimate.
        </blockquote>
      </section>
    </div>
  )
}
