package numtheory

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

// ntVerdict is one statement about a number theory topic and whether it is
// true. The count questions below ask how many of eight statements hold, so a
// learner has to judge each claim instead of recognising a single option.
type ntVerdict struct {
	label string
	yes   bool
	note  string
}

const ntDraw = 8

// countNTVerdicts draws the answer uniformly from 0 to 8, then chooses that
// many true statements and the rest false ones. Drawing the answer first keeps
// every value in range even when the pool is large, which a plain sample of
// eight candidates would not do.
func countNTVerdicts(lead string, pool []ntVerdict) generator.Problem {
	yes := make([]ntVerdict, 0, len(pool))
	no := make([]ntVerdict, 0, len(pool))
	for _, v := range pool {
		if v.yes {
			yes = append(yes, v)
		} else {
			no = append(no, v)
		}
	}

	for attempt := 0; attempt < 40; attempt++ {
		count := rand.Intn(ntDraw + 1)
		if count > len(yes) {
			count = len(yes)
		}
		if ntDraw-count > len(no) {
			count = ntDraw - len(no)
		}

		y := append([]ntVerdict(nil), yes...)
		n := append([]ntVerdict(nil), no...)
		rand.Shuffle(len(y), func(i, j int) { y[i], y[j] = y[j], y[i] })
		rand.Shuffle(len(n), func(i, j int) { n[i], n[j] = n[j], n[i] })

		chosen := append([]ntVerdict(nil), y[:count]...)
		chosen = append(chosen, n[:ntDraw-count]...)
		rand.Shuffle(len(chosen), func(i, j int) { chosen[i], chosen[j] = chosen[j], chosen[i] })

		parts := make([]string, len(chosen))
		notes := make([]string, len(chosen))
		for i, v := range chosen {
			parts[i] = fmt.Sprintf("(%c) %s", rune('a'+i), v.label)
			verdict := "false"
			if v.yes {
				verdict = "true"
			}
			notes[i] = fmt.Sprintf("%s: %s, %s", v.label, verdict, v.note)
		}
		question := fmt.Sprintf("%s %s (enter a number)", lead, strings.Join(parts, " "))
		if strings.Contains(question, strconv.Itoa(count)) {
			continue
		}
		return generator.Problem{
			Question:    question,
			Answer:      strconv.Itoa(count),
			Explanation: fmt.Sprintf("Exactly %d of the 8 statements are true. %s", count, strings.Join(notes, " ")),
			Facts:       map[string]string{"count": strconv.Itoa(count)},
		}
	}

	// The retry loop only rejects questions that leak the answer. Fall back to a
	// fixed count so the generator always returns a problem.
	count := 4
	chosen := append([]ntVerdict(nil), yes[:count]...)
	chosen = append(chosen, no[:ntDraw-count]...)
	parts := make([]string, len(chosen))
	for i, v := range chosen {
		parts[i] = fmt.Sprintf("(%c) %s", rune('a'+i), v.label)
	}
	return generator.Problem{
		Question:    fmt.Sprintf("%s %s (enter a number)", lead, strings.Join(parts, " ")),
		Answer:      strconv.Itoa(count),
		Explanation: "Exactly 4 of the 8 statements are true.",
		Facts:       map[string]string{"count": strconv.Itoa(count)},
	}
}

func quadraticReciprocityProblem() generator.Problem {
	return countNTVerdicts("How many of these statements about quadratic reciprocity are true?", []ntVerdict{
		{"For odd primes p and q, (p/q)(q/p) = (-1)^((p-1)(q-1)/4)", true, "the sign depends on the two primes modulo 4"},
		{"(2/p) = (-1)^((p^2-1)/8)", true, "the supplementary law for 2"},
		{"The Legendre symbol is multiplicative in its numerator", true, "(ab/p) = (a/p)(b/p)"},
		{"(1/p) = 1 for every odd prime p", true, "1 is always a square"},
		{"(-1/p) = (-1)^((p-1)/2)", true, "minus one is a square exactly for p = 1 mod 4"},
		{"The law of quadratic reciprocity applies to odd primes", true, "the hypothesis rules out 2"},
		{"The Legendre symbol takes only the values 1 and -1", false, "it is 0 when p divides the numerator"},
		{"Quadratic reciprocity determines (a/p) without any sign", false, "the sign is the whole content of the law"},
		{"(a/p) = 0 for every a", false, "it is zero only when p divides a"},
		{"The law includes p = 2 as an ordinary odd prime", false, "2 needs the separate supplementary law"},
		{"The Legendre symbol can be any integer", false, "its values are -1, 0 and 1"},
		{"(-1/p) = 1 for every odd prime p", false, "it is -1 for p = 3 mod 4"},
		{"The symbol is multiplicative in the denominator only", false, "it is multiplicative in the numerator"},
		{"(p/q) equals p/q as a rational number", false, "the symbol is a square-class indicator"},
	})
}

func primitiveRootProblem() generator.Problem {
	return countNTVerdicts("How many of these statements about primitive roots are true?", []ntVerdict{
		{"Every prime p has a primitive root modulo p", true, "the group of units is cyclic"},
		{"There are phi(p-1) primitive roots modulo a prime p", true, "each generator of the cyclic unit group"},
		{"A primitive root modulo n has order phi(n)", true, "it generates the whole unit group"},
		{"Powers of a primitive root give every unit modulo n", true, "that is what generating means"},
		{"Primitive roots modulo n exist exactly for n = 2, 4, p^k and 2p^k", true, "the classification of cyclic unit groups"},
		{"2 is a primitive root modulo 3", true, "2 has order 2 = phi(3)"},
		{"3 is a primitive root modulo 7", true, "its powers run through all six units"},
		{"The order of a primitive root divides phi(n)", true, "it equals phi(n), so it divides it"},
		{"Every modulus has a primitive root", false, "8 has none"},
		{"A modulus has exactly one primitive root", false, "there are phi(phi(n)) of them when they exist"},
		{"The order of a primitive root modulo p is p", false, "the order is p-1"},
		{"2 is a primitive root modulo 7", false, "2^3 = 1, so its order is 3"},
		{"Every nonzero residue is a primitive root", false, "only generators of the unit group are"},
		{"8 has a primitive root", false, "the unit group of 8 is not cyclic"},
		{"A primitive root has order phi(n) + 1", false, "the order is phi(n)"},
		{"There are phi(n) primitive roots modulo n", false, "the count is phi(phi(n))"},
	})
}

func primeNumberTheoremProblem() generator.Problem {
	return countNTVerdicts("How many of these statements about the prime number theorem are true?", []ntVerdict{
		{"pi(x) is asymptotic to x/log x", true, "the defining statement of the theorem"},
		{"The n-th prime is asymptotic to n log n", true, "an equivalent form"},
		{"The density of primes near x is about 1/log x", true, "the local reading of the theorem"},
		{"Chebyshev function psi(x) is asymptotic to x", true, "the standard equivalent form"},
		{"Chebyshev proved explicit upper and lower bounds of the right order", true, "his bounds sandwich pi(x) by constants times x/log x"},
		{"The error term is smaller under the Riemann hypothesis", true, "RH gives a square-root error"},
		{"pi(10^6) is 78498", true, "the exact count below one million"},
		{"The prime number theorem is equivalent to psi(x) ~ x", true, "a standard equivalence"},
		{"pi(x) is asymptotic to x", false, "the count is x/log x, not x"},
		{"The density of primes is constant", false, "it falls like 1/log x"},
		{"The n-th prime is asymptotic to n^2", false, "it is n log n"},
		{"The Riemann zeta function has no zeros anywhere", false, "it has trivial and nontrivial zeros"},
		{"The theorem says there are finitely many primes", false, "there are infinitely many"},
		{"There is an exact elementary formula for pi(x)", false, "no such simple closed formula is known"},
		{"pi(x) is asymptotic to log x", false, "that is far too small"},
		{"The Riemann hypothesis has been proved", false, "it remains open"},
	})
}

func sieveProblem() generator.Problem {
	return countNTVerdicts("How many of these statements about sieves are true?", []ntVerdict{
		{"The sieve of Eratosthenes marks multiples of each prime", true, "the marks are the composites"},
		{"To sieve up to n it is enough to use primes at most sqrt(n)", true, "a composite below n has a small factor"},
		{"The sieve of Eratosthenes uses about n log log n operations", true, "the classical complexity estimate"},
		{"A segmented sieve handles ranges too large for memory", true, "it processes blocks"},
		{"Wheel sieves skip numbers divisible by small primes", true, "they reduce the candidate set"},
		{"A sieve finds all primes up to its limit", true, "the unmarked numbers are prime"},
		{"A basic sieve uses O(n) memory", true, "one bit per candidate"},
		{"The sieve can be restricted to an interval", true, "segmented variants do exactly that"},
		{"A sieve tests every number by trial division", false, "it marks multiples instead"},
		{"A sieve can only find even primes", false, "it finds all primes"},
		{"A sieve needs a primality test for every candidate", false, "marking replaces individual tests"},
		{"Every sieve runs in exactly O(n^2) time", false, "the classical sieve is faster"},
		{"A sieve finds composites but not primes", false, "the unmarked entries are the primes"},
		{"A sieve cannot be segmented", false, "segmented sieves are standard"},
		{"A sieve requires factoring every number", false, "it avoids factoring"},
		{"A sieve marks the primes and leaves the composites", false, "it is the other way round"},
	})
}

func zetaProblem() generator.Problem {
	return countNTVerdicts("How many of these statements about the Riemann zeta function are true?", []ntVerdict{
		{"For Re(s) > 1, zeta(s) is the sum of n^(-s)", true, "the Dirichlet series converges there"},
		{"The Euler product expresses zeta as a product over primes", true, "the arithmetic form of the series"},
		{"Zeta has an analytic continuation with a pole at s = 1", true, "the only pole"},
		{"The trivial zeros are at -2, -4, -6 and so on", true, "from the gamma factor"},
		{"The critical strip is 0 < Re(s) < 1", true, "where nontrivial zeros live"},
		{"The Riemann hypothesis places nontrivial zeros on Re(s) = 1/2", true, "the conjecture"},
		{"zeta(2) = pi^2/6", true, "the Basel problem"},
		{"Zeta satisfies a functional equation relating s and 1-s", true, "the completed zeta function is symmetric"},
		{"The Dirichlet series converges for every s with Re(s) > 0", false, "it needs Re(s) > 1"},
		{"Zeta has a pole at s = 0", false, "the pole is at 1"},
		{"Zeta has no zeros at all", false, "it has trivial and nontrivial zeros"},
		{"Every zero of zeta is real", false, "the nontrivial zeros are complex"},
		{"zeta(1) = 0", false, "the pole at 1 is not a zero"},
		{"The Euler product holds for every complex s", false, "it needs Re(s) > 1"},
		{"zeta(3) is rational", false, "it is irrational"},
		{"The Riemann hypothesis has been proved", false, "it is still open"},
	})
}

func classNumberProblem() generator.Problem {
	return countNTVerdicts("How many of these statements about class numbers are true?", []ntVerdict{
		{"The class number measures the size of the ideal class group", true, "it is the group order"},
		{"Class number 1 means unique factorisation", true, "the ring is a UFD"},
		{"Only finitely many imaginary quadratic fields have class number 1", true, "the Heegner list"},
		{"The class number formula relates h to L-functions and regulators", true, "the analytic class number formula"},
		{"A class number above 1 means factorisation need not be unique", true, "nonprincipal ideals appear"},
		{"The ideal class group is abelian", true, "ideals commute"},
		{"The Heegner numbers give the imaginary quadratic class number 1 cases", true, "nine values"},
		{"A regulator appears in the class number formula for real quadratic fields", true, "the unit group contributes"},
		{"Every ring of integers has class number 1", false, "many do not"},
		{"The ideal class group can be nonabelian", false, "it is always abelian"},
		{"Infinitely many imaginary quadratic fields have class number 1", false, "the list is finite"},
		{"The class number equals the number of units", false, "it counts ideal classes"},
		{"A class number of 0 is possible", false, "the group has at least one class"},
		{"Every ideal class group is cyclic", false, "it can be a product of cyclic groups"},
		{"The class number is unchanged by every field extension", false, "it can grow"},
		{"Every class number is prime", false, "composite class numbers occur"},
	})
}

func chebyshevProblem() generator.Problem {
	return countNTVerdicts("How many of these statements about Chebyshev bounds are true?", []ntVerdict{
		{"The Chebyshev function psi(x) is asymptotic to x", true, "equivalent to the prime number theorem"},
		{"The Chebyshev function theta(x) is asymptotic to x", true, "the sum over primes at most x"},
		{"Chebyshev bounded pi(x) between constants times x/log x", true, "his elementary sandwich"},
		{"Bertrand's postulate follows from Chebyshev bounds", true, "a prime in (n, 2n)"},
		{"psi(x) is the sum of the von Mangoldt function up to x", true, "the definition"},
		{"theta(x) is the sum of log p over primes p <= x", true, "the definition"},
		{"Chebyshev bounds are equivalent to the prime number theorem", false, "they give the order, not the limit"},
		{"Chebyshev proved the prime number theorem", false, "he proved bounds only"},
		{"theta(x) is asymptotic to x^2", false, "it is asymptotic to x"},
		{"Bertrand's postulate is false", false, "it is a theorem"},
		{"psi(x) counts primes without weights", false, "it weights by log p"},
		{"theta(x) equals pi(x)", false, "theta weights each prime by log p"},
		{"Chebyshev bounds are exact equalities", false, "they are inequalities"},
		{"psi(x) equals pi(x) log x for every x", false, "only asymptotically"},
		{"There are no elementary bounds for pi(x)", false, "Chebyshev gave elementary ones"},
		{"Bertrand's postulate needs the prime number theorem", false, "it has elementary proofs"},
	})
}

func continuedFractionsProblem() generator.Problem {
	return countNTVerdicts("How many of these statements about continued fractions are true?", []ntVerdict{
		{"A rational number has a finite simple continued fraction", true, "the Euclidean algorithm terminates"},
		{"Convergents are best rational approximations", true, "no fraction with smaller denominator is closer"},
		{"sqrt(2) has the continued fraction [1; 2, 2, 2, ...]", true, "the classical periodic expansion"},
		{"A quadratic irrational has a eventually periodic continued fraction", true, "Lagrange's theorem"},
		{"Convergents alternate above and below the value", true, "the error changes sign"},
		{"A convergent p/q satisfies |alpha - p/q| < 1/q^2", true, "the standard accuracy bound"},
		{"Lagrange proved the periodic continued fraction theorem", true, "for quadratic irrationals"},
		{"Semiconvergents lie between adjacent convergents", true, "they fill the gaps"},
		{"Every continued fraction terminates", false, "irrationals give infinite expansions"},
		{"sqrt(2) has a terminating continued fraction", false, "its expansion is infinite"},
		{"Convergents are worse approximations than arbitrary fractions", false, "they are best approximations"},
		{"A continued fraction is periodic exactly for rational numbers", false, "quadratic irrationals are periodic"},
		{"[1; 2] equals 2/3", false, "it equals 3/2"},
		{"Every continued fraction converges to an integer", false, "it converges to its real value"},
		{"Convergent denominators decrease", false, "they increase"},
		{"Every quadratic irrational has period 1", false, "periods vary"},
	})
}

func dirichletProblem() generator.Problem {
	return countNTVerdicts("How many of these statements about Dirichlet's theorem are true?", []ntVerdict{
		{"There are infinitely many primes in a mod n when gcd(a,n) = 1", true, "Dirichlet's theorem"},
		{"Dirichlet used L-functions to prove the theorem", true, "the analytic method"},
		{"The primes are equidistributed with density 1/phi(n)", true, "among the eligible classes"},
		{"There are infinitely many primes of the form 4k+1", true, "the class 1 mod 4"},
		{"There are infinitely many primes of the form 4k+3", true, "the class 3 mod 4"},
		{"L(1, chi) is nonzero for a nonprincipal character", true, "the key nonvanishing"},
		{"Character orthogonality isolates a residue class", true, "the standard reduction"},
		{"The theorem applies to every residue class coprime to the modulus", true, "that is the hypothesis"},
		{"There are infinitely many primes in every residue class", false, "only the coprime classes"},
		{"The density of primes in a class is 1/n", false, "it is 1/phi(n)"},
		{"Only the class 1 mod n has infinitely many primes", false, "every coprime class does"},
		{"There are finitely many primes of the form 4k+3", false, "there are infinitely many"},
		{"L(1, chi) = 0 for every nonprincipal character", false, "it is nonzero"},
		{"Dirichlet's theorem needs n to be prime", false, "any modulus works"},
		{"Euclid's proof alone covers every residue class", false, "it does not"},
		{"There are finitely many primes of the form 4k+1", false, "there are infinitely many"},
	})
}

func ellipticCurveProblem() generator.Problem {
	return countNTVerdicts("How many of these statements about elliptic curves are true?", []ntVerdict{
		{"y^2 = x^3 + ax + b is nonsingular when its discriminant is nonzero", true, "the defining condition"},
		{"The points of an elliptic curve form an abelian group", true, "the chord and tangent law"},
		{"The Mordell-Weil theorem says the group is finitely generated", true, "over a number field"},
		{"Hasse's bound is |#E(F_p) - (p+1)| <= 2 sqrt(p)", true, "the trace of Frobenius is bounded"},
		{"An elliptic curve over C is a complex torus", true, "a quotient by a lattice"},
		{"The rank measures the free part of the Mordell-Weil group", true, "the number of independent generators"},
		{"Torsion points over Q are classified by Mazur's theorem", true, "fifteen possible groups"},
		{"The group law has a point at infinity as identity", true, "the identity of the group"},
		{"A singular curve still has the elliptic curve group law", false, "the group law needs nonsingularity"},
		{"The elliptic curve group is nonabelian", false, "it is abelian"},
		{"The Mordell-Weil group is always infinitely generated", false, "it is finitely generated"},
		{"#E(F_p) always equals p", false, "it is near p+1"},
		{"A zero discriminant gives a nonsingular curve", false, "zero discriminant means singular"},
		{"Every elliptic curve over Q has rank 1", false, "ranks vary"},
		{"An elliptic curve over Q has finitely many rational points", false, "the group can be infinite"},
		{"The group of points is always cyclic", false, "it can be a product"},
	})
}

func millerRabinProblem() generator.Problem {
	return countNTVerdicts("How many of these statements about Miller-Rabin are true?", []ntVerdict{
		{"Miller-Rabin is a probabilistic primality test", true, "random bases give a confidence level"},
		{"A number passing a base is a strong probable prime to that base", true, "the test condition"},
		{"Composite numbers that pass are called strong pseudoprimes", true, "to the chosen base"},
		{"Deterministic bases suffice for all 64-bit integers", true, "a known small set covers the range"},
		{"Carmichael numbers fool Fermat but not Miller-Rabin", true, "the strong test catches them"},
		{"The error probability falls as 4^(-k) after k rounds", true, "one quarter per random base"},
		{"The test uses modular exponentiation", true, "the core computation"},
		{"Most composites are rejected by the first base", true, "witnesses are common"},
		{"Miller-Rabin is deterministic for every input", false, "random bases make it probabilistic"},
		{"Miller-Rabin is the same as the Fermat test", false, "it is strictly stronger"},
		{"The error probability is 1/2 per round", false, "it is at most 1/4"},
		{"The test needs the factorization of n-1", false, "it only needs n-1 written as 2^s d"},
		{"The test applies only to even numbers", false, "it tests odd candidates"},
		{"A pass proves primality for every size", false, "only bounded ranges are deterministic"},
		{"The test is trial division", false, "it uses exponentiation"},
		{"Miller-Rabin has no false positives", false, "strong pseudoprimes exist"},
	})
}

func bertrandProblem() generator.Problem {
	return countNTVerdicts("How many of these statements about Bertrand's postulate are true?", []ntVerdict{
		{"For every n > 1 there is a prime between n and 2n", true, "the statement of the postulate"},
		{"Chebyshev proved the postulate", true, "his bounds imply it"},
		{"The postulate has elementary proofs", true, "Erdos gave one"},
		{"It implies pi(2n) - pi(n) >= 1", true, "at least one prime in the interval"},
		{"There is a prime between 10 and 20", true, "11, 13, 17 and 19"},
		{"There is a prime between 100 and 200", true, "101 is one"},
		{"It follows from bounds on theta(x)", true, "the Chebyshev function controls primes"},
		{"Stronger interval results hold for large n", true, "the prime gap is smaller than n"},
		{"It guarantees a prime between n and n+1", false, "that interval often has none"},
		{"It holds only for even n", false, "it holds for every n > 1"},
		{"It is false for large n", false, "it holds for all n > 1"},
		{"It requires the prime number theorem", false, "elementary proofs exist"},
		{"It guarantees exactly one prime in the interval", false, "there may be several"},
		{"It concerns primes between 2n and 3n", false, "the interval is n to 2n"},
		{"n+1 is always prime", false, "not in general"},
		{"Euclid's proof of infinitely many primes implies it", false, "it does not give the interval"},
	})
}
