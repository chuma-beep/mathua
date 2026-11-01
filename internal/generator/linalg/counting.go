package linalg

// The structural linear-algebra concepts used to ship as multiple-choice
// recognition questions. They now ask the learner to count how many of eight
// statements are true, which keeps the same ideas but requires a decision on
// every statement instead of one option.

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

type linalgVerdict struct {
	label string
	yes   bool
	note  string
}

const linalgDraw = 8

func countLinalgVerdicts(lead string, pool []linalgVerdict) generator.Problem {
	var yes, no []linalgVerdict
	for _, v := range pool {
		if v.yes {
			yes = append(yes, v)
		} else {
			no = append(no, v)
		}
	}

	// Draw the answer first so every count from 0 to 8 is equally likely; then
	// take that many true statements and fill the rest with false ones.
	answer := rand.Intn(linalgDraw + 1)
	rand.Shuffle(len(yes), func(i, j int) { yes[i], yes[j] = yes[j], yes[i] })
	rand.Shuffle(len(no), func(i, j int) { no[i], no[j] = no[j], no[i] })

	var chosen []linalgVerdict
	chosen = append(chosen, yes[:answer]...)
	chosen = append(chosen, no[:linalgDraw-answer]...)
	rand.Shuffle(len(chosen), func(i, j int) { chosen[i], chosen[j] = chosen[j], chosen[i] })

	var question string
	for attempt := 0; attempt < 40; attempt++ {
		parts := make([]string, len(chosen))
		for i, v := range chosen {
			parts[i] = fmt.Sprintf("(%c) %s", 'a'+i, v.label)
		}
		question = fmt.Sprintf("%s\n\n%s\n\n(enter a number)", lead, strings.Join(parts, "\n"))
		if !strings.Contains(question, strconv.Itoa(answer)) {
			break
		}
		answer = rand.Intn(linalgDraw + 1)
		chosen = chosen[:0]
		chosen = append(chosen, yes[:answer]...)
		chosen = append(chosen, no[:linalgDraw-answer]...)
		rand.Shuffle(len(chosen), func(i, j int) { chosen[i], chosen[j] = chosen[j], chosen[i] })
	}

	notes := make([]string, len(chosen))
	for i, v := range chosen {
		verdict := "false"
		if v.yes {
			verdict = "true"
		}
		notes[i] = fmt.Sprintf("%s: %s, %s", v.label, verdict, v.note)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: fmt.Sprintf("Exactly %d of the eight statements are true. %s", answer, strings.Join(notes, " ")),
		Facts:       map[string]string{"count": strconv.Itoa(answer)},
	}
}

func spaceAxiomsProblem() generator.Problem {
	return countLinalgVerdicts("How many of these statements about real vector spaces are true?", []linalgVerdict{
		{"A vector space is closed under addition and scalar multiplication", true, "both closures are axioms"},
		{"The zero vector of a vector space is unique", true, "two additive identities must be equal"},
		{"Every vector has a unique additive inverse", true, "the inverse is forced by the axioms"},
		{"Scalar multiplication distributes over vector addition", true, "c(u+v) = cu + cv is an axiom"},
		{"1 times v equals v for every vector v", true, "the field unit acts as an identity"},
		{"A subset is a subspace exactly when it contains 0 and is closed under addition and scalar multiplication", true, "this is the subspace test"},
		{"Vector addition is commutative", true, "a vector space is an abelian group under addition"},
		{"The intersection of two subspaces is a subspace", true, "both closures and zero pass through the intersection"},
		{"A vector space must be closed under addition only", false, "scalar multiplication is also required"},
		{"A vector space can have two different zero vectors", false, "additive identities are unique"},
		{"Scalar multiplication need not distribute over vector addition", false, "distributivity is an axiom"},
		{"Every subset containing the zero vector is a subspace", false, "closedness under the operations is also needed"},
		{"The union of two subspaces is always a subspace", false, "the union usually misses sums across the two subspaces"},
		{"A vector space over R must be finite dimensional", false, "polynomials and functions are infinite dimensional examples"},
		{"Vector addition need not be associative", false, "associativity is an axiom"},
		{"Multiplication by 0 can produce a nonzero vector", false, "0v is always the zero vector"},
	})
}

func invertibleMapProblem() generator.Problem {
	return countLinalgVerdicts("How many of these statements about invertible linear maps are true?", []linalgVerdict{
		{"A square matrix is invertible exactly when its determinant is nonzero", true, "the determinant is the product of the eigenvalues"},
		{"A linear map has a trivial kernel exactly when it is injective", true, "T(u) = T(v) means T(u-v) = 0"},
		{"A square matrix that is injective is also surjective", true, "rank n fills an n-dimensional codomain"},
		{"An n by n matrix is invertible exactly when its rank is n", true, "full rank is equivalent to invertibility"},
		{"A square matrix is invertible exactly when zero is not an eigenvalue", true, "a zero eigenvalue gives a nontrivial kernel"},
		{"The inverse of an invertible linear map is linear", true, "the inverse matrix represents it"},
		{"The columns of an invertible matrix form a basis", true, "independence and spanning both hold"},
		{"The inverse of AB is B inverse times A inverse", true, "the order reverses"},
		{"A determinant of zero implies the matrix is invertible", false, "zero determinant means singular"},
		{"A linear map from R^3 to R^2 can be injective", false, "rank-nullity forces a nontrivial kernel"},
		{"A square matrix with a zero eigenvalue is invertible", false, "the zero eigenvalue gives a nonzero kernel vector"},
		{"An injective square map need not be surjective", false, "finite dimension makes injective square maps surjective"},
		{"The inverse of a linear map is always nonlinear", false, "inverses of linear maps are linear"},
		{"A matrix with repeated eigenvalues is never invertible", false, "the identity has repeated eigenvalues and is invertible"},
		{"Every square matrix is invertible", false, "singular matrices exist"},
		{"Row reduction can make a singular matrix invertible", false, "elementary operations preserve singularity"},
	})
}

func spectralProblem() generator.Problem {
	return countLinalgVerdicts("How many of these statements about the spectral theorem are true?", []linalgVerdict{
		{"A real symmetric matrix has only real eigenvalues", true, "symmetry rules out nonreal eigenvalues"},
		{"A real symmetric matrix has an orthonormal basis of eigenvectors", true, "this is the content of the spectral theorem"},
		{"A real symmetric matrix can be written Q D Q transpose with Q orthogonal", true, "orthogonal diagonalization"},
		{"Eigenvectors belonging to distinct eigenvalues of a symmetric matrix are orthogonal", true, "symmetry forces the inner product to vanish"},
		{"A symmetric matrix equals its transpose", true, "that is the definition"},
		{"The spectral theorem applies to every real symmetric matrix", true, "including matrices with repeated eigenvalues"},
		{"A symmetric matrix with repeated eigenvalues can still be diagonalized", true, "the eigenspaces are complete"},
		{"A symmetric matrix with eigenvalues 2 and 3 is positive definite", true, "both eigenvalues are positive"},
		{"Every symmetric matrix is invertible", false, "the zero matrix is symmetric and singular"},
		{"Eigenvectors of a symmetric matrix are always parallel", false, "they are orthogonal for distinct eigenvalues"},
		{"A symmetric matrix must have distinct eigenvalues", false, "the identity is a counterexample"},
		{"A real symmetric matrix may have nonreal eigenvalues", false, "its eigenvalues are always real"},
		{"Only diagonal matrices are orthogonally diagonalizable", false, "every symmetric matrix is"},
		{"Every symmetric matrix is positive definite", false, "the zero matrix is not"},
		{"The spectral theorem requires the matrix to be triangular", false, "symmetry is the hypothesis"},
		{"A symmetric matrix cannot have a repeated eigenvalue", false, "the identity repeats its eigenvalue"},
	})
}

func definiteProblem() generator.Problem {
	return countLinalgVerdicts("How many of these statements about quadratic forms are true?", []linalgVerdict{
		{"A symmetric matrix is positive definite exactly when all its eigenvalues are positive", true, "the form is a sum of eigenvalue-weighted squares"},
		{"The form x^2 + y^2 is positive definite", true, "it vanishes only at the origin"},
		{"The form x^2 - y^2 is indefinite", true, "it takes both signs"},
		{"A 2 by 2 symmetric matrix is positive definite exactly when a is positive and ac - b^2 is positive", true, "Sylvester's criterion for the leading minors"},
		{"The level sets of a positive definite form are ellipses", true, "the form is bounded below away from the origin"},
		{"A symmetric matrix is negative definite exactly when all its eigenvalues are negative", true, "apply the positive definite test to the negative"},
		{"The identity matrix is positive definite", true, "its eigenvalues are all 1"},
		{"A diagonal matrix with positive diagonal entries is positive definite", true, "it is a sum of positive squares"},
		{"A symmetric matrix is positive definite exactly when all its entries are positive", false, "off-diagonal signs can still leave the form indefinite"},
		{"A determinant of zero implies the matrix is positive definite", false, "a zero eigenvalue rules out definiteness"},
		{"A positive definite matrix can have a negative eigenvalue", false, "all eigenvalues must be positive"},
		{"The form x^2 + y^2 is indefinite", false, "it never takes negative values"},
		{"The form x^2 - y^2 is positive definite", false, "it takes negative values"},
		{"A positive determinant alone implies a 2 by 2 matrix is positive definite", false, "both leading minors must be positive"},
		{"Positive definite matrices are always singular", false, "positive eigenvalues make the matrix invertible"},
		{"The level sets of a positive definite form are hyperbolas", false, "they are bounded closed curves"},
	})
}

func cayleyHamiltonProblem() generator.Problem {
	return countLinalgVerdicts("How many of these statements about the Cayley-Hamilton theorem are true?", []linalgVerdict{
		{"Every square matrix satisfies its own characteristic polynomial", true, "that is the theorem"},
		{"The characteristic polynomial of an n by n matrix has degree n", true, "the determinant expands with n factors"},
		{"A matrix with characteristic polynomial lambda^2 - 5 lambda + 6 satisfies A^2 - 5A + 6I = 0", true, "substitute A for lambda"},
		{"The minimal polynomial divides the characteristic polynomial", true, "Cayley-Hamilton forces the divisibility"},
		{"Cayley-Hamilton reduces powers of A to combinations of I through A to the n-1", true, "the relation solves for the top power"},
		{"Cayley-Hamilton holds for matrices that are not diagonalizable", true, "no diagonalizability hypothesis"},
		{"Cayley-Hamilton can compute an inverse when the determinant is nonzero", true, "factor out A from the relation"},
		{"The trace and determinant appear in the characteristic polynomial coefficients", true, "up to signs and position"},
		{"Cayley-Hamilton says every matrix is diagonalizable", false, "it says nothing about diagonalizability"},
		{"The minimal polynomial always equals the characteristic polynomial", false, "it only divides it in general"},
		{"The characteristic polynomial of a 3 by 3 matrix has degree 9", false, "the degree equals the size, 3"},
		{"Cayley-Hamilton only holds for symmetric matrices", false, "it holds for every square matrix"},
		{"Cayley-Hamilton implies every matrix is invertible", false, "singular matrices satisfy it too"},
		{"The minimal polynomial has larger degree than the characteristic polynomial", false, "it divides the characteristic polynomial"},
		{"Cayley-Hamilton says A equals its transpose", false, "that is symmetry, not the theorem"},
		{"The characteristic polynomial degree equals the number of nonzero entries", false, "it equals the matrix size"},
	})
}
