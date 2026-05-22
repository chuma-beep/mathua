package grader

import (
	"regexp"
	"strings"
)

var (
	iToIRe     = regexp.MustCompile(`(\d)i\b`)
	standaloneI = regexp.MustCompile(`\bi\b`)
	polarRe    = regexp.MustCompile(`(\d+)\(cos\s+(\d+)°?\s*\+\s*i\s*sin\s+(\d+)°?\s*\)`)
	degreeRe   = regexp.MustCompile(`°`)
	plusMinusRe = regexp.MustCompile(`±`)
)

type complexGrader struct{}

func (g *complexGrader) grade(expected, answer string) Result {
	hasPlusMinus := strings.Contains(expected, "±") || strings.Contains(answer, "±")

	e := preprocessComplex(expected)
	a := preprocessComplex(answer)

	result := gradeSymPy(e, a)
	if result.Correct {
		return result
	}

	if hasPlusMinus {
		negE := negateComplexExpr(e)
		result2 := gradeSymPy(negE, a)
		if result2.Correct {
			return result2
		}
	}

	return result
}

func preprocessComplex(s string) string {
	s = strings.TrimSpace(s)
	s = degreeRe.ReplaceAllString(s, "")

	// Handle polar form:  r(cos θ + i sin θ) -> r*cos(θ*pi/180) + I*r*sin(θ*pi/180)
	s = polarRe.ReplaceAllStringFunc(s, func(match string) string {
		parts := polarRe.FindStringSubmatch(match)
		if len(parts) == 4 {
			return parts[1] + "*cos(" + parts[2] + "*pi/180) + I*" + parts[1] + "*sin(" + parts[3] + "*pi/180)"
		}
		return match
	})

	// Strip ± prefix for initial comparison (negated check is handled by caller)
	s = plusMinusRe.ReplaceAllString(s, "")

	// Convert 3i -> 3*I and standalone i -> I
	s = iToIRe.ReplaceAllString(s, "$1*I")
	s = standaloneI.ReplaceAllString(s, "I")

	// Normalize +I to +1*I for SymPy (though SymPy handles I fine)
	s = strings.ReplaceAll(s, "+I", "+1*I")
	s = strings.ReplaceAll(s, "-I", "-1*I")

	return s
}

func negateComplexExpr(s string) string {
	// Negate each term in the expression
	// This is a best-effort approach: wrap in -()
	return "-(" + s + ")"
}
