package fractions

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

type fracVerdict struct {
	label string
	yes   bool
	note  string
}

const fracDraw = 8

func countFracVerdicts(lead string, pool []fracVerdict) generator.Problem {
	var yes, no []fracVerdict
	for _, v := range pool {
		if v.yes {
			yes = append(yes, v)
		} else {
			no = append(no, v)
		}
	}
	var count int
	var chosen []fracVerdict
	var question string
	for attempt := 0; attempt < 40; attempt++ {
		count = rand.Intn(fracDraw + 1)
		if count > len(yes) {
			count = len(yes)
		}
		if fracDraw-count > len(no) {
			count = fracDraw - len(no)
		}
		rand.Shuffle(len(yes), func(i, j int) { yes[i], yes[j] = yes[j], yes[i] })
		rand.Shuffle(len(no), func(i, j int) { no[i], no[j] = no[j], no[i] })
		chosen = append(append([]fracVerdict{}, yes[:count]...), no[:fracDraw-count]...)
		rand.Shuffle(len(chosen), func(i, j int) { chosen[i], chosen[j] = chosen[j], chosen[i] })
		parts := make([]string, len(chosen))
		for i, v := range chosen {
			parts[i] = fmt.Sprintf("(%c) %s", 'a'+i, v.label)
		}
		question = fmt.Sprintf("%s %s (enter a number)", lead, strings.Join(parts, " "))
		if !strings.Contains(question, fmt.Sprint(count)) {
			break
		}
	}
	notes := make([]string, len(chosen))
	for i, v := range chosen {
		verdict := "false"
		if v.yes {
			verdict = "true"
		}
		notes[i] = fmt.Sprintf("%s: %s. %s", v.label, verdict, v.note)
	}
	return generator.Problem{
		Question:    question,
		Answer:      fmt.Sprint(count),
		Explanation: fmt.Sprintf("Exactly %d of the %d comparisons are correct. %s", count, fracDraw, strings.Join(notes, " ")),
		Facts:       map[string]string{"count": fmt.Sprint(count)},
	}
}

func benchmarkProblem() generator.Problem {
	lead := "How many of these fractions are less than 1/2?"
	less := true
	if rand.Intn(2) == 1 {
		lead = "How many of these fractions are greater than 1/2?"
		less = false
	}
	low := []fracVerdict{
		{label: "1/3", note: "1/3 is about 0.33, below 0.5."},
		{label: "2/5", note: "2/5 is 0.4, below 0.5."},
		{label: "3/7", note: "3/7 is about 0.43, below 0.5."},
		{label: "4/9", note: "4/9 is about 0.44, below 0.5."},
		{label: "5/11", note: "5/11 is about 0.45, below 0.5."},
		{label: "3/8", note: "3/8 is 0.375, below 0.5."},
		{label: "2/7", note: "2/7 is about 0.29, below 0.5."},
		{label: "4/10", note: "4/10 is 0.4, below 0.5."},
	}
	high := []fracVerdict{
		{label: "2/3", note: "2/3 is about 0.67, above 0.5."},
		{label: "3/5", note: "3/5 is 0.6, above 0.5."},
		{label: "4/7", note: "4/7 is about 0.57, above 0.5."},
		{label: "5/9", note: "5/9 is about 0.56, above 0.5."},
		{label: "6/11", note: "6/11 is about 0.55, above 0.5."},
		{label: "5/8", note: "5/8 is 0.625, above 0.5."},
		{label: "5/7", note: "5/7 is about 0.71, above 0.5."},
		{label: "6/10", note: "6/10 is 0.6, above 0.5."},
	}
	pool := make([]fracVerdict, 0, 16)
	if less {
		for _, v := range low {
			v.yes = true
			pool = append(pool, v)
		}
		for _, v := range high {
			v.yes = false
			pool = append(pool, v)
		}
	} else {
		for _, v := range high {
			v.yes = true
			pool = append(pool, v)
		}
		for _, v := range low {
			v.yes = false
			pool = append(pool, v)
		}
	}
	return countFracVerdicts(lead, pool)
}

func compareWordProblem() generator.Problem {
	// Eight statements are true because Alice ate more, and eight are false
	// because Bob did. The helper draws the answer before choosing, so all nine
	// counts appear.
	truePairs := [][4]int{
		{3, 4, 1, 2},
		{2, 3, 1, 3},
		{5, 8, 1, 2},
		{3, 5, 1, 4},
		{7, 10, 2, 5},
		{4, 5, 1, 3},
		{5, 6, 2, 3},
		{7, 8, 3, 4},
	}
	falsePairs := [][4]int{
		{1, 3, 2, 3},
		{1, 2, 3, 4},
		{2, 5, 3, 5},
		{1, 4, 5, 8},
		{2, 7, 3, 7},
		{3, 10, 4, 10},
		{1, 5, 3, 8},
		{2, 9, 5, 9},
	}
	pool := make([]fracVerdict, 0, 16)
	for _, p := range truePairs {
		pool = append(pool, fracVerdict{
			label: fmt.Sprintf("Alice %d/%d, Bob %d/%d: Alice ate more.", p[0], p[1], p[2], p[3]),
			yes:   true,
			note:  fmt.Sprintf("%d/%d is greater than %d/%d.", p[0], p[1], p[2], p[3]),
		})
	}
	for _, p := range falsePairs {
		pool = append(pool, fracVerdict{
			label: fmt.Sprintf("Alice %d/%d, Bob %d/%d: Alice ate more.", p[0], p[1], p[2], p[3]),
			yes:   false,
			note:  fmt.Sprintf("%d/%d is less than %d/%d, so Bob ate more.", p[0], p[1], p[2], p[3]),
		})
	}
	return countFracVerdicts("How many of these comparisons are correct?", pool)
}
