> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Power_rule) — CC BY-SA 4.0

# Power rule for integration

In calculus, the **power rule** is used to differentiate functions of the form \(f(x) = x^r\), whenever \(r\) is a real number. Since differentiation is a linear operation on the space of differentiable functions, polynomials can also be differentiated using this rule. The power rule underlies the Taylor series as it relates a power series with a function's derivatives.

## Statement
Let \(f\) be a function satisfying \(f(x)=x^r\) for all \(x\), where \(r \in \mathbb{R}\). Then,
\(f'(x) = rx^{r-1} \, .\)
The power rule for integration states that
\(\int\! x^r \, dx=\frac{x^{r+1}}{r+1}+C\)
for any real number \(r \neq -1\). It can be derived by inverting the power rule for differentiation. In this equation C is any constant.

## Proofs
### Proof for real exponents
Let \(f(x) = x^r\), where \(r\) is any real number.

If \(f(x) = e^x\), then \(\ln (f(x)) = x\), where \(\ln\) is the natural logarithm function,
or \(f'(x) = f(x) = e^x\), as was required.
Therefore, applying the chain rule to {{nowrap|\(f(x) = e^{r\ln x}\),}} we see that

\[
f'(x)=\frac{r}{x} e^{r\ln x}= \frac{r}{x}x^r
\]

which simplifies to {{nowrap|\(rx^{r-1}\).}}

When \(x < 0\), we may use the same definition with \(x^r = ((-1)(-x))^r = (-1)^r(-x)^r\), where we now have \(-x > 0\). This necessarily leads to the same result. Note that because \((-1)^r\) does not have a conventional definition when \(r\) is not a rational number, irrational power functions are not well defined for negative bases. In addition, as rational powers of −1 with even denominators (in lowest terms) are not real numbers, these expressions are only real valued for rational powers with odd denominators (in lowest terms).

Finally, whenever the function is differentiable at \(x = 0\), the defining limit for the derivative is:

\[
\lim_{h\to 0} \frac{h^r - 0^r}{h}
\]

which yields 0 only when \(r\) is a rational number with odd denominator (in lowest terms) and \(r > 1\), and 1 when \(r = 1\). For all other values of \(r\), the expression \(h^r\) is not well-defined for \(h < 0\), as was covered above, or is not a real number, so the limit does not exist as a real-valued derivative. For the two cases that do exist, the values agree with the value of the existing power rule at 0, so no exception need be made.

The exclusion of the expression \(0^0\) (the case \(x = 0\)) from our scheme of exponentiation is due to the fact that the function \(f(x, y) = x^y\) has no limit at (0,0), since \(x^0\) approaches 1 as x approaches 0, while \(0^y\) approaches 0 as y approaches 0. Thus, it would be problematic to ascribe any particular value to it, as the value would contradict one of the two cases, dependent on the application. It is traditionally left undefined.

### Proofs for integer exponents
#### Proof by induction (natural numbers)
Let \(n\in\N\). It is required to prove that \(\frac{d}{dx} x^n = nx^{n-1}.\) The base case may be when \(n=0\) or \(1\), depending on how the set of natural numbers is defined.

When \(n=0\), \(\frac{d}{dx} x^0 = \frac{d}{dx} (1) = \lim_{h \to 0}\frac{1-1}{h} = \lim_{h \to 0}\frac{0}{h} = 0 = 0x^{0-1}.\)

When \(n=1\), \(\frac{d}{dx} x^1 = \lim_{h \to 0}\frac{(x+h)-x}{h} = \lim_{h \to 0}\frac{h}{h} = 1 = 1x^{1-1}.\)

Therefore, the base case holds either way.

Suppose the statement holds for some natural number *k*, i.e. \(\frac{d}{dx}x^k = kx^{k-1}.\)

When \(n=k+1\),
\[
\frac{d}{dx}x^{k+1} = \frac{d}{dx}(x^k \cdot x) = x^k \cdot \frac{d}{dx}x + x \cdot \frac{d}{dx}x^k = x^k + x \cdot kx^{k-1} = x^k + kx^k = (k+1)x^k = (k+1)x^{(k+1)-1}
\]
By the principle of mathematical induction, the statement is true for all natural numbers *n*.

#### Proof by binomial theorem (natural number)
Let \(y=x^n\), where \(n\in \mathbb{N}\).

Then,
\[
\begin{align}
\frac{dy}{dx}
&=\lim_{h\to 0}\frac{(x+h)^n-x^n}h\\[4pt]
&=\lim_{h\to 0}\frac{1}{h} \left[x^n+\binom n1 x^{n-1}h+\binom n2 x^{n-2}h^2+\dots+\binom nn h^n-x^n \right]\\[4pt]
&=\lim_{h\to 0}\left[\binom n 1 x^{n-1} + \binom n2 x^{n-2}h+ \dots+\binom nn h^{n-1}\right]\\[4pt]
&=nx^{n-1}
\end{align}
\]


Since n choose 1 is equal to n, and the rest of the terms all contain h, which is 0, the rest of the terms cancel. This proof only works for natural numbers as the binomial theorem only works for natural numbers.


#### Proof by Geometric Series (natural number)
Let \(y=x^n\), where \(n\in \mathbb{N}\)
and \(a,b\in \mathbb{R}\).
Then,
\[
\begin{align}

\left. \frac{dy}{dx} \right|_{x=a}
&=\lim_{b\to a}\frac{b^{n}-a^{n}}{b-a}\\[4pt]
&=\lim_{b\to a} \left[a^{n-1}+a^{n-2}b+ a^{n-3}b^2+\dots+b^{n-1}\right]\\[4pt]
&=na^{n-1}
\end{align}
\]


#### Generalization to negative integer exponents
For a negative integer *n*, let \(n=-m\) so that *m* is a positive integer.
Using the reciprocal rule,
\[
\frac{d}{dx}x^n = \frac{d}{dx} \left(\frac{1}{x^m}\right) = \frac{-\frac{d}{dx}x^m}{(x^m)^2} = -\frac{mx^{m-1}}{x^{2m}} = -mx^{-m-1} = nx^{n-1}.
\]
In conclusion, for any integer \(n\), \(\frac{d}{dx}x^n = nx^{n-1}.\)

### Generalization to rational exponents
Upon proving that the power rule holds for integer exponents, the rule can be extended to rational exponents.
#### Proof by chain rule
This proof is composed of two steps that involve the use of the chain rule for differentiation.
# Let \(y=x^r=x^\frac1n\), where \(n\in\N^+\). Then \(y^n=x\). By the chain rule, \(ny^{n-1}\cdot\frac{dy}{dx}=1\). Solving for \(\frac{dy}{dx}\),
\[
\frac{dy}{dx}
=\frac{1}{ny^{n-1}}
=\frac{1}{n\left(x^\frac1n\right)^{n-1}}
=\frac{1}{nx^{1-\frac1n}}
=\frac{1}{n}x^{\frac1n-1}
=rx^{r-1}
\]
Thus, the power rule applies for rational exponents of the form \(1/n\), where \(n\) is a nonzero natural number. This can be generalized to rational exponents of the form \(p/q\) by applying the power rule for integer exponents using the chain rule, as shown in the next step.
# Let \(y=x^r=x^{p/q}\), where \(p\in\Z, q\in\N^+,\) so that \(r\in\Q\). By the chain rule,
\[
\frac{dy}{dx}
=\frac{d}{dx}\left(x^\frac1q\right)^p
=p\left(x^\frac1q\right)^{p-1}\cdot\frac{1}{q}x^{\frac1q-1}
=\frac{p}{q}x^{p/q-1}=rx^{r-1}
\]

From the above results, we can conclude that when \(r\) is a rational number, \(\frac{d}{dx} x^r=rx^{r-1}.\)

#### Proof by implicit differentiation
A more straightforward generalization of the power rule to rational exponents makes use of implicit differentiation.

Let \(y=x^r=x^{p/q}\), where \(p, q \in \mathbb{Z}\) so that \(r \in \mathbb{Q}\).

Then,
\[
y^q=x^p
\]
Differentiating both sides of the equation with respect to \(x\),
\[
qy^{q-1}\cdot\frac{dy}{dx} = px^{p-1}
\]
Solving for \(\frac{dy}{dx}\),
\[
\frac{dy}{dx} = \frac{px^{p-1}}{qy^{q-1}}.
\]
Since \(y=x^{p/q}\),
\[
\frac d{dx}x^{p/q} = \frac{px^{p-1}}{qx^{p-p/q}}.
\]
Applying laws of exponents,
\[
\frac d{dx}x^{p/q} = \frac{p}{q}x^{p-1}x^{-p+p/q} = \frac{p}{q}x^{p/q-1}.
\]
Thus, letting \(r=\frac{p}{q}\), we can conclude that \(\frac d{dx}x^r = rx^{r-1}\) when \(r\) is a rational number.


#### Proof by Geometric Series (Positive Rationals)
This proof uses specific algebraic formulas to evaluate the definitional form of the derivative exactly for positive rational powers.

If \(y=x^{\frac{p}{q}}\), where \(p,q\in \mathbb{N}\)
and \(a,b\in \mathbb{R}\).
Then,
\[
\begin{align}

\left. \frac{dy}{dx} \right|_{x=a}

&=\lim_{b\to a}\frac{b^{\frac{p}{q}}-a^{\frac{p}{q}}}{b-a}\\[4pt]


\end{align}
\]


Now, consider the geometric sum formula, \(\frac{b^n-a^n}{b-a} = \displaystyle \sum_{i=0}^{n-1} b^{(n-1)-i}a^i\)

Observing the formula, one can show the following:

\[
\begin{align}

(b-a) = (b^{\frac{1}{n}})^{n}-(a^{\frac{1}{n}})^{n}=(b^{\frac{1}{n}}-a^{\frac{1}{n}})\displaystyle \sum_{i=0}^{n-1} (b^{\frac{1}{n}})^{(n-1)-i}(a^{\frac{1}{n}})^i

\end{align}
\]


Thus,


\[
\begin{align}

\frac{(b^{ \frac{1}{n} }-a^{ \frac{1}{n} })}{b-a}= \frac{1}{\displaystyle \sum_{i=0}^{n-1} (b^{\frac{1}{n}})^{(n-1)-i}(a^{\frac{1}{n}})^i}

= \frac{1}{ b^{1-\frac{1}{n}}+b^{1-\frac{2}{n}}a^{\frac{1}{n}}+b^{1-\frac{3}{n}}a^{\frac{2}{n}}+\dots+a^{1-\frac{1}{n}}}

\end{align}
\]


Returning to the first statement, it can now be evaluated in the following steps:

\[
\begin{align}

\lim_{b\to a}\frac{b^{\frac{p}{q}}-a^{\frac{p}{q}}}{b-a}

&= \lim_{b\to a}\frac{(b^{\frac{1}{q}})^{p}-(a^{\frac{1}{q}})^{p}}{b-a}\\[4pt]

&=\lim_{b\to a}\frac{(b^{\frac{1}{q}}-a^{\frac{1}{q}})(\displaystyle \sum_{i=0}^{p-1} (b^{\frac{1}{q}})^{(p-1)-i}(a^{\frac{1}{q}})^i )}{b-a}\\[4pt]

&=\lim_{b\to a}\frac{(b^{\frac{1}{q}}-a^{\frac{1}{q}})}{b-a}(\displaystyle \sum_{i=0}^{p-1} (b^{\frac{1}{q}})^{(p-1)-i}(a^{\frac{1}{q}})^i )\\[4pt]

&=\lim_{b\to a}\frac{1}{\displaystyle \sum_{i=0}^{q-1} (b^{\frac{1}{q}})^{(q-1)-i}(a^{\frac{1}{q}})^i}(\displaystyle \sum_{i=0}^{p-1}(b^{\frac{1}{q}})^{(p-1)-i}(a^{\frac{1}{q}})^i )\\[4pt]

&=(\frac{1}{qa^{1-\frac{1}{q}}})(p(a^{\frac{1}{q}})^{p-1})\\[4pt]

&=\frac{p}{q}a^{\frac{p}{q}-1}\\[4pt]

\end{align}
\]


#### Proof by Geometric Series (Negative Rationals)
If \(y=x^{\frac{p}{q}}\), where \(p,q\in \mathbb{Z/N}\)
and \(a,b\in \mathbb{R}\).

Then,
\[
\begin{align}

\left. \frac{dy}{dx} \right|_{x=a}

&=\lim_{b\to a}\frac{(\frac{1}{b^{\frac{p}{q}}})-(\frac{1}{a^{\frac{p}{q}}})}{b-a}\\[4pt]

&=\lim_{b\to a}(\frac{-1}{(ba)^{\frac{p}{q}}})(\frac{{b^{\frac{p}{q}}-a^{\frac{p}{q}}}}{b-a})\\[4pt]

&= (\frac{-1}{(a^2)^{\frac{p}{q}}})(\frac{p}{q}a^{\frac{p}{q}-1})\\[4pt]

&= -\frac{p}{q}a^{-\frac{p}{q}-1}

\end{align}
\]


## History
The power rule for integrals was first demonstrated in a geometric form by Italian mathematician Bonaventura Cavalieri in the early 17th century for all positive integer values of \({\displaystyle n}\), and during the mid 17th century for all rational powers by the mathematicians Pierre de Fermat, Evangelista Torricelli, Gilles de Roberval, John Wallis, and Blaise Pascal, each working independently. At the time, they were treatises on determining the area between the graph of a rational power function and the horizontal axis. With hindsight, however, it is considered the first general theorem of calculus to be discovered. The power rule for differentiation was derived by Isaac Newton and Gottfried Wilhelm Leibniz, each independently, for rational power functions in the mid 17th century, who both then used it to derive the power rule for integrals as the inverse operation. This mirrors the conventional way the related theorems are presented in modern basic calculus textbooks, where differentiation rules usually precede integration rules.

Although both men stated that their rules, demonstrated only for rational quantities, worked for all real powers, neither sought a proof of such, as at the time the applications of the theory were not concerned with such exotic power functions, and questions of convergence of infinite series were still ambiguous.

The unique case of \(r = -1\) was resolved by Flemish Jesuit and mathematician Grégoire de Saint-Vincent and his student Alphonse Antonio de Sarasa in the mid 17th century, who demonstrated that the associated definite integral,
\(\int_1^x \frac{1}{t}\, dt\)
representing the area between the rectangular hyperbola \(xy = 1\) and the x-axis, was a logarithmic function, whose base was eventually discovered to be the transcendental number e. The modern notation for the value of this definite integral is \(\ln(x)\), the natural logarithm.

## In complex analysis
If we consider functions of the form \(f(z) = z^c\) where \(c\) is any complex number and \(z\) is a complex number in a slit complex plane that excludes the branch point of 0 and any branch cut connected to it, and we use the conventional multivalued definition \(z^c := \exp(c\ln z)\), then it is straightforward to show that, on each branch of the complex logarithm, the same argument used above yields a similar result: \(f'(z) = \frac{c}{z}\exp(c\ln z)\).

In addition, if \(c\) is a positive integer, then there is no need for a branch cut: one may define \(f(0) = 0\), or define positive integral complex powers through complex multiplication, and show that \(f'(z) = cz^{c-1}\) for all complex \(z\), from the definition of the derivative and the binomial theorem.

However, due to the multivalued nature of complex power functions for non-integer exponents, one must be careful to specify the branch of the complex logarithm being used. In addition, no matter which branch is used, if \(c\) is not a positive integer, then the function is not differentiable at 0.

## See also
* Differentiation rules
* General Leibniz rule
* Inverse functions and differentiation
* Linearity of differentiation
* Product rule
* Quotient rule
* Table of derivatives
* Vector calculus identities

## References
### Notes


### Citations


## Further reading
*  Larson, Ron; Hostetler, Robert P.; and Edwards, Bruce H. (2003). *Calculus of a Single Variable: Early Transcendental Functions* (3rd edition). Houghton Mifflin Company. .

