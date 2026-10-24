package graphtheory

import (
	"fmt"
	"math/rand"
)

// graphClaim is the shape the graph-theory generators use for a recognition
// question. It is declared once at package level so a computed claim can be
// appended from another file; the generators previously declared their own
// local copy, which made the pools unreachable from anywhere else.
type graphClaim struct {
	question    string
	correct     string
	distractors []string
	reason      string
}

// graphDistinct drops candidates that repeat the answer or each other and tops
// the list up to three, because generator.MultipleChoice deduplicates its
// options and a collision silently costs the question one choice.
func graphDistinct(correct string, candidates ...string) []string {
	seen := map[string]bool{correct: true}
	out := make([]string, 0, 3)
	for _, c := range candidates {
		if seen[c] {
			continue
		}
		seen[c] = true
		out = append(out, c)
	}
	pool := []string{"0", "1", "2", "none of these"}
	for _, c := range pool {
		if len(out) == 3 {
			break
		}
		if seen[c] {
			continue
		}
		seen[c] = true
		out = append(out, c)
	}
	return out
}

// computedGraphClaims returns parameterised claims for the concept id, so the
// learner meets a drawn graph rather than one of the four sentences every
// textbook prints. Concepts without computed claims return nil and keep their
// written pool unchanged.
func computedGraphClaims(id string) []graphClaim {
	switch id {
	case "graph.brooks":
		n := 3 + rand.Intn(8)
		even := n%2 == 0
		chi := "3"
		why := fmt.Sprintf("C_%d is an odd cycle, which is one of the two exceptions, so it needs 3 colours.", n)
		if even {
			chi = "2"
			why = fmt.Sprintf("C_%d is an even cycle, so two colours alternate around it.", n)
		}
		k := 3 + rand.Intn(6)
		return []graphClaim{
			{
				question:    fmt.Sprintf("What is the chromatic number of the cycle C_%d?", n),
				correct:     chi,
				distractors: graphDistinct(chi, "1", "4", fmt.Sprintf("%d", n)),
				reason:      why,
			},
			{
				question:    fmt.Sprintf("How many colours does the complete graph K_%d need?", k),
				correct:     fmt.Sprintf("%d", k),
				distractors: graphDistinct(fmt.Sprintf("%d", k), fmt.Sprintf("%d", k-1), "2", "3"),
				reason:      fmt.Sprintf("Every pair of the %d vertices is adjacent, so all %d need different colours.", k, k),
			},
		}
	case "graph.subgraph":
		n := 4 + rand.Intn(7)
		edges := n * (n - 1) / 2
		v := 5 + rand.Intn(9)
		return []graphClaim{
			{
				question:    fmt.Sprintf("How many edges does the complete graph K_%d have?", n),
				correct:     fmt.Sprintf("%d", edges),
				distractors: graphDistinct(fmt.Sprintf("%d", edges), fmt.Sprintf("%d", n*(n-1)), fmt.Sprintf("%d", edges-n), fmt.Sprintf("%d", n*n)),
				reason:      fmt.Sprintf("Each of the %d vertices joins the other %d, and every edge is counted twice, so %d * %d / 2 = %d.", n, n-1, n, n-1, edges),
			},
			{
				question:    fmt.Sprintf("A tree on %d vertices has how many edges?", v),
				correct:     fmt.Sprintf("%d", v-1),
				distractors: graphDistinct(fmt.Sprintf("%d", v-1), fmt.Sprintf("%d", v), fmt.Sprintf("%d", v*(v-1)/2), fmt.Sprintf("%d", v+1)),
				reason:      fmt.Sprintf("A tree is connected with no cycle, so it keeps one fewer edge than vertices: %d - 1 = %d.", v, v-1),
			},
		}
	case "graph.random_threshold":
		n := 20 + rand.Intn(30)
		// The threshold for a giant component is 1/n; for connectivity it is
		// log(n)/n. Both are reported as a count of expected edges per vertex.
		pGiant := fmt.Sprintf("1/%d", n)
		pConn := fmt.Sprintf("ln(%d)/%d", n, n)
		return []graphClaim{
			{
				question:    fmt.Sprintf("In G(%d, p), which value of p is the threshold for the emergence of a giant component?", n),
				correct:     pGiant,
				distractors: graphDistinct(pGiant, pConn, "1/2", fmt.Sprintf("%d/%d", n, n)),
				reason:      fmt.Sprintf("The giant component appears around p = 1/n, here 1/%d; connectivity needs the larger ln(n)/n.", n),
			},
			{
				question:    fmt.Sprintf("In G(%d, p), which value of p is the threshold for the graph to become connected?", n),
				correct:     pConn,
				distractors: graphDistinct(pConn, pGiant, "1/2", fmt.Sprintf("%d/ln(%d)", n, n)),
				reason:      fmt.Sprintf("Isolated vertices disappear only above p = ln(n)/n, which is ln(%d)/%d here.", n, n),
			},
		}
	case "graph.kuratowski":
		n := 4 + rand.Intn(6)
		m := 3 + rand.Intn(4)
		return []graphClaim{
			{question: fmt.Sprintf("Does K_%d contain a subdivision of K_5?", n), correct: yn(n >= 5), distractors: graphDistinct(yn(n >= 5), yn(n < 5), "only when n is even"), reason: fmt.Sprintf("K_%d contains K_5 as a subgraph exactly when n is at least 5, and a subgraph is a subdivision of length zero.", n)},
			{question: fmt.Sprintf("Does K_3,3 appear as a subgraph of K_%d,%d?", n, m), correct: yn(n >= 3 && m >= 3), distractors: graphDistinct(yn(n >= 3 && m >= 3), yn(n < 3 || m < 3), "only when n equals m"), reason: "K_3,3 needs at least three vertices on each side."},
			{question: "Which pair of graphs does Kuratowski's theorem name as the obstructions to planarity?", correct: "K_5 and K_3,3", distractors: graphDistinct("K_5 and K_3,3", "K_4 and K_3,3", "K_5 and K_4"), reason: "A graph is planar exactly when it has no subdivision of K_5 or of K_3,3."},
		}
	case "graph.menger":
		k := 2 + rand.Intn(5)
		return []graphClaim{
			{question: fmt.Sprintf("If every pair of non-adjacent vertices in a graph is joined by %d internally disjoint paths, what is the vertex connectivity at least?", k), correct: fmt.Sprintf("%d", k), distractors: graphDistinct(fmt.Sprintf("%d", k), fmt.Sprintf("%d", k-1), fmt.Sprintf("%d", k+1)), reason: "Menger's theorem equates the maximum number of internally disjoint paths with the size of a smallest separating set."},
			{question: "Menger's theorem relates connectivity to which quantity?", correct: "the maximum number of internally disjoint paths", distractors: graphDistinct("the maximum number of internally disjoint paths", "the number of edges", "the chromatic number"), reason: "The min-cut equals the max number of disjoint paths."},
			{question: "For a graph with minimum degree d, what is the edge connectivity at most?", correct: "d", distractors: graphDistinct("d", "d+1", "2d"), reason: "Isolating one vertex cuts exactly its d incident edges, so the edge connectivity cannot exceed the minimum degree."},
		}
	case "graph.minor":
		n := 4 + rand.Intn(6)
		return []graphClaim{
			{question: fmt.Sprintf("How many edges does K_%d have?", n), correct: fmt.Sprintf("%d", n*(n-1)/2), distractors: graphDistinct(fmt.Sprintf("%d", n*(n-1)/2), fmt.Sprintf("%d", n*(n-1)), fmt.Sprintf("%d", n*n)), reason: fmt.Sprintf("Each of the %d vertices joins the other %d, and each edge is counted twice.", n, n-1)},
			{question: "A graph has a K_5 minor exactly when it can be reduced to K_5 by which operations?", correct: "deleting edges, deleting vertices and contracting edges", distractors: graphDistinct("deleting edges, deleting vertices and contracting edges", "subdividing edges", "adding vertices only"), reason: "Minors allow deletion and contraction, not subdivision."},
			{question: "Every planar graph excludes which two minors?", correct: "K_5 and K_3,3", distractors: graphDistinct("K_5 and K_3,3", "K_4 and C_5", "K_6 and K_3,3"), reason: "Wagner's form of Kuratowski's theorem states planarity as the exclusion of these two minors."},
		}
	case "graph.perfect_graph":
		n := 3 + rand.Intn(5)
		return []graphClaim{
			{question: fmt.Sprintf("Is the complete bipartite graph K_%d,%d perfect?", n, n+1), correct: "yes", distractors: graphDistinct("yes", "no", "only when both parts are equal"), reason: "Bipartite graphs are perfect: their chromatic number equals their clique number, which is 2."},
			{question: "Which induced subgraph makes an odd cycle of length five not perfect?", correct: "the cycle itself, whose chromatic number 3 exceeds its clique number 2", distractors: graphDistinct("the cycle itself, whose chromatic number 3 exceeds its clique number 2", "its complement", "its spanning tree"), reason: "An odd hole is the standard obstruction to perfection."},
			{question: "Is the complement of a perfect graph perfect?", correct: "yes", distractors: graphDistinct("yes", "no", "only for bipartite graphs"), reason: "The strong perfect graph theorem gives perfection to the complement as well."},
		}
	case "graph.perfect_matching":
		k := 2 + rand.Intn(5)
		return []graphClaim{
			{question: fmt.Sprintf("How many edges does a perfect matching of K_%d have?", 2*k), correct: fmt.Sprintf("%d", k), distractors: graphDistinct(fmt.Sprintf("%d", k), fmt.Sprintf("%d", 2*k), fmt.Sprintf("%d", k+1)), reason: fmt.Sprintf("The %d vertices pair up, so there are %d edges.", 2*k, k)},
			{question: "Which graph has no perfect matching?", correct: "a path with an odd number of vertices", distractors: graphDistinct("a path with an odd number of vertices", "a cycle of even length", "K_4"), reason: "A perfect matching covers every vertex, which an odd vertex count forbids."},
			{question: "By Tutte's theorem, a perfect matching exists exactly when which condition holds?", correct: "every vertex set S leaves at most |S| odd components", distractors: graphDistinct("every vertex set S leaves at most |S| odd components", "the graph is bipartite", "every vertex has even degree"), reason: "Tutte's condition counts odd components left after deleting S."},
		}
	case "graph.strong_connectivity":
		n := 3 + rand.Intn(6)
		return []graphClaim{
			{question: fmt.Sprintf("Is the cycle C_%d strongly connected?", n), correct: "yes", distractors: graphDistinct("yes", "no", "only when n is even"), reason: "Following the cycle in one direction reaches every vertex from every other."},
			{question: "Is a path with three or more vertices strongly connected?", correct: "no", distractors: graphDistinct("no", "yes", "only when it has an even number of edges"), reason: "The end vertices cannot be reached from the other end, since the edges point one way."},
			{question: fmt.Sprintf("How many edges does the smallest strongly connected digraph on %d vertices have?", n), correct: fmt.Sprintf("%d", n), distractors: graphDistinct(fmt.Sprintf("%d", n), fmt.Sprintf("%d", n-1), fmt.Sprintf("%d", 2*n)), reason: "A directed cycle on all vertices is minimal, and it has one edge per vertex."},
		}
	case "graph.topological_sort":
		n := 3 + rand.Intn(7)
		return []graphClaim{
			{question: "Which graph always admits a topological ordering?", correct: "a directed acyclic graph", distractors: graphDistinct("a directed acyclic graph", "any directed graph", "any connected graph"), reason: "A cycle has no linear order consistent with all its edges."},
			{question: fmt.Sprintf("How many topological orderings does a directed path on %d vertices have?", n), correct: "1", distractors: graphDistinct("1", fmt.Sprintf("%d", n), fmt.Sprintf("%d", n*(n-1)/2)), reason: "Every vertex except the first has a predecessor, so the order is forced."},
			{question: "Every finite DAG has at least one vertex of which kind?", correct: "a source with no incoming edges", distractors: graphDistinct("a source with no incoming edges", "a vertex of degree zero", "a vertex with out-degree one"), reason: "Following incoming edges backwards must stop, and it stops at a source."},
		}
	case "graph.homomorphism":
		n := 3 + rand.Intn(5)
		return []graphClaim{
			{question: fmt.Sprintf("Does K_%d admit a homomorphism to K_%d?", n+1, n), correct: "yes", distractors: graphDistinct("yes", "no", "only when the graphs are isomorphic"), reason: "A homomorphism may identify adjacent vertices, so a larger complete graph maps onto a smaller one."},
			{question: "A graph admits a homomorphism to K_2 exactly when it is...", correct: "bipartite", distractors: graphDistinct("bipartite", "planar", "regular"), reason: "A two-colouring is precisely a homomorphism to K_2."},
			{question: "Is every graph homomorphism injective on vertices?", correct: "no", distractors: graphDistinct("no", "yes", "only between complete graphs"), reason: "Homomorphisms preserve edges, and they may collapse distinct vertices."},
		}
	case "graph.isomorphism":
		n := 3 + rand.Intn(6)
		return []graphClaim{
			{question: fmt.Sprintf("Are K_%d and K_%d isomorphic when the numbers differ?", n, n), correct: "yes", distractors: graphDistinct("yes", "no", "only when n is even"), reason: "Complete graphs of the same order are determined by their order."},
			{question: fmt.Sprintf("For which m is C_%d isomorphic to C_m?", n), correct: fmt.Sprintf("only m = %d", n), distractors: graphDistinct(fmt.Sprintf("only m = %d", n), "every m of the same parity", "no m"), reason: "Cycles are determined by their length up to rotation and reflection."},
			{question: "Two trees with the same degree sequence are...", correct: "not necessarily isomorphic", distractors: graphDistinct("not necessarily isomorphic", "always isomorphic", "isomorphic only when both are paths"), reason: "The degree sequence does not record how the branches are arranged."},
		}
	}
	return nil
}

// yn renders a verdict as the answer string the pools use.
func yn(ok bool) string {
	if ok {
		return "yes"
	}
	return "no"
}
