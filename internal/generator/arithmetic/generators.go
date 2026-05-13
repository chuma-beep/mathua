package arithmetic

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
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
}

type addGen struct {
	minA, maxA  int
	minB, maxB  int
	ensureCarry bool
}

func (g *addGen) Generate(difficulty float64) generator.Problem {
	scale := int(1 + difficulty*float64(g.maxA-g.minA)/10)
	scale = max(scale, 1)
	a := rand.Intn(min(g.maxA-g.minA, scale*10)) + g.minA
	b := rand.Intn(min(g.maxB-g.minB, scale*10)) + g.minB
	if g.ensureCarry {
		a = rand.Intn(90) + 10
		onesA := a % 10
		onesB := 10 - onesA + rand.Intn(9)
		if onesB <= 9 {
			b = (b/10)*10 + onesB
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("%d + %d = ?", a, b),
		Answer:      fmt.Sprintf("%d", a+b),
		Explanation: fmt.Sprintf("%d + %d = %d", a, b, a+b),
	}
}

type addWordGen struct{}

func (g *addWordGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(50) + 1
	b := rand.Intn(50) + 1
	items := []string{"apples", "marbles", "stickers", "crayons", "pencils"}
	item := items[rand.Intn(len(items))]
	return generator.Problem{
		Question:    fmt.Sprintf("You have %d %s. Your friend gives you %d more. How many do you have now?", a, item, b),
		Answer:      fmt.Sprintf("%d", a+b),
		Explanation: fmt.Sprintf("%d + %d = %d %s", a, b, a+b, item),
	}
}

type subGen struct {
	minA, maxA int
	minB, maxB int
}

func (g *subGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(g.maxA-g.minA+1) + g.minA
	bMax := a - g.minB
	if bMax < g.minB {
		bMax = a
	}
	b := rand.Intn(bMax-g.minB+1) + g.minB
	return generator.Problem{
		Question:    fmt.Sprintf("%d - %d = ?", a, b),
		Answer:      fmt.Sprintf("%d", a-b),
		Explanation: fmt.Sprintf("%d - %d = %d", a, b, a-b),
	}
}

type subBorrowGen struct{}

func (g *subBorrowGen) Generate(difficulty float64) generator.Problem {
	tensA := rand.Intn(9) + 1
	onesA := rand.Intn(10)
	onesB := onesA + rand.Intn(10-onesA) + 1
	tensB := rand.Intn(tensA)
	if tensB == tensA {
		tensB--
	}
	a := tensA*10 + onesA
	b := tensB*10 + onesB
	return generator.Problem{
		Question:    fmt.Sprintf("%d - %d = ?", a, b),
		Answer:      fmt.Sprintf("%d", a-b),
		Explanation: fmt.Sprintf("%d - %d = %d (borrowing required)", a, b, a-b),
	}
}

type subWordGen struct{}

func (g *subWordGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(80) + 20
	b := rand.Intn(a-1) + 1
	items := []string{"candies", "pencils", "cards", "coins", "beads"}
	item := items[rand.Intn(len(items))]
	return generator.Problem{
		Question:    fmt.Sprintf("You have %d %s. You give away %d. How many are left?", a, item, b),
		Answer:      fmt.Sprintf("%d", a-b),
		Explanation: fmt.Sprintf("%d - %d = %d %s left", a, b, a-b, item),
	}
}

type placeValueGen struct {
	max   int
	label string
}

func (g *placeValueGen) Generate(difficulty float64) generator.Problem {
	n := rand.Intn(g.max) + 1
	var placeUnit int
	switch g.label {
	case "tens":
		placeUnit = 10
	case "hundreds":
		placeUnit = 100
	case "thousands":
		placeUnit = 1000
	}
	var digit, value int
	switch g.label {
	case "tens":
		digit = (n / 10) % 10
		value = digit * 10
	case "hundreds":
		digit = (n / 100) % 10
		value = digit * 100
	case "thousands":
		digit = (n / 1000) % 10
		value = digit * 1000
	}
	return generator.Problem{
		Question:    fmt.Sprintf("In the number %d, what is the value of the digit in the %s place?", n, g.label),
		Answer:      fmt.Sprintf("%d", value),
		Explanation: fmt.Sprintf("The %s digit is %d, value = %d x %d = %d", g.label, digit, digit, placeUnit, value),
	}
}

type roundGen struct {
	to int
}

func (g *roundGen) Generate(difficulty float64) generator.Problem {
	lo := g.to
	hi := g.to * 10
	n := rand.Intn(hi-lo) + lo
	rounded := int(math.Round(float64(n)/float64(g.to))) * g.to
	return generator.Problem{
		Question:    fmt.Sprintf("Round %d to the nearest %d.", n, g.to),
		Answer:      fmt.Sprintf("%d", rounded),
		Explanation: fmt.Sprintf("%d rounded to nearest %d is %d.", n, g.to, rounded),
	}
}

type multConceptGen struct{}

func (g *multConceptGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(5) + 2
	b := rand.Intn(5) + 2
	return generator.Problem{
		Question:    fmt.Sprintf("%d groups of %d = ?", a, b),
		Answer:      fmt.Sprintf("%d", a*b),
		Explanation: fmt.Sprintf("%d groups of %d means %d x %d = %d", a, b, a, b, a*b),
	}
}

type multByGen struct {
	factor int
}

func (g *multByGen) Generate(difficulty float64) generator.Problem {
	f := []int{2, 5, 10}[rand.Intn(3)]
	n := rand.Intn(12) + 1
	if difficulty > 0.5 {
		n = rand.Intn(20) + 1
	}
	return generator.Problem{
		Question:    fmt.Sprintf("%d x %d = ?", f, n),
		Answer:      fmt.Sprintf("%d", f*n),
		Explanation: fmt.Sprintf("%d x %d = %d", f, n, f*n),
	}
}

type multTablesGen struct{}

func (g *multTablesGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(12) + 1
	b := rand.Intn(12) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("%d x %d = ?", a, b),
		Answer:      fmt.Sprintf("%d", a*b),
		Explanation: fmt.Sprintf("%d x %d = %d", a, b, a*b),
	}
}

type multGen struct {
	minA, maxA int
	minB, maxB int
}

func (g *multGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(g.maxA-g.minA+1) + g.minA
	b := rand.Intn(g.maxB-g.minB+1) + g.minB
	return generator.Problem{
		Question:    fmt.Sprintf("%d x %d = ?", a, b),
		Answer:      fmt.Sprintf("%d", a*b),
		Explanation: fmt.Sprintf("%d x %d = %d", a, b, a*b),
	}
}

type multWordGen struct{}

func (g *multWordGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(12) + 1
	b := rand.Intn(12) + 1
	items := []string{"stickers", "cards", "beads", "coins", "marbles"}
	item := items[rand.Intn(len(items))]
	return generator.Problem{
		Question:    fmt.Sprintf("You have %d bags with %d %s each. How many %s in total?", a, b, item, item),
		Answer:      fmt.Sprintf("%d", a*b),
		Explanation: fmt.Sprintf("%d x %d = %d %s", a, b, a*b, item),
	}
}

type divConceptGen struct{}

func (g *divConceptGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(5) + 2
	b := rand.Intn(5) + 2
	total := a * b
	return generator.Problem{
		Question:    fmt.Sprintf("You have %d items. You split them into %d equal groups. How many in each group?", total, a),
		Answer:      fmt.Sprintf("%d", b),
		Explanation: fmt.Sprintf("%d / %d = %d in each group", total, a, b),
	}
}

type divBasicGen struct{}

func (g *divBasicGen) Generate(difficulty float64) generator.Problem {
	b := rand.Intn(12) + 1
	a := b * (rand.Intn(12) + 1)
	return generator.Problem{
		Question:    fmt.Sprintf("%d / %d = ?", a, b),
		Answer:      fmt.Sprintf("%d", a/b),
		Explanation: fmt.Sprintf("%d / %d = %d", a, b, a/b),
	}
}

type divRemainderGen struct{}

func (g *divRemainderGen) Generate(difficulty float64) generator.Problem {
	b := rand.Intn(10) + 2
	r := rand.Intn(b-1) + 1
	a := b*(rand.Intn(10)+1) + r
	return generator.Problem{
		Question:    fmt.Sprintf("%d / %d = ? (give answer with remainder: Q R)", a, b),
		Answer:      fmt.Sprintf("%d R %d", a/b, r),
		Explanation: fmt.Sprintf("%d / %d = %d remainder %d", a, b, a/b, r),
	}
}

type divLongGen struct{}

func (g *divLongGen) Generate(difficulty float64) generator.Problem {
	b := rand.Intn(12) + 2
	q := rand.Intn(50) + 10
	a := b * q
	return generator.Problem{
		Question:    fmt.Sprintf("%d / %d = ?", a, b),
		Answer:      fmt.Sprintf("%d", q),
		Explanation: fmt.Sprintf("%d / %d = %d", a, b, q),
	}
}

type divWordGen struct{}

func (g *divWordGen) Generate(difficulty float64) generator.Problem {
	b := rand.Intn(10) + 2
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
	}
}

type factorFindGen struct{}

func (g *factorFindGen) Generate(difficulty float64) generator.Problem {
	n := rand.Intn(50) + 10
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
	}
}

type primeGen struct{}

func (g *primeGen) Generate(difficulty float64) generator.Problem {
	primes := []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47}
	composites := []int{4, 6, 8, 9, 10, 12, 14, 15, 16, 18, 20, 21, 22, 24, 25}
	if rand.Intn(2) == 0 {
		p := primes[rand.Intn(len(primes))]
		return generator.Problem{
			Question:    fmt.Sprintf("Is %d prime?", p),
			Answer:      "yes",
			Explanation: fmt.Sprintf("%d is prime — only divisible by 1 and itself.", p),
		}
	}
	c := composites[rand.Intn(len(composites))]
	return generator.Problem{
		Question:    fmt.Sprintf("Is %d prime?", c),
		Answer:      "no",
		Explanation: fmt.Sprintf("%d is composite.", c),
	}
}

type primeFactGen struct{}

func (g *primeFactGen) Generate(difficulty float64) generator.Problem {
	n := rand.Intn(60) + 2
	factors := primeFactors(n)
	fStr := make([]string, len(factors))
	for i, f := range factors {
		fStr[i] = strconv.Itoa(f)
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Find the prime factorization of %d.", n),
		Answer:      strings.Join(fStr, ","),
		Explanation: fmt.Sprintf("%d = %s", n, strings.Join(fStr, " x ")),
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

func (g *gcfGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(50) + 10
	b := rand.Intn(50) + 10
	result := mathutil.GCD(a, b)
	return generator.Problem{
		Question:    fmt.Sprintf("Find the GCF of %d and %d.", a, b),
		Answer:      fmt.Sprintf("%d", result),
		Explanation: fmt.Sprintf("GCF(%d, %d) = %d", a, b, result),
	}
}

type lcmGen struct{}

func (g *lcmGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(20) + 2
	b := rand.Intn(20) + 2
	l := a * b / mathutil.GCD(a, b)
	return generator.Problem{
		Question:    fmt.Sprintf("Find the LCM of %d and %d.", a, b),
		Answer:      fmt.Sprintf("%d", l),
		Explanation: fmt.Sprintf("LCM(%d, %d) = %d", a, b, l),
	}
}

type compositeGen struct{}

func (g *compositeGen) Generate(difficulty float64) generator.Problem {
	n := rand.Intn(50) + 4
	isComposite := false
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			isComposite = true
			break
		}
	}
	ans := "no"
	exp := fmt.Sprintf("%d is prime.", n)
	if isComposite {
		ans = "yes"
		exp = fmt.Sprintf("%d is composite.", n)
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Is %d composite?", n),
		Answer:      ans,
		Explanation: exp,
	}
}

type expConceptGen struct{}

func (g *expConceptGen) Generate(difficulty float64) generator.Problem {
	base := rand.Intn(5) + 2
	exp := rand.Intn(3) + 2
	return generator.Problem{
		Question:    fmt.Sprintf("What does %d^%d mean?", base, exp),
		Answer:      fmt.Sprintf("%d", mathutil.IntPow(base, exp)),
		Explanation: fmt.Sprintf("%d^%d = %s = %d", base, exp, strings.Repeat(fmt.Sprintf("%dx", base), exp-1)+fmt.Sprintf("%d", base), mathutil.IntPow(base, exp)),
	}
}

type expEvalGen struct{}

func (g *expEvalGen) Generate(difficulty float64) generator.Problem {
	base := rand.Intn(9) + 2
	exp := rand.Intn(5) + 1
	if difficulty > 0.5 {
		exp = rand.Intn(5) + 3
	}
	return generator.Problem{
		Question:    fmt.Sprintf("%d^%d = ?", base, exp),
		Answer:      fmt.Sprintf("%d", mathutil.IntPow(base, exp)),
		Explanation: fmt.Sprintf("%d^%d = %d", base, exp, mathutil.IntPow(base, exp)),
	}
}

type expProductRuleGen struct{}

func (g *expProductRuleGen) Generate(difficulty float64) generator.Problem {
	base := rand.Intn(4) + 2
	e1 := rand.Intn(4) + 1
	e2 := rand.Intn(4) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Simplify: %d^%d x %d^%d", base, e1, base, e2),
		Answer:      fmt.Sprintf("%d^%d", base, e1+e2),
		Explanation: fmt.Sprintf("%d^%d x %d^%d = %d^(%d+%d) = %d^%d", base, e1, base, e2, base, e1, e2, base, e1+e2),
	}
}

type expQuotientRuleGen struct{}

func (g *expQuotientRuleGen) Generate(difficulty float64) generator.Problem {
	base := rand.Intn(4) + 2
	e1 := rand.Intn(4) + 3
	e2 := rand.Intn(3) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Simplify: %d^%d / %d^%d", base, e1, base, e2),
		Answer:      fmt.Sprintf("%d^%d", base, e1-e2),
		Explanation: fmt.Sprintf("%d^%d / %d^%d = %d^(%d-%d) = %d^%d", base, e1, base, e2, base, e1, e2, base, e1-e2),
	}
}

type expPowerRuleGen struct{}

func (g *expPowerRuleGen) Generate(difficulty float64) generator.Problem {
	base := rand.Intn(4) + 2
	e1 := rand.Intn(4) + 1
	e2 := rand.Intn(4) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Simplify: (%d^%d)^%d", base, e1, e2),
		Answer:      fmt.Sprintf("%d^%d", base, e1*e2),
		Explanation: fmt.Sprintf("(%d^%d)^%d = %d^(%dx%d) = %d^%d", base, e1, e2, base, e1, e2, base, e1*e2),
	}
}

type sqrtPerfectGen struct{}

func (g *sqrtPerfectGen) Generate(difficulty float64) generator.Problem {
	r := rand.Intn(13) + 1
	if difficulty > 0.5 {
		r = rand.Intn(20) + 1
	}
	return generator.Problem{
		Question:    fmt.Sprintf("sqrt(%d) = ?", r*r),
		Answer:      fmt.Sprintf("%d", r),
		Explanation: fmt.Sprintf("sqrt(%d) = %d because %d x %d = %d", r*r, r, r, r, r*r),
	}
}

type sqrtSimplifyGen struct{}

func (g *sqrtSimplifyGen) Generate(difficulty float64) generator.Problem {
	// sqrt(n) = a*sqrt(b) where n = a^2 * b
	perfects := []int{2, 3, 4, 5, 6}
	square := perfects[rand.Intn(len(perfects))]
	b := rand.Intn(7) + 2
	n := square * square * b
	return generator.Problem{
		Question:    fmt.Sprintf("Simplify: sqrt(%d)", n),
		Answer:      fmt.Sprintf("%d sqrt(%d)", square, b),
		Explanation: fmt.Sprintf("sqrt(%d) = sqrt(%dx%d) = sqrt(%d) x sqrt(%d) = %d sqrt(%d)", n, square*square, b, square*square, b, square, b),
	}
}

type negNumberLineGen struct{}

func (g *negNumberLineGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(10) - 5
	return generator.Problem{
		Question:    fmt.Sprintf("What is the opposite of %d?", a),
		Answer:      fmt.Sprintf("%d", -a),
		Explanation: fmt.Sprintf("The opposite of %d is %d.", a, -a),
	}
}

type negAddSubGen struct{}

func (g *negAddSubGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(20) - 10
	b := rand.Intn(10) - 5
	return generator.Problem{
		Question:    fmt.Sprintf("%d + (%d) = ?", a, b),
		Answer:      fmt.Sprintf("%d", a+b),
		Explanation: fmt.Sprintf("%d + (%d) = %d", a, b, a+b),
	}
}

type negMultDivGen struct{}

func (g *negMultDivGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(12) - 6
	if a == 0 {
		a = 1
	}
	b := rand.Intn(12) - 6
	if b == 0 {
		b = -1
	}
	return generator.Problem{
		Question:    fmt.Sprintf("%d x %d = ?", a, b),
		Answer:      fmt.Sprintf("%d", a*b),
		Explanation: fmt.Sprintf("(%d) x (%d) = %d", a, b, a*b),
	}
}

type orderOpsGen struct {
	parens    bool
	exponents bool
}

func (g *orderOpsGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(9) + 1
	b := rand.Intn(9) + 1
	c := rand.Intn(5) + 1
	var q string
	var result int
	if !g.parens && !g.exponents {
		// 3 + 4 x 2
		q = fmt.Sprintf("%d + %d x %d", a, b, c)
		result = a + b*c
	} else if g.parens && !g.exponents {
		q = fmt.Sprintf("(%d + %d) x %d", a, b, c)
		result = (a + b) * c
	} else {
		e := rand.Intn(3) + 2
		q = fmt.Sprintf("%d + (%d)^%d x %d", a, b, e, c)
		result = a + mathutil.IntPow(b, e)*c
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Evaluate: %s", q),
		Answer:      fmt.Sprintf("%d", result),
		Explanation: fmt.Sprintf("%s = %d (PEMDAS)", q, result),
	}
}

type decIntroGen struct{}

func (g *decIntroGen) Generate(difficulty float64) generator.Problem {
	ones := rand.Intn(9) + 1
	tenths := rand.Intn(10)
	return generator.Problem{
		Question:    fmt.Sprintf("Write %d.%d as a mixed number.", ones, tenths),
		Answer:      fmt.Sprintf("%d %d/10", ones, tenths),
		Explanation: fmt.Sprintf("%d.%d = %d and %d/10", ones, tenths, ones, tenths),
	}
}
