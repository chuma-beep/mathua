package statistics

import (
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

// Report 450: the z-score question used x=70 and mean=50, but the explanation
// wrote 30/10. The answer 2 was right while the working shown was wrong.
func TestNormalZScoreExplanationMatchesQuestion(t *testing.T) {
	g := &normalDistGen{}
	for i := 0; i < 500; i++ {
		p := g.Generate(generator.GeneratorContext{})
		if !strings.Contains(p.Question, "z-score") {
			continue
		}
		if !strings.Contains(p.Explanation, "20/10") {
			t.Fatalf("z-score explanation does not match x=70 and mean=50: %q", p.Explanation)
		}
		return
	}
	t.Fatal("z-score entry was never drawn")
}

func TestMeanAnswerIsExact(t *testing.T) {
	gen := &meanGen{}
	for seed := int64(1); seed <= 200; seed++ {
		p := gen.Generate(generator.GeneratorContext{Difficulty: 0.5, Seed: seed})
		sum, _ := strconv.Atoi(p.Facts["sum"])
		n, _ := strconv.Atoi(p.Facts["n"])
		if p.Answer != p.Facts["answer"] {
			t.Fatalf("seed %d: answer %q does not match fact %q", seed, p.Answer, p.Facts["answer"])
		}
		if strings.Contains(p.Answer, ".") {
			t.Fatalf("seed %d: answer %q is a rounded decimal", seed, p.Answer)
		}
		if sum%n == 0 {
			if p.Answer != strconv.Itoa(sum/n) {
				t.Fatalf("seed %d: exact mean should be %d, got %q", seed, sum/n, p.Answer)
			}
			continue
		}
		num, den, ok := strings.Cut(p.Answer, "/")
		if !ok {
			t.Fatalf("seed %d: non-integer mean %q is not a fraction", seed, p.Answer)
		}
		a, _ := strconv.Atoi(num)
		b, _ := strconv.Atoi(den)
		if a*n != sum*b {
			t.Fatalf("seed %d: %q does not equal %d/%d", seed, p.Answer, sum, n)
		}
	}
}

// Report 468: "at least at the ordinal level" was read as "only ordinal", so a
// learner who counted the single ordinal variable answered 1. The stem now
// names all three levels, and the explanation lists the level of every
// variable so the count can be checked.
func TestMeasurementLevelsStemNamesEveryQualifyingLevel(t *testing.T) {
	gen := &measurementLevelsGen{}
	for seed := int64(1); seed <= 100; seed++ {
		p := gen.Generate(generator.GeneratorContext{Difficulty: 0.5, Seed: seed, KPOrdinal: 2})
		if !strings.Contains(p.Question, "at the ordinal level or higher (ordinal, interval or ratio)") {
			t.Fatalf("seed %d: stem does not name the qualifying levels: %q", seed, p.Question)
		}
		count := 0
		for _, field := range strings.Split(p.Explanation, ";") {
			i := strings.LastIndex(field, ": level ")
			if i < 0 {
				continue
			}
			raw := strings.TrimSuffix(strings.TrimSpace(field[i+len(": level "):]), ".")
			level, err := strconv.Atoi(raw)
			if err != nil {
				t.Fatalf("seed %d: unreadable level %q in %q", seed, raw, field)
			}
			if level >= 2 {
				count++
			}
		}
		if got := strconv.Itoa(count); p.Answer != got {
			t.Fatalf("seed %d: answer %q does not match %s qualifying levels: %q", seed, p.Answer, got, p.Explanation)
		}
	}
}
