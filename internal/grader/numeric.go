package grader

import (
	"math"
	"math/big"
	"regexp"
	"strings"
)

var (
	commaRe      = regexp.MustCompile(`,`)
	plusZeroRe   = regexp.MustCompile(`^\+`)
	andRe        = regexp.MustCompile(`(?i)(\d)\s*and\s*(\d)`)
	mixedRe      = regexp.MustCompile(`^(-?\d+)\s+(\d+/\d+)$`)
	fractionRe   = regexp.MustCompile(`^(-?\d+)/(\d+)$`)
	leadingDotRe = regexp.MustCompile(`^(-?)\.`)
	sciRe        = regexp.MustCompile(`^(-?\d+\.?\d*)[eE](-?\d+)$`)
)

const floatTolerance = 1e-9

type numericGrader struct{}

func (g *numericGrader) grade(expected, answer string) Result {
	e := normalise(expected)
	a := normalise(answer)
	if e == "" || a == "" {
		return emptyResult()
	}
	eVal, eIsFloat, eOK := parseNumeric(e)
	aVal, aIsFloat, aOK := parseNumeric(a)
	if !eOK || !aOK {
		return Result{Correct: false, Score: 0, Feedback: "Could not parse numeric values"}
	}
	if eIsFloat || aIsFloat {
		if math.Abs(eVal-aVal) < floatTolerance {
			return Result{Correct: true, Score: 1}
		}
	} else {
		if eVal == aVal {
			return Result{Correct: true, Score: 1}
		}
	}
	return Result{Correct: false, Score: 0, Feedback: "Incorrect"}
}

func normalise(s string) string {
	s = commaRe.ReplaceAllString(s, "")
	s = plusZeroRe.ReplaceAllString(s, "")
	return strings.TrimSpace(s)
}

func parseNumeric(s string) (float64, bool, bool) {
	if s == "" {
		return 0, false, false
	}

	// Accept the spoken form of a mixed number: "4 and 1/10" / "2and3/10"
	// normalize to "4 1/10" before matching. (The English reading is how
	// students naturally type it, and matches the worked explanation.)
	s = andRe.ReplaceAllString(s, "$1 $2")

	// Mixed number: "1 1/2" — check before stripping whitespace. A negative
	// whole number subtracts its fraction ("-1 1/2" = -(1 + 1/2) = -1.5), not
	// adds it ("-1 1/2" must not parse as -1 + 0.5).
	if m := mixedRe.FindStringSubmatch(s); m != nil {
		whole := new(big.Rat)
		whole.SetString(m[1])
		frac := new(big.Rat)
		if _, ok := frac.SetString(m[2]); !ok {
			return 0, false, false
		}
		total := new(big.Rat)
		if whole.Sign() < 0 {
			total.Sub(whole, frac)
		} else {
			total.Add(whole, frac)
		}
		val, _ := total.Float64() // aislop-ignore-line ai-slop/swallowed-exception -- Float64's 2nd result is an exactness flag, not an error
		return val, true, true
	}

	// Strip whitespace for all other formats
	s = strings.Map(func(r rune) rune {
		if r == ' ' || r == '\t' {
			return -1
		}
		return r
	}, s)

	if strings.HasPrefix(s, "-.") {
		s = "-0." + s[2:]
	}
	if s == "." || strings.HasPrefix(s, ".") {
		s = "0" + s
	}

	// Scientific notation
	if m := sciRe.FindStringSubmatch(s); m != nil {
		coeff := new(big.Rat)
		if _, ok := coeff.SetString(m[1]); !ok {
			return 0, false, false
		}
		exp := new(big.Int)
		if _, ok := exp.SetString(m[2], 10); !ok {
			return 0, false, false
		}
		ten := big.NewInt(10)
		pow := new(big.Int).Exp(ten, new(big.Int).Abs(exp), nil)
		if exp.Sign() >= 0 {
			num := new(big.Int).Mul(coeff.Num(), pow)
			val, _ := new(big.Rat).SetFrac(num, coeff.Denom()).Float64() // aislop-ignore-line ai-slop/swallowed-exception -- Float64's 2nd result is an exactness flag, not an error
			return val, true, true
		}
		denom := new(big.Int).Mul(coeff.Denom(), pow)
		val, _ := new(big.Rat).SetFrac(coeff.Num(), denom).Float64() // aislop-ignore-line ai-slop/swallowed-exception -- Float64's 2nd result is an exactness flag, not an error
		return val, true, true
	}

	// Fraction
	if m := fractionRe.FindStringSubmatch(s); m != nil {
		frac := new(big.Rat)
		if _, ok := frac.SetString(m[1] + "/" + m[2]); !ok {
			return 0, false, false
		}
		val, _ := frac.Float64()
		return val, true, true
	}

	// Integer or float (including big.Rat for exact parsing)
	r := new(big.Rat)
	if _, ok := r.SetString(s); ok {
		val, _ := r.Float64()
		if r.IsInt() {
			return val, false, true
		}
		return val, true, true
	}

	return 0, false, false
}
