package arithmetic

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

func TestEstimateQuotientCaseOneRoundsDividendToHundred(t *testing.T) {
	g := &arithBridgeGen{kind: "estimate_quotient"}
	re := regexp.MustCompile(`Estimate (\d+) divided by (\d+)`)
	for i := 0; i < 200; i++ {
		p := g.Generate(generator.GeneratorContext{KPOrdinal: 1})
		m := re.FindStringSubmatch(p.Question)
		if m == nil {
			t.Fatalf("question does not match estimate shape: %q", p.Question)
		}
		a, _ := strconv.Atoi(m[1])
		b, _ := strconv.Atoi(m[2])
		if a < 100 {
			t.Fatalf("dividend %d cannot be honestly rounded to a non-zero hundred", a)
		}
		a0 := 100 * ((a + 50) / 100)
		b0 := 10 * ((b + 5) / 10)
		want := (a0 + b0/2) / b0
		got, err := strconv.Atoi(p.Answer)
		if err != nil {
			t.Fatalf("answer %q is not an integer", p.Answer)
		}
		if got != want {
			t.Fatalf("%s = %d, want %d from %d/%d", p.Question, got, want, a0, b0)
		}
		if !strings.Contains(p.Question, "nearest hundred") {
			t.Fatalf("prompt does not state nearest hundred: %q", p.Question)
		}
		if !strings.Contains(p.Explanation, fmt.Sprintf("%d rounds to %d", a, a0)) {
			t.Fatalf("explanation does not use rounded dividend %d: %q", a0, p.Explanation)
		}
		if !strings.Contains(p.Explanation, "about") {
			t.Fatalf("explanation still claims an exact quotient: %q", p.Explanation)
		}
	}
}

// Report 458: 508/26 rounds to 500/30, which is about 17, not 16. The old
// generator truncated the division and called the result exact.
func TestEstimateQuotientRoundsToNearest(t *testing.T) {
	g := &arithBridgeGen{kind: "estimate_quotient"}
	re := regexp.MustCompile(`Estimate (\d+) divided by (\d+)`)
	for kp := 1; kp <= 2; kp++ {
		for i := 0; i < 200; i++ {
			p := g.Generate(generator.GeneratorContext{KPOrdinal: kp})
			m := re.FindStringSubmatch(p.Question)
			if m == nil {
				t.Fatalf("unexpected question: %q", p.Question)
			}
			a, _ := strconv.Atoi(m[1])
			b, _ := strconv.Atoi(m[2])
			a0 := 10 * ((a + 5) / 10)
			b0 := 10 * ((b + 5) / 10)
			if kp == 1 {
				a0 = 100 * ((a + 50) / 100)
			}
			want := (a0 + b0/2) / b0
			if p.Answer != strconv.Itoa(want) {
				t.Fatalf("kp%d %q answered %s, want %d", kp, p.Question, p.Answer, want)
			}
		}
	}
}

// Report 452: the explanation used T_n without saying what a triangular number
// is, so the one term the report asked about was never defined.
func TestTriangularExplanationDefinesTerm(t *testing.T) {
	for i := 0; i < 100; i++ {
		p := w8TriangularNumbers(1)
		if !strings.Contains(p.Explanation, "count the dots") {
			t.Fatalf("triangular explanation does not define the term: %q", p.Explanation)
		}
	}
}

// Report 433: 63 + 30 was served under arith.add.carry although the ones
// digits sum to 3, so there was nothing to carry.
func TestAddCarryGenAlwaysCarries(t *testing.T) {
	g := &addGen{minA: 10, maxA: 99, minB: 10, maxB: 99, ensureCarry: true}
	for i := 0; i < 500; i++ {
		p := g.Generate(generator.GeneratorContext{Difficulty: rand.Float64()})
		a, _ := strconv.Atoi(p.Facts["a"])
		b, _ := strconv.Atoi(p.Facts["b"])
		if a%10 == 0 || b%10 == 0 || a%10+b%10 < 10 {
			t.Fatalf("%d + %d does not carry (ones %d + %d)", a, b, a%10, b%10)
		}
	}
}

// Reports 434 and 435: consecutive draws landed on pairs of primes, so the
// answer was 1 and the learner never found a factor.
func TestGcfGenAlwaysSharesAFactor(t *testing.T) {
	g := &gcfGen{}
	for i := 0; i < 500; i++ {
		p := g.Generate(generator.GeneratorContext{Difficulty: rand.Float64()})
		a, _ := strconv.Atoi(p.Facts["a"])
		b, _ := strconv.Atoi(p.Facts["b"])
		if a == b {
			t.Fatalf("GCF of %d and %d repeats the same number", a, b)
		}
		if gcd := mathutil.GCD(a, b); gcd < 2 {
			t.Fatalf("GCF of %d and %d is %d, want a shared factor", a, b, gcd)
		}
	}
}

// Reports 436 and 437: a multiplier of 1 turned the multiplication into a
// no-op, and the nested base was allowed to be 1 as well.
func TestOrderOpsGenUsesRealOperands(t *testing.T) {
	for _, g := range []*orderOpsGen{{parens: true}, {parens: true, exponents: true}} {
		for i := 0; i < 300; i++ {
			p := g.Generate(generator.GeneratorContext{Difficulty: rand.Float64()})
			c, _ := strconv.Atoi(p.Facts["c"])
			if c < 2 {
				t.Fatalf("multiplier %d in %q", c, p.Question)
			}
			if g.exponents {
				b, _ := strconv.Atoi(p.Facts["b"])
				if b < 2 {
					t.Fatalf("exponent base %d in %q", b, p.Question)
				}
			} else if strings.Contains(p.Facts["q"], "(1 + 1)") {
				t.Fatalf("trivial parenthesized sum in %q", p.Question)
			}
		}
	}
}

// Report 445: the count question listed primes and composites but never said
// what either word means, so the per-number notes were the only clue.
func TestCompositeExplanationDefinesComposite(t *testing.T) {
	for i := 0; i < 100; i++ {
		p := compositeProblem()
		if !strings.Contains(p.Explanation, "more than two divisors") {
			t.Fatalf("composite explanation does not define the term: %q", p.Explanation)
		}
		if !strings.Contains(p.Explanation, "exactly two divisors") {
			t.Fatalf("composite explanation does not define prime: %q", p.Explanation)
		}
	}
}

// Report 493: 70×84 was read as 70×80 because the prompt did not say that
// only 71 is split and 84 stays whole.
func TestPartialProductTensPromptKeepsSecondFactorWhole(t *testing.T) {
	g := &w8ArithGen{kind: "partial_products"}
	for i := 0; i < 100; i++ {
		p := g.Generate(generator.GeneratorContext{KPOrdinal: 1})
		a, _ := strconv.Atoi(p.Facts["a"])
		b, _ := strconv.Atoi(p.Facts["b"])
		if !strings.Contains(p.Question, "split only") {
			t.Fatalf("prompt does not say only the first factor splits: %q", p.Question)
		}
		if !strings.Contains(p.Explanation, "stays whole") {
			t.Fatalf("explanation does not keep the second factor whole: %q", p.Explanation)
		}
		want := (a / 10 * 10) * b
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("%q answered %s, want %d", p.Question, p.Answer, want)
		}
	}
}

// Report 492: a wrong divisor count got only the pairing rule, so the learner
// had no list to check their answer against.
func TestDivisorCountExplanationListsDivisors(t *testing.T) {
	g := &w12Gen{kind: "divisor_count"}
	for kp := 1; kp <= 3; kp++ {
		for i := 0; i < 60; i++ {
			p := g.Generate(generator.GeneratorContext{KPOrdinal: kp})
			n, _ := strconv.Atoi(p.Facts["n"])
			var divisors []string
			for d := 1; d <= n; d++ {
				if n%d == 0 {
					divisors = append(divisors, strconv.Itoa(d))
				}
			}
			if p.Answer != strconv.Itoa(len(divisors)) {
				t.Fatalf("%d answered %s, want %d", n, p.Answer, len(divisors))
			}
			want := strings.Join(divisors[:len(divisors)-1], ", ") + " and " + divisors[len(divisors)-1]
			if !strings.Contains(p.Explanation, want) {
				t.Fatalf("explanation does not list the divisors of %d as %q: %q", n, want, p.Explanation)
			}
		}
	}
}

// Report 501: the pool held exactly eight primes, so a draw that took none of
// the composites was the same eight numbers every time. Keep both kinds in
// every draw so the set varies.
func TestCompositeCountKeepsBothKindsInTheDraw(t *testing.T) {
	for i := 0; i < 300; i++ {
		p := compositeProblem()
		n, err := strconv.Atoi(p.Answer)
		if err != nil {
			t.Fatalf("answer %q is not a count", p.Answer)
		}
		if n == 0 || n == 8 {
			t.Fatalf("draw has no mix of primes and composites: answer %d, question %q", n, p.Question)
		}
	}
}

// Report 502: the third ordinal drew a four-digit factor while the lesson is
// titled "Multiplying two three-digit numbers".
func TestThreeByThreeStaysThreeDigits(t *testing.T) {
	g := &w12Gen{kind: "mult_three_three"}
	for _, kp := range []int{1, 2, 3} {
		for i := 0; i < 100; i++ {
			p := g.Generate(generator.GeneratorContext{KPOrdinal: kp})
			a, errA := strconv.Atoi(p.Facts["a"])
			b, errB := strconv.Atoi(p.Facts["b"])
			if errA != nil || errB != nil {
				t.Fatalf("factors are missing from facts: %+v", p.Facts)
			}
			if a < 100 || a > 999 || b < 100 || b > 999 {
				t.Fatalf("kp%d drew %d x %d, which is not two three-digit numbers", kp, a, b)
			}
			if p.Answer != strconv.Itoa(a*b) {
				t.Fatalf("answer %q is not %d x %d", p.Answer, a, b)
			}
		}
	}
}

func TestOrderOfOperationsNegativesNamesTheRule(t *testing.T) {
	// Report 506: the learner asked where the brackets are. There are none by
	// design, because precedence is what the lesson teaches, so the prompt has
	// to name the rule instead of leaving it implicit.
	g := &w12Gen{kind: "order_negatives"}
	re := regexp.MustCompile(`What is -(\d+) \+ (\d+) x (\d+) - (\d+)\?`)
	for i := 0; i < 100; i++ {
		p := g.Generate(generator.GeneratorContext{KPOrdinal: 1})
		if !strings.Contains(p.Question, "order of operations") {
			t.Fatalf("prompt does not name the rule: %q", p.Question)
		}
		m := re.FindStringSubmatch(p.Question)
		if m == nil {
			t.Fatalf("question does not match the shape: %q", p.Question)
		}
		nums := make([]int, 4)
		for j := range nums {
			nums[j], _ = strconv.Atoi(m[j+1])
		}
		want := -nums[0] + nums[1]*nums[2] - nums[3]
		got, err := strconv.Atoi(p.Answer)
		if err != nil {
			t.Fatalf("answer %q is not an integer", p.Answer)
		}
		if got != want {
			t.Fatalf("%s = %d, want %d", p.Question, got, want)
		}
	}
}

func TestEstimateCompatibleShowsTheCancelledDivision(t *testing.T) {
	// Report 507: the explanation said "since 2040 / 30 = 68", which is the
	// step the learner could not do. It has to show the zeros coming out.
	g := &w12Gen{kind: "estimate_compatible"}
	re := regexp.MustCompile(`Estimate (\d+) / (\d+) by using the compatible number (\d+)`)
	for _, kp := range []int{1, 2, 3} {
		for i := 0; i < 100; i++ {
			p := g.Generate(generator.GeneratorContext{KPOrdinal: kp})
			m := re.FindStringSubmatch(p.Question)
			if m == nil {
				t.Fatalf("question does not match the shape: %q", p.Question)
			}
			dividend, _ := strconv.Atoi(m[1])
			friendly, _ := strconv.Atoi(m[3])
			scale := 10
			if friendly >= 100 {
				scale = 100
			}
			want := dividend / friendly
			got, err := strconv.Atoi(p.Answer)
			if err != nil {
				t.Fatalf("answer %q is not an integer", p.Answer)
			}
			if got != want {
				t.Fatalf("%s = %d, want %d", p.Question, got, want)
			}
			cancelled := fmt.Sprintf("%d / %d", dividend/scale, friendly/scale)
			if !strings.Contains(p.Explanation, cancelled) {
				t.Fatalf("explanation hides the mental step %q: %q", cancelled, p.Explanation)
			}
		}
	}
}

func TestCompensationAlwaysLeavesSomethingToSubtract(t *testing.T) {
	// Report 508: when the drawn addend was 10 the prompt read "add 10 ... first
	// add 10", the answer was zero and there was nothing to compensate.
	g := &w8ArithRemainingGen{kind: "compensation"}
	re := regexp.MustCompile(`To add (\d+) to (\d+) by compensation`)
	for i := 0; i < 300; i++ {
		p := g.Generate(generator.GeneratorContext{KPOrdinal: 1})
		m := re.FindStringSubmatch(p.Question)
		if m == nil {
			t.Fatalf("question does not match the shape: %q", p.Question)
		}
		addend, _ := strconv.Atoi(m[1])
		if addend >= 10 {
			t.Fatalf("adding %d leaves nothing to compensate: %q", addend, p.Question)
		}
		want := 10 - addend
		got, err := strconv.Atoi(p.Answer)
		if err != nil {
			t.Fatalf("answer %q is not an integer", p.Answer)
		}
		if got != want || got < 1 || got > 3 {
			t.Fatalf("%s = %d, want %d between 1 and 3", p.Question, got, want)
		}
	}
}
