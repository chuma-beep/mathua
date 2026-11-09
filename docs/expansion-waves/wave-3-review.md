# Wave 3 review

Date: 2026-11-09. Scope: the 50 concepts added by Wave 3, from 1,600 to 1,650.

This file records the inspection result as a metric so later waves can be compared against it.
Passing gates is technical correctness. The counts below come from reading generated questions,
answers and explanations, not from the test suite.

## Defect metrics

| Metric | Wave 3 |
|---|---|
| Concepts inspected | 50 |
| Knowledge-point families inspected | 150 |
| Mathematical defects | 5 |
| Presentation defects | 19 |
| Generator coverage defects | 6 |
| Duplicate concepts replaced before acceptance | 0 |
| Defects fixed | All identified |
| Regression tests added | Yes |
| Singleton units before and after | 25 and 24 |
| Singleton units created by the wave | 1 |
| Singleton units removed by the wave | 2 |

## Mathematical defects

Five defects affected the answer or the well-posedness of the question.

1. `linalg.matrix.block` KP3 asked for the (1,1) entry of \(M^2\) but computed the
   (1,1) entry of \(A B\). The block-diagonal square uses \(A^2\), so the generator now
   multiplies \(A\) by itself. The test recomputes the entry from the two diagonal blocks.
2. `nt.adv.consecutive_sums` counted the one-term representation and returned length 1 as the
   shortest block. The question asks for two or more terms, so both helpers now start at
   \(k=2\). The count is the number of odd divisors greater than 1.
3. `graph.assignment` KP2 could ask for the row-1 column when two permutations tied for the
   optimum. The generator now redraws until the optimum is unique, so the requested column is
   well defined.
4. `graph.bellman_ford` KP2 asked about "the same directed graph" without printing the graph.
   The question now includes the vertex set and the edge list.
5. `graph.stable_matching` KP2 and KP3 asked about "the same instance" without printing the
   preference lists. Both questions now publish the men's and women's preferences.

## Presentation defects

Nineteen questions or explanations rendered a defensible answer in a way a learner should not
have to read. All were fixed before closeout.

1. `nt.adv.perfect_numbers` KP1 listed 28 and 496 under a bound below 28 and then said "the first
   two" when the count was 1.
2. `nt.basics.decimal_period` KP3 said that zero factors of 2 and 5 contribute a terminating
   prefix.
3. `linalg.matrix.block` KP2 rendered the trace as `-1+-1+-2+-1`.
4. `linalg.matrix.norm` KP3 rendered squared entries as `4^2+-1^2+...`.
5. `linalg.matrix.schur_complement` KP1 rendered `3--2\cdot-3/2`.
6. `linalg.matrix.hadamard` KP2 rendered `-2\cdot-2+...`.
7. `linalg.matrix.kronecker` KP2 rendered `2\cdot-2`.
8. `linalg.matrix.powers` KP1 rendered `-3\cdot-3+-1\cdot2`.
9. `linalg.matrix.outer_product` KP1 rendered `-1\cdot-2`.
10. `linalg.transform.householder` KP1 rendered `1\cdot-5`.
11. `linalg.transform.givens` KP2 and KP3 rendered `--12/13` inside a matrix and
    `15x--8y` in the explanation.
12. `linalg.matrix.permutation` KP2 said "There are 1" and KP3 said "1-th component".
13. `linalg.decomp.pseudoinverse` KP1 said "1-th diagonal entry".
14. `graph.bfs` KP2 said "There are 1".
15. `linalg.matrix.rank_factorization` KP2 ended the trace explanation with the stray token
    `dx`.
16. `linalg.matrix.vandermonde` printed a plain node list as a math span, which canonicalised
    into literal dollars.
17. `nt.adv.polygonal` printed a plain polygonal value as a math span, which canonicalised into
    literal dollars.
18. The DFS and BFS graph builder could repeat an undirected edge in the displayed edge list.
    The builder now deduplicates edges.
19. The Floyd-Warshall graph builder could display parallel directed edges with different
    weights. The builder now keeps the smallest weight for each ordered pair.

## Generator coverage defects

Six families initially failed the requirement that a new concept produce at least five distinct
questions per ordinal in twelve draws. They were fixed before the variety baseline was frozen.

1. `linalg.matrix.kronecker` KP1 asked only for the dimension of a fixed 2x2 Kronecker product.
2. `linalg.transform.givens` KP1 had only two variants because the matrix was fixed up to the
   sign of the sine.
3. `linalg.transform.householder` KP2 had four variants.
4. `linalg.transform.householder` KP3 had four variants.
5. `graph.bellman_ford` KP2 printed only the target vertex.
6. `graph.stable_matching` KP2 printed only the number of participants.

The fixed families now draw from larger parameter sets or print the instance they describe. The
new tests check the lower bound as well as the answer.

## Duplicate concepts

No duplicate concept was added. Every proposed id and normalized label was checked against every
`data/concepts/` file before authoring, and `scripts/validate_graph.go` reports the same 19
reviewed collisions as before the wave. The reviewed allowlist was not extended.

## Singleton discrepancy

Wave 3 changed the partition from 139 units over 1,600 concepts to 140 units over 1,650 concepts.
Singleton units fell from 25 to 24.

- The one new singleton is `linear_algebra.b5`, which contains only
  `linalg.transform.householder`.
- `graph_theory.b2`, which contained only `graph.laplacian`, is no longer a singleton because
  the new prerequisite edges changed its topological depth and it moved into a multi-concept
  band.
- `linear_algebra.b2`, which contained only `linalg.vector.dot`, is no longer a singleton for
  the same reason.

The grouping rule and band width are unchanged: `unit = (domain, floor(conceptDepth / bandWidth))`
at width 3. The net change is caused by the added concepts and their prerequisite edges changing
existing depths, not by a change to the algorithm. `data/mastery_units.report.md` owns the live
figures, and `data/mastery_units.json` is regenerated with `go run ./cmd/masteryunits -band 3`.

## Validation

Wave 3 closed with `make validate`, `go test ./...`, `git diff --check` and the documentation
sweep. The additive baseline remains frozen at the 1,500-concept snapshot, no exemption was added
and no learner-domain solution schema was required because all three domains are university
domains. The commit is local. Nothing was pushed or deployed.
