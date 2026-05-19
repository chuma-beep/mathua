#!/usr/bin/env python3
"""
Generate lessons.json mapping and add new concept IDs for scraped Algebrica content.
"""

import json
import os
import glob

BASE = os.path.join(os.path.dirname(__file__), "..")
LESSONS_DIR = os.path.join(BASE, "data", "lessons", "algebrica")
CONCEPTS_DIR = os.path.join(BASE, "data", "concepts")
LESSONS_JSON = os.path.join(BASE, "data", "lessons", "lessons.json")

# ─── Mapping: lesson file → concept IDs to map to ────────────────────────
# Format: "subdirectory/filename.md": ["concept.id.one", "concept.id.two"]

MAPPING = {
    # ── Algebraic Structures (already mapped, keep existing) ──
    "algebraic-structures/groups.md": ["abstract.group.def", "abstract.group.subgroup", "abstract.group.examples"],
    "algebraic-structures/rings.md": ["abstract.rings.def"],
    "algebraic-structures/vector-spaces.md": ["abstract.group.homomorphism"],
    "algebraic-structures/fields.md": [],
    "algebraic-structures/modules.md": [],

    # ── Complex Numbers ──
    "complex-numbers/complex-numbers.md": ["complex.basics.concept"],
    "complex-numbers/operations-with-complex-numbers.md": ["complex.ops.add_sub", "complex.ops.mult", "complex.ops.divide", "complex.ops.conjugate"],
    "complex-numbers/complex-numbers-in-trigonometric-form.md": ["complex.adv.polar"],
    "complex-numbers/de-moivre-theorem.md": ["complex.adv.de_moivre"],
    "complex-numbers/roots-of-unity.md": ["complex.adv.roots"],
    "complex-numbers/complex-numbers-in-exponential-form.md": ["complex.adv.exponential"],
    "complex-numbers/fundamental-inequalities-for-complex-numbers.md": ["complex.adv.inequalities"],

    # ── Trigonometry ──
    "trigonometry/sine-and-cosine.md": ["trig.basics.sin_cos_def", "trig.graph.sin", "trig.graph.cos"],
    "trigonometry/tangent-and-cotangent.md": ["trig.basics.tan_def"],
    "trigonometry/unit-circle.md": ["trig.basics.unit_circle"],
    "trigonometry/pythagorean-identity.md": ["trig.ident.pythagorean"],
    "trigonometry/secant-and-cosecant.md": ["trig.basics.reciprocal"],
    "trigonometry/reduction-formulas-and-reference-angles.md": ["trig.basics.special_angles", "trig.basics.reference_angle"],
    "trigonometry/arcsine-and-arccosine.md": ["trig.adv.inverse"],
    "trigonometry/arctangent-and-arccotangent.md": ["trig.adv.arctan"],
    "trigonometry/the-law-of-sines.md": ["trig.adv.law_sines"],
    "trigonometry/the-law-of-cosines.md": ["trig.adv.law_cosines"],
    "trigonometry/pythagorean-theorem.md": ["geo.triangle.pythagorean"],
    "trigonometry/hyperbolic-sine-and-cosine.md": ["trig.hyperbolic.sinh_cosh"],
    "trigonometry/hyperbolic-tangent-and-cotangent.md": ["trig.hyperbolic.tanh_coth"],
    "trigonometry/right-triangle-trigonometry.md": ["trig.basics.right_triangle"],
    "trigonometry/trigonometric-identities.md": ["trig.ident.identities"],

    # ── Vectors and Matrices ──
    "vectors-and-matrices/vectors.md": ["linalg.vector.concept", "linalg.vector.add", "linalg.vector.dot"],
    "vectors-and-matrices/matrices.md": ["linalg.matrix.concept", "linalg.matrix.add", "linalg.matrix.mult"],
    "vectors-and-matrices/eigenvalues-and-eigenvectors.md": ["linalg.eigen.concept", "linalg.eigen.compute"],
    "vectors-and-matrices/determinant-of-a-square-matrix.md": ["linalg.det.2x2", "linalg.det.3x3"],
    "vectors-and-matrices/linear-combinations.md": ["linalg.lintrans.concept", "linalg.vec.span", "linalg.vec.basis"],
    "vectors-and-matrices/inverse-matrix.md": ["linalg.matrix.identity"],
    "vectors-and-matrices/rank-of-a-matrix.md": ["linalg.matrix.rank"],
    "vectors-and-matrices/matrix-diagonalization.md": ["linalg.eigen.diagonalization"],

    # ── Polynomials ──
    "polynomials/polynomials.md": ["alg.poly.concept", "calc.deriv.power_rule"],
    "polynomials/adding-and-subtracting-polynomials.md": ["alg.poly.add_sub"],
    "polynomials/binomials.md": ["alg.poly.foil"],
    "polynomials/trinomials.md": ["alg.factor.trinomial"],
    "polynomials/monomials.md": ["alg.poly.monomial"],
    "polynomials/notable-products.md": ["alg.poly.special"],
    "polynomials/completing-the-square.md": ["alg.quad.complete_square"],
    "polynomials/factoring-polynomials-ac-method.md": ["alg.factor.ac_method"],
    "polynomials/binomial-theorem.md": ["discrete.combinatorics.binomial_theorem"],
    "polynomials/polynomial-division.md": ["alg.poly.division"],
    "polynomials/roots-of-a-polynomial.md": ["alg.poly.roots"],
    "polynomials/synthetic-division-method.md": ["alg.poly.synthetic_division"],
    "polynomials/vieta-formulas.md": ["alg.poly.vieta"],
    "polynomials/partial-fraction-decomposition.md": [],

    # ── Equations ──
    "equations/equations.md": ["alg.eq.multi_step"],
    "equations/linear-equations.md": ["alg.linear.slope_intercept"],
    "equations/linear-equations-with-parameters.md": ["alg.eq.literal"],
    "equations/factoring-quadratic-equations.md": ["alg.quad.solve_factor"],
    "equations/quadratic-formula.md": ["alg.quad.formula"],
    "equations/quadratic-equations.md": ["alg.quad.quadratic"],
    "equations/incomplete-quadratic-equations.md": ["alg.quad.incomplete"],
    "equations/polynomial-equations.md": ["alg.eq.polynomial"],
    "equations/absolute-value-equations.md": ["alg.eq.absolute_value"],
    "equations/rational-equations.md": ["alg.eq.rational"],
    "equations/irrational-equations.md": ["alg.eq.irrational"],
    "equations/exponential-equations.md": ["alg.eq.exponential"],
    "equations/logarithmic-equations.md": ["alg.eq.logarithmic"],
    "equations/binomial-equations.md": ["alg.eq.binomial"],
    "equations/trinomial-equations.md": ["alg.eq.trinomial"],
    "equations/trigonometric-equations.md": ["trig.eq.basic"],
    "equations/homogeneous-trigonometric-equations.md": ["trig.eq.homogeneous"],
    "equations/quadratic-equations-with-complex-solutions.md": ["alg.quad.complex"],
    "equations/quadratic-equations-with-parameters.md": ["alg.quad.parametric"],
    "equations/equations-with-parameters.md": ["alg.eq.literal"],
    "equations/loss-of-roots.md": ["alg.eq.extraneous_roots"],
    "equations/geometrical-meaning-quadratic-equations.md": ["alg.func.quad"],

    # ── Powers, Radicals and Logarithms ──
    "powers-radicals-logarithms/powers.md": ["arith.exp.concept", "arith.exp.product_rule", "arith.exp.quotient_rule", "arith.exp.power_rule", "arith.exp.evaluate"],
    "powers-radicals-logarithms/radicals.md": ["arith.sqrt.simplify"],
    "powers-radicals-logarithms/logarithms.md": ["alg.log.concept", "alg.log.evaluate", "alg.log.properties"],

    # ── Sets and Numbers ──
    "sets-and-numbers/sets.md": ["discrete.sets.operations"],
    "sets-and-numbers/factorial.md": ["discrete.combinatorics.permutations"],
    "sets-and-numbers/binomial-coefficient.md": ["discrete.combinatorics.pascal", "discrete.combinatorics.combinations"],
    "sets-and-numbers/natural-numbers.md": ["nt.adv.euler_phi"],
    "sets-and-numbers/integers.md": ["nt.basics.divisibility"],
    "sets-and-numbers/modulo-operator.md": ["nt.basics.modular", "nt.basics.congruence"],
    "sets-and-numbers/absolute-value.md": ["arith.neg.abs_value"],
    "sets-and-numbers/intervals.md": ["alg.ineq.interval"],
    "sets-and-numbers/properties-of-real-numbers.md": ["prealg.real.properties"],
    "sets-and-numbers/real-numbers.md": ["prealg.real.concept"],
    "sets-and-numbers/types-of-numbers.md": ["prealg.types"],
    "sets-and-numbers/supremum-and-infimum.md": ["calc.limit.supremum"],

    # ── Limits ──
    "limits/remarkable-limits.md": ["calc.limit.concept"],
    "limits/limits.md": ["calc.limit.properties"],
    "limits/algebra-of-limits.md": ["calc.limit.algebra"],
    "limits/squeeze-theorem.md": ["calc.limit.squeeze"],
    "limits/asymptotes.md": ["calc.limit.asymptotes"],
    "limits/indeterminate-forms.md": ["calc.limit.indeterminate"],
    "limits/little-o-notation.md": ["calc.limit.little_o"],
    "limits/big-o-notation.md": ["calc.limit.big_o"],

    # ── Derivatives ──
    "derivatives/derivatives.md": ["calc.deriv.concept"],
    "derivatives/difference-quotient.md": ["calc.deriv.difference_quotient"],
    "derivatives/differential-of-a-function.md": ["calc.deriv.differential"],
    "derivatives/the-derivative-of-a-composite-function.md": ["calc.deriv.chain_rule"],
    "derivatives/derivative-of-composite-power-functions.md": ["calc.deriv.chain_rule"],
    "derivatives/maximum-minimum-and-inflection-points.md": ["calc.deriv.optimization"],
    "derivatives/partial-derivatives.md": ["calc.deriv.partial"],
    "derivatives/points-of-non-differentiability.md": ["calc.deriv.non_differentiability"],

    # ── Differential Calculus Theorems ──
    "differential-calculus-theorems/rolles-theorem.md": ["calc.deriv.rolle"],
    "differential-calculus-theorems/lagrange-theorem.md": ["calc.deriv.mvt"],
    "differential-calculus-theorems/cauchy-theorem.md": ["calc.deriv.cauchy_mvt"],
    "differential-calculus-theorems/fermat-theorem.md": ["calc.deriv.fermat"],
    "differential-calculus-theorems/hopital-rule.md": ["calc.limit.lhopital"],
    "differential-calculus-theorems/weierstrass-theorem.md": ["calc.limit.weierstrass"],

    # ── Integrals ──
    "integrals/indefinite-integrals.md": ["calc.integral.indefinite"],
    "integrals/integration-by-substitution.md": ["calc.integral.substitution"],
    "integrals/definite-integrals.md": ["calc.integral.definite"],
    "integrals/integration-by-parts.md": ["calc.integral.parts"],
    "integrals/finding-areas-by-integration.md": ["calc.integral.area_between"],
    "integrals/fundamental-theorem-of-calculus.md": ["calc.integral.ftc"],
    "integrals/integral-of-rational-functions.md": ["calc.integral.partial_fractions"],
    "integrals/improper-integrals.md": ["calc.integral.improper"],
    "integrals/integral-of-trigonometric-functions.md": ["calc.integral.trig_integrals"],
    "integrals/trigonometric-substitution-for-integrals.md": ["calc.integral.trig_substitution"],
    "integrals/the-weierstrass-substitution.md": ["calc.integral.weierstrass_sub"],
    "integrals/numerical-integration.md": ["calc.integral.numerical"],
    "integrals/arc-length-of-a-curve.md": ["calc.integral.arc_length"],
    "integrals/integral-of-the-exponential-function.md": ["calc.integral.exp_integral"],
    "integrals/riemann-integrability-criteria.md": ["calc.integral.riemann_criteria"],

    # ── Differential Equations ──
    "differential-equations/differential-equations.md": ["ode.basics.concept"],

    # ── Linear Systems ──
    "linear-systems/rouche-capelli-theorem.md": ["alg.systems.elimination", "alg.systems.substitution"],
    "linear-systems/systems-of-linear-equations.md": ["alg.systems.concept"],
    "linear-systems/cramers-rule.md": ["linalg.sys.cramer"],
    "linear-systems/solving-linear-systems-using-gaussian-elimination.md": ["alg.systems.gaussian"],

    # ── Functions ──
    "functions/functions.md": ["alg.func.concept"],
    "functions/determining-the-domain-of-a-function.md": ["alg.func.domain"],
    "functions/even-and-odd-functions.md": ["alg.func.even_odd"],
    "functions/increasing-and-decreasing-functions.md": ["alg.func.monotonicity"],
    "functions/composite-functions.md": ["alg.func.composite"],
    "functions/inverse-function.md": ["alg.func.inverse"],
    "functions/absolute-value-function.md": ["alg.func.absolute_value"],
    "functions/rational-functions.md": ["alg.func.rational"],
    "functions/polynomial-function.md": ["alg.poly.concept"],
    "functions/exponential-function.md": ["alg.exp.concept"],
    "functions/logarithmic-function.md": ["alg.log.concept"],
    "functions/sine-function.md": [],
    "functions/cosine-function.md": [],
    "functions/tangent-function.md": [],
    "functions/cotangent-function.md": [],
    "functions/secant-function.md": [],
    "functions/cosecant-function.md": [],
    "functions/analyzing-the-graphs-of-functions.md": ["alg.func.graph_analysis"],
    "functions/continuous-functions.md": ["calc.limit.continuity"],
    "functions/uniform-continuity.md": ["calc.limit.uniform_continuity"],
    "functions/discontinuities-of-real-functions.md": ["calc.limit.discontinuity"],
    "functions/convexity-and-concavity-of-functions.md": ["calc.deriv.convexity"],
    "functions/sign-function.md": ["alg.func.sign"],
    "functions/sigmoid-function.md": ["alg.func.sigmoid"],
    "functions/dirichlet-function.md": ["alg.func.dirichlet"],

    # ── Inequalities ──
    "inequalities/linear-inequalities.md": ["alg.ineq.multi_step"],
    "inequalities/quadratic-inequalities.md": ["alg.ineq.quadratic"],
    "inequalities/rational-inequalities.md": ["alg.ineq.rational"],
    "inequalities/irrational-inequalities.md": ["alg.ineq.irrational"],
    "inequalities/logarithmic-inequalities.md": ["alg.ineq.logarithmic"],
    "inequalities/inequalities-with-absolute-value.md": ["alg.ineq.absolute_value"],
    "inequalities/sign-analysis-in-inequalities.md": ["alg.ineq.sign_analysis"],
    "inequalities/systems-of-inequalities.md": ["alg.ineq.systems"],
    "inequalities/trigonometric-inequalities.md": ["trig.ineq.basic"],

    # ── Lines, Planes and Conic Sections ──
    "lines-planes-conic-sections/lines.md": ["geo.coord.lines"],
    "lines-planes-conic-sections/vector-and-parametric-equations-of-a-line.md": ["linalg.vector.parametric"],
    "lines-planes-conic-sections/polar-coordinates.md": ["geo.coord.polar"],
    "lines-planes-conic-sections/parabola.md": ["alg.conic.parabola"],
    "lines-planes-conic-sections/circumference.md": ["alg.conic.circle"],
    "lines-planes-conic-sections/ellipse.md": ["alg.conic.ellipse"],
    "lines-planes-conic-sections/hyperbola.md": ["alg.conic.hyperbola"],

    # ── Probability and Statistics ──
    "probability-and-statistics/introduction-to-the-mean.md": ["stat.data.mean"],
    "probability-and-statistics/arithmetic-mean.md": ["stat.data.mean"],
    "probability-and-statistics/median.md": ["stat.data.median"],
    "probability-and-statistics/geometric-mean.md": ["stat.adv.geometric_mean"],
    "probability-and-statistics/harmonic-mean.md": ["stat.adv.harmonic_mean"],
    "probability-and-statistics/root-mean-square.md": ["stat.adv.rms"],
    "probability-and-statistics/variance.md": ["stat.adv.variance"],
    "probability-and-statistics/variance-and-covariance-of-a-random-variable.md": ["stat.adv.covariance"],
    "probability-and-statistics/discrete-random-variables.md": ["stat.prob.discrete_rv"],
    "probability-and-statistics/continuous-random-variables.md": ["stat.prob.continuous_rv"],
    "probability-and-statistics/mean-or-expected-value-of-a-random-variable.md": ["stat.prob.expected_value"],
    "probability-and-statistics/bernoulli-distribution.md": ["stat.dist.bernoulli"],
    "probability-and-statistics/binomial-distribution.md": ["stat.dist.binomial"],
    "probability-and-statistics/geometric-distribution.md": ["stat.dist.geometric"],
    "probability-and-statistics/hypergeometric-distribution.md": ["stat.dist.hypergeometric"],
    "probability-and-statistics/poisson-distribution.md": ["stat.dist.poisson"],
    "probability-and-statistics/uniform-distribution.md": ["stat.dist.uniform"],
    "probability-and-statistics/normal-distribution.md": ["stat.dist.normal"],
    "probability-and-statistics/standard-normal-z-table.md": ["stat.dist.z_table"],
    "probability-and-statistics/exponential-distribution.md": ["stat.dist.exponential"],
    "probability-and-statistics/gamma-distribution.md": ["stat.dist.gamma"],
    "probability-and-statistics/beta-distribution.md": ["stat.dist.beta"],
    "probability-and-statistics/chi-square-distribution.md": ["stat.dist.chi_square"],
    "probability-and-statistics/student-t-distribution.md": ["stat.dist.student_t"],
    "probability-and-statistics/sampling-distributions.md": ["stat.infer.sampling"],
    "probability-and-statistics/confidence-intervals.md": ["stat.infer.confidence"],
    "probability-and-statistics/bayes-theorem.md": ["stat.prob.bayes"],

    # ── Sequences ──
    "sequences/sequences.md": ["calc.seq.concept"],
    "sequences/arithmetic-sequence.md": ["alg.seq.arithmetic"],
    "sequences/geometric-sequence.md": ["alg.seq.geometric"],
    "sequences/convergent-and-divergent-sequences.md": ["calc.seq.convergence"],
    "sequences/monotone-sequences.md": ["calc.seq.monotone"],
    "sequences/cauchy-sequence.md": ["calc.seq.cauchy"],
    "sequences/principle-of-mathematical-induction.md": ["discrete.proof.induction"],
    "sequences/euler-number-limit-sequence.md": ["calc.seq.euler"],
    "sequences/sequences-of-functions.md": ["calc.seq.function_sequences"],

    # ── Series ──
    "series/series.md": ["calc.series.concept"],
    "series/geometric-series.md": ["alg.seq.sum_geo"],
    "series/harmonic-series.md": ["calc.series.harmonic"],
    "series/series-with-positive-terms.md": ["calc.series.positive_terms"],
    "series/cauchy-convergence-criterion-series.md": ["calc.series.cauchy_criterion"],
    "series/integral-test-for-series-convergence.md": ["calc.series.integral_test"],
    "series/root-test-for-series-convergence.md": ["calc.series.root_test"],
    "series/leibniz-criterion.md": ["calc.series.alternating"],
    "series/function-series.md": ["calc.series.function_series"],
    "series/power-series.md": ["calc.series.power"],
    "series/taylor-series.md": ["calc.series.taylor"],
    "series/fourier-series.md": ["calc.series.fourier"],

    # ── Kinematics ──
    "kinematics/velocity.md": ["calc.deriv.applications"],
    "kinematics/acceleration.md": ["calc.deriv.applications"],
    "kinematics/simple-harmonic-motion.md": ["calc.deriv.applications"],

    # ── Various ──
    "various/propositional-logic.md": ["discrete.logic.propositions"],
    "various/cosine-similarity.md": ["linalg.vector.cosine_similarity"],
    "various/backpropagation.md": ["ml.backpropagation"],
}

# ─── New concept definitions to add ──────────────────────────────────────

NEW_CONCEPTS = [
    # complex_numbers
    {"id": "complex.adv.exponential", "label": "Exponential form of complex numbers (Euler's formula)", "domain": "complex_numbers", "subdomain": "complex_numbers.advanced", "grading_type": "symbolic", "prerequisites": ["complex.adv.polar"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "complex.adv.inequalities", "label": "Fundamental inequalities for complex numbers", "domain": "complex_numbers", "subdomain": "complex_numbers.advanced", "grading_type": "symbolic", "prerequisites": ["complex.basics.concept"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},

    # trigonometry
    {"id": "trig.adv.arctan", "label": "Arctangent and arccotangent functions", "domain": "trigonometry", "subdomain": "trigonometry.advanced", "grading_type": "symbolic", "prerequisites": ["trig.adv.inverse"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "trig.hyperbolic.sinh_cosh", "label": "Hyperbolic sine and cosine functions", "domain": "trigonometry", "subdomain": "trigonometry.hyperbolic", "grading_type": "symbolic", "prerequisites": ["trig.basics.sin_cos_def"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "trig.hyperbolic.tanh_coth", "label": "Hyperbolic tangent and cotangent functions", "domain": "trigonometry", "subdomain": "trigonometry.hyperbolic", "grading_type": "symbolic", "prerequisites": ["trig.hyperbolic.sinh_cosh"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "trig.basics.right_triangle", "label": "Right triangle trigonometry", "domain": "trigonometry", "subdomain": "trigonometry.basics", "grading_type": "symbolic", "prerequisites": ["trig.basics.sin_cos_def"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "trig.ident.identities", "label": "Trigonometric identities", "domain": "trigonometry", "subdomain": "trigonometry.identities", "grading_type": "symbolic", "prerequisites": ["trig.ident.pythagorean"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "trig.eq.basic", "label": "Basic trigonometric equations", "domain": "trigonometry", "subdomain": "trigonometry.identities", "grading_type": "symbolic", "prerequisites": ["trig.basics.sin_cos_def"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "trig.eq.homogeneous", "label": "Homogeneous trigonometric equations", "domain": "trigonometry", "subdomain": "trigonometry.identities", "grading_type": "symbolic", "prerequisites": ["trig.eq.basic"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "trig.ineq.basic", "label": "Trigonometric inequalities", "domain": "trigonometry", "subdomain": "trigonometry.identities", "grading_type": "symbolic", "prerequisites": ["trig.eq.basic"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},

    # calculus - limits
    {"id": "calc.limit.algebra", "label": "Algebra of limits", "domain": "calculus", "subdomain": "calculus.limits", "grading_type": "symbolic", "prerequisites": ["calc.limit.properties"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.limit.squeeze", "label": "Squeeze theorem", "domain": "calculus", "subdomain": "calculus.limits", "grading_type": "symbolic", "prerequisites": ["calc.limit.properties"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.limit.asymptotes", "label": "Asymptotes of functions", "domain": "calculus", "subdomain": "calculus.limits", "grading_type": "symbolic", "prerequisites": ["calc.limit.infinity"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.limit.indeterminate", "label": "Indeterminate forms of limits", "domain": "calculus", "subdomain": "calculus.limits", "grading_type": "symbolic", "prerequisites": ["calc.limit.properties"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.limit.little_o", "label": "Little-o notation", "domain": "calculus", "subdomain": "calculus.limits", "grading_type": "symbolic", "prerequisites": ["calc.limit.properties"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.limit.big_o", "label": "Big O notation", "domain": "calculus", "subdomain": "calculus.limits", "grading_type": "symbolic", "prerequisites": ["calc.limit.properties"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.limit.supremum", "label": "Supremum and infimum", "domain": "calculus", "subdomain": "calculus.limits", "grading_type": "symbolic", "prerequisites": ["calc.limit.concept"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.limit.lhopital", "label": "L'Hopital's rule", "domain": "calculus", "subdomain": "calculus.limits", "grading_type": "symbolic", "prerequisites": ["calc.limit.indeterminate", "calc.deriv.concept"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.limit.weierstrass", "label": "Weierstrass extreme value theorem", "domain": "calculus", "subdomain": "calculus.limits", "grading_type": "symbolic", "prerequisites": ["calc.limit.continuity"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.limit.discontinuity", "label": "Types of discontinuities", "domain": "calculus", "subdomain": "calculus.limits", "grading_type": "symbolic", "prerequisites": ["calc.limit.continuity"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.limit.uniform_continuity", "label": "Uniform continuity", "domain": "calculus", "subdomain": "calculus.limits", "grading_type": "symbolic", "prerequisites": ["calc.limit.continuity"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},

    # calculus - derivatives
    {"id": "calc.deriv.difference_quotient", "label": "Difference quotient", "domain": "calculus", "subdomain": "calculus.derivatives", "grading_type": "symbolic", "prerequisites": ["calc.limit.concept"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.deriv.differential", "label": "Differential of a function", "domain": "calculus", "subdomain": "calculus.derivatives", "grading_type": "symbolic", "prerequisites": ["calc.deriv.concept"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.deriv.partial", "label": "Partial derivatives", "domain": "calculus", "subdomain": "calculus.advanced", "grading_type": "symbolic", "prerequisites": ["calc.deriv.concept"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.deriv.non_differentiability", "label": "Points of non-differentiability", "domain": "calculus", "subdomain": "calculus.derivatives", "grading_type": "symbolic", "prerequisites": ["calc.deriv.concept"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.deriv.rolle", "label": "Rolle's theorem", "domain": "calculus", "subdomain": "calculus.derivatives", "grading_type": "symbolic", "prerequisites": ["calc.limit.continuity", "calc.deriv.concept"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.deriv.mvt", "label": "Mean value theorem (Lagrange)", "domain": "calculus", "subdomain": "calculus.derivatives", "grading_type": "symbolic", "prerequisites": ["calc.deriv.rolle"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.deriv.cauchy_mvt", "label": "Cauchy's mean value theorem", "domain": "calculus", "subdomain": "calculus.advanced", "grading_type": "symbolic", "prerequisites": ["calc.deriv.mvt"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.deriv.fermat", "label": "Fermat's theorem on stationary points", "domain": "calculus", "subdomain": "calculus.derivatives", "grading_type": "symbolic", "prerequisites": ["calc.deriv.concept"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.deriv.convexity", "label": "Convexity and concavity of functions", "domain": "calculus", "subdomain": "calculus.derivatives", "grading_type": "symbolic", "prerequisites": ["calc.deriv.concept"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},

    # calculus - integrals
    {"id": "calc.integral.improper", "label": "Improper integrals", "domain": "calculus", "subdomain": "calculus.integrals", "grading_type": "symbolic", "prerequisites": ["calc.integral.definite"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.integral.trig_integrals", "label": "Integrals of trigonometric functions", "domain": "calculus", "subdomain": "calculus.integrals", "grading_type": "symbolic", "prerequisites": ["calc.integral.indefinite"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.integral.trig_substitution", "label": "Trigonometric substitution for integrals", "domain": "calculus", "subdomain": "calculus.integrals", "grading_type": "symbolic", "prerequisites": ["calc.integral.substitution"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.integral.weierstrass_sub", "label": "Weierstrass substitution (tangent half-angle)", "domain": "calculus", "subdomain": "calculus.advanced", "grading_type": "symbolic", "prerequisites": ["calc.integral.substitution"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.integral.numerical", "label": "Numerical integration", "domain": "calculus", "subdomain": "calculus.integrals", "grading_type": "symbolic", "prerequisites": ["calc.integral.definite"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.integral.arc_length", "label": "Arc length of a curve", "domain": "calculus", "subdomain": "calculus.integrals", "grading_type": "symbolic", "prerequisites": ["calc.integral.definite"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.integral.exp_integral", "label": "Integral of exponential functions", "domain": "calculus", "subdomain": "calculus.integrals", "grading_type": "symbolic", "prerequisites": ["calc.integral.indefinite"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.integral.riemann_criteria", "label": "Riemann integrability criteria", "domain": "calculus", "subdomain": "calculus.integrals", "grading_type": "symbolic", "prerequisites": ["calc.integral.definite"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},

    # calculus - sequences
    {"id": "calc.seq.concept", "label": "Sequences: definition and notation", "domain": "calculus", "subdomain": "calculus.sequences", "grading_type": "symbolic", "prerequisites": ["calc.limit.concept"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.seq.convergence", "label": "Convergent and divergent sequences", "domain": "calculus", "subdomain": "calculus.sequences", "grading_type": "symbolic", "prerequisites": ["calc.seq.concept"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.seq.monotone", "label": "Monotone sequences", "domain": "calculus", "subdomain": "calculus.sequences", "grading_type": "symbolic", "prerequisites": ["calc.seq.convergence"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.seq.cauchy", "label": "Cauchy sequences", "domain": "calculus", "subdomain": "calculus.sequences", "grading_type": "symbolic", "prerequisites": ["calc.seq.convergence"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.seq.euler", "label": "Euler's number as a limit of a sequence", "domain": "calculus", "subdomain": "calculus.sequences", "grading_type": "symbolic", "prerequisites": ["calc.seq.convergence"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.seq.function_sequences", "label": "Sequences of functions", "domain": "calculus", "subdomain": "calculus.sequences", "grading_type": "symbolic", "prerequisites": ["calc.seq.convergence"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},

    # calculus - series
    {"id": "calc.series.concept", "label": "Series: definition and convergence", "domain": "calculus", "subdomain": "calculus.series", "grading_type": "symbolic", "prerequisites": ["calc.seq.convergence"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.series.harmonic", "label": "Harmonic series", "domain": "calculus", "subdomain": "calculus.series", "grading_type": "symbolic", "prerequisites": ["calc.series.concept"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.series.positive_terms", "label": "Series with positive terms", "domain": "calculus", "subdomain": "calculus.series", "grading_type": "symbolic", "prerequisites": ["calc.series.concept"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.series.cauchy_criterion", "label": "Cauchy's convergence criterion for series", "domain": "calculus", "subdomain": "calculus.series", "grading_type": "symbolic", "prerequisites": ["calc.series.concept"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.series.integral_test", "label": "Integral test for series convergence", "domain": "calculus", "subdomain": "calculus.series", "grading_type": "symbolic", "prerequisites": ["calc.series.positive_terms"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.series.root_test", "label": "Root test for series convergence", "domain": "calculus", "subdomain": "calculus.series", "grading_type": "symbolic", "prerequisites": ["calc.series.positive_terms"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.series.alternating", "label": "Alternating series and Leibniz criterion", "domain": "calculus", "subdomain": "calculus.series", "grading_type": "symbolic", "prerequisites": ["calc.series.concept"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.series.function_series", "label": "Series of functions", "domain": "calculus", "subdomain": "calculus.series", "grading_type": "symbolic", "prerequisites": ["calc.series.concept"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.series.power", "label": "Power series", "domain": "calculus", "subdomain": "calculus.series", "grading_type": "symbolic", "prerequisites": ["calc.series.concept"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.series.taylor", "label": "Taylor series", "domain": "calculus", "subdomain": "calculus.series", "grading_type": "symbolic", "prerequisites": ["calc.series.power"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "calc.series.fourier", "label": "Fourier series", "domain": "calculus", "subdomain": "calculus.series", "grading_type": "symbolic", "prerequisites": ["calc.series.function_series"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},

    # algebra - equations
    {"id": "alg.eq.absolute_value", "label": "Equations involving absolute value", "domain": "algebra", "subdomain": "algebra.equations", "grading_type": "symbolic", "prerequisites": ["alg.eq.multi_step"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "alg.eq.rational", "label": "Rational equations", "domain": "algebra", "subdomain": "algebra.equations", "grading_type": "symbolic", "prerequisites": ["alg.eq.multi_step"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "alg.eq.irrational", "label": "Irrational equations (equations with radicals)", "domain": "algebra", "subdomain": "algebra.equations", "grading_type": "symbolic", "prerequisites": ["alg.eq.multi_step"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "alg.eq.exponential", "label": "Exponential equations", "domain": "algebra", "subdomain": "algebra.equations", "grading_type": "symbolic", "prerequisites": ["alg.exp.concept"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "alg.eq.logarithmic", "label": "Logarithmic equations", "domain": "algebra", "subdomain": "algebra.equations", "grading_type": "symbolic", "prerequisites": ["alg.log.properties"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "alg.eq.polynomial", "label": "Polynomial equations (higher degree)", "domain": "algebra", "subdomain": "algebra.equations", "grading_type": "symbolic", "prerequisites": ["alg.quad.formula"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "alg.eq.binomial", "label": "Binomial equations", "domain": "algebra", "subdomain": "algebra.equations", "grading_type": "symbolic", "prerequisites": ["alg.eq.multi_step"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "alg.eq.trinomial", "label": "Trinomial equations", "domain": "algebra", "subdomain": "algebra.equations", "grading_type": "symbolic", "prerequisites": ["alg.quad.formula"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "alg.eq.extraneous_roots", "label": "Extraneous solutions and loss of roots", "domain": "algebra", "subdomain": "algebra.equations", "grading_type": "symbolic", "prerequisites": ["alg.eq.rational", "alg.eq.irrational"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "alg.quad.quadratic", "label": "Quadratic equations: overview and solving methods", "domain": "algebra", "subdomain": "algebra.quadratics", "grading_type": "symbolic", "prerequisites": ["alg.quad.solve_factor", "alg.quad.complete_square", "alg.quad.formula"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "alg.quad.incomplete", "label": "Incomplete quadratic equations", "domain": "algebra", "subdomain": "algebra.quadratics", "grading_type": "symbolic", "prerequisites": ["alg.quad.solve_factor"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "alg.quad.complex", "label": "Quadratic equations with complex solutions", "domain": "algebra", "subdomain": "algebra.quadratics", "grading_type": "symbolic", "prerequisites": ["alg.quad.formula", "complex.basics.concept"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "alg.quad.parametric", "label": "Quadratic equations with parameters", "domain": "algebra", "subdomain": "algebra.quadratics", "grading_type": "symbolic", "prerequisites": ["alg.quad.formula"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},

    # algebra - inequalities
    {"id": "alg.ineq.quadratic", "label": "Quadratic inequalities", "domain": "algebra", "subdomain": "algebra.inequalities", "grading_type": "symbolic", "prerequisites": ["alg.ineq.multi_step", "alg.quad.formula"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "alg.ineq.rational", "label": "Rational inequalities", "domain": "algebra", "subdomain": "algebra.inequalities", "grading_type": "symbolic", "prerequisites": ["alg.ineq.multi_step"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "alg.ineq.irrational", "label": "Irrational inequalities", "domain": "algebra", "subdomain": "algebra.inequalities", "grading_type": "symbolic", "prerequisites": ["alg.ineq.multi_step"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "alg.ineq.logarithmic", "label": "Logarithmic inequalities", "domain": "algebra", "subdomain": "algebra.inequalities", "grading_type": "symbolic", "prerequisites": ["alg.log.properties", "alg.ineq.multi_step"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "alg.ineq.absolute_value", "label": "Inequalities with absolute value", "domain": "algebra", "subdomain": "algebra.inequalities", "grading_type": "symbolic", "prerequisites": ["alg.ineq.multi_step"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "alg.ineq.sign_analysis", "label": "Sign analysis method for inequalities", "domain": "algebra", "subdomain": "algebra.inequalities", "grading_type": "symbolic", "prerequisites": ["alg.ineq.multi_step"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "alg.ineq.systems", "label": "Systems of inequalities", "domain": "algebra", "subdomain": "algebra.inequalities", "grading_type": "symbolic", "prerequisites": ["alg.ineq.multi_step"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},

    # algebra - functions
    {"id": "alg.func.domain", "label": "Determining the domain of a function", "domain": "algebra", "subdomain": "algebra.functions", "grading_type": "symbolic", "prerequisites": ["alg.func.concept"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "alg.func.even_odd", "label": "Even and odd functions", "domain": "algebra", "subdomain": "algebra.functions", "grading_type": "symbolic", "prerequisites": ["alg.func.concept"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "alg.func.monotonicity", "label": "Increasing and decreasing functions", "domain": "algebra", "subdomain": "algebra.functions", "grading_type": "symbolic", "prerequisites": ["alg.func.concept"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "alg.func.composite", "label": "Composite functions", "domain": "algebra", "subdomain": "algebra.functions", "grading_type": "symbolic", "prerequisites": ["alg.func.concept"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "alg.func.inverse", "label": "Inverse functions", "domain": "algebra", "subdomain": "algebra.functions", "grading_type": "symbolic", "prerequisites": ["alg.func.composite"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "alg.func.absolute_value", "label": "Absolute value function", "domain": "algebra", "subdomain": "algebra.functions", "grading_type": "symbolic", "prerequisites": ["alg.func.concept"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "alg.func.rational", "label": "Rational functions", "domain": "algebra", "subdomain": "algebra.functions", "grading_type": "symbolic", "prerequisites": ["alg.func.concept"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "alg.func.graph_analysis", "label": "Analyzing graphs of functions", "domain": "algebra", "subdomain": "algebra.functions", "grading_type": "symbolic", "prerequisites": ["alg.func.concept", "alg.func.monotonicity"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "alg.func.sign", "label": "Sign function", "domain": "algebra", "subdomain": "algebra.functions", "grading_type": "symbolic", "prerequisites": ["alg.func.concept"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "alg.func.sigmoid", "label": "Sigmoid function (logistic function)", "domain": "algebra", "subdomain": "algebra.functions", "grading_type": "symbolic", "prerequisites": ["alg.exp.concept"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "alg.func.dirichlet", "label": "Dirichlet function and pathological functions", "domain": "algebra", "subdomain": "algebra.functions", "grading_type": "symbolic", "prerequisites": ["alg.func.concept"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},

    # algebra - polynomials
    {"id": "alg.poly.division", "label": "Polynomial division", "domain": "algebra", "subdomain": "algebra.polynomials", "grading_type": "symbolic", "prerequisites": ["alg.poly.concept"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "alg.poly.roots", "label": "Roots of a polynomial", "domain": "algebra", "subdomain": "algebra.polynomials", "grading_type": "symbolic", "prerequisites": ["alg.poly.concept", "alg.quad.formula"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "alg.poly.synthetic_division", "label": "Synthetic division method", "domain": "algebra", "subdomain": "algebra.polynomials", "grading_type": "symbolic", "prerequisites": ["alg.poly.division"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "alg.poly.vieta", "label": "Vieta's formulas", "domain": "algebra", "subdomain": "algebra.polynomials", "grading_type": "symbolic", "prerequisites": ["alg.poly.roots"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},

    # algebra - conics
    {"id": "alg.conic.parabola", "label": "Parabola: equation and graph", "domain": "algebra", "subdomain": "algebra.conics", "grading_type": "symbolic", "prerequisites": ["alg.func.quad"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "alg.conic.hyperbola", "label": "Hyperbola: equation and graph", "domain": "algebra", "subdomain": "algebra.conics", "grading_type": "symbolic", "prerequisites": ["alg.conic.circle"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},

    # algebra - systems
    {"id": "alg.systems.concept", "label": "Systems of linear equations: introduction", "domain": "algebra", "subdomain": "algebra.systems", "grading_type": "symbolic", "prerequisites": ["alg.eq.multi_step"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "alg.systems.gaussian", "label": "Gaussian elimination for solving linear systems", "domain": "algebra", "subdomain": "algebra.systems", "grading_type": "symbolic", "prerequisites": ["alg.systems.concept"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},

    # linear algebra
    {"id": "linalg.matrix.rank", "label": "Rank of a matrix", "domain": "linear_algebra", "subdomain": "linear_algebra.matrices", "grading_type": "symbolic", "prerequisites": ["linalg.det.2x2", "linalg.det.3x3"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "linalg.eigen.diagonalization", "label": "Matrix diagonalization", "domain": "linear_algebra", "subdomain": "linear_algebra.eigenvalues", "grading_type": "symbolic", "prerequisites": ["linalg.eigen.compute"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "linalg.vector.parametric", "label": "Vector and parametric equations of a line", "domain": "linear_algebra", "subdomain": "linear_algebra.vectors", "grading_type": "symbolic", "prerequisites": ["linalg.vector.concept"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "linalg.vector.cosine_similarity", "label": "Cosine similarity between vectors", "domain": "linear_algebra", "subdomain": "linear_algebra.vectors", "grading_type": "symbolic", "prerequisites": ["linalg.vector.dot"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},

    # geometry
    {"id": "geo.coord.lines", "label": "Lines in the coordinate plane", "domain": "geometry", "subdomain": "geometry.coordinate", "grading_type": "symbolic", "prerequisites": ["alg.linear.slope_intercept"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "geo.coord.polar", "label": "Polar coordinates", "domain": "geometry", "subdomain": "geometry.coordinate", "grading_type": "symbolic", "prerequisites": ["geo.coord.distance"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},

    # discrete math
    {"id": "discrete.combinatorics.binomial_theorem", "label": "Binomial theorem", "domain": "discrete_math", "subdomain": "discrete_math.combinatorics", "grading_type": "symbolic", "prerequisites": ["discrete.combinatorics.combinations"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},

    # statistics
    # Advanced means/variance
    {"id": "stat.adv.geometric_mean", "label": "Geometric mean", "domain": "statistics", "subdomain": "statistics.advanced", "grading_type": "symbolic", "prerequisites": ["stat.data.mean"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "stat.adv.harmonic_mean", "label": "Harmonic mean", "domain": "statistics", "subdomain": "statistics.advanced", "grading_type": "symbolic", "prerequisites": ["stat.data.mean"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "stat.adv.rms", "label": "Root mean square", "domain": "statistics", "subdomain": "statistics.advanced", "grading_type": "symbolic", "prerequisites": ["stat.data.mean"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "stat.adv.variance", "label": "Variance", "domain": "statistics", "subdomain": "statistics.advanced", "grading_type": "symbolic", "prerequisites": ["stat.data.mean"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "stat.adv.covariance", "label": "Variance and covariance of random variables", "domain": "statistics", "subdomain": "statistics.advanced", "grading_type": "symbolic", "prerequisites": ["stat.adv.variance"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    # Probability - random variables
    {"id": "stat.prob.discrete_rv", "label": "Discrete random variables", "domain": "statistics", "subdomain": "statistics.probability", "grading_type": "symbolic", "prerequisites": ["stat.prob.basic"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "stat.prob.continuous_rv", "label": "Continuous random variables", "domain": "statistics", "subdomain": "statistics.probability", "grading_type": "symbolic", "prerequisites": ["stat.prob.basic"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "stat.prob.expected_value", "label": "Expected value of a random variable", "domain": "statistics", "subdomain": "statistics.probability", "grading_type": "symbolic", "prerequisites": ["stat.prob.discrete_rv"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "stat.prob.bayes", "label": "Bayes' theorem", "domain": "statistics", "subdomain": "statistics.probability", "grading_type": "symbolic", "prerequisites": ["stat.prob.basic"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    # Probability distributions
    {"id": "stat.dist.bernoulli", "label": "Bernoulli distribution", "domain": "statistics", "subdomain": "statistics.distributions", "grading_type": "symbolic", "prerequisites": ["stat.prob.discrete_rv"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "stat.dist.binomial", "label": "Binomial distribution", "domain": "statistics", "subdomain": "statistics.distributions", "grading_type": "symbolic", "prerequisites": ["stat.dist.bernoulli"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "stat.dist.geometric", "label": "Geometric distribution", "domain": "statistics", "subdomain": "statistics.distributions", "grading_type": "symbolic", "prerequisites": ["stat.dist.bernoulli"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "stat.dist.hypergeometric", "label": "Hypergeometric distribution", "domain": "statistics", "subdomain": "statistics.distributions", "grading_type": "symbolic", "prerequisites": ["stat.dist.binomial"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "stat.dist.poisson", "label": "Poisson distribution", "domain": "statistics", "subdomain": "statistics.distributions", "grading_type": "symbolic", "prerequisites": ["stat.dist.binomial"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "stat.dist.uniform", "label": "Uniform distribution", "domain": "statistics", "subdomain": "statistics.distributions", "grading_type": "symbolic", "prerequisites": ["stat.prob.continuous_rv"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "stat.dist.normal", "label": "Normal (Gaussian) distribution", "domain": "statistics", "subdomain": "statistics.distributions", "grading_type": "symbolic", "prerequisites": ["stat.prob.continuous_rv"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "stat.dist.z_table", "label": "Standard normal distribution and Z-table", "domain": "statistics", "subdomain": "statistics.distributions", "grading_type": "symbolic", "prerequisites": ["stat.dist.normal"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "stat.dist.exponential", "label": "Exponential distribution", "domain": "statistics", "subdomain": "statistics.distributions", "grading_type": "symbolic", "prerequisites": ["stat.prob.continuous_rv"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "stat.dist.gamma", "label": "Gamma distribution", "domain": "statistics", "subdomain": "statistics.distributions", "grading_type": "symbolic", "prerequisites": ["stat.dist.exponential"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "stat.dist.beta", "label": "Beta distribution", "domain": "statistics", "subdomain": "statistics.distributions", "grading_type": "symbolic", "prerequisites": ["stat.dist.gamma"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "stat.dist.chi_square", "label": "Chi-square distribution", "domain": "statistics", "subdomain": "statistics.distributions", "grading_type": "symbolic", "prerequisites": ["stat.dist.normal"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "stat.dist.student_t", "label": "Student's t-distribution", "domain": "statistics", "subdomain": "statistics.distributions", "grading_type": "symbolic", "prerequisites": ["stat.dist.normal"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    # Inference
    {"id": "stat.infer.sampling", "label": "Sampling distributions", "domain": "statistics", "subdomain": "statistics.inference", "grading_type": "symbolic", "prerequisites": ["stat.dist.normal"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
    {"id": "stat.infer.confidence", "label": "Confidence intervals", "domain": "statistics", "subdomain": "statistics.inference", "grading_type": "symbolic", "prerequisites": ["stat.infer.sampling"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},

    # prealgebra
    {"id": "prealg.types", "label": "Types of numbers", "domain": "prealgebra", "subdomain": "prealgebra.variables", "grading_type": "symbolic", "prerequisites": ["prealg.real.concept"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},

    # machine learning (new domain)
    {"id": "ml.backpropagation", "label": "Backpropagation in neural networks", "domain": "machine_learning", "subdomain": "machine_learning.neural_networks", "grading_type": "symbolic", "prerequisites": ["calc.deriv.chain_rule", "linalg.matrix.mult"], "mastery_threshold": {"streak": 3, "avg_time_seconds": 120}},
]


def load_json(path):
    with open(path) as f:
        return json.load(f)


def save_json(path, data):
    with open(path, "w") as f:
        json.dump(data, f, indent=2)
        f.write("\n")


def main():
    # Build lessons.json from the mapping
    lessons = []
    for filepath, concept_ids in sorted(MAPPING.items()):
        full_path = os.path.join(LESSONS_DIR, filepath)
        if not os.path.exists(full_path):
            print(f"WARNING: Lesson file not found: {full_path}")
            continue
        for cid in concept_ids:
            lessons.append({
                "concept_id": cid,
                "source": f"algebrica/{filepath}"
            })

    # Save lessons.json
    save_json(LESSONS_JSON, lessons)
    print(f"Saved {len(lessons)} entries to lessons.json")

    # Add new concept entries to existing JSON files
    # Group new concepts by domain file
    domain_files = {
        "complex_numbers.json": [],
        "trigonometry.json": [],
        "calculus.json": [],
        "algebra.json": [],
        "linear_algebra.json": [],
        "geometry.json": [],
        "discrete_math.json": [],
        "statistics.json": [],
        "prealgebra.json": [],
        "machine_learning.json": [],
    }

    for c in NEW_CONCEPTS:
        domain = c["domain"]
        # Map domain name to filename
        domain_to_file = {
            "complex_numbers": "complex_numbers.json",
            "trigonometry": "trigonometry.json",
            "calculus": "calculus.json",
            "algebra": "algebra.json",
            "linear_algebra": "linear_algebra.json",
            "geometry": "geometry.json",
            "discrete_math": "discrete_math.json",
            "statistics": "statistics.json",
            "prealgebra": "prealgebra.json",
            "machine_learning": "machine_learning.json",
        }
        fname = domain_to_file.get(domain)
        if fname:
            domain_files.setdefault(fname, []).append(c)
        else:
            print(f"WARNING: No domain file mapped for {c['id']} (domain={domain})")

    for fname, new_concepts in domain_files.items():
        if not new_concepts:
            continue
        fpath = os.path.join(CONCEPTS_DIR, fname)
        existing = []
        if os.path.exists(fpath):
            existing = load_json(fpath)
        # Filter out any already-present IDs
        existing_ids = {c["id"] for c in existing}
        to_add = [c for c in new_concepts if c["id"] not in existing_ids]
        if not to_add:
            continue
        # Add new concepts
        combined = existing + to_add
        save_json(fpath, combined)
        print(f"Added {len(to_add)} new concepts to {fname} (total: {len(combined)})")

    # Print summary
    existing_lessons = []
    if os.path.exists(LESSONS_JSON):
        existing_lessons = load_json(LESSONS_JSON)
    unique_sources = set(e["source"] for e in existing_lessons)
    print(f"\nSummary:")
    print(f"  Lessons mapped: {len(existing_lessons)}")
    print(f"  Unique lesson files: {len(unique_sources)}")
    print(f"  New concept IDs added: {len(NEW_CONCEPTS)}")


if __name__ == "__main__":
    main()
