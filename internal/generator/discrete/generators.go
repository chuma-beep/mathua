package discrete

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func Register(reg *generator.Registry) {
	reg.Register("discrete.logic.propositions", &propositionsGen{})
	reg.Register("discrete.logic.connectives", &connectivesGen{})
	reg.Register("discrete.logic.truth_tables", &truthTablesGen{})
	reg.Register("discrete.logic.quantifiers", &quantifiersGen{})
	reg.Register("discrete.sets.operations", &setOpsGen{})
	reg.Register("discrete.sets.venn", &vennGen{})
	reg.Register("discrete.combinatorics.permutations", &permutationsGen{})
	reg.Register("discrete.combinatorics.combinations", &combinationsGen{})
	reg.Register("discrete.combinatorics.pascal", &pascalGen{})
	reg.Register("discrete.graphs.basics", &graphBasicsGen{})
	reg.Register("discrete.graphs.paths", &graphPathsGen{})
	reg.Register("discrete.graphs.trees", &treesGen{})
	reg.Register("discrete.sequences.recurrence", &recurrenceGen{})
	reg.Register("discrete.proof.induction", &inductionGen{})

	reg.Register("discrete.combinatorics.binomial_theorem", &binomialTheoremGen{})
}

// ----- 1. propositions -----

type propositionsGen struct{}

func (g *propositionsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		statement string
		isProp    string
		reason    string
	}
	table := []entry{
		{"2 + 2 = 4", "yes", "It is a declarative statement with a definite truth value (true)."},
		{"The sky is blue.", "yes", "It is a declarative statement with a definite truth value (true)."},
		{"x is large.", "no", "It is ambiguous — 'large' is not well-defined without context."},
		{"5 is an even number.", "yes", "It is a declarative statement with a definite truth value (false)."},
		{"What time is it?", "no", "It is a question, not a declarative statement."},
		{"Please sit down.", "no", "It is a command, not a declarative statement."},
		{"1 + 1 = 3", "yes", "It is a declarative statement with a definite truth value (false)."},
		{"This sentence is false.", "no", "It is a paradox — it cannot have a consistent truth value."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{
		Question:    fmt.Sprintf("Is '%s' a proposition? (yes/no)", e.statement),
		Answer:      e.isProp,
		Explanation: e.reason,
	}
}

// ----- 2. connectives -----

type connectivesGen struct{}

func (g *connectivesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type assignment struct {
		p, q  bool
		label string
	}
	assignments := []assignment{
		{true, true, "p=True, q=True"},
		{true, false, "p=True, q=False"},
		{false, true, "p=False, q=True"},
		{false, false, "p=False, q=False"},
	}
	a := assignments[rand.Intn(len(assignments))]

	type conn struct {
		name string
		fn   func(p, q bool) bool
	}
	conns := []conn{
		{"AND", func(p, q bool) bool { return p && q }},
		{"OR", func(p, q bool) bool { return p || q }},
		{"NAND", func(p, q bool) bool { return !(p && q) }},
		{"NOR", func(p, q bool) bool { return !(p || q) }},
		{"XOR", func(p, q bool) bool { return p != q }},
	}
	c := conns[rand.Intn(len(conns))]
	result := c.fn(a.p, a.q)

	// Also sometimes ask about NOT
	if rand.Intn(3) == 0 {
		val := rand.Intn(2) == 0
		label := "p=True"
		if !val {
			label = "p=False"
		}
		return generator.Problem{
			Question:    fmt.Sprintf("Given %s, what is NOT p? (true/false)", label),
			Answer:      boolStr(!val),
			Explanation: fmt.Sprintf("NOT %s = %s", boolStr(val), boolStr(!val)),
		}
	}

	return generator.Problem{
		Question:    fmt.Sprintf("Given %s, what is p %s q? (true/false)", a.label, c.name),
		Answer:      boolStr(result),
		Explanation: fmt.Sprintf("p=%s, q=%s => p %s q = %s", boolStr(a.p), boolStr(a.q), c.name, boolStr(result)),
	}
}

func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// ----- 3. truth_tables -----

type truthTablesGen struct{}

func (g *truthTablesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type expr struct {
		name        string
		description string
	}
	exprs := []expr{
		{"p AND q", "\\(p \\land q\\)"},
		{"p OR q", "\\(p \\lor q\\)"},
		{"NOT p", "\\(\\lnot p\\)"},
		{"p XOR q", "\\(p \\oplus q\\)"},
	}
	e := exprs[rand.Intn(len(exprs))]

	vars := 1
	if e.name != "NOT p" {
		vars = 2
	}

	if rand.Intn(2) == 0 {
		rows := 1 << vars
		return generator.Problem{
			Question:    			fmt.Sprintf("How many rows are in a truth table for \\(%s\\)?", e.description),
			Answer:      fmt.Sprintf("%d", rows),
			Explanation: fmt.Sprintf("With %d variable(s), there are 2^%d = %d rows.", vars, vars, rows),
		}
	}

	if vars == 1 {
		val := rand.Intn(2) == 0
		result := !val
		return generator.Problem{
			Question:    fmt.Sprintf("Given p=%s, what is NOT p? (true/false)", boolStr(val)),
			Answer:      boolStr(result),
			Explanation: fmt.Sprintf("NOT %s = %s", boolStr(val), boolStr(result)),
		}
	}

	p := rand.Intn(2) == 0
	q := rand.Intn(2) == 0
	var result bool
	switch e.name {
	case "p AND q":
		result = p && q
	case "p OR q":
		result = p || q
	case "p XOR q":
		result = p != q
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Given p=%s, q=%s, what is %s? (true/false)", boolStr(p), boolStr(q), e.description),
		Answer:      boolStr(result),
		Explanation: fmt.Sprintf("p=%s, q=%s => %s = %s", boolStr(p), boolStr(q), e.description, boolStr(result)),
	}
}

// ----- 4. quantifiers -----

type quantifiersGen struct{}

func (g *quantifiersGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	sets := [][]int{
		{1, 2, 3},
		{2, 4, 6},
		{1, 3, 5, 7},
		{-2, -1, 0, 1, 2},
		{3, 6, 9, 12},
	}
	s := sets[rand.Intn(len(sets))]
	setStr := formatSet(s)

	type qEntry struct {
		cond     func(x int) bool
		condDesc string
		isForAll bool
		expected string
		reason   string
	}

	entries := []qEntry{
		{func(x int) bool { return x > 0 }, "> 0", true, "true", "All elements are > 0."},
		{func(x int) bool { return x > 0 }, "> 0", false, "true", "At least one element is > 0."},
		{func(x int) bool { return x > 5 }, "> 5", true, "false", "Not all elements are > 5."},
		{func(x int) bool { return x > 5 }, "> 5", false, "",
			""},
		{func(x int) bool { return x < 10 }, "< 10", true, "",
			""},
		{func(x int) bool { return x == 0 }, "= 0", false, "",
			""},
		{func(x int) bool { return x%2 == 0 }, "is even", true, "",
			""},
		{func(x int) bool { return x%2 == 0 }, "is even", false, "",
			""},
	}

	// Compute correct answers for entries with empty expected/reason
	for i, e := range entries {
		switch {
		case e.condDesc == "> 5" && !e.isForAll:
			allFalse := true
			for _, v := range s {
				if e.cond(v) {
					allFalse = false
					break
				}
			}
			if allFalse {
				entries[i].expected = "false"
				entries[i].reason = "No element is > 5."
			} else {
				entries[i].expected = "true"
				entries[i].reason = "At least one element is > 5."
			}
		case e.condDesc == "< 10" && e.isForAll:
			allOk := true
			for _, v := range s {
				if !e.cond(v) {
					allOk = false
					break
				}
			}
			if allOk {
				entries[i].expected = "true"
				entries[i].reason = "All elements are < 10."
			} else {
				entries[i].expected = "false"
				entries[i].reason = "Not all elements are < 10."
			}
		case e.condDesc == "= 0" && !e.isForAll:
			hasZero := false
			for _, v := range s {
				if e.cond(v) {
					hasZero = true
					break
				}
			}
			if hasZero {
				entries[i].expected = "true"
				entries[i].reason = "0 is in the set."
			} else {
				entries[i].expected = "false"
				entries[i].reason = "No element equals 0."
			}
		case e.condDesc == "is even" && e.isForAll:
			allEven := true
			for _, v := range s {
				if !e.cond(v) {
					allEven = false
					break
				}
			}
			if allEven {
				entries[i].expected = "true"
				entries[i].reason = "All elements are even."
			} else {
				entries[i].expected = "false"
				entries[i].reason = "Not all elements are even."
			}
		case e.condDesc == "is even" && !e.isForAll:
			hasEven := false
			for _, v := range s {
				if e.cond(v) {
					hasEven = true
					break
				}
			}
			if hasEven {
				entries[i].expected = "true"
				entries[i].reason = "At least one element is even."
			} else {
				entries[i].expected = "false"
				entries[i].reason = "No element is even."
			}
		}
	}

	e := entries[rand.Intn(len(entries))]
	quant := "\\forall"
	if !e.isForAll {
		quant = "\\exists"
	}

	return generator.Problem{
		Question:    fmt.Sprintf("Is \\(%s x \\in %s: x %s\\) true? (true/false)", quant, setStr, e.condDesc),
		Answer:      e.expected,
		Explanation: fmt.Sprintf("%sx ∈ %s: x %s is %s because %s", quant, setStr, e.condDesc, e.expected, e.reason),
	}
}

// ----- 5. set_operations -----

type setOpsGen struct{}

func (g *setOpsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	allElements := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	type pair struct {
		a, b []int
	}
	pairs := []pair{
		{[]int{1, 2, 3}, []int{3, 4, 5}},
		{[]int{2, 4, 6, 8}, []int{1, 2, 3, 4}},
		{[]int{1, 3, 5}, []int{2, 4, 6}},
		{[]int{1, 2, 3, 4, 5}, []int{4, 5, 6, 7}},
		{[]int{1, 4, 7, 10}, []int{2, 5, 8}},
	}
	p := pairs[rand.Intn(len(pairs))]

	aStr := formatSet(p.a)
	bStr := formatSet(p.b)

	ops := []struct {
		name   string
		desc   string
		answer string
		exp    string
	}{
		{
			"union", fmt.Sprintf("\\(A \\cup B\\) where \\(A=%s\\), \\(B=%s\\)", aStr, bStr),
			formatSet(union(p.a, p.b)),
			fmt.Sprintf("A ∪ B = %s (all elements in A or B)", formatSet(union(p.a, p.b))),
		},
		{
			"intersection", fmt.Sprintf("\\(A \\cap B\\) where \\(A=%s\\), \\(B=%s\\)", aStr, bStr),
			formatSet(intersection(p.a, p.b)),
			fmt.Sprintf("A ∩ B = %s (elements in both A and B)", formatSet(intersection(p.a, p.b))),
		},
		{
			"complement", fmt.Sprintf("the complement of \\(A=%s\\) in \\(U=%s\\)", aStr, formatSet(allElements)),
			formatSet(complement(p.a, allElements)),
			fmt.Sprintf("A' = %s (elements in U but not in A)", formatSet(complement(p.a, allElements))),
		},
		{
			"difference", fmt.Sprintf("\\(A - B\\) where \\(A=%s\\), \\(B=%s\\)", aStr, bStr),
			formatSet(difference(p.a, p.b)),
			fmt.Sprintf("A − B = %s (elements in A but not in B)", formatSet(difference(p.a, p.b))),
		},
	}

	op := ops[rand.Intn(len(ops))]

	return generator.Problem{
		Question:    fmt.Sprintf("What is %s?", op.desc),
		Answer:      op.answer,
		Explanation: op.exp,
	}
}

func union(a, b []int) []int {
	m := make(map[int]bool)
	for _, v := range a {
		m[v] = true
	}
	for _, v := range b {
		m[v] = true
	}
	r := make([]int, 0, len(m))
	for v := range m {
		r = append(r, v)
	}
	sortInts(r)
	return r
}

func intersection(a, b []int) []int {
	mb := make(map[int]bool)
	for _, v := range b {
		mb[v] = true
	}
	var r []int
	for _, v := range a {
		if mb[v] {
			r = append(r, v)
		}
	}
	sortInts(r)
	return r
}

func complement(a, u []int) []int {
	ma := make(map[int]bool)
	for _, v := range a {
		ma[v] = true
	}
	var r []int
	for _, v := range u {
		if !ma[v] {
			r = append(r, v)
		}
	}
	sortInts(r)
	return r
}

func difference(a, b []int) []int {
	mb := make(map[int]bool)
	for _, v := range b {
		mb[v] = true
	}
	var r []int
	for _, v := range a {
		if !mb[v] {
			r = append(r, v)
		}
	}
	sortInts(r)
	return r
}

func sortInts(s []int) {
	for i := 0; i < len(s); i++ {
		for j := i + 1; j < len(s); j++ {
			if s[i] > s[j] {
				s[i], s[j] = s[j], s[i]
			}
		}
	}
}

func formatSet(s []int) string {
	if len(s) == 0 {
		return "∅"
	}
	parts := make([]string, len(s))
	for i, v := range s {
		parts[i] = strconv.Itoa(v)
	}
	return "{" + strings.Join(parts, ",") + "}"
}

// ----- 6. venn -----

type vennGen struct{}

func (g *vennGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type scenario struct {
		a, b, both, universe int
	}
	scenarios := []scenario{
		{10, 15, 5, 30},
		{8, 12, 3, 25},
		{6, 9, 4, 20},
		{15, 10, 6, 30},
		{12, 8, 2, 25},
	}
	s := scenarios[rand.Intn(len(scenarios))]

	unionSize := s.a + s.b - s.both
	compUnion := s.universe - unionSize
	onlyA := s.a - s.both
	onlyB := s.b - s.both

	type qType struct {
		question string
		answer   int
		exp      string
	}
	questions := []qType{
		{
			fmt.Sprintf("Given \\(|A|=%d\\), \\(|B|=%d\\), \\(|A \\cap B|=%d\\), \\(|U|=%d\\), how many elements are in \\(A \\cup B\\)?", s.a, s.b, s.both, s.universe),
			unionSize,
			fmt.Sprintf("|A∪B| = |A| + |B| − |A∩B| = %d + %d − %d = %d", s.a, s.b, s.both, unionSize),
		},
		{
			fmt.Sprintf("Given \\(|A|=%d\\), \\(|B|=%d\\), \\(|A \\cap B|=%d\\), \\(|U|=%d\\), how many elements are in the complement of \\(A \\cup B\\)?", s.a, s.b, s.both, s.universe),
			compUnion,
			fmt.Sprintf("|(A∪B)'| = |U| − |A∪B| = %d − %d = %d", s.universe, unionSize, compUnion),
		},
		{
			fmt.Sprintf("Given \\(|A|=%d\\), \\(|B|=%d\\), \\(|A \\cap B|=%d\\), how many elements are in \\(A\\) only (not in \\(B\\))?", s.a, s.b, s.both),
			onlyA,
			fmt.Sprintf("|A only| = |A| − |A∩B| = %d − %d = %d", s.a, s.both, onlyA),
		},
		{
			fmt.Sprintf("Given \\(|A|=%d\\), \\(|B|=%d\\), \\(|A \\cap B|=%d\\), how many elements are in \\(B\\) only (not in \\(A\\))?", s.a, s.b, s.both),
			onlyB,
			fmt.Sprintf("|B only| = |B| − |A∩B| = %d − %d = %d", s.b, s.both, onlyB),
		},
	}

	q := questions[rand.Intn(len(questions))]
	return generator.Problem{
		Question:    q.question,
		Answer:      fmt.Sprintf("%d", q.answer),
		Explanation: q.exp,
	}
}

// ----- 7. permutations -----

type permutationsGen struct{}

func (g *permutationsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type qType struct {
		question string
		answer   int
		exp      string
	}

	questions := []qType{
		{
			"How many ways to arrange 5 books on a shelf?",
			120,
			"5! = 5 × 4 × 3 × 2 × 1 = 120",
		},
		{
			"How many ways to arrange 4 books on a shelf?",
			24,
			"4! = 4 × 3 × 2 × 1 = 24",
		},
		{
			"How many ways to arrange 3 books on a shelf?",
			6,
			"3! = 3 × 2 × 1 = 6",
		},
		{
			"How many ways can 6 people line up?",
			720,
			"6! = 6 × 5 × 4 × 3 × 2 × 1 = 720",
		},
		{
			"How many ways to pick 1st, 2nd, and 3rd place from 10 runners?",
			720,
			"\\(10 \\times 9 \\times 8 = 720\\) (\\(_{10}P_{3} = 10! / 7!\\))",
		},
		{
			"How many ways to pick 1st and 2nd place from 8 runners?",
			56,
			"8 × 7 = 56 (8P2 = 8! / 6!)",
		},
		{
			"How many ways to arrange the letters A, B, C, D, E?",
			120,
			"5! = 5 × 4 × 3 × 2 × 1 = 120",
		},
		{
			"How many ways to select a president and vice-president from 12 candidates?",
			132,
			"12 × 11 = 132 (12P2 = 12! / 10!)",
		},
	}

	q := questions[rand.Intn(len(questions))]
	return generator.Problem{
		Question:    q.question,
		Answer:      fmt.Sprintf("%d", q.answer),
		Explanation: q.exp,
	}
}

// ----- 8. combinations -----

type combinationsGen struct{}

func (g *combinationsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type combo struct {
		n, k, result int
	}
	combos := []combo{
		{5, 2, 10},
		{5, 3, 10},
		{6, 2, 15},
		{6, 3, 20},
		{7, 2, 21},
		{4, 2, 6},
		{8, 2, 28},
		{7, 3, 35},
		{10, 2, 45},
		{9, 2, 36},
		{4, 3, 4},
		{5, 4, 5},
	}

	type scenario struct {
		question string
		combo    combo
	}
	scenarios := []scenario{
		{"How many ways to choose %d toppings from %d?", combo{}},
		{"How many ways to select %d students from a group of %d?", combo{}},
		{"How many committees of %d can be formed from %d people?", combo{}},
	}

	si := scenarios[rand.Intn(len(scenarios))]
	c := combos[rand.Intn(len(combos))]
	si.combo = c

	return generator.Problem{
		Question:    fmt.Sprintf(si.question, c.k, c.n),
		Answer:      fmt.Sprintf("%d", c.result),
		Explanation: fmt.Sprintf("C(%d,%d) = %d! / (%d!(%d−%d)!) = %d", c.n, c.k, c.n, c.k, c.n, c.k, c.result),
	}
}

// ----- 9. pascal -----

type pascalGen struct{}

func (g *pascalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	n := rand.Intn(max(1, scale)) + 2
	row := pascalRow(n)
	parts := make([]string, len(row))
	for i, v := range row {
		parts[i] = strconv.Itoa(v)
	}
	return generator.Problem{
		Question:    fmt.Sprintf("What is row %d of Pascal's triangle? (answer as comma-separated numbers)", n),
		Answer:      strings.Join(parts, ","),
		Explanation: fmt.Sprintf("Row %d of Pascal's triangle: %s", n, strings.Join(parts, ",")),
	}
}

func pascalRow(n int) []int {
	row := []int{1}
	for i := 1; i <= n; i++ {
		next := make([]int, i+1)
		next[0] = 1
		next[i] = 1
		for j := 1; j < i; j++ {
			next[j] = row[j-1] + row[j]
		}
		row = next
	}
	return row
}

// ----- 10. graph_basics -----

type graphBasicsGen struct{}

func (g *graphBasicsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	vertices := []string{"A", "B", "C", "D", "E"}
	n := len(vertices)

	numEdges := rand.Intn(max(1, scale)) + 3

	type edge struct {
		from, to int
	}
	edges := make([]edge, 0, numEdges)
	adj := make([][]int, n)
	for i := range adj {
		adj[i] = make([]int, n)
	}

	for len(edges) < numEdges {
		a := rand.Intn(n)
		b := rand.Intn(n)
		if a == b || adj[a][b] == 1 {
			continue
		}
		adj[a][b] = 1
		adj[b][a] = 1
		edges = append(edges, edge{a, b})
	}

	edgeStrs := make([]string, len(edges))
	for i, e := range edges {
		edgeStrs[i] = fmt.Sprintf("%s-%s", vertices[e.from], vertices[e.to])
	}
	graphDesc := fmt.Sprintf("Vertices: %s. Edges: %s",
		strings.Join(vertices, ","), strings.Join(edgeStrs, ","))

	type qType struct {
		question string
		answer   string
		exp      string
	}
	questions := []qType{
		{
			fmt.Sprintf("A graph has vertices {%s} and edges {%s}. How many vertices does it have?", strings.Join(vertices, ","), strings.Join(edgeStrs, ",")),
			fmt.Sprintf("%d", n),
			fmt.Sprintf("The graph has %d vertices: %s", n, strings.Join(vertices, ", ")),
		},
		{
			fmt.Sprintf("A graph has vertices {%s} and edges {%s}. How many edges does it have?", strings.Join(vertices, ","), strings.Join(edgeStrs, ",")),
			fmt.Sprintf("%d", len(edges)),
			fmt.Sprintf("The graph has %d edges: %s", len(edges), strings.Join(edgeStrs, ", ")),
		},
	}
	vIdx := rand.Intn(n)
	deg := 0
	for j := 0; j < n; j++ {
		deg += adj[vIdx][j]
	}
	questions = append(questions, qType{
		fmt.Sprintf("In a graph with %s, what is the degree of vertex %s?", graphDesc, vertices[vIdx]),
		fmt.Sprintf("%d", deg),
		fmt.Sprintf("Vertex %s is connected to %d edge(s).", vertices[vIdx], deg),
	})

	q := questions[rand.Intn(len(questions))]
	return generator.Problem{
		Question:    q.question,
		Answer:      q.answer,
		Explanation: q.exp,
	}
}

// ----- 11. graph_paths -----

type graphPathsGen struct{}

func (g *graphPathsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	connected := rand.Intn(2) == 0

	var edges []string
	var adj [4][4]int

	if connected {
		edges = []string{"A-B", "B-C", "C-D"}
		adj[0][1] = 1
		adj[1][0] = 1
		adj[1][2] = 1
		adj[2][1] = 1
		adj[2][3] = 1
		adj[3][2] = 1
	} else {
		edges = []string{"A-B", "C-D"}
		adj[0][1] = 1
		adj[1][0] = 1
		adj[2][3] = 1
		adj[3][2] = 1
	}

	graphStr := fmt.Sprintf("{%s}", strings.Join(edges, ", "))

	if rand.Intn(2) == 0 {
		ans := "yes"
		exp := "There is a path: A-B-C-D"
		if !connected {
			ans = "no"
			exp = "There is no path from A to D because the graph is disconnected."
		}
		return generator.Problem{
			Question:    fmt.Sprintf("In the graph with edges %s, is there a path from A to D? (yes/no)", graphStr),
			Answer:      ans,
			Explanation: exp,
		}
	}

	pathLen := 3
	exp := "The shortest path is A-B-C-D (3 edges)."
	if !connected {
		pathLen = -1
		exp = "There is no path from A to D (the graph is disconnected)."
		return generator.Problem{
			Question:    fmt.Sprintf("In the graph with edges %s, what is the length of the shortest path from A to D? (if none, write 'none')", graphStr),
			Answer:      "none",
			Explanation: exp,
		}
	}

	return generator.Problem{
		Question:    fmt.Sprintf("In the graph with edges %s, what is the length of the shortest path from A to D?", graphStr),
		Answer:      fmt.Sprintf("%d", pathLen),
		Explanation: exp,
	}
}

// ----- 12. trees -----

type treesGen struct{}

func (g *treesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type qType struct {
		question string
		answer   string
		exp      string
	}
	questions := []qType{
		{
			"A tree has 10 vertices. How many edges does it have?",
			"9",
			"A tree with n vertices always has n-1 edges. 10 − 1 = 9.",
		},
		{
			"A tree has 7 vertices. How many edges does it have?",
			"6",
			"A tree with n vertices always has n-1 edges. 7 − 1 = 6.",
		},
		{
			"A tree has 15 vertices. How many edges does it have?",
			"14",
			"A tree with n vertices always has n-1 edges. 15 − 1 = 14.",
		},
		{
			"A tree has 12 edges. How many vertices does it have?",
			"13",
			"A tree with n edges always has n+1 vertices. 12 + 1 = 13.",
		},
		{
			"A tree has 8 edges. How many vertices does it have?",
			"9",
			"A tree with n edges always has n+1 vertices. 8 + 1 = 9.",
		},
		{
			"A connected graph with 6 vertices and no cycles has how many edges?",
			"5",
			"A connected acyclic graph is a tree, so edges = vertices − 1 = 5.",
		},
	}

	q := questions[rand.Intn(len(questions))]
	return generator.Problem{
		Question:    q.question,
		Answer:      q.answer,
		Explanation: q.exp,
	}
}

// ----- 13. recurrence -----

type recurrenceGen struct{}

func (g *recurrenceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type qType struct {
		question string
		answer   string
		exp      string
	}
	questions := []qType{
		{
			"If \\(a_{n} = a_{n-1} + 2\\) and \\(a_{1} = 1\\), what is \\(a_{4}\\)?",
			"7",
			"a_2 = 1+2 = 3, a_3 = 3+2 = 5, a_4 = 5+2 = 7",
		},
		{
			"If \\(a_{n} = a_{n-1} + 3\\) and \\(a_{1} = 2\\), what is \\(a_{4}\\)?",
			"11",
			"a_2 = 2+3 = 5, a_3 = 5+3 = 8, a_4 = 8+3 = 11",
		},
		{
			"If \\(a_{n} = 2a_{n-1}\\) and \\(a_{1} = 3\\), what is \\(a_{4}\\)?",
			"24",
			"\\(a_{2} = 2 \\times 3 = 6\\), \\(a_{3} = 2 \\times 6 = 12\\), \\(a_{4} = 2 \\times 12 = 24\\)",
		},
		{
			"If \\(a_{n} = 2a_{n-1}\\) and \\(a_{1} = 5\\), what is \\(a_{5}\\)?",
			"80",
			"\\(a_{2} = 5 \\times 2 = 10\\), \\(a_{3} = 10 \\times 2 = 20\\), \\(a_{4} = 20 \\times 2 = 40\\), \\(a_{5} = 40 \\times 2 = 80\\)",
		},
		{
			"What is \\(a_{5}\\) if \\(a_{n} = 2n - 1\\)?",
			"9",
			"a_5 = 2(5) − 1 = 10 − 1 = 9",
		},
		{
			"What is \\(a_{6}\\) if \\(a_{n} = 3n + 1\\)?",
			"19",
			"a_6 = 3(6) + 1 = 18 + 1 = 19",
		},
		{
			"What is \\(a_{10}\\) if \\(a_{n} = n^{2}\\)?",
			"100",
			"a_10 = 10^2 = 100",
		},
		{
			"If \\(a_{n} = a_{n-1} + a_{n-2}\\) and \\(a_{1} = 1\\), \\(a_{2} = 1\\), what is \\(a_{6}\\)?",
			"8",
			"a_3=2, a_4=3, a_5=5, a_6=8 (Fibonacci sequence)",
		},
	}

	q := questions[rand.Intn(len(questions))]
	return generator.Problem{
		Question:    q.question,
		Answer:      q.answer,
		Explanation: q.exp,
	}
}

// ----- 14. binomial theorem -----

type binomialTheoremGen struct{}

func (g *binomialTheoremGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		n, k   int
		coeff  int
	}
	entries := []entry{
		{4, 2, 6},
		{5, 2, 10},
		{5, 3, 10},
		{6, 2, 15},
		{6, 3, 20},
		{7, 2, 21},
		{7, 3, 35},
		{8, 2, 28},
		{8, 3, 56},
		{4, 1, 4},
		{4, 3, 4},
		{5, 1, 5},
		{5, 4, 5},
	}
	if rand.Intn(2) == 0 {
		e := entries[rand.Intn(len(entries))]
		return generator.Problem{
			Question:    fmt.Sprintf("Using the binomial theorem, what is the coefficient of \\(x^{%d} y^{%d}\\) in \\((x+y)^{%d}\\)?", e.k, e.n-e.k, e.n),
			Answer:      fmt.Sprintf("%d", e.coeff),
			Explanation: fmt.Sprintf("C(%d,%d) = %d!/(%d!(%d-%d)!) = %d", e.n, e.k, e.n, e.k, e.n, e.k, e.coeff),
		}
	}
	e2 := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    fmt.Sprintf("What is \\(C(%d,%d)\\) in the expansion of \\((x+y)^{%d}\\)?", e2.n, e2.k, e2.n),
		Answer:      fmt.Sprintf("%d", e2.coeff),
		Explanation: fmt.Sprintf("C(%d,%d) = %d", e2.n, e2.k, e2.coeff),
	}
}

// ----- 15. induction -----

type inductionGen struct{}

func (g *inductionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type qType struct {
		question string
		answer   string
		exp      string
	}
	questions := []qType{
		{
			"In induction, to prove P(k) → P(k+1), what step shows P(1) holds? (one word)",
			"base case",
			"The base case verifies that the statement holds for the smallest value (usually n=1).",
		},
		{
			"If P(n): 1+2+...+n = n(n+1)/2, what is P(1)?",
			"1",
			"Left side: 1. Right side: 1(1+1)/2 = 2/2 = 1. Both sides equal 1.",
		},
		{
			"If P(n): 1+2+...+n = n(n+1)/2, what is the right side of P(3)?",
			"6",
			"P(3): 3(3+1)/2 = 3×4/2 = 12/2 = 6",
		},
		{
			"In the inductive step, we assume P(k) is true. What must we prove next?",
			"P(k+1)",
			"The inductive step proves that if P(k) is true, then P(k+1) is also true.",
		},
		{
			"If P(n): 1+3+5+...+(2n-1) = n^2, what is P(1)?",
			"1",
			"Left side: 1. Right side: 1^2 = 1. Both sides equal 1.",
		},
		{
			"If P(n): 1+3+5+...+(2n-1) = n^2, what is P(2)?",
			"4",
			"Left side: 1+3 = 4. Right side: 2^2 = 4.",
		},
		{
			"The principle of mathematical induction is used to prove statements about what set of numbers?",
			"positive integers",
			"Induction is used to prove statements for all positive integers (natural numbers).",
		},
	}

	q := questions[rand.Intn(len(questions))]
	return generator.Problem{
		Question:    q.question,
		Answer:      q.answer,
		Explanation: q.exp,
	}
}
