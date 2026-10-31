# Content expansion

This document records what the corpus is, how a new piece of content is added, and where the
corpus is thin enough that adding content would change what a learner actually meets. It is an
**analytical-tier** doc (AGENTS.md section 6): it carries measurements and the method behind
them, so `scripts/check_ste.py` does not cover it and would report violations if run against it.
That is the expected result, not a defect.

Every figure below is read from `data/`, from `internal/`, or from a script in `scripts/`. None
is hand-written. Section 7 gives the command that reproduces each one, because a number in a doc
and the artifact it describes are one fact (AGENTS.md section 7).

## 1. What the corpus contains

| Layer | Owner | Tool or origin |
|---|---|---|
| Concepts (DAG nodes) | `data/concepts/<domain>.json`, one file per domain | hand-authored |
| Enrichment overlays | `data/concepts/enrichment.json` | `encompasses`, `key_prerequisites`, `interference_group`, `variants` |
| Problems | `internal/generator/<domain>/` | generated on demand, never stored |
| Lessons | `data/lessons/{teaching,algebrica,authored}/**.md` | imports plus authored files |
| Lesson index | `data/lessons/lessons.json` | concept id to lesson source |
| KP shards | `data/lessons/kp/<concept_id>.json` | `scripts/gen_kp_shards.py` from an inline spec |
| Solution schemas | `data/lessons/solutions/<concept_id>.json` | hand-authored, learner domains only |
| Diagrams | `data/diagrams/meta.json` plus assets | Algebrica imports and generated figures |
| Courses | `data/courses.json` | named target sets |
| Destinations | `data/destinations.json` | named bundles of courses |

Two numbers are useful to hold together. There are 1250 concepts, and there are 1250 KP shard files
holding 3750 knowledge points, three per concept.

## 2. The pipeline

Each stage has one owner file, and one guard that fails when the stage is wrong.

| # | Stage | Owner | Guard |
|---|---|---|---|
| 1 | Concept node | `data/concepts/<domain>.json` | `scripts/validate_graph.go`: no cycles over `prerequisites` and `encompasses` (Kahn's), no orphans, no duplicate ids, valid grading enum, `streak` positive, `avg_time_seconds` in range, no singleton interference group |
| 2 | Enrichment | `data/concepts/enrichment.json` | the same loader: a `key_prerequisite` must already be a prerequisite, variant difficulty in range |
| 3 | Generator | `internal/generator/<domain>/`, registered in `internal/generator/all/all.go` | `Registry.Register` self-check (a graded generator must accept its own answer and reject an empty one), `TestFuzz`, and the pool ratchet in `internal/generator/pools_test.go` |
| 4 | Lesson | lesson markdown plus `data/lessons/lessons.json` | `scripts/audit_lessons.py`: every concept has a lesson, every id is live, every file is on disk, no orphan source |
| 5 | KP shard | `data/lessons/kp/`, generated from the `SPEC` in `scripts/gen_kp_shards.py` | `section` must name a real heading in the lesson; the three-KP rule is read from the shards by `scripts/check_docs.py` |
| 6 | Solution schema | `data/lessons/solutions/` | `TestLearnerDomainsHaveSchemas` over the six learner domains, plus `solutions_contract_test` |
| 7 | Diagram | `data/diagrams/meta.json` | `scripts/audit_lessons.py`: a mapping must point at an existing asset and an existing concept |
| 8 | Course | `data/courses.json` | `TestAllCoursesResolve`: every target resolves in the DAG |
| 9 | Documented counts | `README.md` | `scripts/counts.py --check` and `scripts/check_docs.py` |

The order is a dependency order. A concept with no generator cannot serve a problem. A generator
with no lesson has nothing to teach before practice. A KP shard whose `section` does not resolve
leaves the learner a placeholder instead of a worked example, and the audit refuses it.

## 3. What complete means today

All four content guards pass on the current corpus:

```
OK: 687 concepts across 17 domains                    scripts/validate_graph.go
OK: 687 concepts mapped, 687 entries, 586 sources,
    687 kp files, 2061 KPs, kp shards resolved,
    0 issues                                          scripts/audit_lessons.py
counts ok: 687 concepts, 687 KP files, 23 courses     scripts/counts.py --check
OK: docs agree with the code and the corpus           scripts/check_docs.py
```

`scripts/audit_baseline.json` is empty on every key it holds, including `stub_lessons`,
`placeholder_kps`, `thin_pools_single_answer` and `thin_pools_yesno_only`. That file is a
ratchet: known violations live there and shrink to empty as content lands, and a new violation
fails the build. Empty means the corpus is clean and the next wave starts from zero, but it also
means a new thin concept cannot be added quietly.

## 4. Where the corpus is thin

### 4.1 The course catalog reaches 1250 of the 1250 concepts

Courses carry target sets, and a target expands to its transitive prerequisite closure in
`internal/planning/planner.go:64` (`collectChain`). Across the 28 courses those closures reach
1250 of the 1250 concepts: 0 of them have a node, a generator, a lesson and three KP shards and
still sit outside every named course path. The wiring batch of 2026-10-19 moved the 149 that used
to sit there into the courses that teach their subject, and the table stays as the measurement on
the number, because a concept that no course targets is unreachable whatever the diagnostic
offers.

The count of unreachable concepts by domain:

| Domain | Concepts | Unreachable |
|---|---|---|
| calc1 | 76 | 0 |
| calc2 | 63 | 0 |
| calc3 | 61 | 0 |
| algebra | 124 | 0 |
| topology | 42 | 0 |
| statistics | 55 | 0 |
| precalculus | 50 | 0 |
| fractions | 45 | 0 |
| machine_learning | 45 | 0 |
| complex_numbers | 42 | 0 |
| number_theory | 42 | 0 |
| differential_equations | 42 | 0 |
| prealgebra | 60 | 0 |
| trigonometry | 68 | 0 |
| arithmetic | 60 | 0 |
| linear_algebra | 45 | 0 |
| abstract_algebra | 45 | 0 |
| discrete_math | 55 | 0 |
| geometry | 45 | 0 |
| graph_theory | 40 | 0 |
| probability_theory | 50 | 0 |
| numerical_analysis | 40 | 0 |
| real_analysis | 55 | 0 |

Two structural facts sit behind that table. First, a course is a *sampled* path, not a domain
cover: `calc1` targets 50 of the 76 in its domain, `calc2` targets 46 of the 63 in its domain and
`calc3` targets 54 of the 61 in its domain, so the three courses are stages rather than a partition.
Second, one whole
domain still has no course of its own: `trigonometry` is reached only by the closures of `a2` and
`pc`, which bring in 54 and 29 of its 68 concepts, so a learner meets trigonometry as algebra and
precalculus material rather than as a trigonometry syllabus.

### 4.2 Whole domains ask a learner to pick rather than to produce

The `grading_type` mix shows where the corpus is recognition-shaped. The criterion here is a
domain where `multiple_choice` is at least four concepts in ten; these ten qualify,
largest first:

| Domain | Concepts | `multiple_choice` |
|---|---|---|
| topology | 42 | 30 |
| real_analysis | 55 | 33 |
| abstract_algebra | 45 | 19 |
| trigonometry | 68 | 22 |
| probability_theory | 50 | 21 |
| discrete_math | 55 | 16 |
| machine_learning | 45 | 19 |
| complex_numbers | 42 | 19 |
| numerical_analysis | 40 | 18 |
| number_theory | 42 | 16 |

Against that, the three calculus domains carry 72 `symbolic` and 15 `expression` concepts, and `algebra` carries 29 `symbolic` and 13 `expression`. The pool
ratchet in `internal/generator/pools_test.go` exists for exactly this surface, and it runs at two
depths. An easy pool that yields a single distinct answer, or that is yes/no only, is a
violation, and that check is baseline-backed: `scripts/audit_baseline.json` is empty on both keys
because every registered generator clears that bar for its current grading type. A deeper pool is
measured too: the five calculus generators that were widened on 2026-10-19 must each produce at
least eight distinct answers over 300 samples at mid difficulty, because a feed that draws from
four or five answers repeats itself within a session. A new concept that fails either check fails
the build. Turning one of these concepts into `numeric` or `symbolic` means a real generator and a
solution schema, because the explanation stops being one fixed string.

**The first conversion landed.** `topo.dimension` was a table of yes/no rows; it now asks for the
Lebesgue covering dimension of a Euclidean space, a sphere, a torus or the Cantor set, computes
the number from the ambient dimension, and publishes `space`, `dim` and `rule` as facts for the
authored schema in `data/lessons/solutions/topo.dimension.json`. Its `grading_type` is `numeric`,
so topology had 29 `multiple_choice` concepts rather than 30 when the batch landed; the domain has
since grown to 42 concepts, 35 of them multiple choice. The worked example of the path is the four
files it touched: `data/concepts/topology.json`,
`internal/generator/topology/generators.go`, `internal/generator/topology/generators_test.go` and
the schema.

**The abstract-algebra batch converted ten.** The domain was the most recognition-shaped in the
corpus: 43 of its 45 concepts asked a learner to choose a statement. Ten of them have an answer
that is a number, so each keeps its content and loses only the option list: the order of an
element in a cyclic group (`n/gcd(a,n)`), the index of a subgroup, the number of generators of
`Z_n`, the order of `S_n` and of a dihedral group, the order of a finite field, an orbit size, a
kernel size, the number of automorphisms of `Z_n`, and a greatest common divisor. Every one
publishes the numbers its answer follows from, and
`internal/generator/abstract/conversions_test.go` recomputes the answer from those facts rather
than comparing it with itself. The domain now carries 33 `multiple_choice` concepts rather than
43, and the table above is updated to match.

**The machine-learning batch converted ten.** Its procedures are easy to describe and just as
easy to run: mean squared error over three predictions, the sign a perceptron returns, a ReLU
forward pass, one gradient-descent step, an L2 penalty, the size of a validation fold, the width
of a drop-first one-hot encoding, the parameter count of a convolution layer, expected
calibration error over two confidence bins, and one momentum update. Each now asks for the number
rather than for a description of the procedure, publishes its inputs as facts, and
`internal/generator/machinelearning/conversions_test.go` recomputes the answer from them. The
domain carries 19 `multiple_choice` concepts rather than 29.

**The trigonometry batch converted six and authored six schemas.** The six are the definitions
themselves, and a learner now applies each one rather than naming it: given that the sine is the
ratio of the opposite side to the hypotenuse, and a hypotenuse of a chosen length, how long is the
opposite side; the same question for the cosine, tangent, secant, cosecant and cotangent. The
triangles are drawn from eight Pythagorean triples scaled by a factor of three to ten, so every
answer is an exact integer rather than a rounded decimal, and each generator publishes the ratio,
the side it gives and the side it wants. Unlike the two batches above, trigonometry is one of the
six learner domains, so each of the six also needed an authored schema, 
`data/lessons/solutions/trig.func.sine.json` and its five siblings and their table-driven
exemptions are gone from `internal/generator/solutions_coverage_test.go`. Learner-domain schema
coverage moves from 349 of 402 to 355, and trigonometry's `multiple_choice` count falls from 29 to
23, which takes it below the four-in-ten criterion and off the table above.

**The algebra batch converted eight and kept ten exempt.** The eight have an answer a generator
can compute: the smallest integer in the domain of `sqrt(x - k)`, a square or cube root of a
perfect power, `a^b`, the horizontal-axis crossing of a generated line, which of two candidates
survives an extraneous-root check, the value of `x` in a two-equation system, how many of three
listed functions are injective, and the unknown coefficient in a cubic polynomial under the factor theorem. Each publishes its inputs as facts, and
`internal/generator/algebra/conversions_test.go` recomputes the answer from them. The other ten
,  the Dirichlet function's continuity, a function's parity, the sigmoid's three fixed questions,
and the seven multiple-choice concepts, have categorical or definitional answers with no
interpolated parameters, so their exemptions stay: a schema for them would be one fixed paragraph,
which is what the allowlist exists to avoid. Learner-domain schema coverage moves from 355 of 402
to 363.

### 4.3 Lesson depth is no longer thin at the bottom

Lessons come from three places. The `teaching/` tree holds 329 mappings, the Algebrica import
holds 328 and 231 are authored. The median lesson is 5,128 bytes, and no lesson is under 1 KB: the
last six were filled on 2026-10-17.

Those six were `calc.series.ratio_test` (671 bytes), `abstract.group.order` (960), `linalg.decomp.lu`
(1018), `discrete.combinatorics.surjections` (987), `precalc.func.inverse_value` (1,006) and
`calc.vec.divergence` (1,011), and they are between 3.4 and 5.5 KB now. The first four of the
`number_theory` batch, `nt.basics.crt`, `nt.adv.wilson`, `nt.basics.bezout` and
`nt.basics.prime_inf`, went from 853-920 bytes to 3.7-4.7 KB the same day, inside the three
headings their knowledge-point shards already named, so no shard needed regenerating and
`scripts/audit_lessons.py` still resolves all three per concept.

Two things are worth keeping from how the list behaved. It grew while it shrank: three authored
lessons written after the table was first measured were also under 1 KB, so the same measurement
found ten thin lessons where seven had been listed. And one of the six needed more than a longer
body: `calc.series.ratio_test` had a single `##` section with three knowledge points pointing at it,
which is the placeholder shape section 5.3 warns about, so it was restructured into three sections
and its shard rewritten with them.

`scripts/check_docs.py` reads this paragraph from the corpus rather than a copy: it counts the
lessons actually under 1 KB, compares that with the sentence here, and fails on any row that names a
lesson which has outgrown the threshold, the drift this document has now recorded three times.

### 4.4 Ingestion dropped figures and kept their captions

`docs/lesson-media.md` is the owner of this measurement, and it reports 141 paragraphs across
42 lessons that read as a caption whose image was dropped: 49 `certain`, 33 `strong` and 59 likely.
The `certain` tier, captions that open with an auto-numbered `Figure 3.2` label, which the ingestion
pipeline emits for a figure and which survives when the image does not, had been emptied by
reclaiming every case, and the Precalculus 2e trigonometry re-ingest of 2026-10-17 restored 49 of
them along with their lessons' text, as `docs/figure-reclaim.md` records.
`scripts/audit_lesson_media.py` regenerates the list, and the remaining work is per instance: each
candidate is either matched to an existing asset in `data/diagrams/meta.json` or removed.

### 4.5 Diagrams and solution schemas are complete against their own rules

`data/diagrams/meta.json` maps 183 assets, each to at least one concept and a license.
`data/lessons/solutions/` holds 395 schemas, and the six learner domains require one per
computed concept; the table-driven exceptions are an explicit allowlist with a reason per entry
in `internal/generator/solutions_coverage_test.go`. Growing either is additive work, not repair.

### 4.6 Application siblings: the same technique inside a story

Twenty-one concepts apply a technique their parent teaches rather than teach a new one:
`arith.{add,sub,mult,div}.word`, `frac.{add,sub,mult,div,compare,mixed}.word`,
`dec.ops.{add,sub,mult,div}.word`, `ratio.ops.{rate,proportion}.word`, `prealg.eq.word`,
`alg.eq.vars_both_sides.word`, `alg.ineq.multi_step.word`, `alg.systems.word` and
`geo.triangle.pythagorean.word`. They cover 17 subdomains.

Each one is a *dependent*: a parent technique is one of its prerequisites. Twenty of the
twenty-one are also leaves, so nothing else depends on them and no course closure needed them
until the wiring batch of 2026-10-19 named them as targets, which is the only way a leaf enters a
closure. `prealg.eq.word` is the exception on both counts: something depends on it, and that
something is what pulls it into a course closure.

The pattern is cheap because it writes no lesson prose. A sibling shares its parent's lesson file
,  `data/lessons/lessons.json` points both ids at one source, so its three KP sections come from
headings that already exist, and the new files are the shard and the authored solution schema. The
generator and its test extend package files that are already there, and the concept node joins an
existing domain file.

Where it stops, measured against the same two conditions: the parent lesson must carry an
application-shaped heading to draw three KPs from. `trig.basics.sin_cos_def` has one and it is
"Applications in integration", which is not a word problem; `geo.quad.perimeter` and
`geo.solid.volume_rect` have none; `prealg.ineq.two_step` and `alg.quad.discriminant` have one each
and carry two empty headings, so a shard there would need written prose first. Past this point the
work is lesson authoring, not generator writing.

## 5. Options, and what each one costs

### 5.1 Course wiring: sample more, or cover the domain

The gap in section 4.1 has three plausible policies.

- **Add a domain's sinks as targets.** A sink is a concept with no dependents in its domain, and
  every concept reaches some sink, so the sink set covers the domain by construction. This is
  cheap and mechanical, and it is wrong on its own: the calculus sinks are 77 concepts, so
  `calc1` and `calc2` would both become the whole domain and stop being stages.
- **Add the sinks of the depth band a course occupies.** This preserves staging and is the shape
  `cmd/masteryunits` already uses for its draft partition. It needs the band boundary to be a
  decision first.
- **Add one course per uncovered domain.** The only domain with no course at all is
  `complex_numbers`. This is a small, additive change, but it changes the course count, and the
  course count is read from `data/courses.json` by `scripts/counts.py` and quoted in the README.

None of the three is safe to apply silently. All three change the size of a course, and course
size drives the mastered/total figure and the days-remaining estimate a learner sees.

**What the two cover-the-domain policies do, measured.** Both were applied to the real DAG and the
real targets. A *sink* here is a concept that no other concept **in its own domain** depends on (296
of 666), and "the band a course occupies" is `floor(conceptDepth / 3)`, the shape `cmd/masteryunits`
uses:

| Course | Now | A: domain sinks | B: band sinks, same domains |
|---|---:|---:|---:|
| 4th Grade Math | 17 | 107 | 40 |
| 5th Grade Math | 31 | 143 | 98 |
| 6th Grade Math | 33 | 118 | 107 |
| 7th Grade Math | 40 | 121 | 114 |
| 8th Grade Math | 41 | 121 | 115 |
| Prealgebra | 35 | 209 | 147 |
| Algebra I | 49 | 209 | 181 |
| Geometry | 34 | 121 | 109 |
| Algebra II | 66 | 263 | 263 |
| Precalculus | 90 | 481 | 311 |
| Calculus I | 67 | 342 | 305 |
| Calculus II | 70 | 342 | 342 |
| Linear Algebra | 58 | 235 | 223 |
| Discrete Math | 56 | 293 | 293 |
| Discrete Math II | 62 | 293 | 293 |
| Probability & Statistics | 47 | 268 | 262 |
| Number Theory | 28 | 83 | 51 |
| Differential Equations | 77 | 376 | 356 |
| Abstract Algebra | 58 | 296 | 237 |
| Abstract Algebra II | 73 | 296 | 292 |
| Topology | 48 | 405 | 298 |
| Machine Learning | 62 | 452 | 407 |

Three things follow, and together they decide the first two policies.

- **A makes the catalog exhaustive.** The union of the 28 closures goes from 247 to **666**: every
  concept in the corpus becomes reachable. That is the end of a sampled path, the catalog and the
  corpus become the same set and it leaves the diagnostic and a free-choice goal nothing to do.
- **Both policies collapse the staged pairs.** `calc1` and `calc2` are 67 and 70 concepts now and 342
  and 342 under A; `discrete`/`discrete2` become 293/293, and `abstract`/`abstract2` 296/296. Two
  stages that reach the same set are one stage.
- **B needs the domain filter, and even with it inflates.** Unfiltered, `floor(depth / 3)` alone
  hands 4th Grade Math 36 statistics concepts and 22 geometry concepts, because a sink at the same
  depth in an unrelated domain passes the test. Filtered to the domains a course already reaches,
  every course is still 2.5x to 6x its current size.

**So the answer to the first two is "sample more".** The third is separate and additive: a course for
`complex_numbers`, whose 12 sinks close over 113 nodes and reach all 26 of the domain's, of which
20 are reached by no course path today. The union reach goes from 247 to 273 and the course count
from 22 to 23, which the README quotes.

**The narrowest closure rule does not help either.** Narrowing the band to the subdomains a course
already touches, sinks that share a subdomain with something the course reaches, the union still
goes from 247 to 543, and the courses still inflate 1.9x to 3.4x: 4th Grade Math 17 to 32, `calc1`
67 to 230, `calc2` 70 to 265, `diffeq` 77 to 263. `discrete` and `discrete2` still land on 145 each.
No rule that closes over a set of sinks is compatible with staging at any radius, so the targets have
to be *chosen*, and that is curriculum work rather than a policy switch.

**What "sample more" still needs is a number.** Widening the sample has no new invariant to satisfy,
so the decision left is how far: it is a target size for the catalog, and the measurements above are
the boundary that number has to be drawn against. Section 6's first step is that decision, and it is
also the gate on a new domain (section 5.4), because a domain with no course target is unreachable.

**The wiring batch of 2026-10-19 chose the targets.** The 149 concepts that section 4.1 measured as
unreachable were added, subject by subject, to the courses that teach them, so the catalog now
reaches all 1250 and the table in section 4.1 reads zero on every row. That is the curriculum work
the paragraph above says the fix has to be, done deliberately rather than by a closure rule, and it
does not make the courses domain covers: `calc1` reaches all 76 concepts in its own domain and
shares the rest of the calculus corpus with `calc2` and `calc3`, and trigonometry still has no
course of its own.

### 5.2 Producing an answer in a multiple-choice domain

This is the highest-value per-concept change and it needs no catalog decision. It is a generator
plus a solution schema for one concept: change `grading_type`, write a `GradedGenerator` whose
easy pool discriminates, add the schema so a missed answer carries an authored explanation, and
let the pool ratchet and the schema contract test say whether it holds.

### 5.3 Filling a thin lesson

A lesson under 1 KB is a lesson with one heading and little else, and its three KP shards must
name three different headings in it, which is where the placeholder risk lives. Filling one means
writing the body first and regenerating the shards second, because the shard sections resolve
against the body.

### 5.4 What adding a domain costs, itemized from the one we have

The plan's last step adds a domain, with topology as the copy target. Measured off topology, which is
the only domain added since the corpus was built:

| What | Topology | Where |
|---|---:|---|
| Concepts | 30, in 11 subdomains | `data/concepts/topology.json`, 12 KB |
| Generator | 2 files, 52 KB | `internal/generator/topology/`, the package and its test, registered in `internal/generator/all/all.go` and exercised by the pool ratchet in `internal/generator/pools_test.go` |
| KP shards | 30 | `data/lessons/kp/topo.*.json`, three per concept |
| Solution schemas | 1 | `data/lessons/solutions/`, one, because 29 of the 30 are `multiple_choice` and their explanation is a fixed string |
| Course | 1, three targets | `data/courses.json` |
| Domain label | 1 | `DOMAIN_LABELS` in `scripts/counts.py` |
| Regenerated | | `data/mastery_units.json`, `data/lessons/lessons.json`, `data/lessons/mapping.json`, `data/lessons/figure-audit.json`, `scripts/verify_report.json`, the concept graph |

That is six hand-edited things and six regenerated ones per domain, and the count of concepts it adds
is the size of its concept file: the 30 above, or whatever the new domain's is.

**Which of the three counts a domain moves is not the same question as which concepts move.**
`domains` and `courses` are read from the data, `len(by_domain)` over the concepts' own `domain`
field and `len(data/courses.json)`, so a domain added *inside* an existing domain's namespace leaves
the domain count at 17 and only the concept count moves. A genuinely new domain takes it to 18, and a
course for `complex_numbers` (section 5.1) takes the course count to 23 without touching the domain
count, because that domain exists already.

**Nothing checked that a domain was labeled until now.** `DOMAIN_LABELS` holds 17 entries and the
data holds 17 domains, but the two were kept in step by hand, and the failure was silent: the
README's domain table iterates the label list, so an unlabeled domain is absent from the table while
`counts.py` still counts it. `counts.py --check` fails on that mismatch now, in both directions,
because the change this section prepares for is the one that would have tripped it.

### 5.5 The concept target, and where it goes

The catalog holds 1250 concepts. The remaining per-domain targets add up to 0, and this is the table
that says where: the thousand-concept milestone is complete, and graph theory has opened the 1,250
tranche. `Now` is read from the corpus; `Target` is the decision. The rows keep the order step 6 of
section 6 used, largest deficit first.

| Domain | Tier | Now | Target | Add |
|---|---|---:|---:|---:|
| numerical_analysis | university | 40 | 40 | +0 |
| probability_theory | university | 50 | 50 | +0 |
| graph_theory | university | 40 | 40 | +0 |
| arithmetic | school | 60 | 60 | +0 |
| geometry | school | 45 | 45 | +0 |
| fractions | school | 45 | 45 | +0 |
| discrete_math | university | 55 | 55 | +0 |
| complex_numbers | university | 42 | 42 | +0 |
| prealgebra | school | 60 | 60 | +0 |
| abstract_algebra | university | 45 | 45 | +0 |
| machine_learning | university | 45 | 45 | +0 |
| topology | university | 42 | 42 | +0 |
| linear_algebra | university | 45 | 45 | +0 |
| statistics | bridge | 55 | 55 | +0 |
| number_theory | university | 42 | 42 | +0 |
| differential_equations | university | 42 | 42 | +0 |
| precalculus | bridge | 50 | 50 | +0 |
| real_analysis | university | 55 | 55 | +0 |
| calc1 | bridge | 76 | 76 | +0 |
| calc2 | bridge | 63 | 63 | +0 |
| calc3 | bridge | 61 | 61 | +0 |
| algebra | school | 124 | 124 | +0 |
| trigonometry | school | 68 | 68 | +0 |

**The `Now` column had drifted, and this is the correction.** The table was drawn when the corpus held
703 concepts, and the batches since have grown it to 819. `scripts/sweep_counts.py` keeps the sentence
above the table current, because its size is a pattern rule, but it has no rule for the rows, so ten
of the seventeen stayed at their 703-era values while the corpus moved, and the roadmap they publish
named deficits that were larger than the work left: trigonometry read as 16 concepts behind where it
was 8, and calculus as 24 where it is 10. Two domains, geometry and fractions, had in fact reached
their targets and read as though they had not. The remainder is 181 rather than 297, which is the same
seventeen targets measured from a larger starting corpus, and the rows are re-sorted so that "largest
deficit first" still names the largest deficit. `check_docs.py` now reads each `Now` from
`data/concepts/`, each `Add` from `Target - Now`, and the sentence above the table from the two, so the
next drift fails instead of being published.

**Trigonometry reached its target on 2026-10-17.** Eight computed concepts, the sum and difference
formulas, the half-angle formulas, a quadratic in a ratio, the ambiguous case of the law of sines, an
area from two sides and an included angle, a vector's length and its dot product, and a sinusoid's
phase shift, took the domain from 40 to 48, so its row reads +0 and now sits at the foot of the table
with geometry and fractions. The remainder is 173, and the largest deficit left is discrete_math's
sixteen.

**Discrete mathematics reached its target on 2026-10-17.** Sixteen concepts, argument validity,
nested quantifiers, tautology and contradiction, the contrapositive, disproof by counterexample,
structural induction, power sets, partial orders, countable sets, counting by complement, bipartite
matchings, Euler's formula for planar graphs, Cayley's formula, degree sequences, linear recurrences
and growth rates, took the domain from 39 to 55, so its row reads +0 and sits at the foot of the
table with geometry, fractions and trigonometry. Every one is a computed question: the
counterexample concept scans for the first composite value rather than quoting one, the matching
concept searches for a maximum matching, and the recurrence concept derives its closed form from the
roots it chose. The remainder was 157 at that point.

**Complex numbers reached their target on 2026-10-17.** Sixteen concepts, parts and equality,
powers, reciprocals, plane distance, loci, argument arithmetic, complex linear equations, conjugate
root pairs, trigonometric and hyperbolic identities from exponentials, complex geometric series,
Möbius fixed points, residues at simple poles, the cross-ratio, open and closed sets, and domains, 
took the domain from 26 to 42, so its row reads +0. A new college course, `cx` "Complex Numbers and
Analysis", targets all 42: the domain had been reachable only through `pc`'s single polar-form
target, which left twenty of the twenty-six behind it outside every course. The batch also brought
the `complex` grading path its first tests, which is how the polar rounding recorded in
`docs/code-defects.md` came to be documented rather than guessed at. The remainder was 141 at that
point.

**Algebra reached its target on 2026-10-18.** Fifteen concepts, a line through two points,
intercepts, factoring by grouping, sum and difference of cubes, vertex form, a quadratic area word
problem, end behavior, simplifying, multiplying and dividing, adding and subtracting rational
expressions, rational exponents, rationalizing a denominator, direct and inverse variation, change
of base, and classifying a conic, took the domain from 89 to 104, so its row reads +0. Two
subdomains are new, `algebra.rational` and `algebra.radicals`, because the three rational-expression
concepts and the two radical concepts had no existing home. The conic-recognition concept is also
what pulls the four conic prerequisites into a course: targeting it brings `circle`, `ellipse`,
`parabola` and `hyperbola` inside the closure for the first time, and it reaches `func.quad` as
well. The remainder was 126, and the largest deficit left was arithmetic's fourteen.

**Arithmetic reached its target on 2026-10-18.** Fourteen concepts ordering whole numbers,
expanded form, multiplying and dividing by powers of ten, estimating sums and products, a multi-step
word problem, interpreting a remainder, division by zero, integer word problems, powers with
negative bases, the power of a product and of a quotient, cube roots, and estimating square roots, 
took the domain from 46 to 60, so its row reads +0. Two subdomains are new,
`arithmetic.estimation` and `arithmetic.word_problems`, because the estimation and multi-step
concepts had no existing home. The five grade courses 4 through 8 take the batch: the fourth-grade
set gains ordering, expanded form, powers of ten and estimation of sums; the fifth-grade set gains
estimation of products, interpreting remainders, multi-step problems and division by zero; the sixth
grade gains integer word problems, the seventh negative bases, and the eighth the product and
quotient rules, cube roots and square-root bounds. All fourteen new concepts are reachable, and
their prerequisites pull ten older arithmetic concepts into a course path for the first time. Three
older concepts still sit outside every course: `arith.round.thousands`, `arith.factor.composite`
and `arith.order_ops.nested`. The remainder was 112, and the largest deficit left was prealgebra's
fourteen.

**Prealgebra reached its target on 2026-10-18.** Fourteen concepts, percent change, percent error,
commission, markup, simple interest, the percent proportion, unit-price comparison, translating words
into expressions, evaluating formulas, inequality word problems, number-line inequality graphs,
customary units, metric units, and Celsius-Fahrenheit conversion, took the domain from 46 to 60, so
its row reads +0. One subdomain is new, `prealgebra.measurement`, because the three measurement
concepts had no existing home. Each new concept has an authored lesson, three knowledge-point
sections, and a solution schema. The grade courses 5 through 8 take the batch: the fifth-grade set
gains decimal word problems for all four operations and rounding; the sixth-grade set gains decimal
conversion, percent conversion, percent change, commission, markup, the percent proportion,
expression translation and metric conversion; the seventh-grade set gains tax and tip, ratio word
problems, scale, absolute value, signed order of operations, negative and zero exponents, unit-price
comparison, inequality graphs and customary conversion; the eighth-grade set gains percent error,
simple interest, formula evaluation, inequality word problems and temperature conversion. The
prealgebra course adds the capstone applications. All 60 prealgebra concepts are reachable, and the
new targets pull every remaining decimal, percent, ratio, expression and inequality prerequisite into
a course path. The remainder was 98, and the largest deficit left was abstract_algebra's thirteen.

**Abstract algebra reached its target on 2026-10-18.** Thirteen concepts, the center of a group,
automorphisms, simple groups, p-groups, composition series, zero divisors, integral domains, prime
and maximal ideals, Euclidean domains, the characteristic of a field, algebraic elements, submodules,
and module homomorphisms, took the domain from 32 to 45, so its row reads +0. Three candidate
concepts were dropped because an existing lesson already teaches them: conjugacy is covered by
`abstract.group.normal`, so the batch teaches composition series; splitting fields are covered by
`abstract.field.extension`, so it teaches algebraic and transcendental elements; and vector spaces
are covered by `abstract.structures.module`, so it teaches submodules. The prerequisite direction is
preserved: `rings.euclidean` depends on `rings.def`, `rings.ideal` and `rings.polynomial`, not on
`rings.ufd`, because Euclidean implies PID implies UFD and not the reverse. Prime and maximal ideals
are stated for a commutative ring with identity, the setting in which the quotient tests hold, and
the characteristic is kept distinct from the field order: a question on `F_9` asks for its
characteristic of 3 rather than its order of 9. The two computed concepts, p-groups and the
characteristic, are numeric and publish the facts their answers follow from; the other eleven are
multiple choice. All 45 abstract-algebra concepts are reachable, including the eleven that no course
reached before, because `abstract2` now targets the sinks that close over them: the ten targeted
gaps plus the four older sinks `commutator`, `galois_group`, `noetherian` and `exact_sequence`. The
batch adds no authored solution schema, consistent with the university-domain policy. The remainder
was 85, and the largest deficit left was machine_learning's twelve, tied with topology's twelve.

**Machine learning reached its target on 2026-10-18.** Twelve concepts, sigmoid activation,
logistic regression, vanishing and exploding gradients, mean absolute error, the confusion matrix,
F1, k-means clustering, feature scaling, one-hot encoding, class imbalance, Gini impurity, and
bagging and random forests, took the domain from 33 to 45, so its row reads +0. Two subdomains are
new, `machine_learning.data` for the three preprocessing concepts and `machine_learning.trees` for
the two tree concepts. The survey rejected the obvious additions before writing them: the CNN lesson
already teaches output-size arithmetic, the validation lesson already teaches leakage, the
gradient-descent lesson already teaches learning-rate tuning, the overfitting lesson already teaches
capacity, and the transformer lesson already teaches softmax attention, so none of those became a
duplicate concept. Seven of the twelve are numeric, which raises the domain from 24 of 33 multiple
choice to 29 of 45; the numeric generators publish the facts their answers follow from (residuals,
confusion cells, precision and recall, coordinates, mean and standard deviation, class counts) and
the tests recompute each answer from them. All 45 machine-learning concepts are reachable, including
the sixteen that no course reached before, because the `ml` course now targets the advanced sinks
(`bert`, `gan`, `actor_critic`, `adam`, `batchnorm`, `calibration`, `augmentation`) as well as the
new ones. The remainder is 73, and the largest deficit left is topology's twelve.

**Topology reached its target on 2026-10-18.** Twelve concepts, open balls, interior, closure and
boundary, limit points, convergence, completeness, bases and subbases, the subspace topology, dense
sets and separability, homeomorphisms, the Euler characteristic, genus, and Betti numbers, took the
domain from 30 to 42, so its row reads +0. One subdomain is new, `topology.invariants`, for the four
invariant concepts. The survey found that the domain jumped from three basics straight to named
theorems (Urysohn, Tietze, Stone-Čech, Van Kampen, de Rham) with almost no computational middle, so
the batch builds the metric-space toolkit the later lessons assume and adds the computable
invariants the algebraic-topology lessons describe but never name. Six of the twelve are numeric,
which raises the domain from 1 of 30 produced answers to 7 of 42; each numeric generator publishes
the facts its answers follow from (coordinates, boundary points, leading coefficients, cells, and
the Euler characteristic or genus) and the tests recompute the answer from them. All 42 topology
concepts are reachable, including the 34 that no course reached before, because the `topo` course
now targets the ten existing sinks as well as the nine new ones, and its description now covers the
full sequence. The remainder was 61, and the largest deficit left was linear_algebra's eleven, tied
with statistics' eleven.

**Linear algebra reached its target on 2026-10-18.** Eleven concepts, vector-space axioms, kernel
and image, injective, surjective and invertible maps, coordinates and change of basis, orthogonal
complements, projection onto a subspace, the spectral theorem, quadratic forms and definiteness, the
cross product, the condition number, and Cayley-Hamilton, took the domain from 34 to 45, so its row
reads +0. The survey found the domain well supplied with computations but thin on the structural
ideas that connect them, and it rejected the obvious additions that existing lessons already teach:
determinant properties, transpose, linear combinations and coordinate representation, the
characteristic polynomial, algebraic and geometric multiplicity, and Gram-Schmidt behind the QR
lesson. Five of the eleven are numeric and one is a tuple (the cross product), so the batch keeps the
domain's produced-answer emphasis; each numeric generator publishes the facts its answers follow
from (matrix rows, basis vectors, dimensions, singular values) and the tests recompute the rank,
coordinates, dimension, projection, cross product and condition number from them. All 45
linear-algebra concepts are reachable, including the 22 that no course reached before, because the
`linalg` course now targets the five unreachable sinks and the eight new ones, and its description
now covers inner products and matrix decompositions. The remainder was 50, and the largest deficit
left was statistics' eleven.

**Statistics reached its target on 2026-10-18.** Eleven concepts, null and alternative hypotheses,
the one-sample z statistic, the one-sample t statistic, p-values, Type I and II errors and power, the
standard error of a proportion, the chi-square goodness-of-fit test, one-way ANOVA, the coefficient
of determination, histograms, and nonparametric tests, took the domain from 44 to 55, so its row
reads +0. The survey found the domain already strong on computation but with no testing framework:
sampling, confidence intervals, standard error and margin of error existed, but nothing connected
them to a hypothesis test. It rejected the obvious additions that existing lessons already teach:
standard deviation and correlation in the variance/covariance lesson, sample size in confidence
intervals, five-number summaries in the quartiles lesson, the regression line in machine learning,
and the central limit theorem in sampling distributions. Six of the eleven are numeric, and each
numeric generator publishes the facts its answers follow from (the mean, null value, spread and
sample size; the observed and expected counts; the groups; the two sums of squares; the values and
the bin) and the tests recompute the statistic from them. The domain's multiple-choice count rises
from 1 of 44 to 5 of 55, still the corpus's most computational domain. All 55 statistics concepts are
reachable, including the 32 that no course reached before, because the `probstat` course now targets
the sixteen existing sinks and the eight new ones, and its description now covers hypothesis testing
and the tests built on it. The remainder was 39, and the largest deficit left was calculus' ten.

**Calculus reached its target on 2026-10-18.** Ten concepts, derivatives of hyperbolic functions, the
second derivative of a parametric curve, finite differences, arc length in polar coordinates, torsion,
the Frenet frame, the second-order Taylor approximation, hydrostatic force, Pappus's centroid theorem,
and moment of inertia, took the domain from 165 to 175, so its row reads +0. The survey found the
domain already dense and rejected the additions existing lessons already cover: Taylor's theorem with
remainder and the Weierstrass M-test in the series lessons, numerical-integration error bounds, the
Hessian in the partial-derivatives lesson, hyperbolic integrals in the integration tables, and surface
area of graphs and revolutions. Seven of the ten are numeric, one is an expression, one is symbolic
and one is a tuple. The wiring is deliberately scoped: section 5.1 measured sink targeting and
rejected it, because the calculus sinks are most of the domain and would make `calc1` and `calc2`
exhaustive, so the three courses gain only the three new concepts that belong to their stage.
Calculus therefore reads 86 of 175 reachable, with the staged courses preserved; the corpus reaches
731 of 971. The remainder was 29, and the largest deficit left was number_theory's ten, tied with
differential_equations' ten.

**Number theory reached its target on 2026-10-18.** Ten concepts, the Jacobi symbol, discrete
logarithms, square roots modulo a prime, sums of two squares, Pell's equation, Pisano periods, the
Carmichael function, the Miller-Rabin test, prime gaps, and Bertrand's postulate, took the domain
from 32 to 42, so its row reads +0. The survey rejected the additions existing lessons already teach:
modular inverses in the Bézout lesson, Euler's criterion in the Legendre lesson, perfect numbers and
Mersenne primes in the divisor-sum lesson, twin primes and Brun's sieve in the sieve lesson, and
Gaussian integers in the class-number lesson. Eight of the ten are numeric, and each numeric
generator publishes the facts its answers follow from (the Jacobi arguments, the prime and primitive
root, the residue and modulus, n, d, the Fibonacci modulus, the Carmichael input, or the prime before
the gap) and the tests recompute the answer from them. The wiring is scoped plus applied: the `nt`
course gains the ten new concepts together with `crypto` and `diophantine`, the two topics its
description already advertises. Number theory reads 33 of 42 reachable, with the nine analytic and
graduate concepts (Möbius, PNT, zeta, Chebyshev, Dirichlet, elliptic curves, class numbers,
quadratic reciprocity and the sieve) left as enrichment; the corpus reaches 750 of 981. The remainder
is 19, and the largest deficit left was differential_equations' ten.

**Differential equations reached its target on 2026-10-18.** Ten concepts, Bernoulli equations,
Cauchy-Euler equations, reduction of order, power-series solutions, the matrix exponential, Heaviside
steps and impulses, the spring-mass system, mixing problems, predator-prey systems, and Heun's method
,  took the domain from 32 to 42, so its row reads +0. The survey rejected the additions existing
lessons already teach: integrating factors in the linear and exact lessons, the y = vx substitution in
the homogeneous lesson, undetermined coefficients and the modification rule in the nonhomogeneous
lesson, separation of variables and d'Alembert in the heat and wave lessons, linearization in the
stability lesson, and Picard and Fredholm in the boundary and Volterra lessons. Five of the ten are
numeric, and each numeric generator publishes the facts its answers follow from (the Cauchy-Euler
coefficients, the series kind and rate, the spring constants and initial data, the tank volume and
concentration, the Lotka-Volterra coefficients, and the Heun step) and the tests recompute the answer
from them. The wiring is core: the `diffeq` course gains the ten new concepts together with
`first_order.exact`, `second_order.variation_params` and `numerical.rk4`, so exact equations,
variation of parameters and RK4 finally have a course path. Differential equations reads 33 of 42
reachable, with the transforms, PDE, Sturm-Liouville and nonlinear chain left as enrichment; the
corpus reaches 764 of 991. The remainder was 9, and the largest deficit left was precalculus' nine.

**Precalculus reached its target on 2026-10-18, and the corpus reached a thousand.** Nine concepts, 
continuous compounding, Newton's law of cooling, logarithmic scales, polynomial and rational
inequalities, holes in rational graphs, the discriminant of a general conic, projectile motion,
resultant vectors, and decomposition of a function, took the domain from 41 to 50, so its row reads
+0. The survey found the domain's twenty unreachable concepts were all standard precalculus (conics,
sequences and series, vectors, the binomial theorem, induction, parametric and polar graphs), so the
course was wired to cover them rather than leave the subject half-taught. It rejected the additions
existing lessons already teach: half-life in the growth lesson, logarithm properties in the logarithm
lessons, eccentricity in the conic lessons, unit vectors in the vectors lesson, the midpoint in the
distance lesson, even and odd functions in the functions lessons, and one-to-one functions in the
inverse lesson. Six of the nine are numeric, one is an expression, and two are multiple choice; each
numeric generator publishes the facts its answers follow from, and the tests recompute the answer from
them. All 50 precalculus concepts are reachable, so the corpus reaches 794 of 1000, with 206 left
outside the named course paths, the enrichment section 4.1 measures, not a gap in a target. Every
domain on the table was at its target at that milestone.

**The thousand-concept milestone is complete, and the next tranche is 1,250.** The target table is
no longer the whole story: the corpus now exceeds a thousand, and the next phase adds four domains, 
real analysis, graph theory, numerical analysis and probability theory and raises the targets of
calculus, trigonometry and algebra by 70. Graph theory went first and is complete: its 40 concepts
arrived in two batches. The first covered isomorphism, complements and subgraphs, homomorphisms and
products, distance and cuts, Menger's theorem, matchings and covers, Tutte's theorem, and chromatic,
edge and list coloring. The second covered planarity and minors, spectral graph theory, extremal and
random graphs, and flows and digraphs. The `graph` course targets all 40, so the domain reads 40/40
reachable and the corpus reaches 834 of 1,040 with 206 enrichment concepts outside the named paths.
Real analysis, numerical analysis and probability theory follow, with the three target raises. The
raises are recorded in the table above, and each new domain takes its row as it opens.

**Calculus reached its target on 2026-10-19.** Thirty concepts took the domain from 175 to 205, so
its row reads +0. The survey found the domain already dense in one-variable technique, so the batch
deepens the edges of each subdomain instead of restating what the lessons already teach: limits gain
the intermediate value theorem as a computation, limit superior and inferior, and oscillation;
derivatives gain the inverse hyperbolic rules, the mean value theorem as a physical statement,
Newton's method as a convergence question, convexity and Jensen's inequality, marginal analysis, and
the epsilon-delta definition of the derivative; series gain operations, the binomial series, the
Cauchy product, the p-series threshold, Dirichlet's test and the gamma function; integrals gain the
trapezoid error bound, double integrals in polar coordinates, scalar surface integrals, the beta
function, consumer surplus and Fourier coefficients; multivariable calculus gains Lagrange
multipliers with two constraints, the implicit and inverse function theorems, iterated limits and
Euler's homogeneous function theorem; and vector calculus gains the scalar line integral, parametric
surface area, general flux integrals and the vector identities. Nineteen concepts are numeric and
eleven are multiple choice; each numeric generator publishes the facts its answer follows from, and
the tests recompute the answer from those facts. The three staged courses each take the new concepts
that belong to their stage, nine to `calc1`, ten to `calc2` and eleven to `calc3`, so all 205
concepts are reachable and the remaining per-domain total falls to 40.

**Algebra reached its target on 2026-10-19.** Twenty concepts took the domain from 104 to 124, so
its row reads +0. The survey found the groups uneven, functions, exponentials and logarithms,
radicals, rational expressions, systems, sequences and polynomials were each thinner than the
linear and quadratic material, so the batch fills those edges rather than deepening the core
again. Functions gain piecewise evaluation, graph transformations and the domain and range of an
inverse; exponentials and logarithms gain their graphs, compound interest and the solution of a
logarithmic model; radicals gain simplifying, combining like radicals and rationalizing with
conjugates; rational expressions gain complex fractions, proportions and vertical asymptotes;
systems gain a three-variable elimination, the matrix method and a nonlinear pair; sequences gain
sigma notation and induction; and polynomials gain the factor theorem and the rational root
theorem. Twelve concepts are numeric and eight are conceptual multiple choice. Every numeric
generator publishes the facts its answer follows from, and the tests recompute the answer from
those facts. The three school courses each take the concepts of their stage, nine to `a1`, ten to
`a2` and one to `pa` and the corpus reaches 1064 of 1230, so the remaining per-domain total falls
to 20. Algebra is a school domain, so the twelve computed concepts carry authored solution schemas
and the eight fixed-question concepts are named in the exemption list, which the schema gate checks
in both directions.

**Trigonometry reached its target on 2026-10-19.** Twenty concepts took the domain from 48 to 68,
so its row reads +0 and the remaining per-domain total falls to zero. The survey found the domain
already solid on the unit circle, the right triangle and the laws of sines and cosines, but thin at
the edges: nothing taught periodicity as a property, angular velocity, bearings, the product-to-sum
and sum-to-product identities, the triple-angle and power-reduction formulas, the variants of the
Pythagorean identity, general solutions and linear equations, inequalities on an interval, the
tangent and reciprocal graphs, transformations and inverse graphs, the hyperbolic identities and
their inverses and derivatives, polar coordinates or the scalar projection. The batch fills exactly
those gaps. Nine concepts are numeric and eleven are multiple choice; each numeric generator
publishes the facts its answer follows from, and the tests recompute the answer from those facts:
the period from the coefficient, the angular velocity, the bearing from the quadrant, the
triple-angle sine, the power-reduction value, the linear-equation root, the inverse-hyperbolic
cosine, the polar abscissa and the scalar projection. Trigonometry is a school domain, so the nine
computed concepts carry authored solution schemas and the eleven fixed-question concepts are named
in the exemption list, which the gate checks in both directions. There is no trigonometry course;
the domain is reached through the algebra and precalculus courses, so fifteen concepts join `a2`
and five join `pc`, and the corpus reaches 1101 of 1250, leaving 149 outside the named course
paths. Every concept has an authored lesson, three knowledge-point sections and a lesson
registration. Counts, the mastery partition and the web graph are regenerated; the target table is
re-sorted and this paragraph records the change.

**A target per domain, not a rule.** Rule-based widening was measured in section 5.1 and does not
work: closing over a domain's sinks makes the catalog exhaustive, and every smaller radius still
collapses the staged pairs. So the number is a curriculum decision and is written down here rather
than derived, which is why the deltas above are computable and the targets are not.

**What the tiers mean.** *school* is the arithmetic-to-trigonometry sequence a learner meets first,
and its targets are set so each domain carries a year of practice rather than a sampling of one.
*bridge* is precalculus, statistics and calculus, which carry the heaviest course loads and the
longest dependency chains. *university* is the rest, where a target is a floor and not a ceiling:
enough that the domain reads as a subject rather than as a module attached to one course.

**Why the floor matters more than the depth.** Eleven of the seventeen domains sit under 40 concepts,
and five sit at 25 or fewer with most of their content `multiple_choice`. Section 4.2 measured that
as the corpus's weakest surface, so the floor is not only a count: a domain grown from 25 to 42 is
also where the produced-answer conversions of section 5.2 get their room.

## 6. A sequenced plan

1. Decide the course-wiring policy in section 5.1, then apply it with `TestAllCoursesResolve` as the
   guard and `scripts/counts.py --write` for the README. The measurements are in 5.1; what is missing
   is the target size, and the decision gates step 5.
2. Convert one `multiple_choice` concept in an upper domain to a produced answer, as the worked
   example that proves the generator and schema path there. Done: `topo.dimension`, section 4.2.
3. Fill the sub-1 KB lessons, regenerating each shard after its body grows. Done: the last ten were
   filled on 2026-10-17, four of them `number_theory` lessons already in section 4.3's table and six
   that the same measurement found, and no lesson is under 1 KB now. `calc.series.ratio_test` also
   needed three sections instead of one, so its shard was rewritten with it.
4. Reclaim the dropped figures. Done: all 102 landed on 2026-10-15 and 2026-10-16 as lesson figures
   under `web/next-app/public/diagrams/`, which need no `meta.json` entry, so the `certain` tier is
   empty and the remainder is measured by section 4.4.
5. Add a domain only after 1 to 4, using the topology addition as the copy target, because a domain
   costs six hand-edited files, a course target and six regenerated ones (section 5.4).
6. Reach the thousand, in the order of section 5.5's targets, largest deficit first. Each batch
   regenerates everything through `scripts/sweep_counts.py`, which is the difference between a batch
   costing a turn and costing an afternoon.

## 7. Method

Reproduce each measurement. All commands run from the repository root.

```bash
go run scripts/validate_graph.go
python3 scripts/audit_lessons.py
python3 scripts/counts.py --check
python3 scripts/check_docs.py
```

Course reach and the grading-type mix need the DAG and the course targets read together; the
reach number is the union of each course's prerequisite closure, computed the way
`internal/planning/planner.go:64` computes it, and the checker in `scripts/check_docs.py`
recomputes it rather than trusting this file.

```bash
python3 - <<'PY'
import collections, glob, json
concepts, prereq = {}, {}
for path in glob.glob("data/concepts/*.json"):
    if path.endswith("enrichment.json"):
        continue
    for c in json.load(open(path)):
        concepts[c["id"]] = c
        prereq[c["id"]] = set(c.get("prerequisites") or [])

def closure(ids):
    seen, stack = set(), list(ids)
    while stack:
        cid = stack.pop()
        if cid in seen:
            continue
        seen.add(cid)
        stack += list(prereq.get(cid, ()))
    return seen

reach = set()
for course in json.load(open("data/courses.json")):
    reach |= closure(course["targets"])
print("reach", len(reach), "of", len(concepts))
by_domain = collections.Counter(concepts[c]["domain"] for c in set(concepts) - reach)
print("unreachable", dict(by_domain.most_common()))
mix = collections.Counter(concepts[c].get("grading_type") for c in concepts)
print("grading types", dict(mix))
PY
```

Lesson sizes and sources come from `data/lessons/lessons.json`; figure candidates come from
`scripts/audit_lesson_media.py`. The application siblings of section 4.6 are the concept ids ending
in `.word` in `data/concepts/*.json`; whether anything depends on one, and whether a course closure
reaches it, are read from those files and from `data/courses.json`.
