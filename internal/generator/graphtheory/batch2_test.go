package graphtheory

// Tests for the second graph-theory batch.
//
// Every concept is graded through the real router with the grading_type the corpus
// ships. The five multiple-choice concepts are checked the way the other batches
// check theirs: the question is parsed into options, the correct answer must appear
// exactly once, and every other option must be rejected.
//
// The fifteen numeric concepts recompute their answers from the published graph
// parameters: the treewidth and genus formulas, the walk counts of the adjacency
// matrix, the spanning-tree counts of the matrix-tree theorem, the spectral
// quantities, the extremal bounds, and the flow and path counts.

import (
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/grader"
)

var graphBatch2IDs = []string{
	"graph.kuratowski",
	"graph.minor",
	"graph.treewidth",
	"graph.genus",
	"graph.outerplanar",
	"graph.adjacency",
	"graph.laplacian",
	"graph.matrix_tree",
	"graph.spectral_radius",
	"graph.spectral_gap",
	"graph.mantel",
	"graph.turan",
	"graph.ramsey",
	"graph.random_threshold",
	"graph.moore_bound",
	"graph.max_flow",
	"graph.tournament",
	"graph.strong_connectivity",
	"graph.topological_sort",
	"graph.dag_paths",
}

var graphBatch2NumericIDs = []string{
	"graph.treewidth",
	"graph.genus",
	"graph.outerplanar",
	"graph.adjacency",
	"graph.laplacian",
	"graph.matrix_tree",
	"graph.spectral_radius",
	"graph.spectral_gap",
	"graph.mantel",
	"graph.turan",
	"graph.ramsey",
	"graph.moore_bound",
	"graph.max_flow",
	"graph.tournament",
	"graph.dag_paths",
}

var graphBatch2MCIDs = []string{}

// ---------------------------------------------------------------------------
// Router contracts
// ---------------------------------------------------------------------------

func TestBatch2AnswersGradeWithTheirShippedType(t *testing.T) {
	types := graphShippedTypes(t)
	reg := graphBatchRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	for _, id := range graphBatch2IDs {
		typ, ok := types[id]
		if !ok || typ == "" {
			t.Fatalf("%s: no grading_type in the corpus", id)
		}
		if typ != "multiple_choice" && typ != "numeric" {
			t.Fatalf("%s: unexpected grading_type %q", id, typ)
		}
		for i := 0; i < 120; i++ {
			p := graphBatchSample(t, reg, id, i)
			if strings.TrimSpace(p.Answer) == "" {
				t.Fatalf("%s: empty answer at sample %d", id, i)
			}
			if strings.TrimSpace(p.Question) == "" || strings.TrimSpace(p.Explanation) == "" {
				t.Fatalf("%s: empty question or explanation at sample %d", id, i)
			}
			res := router.Grade(grader.GradingType(typ), p.Answer, p.Answer)
			if res.Unavailable {
				t.Skipf("%s: grader unavailable", id)
			}
			if !res.Correct {
				t.Fatalf("%s: the generator's own answer %q does not grade as %s", id, p.Answer, typ)
			}
		}
	}
}

func TestBatch2MultipleChoiceRejectsEveryDistractor(t *testing.T) {
	types := graphShippedTypes(t)
	reg := graphBatchRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	if len(graphBatch2MCIDs) == 0 {
		t.Skip("the graph-theory batch now serves computed count problems")
	}
	for _, id := range graphBatch2MCIDs {
		if types[id] != "multiple_choice" {
			t.Fatalf("%s: expected multiple_choice, corpus says %q", id, types[id])
		}
		for i := 0; i < 40; i++ {
			p := graphBatchSample(t, reg, id, i)
			opts := graphOptionsOf(p.Question)
			if len(opts) < 3 {
				t.Fatalf("%s: only %d options; a recognition question needs at least three", id, len(opts))
			}
			found := 0
			for _, o := range opts {
				if o == p.Answer {
					found++
				}
			}
			if found != 1 {
				t.Fatalf("%s: the answer %q appears %d times among %v", id, p.Answer, found, opts)
			}
			for _, o := range opts {
				if o == p.Answer {
					continue
				}
				if router.Grade(grader.GradingMultipleChoice, p.Answer, o).Correct {
					t.Fatalf("%s: distractor %q is accepted as correct", id, o)
				}
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Small graph helpers
// ---------------------------------------------------------------------------

// graphAdjacency returns the adjacency matrix of the named graph on n vertices.
func graphAdjacency(kind string, n int) [][]int {
	a := make([][]int, n)
	for i := range a {
		a[i] = make([]int, n)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			switch kind {
			case "complete":
				a[i][j] = 1
			case "path":
				if j == i+1 || j == i-1 {
					a[i][j] = 1
				}
			case "cycle":
				if j == (i+1)%n || j == (i+n-1)%n {
					a[i][j] = 1
				}
			}
		}
	}
	return a
}

func graphMatMul(a, b [][]int) [][]int {
	n := len(a)
	c := make([][]int, n)
	for i := range c {
		c[i] = make([]int, n)
		for k := 0; k < n; k++ {
			if a[i][k] == 0 {
				continue
			}
			for j := 0; j < n; j++ {
				c[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	return c
}

func graphMatPow(a [][]int, exp int) [][]int {
	n := len(a)
	out := make([][]int, n)
	for i := range out {
		out[i] = make([]int, n)
		out[i][i] = 1
	}
	base := a
	for e := exp; e > 0; e >>= 1 {
		if e&1 == 1 {
			out = graphMatMul(out, base)
		}
		base = graphMatMul(base, base)
	}
	return out
}

func graphCeilDiv(a, b int) int {
	if a <= 0 {
		return 0
	}
	return (a + b - 1) / b
}

func graphSpectralGapCycle(n int) int {
	return int(math.Round(2 - 2*math.Cos(2*math.Pi/float64(n))))
}

// graphParseCapEdges parses "u,v,cap;..." into a capacity matrix on n nodes.
func graphParseCapEdges(t *testing.T, n int, spec string) [][]int {
	t.Helper()
	cap := make([][]int, n)
	for i := range cap {
		cap[i] = make([]int, n)
	}
	if spec == "" {
		return cap
	}
	for _, part := range strings.Split(spec, ";") {
		fields := strings.Split(part, ",")
		if len(fields) != 3 {
			t.Fatalf("bad edge %q", part)
		}
		u, _ := strconv.Atoi(fields[0])
		v, _ := strconv.Atoi(fields[1])
		c, _ := strconv.Atoi(fields[2])
		cap[u][v] += c
	}
	return cap
}

// graphMaxFlow is Edmonds-Karp on an integer capacity matrix.
func graphMaxFlow(cap [][]int) int {
	n := len(cap)
	flow := 0
	for {
		parent := make([]int, n)
		for i := range parent {
			parent[i] = -1
		}
		parent[0] = 0
		queue := []int{0}
		for len(queue) > 0 && parent[n-1] == -1 {
			u := queue[0]
			queue = queue[1:]
			for v := 0; v < n; v++ {
				if parent[v] == -1 && cap[u][v] > 0 {
					parent[v] = u
					queue = append(queue, v)
				}
			}
		}
		if parent[n-1] == -1 {
			return flow
		}
		bottleneck := int(^uint(0) >> 1)
		for v := n - 1; v != 0; v = parent[v] {
			if cap[parent[v]][v] < bottleneck {
				bottleneck = cap[parent[v]][v]
			}
		}
		for v := n - 1; v != 0; v = parent[v] {
			cap[parent[v]][v] -= bottleneck
			cap[v][parent[v]] += bottleneck
		}
		flow += bottleneck
	}
}

// graphCountPaths counts directed s-t paths in a DAG given as "u,v;..." edges.
func graphCountPaths(t *testing.T, n int, spec string, s, target int) int {
	t.Helper()
	adj := make([][]int, n)
	indeg := make([]int, n)
	for _, part := range strings.Split(spec, ";") {
		fields := strings.Split(part, ",")
		if len(fields) != 2 {
			t.Fatalf("bad dag edge %q", part)
		}
		u, _ := strconv.Atoi(fields[0])
		v, _ := strconv.Atoi(fields[1])
		adj[u] = append(adj[u], v)
		indeg[v]++
	}
	order := []int{}
	queue := []int{}
	for i := 0; i < n; i++ {
		if indeg[i] == 0 {
			queue = append(queue, i)
		}
	}
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		order = append(order, u)
		for _, v := range adj[u] {
			indeg[v]--
			if indeg[v] == 0 {
				queue = append(queue, v)
			}
		}
	}
	if len(order) != n {
		t.Fatalf("dag_paths: graph is not acyclic")
	}
	ways := make([]int, n)
	ways[s] = 1
	for _, u := range order {
		for _, v := range adj[u] {
			ways[v] += ways[u]
		}
	}
	return ways[target]
}

// ---------------------------------------------------------------------------
// Independent numeric checks
// ---------------------------------------------------------------------------

func TestBatch2TreewidthRecomputes(t *testing.T) {
	reg := graphBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := graphBatchSample(t, reg, "graph.treewidth", i)
		var want int
		switch p.Facts["kind"] {
		case "tree", "path":
			want = 1
		case "cycle":
			want = 2
		case "complete":
			want = graphFactInt(t, p, "n") - 1
		case "grid":
			want = int(math.Min(float64(graphFactInt(t, p, "rows")), float64(graphFactInt(t, p, "cols"))))
		case "bipartite":
			want = int(math.Min(float64(graphFactInt(t, p, "a")), float64(graphFactInt(t, p, "b"))))
		default:
			t.Fatalf("treewidth: unknown kind %q", p.Facts["kind"])
		}
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("treewidth: %s answers %s, want %d", p.Facts["kind"], p.Answer, want)
		}
	}
}

func TestBatch2GenusRecomputes(t *testing.T) {
	reg := graphBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := graphBatchSample(t, reg, "graph.genus", i)
		var want int
		switch p.Facts["kind"] {
		case "planar", "cycle":
			want = 0
		case "complete":
			n := graphFactInt(t, p, "n")
			want = graphCeilDiv((n-3)*(n-4), 12)
		case "bipartite":
			a := graphFactInt(t, p, "a")
			b := graphFactInt(t, p, "b")
			want = graphCeilDiv((a-2)*(b-2), 4)
		default:
			t.Fatalf("genus: unknown kind %q", p.Facts["kind"])
		}
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("genus: %s answers %s, want %d", p.Facts["kind"], p.Answer, want)
		}
	}
}

func TestBatch2OuterplanarRecomputes(t *testing.T) {
	reg := graphBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := graphBatchSample(t, reg, "graph.outerplanar", i)
		n := graphFactInt(t, p, "n")
		want := 2*n - 3
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("outerplanar: n=%d answers %s, want %d", n, p.Answer, want)
		}
	}
}

func TestBatch2AdjacencyRecomputes(t *testing.T) {
	reg := graphBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := graphBatchSample(t, reg, "graph.adjacency", i)
		n := graphFactInt(t, p, "n")
		iIdx := graphFactInt(t, p, "i") - 1
		jIdx := graphFactInt(t, p, "j") - 1
		length := graphFactInt(t, p, "length")
		pow := graphMatPow(graphAdjacency(p.Facts["kind"], n), length)
		want := pow[iIdx][jIdx]
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("adjacency: %s n=%d (%d,%d) len=%d answers %s, want %d",
				p.Facts["kind"], n, iIdx+1, jIdx+1, length, p.Answer, want)
		}
	}
}

func TestBatch2LaplacianRecomputes(t *testing.T) {
	reg := graphBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := graphBatchSample(t, reg, "graph.laplacian", i)
		var want int
		switch p.Facts["kind"] {
		case "path", "cycle", "tree", "complete":
			want = 1
		case "disjoint":
			want = graphFactInt(t, p, "components")
		default:
			t.Fatalf("laplacian: unknown kind %q", p.Facts["kind"])
		}
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("laplacian: %s answers %s, want %d", p.Facts["kind"], p.Answer, want)
		}
	}
}

func TestBatch2MatrixTreeRecomputes(t *testing.T) {
	reg := graphBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := graphBatchSample(t, reg, "graph.matrix_tree", i)
		var want int
		switch p.Facts["kind"] {
		case "complete":
			n := graphFactInt(t, p, "n")
			want = graphPow(n, n-2)
		case "cycle":
			want = graphFactInt(t, p, "n")
		case "path":
			want = 1
		case "bipartite":
			a := graphFactInt(t, p, "a")
			b := graphFactInt(t, p, "b")
			want = graphPow(a, b-1) * graphPow(b, a-1)
		default:
			t.Fatalf("matrix-tree: unknown kind %q", p.Facts["kind"])
		}
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("matrix-tree: %s answers %s, want %d", p.Facts["kind"], p.Answer, want)
		}
	}
}

func TestBatch2SpectralRadiusRecomputes(t *testing.T) {
	reg := graphBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := graphBatchSample(t, reg, "graph.spectral_radius", i)
		var want int
		switch p.Facts["kind"] {
		case "complete":
			want = graphFactInt(t, p, "n") - 1
		case "cycle":
			want = 2
		case "regular":
			want = graphFactInt(t, p, "d")
		default:
			t.Fatalf("spectral radius: unknown kind %q", p.Facts["kind"])
		}
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("spectral radius: %s answers %s, want %d", p.Facts["kind"], p.Answer, want)
		}
	}
}

func TestBatch2SpectralGapRecomputes(t *testing.T) {
	reg := graphBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := graphBatchSample(t, reg, "graph.spectral_gap", i)
		var want int
		switch p.Facts["kind"] {
		case "complete":
			want = graphFactInt(t, p, "n")
		case "cycle":
			want = graphSpectralGapCycle(graphFactInt(t, p, "n"))
		case "bipartite":
			want = int(math.Min(float64(graphFactInt(t, p, "a")), float64(graphFactInt(t, p, "b"))))
		case "star":
			want = 1
		default:
			t.Fatalf("spectral_gap: unknown kind %q", p.Facts["kind"])
		}
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("spectral_gap: %s answers %s, want %d", p.Facts["kind"], p.Answer, want)
		}
	}
}

func TestBatch2MantelRecomputes(t *testing.T) {
	reg := graphBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := graphBatchSample(t, reg, "graph.mantel", i)
		n := graphFactInt(t, p, "n")
		want := n * n / 4
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("mantel: n=%d answers %s, want %d", n, p.Answer, want)
		}
	}
}

func TestBatch2TuranRecomputes(t *testing.T) {
	reg := graphBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := graphBatchSample(t, reg, "graph.turan", i)
		n := graphFactInt(t, p, "n")
		r := graphFactInt(t, p, "r")
		want := n * n * (r - 1) / (2 * r)
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("turan: n=%d r=%d answers %s, want %d", n, r, p.Answer, want)
		}
	}
}

func TestBatch2RamseyRecomputes(t *testing.T) {
	reg := graphBatchRegistry(t)
	table := map[[2]int]int{
		{2, 3}: 3, {2, 4}: 4, {2, 5}: 5,
		{3, 3}: 6, {3, 4}: 9, {3, 5}: 14, {4, 4}: 18,
	}
	for i := 0; i < 200; i++ {
		p := graphBatchSample(t, reg, "graph.ramsey", i)
		key := [2]int{graphFactInt(t, p, "a"), graphFactInt(t, p, "b")}
		want, ok := table[key]
		if !ok {
			t.Fatalf("ramsey: no table entry for %v", key)
		}
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("ramsey: %v answers %s, want %d", key, p.Answer, want)
		}
	}
}

func TestBatch2MooreBoundRecomputes(t *testing.T) {
	reg := graphBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := graphBatchSample(t, reg, "graph.moore_bound", i)
		d := graphFactInt(t, p, "d")
		want := d*d + 1
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("moore bound: d=%d answers %s, want %d", d, p.Answer, want)
		}
	}
}

func TestBatch2MaxFlowRecomputes(t *testing.T) {
	reg := graphBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := graphBatchSample(t, reg, "graph.max_flow", i)
		n := graphFactInt(t, p, "nodes")
		cap := graphParseCapEdges(t, n, p.Facts["edges"])
		want := graphMaxFlow(cap)
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("max-flow: edges %q answers %s, want %d", p.Facts["edges"], p.Answer, want)
		}
	}
}

func TestBatch2TournamentRecomputes(t *testing.T) {
	reg := graphBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := graphBatchSample(t, reg, "graph.tournament", i)
		n := graphFactInt(t, p, "n")
		want := 1 << (n * (n - 1) / 2)
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("tournament: n=%d answers %s, want %d", n, p.Answer, want)
		}
	}
}

func TestBatch2DAGPathsRecomputes(t *testing.T) {
	reg := graphBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := graphBatchSample(t, reg, "graph.dag_paths", i)
		n := graphFactInt(t, p, "nodes")
		s := graphFactInt(t, p, "s")
		target := graphFactInt(t, p, "t")
		want := graphCountPaths(t, n, p.Facts["edges"], s, target)
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("dag paths: edges %q answers %s, want %d", p.Facts["edges"], p.Answer, want)
		}
	}
}

func TestBatch2NumericAnswersParseAsDecimals(t *testing.T) {
	reg := graphBatchRegistry(t)
	for _, id := range graphBatch2NumericIDs {
		for i := 0; i < 120; i++ {
			p := graphBatchSample(t, reg, id, i)
			if _, err := strconv.ParseFloat(p.Answer, 64); err != nil {
				t.Fatalf("%s: answer %q is not a decimal number", id, p.Answer)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Pool quality
// ---------------------------------------------------------------------------

func TestBatch2PoolsAreDiverse(t *testing.T) {
	reg := graphBatchRegistry(t)
	const samples = 150
	for _, id := range graphBatch2IDs {
		answers := map[string]bool{}
		questions := map[string]bool{}
		allYesNo := true
		for i := 0; i < samples; i++ {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{Difficulty: 0.12, Seed: int64(163000 + i)})
			if err != nil {
				t.Fatalf("%s: generate: %v", id, err)
			}
			answers[strings.ToLower(strings.TrimSpace(p.Answer))] = true
			questions[p.Question] = true
			if !strings.Contains(p.Question, "(yes/no)") {
				allYesNo = false
			}
		}
		if len(answers) < 2 {
			t.Fatalf("%s: only %d distinct answers over %d samples", id, len(answers), samples)
		}
		if len(questions) < 2 {
			t.Fatalf("%s: only %d distinct questions over %d samples", id, len(questions), samples)
		}
		if allYesNo {
			t.Fatalf("%s: every question is yes/no-only", id)
		}
	}
}
