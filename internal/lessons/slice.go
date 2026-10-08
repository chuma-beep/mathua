package lessons

import (
	"fmt"
	"sort"
	"strings"
)

// A teaching slice is the smallest self-contained piece of instruction a learner needs to
// attempt the next practice item.
//
// It exists because of a measured defect. `/learn` renders a knowledge point's worked
// example, which is the lesson body's section named by the shard's `section` field. When that
// field is empty the server substituted the *whole lesson body*, so 69 of 1,971 steps — 43
// concepts, 33,784 words — served the entire reference article to a learner who had asked for
// one question. `discrete.logic.propositions` served its 2,724-word article three times over,
// 8,172 words in total, against a concept whose mastery threshold is twelve seconds.
//
// The reference article is not a bad piece of teaching. It is simply the wrong surface: it is
// written for a reader who wants the whole topic, and `/learn` is a learner who needs the next
// move. The reference panel (ADR-047) keeps the article available; the slice is what stands
// in front of the practice.
//
// Three levels, in order of preference:
//
//	section resolves -> that section, unchanged. 1,902 steps already take this path.
//	authored slice  -> a bounded rule, a few examples, maybe a misconception.
//	minimal slice   -> an honest statement that the step has no written explanation, and a
//	                  pointer to the reference. Never the article.
//
// The minimal level is deliberately thin. 60 of the 69 empty-section steps have no heading in
// their article at all, so a slice for them has to be *authored* rather than pointed at, and
// authoring 60 of them is content work. Serving an honest placeholder is better than serving
// 2,724 unexpected words, and the audit below makes the remaining backlog countable rather than
// invisible.

// Slice is an authored teaching slice, carried on a knowledge-point shard. Optional: a shard
// with no slice is served the minimal level.
type Slice struct {
	// Rule is the one core rule or definition. One to three sentences.
	Rule string `json:"rule"`
	// Examples are worked examples. Between one and six; see SliceBudget.
	Examples []Example `json:"examples,omitempty"`
	// Note is a short misconception warning. Optional, and at most one.
	Note string `json:"note,omitempty"`
}

// Example is one worked example: the item, and what it demonstrates.
type Example struct {
	Item string `json:"item"`
	Says string `json:"says,omitempty"`
}

// SliceBudget are the product boundaries for a slice. They are curriculum rules, not a
// word-count gate: nothing here rejects a formula or a concise definition, and a slice that
// genuinely needs more should be split into another knowledge point rather than truncated
// mid-sentence.
const (
	// MaxSliceWords bounds instructional prose before the first practice interaction.
	MaxSliceWords = 200
	// MaxSliceExamples bounds examples, per the brief's 3-6 total.
	MaxSliceExamples = 6
	// MinSliceWords is what a usable rule costs. Below this a slice says nothing.
	MinSliceWords = 12
)

// TeachingSlice returns the markdown `/learn` should show for one knowledge point.
//
// The bool reports whether a slice was found at all, which the server treats as "this step has
// no written explanation" rather than as an error. It is not a failure: a minimal slice is a
// valid, honest answer.
func (l *Loader) TeachingSlice(conceptID string, kp KP) (string, bool) {
	// 1. A section that resolves is already the right size — it was authored as a section.
	if body, ok := l.KPSectionBody(conceptID, kp.Section); ok {
		return body, true
	}
	// 2. An authored slice.
	if s := kp.Slice; s != nil && strings.TrimSpace(s.Rule) != "" {
		return renderSlice(kp.Label, *s), true
	}
	// 3. Nothing authored. Say so rather than serving the article.
	return renderMinimalSlice(kp.Label, l.Lesson(conceptID)), true
}

// renderMinimalSlice is the fallback. It must not pretend to teach and must not leak the
// subgoals, which are the procedure the learner is about to be asked to apply unaided.
func renderMinimalSlice(label string, lesson *Lesson) string {
	var b strings.Builder
	if label != "" {
		fmt.Fprintf(&b, "**%s**\n\n", label)
	}
	b.WriteString("There isn't a written explanation for this step yet. ")
	b.WriteString("Work the questions below — every answer explains the reasoning — ")
	b.WriteString("and open **Reference** if you want the full lesson text.")
	return b.String()
}

// renderSlice renders an authored slice to the markdown `/learn` already knows how to display:
// `$…$` for maths via KatexContent, `**bold**`, `-` lists.
func renderSlice(label string, s Slice) string {
	var b strings.Builder
	if label != "" {
		fmt.Fprintf(&b, "**%s**\n\n", label)
	}
	if rule := strings.TrimSpace(s.Rule); rule != "" {
		fmt.Fprintf(&b, "%s\n", rule)
	}
	if len(s.Examples) > 0 {
		b.WriteString("\nExamples\n\n")
		for _, ex := range s.Examples {
			item := strings.TrimSpace(ex.Item)
			if item == "" {
				continue
			}
			if says := strings.TrimSpace(ex.Says); says != "" {
				fmt.Fprintf(&b, "- %s — %s\n", item, says)
			} else {
				fmt.Fprintf(&b, "- %s\n", item)
			}
		}
	}
	if note := strings.TrimSpace(s.Note); note != "" {
		fmt.Fprintf(&b, "\nWatch out: %s\n", note)
	}
	return b.String()
}

// SliceLevel names how a knowledge point reaches the learner, for reporting.
type SliceLevel string

const (
	// SliceLevelSection is a section of the lesson body.
	SliceLevelSection SliceLevel = "section"
	// SliceLevelAuthored is a hand-written slice on the shard.
	SliceLevelAuthored SliceLevel = "authored"
	// SliceLevelMinimal is the honest placeholder.
	SliceLevelMinimal SliceLevel = "minimal"
)

// SliceFinding is one knowledge point's teaching material.
type SliceFinding struct {
	ConceptID string
	KPLabel   string
	Level     SliceLevel
	Words     int
	// FullBodyWords is how many words the learner would have read had the step fallen back to
	// the whole lesson. Zero for a section, since a section is not the article.
	FullBodyWords int
	FirstCard     bool
}

// SliceAudit walks every knowledge point in the corpus and reports how it reaches the learner,
// so the backlog of steps with no written explanation is countable. Nothing gates on this; it
// exists because the original defect was invisible, which is why it survived.
func (l *Loader) SliceAudit() []SliceFinding {
	var out []SliceFinding
	for conceptID, kps := range l.kps {
		lesson := l.concepts[conceptID]
		fullWords := 0
		if lesson != nil {
			fullWords = len(strings.Fields(lesson.Body))
		}
		for i, kp := range kps {
			level := SliceLevelMinimal
			full := 0
			switch {
			case kp.Section != "":
				if _, ok := l.KPSectionBody(conceptID, kp.Section); ok {
					level = SliceLevelSection
				}
			case kp.Slice != nil && strings.TrimSpace(kp.Slice.Rule) != "":
				level = SliceLevelAuthored
			default:
				full = fullWords
			}
			body, _ := l.TeachingSlice(conceptID, kp)
			out = append(out, SliceFinding{
				ConceptID:     conceptID,
				KPLabel:       kp.Label,
				Level:         level,
				Words:         len(strings.Fields(body)),
				FullBodyWords: full,
				FirstCard:     i == 0,
			})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].ConceptID != out[j].ConceptID {
			return out[i].ConceptID < out[j].ConceptID
		}
		return out[i].KPLabel < out[j].KPLabel
	})
	return out
}
