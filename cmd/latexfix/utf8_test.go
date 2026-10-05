package main

import (
	"strings"
	"testing"
)

// These are property guards for the two scanners in this file, not regression tests for a
// fixed bug.
//
// The byte-at-a-time form go-ai-lint AIL070 flags was **not** producing wrong answers, and it
// is worth being precise about why, because the reasoning is load-bearing for anyone tempted to
// "optimise" the rune decoding back out again:
//
//   - `unescapeDelimiters` wrote one byte per iteration via `strings.Builder.WriteByte`, which
//     appends the raw byte, and the loop visited every byte exactly once. That reconstructs the
//     input byte-for-byte. Verified below for em dashes, curly quotes, Greek letters and
//     ellipses.
//   - `textGroup` could not match an opener mid-rune, because a multi-byte UTF-8 sequence
//     contains only bytes ≥ 0x80 while every opener starts with a backslash at 0x5C.
//
// So the rewrite to rune-at-a-time is a robustness change, not a correctness one. These tests
// pin the properties that make it safe, so a future edit that drops or reorders a byte is
// caught even though nothing was broken today.

// Multi-byte content must pass through the delimiter rewrite untouched.
func TestUnescapeDelimitersPreservesMultibyteRunes(t *testing.T) {
	// Inputs are raw strings so the doubled backslashes are literal — that is the form the
	// Algebrica scrape emitted and the only form this function matches.
	cases := []struct {
		name string
		in   string
		want string
		n    int
	}{
		{"em dash passes through", "a—b", "a—b", 0},
		{"curly quotes pass through", "‘x’ and “y”", "‘x’ and “y”", 0},
		{"greek letters pass through", "α + β", "α + β", 0},
		{"ellipsis and thin space pass through", "wait… — now", "wait… — now", 0},
		{"rune after rewritten delimiters", `\\(x\\)—y`, `\(x\)` + "—y", 2},
		{"rune before rewritten delimiters", "α" + `\\(x\\)`, "α" + `\(x\)`, 2},
		{"multibyte inside rewritten delimiters", `\\(π\\)`, `\(π\)`, 2},
		// A row break has an optional length and unit inside the brackets, which is what
		// `rowBreakLen` matches on: `[2pt]` is a length and `[x]` is one too, so neither is
		// a display opener despite looking like one.
		{"row break keeps both backslashes (length)", `\\[2pt]`, `\\[2pt]`, 0},
		{"row break keeps both backslashes (unitless length)", `\\[x]`, `\\[x]`, 0},
		{"display opener is unescaped", `\\[\alpha]`, `\[\alpha]`, 1},
	}
	for _, c := range cases {
		got, n := unescapeDelimiters(c.in)
		if got != c.want {
			t.Errorf("%s: unescapeDelimiters(%q) = %q, want %q", c.name, c.in, got, c.want)
		}
		if n != c.n {
			t.Errorf("%s: count = %d, want %d", c.name, n, c.n)
		}
		if !utf8Valid(got) {
			t.Errorf("%s: result %q is not valid UTF-8", c.name, got)
		}
	}
}

// Whatever the internals, the function must be a pure byte-level identity outside the
// delimiters it rewrites. This is the property the byte-at-a-time version had and the reason
// no data was ever lost.
func TestUnescapeDelimitersIsByteIdentityOutsideDelimiters(t *testing.T) {
	for _, s := range []string{"—", "‘x’", "α", "…", "a—b—c", "∑ π ≈ 3.14"} {
		if got, _ := unescapeDelimiters(s); got != s {
			t.Errorf("unescapeDelimiters(%q) = %q, want it unchanged", s, got)
		}
	}
}

// Group ranges are byte offsets — the caller slices with them — and each must start and end on
// a character boundary. A range starting mid-rune is a perfectly valid slice, so nothing would
// panic; it would just silently begin one byte into a character.
func TestTextGroupRangesLandOnRuneBoundaries(t *testing.T) {
	content := "—prefix \\text{—inner—} —suffix \\mbox{π}"
	groups := textGroup(content)
	if len(groups) != 2 {
		t.Fatalf("got %d groups, want 2: %v", len(groups), groups)
	}
	for _, g := range groups {
		start, end := g[0], g[1]
		if start >= end {
			t.Errorf("group %v is empty or inverted", g)
			continue
		}
		if !isBoundary(content, start) {
			t.Errorf("group %v starts mid-rune at offset %d (byte 0x%02x)", g, start, content[start])
		}
		if end < len(content) && !isBoundary(content, end) {
			t.Errorf("group %v ends mid-rune at offset %d", g, end)
		}
		body := content[start:end]
		if !strings.Contains(body, "inner") && !strings.Contains(body, "π") {
			t.Errorf("group %v does not span its body: %q", g, body)
		}
	}
}

// The reason the byte-wise scan was safe, asserted rather than asserted-about: no opener can
// begin mid-rune because 0x5C is not a UTF-8 continuation byte. If this ever fails, a
// byte-at-a-time scan for `\\text{` could match inside a character and the rewrite above
// becomes load-bearing rather than incidental.
func TestBackslashCannotBeInsideAMultibyteRune(t *testing.T) {
	const backslash = 0x5C
	if backslash >= 0x80 && backslash <= 0xBF {
		t.Fatal("backslash is a valid continuation byte, so a byte-wise opener scan can match " +
			"mid-rune; the reasoning behind textGroup's byte offsets no longer holds")
	}
	// And confirm on a real sample: no continuation byte of a multi-byte rune equals 0x5C.
	const sample = "—‘α…"
	for i, r := range []byte(sample) {
		if i == 0 {
			continue // lead byte
		}
		if r == backslash {
			t.Errorf("continuation byte %d of %q is a backslash", i, sample)
		}
	}
}

func utf8Valid(s string) bool {
	for _, r := range s {
		if r == 0xFFFD {
			return false
		}
	}
	return true
}

func isBoundary(s string, i int) bool {
	if i == len(s) {
		return true
	}
	for j := range s {
		if j == i {
			return true
		}
	}
	return false
}
