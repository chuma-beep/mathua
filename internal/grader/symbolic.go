package grader

import (
	"regexp"
	"strings"
)

// An exponent of digits only. A sign, a letter or an operator inside the
// parentheses keeps them, so x^(-1) and x^(n+1) keep their grouping.
var exponentParensRe = regexp.MustCompile(`\^\((\d+)\)`)

// normaliseSymbolic canonicalises an expression string for the exact-match
// fast path in gradeSymPy: whitespace-insensitive, `**` and `^` unified, and
// MathLive's digit-only exponent parentheses removed.
//
// MathLive parenthesises every exponent longer than one character, so 8x^{31}
// arrives as 8x^(31) from clients that do not run the frontend shim. That is
// the same expression, not a different one, and leaving it in means the answer
// misses the exact fast path. Digits only, so nothing is regrouped.
func normaliseSymbolic(s string) string {
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "**", "^")
	s = exponentParensRe.ReplaceAllString(s, "^$1")
	return strings.TrimSpace(s)
}
