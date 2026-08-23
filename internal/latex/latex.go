// Package latex normalizes mathematics notation from heterogeneous content
// sources into a single KaTeX-friendly dialect ($ inline, $$ display) at
// ingestion time, so neither storage nor the frontend ever sees a source's
// quirks. Validation reports structural problems (unbalanced braces or
// delimiters) as warnings attached to the canonicalized output; callers
// decide whether to reject, quarantine, or serve with logging.
package latex

import (
	"fmt"
	"regexp"
	"strings"
)

// Adapter absorbs surface-syntax quirks specific to one content source.
type Adapter struct {
	Name       string
	Preprocess func(string) string
}

var entityReplacer = strings.NewReplacer(
	"&amp;", "&",
	"&lt;", "<",
	"&gt;", ">",
	"&quot;", "\"",
	"&#39;", "'",
	"&nbsp;", " ",
)

func preprocessGeneric(s string) string {
	s = entityReplacer.Replace(s)
	s = strings.ReplaceAll(s, `\amp`, "&")
	return s
}

var (
	equationEnvRe  = regexp.MustCompile(`(?s)\\begin\{(equation\*?)\}(.*?)\\end\{equation\*?\}`)
	alignEnvRe     = regexp.MustCompile(`(?s)\\begin\{(align\*?)\}(.*?)\\end\{align\*?\}`)
	gatheredEnvRe  = regexp.MustCompile(`(?s)\\begin\{gathered\}(.*?)\\end\{gathered\}`)
	displayDelimRe = regexp.MustCompile(`(?s)\\\[([\s\S]*?)\\\]`)
	inlineDelimRe  = regexp.MustCompile(`\\\(([\s\S]*?)\\\\?\)|\\\(([\s\S]*?)\)`)
	mathTagRe      = regexp.MustCompile(`(?s)<math[^>]*>(.*?)</math>`)
)

// Canonicalize rewrites every recognized math construct to exactly two
// delimiters: inline $...$ and display $$...$$ (KaTeX-native, remark-math
// friendly). Source-specific quirks are absorbed first via the adapter.
func Canonicalize(s string, a Adapter) string {
	if a.Preprocess != nil {
		s = a.Preprocess(s)
	}

	// Raw <math> HTML wrappers (legacy wiki content): unwrap, then let the
	// standard pipeline treat the inner LaTeX like any other.
	for {
		m := mathTagRe.FindStringSubmatchIndex(s)
		if m == nil {
			break
		}
		inner := s[m[2]:m[3]]
		s = s[:m[0]] + "$$" + inner + "$$" + s[m[1]:]
	}

	// Display environments → $$ ... $$
	s = equationEnvRe.ReplaceAllStringFunc(s, func(m string) string {
		g := equationEnvRe.FindStringSubmatch(m)
		return "$$" + strings.TrimSpace(g[2]) + "$$"
	})
	s = alignEnvRe.ReplaceAllStringFunc(s, func(m string) string {
		g := alignEnvRe.FindStringSubmatch(m)
		body := strings.TrimSpace(g[2])
		body = strings.Replace(body, `\begin{aligned}`, "", 1)
		body = strings.Replace(body, `\end{aligned}`, "", 1)
		return "$$\n\\begin{aligned}\n" + strings.TrimSpace(body) + "\n\\end{aligned}\n$$"
	})

	// \[ ...\ ] → $$ ... $$
	s = displayDelimRe.ReplaceAllStringFunc(s, func(m string) string {
		g := displayDelimRe.FindStringSubmatch(m)
		return "$$\n" + strings.TrimSpace(g[1]) + "\n$$"
	})

	// \( ... \) → $ ... $
	s = inlineDelimRe.ReplaceAllStringFunc(s, func(m string) string {
		g := inlineDelimRe.FindStringSubmatch(m)
		inner := g[1]
		if inner == "" {
			inner = g[2]
		}
		return "$" + inner + "$"
	})

	return s
}

// Validate checks canonicalized content for structural problems: unbalanced
// braces inside math regions and stray/unpaired delimiters. It returns
// human-readable warnings; empty slice means clean.
func Validate(canonical string) []string {
	var warnings []string

	for _, m := range regexp.MustCompile(`(?s)\$\$(.*?)\$\$`).FindAllStringSubmatch(canonical, -1) {
		if msg := checkBraces(m[1]); msg != "" {
			warnings = append(warnings, "display math: "+msg)
		}
	}
	outside := regexp.MustCompile(`(?s)\$\$.*?\$\$`).ReplaceAllString(canonical, "")
	for _, m := range regexp.MustCompile(`(?s)\$(.+?)\$`).FindAllStringSubmatch(outside, -1) {
		if msg := checkBraces(m[1]); msg != "" {
			warnings = append(warnings, "inline math: "+msg)
		}
	}

	dollarCount := strings.Count(regexp.MustCompile(`(?s)\$\$.*?\$\$`).ReplaceAllString(canonical, ""), "$")
	if dollarCount%2 != 0 {
		warnings = append(warnings, fmt.Sprintf("odd number of inline delimiters (%d) — possible unpaired '$'", dollarCount))
	}
	return warnings
}

func checkBraces(segment string) string {
	depth := 0
	inEscape := false
	for _, r := range segment {
		switch {
		case inEscape:
			inEscape = false
		case r == '\\':
			inEscape = true
		case r == '{':
			depth++
		case r == '}':
			depth--
			if depth < 0 {
				return "unbalanced '}' (extra closing brace)"
			}
		}
	}
	if depth > 0 {
		return fmt.Sprintf("%d unclosed '{'", depth)
	}
	if inEscape {
		return "dangling backslash"
	}
	return ""
}

// Adapters for the known content sources.
var (
	// Algebrica: \( \) inline, occasional \left/\right sizing — no special
	// preprocessing beyond the shared entity/\amp pass today.
	Algebrica = Adapter{Name: "algebrica", Preprocess: preprocessGeneric}

	// Authored: hand-written lessons already use $ / $$ conventions.
	Authored = Adapter{Name: "authored", Preprocess: preprocessGeneric}

	// Generators: live-generated problem text.
	Generators = Adapter{Name: "generators", Preprocess: preprocessGeneric}

	// Levin extracts: PreTeXt HTML → markdown keeps \begin{equation*} blocks,
	// \amp alignment markers, and full exercise sections.
	Levin = Adapter{Name: "levin", Preprocess: func(s string) string {
		s = preprocessGeneric(s)
		s = stripSections(s, "Exercises")
		return s
	}}

	// OpenStax extracts: html2text flattens tables; nothing structural to fix
	// beyond the shared pass today, but the adapter is the seam where richer
	// extraction-side repair lands later.
	OpenStax = Adapter{Name: "openstax", Preprocess: preprocessGeneric}

	// ORCCA-derived teaching extracts (same family as Algebrica styling).
	ORCCA = Adapter{Name: "orcca", Preprocess: preprocessGeneric}
)

func stripSections(s, heading string) string {
	idx := strings.Index(s, "# "+heading)
	if idx < 0 {
		idx = strings.Index(s, "## "+heading)
	}
	if idx < 0 {
		return s
	}
	return s[:idx]
}

// ForSource picks the adapter for a registered lesson source path, sniffing
// the attribution header that extractors prepend to teaching/ files.
func ForSource(sourcePath, content string) Adapter {
	switch {
	case strings.HasPrefix(sourcePath, "algebrica/"):
		return Algebrica
	case strings.HasPrefix(sourcePath, "authored/"):
		return Authored
	case strings.Contains(content, "OpenStax"):
		return OpenStax
	case strings.Contains(content, "Oscar Levin") || strings.Contains(content, "Open Introduction"):
		return Levin
	default:
		return ORCCA
	}
}
