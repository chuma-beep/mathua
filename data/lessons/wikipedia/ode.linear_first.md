> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Linear_differential_equation) — CC BY-SA 4.0

# First-order linear ODEs

In mathematics, a **linear differential equation** is a differential equation that is linear in the unknown function and its derivatives, so it can be written in the form

\[
a_0(x)y + a_1(x)y' + a_2(x)y'' \cdots + a_n(x)y^{(n)} = b(x)
\]

where *a*\(_{0}\)(*x*),... , and *b*(*x*) are arbitrary differentiable functions that do not need to be linear, and *y*′,... , *y*\(^{(*n*)}\) are the successive derivatives of an unknown function of the variable. Such an equation is an ordinary differential equation (ODE). A *linear differential equation* may also be a linear partial differential equation (PDE), if the unknown function depends on several variables, and the derivatives that appear in the equation are partial derivatives.

## Types of solution

A linear differential equation or a system of linear equations such that the associated homogeneous equations have constant coefficients may be **solved by quadrature**, which means that the solutions may be expressed in terms of integrals. This is also true for a linear equation of order one, with non-constant coefficients. An equation of order two or higher with non-constant coefficients cannot, in general, be solved by quadrature. For order two, Kovacic's algorithm allows deciding whether there are solutions in terms of integrals, and computing them if any.

The solutions of homogeneous linear differential equations with polynomial coefficients are called holonomic functions. This class of functions is stable under sums, products, differentiation, integration, and contains many usual functions and special functions such as exponential function, logarithm, sine, cosine, inverse trigonometric functions, error function, Bessel functions and hypergeometric functions. Their representation by the defining differential equation and initial conditions allows making algorithmic (on these functions) most operations of calculus, such as computation of antiderivatives, limits, asymptotic expansion, and numerical evaluation to any precision, with a certified error bound.

## Basic terminology
The highest order of derivation that appears in a (linear) differential equation is the *order* of the equation. The term *b*(*x*), which does not depend on the unknown function and its derivatives, is sometimes called the *constant term* of the equation (by analogy with algebraic equations), even when this term is a non-constant function. If the constant term is the zero function, then the differential equation is said to be *homogeneous*, as it is a homogeneous polynomial in the unknown function and its derivatives. The equation obtained by replacing, in a linear differential equation, the constant term by the zero function is the **. A differential equation has *constant coefficients* if only constant functions appear as coefficients in the associated homogeneous equation.

A ** of a differential equation is a function that satisfies the equation.
The solutions of a homogeneous linear differential equation form a vector space. In the ordinary case, this vector space has a finite dimension, equal to the order of the equation. All solutions of a linear differential equation are found by adding to a particular solution any solution of the associated homogeneous equation.

## Linear differential operator

A *basic differential operator* of order is a mapping that maps any differentiable function to its th derivative, or, in the case of several variables, to one of its partial derivatives of order. It is commonly denoted

\[
\frac{d^i}{dx^i}
\]

in the case of univariate functions, and

\[
\frac{\partial^{i_1+\cdots +i_n{\partial x_1^{i_1}\cdots \partial x_n^{i_n
\]

in the case of functions of variables. The basic differential operators include the derivative of order 0, which is the identity mapping.

A **linear differential operator** (abbreviated, in this article, as *linear operator* or, simply, *operator*) is a linear combination of basic differential operators, with differentiable functions as coefficients. In the univariate case, a linear operator has thus the form

\[
a_0(x)+a_1(x)\frac{d}{dx} + \cdots +a_n(x)\frac{d^n}{dx^n},
\]

where *a*\(_{0}\)(*x*),... , *a*\(_{*n*}\)(*x*) are differentiable functions, and the nonnegative integer is the *order* of the operator (if *a*\(_{*n*}\)(*x*) is not the zero function).

Let be a linear differential operator. The application of to a function is usually denoted *Lf* or *Lf*(*X*), if one needs to specify the variable (this must not be confused with a multiplication). A linear differential operator is a linear operator; since it maps sums to sums and the product by a scalar to the product by the same scalar.

As the sum of two linear operators is a linear operator, as well as the product (on the left) of a linear operator by a differentiable function, the linear differential operators form a vector space over the real numbers or the complex numbers (depending on the nature of the functions that are considered). They form also a free module over the ring of differentiable functions.

The language of operators allows a compact writing for differentiable equations: if

\[
L = a_0(x)+a_1(x)\frac{d}{dx} + \cdots +a_n(x)\frac{d^n}{dx^n},
\]

is a linear differential operator, then the equation

\[
a_0(x)y +a_1(x)y' + a_2(x)y'' +\cdots +a_n(x)y^{(n)}=b(x)
\]

may be rewritten

\[
Ly=b(x).
\]

There may be several variants to this notation; in particular the variable of differentiation may appear explicitly or not in and the right-hand and of the equation, such as *Ly*(*x*) = *b*(*x*) or *Ly* = *b*.

The *kernel* of a linear differential operator is its kernel as a linear mapping, that is the vector space of the solutions of the (homogeneous) differential equation *Ly* = 0.

In the case of an ordinary differential operator of order , Carathéodory's existence theorem implies that, under very mild conditions, the kernel of is a vector space of dimension , and that the solutions of the equation *Ly*(*x*) = *b*(*x*) have the form

\[
S_0(x) + c_1S_1(x) + \cdots + c_n S_n(x),
\]

where *c*\(_{1}\),... , *c*\(_{*n*}\) are arbitrary numbers. Typically, the hypotheses of Carathéodory's theorem are satisfied in an interval , if the functions *b*, *a*\(_{0}\),... , *a*\(_{*n*}\) are continuous in , and there is a positive real number such that |*a*\(_{*n*}\)(*x*)| > *k* for every in. ## Homogeneous equation with constant coefficients
A homogeneous linear differential equation has *constant coefficients* if it has the form

\[
a_0y + a_1y' + a_2y'' + \cdots + a_n y^{(n)} = 0
\]

where *a*\(_{1}\),... , *a*\(_{*n*}\) are (real or complex) numbers. In other words, it has constant coefficients if it is defined by a linear operator with constant coefficients.

The study of these differential equations with constant coefficients dates back to Leonhard Euler, who introduced the exponential function \(e^{x}\), which is the unique solution of the equation \(f' = f\), such that \(f(0) = 1\). It follows that the th derivative of \(e^{cx}\) is \(c^ne^{cx}\), and this allows solving homogeneous linear differential equations rather easily.

Let

\[
a_0y + a_1y' + a_2y'' + \cdots + a_ny^{(n)} = 0
\]

be a homogeneous linear differential equation with constant coefficients (that is *a*\(_{0}\),... , *a*\(_{*n*}\) are real or complex numbers).

Searching for solutions of this equation that have the form *e*\(^{*αx*}\) is equivalent to searching the constants such that

\[
a_0e^{\alpha x} + a_1\alpha e^{\alpha x} + a_2\alpha^2 e^{\alpha x}+\cdots + a_n\alpha^n e^{\alpha x} = 0.
\]

Factoring out *e*\(^{*αx*}\) (which is never zero), shows that must be a root of the *characteristic polynomial*

\[
a_0 + a_1t + a_2 t^2 + \cdots + a_nt^n
\]

of the differential equation, which is the left-hand side of the characteristic equation

\[
a_0 + a_1t + a_2 t^2 + \cdots + a_nt^n = 0.
\]

When these roots are all distinct, one has distinct solutions that are not necessarily real, even if the coefficients of the equation are real. These solutions can be shown to be linearly independent, by considering the Vandermonde determinant of the values of these solutions at *x* = 0,... , *n* – 1. Together they form a basis of the vector space of solutions of the differential equation (that is, the kernel of the differential operator).
| | Example |
| --- |
| | \[ y**'-2y**+2y''-2y'+y=0 \] has the characteristic equation \[ z^4-2z^3+2z^2-2z+0. \] This has zeros, −*i*, and 1 (multiplicity 2). The solution basis is thus \[ e^{ix},\; e^{-ix},\; e^x,\; xe^x. \] A real basis of solution is thus \[ \cos x,\; \sin x,\; e^x,\; xe^x. \] |

In the case where the characteristic polynomial has only simple roots, the preceding provides a complete basis of the solutions vector space. In the case of multiple roots, more linearly independent solutions are needed for having a basis. These have the form

\[
x^ke^{\alpha x},
\]

where is a nonnegative integer, is a root of the characteristic polynomial of multiplicity , and *k* < *m*. For proving that these functions are solutions, one may remark that if is a root of the characteristic polynomial of multiplicity , the characteristic polynomial may be factored as *P*(*t*)(*t* − *α*)\(^{*m*}\). Thus, applying the differential operator of the equation is equivalent with applying first times the operator.
