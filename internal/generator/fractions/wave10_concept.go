package fractions

// Wave 10 concept and equivalence generators: density, whole numbers as
// fractions, proper and improper forms, equal parts, unit fractions, cross
// multiplication, lowest terms, equivalent chains, benchmarks and percent
// equivalents.

import (
	"fmt"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

func w10Density(kp int) generator.Problem {
	d := w10R(4, 16)
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("How many fractions with denominator %d lie strictly between 0 and 1?", d), d-1,
			fmt.Sprintf("The numerators 1 through %d give %d fractions strictly between 0 and 1.", d-1, d-1), "d", d)
	case 2:
		n := d - 1 - d/2
		return w10Problem(kp, fmt.Sprintf("How many fractions with denominator %d lie strictly between 1/2 and 1?", d), n,
			fmt.Sprintf("The numerators above half are %d through %d, which is %d fractions.", d/2+1, d-1, n), "d", d)
	default:
		d = 4 * w10R(2, 6)
		n := d/2 - 1
		return w10Problem(kp, fmt.Sprintf("How many fractions with denominator %d lie strictly between 1/4 and 3/4?", d), n,
			fmt.Sprintf("The numerators from %d through %d lie strictly inside the interval, which is %d fractions.", d/4+1, 3*d/4-1, n), "d", d)
	}
}

func w10WholeFraction(kp int) generator.Problem {
	w := w10R(2, 9)
	d := w10R(2, 9)
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("Write %d as a fraction with denominator %d. What is the numerator?", w, d), w*d,
			fmt.Sprintf("Multiply the whole number by the denominator: %d times %d gives %d.", w, d, w*d), "w", w, "d", d)
	case 2:
		n := w * d
		return w10Problem(kp, fmt.Sprintf("The fraction %d/d equals %d. What is d?", n, w), d,
			fmt.Sprintf("Divide the numerator by the whole number: %d divided by %d is %d.", n, w, d), "w", w, "n", n)
	default:
		return w10Problem(kp, fmt.Sprintf("How many %d-ths are in %d wholes?", d, w), w*d,
			fmt.Sprintf("Each whole holds %d of them, so %d wholes hold %d times %d, which is %d.", d, w, w, d, w*d), "w", w, "d", d)
	}
}

func w10ProperImproper(kp int) generator.Problem {
	switch kp {
	case 1:
		d := w10R(2, 12)
		n := w10R(1, 2*d)
		if n == d {
			n++
		}
		code := 0
		word := "improper"
		if n < d {
			code = 1
			word = "proper"
		}
		return w10Problem(kp, fmt.Sprintf("Enter 1 if %d/%d is a proper fraction and 0 if it is improper.", n, d), code,
			fmt.Sprintf("The numerator %d is %s the denominator %d, so the fraction is %s.", n, map[bool]string{true: "less than", false: "not less than"}[n < d], d, word), "n", n, "d", d)
	case 2:
		d := w10R(2, 9)
		m := w10R(1, 4)
		r := w10R(0, d-1)
		n := m*d + r
		return w10Problem(kp, fmt.Sprintf("What is the whole-number part of the improper fraction %d/%d?", n, d), m,
			fmt.Sprintf("%d divided by %d is %d with remainder %d, so the whole part is %d.", n, d, m, r, m), "n", n, "d", d)
	default:
		b := w10R(2, 9)
		a := w10R(1, b-1)
		w := w10R(1, 6)
		n := w*b + a
		return w10Problem(kp, fmt.Sprintf("Write %d %d/%d as an improper fraction. What is the numerator over %d?", w, a, b, b), n,
			fmt.Sprintf("Multiply %d by %d to get %d, then add the numerator %d to get %d.", w, b, w*b, a, n), "w", w, "a", a, "b", b)
	}
}

func w10ZeroDenominator(kp int) generator.Problem {
	switch kp {
	case 1:
		zeros := 0
		parts := ""
		for i := 0; i < 3; i++ {
			if w10R(1, 2) == 1 {
				zeros++
				parts += fmt.Sprintf("x%d/0", i+1)
			} else {
				parts += fmt.Sprintf("x%d/%d", i+1, w10R(2, 9))
			}
			if i < 2 {
				parts += ", "
			}
		}
		return w10Problem(kp, fmt.Sprintf("How many of these expressions are undefined: %s?", parts), zeros,
			fmt.Sprintf("A fraction is undefined exactly when its denominator is zero, and %d of the three denominators are zero.", zeros), "zeros", zeros)
	case 2:
		k := w10R(2, 12)
		return w10Problem(kp, fmt.Sprintf("For which value of x is x/(x-%d) undefined?", k), k,
			fmt.Sprintf("The denominator x-%d is zero when x equals %d.", k, k), "k", k)
	default:
		a := w10R(1, 12)
		b := w10R(2, 12)
		for b == a {
			b = w10R(2, 12)
		}
		return w10Problem(kp, fmt.Sprintf("Solve (x-%d)/(x-%d) = 0. Enter the value of x.", a, b), a,
			fmt.Sprintf("A fraction is zero when its numerator is zero and its denominator is not. The numerator is zero at x = %d, and the denominator there is %d.", a, a-b), "a", a, "b", b)
	}
}

func w10EqualParts(kp int) generator.Problem {
	switch kp {
	case 1:
		d := w10R(2, 12)
		return w10Problem(kp, fmt.Sprintf("A whole is cut into %d equal parts. How many parts make the whole?", d), d,
			fmt.Sprintf("All %d equal parts together rebuild the whole.", d), "d", d)
	case 2:
		d := w10R(3, 12)
		s := w10R(1, d-1)
		return w10Problem(kp, fmt.Sprintf("A strip is cut into %d equal parts and %d are shaded. How many parts are unshaded?", d, s), d-s,
			fmt.Sprintf("The unshaded count is the total %d minus the shaded %d, which is %d.", d, s, d-s), "d", d, "s", s)
	default:
		n := w10R(2, 9)
		w := w10R(2, 6)
		return w10Problem(kp, fmt.Sprintf("How many pieces of size 1/%d fit into %d wholes?", n, w), w*n,
			fmt.Sprintf("Each whole holds %d pieces, so %d wholes hold %d times %d, which is %d pieces.", n, w, w, n, w*n), "n", n, "w", w)
	}
}

func w10UnitMeaning(kp int) generator.Problem {
	n := w10R(2, 9)
	k := w10R(2, 9)
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("What is 1/%d of %d?", n, n*k), k,
			fmt.Sprintf("Dividing %d objects into %d equal groups puts %d in each group.", n*k, n, k), "n", n, "total", n*k)
	case 2:
		n = w10R(3, 12)
		k = w10R(3, 12)
		return w10Problem(kp, fmt.Sprintf("What is 1/%d of %d?", n, n*k), k,
			fmt.Sprintf("One part out of %d equal parts of %d is %d.", n, n*k, k), "n", n, "total", n*k)
	default:
		w := w10R(2, 6)
		return w10Problem(kp, fmt.Sprintf("How many unit fractions 1/%d are needed to make %d wholes?", n, w), w*n,
			fmt.Sprintf("Each whole needs %d copies, so %d wholes need %d copies.", n, w, w*n), "n", n, "w", w)
	}
}

func w10CrossCheck(kp int) generator.Problem {
	a := w10R(1, 9)
	b := w10R(2, 9)
	k := w10R(2, 6)
	switch kp {
	case 1:
		equal := w10R(0, 1) == 1
		c := a * k
		d := b * k
		code := 1
		if !equal {
			d++
			code = 0
		}
		return w10Problem(kp, fmt.Sprintf("Do %d/%d and %d/%d name the same fraction? Enter 1 for yes and 0 for no.", a, b, c, d), code,
			fmt.Sprintf("Cross multiply: %d times %d is %d, and %d times %d is %d. The cross products are %s.", a, d, a*d, c, b, c*b, map[bool]string{true: "equal, so the fractions are equivalent", false: "different, so the fractions differ"}[code == 1]), "a", a, "b", b, "c", c, "d", d)
	case 2:
		c := a * k
		d := b*k + 1
		return w10Problem(kp, fmt.Sprintf("For %d/%d and %d/%d, what is the cross product a times d?", a, b, c, d), a*d,
			fmt.Sprintf("Multiply the numerator of the first fraction by the denominator of the second: %d times %d is %d.", a, d, a*d), "a", a, "b", b, "c", c, "d", d)
	default:
		c := a * k
		d := b * k
		return w10Problem(kp, fmt.Sprintf("Find d so that %d/%d = %d/d.", a, b, c), d,
			fmt.Sprintf("The numerator was scaled by %d, so the denominator is %d times %d, which is %d.", k, b, k, d), "a", a, "b", b, "c", c)
	}
}

func w10SimplifyGCF(kp int) generator.Problem {
	g := w10R(2, 9)
	p := w10R(1, 9)
	q := w10R(2, 12)
	for mathutil.GCD(p, q) != 1 || p >= q {
		p = w10R(1, 9)
		q = w10R(2, 12)
	}
	n, d := g*p, g*q
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("What is the greatest common factor of %d and %d?", n, d), g,
			fmt.Sprintf("Both numbers are multiples of %d, and the remaining factors %d and %d are coprime, so the greatest common factor is %d.", g, p, q, g), "n", n, "d", d)
	case 2:
		return w10Problem(kp, fmt.Sprintf("Reduce %d/%d to lowest terms. What is the numerator?", n, d), p,
			fmt.Sprintf("Divide both parts by the greatest common factor %d: %d becomes %d and %d becomes %d.", g, n, p, d, q), "n", n, "d", d)
	default:
		return w10Problem(kp, fmt.Sprintf("Reduce %d/%d to lowest terms. What is the denominator?", n, d), q,
			fmt.Sprintf("Divide both parts by %d: the numerator becomes %d and the denominator becomes %d.", g, p, q), "n", n, "d", d)
	}
}

func w10LowestTerms(kp int) generator.Problem {
	g := w10R(1, 9)
	p := w10R(1, 9)
	q := w10R(2, 12)
	for mathutil.GCD(p, q) != 1 || p >= q {
		p = w10R(1, 9)
		q = w10R(2, 12)
	}
	n, d := g*p, g*q
	switch kp {
	case 1:
		code := 0
		if g == 1 {
			code = 1
		}
		return w10Problem(kp, fmt.Sprintf("Enter 1 if %d/%d is already in lowest terms and 0 otherwise.", n, d), code,
			fmt.Sprintf("The greatest common factor of %d and %d is %d, so the fraction is %s in lowest terms.", n, d, g, map[bool]string{true: "already", false: "not"}[g == 1]), "n", n, "d", d)
	case 2:
		return w10Problem(kp, fmt.Sprintf("What is the greatest common factor of the numerator and denominator of %d/%d?", n, d), g,
			fmt.Sprintf("The shared factor is %d, leaving %d and %d, which share no further factor.", g, p, q), "n", n, "d", d)
	default:
		return w10Problem(kp, fmt.Sprintf("Write %d/%d in lowest terms. What is the numerator?", n, d), p,
			fmt.Sprintf("Dividing numerator and denominator by %d leaves %d/%d.", g, p, q), "n", n, "d", d)
	}
}

func w10EquivChain(kp int) generator.Problem {
	a := w10R(1, 9)
	b := w10R(2, 9)
	k := w10R(2, 6)
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("Fill in the numerator: %d/%d = ?/%d.", a, b, b*k), a*k,
			fmt.Sprintf("The denominator grew by a factor of %d, so the numerator grows by the same factor: %d times %d is %d.", k, a, k, a*k), "a", a, "b", b, "k", k)
	case 2:
		return w10Problem(kp, fmt.Sprintf("Fill in the denominator: %d/%d = %d/?.", a, b, a*k), b*k,
			fmt.Sprintf("The numerator grew by a factor of %d, so the denominator grows by the same factor: %d times %d is %d.", k, b, k, b*k), "a", a, "b", b, "k", k)
	default:
		k2 := k + w10R(1, 4)
		return w10Problem(kp, fmt.Sprintf("In the chain %d/%d = ?/%d = ?/%d, what is the first missing numerator?", a, b, b*k, b*k2), a*k,
			fmt.Sprintf("The first equivalent fraction multiplies both parts by %d, so the numerator is %d times %d, which is %d.", k, a, k, a*k), "a", a, "b", b, "k", k, "k2", k2)
	}
}

func w10MissingNumber(kp int) generator.Problem {
	a := w10R(1, 9)
	b := w10R(2, 9)
	k := w10R(2, 7)
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("Find the missing numerator: %d/%d = x/%d.", a, b, b*k), a*k,
			fmt.Sprintf("Multiply numerator and denominator by %d: %d times %d is %d.", k, a, k, a*k), "a", a, "b", b, "k", k)
	case 2:
		return w10Problem(kp, fmt.Sprintf("Find the missing denominator: %d/%d = %d/x.", a, b, a*k), b*k,
			fmt.Sprintf("Multiply numerator and denominator by %d: %d times %d is %d.", k, b, k, b*k), "a", a, "b", b, "k", k)
	default:
		k = w10R(3, 9)
		return w10Problem(kp, fmt.Sprintf("Find the missing numerator: %d/%d = x/%d.", a, b, b*k), a*k,
			fmt.Sprintf("The denominator is %d times %d, so the numerator is %d times %d, which is %d.", k, b, k, a, a*k), "a", a, "b", b, "k", k)
	}
}

func w10BenchmarkChoice(kp int) generator.Problem {
	switch kp {
	case 1:
		d := w10R(4, 12)
		n := w10R(1, d-1)
		for 2*n == d {
			n = w10R(1, d-1)
		}
		code := 0
		word := "0"
		if 2*n > d {
			code = 1
			word = "1"
		}
		return w10Problem(kp, fmt.Sprintf("Is %d/%d closer to 0 or to 1? Enter 1 for 1 and 0 for 0.", n, d), code,
			fmt.Sprintf("The numerator is %s half the denominator, so the fraction is closer to %s.", map[bool]string{true: "greater than", false: "less than"}[2*n > d], word), "n", n, "d", d)
	case 2:
		d := w10R(4, 12)
		n := w10R(1, d-1)
		d0 := n
		d1 := 2*n - d
		if d1 < 0 {
			d1 = -d1
		}
		d2 := d - n
		code := 0
		if d1 <= d0 && d1 <= d2 {
			code = 1
		} else if d2 < d0 && d2 < d1 {
			code = 2
		}
		return w10Problem(kp, fmt.Sprintf("Which benchmark is %d/%d closest to? Enter 0 for 0, 1 for 1/2 and 2 for 1.", n, d), code,
			fmt.Sprintf("Compare the distances: to 0 is %d/%d, to 1/2 is %d/%d, and to 1 is %d/%d. The smallest is the answer.", d0, d, d1, d, d2, d), "n", n, "d", d)
	default:
		n := w10R(1, 7)
		gap := n - 4
		if gap < 0 {
			gap = -gap
		}
		return w10Problem(kp, fmt.Sprintf("How many eighths separate %d/8 from 1/2?", n), gap,
			fmt.Sprintf("Half is 4/8, and the distance from %d/8 to 4/8 is %d eighths.", n, gap), "n", n)
	}
}

func w10DecimalOrder(kp int) generator.Problem {
	a := w10R(1, 9)
	b := w10R(1, 9)
	c := w10R(1, 9)
	switch kp {
	case 1:
		count := 0
		for _, v := range []int{a, b, c} {
			if v < 5 {
				count++
			}
		}
		return w10Problem(kp, fmt.Sprintf("How many of %d/10, %d/10 and %d/10 are less than 0.5?", a, b, c), count,
			fmt.Sprintf("A tenth is below 0.5 exactly when its numerator is below 5, and %d of the numerators are.", count), "a", a, "b", b, "c", c)
	case 2:
		idx := 1
		if b > a {
			idx = 2
		}
		if c > []int{a, b}[idx-1] {
			idx = 3
		}
		return w10Problem(kp, fmt.Sprintf("Which of %d/10, %d/10 and %d/10 is largest? Enter 1, 2 or 3.", a, b, c), idx,
			fmt.Sprintf("Compare the numerators: the largest is position %d.", idx), "a", a, "b", b, "c", c)
	default:
		vals := []int{a, b, c}
		rank := 1
		for _, v := range vals {
			if v < vals[2] {
				rank++
			}
		}
		return w10Problem(kp, fmt.Sprintf("Arrange %d/10, %d/10 and 0.%d in increasing order. What is the position of 0.%d?", a, b, c, c), rank,
			fmt.Sprintf("As tenths the values are %d, %d and %d, and %d is position %d from the bottom.", a, b, c, c, rank), "a", a, "b", b, "c", c)
	}
}

func w10PercentBench(kp int) generator.Problem {
	denoms := []int{2, 4, 5, 10, 20, 25, 50}
	b := denoms[w10R(0, len(denoms)-1)]
	a := w10R(1, b-1)
	switch kp {
	case 1:
		return w10Problem(kp, fmt.Sprintf("What percent is %d/%d?", a, b), a*100/b,
			fmt.Sprintf("Multiply by 100: %d divided by %d, then times 100, is %d percent.", a, b, a*100/b), "a", a, "b", b)
	case 2:
		p := 5 * w10R(1, 19)
		g := mathutil.GCD(p, 100)
		return w10Problem(kp, fmt.Sprintf("Write %d percent as a fraction in lowest terms. What is the numerator?", p), p/g,
			fmt.Sprintf("%d percent is %d/100, and dividing both parts by %d leaves %d/%d.", p, p, g, p/g, 100/g), "p", p)
	default:
		return w10Problem(kp, fmt.Sprintf("Express %d/%d as a percentage. How many percent is it?", a, b), a*100/b,
			fmt.Sprintf("Scale to hundredths: %d/%d is %d/100, which is %d percent.", a, b, a*100/b, a*100/b), "a", a, "b", b)
	}
}
