package arithmetic

import (
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

var w12Kinds = []string{
	"arith.add.three_addends", "arith.add.four_digit", "arith.add.regroup_twice",
	"arith.sub.four_digit", "arith.sub.regroup_twice", "arith.sub.missing_subtrahend",
	"arith.mult.three_by_two", "arith.mult.powers_ten", "arith.mult.three_by_three",
	"arith.div.three_digit_divisor", "arith.div.average", "arith.div.missing_dividend",
	"arith.estimate.compatible", "arith.estimate.reasonableness",
	"arith.order_ops.exponents", "arith.neg.temperature", "arith.neg.elevation",
	"arith.order_ops.brackets", "arith.order_ops.negatives",
	"arith.word.unit_rate", "arith.word.comparison",
	"arith.dec.compare", "arith.dec.add",
	"arith.money.total", "arith.money.discount", "arith.money.unit_price",
	"arith.time.convert", "arith.time.duration", "arith.time.pace",
	"arith.factor.divisor_count",
}

func w12Fact(t *testing.T, f map[string]string, key string) int {
	t.Helper()
	raw, ok := f[key]
	if !ok {
		t.Fatalf("fact %q missing from %v", key, f)
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		t.Fatalf("fact %q is not an integer: %q", key, raw)
	}
	return v
}

func w12Recompute(t *testing.T, kind string, kp int, f map[string]string) int {
	t.Helper()
	switch kind {
	case "arith.add.three_addends":
		return w12Fact(t, f, "a") + w12Fact(t, f, "b") + w12Fact(t, f, "c")
	case "arith.add.four_digit":
		return w12Fact(t, f, "a") + w12Fact(t, f, "b")
	case "arith.add.regroup_twice":
		return w12Fact(t, f, "a") + w12Fact(t, f, "b")
	case "arith.sub.four_digit":
		return w12Fact(t, f, "a") - w12Fact(t, f, "b")
	case "arith.sub.regroup_twice":
		return w12Fact(t, f, "a") - w12Fact(t, f, "b")
	case "arith.sub.missing_subtrahend":
		return w12Fact(t, f, "a") - w12Fact(t, f, "c")
	case "arith.mult.three_by_two":
		return w12Fact(t, f, "a") * w12Fact(t, f, "b")
	case "arith.mult.powers_ten":
		p := 1
		for i := 0; i < w12Fact(t, f, "k"); i++ {
			p *= 10
		}
		return w12Fact(t, f, "n") * p
	case "arith.mult.three_by_three":
		return w12Fact(t, f, "a") * w12Fact(t, f, "b")
	case "arith.div.three_digit_divisor":
		return w12Fact(t, f, "dividend") / w12Fact(t, f, "divisor")
	case "arith.div.average":
		return (w12Fact(t, f, "n1") + w12Fact(t, f, "n2") + w12Fact(t, f, "n3") + w12Fact(t, f, "n4")) / 4
	case "arith.div.missing_dividend":
		return w12Fact(t, f, "divisor")*w12Fact(t, f, "quotient") + w12Fact(t, f, "remainder")
	case "arith.estimate.compatible":
		return w12Fact(t, f, "dividend") / w12Fact(t, f, "friendly")
	case "arith.estimate.reasonableness":
		unit := w12Fact(t, f, "unit")
		round := func(n int) int { return ((n + unit/2) / unit) * unit }
		return round(w12Fact(t, f, "a")) + round(w12Fact(t, f, "b"))
	case "arith.order_ops.exponents":
		a, b, c, d := w12Fact(t, f, "a"), w12Fact(t, f, "b"), w12Fact(t, f, "c"), w12Fact(t, f, "d")
		switch kp {
		case 1:
			p := 1
			for i := 0; i < c; i++ {
				p *= b
			}
			return a + p
		case 2:
			p := 1
			for i := 0; i < c; i++ {
				p *= b
			}
			return a + p*d
		default:
			base := a + b
			p := 1
			for i := 0; i < c; i++ {
				p *= base
			}
			return p - d
		}
	case "arith.order_ops.negatives":
		a, b, c, d := w12Fact(t, f, "a"), w12Fact(t, f, "b"), w12Fact(t, f, "c"), w12Fact(t, f, "d")
		switch kp {
		case 1:
			return -a + b*c - d
		case 2:
			return (a-b)*c + d
		default:
			return a - (b+c)*d
		}
	case "arith.neg.temperature":
		return w12Fact(t, f, "start") + w12Fact(t, f, "change")
	case "arith.neg.elevation":
		return w12Fact(t, f, "high") - w12Fact(t, f, "low")
	case "arith.order_ops.brackets":
		return (w12Fact(t, f, "a")+w12Fact(t, f, "b"))*w12Fact(t, f, "c") - w12Fact(t, f, "d")
	case "arith.word.unit_rate":
		return w12Fact(t, f, "cost") / w12Fact(t, f, "items")
	case "arith.word.comparison":
		return w12Fact(t, f, "first") - w12Fact(t, f, "second")
	case "arith.dec.compare":
		return w12Fact(t, f, "a") - w12Fact(t, f, "b")
	case "arith.dec.add":
		return w12Fact(t, f, "a") + w12Fact(t, f, "b")
	case "arith.money.total":
		return w12Fact(t, f, "p1") + w12Fact(t, f, "p2") + w12Fact(t, f, "p3")
	case "arith.money.discount":
		return w12Fact(t, f, "price") * w12Fact(t, f, "percent") / 100
	case "arith.money.unit_price":
		return w12Fact(t, f, "total") / w12Fact(t, f, "items")
	case "arith.time.convert":
		return w12Fact(t, f, "hours")*60 + w12Fact(t, f, "minutes")
	case "arith.time.duration":
		return w12Fact(t, f, "end") - w12Fact(t, f, "start")
	case "arith.time.pace":
		return w12Fact(t, f, "minutes") / w12Fact(t, f, "distance")
	case "arith.factor.divisor_count":
		n := w12Fact(t, f, "n")
		count := 0
		for i := 1; i*i <= n; i++ {
			if n%i == 0 {
				count++
				if i*i != n {
					count++
				}
			}
		}
		return count
	}
	t.Fatalf("no recompute rule for %s", kind)
	return 0
}

func TestWave12Answers(t *testing.T) {
	reg := generator.NewRegistry()
	registerWave12(reg)
	for _, id := range w12Kinds {
		gen, err := reg.Get(id)
		if err != nil {
			t.Fatalf("%s is not registered: %v", id, err)
		}
		for kp := 1; kp <= 3; kp++ {
			for i := 0; i < 40; i++ {
				p := gen.Generate(generator.GeneratorContext{KPOrdinal: kp})
				f := map[string]string{}
				for k, v := range p.Facts {
					f[k] = v
				}
				want := w12Recompute(t, id, kp, f)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer %q, recomputed %d (facts %v, question %q)", id, kp, p.Answer, want, p.Facts, p.Question)
				}
				if p.Facts["answer"] != p.Answer {
					t.Fatalf("%s kp%d: facts answer %q, answer %q", id, kp, p.Facts["answer"], p.Answer)
				}
				if !strings.HasSuffix(p.Question, "(enter a number)") {
					t.Fatalf("%s kp%d: question does not ask for a number: %q", id, kp, p.Question)
				}
				if p.KPOrdinal != kp {
					t.Fatalf("%s: KPOrdinal %d, want %d", id, p.KPOrdinal, kp)
				}
			}
		}
	}
}

func TestWave12Diversity(t *testing.T) {
	reg := generator.NewRegistry()
	registerWave12(reg)
	for _, id := range w12Kinds {
		gen, err := reg.Get(id)
		if err != nil {
			t.Fatalf("%s is not registered: %v", id, err)
		}
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 200; i++ {
				seen[gen.Generate(generator.GeneratorContext{KPOrdinal: kp}).Answer] = true
			}
			if len(seen) < 5 {
				t.Fatalf("%s kp%d: only %d distinct answers in 200 draws", id, kp, len(seen))
			}
		}
	}
}
