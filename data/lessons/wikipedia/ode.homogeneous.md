> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Homogeneous_differential_equation) — CC BY-SA 4.0

# Homogeneous second-order ODEs

A differential equation can be **homogeneous** in either of two respects.

A first order differential equation is said to be homogeneous if it may be written

\[
f(x,y) \, dy = g(x,y) \, dx,
\]

where  and  are homogeneous functions of the same degree of  and . In this case, the change of variable 1=*y* = *ux* leads to an equation of the form

\[
\frac{dx}{x} = h(u) \, du,
\]

which is easy to solve by integration of the two members.

Otherwise, a differential equation is homogeneous if it is a homogeneous function of the unknown function and its derivatives. In the case of linear differential equations, this means that there are no constant terms. The solutions of any linear ordinary differential equation of any order may be deduced by integration from the solution of the homogeneous equation obtained by removing the constant term.

## History
The term *homogeneous* was first applied to differential equations by Johann Bernoulli in section 9 of his 1726 article *De integraionibus aequationum differentialium* (On the integration of differential equations).

## Homogeneous first-order differential equations


A first-order ordinary differential equation in the form:


\[
M(x,y)\,dx + N(x,y)\,dy = 0
\]


is a homogeneous type if both functions *M*(*x*, *y*) and *N*(*x*, *y*) are homogeneous functions of the same degree . That is, multiplying each variable by a parameter λ, we find


\[
M(\lambda x, \lambda y) = \lambda^n M(x,y) \quad \text{and} \quad N(\lambda x, \lambda y) = \lambda^n N(x,y)\,.
\]


Thus,

\[
\frac{M(\lambda x, \lambda y)}{N(\lambda x, \lambda y)} = \frac{M(x,y)}{N(x,y)}\,.
\]


### Solution method
In the quotient \(\frac{M(tx,ty)}{N(tx,ty)} = \frac{M(x,y)}{N(x,y)}\), we can let 1=*t* =  to simplify this quotient to a function  of the single variable :


\[
\frac{M(x,y)}{N(x,y)} = \frac{M(tx,ty)}{N(tx,ty)} = \frac{M(1,y/x)}{N(1,y/x)}=f(y/x)\,.
\]

That is

\[
\frac{dy}{dx} = -f(y/x).
\]


Introduce the change of variables 1=*y* = *ux*; differentiate using the product rule:


\[
\frac{dy}{dx}=\frac{d(ux)}{dx} = x\frac{du}{dx} + u\frac{dx}{dx} = x\frac{du}{dx} + u.
\]


This transforms the original differential equation into the separable form

\[
x\frac{du}{dx} = -f(u) - u,
\]

or

\[
\frac 1x\frac{dx}{du} = \frac {-1}{f(u) + u},
\]

which can now be integrated directly: ln *x* equals the antiderivative of the right-hand side (see ordinary differential equation).

### Special case
A first order differential equation of the form (, , , , ,  are all constants)

\[
\left(ax + by + c\right) dx + \left(ex + fy + g\right) dy = 0
\]

where *af* ≠ *be*
can be transformed into a homogeneous type by a linear transformation of both variables ( and  are constants):

\[
t = x + \alpha; \;\; z = y + \beta \,,
\]

where

\[
\alpha=\frac{cf-bg}{af-be}; \;\; \beta=\frac{ag-ce}{af-be}\,.
\]

For cases where 1=*af* = *be*, introduce the change of variables 1=*u* = *ax* + *by* or 1=*u* = *ex* + *fy*; differentiation yields:

\[
\frac{du}{dx}=a-b\bigg(\frac{ac+au}{ag+eu}\bigg),
\]

or

\[
\frac{du}{dx}=e-f\bigg(\frac{ec+au}{eg+eu}\bigg),
\]

for each respective substitution. Both may be solved via separation of variables.

## Homogeneous linear differential equations

A linear differential equation is  **homogeneous** if it is a homogeneous linear equation in the unknown function and its derivatives. It follows that, if *φ*(*x*) is a solution, so is *cφ*(*x*), for any (non-zero) constant . In order for this condition to hold, each nonzero term of the linear differential equation must depend on the unknown function or any derivative of it. A linear differential equation that fails this condition is called **inhomogeneous.**

A linear differential equation can be represented as a linear operator acting on *y*(*x*) where  is usually the independent variable and  is the dependent variable. Therefore, the general form of a linear homogeneous differential equation is


\[
L(y) = 0
\]


where  is a differential operator, a sum of derivatives (defining the "0th derivative" as the original, non-differentiated function), each multiplied by a function *f*<sub>*i*</sub> of :


\[
L = \sum_{i=0}^n f_i(x)\frac{d^i}{dx^i} \, ,
\]

where *f*<sub>*i*</sub> may be constants, but not all *f*<sub>*i*</sub> may be zero.

For example, the following linear differential equation is homogeneous:


\[
\sin(x) \frac{d^2y}{dx^2} + 4 \frac{dy}{dx} + y = 0 \,,
\]


whereas the following two are inhomogeneous:


\[
2 x^2 \frac{d^2y}{dx^2} + 4 x \frac{dy}{dx} + y = \cos(x) \,;
\]


\[
2 x^2 \frac{d^2y}{dx^2} - 3 x \frac{dy}{dx} + y = 2 \,.
\]

The existence of a constant term is a sufficient condition for an equation to be inhomogeneous, as in the above example.

## See also
* Separation of variables

## Notes


## References
* . (This is a good introductory reference on differential equations.)
* . (This is a classic reference on ODEs, first published in 1926.)
*
*

## External links
*Homogeneous differential equations at MathWorld
*Wikibooks: Ordinary Differential Equations/Substitution 1

