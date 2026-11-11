# Wave 11 review

Status: implemented on 2026-11-11.

Wave 11 raises the corpus from 2,000 to 2,050 concepts and creates the Optimization
domain with one university course, `optim`.

## Delivered

- 50 concepts in `data/concepts/optimization.json`, split across five subdomains:
  foundations 10, unconstrained 10, convex and duality 12, constrained 10 and
  discrete 8.
- 150 question families in `internal/generator/optimization`, three per concept.
  The package registers through `internal/generator/all`, so the registry now holds
  23 domain packages.
- 50 lessons in `data/lessons/authored/`, 50 KP shards in `data/lessons/kp/` and 50
  rows in `data/lessons/lessons.json`.
- The `optim` course in `data/courses.json`. Its five targets are
  `optim.bfgs`, `optim.active_set`, `optim.lp_sensitivity`, `optim.integrality_gap`
  and `optim.trust_region_subproblem`. The draft target list named four different
  capstones and left 29 concepts outside every course closure, so the targets were
  chosen from the measured closures instead. All 50 concepts are reachable.

## Measurements

| Gate | Result |
|---|---|
| Concepts | 50 |
| Families | 150 |
| Distinct questions per ordinal in 12 draws | at least 5 in every family |
| Distinct questions per ordinal in 200 draws | at least 20 in every family |
| Variety baseline, minimum of five 60-draw samples | 33 |
| Variety baseline, median | 51 |
| Smallest lesson | 1,326 bytes |
| Largest lesson | 2,107 bytes |
| Lesson headings | one H1 and three H2 sections in every lesson |
| KP shards | 50 files, three entries, three subgoals each |
| Solution schemas | none, because Optimization is a university domain |
| Mastery bandwidth | 3 |
| Backward band edges | 0 |
| Singleton mastery units | 20, unchanged from the pre-wave corpus |
| STE warnings | 37, unchanged advisory output |
| `make validate` | passes |
| `go test ./...` | passes |
| `git diff --check` | clean |

The variety figures are the recorded floors in
`internal/generator/all/variety_baseline.json`, which stores the lowest of five
samples because the generators draw from the global random source.

## Independent recomputation

`TestRepairedFamiliesRecompute` in `internal/generator/optimization/recompute_test.go`
parses each repaired instance back out of the question text and derives the answer a
second time, independently of the generator. It covers eighteen family ordinals,
including every family that was found to be wrong. `TestFactsAndVariety` checks that
`Facts["answer"]` matches the graded answer and that no ordinal collapses to fewer
than three distinct answers in sixty draws.

The remaining families are covered by the parser gates, which reject a non-integer
answer, a question that names a theorem and a missing explanation, and by the
variety floor. They are not independently recomputed. An earlier version of this
document claimed that every answer was recomputed; that claim was wrong, because
the test only parsed the number, and it is retracted here.

## Defects found and fixed during the wave

- Nine families returned a wrong answer and the suite accepted it because the test
  only parsed the number. They are corrected: the doubled midpoint coordinate sum,
  the inverted strict-convexity sign, integer truncation of a stationary point, the
  inverted Newton-to-gradient ratio, the Wolfe slope that omitted the factor of two,
  the Wolfe step that ignored the derivative, a trust-region family that assumed the
  minimizer was outside the ball, the actual reduction computed with the wrong sign
  of step, and the swapped sensitivity coefficient.
- Explanations that contradicted the drawn instance were rewritten, and four
  families that printed their answer in the question or returned a constant were
  replaced with computed quantities.
- The draft course targets did not close over the domain. Replacing them with the
  measured capstones restored full reachability.
- A penalty family looped without bound when its threshold sat below half the
  constraint value. The family now inverts the violation bound directly.
- Twelve families had fewer than 20 instances at one ordinal because a parameter
  range was too narrow. The ranges were widened and the test now enforces the
  design target at every ordinal.
- One question phrased a set as "described as", which the diagnostic meaning
  phrasing lint rejects for numeric grading. It now says "consists of".
- New concepts carry `key_prerequisites`, because 50 bare concepts would have
  pushed curriculum key-prerequisite coverage below its 80 percent floor.

## Boundaries kept

The domain overlaps existing material by design. `optim.first_order` states the
multivariable condition rather than repeating one-variable calculus.
`optim.second_order` uses the Hessian rather than a one-variable test.
`optim.lp_dual` builds a dual the algebra course does not. `optim.kkt` adds
inequalities to the equality-only Lagrange material. `optim.gradient_step`
computes a step rather than describing a training loop. `optim.transportation`
uses supplies and demands rather than the assignment problem.
