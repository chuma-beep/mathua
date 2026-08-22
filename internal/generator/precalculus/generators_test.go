package precalculus

import (
	"math/rand"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/grader"
)

func TestRegisterCoversAllConcepts(t *testing.T) {
	reg := generator.NewRegistry()
	Register(reg)
	for _, e := range all {
		p := e.gen.Generate(generator.GeneratorContext{Difficulty: 0.5})
		if p.Question == "" || p.Answer == "" || p.Explanation == "" {
			t.Errorf("%s produced empty field: q=%q a=%q", e.id, p.Question, p.Answer)
		}
	}
}

func TestFuzz(t *testing.T) {
	router := grader.NewRouter()
	for _, e := range all {
		gen := e.gen
		t.Run(e.id, func(t *testing.T) {
			for i := 0; i < 100; i++ {
				d := rand.Float64()
				p := gen.Generate(generator.GeneratorContext{Difficulty: d})
				if p.Question == "" || p.Answer == "" || p.Explanation == "" {
					t.Fatalf("empty field at difficulty=%.2f: q=%q a=%q e=%q", d, p.Question, p.Answer, p.Explanation)
				}
				correct := false
				for _, gt := range []grader.GradingType{
					grader.GradingNumeric,
					grader.GradingExpression,
					grader.GradingMultipleChoice,
					grader.GradingTuple,
					grader.GradingPolynomial,
				} {
					if res := router.Grade(gt, p.Answer, p.Answer); res.Correct {
						correct = true
						break
					}
				}
				if !correct {
					t.Fatalf("self-grade failed at difficulty=%.2f for question %q (answer=%q)", d, p.Question, p.Answer)
				}
				if gg, ok := gen.(interface {
					Grade(expected, userAnswer string) grader.Result
				}); ok {
					if res := gg.Grade(p.Answer, p.Answer); !res.Correct {
						t.Fatalf("custom Grade rejects own answer %q at difficulty=%.2f: %s", p.Answer, d, res.Feedback)
					}
				}
			}
		})
	}
}

func TestGradeRejectsWrongAnswers(t *testing.T) {
	cases := []struct {
		id       string
		expected string
		wrong    string
	}{
		{"precalc.vector.ops", "(3, -2)", "(3, 2)"},
		{"precalc.systems.nonlinear", "(1, 4)", "(2, 3)"},
		{"precalc.vector.basics", "5", "4"},
		{"precalc.func.piecewise", "7", "8"},
	}
	for _, c := range cases {
		for _, e := range all {
			if e.id != c.id {
				continue
			}
			gg, ok := e.gen.(interface {
				Grade(expected, userAnswer string) grader.Result
			})
			if !ok {
				continue
			}
			if res := gg.Grade(c.expected, c.expected); !res.Correct {
				t.Errorf("%s: expected answer should self-pass", c.id)
			}
			if res := gg.Grade(c.expected, c.wrong); res.Correct {
				t.Errorf("%s: wrong answer (%q) graded correct", c.id, c.wrong)
			}
		}
	}
}

func TestCanonPoly(t *testing.T) {
	cases := []struct {
		in   map[int]int
		want string
	}{
		{map[int]int{2: 1}, "x^2"},
		{map[int]int{1: -1}, "-x"},
		{map[int]int{0: 3}, "3"},
		{map[int]int{2: 2, 0: -3, 1: 1}, "2x^2 + x - 3"},
		{map[int]int{}, "0"},
	}
	for _, c := range cases {
		if got := canonPoly(c.in); got != c.want {
			t.Errorf("canonPoly(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}
