# Wave 11 proposal

Status: proposed on 2026-11-10.

Wave 11 raises the corpus from 2,000 to 2,050 concepts. It is the first half of the 2,100 pilot and the first wave after the approved 2,000 target closed. The wave adds one new university domain, Optimization, with one course, `optim`.

Wave 12 is not authorized by this proposal. Its 50 concepts depend on what the Wave 11 review measures.

## Recommended allocation

| Domain | Tier | Now | Target | Add | Course |
|---|---|---:|---:|---:|---|
| optimization | university | 0 | 50 | +50 | `optim` |
| **Total** | | | | **+50** | |

The 2,100 pilot is therefore two waves: Wave 11 to 2,050 and a later decision about Wave 12 to 2,100.

## Why Optimization

Optimization is the strongest new domain after 2,000 because it connects four domains the corpus already teaches:

- calculus supplies gradients, Hessians and Lagrange multipliers
- linear algebra supplies systems, rank and least squares
- numerical analysis supplies line search, conjugate gradients and Gauss-Newton
- machine learning supplies gradient descent, stochastic methods and regularization

It is also the clearest place where the existing corpus stops short. The current material teaches the derivative test, one constrained example and a set of optimizer names. It does not teach convexity as an optimization structure, linear programming duality, KKT conditions, penalty and barrier methods or integer programming.

The domain is not a rename of the existing calculus or machine-learning material. Every candidate below has a distinct optimization objective, and the overlap audit records the boundary.

## Course

The new course is `optim`, a university course. Its targets are the five subdomain sinks:

- `optim.second_order`
- `optim.hessian`
- `optim.kkt`
- `optim.duality_gap`
- `optim.integrality_gap`

Those targets close over all 50 concepts. The course does not replace the existing calculus, numerical-analysis or machine-learning courses. It gives a learner a path that starts at `calc.deriv.optimization` and reaches the optimization model, algorithm and duality material.

## Quality gates

Wave 11 keeps the full expansion contract.

- Three section-aligned question families per concept.
- At least five distinct questions per ordinal in twelve draws.
- Design target of twenty distinct per ordinal in 200 draws.
- Independent fact-based recomputation of every answer.
- Manual inspection of all 150 families.
- A lesson over 1 KB with one H1 and exactly three H2 sections.
- A three-entry knowledge-point shard.
- A lesson-ledger row and a deliberate course target.
- No solution schema, because Optimization is a university domain and the learner-schema policy covers school and bridge domains only.
- No baseline exemption. The 1,500-concept additive baseline stays frozen.
- Mastery band width stays 3.

## Review order

The first ten cover the foundations and the first algorithm material:

1. `optim.problem.anatomy`
2. `optim.feasible.set`
3. `optim.local_global`
4. `optim.convex_set`
5. `optim.strict_convexity`
6. `optim.subgradient`
7. `optim.epigraph`
8. `optim.coercivity`
9. `optim.first_order`
10. `optim.second_order`

They are authored and inspected before the remaining forty. If they expose a duplicate, a prerequisite error or a generator-shape problem, the process is corrected before the rest of the wave.

## Alternatives considered

**Deepen existing domains only.** Add 50 concepts across differential equations, numerical analysis, number theory and topology. This reuses courses and generators, but it does not answer whether Mathua can add a new mathematical field without losing quality. The pilot is specifically about that question.

**Add a different new domain.** Game theory, stochastic processes, functional analysis and coding theory were considered. Optimization has the strongest prerequisite base, the clearest course path and the least semantic overlap with an existing domain.

**Add 100 concepts now.** A 100-concept first wave would double the review load before the new domain has a measured defect rate. It would not be a pilot.

## Approval sequence

1. Review the 50-candidate ledger in `wave-11-candidates.md`.
2. Author and inspect the first ten.
3. Author the remaining forty only after the first-ten review.
4. Regenerate mastery, graph, flow and documentation artifacts.
5. Run `make validate`, `go test ./...` and `git diff --check`.
6. Record the Wave 11 defect metrics and decide whether Wave 12 to 2,100 is authorized.
