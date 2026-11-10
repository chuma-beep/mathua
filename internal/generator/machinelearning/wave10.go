package machinelearning

// Wave 10 machine learning: AdaGrad. The generator uses gradient pairs whose
// cumulative squared gradients are perfect squares, so every update is an exact
// integer and the test can recompute each step from the published facts.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave10(reg *generator.Registry) {
	reg.Register("ml.optim.adagrad", &w10AdaGradGen{})
	registerWave10Rest(reg)
}

func w10P(o int, q, a, ex string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}

func w10I(m map[string]string, k string, v int) {
	m[k] = strconv.Itoa(v)
}

func w10Ord(ctx generator.GeneratorContext) int {
	o := ctx.KPOrdinal
	if o < 1 || o > 3 {
		o = 1 + rand.Intn(3)
	}
	return o
}

type w10AdaGradGen struct{}

func (g *w10AdaGradGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10Ord(ctx)
	pairs := [][2]int{
		{3, 4}, {4, 3}, {6, 8}, {8, 6}, {5, 12}, {12, 5},
		{8, 15}, {15, 8}, {7, 24}, {24, 7}, {9, 12}, {12, 9},
		{9, 40}, {40, 9}, {20, 21}, {21, 20}, {12, 16}, {16, 12},
		{15, 20}, {20, 15},
	}
	pair := pairs[rand.Intn(len(pairs))]
	g1, g2 := pair[0], pair[1]
	G1 := g1 * g1
	G2 := G1 + g2*g2
	root2 := 0
	for root2*root2 < G2 {
		root2++
	}
	if root2*root2 != G2 {
		panic("wave 10 AdaGrad pair does not have a square cumulative gradient")
	}
	k := 1 + rand.Intn(6)
	eta := root2 * k
	u1 := eta * g1 / g1
	u2 := eta * g2 / root2
	w0 := 300 + rand.Intn(301)
	w1 := w0 - u1
	w2 := w1 - u2
	facts := map[string]string{}
	w10I(facts, "g1", g1)
	w10I(facts, "g2", g2)
	w10I(facts, "eta", eta)
	w10I(facts, "G1", G1)
	w10I(facts, "G2", G2)
	w10I(facts, "u1", u1)
	w10I(facts, "u2", u2)
	w10I(facts, "w0", w0)
	w10I(facts, "w1", w1)
	w10I(facts, "w2", w2)
	switch o {
	case 1:
		return w10P(o,
			fmt.Sprintf("AdaGrad starts at w = %d with accumulator 0, learning rate %d and first gradient %d. What is w after the first update?", w0, eta, g1),
			strconv.Itoa(w1),
			fmt.Sprintf("The accumulator is %d and its square root is %d. The first update is %d x %d / %d = %d, so w becomes %d.", G1, g1, eta, g1, g1, u1, w1),
			facts)
	case 2:
		return w10P(o,
			fmt.Sprintf("An AdaGrad run has first gradient %d and second gradient %d. What is the accumulated squared gradient after the second step?", g1, g2),
			strconv.Itoa(G2),
			fmt.Sprintf("Add the new squared gradient: %d + %d^2 = %d.", G1, g2, G2),
			facts)
	default:
		return w10P(o,
			fmt.Sprintf("An AdaGrad run has first update %d and second update %d. How much smaller is the second update than the first?", u1, u2),
			strconv.Itoa(u1-u2),
			fmt.Sprintf("Subtract the updates: %d - %d = %d. The growing accumulator shrinks the second step.", u1, u2, u1-u2),
			facts)
	}
}
