package calculus

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

// registerBridge registers the calculus concepts added by the content
// expansion. Each generator uses KPOrdinal to match the lesson section.
func registerBridge(reg *generator.Registry) {
	reg.Register("calc.integral.riemann_sums", &riemannSumsGen{})
	reg.Register("calc.series.taylor_error", &taylorErrorGen{})
	reg.Register("calc.series.termwise", &termwiseSeriesGen{})
	reg.Register("calc.multi.functions_several", &functionsSeveralGen{})
	reg.Register("calc.multi.level_curves", &levelCurvesGen{})
	reg.Register("calc.vec.vector_function", &vectorFunctionGen{})
	for id, kind := range map[string]string{
		"calc.integral.properties":    "integral_properties",
		"calc.deriv.definition":       "deriv_definition",
		"calc.integral.initial_value": "initial_value",
		"calc.deriv.rate":             "deriv_rate",
		"calc.deriv.graph_reading":    "graph_reading",
		"calc.integral.symmetry":      "integral_symmetry",
		"calc.limit.trig":             "limit_trig",
	} {
		reg.Register(id, &calc1BridgeGen{kind: kind})
	}
}

func calcBridgeOrdinal(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return rand.Intn(3) + 1
}

func calcBridgeProblem(ordinal int, question, answer, explanation string, facts map[string]string) generator.Problem {
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

func calcBridgeDecimal(v float64) string {
	s := strconv.FormatFloat(v, 'f', 6, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	if s == "" || s == "-0" {
		return "0"
	}
	return s
}

func calcFactorial(n int) int {
	p := 1
	for i := 2; i <= n; i++ {
		p *= i
	}
	return p
}

type riemannSumsGen struct{}

func (g *riemannSumsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := calcBridgeOrdinal(ctx)
	switch o {
	case 1:
		n := 2 + rand.Intn(6)
		a := 1 + rand.Intn(3)
		if rand.Intn(2) == 0 {
			sum := a * n * (n + 1) / 2
			return calcBridgeProblem(o,
				fmt.Sprintf("Use %d unit subintervals on \\([0,%d]\\) and the right endpoint of each. For \\(f(x)=%dx\\), what is the Riemann sum? (enter a number)", n, n, a),
				strconv.Itoa(sum),
				fmt.Sprintf("The right endpoints are 1 through %d, so the sum is %d(1+...+%d) = %d.", n, a, n, sum),
				map[string]string{"n": strconv.Itoa(n), "kind": "right-linear"})
		}
		sum := a * n * (n + 1) * (2*n + 1) / 6
		return calcBridgeProblem(o,
			fmt.Sprintf("Use %d unit subintervals on \\([0,%d]\\) and the right endpoint of each. For \\(f(x)=%dx^{2}\\), what is the Riemann sum? (enter a number)", n, n, a),
			strconv.Itoa(sum),
			fmt.Sprintf("The right endpoints are 1 through %d, so the sum is %d(1^2+...+%d^2) = %d.", n, a, n, sum),
			map[string]string{"n": strconv.Itoa(n), "kind": "right-quadratic"})
	case 2:
		n := 2 + rand.Intn(6)
		a := 1 + rand.Intn(3)
		if rand.Intn(2) == 0 {
			sum := a * n * (n - 1) / 2
			return calcBridgeProblem(o,
				fmt.Sprintf("Use %d unit subintervals on \\([0,%d]\\) and the left endpoint of each. For \\(f(x)=%dx\\), what is the Riemann sum? (enter a number)", n, n, a),
				strconv.Itoa(sum),
				fmt.Sprintf("The left endpoints are 0 through %d, so the sum is %d(0+...+%d) = %d.", n-1, a, n-1, sum),
				map[string]string{"n": strconv.Itoa(n), "kind": "left-linear"})
		}
		sum := a * n * (n + 1) * (2*n + 1) / 6
		sum -= a * n * n
		return calcBridgeProblem(o,
			fmt.Sprintf("Use %d unit subintervals on \\([0,%d]\\) and the left endpoint of each. For \\(f(x)=%dx^{2}\\), what is the Riemann sum? (enter a number)", n, n, a),
			strconv.Itoa(sum),
			fmt.Sprintf("The left endpoints are 0 through %d, so the sum is %d(0^2+...+%d^2) = %d.", n-1, a, n-1, sum),
			map[string]string{"n": strconv.Itoa(n), "kind": "left-quadratic"})
	default:
		a := 1 + rand.Intn(3)
		b := 2 * (1 + rand.Intn(3))
		if rand.Intn(2) == 0 {
			ans := a * b * b / 2
			return calcBridgeProblem(o,
				fmt.Sprintf("As the number of subintervals tends to infinity, the right Riemann sums of \\(f(x)=%dx\\) on \\([0,%d]\\) approach a limit. What is the limit? (enter a number)", a, b),
				strconv.Itoa(ans),
				fmt.Sprintf("The limit is the definite integral: %d times the area of a triangle with base %d and height %d, so the limit is %d.", a, b, a*b, ans),
				map[string]string{"kind": "limit-linear"})
		}
		b = 3 * (1 + rand.Intn(2))
		ans := a * b * b * b / 3
		return calcBridgeProblem(o,
			fmt.Sprintf("As the number of subintervals tends to infinity, the right Riemann sums of \\(f(x)=%dx^{2}\\) on \\([0,%d]\\) approach a limit. What is the limit? (enter a number)", a, b),
			strconv.Itoa(ans),
			fmt.Sprintf("The limit is the definite integral of %dx^2 from 0 to %d, which equals %d.", a, b, ans),
			map[string]string{"kind": "limit-quadratic"})
	}
}

type taylorErrorGen struct{}

func (g *taylorErrorGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := calcBridgeOrdinal(ctx)
	switch o {
	case 1:
		b := []float64{0.5, 1, 1.5}[rand.Intn(3)]
		n := 1 + rand.Intn(3)
		bound := math.Exp(b) * math.Pow(b, float64(n+1)) / float64(calcFactorial(n+1))
		return calcBridgeProblem(o,
			fmt.Sprintf("For \\(f(x)=e^{x}\\) at \\(a=0\\), bound \\(|R_{%d}(x)|\\) on \\([0,%s]\\) using the Lagrange remainder. What is the bound? (enter a decimal)", n, calcBridgeDecimal(b)),
			calcBridgeDecimal(bound),
			fmt.Sprintf("The derivative is bounded by M = e^{%s}. The Lagrange bound is M|b|^{n+1}/(n+1)! with n=%d, which is %s.", calcBridgeDecimal(b), n, calcBridgeDecimal(bound)),
			map[string]string{"kind": "exp-bound", "n": strconv.Itoa(n)})
	case 2:
		b := []float64{0.5, 1, 2}[rand.Intn(3)]
		n := 1 + rand.Intn(4)
		bound := math.Pow(b, float64(n+1)) / float64(calcFactorial(n+1))
		return calcBridgeProblem(o,
			fmt.Sprintf("For \\(f(x)=\\sin x\\) at \\(a=0\\), bound \\(|R_{%d}(x)|\\) on \\([0,%s]\\) using the Lagrange remainder. What is the bound? (enter a decimal)", n, calcBridgeDecimal(b)),
			calcBridgeDecimal(bound),
			fmt.Sprintf("All derivatives of sine have magnitude at most 1, so the bound is |b|^{n+1}/(n+1)! with b=%s and n=%d, which is %s.", calcBridgeDecimal(b), n, calcBridgeDecimal(bound)),
			map[string]string{"kind": "sine-bound", "n": strconv.Itoa(n)})
	default:
		tols := []float64{0.05, 0.01, 0.005, 0.001, 0.0005, 0.0001, 0.00001}
		tol := tols[rand.Intn(len(tols))]
		n := 0
		for math.E/float64(calcFactorial(n+1)) > tol {
			n++
		}
		return calcBridgeProblem(o,
			fmt.Sprintf("For \\(e^{x}\\) on \\([0,1]\\), what is the smallest Taylor degree \\(n\\) whose Lagrange bound is at most %s? (enter a number)", calcBridgeDecimal(tol)),
			strconv.Itoa(n),
			fmt.Sprintf("The bound is e/(n+1)!. Increasing n until e/(n+1)! <= %s gives n=%d.", calcBridgeDecimal(tol), n),
			map[string]string{"kind": "tolerance", "n": strconv.Itoa(n)})
	}
}

type termwiseSeriesGen struct{}

func (g *termwiseSeriesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := calcBridgeOrdinal(ctx)
	switch o {
	case 1:
		n := 2 + rand.Intn(5)
		c := 2 + rand.Intn(5)
		ans := n * c
		return calcBridgeProblem(o,
			fmt.Sprintf("A power series has coefficient %d on \\(x^{%d}\\). What is the coefficient of \\(x^{%d}\\) in its termwise derivative? (enter a number)", c, n, n-1),
			strconv.Itoa(ans),
			fmt.Sprintf("Termwise differentiation multiplies the coefficient of x^n by n, so the new coefficient is %d times %d = %d.", n, c, ans),
			map[string]string{"kind": "first-derivative", "n": strconv.Itoa(n)})
	case 2:
		n := 1 + rand.Intn(5)
		k := 1 + rand.Intn(5)
		a := k * (n + 1)
		return calcBridgeProblem(o,
			fmt.Sprintf("A power series has coefficient %d on \\(x^{%d}\\). What is the coefficient of \\(x^{%d}\\) in its termwise integral? (enter a number)", a, n, n+1),
			strconv.Itoa(k),
			fmt.Sprintf("Termwise integration divides the coefficient of x^n by n+1, so %d divided by %d is %d.", a, n+1, k),
			map[string]string{"kind": "integral", "n": strconv.Itoa(n)})
	default:
		n := 3 + rand.Intn(4)
		c := 2 + rand.Intn(4)
		ans := n * (n - 1) * c
		return calcBridgeProblem(o,
			fmt.Sprintf("A power series has coefficient %d on \\(x^{%d}\\). What is the coefficient of \\(x^{%d}\\) in its second termwise derivative? (enter a number)", c, n, n-2),
			strconv.Itoa(ans),
			fmt.Sprintf("Differentiating twice multiplies by n(n-1), so %d times %d times %d = %d.", n, n-1, c, ans),
			map[string]string{"kind": "second-derivative", "n": strconv.Itoa(n)})
	}
}

type functionsSeveralGen struct{}

func (g *functionsSeveralGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := calcBridgeOrdinal(ctx)
	switch o {
	case 1:
		a := 1 + rand.Intn(3)
		b := 1 + rand.Intn(4)
		c := rand.Intn(5)
		x := 1 + rand.Intn(3)
		y := 1 + rand.Intn(3)
		ans := a*x*x + b*y + c
		return calcBridgeProblem(o,
			fmt.Sprintf("For \\(f(x,y)=%dx^{2}+%dy+%d\\), what is \\(f(%d,%d)\\)? (enter a number)", a, b, c, x, y),
			strconv.Itoa(ans),
			fmt.Sprintf("Substitute x=%d and y=%d: %d(%d)^2+%d(%d)+%d = %d.", x, y, a, x, b, y, c, ans),
			map[string]string{"kind": "evaluate"})
	case 2:
		a := 3 + rand.Intn(6)
		if rand.Intn(2) == 0 {
			return calcBridgeProblem(o,
				fmt.Sprintf("For \\(f(x,y)=\\sqrt{%d-x-y}\\), what is the largest integer value of \\(x+y\\) for which f is defined? (enter a number)", a),
				strconv.Itoa(a),
				fmt.Sprintf("The square root requires %d-x-y >= 0, so x+y <= %d. The largest integer is %d.", a, a, a),
				map[string]string{"kind": "domain-sqrt"})
		}
		return calcBridgeProblem(o,
			fmt.Sprintf("For \\(f(x,y)=\\ln(%d-x-y)\\), what is the largest integer value of \\(x+y\\) for which f is defined? (enter a number)", a),
			strconv.Itoa(a-1),
			fmt.Sprintf("The logarithm requires %d-x-y > 0, so x+y < %d. The largest integer strictly below %d is %d.", a, a, a, a-1),
			map[string]string{"kind": "domain-log"})
	default:
		a := 1 + rand.Intn(3)
		b := 1 + rand.Intn(3)
		x := 1 + rand.Intn(4)
		y := 1 + rand.Intn(4)
		ans := a*x*x + b*y*y
		return calcBridgeProblem(o,
			fmt.Sprintf("For \\(f(x,y)=%dx^{2}+%dy^{2}\\), what is the value of the level curve through \\((%d,%d)\\)? (enter a number)", a, b, x, y),
			strconv.Itoa(ans),
			fmt.Sprintf("The level value is f(%d,%d) = %d(%d)^2+%d(%d)^2 = %d.", x, y, a, x, b, y, ans),
			map[string]string{"kind": "level-value"})
	}
}

type levelCurvesGen struct{}

func (g *levelCurvesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := calcBridgeOrdinal(ctx)
	switch o {
	case 1:
		r := 2 + rand.Intn(5)
		ans := r * r
		return calcBridgeProblem(o,
			fmt.Sprintf("For \\(f(x,y)=x^{2}+y^{2}\\), what level value \\(k\\) gives a circle of radius %d? (enter a number)", r),
			strconv.Itoa(ans),
			fmt.Sprintf("The level curve x^2+y^2=k has radius sqrt(k), so k = %d^2 = %d.", r, ans),
			map[string]string{"kind": "circle-radius"})
	case 2:
		k := 1 + rand.Intn(9)
		x := k + 2 + rand.Intn(5)
		y2 := x*x - k
		return calcBridgeProblem(o,
			fmt.Sprintf("For \\(f(x,y)=x^{2}-y^{2}\\), the point \\((%d,y)\\) lies on the level curve \\(f=%d\\). What is \\(y^{2}\\)? (enter a number)", x, k),
			strconv.Itoa(y2),
			fmt.Sprintf("Substitute x=%d and f=%d: %d - y^2 = %d, so y^2 = %d.", x, k, x*x, k, y2),
			map[string]string{"kind": "hyperbola"})
	default:
		k := 2 + rand.Intn(10)
		x := 1 + rand.Intn(6)
		ans := k - x
		return calcBridgeProblem(o,
			fmt.Sprintf("For \\(f(x,y)=x+y\\), the point \\((%d,y)\\) lies on the level curve \\(f=%d\\). What is \\(y\\)? (enter a number)", x, k),
			strconv.Itoa(ans),
			fmt.Sprintf("Substitute x=%d: %d+y=%d, so y=%d.", x, x, k, ans),
			map[string]string{"kind": "line"})
	}
}

type vectorFunctionGen struct{}

func (g *vectorFunctionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := calcBridgeOrdinal(ctx)
	switch o {
	case 1:
		a := 2 + rand.Intn(4)
		tChoice := rand.Intn(3)
		tName := []string{"0", "\\pi/2", "\\pi"}[tChoice]
		component := rand.Intn(2)
		var ans int
		if component == 0 {
			ans = []int{a, 0, -a}[tChoice]
		} else {
			ans = []int{0, a, 0}[tChoice]
		}
		name := "x"
		if component == 1 {
			name = "y"
		}
		return calcBridgeProblem(o,
			fmt.Sprintf("For \\(\\mathbf r(t)=\\langle %d\\cos t,%d\\sin t,4t\\rangle\\), what is the %s-component at \\(t=%s\\)? (enter a number)", a, a, name, tName),
			strconv.Itoa(ans),
			fmt.Sprintf("At t=%s, the %s-component is %d.", tName, name, ans),
			map[string]string{"kind": "evaluate"})
	case 2:
		a := 2 + rand.Intn(4)
		tChoice := rand.Intn(3)
		tName := []string{"0", "\\pi/2", "\\pi"}[tChoice]
		component := rand.Intn(3)
		var ans int
		switch component {
		case 0:
			ans = []int{0, -a, 0}[tChoice]
		case 1:
			ans = []int{a, 0, -a}[tChoice]
		default:
			ans = 4
		}
		name := []string{"x'", "y'", "z'"}[component]
		return calcBridgeProblem(o,
			fmt.Sprintf("For \\(\\mathbf r(t)=\\langle %d\\cos t,%d\\sin t,4t\\rangle\\), what is \\(%s(t)\\) at \\(t=%s\\)? (enter a number)", a, a, name, tName),
			strconv.Itoa(ans),
			fmt.Sprintf("Differentiating gives the stated component, and at t=%s its value is %d.", tName, ans),
			map[string]string{"kind": "derivative"})
	default:
		triples := [][2]int{{3, 4}, {6, 8}, {5, 12}, {8, 15}, {9, 12}, {7, 24}, {20, 21}, {9, 40}}
		pair := triples[rand.Intn(len(triples))]
		ans := int(math.Sqrt(float64(pair[0]*pair[0] + pair[1]*pair[1])))
		return calcBridgeProblem(o,
			fmt.Sprintf("For \\(\\mathbf r(t)=\\langle %d\\cos t,%d\\sin t,%dt\\rangle\\), what is the speed? (enter a number)", pair[0], pair[0], pair[1]),
			strconv.Itoa(ans),
			fmt.Sprintf("The velocity has magnitude sqrt(%d^2+%d^2) = %d.", pair[0], pair[1], ans),
			map[string]string{"kind": "speed"})
	}
}
