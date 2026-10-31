package topology

// The 2026-10-18 topology batch: twelve concepts that take the domain from 30
// to the target of 42.
//
// The survey found that the domain jumps from three basics straight to named
// theorems (Urysohn, Tietze, Stone-Cech, Van Kampen, de Rham) with almost no
// computational middle, and that 29 of its 30 concepts are multiple choice.
// The batch therefore builds the metric-space toolkit the later lessons assume
// (balls, interior and boundary, limit points, convergence, completeness,
// bases, subspaces, dense sets) and adds the computable invariants the
// algebraic-topology lessons describe but never name (homeomorphism, Euler
// characteristic, genus, Betti numbers).
//
// Six of the twelve are numeric, and every numeric generator publishes the
// facts its answer follows from, so the tests recompute the answer instead of
// comparing a string with itself. The six recognition concepts use
// generator.MultipleChoice, so every distractor is a real option the router
// must reject. The new `topology.invariants` subdomain holds the four invariant
// concepts; the others join the existing `basics` and `constructions` homes.

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerBatch(reg *generator.Registry) {
	reg.Register("topo.metric.ball", &ballGen{})
	reg.Register("topo.metric.interior", &interiorGen{})
	reg.Register("topo.metric.limit_point", &limitPointGen{})
	reg.Register("topo.metric.convergence", &convergenceGen{})
	reg.Register("topo.metric.complete", &completeGen{})
	reg.Register("topo.metric.basis", &basisGen{})
	reg.Register("topo.metric.subspace", &subspaceGen{})
	reg.Register("topo.metric.dense", &denseGen{})
	reg.Register("topo.invariant.homeomorphism", &homeomorphismGen{})
	reg.Register("topo.invariant.euler", &eulerGen{})
	reg.Register("topo.invariant.genus", &genusGen{})
	reg.Register("topo.invariant.betti", &bettiGen{})
}

// ---------------------------------------------------------------------------
// 1. Open balls (numeric)
// ---------------------------------------------------------------------------

type ballGen struct{}

func (g *ballGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	if rand.Intn(2) == 0 {
		triples := [][3]int{
			{3, 4, 5}, {5, 12, 13}, {8, 15, 17}, {7, 24, 25}, {6, 8, 10},
		}
		for attempt := 0; attempt < 50; attempt++ {
			tr := triples[rand.Intn(len(triples))]
			dx, dy, dist := tr[0], tr[1], tr[2]
			x1 := rand.Intn(15) - 7
			y1 := rand.Intn(15) - 7
			x2 := x1 + dx
			y2 := y1 + dy
			if dist == x1 || dist == y1 || dist == x2 || dist == y2 ||
				-dist == x1 || -dist == y1 || -dist == x2 || -dist == y2 {
				continue
			}
			ans := strconv.Itoa(dist)
			q := fmt.Sprintf("In the Euclidean plane, what is the distance d((%d,%d),(%d,%d))? (enter a number)", x1, y1, x2, y2)
			reason := fmt.Sprintf("The differences are %d and %d, so the distance is sqrt(%d^2 + %d^2) = %d.", dx, dy, dx, dy, dist)
			return generator.Problem{
				Question:    q,
				Answer:      ans,
				Explanation: reason,
				Facts: map[string]string{
					"kind":   "distance",
					"x1":     strconv.Itoa(x1),
					"y1":     strconv.Itoa(y1),
					"x2":     strconv.Itoa(x2),
					"y2":     strconv.Itoa(y2),
					"answer": ans,
				},
			}
		}
	}
	for attempt := 0; attempt < 50; attempt++ {
		center := rand.Intn(30) - 15
		radius := 1 + rand.Intn(15)
		side := "lower"
		var val int
		var boundName string
		if rand.Intn(2) == 1 {
			side = "upper"
			val = center + radius
			boundName = "b"
		} else {
			val = center - radius
			boundName = "a"
		}
		if val == center || val == radius || val == -center || val == -radius {
			continue
		}
		ans := strconv.Itoa(val)
		q := fmt.Sprintf("In R with the usual metric, the open ball of radius %d centred at %d is the interval (a, b). What is %s? (enter a number)", radius, center, boundName)
		reason := fmt.Sprintf("The ball is (%d - %d, %d + %d) = (%d, %d), so the %s endpoint is %d.", center, radius, center, radius, center-radius, center+radius, side, val)
		return generator.Problem{
			Question:    q,
			Answer:      ans,
			Explanation: reason,
			Facts: map[string]string{
				"kind":   "endpoint",
				"center": strconv.Itoa(center),
				"radius": strconv.Itoa(radius),
				"side":   side,
				"answer": ans,
			},
		}
	}
	return generator.Problem{
		Question:    "In R with the usual metric, the open ball of radius 3 centred at 5 is the interval (a, b). What is b? (enter a number)",
		Answer:      "8",
		Explanation: "5 + 3 = 8.",
		Facts: map[string]string{
			"kind":   "endpoint",
			"center": "5",
			"radius": "3",
			"side":   "upper",
			"answer": "8",
		},
	}
}

// ---------------------------------------------------------------------------
// 2. Interior, closure and boundary (numeric)
// ---------------------------------------------------------------------------

type interiorGen struct{}

func (g *interiorGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	for attempt := 0; attempt < 50; attempt++ {
		mode := rand.Intn(5)
		var q, ans, reason string
		var pts []int
		kind := "boundary"
		base := rand.Intn(20) - 10

		switch mode {
		case 0:
			// Disjoint closed intervals: 1 to 6 intervals => 2, 4, 6, 8, 10, 12 boundary points
			numIntervals := 1 + rand.Intn(6)
			ansVal := 2 * numIntervals
			ans = strconv.Itoa(ansVal)
			intervals := make([]string, numIntervals)
			cur := base
			pts = make([]int, 0, ansVal)
			leak := false
			for i := 0; i < numIntervals; i++ {
				a := cur
				b := a + rand.Intn(3) + 1
				cur = b + rand.Intn(3) + 2
				intervals[i] = fmt.Sprintf("[%d, %d]", a, b)
				pts = append(pts, a, b)
				if a == ansVal || b == ansVal || -a == ansVal || -b == ansVal {
					leak = true
				}
			}
			if leak {
				continue
			}
			q = fmt.Sprintf("In R with the usual topology, how many boundary points does the set %s have? (enter a number)", strings.Join(intervals, " union "))
			reason = fmt.Sprintf("The boundary consists of the %d endpoints of the pairwise disjoint closed intervals.", ansVal)
		case 1:
			// Adjacent intervals missing division points: 2, 3, 4 intervals touching => 3, 4, 5 boundary points
			m := 2 + rand.Intn(4)
			ansVal := m + 1
			ans = strconv.Itoa(ansVal)
			intervals := make([]string, m)
			cur := base
			pts = make([]int, 0, ansVal)
			leak := false
			for i := 0; i < m; i++ {
				next := cur + rand.Intn(4) + 1
				intervals[i] = fmt.Sprintf("(%d, %d)", cur, next)
				if i == 0 {
					pts = append(pts, cur)
					if cur == ansVal || -cur == ansVal {
						leak = true
					}
				}
				pts = append(pts, next)
				if next == ansVal || -next == ansVal {
					leak = true
				}
				cur = next
			}
			if leak {
				continue
			}
			q = fmt.Sprintf("In R with the usual topology, how many boundary points does the set %s have? (enter a number)", strings.Join(intervals, " union "))
			reason = fmt.Sprintf("The boundary consists of the %d outer and division boundary points.", ansVal)
		case 2:
			// Ray (a, infinity) or (-infinity, a) has 1 boundary point
			a := base
			if a == 1 || a == -1 {
				continue
			}
			pts = []int{a}
			ans = "1"
			if rand.Intn(2) == 0 {
				q = fmt.Sprintf("In R with the usual topology, how many boundary points does the ray (%d, infinity) have? (enter a number)", a)
			} else {
				q = fmt.Sprintf("In R with the usual topology, how many boundary points does the ray (-infinity, %d) have? (enter a number)", a)
			}
			reason = fmt.Sprintf("The only boundary point is %d.", a)
		case 3:
			// Closed ray [a, infinity) has 1 boundary point
			a := base
			if a == 1 || a == -1 {
				continue
			}
			pts = []int{a}
			ans = "1"
			q = fmt.Sprintf("In R with the usual topology, how many boundary points does the closed ray [%d, infinity) have? (enter a number)", a)
			reason = fmt.Sprintf("The only boundary point is %d.", a)
		default:
			// Finite set of isolated points has 0 interior points
			kind = "interior"
			pts = []int{}
			ans = "0"
			n := 2 + rand.Intn(4)
			elems := make([]string, n)
			for i := 0; i < n; i++ {
				elems[i] = strconv.Itoa(base + i*2)
			}
			q = fmt.Sprintf("In R with the usual topology, how many interior points does the set {%s} have? (enter a number)", strings.Join(elems, ", "))
			reason = "A finite set of isolated points has empty interior: no ball around any point stays inside the set."
		}

		parts := make([]string, len(pts))
		for i, p := range pts {
			parts[i] = strconv.Itoa(p)
		}
		return generator.Problem{
			Question:    q,
			Answer:      ans,
			Explanation: reason,
			Facts: map[string]string{
				"points": strings.Join(parts, ","), "kind": kind, "answer": ans,
			},
		}
	}
	return generator.Problem{
		Question:    "In R with the usual topology, how many boundary points does the set [1, 2] union [4, 5] have? (enter a number)",
		Answer:      "4",
		Explanation: "The four endpoints are boundary points.",
		Facts: map[string]string{
			"points": "1,2,4,5", "kind": "boundary", "answer": "4",
		},
	}
}

// ---------------------------------------------------------------------------
// 3. Limit points
// ---------------------------------------------------------------------------

type limitPointGen struct{}

func (g *limitPointGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return limitPointProblem()
}

// ---------------------------------------------------------------------------
// 4. Convergence (numeric)
// ---------------------------------------------------------------------------

type convergenceGen struct{}

func (g *convergenceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	cPool := []int{2, 4, 5, 8, 10}
	for attempt := 0; attempt < 50; attempt++ {
		c := cPool[rand.Intn(len(cPool))]
		a := rand.Intn(16) + 1
		b := rand.Intn(15) + 1
		if rand.Intn(2) == 1 {
			b = -b
		}
		d := rand.Intn(10) + 1

		var ans string
		if a%c == 0 {
			ans = strconv.Itoa(a / c)
		} else {
			ans = strconv.FormatFloat(float64(a)/float64(c), 'f', -1, 64)
		}

		if ans == strconv.Itoa(a) || ans == strconv.Itoa(c) ||
			ans == strconv.Itoa(b) || ans == strconv.Itoa(-b) ||
			ans == strconv.Itoa(d) {
			continue
		}

		var numStr string
		if b > 0 {
			numStr = fmt.Sprintf("%dn + %d", a, b)
		} else {
			numStr = fmt.Sprintf("%dn - %d", a, -b)
		}
		denStr := fmt.Sprintf("%dn + %d", c, d)

		q := fmt.Sprintf("The sequence a_n = (%s)/(%s) converges to what limit? (enter a number)", numStr, denStr)
		reason := fmt.Sprintf("Divide numerator and denominator by n: the leading term dominates, leaving the limit %d/%d = %s.", a, c, ans)
		return generator.Problem{
			Question:    q,
			Answer:      ans,
			Explanation: reason,
			Facts: map[string]string{
				"a": strconv.Itoa(a), "b": strconv.Itoa(b),
				"c": strconv.Itoa(c), "d": strconv.Itoa(d), "answer": ans,
			},
		}
	}
	return generator.Problem{
		Question:    "The sequence a_n = (6n + 1)/(2n + 3) converges to what limit? (enter a number)",
		Answer:      "3",
		Explanation: "The leading coefficients give 6/2 = 3.",
		Facts: map[string]string{
			"a": "6", "b": "1", "c": "2", "d": "3", "answer": "3",
		},
	}
}

// ---------------------------------------------------------------------------
// 5. Completeness and Cauchy sequences
// ---------------------------------------------------------------------------

type completeGen struct{}

func (g *completeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return completeMetricProblem()
}

// ---------------------------------------------------------------------------
// 6. Bases and subbases
// ---------------------------------------------------------------------------

type basisGen struct{}

func (g *basisGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return basisProblem()
}

// ---------------------------------------------------------------------------
// 7. The subspace topology
// ---------------------------------------------------------------------------

type subspaceGen struct{}

func (g *subspaceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return subspaceProblem()
}

// ---------------------------------------------------------------------------
// 8. Dense sets and separability
// ---------------------------------------------------------------------------

type denseGen struct{}

func (g *denseGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return denseProblem()
}

// ---------------------------------------------------------------------------
// 9. Homeomorphisms
// ---------------------------------------------------------------------------

type homeomorphismGen struct{}

func (g *homeomorphismGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return homeomorphismProblem()
}

// ---------------------------------------------------------------------------
// 10. Euler characteristic (numeric)
// ---------------------------------------------------------------------------

type eulerGen struct{}

func (g *eulerGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(3)
	var q, ans, reason string
	var v, e, f int
	switch mode {
	case 0:
		n := 3 + rand.Intn(10)
		if rand.Intn(2) == 0 {
			v = n + 1
			e = 2 * n
			f = n + 1
			ans = "2"
			q = fmt.Sprintf("A %d-gonal pyramid has %d vertices, %d edges and %d faces. What is its Euler characteristic V - E + F? (enter a number)", n, v, e, f)
			reason = fmt.Sprintf("%d - %d + %d = 2, the Euler characteristic of a sphere.", v, e, f)
		} else {
			v = 2 * n
			e = 3 * n
			f = n + 2
			ans = "2"
			q = fmt.Sprintf("A %d-gonal prism has %d vertices, %d edges and %d faces. What is its Euler characteristic V - E + F? (enter a number)", n, v, e, f)
			reason = fmt.Sprintf("%d - %d + %d = 2, as for every polyhedron homeomorphic to a sphere.", v, e, f)
		}
	case 1:
		genus := 1 + rand.Intn(6)
		chi := 2 - 2*genus
		v = 6 + rand.Intn(10)
		f = 2 * (v - chi)
		e = v + f - chi
		ans = strconv.Itoa(chi)
		q = fmt.Sprintf("A triangulated closed surface of genus %d has %d vertices, %d edges and %d faces. What is its Euler characteristic V - E + F? (enter a number)", genus, v, e, f)
		reason = fmt.Sprintf("%d - %d + %d = %d, the Euler characteristic of a surface of genus %d.", v, e, f, chi, genus)
	default:
		chiChoice := []int{0, 1}[rand.Intn(2)]
		v = 6 + rand.Intn(10)
		f = 6 + rand.Intn(10)
		e = v + f - chiChoice
		ans = strconv.Itoa(chiChoice)
		surfName := "torus"
		if chiChoice == 1 {
			surfName = "real projective plane"
		}
		q = fmt.Sprintf("A 2-dimensional cell complex homeomorphic to the %s has %d vertices, %d edges and %d faces. What is its Euler characteristic V - E + F? (enter a number)", surfName, v, e, f)
		reason = fmt.Sprintf("%d - %d + %d = %d.", v, e, f, chiChoice)
	}
	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: reason,
		Facts: map[string]string{
			"v": strconv.Itoa(v), "e": strconv.Itoa(e), "f": strconv.Itoa(f),
			"answer": ans,
		},
	}
}

// ---------------------------------------------------------------------------
// 11. Genus of a surface (numeric)
// ---------------------------------------------------------------------------

type genusGen struct{}

func (g *genusGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kind := "genus"
	if rand.Intn(2) == 1 {
		kind = "chi"
	}
	genus := rand.Intn(12)
	chi := 2 - 2*genus
	if kind == "genus" {
		q := fmt.Sprintf("A closed orientable surface has Euler characteristic %d. What is its genus? (enter a number)", chi)
		ans := strconv.Itoa(genus)
		reason := fmt.Sprintf("For a closed orientable surface chi = 2 - 2g, so g = (2 - (%d))/2 = %d.", chi, genus)
		return generator.Problem{
			Question:    q,
			Answer:      ans,
			Explanation: reason,
			Facts: map[string]string{
				"kind":   "genus",
				"chi":    strconv.Itoa(chi),
				"answer": ans,
			},
		}
	}
	q := fmt.Sprintf("A closed orientable surface has genus %d. What is its Euler characteristic 2 - 2g? (enter a number)", genus)
	ans := strconv.Itoa(chi)
	reason := fmt.Sprintf("2 - 2(%d) = %d.", genus, chi)
	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: reason,
		Facts: map[string]string{
			"kind":   "chi",
			"genus":  strconv.Itoa(genus),
			"answer": ans,
		},
	}
}

// ---------------------------------------------------------------------------
// 12. Betti numbers (numeric)
// ---------------------------------------------------------------------------

type bettiGen struct{}

func (g *bettiGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(6)
	switch mode {
	case 0:
		return generator.Problem{
			Question:    "The first Betti number of a circle is what? (enter a number)",
			Answer:      "1",
			Explanation: "A circle has one independent loop, so H_1 has rank 1.",
			Facts: map[string]string{
				"kind":   "circle",
				"param":  "1",
				"answer": "1",
			},
		}
	case 1:
		genus := 1 + rand.Intn(10)
		b1 := 2 * genus
		ans := strconv.Itoa(b1)
		return generator.Problem{
			Question:    fmt.Sprintf("A closed orientable surface of genus %d has what first Betti number b1? (enter a number)", genus),
			Answer:      ans,
			Explanation: fmt.Sprintf("A closed orientable surface of genus g has 2g independent 1-cycles, so b1 = 2(%d) = %d.", genus, b1),
			Facts: map[string]string{
				"kind":   "surface",
				"param":  strconv.Itoa(genus),
				"answer": ans,
			},
		}
	case 2:
		numWords := []string{
			"zero", "one", "two", "three", "four", "five",
			"six", "seven", "eight", "nine", "ten", "eleven",
			"twelve", "thirteen", "fourteen", "fifteen", "sixteen",
		}
		k := 2 + rand.Intn(10)
		ans := strconv.Itoa(k)
		return generator.Problem{
			Question:    fmt.Sprintf("A wedge sum of %s circles has what first Betti number b1? (enter a number)", numWords[k]),
			Answer:      ans,
			Explanation: fmt.Sprintf("Each circle in a wedge sum contributes one independent generator to the first homology group, giving b1 = %d.", k),
			Facts: map[string]string{
				"kind":   "wedge",
				"param":  strconv.Itoa(k),
				"answer": ans,
			},
		}
	case 3:
		numWords := []string{
			"zero", "one", "two", "three", "four", "five",
			"six", "seven", "eight", "nine", "ten", "eleven",
			"twelve", "thirteen", "fourteen", "fifteen", "sixteen",
		}
		k := 2 + rand.Intn(10)
		ans := strconv.Itoa(k)
		return generator.Problem{
			Question:    fmt.Sprintf("A topological space has %s connected components. What is its zeroth Betti number b0? (enter a number)", numWords[k]),
			Answer:      ans,
			Explanation: fmt.Sprintf("The zeroth Betti number b0 counts path components, so %d components give b0 = %d.", k, k),
			Facts: map[string]string{
				"kind":   "components",
				"param":  strconv.Itoa(k),
				"answer": ans,
			},
		}
	case 4:
		return generator.Problem{
			Question:    "A tree is a connected graph with no cycles. What is its first Betti number b1? (enter a number)",
			Answer:      "0",
			Explanation: "A tree has no independent loop, so H_1 has rank 0.",
			Facts: map[string]string{
				"kind":   "tree",
				"param":  "0",
				"answer": "0",
			},
		}
	default:
		return generator.Problem{
			Question:    "For a sphere, b0 = 1, b1 = 0 and b2 = 1. What is b0 - b1 + b2? (enter a number)",
			Answer:      "2",
			Explanation: "1 - 0 + 1 = 2, which is also the Euler characteristic of the sphere.",
			Facts: map[string]string{
				"kind":   "sphere_euler",
				"param":  "2",
				"answer": "2",
			},
		}
	}
}
