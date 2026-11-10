package fractions

// Wave 9 remaining fractions: ordering mixed numbers and a percentage of a
// fraction. Both publish the operands used by the integer answer.

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

var w9rPureNumber = regexp.MustCompile(`\\\((-?\d+)\\\)`)
var w9rUnitCoefficient = regexp.MustCompile(`(^|[^0-9])1 ?([xyzuvtr])([^A-Za-z]|$)`)
var w9rUnitParen = regexp.MustCompile(`(^|[^0-9])1\(`)
var w9rSpaces = regexp.MustCompile(`[ \t]+`)

// w9rClean removes rendering defects that a learner should not have to read.
func w9rClean(s string) string {
	s = w9rPureNumber.ReplaceAllString(s, "$1")
	s = strings.ReplaceAll(s, "+-", "-")
	s = strings.ReplaceAll(s, "--", "+")
	s = w9rSpaces.ReplaceAllString(s, " ")
	s = w9rUnitCoefficient.ReplaceAllString(s, "${1}${2}${3}")
	s = w9rUnitParen.ReplaceAllString(s, "${1}(")
	return s
}

func registerWave9Remaining(reg *generator.Registry) {
	reg.Register("frac.compare.mixed_order", &w9rMixedOrderGen{})
	reg.Register("frac.percent.of_fraction", &w9rPercentOfFractionGen{})
}

type w9rMixedOrderGen struct{}

func (g *w9rMixedOrderGen) Generate(ctx generator.GeneratorContext) (prob generator.Problem) {
	defer func() {
		prob.Question = w9rClean(prob.Question)
		prob.Explanation = w9rClean(prob.Explanation)
	}()
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	var w1, w2, b, d, a, c, left, right int
	for {
		w1 = w9rFR(1, 6)
		w2 = w9rFR(1, 6)
		b = w9rFR(2, 7)
		d = w9rFR(2, 7)
		a = w9rFR(1, b-1)
		c = w9rFR(1, d-1)
		// Compare w1+a/b and w2+c/d by cross multiplication over b*d.
		left = (w1*b + a) * d
		right = (w2*d + c) * b
		if left != right {
			break
		}
	}
	f := map[string]string{}
	w9rFI(f, "w1", w1)
	w9rFI(f, "a", a)
	w9rFI(f, "b", b)
	w9rFI(f, "w2", w2)
	w9rFI(f, "c", c)
	w9rFI(f, "d", d)
	firstLarger := left > right
	switch kp {
	case 1:
		code := 2
		if firstLarger {
			code = 1
		}
		relation := "the second is larger"
		if firstLarger {
			relation = "the first is larger"
		}
		return w9rFP(kp, fmt.Sprintf("Which mixed number is larger: \\(%d\\frac{%d}{%d}\\) or \\(%d\\frac{%d}{%d}\\)? Enter 1 for the first and 2 for the second.", w1, a, b, w2, c, d), strconv.Itoa(code),
			fmt.Sprintf("Cross multiplication gives \\(%d\\) versus \\(%d\\), so %s.", left, right, relation), f)
	case 2:
		whole := w2
		if firstLarger {
			whole = w1
		}
		return w9rFP(kp, fmt.Sprintf("Between \\(%d\\frac{%d}{%d}\\) and \\(%d\\frac{%d}{%d}\\), what is the whole-number part of the larger mixed number? (enter a number)", w1, a, b, w2, c, d), strconv.Itoa(whole),
			fmt.Sprintf("The larger mixed number has whole part %d.", whole), f)
	default:
		code := 0
		if a*d > c*b {
			code = 1
		}
		return w9rFP(kp, fmt.Sprintf("Compare the fractional parts \\(%d/%d\\) and \\(%d/%d\\). Is the first fractional part larger? Enter 1 for yes and 0 for no.", a, b, c, d), strconv.Itoa(code),
			fmt.Sprintf("Cross multiplication gives \\(%d\\) versus \\(%d\\).", a*d, c*b), f)
	}
}

type w9rPercentOfFractionGen struct{}

func (g *w9rPercentOfFractionGen) Generate(ctx generator.GeneratorContext) (prob generator.Problem) {
	defer func() {
		prob.Question = w9rClean(prob.Question)
		prob.Explanation = w9rClean(prob.Explanation)
	}()
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	p := w9rFR(1, 12)
	a := w9rFR(1, 8)
	b := w9rFR(2, 9)
	f := map[string]string{}
	w9rFI(f, "p", p)
	w9rFI(f, "a", a)
	w9rFI(f, "b", b)
	num := p * a
	den := 100 * b
	switch kp {
	case 1:
		return w9rFP(kp, fmt.Sprintf("What is \\(%d\\%%\\) of \\(%d/%d\\)? Give the numerator before simplification over \\(100\\cdot%d\\). (enter a number)", p, a, b, b), strconv.Itoa(num),
			fmt.Sprintf("The unsimplified product is \\(%d\\cdot%d/(100\\cdot%d)=%d/%d\\).", p, a, b, num, den), f)
	case 2:
		return w9rFP(kp, fmt.Sprintf("What is the denominator before simplification of \\(%d\\%%\\) of \\(%d/%d\\)? (enter a number)", p, a, b), strconv.Itoa(den),
			fmt.Sprintf("The denominator is \\(100\\cdot%d=%d\\).", b, den), f)
	default:
		code := 0
		if 2*num > den {
			code = 1
		}
		return w9rFP(kp, fmt.Sprintf("Is \\(%d\\%%\\) of \\(%d/%d\\) greater than one half? Enter 1 for yes and 0 for no.", p, a, b), strconv.Itoa(code),
			fmt.Sprintf("Compare \\(2\\cdot%d=%d\\) with \\(%d\\).", num, 2*num, den), f)
	}
}

func w9rFR(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func w9rFI(m map[string]string, k string, v int) { m[k] = strconv.Itoa(v) }

func w9rFP(o int, q, a, ex string, facts map[string]string) generator.Problem {
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}
