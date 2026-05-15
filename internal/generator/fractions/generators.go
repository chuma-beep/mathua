package fractions

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

func Register(reg *generator.Registry) {
	reg.Register("frac.concept", &fracConceptGen{})
	reg.Register("frac.parts", &fracPartsGen{})
	reg.Register("frac.on_number_line", &fracNumberLineGen{})
	reg.Register("frac.equivalent", &fracEquivalentGen{})
	reg.Register("frac.simplify", &fracSimplifyGen{})
	reg.Register("frac.compare", &fracCompareGen{})
	reg.Register("frac.benchmark", &fracBenchmarkGen{})
	reg.Register("frac.to_decimal", &fracToDecimalGen{})

	reg.Register("frac.add.same", &fracOpSameDenGen{op: "+"})
	reg.Register("frac.sub.same", &fracOpSameDenGen{op: "-"})
	reg.Register("frac.add.diff", &fracOpDiffDenGen{op: "+"})
	reg.Register("frac.sub.diff", &fracOpDiffDenGen{op: "-"})
	reg.Register("frac.add.word", &fracWordGen{op: "+"})
	reg.Register("frac.mult", &fracMultGen{wholeMul: false})
	reg.Register("frac.mult.whole", &fracMultGen{wholeMul: true})
	reg.Register("frac.div", &fracDivGen{wholeDiv: false})
	reg.Register("frac.div.whole", &fracDivGen{wholeDiv: true})
	reg.Register("frac.mixed.convert", &fracMixedConvertGen{})
	reg.Register("frac.mixed.add", &fracMixedOpGen{op: "+"})
	reg.Register("frac.mixed.sub", &fracMixedOpGen{op: "-"})
	reg.Register("frac.mixed.mult", &fracMixedMultGen{})
}

func fracStr(num, den int) string { return fmt.Sprintf("%d/%d", num, den) }
func mixStr(whole, num, den int) string {
	if whole == 0 {
		return fracStr(num, den)
	}
	return fmt.Sprintf("%d %d/%d", whole, num, den)
}

func lcm(a, b int) int { return a * b / mathutil.GCD(a, b) }

func reduce(num, den int) (int, int) {
	g := mathutil.GCD(num, den)
	return num / g, den / g
}

func toMixed(num, den int) (int, int, int) {
	whole := num / den
	rem := num % den
	return whole, rem, den
}

type fracConceptGen struct{}

func (g *fracConceptGen) Generate(difficulty float64) generator.Problem {
	den := rand.Intn(5) + 3
	num := rand.Intn(den-1) + 1
	bar := ""
	for i := 0; i < den; i++ {
		if i < num {
			bar += "X "
		} else {
			bar += "_ "
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("What fraction of the bar is filled?\n\n%s", bar),
		Answer:      fracStr(num, den),
		Explanation: fmt.Sprintf("%d out of %d equal parts are filled = %d/%d", num, den, num, den),
	}
}

type fracPartsGen struct{}

func (g *fracPartsGen) Generate(difficulty float64) generator.Problem {
	den := rand.Intn(8) + 2
	num := rand.Intn(den-1) + 1
	pick := rand.Intn(2)
	if pick == 0 {
		return generator.Problem{
			Question:    fmt.Sprintf("In the fraction %d/%d, what is the denominator?", num, den),
			Answer:      strconv.Itoa(den),
			Explanation: fmt.Sprintf("The denominator %d tells how many equal parts make the whole.", den),
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("In the fraction %d/%d, what is the numerator?", num, den),
		Answer:      strconv.Itoa(num),
		Explanation: fmt.Sprintf("The numerator %d tells how many parts we have.", num),
	}
}

type fracNumberLineGen struct{}

func (g *fracNumberLineGen) Generate(difficulty float64) generator.Problem {
	den := rand.Intn(4) + 2
	num := rand.Intn(den-1) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Where is %d/%d on a number line from 0 to 1?", num, den),
		Answer:      fmt.Sprintf("%.2f", float64(num)/float64(den)),
		Explanation: fmt.Sprintf("%d/%d = %.2f, located between 0 and 1.", num, den, float64(num)/float64(den)),
	}
}

type fracEquivalentGen struct{}

func (g *fracEquivalentGen) Generate(difficulty float64) generator.Problem {
	den := rand.Intn(5) + 2
	num := rand.Intn(den-1) + 1
	mult := rand.Intn(4) + 2
	goal := fracStr(num*mult, den*mult)
	return generator.Problem{
		Question:    fmt.Sprintf("Find an equivalent fraction to %d/%d by multiplying numerator and denominator by %d.", num, den, mult),
		Answer:      goal,
		Explanation: fmt.Sprintf("%d/%d x %d/%d = %d/%d", num, den, mult, mult, num*mult, den*mult),
	}
}

type fracSimplifyGen struct{}

func (g *fracSimplifyGen) Generate(difficulty float64) generator.Problem {
	den := rand.Intn(10) + 4
	factor := rand.Intn(4) + 2
	for den%factor != 0 {
		den++
	}
	maxNum := den/factor - 1
	if maxNum < 1 {
		maxNum = 1
	}
	num := (rand.Intn(maxNum) + 1) * factor
	n, d := reduce(num, den)
	return generator.Problem{
		Question:    fmt.Sprintf("Simplify: %d/%d", num, den),
		Answer:      fracStr(n, d),
		Explanation: fmt.Sprintf("Divide numerator and denominator by %d: %d/%d = %d/%d", mathutil.GCD(num, den), num, den, n, d),
	}
}

type fracCompareGen struct{}

func (g *fracCompareGen) Generate(difficulty float64) generator.Problem {
	aNum := rand.Intn(9) + 1
	aDen := rand.Intn(8) + 2
	bNum := rand.Intn(9) + 1
	bDen := rand.Intn(8) + 2
	for aNum*bDen == bNum*aDen {
		bNum = rand.Intn(9) + 1
	}
	av := float64(aNum) / float64(aDen)
	bv := float64(bNum) / float64(bDen)
	ans := ">"
	if bv > av {
		ans = "<"
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Compare: %d/%d __ %d/%d  (enter > or <)", aNum, aDen, bNum, bDen),
		Answer:      ans,
		Explanation: fmt.Sprintf("%d/%d = %.3f, %d/%d = %.3f, so %d/%d %s %d/%d", aNum, aDen, av, bNum, bDen, bv, aNum, aDen, ans, bNum, bDen),
	}
}

type fracBenchmarkGen struct{}

func (g *fracBenchmarkGen) Generate(difficulty float64) generator.Problem {
	den := rand.Intn(8) + 2
	num := rand.Intn(den-1) + 1
	val := float64(num) / float64(den)
	ans := "less"
	comp := "<"
	if val > 0.5 {
		ans = "greater"
		comp = ">"
	} else if val == 0.5 {
		ans = "equal"
		comp = "="
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Is %d/%d greater than, less than, or equal to 1/2?", num, den),
		Answer:      ans,
		Explanation: fmt.Sprintf("1/2 = %d/%d, and %d/%d %s 1/2", den/2, den, num, den, comp),
	}
}

type fracToDecimalGen struct{}

func (g *fracToDecimalGen) Generate(difficulty float64) generator.Problem {
	den := []int{2, 4, 5, 8, 10, 20, 25}[rand.Intn(7)]
	num := rand.Intn(den-1) + 1
	dec := float64(num) / float64(den)
	return generator.Problem{
		Question:    fmt.Sprintf("Convert %d/%d to a decimal.", num, den),
		Answer:      fmt.Sprintf("%g", dec),
		Explanation: fmt.Sprintf("%d / %d = %g", num, den, dec),
	}
}

type fracOpSameDenGen struct{ op string }

func (g *fracOpSameDenGen) Generate(difficulty float64) generator.Problem {
	den := rand.Intn(8) + 3
	a := rand.Intn(den-2) + 1
	var b int
	if g.op == "+" {
		b = rand.Intn(den-a) + 1
	} else {
		b = rand.Intn(a) + 1
	}
	var result int
	if g.op == "+" {
		result = a + b
	} else {
		result = a - b
	}
	rn, rd := reduce(result, den)
	return generator.Problem{
		Question:    fmt.Sprintf("%d/%d %s %d/%d = ?", a, den, g.op, b, den),
		Answer:      fmt.Sprintf("%d/%d", rn, rd),
		Explanation: fmt.Sprintf("(%d %s %d)/%d = %d/%d", a, g.op, b, den, rn, rd),
	}
}

type fracOpDiffDenGen struct{ op string }

func (g *fracOpDiffDenGen) Generate(difficulty float64) generator.Problem {
	aDen := rand.Intn(6) + 2
	bDen := rand.Intn(6) + 2
	for aDen == bDen {
		bDen = rand.Intn(6) + 2
	}
	cm := lcm(aDen, bDen)
	aNum := rand.Intn(aDen-1) + 1
	bNum := rand.Intn(bDen-1) + 1
	aScaled := aNum * (cm / aDen)
	bScaled := bNum * (cm / bDen)
	var result int
	if g.op == "+" {
		result = aScaled + bScaled
	} else {
		if aScaled < bScaled {
			aNum, bNum = bNum, aNum
			aDen, bDen = bDen, aDen
			aScaled, bScaled = bScaled, aScaled
		}
		result = aScaled - bScaled
	}
	rn, rd := reduce(result, cm)
	return generator.Problem{
		Question:    fmt.Sprintf("%d/%d %s %d/%d = ?", aNum, aDen, g.op, bNum, bDen),
		Answer:      fmt.Sprintf("%d/%d", rn, rd),
		Explanation: fmt.Sprintf("LCM(%d,%d)=%d: %d/%d %s %d/%d = %d/%d", aDen, bDen, cm, aScaled, cm, g.op, bScaled, cm, rn, rd),
	}
}

type fracWordGen struct{ op string }

func (g *fracWordGen) Generate(difficulty float64) generator.Problem {
	aDen := rand.Intn(6) + 2
	bDen := rand.Intn(6) + 2
	for aDen == bDen {
		bDen = rand.Intn(6) + 2
	}
	cm := lcm(aDen, bDen)
	aNum := rand.Intn(aDen-1) + 1
	bNum := rand.Intn(bDen-1) + 1
	aS := aNum * (cm / aDen)
	bS := bNum * (cm / bDen)
	result := aS + bS
	rn, rd := reduce(result, cm)
	return generator.Problem{
		Question:    fmt.Sprintf("You eat %d/%d of a pizza. Your friend eats %d/%d. How much pizza was eaten total?", aNum, aDen, bNum, bDen),
		Answer:      fmt.Sprintf("%d/%d", rn, rd),
		Explanation: fmt.Sprintf("%d/%d + %d/%d = %d/%d of the pizza.", aNum, aDen, bNum, bDen, rn, rd),
	}
}

type fracMultGen struct{ wholeMul bool }

func (g *fracMultGen) Generate(difficulty float64) generator.Problem {
	if g.wholeMul {
		whole := rand.Intn(9) + 2
		den := rand.Intn(7) + 2
		num := rand.Intn(den-1) + 1
		rn, rd := reduce(num*whole, den)
		w, r, d := toMixed(rn, rd)
		if r == 0 {
			return generator.Problem{
				Question:    fmt.Sprintf("%d x %d/%d = ?", whole, num, den),
				Answer:      strconv.Itoa(w),
				Explanation: fmt.Sprintf("%d x %d/%d = %d/%d = %d", whole, num, den, whole*num, den, w),
			}
		}
		return generator.Problem{
			Question:    fmt.Sprintf("%d x %d/%d = ?", whole, num, den),
			Answer:      fmt.Sprintf("%d/%d", rn, rd),
			Explanation: fmt.Sprintf("%d x %d/%d = %d/%d = %s", whole, num, den, whole*num, den, mixStr(w, r, d)),
		}
	}
	aNum := rand.Intn(6) + 1
	aDen := rand.Intn(8) + 2
	bNum := rand.Intn(6) + 1
	bDen := rand.Intn(8) + 2
	rn, rd := reduce(aNum*bNum, aDen*bDen)
	return generator.Problem{
		Question:    fmt.Sprintf("%d/%d x %d/%d = ?", aNum, aDen, bNum, bDen),
		Answer:      fmt.Sprintf("%d/%d", rn, rd),
		Explanation: fmt.Sprintf("(%d x %d)/(%d x %d) = %d/%d", aNum, bNum, aDen, bDen, rn, rd),
	}
}

type fracDivGen struct{ wholeDiv bool }

func (g *fracDivGen) Generate(difficulty float64) generator.Problem {
	if g.wholeDiv {
		whole := rand.Intn(9) + 2
		aNum := rand.Intn(6) + 1
		aDen := rand.Intn(8) + 2
		rn, rd := reduce(aNum, aDen*whole)
		return generator.Problem{
			Question:    fmt.Sprintf("%d/%d / %d = ?", aNum, aDen, whole),
			Answer:      fmt.Sprintf("%d/%d", rn, rd),
			Explanation: fmt.Sprintf("%d/%d / %d = %d/(%d x %d) = %d/%d", aNum, aDen, whole, aNum, aDen, whole, rn, rd),
		}
	}
	aNum := rand.Intn(6) + 1
	aDen := rand.Intn(8) + 2
	bNum := rand.Intn(6) + 1
	bDen := rand.Intn(8) + 2
	rn, rd := reduce(aNum*bDen, aDen*bNum)
	return generator.Problem{
		Question:    fmt.Sprintf("%d/%d / %d/%d = ?", aNum, aDen, bNum, bDen),
		Answer:      fmt.Sprintf("%d/%d", rn, rd),
		Explanation: fmt.Sprintf("Flip and multiply: %d/%d x %d/%d = %d/%d", aNum, aDen, bDen, bNum, rn, rd),
	}
}

type fracMixedConvertGen struct{}

func (g *fracMixedConvertGen) Generate(difficulty float64) generator.Problem {
	whole := rand.Intn(6) + 1
	den := rand.Intn(7) + 2
	num := rand.Intn(den-1) + 1
	pick := rand.Intn(2)
	if pick == 0 {
		improper := whole*den + num
		return generator.Problem{
			Question:    fmt.Sprintf("Convert %d %d/%d to an improper fraction.", whole, num, den),
			Answer:      fmt.Sprintf("%d/%d", improper, den),
			Explanation: fmt.Sprintf("%d x %d + %d = %d, so %d %d/%d = %d/%d", whole, den, num, improper, whole, num, den, improper, den),
		}
	}
	improper := whole*den + num
	w, r, d := toMixed(improper, den)
	return generator.Problem{
		Question:    fmt.Sprintf("Convert %d/%d to a mixed number.", improper, den),
		Answer:      fmt.Sprintf("%d %d/%d", w, r, d),
		Explanation: fmt.Sprintf("%d / %d = %d remainder %d, so %d/%d = %d %d/%d", improper, den, w, r, improper, den, w, r, d),
	}
}

type fracMixedOpGen struct{ op string }

func (g *fracMixedOpGen) Generate(difficulty float64) generator.Problem {
	w1 := rand.Intn(4) + 1
	w2 := rand.Intn(4) + 1
	den := rand.Intn(6) + 2
	n1 := rand.Intn(den-1) + 1
	n2 := rand.Intn(den-1) + 1
	aImproper := w1*den + n1
	bImproper := w2*den + n2
	var result int
	if g.op == "+" {
		result = aImproper + bImproper
	} else {
		if aImproper < bImproper {
			n1, n2 = n2, n1
			w1, w2 = w2, w1
			aImproper, bImproper = bImproper, aImproper
		}
		result = aImproper - bImproper
	}
	w, r, d := toMixed(result, den)
	if r == 0 {
		return generator.Problem{
			Question:    fmt.Sprintf("%d %d/%d %s %d %d/%d = ?", w1, n1, den, g.op, w2, n2, den),
			Answer:      strconv.Itoa(w),
			Explanation: fmt.Sprintf("%d %d/%d %s %d %d/%d = %d", w1, n1, den, g.op, w2, n2, den, w),
		}
	}
	return generator.Problem{
		Question: fmt.Sprintf("%d %d/%d %s %d %d/%d = ?", w1, n1, den, g.op, w2, n2, den),
		Answer:   fmt.Sprintf("%d %d/%d", w, r, d),
		Explanation: fmt.Sprintf("Convert to improper: %d/%d %s %d/%d = %d/%d = %d %d/%d",
			aImproper, den, g.op, bImproper, den, result, den, w, r, d),
	}
}

type fracMixedMultGen struct{}

func (g *fracMixedMultGen) Generate(difficulty float64) generator.Problem {
	w1 := rand.Intn(3) + 1
	w2 := rand.Intn(3) + 1
	den := rand.Intn(5) + 2
	n1 := rand.Intn(den-1) + 1
	n2 := rand.Intn(den-1) + 1
	aImp := w1*den + n1
	bImp := w2*den + n2
	num := aImp * bImp
	denSq := den * den
	rn, rd := reduce(num, denSq)
	w, r, d := toMixed(rn, rd)
	return generator.Problem{
		Question:    fmt.Sprintf("%d %d/%d x %d %d/%d = ?", w1, n1, den, w2, n2, den),
		Answer:      fmt.Sprintf("%d/%d", rn, rd),
		Explanation: fmt.Sprintf("(%d/%d) x (%d/%d) = %d/%d = %s", aImp, den, bImp, den, rn, rd, mixStr(w, r, d)),
	}
}
