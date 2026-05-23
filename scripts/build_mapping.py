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
    "arith.factor.lcm": "section-comparison-symbols-and-notation-for-intervals",
    "arith.factor.prime": "section-set-notation-and-types-of-numbers",
    "arith.factor.prime_fact": "section-set-notation-and-types-of-numbers",
    "arith.factor.find": "section-set-notation-and-types-of-numbers",
    "arith.factor.composite": "section-set-notation-and-types-of-numbers",
    "frac.concept": "section-fractions-and-fraction-arithmetic",
    "frac.parts": "section-fractions-and-fraction-arithmetic",
    "frac.on_number_line": "section-fractions-and-fraction-arithmetic",
    "frac.equivalent": "section-fractions-and-fraction-arithmetic",
    "frac.simplify": "section-fractions-and-fraction-arithmetic",
    "frac.compare": "section-fractions-and-fraction-arithmetic",
    "frac.benchmark": "section-fractions-and-fraction-arithmetic",
    "frac.add.same": "section-fractions-and-fraction-arithmetic",
    "frac.sub.same": "section-fractions-and-fraction-arithmetic",
    "frac.add.diff": "section-equations-and-inequalities-with-fractions",
    "frac.sub.diff": "section-equations-and-inequalities-with-fractions",
    "frac.mult": "section-fractions-and-fraction-arithmetic",
    "frac.mult.whole": "section-fractions-and-fraction-arithmetic",
    "frac.div": "section-fractions-and-fraction-arithmetic",
    "frac.div.whole": "section-fractions-and-fraction-arithmetic",
    "frac.mixed.convert": "section-fractions-and-fraction-arithmetic",
    "frac.mixed.add": "section-fractions-and-fraction-arithmetic",
    "frac.mixed.sub": "section-fractions-and-fraction-arithmetic",
    "frac.mixed.mult": "section-fractions-and-fraction-arithmetic",
    "frac.add.word": "section-equations-and-inequalities-with-fractions",
    "frac.to_decimal": "section-fractions-and-fraction-arithmetic",
    "dec.add": "section-fractions-and-fraction-arithmetic",
    "dec.sub": "section-fractions-and-fraction-arithmetic",
    "dec.mult": "section-fractions-and-fraction-arithmetic",
    "dec.div": "section-fractions-and-fraction-arithmetic",
    "dec.from_frac": "section-fractions-and-fraction-arithmetic",
    "dec.to_frac": "section-fractions-and-fraction-arithmetic",
    "dec.compare": "section-comparison-symbols-and-notation-for-intervals",
    "dec.round": "section-comparison-symbols-and-notation-for-intervals",
    "pct.concept": "section-percentages",
    "pct.of_number": "section-percentages",
    "pct.find_rate": "section-percentages",
    "pct.increase": "section-percentages",
    "pct.discount": "section-percentages",
    "pct.tax_tip": "section-percentages",
    "pct.from_dec": "section-percentages",
    "pct.to_dec": "section-percentages",
    "ratio.concept": "section-fractions-and-fraction-arithmetic",
    "ratio.simplify": "section-fractions-and-fraction-arithmetic",
    "ratio.proportion": "section-equations-and-inequalities-with-fractions",
    "ratio.rate": "section-exploring-two-variable-data-and-rate-of-change",
    "ratio.scale": "section-geometry-applications",
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
    "geo.basic.points_lines": "section-geometry-of-lines",
    "geo.basic.angles": "section-geometry-formulas",
    "geo.basic.angle_measure": "section-geometry-formulas",
    "geo.basic.complementary": "section-geometry-formulas",
    "geo.basic.vertical": "section-geometry-formulas",
    "geo.coord.plot": "section-cartesian-coordinates",
    "geo.coord.distance": "section-cartesian-coordinates",
    "geo.coord.midpoint": "section-cartesian-coordinates",
    "geo.quad.types": "section-geometry-applications",
    "geo.quad.perimeter": "section-geometry-applications",
    "geo.quad.area": "section-geometry-applications",
    "geo.triangle.types": "section-geometry-formulas",
    "geo.triangle.angles": "section-geometry-formulas",
    "geo.triangle.area": "section-geometry-applications",
    "geo.circle.parts": "section-geometry-applications",
    "geo.circle.circumference": "section-geometry-applications",
    "geo.circle.area": "section-geometry-applications",
    "geo.solid.surface_area": "section-geometry-applications",
    "geo.solid.volume_rect": "section-geometry-applications",
    "stat.data.bar_graph": "section-overview-of-graphing",
    "stat.data.line_plot": "section-exploring-two-variable-data-and-rate-of-change",
    "stat.data.read_table": "section-exploring-two-variable-data-and-rate-of-change",
    "stat.data.mode": "section-exploring-two-variable-data-and-rate-of-change",
    "stat.data.range": "section-exploring-two-variable-data-and-rate-of-change",
    "stat.prob.basic": "section-set-notation-and-types-of-numbers",
    "stat.prob.complement": "section-set-notation-and-types-of-numbers",
    "stat.prob.compound": "section-set-notation-and-types-of-numbers",
    "stat.prob.sample_space": "section-set-notation-and-types-of-numbers",
    "stat.prob.counting": "section-set-notation-and-types-of-numbers",
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
    "geo.basic.points_lines": "9-1-points-lines-and-angles",
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

    print(f"Total concepts: {len(all_concepts)}")

    mapping = {}
    for cid in all_concepts:
        source, slug = all_concepts[cid]
        if source == "manual":
            mapping[cid] = {"source": "manual"}
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
