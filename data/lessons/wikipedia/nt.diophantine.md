> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Diophantine_equation) — CC BY-SA 4.0

# Linear Diophantine equations

In mathematics, a **Diophantine equation** is a polynomial equation with integer coefficients, for which only integer solutions are of interest. A **linear Diophantine equation** equates the sum of two or more unknowns, with coefficients, to a constant. An **exponential Diophantine equation** is one in which unknowns can appear in exponents.

**Diophantine problems** have fewer equations than unknowns and involve finding integers that solve all equations simultaneously. Because such systems of equations define algebraic curves, algebraic surfaces, or, more generally, algebraic sets, their study is a part of algebraic geometry that is called *Diophantine geometry*.

The word *Diophantine* refers to the Hellenistic mathematician of the 3rd century, Diophantus of Alexandria, who made a study of such equations and was one of the first mathematicians to introduce symbolism into algebra. The mathematical study of Diophantine problems that Diophantus initiated is now called **Diophantine analysis**.

While individual equations present a kind of puzzle and have been considered throughout history, the formulation of general theories of Diophantine equations, beyond the case of linear and quadratic equations, was an achievement of the twentieth century.

## Examples
In the following Diophantine equations, and are the unknowns and the other letters are given constants:
<dl>
  <dt>\(ax+by = c\)</dt>
  <dd>This is a linear Diophantine equation, related to Bézout's identity.</dd>
  <dt>\(w^3 + x^3 = y^3 + z^3\)</dt>
  <dd>The smallest nontrivial solution in positive integers is 12\(^{3}\) + 1\(^{3}\) = 9\(^{3}\) + 10\(^{3}\) = 1729. It was famously given as an evident property of 1729, a taxicab number (also named Hardy–Ramanujan number) by Ramanujan to Hardy while meeting in 1917. There are infinitely many nontrivial solutions.</dd>
  <dt>\(x^n + y^n = z^n\)</dt>
  <dd>For *n* = 2 there are infinitely many solutions (*x, y, z*): the Pythagorean triples. For larger integer values of , Fermat's Last Theorem (initially claimed in 1637 by Fermat and proved by Andrew Wiles in 1995) states there are no positive integer solutions (*x, y, z*).</dd>
  <dt>\(x^2 - ny^2 = \pm 1\)</dt>
  <dd>This is Pell's equation, which is named after the English mathematician John Pell. It was studied by Brahmagupta in the 7th century, as well as by Fermat in the 17th century.</dd>
  <dt>\(\frac 4 n = \frac 1 x + \frac 1 y + \frac 1 z\)</dt>
  <dd>The Erdős–Straus conjecture states that, for every positive integer ≥ 2, there exists a solution in , and , all as positive integers. Although not usually stated in polynomial form, this example is equivalent to the polynomial equation \(4xyz = n(yz+xz+xy).\)</dd>
  <dt>\(x^4 + y^4 + z^4 = w^4\)</dt>
  <dd>Conjectured incorrectly by Euler to have no nontrivial solutions. Proved by Elkies to have infinitely many nontrivial solutions, with a computer search by Frye determining the smallest nontrivial solution, 95800\(^{4}\) + 217519\(^{4}\) + 414560\(^{4}\) = 422481\(^{4}\).</dd>
</dl>

## Linear Diophantine equations
### One equation
The simplest linear Diophantine equation takes the form
\(ax+by=c,\)
where , and are given integers. The solutions are described by the following theorem:
*This Diophantine equation has a solution* (where and are integers) *if and only if* *is a multiple of the greatest common divisor of* *and*. *Moreover, if* (*x, y*) *is a solution, then the other solutions have the form* (*x* + *kv, y* − *ku*), *where* *is an arbitrary integer, and* *and* *are the quotients of* *and* *(respectively) by the greatest common divisor of* *and*. **Proof:** If is this greatest common divisor, Bézout's identity asserts the existence of integers and such that *ae* + *bf* = *d*. If is a multiple of , then *c* = *dh* for some integer , and (*eh, fh*) is a solution. On the other hand, for every pair of integers and , the greatest common divisor of and divides *ax* + *by*. Thus, if the equation has a solution, then must be a multiple of. If *a* = *ud* and *b* , then for every solution (*x, y*), we have
\(\begin{aligned}
a(x+kv) + b(y-ku) &= ax+by+k(av-bu) \\
&= ax+by+k(udv-vdu) \\
&= ax+by,
\end{aligned}\)
showing that (*x* + *kv, y* − *ku*) is another solution. Finally, given two solutions such that
\(ax_1 + by_1 = ax_2 + by_2 = c,\)
one deduces that \(u(x_2 - x_1) + v(y_2 - y_1) = 0.\)
As and are coprime, Euclid's lemma shows that divides *x*\(_{2}\) − *x*\(_{1}\), and thus that there exists an integer such that both
\(x_2 - x_1 = kv, \quad y_2 - y_1 = -ku.\)
Therefore,
\(x_2 = x_1 + kv, \quad y_2 = y_1 - ku,\)
which completes the proof.

### Chinese remainder theorem
The Chinese remainder theorem describes an important class of linear Diophantine systems of equations: let \(n_1, \dots, n_k\) be pairwise coprime integers greater than one, \(a_1, \dots, a_k\) be arbitrary integers, and be the product \(n_1 \cdots n_k.\) The Chinese remainder theorem asserts that the following linear Diophantine system has exactly one solution \((x, x_1, \dots, x_k)\) such that 0 ≤ *x* < *N*, and that the other solutions are obtained by adding to a multiple of :
\(\begin{aligned}
x &= a_1 + n_1\,x_1\\
&\;\;\vdots\\
x &= a_k + n_k\,x_k
\end{aligned}\)

### System of linear Diophantine equations
More generally, every system of linear Diophantine equations may be solved by computing the Smith normal form of its matrix, in a way that is similar to the use of the reduced row echelon form to solve a system of linear equations over a field. Using matrix notation every system of linear Diophantine equations may be written
\(AX = C,\)
where is an *m* × *n* matrix of integers, is an *n* × 1 column matrix of unknowns and is an *m* × 1 column matrix of integers.

The computation of the Smith normal form of provides two unimodular matrices (that is matrices that are invertible over the integers and have ±1 as determinant) and of respective dimensions *m* × *m* and *n* × *n*, such that the matrix
\(B = [b_{i,j}] = UAV\)
is such that is not zero for not greater than some integer , and all the other entries are zero. The system to be solved may thus be rewritten as
\(B (V^{-1}X) = UC.\)
Calling the entries of *V* and those of *D* = *UC*, this leads to the system
\(\begin{aligned}
& b_{i,i}y_i = d_i, \quad 1 \leq i \leq k \\
& 0y_i = d_i, \quad k < i \leq n.
\end{aligned}\)

This system is equivalent to the given one in the following sense: A column matrix of integers is a solution of the given system if and only if *x* = *Vy* for some column matrix of integers such that *By* = *D*.

It follows that the system has a solution if and only if divides for *i* ≤ *k* and *d\(_{i}\)* = 0 for *i* > *k*. If this condition is fulfilled, the solutions of the given system are
\(V\,
\begin{bmatrix}
\frac{d_1}{b_{1,1\\
\vdots\\
\frac{d_k}{b_{k,k\\
h_{k+1}\\
\vdots\\
h_n
\end{bmatrix}\, \)
where *h*\(_{*k*+1}\), …, *h\(_{n}\)* are arbitrary integers.

Hermite normal form may also be used for solving systems of linear Diophantine equations. However, Hermite normal form does not directly provide the solutions; to get the solutions from the Hermite normal form, one has to successively solve several linear equations. Nevertheless, Richard Zippel wrote that the Smith normal form "is somewhat more than is actually needed to solve linear diophantine equations. Instead of reducing the equation to diagonal form, we only need to make it triangular, which is called the Hermite normal form. The Hermite normal form is substantially easier to compute than the Smith normal form."

Integer linear programming amounts to finding some integer solutions (optimal in some sense) of linear systems that include also inequations. Thus systems of linear Diophantine equations are basic in this context, and textbooks on integer programming usually have a treatment of systems of linear Diophantine equations.

## Homogeneous equations
A homogeneous Diophantine equation is a Diophantine equation that is defined by a homogeneous polynomial. A typical such equation is the equation of Fermat's Last Theorem
\(x^d+y^d -z^d=0.\)

As a homogeneous polynomial in indeterminates defines a hypersurface in the projective space of dimension *n* − 1, solving a homogeneous Diophantine equation is the same as finding the rational points of a projective hypersurface.

Solving a homogeneous Diophantine equation is generally a very difficult problem, even in the simplest non-trivial case of three indeterminates (in the case of two indeterminates the problem is equivalent with testing if a rational number is the th power of another rational number). A witness of the difficulty of the problem is Fermat's Last Theorem (for *d* > 2, there is no integer solution of the above equation), which needed more than three centuries of mathematicians' efforts before being solved.

For degrees higher than three, most known results are theorems asserting that there are no solutions (for example Fermat's Last Theorem) or that the number of solutions is finite (for example Faltings' theorem).

For the degree three, there are general solving methods, which work on almost all equations that are encountered in practice, but no algorithm is known that works for every cubic equation.

### Degree two
Homogeneous Diophantine equations of degree two are easier to solve. The standard solving method proceeds in two steps. One has first to find one solution, or to prove that there is no solution. When a solution has been found, all solutions are then deduced.

For proving that there is no solution, one may reduce the equation modulo. For example, the Diophantine equation
\(x^2+y^2=3z^2,\)
does not have any other solution than the trivial solution (0, 0, 0). In fact, by dividing , and by their greatest common divisor, one may suppose that they are coprime. The squares modulo 4 are congruent to 0 and 1. Thus the left-hand side of the equation is congruent to 0, 1, or 2, and the right-hand side is congruent to 0 or 3. Thus the equality may be obtained only if , and are all even, and are thus not coprime. Thus the only solution is the trivial solution (0, 0, 0). This shows that there is no rational point on a circle of radius \(\sqrt{3}\), centered at the origin.

More generally, the Hasse principle allows deciding whether a homogeneous Diophantine equation of degree two has an integer solution, and computing a solution if there exist.

If a non-trivial integer solution is known, one may produce all other solutions in the following way.

#### Geometric interpretation
Let
\(Q(x_1, \ldots, x_n)=0\)
be a homogeneous Diophantine equation, where \(Q(x_1, \ldots, x_n)\) is a quadratic form (that is, a homogeneous polynomial of degree 2), with integer coefficients. The *trivial solution* is the solution where all \(x_i\) are zero. If \((a_1, \ldots, a_n)\) is a non-trivial integer solution of this equation, then \(\left(a_1, \ldots, a_n\right)\) are the homogeneous coordinates of a rational point of the hypersurface defined by. Conversely, if \(\left(\frac {p_1}q, \ldots, \frac {p_n}q \right)\) are homogeneous coordinates of a rational point of this hypersurface, where \(q, p_1, \ldots, p_n\) are integers, then \(\left(p_1, \ldots, p_n\right)\) is an integer solution of the Diophantine equation. Moreover, the integer solutions that define a given rational point are all sequences of the form
\(\left(k\frac{p_1}d, \ldots, k\frac{p_n}d\right),\)
where is any integer, and is the greatest common divisor of the \(p_i.\)

It follows that solving the Diophantine equation \(Q(x_1, \ldots, x_n)=0\) is completely reduced to finding the rational points of the corresponding projective hypersurface.

#### Parameterization
Let now \(A=\left(a_1, \ldots, a_n\right)\) be an integer solution of the equation \(Q(x_1, \ldots, x_n)=0.\) As is a polynomial of degree two, a line passing through crosses the hypersurface at a single other point, which is rational if and only if the line is rational (that is, if the line is defined by rational parameters). This allows parameterizing the hypersurface by the lines passing through , and the rational points are those that are obtained from rational lines, that is, those that correspond to rational values of the parameters.

More precisely, one may proceed as follows.

By permuting the indices, one may suppose, without loss of generality that \(a_n\ne 0.\) Then one may pass to the affine case by considering the affine hypersurface defined by
\(q(x_1,\ldots,x_{n-1})=Q(x_1, \ldots, x_{n-1},1),\)
which has the rational point
\(R= (r_1, \ldots, r_{n-1})=\left(\frac{a_1}{a_n}, \ldots, \frac{a_{n-1{a_n}\right).\)

If this rational point is a singular point, that is if all partial derivatives are zero at , all lines passing through are contained in the hypersurface, and one has a cone. The change of variables
\(y_i=x_i-r_i\)
does not change the rational points, and transforms into a homogeneous polynomial in *n* − 1 variables. In this case, the problem may thus be solved by applying the method to an equation with fewer variables.

If the polynomial is a product of linear polynomials (possibly with non-rational coefficients), then it defines two hyperplanes. The intersection of these hyperplanes is a rational flat, and contains rational singular points. This case is thus a special instance of the preceding case.

In the general case, consider the parametric equation of a line passing through :
\(\begin{aligned}
x_2 &= r_2 + t_2(x_1-r_1)\\
&\;\;\vdots\\
x_{n-1} &= r_{n-1} + t_{n-1}(x_1-r_1).
\end{aligned}\)
Substituting this in , one gets a polynomial of degree two in *x*, that is zero for *x*. It is thus divisible by *x*. The quotient is linear in *x*, and may be solved for expressing *x* as a quotient of two polynomials of degree at most two in \(t_2, \ldots, t_{n-1},\) with integer coefficients:
\(x_1=\frac{f_1(t_2, \ldots, t_{n-1})}{f_n(t_2, \ldots, t_{n-1})}.\)
Substituting this in the expressions for \(x_2, \ldots, x_{n-1},\) one gets, for *i* = 1, …, *n* − 1,
\(x_i=\frac{f_i(t_2, \ldots, t_{n-1})}{f_n(t_2, \ldots, t_{n-1})},\)
where \(f_1, \ldots, f_n\) are polynomials of degree at most two with integer coefficients.

Then, one can return to the homogeneous case. Let, for *i* = 1, …, *n*,
\(F_i(t_1, \ldots, t_{n-1})=t_1^2 f_i\left(\frac{t_2}{t_1}, \ldots, \frac{t_{n-1{t_1} \right),\)
be the homogenization of \(f_i.\) These quadratic polynomials with integer coefficients form a parameterization of the projective hypersurface defined by :
\(\begin{aligned}
x_1&= F_1(t_1, \ldots, t_{n-1})\\
&\;\;\vdots\\
x_n&= F_n(t_1, \ldots, t_{n-1}).
\end{aligned}\)

A point of the projective hypersurface defined by is rational if and only if it may be obtained from rational values of \(t_1, \ldots, t_{n-1}.\) As \(F_1, \ldots,F_n\) are homogeneous polynomials, the point is not changed if all are multiplied by the same rational number. Thus, one may suppose that \(t_1, \ldots, t_{n-1}\) are coprime integers. It follows that the integer solutions of the Diophantine equation are exactly the sequences \((x_1, \ldots, x_n)\) where, for *i* = 1,... , *n*,
\(x_i= k\,\frac{F_i(t_1, \ldots, t_{n-1})}{d},\)
where is an integer, \(t_1, \ldots, t_{n-1}\) are coprime integers, and is the greatest common divisor of the integers \(F_i(t_1, \ldots, t_{n-1}).\)

One could hope that the coprimality of could imply that *d* = 1. Unfortunately this is not the case, as shown in the next section.

#### Example of Pythagorean triples
The equation
\(x^2+y^2-z^2=0\)
is probably the first homogeneous Diophantine equation of degree two that has been studied. Its solutions are the Pythagorean triples. This is also the homogeneous equation of the unit circle. In this section, we show how the above method allows retrieving Euclid's formula for generating Pythagorean triples.

For retrieving exactly Euclid's formula, we start from the solution (−1, 0, 1), corresponding to the point (−1, 0) of the unit circle. A line passing through this point may be parameterized by its slope:
\(y=t(x+1).\)
Putting this in the circle equation
\(x^2+y^2-0,\)
one gets
\(x^2-1 +t^2(x+1)^2=0.\)
Dividing by *x* + 1, results in
\(x-1+t^2(x+1)=0,\)
which is easy to solve in :
\(x=\frac{1-t^2}{1+t^2}.\)
It follows
\(y=t(x+1) = \frac{2t}{1+t^2}.\)
Homogenizing as described above one gets all solutions as
\(\begin{aligned}
x&=k\,\frac{s^2-t^2}{d}\\
y&=k\,\frac{2st}{d}\\
z&=k\,\frac{s^2+t^2}{d},
\end{aligned}\)
where is any integer, and are coprime integers, and is the greatest common divisor of the three numerators. In fact, *d* = 2 if and are both odd, and *d* = 1 if one is odd and the other is even.

The *primitive triples* are the solutions where *k* = 1 and *s* > *t* > 0.

This description of the solutions differs slightly from Euclid's formula because Euclid's formula considers only the solutions such that , and are all positive, and does not distinguish between two triples that differ by the exchange of and ,

## Diophantine analysis
### Typical questions
The questions asked in Diophantine analysis include:

#Are there any solutions?
#Are there any solutions beyond some that are easily found by inspection?
#Are there finitely or infinitely many solutions?
#Can all solutions be found in theory?
#Can one in practice compute a full list of solutions?

These traditional problems often lay unsolved for centuries, and mathematicians gradually came to understand their depth (in some cases), rather than treat them as puzzles.

### Typical problem
The given information is that a father's age is 1 less than twice that of his son, and that the digits making up the father's age are reversed in the son's age (i.e. ). This leads to the equation 10*A* + *B* , thus 19*B* − 8*A*. Inspection gives the result *A* , *B* , and thus equals 73 years and equals 37 years. One may easily show that there is not any other solution with and positive integers less than 10.

Many well known puzzles in the field of recreational mathematics lead to diophantine equations. Examples include the cannonball problem, Archimedes's cattle problem and the monkey and the coconuts.

### 17th and 18th centuries
In 1637, Pierre de Fermat scribbled on the margin of his copy of *Arithmetica*: "It is impossible to separate a cube into two cubes, or a fourth power into two fourth powers, or in general, any power higher than the second into two like powers." Stated in more modern language, "The equation *a* = *c*has no solutions for any higher than 2." Following this, he wrote: "I have discovered a truly marvelous proof of this proposition, which this margin is too narrow to contain." Such a proof eluded mathematicians for centuries, however, and as such his statement became famous as Fermat's Last Theorem. It was not until 1995 that it was proven by the British mathematician Andrew Wiles.

In 1657, Fermat attempted to solve the Diophantine equation 61*x*\(^{2}\) + 1 (solved by Brahmagupta over 1000 years earlier). The equation was eventually solved by Euler in the early 18th century, who also solved a number of other Diophantine equations. The smallest solution of this equation in positive integers is *x* , *y* (see Chakravala method).

### Hilbert's tenth problem

In 1900, David Hilbert proposed the solvability of all Diophantine equations as the tenth of his fundamental problems. In 1970, Yuri Matiyasevich solved it negatively, building on work of Julia Robinson, Martin Davis, and Hilary Putnam to prove that a general algorithm for solving all Diophantine equations cannot exist.

### Diophantine geometry
Diophantine geometry, is the application of techniques from algebraic geometry which considers equations that also have a geometric meaning. The central idea of Diophantine geometry is that of a rational point, namely a solution to a polynomial equation or a system of polynomial equations, which is a vector in a prescribed field , when is *not* algebraically closed.

### Modern research
The oldest general method for solving a Diophantine equationor for proving that there is no solution is the method of infinite descent, which was introduced by Pierre de Fermat. Another general method is the Hasse principle that uses modular arithmetic modulo all prime numbers for finding the solutions. Despite many improvements these methods cannot solve most Diophantine equations.

The difficulty of solving Diophantine equations is illustrated by Hilbert's tenth problem, which was set in 1900 by David Hilbert; it was to find an algorithm to determine whether a given polynomial Diophantine equation with integer coefficients has an integer solution. Matiyasevich's theorem implies that such an algorithm cannot exist.

During the 20th century, a new approach has been deeply explored, consisting of using algebraic geometry. In fact, a Diophantine equation can be viewed as the equation of a hypersurface, and the solutions of the equation are the points of the hypersurface that have integer coordinates.

This approach led eventually to the proof by Andrew Wiles in 1994 of Fermat's Last Theorem, stated without proof around 1637. This is another illustration of the difficulty of solving Diophantine equations.

### Infinite Diophantine equations
An example of an infinite Diophantine equation is:
\(n = a^2 + 2b^2 + 3c^2 + 4d^2 + 5e^2 + \cdots,\)
which can be expressed as "How many ways can a given integer be written as the sum of a square plus twice a square plus thrice a square and so on?" The number of ways this can be done for each forms an integer sequence. Infinite Diophantine equations are related to theta functions and infinite dimensional lattices. This equation always has a solution for any positive. Compare this to:
\(n = a^2 + 4b^2 + 9c^2 + 16d^2 + 25e^2 + \cdots,\)
which does not always have a solution for positive. ## Exponential Diophantine equations
If a Diophantine equation has as an additional variable or variables occurring as exponents, it is an exponential Diophantine equation. Examples include:

* the Ramanujan–Nagell equation, 2\(^{*n*}\) − 7 = *x*\(^{2}\)
* the equation of the Fermat–Catalan conjecture and Beal's conjecture, *a\(^{m}\)* + *b\(^{n}\)* = *c\(^{k}\)* with inequality restrictions on the exponents
* the Erdős–Moser equation, 1\(^{*k*}\) + 2\(^{*k*}\) + ⋯ + (*m* − 1)\(^{*k*}\) = *m\(^{k}\)*

A general theory for such equations is not available; particular cases such as Catalan's conjecture and Fermat's Last Theorem have been tackled. However, the majority are solved via ad-hoc methods such as Størmer's theorem or even trial and error.
