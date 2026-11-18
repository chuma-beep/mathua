package lessons

import "strings"

// Algebrica ends an article with a "Selected references" section: a bullet list
// of six to eight links to lecture notes published by universities. It is
// provenance for the source article, not teaching material for the concept, and
// the article body is rendered verbatim in `/study`, so the learner reads a wall
// of external links at the bottom of a page whose job was to explain one idea in
// two minutes. Twenty-four vendored files carry it.
//
// It is removed from the loaded body instead of from the files on disk for three
// reasons: `scripts/scrape_algebrica.py` and `scripts/refetch_algebrica.py` would
// put it back, the files are a vendored copy whose diff against upstream stays
// meaningful, and attribution does not depend on it. The licence obligation is
// met by THIRD_PARTY_NOTICES.md and the per-lesson credits line, not by a reader
// scrolling past a link list.
//
// The cut is deliberately narrow. Only the exact heading is matched, only when
// it is the last section, and never inside a fenced code block, because a wrong
// guess here deletes prose from a lesson and no test can tell you what it said.

const referenceHeading = "## Selected references"

// stripTrailingReferences removes a final "Selected references" section.
//
// It returns the body unchanged when the heading appears anywhere that is not
// the end of the document. Refusing to cut is the safe failure: a reference list
// that stays is untidy, while a cut at a mid-document heading would drop every
// section after it.
func stripTrailingReferences(body string) string {
	lines := strings.Split(body, "\n")
	inFence := false
	cut := -1

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		if cut >= 0 {
			// Any later heading means the list was not the last section, so the
			// body is served as it is.
			if strings.HasPrefix(line, "#") {
				return body
			}
			continue
		}
		if strings.EqualFold(strings.TrimRight(line, " \t\r"), referenceHeading) {
			cut = i
		}
	}

	if cut < 0 {
		return body
	}
	return strings.TrimRight(strings.Join(lines[:cut], "\n"), " \t\r\n") + "\n"
}
