package calculus

// Wave 5 Calculus III: the first five review-order concepts on oriented
// surfaces and differential forms. Every family publishes the exact operands
// its answer follows from, and wave5_test.go recomputes each answer from those
// facts.

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave5(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"calc.vec.orientation":         "orientation",
		"calc.vec.oriented_flux":       "oriented_flux",
		"calc.vec.differential_form":   "differential_form",
		"calc.vec.wedge_product":       "wedge_product",
		"calc.vec.exterior_derivative": "exterior_derivative",
	} {
		reg.Register(id, &w5CalcGen{kind: kind})
	}
}

type w5CalcGen struct{ kind string }

func (g *w5CalcGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "orientation":
		return w5Orientation(kp)
	case "oriented_flux":
		return w5OrientedFlux(kp)
	case "differential_form":
		return w5DifferentialForm(kp)
	case "wedge_product":
		return w5WedgeProduct(kp)
	case "exterior_derivative":
		return w5ExteriorDerivative(kp)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

func w5C(o int, q, a, ex string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}

func w5CI(m map[string]string, k string, v int) { m[k] = strconv.Itoa(v) }

func w5CR(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func w5CNZ(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	for {
		v := w5CR(lo, hi)
		if v != 0 {
			return v
		}
	}
}

func w5Combined(coefs []int, names []string) string {
	parts := make([]string, 0, len(coefs))
	for i, c := range coefs {
		if c == 0 {
			continue
		}
		sign := "+"
		if c < 0 {
			sign = "-"
			c = -c
		}
		core := names[i]
		if c != 1 {
			core = fmt.Sprintf("%d%s", c, core)
		}
		if len(parts) == 0 {
			if sign == "-" {
				parts = append(parts, "-"+core)
			} else {
				parts = append(parts, core)
			}
			continue
		}
		parts = append(parts, sign+" "+core)
	}
	if len(parts) == 0 {
		return "0"
	}
	return strings.Join(parts, " ")
}

func w5Vec(vals ...int) string {
	parts := make([]string, len(vals))
	for i, v := range vals {
		parts[i] = strconv.Itoa(v)
	}
	return "(" + strings.Join(parts, ",") + ")"
}

func w5Sub(a, b int) string {
	if b < 0 {
		return fmt.Sprintf("%d-(%d)", a, b)
	}
	return fmt.Sprintf("%d-%d", a, b)
}

func w5Orientation(kp int) generator.Problem {
	a, b := w5CNZ(-4, 4), w5CNZ(-4, 4)
	f := map[string]string{}
	w5CI(f, "a", a)
	w5CI(f, "b", b)
	switch kp {
	case 1:
		return w5C(kp, fmt.Sprintf("For \\(r(u,v)=(u,v,%s)\\), what is the x-component of \\(r_u\\times r_v\\)? (enter a number)", w5Combined([]int{a, b}, []string{"u", "v"})), strconv.Itoa(-a),
			fmt.Sprintf("\\(r_u=(1,0,%d)\\) and \\(r_v=(0,1,%d)\\). Their cross product is \\(%s\\), so the x-component is %d.", a, b, w5Vec(-a, -b, 1), -a), f)
	case 2:
		return w5C(kp, fmt.Sprintf("For the reparametrised surface \\(r(u,v)=(v,u,%s)\\), what is the z-component of \\(r_u\\times r_v\\)? (enter a number)", w5Combined([]int{a, b}, []string{"u", "v"})), "-1",
			fmt.Sprintf("Swap the parameters: \\(r_u=(0,1,%d)\\) and \\(r_v=(1,0,%d)\\). The cross product is \\(%s\\), so the z-component is -1.", b, a, w5Vec(a, b, -1)), f)
	default:
		c := w5CNZ(-4, 4)
		w5CI(f, "c", c)
		up := 0
		if c > 0 {
			up = 1
		}
		return w5C(kp, fmt.Sprintf("A surface has normal vector \\((%d,%d,%d)\\). Is this the upward orientation? Enter 1 for yes and 0 for no.", a, b, c), strconv.Itoa(up),
			fmt.Sprintf("The upward orientation has a positive z-component. Here the z-component is %d, so the answer is %d.", c, up), f)
	}
}

func w5OrientedFlux(kp int) generator.Problem {
	a, b := w5CNZ(-3, 3), w5CNZ(-3, 3)
	p, q, r := w5CNZ(-4, 4), w5CNZ(-4, 4), w5CNZ(-4, 4)
	f := map[string]string{}
	w5CI(f, "a", a)
	w5CI(f, "b", b)
	w5CI(f, "p", p)
	w5CI(f, "q", q)
	w5CI(f, "r", r)
	switch kp {
	case 1:
		ans := -a*p - b*q + r
		return w5C(kp, fmt.Sprintf("The surface \\(z=%s\\) over the unit square is oriented upward. For \\(F=%s\\), what is the flux? (enter a number)", w5Combined([]int{a, b}, []string{"x", "y"}), w5Vec(p, q, r)), strconv.Itoa(ans),
			fmt.Sprintf("The upward normal is \\(%s\\). Dotting with F gives \\((%d)(%d)+(%d)(%d)+(1)(%d)=%d\\).", w5Vec(-a, -b, 1), -a, p, -b, q, r, ans), f)
	case 2:
		ans := a*p + b*q - r
		return w5C(kp, fmt.Sprintf("The same graph \\(z=%s\\) is now oriented downward. For \\(F=%s\\), what is the flux? (enter a number)", w5Combined([]int{a, b}, []string{"x", "y"}), w5Vec(p, q, r)), strconv.Itoa(ans),
			fmt.Sprintf("The downward normal is \\(%s\\). Dotting with F gives \\((%d)(%d)+(%d)(%d)+(-1)(%d)=%d\\).", w5Vec(a, b, -1), a, p, b, q, r, ans), f)
	default:
		f["surface"] = "vertical"
		return w5C(kp, fmt.Sprintf("The unit square is parametrised by \\(r(u,v)=(u,0,v)\\) with \\(0\\le u,v\\le 1\\). For \\(F=%s\\), what is the oriented flux using \\(r_u\\times r_v\\)? (enter a number)", w5Vec(p, q, r)), strconv.Itoa(-q),
			fmt.Sprintf("\\(r_u\\times r_v=(0,-1,0)\\). Its dot product with F is \\((-1)(%d)=%d\\).", q, -q), f)
	}
}

func w5DifferentialForm(kp int) generator.Problem {
	a, b, c := w5CNZ(-3, 3), w5CNZ(-3, 3), w5CNZ(-3, 3)
	f := map[string]string{}
	w5CI(f, "a", a)
	w5CI(f, "b", b)
	w5CI(f, "c", c)
	switch kp {
	case 1:
		p, q, r := w5CNZ(-3, 3), w5CNZ(-3, 3), w5CNZ(-3, 3)
		w5CI(f, "p", p)
		w5CI(f, "q", q)
		w5CI(f, "r", r)
		ans := a*p + b*q + c*r
		return w5C(kp, fmt.Sprintf("Evaluate the 1-form \\(%s\\) on the vector \\(%s\\). (enter a number)", w5Combined([]int{a, b, c}, []string{"dx", "dy", "dz"}), w5Vec(p, q, r)), strconv.Itoa(ans),
			fmt.Sprintf("A 1-form pairs with a vector by a dot product: \\((%d)(%d)+(%d)(%d)+(%d)(%d)=%d\\).", a, p, b, q, c, r, ans), f)
	case 2:
		u1, u2, u3 := w5CNZ(-3, 3), w5CNZ(-3, 3), w5CNZ(-3, 3)
		v1, v2, v3 := w5CNZ(-3, 3), w5CNZ(-3, 3), w5CNZ(-3, 3)
		f["u"] = fmt.Sprintf("%d,%d,%d", u1, u2, u3)
		f["v"] = fmt.Sprintf("%d,%d,%d", v1, v2, v3)
		d1 := u2*v3 - u3*v2
		d2 := u3*v1 - u1*v3
		d3 := u1*v2 - u2*v1
		ans := a*d1 + b*d2 + c*d3
		return w5C(kp, fmt.Sprintf("Evaluate \\(%s\\) on \\(u=%s\\) and \\(v=%s\\). (enter a number)", w5Combined([]int{a, b, c}, []string{"dy\\wedge dz", "dz\\wedge dx", "dx\\wedge dy"}), w5Vec(u1, u2, u3), w5Vec(v1, v2, v3)), strconv.Itoa(ans),
			fmt.Sprintf("The three coordinate determinants are %d, %d and %d. The value is \\((%d)(%d)+(%d)(%d)+(%d)(%d)=%d\\).", d1, d2, d3, a, d1, b, d2, c, d3, ans), f)
	default:
		n := w5CR(3, 7)
		k := w5CR(1, n-1)
		w5CI(f, "n", n)
		w5CI(f, "k", k)
		ans := 1
		for i := 1; i <= k; i++ {
			ans = ans * (n - k + i) / i
		}
		return w5C(kp, fmt.Sprintf("How many independent \\(k\\)-form components are there on \\(\\mathbb{R}^{%d}\\) when \\(k=%d\\)? (enter a number)", n, k), strconv.Itoa(ans),
			fmt.Sprintf("A \\(k\\)-form on \\(\\mathbb{R}^{%d}\\) is determined by its values on ordered \\(k\\)-tuples, so the count is \\(\\binom{%d}{%d}=%d\\).", n, n, k, ans), f)
	}
}

func w5WedgeProduct(kp int) generator.Problem {
	a, b, c, d := w5CNZ(-3, 3), w5CNZ(-3, 3), w5CNZ(-3, 3), w5CNZ(-3, 3)
	f := map[string]string{}
	w5CI(f, "a", a)
	w5CI(f, "b", b)
	w5CI(f, "c", c)
	w5CI(f, "d", d)
	det := a*d - b*c
	switch kp {
	case 1:
		return w5C(kp, fmt.Sprintf("What is the coefficient of \\(dx\\wedge dy\\) in \\((%s)\\wedge(%s)\\)? (enter a number)", w5Combined([]int{a, b}, []string{"dx", "dy"}), w5Combined([]int{c, d}, []string{"dx", "dy"})), strconv.Itoa(det),
			fmt.Sprintf("The coefficient is the 2x2 determinant \\((%d)(%d)-(%d)(%d)=%d\\).", a, d, b, c, det), f)
	case 2:
		return w5C(kp, fmt.Sprintf("What is the coefficient of \\(dy\\wedge dx\\) in \\((%s)\\wedge(%s)\\)? (enter a number)", w5Combined([]int{a, b}, []string{"dx", "dy"}), w5Combined([]int{c, d}, []string{"dx", "dy"})), strconv.Itoa(-det),
			fmt.Sprintf("The wedge product is antisymmetric, so \\(dy\\wedge dx=-dx\\wedge dy\\). The coefficient is \\(-(%d)=%d\\).", det, -det), f)
	default:
		p, q, r := w5CNZ(-3, 3), w5CNZ(-3, 3), w5CNZ(-3, 3)
		w5CI(f, "p", p)
		w5CI(f, "q", q)
		w5CI(f, "r", r)
		ans := b*r - c*q
		return w5C(kp, fmt.Sprintf("For \\(\\omega=%s\\) and \\(\\eta=%s\\), what is the coefficient of \\(dy\\wedge dz\\) in \\(\\omega\\wedge\\eta\\)? (enter a number)", w5Combined([]int{a, b, c}, []string{"dx", "dy", "dz"}), w5Combined([]int{p, q, r}, []string{"dx", "dy", "dz"})), strconv.Itoa(ans),
			fmt.Sprintf("Collect the \\(dy\\wedge dz\\) terms: \\((%d)(%d)-(%d)(%d)=%d\\).", b, r, c, q, ans), f)
	}
}

func w5ExteriorDerivative(kp int) generator.Problem {
	a, b, c := w5CNZ(-3, 3), w5CNZ(-3, 3), w5CNZ(-3, 3)
	f := map[string]string{}
	w5CI(f, "a", a)
	w5CI(f, "b", b)
	w5CI(f, "c", c)
	switch kp {
	case 1:
		return w5C(kp, fmt.Sprintf("For \\(f=%s\\), what is the coefficient of \\(dx\\) in \\(df\\)? (enter a number)", w5Combined([]int{a, b, c}, []string{"x", "y", "z"})), strconv.Itoa(a),
			fmt.Sprintf("\\(df=f_x\\,dx+f_y\\,dy+f_z\\,dz\\). Differentiate with respect to x: the coefficient is %d.", a), f)
	case 2:
		return w5C(kp, fmt.Sprintf("For \\(\\omega=%s\\), what is the coefficient of \\(dx\\wedge dy\\) in \\(d\\omega\\)? (enter a number)", w5Combined([]int{a, b, c}, []string{"y\\,dx", "x\\,dy", "z\\,dz"})), strconv.Itoa(b-a),
			fmt.Sprintf("\\(d(a y\\,dx)=a\\,dy\\wedge dx=-a\\,dx\\wedge dy\\) and \\(d(b x\\,dy)=b\\,dx\\wedge dy\\). The coefficient is \\(%s=%d\\).", w5Sub(b, a), b-a), f)
	default:
		return w5C(kp, fmt.Sprintf("For \\(\\eta=%s\\), what is the coefficient of \\(dx\\wedge dy\\wedge dz\\) in \\(d\\eta\\)? (enter a number)", w5Combined([]int{a, b, c}, []string{"x\\,dy\\wedge dz", "y\\,dz\\wedge dx", "z\\,dx\\wedge dy"})), strconv.Itoa(a+b+c),
			fmt.Sprintf("Differentiate each coefficient and use cyclic permutations. The three contributions are %d, %d and %d, with sum %d.", a, b, c, a+b+c), f)
	}
}
