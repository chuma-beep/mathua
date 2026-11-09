# Wave 7 candidate ledger

Status: approved, authored and closed on 2026-11-09. The first ten concepts passed independent review, and the remaining forty were authored after that approval.

Date: 2026-11-09. Current corpus: 1,800 concepts. Frozen additive baseline: 1,500 concepts. Wave 7
target: 1,850 concepts.

This ledger fixes the 50 Wave 7 candidates, their subdomains, prerequisite edges and three-family
plans. Every proposed id and normalized label was checked against every file in `data/concepts/`
on 2026-11-09. There is no exact id collision and no normalized-label collision.

The split closes abstract algebra and topology and moves machine learning halfway. Course targets
are `abstract`, `topo` and `ml`. No solution schema is required because none of the three domains is
a learner domain.

## Allocation

| Domain | Tier | Now | Target | Add | Course |
|---|---|---:|---:|---:|---|
| abstract_algebra | university | 56 | 76 | +20 | `abstract`, `abstract2` |
| topology | university | 52 | 72 | +20 | `topo` |
| machine_learning | university | 55 | 75 | +10 | `ml` |
| **Total** | | | | **50** | |

The remaining 150 concepts after this wave are school 88, bridge 45 and university 17. Wave 8 can
carry a full school batch. Wave 9 can finish the bridge and take five school concepts. Wave 10 can
close the remaining school concepts, the last ten machine-learning concepts and graph theory.

## Abstract algebra: 20 candidates

| id | Label | Subdomain | Prerequisite |
|---|---|---|---|
| abstract.group.semidirect | Semidirect products | abstract_algebra.groups | abstract.group.direct_product, abstract.group.automorphism |
| abstract.group.wreath | Wreath products | abstract_algebra.groups | abstract.group.semidirect |
| abstract.group.free_product | Free products | abstract_algebra.groups | abstract.group.free |
| abstract.group.extension | Group extensions | abstract_algebra.groups | abstract.group.quotient_group, abstract.group.semidirect |
| abstract.rep.representation | Group representations | abstract_algebra.representations | abstract.group.homomorphism, abstract.structures.module |
| abstract.rep.character | Characters of group representations | abstract_algebra.representations | abstract.rep.representation |
| abstract.rep.maschke | Maschke's theorem | abstract_algebra.representations | abstract.rep.representation, abstract.field.characteristic |
| abstract.rep.schur | Schur's lemma | abstract_algebra.representations | abstract.rep.representation |
| abstract.rep.induced | Induced representations | abstract_algebra.representations | abstract.rep.character |
| abstract.poly.symmetric | Symmetric polynomials | abstract_algebra.rings | abstract.rings.polynomial |
| abstract.poly.resultant | Resultants | abstract_algebra.rings | abstract.rings.polynomial |
| abstract.poly.discriminant | Discriminants of polynomials | abstract_algebra.rings | abstract.poly.resultant |
| abstract.poly.eisenstein | Eisenstein's criterion | abstract_algebra.rings | abstract.rings.irreducibility |
| abstract.rings.fractional_ideal | Fractional ideals | abstract_algebra.rings | abstract.rings.ideal |
| abstract.rings.dedekind | Dedekind domains | abstract_algebra.rings | abstract.rings.fractional_ideal, abstract.rings.noetherian |
| abstract.rings.valuation | Valuation rings | abstract_algebra.rings | abstract.rings.integral_domain |
| abstract.rings.primary | Primary decomposition | abstract_algebra.rings | abstract.rings.ideal, abstract.rings.prime_ideal |
| abstract.structures.nakayama | Nakayama's lemma | abstract_algebra.structures | abstract.structures.fg_module |
| abstract.structures.ext | Ext and Tor | abstract_algebra.structures | abstract.structures.exact_sequence, abstract.structures.module_hom |
| abstract.structures.exterior | Exterior algebra | abstract_algebra.structures | abstract.structures.tensor |

### Abstract algebra family plans

| id | KP1 | KP2 | KP3 |
|---|---|---|---|
| abstract.group.semidirect | order of a semidirect product | multiply two normal forms | test whether an action is compatible |
| abstract.group.wreath | order of a wreath product | size of the base group | permutation action on coordinates |
| abstract.group.free_product | reduced word length | abelianisation rank | reduce a word to normal form |
| abstract.group.extension | order of an extension | identify the quotient | cocycle condition check |
| abstract.rep.representation | degree of a representation | character at the identity | decompose a direct sum |
| abstract.rep.character | character value on a conjugacy class | inner product of characters | orthogonality check |
| abstract.rep.maschke | characteristic condition | semisimple decomposition | averaging projection |
| abstract.rep.schur | dimension of the endomorphism algebra | scalar condition | count intertwiners |
| abstract.rep.induced | dimension of an induced representation | character at the identity | Frobenius reciprocity count |
| abstract.poly.symmetric | evaluate an elementary symmetric polynomial | convert a power sum | degree of a symmetric expression |
| abstract.poly.resultant | resultant of two small polynomials | common-root test | product formula |
| abstract.poly.discriminant | discriminant of a quadratic or cubic | repeated-root test | sign of the discriminant |
| abstract.poly.eisenstein | prime divisibility checks | irreducibility decision | shifted Eisenstein test |
| abstract.rings.fractional_ideal | inverse of a fractional ideal | product of two fractional ideals | principal-ideal check |
| abstract.rings.dedekind | factorisation exponents | unique factorisation check | ideal-class calculation |
| abstract.rings.valuation | valuation of a product | valuation-ring membership | uniformiser |
| abstract.rings.primary | primary-ideal test | radical of a primary ideal | number of components |
| abstract.structures.nakayama | hypothesis check | minimal number of generators | quotient dimension |
| abstract.structures.ext | Ext^1 dimension | Tor dimension | exact-sequence rank |
| abstract.structures.exterior | dimension of a wedge space | square-zero check | basis count |

## Topology: 20 candidates

| id | Label | Subdomain | Prerequisite |
|---|---|---|---|
| topo.nets | Nets and convergence | topology.basics | topo.basics.spaces |
| topo.filters | Filters and ultrafilters | topology.basics | topo.nets |
| topo.paracompact | Paracompact spaces | topology.separation | topo.separation.normal |
| topo.partition_unity | Partitions of unity | topology.separation | topo.paracompact |
| topo.function_space | The compact-open topology | topology.constructions | topo.basics.continuous |
| topo.uniform_space | Uniform spaces | topology.basics | topo.basics.metric |
| topo.sequential_compact | Sequential compactness | topology.compact | topo.compact |
| topo.limit_point_compact | Limit point compactness | topology.compact | topo.compact |
| topo.continuum | Continua | topology.connected | topo.connected |
| topo.homotopy_group | Higher homotopy groups | topology.algebraic | topo.homotopy |
| topo.fibration | Fibrations | topology.algebraic | topo.homotopy |
| topo.cofibration | Cofibrations | topology.algebraic | topo.homotopy |
| topo.deck | Deck transformations | topology.covering | topo.covering |
| topo.simplicial_homology | Simplicial homology | topology.algebraic | topo.homology.cellular |
| topo.singular_homology | Singular homology | topology.algebraic | topo.simplicial_homology |
| topo.cup_product | The cup product | topology.algebraic | topo.cohomology |
| topo.poincare_duality | Poincare duality | topology.manifolds | topo.cohomology, topo.manifold.orientability |
| topo.jones_polynomial | The Jones polynomial | topology.knot | topo.knot |
| topo.seifert_surface | Seifert surfaces | topology.knot | topo.knot |
| topo.brouwer | The Brouwer fixed point theorem | topology.algebraic | topo.homotopy |

### Topology family plans

| id | KP1 | KP2 | KP3 |
|---|---|---|---|
| topo.nets | net convergence | cluster point | subnet construction |
| topo.filters | filter membership | ultrafilter extension | convergence test |
| topo.paracompact | open-cover refinement count | local finiteness | partition existence |
| topo.partition_unity | subordinate cover | sum of values | support check |
| topo.function_space | compact-open subbasis | evaluation continuity | convergence test |
| topo.uniform_space | entourage membership | uniform continuity | completion |
| topo.sequential_compact | convergent subsequence | sequence count | counterexample |
| topo.limit_point_compact | limit point of an infinite set | count of limit points | equivalence condition |
| topo.continuum | connected compact Hausdorff check | decomposition | cut point |
| topo.homotopy_group | pi_n of a sphere | abelian condition | suspension |
| topo.fibration | homotopy lifting | fiber | long exact sequence |
| topo.cofibration | homotopy extension | mapping cylinder | pushout |
| topo.deck | deck-group order | action check | quotient |
| topo.simplicial_homology | boundary map | Betti number | chain rank |
| topo.singular_homology | H_0 | H_1 of a circle | chain-group rank |
| topo.cup_product | degree | graded commutativity | cohomology ring |
| topo.poincare_duality | Betti pairing | dimension | cap product |
| topo.jones_polynomial | value at 1 | skein relation | crossing change |
| topo.seifert_surface | genus | Euler characteristic | boundary |
| topo.brouwer | fixed-point existence | dimension | counterexample without compactness |

## Machine learning: 10 candidates

| id | Label | Subdomain | Prerequisite |
|---|---|---|---|
| ml.reg.lasso | Lasso regression | machine_learning.regression | ml.linreg.least_squares, ml.reg.ridge_penalty |
| ml.reg.elastic_net | Elastic net | machine_learning.regression | ml.reg.lasso, ml.reg.ridge_penalty |
| ml.reg.polynomial | Polynomial regression | machine_learning.regression | ml.linreg.least_squares |
| ml.tree.gradient_boosting | Gradient boosting | machine_learning.trees | ml.tree.ensemble |
| ml.svm.soft_margin | Soft-margin support vector machines | machine_learning.classification | ml.svm.margin |
| ml.cluster.dbscan | DBSCAN | machine_learning.unsupervised | ml.cluster.kmeans |
| ml.cluster.hierarchical | Hierarchical clustering | machine_learning.unsupervised | ml.cluster.kmeans |
| ml.unsupervised.autoencoder | Autoencoders | machine_learning.neural_networks | ml.backpropagation |
| ml.optim.sgd | Stochastic gradient descent | machine_learning.neural_networks | ml.gradient_descent |
| ml.data.leakage | Data leakage | machine_learning.data | ml.validation, ml.data.scaling |

### Machine learning family plans

| id | KP1 | KP2 | KP3 |
|---|---|---|---|
| ml.reg.lasso | L1 penalty value | sparsity threshold | objective value |
| ml.reg.elastic_net | mixing parameter | coefficient value | combined penalty |
| ml.reg.polynomial | design-matrix entry | degree choice | prediction |
| ml.tree.gradient_boosting | residual | learning-rate update | additive prediction |
| ml.svm.soft_margin | slack sum | margin value | effect of C |
| ml.cluster.dbscan | core-point test | cluster count | noise count |
| ml.cluster.hierarchical | merge distance | dendrogram height | cluster count |
| ml.unsupervised.autoencoder | bottleneck dimension | reconstruction loss | parameter count |
| ml.optim.sgd | one update step | minibatch gradient | variance of the step |
| ml.data.leakage | leakage decision | contaminated feature | train-test separation |

## Overlap audit

The following candidates were rejected or deferred for duplication or for a weak independent
mastery objective.

- Abstract algebra: Burnside's lemma duplicates the existing group-action lesson; Pólya enumeration
  is a combinatorics application rather than a new algebra concept; group cohomology was deferred
  because it would duplicate Ext and Tor at this level; algebraic number theory was deferred because
  the number-theory domain already owns class numbers and ideals of integers.
- Topology: Stone-Cech, one-point compactification, Urysohn, Tietze, van Kampen, Mayer-Vietoris,
  Betti numbers, Euler characteristic, genus and surface classification already exist. The new
  candidates add different methods or different invariants.
- Machine learning: Gaussian processes were rejected because `prob.process.gaussian` already owns
  the label and the concept; AUC was deferred because ROC already exists; quantile regression was
  rejected because Wave 4 added `stat.reg.quantile`; causal inference was deferred for a later
  applied wave; polynomial regression is distinct from the existing linear and ridge lessons.

Every candidate has one observable mastery objective. Every prerequisite exists in the current DAG
or is an earlier candidate in this ledger. No existing prerequisite edge is removed.

## Review order

Author and inspect ten first: four abstract algebra (`semidirect`, `representation`, `symmetric`,
`ext`), three topology (`nets`, `sequential_compact`, `deck`), two machine learning (`lasso`,
`gradient_boosting`) and one more topology or machine-learning concept chosen by the reviewer
(`cup_product` is the suggested tenth).

Run the domain tests, graph validator, lesson audit and shipping contract after the ten. Check
the 5-in-12 hard floor, the 10-in-60 minimum and the 20-in-200 design target. Only then author the
remaining forty.

## Artifact contract

Each concept needs a numeric node in `data/concepts/<domain>.json`; a generator in
`internal/generator/<domain>/wave7.go` registered from the package `Register`; an independent
fact-based `wave7_test.go`; an authored lesson in `data/lessons/authored/<id>.md` with one H1 and
exactly three H2 sections over 1 KB; a three-entry KP shard; a `data/lessons/lessons.json` row; and
a target in `abstract`, `abstract2`, `topo` or `ml`.

No solution schemas are required because none of the three domains is a learner domain.

## Quality gates

- Three section-aligned families per concept.
- Hard floor of 5 distinct questions per ordinal in 12 draws.
- Design target of 20 distinct questions per ordinal in 200 draws.
- Minimum accepted pool of 10 distinct questions per ordinal in 60 draws or a written reason.
- Finite pools report their size.
- Every numeric answer independently recomputed from generator facts.
- At least one family per theorem-heavy concept tests a condition, counterexample, construction or
  hypothesis-sensitive computation.
- No always-yes family and no theorem-name recall.
- Manual inspection of all 150 families, with mathematical and presentation defects counted
  separately.
- Full gate: `make validate`, `go test ./...` and `git diff --check`.
- Regenerate mastery units, graph, flow, docs and variety baseline.
- Keep the additive baseline frozen at 1,500 concepts and keep band width at 3.
- Commit locally only. Do not push or deploy.
