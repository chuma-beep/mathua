# Wave 9 review

Date: 2026-11-10. Scope: the 50 concepts added by Wave 9, from 1,900 to 1,950.

This file records the inspection result as a metric so later waves can be compared against it.
Passing gates is technical correctness. The counts below come from reading generated questions,
answers and explanations, not from the test suite.

## Defect metrics

| Metric | Wave 9 |
|---|---:|
| Concepts inspected | 50 |
| Knowledge-point families inspected | 150 |
| Mathematical and framing defects | 18 |
| Presentation defects | 12 |
| Data-integrity defects | 1 |
| Defects fixed | All identified |
| First-ten 200-draw families below 20 | 1 deliberate finite pool |
| Singleton units before and after | 23 and 22 |
| Singleton units created by the wave | 0 |
| Singleton units removed by the wave | 1 |

## Mathematical and framing defects

Eighteen questions tested the wrong thing, answered from the stem or were always true.

1. `precalc.vector.direction_cosines` KP1 asked for the numerator that the stem already printed.
2. `precalc.systems.nonlinear_inequality` KP2 always said the boundary was included.
3. `precalc.partial_fractions.repeated` KP1 printed the coefficient it asked for.
4. `precalc.partial_fractions.repeated` KP2 and KP3 referred to coefficients that were not fixed
   numerically.
5. `precalc.partial_fractions.quadratic` KP1 and KP2 asked for symbolic coefficients without a
   concrete numerator.
6. `precalc.rational.asymptote_intersection` KP1 always returned zero.
7. `precalc.conic.polar` KP2 asked for a numerator that the equation printed.
8. `calc.deriv.sign` KP2 asked for the jump point that the stem printed.
9. `calc.deriv.related_rates_shadow` KP3 asked for a ratio that was always 3.
10. `calc.deriv.related_rates_angle` KP1 asked for the height that the stem printed.
11. `calc.deriv.optimization_distance` KP3 always returned 1.
12. `calc.series.abel_theorem` KP3 claimed the geometric series converges at its boundary, where
    the ratio is 1 and the series diverges.
13. `calc.series.euler_maclaurin` KP3 always returned 0 for a linear polynomial.
14. `calc.series.taubian` KP3 was always true and was theorem-name recall.
15. `calc.series.lambert` KP3 always returned 1.
16. `calc.limit.continuity_composition` KP2 always said the composition is continuous.
17. `precalc.polar.symmetry` used an abstract curve type instead of naming the polar equation.
18. `precalc.conic.parametric` KP3 asked which conic the equations describe, which is
    meaning-phrasing rather than a numeric computation. It now asks for \(x\) at \(t=\pi\).

All eighteen were replaced with concrete values, conditions that can fail, or uniquely determined
coefficients. Exercises that remained theorem-shaped were given a computational or counterexample
family.

## Presentation defects

Twelve learner-visible rendering, lesson or grammar defects were fixed:

1. Sign joins such as `+-` and `x--3`.
2. Unit and zero coefficients such as `1x`, `0x+0` and `1 y`.
3. Over-escaped pure-number math spans such as `\$41\$` and `\$200\$`.
4. Negative squares written as `-5^2` instead of `(-5)^2`.
5. Products such as `3\cdot-2^2` without parentheses.
6. Singular `1 units` and `1 nearest minority neighbours`.
7. Explanations that used an undefined `a` after the stem had switched to numeric coefficients.
8. Abstract polar symmetry-type wording.
9. A generic closing paragraph shared by the forty new lessons.
10. A sanitizer regression that briefly removed the numeral from “Enter 1 for yes”.
11. `precalc.vector.direction_cosines` used the H1 title as its first H2, so section extraction
    fell back to the whole lesson. The heading is now “The direction-cosine formula”.
12. `calc.series.lambert` used the H1 title as its first H2, so section extraction fell back to
    the whole lesson. The heading is now “The Lambert-series form”.

The final sanitizer dump had no remaining `+-`, `--number`, `\cdot-`, unit-coefficient, pure-number
math-span or “Enter for yes” patterns. The sanitizer was narrowed until it left yes/no instructions
intact.

## Data-integrity defect

The remaining-forty authoring script appended seventeen calculus nodes to `data/concepts/calc1.json`
and `calc2.json` with `"domain": "calculus"`. The graph validator reported 24 domains instead of 23.
The eleven calc1 nodes and six calc2 nodes were corrected to their file domains. The validator then
reported 1,950 concepts across 23 domains and 19 reviewed label collisions.

## Pool and parameterisation work

The first ten were measured at 200 draws per ordinal:

| Concept | KP1 | KP2 | KP3 |
|---|---:|---:|---:|
| precalc.func.one_to_one | 182 | 178 | 44 |
| precalc.poly.complex_factorization | 180 | 179 | 178 |
| precalc.vector.orthogonal_decomposition | 180 | 182 | 185 |
| calc.limit.asymptotic_equivalence | 190 | 178 | 27 |
| calc.deriv.elasticity | 46 | 47 | 47 |
| calc.integral.leibniz | 30 | 30 | 30 |
| calc.series.cesaro | 24 | 69 | 48 |
| calc.series.dirichlet_series | 20 | 24 | 30 |
| frac.mixed.add_carry | 196 | 195 | 193 |
| ml.eval.log_loss | 20 | 114 | 123 |

Every first-ten family reached 20 distinct questions in 200 draws. The remaining forty families
were verified at the 5-in-12 hard floor and the 10-in-60 minimum. They were not individually
profiled at 200 draws. Finite pools are documented in the generators rather than filled with
cosmetic variations.

## Singleton discrepancy

Wave 8 closed with 142 units and 23 singleton units at 1,900 concepts. Wave 9 closes with 143 units
and 22 singleton units at 1,950 concepts.

- No singleton unit was added.
- `precalculus.b3`, which contained only `precalc.analytic.distance`, is no longer a singleton. It
  was absorbed into a multi-concept band because the new precalculus concepts and prerequisite
  edges changed topological depths.
- The rule and band width are unchanged: `unit = (domain, floor(conceptDepth / bandWidth))` at
  width 3.

The net change is caused by the added concepts changing existing depths, not by a change to the
partition algorithm.

## Approved split and artifact scope

Wave 9 finished precalculus at 81, calc1 at 114, calc2 at 97, fractions at 67 and machine learning
at 67. It added fifty concepts: 23 precalculus, 14 calc1, 8 calc2, 3 fractions and 2 machine
learning. Courses are `pc`, `calc1`, `calc2`, `pa` with grade 6, and `ml`.

Precalculus, calc1, calc2 and machine learning are not learner domains, so they add no solution
schemas. The three fractions concepts are in a learner domain, and each has an authored schema.
The first ten were approved after an independent mathematical audit. The remaining forty were
authored after that approval.

## Quality limits

Several theorem-heavy concepts are assessed through finite, polynomial or piecewise examples. This
is deliberate: a proof-recognition question was replaced with a computation, a condition check or a
counterexample. The remaining risk is that a learner can pass a theorem-heavy family by computing a
special case without seeing the full theorem. The calculus series families and the `precalc.polar`
families deserve the most scrutiny in a later quality pass.

The forty remaining lessons use concept-specific introductions and sections plus a shared closing
paragraph. That paragraph is boilerplate. The lessons pass the length and section checks, but a
later editorial pass should replace the closing paragraph with concept-specific prose.

## Validation

Wave 9 closes with 1,950 concepts, 1,950 lessons, 1,950 KP files, 5,850 knowledge points, 28 courses
and all concepts reachable. The graph validator reports 23 domains and 19 reviewed collisions. The
lesson audit reports zero issues. Schema, shipping, sweep, duplicate-registration, counts, additive
and published-docs checks pass. `make validate`, `go test ./...` and `git diff --check` pass.

The additive baseline remains frozen at the 1,500-concept snapshot. The commit is local. Nothing was
pushed or deployed, and the two unrelated server files remain unstaged.
