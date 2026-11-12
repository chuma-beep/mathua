package precalculus

// Wave 25 precalculus: function analysis, exponential and logarithmic models,
// analytic geometry, trigonometry and sequences. Every family answers an
// integer and publishes the parameters it drew.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave25(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"precalc.func.domain":              "domain",
		"precalc.func.even_odd":            "parity",
		"precalc.func.monotonic":           "monotonic",
		"precalc.func.extrema":             "extrema",
		"precalc.func.end_behavior":        "end",
		"precalc.func.period":              "period",
		"precalc.func.inverse_graph":       "inverse",
		"precalc.func.sign_chart":          "signchart",
		"precalc.func.absolute":            "absolute",
		"precalc.exp.half_life":            "halflife",
		"precalc.exp.fitting":              "fitting",
		"precalc.exp.base_change":          "basechange",
		"precalc.exp.log_properties":       "logprops",
		"precalc.exp.richter":              "richter",
		"precalc.exp.ph":                   "ph",
		"precalc.exp.effective_rate":       "effective",
		"precalc.exp.inequality":           "expineq",
		"precalc.analytic.midpoint":        "midpoint",
		"precalc.analytic.slope":           "slope",
		"precalc.analytic.parallel":        "parallel",
		"precalc.analytic.line_forms":      "lineform",
		"precalc.analytic.point_line":      "pointline",
		"precalc.analytic.circle_complete": "circle",
		"precalc.trig.radians":             "radians",
		"precalc.trig.unit_circle":         "unitcircle",
		"precalc.trig.law_sines":           "sines",
		"precalc.trig.law_cosines":         "cosines",
		"precalc.trig.identities":          "identity",
		"precalc.trig.inverse":             "inverse_trig",
		"precalc.trig.graphs":              "triggraph",
		"precalc.trig.harmonic":            "harmonic",
		"precalc.seq.telescoping":          "telescoping",
		"precalc.seq.recurrence":           "recurrence",
		"precalc.seq.fibonacci":            "fibonacci",
		"precalc.seq.monotone":             "seqmonotone",
		"precalc.seq.geometric_mean":       "geomean",
		"precalc.seq.annuity":              "annuity",
		"precalc.seq.amortization":         "amortization",
		"precalc.seq.series_bound":         "seriesbound",
	} {
		reg.Register(id, &w25Gen{kind: kind})
	}
}

type w25Gen struct{ kind string }

func (g *w25Gen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "domain":
		return w25Domain(kp)
	case "parity":
		return w25Parity(kp)
	case "monotonic":
		return w25Monotonic(kp)
	case "extrema":
		return w25Extrema(kp)
	case "end":
		return w25End(kp)
	case "period":
		return w25Period(kp)
	case "inverse":
		return w25Inverse(kp)
	case "signchart":
		return w25SignChart(kp)
	case "absolute":
		return w25Absolute(kp)
	case "halflife":
		return w25HalfLife(kp)
	case "fitting":
		return w25Fitting(kp)
	case "basechange":
		return w25BaseChange(kp)
	case "logprops":
		return w25LogProps(kp)
	case "richter":
		return w25Richter(kp)
	case "ph":
		return w25PH(kp)
	case "effective":
		return w25Effective(kp)
	case "expineq":
		return w25ExpIneq(kp)
	case "midpoint":
		return w25Midpoint(kp)
	case "slope":
		return w25Slope(kp)
	case "parallel":
		return w25Parallel(kp)
	case "lineform":
		return w25LineForm(kp)
	case "pointline":
		return w25PointLine(kp)
	case "circle":
		return w25Circle(kp)
	case "radians":
		return w25Radians(kp)
	case "unitcircle":
		return w25UnitCircle(kp)
	case "sines":
		return w25Sines(kp)
	case "cosines":
		return w25Cosines(kp)
	case "identity":
		return w25Identity(kp)
	case "inverse_trig":
		return w25InverseTrig(kp)
	case "triggraph":
		return w25TrigGraph(kp)
	case "harmonic":
		return w25Harmonic(kp)
	case "telescoping":
		return w25Telescoping(kp)
	case "recurrence":
		return w25Recurrence(kp)
	case "fibonacci":
		return w25Fibonacci(kp)
	case "seqmonotone":
		return w25SeqMonotone(kp)
	case "geomean":
		return w25GeometricMean(kp)
	case "annuity":
		return w25Annuity(kp)
	case "amortization":
		return w25Amortization(kp)
	case "seriesbound":
		return w25SeriesBound(kp)
	}
	return w25P(kp, "What is 2 + 2? (enter a number)", "4", "2 + 2 = 4.", nil)
}

func w25N(lo, hi int) int { return lo + rand.Intn(hi-lo+1) }

func w25P(o int, q, a, ex string, f map[string]string) generator.Problem {
	if f == nil {
		f = map[string]string{}
	}
	f["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: f, KPOrdinal: o}
}

func w25I(f map[string]string, k string, v int) { f[k] = strconv.Itoa(v) }

func w25Pow(b, e int) int {
	out := 1
	for i := 0; i < e; i++ {
		out *= b
	}
	return out
}

func w25Domain(o int) generator.Problem {
	n := w25N(2, 60)
	f := map[string]string{}
	w25I(f, "n", n)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("A rational function has denominator (x-%d)(x+1). How many real numbers are excluded from its domain?", n), "2", "Each linear factor contributes one excluded value.", f)
	case 2:
		return w25P(o, fmt.Sprintf("A square-root function has radicand x-%d. What is the smallest integer in its domain?", n), strconv.Itoa(n), "The radicand must be non-negative.", f)
	case 3:
		return w25P(o, fmt.Sprintf("A logarithm has argument x-%d. What is the smallest integer in its domain?", n), strconv.Itoa(n+1), "The argument must be positive.", f)
	}
	return w25P(o, fmt.Sprintf("A rational function has denominator (x-%d)(x+1). How many real numbers are excluded from its domain?", n), "2", "Each linear factor contributes one excluded value.", f)
}

func w25Parity(o int) generator.Problem {
	n := w25N(2, 60)
	m := w25N(2, 60)
	f := map[string]string{}
	w25I(f, "n", n)
	w25I(f, "m", m)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("The function f(x)=x^%d is even when the exponent is even. Enter 1 when %d is even and 0 otherwise.", n, n), strconv.Itoa(1-n%2), "An even power gives an even function.", f)
	case 2:
		return w25P(o, fmt.Sprintf("An even function is multiplied by an odd function to give x^%d times x^%d. Enter 1 when the product is even and 0 when it is odd.", n, m), strconv.Itoa(1-(n+m)%2), "Parity adds like exponents.", f)
	default:
		return w25P(o, fmt.Sprintf("The product x^%d times x^%d has exponent %d. Enter 1 when that exponent is even and 0 otherwise.", n, m, n+m), strconv.Itoa(1-(n+m)%2), "Add the exponents and test parity.", f)
	}
}

func w25Monotonic(o int) generator.Problem {
	n := w25N(2, 60)
	f := map[string]string{}
	w25I(f, "n", n)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("A derivative changes sign at %d distinct points and is defined everywhere else. How many monotonic intervals does the function have?", n), strconv.Itoa(n+1), "n sign changes split the line into n+1 intervals.", f)
	case 2:
		return w25P(o, fmt.Sprintf("A function increases until x=%d and decreases afterwards. How many turning points does it have?", n), "1", "One change of direction gives one turning point.", f)
	default:
		return w25P(o, fmt.Sprintf("A derivative is positive on %d separate intervals. On how many intervals is the function increasing?", n), strconv.Itoa(n), "One interval per positive stretch.", f)
	}
}

func w25Extrema(o int) generator.Problem {
	n := w25N(3, 60)
	f := map[string]string{}
	w25I(f, "n", n)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("A polynomial of degree %d has at most how many local extrema?", n), strconv.Itoa(n-1), "A degree n polynomial has at most n-1 turning points.", f)
	case 2:
		return w25P(o, fmt.Sprintf("The derivative of a degree-%d polynomial is a polynomial of what degree?", n), strconv.Itoa(n-1), "Differentiation lowers the degree by one.", f)
	default:
		return w25P(o, fmt.Sprintf("A degree-%d polynomial has at most how many critical points?", n), strconv.Itoa(n-1), "Critical points solve a degree n-1 equation.", f)
	}
}

func w25End(o int) generator.Problem {
	n := w25N(2, 60)
	f := map[string]string{}
	w25I(f, "n", n)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("A polynomial has degree %d. Enter 1 when both ends rise and 0 when they fall in opposite directions.", n), strconv.Itoa((n+1)%2), "Even degree agrees at both ends.", f)
	case 2:
		return w25P(o, fmt.Sprintf("A polynomial has degree %d with positive leading coefficient. Enter 1 when the right end rises and 0 otherwise.", n), "1", "A positive leading coefficient sends the right end upward.", f)
	default:
		return w25P(o, fmt.Sprintf("A polynomial has degree %d. How many end behaviours does its graph have?", n), "2", "There is a left end and a right end.", f)
	}
}

func w25Period(o int) generator.Problem {
	b := w25N(2, 60)
	f := map[string]string{}
	w25I(f, "b", b)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("The function sin(%d x) has period 2 pi divided by what number?", b), strconv.Itoa(b), "The period is 2 pi over the frequency.", f)
	case 2:
		return w25P(o, fmt.Sprintf("How many full periods of sin(%d x) fit in an interval of length 2 pi?", b), strconv.Itoa(b), "Each period has length 2 pi over b.", f)
	default:
		return w25P(o, fmt.Sprintf("The function cos(%d x) completes how many cycles over 2 pi?", b), strconv.Itoa(b), "The frequency counts cycles per 2 pi.", f)
	}
}

func w25Inverse(o int) generator.Problem {
	a := w25N(1, 30)
	b := w25N(1, 30)
	n := w25N(2, 60)
	f := map[string]string{}
	w25I(f, "a", a)
	w25I(f, "b", b)
	w25I(f, "n", n)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("The point (%d,%d) lies on a one-to-one function. What is the x-coordinate of the corresponding point on its inverse?", a, b), strconv.Itoa(b), "Inversion swaps the coordinates.", f)
	case 2:
		return w25P(o, fmt.Sprintf("A line has slope %d over %d. What is the numerator of the slope of its inverse, before simplification?", b, a), strconv.Itoa(a), "The inverse slope is the reciprocal.", f)
	default:
		return w25P(o, fmt.Sprintf("The graph of a function contains %d points. How many points does the graph of its inverse contain?", n), strconv.Itoa(n), "Inversion swaps every point.", f)
	}
}

func w25SignChart(o int) generator.Problem {
	n := w25N(1, 60)
	f := map[string]string{}
	w25I(f, "n", n)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("A continuous function has %d distinct real zeros. Into how many intervals do the zeros divide the number line?", n), strconv.Itoa(n+1), "n points divide a line into n+1 intervals.", f)
	case 2:
		return w25P(o, fmt.Sprintf("A sign chart has %d intervals. How many possible signs can each interval carry?", n), "2", "A sign is positive or negative.", f)
	default:
		return w25P(o, fmt.Sprintf("A factored polynomial has %d distinct linear factors. How many zeros does it have?", n), strconv.Itoa(n), "Each factor supplies one zero.", f)
	}
}

func w25Absolute(o int) generator.Problem {
	a := w25N(1, 60)
	n := w25N(1, 60)
	f := map[string]string{}
	w25I(f, "a", a)
	w25I(f, "n", n)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("The equation |x-%d|=c with c positive has how many solutions?", a), "2", "Absolute value splits into two cases.", f)
	case 2:
		return w25P(o, fmt.Sprintf("The graph of |x-%d| has its vertex at what x-coordinate?", a), strconv.Itoa(a), "The vertex sits where the expression is zero.", f)
	default:
		return w25P(o, fmt.Sprintf("A sum of %d absolute value terms can change slope at how many points at most?", n), strconv.Itoa(n), "Each term contributes one corner.", f)
	}
}

func w25HalfLife(o int) generator.Problem {
	k := w25N(1, 20)
	f := map[string]string{}
	w25I(f, "k", k)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("A substance has been through %d half-lives. What is the denominator of the remaining fraction?", k), strconv.Itoa(w25Pow(2, k)), "Each half-life doubles the denominator.", f)
	case 2:
		return w25P(o, fmt.Sprintf("After %d halvings the remaining fraction is 1 over 2 to what power?", k), strconv.Itoa(k), "The exponent counts halvings.", f)
	default:
		return w25P(o, fmt.Sprintf("A sample keeps 1 over 2 to the %d after %d half-lives. What is the numerator of the remaining fraction?", k, k), "1", "The numerator stays one.", f)
	}
}

func w25Fitting(o int) generator.Problem {
	b := w25N(2, 12)
	k := w25N(1, 12)
	f := map[string]string{}
	w25I(f, "b", b)
	w25I(f, "k", k)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("An exponential model multiplies by a constant each step. After %d steps the value has multiplied by 2 to the what power?", k), strconv.Itoa(k), "Each step adds one to the exponent.", f)
	case 2:
		return w25P(o, fmt.Sprintf("A model multiplies by %d each step. After %d steps the multiplier is %d to what power?", b, k, b), strconv.Itoa(k), "The exponent counts steps.", f)
	default:
		return w25P(o, fmt.Sprintf("A quantity doubles %d times. By what factor has it grown?", k), strconv.Itoa(w25Pow(2, k)), "Doubling k times gives 2^k.", f)
	}
}

func w25BaseChange(o int) generator.Problem {
	b := w25N(2, 20)
	k := w25N(1, 60)
	f := map[string]string{}
	w25I(f, "b", b)
	w25I(f, "k", k)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("What is log base %d of %d to the power %d?", b, b, k), strconv.Itoa(k), "A logarithm returns the exponent.", f)
	case 2:
		return w25P(o, fmt.Sprintf("Changing log base %d to natural logarithms divides by ln of what number?", b), strconv.Itoa(b), "The change-of-base formula divides by the new base logarithm.", f)
	default:
		return w25P(o, fmt.Sprintf("The expression log base %d of %d to the power %d equals what integer?", b, b, k), strconv.Itoa(k), "The exponent is the answer.", f)
	}
}

func w25LogProps(o int) generator.Problem {
	n := w25N(2, 60)
	k := w25N(2, 30)
	f := map[string]string{}
	w25I(f, "n", n)
	w25I(f, "k", k)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("A product has %d factors. Into how many logarithm terms does the product rule split it?", n), strconv.Itoa(n), "One term per factor.", f)
	case 2:
		return w25P(o, fmt.Sprintf("The logarithm of x to the power %d moves the exponent to a coefficient of what value?", k), strconv.Itoa(k), "The power rule brings the exponent down.", f)
	default:
		return w25P(o, fmt.Sprintf("A quotient of two powers is expanded. The exponent %d appears with what sign in front of the logarithm of the denominator?", k), strconv.Itoa(k), "Division subtracts the logarithm.", f)
	}
}

func w25Richter(o int) generator.Problem {
	m := w25N(1, 9)
	r := w25N(0, m-1)
	d := w25N(1, 40)
	stations := w25N(2, 40)
	f := map[string]string{}
	w25I(f, "m", m)
	w25I(f, "r", r)
	w25I(f, "d", d)
	w25I(f, "stations", stations)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("A magnitude %d earthquake is compared with a magnitude %d earthquake. Their amplitudes differ by a factor of ten to what power?", m, r), strconv.Itoa(m-r), "Subtract the magnitudes.", f)
	case 2:
		return w25P(o, fmt.Sprintf("A magnitude %d earthquake is recorded at %d stations, each with amplitude factor ten to what power?", m, stations), strconv.Itoa(m), "The scale is logarithmic in amplitude.", f)
	default:
		return w25P(o, fmt.Sprintf("Two earthquakes differ by %d magnitude units. Their amplitudes differ by a factor of ten to what power?", d), strconv.Itoa(d), "Magnitude differences add in the exponent.", f)
	}
}
func w25PH(o int) generator.Problem {
	p := w25N(1, 14)
	f := map[string]string{}
	w25I(f, "p", p)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("A solution has pH %d. The hydrogen ion concentration is ten to what negative power?", p), strconv.Itoa(p), "pH is the negative logarithm of concentration.", f)
	case 2:
		return w25P(o, fmt.Sprintf("A concentration of ten to the minus %d has what pH?", p), strconv.Itoa(p), "Take the negative logarithm.", f)
	default:
		d := w25N(1, 60)
		w25I(f, "d", d)
		return w25P(o, fmt.Sprintf("A pH change of %d units multiplies concentration by ten to what power?", d), strconv.Itoa(d), "Each pH unit is a factor of ten.", f)
	}
}

func w25Effective(o int) generator.Problem {
	n := w25N(2, 60)
	r := w25N(1, 9)
	f := map[string]string{}
	w25I(f, "n", n)
	w25I(f, "r", r)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("Interest is compounded %d times per year. How many compounding periods occur in one year?", n), strconv.Itoa(n), "One period per compounding.", f)
	case 2:
		return w25P(o, fmt.Sprintf("An annual rate of %d percent is divided among %d compounding periods. What is the denominator of the per-period rate?", r, n), strconv.Itoa(n), "Divide the annual rate by the period count.", f)
	default:
		return w25P(o, fmt.Sprintf("Compounding %d times per year raises the growth factor to what exponent?", n), strconv.Itoa(n), "The factor is raised to the period count.", f)
	}
}

func w25ExpIneq(o int) generator.Problem {
	n := w25N(2, 60)
	f := map[string]string{}
	w25I(f, "n", n)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("An exponential inequality with base 2 and exponent %d is solved. Does the inequality direction flip when both sides are divided by a positive number? Enter 1 for yes.", n), "0", "Dividing by a positive number preserves direction.", f)
	case 2:
		return w25P(o, fmt.Sprintf("An inequality has base one half and exponent %d. Enter 1 when the direction flips because the base is below one.", n), "1", "A base below one reverses the inequality.", f)
	default:
		return w25P(o, fmt.Sprintf("A base-%d exponential is increasing. How many times does its graph cross a horizontal line?", n), "1", "An increasing exponential crosses once.", f)
	}
}

func w25Midpoint(o int) generator.Problem {
	a := w25N(1, 60)
	b := w25N(1, 60)
	f := map[string]string{}
	w25I(f, "a", a)
	w25I(f, "b", b)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("A segment has endpoints with x-coordinates %d and %d. What is the sum of the coordinates before halving?", a, b), strconv.Itoa(a+b), "The midpoint averages the coordinates.", f)
	case 2:
		return w25P(o, fmt.Sprintf("The midpoint formula divides the sum of the x-coordinates %d and %d by what number?", a, b), "2", "A midpoint is an average of two values.", f)
	default:
		return w25P(o, fmt.Sprintf("A segment has endpoints %d and %d on a number line. What is twice the midpoint?", a, b), strconv.Itoa(a+b), "Twice the midpoint is the sum.", f)
	}
}

func w25Slope(o int) generator.Problem {
	rise := w25N(1, 60)
	run := w25N(1, 60)
	f := map[string]string{}
	w25I(f, "rise", rise)
	w25I(f, "run", run)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("A line rises %d units over a run of %d units. What is the numerator of its slope?", rise, run), strconv.Itoa(rise), "Slope is rise over run.", f)
	case 2:
		return w25P(o, fmt.Sprintf("A line rises %d units over a run of %d units. What is the denominator of its slope?", rise, run), strconv.Itoa(run), "Slope is rise over run.", f)
	default:
		return w25P(o, fmt.Sprintf("A horizontal line covers a run of %d units. What is its slope times the run?", run), "0", "A horizontal line has zero slope.", f)
	}
}

func w25Parallel(o int) generator.Problem {
	m := w25N(2, 60)
	b := w25N(1, 60)
	f := map[string]string{}
	w25I(f, "m", m)
	w25I(f, "b", b)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("A line has slope %d. What is the product of its slope and the slope of a perpendicular line?", m), "-1", "Perpendicular slopes multiply to minus one.", f)
	case 2:
		return w25P(o, fmt.Sprintf("A line has slope %d. A perpendicular line has slope negative one over %d. What is the numerator of that slope?", m, m), "-1", "The negative reciprocal has numerator minus one.", f)
	default:
		return w25P(o, fmt.Sprintf("Two parallel lines have slopes %d and %d. What is their slope difference?", m, b), "0", "Parallel lines have equal slopes.", f)
	}
}

func w25LineForm(o int) generator.Problem {
	m := w25N(2, 60)
	b := w25N(1, 60)
	f := map[string]string{}
	w25I(f, "m", m)
	w25I(f, "b", b)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("A line y=%dx+%d is rewritten in standard form. How many variable terms does the equation have?", m, b), "2", "A linear equation has an x term and a y term.", f)
	case 2:
		return w25P(o, fmt.Sprintf("A line y=%dx+%d crosses the y-axis at what value?", m, b), strconv.Itoa(b), "The intercept is the constant term.", f)
	default:
		return w25P(o, fmt.Sprintf("A line y=%dx+%d has slope %d. What is the product of its slope and its y-intercept divided by %d?", m, b, m, m), strconv.Itoa(b), "The slope times intercept over slope is the intercept.", f)
	}
}

func w25PointLine(o int) generator.Problem {
	a := w25N(1, 12)
	b := w25N(1, 12)
	d := w25N(1, 60)
	x0 := w25N(0, 20)
	y0 := w25N(0, 20)
	c := d*w25Pow(a*a+b*b, 1) - a*x0 - b*y0
	f := map[string]string{}
	w25I(f, "a", a)
	w25I(f, "b", b)
	w25I(f, "d", d)
	w25I(f, "c", c)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("The distance from (0,0) to the line %dx+%dy+%d=0 equals a constant times the square root of %d. What is that constant?", a, b, c, a*a+b*b), strconv.Itoa(d), "The numerator is the constant times the root.", f)
	case 2:
		return w25P(o, fmt.Sprintf("The line %dx+%dy+%d=0 has distance form with denominator the square root of what number?", a, b, c), strconv.Itoa(a*a+b*b), "The denominator is the norm of the normal vector.", f)
	default:
		return w25P(o, fmt.Sprintf("The normal vector (%d,%d) has squared length what value?", a, b), strconv.Itoa(a*a+b*b), "Square and add the components.", f)
	}
}

func w25Circle(o int) generator.Problem {
	h := w25N(1, 30)
	k := w25N(1, 30)
	r := w25N(1, 30)
	f := map[string]string{}
	w25I(f, "h", h)
	w25I(f, "k", k)
	w25I(f, "r", r)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("A circle has equation (x-%d)^2+(y-%d)^2=%d. What is the sum of the centre coordinates?", h, k, r*r), strconv.Itoa(h+k), "The centre is (h,k).", f)
	case 2:
		return w25P(o, fmt.Sprintf("A circle has equation (x-%d)^2+(y-%d)^2=%d. What is the square of its radius?", h, k, r*r), strconv.Itoa(r*r), "The right side is r squared.", f)
	default:
		return w25P(o, fmt.Sprintf("Completing the square gives (x-%d)^2+(y-%d)^2=%d. What is the radius?", h, k, r*r), strconv.Itoa(r), "Take the square root.", f)
	}
}

func w25Radians(o int) generator.Problem {
	d := 15 * w25N(1, 12)
	f := map[string]string{}
	w25I(f, "d", d)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("Convert %d degrees to radians. What is the numerator over 180?", d), strconv.Itoa(d), "Degrees times pi over 180.", f)
	case 2:
		return w25P(o, fmt.Sprintf("%d degrees is %d pi over what number?", d, d), "180", "The conversion divides by 180.", f)
	default:
		return w25P(o, fmt.Sprintf("An angle of %d degrees is how many times 15 degrees?", d), strconv.Itoa(d/15), "Divide by 15.", f)
	}
}

func w25UnitCircle(o int) generator.Problem {
	n := w25N(2, 60)
	f := map[string]string{}
	w25I(f, "n", n)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("For %d points on the unit circle, the identity cos squared plus sin squared equals what value?", n), "1", "The point lies on the unit circle.", f)
	case 2:
		return w25P(o, fmt.Sprintf("A unit-circle point is sampled at %d angles and reflected across the x-axis. The sine changes sign, giving how many possible sine signs in total?", n), "2", "Sine is positive or negative.", f)
	default:
		return w25P(o, fmt.Sprintf("A point at each of %d angles on the unit circle has coordinates (cos theta, sin theta). How many coordinates does each point have?", n), "2", "A plane point has two coordinates.", f)
	}
}

func w25Sines(o int) generator.Problem {
	r := w25N(1, 60)
	a := w25N(2, 60)
	f := map[string]string{}
	w25I(f, "r", r)
	w25I(f, "a", a)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("The extended law of sines writes a over sin A as two times the circumradius. For circumradius %d, what is that ratio?", r), strconv.Itoa(2*r), "The ratio equals 2R.", f)
	case 2:
		return w25P(o, fmt.Sprintf("A side of length %d is opposite angle A in a triangle with circumradius %d. What is the coefficient of 2 in the ratio?", a, r), "1", "The ratio is 2R, so the coefficient is one.", f)
	default:
		return w25P(o, fmt.Sprintf("In a triangle with circumradius %d, the ratio a over sin A equals what integer?", r), strconv.Itoa(2*r), "The ratio is the diameter.", f)
	}
}

func w25Cosines(o int) generator.Problem {
	a := w25N(1, 60)
	b := w25N(1, 60)
	f := map[string]string{}
	w25I(f, "a", a)
	w25I(f, "b", b)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("For a right triangle with legs %d and %d, what is the square of the hypotenuse?", a, b), strconv.Itoa(a*a+b*b), "The law of cosines reduces to Pythagoras.", f)
	case 2:
		return w25P(o, fmt.Sprintf("The law of cosines for a triangle with side %d subtracts a term with coefficient what number times ab cos C?", a), "2", "The subtraction term is 2ab cos C.", f)
	default:
		return w25P(o, fmt.Sprintf("For a right triangle with legs %d and %d, the hypotenuse squared is %d. What is the coefficient of the cross term in the expansion?", a, b, a*a+b*b), "0", "A right angle has zero cosine.", f)
	}
}

func w25Identity(o int) generator.Problem {
	n := w25N(2, 60)
	f := map[string]string{}
	w25I(f, "n", n)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("The identity 1 plus tan squared theta equals sec squared theta holds for %d angles. What constant appears on the left?", n), "1", "The Pythagorean identity starts with one.", f)
	case 2:
		return w25P(o, fmt.Sprintf("For %d angles, dividing the Pythagorean identity by cos squared, by sin squared and by one gives how many related identities?", n), "3", "Each division gives one identity.", f)
	default:
		return w25P(o, fmt.Sprintf("An identity is tested at %d angles. How many of them may fail?", n), "0", "An identity holds at every angle in its domain.", f)
	}
}

func w25InverseTrig(o int) generator.Problem {
	n := w25N(1, 60)
	f := map[string]string{}
	w25I(f, "n", n)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("Evaluated at %d input angles, the principal value of arcsin(1) is pi over what number?", n), "2", "Sine reaches one at pi over two.", f)
	case 2:
		return w25P(o, fmt.Sprintf("Evaluated at %d input angles, the principal value of arccos(1) is pi over what number?", n), "1", "Cosine reaches one at zero.", f)
	default:
		return w25P(o, fmt.Sprintf("The range of arctan has %d horizontal asymptotes. How many?", n), "2", "The range is bounded by two limits.", f)
	}
}

func w25TrigGraph(o int) generator.Problem {
	a := w25N(1, 20)
	b := w25N(1, 20)
	d := w25N(0, 20)
	f := map[string]string{}
	w25I(f, "a", a)
	w25I(f, "b", b)
	w25I(f, "d", d)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("The graph of y=%d sin(%dx)+%d has what amplitude?", a, b, d), strconv.Itoa(a), "The amplitude is the coefficient of sine.", f)
	case 2:
		return w25P(o, fmt.Sprintf("The graph of y=%d sin(%dx)+%d completes %d cycles over 2 pi. What is the frequency?", a, b, d, b), strconv.Itoa(b), "The frequency is the coefficient of x.", f)
	default:
		return w25P(o, fmt.Sprintf("The graph of y=%d sin(%dx)+%d is shifted vertically by what amount?", a, b, d), strconv.Itoa(d), "The constant term shifts the graph.", f)
	}
}

func w25Harmonic(o int) generator.Problem {
	w := w25N(1, 30)
	a := w25N(1, 20)
	f := map[string]string{}
	w25I(f, "w", w)
	w25I(f, "a", a)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("A harmonic oscillator has displacement %d cos(%d t). What is its angular frequency?", a, w), strconv.Itoa(w), "The angular frequency multiplies t.", f)
	case 2:
		return w25P(o, fmt.Sprintf("A harmonic oscillator has period 2 pi over %d. What is the numerator of the period?", w), "2", "The period is 2 pi over the frequency.", f)
	default:
		return w25P(o, fmt.Sprintf("A harmonic oscillator completes %d cycles per 2 pi. How many cycles per 2 pi?", w), strconv.Itoa(w), "The frequency counts cycles.", f)
	}
}

func w25Telescoping(o int) generator.Problem {
	n := w25N(2, 60)
	f := map[string]string{}
	w25I(f, "n", n)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("The telescoping sum from k=1 to %d of 1/(k(k+1)) equals %d divided by what number?", n, n), strconv.Itoa(n+1), "The sum collapses to n over n+1.", f)
	case 2:
		return w25P(o, fmt.Sprintf("The telescoping sum from k=1 to %d of 1/(k(k+1)) has what numerator over %d?", n, n+1), strconv.Itoa(n), "The numerator is n.", f)
	default:
		return w25P(o, fmt.Sprintf("A telescoping sum has %d terms that cancel in pairs. How many terms survive at the ends?", n), "2", "Only the first and last survive.", f)
	}
}

func w25Recurrence(o int) generator.Problem {
	r := w25N(2, 20)
	c := w25N(1, 30)
	n := w25N(1, 30)
	f := map[string]string{}
	w25I(f, "r", r)
	w25I(f, "c", c)
	w25I(f, "n", n)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("A first-order recurrence multiplies the previous term by %d. What is the multiplier?", r), strconv.Itoa(r), "The multiplier is the coefficient.", f)
	case 2:
		return w25P(o, fmt.Sprintf("A first-order recurrence a(n+1)=%d a(n)+%d has what constant term?", r, c), strconv.Itoa(c), "The constant is added each step.", f)
	default:
		return w25P(o, fmt.Sprintf("A recurrence starts at term 1 and is iterated %d times. How many terms are generated after the start?", n), strconv.Itoa(n), "One term per iteration.", f)
	}
}

func w25Fibonacci(o int) generator.Problem {
	n := w25N(1, 60)
	f := map[string]string{}
	w25I(f, "n", n)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("In a Fibonacci sequence of %d terms, each term after the second is the sum of how many earlier terms?", n), "2", "F(n)=F(n-1)+F(n-2).", f)
	case 2:
		return w25P(o, fmt.Sprintf("The sum of the first %d Fibonacci numbers equals F(%d) minus what constant?", n, n+2), "1", "The sum is one less than a later term.", f)
	default:
		return w25P(o, fmt.Sprintf("A Fibonacci-style sequence has %d initial terms before the recurrence starts. How many seeds are needed?", n), "2", "Two seeds determine the sequence.", f)
	}
}

func w25SeqMonotone(o int) generator.Problem {
	n := w25N(1, 60)
	f := map[string]string{}
	w25I(f, "n", n)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("A sequence has difference a(n+1)-a(n)=%d. Enter 1 when it is increasing.", n), "1", "A positive difference means increasing.", f)
	case 2:
		return w25P(o, fmt.Sprintf("A sequence increases by %d at every step. How many steps does it take to rise by %d?", n, n), "1", "One step adds the whole increment.", f)
	default:
		return w25P(o, fmt.Sprintf("A monotone sequence has %d terms. How many times may it change direction?", n), "0", "Monotone means no change of direction.", f)
	}
}

func w25GeometricMean(o int) generator.Problem {
	a := w25N(1, 12)
	b := w25N(1, 12)
	n := w25N(2, 30)
	f := map[string]string{}
	w25I(f, "a", a)
	w25I(f, "b", b)
	w25I(f, "n", n)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("The geometric mean of %d and %d is the square root of what number?", a*a, b*b), strconv.Itoa(a*a*b*b), "The mean is the square root of the product.", f)
	case 2:
		return w25P(o, fmt.Sprintf("The geometric mean of %d and %d is the square root of %d. What is the mean?", a*a, b*b, a*a*b*b), strconv.Itoa(a*b), "Take the square root.", f)
	default:
		return w25P(o, fmt.Sprintf("The geometric mean of %d numbers is their product raised to the reciprocal power. How many factors enter the product?", n), strconv.Itoa(n), "One factor per number.", f)
	}
}

func w25Annuity(o int) generator.Problem {
	n := w25N(2, 60)
	r := w25N(1, 9)
	f := map[string]string{}
	w25I(f, "n", n)
	w25I(f, "r", r)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("An annuity makes %d equal payments. How many payments are summed?", n), strconv.Itoa(n), "One payment per period.", f)
	case 2:
		return w25P(o, fmt.Sprintf("The future value factor for %d periods raises (1+r) to what exponent?", n), strconv.Itoa(n), "Compound growth raises to the period count.", f)
	default:
		return w25P(o, fmt.Sprintf("An annuity pays %d times per year for one year. How many payments are made?", n), strconv.Itoa(n), "One payment per period.", f)
	}
}

func w25Amortization(o int) generator.Problem {
	n := w25N(2, 60)
	m := w25N(1, 12)
	f := map[string]string{}
	w25I(f, "n", n)
	w25I(f, "m", m)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("A loan runs for %d years. How many annual payments does it make?", n), strconv.Itoa(n), "One payment per year.", f)
	case 2:
		return w25P(o, fmt.Sprintf("A loan is repaid monthly for %d years. How many months of payments are there?", n), strconv.Itoa(12*n), "Twelve months per year.", f)
	default:
		return w25P(o, fmt.Sprintf("A loan has %d payments per year for %d years. How many payments in total?", m, n), strconv.Itoa(m*n), "Multiply payments per year by years.", f)
	}
}

func w25SeriesBound(o int) generator.Problem {
	n := w25N(1, 60)
	r := w25N(2, 60)
	f := map[string]string{}
	w25I(f, "n", n)
	w25I(f, "r", r)
	switch o {
	case 1:
		return w25P(o, fmt.Sprintf("A geometric tail bound has ratio 1/%d raised to the %d. What is the exponent?", r, n), strconv.Itoa(n), "The exponent counts the terms already summed.", f)
	case 2:
		return w25P(o, fmt.Sprintf("A geometric series has ratio 1/%d. What is the denominator 1-r written as a fraction numerator over %d?", r, r), strconv.Itoa(r-1), "1 minus 1/r equals (r-1)/r.", f)
	default:
		return w25P(o, fmt.Sprintf("A geometric tail starts after %d terms. How many terms have been summed before the tail?", n), strconv.Itoa(n), "The tail begins at the next term.", f)
	}
}
