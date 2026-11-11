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

The generator tests parse every answer as an integer, reject a question that names
a theorem and check that the explanation is present. The answers themselves are
computed from drawn parameters in the generator, and the tests sample each ordinal
twenty times, so a solution that ignored its parameters would fail the variety
gate rather than the recomputation assertion. No answer is a stored constant.

## Defects found and fixed during the wave

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
