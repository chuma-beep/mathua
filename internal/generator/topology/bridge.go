package topology

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

var topoBridgeKinds = map[string]string{
	"topo.basics.spaces":                    "spaces",
	"topo.connected.components":             "components",
	"topo.separation.normal":                "normal",
	"topo.algebraic.homotopy_equivalence":   "homotopy_equivalence",
	"topo.covering.universal":               "universal_cover",
	"topo.manifolds.boundary":               "manifold_boundary",
	"topo.invariant.surface_classification": "surfaces",
	"topo.compact.tychonoff":                "tychonoff",
	"topo.algebraic.degree":                 "degree",
	"topo.algebraic.mayer_vietoris":         "mayer_vietoris",
}

type topoBridgeGen struct{ kind string }

func (g *topoBridgeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := topoBridgeOrdinal(ctx)
	switch g.kind {
	case "spaces":
		return topoSpacesProblem(o)
	case "components":
		return topoComponentsProblem(o)
	case "normal":
		return topoNormalProblem(o)
	case "homotopy_equivalence":
		return topoHomotopyEquivalenceProblem(o)
	case "universal_cover":
		return topoUniversalCoverProblem(o)
	case "manifold_boundary":
		return topoManifoldBoundaryProblem(o)
	case "surfaces":
		return topoSurfacesProblem(o)
	case "tychonoff":
		return topoTychonoffProblem(o)
	case "degree":
		return topoDegreeProblem(o)
	case "mayer_vietoris":
		return topoMayerVietorisProblem(o)
	}
	return topoBridgeProblem(o, "What is 2 + 2? (enter a number)", "4", "Two plus two is four.")
}

func topoBridgeOrdinal(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return 1 + rand.Intn(3)
}

func topoBridgeProblem(o int, question, answer, explanation string) generator.Problem {
	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts:       map[string]string{"answer": answer},
		KPOrdinal:   o,
	}
}

func topoBridgePow2(n int) int {
	v := 1
	for i := 0; i < n; i++ {
		v *= 2
	}
	return v
}

type topoBridgeItem struct {
	label string
	yes   bool
	note  string
}

func topoBridgeCount(o int, lead string, pool []topoBridgeItem) generator.Problem {
	const draw = 5
	var yes, no []topoBridgeItem
	for _, it := range pool {
		if it.yes {
			yes = append(yes, it)
		} else {
			no = append(no, it)
		}
	}
	count := rand.Intn(draw + 1)
	if count > len(yes) {
		count = len(yes)
	}
	if draw-count > len(no) {
		count = draw - len(no)
	}
	rand.Shuffle(len(yes), func(i, j int) { yes[i], yes[j] = yes[j], yes[i] })
	rand.Shuffle(len(no), func(i, j int) { no[i], no[j] = no[j], no[i] })
	chosen := append([]topoBridgeItem{}, yes[:count]...)
	chosen = append(chosen, no[:draw-count]...)
	rand.Shuffle(len(chosen), func(i, j int) { chosen[i], chosen[j] = chosen[j], chosen[i] })
	var parts, notes []string
	for i, it := range chosen {
		parts = append(parts, fmt.Sprintf("(%c) %s", 'a'+i, it.label))
		verdict := "no"
		if it.yes {
			verdict = "yes"
		}
		notes = append(notes, fmt.Sprintf("%s: %s, %s", it.label, verdict, it.note))
	}
	q := fmt.Sprintf("%s\n%s\n(enter a number)", lead, strings.Join(parts, "\n"))
	ex := fmt.Sprintf("Exactly %d of the %d items satisfy the condition.\n%s", count, draw, strings.Join(notes, "\n"))
	return topoBridgeProblem(o, q, strconv.Itoa(count), ex)
}

func topoSpacesProblem(o int) generator.Problem {
	switch o {
	case 1:
		n := 2 + rand.Intn(7)
		return topoBridgeProblem(o, fmt.Sprintf("The discrete topology on a set with %d points has how many open sets? (enter a number)", n), strconv.Itoa(topoBridgePow2(n)), fmt.Sprintf("Every subset is open, and a set with %d points has %d subsets.", n, topoBridgePow2(n)))
	case 2:
		n := 3 + rand.Intn(7)
		ans := topoBridgePow2(n-1) + 1
		return topoBridgeProblem(o, fmt.Sprintf("For the particular point topology on a set with %d points, how many open sets are there? (enter a number)", n), strconv.Itoa(ans), fmt.Sprintf("A set is open when it is empty or contains the particular point. There are %d sets containing that point, plus the empty set, so %d open sets.", topoBridgePow2(n-1), ans))
	default:
		return topoBridgeCount(o, "How many of these collections of subsets of a three-point set are topologies?", []topoBridgeItem{
			{"the indiscrete collection", true, "contains the empty set and the whole set and is closed under unions and intersections"},
			{"the discrete collection", true, "contains every subset"},
			{"the particular point collection", true, "empty sets and sets containing one fixed point"},
			{"the excluded point collection", true, "the whole set and sets omitting one fixed point"},
			{"the collection of all subsets of size at most one, with the whole set", true, "closed under finite intersections and arbitrary unions"},
			{"a partition topology", true, "unions of blocks of a partition"},
			{"a collection closed under unions and finite intersections, with empty set and whole set", true, "this is the definition of a topology"},
			{"the collection obtained from a basis by taking all unions", true, "the topology generated by a basis"},
			{"a collection missing the whole set", false, "a topology must contain the whole set"},
			{"a collection missing the empty set", false, "a topology must contain the empty set"},
			{"the collection containing only the empty set and one singleton", false, "the whole set is missing"},
			{"the collection of two singletons and the empty set", false, "the union of the singletons is the whole set, which is missing"},
			{"a collection not closed under union", false, "arbitrary unions must remain in the collection"},
			{"a collection not closed under intersection", false, "finite intersections must remain in the collection"},
			{"the collection of all nonempty subsets", false, "the empty set is missing"},
			{"the collection of all proper subsets", false, "the whole set is missing"},
		})
	}
}

func topoComponentsProblem(o int) generator.Problem {
	switch o {
	case 1:
		k := 2 + rand.Intn(7)
		return topoBridgeProblem(o, fmt.Sprintf("A space is the disjoint union of %d nonempty connected open intervals. How many connected components does it have? (enter a number)", k), strconv.Itoa(k), fmt.Sprintf("Each interval is connected, and disjoint nonempty intervals cannot lie in one component, so there are %d components.", k))
	case 2:
		n := 1 + rand.Intn(8)
		return topoBridgeProblem(o, fmt.Sprintf("How many connected components does the real line with %d distinct points removed have? (enter a number)", n), strconv.Itoa(n+1), fmt.Sprintf("Removing %d points cuts the line into %d open intervals, one more than the number of removed points.", n, n+1))
	default:
		k := 2 + rand.Intn(4)
		m := 2 + rand.Intn(4)
		return topoBridgeProblem(o, fmt.Sprintf("A space has %d connected components and another has %d. How many connected components does their product have? (enter a number)", k, m), strconv.Itoa(k*m), fmt.Sprintf("A product is connected exactly when both factors are connected, so the component count multiplies: %d times %d is %d.", k, m, k*m))
	}
}

func topoNormalProblem(o int) generator.Problem {
	switch o {
	case 1:
		return topoBridgeCount(o, "How many of these spaces are normal?", []topoBridgeItem{
			{"the real line with the usual topology", true, "metric spaces are normal"},
			{"the closed interval", true, "a compact Hausdorff space is normal"},
			{"Euclidean n-space", true, "metric spaces are normal"},
			{"the Cantor set", true, "compact Hausdorff spaces are normal"},
			{"a discrete space", true, "every subset is open and closed, so disjoint closed sets separate"},
			{"the integers with the usual metric", true, "metric spaces are normal"},
			{"any compact Hausdorff space", true, "compact Hausdorff implies normal"},
			{"the real line with the cofinite topology", true, "a finite closed family separates because the space is T1 and the complement of a finite set is open"},
			{"the Sorgenfrey plane", false, "the product of two Sorgenfrey lines is the standard nonnormal example"},
			{"the rationals with the usual topology", false, "the rationals are not a Baire space and normality fails in the usual proof"},
			{"the irrationals with the usual topology", false, "the irrationals fail to be normal in the subspace topology"},
			{"the Niemytzki plane", false, "this is a classical nonnormal Moore plane"},
			{"the space of countable ordinals with the order topology", false, "the product with itself is not normal"},
			{"the real line with the particular point topology", false, "the particular point topology is not Hausdorff and not normal"},
			{"the real line with the excluded point topology", false, "the excluded point topology is not Hausdorff and not normal"},
			{"the Sorgenfrey line squared with the box topology", false, "the box product fails normality"},
		})
	case 2:
		return topoBridgeCount(o, "How many of these statements about normal spaces are true?", []topoBridgeItem{
			{"disjoint closed sets have disjoint open neighbourhoods", true, "this is the definition"},
			{"a normal T1 space admits Urysohn functions", true, "Urysohn's lemma"},
			{"metric spaces are normal", true, "metric spaces are perfectly normal"},
			{"compact Hausdorff spaces are normal", true, "a standard separation theorem"},
			{"closed subspaces of normal spaces are normal", true, "normality is inherited by closed subspaces"},
			{"normal spaces admit Tietze extensions", true, "Tietze's extension theorem holds for normal spaces"},
			{"finite products of normal spaces are normal", false, "normality is not finitely productive"},
			{"arbitrary products of normal spaces are normal", false, "normality is not productive"},
			{"every normal space is Hausdorff", false, "normality alone does not imply Hausdorff"},
			{"every normal space is metrizable", false, "many nonmetrizable normal spaces exist"},
			{"open subspaces of normal spaces are normal", false, "normality is not hereditary for open subspaces"},
			{"quotients of normal spaces are normal", false, "normality is not quotient closed"},
			{"normal spaces are compact", false, "normality and compactness are independent"},
			{"normal spaces are connected", false, "a two-point discrete space is normal and disconnected"},
			{"normal spaces have countable bases", false, "second countability is a stronger assumption"},
			{"normal spaces are path connected", false, "normality says nothing about paths"},
		})
	default:
		return topoBridgeCount(o, "How many of these implications about normality are true?", []topoBridgeItem{
			{"metric implies normal", true, "every metric space is normal"},
			{"compact Hausdorff implies normal", true, "compact Hausdorff spaces are normal"},
			{"normal T1 implies completely regular", true, "normal T1 spaces are T4 and hence completely regular"},
			{"closed subspace of normal implies normal", true, "closed subspaces inherit normality"},
			{"normal plus second countable implies metrizable", true, "the Urysohn metrization theorem"},
			{"normal T1 with a countable base is separable", true, "second countability implies separability"},
			{"normal implies Hausdorff", false, "the indiscrete two-point space is normal and not Hausdorff"},
			{"normal implies compact", false, "the real line is normal and not compact"},
			{"normal implies second countable", false, "the Sorgenfrey line is normal and not second countable"},
			{"normal implies connected", false, "a discrete two-point space is normal and disconnected"},
			{"normal implies path connected", false, "normality does not control paths"},
			{"normal implies locally compact", false, "the rationals are not locally compact"},
			{"normal implies metrizable", false, "the Sorgenfrey line is a counterexample"},
			{"normal implies paracompact", false, "paracompactness is a separate hypothesis"},
			{"normal implies Lindelof", false, "the Sorgenfrey line is normal and not Lindelof"},
			{"normal implies separable", false, "the Sorgenfrey line is normal and not separable"},
		})
	}
}

func topoHomotopyEquivalenceProblem(o int) generator.Problem {
	switch o {
	case 1:
		return topoBridgeCount(o, "How many of these spaces are contractible?", []topoBridgeItem{
			{"the real line", true, "it deformation retracts to a point"},
			{"Euclidean n-space", true, "it is convex"},
			{"the closed interval", true, "it is convex"},
			{"the closed unit disk", true, "it is convex"},
			{"a single point", true, "the identity is already constant"},
			{"a cone on any space", true, "a cone deformation retracts to its apex"},
			{"a tree", true, "a tree deformation retracts to any vertex"},
			{"the closed unit ball in Euclidean space", true, "it is convex"},
			{"the circle", false, "its fundamental group is the integers"},
			{"the two-sphere", false, "its second homology group is nonzero"},
			{"the torus", false, "its fundamental group is free abelian of rank two"},
			{"the real projective plane", false, "its fundamental group has order two"},
			{"a wedge of two circles", false, "its fundamental group is free of rank two"},
			{"two disjoint points", false, "it is disconnected"},
			{"the annulus", false, "it has the homotopy type of a circle"},
			{"the Möbius band", false, "it deformation retracts to a circle"},
		})
	case 2:
		return topoBridgeCount(o, "How many of these pairs of spaces are homotopy equivalent?", []topoBridgeItem{
			{"the real line and a point", true, "the line is contractible"},
			{"the closed interval and a point", true, "the interval is contractible"},
			{"the circle and the punctured plane", true, "the punctured plane deformation retracts to the circle"},
			{"the circle and the annulus", true, "the annulus deformation retracts to its core circle"},
			{"the circle and the cylinder", true, "the cylinder deformation retracts to a circle"},
			{"the Möbius band and the circle", true, "the band deformation retracts to its core"},
			{"the two-sphere and the punctured three-space", true, "punctured three-space deformation retracts to the sphere"},
			{"the torus and the product of two circles", true, "they are homeomorphic"},
			{"the circle and the two-sphere", false, "their fundamental groups differ"},
			{"the torus and the circle", false, "their fundamental groups differ"},
			{"the real projective plane and the two-sphere", false, "their fundamental groups differ"},
			{"a wedge of two circles and a circle", false, "their fundamental groups have different ranks"},
			{"a wedge of two circles and a wedge of three circles", false, "the ranks differ"},
			{"the circle and the closed interval", false, "one has a nontrivial fundamental group"},
			{"two points and one point", false, "their component counts differ"},
			{"the two-sphere and the complement of a point in three-space", true, "the complement deformation retracts to a sphere"},
		})
	default:
		return topoBridgeCount(o, "How many of these properties are preserved by homotopy equivalence?", []topoBridgeItem{
			{"the fundamental group", true, "homotopy equivalences induce isomorphisms"},
			{"the first homology group", true, "homology is homotopy invariant"},
			{"the cohomology ring", true, "cohomology is homotopy invariant"},
			{"the Euler characteristic", true, "it is computed from homology"},
			{"connectedness", true, "a homotopy equivalence is continuous and has a continuous inverse up to homotopy"},
			{"path connectedness", true, "paths can be transported along the equivalence"},
			{"compact support cohomology for proper maps", true, "proper homotopy equivalences preserve it"},
			{"the higher homotopy groups", true, "they are homotopy invariants"},
			{"compactness", false, "the real line is homotopy equivalent to a point"},
			{"the metric diameter", false, "the line and a point have different diameters"},
			{"cardinality", false, "homotopy equivalence need not be injective"},
			{"dimension", false, "the interval and a point have different dimensions"},
			{"metric completeness", false, "the open interval is homotopy equivalent to a point"},
			{"the property of being a manifold", false, "the wedge of two circles is not a manifold"},
			{"the property of being a group", false, "homotopy equivalence is not an algebraic operation"},
			{"the exact value of the distance function", false, "homotopy ignores metric scale"},
		})
	}
}

func topoUniversalCoverProblem(o int) generator.Problem {
	switch o {
	case 1:
		return topoBridgeCount(o, "How many of these universal cover descriptions are correct?", []topoBridgeItem{
			{"the universal cover of the circle is the real line", true, "the exponential map is the covering"},
			{"the universal cover of the torus is the plane", true, "the plane covers the torus by the lattice action"},
			{"the universal cover of the two-sphere is the two-sphere", true, "the sphere is simply connected"},
			{"the universal cover of the projective plane is the two-sphere", true, "the antipodal quotient gives the projective plane"},
			{"the universal cover of the real projective n-space is the n-sphere", true, "the double cover is universal for n at least two"},
			{"the universal cover of a simply connected space is itself", true, "the identity is a universal cover"},
			{"the universal cover of the punctured plane is the plane", true, "the exponential map covers the punctured plane"},
			{"the universal cover of the open annulus is the plane", true, "the annulus has the homotopy type of a circle"},
			{"the universal cover of the circle is the circle", false, "the circle is not simply connected"},
			{"the universal cover of the torus is the torus", false, "the torus is not simply connected"},
			{"the universal cover of the projective plane is the projective plane", false, "the projective plane is not simply connected"},
			{"the universal cover of the two-sphere is the plane", false, "the sphere is already simply connected"},
			{"the universal cover of the Klein bottle is the Klein bottle", false, "the Klein bottle is not simply connected"},
			{"the universal cover of the Möbius band is the Möbius band", false, "the band is not simply connected"},
			{"the universal cover of a point is the circle", false, "a point is simply connected"},
			{"the universal cover of the real line is the circle", false, "the covering map goes the other way"},
		})
	case 2:
		return topoBridgeCount(o, "How many of these statements about covering spaces are true?", []topoBridgeItem{
			{"a covering map is a local homeomorphism", true, "this is part of the definition"},
			{"the fibres of a covering map are discrete", true, "evenly covered neighbourhoods give discrete fibres"},
			{"a universal cover is simply connected", true, "this is the defining property"},
			{"the universal cover exists for connected locally path connected semilocally simply connected spaces", true, "the standard existence theorem"},
			{"deck transformations of a universal cover form a group isomorphic to the fundamental group", true, "the deck group is the fundamental group"},
			{"a covering of a connected base has fibres of the same cardinality", true, "the sheets are locally constant"},
			{"the universal cover of a product is the product of the universal covers", true, "covering maps are compatible with products"},
			{"a connected cover of a simply connected space is a homeomorphism", true, "the fundamental group is trivial"},
			{"every local homeomorphism is a covering map", false, "the inclusion of an open interval into the plane is a counterexample"},
			{"every covering map has finite fibres", false, "the real line covers the circle with infinite fibres"},
			{"every space has a universal cover", false, "the Hawaiian earring is a standard counterexample"},
			{"a covering map is injective", false, "the exponential map is not injective"},
			{"the universal cover of the circle is the circle", false, "the circle is not simply connected"},
			{"coverings preserve the fundamental group", false, "the induced map is injective, not an isomorphism in general"},
			{"every covering is regular", false, "regular coverings correspond to normal subgroups"},
			{"the universal cover of the projective plane is the circle", false, "it is the two-sphere"},
		})
	default:
		return topoBridgeCount(o, "How many of these statements about deck transformations are true?", []topoBridgeItem{
			{"a deck transformation commutes with the covering map", true, "this is the definition"},
			{"deck transformations form a group under composition", true, "the composition of two deck transformations is a deck transformation"},
			{"the deck group acts freely on a universal cover", true, "a nontrivial deck transformation has no fixed point"},
			{"the deck group of a universal cover is isomorphic to the fundamental group", true, "the standard correspondence"},
			{"the quotient by the deck group recovers the base", true, "the covering map is the quotient map"},
			{"a regular covering corresponds to a normal subgroup", true, "the deck group acts transitively on fibres exactly then"},
			{"the deck group of the universal cover of the circle is the integers", true, "translations by integers"},
			{"the deck group of the universal cover of the torus is the integer lattice", true, "translations by lattice vectors"},
			{"every homeomorphism of the cover is a deck transformation", false, "it must commute with the covering map"},
			{"the deck group is always trivial", false, "the circle has a nontrivial deck group"},
			{"the deck group is always finite", false, "the universal cover of the circle has infinite deck group"},
			{"the deck group is always abelian", false, "it is the fundamental group, which need not be abelian"},
			{"deck transformations preserve the fibre pointwise", false, "a free action moves points between sheets"},
			{"the deck group of the universal cover of the two-sphere is the integers", false, "the sphere is simply connected"},
			{"the deck group of the universal cover of the projective plane is trivial", false, "it has order two"},
			{"every covering is universal", false, "the identity cover of a non-simply-connected space is not universal"},
		})
	}
}

func topoManifoldBoundaryProblem(o int) generator.Problem {
	switch o {
	case 1:
		n := 1 + rand.Intn(6)
		return topoBridgeProblem(o, fmt.Sprintf("The closed unit cube in %d-dimensional space is a manifold with boundary of what dimension? (enter a number)", n), strconv.Itoa(n), fmt.Sprintf("The cube is locally homeomorphic to a half-space in %d dimensions, so its dimension is %d.", n, n))
	case 2:
		n := 1 + rand.Intn(8)
		return topoBridgeProblem(o, fmt.Sprintf("How many boundary faces does the closed unit cube in %d-dimensional space have? (enter a number)", n), strconv.Itoa(2*n), fmt.Sprintf("There are two faces for each coordinate direction, so the cube has 2 times %d, which is %d, boundary faces.", n, 2*n))
	default:
		n := 2 + rand.Intn(7)
		return topoBridgeProblem(o, fmt.Sprintf("The boundary of the closed unit ball in %d-dimensional space is a sphere of what dimension? (enter a number)", n), strconv.Itoa(n-1), fmt.Sprintf("The ball is %d-dimensional, so its boundary has dimension one less, namely %d.", n, n-1))
	}
}

func topoSurfacesProblem(o int) generator.Problem {
	switch o {
	case 1:
		g := rand.Intn(9)
		chi := 2 - 2*g
		return topoBridgeProblem(o, fmt.Sprintf("What is the Euler characteristic of the closed orientable surface of genus %d? (enter a number)", g), strconv.Itoa(chi), fmt.Sprintf("The Euler characteristic of an orientable genus-%d surface is 2 minus 2 times %d, which is %d.", g, g, chi))
	case 2:
		k := 1 + rand.Intn(9)
		chi := 2 - k
		return topoBridgeProblem(o, fmt.Sprintf("What is the Euler characteristic of the closed nonorientable surface with %d crosscaps? (enter a number)", k), strconv.Itoa(chi), fmt.Sprintf("The Euler characteristic of a nonorientable surface with %d crosscaps is 2 minus %d, which is %d.", k, k, chi))
	default:
		g := rand.Intn(6)
		h := rand.Intn(6)
		return topoBridgeProblem(o, fmt.Sprintf("The connected sum of an orientable surface of genus %d and one of genus %d has what genus? (enter a number)", g, h), strconv.Itoa(g+h), fmt.Sprintf("Connected sum adds genera, so the result has genus %d plus %d, which is %d.", g, h, g+h))
	}
}

func topoTychonoffProblem(o int) generator.Problem {
	switch o {
	case 1:
		return topoBridgeCount(o, "How many of these products are compact?", []topoBridgeItem{
			{"a finite product of closed intervals", true, "Tychonoff's theorem"},
			{"an arbitrary product of closed intervals", true, "Tychonoff's theorem for arbitrary products"},
			{"the Hilbert cube", true, "a countable product of compact intervals"},
			{"a product of compact Hausdorff spaces", true, "compactness is productive"},
			{"the product of two circles", true, "the torus is compact"},
			{"the product of finitely many compact spaces", true, "the finite case is elementary"},
			{"the product of countably many compact spaces", true, "the countable case follows from Tychonoff"},
			{"the product of any family of compact spaces", true, "the full theorem"},
			{"an infinite product of open intervals", false, "open intervals are not compact"},
			{"the product of the real line with itself", false, "the plane is not compact"},
			{"an infinite product of noncompact spaces", false, "a noncompact factor makes the product noncompact"},
			{"the product of countably many copies of the real line", false, "the product is not compact"},
			{"the box product of countably many compact spaces", false, "the box topology is finer and may fail compactness"},
			{"the product of a compact space with an open interval", false, "the open interval is not compact"},
			{"the product of the integers with the usual topology and a point", false, "the integers are not compact"},
			{"the product of two open disks", false, "open disks are not compact"},
		})
	case 2:
		return topoBridgeCount(o, "How many of these statements about Tychonoff's theorem are true?", []topoBridgeItem{
			{"an arbitrary product of compact spaces is compact", true, "the theorem"},
			{"the product topology is the topology used in the theorem", true, "the box topology does not work"},
			{"the theorem is equivalent to the axiom of choice", true, "over the usual axioms of set theory"},
			{"the finite case is elementary", true, "it follows from the tube lemma"},
			{"the theorem applies to profinite groups", true, "they are inverse limits of finite groups"},
			{"Stone spaces are compact by the theorem", true, "they are closed subspaces of products of two-point spaces"},
			{"the Hilbert cube is compact", true, "a countable product of intervals"},
			{"the theorem gives compactness of arbitrary products of spheres", true, "spheres are compact"},
			{"the theorem holds for the box topology", false, "the box topology is finer"},
			{"the theorem requires the product to be finite", false, "the theorem is specifically about arbitrary products"},
			{"the theorem requires metrizability", false, "no metric is needed"},
			{"the theorem fails for uncountable products", false, "it holds for arbitrary products"},
			{"the theorem is a consequence of compactness of the factors", false, "it is the other way round for products"},
			{"the theorem applies only to Hausdorff spaces", false, "compactness alone is enough"},
			{"the theorem implies every compact space is metrizable", false, "compactness does not imply metrizability"},
			{"the theorem is false for products of two compact spaces", false, "the finite case is true"},
		})
	default:
		return topoBridgeCount(o, "How many of these statements about the finite intersection property are true?", []topoBridgeItem{
			{"compactness is equivalent to the finite intersection property for closed families", true, "the standard criterion"},
			{"a family with the finite intersection property has nonempty intersection in a compact space", true, "this is the criterion"},
			{"the finite intersection property concerns finite subfamilies", true, "only finite intersections are required to be nonempty"},
			{"Tychonoff's theorem can be proved with the finite intersection property", true, "the usual proof extends a family and uses coordinate cluster points"},
			{"a nested sequence of nonempty closed sets in a compact space has nonempty intersection", true, "it is a family with the finite intersection property"},
			{"the criterion applies to arbitrary families of closed sets", true, "no countability is needed"},
			{"the criterion fails in noncompact spaces", false, "the real line has a nested family of closed rays with empty intersection"},
			{"the finite intersection property alone implies compactness", false, "the space must be compact for the conclusion"},
			{"the criterion requires the family to be countable", false, "arbitrary families are allowed"},
			{"the criterion requires open sets", false, "it is stated for closed sets"},
			{"the criterion is equivalent to connectedness", false, "it characterises compactness"},
			{"the criterion fails for Hausdorff spaces", false, "it holds in every compact space"},
			{"the criterion needs the sets to be nested", false, "nesting is only a special case"},
			{"the criterion needs the sets to be open", false, "closed sets are used"},
			{"the criterion is false in compact metric spaces", false, "compact metric spaces are compact"},
			{"the criterion only applies to finite spaces", false, "it applies to all compact spaces"},
		})
	}
}

func topoDegreeProblem(o int) generator.Problem {
	switch o {
	case 1:
		k := 1 + rand.Intn(9)
		return topoBridgeProblem(o, fmt.Sprintf("The map from the circle to itself sending z to z to the %d has what degree? (enter a number)", k), strconv.Itoa(k), fmt.Sprintf("The map wraps the circle %d times in the positive direction, so its degree is %d.", k, k))
	case 2:
		k := 1 + rand.Intn(5)
		m := 1 + rand.Intn(5)
		return topoBridgeProblem(o, fmt.Sprintf("A map of the circle has degree %d and another has degree %d. What is the degree of their composition? (enter a number)", k, m), strconv.Itoa(k*m), fmt.Sprintf("Degree is multiplicative under composition: %d times %d is %d.", k, m, k*m))
	default:
		k := 1 + rand.Intn(8)
		return topoBridgeProblem(o, fmt.Sprintf("The map from the circle to itself sending z to the conjugate of z to the %d has what degree? (enter a number)", k), strconv.Itoa(-k), fmt.Sprintf("Conjugation reverses orientation, so the degree is the negative of %d, which is %d.", k, -k))
	}
}

func topoMayerVietorisProblem(o int) generator.Problem {
	switch o {
	case 1:
		a := rand.Intn(9)
		b := rand.Intn(9)
		c := rand.Intn(9)
		return topoBridgeProblem(o, fmt.Sprintf("A space X is the union of A and B. The Euler characteristics are chi(A)=%d, chi(B)=%d and chi(A intersect B)=%d. What is chi(X)? (enter a number)", a, b, c), strconv.Itoa(a+b-c), fmt.Sprintf("Euler characteristic is additive over the union with the intersection subtracted: %d plus %d minus %d is %d.", a, b, c, a+b-c))
	case 2:
		k := 1 + rand.Intn(6)
		return topoBridgeProblem(o, fmt.Sprintf("A and B are contractible and their intersection has the homology of the %d-sphere. What is the rank of the reduced homology of their union in degree %d? (enter a number)", k, k), "1", fmt.Sprintf("The Mayer-Vietoris sequence identifies the reduced homology of the union in degree %d with the reduced homology of the %d-sphere, which has rank 1.", k, k))
	default:
		return topoBridgeCount(o, "How many of these statements about the Mayer-Vietoris sequence are true?", []topoBridgeItem{
			{"it relates the homology of a union to the homology of the pieces", true, "this is the purpose of the sequence"},
			{"it applies to an open cover of the space", true, "the interiors must cover the space"},
			{"the connecting homomorphism lowers degree by one", true, "it goes from H_n of the union to H_{n-1} of the intersection"},
			{"it is natural with respect to maps of covers", true, "the sequence is natural"},
			{"it computes the homology of the torus from two cylinders", true, "a standard example"},
			{"it computes the homology of a wedge", true, "the reduced homology is the direct sum"},
			{"the sequence is exact", true, "it is a long exact sequence"},
			{"it works for any homology theory satisfying the axioms", true, "the construction is axiomatic"},
			{"it applies to arbitrary closed covers", false, "the sets need to be excisive or open"},
			{"the connecting map raises degree by one", false, "in homology it lowers degree"},
			{"the sequence is short and has only three terms", false, "it is a long exact sequence"},
			{"it computes homotopy groups directly", false, "it is a homology sequence"},
			{"it requires the intersection to be empty", false, "the intersection term is essential"},
			{"it requires the pieces to be contractible", false, "the pieces can be arbitrary"},
			{"it only works for surfaces", false, "it works in all dimensions"},
			{"it is independent of the cover", false, "the maps depend on the chosen cover"},
		})
	}
}
