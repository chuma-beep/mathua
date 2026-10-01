package fractions

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

func Register(reg *generator.Registry) {
	reg.Register("frac.basics.concept", &fracConceptGen{})
	reg.Register("frac.basics.parts", &fracPartsGen{})
	reg.Register("frac.basics.number_line", &fracNumberLineGen{})
	reg.Register("frac.basics.equivalent", &fracEquivalentGen{})
	reg.Register("frac.ops.simplify", &fracSimplifyGen{})
	reg.Register("frac.ops.compare", &fracCompareGen{})
	reg.Register("frac.ops.benchmark", &fracBenchmarkGen{})
	reg.Register("frac.ops.to_decimal", &fracToDecimalGen{})

	reg.Register("frac.add.same", &fracOpSameDenGen{op: "+"})
	reg.Register("frac.sub.same", &fracOpSameDenGen{op: "-"})
	reg.Register("frac.add.diff", &fracOpDiffDenGen{op: "+"})
	reg.Register("frac.sub.diff", &fracOpDiffDenGen{op: "-"})
	reg.Register("frac.add.word", &fracWordGen{op: "+"})
	reg.Register("frac.ops.mult", &fracMultGen{wholeMul: false})
	reg.Register("frac.mult.whole", &fracMultGen{wholeMul: true})
	reg.Register("frac.ops.div", &fracDivGen{wholeDiv: false})
	reg.Register("frac.div.whole", &fracDivGen{wholeDiv: true})
	reg.Register("frac.mixed.convert", &fracMixedConvertGen{})
	reg.Register("frac.mixed.add", &fracMixedOpGen{op: "+"})
	reg.Register("frac.mixed.sub", &fracMixedOpGen{op: "-"})
	reg.Register("frac.mixed.mult", &fracMixedMultGen{})
	reg.Register("frac.sub.word", &fracSubWordGen{})
	reg.Register("frac.mult.word", &fracMultWordGen{})
	reg.Register("frac.div.word", &fracDivWordGen{})
	reg.Register("frac.mixed.word", &fracMixedWordGen{})
	reg.Register("frac.compare.word", &fracCompareWordGen{})
}

func fracStr(num, den int) string { return fmt.Sprintf("%d/%d", num, den) }
func mixStr(whole, num, den int) string {
	if whole == 0 {
		return fracStr(num, den)
	}
	return fmt.Sprintf("%d %d/%d", whole, num, den)
}

func lcm(a, b int) int { return a * b / mathutil.GCD(a, b) }

// opIng is the word form of an operator, so authored prose can say "you are
// adding" without branching on the symbol in the schema.
func opIng(op string) string {
	if op == "-" {
		return "subtracting"
	}
	return "adding"
}

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

func (g *fracConceptGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	den := rand.Intn(max(1, scale*2)) + 3
	num := rand.Intn(max(1, den-2)) + 1
	var bar strings.Builder
	for i := 0; i < den; i++ {
		if i < num {
			bar.WriteString("X ")
		} else {
			bar.WriteString("_ ")
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("What fraction of the bar is filled?\n\n\\(%s\\)", bar.String()),
		Answer:      fracStr(num, den),
		Explanation: fmt.Sprintf("%d out of %d equal parts are filled = %d/%d", num, den, num, den),
		Facts:       map[string]string{"filled": strconv.Itoa(num), "parts": strconv.Itoa(den)},
	}
}

type fracPartsGen struct{}

func (g *fracPartsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	den := rand.Intn(max(1, scale*2)) + 2
	num := rand.Intn(max(1, den-2)) + 1
	pick := rand.Intn(2)
	if pick == 0 {
		return generator.Problem{
			Question:    fmt.Sprintf("In the fraction \\(\\frac{%d}{%d}\\), what is the denominator?", num, den),
			Answer:      strconv.Itoa(den),
			Explanation: fmt.Sprintf("The denominator %d tells how many equal parts make the whole.", den),
			Facts:       map[string]string{"num": strconv.Itoa(num), "den": strconv.Itoa(den)},
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("In the fraction \\(\\frac{%d}{%d}\\), what is the numerator?", num, den),
		Answer:      strconv.Itoa(num),
		Explanation: fmt.Sprintf("The numerator %d tells how many parts we have.", num),
		Facts:       map[string]string{"num": strconv.Itoa(num), "den": strconv.Itoa(den)},
	}
}

type fracNumberLineGen struct{}

func (g *fracNumberLineGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	den := rand.Intn(max(3, scale*2)) + 2
	num := rand.Intn(den-1) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Where is \\(\\frac{%d}{%d}\\) on a number line from \\(0\\) to \\(1\\)?", num, den),
		Answer:      fmt.Sprintf("%.2f", float64(num)/float64(den)),
		Explanation: fmt.Sprintf("%d/%d = %.2f, located between 0 and 1.", num, den, float64(num)/float64(den)),
		Facts:       map[string]string{"num": strconv.Itoa(num), "den": strconv.Itoa(den), "dec": fmt.Sprintf("%.2f", float64(num)/float64(den))},
	}
}

type fracEquivalentGen struct{}

func (g *fracEquivalentGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	den := rand.Intn(max(1, scale*2)) + 2
	num := rand.Intn(max(1, den-2)) + 1
	mult := rand.Intn(max(1, scale)) + 2
	goal := fracStr(num*mult, den*mult)
	return generator.Problem{
		Question:    fmt.Sprintf("Find an equivalent fraction to \\(\\frac{%d}{%d}\\) by multiplying numerator and denominator by \\(%d\\).", num, den, mult),
		Answer:      goal,
		Explanation: fmt.Sprintf("%d/%d x %d/%d = %d/%d", num, den, mult, mult, num*mult, den*mult),
		Facts:       map[string]string{"num": strconv.Itoa(num), "den": strconv.Itoa(den), "mult": strconv.Itoa(mult), "eq_num": strconv.Itoa(num * mult), "eq_den": strconv.Itoa(den * mult)},
	}
}

type fracSimplifyGen struct{}

func (g *fracSimplifyGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	den := rand.Intn(max(1, scale*3)) + 4
	factor := rand.Intn(max(1, scale)) + 2
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
		Question:    fmt.Sprintf("Simplify: \\(\\frac{%d}{%d}\\)", num, den),
		Answer:      fracStr(n, d),
		Explanation: fmt.Sprintf("Divide numerator and denominator by %d: %d/%d = %d/%d", mathutil.GCD(num, den), num, den, n, d),
		Facts:       map[string]string{"num": strconv.Itoa(num), "den": strconv.Itoa(den), "gcd": strconv.Itoa(mathutil.GCD(num, den)), "out_num": strconv.Itoa(n), "out_den": strconv.Itoa(d)},
	}
}

type fracCompareGen struct{}

func (g *fracCompareGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	aNum := rand.Intn(max(1, scale*2)) + 1
	aDen := rand.Intn(max(1, scale*2)) + 2
	bNum := rand.Intn(max(1, scale*2)) + 1
	bDen := rand.Intn(max(1, scale*2)) + 2
	for aNum*bDen == bNum*aDen {
		bNum = rand.Intn(max(1, scale*2)) + 1
	}
	av := float64(aNum) / float64(aDen)
	bv := float64(bNum) / float64(bDen)
	ans := ">"
	if bv > av {
		ans = "<"
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Compare: \\(\\frac{%d}{%d} \\_\\_ \\frac{%d}{%d}\\) (enter > or <)", aNum, aDen, bNum, bDen),
		Answer:      ans,
		Explanation: fmt.Sprintf("%d/%d = %.3f, %d/%d = %.3f, so %d/%d %s %d/%d", aNum, aDen, av, bNum, bDen, bv, aNum, aDen, ans, bNum, bDen),
		Facts:       map[string]string{"a_num": strconv.Itoa(aNum), "a_den": strconv.Itoa(aDen), "a_val": fmt.Sprintf("%.3f", av), "b_num": strconv.Itoa(bNum), "b_den": strconv.Itoa(bDen), "b_val": fmt.Sprintf("%.3f", bv), "op": ans},
	}
}

type fracBenchmarkGen struct{}

func (g *fracBenchmarkGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	den := rand.Intn(max(1, scale*2)) + 2
	num := rand.Intn(max(1, den-2)) + 1
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
		Question:    fmt.Sprintf("Is \\(\\frac{%d}{%d}\\) greater than, less than, or equal to \\(\\frac{1}{2}\\)?", num, den),
		Answer:      ans,
		Explanation: fmt.Sprintf("1/2 = %d/%d, and %d/%d %s 1/2", den/2, den, num, den, comp),
		Facts:       map[string]string{"num": strconv.Itoa(num), "den": strconv.Itoa(den), "dec": fmt.Sprintf("%.2f", val), "word": ans, "comp": comp},
	}
}

type fracToDecimalGen struct{}

func (g *fracToDecimalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	den := []int{2, 4, 5, 8, 10, 20, 25}[rand.Intn(7)]
	num := rand.Intn(max(1, scale*2)) + 1
	dec := float64(num) / float64(den)
	return generator.Problem{
		Question:    fmt.Sprintf("Convert \\(\\frac{%d}{%d}\\) to a decimal.", num, den),
		Answer:      fmt.Sprintf("%g", dec),
		Explanation: fmt.Sprintf("%d / %d = %g", num, den, dec),
		Facts:       map[string]string{"num": strconv.Itoa(num), "den": strconv.Itoa(den), "dec": fmt.Sprintf("%g", dec)},
	}
}

type fracOpSameDenGen struct{ op string }

func (g *fracOpSameDenGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	den := rand.Intn(max(1, scale*2)) + 3
	a := rand.Intn(max(1, den-2)) + 1
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
		Question:    fmt.Sprintf("\\(\\frac{%d}{%d} %s \\frac{%d}{%d} =\\) ?", a, den, g.op, b, den),
		Answer:      fmt.Sprintf("%d/%d", rn, rd),
		Explanation: fmt.Sprintf("(%d %s %d)/%d = %d/%d", a, g.op, b, den, rn, rd),
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "den": strconv.Itoa(den), "op": g.op, "op_ing": opIng(g.op), "raw_num": strconv.Itoa(result), "raw_den": strconv.Itoa(den), "out_num": strconv.Itoa(rn), "out_den": strconv.Itoa(rd)},
	}
}

type fracOpDiffDenGen struct{ op string }

func (g *fracOpDiffDenGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// Staged ladder: small dens → medium → large LCMs that need simplifying.
	denMax := 6
	if ctx.Difficulty >= 0.7 {
		denMax = 12
	} else if ctx.Difficulty >= 0.45 {
		denMax = 9
	}
	aDen := rand.Intn(denMax-2) + 2
	bDen := rand.Intn(denMax-2) + 2
	for aDen == bDen {
		bDen = rand.Intn(denMax-2) + 2
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
		Question:    fmt.Sprintf("\\(\\frac{%d}{%d} %s \\frac{%d}{%d} =\\) ?", aNum, aDen, g.op, bNum, bDen),
		Answer:      fmt.Sprintf("%d/%d", rn, rd),
		Explanation: fmt.Sprintf("Given %d/%d %s %d/%d. Step 1: common denominator LCM(%d,%d)=%d. Step 2: convert to %d/%d %s %d/%d. Step 3: combine = %d/%d. Answer: %d/%d", aNum, aDen, g.op, bNum, bDen, aDen, bDen, cm, aScaled, cm, g.op, bScaled, cm, rn, rd, rn, rd),
		Facts:       map[string]string{"a_num": strconv.Itoa(aNum), "a_den": strconv.Itoa(aDen), "b_num": strconv.Itoa(bNum), "b_den": strconv.Itoa(bDen), "op": g.op, "op_ing": opIng(g.op), "lcm": strconv.Itoa(cm), "a_scaled": strconv.Itoa(aScaled), "b_scaled": strconv.Itoa(bScaled), "raw_num": strconv.Itoa(result), "raw_den": strconv.Itoa(cm), "out_num": strconv.Itoa(rn), "out_den": strconv.Itoa(rd)},
	}
}

type fracWordGen struct{ op string }

func (g *fracWordGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// Staged ladder mirrors frac.add.diff dens so word problems scale too.
	denMax := 6
	if ctx.Difficulty >= 0.7 {
		denMax = 12
	} else if ctx.Difficulty >= 0.45 {
		denMax = 9
	}
	aDen := rand.Intn(denMax-2) + 2
	bDen := rand.Intn(denMax-2) + 2
	for aDen == bDen {
		bDen = rand.Intn(denMax-2) + 2
	}
	cm := lcm(aDen, bDen)
	aNum := rand.Intn(aDen-1) + 1
	bNum := rand.Intn(bDen-1) + 1
	aS := aNum * (cm / aDen)
	bS := bNum * (cm / bDen)
	result := aS + bS
	rn, rd := reduce(result, cm)
	return generator.Problem{
		Question:    fmt.Sprintf("You eat \\(\\frac{%d}{%d}\\) of a pizza. Your friend eats \\(\\frac{%d}{%d}\\). How much pizza was eaten total?", aNum, aDen, bNum, bDen),
		Answer:      fmt.Sprintf("%d/%d", rn, rd),
		Explanation: fmt.Sprintf("Given %d/%d + %d/%d. Step 1: common denominator LCM(%d,%d)=%d → %d/%d + %d/%d. Step 2: add tops = %d/%d. Answer: %d/%d of the pizza", aNum, aDen, bNum, bDen, aDen, bDen, cm, aS, cm, bS, cm, rn, rd, rn, rd),
		Facts:       map[string]string{"a_num": strconv.Itoa(aNum), "a_den": strconv.Itoa(aDen), "b_num": strconv.Itoa(bNum), "b_den": strconv.Itoa(bDen), "lcm": strconv.Itoa(cm), "a_scaled": strconv.Itoa(aS), "b_scaled": strconv.Itoa(bS), "raw_num": strconv.Itoa(result), "raw_den": strconv.Itoa(cm), "out_num": strconv.Itoa(rn), "out_den": strconv.Itoa(rd)},
	}
}

type fracMultGen struct{ wholeMul bool }

func (g *fracMultGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	if g.wholeMul {
		whole := rand.Intn(max(2, scale*2)) + 2
		den := rand.Intn(max(3, scale*2)) + 2
		num := rand.Intn(den-1) + 1
		rn, rd := reduce(num*whole, den)
		w, r, d := toMixed(rn, rd)
		if r == 0 {
			return generator.Problem{
				Question:    fmt.Sprintf("\\(%d \\times \\frac{%d}{%d} =\\) ?", whole, num, den),
				Answer:      strconv.Itoa(w),
				Explanation: fmt.Sprintf("%d x %d/%d = %d/%d = %d", whole, num, den, whole*num, den, w),
				Facts:       map[string]string{"whole": strconv.Itoa(whole), "num": strconv.Itoa(num), "den": strconv.Itoa(den), "raw_num": strconv.Itoa(whole * num), "raw_den": strconv.Itoa(den), "out_num": strconv.Itoa(rn), "out_den": strconv.Itoa(rd), "whole_part": strconv.Itoa(w)},
			}
		}
		return generator.Problem{
			Question:    fmt.Sprintf("\\(%d \\times \\frac{%d}{%d} =\\) ?", whole, num, den),
			Answer:      fmt.Sprintf("%d/%d", rn, rd),
			Explanation: fmt.Sprintf("%d x %d/%d = %d/%d = %s", whole, num, den, whole*num, den, mixStr(w, r, d)),
			Facts:       map[string]string{"whole": strconv.Itoa(whole), "num": strconv.Itoa(num), "den": strconv.Itoa(den), "raw_num": strconv.Itoa(whole * num), "raw_den": strconv.Itoa(den), "out_num": strconv.Itoa(rn), "out_den": strconv.Itoa(rd), "whole_part": strconv.Itoa(w), "rem": strconv.Itoa(r), "rem_den": strconv.Itoa(d)},
		}
	}
	aNum := rand.Intn(max(3, scale*2)) + 1
	aDen := rand.Intn(max(3, scale*2)) + 2
	bNum := rand.Intn(max(3, scale*2)) + 1
	bDen := rand.Intn(max(3, scale*2)) + 2
	rn, rd := reduce(aNum*bNum, aDen*bDen)
	return generator.Problem{
		Question:    fmt.Sprintf("\\(\\frac{%d}{%d} \\times \\frac{%d}{%d} =\\) ?", aNum, aDen, bNum, bDen),
		Answer:      fmt.Sprintf("%d/%d", rn, rd),
		Explanation: fmt.Sprintf("(%d x %d)/(%d x %d) = %d/%d", aNum, bNum, aDen, bDen, rn, rd),
		Facts:       map[string]string{"a_num": strconv.Itoa(aNum), "a_den": strconv.Itoa(aDen), "b_num": strconv.Itoa(bNum), "b_den": strconv.Itoa(bDen), "raw_num": strconv.Itoa(aNum * bNum), "raw_den": strconv.Itoa(aDen * bDen), "out_num": strconv.Itoa(rn), "out_den": strconv.Itoa(rd)},
	}
}

type fracDivGen struct{ wholeDiv bool }

func (g *fracDivGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	if g.wholeDiv {
		whole := rand.Intn(max(2, scale*2)) + 2
		aNum := rand.Intn(max(3, scale*2)) + 1
		aDen := rand.Intn(max(3, scale*2)) + 2
		rn, rd := reduce(aNum, aDen*whole)
		return generator.Problem{
			Question:    fmt.Sprintf("\\(\\frac{%d}{%d} \\div %d =\\) ?", aNum, aDen, whole),
			Answer:      fmt.Sprintf("%d/%d", rn, rd),
			Explanation: fmt.Sprintf("%d/%d / %d = %d/(%d x %d) = %d/%d", aNum, aDen, whole, aNum, aDen, whole, rn, rd),
			Facts:       map[string]string{"num": strconv.Itoa(aNum), "den": strconv.Itoa(aDen), "whole": strconv.Itoa(whole), "raw_num": strconv.Itoa(aNum), "raw_den": strconv.Itoa(aDen * whole), "out_num": strconv.Itoa(rn), "out_den": strconv.Itoa(rd)},
		}
	}
	aNum := rand.Intn(max(3, scale*2)) + 1
	aDen := rand.Intn(max(3, scale*2)) + 2
	bNum := rand.Intn(max(3, scale*2)) + 1
	bDen := rand.Intn(max(3, scale*2)) + 2
	rn, rd := reduce(aNum*bDen, aDen*bNum)
	return generator.Problem{
		Question:    fmt.Sprintf("\\(\\frac{%d}{%d} \\div \\frac{%d}{%d} =\\) ?", aNum, aDen, bNum, bDen),
		Answer:      fmt.Sprintf("%d/%d", rn, rd),
		Explanation: fmt.Sprintf("Flip and multiply: %d/%d x %d/%d = %d/%d", aNum, aDen, bDen, bNum, rn, rd),
		Facts:       map[string]string{"a_num": strconv.Itoa(aNum), "a_den": strconv.Itoa(aDen), "b_num": strconv.Itoa(bNum), "b_den": strconv.Itoa(bDen), "flipped_num": strconv.Itoa(bDen), "flipped_den": strconv.Itoa(bNum), "raw_num": strconv.Itoa(aNum * bDen), "raw_den": strconv.Itoa(aDen * bNum), "out_num": strconv.Itoa(rn), "out_den": strconv.Itoa(rd)},
	}
}

type fracMixedConvertGen struct{}

func (g *fracMixedConvertGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	whole := rand.Intn(max(1, scale*2)) + 1
	den := rand.Intn(max(1, scale*2)) + 2
	num := rand.Intn(den-1) + 1
	pick := rand.Intn(2)
	if pick == 0 {
		improper := whole*den + num
		return generator.Problem{
			Question:    fmt.Sprintf("Convert \\(%d \\frac{%d}{%d}\\) to an improper fraction.", whole, num, den),
			Answer:      fmt.Sprintf("%d/%d", improper, den),
			Explanation: fmt.Sprintf("Given %d %d/%d. Step 1: multiply whole part: %d × %d = %d. Step 2: add numerator: %d + %d = %d. Answer: %d/%d", whole, num, den, whole, den, whole*den, whole*den, num, improper, improper, den),
			Facts:       map[string]string{"whole": strconv.Itoa(whole), "num": strconv.Itoa(num), "den": strconv.Itoa(den), "scaled": strconv.Itoa(whole * den), "improper": strconv.Itoa(improper), "rem": strconv.Itoa(num), "rem_den": strconv.Itoa(den)},
		}
	}
	improper := whole*den + num
	w, r, d := toMixed(improper, den)
	return generator.Problem{
		Question:    fmt.Sprintf("Convert \\(\\frac{%d}{%d}\\) to a mixed number.", improper, den),
		Answer:      fmt.Sprintf("%d %d/%d", w, r, d),
		Explanation: fmt.Sprintf("Given %d/%d. Step 1: divide: %d ÷ %d = %d remainder %d. Step 2: whole part %d, fraction %d/%d. Answer: %d %d/%d", improper, den, improper, den, w, r, w, r, d, w, r, d),
		Facts:       map[string]string{"improper": strconv.Itoa(improper), "den": strconv.Itoa(den), "whole": strconv.Itoa(w), "rem": strconv.Itoa(r), "rem_den": strconv.Itoa(d), "num": strconv.Itoa(r), "scaled": strconv.Itoa(w * den)},
	}
}

type fracMixedOpGen struct{ op string }

func (g *fracMixedOpGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	w1 := rand.Intn(max(1, scale)) + 1
	w2 := rand.Intn(max(1, scale)) + 1
	den := rand.Intn(max(1, scale*2)) + 2
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
			Question:    fmt.Sprintf("\\(%d \\frac{%d}{%d} %s %d \\frac{%d}{%d} = ?\\)", w1, n1, den, g.op, w2, n2, den),
			Answer:      strconv.Itoa(w),
			Explanation: fmt.Sprintf("%d %d/%d %s %d %d/%d = %d", w1, n1, den, g.op, w2, n2, den, w),
			Facts:       map[string]string{"w1": strconv.Itoa(w1), "n1": strconv.Itoa(n1), "w2": strconv.Itoa(w2), "n2": strconv.Itoa(n2), "den": strconv.Itoa(den), "op": g.op, "op_ing": opIng(g.op), "a_improper": strconv.Itoa(aImproper), "b_improper": strconv.Itoa(bImproper), "result": strconv.Itoa(result), "whole": strconv.Itoa(w), "rem": "0", "rem_den": strconv.Itoa(den)},
		}
	}
	return generator.Problem{
		Question: fmt.Sprintf("\\(%d \\frac{%d}{%d} %s %d \\frac{%d}{%d} = ?\\)", w1, n1, den, g.op, w2, n2, den),
		Answer:   fmt.Sprintf("%d %d/%d", w, r, d),
		Explanation: fmt.Sprintf("Convert to improper: %d/%d %s %d/%d = %d/%d = %d %d/%d",
			aImproper, den, g.op, bImproper, den, result, den, w, r, d),
		Facts: map[string]string{"w1": strconv.Itoa(w1), "n1": strconv.Itoa(n1), "w2": strconv.Itoa(w2), "n2": strconv.Itoa(n2), "den": strconv.Itoa(den), "op": g.op, "op_ing": opIng(g.op), "a_improper": strconv.Itoa(aImproper), "b_improper": strconv.Itoa(bImproper), "result": strconv.Itoa(result), "whole": strconv.Itoa(w), "rem": strconv.Itoa(r), "rem_den": strconv.Itoa(d)},
	}
}

type fracMixedMultGen struct{}

func (g *fracMixedMultGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	w1 := rand.Intn(max(2, scale*2)) + 1
	w2 := rand.Intn(max(2, scale*2)) + 1
	den := rand.Intn(max(3, scale*2)) + 2
	n1 := rand.Intn(den-1) + 1
	n2 := rand.Intn(den-1) + 1
	aImp := w1*den + n1
	bImp := w2*den + n2
	num := aImp * bImp
	denSq := den * den
	rn, rd := reduce(num, denSq)
	w, r, d := toMixed(rn, rd)
	return generator.Problem{
		Question:    fmt.Sprintf("\\(%d \\frac{%d}{%d} \\times %d \\frac{%d}{%d} = ?\\)", w1, n1, den, w2, n2, den),
		Answer:      fmt.Sprintf("%d/%d", rn, rd),
		Explanation: fmt.Sprintf("(%d/%d) x (%d/%d) = %d/%d = %s", aImp, den, bImp, den, rn, rd, mixStr(w, r, d)),
		Facts:       map[string]string{"w1": strconv.Itoa(w1), "n1": strconv.Itoa(n1), "w2": strconv.Itoa(w2), "n2": strconv.Itoa(n2), "den": strconv.Itoa(den), "a_improper": strconv.Itoa(aImp), "b_improper": strconv.Itoa(bImp), "raw_num": strconv.Itoa(num), "raw_den": strconv.Itoa(denSq), "out_num": strconv.Itoa(rn), "out_den": strconv.Itoa(rd)},
	}
}

type fracSubWordGen struct{}

func (g *fracSubWordGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	aDen := rand.Intn(max(1, scale*2)) + 3
	bDen := rand.Intn(max(1, scale*2)) + 2
	for aDen == bDen {
		bDen = rand.Intn(max(1, scale*2)) + 2
	}
	cm := lcm(aDen, bDen)
	aNum := rand.Intn(aDen-1) + 1
	bNum := rand.Intn(bDen-1) + 1
	aS := aNum * (cm / aDen)
	bS := bNum * (cm / bDen)
	if aS < bS {
		aNum, bNum = bNum, aNum
		aDen, bDen = bDen, aDen
		aS, bS = bS, aS
	}
	result := aS - bS
	rn, rd := reduce(result, cm)
	return generator.Problem{
		Question:    fmt.Sprintf("You had \\(\\frac{%d}{%d}\\) of a cake and gave away \\(\\frac{%d}{%d}\\). How much is left?", aNum, aDen, bNum, bDen),
		Answer:      fmt.Sprintf("%d/%d", rn, rd),
		Explanation: fmt.Sprintf("%d/%d - %d/%d = %d/%d", aNum, aDen, bNum, bDen, rn, rd),
		Facts:       map[string]string{"a_num": strconv.Itoa(aNum), "a_den": strconv.Itoa(aDen), "b_num": strconv.Itoa(bNum), "b_den": strconv.Itoa(bDen), "lcm": strconv.Itoa(cm), "a_scaled": strconv.Itoa(aS), "b_scaled": strconv.Itoa(bS), "raw_num": strconv.Itoa(result), "raw_den": strconv.Itoa(cm), "out_num": strconv.Itoa(rn), "out_den": strconv.Itoa(rd)},
	}
}

type fracMultWordGen struct{}

func (g *fracMultWordGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	aNum := rand.Intn(max(3, scale*2)) + 1
	aDen := rand.Intn(max(3, scale*2)) + 2
	bNum := rand.Intn(max(3, scale*2)) + 1
	bDen := rand.Intn(max(3, scale*2)) + 2
	rn, rd := reduce(aNum*bNum, aDen*bDen)
	return generator.Problem{
		Question:    fmt.Sprintf("A recipe needs \\(\\frac{%d}{%d}\\) cup of sugar per serving. You make \\(\\frac{%d}{%d}\\) servings. How much sugar total?", aNum, aDen, bNum, bDen),
		Answer:      fmt.Sprintf("%d/%d", rn, rd),
		Explanation: fmt.Sprintf("%d/%d x %d/%d = %d/%d cups", aNum, aDen, bNum, bDen, rn, rd),
		Facts:       map[string]string{"a_num": strconv.Itoa(aNum), "a_den": strconv.Itoa(aDen), "b_num": strconv.Itoa(bNum), "b_den": strconv.Itoa(bDen), "raw_num": strconv.Itoa(aNum * bNum), "raw_den": strconv.Itoa(aDen * bDen), "out_num": strconv.Itoa(rn), "out_den": strconv.Itoa(rd)},
	}
}

type fracDivWordGen struct{}

func (g *fracDivWordGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	aNum := rand.Intn(max(3, scale*2)) + 1
	aDen := rand.Intn(max(3, scale*2)) + 2
	bNum := rand.Intn(max(3, scale*2)) + 1
	bDen := rand.Intn(max(3, scale*2)) + 2
	rn, rd := reduce(aNum*bDen, aDen*bNum)
	return generator.Problem{
		Question:    fmt.Sprintf("You have \\(\\frac{%d}{%d}\\) liters of juice. Each glass holds \\(\\frac{%d}{%d}\\) liters. How many glasses can you fill?", aNum, aDen, bNum, bDen),
		Answer:      fmt.Sprintf("%d/%d", rn, rd),
		Explanation: fmt.Sprintf("%d/%d ÷ %d/%d = %d/%d", aNum, aDen, bNum, bDen, rn, rd),
		Facts:       map[string]string{"a_num": strconv.Itoa(aNum), "a_den": strconv.Itoa(aDen), "b_num": strconv.Itoa(bNum), "b_den": strconv.Itoa(bDen), "flipped_num": strconv.Itoa(bDen), "flipped_den": strconv.Itoa(bNum), "raw_num": strconv.Itoa(aNum * bDen), "raw_den": strconv.Itoa(aDen * bNum), "out_num": strconv.Itoa(rn), "out_den": strconv.Itoa(rd)},
	}
}

type fracMixedWordGen struct{}

func (g *fracMixedWordGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	w1 := rand.Intn(max(1, scale)) + 1
	den := rand.Intn(max(1, scale*2)) + 3
	n1 := rand.Intn(den-1) + 1
	w2 := rand.Intn(max(1, scale)) + 1
	n2 := rand.Intn(den-1) + 1
	aImp := w1*den + n1
	bImp := w2*den + n2
	result := aImp + bImp
	w, r, d := toMixed(result, den)
	if r == 0 {
		return generator.Problem{
			Question:    fmt.Sprintf("You have \\(%d \\frac{%d}{%d}\\) meters of ribbon and buy \\(%d \\frac{%d}{%d}\\) more. Total?", w1, n1, den, w2, n2, den),
			Answer:      strconv.Itoa(w),
			Explanation: fmt.Sprintf("%d %d/%d + %d %d/%d = %d meters", w1, n1, den, w2, n2, den, w),
			Facts:       map[string]string{"w1": strconv.Itoa(w1), "n1": strconv.Itoa(n1), "w2": strconv.Itoa(w2), "n2": strconv.Itoa(n2), "den": strconv.Itoa(den), "whole": strconv.Itoa(w), "rem": "0", "rem_den": strconv.Itoa(den)},
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("You have \\(%d \\frac{%d}{%d}\\) meters of ribbon and buy \\(%d \\frac{%d}{%d}\\) more. Total?", w1, n1, den, w2, n2, den),
		Answer:      fmt.Sprintf("%d %d/%d", w, r, d),
		Explanation: fmt.Sprintf("%d %d/%d + %d %d/%d = %d %d/%d", w1, n1, den, w2, n2, den, w, r, d),
		Facts:       map[string]string{"w1": strconv.Itoa(w1), "n1": strconv.Itoa(n1), "w2": strconv.Itoa(w2), "n2": strconv.Itoa(n2), "den": strconv.Itoa(den), "whole": strconv.Itoa(w), "rem": strconv.Itoa(r), "rem_den": strconv.Itoa(d)},
	}
}

type fracCompareWordGen struct{}

func (g *fracCompareWordGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	aNum := rand.Intn(max(1, scale*2)) + 1
	aDen := rand.Intn(max(1, scale*2)) + 2
	bNum := rand.Intn(max(1, scale*2)) + 1
	bDen := rand.Intn(max(1, scale*2)) + 2
	for aNum*bDen == bNum*aDen {
		bNum = rand.Intn(max(1, scale*2)) + 1
	}
	// Word comparison: who has more?
	av2 := float64(aNum) / float64(aDen)
	bv2 := float64(bNum) / float64(bDen)
	ans := "first"
	if bv2 > av2 {
		ans = "second"
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Alice ate \\(\\frac{%d}{%d}\\) of a pizza, Bob ate \\(\\frac{%d}{%d}\\). Who ate more? (first/second)", aNum, aDen, bNum, bDen),
		Answer:      ans,
		Explanation: fmt.Sprintf("%d/%d = %.3f, %d/%d = %.3f, so %s ate more", aNum, aDen, av2, bNum, bDen, bv2, ans),
		Facts:       map[string]string{"a_num": strconv.Itoa(aNum), "a_den": strconv.Itoa(aDen), "a_val": fmt.Sprintf("%.3f", av2), "b_num": strconv.Itoa(bNum), "b_den": strconv.Itoa(bDen), "b_val": fmt.Sprintf("%.3f", bv2), "who": ans},
	}
}
