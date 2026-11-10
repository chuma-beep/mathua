# Wave 9 candidate ledger

Status: approved and authoring complete on 2026-11-10. The first ten concepts passed review and the remaining forty were authored after that approval. The split finished precalculus at 81, calc1 at 114, calc2 at 97, fractions at 67 and machine learning at 67.

Wave 9 raises the corpus from 1,900 to 1,950 concepts. The approved split closes all remaining bridge domains, closes fractions and extends machine learning by two concepts.

| Domain | Tier | Now | Target | Add | Course |
|---|---|---:|---:|---:|---|
| precalculus | bridge | 58 | 81 | +23 | `pc` |
| calc1 | bridge | 100 | 114 | +14 | `calc1` |
| calc2 | bridge | 89 | 97 | +8 | `calc2` |
| fractions | school | 64 | 67 | +3 | `pa` and grade 6 |
| machine_learning | university | 65 | 67 | +2 | `ml` |
| **Total** | | | | **+50** | |

Wave 10 then closes the remaining 50 exactly: prealgebra +18, geometry +17, graph theory +7 and machine learning +8.

Every id and normalized label below was checked against every file in `data/concepts/` on 2026-11-10. There is no exact id collision and no normalized-label collision. Every prerequisite points to an existing concept or to another Wave 9 candidate. The draft candidate `precalc.poly.synthetic_division` was removed because `alg.poly.synthetic_division` already owns synthetic division.

## Curriculum overlap audit

**Precalculus.** The candidates were checked against algebra, trigonometry, Calculus I and Calculus II. Rejected overlaps include synthetic division (owned by algebra), half-life (owned by algebra exponential decay), vector projection (owned by trigonometry), cross products (owned by linear algebra), power functions (owned by algebra), inverse domain and range (owned by algebra) and the binomial theorem itself (owned by algebra, precalculus and discrete mathematics). The accepted candidates extend the precalculus boundary: factoring over the complex numbers goes beyond algebra factoring, the multinomial theorem goes beyond the binomial theorem, dot-product work is not the line-integral work in Calculus II or III, and parametric and polar conics extend the representation of conics rather than repeating their classification.

**Calculus I.** The candidates were checked against Calculus II, Calculus III, real analysis and numerical analysis. Continuity of compositions and inverses is not the same as continuity on an interval. Asymptotic equivalence is not big-O or little-o. The floor, sign and absolute-value derivatives are new piecewise cases. Elasticity is an economic application of the derivative, not marginal analysis. Shadow, angle, closest-point and cylinder problems are specific related-rates and optimization cases that the generic concepts do not cover. Integrals of floor functions are not integrals of general piecewise functions. Differentiation under the integral sign is a parameter derivative, not the variable-limit FTC.

**Calculus II.** The candidates were checked against Calculus I, Calculus III, real analysis, number theory, probability theory and complex analysis. Cesaro, Borel and Abel summability are not convergence tests. Abel's theorem for power series is not Abel's test. Asymptotic series are not Taylor series or asymptotic equivalence. Euler-Maclaurin is not the trapezoidal rule or Simpson's rule. Tauberian theorems, Dirichlet series and Lambert series are new series structures. None repeats the existing gamma, beta, Fourier or elliptic treatments.

**Fractions.** Adding mixed numbers with carrying is not the existing mixed-number addition because the generated case forces a carry. Ordering mixed numbers is not comparing improper fractions. A percentage of a fraction is not a percentage of a whole number.

**Machine learning.** Log loss is a proper scoring rule, not the cross-entropy training loss or calibration alone. SMOTE is a resampling construction, not the class-imbalance diagnosis or nearest-neighbour distance.

## Precalculus: 23 candidates

| id | Label | Subdomain | Prerequisite | Three families |
|---|---|---|---|---|
| precalc.func.one_to_one | One-to-one functions | precalculus.functions | precalc.func.transformations | horizontal-line test / finite table test / restricted domain |
| precalc.func.inverse_composition | Verifying inverse functions by composition | precalculus.functions | precalc.func.inverse_value, precalc.func.composition_value | f(g(x)) value / g(f(x)) value / verify inverse |
| precalc.func.piecewise_operations | Operating on piecewise functions | precalculus.functions | precalc.func.piecewise, precalc.func.arith_combine | sum value / product value / common domain |
| precalc.analytic.area_polygon | Area of a polygon from coordinates | precalculus.analytic | precalc.analytic.distance | triangle area / quadrilateral area / missing coordinate |
| precalc.poly.complex_factorization | Factoring over the complex numbers | precalculus.polynomials | precalc.poly.complex_conjugate, precalc.poly.rational_roots | linear factors / number of factors / reconstruct polynomial |
| precalc.rational.sign_chart | Sign charts for rational functions | precalculus.rational | precalc.rational.asymptote, precalc.poly.inequality | test-point sign / sign changes / positive interval |
| precalc.rational.asymptote_intersection | Intersections with an oblique asymptote | precalculus.rational | precalc.rational.oblique_asymptote | intersection x / number of intersections / vertical offset |
| precalc.exp.doubling | Doubling time | precalculus.exponential | precalc.exp.growth_model | doubling time / value after k doublings / rate from doubling time |
| precalc.exp.logistic_inverse | The inverse logistic function | precalculus.exponential | precalc.exp.logistic, precalc.log.equation | logit value / inverse probability / restricted domain |
| precalc.seq.summation_properties | Properties of summation notation | precalculus.sequences | precalc.seq.summation | split a sum / constant multiple / combine sums |
| precalc.seq.power_sums | Sums of powers | precalculus.sequences | precalc.seq.summation | sum k / sum k² / sum k³ |
| precalc.conic.parametric | Parametric equations of conics | precalculus.conics | precalc.conic.general, precalc.parametric.equations | parameter to point / eliminate parameter / identify conic |
| precalc.conic.polar | Polar equations of conics | precalculus.conics | precalc.conic.general, precalc.polar.coordinates | conic type / eccentricity / directrix distance |
| precalc.polar.symmetry | Symmetry tests for polar graphs | precalculus.polar | precalc.polar.graph | symmetry code / substitute (-r,θ) / number of symmetries |
| precalc.polar.intersection | Intersections of polar curves | precalculus.polar | precalc.polar.graph | number of intersections / intersection angle / common point |
| precalc.parametric.cycloid | The cycloid | precalculus.parametric | precalc.parametric.equations | point at parameter / period / maximum height |
| precalc.vector.work | Work from the dot product | precalculus.vectors | precalc.vector.dot | work value / cosine of angle / projection magnitude |
| precalc.vector.orthogonal_decomposition | Decomposing a vector into orthogonal components | precalculus.vectors | precalc.vector.dot, precalc.vector.unit | parallel component / perpendicular component / check decomposition |
| precalc.vector.direction_cosines | Direction cosines | precalculus.vectors | precalc.vector.three_d | cosine from components / sum of squares / coordinate from magnitude |
| precalc.systems.nonlinear_inequality | Nonlinear inequalities in two variables | precalculus.systems | precalc.systems.nonlinear, precalc.poly.inequality | test point / boundary type / intersection count |
| precalc.partial_fractions.repeated | Partial fractions with repeated factors | precalculus.analytic | precalc.partial_fractions | repeated-factor coefficient / cleared numerator / evaluate coefficient |
| precalc.partial_fractions.quadratic | Partial fractions with irreducible quadratic factors | precalculus.analytic | precalc.partial_fractions | quadratic coefficient / cleared numerator / discriminant |
| precalc.binomial.multinomial | The multinomial theorem | precalculus.discrete | precalc.binomial_theorem | monomial coefficient / number of terms / evaluate at one |

## Calculus I: 14 candidates

| id | Label | Subdomain | Prerequisite | Three families |
|---|---|---|---|---|
| calc.limit.continuity_composition | Continuity of composite functions | calculus.limits | calc.limit.continuity | limit of composition / continuity test / composition domain |
| calc.limit.continuity_inverse | Continuity of inverse functions | calculus.limits | calc.limit.continuity | inverse value / continuity test / monotonicity condition |
| calc.limit.continuity_extension | Removable discontinuities and continuous extensions | calculus.limits | calc.limit.discontinuity, calc.limit.continuity | limit value / extension value / classify removable |
| calc.limit.asymptotic_equivalence | Asymptotic equivalence | calculus.limits | calc.limit.big_o, calc.limit.little_o | ratio limit / leading coefficient / compare growth |
| calc.deriv.absolute | Derivatives of absolute value functions | calculus.derivatives | calc.deriv.definition, calc.limit.absolute_value | left derivative / right derivative / differentiability |
| calc.deriv.floor | Derivatives of floor and ceiling functions | calculus.derivatives | calc.deriv.piecewise, precalc.func.floor_ceiling | interval derivative / jump location / differentiability |
| calc.deriv.sign | Derivatives of sign functions | calculus.derivatives | calc.deriv.piecewise | interval derivative / jump at zero / product derivative |
| calc.deriv.elasticity | Elasticity of demand | calculus.derivatives | calc.deriv.product_rule, calc.deriv.marginal_analysis | elasticity value / elastic or inelastic / revenue effect |
| calc.deriv.related_rates_shadow | Related rates: shadow length | calculus.advanced | calc.deriv.related_rates | shadow rate / distance rate / ratio |
| calc.deriv.related_rates_angle | Related rates: angle of elevation | calculus.advanced | calc.deriv.related_rates, trig.basics.tan_def | angle rate / tangent value / distance |
| calc.deriv.optimization_distance | Optimization: closest point | calculus.derivatives | calc.deriv.optimization | closest point / squared distance / critical point |
| calc.deriv.optimization_cylinder | Optimization: cylinder dimensions | calculus.derivatives | calc.deriv.optimization | radius / height / volume |
| calc.integral.floor | Integrals of floor functions | calculus.integrals | calc.integral.piecewise, precalc.func.floor_ceiling | integral over interval / sum of floors / jump count |
| calc.integral.leibniz | Differentiation under the integral sign | calculus.integrals | calc.integral.ftc, calc.integral.variable_limits | parameter derivative / value at parameter / integrate after differentiating |

## Calculus II: 8 candidates

| id | Label | Subdomain | Prerequisite | Three families |
|---|---|---|---|---|
| calc.series.cesaro | Cesaro summability | calculus.series | calc.series.concept, calc.series.operations | Cesaro mean / convergence test / Cesaro sum |
| calc.series.abel_theorem | Abel's theorem for power series | calculus.series | calc.series.power, calc.series.abel | boundary limit / series value / condition check |
| calc.series.asymptotic | Asymptotic series | calculus.series | calc.series.taylor, calc.limit.asymptotic_equivalence | leading term / number of terms / error order |
| calc.series.euler_maclaurin | The Euler-Maclaurin formula | calculus.series | calc.series.taylor, calc.integral.numerical | correction term / sum minus integral / error term |
| calc.series.borel | Borel summability | calculus.series | calc.series.power, calc.series.operations | Borel sum / convergence test / coefficient |
| calc.series.taubian | Tauberian theorems | calculus.series | calc.series.cesaro, calc.series.operations | implication test / counterexample / converse condition |
| calc.series.dirichlet_series | Dirichlet series | calculus.series | calc.series.operations, calc.series.power | coefficient / abscissa / value at a point |
| calc.series.lambert | Lambert series | calculus.series | calc.series.power, calc.series.geometric | coefficient / divisor sum / evaluate at a point |

## Fractions: 3 candidates

| id | Label | Subdomain | Prerequisite | Three families |
|---|---|---|---|---|
| frac.mixed.add_carry | Adding mixed numbers with carrying | fractions.mixed | frac.mixed.add, frac.add.diff | carry case / sum / regrouped result |
| frac.compare.mixed_order | Ordering mixed numbers | fractions.equivalence | frac.mixed.convert, frac.ops.compare | larger mixed number / ordered list / benchmark comparison |
| frac.percent.of_fraction | A percentage of a fraction | fractions.equivalence | frac.of_number, pct.ops.of_number | percentage value / remaining fraction / proportional part |

## Machine learning: 2 candidates

| id | Label | Subdomain | Prerequisite | Three families |
|---|---|---|---|---|
| ml.eval.log_loss | Log loss and proper scoring rules | machine_learning.evaluation | ml.loss.cross_entropy, ml.eval.calibration | log loss value / proper score test / compare with accuracy |
| ml.data.smote | Synthetic minority oversampling | machine_learning.data | ml.data.imbalance, ml.knn.nearest_distance | synthetic point / neighbour count / class balance |

## Review order

The first ten cover the riskiest bridge boundary and one representative from each small tail.

1. Precalculus: `precalc.func.one_to_one`, `precalc.poly.complex_factorization`, `precalc.vector.orthogonal_decomposition`.
2. Calculus I: `calc.limit.asymptotic_equivalence`, `calc.deriv.elasticity`, `calc.integral.leibniz`.
3. Calculus II: `calc.series.cesaro`, `calc.series.dirichlet_series`.
4. Fractions: `frac.mixed.add_carry`.
5. Machine learning: `ml.eval.log_loss`.

The first ten are authored, inspected and reported before the remaining forty. If they expose a duplicate, a prerequisite error or a generator-shape problem, the process is corrected before the rest of the wave.

## Artifact contract

Every concept ships with a numeric node in `data/concepts/<domain>.json`, a generator in `internal/generator/<domain>/wave9.go` registered from the package `Register`, an independent fact-based `wave9_test.go`, an authored lesson in `data/lessons/authored/<id>.md` with one H1 and exactly three H2 sections over 1 KB, a three-entry KP shard, a `data/lessons/lessons.json` row and a course target in `pc`, `calc1`, `calc2`, `pa` or `ml`. The three fractions concepts also need authored solution schemas because fractions is a learner domain. No solution schema is required for precalculus, Calculus I, Calculus II or machine learning.

No prerequisite is removed. The 1,500-concept additive baseline stays frozen. Mastery band width stays 3. Commits stay local. Nothing is pushed or deployed.

## Quality gates

- Three section-aligned question families per concept.
- At least five distinct questions per ordinal in twelve draws.
- Design target of twenty distinct per ordinal in 200 draws.
- Minimum accepted pool of ten distinct per ordinal in 60 draws or a written reason.
- Independent fact-based recomputation of every answer.
- Manual inspection of all 150 families with mathematical and presentation defects counted separately.
- No always-yes family and no theorem-name recall. At least one family per theorem-heavy concept must test a condition, counterexample, construction or hypothesis-sensitive computation.
- No baseline exemption, no weakened generator, no lowered variety floor.
