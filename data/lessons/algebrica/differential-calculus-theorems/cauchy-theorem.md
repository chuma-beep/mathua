> Content sourced from [Algebrica](https://algebrica.org/cauchy-theorem/) — CC BY-NC 4.0

## Introduction

**Cauchy’s Theorem** establishes a relationship between the changes of two functions over a given interval. Specifically, if \\( f(x) \\) and \\( g(x) \\) are [continuous](<../continuous-functions/>) on a closed interval \\([a, b]\\) and differentiable in its interior, with \\( g’(x) \neq 0 \\), then there exists at least one point \\( c \\) in \\( (a, b) \\) where the ratio of their [derivatives](<../derivatives>) matches the ratio of their overall change across the interval:

## Statement

The Cauchy’s theorem states the following. Let \\( f(x) \\) and \\( g(x) \\) be two functions such that:

  * \\( f(x) \\) and \\( g(x) \\) are continuous on the interval \\([a, b]\\).
  * \\( f(x) \\) and \\( g(x) \\) are differentiable at every point in the interior of the interval.
  * \\( g’(x) \neq 0 \\) for every \\( x \\) in the interior of \\([a, b]\\).


Then, there exists at least one point \\( c \\) in the interior of the interval \\([a, b]\\) such that:

\\[\frac{f’ \left(c \right)}{g’ \left(c \right)} = \frac{f(b)-f(a)}{g(b)-g(a)} \\]

that is, the ratio of the increments of the functions \\( f(x) \\) and \\( g(x) \\) over the interval \\([a, b]\\) is equal to the ratio of their respective derivatives evaluated at a point \\( c \\) in the interior of the interval.

Setting \\(g(x) = x\\) reduces Cauchy’s theorem directly to [Lagrange’s theorem](<../lagrange-theorem/>). With \\(g’(x) = 1\\) and \\(g(b) - g(a) = b - a\\), the conclusion takes the form: \\[\frac{f’(c)}{1} = \frac{f(b) - f(a)}{b - a} \\] Lagrange’s theorem is therefore a special case of Cauchy’s theorem, obtained when one of the two functions is the identity.

##### Cauchy’s theorem provides the theoretical basis for the proof of [L’Hôpital’s Rule](<../hopital-rule/>).

## Proof of Cauchy’s theorem

To prove the theorem, define a new function \\( \varphi(x) \\) where \\( \lambda \\) is a constant to be determined.:

\\[\varphi(x) = f(x)-\lambda g(x) \\]

We want to choose \\( \lambda \\) such that \\( \varphi(a) = \varphi(b) \\). From this condition, we obtain:

\\[\lambda = \frac{f(b)-f(a)}{g(b)-g(a)} \\]


Now the function \\( \varphi(x) \\) is continuous on \\([a, b]\\) and differentiable on \\((a, b)\\), with \\( \varphi(a) = \varphi(b) \\). We can apply the [Rolle’s Theorem](<../rolles-theorem>), which guarantees that there exists at least one point \\( c \in (a, b) \\) such that \\( \varphi’ \left(c \right) = 0 \\). By calculating the derivative of \\( \varphi(x) \\) and setting \\( \varphi’ \left(c \right) = 0 \\), we obtain:

\\[\varphi’(x) = f’(x) - \lambda g’(x) \\] \\[f’ \left(c \right) = \lambda g’ \left(c \right) \\]


Substituting the value of \\( \lambda \\), we get:

\\[f’ \left(c \right) = \frac{f(b) - f(a)}{g(b) - g(a)} g’ \left(c \right) \\]

By dividing both sides by \\( g’ \left(c \right ) \\), we obtain:

\\[\frac{f’ \left(c \right)}{g’ \left(c \right)} = \frac{f(b) - f(a)}{g(b) - g(a)} \\]

Thus, we have proven the conclusion of the theorem.

## Example 1

Let’s verify, for example, that the theorem is applicable to the functions \\(f(x)\\) and \\(g(x) \\) in the interval [1,3]:

\\[f(x) = 2x^2-4x+2\\] \\[g(x) = x^2\\]

The two functions are [polynomials](<../polynomials>), therefore they are continuous and differentiable for every \\( x \in \mathbb{R} \\). Moreover, \\( g’(x) = 2x \neq 0 \\). This satisfies the hypotheses of the theorem.


Let’s now verify the existence of a point \\( c \in (1,3) \\) such that:

\\[\frac{f’ \left(c \right)}{g’ \left(c \right)} = \frac{f(3)-f(1)}{g(3)-g(1)} \\]

We obtain:

\\[\frac{f(3)-f(1)}{g(3)-g(1)} = \frac{18-12+2-2+4-2}{9-1} = \frac{8}{8} = 1\\]


Now let’s calculate:

\\[\frac{f’ \left(c \right)}{g’ \left(c \right)} = \frac{4c-4}{2c} \\]

From \\(1\\), the equality becomes:

\\[\begin{align} \frac{4c - 4}{2c} &= 1 \\\\[0.5em] \frac{4c - 4}{2c} &= \frac{2c}{2c} \\\\[0.5em] 4c - 4 &= 2c \\\\[0.5em] 2c &= 4 \\\\[0.5em] c &= 2\\\ \end{align} \\]

Since \\( c = 2 \in [1,3] \\), the theorem is verified.

## A note on the hypothesis \\(g’(x) \neq 0\\)

Consider the particular case in which \\(g(b) = g(a)\\). In this situation, the denominator of the following ratio is zero, and the expression is undefined: \\[\frac{f(b)-f(a)}{g(b)-g(a)} \\] For this reason, the constant cannot be introduced: \\[\lambda = \frac{f(b)-f(a)}{g(b)-g(a)} \\] As a result, the proof utilising the following auxiliary function is no longer valid: \\[\varphi(x) = f(x) - \lambda g(x) \\] Under the hypotheses of Cauchy’s theorem, this situation does not arise. The condition \\(g’(x) \neq 0\\) within the interior of the interval ensures that g is strictly monotonic, which in turn guarantees that \\(g(b) \neq g(a)\\). Therefore, the case \\(g(b) = g(a)\\) is precluded by the theorem’s assumptions.

## Selected references

  * **University of Florida J. Keesling**. [The Cauchy Mean Value Theorem](https://people.clas.ufl.edu/kees/files/CauchyMeanValue.pdf)

  * **University of Maryland**. [The Cauchy Mean Value Theorem and Consequences](https://math.umd.edu/~immortal/MATH410/lecturenotes/ch4-4.pdf)

  * **UC Berkeley A. Vizeff**. [Mean Value Theorem and L’Hôpital’s Rule](https://math.berkeley.edu/~avizeff/calculus-I-F23/lecture-16.pdf)
