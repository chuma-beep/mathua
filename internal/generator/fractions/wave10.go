package fractions

// Wave 10 completes the fractions domain at one hundred and twenty concepts.
// Every question has an integer answer so the numeric grader can check it, and
// every generator publishes the operands the answer test recomputes.

import (
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

var w10Kinds = map[string]string{
	"frac.basics.density":           "density",
	"frac.basics.whole":             "whole_fraction",
	"frac.concept.proper_improper":  "proper_improper",
	"frac.concept.zero_denominator": "zero_denominator",
	"frac.concept.equal_parts":      "equal_parts",
	"frac.concept.unit_meaning":     "unit_meaning",
	"frac.equiv.cross_check":        "cross_check",
	"frac.equiv.simplify_gcf":       "simplify_gcf",
	"frac.equiv.lowest_terms_check": "lowest_terms",
	"frac.equiv.chain":              "equiv_chain",
	"frac.equiv.missing_number":     "missing_number",
	"frac.equiv.benchmark_choice":   "benchmark_choice",
	"frac.equiv.decimal_order":      "decimal_order",
	"frac.equiv.percent_bench":      "percent_bench",
	"frac.add.whole":                "add_whole",
	"frac.add.estimate_sum":         "estimate_sum",
	"frac.add.regroup":              "add_regroup",
	"frac.add.unit_fractions":       "add_unit",
	"frac.add.missing_addend":       "missing_addend",
	"frac.add.associative":          "add_associative",
	"frac.add.same_numerator":       "add_same_numerator",
	"frac.add.decimal_mix":          "add_decimal",
	"frac.sub.missing_minuend":      "missing_minuend",
	"frac.sub.missing_subtrahend":   "missing_subtrahend",
	"frac.sub.estimate_diff":        "estimate_diff",
	"frac.sub.regroup_whole":        "sub_regroup_whole",
	"frac.sub.negative_result":      "sub_negative",
	"frac.sub.three":                "sub_three",
	"frac.sub.decimal_mix":          "sub_decimal",
	"frac.sub.number_line_distance": "sub_distance",
	"frac.mult.whole_cancel":        "mult_whole_cancel",
	"frac.mult.estimate_product":    "estimate_product",
	"frac.mult.scaling":             "mult_scaling",
	"frac.mult.unit":                "mult_unit",
	"frac.mult.three":               "mult_three",
	"frac.mult.mixed_improper":      "mult_mixed_improper",
	"frac.mult.decimal":             "mult_decimal",
	"frac.mult.predict":             "mult_predict",
	"frac.div.reciprocal_check":     "div_reciprocal",
	"frac.div.estimate_quotient":    "estimate_quotient",
	"frac.div.complex_simplify":     "complex_simplify",
	"frac.div.unit":                 "div_unit",
	"frac.div.measurement":          "div_measurement",
	"frac.div.partitive":            "div_partitive",
	"frac.div.three":                "div_three",
	"frac.div.zero":                 "div_zero",
	"frac.mixed.compare_whole":      "mixed_compare",
	"frac.mixed.negative_improper":  "mixed_negative",
	"frac.mixed.borrow_twice":       "mixed_borrow_twice",
	"frac.mixed.div_whole":          "mixed_div_whole",
	"frac.mixed.estimate":           "mixed_estimate",
	"frac.mixed.percent":            "mixed_percent",
	"frac.mixed.word_compare":       "mixed_word_compare",
}

func registerWave10(reg *generator.Registry) {
	for id, kind := range w10Kinds {
		reg.Register(id, &w10Gen{kind: kind})
	}
}

type w10Gen struct{ kind string }

func (g *w10Gen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "density":
		return w10Density(kp)
	case "whole_fraction":
		return w10WholeFraction(kp)
	case "proper_improper":
		return w10ProperImproper(kp)
	case "zero_denominator":
		return w10ZeroDenominator(kp)
	case "equal_parts":
		return w10EqualParts(kp)
	case "unit_meaning":
		return w10UnitMeaning(kp)
	case "cross_check":
		return w10CrossCheck(kp)
	case "simplify_gcf":
		return w10SimplifyGCF(kp)
	case "lowest_terms":
		return w10LowestTerms(kp)
	case "equiv_chain":
		return w10EquivChain(kp)
	case "missing_number":
		return w10MissingNumber(kp)
	case "benchmark_choice":
		return w10BenchmarkChoice(kp)
	case "decimal_order":
		return w10DecimalOrder(kp)
	case "percent_bench":
		return w10PercentBench(kp)
	case "add_whole":
		return w10AddWhole(kp)
	case "estimate_sum":
		return w10EstimateSum(kp)
	case "add_regroup":
		return w10AddRegroup(kp)
	case "add_unit":
		return w10AddUnit(kp)
	case "missing_addend":
		return w10MissingAddend(kp)
	case "add_associative":
		return w10AddAssociative(kp)
	case "add_same_numerator":
		return w10AddSameNumerator(kp)
	case "add_decimal":
		return w10AddDecimal(kp)
	case "missing_minuend":
		return w10MissingMinuend(kp)
	case "missing_subtrahend":
		return w10MissingSubtrahend(kp)
	case "estimate_diff":
		return w10EstimateDiff(kp)
	case "sub_regroup_whole":
		return w10SubRegroupWhole(kp)
	case "sub_negative":
		return w10SubNegative(kp)
	case "sub_three":
		return w10SubThree(kp)
	case "sub_decimal":
		return w10SubDecimal(kp)
	case "sub_distance":
		return w10SubDistance(kp)
	case "mult_whole_cancel":
		return w10MultWholeCancel(kp)
	case "estimate_product":
		return w10EstimateProduct(kp)
	case "mult_scaling":
		return w10MultScaling(kp)
	case "mult_unit":
		return w10MultUnit(kp)
	case "mult_three":
		return w10MultThree(kp)
	case "mult_mixed_improper":
		return w10MultMixedImproper(kp)
	case "mult_decimal":
		return w10MultDecimal(kp)
	case "mult_predict":
		return w10MultPredict(kp)
	case "div_reciprocal":
		return w10DivReciprocal(kp)
	case "estimate_quotient":
		return w10EstimateQuotient(kp)
	case "complex_simplify":
		return w10ComplexSimplify(kp)
	case "div_unit":
		return w10DivUnit(kp)
	case "div_measurement":
		return w10DivMeasurement(kp)
	case "div_partitive":
		return w10DivPartitive(kp)
	case "div_three":
		return w10DivThree(kp)
	case "div_zero":
		return w10DivZero(kp)
	case "mixed_compare":
		return w10MixedCompare(kp)
	case "mixed_negative":
		return w10MixedNegative(kp)
	case "mixed_borrow_twice":
		return w10MixedBorrowTwice(kp)
	case "mixed_div_whole":
		return w10MixedDivWhole(kp)
	case "mixed_estimate":
		return w10MixedEstimate(kp)
	case "mixed_percent":
		return w10MixedPercent(kp)
	case "mixed_word_compare":
		return w10MixedWordCompare(kp)
	}
	return w10Problem(kp, "Answer 0. (enter a number)", 0, "This generator is not implemented.", nil, nil)
}

func w10R(min, max int) int {
	if max <= min {
		return min
	}
	return min + rand.Intn(max-min+1)
}

// w10Ints builds the answer and fact map from named operands.
func w10Problem(kp int, question string, answer int, explanation string, pairs ...interface{}) generator.Problem {
	facts := map[string]string{}
	for i := 0; i+1 < len(pairs); i += 2 {
		key, ok := pairs[i].(string)
		if !ok {
			continue
		}
		switch v := pairs[i+1].(type) {
		case int:
			facts[key] = strconv.Itoa(v)
		case string:
			facts[key] = v
		}
	}
	facts["answer"] = strconv.Itoa(answer)
	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       facts,
		KPOrdinal:   kp,
	}
}
