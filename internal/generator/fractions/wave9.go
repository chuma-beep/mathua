package fractions

// Wave 9 fraction concept: adding mixed numbers with carrying. The generator
// publishes the common denominator and the unregrouped fractional numerator so
// the answer can be recomputed independently.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave9(reg *generator.Registry) {
	reg.Register("frac.mixed.add_carry", &w9MixedAddCarryGen{})
}

type w9MixedAddCarryGen struct{}

func (g *w9MixedAddCarryGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	b := 2 + rand.Intn(5)
	d := 2 + rand.Intn(5)
	a := 1 + rand.Intn(b-1)
	c := 1 + rand.Intn(d-1)
	w1 := 1 + rand.Intn(5)
	w2 := 1 + rand.Intn(5)
	l := lcm(b, d)
	fracNum := a*l/b + c*l/d
	carry := fracNum / l
	rem := fracNum % l
	whole := w1 + w2 + carry
	f := map[string]string{}
	w9FI(f, "w1", w1)
	w9FI(f, "a", a)
	w9FI(f, "b", b)
	w9FI(f, "w2", w2)
	w9FI(f, "c", c)
	w9FI(f, "d", d)
	w9FI(f, "l", l)
	w9FI(f, "frac_num", fracNum)
	w9FI(f, "carry", carry)
	w9FI(f, "whole", whole)
	w9FI(f, "rem", rem)
	switch kp {
	case 1:
		return w9FP(kp, fmt.Sprintf("Add %d %d/%d + %d %d/%d. Do the fractional parts produce a carry? Enter 1 for yes and 0 for no.", w1, a, b, w2, c, d), strconv.Itoa(carry),
			fmt.Sprintf("Over the common denominator %d, the fractional parts sum to %d/%d. That is %s.", l, fracNum, l, map[int]string{0: "less than 1, so there is no carry", 1: "at least 1, so there is a carry"}[carry]), f)
	case 2:
		return w9FP(kp, fmt.Sprintf("What is the whole-number part of %d %d/%d + %d %d/%d? (enter a number)", w1, a, b, w2, c, d), strconv.Itoa(whole),
			fmt.Sprintf("The whole parts sum to %d. The fractional parts sum to %d/%d, which contributes a carry of %d, so the whole part is %d.", w1+w2, fracNum, l, carry, whole), f)
	default:
		if carry == 0 {
			return w9FP(kp, fmt.Sprintf("After carrying, what is the numerator of the fractional part of %d %d/%d + %d %d/%d over the denominator %d? (enter a number)", w1, a, b, w2, c, d, l), strconv.Itoa(rem),
				fmt.Sprintf("The fractional parts sum to %d/%d, which is less than one whole, so there is no carry. The numerator is %d over denominator %d.", fracNum, l, rem, l), f)
		}
		carryWord := "whole units"
		if carry == 1 {
			carryWord = "whole unit"
		}
		return w9FP(kp, fmt.Sprintf("After carrying, what is the numerator of the fractional part of %d %d/%d + %d %d/%d over the denominator %d? (enter a number)", w1, a, b, w2, c, d, l), strconv.Itoa(rem),
			fmt.Sprintf("The fractional parts sum to %d/%d. Regrouping %d %s leaves numerator %d over denominator %d.", fracNum, l, carry, carryWord, rem, l), f)
	}
}

func w9FI(m map[string]string, k string, v int) { m[k] = strconv.Itoa(v) }

func w9FP(o int, q, a, ex string, facts map[string]string) generator.Problem {
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}
