package counting

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func Register(reg *generator.Registry) {
	reg.Register("count.objects", &countObjectsGen{})
	reg.Register("count.cardinality", &countCardinalityGen{})
	reg.Register("count.number_line", &countNumberLineGen{})
	reg.Register("count.compare", &countCompareGen{})
	reg.Register("count.skip_2", &skipCountGen{step: 2, max: 20})
	reg.Register("count.skip_5", &skipCountGen{step: 5, max: 50})
	reg.Register("count.skip_10", &skipCountGen{step: 10, max: 100})
	reg.Register("count.objects_20", &countObjects20Gen{})
	reg.Register("count.ordinal", &countOrdinalGen{})
	reg.Register("count.backwards", &countBackwardsGen{})
}

var ordinals = []string{
	"first", "second", "third", "fourth", "fifth",
	"sixth", "seventh", "eighth", "ninth", "tenth",
}

var countSymbols = []string{"#", "@", "%", "&", "+", "=", "~", "$", "!", "?"}

type countObjectsGen struct{}

func (g *countObjectsGen) Generate(difficulty float64) generator.Problem {
	max := int(3 + difficulty*7)
	n := rand.Intn(max) + 1

	scenarios := []struct {
		q   string
		fmt string
	}{
		{"How many %s are there in this row?", "%s"},
		{"A ladybug has %s spots on its back. How many spots?", "%s"},
		{"Count the flowers in the garden:\n\n%s", "%s"},
		{"How many stars are in the sky?\n\n%s", "%s"},
		{"Arrange %s marbles in a line. How many marbles?", "%s"},
	}

	s := scenarios[rand.Intn(len(scenarios))]
	sym := countSymbols[rand.Intn(len(countSymbols))]
	item := "stars"
	switch sym {
	case "#":
		item = "hash marks"
	case "@":
		item = "at symbols"
	case "%":
		item = "percent signs"
	case "&":
		item = "ampersands"
	case "+":
		item = "plus signs"
	case "=":
		item = "equals signs"
	case "~":
		item = "tildes"
	case "$":
		item = "dollar signs"
	case "!":
		item = "exclamation marks"
	case "?":
		item = "question marks"
	}

	display := make([]string, n)
	for i := range display {
		display[i] = sym
	}
	line := strings.Join(display, " ")
	q := fmt.Sprintf(s.q, item)
	if strings.Contains(s.fmt, "%s") {
		q = strings.Replace(q, "%s", line, 1)
	}

	return generator.Problem{
		Question:    q,
		Answer:      fmt.Sprintf("%d", n),
		Explanation: fmt.Sprintf("There %s %d %s.", plural("is", "are", n), n, item),
	}
}

type countCardinalityGen struct{}

func (g *countCardinalityGen) Generate(difficulty float64) generator.Problem {
	max := int(3 + difficulty*7)
	n := rand.Intn(max) + 1

	scenarios := []struct {
		items []string
		unit  string
		scene string
	}{
		{[]string{"ant", "bee", "cricket", "moth", "fly", "worm", "snail", "beetle", "spider", "ladybug"}, "bugs", "in the garden"},
		{[]string{"car", "truck", "bus", "bike", "train", "plane", "boat", "scooter", "van", "taxi"}, "vehicles", "in the parking lot"},
		{[]string{"rose", "tulip", "lily", "daisy", "orchid", "violet", "poppy", "lavender", "jasmine", "sunflower"}, "flowers", "in the vase"},
		{[]string{"pen", "book", "ruler", "eraser", "pencil", "crayon", "sharpener", "marker", "notebook", "sticker"}, "school supplies", "on the desk"},
	}

	s := scenarios[rand.Intn(len(scenarios))]
	if n > len(s.items) {
		n = len(s.items)
	}
	shown := s.items[:n]

	return generator.Problem{
		Question:    fmt.Sprintf("%s %s:\n\n%s\n\nThe last one I count is %d. How many %s are there?",
			s.scene, s.unit, strings.Join(shown, ", "), n, s.unit),
		Answer:      fmt.Sprintf("%d", n),
		Explanation: fmt.Sprintf("The last number you count (%d) tells the total number of %s.", n, s.unit),
	}
}

type countNumberLineGen struct{}

func (g *countNumberLineGen) Generate(difficulty float64) generator.Problem {
	max := int(5 + difficulty*15)
	n := rand.Intn(max-1) + 2
	prev := n - 1
	next := n + 1

	scenarios := []string{
		"Fill in the missing number:",
		"What number belongs here?",
		"Find the missing value on the line:",
	}

	hidden := prev
	display := fmt.Sprintf("... %d, __, %d ...", hidden, next)
	if difficulty > 0.5 {
		hidden = next
		display = fmt.Sprintf("... %d, __, %d ...", prev, next)
	}
	if difficulty > 0.75 {
		hidden = n
		display = fmt.Sprintf("... %d, %d, __ ...", prev, next)
	}

	q := fmt.Sprintf("%s\n\n%s", scenarios[rand.Intn(len(scenarios))], display)
	return generator.Problem{
		Question:    q,
		Answer:      fmt.Sprintf("%d", n),
		Explanation: fmt.Sprintf("The number %d comes between %d and %d on the number line.", n, prev, next),
	}
}

type countCompareGen struct{}

func (g *countCompareGen) Generate(difficulty float64) generator.Problem {
	max := int(5 + difficulty*15)
	a := rand.Intn(max) + 1
	b := rand.Intn(max) + 1
	if a == b {
		b += 1
	}

	answer := ">"
	if b > a {
		answer = "<"
	}

	scenarios := []struct {
		q string
	}{
		{fmt.Sprintf("Compare: %d __ %d\n\nWhich symbol makes this true? Enter > or <.", a, b)},
		{fmt.Sprintf("Is %d greater than or less than %d?\n\nEnter > or <.", a, b)},
		{fmt.Sprintf("Choose the correct symbol: %d __ %d\n\nEnter > if %d is greater, < if %d is less.", a, b, a, a)},
	}

	s := scenarios[rand.Intn(len(scenarios))]
	return generator.Problem{
		Question:    s.q,
		Answer:      answer,
		Explanation: fmt.Sprintf("%d %s %d because %d comes %s %d on the number line.", a, answer, b, a, map[string]string{">": "after", "<": "before"}[answer], b),
	}
}

type skipCountGen struct {
	step int
	max  int
}

func (g *skipCountGen) Generate(difficulty float64) generator.Problem {
	steps := g.max / g.step
	pos := rand.Intn(steps - 2)
	current := (pos + 1) * g.step
	next := current + g.step
	before := current - g.step

	var display string
	var hidden int
	var q string

	scenarios := []string{
		"Fill in the missing number:",
		"What comes next?",
		"Complete the pattern:",
	}

	if difficulty < 0.3 {
		hidden = current
		display = fmt.Sprintf("%d, _, %d", before, next)
	} else if difficulty < 0.6 {
		hidden = next
		display = fmt.Sprintf("%d, %d, _, %d", before, current, next+g.step)
	} else {
		hidden = before
		display = fmt.Sprintf("_, %d, %d, %d", current, next, next+g.step)
	}

	contexts := []string{
		fmt.Sprintf("Counting by %ds", g.step),
		"Skip counting",
		fmt.Sprintf("Every %dth number", g.step),
	}
	q = fmt.Sprintf("%s %s\n\n%s", contexts[rand.Intn(len(contexts))], scenarios[rand.Intn(len(scenarios))], display)

	return generator.Problem{
		Question:    q,
		Answer:      fmt.Sprintf("%d", hidden),
		Explanation: fmt.Sprintf("Counting by %ds: %d, %d, %d, %d", g.step, before, current, next, next+g.step),
	}
}

type countObjects20Gen struct{}

func (g *countObjects20Gen) Generate(difficulty float64) generator.Problem {
	n := rand.Intn(10) + 11

	scenarios := []struct {
		items []string
		unit  string
	}{
		{[]string{"■", "□", "▲", "△", "○", "●", "◆", "◇", "★", "☆", "♠", "♣", "♥", "♦", "☼", "♪", "♫", "☺", "☻", "♀"}, "shapes"},
		{[]string{"red marble", "blue marble", "green marble", "yellow marble", "purple marble", "orange marble", "pink marble", "white marble", "black marble", "gold marble"}, "marbles in a bag"},
	}

	s := scenarios[rand.Intn(len(scenarios))]
	rows := (n + 4) / 5
	var grid []string
	idx := 0
	for r := 0; r < rows; r++ {
		var row []string
		for c := 0; c < 5 && idx < n; c++ {
			sym := s.items[rand.Intn(len(s.items))]
			row = append(row, sym)
			idx++
		}
		grid = append(grid, strings.Join(row, " "))
	}
	layout := strings.Join(grid, "\n")

	questions := []string{
		"Count the items arranged below:\n\n%s\n\nHow many items are there?",
		"This grid shows a collection of items:\n\n%s\n\nEnter the total count.",
		"How many items are in this arrangement?\n\n%s",
	}

	q := fmt.Sprintf(questions[rand.Intn(len(questions))], layout)
	return generator.Problem{
		Question:    q,
		Answer:      fmt.Sprintf("%d", n),
		Explanation: fmt.Sprintf("There are %d items in total.", n),
	}
}

type countOrdinalGen struct{}

func (g *countOrdinalGen) Generate(difficulty float64) generator.Problem {
	n := rand.Intn(10)
	ord := ordinals[n]

	scenarios := []struct {
		name  string
		items []string
	}{
		{"cars in a race: 🏎️", []string{"🏎️", "🚗", "🚙", "🚕", "🚓", "🚑", "🚒", "🚐", "🚜", "🏍️"}},
		{"runners in a race", []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}},
		{"positions on a ladder", []string{"1st", "2nd", "3rd", "4th", "5th", "6th", "7th", "8th", "9th", "10th"}},
	}

	s := scenarios[rand.Intn(len(scenarios))]
	var line strings.Builder
	for i := 1; i <= 10; i++ {
		marker := s.items[i-1]
		if i-1 == n {
			line.WriteString("[")
			line.WriteString(marker)
			line.WriteString("]")
		} else {
			line.WriteString(" ")
			line.WriteString(marker)
			line.WriteString(" ")
		}
		line.WriteString(" ")
	}

	return generator.Problem{
		Question:    fmt.Sprintf("These are the %s:\n\n%s\n\nWhich position does the [%s] hold? (Answer: first, second, third, ...)", s.name, strings.TrimSpace(line.String()), s.items[n]),
		Answer:      ord,
		Explanation: fmt.Sprintf("The [%s] is in the %s position.", s.items[n], ord),
	}
}

type countBackwardsGen struct{}

func (g *countBackwardsGen) Generate(difficulty float64) generator.Problem {
	start := 5 + int(difficulty*15)
	n := rand.Intn(start-3) + 3
	missing := n - 1

	scenarios := []string{
		"Count backwards! Fill in the missing number:",
		"Reverse count — what goes here?",
		"Counting down:",
	}

	display := fmt.Sprintf("%d, _, %d", n, missing-1)
	if difficulty > 0.5 {
		display = fmt.Sprintf("%d, %d, _, %d", n, missing, missing-2)
		missing = missing - 1
	}

	return generator.Problem{
		Question:    fmt.Sprintf("%s\n\n%s", scenarios[rand.Intn(len(scenarios))], display),
		Answer:      fmt.Sprintf("%d", missing),
		Explanation: fmt.Sprintf("Backwards from %d: %d, %d, %d ...", n, n, missing, missing-1),
	}
}

func plural(singular, plural string, n int) string {
	if n == 1 {
		return singular
	}
	return plural
}
