package calculus

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/mathutil"
)

// Helpers for the parameterised derivative generators. Each conversion replaces
// a literal table of three or four rows with a draw whose answer is computed,
// so the variety floor for those concepts rises from the row count to the size
// of the parameter space.

// intPow returns base raised to a small non-negative integer exponent.
func intPow(base, exp int) int {
	out := 1
	for i := 0; i < exp; i++ {
		out *= base
	}
	return out
}

// coefStr renders a leading coefficient the way a learner writes it: 1 is
// silent and -1 is a single minus sign.
func coefStr(c int) string {
	switch c {
	case 1:
		return ""
	case -1:
		return "-"
	}
	return strconv.Itoa(c)
}

// varPow renders name^k, dropping the power at 0 and 1.
func varPow(name string, k int) string {
	switch k {
	case 0:
		return ""
	case 1:
		return name
	}
	return name + "^" + strconv.Itoa(k)
}

// texPow renders name^k for a LaTeX stem, where a first power is the bare name
// rather than name^{1}. Stems that print x^{1} read as machine output.
func texPow(name string, k int) string {
	switch k {
	case 0:
		return ""
	case 1:
		return name
	}
	return name + "^{" + strconv.Itoa(k) + "}"
}

// powTerm renders coef*x^exp as implicit multiplication, keeping the constant
// visible when the exponent is zero. varPow returns an empty string there, so
// coefStr plus varPow alone would render a term as nothing.
func powTerm(coef int, exp int) string {
	if exp == 0 {
		return strconv.Itoa(coef)
	}
	return coefStr(coef) + varPow("x", exp)
}

// monoTerm renders c*x^px*y^py as implicit-multiplication text. The grader
// parses implicit multiplication, so "6xy^2" is a valid symbolic answer.
func monoTerm(c, px, py int) string {
	if c == 0 {
		return "0"
	}
	term := coefStr(c) + varPow("x", px) + varPow("y", py)
	if term == "" || term == "-" {
		return term + "1"
	}
	return term
}

// decimal renders v rounded to the given number of decimal places, without
// trailing zeros, so 2 becomes "2" rather than "2.0000".
func decimal(v float64, places int) string {
	scale := math.Pow(10, float64(places))
	rounded := math.Round(v*scale) / scale
	return strconv.FormatFloat(rounded, 'f', -1, 64)
}

// powSuffix renders ^k for k above one and nothing otherwise, so a factor to
// the first power is not printed as ^{1}.
func powSuffix(k int) string {
	if k <= 1 {
		return ""
	}
	return "^" + strconv.Itoa(k)
}

// linearIn renders ax+by with a silent unit coefficient.
func linearIn(a, b int) string {
	return coefStr(a) + "x" + plusCoef(b) + "y"
}

// plusCoef renders +b for a positive b, with the coefficient omitted at 1.
func plusCoef(b int) string {
	if b == 1 {
		return "+"
	}
	return "+" + strconv.Itoa(b)
}

// quadForm renders a x^2 + b y^2.
func quadForm(a, b int) string {
	return monoTerm(a, 2, 0) + plusCoef(b) + "y^2"
}

// innerLinear renders the inner function ax+b of a chain-rule draw.
func innerLinear(a, b int) string {
	if a == 1 {
		return "x+" + strconv.Itoa(b)
	}
	return strconv.Itoa(a) + "x+" + strconv.Itoa(b)
}

// randRange draws a parameter at or above the difficulty-derived floor. Without
// the random part a generator that only called ScaleInt would return the same
// problem on every draw for a learner whose difficulty is fixed, which is the
// defect the parameterisation is meant to remove.
func randRange(difficulty float64, lo, hi int) int {
	if hi <= lo {
		return lo
	}
	floor := mathutil.ScaleInt(difficulty, lo, hi)
	if floor <= lo || rand.Intn(2) == 0 {
		return lo + rand.Intn(hi-lo+1)
	}
	return floor + rand.Intn(hi-floor+1)
}

// scaledXPow renders c*x^k with a silent unit coefficient, so a chain-rule stem
// reads 3x^2 rather than 1x^2.
func scaledXPow(c, k int) string {
	if c == 1 {
		return varPow("x", k)
	}
	return strconv.Itoa(c) + varPow("x", k)
}

// argX renders the argument of a trigonometric function: a multiplier of one is
// silent, so sin(x) does not read sin(1x).
func argX(a int) string {
	if a == 1 {
		return "x"
	}
	return fmt.Sprintf("%dx", a)
}

// signedTerm renders a term with its sign attached, for building a polynomial
// string where the next coefficient may be negative.
func signedTerm(c int, v string) string {
	switch {
	case c == 1:
		return "+" + v
	case c == -1:
		return "-" + v
	case c > 0:
		return fmt.Sprintf("+%d%s", c, v)
	default:
		return fmt.Sprintf("-%d%s", -c, v)
	}
}
