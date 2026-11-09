# Wave 6 candidate ledger

Status: approved and authoring complete on 2026-11-09. The allocation is approved. The
first ten passed review and the remaining forty were authored after that approval.

Date: 2026-11-09. Current corpus: 1,750 concepts. Frozen additive baseline: 1,500 concepts. Wave 6
target: 1,800 concepts.

This ledger fixes the 50 Wave 6 candidates, their subdomains, prerequisite edges and three-family
plans. Every proposed id and normalized label was checked against every file in `data/concepts/`
with the normalization rule in `internal/concepts/labels.go`. The check found zero exact id
collisions and zero normalized-label collisions.

## Allocation

| Domain | Tier | Now | Target | Add | Course |
|---|---|---:|---:|---:|---|
| discrete_math | university | 64 | 87 | +23 | `discrete` |
| probability_theory | university | 63 | 85 | +22 | `prob` |
| statistics | bridge | 97 | 98 | +1 | `probstat` |
| graph_theory | university | 59 | 70 | +4 | `graph` |
| **Total** | | | | **50** | |

The wave completes discrete mathematics, probability theory and statistics. Graph theory moves to
63 of its 70-concept target. It leaves 200 concepts for Waves 7 through 10.

Four new subdomains are introduced because no existing subdomain names the material:
`discrete_math.algorithms`, `discrete_math.complexity`, `probability_theory.concentration` and
`probability_theory.information`. Each is used by at least two candidates. No existing subdomain
is renamed or removed.

## Curriculum overlap audit

The audit applied the five checks required by the proposal: duplicate identity, genuine
advancement, prerequisite honesty, independent mastery objective and family fit.

- The 50 labels are new under normalized comparison. No candidate is a renamed copy of an existing
  concept.
- `discrete.coding.linear` generalizes the existing Hamming-code lesson into linear codes,
  generator matrices and syndrome decoding. It is not a second Hamming-code lesson.
- `discrete.automata.nfa` and `discrete.automata.regex` extend the existing finite-state-machine
  concept into nondeterminism and regular expressions. They do not repeat the FSM lesson.
- `prob.info.entropy` and `prob.info.kl` add an information-theory spine that is absent from the
  existing moment, limit and process material.
- `prob.process.brownian` is not a relabelled random walk. It asks for continuous-time Gaussian
  structure, covariance and hitting probabilities.
- `stat.infer.multiple_comparisons` is not a second ANOVA or p-value lesson. It tests family-wise
  error and false-discovery control.
- `graph.chordal`, `graph.interval`, `graph.split` and `graph.cograph` are structural graph classes,
  not another treatment of perfect graphs. The existing `graph.perfect_graph` lesson supplies the
  prerequisite, not the answer.
- The audit rejected or deferred compactness, semantic tableaux, Horn clauses, first-order proof
  systems, Lyapunov conditions, large deviations, the Fokker-Planck equation, channel capacity,
  source coding, permutation graphs, pathwidth, branchwidth, clique-width, centrality, network
  reliability, dynamic connectivity, causal inference and regression diagnostics. The reasons are
  duplicate coverage, insufficient slot count or a generator that would test theorem names rather
  than a computation.
- The first prerequisite check found two invalid ids. `discrete.combinatorics.rook` now uses
  `discrete.sets.inclusion_exclusion`, not a nonexistent combinatorics id. `prob.process.renewal`
  and `prob.process.queue` now use `prob.dist.poisson_process`, not a nonexistent process id.

## Discrete mathematics: 23 candidates

All subdomains are existing `discrete_math.*` unless noted. The new `discrete_math.algorithms` and
`discrete_math.complexity` subdomains are introduced by this wave.

| id | Label | Subdomain | Prerequisite | KP1 / KP2 / KP3 |
|---|---|---|---|---|
| discrete.automata.nfa | Nondeterministic finite automata | discrete_math.automata | discrete.automata.fsm | subset-construction state count / accepted string / epsilon closure |
| discrete.automata.regex | Regular expressions and finite automata | discrete_math.automata | discrete.automata.fsm | regex from a DFA / language membership / derivative or size |
| discrete.automata.pumping_lemma | The pumping lemma for regular languages | discrete_math.automata | discrete.automata.nfa | pumping length / decomposition / prove nonregular |
| discrete.automata.cfg | Context-free grammars | discrete_math.automata | discrete.automata.nfa | derive string / parse-tree count / language membership |
| discrete.automata.pda | Pushdown automata | discrete_math.automata | discrete.automata.cfg | stack operation / acceptance / state-stack trace |
| discrete.automata.turing | Turing machines | discrete_math.automata | discrete.automata.pda | transition / tape after steps / accept or reject |
| discrete.automata.decidability | Decidable and undecidable languages | discrete_math.automata | discrete.automata.turing | classify language / reduction direction / closure property |
| discrete.automata.halting | The halting problem | discrete_math.automata | discrete.automata.turing | diagonal argument / reduction / restricted decidable case |
| discrete.complexity.reduction | Polynomial-time reductions | discrete_math.complexity | discrete.automata.turing | map instance / preserve yes or no / compose reductions |
| discrete.complexity.sat | Boolean satisfiability | discrete_math.complexity | discrete.logic.boolean_functions | truth assignment / clause satisfaction count / 3-SAT instance |
| discrete.complexity.np_complete | NP-complete problems | discrete_math.complexity | discrete.complexity.reduction, discrete.complexity.sat | certificate check / reduction direction / class membership |
| discrete.complexity.cook_levin | The Cook-Levin theorem | discrete_math.complexity | discrete.complexity.np_complete | tableau size / clause count / reduction gadget |
| discrete.coding.linear | Linear codes | discrete_math.coding | discrete.coding.hamming | encode message / Hamming distance / correct error |
| discrete.coding.generator | Generator and parity-check matrices | discrete_math.coding | discrete.coding.linear | codeword from matrix / nullspace or parity check / rank and dimension |
| discrete.coding.cyclic | Cyclic codes | discrete_math.coding | discrete.coding.linear | shift / polynomial representation / generator polynomial degree |
| discrete.coding.reed_solomon | Reed-Solomon codes | discrete_math.coding | discrete.coding.linear | evaluation codeword / distance / correct errors |
| discrete.combinatorics.rook | Rook polynomials | discrete_math.combinatorics | discrete.sets.inclusion_exclusion | rook polynomial of a small board / forbidden positions / inclusion-exclusion count |
| discrete.combinatorics.young_tableaux | Standard Young tableaux | discrete_math.combinatorics | discrete.combinatorics.integer_partitions | shape or standard tableau / hook length / count tableaux |
| discrete.combinatorics.q_binomial | Gaussian binomial coefficients | discrete_math.combinatorics | discrete.combinatorics.integer_partitions, discrete.combinatorics.combinations | evaluate at q / lattice-path interpretation / recurrence |
| discrete.combinatorics.permutation_statistics | Inversions and major index | discrete_math.combinatorics | discrete.combinatorics.permutations | inversion count / major index / distribution over S_n |
| discrete.combinatorics.schur | Schur numbers | discrete_math.combinatorics | discrete.combinatorics.pigeonhole | Schur number bound / coloring / small case |
| discrete.algorithms.dynamic_programming | Dynamic programming | discrete_math.algorithms | discrete.sequences.recurrence | recurrence / table value / optimal substructure |
| discrete.algorithms.greedy | Greedy algorithms and exchange arguments | discrete_math.algorithms | discrete.proof.cases | exchange argument / greedy choice / counterexample |

## Probability theory: 22 candidates

| id | Label | Subdomain | Prerequisite | KP1 / KP2 / KP3 |
|---|---|---|---|---|
| prob.process.branching | Branching processes | probability_theory.processes | prob.process.markov_chain | offspring mean / extinction probability / generation distribution |
| prob.process.renewal | Renewal theory | probability_theory.processes | prob.dist.poisson_process | renewal function / interarrival mean / asymptotic rate |
| prob.process.queue | Queueing processes | probability_theory.processes | prob.dist.poisson_process, prob.process.markov_chain | traffic intensity / stationary distribution / expected queue length |
| prob.process.brownian | Brownian motion | probability_theory.processes | prob.limit.clt, prob.process.random_walk | mean and variance / covariance / hitting probability |
| prob.process.ito | Ito integrals | probability_theory.processes | prob.process.brownian | integral mean / quadratic variation / Ito formula |
| prob.process.sde | Stochastic differential equations | probability_theory.processes | prob.process.ito | drift and diffusion / Euler step / solution mean |
| prob.process.gaussian | Gaussian processes | probability_theory.processes | prob.rv.multivariate_transform | covariance kernel / finite-dimensional distribution / conditional mean |
| prob.process.hidden_markov | Hidden Markov models | probability_theory.processes | prob.process.markov_chain | transition and emission / forward probability / Viterbi path |
| prob.process.kalman | The Kalman filter | probability_theory.processes | prob.process.hidden_markov, prob.process.gaussian | predict step / update gain / posterior mean |
| prob.conc.chernoff | Chernoff bounds | probability_theory.concentration | prob.rv.mgf | MGF bound / optimize exponent / tail bound |
| prob.conc.hoeffding | Hoeffding's inequality | probability_theory.concentration | prob.conc.chernoff | range bound / sample-mean deviation / confidence |
| prob.conc.bernstein | Bernstein's inequality | probability_theory.concentration | prob.conc.hoeffding, prob.rv.variance | variance term / bound / compare Hoeffding |
| prob.conc.mcdiarmid | McDiarmid's inequality | probability_theory.concentration | prob.conc.hoeffding | bounded difference / bound / function of independent variables |
| prob.limit.lindeberg | The Lindeberg condition | probability_theory.limits | prob.limit.clt | condition check / triangular array / CLT conclusion |
| prob.limit.berry_esseen | The Berry-Esseen theorem | probability_theory.limits | prob.limit.clt | third moment / rate / sample size |
| prob.limit.local_clt | The local central limit theorem | probability_theory.limits | prob.limit.clt | lattice correction / point probability / approximation |
| prob.bayes.conjugate | Conjugate priors | probability_theory.foundations | prob.bayes.prior | posterior parameters / posterior mean / predictive |
| prob.bayes.predictive | Posterior predictive distributions | probability_theory.foundations | prob.bayes.conjugate | posterior predictive / mean / probability |
| prob.decision.bayes_risk | Bayes risk | probability_theory.foundations | prob.expect.conditional | risk of estimator / Bayes estimator / compare |
| prob.decision.minimax | Minimax decision rules | probability_theory.foundations | prob.decision.bayes_risk | minimax value / least favorable prior / compare Bayes risk |
| prob.info.entropy | Entropy of a random variable | probability_theory.information | prob.rv.expectation | entropy / joint entropy / maximum entropy |
| prob.info.kl | Relative entropy and mutual information | probability_theory.information | prob.info.entropy | KL divergence / mutual information / chain rule |

## Statistics: 1 candidate

| id | Label | Subdomain | Prerequisite | KP1 / KP2 / KP3 |
|---|---|---|---|---|
| stat.infer.multiple_comparisons | Multiple comparisons and false discovery rate | statistics.inference | stat.infer.pvalue, stat.infer.hypothesis | Bonferroni threshold / family-wise error / Benjamini-Hochberg |

## Graph theory: 4 candidates

| id | Label | Subdomain | Prerequisite | KP1 / KP2 / KP3 |
|---|---|---|---|---|
| graph.chordal | Chordal graphs | graph_theory.basics | graph.perfect_graph | perfect elimination ordering / chordless cycle / clique tree |
| graph.interval | Interval graphs | graph_theory.basics | graph.chordal | interval representation / clique number / recognition |
| graph.split | Split graphs | graph_theory.basics | graph.chordal, graph.complement | partition / degree sequence / forbidden induced subgraph |
| graph.cograph | Cographs | graph_theory.basics | graph.complement | cotree / complement / P4-free test |

## Review order and approval sequence

The first ten candidates are authored and reviewed before the remaining forty. The ten cover all
four domains and both computational and theorem-sensitive material.

1. `discrete.automata.nfa`
2. `discrete.automata.regex`
3. `discrete.coding.linear`
4. `discrete.algorithms.dynamic_programming`
5. `prob.process.branching`
6. `prob.process.brownian`
7. `prob.conc.chernoff`
8. `stat.infer.multiple_comparisons`
9. `graph.chordal`
10. `graph.cograph`

For the ten, run the domain tests, lesson audit, shipping contract and an independent mathematical
audit. A family that is answerable by naming a theorem, or that always answers yes, is rejected.
Only after the ten pass review are the remaining forty authored.

## Artifact contract

Every concept ships with:

- a numeric node in the domain's `data/concepts/<domain>.json`;
- a generator in `internal/generator/<domain>/wave6.go`, registered from that package's `Register`;
- an independent `wave6_test.go` that recomputes every answer from Facts;
- a lesson in `data/lessons/authored/<id>.md` with one H1 and exactly three H2 sections over 1 KB;
- a three-entry shard in `data/lessons/kp/<id>.json`;
- a row in `data/lessons/lessons.json`;
- a target in `data/courses.json`, in `discrete`, `prob`, `probstat` or `graph`.

No solution schema is required because none of the four domains is a learner domain. No existing
prerequisite is removed. The band width stays 3. The frozen additive baseline stays at the 1,500
concept snapshot.

## Quality gates

- Three section-aligned question families per concept.
- Hard floor: five distinct questions per ordinal in twelve draws.
- Wave 6 target: twenty distinct questions per ordinal in two hundred draws.
- Minimum accepted pool: ten distinct questions per ordinal in sixty draws, or a written reason.
- At least one family per theorem-heavy concept must test a condition, counterexample, construction
  or hypothesis-sensitive computation. No always-yes family and no theorem-name recall.
- Finite pools must report their size.
- Every numeric answer is recomputed independently from Facts.
- Manual inspection of all 150 families, with mathematical and presentation defects counted
  separately in `docs/expansion-waves/wave-6-review.md`.
- Full gate: `make validate`, `go test ./...`, `git diff --check`; regenerate mastery, graph, flow,
  docs and variety artifacts; commit locally only; no push or deployment.

The corpus milestone is 1,800, not a substitute for the quality gate. A candidate that cannot meet
the family or pool standard is dropped before authoring.
