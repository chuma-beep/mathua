package probability

import (
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

// Report 488: the answer was printed to six places, so a learner who rounded
// 4/11 to 0.3636 was marked wrong. Four places is the precision the question
// asks for, and half a unit there covers the printed value exactly.
func TestBayesAnswerRoundsToFourPlaces(t *testing.T) {
	g := &bayesGen{}
	for i := 0; i < 500; i++ {
		p := g.Generate(generator.GeneratorContext{})
		aA, _ := strconv.Atoi(p.Facts["aA"])
		tA, _ := strconv.Atoi(p.Facts["tA"])
		aB, _ := strconv.Atoi(p.Facts["aB"])
		tB, _ := strconv.Atoi(p.Facts["tB"])
		pa := float64(aA) / float64(tA)
		pb := float64(aB) / float64(tB)
		want := decimal6(math.Round(pa/(pa+pb)*1e4) / 1e4)
		if p.Answer != want {
			t.Fatalf("%q answered %s, want %s", p.Question, p.Answer, want)
		}
		if dot := strings.IndexByte(p.Answer, '.'); dot >= 0 {
			if places := len(p.Answer) - dot - 1; places > 4 {
				t.Fatalf("%q answered to %d places: %s", p.Question, places, p.Answer)
			}
		}
	}
}
