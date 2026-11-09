# Wave 4 proposal

Status: domains and allocation approved on 2026-11-09. Candidate ledger pending review. No concept
has been generated.

Date: 2026-11-09. Current corpus: 1,650 concepts. Frozen additive baseline: 1,500 concepts. Wave 4
target: 1,700 concepts.

This document proposes the domain split, allocations, prerequisite strategy, generator coverage
expectations and quality gates for Wave 4. The exact concept ids and labels belong to a candidate
ledger, which is authored only after this proposal is approved.

## Recommended split

| Domain | Now | Target | Add | Tier | Course |
|---|---:|---:|---:|---|---|
| differential_equations | 46 | 66 | +20 | university | `diffeq` |
| complex_numbers | 51 | 71 | +20 | university | `cx` |
| statistics | 87 | 98 | +10 | bridge | `probstat` |
| **Total** | | | **50** | | |

The split fills the differential-equations target exactly and the complex-numbers target exactly.
Statistics moves from 87 to 97 and leaves one concept for a later wave.

This is a continuous-mathematics and inference wave. It complements Wave 3's number theory, linear
algebra and graph theory instead of repeating them. It avoids precalculus, which carries the
highest duplication risk with algebra and Calculus I. It also avoids the theorem-heavy parts of
real analysis, abstract algebra and topology, where three computational question families are
harder to build and harder to verify independently.

### Why these three domains

- Differential equations has a large target and a well-defined computational spine: first-order
  substitutions, Laplace methods, series solutions, numerical methods, boundary-value problems,
  integral equations and PDE solution techniques.
- Complex numbers has a large target and a clear progression from roots and branches through
  harmonic functions, residue theorems and special functions.
- Statistics has ten remaining inference and regression gaps that do not require duplicating
  probability theory. The proposed topics are tests, intervals, survival, model selection and
  missing-data methods.

### Alternatives, if the recommendation is not approved

1. **Largest-deficit split:** calc3 +26 and real_analysis +24, exactly 50. This fills the two
   largest rows. It is not recommended as the first choice because real analysis is dominated by
   theorem statements and classification questions, which raises the risk of weak generator
   families.
2. **Graph instead of statistics:** differential_equations +20, complex_numbers +20 and
   graph_theory +10. Graph theory would move from 59 to 69 and leave one concept. This is the
   fallback if statistics is judged too close to probability theory.
3. **Discrete instead of statistics:** differential_equations +20, complex_numbers +20 and
   discrete_math +10. Discrete mathematics would move from 64 to 74. This is the fallback if the
   review wants more combinatorics, automata or coding content.

No alternative changes the 1,650-concept corpus until it is approved.

## Candidate topic clusters

These are planning clusters, not approved ids. Every proposed id and normalized label will be
checked against all `data/concepts/` files in the candidate ledger. A cluster that cannot supply
three genuinely distinct families will be replaced before authoring.

### Differential equations: 20

| Cluster | Candidate topics | Expected prerequisite source |
|---|---|---|
| First-order substitutions | Homogeneous substitutions, Clairaut equations, Riccati equations | `ode.first_order.separable`, `ode.first_order.linear`, `ode.first_order.bernoulli` |
| Laplace methods | Inverse transforms, transform of derivatives, periodic forcing, initial and final value theorems | `ode.transforms.laplace`, `ode.transforms.step` |
| Series and special functions | Frobenius method, Bessel functions, Legendre polynomials | `ode.second_order.power_series`, `ode.second_order.cauchy_euler` |
| Numerical methods | Multistep methods, Runge-Kutta-Fehlberg, implicit Euler and stiff equations | `ode.numerical.euler`, `ode.numerical.rk4` |
| Boundary-value methods | Shooting method, finite-difference method | `ode.boundary_value` |
| Integral and functional equations | Fredholm equations, delay differential equations | `ode.integral.volterra`, `ode.basics.existence_uniqueness` |
| Perturbation and PDE methods | Perturbation methods, Duhamel's principle, method of images | `ode.nonlinear.bifurcation`, `ode.pde.heat`, `ode.pde.laplace_eq` |

### Complex numbers: 20

| Cluster | Candidate topics | Expected prerequisite source |
|---|---|---|
| Roots and branches | Roots of unity, principal value and branch cuts, complex powers | `complex.adv.roots`, `complex.adv.logarithm`, `complex.adv.exponential` |
| Zeros and harmonic functions | Zeros of complex trigonometric functions, harmonic conjugates, Poisson kernel | `complex.adv.trig_exponential`, `complex.adv.harmonic`, `complex.integral.cauchy_formula` |
| Theorems of analytic functions | Rouché, open mapping, Casorati-Weierstrass, Picard | `complex.adv.analytic`, `complex.adv.maximum_modulus`, `complex.adv.residue` |
| Continuation and surfaces | Analytic continuation, monodromy, Riemann surfaces | `complex.adv.analytic`, `complex.adv.winding`, `complex.topology.riemann_sphere` |
| Special functions | Elliptic functions, modular forms, gamma function | `complex.series.laurent`, `complex.integral.cauchy_formula`, `complex.adv.liouville` |
| Factorisation and mapping | Weierstrass factorisation, Mittag-Leffler, Schwarz-Christoffel, Jensen's formula | `complex.series.power`, `complex.adv.residue`, `complex.adv.conformal` |

### Statistics: 10

| Cluster | Candidate topics | Expected prerequisite source |
|---|---|---|
| Resampling | Bootstrap confidence intervals, permutation tests | `stat.infer.confidence`, `stat.infer.hypothesis`, `stat.infer.pvalue` |
| Rank tests | Wilcoxon signed-rank, Mann-Whitney U | `stat.infer.nonparametric`, `stat.infer.paired`, `stat.infer.ttest_two` |
| Survival and Bayesian inference | Kaplan-Meier curve, Bayesian credible intervals | `stat.infer.confidence`, `stat.dist.exponential`, `stat.prob.bayes` |
| Model selection and regression | Information criteria, quantile regression, Poisson regression | `stat.reg.multiple`, `stat.reg.r2`, `stat.reg.logistic`, `stat.dist.poisson` |
| Missing data | Imputation and missing-data mechanisms | `stat.data.study_design`, `stat.reg.multiple` |

## Prerequisite strategy

1. Every new concept must have at least one existing prerequisite in its own domain, or a
   prerequisite on an earlier concept in the same wave.
2. No new cross-domain prerequisite is allowed unless the candidate ledger records why the
   dependency is essential. The default is to stay inside the domain.
3. No existing prerequisite edge is removed. The only permitted existing-concept change is the
   addition of a prerequisite that gives an otherwise isolated new concept a real follow-on, and
   only when the existing lesson already teaches the dependency.
4. Every candidate is checked for cycles, missing targets and normalized-label collisions before
   authoring.
5. The candidate ledger will show the depth and band effect of every new prerequisite edge. A new
   singleton unit must be identified before acceptance and either justified or corrected by an
   existing follow-on edge.
6. The mastery band width remains 3. The algorithm is not changed to hide a singleton.

## Generator coverage expectations

Each concept must have three section-aligned question families, one per knowledge point.

| Gate | Requirement |
|---|---|
| Hard floor | At least 5 distinct questions per ordinal in 12 draws |
| Wave 4 design target | At least 20 distinct questions per ordinal in 200 draws |
| Minimum accepted pool | At least 10 distinct questions per ordinal in 60 draws, or a written reason in the candidate ledger |
| Answer verification | Every numeric answer recomputed from published `Facts` by a test independent of the generator answer |
| Section alignment | Different question at each ordinal under the same seed |
| Finite pools | A family with a closed parameter set must report its size and cannot be described as parametric |
| Manual inspection | All 150 families read for answer correctness, explanation correctness and presentation |

The Wave 3 review found that finite pools are acceptable only when their size is explicit. Wave 4
will therefore measure every family rather than rely on the 12-draw floor.

## Quality gates

The wave is not complete until all of the following pass.

- `make validate`.
- `go test ./...`.
- `git diff --check`.
- `go run scripts/validate_graph.go`.
- `go run ./cmd/validate-mastery-units`.
- `python3 scripts/audit_lessons.py`.
- `python3 scripts/audit_solutions.py`.
- `python3 scripts/counts.py --check`.
- `python3 scripts/check_additive.py`.
- `python3 scripts/check_docs.py`.
- `python3 scripts/check_ste.py`.
- `node web/next-app/scripts/build-docs.mjs --check`.

The additive checks must preserve the 1,500-concept baseline. No concept may be removed, no lesson
mapping may be changed, no prerequisite may be removed and no variety floor may be lowered. No
entry may be added to a quality exemption baseline.

Every new concept also needs:

- a concept node with numeric grading and a mastery threshold;
- a generator registered from its domain package;
- an independent fact-based test;
- a lesson over 1 KB with exactly three H2 sections;
- a three-entry KP shard whose sections resolve;
- a `lessons.json` row;
- a target in its course.

The three proposed domains are not learner domains, so no solution schema is required.

## Approval checkpoints

1. Approve the domains and allocation in this proposal.
2. Approve the candidate ledger with exact ids, labels, subdomains, prerequisite edges and
   generator family plan.
3. Author the first ten concepts and run the full domain, lesson and variety checks.
4. Inspect those ten manually. Fix the process before the remaining forty if a defect appears.
5. Author the remaining forty.
6. Regenerate mastery, graph, flow, docs and variety artifacts.
7. Inspect all 150 families, record mathematical and presentation defects separately and commit
   locally.
8. Do not push or deploy. Keep the unrelated server files unstaged.

No concept generation begins before checkpoint 1 and checkpoint 2 are approved.
