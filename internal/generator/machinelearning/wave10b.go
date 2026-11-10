package machinelearning

// Wave 10 machine learning, part two. Every branch publishes a, b, c and answer.

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave10Rest(reg *generator.Registry) {
	reg.Register("ml.optim.rmsprop", &w10gRMSPropGen{})
	reg.Register("ml.optim.adamw", &w10gAdamWGen{})
	reg.Register("ml.ensemble.stacking", &w10gStackingGen{})
	reg.Register("ml.ensemble.voting", &w10gVotingGen{})
	reg.Register("ml.unsupervised.spectral", &w10gSpectralGen{})
	reg.Register("ml.loss.hinge", &w10gHingeGen{})
	reg.Register("ml.loss.huber", &w10gHuberGen{})
}

func w10gP(o int, q, a, ex string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}

func w10gI(m map[string]string, k string, v int) { m[k] = strconv.Itoa(v) }

func w10gF(v float64) string { return strconv.FormatFloat(v, 'f', 3, 64) }

type w10gRMSPropGen struct{}

func (g *w10gRMSPropGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10Ord(ctx)
	g1 := 2 * (1 + rand.Intn(4))
	g2 := 2 * (1 + rand.Intn(4))
	eta := 1 + rand.Intn(5)
	v1 := 0.5 * float64(g1*g1)
	v2 := 0.5*float64(g2*g2) + 0.5*v1
	u1 := float64(eta*g1) / math.Sqrt(v1)
	u2 := float64(eta*g2) / math.Sqrt(v2)
	facts := map[string]string{}
	w10gI(facts, "a", g1)
	w10gI(facts, "b", g2)
	w10gI(facts, "c", eta)
	facts["v1"] = w10gF(v1)
	facts["v2"] = w10gF(v2)
	facts["u1"] = w10gF(u1)
	facts["u2"] = w10gF(u2)
	switch o {
	case 1:
		return w10gP(o, fmt.Sprintf("RMSProp has learning rate %d, decay 0.5, zero initial accumulator and first gradient %d. What is the first update size?", eta, g1), w10gF(u1), fmt.Sprintf("The accumulator is %s and its square root is %.3f. The update is %d x %d / %.3f = %s.", w10gF(v1), math.Sqrt(v1), eta, g1, math.Sqrt(v1), w10gF(u1)), facts)
	case 2:
		return w10gP(o, fmt.Sprintf("An RMSProp run has decay 0.5, first gradient %d and second gradient %d. What is the accumulator after the second step?", g1, g2), w10gF(v2), fmt.Sprintf("Update the decayed average: 0.5 x %d^2 + 0.5 x %s = %s.", g2, w10gF(v1), w10gF(v2)), facts)
	default:
		return w10gP(o, fmt.Sprintf("An RMSProp run has first update %s and second update %s. How much smaller is the second update?", w10gF(u1), w10gF(u2)), w10gF(u1-u2), fmt.Sprintf("Subtract: %s - %s = %s.", w10gF(u1), w10gF(u2), w10gF(u1-u2)), facts)
	}
}

type w10gAdamWGen struct{}

func (g *w10gAdamWGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10Ord(ctx)
	w0 := 100 * (1 + rand.Intn(3))
	eta := 1 + rand.Intn(5)
	lambda := 1 + rand.Intn(10)
	grad := 1 + rand.Intn(5)
	decay := w0 * eta * lambda / 100
	adaptive := eta * grad
	newW := w0 - adaptive - decay
	facts := map[string]string{}
	w10gI(facts, "a", w0)
	w10gI(facts, "b", eta)
	w10gI(facts, "c", lambda)
	w10gI(facts, "decay", decay)
	w10gI(facts, "adaptive", adaptive)
	w10gI(facts, "new", newW)
	w10gI(facts, "grad", grad)
	switch o {
	case 1:
		return w10gP(o, fmt.Sprintf("AdamW has weight w = %d, learning rate %d and weight decay %d%%. What is the decoupled weight-decay term?", w0, eta, lambda), strconv.Itoa(decay), fmt.Sprintf("Multiply w by the learning rate and decay rate: %d x %d x %d/100 = %d.", w0, eta, lambda, decay), facts)
	case 2:
		return w10gP(o, fmt.Sprintf("An AdamW step has learning rate %d and adaptive gradient term %d. What is the Adam adaptive update before weight decay?", eta, grad), strconv.Itoa(adaptive), fmt.Sprintf("Multiply the learning rate by the adaptive term: %d x %d = %d.", eta, grad, adaptive), facts)
	default:
		return w10gP(o, fmt.Sprintf("AdamW starts at w = %d with learning rate %d, adaptive term %d and weight decay %d%%. What is the new weight?", w0, eta, grad, lambda), strconv.Itoa(newW), fmt.Sprintf("Subtract the adaptive step %d and the decay term %d: %d - %d - %d = %d.", adaptive, decay, w0, adaptive, decay, newW), facts)
	}
}

type w10gStackingGen struct{}

func (g *w10gStackingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10Ord(ctx)
	p1, p2 := rand.Intn(11), rand.Intn(11)
	w1, w2 := 1+rand.Intn(5), 1+rand.Intn(5)
	b := rand.Intn(6)
	meta := w1*p1 + w2*p2 + b
	target := rand.Intn(11)
	e1 := p1 - target
	if e1 < 0 {
		e1 = -e1
	}
	e2 := p2 - target
	if e2 < 0 {
		e2 = -e2
	}
	if e1 == e2 {
		target++
		e1 = p1 - target
		if e1 < 0 {
			e1 = -e1
		}
		e2 = p2 - target
		if e2 < 0 {
			e2 = -e2
		}
	}
	facts := map[string]string{}
	w10gI(facts, "a", p1)
	w10gI(facts, "b", p2)
	w10gI(facts, "c", meta)
	facts["w1"] = strconv.Itoa(w1)
	facts["w2"] = strconv.Itoa(w2)
	facts["bias"] = strconv.Itoa(b)
	facts["target"] = strconv.Itoa(target)
	switch o {
	case 1:
		return w10gP(o, fmt.Sprintf("A stacking ensemble has base predictions %d and %d, meta-weights %d and %d and bias %d. What is the meta-feature?", p1, p2, w1, w2, b), strconv.Itoa(meta), fmt.Sprintf("Combine the weighted predictions: %d x %d + %d x %d + %d = %d.", w1, p1, w2, p2, b, meta), facts)
	case 2:
		return w10gP(o, fmt.Sprintf("A stacking ensemble has base predictions %d and %d, meta-weights %d and %d and bias %d. What is the meta-prediction minus the first base prediction?", p1, p2, w1, w2, b), strconv.Itoa(meta-p1), fmt.Sprintf("The meta-prediction is %d, so the difference is %d - %d = %d.", meta, meta, p1, meta-p1), facts)
	default:
		e1 := p1 - target
		if e1 < 0 {
			e1 = -e1
		}
		e2 := p2 - target
		if e2 < 0 {
			e2 = -e2
		}
		better := 1
		if e2 < e1 {
			better = 2
		}
		return w10gP(o, fmt.Sprintf("With target %d, base model 1 predicts %d and base model 2 predicts %d. Which model has the smaller absolute error (1 or 2)?", target, p1, p2), strconv.Itoa(better), fmt.Sprintf("The absolute errors are %d and %d, so model %d is closer.", e1, e2, better), facts)
	}
}

type w10gVotingGen struct{}

func (g *w10gVotingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10Ord(ctx)
	votes := make([]int, 5)
	ones := 0
	for i := range votes {
		votes[i] = rand.Intn(2)
		ones += votes[i]
	}
	if ones == 0 || ones == 5 {
		ones = 3
		votes = []int{1, 1, 1, 0, 0}
	}
	zeros := 5 - ones
	majority := 1
	if zeros > ones {
		majority = 0
	}
	margin := ones - zeros
	if margin < 0 {
		margin = -margin
	}
	facts := map[string]string{}
	w10gI(facts, "a", ones)
	w10gI(facts, "b", zeros)
	w10gI(facts, "c", margin)
	switch o {
	case 1:
		return w10gP(o, fmt.Sprintf("Five classifiers vote %v. How many votes does label 1 receive?", votes), strconv.Itoa(ones), fmt.Sprintf("Count the votes equal to 1: %d.", ones), facts)
	case 2:
		return w10gP(o, fmt.Sprintf("Five classifiers vote %v. What is the majority label (0 or 1)?", votes), strconv.Itoa(majority), fmt.Sprintf("Label 1 has %d votes and label 0 has %d, so the majority is %d.", ones, zeros, majority), facts)
	default:
		hi, lo := ones, zeros
		if lo > hi {
			hi, lo = lo, hi
		}
		return w10gP(o, fmt.Sprintf("Five classifiers vote %v. What is the winning margin in votes?", votes), strconv.Itoa(margin), fmt.Sprintf("The margin is %d - %d = %d.", hi, lo, margin), facts)
	}
}

type w10gSpectralGen struct{}

func (g *w10gSpectralGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10Ord(ctx)
	vals := make([]int, 4)
	for i := range vals {
		v := 1 + rand.Intn(4)
		if rand.Intn(2) == 0 {
			v = -v
		}
		vals[i] = v
	}
	if vals[0] < 0 && vals[1] < 0 && vals[2] < 0 && vals[3] < 0 {
		vals[0] = -vals[0]
	}
	if vals[0] > 0 && vals[1] > 0 && vals[2] > 0 && vals[3] > 0 {
		vals[0] = -vals[0]
	}
	// Path 0-1-2-3; a sign split crosses the middle edge when the first two are negative.
	negative := 0
	for _, v := range vals {
		if v < 0 {
			negative++
		}
	}
	cross := 0
	for i := 0; i < len(vals)-1; i++ {
		if (vals[i] < 0) != (vals[i+1] < 0) {
			cross++
		}
	}
	node := rand.Intn(4)
	cluster := 1
	if vals[node] > 0 {
		cluster = 2
	}
	facts := map[string]string{}
	w10gI(facts, "a", negative)
	w10gI(facts, "b", cross)
	w10gI(facts, "c", cluster)
	facts["vals"] = fmt.Sprintf("%d,%d,%d,%d", vals[0], vals[1], vals[2], vals[3])
	w10gI(facts, "node", node)
	switch o {
	case 1:
		return w10gP(o, fmt.Sprintf("A Fiedler vector is (%d, %d, %d, %d). How many entries are negative?", vals[0], vals[1], vals[2], vals[3]), strconv.Itoa(negative), fmt.Sprintf("Count the negative entries: %d.", negative), facts)
	case 2:
		return w10gP(o, fmt.Sprintf("A Fiedler vector is (%d, %d, %d, %d). Using the sign, which cluster contains node %d (1 for negative, 2 for positive)?", vals[0], vals[1], vals[2], vals[3], node), strconv.Itoa(cluster), fmt.Sprintf("The entry is %d, so the node is in cluster %d.", vals[node], cluster), facts)
	default:
		return w10gP(o, fmt.Sprintf("A Fiedler vector is (%d, %d, %d, %d) on the path 0-1-2-3. How many path edges cross the sign split?", vals[0], vals[1], vals[2], vals[3]), strconv.Itoa(cross), fmt.Sprintf("Check adjacent signs; %d edge(s) cross.", cross), facts)
	}
}

type w10gHingeGen struct{}

func (g *w10gHingeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10Ord(ctx)
	y := 1
	if rand.Intn(2) == 0 {
		y = -1
	}
	f := -2 + rand.Intn(6)
	margin := y * f
	loss := 1 - margin
	if loss < 0 {
		loss = 0
	}
	sub := 0
	if margin < 1 {
		sub = -y
	}
	facts := map[string]string{}
	w10gI(facts, "a", y)
	w10gI(facts, "b", f)
	w10gI(facts, "c", margin)
	w10gI(facts, "loss", loss)
	w10gI(facts, "sub", sub)
	switch o {
	case 1:
		return w10gP(o, fmt.Sprintf("A hinge-loss example has label y = %d and score f = %d. What is the hinge loss max(0, 1 - yf)?", y, f), strconv.Itoa(loss), fmt.Sprintf("The margin is %d, so the loss is max(0, 1 - (%d)) = %d.", margin, margin, loss), facts)
	case 2:
		return w10gP(o, fmt.Sprintf("A classifier has label y = %d and score f = %d. What is the margin yf?", y, f), strconv.Itoa(margin), fmt.Sprintf("Multiply label by score: (%d)(%d) = %d.", y, f, margin), facts)
	default:
		return w10gP(o, fmt.Sprintf("A hinge-loss example has label y = %d and score f = %d. What is the subgradient with respect to f?", y, f), strconv.Itoa(sub), fmt.Sprintf("The margin is %d. The subgradient is -y = %d when the margin is below 1 and 0 otherwise, so here it is %d.", margin, -y, sub), facts)
	}
}

type w10gHuberGen struct{}

func (g *w10gHuberGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10Ord(ctx)
	delta := 6
	r := -12 + rand.Intn(25)
	if o == 1 {
		r = -delta + rand.Intn(2*delta+1)
	}
	if o == 2 {
		for {
			r = -12 + rand.Intn(25)
			if r > delta || r < -delta {
				break
			}
		}
	}
	abs := r
	if abs < 0 {
		abs = -abs
	}
	loss := 0.0
	if abs <= delta {
		loss = float64(r*r) / 2
	} else {
		loss = float64(delta) * (float64(abs) - float64(delta)/2)
	}
	mse := float64(r*r) / 2
	diff := mse - loss
	facts := map[string]string{}
	w10gI(facts, "a", delta)
	w10gI(facts, "b", r)
	facts["c"] = w10gF(loss)
	facts["mse"] = w10gF(mse)
	facts["diff"] = w10gF(diff)
	switch o {
	case 1:
		return w10gP(o, fmt.Sprintf("A Huber loss has threshold delta = %d and residual %d, which is within the threshold. What is the loss?", delta, r), w10gF(loss), fmt.Sprintf("Below the threshold the loss is 0.5r^2 = 0.5 x %d^2 = %s.", r, w10gF(loss)), facts)
	case 2:
		return w10gP(o, fmt.Sprintf("A Huber loss has threshold delta = %d and residual %d, which is beyond the threshold. What is the loss?", delta, r), w10gF(loss), fmt.Sprintf("Above the threshold the loss is delta(|r| - delta/2) = %d(%d - %d) = %s.", delta, abs, delta/2, w10gF(loss)), facts)
	default:
		return w10gP(o, fmt.Sprintf("For residual %d, how much larger is the squared-error loss than the Huber loss with delta = %d?", r, delta), w10gF(diff), fmt.Sprintf("The squared loss is %s and the Huber loss is %s, so the difference is %s.", w10gF(mse), w10gF(loss), w10gF(diff)), facts)
	}
}
