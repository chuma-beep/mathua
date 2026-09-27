package grader

import "strings"

// normaliseSymbolic canonicalises an expression string for the exact-match
// fast path in gradeSymPy: whitespace-insensitive, `**` and `^` unified.
func normaliseSymbolic(s string) string {
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "**", "^")
	return strings.TrimSpace(s)
}
