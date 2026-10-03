package main

import (
	"strings"
	"testing"

	"github.com/chuma-beep/latexnorm"
	"github.com/chuma-beep/mathua/internal/latex"
)

// The first version of this tool stripped the `$…$` off every span it touched,
// turning rendered math into literal text across 33 files. It read like content
// loss and it was invisible in the pass/fail count, because the spans still
// "existed" — they were just prose now. These tests exist to make that impossible
// to reintroduce.

func TestRepairPreservesDelimiters(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"inline math keeps its dollars", `$f✓(x)$`, `$f(x)$`},
		// Removing the tick leaves a bare numeral, which must stop being math
		// altogether or the canonicalizer escapes it to \$42\$ (ADR-016).
		{"bare numeral loses its delimiters", `$42✓$`, `42`},
		{"repaired expression keeps its math", `$2 + 7 = 9✓$`, `$2 + 7 = 9$`},
		{"display math keeps its double dollars", "$$\\frac{1}{2}✓$$", "$$\\frac{1}{2}$$"},
		{"paren math keeps its backslashes", `\(f’(x)\)`, `\(f'(x)\)`},
		{"bracket math keeps its brackets", `\[f’(x)\]`, `\[f'(x)\]`},
		{"two spans on one line both survive", `$a✓$ and $b✓$`, `$a$ and $b$`},
		{"untouched text is byte-identical", `no math here at all`, `no math here at all`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, _ := repair(c.in)
			if got != c.want {
				t.Errorf("repair(%q)\n got %q\nwant %q", c.in, got, c.want)
			}
		})
	}
}

// Prose inside math is prose. `\overset{\text{—}}{16}` means "16 repeats"; a blind
// dash rewrite would silently change what the lesson says.
func TestRepairLeavesTextGroupsAlone(t *testing.T) {
	in := `$\overset{\text{—}}{16} + f’(x)$`
	want := `$\overset{\text{—}}{16} + f'(x)$`
	got, _ := repair(in)
	if got != want {
		t.Errorf("repair(%q)\n got %q\nwant %q", in, got, want)
	}
}

func TestRepairLeavesNestedTextGroupsIntact(t *testing.T) {
	in := `$\text{a {b} ’c} + d’$`
	want := `$\text{a {b} ’c} + d'$`
	got, _ := repair(in)
	if got != want {
		t.Errorf("repair(%q)\n got %q\nwant %q", in, got, want)
	}
}

// The corpus is mostly English. A typographic apostrophe in prose is correct, and
// the whole reason this tool scans regions instead of doing find-and-replace.
func TestRepairDoesNotTouchProse(t *testing.T) {
	in := `John’s theorem states that \(x\) is nice. He said “yes” — twice.`
	got, changes := repair(in)
	if got != in {
		t.Errorf("prose was modified:\n got %q\nwant %q", got, in)
	}
	if len(changes) != 0 {
		t.Errorf("expected no changes in prose, got %d", len(changes))
	}
}

// Repair must not shift byte offsets, because the report quotes them.
func TestRepairIsStable(t *testing.T) {
	in := `Intro $a✓$ middle $b’$ end $c$`
	once, _ := repair(in)
	twice, changes := repair(once)
	if once != twice {
		t.Errorf("repair is not idempotent:\n once %q\ntwice %q", once, twice)
	}
	if len(changes) != 0 {
		t.Errorf("second pass still changed %d spans", len(changes))
	}
}

func TestTextGroupFindsNested(t *testing.T) {
	got := textGroup(`\text{a {b} c} x`)
	if len(got) != 1 {
		t.Fatalf("expected 1 group, got %d: %v", len(got), got)
	}
	if got[0][1] != len(`\text{a {b} c}`) {
		t.Errorf("group end = %d, want %d", got[0][1], len(`\text{a {b} c}`))
	}
}

// Every file must still scan into the same regions afterwards, or the repair
// silently restructured the document rather than editing it.
func TestRepairKeepsRegionCount(t *testing.T) {
	in := "Prose with $x$ and\n\n$$y$$\n\nmore \\(z\\) end $w✓$."
	out, _ := repair(in)
	before := len(latexnorm.Scan(in))
	after := len(latexnorm.Scan(out))
	if before != after {
		t.Errorf("region count changed: %d → %d", before, after)
	}
	if strings.Count(in, "$") != strings.Count(out, "$") {
		t.Errorf("dollar count changed: %d → %d", strings.Count(in, "$"), strings.Count(out, "$"))
	}
}

// The distinction that matters most: a derivative's apostrophe is a typo, a
// quotation mark is not. Rewriting the second to `'` renders primes and silently
// changes what the lesson teaches.
func TestRepairDistinguishesApostropheFromQuotation(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"derivative", `$f’(x) = 18x - 4$`, `$f'(x) = 18x - 4$`},
		{"subscripted derivative", `$D’(x_0)$`, `$D'(x_0)$`},
		{"after a closing brace", `$\sigma’(x)$`, `$\sigma'(x)$`},
		{"closing quotation mark is left alone", `{a < b ‘a\ \text{is less than}\ b’}`, `{a < b ‘a\ \text{is less than}\ b’}`},
		{"quote around a symbol still normalises to ASCII", `$‘x’$`, `$‘x'$`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, _ := repair(c.in)
			if got != c.want {
				t.Errorf("repair(%q)\n got %q\nwant %q", c.in, got, c.want)
			}
		})
	}
}

// The 14,936 doubled delimiters put every Algebrica inline formula in a *prose*
// region, so no region-based tool could see them. Un-escaping is the fix, and it
// is only safe because canonicalization erases the difference. That is the
// invariant, asserted directly rather than assumed.
func TestUnescapingDoesNotChangeWhatIsServed(t *testing.T) {
	cases := []string{
		`denoted by \\( f'(c) \\), is defined`,
		`\\[f'(c) = 5 \\]`,
		`mixed \\( a \\) and $b$ and $$c$$`,
		`no math at all, just prose with a backslash \\( misdirection`,
		`\begin{align} x &= 1 \\[2pt] y &= 2 \end{align}`,
	}
	for _, in := range cases {
		out, _ := repair(in)
		before := latex.Canonicalize(in, latex.ForSource("t.md", in))
		after := latex.Canonicalize(out, latex.ForSource("t.md", out))
		if before != after {
			t.Errorf("served output changed\n in: %q\nout: %q\nbefore: %q\n after: %q", in, out, before, after)
		}
	}
}

func TestUnescapeRewritesInlineDelimitersOnly(t *testing.T) {
	cases := []struct{ in, want string }{
		{`\\( x \\)`, `\( x \)`},
		{`\\[2pt] row break`, `\\[2pt] row break`}, // a row break, not a delimiter
		{`\\ alone`, `\\ alone`},
	}
	for _, c := range cases {
		got, _ := unescapeDelimiters(c.in)
		if got != c.want {
			t.Errorf("unescapeDelimiters(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// Un-escaping is what makes the math visible, so a doubled-delimiter formula must
// come out repaired rather than merely re-spelled.
func TestUnescapingExposesFormulasToTheFixers(t *testing.T) {
	in := `denoted by \\( f’(c) \\), is defined`
	want := `denoted by \( f'(c) \), is defined`
	got, changes := repair(in)
	if got != want {
		t.Errorf("repair(%q)\n got %q\nwant %q", in, got, want)
	}
	if len(changes) == 0 {
		t.Error("expected changes to be reported")
	}
}

func TestUnescapeDistinguishesDisplayFromRowBreak(t *testing.T) {
	cases := []struct{ in, want string }{
		{`\\[ x = 1 \\]`, `\[ x = 1 \]`},
		{`\begin{align} x &= 1 \\[2pt] y &= 2 \end{align}`, `\begin{align} x &= 1 \\[2pt] y &= 2 \end{align}`},
		{`\begin{align} a \\[0.5em] b \end{align}`, `\begin{align} a \\[0.5em] b \end{align}`},
		{`\begin{align} a \\[] b \end{align}`, `\begin{align} a \\[] b \end{align}`},
		// `\\[1em]` is a row break with 1em of space, not a display opener, even
		// though more content follows. The corpus has 1,157 of these.
		{`a \\[1em] b \\[2pt] c`, `a \\[1em] b \\[2pt] c`},
	}
	for _, c := range cases {
		got, _ := unescapeDelimiters(c.in)
		if got != c.want {
			t.Errorf("unescapeDelimiters(%q)\n got %q\nwant %q", c.in, got, c.want)
		}
	}
}
