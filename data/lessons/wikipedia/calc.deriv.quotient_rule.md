> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Quotient_rule) — CC BY-SA 4.0

# Quotient rule

In calculus, the **quotient rule** is a method of finding the derivative of a function that is the ratio of two differentiable functions. Let \(h(x)=\frac{f(x)}{g(x)}\), where both and are differentiable and \(g(x)\neq 0.\) The quotient rule states that the derivative of *h*(*x*) is
\(h'(x) = \frac{f'(x)g(x) - f(x)g'(x)}{(g(x))^2}.\)

It is provable in many ways by using other derivative rules.

## Examples
### Example 1: Basic example
Given \(h(x)=\frac{e^x}{x^2}\), let \(f(x)=e^x, g(x)=x^2\), then using the quotient rule:
\[
\begin{aligned}
 \frac{d}{dx} \left(\frac{e^x}{x^2}\right) &= \frac{\left(\frac{d}{dx}e^x\right)(x^2) - (e^x)\left(\frac{d}{dx} x^2\right)}{(x^2)^2} \\
 &= \frac{(e^x)(x^2) - (e^x)(2x)}{x^4} \\
 &= \frac{x^2 e^x - 2x e^x}{x^4} \\
 &= \frac{x e^x - 2 e^x}{x^3} \\
 &= \frac{e^x(x - 2)}{x^3}.
 \end{aligned}
\]

### Example 2: Derivative of tangent function
The quotient rule can be used to find the derivative of \(\tan x = \frac{\sin x}{\cos x}\) as follows:

\[
\begin{aligned}
 \frac{d}{dx} \tan x &= \frac{d}{dx} \left(\frac{\sin x}{\cos x}\right) \\
 &= \frac{\left(\frac{d}{dx}\sin x\right)(\cos x) - (\sin x)\left(\frac{d}{dx}\cos x\right)}{\cos^2 x} \\
 &= \frac{(\cos x)(\cos x) - (\sin x)(-\sin x)}{\cos^2 x} \\
 &= \frac{\cos^2 x + \sin^2 x}{\cos^2 x} \\
 &= \frac{1}{\cos^2 x} = \sec^2 x.
 \end{aligned}
\]

## Reciprocal rule

The reciprocal rule is a special case of the quotient rule in which the numerator \(f(x)=1\). Applying the quotient rule gives
\[
h'(x)=\frac{d}{dx}\left[\frac{1}{g(x)}\right]=\frac{0 \cdot g(x) - 1 \cdot g'(x)}{g(x)^2}=\frac{-g'(x)}{g(x)^2}.
\]

Utilizing the chain rule yields the same result.

## Proofs
### Proof from derivative definition and limit properties
Let \(h(x) = \frac{f(x)}{g(x)}.\) Applying the definition of the derivative and properties of limits gives the following proof, with the term \(f(x) g(x)\) added and subtracted to allow splitting and factoring in subsequent steps without affecting the value:
\[
\begin{aligned}
 h'(x) &= \lim_{k\to 0} \frac{h(x+k) - h(x)}{k} \\
 &= \lim_{k\to 0} \frac{\frac{f(x+k)}{g(x+k)} - \frac{f(x)}{g(x)}}{k} \\
 &= \lim_{k\to 0} \frac{f(x+k)g(x) - f(x)g(x+k)}{k \cdot g(x)g(x+k)} \\
 &= \lim_{k\to 0} \frac{f(x+k)g(x) - f(x)g(x+k)}{k} \cdot \lim_{k\to 0}\frac{1}{g(x)g(x+k)} \\
 &= \lim_{k\to 0} \left[\frac{f(x+k)g(x) - f(x)g(x) + f(x)g(x) - f(x)g(x+k)}{k} \right] \cdot \frac{1}{[g(x)]^2} \\
 &= \left[\lim_{k\to 0} \frac{f(x+k)g(x) - f(x)g(x)}{k} - \lim_{k\to 0}\frac{f(x)g(x+k) - f(x)g(x)}{k} \right] \cdot \frac{1}{[g(x)]^2} \\
 &= \left[\lim_{k\to 0} \frac{f(x+k) - f(x)}{k} \cdot g(x) - f(x) \cdot \lim_{k\to 0}\frac{g(x+k) - g(x)}{k} \right] \cdot \frac{1}{[g(x)]^2} \\
 &= \frac{f'(x)g(x) - f(x)g'(x)}{[g(x)]^2}.
 \end{aligned}
\]
The limit evaluation \(\lim_{k \to 0}\frac{1}{g(x+k)g(x)}=\frac{1}{[g(x)]^2}\) is justified by the differentiability of \(g(x)\), implying continuity, which can be expressed as \(\lim_{k \to 0}g(x+k) = g(x)\).

### Proof using implicit differentiation
Let \(h(x) = \frac{f(x)}{g(x)},\) so that \(f(x) = g(x)h(x).\)

The product rule then gives \(f'(x)=g'(x)h(x) + g(x)h'(x).\)

Solving for \(h'(x)\) and substituting back for \(h(x)\) gives:

\[
\begin{aligned}
 h'(x) &= \frac{f'(x) -g'(x)h(x)}{g(x)} \\
 &= \frac{f'(x) - g'(x)\cdot\frac{f(x)}{g(x)}}{g(x)} \\
 &= \frac{f'(x)g(x) - f(x)g'(x)}{[g(x)]^2}.
 \end{aligned}
\]

### Proof using the reciprocal rule or chain rule
Let \(h(x) = \frac{f(x)}{g(x)} = f(x) \cdot \frac{1}{g(x)}.\)

Then the product rule gives \(h'(x) = f'(x)\cdot\frac{1}{g(x)} + f(x) \cdot \frac{d}{dx}\left[\frac{1}{g(x)}\right].\)

To evaluate the derivative in the second term, apply the reciprocal rule, or the power rule along with the chain rule:

\[
\frac{d}{dx}\left[\frac{1}{g(x)}\right] = -\frac{1}{g(x)^2} \cdot g'(x) = \frac{-g'(x)}{g(x)^2}.
\]

Substituting the result into the expression gives
\[
\begin{aligned}
 h'(x) &= f'(x)\cdot\frac{1}{g(x)} + f(x)\cdot\left[\frac{-g'(x)}{g(x)^2}\right] \\

 &= \frac{f'(x)}{g(x)} - \frac{f(x)g'(x)}{g(x)^2} \\

 &= \frac{g(x)}{g(x)}\cdot\frac{f'(x)}{g(x)} - \frac{f(x)g'(x)}{g(x)^2} \\

 &= \frac{f'(x)g(x) - f(x)g'(x)}{g(x)^2}.
 \end{aligned}
\]

### Proof by logarithmic differentiation
Let \(h(x)=\frac{f(x)}{g(x)}.\) Taking the absolute value and natural logarithm of both sides of the equation gives

\[
\ln|h(x)|=\ln\left|\frac{f(x)}{g(x)}\right|
\]

Applying properties of the absolute value and logarithms,

\[
\ln|h(x)|=\ln|f(x)|-\ln|g(x)|
\]

Taking the logarithmic derivative of both sides,

\[
\frac{h'(x)}{h(x)}=\frac{f'(x)}{f(x)}-\frac{g'(x)}{g(x)}
\]

Solving for \(h'(x)\) and substituting back \(\tfrac{f(x)}{g(x)}\) for \(h(x)\) gives:

\[
\begin{aligned}
h'(x)&=h(x)\left[\frac{f'(x)}{f(x)}-\frac{g'(x)}{g(x)}\right]\\
&=\frac{f(x)}{g(x)}\left[\frac{f'(x)}{f(x)}-\frac{g'(x)}{g(x)}\right]\\
&=\frac{f'(x)}{g(x)}-\frac{f(x)g'(x)}{g(x)^2}\\
&=\frac{f'(x)g(x)-f(x)g'(x)}{g(x)^2}.
\end{aligned}
\]

Taking the absolute value of the functions is necessary for the logarithmic differentiation of functions that may have negative values, as logarithms are only real-valued for positive arguments. This works because \(\tfrac{d}{dx}(\ln|u|)=\tfrac{u'}{u}\), which justifies taking the absolute value of the functions for logarithmic differentiation.

## Higher order derivatives
Implicit differentiation can be used to compute the th derivative of a quotient (partially in terms of its first *n* − 1 derivatives). For example, differentiating \(f=gh\) twice (resulting in \(f* = g*h + 2g'h' + gh*\)) and then solving for \(h*\) yields
\[
h* = \left(\frac{f}{g}\right)* = \frac{f*-g*h-2g'h'}{g}.
\]
