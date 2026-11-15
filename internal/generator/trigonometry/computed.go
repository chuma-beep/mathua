package trigonometry

// The computed trigonometry concepts the domain was missing. The first wave was a coterminal angle,
// the quadrant of an angle, arc length and sector area as multiples of pi, the sine and cosine of a
// doubled angle, the cofunction of an angle, a triangle's area from its three sides and the
// amplitude of a sine wave. The second wave, which takes the domain to the forty-eight its target
// names, is the sum and difference formulas, the half-angle formulas, a quadratic in a ratio, the
// ambiguous case of the law of sines, an area from two sides and an included angle, a vector's
// length and its dot product, and a sinusoid's phase shift.
//
// The domain was thirty-two concepts and, before this, its `numeric` surface was built from
// Pythagorean triples read off a table: given a ratio, name the other one. These compute an answer
// from parameters the learner has not seen, which is the shift section 4.2 of
// docs/content-expansion.md measures.
//
// Every answer is an integer, an exact reduced fraction or a count by construction: the coterminal
// angle is the remainder the turn count leaves, the quadrant is the count of quarter-turns, the
// radius is a multiple of the angle's denominator so the pi cancels, the doubled angle's fraction
// reduces by a gcd, the cofunction complements to 90 degrees, Heron's formula is used only on
// triples whose area is whole, and the amplitude is the absolute value of a coefficient. The second
// wave reaches the same guarantee by other routes: a product of denominators that cancels, a
// half-angle radicand that is the square of the numerator over the square of the hypotenuse, a
// scaled triple for a vector's length, and factors whose solutions in the first turn are counted.
import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerComputed(reg *generator.Registry) {
	reg.Register("trig.basics.coterminal", &coterminalGen{})
	reg.Register("trig.basics.quadrant", &quadrantGen{})
	reg.Register("trig.basics.arc_length", &arcLengthGen{})
	reg.Register("trig.basics.sector_area", &sectorAreaGen{})
	reg.Register("trig.ident.double_angle", &doubleAngleGen{})
	reg.Register("trig.ident.cofunction", &cofunctionGen{})
	reg.Register("trig.adv.heron_area", &heronAreaGen{})
	reg.Register("trig.graph.amplitude", &amplitudeGen{})
	reg.Register("trig.ident.sum_diff", &sumDiffGen{})
	reg.Register("trig.ident.half_angle", &halfAngleGen{})
	reg.Register("trig.eq.quadratic", &quadraticEqGen{})
	reg.Register("trig.adv.ambiguous_case", &ambiguousCaseGen{})
	reg.Register("trig.adv.area_sas", &areaSASGen{})
	reg.Register("trig.adv.vector_magnitude", &vectorMagnitudeGen{})
	reg.Register("trig.adv.vector_dot", &vectorDotGen{})
	reg.Register("trig.graph.phase_shift", &phaseShiftGen{})
}

func trigSmallInt(lo, hi int) int {
	if hi < lo {
		hi = lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

// piFrac renders k*pi/m in the compact form the corpus uses, with the coefficient only when it is
// not one.
func piFrac(k, m int) string {
	if m == 1 {
		if k == 1 {
			return "\\pi"
		}
		return fmt.Sprintf("%d\\pi", k)
	}
	if k == 1 {
		return fmt.Sprintf("\\pi/%d", m)
	}
	return fmt.Sprintf("%d\\pi/%d", k, m)
}

// coterminalGen reduces an angle the learner is handed outside the first turn.
type coterminalGen struct{}

func (g *coterminalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	turns := 1 + rand.Intn(3+int(ctx.Difficulty*3))
	base := 1 + rand.Intn(359)
	angle := base + 360*turns
	dir := "positive"
	if rand.Intn(2) == 0 {
		angle = base - 360*turns
		dir = "negative"
	}
	return generator.Problem{
		Question: fmt.Sprintf("Find the angle in \\([0^{\\circ}, 360^{\\circ})\\) that is coterminal with \\(%d^{\\circ}\\).", angle),
		Answer:   strconv.Itoa(base),
		Explanation: fmt.Sprintf("A %s angle of %d degrees and %d degrees differ by %d full turns, and coterminal angles differ by whole turns. Adding or subtracting \\(360^{\\circ}\\) repeatedly leaves %d\\(^{\\circ}\\) in the first turn, so that is the coterminal angle.",
			dir, angle, base, turns, base),
		Facts: map[string]string{
			"value": strconv.Itoa(base),
			"given": strconv.Itoa(angle),
			"turns": strconv.Itoa(turns),
		},
	}
}

// quadrantGen counts the quarter-turns an angle crosses.
type quadrantGen struct{}

func (g *quadrantGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	angle := 1 + rand.Intn(359)
	for angle%90 == 0 {
		angle = 1 + rand.Intn(359)
	}
	q := angle/90 + 1
	return generator.Problem{
		Question: fmt.Sprintf("In which quadrant does an angle of \\(%d^{\\circ}\\) terminate? Answer 1, 2, 3 or 4.", angle),
		Answer:   strconv.Itoa(q),
		Explanation: fmt.Sprintf("The quadrant boundaries are the multiples of 90 degrees, and %d degrees lies between %d and %d, which is quadrant %d. The angle is measured counterclockwise from the positive x-axis, so the first quadrant spans 0 to 90 degrees and each next quadrant adds another 90.",
			angle, (q-1)*90, q*90, q),
		Facts: map[string]string{"value": strconv.Itoa(q), "angle": strconv.Itoa(angle)},
	}
}

// arcLengthGen uses s = r * theta with theta a fraction of pi a radius is chosen to cancel.
type arcLengthGen struct{}

func (g *arcLengthGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	dens := []int{2, 3, 4, 6}
	m := dens[rand.Intn(len(dens))]
	k := 1 + rand.Intn(m-1)
	for gcd(k, m) != 1 {
		k = 1 + rand.Intn(m-1)
	}
	t := 1 + rand.Intn(4+int(ctx.Difficulty*4))
	r := m * t // chosen so r*k/m is whole and the pi stays symbolic
	n := k * t
	frac := piFrac(k, m)
	return generator.Problem{
		Question: fmt.Sprintf("A circle has radius %d. A central angle of \\(%s\\) subtends an arc. What is the arc length, written as a multiple of \\(\\pi\\)?", r, frac),
		Answer:   strconv.Itoa(n),
		Explanation: fmt.Sprintf("Arc length is \\(s = r\\theta\\), and in radians that is the radius times the angle. Here \\(s = %d \\times %s = %d\\pi\\), so the arc length is \\(%d\\pi\\) and the multiplier is %d. Radians are the unit that makes this formula free of a conversion factor, which is why arc length is stated in radians and not degrees.",
			r, frac, n, n, n),
		Facts: map[string]string{"value": strconv.Itoa(n), "r": strconv.Itoa(r), "num": strconv.Itoa(k), "den": strconv.Itoa(m), "frac": frac},
	}
}

// sectorAreaGen uses A = r^2 theta / 2 with the same radius trick.
type sectorAreaGen struct{}

func (g *sectorAreaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	dens := []int{2, 4, 6}
	m := dens[rand.Intn(len(dens))]
	k := 1 + rand.Intn(m-1)
	for gcd(k, m) != 1 {
		k = 1 + rand.Intn(m-1)
	}
	t := 1 + rand.Intn(3+int(ctx.Difficulty*3))
	r := m * t
	n := (m / 2) * t * t * k // A = (1/2) r^2 (k pi / m) = m t^2 k pi / 2
	frac := piFrac(k, m)
	return generator.Problem{
		Question: fmt.Sprintf("A circle has radius %d. A sector has central angle \\(%s\\). What is the sector's area, written as a multiple of \\(\\pi\\)?", r, frac),
		Answer:   strconv.Itoa(n),
		Explanation: fmt.Sprintf("Sector area is half the radius squared times the angle in radians: \\(A = \\tfrac{1}{2} r^2 \\theta\\). Here \\(A = \\tfrac{1}{2}(%d)^2 \\cdot %s = %d\\pi\\), so the multiplier is %d. The factor of one half is the same one that makes a triangle half its bounding rectangle, and the whole circle is the special case where the angle is \\(2\\pi\\) and the area is \\(\\pi r^2\\).",
			r, frac, n, n),
		Facts: map[string]string{"value": strconv.Itoa(n), "r": strconv.Itoa(r), "num": strconv.Itoa(k), "den": strconv.Itoa(m), "frac": frac},
	}
}

// doubleAngleGen applies sin(2t) = 2 sin t cos t and cos(2t) = cos^2 t - sin^2 t to a triple.
type doubleAngleGen struct{}

func (g *doubleAngleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type triple struct{ opp, adj, hyp int }
	triples := []triple{{3, 4, 5}, {5, 12, 13}, {8, 15, 17}, {20, 21, 29}, {7, 24, 25}, {9, 40, 41}}
	tr := triples[rand.Intn(len(triples))]
	if rand.Intn(2) == 0 { // sin(2t) = 2 sin t cos t
		num := 2 * tr.opp * tr.adj
		den := tr.hyp * tr.hyp
		d := gcd(num, den)
		return generator.Problem{
			Question: fmt.Sprintf("If \\(\\sin\\theta = %d/%d\\) and \\(\\theta\\) is acute, what is \\(\\sin(2\\theta)\\), as a fraction?", tr.opp, tr.hyp),
			Answer:   fmt.Sprintf("%d/%d", num/d, den/d),
			Explanation: fmt.Sprintf("The double-angle formula is \\(\\sin(2\\theta) = 2\\sin\\theta\\cos\\theta\\). The opposite side is %d and the hypotenuse is %d, so the adjacent side is %d and \\(\\cos\\theta = %d/%d\\). Then \\(\\sin(2\\theta) = 2 \\cdot %d/%d \\cdot %d/%d = %d/%d\\). Because \\(\\theta\\) is acute both ratios are positive, which fixes the sign of the product.",
				tr.opp, tr.hyp, tr.adj, tr.adj, tr.hyp, tr.opp, tr.hyp, tr.adj, tr.hyp, num/d, den/d),
			Facts: map[string]string{"num": strconv.Itoa(num / d), "den": strconv.Itoa(den / d), "sin_num": strconv.Itoa(tr.opp), "sin_den": strconv.Itoa(tr.hyp)},
		}
	}
	// cos(2t) = cos^2 t - sin^2 t
	cn := tr.adj*tr.adj - tr.opp*tr.opp
	sign := ""
	if cn < 0 {
		sign = "-"
		cn = -cn
	}
	den := tr.hyp * tr.hyp
	d := gcd(cn, den)
	cn, den = cn/d, den/d
	return generator.Problem{
		Question: fmt.Sprintf("If \\(\\sin\\theta = %d/%d\\) and \\(\\theta\\) is acute, what is \\(\\cos(2\\theta)\\), as a fraction?", tr.opp, tr.hyp),
		Answer:   fmt.Sprintf("%s%d/%d", sign, cn, den),
		Explanation: fmt.Sprintf("The double-angle formula is \\(\\cos(2\\theta) = \\cos^2\\theta - \\sin^2\\theta\\). With opposite %d and hypotenuse %d the adjacent side is %d, so \\(\\cos\\theta = %d/%d\\). Then \\(\\cos(2\\theta) = %d/%d - %d/%d = %s%d/%d\\). The sign follows from the formula rather than from the acute angle: doubling can carry \\(2\\theta\\) past 90 degrees, where the cosine turns negative.",
			tr.opp, tr.hyp, tr.adj, tr.adj, tr.hyp, tr.adj*tr.adj, tr.hyp*tr.hyp, tr.opp*tr.opp, tr.hyp*tr.hyp, sign, cn, den),
		Facts: map[string]string{"num": sign + strconv.Itoa(cn), "den": strconv.Itoa(den), "sin_num": strconv.Itoa(tr.opp), "sin_den": strconv.Itoa(tr.hyp)},
	}
}

// cofunctionGen uses the fact that the sine of an angle is the cosine of its complement.
type cofunctionGen struct{}

func (g *cofunctionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	theta := 1 + rand.Intn(88) // 1..88 keeps the complement strictly inside (0, 90)
	phi := 90 - theta
	return generator.Problem{
		Question: fmt.Sprintf("Given that \\(\\sin(%d^{\\circ}) = \\cos(\\varphi)\\) with \\(0^{\\circ} < \\varphi < 90^{\\circ}\\), what is \\(\\varphi\\)?", theta),
		Answer:   strconv.Itoa(phi),
		Explanation: fmt.Sprintf("The cofunction identity says a trigonometric function of an angle equals the cofunction of its complement: \\(\\sin\\theta = \\cos(90^{\\circ} - \\theta)\\). Here \\(90 - %d = %d\\), so \\(\\varphi = %d^{\\circ}\\). That complement is why cosine is the co-sine of sine, and it is a statement about the two acute angles of a right triangle summing to 90 degrees.",
			theta, phi, phi),
		Facts: map[string]string{"value": strconv.Itoa(phi), "theta": strconv.Itoa(theta)},
	}
}

// pythagoreanPair returns a primitive Pythagorean triple by Euclid's formula.
// The legs are m^2-n^2 and 2mn and the hypotenuse is m^2+n^2, with m and n
// coprime and of opposite parity so that the triple is primitive. The bounded
// retry is a guard rather than the normal path: n = m-1 always satisfies both
// conditions.
func pythagoreanPair(maxM int) (int, int, int) {
	if maxM < 2 {
		maxM = 2
	}
	for attempt := 0; attempt < 40; attempt++ {
		m := 2 + rand.Intn(maxM-1)
		n := 1 + rand.Intn(m-1)
		if gcd(m, n) != 1 || (m-n)%2 == 0 {
			continue
		}
		return m*m - n*n, 2 * m * n, m*m + n*n
	}
	return 3, 4, 5
}

// heronAreaGen uses Heronian triangles: integer sides with an integer area. The
// three knowledge points build them from three different constructions rather
// than from a fixed table, because a table of nine was the whole supply a
// learner could meet.
type heronAreaGen struct{}

func (g *heronAreaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := trigOrdinal(ctx)
	var a, b, c, area int
	switch o {
	case 1:
		// A primitive right triangle. One leg is even, so half their product,
		// which is the area, is already an integer.
		a, b, c = pythagoreanPair(9)
		area = a * b / 2
	case 2:
		// The same construction scaled by k. Scaling multiplies the area by
		// k squared, which keeps it whole.
		x, y, z := pythagoreanPair(7)
		k := 1 + rand.Intn(5)
		a, b, c = k*x, k*y, k*z
		area = k * k * x * y / 2
	default:
		// Gluing a right triangle to its mirror image across a leg gives an
		// isosceles triangle with sides c, c and 2a, height b and area ab.
		x, y, z := pythagoreanPair(9)
		if rand.Intn(2) == 0 {
			a, b, c = z, z, 2*x
			area = x * y
		} else {
			a, b, c = z, z, 2*y
			area = x * y
		}
	}
	semi := (a + b + c) / 2
	return generator.Problem{
		Question: fmt.Sprintf("A triangle has side lengths %d, %d and %d. What is its area?", a, b, c),
		Answer:   strconv.Itoa(area),
		Explanation: fmt.Sprintf("Heron's formula is \\(A = \\sqrt{s(s-a)(s-b)(s-c)}\\) with \\(s\\) the semiperimeter. Here the semiperimeter is %d, so \\(A = \\sqrt{%d \\cdot %d \\cdot %d \\cdot %d} = %d\\). Three sides determine the triangle, so an area is recoverable from them even though no height was given; the law of cosines is what lets Heron's formula be derived from the sides.",
			semi, semi, semi-a, semi-b, semi-c, area),
		Facts:     map[string]string{"value": strconv.Itoa(area), "a": strconv.Itoa(a), "b": strconv.Itoa(b), "c": strconv.Itoa(c)},
		KPOrdinal: o,
	}
}

// amplitudeGen reads the amplitude off the coefficient, sign and all.
type amplitudeGen struct{}

func (g *amplitudeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := trigSmallInt(2, 8+int(ctx.Difficulty*4))
	sign := ""
	if rand.Intn(2) == 0 {
		a = -a
		sign = "-"
	}
	fn := "sin"
	if rand.Intn(2) == 0 {
		fn = "cos"
	}
	amp := a
	if amp < 0 {
		amp = -amp
	}
	return generator.Problem{
		Question: fmt.Sprintf("What is the amplitude of \\(y = %s%d\\,%s(x)\\)?", sign, amp, fn),
		Answer:   strconv.Itoa(amp),
		Explanation: fmt.Sprintf("The amplitude is the distance from the midline to a peak, which is the absolute value of the coefficient: \\(|%d| = %d\\). The sign is not part of the amplitude because a reflection over the axis changes where the wave starts, not how far it swings. The coefficient never changes the period either, which is set by the multiplier on \\(x\\) alone.",
			a, amp),
		Facts: map[string]string{"value": strconv.Itoa(amp), "a": strconv.Itoa(a), "fn": fn, "fnmath": "\\" + fn},
	}
}

// -------------------------------------------------------------------------------------------------
// The second wave: eight more computed concepts, each needing a formula rather than a look-up, two
// angles combined, an angle halved, a quadratic in a ratio, the triangles the law of sines leaves
// open, an area from two sides and an included angle, a vector's length and dot product, and a
// sinusoid's horizontal shift.

// computedTriples are the primitive Pythagorean triples the second wave draws from. A triple is what
// gives an acute angle an exact rational sine and cosine at once, and every one of these answers
// needs both: a sum of two ratios, a half-angle radicand, or a vector's length.
var computedTriples = [][3]int{{3, 4, 5}, {5, 12, 13}, {8, 15, 17}, {7, 24, 25}, {20, 21, 29}, {9, 40, 41}}

// sumDiffGen applies the sum and difference formulas to two triples. The two stated ratios are a sine
// and a cosine of an acute angle, so the formulas' other two terms come from the same triples and
// every answer is a fraction over the product of the two hypotenuses that reduces by a gcd.
type sumDiffGen struct{}

func (g *sumDiffGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// Bounded rather than open ended: see generator.MaxDrawAttempts. A draw that
	// still prints its answer is retried, and after the ceiling the problem is
	// returned as generated so that a leak stays a leak instead of a hang.
	for attempt := 0; ; attempt++ {
		t1 := computedTriples[rand.Intn(len(computedTriples))]
		t2 := computedTriples[rand.Intn(len(computedTriples))]
		sinA, cosA := t1[0], t1[1]
		sinB, cosB := t2[0], t2[1]
		fn, expr := "\\sin", "\\alpha+\\beta"
		term1, term2 := sinA*cosB, cosA*sinB
		if rand.Intn(2) == 0 {
			fn, expr = "\\cos", "\\alpha-\\beta"
			term1, term2 = cosA*cosB, sinA*sinB
		}
		sum, rawDen := term1+term2, t1[2]*t2[2]
		d := gcd(sum, rawDen)
		num, den := sum/d, rawDen/d
		ans := fmt.Sprintf("%d/%d", num, den)
		q := fmt.Sprintf("If \\(\\sin\\alpha = %d/%d\\) and \\(\\cos\\beta = %d/%d\\), with both angles acute, what is \\(%s(%s)\\), as a fraction?", sinA, t1[2], cosB, t2[2], fn, expr)
		if strings.Contains(q, ans) && attempt < generator.MaxDrawAttempts {
			continue
		}
		return generator.Problem{
			Question: q,
			Answer:   ans,
			Explanation: fmt.Sprintf("The formulas for a sum and a difference expand a combined angle into products of the individual ratios. The first triple gives \\(\\sin\\alpha = %d/%d\\) and \\(\\cos\\alpha = %d/%d\\); the second gives \\(\\cos\\beta = %d/%d\\) and \\(\\sin\\beta = %d/%d\\). Then \\(%s(%s) = %d/%d + %d/%d = %d/%d\\), which is \\(%d/%d\\) in lowest terms. Both angles being acute is what keeps every term positive and the signs untroubled.",
				sinA, t1[2], cosA, t1[2], cosB, t2[2], sinB, t2[2], fn, expr, term1, rawDen, term2, rawDen, sum, rawDen, num, den),
			Facts: map[string]string{
				"num": strconv.Itoa(num), "den": strconv.Itoa(den),
				"sum": strconv.Itoa(sum), "raw_den": strconv.Itoa(rawDen),
				"term1": strconv.Itoa(term1), "term2": strconv.Itoa(term2),
				"fn": fn, "expr": expr,
				"sin_a_num": strconv.Itoa(sinA), "sin_a_den": strconv.Itoa(t1[2]),
				"cos_b_num": strconv.Itoa(cosB), "cos_b_den": strconv.Itoa(t2[2]),
			},
		}
	}
}

// halfAngleGen applies the half-angle formulas to a cosine a triple determines. Writing the cosine
// as (q^2 - 2p^2)/q^2 makes (1 - cos)/2 the square of p/q and (1 + cos)/2 the square of the triple's
// other leg over q, so the root is exact and no approximation enters.
type halfAngleGen struct{}

func (g *halfAngleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	t := computedTriples[rand.Intn(len(computedTriples))]
	p, b, q := t[0], t[1], t[2]
	cosNum, cosDen := q*q-2*p*p, q*q
	op, radNum, num := "-", p*p, p
	half := "\\sin(\\theta/2)"
	if rand.Intn(2) == 0 {
		op, radNum, num = "+", b*b, b
		half = "\\cos(\\theta/2)"
	}
	return generator.Problem{
		Question: fmt.Sprintf("If \\(\\cos\\theta = %d/%d\\) with \\(0^{\\circ} < \\theta < 90^{\\circ}\\), what is \\(%s\\), as a fraction?", cosNum, cosDen, half),
		Answer:   fmt.Sprintf("%d/%d", num, q),
		Explanation: fmt.Sprintf("The half-angle formulas give the sine of half an angle as the square root of \\((1-\\cos\\theta)/2\\) and the cosine as the square root of \\((1+\\cos\\theta)/2\\). Halving an angle between 0 and 90 degrees leaves the half between 0 and 45, still in the first quadrant, so the positive root is the one to take. Here \\((1 %s %d/%d)/2 = %d/%d\\), and the square root of that is \\(%d/%d\\). A proportion is a triangle, which is why the answer is a ratio the triple already contained.",
			op, cosNum, cosDen, radNum, cosDen, num, q),
		Facts: map[string]string{
			"num": strconv.Itoa(num), "den": strconv.Itoa(q),
			"cos_num": strconv.Itoa(cosNum), "cos_den": strconv.Itoa(cosDen),
			"op": op, "half": half,
			"rad_num": strconv.Itoa(radNum), "rad_den": strconv.Itoa(cosDen),
		},
	}
}

// trigQuadratic renders a quadratic in sin x the way the corpus writes one: the square first, then
// the linear term and the constant only when they are there, and a coefficient of one unwritten. A
// builder rather than repeated concatenation, because the term loop is what make lint-go rejects.
func trigQuadratic(a, b, c int) string {
	var s strings.Builder
	if a == 1 {
		s.WriteString("\\sin^{2}x")
	} else {
		fmt.Fprintf(&s, "%d\\sin^{2}x", a)
	}
	for _, t := range []struct {
		coef int
		body string
	}{{b, "\\sin x"}, {c, ""}} {
		if t.coef == 0 {
			continue
		}
		sign := " + "
		n := t.coef
		if n < 0 {
			sign, n = " - ", -n
		}
		s.WriteString(sign)
		if n == 1 && t.body != "" {
			s.WriteString(t.body)
			continue
		}
		fmt.Fprintf(&s, "%d%s", n, t.body)
	}
	return s.String()
}

// trigSineFactor renders one factor of a factored quadratic: the sine term, then the constant with
// the sign that makes the root come out.
func trigSineFactor(n, d int) string {
	var s strings.Builder
	if d == 1 {
		s.WriteString("\\sin x")
	} else {
		fmt.Fprintf(&s, "%d\\sin x", d)
	}
	switch {
	case n == 0:
	case n > 0:
		fmt.Fprintf(&s, " - %d", n)
	default:
		fmt.Fprintf(&s, " + %d", -n)
	}
	return "(" + s.String() + ")"
}

// trigSineRootString writes a root as the corpus writes a ratio.
func trigSineRootString(n, d int) string {
	if d == 1 {
		return strconv.Itoa(n)
	}
	return fmt.Sprintf("%d/%d", n, d)
}

// quadraticEqGen asks how many angles solve a quadratic in sin x, which is the shape the factoring
// decides: a root whose sine is a turning point is reached once in the first turn and every other
// root twice.
type quadraticEqGen struct{}

func (g *quadraticEqGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type root struct {
		n, d  int
		angle string
		count int
	}
	roots := []root{
		{-1, 1, "270^{\\circ}", 1},
		{-1, 2, "210^{\\circ} and 330^{\\circ}", 2},
		{0, 1, "0^{\\circ} and 180^{\\circ}", 2},
		{1, 2, "30^{\\circ} and 150^{\\circ}", 2},
		{1, 1, "90^{\\circ}", 1},
	}
	r1 := roots[rand.Intn(len(roots))]
	r2 := roots[rand.Intn(len(roots))]
	for r2.n == r1.n && r2.d == r1.d {
		r2 = roots[rand.Intn(len(roots))]
	}
	a := r1.d * r2.d
	b := -(r1.n*r2.d + r2.n*r1.d)
	c := r1.n * r2.n
	if d := gcd(gcd(a, b), c); d != 1 {
		a, b, c = a/d, b/d, c/d
	}
	eq := trigQuadratic(a, b, c)
	factored := trigSineFactor(r1.n, r1.d) + trigSineFactor(r2.n, r2.d)
	count := r1.count + r2.count
	noun := "solutions"
	if count == 1 {
		noun = "solution"
	}
	return generator.Problem{
		Question: fmt.Sprintf("How many solutions does \\(%s = 0\\) have for \\(0^{\\circ} \\le x < 360^{\\circ}\\)?", eq),
		Answer:   strconv.Itoa(count),
		Explanation: fmt.Sprintf("A quadratic in a ratio factors like any other, so \\(%s = 0\\) is \\(%s = 0\\). That needs \\(\\sin x = %s\\) or \\(\\sin x = %s\\). In the first turn those are reached at %s, and at %s, which is %d %s in all. A root the sine reaches at a turning point gives one angle and every other root gives two, because a horizontal line below the peak crosses the wave twice.",
			eq, factored, trigSineRootString(r1.n, r1.d), trigSineRootString(r2.n, r2.d), r1.angle, r2.angle, count, noun),
		Facts: map[string]string{
			"value": strconv.Itoa(count), "noun": noun, "eq": eq, "factored": factored,
			"a": strconv.Itoa(a), "b": strconv.Itoa(b), "c": strconv.Itoa(c),
			"r1": trigSineRootString(r1.n, r1.d), "r2": trigSineRootString(r2.n, r2.d),
			"s1": r1.angle, "s2": r2.angle,
		},
	}
}

// ambiguousCaseGen counts the triangles a side-side-angle triple leaves open. The angle is fixed at
// 30 degrees so the height is exactly half of the second side, and the side opposite the angle is
// then drawn below, at, between or above the height and that second side.
type ambiguousCaseGen struct{}

func (g *ambiguousCaseGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	b := 2 * (3 + rand.Intn(4+int(ctx.Difficulty*6))) // 6..20, so b sin 30 is a whole number
	h := b / 2
	var a int
	switch rand.Intn(4) {
	case 0:
		a = 1 + rand.Intn(h-1) // shorter than the height
	case 1:
		a = h // exactly the height
	case 2:
		a = h + 1 + rand.Intn(b-h-1) // between the height and b
	default:
		a = b + rand.Intn(5) // at least b
	}
	count := 0
	switch {
	case a < h:
		count = 0
	case a == h:
		count = 1
	case a < b:
		count = 2
	default:
		count = 1
	}
	return generator.Problem{
		Question: fmt.Sprintf("In a triangle \\(A = 30^{\\circ}\\), the side \\(a = %d\\) opposite \\(A\\) and the side \\(b = %d\\) opposite \\(B\\). How many triangles have these measurements? Answer 0, 1 or 2.", a, b),
		Answer:   strconv.Itoa(count),
		Explanation: fmt.Sprintf("The law of sines leaves a side-side-angle triple open because two different angles can share a sine. The pair \\(b = %d\\) and \\(A = 30^{\\circ}\\) fixes a height \\(h = b\\sin A = %d\\) above the far ray, and the side \\(a\\) meets that ray according to how it compares with \\(h\\) and with \\(b\\): below the height it never reaches, between them it meets twice, and at or above \\(b\\) it meets once. Here \\(a = %d\\), so there %s.",
			b, h, a, map[int]string{0: "is no triangle", 1: "is exactly one triangle", 2: "are two triangles"}[count]),
		Facts: map[string]string{
			"value": strconv.Itoa(count), "a": strconv.Itoa(a), "b": strconv.Itoa(b),
			"h": strconv.Itoa(h), "angle": "30",
		},
	}
}

// areaSASGen uses A = ab sin C / 2 with the sides chosen to clear the sine, which is 1/2 at the two
// angles where the sine of a whole number of degrees is rational and beyond that only at 90.
type areaSASGen struct{}

func (g *areaSASGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	angle, sinNum, sinDen := 30+120*rand.Intn(2), 1, 2
	if rand.Intn(3) == 0 {
		angle, sinNum, sinDen = 90, 1, 1
	}
	a := 2 * (2 + rand.Intn(3+int(ctx.Difficulty*4)))
	if sinDen == 2 {
		a *= 2 // so that a*b/4 stays whole
	}
	b := 3 + rand.Intn(3+int(ctx.Difficulty*5))
	area := a * b * sinNum / (2 * sinDen)
	sinVal := fmt.Sprintf("%d/%d", sinNum, sinDen)
	if sinDen == 1 {
		sinVal = "1"
	}
	return generator.Problem{
		Question: fmt.Sprintf("A triangle has two sides of length %d and %d with an included angle of \\(%d^{\\circ}\\). What is its area?", a, b, angle),
		Answer:   strconv.Itoa(area),
		Explanation: fmt.Sprintf("Two sides and the angle between them are enough for the area because that angle pins the height: \\(A = \\tfrac{1}{2}ab\\sin C\\). Here \\(A = \\tfrac{1}{2} \\cdot %d \\cdot %d \\cdot %s = %d\\). Both sides are used, and only one angle, which is why this formula does not care about the third side or the other two angles.",
			a, b, sinVal, area),
		Facts: map[string]string{
			"value": strconv.Itoa(area), "a": strconv.Itoa(a), "b": strconv.Itoa(b),
			"angle": strconv.Itoa(angle), "sin_val": sinVal,
		},
	}
}

// vectorMagnitudeGen measures a scaled triple, so the length the Pythagorean theorem produces is
// exact and the signs of the components are removed by the squaring.
type vectorMagnitudeGen struct{}

func (g *vectorMagnitudeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	t := computedTriples[rand.Intn(len(computedTriples))]
	k := 1 + rand.Intn(2+int(ctx.Difficulty*3))
	x, y := k*t[0], k*t[1]
	if rand.Intn(2) == 0 {
		x = -x
	}
	if rand.Intn(2) == 0 {
		y = -y
	}
	mag := k * t[2]
	return generator.Problem{
		Question: fmt.Sprintf("What is the magnitude of the vector \\(\\langle %d, %d \\rangle\\)?", x, y),
		Answer:   strconv.Itoa(mag),
		Explanation: fmt.Sprintf("The magnitude of a vector is the length of the arrow from the origin to its tip, which is the hypotenuse of a right triangle on its components: \\(|\\langle x, y \\rangle| = \\sqrt{x^2 + y^2}\\). Here that is \\(\\sqrt{%d + %d} = \\sqrt{%d} = %d\\). Squaring first is what makes the components' signs irrelevant: a vector pointing into a different quadrant is the same distance from the origin.",
			x*x, y*y, x*x+y*y, mag),
		Facts: map[string]string{
			"value": strconv.Itoa(mag), "x": strconv.Itoa(x), "y": strconv.Itoa(y),
			"x2": strconv.Itoa(x * x), "y2": strconv.Itoa(y * y), "sum": strconv.Itoa(x*x + y*y),
		},
	}
}

// vectorDotGen multiplies matching components and adds, which needs no triple and no formula: the
// answer is an integer for any pair of integer vectors, and it is negative exactly when the two
// point more away from each other than along.
type vectorDotGen struct{}

func (g *vectorDotGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	rng := 2 + 2*rand.Intn(1+int(ctx.Difficulty*3))
	a := trigSmallInt(-rng, rng)
	b := trigSmallInt(-rng, rng)
	if a == 0 && b == 0 {
		a = 1 // the question asks about two vectors, and one of them has to be there
	}
	c := trigSmallInt(-rng, rng)
	d := trigSmallInt(-rng, rng)
	// A zero dot product is left out: the vectors would be perpendicular, which the schema has no
	// step for. Both of the second vector's components have to move, because a zero product on
	// either side is enough to keep the sum at zero whatever the other one does.
	for i := 0; i < 20 && a*c+b*d == 0; i++ {
		c = trigSmallInt(-rng, rng)
		d = trigSmallInt(-rng, rng)
	}
	t1, t2 := a*c, b*d
	return generator.Problem{
		Question: fmt.Sprintf("What is the dot product \\(\\langle %d, %d \\rangle \\cdot \\langle %d, %d \\rangle\\)?", a, b, c, d),
		Answer:   strconv.Itoa(t1 + t2),
		Explanation: fmt.Sprintf("The dot product multiplies matching components and adds the products: \\(\\langle a, b \\rangle \\cdot \\langle c, d \\rangle = ac + bd\\). Here \\((%d)(%d) + (%d)(%d) = %d + (%d) = %d\\). It is a single number rather than a vector, and its sign says whether the two point broadly the same way or broadly opposite ways.",
			a, c, b, d, t1, t2, t1+t2),
		Facts: map[string]string{
			"value": strconv.Itoa(t1 + t2), "a": strconv.Itoa(a), "b": strconv.Itoa(b),
			"c": strconv.Itoa(c), "d": strconv.Itoa(d), "t1": strconv.Itoa(t1), "t2": strconv.Itoa(t2),
		},
	}
}

// phaseShiftGen reads a horizontal shift off a sinusoid whose constant is a whole multiple of the
// coefficient on x, so the quotient is exact and the answer is an integer number of degrees.
type phaseShiftGen struct{}

func (g *phaseShiftGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	b := 2 + rand.Intn(4+int(ctx.Difficulty*4))
	shift := 5 * (1 + rand.Intn(11))
	c := b * shift
	cc, dir, signed := strconv.Itoa(c), "right", shift
	inner := fmt.Sprintf("%dx - %d^{\\circ}", b, c)
	if rand.Intn(2) == 0 {
		cc, dir, signed = "-"+strconv.Itoa(c), "left", -shift
		inner = fmt.Sprintf("%dx + %d^{\\circ}", b, c)
	}
	a := 2 + rand.Intn(4)
	return generator.Problem{
		Question: fmt.Sprintf("What is the phase shift of \\(y = %d\\sin(%s)\\), in degrees? Give a shift to the right as positive and a shift to the left as negative.", a, inner),
		Answer:   strconv.Itoa(signed),
		Explanation: fmt.Sprintf("The phase shift is the horizontal displacement of the wave, and it is where the argument of the sine reaches zero. Setting \\(%s = 0\\) leaves \\(%dx = %s^{\\circ}\\), so \\(x = %d^{\\circ}\\): the graph is %d degrees to the %s. Dividing the constant by the coefficient on \\(x\\) is what separates the shift from the period, which the coefficient alone sets.",
			inner, b, cc, signed, shift, dir),
		Facts: map[string]string{
			"value": strconv.Itoa(signed), "a": strconv.Itoa(a), "b": strconv.Itoa(b),
			"c_signed": cc, "shift": strconv.Itoa(shift), "dir": dir, "inner": inner,
		},
	}
}
