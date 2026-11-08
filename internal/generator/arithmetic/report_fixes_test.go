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
		want := a0 / b0
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
