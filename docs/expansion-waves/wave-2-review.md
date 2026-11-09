# Wave 2 review

Date: 2026-11-09. Scope: the 50 concepts added by Wave 2, from 1,550 to 1,600.

This file records the inspection result as a metric so later waves can be compared against it.
Passing gates is technical correctness. The counts below come from reading generated questions,
answers and explanations, not from the test suite.

## Defect metrics

| Metric | Wave 2 |
|---|---|
| Concepts inspected | 50 |
| Knowledge-point families inspected | 150 |
| Mathematical answer defects | 1 |
| Presentation defects | 13 |
| Duplicate concepts replaced before acceptance | 3 |
| Defects fixed | All identified |
| Regression tests added | Yes |
| Singleton units before and after | 24 and 25 |
| Singleton units created by the wave | 1 |

## Mathematical defect

One question had a wrong answer, not a wrong render.

1. `calc.series.cauchy_condensation` KP1 asked for a constant `k` in `b_n = 2^k` while the
   condensed term was `b_n = 2^{-17n}`. No constant `k` describes that term. The question now
   asks for `k` in `b_n = 2^{kn}`, and the explanation ends at `2^{(1-p)n}`. The test recomputes
   `1-p` from the published `p` fact rather than reading the answer.

The new tests independently recompute every answer family from generator facts. They are in
`internal/generator/algebra/wave3_test.go`, `internal/generator/calculus/wave3_test.go` and
`internal/generator/statistics/wave2_test.go`.

## Duplicate concepts replaced

The graph validator found three algebra labels that duplicated precalculus concepts. Duplicating
precalc content would have weakened the expansion, so all three were replaced with genuinely new
algebra concepts before acceptance.

| Removed | Replacement |
|---|---|
| `alg.poly.descartes` | `alg.poly.reciprocal` |
| `alg.conic.eccentricity` | `alg.conic.intersection` |
| `alg.conic.rotated` | `alg.conic.translation` |

The replacement updated the concept nodes, generators, tests, lessons, knowledge-point shards,
solution schemas, course targets, lesson ledger, variety baseline and candidate ledger.

## Presentation defects

Thirteen questions or lessons rendered a defensible answer in a way a learner should not have to
read. All were fixed before closeout.

1. `calc.series.cauchy_condensation` KP3 printed the final fraction twice.
2. `calc.series.raabe` KP2 invented an extra factor that was not in the series.
3. `calc.series.kummer` KP1 printed `4n+44` where the expansion gives `4n+16`.
4. `calc.series.kummer` KP2 printed `1n` and `n^1`.
5. `calc.series.kummer` KP3 stated the limit without a readable expansion.
6. `calc.series.uniform_convergence` KP1 explained the supremum without stating the yes or no
   answer.
7. `calc.series.uniform_convergence` KP2 left `\sum` outside math delimiters.
8. `calc.series.uniform_convergence` KP2 mentioned geometric and alternating cases that did not
   belong to the sampled series.
9. `calc.integral.gauss_legendre` KP2 contained the nonsense term `926^4/25`.
10. `calc.integral.romberg` KP3 rendered `31/3--1/3`.
11. `calc.integral.fresnel` KP3 rounded a tiny coefficient to `0.000000`.
12. `alg.conic.translation` KP3 rendered `-4 + -4`.
13. The `calc.series.abel` lesson used `Abel's test` for both its H1 title and its first section,
    so section extraction matched the title and returned the whole body.

## Singleton discrepancy

Wave 2 increased the corpus by 50 concepts and the band-width-3 partition from 138 units to 139.
Singleton units rose from 24 to 25. The extra singleton comes from the new subdomain and band
distribution, not from a change to the mastery algorithm. `data/mastery_units.report.md` owns the
live figure, the band width remains 3, and the singleton anomaly is reported without failing
`make validate`.

## Validation

Wave 2 closed with `go test ./...`, `make validate`, `git diff --check` and the documentation
sweep passing. The additive baseline remains frozen at the 1,500-concept snapshot, and no
exemption was added.
