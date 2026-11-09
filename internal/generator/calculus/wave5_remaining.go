package calculus

// Wave 5 Calculus III, remaining twenty-one concepts. The finite and polynomial
// examples keep the theorem-heavy topics computational: each family asks for a
// coefficient, an area, a boundary value or a finite sum, and publishes the
// operands in Facts for independent recomputation.

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave5Remaining(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"calc.vec.general_stokes":         "general_stokes",
		"calc.vec.pullback":               "pullback",
		"calc.vec.closed_exact":           "closed_exact",
		"calc.vec.vector_potential":       "vector_potential",
		"calc.vec.helmholtz":              "helmholtz",
		"calc.vec.gauge":                  "gauge",
		"calc.multi.harmonic":             "harmonic",
		"calc.multi.mean_value":           "mean_value",
		"calc.multi.maximum_principle":    "maximum_principle",
		"calc.multi.poisson":              "poisson",
		"calc.multi.dirichlet":            "dirichlet",
		"calc.multi.neumann":              "neumann",
		"calc.vec.greens_identity":        "greens_identity",
		"calc.multi.fundamental_solution": "fundamental_solution",
		"calc.vec.euler_lagrange":         "euler_lagrange",
		"calc.vec.first_variation":        "first_variation",
		"calc.vec.geodesic":               "geodesic",
		"calc.vec.minimal_surface":        "minimal_surface",
		"calc.multi.surface_measure":      "surface_measure",
		"calc.multi.fourier":              "fourier",
		"calc.multi.convolution":          "convolution",
	} {
		reg.Register(id, &w5CalcRemainingGen{kind: kind})
	}
}

type w5CalcRemainingGen struct{ kind string }

func (g *w5CalcRemainingGen) Generate(ctx generator.GeneratorContext) (p generator.Problem) {
	defer func() {
		p.Question = w5Clean(p.Question)
		p.Explanation = w5Clean(p.Explanation)
	}()
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "general_stokes":
		return w5GeneralStokes(kp)
	case "pullback":
		return w5Pullback(kp)
	case "closed_exact":
		return w5ClosedExact(kp)
	case "vector_potential":
		return w5VectorPotential(kp)
	case "helmholtz":
		return w5Helmholtz(kp)
	case "gauge":
		return w5Gauge(kp)
	case "harmonic":
		return w5Harmonic(kp)
	case "mean_value":
		return w5MeanValue(kp)
	case "maximum_principle":
		return w5MaximumPrinciple(kp)
	case "poisson":
		return w5Poisson(kp)
	case "dirichlet":
		return w5Dirichlet(kp)
	case "neumann":
		return w5Neumann(kp)
	case "greens_identity":
		return w5GreensIdentity(kp)
	case "fundamental_solution":
		return w5FundamentalSolution(kp)
	case "euler_lagrange":
		return w5EulerLagrange(kp)
	case "first_variation":
		return w5FirstVariation(kp)
	case "geodesic":
		return w5Geodesic(kp)
	case "minimal_surface":
		return w5MinimalSurface(kp)
	case "surface_measure":
		return w5SurfaceMeasure(kp)
	case "fourier":
		return w5Fourier(kp)
	case "convolution":
		return w5Convolution(kp)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

var w5PureNumber = regexp.MustCompile(`\\\((-?\d+)\\\)`)
var w5UnitCoefficient = regexp.MustCompile(`(^|[^0-9])1 ([a-zA-Z])`)
var w5UnitParen = regexp.MustCompile(`(^|[^0-9])1\(`)

// w5Clean removes presentation artefacts that are not part of the mathematics:
// pure numbers wrapped as math, unit coefficients and doubled signs.
func w5Clean(s string) string {
	s = w5PureNumber.ReplaceAllString(s, "$1")
	s = strings.ReplaceAll(s, "+-", "-")
	s = strings.ReplaceAll(s, "--", "+")
	s = w5UnitCoefficient.ReplaceAllString(s, "${1}${2}")
	s = w5UnitParen.ReplaceAllString(s, "${1}(")
	s = strings.ReplaceAll(s, "  ", " ")
	return s
}

type w5Term struct {
	coef int
	body string
}

func w5SignedTerm(coef int, body string) string {
	if coef == 0 {
		return ""
	}
	sign := "+"
	mag := coef
	if mag < 0 {
		sign = "-"
		mag = -mag
	}
	core := body
	if mag != 1 || core == "" {
		core = fmt.Sprintf("%d%s", mag, core)
	}
	return sign + " " + core
}

func w5Poly(terms ...w5Term) string {
	var b strings.Builder
	for _, t := range terms {
		if t.coef == 0 {
			continue
		}
		part := w5SignedTerm(t.coef, t.body)
		if b.Len() == 0 {
			part = strings.TrimPrefix(part, "+ ")
			b.WriteString(part)
		} else {
			b.WriteString(" " + part)
		}
	}
	if b.Len() == 0 {
		return "0"
	}
	return b.String()
}

func w5Sum(vals ...int) string {
	parts := make([]string, len(vals))
	for i, v := range vals {
		if v < 0 {
			parts[i] = fmt.Sprintf("(%d)", v)
		} else {
			parts[i] = strconv.Itoa(v)
		}
	}
	return strings.Join(parts, " + ")
}

func w5Prod(a, b int) string { return fmt.Sprintf("(%d)(%d)", a, b) }

func w5SumParts(parts ...string) string { return strings.Join(parts, " + ") }

func w5PlusTerm(coef int, body string) string {
	if coef == 0 {
		return ""
	}
	return " " + w5SignedTerm(coef, body)
}

func w5GeneralStokes(kp int) generator.Problem {
	a, b := w5CR(1, 4), w5CR(1, 4)
	p, q := w5CNZ(-3, 3), w5CNZ(-3, 3)
	f := map[string]string{}
	w5CI(f, "a", a)
	w5CI(f, "b", b)
	w5CI(f, "p", p)
	w5CI(f, "q", q)
	form := w5Poly(w5Term{coef: p, body: " y\\,dx"}, w5Term{coef: q, body: " x\\,dy"})
	switch kp {
	case 1:
		return w5C(kp, fmt.Sprintf("For the rectangle \\([0,%d]\\times[0,%d]\\) with the counterclockwise boundary, what is \\(\\int_{\\partial R} x\\,dy\\)? (enter a number)", a, b), strconv.Itoa(a*b),
			fmt.Sprintf("Stokes gives the double integral of \\(dx\\wedge dy\\) over the rectangle, which is its area: %s.", w5Prod(a, b)), f)
	case 2:
		return w5C(kp, fmt.Sprintf("For \\(\\omega=%s\\), what is the coefficient of \\(dx\\wedge dy\\) in \\(d\\omega\\)? (enter a number)", form), strconv.Itoa(q-p),
			fmt.Sprintf("Differentiate the coefficients and alternate. The coefficient is %d.", q-p), f)
	default:
		ans := (q - p) * a * b
		return w5C(kp, fmt.Sprintf("Verify Stokes on \\([0,%d]\\times[0,%d]\\) for \\(\\omega=%s\\). What is \\(\\int_{\\partial R}\\omega\\)? (enter a number)", a, b, form), strconv.Itoa(ans),
			fmt.Sprintf("The coefficient of \\(dx\\wedge dy\\) in \\(d\\omega\\) is %d and the rectangle has area %d, so the boundary integral is %d.", q-p, a*b, ans), f)
	}
}

func w5Pullback(kp int) generator.Problem {
	a, b := w5CR(1, 5), w5CR(1, 5)
	f := map[string]string{}
	w5CI(f, "a", a)
	w5CI(f, "b", b)
	switch kp {
	case 1:
		return w5C(kp, fmt.Sprintf("For \\(\\phi(u,v)=(%d u,%d v)\\), what is the coefficient of \\(du\\) in \\(\\phi^*(dx)\\)? (enter a number)", a, b), strconv.Itoa(a),
			fmt.Sprintf("Pull back \\(x=%d u\\), so \\(dx=%d\\,du\\). The coefficient is %d.", a, a, a), f)
	case 2:
		return w5C(kp, fmt.Sprintf("For \\(\\phi(u,v)=(%d u,%d v)\\), what is the coefficient of \\(du\\wedge dv\\) in \\(\\phi^*(dx\\wedge dy)\\)? (enter a number)", a, b), strconv.Itoa(a*b),
			fmt.Sprintf("The pullback multiplies the coordinate differentials: \\(dx\\wedge dy=(%d\\,du)\\wedge(%d\\,dv)=%d\\,du\\wedge dv\\).", a, b, a*b), f)
	default:
		return w5C(kp, fmt.Sprintf("For the linear map \\(\\phi(u,v)=(%d u,%d v)\\), by what factor does it scale area? (enter a number)", a, b), strconv.Itoa(a*b),
			fmt.Sprintf("The Jacobian determinant is \\(%d\\cdot%d=%d\\), so areas scale by %d.", a, b, a*b, a*b), f)
	}
}

func w5ClosedExact(kp int) generator.Problem {
	a, b := w5CNZ(-4, 4), w5CNZ(-4, 4)
	x, y := w5CNZ(-4, 4), w5CNZ(-4, 4)
	f := map[string]string{}
	w5CI(f, "a", a)
	w5CI(f, "b", b)
	w5CI(f, "x", x)
	w5CI(f, "y", y)
	form := w5Poly(w5Term{coef: a, body: " y\\,dx"}, w5Term{coef: b, body: " x\\,dy"})
	potential := w5Poly(w5Term{coef: a, body: " x"}, w5Term{coef: b, body: " y"})
	switch kp {
	case 1:
		is := 0
		if a == b {
			is = 1
		}
		return w5C(kp, fmt.Sprintf("For \\(\\omega=%s\\), is \\(\\omega\\) closed? Enter 1 for yes and 0 for no.", form), strconv.Itoa(is),
			fmt.Sprintf("The coefficient of \\(dx\\wedge dy\\) in \\(d\\omega\\) is %d. A form is closed exactly when that coefficient is zero, so the answer is %d.", b-a, is), f)
	case 2:
		return w5C(kp, fmt.Sprintf("For the exact form \\(\\omega=%d\\,dx+%d\\,dy\\), a potential is \\(f=%s\\). What is \\(f\\) at \\((%d,%d)\\)? (enter a number)", a, b, potential, x, y), strconv.Itoa(a*x+b*y),
			fmt.Sprintf("Substitute the point into \\(%s\\): the value is %s.", potential, w5SumParts(w5Prod(a, x), w5Prod(b, y))), f)
	default:
		return w5C(kp, fmt.Sprintf("For \\(\\omega=%s\\), what is the coefficient of \\(dx\\wedge dy\\) in \\(d\\omega\\)? (enter a number)", form), strconv.Itoa(b-a),
			fmt.Sprintf("Differentiate the coefficients and alternate: the coefficient is %d.", b-a), f)
	}
}

func w5VectorPotential(kp int) generator.Problem {
	c := w5CNZ(-6, 6)
	x := w5CR(1, 4)
	lam := w5CNZ(-3, 3)
	f := map[string]string{}
	w5CI(f, "c", c)
	w5CI(f, "x", x)
	w5CI(f, "lambda", lam)
	switch kp {
	case 1:
		ay := w5Poly(w5Term{coef: c, body: " x"})
		return w5C(kp, fmt.Sprintf("For \\(F=(0,0,%d)\\), the vector potential has y-component \\(%s\\). What is that component at \\(x=%d\\)? (enter a number)", c, ay, x), strconv.Itoa(c*x),
			fmt.Sprintf("The y-component is %s, which equals %d at \\(x=%d\\).", ay, c*x, x), f)
	case 2:
		return w5C(kp, fmt.Sprintf("For \\(A=(0,%d x,0)\\), what is the z-component of \\(\\operatorname{curl} A\\)? (enter a number)", c), strconv.Itoa(c),
			fmt.Sprintf("The z-component is \\(\\partial_x A_y=%d\\).", c), f)
	default:
		ans := (c + lam) * x
		return w5C(kp, fmt.Sprintf("Add the gauge term \\(\\nabla(%d x y)\\) to \\(A=(0,%d x,0)\\). What is the new y-component at \\(x=%d\\)? (enter a number)", lam, c, x), strconv.Itoa(ans),
			fmt.Sprintf("The gradient adds \\(%d x\\) to the y-component, so the new value is %s.", lam, w5Prod(c+lam, x)), f)
	}
}

func w5Helmholtz(kp int) generator.Problem {
	a, b, c := w5CNZ(-4, 4), w5CNZ(-4, 4), w5CNZ(-4, 4)
	f := map[string]string{}
	w5CI(f, "a", a)
	w5CI(f, "b", b)
	w5CI(f, "c", c)
	switch kp {
	case 1:
		return w5C(kp, fmt.Sprintf("For \\(F=(%d x,%d y,%d z)\\), what is \\(\\nabla\\cdot F\\)? (enter a number)", a, b, c), strconv.Itoa(a+b+c),
			fmt.Sprintf("Add the coefficients: %s = %d.", w5Sum(a, b, c), a+b+c), f)
	case 2:
		return w5C(kp, fmt.Sprintf("For \\(F=(%d y,%d x,0)\\), what is the z-component of \\(\\operatorname{curl} F\\)? (enter a number)", a, b), strconv.Itoa(b-a),
			fmt.Sprintf("The z-component is the x-derivative of the y component minus the y-derivative of the x component: %d.", b-a), f)
	default:
		return w5C(kp, fmt.Sprintf("For \\(F=(%d y,%d x,0)\\), is \\(F\\) divergence-free? Enter 1 for yes and 0 for no.", a, b), "1",
			"The divergence is zero, so the answer is 1.", f)
	}
}

func w5Gauge(kp int) generator.Problem {
	a, b, c := w5CNZ(-3, 3), w5CNZ(-3, 3), w5CNZ(-3, 3)
	lam := w5CNZ(-3, 3)
	x, y, z := w5CR(1, 4), w5CR(1, 4), w5CR(1, 4)
	f := map[string]string{}
	w5CI(f, "a", a)
	w5CI(f, "b", b)
	w5CI(f, "c", c)
	w5CI(f, "lambda", lam)
	w5CI(f, "x", x)
	w5CI(f, "y", y)
	w5CI(f, "z", z)
	switch kp {
	case 1:
		return w5C(kp, fmt.Sprintf("For \\(A=(%d x,%d y,%d z)\\) and gauge function \\(g=%d(x+y+z)\\), what is the new x-component? (enter a number)", a, b, c, lam), strconv.Itoa(a+lam),
			fmt.Sprintf("The gradient adds %d to the x-component, so it becomes %s.", lam, w5Sum(a, lam)), f)
	case 2:
		return w5C(kp, fmt.Sprintf("Does adding \\(\\nabla g\\) to \\(A=(%d x,%d y,%d z)\\) change the magnetic field \\(\\nabla\\times A\\)? Enter 1 for yes and 0 for no.", a, b, c), "0",
			"A gradient is curl-free, so the curl of the potential is unchanged. The answer is 0.", f)
	default:
		return w5C(kp, fmt.Sprintf("For \\(g=%d(x+y+z)\\), what is \\(g(%d,%d,%d)\\)? (enter a number)", lam, x, y, z), strconv.Itoa(lam*(x+y+z)),
			fmt.Sprintf("Substitute the coordinates: %s = %d.", w5Prod(lam, x+y+z), lam*(x+y+z)), f)
	}
}

func w5Harmonic(kp int) generator.Problem {
	a, b := w5CNZ(-6, 6), w5CNZ(-6, 6)
	c := w5CNZ(-6, 6)
	f := map[string]string{}
	w5CI(f, "a", a)
	w5CI(f, "b", b)
	w5CI(f, "c", c)
	poly := w5Poly(w5Term{coef: a, body: " x^2"}, w5Term{coef: b, body: " y^2"}, w5Term{coef: c, body: " z^2"})
	switch kp {
	case 1:
		return w5C(kp, fmt.Sprintf("For \\(f=%s\\), what is \\(\\Delta f\\)? (enter a number)", poly), strconv.Itoa(2*(a+b+c)),
			fmt.Sprintf("The Laplacian is twice the coefficient sum: 2(%s) = %d.", w5Sum(a, b, c), 2*(a+b+c)), f)
	case 2:
		partial := w5Poly(w5Term{coef: a, body: " x^2"}, w5Term{coef: b, body: " y^2"}) + " + c z^2"
		return w5C(kp, fmt.Sprintf("For \\(f=%s\\), what value of \\(c\\) makes \\(f\\) harmonic? (enter a number)", partial), strconv.Itoa(-(a + b)),
			fmt.Sprintf("Harmonic means 2(%s + c) = 0, so \\(c=%d\\).", w5Sum(a, b), -(a+b)), f)
	default:
		plane := w5Poly(w5Term{coef: a, body: " x^2"}, w5Term{coef: 1, body: " b y^2"})
		return w5C(kp, fmt.Sprintf("For \\(f=%s\\) in the plane, what value of \\(b\\) makes \\(f\\) harmonic? (enter a number)", plane), strconv.Itoa(-a),
			fmt.Sprintf("The Laplacian is 2(%s), which is zero when \\(b=%d\\).", w5Sum(a, 0), -a), f)
	}
}

func w5MeanValue(kp int) generator.Problem {
	a, b, c := w5CNZ(-3, 3), w5CNZ(-3, 3), w5CNZ(-3, 3)
	x0, y0 := w5CNZ(-3, 3), w5CNZ(-3, 3)
	r := w5CR(1, 3)
	f := map[string]string{}
	w5CI(f, "a", a)
	w5CI(f, "b", b)
	w5CI(f, "c", c)
	w5CI(f, "x0", x0)
	w5CI(f, "y0", y0)
	w5CI(f, "r", r)
	poly := w5Poly(w5Term{coef: a, body: " x"}, w5Term{coef: b, body: " y"}, w5Term{coef: c})
	ans := a*x0 + b*y0 + c
	value := w5SumParts(w5Prod(a, x0), w5Prod(b, y0), strconv.Itoa(c)) + " = " + strconv.Itoa(ans)
	switch kp {
	case 1:
		return w5C(kp, fmt.Sprintf("For the harmonic function \\(f=%s\\), what is its average on the circle of radius %d centered at \\((%d,%d)\\)? (enter a number)", poly, r, x0, y0), strconv.Itoa(ans),
			fmt.Sprintf("The mean value property gives the centre value: %s.", value), f)
	case 2:
		return w5C(kp, fmt.Sprintf("For \\(f=%s\\), what is \\(f(%d,%d)\\)? (enter a number)", poly, x0, y0), strconv.Itoa(ans),
			fmt.Sprintf("Substitute the point: %s.", value), f)
	default:
		return w5C(kp, fmt.Sprintf("For \\(f=%s\\), what is the average of its values at the four cardinal points of the circle radius %d centered at \\((%d,%d)\\)? (enter a number)", poly, r, x0, y0), strconv.Itoa(ans),
			fmt.Sprintf("The four deviations cancel, so the average is the centre value %d.", ans), f)
	}
}

func w5MaximumPrinciple(kp int) generator.Problem {
	bnd := w5CR(1, 10)
	c := w5CNZ(-5, 5)
	f := map[string]string{}
	w5CI(f, "boundary", bnd)
	w5CI(f, "constant", c)
	quad := w5Poly(w5Term{coef: -4, body: " x^2"}, w5Term{coef: -2, body: " y^2"})
	switch kp {
	case 1:
		return w5C(kp, fmt.Sprintf("A nonconstant harmonic function on a disk has maximum %d on the boundary. What is its maximum on the closed disk? (enter a number)", bnd), strconv.Itoa(bnd),
			"The maximum principle says the maximum on the closed disk is attained on the boundary, so it is the given boundary maximum.", f)
	case 2:
		return w5C(kp, fmt.Sprintf("A harmonic function on a disk has boundary maximum %d and attains its maximum at an interior point. Is it constant? Enter 1 for yes and 0 for no.", bnd), "1",
			"If a harmonic function attains an interior maximum, the maximum principle forces it to be constant. The answer is 1.", f)
	default:
		return w5C(kp, fmt.Sprintf("A constant harmonic function equals %d everywhere. What is its maximum on the closed disk? (enter a number)", c), strconv.Itoa(c),
			fmt.Sprintf("A constant function has the same value everywhere, so the maximum is %d. The quadratic example %s is not involved here.", c, quad), f)
	}
}

func w5Poisson(kp int) generator.Problem {
	a, b := w5CNZ(-4, 4), w5CNZ(-4, 4)
	x, y := w5CNZ(-3, 3), w5CNZ(-3, 3)
	f := map[string]string{}
	w5CI(f, "a", a)
	w5CI(f, "b", b)
	w5CI(f, "x", x)
	w5CI(f, "y", y)
	switch kp {
	case 1:
		return w5C(kp, fmt.Sprintf("For \\(u=%d x^2+%d y^2\\), what is \\(\\Delta u\\)? (enter a number)", a, b), strconv.Itoa(2*(a+b)),
			fmt.Sprintf("The Laplacian is \\(2(%d+%d)=%d\\).", a, b, 2*(a+b)), f)
	case 2:
		c := 4 * w5CR(1, 10)
		w5CI(f, "constant", c)
		return w5C(kp, fmt.Sprintf("For \\(\\Delta u=%d\\), a particular solution is \\(u=k(x^2+y^2)\\). What is \\(k\\)? (enter a number)", c), strconv.Itoa(c/4),
			fmt.Sprintf("\\(\\Delta(k(x^2+y^2))=4k\\), so \\(k=%d/4=%d\\).", c, c/4), f)
	default:
		return w5C(kp, fmt.Sprintf("For \\(u=%d x^2+%d y^2\\), what is \\(u_{xx}+u_{yy}\\) at \\((%d,%d)\\)? (enter a number)", a, b, x, y), strconv.Itoa(2*(a+b)),
			fmt.Sprintf("The second derivatives are constant, so the value is \\(2(%d+%d)=%d\\).", a, b, 2*(a+b)), f)
	}
}

func w5Dirichlet(kp int) generator.Problem {
	c := w5CNZ(-5, 5)
	k := w5CNZ(-3, 3)
	x, y := w5CR(1, 4), w5CR(1, 4)
	f := map[string]string{}
	w5CI(f, "boundary", c)
	w5CI(f, "shift", k)
	w5CI(f, "x", x)
	w5CI(f, "y", y)
	switch kp {
	case 1:
		return w5C(kp, fmt.Sprintf("The Dirichlet problem on a disk has constant boundary value %d. What is the value at the centre? (enter a number)", c), strconv.Itoa(c),
			fmt.Sprintf("The unique harmonic extension of a constant boundary value is the same constant, %d.", c), f)
	case 2:
		return w5C(kp, fmt.Sprintf("The Dirichlet problem on a disk has constant boundary value %d. What is the value at the interior point \\((%d,%d)\\)? (enter a number)", c, x, y), strconv.Itoa(c),
			fmt.Sprintf("The solution is constant %d, so every interior value is %d.", c, c), f)
	default:
		return w5C(kp, fmt.Sprintf("Two harmonic functions on a disk have the same boundary data. If one equals %d at a point, what must the other equal there? (enter a number)", c), strconv.Itoa(c),
			fmt.Sprintf("Uniqueness forces the two solutions to be equal, so the other value is %d.", c), f)
	}
}

func w5Neumann(kp int) generator.Problem {
	a := w5CNZ(-6, 6)
	c, k := w5CNZ(-4, 4), w5CNZ(-3, 3)
	f := map[string]string{}
	w5CI(f, "a", a)
	w5CI(f, "c", c)
	w5CI(f, "shift", k)
	switch kp {
	case 1:
		return w5C(kp, fmt.Sprintf("A homogeneous Neumann problem on \\([0,1]\\) has \\(u'(0)=%d\\). For a solution to exist, what must \\(u'(1)\\) be? (enter a number)", a), strconv.Itoa(a),
			fmt.Sprintf("The compatibility condition is \\(u'(1)-u'(0)=0\\), so \\(u'(1)=%d\\).", a), f)
	case 2:
		return w5C(kp, fmt.Sprintf("For the Neumann data \\(u'(0)=%d\\) and \\(u'(1)=%d\\), what is the compatibility difference \\(u'(1)-u'(0)\\)? (enter a number)", a, a), "0",
			fmt.Sprintf("The two boundary derivatives are equal, so their difference is zero. The answer is 0."), f)
	default:
		return w5C(kp, fmt.Sprintf("A Neumann problem has solutions differing by a constant. If one solution has value %d at 0 and the shift is %d, what is the other value at 0? (enter a number)", c, k), strconv.Itoa(c+k),
			fmt.Sprintf("Adding the constant shift gives %s.", w5Sum(c, k)), f)
	}
}

func w5GreensIdentity(kp int) generator.Problem {
	a, b := w5CNZ(-4, 4), w5CNZ(-4, 4)
	f := map[string]string{}
	w5CI(f, "a", a)
	w5CI(f, "b", b)
	switch kp {
	case 1:
		return w5C(kp, fmt.Sprintf("For \\(u=%d x\\) and \\(v=%d x\\) on the unit square, what is \\(\\int_{\\partial R}u\\,\\partial v/\\partial n\\)? (enter a number)", a, b), strconv.Itoa(a*b),
			fmt.Sprintf("Green's first identity gives the area integral of \\(\\nabla u\\cdot\\nabla v=%d\\), so the boundary term is %d.", a*b, a*b), f)
	case 2:
		return w5C(kp, fmt.Sprintf("For the same \\(u=%d x\\) and \\(v=%d x\\), what is the boundary term in Green's first identity? (enter a number)", a, b), strconv.Itoa(a*b),
			fmt.Sprintf("The boundary term equals \\(\\int_R\\nabla u\\cdot\\nabla v=%d\\).", a*b), f)
	default:
		return w5C(kp, fmt.Sprintf("For \\(u=%d x\\) and \\(v=%d x\\), is the boundary term symmetric when u and v are swapped? Enter 1 for yes and 0 for no.", a, b), "1",
			"The product \\(\\nabla u\\cdot\\nabla v\\) is symmetric, so the boundary term is unchanged. The answer is 1.", f)
	}
}

func w5FundamentalSolution(kp int) generator.Problem {
	n := w5CR(3, 12)
	q := w5CNZ(-5, 5)
	f := map[string]string{}
	w5CI(f, "n", n)
	w5CI(f, "source", q)
	switch kp {
	case 1:
		return w5C(kp, fmt.Sprintf("In dimension \\(n=%d\\), the fundamental solution of the Laplacian is proportional to \\(r^p\\). What is \\(p=2-n\\)? (enter a number)", n), strconv.Itoa(2-n),
			fmt.Sprintf("\\(p=2-%d=%d\\).", n, 2-n), f)
	case 2:
		return w5C(kp, fmt.Sprintf("A point source has strength %d. What is the total flux through a small circle around it? (enter a number)", q), strconv.Itoa(q),
			fmt.Sprintf("The flux through any closed curve around a point source equals the source strength, %d.", q), f)
	default:
		return w5C(kp, fmt.Sprintf("A point source has strength %d. What is the total flux through a small sphere around it in three dimensions? (enter a number)", q), strconv.Itoa(q),
			fmt.Sprintf("The flux through any closed surface around a point source equals the source strength, %d.", q), f)
	}
}

func w5EulerLagrange(kp int) generator.Problem {
	a := 2 * w5CNZ(-6, 6)
	x := w5CR(1, 4)
	f := map[string]string{}
	w5CI(f, "a", a)
	w5CI(f, "x", x)
	lag := "\\tfrac12(y')^2" + w5PlusTerm(a, " y")
	switch kp {
	case 1:
		return w5C(kp, fmt.Sprintf("For \\(L=%s\\), the Euler-Lagrange equation is \\(y''=a\\). What is \\(a\\)? (enter a number)", lag), strconv.Itoa(a),
			fmt.Sprintf("The Euler-Lagrange equation gives \\(y''=%d\\).", a), f)
	case 2:
		return w5C(kp, fmt.Sprintf("Solve \\(y''=%d\\) with \\(y(0)=0\\) and \\(y'(0)=0\\). The solution is \\(y=kx^2\\). What is \\(k\\)? (enter a number)", a), strconv.Itoa(a/2),
			fmt.Sprintf("Integrating twice gives \\(y=(%d/2)x^2\\), so \\(k=%d\\).", a, a/2), f)
	default:
		is := 0
		if a == 0 {
			is = 1
		}
		return w5C(kp, fmt.Sprintf("For \\(L=%s\\), is every straight line \\(y=mx+b\\) a solution of the Euler-Lagrange equation? Enter 1 for yes and 0 for no.", lag), strconv.Itoa(is),
			fmt.Sprintf("The Euler-Lagrange equation is \\(y''=%d\\). A straight line has second derivative zero, so the answer is %d.", a, is), f)
	}
}

func w5FirstVariation(kp int) generator.Problem {
	a := 2 * w5CNZ(-6, 6)
	m := w5CNZ(-3, 3)
	x1 := w5CR(1, 4)
	y1 := m * x1
	f := map[string]string{}
	w5CI(f, "a", a)
	w5CI(f, "slope", m)
	w5CI(f, "x1", x1)
	w5CI(f, "y1", y1)
	integrand := "(y')^2" + w5PlusTerm(a, " y")
	switch kp {
	case 1:
		return w5C(kp, fmt.Sprintf("For \\(J[y]=\\int(%s)\\,dx\\), the first variation gives \\(2y''=a\\). What is \\(a\\)? (enter a number)", integrand), strconv.Itoa(a),
			fmt.Sprintf("The Euler-Lagrange equation is \\(2y''=%d\\), so \\(a=%d\\).", a, a), f)
	case 2:
		return w5C(kp, fmt.Sprintf("A stationary path for \\(J[y]=\\int(y')^2\\,dx\\) runs from \\((0,0)\\) to \\((%d,%d)\\). What is its constant slope? (enter a number)", x1, y1), strconv.Itoa(y1/x1),
			fmt.Sprintf("The stationary path is straight, so the slope is %d divided by %d, which is %d.", y1, x1, y1/x1), f)
	default:
		return w5C(kp, fmt.Sprintf("For \\(J[y]=\\int(%s)\\,dx\\), the first variation has the boundary term \\([y'\\delta y]\\) at the endpoints. If \\(\\delta y\\) vanishes at both endpoints, what is that boundary term? (enter a number)", integrand), "0",
			"The variation vanishes at the endpoints, so the boundary term is zero. The answer is 0.", f)
	}
}

func w5Geodesic(kp int) generator.Problem {
	triples := [][3]int{{3, 4, 5}, {4, 3, 5}, {5, 12, 13}, {12, 5, 13}, {8, 15, 17}, {15, 8, 17}}
	t := triples[rand.Intn(len(triples))]
	p, q := t[0], t[1]
	hyp := t[2]
	k := w5CR(1, 4)
	f := map[string]string{}
	w5CI(f, "leg1", p)
	w5CI(f, "leg2", q)
	w5CI(f, "hypotenuse", hyp)
	w5CI(f, "scale", k)
	switch kp {
	case 1:
		return w5C(kp, fmt.Sprintf("What is the straight-line distance from \\((0,0)\\) to \\((%d,%d)\\)? (enter a number)", p*k, q*k), strconv.Itoa(hyp*k),
			fmt.Sprintf("The displacement has legs %d and %d, so the distance is %d.", p*k, q*k, hyp*k), f)
	case 2:
		return w5C(kp, fmt.Sprintf("For the metric \\(ds^2=%d^2 dx^2+%d^2 dy^2\\), what is the length of the horizontal segment from \\((0,0)\\) to \\((1,0)\\)? (enter a number)", p*k, q*k), strconv.Itoa(p*k),
			fmt.Sprintf("Along the horizontal segment \\(dy=0\\), so the length is %d.", p*k), f)
	default:
		return w5C(kp, fmt.Sprintf("For the metric \\(ds^2=%d^2 dx^2+%d^2 dy^2\\), what is the length of the segment from \\((0,0)\\) to \\((1,1)\\)? (enter a number)", p*k, q*k), strconv.Itoa(hyp*k),
			fmt.Sprintf("The length is the squared root of the sum of the squares of %d and %d, which is %d.", p*k, q*k, hyp*k), f)
	}
}

func w5MinimalSurface(kp int) generator.Problem {
	legs := [][3]int{{2, 2, 3}, {2, -2, 3}, {4, 8, 9}, {8, 4, 9}}
	pair := legs[rand.Intn(len(legs))]
	a, b, factor := pair[0], pair[1], pair[2]
	if kp != 1 {
		a, b, factor = w5CNZ(-7, 7), w5CNZ(-7, 7), 0
	}
	w, h := w5CR(1, 4), w5CR(1, 4)
	f := map[string]string{}
	w5CI(f, "a", a)
	w5CI(f, "b", b)
	w5CI(f, "factor", factor)
	w5CI(f, "w", w)
	w5CI(f, "h", h)
	plane := w5Poly(w5Term{coef: a, body: " x"}, w5Term{coef: b, body: " y"})
	switch kp {
	case 1:
		return w5C(kp, fmt.Sprintf("For the graph \\(z=%s\\) over a rectangle of area %d, what is the surface area? (enter a number)", plane, w*h), strconv.Itoa(factor*w*h),
			fmt.Sprintf("The area element has coefficient %d, so the area is %d times %d, which is %d.", factor, factor, w*h, factor*w*h), f)
	case 2:
		return w5C(kp, fmt.Sprintf("For the minimal-surface equation, the coefficient of \\(z_{xx}\\) is \\(1+z_y^2\\). For \\(z=%s\\), what is that coefficient? (enter a number)", plane), strconv.Itoa(1+b*b),
			fmt.Sprintf("Here \\(z_y=%d\\), so the coefficient is 1 plus the square of %d, which is %d.", b, b, 1+b*b), f)
	default:
		return w5C(kp, fmt.Sprintf("Is the plane \\(z=%s\\) a minimal surface? Enter 1 for yes and 0 for no.", plane), "1",
			"A plane has zero mean curvature, so it is a minimal surface. The answer is 1.", f)
	}
}

func w5SurfaceMeasure(kp int) generator.Problem {
	a, b := w5CR(1, 4), w5CR(1, 4)
	m, n := w5CR(1, 4), w5CR(1, 4)
	c, d := w5CR(1, 3), w5CR(1, 3)
	f := map[string]string{}
	w5CI(f, "a", a)
	w5CI(f, "b", b)
	w5CI(f, "m", m)
	w5CI(f, "n", n)
	w5CI(f, "c", c)
	w5CI(f, "d", d)
	switch kp {
	case 1:
		return w5C(kp, fmt.Sprintf("For \\(r(u,v)=(%d u,%d v,0)\\), what is the area element coefficient \\(|r_u\\times r_v|\\)? (enter a number)", a, b), strconv.Itoa(a*b),
			fmt.Sprintf("The cross product has length %d.", a*b), f)
	case 2:
		return w5C(kp, fmt.Sprintf("For \\(r(u,v)=(%d u,%d v,0)\\) over \\([0,%d]\\times[0,%d]\\), what is the surface area? (enter a number)", a, b, m, n), strconv.Itoa(a*b*m*n),
			fmt.Sprintf("The area is %d times the parameter-domain area %d, which is %d.", a*b, m*n, a*b*m*n), f)
	default:
		return w5C(kp, fmt.Sprintf("If the parameters are rescaled by \\(u\\mapsto %d u\\) and \\(v\\mapsto %d v\\), by what factor does the area element scale? (enter a number)", c, d), strconv.Itoa(c*d),
			fmt.Sprintf("Both tangent vectors scale, so the area element scales by %d.", c*d), f)
	}
}

func w5Fourier(kp int) generator.Problem {
	f := [2][2]int{}
	for i := range f {
		for j := range f[i] {
			f[i][j] = w5CR(-3, 3)
		}
	}
	ff := map[string]string{}
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			w5CI(ff, fmt.Sprintf("f%d%d", i, j), f[i][j])
		}
	}
	switch kp {
	case 1:
		s := f[0][0] + f[0][1] + f[1][0] + f[1][1]
		return w5C(kp, fmt.Sprintf("For the 2x2 array \\(f=%v\\), what is the Fourier coefficient \\(F(0,0)\\) equal to the sum of all entries? (enter a number)", f), strconv.Itoa(s),
			fmt.Sprintf("The zero-frequency coefficient is the sum of the entries, which is %d.", s), ff)
	case 2:
		s := f[0][0] - f[0][1] + f[1][0] - f[1][1]
		return w5C(kp, fmt.Sprintf("For the 2x2 array \\(f=%v\\), what is the coefficient \\(F(0,1)=f_{00}-f_{01}+f_{10}-f_{11}\\)? (enter a number)", f), strconv.Itoa(s),
			fmt.Sprintf("Substitute the entries: the alternating sum is %d.", s), ff)
	default:
		s := f[0][0] - f[0][1] - f[1][0] + f[1][1]
		return w5C(kp, fmt.Sprintf("For the 2x2 array \\(f=%v\\), what is the coefficient \\(F(1,1)=f_{00}-f_{01}-f_{10}+f_{11}\\)? (enter a number)", f), strconv.Itoa(s),
			fmt.Sprintf("Substitute the entries: the result is %d.", s), ff)
	}
}

func w5Convolution(kp int) generator.Problem {
	f := [2][2]int{}
	g := [2][2]int{}
	for i := range f {
		for j := range f[i] {
			f[i][j] = w5CR(-2, 2)
			g[i][j] = w5CR(-2, 2)
		}
	}
	facts := map[string]string{}
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			w5CI(facts, fmt.Sprintf("f%d%d", i, j), f[i][j])
			w5CI(facts, fmt.Sprintf("g%d%d", i, j), g[i][j])
		}
	}
	h := [2][2]int{}
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for a := 0; a < 2; a++ {
				for b := 0; b < 2; b++ {
					h[i][j] += f[a][b] * g[(i-a+2)%2][(j-b+2)%2]
				}
			}
		}
	}
	switch kp {
	case 1:
		return w5C(kp, fmt.Sprintf("For the 2x2 arrays \\(f=%v\\) and \\(g=%v\\), what is the circular convolution \\(h(0,0)\\)? (enter a number)", f, g), strconv.Itoa(h[0][0]),
			fmt.Sprintf("Sum the four products that contribute to \\(h(0,0)\\): the value is %d.", h[0][0]), facts)
	case 2:
		return w5C(kp, fmt.Sprintf("For the same arrays \\(f=%v\\) and \\(g=%v\\), what is \\(h(1,1)\\)? (enter a number)", f, g), strconv.Itoa(h[1][1]),
			fmt.Sprintf("The wrap-around convolution gives \\(h(1,1)=%d\\).", h[1][1]), facts)
	default:
		s := h[0][0] + h[0][1] + h[1][0] + h[1][1]
		return w5C(kp, fmt.Sprintf("For the same arrays \\(f=%v\\) and \\(g=%v\\), what is the sum of all four convolution entries? (enter a number)", f, g), strconv.Itoa(s),
			fmt.Sprintf("The total convolution sum is the product of the total sums, which is %d.", s), facts)
	}
}
