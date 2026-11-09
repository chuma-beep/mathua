package abstract

// Wave 7 abstract algebra, first batch: semidirect products, group
// representations, symmetric polynomials and Ext/Tor. Every answer is an
// integer computed from published facts.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

func registerWave7(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"abstract.group.semidirect":   "semidirect",
		"abstract.rep.representation": "representation",
		"abstract.poly.symmetric":     "symmetric",
		"abstract.structures.ext":     "ext",
	} {
		reg.Register(id, &w7AbstractGen{kind: kind})
	}
}

type w7AbstractGen struct{ kind string }

func (g *w7AbstractGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "semidirect":
		return w7Semidirect(kp)
	case "representation":
		return w7Representation(kp)
	case "symmetric":
		return w7Symmetric(kp)
	case "ext":
		return w7Ext(kp)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

func w7A(o int, q, a, ex string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}

func w7AI(m map[string]string, k string, v int) { m[k] = strconv.Itoa(v) }

func w7AR(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func w7APhi(n int) int {
	c := 0
	for k := 1; k <= n; k++ {
		if mathutil.GCD(k, n) == 1 {
			c++
		}
	}
	return c
}

func w7Semidirect(kp int) generator.Problem {
	primes := []int{3, 5, 7, 11}
	n := primes[rand.Intn(len(primes))]
	f := map[string]string{}
	w7AI(f, "n", n)
	switch kp {
	case 1:
		n = w7AR(3, 15)
		w7AI(f, "n", n)
		return w7A(kp, fmt.Sprintf("Let \\(N=C_{%d}\\) and \\(H=C_2\\), with \\(H\\) acting by inversion. What is \\(|N\\rtimes H|\\)?", n), strconv.Itoa(2*n),
			fmt.Sprintf("The semidirect product has \\(|N||H|=%d\\cdot2=%d\\) elements.", n, 2*n), f)
	case 2:
		a := w7AR(0, n-1)
		b := w7AR(0, n-1)
		h1 := rand.Intn(2)
		h2 := rand.Intn(2)
		w7AI(f, "a", a)
		w7AI(f, "b", b)
		w7AI(f, "h1", h1)
		w7AI(f, "h2", h2)
		comp := (a + b) % n
		if h1 == 1 {
			comp = (a - b + n) % n
		}
		return w7A(kp, fmt.Sprintf("In \\(C_{%d}\\rtimes C_2\\) with inversion action, multiply \\((a,h_1)=(%d,%d)\\) by \\((b,h_2)=(%d,%d)\\). What is the \\(N\\)-component of the product?", n, a, h1, b, h2), strconv.Itoa(comp),
			fmt.Sprintf("The product is \\((a+(-1)^{h_1}b,h_1+h_2)\\). The \\(N\\)-component is %d.", comp), f)
	default:
		n = w7AR(3, 15)
		w7AI(f, "n", n)
		fixed := 1
		if n%2 == 0 {
			fixed = 2
		}
		return w7A(kp, fmt.Sprintf("The group \\(C_2\\) acts on \\(C_{%d}\\) by inversion. How many elements of \\(C_{%d}\\) are fixed by the inversion?", n, n), strconv.Itoa(fixed),
			fmt.Sprintf("Solve \\(-x=x\\pmod{%d}\\), equivalently \\(2x=0\\). The number of solutions is \\(\\gcd(2,%d)=%d\\).", n, n, fixed), f)
	}
}

func w7Representation(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		n := w7AR(2, 12)
		w7AI(f, "order", n)
		return w7A(kp, fmt.Sprintf("The regular representation of a finite group of order %d acts on the group algebra. What is its degree?", n), strconv.Itoa(n),
			"The regular representation acts on a vector space with one basis element per group element, so its degree is the group order.", f)
	case 2:
		a, b, c := w7AR(1, 5), w7AR(1, 5), w7AR(1, 5)
		w7AI(f, "a", a)
		w7AI(f, "b", b)
		w7AI(f, "c", c)
		return w7A(kp, fmt.Sprintf("A representation is the direct sum of three representations of degrees %d, %d and %d. What is the character of the direct sum at the identity?", a, b, c), strconv.Itoa(a+b+c),
			"The character at the identity is the degree, and degrees add under direct sum.", f)
	default:
		n := w7AR(3, 15)
		w7AI(f, "n", n)
		return w7A(kp, fmt.Sprintf("In the natural permutation representation of \\(S_{%d}\\) on %d points, what is the character value at a transposition?", n, n), strconv.Itoa(n-2),
			fmt.Sprintf("A transposition fixes %d points and swaps two, so the permutation character is %d.", n-2, n-2), f)
	}
}

func w7Symmetric(kp int) generator.Problem {
	a := w7AR(-4, 4)
	b := w7AR(-4, 4)
	c := w7AR(-4, 4)
	f := map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "c": strconv.Itoa(c)}
	switch kp {
	case 1:
		ans := a*b + a*c + b*c
		return w7A(kp, fmt.Sprintf("For variables \\(a=%d\\), \\(b=%d\\) and \\(c=%d\\), what is the elementary symmetric polynomial \\(e_2=ab+ac+bc\\)?", a, b, c), strconv.Itoa(ans),
			fmt.Sprintf("The pairwise products are %d, %d and %d; their sum is %d.", a*b, a*c, b*c, ans), f)
	case 2:
		ans := a*a + b*b + c*c
		return w7A(kp, fmt.Sprintf("For \\(a=%d\\), \\(b=%d\\) and \\(c=%d\\), what is the power sum \\(p_2=a^2+b^2+c^2\\)?", a, b, c), strconv.Itoa(ans),
			fmt.Sprintf("The squares are %d, %d and %d; their sum is %d.", a*a, b*b, c*c, ans), f)
	default:
		ans := a * b * c
		return w7A(kp, fmt.Sprintf("For \\(a=%d\\), \\(b=%d\\) and \\(c=%d\\), what is \\(e_3=abc\\)?", a, b, c), strconv.Itoa(ans),
			fmt.Sprintf("The product is \\((%d)(%d)(%d)=%d\\).", a, b, c, ans), f)
	}
}

func w7Ext(kp int) generator.Problem {
	m := w7AR(2, 9)
	n := w7AR(2, 9)
	f := map[string]string{"m": strconv.Itoa(m), "n": strconv.Itoa(n)}
	switch kp {
	case 1:
		ans := mathutil.GCD(m, n)
		return w7A(kp, fmt.Sprintf("Over \\(\\mathbb{Z}\\), what is the size of \\(\\operatorname{Ext}^1(\\mathbb{Z}/%d\\mathbb{Z},\\mathbb{Z}/%d\\mathbb{Z})\\)?", m, n), strconv.Itoa(ans),
			fmt.Sprintf("The Ext group is cyclic of order \\(\\gcd(%d,%d)=%d\\).", m, n, ans), f)
	case 2:
		ans := mathutil.GCD(m, n)
		return w7A(kp, fmt.Sprintf("Over \\(\\mathbb{Z}\\), what is the size of \\(\\operatorname{Tor}_1(\\mathbb{Z}/%d\\mathbb{Z},\\mathbb{Z}/%d\\mathbb{Z})\\)?", m, n), strconv.Itoa(ans),
			fmt.Sprintf("The Tor group is cyclic of order \\(\\gcd(%d,%d)=%d\\).", m, n, ans), f)
	default:
		a, b := w7AR(1, 6), w7AR(1, 6)
		f["dimA"] = strconv.Itoa(a)
		f["dimC"] = strconv.Itoa(b)
		return w7A(kp, fmt.Sprintf("For a short exact sequence \\(0\\to A\\to B\\to C\\to 0\\) of vector spaces with \\(\\dim A=%d\\) and \\(\\dim C=%d\\), what is \\(\\dim B\\)?", a, b), strconv.Itoa(a+b),
			"Dimensions add in a short exact sequence of vector spaces.", f)
	}
}
