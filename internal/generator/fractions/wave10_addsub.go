package fractions

// Wave 10 addition and subtraction generators. Answers are numerators,
// denominators, whole parts or yes/no codes so every result is an integer.

import (
	"fmt"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

func w10AddWhole(kp int) generator.Problem {
	w := w10R(1, 9)
	b := w10R(2, 9)
	a := w10R(1, b-1)
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("Write %d + %d/%d as an improper fraction over %d. What is the numerator?", w, a, b, b), w*b+a,
			fmt.Sprintf("Convert the whole number to %d/%d, then add the numerator: %d plus %d is %d.", w*b, b, w*b, a, w*b+a), "w", w, "a", a, "b", b)
	case 2:
		return w10Problem(kp, fmt.Sprintf("What is the whole-number part of %d + %d/%d?", w, a, b), w,
			fmt.Sprintf("The fractional part is less than one, so the whole part stays %d.", w), "w", w, "a", a, "b", b)
	default:
		n := w*b + 2*b + a
		return w10Problem(kp, fmt.Sprintf("Write %d + %d/%d, where the fraction is improper, as an improper fraction over %d. What is the numerator?", w, 2*b+a, b, b), n,
			fmt.Sprintf("The fractional part is %d/%d, which is 2 and %d/%d, so the total numerator over %d is %d.", 2*b+a, b, a, b, b, n), "w", w, "b", b, "a", a)
	}
}

func w10EstimateSum(kp int) generator.Problem {
	b := []int{2, 3, 4, 6}[w10R(0, 3)]
	d := []int{2, 3, 4, 6}[w10R(0, 3)]
	a := w10R(1, b-1)
	c := w10R(1, d-1)
	num := a*d + c*b
	den := b * d
	switch kp {
	case 1:
		q := (num + den/2) / den
		return w10Problem(kp, fmt.Sprintf("Estimate %d/%d + %d/%d by rounding to the nearest whole number.", a, b, c, d), q,
			fmt.Sprintf("The exact sum is %d/%d, which is closest to %d.", num, den, q), "a", a, "b", b, "c", c, "d", d)
	case 2:
		code := 0
		if num > den {
			code = 1
		}
		return w10Problem(kp, fmt.Sprintf("Is %d/%d + %d/%d greater than 1? Enter 1 for yes and 0 for no.", a, b, c, d), code,
			fmt.Sprintf("The sum is %d/%d, which is %s 1.", num, den, map[bool]string{true: "greater than", false: "not greater than"}[num > den]), "a", a, "b", b, "c", c, "d", d)
	default:
		b = w10R(4, 9)
		d = w10R(4, 9)
		a = w10R(1, b-1)
		c = w10R(1, d-1)
		num = a*d + c*b
		den = b * d
		q := (num + den/2) / den
		return w10Problem(kp, fmt.Sprintf("Estimate %d/%d + %d/%d by rounding to the nearest whole number.", a, b, c, d), q,
			fmt.Sprintf("The exact sum is %d/%d, and its nearest whole number is %d.", num, den, q), "a", a, "b", b, "c", c, "d", d)
	}
}

func w10AddRegroup(kp int) generator.Problem {
	b := w10R(3, 9)
	a := w10R(2, b-1)
	c := w10R(b-a+1, b-1)
	carry := (a + c) / b
	rem := (a + c) % b
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("%d/%d + %d/%d is an improper sum. What is the whole-number carry?", a, b, c, b), carry,
			fmt.Sprintf("The numerators add to %d, which contains %d whole %ds.", a+c, carry, b), "a", a, "b", b, "c", c)
	case 2:
		return w10Problem(kp, fmt.Sprintf("After regrouping %d/%d + %d/%d, what is the numerator of the fractional part over %d?", a, b, c, b, b), rem,
			fmt.Sprintf("The numerators add to %d. Regrouping %d whole units leaves %d over %d.", a+c, carry, rem, b), "a", a, "b", b, "c", c)
	default:
		return w10Problem(kp, fmt.Sprintf("Write %d/%d + %d/%d as a single improper fraction over %d. What is the numerator?", a, b, c, b, b), a+c,
			fmt.Sprintf("With a common denominator, add the numerators: %d plus %d is %d.", a, c, a+c), "a", a, "b", b, "c", c)
	}
}

func w10AddUnit(kp int) generator.Problem {
	a := w10R(2, 8)
	b := a * w10R(2, 6)
	num := b/a + 1
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("Add 1/%d + 1/%d over the common denominator %d. What is the numerator?", a, b, b), num,
			fmt.Sprintf("The first fraction becomes %d/%d, and adding 1 gives %d.", b/a, b, num), "a", a, "b", b)
	case 2:
		return w10Problem(kp, fmt.Sprintf("What is the least common denominator of 1/%d and 1/%d?", a, b), b,
			fmt.Sprintf("Since %d is a multiple of %d, the least common denominator is %d.", b, a, b), "a", a, "b", b)
	default:
		n := w10R(2, 9)
		return w10Problem(kp, fmt.Sprintf("How many copies of the unit fraction 1/%d make one whole?", n), n,
			fmt.Sprintf("A whole splits into %d equal parts of size 1/%d.", n, n), "n", n)
	}
}

func w10MissingAddend(kp int) generator.Problem {
	switch kp {
	case 1:
		b := w10R(2, 9)
		a := w10R(1, b-1)
		c := w10R(a+1, b)
		return w10Problem(kp, fmt.Sprintf("Find the missing numerator: %d/%d + x/%d = %d/%d.", a, b, b, c, b), c-a,
			fmt.Sprintf("The denominators match, so subtract the known numerator: %d minus %d is %d.", c, a, c-a), "a", a, "b", b, "c", c)
	case 2:
		b := w10R(2, 7)
		d := b + w10R(1, 4)
		a := w10R(1, b-1)
		c := w10R(1, d-1)
		num := c*b - a*d
		for num <= 0 {
			c = w10R(1, d-1)
			num = c*b - a*d
		}
		return w10Problem(kp, fmt.Sprintf("Solve %d/%d + x = %d/%d. What is the numerator of x over the common denominator %d?", a, b, c, d, b*d), num,
			fmt.Sprintf("Convert both fractions to denominator %d and subtract: %d minus %d is %d.", b*d, c*b, a*d, num), "a", a, "b", b, "c", c, "d", d)
	default:
		b := w10R(2, 9)
		a := w10R(1, b-1)
		w := w10R(2, 6)
		c := a + w*b
		return w10Problem(kp, fmt.Sprintf("Find the whole number: w + %d/%d = %d/%d.", a, b, c, b), w,
			fmt.Sprintf("Subtract the fractional numerator: %d minus %d is %d, which is %d wholes of size %d.", c, a, c-a, w, b), "a", a, "b", b, "c", c, "w", w)
	}
}

func w10AddAssociative(kp int) generator.Problem {
	b := w10R(3, 9)
	a := w10R(1, b-1)
	c := w10R(1, b-1)
	e := w10R(1, b-1)
	total := a + c + e
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("Add %d/%d + %d/%d + %d/%d. What is the numerator over %d?", a, b, c, b, e, b, b), total,
			fmt.Sprintf("The denominators match, so add the numerators in any order: %d + %d + %d is %d.", a, c, e, total), "a", a, "b", b, "c", c, "e", e)
	case 2:
		return w10Problem(kp, fmt.Sprintf("Compute %d/%d + %d/%d + %d/%d. What is the whole-number part?", a, b, c, b, e, b), total/b,
			fmt.Sprintf("The numerators total %d, which contains %d whole units of size %d.", total, total/b, b), "a", a, "b", b, "c", c, "e", e)
	default:
		return w10Problem(kp, fmt.Sprintf("Group the first two terms of %d/%d + %d/%d + %d/%d. What is the numerator of that first sum over %d?", a, b, c, b, e, b, b), a+c,
			fmt.Sprintf("Grouping does not change the value, and the first pair adds to %d/%d.", a+c, b), "a", a, "b", b, "c", c, "e", e)
	}
}

func w10AddSameNumerator(kp int) generator.Problem {
	a := w10R(1, 8)
	b := w10R(2, 9)
	d := w10R(2, 9)
	for d == b {
		d = w10R(2, 9)
	}
	l := lcm(b, d)
	num := a * (l/b + l/d)
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("Add %d/%d + %d/%d. What is the numerator over the common denominator %d?", a, b, a, d, l), num,
			fmt.Sprintf("The shared numerator %d scales by %d and %d, giving %d over %d.", a, l/b, l/d, num, l), "a", a, "b", b, "d", d, "l", l)
	case 2:
		return w10Problem(kp, fmt.Sprintf("What is the least common denominator of %d/%d and %d/%d?", a, b, a, d), l,
			fmt.Sprintf("The least common multiple of %d and %d is %d.", b, d, l), "a", a, "b", b, "d", d)
	default:
		code := 0
		if num > l {
			code = 1
		}
		return w10Problem(kp, fmt.Sprintf("Is %d/%d + %d/%d greater than 1? Enter 1 for yes and 0 for no.", a, b, a, d), code,
			fmt.Sprintf("The sum is %d/%d, which is %s 1.", num, l, map[bool]string{true: "greater than", false: "not greater than"}[num > l]), "a", a, "b", b, "d", d)
	}
}

func w10AddDecimal(kp int) generator.Problem {
	b := []int{2, 4, 5, 10, 20, 25, 50}[w10R(0, 6)]
	a := w10R(1, b-1)
	fracHundred := a * 100 / b
	dec := w10R(1, 99)
	sum := fracHundred + dec
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("Add %d/%d + 0.%02d. Give the numerator over 100.", a, b, dec), sum,
			fmt.Sprintf("The fraction is %d/100, and adding 0.%02d gives %d/100.", fracHundred, dec, sum), "a", a, "b", b, "dec", dec)
	case 2:
		g := mathutil.GCD(sum, 100)
		return w10Problem(kp, fmt.Sprintf("Add %d/%d + 0.%02d, then reduce. What is the numerator?", a, b, dec), sum/g,
			fmt.Sprintf("The sum is %d/100, and dividing both parts by %d leaves %d/%d.", sum, g, sum/g, 100/g), "a", a, "b", b, "dec", dec)
	default:
		code := 0
		if sum > 100 {
			code = 1
		}
		return w10Problem(kp, fmt.Sprintf("Is %d/%d + 0.%02d greater than 1? Enter 1 for yes and 0 for no.", a, b, dec), code,
			fmt.Sprintf("The sum is %d/100, which is %s 1.", sum, map[bool]string{true: "greater than", false: "not greater than"}[sum > 100]), "a", a, "b", b, "dec", dec)
	}
}

func w10MissingMinuend(kp int) generator.Problem {
	switch kp {
	case 1:
		b := w10R(2, 9)
		a := w10R(1, b-1)
		c := w10R(1, b-1)
		return w10Problem(kp, fmt.Sprintf("Solve x - %d/%d = %d/%d. What is the numerator of x over %d?", a, b, c, b, b), a+c,
			fmt.Sprintf("Add the subtrahend and the difference: %d plus %d is %d.", a, c, a+c), "a", a, "b", b, "c", c)
	case 2:
		b := w10R(2, 8)
		d := w10R(2, 8)
		a := w10R(1, b-1)
		c := w10R(1, d-1)
		return w10Problem(kp, fmt.Sprintf("Solve x - %d/%d = %d/%d. What is the numerator of x over %d?", a, b, c, d, b*d), a*d+c*b,
			fmt.Sprintf("Convert both fractions to denominator %d and add: %d plus %d is %d.", b*d, a*d, c*b, a*d+c*b), "a", a, "b", b, "c", c, "d", d)
	default:
		b := w10R(2, 9)
		a := w10R(1, b-1)
		w := w10R(1, 6)
		return w10Problem(kp, fmt.Sprintf("Solve x - %d/%d = %d. What is the numerator of x over %d?", a, b, w, b), w*b+a,
			fmt.Sprintf("Add the whole number and the fraction: %d wholes is %d/%d, and adding %d gives %d.", w, w*b, b, a, w*b+a), "a", a, "b", b, "w", w)
	}
}

func w10MissingSubtrahend(kp int) generator.Problem {
	switch kp {
	case 1:
		b := w10R(3, 9)
		a := w10R(2, b-1)
		c := w10R(1, a-1)
		return w10Problem(kp, fmt.Sprintf("Solve %d/%d - x = %d/%d. What is the numerator of x over %d?", a, b, c, b, b), a-c,
			fmt.Sprintf("The denominators match, so subtract the difference numerator: %d minus %d is %d.", a, c, a-c), "a", a, "b", b, "c", c)
	case 2:
		b := w10R(2, 7)
		d := b + w10R(1, 4)
		a := w10R(2, b-1)
		c := w10R(1, d-1)
		num := a*d - c*b
		for num <= 0 {
			c = w10R(1, d-1)
			num = a*d - c*b
		}
		return w10Problem(kp, fmt.Sprintf("Solve %d/%d - x = %d/%d. What is the numerator of x over %d?", a, b, c, d, b*d), num,
			fmt.Sprintf("Common denominator %d gives %d minus %d, which is %d.", b*d, a*d, c*b, num), "a", a, "b", b, "c", c, "d", d)
	default:
		b := w10R(2, 9)
		a := w10R(2, b-1)
		w := w10R(1, 4)
		num := a - w*b
		for num <= 0 {
			w = w10R(0, 3)
			num = a - w*b
		}
		return w10Problem(kp, fmt.Sprintf("Solve %d/%d - x = %d. What is the numerator of x over %d?", a, b, w, b), num,
			fmt.Sprintf("Write %d as %d/%d, then subtract numerators: %d minus %d is %d.", w, w*b, b, a, w*b, num), "a", a, "b", b, "w", w)
	}
}

func w10EstimateDiff(kp int) generator.Problem {
	b := []int{2, 3, 4, 6}[w10R(0, 3)]
	d := []int{2, 3, 4, 6}[w10R(0, 3)]
	a := w10R(1, b-1)
	c := w10R(1, d-1)
	num := a*d - c*b
	if num < 0 {
		num = -num
	}
	den := b * d
	switch kp {
	case 1:
		code := 0
		if num*2 > den {
			code = 2
		} else if num*2 == den {
			code = 1
		}
		return w10Problem(kp, fmt.Sprintf("Is the difference %d/%d - %d/%d closest to 0, 1/2 or 1? Enter 0, 1 or 2.", a, b, c, d), code,
			fmt.Sprintf("The absolute difference is %d/%d, which is closest to the benchmark at code %d.", num, den, code), "a", a, "b", b, "c", c, "d", d)
	case 2:
		q := (num + den/2) / den
		return w10Problem(kp, fmt.Sprintf("Estimate the difference %d/%d - %d/%d as a whole number.", a, b, c, d), q,
			fmt.Sprintf("The absolute difference is %d/%d, which rounds to %d.", num, den, q), "a", a, "b", b, "c", c, "d", d)
	default:
		code := 0
		if num*2 < den {
			code = 1
		}
		return w10Problem(kp, fmt.Sprintf("Is the difference %d/%d - %d/%d less than 1/2? Enter 1 for yes and 0 for no.", a, b, c, d), code,
			fmt.Sprintf("The absolute difference is %d/%d, which is %s 1/2.", num, den, map[bool]string{true: "less than", false: "not less than"}[num*2 < den]), "a", a, "b", b, "c", c, "d", d)
	}
}

func w10SubRegroupWhole(kp int) generator.Problem {
	w := w10R(2, 9)
	b := w10R(2, 9)
	a := w10R(1, b-1)
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("Rewrite %d - %d/%d by borrowing one whole. What is the new whole-number part?", w, a, b), w-1,
			fmt.Sprintf("Borrow one from %d and turn it into %d/%d, leaving the whole part %d.", w, b, b, w-1), "w", w, "a", a, "b", b)
	case 2:
		return w10Problem(kp, fmt.Sprintf("After borrowing one whole, what is the numerator of the fractional part of %d - %d/%d over %d?", w, a, b, b), b-a,
			fmt.Sprintf("The borrowed whole becomes %d/%d, and subtracting %d leaves %d.", b, b, a, b-a), "w", w, "a", a, "b", b)
	default:
		return w10Problem(kp, fmt.Sprintf("Write %d - %d/%d as an improper fraction over %d. What is the numerator?", w, a, b, b), w*b-a,
			fmt.Sprintf("Convert %d to %d/%d, then subtract %d to get %d.", w, w*b, b, a, w*b-a), "w", w, "a", a, "b", b)
	}
}

func w10SubNegative(kp int) generator.Problem {
	b := w10R(2, 7)
	d := b + w10R(1, 4)
	a := w10R(1, b-1)
	c := w10R(1, d-1)
	for a*d >= c*b {
		c = w10R(1, d-1)
	}
	num := a*d - c*b
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("Compute %d/%d - %d/%d over the common denominator %d. What is the (signed) numerator?", a, b, c, d, b*d), num,
			fmt.Sprintf("Convert both fractions: %d minus %d is %d.", a*d, c*b, num), "a", a, "b", b, "c", c, "d", d)
	case 2:
		return w10Problem(kp, fmt.Sprintf("Is %d/%d - %d/%d negative? Enter 1 for yes and 0 for no.", a, b, c, d), 1,
			fmt.Sprintf("The numerator %d is smaller than %d, so the difference is negative.", a*d, c*b), "a", a, "b", b, "c", c, "d", d)
	default:
		return w10Problem(kp, fmt.Sprintf("What is the magnitude of %d/%d - %d/%d as a numerator over %d?", a, b, c, d, b*d), c*b-a*d,
			fmt.Sprintf("The signed numerator is %d, so the magnitude is %d.", num, c*b-a*d), "a", a, "b", b, "c", c, "d", d)
	}
}

func w10SubThree(kp int) generator.Problem {
	b := w10R(4, 9)
	a := w10R(3, b-1)
	c := w10R(1, a-2)
	e := w10R(1, a-c-1)
	num := a - c - e
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("Compute %d/%d - %d/%d - %d/%d. What is the numerator over %d?", a, b, c, b, e, b, b), num,
			fmt.Sprintf("The denominators match, so subtract in order: %d minus %d minus %d is %d.", a, c, e, num), "a", a, "b", b, "c", c, "e", e)
	case 2:
		return w10Problem(kp, fmt.Sprintf("Compute %d/%d - %d/%d - %d/%d. What is the whole-number part?", a, b, c, b, e, b), num/b,
			fmt.Sprintf("The numerator %d contains %d whole units of size %d.", num, num/b, b), "a", a, "b", b, "c", c, "e", e)
	default:
		return w10Problem(kp, fmt.Sprintf("After %d/%d - %d/%d - %d/%d, what is the remainder numerator over %d?", a, b, c, b, e, b, b), num%b,
			fmt.Sprintf("The numerator %d leaves remainder %d after removing whole units of size %d.", num, num%b, b), "a", a, "b", b, "c", c, "e", e)
	}
}

func w10SubDecimal(kp int) generator.Problem {
	b := []int{2, 4, 5, 10, 20, 25, 50}[w10R(0, 6)]
	a := w10R(1, b-1)
	fracHundred := a * 100 / b
	dec := w10R(1, fracHundred-1)
	diff := fracHundred - dec
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("Compute %d/%d - 0.%02d. Give the numerator over 100.", a, b, dec), diff,
			fmt.Sprintf("The fraction is %d/100, and subtracting %d leaves %d.", fracHundred, dec, diff), "a", a, "b", b, "dec", dec)
	case 2:
		g := mathutil.GCD(diff, 100)
		return w10Problem(kp, fmt.Sprintf("Compute %d/%d - 0.%02d and reduce. What is the numerator?", a, b, dec), diff/g,
			fmt.Sprintf("The difference is %d/100, and dividing both parts by %d leaves %d/%d.", diff, g, diff/g, 100/g), "a", a, "b", b, "dec", dec)
	default:
		return w10Problem(kp, fmt.Sprintf("Is %d/%d - 0.%02d positive? Enter 1 for yes and 0 for no.", a, b, dec), 1,
			fmt.Sprintf("The fraction is %d/100 and the decimal is %d/100, so the difference is positive.", fracHundred, dec), "a", a, "b", b, "dec", dec)
	}
}

func w10SubDistance(kp int) generator.Problem {
	a := w10R(1, 9)
	b := w10R(2, 9)
	c := w10R(1, 9)
	d := w10R(2, 9)
	l := lcm(b, d)
	num := a*(l/b) - c*(l/d)
	if num < 0 {
		num = -num
	}
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("What is the least common denominator of %d/%d and %d/%d?", a, b, c, d), l,
			fmt.Sprintf("The least common multiple of %d and %d is %d.", b, d, l), "a", a, "b", b, "c", c, "d", d)
	case 2:
		return w10Problem(kp, fmt.Sprintf("What is the distance between %d/%d and %d/%d as a numerator over %d?", a, b, c, d, l), num,
			fmt.Sprintf("Convert both to denominator %d and subtract the smaller numerator from the larger: the distance numerator is %d.", l, num), "a", a, "b", b, "c", c, "d", d)
	default:
		code := 0
		if num*2 > l {
			code = 1
		}
		return w10Problem(kp, fmt.Sprintf("Is the distance between %d/%d and %d/%d greater than 1/2? Enter 1 for yes and 0 for no.", a, b, c, d), code,
			fmt.Sprintf("The distance is %d/%d, which is %s 1/2.", num, l, map[bool]string{true: "greater than", false: "not greater than"}[num*2 > l]), "a", a, "b", b, "c", c, "d", d)
	}
}
