> Content sourced from [Algebrica](https://algebrica.org/lagrange-theorem/) — CC BY-NC 4.0

## Statement

The Lagrange’s theorem, also known as the mean value theorem, states the following. Consider a function \\( f(x) \\), [continuous](<../continuous-functions/>) in the closed and bounded interval \\([a, b]\\) and differentiable at every point inside the interval. Then, there exists at least one point \\( c \\) inside the interval such that the following relation holds:

\\[f’ \left (c \right ) = \frac{f(b)-f(a)}{b-a} \\]

##### This means that there exists at least one point where the derivative of the function is equal to the slope of the secant line connecting \\(a\\) and \\(b\\). In other words, at some point in the interval, the instantaneous rate of change of the function matches its average rate of change.

## A geometric view of Lagrange’s theorem

From a geometric point of view, the theorem states that there exists at least one point \\( c \\) where the tangent line at that point is parallel to the secant line connecting points \\( A \\) and \\( B \\) on the graph.

![Lagrange’s theorem.](/diagrams/algebrica/lagrange-theoreme-4.png)

In the right-angled triangle \\( ABH \\), we have \\(\overline{BH} = \overline{AH} \cdot \tan{\alpha}\\) that is:

\\[\tan{\alpha} = \frac{\overline{BH}}{\overline{AH}} \\]

We have:

\\[\overline{BH} = f(b)-f(a), \quad \overline{AH} = b-a \\]

The slope of segment \\( AB \\) is equal to \\(\tan\alpha\\) that is:

\\[\tan{\alpha} = \frac{f(b)-f(a)}{b-a} \\]

Since the tangent at \\( c \\) to the curve is parallel to \\( AB \\), both have the same slope, hence:

\\[f’ \left(c \right) = \frac{f(b)-f(a)}{b-a} \\]

## Proof

To prove Lagrange’s theorem we define an auxiliary function \\( \varphi(x) \\) as follows:

\\[\varphi(x) = f(x)-f(a)-\frac{f(b)-f(a)}{b-a} \cdot (x-a) \\]

We verify that \\( \varphi(x) \\) satisfies the hypotheses of [Rolle’s Theorem](<../rolles-theorem/>):

  * \\(f(x)\\) is continuous on the closed interval \\([a, b]\\).
  * \\(f(x)\\) is differentiable on the open interval \\((a, b)\\).
  * \\(\varphi(a) = \varphi(b)\\).


By calculating \\( \varphi(a) \\) and \\( \varphi(b) \\), we obtain:

\\[\varphi(a) = f(a)-f(a)-\frac{f(b)-f(a)}{b-a} \cdot (a-a) = 0 \\]

\\[\begin{align} \varphi(b) &= f(b)-f(a)-\frac{f(b)-f(a)}{b-a} \cdot (b-a)\\\\[0.5em] &= f(b)-(f(b)-f(a)) = 0 \end{align} \\]

Therefore, \\( \varphi(a) = \varphi(b) = 0 \\).


Applying Rolle’s Theorem to \\( \varphi(x) \\), we find that there exists at least one point \\( c \\) within the interval \\( ]a, b[ \\) such that \\( \varphi’ \left(c \right) = 0 \\). Calculating the derivative of \\( \varphi(x) \\), we have:

\\[\varphi’(x) = f’(x)-\frac{f(b)-f(a)}{b-a}\\]


Now, we calculate the derivative of \\( \varphi(x) \\) at the point \\( c \\) and set it equal to 0. From this, we obtain:

\\[\varphi’ \left(c \right) = f’ \left(c \right)-\frac{f(b)-f(a)}{b-a} = 0\\]

That is:

\\[f’ \left(c \right)=\frac{f(b)-f(a)}{b-a}\\]

which corresponds exactly to the thesis we wanted to prove.

## A special case: when the derivative is zero everywhere

From Lagrange’s theorem, it follows that if a function \\( f(x) \\) is continuous on the interval \\([a,b]\\), differentiable on the interval \\( ]a,b[ \\), and \\( f’(x) \\) is zero at every point in the interior of the interval, then \\( f(x) \\) is constant on the entire interval \\([a,b]\\).

![](/diagrams/algebrica/lagrange-theorem-2.png)

Indeed, if we take a point \\(\overline{x} \in [a,b]\\) and apply the theorem to the interval \\([a, \overline{x}]\\), then there exists a point \\(c \in ]a, \overline{x}[\\) such that: \\[f’\left( c \right) = \frac{f(\overline{x})-f(a)}{\overline{x}-a} \\]

Since \\( f’(\overline{x}) = 0 \\) for every point in \\( ]a,b[ \\), it follows that \\( f’\left( c \right) = 0 \\) as well. For this reason, it must be:

\\[f(\overline{x}) – f(a) = 0 \rightarrow f(\overline{x}) = f(a) \, \forall \, \overline{x} \in [a,b] \\]

For this reason, \\( f \\) is constant on the entire interval \\([a,b]\\).

## A numerical example

To see the theorem at work, let us examine a concrete case. We want to identify a point \\( c \\) inside a given interval where the instantaneous rate of change of a function coincides with its average rate of change on that interval. The computation also shows that such a point is not something one can usually predict by inspection. Consider the [polynomial](<../polynomials/>):

\\[f(x) = x^3 - 4x^2 + x + 6 \\]

on the closed interval \\( [1,4] \\). Being a polynomial, \\( f \\) is [continuous](<../continuous-functions/>) on \\( [1,4] \\) and differentiable on \\( (1,4) \\). The assumptions of Lagrange’s theorem are therefore satisfied. We begin by evaluating the function at the endpoints:

\\[f(1) = 1 - 4 + 1 + 6 = 4 \\] \\[f(4) = 64 - 64 + 4 + 6 = 10 \\]

The slope of the secant line through the points \\( (1,4) \\) and \\( (4,10) \\) is

\\[\frac{f(4) - f(1)}{4 - 1} = \frac{10 - 4}{3} = 2 \\]

This number represents the average rate of change of \\( f \\) over \\( [1,4] \\). Next we compute the derivative:

\\[f’(x) = 3x^2 - 8x + 1 \\]

We look for values of \\( c \\) such that \\(f’(c) = 2.\\) This leads to the equation:

\\[3c^2 - 8c + 1 = 2 \\] \\[3c^2 - 8c - 1 = 0 \\]

Applying the [quadratic formula](<../quadratic-formula/>) gives:

\\[c = \frac{8 \pm \sqrt{64 + 12}}{6} = \frac{8 \pm \sqrt{76}}{6} = \frac{4 \pm \sqrt{19}}{3} \\]

We obtain two candidates:

\\[c_1 = \frac{4 - \sqrt{19}}{3} \approx 0.21 \\] \\[c_2 = \frac{4 + \sqrt{19}}{3} \approx 2.79 \\]

Lagrange’s theorem guarantees a solution inside the open interval \\( ]1,4[ \\). Among the two values found, only:

\\[c = \frac{4 + \sqrt{19}}{3} \approx 2.79 \\]

lies in \\( ]1,4[ \\). The other root falls outside the interval and is therefore not relevant in this context. At this point, the tangent line to the graph of \\( f \\) is parallel to the secant line joining \\( (1,48) \\) and \\( (4,10) \\).

###### This example shows that the equation \\( f’(c)=2 \\) may admit more than one solution, yet only those inside the interval are meaningful for the theorem.

The solution is: \\[c = \dfrac{4+\sqrt{19}}{3} \\]

## Selected references

  * **Harvard University, O. Knill**. [The Mean Value Theorem](https://people.math.harvard.edu/~knill/teaching/math1a_2014/handouts/26-rolle.pdf)

  * **Stony Brook University**. [Lecture 18: Mean Value Theorem](https://www.math.stonybrook.edu/Videos/MAT131Online/Handouts/Lecture-18-Handout.pdf)

  * **University of California Davis, D. Kouba**. [Mean Value Theorem – Problems and Proofs](https://www.math.ucdavis.edu/~kouba/CalcOneDIRECTORY/meanvaluetheoremdirectory/MeanValueTheorem.html)

  * **University of Cambridge, W. T. Gowers**. [What is the point of the Mean Value Theorem?](https://www.dpmms.cam.ac.uk/~wtg10/meanvalue.html)
