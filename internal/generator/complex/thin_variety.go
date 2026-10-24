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
	"fmt"
	"math/rand"

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
	shapes := []hypShape{
		{
			question: "For real θ, which real function equals cosh(iθ)?",
			fn:       "cosh", argument: "iθ", value: "cos θ",
			why: "cosh(iθ) = (e^(iθ) + e^(-iθ))/2, which is the identity for cos θ, so the imaginary argument turns cosh into cos.",
			d1:  "sin θ", d2: "i sin θ", d3: "cosh θ",
		},
		{
			question: "For real θ, which function equals sinh(iθ)?",
			fn:       "sinh", argument: "iθ", value: "i sin θ",
			why: "sinh(iθ) = (e^(iθ) - e^(-iθ))/2, and e^(iθ) - e^(-iθ) = 2i sin θ, so the imaginary argument leaves a factor of i.",
			d1:  "sin θ", d2: "cos θ", d3: "i cosh θ",
		},
		{
			question: "For real θ, which function equals tanh(iθ)?",
			fn:       "tanh", argument: "iθ", value: "i tan θ",
			why: "tanh(iθ) = sinh(iθ)/cosh(iθ) = (i sin θ)/(cos θ) = i tan θ, so the quotient inherits the factor of i from the sine.",
			d1:  "tan θ", d2: "-i tan θ", d3: "i tanh θ",
		},
		{
			question: "For real θ, what is cosh(iπ)?",
			fn:       "cosh", argument: "iπ", value: "-1",
			why: "cosh(iπ) = cos π by the imaginary-argument identity, and cos π = -1.",
			d1:  "1", d2: "i", d3: "-i",
		},
		{
			question: "For real θ, what is sinh(iπ)?",
			fn:       "sinh", argument: "iπ", value: "0",
			why: "sinh(iπ) = i sin π by the imaginary-argument identity, and sin π = 0.",
			d1:  "1", d2: "-1", d3: "iπ",
		},
	}
	for n := 1; n <= 5; n++ {
		shapes = append(shapes,
			hypShape{
				question: fmt.Sprintf("Which real quantity equals cosh(%di)?", n),
				fn:       "cosh", argument: fmt.Sprintf("%di", n), value: fmt.Sprintf("cos %d", n),
				why: fmt.Sprintf("cosh(%di) = (e^(%di) + e^(-%di))/2 = cos %d, because the two exponentials add to twice the cosine.", n, n, n, n),
				d1:  fmt.Sprintf("-cos %d", n),
				d2:  fmt.Sprintf("cosh %d", n),
				d3:  fmt.Sprintf("i sin %d", n),
			},
			hypShape{
				question: fmt.Sprintf("Which quantity equals sinh(%di)?", n),
				fn:       "sinh", argument: fmt.Sprintf("%di", n), value: fmt.Sprintf("i sin %d", n),
				why: fmt.Sprintf("sinh(%di) = (e^(%di) - e^(-%di))/2 = i sin %d, because the difference of the exponentials is twice i times the sine.", n, n, n, n),
				d1:  fmt.Sprintf("-i sin %d", n),
				d2:  fmt.Sprintf("sinh %d", n),
				d3:  fmt.Sprintf("cos %d", n),
			},
		)
	}
	s := shapes[rand.Intn(len(shapes))]
	p := complexChoice(s.question, s.value, s.why, s.d1, s.d2, s.d3)
	p.Facts = map[string]string{"fn": s.fn, "argument": s.argument, "value": s.value}
	return p
}

// ----- trigonometric functions in exponential form -----

type trigExponentialGen struct{}

func (g *trigExponentialGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	n := 1 + rand.Intn(4)
	arg, inverse, label := "iθ", "-iθ", "θ"
	if n > 1 {
		arg = fmt.Sprintf("%diθ", n)
		inverse = fmt.Sprintf("-%diθ", n)
		label = fmt.Sprintf("%dθ", n)
	}
	cosForm := fmt.Sprintf("(e^(%s) + e^(%s))/2", arg, inverse)
	sinForm := fmt.Sprintf("(e^(%s) - e^(%s))/(2i)", arg, inverse)
	cosSigned := fmt.Sprintf("(e^(%s) - e^(%s))/2", arg, inverse)
	cosOverI := fmt.Sprintf("(e^(%s) + e^(%s))/(2i)", arg, inverse)

	if rand.Intn(2) == 0 {
		p := complexChoice(
			fmt.Sprintf("Which identity expresses %s in the form cos %s = ...?", label, label),
			cosForm,
			fmt.Sprintf("Adding the Euler formulas cancels the sine terms: e^(%s) + e^(%s) = 2 cos %s, so cos %s = %s.", arg, inverse, label, label, cosForm),
			sinForm, cosSigned, cosOverI,
		)
		p.Facts = map[string]string{"asked": "cos", "cos": cosForm, "sin": sinForm, "value": cosForm}
		return p
	}
	p := complexChoice(
		fmt.Sprintf("Which identity expresses %s in the form sin %s = ...?", label, label),
		sinForm,
		fmt.Sprintf("Subtracting the Euler formulas leaves the sine terms with a factor of i: e^(%s) - e^(%s) = 2i sin %s, so sin %s = %s.", arg, inverse, label, label, sinForm),
		cosForm, cosSigned, cosOverI,
	)
	p.Facts = map[string]string{"asked": "sin", "cos": cosForm, "sin": sinForm, "value": sinForm}
	return p
}

// ----- Schwarz lemma -----

type schwarzLemmaGen struct{}

func (g *schwarzLemmaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(5) {
	case 0:
		return complexChoice(
			"If \\(f\\) maps the unit disc to itself and \\(f(0)=0\\), what is the largest possible value of \\(|f'(0)|\\)?",
			"1",
			"Schwarz lemma gives \\(|f'(0)|\\) at most 1, and equality happens exactly for a rotation \\(f(z)=cz\\) with \\(|c|=1\\).",
			"2", "1/2", "0",
		)
	case 1:
		return complexChoice(
			"Under Schwarz lemma, when does \\(|f'(0)|=1\\) hold?",
			"when f is a rotation",
			"Equality forces the function to be a rotation of the disc, which is the lemma's equality case.",
			"never",
			"when f is constant",
			"when f has a pole in the disc",
		)
	case 2:
		return complexChoice(
			"If \\(f\\) maps the unit disc to itself and \\(f(0)=0\\), what does Schwarz lemma say about \\(|f(z)|\\) for \\(|z|<1\\)?",
			"|f(z)| is at most |z|",
			"The lemma bounds the image of every point, not only the derivative at the centre: \\(|f(z)| \\le |z|\\), which is the derivative bound integrated along the radius.",
			"|f(z)| equals 1",
			"|f(z)| is at least |z|",
			"|f(z)| is unbounded",
		)
	case 3:
		return complexChoice(
			"If \\(f\\) maps the unit disc to itself with \\(f(0)=a\\) and \\(|a|<1\\), what is the largest possible \\(|f'(0)|\\)?",
			"1 - |a|^2",
			"Schwarz-Pick applies the lemma to the disc automorphism sending \\(a\\) to 0, which leaves the derivative bound \\(|f'(0)| \\le 1-|a|^2\\).",
			"1",
			"1 + |a|^2",
			"1 - |a|",
		)
	default:
		return complexChoice(
			"If \\(f\\) maps the unit disc to itself, \\(f(0)=0\\), and \\(f\\) is not a rotation, where is \\(|f(z)|<|z|\\) guaranteed?",
			"at every z with 0<|z|<1",
			"A non-rotation cannot attain the bound anywhere, because equality at a single interior point already forces the function to be a rotation.",
			"nowhere",
			"only near the boundary",
			"only at the centre",
		)
	}
}
