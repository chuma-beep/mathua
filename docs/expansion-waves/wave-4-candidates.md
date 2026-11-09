# Wave 4 candidate ledger

Status: approved for authoring. Authored and closed at 1,700 concepts.

Date: 2026-11-09. Current corpus: 1,650 concepts. Wave 4 target: 1,700 concepts.

Every id and normalized label below was checked against every file in `data/concepts/` on
2026-11-09 with the same normalization rule as `internal/concepts/labels.go`. There is no exact id
collision and no normalized-label collision. The complex gamma candidate was renamed from "The
gamma function" to "The complex gamma function" because `calc.series.gamma` already owns the real
gamma-series label.

A prerequisite marked `*` is another candidate in this ledger. All other prerequisites exist in
the current corpus. The candidate ledger is checkpoint 2: authoring starts only after this file
is approved.

## Allocation

| Domain | Now | Target | Add | Course |
|---|---:|---:|---:|---|
| differential_equations | 46 | 66 | +20 | `diffeq` |
| complex_numbers | 51 | 71 | +20 | `cx` |
| statistics | 87 | 97 | +10 | `probstat` |
| **Total** | | | **50** | |

The wave fills differential equations and complex numbers exactly. Statistics moves to 97 and
leaves one concept for a later wave. The three domains are not learner domains, so no authored
solution schema is required.

## Differential equations: 20 candidates

| id | Label | Subdomain | Prerequisite |
|---|---|---|---|
| ode.first_order.homogeneous_substitution | Homogeneous substitutions for first-order ODEs | differential_equations.first_order | ode.first_order.separable |
| ode.first_order.clairaut | Clairaut equations | differential_equations.first_order | ode.first_order.linear |
| ode.first_order.riccati | Riccati equations | differential_equations.first_order | ode.first_order.bernoulli |
| ode.transforms.inverse_laplace | Inverse Laplace transforms | differential_equations.transforms | ode.transforms.laplace |
| ode.transforms.derivative | Laplace transform of derivatives | differential_equations.transforms | ode.transforms.laplace |
| ode.transforms.periodic | Laplace transforms of periodic functions | differential_equations.transforms | ode.transforms.laplace, ode.transforms.step |
| ode.transforms.initial_final | Initial and final value theorems | differential_equations.transforms | ode.transforms.laplace |
| ode.second_order.frobenius | The Frobenius method | differential_equations.second_order | ode.second_order.power_series |
| ode.second_order.bessel | Bessel functions | differential_equations.second_order | ode.second_order.power_series |
| ode.second_order.legendre | Legendre polynomials | differential_equations.second_order | ode.second_order.power_series |
| ode.numerical.multistep | Multistep methods | differential_equations.numerical | ode.numerical.euler |
| ode.numerical.rkf | Runge-Kutta-Fehlberg methods | differential_equations.numerical | ode.numerical.rk4 |
| ode.numerical.implicit | Implicit Euler and stiff equations | differential_equations.numerical | ode.numerical.euler |
| ode.boundary.shooting | The shooting method | differential_equations.boundary | ode.boundary_value |
| ode.boundary.finite_difference | Finite-difference methods for boundary-value problems | differential_equations.boundary | ode.boundary_value |
| ode.integral.fredholm | Fredholm integral equations | differential_equations.integral | ode.integral.volterra |
| ode.basics.delay | Delay differential equations | differential_equations.basics | ode.basics.existence_uniqueness |
| ode.nonlinear.perturbation | Perturbation methods | differential_equations.nonlinear | ode.nonlinear.bifurcation |
| ode.pde.duhamel | Duhamel's principle | differential_equations.pde | ode.pde.heat, ode.transforms.step |
| ode.pde.images | The method of images | differential_equations.pde | ode.pde.laplace_eq |

### Differential-equation generator families

| id | KP1 | KP2 | KP3 |
|---|---|---|---|
| homogeneous_substitution | Identify the substitution | Transform the equation | Solve the separated equation |
| clairaut | Differentiate the Clairaut form | Find the general solution | Verify the singular solution |
| riccati | Identify a particular solution | Apply the Riccati substitution | Solve the resulting linear equation |
| inverse_laplace | Match a transform pair | Invert a linear combination | Invert a shifted transform |
| laplace_derivative | Transform a derivative | Insert the initial condition | Solve the transformed equation |
| periodic | Compute a periodic transform | Identify the period | Invert the periodic response |
| initial_final | Read the initial value | Read the final value | Check the limit conditions |
| frobenius | Find the indicial roots | Build the recurrence | State the solution near the singularity |
| bessel | Identify the Bessel order | Use the recurrence | Evaluate a Bessel value |
| legendre | Evaluate a Legendre polynomial | Use orthogonality | Apply the recurrence |
| multistep | Apply the multistep formula | Determine the order | Check zero stability |
| rkf | Compute the embedded pair | Estimate the local error | Choose the next step |
| implicit_euler | Apply the implicit step | Solve the implicit equation | Test stiffness |
| shooting | Choose an initial slope | Integrate to the boundary | Update the slope |
| finite_difference | Build the grid equations | Insert the boundary values | Solve the linear system |
| fredholm | Identify the kernel | Reduce a separable kernel | Evaluate the resolvent |
| delay | Evaluate the history function | Advance one delay step | Find the characteristic roots |
| perturbation | Find the leading term | Compute the first correction | Check the validity range |
| duhamel | Identify the impulse response | Write the convolution | Evaluate the response |
| images | Place the image source | Enforce the boundary condition | Evaluate the potential |

## Complex numbers: 20 candidates

| id | Label | Subdomain | Prerequisite |
|---|---|---|---|
| complex.adv.roots_unity | Roots of unity | complex_numbers.advanced | complex.adv.roots |
| complex.adv.branches | Principal values and branch cuts | complex_numbers.advanced | complex.adv.logarithm |
| complex.adv.complex_powers | Complex exponentiation | complex_numbers.advanced | complex.adv.exponential, complex.adv.logarithm |
| complex.adv.trig_zeros | Zeros of complex trigonometric functions | complex_numbers.advanced | complex.adv.trig_exponential |
| complex.adv.harmonic_conjugate | Harmonic conjugates | complex_numbers.advanced | complex.adv.harmonic, complex.adv.cauchy_riemann |
| complex.adv.poisson_kernel | The Poisson kernel | complex_numbers.advanced | complex.adv.harmonic, complex.integral.cauchy_formula |
| complex.adv.rouche | Rouché's theorem | complex_numbers.advanced | complex.adv.argument_principle, complex.adv.residue |
| complex.adv.open_mapping | The open mapping theorem | complex_numbers.advanced | complex.adv.analytic, complex.adv.maximum_modulus |
| complex.adv.casorati | The Casorati-Weierstrass theorem | complex_numbers.advanced | complex.series.laurent, complex.adv.residue |
| complex.adv.picard | Picard's theorem | complex_numbers.advanced | complex.adv.liouville, complex.adv.maximum_modulus |
| complex.adv.analytic_continuation | Analytic continuation | complex_numbers.advanced | complex.adv.analytic, complex.series.power |
| complex.adv.monodromy | Monodromy | complex_numbers.advanced | complex.adv.analytic_continuation*, complex.adv.winding |
| complex.adv.riemann_surface | Riemann surfaces | complex_numbers.topology | complex.adv.analytic_continuation*, complex.topology.riemann_sphere |
| complex.adv.elliptic | Elliptic functions | complex_numbers.advanced | complex.series.laurent, complex.adv.liouville |
| complex.adv.modular_forms | Modular forms | complex_numbers.advanced | complex.adv.elliptic* |
| complex.adv.gamma | The complex gamma function | complex_numbers.advanced | complex.integral.cauchy_formula, complex.adv.analytic |
| complex.adv.weierstrass_factorization | The Weierstrass factorization theorem | complex_numbers.advanced | complex.adv.analytic, complex.series.power |
| complex.adv.mittag_leffler | The Mittag-Leffler theorem | complex_numbers.advanced | complex.adv.residue, complex.series.laurent |
| complex.adv.schwarz_christoffel | The Schwarz-Christoffel mapping | complex_numbers.advanced | complex.adv.conformal, complex.adv.mobius |
| complex.adv.jensen | Jensen's formula | complex_numbers.advanced | complex.adv.harmonic, complex.adv.maximum_modulus |

### Complex-number generator families

| id | KP1 | KP2 | KP3 |
|---|---|---|---|
| roots_unity | List the roots | Sum the roots | Identify primitive roots |
| branches | Choose the principal value | Locate the branch cut | Test continuity across the cut |
| complex_powers | Evaluate a complex power | Apply the principal branch | Use the exponential identity |
| trig_zeros | Solve a trigonometric zero | State the period | Count zeros in a range |
| harmonic_conjugate | Find the conjugate | Verify the Cauchy-Riemann equations | Check harmonicity |
| poisson_kernel | Evaluate the kernel | Integrate against a boundary value | Recover the boundary limit |
| rouche | Compare two functions on a contour | Count zeros | Apply Rouché's conclusion |
| open_mapping | Identify an open image | Test the hypothesis | Apply the theorem |
| casorati | Identify an essential singularity | Describe the local image | Apply Casorati-Weierstrass |
| picard | Identify the exceptional value | Count omitted values | Apply Picard's theorem |
| analytic_continuation | Continue a power series | Find the new radius | Identify the obstruction |
| monodromy | Follow a loop | Compute the change of branch | Describe the monodromy |
| riemann_surface | Identify the sheets | Locate a branch point | Compute the local model |
| elliptic | Find the periods | Locate the poles | Use the addition relation |
| modular_forms | Identify the weight | Apply the transformation | Read the q-expansion |
| gamma | Evaluate the gamma function | Use the reflection formula | Locate the poles |
| weierstrass_factorization | List the zeros | Determine the order | Build the product |
| mittag_leffler | List the poles | Find the principal parts | Build the series |
| schwarz_christoffel | Identify the vertices | Apply the exponent formula | Check the mapping |
| jensen | Evaluate the integral | Count zeros inside | Apply Jensen's formula |

## Statistics: 10 candidates

| id | Label | Subdomain | Prerequisite |
|---|---|---|---|
| stat.infer.bootstrap | Bootstrap confidence intervals | statistics.inference | stat.infer.confidence, stat.sampling.se |
| stat.test.permutation | Permutation tests | statistics.inference | stat.infer.hypothesis, stat.infer.pvalue |
| stat.test.wilcoxon | The Wilcoxon signed-rank test | statistics.inference | stat.infer.nonparametric, stat.infer.paired |
| stat.test.mann_whitney | The Mann-Whitney U test | statistics.inference | stat.infer.nonparametric, stat.infer.ttest_two |
| stat.infer.kaplan_meier | The Kaplan-Meier estimator | statistics.inference | stat.infer.confidence, stat.dist.exponential |
| stat.infer.credible_interval | Bayesian credible intervals | statistics.inference | stat.prob.bayes, stat.infer.confidence |
| stat.reg.information_criteria | Information criteria for model selection | statistics.regression | stat.reg.multiple, stat.reg.r2 |
| stat.reg.quantile | Quantile regression | statistics.regression | stat.reg.least_squares, stat.reg.multiple |
| stat.reg.poisson | Poisson regression | statistics.regression | stat.reg.logistic, stat.dist.poisson |
| stat.data.missing | Missing-data mechanisms and imputation | statistics.data_collection | stat.data.study_design, stat.reg.multiple |

### Statistics generator families

| id | KP1 | KP2 | KP3 |
|---|---|---|---|
| bootstrap | Resample a statistic | Build a percentile interval | Estimate the standard error |
| permutation | Compute the observed statistic | Permute the labels | Read the permutation p-value |
| wilcoxon | Rank the absolute differences | Sum the signed ranks | Read the test conclusion |
| mann_whitney | Rank the pooled sample | Compute U | Read the test conclusion |
| kaplan_meier | Compute a survival step | Account for censoring | Find the median survival |
| credible_interval | Form the posterior | Find the interval | Interpret the prior effect |
| information_criteria | Compute AIC | Compute BIC | Compare two models |
| quantile_regression | Evaluate the quantile loss | Estimate a coefficient | Predict a quantile |
| poisson_regression | Apply the log link | Compute a rate | Interpret a coefficient |
| missing_data | Identify the missingness mechanism | Impute a value | Measure the bias |

## Review order

Author and inspect the first ten before the remaining forty:

1. Four differential-equation concepts: homogeneous substitution, inverse Laplace, Frobenius method and shooting method.
2. Three complex concepts: roots of unity, harmonic conjugates and Rouché's theorem.
3. Three statistics concepts: bootstrap intervals, Wilcoxon signed-rank and Kaplan-Meier.

The sample covers the ODE first-order, transform, series and boundary paths; the complex roots,
harmonic and theorem paths; and the statistics resampling, rank and survival paths. If the ten
expose a generator or schema problem, fix the process before the remaining forty.

## Artifact contract

Every concept must ship with:

- a node in `data/concepts/<domain>.json` with numeric grading and a mastery threshold;
- a generator in `internal/generator/<domain>/wave4.go`, registered from that package's `Register`;
- a `wave4_test.go` that recomputes every numeric answer from published facts;
- a lesson in `data/lessons/authored/<id>.md` with one H1 and exactly three H2 sections over 1 KB;
- a three-entry KP shard in `data/lessons/kp/<id>.json` whose sections resolve;
- a `lessons.json` row;
- a target in `diffeq`, `cx` or `probstat`.

No solution schema is required because none of the three domains is a learner domain.

## Coverage and quality expectations

- Three section-aligned families per concept.
- Hard floor: at least 5 distinct questions per ordinal in 12 draws.
- Wave 4 target: at least 20 distinct questions per ordinal in 200 draws.
- Minimum accepted pool: at least 10 distinct questions per ordinal in 60 draws, or a written
  reason in the final review.
- Every family publishes facts and is independently recomputed.
- No existing prerequisite is removed. The only existing-concept edge allowed is an added
  prerequisite that gives an isolated new concept a real follow-on.
- The 1,500-concept additive baseline stays frozen. No variety floor is lowered and no exemption
  is added.
- The mastery band width stays 3. A new singleton is identified before acceptance.
- Manual inspection covers all 150 families, with mathematical and presentation defects counted
  separately.
- `make validate`, `go test ./...`, `git diff --check` and the regenerated documentation bundle
  must pass before the wave closes.
- The commit is local. No push or deployment. The unrelated server files remain unstaged.

## Approval

This ledger is the second checkpoint. Authoring does not begin until the ids, labels, subdomains,
prerequisite edges and generator plans above are approved. Any replacement after approval must be
recorded before the replacement concept is added.
