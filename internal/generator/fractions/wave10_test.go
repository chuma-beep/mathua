package fractions

import (
	"math"
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

func generatorContext(kp int) generator.GeneratorContext {
	return generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp}
}

func w10FactInt(t *testing.T, facts map[string]string, key string) int {
	t.Helper()
	v, err := strconv.Atoi(facts[key])
	if err != nil {
		t.Fatalf("fact %q = %q is not an integer", key, facts[key])
	}
	return v
}

func w10RoundDiv(num, den int) int {
	return (num + den/2) / den
}

// w10Recompute derives the answer from the published facts alone, so a change
// in a generator cannot silently pass with a stale answer.
func w10Recompute(kind string, kp int, f map[string]string) int {
	get := func(k string) int {
		v, _ := strconv.Atoi(f[k])
		return v
	}
	a, b, c, d, e, fct := get("a"), get("b"), get("c"), get("d"), get("e"), get("f")
	switch kind {
	case "density":
		den := get("d")
		if kp == 1 {
			return den - 1
		}
		if kp == 2 {
			return den - 1 - den/2
		}
		return den/2 - 1
	case "whole_fraction":
		if kp == 2 {
			return get("n") / get("w")
		}
		return get("w") * get("d")
	case "proper_improper":
		if kp == 1 {
			if get("n") < get("d") {
				return 1
			}
			return 0
		}
		if kp == 2 {
			return get("n") / get("d")
		}
		return get("w")*get("b") + get("a")
	case "zero_denominator":
		if kp == 1 {
			return get("zeros")
		}
		if kp == 2 {
			return get("k")
		}
		return get("a")
	case "equal_parts":
		if kp == 1 {
			return get("d")
		}
		if kp == 2 {
			return get("d") - get("s")
		}
		return get("w") * get("n")
	case "unit_meaning":
		if kp == 3 {
			return get("w") * get("n")
		}
		return get("total") / get("n")
	case "cross_check":
		if kp == 1 {
			if get("a")*get("d") == get("c")*get("b") {
				return 1
			}
			return 0
		}
		if kp == 2 {
			return get("a") * get("d")
		}
		return get("b") * get("c") / get("a")
	case "simplify_gcf":
		g := mathutil.GCD(get("n"), get("d"))
		if kp == 1 {
			return g
		}
		if kp == 2 {
			return get("n") / g
		}
		return get("d") / g
	case "lowest_terms":
		g := mathutil.GCD(get("n"), get("d"))
		if kp == 1 {
			if g == 1 {
				return 1
			}
			return 0
		}
		if kp == 2 {
			return g
		}
		return get("n") / g
	case "equiv_chain", "missing_number":
		if kp == 2 {
			return get("b") * get("k")
		}
		return get("a") * get("k")
	case "benchmark_choice":
		n, den := get("n"), get("d")
		if kp == 1 {
			if 2*n > den {
				return 1
			}
			return 0
		}
		if kp == 2 {
			d0 := n
			d1 := 2*n - den
			if d1 < 0 {
				d1 = -d1
			}
			d2 := den - n
			if d1 <= d0 && d1 <= d2 {
				return 1
			}
			if d2 < d0 && d2 < d1 {
				return 2
			}
			return 0
		}
		gap := n - 4
		if gap < 0 {
			gap = -gap
		}
		return gap
	case "decimal_order":
		a, b, c := get("a"), get("b"), get("c")
		if kp == 1 {
			count := 0
			for _, v := range []int{a, b, c} {
				if v < 5 {
					count++
				}
			}
			return count
		}
		if kp == 2 {
			idx := 1
			if b > a {
				idx = 2
			}
			if c > []int{a, b}[idx-1] {
				idx = 3
			}
			return idx
		}
		rank := 1
		for _, v := range []int{a, b} {
			if v < c {
				rank++
			}
		}
		return rank
	case "percent_bench":
		if kp == 2 {
			p := get("p")
			return p / mathutil.GCD(p, 100)
		}
		return get("a") * 100 / get("b")
	case "add_whole":
		w, a, b := get("w"), get("a"), get("b")
		if kp == 1 {
			return w*b + a
		}
		if kp == 2 {
			return w
		}
		return w*b + 2*b + a
	case "estimate_sum":
		num := a*d + c*b
		den := b * d
		if kp == 2 {
			if num > den {
				return 1
			}
			return 0
		}
		return w10RoundDiv(num, den)
	case "add_regroup":
		if kp == 1 {
			return (a + c) / b
		}
		if kp == 2 {
			return (a + c) % b
		}
		return a + c
	case "add_unit":
		if kp == 1 {
			return b/a + 1
		}
		if kp == 2 {
			return b
		}
		return get("n")
	case "missing_addend":
		if kp == 1 {
			return get("c") - a
		}
		if kp == 2 {
			return get("c")*b - a*get("d")
		}
		return get("w")
	case "add_associative":
		total := a + c + e
		if kp == 1 {
			return total
		}
		if kp == 2 {
			return total / b
		}
		return a + c
	case "add_same_numerator":
		l := lcm(b, get("d"))
		num := a * (l/b + l/get("d"))
		if kp == 1 {
			return num
		}
		if kp == 2 {
			return l
		}
		if num > l {
			return 1
		}
		return 0
	case "add_decimal":
		sum := get("a")*100/get("b") + get("dec")
		if kp == 1 {
			return sum
		}
		if kp == 2 {
			return sum / mathutil.GCD(sum, 100)
		}
		if sum > 100 {
			return 1
		}
		return 0
	case "missing_minuend":
		if kp == 1 {
			return a + get("c")
		}
		if kp == 2 {
			return a*get("d") + get("c")*b
		}
		return get("w")*b + a
	case "missing_subtrahend":
		if kp == 1 {
			return a - get("c")
		}
		if kp == 2 {
			return a*get("d") - get("c")*b
		}
		return a - get("w")*b
	case "estimate_diff":
		num := a*get("d") - get("c")*b
		if num < 0 {
			num = -num
		}
		den := b * get("d")
		if kp == 1 {
			if num*2 > den {
				return 2
			}
			if num*2 == den {
				return 1
			}
			return 0
		}
		if kp == 2 {
			return w10RoundDiv(num, den)
		}
		if num*2 < den {
			return 1
		}
		return 0
	case "sub_regroup_whole":
		if kp == 1 {
			return get("w") - 1
		}
		if kp == 2 {
			return b - a
		}
		return get("w")*b - a
	case "sub_negative":
		if kp == 1 {
			return a*get("d") - get("c")*b
		}
		if kp == 2 {
			return 1
		}
		return get("c")*b - a*get("d")
	case "sub_three":
		num := a - get("c") - e
		if kp == 1 {
			return num
		}
		if kp == 2 {
			return num / b
		}
		return num % b
	case "sub_decimal":
		diff := get("a")*100/get("b") - get("dec")
		if kp == 1 {
			return diff
		}
		if kp == 2 {
			return diff / mathutil.GCD(diff, 100)
		}
		if diff > 0 {
			return 1
		}
		return 0
	case "sub_distance":
		l := lcm(b, get("d"))
		if kp == 1 {
			return l
		}
		num := a*(l/b) - get("c")*(l/get("d"))
		if num < 0 {
			num = -num
		}
		if kp == 2 {
			return num
		}
		if num*2 > l {
			return 1
		}
		return 0
	case "mult_whole_cancel":
		n := get("w") * a
		g := mathutil.GCD(n, b)
		if kp == 1 {
			return n / g
		}
		if kp == 2 {
			return b / g
		}
		return n / b
	case "estimate_product":
		num, den := a*get("c"), b*get("d")
		if kp == 1 {
			if num*2 > den {
				return 2
			}
			if num*2 == den {
				return 1
			}
			return 0
		}
		if kp == 2 {
			if num*2 < den {
				return 1
			}
			return 0
		}
		return w10RoundDiv(num, den)
	case "mult_scaling":
		n := get("w") * a
		g := mathutil.GCD(n, b)
		if kp == 1 {
			return n / g
		}
		if kp == 2 {
			return n / b
		}
		if a > b {
			return 1
		}
		return 0
	case "mult_unit":
		g := mathutil.GCD(a, b*get("c"))
		if kp == 1 {
			return b * get("c") / g
		}
		if kp == 2 {
			return a / g
		}
		return b * get("c")
	case "mult_three":
		num, den := a*c*e, b*get("d")*fct
		g := mathutil.GCD(num, den)
		if kp == 1 {
			return num / g
		}
		if kp == 2 {
			return den / g
		}
		if num*2 < den {
			return 1
		}
		return 0
	case "mult_mixed_improper":
		imp := get("w")*b + a
		if kp == 1 {
			return imp
		}
		if kp == 2 {
			return imp * get("c")
		}
		return imp * get("c") / (b * get("d"))
	case "mult_decimal":
		num, den := a*get("dec"), b*100
		g := mathutil.GCD(num, den)
		if kp == 1 {
			return num / g
		}
		if kp == 2 {
			return den / g
		}
		if num*2 > den {
			return 1
		}
		return 0
	case "mult_predict":
		if kp == 1 {
			if get("c") > get("d") {
				return 1
			}
			return 0
		}
		if kp == 2 {
			if a*get("c") > b*get("d") {
				return 1
			}
			return 0
		}
		if a*get("c") > b*get("d") {
			return 2
		}
		if a*get("c") == b*get("d") {
			return 1
		}
		return 0
	case "div_reciprocal":
		if kp == 1 {
			return a * get("d")
		}
		if kp == 2 {
			return b * get("c")
		}
		return a * get("d") / (b * get("c"))
	case "estimate_quotient":
		num, den := a*get("d"), b*get("c")
		if kp == 1 {
			return w10RoundDiv(num, den)
		}
		if kp == 2 {
			if num > den {
				return 1
			}
			return 0
		}
		if num*2 > den {
			return 2
		}
		if num*2 == den {
			return 1
		}
		return 0
	case "complex_simplify":
		num, den := a*get("d"), b*get("c")
		g := mathutil.GCD(num, den)
		if kp == 1 {
			return num / g
		}
		if kp == 2 {
			return den / g
		}
		return num / den
	case "div_unit":
		n := a * get("c")
		g := mathutil.GCD(n, b)
		if kp == 1 {
			return n / g
		}
		if kp == 2 {
			return b / g
		}
		return n
	case "div_measurement":
		total := get("w")*b + a
		c := get("c")
		if kp == 2 {
			return total * c / (2 * b)
		}
		return total * c / b
	case "div_partitive":
		g := mathutil.GCD(a, b*get("c"))
		if kp == 1 {
			return a / g
		}
		if kp == 2 {
			return b * get("c") / g
		}
		return b * get("c")
	case "div_three":
		num, den := a*get("d")*fct, b*get("c")*e
		g := mathutil.GCD(num, den)
		if kp == 1 {
			return num / g
		}
		if kp == 2 {
			return den / g
		}
		if num > den {
			return 1
		}
		return 0
	case "div_zero":
		if kp == 1 {
			return get("zeros")
		}
		if kp == 2 {
			return get("b")
		}
		return 0
	case "mixed_compare":
		if kp == 1 {
			if get("w2") > get("w1") {
				return 2
			}
			return 1
		}
		if kp == 2 {
			left := (get("w1")*b + a) * get("d")
			right := (get("w2")*get("d") + get("c")) * b
			if right > left {
				return 2
			}
			return 1
		}
		diff := get("w1") - get("w2")
		if diff < 0 {
			diff = -diff
		}
		return diff
	case "mixed_negative":
		imp := get("w")*b + a
		if kp == 1 {
			return imp
		}
		if kp == 2 {
			return -imp
		}
		return get("w")
	case "mixed_borrow_twice":
		whole := get("w1") - get("w2") - 1
		num := a + b - get("c")
		if kp == 1 {
			return whole
		}
		if kp == 2 {
			return num
		}
		return (get("w1")-get("w2"))*b + num
	case "mixed_div_whole":
		imp := get("w")*b + a
		den := b * get("c")
		g := mathutil.GCD(imp, den)
		if kp == 1 {
			return imp / g
		}
		if kp == 2 {
			return den / g
		}
		return imp / den
	case "mixed_estimate":
		sum := float64(get("w1")+get("w2")) + float64(a)/float64(b) + float64(get("c"))/float64(get("d"))
		diff := float64(get("w1")-get("w2")) + float64(a)/float64(b) - float64(get("c"))/float64(get("d"))
		if diff < 0 {
			diff = -diff
		}
		if kp == 1 {
			return int(math.Floor(sum + 0.5))
		}
		if kp == 2 {
			return int(math.Floor(diff + 0.5))
		}
		if sum > float64(get("w1")+get("w2")+1) {
			return 1
		}
		return 0
	case "mixed_percent":
		if kp == 2 {
			return get("a") * 100 / b
		}
		return (get("w")*b + a) * 100 / b
	case "mixed_word_compare":
		if kp == 3 {
			diff := get("w1") - get("w2")
			if diff < 0 {
				diff = -diff
			}
			return diff
		}
		left := float64(get("w1")) + float64(a)/float64(b)
		right := float64(get("w2")) + float64(get("c"))/float64(get("d"))
		if right > left {
			return 2
		}
		return 1
	}
	return -999999
}

func TestWave10Answers(t *testing.T) {
	for id, kind := range w10Kinds {
		gen := &w10Gen{kind: kind}
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 20; draw++ {
				p := gen.Generate(generatorContext(kp))
				if p.KPOrdinal != kp {
					t.Fatalf("%s kp%d: ordinal = %d", id, kp, p.KPOrdinal)
				}
				answer, err := strconv.Atoi(p.Answer)
				if err != nil {
					t.Fatalf("%s kp%d: answer %q is not an integer", id, kp, p.Answer)
				}
				if p.Facts["answer"] != p.Answer {
					t.Fatalf("%s kp%d: fact answer %q differs from answer %q", id, kp, p.Facts["answer"], p.Answer)
				}
				want := w10Recompute(kind, kp, p.Facts)
				if answer != want {
					t.Fatalf("%s kp%d: answer %d, want %d\nquestion: %s\nfacts: %v", id, kp, answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func TestWave10Diversity(t *testing.T) {
	for id, kind := range w10Kinds {
		gen := &w10Gen{kind: kind}
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 100; i++ {
				seen[gen.Generate(generatorContext(kp)).Question] = true
			}
			if len(seen) < 5 {
				t.Fatalf("%s kp%d: only %d distinct questions in 100 draws", id, kp, len(seen))
			}
		}
	}
}
