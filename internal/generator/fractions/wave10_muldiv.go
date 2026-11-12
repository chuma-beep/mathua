package fractions

// Wave 10 multiplication, division and mixed-number generators. Every answer is
// an integer: a reduced numerator, a reduced denominator, a whole part, a count
// or a yes/no code.

import (
	"fmt"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

func w10MultWholeCancel(kp int) generator.Problem {
	w := w10R(2, 9)
	b := w10R(2, 9)
	a := w10R(1, b-1)
	n := w * a
	g := mathutil.GCD(n, b)
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("Multiply %d times %d/%d and reduce. What is the numerator?", w, a, b), n/g,
			fmt.Sprintf("Multiply first: %d times %d is %d. Dividing numerator and denominator by %d gives %d/%d.", w, a, n, g, n/g, b/g), "w", w, "a", a, "b", b)
	case 2:
		return w10Problem(kp, fmt.Sprintf("Multiply %d times %d/%d and reduce. What is the denominator?", w, a, b), b/g,
			fmt.Sprintf("The unreduced product is %d/%d, and dividing both parts by %d gives denominator %d.", n, b, g, b/g), "w", w, "a", a, "b", b)
	default:
		return w10Problem(kp, fmt.Sprintf("What is the whole-number part of the product %d times %d/%d?", w, a, b), n/b,
			fmt.Sprintf("The product is %d/%d, which contains %d whole units of size %d.", n, b, n/b, b), "w", w, "a", a, "b", b)
	}
}

func w10EstimateProduct(kp int) generator.Problem {
	b := []int{2, 3, 4, 6}[w10R(0, 3)]
	d := []int{2, 3, 4, 6}[w10R(0, 3)]
	a := w10R(1, b-1)
	c := w10R(1, d-1)
	num := a * c
	den := b * d
	switch kp {
	case 1:
		code := 0
		if num*2 > den {
			code = 2
		} else if num*2 == den {
			code = 1
		}
		return w10Problem(kp, fmt.Sprintf("Is %d/%d times %d/%d closest to 0, 1/2 or 1? Enter 0, 1 or 2.", a, b, c, d), code,
			fmt.Sprintf("The product is %d/%d, which is closest to benchmark code %d.", num, den, code), "a", a, "b", b, "c", c, "d", d)
	case 2:
		code := 0
		if num*2 < den {
			code = 1
		}
		return w10Problem(kp, fmt.Sprintf("Is %d/%d times %d/%d less than 1/2? Enter 1 for yes and 0 for no.", a, b, c, d), code,
			fmt.Sprintf("The product is %d/%d, which is %s 1/2.", num, den, map[bool]string{true: "less than", false: "not less than"}[num*2 < den]), "a", a, "b", b, "c", c, "d", d)
	default:
		q := (num + den/2) / den
		return w10Problem(kp, fmt.Sprintf("Estimate %d/%d times %d/%d as a whole number.", a, b, c, d), q,
			fmt.Sprintf("The product is %d/%d, whose nearest whole number is %d.", num, den, q), "a", a, "b", b, "c", c, "d", d)
	}
}

func w10MultScaling(kp int) generator.Problem {
	w := w10R(2, 9)
	a := w10R(1, 9)
	b := w10R(2, 9)
	n := w * a
	g := mathutil.GCD(n, b)
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("Scale %d by %d/%d and reduce. What is the numerator?", w, a, b), n/g,
			fmt.Sprintf("%d times %d is %d, and both parts divide by %d.", w, a, n, g), "w", w, "a", a, "b", b)
	case 2:
		return w10Problem(kp, fmt.Sprintf("Scale %d by %d/%d. What is the whole-number part of the result?", w, a, b), n/b,
			fmt.Sprintf("The result is %d/%d, which contains %d whole units.", n, b, n/b), "w", w, "a", a, "b", b)
	default:
		code := 0
		if a > b {
			code = 1
		}
		return w10Problem(kp, fmt.Sprintf("Scaling %d by %d/%d gives a result larger than %d. Enter 1 if this is true and 0 if false.", w, a, b, w), code,
			fmt.Sprintf("The scale factor %d/%d is %s 1, so the result is %s the original.", a, b, map[bool]string{true: "greater than", false: "not greater than"}[a > b], map[bool]string{true: "larger than", false: "not larger than"}[a > b]), "w", w, "a", a, "b", b)
	}
}

func w10MultUnit(kp int) generator.Problem {
	a := w10R(1, 9)
	b := w10R(2, 9)
	c := w10R(2, 9)
	num := a
	g := mathutil.GCD(num, b*c)
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("Multiply %d/%d by 1/%d and reduce. What is the denominator?", a, b, c), b*c/g,
			fmt.Sprintf("The product is %d/%d, and dividing both parts by %d gives denominator %d.", a, b*c, g, b*c/g), "a", a, "b", b, "c", c)
	case 2:
		return w10Problem(kp, fmt.Sprintf("Multiply %d/%d by 1/%d and reduce. What is the numerator?", a, b, c), num/g,
			fmt.Sprintf("The product is %d/%d, which reduces to numerator %d.", a, b*c, num/g), "a", a, "b", b, "c", c)
	default:
		return w10Problem(kp, fmt.Sprintf("What is the denominator of the unreduced product %d/%d times 1/%d?", a, b, c), b*c,
			fmt.Sprintf("Multiplying denominators gives %d times %d, which is %d.", b, c, b*c), "a", a, "b", b, "c", c)
	}
}

func w10MultThree(kp int) generator.Problem {
	a := w10R(1, 8)
	b := w10R(2, 9)
	c := w10R(1, 8)
	d := w10R(2, 9)
	e := w10R(1, 8)
	f := w10R(2, 9)
	num := a * c * e
	den := b * d * f
	g := mathutil.GCD(num, den)
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("Multiply %d/%d times %d/%d times %d/%d and reduce. What is the numerator?", a, b, c, d, e, f), num/g,
			fmt.Sprintf("The numerator product is %d and the denominator product is %d. Dividing both by %d gives %d/%d.", num, den, g, num/g, den/g), "a", a, "b", b, "c", c, "d", d, "e", e, "f", f)
	case 2:
		return w10Problem(kp, fmt.Sprintf("Multiply %d/%d times %d/%d times %d/%d and reduce. What is the denominator?", a, b, c, d, e, f), den/g,
			fmt.Sprintf("The unreduced product is %d/%d, which reduces to denominator %d.", num, den, den/g), "a", a, "b", b, "c", c, "d", d, "e", e, "f", f)
	default:
		code := 0
		if num*2 < den {
			code = 1
		}
		return w10Problem(kp, fmt.Sprintf("Is the product %d/%d times %d/%d times %d/%d less than 1/2? Enter 1 for yes and 0 for no.", a, b, c, d, e, f), code,
			fmt.Sprintf("The product is %d/%d, which is %s 1/2.", num, den, map[bool]string{true: "less than", false: "not less than"}[num*2 < den]), "a", a, "b", b, "c", c, "d", d, "e", e, "f", f)
	}
}

func w10MultMixedImproper(kp int) generator.Problem {
	b := w10R(2, 9)
	a := w10R(1, b-1)
	w := w10R(1, 6)
	c := w10R(1, 8)
	d := w10R(2, 9)
	improper := w*b + a
	num := improper * c
	den := b * d
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("Before multiplying %d %d/%d by %d/%d, write the mixed number as an improper fraction. What is its numerator?", w, a, b, c, d), improper,
			fmt.Sprintf("%d times %d is %d, and adding %d gives %d.", w, b, w*b, a, improper), "w", w, "a", a, "b", b, "c", c, "d", d)
	case 2:
		return w10Problem(kp, fmt.Sprintf("Multiply %d %d/%d by %d/%d over the denominator %d. What is the unreduced numerator?", w, a, b, c, d, den), num,
			fmt.Sprintf("The improper numerator is %d, and multiplying by %d gives %d.", improper, c, num), "w", w, "a", a, "b", b, "c", c, "d", d)
	default:
		return w10Problem(kp, fmt.Sprintf("What is the whole-number part of %d %d/%d times %d/%d?", w, a, b, c, d), num/den,
			fmt.Sprintf("The product is %d/%d, which contains %d whole units of size %d.", num, den, num/den, den), "w", w, "a", a, "b", b, "c", c, "d", d)
	}
}

func w10MultDecimal(kp int) generator.Problem {
	b := []int{2, 4, 5, 10, 20, 25, 50}[w10R(0, 6)]
	a := w10R(1, b-1)
	dec := w10R(1, 99)
	num := a * dec
	den := b * 100
	g := mathutil.GCD(num, den)
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("Multiply %d/%d by 0.%02d and reduce. What is the numerator?", a, b, dec), num/g,
			fmt.Sprintf("The product is %d/%d, and dividing both parts by %d gives numerator %d.", num, den, g, num/g), "a", a, "b", b, "dec", dec)
	case 2:
		return w10Problem(kp, fmt.Sprintf("Multiply %d/%d by 0.%02d and reduce. What is the denominator?", a, b, dec), den/g,
			fmt.Sprintf("The unreduced product is %d/%d, which reduces to denominator %d.", num, den, den/g), "a", a, "b", b, "dec", dec)
	default:
		code := 0
		if num*2 > den {
			code = 1
		}
		return w10Problem(kp, fmt.Sprintf("Is %d/%d times 0.%02d greater than 1/2? Enter 1 for yes and 0 for no.", a, b, dec), code,
			fmt.Sprintf("The product is %d/%d, which is %s 1/2.", num, den, map[bool]string{true: "greater than", false: "not greater than"}[num*2 > den]), "a", a, "b", b, "dec", dec)
	}
}

func w10MultPredict(kp int) generator.Problem {
	a := w10R(1, 9)
	b := w10R(2, 9)
	c := w10R(1, 9)
	d := w10R(2, 9)
	for c == d {
		c = w10R(1, 9)
	}
	switch kp {
	case 1:
		code := 0
		if c > d {
			code = 1
		}
		return w10Problem(kp, fmt.Sprintf("Without computing, is %d/%d times %d/%d larger than %d/%d? Enter 1 for yes and 0 for no.", a, b, c, d, a, b), code,
			fmt.Sprintf("The second factor is %s 1, so the product is %s the first factor.", map[bool]string{true: "greater than", false: "less than"}[c > d], map[bool]string{true: "larger than", false: "smaller than"}[c > d]), "a", a, "b", b, "c", c, "d", d)
	case 2:
		code := 0
		if a*c > b*d {
			code = 1
		}
		return w10Problem(kp, fmt.Sprintf("Is %d/%d times %d/%d greater than 1? Enter 1 for yes and 0 for no.", a, b, c, d), code,
			fmt.Sprintf("The product is %d/%d, which is %s 1.", a*c, b*d, map[bool]string{true: "greater than", false: "not greater than"}[a*c > b*d]), "a", a, "b", b, "c", c, "d", d)
	default:
		code := 0
		if a*c > b*d {
			code = 2
		} else if a*c == b*d {
			code = 1
		}
		return w10Problem(kp, fmt.Sprintf("Classify the product %d/%d times %d/%d: enter 0 if below 1, 1 if equal to 1 and 2 if above 1.", a, b, c, d), code,
			fmt.Sprintf("Compare %d with %d: the product is code %d.", a*c, b*d, code), "a", a, "b", b, "c", c, "d", d)
	}
}

func w10DivReciprocal(kp int) generator.Problem {
	a := w10R(1, 9)
	b := w10R(2, 9)
	c := w10R(1, 9)
	d := w10R(2, 9)
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("To divide %d/%d by %d/%d, invert the second fraction. What is the numerator of the resulting product?", a, b, c, d), a*d,
			fmt.Sprintf("Inverting gives %d/%d, and the numerator is %d times %d, which is %d.", d, c, a, d, a*d), "a", a, "b", b, "c", c, "d", d)
	case 2:
		return w10Problem(kp, fmt.Sprintf("After inverting %d/%d, what is the denominator of %d/%d times %d/%d?", c, d, a, b, d, c), b*c,
			fmt.Sprintf("Multiplying denominators gives %d times %d, which is %d.", b, c, b*c), "a", a, "b", b, "c", c, "d", d)
	default:
		num, den := a*d, b*c
		return w10Problem(kp, fmt.Sprintf("What is the whole-number part of %d/%d divided by %d/%d?", a, b, c, d), num/den,
			fmt.Sprintf("The quotient is %d/%d, which contains %d whole units.", num, den, num/den), "a", a, "b", b, "c", c, "d", d)
	}
}

func w10EstimateQuotient(kp int) generator.Problem {
	a := w10R(1, 9)
	b := w10R(2, 9)
	c := w10R(1, 9)
	d := w10R(2, 9)
	num, den := a*d, b*c
	switch kp {
	case 1:
		q := (num + den/2) / den
		return w10Problem(kp, fmt.Sprintf("Estimate %d/%d divided by %d/%d as a whole number.", a, b, c, d), q,
			fmt.Sprintf("The quotient is %d/%d, whose nearest whole number is %d.", num, den, q), "a", a, "b", b, "c", c, "d", d)
	case 2:
		code := 0
		if num > den {
			code = 1
		}
		return w10Problem(kp, fmt.Sprintf("Is %d/%d divided by %d/%d greater than 1? Enter 1 for yes and 0 for no.", a, b, c, d), code,
			fmt.Sprintf("The quotient is %d/%d, which is %s 1.", num, den, map[bool]string{true: "greater than", false: "not greater than"}[num > den]), "a", a, "b", b, "c", c, "d", d)
	default:
		code := 0
		if num*2 > den {
			code = 2
		} else if num*2 == den {
			code = 1
		}
		return w10Problem(kp, fmt.Sprintf("Is the quotient %d/%d divided by %d/%d closest to 0, 1/2 or 1? Enter 0, 1 or 2.", a, b, c, d), code,
			fmt.Sprintf("The quotient %d/%d is closest to benchmark code %d.", num, den, code), "a", a, "b", b, "c", c, "d", d)
	}
}

func w10ComplexSimplify(kp int) generator.Problem {
	a := w10R(1, 9)
	b := w10R(2, 9)
	c := w10R(1, 9)
	d := w10R(2, 9)
	num, den := a*d, b*c
	g := mathutil.GCD(num, den)
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("Simplify the complex fraction (%d/%d)/(%d/%d). What is the numerator?", a, b, c, d), num/g,
			fmt.Sprintf("Rewrite as %d/%d times %d/%d, giving %d/%d, which reduces to numerator %d.", a, b, d, c, num, den, num/g), "a", a, "b", b, "c", c, "d", d)
	case 2:
		return w10Problem(kp, fmt.Sprintf("Simplify the complex fraction (%d/%d)/(%d/%d). What is the denominator?", a, b, c, d), den/g,
			fmt.Sprintf("The quotient %d/%d reduces to denominator %d.", num, den, den/g), "a", a, "b", b, "c", c, "d", d)
	default:
		return w10Problem(kp, fmt.Sprintf("What is the whole-number part of (%d/%d)/(%d/%d)?", a, b, c, d), num/den,
			fmt.Sprintf("The quotient is %d/%d, which contains %d whole units.", num, den, num/den), "a", a, "b", b, "c", c, "d", d)
	}
}

func w10DivUnit(kp int) generator.Problem {
	a := w10R(1, 9)
	b := w10R(2, 9)
	c := w10R(2, 9)
	num, den := a*c, b
	g := mathutil.GCD(num, den)
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("Divide %d/%d by 1/%d and reduce. What is the numerator?", a, b, c), num/g,
			fmt.Sprintf("Multiplying by the reciprocal %d gives %d/%d, which reduces to numerator %d.", c, num, den, num/g), "a", a, "b", b, "c", c)
	case 2:
		return w10Problem(kp, fmt.Sprintf("Divide %d/%d by 1/%d and reduce. What is the denominator?", a, b, c), den/g,
			fmt.Sprintf("The quotient %d/%d reduces to denominator %d.", num, den, den/g), "a", a, "b", b, "c", c)
	default:
		return w10Problem(kp, fmt.Sprintf("How many one-%d-ths fit into %d/%d? Give the pile count over the denominator %d.", c, a, b, b), num,
			fmt.Sprintf("Dividing by 1/%d multiplies by %d, so the numerator is %d times %d, which is %d over %d.", c, c, a, c, num, b), "a", a, "b", b, "c", c)
	}
}

func w10DivMeasurement(kp int) generator.Problem {
	b := w10R(2, 9)
	a := w10R(1, b-1)
	w := w10R(1, 5)
	c := w10R(2, 6)
	total := w*b + a
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("How many 1/%d pieces fit into %d %d/%d? Enter the count.", c, w, a, b), total*c/b,
			fmt.Sprintf("The mixed number is %d/%d, and dividing by 1/%d multiplies by %d, giving %d pieces.", total, b, c, c, total*c/b), "w", w, "a", a, "b", b, "c", c)
	case 2:
		return w10Problem(kp, fmt.Sprintf("How many 2/%d pieces fit into %d %d/%d? Enter the count.", c, w, a, b), total*c/(2*b),
			fmt.Sprintf("Dividing %d/%d by 2/%d multiplies by %d/2, giving %d pieces.", total, b, c, c, total*c/(2*b)), "w", w, "a", a, "b", b, "c", c)
	default:
		return w10Problem(kp, fmt.Sprintf("After fitting as many 1/%d pieces as possible into %d %d/%d, how many pieces fit?", c, w, a, b), total*c/b,
			fmt.Sprintf("The total numerator over %d is %d, so %d pieces of size 1/%d fit.", b, total, total*c/b, c), "w", w, "a", a, "b", b, "c", c)
	}
}

func w10DivPartitive(kp int) generator.Problem {
	a := w10R(1, 9)
	b := w10R(2, 9)
	c := w10R(2, 9)
	num, den := a, b*c
	g := mathutil.GCD(num, den)
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("Share %d/%d equally among %d people and reduce. What is the numerator of one share?", a, b, c), num/g,
			fmt.Sprintf("Each share is %d/(%d times %d), which is %d/%d, reducing to numerator %d.", a, b, c, num, den, num/g), "a", a, "b", b, "c", c)
	case 2:
		return w10Problem(kp, fmt.Sprintf("Share %d/%d equally among %d people and reduce. What is the denominator of one share?", a, b, c), den/g,
			fmt.Sprintf("One share is %d/%d, which reduces to denominator %d.", num, den, den/g), "a", a, "b", b, "c", c)
	default:
		return w10Problem(kp, fmt.Sprintf("Share %d/%d equally among %d people using the unreduced denominator. What is the denominator of one share?", a, b, c), b*c,
			fmt.Sprintf("Dividing by %d multiplies the denominator: %d times %d is %d.", c, b, c, b*c), "a", a, "b", b, "c", c)
	}
}

func w10DivThree(kp int) generator.Problem {
	a := w10R(1, 8)
	b := w10R(2, 9)
	c := w10R(1, 8)
	d := w10R(2, 9)
	e := w10R(1, 8)
	f := w10R(2, 9)
	num, den := a*d*f, b*c*e
	g := mathutil.GCD(num, den)
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("Compute %d/%d divided by %d/%d divided by %d/%d and reduce. What is the numerator?", a, b, c, d, e, f), num/g,
			fmt.Sprintf("Invert both divisors and multiply: the numerator is %d times %d times %d, which is %d.", a, d, f, num), "a", a, "b", b, "c", c, "d", d, "e", e, "f", f)
	case 2:
		return w10Problem(kp, fmt.Sprintf("Compute %d/%d divided by %d/%d divided by %d/%d and reduce. What is the denominator?", a, b, c, d, e, f), den/g,
			fmt.Sprintf("The unreduced denominator is %d times %d times %d, which is %d, and it reduces to %d.", b, c, e, den, den/g), "a", a, "b", b, "c", c, "d", d, "e", e, "f", f)
	default:
		code := 0
		if num > den {
			code = 1
		}
		return w10Problem(kp, fmt.Sprintf("Is %d/%d divided by %d/%d divided by %d/%d greater than 1? Enter 1 for yes and 0 for no.", a, b, c, d, e, f), code,
			fmt.Sprintf("The quotient is %d/%d, which is %s 1.", num, den, map[bool]string{true: "greater than", false: "not greater than"}[num > den]), "a", a, "b", b, "c", c, "d", d, "e", e, "f", f)
	}
}

func w10DivZero(kp int) generator.Problem {
	switch kp {
	case 1:
		zeros := 1 + w10R(0, 2)
		parts := ""
		for i := 0; i < 3; i++ {
			if i < zeros {
				parts += fmt.Sprintf("x%d/0", i+1)
			} else {
				parts += fmt.Sprintf("x%d/%d", i+1, w10R(2, 9))
			}
			if i < 2 {
				parts += ", "
			}
		}
		return w10Problem(kp, fmt.Sprintf("How many of these expressions are undefined: %s?", parts), zeros,
			fmt.Sprintf("Each expression with denominator zero is undefined, and %d of the three qualify.", zeros), "zeros", zeros)
	case 2:
		b := w10R(2, 12)
		return w10Problem(kp, fmt.Sprintf("For which value of x is 5/(x-%d) undefined?", b), b,
			fmt.Sprintf("The denominator is zero when x equals %d.", b), "b", b)
	default:
		b := w10R(2, 9)
		return w10Problem(kp, fmt.Sprintf("Compute 0 divided by %d/%d. Enter the result.", b, w10R(2, 9)), 0,
			fmt.Sprintf("Zero divided by any nonzero fraction is zero."), "b", b)
	}
}

func w10MixedCompare(kp int) generator.Problem {
	w1 := w10R(1, 8)
	w2 := w10R(1, 8)
	b := w10R(2, 9)
	d := w10R(2, 9)
	a := w10R(1, b-1)
	c := w10R(1, d-1)
	for w1 == w2 && a*d == c*b {
		w2 = w10R(1, 8)
	}
	switch kp {
	case 1:
		code := 1
		if w2 > w1 {
			code = 2
		}
		return w10Problem(kp, fmt.Sprintf("Compare %d %d/%d and %d %d/%d by whole part first. Which whole part is larger? Enter 1 or 2.", w1, a, b, w2, c, d), code,
			fmt.Sprintf("The whole parts are %d and %d, so position %d is larger.", w1, w2, code), "w1", w1, "a", a, "b", b, "w2", w2, "c", c, "d", d)
	case 2:
		left := (w1*b + a) * d
		right := (w2*d + c) * b
		code := 1
		if right > left {
			code = 2
		}
		return w10Problem(kp, fmt.Sprintf("Which is larger, %d %d/%d or %d %d/%d? Enter 1 or 2.", w1, a, b, w2, c, d), code,
			fmt.Sprintf("Cross multiplying gives %d against %d, so position %d is larger.", left, right, code), "w1", w1, "a", a, "b", b, "w2", w2, "c", c, "d", d)
	default:
		diff := w1 - w2
		if diff < 0 {
			diff = -diff
		}
		return w10Problem(kp, fmt.Sprintf("How many whole units separate %d %d/%d and %d %d/%d?", w1, a, b, w2, c, d), diff,
			fmt.Sprintf("Subtract the whole parts: the distance is %d whole units.", diff), "w1", w1, "a", a, "b", b, "w2", w2, "c", c, "d", d)
	}
}

func w10MixedNegative(kp int) generator.Problem {
	b := w10R(2, 9)
	a := w10R(1, b-1)
	w := w10R(1, 6)
	improper := w*b + a
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("Write -(%d %d/%d) as a negative improper fraction over %d. What is the magnitude of the numerator?", w, a, b, b), improper,
			fmt.Sprintf("The positive mixed number has numerator %d, so the negative form has magnitude %d.", improper, improper), "w", w, "a", a, "b", b)
	case 2:
		return w10Problem(kp, fmt.Sprintf("Write -(%d %d/%d) as an improper fraction over %d. What is the signed numerator?", w, a, b, b), -improper,
			fmt.Sprintf("The mixed number is %d/%d, so the negative is -%d/%d.", improper, b, improper, b), "w", w, "a", a, "b", b)
	default:
		return w10Problem(kp, fmt.Sprintf("How many whole units are in -(%d %d/%d) when written as a mixed number?", w, a, b), w,
			fmt.Sprintf("The magnitude is %d and %d/%d, so the negative mixed number has whole part -%d.", w, a, b, w), "w", w, "a", a, "b", b)
	}
}

func w10MixedBorrowTwice(kp int) generator.Problem {
	b := w10R(3, 9)
	a := w10R(1, b-2)
	d := w10R(a+2, a+b)
	c := w10R(a+1, d-1)
	w2 := w10R(1, 5)
	w1 := w2 + w10R(1, 3)
	whole := w1 - w2 - 1
	num := a + b - c
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("After borrowing one whole, what is the whole-number part of %d %d/%d - %d %d/%d?", w1, a, b, w2, c, d), whole,
			fmt.Sprintf("Borrowing turns the fractional part into %d/%d, and the whole part becomes %d minus %d minus 1, which is %d.", a+b, b, w1, w2, whole), "w1", w1, "a", a, "b", b, "w2", w2, "c", c, "d", d)
	case 2:
		return w10Problem(kp, fmt.Sprintf("What is the numerator of the fractional part of %d %d/%d - %d %d/%d over the denominator %d?", w1, a, b, w2, c, d, b), num,
			fmt.Sprintf("The borrowed numerator is %d, and subtracting %d gives %d.", a+b, c, num), "w1", w1, "a", a, "b", b, "w2", w2, "c", c, "d", d)
	default:
		return w10Problem(kp, fmt.Sprintf("If %d %d/%d and %d %d/%d were written over the common denominator %d, what would the numerator difference be?", w1, a, b, w2, c, d, b), (w1-w2)*b+num,
			fmt.Sprintf("The difference is %d whole units and %d/%d, giving numerator %d over %d.", whole, num, b, (w1-w2)*b+num, b), "w1", w1, "a", a, "b", b, "w2", w2, "c", c, "d", d)
	}
}

func w10MixedDivWhole(kp int) generator.Problem {
	b := w10R(2, 9)
	a := w10R(1, b-1)
	w := w10R(1, 6)
	c := w10R(2, 6)
	improper := w*b + a
	num, den := improper, b*c
	g := mathutil.GCD(num, den)
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("Divide %d %d/%d by %d and reduce. What is the numerator?", w, a, b, c), num/g,
			fmt.Sprintf("The mixed number is %d/%d, so the quotient is %d/%d, reducing to numerator %d.", improper, b, num, den, num/g), "w", w, "a", a, "b", b, "c", c)
	case 2:
		return w10Problem(kp, fmt.Sprintf("Divide %d %d/%d by %d and reduce. What is the denominator?", w, a, b, c), den/g,
			fmt.Sprintf("The quotient %d/%d reduces to denominator %d.", num, den, den/g), "w", w, "a", a, "b", b, "c", c)
	default:
		return w10Problem(kp, fmt.Sprintf("What is the whole-number part of %d %d/%d divided by %d?", w, a, b, c), num/den,
			fmt.Sprintf("The quotient is %d/%d, which contains %d whole units.", num, den, num/den), "w", w, "a", a, "b", b, "c", c)
	}
}

func w10MixedEstimate(kp int) generator.Problem {
	w1 := w10R(1, 8)
	w2 := w10R(1, 8)
	b := w10R(2, 6)
	d := w10R(2, 6)
	a := w10R(1, b-1)
	c := w10R(1, d-1)
	sum := float64(w1+w2) + float64(a)/float64(b) + float64(c)/float64(d)
	diff := float64(w1-w2) + float64(a)/float64(b) - float64(c)/float64(d)
	if diff < 0 {
		diff = -diff
	}
	round := func(v float64) int { return int(v + 0.5) }
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("Estimate %d %d/%d + %d %d/%d as a whole number.", w1, a, b, w2, c, d), round(sum),
			fmt.Sprintf("The sum is about %.2f, which rounds to %d.", sum, round(sum)), "w1", w1, "a", a, "b", b, "w2", w2, "c", c, "d", d)
	case 2:
		return w10Problem(kp, fmt.Sprintf("Estimate the absolute difference between %d %d/%d and %d %d/%d as a whole number.", w1, a, b, w2, c, d), round(diff),
			fmt.Sprintf("The difference is about %.2f, which rounds to %d.", diff, round(diff)), "w1", w1, "a", a, "b", b, "w2", w2, "c", c, "d", d)
	default:
		code := 0
		if sum > float64(w1+w2+1) {
			code = 1
		}
		return w10Problem(kp, fmt.Sprintf("Is %d %d/%d + %d %d/%d greater than %d? Enter 1 for yes and 0 for no.", w1, a, b, w2, c, d, w1+w2+1), code,
			fmt.Sprintf("The sum is about %.2f, which is %s %d.", sum, map[bool]string{true: "greater than", false: "not greater than"}[code == 1], w1+w2+1), "w1", w1, "a", a, "b", b, "w2", w2, "c", c, "d", d)
	}
}

func w10MixedPercent(kp int) generator.Problem {
	b := []int{2, 4, 5, 10, 20, 25, 50}[w10R(0, 6)]
	a := w10R(1, b-1)
	w := w10R(1, 6)
	improper := w*b + a
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("Write %d %d/%d as a percent.", w, a, b), improper*100/b,
			fmt.Sprintf("The mixed number is %d/%d, and multiplying by 100 gives %d percent.", improper, b, improper*100/b), "w", w, "a", a, "b", b)
	case 2:
		return w10Problem(kp, fmt.Sprintf("What percent is the fractional part %d/%d?", a, b), a*100/b,
			fmt.Sprintf("The fraction %d/%d is %d/100, which is %d percent.", a, b, a*100/b, a*100/b), "a", a, "b", b)
	default:
		return w10Problem(kp, fmt.Sprintf("Write %d %d/%d as a fraction over 100. What is the numerator?", w, a, b), improper*100/b,
			fmt.Sprintf("The percent is %d, so the numerator over 100 is %d.", improper*100/b, improper*100/b), "w", w, "a", a, "b", b)
	}
}

func w10MixedWordCompare(kp int) generator.Problem {
	w1 := w10R(1, 6)
	w2 := w10R(1, 6)
	b := w10R(2, 8)
	d := w10R(2, 8)
	a := w10R(1, b-1)
	c := w10R(1, d-1)
	for w1 == w2 && a*d == c*b {
		w2 = w10R(1, 6)
	}
	code := 1
	if float64(w2)+float64(c)/float64(d) > float64(w1)+float64(a)/float64(b) {
		code = 2
	}
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("A runner covers %d %d/%d laps and a second covers %d %d/%d. Who covers more? Enter 1 or 2.", w1, a, b, w2, c, d), code,
			fmt.Sprintf("Comparing the mixed numbers, position %d is larger.", code), "w1", w1, "a", a, "b", b, "w2", w2, "c", c, "d", d)
	case 2:
		return w10Problem(kp, fmt.Sprintf("Two jars hold %d %d/%d litres and %d %d/%d litres. Which jar holds more? Enter 1 or 2.", w1, a, b, w2, c, d), code,
			fmt.Sprintf("Cross multiplying the fractional forms shows that position %d is larger.", code), "w1", w1, "a", a, "b", b, "w2", w2, "c", c, "d", d)
	default:
		diff := w1 - w2
		if diff < 0 {
			diff = -diff
		}
		return w10Problem(kp, fmt.Sprintf("Two hikers walk %d %d/%d km and %d %d/%d km. How many whole kilometres separate them?", w1, a, b, w2, c, d), diff,
			fmt.Sprintf("The whole parts differ by %d, so the distances are %d whole kilometres apart.", diff, diff), "w1", w1, "a", a, "b", b, "w2", w2, "c", c, "d", d)
	}
}
