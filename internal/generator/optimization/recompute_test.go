package optimization

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

// A recomputeCase parses a drawn question back into the parameters it printed
// and derives the answer a second time, independently of the generator. It
// exists because TestEveryAnswerIsNumeric accepts any number that parses: the
// only way to catch a wrong answer is to work it out again from the instance.
type recomputeCase struct {
	id     string
	kp     int
	parse  *regexp.Regexp
	expect func(p []int) (int, bool)
}

func atoiOrZero(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

func parseGroups(re *regexp.Regexp, q string) []int {
	m := re.FindStringSubmatch(q)
	if m == nil {
		return nil
	}
	out := make([]int, 0, len(m)-1)
	for _, g := range m[1:] {
		out = append(out, atoiOrZero(g))
	}
	return out
}

func clampInt(v, low, high int) int {
	if v < low {
		return low
	}
	if v > high {
		return high
	}
	return v
}

func maxOf(a, b int) int {
	if a > b {
		return a
	}
	return b
}

var recomputeCases = []recomputeCase{
	{
		id:    "optim.convex_set",
		kp:    1,
		parse: regexp.MustCompile(`convex combination of \((-?\d+), (-?\d+)\) and \((-?\d+), (-?\d+)\)`),
		expect: func(p []int) (int, bool) {
			sum := p[0] + p[1] + p[2] + p[3]
			if sum%2 != 0 {
				return 0, false
			}
			return sum / 2, true
		},
	},
	{
		id:    "optim.strict_convexity",
		kp:    3,
		parse: regexp.MustCompile(`interval from (-?\d+) to (-?\d+)`),
		expect: func(p []int) (int, bool) {
			d := p[1] - p[0]
			if d < 0 || d%2 != 0 {
				return 0, false
			}
			return -(d * d) / 4, true
		},
	},
	{
		id:    "optim.first_order",
		kp:    2,
		parse: regexp.MustCompile(`f\(x\) = (\d+) x\^2 - (\d+) x`),
		expect: func(p []int) (int, bool) {
			a, b := p[0], p[1]
			if b%(2*a) != 0 {
				return 0, false
			}
			return b / (2 * a), true
		},
	},
	{
		id:    "optim.second_order",
		kp:    3,
		parse: regexp.MustCompile(`f\(x, y\) = (\d+) x\^2 - (\d+) y\^2, evaluate f at the point \((-?\d+), (-?\d+)\)`),
		expect: func(p []int) (int, bool) {
			return p[0]*p[2]*p[2] - p[1]*p[3]*p[3], true
		},
	},
	{
		id:    "optim.newton_optim",
		kp:    3,
		parse: regexp.MustCompile(`second derivative (\d+) and gradient (-?\d+) at x = (-?\d+)`),
		expect: func(p []int) (int, bool) {
			a, grad, x := p[0], p[1], p[2]
			if grad != a*x {
				return 0, false
			}
			// Gradient step of size 1 moves by grad, Newton by grad/a.
			return a, true
		},
	},
	{
		id:    "optim.wolfe",
		kp:    2,
		parse: regexp.MustCompile(`f\(x\) = (\d+) x\^2 at x = (-?\d+)`),
		expect: func(p []int) (int, bool) {
			return 2 * p[0] * p[1], true
		},
	},
	{
		id:    "optim.wolfe",
		kp:    3,
		parse: regexp.MustCompile(`f\(x\) = x\^2 with x = (-?\d+) and step size (-?\d+)`),
		expect: func(p []int) (int, bool) {
			x, t := p[0], p[1]
			return x - 2*t*x, true
		},
	},
	{
		id:    "optim.trust_region_ratio",
		kp:    1,
		parse: regexp.MustCompile(`f\(x\) = (\d+) x\^2 at x = (-?\d+), a step of (-?\d+)`),
		expect: func(p []int) (int, bool) {
			a, x, s := p[0], p[1], p[2]
			return a * (x*x - (x+s)*(x+s)), true
		},
	},
	{
		id:    "optim.trust_region_subproblem",
		kp:    1,
		parse: regexp.MustCompile(`Minimize \((\d+)/2\)x\^2 \+ (-?\d+) x over \|x\| ≤ (\d+)`),
		expect: func(p []int) (int, bool) {
			a, grad, delta := p[0], p[1], p[2]
			if grad%a != 0 {
				return 0, false
			}
			return clampInt(-grad/a, -delta, delta), true
		},
	},
	{
		id:    "optim.active_set",
		kp:    1,
		parse: regexp.MustCompile(`At \((-?\d+), (-?\d+)\) the working set tests x ≤ (-?\d+) and y ≤ (-?\d+)`),
		expect: func(p []int) (int, bool) {
			return maxOf(p[0]-p[2], p[1]-p[3]), true
		},
	},
	{
		id:    "optim.lp_standard_form",
		kp:    2,
		parse: regexp.MustCompile(`(\d+) decision variables and (\d+) inequality constraints`),
		expect: func(p []int) (int, bool) {
			return p[0] + p[1], true
		},
	},
	{
		id:    "optim.weak_duality",
		kp:    2,
		parse: regexp.MustCompile(`a = (-?\d+), y = (-?\d+) and c = (-?\d+)`),
		expect: func(p []int) (int, bool) {
			a, y, c := p[0], p[1], p[2]
			return c - a*y, true
		},
	},
	{
		id:    "optim.complementary_slackness",
		kp:    3,
		parse: regexp.MustCompile(`y = (-?\d+), the constraint is (-?\d+) x ≤ (-?\d+) and x = (-?\d+)`),
		expect: func(p []int) (int, bool) {
			return p[0] * (p[2] - p[1]*p[3]), true
		},
	},
	{
		id:    "optim.lp_sensitivity",
		kp:    2,
		parse: regexp.MustCompile(`c = (-?\d+), a = (-?\d+) and y = (-?\d+)\. After c rises to (-?\d+)`),
		expect: func(p []int) (int, bool) {
			c, a, y, newC := p[0], p[1], p[2], p[3]
			if newC <= c {
				return 0, false
			}
			return newC - a*y, true
		},
	},
	{
		id:    "optim.slater",
		kp:    3,
		parse: regexp.MustCompile(`contains the two points (-?\d+) and (-?\d+)`),
		expect: func(p []int) (int, bool) {
			if (p[0]+p[1])%2 != 0 {
				return 0, false
			}
			return (p[0] + p[1]) / 2, true
		},
	},
	{
		id:    "optim.branch_bound",
		kp:    2,
		parse: regexp.MustCompile(`relaxation value (-?\d+)\.(-?\d+)`),
		expect: func(p []int) (int, bool) {
			return p[0], true
		},
	},
	{
		id:    "optim.kkt",
		kp:    3,
		parse: regexp.MustCompile(`multiplier (-?\d+) and the matching constraint slack (-?\d+)`),
		expect: func(p []int) (int, bool) {
			return p[0] * p[1], true
		},
	},
}

// Each repaired family is drawn repeatedly and checked against an independent
// derivation from the text the learner actually sees.
func TestRepairedFamiliesRecompute(t *testing.T) {
	reg := wave11Registry(t)
	for _, c := range recomputeCases {
		for draw := 0; draw < 40; draw++ {
			p, err := reg.GenerateContext(c.id, generator.GeneratorContext{KPOrdinal: c.kp})
			if err != nil {
				t.Fatalf("%s kp%d: %v", c.id, c.kp, err)
			}
			params := parseGroups(c.parse, p.Question)
			if params == nil {
				t.Fatalf("%s kp%d: question did not match the parser: %s", c.id, c.kp, p.Question)
			}
			want, ok := c.expect(params)
			if !ok {
				t.Fatalf("%s kp%d: drawn instance has no exact answer: %s", c.id, c.kp, p.Question)
			}
			if p.Answer != strconv.Itoa(want) {
				t.Fatalf("%s kp%d: answer %s, recomputed %d from %s", c.id, c.kp, p.Answer, want, p.Question)
			}
		}
	}
}

// Every family records its answer as a fact, and no ordinal collapses to a
// constant answer that a learner could memorise without doing the work.
func TestFactsAndVariety(t *testing.T) {
	reg := wave11Registry(t)
	for _, id := range wave11IDs {
		for _, kp := range []int{1, 2, 3} {
			seen := make(map[string]struct{})
			for draw := 0; draw < 60; draw++ {
				p, err := reg.GenerateContext(id, generator.GeneratorContext{KPOrdinal: kp})
				if err != nil {
					t.Fatalf("%s kp%d: %v", id, kp, err)
				}
				if got := p.Facts["answer"]; got != p.Answer {
					t.Fatalf("%s kp%d fact %q does not match answer %q", id, kp, got, p.Answer)
				}
				seen[p.Answer] = struct{}{}
			}
			if len(seen) < 3 {
				t.Fatalf("%s kp%d produced %d distinct answers in 60 draws, want at least 3", id, kp, len(seen))
			}
		}
	}
}
