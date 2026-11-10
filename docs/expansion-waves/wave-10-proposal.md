# Wave 10 proposal

Status: approved and authoring complete on 2026-11-10. The allocation and candidate ledger were approved. The first ten passed review and the remaining forty were authorized. Wave 10 closed at 2,000 concepts.

Wave 9 closed at 1,950 concepts. Exactly 50 concepts remain before the 2,000 target. This proposal recommends a final split that closes every remaining domain target exactly.

## Recommended allocation

| Domain | Tier | Now | Target | Add | Course |
|---|---|---:|---:|---:|---|
| prealgebra | school | 71 | 89 | +18 | `pa` and grades 6 to 8 |
| geometry | school | 70 | 87 | +17 | `geo` |
| graph_theory | university | 63 | 70 | +7 | `graph` |
| machine_learning | university | 67 | 75 | +8 | `ml` |
| **Total** | | | | **+50** | |

Wave 10 would take the corpus to 2,000 concepts. No domain would be left below its target, so the expansion would close exactly rather than leave a residual row for another wave.

## Why this split

The remaining 50 are already fixed by the target table. There are only four domains left, and their deficits sum to 50. The recommendation does not invent a new allocation; it finishes the existing one.

The split also separates the two kinds of work cleanly. Prealgebra and geometry are learner domains with authored solution schemas, concrete diagrams and step-by-step explanations. Graph theory and machine learning are university domains with exact graph algorithms, small matrices, numeric losses and optimization examples. Both groups can support three genuinely different question families without theorem-name recall.

This is the final wave, so the checkpoint is quality rather than volume. Reaching 2,000 is a corpus milestone, not a reason to stop improving the existing material. Every new concept must still earn its place and every lesson must teach the operation its generator tests.

## Alternatives considered

**Reopen an advanced domain.** Add more abstract algebra, topology or real analysis and move part of the school allocation to a later wave. This is rejected because it leaves the approved target table unfinished and repeats the theorem-heavy workload that Wave 7 already flagged as a quality risk.

**Split the school content across two waves.** Take prealgebra and geometry now and leave graph and machine learning for a later wave. This does not close the target in 50 concepts and would require an extra wave for 15 university concepts.

**Take graph and machine learning first.** Finish the university remainder now and leave school for a later wave. This is mathematically possible, but it delays the learner-facing school content that still needs schemas and diagrams. The final wave would then be almost entirely school work, with no university balance.

The recommendation is the first table: prealgebra 18, geometry 17, graph theory 7 and machine learning 8.

## Overlap audit

Before authoring the first ten, the candidate ledger must record an audit for every proposed concept:

1. **Duplicate identity.** Check every id and normalized label against all 23 domains. Check semantic overlap as well as exact string collision.
2. **Genuine advancement.** Prealgebra must not repeat arithmetic or early algebra under a new name. Geometry must not repeat an existing theorem, construction or formula. Graph theory must not repeat a graph class or algorithm already in the corpus. Machine learning must not repeat an existing loss, model or evaluation concept.
3. **Prerequisite honesty.** Every new edge must point to an existing concept or another same-wave candidate. No missing prerequisite may be invented at authoring time.
4. **Independent mastery objective.** Each concept must be teachable and assessable on its own. A second worked example of an existing concept is not a new concept.
5. **Family fit.** The three families must measure the stated objective. A family that asks for the title, the theorem name or a definition without a computation fails the audit.

Specific overlap areas to check: prealgebra ratio, percent and integer material against arithmetic, fractions and algebra; geometry triangle, circle and solid material against existing geometry and trigonometry; graph classes and algorithms against the existing graph-theory and discrete-mathematics sets; machine-learning losses and optimizers against the existing ML and statistics material.

## Quality gates

Every Wave 10 concept must meet the existing expansion contract.

- Three section-aligned question families.
- At least five distinct questions per ordinal in twelve draws.
- Design target of twenty distinct per ordinal in 200 draws.
- Minimum accepted pool of ten distinct per ordinal in 60 draws or a written reason.
- Finite pools must report their size.
- Independent fact-based recomputation of every answer.
- No always-yes family and no theorem-name recall.
- A lesson over 1 KB with one H1 and exactly three H2 sections.
- A three-entry knowledge-point shard.
- A lesson-ledger row and a deliberate course target.
- An authored solution schema in every learner domain. Prealgebra and geometry need schemas. Graph theory and machine learning do not.

The wave will be inspected manually across all 150 families. Mathematical defects and presentation defects will be counted separately. No baseline exemption will be added. The 1,500-concept additive baseline stays frozen. Mastery band width stays 3. Commits stay local. Nothing will be pushed or deployed.

## Review order

The first ten will cover both learner domains and both university domains:

- Four prealgebra concepts.
- Three geometry concepts.
- Two graph-theory concepts.
- One machine-learning concept.

The first ten are authored and inspected before the remaining forty. If they expose a duplicate, a prerequisite error, a schema problem or a generator-shape problem, the process is corrected before the rest of the wave.

## Approval sequence

1. Approve or change the allocation in this proposal.
2. Create `docs/expansion-waves/wave-10-candidates.md` with exactly 50 ids, labels, subdomains, prerequisites, three-family plans and the overlap audit.
3. Review the candidate ledger.
4. Author and inspect the first ten.
5. Author the remaining forty only after the first ten are approved.
6. Run the full gate, record the final review and commit locally.

The proposal is complete. The next decision is whether to approve the recommended allocation or choose one of the alternatives.
