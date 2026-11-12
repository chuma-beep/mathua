package complex

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

// Wave 19 adds Gaussian arithmetic and plane geometry. Every family publishes
// the drawn parameters in Facts, so wave19_test.go can recompute the answer
// instead of trusting the generator.

func registerWave19(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"complex.gauss.norm":         "gauss_norm",
		"complex.gauss.units":        "gauss_units",
		"complex.gauss.division":     "gauss_division",
		"complex.gauss.gcd":          "gauss_gcd",
		"complex.gauss.primes":       "gauss_primes",
		"complex.gauss.two_squares":  "gauss_two_squares",
		"complex.gauss.associates":   "gauss_associates",
		"complex.gauss.pythagorean":  "gauss_pythagorean",
		"complex.gauss.theta":        "gauss_theta",
		"complex.gauss.eisenstein":   "gauss_eisenstein",
		"complex.geom.rotation":      "geom_rotation",
		"complex.geom.area":          "geom_area",
		"complex.geom.circle":        "geom_circle",
		"complex.geom.reflection":    "geom_reflection",
		"complex.geom.spiral":        "geom_spiral",
		"complex.geom.ptolemy":       "geom_ptolemy",
		"complex.geom.triangle":      "geom_triangle",
		"complex.geom.centroid":      "geom_centroid",
		"complex.geom.perpendicular": "geom_perpendicular",
		"complex.geom.lattice":       "geom_lattice",
	} {
		reg.Register(id, &w19Gen{kind: kind})
	}
	registerWave19Rest(reg)
}

type w19Gen struct{ kind string }

func (g *w19Gen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := ctx.KPOrdinal
	if o < 1 || o > 3 {
		o = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "gauss_norm":
		return w19GaussNorm(o)
	case "gauss_units":
		return w19GaussUnits(o)
	case "gauss_division":
		return w19GaussDivision(o)
	case "gauss_gcd":
		return w19GaussGCD(o)
	case "gauss_primes":
		return w19GaussPrimes(o)
	case "gauss_two_squares":
		return w19GaussTwoSquares(o)
	case "gauss_associates":
		return w19GaussAssociates(o)
	case "gauss_pythagorean":
		return w19GaussPythagorean(o)
	case "gauss_theta":
		return w19GaussTheta(o)
	case "gauss_eisenstein":
		return w19GaussEisenstein(o)
	case "geom_rotation":
		return w19GeomRotation(o)
	case "geom_area":
		return w19GeomArea(o)
	case "geom_circle":
		return w19GeomCircle(o)
	case "geom_reflection":
		return w19GeomReflection(o)
	case "geom_spiral":
		return w19GeomSpiral(o)
	case "geom_ptolemy":
		return w19GeomPtolemy(o)
	case "geom_triangle":
		return w19GeomTriangle(o)
	case "geom_centroid":
		return w19GeomCentroid(o)
	case "geom_perpendicular":
		return w19GeomPerpendicular(o)
	case "geom_lattice":
		return w19GeomLattice(o)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

func w19N(lo, hi int) int {
	return lo + rand.Intn(hi-lo+1)
}

func w19P(o int, q, a, ex string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}

func w19I(m map[string]string, k string, v int) { m[k] = strconv.Itoa(v) }

func w19I64(m map[string]string, k string, v int64) { m[k] = strconv.FormatInt(v, 10) }

func w19IsPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func w19Prime(residue, lo, hi int, allowTwo bool) int {
	var pool []int
	for n := lo; n <= hi; n++ {
		if !w19IsPrime(n) {
			continue
		}
		if n == 2 && !allowTwo {
			continue
		}
		if residue == 0 || n%4 == residue {
			pool = append(pool, n)
		}
	}
	return pool[rand.Intn(len(pool))]
}

// w19DivisorCounts returns the number of divisors of n congruent to 1 and to 3
// modulo four.
func w19DivisorCounts(n int) (int, int) {
	var one, three int
	for d := 1; d <= n; d++ {
		if n%d != 0 {
			continue
		}
		switch d % 4 {
		case 1:
			one++
		case 3:
			three++
		}
	}
	return one, three
}

func w19MinSquares(n int) int {
	for a := 0; a*a <= n; a++ {
		if a*a == n {
			return 1
		}
	}
	for a := 0; a*a <= n; a++ {
		for b := 0; a*a+b*b <= n; b++ {
			if a*a+b*b == n {
				return 2
			}
		}
	}
	for a := 0; a*a <= n; a++ {
		for b := 0; a*a+b*b <= n; b++ {
			for c := 0; a*a+b*b+c*c <= n; c++ {
				if a*a+b*b+c*c == n {
					return 3
				}
			}
		}
	}
	return 4
}

func w19GaussNorm(o int) generator.Problem {
	f := map[string]string{}
	switch o {
	case 1:
		a, b := w19N(1, 15), w19N(1, 15)
		w19I(f, "a", a)
		w19I(f, "b", b)
		return w19P(o, fmt.Sprintf("For z = %d+%di, what is N(z)? (enter an integer)", a, b),
			strconv.Itoa(a*a+b*b), fmt.Sprintf("N(z)=a^2+b^2=%d^2+%d^2=%d.", a, b, a*a+b*b), f)
	case 2:
		a, b, c, d := w19N(1, 9), w19N(1, 9), w19N(1, 9), w19N(1, 9)
		w19I(f, "a", a)
		w19I(f, "b", b)
		w19I(f, "c", c)
		w19I(f, "d", d)
		ans := (a*a + b*b) * (c*c + d*d)
		return w19P(o, fmt.Sprintf("What is N((%d+%di)(%d+%di))? (enter an integer)", a, b, c, d),
			strconv.Itoa(ans), "The norm is multiplicative, so multiply the two norms.", f)
	default:
		a, b := w19N(1, 12), w19N(1, 12)
		w19I(f, "a", a)
		w19I(f, "b", b)
		return w19P(o, fmt.Sprintf("What is N(z)+N(conj z) for z = %d+%di? (enter an integer)", a, b),
			strconv.Itoa(2*(a*a+b*b)), "Conjugation preserves the norm, so the sum is twice the norm.", f)
	}
}

func w19GaussUnits(o int) generator.Problem {
	a, b := w19N(1, 10), w19N(1, 10)
	k := rand.Intn(4)
	f := map[string]string{}
	w19I(f, "a", a)
	w19I(f, "b", b)
	w19I(f, "k", k)
	re := [4]int{a, -b, -a, b}
	im := [4]int{b, a, -b, -a}
	switch o {
	case 1:
		return w19P(o, fmt.Sprintf("For z = %d+%di, what is the real part of i^%d z? (enter an integer)", a, b, k),
			strconv.Itoa(re[k]), "Multiplication by i is a quarter turn, so apply it k times.", f)
	case 2:
		return w19P(o, fmt.Sprintf("For z = %d+%di, what is the imaginary part of i^%d z? (enter an integer)", a, b, k),
			strconv.Itoa(im[k]), "Track the coordinates through k quarter turns.", f)
	default:
		return w19P(o, fmt.Sprintf("What is the norm of i^%d(%d+%di)? (enter an integer)", k, a, b),
			strconv.Itoa(a*a+b*b), "A unit has norm one, so multiplication by it preserves the norm.", f)
	}
}

func w19GaussDivision(o int) generator.Problem {
	f := map[string]string{}
	switch o {
	case 1:
		n := w19N(2, 9)
		a := w19N(-20, 20)
		w19I(f, "n", n)
		w19I(f, "a", a)
		q := int(math.Round(float64(a) / float64(n)))
		return w19P(o, fmt.Sprintf("Round %d/%d to the nearest integer. (enter an integer)", a, n),
			strconv.Itoa(q), "Gaussian division rounds both coordinates to the nearest integer.", f)
	case 2:
		n := w19N(2, 9)
		a, b := w19N(-20, 20), w19N(-20, 20)
		w19I(f, "n", n)
		w19I(f, "a", a)
		w19I(f, "b", b)
		q := int(math.Round(float64(a) / float64(n)))
		rem := (a-q*n)*(a-q*n) + b*b
		return w19P(o, fmt.Sprintf("After rounding the quotient of %d+%di by %d, what is N(remainder)? (enter an integer)", a, b, n),
			strconv.Itoa(rem), "The remainder is z minus the rounded quotient times the divisor.", f)
	default:
		n := w19N(2, 40)
		w19I(f, "n", n)
		return w19P(o, fmt.Sprintf("For a divisor of norm %d, what is the largest integer bound on N(remainder)? (enter an integer)", n),
			strconv.Itoa(n/2), "The Gaussian remainder satisfies N(r) <= N(w)/2.", f)
	}
}

func w19GaussGCD(o int) generator.Problem {
	f := map[string]string{}
	switch o {
	case 1:
		a, b := w19N(1, 12), w19N(1, 12)
		w19I(f, "a", a)
		w19I(f, "b", b)
		return w19P(o, fmt.Sprintf("What is the norm of gcd(%d+%di, i(%d+%di))? (enter an integer)", a, b, a, b),
			strconv.Itoa(a*a+b*b), "z and iz are associates, so the gcd is z up to a unit.", f)
	case 2:
		p := w19Prime(3, 3, 199, false)
		w19I(f, "p", p)
		return w19P(o, fmt.Sprintf("What is the norm of gcd(%d, %d(1+i))? (enter an integer)", p, p),
			strconv.Itoa(p*p), "A prime 3 mod 4 is a Gaussian prime, so it is the gcd.", f)
	default:
		a, b := w19N(1, 12), w19N(1, 12)
		for w19GCD(a, b) != 1 {
			a, b = w19N(1, 12), w19N(1, 12)
		}
		w19I(f, "a", a)
		w19I(f, "b", b)
		return w19P(o, fmt.Sprintf("What is the norm of gcd(%d+%di, %d-%di)? (enter an integer)", a, b, a, b),
			"1", "Coprime coordinates make the two conjugates coprime, so the gcd is a unit of norm one.", f)
	}
}

func w19GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func w19GaussPrimes(o int) generator.Problem {
	f := map[string]string{}
	switch o {
	case 1:
		n := w19N(20, 200)
		w19I(f, "n", n)
		count := 0
		for p := 3; p <= n; p += 4 {
			if w19IsPrime(p) {
				count++
			}
		}
		return w19P(o, fmt.Sprintf("How many primes congruent to 3 mod 4 are at most %d? (enter an integer)", n),
			strconv.Itoa(count), "Count the primes p <= n with p = 3 mod 4.", f)
	case 2:
		p := w19Prime(0, 5, 200, false)
		w19I(f, "p", p)
		split := 0
		if p%4 == 1 {
			split = 1
		}
		return w19P(o, fmt.Sprintf("Does the prime %d split in the Gaussian integers? (enter 1 for yes, 0 for no)", p),
			strconv.Itoa(split), "A prime splits exactly when it is 1 mod 4.", f)
	default:
		p := w19Prime(1, 5, 200, false)
		w19I(f, "p", p)
		return w19P(o, fmt.Sprintf("Each Gaussian prime over %d has what norm? (enter an integer)", p),
			strconv.Itoa(p), "The two conjugate factors of a split prime both have norm p.", f)
	}
}

func w19GaussTwoSquares(o int) generator.Problem {
	f := map[string]string{}
	switch o {
	case 1:
		n := w19N(1, 200)
		w19I(f, "n", n)
		one, three := w19DivisorCounts(n)
		return w19P(o, fmt.Sprintf("What is r2(%d), the number of signed ordered representations as a sum of two squares? (enter an integer)", n),
			strconv.Itoa(4*(one-three)), "r2(n) is four times the excess of divisors 1 mod 4 over divisors 3 mod 4.", f)
	case 2:
		n := w19N(1, 200)
		w19I(f, "n", n)
		return w19P(o, fmt.Sprintf("What is the smallest number of squares needed to write %d? (enter an integer)", n),
			strconv.Itoa(w19MinSquares(n)), "Check one, two and three squares, then use Lagrange's four-square theorem.", f)
	default:
		p := w19Prime(1, 5, 150, false)
		w19I(f, "p", p)
		return w19P(o, fmt.Sprintf("How many signed ordered representations of %d as a sum of two squares are there? (enter an integer)", p),
			"8", "A split prime has one unordered positive representation, giving eight signed ordered ones.", f)
	}
}

func w19GaussAssociates(o int) generator.Problem {
	f := map[string]string{}
	a, b := w19N(1, 12), w19N(1, 12)
	w19I(f, "a", a)
	w19I(f, "b", b)
	switch o {
	case 1:
		return w19P(o, fmt.Sprintf("How many associates does %d+%di have? (enter an integer)", a, b),
			"4", "The four units give four associates.", f)
	case 2:
		return w19P(o, fmt.Sprintf("What is the norm of every associate of %d+%di? (enter an integer)", a, b),
			strconv.Itoa(a*a+b*b), "Units have norm one, so associates share the norm.", f)
	default:
		return w19P(o, fmt.Sprintf("What is the real part of (1+i)(%d+%di)? (enter an integer)", a, b),
			strconv.Itoa(a-b), "Multiply out and collect the real part.", f)
	}
}

func w19GaussPythagorean(o int) generator.Problem {
	a, b := w19N(1, 12), w19N(1, 12)
	f := map[string]string{}
	w19I(f, "a", a)
	w19I(f, "b", b)
	switch o {
	case 1:
		return w19P(o, fmt.Sprintf("The square of %d+%di gives a Pythagorean triple. What is its hypotenuse? (enter an integer)", a, b),
			strconv.Itoa(a*a+b*b), "The hypotenuse is N(a+bi).", f)
	case 2:
		leg := a*a - b*b
		if leg < 0 {
			leg = -leg
		}
		return w19P(o, fmt.Sprintf("The square of %d+%di gives a Pythagorean triple. What is the absolute difference of the legs? (enter an integer)", a, b),
			strconv.Itoa(leg), "One leg is |a^2-b^2| and the other is 2ab.", f)
	default:
		leg := a*a - b*b
		if leg < 0 {
			leg = -leg
		}
		perim := a*a + b*b + leg + 2*a*b
		return w19P(o, fmt.Sprintf("What is the perimeter of the triple from squaring %d+%di? (enter an integer)", a, b),
			strconv.Itoa(perim), "Add the two legs and the hypotenuse.", f)
	}
}

func w19GaussTheta(o int) generator.Problem {
	f := map[string]string{}
	switch o {
	case 1:
		k := w19N(1, 40)
		w19I(f, "k", k)
		return w19P(o, fmt.Sprintf("How many terms of the theta series have |n| <= %d? (enter an integer)", k),
			strconv.Itoa(2*k+1), "The indices run from -k to k inclusive.", f)
	case 2:
		n := w19N(1, 15)
		w19I(f, "n", n)
		return w19P(o, fmt.Sprintf("Which exponent of q does the term q^(n^2) carry when n = %d? (enter an integer)", n),
			strconv.Itoa(n*n), "The theta series sums q^{n^2}.", f)
	default:
		n := w19N(1, 200)
		w19I(f, "n", n)
		one, three := w19DivisorCounts(n)
		return w19P(o, fmt.Sprintf("What is the coefficient of q^%d in the square-lattice theta series? (enter an integer)", n),
			strconv.Itoa(4*(one-three)), "The coefficient is r2(n), computed from divisors modulo four.", f)
	}
}

func w19GaussEisenstein(o int) generator.Problem {
	f := map[string]string{}
	a, b := w19N(-12, 12), w19N(-12, 12)
	for a == 0 && b == 0 {
		a, b = w19N(-12, 12), w19N(-12, 12)
	}
	w19I(f, "a", a)
	w19I(f, "b", b)
	switch o {
	case 1:
		n := a*a - a*b + b*b
		return w19P(o, fmt.Sprintf("What is the norm of %d+%dω in the Eisenstein integers? (enter an integer)", a, b),
			strconv.Itoa(n), "The norm is a^2-ab+b^2.", f)
	case 2:
		return w19P(o, fmt.Sprintf("How many units does the Eisenstein ring have? (enter an integer) [a=%d, b=%d]", a, b),
			"6", "The six units are the sixth roots of unity.", f)
	default:
		aa := w19N(1, 15)
		w19I(f, "aa", aa)
		return w19P(o, fmt.Sprintf("What is the norm of %d+%dω? (enter an integer)", aa, aa),
			strconv.Itoa(aa*aa), "Substitute a=b in a^2-ab+b^2.", f)
	}
}

func w19GeomRotation(o int) generator.Problem {
	a, b := w19N(1, 12), w19N(1, 12)
	k := rand.Intn(5)
	f := map[string]string{}
	w19I(f, "a", a)
	w19I(f, "b", b)
	w19I(f, "k", k)
	re := [5]int{a, -b, -a, b, a}
	im := [5]int{b, a, -b, -a, b}
	switch o {
	case 1:
		return w19P(o, fmt.Sprintf("What is the real part of i^%d(%d+%di)? (enter an integer)", k, a, b),
			strconv.Itoa(re[k]), "Each factor of i is a quarter turn.", f)
	case 2:
		return w19P(o, fmt.Sprintf("What is the imaginary part of i^%d(%d+%di)? (enter an integer)", k, a, b),
			strconv.Itoa(im[k]), "Track the coordinates through k quarter turns.", f)
	default:
		return w19P(o, fmt.Sprintf("After %d quarter turns, what is the real part of the image of %d+%di? (enter an integer)", k, a, b),
			strconv.Itoa(re[k]), "Four turns restore the point, so reduce k modulo four.", f)
	}
}

func w19GeomArea(o int) generator.Problem {
	f := map[string]string{}
	switch o {
	case 1:
		a, b := w19N(-6, 6), w19N(-6, 6)
		c, d := w19N(-6, 6), w19N(-6, 6)
		e, g := w19N(-6, 6), w19N(-6, 6)
		w19I(f, "a", a)
		w19I(f, "b", b)
		w19I(f, "c", c)
		w19I(f, "d", d)
		w19I(f, "e", e)
		w19I(f, "f", g)
		det := (c-a)*(g-b) - (d-b)*(e-a)
		if det < 0 {
			det = -det
		}
		return w19P(o, fmt.Sprintf("What is twice the area of the triangle with vertices (%d,%d), (%d,%d), (%d,%d)? (enter an integer)", a, b, c, d, e, g),
			strconv.Itoa(det), "Twice the area is the absolute determinant of the two edge vectors.", f)
	case 2:
		p := 2 * w19N(1, 8)
		q := w19N(1, 10)
		w19I(f, "p", p)
		w19I(f, "q", q)
		return w19P(o, fmt.Sprintf("A right triangle has legs %d and %d. What is its area? (enter an integer)", p, q),
			strconv.Itoa(p*q/2), "The area is half the product of the legs.", f)
	default:
		a, b := w19N(-5, 5), w19N(-5, 5)
		c, d := w19N(-5, 5), w19N(-5, 5)
		e, g := a+2*(c-a), b+2*(d-b)
		w19I(f, "a", a)
		w19I(f, "b", b)
		w19I(f, "c", c)
		w19I(f, "d", d)
		w19I(f, "e", e)
		w19I(f, "f", g)
		return w19P(o, fmt.Sprintf("What is twice the area of the triangle with vertices (%d,%d), (%d,%d), (%d,%d)? (enter an integer)", a, b, c, d, e, g),
			"0", "The third point lies on the line through the first two, so the points are collinear.", f)
	}
}

func w19GeomCircle(o int) generator.Problem {
	f := map[string]string{}
	switch o {
	case 1:
		a, b := w19N(-8, 8), w19N(-8, 8)
		c, d := w19N(-8, 8), w19N(-8, 8)
		w19I(f, "a", a)
		w19I(f, "b", b)
		w19I(f, "c", c)
		w19I(f, "d", d)
		return w19P(o, fmt.Sprintf("What is the squared radius of the circle centred at (%d,%d) through (%d,%d)? (enter an integer)", a, b, c, d),
			strconv.Itoa((c-a)*(c-a)+(d-b)*(d-b)), "The squared radius is the squared distance to the centre.", f)
	case 2:
		p, q := w19N(-8, 8), w19N(-8, 8)
		r, s := w19N(-8, 8), w19N(-8, 8)
		w19I(f, "p", p)
		w19I(f, "q", q)
		w19I(f, "r", r)
		w19I(f, "s", s)
		return w19P(o, fmt.Sprintf("A diameter has endpoints (%d,%d) and (%d,%d). What is twice the real coordinate of the centre? (enter an integer)", p, q, r, s),
			strconv.Itoa(p+r), "The centre is the midpoint, so its real coordinate is (p+r)/2.", f)
	default:
		a, b := w19N(-8, 8), w19N(-8, 8)
		c, d := w19N(-8, 8), w19N(-8, 8)
		w19I(f, "a", a)
		w19I(f, "b", b)
		w19I(f, "c", c)
		w19I(f, "d", d)
		return w19P(o, fmt.Sprintf("What is |z-c|^2 for z = (%d,%d) and c = (%d,%d)? (enter an integer)", c, d, a, b),
			strconv.Itoa((c-a)*(c-a)+(d-b)*(d-b)), "Square the coordinate differences and add.", f)
	}
}

func w19GeomReflection(o int) generator.Problem {
	a, b := w19N(1, 12), w19N(1, 12)
	f := map[string]string{}
	w19I(f, "a", a)
	w19I(f, "b", b)
	switch o {
	case 1:
		return w19P(o, fmt.Sprintf("What is the imaginary part of the reflection of %d+%di in the real axis? (enter an integer)", a, b),
			strconv.Itoa(-b), "Reflection in the real axis conjugates the point.", f)
	case 2:
		return w19P(o, fmt.Sprintf("What is the real part of the reflection of %d+%di in the imaginary axis? (enter an integer)", a, b),
			strconv.Itoa(-a), "Reflection in the imaginary axis sends z to -conj(z).", f)
	default:
		return w19P(o, fmt.Sprintf("What is the real part of the reflection of %d+%di in the line y=x? (enter an integer)", a, b),
			strconv.Itoa(b), "Reflection in y=x swaps the coordinates.", f)
	}
}

func w19GeomSpiral(o int) generator.Problem {
	f := map[string]string{}
	switch o {
	case 1:
		p, q := w19N(1, 9), w19N(1, 9)
		k := w19N(1, 5)
		w19I(f, "p", p)
		w19I(f, "q", q)
		w19I(f, "k", k)
		return w19P(o, fmt.Sprintf("The segment from 0 to %d+%di is multiplied by %di. What is the scale factor? (enter an integer)", p, q, k),
			strconv.Itoa(k), "Multiplication by ki scales by |k|.", f)
	case 2:
		k := w19N(1, 15)
		w19I(f, "k", k)
		return w19P(o, fmt.Sprintf("Multiplication by %di turns the plane by how many degrees? (enter an integer)", k),
			strconv.Itoa(90*k), "Each factor of i is a quarter turn, so k factors give 90k degrees.", f)
	default:
		k1, k2 := w19N(1, 4), w19N(1, 4)
		w19I(f, "k1", k1)
		w19I(f, "k2", k2)
		return w19P(o, fmt.Sprintf("A spiral of scale %d is followed by one of scale %d. What is the scale of the composition? (enter an integer)", k1, k2),
			strconv.Itoa(k1*k2), "Composing similarities multiplies their scale factors.", f)
	}
}

func w19GeomPtolemy(o int) generator.Problem {
	p, q := w19N(1, 12), w19N(1, 12)
	f := map[string]string{}
	w19I(f, "p", p)
	w19I(f, "q", q)
	switch o {
	case 1:
		return w19P(o, fmt.Sprintf("A rectangle has sides %d and %d. What is the product of its diagonals? (enter an integer)", p, q),
			strconv.Itoa(p*p+q*q), "Each diagonal has squared length p^2+q^2, so the product is that value.", f)
	case 2:
		return w19P(o, fmt.Sprintf("A rectangle has sides %d and %d. What is AB*CD + BC*DA in Ptolemy's inequality? (enter an integer)", p, q),
			strconv.Itoa(2*p*q), "The two adjacent products are each pq.", f)
	default:
		return w19P(o, fmt.Sprintf("For the same rectangle with sides %d and %d, what is the difference between the two sides of Ptolemy's relation? (enter an integer)", p, q),
			strconv.Itoa((p-q)*(p-q)), "The equality gap is the squared difference of the sides.", f)
	}
}

func w19GeomTriangle(o int) generator.Problem {
	f := map[string]string{}
	switch o {
	case 1:
		a, b := w19N(1, 20), w19N(1, 20)
		w19I(f, "a", a)
		w19I(f, "b", b)
		return w19P(o, fmt.Sprintf("For z = %d and w = %di, what is |z|+|w|? (enter an integer)", a, b),
			strconv.Itoa(a+b), "The moduli are the two positive coordinates.", f)
	case 2:
		a, b := w19N(1, 10), w19N(1, 10)
		c, d := w19N(1, 10), w19N(1, 10)
		w19I(f, "a", a)
		w19I(f, "b", b)
		w19I(f, "c", c)
		w19I(f, "d", d)
		return w19P(o, fmt.Sprintf("For z = %d+%di and w = %d+%di, what is |z+w|^2? (enter an integer)", a, b, c, d),
			strconv.Itoa((a+c)*(a+c)+(b+d)*(b+d)), "Add the coordinates, then square and add.", f)
	default:
		a, b := w19N(1, 12), w19N(1, 12)
		w19I(f, "a", a)
		w19I(f, "b", b)
		return w19P(o, fmt.Sprintf("For z = %d+%di, what is |z|^2? (enter an integer)", a, b),
			strconv.Itoa(a*a+b*b), "The squared modulus is the sum of the squared parts.", f)
	}
}

func w19GeomCentroid(o int) generator.Problem {
	a, b, c := w19N(1, 9), w19N(1, 9), w19N(1, 9)
	d, e, g := w19N(1, 9), w19N(1, 9), w19N(1, 9)
	f := map[string]string{}
	w19I(f, "a", a)
	w19I(f, "b", b)
	w19I(f, "c", c)
	w19I(f, "d", d)
	w19I(f, "e", e)
	w19I(f, "f", g)
	switch o {
	case 1:
		return w19P(o, fmt.Sprintf("What is three times the real coordinate of the centroid of (%d,%d), (%d,%d), (%d,%d)? (enter an integer)", a, b, c, d, e, g),
			strconv.Itoa(a+c+e), "The centroid averages the coordinates, so three times it is their sum.", f)
	case 2:
		return w19P(o, fmt.Sprintf("What is three times the imaginary coordinate of the centroid of (%d,%d), (%d,%d), (%d,%d)? (enter an integer)", a, b, c, d, e, g),
			strconv.Itoa(b+d+g), "Average the imaginary coordinates and multiply by three.", f)
	default:
		return w19P(o, fmt.Sprintf("What is the sum of the real coordinates of the vertices (%d,%d), (%d,%d), (%d,%d)? (enter an integer)", a, b, c, d, e, g),
			strconv.Itoa(a+c+e), "Sum the real parts before dividing.", f)
	}
}

func w19GeomPerpendicular(o int) generator.Problem {
	f := map[string]string{}
	switch o {
	case 1:
		a, b := w19N(1, 10), w19N(1, 10)
		c, d := w19N(1, 10), w19N(1, 10)
		w19I(f, "a", a)
		w19I(f, "b", b)
		w19I(f, "c", c)
		w19I(f, "d", d)
		return w19P(o, fmt.Sprintf("What is the dot product of (%d,%d) and (%d,%d)? (enter an integer)", a, b, c, d),
			strconv.Itoa(a*c+b*d), "The dot product is ac+bd.", f)
	case 2:
		a, b := w19N(1, 10), w19N(1, 10)
		c, d := w19N(1, 10), w19N(1, 10)
		w19I(f, "a", a)
		w19I(f, "b", b)
		w19I(f, "c", c)
		w19I(f, "d", d)
		return w19P(o, fmt.Sprintf("What is the determinant of (%d,%d) and (%d,%d)? (enter an integer)", a, b, c, d),
			strconv.Itoa(a*d-b*c), "The determinant is ad-bc.", f)
	default:
		p, q := w19N(1, 12), w19N(1, 12)
		w19I(f, "p", p)
		w19I(f, "q", q)
		return w19P(o, fmt.Sprintf("How many right angles does a %d by %d rectangle have? (enter an integer)", p, q),
			"4", "Every rectangle has four right angles.", f)
	}
}

func w19GeomLattice(o int) generator.Problem {
	f := map[string]string{}
	switch o {
	case 1:
		n := w19N(1, 30)
		w19I(f, "n", n)
		return w19P(o, fmt.Sprintf("How many lattice points lie in the closed square [0,%d] x [0,%d]? (enter an integer)", n, n),
			strconv.Itoa((n+1)*(n+1)), "Each axis has n+1 integer coordinates.", f)
	case 2:
		a, b := w19N(0, 12), w19N(0, 12)
		c, d := w19N(0, 12), w19N(0, 12)
		w19I(f, "a", a)
		w19I(f, "b", b)
		w19I(f, "c", c)
		w19I(f, "d", d)
		return w19P(o, fmt.Sprintf("What is the squared distance between (%d,%d) and (%d,%d)? (enter an integer)", a, b, c, d),
			strconv.Itoa((c-a)*(c-a)+(d-b)*(d-b)), "The squared distance is the norm of the difference.", f)
	default:
		w, h := w19N(1, 15), w19N(1, 15)
		w19I(f, "w", w)
		w19I(f, "h", h)
		return w19P(o, fmt.Sprintf("What is the area of the lattice rectangle [0,%d] x [0,%d]? (enter an integer)", w, h),
			strconv.Itoa(w*h), "A lattice rectangle has integer area w*h.", f)
	}
}
