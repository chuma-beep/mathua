package odes

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

// odeBridgeOrdinal returns the section the learner is working in, or a random
// section when the engine asks for an unguided practice question.
func odeBridgeOrdinal(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return 1 + rand.Intn(3)
}

// odeDecimal keeps computed answers readable: two decimal places at most, and
// no trailing zeros or trailing decimal point.
func odeDecimal(v float64) string {
	s := strconv.FormatFloat(v, 'f', 2, 64)
	s = strings.TrimRight(s, "0")
	return strings.TrimRight(s, ".")
}

func odeBridgeProblem(ordinal int, question, answer, explanation string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = answer
	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts:       facts,
		KPOrdinal:   ordinal,
	}
}

type slopeFieldGen struct{}

func (g *slopeFieldGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := odeBridgeOrdinal(ctx)
	switch o {
	case 1:
		a := 1 + rand.Intn(4)
		b := -3 + rand.Intn(7)
		x0 := -3 + rand.Intn(7)
		y0 := -3 + rand.Intn(7)
		slope := a*x0 + b*y0
		return odeBridgeProblem(o,
			fmt.Sprintf("For y' = %d x + (%d) y, what is the slope at the point (%d, %d)? (enter a number)", a, b, x0, y0),
			strconv.Itoa(slope),
			fmt.Sprintf("Substitute x = %d and y = %d: %d(%d) + (%d)(%d) = %d.", x0, y0, a, x0, b, y0, slope),
			map[string]string{"kind": "slope", "x": strconv.Itoa(x0), "y": strconv.Itoa(y0)})
	case 2:
		switch rand.Intn(3) {
		case 0:
			k := 10 * (2 + rand.Intn(8))
			a := 1 + rand.Intn(4)
			return odeBridgeProblem(o,
				fmt.Sprintf("For y' = %d y (1 - y/%d), what is the nonzero equilibrium? (enter a number)", a, k),
				strconv.Itoa(k),
				fmt.Sprintf("The nonzero equilibrium is y = %d, where the factor 1 - y/%d is zero.", k, k),
				map[string]string{"kind": "equilibrium", "K": strconv.Itoa(k)})
		case 1:
			p := -5 + rand.Intn(5)
			q := p + 2 + rand.Intn(5)
			a := 1 + rand.Intn(3)
			return odeBridgeProblem(o,
				fmt.Sprintf("For y' = %d (y - (%d))(%d - y), what is the larger equilibrium? (enter a number)", a, p, q),
				strconv.Itoa(q),
				fmt.Sprintf("The equilibria are y = %d and y = %d, so the larger one is %d.", p, q, q),
				map[string]string{"kind": "equilibrium", "p": strconv.Itoa(p), "q": strconv.Itoa(q)})
		default:
			p := -5 + rand.Intn(5)
			q := p + 2 + rand.Intn(5)
			a := 1 + rand.Intn(3)
			return odeBridgeProblem(o,
				fmt.Sprintf("For y' = %d (y - (%d))(y - (%d)), what is the smaller equilibrium? (enter a number)", a, p, q),
				strconv.Itoa(p),
				fmt.Sprintf("The equilibria are y = %d and y = %d, so the smaller one is %d.", p, q, p),
				map[string]string{"kind": "equilibrium", "p": strconv.Itoa(p), "q": strconv.Itoa(q)})
		}
	default:
		a := 1 + rand.Intn(4)
		b := -3 + rand.Intn(7)
		x0 := -2 + rand.Intn(5)
		y0 := -2 + rand.Intn(5)
		slope := a*x0 + b*y0
		second := a + b*slope
		return odeBridgeProblem(o,
			fmt.Sprintf("For y' = %d x + (%d) y, what is y'' at (%d, %d)? (enter a number) (hint: y'' = f_x + f_y f)", a, b, x0, y0),
			strconv.Itoa(second),
			fmt.Sprintf("Here f_x = %d and f_y = %d, while f(%d, %d) = %d. So y'' = %d + (%d)(%d) = %d.",
				a, b, x0, y0, slope, a, b, slope, second),
			map[string]string{"kind": "concavity", "x": strconv.Itoa(x0), "y": strconv.Itoa(y0)})
	}
}

type logisticGen struct{}

func (g *logisticGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := odeBridgeOrdinal(ctx)
	k := []int{20, 25, 40, 50, 100}[rand.Intn(5)]
	r := []float64{0.1, 0.2, 0.25, 0.5}[rand.Intn(4)]
	switch o {
	case 1:
		return odeBridgeProblem(o,
			fmt.Sprintf("For dP/dt = %.2f P (1 - P/%d), what is the carrying capacity? (enter a number)", r, k),
			strconv.Itoa(k),
			fmt.Sprintf("The carrying capacity is K = %d, the value where the growth factor 1 - P/K becomes zero.", k),
			map[string]string{"kind": "capacity", "K": strconv.Itoa(k)})
	case 2:
		p := (k / 10) * (1 + rand.Intn(9))
		growth := r * float64(p) * (1 - float64(p)/float64(k))
		return odeBridgeProblem(o,
			fmt.Sprintf("For dP/dt = %.2f P (1 - P/%d), what is the growth rate when P = %d? (enter a decimal)", r, k, p),
			odeDecimal(growth),
			fmt.Sprintf("Substitute P = %d: %.2f(%d)(1 - %d/%d) = %s.", p, r, p, p, k, odeDecimal(growth)),
			map[string]string{"kind": "growth", "P": strconv.Itoa(p), "K": strconv.Itoa(k)})
	default:
		switch rand.Intn(3) {
		case 0:
			return odeBridgeProblem(o,
				fmt.Sprintf("For dP/dt = %.2f P (1 - P/%d), at what population is the growth rate largest? (enter a number)", r, k),
				strconv.Itoa(k/2),
				fmt.Sprintf("The growth rate is a downward parabola in P with maximum at P = K/2 = %d.", k/2),
				map[string]string{"kind": "max_growth_population", "K": strconv.Itoa(k)})
		case 1:
			rate := r * float64(k) / 4
			return odeBridgeProblem(o,
				fmt.Sprintf("For dP/dt = %.2f P (1 - P/%d), what is the maximum growth rate? (enter a decimal)", r, k),
				odeDecimal(rate),
				fmt.Sprintf("The maximum growth rate is rK/4 = %.2f x %d / 4 = %s.", r, k, odeDecimal(rate)),
				map[string]string{"kind": "max_growth_rate", "K": strconv.Itoa(k)})
		default:
			return odeBridgeProblem(o,
				fmt.Sprintf("For dP/dt = %.2f P (1 - P/%d), what is the long-run population for a positive initial population? (enter a number)", r, k),
				strconv.Itoa(k),
				fmt.Sprintf("Every positive solution increases or decreases toward the carrying capacity K = %d.", k),
				map[string]string{"kind": "long_run", "K": strconv.Itoa(k)})
		}
	}
}

type existenceUniquenessGen struct{}

func (g *existenceUniquenessGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := odeBridgeOrdinal(ctx)
	a := rand.Intn(7) - 3
	switch o {
	case 1:
		return odeBridgeProblem(o,
			fmt.Sprintf("For y' = sqrt(y - (%d)), what is the smallest integer y at which f(x,y) is continuous? (enter a number)", a),
			strconv.Itoa(a),
			fmt.Sprintf("The square root sqrt(y - (%d)) is defined and continuous when y >= %d, so the smallest integer is %d.", a, a, a),
			map[string]string{"kind": "f_continuity", "a": strconv.Itoa(a)})
	case 2:
		return odeBridgeProblem(o,
			fmt.Sprintf("For y' = sqrt(y - (%d)), at what integer y is the partial derivative with respect to y undefined? (enter a number)", a),
			strconv.Itoa(a),
			fmt.Sprintf("The partial derivative is 1/(2 sqrt(y - (%d))), which is undefined at y = %d.", a, a),
			map[string]string{"kind": "fy_continuity", "a": strconv.Itoa(a)})
	default:
		start := -3 + rand.Intn(9)
		values := make([]int, 5)
		for i := range values {
			values[i] = start + i
		}
		count := 0
		var ok []string
		for _, y := range values {
			if y > a {
				count++
				ok = append(ok, strconv.Itoa(y))
			}
		}
		return odeBridgeProblem(o,
			fmt.Sprintf("For y' = sqrt(y - (%d)), both f and its y-partial are continuous at how many of y = %d, %d, %d, %d, %d? (enter a number)",
				a, values[0], values[1], values[2], values[3], values[4]),
			strconv.Itoa(count),
			fmt.Sprintf("Both functions are continuous exactly when y > %d. That holds for %d of the five values: %s.",
				a, count, strings.Join(ok, ", ")),
			map[string]string{"kind": "count", "a": strconv.Itoa(a), "count": strconv.Itoa(count)})
	}
}
