package complex

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

func Register(reg *generator.Registry) {
	reg.Register("complex.concept", &conceptGen{})
	reg.Register("complex.add_sub", &addSubGen{})
	reg.Register("complex.mult", &multGen{})
	reg.Register("complex.conjugate", &conjugateGen{})
	reg.Register("complex.divide", &divideGen{})
	reg.Register("complex.polar", &polarGen{})
	reg.Register("complex.de_moivre", &deMoivreGen{})
	reg.Register("complex.roots", &rootsGen{})
}

func fmtComplex(r, i int) string {
	if i >= 0 {
		return fmt.Sprintf("%d+%di", r, i)
	}
	return fmt.Sprintf("%d%di", r, i)
}

type conceptGen struct{}

func (g *conceptGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(8) + 1
	b := rand.Intn(8) + 1
	if rand.Intn(2) == 0 {
		return generator.Problem{
			Question:    fmt.Sprintf("What is the real part of %s?", fmtComplex(a, b)),
			Answer:      strconv.Itoa(a),
			Explanation: fmt.Sprintf("The real part of %s is %d.", fmtComplex(a, b), a),
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("What is the imaginary part of %s?", fmtComplex(a, b)),
		Answer:      strconv.Itoa(b),
		Explanation: fmt.Sprintf("The imaginary part of %s is %d.", fmtComplex(a, b), b),
	}
}

type addSubGen struct{}

func (g *addSubGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(8) + 1
	b := rand.Intn(8) + 1
	c := rand.Intn(8) + 1
	d := rand.Intn(8) + 1
	if rand.Intn(2) == 0 {
		r := a + c
		i := b + d
		return generator.Problem{
			Question:    fmt.Sprintf("(%s) + (%s) = ?", fmtComplex(a, b), fmtComplex(c, d)),
			Answer:      fmtComplex(r, i),
			Explanation: fmt.Sprintf("(%s) + (%s) = (%d+%d) + (%d+%d)i = %s", fmtComplex(a, b), fmtComplex(c, d), a, c, b, d, fmtComplex(r, i)),
		}
	}
	r := a - c
	i := b - d
	return generator.Problem{
		Question:    fmt.Sprintf("(%s) - (%s) = ?", fmtComplex(a, b), fmtComplex(c, d)),
		Answer:      fmtComplex(r, i),
		Explanation: fmt.Sprintf("(%s) - (%s) = (%d-%d) + (%d-%d)i = %s", fmtComplex(a, b), fmtComplex(c, d), a, c, b, d, fmtComplex(r, i)),
	}
}

type multGen struct{}

func (g *multGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(5) + 1
	b := rand.Intn(5) + 1
	c := rand.Intn(5) + 1
	d := rand.Intn(5) + 1
	r := a*c - b*d
	i := a*d + b*c
	return generator.Problem{
		Question:    fmt.Sprintf("(%s)(%s) = ?", fmtComplex(a, b), fmtComplex(c, d)),
		Answer:      fmtComplex(r, i),
		Explanation: fmt.Sprintf("(%s)(%s) = (%d)(%d) - (%d)(%d) + [(%d)(%d)+(%d)(%d)]i = %s", fmtComplex(a, b), fmtComplex(c, d), a, c, b, d, a, d, b, c, fmtComplex(r, i)),
	}
}

type conjugateGen struct{}

func (g *conjugateGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(8) + 1
	b := rand.Intn(8) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("What is the conjugate of %s?", fmtComplex(a, b)),
		Answer:      fmtComplex(a, -b),
		Explanation: fmt.Sprintf("The conjugate of %s is %s.", fmtComplex(a, b), fmtComplex(a, -b)),
	}
}

type divideGen struct{}

func (g *divideGen) Generate(difficulty float64) generator.Problem {
	r := rand.Intn(4) + 1
	s := rand.Intn(4) + 1
	c := rand.Intn(3) + 1
	d := rand.Intn(3) + 1
	numR := r*c - s*d
	numI := r*d + s*c
	return generator.Problem{
		Question:    fmt.Sprintf("(%s) / (%s) = ?", fmtComplex(numR, numI), fmtComplex(c, d)),
		Answer:      fmtComplex(r, s),
		Explanation: fmt.Sprintf("(%s) / (%s) = %s", fmtComplex(numR, numI), fmtComplex(c, d), fmtComplex(r, s)),
	}
}

type polarGen struct{}

func (g *polarGen) Generate(difficulty float64) generator.Problem {
	if rand.Intn(2) == 0 {
		type angle struct {
			deg     int
			cosReal int
			sinImag int
		}
		angles := []angle{
			{0, 1, 0},
			{90, 0, 1},
			{180, -1, 0},
			{270, 0, -1},
		}
		a := angles[rand.Intn(len(angles))]
		r := rand.Intn(4) + 2
		realPart := r * a.cosReal
		imagPart := r * a.sinImag
		return generator.Problem{
			Question:    fmt.Sprintf("Convert %d(cos %d° + i sin %d°) to rectangular form.", r, a.deg, a.deg),
			Answer:      fmtComplex(realPart, imagPart),
			Explanation: fmt.Sprintf("%d(cos %d° + i sin %d°) = %d(%d) + %d(%d)i = %s", r, a.deg, a.deg, r, a.cosReal, r, a.sinImag, fmtComplex(realPart, imagPart)),
		}
	}

	type triple struct{ a, b, r int }
	triples := []triple{
		{3, 4, 5},
		{4, 3, 5},
		{5, 12, 13},
		{12, 5, 13},
		{6, 8, 10},
		{8, 6, 10},
	}
	t := triples[rand.Intn(len(triples))]
	theta := math.Atan2(float64(t.b), float64(t.a)) * 180 / math.Pi
	thetaStr := strconv.Itoa(int(math.Round(theta)))
	return generator.Problem{
		Question:    fmt.Sprintf("Convert %s to polar form.", fmtComplex(t.a, t.b)),
		Answer:      fmt.Sprintf("%d(cos %s° + i sin %s°)", t.r, thetaStr, thetaStr),
		Explanation: fmt.Sprintf("r = √(a²+b²) = √(%d²+%d²) = %d. θ = arctan(b/a) = arctan(%d/%d) ≈ %s°", t.a, t.b, t.r, t.b, t.a, thetaStr),
	}
}

type deMoivreGen struct{}

func (g *deMoivreGen) Generate(difficulty float64) generator.Problem {
	type stdAngle struct{ deg int }
	angles := []stdAngle{
		{30}, {45}, {60}, {90},
	}
	angle := angles[rand.Intn(len(angles))]
	r := rand.Intn(4) + 2
	n := rand.Intn(2) + 2

	newR := mathutil.IntPow(r, n)
	newDeg := (angle.deg * n) % 360

	var ans string
	switch newDeg {
	case 0:
		ans = fmtComplex(newR, 0)
	case 90:
		ans = fmtComplex(0, newR)
	case 180:
		ans = fmtComplex(-newR, 0)
	case 270:
		ans = fmtComplex(0, -newR)
	default:
		ans = fmt.Sprintf("%d(cos %d° + i sin %d°)", newR, newDeg, newDeg)
	}

	return generator.Problem{
		Question:    fmt.Sprintf("[%d(cos %d° + i sin %d°)]^%d = ?", r, angle.deg, angle.deg, n),
		Answer:      ans,
		Explanation: fmt.Sprintf("[%d(cos %d° + i sin %d°)]^%d = %d(cos %d×%d° + i sin %d×%d°) = %d(cos %d° + i sin %d°) = %s", r, angle.deg, angle.deg, n, r, angle.deg, n, angle.deg, n, newR, newDeg, newDeg, ans),
	}
}

type rootsGen struct{}

func (g *rootsGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(4) + 1
	b := rand.Intn(4) + 1
	c := a*a - b*b
	d := 2 * a * b

	return generator.Problem{
		Question:    fmt.Sprintf("Find the square roots of %s.", fmtComplex(c, d)),
		Answer:      fmt.Sprintf("±%s", fmtComplex(a, b)),
		Explanation: fmt.Sprintf("(%s)² = %s, so the square roots are ±%s.", fmtComplex(a, b), fmtComplex(c, d), fmtComplex(a, b)),
	}
}
