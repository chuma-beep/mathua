package trigonometry

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

// The four graph concepts asked three fixed questions each, in prose, with the
// answer written into the option list. Period, asymptotes, amplitude and range
// are all functions of the coefficients the graph is written with, so the shapes
// below draw those coefficients and compute the answer instead of reciting it.
// Each pool keeps one definitional question, because a learner still has to know
// what a reciprocal is.

// degStr renders an angle in degrees without a trailing .0.
func degStr(v float64) string {
	return strconv.FormatFloat(math.Round(v*100)/100, 'f', -1, 64)
}

// distinctDistractors returns three options that differ from the answer and from
// each other.
//
// trigPickMC hands the candidates to generator.MultipleChoice, which caps the
// pool at four letters and drops duplicates. A claim whose candidates collide,
// for example an amplitude equal to the period coefficient, would ship three
// options, and the batch test requires four.
func distinctDistractors(correct string, candidates ...string) []string {
	seen := map[string]bool{correct: true}
	out := make([]string, 0, 3)
	consider := func(candidate string) {
		if len(out) == 3 || seen[candidate] {
			return
		}
		seen[candidate] = true
		out = append(out, candidate)
	}
	for _, candidate := range candidates {
		consider(candidate)
	}
	// These are wrong for every shape here, so they restore a fourth option
	// when the collisions above remove one, and a fifth when they remove two.
	for _, fallback := range []string{"0", "180", "360", "-1", "[0, 0]"} {
		consider(fallback)
	}
	return out
}

type tanGraphGen struct{}

func (g *tanGraphGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	k := 2 + rand.Intn(5) // 2..6, and every one divides 180
	period := degStr(180.0 / float64(k))
	first := degStr(90.0 / float64(k))
	step := degStr(180.0 / float64(k))
	switch rand.Intn(4) {
	case 0:
		return trigPickMC([]trigMCClaim{{
			q:       fmt.Sprintf("The period of the graph of y = tan(%dx) is how many degrees?", k),
			correct: period,
			distractors: distinctDistractors(period,
				degStr(180.0*float64(k)),
				degStr(360.0/float64(k)),
				first,
				degStr(90.0*float64(k)),
			),
			reason: fmt.Sprintf("Compressing the input by %d repeats the pattern %d times as often, so the period is 180 divided by %d, which is %s degrees.", k, k, k, period),
		}})
	case 1:
		return trigPickMC([]trigMCClaim{{
			q:       fmt.Sprintf("Where does the graph of y = tan(%dx) have its vertical asymptotes?", k),
			correct: fmt.Sprintf("x = %s + %sj", first, step),
			distractors: distinctDistractors(
				fmt.Sprintf("x = %s + %sj", first, step),
				fmt.Sprintf("x = %s + %sj", step, step),
				fmt.Sprintf("x = %s + 180j", first),
				fmt.Sprintf("x = %s + %sj", degStr(90.0*float64(k)), step),
			),
			reason: fmt.Sprintf("tan(%dx) is undefined where cos(%dx) is zero, so %dx must be an odd multiple of 90 degrees and x = %s + %sj.", k, k, k, first, step),
		}})
	case 2:
		return trigPickMC([]trigMCClaim{{
			q:       fmt.Sprintf("As x approaches %s degrees from below, tan(%dx)...", first, k),
			correct: "increases without bound",
			distractors: distinctDistractors("increases without bound",
				"approaches zero",
				"decreases without bound",
				"approaches one",
			),
			reason: fmt.Sprintf("The sine stays near 1 while the cosine of %dx approaches zero through positive values, so the quotient grows without bound.", k),
		}})
	default:
		return trigPickMC([]trigMCClaim{{
			q:       fmt.Sprintf("On the interval from 0 to 180 degrees, how many complete periods does the graph of y = tan(%dx) show?", k),
			correct: strconv.Itoa(k),
			distractors: distinctDistractors(strconv.Itoa(k),
				strconv.Itoa(2*k),
				strconv.Itoa(k+1),
				"1",
				"180",
			),
			reason: fmt.Sprintf("One period occupies 180 divided by %d degrees, which fits %d times into 180 degrees.", k, k),
		}})
	}
}

// sineGraph renders a sin(b(x - c)) + d the way a textbook writes it: no
// coefficient of one, no "x - -1" when the shift is negative, and no "+ 0".
func sineGraph(a, b, c, d int) string {
	inner := "x"
	switch {
	case c > 0:
		inner = fmt.Sprintf("x - %d", c)
	case c < 0:
		inner = fmt.Sprintf("x + %d", -c)
	}
	arg := inner
	if b != 1 {
		arg = fmt.Sprintf("%d(%s)", b, inner)
	}
	body := fmt.Sprintf("sin(%s)", arg)
	if a != 1 {
		body = fmt.Sprintf("%d %s", a, body)
	}
	switch {
	case d > 0:
		body = fmt.Sprintf("%s + %d", body, d)
	case d < 0:
		body = fmt.Sprintf("%s - %d", body, -d)
	}
	return "y = " + body
}

type graphTransformGen struct{}

func (g *graphTransformGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := 1 + rand.Intn(4)
	b := 1 + rand.Intn(4) // 1..4, every one divides 360
	c := rand.Intn(7) - 3
	d := rand.Intn(11) - 5
	period := degStr(360.0 / float64(b))
	amp := strconv.Itoa(a)
	shift := fmt.Sprintf("%+d", d)
	low := d - a
	high := d + a
	switch rand.Intn(5) {
	case 0:
		return trigPickMC([]trigMCClaim{{
			q:       fmt.Sprintf("For %s, the amplitude of the graph is...", sineGraph(a, b, c, d)),
			correct: amp,
			distractors: distinctDistractors(amp,
				strconv.Itoa(a*b),
				period,
				shift,
				strconv.Itoa(a+d),
			),
			reason: fmt.Sprintf("Amplitude is the size of the coefficient in front of the sine, so it is %d regardless of the period and the shift.", a),
		}})
	case 1:
		return trigPickMC([]trigMCClaim{{
			q:       fmt.Sprintf("For %s, the period of the graph is how many degrees?", sineGraph(a, b, c, d)),
			correct: period,
			distractors: distinctDistractors(period,
				strconv.Itoa(360*b),
				degStr(360.0/float64(a)),
				degStr(180.0/float64(b)),
				strconv.Itoa(360+b),
			),
			reason: fmt.Sprintf("The coefficient %d compresses the input, so the period is 360 divided by %d, which is %s degrees, while the amplitude and the shift leave it alone.", b, b, period),
		}})
	case 2:
		return trigPickMC([]trigMCClaim{{
			q:       fmt.Sprintf("For %s, the vertical shift of the graph is...", sineGraph(a, b, c, d)),
			correct: shift,
			distractors: distinctDistractors(shift,
				fmt.Sprintf("%+d", c),
				fmt.Sprintf("%+d", a),
				fmt.Sprintf("%+d", b),
				fmt.Sprintf("%+d", d+a),
			),
			reason: fmt.Sprintf("The constant added after the sine moves every point by that amount, so the shift is %+d, while %+d is the horizontal shift.", d, c),
		}})
	case 3:
		return trigPickMC([]trigMCClaim{{
			q:       fmt.Sprintf("For %s, the range of the graph is...", sineGraph(a, b, c, d)),
			correct: fmt.Sprintf("[%d, %d]", low, high),
			distractors: distinctDistractors(fmt.Sprintf("[%d, %d]", low, high),
				fmt.Sprintf("[%d, %d]", d, d+a),
				fmt.Sprintf("[%d, %d]", low-a, high+a),
				fmt.Sprintf("[%d, %d]", c, c+a),
			),
			reason: fmt.Sprintf("The sine runs from -1 to 1, so %d sin(...) runs from %d to %d, and adding %d moves that interval to [%d, %d].", a, -a, a, d, low, high),
		}})
	default:
		return trigPickMC([]trigMCClaim{{
			q:       fmt.Sprintf("For %s, the horizontal shift of the graph is...", sineGraph(a, b, c, d)),
			correct: fmt.Sprintf("%+d", c),
			distractors: distinctDistractors(fmt.Sprintf("%+d", c),
				fmt.Sprintf("%+d", d),
				fmt.Sprintf("%+d", c*b),
				fmt.Sprintf("%+d", -c),
				fmt.Sprintf("%+d", a),
			),
			reason: fmt.Sprintf("The form b(x - c) marks the shift inside the brackets, so the graph moves %+d horizontally, and dividing by %d is already accounted for.", c, b),
		}})
	}
}

type reciprocalGraphGen struct{}

func (g *reciprocalGraphGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	k := 1 + rand.Intn(4)
	a := 1 + rand.Intn(4)
	period := degStr(360.0 / float64(k))
	step := degStr(180.0 / float64(k))
	switch rand.Intn(4) {
	case 0:
		return trigPickMC([]trigMCClaim{{
			q:       fmt.Sprintf("Where does the graph of y = csc(%dx) have its vertical asymptotes?", k),
			correct: fmt.Sprintf("x = %sj", step),
			distractors: distinctDistractors(fmt.Sprintf("x = %sj", step),
				fmt.Sprintf("x = %s + %sj", step, step),
				fmt.Sprintf("x = %s + 180j", degStr(90.0/float64(k))),
				"nowhere",
			),
			reason: fmt.Sprintf("The cosecant is undefined where its sine is zero, and sin(%dx) = 0 at x = %sj.", k, step),
		}})
	case 1:
		return trigPickMC([]trigMCClaim{{
			q:       fmt.Sprintf("What is the range of y = %d sec(x)?", a),
			correct: fmt.Sprintf("(-infinity, -%d] together with [%d, infinity)", a, a),
			distractors: distinctDistractors(fmt.Sprintf("(-infinity, -%d] together with [%d, infinity)", a, a),
				fmt.Sprintf("[-%d, %d]", a, a),
				fmt.Sprintf("[-1, 1]"),
				fmt.Sprintf("(-infinity, -1] together with [1, infinity)"),
			),
			reason: fmt.Sprintf("The secant never lies between -1 and 1, so scaling by %d widens the excluded gap to (-%d, %d).", a, a, a),
		}})
	case 2:
		return trigPickMC([]trigMCClaim{{
			q:       fmt.Sprintf("The period of the graph of y = sec(%dx) is how many degrees?", k),
			correct: period,
			distractors: distinctDistractors(period,
				step,
				degStr(360.0*float64(k)),
				"180",
			),
			reason: fmt.Sprintf("sec(%dx) = 1/cos(%dx), so the period is that of cos(%dx), which is 360 divided by %d, or %s degrees.", k, k, k, k, period),
		}})
	default:
		return trigPickMC([]trigMCClaim{{
			q:       "Which identity defines the cosecant in terms of a sine?",
			correct: "csc(x) = 1/sin(x)",
			distractors: distinctDistractors("csc(x) = 1/sin(x)",
				"csc(x) = -sin(x)",
				"csc(x) = sin(x)/cos(x)",
				"csc(x) = 1/cos(x)",
			),
			reason: "The cosecant is the reciprocal of the sine, the same way the secant is the reciprocal of the cosine.",
		}})
	}
}

type inverseGraphGen struct{}

func (g *inverseGraphGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := 2 + rand.Intn(5) // 2..6
	switch rand.Intn(4) {
	case 0:
		return trigPickMC([]trigMCClaim{{
			q:       fmt.Sprintf("What is the domain of y = arcsin(x/%d)?", a),
			correct: fmt.Sprintf("[-%d, %d]", a, a),
			distractors: distinctDistractors(fmt.Sprintf("[-%d, %d]", a, a),
				fmt.Sprintf("[-1, 1]"),
				fmt.Sprintf("[0, %d]", a),
				fmt.Sprintf("(-infinity, infinity)"),
			),
			reason: fmt.Sprintf("The sine takes values only in [-1, 1], so x/%d must lie there, which puts x in [-%d, %d].", a, a, a),
		}})
	case 1:
		return trigPickMC([]trigMCClaim{{
			q:       fmt.Sprintf("What is the domain of y = arccos(%dx)?", a),
			correct: fmt.Sprintf("[-(1/%d), 1/%d]", a, a),
			distractors: distinctDistractors(fmt.Sprintf("[-(1/%d), 1/%d]", a, a),
				fmt.Sprintf("[-%d, %d]", a, a),
				fmt.Sprintf("[0, 1/%d]", a),
				fmt.Sprintf("(-infinity, infinity)"),
			),
			reason: fmt.Sprintf("The cosine takes values only in [-1, 1], so %dx must lie there, which puts x in [-(1/%d), 1/%d].", a, a, a),
		}})
	case 2:
		return trigPickMC([]trigMCClaim{{
			q:       "At which angles does the tangent have vertical asymptotes in its graph?",
			correct: "at the odd multiples of 90 degrees",
			distractors: distinctDistractors("at the odd multiples of 90 degrees",
				"at the multiples of 180 degrees",
				"at the even multiples of 90 degrees",
				"nowhere",
			),
			reason: "tan(x) = sin(x)/cos(x) and cos(x) is zero at 90, 270, 450 degrees and their negatives, which are the odd multiples of 90.",
		}})
	default:
		return trigPickMC([]trigMCClaim{{
			q:       "What is the range of y = arccos(x)?",
			correct: "from 0 to 180 degrees",
			distractors: distinctDistractors("from 0 to 180 degrees",
				"from -90 to 90 degrees",
				"from -180 to 180 degrees",
				"from 0 to 360 degrees",
			),
			reason: "The arccosine is defined to return the angle in [0, 180] degrees, which is exactly where the cosine is one-to-one.",
		}})
	}
}
