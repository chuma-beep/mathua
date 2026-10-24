package trigonometry

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

type linearEqGen struct{}

func (g *linearEqGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// Only exact rational or radical sine values give an integer degree answer.
	// The ratio is carried as text so the solution walkthrough can cite it.
	type ref struct {
		num, den int
		ratio    string
		answer   int
	}
	// The radicals are absent on purpose: a sine of sqrt(2)/2 cannot be written
	// as an integer equation a sin x = b, so those angles would break the form
	// the solution walkthrough cites.
	refs := []ref{
		{1, 2, "1/2", 30},
		{-1, 2, "-1/2", 210},
		{1, 1, "1", 90},
		{-1, 1, "-1", 270},
		{0, 1, "0", 0},
	}
	r := refs[rand.Intn(len(refs))]
	k := 1 + rand.Intn(4)
	// The equation is (k*den) sin x = k*num, so dividing both sides leaves the
	// exact sine value the answer was built from.
	return generator.Problem{
		Question: fmt.Sprintf("Solve \\(%d\\sin x = %d\\) for the smallest nonnegative solution in degrees.", k*r.den, k*r.num),
		Answer:   strconv.Itoa(r.answer),
		Explanation: fmt.Sprintf("Dividing both sides by %d leaves \\(\\sin x = %s\\). The smallest nonnegative angle with that sine is \\(%d^{\\circ}\\), read from the unit circle.",
			k*r.den, r.ratio, r.answer),
		Facts: map[string]string{
			"a":      strconv.Itoa(k * r.den),
			"b":      strconv.Itoa(k * r.num),
			"ratio":  r.ratio,
			"answer": strconv.Itoa(r.answer),
		},
	}
}

type tripleAngleGen struct{}

func (g *tripleAngleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// A Pythagorean triple from Euclid's formula, scaled so the hypotenuse stays
	// printable. The cosine leg is kept for the facts even though the question
	// asks for the sine identity the solution walkthrough derives.
	m := 2 + rand.Intn(6)
	n := 1 + rand.Intn(m-1)
	scale := 1 + rand.Intn(3)
	opp := scale * (m*m - n*n)
	adj := scale * 2 * m * n
	hyp := scale * (m*m + n*n)

	// sin(3t) = 3 sin t - 4 sin^3 t.
	num := 3*opp*hyp*hyp - 4*opp*opp*opp
	den := hyp * hyp * hyp
	answer := trigFrac(num, den)
	return generator.Problem{
		Question: fmt.Sprintf("If \\(\\sin\\theta = %d/%d\\) and \\(\\theta\\) is acute, what is \\(\\sin(3\\theta)\\)? Give an exact fraction.", opp, hyp),
		Answer:   answer,
		Explanation: fmt.Sprintf("The triple-angle identity is \\(\\sin(3\\theta) = 3\\sin\\theta - 4\\sin^{3}\\theta\\). With \\(\\sin\\theta = %d/%d\\) each term shares the denominator %d, so the result is \\(%s\\).",
			opp, hyp, den, answer),
		Facts: map[string]string{
			"opp":    strconv.Itoa(opp),
			"adj":    strconv.Itoa(adj),
			"hyp":    strconv.Itoa(hyp),
			"answer": answer,
		},
	}
}
