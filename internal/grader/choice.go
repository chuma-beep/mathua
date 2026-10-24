package grader

import (
	"regexp"
	"strings"
)

var singleLetterRe = regexp.MustCompile(`\b[A-Za-z]\b`)
var choicePrefixRe = regexp.MustCompile(`^(?i)(?:option|choice)?\s*\(?([A-E])(?:\)|\.|\:)?$`)

// A plain multi-word phrase, the only shape where the answer field can lose a
// word boundary without changing meaning.
var wordPhraseRe = regexp.MustCompile(`^[A-Za-z]+(?:\s+[A-Za-z]+)+$`)

// samePhrase accepts a phrase whose spaces were lost on the way out of the
// math field.
//
// MathLive serialises `line segment` as `l i n e s e g m e n t`: it treats
// adjacent letters in math mode as a product and spaces every one. The client
// rejoins single-letter runs, which turns the phrase into `linesegment`, so an
// exact comparison against `line segment` marks a correct answer wrong.
//
// Two guards keep this narrow. The expected answer must be letters and spaces
// only, so `5 R 3` (remainder notation) and `2pi r` (a product) never come
// here. The answer must be letters only, so a collapsed phrase cannot absorb
// digits, operators or a fraction.
func samePhrase(expected, answer string) bool {
	e := strings.TrimRight(strings.ToLower(strings.TrimSpace(expected)), ".!?,;")
	a := strings.TrimRight(strings.ToLower(strings.TrimSpace(answer)), ".!?,;")
	if !wordPhraseRe.MatchString(e) || a == "" {
		return false
	}
	for _, r := range a {
		if (r < 'a' || r > 'z') && r != ' ' {
			return false
		}
	}
	joined := strings.Join(strings.Fields(e), "")
	return strings.Join(strings.Fields(a), "") == joined
}

type choiceGrader struct{}

func (g *choiceGrader) grade(expected, answer string) Result {
	e := strings.TrimSpace(expected)
	a := strings.TrimSpace(answer)
	if e == "" || a == "" {
		return emptyResult()
	}
	if strings.EqualFold(e, a) {
		return Result{Correct: true, Score: 1}
	}
	if strings.EqualFold(strings.TrimRight(e, ".!?,;"), strings.TrimRight(a, ".!?,;")) {
		return Result{Correct: true, Score: 1}
	}
	if samePhrase(e, a) {
		return Result{Correct: true, Score: 1}
	}

	// If answer is a single letter or labeled choice, look for that letter in expected
	answerLetter := a
	if m := choicePrefixRe.FindStringSubmatch(a); len(m) > 1 {
		answerLetter = m[1]
	}
	if singleLetterRe.MatchString(answerLetter) && len(strings.TrimSpace(answerLetter)) == 1 {
		letters := singleLetterRe.FindAllString(e, -1)
		for _, l := range letters {
			if strings.EqualFold(l, answerLetter) {
				return Result{Correct: true, Score: 1}
			}
		}
	}

	// If expected is a single letter, look for it in the answer's letters
	if singleLetterRe.MatchString(e) && len(strings.TrimSpace(e)) == 1 {
		letters := singleLetterRe.FindAllString(a, -1)
		for _, l := range letters {
			if strings.EqualFold(l, e) {
				return Result{Correct: true, Score: 1}
			}
		}
	}

	return Result{Correct: false, Score: 0, Feedback: "Incorrect choice"}
}
