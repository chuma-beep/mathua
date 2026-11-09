# Wave 6 proposal

Status: approved and authoring complete on 2026-11-09. The allocation is approved. The candidate ledger is approved. The first ten passed review and the remaining forty were authored after that approval.

Date: 2026-11-09. Current corpus: 1,750 concepts. Frozen additive baseline: 1,500 concepts. Wave 6
target: 1,800 concepts.

This document proposes the domain split, candidate themes, prerequisite strategy and quality gates
for Wave 6. The exact ids and labels belong to a candidate ledger, which is authored only after the
allocation is approved.

## Remaining allocation

The remaining 250 concepts split into 88 school, 46 bridge and 116 university concepts.

| Domain | Now | Target | Add | Tier | Course |
|---|---:|---:|---:|---|---|
| discrete_math | 64 | 87 | +23 | university | `discrete` |
| probability_theory | 63 | 85 | +22 | university | `prob` |
| statistics | 97 | 98 | +1 | bridge | `probstat` |
| graph_theory | 59 | 70 | +4 | university | `graph` |
| **Total** | | | **50** | | |

This split completes discrete mathematics, probability theory and statistics, and moves graph
theory to 63 of its 70-concept target. It leaves 200 concepts for Waves 7 through 10.

## Why this split

- **Discrete mathematics and probability theory are the two largest remaining rows, but the
  recommendation does not rest on size alone.** Both domains have many advanced topics that are
  exact, finite and independently checkable. They support three genuinely different question
  families without falling back on theorem-name recall.
- **Statistics needs one concept, not a token batch.** One well-chosen capstone closes the row
  instead of leaving a lone concept for a later wave. The candidate must add a real statistical
  method, not repeat an existing test under a new label.
- **Graph theory gets four concepts rather than a full batch.** The domain is already dense and
  overlaps discrete mathematics. Four carefully chosen structural or algorithmic additions can
  fill real gaps without repeating graph material already taught in the discrete course.
- **The split avoids the highest duplication risks.** Precalculus is a consolidation domain and
  already overlaps algebra and trigonometry. Abstract algebra and topology are mathematically
  important but theorem-heavy, so they are better served by a later wave with a dedicated
  computational design rather than by a rushed allocation here.
- **The school rows can wait without harming prerequisites.** All proposed Wave 6 concepts depend
  on existing university material, so no learner is blocked by the later school expansion.

## Alternatives considered

| Option | Split | Strength | Risk |
|---|---|---|---|
| Applied data wave | probability +22, machine learning +20, statistics +1, graph +7 | completes three applied rows | probability and machine learning overlap statistics and linear algebra |
| Pure advanced wave | abstract algebra +20, topology +20, graph +10 | closes two difficult rows | theorem-heavy, high recognition risk and weaker generator feasibility |
| School foundations wave | trigonometry +20, arithmetic +18, fractions +12 | balances tiers and has easy generators | high duplication risk in already dense school domains |
| Bridge completion wave | precalculus +23, calc1 +14, calc2 +8, statistics +1, graph +4 | closes four bridge rows | precalculus overlaps algebra and trigonometry heavily |

The recommended split is the best balance of mathematical importance, exact generator feasibility,
course completion and overlap risk. It also leaves the difficult theorem-heavy domains for a wave
where the generator design can be built around finite models rather than definitions.

## Candidate clusters

### Discrete mathematics: 23 candidates

The domain already covers propositional logic, quantifiers, set operations, relations, counting,
permutations, combinations, Pascal, binomial identities, stars and bars, inclusion-exclusion,
derangements, Catalan numbers, Stirling numbers, Bell numbers, integer partitions, lattice paths,
generating functions, recurrences, master theorem, Boolean functions, Burnside, resolution, finite
state machines, Hamming codes and several proof techniques.

The new candidates should extend the advanced spine rather than repeat those topics:

- automata and formal languages: nondeterministic automata, regular expressions, pumping lemma,
  context-free grammars, pushdown automata, Turing machines, decidability and the halting problem;
- computability and complexity: reductions, NP-completeness, satisfiability, Cook-Levin, complexity
  classes and approximation;
- coding theory: linear codes, generator and parity-check matrices, Hamming distance, cyclic codes
  and Reed-Solomon codes;
- advanced combinatorics: rook polynomials, Young tableaux, q-binomial coefficients, Schur numbers
  and permutation statistics;
- algorithms and recurrences: dynamic programming, greedy correctness, divide-and-conquer analysis
  and amortised analysis;
- logic and proof: compactness, semantic tableaux, Horn clauses and first-order proof systems.

The ledger must reject any candidate whose lesson would mostly repeat an existing discrete lesson,
and any candidate that belongs more naturally to graph theory or number theory.

### Probability theory: 22 candidates

The domain already covers axioms, conditional probability, independence, random variables, joint
and marginal distributions, expectation, variance, covariance, generating functions, common
distributions, the Poisson process, Markov chains, random walks, martingales, optional stopping,
characteristic functions, convergence modes and the main limit theorems.

The new candidates should extend the stochastic-process and inference spine:

- stochastic processes: branching processes, renewal theory, queues, Gaussian processes, hidden
  Markov models and the Kalman filter;
- continuous-time processes: Brownian motion, Itô integrals, stochastic differential equations and
  the Fokker-Planck equation;
- concentration and large deviations: Chernoff bounds, Hoeffding, Bernstein, McDiarmid and the
  large-deviation principle;
- Bayesian and decision theory: conjugate priors, posterior predictive distributions, Bayes risk,
  minimax rules and admissible estimators;
- information theory: entropy, relative entropy, mutual information, channel capacity and source
  coding;
- limit refinements: Lindeberg and Lyapunov conditions, Berry-Esseen and the local central limit
  theorem.

Every concept must publish a finite distribution, exact moment, finite Markov chain or closed-form
bound. A family that only asks a learner to name a theorem is not acceptable.

### Statistics: 1 candidate

The single concept should close a real gap rather than duplicate an existing test. Candidate areas
are multiple-comparison correction, causal inference with potential outcomes, or a regression
diagnostic that is not already covered. The ledger must compare the candidate against every
existing statistics label before choosing it.

### Graph theory: 4 candidates

The domain already covers isomorphism, complements, subgraphs, homomorphisms, graph products,
distance, cuts, connectivity, matchings, coloring, planarity, spectral graph theory, extremal
graph theory, flows and shortest paths. The four additions should be one coherent gap, not four
unrelated novelties.

Candidate clusters:

- structural graph classes: chordal, interval, split, cograph or permutation graphs, with
  recognition and invariant questions;
- graph width and decomposition: pathwidth, branchwidth or clique-width, with small examples;
- random and extremal refinements: expansion, graph limits or regularity, only if the generator can
  ask exact finite questions;
- graph algorithms: centrality, network reliability or dynamic connectivity, provided they do not
  duplicate the existing shortest-path, flow and matching material.

## Prerequisite strategy

- Every new concept needs at least one existing or same-wave prerequisite.
- No existing prerequisite edge may be removed.
- The only permitted change to an existing concept is an added follow-on prerequisite when an
  isolated new concept needs a real successor.
- Every candidate must be checked for cycles, missing ids and course reachability.
- Cross-domain edges need a written reason. The ledger should prefer an existing edge over a new
  one when both express the same dependency.
- The ledger must show the depth and band effect of the new edges and identify every new singleton
  before authoring.

## Generator and review gates

- Three genuinely different section-aligned families per concept.
- Hard floor: five distinct valid questions per ordinal in twelve draws.
- Wave 6 target: twenty distinct questions per ordinal in two hundred draws.
- Minimum accepted pool: ten distinct questions per ordinal in sixty draws, or a written reason.
- Every numeric answer recomputed independently from published facts.
- Every theorem-heavy concept needs at least one family that distinguishes understanding from
  recognition: a condition check, a counterexample, a construction or a computation that changes
  when the hypothesis fails.
- No family may be answerable by naming a theorem. No always-yes family is acceptable.
- Finite pools must report their size.
- Manual inspection of all 150 families, with mathematical and presentation defects counted
  separately.

## Approval sequence

1. Approve the allocation in this proposal.
2. Author the candidate ledger with exactly 50 ids, labels, subdomains, prerequisites and
   three-family plans, plus the overlap audit.
3. Author the first ten review concepts only.
4. Run the domain tests, lesson audit, shipping contract and an independent mathematical audit.
5. Author the remaining forty only after the first ten are approved.
6. Run the full gate, regenerate mastery, graph, flow, docs and variety artifacts, record the
   review and commit locally.

Reaching 1,800 is a corpus milestone, not automatically a learning-quality milestone. A candidate
that cannot support three honest families is dropped before authoring, not shipped weak.

## Validation expectations

The full wave must pass `make validate`, `go test ./...` and `git diff --check`; regenerate mastery
units, graph, flow, documentation and variety artifacts; preserve the frozen 1,500-concept additive
baseline; add no solution schemas because these domains are not learner domains; keep band width at
3; and leave `internal/server/server.go` and `internal/server/google_login_host_test.go` unstaged.
No push or deployment is part of this wave.
