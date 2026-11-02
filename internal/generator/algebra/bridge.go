package algebra

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

// algBridgeGen serves the Algebra I and Algebra II concepts added in the
// 2026 content expansion. Each knowledge point gets its own numeric shape so
// the question matches the section of the lesson the learner is reading.
type algBridgeGen struct{ kind string }

func algOrdinal(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return 1 + rand.Intn(3)
}

func algProb(o int, q, a, e string) generator.Problem {
	return generator.Problem{
		Question:    q,
		Answer:      a,
		Explanation: e,
		Facts:       map[string]string{"answer": a},
		KPOrdinal:   o,
	}
}

func (g *algBridgeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := algOrdinal(ctx)
	switch g.kind {
	case "rate_of_change":
		switch o {
		case 1:
			m := 2 + rand.Intn(8)
			a := rand.Intn(7) - 3
			b := a + 1 + rand.Intn(5)
			ans := m * (b - a)
			return algProb(o, fmt.Sprintf("For f(x) = %dx, what is the change in output f(b) - f(a) on the interval from x = %d to x = %d? (enter a number)", m, a, b), strconv.Itoa(ans), fmt.Sprintf("The outputs are %d and %d, so the change is %d.", m*a, m*b, ans))
		case 2:
			m := 2 + rand.Intn(8)
			c := rand.Intn(9) - 4
			a := rand.Intn(7) - 3
			b := a + 1 + rand.Intn(5)
			return algProb(o, fmt.Sprintf("Find the average rate of change of f(x) = %dx + %d on the interval from x = %d to x = %d. (enter a number)", m, c, a, b), strconv.Itoa(m), fmt.Sprintf("A linear function has the same average rate on every interval, namely its slope %d.", m))
		default:
			pairs := [][2]int{{120, 3}, {150, 5}, {240, 4}, {180, 6}, {200, 8}, {300, 5}, {90, 2}}
			p := pairs[rand.Intn(len(pairs))]
			return algProb(o, fmt.Sprintf("A car travels %d kilometres in %d hours. What is its average speed in kilometres per hour? (enter a number)", p[0], p[1]), strconv.Itoa(p[0]/p[1]), fmt.Sprintf("Distance divided by time is %d / %d = %d kilometres per hour.", p[0], p[1], p[0]/p[1]))
		}
	case "poly_graph_zeros":
		switch o {
		case 1:
			r1 := rand.Intn(9) - 4
			r2 := r1 + 1 + rand.Intn(4)
			r3 := r2 + 1 + rand.Intn(4)
			return algProb(o, fmt.Sprintf("A polynomial has zeros at x = %d, x = %d and x = %d. What is the largest zero? (enter a number)", r1, r2, r3), strconv.Itoa(r3), fmt.Sprintf("The largest of the three zeros is %d.", r3))
		case 2:
			r := rand.Intn(9) - 4
			k := 1 + rand.Intn(4)
			cross := 0
			if k%2 == 1 {
				cross = 1
			}
			note := "even, so the graph touches the axis and turns back"
			if k%2 == 1 {
				note = "odd, so the graph crosses the axis"
			}
			return algProb(o, fmt.Sprintf("A polynomial has the factor (x - %d)^%d. At x = %d, enter 1 if the graph crosses the axis and 0 if it touches and turns back. (enter a number)", r, k, r), strconv.Itoa(cross), fmt.Sprintf("The exponent %d is %s.", k, note))
		default:
			r1 := rand.Intn(7) - 3
			r2 := r1 + 1 + rand.Intn(4)
			m1 := 1 + rand.Intn(3)
			m2 := 1 + rand.Intn(3)
			deg := m1 + m2
			return algProb(o, fmt.Sprintf("A polynomial is written as (x - %d)^%d (x - %d)^%d. What is its degree? (enter a number)", r1, m1, r2, m2), strconv.Itoa(deg), fmt.Sprintf("Add the exponents: %d + %d = %d.", m1, m2, deg))
		}
	case "quadratic_form":
		switch o {
		case 1:
			a := 2 + rand.Intn(4)
			b := rand.Intn(9) - 4
			c := rand.Intn(9) - 4
			k := rand.Intn(9) - 4
			if k == c {
				k++
			}
			return algProb(o, fmt.Sprintf("Write %dx^2 + %dx = %d in standard form. What is the constant term c? (enter a number)", a, b, k), strconv.Itoa(-k), fmt.Sprintf("Subtract %d from both sides, so c = %d.", k, -k))
		case 2:
			a := 1 + rand.Intn(4)
			b := rand.Intn(11) - 5
			c := rand.Intn(7) - 3
			ans := b*b - 4*a*c
			return algProb(o, fmt.Sprintf("For %dx^2 + %dx + %d = 0, what is the discriminant b^2 - 4ac? (enter a number)", a, b, c), strconv.Itoa(ans), fmt.Sprintf("Substitute: (%d)^2 - 4(%d)(%d) = %d.", b, a, c, ans))
		default:
			a := 1 + rand.Intn(4)
			b := rand.Intn(11) - 5
			c := rand.Intn(7) - 3
			d := b*b - 4*a*c
			n := 2
			if d == 0 {
				n = 1
			} else if d < 0 {
				n = 0
			}
			return algProb(o, fmt.Sprintf("How many distinct real solutions does %dx^2 + %dx + %d = 0 have? (enter a number)", a, b, c), strconv.Itoa(n), fmt.Sprintf("The discriminant is %d, so there are %d distinct real solutions.", d, n))
		}
	case "linear_programming":
		switch o {
		case 1:
			s := 6 + rand.Intn(10)
			d := rand.Intn(4) - 2
			if (s+d)%2 != 0 {
				d++
			}
			x := (s + d) / 2
			return algProb(o, fmt.Sprintf("Two boundary lines are x + y = %d and x - y = %d. What is the x-coordinate of their intersection? (enter a number)", s, d), strconv.Itoa(x), fmt.Sprintf("Add the equations: 2x = %d, so x = %d.", s+d, x))
		case 2:
			x := 1 + rand.Intn(6)
			y := 1 + rand.Intn(6)
			a := 2 + rand.Intn(5)
			b := 2 + rand.Intn(5)
			ans := a*x + b*y
			return algProb(o, fmt.Sprintf("A feasible region has a corner at (%d, %d). What is the value of P = %dx + %dy there? (enter a number)", x, y, a, b), strconv.Itoa(ans), fmt.Sprintf("Substitute: %d(%d) + %d(%d) = %d.", a, x, b, y, ans))
		default:
			n := 4 + rand.Intn(8)
			a := 2 + rand.Intn(5)
			b := 2 + rand.Intn(5)
			best := a
			if b > a {
				best = b
			}
			ans := best * n
			return algProb(o, fmt.Sprintf("Maximise P = %dx + %dy subject to x >= 0, y >= 0 and x + y <= %d. What is the maximum value? (enter a number)", a, b, n), strconv.Itoa(ans), fmt.Sprintf("The maximum is reached at the corner on the axis with the larger coefficient, giving %d x %d = %d.", best, n, ans))
		}
	case "natural_base":
		switch o {
		case 1:
			k := 1 + rand.Intn(9)
			return algProb(o, fmt.Sprintf("What is ln(e^%d)? (enter a number)", k), strconv.Itoa(k), fmt.Sprintf("The natural logarithm undoes the natural exponential, so ln(e^%d) = %d.", k, k))
		case 2:
			k := 2 + rand.Intn(3)
			m := k * (2 + rand.Intn(4))
			x := m / k
			return algProb(o, fmt.Sprintf("Solve e^(%dx) = e^%d for x. (enter a number)", k, m), strconv.Itoa(x), fmt.Sprintf("Equal bases give equal exponents: %dx = %d, so x = %d.", k, m, x))
		default:
			k := 2 + rand.Intn(8)
			return algProb(o, fmt.Sprintf("What is e^(ln %d)? (enter a number)", k), strconv.Itoa(k), fmt.Sprintf("The natural exponential and the natural logarithm undo each other, so e^(ln %d) = %d.", k, k))
		}
	case "point_slope":
		switch o {
		case 1:
			m := 2 + rand.Intn(7)
			x1 := rand.Intn(7) - 3
			y1 := rand.Intn(9) - 4
			return algProb(o, fmt.Sprintf("A line is written as y - %d = %d(x - %d). What is its slope? (enter a number)", y1, m, x1), strconv.Itoa(m), fmt.Sprintf("Point-slope form is y - y1 = m(x - x1), so the slope is the coefficient %d.", m))
		case 2:
			m := 1 + rand.Intn(5)
			x1 := rand.Intn(5)
			y1 := rand.Intn(7) - 2
			x2 := x1 + 1 + rand.Intn(4)
			ans := y1 + m*(x2-x1)
			return algProb(o, fmt.Sprintf("A line passes through (%d, %d) with slope %d. What is its y-value at x = %d? (enter a number)", x1, y1, m, x2), strconv.Itoa(ans), fmt.Sprintf("Move %d steps in x and add %d each step: %d + %d x %d = %d.", x2-x1, m, y1, m, x2-x1, ans))
		default:
			m := 2 + rand.Intn(5)
			x1 := rand.Intn(5)
			y1 := rand.Intn(9) - 3
			b := y1 - m*x1
			return algProb(o, fmt.Sprintf("The line y - %d = %d(x - %d) is rewritten as y = mx + b. What is b? (enter a number)", y1, m, x1), strconv.Itoa(b), fmt.Sprintf("Distribute and add: b = %d - %d x %d = %d.", y1, m, x1, b))
		}
	case "seq_recursive":
		switch o {
		case 1:
			a1 := 2 + rand.Intn(8)
			k := 2 + rand.Intn(6)
			ans := a1 + 4*k
			return algProb(o, fmt.Sprintf("A sequence has a1 = %d and a(n+1) = a(n) + %d. What is a5? (enter a number)", a1, k), strconv.Itoa(ans), fmt.Sprintf("Four steps of size %d from %d give %d.", k, a1, ans))
		case 2:
			a1 := 1 + rand.Intn(4)
			r := 2 + rand.Intn(3)
			ans := a1 * r * r * r
			return algProb(o, fmt.Sprintf("A sequence has a1 = %d and a(n+1) = %d a(n). What is a4? (enter a number)", a1, r), strconv.Itoa(ans), fmt.Sprintf("Three multiplications by %d give %d x %d^3 = %d.", r, a1, r, ans))
		default:
			a1 := 1 + rand.Intn(4)
			a2 := a1 + 1 + rand.Intn(4)
			a3 := a1 + a2
			a4 := a2 + a3
			a5 := a3 + a4
			return algProb(o, fmt.Sprintf("A sequence has a1 = %d, a2 = %d and a(n+2) = a(n+1) + a(n). What is a5? (enter a number)", a1, a2), strconv.Itoa(a5), fmt.Sprintf("The terms are %d, %d, %d, %d and %d.", a1, a2, a3, a4, a5))
		}
	case "nth_root":
		switch o {
		case 1:
			n := 2 + rand.Intn(4)
			k := 2 + rand.Intn(5)
			v := 1
			for i := 0; i < n; i++ {
				v *= k
			}
			return algProb(o, fmt.Sprintf("What is the %dth root of %d? (enter a number)", n, v), strconv.Itoa(k), fmt.Sprintf("%d^%d = %d, so the root is %d.", k, n, v, k))
		case 2:
			n := 3 + rand.Intn(3)
			k := 2 + rand.Intn(3)
			v := 1
			for i := 0; i < n; i++ {
				v *= k
			}
			m := 1 + rand.Intn(n-1)
			pow := 1
			for i := 0; i < m; i++ {
				pow *= k
			}
			return algProb(o, fmt.Sprintf("Rewrite the expression %d^(%d/%d) as a whole number. (enter a number)", v, m, n), strconv.Itoa(pow), fmt.Sprintf("The denominator is the root and the numerator is the power: the %dth root of %d is %d, and %d^%d = %d.", n, v, k, k, m, pow))
		default:
			n := 3 + 2*rand.Intn(2)
			k := 1 + rand.Intn(4)
			v := 1
			for i := 0; i < n; i++ {
				v *= k
			}
			return algProb(o, fmt.Sprintf("What is the %dth root of -%d? (enter a number)", n, v), strconv.Itoa(-k), fmt.Sprintf("An odd root of a negative number is negative, and %d^%d = %d, so the root is -%d.", k, n, v, k))
		}
	case "fundamental_theorem":
		switch o {
		case 1:
			n := 3 + rand.Intn(8)
			return algProb(o, fmt.Sprintf("How many complex roots does a polynomial of degree %d have, counted with multiplicity? (enter a number)", n), strconv.Itoa(n), fmt.Sprintf("The fundamental theorem gives exactly %d roots over the complex numbers.", n))
		case 2:
			n := 3 + rand.Intn(7)
			k := 1 + rand.Intn(n-1)
			ans := n - k
			return algProb(o, fmt.Sprintf("A real polynomial of degree %d has exactly %d real roots. How many nonreal roots does it have, counted with multiplicity? (enter a number)", n, k), strconv.Itoa(ans), fmt.Sprintf("The remaining %d - %d = %d roots are nonreal.", n, k, ans))
		default:
			a := 1 + rand.Intn(4)
			b := 2 + rand.Intn(5)
			return algProb(o, fmt.Sprintf("A real polynomial has the root %d + %di. What is the imaginary part of the other nonreal root? (enter a number)", a, b), strconv.Itoa(-b), fmt.Sprintf("Nonreal roots of a real polynomial come in conjugate pairs, so the other root is %d - %di.", a, b))
		}
	case "half_life":
		switch o {
		case 1:
			a0 := 2 * (10 + rand.Intn(40))
			return algProb(o, fmt.Sprintf("A sample starts at %d grams and has a half-life of 1 hour. How many grams remain after 1 hour? (enter a number)", a0), strconv.Itoa(a0/2), fmt.Sprintf("One half-life leaves half: %d / 2 = %d.", a0, a0/2))
		case 2:
			a0 := 64 * (1 + rand.Intn(4))
			k := 1 + rand.Intn(4)
			ans := a0
			for i := 0; i < k; i++ {
				ans /= 2
			}
			return algProb(o, fmt.Sprintf("A sample starts at %d grams. How many grams remain after %d half-lives? (enter a number)", a0, k), strconv.Itoa(ans), fmt.Sprintf("Halve %d a total of %d times to get %d.", a0, k, ans))
		default:
			h := 1 + rand.Intn(6)
			p := 1 + rand.Intn(4)
			ans := h * p
			return algProb(o, fmt.Sprintf("A substance has a half-life of %d hours. How many hours until 1/%d of the original amount remains? (enter a number)", h, 1<<p), strconv.Itoa(ans), fmt.Sprintf("1/%d needs %d half-lives, so the time is %d x %d = %d hours.", 1<<p, p, p, h, ans))
		}
	case "rational_exponent":
		switch o {
		case 1:
			n := 2 + rand.Intn(3)
			k := 2 + rand.Intn(3)
			x := 1
			for i := 0; i < n; i++ {
				x *= k
			}
			return algProb(o, fmt.Sprintf("Solve x^(1/%d) = %d for x. (enter a number)", n, k), strconv.Itoa(x), fmt.Sprintf("Raise both sides to the %d power: x = %d^%d = %d.", n, k, n, x))
		case 2:
			n := 2 + rand.Intn(3)
			x := 2 + rand.Intn(4)
			m := 1 + rand.Intn(n-1)
			k := 1
			for i := 0; i < m; i++ {
				k *= x
			}
			return algProb(o, fmt.Sprintf("Solve the equation x^(%d/%d) = %d for x. (enter a number)", m, n, k), strconv.Itoa(x), fmt.Sprintf("Raise both sides to the %d/%d power: x = %d.", n, m, x))
		default:
			n := 2 + rand.Intn(2)
			k := 2 + rand.Intn(5)
			x := 1
			for i := 0; i < n; i++ {
				x *= k
			}
			return algProb(o, fmt.Sprintf("Solve x^(-1/%d) = 1/%d for x. (enter a number)", n, k), strconv.Itoa(x), fmt.Sprintf("A negative exponent is a reciprocal, so x^(1/%d) = %d and x = %d.", n, k, x))
		}
	case "rational_word":
		switch o {
		case 1:
			pairs := [][2]int{{2, 2}, {3, 6}, {4, 4}, {6, 3}, {5, 20}, {8, 8}, {10, 15}}
			p := pairs[rand.Intn(len(pairs))]
			ans := p[0] * p[1] / (p[0] + p[1])
			return algProb(o, fmt.Sprintf("One machine finishes a job in %d hours and a second in %d hours. Working together, how many hours do they take? (enter a number)", p[0], p[1]), strconv.Itoa(ans), fmt.Sprintf("The combined rate is 1/%d + 1/%d, so the time is %d x %d / (%d + %d) = %d hours.", p[0], p[1], p[0], p[1], p[0], p[1], ans))
		case 2:
			a := 2 + rand.Intn(5)
			b := 2 + rand.Intn(5)
			c := 1 + rand.Intn(b-1)
			x := a * b / c
			for x*c != a*b {
				b++
				x = a * b / c
			}
			return algProb(o, fmt.Sprintf("The proportion x/%d = %d/%d is true. What is x? (enter a number)", a, b, c), strconv.Itoa(x), fmt.Sprintf("Cross multiply: %d x = %d x %d, so x = %d.", c, a, b, x))
		default:
			pairs := [][2]int{{6, 3}, {12, 4}, {20, 5}, {8, 4}, {10, 5}, {15, 5}}
			p := pairs[rand.Intn(len(pairs))]
			ans := p[0] * p[1] / (p[0] - p[1])
			return algProb(o, fmt.Sprintf("A pipe fills a tank in %d hours and a drain empties it in %d hours. With both open, how many hours does the tank take to fill? (enter a number)", p[0], p[1]), strconv.Itoa(ans), fmt.Sprintf("The net rate is 1/%d - 1/%d, so the time is %d x %d / (%d - %d) = %d hours.", p[0], p[1], p[0], p[1], p[0], p[1], ans))
		}
	case "composite_domain":
		switch o {
		case 1:
			a := rand.Intn(7) - 2
			return algProb(o, fmt.Sprintf("Let f(x) = sqrt(x) and g(x) = x - %d. What is the smallest integer in the domain of f(g(x))? (enter a number)", a), strconv.Itoa(a), fmt.Sprintf("The composite needs x - %d >= 0, so the smallest integer is %d.", a, a))
		case 2:
			a := rand.Intn(6)
			return algProb(o, fmt.Sprintf("Let f(x) = sqrt(x) and g(x) = 1/(x - %d). What is the smallest integer in the domain of f(g(x))? (enter a number)", a), strconv.Itoa(a+1), fmt.Sprintf("The inner function needs x > %d and its output is positive there, so the smallest integer is %d.", a, a+1))
		default:
			a := rand.Intn(6) - 2
			return algProb(o, fmt.Sprintf("Let f(x) = ln(x) and g(x) = x - %d. What is the smallest integer in the domain of f(g(x))? (enter a number)", a), strconv.Itoa(a+1), fmt.Sprintf("The logarithm needs x - %d > 0, so the smallest integer is %d.", a, a+1))
		}
	case "linear_word":
		switch o {
		case 1:
			r1 := 2 + rand.Intn(5)
			r2 := r1 + 1 + rand.Intn(4)
			f1 := 10 * (1 + rand.Intn(6))
			f2 := f1 - (r2-r1)*(1+rand.Intn(5))
			if f2 < 0 {
				f2 += 20
			}
			x := (f1 - f2) / (r2 - r1)
			return algProb(o, fmt.Sprintf("Plan A costs %d plus %d per month. Plan B costs %d plus %d per month. After how many months do they cost the same? (enter a number)", f1, r1, f2, r2), strconv.Itoa(x), fmt.Sprintf("Set %d + %dx = %d + %dx and solve: x = %d.", f1, r1, f2, r2, x))
		case 2:
			base := 10 * (1 + rand.Intn(8))
			rate := 2 + rand.Intn(9)
			q := 2 + rand.Intn(9)
			ans := base + rate*q
			return algProb(o, fmt.Sprintf("A repair charges %d dollars plus %d dollars per hour. What is the cost for %d hours? (enter a number)", base, rate, q), strconv.Itoa(ans), fmt.Sprintf("Substitute: %d + %d x %d = %d.", base, rate, q, ans))
		default:
			k := 2 + rand.Intn(3)
			b := 3 + rand.Intn(9)
			s := (k + 1) * b
			return algProb(o, fmt.Sprintf("An older sibling is %d times as old as a younger one. Their ages add to %d. How old is the younger sibling? (enter a number)", k, s), strconv.Itoa(b), fmt.Sprintf("If the younger age is x, then %dx + x = %d, so x = %d.", k, s, b))
		}
	case "horizontal_vertical":
		switch o {
		case 1:
			x := rand.Intn(9) - 4
			y := rand.Intn(11) - 5
			return algProb(o, fmt.Sprintf("A horizontal line passes through (%d, %d). What is its y-value at every point? (enter a number)", x, y), strconv.Itoa(y), fmt.Sprintf("A horizontal line keeps y constant, so the equation is y = %d.", y))
		case 2:
			x := rand.Intn(9) - 4
			y := rand.Intn(11) - 5
			return algProb(o, fmt.Sprintf("A vertical line passes through (%d, %d). What is its x-value at every point? (enter a number)", x, y), strconv.Itoa(x), fmt.Sprintf("A vertical line keeps x constant, so the equation is x = %d.", x))
		default:
			y := rand.Intn(9) - 4
			x1 := rand.Intn(5)
			x2 := x1 + 1 + rand.Intn(4)
			return algProb(o, fmt.Sprintf("The points (%d, %d) and (%d, %d) lie on one line. Enter 0 if the line is horizontal and 1 if it is vertical. (enter a number)", x1, y, x2, y), "0", fmt.Sprintf("Both points have y = %d, so the line is horizontal.", y))
		}
	}
	return algProb(o, "Evaluate 2 + 2. (enter a number)", "4", "Two plus two is four.")
}
