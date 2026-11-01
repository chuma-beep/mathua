package complex

import (
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

// complexChoice builds a multiple choice problem with an explanation. These
// concepts are graded multiple_choice, and MultipleChoice returns no
// explanation, so every call site supplies one.
func complexChoice(question, correct, explanation string, distractors ...string) generator.Problem {
	p := generator.MultipleChoice(question, correct, distractors...)
	p.Explanation = explanation
	return p
}

type analyticGen struct{}

func (g *analyticGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return analyticProblem()
}

type cauchyRiemannGen struct{}

func (g *cauchyRiemannGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return cauchyRiemannProblem()
}

type conformalGen struct{}

func (g *conformalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return conformalProblem()
}

// intPowLocal is a small integer power used to print derivative values.
func intPowLocal(base, exp int) int {
	out := 1
	for i := 0; i < exp; i++ {
		out *= base
	}
	return out
}

type eulerIdentityGen struct{}

func (g *eulerIdentityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return eulerIdentityProblem()
}

// unitCoef renders a coefficient for a term x, dropping a leading one.
func unitCoef(c int) string {
	if c == 1 {
		return ""
	}
	return strconv.Itoa(c)
}

type harmonicGen struct{}

func (g *harmonicGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return harmonicProblem()
}

type mobiusTransformGen struct{}

func (g *mobiusTransformGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return mobiusProblem()
}

type residueGen struct{}

func (g *residueGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return residueProblem()
}
