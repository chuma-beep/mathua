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
	reg.Register("discrete.logic.equivalence", &equivalenceGen{})
	reg.Register("discrete.sets.relations", &relationsGen{})
	reg.Register("discrete.combinatorics.pigeonhole", &pigeonholeConvertedGen{})
	reg.Register("discrete.graphs.eulerian", &eulerianGen{})
	reg.Register("discrete.proof.strong_induction", &strongInductionGen{})
	reg.Register("discrete.combinatorics.stars_bars", &starsBarsConvertedGen{})
	reg.Register("discrete.graphs.coloring", &coloringGen{})
	reg.Register("discrete.proof.contradiction", &contradictionGen{})
	reg.Register("discrete.sets.inclusion_exclusion", &inclusionExclusionConvertedGen{})
	reg.Register("discrete.sequences.generating_functions", &generatingFunctionsGen{})

	registerBatch(reg)
	registerCombinatoricsBatch(reg)
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
		{"x is large.", "no", "It is ambiguous, 'large' is not well-defined without context."},
		{"5 is an even number.", "yes", "It is a declarative statement with a definite truth value (false)."},
		{"What time is it?", "no", "It is a question, not a declarative statement."},
		{"Please sit down.", "no", "It is a command, not a declarative statement."},
		{"1 + 1 = 3", "yes", "It is a declarative statement with a definite truth value (false)."},
		{"This sentence is false.", "no", "It is a paradox, it cannot have a consistent truth value."},
	}
	// Production question: count propositions in a random sample of 3-5
	// statements, so every draw has a typed numeric answer.
	k := 3 + rand.Intn(3)
	idx := rand.Perm(len(table))[:k]
	count := 0
	parts := make([]string, k)
	for i, j := range idx {
		if table[j].isProp == "yes" {
			count++
		}
		parts[i] = fmt.Sprintf("(%s) '%s'", string(rune('a'+i)), table[j].statement)
	}
	return generator.Problem{
		Question:    fmt.Sprintf("How many of these are propositions? %s (enter a number)", strings.Join(parts, " ")),
		Answer:      fmt.Sprintf("%d", count),
		Explanation: fmt.Sprintf("%d of the %d statements are declarative with a definite truth value.", count, k),
	}
}

// ----- 2. connectives -----

type connectivesGen struct{}

func (g *connectivesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type conn struct {
		expr string
		fn   func(p, q bool) bool
	}
	conns := []conn{
		{"p AND q", func(p, q bool) bool { return p && q }},
		{"p OR q", func(p, q bool) bool { return p || q }},
		{"p NAND q", func(p, q bool) bool { return !(p && q) }},
		{"p NOR q", func(p, q bool) bool { return !(p || q) }},
		{"p XOR q", func(p, q bool) bool { return p != q }},
		{"NOT p", func(p, q bool) bool { return !p }},
		{"p -> q", func(p, q bool) bool { return !p || q }},
		{"p <-> q", func(p, q bool) bool { return p == q }},
		{"NOT q", func(p, q bool) bool { return !q }},
		{"p AND NOT q", func(p, q bool) bool { return p && !q }},
		{"NOT p AND q", func(p, q bool) bool { return !p && q }},
		{"NOT p OR q", func(p, q bool) bool { return !p || q }},
		{"p OR NOT q", func(p, q bool) bool { return p || !q }},
		{"p XNOR q", func(p, q bool) bool { return p == q }},
	}
	c := conns[rand.Intn(len(conns))]
	count := 0
	for p := 0; p < 2; p++ {
		for q := 0; q < 2; q++ {
			if c.fn(p == 1, q == 1) {
				count++
			}
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Of the four assignments of p and q, for how many is %s true? (enter a number)", c.expr),
		Answer:      strconv.Itoa(count),
		Explanation: fmt.Sprintf("Exactly %d of the four assignments make %s true.", count, c.expr),
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

	rows := 1 << vars
	trueRows := 0
	for mask := 0; mask < rows; mask++ {
		p := mask&1 != 0
		q := mask&2 != 0
		var v bool
		switch e.name {
		case "p AND q":
			v = p && q
		case "p OR q":
			v = p || q
		case "NOT p":
			v = !p
		case "p XOR q":
			v = p != q
		}
		if v {
			trueRows++
		}
	}
	switch rand.Intn(3) {
	case 0:
		return generator.Problem{
			Question:    fmt.Sprintf("How many rows are in a truth table for %s?", e.description),
			Answer:      fmt.Sprintf("%d", rows),
			Explanation: fmt.Sprintf("With %d variable(s), there are 2^%d = %d rows.", vars, vars, rows),
		}
	case 1:
		return generator.Problem{
			Question:    fmt.Sprintf("How many of the %d rows in the truth table for %s are true?", rows, e.description),
			Answer:      fmt.Sprintf("%d", trueRows),
			Explanation: fmt.Sprintf("%s is true in %d of the %d rows.", e.description, trueRows, rows),
		}
	default:
		falseRows := rows - trueRows
		return generator.Problem{
			Question:    fmt.Sprintf("How many of the %d rows in the truth table for %s are false?", rows, e.description),
			Answer:      fmt.Sprintf("%d", falseRows),
			Explanation: fmt.Sprintf("%s is false in %d of the %d rows.", e.description, falseRows, rows),
		}
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
		cond func(x int) bool
		desc string
	}
	entries := []qEntry{
		{func(x int) bool { return x > 0 }, "> 0"},
		{func(x int) bool { return x > 5 }, "> 5"},
		{func(x int) bool { return x < 10 }, "< 10"},
		{func(x int) bool { return x == 0 }, "= 0"},
		{func(x int) bool { return x%2 == 0 }, "is even"},
	}
	e := entries[rand.Intn(len(entries))]
	count := 0
	for _, v := range s {
		if e.cond(v) {
			count++
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("How many elements of %s satisfy x %s? (enter a number)", setStr, e.desc),
		Answer:      strconv.Itoa(count),
		Explanation: fmt.Sprintf("Counting the elements of %s with x %s gives %d.", setStr, e.desc, count),
	}
}

// ----- 5. set_operations -----

type setOpsGen struct{}

func (g *setOpsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return setOperationsProblem()
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
	// Easy rows 2..4: at Difficulty 0.12 scale=1, so all three rows are
	// reachable and the pool discriminates (1,2,1 vs 1,3,3,1 vs 1,4,6,4,1).
	n := 2 + rand.Intn(3)
	if scale > 3 {
		n = 2 + rand.Intn(6) // rows 2..7
	}
	row := pascalRow(n)
	parts := make([]string, len(row))
	for i, v := range row {
		parts[i] = strconv.Itoa(v)
	}
	return generator.Problem{
		Question:    fmt.Sprintf("What is row %d of Pascal's triangle, counting the top 1 as row 0? (answer as comma-separated numbers)", n),
		Answer:      strings.Join(parts, ","),
		Explanation: fmt.Sprintf("Row %d, counting the top 1 as row 0: %s", n, strings.Join(parts, ",")),
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
	n := 4
	if ctx.Difficulty > 0.6 {
		n = 5
	}
	verts := []string{"A", "B", "C", "D", "E"}[:n]
	type edge struct{ u, v int }
	var pairs []edge
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			pairs = append(pairs, edge{i, j})
		}
	}
	rand.Shuffle(len(pairs), func(a, b int) { pairs[a], pairs[b] = pairs[b], pairs[a] })
	adj := make([][]int, n)
	addEdge := func(u, v int) {
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	hasEdge := func(u, v int) bool {
		for _, w := range adj[u] {
			if w == v {
				return true
			}
		}
		return false
	}
	if rand.Intn(2) == 0 {
		// Connected: random spanning tree over all vertices plus extras.
		perm := rand.Perm(n)
		for i := 1; i < n; i++ {
			addEdge(perm[i-1], perm[i])
		}
		for _, e := range pairs {
			if rand.Intn(3) == 0 && !hasEdge(e.u, e.v) {
				addEdge(e.u, e.v)
			}
		}
	} else {
		// Disconnected: split vertices into two nonempty groups,
		// edges only within groups so no path crosses the cut.
		cut := 1 + rand.Intn(n-1)
		var within []edge
		for _, e := range pairs {
			if (e.u < cut) == (e.v < cut) {
				within = append(within, e)
			}
		}
		for _, e := range within[:1+rand.Intn(len(within))] {
			addEdge(e.u, e.v)
		}
	}
	var edgeStrs []string
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if hasEdge(i, j) {
				edgeStrs = append(edgeStrs, verts[i]+"-"+verts[j])
			}
		}
	}
	graphStr := fmt.Sprintf("{%s}", strings.Join(edgeStrs, ", "))
	// BFS distances and parents from verts[0].
	dist := make([]int, n)
	parent := make([]int, n)
	for i := range dist {
		dist[i] = -1
		parent[i] = -1
	}
	dist[0] = 0
	queue := []int{0}
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		for _, w := range adj[v] {
			if dist[w] == -1 {
				dist[w] = dist[v] + 1
				parent[w] = v
				queue = append(queue, w)
			}
		}
	}
	// Connected components via BFS.
	comp := make([]int, n)
	ncomp := 0
	for i := 0; i < n; i++ {
		if comp[i] != 0 {
			continue
		}
		ncomp++
		comp[i] = ncomp
		q := []int{i}
		for len(q) > 0 {
			v := q[0]
			q = q[1:]
			for _, w := range adj[v] {
				if comp[w] == 0 {
					comp[w] = ncomp
					q = append(q, w)
				}
			}
		}
	}
	src, dst := verts[0], verts[n-1]
	var pathStr string
	if dist[n-1] != -1 {
		var rev []string
		for v := n - 1; v != -1; v = parent[v] {
			rev = append([]string{verts[v]}, rev...)
		}
		pathStr = strings.Join(rev, "-")
	}
	compSize := 0
	for i := 0; i < n; i++ {
		if comp[i] == comp[0] {
			compSize++
		}
	}
	switch rand.Intn(3) {
	case 0:
		if dist[n-1] == -1 {
			return generator.Problem{
				Question:    fmt.Sprintf("In the graph with edges %s, what is the length of the shortest path from %s to %s? (enter -1 if there is no path)", graphStr, src, dst),
				Answer:      "-1",
				Explanation: fmt.Sprintf("There is no path from %s to %s (the graph is disconnected).", src, dst),
			}
		}
		return generator.Problem{
			Question:    fmt.Sprintf("In the graph with edges %s, what is the length of the shortest path from %s to %s?", graphStr, src, dst),
			Answer:      fmt.Sprintf("%d", dist[n-1]),
			Explanation: fmt.Sprintf("The shortest path is %s (%d edges).", pathStr, dist[n-1]),
		}
	case 1:
		var parts []string
		for c := 1; c <= ncomp; c++ {
			var vs []string
			for i := 0; i < n; i++ {
				if comp[i] == c {
					vs = append(vs, verts[i])
				}
			}
			parts = append(parts, "{"+strings.Join(vs, ", ")+"}")
		}
		return generator.Problem{
			Question:    fmt.Sprintf("How many connected components does the graph with edges %s have? (enter a number)", graphStr),
			Answer:      fmt.Sprintf("%d", ncomp),
			Explanation: fmt.Sprintf("The components are %s, so there are %d.", strings.Join(parts, " and "), ncomp),
		}
	default:
		return generator.Problem{
			Question:    fmt.Sprintf("In the graph with edges %s, how many vertices are in the connected component that contains %s?", graphStr, src),
			Answer:      fmt.Sprintf("%d", compSize),
			Explanation: fmt.Sprintf("The component containing %s has %d vertices.", src, compSize),
		}
	}
}

// ----- 12. trees -----

type treesGen struct{}

func (g *treesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	n := 2 + rand.Intn(39)
	if rand.Intn(2) == 0 {
		return generator.Problem{
			Question:    fmt.Sprintf("A tree has %d vertices. How many edges does it have?", n),
			Answer:      strconv.Itoa(n - 1),
			Explanation: fmt.Sprintf("A tree with n vertices always has n-1 edges: %d - 1 = %d, because a connected acyclic graph on n vertices has exactly n-1 edges.", n, n-1),
			Facts:       map[string]string{"vertices": strconv.Itoa(n), "edges": strconv.Itoa(n - 1)},
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("A tree has %d edges. How many vertices does it have?", n-1),
		Answer:      strconv.Itoa(n),
		Explanation: fmt.Sprintf("A tree with m edges has m+1 vertices: %d + 1 = %d.", n-1, n),
		Facts:       map[string]string{"vertices": strconv.Itoa(n), "edges": strconv.Itoa(n - 1)},
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
		n, k  int
		coeff int
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
	return inductionProblem()
}

// ----- equivalence -----

type equivalenceGen struct{}

func (g *equivalenceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return equivalenceProblem()
}

// ----- relations -----

type relationsGen struct{}

func (g *relationsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return relationsProblem()
}

// ----- pigeonhole -----

type pigeonholeGen struct{}

func (g *pigeonholeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q string
		a string
		e string
	}
	easy := []entry{
		{"Among 13 people, how many must share a birth month? (enter a number)", "2", "13 pigeons, 12 holes: ceil(13/12) = 2 share a month."},
		{"Among 25 people, how many must share a birth month? (enter a number)", "3", "ceil(25/12) = 3: 2 per month holds only 24 < 25."},
		{"If 30 students get 4 possible grades, how many must share a grade? (enter a number)", "8", "ceil(30/4) = 8: 7 per grade holds only 28 < 30."},
		{"If 10 socks are black or white, what is the largest k so that k socks must share a color? (enter a number)", "5", "Worst case splits 5-5, so 5 of one color is guaranteed but 6 is not."},
		{"Among 5 points in a unit square, must two be within √2/2? (yes/no)", "yes", "Partition square into 4 quarters (side 0.5); 5 pigeons → one quarter has 2 points, diagonal √2/2."},
		{"If 10 socks are either black or white, must 6 be same color? (yes/no)", "no", "Worst case 5-5 split, so 10 does not guarantee 6 of one color (need 11)."},
	}
	hard := []entry{
		{"What is ⌈100/12⌉? (enter a number)", "9", "Ceiling of 100/12 ≈8.33 →9 pigeons in one hole."},
		{"Among n+1 integers, must two have same remainder mod n? (yes/no)", "yes", "n possible remainders, n+1 numbers → pigeonhole."},
		{"If 30 students have 4 possible grades, must ≥8 share a grade? (yes/no)", "yes", "⌈30/4⌉=8."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

// ----- eulerian -----

type eulerianGen struct{}

func (g *eulerianGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return eulerianProblem()
}

// ----- strong induction -----

type strongInductionGen struct{}

func (g *strongInductionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return strongInductionProblem()
}

type starsBarsGen struct{}

func (g *starsBarsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q string
		a string
		e string
	}
	easy := []entry{
		{"How many nonnegative solutions to x1+x2=3? (enter a number)", "4", "C(3+2-1,2-1)=C(4,1)=4: (3,0),(2,1),(1,2),(0,3)."},
		{"How many nonnegative solutions to x1+x2=4? (enter a number)", "5", "C(4+2-1,2-1)=C(5,1)=5: x1 = 0..4 with x2 forced."},
		{"How many solutions to x1+x2+x3=2 with xi≥0? (enter a number)", "6", "C(2+3-1,3-1)=C(4,2)=6."},
		{"How many ways to put 2 identical balls into 2 boxes? (enter a number)", "3", "C(2+2-1,2-1)=C(3,1)=3: (2,0),(1,1),(0,2)."},
		{"Does stars and bars count distributions of n identical objects into k boxes? (yes/no)", "yes", "C(n+k-1,k-1)."},
	}
	hard := []entry{
		{"How many solutions to x1+x2+x3=5 with xi≥1? (enter a number)", "6", "Set yi=xi-1: y1+y2+y3=2 → C(4,2)=6."},
		{"How many ways to put 4 identical balls into 3 boxes? (enter a number)", "15", "C(4+3-1,3-1)=C(6,2)=15."},
		{"Does stars and bars use C(n+k-1,k-1)? (yes/no)", "yes", "Formula."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type coloringGen struct{}

func (g *coloringGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return coloringProblem()
}

type contradictionGen struct{}

func (g *contradictionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return contradictionProblem()
}

type inclusionExclusionGen struct{}

func (g *inclusionExclusionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q string
		a string
		e string
	}
	easy := []entry{
		{"For |A|=3,|B|=4,|A∩B|=1, what is |A∪B|? (enter a number)", "6", "|A∪B|=|A|+|B|-|A∩B|=3+4-1=6."},
		{"For |A|=5,|B|=4,|A∩B|=2, what is |A∪B|? (enter a number)", "7", "|A∪B|=5+4-2=7."},
		{"For |A|=8,|B|=12,|A∩B|=3, what is |A∪B|? (enter a number)", "17", "|A∪B|=8+12-3=17."},
		{"How many integers ≤6 are divisible by neither 2 nor 3? (enter a number)", "2", "6-floor(6/2)-floor(6/3)+floor(6/6)=6-3-2+1=2 (1 and 5)."},
		{"Does |A∪B|=|A|+|B|-|A∩B|? (yes/no)", "yes", "Two sets: add both, subtract the double-counted overlap."},
		{"Is |A∪B| = |A|+|B| for overlapping sets? (yes/no)", "no", "Overlapping sets double-count the intersection; subtract it once."},
	}
	hard := []entry{
		{"How many integers ≤10 not divisible by 2 or 3? (enter a number)", "3", "10 - floor(10/2)-floor(10/3)+floor(10/6)=3 (1,5,7)."},
		{"Does derangement count use inclusion-exclusion? (yes/no)", "yes", "Count permutations with no fixed point."},
		{"Is principle: |∪ Ai| = Σ|Ai| - Σ|Ai∩Aj| + ...? (yes/no)", "yes", "General formula."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type generatingFunctionsGen struct{}

func (g *generatingFunctionsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return generatingFunctionsProblem()
}
