package algebra

import (
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func TestTurningPointsLeastDegreeMatchesAnswer(t *testing.T) {
	pattern := regexp.MustCompile(`exactly ([0-9]+) turning points`)
	for i := 0; i < 300; i++ {
		p := (&turningPointsGen{}).Generate(generator.GeneratorContext{KPOrdinal: 3})
		m := pattern.FindStringSubmatch(p.Question)
		if m == nil {
			t.Fatalf("question does not state the turning-point count: %q", p.Question)
		}
		turns, err := strconv.Atoi(m[1])
		if err != nil {
			t.Fatalf("turning-point count %q: %v", m[1], err)
		}
		want := strconv.Itoa(turns + 1)
		if p.Answer != want {
			t.Fatalf("%d turning points should need degree %s, got answer %s: %q", turns, want, p.Answer, p.Question)
		}
		if !strings.Contains(p.Explanation, "at least "+want) {
			t.Fatalf("explanation does not state the least degree %s: %q", want, p.Explanation)
		}
	}
}

func TestRemainderDivisorHasNoDoubleSign(t *testing.T) {
	for _, kp := range []int{1, 3} {
		for i := 0; i < 300; i++ {
			p := (&polyRemainderGen{}).Generate(generator.GeneratorContext{KPOrdinal: kp})
			if strings.Contains(p.Question, "- -") {
				t.Fatalf("divisor rendered a double sign: %q", p.Question)
			}
		}
	}
}

func TestCramerEquationOmitsUnitCoefficient(t *testing.T) {
	for i := 0; i < 300; i++ {
		p := (&cramerGen{}).Generate(generator.GeneratorContext{KPOrdinal: 1})
		if strings.Contains(p.Question, "1x") || strings.Contains(p.Question, "1y") {
			t.Fatalf("equation rendered a unit coefficient: %q", p.Question)
		}
	}
}
