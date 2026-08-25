"""
Build mapping of Wikipedia concept IDs → source sections (ORCCA or OpenStax Prealgebra 1e)
"""
import json

ORCCA_SECTIONS = {
    "arith.abs_value": "section-absolute-value-and-square-root",
    "arith.sqrt.perfect": "section-absolute-value-and-square-root",
    "arith.neg.number_line": "section-arithmetic-with-negative-numbers",
    "arith.neg.add_sub": "section-arithmetic-with-negative-numbers",
    "arith.neg.mult_div": "section-arithmetic-with-negative-numbers",
    "arith.neg.order_ops": "section-arithmetic-with-negative-numbers",
    "arith.order_ops.basic": "section-order-of-operations",
    "arith.order_ops.full": "section-order-of-operations",
    "arith.order_ops.nested": "section-order-of-operations",
    "arith.sci_notation": "section-scientific-notation",
    "arith.sci_notation.ops": "section-scientific-notation",
    "arith.exp.neg": "section-introduction-to-exponent-properties",
    "arith.exp.zero": "section-introduction-to-exponent-properties",
    "arith.factor.gcf": "section-factoring-out-the-common-factor",
    "prealg.var.concept": "section-variables-and-evaluating-expressions",
    "prealg.expr.evaluate": "section-variables-and-evaluating-expressions",
    "prealg.expr.distribute": "section-algebraic-properties-and-simplifying-expressions",
    "prealg.expr.like_terms": "section-combining-like-terms",
    "prealg.eq.one_step_add": "section-solving-one-step-equations",
    "prealg.eq.one_step_mult": "section-solving-one-step-equations",
    "prealg.eq.two_step": "section-solving-multistep-linear-equations",
    "prealg.eq.word": "section-modeling-with-equations-and-inequalities",
    "prealg.ineq.one_step": "section-solving-one-step-inequalities",
    "prealg.ineq.two_step": "section-solving-multistep-linear-inequalities",
    "alg.linear.slope": "section-slope",
    "alg.linear.graph": "section-slope-intercept-form",
    "alg.linear.standard_form": "section-standard-form",
    "alg.linear.parallel_perp": "section-geometry-of-lines",
    "alg.func.notation": "section-function-basics",
    "alg.func.evaluate": "section-function-basics",
    "alg.func.linear": "section-technical-definition-of-a-function",
    "alg.exp.evaluate": "section-introduction-to-exponent-properties",
    "alg.eq.vars_both_sides": "section-solving-multistep-linear-equations",
    "alg.factor.gcf": "section-factoring-out-the-common-factor",
    "alg.factor.diff_squares": "section-factoring-special-polynomials",
    "alg.poly.mult_mono": "section-multiplying-polynomials",
    "alg.ineq.compound": "section-compound-inequalities",
    "alg.ineq.two_var": "section-solving-inequalities-graphically",
    "alg.quad.discriminant": "section-the-quadratic-formula",
    "alg.systems.word": "section-modeling-with-equations-and-inequalities",
    "alg.seq.sum_arith": "section-function-basics",
    "geo.coord.plot": "section-cartesian-coordinates",
    "geo.coord.distance": "section-cartesian-coordinates",
    "geo.coord.midpoint": "section-cartesian-coordinates",
    "stat.data.bar_graph": "section-overview-of-graphing",
    "stat.data.line_plot": "section-exploring-two-variable-data-and-rate-of-change",
    "stat.data.read_table": "section-exploring-two-variable-data-and-rate-of-change",
    "stat.data.mode": "section-exploring-two-variable-data-and-rate-of-change",
    "stat.data.range": "section-exploring-two-variable-data-and-rate-of-change",
}

OPENSTAX_SECTIONS = {
    "arith.add.single": "1-2-add-whole-numbers",
    "arith.add.double": "1-2-add-whole-numbers",
    "arith.add.carry": "1-2-add-whole-numbers",
    "arith.add.triple": "1-2-add-whole-numbers",
    "arith.add.word": "1-2-add-whole-numbers",
    "arith.sub.single": "1-3-subtract-whole-numbers",
    "arith.sub.double": "1-3-subtract-whole-numbers",
    "arith.sub.borrow": "1-3-subtract-whole-numbers",
    "arith.sub.word": "1-3-subtract-whole-numbers",
    "arith.mult.concept": "1-4-multiply-whole-numbers",
    "arith.mult.tables": "1-4-multiply-whole-numbers",
    "arith.mult.double": "1-4-multiply-whole-numbers",
    "arith.mult.triple": "1-4-multiply-whole-numbers",
    "arith.mult.word": "1-4-multiply-whole-numbers",
    "arith.mult.2_5_10": "1-4-multiply-whole-numbers",
    "arith.div.concept": "1-5-divide-whole-numbers",
    "arith.div.basic": "1-5-divide-whole-numbers",
    "arith.div.long": "1-5-divide-whole-numbers",
    "arith.div.remainder": "1-5-divide-whole-numbers",
    "arith.div.word": "1-5-divide-whole-numbers",
    "arith.place.tens": "1-1-introduction-to-whole-numbers",
    "arith.place.hundreds": "1-1-introduction-to-whole-numbers",
    "arith.place.thousands": "1-1-introduction-to-whole-numbers",
    "arith.round.tens": "1-1-introduction-to-whole-numbers",
    "arith.round.hundreds": "1-1-introduction-to-whole-numbers",
    "arith.round.thousands": "1-1-introduction-to-whole-numbers",
    "arith.dec.intro": "5-1-decimals",
    # Fractions: split across Prealgebra ch.4 topic sections.
    "frac.concept": "4-1-visualize-fractions",
    "frac.parts": "4-1-visualize-fractions",
    "frac.on_number_line": "4-1-visualize-fractions",
    "frac.equivalent": "4-1-visualize-fractions",
    "frac.simplify": "4-1-visualize-fractions",
    "frac.compare": "4-1-visualize-fractions",
    "frac.benchmark": "4-1-visualize-fractions",
    "frac.mult": "4-2-multiply-and-divide-fractions",
    "frac.mult.whole": "4-2-multiply-and-divide-fractions",
    "frac.div": "4-2-multiply-and-divide-fractions",
    "frac.div.whole": "4-2-multiply-and-divide-fractions",
    "frac.mixed.convert": "4-3-multiply-and-divide-mixed-numbers-and-complex-fractions",
    "frac.mixed.mult": "4-3-multiply-and-divide-mixed-numbers-and-complex-fractions",
    "frac.add.same": "4-4-add-and-subtract-fractions-with-common-denominators",
    "frac.sub.same": "4-4-add-and-subtract-fractions-with-common-denominators",
    "frac.add.diff": "4-5-add-and-subtract-fractions-with-different-denominators",
    "frac.sub.diff": "4-5-add-and-subtract-fractions-with-different-denominators",
    "frac.mixed.add": "4-6-add-and-subtract-mixed-numbers",
    "frac.mixed.sub": "4-6-add-and-subtract-mixed-numbers",
    "frac.to_decimal": "5-3-decimals-and-fractions",
    "dec.add": "5-2-decimal-operations",
    "dec.sub": "5-2-decimal-operations",
    "dec.mult": "5-2-decimal-operations",
    "dec.div": "5-2-decimal-operations",
    "dec.from_frac": "5-3-decimals-and-fractions",
    "dec.to_frac": "5-3-decimals-and-fractions",
    "dec.compare": "5-1-decimals",
    "dec.round": "5-1-decimals",
    # Percents: split across Prealgebra ch.6.
    "pct.concept": "6-1-understand-percent",
    "pct.from_dec": "6-1-understand-percent",
    "pct.to_dec": "6-1-understand-percent",
    "pct.of_number": "6-2-solve-general-applications-of-percent",
    "pct.find_rate": "6-2-solve-general-applications-of-percent",
    "pct.increase": "6-2-solve-general-applications-of-percent",
    "pct.discount": "6-3-solve-sales-tax-commission-and-discount-applications",
    "pct.tax_tip": "6-3-solve-sales-tax-commission-and-discount-applications",
    # Ratios and proportions.
    "ratio.concept": "5-6-ratios-and-rate",
    "ratio.simplify": "5-6-ratios-and-rate",
    "ratio.rate": "5-6-ratios-and-rate",
    "ratio.scale": "5-6-ratios-and-rate",
    "ratio.proportion": "6-5-solve-proportions-and-their-applications",
    # Factors, primes, LCM.
    "arith.factor.prime": "2-4-find-multiples-and-factors",
    "arith.factor.prime_fact": "2-4-find-multiples-and-factors",
    "arith.factor.find": "2-4-find-multiples-and-factors",
    "arith.factor.composite": "2-4-find-multiples-and-factors",
    "arith.factor.lcm": "2-5-prime-factorization-and-the-least-common-multiple",
    # Probability basics.
    "stat.prob.basic": "5-5-averages-and-probability",
    "stat.prob.complement": "5-5-averages-and-probability",
    "stat.prob.compound": "5-5-averages-and-probability",
    "stat.prob.sample_space": "5-5-averages-and-probability",
    "stat.prob.counting": "5-5-averages-and-probability",
}

# Geometry lessons are hand-authored (ORCCA has no dedicated geometry
# coverage; see data/lessons/teaching/geo.*.md).
AUTHORED_SECTIONS = {
    "geo.basic.points_lines",
    "geo.basic.angles",
    "geo.basic.angle_measure",
    "geo.basic.complementary",
    "geo.basic.vertical",
    "geo.triangle.types",
    "geo.triangle.angles",
    "geo.triangle.area",
    "geo.quad.types",
    "geo.quad.perimeter",
    "geo.quad.area",
    "geo.circle.parts",
    "geo.circle.circumference",
    "geo.circle.area",
    "geo.solid.surface_area",
    "geo.solid.volume_rect",
}

APPLIED_CALC_SECTIONS = {
    "calc.limit.numeric": "limits",
    "calc.limit.infinity": "limits",
    "calc.deriv.sum_rule": "derivative-rules",
    "calc.deriv.product_rule": "derivative-rules",
    "calc.deriv.quotient_rule": "derivative-rules",
    "calc.deriv.exp_log": "derivative-rules",
    "calc.deriv.trig": "derivative-rules",
    "calc.deriv.implicit": "implicit-related-rates",
    "calc.deriv.related_rates": "implicit-related-rates",
    "calc.integral.power_rule": "integrals",
    "calc.integral.volume": "integrals",
}

DISCRETE_MATH_SECTIONS = {
    "discrete.logic.connectives": "prop-logic",
    "discrete.logic.truth_tables": "prop-logic",
    "discrete.logic.quantifiers": "prop-logic",
    "discrete.sets.venn": "sets",
    "discrete.graphs.basics": "graph-theory",
    "discrete.graphs.paths": "graph-theory",
    "discrete.graphs.trees": "graph-theory",
    "discrete.recurrence": "sequences",
}

OPENSTAX_PRECALC_SECTIONS = {
    "trig.radians": "5-1-angles",
    "trig.period": "5-4-right-triangle-trigonometry",
}

HEFFERON_LINALG_SECTIONS = {
    "linalg.systems.matrix": "linear-systems",
}

UNCOVERED = {
    "count.backwards": "manual",
    "count.cardinality": "manual",
    "count.compare": "manual",
    "count.number_line": "manual",
    "count.objects": "manual",
    "count.objects_20": "manual",
    "count.ordinal": "manual",
    "count.skip_2": "manual",
    "count.skip_5": "manual",
    "count.skip_10": "manual",
    "ode.separable": "manual",
    "ode.linear_first": "manual",
    "ode.homogeneous": "manual",
    "ode.nonhomogeneous": "manual",
    "ode.exact": "manual",
    "ode.laplace": "manual",
    "ode.systems": "manual",
    "nt.gcd_euclidean": "manual",
    "nt.diophantine": "manual",
    "nt.fermat_little": "manual",
    "nt.crypto": "manual",
    "topo.continuous": "manual",
    "topo.metric": "manual",
    "topo.open_closed": "manual",
}

def main():
    section_dicts = [
        ("orcca", ORCCA_SECTIONS),
        ("openstax", OPENSTAX_SECTIONS),
        ("applied_calc", APPLIED_CALC_SECTIONS),
        ("discrete_math", DISCRETE_MATH_SECTIONS),
        ("openstax_precalc", OPENSTAX_PRECALC_SECTIONS),
        ("hefferon_linalg", HEFFERON_LINALG_SECTIONS),
    ]

    all_concepts = {}
    for name, d in section_dicts:
        for cid, slug in d.items():
            all_concepts[cid] = (name, slug)
    all_concepts.update({k: ("manual", None) for k in UNCOVERED})
    all_concepts.update({k: ("authored", None) for k in AUTHORED_SECTIONS})

    print(f"Total concepts: {len(all_concepts)}")

    mapping = {}
    for cid in all_concepts:
        source, slug = all_concepts[cid]
        if source in ("manual", "authored"):
            mapping[cid] = {"source": source}
        else:
            mapping[cid] = {"source": source, "slug": slug}

    counts = {}
    for v in mapping.values():
        s = v["source"]
        counts[s] = counts.get(s, 0) + 1
    for src, cnt in sorted(counts.items()):
        print(f"  {src}: {cnt}")
    print(f"  Total: {len(mapping)}")

    with open("data/lessons/mapping.json", "w") as f:
        json.dump(mapping, f, indent=2)

    print("Mapping written to data/lessons/mapping.json")

if __name__ == "__main__":
    main()
