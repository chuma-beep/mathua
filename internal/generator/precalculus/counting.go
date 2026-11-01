package precalculus

// Count and parameter questions for the precalculus concepts that used to be
// recognition questions. The count form asks how many of eight listed
// statements are true, with the answer drawn uniformly from zero to eight so
// every value appears. The parameter forms ask for one number that follows
// from the stated model, so the learner computes rather than recognises.

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

type precalcVerdict struct {
	label string
	yes   bool
	note  string
}

const precalcDraw = 8

func countPrecalcVerdicts(lead string, pool []precalcVerdict) generator.Problem {
	var trueItems, falseItems []precalcVerdict
	for _, v := range pool {
		if v.yes {
			trueItems = append(trueItems, v)
		} else {
			falseItems = append(falseItems, v)
		}
	}

	var chosen []precalcVerdict
	var question string
	count := 0
	for attempt := 0; attempt < 40; attempt++ {
		count = rand.Intn(precalcDraw + 1)
		rand.Shuffle(len(trueItems), func(i, j int) {
			trueItems[i], trueItems[j] = trueItems[j], trueItems[i]
		})
		rand.Shuffle(len(falseItems), func(i, j int) {
			falseItems[i], falseItems[j] = falseItems[j], falseItems[i]
		})
		chosen = append(chosen[:0], trueItems[:count]...)
		chosen = append(chosen, falseItems[:precalcDraw-count]...)
		rand.Shuffle(len(chosen), func(i, j int) {
			chosen[i], chosen[j] = chosen[j], chosen[i]
		})

		var b strings.Builder
		b.WriteString(lead)
		for i, v := range chosen {
			fmt.Fprintf(&b, "\n\n(%c) %s", 'a'+i, v.label)
		}
		b.WriteString("\n\n(enter a number)")
		question = b.String()
		if !strings.Contains(question, strconv.Itoa(count)) {
			break
		}
	}

	var explanation strings.Builder
	fmt.Fprintf(&explanation, "Exactly %d of the %d statements are true.", count, precalcDraw)
	for _, v := range chosen {
		verdict := "false"
		if v.yes {
			verdict = "true"
		}
		fmt.Fprintf(&explanation, "\n- %s: %s. %s", v.label, verdict, v.note)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(count),
		Explanation: explanation.String(),
		Facts:       map[string]string{"count": strconv.Itoa(count)},
	}
}

func polarCircleProblem() generator.Problem {
	pool := []precalcVerdict{
		{`\(r = a\)`, true, "This is the circle centered at the pole with radius a."},
		{`\(r = -a\)`, true, "Negative radius traces the same circle as r = a."},
		{`\(r = a\cos\theta\)`, true, "This is a circle of diameter a through the pole, centered on the polar axis."},
		{`\(r = a\sin\theta\)`, true, "This is a circle of diameter a through the pole, centered on the line theta = 90 degrees."},
		{`\(r = 2a\cos\theta\)`, true, "This is a circle of diameter 2a through the pole, centered on the polar axis."},
		{`\(r = 2a\sin\theta\)`, true, "This is a circle of diameter 2a through the pole, centered on the line theta = 90 degrees."},
		{`\(r = a\cos(\theta - \alpha)\)`, true, "This is a circle through the pole rotated by the fixed angle alpha."},
		{`\(r = a\sin(\theta + \alpha)\)`, true, "This is a circle through the pole rotated by the fixed angle alpha."},
		{`\(\theta = \alpha\)`, false, "A fixed angle is a straight line through the pole."},
		{`\(r = a\theta\)`, false, "This is an Archimedean spiral."},
		{`\(r = a(1 + \cos\theta)\)`, false, "This is a cardioid, not a circle."},
		{`\(r = a\cos 2\theta\)`, false, "This is a four-petal rose."},
		{`\(r^2 = a^2\cos 2\theta\)`, false, "This is a lemniscate."},
		{`\(r = a\sec\theta\)`, false, "This is the vertical line x = a."},
		{`\(r = a\csc\theta\)`, false, "This is the horizontal line y = a."},
		{`\(r = e^{a\theta}\)`, false, "This is a logarithmic spiral."},
	}
	return countPrecalcVerdicts(
		`In these equations \(a > 0\) and \(\alpha\) is a fixed angle. How many represent a circle?`,
		pool)
}

func parametricParabolaProblem() generator.Problem {
	pool := []precalcVerdict{
		{`\(x = t,\ y = t^2\)`, true, "Eliminating t gives y = x^2, a parabola opening upward."},
		{`\(x = t^2,\ y = t\)`, true, "Eliminating t gives x = y^2, a parabola opening rightward."},
		{`\(x = t + h,\ y = (t + h)^2\)`, true, "With u = t + h this is y = u^2, a parabola opening upward."},
		{`\(x = t^2 + h,\ y = t\)`, true, "With u = t and x = u^2 + h this is a parabola opening rightward."},
		{`\(x = kt,\ y = k^2 t^2\)`, true, "Eliminating t gives y = x^2, a parabola opening upward."},
		{`\(x = k t^2,\ y = k t\)`, true, "Eliminating t gives x = y^2/k, a parabola opening rightward."},
		{`\(x = (t + h)^2,\ y = t + h\)`, true, "With u = t + h this is x = u^2, a parabola opening rightward."},
		{`\(x = t^2,\ y = kt\)`, true, "Eliminating t gives y^2 = k^2 x, a parabola opening rightward."},
		{`\(x = t,\ y = kt\)`, false, "Eliminating t gives y = kx, a straight line."},
		{`\(x = a\cos t,\ y = a\sin t\)`, false, "This is the circle x^2 + y^2 = a^2."},
		{`\(x = a\cos t,\ y = b\sin t\)`, false, "This is an ellipse, not a parabola."},
		{`\(x = a\sec t,\ y = b\tan t\)`, false, "This is a hyperbola."},
		{`\(x = t,\ y = t^n\)`, false, "With n greater than two this is a higher-degree curve, not a parabola."},
		{`\(x = t,\ y = e^t\)`, false, "This is an exponential curve."},
		{`\(x = t\cos t,\ y = t\sin t\)`, false, "This is a spiral."},
		{`\(x = a,\ y = t\)`, false, "This is the vertical line x = a."},
	}
	return countPrecalcVerdicts(
		`In these parameterizations \(h\), \(k\), \(a\) and \(b\) are nonzero constants and \(n\) is an integer greater than two. How many represent a parabola?`,
		pool)
}

func transformationPointProblem(ctx generator.GeneratorContext) generator.Problem {
	limit := 2 + int(ctx.Difficulty*4)
	h := 1 + rand.Intn(limit)
	k := 1 + rand.Intn(limit)
	a := 2 + rand.Intn(3)
	p := -4 + rand.Intn(9)
	q := -4 + rand.Intn(9)

	var form, coordinate, reason string
	answer := 0
	switch rand.Intn(6) {
	case 0:
		form = fmt.Sprintf(`\(g(x) = f(x - %d)\)`, h)
		coordinate = "x"
		answer = p + h
		reason = fmt.Sprintf("The graph shifts right by %d, so the x-coordinate becomes %d + %d = %d.", h, p, h, answer)
	case 1:
		form = fmt.Sprintf(`\(g(x) = f(x + %d)\)`, h)
		coordinate = "x"
		answer = p - h
		reason = fmt.Sprintf("The graph shifts left by %d, so the x-coordinate becomes %d - %d = %d.", h, p, h, answer)
	case 2:
		form = fmt.Sprintf(`\(g(x) = f(x) + %d\)`, k)
		coordinate = "y"
		answer = q + k
		reason = fmt.Sprintf("The graph shifts up by %d, so the y-coordinate becomes %d + %d = %d.", k, q, k, answer)
	case 3:
		form = fmt.Sprintf(`\(g(x) = f(x) - %d\)`, k)
		coordinate = "y"
		answer = q - k
		reason = fmt.Sprintf("The graph shifts down by %d, so the y-coordinate becomes %d - %d = %d.", k, q, k, answer)
	case 4:
		form = fmt.Sprintf(`\(g(x) = %d f(x)\)`, a)
		coordinate = "y"
		answer = a * q
		reason = fmt.Sprintf("A vertical stretch by %d multiplies the y-coordinate: %d times %d is %d.", a, a, q, answer)
	default:
		form = `\(g(x) = f(-x)\)`
		coordinate = "x"
		answer = -p
		reason = fmt.Sprintf("Reflecting across the y-axis sends the x-coordinate %d to %d.", p, answer)
	}

	return generator.Problem{
		Question: fmt.Sprintf(
			`The point \((%d, %d)\) lies on the graph of \(f\). The graph of %s is a transformation of the graph of \(f\). What is the %s-coordinate of the corresponding point on the graph of \(g\)? (enter a number)`,
			p, q, form, coordinate),
		Answer:      strconv.Itoa(answer),
		Explanation: reason,
		Facts: map[string]string{
			"x": strconv.Itoa(p), "y": strconv.Itoa(q), "answer": strconv.Itoa(answer),
		},
	}
}

func precalcDecimal(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func continuousCompoundingProblem(ctx generator.GeneratorContext) generator.Problem {
	principal := 500 * (1 + rand.Intn(8))
	rate := float64(1+rand.Intn(9)) / 100
	years := 1 + rand.Intn(10)
	rateText := precalcDecimal(rate)

	switch rand.Intn(3) {
	case 0:
		exponent := rate * float64(years)
		return generator.Problem{
			Question: fmt.Sprintf(
				`An account grows according to \(A(t) = P e^{rt}\). With \(r = %s\) and \(t = %d\), what is the exponent \(rt\)? (enter a decimal)`,
				rateText, years),
			Answer: precalcDecimal(exponent),
			Explanation: fmt.Sprintf(
				"Multiply the rate by the time: %s times %d equals %s.", rateText, years, precalcDecimal(exponent)),
			Facts: map[string]string{
				"principal": strconv.Itoa(principal), "rate": rateText, "years": strconv.Itoa(years),
			},
		}
	case 1:
		return generator.Problem{
			Question: fmt.Sprintf(
				`An account grows according to \(A(t) = P e^{rt}\) with principal \(P = %d\). What is \(A(0)\)? (enter a number)`,
				principal),
			Answer: strconv.Itoa(principal),
			Explanation: fmt.Sprintf(
				"At t = 0 the exponential factor is 1, so A(0) equals the principal, %d.", principal),
			Facts: map[string]string{
				"principal": strconv.Itoa(principal), "rate": rateText, "years": strconv.Itoa(years),
			},
		}
	default:
		return generator.Problem{
			Question: fmt.Sprintf(
				`An investment is modeled by \(A(t) = P e^{rt}\), and the exponent is \(%st\). What is the annual rate \(r\) as a decimal? (enter a decimal)`,
				rateText),
			Answer: rateText,
			Explanation: fmt.Sprintf(
				"The exponent has the form rt, so the coefficient of t is the rate, %s.", rateText),
			Facts: map[string]string{
				"principal": strconv.Itoa(principal), "rate": rateText, "years": strconv.Itoa(years),
			},
		}
	}
}

func coolingTemperatureProblem(ctx generator.GeneratorContext) generator.Problem {
	ambient := 10 + 5*rand.Intn(5)
	excess := 40 + 10*rand.Intn(6)
	initial := ambient + excess
	k := float64(1+rand.Intn(5)) / 10
	formula := fmt.Sprintf(`\(T(t) = %d + %d e^{-%st}\)`, ambient, excess, precalcDecimal(k))

	var ask, reason string
	answer := 0
	switch rand.Intn(5) {
	case 0:
		ask = "What is the ambient temperature?"
		answer = ambient
		reason = fmt.Sprintf("The constant term is the ambient temperature, %d.", ambient)
	case 1:
		ask = `What is the temperature at \(t = 0\)?`
		answer = initial
		reason = fmt.Sprintf("At t = 0 the exponential is 1, so T(0) = %d + %d = %d.", ambient, excess, initial)
	case 2:
		ask = "What is the initial temperature excess over the ambient temperature?"
		answer = excess
		reason = fmt.Sprintf("The coefficient of the exponential is the initial excess, %d.", excess)
	case 3:
		ask = "What temperature does the object approach after a long time?"
		answer = ambient
		reason = fmt.Sprintf("The exponential term decays to zero, leaving the ambient temperature %d.", ambient)
	default:
		ask = "What is the coefficient multiplying the exponential?"
		answer = excess
		reason = fmt.Sprintf("The coefficient in front of the exponential is %d.", excess)
	}

	return generator.Problem{
		Question: fmt.Sprintf(
			`A body cools according to %s, where t is in minutes. %s (enter a number)`,
			formula, ask),
		Answer:      strconv.Itoa(answer),
		Explanation: reason,
		Facts: map[string]string{
			"ambient": strconv.Itoa(ambient), "initial": strconv.Itoa(initial),
			"excess": strconv.Itoa(excess), "k": precalcDecimal(k),
		},
	}
}
