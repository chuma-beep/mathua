package latex

import (
	"strings"
	"testing"
)

func TestCanonicalizeInlineDelims(t *testing.T) {
	in := `Euler said \(e^{i\pi} + 1 = 0\) and meant it.`
	want := `Euler said $e^{i\pi} + 1 = 0$ and meant it.`
	if got := Canonicalize(in, Authored); got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestCanonicalizeDisplayDelims(t *testing.T) {
	in := "Then:\n\\[ x^2 + y^2 = r^2 \\]\nDone."
	if got := Canonicalize(in, Authored); !strings.Contains(got, "$$\nx^2 + y^2 = r^2\n$$") {
		t.Errorf("display not canonicalized: %q", got)
	}
}

func TestCanonicalizeEquationEnv(t *testing.T) {
	in := "Consider\n\\begin{equation*}\n a_n = 5a_{n-1} - 6a_{n-2}\n\\end{equation*}\nas the recurrence."
	got := Canonicalize(in, Levin)
	if !strings.Contains(got, "$$a_n = 5a_{n-1} - 6a_{n-2}$$") {
		t.Errorf("equation env not unwrapped: %q", got)
	}
}

func TestCanonicalizeAlignEnv(t *testing.T) {
	in := "\\begin{align}\na &= b \\\\ c &= d\n\\end{align}"
	got := Canonicalize(in, Levin)
	if !strings.Contains(got, "\\begin{aligned}") {
		t.Errorf("align body not converted to aligned: %q", got)
	}
}

func TestPreprocessAmpAndEntities(t *testing.T) {
	in := "&amp; \\amp= 3 &gt; 2"
	got := Canonicalize(in, ORCCA)
	if strings.Contains(got, "\\amp") || strings.Contains(got, "&amp;") {
		t.Errorf("markers survived: %q", got)
	}
}

func TestValidateClean(t *testing.T) {
	if w := Validate("$\\frac{1}{2}$ and $$\\sqrt{3}$$", Structured); len(w) != 0 {
		t.Errorf("expected no warnings, got %v", w)
	}
}

func TestValidateUnbalancedBrace(t *testing.T) {
	w := Validate("$\\frac{1}{2$", Structured)
	found := false
	for _, m := range w {
		if strings.Contains(m, "unclosed '{'") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected unclosed-brace warning, got %v", w)
	}
}

func TestValidateUnpairedDollar(t *testing.T) {
	w := Validate("costs $5 and up", ORCCA)
	if len(w) == 0 {
		t.Error("expected odd-dollar warning for stray '$'")
	}
}


func TestValidateEscapedDollarClean(t *testing.T) {
	// The real arith.factor.find.md case: '\$' inside inline math is an
	// escaped literal dollar, not a delimiter.
	w := Validate("such as $\\$631{,}897.15$ , which is a terminating decimal", ORCCA)
	if len(w) != 0 {
		t.Errorf("escaped '\\$' should not warn, got %v", w)
	}
}

func TestValidateEscapedDollarInDisplay(t *testing.T) {
	w := Validate("$$a \\text{ costs } \\$5$$ tail", ORCCA)
	if len(w) != 0 {
		t.Errorf("escaped '\\$' inside display math should not warn, got %v", w)
	}
}

func TestValidateRowBreakThenRealDelimiter(t *testing.T) {
	// Even backslash run: '\\' row break followed by a REAL unpaired delimiter.
	// Must still warn — guards against a naive ReplaceAll("\\$", ...) fix.
	w := Validate("row one \\\\ $x tail", ORCCA)
	if len(w) == 0 {
		t.Error("real '$' after \\\\ row break must still be counted")
	}
}

func TestForSourceRouting(t *testing.T) {
	header := "> Content sourced from [OpenStax Prealgebra 1e](x) — CC BY 4.0\nbody"
	if a := ForSource("teaching/x.md", header); a.Name != "structured" {
		t.Errorf("openstax sniff failed: %s", a.Name)
	}
	if a := ForSource("teaching/geo.triangle.types.md", "# Types of Triangles\nhand-authored geometry lesson"); a.Name != "structured" {
		t.Errorf("authored-geometry sniff failed: %s", a.Name)
	}
	levinHeader := "> Content sourced from [Discrete Mathematics: An Open Introduction, 3e](y) by Oscar Levin"
	if a := ForSource("teaching/y.md", levinHeader); a.Name != "levin" {
		t.Errorf("levin sniff failed: %s", a.Name)
	}
	if a := ForSource("algebrica/functions/functions.md", ""); a.Name != "algebrica" {
		t.Errorf("algebrica route failed: %s", a.Name)
	}
	if a := ForSource("authored/z.md", ""); a.Name != "authored" {
		t.Errorf("authored route failed: %s", a.Name)
	}
}

func TestGeneratorsAdapterKeepsAnswersUntouched(t *testing.T) {
	in := `\(x^2\) means $x \cdot x$`
	got := Canonicalize(in, Generators)
	if !strings.HasPrefix(strings.TrimSpace(got), "$x^2$") {
		t.Errorf("generator inline conversion failed: %q", got)
	}
}

func TestAlgebricaDoubledDisplayDelims(t *testing.T) {
	in := `their integrals:\n\n\\[\int f(x) \\, dx = \\tag{1}\\]\n\ndone.`
	got := Canonicalize(in, Algebrica)
	if !strings.Contains(got, "$$") {
		t.Fatalf("doubled \\\\[] display delimiters not converted: %q", got)
	}
	if strings.Contains(got, "\\\\[") {
		t.Errorf("raw doubled delimiter survived: %q", got)
	}
}

func TestAlgebricaDoubledSpacing(t *testing.T) {
	in := `\\(\\int f(x) \\, dx\\)`
	got := Canonicalize(in, Algebrica)
	if !strings.Contains(got, "\\, dx") {
		t.Errorf("thin-space should collapse to single-backslash form: %q", got)
	}
	if strings.Contains(got, ", dx =") && !strings.Contains(in, ", dx") {
		t.Errorf("comma artifact introduced: %q", got)
	}
}

func TestAlgebricaDoubledSetBraces(t *testing.T) {
	in := `\\(\\mathbb{C} := \\{\\, z = a + bi \\mid a,\\, b \\in \\mathbb{R} \\,\\}\\)`
	got := Canonicalize(in, Algebrica)
	if !strings.Contains(got, `\{`) || !strings.Contains(got, `\}`) {
		t.Errorf("set braces not normalized to single backslash: %q", got)
	}
	if strings.Contains(got, `\\\\`) {
		t.Errorf("doubled braces survived: %q", got)
	}
}

func TestRepairBrokenSqrtRealCase(t *testing.T) {
	in := `\\( \\int ( 4x^3 - \\frac{3}{\\sqrt}{x}} + 2\\cos x ) \\, dx \\)`
	got := Canonicalize(in, Algebrica)
	if strings.Contains(got, `sqrt}{`) {
		t.Errorf("broken sqrt survived: %q", got)
	}
	if !strings.Contains(got, `\sqrt{x}`) {
		t.Errorf("expected repaired sqrt{x}: %q", got)
	}
}

func TestHeadingsPreservedIdempotently(t *testing.T) {
	in := "# Operations with complex numbers\n\n## Sum and difference"
	got := Canonicalize(in, Algebrica)
	if !strings.Contains(got, "# Operations with complex numbers") ||
		!strings.Contains(got, "## Sum and difference") {
		t.Errorf("headings were rewritten: %q", got)
	}
	// Canonicalize must be idempotent: heading levels never shift across calls.
	if again := Canonicalize(got, Algebrica); again != got {
		t.Errorf("not idempotent on headings:\nfirst:  %q\nsecond: %q", got, again)
	}
}

func TestStripSpanWrappersKeepsInner(t *testing.T) {
	in := `<span class="math-display">\int x\,dx</span>`
	got := Canonicalize(in, ORCCA)
	if strings.Contains(got, "<span") {
		t.Errorf("span wrapper survived: %q", got)
	}
	if !strings.Contains(got, `\int x`) {
		t.Errorf("inner content lost: %q", got)
	}
}

func TestAlgebricaRowbreakSpacingSurvives(t *testing.T) {
	in := "\\[\n\\begin{align}\n\\frac{3}{2}(2)^2 + c &= 1 \\\\\\\\\\[6pt]\nc &= -5\n\\end{align}\n\\]"
	got := Canonicalize(in, Algebrica)
	if strings.Contains(got, "%%MUARB") {
		t.Errorf("rowbreak token leaked into output: %q", got)
	}
	if !strings.Contains(got, "\\[6pt]") {
		t.Errorf("row-break spacing not restored: %q", got)
	}
	if strings.Contains(got, "$$\n$$") || strings.Contains(got, "\n$$\n$$") {
		t.Errorf("nested/empty display math produced: %q", got)
	}
}

func TestProtectRestoreRoundtrip(t *testing.T) {
	in := "x \\\\\\\\[6pt] y"
	protected := protectRowbreaks(in)
	if !strings.Contains(protected, "%%MUARB:") {
		t.Fatalf("rowbreak not tokenized: %q", protected)
	}
	// Any 2+ backslash variant normalizes to the canonical doubled form.
	out := restoreRowbreaks(protected)
	if want := "x \\\\[6pt] y"; out != want {
		t.Errorf("canonical form mismatch: got %q want %q", out, want)
	}
}

func TestAlgebricaEmRowbreakCasesEnv(t *testing.T) {
	in := "\\[\n\\\\log_af(x) = \\begin{cases}\na > 0 \\\\[0.6em]\na \\neq 1\n\\end{cases}\n\\]"
	got := Canonicalize(in, Algebrica)
	if strings.Contains(got, "%%MUARB") {
		t.Errorf("rowbreak token leaked: %q", got)
	}
	if !strings.Contains(got, "\\[0.6em]") {
		t.Errorf("em rowbreak not restored: %q", got)
	}
	if !strings.Contains(got, "\\begin{cases}") {
		t.Errorf("cases env lost: %q", got)
	}
	if strings.Contains(got, "\n$$") && strings.Count(got, "$$")%2 != 0 {
		t.Errorf("unpaired delimiters: %q", got)
	}
}

func TestCanonicalizeCurrencyPairs(t *testing.T) {
	// Scanner-validated pairs with pure-numeric content are prices, not math.
	in := "costs $5, saves $10 and pays $20 total"
	want := `costs \$5, saves \$10 and pays \$20 total`
	if got := Canonicalize(in, ORCCA); got != want {
		t.Errorf("got %q want %q", got, want)
	}
	// Real math sharing a line with escaped currency stays math.
	in2 := "If $x^2$ applies the cost is \\$3.99 each"
	got2 := Canonicalize(in2, ORCCA)
	if !strings.Contains(got2, "$x^2$") || !strings.Contains(got2, `\$3.99`) {
		t.Errorf("math/currency mix wrong: %q", got2)
	}
}
