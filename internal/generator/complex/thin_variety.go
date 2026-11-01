package complex

// Authored variety for the three concepts the canonical ratchet measured at two
// distinct situations each. These could not be widened by parameter: each was a
// two-way switch over a fixed pair of questions, so the supply of problems was
// the size of the switch.
//
// The two schema-backed concepts stay inside the narrative their authored
// walkthrough tells. data/lessons/solutions/complex.adv.hyperbolic.json says the
// function is evaluated at an argument and the exponentials recombine into a
// trigonometric function; every shape here is exactly that. The
// trig_exponential schema names the cosine and sine identities as the facts it
// cites, so those two identities are asked at every small multiple of the angle
// rather than replaced by a different theorem.

import (
	"github.com/chuma-beep/mathua/internal/generator"
)

// hypShape is one question about a hyperbolic function at an imaginary
// argument, carrying the facts its authored walkthrough cites.
type hypShape struct {
	question   string
	fn         string
	argument   string
	value      string
	why        string
	d1, d2, d3 string
}

// ----- hyperbolic functions on imaginary arguments -----

type hyperbolicGen struct{}

func (g *hyperbolicGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return hyperbolicProblem()
}

// ----- trigonometric functions in exponential form -----

type trigExponentialGen struct{}

func (g *trigExponentialGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return trigExponentialProblem()
}

// ----- Schwarz lemma -----

type schwarzLemmaGen struct{}

func (g *schwarzLemmaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return schwarzLemmaProblem()
}
