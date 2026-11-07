package arithmetic

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
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
