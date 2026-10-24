package calculus

import (
	"fmt"
	"math/rand"
)

// calcDistinct returns three distractors for correct that differ from it and
// from each other.
//
// generator.MultipleChoice caps at four options and dedupes, so a candidate
// equal to the answer drops the question to three options and fails the
// recognition-question test.
func calcDistinct(correct string, candidates ...string) []string {
	out := make([]string, 0, 3)
	seen := map[string]bool{correct: true}
	for _, c := range candidates {
		if len(out) == 3 {
			break
		}
		if seen[c] {
			continue
		}
		seen[c] = true
		out = append(out, c)
	}
	for _, c := range []string{"0", "1", "-1", "none of these"} {
		if len(out) == 3 {
			break
		}
		if seen[c] {
			continue
		}
		seen[c] = true
		out = append(out, c)
	}
	return out
}

// computedMultiClaims returns parameterised claims for the concepts whose pools
// held three fixed sentences.
func computedMultiClaims(id string) []calcMCClaim {
	switch id {
	case "calc.multi.iterated_limits":
		return iteratedLimitClaims()
	case "calc.vec.identities":
		return vectorIdentityClaims()
	}
	return nil
}

func iteratedLimitClaims() []calcMCClaim {
	a := 1 + rand.Intn(4)
	b := 1 + rand.Intn(4)
	p := 2 + rand.Intn(2)
	q := 2 + rand.Intn(2)
	sum := fmt.Sprintf("%d", powIntCalc(a, p)+powIntCalc(b, q))
	linear := fmt.Sprintf("%d", a+b)
	product := fmt.Sprintf("%d", a*b)
	return []calcMCClaim{
		{
			q:           fmt.Sprintf("What is the limit of f(x,y) = x^%d + y^%d as (x,y) tends to (%d, %d)?", p, q, a, b),
			correct:     sum,
			distractors: calcDistinct(sum, fmt.Sprintf("%d", a+b), fmt.Sprintf("%d", a*p+b*q), fmt.Sprintf("%d", a*b)),
			reason:      fmt.Sprintf("A sum of powers is continuous, so substituting gives %d^%d + %d^%d = %s.", a, p, b, q, sum),
		},
		{
			q:           fmt.Sprintf("What is the limit of (x^2 - y^2)/(x - y) as (x,y) tends to (%d, %d)?", a, b),
			correct:     linear,
			distractors: calcDistinct(linear, fmt.Sprintf("%d", a-b), "0", fmt.Sprintf("%d", a*b)),
			reason:      fmt.Sprintf("The numerator factors as (x - y)(x + y), so along the punctured domain the function is x + y, which tends to %s.", linear),
		},
		{
			q:           fmt.Sprintf("What is the limit of f(x,y) = x y as (x,y) tends to (%d, %d)?", a, b),
			correct:     product,
			distractors: calcDistinct(product, fmt.Sprintf("%d", a+b), fmt.Sprintf("%d", a-1), fmt.Sprintf("%d", 2*a*b)),
			reason:      fmt.Sprintf("The product is continuous, so the limit is %d * %d = %s.", a, b, product),
		},
	}
}

func vectorIdentityClaims() []calcMCClaim {
	p := 1 + rand.Intn(5)
	q := 1 + rand.Intn(5)
	r := 1 + rand.Intn(5)
	s := 1 + rand.Intn(5)
	k := 2 + rand.Intn(4)

	normSquared := fmt.Sprintf("%d", p*p+q*q)
	dotValue := p*r + q*s
	dot := fmt.Sprintf("%d", dotValue)
	scaled := fmt.Sprintf("%d", k*dotValue)
	return []calcMCClaim{
		{
			q:           fmt.Sprintf("For the vector v = (%d, %d), what is |v|^2?", p, q),
			correct:     normSquared,
			distractors: calcDistinct(normSquared, fmt.Sprintf("%d", p+q), fmt.Sprintf("%d", (p+q)*(p+q)), fmt.Sprintf("%d", p*q)),
			reason:      fmt.Sprintf("The squared norm is the sum of squared components: %d + %d = %s.", p*p, q*q, normSquared),
		},
		{
			q:           fmt.Sprintf("For u = (%d, %d) and v = (%d, %d), what is u · v?", p, q, r, s),
			correct:     dot,
			distractors: calcDistinct(dot, fmt.Sprintf("%d", p*s+q*r), fmt.Sprintf("%d", p*r-q*s), fmt.Sprintf("%d", (p+r)*(q+s))),
			reason:      fmt.Sprintf("The dot product multiplies matching components and adds: %d*%d + %d*%d = %s.", p, r, q, s, dot),
		},
		{
			q:           fmt.Sprintf("If u · v = %d, what is (%du) · v?", dotValue, k),
			correct:     scaled,
			distractors: calcDistinct(scaled, fmt.Sprintf("%d", dotValue), fmt.Sprintf("%d", dotValue+k), fmt.Sprintf("%d", dotValue*dotValue)),
			reason:      fmt.Sprintf("The dot product is linear in each argument, so the factor %d multiplies the result: %s.", k, scaled),
		},
	}
}

// powIntCalc is base raised to exp for small non-negative exponents.
func powIntCalc(base, exp int) int {
	out := 1
	for i := 0; i < exp; i++ {
		out *= base
	}
	return out
}
