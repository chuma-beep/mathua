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
	if !strings.Contains(got, "\\begin{aligned}") || !strings.Contains(got, "$$") {
		t.Errorf("align env not converted to aligned-in-display: %q", got)
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
	if w := Validate("$\\frac{1}{2}$ and $$\\sqrt{3}$$"); len(w) != 0 {
		t.Errorf("expected no warnings, got %v", w)
	}
}

func TestValidateUnbalancedBrace(t *testing.T) {
	w := Validate("$\\frac{1}{2$")
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
	w := Validate("costs $5 and up")
	if len(w) == 0 {
		t.Error("expected odd-dollar warning for stray '$'")
	}
}

func TestForSourceRouting(t *testing.T) {
	header := "> Content sourced from [OpenStax Prealgebra 1e](x) — CC BY 4.0\nbody"
	if a := ForSource("teaching/x.md", header); a.Name != "openstax" {
		t.Errorf("openstax sniff failed: %s", a.Name)
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
