# Wave 10 candidate ledger

Status: approved and authoring complete on 2026-11-10. The first ten passed review and the remaining forty were authorized. All fifty concepts are authored.

Wave 10 raises the corpus from 1,950 to 2,000 concepts. It is the final wave and it closes every remaining domain target exactly.

Every id and normalized label below was checked against every file in `data/concepts/` on 2026-11-10. There is no exact id collision and no normalized-label collision. Every prerequisite points to an existing concept.

## Allocation

| Domain | Tier | Now | Target | Add | Course |
|---|---|---:|---:|---:|---|
| prealgebra | school | 71 | 89 | +18 | `pa` and grades 6 to 8 |
| geometry | school | 70 | 87 | +17 | `geo` |
| graph_theory | university | 63 | 70 | +7 | `graph` |
| machine_learning | university | 67 | 75 | +8 | `ml` |
| **Total** | | | | **+50** | |

Prealgebra and geometry are learner domains, so their concepts need authored solution schemas. Graph theory and machine learning are not learner domains and need none.

## Prealgebra: 18 candidates

| id | Label | Subdomain | Prerequisite | Families |
|---|---|---|---|---|
| prealg.ratio.three_term | Three-term ratios | prealgebra.ratios | ratio.basics.concept | Combine three quantities into one ratio; simplify a three-term ratio; scale a three-term ratio to a total |
| prealg.ratio.part_whole | Part-to-whole ratios | prealgebra.ratios | ratio.basics.concept | Write a part-to-whole ratio from counts; find the whole from a part and its ratio; compare part-to-part with part-to-whole |
| prealg.ratio.tape | Tape diagrams for ratios | prealgebra.ratios | ratio.basics.concept | Read a tape diagram; find the total from equal units; recover a missing unit count |
| prealg.ratio.constant | The constant of proportionality | prealgebra.ratios | ratio.ops.proportion | Find the constant from a table; use it to predict a value; identify a proportional table |
| pct.ops.original | Finding the original amount | prealgebra.percentages | pct.ops.change | Recover the original after an increase; recover it after a decrease; check the result |
| pct.ops.percent_points | Percentage points | prealgebra.percentages | pct.basics.concept | Compute a change in percentage points; convert points to percent change; distinguish the two |
| pct.ops.successive | Successive percentage changes | prealgebra.percentages | pct.ops.change | Apply two increases; apply an increase then a decrease; find the net multiplier |
| pct.ops.reverse | Reverse percentages | prealgebra.percentages | pct.ops.of_number | Find a number from a percentage of it; recover the base; check the percentage |
| dec.ops.div_decimal | Dividing by a decimal | prealgebra.decimals | dec.ops.div | Divide by a decimal; scale the divisor; check the quotient |
| dec.ops.estimate | Estimating with decimals | prealgebra.decimals | dec.basics.compare | Estimate a sum; estimate a product; compare with the exact value |
| dec.basics.expanded | Expanded form of decimals | prealgebra.decimals | dec.basics.place_value | Write expanded form; sum the place values; return to standard form |
| prealg.expr.two_variables | Evaluating expressions with two variables | prealgebra.variables | prealg.expr.evaluate | Evaluate at two values; substitute and simplify; recover a value from a table |
| prealg.expr.coefficient | Coefficients and constant terms | prealgebra.variables | prealg.expr.like_terms | Read the coefficient and constant; combine like terms; count terms |
| prealg.expr.distribute_simplify | Distributing and simplifying | prealgebra.variables | prealg.expr.distribute | Distribute; combine like terms; evaluate the result |
| prealg.eq.literal | Solving a formula for a variable | prealgebra.equations | prealg.eq.two_step, prealg.expr.formula | Solve for a named variable; substitute values; check the rearranged formula |
| prealg.eq.check | Checking a solution | prealgebra.equations | prealg.eq.one_step_add | Substitute into both sides; identify an incorrect solution; verify the check |
| prealg.measure.time | Converting units of time | prealgebra.measurement | prealg.measure.convert | Convert hours, minutes and seconds; find elapsed time; combine mixed units |
| prealg.measure.precision | Precision and significant digits | prealgebra.measurement | prealg.measure.metric | Count significant digits; round to a stated precision; compare two measurements |

## Geometry: 17 candidates

| id | Label | Subdomain | Prerequisite | Families |
|---|---|---|---|---|
| geo.triangle.altitude | Altitudes of a triangle | geometry.triangles | geo.triangle.area | Find an altitude from area; find area from an altitude; locate the altitude in an acute or obtuse triangle |
| geo.triangle.orthocenter | The orthocenter | geometry.triangles | geo.triangle.altitude | Find the orthocenter from altitudes; identify it in a right triangle; use a coordinate triangle |
| geo.triangle.circumcenter | The circumcenter | geometry.triangles | geo.triangle.perpendicular_bisector | Intersect perpendicular bisectors; find the circumradius; use the right-triangle midpoint |
| geo.triangle.incenter | The incenter | geometry.triangles | geo.triangle.angle_bisector | Intersect angle bisectors; find the inradius from area and semiperimeter; find a distance to a side |
| geo.triangle.cevian | Ceva's theorem | geometry.triangles | geo.triangle.median_centroid | Compute the Ceva product; solve for a missing ratio; test concurrency |
| geo.triangle.apollonius | Apollonius's theorem | geometry.triangles | geo.triangle.median_centroid | Find a median length; use the side relation; check a median value |
| geo.circle.power_point | Power of a point | geometry.circles | geo.circle.chord_segments | Multiply chord segments; find a tangent length; use a secant |
| geo.circle.sagitta | Sagitta and chord height | geometry.circles | geo.circle.chord | Find the radius from chord and sagitta; find the sagitta; use the central angle |
| geo.circle.radical_axis | The radical axis | geometry.circles | geo.circle.power_point | Find an equal-power line; test two circle equations; count intersections |
| geo.polygon.diagonal_count | Counting polygon diagonals | geometry.polygons | geo.polygon.interior_angles | Count diagonals; recover the side count; count all segments |
| geo.polygon.tessellation | Tessellations and angle sums | geometry.polygons | geo.polygon.interior_angles | Find the angle sum; test a regular tessellation; count polygons around a vertex |
| geo.coord.section | The section formula | geometry.coordinate | geo.coord.midpoint | Divide a segment in a ratio; find the midpoint; recover a ratio from coordinates |
| geo.coord.sphere_equation | Equations of spheres | geometry.coordinate | geo.coord.distance_3d | Read centre and radius; test a point; build an equation from a diameter |
| geo.transform.composition | Composition of transformations | geometry.transform | geo.transform.translation, geo.transform.rotation | Apply two transformations; invert a composition; find a fixed point |
| geo.transform.isometry | Isometries and their invariants | geometry.transform | geo.transform.reflection | Test distance preservation; test orientation; classify an isometry |
| geo.solid.hemisphere | Hemispheres | geometry.solid | geo.solid.sphere | Find volume; find surface area; compare with the full sphere |
| geo.solid.density | Density and volume | geometry.solid | geo.solid.volume_rect | Find density from mass and volume; find mass; find volume |

## Graph theory: 7 candidates

| id | Label | Subdomain | Prerequisite | Families |
|---|---|---|---|---|
| graph.feedback_vertex_set | Minimum feedback vertex sets | graph_theory.connectivity | graph.cut_vertex | Find the minimum size; test a proposed set; count remaining forest edges |
| graph.feedback_arc_set | Minimum feedback arc sets | graph_theory.connectivity | graph.strong_connectivity | Find the minimum size; test a proposed set; test whether the remainder is a DAG |
| graph.arboricity | Arboricity | graph_theory.planarity | graph.degeneracy | Compute the edge bound; find a forest decomposition count; test a density bound |
| graph.thickness | Graph thickness | graph_theory.planarity | graph.kuratowski | Test a planar subgraph decomposition; find a lower bound; count planar layers |
| graph.crossing_number | Crossing number | graph_theory.planarity | graph.genus | Read a small drawing; compute K5 or K3,3 crossing counts; test a claimed drawing |
| graph.edge_cover | Edge covers | graph_theory.matching | graph.vertex_cover | Find the minimum size; test a cover; relate it to a matching |
| graph.clique_cover | Clique covers | graph_theory.matching | graph.clique | Find the minimum clique cover; test a partition; use the complement colouring |

## Machine learning: 8 candidates

| id | Label | Subdomain | Prerequisite | Families |
|---|---|---|---|---|
| ml.optim.adagrad | AdaGrad optimization | machine_learning.neural_networks | ml.optim.sgd | Apply one update; track accumulated squared gradient; compare a second step |
| ml.optim.rmsprop | RMSProp optimization | machine_learning.neural_networks | ml.optim.adagrad | Apply one update; update the decayed average; compare with AdaGrad |
| ml.optim.adamw | AdamW optimization | machine_learning.neural_networks | ml.optim.adam | Apply the decoupled weight decay; compute the update; compare with Adam |
| ml.ensemble.stacking | Stacking ensembles | machine_learning.trees | ml.tree.ensemble | Build the meta-feature; make the meta-prediction; compare with a base model |
| ml.ensemble.voting | Voting classifiers | machine_learning.trees | ml.tree.ensemble | Count votes; choose the majority label; apply the tie rule |
| ml.unsupervised.spectral | Spectral clustering | machine_learning.unsupervised | ml.cluster.kmeans | Read the Fiedler sign; assign clusters; compute a normalized cut |
| ml.loss.hinge | Hinge loss | machine_learning.neural_networks | ml.svm.margin | Evaluate the hinge loss; compute the margin; read the subgradient |
| ml.loss.huber | Huber loss | machine_learning.neural_networks | ml.loss.mae, ml.loss.mse | Evaluate the loss below the threshold; evaluate above it; compare with MSE and MAE |

## Review order

The first ten cover every domain and the riskiest arithmetic and schematic paths.

- Prealgebra: three-term ratios, part-to-whole ratios, finding the original amount, solving a formula for a variable.
- Geometry: altitudes, power of a point, the section formula.
- Graph theory: minimum feedback vertex sets, edge covers.
- Machine learning: AdaGrad.

The first ten are authored and inspected before the remaining forty. If they expose a duplicate, a prerequisite error, a schema problem or a generator-shape problem, the process is corrected before the rest of the wave.

## Artifact contract

Every concept needs a numeric node in its domain file, a generator in the domain package registered from `Register`, an independent fact-based test, a lesson over 1 KB with one H1 and exactly three H2 sections, a three-entry KP shard, a lessons.json row and a course target. Prealgebra and geometry also need authored solution schemas. Graph theory and machine learning do not.

Each concept must produce three section-aligned question families, at least five distinct questions per ordinal in twelve draws, and at least ten per ordinal in sixty draws or a written reason. The design target is twenty per ordinal in 200 draws. Every answer is recomputed independently from Facts.

No prerequisite is removed. Baseline exemptions are not allowed. The additive baseline stays frozen at 1,500 concepts. Band width stays 3. Commits are local. Nothing is pushed or deployed.
