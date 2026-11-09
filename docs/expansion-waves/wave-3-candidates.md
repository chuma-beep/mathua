# Wave 3 candidate ledger

Wave 3 raises the corpus from 1,600 to 1,650 concepts. This ledger fixes the 50 candidates and
their provisional prerequisite edges before authoring. Every id and label below was checked against
every file in `data/concepts/` on 2026-11-09 with the same normalization rule as
`internal/concepts/labels.go`: there is no exact id collision and no normalized-label collision.

The split finishes two target rows exactly and opens a third. Number theory gains 20 concepts,
from 51 to 71, which is its whole remaining allocation. Linear algebra gains 21, from 59 to 80,
which is also its whole remaining allocation. Graph theory gains 9, from 50 to 59, leaving 11 of
its 20 for a later wave. The three domains are university domains. The six learner domains that
require authored solution schemas are arithmetic, fractions, prealgebra, algebra, geometry and
trigonometry, so this wave adds no solution schema.

The three domains were chosen because their remaining targets are large and their additions are
exact computation rather than a second treatment of algebra, precalculus or calculus. No candidate
repeats an existing lesson on modular inverses, perfect numbers, Mersenne primes, twin primes,
matrix powers or traversal, so the wave adds coverage rather than a synonym.

Course targets are `nt`, `linalg` and `graph`. Every new concept is added to the target list of its
domain's course, and every target resolves through the DAG.

## Number theory: 20 candidates

| id | Label | Subdomain | Prerequisite |
|---|---|---|---|
| nt.basics.mod_power | Modular exponentiation | number_theory.basics | nt.basics.modular |
| nt.basics.mod_inverse | Modular inverses | number_theory.basics | nt.basics.bezout |
| nt.basics.euler_theorem | Euler's theorem | number_theory.basics | nt.adv.euler_phi, nt.basics.mod_power |
| nt.basics.decimal_period | Repeating decimal periods | number_theory.basics | nt.basics.mod_inverse |
| nt.basics.digital_root | Digital roots and casting out nines | number_theory.basics | nt.basics.divisibility |
| nt.adv.aliquot | Aliquot sums | number_theory.advanced | nt.divisor.sum |
| nt.adv.perfect_numbers | Perfect, abundant and deficient numbers | number_theory.advanced | nt.adv.aliquot |
| nt.adv.mersenne | Mersenne numbers and primes | number_theory.advanced | nt.basics.prime_inf |
| nt.adv.fermat_numbers | Fermat numbers | number_theory.advanced | nt.basics.mod_power |
| nt.adv.repunit | Repunits | number_theory.advanced | nt.basics.modular |
| nt.adv.prime_counting | The prime-counting function | number_theory.advanced | nt.adv.sieve, nt.analytic.prime_number_theorem |
| nt.adv.twin_primes | Twin primes | number_theory.advanced | nt.adv.prime_gap |
| nt.adv.goldbach | Goldbach partitions | number_theory.advanced | nt.basics.prime_inf |
| nt.adv.polygonal | Polygonal numbers | number_theory.advanced | nt.basics.divisibility |
| nt.adv.pentagonal | Pentagonal numbers | number_theory.advanced | nt.adv.polygonal |
| nt.adv.egyptian | Egyptian fractions | number_theory.advanced | nt.basics.gcd_euclidean |
| nt.adv.highly_composite | Highly composite numbers | number_theory.advanced | nt.divisor.count |
| nt.adv.amicable | Amicable numbers | number_theory.advanced | nt.adv.aliquot |
| nt.adv.happy | Happy numbers | number_theory.advanced | nt.basics.digital_root |
| nt.adv.consecutive_sums | Sums of consecutive integers | number_theory.advanced | nt.basics.divisibility |

## Linear algebra: 21 candidates

| id | Label | Subdomain | Prerequisite |
|---|---|---|---|
| linalg.matrix.powers | Matrix powers | linear_algebra.matrices | linalg.matrix.mult |
| linalg.matrix.nilpotent | Nilpotent matrices | linear_algebra.matrices | linalg.matrix.powers |
| linalg.matrix.idempotent | Idempotent matrices | linear_algebra.matrices | linalg.matrix.powers |
| linalg.matrix.elementary | Elementary matrices | linear_algebra.matrices | linalg.systems.row_reduction |
| linalg.matrix.block | Block matrices | linear_algebra.matrices | linalg.matrix.mult |
| linalg.matrix.permutation | Permutation matrices | linear_algebra.matrices | linalg.det.properties, linalg.matrix.identity |
| linalg.matrix.norm | Matrix norms | linear_algebra.matrices | linalg.vec.norm_squared |
| linalg.decomp.pseudoinverse | Moore-Penrose pseudoinverse | linear_algebra.decompositions | linalg.decomp.svd |
| linalg.decomp.cholesky | Cholesky decomposition | linear_algebra.decompositions | linalg.quadratic.definite |
| linalg.matrix.schur_complement | Schur complements | linear_algebra.matrices | linalg.matrix.block |
| linalg.matrix.trace_identities | Trace identities | linear_algebra.matrices | linalg.eigen.trace |
| linalg.matrix.rank_factorization | Rank factorization | linear_algebra.matrices | linalg.matrix.rank |
| linalg.matrix.outer_product | Outer products | linear_algebra.matrices | linalg.matrix.mult |
| linalg.matrix.hadamard | Hadamard products | linear_algebra.matrices | linalg.matrix.add |
| linalg.matrix.kronecker | Kronecker products | linear_algebra.matrices | linalg.matrix.mult |
| linalg.matrix.vandermonde | Vandermonde matrices | linear_algebra.matrices | linalg.det.3x3 |
| linalg.matrix.companion | Companion matrices | linear_algebra.matrices | linalg.eigen.cayley_hamilton |
| linalg.matrix.circulant | Circulant matrices | linear_algebra.matrices | linalg.eigen.compute |
| linalg.transform.householder | Householder reflections | linear_algebra.transformations | linalg.orth.projection_matrix |
| linalg.transform.givens | Givens rotations | linear_algebra.transformations | linalg.transform.det_scale |
| linalg.matrix.determinant_lemma | The matrix determinant lemma | linear_algebra.matrices | linalg.det.properties |

## Graph theory: 9 candidates

| id | Label | Subdomain | Prerequisite |
|---|---|---|---|
| graph.dfs | Depth-first search | graph_theory.connectivity | graph.adjacency |
| graph.bfs | Breadth-first search | graph_theory.connectivity | graph.adjacency |
| graph.bellman_ford | Bellman-Ford shortest paths | graph_theory.flows | graph.shortest_path |
| graph.floyd_warshall | All-pairs shortest paths | graph_theory.flows | graph.shortest_path |
| graph.union_find | Disjoint-set union | graph_theory.connectivity | graph.connectivity |
| graph.assignment | The assignment problem | graph_theory.matching | graph.matching |
| graph.stable_matching | Stable matching | graph_theory.matching | graph.matching |
| graph.clique | Clique number | graph_theory.matching | graph.independent_set |
| graph.dominating_set | Domination number | graph_theory.matching | graph.vertex_cover |

## Review order

The wave is authored in three batches so a schema or generator defect is found before it is copied
forty times.

1. Four number-theory concepts: modular exponentiation, modular inverses, Euler's theorem and
   repeating decimal periods.
2. Three linear-algebra concepts: matrix powers, nilpotent matrices and idempotent matrices.
3. Three graph concepts: depth-first search, breadth-first search and disjoint-set union.
4. The remaining 40 concepts, in the table order above.

Each batch runs `go test ./internal/generator/...`, `go run scripts/validate_graph.go`,
`python3 scripts/audit_lessons.py` and the new domain tests before the next batch begins.

## Artifact contract

Every concept ships with all of the following.

- A node in `data/concepts/<domain>.json` with a numeric grading type and a mastery threshold.
- A generator in `internal/generator/<domain>/wave3.go`, registered from that package's
  `Register` function.
- A generator test in `internal/generator/<domain>/wave3_test.go` that recomputes every numeric
  answer from the facts the generator publishes, independently of the answer string.
- A lesson in `data/lessons/authored/<concept>.md` with one H1 title and exactly three H2 sections,
  each section over the lesson corpus threshold.
- A knowledge-point shard in `data/lessons/kp/<concept>.json` with exactly three entries whose
  `section` fields match the lesson headings.
- A `{concept_id, source}` row in `data/lessons/lessons.json`.
- A target entry in `data/courses.json`.

The wave also regenerates the mastery units, the graph and flow artifacts, the documented counts
and the variety baseline, then runs `make validate`, `go test ./...` and `git diff --check`.
