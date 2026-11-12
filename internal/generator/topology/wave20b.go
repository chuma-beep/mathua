package topology

// Wave 20 topology, second batch: homology, knot invariants and geometric
// group theory. Each family publishes its parameters in Facts and returns a
// single integer.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave20Rest(reg *generator.Registry) {
	reg.Register("topo.hom.chain_complex", &w20Family{gen: w20HomChainComplex})
	reg.Register("topo.hom.exact_sequence", &w20Family{gen: w20HomExactSequence})
	reg.Register("topo.hom.snake_lemma", &w20Family{gen: w20HomSnakeLemma})
	reg.Register("topo.hom.five_lemma", &w20Family{gen: w20HomFiveLemma})
	reg.Register("topo.hom.kunneth", &w20Family{gen: w20HomKunneth})
	reg.Register("topo.hom.universal_coefficient", &w20Family{gen: w20HomUniversalCoefficient})
	reg.Register("topo.hom.euler_poincare", &w20Family{gen: w20HomEulerPoincare})
	reg.Register("topo.hom.sphere", &w20Family{gen: w20HomSphere})
	reg.Register("topo.hom.projective", &w20Family{gen: w20HomProjective})
	reg.Register("topo.hom.torus", &w20Family{gen: w20HomTorus})
	reg.Register("topo.knot.crossing", &w20Family{gen: w20KnotCrossing})
	reg.Register("topo.knot.unknotting", &w20Family{gen: w20KnotUnknotting})
	reg.Register("topo.knot.seifert_genus", &w20Family{gen: w20KnotSeifertGenus})
	reg.Register("topo.knot.alexander", &w20Family{gen: w20KnotAlexander})
	reg.Register("topo.knot.skein", &w20Family{gen: w20KnotSkein})
	reg.Register("topo.knot.torus_knot", &w20Family{gen: w20KnotTorusKnot})
	reg.Register("topo.knot.connected_sum", &w20Family{gen: w20KnotConnectedSum})
	reg.Register("topo.knot.bracket", &w20Family{gen: w20KnotBracket})
	reg.Register("topo.knot.writhe", &w20Family{gen: w20KnotWrithe})
	reg.Register("topo.knot.linking", &w20Family{gen: w20KnotLinking})
	reg.Register("topo.geo.cayley", &w20Family{gen: w20GeoCayley})
	reg.Register("topo.geo.growth", &w20Family{gen: w20GeoGrowth})
	reg.Register("topo.geo.quasi_isometry", &w20Family{gen: w20GeoQuasiIsometry})
	reg.Register("topo.geo.ends", &w20Family{gen: w20GeoEnds})
	reg.Register("topo.geo.svarc_milnor", &w20Family{gen: w20GeoSvarcMilnor})
	reg.Register("topo.geo.hyperbolic", &w20Family{gen: w20GeoHyperbolic})
	reg.Register("topo.geo.boundary", &w20Family{gen: w20GeoBoundary})
	reg.Register("topo.geo.dehn", &w20Family{gen: w20GeoDehn})
	reg.Register("topo.geo.amenable", &w20Family{gen: w20GeoAmenable})
	reg.Register("topo.geo.property_t", &w20Family{gen: w20GeoPropertyT})
}

func w20KnotPrimes() []int { return []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47} }

func w20KnotPair() (int, int) {
	primes := w20KnotPrimes()
	p := primes[rand.Intn(len(primes))]
	q := primes[rand.Intn(len(primes))]
	for q == p {
		q = primes[rand.Intn(len(primes))]
	}
	return p, q
}

func w20Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func w20Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ----- homology -----

func w20HomChainComplex(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		n := w20R(2, 30)
		w20I(f, "n", n)
		maps := n
		return w20T(kp, fmt.Sprintf("A chain complex has groups \\(C_0,\\dots,C_{%d}\\). How many boundary maps does it have?", n),
			strconv.Itoa(maps), fmt.Sprintf("There is one map \\(\\partial_k:C_k\\to C_{k-1}\\) for each \\(k=1,\\dots,%d\\), so %d maps.", n, maps), f)
	case 2:
		n := w20R(2, 8)
		w20I(f, "n", n)
		ranks := make([]int, n+1)
		sum := 0
		for k := 0; k <= n; k++ {
			ranks[k] = w20R(1, 6)
			w20I(f, fmt.Sprintf("c%d", k), ranks[k])
			if k%2 == 0 {
				sum += ranks[k]
			} else {
				sum -= ranks[k]
			}
		}
		return w20T(kp, fmt.Sprintf("A chain complex has groups of ranks %v. What is the alternating sum \\(\\sum_k (-1)^k \\operatorname{rank} C_k\\)?", ranks),
			strconv.Itoa(sum), fmt.Sprintf("Add even ranks and subtract odd ranks: the alternating sum is %d.", sum), f)
	default:
		c := w20R(5, 20)
		r1 := w20R(0, c-1)
		r2 := w20R(0, c-r1)
		w20I(f, "c", c)
		w20I(f, "r1", r1)
		w20I(f, "r2", r2)
		b := c - r1 - r2
		return w20T(kp, fmt.Sprintf("In a chain complex, \\(C_k\\) has rank %d, the map into it has rank %d and the map out of it has rank %d. What is the Betti number \\(b_k\\)?", c, r1, r2),
			strconv.Itoa(b), fmt.Sprintf("The homology is \\(\\ker/\\operatorname{im}\\), so \\(b_k=c_k-\\operatorname{rank} d_k-\\operatorname{rank} d_{k+1}=%d-%d-%d=%d\\).", c, r1, r2, b), f)
	}
}

func w20HomExactSequence(kp int) generator.Problem {
	f := map[string]string{}
	a := w20R(1, 30)
	c := w20R(1, 30)
	w20I(f, "a", a)
	w20I(f, "c", c)
	switch kp {
	case 1:
		b := a + c
		return w20T(kp, fmt.Sprintf("A short exact sequence \\(0\\to A\\to B\\to C\\to 0\\) of finite-dimensional vector spaces has \\(\\dim A=%d\\) and \\(\\dim C=%d\\). What is \\(\\dim B\\)?", a, c),
			strconv.Itoa(b), fmt.Sprintf("Exactness gives \\(\\dim B=\\dim A+\\dim C=%d+%d=%d\\).", a, c, b), f)
	case 2:
		b := a + c
		w20I(f, "b", b)
		return w20T(kp, fmt.Sprintf("A short exact sequence \\(0\\to A\\to B\\to C\\to 0\\) has \\(\\dim B=%d\\) and \\(\\dim A=%d\\). What is \\(\\dim C\\)?", b, a),
			strconv.Itoa(c), fmt.Sprintf("Rearrange \\(\\dim B=\\dim A+\\dim C\\) to get \\(\\dim C=%d-%d=%d\\).", b, a, c), f)
	default:
		b := a + c
		w20I(f, "b", b)
		return w20T(kp, fmt.Sprintf("A short exact sequence \\(0\\to A\\to B\\to C\\to 0\\) has \\(\\dim B=%d\\) and \\(\\dim C=%d\\). What is \\(\\dim A\\)?", b, c),
			strconv.Itoa(a), fmt.Sprintf("Rearrange \\(\\dim B=\\dim A+\\dim C\\) to get \\(\\dim A=%d-%d=%d\\).", b, c, a), f)
	}
}

func w20HomSnakeLemma(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		k := w20R(1, 30)
		w20I(f, "k", k)
		terms := 3 * k
		return w20T(kp, fmt.Sprintf("The long exact sequence of a pair is written out for %d consecutive degrees, three groups per degree. How many groups are displayed?", k),
			strconv.Itoa(terms), fmt.Sprintf("Each degree contributes three groups, so \\(3\\cdot %d=%d\\).", k, terms), f)
	case 2:
		n := w20R(5, 15)
		k := w20R(1, 4)
		w20I(f, "n", n)
		w20I(f, "k", k)
		degree := n - k
		return w20T(kp, fmt.Sprintf("Starting from \\(H_{%d}\\), the connecting map is applied %d times. Which degree is reached?", n, k),
			strconv.Itoa(degree), fmt.Sprintf("Each connecting map lowers the degree by one, so \\(%d-%d=%d\\).", n, k, degree), f)
	default:
		k := w20R(1, 30)
		w20I(f, "k", k)
		maps := k
		return w20T(kp, fmt.Sprintf("A long exact sequence is split into %d segments, each with one connecting map. How many connecting maps are there?", k),
			strconv.Itoa(maps), fmt.Sprintf("One connecting map per segment gives %d.", maps), f)
	}
}

func w20HomFiveLemma(kp int) generator.Problem {
	f := map[string]string{}
	n := w20R(3, 40)
	w20I(f, "n", n)
	switch kp {
	case 1:
		objects := 2 * n
		return w20T(kp, fmt.Sprintf("A commutative ladder has %d columns, each with a top and a bottom object. How many objects are in the diagram?", n),
			strconv.Itoa(objects), fmt.Sprintf("Two rows of %d objects give %d objects.", n, objects), f)
	case 2:
		squares := n - 1
		return w20T(kp, fmt.Sprintf("A commutative ladder has %d columns of objects. How many square faces does the diagram have?", n),
			strconv.Itoa(squares), fmt.Sprintf("Each adjacent pair of columns forms one square, so there are \\(%d-1=%d\\).", n, squares), f)
	default:
		arrows := 2 * (n - 1)
		return w20T(kp, fmt.Sprintf("A commutative ladder has %d columns. How many horizontal arrows does it contain?", n),
			strconv.Itoa(arrows), fmt.Sprintf("Each of the two rows has %d arrows, giving %d in total.", n-1, arrows), f)
	}
}

func w20HomKunneth(kp int) generator.Problem {
	f := map[string]string{}
	m, n := w20R(1, 10), w20R(1, 10)
	w20I(f, "m", m)
	w20I(f, "n", n)
	switch kp {
	case 1:
		b1 := m + n
		return w20T(kp, fmt.Sprintf("The Kunneth formula gives the homology of \\(T^{%d}\\times T^{%d}=T^{%d}\\). What is its first Betti number?", m, n, m+n),
			strconv.Itoa(b1), fmt.Sprintf("The first Betti number of \\(T^{m+n}\\) is \\(m+n=%d+%d=%d\\).", m, n, b1), f)
	case 2:
		total := 1 << (m + n)
		return w20T(kp, fmt.Sprintf("What is the total rank of the homology of the torus \\(T^{%d}\\times T^{%d}\\)?", m, n),
			strconv.Itoa(total), fmt.Sprintf("The product is \\(T^{%d}\\), whose Betti numbers sum to \\(2^{%d}=%d\\).", m+n, m+n, total), f)
	default:
		b2 := w20Binom(m+n, 2)
		return w20T(kp, fmt.Sprintf("What is the second Betti number of the torus \\(T^{%d}\\times T^{%d}\\)?", m, n),
			strconv.Itoa(b2), fmt.Sprintf("The second Betti number of \\(T^{%d}\\) is \\(\\binom{%d}{2}=%d\\).", m+n, m+n, b2), f)
	}
}

func w20HomUniversalCoefficient(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		r := w20R(1, 20)
		w20I(f, "r", r)
		return w20T(kp, fmt.Sprintf("A space has \\(H_n(X;\\mathbb{Z})\\cong\\mathbb{Z}^{%d}\\). What is the free rank of \\(H^n(X;\\mathbb{Z})\\)?", r),
			strconv.Itoa(r), fmt.Sprintf("The universal coefficient theorem preserves the free part, so the rank is %d.", r), f)
	case 2:
		k := w20R(2, 40)
		w20I(f, "k", k)
		return w20T(kp, fmt.Sprintf("What is the order of the group \\(\\operatorname{Ext}(\\mathbb{Z}/%d\\mathbb{Z},\\mathbb{Z})\\)?", k),
			strconv.Itoa(k), fmt.Sprintf("The universal coefficient sequence gives \\(\\operatorname{Ext}(\\mathbb{Z}/k,\\mathbb{Z})\\cong\\mathbb{Z}/k\\), of order \\(k=%d\\).", k), f)
	default:
		k := w20R(2, 40)
		w20I(f, "k", k)
		return w20T(kp, fmt.Sprintf("The lens space \\(L(%d,1)\\) has \\(H_1\\cong\\mathbb{Z}/%d\\mathbb{Z}\\). What is the order of \\(H^2(L(%d,1);\\mathbb{Z})\\)?", k, k, k),
			strconv.Itoa(k), fmt.Sprintf("The universal coefficient theorem puts the torsion of \\(H_1\\) into \\(H^2\\), so the order is %d.", k), f)
	}
}

func w20HomEulerPoincare(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		n := w20R(1, 12)
		w20I(f, "n", n)
		chi := 1
		if n%2 == 0 {
			chi = 2
		}
		return w20T(kp, fmt.Sprintf("What is the Euler characteristic \\(\\sum_k(-1)^k b_k\\) of the sphere \\(S^{%d}\\)?", n),
			strconv.Itoa(chi), fmt.Sprintf("Only \\(b_0\\) and \\(b_n\\) are nonzero, so \\(\\chi=1+(-1)^n=%d\\) for \\(n=%d\\).", chi, n), f)
	case 2:
		n := w20R(1, 30)
		w20I(f, "n", n)
		return w20T(kp, fmt.Sprintf("What is the Euler characteristic of the torus \\(T^{%d}\\)?", n),
			"0", fmt.Sprintf("The Betti numbers of \\(T^{%d}\\) are \\(\\binom{%d}{k}\\), and their alternating sum vanishes.", n, n), f)
	default:
		g := w20R(1, 30)
		w20I(f, "g", g)
		chi := 2 - 2*g
		return w20T(kp, fmt.Sprintf("What is the Euler characteristic of the closed orientable surface of genus %d?", g),
			strconv.Itoa(chi), fmt.Sprintf("Its Betti numbers are \\(1, 2g, 1\\), so \\(\\chi=2-2g=2-2\\cdot %d=%d\\).", g, chi), f)
	}
}

func w20HomSphere(kp int) generator.Problem {
	f := map[string]string{}
	n := w20R(2, 12)
	w20I(f, "n", n)
	switch kp {
	case 1:
		return w20T(kp, fmt.Sprintf("What is the total rank \\(\\sum_k \\operatorname{rank} H_k(S^{%d};\\mathbb{Z})\\)?", n),
			"2", fmt.Sprintf("Only \\(H_0\\) and \\(H_{%d}\\) are nonzero, each of rank one, so the total is 2.", n), f)
	case 2:
		return w20T(kp, fmt.Sprintf("What is the rank of \\(H_{%d}(S^{%d};\\mathbb{Z})\\)?", n, n),
			"1", fmt.Sprintf("The top homology of a connected closed orientable %d-manifold is \\(\\mathbb{Z}\\), of rank one.", n), f)
	default:
		chi := 1
		if n%2 == 0 {
			chi = 2
		}
		return w20T(kp, fmt.Sprintf("What is the Euler characteristic of \\(S^{%d}\\)?", n),
			strconv.Itoa(chi), fmt.Sprintf("\\(\\chi(S^n)=1+(-1)^n\\), which equals %d for \\(n=%d\\).", chi, n), f)
	}
}

func w20HomProjective(kp int) generator.Problem {
	f := map[string]string{}
	n := w20R(2, 12)
	w20I(f, "n", n)
	switch kp {
	case 1:
		summands := 1
		if n%2 == 1 {
			summands = 2
		}
		return w20T(kp, fmt.Sprintf("How many infinite cyclic summands does \\(H_*(\\mathbb{RP}^{%d};\\mathbb{Z})\\) contain?", n),
			strconv.Itoa(summands), fmt.Sprintf("The group is \\(\\mathbb{Z}\\) in degree 0, and also \\(\\mathbb{Z}\\) in degree %d when that degree is odd; here the count is %d.", n, summands), f)
	case 2:
		chi := 0
		if n%2 == 0 {
			chi = 1
		}
		return w20T(kp, fmt.Sprintf("What is the Euler characteristic of \\(\\mathbb{RP}^{%d}\\)?", n),
			strconv.Itoa(chi), fmt.Sprintf("The rational homology vanishes in positive degrees when \\(n\\) is odd, so \\(\\chi\\) is %d for \\(n=%d\\).", chi, n), f)
	default:
		torsion := n / 2
		return w20T(kp, fmt.Sprintf("How many \\(\\mathbb{Z}/2\\mathbb{Z}\\) summands does \\(H_*(\\mathbb{RP}^{%d};\\mathbb{Z})\\) contain?", n),
			strconv.Itoa(torsion), fmt.Sprintf("The torsion appears in the odd degrees below %d, giving \\(\\lfloor %d/2\\rfloor=%d\\) summands.", n, n, torsion), f)
	}
}

func w20HomTorus(kp int) generator.Problem {
	f := map[string]string{}
	n := w20R(1, 12)
	w20I(f, "n", n)
	switch kp {
	case 1:
		return w20T(kp, fmt.Sprintf("What is the first Betti number of the torus \\(T^{%d}\\)?", n),
			strconv.Itoa(n), fmt.Sprintf("The torus has \\(n\\) independent circle directions, so \\(b_1=n=%d\\).", n), f)
	case 2:
		total := 1 << n
		return w20T(kp, fmt.Sprintf("What is the total rank of \\(H_*(T^{%d};\\mathbb{Z})\\)?", n),
			strconv.Itoa(total), fmt.Sprintf("The Betti numbers are \\(\\binom{n}{k}\\), whose sum is \\(2^n=2^{%d}=%d\\).", n, total), f)
	default:
		b2 := w20Binom(n, 2)
		return w20T(kp, fmt.Sprintf("What is the second Betti number of \\(T^{%d}\\)?", n),
			strconv.Itoa(b2), fmt.Sprintf("Choose two of the %d circle directions: \\(\\binom{%d}{2}=%d\\).", n, n, b2), f)
	}
}

// ----- knots -----

func w20KnotCrossing(kp int) generator.Problem {
	f := map[string]string{}
	p, q := w20KnotPair()
	w20I(f, "p", p)
	w20I(f, "q", q)
	switch kp {
	case 1:
		crossings := p*q - w20Max(p, q)
		return w20T(kp, fmt.Sprintf("What is the crossing number of the torus knot \\(T(%d,%d)\\)?", p, q),
			strconv.Itoa(crossings), fmt.Sprintf("For \\(T(p,q)\\) the crossing number is \\(pq-\\max(p,q)=%d\\cdot %d-%d=%d\\).", p, q, w20Max(p, q), crossings), f)
	case 2:
		return w20T(kp, fmt.Sprintf("What is the crossing number of the torus knot \\(T(2,%d)\\)?", q),
			strconv.Itoa(q), fmt.Sprintf("\\(T(2,q)\\) is the \\((2,q)\\) torus knot with crossing number \\(q=%d\\).", q), f)
	default:
		return w20T(kp, fmt.Sprintf("What is the crossing number of the torus knot \\(T(%d,2)\\)?", p),
			strconv.Itoa(p), fmt.Sprintf("\\(T(p,2)\\) has crossing number \\(p=%d\\).", p), f)
	}
}

func w20KnotUnknotting(kp int) generator.Problem {
	f := map[string]string{}
	p, q := w20KnotPair()
	w20I(f, "p", p)
	w20I(f, "q", q)
	switch kp {
	case 1:
		u := (p - 1) * (q - 1) / 2
		return w20T(kp, fmt.Sprintf("What is the unknotting number of the torus knot \\(T(%d,%d)\\)?", p, q),
			strconv.Itoa(u), fmt.Sprintf("For \\(T(p,q)\\) the unknotting number is \\((p-1)(q-1)/2=%d\\).", u), f)
	case 2:
		u := (q - 1) / 2
		return w20T(kp, fmt.Sprintf("What is the unknotting number of \\(T(2,%d)\\)?", q),
			strconv.Itoa(u), fmt.Sprintf("Substitute \\(p=2\\): \\((2-1)(%d-1)/2=%d\\).", q, u), f)
	default:
		u := q - 1
		return w20T(kp, fmt.Sprintf("What is the unknotting number of \\(T(3,%d)\\)?", q),
			strconv.Itoa(u), fmt.Sprintf("Substitute \\(p=3\\): \\((3-1)(%d-1)/2=%d\\).", q, u), f)
	}
}

func w20KnotSeifertGenus(kp int) generator.Problem {
	f := map[string]string{}
	p, q := w20KnotPair()
	w20I(f, "p", p)
	w20I(f, "q", q)
	switch kp {
	case 1:
		g := (p - 1) * (q - 1) / 2
		return w20T(kp, fmt.Sprintf("What is the Seifert genus of the torus knot \\(T(%d,%d)\\)?", p, q),
			strconv.Itoa(g), fmt.Sprintf("The genus is \\((p-1)(q-1)/2=%d\\).", g), f)
	case 2:
		g := (q - 1) / 2
		return w20T(kp, fmt.Sprintf("What is the Seifert genus of \\(T(2,%d)\\)?", q),
			strconv.Itoa(g), fmt.Sprintf("With \\(p=2\\), the genus is \\((2-1)(%d-1)/2=%d\\).", q, g), f)
	default:
		twice := (p - 1) * (q - 1)
		return w20T(kp, fmt.Sprintf("For \\(T(%d,%d)\\), what is twice the Seifert genus?", p, q),
			strconv.Itoa(twice), fmt.Sprintf("Twice the genus is \\((p-1)(q-1)=%d\\).", twice), f)
	}
}

func w20KnotAlexander(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		p, q := w20KnotPair()
		w20I(f, "p", p)
		w20I(f, "q", q)
		degree := (p - 1) * (q - 1)
		return w20T(kp, fmt.Sprintf("What is the degree of the Alexander polynomial of \\(T(%d,%d)\\)?", p, q),
			strconv.Itoa(degree), fmt.Sprintf("The Alexander polynomial has degree \\((p-1)(q-1)=%d\\).", degree), f)
	case 2:
		q := w20R(2, 13)
		w20I(f, "q", q)
		degree := q - 1
		return w20T(kp, fmt.Sprintf("What is the degree of the Alexander polynomial of \\(T(2,%d)\\)?", q),
			strconv.Itoa(degree), fmt.Sprintf("With \\(p=2\\), the degree is \\((2-1)(%d-1)=%d\\).", q, degree), f)
	default:
		d1, d2 := w20R(1, 12), w20R(1, 12)
		w20I(f, "d1", d1)
		w20I(f, "d2", d2)
		sum := d1 + d2
		return w20T(kp, fmt.Sprintf("Two knots have Alexander polynomials of degrees %d and %d. What is the degree of the Alexander polynomial of their connected sum?", d1, d2),
			strconv.Itoa(sum), fmt.Sprintf("Alexander polynomials multiply, so their degrees add: \\(%d+%d=%d\\).", d1, d2, sum), f)
	}
}

func w20KnotSkein(kp int) generator.Problem {
	f := map[string]string{}
	c := w20R(3, 20)
	w20I(f, "c", c)
	switch kp {
	case 1:
		states := 1 << c
		return w20T(kp, fmt.Sprintf("A skein tree smooths each crossing in two ways. A diagram has %d crossings. How many leaves does the tree have?", c),
			strconv.Itoa(states), fmt.Sprintf("Each crossing doubles the number of leaves: \\(2^{%d}=%d\\).", c, states), f)
	case 2:
		k := w20R(0, c)
		w20I(f, "k", k)
		terms := w20Binom(c, k)
		return w20T(kp, fmt.Sprintf("In a skein tree on %d crossings, how many leaves have exactly %d crossings smoothed one way?", c, k),
			strconv.Itoa(terms), fmt.Sprintf("Choose which %d of the %d crossings are smoothed that way: \\(\\binom{%d}{%d}=%d\\).", k, c, c, k, terms), f)
	default:
		nodes := (1 << c) - 1
		return w20T(kp, fmt.Sprintf("A full binary skein tree has %d levels of crossings. How many internal nodes does it contain?", c),
			strconv.Itoa(nodes), fmt.Sprintf("A full binary tree with \\(2^{%d}\\) leaves has \\(2^{%d}-1=%d\\) internal nodes.", c, c, nodes), f)
	}
}

func w20KnotTorusKnot(kp int) generator.Problem {
	f := map[string]string{}
	p, q := w20KnotPair()
	w20I(f, "p", p)
	w20I(f, "q", q)
	switch kp {
	case 1:
		braid := w20Min(p, q)
		return w20T(kp, fmt.Sprintf("What is the braid index of the torus knot \\(T(%d,%d)\\)?", p, q),
			strconv.Itoa(braid), fmt.Sprintf("The braid index of \\(T(p,q)\\) is \\(\\min(p,q)=%d\\).", braid), f)
	case 2:
		crossings := p*q - w20Max(p, q)
		return w20T(kp, fmt.Sprintf("How many crossings does the standard diagram of \\(T(%d,%d)\\) have?", p, q),
			strconv.Itoa(crossings), fmt.Sprintf("Crossings are \\(pq-\\max(p,q)=%d\\).", crossings), f)
	default:
		g := (p - 1) * (q - 1) / 2
		return w20T(kp, fmt.Sprintf("What is the genus of \\(T(%d,%d)\\)?", p, q),
			strconv.Itoa(g), fmt.Sprintf("The genus is \\((p-1)(q-1)/2=%d\\).", g), f)
	}
}

func w20KnotConnectedSum(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		g1, g2 := w20R(1, 8), w20R(1, 8)
		w20I(f, "g1", g1)
		w20I(f, "g2", g2)
		g := g1 + g2
		return w20T(kp, fmt.Sprintf("Knots of genera %d and %d are connected-summed. What is the genus of the result?", g1, g2),
			strconv.Itoa(g), fmt.Sprintf("Genus is additive under connected sum: \\(%d+%d=%d\\).", g1, g2, g), f)
	case 2:
		c1, c2 := w20R(3, 20), w20R(3, 20)
		w20I(f, "c1", c1)
		w20I(f, "c2", c2)
		c := c1 + c2
		return w20T(kp, fmt.Sprintf("Knots with crossing numbers %d and %d are connected-summed. What is the crossing number of the result?", c1, c2),
			strconv.Itoa(c), fmt.Sprintf("Crossing number is additive for connected sums of knots: \\(%d+%d=%d\\).", c1, c2, c), f)
	default:
		d1, d2 := w20R(1, 12), w20R(1, 12)
		w20I(f, "d1", d1)
		w20I(f, "d2", d2)
		d := d1 + d2
		return w20T(kp, fmt.Sprintf("Alexander polynomials of degrees %d and %d are multiplied under connected sum. What is the degree of the product?", d1, d2),
			strconv.Itoa(d), fmt.Sprintf("Degrees add when polynomials multiply: \\(%d+%d=%d\\).", d1, d2, d), f)
	}
}

func w20KnotBracket(kp int) generator.Problem {
	f := map[string]string{}
	c := w20R(3, 20)
	w20I(f, "c", c)
	switch kp {
	case 1:
		states := 1 << c
		return w20T(kp, fmt.Sprintf("The Kauffman bracket is a state sum over smoothings. A diagram with %d crossings has how many states?", c),
			strconv.Itoa(states), fmt.Sprintf("Two choices per crossing give \\(2^{%d}=%d\\) states.", c, states), f)
	case 2:
		p, n := w20R(0, 12), w20R(0, 12)
		w20I(f, "P", p)
		w20I(f, "N", n)
		writhe := p - n
		return w20T(kp, fmt.Sprintf("A diagram has %d positive and %d negative crossings. What is its writhe?", p, n),
			strconv.Itoa(writhe), fmt.Sprintf("Writhe is positive minus negative: \\(%d-%d=%d\\).", p, n, writhe), f)
	default:
		p, n := w20R(0, 12), w20R(0, 12)
		w20I(f, "P", p)
		w20I(f, "N", n)
		return w20T(kp, fmt.Sprintf("A diagram has %d positive and %d negative crossings. How many crossings does it have in total?", p, n),
			strconv.Itoa(p+n), fmt.Sprintf("The total is the sum of the two counts: \\(%d+%d=%d\\).", p, n, p+n), f)
	}
}

func w20KnotWrithe(kp int) generator.Problem {
	f := map[string]string{}
	p, n := w20R(0, 15), w20R(0, 15)
	w20I(f, "P", p)
	w20I(f, "N", n)
	switch kp {
	case 1:
		writhe := p - n
		return w20T(kp, fmt.Sprintf("A knot diagram has %d positive and %d negative crossings. What is its writhe?", p, n),
			strconv.Itoa(writhe), fmt.Sprintf("Writhe is \\(P-N=%d-%d=%d\\).", p, n, writhe), f)
	case 2:
		total := p + n
		return w20T(kp, fmt.Sprintf("A knot diagram has %d positive and %d negative crossings. How many crossings are there in total?", p, n),
			strconv.Itoa(total), fmt.Sprintf("Add the two counts: \\(%d+%d=%d\\).", p, n, total), f)
	default:
		k := w20R(1, 10)
		w20I(f, "k", k)
		writhe := p + k - n
		return w20T(kp, fmt.Sprintf("A diagram has %d positive and %d negative crossings. Adding %d positive Reidemeister I kinks gives what writhe?", p, n, k),
			strconv.Itoa(writhe), fmt.Sprintf("Each positive kink adds one: \\(%d+%d-%d=%d\\).", p, k, n, writhe), f)
	}
}

func w20KnotLinking(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		n := w20R(1, 12)
		w20I(f, "n", n)
		return w20T(kp, fmt.Sprintf("What is the linking number of the torus link \\(T(2,2\\cdot %d)\\)?", n),
			strconv.Itoa(n), fmt.Sprintf("The \\((2,2n)\\) torus link has linking number \\(n=%d\\).", n), f)
	case 2:
		l := w20R(-6, 6)
		p := w20R(0, 12)
		n := p - 2*l
		if n < 0 {
			n = 0
			p = 2 * l
			if p < 0 {
				p = 0
				l = 0
			}
		}
		w20I(f, "P", p)
		w20I(f, "N", n)
		w20I(f, "L", l)
		return w20T(kp, fmt.Sprintf("A two-component link diagram has %d positive and %d negative crossings between the components. What is its linking number?", p, n),
			strconv.Itoa(l), fmt.Sprintf("Linking number is half the signed sum: \\((P-N)/2=(%d-%d)/2=%d\\).", p, n, l), f)
	default:
		n := w20R(0, 12)
		l := w20R(-5, 5)
		w20I(f, "N", n)
		w20I(f, "L", l)
		p := n + 2*l
		return w20T(kp, fmt.Sprintf("A link diagram has %d negative and \\(P\\) positive inter-component crossings, with linking number %d. What is \\(P\\)?", n, l),
			strconv.Itoa(p), fmt.Sprintf("Solve \\((P-N)/2=L\\): \\(P=N+2L=%d+2\\cdot %d=%d\\).", n, l, p), f)
	}
}

// ----- geometric group theory -----

func w20GeoCayley(kp int) generator.Problem {
	f := map[string]string{}
	n, k := w20R(3, 30), w20R(1, 6)
	w20I(f, "n", n)
	w20I(f, "k", k)
	switch kp {
	case 1:
		degree := 2 * k
		return w20T(kp, fmt.Sprintf("The Cayley graph of \\(\\mathbb{Z}/%d\\mathbb{Z}\\) uses %d generators and their inverses, with no generator of order two. What is its degree?", n, k),
			strconv.Itoa(degree), fmt.Sprintf("Each generator contributes two neighbours, so the degree is \\(2k=2\\cdot %d=%d\\).", k, degree), f)
	case 2:
		return w20T(kp, fmt.Sprintf("How many vertices does the Cayley graph of \\(\\mathbb{Z}/%d\\mathbb{Z}\\) with %d generators have?", n, k),
			strconv.Itoa(n), fmt.Sprintf("The Cayley graph has one vertex per group element, so it has %d vertices.", n), f)
	default:
		edges := n * k
		return w20T(kp, fmt.Sprintf("The Cayley graph of \\(\\mathbb{Z}/%d\\mathbb{Z}\\) with %d inverse-closed generator pairs has how many edges?", n, k),
			strconv.Itoa(edges), fmt.Sprintf("The degree is \\(2k\\) on %d vertices, so the edge count is \\(n\\cdot 2k/2=nk=%d\\cdot %d=%d\\).", n, n, k, edges), f)
	}
}

func w20GeoGrowth(kp int) generator.Problem {
	f := map[string]string{}
	k := w20R(2, 30)
	w20I(f, "k", k)
	switch kp {
	case 1:
		words := 2 * k
		return w20T(kp, fmt.Sprintf("The free group \\(F_{%d}\\) has %d generators. How many reduced words of length 1 are there?", k, k),
			strconv.Itoa(words), fmt.Sprintf("Each of the %d generators and its inverse gives one word, so \\(2\\cdot %d=%d\\).", k, k, words), f)
	case 2:
		words := 2 * k * (2*k - 1)
		return w20T(kp, fmt.Sprintf("How many reduced words of length 2 does the free group \\(F_{%d}\\) have?", k),
			strconv.Itoa(words), fmt.Sprintf("After the first letter, the inverse of the previous letter is forbidden, leaving \\(2k-1\\) choices: \\(2\\cdot %d\\cdot %d=%d\\).", k, 2*k-1, words), f)
	default:
		rate := 2*k - 1
		return w20T(kp, fmt.Sprintf("The number of reduced words of length \\(n\\) in \\(F_{%d}\\) grows like \\(C\\lambda^n\\). What is \\(\\lambda\\)?", k),
			strconv.Itoa(rate), fmt.Sprintf("Each step has \\(2k-1\\) choices, so the growth rate is \\(2\\cdot %d-1=%d\\).", k, rate), f)
	}
}

func w20GeoQuasiIsometry(kp int) generator.Problem {
	f := map[string]string{}
	n, k := w20R(1, 12), w20R(2, 30)
	w20I(f, "n", n)
	w20I(f, "k", k)
	switch kp {
	case 1:
		return w20T(kp, fmt.Sprintf("The subgroup \\(%d\\mathbb{Z}\\) has finite index in \\(\\mathbb{Z}\\). What is that index?", n),
			strconv.Itoa(n), fmt.Sprintf("The quotient \\(\\mathbb{Z}/%d\\mathbb{Z}\\) has %d elements, so the index is %d.", n, n, n), f)
	case 2:
		return w20T(kp, fmt.Sprintf("The group \\(\\mathbb{Z}^{%d}\\) is quasi-isometric to \\(\\mathbb{R}^{%d}\\). What is the quasi-isometry dimension?", n, n),
			strconv.Itoa(n), fmt.Sprintf("The dimension is preserved by quasi-isometry, so it is %d.", n), f)
	default:
		return w20T(kp, fmt.Sprintf("The regular tree \\(T_{%d}\\) has %d branches at every vertex. What is its branching factor?", k, k),
			strconv.Itoa(k), fmt.Sprintf("Each vertex has %d neighbours beyond the one it came from, so the branching factor is %d.", k, k), f)
	}
}

func w20GeoEnds(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		n := w20R(1, 30)
		w20I(f, "n", n)
		ends := 1
		if n == 1 {
			ends = 2
		}
		return w20T(kp, fmt.Sprintf("How many ends does \\(\\mathbb{Z}^{%d}\\) have?", n),
			strconv.Itoa(ends), fmt.Sprintf("A line has two ends, while \\(\\mathbb{Z}^n\\) for \\(n\\ge 2\\) has one; here \\(n=%d\\), so the answer is %d.", n, ends), f)
	case 2:
		k := w20R(2, 30)
		w20I(f, "k", k)
		ends := 2 * k
		return w20T(kp, fmt.Sprintf("How many ends does the free group \\(F_{%d}\\) have?", k),
			strconv.Itoa(ends), fmt.Sprintf("A free group of rank \\(k\\ge 2\\) has \\(2k\\) ends, so \\(2\\cdot %d=%d\\).", k, ends), f)
	default:
		m, n := w20R(1, 8), w20R(1, 8)
		w20I(f, "m", m)
		w20I(f, "n", n)
		ends := 1
		if m+n == 1 {
			ends = 2
		}
		return w20T(kp, fmt.Sprintf("How many ends does \\(\\mathbb{Z}^{%d}\\times\\mathbb{Z}^{%d}\\) have?", m, n),
			strconv.Itoa(ends), fmt.Sprintf("The product is \\(\\mathbb{Z}^{%d}\\), which has %d end(s).", m+n, ends), f)
	}
}

func w20GeoSvarcMilnor(kp int) generator.Problem {
	f := map[string]string{}
	m, n := w20R(1, 20), w20R(1, 20)
	w20I(f, "m", m)
	w20I(f, "n", n)
	switch kp {
	case 1:
		return w20T(kp, fmt.Sprintf("What is the growth degree of \\(\\mathbb{Z}^{%d}\\)?", m),
			strconv.Itoa(m), fmt.Sprintf("A ball of radius \\(r\\) in \\(\\mathbb{Z}^{%d}\\) has about \\(r^{%d}\\) points, so the growth degree is %d.", m, m, m), f)
	case 2:
		degree := m + n
		return w20T(kp, fmt.Sprintf("What is the growth degree of \\(\\mathbb{Z}^{%d}\\times\\mathbb{Z}^{%d}\\)?", m, n),
			strconv.Itoa(degree), fmt.Sprintf("The product is \\(\\mathbb{Z}^{%d}\\), of growth degree \\(%d+%d=%d\\).", m+n, m, n, degree), f)
	default:
		return w20T(kp, fmt.Sprintf("\\(\\mathbb{Z}^{%d}\\) is quasi-isometric to \\(\\mathbb{R}^{n}\\) only for one value of \\(n\\). What is that value?", m),
			strconv.Itoa(m), fmt.Sprintf("A quasi-isometry preserves the dimension, so \\(\\mathbb{Z}^{%d}\\) matches \\(\\mathbb{R}^{%d}\\), and the value is %d.", m, m, m), f)
	}
}

func w20GeoHyperbolic(kp int) generator.Problem {
	f := map[string]string{}
	n := w20R(2, 12)
	w20I(f, "n", n)
	switch kp {
	case 1:
		dim := n - 1
		return w20T(kp, fmt.Sprintf("The visual boundary of hyperbolic space \\(\\mathbb{H}^{%d}\\) is a sphere. What is its dimension?", n),
			strconv.Itoa(dim), fmt.Sprintf("The boundary of \\(\\mathbb{H}^n\\) is \\(S^{n-1}\\), of dimension \\(n-1=%d\\).", dim), f)
	case 2:
		rate := n - 1
		return w20T(kp, fmt.Sprintf("Volumes of balls in \\(\\mathbb{H}^{%d}\\) grow exponentially. What is the exponent \\(\\lambda\\) in \\(\\operatorname{Vol}(B_r)\\sim e^{\\lambda r}\\)?", n),
			strconv.Itoa(rate), fmt.Sprintf("The volume growth exponent of \\(\\mathbb{H}^n\\) is \\(n-1=%d\\).", rate), f)
	default:
		dim := n - 1
		return w20T(kp, fmt.Sprintf("What is the covering dimension of the boundary sphere of \\(\\mathbb{H}^{%d}\\)?", n),
			strconv.Itoa(dim), fmt.Sprintf("The boundary is \\(S^{%d}\\), whose covering dimension is %d.", n-1, dim), f)
	}
}

func w20GeoBoundary(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		n := w20R(1, 12)
		w20I(f, "n", n)
		dim := n - 1
		return w20T(kp, fmt.Sprintf("The visual boundary of \\(\\mathbb{H}^{%d}\\) is a sphere. What is its dimension?", n),
			strconv.Itoa(dim), fmt.Sprintf("The boundary is \\(S^{%d}\\), so the dimension is %d.", n-1, dim), f)
	case 2:
		n := w20R(1, 30)
		w20I(f, "n", n)
		ends := 1
		if n == 1 {
			ends = 2
		}
		return w20T(kp, fmt.Sprintf("How many ends does the boundary of \\(\\mathbb{Z}^{%d}\\) detect?", n),
			strconv.Itoa(ends), fmt.Sprintf("The number of ends of \\(\\mathbb{Z}^{%d}\\) is %d.", n, ends), f)
	default:
		m, n := w20R(1, 8), w20R(1, 8)
		w20I(f, "m", m)
		w20I(f, "n", n)
		ends := 1
		if m+n == 1 {
			ends = 2
		}
		return w20T(kp, fmt.Sprintf("How many ends does \\(\\mathbb{Z}^{%d}\\times\\mathbb{Z}^{%d}\\) have?", m, n),
			strconv.Itoa(ends), fmt.Sprintf("The product is \\(\\mathbb{Z}^{%d}\\), with %d end(s).", m+n, ends), f)
	}
}

func w20GeoDehn(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		n := w20R(1, 30)
		w20I(f, "n", n)
		degree := 1
		return w20T(kp, fmt.Sprintf("What is the polynomial degree of the Dehn function of \\(\\mathbb{Z}^{%d}\\)?", n),
			strconv.Itoa(degree), "The Dehn function of a hyperbolic group, and of \\(\\mathbb{Z}\\), is linear, of degree 1.", f)
	case 2:
		n := w20R(2, 30)
		w20I(f, "n", n)
		return w20T(kp, fmt.Sprintf("What is the polynomial degree of the Dehn function of \\(\\mathbb{Z}^{%d}\\)?", n),
			"2", fmt.Sprintf("For \\(n\\ge 2\\), \\(\\mathbb{Z}^n\\) has a quadratic Dehn function, of degree 2; here \\(n=%d\\).", n), f)
	default:
		n := w20R(3, 30)
		w20I(f, "n", n)
		return w20T(kp, fmt.Sprintf("What is the polynomial degree of the Dehn function of \\(\\mathbb{Z}^{%d}\\)?", n),
			"2", fmt.Sprintf("For \\(n\\ge 2\\), \\(\\mathbb{Z}^n\\) has a quadratic Dehn function, of degree 2; here \\(n=%d\\).", n), f)
	}
}

func w20GeoAmenable(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		n := w20R(1, 30)
		w20I(f, "n", n)
		return w20T(kp, fmt.Sprintf("Is \\(\\mathbb{Z}^{%d}\\) amenable? Enter 1 for yes and 0 for no.", n),
			"1", "Every abelian group is amenable, and \\(\\mathbb{Z}^n\\) is abelian.", f)
	case 2:
		k := w20R(2, 30)
		w20I(f, "k", k)
		return w20T(kp, fmt.Sprintf("Is the free group \\(F_{%d}\\) amenable? Enter 1 for yes and 0 for no.", k),
			"0", "A free group of rank at least two contains a free subgroup on two generators, and \\(F_2\\) is not amenable.", f)
	default:
		m, k := w20R(1, 12), w20R(2, 12)
		w20I(f, "m", m)
		w20I(f, "k", k)
		return w20T(kp, fmt.Sprintf("Is \\(\\mathbb{Z}^{%d}\\times F_{%d}\\) amenable? Enter 1 for yes and 0 for no.", m, k),
			"0", "A group containing \\(F_2\\) is not amenable, and \\(F_k\\) for \\(k\\ge 2\\) contains a copy of \\(F_2\\).", f)
	}
}

func w20GeoPropertyT(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		n := w20R(1, 30)
		w20I(f, "n", n)
		return w20T(kp, fmt.Sprintf("Does \\(\\mathbb{Z}^{%d}\\) have property (T)? Enter 1 for yes and 0 for no.", n),
			"0", "Abelian groups do not have property (T); they have almost invariant vectors for their regular representations.", f)
	case 2:
		n := w20R(3, 30)
		w20I(f, "n", n)
		return w20T(kp, fmt.Sprintf("Does \\(\\mathrm{SL}(%d,\\mathbb{Z})\\) have property (T)? Enter 1 for yes and 0 for no.", n),
			"1", fmt.Sprintf("For \\(n\\ge 3\\), \\(\\mathrm{SL}(n,\\mathbb{Z})\\) has property (T); here \\(n=%d\\).", n), f)
	default:
		n := w20R(2, 24)
		w20I(f, "n", n)
		return w20T(kp, fmt.Sprintf("Does a finite group of order %d have property (T)? Enter 1 for yes and 0 for no.", n),
			"1", "Every finite group has property (T), because its regular representation has no almost invariant vectors.", f)
	}
}
