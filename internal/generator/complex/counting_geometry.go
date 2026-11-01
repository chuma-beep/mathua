package complex

import "github.com/chuma-beep/mathua/internal/generator"

func conformalProblem() generator.Problem {
	pool := []complexVerdict{
		{"an analytic function with nonzero derivative is conformal", true, "the derivative gives an invertible linear map"},
		{"w = 1/z is conformal on C minus {0}", true, "its derivative is -1/z^2"},
		{"w = z^2 is conformal on C minus {0}", true, "the derivative 2z vanishes only at 0"},
		{"w = e^z is conformal everywhere", true, "its derivative never vanishes"},
		{"every Mobius map is conformal", true, "the derivative is nonzero"},
		{"a conformal map preserves angles between curves", true, "the definition"},
		{"a conformal map is locally invertible", true, "the inverse function theorem"},
		{"an analytic map with f'(a) nonzero is conformal at a", true, "the local criterion"},
		{"w = z^2 is conformal at 0", false, "the derivative vanishes there"},
		{"w = z bar is conformal", false, "it reverses orientation"},
		{"w = |z|^2 is conformal", false, "it is not analytic away from 0"},
		{"w = 1/z is conformal at 0", false, "it has a pole there"},
		{"every analytic function is conformal", false, "the derivative may vanish"},
		{"every conformal map is constant", false, "Mobius maps are nonconstant"},
		{"w = Re z is conformal", false, "it is real valued"},
		{"w = z^3 is conformal at 0", false, "the derivative vanishes there"},
	}
	return countComplexVerdicts("How many of these statements about conformal maps are true?", pool)
}

func cauchyGoursatProblem() generator.Problem {
	pool := []complexVerdict{
		{"the integral of an analytic function around a simple closed contour is 0", true, "Cauchy's theorem"},
		{"Goursat's proof does not assume the derivative is continuous", true, "that is the improvement over Cauchy"},
		{"the theorem applies on a simply connected domain", true, "no holes inside the contour"},
		{"an integral is unchanged by deformation avoiding singularities", true, "deformation invariance"},
		{"the integral of 1/z around |z|=1 is 2 pi i", true, "the singularity is inside"},
		{"the integral of z^2 around any closed contour is 0", true, "z^2 has an antiderivative"},
		{"a function analytic in a disc has a primitive there", true, "the disc is simply connected"},
		{"the Cauchy integral formula recovers f from its boundary values", true, "for analytic f inside the contour"},
		{"the theorem needs the derivative to be continuous", false, "Goursat removed that hypothesis"},
		{"the theorem applies through a singularity", false, "the integrand must be analytic inside"},
		{"the integral of 1/z around |z|=1 is 0", false, "it is 2 pi i"},
		{"the integral of z^2 around |z|=1 is 2 pi i", false, "it is 0"},
		{"the theorem holds only for circles", false, "any simple closed contour works"},
		{"the theorem holds only for polynomials", false, "all analytic functions qualify"},
		{"every closed contour integral is 0", false, "singularities can contribute"},
		{"analyticity outside the contour is enough", false, "analyticity inside is required"},
	}
	return countComplexVerdicts("How many of these statements about the Cauchy-Goursat theorem are true?", pool)
}

func powerSeriesProblem() generator.Problem {
	pool := []complexVerdict{
		{"the radius follows from the ratio or root test", true, "the coefficients determine it"},
		{"1/(1-z) has radius 1", true, "the singularity at 1 is on the unit circle"},
		{"e^z has infinite radius", true, "its series converges everywhere"},
		{"the coefficient is f^(n)(0)/n!", true, "the Taylor coefficient formula"},
		{"a power series expansion is unique", true, "coefficients are determined"},
		{"convergence is uniform on compact subsets inside the radius", true, "the standard theorem"},
		{"the geometric series is the model example", true, "sum z^n = 1/(1-z)"},
		{"a power series may be differentiated term by term inside its radius", true, "analyticity of the sum"},
		{"the radius is always 1", false, "e^z has infinite radius"},
		{"e^z has radius 0", false, "it is entire"},
		{"the coefficient is f(n)", false, "derivatives, not values, give coefficients"},
		{"the series converges outside its radius", false, "it diverges there in general"},
		{"a power series is unique only at the centre", false, "uniqueness holds on the disc"},
		{"1/(1-z) is entire", false, "it has a pole at 1"},
		{"term-by-term integration is invalid", false, "it is valid inside the radius"},
		{"every power series is entire", false, "the radius may be finite"},
	}
	return countComplexVerdicts("How many of these statements about power series are true?", pool)
}

func mobiusProblem() generator.Problem {
	pool := []complexVerdict{
		{"a Mobius map sends generalized circles to generalized circles", true, "circles and lines are preserved"},
		{"a Mobius map is conformal", true, "its derivative is nonzero"},
		{"a Mobius map is invertible", true, "the inverse is also Mobius"},
		{"a composition of Mobius maps is Mobius", true, "they form a group"},
		{"a Mobius map is determined by three points", true, "three images fix it"},
		{"the cross ratio is invariant under Mobius maps", true, "the defining property"},
		{"a Mobius map extends to the Riemann sphere", true, "it is a homeomorphism of the sphere"},
		{"fixed points satisfy a quadratic equation", true, "unless the map is the identity"},
		{"a Mobius map sends every circle to a circle", false, "a circle through the pole can become a line"},
		{"a Mobius map is not invertible", false, "its determinant is nonzero"},
		{"a Mobius map changes angles arbitrarily", false, "it is conformal"},
		{"a Mobius map is determined by two points", false, "three points are needed"},
		{"a Mobius map preserves Euclidean distance", false, "it distorts distance"},
		{"every Mobius map has two distinct fixed points", false, "parabolic maps have one"},
		{"every Mobius map is a rotation", false, "translations and inversions are included"},
		{"a Mobius map always sends infinity to infinity", false, "only when the denominator degree matches"},
	}
	return countComplexVerdicts("How many of these statements about Mobius maps are true?", pool)
}

func schwarzLemmaProblem() generator.Problem {
	pool := []complexVerdict{
		{"|f(z)| is at most |z| on the unit disc", true, "the Schwarz lemma bound"},
		{"if equality holds at a nonzero point then f is a rotation", true, "the equality case"},
		{"|f'(0)| is at most 1", true, "differentiate the bound at 0"},
		{"it applies to maps of the disc to itself with f(0)=0", true, "the hypotheses"},
		{"the Schwarz-Pick lemma generalises it", true, "to arbitrary disc points"},
		{"disc automorphisms are Blaschke factors", true, "the classification"},
		{"if f'(0)=1 then f is the identity", true, "the rigidity case"},
		{"the bound is strict unless f is a rotation", true, "equality forces a rotation"},
		{"|f(z)| is at least |z|", false, "the inequality goes the other way"},
		{"equality implies f is constant", false, "equality gives a rotation"},
		{"|f'(0)| is at least 1", false, "it is at most 1"},
		{"it maps the disc onto the whole plane", false, "the image stays in the disc"},
		{"every holomorphic disc map is a rotation", false, "only the equality case is"},
		{"the lemma needs f(0)=1", false, "the hypothesis is f(0)=0"},
		{"the bound is equality everywhere", false, "strict inequality is typical"},
		{"the lemma applies to maps of the plane to the disc", false, "it is a disc-to-disc statement"},
	}
	return countComplexVerdicts("How many of these statements about the Schwarz lemma are true?", pool)
}

func lociProblem() generator.Problem {
	pool := []complexVerdict{
		{"|z-a| = r is a circle", true, "points at distance r from a"},
		{"|z-a| < r is an open disk", true, "the strict inequality gives openness"},
		{"|z-a| = |z-b| is a perpendicular bisector", true, "equal distances to two points"},
		{"Re z = c is a vertical line", true, "a fixed real part"},
		{"Im z = c is a horizontal line", true, "a fixed imaginary part"},
		{"|z-a| <= r is a closed disk", true, "the boundary is included"},
		{"arg z = theta is a ray", true, "a fixed direction from the origin"},
		{"|z-a| = 2|z-b| is an Apollonius circle", true, "a fixed distance ratio"},
		{"|z-a| = r is a line", false, "it is a circle"},
		{"|z-a| = |z-b| is a circle", false, "it is a line"},
		{"Re z = c is a circle", false, "it is a line"},
		{"|z-a| < r is a circle", false, "it is a disk"},
		{"Im z = c is a vertical line", false, "it is horizontal"},
		{"arg z = theta is a line", false, "it is a ray"},
		{"|z-a| = r is a disk", false, "it is the boundary circle"},
		{"|z-a| = 2|z-b| is a line", false, "it is a circle"},
	}
	return countComplexVerdicts("How many of these descriptions of loci are correct?", pool)
}

func conjugateRootsProblem() generator.Problem {
	pool := []complexVerdict{
		{"a real polynomial has roots in conjugate pairs", true, "conjugating the equation preserves it"},
		{"an odd degree real polynomial has a real root", true, "nonreal roots pair off"},
		{"conjugate roots have equal multiplicity", true, "the conjugate factor has the same power"},
		{"the coefficients of a real polynomial are real", true, "that is the hypothesis"},
		{"x^2 + 1 has roots i and -i", true, "the two square roots of -1"},
		{"a root plus its conjugate has a real sum", true, "the imaginary parts cancel"},
		{"a root times its conjugate is real", true, "it is the squared modulus"},
		{"if 3+i is a root then 3-i is a root", true, "the conjugate pair"},
		{"every complex polynomial has conjugate root pairs", false, "complex coefficients need not be real"},
		{"if 3+i is a root then -3+i is a root", false, "the conjugate is 3-i"},
		{"a real polynomial can have exactly one nonreal root", false, "nonreal roots come in pairs"},
		{"conjugate roots can have different multiplicities", false, "the multiplicities match"},
		{"x^2 + 1 has real roots", false, "its roots are imaginary"},
		{"a complex coefficient polynomial always has conjugate roots", false, "only real coefficients force that"},
		{"the conjugate of a root is never a root", false, "it is for real polynomials"},
		{"a real polynomial of even degree must have no real roots", false, "x^2 - 1 has two"},
	}
	return countComplexVerdicts("How many of these statements about conjugate roots are true?", pool)
}

func trigExponentialProblem() generator.Problem {
	pool := []complexVerdict{
		{"cos theta = (e^(i theta) + e^(-i theta))/2", true, "the cosine identity"},
		{"sin theta = (e^(i theta) - e^(-i theta))/(2i)", true, "the sine identity"},
		{"e^(i theta) = cos theta + i sin theta", true, "Euler's formula"},
		{"cosh(i theta) = cos theta", true, "hyperbolic cosine in terms of cosine"},
		{"sinh(i theta) = i sin theta", true, "hyperbolic sine in terms of sine"},
		{"|e^(i theta)| = 1", true, "it lies on the unit circle"},
		{"de Moivre follows from the exponential form", true, "(e^(i theta))^n = e^(i n theta)"},
		{"e^(i theta) e^(i phi) = e^(i(theta+phi))", true, "the exponent law"},
		{"cos theta = (e^(i theta) - e^(-i theta))/(2i)", false, "that is sine"},
		{"sin theta = (e^(i theta) + e^(-i theta))/2", false, "that is cosine"},
		{"cosh(i theta) = i cos theta", false, "the extra i is wrong"},
		{"sinh(i theta) = sin theta", false, "the factor i is required"},
		{"e^(i theta) = cos theta - i sin theta", false, "the sign of the sine term is wrong"},
		{"|e^(i theta)| = theta", false, "the modulus is always 1"},
		{"e^(i theta) is always real", false, "it is real only at multiples of pi"},
		{"cos theta = e^(i theta) + e^(-i theta)", false, "the factor 1/2 is missing"},
	}
	return countComplexVerdicts("How many of these identities are correct?", pool)
}

func hyperbolicProblem() generator.Problem {
	pool := []complexVerdict{
		{"cosh^2 x - sinh^2 x = 1", true, "the hyperbolic Pythagorean identity"},
		{"sinh(x+y) = sinh x cosh y + cosh x sinh y", true, "the addition formula"},
		{"cosh(x+y) = cosh x cosh y + sinh x sinh y", true, "the addition formula"},
		{"d/dx sinh x = cosh x", true, "the basic derivative"},
		{"d/dx cosh x = sinh x", true, "the basic derivative"},
		{"cosh is even", true, "cosh(-x) = cosh x"},
		{"sinh is odd", true, "sinh(-x) = -sinh x"},
		{"cosh(i theta) = cos theta", true, "the circular connection"},
		{"cosh^2 x + sinh^2 x = 1", false, "the sign should be minus"},
		{"d/dx sinh x = -cosh x", false, "the derivative is positive"},
		{"cosh is odd", false, "it is even"},
		{"sinh is even", false, "it is odd"},
		{"sinh(i theta) = sin theta", false, "the factor i is missing"},
		{"cosh^2 x - sinh^2 x = 0", false, "the difference is 1"},
		{"d/dx tanh x = sech x", false, "the derivative is sech^2 x"},
		{"cosh x is always positive and never 1", false, "cosh 0 = 1"},
	}
	return countComplexVerdicts("How many of these statements about hyperbolic functions are true?", pool)
}

func openSetsProblem() generator.Problem {
	pool := []complexVerdict{
		{"{z : Re z > a}", true, "an open half-plane"},
		{"{z : |z-a| < r}", true, "an open disk"},
		{"C minus {0}", true, "a punctured plane"},
		{"the upper half-plane", true, "an open half-plane"},
		{"a union of open disks", true, "unions of open sets are open"},
		{"the complement of a finite set", true, "finite sets are closed"},
		{"C", true, "the whole plane is open"},
		{"the empty set", true, "vacuously open"},
		{"{z : |z-a| <= r}", false, "it includes its boundary"},
		{"a single point {a}", false, "no neighbourhood fits"},
		{"the real axis", false, "it has no planar neighbourhood"},
		{"the closed unit disk", false, "the boundary is included"},
		{"the integers Z", false, "it is discrete and closed"},
		{"the interval [0,1] on the real axis", false, "it is not open in C"},
		{"{z : Re z >= a}", false, "the boundary line is included"},
		{"the unit circle |z|=1", false, "it is a closed curve"},
	}
	return countComplexVerdicts("How many of these subsets of C are open?", pool)
}

func domainProblem() generator.Problem {
	pool := []complexVerdict{
		{"the open unit disk", true, "open and connected"},
		{"the upper half-plane", true, "open and connected"},
		{"the punctured plane C minus {0}", true, "open and connected"},
		{"C", true, "the whole plane"},
		{"{z : |z| > 1}", true, "the exterior of the closed disk"},
		{"the annulus 1 < |z| < 2", true, "open and connected"},
		{"{z : Re z > 0}", true, "an open half-plane"},
		{"C minus the negative real axis", true, "a slit plane"},
		{"the closed unit disk", false, "it is not open"},
		{"the real axis", false, "it has empty interior in C"},
		{"two disjoint open disks", false, "it is disconnected"},
		{"{z : |z| >= 1}", false, "it is closed and not open"},
		{"the unit circle |z| = 1", false, "it has empty interior"},
		{"{z : Re z > 1} union {z : Re z < -1}", false, "it is disconnected"},
		{"a single point", false, "it is not open"},
		{"the integers Z", false, "it is not open"},
	}
	return countComplexVerdicts("How many of these sets are domains, meaning open connected subsets of C?", pool)
}
