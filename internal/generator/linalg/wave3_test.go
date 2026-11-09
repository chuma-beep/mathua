package linalg

import (
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func w3TestInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func w3TestMatrix(t *testing.T, p generator.Problem, prefix string) (int, int, int, int) {
	t.Helper()
	return w3TestInt(t, p, prefix+"a"), w3TestInt(t, p, prefix+"b"), w3TestInt(t, p, prefix+"c"), w3TestInt(t, p, prefix+"d")
}

func w3TestMul(a, b, c, d, e, f, g, h int) (int, int, int, int) {
	return a*e + b*g, a*f + b*h, c*e + d*g, c*f + d*h
}

func w3TestPow(a, b, c, d, k int) (int, int, int, int) {
	ra, rb, rc, rd := 1, 0, 0, 1
	for i := 0; i < k; i++ {
		ra, rb, rc, rd = w3TestMul(ra, rb, rc, rd, a, b, c, d)
	}
	return ra, rb, rc, rd
}

func w3TestDet(a, b, c, d int) int { return a*d - b*c }

func w3TestAbs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func TestWave3LinalgAnswerRecomputation(t *testing.T) {
	kinds := []string{
		"powers", "nilpotent", "idempotent", "elementary", "block", "permutation",
		"norm", "pseudoinverse", "cholesky", "schur_complement", "trace_identities",
		"rank_factorization", "outer_product", "hadamard", "kronecker", "vandermonde",
		"companion", "circulant", "householder", "givens", "determinant_lemma",
	}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 30; draw++ {
				p := (&w3LinalgGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				want := w3TestLinalgAnswer(t, kind, kp, p)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func w3TestLinalgAnswer(t *testing.T, kind string, kp int, p generator.Problem) int {
	t.Helper()
	switch kind {
	case "powers":
		a, b, c, d := w3TestMatrix(t, p, "")
		k := 2
		if kp == 2 {
			k = 3
		} else if kp == 3 {
			k = w3TestInt(t, p, "power")
		}
		ra, rb, _, rd := w3TestPow(a, b, c, d, k)
		if kp == 1 {
			return ra
		}
		if kp == 2 {
			return ra + rd
		}
		return rb
	case "nilpotent":
		a, b, c, d := w3TestMatrix(t, p, "")
		ra, rb, rc, rd := w3TestPow(a, b, c, d, 2)
		if ra != 0 || rb != 0 || rc != 0 || rd != 0 {
			t.Fatalf("nilpotent matrix does not square to zero: %v", p.Facts)
		}
		if kp == 1 {
			return ra
		}
		if kp == 2 {
			if a == 0 && b == 0 && c == 0 && d == 0 {
				return 1
			}
			return 2
		}
		return 1
	case "idempotent":
		a, b, c, d := w3TestMatrix(t, p, "")
		ra, rb, rc, rd := w3TestPow(a, b, c, d, 2)
		if ra != a || rb != b || rc != c || rd != d {
			t.Fatalf("idempotent matrix does not square to itself: %v", p.Facts)
		}
		if kp == 1 {
			return rd
		}
		if kp == 2 {
			return a + d
		}
		return 1
	case "elementary":
		a, b, c, d := w3TestMatrix(t, p, "")
		if kp == 1 {
			return w3TestDet(a, b, c, d)
		}
		if kp == 2 {
			return c
		}
		return a + d
	case "block":
		a, b, c, d := w3TestMatrix(t, p, "A")
		e, f, g, h := w3TestMatrix(t, p, "B")
		if kp == 1 {
			return w3TestDet(a, b, c, d) * w3TestDet(e, f, g, h)
		}
		if kp == 2 {
			return a + d + e + h
		}
		ra, _, _, _ := w3TestMul(a, b, c, d, a, b, c, d)
		return ra
	case "permutation":
		permParts := strings.Split(p.Facts["permutation"], ",")
		perm := make([]int, len(permParts))
		for i, s := range permParts {
			perm[i], _ = strconv.Atoi(s)
		}
		if kp == 1 {
			inv := 0
			for i := 0; i < len(perm); i++ {
				for j := i + 1; j < len(perm); j++ {
					if perm[i] > perm[j] {
						inv++
					}
				}
			}
			if inv%2 == 0 {
				return 1
			}
			return -1
		}
		if kp == 2 {
			inv := 0
			for i := 0; i < len(perm); i++ {
				for j := i + 1; j < len(perm); j++ {
					if perm[i] > perm[j] {
						inv++
					}
				}
			}
			return inv
		}
		row := w3TestInt(t, p, "row")
		vecParts := strings.Split(p.Facts["vector"], ",")
		vec := make([]int, len(vecParts))
		for i, s := range vecParts {
			vec[i], _ = strconv.Atoi(s)
		}
		return vec[perm[row]]
	case "norm":
		a, b, c, d := w3TestMatrix(t, p, "")
		if kp == 1 {
			col1, col2 := w3TestAbs(a)+w3TestAbs(c), w3TestAbs(b)+w3TestAbs(d)
			if col1 > col2 {
				return col1
			}
			return col2
		}
		if kp == 2 {
			row1, row2 := w3TestAbs(a)+w3TestAbs(b), w3TestAbs(c)+w3TestAbs(d)
			if row1 > row2 {
				return row1
			}
			return row2
		}
		return a*a + b*b + c*c + d*d
	case "pseudoinverse":
		if kp == 1 {
			ans, err := strconv.Atoi(p.Answer)
			if err != nil {
				t.Fatalf("pseudoinverse answer %q is not an integer", p.Answer)
			}
			if ans == 0 {
				return 0
			}
			for i := 1; i <= 3; i++ {
				v := w3TestInt(t, p, "d"+strconv.Itoa(i))
				if v != 0 && v*ans == 1 {
					return ans
				}
			}
			t.Fatalf("pseudoinverse answer %d is not the reciprocal of a diagonal entry: %v", ans, p.Facts)
		}
		r := 0
		for i := 1; i <= 3; i++ {
			if w3TestInt(t, p, "d"+strconv.Itoa(i)) != 0 {
				r++
			}
		}
		return r
	case "cholesky":
		a, b, c, d := w3TestMatrix(t, p, "")
		pv := w3TestInt(t, p, "p")
		q := w3TestInt(t, p, "q")
		r := w3TestInt(t, p, "r")
		if a != pv*pv || b != pv*q || c != pv*q || d != q*q+r*r {
			t.Fatalf("Cholesky facts do not match the matrix: %v", p.Facts)
		}
		if kp == 1 {
			return pv
		}
		if kp == 2 {
			return q
		}
		return w3TestDet(a, b, c, d)
	case "schur_complement":
		a, b, c, d := w3TestMatrix(t, p, "")
		comp := d - c*b/a
		if kp == 1 {
			return comp
		}
		if kp == 2 {
			return a * comp
		}
		if comp != 0 {
			return 1
		}
		return 0
	case "trace_identities":
		a, b, c, d := w3TestMatrix(t, p, "A")
		e, f, g, h := w3TestMatrix(t, p, "B")
		if kp == 1 {
			ra, _, _, rd := w3TestMul(a, b, c, d, e, f, g, h)
			return ra + rd
		}
		if kp == 2 {
			ra, _, _, rd := w3TestMul(a, b, c, d, a, b, c, d)
			return ra + rd
		}
		return a + d + e + h
	case "rank_factorization":
		u1, u2 := w3TestInt(t, p, "u1"), w3TestInt(t, p, "u2")
		v1, v2 := w3TestInt(t, p, "v1"), w3TestInt(t, p, "v2")
		if kp == 1 {
			return 1
		}
		if kp == 2 {
			return u1*v1 + u2*v2
		}
		return (u1 + u2) * (v1 + v2)
	case "outer_product":
		u := []int{w3TestInt(t, p, "u1"), w3TestInt(t, p, "u2"), w3TestInt(t, p, "u3")}
		v := []int{w3TestInt(t, p, "v1"), w3TestInt(t, p, "v2"), w3TestInt(t, p, "v3")}
		if kp == 1 {
			i, j := w3TestInt(t, p, "row")-1, w3TestInt(t, p, "col")-1
			return u[i] * v[j]
		}
		if kp == 2 {
			return u[0]*v[0] + u[1]*v[1] + u[2]*v[2]
		}
		return 1
	case "hadamard":
		a, b, c, d := w3TestMatrix(t, p, "A")
		e, f, g, h := w3TestMatrix(t, p, "B")
		if kp == 1 {
			// The question picks a position, but the facts do not publish it.
			// Verify the answer is one of the four corresponding products.
			for _, v := range []int{a * e, b * f, c * g, d * h} {
				if p.Answer == strconv.Itoa(v) {
					return v
				}
			}
			t.Fatalf("Hadamard answer %s is not a product of corresponding entries: %v", p.Answer, p.Facts)
		}
		if kp == 2 {
			return a*e + d*h
		}
		return a*e + b*f + c*g + d*h
	case "kronecker":
		a, _, _, d := w3TestMatrix(t, p, "A")
		e, _, _, h := w3TestMatrix(t, p, "B")
		if kp == 1 {
			return 4
		}
		if kp == 2 {
			return a * e
		}
		return (a + d) * (e + h)
	case "vandermonde":
		x, y, z := w3TestInt(t, p, "x"), w3TestInt(t, p, "y"), w3TestInt(t, p, "z")
		if kp == 1 {
			return (y - x) * (z - x) * (z - y)
		}
		if kp == 2 {
			return z * z
		}
		return 3
	case "companion":
		b, c := w3TestInt(t, p, "b"), w3TestInt(t, p, "c")
		if kp == 1 || kp == 3 {
			return c
		}
		return -b
	case "circulant":
		a, b, c := w3TestInt(t, p, "a"), w3TestInt(t, p, "b"), w3TestInt(t, p, "c")
		if kp == 1 {
			return a*a*a + b*b*b + c*c*c - 3*a*b*c
		}
		if kp == 2 {
			return 3 * a
		}
		return b
	case "householder":
		if kp == 2 {
			return -1
		}
		if kp == 3 {
			return 1
		}
		a, b, c, d := w3TestMatrix(t, p, "")
		x, y := w3TestInt(t, p, "x"), w3TestInt(t, p, "y")
		_ = c
		_ = d
		return a*x + b*y
	case "givens":
		c, s := w3TestInt(t, p, "c"), w3TestInt(t, p, "s")
		x, y := w3TestInt(t, p, "x"), w3TestInt(t, p, "y")
		if kp == 1 {
			return 1
		}
		if kp == 2 {
			return s
		}
		return c*x - s*y
	case "determinant_lemma":
		a, d := w3TestInt(t, p, "a"), w3TestInt(t, p, "d")
		u1, u2 := w3TestInt(t, p, "u1"), w3TestInt(t, p, "u2")
		v1, v2 := w3TestInt(t, p, "v1"), w3TestInt(t, p, "v2")
		scalar := v1*u1/a + v2*u2/d
		detA := a * d
		if kp == 1 {
			return scalar
		}
		if kp == 2 {
			return detA * (1 + scalar)
		}
		return detA * scalar
	}
	t.Fatalf("unknown kind %q", kind)
	return 0
}

func TestWave3LinalgQuestionDiversity(t *testing.T) {
	kinds := []string{
		"powers", "nilpotent", "idempotent", "elementary", "block", "permutation",
		"norm", "pseudoinverse", "cholesky", "schur_complement", "trace_identities",
		"rank_factorization", "outer_product", "hadamard", "kronecker", "vandermonde",
		"companion", "circulant", "householder", "givens", "determinant_lemma",
	}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 60; i++ {
				p := (&w3LinalgGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				seen[p.Question] = true
			}
			if len(seen) < 5 {
				t.Fatalf("%s kp%d: only %d distinct questions in 60 draws", kind, kp, len(seen))
			}
		}
	}
}
