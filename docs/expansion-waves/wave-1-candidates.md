# Wave 1 candidate ledger

Wave 1 is the mandatory 50-concept pilot for the 2,000-concept target. This ledger is the survey
before authoring: it fixes the candidate list, the domain split and the working labels so that two
authors do not propose the same concept. It is not an authored wave. No corpus file outside this
ledger changes until the first ten candidates are reviewed.

Every candidate id and label below was checked against `data/concepts/` on 2026-11-08 and has no
exact collision. Prerequisites are deliberately left to authoring: the correct edge depends on the
lesson and the generator, and a guessed edge is worse than a missing one. The scaffold validator
will refuse a manifest that omits them.

## Split

| Domain | Candidates | Target row |
|---|---:|---|
| algebra | 14 | 139 to 171 (+32) |
| numerical_analysis | 20 | 49 to 69 (+20) |
| calc1 | 16 | 84 to 114 (+30) |
| **Total** | **50** | |

The split is intentional. Algebra is schema-heavy and prose-heavy, numerical analysis is computed
and method-heavy, and calculus is application-heavy. If the pilot cannot sustain all three, the
remaining waves do not start.

## Review order

1. Author the first ten: five algebra, three numerical analysis and two calculus.
2. Run the full gates on those ten, including the additive check, the variety ratchet and the
   mastery-unit partition. Measure the singleton count; do not normalize it.
3. Review the ten against the artifact contract. Only then author the remaining forty.
4. If the ten expose a scaffolding or schema problem, fix the process before the forty. A failed
   pilot is a cheaper outcome than 500 concepts built on the wrong shape.

## Algebra: 14 candidates

| Concept id | Working label | Draft family |
|---|---|---|
| `alg.poly.turning_points` | Turning points of a polynomial graph | graph/read |
| `alg.func.piecewise_graph` | Graph piecewise-defined functions | graph/read |
| `alg.systems.graphing` | Solve systems of linear equations by graphing | graph/read |
| `alg.systems.cramer` | Cramer's rule for 2x2 systems | computed |
| `alg.poly.remainder` | The remainder theorem | computed |
| `alg.poly.binomial` | The binomial theorem | symbolic |
| `alg.seq.geometric_infinite` | Sum of an infinite geometric series | computed |
| `alg.exp.from_points` | Write an exponential function from two points | symbolic |
| `alg.log.scale` | Logarithmic scales: decibels and pH | word |
| `alg.regression.linear` | Line of best fit by least squares | computed |
| `alg.conic.parabola_focus` | Focus and directrix of a parabola | symbolic |
| `alg.conic.ellipse_foci` | Foci and eccentricity of an ellipse | computed |
| `alg.ineq.absolute_graph` | Graph absolute value inequalities in two variables | graph/read |
| `alg.poly.zeros_multiplicity` | Multiplicity of zeros and graph behavior | graph/read |

## Numerical analysis: 20 candidates

| Concept id | Working label | Draft family |
|---|---|---|
| `numa.errors.backward` | Backward error and backward stability | computed |
| `numa.errors.floating` | Floating-point arithmetic and rounding modes | computed |
| `numa.roots.brent` | Brent's method | method |
| `numa.roots.aitken` | Aitken acceleration | computed |
| `numa.roots.muller` | Muller's method | method |
| `numa.interp.barycentric` | Barycentric interpolation | computed |
| `numa.interp.runge` | Runge's phenomenon | graph/read |
| `numa.interp.piecewise_linear` | Piecewise linear interpolation | computed |
| `numa.integration.boole` | Boole's rule | computed |
| `numa.integration.gauss_kronrod` | Gauss-Kronrod quadrature | computed |
| `numa.linear.gauss_jordan` | Gauss-Jordan elimination | method |
| `numa.linear.thomas` | The Thomas algorithm for tridiagonal systems | method |
| `numa.linear.matrix_condition` | Condition number of a matrix | computed |
| `numa.linear.inverse_power` | The inverse power method | method |
| `numa.fourier.dft` | The discrete Fourier transform | computed |
| `numa.fourier.aliasing` | Aliasing and the Nyquist rate | word |
| `numa.approx.pade` | Pade approximants | symbolic |
| `numa.diff.finite` | Finite-difference derivatives | computed |
| `numa.diff.richardson` | Richardson extrapolation | computed |
| `numa.optimization.line_search` | Line search methods | method |

## Calculus 1: 16 candidates

| Concept id | Working label | Draft family |
|---|---|---|
| `calc.limit.piecewise` | Limits of piecewise functions | computed |
| `calc.limit.continuity_interval` | Continuity on an interval | graph/read |
| `calc.limit.conjugate` | Limits by rationalizing | symbolic |
| `calc.limit.complex_fraction` | Limits of complex fractions | symbolic |
| `calc.limit.absolute_value` | Limits involving absolute value | symbolic |
| `calc.deriv.normal_line` | Equations of normal lines | symbolic |
| `calc.deriv.piecewise` | Differentiability of piecewise functions | computed |
| `calc.deriv.optimization_fence` | Optimization: fencing and area | word |
| `calc.deriv.optimization_revenue` | Optimization: revenue and cost | word |
| `calc.deriv.related_rates_ladder` | Related rates: the sliding ladder | word |
| `calc.deriv.related_rates_cone` | Related rates: filling a cone | word |
| `calc.integral.total_distance` | Total distance and displacement | word |
| `calc.integral.piecewise` | Integrals of piecewise functions | computed |
| `calc.integral.absolute_value` | Integrals involving absolute value | symbolic |
| `calc.integral.accumulation` | Accumulation functions | graph/read |
| `calc.integral.motion` | Motion problems with definite integrals | word |

## What each candidate must become

Each row becomes a concept node, a three-family generator, a lesson over 1 KB with three
section-aligned headings, three KP shards, an authored solution schema in the learner domains, a
deliberate course target and the regenerated derived artifacts. The contract is section 5.3 of
`docs/content-expansion.md`. A candidate that cannot meet it is dropped before authoring, not
shipped weak.
