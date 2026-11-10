package prealgebra

// Wave 10 prealgebra: three-term ratios, part-to-whole ratios, recovery of an
// original amount after a percent change, and literal equations. Every branch
// publishes the same fact keys so one authored schema resolves for all three
// knowledge points.

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/grader"
)

func registerWave10(reg *generator.Registry) {
	reg.Register("prealg.ratio.three_term", &w10ThreeTermRatioGen{})
	reg.Register("prealg.ratio.part_whole", &w10PartWholeGen{})
	reg.Register("pct.ops.original", &w10OriginalAmountGen{})
	reg.Register("prealg.eq.literal", &w10LiteralEquationGen{})
	registerWave10Rest(reg)
}

func w10P(o int, q, a, ex string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}

func w10I(m map[string]string, k string, v int) {
	m[k] = strconv.Itoa(v)
}

func w10Ord(ctx generator.GeneratorContext) int {
	o := ctx.KPOrdinal
	if o < 1 || o > 3 {
		o = 1 + rand.Intn(3)
	}
	return o
}

func w10GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func w10PartWord(n int) string {
	if n == 1 {
		return "part"
	}
	return "parts"
}

func w10NormaliseTriple(s string) (int, int, int, bool) {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " to ", ":")
	s = strings.ReplaceAll(s, "/", ":")
	s = strings.ReplaceAll(s, "÷", ":")
	parts := strings.Split(s, ":")
	if len(parts) != 3 {
		return 0, 0, 0, false
	}
	a, errA := strconv.Atoi(strings.TrimSpace(parts[0]))
	b, errB := strconv.Atoi(strings.TrimSpace(parts[1]))
	c, errC := strconv.Atoi(strings.TrimSpace(parts[2]))
	if errA != nil || errB != nil || errC != nil || b == 0 || c == 0 {
		return 0, 0, 0, false
	}
	return a, b, c, true
}

func w10NumericEqual(expected, answer string) bool {
	e, errE := strconv.ParseFloat(strings.TrimSpace(expected), 64)
	a, errA := strconv.ParseFloat(strings.TrimSpace(answer), 64)
	return errE == nil && errA == nil && e == a
}

// ---------------------------------------------------------------------------
// 1. Three-term ratios
// ---------------------------------------------------------------------------

type w10ThreeTermRatioGen struct{}

func (g *w10ThreeTermRatioGen) AnswerFormat() generator.AnswerFormat {
	return generator.AnswerFormat{Hint: "Ratio as a:b:c, e.g. 2:3:4", InputMode: "text"}
}

func (g *w10ThreeTermRatioGen) Grade(expected, userAnswer string) grader.Result {
	if strings.Contains(expected, ":") {
		e1, e2, e3, eOK := w10NormaliseTriple(expected)
		a1, a2, a3, aOK := w10NormaliseTriple(userAnswer)
		if !eOK || !aOK {
			return grader.Result{Correct: false, Score: 0, Feedback: "Enter the ratio as a:b:c"}
		}
		if e1 == a1 && e2 == a2 && e3 == a3 {
			return grader.Result{Correct: true, Score: 1}
		}
		return grader.Result{Correct: false, Score: 0, Feedback: "Incorrect ratio"}
	}
	if w10NumericEqual(expected, userAnswer) {
		return grader.Result{Correct: true, Score: 1}
	}
	return grader.Result{Correct: false, Score: 0}
}

func (g *w10ThreeTermRatioGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10Ord(ctx)
	for {
		r1 := 1 + rand.Intn(5)
		r2 := 1 + rand.Intn(5)
		r3 := 1 + rand.Intn(5)
		if r1 == r2 && r2 == r3 {
			continue
		}
		if w10GCD(w10GCD(r1, r2), r3) != 1 {
			continue
		}
		scale := 2 + rand.Intn(5)
		a, b, c := r1*scale, r2*scale, r3*scale
		total := a + b + c
		facts := map[string]string{}
		w10I(facts, "r1", r1)
		w10I(facts, "r2", r2)
		w10I(facts, "r3", r3)
		w10I(facts, "scale", scale)
		w10I(facts, "total", total)
		switch o {
		case 1:
			ans := fmt.Sprintf("%d:%d:%d", r1, r2, r3)
			return w10P(o,
				fmt.Sprintf("A bag has %d red, %d blue and %d green counters. Write the ratio red:blue:green in simplest form.", a, b, c),
				ans,
				fmt.Sprintf("The common factor is %d, so divide every count by %d: %d:%d:%d.", scale, scale, r1, r2, r3),
				facts)
		case 2:
			ans := fmt.Sprintf("%d:%d:%d", r1, r2, r3)
			return w10P(o,
				fmt.Sprintf("Simplify the three-term ratio %d:%d:%d.", a, b, c),
				ans,
				fmt.Sprintf("Divide all three terms by their common factor %d to get %d:%d:%d.", scale, r1, r2, r3),
				facts)
		default:
			unit := 1 + rand.Intn(10)
			total = (r1 + r2 + r3) * unit
			largest := r3
			if r1 > largest {
				largest = r1
			}
			if r2 > largest {
				largest = r2
			}
			ans := largest * unit
			w10I(facts, "scale", unit)
			w10I(facts, "total", total)
			return w10P(o,
				fmt.Sprintf("Three quantities are in the ratio %d:%d:%d and their total is %d. What is the largest quantity?", r1, r2, r3, total),
				strconv.Itoa(ans),
				fmt.Sprintf("One part is %d / (%d+%d+%d) = %d. The largest share is %d parts, or %d.", total, r1, r2, r3, unit, largest, ans),
				facts)
		}
	}
}

// ---------------------------------------------------------------------------
// 2. Part-to-whole ratios
// ---------------------------------------------------------------------------

type w10PartWholeGen struct{}

func (g *w10PartWholeGen) AnswerFormat() generator.AnswerFormat {
	return generator.AnswerFormat{Hint: "Ratio as a:b, or a number when asked", InputMode: "text"}
}

func (g *w10PartWholeGen) Grade(expected, userAnswer string) grader.Result {
	if strings.Contains(expected, ":") {
		e1, e2, eOK := w10NormalisePair(expected)
		a1, a2, aOK := w10NormalisePair(userAnswer)
		if !eOK || !aOK {
			return grader.Result{Correct: false, Score: 0, Feedback: "Enter the ratio as a:b"}
		}
		if e1 == a1 && e2 == a2 {
			return grader.Result{Correct: true, Score: 1}
		}
		return grader.Result{Correct: false, Score: 0, Feedback: "Incorrect ratio"}
	}
	if w10NumericEqual(expected, userAnswer) {
		return grader.Result{Correct: true, Score: 1}
	}
	return grader.Result{Correct: false, Score: 0}
}

func w10NormalisePair(s string) (int, int, bool) {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " to ", ":")
	s = strings.ReplaceAll(s, "/", ":")
	s = strings.ReplaceAll(s, "÷", ":")
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return 0, 0, false
	}
	a, errA := strconv.Atoi(strings.TrimSpace(parts[0]))
	b, errB := strconv.Atoi(strings.TrimSpace(parts[1]))
	if errA != nil || errB != nil || b == 0 {
		return 0, 0, false
	}
	return a, b, true
}

func (g *w10PartWholeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10Ord(ctx)
	for {
		p := 1 + rand.Intn(5)
		q := 1 + rand.Intn(5)
		if p == q || w10GCD(p, q) != 1 {
			continue
		}
		partRatio, wholeRatio := p, p+q
		if w10GCD(partRatio, wholeRatio) != 1 {
			continue
		}
		scale := 2 + rand.Intn(5)
		part := partRatio * scale
		whole := wholeRatio * scale
		other := q * scale
		facts := map[string]string{}
		w10I(facts, "ratio_part", partRatio)
		w10I(facts, "ratio_whole", wholeRatio)
		w10I(facts, "part", part)
		w10I(facts, "whole", whole)
		w10I(facts, "scale", scale)
		switch o {
		case 1:
			ans := fmt.Sprintf("%d:%d", partRatio, wholeRatio)
			return w10P(o,
				fmt.Sprintf("A class has %d boys and %d girls. Write the ratio of boys to all students in simplest form.", part, other),
				ans,
				fmt.Sprintf("There are %d + %d = %d students. The ratio %d:%d simplifies by %d to %d:%d.", part, other, whole, part, whole, scale, partRatio, wholeRatio),
				facts)
		case 2:
			return w10P(o,
				fmt.Sprintf("The ratio of a part to the whole is %d:%d. The part is %d. What is the whole?", partRatio, wholeRatio, part),
				strconv.Itoa(whole),
				fmt.Sprintf("The scale factor is %d / %d = %d, so the whole is %d x %d = %d.", part, partRatio, scale, wholeRatio, scale, whole),
				facts)
		default:
			return w10P(o,
				fmt.Sprintf("The ratio of red counters to all counters is %d:%d. There are %d blue counters and no other colours. How many red counters are there?", partRatio, wholeRatio, other),
				strconv.Itoa(part),
				fmt.Sprintf("Blue is %d - %d = %d %s, so one ratio part is %d / %d = %d counters. Red is %d %s, or %d.", wholeRatio, partRatio, q, w10PartWord(q), other, q, scale, partRatio, w10PartWord(partRatio), part),
				facts)
		}
	}
}

// ---------------------------------------------------------------------------
// 3. Finding the original amount
// ---------------------------------------------------------------------------

type w10OriginalAmountGen struct{}

func (g *w10OriginalAmountGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10Ord(ctx)
	original := (1 + rand.Intn(10)) * 100
	percent := (1 + rand.Intn(10)) * 5
	change := original * percent / 100
	direction := "increase"
	newValue := original + change
	if o == 2 {
		direction = "decrease"
		newValue = original - change
	}
	facts := map[string]string{}
	w10I(facts, "percent", percent)
	w10I(facts, "original", original)
	w10I(facts, "new", newValue)
	w10I(facts, "change", change)
	facts["direction"] = direction
	switch o {
	case 1:
		return w10P(o,
			fmt.Sprintf("After a %d%% increase, a price is %d. What was the original price?", percent, newValue),
			strconv.Itoa(original),
			fmt.Sprintf("The new price is %d%% of the original, so the original is %d / 1.%02d = %d.", 100+percent, newValue, percent, original),
			facts)
	case 2:
		return w10P(o,
			fmt.Sprintf("After a %d%% decrease, a price is %d. What was the original price?", percent, newValue),
			strconv.Itoa(original),
			fmt.Sprintf("The new price is %d%% of the original, so the original is %d / 0.%02d = %d.", 100-percent, newValue, 100-percent, original),
			facts)
	default:
		return w10P(o,
			fmt.Sprintf("After a %d%% increase, a price is %d. How much was added to the original price?", percent, newValue),
			strconv.Itoa(change),
			fmt.Sprintf("The original price is %d, and the new price is %d. The increase is %d - %d = %d.", original, newValue, newValue, original, change),
			facts)
	}
}

// ---------------------------------------------------------------------------
// 4. Literal equations
// ---------------------------------------------------------------------------

type w10LiteralEquationGen struct{}

func (g *w10LiteralEquationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10Ord(ctx)
	facts := map[string]string{}
	switch o {
	case 1:
		l := 2 + rand.Intn(11)
		w := 2 + rand.Intn(11)
		perimeter := 2 * (l + w)
		facts["formula"] = "P = 2l + 2w"
		facts["target"] = "w"
		w10I(facts, "value1", perimeter)
		w10I(facts, "value2", l)
		w10I(facts, "value3", 2)
		return w10P(o,
			fmt.Sprintf("The perimeter of a rectangle is P = 2l + 2w. If P = %d and l = %d, what is w?", perimeter, l),
			strconv.Itoa(w),
			fmt.Sprintf("Subtract 2l from P: %d - 2(%d) = %d. Divide by 2: %d / 2 = %d.", perimeter, l, 2*w, 2*w, w),
			facts)
	case 2:
		l := 2 + rand.Intn(9)
		w := 2 + rand.Intn(9)
		h := 2 + rand.Intn(9)
		volume := l * w * h
		facts["formula"] = "V = lwh"
		facts["target"] = "h"
		w10I(facts, "value1", volume)
		w10I(facts, "value2", l)
		w10I(facts, "value3", w)
		return w10P(o,
			fmt.Sprintf("The volume of a rectangular prism is V = lwh. If V = %d, l = %d and w = %d, what is h?", volume, l, w),
			strconv.Itoa(h),
			fmt.Sprintf("Divide the volume by the product of the other two dimensions: %d / (%d x %d) = %d.", volume, l, w, h),
			facts)
	default:
		b1 := 2 + rand.Intn(9)
		b2 := 2 + rand.Intn(9)
		h := 2 * (1 + rand.Intn(5))
		area := (b1 + b2) * h / 2
		facts["formula"] = "A = (b1 + b2)h / 2"
		facts["target"] = "h"
		w10I(facts, "value1", area)
		w10I(facts, "value2", b1)
		w10I(facts, "value3", b2)
		return w10P(o,
			fmt.Sprintf("The area of a trapezoid is A = (b1 + b2)h / 2. If A = %d, b1 = %d and b2 = %d, what is h?", area, b1, b2),
			strconv.Itoa(h),
			fmt.Sprintf("Multiply the area by 2 and divide by the sum of the bases: 2(%d) / (%d + %d) = %d.", area, b1, b2, h),
			facts)
	}
}
