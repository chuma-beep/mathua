# Wave 8 candidate ledger

Status: approved and authoring complete on 2026-11-10. The first ten passed review and the remaining forty were authored after that approval.

Wave 8 raises the corpus from 1,850 to 1,900 concepts. This ledger fixes the 50 candidates and their three-family plans before authoring. Every id and label was checked against every file in `data/concepts/` on 2026-11-10 with the same normalization rule as `internal/concepts/labels.go`: there is no exact id collision and no normalized-label collision.

The split is trigonometry +20, arithmetic +18 and fractions +12. All three are school domains. Course wiring is `pc` for trigonometry, `pa` for arithmetic and fractions, and grade courses 4 to 6 for the arithmetic and fraction concepts that fit their level. Every candidate is reachable through at least `pc` or `pa`.

## Overlap audit

The five checks in the Wave 8 proposal were run.

1. **Duplicate identity.** No id or normalized-label collision.
2. **Genuine advancement.** The arithmetic candidates are mental strategies, alternative number systems and representation models that the existing 72 arithmetic concepts do not teach. The fraction candidates are representations and one-step reasoning moves that the existing 52 fraction concepts do not teach. The trigonometry candidates extend triangle trigonometry into spherical geometry, orthogonal polynomials, harmonic addition and phasors; none repeats the existing law of sines, law of cosines, Heron, sum and difference or hyperbolic material.
3. **Prerequisite honesty.** Every prerequisite below exists in the current corpus. No candidate depends on a later Wave 8 candidate.
4. **Independent mastery objective.** Each candidate has a computation or construction that the existing corpus does not measure. Candidates that only renamed an existing lesson were rejected.
5. **Family fit.** Each concept has three distinct families in the tables below. Every family asks for an integer or an exact rational value and publishes the operands in `Facts`.

Rejected or deferred overlaps include angle bisector theorem, Stewart's theorem, Ceva and Menelaus in the plane, sum of arctangents already implied by complex arguments, generic fraction word problems and generic multiplication drill variants.

## Trigonometry: 20 candidates

Primary target `pc`; secondary target `a2`.

| id | Label | Subdomain | Prerequisite | KP1 / KP2 / KP3 |
|---|---|---|---|---|
| `trig.ident.law_tangents` | The law of tangents | trigonometry.identities | `trig.adv.law_sines`, `trig.adv.law_cosines` | half-sum tangent / side ratio / angle |
| `trig.ident.mollweide` | Mollweide's formulas | trigonometry.identities | `trig.adv.law_sines`, `trig.adv.law_cosines` | side ratio / angle sum / check identity |
| `trig.adv.half_side` | The half-side formulas | trigonometry.advanced | `trig.adv.law_sines`, `trig.adv.law_cosines` | semiperimeter / half-angle / side |
| `trig.adv.napier` | Napier's analogies | trigonometry.advanced | `trig.adv.law_sines`, `trig.adv.law_cosines` | angle pair / side pair / ratio |
| `trig.spherical.law_sines` | The spherical law of sines | trigonometry.spherical | `trig.adv.law_sines` | side from angles / angle from sides / great-circle |
| `trig.spherical.law_cosines` | The spherical law of cosines | trigonometry.spherical | `trig.adv.law_cosines` | side / angle / distance |
| `trig.spherical.excess` | Spherical excess and area | trigonometry.spherical | `trig.spherical.law_cosines` | area / excess / radius |
| `trig.spherical.haversine` | The haversine formula | trigonometry.spherical | `trig.spherical.law_cosines`, `trig.basics.radians` | distance / central angle / radius |
| `trig.basics.versine` | Versine and coversine | trigonometry.basics | `trig.basics.sin_cos_def` | value / identity / inverse |
| `trig.basics.exsecant` | Exsecant and excosecant | trigonometry.basics | `trig.basics.reciprocal` | value / identity / domain |
| `trig.adv.gudermannian` | The Gudermannian function | trigonometry.advanced | `trig.hyperbolic.tanh_coth`, `trig.adv.arctan` | angle / hyperbolic value / inverse |
| `trig.adv.weierstrass` | Tangent half-angle substitution | trigonometry.advanced | `trig.ident.half_angle`, `trig.ident.sum_diff` | substitution / rational form / parameter |
| `trig.adv.chebyshev` | Chebyshev polynomials of the first kind | trigonometry.advanced | `trig.ident.double_angle`, `trig.ident.triple_angle` | value / recurrence / coefficient |
| `trig.adv.chebyshev_second` | Chebyshev polynomials of the second kind | trigonometry.advanced | `trig.adv.chebyshev`, `trig.ident.sum_to_product` | value / recurrence / coefficient |
| `trig.adv.ptolemy` | Ptolemy's theorem | trigonometry.advanced | `trig.adv.law_cosines`, `trig.ident.sum_diff` | diagonal / side / cyclic test |
| `trig.adv.brahmagupta` | Brahmagupta's formula | trigonometry.advanced | `trig.adv.heron_area`, `trig.adv.ptolemy` | area / side / cyclic test |
| `trig.adv.arctan_sum` | The arctangent addition formula | trigonometry.advanced | `trig.adv.arctan`, `trig.ident.sum_diff` | sum / difference / quadrant |
| `trig.adv.machin` | Machin's formula | trigonometry.advanced | `trig.adv.arctan_sum` | pi approximation / arctan / term count |
| `trig.adv.harmonic_addition` | Combining sine and cosine into one wave | trigonometry.advanced | `trig.ident.sum_diff`, `trig.graph.period` | amplitude / phase / value |
| `trig.adv.phasor` | Phasors and phase angle | trigonometry.advanced | `trig.adv.harmonic_addition`, `trig.graph.phase_shift` | magnitude / angle / sum |

## Arithmetic: 18 candidates

Primary target `pa`; grade-course targets are shown in the last column.

| id | Label | Subdomain | Prerequisite | KP1 / KP2 / KP3 | Course |
|---|---|---|---|---|---|
| `arith.place.binary` | Binary place value | arithmetic.place_value | `arith.place.tens` | value / place / conversion | 4, `pa` |
| `arith.place.base_five` | Base-five place value | arithmetic.place_value | `arith.place.tens` | value / place / conversion | 4, `pa` |
| `arith.place.roman` | Roman numerals | arithmetic.place_value | `arith.place.tens` | value / numeral / subtractive | 4, `pa` |
| `arith.place.tally` | Tally marks and unary counting | arithmetic.place_value | `arith.add.single` | count / group / value | 4, `pa` |
| `arith.ops.making_ten` | Making ten | arithmetic.addition | `arith.add.single` | missing addend / sum / strategy | 4, `pa` |
| `arith.ops.doubling_halving` | Doubling and halving | arithmetic.multiplication | `arith.mult.tables` | double / half / product | 4, `pa` |
| `arith.ops.compensation` | Compensation in addition and subtraction | arithmetic.addition | `arith.add.double` | adjust / result / difference | 4, `pa` |
| `arith.ops.near_doubles` | Near doubles | arithmetic.addition | `arith.add.single` | sum / double / difference | 4, `pa` |
| `arith.mult.partial_products` | Partial products | arithmetic.multiplication | `arith.mult.two_by_two` | products / sum / area | 5, `pa` |
| `arith.div.partial_quotients` | Partial quotients | arithmetic.division | `arith.div.long` | quotients / sum / remainder | 5, `pa` |
| `arith.mult.area_model` | The area model of multiplication | arithmetic.multiplication | `arith.mult.concept` | area / factor / product | 4, `pa` |
| `arith.mult.lattice` | Lattice multiplication | arithmetic.multiplication | `arith.mult.two_by_two` | diagonal sum / product / grid | 5, `pa` |
| `arith.factor.square_numbers` | Square numbers | arithmetic.factors | `arith.mult.tables` | square / root / pattern | 5, `pa` |
| `arith.factor.triangular_numbers` | Triangular numbers | arithmetic.factors | `arith.add.double` | triangular / index / sum | 5, `pa` |
| `arith.ops.factorial_notation` | Factorial notation | arithmetic.exponents | `arith.mult.tables` | factorial / ratio / expression | 6, `pa` |
| `arith.factor.gcf_lcm_word` | GCF and LCM word problems | arithmetic.factors | `arith.factor.gcf`, `arith.factor.lcm` | gcf / lcm / word | 6, `pa` |
| `arith.time.elapsed` | Elapsed time | arithmetic.word_problems | `arith.sub.double` | duration / end / start | 4, `pa` |
| `arith.money.change` | Making change | arithmetic.word_problems | `arith.sub.double`, `arith.dec.intro` | change / coins / total | 4, `pa` |

## Fractions: 12 candidates

Primary target `pa`.

| id | Label | Subdomain | Prerequisite | KP1 / KP2 / KP3 | Course |
|---|---|---|---|---|---|
| `frac.basics.set` | A fraction of a set | fractions.concept | `frac.of_number` | fraction of set / count / remainder | 5, `pa` |
| `frac.basics.improper_number_line` | Improper fractions on the number line | fractions.concept | `frac.basics.number_line`, `frac.mixed.convert` | position / convert / compare | 5, `pa` |
| `frac.compare.same_numerator` | Comparing fractions with the same numerator | fractions.equivalence | `frac.ops.compare` | compare / order / reason | 5, `pa` |
| `frac.order.three` | Ordering three fractions | fractions.equivalence | `frac.ops.compare` | order / compare / common denominator | 5, `pa` |
| `frac.add.number_line` | Adding fractions on the number line | fractions.addition | `frac.add.diff` | sum / missing / distance | 5, `pa` |
| `frac.sub.number_line` | Subtracting fractions on the number line | fractions.subtraction | `frac.sub.diff` | difference / missing / distance | 5, `pa` |
| `frac.mult.area_model` | Multiplying fractions with an area model | fractions.multiplication | `frac.ops.mult` | area / product / factor | 5, `pa` |
| `frac.div.common_denominator` | Dividing fractions with a common denominator | fractions.division | `frac.ops.div` | quotient / missing factor / compare | 6, `pa` |
| `frac.div.missing_factor` | Division as a missing factor | fractions.division | `frac.ops.div` | missing factor / quotient / equation | 6, `pa` |
| `frac.mixed.to_decimal` | Mixed numbers as decimals | fractions.mixed | `frac.ops.to_decimal`, `frac.mixed.convert` | decimal / fraction / compare | 6, `pa` |
| `frac.percent.convert` | Fractions and percentages | fractions.equivalence | `frac.ops.to_decimal` | percent / fraction / compare | 6, `pa` |
| `frac.word.remaining_two_step` | Two-step fraction remaining problems | fractions.concept | `frac.remaining`, `frac.word.multi_step` | remaining / consumed / fraction | 6, `pa` |

## Review order

The first ten are reviewed before the remaining forty.

1. Trigonometry: `trig.basics.versine`, `trig.adv.chebyshev`, `trig.adv.harmonic_addition`, `trig.adv.arctan_sum`.
2. Arithmetic: `arith.place.binary`, `arith.ops.making_ten`, `arith.mult.partial_products`.
3. Fractions: `frac.basics.set`, `frac.compare.same_numerator`, `frac.div.common_denominator`.

The first ten cover a theorem-adjacent identity, a spherical computation, an orthogonal-polynomial recurrence, a graph-to-algebra conversion, two number systems, two mental models and three fraction representations. If they expose a schema or generator problem, fix the process before authoring the remaining forty.

## Artifact contract

Every concept ships with a numeric node in `data/concepts/<domain>.json`, a generator in `internal/generator/<domain>/wave8.go` registered from the package `Register`, an independent fact-based `wave8_test.go`, an authored lesson in `data/lessons/authored/<id>.md` with one H1 and exactly three H2 sections over 1 KB, a three-entry shard in `data/lessons/kp/<id>.json`, a row in `data/lessons/lessons.json`, and the course targets listed above.

Arithmetic and fractions are learner domains, so the computed concepts follow the existing solution-schema policy. The ledger will mark schema-bearing concepts before authoring; no new exemption is added.

The wave also runs the five-check overlap audit, updates the variety baseline, regenerates mastery units, graph, flow and docs, and runs `make validate`, `go test ./...` and `git diff --check`. The additive baseline remains frozen at 1,500 concepts and the mastery band width remains 3. The commit is local. Nothing is pushed or deployed.
