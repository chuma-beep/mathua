# Wave 5 proposal

Status: approved and authoring complete on 2026-11-09. The allocation and candidate ledger are
approved. The first ten passed the independent audit; the remaining forty were authored after that
approval. The final corpus is 1,750 concepts.

Date: 2026-11-09. Starting corpus: 1,700 concepts. Frozen additive baseline: 1,500 concepts. Wave 5
target: 1,750 concepts.

This document proposes the domain split, allocation, prerequisite strategy, generator coverage
expectations and quality gates for Wave 5. The exact concept ids and labels belong to a candidate
ledger, which is authored only after this proposal is approved.

## Recommended split

| Domain | Now | Target | Add | Tier | Course |
|---|---:|---:|---:|---|---|
| calc3 | 71 | 97 | +26 | bridge | `calc3` |
| real_analysis | 65 | 89 | +24 | university | `ra` |
| **Total** | | | **50** | | |

This is the largest-deficit split. It closes the two largest remaining rows exactly and leaves 250
concepts for Waves 6 through 10. It also completes two domains whose targets have been open since
the 2,000-concept allocation was agreed.

### Why these two domains

- Calc III has the largest remaining deficit. Its existing 71 concepts cover the standard
  multivariable and vector-calculus spine, but not the advanced integral theorems, differential
  forms, vector potentials, harmonic functions, variational methods and weak formulations that
  naturally follow.
- Real analysis has the second-largest deficit. Its existing 65 concepts cover order, sequences,
  series, limits, continuity, differentiation, Riemann integration, metric spaces, measure zero
  and several approximation theorems. The remaining 24 can extend that spine into limsup and
  liminf, Lebesgue integration, Lp spaces, Fourier analysis and functional-analytic examples.
- The other large deficits are discrete mathematics +23, precalculus +23, probability theory +22,
  abstract algebra +20, machine learning +20, topology +20, trigonometry +20, arithmetic +18,
  prealgebra +18, geometry +17 and fractions +15. None of them closes two target rows as cleanly.

### The risk, stated plainly

This is the most theorem-heavy wave so far. Real analysis in particular can collapse into
classification questions whose three families are the same question with different wording. The
candidate ledger must therefore prefer concepts that permit exact computation on finite examples,
explicit functions, small metric spaces or published facts. A candidate that cannot support three
genuinely distinct computational families is dropped before authoring, not shipped weak.

Calc III also carries a duplicate risk. Its existing corpus already has line integrals, flux,
Green's theorem, Stokes' theorem, the divergence theorem, Jacobians, Lagrange multipliers, Taylor
in two variables, implicit and inverse function theorems and vector identities. The ledger must
check every candidate against those labels and against the differential-equations and topology
domains before authoring.

## Alternatives, if the recommendation is not approved

1. **School foundations:** trigonometry +20 (78 to 98), arithmetic +18 (72 to 90) and fractions
   +12 (52 to 64). Total 50. This is easier to generate and strengthens the school tier, which has
   88 concepts left. It leaves the two largest advanced deficits open and uses three domains.
2. **Mixed completion:** calc3 +26, real_analysis +12, graph_theory +11 and statistics +1. Total
   50. This completes calc3, graph theory and statistics, and moves real analysis halfway. It
   spreads the wave over four domains and leaves real analysis at 77 of 89.
3. **Discrete and probability:** discrete_math +23, probability_theory +22 and statistics +5.
   Total 50. This is the most computational alternative and avoids real analysis, but it leaves
   calc3 and real analysis untouched.
4. **Bridge foundations:** precalculus +23, calc1 +14, calc2 +8 and statistics +1. Total 46; four
   more would need a second domain. This is not recommended because precalculus duplicates algebra
   and Calculus I more often than any other domain.

## Candidate topic clusters

These are planning clusters, not approved ids. Every proposed id and normalized label will be
checked against all `data/concepts/` files in the candidate ledger. A cluster that cannot supply
three genuinely distinct families will be replaced.

### Calc III: 26

| Cluster | Candidate topics | Expected prerequisite source |
|---|---|---|
| Integral theorems on surfaces | Orientation, surface flux, circulation, generalized Stokes, divergence on regions | `calc.vec.stokes`, `calc.vec.divergence_theorem` |
| Differential forms | Differential forms, wedge product, exterior derivative, pullback, closed and exact forms | `calc.vec.identities`, `calc.vec.curl_3d` |
| Vector potentials | Vector potentials, Helmholtz decomposition, gauge freedom, solenoidal fields | `calc.vec.conservative_test`, `calc.vec.potential_value` |
| Harmonic functions | Harmonic functions, mean value property, maximum principle, Poisson equation, Dirichlet and Neumann problems | `calc.deriv.partial_higher`, `calc.vec.divergence` |
| Green identities | Green's identities, Green's functions in several variables, method of images, Dirichlet kernels | `calc.vec.greens_theorem`, `calc.vec.divergence_theorem` |
| Variational calculus | First variation, Euler-Lagrange equation, geodesics, minimal surfaces, constrained variations | `calc.multi.lagrange`, `calc.vec.arc_length` |
| Fourier analysis in several variables | Multivariable Fourier transform, convolution, Parseval, heat kernel, spherical harmonics | `calc.seq.function_sequences`, `calc.integral.double_polar` |
| Weak formulations | Distributions, delta function, weak derivatives, Sobolev spaces, integration by parts for weak derivatives | `calc.deriv.partial_higher`, `calc.multi.triple_integral_region` |
| Topology of level sets | Regular values, level surfaces, tangent spaces, orientation of manifolds, boundary orientation | `calc.multi.level_curves`, `calc.vec.surface_tangent` |
| Advanced change of variables | Change of variables for manifolds, Jacobian of a parametrisation, surface measure, volume forms | `calc.multi.transform_jacobian`, `calc.multi.coordinates_jacobian` |

The exact 26 will be chosen only after a collision check. Existing labels already own many basic
topics, so the ledger should prefer the advanced entries in these clusters.

### Real analysis: 24

| Cluster | Candidate topics | Expected prerequisite source |
|---|---|---|
| Limsup and liminf | Limsup and liminf of sequences, subsequential limits, cluster points, limsup of sets | `ra.seq.subsequence`, `ra.seq.bolzano_weierstrass` |
| Series refinements | Cauchy products, Dirichlet's test, double series, Fubini for series, convergence of rearrangements | `ra.seq.series`, `ra.seq.rearrangement` |
| Measure and integration | Lebesgue measure, measurable sets, measurable functions, Lebesgue integral, simple functions | `ra.measure.null`, `ra.measure.cantor_set` |
| Convergence theorems | Monotone convergence, Fatou's lemma, dominated convergence, uniform integrability | `ra.deriv.uniform_conv`, `ra.measure.null` |
| Function spaces | Lp spaces, Holder's inequality, Minkowski's inequality, completeness of Lp, essential supremum | `ra.metric.complete`, `ra.metric.compact` |
| Functional analysis examples | Banach spaces, Hilbert spaces, bounded linear operators, compact operators, Riesz representation | `ra.metric.complete`, `ra.deriv.arzela_ascoli` |
| Fourier analysis | Fourier series, convergence in L2, convolution, approximate identities, Fejer's theorem | `ra.deriv.weierstrass_approx`, `ra.deriv.uniform_conv` |
| Measure refinements | Egorov's theorem, Lusin's theorem, Radon-Nikodym theorem, product measures, Fubini-Tonelli | `ra.measure.null`, `ra.measure.cantor_set` |
| Differentiation refinements | Differentiation under the integral sign, absolutely continuous functions, Lebesgue differentiation | `ra.deriv.ftc`, `ra.deriv.bounded_variation` |
| Special functions | Gamma function, Beta function, Stirling's formula, asymptotic estimates | `ra.deriv.taylor`, `ra.deriv.improper` |

The exact 24 will be chosen only after a collision check against `complex.adv.gamma`, the existing
real-analysis series concepts and the topology domain.

## Curriculum overlap audit gate

Before the first ten concepts are authored, the candidate ledger must pass five checks for every
candidate. This gate is mandatory. A candidate that fails any check is replaced before authoring.

1. **Duplicate identity.** Check the exact id and normalized label against every file in
   `data/concepts/`, not only the two target domains. Check the same topic under a different id or
   a different wording.
2. **Genuine advancement.** Check that the concept is genuinely more advanced than the existing
   Calculus II and Calculus III material, or than the existing real-analysis sequence and
   integration spine. A renamed review of an existing idea is rejected.
3. **Prerequisite honesty.** Check that every prerequisite is a real mathematical dependency. A
   missing prerequisite is added to the graph only when the dependency is essential. Do not treat
   Fourier analysis as an automatic extension of Lebesgue integration if the actual prerequisite
   chain does not support it.
4. **Independent mastery objective.** State what a learner can do after the concept that they could
   not do before. If the objective is only to recognise a theorem name, the candidate is rejected.
5. **Family fit.** Confirm that the three proposed families measure that objective with distinct
   mathematics. A family that changes only wording, coefficients or option order is not a second
   family. Prefer explicit functions, finite counting measures, small metric spaces and exact
   numerical invariants.

The ledger records the audit result for each candidate, not only for the wave as a whole. A
candidate with a genuine overlap is replaced rather than admitted to a quality exemption.

## Approval sequence

The wave follows the user-approved sequence.

| Stage | Decision |
|---|---|
| Proposal | Allocation approved |
| Candidate ledger | Inspect all 50 ids and prerequisite edges |
| First ten concepts | Manual quality review |
| Remaining forty | Author only after first-ten approval |
| Full validation | Required |
| Commit | After every gate passes |

Reaching 1,750 concepts is a corpus milestone. It is not automatically a learning-quality
milestone. The quality of the 50 additions takes precedence over hitting the number.

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
| Wave 5 design target | At least 20 distinct questions per ordinal in 200 draws |
| Minimum accepted pool | At least 10 distinct questions per ordinal in 60 draws, or a written reason in the candidate ledger |
| Answer verification | Every numeric answer recomputed from published `Facts` by a test independent of the generator answer |
| Section alignment | Different question at each ordinal under the same seed |
| Finite pools | A family with a closed parameter set must report its size and cannot be described as parametric |
| Manual inspection | All 150 families read for answer correctness, explanation correctness and presentation |

The real-analysis risk is handled by choosing families that can be checked on finite examples. A
sequence can have an explicit formula. A metric space can have finitely many points. A measure can
be a finite counting measure. A function space can be finite-dimensional. A theorem can be asked
through a small counterexample or a numeric invariant. If a candidate cannot be made exact in that
way, it is not a Wave 5 candidate.

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
mapping may be changed, no prerequisite may be removed and no variety floor may be lowered. No entry
may be added to a quality exemption baseline.

Every new concept also needs:

- a concept node with numeric grading and a mastery threshold;
- a generator registered from its domain package;
- an independent fact-based test;
- a lesson over 1 KB with exactly three H2 sections;
- a three-entry KP shard whose sections resolve;
- a `lessons.json` row;
- a target in `calc3` or `ra`.

Both domains are non-learner domains, so no solution schema is required.

## Review order

1. Author five calc3 concepts and five real-analysis concepts.
2. Run the domain tests, the lesson audit and the shipping contract on those ten.
3. Inspect the ten manually, including every explanation and every family's 60-draw pool.
4. Fix the process before authoring the remaining forty if a defect appears.
5. Author the remaining forty.
6. Regenerate mastery, graph, flow, docs and variety artifacts.
7. Inspect all 150 families, record mathematical and presentation defects separately and commit
   locally.
8. Do not push or deploy. Keep the unrelated server files unstaged.

## Approval checkpoints

1. Approve the recommended split or choose an alternative.
2. Approve the candidate ledger with exact ids, labels, subdomains, prerequisite edges and generator
   family plans.
3. Author and inspect the first ten.
4. Author the remaining forty.
5. Close out with the full gate and a local commit.

No concept generation begins before checkpoint 1 and checkpoint 2 are approved.
