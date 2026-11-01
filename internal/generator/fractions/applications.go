package fractions

// The seven concepts that finish the fractions domain: subtracting a fraction from a whole number,
// a fraction as a division, an ordering question against the half benchmark, a fraction of a unit of
// measure, scaling a recipe, the fraction remaining, and writing an equivalent fraction with a
// chosen denominator.
//
// These are the applications the domain's arithmetic was missing: the operations existed, but the
// situations that call for them were represented only as generic word problems.
//
// Answers follow the package's convention, a reduced fraction written as `num/den`.
import (
	"fmt"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerApplications(reg *generator.Registry) {
	reg.Register("frac.sub.whole", &subWholeGen{})
	reg.Register("frac.as_division", &asDivisionGen{})
	reg.Register("frac.order.half", &orderHalfGen{})
	reg.Register("frac.of_unit", &ofUnitGen{})
	reg.Register("frac.recipe.scale", &recipeScaleGen{})
	reg.Register("frac.remaining", &remainingGen{})
	reg.Register("frac.equivalent.scale", &equivalentScaleGen{})
}

// subWholeGen subtracts a fraction from a whole number, which needs the whole written as a fraction
// first.
type subWholeGen struct{}

func (g *subWholeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	w := smallInt(1, 3+int(ctx.Difficulty*4))
	b := smallInt(2, 5+int(ctx.Difficulty*5))
	a := smallInt(1, b-1)
	n, dd := reduce(w*b-a, b)
	return generator.Problem{
		Question: fmt.Sprintf("Subtract: \\(%d - \\frac{%d}{%d}\\). Give your answer as a fraction in lowest terms.", w, a, b),
		Answer:   fmt.Sprintf("%d/%d", n, dd),
		Explanation: fmt.Sprintf("Write the whole number over 1 and then over the common denominator: %d = %d/%d. Then %d/%d - %d/%d = %d/%d, which reduces to %d/%d.",
			w, w*b, b, w*b, b, a, b, w*b-a, b, n, dd),
		Facts: map[string]string{"out_num": strconv.Itoa(n), "out_den": strconv.Itoa(dd)},
	}
}

// metresWord writes "1 metre" rather than "1 metres".
func metresWord(n int) string {
	if n == 1 {
		return "1 metre"
	}
	return fmt.Sprintf("%d metres", n)
}

// asDivisionGen is a fraction read as a division: a length cut into equal pieces.
type asDivisionGen struct{}

func (g *asDivisionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := smallInt(1, 3+int(ctx.Difficulty*4))
	b := smallInt(2, 5+int(ctx.Difficulty*5))
	// A whole metre cut into an equal number of metres is the fraction 1/1 and
	// reads as nonsense, so the ribbon is never exactly as long as it has pieces.
	for b == a {
		b = smallInt(2, 5+int(ctx.Difficulty*5))
	}
	n, dd := reduce(a, b)
	return generator.Problem{
		Question: fmt.Sprintf("A ribbon %s long is cut into %d equal pieces. How long is each piece, as a fraction of a metre in lowest terms? (enter a fraction such as 3/2)", metresWord(a), b),
		Answer:   fmt.Sprintf("%d/%d", n, dd),
		Explanation: fmt.Sprintf("Sharing %d equally among %d is the division %d ÷ %d, and a division can be written as a fraction: %d/%d, which reduces to %d/%d. That is why the fraction bar is a division sign.",
			a, b, a, b, a, b, n, dd),
		Facts: map[string]string{"out_num": strconv.Itoa(n), "out_den": strconv.Itoa(dd),
			"length": strconv.Itoa(a), "pieces": strconv.Itoa(b)},
	}
}

// orderHalfGen compares several fractions to one half, which is the benchmark that needs no common
// denominator.
type orderHalfGen struct{}

func (g *orderHalfGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type frac struct{ n, d int }
	fs := make([]frac, 3)
	for i := range fs {
		fs[i] = frac{n: smallInt(1, 4+int(ctx.Difficulty*5)), d: smallInt(2, 5+int(ctx.Difficulty*5))}
	}
	count := 0
	for _, f := range fs {
		if 2*f.n > f.d {
			count++
		}
	}
	return generator.Problem{
		Question: fmt.Sprintf("How many of \\(\\frac{%d}{%d}\\), \\(\\frac{%d}{%d}\\) and \\(\\frac{%d}{%d}\\) are greater than one half?",
			fs[0].n, fs[0].d, fs[1].n, fs[1].d, fs[2].n, fs[2].d),
		Answer: fmt.Sprintf("%d/%d", count, 1),
		Explanation: fmt.Sprintf("Comparing to one half needs no common denominator: a fraction is greater than a half exactly when twice its numerator beats its denominator. Counting those gives %d of the three.",
			count),
		Facts: map[string]string{"count": strconv.Itoa(count), "out_num": strconv.Itoa(count), "out_den": "1"},
	}
}

// ofUnitGen asks what fraction of a unit a number of its parts is.
type ofUnitGen struct{}

func (g *ofUnitGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	units := []struct {
		article string
		name    string
		one     string
		many    string
		size    int
	}{
		{"an", "hour", "minute", "minutes", 60},
		{"a", "dollar", "cent", "cents", 100},
		{"a", "day", "hour", "hours", 24},
		{"a", "metre", "centimetre", "centimetres", 100},
	}
	u := units[smallInt(0, len(units)-1)]
	k := smallInt(1, u.size-1)
	n, dd := reduce(k, u.size)
	part := u.many
	if k == 1 {
		part = u.one
	}
	return generator.Problem{
		Question: fmt.Sprintf("What fraction of %s %s is %d %s? Give your answer in lowest terms.", u.article, u.name, k, part),
		Answer:   fmt.Sprintf("%d/%d", n, dd),
		Explanation: fmt.Sprintf("There are %d %s in %s %s, so %d %s is %d/%d, which reduces to %d/%d.",
			u.size, u.many, u.article, u.name, k, part, k, u.size, n, dd),
		Facts: map[string]string{"out_num": strconv.Itoa(n), "out_den": strconv.Itoa(dd),
			"parts": strconv.Itoa(k), "total": strconv.Itoa(u.size)},
	}
}

// recipeScaleGen scales a fractional quantity from one number of servings to another.
type recipeScaleGen struct{}

func (g *recipeScaleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := smallInt(1, 3+int(ctx.Difficulty*3))
	b := smallInt(2, 4+int(ctx.Difficulty*4))
	servings := smallInt(2, 4+int(ctx.Difficulty*3))
	want := servings * smallInt(1, 3)
	n, dd := reduce(a*want, b*servings)
	return generator.Problem{
		Question: fmt.Sprintf("A recipe needs \\(\\frac{%d}{%d}\\) of a cup of flour for %d servings. How much flour is needed for %d servings? Give your answer as a fraction of a cup in lowest terms.",
			a, b, servings, want),
		Answer: fmt.Sprintf("%d/%d", n, dd),
		Explanation: fmt.Sprintf("The amount per serving is %d/%d divided by %d, and multiplying by %d servings gives %d/%d, which reduces to %d/%d. Scaling a recipe is a multiplication by the ratio of the servings.",
			a, b, servings, want, a*want, b*servings, n, dd),
		Facts: map[string]string{"out_num": strconv.Itoa(n), "out_den": strconv.Itoa(dd)},
	}
}

// remainingGen asks what fraction is left after part of a whole is used.
type remainingGen struct{}

func (g *remainingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	b := smallInt(2, 5+int(ctx.Difficulty*5))
	a := smallInt(1, b-1)
	n, dd := reduce(b-a, b)
	return generator.Problem{
		Question: fmt.Sprintf("A tank is \\(\\frac{%d}{%d}\\) full. What fraction of the tank is empty? Give your answer in lowest terms.", a, b),
		Answer:   fmt.Sprintf("%d/%d", n, dd),
		Explanation: fmt.Sprintf("The whole tank is %d/%d, so the empty part is %d/%d - %d/%d = %d/%d, which reduces to %d/%d. Subtracting from one is how a remainder is written as a fraction.",
			b, b, b, b, a, b, b-a, b, n, dd),
		Facts: map[string]string{"out_num": strconv.Itoa(n), "out_den": strconv.Itoa(dd)},
	}
}

// equivalentScaleGen writes an equivalent fraction with a chosen denominator.
type equivalentScaleGen struct{}

func (g *equivalentScaleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := smallInt(1, 5)
	b := smallInt(2, 6)
	k := smallInt(2, 3+int(ctx.Difficulty*4))
	d := b * k
	n := a * k
	return generator.Problem{
		Question: fmt.Sprintf("Write \\(\\frac{%d}{%d}\\) as a fraction with denominator %d. What is its numerator?", a, b, d),
		Answer:   fmt.Sprintf("%d/%d", n, 1),
		Explanation: fmt.Sprintf("The denominator was multiplied by %d, so the numerator must be multiplied by the same number to keep the value: %d * %d = %d. Multiplying top and bottom by the same factor is what makes an equivalent fraction.",
			k, a, k, n),
		Facts: map[string]string{"numerator": strconv.Itoa(n), "out_num": strconv.Itoa(n),
			"out_den": "1", "factor": strconv.Itoa(k)},
	}
}
