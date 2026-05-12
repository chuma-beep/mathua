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
		return Result{Correct: false, Score: 0, Feedback: "Answer must not be empty"}
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

	// Mixed number: "1 1/2" — check before stripping whitespace
	if m := mixedRe.FindStringSubmatch(s); m != nil {
		whole := new(big.Rat)
		whole.SetString(m[1])
		frac := new(big.Rat)
		if _, ok := frac.SetString(m[2]); !ok {
			return 0, false, false
		}
		val, _ := new(big.Rat).Add(whole, frac).Float64()
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
			val, _ := new(big.Rat).SetFrac(num, coeff.Denom()).Float64()
			return val, true, true
		}
		denom := new(big.Int).Mul(coeff.Denom(), pow)
		val, _ := new(big.Rat).SetFrac(coeff.Num(), denom).Float64()
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
