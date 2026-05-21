> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Exact_differential_equation) — CC BY-SA 4.0

# Exact equations

In mathematics, an **exact differential equation** or **total differential equation** is a certain kind of ordinary differential equation which is widely used in physics and engineering.

## Definition
Given a simply connected and open subset *D* of \(\mathbb{R}^2\) and two functions *I* and *J* which are continuous on *D*, an implicit first-order ordinary differential equation of the form

\(I(x, y)\, dx + J(x, y)\, dy = 0,\)

is called an **exact differential equation** if there exists a continuously differentiable function *F*, called the **potential function**, so that
\(\frac{\partial F}{\partial x} = I\)
and
\(\frac{\partial F}{\partial y} = J.\)

An exact equation may also be presented in the following form:
\(I(x, y) + J(x, y) \, y'(x) = 0\)
where the same constraints on *I* and *J* apply for the differential equation to be exact.

The nomenclature of "exact differential equation" refers to the exact differential of a function. For a function \(F(x_0, x_1,...,x_{n-1},x_n)\), the exact or total derivative with respect to \(x_0\) is given by
\(\frac{dF}{dx_0}=\frac{\partial F}{\partial x_0}+\sum_{i=1}^{n}\frac{\partial F}{\partial x_i}\frac{dx_i}{dx_0}.\)

### Example
The function \(F:\mathbb{R}^{2}\to\mathbb{R}\) given by

\(F(x,y) = \frac{1}{2}(x^2 + y^2)+c\)

is a potential function for the differential equation

\(x\,dx + y\,dy = 0.\,\)

## First-order exact differential equations
### Identifying first-order exact differential equations
Let the functions \(M\), \(N\), \(M_y\), and \(N_x\), where the subscripts denote the partial derivative with respect to the relative variable, be continuous in the region \(R: \alpha < x < \beta, \gamma < y < \delta\). Then the differential equation

\(M(x, y) + N(x, y)\frac{dy}{dx} = 0\)

is exact if and only if

\(M_y(x, y) = N_x(x, y)\)

That is, there exists a function \(\psi(x, y)\), called a * potential function*, such that

\(\psi _x(x, y) = M(x, y) \text{ and } \psi_y(x, y) = N(x, y)\)

So, in general:

\(M_y(x, y) = N_x(x, y) \iff
\begin{cases}
\exists \psi(x, y)\\
\psi_x(x, y) = M(x, y)\\
\psi_y(x, y) = N(x, y)
\end{cases}\)

#### Proof
The proof has two parts.

First, suppose there is a function \(\psi(x,y)\) such that \(\psi_x(x, y) = M(x, y) \text{ and } \psi_y(x, y) = N(x, y)\)

It then follows that \(M_y(x, y) = \psi _{xy}(x, y) \text{ and } N_x(x, y) = \psi _{yx}(x, y)\)

Since \(M_y\) and \(N_x\) are continuous, then \(\psi _{xy}\) and \(\psi _{yx}\) are also continuous which guarantees their equality.

The second part of the proof involves the construction of \(\psi(x, y)\) and can also be used as a procedure for solving first-order exact differential equations. Suppose that \(M_y(x, y) = N_x(x, y)\) and let there be a function \(\psi(x, y)\) for which
\(\psi _x(x, y) = M(x, y) \text{ and } \psi _y(x, y) = N(x, y)\)

Begin by integrating the first equation with respect to \(x\). In practice, it doesn't matter if you integrate the first or the second equation, so long as the integration is done with respect to the appropriate variable.

\(\frac{\partial \psi}{\partial x}(x, y) = M(x, y)\)
\(\psi(x, y) = \int M(x, y) \, dx + h(y)\)
\(\psi(x, y) = Q(x, y) + h(y)\)

where \(Q(x, y)\) is any differentiable function such that \(Q_x = M\). The function \(h(y)\) plays the role of a constant of integration, but instead of just a constant, it is a function of \(y\); since \(M\) is a function of both \(x\) and \(y\) and we are only integrating with respect to \(x\).

Now to show that it is always possible to find an \(h(y)\) such that \(\psi _y = N\).
\(\psi(x, y) = Q(x, y) + h(y)\)

Differentiate both sides with respect to \(y\).
\(\frac{\partial \psi}{\partial y}(x, y) = \frac{\partial Q}{\partial y}(x, y) + h'(y)\)

Set the result equal to \(N\) and solve for \(h'(y)\).
\(h'(y) = N(x, y) - \frac{\partial Q}{\partial y}(x, y)\)

In order to determine \(h'(y)\) from this equation, the right-hand side must depend only on \(y\). This can be proven by showing that its derivative with respect to \(x\) is always zero, so differentiate the right-hand side with respect to \(x\).
\(\frac{\partial N}{\partial x}(x, y) - \frac{\partial}{\partial x}\frac{\partial Q}{\partial y}(x, y) \iff \frac{\partial N}{\partial x}(x, y) - \frac{\partial}{\partial y}\frac{\partial Q}{\partial x}(x, y)\)

Since \(Q_x = M\),
\(\frac{\partial N}{\partial x}(x, y) - \frac{\partial M}{\partial y}(x, y)\)
Now, this is zero based on our initial supposition that \(M_y(x, y) = N_x(x, y)\)

Therefore,
\(h'(y) = N(x, y) - \frac{\partial Q}{\partial y}(x, y)\)
\(h(y) = \int{\left(N(x, y) - \frac{\partial Q}{\partial y}(x, y)\right) dy}\)

\(\psi(x, y) = Q(x, y) + \int \left(N(x, y) - \frac{\partial Q}{\partial y}(x, y)\right) \, dy + C\)

And this completes the proof.

### Solutions to first-order exact differential equations
First-order exact differential equations of the form
\(M(x, y) + N(x, y)\frac{dy}{dx} = 0\)

can be written in terms of the potential function \(\psi(x, y)\)
\(\frac{\partial \psi}{\partial x} + \frac{\partial \psi}{\partial y}\frac{dy}{dx} = 0\)

where
\(\begin{cases}
\psi _x(x, y) = M(x, y)\\
\psi _y(x, y) = N(x, y)
\end{cases}\)

This is equivalent to taking the total derivative of \(\psi(x,y)\).
\(\frac{\partial \psi}{\partial x} + \frac{\partial \psi}{\partial y}\frac{dy}{dx} = 0 \iff \frac{d}{dx}\psi(x, y(x)) = 0\)

The solutions to an exact differential equation are then given by
\(\psi(x, y(x)) = c\)

and the problem reduces to finding \(\psi(x, y)\).

This can be done by integrating the two expressions \(M(x, y) \, dx\) and \(N(x, y) \, dy\) and then writing down each term in the resulting expressions only once and summing them up in order to get \(\psi(x, y)\).

The reasoning behind this is the following. Since
\(\begin{cases}
\psi _x(x, y) = M(x, y)\\
\psi _y(x, y) = N(x, y)
\end{cases}\)

it follows, by integrating both sides, that
\(\begin{cases}
\psi(x, y) = \int M(x, y) \, dx + h(y) = Q(x, y) + h(y)\\
\psi(x, y) = \int N(x, y) \, dy + g(x) = P(x, y) + g(x)
\end{cases}\)

Therefore,
\(Q(x, y) + h(y) = P(x, y) + g(x)\)

where \(Q(x, y)\) and \(P(x, y)\) are differentiable functions such that \(Q_x = M\) and \(P_y = N\).

In order for this to be true and for both sides to result in the exact same expression, namely \(\psi(x, y)\), then \(h(y)\) *must* be contained within the expression for \(P(x, y)\) because it cannot be contained within \(g(x)\); since it is entirely a function of \(y\) and not \(x\) and is therefore not allowed to have anything to do with \(x\). By analogy, \(g(x)\) *must* be contained within the expression \(Q(x, y)\).

Ergo,
\(Q(x, y) = g(x) + f(x, y) \text{ and } P(x, y) = h(y) + d(x, y)\)

for some expressions \(f(x, y)\) and \(d(x, y)\).
Plugging in into the above equation, we find that
\(g(x) + f(x, y) + h(y) = h(y) + d(x, y) + g(x) \Rightarrow f(x, y) = d(x, y)\)
and so \(f(x, y)\) and \(d(x, y)\) turn out to be the same function. Therefore,
\(Q(x, y) = g(x) + f(x, y) \text { and } P(x, y) = h(y) + f(x, y)\)

Since we already showed that
\(\begin{cases}
\psi(x, y) = Q(x, y) + h(y)\\
\psi(x, y) = P(x, y) + g(x)
\end{cases}\)

it follows that
\(\psi(x, y) = g(x) + f(x, y) + h(y)\)

So, we can construct \(\psi(x, y)\) by doing \(\int M(x,y) \, dx\) and \(\int N(x, y) \, dy\) and then taking the common terms we find within the two resulting expressions (that would be \(f(x, y)\) ) and then adding the terms which are uniquely found in either one of them – \(g(x)\) and \(h(y)\).

## Second-order exact differential equations
The concept of exact differential equations can be extended to second-order equations. Consider starting with the first-order exact equation:

\(I(x,y)+J(x,y){dy \over dx}=0\)

Since both functions \(I(x,y)\), \(J(x,y)\) are functions of two variables, implicitly differentiating the multivariate function yields

\({dI \over dx} +\left({ dJ\over dx}\right){dy \over dx}+{d^2y \over dx^2} (J(x,y))=0\)

Expanding the total derivatives gives that

\({dI \over dx}={\partial I\over\partial x}+{\partial I\over\partial y}{dy \over dx}\)

and that

\({dJ \over dx}={\partial J\over\partial x}+{\partial J\over\partial y}{dy \over dx}\)

Combining the \({dy \over dx}\) terms gives

\({\partial I\over\partial x}+{dy \over dx}\left({\partial I\over\partial y}+{\partial J\over\partial x}+{\partial J\over\partial y}{dy \over dx}\right)+{d^2y \over dx^2} (J(x,y))=0\)

If the equation is exact, then...
