package machinelearning

// Wave 9 remaining machine learning: SMOTE as a synthetic interpolation
// construction. Every answer is exact and the operands are published.

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

var w9rPureNumber = regexp.MustCompile(`\\\((-?\d+)\\\)`)
var w9rUnitCoefficient = regexp.MustCompile(`(^|[^0-9])1 ?([xyzuvtr])([^A-Za-z]|$)`)
var w9rUnitParen = regexp.MustCompile(`(^|[^0-9])1\(`)
var w9rSpaces = regexp.MustCompile(`[ \t]+`)

// w9rClean removes rendering defects that a learner should not have to read.
func w9rClean(s string) string {
	s = w9rPureNumber.ReplaceAllString(s, "$1")
	s = strings.ReplaceAll(s, "+-", "-")
	s = strings.ReplaceAll(s, "--", "+")
	s = w9rSpaces.ReplaceAllString(s, " ")
	s = w9rUnitCoefficient.ReplaceAllString(s, "${1}${2}${3}")
	s = w9rUnitParen.ReplaceAllString(s, "${1}(")
	return s
}

func registerWave9Remaining(reg *generator.Registry) {
	reg.Register("ml.data.smote", &w9rSMOTEGen{})
}

type w9rSMOTEGen struct{}

func (g *w9rSMOTEGen) Generate(ctx generator.GeneratorContext) (prob generator.Problem) {
	defer func() {
		prob.Question = w9rClean(prob.Question)
		prob.Explanation = w9rClean(prob.Explanation)
	}()
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	x1 := 2 * w9rMR(-5, 5)
	x2 := 2 * w9rMR(-5, 5)
	y1 := 2 * w9rMR(-5, 5)
	y2 := 2 * w9rMR(-5, 5)
	k := w9rMR(2, 6)
	n := w9rMR(1, 8)
	m := w9rMR(5, 12)
	if x1 == x2 && y1 == y2 {
		x2 += 2
	}
	f := map[string]string{}
	w9rMI(f, "x1", x1)
	w9rMI(f, "y1", y1)
	w9rMI(f, "x2", x2)
	w9rMI(f, "y2", y2)
	w9rMI(f, "k", k)
	w9rMI(f, "n", n)
	w9rMI(f, "m", m)
	switch kp {
	case 1:
		xnew := (x1 + x2) / 2
		return w9rMP(kp, fmt.Sprintf("SMOTE averages the minority point \\((%d,%d)\\) and a neighbour \\((%d,%d)\\). What is the x-coordinate of the synthetic point? (enter a number)", x1, y1, x2, y2), strconv.Itoa(xnew),
			fmt.Sprintf("The midpoint x-coordinate is \\((%d+%d)/2=%d\\).", x1, x2, xnew), f)
	case 2:
		return w9rMP(kp, fmt.Sprintf("SMOTE chooses from the %d nearest minority neighbours. How many neighbours are considered? (enter a number)", k), strconv.Itoa(k),
			fmt.Sprintf("The neighbour count is %d.", k), f)
	default:
		return w9rMP(kp, fmt.Sprintf("A minority class has %d examples. SMOTE creates %d synthetic examples. How many minority examples are there now? (enter a number)", m, n), strconv.Itoa(m+n),
			fmt.Sprintf("The new total is \\(%d+%d=%d\\).", m, n, m+n), f)
	}
}

func w9rMR(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func w9rMI(m map[string]string, k string, v int) { m[k] = strconv.Itoa(v) }

func w9rMP(o int, q, a, ex string, facts map[string]string) generator.Problem {
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}
