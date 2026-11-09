# Wave 3 follow-up review

Date: 2026-11-09. Status: complete. Scope: the pre-Wave-4 review requested after the Wave 3
closeout. No concept was added. The corpus remains at 1,650.

The review covered four questions: the readability impact of the 19 presentation defects, the
residual variety of the six repaired generator families, the isolation of
`linalg.transform.householder`, and the progression created by the new prerequisite edges.

## Presentation defects

The 19 defects split into three classes.

| Class | Count | Examples |
|---|---:|---|
| Substantive readability | 3 | `nt.adv.perfect_numbers` described numbers outside the bound; `nt.basics.decimal_period` described zero factors as a terminating prefix; `linalg.matrix.rank_factorization` ended with the stray token `dx` |
| Minor grammar | 3 | `There are 1`, `1-th component` and `1-th diagonal entry` |
| Formatting only | 13 | `+-` and `--` joins, `\cdot-` products, over-escaped plain numbers and repeated graph edges |

The three substantive cases changed what a reader understood, not only how the expression looked.
They are fixed. The three grammar cases were readable but unpolished. The remaining thirteen were
typographic: the answer and the mathematics were correct before and after the repair. No unresolved
presentation defect now changes mathematical meaning.

## Generator coverage

The six repaired families were sampled for 200 draws at each ordinal. The count is distinct
questions, not distinct option orders.

| Concept | KP1 | KP2 | KP3 |
|---|---:|---:|---:|
| `linalg.matrix.kronecker` | 200 | 200 | 200 |
| `linalg.transform.givens` | 28 | 28 | 195 |
| `linalg.transform.householder` | 158 | 24 | 24 |
| `graph.bellman_ford` | 200 | 200 | 200 |
| `graph.stable_matching` | 200 | 200 | 200 |

The minimum repaired supply is 24 distinct questions at the Householder determinant and symmetry
families. That is well above the required five, so the six named coverage defects leave no
contract-level gap. The Givens and Householder families are still finite pools rather than fully
parametric families. They are the thinnest Wave 3 supplies and should not be treated as unbounded
in a later variety audit.

## Householder isolation

Before this review, `linalg.transform.householder` had one prerequisite,
`linalg.orth.projection_matrix`, and no dependents. Its concept depth was 15, and no other linear
algebra concept had depth 15, so it formed the singleton unit `linear_algebra.b5`.

That isolation was real, but the concept is not redundant. The existing QR lesson already says
that Householder reflections compute the same factorisation more stably than classical
Gram-Schmidt. The correction was to make that relation explicit:

`linalg.decomp.qr` now has `linalg.transform.householder` as a prerequisite.

The effect is narrow. Householder keeps depth 15. QR moves from depth 14 to depth 16. The band-5
unit now contains `linalg.decomp.qr` and `linalg.transform.householder`. No other concept changes
band and no new singleton is created. The partition remains 140 units, and singleton units fall
from 24 to 23.

The resulting order is defensible: the learner meets projection, then the reflection that performs
a stable orthogonal transformation, then the QR decomposition that can use it. The QR lesson
already carried the Householder explanation, so this is not a new topic edge invented for the
partition.

## Prerequisite progression

All 50 new concepts have at least one existing prerequisite. No prerequisite target is missing.
The new edges point to established foundations in the same domain: modular arithmetic, divisor
functions, matrix multiplication, projections and matchings.

Only the existing concept `linalg.vector.dot` changed depth range, from depth 8 to depth 6, and it
stayed in band 2. Ten later linear algebra concepts changed their depth maximum from 28 to 29 and
stayed in band 9. No existing concept crossed a band because of Wave 3. The QR edge above changes
only QR, and it was added by this review rather than by the original wave.

One redundancy is worth recording. `nt.adv.prime_counting` lists both `nt.adv.sieve` and
`nt.analytic.prime_number_theorem` as prerequisites, but `nt.adv.sieve` already depends on the
prime number theorem. The direct theorem edge does not change the depth: the effective order is
prime number theorem at depth 14, sieve at depth 15 and prime counting at depth 16. It is a
redundant edge, not a reversed progression in the DAG. If elementary prime counting is to move
earlier, Wave 4 needs a separate elementary counting concept rather than a rewiring of this one.

## Decision and status

The direct prerequisite edge from `nt.adv.prime_counting` to
`nt.analytic.prime_number_theorem` is retained. It is transitively redundant through
`nt.adv.sieve`, but it is intentional, explicit and does not affect the depth or band assignment
of `nt.adv.prime_counting` or any other concept.

The existing concept is not rewired to introduce elementary prime counting. That work is deferred
as a separate potential concept for Wave 4, where it can be proposed with its own id, lesson,
generator and prerequisite edge. It is not part of this follow-up.

The follow-up review is closed. The live partition is 140 units and 23 singleton units, with the
Householder singleton resolved by the QR prerequisite edge.

## Verification

The follow-up changes are a prerequisite edge, two generator parameter sets, regenerated mastery
artifacts and documentation. No concept, lesson mapping or course target changed. The full gate is
rerun after the edits. The commit is local. Nothing is pushed or deployed, and the unrelated
server files remain unstaged.
