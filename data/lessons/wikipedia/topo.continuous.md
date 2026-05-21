> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Continuous_function) — CC BY-SA 4.0

# Continuity in topological spaces

In mathematics, a **continuous function** is a function such that a small variation of the argument induces a small variation of the value of the function. This implies there are no abrupt changes in value, known as *discontinuities*. More precisely, a function is continuous if arbitrarily small changes in its value can be assured by restricting to sufficiently small changes of its argument. A **discontinuous function** is a function that Until the 19th century, mathematicians largely relied on intuitive notions of continuity and considered only continuous functions. The epsilon–delta definition of a limit was introduced to formalize the definition of continuity.

Continuity is one of the core concepts of calculus and mathematical analysis, where arguments and values of functions are real and complex numbers. The concept has been generalized to functions between metric spaces and between topological spaces. The latter are the most general continuous functions, and their definition is the basis of topology.

A stronger form of continuity is uniform continuity. In order theory, especially in domain theory, a related concept of continuity is Scott continuity.

As a practical example, the function *H*(*t*) denoting the height of a growing flower at time would be considered continuous. In contrast, the function *M*(*t*) denoting the amount of money in a bank account at time would be considered discontinuous since it "jumps" at each point in time when money is deposited or withdrawn.

## History
A form of the epsilon–delta definition of continuity was first given by Bernard Bolzano in 1817. Augustin-Louis Cauchy defined continuity of \(y = f(x)\) as follows: an infinitely small increment \(\alpha\) of the independent variable \(x\) always produces an infinitely small change \(f(x+\alpha)-f(x)\) of the dependent variable \(y\) (see e.g. *Cours d'Analyse*, p. 34). Cauchy defined infinitely small quantities in terms of variable quantities, and his definition of continuity closely parallels the infinitesimal definition used today (see microcontinuity). The formal definition and the distinction between pointwise continuity and uniform continuity were first given by Bolzano in the 1830s, but the work wasn't published until the 1930s. Like Bolzano, Karl Weierstrass considered that a function \(y = f(x)\) at a point \(x=c\), that is \(f(x)\big|_{x=c}\), is continuous the value of \(f(c)\) is defined at both sides \(f(x)\big|_{x\to c^+}\) and \(f(x)\big|_{x\to c^-}\) the function is defined too. Édouard Goursat assumed continuity provided that the function is defined at \(f(c)\) and (at least) on one side of \(f(x)\big|_{x\to c}\), whereas Camille Jordan "allowed" it even if the function was defined only at \(x=c\). All three of those nonequivalent definitions of pointwise continuity are still in use. Eduard Heine provided the first published definition of uniform continuity in 1872, but based these ideas on lectures given by Peter Gustav Lejeune Dirichlet in 1854.

### Definition

A real function that is a function from real numbers to real numbers can be represented by a graph in the Cartesian plane; such a function is continuous if, roughly speaking, the graph is a single unbroken curve whose domain is the entire real line. A more mathematically rigorous definition is given below.

Continuity of real functions is usually defined in terms of limits. A function *f* with variable is *continuous at* the real number , if the limit of \(f(x),\) as tends to , is equal to \(f(c).\)

There are several different definitions of the (global) continuity of a function, which depend on the nature of its domain.

A function is continuous on an open interval if the interval is contained in the function's domain and the function is continuous at every interval point. A function that is continuous on the interval \((-\infty, +\infty)\) (the whole real line) is often called simply a continuous function; one also says that such a function is *continuous everywhere*. For example, all polynomial functions are continuous everywhere.

A function is continuous on a semi-open or a closed interval; if the interval is contained in the domain of the function, the function is continuous at every interior point of the interval, and the value of the function at each endpoint that belongs to the interval is the limit of the values of the function when the variable tends to the endpoint from the interior of the interval. For example, the function \(f(x) = \sqrt{x}\) is continuous on its whole domain, which is the semi-open interval \([0,+\infty).\)

Many commonly encountered functions are partial functions that have a domain formed by all real numbers, except some isolated points. Examples include the reciprocal function \(x \mapsto \frac {1}{x}\) and the tangent function \(x\mapsto \tan x.\) When they are continuous on their domain, one says, in some contexts, that they are continuous, although they are not continuous everywhere. In other contexts, mainly when one is interested in their behavior near the exceptional points, one says they are discontinuous.

A partial function is *discontinuous* at a point if the point belongs to the topological closure of its domain, and either the point does not belong to the domain of the function or the function is not continuous at the point. For example, the functions \(x\mapsto \frac {1}{x}\) and \(x\mapsto \sin(\frac {1}{x})\) are discontinuous at 0, and remain discontinuous whichever value is chosen for defining them at 0. A point where a function is discontinuous is called a *discontinuity*.

Using mathematical notation, several ways exist to define continuous functions in the three senses mentioned above.

Let \(f : D \to \R\) be a function whose domain \(D\) is contained in \(\R\) of real numbers.

Some (but not all) possibilities for \(D\) are:
*\(D\) is the whole real line; that is, \(D = \R\)
*\(D\) is a closed interval of the form \(D = [a, b] = \{x \in \R \mid a \leq x \leq b \} ,\) where and are real numbers
*\(D\) is an open interval of the form \(D = (a, b) = \{x \in \R \mid a < x < b \},\) where and are real numbers

In the case of an open interval, \(a\) and \(b\) do not belong to \(D\), and the values \(f(a)\) and \(f(b)\) are not defined, and if they are, they do not matter for continuity on \(D\).

#### Definition in terms of limits of functions
The function *f* is *continuous at some point* *c* of its domain if the limit of \(f(x),\) as *x* approaches *c* through the domain of *f*, exists and is equal to \(f(c).\) In mathematical notation, this is written as

\[
\lim_{x \to c}{f(x)} = f(c).
\]

In detail this means three conditions: first, *f* has to be defined at *c* (guaranteed by the requirement that *c* be in the domain of *f*). Second, the limit of that equation has to exist. Third, the value of this limit must equal \(f(c).\)

(Here, we have assumed that the domain of *f* does not have any isolated points.)

#### Definition in terms of neighborhoods
A neighborhood of a point *c* is a set that contains, at least, all points within some fixed distance of *c*. Intuitively, a function is continuous at a point *c* if the range of *f* over the neighborhood of *c* shrinks to a single point \(f(c)\) as the width of the neighborhood around *c* shrinks to zero. More precisely, a function *f* is continuous at a point *c* of its domain if, for any neighborhood \(N_1(f(c))\) there is a neighborhood \(N_2(c)\) in its domain such that \(f(x) \in N_1(f(c))\) whenever \(x\in N_2(c).\)

As neighborhoods are defined in any topological space, this definition of a continuous function applies not only for real functions but also when the domain and the codomain are topological spaces and is thus the most general definition. It follows that a function is automatically continuous at every isolated point of its domain. For example, every real-valued function on the integers is continuous.

#### Definition in terms of limits of sequences

One can instead require that for any sequence \((x_n)_{n \in \N}\) of points in the domain which converges to *c*, the corresponding sequence \(\left(f(x_n)\right)_{n\in \N}\) converges to \(f(c).\) In mathematical notation,
\[
\forall (x_n)_{n \in \N} \subset D:\lim_{n\to\infty} x_n = c \Rightarrow \lim_{n\to\infty} f(x_n) = f(c)\,.
\]

#### Weierstrass and Jordan definitions (epsilon–delta) of continuous functions

Explicitly including the definition of the limit of a function, we obtain a self-contained definition: Given a function \(f : D \to \mathbb{R}\) as above and an element \(x_0\) of the domain \(D\), \(f\) is said to be continuous at the point \(x_0\) when the following holds: For any positive real number \(\varepsilon > 0,\) however small, there exists some positive real number \(\delta > 0\) such that for all \(x\) in the domain of \(f\) with \(x_0 - \delta < x < x_0 + \delta,\) the value of \(f(x)\) satisfies

\[
f\left(x_0\right) - \varepsilon < f(x) < f(x_0) + \varepsilon.
\]

Alternatively written, continuity of \(f : D \to \mathbb{R}\) at \(x_0 \in D\) means that for every \(\varepsilon > 0,\) there exists a \(\delta > 0\) such that for all \(x \in D\):

\[
\left|x - x_0\right| < \delta ~~\text{ implies }~~ |f(x) - f(x_0)| < \varepsilon.
\]

More intuitively, we can say that if we want to get all the \(f(x)\) values to stay in some small neighborhood around \(f\left(x_0\right),\) we need to choose a small enough neighborhood for the \(x\) values around \(x_0.\) If we can do that no matter how small the \(f(x_0)\) neighborhood is, then \(f\) is continuous at \(x_0.\)

In modern terms, this is generalized by the definition of continuity of a function with respect to a basis for the topology, here the metric topology.

Weierstrass had required that the interval \(x_0 - \delta < x < x_0 + \delta\) be entirely within the domain \(D\), but Jordan removed that restriction.

#### Definition in terms of control of the remainder
In proofs and numerical analysis, we often need to know how fast limits are converging, or in other words, control of the remainder. We can formalize this to a definition of continuity.
A function \(C: [0,\infty) \to [0,\infty]\) is called a control function if
* *C* is non-decreasing
*\(\inf_{\delta > 0} C(\delta) = 0\)

A function \(f : D \to R\) is *C*-continuous at \(x_0\) if there exists such a neighbourhood \(N(x_0)\) that

\[
|f(x) - f(x_0)| \leq C\left(\left|x - x_0\right|\right) \text{ for all } x \in D \cap N(x_0)
\]

A function is continuous in \(x_0\) if it is *C*-continuous for some control function *C*.

This approach leads naturally to refining the notion of continuity by restricting the set of admissible control functions. For a given set of control functions \(\mathcal{C}\) a function is for every subset \(B \subseteq Y,\) \(f^{-1}\left(\operatorname{int}_Y B\right) \subseteq \operatorname{int}_X\left(f^{-1}(B)\right);\)
