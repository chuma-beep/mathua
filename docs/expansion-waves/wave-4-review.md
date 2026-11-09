# Wave 4 review

Date: 2026-11-09. Scope: the 50 concepts added by Wave 4, from 1,650 to 1,700.

This file records the inspection result as a metric so later waves can be compared against it.
Passing gates is technical correctness. The counts below come from reading generated questions,
answers and explanations, not from the test suite.

## Defect metrics

| Metric | Wave 4 |
|---|---:|
| Concepts inspected | 50 |
| Knowledge-point families inspected | 150 |
| Mathematical and framing defects | 4 |
| Presentation defects | 8 |
| Pool and parameterization defects | 6 |
| Duplicate concepts replaced before acceptance | 0 |
| Defects fixed | All identified |
| Regression tests added | Yes |
| Singleton units before and after | 23 and 24 |
| Singleton units created by the wave | 1 |
| Singleton units removed by the wave | 0 |

## Approved split

Differential equations gains 20 concepts, from 46 to 66. Complex numbers gains 20, from 51 to 71.
Statistics gains 10, from 87 to 97. The courses are `diffeq`, `cx` and `probstat`. None of the
three domains is a learner domain, so no authored solution schema was required.

## Mathematical and framing defects

Four questions asked for a quantity that was already stated, or stated an order that did not match
the computation.

1. `complex.adv.trig_zeros` KP2 said that \(\sin(kz)\) still has period \(2\pi\) and answered 2.
   It now asks for \(k\) times the period, which is 2 units of \(\pi\), with a correct
   explanation.
2. `complex.adv.schwarz_christoffel` KP1 asked for the number of vertices after stating the number
   of vertices. It now asks for the interior-angle sum in units of \(\pi\), answer \(n-2\), and
   the test was updated.
3. `ode.numerical.rkf` KP1 stated fourth and fifth order but asked for the lower order. It now
   fixes the facts at 4 and 5 and asks for their sum, answer 9.
4. `ode.transforms.periodic` KP1 asked for the period after stating it. It now asks for the length
   of three periods, answer \(3T\), and the test was updated.

## Presentation defects

Eight questions or explanations rendered a defensible answer in a way a learner should not have
to read.

1. `complex.adv.roots_unity` KP3 printed literal `\u03c6`; it now says `phi(n)`.
2. `complex.adv.schwarz_christoffel` used "an 7-gon" and "an 3-gon"; it now says "a polygon with
   n sides" or equivalent.
3. `complex.adv.jensen` KP1 and KP3 said "1 zeros"; it now uses the singular "zero" when the count
   is 1.
4. `complex.adv.weierstrass_factorization` KP2 said "1 zeros"; it now uses the singular "zero".
5. `ode.first_order.riccati` KP1 printed "1 y^2"; it now omits the unit coefficient.
6. `ode.transforms.periodic` KP2 and KP3 printed "1 units"; it now uses a singular/plural helper.
7. `ode.transforms.initial_final` KP3 had an awkward explanation of the condition \(1>0\); it now
   says the pole is negative and the theorem applies.
8. `stat.test.mann_whitney` KP2 omitted the samples; it now includes both samples and the pooled
   rank sum.

After the final sample dump no remaining `+-`, `--number`, `\cdot-`, `1 units`, `1 zeros`,
literal unicode escape, `an N-gon` or LaTeX over-escape pattern was found in the 150 samples.

## Pool and parameterization defects

Six families were too narrow or published a fact that did not describe the question.

1. `stat.infer.bootstrap` KP3 published `variance` while the question gave a standard deviation.
   The fact is now `sd`.
2. `stat.test.permutation` could produce a non-integer difference of means. It now redraws until
   both group sums divide by the group size.
3. `stat.test.mann_whitney` could produce tied pooled values. It now redraws until all pooled
   values are distinct.
4. `stat.reg.quantile` KP3 did not include the line and the input in the question. It now does.
5. `stat.data.missing` KP2 could produce a non-integer mean. It now uses an even step.
6. `complex.adv.casorati` had a single-answer easy pool because all three families answered 1.
   KP1 now asks for the residue coefficient \(a\), which varies and satisfies the pool ratchet.

Earlier ODE and complex inspection also widened several parameter ranges and replaced constant
families with parameterized ones. Those changes are part of the same pool work and are covered by
the three `wave4_test.go` diversity tests.

## Duplicate concepts

No duplicate concept was added. Every id and normalized label was checked against every
`data/concepts/` file before authoring. The graph validator reports the same 19 reviewed label
collisions as before the wave. The reviewed allowlist was not extended.

## Singleton discrepancy

Wave 4 changed the partition from 140 units over 1,650 concepts to 141 units over 1,700 concepts.
Singleton units rose from 23 to 24.

- The one new singleton is `differential_equations.b15`, which contains only
  `ode.integral.fredholm`.
- No singleton unit was removed.

The grouping rule and band width are unchanged: `unit = (domain, floor(conceptDepth / bandWidth))`
at width 3. The new unit is caused by the added Fredholm concept landing in a new depth band, not
by a change to the algorithm. `data/mastery_units.report.md` owns the live figures, and
`data/mastery_units.json` is regenerated with `go run ./cmd/masteryunits -band 3`.

## Validation

Wave 4 closed with `make validate`, `go test ./...`, `git diff --check` and the documentation
sweep. The additive baseline remains frozen at the 1,500-concept snapshot, no exemption was added
and no review allowlist was extended. The commit is local. Nothing was pushed or deployed.
