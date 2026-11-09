# Wave 1 review

Date: 2026-11-09. Scope: the 50 concepts added by Wave 1, from 1,500 to 1,550.

This file records the inspection result as a metric so later waves can be compared against it.
Passing gates is technical correctness. The counts below come from reading generated questions,
answers and explanations, not from the test suite.

## Defect metrics

| Metric | Wave 1 |
|---|---|
| Concepts inspected | 50 |
| Knowledge-point families inspected | 150 |
| Mathematical answer defects | 5 |
| Presentation defects | 12 |
| Defects fixed | All identified |
| Regression tests added | Yes |
| Singleton units before and after | 24 and 24 |
| Singleton units created by the wave | 0 |

## Mathematical defects

Five questions had a wrong answer, not a wrong render. One was in the first ten and four were in
the remaining forty.

1. `alg.poly.turning_points` KP3 asked for the least degree given a number of turning points and
   answered `degree - 1` instead of `turning points + 1`.
2. `calc.deriv.piecewise` KP3 asked for the intercept that makes the branches meet and answered
   the slope.
3. `numa.interp.piecewise_linear` answered a slope that assumed an interval width of 2.
4. `numa.linear.matrix_condition` could truncate a non-integer condition number.
5. `numa.linear.inverse_power` could choose a shift equidistant from two eigenvalues.

Each one is covered by a test that independently recomputes the answer from the generator facts.
The tests are in `internal/generator/algebra/wave1_test.go` and `wave2_test.go`,
`internal/generator/numerical/wave1_test.go` and `wave2_test.go`, and
`internal/generator/calculus/wave2_test.go`.

## Presentation defects

Twelve questions rendered a defensible answer in a way a learner should not have to read: double
signs such as `x - -1`, unit coefficients such as `+ 1y`, `1 bisection steps`, singular nouns
after a count, a wrong question word, and explanations that stated a rule without showing the
arithmetic. All were fixed before closeout.

## Singleton discrepancy

Wave 1 changed the corpus but not the partition's singleton set: 24 before, 24 after, the same
24 units. The ADR-050 figure of 19 was a 1,000-concept measurement that a count rewrite made look
current. The correction is recorded in ADR-050 under the 2026-11-09 amendment, and the generated
`data/mastery_units.report.md` owns the live figure.
