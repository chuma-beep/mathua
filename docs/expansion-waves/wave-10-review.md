# Wave 10 review

Date: 2026-11-10. Scope: the 50 concepts added by Wave 10, from 1,950 to 2,000.

This file records the inspection result as a metric so later work can be compared against it.
Passing gates is technical correctness. The counts below come from reading generated questions,
answers and explanations, not from the test suite.

## Defect metrics

| Metric | Wave 10 |
|---|---:|
| Concepts inspected | 50 |
| Knowledge-point families inspected | 150 |
| Mathematical defects | 1 |
| Framing defects | 2 |
| Presentation defects | 10 |
| Data-integrity defects | 1 |
| Defects fixed | All identified |
| Families below 10 distinct in 200 draws | 0 |
| Singleton units before and after | 22 and 20 |
| Singleton units created by the wave | 0 |
| Singleton units removed by the wave | 2 |

## Mathematical defect

One family computed the wrong quantity, and the independent test caught it before review.

1. `pct.ops.successive` KP2 applied the second percentage change to the original base instead of
   the value after the first change. The generator now applies the decrease to the intermediate
   value, so an increase followed by an equal decrease produces a net decrease rather than the
   original amount.

## Framing defects

Two families asked a question whose wording did not match the signed or empty answer.

1. `ml.ensemble.stacking` KP2 asked how far the meta-prediction was from a base prediction, but
   returned a signed difference. It now asks for the meta-prediction minus the base prediction.
2. `graph.feedback_arc_set` KP2 could print an empty proposed set as a blank after the word "is".
   The generator now displays an empty set as the word `empty`.

## Presentation defects

Ten learner-visible rendering defects were fixed.

1. `dec.ops.div_decimal` printed a floating-point artefact such as `6.800000000000001`. The
   dividend is now computed in integer tenths and formatted once.
2. `prealg.expr.two_variables` printed unit coefficients such as `1y`. Both coefficients are now
   at least two.
3. `prealg.expr.coefficient` printed expressions such as `2x + 1x`. Both variable coefficients
   are now at least two.
4. `prealg.expr.distribute_simplify` printed `1(x + 8)`. The outside factor is now at least two.
5. `geo.transform.composition` printed `1 - -4`. Subtracted values are parenthesised.
6. `geo.transform.isometry` printed `-1^2`. Coordinate differences are parenthesised before
   squaring.
7. `geo.coord.sphere_equation` printed `(z - 0)^2`. A zero center now renders as `z^2`.
8. `ml.loss.hinge` printed `1 - -3`. The margin is parenthesised inside the loss expression.
9. `ml.ensemble.stacking` printed products such as `1 x -3`. Base predictions and bias are now
   non-negative.
10. `ml.ensemble.voting` explained an absolute margin as `2 - 3 = 1`. The explanation now
    subtracts the smaller count from the larger one.

## Data-integrity defect

The first full validation run found stale corpus counts across `README.md`, `AGENTS.md`,
`CONTRIBUTING.md`, the architecture and mastery documents, ADR-050 and the published docs. The
counts script rewrote the README and the affected prose was updated to 2,000 concepts, 2,000 KP
shards and 6,000 knowledge points. The target table in `docs/content-expansion.md` now shows
prealgebra, geometry, graph theory and machine learning at their targets with no remaining add.

A later artifact audit found that the first ten concepts did not have `key_prerequisites` entries
in `data/concepts/enrichment.json`, while the remaining forty did. All ten entries were added so
every Wave 10 concept carries the same enrichment metadata.

The lesson `geo.solid.hemisphere` was 952 bytes after authoring. It was extended to 1,457 bytes so
the thin-lesson guard remains satisfied without an exemption.

## Pool and parameterisation work

Every new family is tested at 200 draws per ordinal. The tests require at least ten distinct
questions at that size and at least five in twelve draws. All 150 families pass both thresholds.
The strongest pools are the arithmetic and geometry families with large integer parameter ranges.
The smallest pools are the finite-table graph and machine-learning families, which still exceed
ten distinct stems through their claimed-drawing, vote-vector or residual parameters.

## Singleton discrepancy

Wave 9 closed with 143 units and 22 singleton units at 1,950 concepts. Wave 10 closes with 143
units and 20 singleton units at 2,000 concepts.

- No singleton unit was added.
- `geometry.b5`, which contained only `geo.solid.frustum`, is no longer a singleton. New geometry
  concepts and prerequisite edges moved depth neighbours into its band.
- `machine_learning.b4`, which contained only `ml.data.scaling`, is no longer a singleton. New
  machine-learning concepts moved depth neighbours into its band.
- The rule and band width are unchanged: `unit = (domain, floor(conceptDepth / bandWidth))` at
  width 3.

The net change is caused by the added concepts changing topological depths, not by a change to
the partition algorithm.

## Artifact scope

Wave 10 added 50 concepts: 18 prealgebra, 17 geometry, 7 graph theory and 8 machine learning.
It closes every remaining domain target: prealgebra 89, geometry 87, graph theory 70 and machine
learning 75. Courses are `pa` with grades 6 to 8, `geo`, `graph` and `ml`.

Prealgebra and geometry are learner domains. All 28 new concepts in those domains have authored
solution schemas. Graph theory and machine learning are not learner domains and add none. Learner
schema coverage is 467 of 602 concepts.

## Quality limits

The graph-theory and machine-learning families are intentionally small and numerical. They test
the defining computation or invariant, not a proof or a full algorithm trace. The crossing-number
family uses the established values for K4 through K8 and K3,3; it does not claim values for larger
complete graphs. The spectral-clustering family uses sign assignments and path cuts rather than a
full eigenvector computation.

The lessons are concept-specific but several use the same three-part rhythm: define the object,
state the relation, then check the result. That is deliberate for a final school-content wave, but
a later editorial pass should vary the openings more.

## Validation

Wave 10 closes with 2,000 concepts, 2,000 lessons, 2,000 KP files, 6,000 knowledge points, 28
courses and all concepts reachable. The graph validator reports 23 domains and 19 reviewed
collisions. The lesson audit reports zero issues. Schema, shipping, variety, duplicate-registration,
counts, additive and published-docs checks pass. `make validate`, `go test ./...` and
`git diff --check` pass.

The additive baseline remains frozen at the 1,500-concept snapshot. The commit is local. Nothing was
pushed or deployed, and the two unrelated server files remain unstaged.
