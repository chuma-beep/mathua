package verify

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"
)

// Eval evaluates a mathematical expression and returns the result.
func Eval(expr string) (float64, error) {
	e, err := govaluate.NewEvaluableExpression(expr)
	if err != nil {
		return 0, fmt.Errorf("parse %q: %w", expr, err)
	}
	result, err := e.Evaluate(nil)
	if err != nil {
		return 0, fmt.Errorf("eval %q: %w", expr, err)
	}
	v, ok := result.(float64)
	if !ok {
		return 0, fmt.Errorf("unexpected result type %T for %q", result, expr)
	}
	return v, nil
}

// CheckExpr evaluates expr and verifies the result matches expected (within 1e-9 tolerance).
func CheckExpr(expr, expected string) error {
	got, err := Eval(expr)
	if err != nil {
		return err
	}
	e, err := strconv.ParseFloat(expected, 64)
	if err != nil {
		return fmt.Errorf("parse expected %q: %w", expected, err)
	}
	if math.Abs(got-e) > 1e-9 {
		return fmt.Errorf("%s = %g, want %g", expr, got, e)
	}
	return nil
}

// ExtractNumeric attempts to extract a numeric expression from a question string
// by matching common patterns. Returns the expression and true on success.
var patterns = []struct {
	re   *regexp.Regexp
	expr func([]string) string
}{
	// "a + b = ?"  or  "a + b = ?"  or  "a + b = ____"
	{regexp.MustCompile(`^(-?\d+)\s*\+\s*(-?\d+)\s*=\s*[?_\s]+$`),
		func(m []string) string { return m[1] + "+" + m[2] }},

	// "a - b = ?"
	{regexp.MustCompile(`^(-?\d+)\s*-\s*(-?\d+)\s*=\s*[?_\s]+$`),
		func(m []string) string { return m[1] + "-(" + m[2] + ")" }},

	// "a × b = ?"  or  "a x b = ?"  or  "a * b = ?"
	{regexp.MustCompile(`^(-?\d+)\s*[×x*]\s*(-?\d+)\s*=\s*[?_\s]+$`),
		func(m []string) string { return m[1] + "*" + m[2] }},

	// "count the Xs" (return the numeric token)
	{regexp.MustCompile(`^count\b`),
		func(m []string) string { return "" }},

	// "What is X + Y?"  or  "X + Y = ?"
	{regexp.MustCompile(`(?:^What is\s+)?(-?\d+)\s*\+\s*(-?\d+)(?:\s*=\s*\?)??$`),
		func(m []string) string { return m[1] + "+" + m[2] }},

	// "What is X% of Y?"  (percentage)
	{regexp.MustCompile(`(-?\d+)%\s*of\s*(-?\d+)`),
		func(m []string) string { return m[1] + "/100*" + m[2] }},
}

var unaryMinusRe = regexp.MustCompile(`([+*/(^-]\s*)-(\d+)`)

// ExtractExpr tries to extract a verifiable expression from a question string.
// Returns the normalised expression and true on success.
func ExtractExpr(question string) (string, bool) {
	q := strings.TrimSpace(question)
	for _, p := range patterns {
		if m := p.re.FindStringSubmatch(q); m != nil {
			e := p.expr(m)
			if e == "" {
				continue
			}
			// Normalise for govaluate: replace "*-3" with "*(0-3)" etc.
			e = unaryMinusRe.ReplaceAllString(e, "${1}(0-${2})")
			return e, true
		}
	}
	return "", false
}
