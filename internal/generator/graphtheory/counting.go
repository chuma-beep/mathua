package graphtheory

// The graph-theory counting generators. Each pool holds sixteen statements,
// eight true and eight false, so the answer is a count from 0 to 8 rather than
// a choice among four labels. The count is drawn first and matching statements
// are selected afterwards, which keeps all nine answers equally likely.

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

// graphVerdict is one statement in a counting pool.
type graphVerdict struct {
	label string
	yes   bool
	note  string
}

const graphDraw = 8

// countGraphVerdicts draws the answer uniformly from 0 to 8, then takes that
// many true statements and the rest false ones.
func countGraphVerdicts(lead string, pool []graphVerdict) generator.Problem {
	var yes, no []graphVerdict
	for _, v := range pool {
		if v.yes {
			yes = append(yes, v)
		} else {
			no = append(no, v)
		}
	}
	rand.Shuffle(len(yes), func(i, j int) { yes[i], yes[j] = yes[j], yes[i] })
	rand.Shuffle(len(no), func(i, j int) { no[i], no[j] = no[j], no[i] })

	count := rand.Intn(graphDraw + 1)
	chosen := make([]graphVerdict, 0, graphDraw)
	var question string
	for attempt := 0; attempt < 40; attempt++ {
		if count > len(yes) {
			count = len(yes)
		}
		if graphDraw-count > len(no) {
			count = graphDraw - len(no)
		}
		chosen = chosen[:0]
		chosen = append(chosen, yes[:count]...)
		chosen = append(chosen, no[:graphDraw-count]...)
		rand.Shuffle(len(chosen), func(i, j int) { chosen[i], chosen[j] = chosen[j], chosen[i] })

		parts := make([]string, 0, graphDraw)
		for i, v := range chosen {
			parts = append(parts, fmt.Sprintf("(%c) %s", 'a'+i, v.label))
		}
		question = lead + " " + strings.Join(parts, " ") + " (enter a number)"
		if !strings.Contains(question, strconv.Itoa(count)) {
			break
		}
		count = rand.Intn(graphDraw + 1)
	}

	notes := make([]string, 0, len(chosen))
	for _, v := range chosen {
		verdict := "false"
		if v.yes {
			verdict = "true"
		}
		notes = append(notes, fmt.Sprintf("%s: %s, %s", v.label, verdict, v.note))
	}
	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(count),
		Explanation: fmt.Sprintf("Exactly %d of these statements are true. %s", count, strings.Join(notes, " ")),
		Facts:       map[string]string{"count": strconv.Itoa(count)},
	}
}

func isomorphismProblem() generator.Problem {
	return countGraphVerdicts("How many of these statements about graph isomorphism are true?", []graphVerdict{
		{"An isomorphism is a bijection on vertices", true, "it matches the two vertex sets one to one"},
		{"An isomorphism preserves adjacency", true, "edges map to edges"},
		{"An isomorphism preserves non-adjacency", true, "non-edges map to non-edges"},
		{"The degree multiset is an isomorphism invariant", true, "each vertex keeps its degree"},
		{"Isomorphism is an equivalence relation", true, "identity, inverse and composition all work"},
		{"Isomorphic graphs have the same number of components", true, "components map to components"},
		{"Isomorphic graphs have the same number of triangles", true, "triangles map to triangles"},
		{"A self-complementary graph is isomorphic to its complement", true, "P_4 is the smallest example"},
		{"Graphs of the same order are isomorphic", false, "many graphs of the same order differ"},
		{"Graphs with the same degree sequence are isomorphic", false, "the sequence is necessary but not sufficient"},
		{"An isomorphism must preserve the vertex names", false, "it relabels vertices freely"},
		{"The chromatic number determines the isomorphism class", false, "many non-isomorphic graphs share it"},
		{"An isomorphism may map an edge to a non-edge", false, "that would break the adjacency relation"},
		{"Isomorphic graphs always have the same drawing", false, "drawings are not part of the graph"},
		{"A bijection preserving edge counts is an isomorphism", false, "the adjacency relation itself must be preserved"},
		{"The complement of a graph is always isomorphic to it", false, "only self-complementary graphs are"},
	})
}

func subgraphProblem() generator.Problem {
	return countGraphVerdicts("How many of these statements about subgraphs are true?", []graphVerdict{
		{"An induced subgraph keeps every edge between chosen vertices", true, "that is the definition of induced"},
		{"A spanning subgraph keeps the whole vertex set", true, "only edges are removed"},
		{"Deleting a vertex also deletes its incident edges", true, "no edge may dangle"},
		{"Every graph is a subgraph of itself", true, "the relation is reflexive"},
		{"A spanning tree is a spanning subgraph", true, "it keeps the vertices and drops edges"},
		{"A subgraph has at most as many edges as the original", true, "edges are only removed"},
		{"Every induced subgraph is a subgraph", true, "an induced subgraph is a special subgraph"},
		{"The empty graph is a subgraph of every graph", true, "take no vertices and no edges"},
		{"An induced subgraph keeps every edge of the original graph", false, "only edges inside the chosen set remain"},
		{"Every subgraph is induced", false, "a subgraph may omit an internal edge"},
		{"Every subgraph is spanning", false, "a subgraph may drop vertices"},
		{"Deleting a vertex leaves its incident edges in place", false, "those edges lose an endpoint"},
		{"A subgraph always has fewer vertices", false, "it may have the same vertices"},
		{"A subgraph always has fewer edges", false, "it may have the same edges"},
		{"Adding edges always produces a subgraph", false, "a subgraph cannot gain edges"},
		{"The subgraph relation is symmetric", false, "it is reflexive and transitive but not symmetric"},
	})
}

func homomorphismProblem() generator.Problem {
	return countGraphVerdicts("How many of these statements about graph homomorphisms are true?", []graphVerdict{
		{"A homomorphism maps every edge to an edge", true, "that is the defining condition"},
		{"A proper k-colouring is a homomorphism to K_k", true, "colours are the target vertices"},
		{"A graph is bipartite exactly when it maps to K_2", true, "that is a proper two-colouring"},
		{"An isomorphism is a bijective homomorphism with a homomorphism inverse", true, "this is the structural characterisation"},
		{"A homomorphism need not be injective", true, "several vertices may share an image"},
		{"A homomorphism need not be surjective", true, "the image may miss target vertices"},
		{"The composition of two homomorphisms is a homomorphism", true, "edges map through both maps"},
		{"Every graph maps to K_1", true, "send every vertex to the only vertex"},
		{"A homomorphism preserves non-edges", false, "only edges are required to map to edges"},
		{"Every homomorphism is injective", false, "colourings collapse vertices"},
		{"Every homomorphism is surjective", false, "the image may be a proper part"},
		{"Every homomorphism is an isomorphism", false, "bijectivity is an extra condition"},
		{"A homomorphism preserves vertex degrees", false, "several vertices may collapse onto one"},
		{"An odd cycle maps to K_2", false, "odd cycles are not bipartite"},
		{"A homomorphism may map an edge to a non-edge", false, "that violates the definition"},
		{"The target of a homomorphism must be complete", false, "any graph can be a target"},
	})
}

func mengerProblem() generator.Problem {
	return countGraphVerdicts("How many of these statements about Menger's theorem are true?", []graphVerdict{
		{"The maximum number of internally vertex-disjoint u-v paths equals the minimum vertex cut", true, "this is the vertex form"},
		{"The edge form counts edge-disjoint paths and edge cuts", true, "the same min-max duality holds"},
		{"The vertex connectivity is at most the edge connectivity", true, "kappa <= lambda"},
		{"The edge connectivity is at most the minimum degree", true, "lambda <= delta"},
		{"A graph is k-connected exactly when it has k internally disjoint paths between every pair", true, "this is the path form of k-connectivity"},
		{"Menger's theorem is a min-max duality", true, "a maximum packing equals a minimum cover"},
		{"A complete graph K_n has vertex connectivity n-1", true, "removing n-2 vertices leaves a connected pair"},
		{"A cycle has vertex connectivity 2", true, "two vertices separate it"},
		{"The maximum number of disjoint paths equals the maximum degree", false, "degree is only an upper bound"},
		{"Menger's theorem applies only to trees", false, "it applies to all graphs"},
		{"Menger's theorem applies only to bipartite graphs", false, "it is general"},
		{"Vertex connectivity always equals edge connectivity", false, "they can differ"},
		{"The edge connectivity always equals the minimum degree", false, "it can be smaller"},
		{"k-connectivity only requires minimum degree at least k", false, "the graph must stay connected after removing k-1 vertices"},
		{"A disconnected graph is 1-connected", false, "it has connectivity 0"},
		{"The minimum vertex cut may be larger than the maximum number of disjoint paths", false, "Menger's theorem makes them equal"},
	})
}

func perfectMatchingProblem() generator.Problem {
	return countGraphVerdicts("How many of these statements about perfect matchings are true?", []graphVerdict{
		{"A perfect matching covers every vertex exactly once", true, "that is the definition"},
		{"A graph with a perfect matching has even order", true, "vertices are paired"},
		{"A perfect matching is a maximum matching", true, "it saturates every vertex"},
		{"Every regular bipartite graph has a perfect matching", true, "Hall's condition follows from regularity"},
		{"A path with an odd number of vertices has no perfect matching", true, "one vertex would remain"},
		{"K_4 has a perfect matching", true, "pair the four vertices"},
		{"A perfect matching on n vertices has n/2 edges", true, "each edge covers two vertices"},
		{"A graph with an isolated vertex has no perfect matching", true, "the isolated vertex cannot be covered"},
		{"A perfect matching covers every edge exactly once", false, "it covers vertices, not edges"},
		{"A graph of odd order can have a perfect matching", false, "pairing leaves one vertex out"},
		{"Every graph has a perfect matching", false, "a star with three leaves does not"},
		{"A perfect matching is a minimum vertex cover", false, "those are different objects"},
		{"A star K_{1,3} has a perfect matching", false, "its leaves cannot all be paired"},
		{"A perfect matching on n vertices has n edges", false, "it has n/2 edges"},
		{"A graph with a perfect matching must be connected", false, "two disjoint edges form a matching"},
		{"A perfect matching must use every edge", false, "it is a set of disjoint edges"},
	})
}

func tutteProblem() generator.Problem {
	return countGraphVerdicts("How many of these statements about Tutte's theorem are true?", []graphVerdict{
		{"Tutte's theorem characterises graphs with a perfect matching", true, "it gives a necessary and sufficient condition"},
		{"The condition counts odd components of G - S", true, "odd components must be bounded by |S|"},
		{"Every set S must leave at most |S| odd components", true, "this is the Tutte condition"},
		{"A graph of odd order fails the condition at the empty set", true, "the empty set leaves one odd component"},
		{"Tutte's theorem generalises Hall's marriage theorem", true, "Hall is the bipartite special case"},
		{"A bridge can force odd components on both sides", true, "removing an endpoint exposes them"},
		{"The condition is necessary and sufficient", true, "that is the content of the theorem"},
		{"Tutte's theorem applies to all finite graphs", true, "no bipartite hypothesis is needed"},
		{"Tutte's theorem characterises Eulerian graphs", false, "it is about perfect matchings"},
		{"The condition counts even components of G - S", false, "only odd components matter"},
		{"The condition is only necessary", false, "it is also sufficient"},
		{"Every graph of odd order satisfies the condition", false, "it fails at S empty"},
		{"Tutte's theorem generalises Brooks' theorem", false, "it generalises Hall's theorem"},
		{"A graph with a bridge can never have a perfect matching", false, "a bridge can still sit inside one"},
		{"Tutte's theorem applies only to bipartite graphs", false, "it is general"},
		{"The empty set is not allowed as S", false, "S may be empty"},
	})
}

func brooksProblem() generator.Problem {
	return countGraphVerdicts("How many of these statements about Brooks' theorem are true?", []graphVerdict{
		{"The chromatic number is at most the maximum degree except for complete graphs and odd cycles", true, "those are the only exceptions"},
		{"The exceptions are exactly complete graphs and odd cycles", true, "Brooks' theorem is sharp"},
		{"A connected graph of maximum degree 3 has chromatic number at most 3 unless it is K_4", true, "the Delta = 3 case"},
		{"Brooks' theorem strengthens the greedy bound Delta + 1", true, "it saves one colour outside the exceptions"},
		{"A connected graph of maximum degree 2 is a path or a cycle", true, "that is the classification"},
		{"An even cycle has chromatic number 2", true, "it is bipartite"},
		{"An odd cycle has chromatic number 3", true, "it is not bipartite"},
		{"K_n has chromatic number n", true, "every pair is adjacent"},
		{"The bound uses the minimum degree", false, "it uses the maximum degree"},
		{"Trees and paths are exceptions to Brooks' theorem", false, "they satisfy the bound"},
		{"A graph of maximum degree 3 always has chromatic number 2", false, "K_4 needs 4"},
		{"The chromatic number always equals the maximum degree", false, "cycles show the gap"},
		{"A connected graph of maximum degree 2 is always a tree", false, "it may be a cycle"},
		{"An odd cycle has chromatic number 2", false, "it needs 3"},
		{"K_4 has chromatic number 3", false, "it needs 4"},
		{"Brooks' theorem applies only to bipartite graphs", false, "it is general"},
	})
}

func perfectGraphProblem() generator.Problem {
	return countGraphVerdicts("How many of these statements about perfect graphs are true?", []graphVerdict{
		{"A graph is perfect when clique number equals chromatic number on every induced subgraph", true, "this is the definition"},
		{"Every bipartite graph is perfect", true, "its clique number and chromatic number are both at most 2"},
		{"An odd hole is not perfect", true, "it needs 3 colours but has clique number 2"},
		{"The strong perfect graph theorem forbids odd holes and odd antiholes", true, "those are the only obstructions"},
		{"Every chordal graph is perfect", true, "chordal graphs have simplicial vertices"},
		{"The complement of a perfect graph is perfect", true, "this is the perfect graph theorem"},
		{"Every complete graph is perfect", true, "clique number equals chromatic number"},
		{"Every tree is perfect", true, "trees are bipartite"},
		{"A perfect graph is defined by chromatic number at most 2", false, "the definition compares clique and chromatic numbers"},
		{"Every graph is perfect", false, "odd holes are counterexamples"},
		{"Every odd hole is perfect", false, "an odd hole needs 3 colours with clique number 2"},
		{"Perfect graphs are exactly the planar graphs", false, "the classes differ"},
		{"Chordal graphs are never perfect", false, "they always are"},
		{"Bipartite graphs are not perfect", false, "they are perfect"},
		{"The complement of a perfect graph is never perfect", false, "it always is"},
		{"Perfectness is not inherited by induced subgraphs", false, "the definition requires every induced subgraph"},
	})
}

func kuratowskiProblem() generator.Problem {
	return countGraphVerdicts("How many of these statements about Kuratowski's theorem are true?", []graphVerdict{
		{"A finite graph is planar exactly when it has no subdivision of K_5 or K_{3,3}", true, "this is Kuratowski's theorem"},
		{"A subdivision replaces an edge by a path", true, "new degree-two vertices are inserted"},
		{"A subdivision of K_{3,3} is nonplanar", true, "subdivision preserves planarity"},
		{"Wagner's theorem uses K_5 and K_{3,3} minors", true, "the minor version of the same obstruction"},
		{"K_{3,3} is nonplanar", true, "it is one of the forbidden graphs"},
		{"Every tree is planar", true, "trees have no cycles to obstruct embedding"},
		{"Every cycle is planar", true, "a cycle bounds a disc"},
		{"Subdivision preserves planarity", true, "inserting degree-two vertices changes nothing"},
		{"Kuratowski's theorem forbids K_4 and C_5", false, "both are planar"},
		{"A subdivision replaces an edge by a triangle", false, "it replaces it by a path"},
		{"A subdivision of K_{3,3} is planar", false, "subdivision preserves nonplanarity"},
		{"Wagner's theorem uses K_4 and C_5", false, "it uses K_5 and K_{3,3}"},
		{"K_5 is planar", false, "it is the other forbidden graph"},
		{"Every graph on five vertices is planar", false, "K_5 is a counterexample"},
		{"Subdivision can destroy planarity", false, "it preserves it"},
		{"Every subgraph of K_{3,3} is planar", false, "K_{3,3} itself is a subgraph"},
	})
}

func minorProblem() generator.Problem {
	return countGraphVerdicts("How many of these statements about graph minors are true?", []graphVerdict{
		{"A minor is obtained by deleting vertices and edges and contracting edges", true, "those are the three operations"},
		{"Subdivision is not a minor operation", true, "it adds vertices rather than removing them"},
		{"Planarity is closed under taking minors", true, "a minor of a planar graph is planar"},
		{"Hadwiger's conjecture says chromatic number k forces a K_k minor", true, "it is open in general"},
		{"A graph with a K_5 minor is nonplanar", true, "K_5 is nonplanar and planarity is minor-closed"},
		{"Every subgraph is a minor", true, "delete the edges outside it"},
		{"Contracting edges can create a K_5 minor", true, "contraction can merge vertices"},
		{"The minor relation is transitive", true, "a minor of a minor is a minor"},
		{"A minor is obtained by subdivision only", false, "subdivision is not one of the operations"},
		{"A minor is obtained by adding edges", false, "minors never add edges"},
		{"Every minor is a subgraph", false, "contraction can create edges that were not present"},
		{"A graph with a K_5 minor is planar", false, "K_5 is a forbidden minor"},
		{"Planarity is not closed under minors", false, "it is minor-closed"},
		{"Hadwiger's conjecture is about independent sets", false, "it is about complete minors"},
		{"Every graph with a K_3 minor is nonplanar", false, "K_3 is planar"},
		{"Vertex deletion is not allowed in a minor", false, "it is one of the operations"},
	})
}

func randomThresholdProblem() generator.Problem {
	return countGraphVerdicts("How many of these statements about random graph thresholds are true?", []graphVerdict{
		{"In G(n,p) every edge appears independently with probability p", true, "this is the Erdős-Rényi model"},
		{"The connectivity threshold is log n / n", true, "above it the graph is connected with high probability"},
		{"The giant component threshold is 1/n", true, "a linear-order component appears above it"},
		{"Above 1/n the largest component has linear order", true, "the giant component emerges"},
		{"Below log n / n there are isolated vertices with high probability", true, "connectivity fails"},
		{"At p = 0 the graph is empty", true, "no edge appears"},
		{"At p = 1 the graph is complete", true, "every edge appears"},
		{"The expected number of edges is C(n,2) p", true, "sum the indicator expectations"},
		{"The model is defined by edge probability 1/n", false, "p is a parameter"},
		{"The connectivity threshold is 1/n", false, "it is log n / n"},
		{"The giant component threshold is log n / n", false, "it is 1/n"},
		{"Above 1/n the largest component stays constant", false, "it grows linearly"},
		{"Below 1/n the graph is always connected", false, "it is usually disconnected"},
		{"At p = 1 the graph is empty", false, "it is complete"},
		{"The expected number of edges is n p^2", false, "it is C(n,2) p"},
		{"Edges in G(n,p) are dependent", false, "they are independent by definition"},
	})
}

func strongConnectivityProblem() generator.Problem {
	return countGraphVerdicts("How many of these statements about strong connectivity are true?", []graphVerdict{
		{"A digraph is strongly connected when every vertex reaches every other", true, "this is the definition"},
		{"The condensation of a digraph is a DAG", true, "a cycle of SCCs would merge them"},
		{"The vertices of the condensation are the strongly connected components", true, "each component becomes one vertex"},
		{"Reversing every edge preserves strong connectivity", true, "paths reverse with the edges"},
		{"A digraph has more than one SCC exactly when its condensation has more than one vertex", true, "the condensation records the components"},
		{"A directed cycle is strongly connected", true, "follow the cycle in either direction"},
		{"A single vertex is strongly connected", true, "it reaches itself by the empty path"},
		{"A strongly connected digraph has a connected underlying graph", true, "reachability implies weak connectivity"},
		{"A connected underlying graph is enough for strong connectivity", false, "a one-way path is a counterexample"},
		{"The condensation may contain a directed cycle", false, "it is always a DAG"},
		{"Reversing all edges can destroy strong connectivity", false, "it preserves it"},
		{"A directed path is strongly connected", false, "the last vertex cannot reach the first"},
		{"Positive outdegree at every vertex is enough for strong connectivity", false, "two disjoint cycles are a counterexample"},
		{"Every DAG is strongly connected", false, "a nontrivial DAG has no directed cycle"},
		{"The condensation is always a tree", false, "it is a DAG with no such restriction"},
		{"Edge directions can be ignored when testing strong connectivity", false, "directions define reachability"},
	})
}

func topologicalSortProblem() generator.Problem {
	return countGraphVerdicts("How many of these statements about topological orderings are true?", []graphVerdict{
		{"A topological ordering exists exactly when the digraph is acyclic", true, "a cycle admits no linear order"},
		{"Every edge points from an earlier vertex to a later one", true, "this is the defining property"},
		{"Kahn's algorithm repeatedly removes a vertex of in-degree zero", true, "that is the peeling process"},
		{"A topological ordering exists exactly when there is no directed cycle", true, "the two conditions are equivalent"},
		{"Ties between available vertices give different valid orders", true, "any zero in-degree choice works"},
		{"A directed path has exactly one topological ordering", true, "the path fixes the order"},
		{"Every finite DAG has a source", true, "follow incoming edges backwards"},
		{"A topological order lists every vertex", true, "it is a linear order of the vertex set"},
		{"A topological ordering exists exactly when the graph is connected", false, "acyclicity is the condition"},
		{"Edges point from later vertices to earlier ones", false, "that reverses the definition"},
		{"Kahn's algorithm removes a vertex of maximum degree", false, "it removes in-degree zero"},
		{"A digraph with a directed cycle has a topological ordering", false, "the cycle cannot be linearised"},
		{"A topological ordering is always unique", false, "independent vertices can be ordered freely"},
		{"A DAG with a unique source has a unique topological ordering", false, "branches still allow choices"},
		{"A topological ordering may omit a vertex", false, "it lists the whole vertex set"},
		{"Every digraph has a topological ordering", false, "cyclic digraphs do not"},
	})
}
