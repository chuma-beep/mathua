package probability

// Wave 6 probability theory, first batch: branching processes, Brownian motion
// and Chernoff bounds. The answers are exact integers; the facts publish the
// distribution, times, barriers or Bernoulli parameter so wave6_test.go can
// recompute every value.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave6(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"prob.process.branching": "branching",
		"prob.process.brownian":  "brownian",
		"prob.conc.chernoff":     "chernoff",
	} {
		reg.Register(id, &w6ProbGen{kind: kind})
	}
}

type w6ProbGen struct{ kind string }

func (g *w6ProbGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "branching":
		return w6Branching(kp)
	case "brownian":
		return w6Brownian(kp)
	case "chernoff":
		return w6Chernoff(kp)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

func w6P(o int, q, a, ex string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}

func w6PI(m map[string]string, k string, v int) { m[k] = strconv.Itoa(v) }

func w6PR(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

// Branching distributions use denominator 8. Each tuple is p0,p1,p2.
var w6BranchDist = [][3]int{
	{2, 2, 4},
	{1, 3, 4},
	{0, 0, 8},
	{4, 4, 0},
	{2, 4, 2},
	{3, 3, 2},
	{4, 2, 2},
	{8, 0, 0},
	{0, 4, 4},
	{0, 8, 0},
	{1, 7, 0},
	{7, 1, 0},
	{5, 3, 0},
	{6, 2, 0},
	{2, 6, 0},
	{3, 5, 0},
	{0, 3, 5},
	{0, 5, 3},
	{0, 2, 6},
	{0, 6, 2},
}

func w6BranchExt4(p [3]int) int {
	if p[0] == 0 {
		return 0
	}
	if p[2] == 0 {
		if p[1] == 8 {
			return 0
		}
		return 4
	}
	switch p {
	case [3]int{2, 2, 4}:
		return 2
	case [3]int{1, 3, 4}:
		return 1
	default:
		return 4
	}
}

var w6BranchExtPool = [][3]int{{2, 2, 4}, {1, 3, 4}, {2, 4, 2}, {3, 3, 2}, {4, 2, 2}}

func w6Branching(kp int) generator.Problem {
	pool := w6BranchDist
	if kp == 2 {
		pool = w6BranchExtPool
	}
	p := pool[rand.Intn(len(pool))]
	p0, p1, p2 := p[0], p[1], p[2]
	f := map[string]string{}
	w6PI(f, "p0", p0)
	w6PI(f, "p1", p1)
	w6PI(f, "p2", p2)
	w6PI(f, "mean8", p1+2*p2)
	w6PI(f, "ext4", w6BranchExt4(p))
	gen512 := 64*p0 + 8*p1*p0 + p2*p0*p0
	w6PI(f, "gen2_512", gen512)
	switch kp {
	case 1:
		return w6P(kp, fmt.Sprintf("A branching process has offspring distribution P(0)=%d/8, P(1)=%d/8, P(2)=%d/8. What is 8 times the mean offspring count?", p0, p1, p2), strconv.Itoa(p1+2*p2),
			fmt.Sprintf("The mean is P(1)+2P(2) = %d/8 + 2(%d/8) = %d/8. Eight times it is %d.", p1, p2, p1+2*p2, p1+2*p2), f)
	case 2:
		return w6P(kp, fmt.Sprintf("For the branching process with P(0)=%d/8, P(1)=%d/8 and P(2)=%d/8, what is 4 times the extinction probability?", p0, p1, p2), strconv.Itoa(w6BranchExt4(p)),
			fmt.Sprintf("The extinction probability is the smallest fixed point of p0 + p1 s + p2 s^2. For this distribution four times it is %d.", w6BranchExt4(p)), f)
	default:
		return w6P(kp, fmt.Sprintf("For the branching process with P(0)=%d/8, P(1)=%d/8 and P(2)=%d/8, what is 512 times the probability that the process is extinct by generation 2?", p0, p1, p2), strconv.Itoa(gen512),
			fmt.Sprintf("Start with generation 0 extinct probability 0. Apply f(s)=p0+p1 s+p2 s^2 twice, starting from 0. The requested scaled value is %d.", gen512), f)
	}
}

func w6Brownian(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		t := w6PR(1, 6)
		c := w6PR(1, 4)
		w6PI(f, "t", t)
		w6PI(f, "c", c)
		ans := c * c * t
		proc := fmt.Sprintf("B_%d", t)
		if c != 1 {
			proc = fmt.Sprintf("(%d B_%d)", c, t)
		}
		return w6P(kp, fmt.Sprintf("For standard Brownian motion B_t, what is E[%s^2]?", proc), strconv.Itoa(ans),
			fmt.Sprintf("E[B_t^2]=t=%d. Multiplying the process by %d squares the scale, so the answer is %d.", t, c, ans), f)
	case 2:
		s := w6PR(1, 4)
		t := s + w6PR(1, 5)
		w6PI(f, "s", s)
		w6PI(f, "t", t)
		return w6P(kp, fmt.Sprintf("For standard Brownian motion, what is Cov(B_%d,B_%d)?", s, t), strconv.Itoa(s),
			fmt.Sprintf("The covariance is min(s,t). Since %d is smaller, the covariance is %d.", s, s), f)
	default:
		a := w6PR(1, 6)
		b := w6PR(1, 6)
		w6PI(f, "a", a)
		w6PI(f, "b", b)
		return w6P(kp, fmt.Sprintf("Standard Brownian motion starts at 0. What is (%d+%d) times the probability that it hits level %d before level -%d?", a, b, a, b), strconv.Itoa(b),
			fmt.Sprintf("The hitting probability is b/(a+b) = %d/%d. Multiplying by a+b gives %d.", b, a+b, b), f)
	}
}

func w6Chernoff(kp int) generator.Problem {
	k := w6PR(1, 7)
	m := w6PR(1, 3)
	pow := func(base, exp int) int {
		r := 1
		for i := 0; i < exp; i++ {
			r *= base
		}
		return r
	}
	mgf := pow(8+k, m)
	f := map[string]string{}
	w6PI(f, "p_num", k)
	w6PI(f, "p_den", 8)
	w6PI(f, "trials", m)
	w6PI(f, "eight_mgf", mgf)
	w6PI(f, "eight_mean", m*k)
	w6PI(f, "sixteen_bound", mgf)
	scaleText := "8"
	if m > 1 {
		scaleText = fmt.Sprintf("8^%d", m)
	}
	switch kp {
	case 1:
		return w6P(kp, fmt.Sprintf("Let X be the sum of %d independent Bernoulli(%d/8) variables. What is %s times the moment generating function M(t) at t=ln 2?", m, k, scaleText), strconv.Itoa(mgf),
			fmt.Sprintf("M(ln 2)=(1+p)^%d=((8+%d)/8)^%d. Multiplying by 8^%d gives %d.", m, k, m, m, mgf), f)
	case 2:
		return w6P(kp, fmt.Sprintf("Let X be the sum of %d independent Bernoulli(%d/8) variables. What is 8 times the mean of X?", m, k), strconv.Itoa(m*k),
			fmt.Sprintf("The mean is %d times %d/8. Eight times it is %d.", m, k, m*k), f)
	default:
		return w6P(kp, fmt.Sprintf("Let X be the sum of %d independent Bernoulli(%d/8) variables. Use the Chernoff factor at t=ln 2. What is 16^%d times the bound on P(X>=%d)?", m, k, m, m), strconv.Itoa(mgf),
			fmt.Sprintf("The factor is ((8+%d)/16)^%d. Multiplying by 16^%d gives %d.", k, m, m, mgf), f)
	}
}
