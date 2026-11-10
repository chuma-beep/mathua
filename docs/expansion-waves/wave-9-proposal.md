# Wave 9 proposal

Status: approved and authoring complete on 2026-11-10. The allocation was approved, the first ten concepts passed review, and the remaining forty were authored after that approval. The corpus is now 1,950 concepts.

Wave 8 closed at 1,900 concepts. The remaining 100 concepts are unevenly distributed. This proposal recommends a Wave 9 split of 50 that closes every bridge domain and leaves a clean 50-concept Wave 10.

## Recommended allocation

| Domain | Tier | Now | Target | Add | Course |
|---|---|---:|---:|---:|---|
| precalculus | bridge | 58 | 81 | +23 | `pc` |
| calc1 | bridge | 100 | 114 | +14 | `calc1` |
| calc2 | bridge | 89 | 97 | +8 | `calc2` |
| fractions | school | 64 | 67 | +3 | `pa` and grade 6 |
| machine_learning | university | 65 | 67 | +2 | `ml` |
| **Total** | | | | **+50** | |

Wave 9 would take the corpus to 1,950 concepts. Wave 10 would then close the remaining 50 exactly: prealgebra +18, geometry +17, graph theory +7 and machine learning +8.

## Why this split

The remaining bridge allocation is 45 concepts. It is the largest coherent block left, and it is also the block with the highest overlap risk: precalculus sits between algebra and calculus, while Calculus I and Calculus II extend material the corpus already teaches. Concentrating that work in one wave makes the curriculum overlap audit a real gate rather than a formality. The audit can examine the whole bridge boundary at once.

The three extra concepts are deliberately small. Fractions needs three concepts to close its target, and machine learning needs a small extension rather than a full wave. Adding them here keeps every remaining row exact and leaves Wave 10 as a school-foundations wave with the rest of graph theory and machine learning.

This split also keeps the final wave manageable. Wave 10 would contain prealgebra, geometry, graph theory and machine learning. It would not mix the bridge overlap audit with the closing school work.

## Alternatives considered

**School-first wave.** Prealgebra 18, geometry 17, fractions 3, graph theory 7 and machine learning 5 would reach 50. It would close all school content and leave all bridge work for Wave 10. It is attractive because school foundations are computationally straightforward. It loses the chance to concentrate the bridge overlap audit now, and Wave 10 would carry all 45 bridge concepts plus five machine-learning concepts.

**Mixed bridge and school wave.** Precalculus 23, Calculus I 14, Calculus II 8 and five school concepts would also reach 50. It spreads the remaining school work across both waves. It makes the Wave 10 allocation less clean and does not improve the bridge audit.

**Bridge-only wave with a different tail.** Precalculus 23, Calculus I 14, Calculus II 8 and five machine-learning concepts would close the bridge and leave school 38 plus graph 7 plus machine learning 5 for Wave 10. It is defensible, but it delays fractions closure and gives machine learning a larger share than its remaining target justifies.

The recommendation is the first table: bridge closure, fractions closure and two machine-learning concepts.

## Curriculum overlap audit

Before authoring the first ten, the candidate ledger must record an audit for every proposed concept:

1. **Duplicate identity.** Does another domain already teach the same label or the same mathematical object? Normalized label collisions and semantic overlaps are both checked.
2. **Genuine advancement.** Is the concept actually beyond the existing Algebra II, precalculus, Calculus I or Calculus II treatment, or is it a renamed prerequisite?
3. **Prerequisite honesty.** Did the concept introduce a dependency that should already exist? Every new edge must point to an existing concept or another same-wave candidate.
4. **Independent mastery objective.** Can the concept be mastered separately from its neighbours? A concept that is only a second worked example of an existing one is not a new concept.
5. **Family fit.** Do the three question families measure the stated objective, or do they test the title rather than the mathematics?

The audit will also check the specific boundary areas: function transformations, inverse functions, rational functions, sequences and series, limits, derivatives, integrals and series convergence. Any candidate that fails the audit is dropped before authoring, not shipped weak.

## Quality gates

Every Wave 9 concept must meet the existing expansion contract:

- Three section-aligned question families.
- At least five distinct questions per ordinal in twelve draws.
- Design target of twenty distinct per ordinal in 200 draws.
- Minimum accepted pool of ten distinct per ordinal in 60 draws or a written reason.
- Independent fact-based recomputation of every answer.
- A lesson over 1 KB with one H1 and exactly three H2 sections.
- A three-entry knowledge-point shard.
- A lesson-ledger row and a deliberate course target.
- An authored solution schema in every learner domain. Fractions is a learner domain. Precalculus, Calculus I, Calculus II and machine learning are not learner domains and do not need schemas.

The wave will be inspected manually across all 150 families. Mathematical defects and presentation defects will be counted separately. No baseline exemption will be added. The 1,500-concept additive baseline stays frozen. Mastery band width stays 3. Commits stay local. Nothing will be pushed or deployed.

## Review order

The first ten will be chosen to cover the riskiest boundary first:

- Three precalculus concepts.
- Three Calculus I concepts.
- Two Calculus II concepts.
- One fractions concept.
- One machine-learning concept.

The first ten are authored and inspected before the remaining forty. If they expose a duplicate, a prerequisite error or a generator-shape problem, the process is corrected before the rest of the wave.

## Approval sequence

1. Approve or change the allocation in this proposal.
2. Create `docs/expansion-waves/wave-9-candidates.md` with exactly 50 ids, labels, subdomains, prerequisites, three-family plans and the overlap audit.
3. Review the candidate ledger.
4. Author and inspect the first ten.
5. Author the remaining forty only after the first ten are approved.
6. Run the full gate and commit locally.

The proposal is complete. The next decision is whether to approve the recommended allocation or choose one of the alternatives.
