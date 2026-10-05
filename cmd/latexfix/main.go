// latexfix repairs the LaTeX defects that `make latex-normalize` finds, and only
// those that are mechanical. Anything needing a human decision is left alone and
// reported, so this tool can be run over the corpus without risk.
//
// Three rules, all of them learned by getting it wrong first:
//
//  1. Only regions that latexnorm.Scan classifies as math are edited. The corpus is
//     mostly English prose, where a typographic apostrophe or an en dash is
//     *correct*. A blind find-and-replace would fix the math and corrupt the prose.
//
//  2. A region's Content EXCLUDES its delimiters — input[Start:End] is
//     Opener + Content + Closer. Writing Content back over Start:End strips the
//     `$…$` off every span it touches, which turns rendered math into literal text
//     in 33 files and reads as content loss. The delimiters go back on.
//
//  3. Nothing inside \text{…} is touched. `\overset{\text{—}}{16}` means "16
//     repeats"; the em dash is a notation glyph, and rewriting it to `-` changes
//     what the lesson says. Prose inside math is prose, even when the surrounding
//     dollars make it look like an expression.
//
//  4. The Algebrica scrape is double-escaped: the source holds the literal
//     characters \( and \), which latexnorm reads as an escaped backslash followed
//     by a paren — so an inline formula sits in *prose* as far as any region-based
//     tool is concerned, and the canonicalizer is what turns it into `$ … $`.
//     That put 14,936 formulas out of reach of the rules above. Un-escaping fixes
//     the honesty of the file and lets the math be found at all.
//
//     Only \( and \) are rewritten. \[ is left alone on purpose: `\\[2pt]` is a
//     row break inside an align environment, not a display-math opener, and 1,157
//     of the 5,657 occurrences really are row breaks.
//
// Usage:
//
//	go run ./cmd/latexfix            # dry run, report only
//	go run ./cmd/latexfix -w         # write fixes in place
//	go run ./cmd/latexfix data/lessons/authored
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/chuma-beep/latexnorm"
)

type fix struct {
	old string
	new string
	// skipInText keeps this rule out of \text{…} and friends, where the
	// characters are prose rather than notation.
	skipInText bool
	// attachedOnly requires the character to follow a symbol, which is what
	// separates a derivative's apostrophe from a quotation mark.
	attachedOnly bool
	why          string
}

var fixes = []fix{
	// KaTeX renders a Unicode apostrophe happily, which is exactly why no
	// existing gate noticed that derivatives were written `f’(x)`.
	//
	// Attached-only, and this is the whole subtlety. `f’(x)` wants ASCII. But
	// `{a < b ‘a\ \text{is less than}\ b’}` uses the same characters as quotation
	// marks around an English phrase, and rewriting those to `'` makes LaTeX render
	// primes on a and b — a silent change to what the lesson says. An apostrophe is
	// attached to the thing it qualifies; a quotation mark is not. There is no rule
	// for `‘`: an opening curly quote in math is always a quotation mark, never an
	// apostrophe, and the right repair is a \text{} group that needs a human.
	{old: "’", new: "'", skipInText: true, attachedOnly: true, why: "attached Unicode apostrophe from the scrape; LaTeX math wants ASCII"},
	{old: "✓", new: "", why: "scrape artifact: a correctness tick that landed inside the math span"},
	{old: "&gt;", new: ">", why: "HTML entity leaked into math from the scrape"},
	{old: "&lt;", new: "<", why: "HTML entity leaked into math from the scrape"},
	{old: "&#39;", new: "'", why: "HTML entity leaked into math from the scrape"},
	{old: "&nbsp;", new: " ", why: "HTML entity leaked into math from the scrape"},
	{old: "&amp;", new: `\&`, why: "HTML entity leaked into math; an alignment tab is not intended here"},
}

// bareNumeral matches the content the canonicalizer refuses to serve as math.
// ADR-016: a bare numeral set as inline math gets escaped to \$42\$, so the learner
// sees literal dollar signs. Bare numerals must be plain text.
var bareNumeral = regexp.MustCompile(`^[0-9][0-9.,]*$`)

// textGroup finds the byte ranges of \text{…}-style groups inside a math region.
// Nested braces are tracked so `\text{a {b} c}` closes at the right place.
func textGroup(content string) [][2]int {
	var out [][2]int
	openers := []string{`\text{`, `\mbox{`, `\emph{`, `\textit{`, `\textbf{`}
	// Advance by rune width rather than by one byte (go-ai-lint AIL070). The returned ranges
	// are byte offsets, which is what the caller slices with, so `i` remains a byte index.
	//
	// The byte-at-a-time version could not produce a wrong answer, and the reason is worth
	// recording because it is not obvious: a multi-byte UTF-8 sequence is made only of bytes
	// ≥ 0x80, and every opener here begins with a backslash at 0x5C, so a fragment starting
	// mid-rune can never begin with `\text{`. The scan was safe because of a property of
	// UTF-8, not because of anything in this code — and it would stop being safe the moment
	// an opener were matched on a byte that can appear inside a rune.
	//
	// Decoding the rune makes that guarantee local and explicit instead of resting on the
	// encoding.
	// The index is managed by hand rather than with `for i := range content`, because this
	// loop jumps `i` forward to skip a matched group. A range loop reassigns its index at the
	// top of every iteration, so that jump would be discarded and the scan would restart
	// inside the group it just consumed.
	for i := 0; i < len(content); {
		matched := ""
		for _, o := range openers {
			if strings.HasPrefix(content[i:], o) {
				matched = o
				break
			}
		}
		if matched == "" {
			_, w := utf8.DecodeRuneInString(content[i:])
			i += w
			continue
		}
		depth, j := 1, i+len(matched)
		for ; j < len(content) && depth > 0; j++ {
			switch content[j] {
			case '{':
				depth++
			case '}':
				depth--
			}
		}
		out = append(out, [2]int{i, j})
		i = j
	}
	return out
}

// attachedBefore reports whether the byte at pos continues a symbol — a letter,
// digit, closing brace or closing paren. That is what makes `f’` an apostrophe
// while ` b’` is a closing quotation mark.
func attachedBefore(s string, pos int) bool {
	if pos == 0 {
		return false
	}
	c := s[pos-1]
	return c == '}' || c == ')' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

func inAny(ranges [][2]int, pos int) bool {
	for _, r := range ranges {
		if pos >= r[0] && pos < r[1] {
			return true
		}
	}
	return false
}

// applyFixes rewrites the rules it can inside one math region and reports what it
// did. It returns the region's Content, never the delimited slice.
func applyFixes(content string) (string, []string) {
	out := content
	var reasons []string
	for _, f := range fixes {
		if !strings.Contains(out, f.old) {
			continue
		}
		if f.skipInText {
			// Recomputed per rule: an earlier rule may have changed the length, and
			// a stale mask would then guard the wrong bytes.
			ranges := textGroup(out)
			var b strings.Builder
			prev := 0
			for i := 0; i < len(out); i++ {
				if !strings.HasPrefix(out[i:], f.old) || inAny(ranges, i) {
					continue
				}
				if f.attachedOnly && !attachedBefore(out, i) {
					continue
				}
				b.WriteString(out[prev:i])
				b.WriteString(f.new)
				prev = i + len(f.old)
				i = prev - 1
			}
			if prev == 0 {
				continue
			}
			b.WriteString(out[prev:])
			if b.String() != out {
				out = b.String()
				reasons = append(reasons, fmt.Sprintf("%q→%q: %s", f.old, f.new, f.why))
			}
			continue
		}
		next := strings.ReplaceAll(out, f.old, f.new)
		if next != out {
			out = next
			reasons = append(reasons, fmt.Sprintf("%q→%q: %s", f.old, f.new, f.why))
		}
	}
	return out, reasons
}

// unescapeDelimiters collapses the scrape's doubled backslash on inline-math
// delimiters. It is deliberately not part of applyFixes: `\(` sits in a *prose*
// region, which is the whole problem.
// rowBreakLen matches what follows a `\\[` that is a row break rather than a
// display-math opener: an optional length and unit, then the closing bracket.
var rowBreakLen = regexp.MustCompile(`^\[?[0-9.]*[a-z]*\]`)

func unescapeDelimiters(src string) (string, int) {
	var b strings.Builder
	n := 0
	// Rune-at-a-time rather than byte-at-a-time (go-ai-lint AIL070).
	//
	// To be accurate about what this changes: the byte-at-a-time version was **not** losing
	// anything. `strings.Builder.WriteByte` appends the raw byte and the loop visited every
	// byte exactly once, so it reconstructed the input byte-for-byte — verified for em
	// dashes, curly quotes, Greek letters and ellipses. The lint is a robustness finding, not
	// a correctness one, and this rewrite should not be described as a bug fix.
	//
	// What it does buy is that the invariant stops being incidental. Writing the whole rune
	// in one step means a future branch that jumps `i` forward cannot drop a byte by
	// accident, because the copy and the advance are now the same operation instead of two
	// that have to agree.
	//
	// The index is advanced explicitly rather than with a range loop, because the branches
	// below consume a fixed number of bytes; a range loop would discard those jumps.
	for i := 0; i < len(src); {
		if strings.HasPrefix(src[i:], `\\(`) || strings.HasPrefix(src[i:], `\\)`) {
			b.WriteByte('\\')
			b.WriteByte(src[i+2])
			// Consume i, i+1 and i+2 in one step.
			i += 3
			n++
			continue
		}
		if strings.HasPrefix(src[i:], `\\]`) {
			b.WriteString(`\]`)
			i += 3
			n++
			continue
		}
		if strings.HasPrefix(src[i:], `\\[`) {
			if rowBreakLen.MatchString(src[i+3:]) {
				// A row break: leave both backslashes alone. Written together here because
				// the next iteration no longer starts on a backslash pair.
				b.WriteString(src[i : i+2])
				i += 2
				continue
			}
			b.WriteString(`\[`)
			i += 3
			n++
			continue
		}
		_, w := utf8.DecodeRuneInString(src[i:])
		b.WriteString(src[i : i+w])
		i += w
	}
	return b.String(), n
}

type change struct {
	Offset  int
	Before  string
	After   string
	Reasons []string
}

// repair applies both passes: un-escape the doubled delimiters, then fix the math
// regions the un-escaping made visible. It always returns a usable string — the
// caller writes when out != src — so a file that needs only un-escaping is still
// written even though no region changed.
func repair(src string) (string, []change) {
	out, unescaped := unescapeDelimiters(src)
	changes := scanAndFix(&out)

	if unescaped > 0 {
		// Reported without byte offsets: after the rewrite they refer to a
		// different string, and a stale offset is worse than none.
		changes = append([]change{{
			Before:  snippet(src, 120),
			After:   snippet(out, 120),
			Reasons: []string{fmt.Sprintf("unescaped %d doubled \\( / \\) delimiters: the scrape wrote \\( which latexnorm reads as prose, so the formula was invisible to every region-based tool", unescaped)},
		}}, changes...)
	}
	return out, changes
}

// scanAndFix rewrites the rules it can inside each math region of *src, in place.
func scanAndFix(src *string) []change {
	var b strings.Builder
	prev := 0
	var changes []change
	current := *src

	for _, r := range latexnorm.Scan(current) {
		if r.Kind == latexnorm.Prose {
			continue
		}
		content, reasons := applyFixes(r.Content)
		if len(reasons) == 0 {
			continue
		}

		// A repaired span that is now a bare numeral must lose its delimiters, not
		// keep them. `$42✓` -> `$42` reads as math and the canonicalizer then
		// escapes it to `\$42\$`, which the learner sees as literal dollar signs --
		// so removing the tick would have *created* a defect while fixing another.
		// ADR-016 puts this in content, and this is content.
		if bareNumeral.MatchString(strings.TrimSpace(content)) {
			content = strings.TrimSpace(content)
			reasons = append(reasons, "span is now a bare numeral, so its $ delimiters go too: the canonicalizer escapes a bare numeral set as math and the learner would see literal dollar signs (ADR-016)")
			b.WriteString(current[prev:r.Start])
			b.WriteString(content)
			prev = r.End
			changes = append(changes, change{r.Start, r.Content, content, reasons})
			continue
		}

		changes = append(changes, change{r.Start, r.Content, content, reasons})
		// Opener + fixed Content + Closer. r.Start:r.End is the delimited slice, so
		// the delimiters must be re-attached or the math stops being math.
		b.WriteString(current[prev:r.Start])
		b.WriteString(r.Opener)
		b.WriteString(content)
		b.WriteString(r.Closer)
		prev = r.End
	}
	if len(changes) > 0 {
		b.WriteString(current[prev:])
		*src = b.String()
	}
	return changes
}

func main() {
	write := flag.Bool("w", false, "write fixes in place (default: report only)")
	flag.Parse()

	root := "data/lessons"
	if flag.NArg() > 0 {
		root = flag.Arg(0)
	}

	enc := json.NewEncoder(os.Stdout)
	touched, changed := 0, 0

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		touched++

		src := string(raw)
		out, changes := repair(src)
		if out == src {
			return nil
		}
		for _, c := range changes {
			_ = enc.Encode(map[string]any{
				"level":  "WARN",
				"msg":    "repaired math span",
				"file":   filepath.ToSlash(rel),
				"offset": c.Offset,
				"before": snippet(c.Before, 120),
				"after":  snippet(c.After, 120),
				"rules":  c.Reasons,
			})
		}
		changed++
		if *write {
			if err := os.WriteFile(path, []byte(out), 0o644); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "latexfix:", err)
		os.Exit(1)
	}

	mode := "dry run (pass -w to write)"
	if *write {
		mode = "written"
	}
	fmt.Fprintf(os.Stderr, "latexfix: %d files scanned, %d changed, %s\n", touched, changed, mode)
}

func snippet(s string, max int) string {
	flat := strings.Join(strings.Fields(s), " ")
	if len(flat) <= max {
		return flat
	}
	head := (max - 1) * 7 / 10
	return flat[:head] + "…" + flat[len(flat)-(max-1-head):]
}
