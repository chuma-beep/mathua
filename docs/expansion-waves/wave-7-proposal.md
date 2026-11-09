# Wave 7 proposal

Status: approved, authored and closed on 2026-11-09. The allocation is complete. The first ten concepts passed independent review, and the remaining forty were authored after that approval.

Date: 2026-11-09. Current corpus: 1,800 concepts. Frozen additive baseline: 1,500 concepts. Wave 7
target: 1,850 concepts.

## Recommendation

Split Wave 7 as follows.

| Domain | Tier | Now | Target | Add | Course |
|---|---|---:|---:|---:|---|
| abstract_algebra | university | 56 | 76 | +20 | `abstract` |
| topology | university | 52 | 72 | +20 | `topo` |
| machine_learning | university | 55 | 75 | +10 | `ml` |
| **Total** | | | | **50** | |

This closes abstract algebra and topology exactly and moves machine learning halfway. It leaves 150
concepts for Waves 8 to 10: school 88, bridge 45 and university 17. That tail is unusually clean.
Wave 8 can carry a full school batch, Wave 9 can finish the bridge and five school concepts, and
Wave 10 can close the remaining school concepts, the last machine-learning concepts and graph
theory.

The alternative splits are weaker.

- A school wave of trigonometry +20, arithmetic +18 and fractions +12 is easy to generate but
  leaves every advanced deficit for the final three waves.
- A bridge wave of precalculus +23, calc1 +14, calc2 +8 and graph +5 closes the calculus targets but
  concentrates the duplication risk in precalculus, algebra and trigonometry, the overlap the review
  has warned about since Wave 2.
- A machine-learning +20 and graph +7 wave with 23 school concepts mixes two registers and makes the
  final school batch harder to stage.

## Why this split

Abstract algebra and topology are the two largest university deficits. Both already have substantial
curricula, so the new concepts must be genuinely more advanced or genuinely different in method.
Machine learning is the only applied university domain with a large deficit. Ten concepts are enough
to test whether the new material can be taught through computation and model behaviour rather than
through terminology.

The allocation does not touch precalculus. That is deliberate. Precalculus is the highest-risk row
for duplicate concepts, and this wave should not spend its review budget there.

## Candidate clusters

Exact ids, labels, subdomains, prerequisites and three-family plans belong in the candidate ledger.

**Abstract algebra.** Group actions and orbit counting, semidirect products, free groups and
presentations, abelianisation, group extensions, representation theory and characters, localisation
of rings, Euclidean and principal ideal examples, unique factorisation, polynomial irreducibility,
finite fields, Galois groups, solvability by radicals, module decompositions, exact sequences,
tensor products, category language and a first homological construction. Every candidate must be
checked against the existing 56 concepts and against the Wave 4 and Wave 5 complex and real-analysis
material.

**Topology.** Quotient spaces, product and box topologies, nets and filters, function spaces,
separation axioms, compactifications, connectedness refinements, homotopy, covering spaces,
fundamental-group calculations, van Kampen applications, homology and cohomology, manifolds,
surfaces, knot invariants and dimension. The existing 52 concepts already own many of these themes,
so the ledger must reject any concept that is only a renamed version of an existing one.

**Machine learning.** Regularisation, model selection, bias and variance, decision trees and
ensembles, support vector machines, clustering, principal components, neural networks and
backpropagation, sequence models and attention, reinforcement learning, causal inference, model
calibration and Bayesian machine learning. Prefer concepts with exact small-data computations and
observable model changes.

## Risks

1. **Theorem-heavy domains.** Abstract algebra and topology can collapse into theorem-name recall.
   Every concept needs at least one family that tests a condition, a counterexample, a construction
   or a hypothesis-sensitive computation. No family may be always yes.
2. **Existing density.** Abstract algebra and topology are already dense. The duplicate check must
   compare normalized labels and the actual mathematical content, not only the id.
3. **Machine learning overlap.** Machine learning shares linear algebra and statistics tools. The
   ledger must show that each concept has an independent mastery objective and is not a renamed
   regression or matrix topic.
4. **Course reachability.** Abstract algebra has two courses, `abstract` and `abstract2`; topology
   has `topo`; machine learning has `ml`. The ledger must choose the correct target course and show
   the prerequisite closure.

## Overlap audit gate

Before any authoring, the candidate ledger must pass five checks for every concept.

1. **Duplicate identity.** No exact id collision and no normalized-label collision against all 23
   domains.
2. **Genuine advancement.** The concept is not already taught under another id or as a section of an
   existing lesson.
3. **Prerequisite honesty.** Every prerequisite exists or is a same-wave candidate, the edge is
   mathematically necessary and no existing edge is removed.
4. **Independent mastery objective.** The concept has one observable skill that a learner can
   practise separately from its neighbours.
5. **Family fit.** Three genuinely different families can test that skill with exact answers and
   enough distinct questions.

## Quality gates

- Three section-aligned question families per concept.
- Hard floor of 5 distinct questions per ordinal in 12 draws.
- Design target of 20 distinct questions per ordinal in 200 draws.
- Minimum accepted pool of 10 distinct questions per ordinal in 60 draws or a written reason.
- Finite pools must report their size.
- Every numeric answer independently recomputed from generator facts.
- At least one family per theorem-heavy concept must test a condition, counterexample, construction
  or hypothesis-sensitive computation.
- No always-yes family and no theorem-name recall.
- Manual inspection of all 150 families, with mathematical and presentation defects counted
  separately.
- Full gate: `make validate`, `go test ./...` and `git diff --check`.
- Regenerate mastery units, graph, flow, docs and variety baseline.
- Keep the 1,500-concept additive baseline frozen and keep band width at 3.
- Local commits only. Do not push or deploy.

## Approval sequence

1. Approve the allocation.
2. Inspect the candidate ledger with all 50 ids and prerequisite edges.
3. Manually review the first ten concepts.
4. Author the remaining forty only after the first-ten review.
5. Run the full validation and artifact regeneration.
6. Commit locally after every gate passes.

Reaching 1,850 is a corpus milestone, not automatically a learning-quality milestone. The quality of
the 50 additions matters more than the number.
