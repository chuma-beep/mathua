package topology

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

// topoVerdict is one candidate in a counting question. yes records whether the
// candidate has the property the question asks about, and note says why in one
// short clause. Every pool holds eight yes cases and eight no cases, so drawing
// eight candidates gives a count anywhere from zero to eight and the answer
// stays a small integer while the stem changes with every draw.
type topoVerdict struct {
	label string
	yes   bool
	note  string
}

const topoDraw = 8

// countTopoVerdicts draws eight distinct candidates from pool and asks how many
// have the property. The answer is the count, so the question is production
// rather than recognition; the explanation still names every verdict.
func countTopoVerdicts(lead string, pool []topoVerdict) generator.Problem {
	var yes, no []topoVerdict
	for _, v := range pool {
		if v.yes {
			yes = append(yes, v)
		} else {
			no = append(no, v)
		}
	}
	if len(yes) > topoDraw {
		yes = yes[:topoDraw]
	}
	if len(no) > topoDraw {
		no = no[:topoDraw]
	}
	// Choose the answer first, then fill the list to match it. Sampling eight
	// candidates and counting would reach 0 and 8 only when every draw lands on
	// the same side, which is too rare to keep all nine answers alive.
	count := rand.Intn(topoDraw + 1)
	if count > len(yes) {
		count = len(yes)
	}
	if topoDraw-count > len(no) {
		count = topoDraw - len(no)
	}
	rand.Shuffle(len(yes), func(i, j int) { yes[i], yes[j] = yes[j], yes[i] })
	rand.Shuffle(len(no), func(i, j int) { no[i], no[j] = no[j], no[i] })
	chosen := append([]topoVerdict{}, yes[:count]...)
	chosen = append(chosen, no[:topoDraw-count]...)
	rand.Shuffle(len(chosen), func(i, j int) { chosen[i], chosen[j] = chosen[j], chosen[i] })

	parts := make([]string, 0, topoDraw)
	notes := make([]string, 0, topoDraw)
	for i, v := range chosen {
		parts = append(parts, fmt.Sprintf("(%s) %s", string(rune('a'+i)), v.label))
		verdict := "no"
		if v.yes {
			verdict = "yes"
		}
		notes = append(notes, fmt.Sprintf("%s: %s, %s", v.label, verdict, v.note))
	}
	return generator.Problem{
		Question:    fmt.Sprintf("%s %s (enter a number)", lead, strings.Join(parts, " ")),
		Answer:      strconv.Itoa(count),
		Explanation: strings.Join(notes, "; ") + ".",
		Facts:       map[string]string{"count": strconv.Itoa(count)},
	}
}

func openClosedProblem() generator.Problem {
	a := rand.Intn(20) - 10
	b := a + 1 + rand.Intn(8)
	c := b + 1 + rand.Intn(5)
	d := c + 1 + rand.Intn(5)
	kind := "open"
	if rand.Intn(2) == 0 {
		kind = "closed"
	}
	pick := func(open, closed bool) bool {
		if kind == "open" {
			return open
		}
		return closed
	}
	pool := []topoVerdict{
		{fmt.Sprintf("(%d, %d)", a, b), pick(true, false), "an open interval is open and not closed"},
		{fmt.Sprintf("(-inf, %d)", b), pick(true, false), "an open ray is open and not closed"},
		{fmt.Sprintf("R \\ {%d}", a), pick(true, false), "removing a point leaves an open set that is not closed"},
		{fmt.Sprintf("(%d, %d) union (%d, %d)", a, b, c, d), pick(true, false), "a union of open intervals is open and not closed"},
		{fmt.Sprintf("(%d, inf)", b), pick(true, false), "an open ray is open and not closed"},
		{fmt.Sprintf("(-inf, %d) union (%d, inf)", a, b), pick(true, false), "a union of open rays is open and not closed"},
		{fmt.Sprintf("[%d, %d]", a, b), pick(false, true), "a closed interval is closed and not open"},
		{fmt.Sprintf("[%d, inf)", a), pick(false, true), "a closed ray is closed and not open"},
		{fmt.Sprintf("{%d, %d}", a, b), pick(false, true), "a finite set is closed and not open"},
		{fmt.Sprintf("[%d, %d] union [%d, %d]", a, b, c, d), pick(false, true), "a finite union of closed intervals is closed and not open"},
		{"N", pick(false, true), "the natural numbers form a closed set that is not open"},
		{fmt.Sprintf("{%d, %d, %d}", a, b, c), pick(false, true), "a finite set is closed and not open"},
		{"R", pick(true, true), "R is both open and closed"},
		{"the empty set", pick(true, true), "the empty set is both open and closed"},
		{fmt.Sprintf("(%d, %d]", a, b), pick(false, false), "a half-open interval is neither open nor closed"},
		{fmt.Sprintf("[%d, %d)", a, b), pick(false, false), "a half-open interval is neither open nor closed"},
		{"Q", pick(false, false), "the rationals have empty interior and closure R"},
		{"R minus Q", pick(false, false), "the irrationals have empty interior and closure R"},
	}
	return countTopoVerdicts(fmt.Sprintf("How many of these subsets of R are %s?", kind), pool)
}

func compactProblem() generator.Problem {
	a := rand.Intn(12) - 6
	b := a + 1 + rand.Intn(6)
	c := b + 1 + rand.Intn(4)
	d := c + 1 + rand.Intn(4)
	pool := []topoVerdict{
		{fmt.Sprintf("[%d, %d]", a, b), true, "a closed bounded interval is compact"},
		{fmt.Sprintf("{%d, %d}", a, b), true, "a finite set is compact"},
		{fmt.Sprintf("[%d, %d] union [%d, %d]", a, b, c, d), true, "a finite union of compact sets is compact"},
		{"the empty set", true, "the empty set is compact"},
		{fmt.Sprintf("{%d}", a), true, "a one-point set is compact"},
		{fmt.Sprintf("[%d, %d] union {%d}", a, b, c), true, "a finite union of compact sets is compact"},
		{"{0} union {1/n : n >= 1}", true, "the set is closed and bounded"},
		{fmt.Sprintf("[%d, %d] intersect Z", a, b), true, "a finite set is compact"},
		{fmt.Sprintf("(%d, %d)", a, b), false, "the open interval is bounded but not closed"},
		{fmt.Sprintf("[%d, inf)", a), false, "the ray is closed but not bounded"},
		{fmt.Sprintf("(%d, %d]", a, b), false, "the half-open interval is not closed"},
		{"Z", false, "the integers are closed but not bounded"},
		{"R", false, "R is closed but not bounded"},
		{"Q", false, "the rationals are not closed in R"},
		{"R minus Q", false, "the irrationals are not closed in R"},
		{fmt.Sprintf("(%d, inf) intersect Q", a), false, "the set is not closed in R"},
	}
	return countTopoVerdicts("How many of these subsets of R are compact?", pool)
}

func connectedProblem() generator.Problem {
	a := rand.Intn(12) - 6
	b := a + 1 + rand.Intn(6)
	c := b + 1 + rand.Intn(4)
	d := c + 1 + rand.Intn(4)
	pool := []topoVerdict{
		{fmt.Sprintf("(%d, %d)", a, b), true, "an interval is connected"},
		{fmt.Sprintf("{%d}", a), true, "a one-point set is connected"},
		{fmt.Sprintf("[%d, %d]", a, b), true, "an interval is connected"},
		{"R", true, "R is connected"},
		{fmt.Sprintf("[%d, %d] union [%d, %d]", a, b, b, c), true, "the intervals overlap, so the union is connected"},
		{fmt.Sprintf("(%d, inf)", a), true, "a ray is an interval and is connected"},
		{fmt.Sprintf("(-inf, %d)", b), true, "a ray is an interval and is connected"},
		{fmt.Sprintf("[%d, inf)", a), true, "a ray is an interval and is connected"},
		{fmt.Sprintf("(%d, %d) union (%d, %d)", a, b, c, d), false, "the two intervals are separated"},
		{fmt.Sprintf("{%d, %d}", a, b), false, "a two-point set is disconnected"},
		{"Q", false, "the rationals are totally disconnected"},
		{"Z", false, "the integers are totally disconnected"},
		{"R minus Q", false, "the irrationals are totally disconnected"},
		{"the Cantor set", false, "the Cantor set is totally disconnected"},
		{fmt.Sprintf("{%d, %d, %d}", a, b, c), false, "a three-point set is disconnected"},
		{fmt.Sprintf("[%d, %d] union [%d, %d]", a, b, d, d+1), false, "the two closed intervals are separated"},
	}
	return countTopoVerdicts("How many of these subsets of R are connected?", pool)
}

func continuousProblem() generator.Problem {
	a := rand.Intn(12) - 6
	k := 1 + rand.Intn(5)
	m := 1 + rand.Intn(5)
	pool := []topoVerdict{
		{fmt.Sprintf("x^2 + %dx + %d", k, m), true, "a polynomial is continuous"},
		{fmt.Sprintf("|x - %d|", a), true, "the absolute value of a continuous function is continuous"},
		{fmt.Sprintf("1/(x^2 + %d)", k), true, "the denominator never vanishes"},
		{"sin x", true, "the sine function is continuous"},
		{"e^x", true, "the exponential is continuous"},
		{"cos x", true, "the cosine function is continuous"},
		{"x^3", true, "a polynomial is continuous"},
		{fmt.Sprintf("x + 1/(x^2 + %d)", k), true, "a sum of continuous functions is continuous"},
		{"floor(x)", false, "the floor function jumps at every integer"},
		{fmt.Sprintf("1/(x - %d)", a), false, "the function has a pole at x = a"},
		{"1/x", false, "the function has a pole at x = 0"},
		{fmt.Sprintf("sqrt(x - %d)", a), false, "the function is not defined for x < a"},
		{"tan x", false, "the tangent has poles at odd multiples of pi/2"},
		{fmt.Sprintf("1/(x - %d)^2", a), false, "the function has a pole at x = a"},
		{"the step function", false, "the step function jumps at the origin"},
		{fmt.Sprintf("ln|x - %d|", a), false, "the logarithm has a pole at x = a"},
	}
	return countTopoVerdicts("How many of these functions are continuous on R?", pool)
}
