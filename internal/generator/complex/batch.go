package complex

// The 2026-10-17 batch: sixteen concepts that take complex_numbers from 26 to the domain target of 42.
//
// Every answer is computed from the instance. Seven concepts answer with the `complex` grading type,
// whose path rewrites `i` to SymPy's `I` and compares expressions; the rest are exact integers or
// recognitions. Answers avoid polar form except where the angle is a whole number of degrees, because
// the grader's polar conversion only accepts integers.

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

func registerBatch(reg *generator.Registry) {
	reg.Register("complex.basics.parts", &partsGen{})
	reg.Register("complex.basics.powers", &powersGen{})
	reg.Register("complex.basics.reciprocal", &reciprocalGen{})
	reg.Register("complex.basics.geometry", &planeDistanceGen{})
	reg.Register("complex.basics.loci", &lociGen{})
	reg.Register("complex.basics.arg_arithmetic", &argArithmeticGen{})
	reg.Register("complex.basics.equations", &complexEquationsGen{})
	reg.Register("complex.adv.conjugate_roots", &conjugateRootsGen{})
	reg.Register("complex.adv.trig_exponential", &trigExponentialGen{})
	reg.Register("complex.adv.hyperbolic", &hyperbolicGen{})
	reg.Register("complex.adv.geometric_series", &complexSeriesGen{})
	reg.Register("complex.adv.mobius_fixed", &mobiusFixedGen{})
	reg.Register("complex.adv.residue_simple", &simpleResidueGen{})
	reg.Register("complex.adv.cross_ratio", &crossRatioGen{})
	reg.Register("complex.topology.open_sets", &openSetGen{})
	reg.Register("complex.topology.domain", &domainGen{})
}

// ----- helpers -----

func gcdInt(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// fracStr renders n/d in lowest terms with a positive denominator.
func fracStr(n, d int) string {
	if d < 0 {
		n, d = -n, -d
	}
	g := gcdInt(n, d)
	if g > 1 {
		n, d = n/g, d/g
	}
	if d == 1 {
		return strconv.Itoa(n)
	}
	return fmt.Sprintf("%d/%d", n, d)
}

// complexFrac renders (reN/reD) + (imN/imD)i, which is what a conjugate division produces.
func complexFrac(reN, reD, imN, imD int) string {
	re := fracStr(reN, reD)
	im := fracStr(imN, imD)
	im = strings.TrimPrefix(im, "-")
	if imN >= 0 {
		return re + "+" + im + "i"
	}
	return re + "-" + im + "i"
}

func absInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// signed renders a leading sign for use inside a sum.
func signed(n int) string {
	if n >= 0 {
		return "+" + strconv.Itoa(n)
	}
	return strconv.Itoa(n)
}

// ----- 1. parts -----

type partsGen struct{}

func (g *partsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	span := mathutil.ScaleInt(ctx.Difficulty, 8, 30)
	a := rand.Intn(2*span+1) - span
	if a == 0 {
		a = 3
	}
	b := rand.Intn(2*span+1) - span
	if b == 0 {
		b = 5
	}
	left := fmt.Sprintf("x %s %di", map[bool]string{true: "+", false: "-"}[a >= 0], absInt(a))
	q := fmt.Sprintf("The complex numbers \\(%s\\) and \\(%d + yi\\) are equal. What is \\(x + y\\)?", left, b)
	ans := strconv.Itoa(b + a)
	p := generator.Problem{
		Question: q,
		Answer:   ans,
		Facts: map[string]string{
			"a":         strconv.Itoa(a),
			"b":         strconv.Itoa(b),
			"statement": fmt.Sprintf("%s = %d + yi", left, b),
			"x":         strconv.Itoa(b),
			"y":         strconv.Itoa(a),
			"value":     ans,
			"answer":    ans,
		},
	}
	p.Explanation = fmt.Sprintf("Subgoal 1: State governing formula. Two complex numbers are equal if and only if their real parts agree and their imaginary parts agree: Re(z₁) = Re(z₂) and Im(z₁) = Im(z₂).\n\nSubgoal 2: Substitute parameters. From the equality %s = %d + yi, equate the real components to obtain x = %d and the imaginary components to obtain y = %d.\n\nSubgoal 3: Evaluate and simplify. Compute the required sum: x + y = %d + (%d) = %d.",
		left, b, b, a, b, a, b+a)
	return p
}

// ----- 2. powers -----

type powersGen struct{}

func (g *powersGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	if rand.Intn(2) == 0 {
		k := mathutil.ScaleInt(ctx.Difficulty, 12, 320) + rand.Intn(4)
		answers := []string{"1", "i", "-1", "-i"}
		rem := k % 4
		ans := answers[rem]
		p := generator.Problem{
			Question: fmt.Sprintf("What is \\(i^{%d}\\)?", k),
			Answer:   ans,
			Facts: map[string]string{
				"base":     "i",
				"exponent": strconv.Itoa(k),
				"value":    ans,
				"answer":   ans,
				"method":   fmt.Sprintf("The powers of i repeat with period 4; %d leaves remainder %d on division by 4.", k, rem),
			},
		}
		p.Explanation = fmt.Sprintf("Subgoal 1: State governing formula. The powers of the imaginary unit i repeat periodically with period 4: i⁰ = 1, i¹ = i, i² = -1 and i³ = -i, so iᵏ = i^(k mod 4).\n\nSubgoal 2: Substitute parameters. Divide the exponent %d by 4 to determine the remainder: %d = 4(%d) + %d, yielding remainder %d.\n\nSubgoal 3: Evaluate and simplify. Evaluate the periodic power: i^{%d} = i^{%d} = %s.",
			k, k, k/4, rem, rem, k, rem, ans)
		return p
	}
	maxCoeff := mathutil.ScaleInt(ctx.Difficulty, 4, 12)
	a := 1 + rand.Intn(maxCoeff)
	b := 1 + rand.Intn(maxCoeff)
	re := a*a - b*b
	im := 2 * a * b
	ans := fmtComplex(re, im)
	zStr := fmtComplex(a, b)
	p := generator.Problem{
		Question: fmt.Sprintf("What is \\((%s)^2\\)?", zStr),
		Answer:   ans,
		Facts: map[string]string{
			"base":     zStr,
			"exponent": "2",
			"value":    ans,
			"answer":   ans,
			"method":   fmt.Sprintf("Expand the square: (a+bi)^2 = a^2 - b^2 + 2abi, so with a = %d and b = %d the real part is %d - %d and the imaginary part is %d.", a, b, a*a, b*b, im),
		},
	}
	p.Explanation = fmt.Sprintf("Subgoal 1: State governing formula. Expand the square of a complex binomial using the identity (a + bi)² = a² - b² + 2abi, derived from i² = -1.\n\nSubgoal 2: Substitute parameters. Substitute a = %d and b = %d into the expansion: (%d)² - (%d)² + 2(%d)(%d)i.\n\nSubgoal 3: Evaluate and simplify. Compute the real part %d - %d = %d and imaginary part %d to obtain %s.",
		a, b, a, b, a, b, a*a, b*b, re, im, ans)
	return p
}

// ----- 3. reciprocal -----

type reciprocalGen struct{}

func (g *reciprocalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	maxCoeff := mathutil.ScaleInt(ctx.Difficulty, 4, 10)
	a := 1 + rand.Intn(maxCoeff)
	b := 1 + rand.Intn(maxCoeff)
	den := a*a + b*b
	ans := complexFrac(a, den, -b, den)
	zStr := fmtComplex(a, b)
	conjStr := fmtComplex(a, -b)
	p := generator.Problem{
		Question: fmt.Sprintf("Write \\(\\dfrac{1}{%s}\\) in the form \\(a + bi\\).", zStr),
		Answer:   ans,
		Facts: map[string]string{
			"a":         strconv.Itoa(a),
			"b":         strconv.Itoa(b),
			"den":       strconv.Itoa(den),
			"z":         zStr,
			"conjugate": conjStr,
			"real_num":  strconv.Itoa(a),
			"imag_num":  strconv.Itoa(-b),
			"value":     ans,
			"answer":    ans,
		},
	}
	p.Explanation = fmt.Sprintf("Subgoal 1: State governing formula. To write 1/z in Cartesian form a + bi, multiply the numerator and denominator by the complex conjugate z̄ = a - bi: 1/z = z̄ / |z|² = (a - bi) / (a² + b²).\n\nSubgoal 2: Substitute parameters. For z = %s, conjugate is z̄ = %s and the squared modulus is |z|² = %d² + %d² = %d.\n\nSubgoal 3: Evaluate and simplify. Distribute the denominator and reduce real and imaginary fractions to lowest terms to obtain %s.",
		zStr, conjStr, a, b, den, ans)
	return p
}

// ----- 4. distance in the plane -----

type planeDistanceGen struct{}

func (g *planeDistanceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	maxM := mathutil.ScaleInt(ctx.Difficulty, 3, 5)
	m := 2 + rand.Intn(maxM-1)
	n := 1 + rand.Intn(m-1)
	for gcdInt(m, n) > 1 || (m%2 == 1 && n%2 == 1) {
		m = 2 + rand.Intn(maxM-1)
		n = 1 + rand.Intn(m-1)
	}
	da := m*m - n*n
	db := 2 * m * n
	dc := m*m + n*n
	k := 1 + rand.Intn(mathutil.ScaleInt(ctx.Difficulty, 1, 2))
	da *= k
	db *= k
	dc *= k
	if rand.Intn(2) == 0 {
		da, db = db, da
	}
	if rand.Intn(2) == 0 {
		da = -da
	}
	if rand.Intn(2) == 0 {
		db = -db
	}
	coordSpan := mathutil.ScaleInt(ctx.Difficulty, 6, 15)
	r1 := rand.Intn(2*coordSpan+1) - coordSpan
	i1 := rand.Intn(2*coordSpan+1) - coordSpan
	r2, i2 := r1+da, i1+db
	z1Str := fmtComplex(r1, i1)
	z2Str := fmtComplex(r2, i2)
	q := fmt.Sprintf("Find the distance in the complex plane between \\(%s\\) and \\(%s\\).",
		z1Str, z2Str)
	ans := strconv.Itoa(dc)
	p := generator.Problem{
		Question: q,
		Answer:   ans,
		Facts: map[string]string{
			"z1":     z1Str,
			"z2":     z2Str,
			"dx":     strconv.Itoa(da),
			"dy":     strconv.Itoa(db),
			"value":  ans,
			"answer": ans,
		},
	}
	p.Explanation = fmt.Sprintf("Subgoal 1: State governing formula. The distance between two complex points z₁ = x₁ + y₁i and z₂ = x₂ + y₂i is the modulus of their difference: d = |z₂ - z₁| = √((x₂ - x₁)² + (y₂ - y₁)²).\n\nSubgoal 2: Substitute parameters. Substitute the coordinates to compute coordinate differences: Δx = %d - (%d) = %d and Δy = %d - (%d) = %d.\n\nSubgoal 3: Evaluate and simplify. Evaluate the Euclidean norm: √((%d)² + (%d)²) = √(%d) = %d.",
		r2, r1, da, i2, i1, db, da, db, da*da+db*db, dc)
	return p
}

// ----- 5. loci -----

type lociGen struct{}

func (g *lociGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	p := 1 + rand.Intn(4)
	q := 1 + rand.Intn(4)
	if rand.Intn(2) == 0 {
		q = -q
	}
	r := 2 + rand.Intn(4)
	z0 := fmtComplex(p, q)

	if rand.Intn(2) == 0 {
		statement := fmt.Sprintf("|z - (%s)| = %d", z0, r)
		correct := fmt.Sprintf("a circle centred at %s with radius %d", z0, r)
		wrongCentre := fmt.Sprintf("a circle centred at %d with radius %d", r, p+q)
		wrongRadius := fmt.Sprintf("a circle centred at %s with radius %d", z0, r+1)
		wrongType := fmt.Sprintf("the perpendicular bisector of the segment from 0 to %s", z0)
		pr := generator.MultipleChoice(
			fmt.Sprintf("Describe the set of complex numbers z with \\(|z - (%s)| = %d\\).", z0, r),
			correct, wrongCentre, wrongRadius, wrongType)
		pr.Facts = map[string]string{
			"statement": statement, "kind": "circle", "value": correct,
			"reason": fmt.Sprintf("|z - z0| is the distance from z to %s, so the points at the fixed distance %d from it form a circle centred there.", z0, r),
		}
		pr.Explanation = fmt.Sprintf("|z - z0| is the distance from z to z0, so the set of points at a fixed distance %d from %s is a circle centred there.", r, z0)
		return pr
	}

	p2 := p + 1 + rand.Intn(3)
	q2 := q + 1 + rand.Intn(3)
	z1 := fmtComplex(p, q)
	z2 := fmtComplex(p2, q2)
	statement := fmt.Sprintf("|z - (%s)| = |z - (%s)|", z1, z2)
	correct := fmt.Sprintf("the perpendicular bisector of the segment from %s to %s", z1, z2)
	wrongType := fmt.Sprintf("a circle centred at %s with radius %d", z1, p2)
	wrongPoint := fmt.Sprintf("the single point %s", z1)
	wrongHalf := fmt.Sprintf("the perpendicular bisector of the segment from 0 to %s", z2)
	pr := generator.MultipleChoice(
		fmt.Sprintf("Describe the set of complex numbers z with \\(|z - (%s)| = |z - (%s)|\\).", z1, z2),
		correct, wrongType, wrongPoint, wrongHalf)
	pr.Facts = map[string]string{
		"statement": statement, "kind": "bisector", "value": correct,
		"reason": fmt.Sprintf("Both distances are equal, so z is equidistant from %s and %s, and the points equidistant from two points form their perpendicular bisector.", z1, z2),
	}
	pr.Explanation = fmt.Sprintf("Both distances are equal, so z is equidistant from %s and %s. The locus of points equidistant from two points is their perpendicular bisector.", z1, z2)
	return pr
}

// ----- 6. argument arithmetic -----

type argArithmeticGen struct{}

func (g *argArithmeticGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	step := 30
	if ctx.Difficulty > 0.6 {
		step = 5
	} else if ctx.Difficulty > 0.3 {
		step = 15
	}
	numSteps := 360 / step
	a1 := step * rand.Intn(numSteps)
	a2 := step * rand.Intn(numSteps)
	sum := (a1 + a2) % 360
	q := fmt.Sprintf("The complex number z has argument %d\u00b0 and w has argument %d\u00b0. Give the argument of \\(zw\\) in degrees, in the interval \\([0, 360)\\).", a1, a2)
	ans := strconv.Itoa(sum)
	p := generator.Problem{
		Question: q,
		Answer:   ans,
		Facts: map[string]string{
			"a1":     strconv.Itoa(a1),
			"a2":     strconv.Itoa(a2),
			"sum":    strconv.Itoa(a1 + a2),
			"value":  ans,
			"answer": ans,
		},
	}
	p.Explanation = fmt.Sprintf("Subgoal 1: State governing formula. When multiplying two complex numbers z and w, their moduli multiply and their arguments add modulo 360°: arg(zw) ≡ arg(z) + arg(w) (mod 360°).\n\nSubgoal 2: Substitute parameters. Substitute the given arguments arg(z) = %d° and arg(w) = %d° to find their raw sum: %d° + %d° = %d°.\n\nSubgoal 3: Evaluate and simplify. Reduce modulo 360° into the standard interval [0, 360)°: %d mod 360 = %d.",
		a1, a2, a1, a2, a1+a2, a1+a2, sum)
	return p
}

// ----- 7. linear equations -----

type complexEquationsGen struct{}

func (g *complexEquationsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	coeffMax := mathutil.ScaleInt(ctx.Difficulty, 3, 6)
	rootMax := mathutil.ScaleInt(ctx.Difficulty, 3, 7)
	a := 1 + rand.Intn(coeffMax)
	b := 1 + rand.Intn(coeffMax)
	if rand.Intn(2) == 0 {
		b = -b
	}
	// Choose the answer first so the coefficient equation is exact. Resample a
	// zero root: "(2 + 3i)z = 0" has the answer 0 + 0i for every coefficient,
	// which tests nothing and prints an answer a learner reads as a typo.
	pr := rand.Intn(2*rootMax+1) - rootMax
	pi := rand.Intn(2*rootMax+1) - rootMax
	for pr == 0 && pi == 0 {
		pr = rand.Intn(2*rootMax+1) - rootMax
		pi = rand.Intn(2*rootMax+1) - rootMax
	}
	c := a*pr - b*pi
	d := a*pi + b*pr
	ans := fmtComplex(pr, pi)
	lhsStr := fmtComplex(a, b)
	rhsStr := fmtComplex(c, d)
	conjStr := fmtComplex(a, -b)
	den := a*a + b*b
	q := fmt.Sprintf("Solve \\((%s)z = %s\\) for the complex number z.", lhsStr, rhsStr)
	p := generator.Problem{
		Question: q,
		Answer:   ans,
		Facts: map[string]string{
			"a":      strconv.Itoa(a),
			"b":      strconv.Itoa(b),
			"c":      strconv.Itoa(c),
			"d":      strconv.Itoa(d),
			"lhs":    lhsStr,
			"rhs":    rhsStr,
			"re":     strconv.Itoa(pr),
			"im":     strconv.Itoa(pi),
			"value":  ans,
			"answer": ans,
		},
	}
	p.Explanation = fmt.Sprintf("Subgoal 1: State governing formula. To solve the linear complex equation (a + bi)z = c + di, divide both sides by (a + bi) by multiplying numerator and denominator by the complex conjugate (a - bi): z = (c + di)(a - bi) / (a² + b²).\n\nSubgoal 2: Substitute parameters. With a + bi = %s and c + di = %s, conjugate is %s and squared modulus is %d² + (%d)² = %d.\n\nSubgoal 3: Evaluate and simplify. Expand numerator (%s)(%s) = %d + (%d)i and divide by %d to obtain z = %s.",
		lhsStr, rhsStr, conjStr, a, b, den, rhsStr, conjStr, den*pr, den*pi, den, ans)
	return p
}

// ----- 8. conjugate roots -----

type conjugateRootsGen struct{}

func (g *conjugateRootsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	p := 1 + rand.Intn(5)
	q := 1 + rand.Intn(5)
	if rand.Intn(2) == 0 {
		q = -q
	}
	root := fmtComplex(p, q)
	correct := fmtComplex(p, -q)
	wrong1 := fmtComplex(-p, q)
	wrong2 := fmtComplex(-p, -q)
	wrong3 := fmtComplex(q, -p)
	pr := generator.MultipleChoice(
		fmt.Sprintf("A polynomial with real coefficients has the complex root %s. Which number must also be a root?", root),
		correct, wrong1, wrong2, wrong3)
	pr.Facts = map[string]string{
		"root": root, "conjugate": correct, "value": correct,
	}
	pr.Explanation = fmt.Sprintf("Real coefficients make the polynomial's values at conjugate points conjugate, so if %s is a root then %s is too: the pair appears together.", root, correct)
	return pr
}

// ----- 9. trig from exponentials -----

// ----- 10. hyperbolic on imaginary arguments -----

// ----- 11. complex geometric series -----

type complexSeriesGen struct{}

func (g *complexSeriesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// The answer is determined by k alone, so the number of distinct answers is
	// the size of this range. It used to stop at 8, which left the series as one
	// of the few concepts a learner could meet twice.
	maxK := mathutil.ScaleInt(ctx.Difficulty, 6, 14)
	k := 2 + rand.Intn(maxK-1) // ratio i/k, modulus 1/k < 1
	den := k*k + 1
	reN := k * k
	imN := k
	ans := complexFrac(reN, den, imN, den)
	r := fmt.Sprintf("i/%d", k)
	q := fmt.Sprintf("The series \\(1 + r + r^2 + \\cdots\\) has ratio \\(r = %s\\). Find its sum as a complex number.", r)
	p := generator.Problem{
		Question: q,
		Answer:   ans,
		Facts: map[string]string{
			"k":        strconv.Itoa(k),
			"ratio":    r,
			"den":      strconv.Itoa(den),
			"real_num": strconv.Itoa(reN),
			"imag_num": strconv.Itoa(imN),
			"value":    ans,
			"answer":   ans,
		},
	}
	p.Explanation = fmt.Sprintf("Subgoal 1: State governing formula. An infinite geometric series 1 + r + r² + ⋯ with complex common ratio |r| < 1 converges to S = 1 / (1 - r).\n\nSubgoal 2: Substitute parameters. For ratio r = %s, modulus is |r| = 1/%d < 1, so S = 1 / (1 - i/%d) = %d / (%d - i).\n\nSubgoal 3: Evaluate and simplify. Multiply numerator and denominator by conjugate (%d + i): S = %d(%d + i) / (%d² + 1²) = (%d + %di)/%d = %s.",
		r, k, k, k, k, k, k, k, k, reN, imN, den, ans)
	return p
}

// ----- 12. Mobius fixed points -----

type mobiusFixedGen struct{}

func (g *mobiusFixedGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	span := mathutil.ScaleInt(ctx.Difficulty, 3, 7)
	c := 1 + rand.Intn(mathutil.ScaleInt(ctx.Difficulty, 1, 2))
	rootA := rand.Intn(2*span+1) - span
	rootB := rand.Intn(2*span+1) - span
	for rootB == rootA {
		rootB = rand.Intn(2*span+1) - span
	}
	a := rand.Intn(2*span+1) - span
	b := -c * rootA * rootB // fixes the constant term of cz^2 + (d-a)z - b = 0
	d := a - c*(rootA+rootB)
	for a*d-b*c == 0 { // degenerate maps are not Mobius transformations
		a++
		d = a - c*(rootA+rootB)
	}
	lo, hi := rootA, rootB
	if lo > hi {
		lo, hi = hi, lo
	}
	ans := strconv.Itoa(hi)
	question := fmt.Sprintf("The Moebius transformation is \\(f(z) = \\dfrac{%dz %s}{%dz %s}\\). Its fixed points solve \\(z = f(z)\\). Give the fixed point with the larger real part.",
		a, signed(b), c, signed(d))
	pb := generator.Problem{
		Question: question,
		Answer:   ans,
		Facts: map[string]string{
			"a":      strconv.Itoa(a),
			"b":      strconv.Itoa(b),
			"c":      strconv.Itoa(c),
			"d":      strconv.Itoa(d),
			"fz":     fmt.Sprintf("(%dz %s)/(%dz %s)", a, signed(b), c, signed(d)),
			"root1":  strconv.Itoa(lo),
			"root2":  strconv.Itoa(hi),
			"value":  ans,
			"answer": ans,
		},
	}
	pb.Explanation = fmt.Sprintf("Subgoal 1: State governing formula. A fixed point of the Möbius transformation f(z) = (az + b)/(cz + d) satisfies z = f(z), which cross-multiplies to cz² + (d - a)z - b = 0.\n\nSubgoal 2: Substitute parameters. Substitute coefficients a = %d, b = %d, c = %d and d = %d to form %dz² %sz %s = 0, whose roots are z₁ = %d and z₂ = %d.\n\nSubgoal 3: Evaluate and simplify. Compare the real roots: the larger fixed point is %d.",
		a, b, c, d, c, signed(d-a), signed(-b), lo, hi, hi)
	return pb
}

// ----- 13. residue at a simple pole -----

type simpleResidueGen struct{}

func (g *simpleResidueGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	span := mathutil.ScaleInt(ctx.Difficulty, 3, 7)
	a := 1 + rand.Intn(mathutil.ScaleInt(ctx.Difficulty, 3, 6))
	p := rand.Intn(2*span+1) - span
	q := 1 + rand.Intn(span)
	u := rand.Intn(2*span+1) - span
	v := rand.Intn(2*span+1) - span
	z0 := fmtComplex(p, q)
	re := a*p + u
	im := a*q + v
	ans := fmtComplex(re, im)
	q2 := fmt.Sprintf("The function is \\(f(z) = \\dfrac{%dz %s}{z - (%s)}\\). Find the residue of f at \\(z = %s\\).",
		a, fmtSignedComplex(u, v), z0, z0)
	pb := generator.Problem{
		Question: q2,
		Answer:   ans,
		Facts: map[string]string{
			"a":        strconv.Itoa(a),
			"z0":       z0,
			"constant": fmtSignedComplex(u, v),
			"fz":       fmt.Sprintf("(%dz %s)/(z - (%s))", a, fmtSignedComplex(u, v), z0),
			"real":     strconv.Itoa(re),
			"imag":     strconv.Itoa(im),
			"value":    ans,
			"answer":   ans,
		},
	}
	pb.Explanation = fmt.Sprintf("Subgoal 1: State governing formula. At a simple pole z₀, the residue of f(z) = g(z)/(z - z₀) is given by Res(f, z₀) = lim_{z → z₀} (z - z₀)f(z) = g(z₀).\n\nSubgoal 2: Substitute parameters. Here g(z) = %dz %s where the pole is located at z₀ = %s.\n\nSubgoal 3: Evaluate and simplify. Evaluate the numerator at z₀: %d(%s) %s = %s.",
		a, fmtSignedComplex(u, v), z0, a, z0, fmtSignedComplex(u, v), ans)
	return pb
}

// fmtSignedComplex renders a complex constant with an explicit leading sign, for use inside a sum.
func fmtSignedComplex(r, i int) string {
	if i >= 0 {
		return fmt.Sprintf("%s + %di", signed(r), i)
	}
	return fmt.Sprintf("%s - %di", signed(r), -i)
}

// ----- 14. cross ratio -----

type crossRatioGen struct{}

func (g *crossRatioGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	span := mathutil.ScaleInt(ctx.Difficulty, 2, 6)
	z1 := rand.Intn(span) - span/2
	gap1 := 1 + rand.Intn(span)
	gap2 := 1 + rand.Intn(span)
	gap3 := 1 + rand.Intn(span)
	z2 := z1 + gap1
	z3 := z2 + gap2
	z4 := z3 + gap3
	num := (z3 - z1) * (z4 - z2)
	den := (z3 - z2) * (z4 - z1)
	ans := fracStr(num, den)
	q := fmt.Sprintf("For the real points \\(z_1 = %d\\), \\(z_2 = %d\\), \\(z_3 = %d\\) and \\(z_4 = %d\\), compute the cross-ratio \\(\\dfrac{(z_3 - z_1)(z_4 - z_2)}{(z_3 - z_2)(z_4 - z_1)}\\).",
		z1, z2, z3, z4)
	p := generator.Problem{
		Question: q,
		Answer:   ans,
		Facts: map[string]string{
			"z1":          strconv.Itoa(z1),
			"z2":          strconv.Itoa(z2),
			"z3":          strconv.Itoa(z3),
			"z4":          strconv.Itoa(z4),
			"numerator":   strconv.Itoa(num),
			"denominator": strconv.Itoa(den),
			"value":       ans,
			"answer":      ans,
		},
	}
	p.Explanation = fmt.Sprintf("Subgoal 1: State governing formula. The cross-ratio (z₁, z₂; z₃, z₄) is defined as (z₃ - z₁)(z₄ - z₂) / ((z₃ - z₂)(z₄ - z₁)).\n\nSubgoal 2: Substitute parameters. Substitute the given points: numerator is (%d - %d)(%d - %d) = (%d)(%d) = %d while denominator is (%d - %d)(%d - %d) = (%d)(%d) = %d.\n\nSubgoal 3: Evaluate and simplify. Divide terms and reduce fraction to lowest terms: %d/%d = %s.",
		z3, z1, z4, z2, z3-z1, z4-z2, num, z3, z2, z4, z1, z3-z2, z4-z1, den, num, den, ans)
	return p
}

// ----- 15. open and closed sets -----

type openSetGen struct{}

func (g *openSetGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct{ set, class, why string }
	entries := []entry{
		{"{z : |z| < 3}", "open, not closed", "the boundary |z| = 3 is excluded, and every point inside has a small disk that stays inside"},
		{"{z : |z| \u2264 3}", "closed, not open", "the boundary circle is included as well as the interior, and an interior point can never move in every direction and stay inside"},
		{"{z : Re z > 0}", "open, not closed", "the boundary line Re z = 0 is excluded, and each point has room to its left and right"},
		{"{z : Re z \u2265 0}", "closed, not open", "the boundary line is included, and a point on it cannot move left at all"},
	}
	e := entries[rand.Intn(len(entries))]
	pr := generator.MultipleChoice(
		fmt.Sprintf("Classify the set \\(%s\\) in the complex plane.", e.set),
		e.class, "open, not closed", "closed, not open", "both open and closed", "neither open nor closed")
	pr.Facts = map[string]string{"set": e.set, "class": e.class, "why": e.why, "value": e.class}
	pr.Explanation = fmt.Sprintf("The set is %s: %s. Openness is about whether the boundary is included, which is a separate question from whether the set is all in one piece.", e.class, e.why)
	return pr
}

// ----- 16. domains -----

type domainGen struct{}

func (g *domainGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct{ set, verdict, why string }
	domains := []entry{
		{"the punctured plane C \\ {0}", "a domain",
			"it is open, because every point has a small disk avoiding 0, and connected, because any two points can be joined by a path that bends around 0"},
		{"the open unit disk", "a domain",
			"it is open, because its boundary |z| = 1 is excluded, and connected, because any two interior points can be joined by a straight segment inside it"},
		{"the upper half-plane", "a domain",
			"it is open, because its boundary, the real axis Im z = 0, is excluded, and connected, because a straight segment between two points with positive imaginary part stays there"},
		{"the annulus {z : 1 < |z| < 2}", "a domain",
			"it is open, because both boundary circles are excluded, and connected, because a path can go around the hole between radius 1 and radius 2"},
	}
	notDomains := []entry{
		{"the union of two disjoint open disks", "not a domain",
			"it is open, but it is not connected: it splits into two disjoint non-empty open pieces"},
		{"the closed unit disk", "not a domain",
			"it is connected, but it is not open: a point on |z| = 1 has no disk that stays inside"},
		{"the real axis", "not a domain",
			"it is connected, but it is not open in the complex plane: a point on the line has no disk that stays on the line"},
		{"{z : |z| \u2265 1}", "not a domain",
			"it is connected, but it is not open: the circle |z| = 1 belongs to it and has no disk that stays inside"},
	}
	answer := domains[rand.Intn(len(domains))]
	rand.Shuffle(len(notDomains), func(i, j int) { notDomains[i], notDomains[j] = notDomains[j], notDomains[i] })
	pr := generator.MultipleChoice(
		"Which of these sets is a domain, meaning an open connected subset of the complex plane?",
		answer.set, notDomains[0].set, notDomains[1].set, notDomains[2].set)
	pr.Facts = map[string]string{
		"value": answer.set, "verdict": answer.verdict, "why": answer.why,
		"distractor1": notDomains[0].set, "distractor2": notDomains[1].set, "distractor3": notDomains[2].set,
	}
	pr.Explanation = fmt.Sprintf("%s is %s: %s. The two requirements are independent, so a set that fails either one is not a domain.", answer.set, answer.verdict, answer.why)
	return pr
}
