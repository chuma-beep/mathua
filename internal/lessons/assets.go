package lessons

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// Visual and tabular instructional content, extracted from a lesson body.
//
// This exists because the lesson model had nowhere to put a diagram. `Lesson` was
// `{Title, Body, Concepts}` and `KP` was `{label, section, subgoals}`, so the only way a
// figure could travel was as literal `![alt](path)` text inside one opaque string. That
// works for `/study`, which renders the body verbatim, and it works *by accident* for
// `/learn`, which renders one section of the body. It fails for everything else:
//
//   - Nothing can select the assets relevant to a specific learning step, because they are
//     not addressable.
//   - Nothing can be tested. A missing figure and a correctly-rendered one produce the same
//     body, so `scripts/audit_lessons.py` could only ever check the hand-maintained
//     concept->diagram map in `internal/engine`, never what the bodies actually reference.
//   - `KatexContent` rewrote relative image paths on the way to the browser
//     (`svg/x.svg` -> `/diagrams/algebrica/x.svg`) because there was no other place for that
//     rule to live, and nothing tested it.
//
// Assets are therefore parsed out of the body deterministically and exposed as data. The
// body is left byte-for-byte intact: it remains the canonical reference prose that `/study`
// renders, and it is the single source these are derived from, so the two cannot drift.

type AssetKind string

const (
	// AssetImage is a figure: a diagram, graph, photograph or plot.
	AssetImage AssetKind = "image"
	// AssetTable is a tabular block. Tables survive ingestion as GFM markdown inside the
	// body, so they are recognised as blocks rather than parsed into cells.
	AssetTable AssetKind = "table"
)

// Asset is one piece of instructional visual content, with everything needed to render it
// and to reason about it.
type Asset struct {
	// ID is stable for a given body: the source file's basename, disambiguated by a counter
	// when a body uses the same file twice. Referenced by `KP.AssetIDs`.
	ID string `json:"id"`
	// Kind is image or table.
	Kind AssetKind `json:"kind"`
	// Src is the reference exactly as written in the body. Relative paths are resolved by
	// the client, because only it knows the serving base.
	Src string `json:"src,omitempty"`
	// Alt is the markdown alt text. Frequently empty in the vendored corpus — 133 of 254
	// algebrica references are `![]` — which is why Caption matters and why a figure with
	// neither is worth reporting.
	Alt string `json:"alt,omitempty"`
	// Caption is a standalone italic line immediately after the figure, which is how both
	// corpora carry attribution and description.
	Caption string `json:"caption,omitempty"`
	// Line is the 1-based line in Body where the asset appears, so a renderer can interleave
	// assets with the prose that surrounds them.
	Line int `json:"line"`
	// Section is the heading the asset sits under, normalised the same way KP.Section is
	// matched. This is what makes "the assets relevant to this learning step" answerable.
	Section string `json:"section,omitempty"`
	// ConceptIDs are the concepts the lesson serves. Assets are lesson-scoped, so this is
	// constant within a lesson and exists so a consumer can filter without loading the map.
	ConceptIDs []string `json:"concept_ids,omitempty"`
	// Source is the attribution carried by the asset, if any. Never synthesised: an asset
	// with no attribution has none, which is different from having the wrong one.
	Source string `json:"source,omitempty"`
}

// imageRe matches a markdown image anywhere in the body, including mid-paragraph and with
// alt text wrapped across lines.
//
// Both shapes are real in the vendored corpus and both were missed by a line-anchored
// matcher, which is why the completeness test exists: `algebrica/trigonometry/
// sine-and-cosine.md:24` embeds a figure inside a sentence, and three
// probability/geometry lessons wrap the alt text over several lines. The alt pattern
// tolerates escaped brackets because one lesson's alt text literally contains `\[a,b\]`.
//
// A reference-style definition (`[label]: path`) is *not* matched. No corpus uses one, and
// supporting it would mean resolving document-scoped labels — the completeness test fails
// if a reference appears, so it cannot arrive unnoticed.
// Alt text may itself contain `]`, and the closing delimiter is the `]` followed by `(`.
// CommonMark allows brackets inside alt text and the client renderer follows it, so a pattern
// that stops at the first `]` is a parser that disagrees with the renderer it exists to
// describe: `numerical-integration-1.svg` carries the caption "The interval [a, b] has been
// divided into 6 subintervals", and `floor-and-ceiling-functions-1.svg` carries "[n, n+1)",
// whose bracket is never closed at all. The `]` alternative is what makes both work — the
// closing delimiter is tried first and falls back to consuming the bracket. Go's regexp
// prefers the leftmost-first alternative, which is the behaviour needed; the tempting
// `\](?!\()` is Perl syntax that does not compile here.
//
// The path must end in an image extension, which is what keeps this from eating the
// cross-reference links the algebrica corpus is full of: an exclamation mark in prose
// followed by `[continuous](<../continuous-functions/>)` is otherwise a perfect match for an
// image, and the corpus has 39 such links. Every figure reference in the corpus is a .svg or
// a .png, so requiring that costs nothing and removes the whole class.
var imageRe = regexp.MustCompile(
	`!\[((?:[^\]\\]|\\.|\])*)\]\(\s*([^)\s<]*\.(?:svg|png|jpe?g|webp|gif))\s*\)`)

// A paragraph break. A markdown image reference cannot span one, and matching per paragraph
// is what stops an exclamation mark in one sentence from pairing with a link's `](` several
// paragraphs later and swallowing the figures in between into its alt text.
var blankLineRe = regexp.MustCompile(`\n[ \t]*\n`)

// captionRe matches a standalone italic line, which is how both corpora write figure
// captions and attributions.
var captionRe = regexp.MustCompile(`^\*([^*].*?)\*$`)

// headingRe matches an ATX heading and captures its level and text.
var headingRe = regexp.MustCompile(`^(#{1,6})\s+(.+?)\s*$`)

// ExtractAssets walks a lesson body and returns every figure it references, in document
// order, each tagged with the section it sits under.
//
// Scans the whole body rather than line by line, because figures appear mid-sentence and
// their alt text wraps. Sections come from the nearest preceding heading, which requires
// tracking byte offset against the heading positions — hence the single pass with an index.
func ExtractAssets(body string, conceptIDs []string) []Asset {
	lines := strings.Split(body, "\n")
	// Offset of each line, so a byte index can be turned into a line number.
	starts := make([]int, len(lines)+1)
	off := 0
	for i, ln := range lines {
		starts[i] = off
		off += len(ln) + 1
	}
	starts[len(lines)] = off
	lineAt := func(idx int) int {
		lo, hi := 0, len(starts)-1
		for lo < hi {
			mid := (lo + hi + 1) / 2
			if starts[mid] <= idx {
				lo = mid
			} else {
				hi = mid - 1
			}
		}
		return lo + 1
	}

	var out []Asset
	seen := map[string]int{}

	// Current section, updated by the most recent heading at or before the match offset.
	sectionAt := func(idx int) string {
		ln := lineAt(idx)
		for i := ln - 1; i >= 0; i-- {
			if m := headingRe.FindStringSubmatch(lines[i]); m != nil {
				return normSectionKey(m[2])
			}
		}
		return ""
	}

	for _, chunk := range paragraphChunks(body) {
		for _, loc := range imageRe.FindAllStringSubmatchIndex(chunk.text, -1) {
			alt := strings.TrimSpace(chunk.text[loc[2]:loc[3]])
			src := strings.TrimSpace(chunk.text[loc[4]:loc[5]])
			section := sectionAt(chunk.offset + loc[0])
			name := path.Base(src)
			if name == "" || name == "." || name == "/" {
				name = "asset"
			}
			if n := seen[name]; n > 0 {
				seen[name] = n + 1
				name = fmt.Sprintf("%s-%d", name, n+1)
			} else {
				seen[name] = 1
			}
			a := Asset{
				ID:         name,
				Kind:       AssetImage,
				Src:        src,
				Alt:        alt,
				Line:       lineAt(chunk.offset + loc[0]),
				Section:    section,
				ConceptIDs: conceptIDs,
			}
			if cap := captionAfter(lines, a.Line-1); cap != "" {
				a.Caption = cap
			}
			out = append(out, a)
		}
	}
	return out
}

// paragraphChunks splits a body on blank lines, keeping each chunk's offset so line numbers
// stay absolute.
func paragraphChunks(body string) []struct {
	text   string
	offset int
} {
	var out []struct {
		text   string
		offset int
	}
	pos := 0
	for {
		end := len(body)
		next := len(body)
		if loc := blankLineRe.FindStringIndex(body[pos:]); loc != nil {
			end = pos + loc[0]
			next = pos + loc[1]
		}
		if end > pos {
			out = append(out, struct {
				text   string
				offset int
			}{body[pos:end], pos})
		}
		if next >= len(body) {
			return out
		}
		pos = next
	}
}

// BrokenImageRefs returns the offsets of `![` occurrences that are not part of a
// well-formed image reference.
//
// These exist in the vendored corpus: four lessons contain a bare `![` with the alt text and
// path lost by the scrape. They render as a literal `![`, so they are a data defect rather
// than a shape to parse, and they are surfaced so the gate can report them instead of the
// extractor quietly skipping them.
func BrokenImageRefs(body string) []int {
	// Ranges covered by a well-formed reference, so that a `![` consumed as part of someone's
	// alt text is not reported as broken in its own right. Matching an unanchored regex at
	// each `![` is not equivalent: the first version of this did that, and a bare `![` line
	// was "rescued" by a well-formed figure further down the same file, which hid two of the
	// four real defects.
	covered := make([][2]int, 0, 8)
	for _, chunk := range paragraphChunks(body) {
		for _, loc := range imageRe.FindAllStringSubmatchIndex(chunk.text, -1) {
			covered = append(covered, [2]int{chunk.offset + loc[0], chunk.offset + loc[1]})
		}
	}
	starts := make(map[int]bool, len(covered))
	for _, c := range covered {
		starts[c[0]] = true
	}

	var out []int
	// Rune-at-a-time rather than byte-at-a-time, which is lint compliance and also the
	// clearer expression of intent: `range` over a string still yields the byte offset of each
	// rune's first byte, so `i` remains the offset `covered` and the caller's line arithmetic
	// both expect.
	//
	// Testing `body[i+1] == '['` after a `!` is safe even though `i` steps by rune width: `!`
	// is one byte, so `i+1` is the start of the next character, and if that character is
	// multi-byte then `body[i+1]` is a continuation byte — always >= 0x80, never `[` (0x5B).
	for i, r := range body {
		if r != '!' {
			continue
		}
		if i+1 >= len(body) || body[i+1] != '[' {
			continue
		}
		if starts[i] {
			continue
		}
		inside := false
		for _, c := range covered {
			if i > c[0] && i < c[1] {
				inside = true
				break
			}
		}
		if !inside {
			out = append(out, i)
		}
	}
	return out
}

// captionAfter returns the caption following line `i`, if the next non-blank line is a
// standalone italic. It does not skip over other content, so a caption is never stolen from
// a paragraph that merely happens to be italic further down.
func captionAfter(lines []string, i int) string {
	for j := i + 1; j < len(lines); j++ {
		t := strings.TrimSpace(lines[j])
		if t == "" {
			continue
		}
		if m := captionRe.FindStringSubmatch(t); m != nil {
			return strings.TrimSpace(m[1])
		}
		return ""
	}
	return ""
}

// corpusOf derives the corpus name from a lesson file path, so `Source` says where the
// prose came from rather than repeating the filename.
//
// The first path segment under data/lessons, not the containing directory: the sources are
// `algebrica/functions/floor-and-ceiling-functions.md`, and returning `algebrica/functions`
// would make the field a per-topic value that happens to look like a corpus. A lesson with
// no directory reports "authored", which is what a hand-written lesson is.
//
// This is the corpus, not the attribution of an individual claim. A figure's own
// attribution is whatever its caption carries; nothing is invented to fill the field.
func corpusOf(p string) string {
	p = filepath.ToSlash(p)
	i := strings.Index(p, "lessons/")
	if i < 0 {
		if filepath.Base(p) == "" || filepath.Dir(p) == "." {
			return "authored"
		}
		return "authored"
	}
	rest := p[i+len("lessons/"):]
	if j := strings.Index(rest, "/"); j >= 0 {
		return rest[:j]
	}
	return "authored"
}

// AssetsFor returns every asset in the concept's lesson.
func (l *Loader) AssetsFor(conceptID string) []Asset {
	lesson := l.concepts[conceptID]
	if lesson == nil {
		return nil
	}
	return lesson.Assets
}

// KPAssets returns the assets relevant to one learning step: those in the section that KP's
// worked example is drawn from.
//
// This is the whole point of extracting them. `/study` shows every asset in document order,
// which is the reference view; `/learn` shows only these, so a shorter surface means less
// irrelevant text rather than less instructional information.
//
// An authored `asset_ids` list on the KP shard wins, so Pass 3 can point a step at a figure
// that lives outside its section without moving prose around.
func (l *Loader) KPAssets(conceptID string, kp KP) []Asset {
	lesson := l.concepts[conceptID]
	if lesson == nil {
		return nil
	}
	if len(kp.AssetIDs) > 0 {
		byID := make(map[string]Asset, len(lesson.Assets))
		for _, a := range lesson.Assets {
			byID[a.ID] = a
		}
		out := make([]Asset, 0, len(kp.AssetIDs))
		for _, id := range kp.AssetIDs {
			if a, ok := byID[id]; ok {
				out = append(out, a)
			}
		}
		return out
	}
	want := normSectionKey(kp.Section)
	if want == "" {
		// No section means no figures either. Returning the whole lesson's asset list was the
		// same mistake as serving the whole lesson body, in a different form: 40 figures were
		// attached to fallback steps that have no section to place them in, so a learner
		// meeting one concept's opening card was shown every figure the article contains. A
		// shard that genuinely wants a figure names it in asset_ids above.
		return nil
	}
	var out []Asset
	for _, a := range lesson.Assets {
		if a.Section == want {
			out = append(out, a)
		}
	}
	return out
}
