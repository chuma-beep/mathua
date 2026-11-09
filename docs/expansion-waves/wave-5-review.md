# Wave 5 review

Date: 2026-11-09. Scope: the 50 concepts added by Wave 5, from 1,650 to 1,750.

This file records the inspection result as a metric so later waves can be compared against it.
Passing gates is technical correctness. The counts below come from reading generated questions,
answers and explanations, not from the test suite.

## Defect metrics

| Metric | Wave 5 |
|---|---:|
| Concepts inspected | 50 |
| Knowledge-point families inspected | 150 |
| Mathematical or framing defects | 5 |
| Lesson defects | 4 |
| Presentation defects | 24 |
| Duplicate concepts replaced before acceptance | 0 |
| Defects fixed | All identified |
| Regression tests added | Yes |
| Families below 10 distinct questions in 60 draws after fixes | 0 |
| Singleton units before and after | 24 and 23 |
| Singleton units created by the wave | 0 |
| Singleton units removed by the wave | 1 |

The first ten were audited separately before the remaining forty were authorized. That audit found
no wrong answers and no sign-convention defects. It found four lesson gaps, recorded below. All 150
families were then inspected through generated samples. The final sample dump contained no `+-`,
`--number`, `\cdot-`, unit-coefficient, literal unicode or over-escaped bare-number pattern.

## Mathematical and framing defects

Five defects affected the answer, the assumptions or the well-posedness of the question.

1. `calc.multi.neumann` stated the wrong compatibility condition for the one-dimensional Neumann
   problem. For \(u''=f\) on \([0,1]\), the condition is
   \(u'(1)-u'(0)=\int f\); with outward normals the homogeneous condition is
   \(u'(1)=u'(0)\). The family now uses that condition.
2. `calc.multi.dirichlet` asked for two harmonic functions with the same boundary data to have
   different interior values. The question now tests uniqueness: if one function equals \(c\) at a
   point, the other must equal \(c\) at that point.
3. `calc.vec.euler_lagrange` claimed every straight line solves the Euler-Lagrange equation for
   \(L=\frac12(y')^2+a y\). The equation is \(y''=a\), so a straight line is a solution exactly
   when \(a=0\). The answer is now \(1\) only in that case.
4. `calc.multi.fundamental_solution` used \(p=2-n\) for every dimension, including \(n=2\), where
   the fundamental solution is logarithmic. The generator now draws \(n=3,\ldots,12\).
5. `ra.seq.fourier_convergence` stated the exponent \(p\) and then asked for \(p\). It now asks for
   \(a_n n^p\), which is the coefficient, so the question is not tautological.

## Lesson defects

The independent audit of the first ten found four lesson gaps. All were fixed before the remaining
forty were authored.

1. `ra.seq.limsup` did not define the tail suprema that the limit superior uses.
2. `ra.seq.liminf` did not define the tail infima that the limit inferior uses.
3. `ra.measure.lebesgue` did not distinguish the measure of a set from the integral of a simple
   function supported on it.
4. `ra.measure.simple_integral` did not separate the support measure from the weighted sum that
   defines the integral.

## Presentation defects

Twenty-four families required presentation fixes. The defects were sign joins, unit coefficients,
over-escaped bare numbers, singular/plural agreement and self-contained-instance problems.

- Sign joins and unit coefficients: `calc.vec.general_stokes`, `calc.vec.helmholtz`,
  `calc.vec.gauge`, `calc.vec.vector_potential`, `calc.vec.closed_exact`,
  `calc.multi.harmonic`, `calc.multi.mean_value`, `calc.multi.maximum_principle`,
  `calc.vec.euler_lagrange`, `calc.vec.first_variation` and `calc.multi.minimal_surface`.
- Over-escaped bare numbers: `calc.vec.geodesic`, `calc.multi.surface_measure`,
  `calc.multi.fourier` and `ra.measure.measurable_functions`.
- Signed factors and squares: `ra.metric.lp_space`, `ra.metric.holder`, `ra.metric.hilbert`,
  `ra.metric.riesz_representation`, `ra.deriv.absolute_continuity` and
  `ra.deriv.lebesgue_differentiation`.
- Grammar or self-contained instances: `ra.measure.measurable_sets` said "there are 1" and
  `ra.measure.fubini_tonelli` referred to "the same constant function" without stating it.
- One remaining sign join: `ra.metric.bounded_operator`.

A small `w5Clean` text sanitizer now strips math delimiters around pure integers, normalizes
`+-` and `--`, removes unit coefficients and collapses repeated spaces for the remaining
calculus and real-analysis families. The help is textual only; it does not alter answers or facts.
The final dump found no remaining instance of these patterns.

## Duplicate concepts

No duplicate concept was added. The ledger checked every id and normalized label against all 23
domains, and `scripts/validate_graph.go` reports the same 19 reviewed collisions as before the
wave. The reviewed allowlist was not extended.

## Singleton discrepancy

Wave 5 changed the partition from 141 units and 24 singletons over 1,700 concepts to 141 units and
23 singletons over 1,750 concepts.

- No singleton was added.
- `real_analysis.b13`, which previously contained only `ra.deriv.riemann_stieltjes`, is no longer
  a singleton. It now contains nine concepts at depths 39 to 41: `ra.deriv.riemann_stieltjes`,
  `ra.deriv.absolute_continuity`, `ra.measure.fatou`, `ra.metric.banach`,
  `ra.metric.bounded_operator`, `ra.metric.hilbert`, `ra.metric.holder`,
  `ra.metric.minkowski` and `ra.metric.riesz_representation`.

The new concepts and prerequisite edges changed the topological depth of existing real-analysis
concepts, which moved `ra.deriv.riemann_stieltjes` into a populated band. The rule and band width
are unchanged: `unit = (domain, floor(conceptDepth / bandWidth))` at width 3.
`data/mastery_units.report.md` owns the live figures.

## Quality risks

Many theorem-heavy concepts are implemented through finite, polynomial or small-metric examples.
That keeps every answer exact and independently checkable, but it means some families test a
consequence of the theorem rather than a full proof. Three families are deliberately recognized as
thinner than the best work in the wave: the Helmholtz divergence-free yes/no question, the
Holder and Minkowski scaling questions, and the bounded-operator always-yes question. They meet the
family and pool gates, but later work should deepen them before they are treated as model examples.

## Validation

The wave closed with `make validate`, `go test ./...` and `git diff --check` passing. The graph,
flow, documentation, counts and published-docs artifacts were regenerated at 1,750 concepts. The
additive baseline remains frozen at the 1,500-concept snapshot, no exemption was added and no
solution schema was required because Calculus III and real analysis are not learner domains. The
commit is local. Nothing was pushed or deployed.
