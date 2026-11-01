package trigonometry

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

// A multiple-choice question asks the learner to recognise one value among
// several. A count question asks for every listed statement to be judged, so
// the same exact values become work rather than a choice. The pools below hold
// sixteen statements, eight true and eight false, and the answer is drawn
// uniformly from zero to eight before the statements are chosen. Drawing the
// answer first keeps the rare all-true and all-false cases as frequent as the
// middle ones, which a random sample of eight cannot do.

type trigVerdict struct {
	label string
	yes   bool
	note  string
}

const trigDraw = 8

func countTrigVerdicts(lead string, pool []trigVerdict) generator.Problem {
	var yes, no []trigVerdict
	for _, v := range pool {
		if v.yes {
			yes = append(yes, v)
		} else {
			no = append(no, v)
		}
	}
	count := rand.Intn(trigDraw + 1)
	if count > len(yes) {
		count = len(yes)
	}
	if trigDraw-count > len(no) {
		count = trigDraw - len(no)
	}
	var chosen []trigVerdict
	var question string
	for attempt := 0; attempt < 40; attempt++ {
		yesCopy := append([]trigVerdict(nil), yes...)
		noCopy := append([]trigVerdict(nil), no...)
		rand.Shuffle(len(yesCopy), func(i, j int) { yesCopy[i], yesCopy[j] = yesCopy[j], yesCopy[i] })
		rand.Shuffle(len(noCopy), func(i, j int) { noCopy[i], noCopy[j] = noCopy[j], noCopy[i] })
		chosen = append(append([]trigVerdict(nil), yesCopy[:count]...), noCopy[:trigDraw-count]...)
		rand.Shuffle(len(chosen), func(i, j int) { chosen[i], chosen[j] = chosen[j], chosen[i] })
		parts := make([]string, len(chosen))
		for i, v := range chosen {
			parts[i] = fmt.Sprintf("(%c) %s", 'a'+i, v.label)
		}
		question = fmt.Sprintf("%s %s (enter a number)", lead, strings.Join(parts, " "))
		if !strings.Contains(question, fmt.Sprintf("%d", count)) {
			break
		}
	}
	notes := make([]string, len(chosen))
	for i, v := range chosen {
		verdict := "no"
		if v.yes {
			verdict = "yes"
		}
		notes[i] = fmt.Sprintf("%s: %s, %s", v.label, verdict, v.note)
	}
	return generator.Problem{
		Question:    question,
		Answer:      fmt.Sprintf("%d", count),
		Explanation: fmt.Sprintf("Exactly %d of the %d statements are true.\n%s", count, trigDraw, strings.Join(notes, "\n")),
		Facts:       map[string]string{"count": fmt.Sprintf("%d", count)},
	}
}

func angleConversionProblem() generator.Problem {
	pool := []trigVerdict{
		{"30 degrees = pi/6", true, "180 degrees is pi, so 30 degrees is one sixth of pi."},
		{"45 degrees = pi/4", true, "45 degrees is one eighth of 360 degrees."},
		{"60 degrees = pi/3", true, "60 degrees is one sixth of 360 degrees."},
		{"90 degrees = pi/2", true, "90 degrees is a quarter turn."},
		{"120 degrees = 2pi/3", true, "120 degrees is one third of 360 degrees."},
		{"180 degrees = pi", true, "A half turn is pi radians."},
		{"270 degrees = 3pi/2", true, "270 degrees is three quarters of a turn."},
		{"360 degrees = 2pi", true, "A full turn is 2pi radians."},
		{"30 degrees = pi/3", false, "30 degrees is pi/6; pi/3 is 60 degrees."},
		{"45 degrees = pi/2", false, "45 degrees is pi/4; pi/2 is 90 degrees."},
		{"60 degrees = pi/6", false, "60 degrees is pi/3; pi/6 is 30 degrees."},
		{"90 degrees = pi", false, "90 degrees is pi/2; pi is 180 degrees."},
		{"120 degrees = 3pi/2", false, "120 degrees is 2pi/3; 3pi/2 is 270 degrees."},
		{"135 degrees = 2pi/3", false, "135 degrees is 3pi/4; 2pi/3 is 120 degrees."},
		{"210 degrees = 5pi/6", false, "210 degrees is 7pi/6; 5pi/6 is 150 degrees."},
		{"300 degrees = 4pi/3", false, "300 degrees is 5pi/3; 4pi/3 is 240 degrees."},
	}
	return countTrigVerdicts("How many of these angle conversions are correct?", pool)
}

func referenceAngleProblem() generator.Problem {
	pool := []trigVerdict{
		{"the reference angle for 150 degrees is 30 degrees", true, "150 is 30 short of 180."},
		{"the reference angle for 210 degrees is 30 degrees", true, "210 is 30 past 180."},
		{"the reference angle for 300 degrees is 60 degrees", true, "300 is 60 short of 360."},
		{"the reference angle for 240 degrees is 60 degrees", true, "240 is 60 past 180."},
		{"the reference angle for 135 degrees is 45 degrees", true, "135 is 45 short of 180."},
		{"the reference angle for 225 degrees is 45 degrees", true, "225 is 45 past 180."},
		{"the reference angle for 120 degrees is 60 degrees", true, "120 is 60 short of 180."},
		{"the reference angle for 330 degrees is 30 degrees", true, "330 is 30 short of 360."},
		{"the reference angle for 150 degrees is 60 degrees", false, "The acute angle to the x-axis is 30 degrees."},
		{"the reference angle for 210 degrees is 60 degrees", false, "The acute angle to the x-axis is 30 degrees."},
		{"the reference angle for 300 degrees is 30 degrees", false, "The acute angle to the x-axis is 60 degrees."},
		{"the reference angle for 240 degrees is 30 degrees", false, "The acute angle to the x-axis is 60 degrees."},
		{"the reference angle for 135 degrees is 30 degrees", false, "The acute angle to the x-axis is 45 degrees."},
		{"the reference angle for 225 degrees is 60 degrees", false, "The acute angle to the x-axis is 45 degrees."},
		{"the reference angle for 120 degrees is 30 degrees", false, "The acute angle to the x-axis is 60 degrees."},
		{"the reference angle for 330 degrees is 60 degrees", false, "The acute angle to the x-axis is 30 degrees."},
	}
	return countTrigVerdicts("How many of these reference angles are correct?", pool)
}

func sineGraphProblem() generator.Problem {
	pool := []trigVerdict{
		{"y = sin x has period 360 degrees", true, "The sine function repeats every full turn."},
		{"y = sin x has amplitude 1", true, "It oscillates one unit either side of the axis."},
		{"y = sin x has range [-1, 1]", true, "The maximum is 1 and the minimum is -1."},
		{"y = sin x is an odd function", true, "sin(-x) = -sin x."},
		{"y = sin x passes through the origin", true, "sin 0 = 0."},
		{"y = sin x has zeros at multiples of 180 degrees", true, "sin k pi = 0 for every integer k."},
		{"y = sin x reaches its maximum at 90 degrees", true, "sin 90 degrees = 1."},
		{"y = sin x has a minimum of -1 at 270 degrees", true, "sin 270 degrees = -1."},
		{"y = sin x has period 180 degrees", false, "That is the period of tangent."},
		{"y = sin x has amplitude 2", false, "The amplitude is 1."},
		{"y = sin x has range [0, 1]", false, "It also takes negative values."},
		{"y = sin x is an even function", false, "Cosine is the even one."},
		{"y = sin x has zeros at multiples of 90 degrees", false, "Only the multiples of 180 degrees are zeros."},
		{"y = sin x reaches its maximum at 0 degrees", false, "The maximum is at 90 degrees."},
		{"y = sin x has no minimum", false, "The minimum is -1."},
		{"y = sin x has period pi radians", false, "The period is 2pi radians."},
	}
	return countTrigVerdicts("How many of these statements about the sine graph are true?", pool)
}

func cosineGraphProblem() generator.Problem {
	pool := []trigVerdict{
		{"y = cos x has period 360 degrees", true, "The cosine function repeats every full turn."},
		{"y = cos x has amplitude 1", true, "It oscillates one unit either side of the axis."},
		{"y = cos x has range [-1, 1]", true, "The maximum is 1 and the minimum is -1."},
		{"y = cos x is an even function", true, "cos(-x) = cos x."},
		{"y = cos x has value 1 at 0 degrees", true, "cos 0 = 1."},
		{"y = cos x has zeros at 90 degrees plus multiples of 180 degrees", true, "cos 90 degrees = 0."},
		{"y = cos x reaches its maximum at 0 degrees", true, "The value there is 1."},
		{"y = cos x has a minimum of -1 at 180 degrees", true, "cos 180 degrees = -1."},
		{"y = cos x has period 180 degrees", false, "That is the period of tangent."},
		{"y = cos x has amplitude 2", false, "The amplitude is 1."},
		{"y = cos x has range [0, 1]", false, "It also takes negative values."},
		{"y = cos x is an odd function", false, "Sine is the odd one."},
		{"y = cos x has value 0 at 0 degrees", false, "cos 0 = 1."},
		{"y = cos x has zeros at multiples of 180 degrees", false, "The zeros are 90 degrees plus multiples of 180 degrees."},
		{"y = cos x reaches its maximum at 90 degrees", false, "The maximum is at 0 degrees."},
		{"y = cos x has period pi radians", false, "The period is 2pi radians."},
	}
	return countTrigVerdicts("How many of these statements about the cosine graph are true?", pool)
}

func tangentGraphProblem() generator.Problem {
	pool := []trigVerdict{
		{"y = tan x has period 180 degrees", true, "Tangent repeats every half turn."},
		{"y = tan x has range all real numbers", true, "It takes every real value."},
		{"y = tan x is an odd function", true, "tan(-x) = -tan x."},
		{"y = tan x has asymptotes at 90 degrees plus multiples of 180 degrees", true, "Cosine is zero there."},
		{"y = tan x has value 0 at 0 degrees", true, "tan 0 = 0."},
		{"y = tan x is increasing on (-90 degrees, 90 degrees)", true, "Its derivative is sec squared, which is positive."},
		{"y = tan x has no amplitude", true, "It is unbounded, so no amplitude is defined."},
		{"y = tan x has zeros at multiples of 180 degrees", true, "tan k pi = 0."},
		{"y = tan x has period 360 degrees", false, "The period is 180 degrees."},
		{"y = tan x has range [-1, 1]", false, "It is unbounded."},
		{"y = tan x is an even function", false, "Tangent is odd."},
		{"y = tan x has asymptotes at multiples of 180 degrees", false, "The asymptotes are at 90 degrees plus multiples of 180 degrees."},
		{"y = tan x has value 1 at 0 degrees", false, "tan 0 = 0."},
		{"y = tan x is decreasing on (-90 degrees, 90 degrees)", false, "It is increasing on each branch."},
		{"y = tan x has amplitude 1", false, "It is unbounded, so no amplitude is defined."},
		{"y = tan x has zeros at 90 degrees plus multiples of 180 degrees", false, "Those are its asymptotes."},
	}
	return countTrigVerdicts("How many of these statements about the tangent graph are true?", pool)
}

func graphTransformationProblem() generator.Problem {
	pool := []trigVerdict{
		{"y = 2 sin x has amplitude 2", true, "The coefficient of sine is the amplitude."},
		{"y = 3 cos x has amplitude 3", true, "The coefficient of cosine is the amplitude."},
		{"y = sin(2x) has period 180 degrees", true, "The period is 360 divided by 2."},
		{"y = cos(3x) has period 120 degrees", true, "The period is 360 divided by 3."},
		{"y = sin x + 4 has vertical shift 4", true, "The constant added outside the function shifts it up."},
		{"y = sin(x - 30 degrees) has phase shift 30 degrees right", true, "Subtracting inside moves the graph right."},
		{"y = -sin x is a reflection of y = sin x in the x-axis", true, "The negative sign flips every value."},
		{"y = sin x + 2 has range [1, 3]", true, "Adding 2 moves [-1, 1] up to [1, 3]."},
		{"y = 2 sin x has amplitude 1", false, "The amplitude is 2."},
		{"y = sin(2x) has period 360 degrees", false, "The period is 180 degrees."},
		{"y = cos(3x) has period 360 degrees", false, "The period is 120 degrees."},
		{"y = sin x + 4 has vertical shift -4", false, "The shift is 4 up."},
		{"y = sin(x - 30 degrees) has phase shift 30 degrees left", false, "It shifts 30 degrees right."},
		{"y = sin x + 2 has range [-1, 1]", false, "Adding 2 changes the range to [1, 3]."},
		{"y = -sin x has amplitude -1", false, "Amplitude is positive; it is 1."},
		{"y = 3 cos x has amplitude 1/3", false, "The amplitude is 3."},
	}
	return countTrigVerdicts("How many of these statements about transformed sine and cosine graphs are true?", pool)
}

func inverseFunctionProblem() generator.Problem {
	pool := []trigVerdict{
		{"arcsin has domain [-1, 1]", true, "Sine only takes values between -1 and 1."},
		{"arccos has range [0, pi]", true, "The principal cosine inverse returns angles from 0 to pi."},
		{"arctan has range (-pi/2, pi/2)", true, "The principal tangent inverse stays between the two vertical asymptotes."},
		{"arcsin(sin x) = x only on [-pi/2, pi/2]", true, "Outside that interval the inverse returns the principal angle instead."},
		{"arctan is defined for every real number", true, "Tangent covers every real value on (-pi/2, pi/2)."},
		{"arccos is decreasing on [-1, 1]", true, "Larger cosine values correspond to smaller angles."},
		{"arcsin is an odd function", true, "arcsin(-x) = -arcsin(x)."},
		{"arctan has horizontal asymptotes y = pi/2 and y = -pi/2", true, "The graph approaches those levels but never reaches them."},
		{"arcsin has range [0, pi]", false, "Its range is [-pi/2, pi/2]."},
		{"arccos has domain [-pi/2, pi/2]", false, "Its domain is [-1, 1]."},
		{"arctan has range all real numbers", false, "Its range is (-pi/2, pi/2)."},
		{"arcsin is defined for every real number", false, "Its domain is [-1, 1]."},
		{"arccos is an odd function", false, "arccos(-x) = pi - arccos(x), so it is neither odd nor even."},
		{"every inverse trigonometric function has period 2pi", false, "Principal inverse functions are not periodic."},
		{"arcsin(sin x) = x for every real x", false, "It fails outside [-pi/2, pi/2]."},
		{"arctan has a vertical asymptote at x = pi/2", false, "It has horizontal asymptotes instead."},
	}
	return countTrigVerdicts("How many of these statements about inverse trigonometric functions are true?", pool)
}

func arcsineProblem() generator.Problem {
	pool := []trigVerdict{
		{"y = arcsin x has domain [-1, 1]", true, "Only those inputs are sine values."},
		{"y = arcsin x has range [-pi/2, pi/2]", true, "This is the principal branch."},
		{"y = arcsin x is increasing", true, "Sine increases on the principal interval."},
		{"y = arcsin x is odd", true, "arcsin(-x) = -arcsin(x)."},
		{"arcsin 0 = 0", true, "sin 0 = 0."},
		{"arcsin 1 = pi/2", true, "sin(pi/2) = 1."},
		{"arcsin(-1) = -pi/2", true, "sin(-pi/2) = -1."},
		{"the graph of arcsin passes through the origin", true, "arcsin 0 = 0."},
		{"y = arcsin x has range [0, pi]", false, "That is the range of arccos."},
		{"y = arcsin x has domain all real numbers", false, "Its domain is [-1, 1]."},
		{"y = arcsin x is decreasing", false, "It is increasing."},
		{"y = arcsin x is even", false, "It is odd."},
		{"arcsin 1 = 0", false, "arcsin 1 = pi/2."},
		{"arcsin 0 = pi/2", false, "arcsin 0 = 0."},
		{"arcsin(-1) = pi/2", false, "arcsin(-1) = -pi/2."},
		{"y = arcsin x has period 2pi", false, "The principal branch is not periodic."},
	}
	return countTrigVerdicts("How many of these statements about y = arcsin x are true?", pool)
}

func arccosineProblem() generator.Problem {
	pool := []trigVerdict{
		{"y = arccos x has domain [-1, 1]", true, "Only those inputs are cosine values."},
		{"y = arccos x has range [0, pi]", true, "This is the principal branch."},
		{"y = arccos x is decreasing", true, "Larger cosine values give smaller angles."},
		{"arccos 1 = 0", true, "cos 0 = 1."},
		{"arccos 0 = pi/2", true, "cos(pi/2) = 0."},
		{"arccos(-1) = pi", true, "cos pi = -1."},
		{"arccos(-x) = pi - arccos x", true, "This identity follows from the unit circle."},
		{"the y-intercept of y = arccos x is pi/2", true, "arccos 0 = pi/2."},
		{"y = arccos x has range [-pi/2, pi/2]", false, "That is the range of arcsin."},
		{"y = arccos x is increasing", false, "It is decreasing."},
		{"y = arccos x is odd", false, "It is neither odd nor even."},
		{"arccos 0 = 0", false, "arccos 0 = pi/2."},
		{"arccos 1 = pi/2", false, "arccos 1 = 0."},
		{"arccos(-1) = 0", false, "arccos(-1) = pi."},
		{"y = arccos x has domain all real numbers", false, "Its domain is [-1, 1]."},
		{"y = arccos x has period 2pi", false, "The principal branch is not periodic."},
	}
	return countTrigVerdicts("How many of these statements about y = arccos x are true?", pool)
}

func arctangentProblem() generator.Problem {
	pool := []trigVerdict{
		{"y = arctan x has domain all real numbers", true, "Every real number is a tangent value."},
		{"y = arctan x has range (-pi/2, pi/2)", true, "This is the principal branch."},
		{"y = arctan x is increasing", true, "Tangent increases on the principal interval."},
		{"y = arctan x is odd", true, "arctan(-x) = -arctan(x)."},
		{"arctan 0 = 0", true, "tan 0 = 0."},
		{"arctan 1 = pi/4", true, "tan(pi/4) = 1."},
		{"y = arctan x has horizontal asymptotes y = pi/2 and y = -pi/2", true, "The graph approaches those levels."},
		{"y = arctan x passes through the origin", true, "arctan 0 = 0."},
		{"y = arctan x has range [0, pi]", false, "That is the range of arccos."},
		{"y = arctan x is decreasing", false, "It is increasing."},
		{"y = arctan x is even", false, "It is odd."},
		{"y = arctan x has period pi", false, "The principal branch is not periodic."},
		{"y = arctan x has vertical asymptotes", false, "It has horizontal asymptotes."},
		{"arctan 1 = pi/2", false, "arctan 1 = pi/4."},
		{"arctan 0 = pi/2", false, "arctan 0 = 0."},
		{"y = arctan x has a maximum at x = 1", false, "It is increasing and has no maximum."},
	}
	return countTrigVerdicts("How many of these statements about y = arctan x are true?", pool)
}

func arccotangentProblem() generator.Problem {
	pool := []trigVerdict{
		{"y = arccot x has domain all real numbers", true, "Every real number is a cotangent value."},
		{"y = arccot x has range (0, pi)", true, "This is the standard principal branch."},
		{"y = arccot x is decreasing", true, "Cotangent decreases on (0, pi)."},
		{"arccot 0 = pi/2", true, "cot(pi/2) = 0."},
		{"arccot 1 = pi/4", true, "cot(pi/4) = 1."},
		{"y = arccot x has horizontal asymptotes y = 0 and y = pi", true, "The graph approaches those levels."},
		{"arccot(-x) = pi - arccot x", true, "This identity follows from the principal branch."},
		{"the graph of arccot never touches the x-axis", true, "It approaches 0 only as an asymptote."},
		{"y = arccot x has range (-pi/2, pi/2)", false, "That is the range of arctan."},
		{"y = arccot x is increasing", false, "It is decreasing."},
		{"y = arccot x is odd", false, "It is neither odd nor even."},
		{"arccot 0 = 0", false, "arccot 0 = pi/2."},
		{"arccot 1 = pi/2", false, "arccot 1 = pi/4."},
		{"y = arccot x has period pi", false, "The principal branch is not periodic."},
		{"y = arccot x has vertical asymptotes", false, "It has horizontal asymptotes."},
		{"y = arccot x passes through the origin", false, "It has y-intercept pi/2."},
	}
	return countTrigVerdicts("How many of these statements about y = arccot x are true?", pool)
}

func inverseGraphProblem() generator.Problem {
	pool := []trigVerdict{
		{"the graph of arcsin passes through the origin", true, "arcsin 0 = 0."},
		{"the graph of arccos runs from (-1, pi) to (1, 0)", true, "It decreases across its domain."},
		{"the graph of arctan has two horizontal asymptotes", true, "They are y = pi/2 and y = -pi/2."},
		{"the graph of arccot is decreasing", true, "Cotangent decreases on (0, pi)."},
		{"an inverse trigonometric graph is the reflection of the restricted original across y = x", true, "Inverse graphs swap x and y."},
		{"the domain of an inverse graph is the range of the restricted original", true, "Inputs and outputs are exchanged."},
		{"the graph of arctan crosses the origin", true, "arctan 0 = 0."},
		{"the y-intercept of arccos is pi/2", true, "arccos 0 = pi/2."},
		{"the graph of arcsin has vertical asymptotes", false, "It is defined on the closed interval [-1, 1]."},
		{"the graph of arccos is increasing", false, "It is decreasing."},
		{"the graph of arctan has range [0, pi]", false, "Its range is (-pi/2, pi/2)."},
		{"the graph of arccot passes through the origin", false, "Its y-intercept is pi/2."},
		{"the graph of arcsin has range [0, pi]", false, "Its range is [-pi/2, pi/2]."},
		{"the graph of arccos has a vertical asymptote at x = 0", false, "It is defined at x = 0."},
		{"the graph of arctan is periodic", false, "The principal branch is not periodic."},
		{"the graph of arccot is the same as the graph of arctan shifted up", false, "Their shapes and asymptotes differ."},
	}
	return countTrigVerdicts("How many of these statements about inverse trigonometric graphs are true?", pool)
}

func reciprocalGraphProblem() generator.Problem {
	pool := []trigVerdict{
		{"y = csc x has vertical asymptotes at multiples of pi", true, "Sine is zero at those points."},
		{"y = sec x has vertical asymptotes at pi/2 plus multiples of pi", true, "Cosine is zero there."},
		{"y = cot x has vertical asymptotes at multiples of pi", true, "Tangent is zero at those points."},
		{"the range of csc is (-infinity, -1] union [1, infinity)", true, "It is the reciprocal of sine."},
		{"sec is an even function", true, "sec(-x) = sec(x)."},
		{"csc is an odd function", true, "csc(-x) = -csc(x)."},
		{"cot is decreasing on each interval between its asymptotes", true, "Its derivative is -csc^2 x."},
		{"the period of cot is pi", true, "Cotangent repeats every pi."},
		{"the range of csc is [-1, 1]", false, "Cosecant never takes values between -1 and 1."},
		{"sec is an odd function", false, "It is even."},
		{"the period of cot is 2pi", false, "It is pi."},
		{"csc is defined at x = 0", false, "Sine is zero at 0, so cosecant is undefined."},
		{"sec has zeros at pi/2 plus multiples of pi", false, "Those are its asymptotes."},
		{"the range of cot is [-1, 1]", false, "Its range is all real numbers."},
		{"csc has amplitude 1", false, "Reciprocal graphs are unbounded, so they have no amplitude."},
		{"cot is an even function", false, "It is odd."},
	}
	return countTrigVerdicts("How many of these statements about reciprocal trigonometric graphs are true?", pool)
}

func hyperbolicDerivativeProblem() generator.Problem {
	pool := []trigVerdict{
		{"d/dx sinh x = cosh x", true, "The hyperbolic sine and cosine differentiate into each other."},
		{"d/dx cosh x = sinh x", true, "This is the standard derivative."},
		{"d/dx tanh x = sech^2 x", true, "The derivative is positive and peaks at zero."},
		{"d/dx sech x = -sech x tanh x", true, "This mirrors the derivative of secant."},
		{"d/dx coth x = -csch^2 x", true, "The derivative is negative everywhere it is defined."},
		{"d/dx csch x = -csch x coth x", true, "This mirrors the derivative of cosecant."},
		{"d/dx arcsinh x = 1/sqrt(x^2 + 1)", true, "The inverse hyperbolic sine has this derivative."},
		{"d/dx arctanh x = 1/(1 - x^2)", true, "This holds for |x| < 1."},
		{"d/dx sinh x = sech x", false, "The derivative is cosh x."},
		{"d/dx cosh x = -sinh x", false, "The derivative is sinh x."},
		{"d/dx tanh x = tanh^2 x", false, "The derivative is sech^2 x."},
		{"d/dx sech x = sech x tanh x", false, "The correct derivative has a minus sign."},
		{"d/dx coth x = csch^2 x", false, "The correct derivative has a minus sign."},
		{"d/dx arcsinh x = 1/sqrt(1 - x^2)", false, "That is the derivative of arcsin."},
		{"d/dx arctanh x = 1/(x^2 - 1)", false, "The denominator is 1 - x^2."},
		{"d/dx csch x = csch x coth x", false, "The correct derivative has a minus sign."},
	}
	return countTrigVerdicts("How many of these hyperbolic derivative statements are true?", pool)
}
