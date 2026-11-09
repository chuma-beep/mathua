# Wave 7 review

Date: 2026-11-09. Scope: the 50 concepts added by Wave 7, from 1,800 to 1,850.

This file records the inspection result as a metric so later waves can be compared against it.
Passing gates is technical correctness. The counts below come from reading generated questions,
answers and explanations, not from the test suite.

## Defect metrics

| Metric | Wave 7 |
|---|---:|
| Concepts inspected | 50 |
| Knowledge-point families inspected | 150 |
| Mathematical answer defects | 0 |
| Framing and recognition defects fixed | 7 |
| Presentation defects fixed | 5 |
| Pool and parameterisation defects fixed | 17 |
| Curriculum-density defects fixed | 1 |
| Duplicate concepts replaced before acceptance | 0 |
| Defects fixed | All identified |
| Regression tests added | Yes |
| Singleton units before and after | 22 and 22 |
| Singleton units created by the wave | 0 |
| Singleton units removed by the wave | 0 |

## Mathematical and framing defects

No wrong answer was found in the 150 families. Seven theorem-heavy families nevertheless tested a
constant or theorem-name conclusion rather than a condition or a computation. All seven were
replaced.

1. `abstract.rings.dedekind` KP2 asked whether unique factorisation holds, which is always yes in
   the model. It now asks for the number of distinct prime ideals in the factorisation.
2. `abstract.structures.nakayama` KP2 asked whether \(JM\) is a proper submodule while the
   generator always chose \(j<d\). The generator now allows \(j=d\), so the hypothesis can fail.
3. `abstract.structures.exterior` KP2 asked for the square of a form, which is always zero. It now
   asks how many basis \(k\)-vectors contain the first coordinate.
4. `topo.paracompact` KP2 asked whether a locally finite refinement is locally finite. It now asks
   for the size of a refinement with a stated number of extra sets.
5. `topo.function_space` KP2 asked whether evaluation is continuous. It now asks for the number of
   constant functions between two finite sets.
6. `topo.uniform_space` KP2 asked whether every function on a finite discrete uniform space is
   uniformly continuous. It now counts the entourage pairs with a fixed first coordinate.
7. `ml.svm.soft_margin` KP3 asked whether increasing \(C\) allows more margin violations. It now
   asks which of two values of \(C\) gives the smaller margin.

The first-ten review had already checked the first ten concepts. It found no wrong answers and no
sign-convention defects, but strengthened the limsup/liminf and measure lessons.

## Remaining recognition risks

Some theorem-heavy concepts still contain a definitional family. These are recorded rather than
hidden: `abstract.rep.schur` (endomorphism dimension), `abstract.rep.induced` (Frobenius
multiplicity), `topo.partition_unity` (the sum is 1), `topo.continuum` (connectedness), and
`topo.brouwer` (the closed disk has a fixed point). Each of those concepts has at least one other
family that computes a dimension, count, characteristic or counterexample. A later wave should
replace the definitional family if a genuinely computational variant exists.

## Presentation defects

Four questions rendered polynomial data badly. All were fixed.

1. `abstract.poly.resultant` KP2 printed the second linear polynomial as `1 x+-4`.
2. `abstract.poly.discriminant` KP1 printed the quadratic as `1 x^2+3 x+-2`.
3. `abstract.poly.discriminant` KP3 printed `4 x^2+-4 x+-1`.
4. `abstract.poly.discriminant` KP2 repeated the same unit-coefficient and sign-join problem.
5. `topo.deck` used "Deck transformations" for both the lesson title and its first section, so
   section extraction returned the whole body.

A small polynomial formatter now suppresses unit coefficients and renders signed terms without
`+-`. The deck lesson now uses "The deck transformation group" as its first section and the KP
shard matches it. The final dump found no remaining `+-`, `--number`, `\cdot-`, `1 x`, `1 units`, `1 zeros`,
literal unicode or pure-number math-span pattern in the 120 remaining families.

## Pool and parameterisation defects

Seventeen families initially produced fewer distinct questions than the Wave 7 design target.
The ranges were widened or the stem was made to print the instance. The affected concepts include
`abstract.group.free_product`, `abstract.rep.character`, `abstract.rep.maschke`,
`abstract.rep.schur`, `abstract.rep.induced`, `topo.filters`, `topo.paracompact`,
`topo.partition_unity`, `topo.uniform_space`, `topo.limit_point_compact`, `topo.continuum`,
`topo.homotopy_group`, `topo.fibration`, `topo.cofibration`, `topo.simplicial_homology`,
`topo.jones_polynomial`, `topo.seifert_surface`, `ml.reg.polynomial`, `ml.svm.soft_margin` and
`ml.optim.sgd`. Every family now clears the hard five-in-twelve floor and the ten-in-sixty design
floor. The first-ten families measured at least 11 distinct questions in 200 draws except the
intentionally finite pools recorded in the candidate ledger.

## Curriculum density

The first graph validation after authoring failed because the new nodes had no
`key_prerequisites` entries, taking coverage to 79.8% against the 80.0% floor. Fifty enrichment
entries were added, one per Wave 7 concept, using the concept's prerequisite list. The validator
then reported `OK: 1850 concepts across 23 domains, 19 reviewed label collisions`.

## Singleton impact

Wave 6 closed with 141 mastery units and 22 singleton units. Wave 7 also closes with 141 units and
22 singletons. No singleton unit was added or removed. The new concepts fit existing bands and did
not change the singleton set. The band width remains 3 and the grouping algorithm is unchanged.

## Validation

Wave 7 closed with `make validate`, `go test ./...`, `git diff --check`, the graph validator, the
mastery-unit validator, the lesson audit, the additive check, the counts check and the published
docs bundle check. The frozen additive baseline remains the 1,500-concept snapshot. No solution
schemas were added because abstract algebra, topology and machine learning are not learner
domains. The commit is local. Nothing was pushed or deployed.
