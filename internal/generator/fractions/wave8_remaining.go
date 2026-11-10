package fractions

// Wave 8 fractions remaining concepts. Every answer is an integer and every
// operand is published in Facts for independent recomputation.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave8Remaining(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"frac.basics.improper_number_line": "improper_number_line",
		"frac.order.three":                 "order_three",
		"frac.add.number_line":             "add_number_line",
		"frac.sub.number_line":             "sub_number_line",
		"frac.mult.area_model":             "mult_area_model",
		"frac.div.missing_factor":          "missing_factor",
		"frac.mixed.to_decimal":            "mixed_decimal",
		"frac.percent.convert":             "percent_convert",
		"frac.word.remaining_two_step":     "remaining_two_step",
	} {
		reg.Register(id, &w8FracRemainingGen{kind: kind})
	}
}

type w8FracRemainingGen struct{ kind string }

func (g *w8FracRemainingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "improper_number_line":
		return w8ImproperNumberLine(kp)
	case "order_three":
		return w8OrderThree(kp)
	case "add_number_line":
		return w8AddNumberLine(kp)
	case "sub_number_line":
		return w8SubNumberLine(kp)
	case "mult_area_model":
		return w8MultAreaModel(kp)
	case "missing_factor":
		return w8MissingFactor(kp)
	case "mixed_decimal":
		return w8MixedDecimal(kp)
	case "percent_convert":
		return w8PercentConvert(kp)
	case "remaining_two_step":
		return w8RemainingTwoStep(kp)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

func w8RP(o int, q, a, ex string, f map[string]string) generator.Problem {
	f["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: f, KPOrdinal: o}
}

func w8RPI(f map[string]string, k string, v int) { f[k] = strconv.Itoa(v) }

func w8RPR(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func w8RPGCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func w8ImproperNumberLine(kp int) generator.Problem {
	b := w8RPR(2, 9)
	whole := w8RPR(1, 5)
	rem := w8RPR(1, b-1)
	a := whole*b + rem
	f := map[string]string{}
	w8RPI(f, "a", a)
	w8RPI(f, "b", b)
	w8RPI(f, "whole", whole)
	w8RPI(f, "remainder", rem)
	switch kp {
	case 1:
		return w8RP(kp, fmt.Sprintf("Write \\(%d/%d\\) as a mixed number. What is the whole-number part? (enter a number)", a, b), strconv.Itoa(whole),
			fmt.Sprintf("\\(%d=%d\\times%d+%d\\), so the whole part is %d.", a, whole, b, rem, whole), f)
	case 2:
		return w8RP(kp, fmt.Sprintf("Write \\(%d/%d\\) as a mixed number. What is the numerator of the fractional part? (enter a number)", a, b), strconv.Itoa(rem),
			fmt.Sprintf("\\(%d=%d\\times%d+%d\\), so the fractional numerator is %d.", a, whole, b, rem, rem), f)
	default:
		return w8RP(kp, fmt.Sprintf("The point \\(%d/%d\\) lies between which two consecutive whole numbers? Enter the lower whole number.", a, b), strconv.Itoa(whole),
			fmt.Sprintf("\\(%d/%d\\) is between %d and %d, so the lower whole number is %d.", a, b, whole, whole+1, whole), f)
	}
}

func w8OrderThree(kp int) generator.Problem {
	for {
		b1 := w8RPR(2, 6)
		b2 := w8RPR(2, 6)
		b3 := w8RPR(2, 6)
		a1 := w8RPR(1, b1-1)
		a2 := w8RPR(1, b2-1)
		a3 := w8RPR(1, b3-1)
		v1, v2, v3 := a1*b2*b3, a2*b1*b3, a3*b1*b2
		if v1 == v2 || v2 == v3 || v1 == v3 {
			continue
		}
		small, large := 1, 1
		if v2 < v1 && v2 < v3 {
			small = 2
		} else if v3 < v1 && v3 < v2 {
			small = 3
		}
		if v2 > v1 && v2 > v3 {
			large = 2
		} else if v3 > v1 && v3 > v2 {
			large = 3
		}
		lcm := b1 * b2 / w8RPGCD(b1, b2)
		lcm = lcm * b3 / w8RPGCD(lcm, b3)
		f := map[string]string{}
		w8RPI(f, "a1", a1)
		w8RPI(f, "b1", b1)
		w8RPI(f, "a2", a2)
		w8RPI(f, "b2", b2)
		w8RPI(f, "a3", a3)
		w8RPI(f, "b3", b3)
		w8RPI(f, "lcm", lcm)
		switch kp {
		case 1:
			return w8RP(kp, fmt.Sprintf("Which fraction is smallest: \\(%d/%d\\), \\(%d/%d\\) or \\(%d/%d\\)? Enter 1, 2 or 3.", a1, b1, a2, b2, a3, b3), strconv.Itoa(small),
				fmt.Sprintf("Compare with a common denominator or cross products. The smallest is fraction %d.", small), f)
		case 2:
			return w8RP(kp, fmt.Sprintf("Which fraction is largest: \\(%d/%d\\), \\(%d/%d\\) or \\(%d/%d\\)? Enter 1, 2 or 3.", a1, b1, a2, b2, a3, b3), strconv.Itoa(large),
				fmt.Sprintf("The largest is fraction %d.", large), f)
		default:
			return w8RP(kp, fmt.Sprintf("What is the least common denominator of %d, %d and %d? (enter a number)", b1, b2, b3), strconv.Itoa(lcm),
				fmt.Sprintf("The least common multiple of %d, %d and %d is %d.", b1, b2, b3, lcm), f)
		}
	}
}

func w8AddNumberLine(kp int) generator.Problem {
	b := w8RPR(3, 12)
	a := w8RPR(1, b-1)
	c := w8RPR(1, b-1)
	num := a + c
	g := w8RPGCD(num, b)
	f := map[string]string{}
	w8RPI(f, "a", a)
	w8RPI(f, "b", b)
	w8RPI(f, "c", c)
	w8RPI(f, "gcd", g)
	switch kp {
	case 1:
		return w8RP(kp, fmt.Sprintf("Add \\(%d/%d+%d/%d\\) using the number line. What is the numerator before simplifying? (enter a number)", a, b, c, b), strconv.Itoa(num),
			fmt.Sprintf("The denominators match, so add the numerators: \\(%d+%d=%d\\).", a, c, num), f)
	case 2:
		return w8RP(kp, fmt.Sprintf("For \\(%d/%d+%d/%d\\), what is the common denominator? (enter a number)", a, b, c, b), strconv.Itoa(b),
			fmt.Sprintf("Both fractions already have denominator %d.", b), f)
	default:
		if g == 1 {
			return w8RP(kp, fmt.Sprintf("Simplify \\(%d/%d\\). What is the numerator? (enter a number)", num, b), strconv.Itoa(num),
				fmt.Sprintf("The fraction \\(%d/%d\\) is already in lowest terms.", num, b), f)
		}
		if g == 1 {
			return w8RP(kp, fmt.Sprintf("Simplify \\(%d/%d\\). What is the numerator? (enter a number)", num, b), strconv.Itoa(num),
				fmt.Sprintf("The fraction \\(%d/%d\\) is already in lowest terms.", num, b), f)
		}
		return w8RP(kp, fmt.Sprintf("Simplify \\(%d/%d\\). What is the numerator? (enter a number)", num, b), strconv.Itoa(num/g),
			fmt.Sprintf("Divide numerator and denominator by %d: \\(%d/%d\\).", g, num/g, b/g), f)
	}
}

func w8SubNumberLine(kp int) generator.Problem {
	b := w8RPR(3, 12)
	a := w8RPR(2, b)
	c := w8RPR(1, a-1)
	num := a - c
	g := w8RPGCD(num, b)
	f := map[string]string{}
	w8RPI(f, "a", a)
	w8RPI(f, "b", b)
	w8RPI(f, "c", c)
	w8RPI(f, "gcd", g)
	switch kp {
	case 1:
		return w8RP(kp, fmt.Sprintf("Subtract \\(%d/%d-%d/%d\\) using the number line. What is the numerator before simplifying? (enter a number)", a, b, c, b), strconv.Itoa(num),
			fmt.Sprintf("The denominators match, so subtract the numerators: \\(%d-%d=%d\\).", a, c, num), f)
	case 2:
		return w8RP(kp, fmt.Sprintf("For \\(%d/%d-%d/%d\\), what is the common denominator? (enter a number)", a, b, c, b), strconv.Itoa(b),
			fmt.Sprintf("Both fractions already have denominator %d.", b), f)
	default:
		if g == 1 {
			return w8RP(kp, fmt.Sprintf("Simplify \\(%d/%d\\). What is the numerator? (enter a number)", num, b), strconv.Itoa(num),
				fmt.Sprintf("The fraction \\(%d/%d\\) is already in lowest terms.", num, b), f)
		}
		if g == 1 {
			return w8RP(kp, fmt.Sprintf("Simplify \\(%d/%d\\). What is the numerator? (enter a number)", num, b), strconv.Itoa(num),
				fmt.Sprintf("The fraction \\(%d/%d\\) is already in lowest terms.", num, b), f)
		}
		return w8RP(kp, fmt.Sprintf("Simplify \\(%d/%d\\). What is the numerator? (enter a number)", num, b), strconv.Itoa(num/g),
			fmt.Sprintf("Divide numerator and denominator by %d: \\(%d/%d\\).", g, num/g, b/g), f)
	}
}

func w8MultAreaModel(kp int) generator.Problem {
	b := w8RPR(2, 9)
	d := w8RPR(2, 9)
	a := w8RPR(1, b-1)
	c := w8RPR(1, d-1)
	num := a * c
	den := b * d
	g := w8RPGCD(num, den)
	f := map[string]string{}
	w8RPI(f, "a", a)
	w8RPI(f, "b", b)
	w8RPI(f, "c", c)
	w8RPI(f, "d", d)
	w8RPI(f, "gcd", g)
	switch kp {
	case 1:
		return w8RP(kp, fmt.Sprintf("Use an area model to multiply \\(%d/%d\\times%d/%d\\). What is the numerator before simplifying? (enter a number)", a, b, c, d), strconv.Itoa(num),
			fmt.Sprintf("The rectangle has \\(%d\\times%d=%d\\) shaded cells in the numerator.", a, c, num), f)
	case 2:
		return w8RP(kp, fmt.Sprintf("For \\(%d/%d\\times%d/%d\\), what is the denominator before simplifying? (enter a number)", a, b, c, d), strconv.Itoa(den),
			fmt.Sprintf("The grid has \\(%d\\times%d=%d\\) cells in the denominator.", b, d, den), f)
	default:
		if g == 1 {
			return w8RP(kp, fmt.Sprintf("Simplify the product \\(%d/%d\\). What is the numerator? (enter a number)", num, den), strconv.Itoa(num),
				fmt.Sprintf("The product \\(%d/%d\\) is already in lowest terms.", num, den), f)
		}
		return w8RP(kp, fmt.Sprintf("Simplify the product \\(%d/%d\\). What is the numerator? (enter a number)", num, den), strconv.Itoa(num/g),
			fmt.Sprintf("Divide numerator and denominator by %d: \\(%d/%d\\).", g, num/g, den/g), f)
	}
}

func w8MissingFactor(kp int) generator.Problem {
	b := w8RPR(2, 9)
	d := w8RPR(2, 9)
	a := w8RPR(1, 12)
	c := w8RPR(1, d-1)
	// x * (c/d) = a/b, so x = a*d/(b*c).
	num := a * d
	den := b * c
	g := w8RPGCD(num, den)
	f := map[string]string{}
	w8RPI(f, "a", a)
	w8RPI(f, "b", b)
	w8RPI(f, "c", c)
	w8RPI(f, "d", d)
	w8RPI(f, "gcd", g)
	switch kp {
	case 1:
		return w8RP(kp, fmt.Sprintf("Solve \\(x\\times%d/%d=%d/%d\\). What is the numerator of \\(x\\) in lowest terms? (enter a number)", c, d, a, b), strconv.Itoa(num/g),
			fmt.Sprintf("\\(x=(%d\\times%d)/(%d\\times%d)=%d/%d\\), %s.", a, d, b, c, num, den, map[bool]string{true: "which is already in lowest terms", false: fmt.Sprintf("which simplifies by %d", g)}[g == 1]), f)
	case 2:
		return w8RP(kp, fmt.Sprintf("Solve \\(x\\times%d/%d=%d/%d\\). What is the denominator of \\(x\\) in lowest terms? (enter a number)", c, d, a, b), strconv.Itoa(den/g),
			fmt.Sprintf("\\(x=(%d\\times%d)/(%d\\times%d)=%d/%d\\), %s.", a, d, b, c, num, den, map[bool]string{true: "which is already in lowest terms", false: fmt.Sprintf("which simplifies by %d", g)}[g == 1]), f)
	default:
		return w8RP(kp, fmt.Sprintf("In \\(x\\times%d/%d=%d/%d\\), what is the numerator of \\(x\\) before simplifying? (enter a number)", c, d, a, b), strconv.Itoa(num),
			fmt.Sprintf("Multiply both sides by \\(%d/%d\\): the numerator is \\(%d\\times%d=%d\\).", d, c, a, d, num), f)
	}
}

func w8MixedDecimal(kp int) generator.Problem {
	w := w8RPR(1, 9)
	n := w8RPR(1, 9)
	f := map[string]string{}
	w8RPI(f, "whole", w)
	w8RPI(f, "numerator", n)
	w8RPI(f, "denominator", 10)
	switch kp {
	case 1:
		return w8RP(kp, fmt.Sprintf("Write \\(%d\\,%d/10\\) as a decimal. What is its value in hundredths? (enter a number)", w, n), strconv.Itoa(w*100+n*10),
			fmt.Sprintf("\\(%d\\,%d/10=%d.%d=%d\\) hundredths.", w, n, w, n, w*100+n*10), f)
	case 2:
		return w8RP(kp, fmt.Sprintf("Write \\(%d\\,%d/10\\) as a decimal. What is its value in tenths? (enter a number)", w, n), strconv.Itoa(w*10+n),
			fmt.Sprintf("\\(%d\\,%d/10=%d\\) tenths.", w, n, w*10+n), f)
	default:
		return w8RP(kp, fmt.Sprintf("Is the decimal %d.%d greater than %d? Enter 1 for yes and 0 for no.", w, n, w), "1",
			fmt.Sprintf("The decimal has whole part %d and a positive tenths digit, so it is greater than %d.", w, w), f)
	}
}

func w8PercentConvert(kp int) generator.Problem {
	den := []int{2, 4, 5, 10, 20, 25, 50}[w8RPR(0, 6)]
	a := w8RPR(1, den-1)
	pct := a * 100 / den
	over100 := a * 100 / den
	f := map[string]string{}
	w8RPI(f, "a", a)
	w8RPI(f, "b", den)
	w8RPI(f, "percent", pct)
	switch kp {
	case 1:
		return w8RP(kp, fmt.Sprintf("Write \\(%d/%d\\) as a percent. What is the percent value? (enter a number)", a, den), strconv.Itoa(pct),
			fmt.Sprintf("\\(%d/%d=%d/100=%d\\%%\\).", a, den, over100, pct), f)
	case 2:
		return w8RP(kp, fmt.Sprintf("Write \\(%d/%d\\) with denominator 100. What is the numerator? (enter a number)", a, den), strconv.Itoa(over100),
			fmt.Sprintf("Multiply numerator and denominator by %d to reach denominator 100.", 100/den), f)
	default:
		return w8RP(kp, fmt.Sprintf("Is \\(%d/%d\\) greater than 50 percent? Enter 1 for yes and 0 for no.", a, den), map[bool]string{true: "1", false: "0"}[pct > 50],
			fmt.Sprintf("The fraction is %d percent, so the answer is %s.", pct, map[bool]string{true: "1", false: "0"}[pct > 50]), f)
	}
}

func w8RemainingTwoStep(kp int) generator.Problem {
	den1 := w8RPR(2, 6)
	den2 := w8RPR(2, 6)
	num1 := w8RPR(1, den1-1)
	num2 := w8RPR(1, den2-1)
	n := den1 * den2 * w8RPR(2, 5)
	afterFirst := n * (den1 - num1) / den1
	afterSecond := afterFirst * (den2 - num2) / den2
	f := map[string]string{}
	w8RPI(f, "n", n)
	w8RPI(f, "num1", num1)
	w8RPI(f, "den1", den1)
	w8RPI(f, "num2", num2)
	w8RPI(f, "den2", den2)
	switch kp {
	case 1:
		return w8RP(kp, fmt.Sprintf("A jar has %d candies. \\(%d/%d\\) are eaten. How many remain? (enter a number)", n, num1, den1), strconv.Itoa(afterFirst),
			fmt.Sprintf("Eaten: \\(%d\\times%d/%d=%d\\). Remaining: \\(%d-%d=%d\\).", n, num1, den1, n*num1/den1, n, n*num1/den1, afterFirst), f)
	case 2:
		return w8RP(kp, fmt.Sprintf("After \\(%d/%d\\) of %d is eaten, \\(%d/%d\\) of the remainder is eaten. How many candies remain? (enter a number)", num1, den1, n, num2, den2), strconv.Itoa(afterSecond),
			fmt.Sprintf("After the first step %d remain. Then \\(%d\\times%d/%d=%d\\) are eaten, leaving %d.", afterFirst, afterFirst, num2, den2, afterFirst*num2/den2, afterSecond), f)
	default:
		return w8RP(kp, fmt.Sprintf("A jar has %d candies. \\(%d/%d\\) are eaten, then \\(%d/%d\\) of the remainder is eaten. What is the numerator of the remaining fraction in lowest terms? (enter a number)", n, num1, den1, num2, den2), strconv.Itoa(afterSecond/w8RPGCD(afterSecond, n)),
			fmt.Sprintf("The remaining fraction is \\(%d/%d\\), which reduces to \\(%d/%d\\).", afterSecond, n, afterSecond/w8RPGCD(afterSecond, n), n/w8RPGCD(afterSecond, n)), f)
	}
}
