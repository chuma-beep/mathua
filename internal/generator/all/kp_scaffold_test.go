package all

import (
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func TestGenerators_DeterministicKPOrdinalVariants(t *testing.T) {
	reg := generator.NewRegistry()
	Register(reg)

	t.Run("arithmetic div zero count", func(t *testing.T) {
		for _, kp := range []int{0, 1, 2, 3} {
			p, err := reg.GenerateContext("arith.div.zero", generator.GeneratorContext{KPOrdinal: kp, Seed: 42})
			if err != nil {
				t.Fatalf("KP %d: unexpected error: %v", kp, err)
			}
			if p.Answer < "0" || p.Answer > "8" || len(p.Answer) != 1 {
				t.Errorf("KP %d: expected a single-digit count, got %q", kp, p.Answer)
			}
			if p.Facts["count"] != p.Answer {
				t.Errorf("KP %d: count fact %q does not match answer %q", kp, p.Facts["count"], p.Answer)
			}
		}
	})

	t.Run("algebra multi step equation shapes", func(t *testing.T) {
		p1, err := reg.GenerateContext("alg.eq.multi_step", generator.GeneratorContext{KPOrdinal: 1, Seed: 101})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p1.Facts["which"] != "axb" {
			t.Errorf("expected KP 1 which=axb, got %q", p1.Facts["which"])
		}

		p2, err := reg.GenerateContext("alg.eq.multi_step", generator.GeneratorContext{KPOrdinal: 2, Seed: 102})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p2.Facts["which"] != "axb" {
			t.Errorf("expected KP 2 which=axb, got %q", p2.Facts["which"])
		}

		p3, err := reg.GenerateContext("alg.eq.multi_step", generator.GeneratorContext{KPOrdinal: 3, Seed: 103})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p3.Facts["which"] != "distributed" {
			t.Errorf("expected KP 3 which=distributed, got %q", p3.Facts["which"])
		}
		if !strings.Contains(p3.Question, "(") {
			t.Errorf("expected distributed equation to contain parentheses, got %q", p3.Question)
		}
	})

	t.Run("algebra slope intercept variants", func(t *testing.T) {
		p1, err := reg.GenerateContext("alg.linear.slope_intercept", generator.GeneratorContext{KPOrdinal: 1, Seed: 201})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p1.Facts["which"] != "integer" {
			t.Errorf("expected KP 1 which=integer, got %q", p1.Facts["which"])
		}

		p3, err := reg.GenerateContext("alg.linear.slope_intercept", generator.GeneratorContext{KPOrdinal: 3, Seed: 203})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p3.Facts["which"] != "fraction" {
			t.Errorf("expected KP 3 which=fraction, got %q", p3.Facts["which"])
		}
	})

	t.Run("calculus area between curves modes", func(t *testing.T) {
		p1, err := reg.GenerateContext("calc.integral.area_between", generator.GeneratorContext{KPOrdinal: 1, Seed: 301})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(p1.Question, "and the x-axis") {
			t.Errorf("expected KP 1 to ask for line and x-axis area, got %q", p1.Question)
		}

		p2, err := reg.GenerateContext("calc.integral.area_between", generator.GeneratorContext{KPOrdinal: 2, Seed: 302})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(p2.Question, "- x^2") {
			t.Errorf("expected KP 2 to ask for parabola area, got %q", p2.Question)
		}

		p3, err := reg.GenerateContext("calc.integral.area_between", generator.GeneratorContext{KPOrdinal: 3, Seed: 303})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(p3.Question, "bounded by") || !strings.Contains(p3.Question, "y = x^2") {
			t.Errorf("expected KP 3 to ask for line and parabola intersection area, got %q", p3.Question)
		}
	})

	t.Run("calculus mean value theorem modes", func(t *testing.T) {
		p1, err := reg.GenerateContext("calc.deriv.mvt", generator.GeneratorContext{KPOrdinal: 1, Seed: 401})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(p1.Question, "f(x) = x^2") {
			t.Errorf("expected KP 1 to use base case f(x) = x^2, got %q", p1.Question)
		}

		p2, err := reg.GenerateContext("calc.deriv.mvt", generator.GeneratorContext{KPOrdinal: 2, Seed: 402})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(p2.Question, "x^2 +") {
			t.Errorf("expected KP 2 to use general quadratic, got %q", p2.Question)
		}

		p3, err := reg.GenerateContext("calc.deriv.mvt", generator.GeneratorContext{KPOrdinal: 3, Seed: 403})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(p3.Question, "average rate of change") {
			t.Errorf("expected KP 3 to ask for average rate of change, got %q", p3.Question)
		}
	})
}

func TestGenerators_BackwardsCompatibilityKPOrdinalZero(t *testing.T) {
	reg := generator.NewRegistry()
	Register(reg)

	conceptsToCheck := []string{
		"arith.div.zero",
		"arith.div.remainder",
		"arith.add.word",
		"alg.eq.multi_step",
		"alg.linear.slope_intercept",
		"alg.eq.vars_both_sides",
		"calc.integral.area_between",
		"calc.deriv.mvt",
		"calc.deriv.rolle",
	}

	for _, cid := range conceptsToCheck {
		p, err := reg.GenerateContext(cid, generator.GeneratorContext{Difficulty: 0.5, Seed: 777})
		if err != nil {
			t.Fatalf("concept %q failed at KPOrdinal 0: %v", cid, err)
		}
		if p.Question == "" || p.Answer == "" || p.Explanation == "" {
			t.Errorf("concept %q returned empty fields at KPOrdinal 0: %+v", cid, p)
		}
	}
}
