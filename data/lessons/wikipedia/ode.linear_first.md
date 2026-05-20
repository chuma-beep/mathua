> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Linear_differential_equation) — CC BY-SA 4.0

# First-order linear ODEs

In mathematics, a **linear differential equation** is a differential equation that is linear in the unknown function and its derivatives, so it can be written in the form

\[
a_0(x)y + a_1(x)y' + a_2(x)y'' \cdots + a_n(x)y^{(n)} = b(x)
\]

where 1=*a*<sub>0</sub>(*x*), ...,  and *b*(*x*) are arbitrary differentiable functions that do not need to be linear, and *y*′, ..., *y*<sup>(*n*)</sup>  are the successive derivatives of an unknown function  of the variable .

Such an equation is an ordinary differential equation (ODE). A *linear differential equation* may also be a linear partial differential equation (PDE), if the unknown function depends on several variables, and the derivatives that appear in the equation are partial derivatives.

## Types of solution

A linear differential equation or a system of linear equations such that the associated homogeneous equations have constant coefficients may be **solved by quadrature**, which means that the solutions may be expressed in terms of integrals. This is also true for a linear equation of order one, with non-constant coefficients. An equation of order two or higher with non-constant coefficients cannot, in general, be solved by quadrature. For order two, Kovacic's algorithm allows deciding whether there are solutions in terms of integrals, and computing them if any.

The solutions of homogeneous linear differential equations with polynomial coefficients are called holonomic functions. This class of functions is stable under sums, products, differentiation, integration, and contains many usual functions and special functions such as exponential function, logarithm, sine, cosine, inverse trigonometric functions, error function, Bessel functions and hypergeometric functions. Their representation by the defining differential equation and initial conditions allows making algorithmic (on these functions) most operations of calculus, such as computation of antiderivatives, limits, asymptotic expansion, and numerical evaluation to any precision, with a certified error bound.

## Basic terminology
The highest order of derivation that appears in a (linear) differential equation is the *order* of the equation. The term *b*(*x*), which does not depend on the unknown function and its derivatives, is sometimes called the *constant term* of the equation (by analogy with algebraic equations), even when this term is a non-constant function. If the constant term is the zero function, then the differential equation is said to be *homogeneous*, as it is a homogeneous polynomial in the unknown function and its derivatives. The equation obtained by replacing, in a linear differential equation, the constant term by the zero function is the **. A differential equation has *constant coefficients* if only constant functions appear as coefficients in the associated homogeneous equation.

A ** of a differential equation is a function that satisfies the equation.
The solutions of a homogeneous linear differential equation form a vector space. In the ordinary case, this vector space has a finite dimension, equal to the order of the equation. All solutions of a linear differential equation are found by adding to a particular solution any solution of the associated homogeneous equation.

## Linear differential operator

A *basic differential operator* of order  is a mapping that maps any differentiable function to its th derivative, or, in the case of several variables, to one of its partial derivatives of order . It is commonly denoted

\[
\frac{d^i}{dx^i}
\]

in the case of univariate functions, and

\[
\frac{\partial^{i_1+\cdots +i_n}}{\partial x_1^{i_1}\cdots \partial x_n^{i_n}}
\]

in the case of functions of  variables. The basic differential operators include the derivative of order 0, which is the identity mapping.

A **linear differential operator** (abbreviated, in this article, as *linear operator* or, simply, *operator*) is a linear combination of basic differential operators, with differentiable functions as coefficients. In the univariate case, a linear operator has thus the form

\[
a_0(x)+a_1(x)\frac{d}{dx} + \cdots +a_n(x)\frac{d^n}{dx^n},
\]

where *a*<sub>0</sub>(*x*), ..., *a*<sub>*n*</sub>(*x*) are differentiable functions, and the nonnegative integer  is the *order* of the operator (if *a*<sub>*n*</sub>(*x*) is not the zero function).

Let  be a linear differential operator. The application of  to a function  is usually denoted *Lf* or *Lf*(*X*), if one needs to specify the variable (this must not be confused with a multiplication). A linear differential operator is a linear operator, since it maps sums to sums and the product by a scalar to the product by the same scalar.

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


There may be several variants to this notation; in particular the variable of differentiation may appear explicitly or not in  and the right-hand and of the equation, such as 1=*Ly*(*x*) = *b*(*x*) or 1=*Ly* = *b*.

The *kernel* of a linear differential operator is its kernel as a linear mapping, that is the vector space of the solutions of the (homogeneous) differential equation 1=*Ly* = 0.

In the case of an ordinary differential operator of order , Carathéodory's existence theorem implies that, under very mild conditions, the kernel of  is a vector space of dimension , and that the solutions of the equation 1=*Ly*(*x*) = *b*(*x*) have the form

\[
S_0(x) + c_1S_1(x) + \cdots + c_n S_n(x),
\]

where *c*<sub>1</sub>, ..., *c*<sub>*n*</sub> are arbitrary numbers. Typically, the hypotheses of Carathéodory's theorem are satisfied in an interval , if the functions *b*, *a*<sub>0</sub>, ..., *a*<sub>*n*</sub> are continuous in , and there is a positive real number  such that 1=|*a*<sub>*n*</sub>(*x*)| > *k* for every  in .

## Homogeneous equation with constant coefficients
A homogeneous linear differential equation has *constant coefficients* if it has the form

\[
a_0y + a_1y' + a_2y'' + \cdots + a_n y^{(n)} = 0
\]

where *a*<sub>1</sub>, ..., *a*<sub>*n*</sub> are (real or complex) numbers. In other words, it has constant coefficients if it is defined by a linear operator with constant coefficients.

The study of these differential equations with constant coefficients dates back to Leonhard Euler, who introduced the exponential function \(e^{x}\), which is the unique solution of the equation \(f' = f\), such that \(f(0) = 1\). It follows that the th derivative of \(e^{cx}\) is \(c^ne^{cx}\), and this allows solving homogeneous linear differential equations rather easily.

Let

\[
a_0y + a_1y' + a_2y'' + \cdots + a_ny^{(n)} = 0
\]

be a homogeneous linear differential equation with constant coefficients (that is *a*<sub>0</sub>, ..., *a*<sub>*n*</sub> are real or complex numbers).

Searching for solutions of this equation that have the form *e*<sup>*αx*</sup> is equivalent to searching the constants  such that

\[
a_0e^{\alpha x} + a_1\alpha e^{\alpha x} + a_2\alpha^2 e^{\alpha x}+\cdots + a_n\alpha^n e^{\alpha x} = 0.
\]

Factoring out *e*<sup>*αx*</sup> (which is never zero), shows that  must be a root of the *characteristic polynomial*

\[
a_0 + a_1t + a_2 t^2 + \cdots + a_nt^n
\]

of the differential equation, which is the left-hand side of the characteristic equation

\[
a_0 + a_1t + a_2 t^2 + \cdots + a_nt^n = 0.
\]


When these roots are all distinct, one has  distinct solutions that are not necessarily real, even if the coefficients of the equation are real. These solutions can be shown to be linearly independent, by considering the Vandermonde determinant of the values of these solutions at 1=*x* = 0, ..., *n* – 1. Together they form a basis of the vector space of solutions of the differential equation (that is, the kernel of the differential operator).
{| class="toccolours floatright" style="width:35%; margin: 0.5em 0 0.5em 1em;"
! style="background:#ffffaa; padding: 3px 5px 3px 5px; font-size:larger;" | Example
|-
| style="font-size:100%; padding:0 5px 0 5px;" |

\[
y**'-2y**+2y''-2y'+y=0
\]

has the characteristic equation

\[
z^4-2z^3+2z^2-2z+1=0.
\]

This has zeros, , −*i*, and 1 (multiplicity 2). The solution basis is thus

\[
e^{ix},\; e^{-ix},\; e^x,\; xe^x.
\]

A real basis of solution is thus

\[
\cos x,\; \sin x,\; e^x,\; xe^x.
\]


|}
In the case where the characteristic polynomial has only simple roots, the preceding provides a complete basis of the solutions vector space. In the case of multiple roots, more linearly independent solutions are needed for having a basis. These have the form

\[
x^ke^{\alpha x},
\]

where  is a nonnegative integer,  is a root of the characteristic polynomial of multiplicity , and *k* < *m*. For proving that these functions are solutions, one may remark that if  is a root of the characteristic polynomial of multiplicity , the characteristic polynomial may be factored as *P*(*t*)(*t* − *α*)<sup>*m*</sup>. Thus, applying the differential operator of the equation is equivalent with applying first  times the operator {{nowrap|\(\frac{d}{dx} - \alpha\),}} and then the operator that has  as characteristic polynomial. By the exponential shift theorem,

\[
\left(\frac{d}{dx}-\alpha\right)\left(x^ke^{\alpha x}\right)= kx^{k-1}e^{\alpha x},
\]


and thus one gets zero after *k* + 1 application of {{nowrap|1=\(\frac{d}{dx} - \alpha\).}}

As, by the fundamental theorem of algebra, the sum of the multiplicities of the roots of a polynomial equals the degree of the polynomial, the number of above solutions equals the order of the differential equation, and these solutions form a basis of the vector space of the solutions.

In the common case where the coefficients of the equation are real, it is generally more convenient to have a basis of the solutions consisting of real-valued functions. Such a basis may be obtained from the preceding basis by remarking that, if *a* + *ib* is a root of the characteristic polynomial, then *a* – *ib* is also a root, of the same multiplicity. Thus a real basis is obtained by using Euler's formula, and replacing \(x^ke^{(a+ib)x}\) and \(x^ke^{(a-ib)x}\) by \(x^ke^{ax} \cos(bx)\) and \(x^ke^{ax} \sin(bx)\).

### Second-order case
A homogeneous linear differential equation of the second order may be written

\[
y'' + ay' + by = 0,
\]

and its characteristic polynomial is

\[
r^2 + ar + b.
\]


If  and  are real, there are three cases for the solutions, depending on the discriminant 1=*D* = *a*<sup>2</sup> − 4*b*. In all three cases, the general solution depends on two arbitrary constants *c*<sub>1</sub> and *c*<sub>2</sub>.

* If *D* > 0, the characteristic polynomial has two distinct real roots , and . In this case, the general solution is
\[
c_1 e^{\alpha x} + c_2 e^{\beta x}.
\]

* If 1=*D* = 0, the characteristic polynomial has a double root −*a*/2, and the general solution is
\[
(c_1 + c_2 x) e^{-ax/2}.
\]

* If *D* < 0, the characteristic polynomial has two complex conjugate roots *α* ± *βi*, and the general solution is
\[
c_1 e^{(\alpha + \beta i)x} + c_2 e^{(\alpha - \beta i)x},
\]
 which may be rewritten in real terms, using Euler's formula as
\[
e^{\alpha x} (c_1\cos(\beta x) + c_2 \sin(\beta x)).
\]


Finding the solution *y*(*x*) satisfying 1=*y*(0) = *d*<sub>1</sub> and 1=*y*′(0) = *d*<sub>2</sub>, one equates the values of the above general solution at 0 and its derivative there to *d*<sub>1</sub> and *d*<sub>2</sub>, respectively. This results in a linear system of two linear equations in the two unknowns *c*<sub>1</sub> and *c*<sub>2</sub>. Solving this system gives the solution for a so-called Cauchy problem, in which the values at 0 for the solution of the DEQ and its derivative are specified.

## Non-homogeneous equation with constant coefficients
A non-homogeneous equation of order  with constant coefficients may be written

\[
y^{(n)}(x) + a_1 y^{(n-1)}(x) + \cdots + a_{n-1} y'(x)+ a_ny(x) = f(x),
\]

where *a*<sub>1</sub>, ..., *a*<sub>*n*</sub> are real or complex numbers,  is a given function of , and  is the unknown function (for sake of simplicity, "(*x*)" will be omitted in the following).

There are several methods for solving such an equation. The best method depends on the nature of the function  that makes the equation non-homogeneous. If  is a linear combination of exponential and sinusoidal functions, then the exponential response formula may be used. If, more generally,  is a linear combination of functions of the form *x*<sup>*n*</sup>*e*<sup>*ax*</sup>, *x*<sup>*n*</sup> cos(*ax*), and *x*<sup>*n*</sup> sin(*ax*), where  is a nonnegative integer, and  a constant (which need not be the same in each term), then the method of undetermined coefficients may be used. Still more general, the annihilator method applies when  satisfies a homogeneous linear differential equation, typically, a holonomic function.

The most general method is the variation of constants, which is presented here.

The general solution of the associated homogeneous equation

\[
y^{(n)} + a_1 y^{(n-1)} + \cdots + a_{n-1} y'+ a_ny = 0
\]

is

\[
y=u_1y_1+\cdots+ u_ny_n,
\]

where (*y*<sub>1</sub>, ..., *y*<sub>*n*</sub>) is a basis of the vector space of the solutions and *u*<sub>1</sub>, ..., *u*<sub>*n*</sub> are arbitrary constants. The method of variation of constants takes its name from the following idea. Instead of considering *u*<sub>1</sub>, ..., *u*<sub>*n*</sub> as constants, they can be considered as unknown functions that have to be determined for making  a solution of the non-homogeneous equation. For this purpose, one adds the constraints

\[
\begin{align}
0 &= u'_1y_1 + u'_2y_2 + \cdots+u'_ny_n \\
0 &= u'_1y'_1 + u'_2y'_2 + \cdots + u'_n y'_n \\
  &\;\;\vdots \\
0 &= u'_1y^{(n-2)}_1+u'_2y^{(n-2)}_2 + \cdots + u'_n y^{(n-2)}_n,
\end{align}
\]

which imply (by product rule and induction)

\[
y^{(i)} = u_1 y_1^{(i)} + \cdots + u_n y_n^{(i)}
\]

for 1=*i* = 1, ..., *n* – 1, and

\[
y^{(n)} = u_1 y_1^{(n)} + \cdots + u_n y_n^{(n)} +u'_1y_1^{(n-1)}+u'_2y_2^{(n-1)}+\cdots+u'_ny_n^{(n-1)}.
\]


Replacing in the original equation  and its derivatives by these expressions, and using the fact that *y*<sub>1</sub>, ..., *y*<sub>*n*</sub> are solutions of the original homogeneous equation, one gets

\[
f=u'_1y_1^{(n-1)} + \cdots + u'_ny_n^{(n-1)}.
\]


This equation and the above ones with 0 as left-hand side form a system of  linear equations in *u*′<sub>1</sub>, ..., *u*′<sub>*n*</sub> whose coefficients are known functions (, the *y*, and their derivatives). This system can be solved by any method of linear algebra. The computation of antiderivatives gives *u*<sub>1</sub>, ..., *u*<sub>*n*</sub>, and then 1=*y* = *u*<sub>1</sub>*y*<sub>1</sub> + ⋯ + *u*<sub>*n*</sub>*y*<sub>*n*</sub>.

As antiderivatives are defined up to the addition of a constant, one finds again that the general solution of the non-homogeneous equation is the sum of an arbitrary solution and the general solution of the associated homogeneous equation.

## First-order equation with variable coefficients
The general form of a linear ordinary differential equation of order 1, after dividing out the coefficient of *y*′(*x*), is:

\[
y'(x) = f(x) y(x) + g(x).
\]


If the equation is homogeneous, i.e. 1=*g*(*x*) = 0, one may rewrite and integrate:

\[
\frac{y'}{y}= f, \qquad \log y = k +F,
\]

where  is an arbitrary constant of integration and \(F=\textstyle\int f\,dx\) is any antiderivative of . Thus, the general solution of the homogeneous equation is

\[
y=ce^F,
\]

where 1=*c* = *e*<sup>*k*</sup> is an arbitrary constant.

For the general non-homogeneous equation, it is useful to multiply both sides of the equation by the reciprocal *e*<sup>−*F*</sup> of a solution of the homogeneous equation. This gives

\[
y'e^{-F}-yfe^{-F}= ge^{-F}.
\]

As {{tmath|1=-fe^{-F} = \tfrac{d}{dx} \left(e^{-F}\right),}} the product rule allows rewriting the equation as

\[
\frac{d}{dx}\left(ye^{-F}\right)= ge^{-F}.
\]

Thus, the general solution is

\[
y=ce^F + e^F\int ge^{-F}dx,
\]

where  is a constant of integration, and  is any antiderivative of  (changing of antiderivative amounts to change the constant of integration).

### Alternate Approach
Let us multiply both sides of equation
\(y' - f(x) y = g(x).\) by a function \(h(x)\)

\[
h(x)y' - h(x)f(x) y = h(x) g(x).
\]


We want left had side to be a total derivative: \(d(h(x)y(x))/dx\).
Using product rule, \(d(h(x)y)= h(x)dy+ydh(x)\) comparing with the original equation, \(-h(x)f(x)y= yh'(x)\) or \(dh(x)/h(x)=-f(x)dx\) or integrating we get

\[
h(x)=e^{-\textstyle\int f\,dx}
\]

And integrating \(d(h(x)y(x))/dx=h(x)g(x)\)
we get

\[
y(x)=\frac{C+\textstyle\int h(x)g(x)\,dx}{h(x)}
\]


### Example
Solving the equation

\[
y'(x) + \frac{y(x)}{x} = 3x.
\]

The associated homogeneous equation \(y'(x) + \frac{y(x)}{x} = 0\) gives

\[
\frac{y'}{y}=-\frac{1}{x},
\]

that is

\[
y=\frac{c}{x}.
\]


Dividing the original equation by one of these solutions gives

\[
xy'+y=3x^2.
\]

That is

\[
(xy)'=3x^2,
\]


\[
xy=x^3 +c,
\]

and

\[
y(x)=x^2+c/x.
\]

For the initial condition

\[
y(1)=\alpha,
\]

one gets the particular solution

\[
y(x)=x^2+\frac{\alpha-1}{x}.
\]


## System of linear differential equations


A system of linear differential equations consists of several linear differential equations that involve several unknown functions. In general one restricts the study to systems such that the number of unknown functions equals the number of equations.

An arbitrary linear ordinary differential equation and a system of such equations can be converted into a first order system of linear differential equations by adding variables for all but the highest order derivatives. That is, if {{tmath| y', y*, \ldots, y^{(k)} }} appear in an equation, one may replace them by new unknown functions  that must satisfy the equations  and {{tmath|1=y_i'=y_{i+1},}} for 1=*i* = 1, ..., *k'' – 1.

A linear system of the first order, which has  unknown functions and  differential equations may normally be solved for the derivatives of the unknown functions. If it is not the case this is a differential-algebraic system, and this is a different theory. Therefore, the systems that are considered here have the form

\[
\begin{align}
y_1'(x) &= b_1(x) +a_{1,1}(x)y_1+\cdots+a_{1,n}(x)y_n\\[1ex]
&\;\;\vdots\\[1ex]
y_n'(x) &= b_n(x) +a_{n,1}(x)y_1+\cdots+a_{n,n}(x)y_n,
\end{align}
\]

where  and the {{tmath|a_{i,j} }} are functions of . In matrix notation, this system may be written (omitting "(*x*)")

\[
\mathbf{y}' = A\mathbf{y}+\mathbf{b}.
\]


The solving method is similar to that of a single first order linear differential equations, but with complications stemming from noncommutativity of matrix multiplication.

Let

\[
\mathbf{u}' = A\mathbf{u}.
\]

be the homogeneous equation associated to the above matrix equation.
Its solutions form a vector space of dimension , and are therefore the columns of a square matrix of functions , whose determinant is not the zero function. If 1=*n* = 1, or  is a matrix of constants, or, more generally, if  commutes with its antiderivative , then one may choose  equal the exponential of . In fact, in these cases, one has

\[
\frac{d}{dx}\exp(B) = A\exp (B).
\]

In the general case there is no closed-form solution for the homogeneous equation, and one has to use either a numerical method, or an approximation method such as Magnus expansion.

Knowing the matrix , the general solution of the non-homogeneous equation is

\[
\mathbf{y}(x) = U(x)\mathbf{y_0} + U(x)\int U^{-1}(x)\mathbf{b}(x)\,dx,
\]

where the column matrix \(\mathbf{y_0}\) is an arbitrary constant of integration.

If initial conditions are given as

\[
\mathbf y(x_0)=\mathbf y_0,
\]

the solution that satisfies these initial conditions is

\[
\mathbf{y}(x) = U(x)U^{-1}(x_0)\mathbf{y_0} + U(x)\int_{x_0}^x U^{-1}(t)\mathbf{b}(t)\,dt.
\]


## Higher order with variable coefficients
A linear ordinary equation of order one with variable coefficients may be solved by quadrature, which means that the solutions may be expressed in terms of integrals. This is not the case for order at least two. This is the main result of Picard–Vessiot theory which was initiated by Émile Picard and Ernest Vessiot, and whose recent developments are called differential Galois theory.

The impossibility of solving by quadrature can be compared with the Abel–Ruffini theorem, which states that an algebraic equation of degree at least five cannot, in general, be solved by radicals. This analogy extends to the proof methods and motivates the denomination of differential Galois theory.

Similarly to the algebraic case, the theory allows deciding which equations may be solved by quadrature, and if possible solving them. However, for both theories, the necessary computations are extremely difficult, even with the most powerful computers.

Nevertheless, the case of order two with rational coefficients has been completely solved by Kovacic's algorithm.

### Cauchy–Euler equation
Cauchy–Euler equations are examples of equations of any order, with variable coefficients, that can be solved explicitly. These are the equations of the form

\[
x^n y^{(n)}(x) + a_{n-1} x^{n-1} y^{(n-1)}(x) + \cdots + a_0 y(x) = 0,
\]

where {{tmath|a_0, \ldots, a_{n-1} }} are constant coefficients.

## Holonomic functions

A holonomic function, also called  a *D-finite function*, is a function that is a solution of a homogeneous linear differential equation with polynomial coefficients.

Most functions that are commonly considered in mathematics are holonomic or quotients of holonomic functions. In fact, holonomic functions include polynomials, algebraic functions, logarithm, exponential function, sine, cosine, hyperbolic sine, hyperbolic cosine, inverse trigonometric and inverse hyperbolic functions, and many special functions such as Bessel functions and hypergeometric functions.

Holonomic functions have several closure properties; in particular, sums, products, derivative and integrals of holonomic functions are holonomic. Moreover, these closure properties are effective, in the sense that there are algorithms for computing the differential equation of the result of any of these operations, knowing the differential equations of the input.

Usefulness of the concept of holonomic functions results of Zeilberger's theorem, which follows.

A *holonomic sequence* is a sequence of numbers that may be generated by a recurrence relation with polynomial coefficients. The coefficients of the Taylor series at a point of a holonomic function form a holonomic sequence. Conversely, if the sequence of the coefficients of a power series is holonomic, then the series defines a holonomic function (even if the radius of convergence is zero). There are efficient algorithms for both conversions, that is for computing the recurrence relation from the differential equation, and *vice versa*.


It follows that, if one represents (in a computer) holonomic functions by their defining differential equations and initial conditions, most calculus operations can be done automatically on these functions, such as derivative, indefinite and definite integral, fast computation of Taylor series (thanks of the recurrence relation on its coefficients), evaluation to a high precision with certified bound of the approximation error, limits, localization of singularities, asymptotic behavior at infinity and near singularities, proof of identities, etc.

## See also
* Continuous-repayment mortgage
* Fourier transform
* Laplace transform
* Linear difference equation
* Variation of parameters

## References

*
*
*

## External links
* https://eqworld.ipmnet.ru/en/solutions/ode.htm
* Dynamic Dictionary of Mathematical Function . Automatic and interactive study of many holonomic functions.

