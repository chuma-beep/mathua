package machinelearning

// Wave 7 machine learning, first batch: lasso regression and gradient
// boosting. Every answer is an integer computed from published facts.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave7(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"ml.reg.lasso":              "lasso",
		"ml.tree.gradient_boosting": "gradient_boosting",
	} {
		reg.Register(id, &w7MLGen{kind: kind})
	}
}

type w7MLGen struct{ kind string }

func (g *w7MLGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "lasso":
		return w7Lasso(kp)
	case "gradient_boosting":
		return w7GradientBoosting(kp)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

func w7M(o int, q, a, ex string, f map[string]string) generator.Problem {
	if f == nil {
		f = map[string]string{}
	}
	f["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: f, KPOrdinal: o}
}

func w7MI(m map[string]string, k string, v int) { m[k] = strconv.Itoa(v) }

func w7MR(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func w7Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func w7Lasso(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		lambda := w7MR(2, 5)
		w1 := w7MR(-3, 3)
		w2 := w7MR(-3, 3)
		for w1 == 0 {
			w1 = w7MR(-3, 3)
		}
		f["lambda"] = strconv.Itoa(lambda)
		f["w1"] = strconv.Itoa(w1)
		f["w2"] = strconv.Itoa(w2)
		ans := lambda * (w7Abs(w1) + w7Abs(w2))
		return w7M(kp, fmt.Sprintf("A lasso model has penalty \\(\\lambda=%d\\) and weights \\(w_1=%d\\), \\(w_2=%d\\). What is the L1 penalty \\(\\lambda(|w_1|+|w_2|)\\)?", lambda, w1, w2), strconv.Itoa(ans),
			"Add the absolute weights and multiply by lambda.", f)
	case 2:
		c := w7MR(-6, 6)
		lambda := w7MR(1, 6)
		f["c"] = strconv.Itoa(c)
		f["lambda"] = strconv.Itoa(lambda)
		ans := c - lambda
		if c < 0 {
			ans = c + lambda
		}
		if w7Abs(c) <= lambda {
			ans = 0
		}
		return w7M(kp, fmt.Sprintf("In a one-dimensional lasso problem the least-squares correlation is %d and the penalty is \\(\\lambda=%d\\). What is the soft-thresholded coefficient?", c, lambda), strconv.Itoa(ans),
			"The lasso coefficient is the soft threshold \\(\\operatorname{sign}(c)\\max(|c|-\\lambda,0)\\).", f)
	default:
		x := w7MR(1, 4)
		w := w7MR(-3, 3)
		y := w7MR(0, 10)
		lambda := w7MR(1, 4)
		f["x"] = strconv.Itoa(x)
		f["w"] = strconv.Itoa(w)
		f["y"] = strconv.Itoa(y)
		f["lambda"] = strconv.Itoa(lambda)
		r := y - x*w
		ans := r*r + 2*lambda*w7Abs(w)
		return w7M(kp, fmt.Sprintf("A one-dimensional lasso has \\(x=%d\\), weight \\(w=%d\\), target \\(y=%d\\) and penalty \\(\\lambda=%d\\). What is twice the objective \\(\\tfrac12(y-xw)^2+\\lambda|w|\\)?", x, w, y, lambda), strconv.Itoa(ans),
			"Twice the objective is \\((y-xw)^2+2\\lambda|w|\\).", f)
	}
}

func w7GradientBoosting(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		y := w7MR(5, 15)
		pred := w7MR(0, 5)
		f["y"] = strconv.Itoa(y)
		f["pred"] = strconv.Itoa(pred)
		return w7M(kp, fmt.Sprintf("A gradient-boosting model predicts %d for a target value %d. What is the residual \\(y-\\hat y\\)?", pred, y), strconv.Itoa(y-pred),
			"The residual is the target minus the current prediction.", f)
	case 2:
		y := w7MR(6, 16)
		pred := w7MR(0, 5)
		den := 2 + rand.Intn(2) // 2 or 3
		for (y-pred)%den != 0 {
			y = w7MR(6, 16)
		}
		f["y"] = strconv.Itoa(y)
		f["pred"] = strconv.Itoa(pred)
		f["den"] = strconv.Itoa(den)
		res := y - pred
		ans := pred + res/den
		return w7M(kp, fmt.Sprintf("A gradient-boosting model has prediction %d, target %d and learning rate \\(1/%d\\). What is the updated prediction after fitting the residual exactly?", pred, y, den), strconv.Itoa(ans),
			"Add the residual divided by the learning-rate denominator.", f)
	default:
		base := w7MR(0, 5)
		r1 := w7MR(-4, 4)
		r2 := w7MR(-4, 4)
		if (r1+r2)%2 != 0 {
			r2++
		}
		f["base"] = strconv.Itoa(base)
		f["r1"] = strconv.Itoa(r1)
		f["r2"] = strconv.Itoa(r2)
		ans := base + (r1+r2)/2
		return w7M(kp, fmt.Sprintf("A gradient-boosting model starts at %d. Two trees fit residuals %d and %d with learning rate \\(1/2\\). What is the final prediction?", base, r1, r2), strconv.Itoa(ans),
			"Add half of each residual to the base prediction.", f)
	}
}
