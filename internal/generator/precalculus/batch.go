package precalculus

// The 2026-10-18 precalculus batch: nine concepts that take the domain from 41
// to the target of 50 and complete the thousand-concept milestone.
//
// The survey found the domain's unreachable concepts were all standard
// precalculus (conics, sequences and series, vectors, the binomial theorem,
// induction, parametric and polar graphs), so the wiring covers the subject
// rather than leaving it half-taught. The batch itself fills nine remaining
// gaps: continuous compounding, Newton's law of cooling, logarithmic scales,
// polynomial and rational inequalities, holes in rational graphs, the
// discriminant of a general conic, projectile motion, resultant vectors, and
// decomposition of a function.
//
// Several obvious additions were already taught and were not repeated:
// half-life in the growth lesson, logarithm properties in the logarithm
// lessons, eccentricity in the conic lessons, unit vectors in the vectors
// lesson, the midpoint in the distance lesson, even and odd functions in the
// functions lessons, and one-to-one functions in the inverse lesson.
//
// Six concepts are numeric, one is an expression, and two are multiple choice.
// Every numeric generator publishes the facts its answer follows from, so the
// tests recompute the answer from those facts.

import (
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerBatch(reg *generator.Registry) {
	reg.Register("precalc.exp.continuous", &continuousCompoundingGen{})
	reg.Register("precalc.exp.cooling", &coolingGen{})
	reg.Register("precalc.log.scale", &logScaleGen{})
	reg.Register("precalc.poly.inequality", &polyInequalityGen{})
	reg.Register("precalc.rational.holes", &rationalHoleGen{})
	reg.Register("precalc.conic.general", &generalConicGen{})
	reg.Register("precalc.parametric.projectile", &projectileGen{})
	reg.Register("precalc.vector.resultant", &resultantGen{})
	reg.Register("precalc.func.decomposition", &decompositionGen{})
}

// ---------------------------------------------------------------------------
// 1. Continuous compounding
// ---------------------------------------------------------------------------

type continuousCompoundingGen struct{}

func (g *continuousCompoundingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question    string
		correct     string
		distractors []string
		reason      string
	}
	claims := []claim{
		{
			question:    "The amount after t years for principal P at annual rate r compounded continuously is...",
			correct:     "A = P e^(rt)",
			distractors: []string{"A = P(1+r)^t", "A = P(1+r/n)^(nt)", "A = P + rt"},
			reason:      "Continuous compounding is the limit of P(1+r/n)^(nt) as n tends to infinity, which is P e^(rt).",
		},
		{
			question:    "As the number of compounding periods per year increases without bound, the amount approaches...",
			correct:     "P e^(rt)",
			distractors: []string{"P(1+rt)", "the principal alone", "the simple-interest amount"},
			reason:      "The limit of (1+r/n)^(nt) is e^(rt), so continuous compounding is the limiting case.",
		},
		{
			question:    "Continuous compounding is the limit of...",
			correct:     "compound interest as the compounding frequency tends to infinity",
			distractors: []string{"simple interest as t grows", "the doubling time", "linear growth"},
			reason:      "More frequent compounding converges to the exponential formula rather than growing without bound.",
		},
		{
			question:    "If 1000 is invested at 5% compounded continuously, the amount after t years is...",
			correct:     "1000 e^(0.05t)",
			distractors: []string{"1000(1.05)^t", "1000(1+0.05t)", "1000 e^(5t)"},
			reason:      "Substitute P = 1000 and r = 0.05 into A = P e^(rt).",
		},
		{
			question:    "For the same positive nominal rate, which grows faster: annual compounding or continuous compounding?",
			correct:     "continuous compounding",
			distractors: []string{"annual compounding", "they are equal", "it depends on the principal"},
			reason:      "e^r exceeds 1+r for positive r, so continuous compounding is the upper limit of more frequent compounding.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}

// ---------------------------------------------------------------------------
// 2. Newton's law of cooling
// ---------------------------------------------------------------------------

type coolingGen struct{}

func (g *coolingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question    string
		correct     string
		distractors []string
		reason      string
	}
	claims := []claim{
		{
			question:    "Newton's law of cooling says the rate of cooling is proportional to...",
			correct:     "the difference between the object and ambient temperatures",
			distractors: []string{"the object temperature alone", "the ambient temperature alone", "the time elapsed"},
			reason:      "The rate is k times the temperature excess, which is why the excess decays exponentially.",
		},
		{
			question:    "The solution of Newton's law of cooling has the form...",
			correct:     "T(t) = T_env + (T0 - T_env)e^(-kt)",
			distractors: []string{"T(t) = T0 e^(kt)", "T(t) = T_env + kt", "T(t) = T0 + T_env t"},
			reason:      "The excess T - T_env decays exponentially, so the temperature approaches the ambient value.",
		},
		{
			question:    "As t grows, the object's temperature approaches...",
			correct:     "the ambient temperature",
			distractors: []string{"zero", "the initial temperature", "infinity"},
			reason:      "The exponential term vanishes, leaving T = T_env.",
		},
		{
			question:    "A body at 90 C in a 20 C room is modeled by T(t) = 20 + 70e^(-0.1t). What is the ambient temperature?",
			correct:     "20",
			distractors: []string{"70", "90", "0"},
			reason:      "The ambient temperature is the constant term in the model, 20.",
		},
		{
			question:    "The cooling constant k is positive because...",
			correct:     "the temperature difference decays",
			distractors: []string{"the temperature grows", "the ambient temperature changes", "the object heats up"},
			reason:      "A positive k makes e^(-kt) decrease, so the excess shrinks toward zero.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}

// ---------------------------------------------------------------------------
// 3. Logarithmic scales (numeric)
// ---------------------------------------------------------------------------

type logScaleGen struct{}

func (g *logScaleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		kind     string
		e        int
		reason   string
	}
	claims := []claim{
		{
			question: "The pH of a solution is -log10[H+]. What is the pH when [H+] = 10^-3? (enter a number)",
			answer:   "3", kind: "ph", e: 3,
			reason: "pH = -log10(10^-3) = 3.",
		},
		{
			question: "The pH of a solution is -log10[H+]. What is the pH when [H+] = 10^-5? (enter a number)",
			answer:   "5", kind: "ph", e: 5,
			reason: "pH = -log10(10^-5) = 5.",
		},
		{
			question: "The pH of a solution is -log10[H+]. What is the pH when [H+] = 10^-7? (enter a number)",
			answer:   "7", kind: "ph", e: 7,
			reason: "pH = -log10(10^-7) = 7.",
		},
		{
			question: "The Richter magnitude of an earthquake is log10(A/A0). What is the magnitude when A/A0 = 10^4? (enter a number)",
			answer:   "4", kind: "richter", e: 4,
			reason: "log10(10^4) = 4.",
		},
		{
			question: "Sound level in decibels is 10 log10(I/I0). What is the level when I/I0 = 10^6? (enter a number)",
			answer:   "60", kind: "db", e: 6,
			reason: "10 log10(10^6) = 10(6) = 60.",
		},
		{
			question: "Sound level in decibels is 10 log10(I/I0). What is the level when I/I0 = 10^3? (enter a number)",
			answer:   "30", kind: "db", e: 3,
			reason: "10 log10(10^3) = 10(3) = 30.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts:       map[string]string{"kind": c.kind, "e": strconv.Itoa(c.e), "answer": c.answer},
	}
}

// ---------------------------------------------------------------------------
// 4. Polynomial and rational inequalities (numeric)
// ---------------------------------------------------------------------------

type polyInequalityGen struct{}

func (g *polyInequalityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		p, q     int
		reason   string
	}
	claims := []claim{
		{
			question: "Solve (x - 1)(x - 3) < 0. The solution is the open interval (a, b). What is b - a? (enter a number)",
			answer:   "2", p: 1, q: 3,
			reason: "The product is negative between the roots, so the interval is (1,3) and its width is 2.",
		},
		{
			question: "Solve (x + 2)(x - 4) < 0. The solution is the open interval (a, b). What is b - a? (enter a number)",
			answer:   "6", p: -2, q: 4,
			reason: "The interval is (-2,4), whose width is 6.",
		},
		{
			question: "Solve x^2 - 4 < 0. The solution is the open interval (a, b). What is b - a? (enter a number)",
			answer:   "4", p: -2, q: 2,
			reason: "The roots are -2 and 2, so the interval is (-2,2) with width 4.",
		},
		{
			question: "Solve (x - 2)(x - 5) < 0. The solution is the open interval (a, b). What is b - a? (enter a number)",
			answer:   "3", p: 2, q: 5,
			reason: "The interval is (2,5), whose width is 3.",
		},
		{
			question: "Solve (x + 3)(x - 1) < 0. The solution is the open interval (a, b). What is b - a? (enter a number)",
			answer:   "4", p: -3, q: 1,
			reason: "The interval is (-3,1), whose width is 4.",
		},
		{
			question: "Solve x^2 - 9 < 0. The solution is the open interval (a, b). What is b - a? (enter a number)",
			answer:   "6", p: -3, q: 3,
			reason: "The roots are -3 and 3, so the interval is (-3,3) with width 6.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts:       map[string]string{"p": strconv.Itoa(c.p), "q": strconv.Itoa(c.q), "answer": c.answer},
	}
}

// ---------------------------------------------------------------------------
// 5. Holes in rational graphs (numeric)
// ---------------------------------------------------------------------------

type rationalHoleGen struct{}

func (g *rationalHoleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		root     int
		reason   string
	}
	claims := []claim{
		{
			question: "The function f(x) = (x^2 - 4)/(x - 2) has a hole at x = a. What is a? (enter a number)",
			answer:   "2", root: 2,
			reason: "The numerator factors as (x-2)(x+2), so x-2 cancels and leaves a hole at x = 2.",
		},
		{
			question: "The function f(x) = (x^2 - 9)/(x - 3) has a hole at x = a. What is a? (enter a number)",
			answer:   "3", root: 3,
			reason: "The numerator factors as (x-3)(x+3), so the cancelled factor gives a hole at x = 3.",
		},
		{
			question: "The function f(x) = (x^2 - 1)/(x + 1) has a hole at x = a. What is a? (enter a number)",
			answer:   "-1", root: -1,
			reason: "The numerator factors as (x+1)(x-1), so the cancelled factor gives a hole at x = -1.",
		},
		{
			question: "The function f(x) = (x^2 - 16)/(x - 4) has a hole at x = a. What is a? (enter a number)",
			answer:   "4", root: 4,
			reason: "The numerator factors as (x-4)(x+4), so x-4 cancels and leaves a hole at x = 4.",
		},
		{
			question: "The function f(x) = (x^2 - 25)/(x + 5) has a hole at x = a. What is a? (enter a number)",
			answer:   "-5", root: -5,
			reason: "The numerator factors as (x+5)(x-5), so the cancelled factor gives a hole at x = -5.",
		},
		{
			question: "The function f(x) = (x^2 - 36)/(x - 6) has a hole at x = a. What is a? (enter a number)",
			answer:   "6", root: 6,
			reason: "The numerator factors as (x-6)(x+6), so x-6 cancels and leaves a hole at x = 6.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts:       map[string]string{"root": strconv.Itoa(c.root), "answer": c.answer},
	}
}

// ---------------------------------------------------------------------------
// 6. Classifying a general conic (numeric)
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// 7. Projectile motion (numeric)
// ---------------------------------------------------------------------------

type projectileGen struct{}

func (g *projectileGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		v        int
		reason   string
	}
	claims := []claim{
		{
			question: "A projectile is launched vertically with initial velocity 20 m/s. Taking g = 10 m/s^2, its height is y = 20t - 5t^2. What is its maximum height? (enter a number)",
			answer:   "20", v: 20,
			reason: "The vertex is at t = v/g = 2, and y(2) = 40 - 20 = 20.",
		},
		{
			question: "A projectile is launched vertically with initial velocity 10 m/s. Taking g = 10 m/s^2, its height is y = 10t - 5t^2. What is its maximum height? (enter a number)",
			answer:   "5", v: 10,
			reason: "The maximum height is v^2/(2g) = 100/20 = 5.",
		},
		{
			question: "A projectile is launched vertically with initial velocity 30 m/s. Taking g = 10 m/s^2, its height is y = 30t - 5t^2. What is its maximum height? (enter a number)",
			answer:   "45", v: 30,
			reason: "The maximum height is v^2/(2g) = 900/20 = 45.",
		},
		{
			question: "A projectile is launched vertically with initial velocity 40 m/s. Taking g = 10 m/s^2, its height is y = 40t - 5t^2. What is its maximum height? (enter a number)",
			answer:   "80", v: 40,
			reason: "The maximum height is v^2/(2g) = 1600/20 = 80.",
		},
		{
			question: "A projectile is launched vertically with initial velocity 5 m/s. Taking g = 10 m/s^2, its height is y = 5t - 5t^2. What is its maximum height? (enter a decimal)",
			answer:   "1.25", v: 5,
			reason: "The maximum height is v^2/(2g) = 25/20 = 1.25.",
		},
		{
			question: "A projectile is launched vertically with initial velocity 15 m/s. Taking g = 10 m/s^2, its height is y = 15t - 5t^2. What is its maximum height? (enter a decimal)",
			answer:   "11.25", v: 15,
			reason: "The maximum height is v^2/(2g) = 225/20 = 11.25.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts:       map[string]string{"v": strconv.Itoa(c.v), "answer": c.answer},
	}
}

// ---------------------------------------------------------------------------
// 8. Resultant vectors (numeric)
// ---------------------------------------------------------------------------

type resultantGen struct{}

func (g *resultantGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		a, b     int
		reason   string
	}
	claims := []claim{
		{
			question: "A vector has perpendicular components 3 and 4. What is its magnitude? (enter a number)",
			answer:   "5", a: 3, b: 4,
			reason: "The magnitude is sqrt(3^2 + 4^2) = sqrt(25) = 5.",
		},
		{
			question: "A vector has perpendicular components 6 and 8. What is its magnitude? (enter a number)",
			answer:   "10", a: 6, b: 8,
			reason: "The magnitude is sqrt(36 + 64) = 10.",
		},
		{
			question: "A vector has perpendicular components 5 and 12. What is its magnitude? (enter a number)",
			answer:   "13", a: 5, b: 12,
			reason: "The magnitude is sqrt(25 + 144) = 13.",
		},
		{
			question: "A vector has perpendicular components 8 and 15. What is its magnitude? (enter a number)",
			answer:   "17", a: 8, b: 15,
			reason: "The magnitude is sqrt(64 + 225) = 17.",
		},
		{
			question: "A vector has perpendicular components 7 and 24. What is its magnitude? (enter a number)",
			answer:   "25", a: 7, b: 24,
			reason: "The magnitude is sqrt(49 + 576) = 25.",
		},
		{
			question: "A vector has perpendicular components 9 and 12. What is its magnitude? (enter a number)",
			answer:   "15", a: 9, b: 12,
			reason: "The magnitude is sqrt(81 + 144) = 15.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts:       map[string]string{"a": strconv.Itoa(c.a), "b": strconv.Itoa(c.b), "answer": c.answer},
	}
}

// ---------------------------------------------------------------------------
// 9. Decomposing a function (expression)
// ---------------------------------------------------------------------------
