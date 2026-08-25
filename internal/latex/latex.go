// Package latex normalizes mathematics notation from heterogeneous content
// sources into a single KaTeX-friendly dialect ($ inline, $$ display) at
// ingestion time, so neither storage nor the frontend ever sees a source's
// quirks.
//
// The pipeline is region-based: latexnorm.Scan segments the text into Prose
// and math regions with full contextual rules (escapes, text groups,
// line-anchored blocks), and normalization is applied where it is valid —
// inside math — while Prose dollars are escaped so downstream markdown
// parsers never mistake prices for delimiters. Validation reports structural
// problems as warnings attached to the canonicalized output; callers decide
// whether to reject, quarantine, or serve with logging.
package latex

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/chuma-beep/latexnorm"
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

var doubledPunctRe = regexp.MustCompile(`\\\\([{},;|])`)
var doubledBeginEndRe = regexp.MustCompile(`\\\\(begin|end)\{`)
var brokenSqrtRe = regexp.MustCompile(`\\sqrt\}\{`)

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

var (
	doubledDisplayDelimOpen  = regexp.MustCompile(`\\\\\[`)
	doubledDisplayDelimClose = regexp.MustCompile(`\\\\\]`)
	doubledInlineDelim       = regexp.MustCompile(`\\\\([()])`)
)

// repairBrokenSqrt fixes the scraper artifact \sqrt}{x}} -> \sqrt{x}}.
func repairBrokenSqrt(s string) string {
	return brokenSqrtRe.ReplaceAllString(s, `\sqrt{`)
}

var rbProtectRe = regexp.MustCompile(`\\{2,}(\[[0-9]+(?:\.[0-9]+)?\s?(?:pt|em|ex|mm|cm)\])`)
var rbRestoreRe = regexp.MustCompile(`%%MUARB:(\[[0-9]+(?:\.[0-9]+)?\s?(?:pt|em|ex|mm|cm)\])%%`)

const rbToken = "%%MUARB:%s%%"

// protectRowbreaks freezes LaTeX row-break spacing (e.g. \\[6pt], written
// with 2-4 leading backslashes in scraped sources) behind inert tokens so
// the doubled-delimiter normalization cannot mistake their brackets for
// display math delimiters. restoreRowbreaks runs at the end of Canonicalize.
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

// pureCurrencyRe matches content that is a price or measurement — digits,
// separators, optional ÷/× chains, trailing punctuation — and nothing else.
var pureCurrencyRe = regexp.MustCompile(`^[\d.,]+([ \t]?[÷×][ \t]?\d[\d.,]*)*[.,;:!?]?$`)

// Brace-corruption repairs for scrape artifacts where an argument's opening
// brace was written as a closing one: \frac}{a}}{b}}, \sqrt{}{x}},
// e^}{x^2}, \mathrm}{d}. Each pass handles one shape; run via
// repairBraces until stable.
var (
	fracCorruptRe   = regexp.MustCompile(`\\([dt]?frac)\}\{([^{}]*)\}\{([^{}]*)\}\}`)
	sqrtCorruptRe   = regexp.MustCompile(`\\sqrt\{\}\{([^{}]*)\}\}`)
	scriptCorruptRe = regexp.MustCompile(`([_^])\}\{`)
	argCorruptRe    = regexp.MustCompile(`\\(mathrm|text|mbox|mathbf|mathit|mathsf|mathtt|textrm|textbf|textit|textup|textnormal|operatorname)\}\{`)
	punctEscapeRe   = regexp.MustCompile(`\\([+.=])`)
)

func repairBraces(s string) string {
	for i := 0; i < 5; i++ {
		prev := s
		s = fracCorruptRe.ReplaceAllString(s, `\$1{$2}{$3}`)
		s = sqrtCorruptRe.ReplaceAllString(s, `\sqrt{$1}`)
		s = scriptCorruptRe.ReplaceAllString(s, `$1{`)
		s = argCorruptRe.ReplaceAllString(s, `\$1{`)
		if s == prev {
			break
		}
	}
	return s
}

// rawDollarRe finds a dollar not already escaped (no backslash before it).
var rawDollarRe = regexp.MustCompile(`(^|[^\\])\$`)

// hasDigitRe / proseMathMarkerRe support the Inline prose-shape rule.
var (
	hasDigitRe        = regexp.MustCompile(`\d`)
	proseMathMarkerRe = regexp.MustCompile(`[\\^_{}\[\]+\-=*/<>()]`)
)

// escapeDollars escapes every unescaped dollar in a Prose region. Already
// escaped ones (from earlier fixed-point passes) are left alone.
func escapeDollars(s string) string {
	return rawDollarRe.ReplaceAllStringFunc(s, func(m string) string {
		return m[:len(m)-1] + `\$`
	})
}

// preprocessGeneric runs for every source: HTML entity cleanup, leaked
// wrapper spans, row-break protection, and doubled-backslash dialect
// collapse (the doubled dialect is not Algebrica-exclusive — authored
// lessons use it too).
func preprocessGeneric(s string) string {
	s = entityReplacer.Replace(s)
	s = strings.ReplaceAll(s, `\amp`, "&")
	s = stripSpanWrappers(s)
	s = protectRowbreaks(s)
	s = collapseDoubled(s)
	s = normalizeDoubledDelims(s)
	return s
}

var tikzRe = regexp.MustCompile(`(?s)\\begin\{tikzpicture\}.*?\\end\{tikzpicture\}`)

// mathNorm normalizes the content of a math region. Prose rules never apply
// here; this is the only place dialect rewrites are valid. align is
// rewritten to aligned because KaTeX renders align only at the top of a
// display block, while the corpus nests it inside $$...$$. TikZ diagrams
// cannot render in KaTeX and are stripped.
func mathNorm(s string) string {
	s = protectRowbreaks(s)
	s = collapseDoubled(s)
	s = repairBraces(s)
	s = strings.ReplaceAll(s, `\begin{align}`, `\begin{aligned}`)
	s = strings.ReplaceAll(s, `\end{align}`, `\end{aligned}`)
	s = tikzRe.ReplaceAllString(s, "")
	s = punctEscapeRe.ReplaceAllString(s, "$1")
	return restoreRowbreaks(s)
}

var mathTagRe = regexp.MustCompile(`(?s)<math[^>]*>(.*?)</math>`)

// Canonicalize rewrites every recognized math construct to exactly two
// delimiters: inline $...$ and display $$...$$ (KaTeX-native). Source
// quirks are absorbed first via the adapter, then the text is segmented
// into regions and each region is normalized by its own rules:
//
//	Prose            every raw $ is escaped — prices stay text
//	math regions     dialect collapse only; delimiters are re-emitted
//	Inline (paired)  pure-numeric content is treated as currency and escaped
//
// Deeply-escaped scrape artifacts shed one backslash level per pass, so the
// pipeline runs to a fixed point: the returned string is stable under
// Canonicalize (C(C(x)) == C(x)), capped at maxCanonPasses.
func Canonicalize(s string, a Adapter) string {
	prev := s
	for i := 0; i < maxCanonPasses; i++ {
		next := canonicalizeOnce(prev, a)
		if next == prev {
			return next
		}
		prev = next
	}
	return prev
}

const maxCanonPasses = 8

func canonicalizeOnce(s string, a Adapter) string {
	if a.Preprocess != nil {
		s = a.Preprocess(s)
	}

	// Legacy wiki content: unwrap raw <math> wrappers into display math
	// before scanning.
	for {
		m := mathTagRe.FindStringSubmatchIndex(s)
		if m == nil {
			break
		}
		inner := s[m[2]:m[3]]
		s = s[:m[0]] + "$$" + inner + "$$" + s[m[1]:]
	}

	var b strings.Builder
	for _, r := range latexnorm.Scan(s) {
		switch r.Kind {
		case latexnorm.Prose:
			// Dollars in prose are prices or punctuation, never delimiters.
			b.WriteString(escapeDollars(r.Content))
		case latexnorm.Inline:
			// A scanner-validated pair can still be prose: word-problem
			// prices ($5, saves $10) pair on one line but carry no math
			// markers. Digits + no markers + (pure numeric or containing a
			// space) means currency/prose; anything else is math.
			content := r.Content
			if hasDigitRe.MatchString(content) &&
				!proseMathMarkerRe.MatchString(content) &&
				(pureCurrencyRe.MatchString(content) || strings.Contains(content, " ")) {
				b.WriteString(`\$` + content + `\$`)
				continue
			}
			b.WriteString("$" + mathNorm(content) + "$")
		case latexnorm.InlineBlock:
			// Single-$ display dialect (line-anchored): promote to $$.
			b.WriteString("$$\n" + strings.TrimSpace(mathNorm(r.Content)) + "\n$$")
		case latexnorm.Display:
			b.WriteString("$$" + mathNorm(r.Content) + "$$")
		case latexnorm.Paren:
			b.WriteString("$" + mathNorm(r.Content) + "$")
		case latexnorm.Bracket:
			b.WriteString("$$\n" + strings.TrimSpace(mathNorm(r.Content)) + "\n$$")
		case latexnorm.Env:
			b.WriteString(envOutput(r))
		}
	}
	return restoreRowbreaks(b.String())
}

// envOutput rewrites environment regions. Equation environments become
// display math; align environments become aligned blocks (the corpus's align
// envs conventionally sit inside display delimiters, so no $$ is added);
// tikzpicture diagrams cannot render in KaTeX and are stripped; everything
// else passes through with normalized content.
func envOutput(r latexnorm.Region) string {
	name := r.Opener
	name = strings.TrimPrefix(name, `\begin{`)
	name = strings.TrimSuffix(name, `}`)
	inner := strings.TrimSpace(r.Content)
	switch name {
	case "tikzpicture":
		return ""
	case "equation", "equation*":
		return "$$" + mathNorm(inner) + "$$"
	case "align", "align*":
		inner = mathNorm(inner)
		inner = strings.Replace(inner, `\begin{aligned}`, "", 1)
		inner = strings.Replace(inner, `\end{aligned}`, "", 1)
		return "\n\\begin{aligned}\n" + strings.TrimSpace(inner) + "\n\\end{aligned}\n"
	default:
		return r.Opener + mathNorm(r.Content) + r.Closer
	}
}

// Validate checks canonicalized content for structural problems: unbalanced
// braces inside math regions and stray unescaped dollars in prose. It
// returns human-readable warnings; empty slice means clean.
func Validate(canonical string) []string {
	var warnings []string
	for _, r := range latexnorm.Scan(canonical) {
		switch r.Kind {
		case latexnorm.Prose:
			if n := len(rawDollarRe.FindAllString(r.Content, -1)); n > 0 {
				warnings = append(warnings, fmt.Sprintf("odd number of inline delimiters (%d) — possible unpaired '$'", n))
			}
		case latexnorm.Inline, latexnorm.Paren:
			if msg := checkBraces(r.Content); msg != "" {
				warnings = append(warnings, "inline math: "+msg)
			}
		case latexnorm.InlineBlock, latexnorm.Display, latexnorm.Bracket, latexnorm.Env:
			if msg := checkBraces(r.Content); msg != "" {
				warnings = append(warnings, "display math: "+msg)
			}
		}
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
	// Algebrica: row-break protection and broken \sqrt}{ artifacts; the
	// doubled-backslash collapse runs for every source (preprocessGeneric).
	Algebrica = Adapter{Name: "algebrica", Preprocess: func(s string) string {
		return repairBrokenSqrt(preprocessGeneric(s))
	}}

	// Authored: hand-written lessons already use $ / $$ conventions.
	Authored = Adapter{Name: "authored", Preprocess: preprocessGeneric}

	// Generators: live-generated problem text.
	Generators = Adapter{Name: "generators", Preprocess: preprocessGeneric}

	// Levin extracts: PreTeXt HTML → markdown keeps \begin{equation*} blocks,
	// \amp alignment markers, and full exercise sections.
	Levin = Adapter{Name: "levin", Preprocess: func(s string) string {
		return stripSections(preprocessGeneric(s), "Exercises")
	}}

	// OpenStax extracts: html2text flattens tables.
	OpenStax = Adapter{Name: "openstax", Preprocess: preprocessGeneric}

	// ORCCA-derived teaching extracts (same family as Algebrica styling).
	ORCCA = Adapter{Name: "orcca", Preprocess: func(s string) string {
		return repairBrokenSqrt(preprocessGeneric(s))
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
