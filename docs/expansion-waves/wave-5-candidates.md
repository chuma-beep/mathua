# Wave 5 candidate ledger

Status: approved and authoring complete on 2026-11-09. The first ten passed the independent audit,
and the remaining forty were authored after that approval. The final corpus is 1,750 concepts.

Date: 2026-11-09. Starting corpus: 1,700 concepts. Wave 5 target: 1,750 concepts. Frozen additive
baseline: 1,500 concepts.

This ledger fixes the 50 candidates and their provisional prerequisite edges before authoring. Every
id and normalized label was checked against every file in `data/concepts/` with the normalization
rule in `internal/concepts/labels.go`: there is no exact id collision and no normalized-label
collision.

The split is Calculus III plus 26 and real analysis plus 24. This file is the Wave 5 ledger, so the
table below is the authoritative allocation.

| Domain | Now | Target | Add | Tier | Course |
|---|---:|---:|---:|---|---|
| calc3 | 71 | 97 | +26 | bridge | `calc3` |
| real_analysis | 65 | 89 | +24 | university | `ra` |
| **Total** | | | **50** | | |

The two target rows close exactly. The wave leaves 250 concepts for Waves 6 through 10.

## Curriculum overlap audit

The five checks in `docs/expansion-waves/wave-5-proposal.md` were applied to every candidate.

1. **Duplicate identity.** No exact id or normalized label collides with the existing corpus. The
   check covered all 23 domains, not only Calculus III and real analysis.
2. **Genuine advancement.** Every candidate is beyond the existing Calculus II and Calculus III
   spine. The ledger does not repeat Cauchy products, Dirichlet's test, Fourier series, the gamma
   function, the beta function, integration by parts, Green's theorem, Stokes' theorem, the
   divergence theorem, Jacobians, Lagrange multipliers or the Riemann integral.
3. **Prerequisite honesty.** Every prerequisite is an existing concept or an earlier concept in
   this wave. The Fourier convergence candidate depends on `calc.series.fourier` and
   `ra.deriv.uniform_conv`, not on Lebesgue integration. The multivariable Fourier transform depends
   on `calc.seq.function_sequences`, not on the real-analysis measure sequence.
4. **Independent mastery objective.** Each row states three families that measure a distinct
   computation, construction or comparison. No row is only a theorem-name recall.
5. **Family fit.** The three families are separated by the mathematical object they manipulate, not
   by wording or coefficient changes. Finite counting measures, explicit polynomials, small metric
   spaces and exact matrices are preferred where the theorem is abstract.

### Cross-domain decisions

- `calc.multi.fundamental_solution` is about the fundamental solution of the Laplacian in several
  variables. It does not repeat `ode.green.function`, which is the one-dimensional Green function.
- `calc.multi.fourier` and `calc.multi.convolution` are several-variable transforms. They do not
  repeat `ode.transforms.fourier` or `ode.transforms.convolution`.
- `ra.seq.fourier_convergence` starts from the existing `calc.series.fourier` and studies
  convergence. It does not repeat the Calculus II Fourier-series concept.
- `ra.metric.lp_space`, `ra.metric.banach` and `ra.metric.hilbert` are finite-dimensional and
  sequence-space examples. They do not repeat the existing metric-space concepts.
- `ra.measure.*` extends the existing measure-zero and Cantor-set material into measure, measurable
  functions, integration and convergence theorems.

### Candidates rejected by the audit

These topics were considered and dropped before the ledger was fixed:

- The method of images, because `ode.pde.images` already owns that label and the idea.
- The Poisson kernel, because `complex.adv.poisson_kernel` was added in Wave 4.
- Fourier series, because `calc.series.fourier` already owns the label.
- Cauchy product and Dirichlet's test, because `calc.series.cauchy_product` and
  `calc.series.dirichlet_test` already own them.
- The gamma and beta functions, because `calc.series.gamma` and `calc.integral.beta` already own
  them.
- Green's functions in several variables, because the overlap with `ode.green.function` is too
  close without a separate PDE course.
- Radon-Nikodym, Egorov and Lusin, because the remaining 24 slots are better spent on the
  measure-integration-functional sequence that supports exact finite examples.

## Calculus III: 26 candidates

| id | Label | Subdomain | Prerequisite |
|---|---|---|---|
| calc.vec.orientation | Orientation of surfaces | calculus.vector | calc.vec.parametric_surface |
| calc.vec.oriented_flux | Oriented surface flux | calculus.vector | calc.vec.flux_general, calc.vec.orientation |
| calc.vec.general_stokes | Generalized Stokes' theorem | calculus.vector | calc.vec.stokes, calc.vec.exterior_derivative |
| calc.vec.differential_form | Differential forms | calculus.vector | calc.vec.identities |
| calc.vec.wedge_product | The wedge product | calculus.vector | calc.vec.differential_form |
| calc.vec.exterior_derivative | The exterior derivative | calculus.vector | calc.vec.differential_form |
| calc.vec.pullback | Pullback of a differential form | calculus.vector | calc.vec.differential_form |
| calc.vec.closed_exact | Closed and exact forms | calculus.vector | calc.vec.exterior_derivative |
| calc.vec.vector_potential | Vector potentials | calculus.vector | calc.vec.conservative_test |
| calc.vec.helmholtz | Helmholtz decomposition | calculus.vector | calc.vec.vector_potential, calc.vec.divergence |
| calc.vec.gauge | Gauge transformations | calculus.vector | calc.vec.vector_potential |
| calc.multi.harmonic | Harmonic functions | calculus.multivariable | calc.deriv.partial_higher |
| calc.multi.mean_value | The mean value property | calculus.multivariable | calc.multi.harmonic |
| calc.multi.maximum_principle | The maximum principle | calculus.multivariable | calc.multi.harmonic |
| calc.multi.poisson | The Poisson equation | calculus.multivariable | calc.multi.harmonic |
| calc.multi.dirichlet | The Dirichlet problem | calculus.multivariable | calc.multi.poisson |
| calc.multi.neumann | The Neumann problem | calculus.multivariable | calc.multi.poisson |
| calc.vec.greens_identity | Green's identities | calculus.vector | calc.vec.greens_theorem, calc.vec.divergence_theorem |
| calc.multi.fundamental_solution | The fundamental solution | calculus.multivariable | calc.multi.poisson |
| calc.vec.euler_lagrange | The Euler-Lagrange equation | calculus.advanced | calc.multi.lagrange, calc.vec.arc_length |
| calc.vec.first_variation | The first variation | calculus.advanced | calc.vec.euler_lagrange |
| calc.vec.geodesic | Geodesics | calculus.advanced | calc.vec.first_variation |
| calc.vec.minimal_surface | Minimal surfaces | calculus.advanced | calc.vec.first_variation |
| calc.multi.surface_measure | Surface measure | calculus.multivariable | calc.multi.transform_jacobian, calc.vec.parametric_surface |
| calc.multi.fourier | Multivariable Fourier transform | calculus.multivariable | calc.seq.function_sequences |
| calc.multi.convolution | Multivariable convolution | calculus.multivariable | calc.multi.fourier |

### Calculus III family plan

| id | KP1 | KP2 | KP3 |
|---|---|---|---|
| calc.vec.orientation | Sign of the normal from a parametrisation | Orientation induced on the boundary | Effect of reversing the parametrisation |
| calc.vec.oriented_flux | Flux with the chosen normal | Flux with the opposite normal | Outward flux through a closed surface |
| calc.vec.general_stokes | Boundary orientation of a square | Evaluate a 1-form on the oriented boundary | Verify the theorem on a square |
| calc.vec.differential_form | Evaluate a 1-form on a vector | Evaluate a 2-form on a pair of vectors | Degree of a form |
| calc.vec.wedge_product | Coefficient of dx wedge dy | Antisymmetry sign | Wedge of two explicit 1-forms |
| calc.vec.exterior_derivative | Exterior derivative of a 0-form | Exterior derivative of a 1-form | Check d squared equals zero |
| calc.vec.pullback | Pullback of dx | Pullback of dx wedge dy | Jacobian factor in a pullback |
| calc.vec.closed_exact | Closed-form test | Exact-form potential | Closed but not exact on the punctured plane |
| calc.vec.vector_potential | Find a potential for a simple field | Verify the curl | Gauge shift of the potential |
| calc.vec.helmholtz | Divergence-free part | Curl-free part | Component count in a finite example |
| calc.vec.gauge | Add a gradient | Show the field is unchanged | Evaluate the gauge function |
| calc.multi.harmonic | Laplacian zero test | Coefficient that makes a polynomial harmonic | Sum of two harmonic functions |
| calc.multi.mean_value | Average on a circle | Value at the centre | Average over a boundary |
| calc.multi.maximum_principle | Boundary maximum | Interior maximum test | Constant function case |
| calc.multi.poisson | Solve a polynomial Poisson problem | Particular solution | Compatibility condition |
| calc.multi.dirichlet | Boundary data | Value at an interior point | Uniqueness check |
| calc.multi.neumann | Normal derivative | Compatibility integral | Constant shift |
| calc.vec.greens_identity | Evaluate an identity | Boundary term | Symmetry of the identity |
| calc.multi.fundamental_solution | Laplacian of log r | Normalising coefficient | Delta source in a finite model |
| calc.vec.euler_lagrange | Functional derivative | Solve a simple Euler-Lagrange equation | Geodesic check |
| calc.vec.first_variation | Variation of length | Stationary path | Boundary term |
| calc.vec.geodesic | Shortest path in the plane | Great circle on a sphere | Geodesic equation |
| calc.vec.minimal_surface | Area functional | Minimal surface equation | Catenoid check |
| calc.multi.surface_measure | Area element | Surface area of a graph | Scaling under a parametrisation |
| calc.multi.fourier | Gaussian transform | Scaling rule | Inverse transform at a point |
| calc.multi.convolution | Convolution of two boxes | Fourier product | Value at a point |

## Real analysis: 24 candidates

| id | Label | Subdomain | Prerequisite |
|---|---|---|---|
| ra.seq.limsup | The limit superior of a sequence | real_analysis.seq | ra.seq.subsequence |
| ra.seq.liminf | The limit inferior of a sequence | real_analysis.seq | ra.seq.subsequence |
| ra.seq.cluster | Cluster points and subsequential limits | real_analysis.seq | ra.seq.bolzano_weierstrass |
| ra.seq.double_series | Double series | real_analysis.seq | ra.seq.series |
| ra.seq.fourier_convergence | Convergence of Fourier series | real_analysis.seq | calc.series.fourier, ra.deriv.uniform_conv |
| ra.measure.lebesgue | Lebesgue measure | real_analysis.measure | ra.measure.null |
| ra.measure.measurable_sets | Measurable sets | real_analysis.measure | ra.measure.lebesgue |
| ra.measure.measurable_functions | Measurable functions | real_analysis.measure | ra.measure.measurable_sets |
| ra.measure.simple_integral | The integral of a simple function | real_analysis.measure | ra.measure.measurable_functions |
| ra.measure.lebesgue_integral | The Lebesgue integral | real_analysis.measure | ra.measure.simple_integral |
| ra.measure.monotone_convergence | The monotone convergence theorem | real_analysis.measure | ra.measure.lebesgue_integral |
| ra.measure.fatou | Fatou's lemma | real_analysis.measure | ra.measure.monotone_convergence |
| ra.measure.dominated_convergence | The dominated convergence theorem | real_analysis.measure | ra.measure.lebesgue_integral |
| ra.measure.product | Product measures | real_analysis.measure | ra.measure.lebesgue |
| ra.measure.fubini_tonelli | The Fubini-Tonelli theorem | real_analysis.measure | ra.measure.product |
| ra.metric.lp_space | Lp spaces | real_analysis.metric | ra.measure.lebesgue_integral, ra.metric.complete |
| ra.metric.holder | Holder's inequality | real_analysis.metric | ra.metric.lp_space |
| ra.metric.minkowski | Minkowski's inequality | real_analysis.metric | ra.metric.lp_space |
| ra.metric.banach | Banach spaces | real_analysis.metric | ra.metric.lp_space |
| ra.metric.hilbert | Hilbert spaces | real_analysis.metric | ra.metric.banach |
| ra.metric.bounded_operator | Bounded linear operators | real_analysis.metric | ra.metric.banach |
| ra.metric.riesz_representation | The Riesz representation theorem | real_analysis.metric | ra.metric.hilbert |
| ra.deriv.absolute_continuity | Absolutely continuous functions | real_analysis.deriv | ra.deriv.bounded_variation |
| ra.deriv.lebesgue_differentiation | The Lebesgue differentiation theorem | real_analysis.deriv | ra.measure.lebesgue_integral |

### Real analysis family plan

| id | KP1 | KP2 | KP3 |
|---|---|---|---|
| ra.seq.limsup | Compute the limsup of an explicit sequence | Compare limsup with the limit | Limsup of a sum |
| ra.seq.liminf | Compute the liminf of an explicit sequence | Compare liminf with the limit | Liminf of a product |
| ra.seq.cluster | Identify a cluster point | Count the cluster points | Limit of a subsequence |
| ra.seq.double_series | Sum over a finite grid | Convergence by rows | Absolute convergence |
| ra.seq.fourier_convergence | Coefficient decay | Dirichlet kernel value | Pointwise value at a point |
| ra.measure.lebesgue | Measure of an interval | Measure of a union | Measure-zero example |
| ra.measure.measurable_sets | Sigma-algebra closure | Preimage of a set | Countable union |
| ra.measure.measurable_functions | Indicator function | Composition with a continuous function | Supremum of measurable functions |
| ra.measure.simple_integral | Sum of c_i m(A_i) | Integral under a counting measure | Integral of an indicator |
| ra.measure.lebesgue_integral | Integral of a step function | Compare with the Riemann integral | Counting measure example |
| ra.measure.monotone_convergence | Increasing sequence of integrals | Indicator sequence | Finite measure case |
| ra.measure.fatou | Lower bound | Strict inequality example | Nonnegative sequence |
| ra.measure.dominated_convergence | Limit of integrals | Dominating function | Indicator sequence |
| ra.measure.product | Rectangle measure | Product sigma-algebra | Area of a product set |
| ra.measure.fubini_tonelli | Iterated integral | Swap the order | Finite grid check |
| ra.metric.lp_space | Compute an Lp norm | The p = 2 case | Finite sequence space |
| ra.metric.holder | Verify the inequality | Equality case | Conjugate exponents |
| ra.metric.minkowski | Triangle inequality | Norm of a sum | Equality case |
| ra.metric.banach | Completeness test | Cauchy sequence | Finite-dimensional example |
| ra.metric.hilbert | Inner product | Orthogonality | Finite Parseval identity |
| ra.metric.bounded_operator | Operator norm | Boundedness test | Matrix example |
| ra.metric.riesz_representation | Represent a functional | Inner-product form | Norm of the representative |
| ra.deriv.absolute_continuity | Fundamental theorem check | Cantor function | Epsilon-delta estimate |
| ra.deriv.lebesgue_differentiation | Average over an interval | Derivative of an integral | Lebesgue point |

## Prerequisite and band audit

- No existing prerequisite is removed or changed in this ledger.
- Every new edge points to an existing concept or to an earlier concept in the same wave.
- The only cross-domain prerequisite is `ra.seq.fourier_convergence` to `calc.series.fourier`. It is
  essential because the candidate studies convergence of the series that concept defines.
- The graph must be acyclic. The same-wave edges are ordered from forms to Stokes, from harmonic
  functions to Poisson problems, from Euler-Lagrange to geodesics and from measure to integration
  to Lp spaces.
- The mastery band width remains 3. The depth and band effect, including every new singleton unit,
  will be measured when the nodes are added. A new singleton is identified before acceptance and
  either justified or given a real follow-on edge.
- No existing concept is rewired to make the allocation work.

## Risks and mitigations

| Risk | Mitigation |
|---|---|
| Real-analysis theorem recognition | Every family asks for an explicit value, set, sequence or finite counterexample |
| Abstract integration without computation | Use simple functions, counting measures and step functions before general theorems |
| Calc III differential forms becoming notation-only | Ask for coefficients, signs, pullbacks and boundary evaluations |
| Cross-domain overlap with ODE or complex analysis | The rejected-candidate list records the overlaps and the replacements |
| Thin calculus definitions | A candidate without three distinct computations is dropped before authoring |
| New singleton units | Measure the partition after authoring and identify every singleton before acceptance |

## Review order

1. Author five Calculus III concepts: orientation, oriented flux, differential forms, the wedge
   product and the exterior derivative.
2. Author five real-analysis concepts: limsup, liminf, cluster points, Lebesgue measure and the
   integral of a simple function.
3. Run the domain tests, the lesson audit and the shipping contract on those ten.
4. Inspect the ten manually, including every family's 60-draw pool.
5. If the process is sound, author the remaining forty.
6. Regenerate mastery, graph, flow, docs and variety artifacts.
7. Inspect all 150 families, record mathematical and presentation defects separately and commit
   locally.
8. Do not push or deploy. Keep the two unrelated server files unstaged.

## Artifact contract

Every concept needs:

- a numeric concept node in `data/concepts/<domain>.json`;
- a generator in `internal/generator/<domain>/wave5.go`, registered from that package's `Register`;
- an independent `wave5_test.go` that recomputes every answer from `Facts`;
- a lesson in `data/lessons/authored/<id>.md` with one H1 and exactly three H2 sections over 1 KB;
- a three-entry KP shard in `data/lessons/kp/<id>.json`;
- a `data/lessons/lessons.json` row;
- a target in `calc3` or `ra`.

Neither domain is a learner domain, so no solution schema is required. The frozen additive baseline
stays at 1,500 concepts. No quality exemption is added.

## Approval checkpoints

1. The allocation is approved.
2. This ledger is reviewed with all 50 ids and prerequisite edges.
3. The first ten concepts are authored and manually reviewed.
4. The remaining forty are authored only after the first-ten approval.
5. Full validation is run and the work is committed locally.
