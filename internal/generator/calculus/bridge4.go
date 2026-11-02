package calculus

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

// calc3BridgeGen serves the calculus III concepts added by the content
// expansion. The kind selects the concept and KPOrdinal selects the lesson
// section.
type calc3BridgeGen struct{ kind string }

func (g *calc3BridgeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := calcBridgeOrdinal(ctx)
	switch g.kind {
	case "quadric_surfaces":
		switch o {
		case 1:
			type entry struct {
				eq  string
				yes bool
			}
			pool := []entry{
				{"\\(x^2+y^2+z^2=1\\)", true},
				{"\\(x^2/4+y^2/9+z^2=1\\)", true},
				{"\\(x^2/9+y^2+z^2/4=1\\)", true},
				{"\\(x^2+y^2/4+z^2/9=4\\)", true},
				{"\\(x^2+y^2-z^2=1\\)", false},
				{"\\(z=x^2+y^2\\)", false},
				{"\\(z=x^2-y^2\\)", false},
				{"\\(x^2+y^2=z^2\\)", false},
			}
			perm := rand.Perm(len(pool))[:4]
			count := 0
			parts := make([]string, 0, 4)
			for i, j := range perm {
				if pool[j].yes {
					count++
				}
				parts = append(parts, fmt.Sprintf("(%c) %s", 'a'+i, pool[j].eq))
			}
			return calcBridgeProblem(o,
				fmt.Sprintf("How many of these equations represent an ellipsoid?\n%s\n(enter a number)", strings.Join(parts, "\n")),
				strconv.Itoa(count),
				fmt.Sprintf("Exactly %d of the four equations represent an ellipsoid; the others are hyperboloids, paraboloids or a cone.", count),
				map[string]string{"kind": "ellipsoid-count"})
		case 2:
			type entry struct {
				eq  string
				yes bool
			}
			pool := []entry{
				{"\\(z=x^2+y^2\\)", true},
				{"\\(z=x^2/4+y^2/9\\)", true},
				{"\\(z=2x^2+3y^2\\)", true},
				{"\\(z=x^2+y^2+1\\)", true},
				{"\\(z=x^2-y^2\\)", false},
				{"\\(x^2+y^2-z^2=1\\)", false},
				{"\\(x^2+y^2=z^2\\)", false},
				{"\\(x^2+y^2+z^2=1\\)", false},
			}
			perm := rand.Perm(len(pool))[:4]
			count := 0
			parts := make([]string, 0, 4)
			for i, j := range perm {
				if pool[j].yes {
					count++
				}
				parts = append(parts, fmt.Sprintf("(%c) %s", 'a'+i, pool[j].eq))
			}
			return calcBridgeProblem(o,
				fmt.Sprintf("How many of these equations represent an elliptic paraboloid?\n%s\n(enter a number)", strings.Join(parts, "\n")),
				strconv.Itoa(count),
				fmt.Sprintf("Exactly %d of the four equations represent an elliptic paraboloid; the hyperbolic paraboloid is not one of them.", count),
				map[string]string{"kind": "paraboloid-count"})
		default:
			a := 1 + rand.Intn(3)
			b := 1 + rand.Intn(3)
			p := rand.Intn(5) - 2
			q := rand.Intn(5) - 2
			ans := a*p*p + b*q*q
			return calcBridgeProblem(o,
				fmt.Sprintf("For the elliptic paraboloid \\(z=%d x^2+%d y^2\\), what is \\(z\\) at \\((x,y)=(%d,%d)\\)? (enter a number)", a, b, p, q),
				strconv.Itoa(ans),
				fmt.Sprintf("Substitute the point: z = %d(%d)^2 + %d(%d)^2 = %d.", a, p, b, q, ans),
				map[string]string{"kind": "evaluate"})
		}
	case "triple_integral_region":
		switch o {
		case 1:
			a := 2 + rand.Intn(5)
			b := 2 + rand.Intn(5)
			c := 2 + rand.Intn(5)
			ans := a * b * c
			return calcBridgeProblem(o,
				fmt.Sprintf("What is the volume of the box \\(0\\le x\\le %d\\), \\(0\\le y\\le %d\\), \\(0\\le z\\le %d\\)? (enter a number)", a, b, c),
				strconv.Itoa(ans),
				fmt.Sprintf("The volume of a box is the product of its side lengths: %d x %d x %d = %d.", a, b, c, ans),
				map[string]string{"kind": "box"})
		case 2:
			m := 1 + rand.Intn(3)
			a := 6 * m
			b := 2 + rand.Intn(4)
			c := 2 + rand.Intn(4)
			ans := a * b * c / 6
			return calcBridgeProblem(o,
				fmt.Sprintf("Find the volume of the tetrahedron in the first octant bounded by \\(x/%d+y/%d+z/%d=1\\) and the coordinate planes. (enter a number)", a, b, c),
				strconv.Itoa(ans),
				fmt.Sprintf("The volume of this tetrahedron is abc/6 = (%d)(%d)(%d)/6 = %d.", a, b, c, ans),
				map[string]string{"kind": "tetrahedron"})
		default:
			a := 2 + rand.Intn(4)
			b := 2 + rand.Intn(4)
			c := 2 + rand.Intn(4)
			d := 2 + rand.Intn(4)
			ans := d * a * b * c
			return calcBridgeProblem(o,
				fmt.Sprintf("A box \\(0\\le x\\le %d\\), \\(0\\le y\\le %d\\), \\(0\\le z\\le %d\\) has constant density %d. What is its mass? (enter a number)", a, b, c, d),
				strconv.Itoa(ans),
				fmt.Sprintf("Mass is density times volume: %d x %d x %d x %d = %d.", d, a, b, c, ans),
				map[string]string{"kind": "mass"})
		}
	case "vector_field":
		switch o {
		case 1:
			a := 1 + rand.Intn(4)
			b := 1 + rand.Intn(4)
			c := 1 + rand.Intn(4)
			ans := a + b + c
			return calcBridgeProblem(o,
				fmt.Sprintf("For \\(\\mathbf F=\\langle %d x,%d y,%d z\\rangle\\), what is the divergence? (enter a number)", a, b, c),
				strconv.Itoa(ans),
				fmt.Sprintf("The divergence is the sum of the partial derivatives: %d + %d + %d = %d.", a, b, c, ans),
				map[string]string{"kind": "divergence"})
		case 2:
			a := 1 + rand.Intn(4)
			b := 1 + rand.Intn(4)
			c := 1 + rand.Intn(4)
			vals := []int{-b, -c, -a}
			names := []string{"x", "y", "z"}
			k := rand.Intn(3)
			return calcBridgeProblem(o,
				fmt.Sprintf("For \\(\\mathbf F=\\langle %d y,%d z,%d x\\rangle\\), what is the %s-component of the curl? (enter a number)", a, b, c, names[k]),
				strconv.Itoa(vals[k]),
				fmt.Sprintf("The curl is \\(\\langle %d,%d,%d\\rangle\\), so the %s-component is %d.", -b, -c, -a, names[k], vals[k]),
				map[string]string{"kind": "curl"})
		default:
			a := 1 + rand.Intn(4)
			b := 1 + rand.Intn(4)
			c := 1 + rand.Intn(4)
			p := rand.Intn(5) - 2
			q := rand.Intn(5) - 2
			r := rand.Intn(5) - 2
			vals := []int{a * q, b * r, c * p}
			names := []string{"x", "y", "z"}
			k := rand.Intn(3)
			return calcBridgeProblem(o,
				fmt.Sprintf("For \\(\\mathbf F=\\langle %d y,%d z,%d x\\rangle\\), what is the %s-component at \\((%d,%d,%d)\\)? (enter a number)", a, b, c, names[k], p, q, r),
				strconv.Itoa(vals[k]),
				fmt.Sprintf("Substitute the point into the %s-component, giving %d.", names[k], vals[k]),
				map[string]string{"kind": "evaluate"})
		}
	case "surface_tangent":
		a := 1 + rand.Intn(3)
		b := 1 + rand.Intn(3)
		p := rand.Intn(4) - 1
		q := rand.Intn(4) - 1
		switch o {
		case 1:
			ans := 2 * a * p
			return calcBridgeProblem(o,
				fmt.Sprintf("For the surface \\(z=%d x^2+%d y^2\\), what is the x-component of a normal vector at \\((x,y)=(%d,%d)\\)? (enter a number)", a, b, p, q),
				strconv.Itoa(ans),
				fmt.Sprintf("A normal is \\(\\langle f_x,f_y,-1\\rangle\\), and f_x = %d x, so at x=%d the x-component is %d.", 2*a, p, ans),
				map[string]string{"kind": "normal-x"})
		case 2:
			ans := 2 * b * q
			return calcBridgeProblem(o,
				fmt.Sprintf("For the surface \\(z=%d x^2+%d y^2\\), what is the y-component of a normal vector at \\((x,y)=(%d,%d)\\)? (enter a number)", a, b, p, q),
				strconv.Itoa(ans),
				fmt.Sprintf("A normal is \\(\\langle f_x,f_y,-1\\rangle\\), and f_y = %d y, so at y=%d the y-component is %d.", 2*b, q, ans),
				map[string]string{"kind": "normal-y"})
		default:
			z0 := a*p*p + b*q*q
			ans := z0 + 2*a*p
			return calcBridgeProblem(o,
				fmt.Sprintf("For \\(z=%d x^2+%d y^2\\), use the tangent plane at \\((x,y)=(%d,%d)\\) to estimate \\(z\\) at \\((%d,%d)\\). (enter a number)", a, b, p, q, p+1, q),
				strconv.Itoa(ans),
				fmt.Sprintf("The tangent-plane estimate is f(%d,%d)+f_x(%d,%d) = %d + %d = %d.", p, q, p, q, z0, 2*a*p, ans),
				map[string]string{"kind": "estimate"})
		}
	case "multi_continuity":
		switch o {
		case 1:
			a := rand.Intn(5) - 2
			b := rand.Intn(5) - 2
			ans := a*a + b*b
			return calcBridgeProblem(o,
				fmt.Sprintf("For \\(f(x,y)=x^2+y^2\\), what is \\(\\lim_{(x,y)\\to(%d,%d)} f(x,y)\\)? (enter a number)", a, b),
				strconv.Itoa(ans),
				fmt.Sprintf("The function is a polynomial, so substitute the point: %d^2 + %d^2 = %d.", a, b, ans),
				map[string]string{"kind": "polynomial"})
		case 2:
			m := []int{0, 2, 3, 4, 5}[rand.Intn(5)]
			ans := calcBridgeDecimal(float64(1-m*m) / float64(1+m*m))
			return calcBridgeProblem(o,
				fmt.Sprintf("For \\(f(x,y)=(x^2-y^2)/(x^2+y^2)\\), what is the limit along the path \\(y=%d x\\) as \\((x,y)\\to(0,0)\\)? (enter a decimal)", m),
				ans,
				fmt.Sprintf("Substituting y=%dx gives (1-%d^2)/(1+%d^2) = %s.", m, m, m, ans),
				map[string]string{"kind": "path"})
		default:
			items := []string{
				"\\(x^2+y^2\\)", "\\(xy\\)", "\\(3x+2y\\)", "\\(x^2y/(x^2+y^2)\\) with value 0 at the origin",
				"\\(xy/(x^2+y^2)\\) with value 0 at the origin", "\\((x-y)/(x+y)\\) with value 0 at the origin",
				"\\(x^2/(x^2+y^2)\\) with value 0 at the origin", "\\(|x|/(|x|+|y|)\\) with value 0 at the origin",
			}
			perm := rand.Perm(len(items))[:5]
			count := 0
			parts := make([]string, 0, 5)
			for i, j := range perm {
				if j < 4 {
					count++
				}
				parts = append(parts, fmt.Sprintf("(%c) %s", 'a'+i, items[j]))
			}
			return calcBridgeProblem(o,
				fmt.Sprintf("How many of these functions are continuous at the origin?\n%s\n(enter a number)", strings.Join(parts, "\n")),
				strconv.Itoa(count),
				fmt.Sprintf("Exactly %d of the five functions are continuous at the origin; the others have path-dependent limits.", count),
				map[string]string{"kind": "count"})
		}
	case "moments_3d":
		a := 2 + rand.Intn(4)
		b := 2 + rand.Intn(4)
		switch o {
		case 1:
			c := 2 + rand.Intn(4)
			d := 2 + rand.Intn(4)
			ans := d * a * b * c
			return calcBridgeProblem(o,
				fmt.Sprintf("A box has side lengths %d, %d and %d and constant density %d. What is its mass? (enter a number)", a, b, c, d),
				strconv.Itoa(ans),
				fmt.Sprintf("Mass is density times volume: %d x %d x %d x %d = %d.", d, a, b, c, ans),
				map[string]string{"kind": "mass"})
		case 2:
			c := 2 * (1 + rand.Intn(3))
			d := 2 + rand.Intn(4)
			ans := d * a * b * c * c / 2
			return calcBridgeProblem(o,
				fmt.Sprintf("A box has side lengths %d, %d and %d and constant density %d. What is its first moment with respect to the xy-plane? (enter a number)", a, b, c, d),
				strconv.Itoa(ans),
				fmt.Sprintf("The z-coordinate of the centroid is %d/2, so the moment is mass times that height: %d.", c, ans),
				map[string]string{"kind": "first-moment"})
		default:
			p := 1 + rand.Intn(3)
			q := 1 + rand.Intn(3)
			x := 6 * p
			y := 6 * q
			z := 2 + rand.Intn(4)
			d := 2 + rand.Intn(4)
			ans := d * x * y * z * (x*x + y*y) / 12
			return calcBridgeProblem(o,
				fmt.Sprintf("A box has side lengths %d, %d and %d and constant density %d. What is its moment of inertia about the z-axis? (enter a number)", x, y, z, d),
				strconv.Itoa(ans),
				fmt.Sprintf("The moment is density times volume times (a^2+b^2)/12: %d x %d x %d x %d x (%d^2+%d^2)/12 = %d.", d, x, y, z, x, y, ans),
				map[string]string{"kind": "inertia"})
		}
	default: // normal_line
		a := 1 + rand.Intn(3)
		b := 1 + rand.Intn(3)
		p := rand.Intn(4) - 1
		q := rand.Intn(4) - 1
		switch o {
		case 1:
			ans := 2 * a * p
			return calcBridgeProblem(o,
				fmt.Sprintf("For the surface \\(z=%d x^2+%d y^2\\), what is the x-component of the direction of the normal line at \\((x,y)=(%d,%d)\\)? (enter a number)", a, b, p, q),
				strconv.Itoa(ans),
				fmt.Sprintf("The normal direction is \\(\\langle 2ax,2by,-1\\rangle\\); at x=%d the x-component is %d.", p, ans),
				map[string]string{"kind": "direction-x"})
		case 2:
			ans := q * (1 + 2*b)
			return calcBridgeProblem(o,
				fmt.Sprintf("The normal line to \\(z=%d x^2+%d y^2\\) at \\((x,y)=(%d,%d)\\) is parametrized with \\(t=0\\) at the surface point. What is its y-coordinate at \\(t=1\\)? (enter a number)", a, b, p, q),
				strconv.Itoa(ans),
				fmt.Sprintf("The y-coordinate is %d + (2 x %d x %d) = %d.", q, b, q, ans),
				map[string]string{"kind": "y-at-one"})
		default:
			z0 := a*p*p + b*q*q
			ans := z0 - 1
			return calcBridgeProblem(o,
				fmt.Sprintf("The normal line to \\(z=%d x^2+%d y^2\\) at \\((x,y)=(%d,%d)\\) is parametrized with \\(t=0\\) at the surface point. What is its z-coordinate at \\(t=1\\)? (enter a number)", a, b, p, q),
				strconv.Itoa(ans),
				fmt.Sprintf("The z-coordinate is f(%d,%d) - 1 = %d - 1 = %d.", p, q, z0, ans),
				map[string]string{"kind": "z-at-one"})
		}
	}
}
