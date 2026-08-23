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

var leakedSpanRe = regexp.MustCompile(`(?s)<span class="math-(display|inline)">([\s\S]*?)</span>`)

func stripSpanWrappers(s string) string {
	return leakedSpanRe.ReplaceAllString(s, "$2")
}

func preprocessGeneric(s string) string {
	s = entityReplacer.Replace(s)
	s = strings.ReplaceAll(s, `\amp`, "&")
	// Leaked wrapper spans from scrapers: keep inner content, drop the tags.
	s = stripSpanWrappers(s)
	return s
}

var doubledPunctRe = regexp.MustCompile(`\\\\([{},;|])`)
var doubledBeginEndRe = regexp.MustCompile(`\\\\(begin|end)\{`)
var brokenSqrtRe = regexp.MustCompile(`\\sqrt\}\{`)
var headingRe = regexp.MustCompile(`(?m)^(#{1,5})( )`)

// collapseDoubled normalizes the Algebrica scraping dialect where command
// backslashes were themselves escaped: \\{ -> \{, \\, -> \,, \\begin -> \begin.
// Row-break spacing like \\[6pt] is intentionally NOT touched here ([ and ]
// are excluded), and display delimiters \\[ \\] are consumed later by the
// delimiter pass after this collapse makes them single-form.
func collapseDoubled(s string) string {
	s = doubledBeginEndRe.ReplaceAllString(s, "\\$1{")
	s = doubledPunctRe.ReplaceAllString(s, "\\$1")
	return s
}

// normalizeDoubledDelims converts the doubled delimiter forms \\[ \\]
// \\( \\) into their single-backslash equivalents. Run AFTER row-break
// protection and BEFORE the standard single-form passes; safe because any
// bracket preceded by 2+ backslashes at this point is a real delimiter.
func normalizeDoubledDelims(s string) string {
	s = doubledDisplayDelimOpen.ReplaceAllString(s, "\\[")
	s = doubledDisplayDelimClose.ReplaceAllString(s, "\\]")
	s = doubledInlineDelim.ReplaceAllString(s, "\\$1")
	return s
}

// repairBrokenSqrt fixes the scraper artifact \sqrt}{x}} -> \sqrt{x}}.
func repairBrokenSqrt(s string) string {
	return brokenSqrtRe.ReplaceAllString(s, `\sqrt{`)
}

var rbProtectRe = regexp.MustCompile(`\\{2,}(\[[0-9]+(?:\.[0-9]+)?(?:pt|em|ex|mm|cm)\])`)
var rbRestoreRe = regexp.MustCompile(`%%MUARB:(\[[0-9]+(?:\.[0-9]+)?(?:pt|em|ex|mm|cm)\])%%`)

const rbToken = "%%MUARB:%s%%"

// protectRowbreaks freezes LaTeX row-break spacing (e.g. \\[6pt], written
// with 2-4 leading backslashes in scraped sources) behind inert tokens so the
// doubled-delimiter normalization cannot mistake their brackets for display
// math delimiters. restoreRowbreaks runs at the end of Canonicalize.
func protectRowbreaks(s string) string {
	return rbProtectRe.ReplaceAllStringFunc(s, func(m string) string {
		sub := rbProtectRe.FindStringSubmatch(m)
		if sub == nil {
			return m
		}
		return strings.Replace(rbToken, "%s", sub[1], 1)
	})
}

func restoreRowbreaks(s string) string {
	return rbRestoreRe.ReplaceAllString(s, "\\\\$1")
}

// demoteHeadings shifts every ATX heading one level deeper so scraped section
// titles stop competing with the page's own title typography.
func demoteHeadings(s string) string {
	return headingRe.ReplaceAllString(s, "#$1$2")
}

var (
	equationEnvRe  = regexp.MustCompile(`(?s)\\begin\{(equation\*?)\}(.*?)\\end\{equation\*?\}`)
	alignEnvRe     = regexp.MustCompile(`(?s)\\begin\{(align\*?)\}(.*?)\\end\{align\*?\}`)
	gatheredEnvRe  = regexp.MustCompile(`(?s)\\begin\{gathered\}(.*?)\\end\{gathered\}`)
	displayDelimRe = regexp.MustCompile(`\\\[([\s\S]*?)\\\]`)
	inlineDelimRe  = regexp.MustCompile(`\\\(([\s\S]*?)\\\\?\)|\\\(([\s\S]*?)\)`)
	mathTagRe      = regexp.MustCompile(`(?s)<math[^>]*>(.*?)</math>`)

	// Doubled-backslash delimiter dialect (Algebrica scrape artifact):
	// \\[ ... \\] display pairs and \\( ... \\) inline pairs.
	doubledDisplayDelimOpen  = regexp.MustCompile(`\\\\\[`)
	doubledDisplayDelimClose = regexp.MustCompile(`\\\\\]`)
	doubledInlineDelim       = regexp.MustCompile(`\\\\[()]`)
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
		// Align envs in this corpus always sit inside display delimiters;
		// the parent wrapper provides $$, so emit only the aligned body.
		return "\n\\begin{aligned}\n" + strings.TrimSpace(body) + "\n\\end{aligned}\n"
	})

	// \[ ...\ ] → $$ ... $$
	s = convertDisplays(s)

	// \( ... \) → $ ... $
	s = inlineDelimRe.ReplaceAllStringFunc(s, func(m string) string {
		g := inlineDelimRe.FindStringSubmatch(m)
		inner := g[1]
		if inner == "" {
			inner = g[2]
		}
		return "$" + inner + "$"
	})

	return restoreRowbreaks(s)
}

// escToken stands in for an escaped literal '\$' during validation so the
// delimiter-pairing and -counting passes never mistake it for a real inline
// math boundary. Chosen to be invisible to every structural check below.
const escToken = "@@MU-ESC-DOLLAR@@"

// neutralizeEscapedDelims rewrites '\$' (an escaped literal dollar sign) to
// escToken, honoring backslash parity: an odd run of backslashes before '$'
// means the dollar is escaped, while an even run means those are '\\' row
// breaks followed by a real delimiter.
func neutralizeEscapedDelims(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	run := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '\\' {
			run++
			continue
		}
		if c == '$' && run%2 == 1 {
			b.WriteString(escToken)
		} else {
			for j := 0; j < run; j++ {
				b.WriteByte('\\')
			}
			b.WriteByte(c)
		}
		run = 0
	}
	for j := 0; j < run; j++ {
		b.WriteByte('\\')
	}
	return b.String()
}

// Validate checks canonicalized content for structural problems: unbalanced
// braces inside math regions and stray/unpaired delimiters. It returns
// human-readable warnings; empty slice means clean. Escaped dollars (\$) are
// honored: they are neither delimiters nor pairing hazards.
func Validate(canonical string) []string {
	canonical = neutralizeEscapedDelims(canonical)

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

// convertDisplays rewrites single-form display delimiters. Inside markdown
// pipe-table rows it emits one-line inline math so GFM tables stay intact;
// everywhere else it emits standard block display math.
func convertDisplays(s string) string {
	matches := displayDelimRe.FindAllStringSubmatchIndex(s, -1)
	if len(matches) == 0 {
		return s
	}
	var b strings.Builder
	last := 0
	for _, m := range matches {
		start, end := m[0], m[1]
		b.WriteString(s[last:start])
		inner := strings.ReplaceAll(strings.TrimSpace(s[m[2]:m[3]]), "\n", " ")
		if onTableRowLine(s, start) {
			b.WriteString("$" + inner + "$")
		} else {
			b.WriteString("$$\n" + inner + "\n$$")
		}
		last = end
	}
	b.WriteString(s[last:])
	return b.String()
}

func onTableRowLine(s string, pos int) bool {
	lineStart := strings.LastIndexByte(s[:pos], '\n') + 1
	if strings.HasPrefix(strings.TrimSpace(s[lineStart:pos]), "|") {
		return true
	}
	rest := s[pos:]
	if nl := strings.IndexByte(rest, '\n'); nl >= 0 {
		rest = rest[:nl]
	}
	return strings.Contains(rest, "|")
}

// Adapters for the known content sources.
var (
	// Algebrica: doubled-backslash scrape dialect (\\{, \\,, \\[, \\begin),
	// broken \sqrt}{ artifacts, and oversized section headings.
	Algebrica = Adapter{Name: "algebrica", Preprocess: func(s string) string {
		s = protectRowbreaks(s)
		s = preprocessGeneric(s)
		s = stripSpanWrappers(s)
		s = collapseDoubled(s)
		s = normalizeDoubledDelims(s)
		s = repairBrokenSqrt(s)
		return demoteHeadings(s)
	}}

	// Authored: hand-written lessons already use $ / $$ conventions.
	Authored = Adapter{Name: "authored", Preprocess: preprocessGeneric}

	// Generators: live-generated problem text.
	Generators = Adapter{Name: "generators", Preprocess: preprocessGeneric}

	// Levin extracts: PreTeXt HTML → markdown keeps \begin{equation*} blocks,
	// \amp alignment markers, and full exercise sections.
	Levin = Adapter{Name: "levin", Preprocess: func(s string) string {
		s = preprocessGeneric(s)
		s = stripSections(s, "Exercises")
		return demoteHeadings(s)
	}}

	// OpenStax extracts: html2text flattens tables; demote headings so chapter
	// titles don't render at page-title scale.
	OpenStax = Adapter{Name: "openstax", Preprocess: func(s string) string {
		s = preprocessGeneric(s)
		return demoteHeadings(s)
	}}

	// ORCCA-derived teaching extracts (same family as Algebrica styling).
	ORCCA = Adapter{Name: "orcca", Preprocess: func(s string) string {
		s = preprocessGeneric(s)
		s = repairBrokenSqrt(s)
		return demoteHeadings(s)
	}}
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
