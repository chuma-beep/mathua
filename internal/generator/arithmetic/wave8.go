package arithmetic

// Wave 8 arithmetic: the first three review concepts. Every answer is an
// integer and every operand is published in Facts so wave8_test.go can
// recompute it independently.

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave8(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"arith.place.binary":          "binary",
		"arith.ops.making_ten":        "making_ten",
		"arith.mult.partial_products": "partial_products",
	} {
		reg.Register(id, &w8ArithGen{kind: kind})
	}
}

type w8ArithGen struct{ kind string }

func (g *w8ArithGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "binary":
		return w8Binary(kp)
	case "making_ten":
		return w8MakingTen(kp)
	case "partial_products":
		return w8PartialProducts(kp)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

func w8A(o int, q, a, ex string, f map[string]string) generator.Problem {
	f["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: f, KPOrdinal: o}
}

func w8AI(f map[string]string, k string, v int) { f[k] = strconv.Itoa(v) }

func w8AR(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func w8BinaryString(n int) string {
	if n == 0 {
		return "0"
	}
	var b strings.Builder
	// Build from the highest power down.
	start := 1
	for start<<1 <= n {
		start <<= 1
	}
	for p := start; p > 0; p >>= 1 {
		if n&p != 0 {
			b.WriteByte('1')
		} else {
			b.WriteByte('0')
		}
	}
	return b.String()
}

func w8Binary(kp int) generator.Problem {
	n := w8AR(5, 255)
	bin := w8BinaryString(n)
	f := map[string]string{"binary": bin}
	w8AI(f, "n", n)
	switch kp {
	case 1:
		return w8A(kp, fmt.Sprintf("What is the binary numeral \\(%s_2\\) in decimal? (enter a number)", bin), strconv.Itoa(n),
			fmt.Sprintf("The place values from the left are powers of two. The digits of %s add to %d.", bin, n), f)
	case 2:
		p := 1
		for p<<1 <= n {
			p <<= 1
		}
		w8AI(f, "power", p)
		return w8A(kp, fmt.Sprintf("What is the largest power of two that is at most %d? (enter a number)", n), strconv.Itoa(p),
			fmt.Sprintf("Double 1 until the next double would exceed %d. The largest power is %d.", n, p), f)
	default:
		bits := 0
		for m := n; m > 0; m >>= 1 {
			bits += m & 1
		}
		w8AI(f, "ones", bits)
		return w8A(kp, fmt.Sprintf("How many 1 bits are in the binary expansion of %d? (enter a number)", n), strconv.Itoa(bits),
			fmt.Sprintf("The binary expansion is %s, which contains %d ones.", bin, bits), f)
	}
}

func w8MakingTen(kp int) generator.Problem {
	a := w8AR(1, 10)
	b := w8AR(11-a, 18)
	f := map[string]string{}
	w8AI(f, "a", a)
	w8AI(f, "b", b)
	switch kp {
	case 1:
		return w8A(kp, fmt.Sprintf("How many more than %d are needed to make 10? (enter a number)", a), strconv.Itoa(10-a),
			fmt.Sprintf("The missing addend is \\(10-%d=%d\\).", a, 10-a), f)
	case 2:
		return w8A(kp, fmt.Sprintf("What is \\(%d+%d\\)? (enter a number)", a, b), strconv.Itoa(a+b),
			fmt.Sprintf("Make ten first: \\(%d+%d=10\\), then add the remaining \\(%d-%d=%d\\). The total is %d.", a, 10-a, b, 10-a, b-(10-a), a+b), f)
	default:
		return w8A(kp, fmt.Sprintf("After making ten from %d, how much of %d remains above ten? (enter a number)", a, b), strconv.Itoa(a+b-10),
			fmt.Sprintf("The amount above ten is \\(%d+%d-10=%d\\).", a, b, a+b-10), f)
	}
}

func w8PartialProducts(kp int) generator.Problem {
	t := w8AR(1, 8)
	u := w8AR(1, 9)
	v := w8AR(1, 8)
	w := w8AR(1, 9)
	a := 10*t + u
	b := 10*v + w
	f := map[string]string{}
	w8AI(f, "a", a)
	w8AI(f, "b", b)
	w8AI(f, "t", t)
	w8AI(f, "u", u)
	w8AI(f, "v", v)
	w8AI(f, "w", w)
	switch kp {
	case 1:
		return w8A(kp, fmt.Sprintf("For \\(%d\\times%d\\), what is the tens partial product \\(%d0\\times%d\\)? (enter a number)", a, b, t, b), strconv.Itoa(10*t*b),
			fmt.Sprintf("The tens part of %d is %d0, so the partial product is \\(%d0\\times%d=%d\\).", a, t, t, b, 10*t*b), f)
	case 2:
		return w8A(kp, fmt.Sprintf("For \\(%d\\times%d\\), what is the ones partial product \\(%d\\times%d\\)? (enter a number)", a, b, u, b), strconv.Itoa(u*b),
			fmt.Sprintf("The ones digit of %d is %d, so the partial product is \\(%d\\times%d=%d\\).", a, u, u, b, u*b), f)
	default:
		return w8A(kp, fmt.Sprintf("What is \\(%d\\times%d\\)? (enter a number)", a, b), strconv.Itoa(a*b),
			fmt.Sprintf("Add the partial products: \\(%d0\\times%d=%d\\) and \\(%d\\times%d=%d\\), giving %d.", t, b, 10*t*b, u, b, u*b, a*b), f)
	}
}
