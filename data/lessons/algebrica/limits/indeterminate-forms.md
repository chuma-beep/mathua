> Content sourced from [Algebrica](https://algebrica.org/indeterminate-forms/) — CC BY-NC 4.0

## What are indeterminate forms?

When working with limits, we will often come across expressions like these:

\\[\frac{0}{0} \qquad \frac{\infty}{\infty} \qquad 0 \cdot \infty \qquad \infty - \infty \\]

\\[1^{\infty} \qquad 0^{0} \qquad \infty^{0} \\]

These are called **indeterminate forms** because the value of the limit cannot be determined just by looking at them. Different [functions](<../functions/>) can produce the same symbolic form yet lead to completely different results. In this section, we will see what these forms mean, why they are called indeterminate, and how to work with them.

In simple terms, an indeterminate form does not stand for a specific number. It indicates that direct substitution is not sufficient to find the limit. For example, suppose we have two functions, \\(f(x)\\) and \\(g(x)\\), both of which tend to zero as \\(x\\) approaches some value \\(a\\):

\\[\lim_{x \to a} f(x) = 0 \qquad \text{and} \qquad \lim_{x \to a} g(x) = 0 \\]

If we try to find the limit of their ratio, direct substitution leads to:

\\[\lim_{x \to a} \frac{f(x)}{g(x)} \to \frac{0}{0} \\]

This is where the problem begins. That notation alone does not tell us how the ratio will behave. Depending on how \\(f(x)\\) and \\(g(x)\\) approach zero, the limit could be any finite number, zero, infinity, or it might not exist at all. The indeterminate form does not resolve the problem, it signals that further analysis is required.


Prior to analyzing each form, three key considerations should be noted:

  * An indeterminate form does not represent a specific number.
  * It does not signify an equality.
  * It serves as a label describing the limiting behaviour of the involved expressions.


## The form \\( \frac{0}{0} \\)

This is the indeterminate form you will encounter most often. It usually shows up when both the numerator and denominator become zero at the same value. Often, you can simplify the expression using algebra before finding the limit. For example:

\\[\lim_{x \to 2} \frac{x^2 - 4}{x - 2} \\]

Direct substitution gives \\(\dfrac{0}{0}\\). Factoring the numerator, we obtain:

\\[x^2 - 4 = (x - 2)(x + 2) \\]

so for \\(x \neq 2\\):

\\[\frac{x^2 - 4}{x - 2} = x + 2 \\]

Thus we obtain:

\\[\lim_{x \to 2} \frac{x^2 - 4}{x - 2} = 4 \\]

###### Once we simplify and see the underlying structure, the indeterminate form goes away. In more advanced problems, we can use [L’Hôpital’s Rule](<../hopital-rule/>) to resolve this type of limit, as long as the conditions for the rule are met.

## The form \\( \dfrac{\infty}{\infty} \\)

This form arises when both numerator and denominator diverge. For example, consider the following limit:

\\[\lim_{x \to \infty} \frac{3x^2 + 1}{2x^2 - 5} \\]

Direct substitution gives \\(\dfrac{\infty}{\infty}\\). To resolve this, we factor out the highest power of \\(x\\):

\\[\begin{align} \frac{3x^2 + 1}{2x^2 - 5} &= \frac{x^2(3 + 1/x^2)}{x^2(2 - 5/x^2)} \\\\[6pt] &= \frac{3 + 1/x^2}{2 - 5/x^2} \end{align} \\]

As \\(x \to \infty\\), the terms \\(1/x^2\\) and \\(5/x^2\\) vanish, so the limit becomes:

\\[\lim_{x \to \infty} \frac{3x^2 + 1}{2x^2 - 5} = \frac{3}{2} \\]

###### Once again, just looking at the symbols is not enough. The way the terms behave for large values of \\(x\\) determines the answer.

## The form \\( 0 \cdot \infty \\)

This form arises when one factor approaches zero while the other increases without bound. For example:

\\[\lim_{x \to 0^+} x \ln x \\]

As \\(x\\) approaches zero from the right, \\(\ln x\\) becomes increasingly negative, so the product takes the form \\(0 \cdot (-\infty)\\). To resolve this, we rewrite the expression as a ratio:

\\[x \ln x = \frac{\ln x}{1/x} \\]

Now, as \\(x \to 0^+\\) we have that:

\\[\ln x \to -\infty \qquad \frac{1}{x} \to +\infty \\]

so we obtain the form \\(\dfrac{-\infty}{+\infty}\\), which can be resolved using asymptotic comparison or [L’Hôpital’s Rule.](<../hopital-rule/>)

The result is:

\\[\lim_{x \to 0^+} x \ln x = 0 \\]

###### The key step is to transform the product into a quotient, which brings the problem back to a more familiar form.

## The form \\( \infty - \infty \\)

This form can hide how divergent quantities cancel each other out. For example, consider the following limit:

\\[\lim_{x \to \infty} \left( \sqrt{x^2 + x} - x \right) \\]

Both terms grow without bound, which gives the form \\(\infty - \infty\\). To resolve this, we rationalise the expression by multiplying and dividing by the conjugate:

\\[\begin{align} \sqrt{x^2 + x} - x &= \frac{\left(\sqrt{x^2 + x} - x\right)\left(\sqrt{x^2 + x} + x\right)}{\sqrt{x^2 + x} + x} \\\\[6pt] &= \frac{x^2 + x - x^2}{\sqrt{x^2 + x} + x} \\\\[6pt] &= \frac{x}{\sqrt{x^2 + x} + x} \end{align} \\]

Next, we factor \\(x\\) from the denominator:

\\[\begin{align} \frac{x}{\sqrt{x^2 + x} + x} &= \frac{x}{x\left(\sqrt{1 + 1/x} + 1\right)} \\\ &= \frac{1}{\sqrt{1 + 1/x} + 1} \end{align} \\]

Taking the limit as \\(x \to \infty\\):

\\[\lim_{x \to \infty} \frac{1}{\sqrt{1 + 1/x} + 1} = \frac{1}{\sqrt{1} + 1} = \frac{1}{2} \\]

###### After simplifying, the apparent divergence disappears and the limit turns out to be finite.

## Rationale for using the conjugate

In the analysis of the form \\(\infty - \infty\\), the conjugate is employed — an expression identical to the original except that the sign between its terms is reversed. The rationale for this approach is examined below. Consider the limit:

\\[\lim_{x \to \infty} \left(\sqrt{x^2 + 2x + 3} - x\right) \\]

Both terms increase without bound; thus, direct substitution results in the indeterminate form \\(\infty - \infty\\). The primary issue is not divergence itself, but rather that the two quantities grow at the same leading order. Their subtraction conceals a potential cancellation of dominant terms. If the expression remains unaltered, it is not possible to determine what remains after the dominant components offset each other. Multiplying and dividing by the conjugate makes this cancellation explicit:

\\[\begin{align} \frac{\left(\sqrt{x^2 + 2x + 3} - x\right)\left(\sqrt{x^2 + 2x + 3} + x\right)}{\sqrt{x^2 + 2x + 3} + x} &= \frac{x^2 + 2x + 3 - x^2}{\sqrt{x^2 + 2x + 3} + x} \\\\[6pt] &= \frac{2x + 3}{\sqrt{x^2 + 2x + 3} + x} \end{align} \\]

The original difference is now expressed as a quotient. The leading term \\(x^2\\) is eliminated, and the remaining expression can be further analysed by factoring \\(x\\) from both numerator and denominator:

\\[\begin{align} \frac{2x + 3}{\sqrt{x^2 + 2x + 3} + x} &= \frac{x\left(2 + 3/x\right)}{x\left(\sqrt{1 + 2/x + 3/x^2} + 1\right)} \\\\[6pt] &= \frac{2 + 3/x}{\sqrt{1 + 2/x + 3/x^2} + 1} \end{align} \\]

As \\(x \to \infty\\), the terms \\(3/x\\), \\(2/x\\) and \\(3/x^2\\) all approach zero, and the limit simplifies to:

\\[\lim_{x \to \infty} \left(\sqrt{x^2 + 2x + 3} - x\right) = \frac{2}{1 + 1} = 1 \\]

###### The use of the conjugate exposes the cancellation of dominant terms and reveals the true order of growth of the expression. The leading \\(x^2\\) terms cancel, and the behaviour is governed by lower-order terms. This perspective connects naturally with [Big O](<../big-o-notation/>) and [little-o](<../little-o-notation/>) notation, where limits are understood by comparing relative growth rates.

## Exponential indeterminate forms

When working with limits that involve exponential expressions, three additional indeterminate forms can appear:

\\[1^{\infty} \qquad 0^{0} \qquad \infty^{0} \\]

These arise in expressions of the form:

\\[\lim_{x \to a} F(x)^{G(x)} \\]

where base and exponent behave in conflicting ways. As with the other indeterminate forms, the symbolic appearance alone is not enough to determine the outcome. To handle these cases, we use the logarithmic transformation. Taking the [logarithm](<../logarithms/>) converts the exponential structure into a product, which can then be rewritten as a ratio and resolved using methods such as L’Hôpital’s Rule. The procedure is as follows. Suppose we want to compute:

\\[L = \lim_{x \to a} F(x)^{G(x)} \\]

Taking the natural logarithm of both sides we obtain:

\\[\ln L = \lim_{x \to a} G(x) \ln F(x) \\]

This reduces the problem to a limit of the form \\(0 \cdot \infty\\) or \\(\dfrac{0}{0}\\), both of which we already know how to handle. Once \\(\ln L\\) is found, we recover the original limit as \\(L = e^{\ln L}\\). As a concrete example, consider:

\\[\lim_{x \to 0} (1 + x)^{1/x} \\]

This has the form \\(1^{\infty}\\). Setting \\(L\\) equal to the limit and taking the logarithm:

\\[\ln L = \lim_{x \to 0} \frac{\ln(1 + x)}{x} \\]

This is now a \\(\dfrac{0}{0}\\) form. Applying L’Hôpital’s Rule, we differentiate numerator and denominator separately:

\\[\lim_{x \to 0} \frac{\ln(1 + x)}{x} = \lim_{x \to 0} \frac{1/(1+x)}{1} = 1 \\]

Therefore we obtain:

\\[L = e^{1} = e \\]

###### The exponential form contains a hidden quotient structure. Applying a logarithmic transformation makes this structure easier to see.

All seven indeterminate forms share a fundamental limitation: the symbolic expression alone does not determine the value of the limit. The critical factor is always the relative rate at which the quantities involved grow or diminish. Recognising the indeterminate form is only the initial step, determining the limit requires a closer analysis of the functions involved.

## Structural reduction of indeterminate forms

|   
---|---  
\\[\dfrac{0}{0}\\]| Factor and simplify, or apply [L’Hôpital’s Rule](<../hopital-rule/>)  
\\[\dfrac{\infty}{\infty}\\]| Factor out the dominant term, or apply L’Hôpital’s Rule  
\\[0 \cdot \infty\\]| Rewrite as \\(\dfrac{0}{1/\infty}\\) or \\(\dfrac{\infty}{1/0}\\) to obtain \\(\dfrac{0}{0}\\) or \\(\dfrac{\infty}{\infty}\\)  
\\[\infty - \infty\\]| Multiply by the conjugate, or find a common denominator  
\\[1^{\infty}\\]| Take the logarithm and reduce to \\(0 \cdot \infty\\)  
\\[0^{0}\\]| Take the logarithm and reduce to \\(0 \cdot \infty\\)  
\\[\infty^{0}\\]| Take the logarithm and reduce to \\(0 \cdot \infty\\)  
  
## Selected references

  * **MIT, D.J. Kleitman**. [Indeterminate Forms and L’Hôpital’s Rule](https://math.mit.edu/~djk/18_01/chapter26/contents.html)

  * **UC Davis, D. Kouba**. [Determining Limits Using L’Hôpital’s Rule](https://www.math.ucdavis.edu/~kouba/CalcOneDIRECTORY/lhopitaldirectory/LHopital.html)

  * **UC Davis, R. Marx**. [Indeterminate Forms and Limits](https://www.math.ucdavis.edu/~marx/Sec.%208.6.pdf)

  * **Harvard University, O. Knill**. [Infinity and Indeterminate Forms](https://people.math.harvard.edu/~knill/teaching/math1a2024/handouts/lecture10.pdf)

  * **University of Toronto – J. Campesato**. [Indeterminate Forms and L’Hôpital’s Rule](https://www.math.toronto.edu/campesat/ens/1819/lec21-1128.pdf)
