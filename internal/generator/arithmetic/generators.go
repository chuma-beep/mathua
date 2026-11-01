package arithmetic

import (
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/grader"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

func Register(reg *generator.Registry) {
	reg.Register("arith.add.single", &addGen{minA: 1, maxA: 9, minB: 1, maxB: 9})
	reg.Register("arith.add.double", &addGen{minA: 10, maxA: 99, minB: 10, maxB: 99})
	reg.Register("arith.add.carry", &addGen{minA: 10, maxA: 99, minB: 10, maxB: 99, ensureCarry: true})
	reg.Register("arith.add.triple", &addGen{minA: 100, maxA: 999, minB: 100, maxB: 999})
	reg.Register("arith.add.word", &addWordGen{})

	reg.Register("arith.sub.single", &subGen{minA: 5, maxA: 9, minB: 1, maxB: 4})
	reg.Register("arith.sub.double", &subGen{minA: 20, maxA: 99, minB: 10, maxB: 50})
	reg.Register("arith.sub.borrow", &subBorrowGen{})
	reg.Register("arith.sub.word", &subWordGen{})

	reg.Register("arith.place.tens", &placeValueGen{max: 99, label: "tens"})
	reg.Register("arith.place.hundreds", &placeValueGen{max: 999, label: "hundreds"})
	reg.Register("arith.place.thousands", &placeValueGen{max: 9999, label: "thousands"})

	reg.Register("arith.round.tens", &roundGen{to: 10})
	reg.Register("arith.round.hundreds", &roundGen{to: 100})
	reg.Register("arith.round.thousands", &roundGen{to: 1000})

	reg.Register("arith.mult.concept", &multConceptGen{})
	reg.Register("arith.mult.2_5_10", &multByGen{factor: 0})
	reg.Register("arith.mult.tables", &multTablesGen{})
	reg.Register("arith.mult.double", &multGen{minA: 10, maxA: 99, minB: 2, maxB: 9})
	reg.Register("arith.mult.triple", &multGen{minA: 100, maxA: 999, minB: 2, maxB: 9})
	reg.Register("arith.mult.word", &multWordGen{})

	reg.Register("arith.div.concept", &divConceptGen{})
	reg.Register("arith.div.basic", &divBasicGen{})
	reg.Register("arith.div.remainder", &divRemainderGen{})
	reg.Register("arith.div.long", &divLongGen{})
	reg.Register("arith.div.word", &divWordGen{})

	reg.Register("arith.factor.find", &factorFindGen{})
	reg.Register("arith.factor.prime", &primeGen{})
	reg.Register("arith.factor.prime_fact", &primeFactGen{})
	reg.Register("arith.factor.gcf", &gcfGen{})
	reg.Register("arith.factor.lcm", &lcmGen{})
	reg.Register("arith.factor.composite", &compositeGen{})

	reg.Register("arith.exp.concept", &expConceptGen{})
	reg.Register("arith.exp.evaluate", &expEvalGen{})
	reg.Register("arith.exp.product_rule", &expProductRuleGen{})
	reg.Register("arith.exp.quotient_rule", &expQuotientRuleGen{})
	reg.Register("arith.exp.power_rule", &expPowerRuleGen{})

	reg.Register("arith.sqrt.perfect", &sqrtPerfectGen{})
	reg.Register("arith.sqrt.simplify", &sqrtSimplifyGen{})

	reg.Register("arith.neg.number_line", &negNumberLineGen{})
	reg.Register("arith.neg.add_sub", &negAddSubGen{})
	reg.Register("arith.neg.mult_div", &negMultDivGen{})

	reg.Register("arith.order_ops.basic", &orderOpsGen{parens: false, exponents: false})
	reg.Register("arith.order_ops.full", &orderOpsGen{parens: true, exponents: false})
	reg.Register("arith.order_ops.nested", &orderOpsGen{parens: true, exponents: true})

	reg.Register("arith.dec.intro", &decIntroGen{})

	registerBatch(reg)
}

type addGen struct {
	minA, maxA  int
	minB, maxB  int
	ensureCarry bool
}

func (g *addGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	diffA := g.maxA - g.minA
	if diffA <= 0 {
		diffA = 1
	}
	diffB := g.maxB - g.minB
	if diffB <= 0 {
		diffB = 1
	}
	scale := int(1 + ctx.Difficulty*float64(diffA)/10)
	scale = max(scale, 1)
	a := rand.Intn(min(diffA, scale*10)) + g.minA
	b := rand.Intn(min(diffB, scale*10)) + g.minB
	if g.ensureCarry {
		a = rand.Intn(90) + 10
		onesA := a % 10
		onesB := 10 - onesA + rand.Intn(9)
		if onesB <= 9 {
			b = (b/10)*10 + onesB
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("\\(%d + %d = ?\\)", a, b),
		Answer:      fmt.Sprintf("%d", a+b),
		Explanation: fmt.Sprintf("Given %d + %d. Step 1: start at %d. Step 2: count up %d → %d. Answer: %d", a, b, a, b, a+b, a+b),
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "sum": strconv.Itoa(a + b), "ones_a": strconv.Itoa(a % 10), "ones_b": strconv.Itoa(b % 10)},
	}
}

type addWordGen struct{}

func (g *addWordGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// Staged ladder: small counts → two-digit sums.
	maxN := 20
	if ctx.KPOrdinal == 1 {
		maxN = 10
	} else if ctx.KPOrdinal == 2 {
		maxN = 25
	} else if ctx.KPOrdinal == 3 {
		maxN = 50
	} else if ctx.Difficulty >= 0.7 {
		maxN = 50
	} else if ctx.Difficulty >= 0.45 {
		maxN = 35
	}
	a := rand.Intn(maxN) + 1
	b := rand.Intn(maxN) + 1
	items := []string{"apples", "marbles", "stickers", "crayons", "pencils"}
	item := items[rand.Intn(len(items))]
	return generator.Problem{
		Question:    fmt.Sprintf("You have %d %s. Your friend gives you %d more. How many do you have now?", a, item, b),
		Answer:      fmt.Sprintf("%d", a+b),
		Explanation: fmt.Sprintf("Given %d + %d %s. Step 1: start at %d. Step 2: count up %d → %d. Answer: %d %s", a, b, item, a, b, a+b, a+b, item),
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "sum": strconv.Itoa(a + b), "item": item},
	}
}

type subGen struct {
	minA, maxA int
	minB, maxB int
}

func (g *subGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := rand.Intn(g.maxA-g.minA+1) + g.minA
	bMax := a - g.minB
	if bMax < g.minB {
		bMax = a
	}
	b := rand.Intn(bMax-g.minB+1) + g.minB
	return generator.Problem{
		Question:    fmt.Sprintf("\\(%d - %d = ?\\)", a, b),
		Answer:      fmt.Sprintf("%d", a-b),
		Explanation: fmt.Sprintf("%d - %d = %d", a, b, a-b),
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "diff": strconv.Itoa(a - b)},
	}
}

type subBorrowGen struct{}

func (g *subBorrowGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	tensA := rand.Intn(max(1, scale*2)) + 1
	// onesA stays under 9 so that 10-onesA is always a positive bound. A zero
	// or negative bound panics rand.Intn, and the old range reached 11 at
	// difficulty 1.0, which the Learn stepper does request, taking the handler
	// down with a panic instead of serving a question.
	onesA := rand.Intn(max(1, min(8, scale*2-1))) + 1
	onesB := onesA + rand.Intn(10-onesA) + 1
	tensB := rand.Intn(tensA)
	if tensB == tensA {
		tensB--
	}
	a := tensA*10 + onesA
	b := tensB*10 + onesB
	return generator.Problem{
		Question:    fmt.Sprintf("\\(%d - %d = ?\\)", a, b),
		Answer:      fmt.Sprintf("%d", a-b),
		Explanation: fmt.Sprintf("%d - %d = %d (borrowing required)", a, b, a-b),
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "diff": strconv.Itoa(a - b), "tens_a": strconv.Itoa(tensA), "ones_a": strconv.Itoa(onesA), "tens_b": strconv.Itoa(tensB), "ones_b": strconv.Itoa(onesB), "ones_after_borrow": strconv.Itoa(onesA + 10), "ones_diff": strconv.Itoa(onesA + 10 - onesB), "tens_left": strconv.Itoa(tensA - 1 - tensB)},
	}
}

type subWordGen struct{}

func (g *subWordGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(max(1, scale*16)) + 20
	b := rand.Intn(a-1) + 1
	items := []string{"candies", "pencils", "cards", "coins", "beads"}
	item := items[rand.Intn(len(items))]
	return generator.Problem{
		Question:    fmt.Sprintf("You have %d %s. You give away %d. How many are left?", a, item, b),
		Answer:      fmt.Sprintf("%d", a-b),
		Explanation: fmt.Sprintf("%d - %d = %d %s left", a, b, a-b, item),
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "diff": strconv.Itoa(a - b), "item": item},
	}
}

type placeValueGen struct {
	max   int
	label string
}

func (g *placeValueGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	max := g.max
	if max <= 1 {
		max = 2
	}
	n := rand.Intn(max) + 1
	// The zero-value generator carries no label, and an empty label means the
	// tens place, so the default has to be usable rather than zero.
	placeUnit := 10
	switch g.label {
	case "hundreds":
		placeUnit = 100
	case "thousands":
		placeUnit = 1000
	}
	// A number below the place unit has no digit there: "the digit 0 is in the
	// tens place of 7" is not a question, so lift it into the place.
	if n < placeUnit {
		n += placeUnit
	}
	// A zero in the place makes the answer zero, which teaches nothing. Shift
	// to the next number whose digit there is nonzero: 7308 becomes 7310 for
	// tens, 7400 for hundreds.
	if (n/placeUnit)%10 == 0 {
		n += placeUnit - n%placeUnit
	}
	digit := (n / placeUnit) % 10
	value := digit * placeUnit
	// The stem names the digit and asks for its value. Asking for "the value of
	// the digit in the tens place" of 87 let a learner answer 8, which is the
	// digit, and be marked wrong for a question the stem did not settle.
	return generator.Problem{
		Question:    fmt.Sprintf("In the number %d, the digit %d is in the %s place. What is the value of that digit?", n, digit, g.label),
		Answer:      fmt.Sprintf("%d", value),
		Explanation: fmt.Sprintf("The %s digit is %d, value = %d x %d = %d", g.label, digit, digit, placeUnit, value),
		Facts:       map[string]string{"n": strconv.Itoa(n), "digit": strconv.Itoa(digit), "place_unit": strconv.Itoa(placeUnit), "value": strconv.Itoa(value), "label": g.label},
	}
}

type roundGen struct {
	to int
}

func (g *roundGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	to := g.to
	if to <= 0 {
		to = 1
	}
	lo := to
	hi := to * 10
	n := rand.Intn(hi-lo) + lo
	rounded := int(math.Round(float64(n)/float64(to))) * to
	return generator.Problem{
		Question:    fmt.Sprintf("Round %d to the nearest %d.", n, to),
		Answer:      fmt.Sprintf("%d", rounded),
		Explanation: fmt.Sprintf("%d rounded to nearest %d is %d.", n, to, rounded),
		Facts:       map[string]string{"n": strconv.Itoa(n), "to": strconv.Itoa(to), "rounded": strconv.Itoa(rounded), "remainder": strconv.Itoa(n % to)},
	}
}

type multConceptGen struct{}

func (g *multConceptGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(max(4, scale*3)) + 2
	b := rand.Intn(max(4, scale*3)) + 2
	return generator.Problem{
		Question:    fmt.Sprintf("%d groups of %d = ?", a, b),
		Answer:      fmt.Sprintf("%d", a*b),
		Explanation: fmt.Sprintf("%d groups of %d means %d x %d = %d", a, b, a, b, a*b),
		Facts:       map[string]string{"groups": strconv.Itoa(a), "per_group": strconv.Itoa(b), "total": strconv.Itoa(a * b)},
	}
}

type multByGen struct {
	factor int
}

func (g *multByGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	f := []int{2, 5, 10}[rand.Intn(3)]
	n := rand.Intn(max(1, scale*2)) + 1
	if ctx.Difficulty > 0.5 {
		n = rand.Intn(20) + 1
	}
	return generator.Problem{
		Question:    fmt.Sprintf("\\(%d \\times %d = ?\\)", f, n),
		Answer:      fmt.Sprintf("%d", f*n),
		Explanation: fmt.Sprintf("%d x %d = %d", f, n, f*n),
		Facts:       map[string]string{"factor": strconv.Itoa(f), "n": strconv.Itoa(n), "total": strconv.Itoa(f * n)},
	}
}

type multTablesGen struct{}

func (g *multTablesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(max(1, scale*2)) + 1
	b := rand.Intn(max(1, scale*2)) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("\\(%d \\times %d = ?\\)", a, b),
		Answer:      fmt.Sprintf("%d", a*b),
		Explanation: fmt.Sprintf("%d x %d = %d", a, b, a*b),
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "total": strconv.Itoa(a * b)},
	}
}

type multGen struct {
	minA, maxA int
	minB, maxB int
}

func (g *multGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := rand.Intn(g.maxA-g.minA+1) + g.minA
	b := rand.Intn(g.maxB-g.minB+1) + g.minB
	return generator.Problem{
		Question:    fmt.Sprintf("\\(%d \\times %d = ?\\)", a, b),
		Answer:      fmt.Sprintf("%d", a*b),
		Explanation: fmt.Sprintf("%d x %d = %d", a, b, a*b),
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "total": strconv.Itoa(a * b)},
	}
}

type multWordGen struct{}

func (g *multWordGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(max(1, scale*2)) + 1
	b := rand.Intn(max(1, scale*2)) + 1
	items := []string{"stickers", "cards", "beads", "coins", "marbles"}
	item := items[rand.Intn(len(items))]
	return generator.Problem{
		Question:    fmt.Sprintf("You have %d bags with %d %s each. How many %s in total?", a, b, item, item),
		Answer:      fmt.Sprintf("%d", a*b),
		Explanation: fmt.Sprintf("%d x %d = %d %s", a, b, a*b, item),
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "total": strconv.Itoa(a * b), "item": item},
	}
}

type divConceptGen struct{}

func (g *divConceptGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := rand.Intn(5) + 2
	b := rand.Intn(5) + 2
	total := a * b
	return generator.Problem{
		Question:    fmt.Sprintf("You have %d items. You split them into %d equal groups. How many in each group?", total, a),
		Answer:      fmt.Sprintf("%d", b),
		Explanation: fmt.Sprintf("%d / %d = %d in each group", total, a, b),
		Facts:       map[string]string{"total": strconv.Itoa(total), "groups": strconv.Itoa(a), "per_group": strconv.Itoa(b)},
	}
}

type divBasicGen struct{}

func (g *divBasicGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	b := rand.Intn(max(1, scale*2)) + 1
	a := b * (rand.Intn(max(1, scale*2)) + 1)
	return generator.Problem{
		Question:    fmt.Sprintf("\\(%d \\div %d = ?\\)", a, b),
		Answer:      fmt.Sprintf("%d", a/b),
		Explanation: fmt.Sprintf("%d / %d = %d", a, b, a/b),
		Facts:       map[string]string{"dividend": strconv.Itoa(a), "divisor": strconv.Itoa(b), "quotient": strconv.Itoa(a / b)},
	}
}

type divRemainderGen struct{}

// The question asks for the remainder alone. It used to ask for "Q R" and
// grade it with a regex, which needs a letter between two numbers: a phone
// keypad cannot type R, so the answer box collected nothing and the report
// arrived with an empty answer. A remainder is one integer, so the numeric
// keypad is enough and this generator needs no Grade method of its own.
func (g *divRemainderGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// Staged ladder: small divisor → larger divisor and quotient.
	divMax := 5
	if ctx.KPOrdinal == 1 {
		divMax = 5
	} else if ctx.KPOrdinal == 2 {
		divMax = 8
	} else if ctx.KPOrdinal == 3 {
		divMax = 12
	} else if ctx.Difficulty >= 0.7 {
		divMax = 12
	} else if ctx.Difficulty >= 0.45 {
		divMax = 8
	}
	b := rand.Intn(divMax-1) + 2
	r := rand.Intn(b-1) + 1
	a := b*(rand.Intn(divMax)+1) + r
	return generator.Problem{
		Question:    fmt.Sprintf("What is the remainder when %d is divided by %d?", a, b),
		Answer:      strconv.Itoa(r),
		Explanation: fmt.Sprintf("%d divided by %d is %d with %d left over: %d × %d = %d, and %d − %d = %d.", a, b, a/b, r, b, a/b, (a/b)*b, a, (a/b)*b, r),
		Facts:       map[string]string{"dividend": strconv.Itoa(a), "divisor": strconv.Itoa(b), "quotient": strconv.Itoa(a / b), "remainder": strconv.Itoa(r), "used": strconv.Itoa((a / b) * b)},
	}
}

type divLongGen struct{}

func (g *divLongGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	b := rand.Intn(max(1, scale*2)) + 2
	q := rand.Intn(50) + 10
	a := b * q
	return generator.Problem{
		Question:    fmt.Sprintf("\\(%d \\div %d = ?\\)", a, b),
		Answer:      fmt.Sprintf("%d", q),
		Explanation: fmt.Sprintf("Given %d ÷ %d. Step 1: ask how many %ds fit in %d. Step 2: %d × %d = %d. Answer: %d", a, b, b, a, b, q, a, q),
		Facts:       map[string]string{"dividend": strconv.Itoa(a), "divisor": strconv.Itoa(b), "quotient": strconv.Itoa(q)},
	}
}

type divWordGen struct{}

func (g *divWordGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	b := rand.Intn(max(1, scale*2)) + 2
	q := rand.Intn(12) + 1
	a := b * q
	items := []string{"cookies", "cards", "pencils", "marbles", "coins"}
	item := items[rand.Intn(len(items))]
	people := []string{"friends", "classmates", "teammates", "siblings"}
	person := people[rand.Intn(len(people))]
	return generator.Problem{
		Question:    fmt.Sprintf("You have %d %s shared equally among %d %s. How many does each person get?", a, item, b, person),
		Answer:      fmt.Sprintf("%d", q),
		Explanation: fmt.Sprintf("%d / %d = %d %s per person", a, b, q, item),
		Facts:       map[string]string{"total": strconv.Itoa(a), "people": strconv.Itoa(b), "quotient": strconv.Itoa(q), "item": item, "person": person},
	}
}

type factorFindGen struct{}

func parseFactorSet(s string) map[int]bool {
	out := make(map[int]bool)
	parts := strings.Split(s, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil
		}
		out[n] = true
	}
	return out
}

func (g *factorFindGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	n := rand.Intn(max(10, scale*10)) + 10
	var factors []int
	for i := 1; i <= n; i++ {
		if n%i == 0 {
			factors = append(factors, i)
		}
	}
	fStr := make([]string, len(factors))
	for i, f := range factors {
		fStr[i] = strconv.Itoa(f)
	}
	return generator.Problem{
		Question:    fmt.Sprintf("List all factors of %d. Answer with numbers separated by commas.", n),
		Answer:      strings.Join(fStr, ","),
		Explanation: fmt.Sprintf("Factors of %d: %v", n, factors),
		Facts:       map[string]string{"n": strconv.Itoa(n), "factors": strings.Join(fStr, ", ")},
	}
}

func (g *factorFindGen) Grade(expected, answer string) grader.Result {
	eSet := parseFactorSet(expected)
	aSet := parseFactorSet(answer)
	if eSet == nil || aSet == nil {
		return grader.Result{Correct: false, Score: 0, Feedback: "Enter factors separated by commas, e.g. 1,2,3,6"}
	}
	if len(eSet) != len(aSet) {
		return grader.Result{Correct: false, Score: 0, Feedback: "Incorrect number of factors"}
	}
	for k := range eSet {
		if !aSet[k] {
			return grader.Result{Correct: false, Score: 0, Feedback: "Incorrect"}
		}
	}
	return grader.Result{Correct: true, Score: 1}
}

type primeGen struct{}

func (g *primeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return primeProblem()
}

type primeFactGen struct{}

func (g *primeFactGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	n := rand.Intn(scale*20) + 2
	factors := primeFactors(n)
	fStr := make([]string, len(factors))
	for i, f := range factors {
		fStr[i] = strconv.Itoa(f)
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Find the prime factorization of %d. Enter the factors separated by commas, smallest first.", n),
		Answer:      strings.Join(fStr, ","),
		Explanation: fmt.Sprintf("%d = %s", n, strings.Join(fStr, " x ")),
		Facts:       map[string]string{"n": strconv.Itoa(n), "factors": strings.Join(fStr, " × ")},
	}
}

func primeFactors(n int) []int {
	var f []int
	for i := 2; i*i <= n; i++ {
		for n%i == 0 {
			f = append(f, i)
			n /= i
		}
	}
	if n > 1 {
		f = append(f, n)
	}
	return f
}

type gcfGen struct{}

func (g *gcfGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(max(10, scale*10)) + 10
	b := rand.Intn(max(10, scale*10)) + 10
	result := mathutil.GCD(a, b)
	return generator.Problem{
		Question:    fmt.Sprintf("Find the GCF of %d and %d.", a, b),
		Answer:      fmt.Sprintf("%d", result),
		Explanation: fmt.Sprintf("GCF(%d, %d) = %d", a, b, result),
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "gcd": strconv.Itoa(result)},
	}
}

type lcmGen struct{}

func (g *lcmGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(max(4, scale*4)) + 2
	b := rand.Intn(max(4, scale*4)) + 2
	l := a * b / mathutil.GCD(a, b)
	return generator.Problem{
		Question:    fmt.Sprintf("Find the LCM of %d and %d.", a, b),
		Answer:      fmt.Sprintf("%d", l),
		Explanation: fmt.Sprintf("LCM(%d, %d) = %d", a, b, l),
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "lcm": strconv.Itoa(l), "gcd": strconv.Itoa(mathutil.GCD(a, b))},
	}
}

type compositeGen struct{}

func (g *compositeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return compositeProblem()
}

type expConceptGen struct{}

func (g *expConceptGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// Difficulty still sets the size of the base, but the exponent reaches 5 so
	// that equal bases are not the only thing separating two draws.
	scale := int(1 + ctx.Difficulty*5)
	base := rand.Intn(3+scale*2) + 2
	exp := rand.Intn(4) + 2
	return generator.Problem{
		Question:    fmt.Sprintf("What is the value of \\(%d^{%d}\\)?", base, exp),
		Answer:      fmt.Sprintf("%d", mathutil.IntPow(base, exp)),
		Explanation: fmt.Sprintf("\\(%d^{%d} = %s = %d\\)", base, exp, strings.Repeat(fmt.Sprintf("%d \\times ", base), exp-1)+fmt.Sprintf("%d", base), mathutil.IntPow(base, exp)),
		Facts:       map[string]string{"base": strconv.Itoa(base), "exp": strconv.Itoa(exp), "value": strconv.Itoa(mathutil.IntPow(base, exp))},
	}
}

type expEvalGen struct{}

func (g *expEvalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	base := rand.Intn(max(1, scale*2)) + 2
	exp := rand.Intn(max(1, scale)) + 1
	if ctx.Difficulty > 0.5 {
		exp = rand.Intn(max(1, scale)) + 3
	}
	return generator.Problem{
		Question:    fmt.Sprintf("\\(%d^{%d} =\\) ?", base, exp),
		Answer:      fmt.Sprintf("%d", mathutil.IntPow(base, exp)),
		Explanation: fmt.Sprintf("\\(%d^{%d} = %d\\)", base, exp, mathutil.IntPow(base, exp)),
		Facts:       map[string]string{"base": strconv.Itoa(base), "exp": strconv.Itoa(exp), "value": strconv.Itoa(mathutil.IntPow(base, exp))},
	}
}

type expProductRuleGen struct{}

func (g *expProductRuleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	base := rand.Intn(max(3, scale*2)) + 2
	e1 := rand.Intn(max(3, scale*2)) + 1
	e2 := rand.Intn(max(3, scale*2)) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Simplify: \\(%d^{%d} \\times %d^{%d}\\)", base, e1, base, e2),
		Answer:      fmt.Sprintf("%d^%d", base, e1+e2),
		Explanation: fmt.Sprintf("\\(%d^{%d} \\times %d^{%d} = %d^{%d+%d} = %d^{%d}\\)", base, e1, base, e2, base, e1, e2, base, e1+e2),
		Facts:       map[string]string{"base": strconv.Itoa(base), "e1": strconv.Itoa(e1), "e2": strconv.Itoa(e2), "sum_exp": strconv.Itoa(e1 + e2)},
	}
}

type expQuotientRuleGen struct{}

func (g *expQuotientRuleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	base := rand.Intn(max(3, scale*2)) + 2
	e1 := rand.Intn(max(3, scale*2)) + 3
	e2 := rand.Intn(max(2, scale/2+1)) + 1
	if e2 >= e1 {
		e1 = e2 + 1
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Simplify: \\(\\frac{%d^{%d}}{%d^{%d}}\\)", base, e1, base, e2),
		Answer:      fmt.Sprintf("%d^%d", base, e1-e2),
		Explanation: fmt.Sprintf("\\(\\frac{%d^{%d}}{%d^{%d}} = %d^{%d-%d} = %d^{%d}\\)", base, e1, base, e2, base, e1, e2, base, e1-e2),
		Facts:       map[string]string{"base": strconv.Itoa(base), "e1": strconv.Itoa(e1), "e2": strconv.Itoa(e2), "diff_exp": strconv.Itoa(e1 - e2)},
	}
}

type expPowerRuleGen struct{}

func (g *expPowerRuleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	base := rand.Intn(max(3, scale*2)) + 2
	e1 := rand.Intn(max(3, scale*2)) + 1
	e2 := rand.Intn(max(3, scale*2)) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Simplify: \\((%d^{%d})^{%d}\\)", base, e1, e2),
		Answer:      fmt.Sprintf("%d^%d", base, e1*e2),
		Explanation: fmt.Sprintf("\\((%d^{%d})^{%d} = %d^{%d \\times %d} = %d^{%d}\\)", base, e1, e2, base, e1, e2, base, e1*e2),
		Facts:       map[string]string{"base": strconv.Itoa(base), "e1": strconv.Itoa(e1), "e2": strconv.Itoa(e2), "prod_exp": strconv.Itoa(e1 * e2)},
	}
}

type sqrtPerfectGen struct{}

func (g *sqrtPerfectGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	r := rand.Intn(max(1, scale*2)) + 1
	if ctx.Difficulty > 0.5 {
		r = rand.Intn(20) + 1
	}
	return generator.Problem{
		Question:    fmt.Sprintf("\\(\\sqrt{%d} =\\) ?", r*r),
		Answer:      fmt.Sprintf("%d", r),
		Explanation: fmt.Sprintf("\\(\\sqrt{%d} = %d\\) because \\(%d \\times %d = %d\\)", r*r, r, r, r, r*r),
		Facts:       map[string]string{"r": strconv.Itoa(r), "radicand": strconv.Itoa(r * r)},
	}
}

type sqrtSimplifyGen struct{}

// The answer is a coefficient and a radical, e.g. 5 sqrt(2). A numeric keypad
// cannot produce it, which is what the two "formatting" reports were about.
func (g *sqrtSimplifyGen) AnswerFormat() generator.AnswerFormat {
	return generator.AnswerFormat{Hint: "Coefficient and radical, e.g. 5 sqrt(2)", InputMode: "text"}
}

var sqrtPartsRe = regexp.MustCompile(`^(-?\d*)sqrt\(?(\d+)\)?$`)

// sqrtParts reads the forms a learner actually writes: "5 sqrt(2)",
// "5sqrt(2)", "5√2", "5\sqrt{2}" and "sqrt(2)". The grader used to require
// exactly "a sqrt(b)", so a learner who computed a perfect square and typed the
// integer, 12 for sqrt(144), was told they were wrong (reports 138 and 221).
func sqrtParts(s string) (coef, rad int, ok bool) {
	v := strings.ToLower(strings.TrimSpace(s))
	v = strings.NewReplacer("\\", "", "{", "", "}", "", " ", "", "\u221a", "sqrt").Replace(v)
	// MathLive writes a radical as sqrt() while the radicand is being typed, so
	// a learner who tapped the radical template and then the digits submitted
	// "3sqrt()2" and was told their answer was malformed (report 233).
	v = strings.ReplaceAll(v, "sqrt()", "sqrt")
	m := sqrtPartsRe.FindStringSubmatch(v)
	if m == nil {
		return 0, 0, false
	}
	coef = 1
	if m[1] != "" && m[1] != "-" {
		parsed, err := strconv.Atoi(m[1])
		if err != nil {
			return 0, 0, false
		}
		coef = parsed
	} else if m[1] == "-" {
		coef = -1
	}
	rad, err := strconv.Atoi(m[2])
	if err != nil || rad < 2 {
		return 0, 0, false
	}
	return coef, rad, true
}

// squareFreeFactors are the radicands this generator may leave inside the
// radical. Every one is square-free, so the answer cannot be simplified
// further. Before this the remainder was drawn at random, so sqrt(100) was
// served with the answer "5 sqrt(4)": not simplified, and 10, the value a
// learner computes, was rejected as unparseable.
var squareFreeFactors = []int{2, 3, 5, 6, 7, 10, 11, 13}

func (g *sqrtSimplifyGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	perfects := []int{2, 3, 4, 5, 6}
	square := perfects[rand.Intn(len(perfects))]
	pool := minInt(len(squareFreeFactors), 2+scale)
	b := squareFreeFactors[rand.Intn(pool)]
	n := square * square * b
	return generator.Problem{
		Question:    fmt.Sprintf("Simplify: \\(\\sqrt{%d}\\)", n),
		Answer:      fmt.Sprintf("%d sqrt(%d)", square, b),
		Explanation: fmt.Sprintf("\\(\\sqrt{%d} = \\sqrt{%d \\times %d} = \\sqrt{%d} \\times \\sqrt{%d} = %d \\sqrt{%d}\\)", n, square*square, b, square*square, b, square, b),
		Facts:       map[string]string{"n": strconv.Itoa(n), "square": strconv.Itoa(square), "b": strconv.Itoa(b)},
	}
}

func (g *sqrtSimplifyGen) Grade(expected, answer string) grader.Result {
	eCoef, eRad, eOK := sqrtParts(expected)
	aCoef, aRad, aOK := sqrtParts(answer)
	if !eOK || !aOK {
		return grader.Result{Correct: false, Score: 0, Feedback: "Expected format: a sqrt(b) (e.g. 3 sqrt(2))"}
	}
	if eCoef == aCoef && eRad == aRad {
		return grader.Result{Correct: true, Score: 1}
	}
	return grader.Result{Correct: false, Score: 0, Feedback: "Incorrect"}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type negNumberLineGen struct{}

func (g *negNumberLineGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(scale*2) - scale
	// 0 is its own opposite, trivial and confusing; re-roll to non-zero.
	for a == 0 {
		a = rand.Intn(scale*2) - scale
		if scale <= 1 {
			a = rand.Intn(4) - 2
			if a == 0 {
				a = 1
			}
			break
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("What is the opposite of %d?", a),
		Answer:      fmt.Sprintf("%d", -a),
		Explanation: fmt.Sprintf("The opposite of %d is %d.", a, -a),
		Facts:       map[string]string{"a": strconv.Itoa(a), "opposite": strconv.Itoa(-a)},
	}
}

type negAddSubGen struct{}

func (g *negAddSubGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	if ctx.KPOrdinal == 1 {
		scale = 1
	} else if ctx.KPOrdinal == 2 {
		scale = 3
	} else if ctx.KPOrdinal == 3 {
		scale = 5
	}
	a := rand.Intn(scale*4) - scale*2
	b := rand.Intn(scale*2) - scale
	return generator.Problem{
		Question:    fmt.Sprintf("\\(%d + (%d) = ?\\)", a, b),
		Answer:      fmt.Sprintf("%d", a+b),
		Explanation: fmt.Sprintf("%d + (%d) = %d", a, b, a+b),
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "sum": strconv.Itoa(a + b), "b_neg": strconv.Itoa(-b)},
	}
}

type negMultDivGen struct{}

func (g *negMultDivGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(scale*2) - scale
	if a == 0 {
		a = 1
	}
	b := rand.Intn(scale*2) - scale
	if b == 0 {
		b = -1
	}
	return generator.Problem{
		Question:    fmt.Sprintf("\\(%d \\times %d = ?\\)", a, b),
		Answer:      fmt.Sprintf("%d", a*b),
		Explanation: fmt.Sprintf("(%d) x (%d) = %d", a, b, a*b),
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "product": strconv.Itoa(a * b)},
	}
}

type orderOpsGen struct {
	parens    bool
	exponents bool
}

func (g *orderOpsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	if ctx.KPOrdinal == 1 {
		scale = 1
	} else if ctx.KPOrdinal == 2 {
		scale = 3
	} else if ctx.KPOrdinal == 3 {
		scale = 5
	}
	a := rand.Intn(scale*2) + 1
	b := rand.Intn(scale*2) + 1
	c := rand.Intn(max(1, scale)) + 1
	var q string
	var result int
	// e and shape are hoisted so one authored schema can cover all three
	// shapes; e is 0 when the drawn expression has no power in it.
	e := 0
	shape := "flat"
	if !g.parens && !g.exponents {
		// 3 + 4 x 2
		q = fmt.Sprintf("%d + %d \\times %d", a, b, c)
		result = a + b*c
	} else if g.parens && !g.exponents {
		shape = "brackets"
		q = fmt.Sprintf("(%d + %d) \\times %d", a, b, c)
		result = (a + b) * c
	} else {
		shape = "nested"
		e = rand.Intn(3) + 2
		q = fmt.Sprintf("%d + (%d)^%d \\times %d", a, b, e, c)
		result = a + mathutil.IntPow(b, e)*c
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Evaluate: \\(%s\\)", q),
		Answer:      fmt.Sprintf("%d", result),
		Explanation: fmt.Sprintf("%s = %d (PEMDAS)", q, result),
		Facts:       map[string]string{"q": q, "result": strconv.Itoa(result), "a": strconv.Itoa(a), "b": strconv.Itoa(b), "c": strconv.Itoa(c), "e": strconv.Itoa(e), "shape": shape},
	}
}

type decIntroGen struct{}

func (g *decIntroGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	ones := rand.Intn(max(1, scale)) + 1
	// Tenths digit is never zero: a mixed number needs a real fractional
	// part ("Write 2.0 as a mixed number" with answer "2 0/10" is degenerate).
	tenths := rand.Intn(max(1, scale*2)) + 1
	if tenths > 9 {
		tenths = 9
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Write %d.%d as a mixed number.", ones, tenths),
		Answer:      fmt.Sprintf("%d %d/10", ones, tenths),
		Explanation: fmt.Sprintf("\\(%d.%d = %d\\ \\text{and}\\ \\frac{%d}{10}\\)", ones, tenths, ones, tenths),
		Facts:       map[string]string{"ones": strconv.Itoa(ones), "tenths": strconv.Itoa(tenths)},
	}
}
