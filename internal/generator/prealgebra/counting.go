package prealgebra

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

// rationalCountProblem asks how many of eight listed numbers are rational.
// The pool holds eight rational and eight irrational values, so the count is
// drawn uniformly from zero to eight and every answer appears.
func rationalCountProblem() generator.Problem {
	type value struct {
		text     string
		rational bool
		note     string
	}
	pool := []value{
		{`\sqrt{2}`, false, "the square root of a non-square integer is irrational"},
		{`\pi`, false, "pi has a non-repeating, non-terminating decimal expansion"},
		{`\sqrt{5}`, false, "5 is not a perfect square"},
		{`e`, false, "e is transcendental, so it is irrational"},
		{`\sqrt{7}`, false, "7 is not a perfect square"},
		{`0.101101110\ldots`, false, "the decimal neither terminates nor repeats"},
		{`\sqrt[3]{2}`, false, "the cube root of 2 is irrational"},
		{`\sqrt{11}`, false, "11 is not a perfect square"},
		{`\frac{3}{4}`, true, "a ratio of two integers"},
		{`0.75`, true, "it terminates, so it equals 3/4"},
		{`0.333\ldots`, true, "it repeats, so it equals 1/3"},
		{`-7`, true, "an integer is a ratio over 1"},
		{`\sqrt{9}`, true, "the square root is 3, an integer"},
		{`\frac{5}{2}`, true, "a ratio of two integers"},
		{`0.121212\ldots`, true, "it repeats, so it is rational"},
		{`\sqrt{16}`, true, "the square root is 4, an integer"},
	}
	answer := rand.Intn(9)
	var yes, no []value
	for _, v := range pool {
		if v.rational {
			yes = append(yes, v)
		} else {
			no = append(no, v)
		}
	}
	rand.Shuffle(len(yes), func(i, j int) { yes[i], yes[j] = yes[j], yes[i] })
	rand.Shuffle(len(no), func(i, j int) { no[i], no[j] = no[j], no[i] })
	chosen := append(append([]value{}, yes[:answer]...), no[:8-answer]...)
	rand.Shuffle(len(chosen), func(i, j int) { chosen[i], chosen[j] = chosen[j], chosen[i] })
	parts := make([]string, len(chosen))
	notes := make([]string, len(chosen))
	for i, v := range chosen {
		text := v.text
		plain := true
		for _, r := range text {
			if (r < '0' || r > '9') && r != '.' && r != '-' {
				plain = false
				break
			}
		}
		if plain {
			parts[i] = fmt.Sprintf("(%c) %s", 'a'+i, text)
		} else {
			parts[i] = fmt.Sprintf("(%c) \\(%s\\)", 'a'+i, text)
		}
		verdict := "irrational"
		if v.rational {
			verdict = "rational"
		}
		notes[i] = fmt.Sprintf("%s: %s, %s.", v.text, verdict, v.note)
	}
	var b strings.Builder
	b.WriteString("How many of these numbers are rational? ")
	b.WriteString(strings.Join(parts, " "))
	b.WriteString(" (enter a number)")
	return generator.Problem{
		Question:    b.String(),
		Answer:      fmt.Sprintf("%d", answer),
		Explanation: fmt.Sprintf("Exactly %d of the eight numbers are rational. %s", answer, strings.Join(notes, " ")),
		Facts:       map[string]string{"count": fmt.Sprintf("%d", answer)},
	}
}

// propertyValueProblem asks for a value that follows from one algebraic
// property. Each family keeps the property visible in the reasoning.
func propertyValueProblem(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*8)
	switch rand.Intn(4) {
	case 0:
		k := 2 + rand.Intn(10+scale)
		return generator.Problem{
			Question:    fmt.Sprintf("What number must be added to %d to get 0? (enter a number)", k),
			Answer:      fmt.Sprintf("%d", -k),
			Explanation: fmt.Sprintf("The additive inverse of %d is %d, because %d + (%d) = 0.", k, -k, k, -k),
			Facts:       map[string]string{"k": fmt.Sprintf("%d", k), "answer": fmt.Sprintf("%d", -k)},
		}
	case 1:
		k := 2 + rand.Intn(10+scale)
		return generator.Problem{
			Question:    fmt.Sprintf("Using the multiplicative identity, what is %d times 1? (enter a number)", k),
			Answer:      fmt.Sprintf("%d", k),
			Explanation: fmt.Sprintf("Multiplying by 1 leaves the number unchanged, so %d times 1 is %d.", k, k),
			Facts:       map[string]string{"k": fmt.Sprintf("%d", k), "answer": fmt.Sprintf("%d", k)},
		}
	case 2:
		a := 2 + rand.Intn(6+scale)
		b := 2 + rand.Intn(6+scale)
		c := 2 + rand.Intn(6+scale)
		value := a * (b + c)
		return generator.Problem{
			Question:    fmt.Sprintf("Using the distributive property, what is %d(%d + %d)? (enter a number)", a, b, c),
			Answer:      fmt.Sprintf("%d", value),
			Explanation: fmt.Sprintf("Distribute %d over the sum: %d times %d plus %d times %d equals %d.", a, a, b, a, c, value),
			Facts:       map[string]string{"a": fmt.Sprintf("%d", a), "b": fmt.Sprintf("%d", b), "c": fmt.Sprintf("%d", c), "answer": fmt.Sprintf("%d", value)},
		}
	default:
		a := 2 + rand.Intn(15+scale)
		b := 2 + rand.Intn(15+scale)
		return generator.Problem{
			Question:    fmt.Sprintf("If %d + %d = %d + x, what is x? (enter a number)", a, b, b),
			Answer:      fmt.Sprintf("%d", a),
			Explanation: fmt.Sprintf("Addition is commutative, so %d + %d equals %d + %d and x is %d.", a, b, b, a, a),
			Facts:       map[string]string{"a": fmt.Sprintf("%d", a), "b": fmt.Sprintf("%d", b), "answer": fmt.Sprintf("%d", a)},
		}
	}
}

// betterBuyProblem asks how many cents per ounce separate the two unit
// prices, which is the quantity a shopper compares.
func betterBuyProblem() generator.Problem {
	sizes := []int{8, 12, 16, 20, 24, 32}
	unitPrices := []int{10, 15, 20, 25, 30, 40}
	i := rand.Intn(len(sizes))
	j := rand.Intn(len(sizes))
	for j == i {
		j = rand.Intn(len(sizes))
	}
	k := rand.Intn(len(unitPrices))
	l := rand.Intn(len(unitPrices))
	for l == k {
		l = rand.Intn(len(unitPrices))
	}
	size1, size2 := sizes[i], sizes[j]
	unit1, unit2 := unitPrices[k], unitPrices[l]
	price1 := moneyCents(size1 * unit1)
	price2 := moneyCents(size2 * unit2)
	diff := unit1 - unit2
	if diff < 0 {
		diff = -diff
	}
	better := "Package 1"
	if unit2 < unit1 {
		better = "Package 2"
	}
	return generator.Problem{
		Question: fmt.Sprintf("Package 1: %d oz for %s dollars. Package 2: %d oz for %s dollars. How many cents per ounce cheaper is the better buy? (enter a number)",
			size1, price1, size2, price2),
		Answer: fmt.Sprintf("%d", diff),
		Explanation: fmt.Sprintf("Package 1 costs %s cents per ounce and Package 2 costs %s cents per ounce, so %s is cheaper by %d cents per ounce.",
			moneyCents(unit1), moneyCents(unit2), better, diff),
		Facts: map[string]string{
			"size_1": fmt.Sprintf("%d", size1), "price_1": price1, "unit_1": moneyCents(unit1),
			"size_2": fmt.Sprintf("%d", size2), "price_2": price2, "unit_2": moneyCents(unit2),
			"answer": fmt.Sprintf("%d", diff),
		},
	}
}

// inequalityIntegerProblem counts the integers in a bounded window that
// satisfy the inequality, so the direction and the open or closed endpoint
// both change the answer.
func inequalityIntegerProblem() generator.Problem {
	bound := rand.Intn(7) - 3
	symbols := []string{">", ">=", "<", "<="}
	symbol := symbols[rand.Intn(len(symbols))]
	low, high := -6, 6
	count := 0
	for x := low; x <= high; x++ {
		ok := false
		switch symbol {
		case ">":
			ok = x > bound
		case ">=":
			ok = x >= bound
		case "<":
			ok = x < bound
		case "<=":
			ok = x <= bound
		}
		if ok {
			count++
		}
	}
	strict := symbol == ">" || symbol == "<"
	endpoint := "closed"
	if strict {
		endpoint = "open"
	}
	return generator.Problem{
		Question: fmt.Sprintf("How many integers from %d to %d inclusive satisfy \\(x %s %d\\)? (enter a number)",
			low, high, map[string]string{">": ">", ">=": `\ge`, "<": "<", "<=": `\le`}[symbol], bound),
		Answer: fmt.Sprintf("%d", count),
		Explanation: fmt.Sprintf("The graph uses an %s circle at %d and shades to the %s. Counting the integers from %d to %d that satisfy the inequality gives %d.",
			endpoint, bound, map[bool]string{true: "right", false: "left"}[symbol == ">" || symbol == ">="], low, high, count),
		Facts: map[string]string{
			"bound": fmt.Sprintf("%d", bound), "symbol": symbol, "low": fmt.Sprintf("%d", low),
			"high": fmt.Sprintf("%d", high), "answer": fmt.Sprintf("%d", count),
		},
	}
}
