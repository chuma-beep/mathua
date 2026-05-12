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

type countObjectsGen struct{}

func (g *countObjectsGen) Generate(difficulty float64) generator.Problem {
	max := int(3 + difficulty*7)
	n := rand.Intn(max) + 1
	dots := strings.Repeat("* ", n)
	dots = strings.TrimSpace(dots)
	return generator.Problem{
		Question:    fmt.Sprintf("Count the dots:\n\n%s", dots),
		Answer:      fmt.Sprintf("%d", n),
		Explanation: fmt.Sprintf("There %s %d dot%s.", plural("is", "are", n), n, plural("", "s", n)),
	}
}

type countCardinalityGen struct{}

func (g *countCardinalityGen) Generate(difficulty float64) generator.Problem {
	max := int(3 + difficulty*7)
	n := rand.Intn(max) + 1
	objects := []string{"apple", "star", "cat", "car", "balloon", "dog", "frog", "book", "flower", "note"}
	if n > len(objects) {
		n = len(objects)
	}
	shown := objects[:n]
	return generator.Problem{
		Question:    fmt.Sprintf("I count: %s.\n\nThe last number I say is %d. How many objects are there?", strings.Join(shown, ", "), n),
		Answer:      fmt.Sprintf("%d", n),
		Explanation: fmt.Sprintf("The last number you count (%d) tells you the total.", n),
	}
}

type countNumberLineGen struct{}

func (g *countNumberLineGen) Generate(difficulty float64) generator.Problem {
	max := int(5 + difficulty*15)
	n := rand.Intn(max-1) + 2
	prev := n - 1
	next := n + 1
	hidden := prev
	display := fmt.Sprintf("%d, __, %d", hidden, next)
	if difficulty > 0.5 {
		hidden = next
		display = fmt.Sprintf("%d, __, %d", prev, next)
	}
	return generator.Problem{
		Question:    fmt.Sprintf("What number belongs on the number line?\n\n%s", display),
		Answer:      fmt.Sprintf("%d", n),
		Explanation: fmt.Sprintf("%d comes between %d and %d.", n, prev, next),
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
	return generator.Problem{
		Question:    fmt.Sprintf("Compare: %d __ %d\nEnter > or <", a, b),
		Answer:      answer,
		Explanation: fmt.Sprintf("Since %d vs %d, we write %d %s %d.", a, b, a, answer, b),
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
	display := fmt.Sprintf("%d, _, %d", before, next)
	hidden := current
	if difficulty > 0.5 {
		display = fmt.Sprintf("%d, %d, _, %d", before, current, next+g.step)
		hidden = next
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Skip count by %d:\n\n%s", g.step, display),
		Answer:      fmt.Sprintf("%d", hidden),
		Explanation: fmt.Sprintf("Counting by %ds: %d, %d, %d, %d", g.step, before, current, next, next+g.step),
	}
}

type countObjects20Gen struct{}

func (g *countObjects20Gen) Generate(difficulty float64) generator.Problem {
	n := rand.Intn(10) + 11
	dots := strings.Repeat("* ", n)
	dots = strings.TrimSpace(dots)
	return generator.Problem{
		Question:    fmt.Sprintf("Count the dots:\n\n%s", dots),
		Answer:      fmt.Sprintf("%d", n),
		Explanation: fmt.Sprintf("There are %d dots.", n),
	}
}

type countOrdinalGen struct{}

func (g *countOrdinalGen) Generate(difficulty float64) generator.Problem {
	n := rand.Intn(10)
	ord := ordinals[n]
	line := ""
	for i := 1; i <= 10; i++ {
		marker := "O"
		if i-1 == n {
			marker = "X"
		}
		line += marker + " "
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Which position is the X in?\n\n%s\n\n(Answer: first, second, third, ...)", strings.TrimSpace(line)),
		Answer:      ord,
		Explanation: fmt.Sprintf("The X is in the %s position.", ord),
	}
}

type countBackwardsGen struct{}

func (g *countBackwardsGen) Generate(difficulty float64) generator.Problem {
	start := 5 + int(difficulty*15)
	n := rand.Intn(start-3) + 3
	missing := n - 1
	display := fmt.Sprintf("%d, _, %d", n, missing-1)
	return generator.Problem{
		Question:    fmt.Sprintf("Count backwards! Fill in the missing number:\n\n%s", display),
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
