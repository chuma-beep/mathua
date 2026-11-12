package trigonometry

import (
	"math"
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func TestWave9TrigAnswers(t *testing.T) {
	for _, kind := range w9TestKinds {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 40; draw++ {
				p := (&w9TrigGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				if p.KPOrdinal != kp {
					t.Fatalf("%s kp%d: ordinal %d", kind, kp, p.KPOrdinal)
				}
				want := w9TestAnswer(t, kind, kp, p)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer %s, want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func TestWave9TrigDiversity(t *testing.T) {
	for _, kind := range w9TestKinds {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for draw := 0; draw < 200; draw++ {
				p := (&w9TrigGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				seen[p.Answer] = true
			}
			if len(seen) < 5 {
				t.Errorf("%s kp%d: only %d distinct answers in 200 draws: %v", kind, kp, len(seen), seen)
			}
		}
	}
}

var w9TestKinds = []string{
	"factoring", "reciprocal", "mixed_angles", "extraneous", "absolute",
	"higher_derivatives", "taylor", "power_integrals", "damping", "envelope",
	"beats", "catenary", "gudermannian", "hyperbolic_addition",
	"hyperbolic_integral", "spherical_right", "spherical_polar",
	"spherical_polygon", "gradian", "signed", "conditional", "tangent_triple",
}

func w9TestInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func w9TestDegrees(rad float64) int { return int(math.Round(rad * 180 / math.Pi)) }

func w9TestAnswer(t *testing.T, kind string, kp int, p generator.Problem) int {
	t.Helper()
	switch kind {
	case "factoring":
		n := w9TestInt(t, p, "N")
		switch kp {
		case 1:
			return 4 * n
		case 2:
			return 2 * n
		default:
			return n
		}
	case "reciprocal":
		if kp == 1 {
			return w9TestInt(t, p, "a")
		}
		return w9TestInt(t, p, "b")
	case "mixed_angles":
		m, n := w9TestInt(t, p, "m"), w9TestInt(t, p, "N")
		if kp == 3 {
			return m * n
		}
		return 2 * m * n
	case "extraneous":
		return 2 * w9TestInt(t, p, "N")
	case "absolute":
		n := w9TestInt(t, p, "N")
		if kp == 1 {
			return 4 * n
		}
		return 2 * n
	case "higher_derivatives":
		a := w9TestInt(t, p, "a")
		switch kp {
		case 1:
			return -a * a * a
		case 2:
			return a * a * a * a
		default:
			return 2 * a
		}
	case "taylor":
		if kp == 3 {
			a := w9TestInt(t, p, "a")
			return a * a * a
		}
		n := w9TestInt(t, p, "n")
		if kp == 1 {
			return w9Fact(2*n + 1)
		}
		return w9Fact(2 * n)
	case "power_integrals":
		n := w9TestInt(t, p, "n")
		if kp == 2 {
			return w9DoubleFact(n)
		}
		return w9DoubleFact(n - 1)
	case "damping":
		switch kp {
		case 1:
			return w9TestInt(t, p, "m")
		case 2:
			return 1 << w9TestInt(t, p, "m")
		default:
			return w9TestInt(t, p, "k") * w9TestInt(t, p, "t")
		}
	case "envelope":
		a, b := w9TestInt(t, p, "a"), w9TestInt(t, p, "b")
		switch kp {
		case 1:
			return a + b
		case 2:
			return a - b
		default:
			return 2 * b
		}
	case "beats":
		switch kp {
		case 1:
			return w9TestInt(t, p, "f2") - w9TestInt(t, p, "f1")
		case 2:
			return w9TestInt(t, p, "f") * w9TestInt(t, p, "m")
		default:
			return 10 * (w9TestInt(t, p, "f2") - w9TestInt(t, p, "f1"))
		}
	case "catenary":
		a := float64(w9TestInt(t, p, "a"))
		switch kp {
		case 1:
			return int(math.Round(a * math.Cosh(1)))
		case 2:
			return int(math.Round(a * math.Cosh(0.5)))
		default:
			return int(math.Round(2 * a * math.Sinh(1)))
		}
	case "gudermannian":
		x := float64(w9TestInt(t, p, "x"))
		switch kp {
		case 1:
			return int(math.Round(math.Sinh(x)))
		case 2:
			return w9TestDegrees(2*math.Atan(math.Exp(x)) - math.Pi/2)
		default:
			return w9TestInt(t, p, "x")
		}
	case "hyperbolic_addition":
		if kp == 3 {
			s1, c1 := float64(w9TestInt(t, p, "s1")), float64(w9TestInt(t, p, "c1"))
			s2, c2 := float64(w9TestInt(t, p, "s2")), float64(w9TestInt(t, p, "c2"))
			return int(s1*c2 + c1*s2)
		}
		s, c := float64(w9TestInt(t, p, "s")), float64(w9TestInt(t, p, "c"))
		if kp == 1 {
			return int(2 * s * c)
		}
		return int(c*c + s*s)
	case "hyperbolic_integral":
		a := float64(w9TestInt(t, p, "a"))
		switch kp {
		case 1:
			return int(math.Round(math.Sinh(a)))
		case 2:
			return int(math.Round(math.Cosh(a) - 1))
		default:
			return w9TestInt(t, p, "a")
		}
	case "spherical_right":
		switch kp {
		case 1:
			return w9TestInt(t, p, "A") + w9TestInt(t, p, "B")
		case 2:
			a := float64(w9TestInt(t, p, "a")) * math.Pi / 180
			b := float64(w9TestInt(t, p, "b")) * math.Pi / 180
			return w9TestDegrees(math.Acos(math.Cos(a) * math.Cos(b)))
		default:
			a := float64(w9TestInt(t, p, "a")) * math.Pi / 180
			b := float64(w9TestInt(t, p, "B")) * math.Pi / 180
			return w9TestDegrees(math.Acos(math.Cos(a) * math.Sin(b)))
		}
	case "spherical_polar":
		switch kp {
		case 1:
			return 180 - w9TestInt(t, p, "a")
		case 2:
			return 180 - w9TestInt(t, p, "A")
		default:
			return 540 - (w9TestInt(t, p, "A") + w9TestInt(t, p, "B") + w9TestInt(t, p, "C"))
		}
	case "spherical_polygon":
		n := w9TestInt(t, p, "n")
		return n*w9TestInt(t, p, "A") - (n-2)*180
	case "gradian":
		k := w9TestInt(t, p, "k")
		switch kp {
		case 1:
			return 9 * k
		case 2:
			return 10 * k
		default:
			return 10 * (45 * k) / 9
		}
	case "signed":
		k := w9TestInt(t, p, "k")
		switch kp {
		case 1:
			return 360 - k
		case 2:
			return k
		default:
			return (360 - k%360) % 360
		}
	case "conditional":
		switch kp {
		case 1:
			return 90 - w9TestInt(t, p, "A")
		case 2:
			return 180 - w9TestInt(t, p, "A") - w9TestInt(t, p, "B")
		default:
			return 90 - w9TestInt(t, p, "B")
		}
	case "tangent_triple":
		n := w9TestInt(t, p, "n")
		switch kp {
		case 1:
			return 3*n*n - 1
		case 2:
			return n * (n*n - 3)
		default:
			return 3*n - n*n*n
		}
	}
	return 0
}
