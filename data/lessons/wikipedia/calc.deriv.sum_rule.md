> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Linearity_of_differentiation) — CC BY-SA 4.0

# Sum and difference rules

In calculus, the derivative of any linear combination of functions equals the same linear combination of the derivatives of the functions; this property is known as **linearity of differentiation**, the **rule of linearity**, or the superposition rule for differentiation. It is a fundamental property of the derivative that encapsulates in a single rule two simpler rules of differentiation, the sum rule (the derivative of the sum of two functions is the sum of the derivatives) and the constant factor rule (the derivative of a constant multiple of a function is the same constant multiple of the derivative). Thus it can be said that differentiation is linear, or the differential operator is a linear operator.

## Statement and derivation
Let *f* and *g* be functions, with *α* and *β* constants. Now consider

\(\frac{\mbox{d{\mbox{d} x} ( \alpha \cdot f(x) + \beta \cdot g(x) ).\)

By the sum rule in differentiation, this is

\(\frac{\mbox{d{\mbox{d} x} ( \alpha \cdot f(x) ) + \frac{\mbox{d{\mbox{d} x} (\beta \cdot g(x)),\)

and by the constant factor rule in differentiation, this reduces to

\(\alpha \cdot f'(x) + \beta \cdot g'(x).\)

Therefore,

\(\frac{\mbox{d{\mbox{d} x}(\alpha \cdot f(x) + \beta \cdot g(x)) = \alpha \cdot f'(x) + \beta \cdot g'(x).\)

Omitting the brackets, this is often written as:

\((\alpha \cdot f + \beta \cdot g)' = \alpha \cdot f'+ \beta \cdot g'.\)

## Detailed proofs/derivations from definition
We can prove the entire linearity principle at once, or, we can prove the individual steps (of constant factor and adding) individually. Here, both will be shown.

Proving linearity directly also proves the constant factor rule, the sum rule, and the difference rule as special cases. The sum rule is obtained by setting both constant coefficients to \(1\). The difference rule is obtained by setting the first constant coefficient to \(1\) and the second constant coefficient to \(-1\). The constant factor rule is obtained by setting either the second constant coefficient or the second function to \(0\). (From a technical standpoint, the domain of the second function must also be considered - one way to avoid issues is setting the second function equal to the first function and the second constant coefficient equal to \(0\). One could also define both the second constant coefficient and the second function to be 0, where the domain of the second function is a superset of the first function, among other possibilities.)

On the contrary, if we first prove the constant factor rule and the sum rule, we can prove linearity and the difference rule. Proving linearity is done by defining the first and second functions as being two other functions being multiplied by constant coefficients. Then, as shown in the derivation from the previous section, we can first use the sum law while differentiation, and then use the constant factor rule, which will reach our conclusion for linearity. In order to prove the difference rule, the second function can be redefined as another function multiplied by the constant coefficient of \(-1\). This would, when simplified, give us the difference rule for differentiation.

In the proofs/derivations below, the coefficients \(a, b\) are used; they correspond to the coefficients \(\alpha, \beta\) above.

### Linearity (directly)
Let \(a, b \in \mathbb{R}\). Let \(f, g\) be functions. Let \(j\) be a function, where \(j\) is defined only where \(f\) and \(g\) are both defined. (In other words, the domain of \(j\) is the intersection of the domains of \(f\) and \(g\).) Let \(x\) be in the domain of \(j\). Let \(j(x) = af(x) + bg(x)\).

We want to prove that \(j^{\prime}(x) = af^{\prime}(x) + bg^{\prime}(x)\).

By definition, we can see that

\[
\begin{align}
j^{\prime}(x) &= \lim_{h \rightarrow 0} \frac{j(x + h) - j(x)}{h} \\
&= \lim_{h \rightarrow 0} \frac{\left( af(x + h) + bg(x + h) \right) - \left( af(x) + bg(x) \right)}{h} \\
&= \lim_{h \rightarrow 0} \left( a\frac{f(x + h) - f(x)}{h} + b\frac{g(x + h) - g(x)}{h} \right) \\
\end{align}
\]

In order to use the limits law for the sum of limits, we need to know that \(\lim_{h \to 0} a\frac{f(x + h) - f(x)}{h}\) and \(\lim_{h \to 0} b\frac{g(x + h) - g(x)}{h}\) both individually exist. For these smaller limits, we need to know that \(\lim_{h \to 0} \frac{f(x + h) - f(x)}{h}\) and \(\lim_{h \to 0} \frac{g(x + h) - g(x)}{h}\) both individually exist to use the coefficient law for limits. By definition, \(f^{\prime}(x) = \lim_{h \to 0} \frac{f(x + h) - f(x)}{h}\) and \(g^{\prime}(x) = \lim_{h \to 0} \frac{g(x + h) - g(x)}{h}\). So, if we know that \(f^{\prime}(x)\) and \(g^{\prime}(x)\) both exist, we will know that \(\lim_{h \to 0} \frac{f(x + h) - f(x)}{h}\) and \(\lim_{h \to 0} \frac{g(x + h) - g(x)}{h}\) both individually exist. This allows us to use the coefficient law for limits to write

\[
\lim_{h \to 0} a\frac{f(x + h) - f(x)}{h}
= a\lim_{h \to 0}\frac{f(x + h) - f(x)}{h}
\]

and

\[
\lim_{h \to 0} b\frac{g(x + h) - g(x)}{h}
= b\lim_{h \to 0}\frac{g(x + h) - g(x)}{h}.
\]

With this, we can go back to apply the limit law for the sum of limits; since we know that \(\lim_{h \rightarrow 0} a\frac{f(x + h) - f(x)}{h}\) and \(\lim_{h \rightarrow 0} b\frac{g(x + h) - g(x)}{h}\) both individually exist. From here, we can directly go back to the derivative we were working on.
\[
\begin{align}
j^{\prime}(x) &= \lim_{h \rightarrow 0} \left( a\frac{f(x + h) - f(x)}{h} + b\frac{g(x + h) - g(x)}{h} \right) \\
&= \lim_{h \rightarrow 0} \left( a\frac{f(x + h) - f(x)}{h}\right) + \lim_{h \rightarrow 0} \left(b\frac{g(x + h) - g(x)}{h} \right) \\
&= a\lim_{h \rightarrow 0} \left( \frac{f(x + h) - f(x)}{h}\right) + b\lim_{h \rightarrow 0} \left(\frac{g(x + h) - g(x)}{h} \right) \\
&= af^{\prime}(x) + bg^{\prime}(x)
\end{align}
\]
Finally, we have shown what we claimed in the beginning: \(j^{\prime}(x) = af^{\prime}(x) + bg^{\prime}(x)\).

### Sum
Let \(f, g\) be functions. Let \(j\) be a function, where \(j\) is defined only where \(f\) and \(g\) are both defined.
(In other words, the domain of \(j\) is the intersection of the domains of \(f\) and \(g\).) Let \(x\) be in the domain of \(j\). Let \(j(x) = f(x) + g(x)\).

We want to prove that \(j^{\prime}(x) = f^{\prime}(x) + g^{\prime}(x)\).

By definition, we can see that

\[
\begin{align}
j^{\prime}(x) &= \lim_{h \rightarrow 0} \frac{j(x + h) - j(x)}{h} \\
&= \lim_{h \rightarrow 0} \frac{\left( f(x + h) + g(x + h) \right) - \left( f(x) + g(x) \right)}{h} \\
&= \lim_{h \rightarrow 0} \left( \frac{f(x + h) - f(x)}{h} + \frac{g(x + h) - g(x)}{h} \right) \\
\end{align}
\]
In order to use the law for the sum of limits here, we need to show that the individual limits, \(\lim_{h \rightarrow 0} \frac{f(x + h) - f(x)}{h}\) and \(\lim_{h \rightarrow 0} \frac{g(x + h) - g(x)}{h}\) both exist. By definition, \(f^{\prime}(x) = \lim_{h \rightarrow 0} \frac{f(x + h) - f(x)}{h}\)and \(g^{\prime}(x) = \lim_{h \rightarrow 0} \frac{g(x + h) - g(x)}{h}\), so the limits exist whenever the derivatives \(f^{\prime}(x)\) and \(g^{\prime}(x)\) exist. So, assuming that the derivatives exist, we can continue the above derivation

\[
\begin{align}
j^{\prime}(x) &= \lim_{h \rightarrow 0} \left( \frac{f(x + h) - f(x)}{h} + \frac{g(x + h) - g(x)}{h} \right) \\
&= \lim_{h \rightarrow 0} \frac{f(x + h) - f(x)}{h} + \lim_{h \rightarrow 0} \frac{g(x + h) - g(x)}{h} \\
&= f^{\prime}(x) + g^{\prime}(x)
\end{align}
\]

Thus, we have shown what we wanted to show, that: \(j^{\prime}(x) = f^{\prime}(x) + g^{\prime}(x)\).

### Difference
Let \(f, g\) be functions. Let \(j\) be a function, where \(j\) is defined only where \(f\) and \(g\) are both defined. (In other words, the domain of \(j\) is the intersection of the domains of \(f\) and \(g\).) Let \(x\) be in the domain of \(j\). Let \(j(x) = f(x) - g(x)\).

We want to prove that \(j^{\prime}(x) = f^{\prime}(x) - g^{\prime}(x)\).

By definition, we can see that:

\[
\begin{align}
j^{\prime}(x) &= \lim_{h \rightarrow 0} \frac{j(x + h) - j(x)}{h} \\
&= \lim_{h \rightarrow 0} \frac{\left( f(x + h) - (g(x + h) \right) - \left( f(x) - g(x) \right)}{h} \\
&= \lim_{h \rightarrow 0} \left( \frac{f(x + h) - f(x)}{h} - \frac{g(x + h) - g(x)}{h} \right) \\
\end{align}
\]

In order to use the law for the difference of limits here, we need to show that the individual limits, \(\lim_{h \rightarrow 0} \frac{f(x + h) - f(x)}{h}\) and \(\lim_{h \rightarrow 0} \frac{g(x + h) - g(x)}{h}\) both exist. By definition, \(f^{\prime}(x) = \lim_{h \rightarrow 0} \frac{f(x + h) - f(x)}{h}\) and that \(g^{\prime}(x) = \lim_{h \rightarrow 0} \frac{g(x + h) - g(x)}{h}\), so these limits exist whenever the derivatives \(f^{\prime}(x)\) and \(g^{\prime}(x)\) exist. So, assuming that the derivatives exist, we can continue the above derivation

\[
\begin{align}
j^{\prime}(x) &= \lim_{h \rightarrow 0} \left( \frac{f(x + h) - f(x)}{h} - \frac{g(x + h) - g(x)}{h} \right) \\
&= \lim_{h \rightarrow 0} \frac{f(x + h) - f(x)}{h} - \lim_{h \rightarrow 0} \frac{g(x + h) - g(x)}{h} \\
&= f^{\prime}(x) - g^{\prime}(x)
\end{align}
\]

Thus, we have shown what we wanted to show, that: \(j^{\prime}(x) = f^{\prime}(x) - g^{\prime}(x)\).

### Constant coefficient
Let \(f\) be a function. Let \(a \in \mathbb{R}\); \(a\) will be the constant coefficient. Let \(j\) be a function, where j is defined only where \(f\) is defined. (In other words, the domain of \(j\) is equal to the domain of \(f\).) Let \(x\) be in the domain of \(j\). Let \(j(x) = af(x)\).

We want to prove that \(j^{\prime}(x) = af^{\prime}(x)\).

By definition, we can see that:

\[
\begin{align}
j^{\prime}(x) &= \lim_{h \rightarrow 0} \frac{j(x + h) - j(x)}{h} \\
&= \lim_{h \rightarrow 0} \frac{af(x + h) - af(x)}{h} \\
&= \lim_{h \rightarrow 0} a\frac{f(x + h) - f(x)}{h} \\
\end{align}
\]

Now, in order to use a limit law for constant coefficients to show that

\[
\lim_{h \rightarrow 0} a\frac{f(x + h) - f(x)}{h} = a\lim_{h \rightarrow 0} \frac{f(x + h) - f(x)}{h}
\]
 we need to show that \(\lim_{h \rightarrow 0} \frac{f(x + h) - f(x)}{h}\) exists.
However, \(f^{\prime}(x) = \lim_{h \rightarrow 0} \frac{f(x + h) - f(x)}{h}\), by the definition of the derivative. So, if \(f^{\prime}(x)\) exists, then \(\lim_{h \rightarrow 0} \frac{f(x + h) - f(x)}{h}\) exists.

Thus, if we assume that \(f^{\prime}(x)\) exists, we can use the limit law and continue our proof.

\[
\begin{align}
j^{\prime}(x) &= \lim_{h \rightarrow 0} a\frac{f(x + h) - f(x)}{h} \\
&= a\lim_{h \rightarrow 0} \frac{f(x + h) - f(x)}{h} \\
&= af^{\prime}(x) \\
\end{align}
\]

Thus, we have proven that when \(j(x) = af(x)\), we have \(j^{\prime}(x) = af^{\prime}(x)\).
