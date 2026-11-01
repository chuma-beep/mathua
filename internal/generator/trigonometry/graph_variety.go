package trigonometry

import (
	"fmt"
	"math"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

// The four graph concepts asked three fixed questions each, in prose, with the
// answer written into the option list. Period, asymptotes, amplitude and range
// are all functions of the coefficients the graph is written with, so the shapes
// below draw those coefficients and compute the answer instead of reciting it.
// Each pool keeps one definitional question, because a learner still has to know
// what a reciprocal is.

// degStr renders an angle in degrees without a trailing .0.
func degStr(v float64) string {
	return strconv.FormatFloat(math.Round(v*100)/100, 'f', -1, 64)
}

// distinctDistractors returns three options that differ from the answer and from
// each other.
//
// trigPickMC hands the candidates to generator.MultipleChoice, which caps the
// pool at four letters and drops duplicates. A claim whose candidates collide,
// for example an amplitude equal to the period coefficient, would ship three
// options, and the batch test requires four.
func distinctDistractors(correct string, candidates ...string) []string {
	seen := map[string]bool{correct: true}
	out := make([]string, 0, 3)
	consider := func(candidate string) {
		if len(out) == 3 || seen[candidate] {
			return
		}
		seen[candidate] = true
		out = append(out, candidate)
	}
	for _, candidate := range candidates {
		consider(candidate)
	}
	// These are wrong for every shape here, so they restore a fourth option
	// when the collisions above remove one, and a fifth when they remove two.
	for _, fallback := range []string{"0", "180", "360", "-1", "[0, 0]"} {
		consider(fallback)
	}
	return out
}

type tanGraphGen struct{}

func (g *tanGraphGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return tangentGraphProblem()
}

// sineGraph renders a sin(b(x - c)) + d the way a textbook writes it: no
// coefficient of one, no "x - -1" when the shift is negative, and no "+ 0".
func sineGraph(a, b, c, d int) string {
	inner := "x"
	switch {
	case c > 0:
		inner = fmt.Sprintf("x - %d", c)
	case c < 0:
		inner = fmt.Sprintf("x + %d", -c)
	}
	arg := inner
	if b != 1 {
		arg = fmt.Sprintf("%d(%s)", b, inner)
	}
	body := fmt.Sprintf("sin(%s)", arg)
	if a != 1 {
		body = fmt.Sprintf("%d %s", a, body)
	}
	switch {
	case d > 0:
		body = fmt.Sprintf("%s + %d", body, d)
	case d < 0:
		body = fmt.Sprintf("%s - %d", body, -d)
	}
	return "y = " + body
}

type graphTransformGen struct{}

func (g *graphTransformGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return graphTransformationProblem()
}

type reciprocalGraphGen struct{}

func (g *reciprocalGraphGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return reciprocalGraphProblem()
}

type inverseGraphGen struct{}

func (g *inverseGraphGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return inverseGraphProblem()
}
