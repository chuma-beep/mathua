package trigonometry

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

// trigBridgeGen serves the concepts added by the 2026 content expansion. Each
// knowledge point asks for a different kind of computation, so the generator
// switches on KPOrdinal rather than drawing one shape for the whole concept.
type trigBridgeGen struct{ kind string }

func trigOrdinal(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return 1 + rand.Intn(3)
}

func trigProb(o int, q, a, e string) generator.Problem {
	return generator.Problem{
		Question:    q,
		Answer:      a,
		Explanation: e,
		Facts:       map[string]string{"answer": a},
		KPOrdinal:   o,
	}
}

func trigDecimal(v float64) string {
	return strconv.FormatFloat(v, 'f', 2, 64)
}

func (g *trigBridgeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := trigOrdinal(ctx)
	triples := [][2]int{{3, 4}, {5, 12}, {8, 15}, {7, 24}, {20, 21}, {9, 40}}
	switch g.kind {
	case "multiple_angle":
		switch o {
		case 1:
			t := triples[rand.Intn(len(triples))]
			a, b := t[0], t[1]
			sin2x := float64(2*a*b) / float64(a*a+b*b)
			return trigProb(o, fmt.Sprintf("An angle x in the first quadrant has sin x = %d/%d and cos x = %d/%d. Use the double-angle identity to find sin 2x. (enter a decimal)", a, a*a+b*b, b, a*a+b*b), trigDecimal(sin2x), fmt.Sprintf("sin 2x = 2 sin x cos x = 2 x %d/%d x %d/%d = %.2f.", a, a*a+b*b, b, a*a+b*b, sin2x))
		case 2:
			b := 2 + rand.Intn(5)
			fn := []string{"sin", "cos"}[rand.Intn(2)]
			return trigProb(o, fmt.Sprintf("How many solutions does %s(%dx) = 0 have in the interval [0, 360) degrees? (enter a number)", fn, b), strconv.Itoa(b), fmt.Sprintf("One full turn contains %d periods of %s(%dx), and each period contributes one zero in the half-open interval, so there are %d solutions.", b, fn, b, b))
		default:
			b := 2 + rand.Intn(5)
			c := 1 + rand.Intn(4)
			return trigProb(o, fmt.Sprintf("How many solutions does sin(%dx) = %d have in the interval [0, 360) degrees? (enter a number)", b, c), strconv.Itoa(b), fmt.Sprintf("The equation sin(u) = %d has two solutions in each full period when |%d| < 1, and there are %d periods in [0, 360), so there are %d solutions.", c, c, b, 2*b))
		}
	case "elevation":
		switch o {
		case 1:
			type scenario struct {
				text string
				ans  int
			}
			pool := []scenario{
				{"A ladder leans against a wall. You know its length and the angle at the ground, and want the height reached. Which ratio gives the height? (enter 1 for sine, 2 for cosine, 3 for tangent)", 1},
				{"A ladder leans against a wall. You know its length and the angle at the ground, and want the distance from the wall to its foot. Which ratio gives that distance? (enter 1 for sine, 2 for cosine, 3 for tangent)", 2},
				{"From the ground you know the horizontal distance to a tower and its height, and want the angle of elevation. Which ratio gives the angle? (enter 1 for sine, 2 for cosine, 3 for tangent)", 3},
				{"A tree casts a shadow and you know the sun's angle of elevation. You want the tree's height from the shadow. Which ratio gives the height? (enter 1 for sine, 2 for cosine, 3 for tangent)", 3},
				{"A kite string makes a known angle with the ground and you know the kite's height. You want the string length. Which ratio gives the length? (enter 1 for sine, 2 for cosine, 3 for tangent)", 1},
				{"A ramp of known length rises at a known angle. You want its horizontal run. Which ratio gives the run? (enter 1 for sine, 2 for cosine, 3 for tangent)", 2},
			}
			s := pool[rand.Intn(len(pool))]
			return trigProb(o, s.text, strconv.Itoa(s.ans), "Use sine for opposite over hypotenuse, cosine for adjacent over hypotenuse, and tangent for opposite over adjacent.")
		case 2:
			d := 3 + rand.Intn(15)
			return trigProb(o, fmt.Sprintf("A tree casts a shadow %d metres long when the angle of elevation of the sun is 45 degrees. How tall is the tree in metres? (enter a number)", d), strconv.Itoa(d), fmt.Sprintf("tan 45 = 1, so the height equals the shadow length, %d metres.", d))
		default:
			d := 1 + rand.Intn(9)
			return trigProb(o, fmt.Sprintf("A tower viewed from %d metres away has an angle of elevation of 60 degrees. What is its height divided by the square root of 3? (enter a number)", d), strconv.Itoa(d), fmt.Sprintf("tan 60 = sqrt(3), so height = %d sqrt(3) and height/sqrt(3) = %d.", d, d))
		}
	case "sinusoidal_model":
		switch o {
		case 1:
			lo := 2 + rand.Intn(20)
			hi := lo + 2*(2+rand.Intn(15))
			mid := (lo + hi) / 2
			return trigProb(o, fmt.Sprintf("A sinusoidal quantity ranges from %d to %d. What is its midline value? (enter a number)", lo, hi), strconv.Itoa(mid), fmt.Sprintf("The midline is the average of the extremes: (%d + %d)/2 = %d.", lo, hi, mid))
		case 2:
			lo := 2 + rand.Intn(20)
			hi := lo + 2*(2+rand.Intn(15))
			amp := (hi - lo) / 2
			return trigProb(o, fmt.Sprintf("A sinusoidal quantity ranges from %d to %d. What is its amplitude? (enter a number)", lo, hi), strconv.Itoa(amp), fmt.Sprintf("The amplitude is half the range: (%d - %d)/2 = %d.", hi, lo, amp))
		default:
			b := []int{1, 2, 3, 4, 5, 6, 8, 9, 10, 12}[rand.Intn(10)]
			return trigProb(o, fmt.Sprintf("A sinusoidal model is written y = A sin(%dx) + D with x in degrees. What is its period in degrees? (enter a number)", b), strconv.Itoa(360/b), fmt.Sprintf("The period is 360/b = 360/%d = %d degrees.", b, 360/b))
		}
	case "inverse_composition":
		switch o {
		case 1:
			angles := []int{0, 30, 45, 60, 90}
			a := angles[rand.Intn(len(angles))]
			return trigProb(o, fmt.Sprintf("What is arcsin(sin %d degrees), evaluated in the principal range of arcsin? (enter a number of degrees)", a), strconv.Itoa(a), fmt.Sprintf("For an angle in the principal range, arcsin(sin %d) = %d.", a, a))
		case 2:
			angles := []int{0, 30, 45, 60, 90, 120, 135, 150, 180}
			a := angles[rand.Intn(len(angles))]
			return trigProb(o, fmt.Sprintf("What is arccos(cos %d degrees), evaluated in the principal range of arccos? (enter a number of degrees)", a), strconv.Itoa(a), fmt.Sprintf("For an angle in the principal range, arccos(cos %d) = %d.", a, a))
		default:
			pool := []struct {
				expr string
				yes  bool
			}{{"sin(arcsin x) = x for x in [-1, 1]", true}, {"arcsin(sin 200 degrees) = 200 degrees", false}, {"cos(arccos x) = x for x in [-1, 1]", true}, {"arccos(cos 60 degrees) = 60 degrees", true}, {"tan(arctan x) = x for every real x", true}, {"arcsin(sin 150 degrees) = 150 degrees", false}, {"arccos(cos 240 degrees) = 240 degrees", false}}
			perm := rand.Perm(len(pool))[:5]
			count := 0
			var parts []string
			for i, j := range perm {
				parts = append(parts, fmt.Sprintf("(%s) %s", string(rune('a'+i)), pool[j].expr))
				if pool[j].yes {
					count++
				}
			}
			return trigProb(o, fmt.Sprintf("How many of these statements are true? %s (enter a number)", joinSpace(parts)), strconv.Itoa(count), fmt.Sprintf("%d of the five statements are true. An inverse composition returns the original value only when that value lies in the principal range of the inverse function.", count))
		}
	case "vector_components":
		switch o {
		case 1:
			r := 2 + rand.Intn(19)
			angle := []int{0, 90, 180, 270}[rand.Intn(4)]
			x := 0
			switch angle {
			case 0:
				x = r
			case 180:
				x = -r
			}
			return trigProb(o, fmt.Sprintf("A vector of magnitude %d points at %d degrees from the positive x-axis. What is its x-component? (enter a number)", r, angle), strconv.Itoa(x), fmt.Sprintf("The x-component is r cos(theta) = %d cos %d = %d.", r, angle, x))
		case 2:
			r := 2 + rand.Intn(19)
			angle := []int{0, 90, 180, 270}[rand.Intn(4)]
			y := 0
			switch angle {
			case 90:
				y = r
			case 270:
				y = -r
			}
			return trigProb(o, fmt.Sprintf("A vector of magnitude %d points at %d degrees from the positive x-axis. What is its y-component? (enter a number)", r, angle), strconv.Itoa(y), fmt.Sprintf("The y-component is r sin(theta) = %d sin %d = %d.", r, angle, y))
		default:
			t := triples[rand.Intn(len(triples))]
			return trigProb(o, fmt.Sprintf("A vector has components %d and %d. What is its magnitude? (enter a number)", t[0], t[1]), strconv.Itoa(t[0]*t[0]+t[1]*t[1]), "The magnitude squared is the sum of the squares; here it is a Pythagorean triple.")
		}
	case "even_odd":
		switch o {
		case 1:
			fns := []struct {
				name string
				odd  bool
			}{{"sin", true}, {"cos", false}, {"tan", true}, {"sec", false}, {"csc", true}, {"cot", true}}
			f := fns[rand.Intn(len(fns))]
			k := 1
			word := "even"
			if f.odd {
				k = -1
				word = "odd"
			}
			return trigProb(o, fmt.Sprintf("The function %s is %s, so %s(-x) = k %s(x) for every x in its domain. What is k? (enter -1 or 1)", f.name, word, f.name, f.name), strconv.Itoa(k), fmt.Sprintf("%s is %s, so the coefficient is %d.", f.name, word, k))
		case 2:
			pool := []struct {
				name string
				even bool
			}{{"cos", true}, {"sec", true}, {"sin", false}, {"tan", false}, {"csc", false}, {"cot", false}}
			perm := rand.Perm(len(pool))[:5]
			count := 0
			var names []string
			for _, i := range perm {
				names = append(names, pool[i].name)
				if pool[i].even {
					count++
				}
			}
			return trigProb(o, fmt.Sprintf("How many of these functions are even: %s? (enter a number)", joinComma(names)), strconv.Itoa(count), fmt.Sprintf("%d of the five functions are even. Cosine and secant are even; the other four are odd.", count))
		default:
			p := 1 + rand.Intn(9)
			fn := []string{"sin", "tan", "csc", "cot"}[rand.Intn(4)]
			return trigProb(o, fmt.Sprintf("Suppose %s(a) = %d/10. What is %s(-a)? (enter a decimal)", fn, p, fn), trigDecimal(-float64(p)/10), fmt.Sprintf("%s is odd, so %s(-a) = -%s(a) = -%.1f.", fn, fn, fn, float64(p)/10))
		}
	case "frequency":
		switch o {
		case 1:
			b := []int{1, 2, 3, 4, 5, 6, 8, 9, 10, 12, 15, 18, 20, 24, 30}[rand.Intn(15)]
			return trigProb(o, fmt.Sprintf("What is the period in degrees of y = sin(%dx)? (enter a number)", b), strconv.Itoa(360/b), fmt.Sprintf("The period is 360/%d = %d degrees.", b, 360/b))
		case 2:
			p := []int{360, 180, 120, 90, 72, 60, 45, 40, 36, 30, 24, 20, 18, 15, 12}[rand.Intn(15)]
			return trigProb(o, fmt.Sprintf("A sine graph completes one cycle every %d degrees. What is the coefficient b in y = sin(bx)? (enter a number)", p), strconv.Itoa(360/p), fmt.Sprintf("b = 360/period = 360/%d = %d.", p, 360/p))
		default:
			b := []int{1, 2, 3, 4, 5, 6, 8, 9, 10, 12}[rand.Intn(10)]
			return trigProb(o, fmt.Sprintf("How many complete cycles does y = cos(%dx) make over 360 degrees? (enter a number)", b), strconv.Itoa(b), fmt.Sprintf("The graph completes %d cycles in 360 degrees because the coefficient is %d.", b, b))
		}
	case "using_identities":
		switch o {
		case 1:
			t := triples[rand.Intn(len(triples))]
			a, b := t[0], t[1]
			return trigProb(o, fmt.Sprintf("An angle x in the first quadrant has sin x = %d/%d. What is cos x? (enter a decimal)", a, a*a+b*b), trigDecimal(float64(b)/float64(a*a+b*b)), fmt.Sprintf("cos x = sqrt(1 - sin^2 x) = %d/%d = %.2f.", b, a*a+b*b, float64(b)/float64(a*a+b*b)))
		case 2:
			t := triples[rand.Intn(len(triples))]
			a, b := t[0], t[1]
			return trigProb(o, fmt.Sprintf("An angle x in the first quadrant has sin x = %d/%d and cos x = %d/%d. What is tan x? (enter a decimal)", a, a*a+b*b, b, a*a+b*b), trigDecimal(float64(a)/float64(b)), fmt.Sprintf("tan x = sin x / cos x = %d/%d = %.2f.", a, b, float64(a)/float64(b)))
		default:
			t := triples[rand.Intn(len(triples))]
			a := t[0]
			s := float64(a) / float64(a*a+t[1]*t[1])
			v := 1 - 2*s*s
			return trigProb(o, fmt.Sprintf("An angle x in the first quadrant has sin x = %d/%d. Use cos 2x = 1 - 2 sin^2 x to find cos 2x. (enter a decimal)", a, a*a+t[1]*t[1]), trigDecimal(v), fmt.Sprintf("cos 2x = 1 - 2(%d/%d)^2 = %.2f.", a, a*a+t[1]*t[1], v))
		}
	case "harmonic_motion":
		switch o {
		case 1:
			a := 2 + rand.Intn(15)
			return trigProb(o, fmt.Sprintf("A particle moves as y = %d sin(3t). What is its amplitude? (enter a number)", a), strconv.Itoa(a), fmt.Sprintf("The amplitude is the coefficient %d.", a))
		case 2:
			b := []int{1, 2, 3, 4, 5, 6, 8, 9, 10, 12}[rand.Intn(10)]
			return trigProb(o, fmt.Sprintf("A particle moves as y = 5 sin(%dt) with t in degrees. What is its period? (enter a number)", b), strconv.Itoa(360/b), fmt.Sprintf("The period is 360/%d = %d.", b, 360/b))
		default:
			a := 2 + rand.Intn(15)
			d := 1 + rand.Intn(20)
			return trigProb(o, fmt.Sprintf("A particle moves as y = %d + %d sin(2t). What is its maximum displacement? (enter a number)", d, a), strconv.Itoa(d+a), fmt.Sprintf("The maximum is the offset plus the amplitude: %d + %d = %d.", d, a, d+a))
		}
	case "segment_area":
		switch o {
		case 1:
			r := 2 * (1 + rand.Intn(9))
			angle := []int{90, 180, 270}[rand.Intn(3)]
			ans := angle * r * r / 360
			return trigProb(o, fmt.Sprintf("A sector of a circle of radius %d has central angle %d degrees. What is its area divided by pi? (enter a number)", r, angle), strconv.Itoa(ans), fmt.Sprintf("Area is (theta/360) pi r^2, so the coefficient is %d x %d^2 / 360 = %d.", angle, r, ans))
		case 2:
			r := 2 * (1 + rand.Intn(9))
			angle := []int{90, 180, 270}[rand.Intn(3)]
			ans := r * r / 2
			return trigProb(o, fmt.Sprintf("A triangle has two sides of length %d meeting at %d degrees. What is its area? (enter a number)", r, angle), strconv.Itoa(ans), fmt.Sprintf("Area is (1/2) r^2 sin(theta); here sin is 0, 1 or -1, and the magnitude gives %d.", ans))
		default:
			r := 2 * (1 + rand.Intn(9))
			ans := r * r / 2
			return trigProb(o, fmt.Sprintf("A segment is cut off by a diameter of a circle of radius %d. What is the segment's area divided by pi? (enter a number)", r), strconv.Itoa(ans), fmt.Sprintf("A diameter cuts a semicircle, whose area is (1/2) pi r^2, so the coefficient is %d^2/2 = %d.", r, ans))
		}
	}
	return trigProb(o, "What is 1 + 1? (enter a number)", "2", "One plus one is two.")
}

func joinSpace(parts []string) string {
	return strings.Join(parts, " ")
}

func joinComma(parts []string) string {
	if len(parts) < 2 {
		return strings.Join(parts, "")
	}
	return strings.Join(parts[:len(parts)-1], ", ") + " and " + parts[len(parts)-1]
}
