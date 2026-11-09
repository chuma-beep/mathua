# Wave 6 review

Date: 2026-11-09. Scope: the 50 concepts added by Wave 6, from 1,650 to 1,800.

This file records the inspection result as a metric so later waves can be compared against it.
Passing gates is technical correctness. The counts below come from reading generated questions,
answers and explanations, not from the test suite.

## Defect metrics

| Metric | Wave 6 |
|---|---|
| Concepts inspected | 50 |
| Knowledge-point families inspected | 150 |
| Mathematical answer defects | 2 |
| Framing and well-posedness defects | 9 |
| Presentation defects | 1 |
| Pool or parameterization defects | 5 |
| Duplicate concepts replaced before acceptance | 0 |
| Defects fixed | All identified |
| Regression tests added | Yes |
| Singleton units before and after | 23 and 22 |
| Singleton units created by the wave | 0 |
| Singleton units removed by the wave | 1 |

## Mathematical answer defects

Two defects produced a wrong answer.

1. `prob.conc.bernstein` KP1 scaled the Bernstein denominator by 12 but used
   \(6v+2\). The denominator is \((3v+2)/6\), so 12 times it is \(6v+4\). The
   generator and its independent test now use \(6v+4\).
2. `prob.info.entropy` KP1 reported the total count times the entropy instead of
   twice the entropy. The weighted sum is divided by the total before multiplying
   by 2. The generator and test now agree on the requested \(2H\).

## Framing and well-posedness defects

Nine families asked for a fact that had just been stated, asked for a theorem
name or prefactor, or forced an always-yes answer. They were replaced with
condition checks, constructions or hypothesis-sensitive computations.

1. `prob.conc.chernoff` used an artificial scale factor. The scale is gone and
   the three questions ask for \(8M(\ln 2)\), \(8M'(0)\) and 16 times the bound.
2. `prob.conc.hoeffding` KP2 asked for the prefactor 2. It now asks for the
   exponent after the sample size is doubled.
3. `prob.conc.bernstein` KP3 was always yes. It now asks for the numerator over
   6 after the variance is doubled.
4. `prob.conc.mcdiarmid` KP3 asked for the prefactor 2. It now asks for the
   factor by which the exponent changes when \(t\) doubles.
5. `prob.limit.lindeberg` KP1 stated the maximum and asked for it, and KP3 was
   always yes. The families now ask for the condition, the ratio and the
   failure condition.
6. `prob.limit.berry_esseen` KP1 stated rho and KP3 stated n. They now ask for
   the integer multiplier in rho and for \(n\sigma^3\).
7. `prob.limit.local_clt` KP1 stated the span and KP2 stated the variance. They
   now ask for the normalising exponent denominator and the variance of a sum.
8. `prob.process.branching` could draw a distribution with \(p_0=0\), making
   extinction trivial. KP2 now draws only from a pool with \(p_0>0\) and
   \(p_2>0\).
9. `discrete.complexity.sat` KP3 subtracted from a clause count and could return
   a negative value. It now asks for the maximum number of clauses satisfied by
   any assignment.

## Presentation defect

One question had a grammar defect. `discrete.combinatorics.rook` said "On an
5x5 board"; it now says "On a 5 by 5 board".

## Pool and parameterization defects

Four families had a pool too small for the Wave 6 design target, or a parameter
that made the family read as a single fixed question. They were widened.

1. `prob.conc.chernoff` had an artificial scale that made the stems awkward. The
   scale was removed and the family now has seven Bernoulli parameters.
2. `prob.info.entropy` had only four weighted distributions and was reduced
   further when the artificial scale was removed. The pool now has twelve
   distributions with lengths 2, 3, 4, 5 and 8.
3. `prob.process.branching` had a trivial extinction case. The KP2 pool now has
   five nontrivial distributions.
4. `prob.limit.local_clt` had two self-answering families. The rewritten
   families now vary with n and v.
5. `discrete.automata.halting` had a single-answer easy pool: every family
   answered 1. KP2 now asks for the composed reduction cost, so the easy pool
   has more than one answer.

## Accepted pool exceptions

`discrete.coding.linear` KP1 has exactly 16 distinct messages because the
Hamming [4,7] code has 16 messages. It is below the 20-in-200 design target and
above the 10-in-60 minimum. The finite pool is the code itself, not a sampler
defect.

The remaining families clear the 5-in-12 hard floor and the 10-in-60 minimum.
Some theorem-heavy families sit at nine distinct stems in 60 draws because their
integer parameter ranges are deliberately small; the review records that as a
quality risk rather than lowering the floor.

## Quality risks

Many theorem-heavy concepts use finite, polynomial or small-integer examples.
That is deliberate, because the alternative is theorem-name recall. The risk is
that a family reads as recognition rather than computation. The following
families remain the most recognition-shaped and should be watched in later
reviews: `prob.conc.hoeffding`, `prob.conc.bernstein`,
`prob.conc.mcdiarmid`, `prob.limit.lindeberg`, `prob.limit.berry_esseen` and
`prob.limit.local_clt`.

The first ten were approved after review. The remaining forty were authored
after that approval. All 150 families were sampled and inspected.

## Singleton discrepancy

Wave 6 changed the partition from 141 units over 1,750 concepts to 141 units
over 1,800 concepts. Singleton units fell from 23 to 22.

- No singleton was added.
- `discrete_math.b8`, which contained only `discrete.sequences.recurrence`, is
  no longer a singleton. The new discrete concepts and prerequisite edges
  changed its topological depth, so it moved into a multi-concept band.

The grouping rule and band width are unchanged at width 3. The change is caused
by the added concepts changing existing depths, not by a change to the
algorithm. `data/mastery_units.report.md` owns the live figures.

## Validation

Wave 6 closed with `make validate`, `go test ./...`, `git diff --check` and the
documentation sweep. The additive baseline remains frozen at the 1,500-concept
snapshot, no exemption was added and no learner-domain solution schema was
required because none of the four domains is a learner domain. The commit is
local. Nothing was pushed or deployed.
