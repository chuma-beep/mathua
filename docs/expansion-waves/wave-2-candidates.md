# Wave 2 candidate ledger

Wave 2 raises the corpus from 1,550 to 1,600 concepts. This ledger fixes the 50 candidates and
their provisional prerequisite edges before authoring. Every id and label below was checked against
`data/concepts/` on 2026-11-09: there is no exact id collision and no normalized-label collision.

The split finishes algebra (+18), extends Calculus II into the convergence tests and quadrature
methods its target calls for (+18), and adds the inference and distribution concepts statistics is
missing (+14). Course targets are `a2`, `calc2` and `probstat`.

## Algebra: 18 candidates

| id | Label | Subdomain | Prerequisite |
|---|---|---|---|
| alg.poly.reciprocal | Reciprocal polynomials | algebra.polynomials | alg.poly.vieta |
| alg.poly.bounds | Upper and lower bounds for roots | algebra.polynomials | alg.poly.roots |
| alg.func.operations | Combining functions | algebra.func | alg.func.notation |
| alg.func.inverse_graph | Graphs of inverse functions | algebra.func | alg.func.inverse |
| alg.func.composition_graph | Composition from graphs | algebra.func | alg.func.composite, alg.func.graph_analysis |
| alg.exp.regression | Exponential regression | algebra.linear | alg.regression.linear, alg.exp.from_points |
| alg.log.domain | Domain of logarithmic functions | algebra.logarithms | alg.log.concept, alg.func.domain |
| alg.log.base_equation | Solving for a logarithm's base | algebra.logarithms | alg.eq.logarithmic |
| alg.seq.arithmetic_mean | Arithmetic means between terms | algebra.sequences | alg.seq.arithmetic |
| alg.seq.geometric_mean | Geometric means between terms | algebra.sequences | alg.seq.geometric |
| alg.seq.fibonacci | Fibonacci and Binet's formula | algebra.sequences | alg.seq.recursive |
| alg.conic.intersection | Line and conic intersections | algebra.conics | alg.conic.parabola |
| alg.conic.translation | Translated conic equations | algebra.conics | alg.conic.circle |
| alg.conic.tangent | Tangent lines to conics | algebra.conics | alg.conic.classify |
| alg.systems.matrix_inverse | Solving systems with inverse matrices | algebra.systems | alg.systems.matrices |
| alg.regression.quadratic | Quadratic regression | algebra.linear | alg.regression.linear |
| alg.regression.power | Power regression | algebra.linear | alg.regression.linear |
| alg.regression.logarithmic | Logarithmic regression | algebra.linear | alg.regression.linear |

## Calculus II: 18 candidates

| id | Label | Subdomain | Prerequisite |
|---|---|---|---|
| calc.series.cauchy_condensation | Cauchy condensation test | calculus.series | calc.series.positive_terms |
| calc.series.raabe | Raabe's test | calculus.series | calc.series.ratio_test |
| calc.series.bertrand | Bertrand's test | calculus.series | calc.series.p_series |
| calc.series.kummer | Kummer's test | calculus.series | calc.series.positive_terms |
| calc.series.gauss | Gauss's test | calculus.series | calc.series.raabe |
| calc.series.abel | Abel's test | calculus.series | calc.series.alternating |
| calc.series.uniform_convergence | Uniform convergence | calculus.series | calc.series.function_series |
| calc.series.euler_transform | Euler's transformation of series | calculus.series | calc.series.alternating |
| calc.series.power_integration | Integrating power series | calculus.series | calc.series.power |
| calc.series.power_differentiation | Differentiating power series | calculus.series | calc.series.power |
| calc.series.lagrange_inversion | Lagrange inversion | calculus.series | calc.series.taylor |
| calc.integral.romberg | Romberg integration | calculus.integrals | calc.integral.trapezoid |
| calc.integral.monte_carlo | Monte Carlo integration | calculus.integrals | calc.integral.definite |
| calc.integral.gauss_legendre | Gauss-Legendre quadrature | calculus.integrals | calc.integral.simpson |
| calc.integral.fresnel | Fresnel integrals | calculus.integrals | calc.integral.improper |
| calc.integral.elliptic | Elliptic integrals | calculus.integrals | calc.integral.trig_substitution |
| calc.integral.wallis | Wallis integrals | calculus.integrals | calc.integral.trig_integrals |
| calc.integral.centroid_curve | Centroid of a curve | calculus.integrals | calc.integral.center_mass |

## Statistics: 14 candidates

| id | Label | Subdomain | Prerequisite |
|---|---|---|---|
| stat.data.weighted_mean | Weighted mean | statistics.central_tendency | stat.data.mean |
| stat.adv.coefficient_variation | Coefficient of variation | statistics.spread | stat.spread.std_dev |
| stat.prob.independence | Independent events | statistics.probability | stat.prob.conditional |
| stat.prob.mutually_exclusive | Mutually exclusive events | statistics.probability | stat.prob.compound |
| stat.dist.negative_binomial | Negative binomial distribution | statistics.distributions | stat.dist.binomial |
| stat.dist.multinomial | Multinomial distribution | statistics.distributions | stat.dist.binomial |
| stat.dist.lognormal | Lognormal distribution | statistics.distributions | stat.dist.normal |
| stat.infer.ci_mean_t | Confidence interval for a mean with t | statistics.inference | stat.dist.student_t, stat.infer.confidence |
| stat.infer.ci_proportion | Confidence interval for a proportion | statistics.inference | stat.infer.proportion |
| stat.infer.power | Power of a hypothesis test | statistics.inference | stat.infer.type_error |
| stat.infer.effect_size | Effect size | statistics.inference | stat.infer.ttest |
| stat.reg.multiple | Multiple regression | statistics.regression | stat.reg.least_squares |
| stat.reg.prediction_interval | Prediction intervals | statistics.regression | stat.reg.least_squares |
| stat.test.kruskal_wallis | Kruskal-Wallis test | statistics.inference | stat.infer.nonparametric |

## Review order

Author and inspect ten first: four algebra, three Calculus II and three statistics. The algebra
sample covers the schema-heavy path, the calculus sample covers numeric quadrature, and the
statistics sample covers a distribution, an interval and a test. If the ten expose a generator or
schema defect, fix the process before authoring the remaining forty.

Every concept must meet the artifact contract in `docs/content-expansion.md`: three question
families aligned to three knowledge points, a lesson over 1 KB with matching `##` sections, three
KP shards with three subgoals each, a course target, and independent fact-based answer tests.
Algebra also needs a solution schema. Calculus II and statistics are university domains and need
none.

The expansion is additive. No candidate is added to `section_alignment_baseline.json`,
`audit_baseline.json` or any other baseline, and no variety floor is lowered.
