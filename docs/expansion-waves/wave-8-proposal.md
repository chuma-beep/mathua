# Wave 8 proposal

Status: approved and authoring complete on 2026-11-10. The first ten were approved after review and the remaining forty were authored after that approval.

Wave 8 is the school-foundations wave that follows the theorem-heavy Wave 7. It adds 50 concepts and moves the corpus from 1,850 to 1,900. The split is trigonometry +20, arithmetic +18 and fractions +12. All three domains are school domains.

## Allocation

| Domain | Tier | Now | Target | Add | Course |
|---|---|---:|---:|---:|---|
| trigonometry | school | 78 | 98 | +20 | `pc` |
| arithmetic | school | 72 | 90 | +18 | `pa` and the grade courses |
| fractions | school | 52 | 64 | +12 | `pa` and the grade courses |
| **Total** | | | | **50** | |

The remaining 100 concepts after Wave 8 are school 38, bridge 45 and university 17. Waves 9 and 10 can then close the corpus without reopening the advanced theorem-heavy domains immediately after Wave 7.

## Why this split

Wave 7 produced zero mathematical answer defects but left five recorded recognition risks in Schur's lemma, induced representations, partitions of unity, continua and Brouwer. Wave 8 deliberately does not add more theorem-heavy material while those risks are being monitored. It uses foundational domains where exact computation and genuinely different question families are natural.

The split also leaves precalculus and calculus for Wave 9. That is the wave where the algebra/precalculus/calculus overlap audit matters most, so the audit gets a dedicated batch instead of being folded into a school wave.

## Alternatives considered

- Complete the bridge first: precalculus +23, calc1 +14, calc2 +8, machine_learning +5.
- A different school mix: trigonometry +20, geometry +17, fractions +13.
- A mixed close: graph_theory +7, machine_learning +10 and 33 school concepts.

The approved allocation is trigonometry, arithmetic and fractions because it gives the largest coherent school batch and keeps the bridge work together for Wave 9.

## Overlap audit gate

Wave 8 must run the same five checks as Wave 7 before authoring the first ten.

1. **Duplicate identity.** No exact id or normalized-label collision against all 23 domains.
2. **Genuine advancement.** A candidate must not repeat an existing lesson under another name. Arithmetic mental strategies, alternative number bases and fraction representations must be checked against the existing arithmetic and prealgebra material.
3. **Prerequisite honesty.** Every prerequisite must exist. No candidate may reference a later Wave 8 candidate unless that candidate is authored first.
4. **Independent mastery objective.** Each concept needs a learning objective that the existing corpus does not already measure.
5. **Family fit.** Each concept needs three question families that measure that objective, not three rewordings of one computation.

Any candidate that fails a check is replaced before authoring, not shipped weak.

## Quality gates

- Three section-aligned question families per concept.
- Hard floor 5 distinct questions per ordinal in 12 draws.
- Working target 10 distinct per ordinal in 60 draws and 20 per ordinal in 200 draws.
- Every numeric answer independently recomputed from `Facts`.
- Manual inspection of all 150 families, with mathematical and presentation defects counted separately.
- No always-true family and no theorem-name recall. Theorem-heavy content is not planned for Wave 8, but any definitional family must still be paired with an application or failure case.
- Lessons over 1 KB with one H1 and exactly three H2 sections.
- Three-entry knowledge-point shards whose sections resolve in the lesson.
- Course reachability for every concept.
- No solution schemas because all three domains are school domains and the existing arithmetic and fraction schema policy still applies to their computed concepts. The candidate ledger will identify the schema-bearing concepts explicitly.
- The 1,500-concept additive baseline stays frozen. Band width stays 3.
- Local commits only. No push and no deployment.

## Review order

Author and inspect ten first: four trigonometry, three arithmetic and three fractions. Run the domain tests, lesson audit, graph validator and shipping contract on those ten. If they expose a generator or schema problem, fix the process before authoring the remaining forty.

## Closeout

The full wave runs `make validate`, `go test ./...` and `git diff --check`, regenerates mastery units, graph, flow, docs and the variety baseline, records defects in `docs/expansion-waves/wave-8-review.md`, and commits locally with a message such as `Add Wave 8: 50 concepts to 1,900`. Nothing is pushed or deployed.
