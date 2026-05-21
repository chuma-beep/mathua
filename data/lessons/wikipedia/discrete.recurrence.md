> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Recurrence_relation) — CC BY-SA 4.0

# Recurrence relations

In mathematics, a **recurrence relation** is an equation according to which the \(n\)th term of a sequence of numbers is equal to some combination of the previous terms. Often, only \(k\) previous terms of the sequence appear in the equation, for a parameter \(k\) that is independent of \(n\); this number \(k\) is called the *order* of the relation. If the values of the first \(k\) numbers in the sequence have been given, the rest of the sequence can be calculated by repeatedly applying the equation.

In *linear recurrences*, the th term is equated to a linear function of the \(k\) previous terms. A famous example is the recurrence for the Fibonacci numbers,
\(F_n=F_{n-1}+F_{n-2}\)
where the order \(k\) is two and the linear function merely adds the two previous terms. This example is a linear recurrence with constant coefficients, because the coefficients of the linear function (1 and 1) are constants that do not depend on \(n.\) For these recurrences, one can express the general term of the sequence as a closed-form expression of \(n\). As well, linear recurrences with polynomial coefficients depending on \(n\) are also important, because many common elementary functions and special functions have a Taylor series whose coefficients satisfy such a recurrence relation (see holonomic function).

Solving a recurrence relation means obtaining a closed-form solution: a non-recursive function of \(n\).

The concept of a recurrence relation can be extended to multidimensional arrays, that is, indexed families that are indexed by tuples of natural numbers.

## Definition
A *recurrence relation* is an equation that expresses each element of a sequence as a function of the preceding ones. More precisely, in the case where only the immediately preceding element is involved, a recurrence relation has the form
\(u_n=\varphi(n, u_{n-1})\quad\text{for}\quad n>0,\)
where
\(\varphi:\mathbb N\times X \to X\)
is a function, where is a set to which the elements of a sequence must belong. For any \(u_0\in X\), this defines a unique sequence with \(u_0\) as its first element, called the *initial value*.

It is easy to modify the definition for getting sequences starting from the term of index 1 or higher.

This defines recurrence relation of *first order*. A recurrence relation of *order* has the form
\(u_n=\varphi(n, u_{n-1}, u_{n-2}, \ldots, u_{n-k})\quad\text{for}\quad n\ge k,\)

where \(\varphi: \mathbb N\times X^k \to X\) is a function that involves consecutive elements of the sequence.
In this case, initial values are needed for defining a sequence.

## Examples
### Factorial
The factorial is defined by the recurrence relation
\(n!=n\cdot (n-1)!\quad\text{for}\quad n>0,\)
and the initial condition
\(0!=1.\)
This is an example of a *linear recurrence with polynomial coefficients* of order 1, with the simple polynomial (in )
\(n\)
as its only coefficient.

### Logistic map
An example of a recurrence relation is the logistic map defined by

\(x_{n+1} = r x_n (1 - x_n),\)

for a given constant \(r.\) The behavior of the sequence depends dramatically on \(r,\) but is stable when the initial condition \(x_0\) varies.

### Fibonacci numbers
The recurrence of order two satisfied by the Fibonacci numbers is the canonical example of a homogeneous linear recurrence relation with constant coefficients (see below). The Fibonacci sequence is defined using the recurrence

\(F_n = F_{n-1}+F_{n-2}\)

with initial conditions

\(F_0 = 0\)
\(F_1 = 1.\)

Explicitly, the recurrence yields the equations
\(F_2 = F_1 + F_0\)
\(F_3 = F_2 + F_1\)
\(F_4 = F_3 + F_2\)
etc.

We obtain the sequence of Fibonacci numbers, which begins
0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89,... The recurrence can be solved by methods described below yielding Binet's formula, which involves powers of the two roots of the characteristic polynomial \(t^2 = t + 1\); the generating function of the sequence is the rational function
\(\frac{t}{1-t-t^2}.\)

### Binomial coefficients
A simple example of a multidimensional recurrence relation is given by the binomial coefficients \(\tbinom{n}{k}\), which count the ways of selecting \(k\) elements out of a set of \(n\) elements.
They can be computed by the recurrence relation
\(\binom{n}{k}=\binom{n-1}{k-1}+\binom{n-1}{k},\)
with the base cases \(\tbinom{n}{0}=\tbinom{n}{n}=1\). Using this formula to compute the values of all binomial coefficients generates an infinite array called Pascal's triangle. The same values can also be computed directly by a different formula that is not a recurrence, but uses factorials, multiplication and division, not just additions:
\(\binom{n}{k}=\frac{n!}{k!(n-k)!}.\)

The binomial coefficients can also be computed with a uni-dimensional recurrence:
\(\binom n k = \binom n{k-1}(n-k+1)/k,\)
with the initial value \(\binom n 0 =1\) (The division is not displayed as a fraction for emphasizing that it must be computed after the multiplication, for not introducing fractional numbers).
This recurrence is widely used in computers because it does not require to build a table as does the bi-dimensional recurrence, and does not involve very large integers as does the formula with factorials (if one uses \(\binom nk= \binom n{n-k},\) all involved integers are smaller than the final result).

## Difference operator and difference equations
The **** is an operator that maps sequences to sequences, and, more generally, functions to functions. It is commonly denoted \(\Delta,\) and is defined, in functional notation, as
\((\Delta f)(x)=f(x+1)-f(x).\)
It is thus a special case of finite difference.

When using the index notation for sequences, the definition becomes
\((\Delta a)_n= a_{n+1} - a_n.\)
The parentheses around \(\Delta f\) and \(\Delta a\) are generally omitted, and \(\Delta a_n\) must be understood as the term of index in the sequence \(\Delta a,\) and not \(\Delta\) applied to the element \(a_n.\)

Given sequence \(a=(a_n)_{n\in \N},\) the **** of is \(\Delta a.\)

The **** is
\(\Delta^2 a=(\Delta\circ\Delta)a= \Delta(\Delta a).\) A simple computation shows that
\(\Delta^2 a_n= a_{n+2} - 2a_{n+1} + a_n.\)

More generally: the *th difference* is defined recursively as \(\Delta^k=\Delta\circ \Delta^{k-1},\) and one has
\(\Delta^k a_n = \sum_{t=0}^k (-1)^t \binom{k}{t} a_{n+k-t}.\)

This relation can be inverted, giving
\(a_{n+k} = a_n + {k\choose 1} \Delta a_n + \cdots + {k\choose k} \Delta^k(a_n).\)

A **** of order is an equation that involves the first differences of a sequence or a function, in the same way as a differential equation of order relates the first derivatives of a function.

The two above relations allow transforming a recurrence relation of order into a difference equation of order , and, conversely, a difference equation of order into recurrence relation of order. Each transformation is the inverse of the other, and the sequences that are solution of the difference equation are exactly those that satisfies the recurrence relation.

For example, the difference equation
\(3\Delta^2 a_n + 2\Delta a_n + 7a_n = 0\)
is equivalent to the recurrence relation
\(3a_{n+2} = 4a_{n+1} - 8a_n,\)
in the sense that the two equations are satisfied by the same sequences.

As it is equivalent for a sequence to satisfy a recurrence relation or to be the solution of a difference equation, the use of the term "difference equation" is not limited to equations using a difference operator, and the two terms "recurrence relation" and "difference equation" can be used interchangeably. See Rational difference equation, Linear constant-coefficient difference equation and Matrix difference equation for examples of using "difference equation" instead of "recurrence relation".

Difference equations resemble differential equations, and this resemblance is often used to mimic methods for solving differentiable equations to apply to solving difference equations, and therefore recurrence relations.

Summation equations relate to difference equations as integral equations relate to differential equations. See time scale calculus for a unification of the theory of difference equations with that of differential equations.

### From sequences to grids
Single-variable or one-dimensional recurrence relations are about sequences (i.e. functions defined on one-dimensional grids). Multi-variable or n-dimensional recurrence relations are about \(n\)-dimensional grids. Functions defined on \(n\)-grids can also be studied with partial difference equations.

## Solving
### Solving linear recurrence relations with constant coefficients

### Solving first-order non-homogeneous recurrence relations with variable coefficients
Moreover, for the general first-order non-homogeneous linear recurrence relation with variable coefficients:

\(a_{n+1} = f_n a_n + g_n, \qquad f_n \neq 0,\)
there is also a nice method to solve it:
\(a_{n+1} - f_n a_n = g_n\)
\(\frac{a_{n+1{\prod_{k=0}^n f_k} - \frac{f_n a_n}{\prod_{k=0}^n f_k} = \frac{g_n}{\prod_{k=0}^n f_k}\)

\(\frac{a_{n+1{\prod_{k=0}^n f_k} - \frac{a_n}{\prod_{k=0}^{n-1} f_k} = \frac{g_n}{\prod_{k=0}^n f_k}\)

Let
\(A_n = \frac{a_n}{\prod_{k=0}^{n-1} f_k},\)
Then
\(A_{n+1} - A_n = \frac{g_n}{\prod_{k=0}^n f_k}\)
\(\sum_{m=0}^{n-1}(A_{m+1} - A_m) = A_n - A_0 = \sum_{m=0}^{n-1}\frac{g_m}{\prod_{k=0}^m f_k}\)
\(\frac{a_n}{\prod_{k=0}^{n-1} f_k} = A_0 + \sum_{m=0}^{n-1}\frac{g_m}{\prod_{k=0}^m f_k}\)
\(a_n = \left(\prod_{k=0}^{n-1} f_k \right) \left(A_0 + \sum_{m=0}^{n-1}\frac{g_m}{\prod_{k=0}^m f_k}\right)\)

If we apply the formula to \(a_{n+1} = (1 + h f_{nh}) a_n + hg_{nh}\) and take the limit \(h \to 0\), we get the formula for first order linear differential equations with variable coefficients; the sum becomes an integral, and the product becomes the exponential function of an integral.

### Solving general homogeneous linear recurrence relations
Many homogeneous linear recurrence relations may be solved by means of the generalized hypergeometric series. Special cases of these lead to recurrence relations for the orthogonal polynomials, and many special functions. For example, the solution to

\(J_{n+1}=\frac{2n}{z}J_n-J_{n-1}\)

is given by

\(J_n=J_n(z),\)

the Bessel function, while

\((b-n)M_{n-1} +(2n-b+z)M_n - nM_{n+1}=0\)

is solved by

\(M_n=M(n,b;z)\)

the confluent hypergeometric series. Sequences which are the solutions of linear difference equations with polynomial coefficients are called P-recursive. For these specific recurrence equations algorithms are known which find polynomial, rational or hypergeometric solutions.

### Solving general non-homogeneous linear recurrence relations with constant coefficients
Furthermore, for the general non-homogeneous linear recurrence relation with constant coefficients, one can solve it based on variation of parameter.
### Solving first-order rational difference equations

A first order rational difference equation has the form \(w_{t+1} = \tfrac{aw_t+b}{cw_t+d}\). Such an equation can be solved by writing \(w_t\) as a nonlinear transformation of another variable \(x_t\) which itself evolves linearly. Then standard methods can be used to solve the linear difference equation in \(x_t\).

## Stability
### Stability of linear higher-order recurrences
The linear recurrence of order \(d\),

\(a_n = c_1a_{n-1} + c_2a_{n-2}+\cdots+c_da_{n-d},\)

has the characteristic equation

\(\lambda^d - c_1 \lambda^{d-1} - c_2 \lambda^{d-2} - \cdots - c_d \lambda^0 =0.\)

The recurrence is stable, meaning that the iterates converge asymptotically to a fixed value, if and only if the eigenvalues (i.e., the roots of the characteristic equation), whether real or complex, are all less than unity in absolute value.

### Stability of linear first-order matrix recurrences

In the first-order matrix difference equation

\([x_t - x^*] = A[x_{t-1}-x^*]\)

with state vector \(x\) and transition matrix \(A\), \(x\) converges asymptotically to the steady state vector \(x^*\) if and only if all eigenvalues of the transition matrix \(A\) (whether real or complex) have an absolute value which is less than 1.

### Stability of nonlinear first-order recurrences
Consider the nonlinear first-order recurrence

\(x_n=f(x_{n-1}).\)

This recurrence is locally stable, meaning that it converges to a fixed point \(x^*\) from points sufficiently close to \(x^*\), if the slope of \(f\) in the neighborhood of \(x^*\) is smaller than unity in absolute value: that is,

\(| f' (x^*) | < 1.\)

A nonlinear recurrence could have multiple fixed points, in which case some fixed points may be locally stable and others locally unstable; for continuous *f* two adjacent fixed points cannot both be locally stable.

A nonlinear recurrence relation could also have a cycle of period \(k\) for \(k > 1\). Such a cycle is stable, meaning that it attracts a set of initial conditions of positive measure, if the composite function

\(g(x) := f \circ f \circ \cdots \circ f(x)\)

with \(f\) appearing \(k\) times is locally stable according to the same criterion:

\(| g' (x^*) | < 1,\)

where \(x^*\) is any point on the cycle.

In a chaotic recurrence relation, the variable \(x\) stays in a bounded region but never converges to a fixed point or an attracting cycle; any fixed points or cycles of the equation are unstable. See also logistic map, dyadic transformation, and tent map.

## Relationship to differential equations
When solving an ordinary differential equation numerically, one typically encounters a recurrence relation. For example, when solving the initial value problem

\(y'(t) = f(t,y(t)), \ \ y(t_0)=y_0,\)

with Euler's method and a step size \(h\), one calculates the values

\(y_0=y(t_0), \ \ y_1=y(t_0+h), \ \ y_2=y(t_0+2h), \ \dots\)

by the recurrence

\(\, y_{n+1} = y_n + hf(t_n,y_n), t_n = t_0 + nh\)

Systems of linear first order differential equations can be discretized exactly analytically using the methods shown in the discretization article.

## Applications
### Mathematical biology
Some of the best-known difference equations have their origins in the attempt to model population dynamics. For example, the Fibonacci numbers were once used as a model for the growth of a rabbit population.

The logistic map is used either directly to model population growth, or as a starting point for more detailed models of population dynamics. In this context, coupled difference equations are often used to model the interaction of two or more populations. For example, the Nicholson–Bailey model for a host-parasite interaction is given by

\(N_{t+1} = \lambda N_t e^{-aP_t}\)
\(P_{t+1} = N_t(1-e^{-aP_t}),\)

with \(N_t\) representing the hosts, and \(P_t\) the parasites, at time \(t\).

Integrodifference equations are a form of recurrence relation important to spatial ecology. These and other difference equations are particularly suited to modeling univoltine populations.

### Computer science
Recurrence relations are also of fundamental importance in analysis of algorithms. If an algorithm is designed so that it will break a problem into smaller subproblems (divide and conquer), its running time is described by a recurrence relation.

A simple example is the time an algorithm takes to find an element in an ordered vector with \(n\) elements, in the worst case.

A naive algorithm will search from left to right, one element at a time. The worst possible scenario is when the required element is the last, so the number of comparisons is \(n\).

A better algorithm is called binary search. However, it requires a sorted vector. It will first check if the element is at the middle of the vector. If not, then it will check if the middle element is greater or lesser than the sought element. At this point, half of the vector can be discarded, and the algorithm can be run again on the other half. The number of comparisons will be given by

\(c_1=1\)
\(c_n=1+c_{n/2}\)

the time complexity of which will be \(O(\log_2(n))\).

### Digital signal processing
In digital signal processing, recurrence relations can model feedback in a system, where outputs at one time become inputs for future time. They thus arise in infinite impulse response (IIR) digital filters.

For example, the equation for a "feedforward" IIR comb filter of delay \(T\) is:

\(y_t = (1 - \alpha) x_t + \alpha y_{t - T},\)

where \(x_t\) is the input at time \(t\), \(y_t\) is the output at time \(t\), and \(\alpha\) controls how much of the delayed signal is fed back into the output. From this we can see that

\(y_t = (1 - \alpha) x_t + \alpha ((1-\alpha) x_{t-T} + \alpha y_{t - 2T})\)
\(y_t = (1 - \alpha) x_t + (\alpha-\alpha^2) x_{t-T} + \alpha^2 y_{t - 2T}\)

etc.

### Economics

Recurrence relations, especially linear recurrence relations, are used extensively in both theoretical and empirical economics. In particular, in macroeconomics one might develop a model of various broad sectors of the economy (the financial sector, the goods sector, the labor market, etc.) in which some agents' actions depend on lagged variables. The model would then be solved for current values of key variables (interest rate, real GDP, etc.) in terms of past and current values of other variables.
