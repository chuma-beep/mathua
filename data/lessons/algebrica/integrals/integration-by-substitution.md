# Integration by Substitution

Source: algebrica.org — CC BY-NC 4.0
https://algebrica.org/integration-by-substitution/

## How substitution simplifies integration

Integration by substitution is a technique used to simplify an integral by introducing a suitable substitution. When the [integral](../indefinite-integrals/) is not straightforward to compute, this method proves highly useful as it allows rewriting the integral of a function \\(f(x)\\) in terms of a new variable \\(u\\), simplifying the computation:

\\[\int f(g(x))\\,g'(x)\\,dx = \int f(u)\\,du\\]

The process involves the following steps:

+ Introduce a change of variable by defining \\( u = g(x) \\), where \\( g(x) \\) is an appropriately chosen function.
+ Compute the differential transformation, given by \\( du = g'(x)dx. \\)
+ Rewrite the integral in terms of \\( u \\), replacing \\( x \\) and \\( dx \\) accordingly, to obtain an equivalent expression that is often more straightforward to solve.
+ Once the integral is evaluated, revert to the original variable \\( x \\) to express the final result in its initial form.

> The key insight is that substitution reverses the chain rule: recognizing this connection makes it easier to identify when and how to apply the technique.

- - -

The method of substitution is a direct consequence of the [chain rule](../the-derivative-of-a-composite-function/) for derivatives. If \\( F(x) = H(g(x)) \\), then by the chain rule:

\\[
F'(x) = H'(g(x))\\, g'(x)
\\]

Therefore, whenever an integrand has the form \\(H'(g(x))\\, g'(x)\\) it is the derivative of the composite function \\( H(g(x)) \\). Integration by substitution simply reverses this process by introducing \\( u = g(x) \\), reducing the integral to:

\\[
\int H'(u)\\, du = H(u) + c
\\]

- - -
## Recognizing when to use substitution

Before proceeding to concrete examples, it is useful to understand when a substitution is likely to be effective. The technique is most natural when the integrand contains a [composite function](../composite-functions/). In many cases, the integral has the general form:

\\[
f(g(x))\\,g'(x)
\\]

or differs from it only by a constant factor. When this pattern appears, choosing \\( u = g(x) \\) simplifies the expression by reducing the composite structure to a single variable. A common signal is the presence of expressions such as \\( (ax+b)^n \\), \\( \sqrt{ax+b} \\), \\( \ln(ax+b) \\), or \\( e^{ax+b} \\). In these cases, the inner linear function \\( ax+b \\) is often a natural candidate for substitution. Similarly, in rational expressions of the form:

\\[
\frac{g'(x)}{g(x)}
\\]

the [derivative](../derivatives) of the denominator suggests the substitution \\( u = g(x) \\).

> In practice, the key idea is to look for an inner expression whose derivative also appears, exactly or up to a multiplicative constant, elsewhere in the integrand. When such a relationship is present, substitution typically transforms the integral into a simpler form.

- - -
## Substitution patterns

|                                        |                  |
| -------------------------------------- | ---------------- |
| \\[ \int f(g(x))\\, g'(x)\\, dx \\]    | \\[ u = g(x) \\] |
| \\[ \int (ax+b)^n\\, dx \\]            | \\[ u = ax+b \\] |
| \\[ \int e^{ax+b}\\]                   | \\[ u = ax+b \\] |
| \\[\int \ln(ax+b)\\, dx \\]            | \\[ u = ax+b \\] |
| \\[ \int \dfrac{g'(x)}{g(x)}\\, dx \\] | \\[ u = g(x) \\] |
- - -
## Example 1

Consider the following integral:

\\[\int (2x+1)^3 \\,dx\\]

---

Let \\( u = 2x + 1 \\), which simplifies the exponentiation. Differentiating both sides with respect to \\( x \\) we have:

\\[du = 2\\,dx\\]

Solving for \\( dx \\) we obtain:

\\[dx = \frac{du}{2}\\]

---

Expressing the integral entirely in terms of \\( u \\):

\\[\int u^3 \cdot \frac{du}{2} = \frac{1}{2} \int u^3 \\,du\\]

We now proceed to solve the integral in \\( u \\), which has been reduced to a power integral of the form \\( u^a \\). We obtain:

\\[\frac{1}{2} \cdot \frac{u^4}{4} + c = \frac{1}{8} u^4 + c\\]

Substituting back \\( u = 2x+1 \\), we get:

\\[\frac{1}{8} (2x+1)^4 + c\\]

- - -
## Example 2

Evaluate the following integral:

\\[\int \frac{1}{3x-5} \\,dx\\]

---

Let \\( u = 3x - 5 \\), which simplifies the denominator. Differentiating both sides with respect to \\( x \\), we have:

\\[du = 3\\,dx\\]

Solving for \\( dx \\), we obtain:

\\[dx = \frac{du}{3}\\]

---

Expressing the integral entirely in terms of \\( u \\):

\\[\int \frac{1}{u} \cdot \frac{du}{3} = \frac{1}{3} \int \frac{du}{u}\\]

We now proceed to solve the integral in \\( u \\), which has been reduced to a logarithmic integral of the form \\( 1/u \\). We obtain:

\\[\frac{1}{3} \ln |u| + c\\]

Substituting back \\( u = 3x - 5 \\), we get:

\\[\frac{1}{3} \ln |3x - 5| + c\\]

> Integration by substitution is an effective technique, but selecting the right substitution requires practice and the ability to recognize the structure of the integrand.

- - -
## Example 3

Evaluate the following integral:

\\[\int x \sin(x^2)\\,dx\\]

---

The substitution is less immediate here, since the integrand does not match the standard pattern as directly as in the previous examples. Let \\( u = x^2 \\), which simplifies the argument of the sine function. Differentiating both sides with respect to \\( x \\), we get:

\\[du = 2x\\,dx\\]

Solving for \\( dx \\) we obtain:

\\[dx = \frac{du}{2x}\\]

---

Rewriting everything in terms of \\( u \\), and since \\( du = 2x\\,dx \\), we have:

\\[\int x \sin(u) \cdot \frac{du}{2x} = \frac{1}{2} \int \sin(u)\\,du\\]

---

We now proceed to solve the integral in \\( u \\):

\\[\int \sin u\\,du = -\cos u\\]

Thus:

\\[\frac{1}{2} (-\cos u) + c = -\frac{1}{2} \cos u + c\\]

Substituting back \\( u = x^2 \\), we obtain:

\\[-\frac{1}{2} \cos(x^2) + c\\]

- - -
## Example 4

Evaluate the following integral:

\\[\int \cos x \sqrt{\sin x}\\,dx\\]

---

Let \\( u = \sin x \\), which simplifies the square root term. Differentiating both sides with respect to \\( x \\), we get:

\\[du = \cos x\\,dx\\]

Since \\( du = \cos x\\,dx \\), we can substitute directly into the integral.

---

Substituting \\( u = \sin x \\) and \\( du = \cos x\\,dx \\) into the integral, we obtain:

\\[\int \sqrt{u}\\,du = \int u^{1/2}\\,du\\]

---

We now compute the integral:

\\[\int u^{1/2}\\,du = \frac{u^{3/2}}{\frac{3}{2}} = \frac{2}{3}u^{3/2} + c\\]

Substituting back \\( u = \sin x \\), we obtain:

\\[\frac{2}{3} (\sin x)^{3/2} + c\\]

- - -
## Trigonometric substitutions

Trigonometric substitution applies when an integral involves [polynomial](../polynomials), [rational](../rational-functions/), or algebraic expressions that can be simplified using the [fundamental trigonometric identity](../pythagorean-identity/):

\\[\sin^2 x + \cos^2 x = 1\\]

which can be expressed in the following forms:

\\[
\begin{align*}
\cos^2 x &= 1 - \sin^2 x \\\\[0.5em]
\sec^2 x &= 1 + \tan^2 x \\\\[0.5em]
\tan^2 x &= \sec^2 x - 1
\end{align*}
\\]

To simplify an integral, choose an appropriate substitution based on the expression present in the integrand:

- If the integrand contains \\( 1 - x^2 \\), use \\( x = \sin u \\).
- If the integrand contains \\( 1 + x^2 \\), use \\( x = \tan u \\).
- If the integrand contains \\( x^2 - 1 \\), use \\( x = \sec u \\).

> A complete discussion of trigonometric substitution, including the geometric rationale and fully worked examples, is presented in the dedicated section [Trigonometric Substitution for Integrals](../trigonometric-substitution-for-integrals/).

- - -
## Example 5

Evaluate the following integral:

\\[\int \frac{1}{\sqrt{9-x^2}}\\,dx\\]

---

For expressions of the form \\( a^2 - x^2 \\), a natural substitution is:

\\[x = 3\sin u\\]

Differentiating both sides:

\\[dx = 3\cos u\\,du\\]

---

Substituting \\( x = 3\sin u \\) in the denominator:

\\[\sqrt{9-x^2} = \sqrt{9-9\sin^2 u} = \sqrt{9(1-\sin^2 u)}\\]

Since \\( \sin^2 u + \cos^2 u = 1 \\), we obtain:

\\[\sqrt{9(1-\sin^2 u)} = \sqrt{9\cos^2 u} = 3\cos u\\]

Thus, the integral transforms into:

\\[\int \frac{3\cos u\\,du}{3\cos u} = \int du = u+c\\]

> This step assumes \\( \cos u \geq 0 \\), which holds since the substitution \\( x = 3\sin u \\) implies \\( u \in [-\pi/2,\\,\pi/2] \\).

---

From the substitution \\( x = 3\sin u \\), solving for \\( u \\) via the [arcsine](../arcsine-and-arccosine/) function gives:

\\[u = \arcsin\left(\frac{x}{3}\right)\\]

Thus:

\\[\arcsin\left(\frac{x}{3}\right) + c\\]

- - -
## Substitution rule for definite integrals

When applying substitution to evaluate [definite integrals](../definite-integrals), the limits of integration must be adjusted to reflect the new variable. If the limits are not changed, the result will be incorrect. Given the substitution \\( u = g(x) \\), we have:

\\[\int_{a}^{b} f(g(x))\\,g'(x)\\,dx = \int_{g(a)}^{g(b)} f(u)\\,du\\]

- - -
## Example 6

Evaluate the following definite integral:

\\[\int_{0}^{1} x\cos(x^2)\\,dx\\]

---

Using the substitution \\( u = x^2 \\), we get:

\\[du = 2x\\,dx \qquad dx = \frac{du}{2x}\\]

---

The limits of integration must be updated. When \\( x = 0 \\), then \\( u = 0 \\); when \\( x = 1 \\), then \\( u = 1 \\). In this case the transformed limits coincide with the original ones, though this is not generally the case. We obtain:

\\[\int_{0}^{1} x\cos(x^2)\\,dx = \int_{0}^{1} \cos u \cdot \frac{du}{2} = \frac{1}{2}\int_{0}^{1} \cos u\\,du\\]

Evaluating the integral we obtain:

\\[\frac{1}{2}\sin u\Big|_{0}^{1} = \frac{1}{2}(\sin 1 - \sin 0) = \frac{\sin 1}{2}\\]

- - -
## Flowchart

- `Integral to solve`
  - **IF** the integrand matches a known form
    - _integrate directly_ apply the corresponding standard formula (power, exponential, trigonometric…)
  - `ELSE IF` the integrand contains \\( a^2 - x^2 \\) or \\( a^2 + x^2 \\) or \\( x^2 - a^2 \\)
    - _choose the trigonometric substitution_
       \\[\begin{align*}
      a^2 - x^2 &\rightarrow x = a\sin u \\\\[3pt]
      a^2 + x^2 &\rightarrow x = a\tan u \\\\[3pt]
      x^2 - a^2 &\rightarrow x = a\sec u
      \end{align*}\\]
    - _compute \\( dx \\) from the substitution_
      differentiate both sides: e.g. \\( x = a\sin u \rightarrow dx = a\cos u\\,du \\)
    - _apply the Pythagorean identity_
      the radical collapses to a single trigonometric function
    - _integrate in \\( u \\)_
      the result is now a standard trigonometric integral
    - **IF** the integral is definite
      - _update the limits of integration_
        replace \\( a \\) and \\( b \\) with the corresponding values of \\( u \\) before evaluating
      - _evaluate and stop_
        the result is a number, no back-substitution needed
    - `ELSE`
      - _back-substitute using the inverse function_
        express the result in terms of the original variable \\( x \\)
  - `ELSE IF` the integrand has the form \\( f(g(x)) \cdot g'(x) \\)
    - _set \\( u = g(x) \\)_
      choose the inner function whose derivative appears in the integrand
    - _compute \\( du = g'(x)\\,dx \\)_
      if needed, solve for \\( dx \\) and adjust for any constant factor
    - _rewrite the integral entirely in \\( u \\)_
      every occurrence of \\( x \\) and \\( dx \\) must be replaced
    - _integrate in \\( u \\)_
      apply the appropriate standard formula
    - **IF** the integral is definite
      - _update the limits of integration_
         replace \\( a \\) and \\( b \\) with \\( g(a) \\) and \\( g(b) \\) before evaluating
      - _evaluate and stop_
        the result is a number, no back-substitution needed
    - `ELSE`
      - _substitute back \\( u = g(x) \\)_
        express the antiderivative in terms of the original variable \\( x \\)
  - `ELSE`
    - _simplify the integrand first_
      expand, factor, use trigonometric identities, or split into partial fractions
    - _restart_