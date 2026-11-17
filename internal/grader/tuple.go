package grader

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var (
	tupleRe = regexp.MustCompile(`[();,]+`)
	// A degree symbol is hard to type from a math keyboard, and a learner who
	// writes 90 instead of 90° is not making a mathematical error. The same
	// applies to "deg" and "degrees", which are how the value is usually typed
	// when the symbol is not available.
	tupleDegreeRe = regexp.MustCompile(`(?i)\s*(?:°|deg(?:rees?)?)\s*$`)
)

type tupleGrader struct{}

func (g *tupleGrader) grade(expected, answer string) Result {
	e := normaliseTuple(expected)
	a := normaliseTuple(answer)
	if len(e) == 0 || len(a) == 0 {
		return emptyResult()
	}
	if len(e) != len(a) {
		return Result{Correct: false, Score: 0, Feedback: fmt.Sprintf("Expected %s, got %s", strings.Join(e, ","), strings.Join(a, ","))}
	}
	for i := range e {
		if !tuplePartEqual(e[i], a[i]) {
			return Result{Correct: false, Score: 0, Feedback: fmt.Sprintf("Mismatch at position %d", i+1)}
		}
	}
	return Result{Correct: true, Score: 1}
}

// tuplePartEqual ignores a degree marker, because it is a unit rather than
// part of the value, and compares everything else exactly. The exactness
// matters: a tuple such as (4000,300,5) lists three separate quantities, so
// accepting 4000.0 for 4000 would quietly grade a different answer as correct.
//
// Equations are the exception. A parametric line is the same line when the
// terms are reordered, so `x=5+3t` and `x=3t+5` are one answer written two
// ways. Each side keeps its label, then the right-hand sides are compared by
// coefficient, which accepts `3t+5`, `5+3t` and `5+3*t` without weakening the
// exact comparison for the numeric tuples above.
func tuplePartEqual(expected, answer string) bool {
	expected = strings.TrimSpace(tupleDegreeRe.ReplaceAllString(expected, ""))
	answer = strings.TrimSpace(tupleDegreeRe.ReplaceAllString(answer, ""))
	if strings.EqualFold(expected, answer) {
		return true
	}
	eLabel, eBody, eOK := splitEquation(expected)
	aLabel, aBody, aOK := splitEquation(answer)
	if !eOK || !aOK || !strings.EqualFold(eLabel, aLabel) {
		return false
	}
	eCoeffs, eOK := polynomialCoefficients(eBody)
	aCoeffs, aOK := polynomialCoefficients(aBody)
	if !eOK || !aOK || len(eCoeffs) != len(aCoeffs) {
		return false
	}
	for key, want := range eCoeffs {
		if math.Abs(want-aCoeffs[key]) > 1e-9 {
			return false
		}
	}
	return true
}

func splitEquation(s string) (label, body string, ok bool) {
	i := strings.IndexByte(s, '=')
	if i < 0 || strings.IndexByte(s[i+1:], '=') >= 0 {
		return "", "", false
	}
	label = strings.TrimSpace(s[:i])
	body = strings.TrimSpace(s[i+1:])
	if label == "" || body == "" {
		return "", "", false
	}
	return label, body, true
}

// polynomialCoefficients reads a sum of monomials such as `5+3t` or `-4+t`
// into a map keyed by variable and degree. It is deliberately small: the
// parametric generators emit integer coefficients and a single variable, and
// anything outside that shape falls back to exact comparison.
func polynomialCoefficients(s string) (map[string]float64, bool) {
	s = strings.ToLower(s)
	s = strings.NewReplacer(" ", "", "*", "", "{", "", "}", "", "\\", "", "\u2212", "-").Replace(s)
	if s == "" {
		return nil, false
	}
	terms := splitTerms(s)
	coeffs := make(map[string]float64, len(terms))
	for _, term := range terms {
		if term == "" {
			return nil, false
		}
		sign := 1.0
		switch term[0] {
		case '+':
			term = term[1:]
		case '-':
			sign = -1
			term = term[1:]
		}
		if term == "" {
			return nil, false
		}
		i := 0
		for i < len(term) && (term[i] >= '0' && term[i] <= '9' || term[i] == '.') {
			i++
		}
		coef := 1.0
		if i > 0 {
			parsed, err := strconv.ParseFloat(term[:i], 64)
			if err != nil {
				return nil, false
			}
			coef = parsed
		}
		rest := term[i:]
		key := "1"
		if rest != "" {
			j := 0
			for j < len(rest) && rest[j] >= 'a' && rest[j] <= 'z' {
				j++
			}
			name := rest[:j]
			if name == "" {
				return nil, false
			}
			degree := len(name)
			tail := rest[j:]
			if tail != "" {
				if len(name) != 1 || !strings.HasPrefix(tail, "^") {
					return nil, false
				}
				parsed, err := strconv.Atoi(tail[1:])
				if err != nil || parsed < 0 {
					return nil, false
				}
				degree = parsed
				name = name[:1]
			} else if len(name) > 1 && strings.Trim(name, name[:1]) != "" {
				return nil, false
			} else {
				name = name[:1]
			}
			key = fmt.Sprintf("%s^%d", name, degree)
		}
		coeffs[key] += sign * coef
	}
	for key, coef := range coeffs {
		if coef == 0 {
			delete(coeffs, key)
		}
	}
	return coeffs, true
}

func splitTerms(s string) []string {
	var terms []string
	start := 0
	prev := rune(0)
	for i, r := range s {
		if i > 0 && (r == '+' || r == '-') && prev != 'e' {
			terms = append(terms, s[start:i])
			start = i
		}
		prev = r
	}
	return append(terms, s[start:])
}

func normaliseTuple(s string) []string {
	s = strings.TrimSpace(s)
	s = tupleRe.ReplaceAllString(s, ",")
	parts := strings.Split(s, ",")
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
