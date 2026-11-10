package machinelearning

// Wave 9 machine-learning concept: log loss as a proper scoring rule. The
// generator uses base-2 probabilities of the form 1/2^k so every answer is an
// exact integer number of bits.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave9(reg *generator.Registry) {
	reg.Register("ml.eval.log_loss", &w9LogLossGen{})
}

type w9LogLossGen struct{}

func (g *w9LogLossGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	k := 1 + rand.Intn(20)
	n := 2 + rand.Intn(8)
	f := map[string]string{}
	w9MI(f, "k", k)
	w9MI(f, "n", n)
	switch kp {
	case 1:
		return w9MP(kp, fmt.Sprintf("A positive example is predicted with probability 1/%d. What is its log loss in bits? (enter a number)", 1<<k), strconv.Itoa(k),
			fmt.Sprintf("For a positive label, log loss is -log2(p). With p=1/%d, the value is %d bits.", 1<<k, k), f)
	case 2:
		bitWord := "bits"
		if k == 1 {
			bitWord = "bit"
		}
		return w9MP(kp, fmt.Sprintf("%d positive examples are each predicted with probability 1/%d. What is the total log loss in bits? (enter a number)", n, 1<<k), strconv.Itoa(n*k),
			fmt.Sprintf("Each example contributes %d %s, so %d examples contribute %d bits.", k, bitWord, n, n*k), f)
	default:
		return w9MP(kp, fmt.Sprintf("On %d positive examples, model A always predicts 1/2 and model B predicts 1/%d. How many extra bits does model B use in total? (enter a number)", n, 1<<k), strconv.Itoa((k-1)*n),
			fmt.Sprintf("Model A uses %d bits and model B uses %d bits. The extra is %d bits.", n, k*n, (k-1)*n), f)
	}
}

func w9MI(m map[string]string, k string, v int) { m[k] = strconv.Itoa(v) }

func w9MP(o int, q, a, ex string, facts map[string]string) generator.Problem {
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}
