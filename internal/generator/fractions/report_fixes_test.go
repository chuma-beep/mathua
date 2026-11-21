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

// Reports 495-497: anything below one half was labelled closest to 0, so 3/8
// was marked as closer to 0 than to 1/2.
func TestEstimateProductUsesNearestBenchmark(t *testing.T) {
	re := regexp.MustCompile(`Is (\d+)/(\d+) times (\d+)/(\d+) closest to 0, 1/2 or 1`)
	for i := 0; i < 500; i++ {
		p := w10EstimateProduct(1)
		m := re.FindStringSubmatch(p.Question)
		if m == nil {
			t.Fatalf("unexpected estimate product question: %q", p.Question)
		}
		a, _ := strconv.Atoi(m[1])
		b, _ := strconv.Atoi(m[2])
		c, _ := strconv.Atoi(m[3])
		d, _ := strconv.Atoi(m[4])
		num := a * c
		den := b * d
		if 4*num == den || 4*num == 3*den {
			t.Fatalf("%q lands on a tie between benchmarks", p.Question)
		}
		want := 0
		if 4*num > den {
			want = 1
		}
		if 4*num > 3*den {
			want = 2
		}
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("%q answered %s, want %d", p.Question, p.Answer, want)
		}
	}
}

// Report 489: 114/100 was rejected when the prompt asked for the numerator
// 114, though 114/100 is the sum the learner was told to write over 100.
func TestAddDecimalAcceptsFractionOverHundred(t *testing.T) {
	g := &w10Gen{kind: "add_decimal"}
	for i := 0; i < 200; i++ {
		p := g.Generate(generatorContext(1))
		if res := g.Grade(p.Answer, p.Answer+"/100"); !res.Correct {
			t.Fatalf("%q rejected the fraction %s/%d", p.Question, p.Answer, 100)
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
