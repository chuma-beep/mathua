package topology

// Wave 20 topology, first batch: surfaces and smooth manifolds. Every answer is
// an integer computed from the parameters published in Facts, so
// wave20_test.go can rebuild the formula without trusting the generator.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave20(reg *generator.Registry) {
	reg.Register("topo.surface.connected_sum", &w20Family{gen: w20SurfConnectedSum})
	reg.Register("topo.surface.crosscap", &w20Family{gen: w20SurfCrosscap})
	reg.Register("topo.surface.handle", &w20Family{gen: w20SurfHandle})
	reg.Register("topo.surface.euler", &w20Family{gen: w20SurfEuler})
	reg.Register("topo.surface.punctured", &w20Family{gen: w20SurfPunctured})
	reg.Register("topo.surface.boundary", &w20Family{gen: w20SurfBoundary})
	reg.Register("topo.surface.double_cover", &w20Family{gen: w20SurfDoubleCover})
	reg.Register("topo.surface.genus_formula", &w20Family{gen: w20SurfGenusFormula})
	reg.Register("topo.surface.boundary_connect", &w20Family{gen: w20SurfBoundaryConnect})
	reg.Register("topo.surface.triangulation", &w20Family{gen: w20SurfTriangulation})
	reg.Register("topo.diff.tangent_bundle", &w20Family{gen: w20DiffTangentBundle})
	reg.Register("topo.diff.vector_field", &w20Family{gen: w20DiffVectorField})
	reg.Register("topo.diff.critical_point", &w20Family{gen: w20DiffCriticalPoint})
	reg.Register("topo.diff.morse", &w20Family{gen: w20DiffMorse})
	reg.Register("topo.diff.sard", &w20Family{gen: w20DiffSard})
	reg.Register("topo.diff.transversality", &w20Family{gen: w20DiffTransversality})
	reg.Register("topo.diff.whitney", &w20Family{gen: w20DiffWhitney})
	reg.Register("topo.diff.tubular", &w20Family{gen: w20DiffTubular})
	reg.Register("topo.diff.degree", &w20Family{gen: w20DiffDegree})
	reg.Register("topo.diff.flow", &w20Family{gen: w20DiffFlow})
	registerWave20Rest(reg)
}

type w20Family struct {
	gen func(kp int) generator.Problem
}

func (g *w20Family) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	return g.gen(kp)
}

func w20T(kp int, q, a, ex string, f map[string]string) generator.Problem {
	if f == nil {
		f = map[string]string{}
	}
	f["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: f, KPOrdinal: kp}
}

func w20I(f map[string]string, k string, v int) { f[k] = strconv.Itoa(v) }

func w20R(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func w20ChiOrientable(g int) int { return 2 - 2*g }

func w20ChiNonorientable(k int) int { return 2 - k }

func w20Binom(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	if k > n-k {
		k = n - k
	}
	out := 1
	for i := 1; i <= k; i++ {
		out = out * (n - k + i) / i
	}
	return out
}

// ----- surfaces -----

func w20SurfConnectedSum(kp int) generator.Problem {
	f := map[string]string{}
	g1, g2 := w20R(1, 9), w20R(1, 9)
	w20I(f, "g1", g1)
	w20I(f, "g2", g2)
	switch kp {
	case 1:
		chi := 2 - 2*(g1+g2)
		return w20T(kp, fmt.Sprintf("Two closed orientable surfaces have genera %d and %d. What is the Euler characteristic of their connected sum?", g1, g2),
			strconv.Itoa(chi), fmt.Sprintf("The connected sum satisfies \\(\\chi(S_1 \\# S_2)=\\chi(S_1)+\\chi(S_2)-2\\). With \\(\\chi=2-2g\\) for each surface, the genus adds: \\(\\chi=2-2(%d+%d)=%d\\).", g1, g2, chi), f)
	case 2:
		chi1, chi2 := 2-2*g1, 2-2*g2
		w20I(f, "chi1", chi1)
		w20I(f, "chi2", chi2)
		chi := chi1 + chi2 - 2
		return w20T(kp, fmt.Sprintf("A surface with \\(\\chi=%d\\) is glued to one with \\(\\chi=%d\\) by a connected sum. What is the Euler characteristic of the result?", chi1, chi2),
			strconv.Itoa(chi), fmt.Sprintf("The connected sum subtracts one disk from each surface before gluing, so \\(\\chi=\\chi_1+\\chi_2-2=%d+%d-2=%d\\).", chi1, chi2, chi), f)
	default:
		genus := g1 + g2
		return w20T(kp, fmt.Sprintf("A surface of genus %d is connected-summed with one of genus %d. What is the genus of the result?", g1, g2),
			strconv.Itoa(genus), fmt.Sprintf("Connected sum adds orientable genera: \\(%d+%d=%d\\).", g1, g2, genus), f)
	}
}

func w20SurfCrosscap(kp int) generator.Problem {
	f := map[string]string{}
	k := w20R(1, 20)
	w20I(f, "k", k)
	switch kp {
	case 1:
		chi := 2 - k
		return w20T(kp, fmt.Sprintf("The nonorientable surface \\(N_{%d}\\) is the connected sum of %d projective planes. What is its Euler characteristic?", k, k),
			strconv.Itoa(chi), fmt.Sprintf("A connected sum of \\(k\\) projective planes has \\(\\chi=2-k=2-%d=%d\\).", k, chi), f)
	case 2:
		chi := 2 - k
		w20I(f, "chi", chi)
		return w20T(kp, fmt.Sprintf("A connected sum of projective planes has Euler characteristic %d. How many projective planes does it contain?", chi),
			strconv.Itoa(k), fmt.Sprintf("Solve \\(2-k=%d\\), so \\(k=2-%d=%d\\).", chi, chi, k), f)
	default:
		genus := k - 1
		return w20T(kp, fmt.Sprintf("What is the genus of the orientable double cover of \\(N_{%d}\\), the connected sum of %d projective planes?", k, k),
			strconv.Itoa(genus), fmt.Sprintf("The orientable double cover of \\(N_k\\) is the orientable surface of genus \\(k-1=%d-1=%d\\).", k, genus), f)
	}
}

func w20SurfHandle(kp int) generator.Problem {
	f := map[string]string{}
	g := w20R(1, 15)
	w20I(f, "g", g)
	switch kp {
	case 1:
		oneHandles := 2 * g
		return w20T(kp, fmt.Sprintf("A closed orientable surface of genus %d is built from one 0-handle, some 1-handles and one 2-handle. How many 1-handles are needed?", g),
			strconv.Itoa(oneHandles), fmt.Sprintf("Each handle adds one to the genus and contributes two 1-handles, so \\(2g=2\\cdot %d=%d\\).", g, oneHandles), f)
	case 2:
		total := 2*g + 2
		return w20T(kp, fmt.Sprintf("A genus-%d handle decomposition uses one 0-handle, \\(2g\\) 1-handles and one 2-handle. How many handles are there in total?", g),
			strconv.Itoa(total), fmt.Sprintf("The total is \\(1+2g+1=2g+2=2\\cdot %d+2=%d\\).", g, total), f)
	default:
		chi := 1 - 2*g
		return w20T(kp, fmt.Sprintf("A surface is built from one 0-handle and %d 1-handles, with no 2-handle. What is its Euler characteristic?", 2*g),
			strconv.Itoa(chi), fmt.Sprintf("Handles contribute with sign: \\(\\chi=1-(2g)=1-%d=%d\\).", 2*g, chi), f)
	}
}

func w20SurfEuler(kp int) generator.Problem {
	f := map[string]string{}
	g := w20R(1, 20)
	w20I(f, "g", g)
	switch kp {
	case 1:
		chi := 2 - 2*g
		return w20T(kp, fmt.Sprintf("What is the Euler characteristic of the closed orientable surface of genus %d?", g),
			strconv.Itoa(chi), fmt.Sprintf("For a closed orientable genus-\\(g\\) surface, \\(\\chi=2-2g=2-2\\cdot %d=%d\\).", g, chi), f)
	case 2:
		chi := 2 - 2*g
		w20I(f, "chi", chi)
		return w20T(kp, fmt.Sprintf("A closed orientable surface has Euler characteristic %d. What is its genus?", chi),
			strconv.Itoa(g), fmt.Sprintf("Solve \\(2-2g=%d\\), so \\(g=(2-%d)/2=%d\\).", chi, chi, g), f)
	default:
		b1 := 2 * g
		return w20T(kp, fmt.Sprintf("What is the first Betti number \\(b_1\\) of the closed orientable surface of genus %d?", g),
			strconv.Itoa(b1), fmt.Sprintf("The first Betti number of a genus-\\(g\\) surface is \\(2g=2\\cdot %d=%d\\).", g, b1), f)
	}
}

func w20SurfPunctured(kp int) generator.Problem {
	f := map[string]string{}
	g, p := w20R(0, 8), w20R(1, 30)
	w20I(f, "g", g)
	w20I(f, "p", p)
	switch kp {
	case 1:
		chi := 2 - 2*g - p
		return w20T(kp, fmt.Sprintf("An orientable surface of genus %d has %d punctures. What is its Euler characteristic?", g, p),
			strconv.Itoa(chi), fmt.Sprintf("Each puncture lowers \\(\\chi\\) by one: \\(\\chi=2-2g-p=2-2\\cdot %d-%d=%d\\).", g, p, chi), f)
	case 2:
		b1 := 2*g + p - 1
		return w20T(kp, fmt.Sprintf("An orientable surface of genus %d with %d punctures retracts onto a wedge of circles. What is the first Betti number?", g, p),
			strconv.Itoa(b1), fmt.Sprintf("With \\(H_2=0\\), \\(\\chi=1-b_1\\), so \\(b_1=2g+p-1=2\\cdot %d+%d-1=%d\\).", g, p, b1), f)
	default:
		chi := 2 - p
		return w20T(kp, fmt.Sprintf("The sphere is punctured at %d points. What is the Euler characteristic of the punctured surface?", p),
			strconv.Itoa(chi), fmt.Sprintf("A sphere has \\(\\chi=2\\) and each of the %d punctures subtracts one, giving \\(2-%d=%d\\).", p, p, chi), f)
	}
}

func w20SurfBoundary(kp int) generator.Problem {
	f := map[string]string{}
	g, b := w20R(0, 8), w20R(1, 12)
	w20I(f, "g", g)
	w20I(f, "b", b)
	switch kp {
	case 1:
		chi := 2 - 2*g - b
		return w20T(kp, fmt.Sprintf("An orientable surface of genus %d has %d boundary components. What is its Euler characteristic?", g, b),
			strconv.Itoa(chi), fmt.Sprintf("Boundary components subtract like punctures: \\(\\chi=2-2g-b=2-2\\cdot %d-%d=%d\\).", g, b, chi), f)
	case 2:
		b1 := 2*g + b - 1
		return w20T(kp, fmt.Sprintf("An orientable surface of genus %d with %d boundary components has the homotopy type of a wedge of circles. What is its first Betti number?", g, b),
			strconv.Itoa(b1), fmt.Sprintf("Since \\(\\chi=1-b_1\\) for this surface, \\(b_1=2g+b-1=2\\cdot %d+%d-1=%d\\).", g, b, b1), f)
	default:
		chi := 2 - b
		return w20T(kp, fmt.Sprintf("A sphere has %d open disks removed. What is the Euler characteristic of the resulting surface with boundary?", b),
			strconv.Itoa(chi), fmt.Sprintf("Each removed disk lowers \\(\\chi\\) by one, so \\(\\chi=2-b=2-%d=%d\\).", b, chi), f)
	}
}

func w20SurfDoubleCover(kp int) generator.Problem {
	f := map[string]string{}
	k := w20R(1, 20)
	w20I(f, "k", k)
	switch kp {
	case 1:
		genus := k - 1
		return w20T(kp, fmt.Sprintf("What is the genus of the orientable double cover of the nonorientable surface \\(N_{%d}\\)?", k),
			strconv.Itoa(genus), fmt.Sprintf("The double cover of \\(N_k\\) is orientable of genus \\(k-1=%d-1=%d\\).", k, genus), f)
	case 2:
		chi := 4 - 2*k
		return w20T(kp, fmt.Sprintf("The nonorientable surface \\(N_{%d}\\) has Euler characteristic \\(2-%d\\). What is the Euler characteristic of its orientable double cover?", k, k),
			strconv.Itoa(chi), fmt.Sprintf("A double cover multiplies \\(\\chi\\) by two: \\(2(2-k)=4-2\\cdot %d=%d\\).", k, chi), f)
	default:
		genus := k
		return w20T(kp, fmt.Sprintf("The orientable double cover of a nonorientable surface has genus %d. How many projective planes make up the original surface?", genus),
			strconv.Itoa(k), fmt.Sprintf("The cover of \\(N_k\\) has genus \\(k-1\\), so \\(k=%d+1=%d\\).", genus, k), f)
	}
}

func w20SurfGenusFormula(kp int) generator.Problem {
	f := map[string]string{}
	g, k := w20R(0, 8), w20R(1, 12)
	w20I(f, "g", g)
	w20I(f, "k", k)
	switch kp {
	case 1:
		chi := 2 - 2*g - k
		return w20T(kp, fmt.Sprintf("The orientable surface of genus %d is connected-summed with \\(N_{%d}\\). What is the Euler characteristic of the result?", g, k),
			strconv.Itoa(chi), fmt.Sprintf("Use \\(\\chi_1+\\chi_2-2=(2-2g)+(2-k)-2=2-2g-k=2-2\\cdot %d-%d=%d\\).", g, k, chi), f)
	case 2:
		genus := 2*g + k
		return w20T(kp, fmt.Sprintf("The orientable surface of genus %d is connected-summed with the connected sum of %d projective planes. What is the nonorientable genus of the result?", g, k),
			strconv.Itoa(genus), fmt.Sprintf("The sum is nonorientable with genus \\(2g+k=2\\cdot %d+%d=%d\\).", g, k, genus), f)
	default:
		chi := 2 - k
		return w20T(kp, fmt.Sprintf("What is the Euler characteristic of the connected sum of %d projective planes?", k),
			strconv.Itoa(chi), fmt.Sprintf("Each projective plane contributes one crosscap, so \\(\\chi=2-k=2-%d=%d\\).", k, chi), f)
	}
}

func w20SurfBoundaryConnect(kp int) generator.Problem {
	f := map[string]string{}
	g1, g2 := w20R(1, 8), w20R(1, 8)
	b1, b2 := w20R(1, 4), w20R(1, 4)
	w20I(f, "g1", g1)
	w20I(f, "g2", g2)
	w20I(f, "b1", b1)
	w20I(f, "b2", b2)
	switch kp {
	case 1:
		chi1, chi2 := 2-2*g1-b1, 2-2*g2-b2
		w20I(f, "chi1", chi1)
		w20I(f, "chi2", chi2)
		chi := chi1 + chi2 - 1
		return w20T(kp, fmt.Sprintf("Two surfaces with boundary have \\(\\chi=%d\\) and \\(\\chi=%d\\). They are joined by a boundary connected sum. What is the Euler characteristic?", chi1, chi2),
			strconv.Itoa(chi), fmt.Sprintf("A boundary connected sum glues one boundary interval from each, so \\(\\chi=\\chi_1+\\chi_2-1=%d+%d-1=%d\\).", chi1, chi2, chi), f)
	case 2:
		genus := g1 + g2
		return w20T(kp, fmt.Sprintf("A genus-%d surface with boundary is boundary-connected to a genus-%d surface with boundary. What is the genus of the result?", g1, g2),
			strconv.Itoa(genus), fmt.Sprintf("Boundary connected sum adds genera: \\(%d+%d=%d\\).", g1, g2, genus), f)
	default:
		chi := 2 - 2*g1 - b1
		return w20T(kp, fmt.Sprintf("What is the Euler characteristic of an orientable surface of genus %d with %d boundary components?", g1, b1),
			strconv.Itoa(chi), fmt.Sprintf("Use \\(\\chi=2-2g-b=2-2\\cdot %d-%d=%d\\).", g1, b1, chi), f)
	}
}

func w20SurfTriangulation(kp int) generator.Problem {
	f := map[string]string{}
	g := w20R(0, 6)
	chi := 2 - 2*g
	w20I(f, "g", g)
	w20I(f, "chi", chi)
	switch kp {
	case 1:
		v := w20R(5, 40)
		w20I(f, "V", v)
		edges := 3 * (v - chi)
		return w20T(kp, fmt.Sprintf("A closed orientable surface of Euler characteristic %d is triangulated with %d vertices and every face a triangle. How many edges does the triangulation have?", chi, v),
			strconv.Itoa(edges), fmt.Sprintf("Since \\(3F=2E\\), substitution into \\(V-E+F=\\chi\\) gives \\(E=3(V-\\chi)=3(%d-%d)=%d\\).", v, chi, edges), f)
	case 2:
		v := w20R(5, 40)
		w20I(f, "V", v)
		faces := 2 * (v - chi)
		return w20T(kp, fmt.Sprintf("A triangulation of a surface with \\(\\chi=%d\\) has %d vertices and all faces triangular. How many faces are there?", chi, v),
			strconv.Itoa(faces), fmt.Sprintf("From \\(3F=2E\\) and \\(V-E+F=\\chi\\) one gets \\(F=2(V-\\chi)=2(%d-%d)=%d\\).", v, chi, faces), f)
	default:
		start := 4 - chi
		if start < 4 {
			start = 4
		}
		t := start + w20R(0, 15)
		v := t + chi
		w20I(f, "V", v)
		edges := 3 * t
		w20I(f, "E", edges)
		return w20T(kp, fmt.Sprintf("A triangulated surface with \\(\\chi=%d\\) has %d edges. How many vertices does the triangulation have?", chi, edges),
			strconv.Itoa(v), fmt.Sprintf("Triangular faces give \\(E=3(V-\\chi)\\), so \\(V=E/3+\\chi=%d/3+%d=%d\\).", edges, chi, v), f)
	}
}

// ----- differential topology -----

func w20DiffTangentBundle(kp int) generator.Problem {
	f := map[string]string{}
	n := w20R(1, 12)
	w20I(f, "n", n)
	switch kp {
	case 1:
		return w20T(kp, fmt.Sprintf("What is the rank of the tangent bundle \\(TM\\) of a smooth %d-manifold?", n),
			strconv.Itoa(n), fmt.Sprintf("The fibre \\(T_pM\\) is an %d-dimensional vector space, so the bundle has rank %d.", n, n), f)
	case 2:
		dim := 2 * n
		return w20T(kp, fmt.Sprintf("The tangent bundle of a smooth %d-manifold is itself a smooth manifold. What is its dimension?", n),
			strconv.Itoa(dim), fmt.Sprintf("The total space has dimension \\(2n=2\\cdot %d=%d\\).", n, dim), f)
	default:
		dim := 2*n - 1
		return w20T(kp, fmt.Sprintf("What is the dimension of the unit tangent bundle of a smooth %d-manifold?", n),
			strconv.Itoa(dim), fmt.Sprintf("The unit tangent bundle is a sphere bundle of rank \\(n-1\\), so its total dimension is \\(n+(n-1)=2n-1=2\\cdot %d-1=%d\\).", n, dim), f)
	}
}

func w20DiffVectorField(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		n := w20R(2, 12)
		w20I(f, "n", n)
		index := 1 + (n % 2)
		if n%2 == 0 {
			index = 2
		} else {
			index = 0
		}
		// chi(S^n) = 1 + (-1)^n: 2 in even dimensions, 0 in odd.
		return w20T(kp, fmt.Sprintf("What is the sum of the indices of the zeros of a generic vector field on the sphere \\(S^{%d}\\)?", n),
			strconv.Itoa(index), fmt.Sprintf("By Poincaré-Hopf the index sum equals \\(\\chi(S^n)=1+(-1)^n\\), which is %d for \\(n=%d\\).", index, n), f)
	case 2:
		sources, saddles := w20R(1, 8), w20R(1, 8)
		w20I(f, "sources", sources)
		w20I(f, "saddles", saddles)
		sum := sources - saddles
		return w20T(kp, fmt.Sprintf("A vector field on a closed surface has %d zeros of index 1 and %d zeros of index -1, and no other zeros. What is the index sum?", sources, saddles),
			strconv.Itoa(sum), fmt.Sprintf("Add the indices: \\(%d\\cdot 1+%d\\cdot(-1)=%d\\).", sources, saddles, sum), f)
	default:
		g := w20R(1, 30)
		w20I(f, "g", g)
		sum := 2 - 2*g
		return w20T(kp, fmt.Sprintf("What is the sum of the indices of the zeros of a generic vector field on the closed orientable surface of genus %d?", g),
			strconv.Itoa(sum), fmt.Sprintf("The index sum equals \\(\\chi=2-2g=2-2\\cdot %d=%d\\).", g, sum), f)
	}
}

func w20DiffCriticalPoint(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		n := w20R(1, 20)
		w20I(f, "n", n)
		count := 1 << n
		return w20T(kp, fmt.Sprintf("What is the minimum number of critical points of a Morse function on the torus \\(T^{%d}\\)?", n),
			strconv.Itoa(count), fmt.Sprintf("A Morse function on \\(T^n\\) has at least \\(2^n\\) critical points, one for each cell of the torus, so \\(2^{%d}=%d\\).", n, count), f)
	case 2:
		g := w20R(1, 30)
		w20I(f, "g", g)
		count := 2*g + 2
		return w20T(kp, fmt.Sprintf("What is the minimum number of critical points of a Morse function on the closed orientable surface of genus %d?", g),
			strconv.Itoa(count), fmt.Sprintf("Such a surface needs at least \\(2g+2\\) critical points, one maximum, one minimum and two saddles per handle: \\(2\\cdot %d+2=%d\\).", g, count), f)
	default:
		n := w20R(2, 30)
		w20I(f, "n", n)
		count := 2
		return w20T(kp, fmt.Sprintf("What is the minimum number of critical points of the height function on the sphere \\(S^{%d}\\)?", n),
			strconv.Itoa(count), "A height function on a sphere has one maximum and one minimum, so the minimum is 2.", f)
	}
}

func w20DiffMorse(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		n := w20R(1, 20)
		w20I(f, "n", n)
		total := 1 << n
		return w20T(kp, fmt.Sprintf("What is the sum of the Betti numbers of the torus \\(T^{%d}\\)?", n),
			strconv.Itoa(total), fmt.Sprintf("The Betti numbers of \\(T^n\\) are \\(\\binom{n}{k}\\), which sum to \\(2^n=2^{%d}=%d\\).", n, total), f)
	case 2:
		g := w20R(1, 30)
		w20I(f, "g", g)
		total := 2 + 2*g
		return w20T(kp, fmt.Sprintf("What is the sum of the Betti numbers of the closed orientable surface of genus %d?", g),
			strconv.Itoa(total), fmt.Sprintf("The Betti numbers are \\(1, 2g, 1\\), whose sum is \\(2+2g=2+2\\cdot %d=%d\\).", g, total), f)
	default:
		n := w20R(1, 20)
		w20I(f, "n", n)
		count := n + 1
		return w20T(kp, fmt.Sprintf("A Morse function on a compact %d-manifold satisfies the Morse inequalities at how many degrees?", n),
			strconv.Itoa(count), fmt.Sprintf("There is one inequality for each degree from 0 to %d, so there are %d.", n, count), f)
	}
}

func w20DiffSard(kp int) generator.Problem {
	f := map[string]string{}
	n := w20R(2, 30)
	w20I(f, "n", n)
	switch kp {
	case 1:
		m := w20R(1, n)
		w20I(f, "m", m)
		dim := m - 1
		return w20T(kp, fmt.Sprintf("A smooth map \\(f:\\mathbb{R}^{%d}\\to\\mathbb{R}^{%d}\\) has critical set of dimension at most %d. What is the largest possible dimension of its set of critical values?", m, n, dim),
			strconv.Itoa(dim), fmt.Sprintf("Sard's theorem allows the critical values to have dimension at most \\(m-1=%d-1=%d\\).", m, dim), f)
	case 2:
		m := w20R(1, n)
		w20I(f, "m", m)
		codim := n - m + 1
		return w20T(kp, fmt.Sprintf("For a smooth map \\(f:\\mathbb{R}^{%d}\\to\\mathbb{R}^{%d}\\), the set of critical values has measure zero. What is its codimension in the target?", m, n),
			strconv.Itoa(codim), fmt.Sprintf("The image has dimension at most \\(m-1\\), so the codimension is at least \\(n-(m-1)=%d-%d+1=%d\\).", n, m, codim), f)
	default:
		dim := n
		return w20T(kp, fmt.Sprintf("For a smooth map \\(f:\\mathbb{R}^{%d}\\to\\mathbb{R}^{%d}\\), Sard's theorem says the set of critical values has measure zero. What is the dimension of the target space?", n, n),
			strconv.Itoa(dim), fmt.Sprintf("The target is \\(\\mathbb{R}^{%d}\\), of dimension %d.", n, n), f)
	}
}

func w20DiffTransversality(kp int) generator.Problem {
	f := map[string]string{}
	n := w20R(4, 12)
	w20I(f, "n", n)
	switch kp {
	case 1:
		a := w20R(1, n-1)
		b := w20R(1, n-1)
		w20I(f, "a", a)
		w20I(f, "b", b)
		dim := a + b - n
		if dim < 0 {
			dim = 0
		}
		return w20T(kp, fmt.Sprintf("Submanifolds \\(A\\) and \\(B\\) of dimensions %d and %d meet transversely in an %d-manifold. What is the expected dimension of \\(A\\cap B\\)?", a, b, n),
			strconv.Itoa(dim), fmt.Sprintf("Transversality gives \\(\\dim(A\\cap B)=\\dim A+\\dim B-n=%d+%d-%d=%d\\), and a negative value means a generic intersection is empty.", a, b, n, dim), f)
	case 2:
		a := w20R(1, n-1)
		w20I(f, "a", a)
		codim := n - a
		return w20T(kp, fmt.Sprintf("A submanifold \\(A\\) of dimension %d lies in an %d-manifold. What is its codimension?", a, n),
			strconv.Itoa(codim), fmt.Sprintf("Codimension is \\(n-\\dim A=%d-%d=%d\\).", n, a, codim), f)
	default:
		a := w20R(1, n-1)
		b := w20R(1, n-1)
		w20I(f, "a", a)
		w20I(f, "b", b)
		codim := 2*n - a - b
		return w20T(kp, fmt.Sprintf("Submanifolds of dimensions %d and %d in an %d-manifold meet transversely. What is the codimension of their intersection?", a, b, n),
			strconv.Itoa(codim), fmt.Sprintf("The codimensions add: \\((n-a)+(n-b)=2n-a-b=2\\cdot %d-%d-%d=%d\\).", n, a, b, codim), f)
	}
}

func w20DiffWhitney(kp int) generator.Problem {
	f := map[string]string{}
	n := w20R(1, 12)
	w20I(f, "n", n)
	switch kp {
	case 1:
		dim := 2 * n
		return w20T(kp, fmt.Sprintf("By the Whitney embedding theorem, every smooth %d-manifold embeds in some Euclidean space. What is the guaranteed embedding dimension?", n),
			strconv.Itoa(dim), fmt.Sprintf("The theorem gives an embedding into \\(\\mathbb{R}^{2n}\\), here \\(\\mathbb{R}^{%d}\\).", dim), f)
	case 2:
		dim := 2*n - 1
		return w20T(kp, fmt.Sprintf("By the Whitney immersion theorem, every smooth %d-manifold immerses in some Euclidean space. What is the guaranteed immersion dimension?", n),
			strconv.Itoa(dim), fmt.Sprintf("The theorem gives an immersion into \\(\\mathbb{R}^{2n-1}\\), here \\(\\mathbb{R}^{%d}\\).", dim), f)
	default:
		codim := n
		return w20T(kp, fmt.Sprintf("A smooth %d-manifold is embedded in \\(\\mathbb{R}^{2n}\\) by the Whitney theorem. What is its codimension in that embedding?", n),
			strconv.Itoa(codim), fmt.Sprintf("Codimension is \\(2n-n=n=%d\\).", codim), f)
	}
}

func w20DiffTubular(kp int) generator.Problem {
	f := map[string]string{}
	m := w20R(1, 8)
	n := m + w20R(1, 4)
	w20I(f, "m", m)
	w20I(f, "n", n)
	switch kp {
	case 1:
		rank := n - m
		return w20T(kp, fmt.Sprintf("An %d-dimensional submanifold sits in an %d-dimensional manifold. What is the rank of its normal bundle?", m, n),
			strconv.Itoa(rank), fmt.Sprintf("The normal bundle has rank equal to the codimension, \\(n-m=%d-%d=%d\\).", n, m, rank), f)
	case 2:
		dim := n
		return w20T(kp, fmt.Sprintf("A tubular neighbourhood of an %d-dimensional submanifold inside an %d-dimensional manifold has what dimension?", m, n),
			strconv.Itoa(dim), fmt.Sprintf("A tubular neighbourhood is open in the ambient manifold, so its dimension is the ambient dimension %d.", n), f)
	default:
		dim := n
		return w20T(kp, fmt.Sprintf("A normal bundle of rank %d sits over an %d-dimensional submanifold. What is the dimension of the total space of the bundle?", n-m, m),
			strconv.Itoa(dim), fmt.Sprintf("The total space has dimension \\(m+\\operatorname{rank}=%d+(%d-%d)=%d\\).", m, n, m, dim), f)
	}
}

func w20DiffDegree(kp int) generator.Problem {
	f := map[string]string{}
	k := w20R(1, 12)
	w20I(f, "k", k)
	switch kp {
	case 1:
		return w20T(kp, fmt.Sprintf("What is the degree of the circle map \\(z\\mapsto z^{%d}\\)?", k),
			strconv.Itoa(k), fmt.Sprintf("The map wraps the circle %d times, so its degree is %d.", k, k), f)
	case 2:
		n := w20R(1, 30)
		w20I(f, "n", n)
		degree := 1
		if n%2 == 0 {
			degree = -1
		}
		return w20T(kp, fmt.Sprintf("What is the degree of the antipodal map on the sphere \\(S^{%d}\\)?", n),
			strconv.Itoa(degree), fmt.Sprintf("The antipodal map on \\(S^n\\) has degree \\((-1)^{n+1}\\), which is %d for \\(n=%d\\).", degree, n), f)
	default:
		degree := -k
		return w20T(kp, fmt.Sprintf("What is the degree of the circle map \\(z\\mapsto z^{-%d}\\)?", k),
			strconv.Itoa(degree), fmt.Sprintf("Reversing orientation negates the winding, so the degree is \\(-%d\\).", k), f)
	}
}

func w20DiffFlow(kp int) generator.Problem {
	f := map[string]string{}
	n := w20R(2, 30)
	w20I(f, "n", n)
	switch kp {
	case 1:
		return w20T(kp, fmt.Sprintf("A smooth vector field on an %d-manifold generates a flow. What is the dimension of the phase space of the flow?", n),
			strconv.Itoa(n), fmt.Sprintf("The flow moves points inside the %d-manifold, so the phase space has dimension %d.", n, n), f)
	case 2:
		k := w20R(1, n-1)
		w20I(f, "k", k)
		return w20T(kp, fmt.Sprintf("A linear vector field on \\(\\mathbb{R}^{%d}\\) has %d positive eigenvalues and no zero eigenvalues. What is the dimension of its unstable manifold?", n, k),
			strconv.Itoa(k), fmt.Sprintf("The unstable manifold is spanned by the eigenvectors of the positive eigenvalues, so its dimension is %d.", k), f)
	default:
		k := w20R(1, n-1)
		w20I(f, "k", k)
		dim := n - k
		return w20T(kp, fmt.Sprintf("A linear vector field on \\(\\mathbb{R}^{%d}\\) has %d positive eigenvalues and no zero eigenvalues. What is the dimension of its stable manifold?", n, k),
			strconv.Itoa(dim), fmt.Sprintf("The remaining \\(n-k=%d-%d=%d\\) eigenvalues are negative, which span the stable manifold.", n, k, dim), f)
	}
}
