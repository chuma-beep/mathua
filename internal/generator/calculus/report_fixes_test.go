package calculus

import (
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func TestReducedFraction(t *testing.T) {
	cases := map[[2]int]string{
		{2, 6}: "1/3",
		{3, 7}: "3/7",
		{1, 5}: "1/5",
		{4, 8}: "1/2",
	}
	for in, want := range cases {
		if got := reducedFraction(in[0], in[1]); got != want {
			t.Errorf("reducedFraction(%d,%d) = %q, want %q", in[0], in[1], got, want)
		}
	}
}

func TestSequenceTermsAreReducedAndStartAtOne(t *testing.T) {
	g := &seqConceptGen{}
	re := regexp.MustCompile(`(\d+)/(\d+)`)
	seen := 0
	for i := 0; i < 200 && seen < 20; i++ {
		p := g.Generate(generator.GeneratorContext{Difficulty: 0.5})
		if !strings.Contains(p.Question, `a_{n} = n/(n+`) {
			continue
		}
		seen++
		if !strings.Contains(p.Question, "starting at \\(n=1\\)") {
			t.Fatalf("sequence prompt omits starting index: %q", p.Question)
		}
		for _, m := range re.FindAllStringSubmatch(p.Answer, -1) {
			n, _ := strconv.Atoi(m[1])
			d, _ := strconv.Atoi(m[2])
			if reducedFraction(n, d) != m[0] {
				t.Fatalf("answer contains unreduced fraction %q in %q", m[0], p.Answer)
			}
		}
	}
	if seen == 0 {
		t.Fatal("sequence case 0 did not generate in 200 draws")
	}
}
