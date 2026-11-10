package fractions

// Wave 8 fractions: the first three review concepts. Every answer is an
// integer and every operand is published in Facts so wave8_test.go can
// recompute it independently.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave8(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"frac.basics.set":             "set",
		"frac.compare.same_numerator": "same_numerator",
		"frac.div.common_denominator": "common_denominator",
	} {
		reg.Register(id, &w8FracGen{kind: kind})
	}
}

type w8FracGen struct{ kind string }

func (g *w8FracGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "set":
		return w8FractionSet(kp)
	case "same_numerator":
		return w8SameNumerator(kp)
	case "common_denominator":
		return w8CommonDenominator(kp)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

func w8F(o int, q, a, ex string, f map[string]string) generator.Problem {
	f["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: f, KPOrdinal: o}
}

func w8FI(f map[string]string, k string, v int) { f[k] = strconv.Itoa(v) }

func w8FR(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func w8FGCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func w8FractionSet(kp int) generator.Problem {
	den := w8FR(3, 8)
	num := w8FR(1, den-1)
	groups := w8FR(2, 9)
	n := den * groups
	f := map[string]string{}
	w8FI(f, "n", n)
	w8FI(f, "num", num)
	w8FI(f, "den", den)
	switch kp {
	case 1:
		ans := groups * num
		return w8F(kp, fmt.Sprintf("A set has %d objects. \\(%d/%d\\) of them are red. How many are red? (enter a number)", n, num, den), strconv.Itoa(ans),
			fmt.Sprintf("The set divides into %d groups of %d. The fraction selects %d of those groups, so %d objects are red.", den, groups, num, ans), f)
	case 2:
		red := groups * num
		return w8F(kp, fmt.Sprintf("A set has %d objects and \\(%d/%d\\) are red. How many are not red? (enter a number)", n, num, den), strconv.Itoa(n-red),
			fmt.Sprintf("The red part has %d objects, so the rest is \\(%d-%d=%d\\).", red, n, red, n-red), f)
	default:
		return w8F(kp, fmt.Sprintf("A set of %d objects is split into %d equal groups. How many objects are in each group? (enter a number)", n, den), strconv.Itoa(groups),
			fmt.Sprintf("\\(%d/%d=%d\\), so each group contains %d objects.", n, den, groups, groups), f)
	}
}

func w8SameNumerator(kp int) generator.Problem {
	a := w8FR(1, 9)
	b := w8FR(2, 12)
	c := w8FR(2, 12)
	for c == b {
		c = w8FR(2, 12)
	}
	larger := b
	if c < b {
		larger = c
	}
	f := map[string]string{}
	w8FI(f, "a", a)
	w8FI(f, "b", b)
	w8FI(f, "c", c)
	switch kp {
	case 1:
		code := 1
		if c < b {
			code = 2
		}
		return w8F(kp, fmt.Sprintf("Which fraction is larger: \\(%d/%d\\) or \\(%d/%d\\)? Enter 1 for the first and 2 for the second.", a, b, a, c), strconv.Itoa(code),
			fmt.Sprintf("With the same numerator, the fraction with the smaller denominator is larger. The smaller denominator is %d, so the answer is %d.", larger, code), f)
	case 2:
		return w8F(kp, fmt.Sprintf("For \\(%d/%d\\) and \\(%d/%d\\), what is the denominator of the larger fraction? (enter a number)", a, b, a, c), strconv.Itoa(larger),
			fmt.Sprintf("The denominators are %d and %d. The smaller one, %d, gives the larger fraction.", b, c, larger), f)
	default:
		ans := 0
		if a < b {
			ans = 1
		}
		return w8F(kp, fmt.Sprintf("Is \\(%d/%d\\) less than 1? Enter 1 for yes and 0 for no.", a, b), strconv.Itoa(ans),
			fmt.Sprintf("A fraction is less than 1 exactly when its numerator is smaller than its denominator. Here %d is %s than %d, so the answer is %d.", a, map[bool]string{true: "smaller", false: "not smaller"}[a < b], b, ans), f)
	}
}

func w8CommonDenominator(kp int) generator.Problem {
	a := w8FR(1, 9)
	b := w8FR(1, 9)
	c := w8FR(2, 12)
	g := w8FGCD(a, b)
	f := map[string]string{}
	w8FI(f, "a", a)
	w8FI(f, "b", b)
	w8FI(f, "c", c)
	w8FI(f, "gcd", g)
	simplify := fmt.Sprintf("The quotient is \\(%d/%d\\). Divide both parts by %d to get \\(%d/%d\\).", a, b, g, a/g, b/g)
	if g == 1 {
		simplify = fmt.Sprintf("The quotient is \\(%d/%d\\), which is already in lowest terms.", a, b)
	}
	switch kp {
	case 1:
		return w8F(kp, fmt.Sprintf("Compute \\(( %d/%d ) / ( %d/%d )\\) and simplify. What is the numerator of the simplified quotient? (enter a number)", a, c, b, c), strconv.Itoa(a/g),
			fmt.Sprintf("The common denominator cancels. %s The numerator is %d.", simplify, a/g), f)
	case 2:
		return w8F(kp, fmt.Sprintf("Compute \\(( %d/%d ) / ( %d/%d )\\) and simplify. What is the denominator of the simplified quotient? (enter a number)", a, c, b, c), strconv.Itoa(b/g),
			fmt.Sprintf("The common denominator cancels. %s The denominator is %d.", simplify, b/g), f)
	default:
		return w8F(kp, fmt.Sprintf("What is the greatest common divisor of %d and %d? (enter a number)", a, b), strconv.Itoa(g),
			fmt.Sprintf("The largest integer dividing both %d and %d is %d, so the quotient simplifies by a factor of %d.", a, b, g, g), f)
	}
}
