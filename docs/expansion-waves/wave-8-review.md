# Wave 8 review

Date: 2026-11-10. Scope: the 50 concepts added by Wave 8, from 1,850 to 1,900.

This file records the inspection result as a metric so later waves can be compared against it.
Passing gates is technical correctness. The counts below come from reading generated questions,
answers and explanations, not from the test suite alone.

## Defect metrics

| Metric | Wave 8 |
|---|---|
| Concepts inspected | 50 |
| Knowledge-point families inspected | 150 |
| Mathematical or framing defects | 4 |
| Presentation defects | 10 |
| Pool or parameterisation defects | 12 |
| Duplicate concepts replaced before acceptance | 0 |
| Defects fixed | All identified |
| Regression tests added | Yes |
| Singleton units before and after | 22 and 23 |
| Singleton units created by the wave | 1 |
| Singleton units removed by the wave | 0 |

## Mathematical and framing defects

Four defects affected the mathematical framing of a family. None survived into the final
generator.

1. `trig.spherical.law_sines` initially multiplied and divided the wrong sine numerators in the
   draft, so the reduced ratio was not the spherical sine rule. The generator now reduces
   \(\sin b=\sin B\sin a/\sin A\) and the test recomputes the reduced numerator and denominator.
2. `trig.spherical.law_cosines` initially paired cosine inputs with unrelated result values. The
   case table now stores the cosine, sine and result components together and every question
   prints the instance it uses.
3. `trig.adv.machin` initially printed a minus sign for Machin-like identities that add their
   two arctangent terms. The table now records the sign and the stem prints the correct
   expression.
4. `arith.ops.compensation` could ask how much to subtract when the overshoot was negative. The
   delta range is now 7 through 10, so the compensation is always a nonnegative subtraction.

No wrong final answer was found in the 150 families.

## Presentation defects

Ten questions or explanations rendered a defensible answer in a way a learner should not have to
read. All were fixed before closeout.

1. `trig.basics.versine` KP3 rendered an over-escaped bare denominator as literal dollars.
2. `frac.div.common_denominator` said "divide both parts by 1" when the quotient was already in
   lowest terms.
3. `arith.place.roman` KP2 said "there are 1" when the numeral contained one X.
4. `arith.place.tally` KP2 said "1 marks" when one mark remained.
5. `frac.add.number_line` said "divide by 1" when the sum was already reduced.
6. `frac.sub.number_line` had the same already-reduced wording.
7. `frac.mult.area_model` had the same already-reduced wording.
8. `frac.div.missing_factor` said "which simplifies by 1" when no simplification was needed.
9. `arith.factor.triangular_numbers`, `frac.mixed.to_decimal` and `trig.adv.machin` each wrapped
   a bare number or decimal in math delimiters, producing literal-dollar warnings in the latex
   check.
10. `arith.mult.lattice` KP3 asked for the product without framing it as lattice multiplication.

## Pool and parameterisation defects

Twelve families began below the Wave 8 design target of 20 distinct questions in 200 draws.
They were widened before the baseline was frozen.

- Versine and harmonic addition moved from 8 and 6 variants to 24 Pythagorean-angle variants.
- Chebyshev first kind moved from 5 to 20 values of \(n\).
- Making ten KP1 is a deliberate finite pool of 10 choices for the addend.
- Haversine was widened from 5 rational special angles to 9, with a scale factor giving 18
  stems per family.
- Gudermannian moved from 8 to 19 values of \(k\).
- Weierstrass moved from 6 to 20 values of \(t\).
- Chebyshev second kind moved from 8 to 20 values of \(n\).
- Machin moved from 3 formula variants to 3 formulas with a scale factor, giving 12 stems.
- Spherical sines and spherical cosines moved from 5 and 4 cases to 10 cases each.
- The two-step remaining problem no longer has a constant KP3 stem.

Every family clears the hard floor of 5 distinct questions per ordinal in 12 draws and the
minimum of 10 in 60 draws.

## Duplicate concepts

No duplicate concept was added. Every id and normalized label was checked against every
`data/concepts/` file before authoring, and `scripts/validate_graph.go` reports the same 19
reviewed collisions as before the wave. The reviewed allowlist was not extended.

The overlap audit rejected generic fraction word problems, ordinary multiplication drill
variants, planar angle-bisector and Ceva material, and sum-of-arctangent material already
implied by complex arguments.

## Singleton discrepancy

Wave 8 changed the partition from 141 units over 1,850 concepts to 142 units over 1,900
concepts. Singleton units moved from 22 to 23.

- The only new singleton is `arithmetic.b4`, containing only `arith.factor.gcf_lcm_word`.
- No singleton was removed.

The grouping rule and band width are unchanged: `unit = (domain, floor(conceptDepth / bandWidth))`
at width 3. The new unit is caused by the added concept landing in a new depth band.

## Curriculum and artifacts

The split is trigonometry 78 to 98, arithmetic 72 to 90 and fractions 52 to 64. Course targets
are `pc` and `a2` for trigonometry, `pa` and the grade courses for arithmetic and fractions.
All three domains are learner domains, so all 50 concepts have authored solution schemas. No
schema exemption was added.

The corpus now has 1,900 concepts across 23 domains, 1,900 lessons, 1,900 KP files and 5,700
knowledge points. All 1,900 concepts are reachable through a course.

## Validation

Wave 8 closed with `make validate`, `go test ./...`, `git diff --check`, the graph validator,
the lesson audit, the solution audit, the additive check, the documentation checks and the
published-docs bundle. The additive baseline remains frozen at the 1,500-concept snapshot and
the mastery band width remains 3.

Known risks recorded for later review: the law-of-tangents and Napier families share a
right-triangle formula, the making-ten KP1 pool is deliberately finite and the haversine family
uses only rational special angles. These are recorded as limits, not accepted defects.

The commit is local. Nothing was pushed or deployed.
