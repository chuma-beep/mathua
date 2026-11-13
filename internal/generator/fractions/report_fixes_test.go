package fractions

import (
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

// Report 462: the question printed 4/10, 3/10 and 4/10, so "which is largest"
// had two correct positions and the stored answer was arbitrary.
func TestDecimalOrderUsesDistinctTenths(t *testing.T) {
	re := regexp.MustCompile(`Which of (\d+)/10, (\d+)/10 and (\d+)/10 is largest`)
	for i := 0; i < 300; i++ {
		p := w10DecimalOrder(2)
		m := re.FindStringSubmatch(p.Question)
		if m == nil {
			t.Fatalf("unexpected decimal order question: %q", p.Question)
		}
		vals := make([]int, 3)
		for j := range vals {
			vals[j], _ = strconv.Atoi(m[j+1])
		}
		if vals[0] == vals[1] || vals[0] == vals[2] || vals[1] == vals[2] {
			t.Fatalf("repeated tenths in %q", p.Question)
		}
		want := 1
		if vals[1] > vals[want-1] {
			want = 2
		}
		if vals[2] > vals[want-1] {
			want = 3
		}
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("%q answered %s, want position %d", p.Question, p.Answer, want)
		}
	}
}

// Reports 454 and 460: 57/100 and 1 were marked wrong because the generator
// counted every fraction above 1/2 instead of every fraction above 3/4.
func TestFractionEstimateCountsCloserToOne(t *testing.T) {
	g := &fracBridgeGen{kind: "estimate"}
	re := regexp.MustCompile(`closer to 1 than to 1/2: (\d+)/100, (\d+)/100, (\d+)/100 and (\d+)/100`)
	for i := 0; i < 500; i++ {
		p := g.Generate(generator.GeneratorContext{KPOrdinal: 3})
		m := re.FindStringSubmatch(p.Question)
		if m == nil {
			t.Fatalf("unexpected estimate question: %q", p.Question)
		}
		want := 0
		for j := 1; j <= 4; j++ {
			v, _ := strconv.Atoi(m[j])
			if 4*v > 300 {
				want++
			}
		}
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("%q answered %s, want %d", p.Question, p.Answer, want)
		}
		if !strings.Contains(p.Explanation, "3/4") {
			t.Fatalf("explanation does not state the midpoint: %q", p.Explanation)
		}
	}
}
