# Reclaiming the dropped figures

This is the work plan for the figures the lesson ingestion dropped while keeping their captions.
It is an **analytical-tier** doc (AGENTS.md section 6), so `scripts/check_ste.py` does not cover
it. The measurement belongs to `scripts/audit_lesson_media.py`, which writes
`data/lessons/figure-audit.json` and the report in `docs/lesson-media.md`; this file turns the
`certain` tier of that report into a plan.

## 1. What is missing

141 paragraphs across 42 lessons read as a caption whose image was dropped: 49 `certain`, 33 `strong`
and 59 `likely`. Every one sits
in `data/lessons/teaching/`. The `certain` tier — the one that opened with an
auto-numbered `Figure 3.2` label, which the PreTeXt conversion emits for a figure and keeps when it
drops the image, so that it needed no judgement — holds 49 rows again, every one of them restored by
the Precalculus 2e re-ingest of 2026-10-17, and section 3 records that against the batches that had
emptied it. The rest is counted by weaker evidence: `strong` describes a picture or
is duplicated in its lesson, and `likely` is the teaching corpus, where a standalone italic paragraph
is the ingestion pipeline's own signature.

## 2. The upstream is OpenStax, and the palette is still the dialect

Every teaching lesson opens with a credit line naming its source, for example:

```
> Content sourced from [OpenStax Prealgebra](https://openstax.org/books/prealgebra/pages/1-introduction) by Marecek & Anthony-Smith - CC BY-NC-SA 4.0
```

Three books account for every `certain` caption: **OpenStax Prealgebra** (36), **OpenStax Calculus
Volume 1** (36) and, since 2026-10-17, **OpenStax Precalculus 2e** (49). The first seventy-two have
landed, as section 7 records. All three are CC BY-NC-SA 4.0, which permits reuse and adaptation with
attribution but carries the terms on to the adapted work, and the attribution is already the lesson's
own header line. The first two headers said CC BY 4.0 until 2026-10-17, which dropped the
NonCommercial and ShareAlike terms from the record; `THIRD_PARTY_NOTICES.md` now quotes what each
book's own panel says.

That changes what "a direct look-alike" means for these figures. Their *drawing* is upstream's, not
Mathua's: they are fetched rather than redrawn, so the line work, the labels and the layout are
OpenStax's and none of it is a redraw.

**Their color, though, goes, and this reverses an earlier decision.** The paragraph that stood here
said a re-fetched figure keeps the OpenStax look and does not need to match Algebrica. That was
wrong for the reason the style doc already gives about redrawn figures: a lesson page shows several
figures from the same source, and one figure that does not match its neighbours is worse than no
figure at all. A full-color raster beside greyscale line art is exactly that outlier, so every
fetched asset is converted at fetch time — grey, its tone stretched to the full ink-to-paper range,
then mapped onto the greys that section 4 of `docs/diagram-style.md` fixes. The drawing survives; the
saturated color was the only part of it that never belonged on the page.

## 3. The certain tier was emptied, and the 2026-10-17 re-ingest put 49 rows back

Every caption that opened with an auto-numbered label had landed by 2026-10-16. The tier stood at
119 rows over 49 distinct figures when the fetch batches began, and it went 119 to 72 to 41 to 8 to
5 to 0, one batch at a time, which section 7 records.

It holds 49 rows again, and they returned in one movement rather than one batch. Re-ingesting
`trig.period` and `trig.radians` from OpenStax Precalculus 2e restored the two sections' figure
captions along with their text: 25 in §6.1 Graphs of the Sine and Cosine Functions and 24 in §5.1
Angles. They are all 49 of section 1's `certain` rows; the 33 `strong` and 71 `likely` rows beside
them are the ones no batch has touched. The retired hand-rolled converter had dropped these lessons'
captions, so the committed pair never carried them and the tier could be — and was — driven to zero
without the loss showing.

These 49 are work rather than a judgement, because a caption the lesson's prose cites by number is
one that section 4 assigns a disposition to and section 6 lists, and every such caption before them
was reclaimed rather than left standing. Both lessons cite figures by number throughout, and
`scripts/prune_orphan_captions.py` keeps a caption in a lesson that cites a figure, so pruning them
is closed off as well. What they needed first was the disposition step section 4 describes: open the
upstream section and read the figure's `alt` text, because the caption is the wrong text to classify
from — that is the correction the first fetch batch forced. That step is done. Section 6 lists all
forty-nine, and forty-seven of them are redraws.

## 4. Dispositions

| Disposition | Figures | Action |
| --- | --- | --- |
| schematic, redraw | 47 | Draw in the Algebrica cool dialect from section 4 of the style doc. Graphs, number lines, regions and solids are fully determined by their captions, so no upstream fetch is needed and the result is Mathua's own work. |
| illustration, fetch | 1 | Fetch the OpenStax original into `web/next-app/public/diagrams/<book>/` and put it on the dialect's greys. The license is CC BY-NC-SA 4.0 and the credit line is already in the lesson header. |
| unidentified | 0 | Empty since 2026-10-16: every bare caption was classified by opening its upstream section and looking at the image, and each of the seventeen is a `schematic, redraw` or an `illustration, fetch` now. |
| photo, third-party credit | 1 | Fetch under the credit line the caption already carries, which names a photographer or agency (`credit: modification of work by Steve Jurvetson, Wikimedia Commons` and similar); a photograph cannot be redrawn as a line figure, so a replacement was never the alternative. The five rows it held were fetched by 2026-10-16, and section 6 records the open question that keeps the new row unlanded. |

**One correction the first batch forced.** Figures 1.7, 1.8 and 1.9 were assigned
`illustration, fetch` when this table was drawn, because the assignment read the caption and
nothing else: *76 is closer to 80 than to 70* describes a picture of a number line, and a picture
of a number line is an illustration. Opening Prealgebra 1.1 shows all three are the same 70..80
number line as Figure 1.2 under a different mark, which makes them schematics. They were redrawn
(section 7), so the redraw row gained three figures and then lost the same three, and the fetch
row lost three: when that batch landed the table above read 43 redraw and 37 fetch, where reading
the caption alone had suggested 35 and 38, and the classification below moved those counts on
again to 53 and 44, of which section 7 has since landed one hundred and two.

**The column is a guess until someone opens upstream, and the caption is the wrong thing to
read.** The reliable text is the upstream `alt` attribute, which describes the picture rather than
the argument the picture supports: Figure 1.7's caption is a sentence about rounding, while its
alt is *An image of a number line from 70 to 80 with increments of one*. That is how Figure 1.6 was
placed here as a `schematic, redraw` even though its caption is a bare `Figure 1.6` — Prealgebra
1.1 gives it an alt of *A chart titled 'Place value' with fifteen columns and 4 rows* — and it is
how the seventeen `unidentified` rows were classified on 2026-10-16.

**The seventeen were read twice, and the first reading was not enough.** The disposition is a
reading of the description, so the description has to be the whole one: the scrape behind this plan
truncated every alt at 110 characters, and the 23 upstream sections were re-fetched to read
`data-alt` in full (Figure 6.34's is 1,244 characters). The text alone still does not settle it,
because an alt describes a picture in words and the words can hide the shape: Figure 1.4 opens *An
image with three items* and the three items are a cube, a rod and a flat, which is what the lesson
calls base-10 blocks, and nothing in *An image of 28 cookies placed at random* says whether that is
a drawing or a photograph. So each of the seventeen images was downloaded and looked at, and the
disposition below is the reading of the picture with the alt as its check.

**Where the line falls: ten redraw, seven fetch.** A figure is a schematic when its content is a
number line, a graph, a geometric block, tile, circle or region, or a table of words and numbers —
a redraw of that is faithful and is Mathua's own work. It is an illustration when the subject is a
photograph (three stacks of dollar bills, 24 pennies, 28 cookies in three of the figures, two
pizzas) or when it is a table whose cells are formulas and small plots (6.34), which an SVG cannot
typeset. Two calls are worth recording because a later reader will meet them again: Figure 5.7 is a
calendar page with seven distances written on it, which is a table of words and numbers and is
therefore a redraw even though the upstream image is a picture of a calendar; and 1.3 is a
photograph of currency with no credit line, so it is `illustration, fetch` and not the
`photo, third-party credit` row, which exists for a third party whose name must travel with the
image rather than for every photograph.

**What lands with them.** All seventeen keep `(none)` in the Caption column, because upstream's
`figcaption` holds only the label for these figures — the alt is the only description there is. It
is also what a landed figure will carry: the rule that a lesson figure's alt is its caption's
words, never a paraphrase, has nothing to copy here, and the same document's alt is still the
author's words.

**Ten of the figures below were invisible until 2026-10-15**, because of a defect in this plan's own
measurement rather than in the corpus. `section_a_candidates` skipped any standalone-italic
paragraph whose neighbour started with `-`, `*`, `|` or `>`, which is right for a line inside a
list and wrong for a caption that simply sits above one: Figure 1.5 is above the place-value table,
Figure 1.6 above the bulleted list of digit values, Figures 4.50 and 4.51 are two captions in a
row. Exempting auto-numbered labels from that test added 24 rows and these ten figures — 1.5, 1.6,
3.20, 4.40, 4.50, 4.51, 4.52, 4.58, 4.59 and 5.3 — every one of which rates `certain` and none of
which is a list item. The tier is 220 rows now rather than 196, which is a fix and not a
regression: the same corpus, counted without the exclusion.

## 5. The captions were damaged too, and are repaired

The caption text was ingested with `figcaption.get_text(' ', strip=True)`. OpenStax writes a
mathematical expression as `<math><semantics><mrow>…</mrow><annotation-xml>…</annotation-xml>
</semantics></math>`: the expression is present **twice**, and the flattening joins every token
with a space. So a caption arrived as `… the number 138 . 138 .` — the expression, then the
expression again. The fault is in `scripts/ingest/ingest_openstax.py` and not in the document.
A figure restored under a caption like that still reads as broken, so repairing the caption was
part of the work rather than a separate cleanup.

**Scale.** The corpus has 244 lines that carry a figure caption; 160 of them come from a book the
modern ingester fetches. Of those 160, 119 were repaired, 40 already matched the pipeline's output,
and 1 was held back for a hand. The repair landed on 2026-10-15, and the re-run is idempotent:
159 already correct, 0 to repair, 1 held back.

**The repair took the text from upstream**, which is the only reliable source for it: each caption
was regenerated through the fixed ingester, so it arrives in the same form as the lesson body
around it, and then passed through the corpus `dollar-number` rule
(`scripts/lint_math_layout.py`) that bare numerals are upright text rather than inline math. That
rule is why the repaired captions read `… worth 35 cents, or $\frac{35}{100}$ of a dollar.` and not
`… worth $35$ cents`. The test that the repair is complete is a comparison against upstream, not a
pattern: the numbers alone were never the whole damage, and letters, expressions, percentages and
dollar amounts were doubled in the same way.

Two things the repair had to know:

- **A bare caption is not damage.** `Figure 1.3`, `1.4` and `1.6` have no caption upstream either —
  their `figcaption` holds only the label — so they are reclaim targets rather than repair targets,
  and section 4 records what each of them turned out to be.
- **One caption needed a hand.** `Figure 2.18` carries a piecewise function as an inline MathML
  table. Pandoc renders it with literal newlines, and the line-scoped dollar balancing then escapes
  the span's opener, so the ingester cannot emit it correctly today; the text in the lesson is the
  pipeline's intended output written on one line.

**The guard** is rule 16 in `scripts/audit_lessons.py`: a caption may not contain an immediately
repeated token run (`doubled_caption_run`). Measured against the damaged corpus it catches 63 of the
67 distinct damaged captions and fires on none of the repaired ones, and it has no baseline, so
re-ingesting with the old figure handler fails the audit. Rule 16 covers captions only, because a
corpus-wide doubled-token test is not usable: it reports 1009 legitimately repeated pairs, mostly
`& &` inside math and paired table cells.

**One measurement it moved.** The repaired captions that gained inner `*x*` emphasis — OpenStax
marks a variable with `<em>` — stopped matching `audit_lesson_media.py`'s standalone-emphasis test,
which removed eight imageless captions from this plan. The test now allows a balanced inner span
and still rejects a thematic break (`* * *`); the report stood at 324 rows and 220 in the
certain tier on 2026-10-15, after two reclaims in section 7 had removed eighteen rows and the label
exemption in section 4 had restored twenty-four. Section 7 records what each batch since has
removed, and section 1 states what is left.

**One glyph normalisation, in the same text.** OpenStax marks the parts of a two-part example with
circled letters (`ⓐ` … `ⓔ`), and the ingester copied the codepoint through. At lesson size `ⓐ` is an
`a` inside a circle — which is what `@` is — so a caption reads as if an email sign preceded it, and
no theme font carries the codepoint, so the browser substitutes a face at a different size and
baseline. The corpus already spelled the parts out in prose (`part (a)`, 22 times, next to a marker
that said `ⓐ`), so all 2340 markers across 72 lessons, several of them inside captions, were written
as `(a)` … `(e)` on 2026-10-15. AGENTS.md section 1 records the reasoning and rule 17 of
`scripts/audit_lessons.py` enforces it; the captions in section 6 are the repaired text with that
normalisation applied.

**Repaired on 2026-10-17, and a different fault.** `trig.period` and `trig.radians` were the only
two lessons that came from the retired hand-rolled converter (`scripts/extract_openstax.py`), and
their text was mangled rather than merely doubled: `_underline Send underline_ine` stood for `Sine`,
`sin(t)= opposite hypotenuse` stood for $\sin t = \frac{\text{opposite}}{\text{hypotenuse}}$, and
every expression appeared twice. The modern pipeline renders both sections correctly, but it produced
844 and 1073 lines raw against the committed 279 and 557: the committed files were curated, with back
matter stripped and the wiki-lesson fixes applied, so re-ingesting them would replace the lessons
rather than repair them. That was a content decision, and it was taken rather than left open: the
mangling is not mechanical, so re-authoring every expression by hand is the same work with more ways
to be wrong, and both lessons were re-ingested instead. `trig.radians` comes from §5.1 Angles,
unchanged. `trig.period` comes from §6.1 Graphs of the Sine and Cosine Functions, because that is the
section carrying amplitude, period and phase shift; it had been mapped to §5.4 Right Triangle
Trigonometry, whose lesson has none of the three, and `docs/content-expansion.md` is where the
mapping target is recorded. Back matter was stripped at the first apparatus heading, as rule 12
requires, and the circled part markers became `(a)` … `(e)` through the same normalisation the
paragraph above records — this time written into the ingester, so every future run gets it.

**The re-ingest found two faults in the ingester, and both are fixed.** A MathML table cell that
pandoc renders as display math ends its row early: the cell continues on later lines as a fenced
`math` block, whose closing line also carries the rest of the row, so the table loses its columns.
The worse damage is a step later: the fence pass turns that block into a `$$` pair that never closes
on a line of its own, so the display-block pass stays open past it and strips every inline `$` in the
rest of the page. §6.1's prose after its two tables read `the range is \left\lbrack {- 1,1}
\right\rbrack.` as literal text on that account. `fold_table_math` in
`scripts/ingest/ingest_openstax.py` folds such a cell back into its row as inline math before either
pass runs, and §6.1's two tables are the corpus's first exercise of it. Restoring the sections
restored their figure captions as well, which is the 49 `certain` rows section 3 accounts for.

## 6. The figures

**The table holds forty-nine rows, and they are the ones the 2026-10-17 re-ingest restored.** It has
one row per figure with a disposition — the `schematic, redraw`, `illustration, fetch`,
`unidentified` and `photo, third-party credit` captions section 4 defines — and every row that stood
in it before had landed, so each of those now lives in section 7's record of its batch instead.
`Lessons` is how many concept lessons carry the caption, because one figure is restored once and
every lesson that cites it gains the image; the caption is the text as repaired in section 5,
truncated, so it is what the figure carries now.

The forty-nine are the `certain` rows section 3 records: the captions the Precalculus 2e re-ingest
restored, twenty-five in §6.1 Graphs of the Sine and Cosine Functions and twenty-four in §5.1 Angles.
Each disposition was assigned from upstream rather than from the caption, which is the correction
section 4 records: both sections were opened and every figure's `alt` attribute read in full, and the
five whose alt named a photograph or an illustration rather than a graph — the rainbow, the lifted
circle, the spring, the sprinkler and the water wheel — were downloaded and looked at. Section 1's
other rows are not listed here because no disposition was ever assigned to them either: a `strong` or
`likely` candidate is judged one instance at a time — matched to an asset that already exists, or
removed — and that measurement belongs to `scripts/audit_lesson_media.py`.

**The Figure column repeats itself because the book does.** Precalculus 2e's web pages number figures
per section, so both §6.1 and §5.1 open at `Figure 1`, and the caption text the corpus carries uses
the same page-local label. In print the book numbers them per chapter (`Figure 6.1`, `Figure 5.1`),
which is the form to give a reader who wants the upstream image.

| Figure | Caption (repaired) | Book | Section | Lessons | Disposition |
| --- | --- | --- | --- | --- | --- |
| Figure 1 | Light can be separated into colors because of its wavelike properties. (credit: "wonderferret"/ Flickr) | Precalculus 2e | §6.1 | 1 | photo, third-party credit |
| Figure 2 | The sine function | Precalculus 2e | §6.1 | 1 | schematic, redraw |
| Figure 3 | Plotting values of the sine function | Precalculus 2e | §6.1 | 1 | schematic, redraw |
| Figure 4 | The cosine function | Precalculus 2e | §6.1 | 1 | schematic, redraw |
| Figure 5 | (none) | Precalculus 2e | §6.1 | 1 | schematic, redraw |
| Figure 6 | Odd symmetry of the sine function | Precalculus 2e | §6.1 | 1 | schematic, redraw |
| Figure 7 | Even symmetry of the cosine function | Precalculus 2e | §6.1 | 1 | schematic, redraw |
| Figure 8 | (none) | Precalculus 2e | §6.1 | 1 | schematic, redraw |
| Figure 9 | (none) | Precalculus 2e | §6.1 | 1 | schematic, redraw |
| Figure 10 | (none) | Precalculus 2e | §6.1 | 1 | schematic, redraw |
| Figure 11 | (none) | Precalculus 2e | §6.1 | 1 | schematic, redraw |
| Figure 12 | (none) | Precalculus 2e | §6.1 | 1 | schematic, redraw |
| Figure 13 | (none) | Precalculus 2e | §6.1 | 1 | schematic, redraw |
| Figure 14 | (none) | Precalculus 2e | §6.1 | 1 | schematic, redraw |
| Figure 15 | (none) | Precalculus 2e | §6.1 | 1 | schematic, redraw |
| Figure 16 | (none) | Precalculus 2e | §6.1 | 1 | schematic, redraw |
| Figure 17 | (none) | Precalculus 2e | §6.1 | 1 | schematic, redraw |
| Figure 18 | (none) | Precalculus 2e | §6.1 | 1 | schematic, redraw |
| Figure 19 | (none) | Precalculus 2e | §6.1 | 1 | schematic, redraw |
| Figure 20 | A horizontally compressed, vertically stretched, and horizontally… | Precalculus 2e | §6.1 | 1 | schematic, redraw |
| Figure 21 | (none) | Precalculus 2e | §6.1 | 1 | schematic, redraw |
| Figure 22 | (none) | Precalculus 2e | §6.1 | 1 | schematic, redraw |
| Figure 23 | (none) | Precalculus 2e | §6.1 | 1 | schematic, redraw |
| Figure 24 | (none) | Precalculus 2e | §6.1 | 1 | schematic, redraw |
| Figure 25 | (none) | Precalculus 2e | §6.1 | 1 | schematic, redraw |
| Figure 1 | (none) | Precalculus 2e | §5.1 | 1 | schematic, redraw |
| Figure 2 | (none) | Precalculus 2e | §5.1 | 1 | schematic, redraw |
| Figure 3 | Angle theta, shown as ∠θ | Precalculus 2e | §5.1 | 1 | schematic, redraw |
| Figure 4 | (none) | Precalculus 2e | §5.1 | 1 | schematic, redraw |
| Figure 5 | (none) | Precalculus 2e | §5.1 | 1 | schematic, redraw |
| Figure 6 | (none) | Precalculus 2e | §5.1 | 1 | schematic, redraw |
| Figure 7 | Quadrantal angles are angles in standard position whose terminal… | Precalculus 2e | §5.1 | 1 | schematic, redraw |
| Figure 8 | (none) | Precalculus 2e | §5.1 | 1 | schematic, redraw |
| Figure 9 | (none) | Precalculus 2e | §5.1 | 1 | schematic, redraw |
| Figure 10 | (none) | Precalculus 2e | §5.1 | 1 | schematic, redraw |
| Figure 11 | The angle t sweeps out a measure of one radian. Note that the… | Precalculus 2e | §5.1 | 1 | schematic, redraw |
| Figure 12 | (a) In an angle of 1 radian, the arc length s equals the radius… | Precalculus 2e | §5.1 | 1 | schematic, redraw |
| Figure 13 | A 45° angle contains one-eighth of the circumference of a circle… | Precalculus 2e | §5.1 | 1 | schematic, redraw |
| Figure 14 | Commonly encountered angles measured in degrees | Precalculus 2e | §5.1 | 1 | schematic, redraw |
| Figure 15 | Commonly encountered angles measured in radians | Precalculus 2e | §5.1 | 1 | schematic, redraw |
| Figure 16 | An angle of 140° and an angle of –220° are coterminal angles. | Precalculus 2e | §5.1 | 1 | schematic, redraw |
| Figure 17 | (none) | Precalculus 2e | §5.1 | 1 | schematic, redraw |
| Figure 18 | (none) | Precalculus 2e | §5.1 | 1 | schematic, redraw |
| Figure 19 | (none) | Precalculus 2e | §5.1 | 1 | schematic, redraw |
| Figure 20 | (none) | Precalculus 2e | §5.1 | 1 | schematic, redraw |
| Figure 21 | (none) | Precalculus 2e | §5.1 | 1 | schematic, redraw |
| Figure 22 | The area of the sector equals half the square of the radius times… | Precalculus 2e | §5.1 | 1 | schematic, redraw |
| Figure 23 | The sprinkler sprays 20 ft within an arc of 30°. | Precalculus 2e | §5.1 | 1 | schematic, redraw |
| Figure 24 | (none) | Precalculus 2e | §5.1 | 1 | illustration, fetch |

**The tally is forty-seven redraws, one fetch and one third-party photograph, and no `unidentified`
row.** Three of the calls are worth recording, because a later reader will meet them again:

* **§5.1 Figure 24, the water wheel, is `illustration, fetch`** even though it is a drawing and not a
  photograph, because its content is a specific object rather than a relation: the prose calls it *the
  design of the 3rd century Roman water wheel in Hierapolis*, which the caption does not describe and
  a redraw would have to invent. Section 4's line is that a redraw needs content fully determined by
  its caption, and this one is not determined by anything but the picture.
* **§6.1 Figure 25, the spring, is `schematic, redraw`** for the reason the water wheel is not: the
  render is shaded and pictorial, but what it carries is one measured length, `y`, and the exercise
  states the setup in prose — a weight on a spring hung from a board. A coil between two boards with a
  dimension arrow is the entire content.
* **§6.1 Figure 1 is the only `photo, third-party credit`**, and its caption already carries
  `(credit: "wonderferret"/ Flickr)`. Where that credit renders is not settled: the decision of
  2026-10-17 is that attribution for a third-party photograph does not sit on the lesson page until
  the placement question is answered. The row is fetch-ready and placement-blocked, and it blocks
  nothing else.

Section 8's order applies to what is left as it did to the last four batches: the forty-seven
schematics are one batch and come first, the water wheel is a fetch, and the photograph waits on the
open question.

## 7. Progress

**Figure 1.2** ("The numbers on the number line increase from left to right, and decrease from
right to left") landed on 2026-10-15 as the first reclaim and the proof of the pipeline. It is a
`schematic, redraw`: drawn as `web/next-app/public/diagrams/arithmetic/number-line-counting.svg`,
embedded in all six lessons that cite it (`arith.place.hundreds`, `arith.place.tens`,
`arith.place.thousands`, `arith.round.hundreds`, `arith.round.tens`, `arith.round.thousands`), and
the regenerated audit dropped those six rows from the `certain` tier. It is a lesson figure, so it
needs no `data/diagrams/meta.json` entry (section 5 of the style doc explains why).

**Figures 1.7, 1.8 and 1.9** (the same 70..80 number line under three marks: 76 nearer 80, 72
nearer 70, 75 exactly midway) landed on 2026-10-15 as one batch, in those six lessons. They are
`number-line-round-76.svg`, `-72.svg` and `-75.svg`, beside Figure 1.2's
`number-line-counting.svg` and on its canvas to the pixel — 460x130, the same ink, the same
tick style — because the four sit within a few paragraphs of each other and a figure that does
not match its neighbours is worse than no figure. Each draws the distance to both neighbouring
tens to scale, so the shorter span is the argument for which way the number rounds; the ink weight
on that span's arrow agrees with it, and at 75 both spans are the same length and both are inked.
The alt text is the caption minus its `Figure 1.7 ` label, which is the rule
`scripts/place_lesson_figures.py` follows and the exact text Figure 1.2 carries.

Those three captions held eighteen rows, and the regenerated audit dropped all of them: 318
candidates to 300, and 214 in the certain tier to 196. They were the first batch drawn from a
disposition that had been assigned by reading a caption, and the reading was wrong — section 4
records the correction.

**Figures 1.4, 1.5 and 1.6** (the base-10 vocabulary, then 138 modeled with it, then 5,278,194 in a
place value chart) landed on 2026-10-16 as one batch, in those six lessons, as
`base-ten-blocks.svg`, `base-ten-138.svg` and `place-value-chart.svg`. 1.5 and 1.6 had been
`schematic, redraw` since the plan was written, read out of their upstream `alt`; 1.4 was one of the
seventeen that section 4 classified, and the three are one subject, so they land together. 1.4 and
1.5 share a cell size of 9px, which is what makes a rod ten cubes and a flat ten rods by
construction rather than by three separate drawings that happen to agree.

The upstream 1.4 stacks its three items vertically with each label beside its item, and the lesson
figure is bounded to 180px tall (section 3 of the style doc), so the batch is landscape and each
label sits under its item in two lines instead. That is a transposition of the layout and not of the
content: the same three blocks, the same three labels.

The alt text differs by figure, and the difference is the point. 1.5 has a caption, so its alt is
the caption minus its `Figure 1.5 ` label, the rule Figure 1.2 follows. 1.4 and 1.6 have only the
bare label upstream, so there is nothing of the author's to reuse except the `alt` itself: both
carry upstream's `data-alt` verbatim, which is the same document's accessibility text and not a
paraphrase of it.

Those three captions held eighteen rows, and the regenerated audit dropped all of them: 324
candidates to 306, and 220 in the certain tier to 202. The ceiling for orphan captions in
`internal/lessons/assets_test.go` ratchets down with it, to 306.

**Figure 1.12** (a 2-by-3 rectangle of unit squares, its area 6 square feet) landed on 2026-10-16 in
the six `arith.mult.*` lessons, as `rectangle-area.svg`. It is upstream's layout kept whole: `2 ft`
against the left edge, `3 ft` under the bottom edge, and `2 · 3 = 6 ft²` to the right, which is where
upstream puts it. Upstream's words are set at the family's 13px label size and the blocks are sized
to fill the canvas, so the blocks are larger than upstream's against the same words — a re-drawn
figure cannot carry upstream's text-to-shape ratio without breaking the size the other arithmetic
figures use, and the dial is turned towards the drawing. The cell size is one number, so the
rectangle is six unit squares and the two labeled edges are the drawing's own edges rather than a
coincidence of it. The alt is the caption minus its `Figure 1.12 ` label, the rule Figure 1.2
follows.

These six lessons are the first the plan has put a figure into that had none: the `arith.mult.*`
bodies carried no image at all, where the earlier batches went into lessons that already carried
Figure 1.2 or 1.7-1.9. So this is also the first reclaim whose lesson `lessonFigures.spec.ts` could
not measure before, and the case for it is added with the figure rather than after it.

Those six captions held six rows, and the regenerated audit dropped all of them: 306 rows to 300,
and 202 in the certain tier to 196. The ceiling for orphan captions in
`internal/lessons/assets_test.go` follows, to 300.

**Figures 4.2, 4.3 and 4.5** (a circle cut into thirds; five tiles of one length showing a whole,
halves, thirds, fourths and sixths; and a half above the same rectangle in sixths) landed on
2026-10-16 in the eleven `frac.*` lessons that cite them, as `fraction-circle.svg`,
`fraction-tiles.svg` and `fraction-sixths.svg`. 4.2 and 4.3 are the same seven lessons, where the
two figures sit a few paragraphs apart, so they are drawn to the same rules and the same label
size; 4.5 is in the four `frac.ops.*` and `frac.*.whole` lessons.

These are the first figures in the plan whose labels *are* mathematics. Upstream writes every label
as a stacked fraction — numerator, rule, denominator — and a figure that teaches fractions while
writing `1/3` on one line is showing the answer instead of the notation the lesson is teaching, so
each label here is three primitives, two `<text>` and a `<line>`, on a shared vertical rhythm. There
is no KaTeX inside an SVG file and no `foreignObject` survives the sanitizer, so the fraction is
drawn rather than typeset.

One deliberate departure, and it is worth naming because it is visible. Upstream colors the five
tiles differently — blue, grey, yellow, grey, peach — and the house palette has one fill,
`#e6e6e6`, with ink at `#333333`. A tile is distinguished by how many pieces it is cut into, which
is what the figure is about; color would carry no information the division count does not, so all
five tiles are the same fill. The division count is also what makes the tiles one figure rather than
five: they share a length, so the pieces of the sixths row are visibly half the width of the thirds.

The canvases are sized to their content, which for these three does not mean one width. The circle
is `240x180` — a lone circle in a `460`-wide canvas would be mostly empty panel — and the two bar
figures are `420x180`, the tiles and the rectangles being `360` wide with the same margin the number
lines leave. All three carry the opaque white background rect, so they read on the dark panel.

Their alt text is upstream's `data-alt` verbatim, all three being bare-label captions, and it was
read out of the scrape rather than retyped. That scrape is double-encoded — decoded as latin-1 while
holding UTF-8, so a curly quote arrives as three characters — and the wiring inverts exactly that
error rather than transcribing around it, which is also why Figure 4.3's alt is one 530-character
sentence that no hand copy would have survived intact.

Those eighteen captions held eighteen rows, and the regenerated audit dropped all of them: 300 to
282, and the tier that opens on an auto-numbered label from 196 to 178. The ceiling for orphan
captions in `internal/lessons/assets_test.go` follows, to 282. This batch leaves no lesson in the
`frac.*` set without a figure, since all eleven had none.

**Figures 6.3 and 6.5** (one mechanism shown in both directions: the percents of Table 6.1 as
decimals, and the decimals of Table 6.2 as percents) landed on 2026-10-16 in the three `pct.*`
lessons that cite them, as `percent-to-decimal.svg` and `decimal-to-percent.svg`.

This is the first place in the plan where the upstream is inconsistent with itself rather than
merely silent. The section page references **one image hash for both figures** — the same table — and
that table is the `Percent | Decimal` one with `006.%` through `012.5%`. Figure 6.5's alt names the
columns in the other order and lists five decimals where that image's decimal column has four
values, so the alt and the image cannot both be describing Figure 6.5. The lesson prose decides it:
"Figure 6.5 uses the decimal numbers in Table 6.2 and shows visually to convert them to percents by
moving the decimal point two places to the right". Two tables are drawn, each following its own
figure's alt, because for 6.5 the alt is the only description of what the figure is for. The four
values of Table 6.2 are in it and so is the fifth the alt names, and the prose stays true either way
since it claims the figure uses those numbers rather than only those.

The notation `006.%` is upstream's and not a slip: it writes the percent as a number with its
decimal point visible at the end, so the two hops under it have something to move. The hops are two
small arcs under the right end of the value, where the point arrives; they are not themselves
directional, and the prose says which way the point travels.

Each table is `250` wide, which is the width of its content: two columns of `115` holding labels of
up to six characters at the family's 13px. The header row is filled `#e6e6e6` and the body cells are
white, ruled with guides and framed in ink, which is the same division of structure the place value
chart uses. The white background rect is there as in every local figure, so the tables read on the
dark panel.

Their alt text is upstream's `data-alt` verbatim, read out of the scrape with its encoding inverted.
Those six captions held six rows, and the regenerated audit dropped all of them: 282 to 276, and 172
in the certain tier. The ceiling for orphan captions in `internal/lessons/assets_test.go` follows, to
276.

**Figure 5.5** (a number line from 0 to 1 marked every 0.2, with a measurement bracket over each
step) landed on 2026-10-16 in the four `dec.ops.*` lessons, as `decimal-number-line.svg`. It is the
first reclaim in the plan that reuses a dialect already in the tree rather than starting one: the
`number-line-*.svg` files have drawn the same axis, ticks, arrowheads and 13px labels since the
arithmetic set was placed, so this adds one primitive to that vocabulary, the bracket, and inherits
the rest.

The bracket is where the drawing and upstream disagree about the figure's own text. Upstream's image
braces four of the five steps -- the one from 0.8 to 1 has none -- while its alt says the braces show
"a distance of 0.2 between each adjacent set of 2 numbers", which is five. Five are drawn. That is
the rule section 7 records for Figure 6.5, and the reason is the same: the alt is the description
that ships with the figure, and a reader told about a bracket that is not there has been told
something false. It also happens to be the better figure, since the claim being made is that every
step is worth the same amount.

The bracket itself is one unfilled path: a bar whose ends turn down toward the axis, with a small
peak at the center, drawn in `#636363` because it annotates the axis rather than being part of it.
The five `0.2` labels over it are ink, not annotation, because they are the measurement the figure
exists to make; the six axis numbers stay in the annotation color, as they do in every other number
line here. The canvas is 460 by 130, the size the other number lines use, so the figures that share
lessons with this one line up.

Its alt text is upstream's `data-alt` verbatim, read out of the scrape with its encoding inverted.
Those four captions held four rows, and the regenerated audit dropped all of them: 276 to 272, and
168 in the certain tier. The ceiling for orphan captions in `internal/lessons/assets_test.go`
follows, to 272.

**Figure 5.6** (a circle with the three parts named: a dashed diameter across the widest portion, a
dashed radius from the center to the rim, and the circumference along the edge) landed on 2026-10-16
in the three lessons that cite it, as `circle-parts.svg`. Two decisions are worth naming because
both are places the dialect did not already have an answer.

The dashed lines are ink at `stroke-dasharray="4"`. The style doc ties dashes to guides, and a guide
is `#BABABA`; the Algebrica set also uses an ink dashed line for an annotation that is not
scaffolding, which is what these two are -- they are what the words point at, so drawing them in the
guide color would make the two labeled objects the faintest thing in the figure.

`Circumference` is drawn one glyph at a time: thirteen `<text>` elements, each positioned on a
radius-76 arc under the circle and rotated to that arc's tangent. The word is placed *along* the
edge, and that is what along means. `<textPath>` would do it in one element, but it needs an `href`
to a path `id`, and an `href` on a text element is exactly the kind of reference a sanitizer drops --
the failure mode would be a figure with no such word in it at all, not a figure with a straight one.
The host font is monospace and every advance is the same 7.8px, so thirteen positions computed on an
arc are the layout rather than an approximation of the layout. The canvas is `220x180`: a lone
circle, sized to itself, with the height the renderer's bound allows.

Its alt text is upstream's `data-alt` verbatim, read out of the scrape with its encoding inverted.
Those three captions held three rows, and the regenerated audit dropped all of them: 272 to 269, and
165 in the certain tier. Three lessons leave the orphan list, so the report's lesson count goes 109
to 106, and the ceiling for orphan captions in `internal/lessons/assets_test.go` follows, to 269.

**Figure 5.7** (a page of a January calendar with the week's distances written into the days: 2, 15,
8, 3, 8, 5 and 8 miles) landed on 2026-10-16 in the five `stat.prob.*` lessons, as
`statistics/calendar-distances.svg`. It is the first reclaim outside the arithmetic directory, and
that is the only part of the wiring that is new: a lesson figure is a markdown image on an absolute
path, so a subject directory costs nothing and the figure lands where a reader looking for it would
look.

Upstream's calendar is blue-ruled, its header row is the darker of two shades and its distances are
grey italic. The house has one fill, so the header takes `#e6e6e6` and the body stays white, which is
the division the percent tables landed with -- the relationship survives (the header is the darker
strip) without a second fill. The distances are ink italic rather than grey, because they are what
the lesson is about: the mean of those seven numbers is the exercise, and the dates are scaffolding
the page around them. `New Year's day` is the one `#636363` label, because it names a day rather than
contributing a value.

The canvas is `460x170`, and the width decides the layout. Seven columns of 64 hold `Wednesday` at
the 11px the stylesheet forces (59.4px of a monospace advance), and the tallest row is the first,
which carries a date, a distance and the event name. That name is set on two lines -- `New Year's`
and `day` -- because at 8.5px it needs 76 of the 64 columns available, and a label that overflows its
cell in a calendar is a label a reader has to guess at.

Its alt text is upstream's `data-alt` verbatim, read out of the scrape with its encoding inverted.
Those five captions held five rows, and the regenerated audit dropped all of them: 269 to 264, and
160 in the certain tier. Five lessons leave the orphan list, so the report's lesson count goes 106 to
101, and the ceiling for orphan captions in `internal/lessons/assets_test.go` follows, to 264.

**Figure 2.12** (three graphs of three functions at `x = 2`: a line with a hole where its denominator
vanishes, the half-lines `y = 1` and `y = -1` whose ends are both holes, and `1/(x - 2)^2` running off
the top either side of the asymptote) landed on 2026-10-16 in `calc.limit.numeric`, as
`calculus/limit-behavior.svg`. It is the first graph and the first y axis in the local set, and the
first figure under the calculus directory.

The dialect had what a graph needs. Axes and their arrowheads are `#BABABA` at 1.5, their unit ticks
belong to the axis so they are the same color, the curve is ink at 2, and the hollow marker is a
`fill="white"` circle with an ink stroke -- the shape the vendored set already uses for a point that
is not on the graph, and the one place a figure carries white inside itself.

Three things upstream does not settle, and the choices are visible. The three panels share one y
range, so their x axes land on one row and they read as small multiples of a single question. Each
panel's x range is its own, because panel (b) has to carry the segment out to `x = 4` without leaving
its hole at the frame's edge, and panel (c) has nothing to show left of `x = -2.5`. And no tick
carries a number: the caption names the functions, the holes are drawn where they are, and nine
labeled ticks in a 153px panel would be a table rather than a picture.

This figure has a caption, so its alt is that caption minus its `Figure 2.12 ` label, the rule Figures
1.2 and 1.5 follow. Upstream's `data-alt` also exists, is longer, and is wrong in a way that matters:
it puts the second graph's holes at `(2, 1)` and `(-2, 1)`, and the second of those is not a point of
that graph, whose lower half-line is `y = -1`. Figures 5.5 and 6.5 record the same disagreement
between a drawing and its own description and resolve it toward the description, because there the
description was a figure that could be drawn. Here it cannot: no hole belongs at `(-2, 1)`. The
drawing follows the mathematics and the alt is the author's other sentence about it.

That one caption held one row, and the regenerated audit dropped it: 264 to 263, and 159 in the
certain tier. The lesson keeps other orphan captions, so the report's lesson count stays at 101, and
the ceiling for orphan captions in `internal/lessons/assets_test.go` follows the certain tier, to
263.

**Figures 2.13 and 2.14** (the two graphs whose captions say they confirm the table estimates:
`sin(x)/x` over `[-6, 6]`, empty at the origin because the function is not defined there and 1 is the
value it approaches; and `(sqrt(x) - 2)/(x - 4)` over `[0, 8]`, empty at `(4, 0.25)` for the same
reason) landed on 2026-10-16 in `calc.limit.numeric`, as `calculus/sine-over-x.svg` and
`calculus/radical-over-x.svg`. They are the first local figures whose ticks carry numbers, so the
labels are the only new thing the dialect needed: 11px, ink, every even integer on the x axis and
every twentieth on the y axis of 2.13.

The break in 2.14's y axis is the figure's real decision. Upstream runs 0.25 to 0.50 through about
two thirds of the frame, jogs under a zigzag, and gives the last thirteenth to 0.25 down to 0; the
curve is drawn through the same two-piece mapping, 364px per unit above the break and 72 below it.
The alternative, one linear axis from 0 to 0.5, puts the whole curve in the top sixth of an empty
frame and hides the flat tail -- which is the part that says the limit is 0. A break says the scale
changes here, which is true, and the zigzag is where it says it.

Both captions are mostly math, and an alt is plain text, so the caption's math is rendered into the
words it is made of by a small converter: `\frac{a}{b}` becomes `(a)/(b)`, `\sqrt{x}` becomes
`sqrt(x)`, `\text{sin}` becomes `sin`, and grouping braces and `$` come off. The converter asserts
on a macro it does not know, so no caption reaches a lesson as a backslash. For 2.14 the result is
`(sqrt(x) - 2)/(x - 4)`, which is how upstream's own `data-alt` writes that same function, and an
independent rendering that agrees with the author's is the check that this is a rendering and not a
paraphrase. 2.13's alt reads "The graph of f(x) = (sin x)/x confirms the estimate from Table 2.2."

Those two captions held two rows, and the regenerated audit dropped both: 263 to 261, and 157 in the
certain tier. The lesson keeps other orphan captions, so the report's lesson count stays at 101, and
the ceiling for orphan captions in `internal/lessons/assets_test.go` follows, to 261.

**Figures 2.15 and 2.16** (a curve whose value at `x = -1` is 4 while the curve itself passes through
3 there, drawn as a filled marker one unit above a hollow one; and a parabola whose vertex at
`(2, -1)` is the one point missing from it) landed on 2026-10-16 in `calc.limit.numeric`, as
`calculus/one-value-off-the-curve.svg` and `calculus/missing-vertex.svg`.

Neither caption gives a formula, which is the release to draw the shape the caption describes rather
than a function it does not name. 2.15's curve is therefore points chosen by hand and smoothed as a
Catmull-Rom chain -- steep out of the lower left, flat through the marker pair, steep again to the
upper right, which is the shape upstream's has -- and 2.16's is `0.25(x - 2)^2 - 1`, a parabola wide
enough that both arms and the vertex are inside a frame 180px tall. They are the first local figures
with a filled marker beside a hollow one, and their widths are set by the ratio the curve wants, 300
and 330, because a curve that rises five units over five units of x should not be drawn 3:1 wide.

**The batch's finding is that the audit could not see 2.16.** Its alt is upstream's `data-alt`
verbatim and that text contains "over [-2.5, 5]"; `IMAGE_OPEN_RE` in
`scripts/audit_lesson_media.py` stopped at the first `]`, so the audit did not read the line above
the caption as an image and counted the caption as an orphan anyway. `imageRe` in
`internal/lessons/assets.go` has carried a `\]` alternative since it was written -- its comment names
`numerical-integration.md`'s "The interval [a, b]" and `[n, n+1)`, whose bracket is never closed at
all -- so the Go test counted the same caption as wired. The two owners of one number disagreed by
one: 259 there and 260 here, on the same corpus. Section 1 quotes the audit, so the report was the
wrong one. The pattern now mirrors `imageRe` in both `IMAGE_RE` and `IMAGE_OPEN_RE`, and the four
vendored bracketed alts move no count at all: none of their captions is a row that a bracket is the
reason for. That the drift was invisible is the other half of the finding -- the ceiling is an upper
bound, so 259 and 260 both pass it, and the disagreement was only found by recomputing both plans
while adding a figure.

2.15 has a caption, so its alt is that caption minus its `Figure 2.15 ` label, with the `$g(x)$`
rendered as plain text by the converter 2.13 and 2.14 needed. 2.16 has the bare label upstream, so
its alt is upstream's `data-alt` verbatim, brackets and all.

Those two captions held two rows, and with the pattern fixed the regenerated audit drops both: 261 to
259, and 155 in the certain tier. The lesson keeps other orphan captions, so the report's lesson
count stays at 101, and the ceiling for orphan captions in `internal/lessons/assets_test.go` follows,
to 259.

**Figures 2.17, 2.18 and 2.19** (the last three graphs of the section: `sin(1/x)`, which oscillates
without settling as `x` approaches 0; a piecewise function with a break at `x = 2`, whose first piece
ends at an open circle at `(2, 3)` and whose second starts at a filled dot at `(2, 0)`; and `1/x` in
two branches) landed on 2026-10-16 in `calc.limit.numeric`, as `calculus/oscillation.svg`,
`calculus/piecewise-jump.svg` and `calculus/reciprocal.svg`.

The helper had to grow for them, and the growth is the reason they are one batch: this is the first
group needing ticks at tenths with labels at fifths (2.17), no labels at all (2.18, which upstream
also leaves unnumbered), and labels on both axes (2.19). The axes primitives take a tick list now
instead of looping over integers.

2.17 is the one with a method in it. `sin(1/x)` has no last oscillation, so sampling uniformly in `x`
spends its points where the curve is already smooth and starves the band. It is sampled uniformly in
`1/x` from 2 to 500 instead: every oscillation gets the same number of samples, the band closes into
the solid block upstream's has, and `x` down to 0.002 is as far in as a 280px frame shows anything
but a block. The count is 1,500 points per branch rather than 4,000 because the figure is inlined
into the lesson page — 4,000 was 124KB of path data and 1,500 is 47KB for the same picture. 2.18
draws its two pieces as two paths and marks both endpoints, which is the whole content of its
caption: the formula with the break is not at the left piece's end and is at the right piece's start.
2.19's branches leave the frame at `y = ±3.6` instead of being clipped into a chord, the same trick
2.12's third panel used.

The alts split three ways here, and the third way is new. 2.17 and 2.19 have captions that are prose
with a formula in them, so their alts are the caption minus the label with the formula rendered by
the converter 2.13 and 2.14 needed; the `*x*` stays, because the emphasis is the author's and the
renderer flattens it in the attribute. 2.18's caption is not prose at all: it is
`f(x) = \left\{ ... \right.`, a piecewise definition inside a matrix, and the converter refuses a
`\begin` rather than inventing words for a formula, which is the assert doing what it was written to
do. So 2.18 takes upstream's `data-alt` verbatim, and that alt describes the figure rather than
restating it: "The first piece is f(x) = x+1 if x < 2. ... There is an open circle at (2,3), where the
endpoint would be." The rule that follows is that a caption is a source for an alt while it is a
sentence about the figure; when the caption *is* the figure's formula, the `data-alt` is the sentence
and the caption is not.

Those three captions held three rows, and the regenerated audit drops all three: 259 to 256, and 152
in the certain tier. The lesson keeps other orphan captions, so the report's lesson count stays at
101, and the ceiling for orphan captions in `internal/lessons/assets_test.go` follows, to 256.

**Figures 2.21 and 2.23** (the figure an exercise is read off -- four limits and four values asked at
`x = -4, -2, 1` and `3`, with a curve interrupted twice and a vertical asymptote in the middle of it
-- and Einstein's mass ratio against speed ratio) landed on 2026-10-16 in `calc.limit.numeric`, as
`calculus/interrupted-curve.svg` and `calculus/mass-ratio.svg`. They are the last two redrawable
figures of Calculus 2.2.

2.21 is drawn in the alt's order rather than as one function, because that is what it is: a spline
below `x = -2` ending at an open circle at `(-2, 3)`, a second spline from there to a filled dot at
`(1, 6)`, an open circle at `(1, 3)` with a branch running down to `-infinity` along `x = 3`, and the
mirror branch on the other side of the asymptote, which comes back up through the concavity change
the alt names in the first quadrant. Only the two interesting y values are numbered, 3 and 6, and
only the four x values the exercise asks about -- upstream's choice, and the right one: numbers
everywhere would have a reader checking the graph instead of reading it.

2.23's axes are named for its quantities, `v/c` and `m/m0`, and the formula is written inside the
frame the way upstream writes it, as a fraction with a bar and the root below it. That covers the
whole figure: there is no tick label above 5 and no second series. `m0` is written with a plain
zero and not U+2080, because no theme font carries a subscript zero and a substituted face would sit
at a different size and baseline inside the drawing -- which is the objection AGENTS.md section 1
makes about a codepoint the fonts do not carry, and it applies to a subscript as much as to a
pictograph. The two superscripts are U+00B2, the one the arithmetic figures already use.

Both have captions that are prose, so both alts are the caption minus its label: 2.21's is "The graph
shows f(x)." and 2.23's is the sentence about Einstein's equation, kept with the author's
typographic apostrophe.

Those two captions held two rows, and the regenerated audit drops both: 256 to 254, and 150 in the
certain tier. The lesson keeps other orphan captions, so the report's lesson count stays at 101, and
the ceiling for orphan captions in `internal/lessons/assets_test.go` follows, to 254. With 2.23
landed, every redrawable figure in Calculus 2.2 is drawn: 2.12 through 2.19 and 2.21, with 2.20 and
2.22 remaining as fetched illustrations and a photograph.

**Figures 3.25, 3.26, 3.27 and 3.33** (the difference quotient the sine rule is guessed from; the
two limits its proof needs; sine and cosine in one frame with dots where the tangent is
horizontal; and the family `2^x`, `2.7^x`, `2.8^x`, `3^x` that places e) landed on 2026-10-16 in
`calc.deriv.trig` and `calc.deriv.exp_log`, as `calculus/difference-quotient.svg`,
`calculus/two-limits.svg`, `calculus/sine-and-cosine.svg` and `calculus/exponential-family.svg`.
They are the first figures outside section 2, and the first two lessons here to gain a figure at
all.

Three of the four needed the same layout fix, which makes it a rule: the frame has to run past its
last tick label, or the label, the arrowhead and the axis name pile up on top of each other. 3.25
runs to `2pi + 0.7`, 3.27 to `2pi + 0.55` and 3.33 to `0.72`.

3.25 draws `D(x) = (sin(x + 0.01) - sin x)/0.01` over `[-2pi, 2pi]`, and repeats the formula below
the frame rather than inside it, because the lesson writes it as a display equation and upstream
prints it under the plot.

3.26 is two panels with one curve each. Upstream gives each panel its own y range; this gives them
one, `[-0.85, 1.05]`, so that their x axes land in one row and the two limits can be read against
each other, which is the only reason to show them together. The cost is real and worth naming:
`(cos h - 1)/h` has amplitude 0.73 and is drawn at 82% of the size its own range would give it, and
the shared range then holds labels every 0.5 instead of upstream's every 0.2. The alt says the two
curves "have discontinuities on the y-axis", which is what the two hollow markers at `(0, 1)` and
`(0, 0)` are.

3.27 is the first local figure to use the trig accents section 4 of the style doc reserves for a
function family, `#2980b9` for the sine and `#e74c3c` for the cosine: two curves that cross four
times in one frame cannot be told apart by weight alone. The dots are the caption's statement --
the two extrema of the sine and the two zeros of the cosine at those same x values -- plus the
cosine's maxima at `x = 0` and `x = 2pi`, which are upstream's and are kept because a look-alike is
the goal. Upstream's dot set is not symmetric, and this record says so rather than inventing a dot
at `(pi, -1)` to match: there is no dot there.

3.33 has the x axis on the frame's bottom edge instead of at `y = 0`, because its y range starts at
0.5, and the y axis carries a break glyph to say so -- the glyph 2.14 introduced for a two-piece
scale, here for one that simply starts above zero. Its four curves all end at the right edge so
each label sits at its own curve's end, and the two the caption names (`2^x` and `3^x`) are ink at
width 2 while the pair that brackets e is annotation grey at 1.5. The labels read `3^x`, not a
raised superscript, for the same reason 2.23's `m0` is a plain zero.

**No new e2e case.** None of the four introduces a primitive the sanitizer has not already been
shown to pass: a colored stroke is the same `stroke` attribute the ink curves carry, and the break
glyph and the bottom x axis are shapes 2.14 and 2.12 already draw. The screenshot check carries
them, and `calc.deriv.trig` is a lesson that had no figure before today, which the
`arith.mult.tables` case already covers as a pattern.

**The lesson count moved, not only the row count.** Figure 3.33 was `calc.deriv.exp_log`'s only
orphan caption, so with it wired that lesson leaves the report entirely: **101 lessons to 100**.
Every earlier batch took rows off the list and left the lesson count alone, so this is the first
time the third number in section 1 has moved, and the first batch where the Go test's own log line
and the audit disagree with the docs unless all three are updated together.

Those four captions held four rows: the report went from 254 to 250, and the certain tier from
150 to 146. The ceiling for orphan captions in `internal/lessons/assets_test.go` follows, to
250. (The numbers are written in this order rather than as one live claim because the checker
reads a count followed by the word it counts: a dated record must not look like today's
measurement, or every later batch has to rewrite history to keep the checker quiet.)

**Figures 4.40, 4.42 and 4.45** (a curve approaching the horizontal asymptote `y = 2`; a curve that
crosses its horizontal asymptote an infinite number of times; and `tan x` with vertical asymptotes
at `x = +-pi/2`) landed on 2026-10-16 in `calc.limit.infinity`, as
`calculus/horizontal-asymptote.svg`, `calculus/crossing-asymptote.svg` and
`calculus/tangent-asymptotes.svg`. This is the first batch whose subject is the asymptote itself,
and the first whose third figure goes into a numbered list item rather than a paragraph.

**The dashed guide is new to the local set, and already proven.** No local calculus figure used
`stroke-dasharray` before today; the vendored Algebrica set does, and three of the lessons the e2e
list already checks carry one (`alg.log.concept`, `calc.deriv.cauchy_mvt`, `calc.seq.euler`), so the
attribute is one the sanitizer has been shown to pass for months. That is why this batch adds no
e2e case: the new shape is new *here* and not new to the sanitizer, which is the distinction the
rule was written for. The guides themselves are annotation grey `#636363` rather than the axes'
`#BABABA`: a labeled asymptote is an annotation, not an axis.

**4.40's caption does not name its function, and the figure still had to be drawable.** "The
function approaches the asymptote `y = 2` as `x` approaches `+-oo`" is prose about a curve, so the
alt is that sentence and the curve has to come from somewhere else: the lesson's own prose later in
the section ("we used graphical evidence in Figure 4.40 and numerical evidence in Table 4.2 to
conclude that `lim (2 + 1/x) = 2`) and upstream's `data-alt`, which reads "The function f(x) 2 +
1/x is graphed". Taking the function from the lesson rather than the caption is a decision, and it
is the reason this row was classified redraw while its neighbour 4.43, whose caption names no
function *and* whose prose names it only in the surrounding paragraph, was classified fetch.

**Both 4.40 and 4.42 draw one path per branch, from analytic domains, with no clip path.** The
right branch of `2 + 1/x` starts where the curve leaves the top of the frame, at `x = 1/3`, and the
left branch ends where it leaves the bottom, at `x = -1/7`; 4.42's branches end where the curve
leaves the frame at `y = 2.2` and `y = -1.4`, found by bisection. That is how `reciprocal.svg` drew
`1/x` and it means no path ever leaves the box, which is visible rather than merely contained.
4.40 also needed a wider x range than its ticks: labels every 2 instead of upstream's every 1,
because at 14.4px per unit a two-character label does not fit between two ticks. 4.26's record
already made that trade; this is the second time.

**4.45 is the first figure in the set that a numbered list item carries.** Its caption is indented
four spaces upstream, so the insertion preserves the indentation; an image at column 0 would have
ended the list item and left the caption orphaned in the markdown in a way the audit would still
count as wired. The label placement is worth recording too: `y = tan x` sits *left* of the y axis
between the two asymptotes, not right of it as upstream has it, because at 45.9px per unit the free
span right of the axis ends where the curve crosses the label's own row, and the label is 59px
wide. The `x = pi/2` label had to drop below the x axis for the same reason in reverse: at row 97 it
landed on the "2" and "3" tick labels, which sit at row+15.

**Three new macros and one rule in the alt converter.** 4.40's caption ends
`$\text{\pm}\infty$` and 4.45's reads `$x = \text{\pm}\frac{\pi}{2}$`, so the converter now
renders `\pm`, `\infty` and `\pi` instead of asserting on them, and it no longer brackets a
fraction whose numerator and denominator are each a single token: 4.45's alt reads `x = +-pi/2`,
not `x = +-(pi)/(2)`. A part with an operator inside keeps its parentheses, so 2.14's committed alt
is unchanged and the wire script checks that against the string in the lesson rather than against a
comment. One further fix came out of the same caption: upstream writes `\left( {\text{cos} x}`,
and the space after the opening parenthesis leaked into 4.42's alt as `( cos x )/x`, so the
converter now closes up whitespace inside parentheses.

**The checker caught a stale count in the previous record, which is worth writing down.** After
this batch's numbers moved, `scripts/check_docs.py` failed on the 3.25-3.33 record, which ended its
row count with the number 250 and the word "candidates" straight after it: the checker reads any
count followed by the word it counts as a live claim, and a dated record is not one. The fix was
to phrase the history as a sequence ("went from 254 to 250") rather than as a measurement, and then
to write this paragraph so that it did not re-create the pattern while describing it. Every later batch has to keep that distinction, or the records have
to be rewritten to match today's corpus.

Those three captions held three rows: the report went from 250 to 247, and the certain tier from
146 to 143. The lesson keeps nineteen other orphan captions, so the lesson count stays at 100, and
the ceiling in `internal/lessons/assets_test.go` follows, to 247.

**Figures 4.50, 4.51 and 4.52** (the even powers `x^2`, `x^4` and `x^6`; the odd powers `x`, `x^3`
and `x^5`; and a polynomial over its leading term) landed on 2026-10-16 in `calc.limit.infinity`,
as `calculus/power-even.svg`, `calculus/power-odd.svg` and `calculus/polynomial-leading-term.svg`.

**4.50 and 4.51 are the first local figures to use the accents for something that is not
trigonometric.** Section 4 of the style doc gives the rule's scope as "where a function is compared
across a family" and names the three colors after the functions that first needed them. This batch
needs them for a measured reason: at the frame's top edge 4.50's `x^4` and `x^6` curves exit 12px
apart, and they are closer at every row below it, so no set of positional labels can identify them —
and their left-to-right order flips at `x = +-1`, so a stacked set cannot either. 3.33 manages four
exponential curves without color because its labels sit in the right margin, each at the height
where its own curve leaves the frame; these curves all leave through the top. Upstream encodes the
same three curves by color, and the vendored set colors text as well as strokes (`gen-diagrams.py`
writes `#2980b9` labels), so here the label and its curve carry the same accent.

**4.52 takes no accent, and that is the same rule read the other way.** Its two curves coincide at
both ends — that is the caption's statement — so they are told apart by the ladder
`exponential-family.svg` uses: `#333333` at width 2 for `f(x) = 5x^3 - 3x^2 + 4`, and `#636363` at
1.5 for `g(x) = 5x^3`, the term it behaves like. A reader who cannot see the color still has two
labels, and each label says which function it belongs to.

**The pair shares one x scale, and that is what buys the labels their room.** Both frames are 420px
wide over `[-5, 5]`, or 36.6px per unit against upstream's 38.4. 4.51 needed it: at `y = 10` the
`x^5` curve's right branch is 58px from the y axis, so a 46px label fits beside it, where at a
360px canvas the gap was 45px and the label had to overlap the axis's own name. 4.50's three labels
are stacked in the bottom-right wedge instead of upstream's top-right corner, by the same kind of
measurement: that corner has 14px between the `x^2` curve's exit and the frame's edge, and the wedge
has 35px at the top row of the stack. One more row moved after looking at the render, which is the
only way this class of defect shows up: 4.51's `y = x^5` label sat at baseline 27, on the rows the y
axis's own name and arrowhead occupy, and it is at 34 now.

**The converter learned two macros, and one of them was silently wrong first.** 4.50 and 4.51 state
their end behavior in OpenStax's display notation, `\underset{x\rightarrow\infty}{\text{lim}}x^{n}`.
`\underset{A}{B}` renders as `B(A)`: dropping A would lose which limit is meant, and appending it
bare reads as `lim x -> oo x^n`, ambiguous about where the argument starts. `\rightarrow` renders as
` -> ` with LaTeX's own relation spacing, which the captions do not write. The first attempt put the
`\rightarrow` branch after the `\left`/`\right` branch, and since `\rightarrow` starts with
`\right`, that branch ate it and produced `xarrowoo` in all three alts. The wire script asserts the
three expected strings rather than printing them, which is why the defect did not reach the corpus;
that assert is the unknown-macro assert one level up.

Those three captions held three rows: the report went from 247 to 244, and the certain tier from
143 to 140. The lesson keeps sixteen other orphan captions, so the lesson count stays at 100, and
the ceiling in `internal/lessons/assets_test.go` follows, to 244.

**Both frames ended on a labeled tick, and the last label ran into the axis name.** 4.50's and
4.51's x axes were drawn over `[-5, 5]` with every integer labeled, so the last label sat on the
frame's right edge, centred under the arrowhead, and the axis's own `x` -- 3px past the arrow --
landed on it. At 600% the render reads `5x`, one word. The rule the 3.25-3.33 batch derived is that
the frame runs past its last tick label, and 4.40 and 4.42 obey it: their frames reach `+-9` and
`31.5` while their labels stop at `+-8` and `30`. Both figures now label `-4` to `4` and leave
`+-5` unlabeled, which is 4.40's shape. Nothing else moved; the range, the scale and the label
stack are as recorded above.

**Figures 4.53, 4.54 and 4.56** (a linear over linear with a horizontal asymptote; a quadratic over
a cubic whose asymptote is the axis itself; a quadratic over a linear with an oblique asymptote)
landed on 2026-10-16 in `calc.limit.infinity`, as `calculus/linear-over-linear.svg`,
`calculus/quadratic-over-cubic.svg` and `calculus/quadratic-over-linear.svg`.

**Two of the three sit inside a numbered list, and the insertion keeps the indent.** Upstream's
Example 4.25 numbers three functions and puts each figure inside its item, so 4.53's and 4.54's
captions are indented four spaces and their images are too. That is 4.45's case again, and the rule
is the same: the indent is read off the caption line and written back, never assumed.

**4.54 carries no guide, because its asymptote is the x axis.** The caption says the asymptote is
`y = 0`, and a dashed line drawn on the axis would be a second line drawn on the first. Upstream
draws none either. The curve's visible features are the pole at the cubic's only real root (about
`-1.546`), where the left branch leaves the bottom and the right branch leaves the top, and a hump
of about `0.85` at `x = 1`.

**4.56 is the first asymptote here that is not axis-aligned**, so its guide is a dashed path along
`y = 3x - 2` between the rows it leaves the frame rather than `guide_h` or `guide_v`. Its label is
ink, as 4.40's `y = 2` and 4.45's `x = pi/2` are; only the guide itself is annotation grey. The
label sits right of both exits rather than in the wedge between the curve and the line, and that is
a measurement: the gap is `4/(x + 2)`, which is 6px at `x = 0` and 3px at `x = 3` at this y scale.
The scale is small because the figure's point is that the curve hugs the line far out, so the frame
reaches `+-28`. Upstream's canvas is 347px tall and this one is 180, and the 6px wedge is what that
costs.

**One analytic domain was wrong first, and the generator caught it rather than the eye.** 4.56's
right branch does not leave the bottom: `4/(x + 2)` is positive there, so the branch enters the
frame at the top just right of the pole and leaves it at the top again, one U-shaped path rather
than two branches. The first attempt bisected an interval that crossed the pole and divided by zero.

**The converter learned one macro, `\ `, a control space.** 4.56's caption reads
`... 2\ \text{as}\ x ...`, and the unknown-macro assert fired on it before any alt was written. It
renders as a space, which cannot change an alt that is already shipped, because those are written
rather than derived.

Those three captions held three rows: the report went from 244 to 241, and the certain tier from
140 to 137. The lesson keeps thirteen other orphan captions, so the lesson count stays at 100, and
the ceiling in `internal/lessons/assets_test.go` follows, to 241.

**Figures 4.58 and 4.59** (sine over a little more than two periods; tangent over
`[-3pi/2, 3pi/2]`, where the limit fails) landed on 2026-10-16 in `calc.limit.infinity`, as
`calculus/sine-wave.svg` and `calculus/tangent-periodic.svg`.

**4.59's x tick labels are stacked fractions, and that is a layout decision rather than
decoration.** Upstream writes `-3pi/2` as a fraction with a bar, and inline it is five characters.
The frame has to stop at about `+-4.93`: left of `-4.931` the branch before `-3pi/2` is above the
frame, and right of `5.093` the branch after `3pi/2` is below it, so a wider frame shows a fourth
branch that upstream's figure does not have. At that width a five-character label on the last tick
runs into the axis's own `x`. Stacked, the widest numerator is `3pi` plus its sign, and the bar is
as wide as the numerator plus a little, which is how upstream's bars look. The single-line labels
`-pi`, `pi` and `0` sit on the fraction bars' row, which is also how upstream aligns them.

**4.59 draws no asymptote guides.** A full-height guide at `pi/2` runs through the stacked label
that names it, and the branches' own vertical runs already mark the asymptotes; upstream draws none
either. 4.45, the lesson's other tangent figure, does draw them, and the difference is not
inconsistency: its tick labels are integers, so its guides are clear, and its caption is about the
asymptotes where this one is about the limit.

**4.58 labels the axes in `pi`, not in integers.** Upstream labels `-4` to `4` under a sine of
period `2pi`, which reads as a period of six. `sine-and-cosine.svg`, the tree's other sine figure,
labels `pi` and `2pi`, and this figure is more useful agreeing with it.

Both captions carry `\ ` and neither ends in a full stop, so neither alt does either.

Those two captions held two rows: the report went from 241 to 239, and the certain tier from 137 to
135. The lesson keeps eleven other orphan captions, so the lesson count stays at 100, and the
ceiling in `internal/lessons/assets_test.go` follows, to 239.

**Figures 4.3, 5.32 and 5.33** (the airplane triangle; speed against time with the two rectangles of
area; a linear velocity whose area changes sign at `t = 5/3`) landed on 2026-10-16, in
`calc.deriv.related_rates` and `calc.integral.power_rule`, as `calculus/airplane-triangle.svg`,
`calculus/speed-time-areas.svg` and `calculus/linear-velocity-areas.svg`.

**4.3 is a schematic, not a copy.** The upstream image is a drawing of a person, an airplane and a
radio tower, and the disposition for this tier is Mathua's own work, so the figure is the right
triangle the caption describes: `x` on the ground, `4000 ft` up the tower, `s` along the
hypotenuse, the ground drawn past the triangle because it is the ground, and an arrow at the
airplane's vertex because the height is constant and the horizontal distance is not. The three
quantities are named in the caption, so the caption is the alt.

**The two area figures shade with a filled path and no stroke, drawn before the axes.** That is the
vendored idiom: `algebrica/definite-integrals-*.svg` closes the region, fills it, leaves it
unstroked, then draws the axes over it and the function over both. The fill is the dialect's
`#E6E6E6`; upstream's blue is not.

**Their `0` is on the y axis, not the x axis.** 5.32's y labels are 16.5px apart and 5.33's are
11.8px, because the dialect's canvas is 180px tall and neither range is small. A `0` on the x axis
sits to the left of the y axis, which lands it on the `-10` label in 5.32 and the `-1` label in
5.33. Upstream puts `0` on the y axis in both figures and `sine-and-cosine.svg` does the same, so
these two do too. The y tick at zero draws a horizontal tick along the x axis, where it is
invisible.

Those three captions held three rows: the report went from 239 to 236, and the certain tier from 135
to 132. Both lessons keep other orphan captions, four in the related-rates lesson and three in the
power-rule lesson, so the lesson count stays at 100, and the ceiling in
`internal/lessons/assets_test.go` follows, to 236.

**Figures 5.35 and 5.36** (3x^8 - 2 at two scales, where the area is positive outside its roots and
negative between them; 5 sin x over `[-pi, pi]`, whose areas cancel) landed on 2026-10-16 in
`calc.integral.power_rule`, as `calculus/power-positive-negative.svg` and
`calculus/odd-function-areas.svg`.

**5.35 is the tree's second two-panel figure.** `two-limits.svg` set the layout: a 580-wide canvas,
two panels, `(a)` and `(b)` centred under each at row 158. The curve is drawn only where it is
inside each panel's frame, stopping at `y = 640` in (a) and `y = 5.6` in (b), because past those
points it is off-frame and a path would paint over the margins. The panels show the same curve at
two y scales, which is what makes the caption's claim about positive and negative area visible at
once, and the dashed guides at `x = +-2` in (b) mark the interval (a) is drawn over.

**5.35 carries one function label, in (a).** (b) has no room: its y tick labels march down the axis
at 15.5px, the `y` name sits at 454, and the guides at `-2` and `2` leave 64px between them for a
92px label. The label in (a) sits at row 44, which clears the `y` name above it and the curve's
exit beside it, the curve reaching row 44 at `x = 1.89`.

**5.36 draws 5 sin x, and upstream's figure is mirrored relative to the function its own label
names.** Upstream's shaded positive area lies between `-pi` and 0, where 5 sin x is negative; the
lesson computes the integral with 5 cos x as the antiderivative over `[-pi, pi]`, so the positive
area is between 0 and `pi`. The figure follows the lesson, as 4.54 followed the lesson's `4x^3`
over upstream's `4x^2` in its alt text.

Those two captions held two rows: the report went from 236 to 234, and the certain tier from 132 to
130. The lesson keeps one other orphan caption, so the lesson count stays at 100, and the ceiling
in `internal/lessons/assets_test.go` follows, to 234.

**Figures 6.25 and 6.27** (the region and the solid it sweeps, and the strip that generates one
shell) landed on 2026-10-16 in `calc.integral.volume`, as `calculus/revolved-cup.svg` and
`calculus/shell-strip.svg`. They are the first figures in the tree whose second panel is a solid
rather than a graph, so they establish the tree's 3D idiom.

**A solid of revolution is drawn in axonometric, in the same dialect as a flat shape.** Its circular
ends are ellipses with a fixed `rx:ry` ratio of 7:2. The rim is a full ellipse, because we look down
into the solid; the base is a near arc stroked in ink with the far arc dashed `#636363`, because the
far half is hidden. The silhouette is one closed path filled `#E6E6E6`, the same grey as every other
filled shape in the tree, and there is no gradient and no shading.

**A curved top surface is a bowl, and a bowl is nested level ellipses.** The curve `y = f(x)`
revolved about the y-axis gives the surface `y = f(r)`, so the circle of radius `r` sits at height
`f(r)`: its ellipse is centered at `row(f(r))` with both radii scaled by `r/b`. Two interior levels
are enough. They nest because both their centers descend and their radii shrink, and a smaller
level's near arc can fall below the rim's near arc, which is correct for a bowl and is how the depth
reads.

**6.25's `a` is 3.0, not 1.6.** `a` is where the curve starts on the y-axis, so it is the depth of
the bowl: at 1.6 the revolved top surface was a deep cup, and upstream's is a shallow dimple.

**6.27's strip is filled and stroked, where the area figures' regions are fill-only.** Its top edge
is the sample height `f(x_i*)` and not the curve, so it is a named rectangle rather than an area
under a curve. 5.32's rectangles stay fill-only because their top edges *are* the function. 6.27
also draws no ticks: its two ends are named `x_(i-1)` and `x_i`, so there are no numbers to place.

Those two captions held two rows: the report went from 234 to 232, and the certain tier from 130 to
128. The lesson keeps ten other orphan captions, so the lesson count stays at 100, and the
ceiling in `internal/lessons/assets_test.go` follows, to 232.

**Figure 6.26** (the rectangle, the shell it sweeps, and all the shells together) landed on
2026-10-16 in `calc.integral.volume` as `calculus/shell-approximation.svg`, and it adds the tree's
second rule for a solid.

**A tube's hole is white, and the inner wall between its two ellipses is not filled.** Upstream
distinguishes the inner wall from the body with a shade, and the dialect has no shade: a grey inner
wall is indistinguishable from a grey body, so the tube stops reading as hollow, which contradicts
its own alt text ("a cylindrical shell, hollow in the middle"). So the hole is the inner cylinder's
whole silhouette filled white -- not just the rim ellipse, because the region between the inner
rim's near arc and the inner base's near arc is part of the hole -- and two ellipses are drawn in
it: the inner rim in ink, because it is the edge of the hole, and the inner base in annotation grey,
because it is the bottom opening seen through the hole. So (b) is the outer silhouette, one white
fill, and two ellipses.

**A panel label moves to row 172 when the panel holds a solid.** `two-limits.svg` puts `(a)` and
`(b)` on row 158, and that is right for a figure made of curves. A solid's base is a near arc that
reaches `cy_base + ry`, which at `cy_base = 140` is row 158, so the label lands on the arc. The same
defect was in Figure 6.25 and is fixed alongside this one. The rule is about the arc and not the
row: 158 stays free for a figure with nothing drawn below its axis.

**Panel (c) is stepped, not flat.** Each shell of radius `r_k` has height `f(r_k)`, so the stack's
top is a staircase descending inward, which is what the caption means by "an approximation of the
original solid". Upstream's own (c) is hard to read at this size and may draw its rings level; the
caption decides it.

That caption held one row: the report went from 232 to 231, and the certain tier from 128 to 127.
The lesson keeps nine other orphan captions, so the lesson count stays at 100, and the ceiling in
`internal/lessons/assets_test.go` follows, to 231.

**Figures 6.28 and 6.29** (the shell cut open and flattened, and the solid that `1/x` generates about
the y-axis) landed on 2026-10-16 in `calc.integral.volume`, as `calculus/shell-cut-plate.svg` and
`calculus/one-over-x-solid.svg`.

**6.28(a) reuses Figure 6.26(b)'s tube rather than redrawing one.** The two figures sit four
paragraphs apart, so the shell is the same shape in both -- rim radius 62, ellipse ratio 7:2, base
row 140, height 78 -- and the top row lands within two pixels of 6.26's, because each figure samples
the curve at its own point. A redrawn tube would not be. Upstream labels the curve on the top `f(x)`
where the tree's other figures write `y = f(x)`, and this one follows the tree.

**A dimension line is annotation grey, carries an arrowhead at each end pointing out at the object,
and each end is met by a short extension line.** 6.28(b) measures the plate's three sides, and a
measured side needs the marks of a ruler rather than the ticks of a graph: upstream draws no axes
there, and the plate is named `2πx_i*`, `f(x_i*)` and `Δx` rather than by coordinates. The extension
lines are what say which two points the dimension spans; an arrowhead on its own floats.

**A closed silhouette that runs right to left along an arc is built from explicit arc points, not
from an `A` command.** An arc's sweep flag is read against the direction of travel, so the flag that
draws the outer wall left to right draws the inner wall inside out when the same path comes back.
6.29(b) is one closed path around a wall with a hole through it, and both its outer arc and its inner
arc are sampled.

**The cut line is ink, not upstream's blue.** The dialect has no accent that means "this is a cut",
and a color used once is decoration; the line is stroked in ink and a leader carries the words `Cut
line`, which is how upstream names it too.

**The solid closes on a circular edge rather than on a rim.** At `y = 1` the outer radius equals
the hole's radius, because `1/x` at `x = 1` is the region's height and `x = 1` is also the hole, so
the outer wall meets the hole's wall on one circle. The top ellipse is therefore both the outer
surface's edge and the hole's rim, and Figure 6.26's rule -- a full rim ellipse plus a white hole --
collapses to a single ellipse here. The outer wall is not a cylinder either: it is vertical to
`y = 1/3`, where the region reaches `x = 3`, and then follows `1/y` in to the hole.

**6.29 is the first caption in the corpus to close an interval**, and its alt spells `[1,3]` where
the LaTeX writes `\left\lbrack {1,3} \right\rbrack`.

Those two captions held two rows: the report went from 231 to 229, and the certain tier from 127 to
125. Distinct figures go from 57 to 55 and redraws from 8 to 6. The lesson keeps seven other orphan
captions, so the lesson count stays at 100, and the ceiling in `internal/lessons/assets_test.go`
follows, to 229.

**Figure 6.30** (the region under `f(x) = 2x - x^2` and the solid it sweeps about the y-axis)
landed on 2026-10-16 in `calc.integral.volume`, as `calculus/two-x-minus-x-squared-solid.svg`.

**The parabola's two branches are the solid's two walls.** The right branch is the outer wall,
because `x = 1 + sqrt(1 - y)` is the region's far edge at height `y`, and the left branch is the
funnel, because `x = 1 - sqrt(1 - y)` is its near edge. So the walls meet at `y = 1`, where the
solid's top is a circular knife edge and not a rim, and the hole closes on the base's center at
`y = 0`, where the inner radius is 0 and the base is a full disc.

**A hole that closes inside the solid is drawn as the whole silhouette of the void, exactly as a
through-hole is.** 6.30(b)'s funnel is the apex, the two generatrices and the rim's far arc, filled
white. Its lower half is hidden behind the dome's near surface in a real view, and 6.26 already
settled that the tree draws a void's whole silhouette rather than the part that is not occluded, so
6.30 follows 6.26 and not the occlusion. The rim is still a full ellipse in ink, so its near arc
crosses the white.

That caption held one row: the report went from 229 to 228, and the certain tier from 125 to 124.
Distinct figures go from 55 to 54 and redraws from 6 to 5. The lesson keeps six other orphan
captions, so the lesson count stays at 100, and the ceiling in `internal/lessons/assets_test.go`
follows, to 228.

**Figure 6.31** (the region Q to the left of `g(y) = 2 sqrt(y)` and the solid it sweeps about the
x-axis) landed on 2026-10-16 in `calc.integral.volume`, as `calculus/paraboloid-hole-solid.svg`.

**The tree's oblique projection is one projection, with the depth drawn along either of the other
two axes.** A y-axis solid's circular ends are horizontal circles, so the depth is drawn vertically
and an end comes out as an ellipse 7 wide to 2 tall; 6.25 to 6.30 all read that way. An x-axis
solid's ends are vertical circles, so the depth is drawn horizontally and the same end comes out 2
wide to 7 tall -- the same projection, rotated with the axis it belongs to. Without that, a circle
perpendicular to x projects to a vertical segment, which is why a horizontal cylinder cannot be
drawn at all in the first reading, and 6.31(b) is the tree's first figure with a horizontal axis.

**The hole's whole silhouette is filled white, and the far end's disc with it.** The hole is the
paraboloid of radius `x^2/4`: it closes at `x = 0`, where the inner radius is 0, and its radius at
`x = 4` equals the cylinder's, so the far end is a knife ring whose disc is the hole's mouth and not
material, while the near end is a full disc and stays grey. The two silhouette curves are
`y = +-x^2/4`; upstream dashes the upper and strokes the lower, and that is kept -- the upper is the
hole's far edge and the lower its near edge.

**Both panels keep equal scales, so the canvas is 400 wide and not 580.** The region is a 4 by 4
square, and a parabola drawn with unequal scales is a different parabola; the solid's diameter is
twice its length, and with a panel filled it would read as a long cylinder rather than a squat one.
The style doc sizes a canvas to its content and the calculus figures already range from 210 to 780
wide, so 400 is the honest width here and not 580 with 118px of margin on each side.

The caption is terse -- it says `$g(y)$` and not the function -- so the alt says "the function
g(y)" and the figure's own curve label carries `g(y) = 2` and a root. That caption held one row:
the report went from 228 to 227, and the certain tier from 124 to 123. Distinct figures go from 54
to 53 and redraws from 5 to 4. The lesson keeps five other orphan captions, so the lesson count
stays at 100, and the ceiling in `internal/lessons/assets_test.go` follows, to 227.

**Figure 6.32** (R between `f(x) = x` and the x-axis over [1, 2], revolved about the line
`x = -1`) landed on 2026-10-16 in `calc.integral.volume`, as `calculus/bowl-offset-axis-solid.svg`.

**The axis of revolution need not be a coordinate axis.** Panel (b) is centred on the line
`x = -1` rather than on the y-axis, so that line is drawn as an annotation -- annotation grey,
solid where it is outside the solid and dashed where the solid hides it, with no arrowhead,
because it is not an axis -- and the coordinate axes are still drawn at `x = 0` and `y = 0`. The
y-axis crosses the solid and keeps its dashed hidden span, and so does the axis of revolution. The
label that names it goes inside the cavity, on the white, because the top of the panel is where the
rim's far arc already is.

**A cavity can have a floor, which 6.26 and 6.30 did not have.** The solid is an outer cylinder of
radius 3 from `y = 0` to `y = 2` whose cavity is a cylinder of radius 2 from `y = 0` to `y = 1` and
then opens as `y + 1`, so the top is a knife edge and the cavity's floor is a real disc. Its rim is
one annotation-grey ellipse, the way 6.26 draws a hole's inner base, because at this projection's
angle the floor is not visible: a ray from the deepest point of it leaves the cavity through the
wall. The cavity's right slant edge is the line `x = y`, which is why upstream draws `f(x) = x`
through panel (b) as well, and that edge carries the label `y = f(x)`.

Both panels keep equal scales, because the region's upper boundary is the line `y = x` and a line
drawn with unequal scales is a different line, and because the solid is 6 units across and 2 tall,
which a filled panel would redraw as a deep bowl.

That caption held one row: the report went from 227 to 226, and the certain tier from 123 to 122.
Distinct figures go from 53 to 52 and redraws from 4 to 3. The lesson keeps four other orphan
captions, so the lesson count stays at 100, and the ceiling in `internal/lessons/assets_test.go`
follows, to 226.

**Figure 6.33** (R between `f(x) = sqrt(x)` and `g(x) = 1/x` over [1, 4], revolved about the
y-axis) landed on 2026-10-16 in `calc.integral.volume`, as
`calculus/two-curve-waisted-solid.svg`.

**A cavity bounded by two functions is a waisted silhouette, not a funnel.** The region's right edge
is `x = 4` at every height, so the outer surface is a plain cylinder of radius 4 and both ends are
knife edges: the inner radius reaches the outer wall at `y = 1/4` and again at `y = 2`. The inner
radius is `max(1/y, y^2)`, which is `1/y` below `y = 1` and `y^2` above it, so the cavity is widest
at both ends and narrowest at the waist, and the material is two wings that are widest at the waist.
There is no floor and no inner rim to draw: the void's whole silhouette is one white closed path
bounded by the two cavity curves, the rim's far arc and the base's near arc. This is the figure
6.26 and 6.30 had no counterpart for, because their voids were a straight hole and a funnel that
closes at a point.

**The white fill goes in after the hidden far base arc and before the near base arc and the rim.**
The void's silhouette meets the rim's ellipse and the base's near arc exactly, so a white fill drawn
last shaves half the stroke off each and they read thin. Drawing the hidden far base arc first keeps
the void's white over it, which is what hides it.

Both panels keep equal scales, because panel (a) is read against its axes and panel (b) is 8 units
across and 1 3/4 tall, which a filled panel would redraw as a tall shell.

That caption held one row: the report went from 226 to 225, and the certain tier from 122 to 121.
Distinct figures go from 52 to 51 and redraws from 3 to 2. The lesson keeps three other orphan
captions, so the lesson count stays at 100, and the ceiling in `internal/lessons/assets_test.go`
follows, to 225.

**Figures 6.35 and 6.36** (6.35: the region between `y = x` and `y = 2 - x`, a triangle whose apexes
lie on the axis; 6.36: the region between `y = 4x - x^2` and the x-axis over [0, 4]) landed on
2026-10-16 in `calc.integral.volume`, as `calculus/bicone-x-axis-solid.svg` and
`calculus/lens-x-axis-solid.svg`.

**An x-axis solid is the region and its mirror.** Neither region has an end of non-zero radius, so
neither solid has an end ellipse: 6.35's triangle has two apexes on the axis, and 6.36's lens comes
to a point at each end. The silhouette is therefore one closed filled path, the region's profile and
its reflection below the axis, and the axis of revolution is the x-axis, drawn solid outside the
solid and dashed inside it. Both panels keep equal scales, because 6.36's solid is 4 wide and 8
tall, which a filled panel would redraw as a squat lens.

**The bicone is the one solid here with a circle to draw.** Its surface turns at `x = 1`, so that
cross-section's circle is drawn: a vertical circle of radius 1 projects to an ellipse with
`rx = (2/7) s` and `ry = s`, its near half is the left half on screen and stays ink, and its far half
is inside the solid and goes dashed in annotation grey. 6.36's lens has no such turn, so it has no
cross-section to draw.

**6.35 (a) needed a taller frame below the axis than above.** The two lines cross at (1, 1) and their
labels belong at the lines' ends, but a frame centred on the region put the lower label on the tick
and the axis name. The frame runs `y` from -1.2 to 2, which keeps equal scales and leaves room below
the axis for the second label, and the `y` axis name moves clear of the upper line's end at the top
of the frame.

The two captions held one row each: the report went from 225 to 223, the certain tier from 121 to
119, distinct figures from 51 to 49 and redraws from 2 to 0, which empties the tier. The lesson
keeps one other orphan caption, so the lesson count stays at 100, and the ceiling in
`internal/lessons/assets_test.go` follows, to 223. The lesson's e2e case count moves to eleven.

**Ten Prealgebra upstream drawings** (Figures 1.3, 1.11, 1.13, 1.14, 1.15 and 2.6 to 2.10) landed
on 2026-10-16 in twenty-one arithmetic lessons, as
`web/next-app/public/diagrams/prealgebra/<number>.png`. They are the first figures in the tree that
are fetched rather than drawn: their disposition is `illustration, fetch`, so upstream's own drawing
is the look-alike, and the credit line is already the lesson's header.

**Fetched figures live under `diagrams/`, beside the drawn ones.** That is where the tree already
keeps a lesson figure — the reclaimed OpenStax SVGs are in `diagrams/calculus/` and
`diagrams/arithmetic/` — and it is what puts them under the e2e spec's `img[src*="/diagrams/"]`
selector and its 404 watch. A separate `public/figures/` directory would have rendered fine and been
measured by nothing. The diagram audit's metadata rules iterate only the assets the engine's
`conceptDiagrams` map names, so a fetched lesson figure needs no `meta.json` entry.

**A dollar sign in an alt is escaped.** Upstream's alt for Figure 1.3 describes three stacks of
currency and names `$100`, `$10` and `$1` bills. The app's math pass runs on the raw markdown before
parsing, so an unescaped `$100 bills ... $300` pairs two prices into an inline math span and eats
the rest of the alt; `\$` is the corpus's own spelling for a literal dollar and survives the pass.

**Both alt cases appear in this batch.** Figures 1.3, 1.11, 1.13, 1.14 and 1.15 have a bare label
for a caption, so their alt is upstream's `data-alt` verbatim; Figures 2.6 to 2.10 have a caption
that is a sentence, so their alt is that sentence with the label stripped. Every caption here is
plain text, so `plain` had nothing to render.

The ten captions held 47 rows between them: the report went from 223 to 176, the certain tier from
119 to 72, distinct figures from 49 to 39 and the fetch tier from 44 to 34. The audit's shipped and
referenced counts both move by ten, to 653 and 347, because the figures are files under
`public/diagrams/` that lesson bodies now reference. The lesson count drops from 100 to 79, since
twenty-one lessons no longer carry an orphan. The ceiling in `internal/lessons/assets_test.go`
follows, to 176, and the e2e case counts move to eight and two.

**Nine more Prealgebra upstream drawings** (Figures 4.4, 4.6, 4.7, 5.2, 5.3, 5.4, 6.2, 6.6 and 6.7)
landed on 2026-10-16 in twenty-five lessons, as
`web/next-app/public/diagrams/prealgebra/<number>.png`.

**The disposition follows the caption, not the picture.** Figure 6.6 is a receipt and a fan of
five-dollar bills rather than a mathematical object, and 5.2 is a chart of column headings, but both
are `illustration, fetch`: no caption in this batch names a photographer or an agency, which is the
only signal that separates the photo tier, and none of the nine is a number line, a graph or a
region that a redraw would have had to determine.

**Figure 4.7's caption is the reason the caption in the lesson is the source of truth.** Upstream's
copy of it loses the fraction and the escaped dollars that the lesson's repaired text carries, and
upstream's copy of Figure 5.4 is mojibake, so an alt read off upstream would have disagreed with the
caption printed above the image. Every lesson sharing a figure carries the same caption, and the
fetch script asserts that rather than trusting the first one it reads.

The nine captions held 31 rows between them: the report went from 176 to 145, the certain tier from
72 to 41, distinct figures from 39 to 30 and the fetch tier from 34 to 25. The audit's shipped and
referenced counts both move by nine, to 662 and 356. The lesson count drops from 79 to 59, since
twenty of the twenty-five lessons no longer carry an orphan and the five percent lessons still
carry another. The ceiling in `internal/lessons/assets_test.go` follows, to 145, and the e2e case
counts move to three and three, because Figure 4.4 landed in `frac.basics.parts` and Figure 6.2 in
`pct.convert.to_dec`, which are two of the spec's cases.

Section 2 records that a fetched figure is put on the dialect's greys, and these nine are converted
as they land.

**The conversion is applied at fetch time, because the renderer applies nothing.** A lesson figure
reaches the page as a plain `<img>` with no filter (`KatexContent.tsx:208`), so a tone that is not
baked into the asset does not exist. Each fetched asset is flattened to `sRGB`, dropped to grey,
stretched with `-auto-level`, and mapped through a 256x1 lookup table of the dialect's five levels —
`#333333`, `#636363`, `#BABABA`, `#E6E6E6` and white, which is a dialect color too.

**The stretch is not decoration.** Five of the nineteen fetched assets bottom out at a mid grey: the
micrograph at 1.11 at 0.48, the quarter at 4.6 at 0.43, the three cell-counting sets at 0.27 and the
hundredths grid at 6.2 at 0.22. Without the stretch they map entirely into the light half of the
ramp, so they carry no ink at all and read as washed out rather than as line art. The other fourteen
bottom out at 0.14, which is close enough to black that the stretch barely moves them.

**The paper stays white.** Mapping white to `#E6E6E6` instead would put every fetched figure on a
grey card that nothing else on the page shares.

The ten drawings that landed before this batch are still in OpenStax's color, and the commit that
follows converts them; that is also where the format the CDN actually serves is recorded.

**The ten drawings fetched before this batch are now real PNGs on the same greys.** They were fetched
as `.png` and sixteen of the nineteen were not PNGs at all: byte-identical to OpenStax's raw archive
resource, which is a CMYK JPEG with an Adobe APP14 marker, named `.png`. The image CDN serves that
original while its `f=png` transform is cold and answers `content-type: image/jpeg`, and the guard in
the fetch script tested the header's prefix, `image/`, which `image/jpeg` satisfies — so the lie was
stored. The guard now tests the PNG signature, the download retries six times, and it raises rather
than storing a fallback. All nineteen re-downloaded through `f=png` as lossless palette PNGs: 1896 KB
of CMYK became 387 KB of PNG, and the sixteen that were JPEGs went 1819 KB to 311 KB.

**One hash does not transform at all**, which is why the guard raises instead of falling back
silently: Figure 1.2's number line, `5ac935f7de7aea3f4f9f847d30f647fca6a89942`, still answered
`image/jpeg` after six retries, because the CDN 302s to its image API and that 302s again to the
archive resource. It is not in the fetched set yet, so nothing here depends on it, and a later batch
that meets it needs a deliberate local CMYK-to-sRGB conversion with a note saying so.

**A review montage is not evidence about a CMYK file.** ImageMagick's `montage` renders an Adobe CMYK
JPEG as black with red rules while a plain `magick in.png out.png` converts it correctly, and it
renders 16-bit grayscale PNGs as blank white. Montages are therefore built from 8-bit sRGB copies,
and Chromium — which drew the CMYK originals correctly, and is why the first batch looked right in
the app the whole time — stays the authority on what a figure looks like.

**Twenty-five Calculus Volume 1 illustrations** (2.20; 3.18, 3.19, 3.20, 3.21; 3.30, 3.31, 3.32;
3.34, 3.35; 4.2, 4.5, 4.6; 4.41, 4.43, 4.44, 4.46, 4.47, 4.48, 4.49, 4.55, 4.57, 4.60, 4.61; and
6.34) landed on 2026-10-16 in nine lessons, as
`web/next-app/public/diagrams/calculus/<number>.png`. That empties the fetch tier, which section 4
records, and leaves only the five photo captions.

**Two captions carry emphasis, so the caption paragraph is matched by its whole line.** Figure 2.20
ends "has infinite limits at *a*" and Figure 3.19 reads "horizontal tangent lines at *x* = 2/3 and
*x* = 4", so the Prealgebra rule — a run of characters containing no asterisk — matches neither. The
emphasis markers stay in the alt, which is what the four redrawn figures whose captions carry them
already do.

**The `\|` this batch was told to add to `plain` does not exist.** Section 6 renders Figure 4.48's
caption as `$\left\| {f(x) - L} \right\| < \varepsilon$`, and that `\|` is the table escaping its own
column separator: the caption in the lesson is `$\left| {f(x) - L} \right| < \varepsilon$`, a single
bar, which `\left` and `\right` already skip. What the caption did need is `\varepsilon`, and Figure
4.5's needs `\theta`; the assertion on an unknown macro is what found both.

**The CDN generates the `f=png` transform on demand, and twenty-five cold assets need more patience
than one asset's.** `fetch_png`'s six retries were written for the Prealgebra repair, where an asset
warmed after a request or two. Asking for Figure 2.20 six times spends the whole wait on it while the
other twenty-four transforms have not been requested at all, so the batch fetcher asks for every
asset first and then re-asks every asset still cold each round. That converges — 21 cold after round
one, 16 by round six, 3 by round ten — but it took two runs of the script to land all twenty-five,
and the transform was still cold for all twenty-five at the first request.

**A batch is fetched before any lesson is written.** The fetcher runs after every alt has been read
and validated and before the first lesson is edited, so a fetch that fails leaves the corpus
untouched rather than half-wired, which is what let the first attempt fail on Figure 2.20 with no
lesson changed.

The twenty-five captions held 33 rows across nine lessons: the report went from 145 to 112, the
certain tier from 41 to 8, distinct figures from 30 to 5, and the fetch tier from 25 to 0. The
audit's shipped count moves by twenty-five, to 687, and the referenced count by the same twenty-five
distinct URLs, to 381, because 33 rows reference 25 files. The lesson count drops from 59 to 52:
seven of the nine lessons no longer carry an orphan, and `calc.limit.numeric` and
`calc.deriv.related_rates` still carry the photo captions 2.22 and 4.4. The ceiling in
`internal/lessons/assets_test.go` follows, to 112, and the e2e case counts to eleven and twelve,
because Figure 2.20 landed in `calc.limit.numeric` and Figure 6.34 in `calc.integral.volume`.

**Three Calculus photographs** (2.22, `credit: NASA`; 4.4, Steve Jurvetson; 5.34, Carter Brown)
landed on 2026-10-16 in three lessons, as `web/next-app/public/diagrams/calculus/<number>.png`,
duotoned like every other fetched asset.

**A caption that is only a credit is not a description, so the alt is upstream's.** None of the
three says anything about its picture — `Figure 2.22 (credit: NASA)` — and the fourth one, 6.4, is
the same. So the alt rule has three branches rather than two, and this is the third: a bare label, a
sentence, and a credit. Upstream's own `data-alt` is then the only description that exists ("A
picture of a futuristic spaceship speeding through deep space"), and `plain` is not run on it,
because it is the author's words rather than the caption's mathematics.

**The credit stays in the caption, and nothing is added to the lesson.** The credit line is
upstream's own attribution and it travels with the figure, inside the caption a reader already sees.
Where a global third-party notice belongs is not decided here; section 8's item 3 keeps the per
figure question open, and this is the decision it was waiting for on the fetch side.

The three captions held 3 rows across three lessons: the report went from 112 to 109, the certain
tier from 8 to 5, and the audit's shipped and referenced counts from 687 to 690 and 381 to 384. All
three lessons lose their last orphan, so the lesson count drops from 52 to 49. The ceiling in
`internal/lessons/assets_test.go` follows, to 109, and `calc.limit.numeric`'s e2e count to twelve,
because Figure 2.22 is the only fetched photograph in the corpus so far.

**Two Prealgebra photographs** (6.4, `credit: Riles32807`; 6.8, `credit: Charleston's TheDigitel`)
landed on 2026-10-16 in five lessons, as `web/next-app/public/diagrams/prealgebra/<number>.png`,
duotoned like every other fetched asset. Figure 6.8's caption is the tier's second shape — it
describes its picture ("Applications of discounts are common in everyday life") and carries a credit
after it — so its alt is the caption's own words, where 6.4's is upstream's `data-alt`, because its
caption is a credit and nothing else.

**That empties the `certain` tier, and with it sections 3 and 6.** The audit's report named no
certain tier at all through that day: the 104 rows that remained were 33 `strong` and 71 `likely`, and
nothing else. The table in section 6 was empty that day, because a disposition was only ever
assigned to a certain caption, and section 3 recorded a tier that was gone rather than a count of
work left. What remained was a different kind of task — a `strong` or `likely` candidate is matched
to an asset that already exists or removed — and it is measured rather than planned here.

The two captions held 5 rows across five lessons: the report went from 109 to 104, the certain tier
from 5 to 0, and the audit's shipped and referenced counts from 690 to 692 and 384 to 386. All five
lessons lose their last orphan, so the lesson count drops from 49 to 44, and the ceiling in
`internal/lessons/assets_test.go` follows, to 104. `pct.convert.to_dec`'s e2e count moves to four.

**The ceiling's first bucket is empty too.** The comment above `TestOrphanCaptionsDoNotGrow` used to
give three reasons for keeping an orphan, and the first was that the lesson's prose cites the figure,
by number or by position. No row was kept for that reason in the 104, because every caption a lesson
cites by number is a caption with an auto-numbered label and those had all landed, so the 104 split
as 74 (no picture vocabulary) and 30 (bare section labels). The bucket filled again on 2026-10-17
with the 49 rows section 3 records, and the comment now carries them as its first bullet.

**The Precalculus 2e trigonometry re-ingest, 2026-10-17.** `trig.period` and `trig.radians` are the
only two lessons the retired hand-rolled converter produced, and it dropped their figure captions
along with their images; re-ingesting both from the curated source restores both halves, so the
captions came back with the text. The regenerated audit moved 104 rows to 153 and 0 in the certain
tier to 49 — 25 from §6.1 and 24 from §5.1 — and the ceiling in `internal/lessons/assets_test.go`
follows, to 153. That is the one movement of that ceiling which is up rather than down, and its
reason is a content decision rather than an ingestion accident: the record in section 5 says why the
sections were replaced rather than repaired, and section 3 says why their captions are work and not a
judgement. The row count also rises because the two lessons are longer than the mangled pair, so
`docs/lesson-media.md`'s per-lesson line numbers and its lesson count move with it. Section 6 assigns
each of the forty-nine its disposition: forty-seven are `schematic, redraw`, the water wheel is
`illustration, fetch`, and the chapter's rainbow photograph is the one `photo, third-party credit`,
which waits on where a credit line renders rather than on the fetch.

## 8. Order of work

**Every item below was done as of 2026-10-16 and is owed again since 2026-10-17**, and the list stays
because it is the order the 49 rows section 3 records will follow, in the same sequence the four
earlier batches took. It is three tiers now rather than four: the `unidentified` tier used to lead
here as its own step, and section 4 records that it was emptied on 2026-10-16, so those seventeen
figures sit inside items 1 and 2 by their new dispositions.

1. **The schematics.** Highest value and no license question: they are redrawn as Mathua's own
   work, in the dialect the style doc fixes, and they land under `source: local`. This is forty-seven
   of the forty-nine.
2. **The illustrations.** Fetch, put on the dialect's greys (section 2) and credit. This is the §5.1
   water wheel.
3. **The photos.** Decide per figure between a licensed fetch and a replacement, because a
   photograph cannot be redrawn as a line figure. All five of the last batch were fetched, which
   section 7 records, so a licensed fetch was both possible and faithful every time; the credit line
   the caption already carries is what made the decision easy. The one waiting now, §6.1's rainbow,
   is blocked on where that line renders rather than on the fetch.

Each figure is wired by the paths in section 5 of `docs/diagram-style.md`, and the published
`web/next-app/public/diagrams/meta.json` copy is refreshed before the change is called done.

## 9. Method

```bash
python3 scripts/audit_lesson_media.py --json data/lessons/figure-audit.json > docs/lesson-media.md
```
