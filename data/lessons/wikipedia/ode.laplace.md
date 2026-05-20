> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Laplace_transform) — CC BY-SA 4.0

# Laplace transforms

In mathematics, the **Laplace transform**, named after Pierre-Simon Laplace (), is an integral transform that converts a function of a real variable (usually , in the *time domain*) to a function of a complex variable \(s\) (in the complex-valued frequency domain, also known as ***s*-domain** or ***s*-plane**).  The functions are often denoted using a lowercase symbol for the time-domain function and the corresponding uppercase symbol for the frequency-domain function, e.g. \(x(t)\) and .

The transform is useful for converting differentiation and integration in the time domain into the algebraic operations multiplication and division in the Laplace domain (analogous to how logarithms are useful for simplifying multiplication and division into addition and subtraction). This gives the transform many applications in science and engineering, mostly as a tool for solving linear differential equations and dynamical systems by replacing ordinary differential equations and integral equations with algebraic polynomial equations, and by replacing convolution with multiplication.

For example, through the Laplace transform, the equation of the simple harmonic oscillator (Hooke's law) \(x''(t) + kx(t)=0\) is converted into the algebraic equation \(s^2X(s) - s x(0) - x'(0) +kX(s) = 0,\) which incorporates the initial conditions \(x(0)\) and , and can be solved for the unknown function .  Once solved, the inverse Laplace transform can be used to transform it to the original domain. This is often aided by referencing tables such as that given below.

The Laplace transform is defined (for suitable functions ) by the integral

\[
\mathcal{L}\{f\}(s) = \int_0^\infty f(t) e^{-st} \, dt,
\]

where  is a complex number.

The Laplace transform is related to many other transforms. It is essentially the same as the Mellin transform and is closely related to the Fourier transform. Unlike for the Fourier transform, the Laplace transform of a function is often an analytic function, meaning that it can be expressed as a power series that converges locally, the coefficients of which represent the moments of the original function. Moreover, the techniques of complex analysis, especially contour integrals, can be used for simplifying calculations.

## History

The Laplace transform is named after mathematician and astronomer Pierre-Simon, Marquis de Laplace, who used a similar transform in his work on probability theory. Laplace wrote extensively about the use of generating functions  (1814), and the integral form of the Laplace transform evolved naturally as a result.

Laplace's use of generating functions was similar to what is now known as the z-transform, and he gave little attention to the continuous variable case which was discussed by Niels Henrik Abel.

From 1744, Leonhard Euler investigated integrals of the form
\(z = \int X(x) e^{ax}\, dx \quad\text{ and }\quad z = \int X(x) x^A \, dx\)
as solutions of differential equations, introducing in particular the gamma function. Joseph-Louis Lagrange was an admirer of Euler and, in his work on integrating probability density functions, investigated expressions of the form
\(\int X(x) e^{- a x } a^x\, dx,\)
which resembles a Laplace transform.

These types of integrals seem first to have attracted Laplace's attention in 1782, where he was following in the spirit of Euler in using the integrals themselves as solutions of equations. However, in 1785, Laplace took the critical step forward when, rather than simply looking for a solution in the form of an integral, he started to apply the transforms in the sense that was later to become popular. He used an integral of the form
\(\int x^s \varphi (x)\, dx,\)
akin to a Mellin transform, to transform the whole of a difference equation, in order to look for solutions of the transformed equation. He then went on to apply the Laplace transform in the same way and started to derive some of its properties, beginning to appreciate its potential power.

Laplace also recognised that Joseph Fourier's method of Fourier series for solving the diffusion equation could only apply to a limited region of space, because those solutions were periodic. In 1809, Laplace applied his transform to find solutions that diffused indefinitely in space.  In 1821, Cauchy developed an operational calculus for the Laplace transform that could be used to study linear differential equations in much the same way the transform is now used in basic engineering.  This method was popularized, and perhaps rediscovered, by Oliver Heaviside around the turn of the century.

Bernhard Riemann used the Laplace transform in his 1859 paper *On the number of primes less than a given magnitude*, in which he also developed the inversion theorem.  Riemann used the Laplace transform to develop the functional equation of the Riemann zeta function, and his method is still used to relate the modular transformation law of the Jacobi theta function, which is readily proved via Poisson summation, to the functional equation.

Hjalmar Mellin was among the first to study the Laplace transform rigorously in the Karl Weierstrass school of analysis, and apply it to the study of differential equations and special functions, at the turn of the 20th century.  At around the same time, Heaviside was busy with his operational calculus.  Thomas Joannes Stieltjes considered a generalization of the Laplace transform connected to his work on moments.  Other contributors in this time period included Mathias Lerch, Oliver Heaviside, and Thomas Bromwich.

In 1929, Vannevar Bush and Norbert Wiener published *Operational Circuit Analysis* as a text for engineering analysis of electrical circuits, applying both Fourier transforms and operational calculus, and in which they included one of the first predecessors of the modern table of Laplace transforms.
In 1934, Raymond Paley and Norbert Wiener published the important work *Fourier transforms in the complex domain*, about what is now called the Laplace transform (see below).   Also during the 30s, the Laplace transform was instrumental in Godfrey Harold Hardy and John Edensor Littlewood's study of tauberian theorems, and this application was later expounded on by , who developed other aspects of the theory such as a new method for inversion.  Edward Charles Titchmarsh wrote the influential *Introduction to the theory of the Fourier integral* (1937).

The current widespread use of the transform (mainly in engineering) came about during and soon after World War II, replacing the earlier Heaviside operational calculus. The advantages of the Laplace transform had been emphasized by Gustav Doetsch.

## Formal definition


The Laplace transform of a function *f*(*t*), defined for all real numbers *t* ≥ 0, is the function *F*(*s*) defined by
{{Equation box 1
|indent = :
|equation = \(F(s) = \int_0^\infty f(t)e^{-st} \, dt,\)
|ref = Eq. 1
|border
|border colour = #0073CF
|background colour=#F5FFFA}}
where \(s\) is a complex frequency-domain parameter \(s = \sigma + i \omega\) with real numbers  and .

An alternate notation for the Laplace transform is \(\mathcal{L}\{f\}\) instead of *F*. Thus \(F(s) = \mathcal L\{f\}(s)\) in functional notation. This is often written, especially in engineering settings, as {{tmath|1= F(s) = \mathcal{L}\{f(t)\} }}, with the understanding that the dummy variable \(t\) does not appear in the function .

The meaning of the integral depends on types of functions of interest. A necessary condition for existence of the integral is that  must be locally integrable on . For locally integrable functions that decay at infinity or are of exponential type ({{tmath| \vert f(t) \vert \le Ae^{B \vert t \vert} }}), the integral can be understood to be a (proper) Lebesgue integral. However, for many applications it is necessary to regard it as a conditionally convergent improper integral at ∞. Still more generally, the integral can be understood in a weak sense, and this is dealt with below.

One can define the Laplace transform of a finite Borel measure  by the Lebesgue integral
\(\mathcal{L}\{\mu\}(s) = \int_{[0,\infty)} e^{-st}\, d\mu(t).\)

An important special case is where  is a probability measure, for example, the Dirac delta function. In operational calculus, the Laplace transform of a measure is often treated as though the measure came from a probability density function . In that case, to avoid potential confusion, one often writes
\(\mathcal{L}\{f\}(s) = \int_{0^-}^\infty f(t)e^{-st} \, dt,\)
where the lower limit of 0<sup>−</sup> is shorthand notation for
\(\lim_{\varepsilon \to 0^+}\int_{-\varepsilon}^\infty.\)

This limit emphasizes that any point mass located at 0 is entirely captured by the Laplace transform. Although with the Lebesgue integral, it is not necessary to take such a limit, it does appear more naturally in connection with the Laplace–Stieltjes transform.

### Bilateral Laplace transform


When one says "the Laplace transform" without qualification, the unilateral or one-sided transform is usually intended. The Laplace transform can be alternatively defined as the *bilateral Laplace transform*, or two-sided Laplace transform, by extending the limits of integration to be the entire real axis. If that is done, the common unilateral transform becomes a special case of the bilateral transform, where the definition of the function being transformed includes being multiplied by the Heaviside step function.

The bilateral Laplace transform *F*(*s*) is defined as follows:
{{equation box 1
|indent = :
|equation = \(F(s) = \int_{-\infty}^\infty e^{-st} f(t)\, dt.\)
|ref = Eq. 2
|border
|border colour = #0073CF
|background colour=#F5FFFA}}
An alternate notation for the bilateral Laplace transform is {{tmath| \mathcal{B}\{f\} }}, instead of .

### Inverse Laplace transform


Two integrable functions have the same Laplace transform only if they differ on a set of Lebesgue measure zero. This means that, on the range of the transform, there is an inverse transform. In fact, besides integrable functions, the Laplace transform is a one-to-one mapping from one function space into another in many other function spaces as well, although there is usually no easy characterization of the range.

Typical function spaces in which this is true include the spaces of bounded continuous functions, the space *L*<sup>∞</sup>(0, ∞), or more generally tempered distributions on . The Laplace transform is also defined and injective for suitable spaces of tempered distributions.

In these cases, the image of the Laplace transform lives in a space of analytic functions in the region of convergence.  The inverse Laplace transform is given by the following complex integral, which is known by various names (the **Bromwich integral**, the **Fourier–Mellin integral**, and **Mellin's inverse formula**):
{{equation box 1
|indent = :
|equation = \(f(t) = \mathcal{L}^{-1}\{F\}(t) = \frac{1}{2 \pi i} \lim_{T\to\infty} \int_{\gamma - i T}^{\gamma + i T} e^{st} F(s)\, ds,\)
|ref = Eq. 3
|border
|border colour = #0073CF
|background colour=#F5FFFA}}
where  is a real number so that the contour path of integration is in the region of convergence of *F*(*s*). In most applications, the contour can be closed, allowing the use of the residue theorem. An alternative formula for the inverse Laplace transform is given by Post's inversion formula. The limit here is interpreted in the weak-* topology.

In practice, it is typically more convenient to decompose a Laplace transform into known transforms of functions obtained from a table and construct the inverse by inspection.

### Probability theory
In pure and applied probability, the Laplace transform is defined as an expected value. If  is a random variable with probability density function , then the Laplace transform of  is given by the expectation
\(\mathcal{L}\{f\}(s) = \operatorname{E}\left[e^{-sX}\right],\)
where \(\operatorname{E}[r]\) is the expectation of random variable .

By convention, this is referred to as the Laplace transform of the random variable  itself. Here, replacing  by −*t* gives the moment generating function of . The Laplace transform has applications throughout probability theory, including first passage times of stochastic processes such as Markov chains, and renewal theory.

Of particular use is the ability to recover the cumulative distribution function of a continuous random variable  by means of the Laplace transform as follows:
\(F_X(x) = \mathcal{L}^{-1}\left\{\frac{1}{s} \operatorname{E}\left[e^{-sX}\right]\right\}(x) = \mathcal{L}^{-1}\left\{\frac{1}{s} \mathcal{L}\{f\}(s)\right\}(x).\)

### Algebraic construction
The Laplace transform can be alternatively defined in a purely algebraic manner by applying a field of fractions construction to the convolution ring of functions on the positive half-line. The resulting space of abstract operators is exactly equivalent to Laplace space, but in this construction the forward and reverse transforms never need to be explicitly defined (avoiding the related difficulties with proving convergence).

## Region of convergence

If *f* is a locally integrable function (or more generally a Borel measure locally of bounded variation), then the Laplace transform *F*(*s*) of *f* converges provided that the limit
\(\lim_{R\to\infty}\int_0^R f(t)e^{-st}\,dt\)
exists.

The Laplace transform converges absolutely if the integral
\(\int_0^\infty \left|f(t)e^{-st}\right|\,dt\)
exists as a proper Lebesgue integral. The Laplace transform is usually understood as conditionally convergent, meaning that it converges in the former but not in the latter sense.

The set of values for which *F*(*s*) converges absolutely is either of the form Re(*s*) > *a* or Re(*s*) ≥ *a*, where *a* is an extended real constant with −∞ ≤ *a* ≤ ∞ (a consequence of the dominated convergence theorem). The constant *a* is known as the abscissa of absolute convergence, and depends on the growth behavior of *f*(*t*). Analogously, the two-sided transform converges absolutely in a strip of the form  *a* < Re(*s*) < *b*, and possibly including the lines 1=Re(*s*) = *a* or 1=Re(*s*) = *b*. The subset of values of *s* for which the Laplace transform converges absolutely is called the region of absolute convergence, or the domain of absolute convergence.  In the two-sided case, it is sometimes called the strip of absolute convergence. The Laplace transform is analytic in the region of absolute convergence: this is a consequence of Fubini's theorem and Morera's theorem.

Similarly, the set of values for which *F*(*s*) converges (conditionally or absolutely) is known as the region of conditional convergence, or simply the **region of convergence** (ROC). If the Laplace transform converges (conditionally) at 1=*s* = *s*<sub>0</sub>, then it automatically converges for all *s* with Re(*s*) > Re(*s*<sub>0</sub>). Therefore, the region of convergence is a half-plane of the form Re(*s*) > *a*, possibly including some points of the boundary line 1=Re(*s*) = *a*.

In the region of convergence Re(*s*) > Re(*s*<sub>0</sub>), the Laplace transform of *f* can be expressed by integrating by parts as the integral
\(F(s) = (s-s_0)\int_0^\infty e^{-(s-s_0)t}\beta(t)\,dt, \quad \beta(u) = \int_0^u e^{-s_0t}f(t)\,dt.\)

That is, *F*(*s*) can effectively be expressed, in the region of convergence, as the absolutely convergent Laplace transform of some other function.  In particular, it is analytic.  In its most general form, the Laplace transform gives a one-to-one correspondence between the holomorphic functions which, for some , are defined on \(\{s\in\mathbb{C}\ \vert \ \mathrm{Re}(s)>\sigma\}\) and are bounded there in absolute value by a polynomial, and the distributions on the real line supported on  which become tempered distributions after multiplied by \(e^{-\sigma t}\) for some .

There are several Paley–Wiener theorems concerning the relationship between the decay properties of *f*, and the properties of the Laplace transform within the region of convergence.

In engineering applications, a function corresponding to a linear time-invariant (LTI) system is *stable* if every bounded input produces a bounded output. This is equivalent to the absolute convergence of the Laplace transform of the impulse response function in the region Re(*s*) ≥ 0. As a result, LTI systems are stable, provided that the poles of the Laplace transform of the impulse response function have negative real part.

This ROC is used in knowing about the causality and stability of a system.

## Properties and theorems
The Laplace transform's key property is that it converts differentiation and integration in the time domain into multiplication and division by *s* in the Laplace domain. Thus, the Laplace variable *s* is also known as an *operator variable* in the Laplace domain: either the *derivative operator* or (for *s*<sup>−1</sup>) the *integration operator*.

Given the functions *f*(*t*) and *g*(*t*), and their respective Laplace transforms *F*(*s*) and *G*(*s*),
\(\begin{align}
f(t) &= \mathcal{L}^{-1}\{F(s)\},\\
g(t) &= \mathcal{L}^{-1}\{G(s)\},
\end{align}\)

the following table is a list of properties of unilateral Laplace transform:

{| class="wikitable" id="291017_tableid"
|+ Properties of the unilateral Laplace transform
|-
 ! scope="col" | Property
 ! scope="col" | Time domain
 ! scope="col" | *s* domain
 ! scope="col" | Comment
|-
 ! scope="row" | Linearity
 | \(a f(t) + b g(t) \\)
 | \(a F(s) + b G(s) \\)
 | Can be proved using basic rules of integration.
|-
 ! scope="row" | Frequency-domain derivative
 | \(t f(t) \\)
 | \(-F'(s) \\)
 | *F*′ is the first derivative of *F* with respect to *s*.
|-
 ! scope="row" | Frequency-domain general derivative
 | \(t^{n} f(t) \\)
 | \((-1)^{n} F^{(n)}(s) \\)
 | More general form, *n*th derivative of *F*(*s*).
|-
 ! scope="row" | Derivative
 | \(f'(t) \\)
 | \(s F(s) - f(0^{-}) \\)
 | *f* is assumed to be a differentiable function, and its derivative is assumed to be of exponential type.  This can then be obtained by integration by parts
|-
 ! scope="row" | Second derivative
 | \(f''(t) \\)
 | \(s^2 F(s) - s f(0^{-}) - f'(0^{-}) \\)
 | *f* is assumed twice differentiable and the second derivative to be of exponential type. Follows by applying the Differentiation property to *f*′(*t*).
|-
 ! scope="row" | General derivative
 | \(f^{(n)}(t)  \\)
 | \(s^n F(s) - \sum_{k=1}^{n} s^{n-k} f^{(k-1)}(0^{-}) \\)
 | *f* is assumed to be *n*-times differentiable, with *n*th derivative of exponential type.  Follows by mathematical induction.
|-
 ! scope="row" | Frequency-domain integration
 | \(\frac{1}{t}f(t)  \\)
 | \(\int_s^\infty F(\sigma)\, d\sigma \\)
 | This is deduced using the nature of frequency differentiation and conditional convergence.
|-
 ! scope="row" | Time-domain integration
 | \(\int_0^t f(\tau)\, d\tau = (u * f)(t)\)
 | \({1 \over s} F(s)\)
 | *u*(*t*) is the Heaviside step function and (*u* ∗ *f*)(*t*) is the convolution of *u*(*t*) and *f*(*t*).
|-
 ! scope="row" | Frequency shifting
 | \(e^{at} f(t)\)
 | \(F(s - a)\)
 |
|-
 ! scope="row" | Time shifting
 | \(f(t - a) u(t - a)\)
\(f(t) u(t - a) \\)
 | \(e^{-as} F(s) \\)
\(e^{-as} \mathcal{L}\{f(t + a)\}\)
 | *a* > 0, *u*(*t*) is the Heaviside step function
|-
 ! scope="row" | Time scaling
 | \(f(at)\)
 | \(\frac{1}{a} F \left ({s \over a} \right)\)
 | *a* > 0
|-
 ! scope="row" | Multiplication
 | \(f(t)g(t)\)
 | \(\frac{1}{2\pi i}\lim_{T\to\infty}\int_{c - iT}^{c + iT}F(\sigma)G(s - \sigma)\,d\sigma \\)
 | The integration is done along the vertical line 1=Re(*σ*) = *c* that lies entirely within the region of convergence of *F*.
|-
 ! scope="row" | Convolution
 | \((f * g)(t) = \int_{0}^{t} f(\tau)g(t - \tau)\,d\tau\)
 | \(F(s) \cdot G(s) \\)
 |
|-
 ! scope="row" | Circular convolution
 | \((f * g)(t) = \int_{0}^T f(\tau)g(t - \tau)\,d\tau\)
 | \(F(s) \cdot G(s) \\)
 | For periodic functions with period *T*.
|-
 ! scope="row" | Complex conjugation
 | \(f^*(t)\)
 | \(F^*(s^*)\)
 |
|-
 ! scope="row" | Periodic function
 | \(f(t)\)
 | \({1 \over 1 - e^{-Ts}} \int_0^T e^{-st} f(t)\,dt\)
 | *f*(*t*) is a periodic function of period *T* so that 1=*f*(*t*) = *f*(*t* + *T*), for all *t* ≥ 0. This is the result of the time shifting property and the geometric series.
|-
! scope="row" | Periodic summation
| \(f_P(t) = \sum_{n=0}^{\infty} f(t-Tn)\)
\(f_P(t) = \sum_{n=0}^{\infty} (-1)^n f(t-Tn)\)
| \(F_P(s) = \frac{1}{1-e^{-Ts}} F(s)\)
\(F_P(s) = \frac{1}{1+e^{-Ts}} F(s)\)
|
|}

; Initial value theorem :
\(f(0^+)=\lim_{s\to \infty}{sF(s)}.\)
; Final value theorem :
{{tmath|1= f(\infty)=\lim_{s\to 0}{sF(s)} }}, if all poles of \(sF(s)\) are in the left half-plane.
The final value theorem is useful because it gives the long-term behaviour without having to perform partial fraction decompositions (or other difficult algebra). If *F*(*s*) has a pole in the right-hand plane or poles on the imaginary axis (e.g., if \(f(t) = e^t\) or ), then the behaviour of this formula is undefined.

### Relation to power series
The Laplace transform can be viewed as a continuous analogue of a power series. If  *a*(*n*) is a discrete function of a positive integer *n*, then the power series associated to *a*(*n*) is the series
\(\sum_{n=0}^{\infty} a(n) x^n\)
where *x* is a real variable (see *Z-transform*). Replacing summation over *n* with integration over  *t*, a continuous version of the power series becomes
\(\int_{0}^{\infty} f(t) x^t\, dt\)
where the discrete function *a*(*n*) is replaced by the continuous one *f*(*t*).

Changing the base of the power from *x* to *e* gives
\(\int_{0}^{\infty} f(t) \left(e^{\ln{x}}\right)^t\, dt\)

For this to converge for, say, all bounded functions *f*, it is necessary to require that ln *x* < 0. Making the substitution 1=−*s* = ln *x* gives just the Laplace transform:
\(\int_{0}^{\infty} f(t) e^{-st}\, dt\)

In other words, the Laplace transform is a continuous analog of a power series, in which the discrete parameter *n* is replaced by the continuous parameter *t*, and *x* is replaced by *e*<sup>−*s*</sup>.

Analogously to a power series, if {{tmath|1= a(n)=O(\rho^{-n}) }}, then the power series converges to an analytic function in , if {{tmath|1= f(t)=O(e^{-\sigma t}) }}, the Laplace transform converges to an analytic function for .

### Relation to moments


The quantities
\(\mu_n = \int_0^\infty t^nf(t)\, dt\)
are the *moments* of the function *f*.  If the first *n* moments of *f* converge absolutely, then by repeated differentiation under the integral,
\((-1)^n(\mathcal L f)^{(n)}(0) = \mu_n .\)
This is of special significance in probability theory, where the moments of a random variable *X* are given by the expectation values {{tmath|1= \mu_n=\operatorname{E}[X^n] }}.  Then, the relation holds
\(\mu_n = (-1)^n\frac{d^n}{ds^n}\operatorname{E}\left[e^{-sX}\right](0).\)

### Transform of a function's derivative
It is often convenient to use the differentiation property of the Laplace transform to find the transform of a function's derivative. This can be derived from the basic expression for a Laplace transform as follows:
\(\begin{align}
  \mathcal{L} \left\{f(t)\right\} &= \int_{0^-}^\infty e^{-st} f(t)\, dt \\[6pt]
                                  &= \left[\frac{f(t)e^{-st}}{-s} \right]_{0^-}^\infty -
                                       \int_{0^-}^\infty \frac{e^{-st}}{-s} f'(t) \, dt\quad \text{(by parts)} \\[6pt]
                                  &= \left[-\frac{f(0^-)}{-s}\right] + \frac 1 s \mathcal{L} \left\{f'(t)\right\},
\end{align}\)
yielding
\(\mathcal{L} \{ f'(t) \} = s\cdot\mathcal{L} \{ f(t) \}-f(0^-),\)
and in the bilateral case,
\(\mathcal{L} \{ f'(t) \} = s \int_{-\infty}^\infty e^{-st} f(t)\,dt  = s \cdot \mathcal{L} \{ f(t) \}.\)

The general result
\(\mathcal{L} \left\{ f^{(n)}(t) \right\} = s^n \cdot \mathcal{L} \{ f(t) \} - s^{n - 1} f(0^-) - \cdots - f^{(n - 1)}(0^-),\)
where \(f^{(n)}\) denotes the *n*th derivative of *f*, can then be established with an inductive argument.

### Evaluating integrals over the positive real axis
A useful property of the Laplace transform is the following:
\(\int_0^\infty f(x)g(x)\,dx = \int_0^\infty(\mathcal{L} f)(s)\cdot(\mathcal{L}^{-1}g)(s)\,ds\)
under suitable assumptions on the behaviour of  and  in a right neighbourhood of \(0\) and on the decay rate of  and  in a left neighbourhood of . The above formula is a variation of integration by parts, with the operators
\(\frac{d}{dx}\) and \(\int \,dx\) being replaced by \(\mathcal{L}\) and {{tmath| \mathcal{L}^{-1} }}. Let us prove the equivalent formulation:
\(\int_0^\infty(\mathcal{L} f)(x)g(x)\,dx = \int_0^\infty f(s)(\mathcal{L}g)(s)\,ds.\)

By plugging in \((\mathcal{L}f)(x)=\int_0^\infty f(s)e^{-sx}\,ds\) the left-hand side turns into:
\(\int_0^\infty\int_0^\infty f(s)g(x) e^{-sx}\,ds\,dx,\)
but assuming Fubini's theorem holds, by reversing the order of integration we get the wanted right-hand side.

This method can be used to compute integrals that would otherwise be difficult to compute using elementary methods of real calculus. For example,
\(\int_0^\infty\frac{\sin x}{x}dx = \int_0^\infty \mathcal{L}(1)(x)\sin x dx = \int_0^\infty 1 \cdot \mathcal{L}(\sin)(x)dx = \int_0^\infty \frac{dx}{x^2 + 1} = \frac{\pi}{2}.\)

## Relationship to other transforms
### Laplace–Stieltjes transform
The (unilateral) Laplace–Stieltjes transform of a function *g* : ℝ → ℝ is defined by the Lebesgue–Stieltjes integral
\(\{ \mathcal{L}^*g \}(s) = \int_0^\infty e^{-st} \, d\,g(t) ~.\)

The function *g* is assumed to be of bounded variation.  If *g* is the antiderivative of *f*:
\(g(x) = \int_0^x f(t)\,d\,t\)

then the Laplace–Stieltjes transform of  and the Laplace transform of  coincide. In general, the Laplace–Stieltjes transform is the Laplace transform of the Stieltjes measure associated to . So in practice, the only distinction between the two transforms is that the Laplace transform is thought of as operating on the density function of the measure, whereas the Laplace–Stieltjes transform is thought of as operating on its cumulative distribution function.

### Fourier transform


Let \(f\) be a complex-valued Lebesgue integrable function supported on , and let \(F(s) = \mathcal Lf(s)\) be its Laplace transform.  Then, within the region of convergence, we have

\[
F(\sigma + i\tau) = \int_0^\infty f(t)e^{-\sigma t}e^{-i\tau t}\,dt,
\]

which is the Fourier transform of the function {{tmath| f(t)e^{-\sigma t} }}.

Indeed, the Fourier transform is a special case (under certain conditions) of the bilateral Laplace transform. The main difference is that the Fourier transform of a function is a complex function of a *real* variable (frequency ), the Laplace transform of a function is a complex function of a *complex* variable (damping factor \(\sigma\) and frequency ). The Laplace transform is usually restricted to transformation of functions of *t* with *t* ≥ 0. A consequence of this restriction is that the Laplace transform of a function is a holomorphic function of the variable *s*. Unlike the Fourier transform, the Laplace transform of a distribution is generally a well-behaved function. Techniques of complex variables can also be used to directly study Laplace transforms. As a holomorphic function, the Laplace transform has a power series representation. This power series expresses a function as a linear superposition of moments of the function. This perspective has applications in probability theory.

Formally, the Fourier transform is equivalent to evaluating the bilateral Laplace transform with imaginary argument 1=*s* = *iω* when the condition explained below is fulfilled,

\[
\begin{align}
  \hat{f}(\omega) &= \mathcal{F}\{f(t)\} \\[4pt]
                  &= \mathcal{L}\{f(t)\}|_{s = i \omega}  =  F(s)|_{s = i \omega} \\[4pt]
                  &= \int_{-\infty}^\infty e^{-i \omega t} f(t)\,dt~.
\end{align}
\]


This convention of the Fourier transform ( in **) requires a factor of  on the inverse Fourier transform. This relationship between the Laplace and Fourier transforms is often used to determine the frequency spectrum of a signal or dynamical system.

The above relation is valid as stated if and only if the region of convergence (ROC) of *F*(*s*) contains the imaginary axis, 1=*σ* = 0.

For example, the function 1=*f*(*t*) = cos(*ω*<sub>0</sub>*t*) has a Laplace transform 1=*F*(*s*) =  *s*/(*s*<sup>2</sup> + *ω*<sub>0</sub><sup>2</sup>) whose ROC is Re(*s*) > 0. As 1=*s* = *iω*<sub>0</sub> is a pole of  *F*(*s*), substituting 1=*s* = *iω* in *F*(*s*) does not yield the Fourier transform of *f*(*t*)*u*(*t*), which contains terms proportional to the Dirac delta functions *δ*(*ω* ± *ω*<sub>0</sub>).

However, a relation of the form

\[
\lim_{\sigma\to 0^+} F(\sigma+i\omega) = \hat{f}(\omega)
\]

holds under much weaker conditions. For instance, this holds for the above example provided that the limit is understood as a weak limit of measures (see vague topology). General conditions relating the limit of the Laplace transform of a function on the boundary to the Fourier transform take the form of Paley–Wiener theorems.

### Mellin transform


The Mellin transform and its inverse are related to the two-sided Laplace transform by a change of variables.

If in the Mellin transform
\(G(s) = \mathcal{M}\{g(\theta)\} = \int_0^\infty \theta^s g(\theta) \, \frac{d\theta} \theta\)
we set 1=*θ* = *e*<sup>−*t*</sup> we get a two-sided Laplace transform.

### Z-transform


The unilateral or one-sided Z-transform is the Laplace transform of an ideally sampled signal with the substitution of
\(z \stackrel{\mathrm{def} }{ {}={} } e^{sT} ,\)
where 1=*T* = 1/*f<sub>s</sub>* is the sampling interval (in units of time e.g., seconds) and  *f<sub>s</sub>* is the sampling rate (in samples per second or hertz).

Let
\(\Delta_T(t) \ \stackrel{\mathrm{def}}{=}\  \sum_{n=0}^{\infty}  \delta(t - n T)\)
be a sampling impulse train (also called a Dirac comb) and
\(\begin{align}
  x_q(t)   &\stackrel{\mathrm{def} }{ {}={} }  x(t) \Delta_T(t) = x(t) \sum_{n=0}^{\infty}  \delta(t - n T) \\
           &= \sum_{n=0}^{\infty} x(n T) \delta(t - n T) = \sum_{n=0}^{\infty} x[n] \delta(t - n T)
\end{align}\)
be the sampled representation of the continuous-time *x*(*t*)
\(x[n] \stackrel{\mathrm{def} }{ {}={} }  x(nT) ~.\)

The Laplace transform of the sampled signal *x*<sub>*q*</sub>(*t*)  is
\(\begin{align}
  X_q(s) &= \int_{0^-}^\infty x_q(t) e^{-s t} \,dt \\
         &= \int_{0^-}^\infty \sum_{n=0}^\infty x[n] \delta(t - n T) e^{-s t} \, dt \\
         &= \sum_{n=0}^\infty x[n] \int_{0^-}^\infty \delta(t - n T) e^{-s t} \, dt \\
         &= \sum_{n=0}^\infty x[n] e^{-n s T}~.
\end{align}\)

This is the precise definition of the unilateral Z-transform of the discrete function *x*[*n*]
\(X(z) = \sum_{n=0}^{\infty} x[n] z^{-n}\)
with the substitution of *z* → *e*<sup>*sT*</sup>.

Comparing the last two equations, we find the relationship between the unilateral Z-transform and the Laplace transform of the sampled signal,
\(X_q(s) = X(z) \Big|_{z=e^{sT}}.\)

The similarity between the Z- and Laplace transforms is expanded upon in the theory of time scale calculus.

### Borel transform
The integral form of the Borel transform
\(F(s) = \int_0^\infty f(z)e^{-sz}\, dz\)
is a special case of the Laplace transform for *f* an entire function of exponential type, meaning that
\(|f(z)|\le Ae^{B|z|}\)
for some constants *A* and *B*.  The generalized Borel transform allows a different weighting function to be used, rather than the exponential function, to transform functions not of exponential type. Nachbin's theorem gives necessary and sufficient conditions for the Borel transform to be well defined.

### Fundamental relationships
Since an ordinary Laplace transform can be written as a special case of a two-sided transform, and since the two-sided transform can be written as the sum of two one-sided transforms, the theory of the Laplace-, Fourier-, Mellin-, and Z-transforms are at bottom the same subject. However, a different point of view and different characteristic problems are associated with each of these four major integral transforms.

## Table of selected Laplace transforms


The following table provides Laplace transforms for many common functions of a single variable. For definitions and explanations, see the *Explanatory Notes* at the end of the table.

Because the Laplace transform is a linear operator,
* The Laplace transform of a sum is the sum of Laplace transforms of each term.\(\mathcal{L}\{f(t) + g(t)\}  = \mathcal{L}\{f(t)\} + \mathcal{L}\{ g(t)\}\)
* The Laplace transform of a multiple of a function is that multiple times the Laplace transformation of that function.\(\mathcal{L}\{a f(t)\}  = a \mathcal{L}\{ f(t)\}\)

Using this linearity, and various trigonometric, hyperbolic, and complex number (etc.) properties and/or identities, some Laplace transforms can be obtained from others more quickly than by using the definition directly.

The unilateral Laplace transform takes as input a function whose time domain is the non-negative reals, which is why all of the time domain functions in the table below are multiples of the Heaviside step function, *u*(*t*).

The entries of the table that involve a time delay *τ* are required to be causal (meaning that *τ* > 0).  A causal system is a system where the impulse response *h*(*t*) is zero for all time  prior to 1=*t* = 0. In general, the region of convergence for causal systems is not the same as that of anticausal systems.

{| class="wikitable" style="text-align: center;"
|+ Selected Laplace transforms
|-
 ! scope="col" | Function
 ! scope="col" | Time domain  \(f(t) = \mathcal{L}^{-1}\{F(s)\}\)
 ! scope="col" | Laplace s-domain  \(F(s) = \mathcal{L}\{f(t)\}\)
 ! scope="col" | Region of convergence
 ! scope="col" | Reference
|-
 ! scope="row" | unit impulse
 | \(\delta(t) \\)
 | \(1\)
 | all *s*
 | inspection
|-
 ! scope="row" | delayed impulse
 | \(\delta(t - \tau) \\)
 | \(e^{-\tau s} \\)
 | all *s*
 | time shift ofunit impulse
|-
 ! scope="row"| unit step
 | \(u(t) \\)
 | \({ 1 \over s }\)
 | \(\operatorname{Re}(s) > 0\)
 | integrate unit impulse
|-
 ! scope="row" | delayed unit step
 | \(u(t - \tau) \\)
 | \(\frac 1 s e^{-\tau s}\)
 | \(\operatorname{Re}(s) > 0\)
 | time shift ofunit step
|-
 ! scope="row" | product of delayed function and delayed step
 | \(f(t-\tau)u(t-\tau)\)
 | \(e^{-s\tau}\mathcal{L}\{f(t)\}\)
 |
 | u-substitution, \(u=t-\tau\)
|-
 !rectangular impulse
 | \(u (t) - u(t - \tau)\)
 | \(\frac{1}{s}(1 - e^{-\tau s})\)
 | \(\operatorname{Re}(s) > 0\)
 |
|-
 ! scope="row" | ramp
 | \(t \cdot u(t)\\)
 | \(\frac 1 {s^2}\)
 | \(\operatorname{Re}(s) > 0\)
 | integrate unitimpulse twice
|-
 ! scope="row" | *n*th power  (for integer *n*)
 | \(t^n \cdot u(t)\)
 | \({ n! \over s^{n + 1} }\)
 | \(\operatorname{Re}(s) > 0\)  (*n* > −1)
 | integrate unitstep *n* times
|-
 ! scope="row" | *q*th power  (for complex *q*)
 | \(t^q \cdot u(t)\)
 | \({ \operatorname{\Gamma}(q + 1) \over s^{q + 1} }\)
 | \(\operatorname{Re}(s) > 0\)    \(\operatorname{Re}(q) > -1\)
 |
|-
 ! scope="row" | *n*th root
 | \(\sqrt[n]{t} \cdot u(t)\)
 | \({ 1 \over s^{\frac 1 n + 1} } \operatorname{\Gamma}\left(\frac 1 n + 1\right)\)
 | \(\operatorname{Re}(s) > 0\)
 | Set *q*  above.
|-
 ! scope="row" | *n*th power with frequency shift
 | \(t^{n} e^{-\alpha t} \cdot u(t)\)
 | \(\frac{n!}{(s+\alpha)^{n+1}}\)
 | \(\operatorname{Re}(s) > -\alpha\)
 | Integrate unit step,apply frequency shift
|-
 ! scope="row" | delayed *n*th power  with frequency shift
 | \((t-\tau)^n e^{-\alpha (t-\tau)} \cdot u(t-\tau)\)
 | \(\frac{n! \cdot e^{-\tau s}}{(s+\alpha)^{n+1}}\)
 |  \(\operatorname{Re}(s) > -\alpha\)
 | integrate unit step,apply frequency shift,apply time shift
|-
 ! scope="row" | exponential decay
 | \(e^{-\alpha t} \cdot u(t)\)
 | \({ 1 \over s+\alpha }\)
 | \(\operatorname{Re}(s) > -\alpha\)
 | Frequency shift ofunit step
|-
 ! scope="row" | two-sided exponential decay (only for bilateral transform)
 | \(e^{-\alpha|t|}  \\)
 | \({ 2\alpha \over \alpha^2 - s^2 }\)
 | \(-\alpha < \operatorname{Re}(s) < \alpha\)
 | Frequency shift ofunit step
|-
 ! scope="row" | exponential approach
 | \((1-e^{-\alpha t})  \cdot u(t)  \\)
 | \(\frac{\alpha}{s(s+\alpha)}\)
 | \(\operatorname{Re}(s) > 0\)
 | unit step minusexponential decay
|-
 ! scope="row" | sine
 | \(\sin(\omega t) \cdot u(t) \\)
 | \({ \omega \over s^2 + \omega^2  }\)
 | \(\operatorname{Re}(s) > 0\)
 |
|-
 ! scope="row" | cosine
 | \(\cos(\omega t) \cdot u(t) \\)
 | \({ s \over s^2 + \omega^2  }\)
 | \(\operatorname{Re}(s) > 0\)
 |
|-
 ! scope="row" | hyperbolic sine
 | \(\sinh(\alpha t) \cdot u(t) \\)
 | \({ \alpha \over s^2 - \alpha^2 }\)
 | \(\operatorname{Re}(s) > \left| \alpha \right|\)
 |
|-
 ! scope="row" | hyperbolic cosine
 | \(\cosh(\alpha t) \cdot u(t) \\)
 | \({ s \over s^2 - \alpha^2  }\)
 |  \(\operatorname{Re}(s) > \left| \alpha \right|\)
 |
|-
 ! scope="row" | exponentially decaying  sine wave
 | \(e^{-\alpha t}  \sin(\omega t) \cdot u(t) \\)
 | \({ \omega \over (s+\alpha)^2 + \omega^2  }\)
 | \(\operatorname{Re}(s) > - \alpha\)
 |
|-
 ! scope="row" | exponentially decaying  cosine wave
 | \(e^{-\alpha t}  \cos(\omega t) \cdot u(t) \\)
 | \({ s+\alpha \over (s+\alpha)^2 + \omega^2  }\)
 | \(\operatorname{Re}(s) > - \alpha\)
 |
|-
 ! scope="row" | natural logarithm
 | \(\ln(t) \cdot u(t)\)
 | \(-{1 \over s} \left[ \ln(s)+\gamma \right]\)
 |  \(\operatorname{Re}(s) > 0\)
 |
|-
 ! scope="row" | Bessel function  of the first kind,  of order *n*
 | \(J_n(\omega t) \cdot u(t)\)
 | \(\frac{ \left(\sqrt{s^2+ \omega^2}-s\right)^{\!n}}{\omega^n \sqrt{s^2 + \omega^2}}\)
 | \(\operatorname{Re}(s) > 0\)  (*n* > −1)
 |
|-
 ! scope="row" | Error function
 | \(\operatorname{erf}(t) \cdot u(t)\)
 | \(\frac{1}{s} e^{s^2 / 4} \!\left(1 - \operatorname{erf} \frac{s}{2} \right)\)
 | \(\operatorname{Re}(s) > 0\)
 |
|-
 | colspan=5 style="text-align: left;" |**Explanatory notes:**


* *u*(*t*) represents the Heaviside step function.
* *δ*  represents the Dirac delta function.
* Γ(*z*) represents the gamma function.
* *γ* is the Euler–Mascheroni constant.

* *t*, a real number, typically represents *time*, although it can represent *any* independent dimension.
* *s* is the complex frequency domain parameter, and  Re(*s*) is its real part.
* *α*, *β*, *τ*, and *ω* are real numbers.
* *n* is an integer.

|}

## *s*-domain equivalent circuits and impedances
The Laplace transform is often used in circuit analysis by conversions to the *s*-domain of circuit elements. Circuit elements can be transformed into impedances, very similar to phasor impedances.

Here is a summary of equivalents:


Note that the resistor is exactly the same in the time domain and the *s*-domain. The sources are put in if there are initial conditions on the circuit elements. For example, if a capacitor has an initial voltage across it, or if the inductor has an initial current through it, the sources inserted in the *s*-domain account for that.

The equivalents for current and voltage sources are derived from the transformations in the table above.

## Examples and applications


The Laplace transform is used frequently in engineering and physics; the output of a linear time-invariant system can be calculated by convolving its unit impulse response with the input signal. Performing this calculation in Laplace space turns the convolution into a multiplication; the latter being easier to solve because of its algebraic form. For more information, see control theory. The Laplace transform is invertible on a large class of functions.  Given a mathematical or functional description of an input or output to a system, the Laplace transform provides an alternative functional description that often simplifies the process of analyzing the behavior of the system, or in synthesizing a new system based on a set of specifications.

The Laplace transform can also be used to solve differential equations and is used extensively in mechanical engineering and electrical engineering.  The Laplace transform reduces a linear differential equation to an algebraic equation, which can then be solved by the formal rules of algebra.  The original differential equation can then be solved by applying the inverse Laplace transform. English electrical engineer Oliver Heaviside first proposed a similar scheme, although without using the Laplace transform; and the resulting operational calculus is credited as the Heaviside calculus.

### Evaluating improper integrals
Let {{tmath|1= \mathcal{L}\left\{f(t)\right\} = F(s) }}. Then (see the table above)

\[
\partial_s\mathcal{L} \left\{\frac{f(t)} t \right\} = \partial_s\int_0^\infty \frac{f(t)}{t}e^{-st}\, dt = -\int_0^\infty f(t)e^{-st}dt = - F(s)
\]


From which one gets:
\(\mathcal{L} \left\{\frac{f(t)} t \right\} = \int_s^\infty F(p)\, dp.\)

In the limit , one gets
\(\int_0^\infty \frac{f(t)} t \, dt = \int_0^\infty F(p)\, dp,\)
provided that the interchange of limits can be justified. This is often possible as a consequence of the final value theorem. Even when the interchange cannot be justified the calculation can be suggestive. For example, with *a* ≠ 0 ≠ *b*, proceeding formally one has
\(\begin{align}
\int_0^\infty \frac{ \cos(at) - \cos(bt) }{t} \, dt
&=\int_0^\infty \left(\frac p {p^2 + a^2} - \frac{p}{p^2 + b^2}\right)\, dp \\[6pt]
&=\left[ \frac{1}{2}  \ln\frac{p^2 + a^2}{p^2 + b^2} \right]_0^\infty = \frac{1}{2} \ln \frac{b^2}{a^2} = \ln \left| \frac {b}{a} \right|.
\end{align}\)

### Complex impedance of a capacitor
In the theory of electrical circuits, the current flow in a capacitor is proportional to the capacitance and rate of change in the electrical potential (with equations as for the SI unit system). Symbolically, this is expressed by the differential equation
\(i = C { dv \over dt} ,\)
where *C* is the capacitance of the capacitor, 1=*i* = *i*(*t*) is the electric current through the capacitor as a function of time, and 1=*v* = *v*(*t*) is the voltage across the terminals of the capacitor, also as a function of time.

Taking the Laplace transform of this equation, we obtain
\(I(s) = C(s V(s) - V_0),\)
where
\(\begin{align}
  I(s) &= \mathcal{L} \{ i(t) \},\\
  V(s) &= \mathcal{L} \{ v(t) \},
\end{align}\)
and
\(V_0 = v(0).\)

Solving for *V*(*s*) we have
\(V(s) = { I(s) \over sC } + { V_0 \over s }.\)

The definition of the complex impedance *Z* (in ohms) is the ratio of the complex voltage *V* divided by the complex current *I* while holding the initial state *V*<sub>0</sub> at zero:
\(Z(s) = \left. { V(s) \over I(s) } \right|_{V_0 = 0}.\)

Using this definition and the previous equation, we find:
\(Z(s) = \frac{1}{sC},\)
which is the correct expression for the complex impedance of a capacitor. In addition, the Laplace transform has large applications in control theory.

### Impulse response
Consider a linear time-invariant system with transfer function
\(H(s) = \frac{1}{(s + \alpha)(s + \beta)}.\)

The impulse response is the inverse Laplace transform of this transfer function:
\(h(t) = \mathcal{L}^{-1}\{H(s)\}.\)

; Partial fraction expansion :

To evaluate this inverse transform, we begin by expanding *H*(*s*) using the method of partial fraction expansion,
\(\frac{1}{(s + \alpha)(s + \beta)} = { P \over s + \alpha } + { R \over s+\beta }.\)

The unknown constants *P* and *R* are the residues located at the corresponding poles of the transfer function. Each residue represents the relative contribution of that singularity to the transfer function's overall shape.

By the residue theorem, the inverse Laplace transform depends only upon the poles and their residues. To find the residue *P*, we multiply both sides of the equation by *s* + *α* to get
\(\frac{1}{s + \beta} = P  + { R (s + \alpha) \over s + \beta }.\)

Then by letting 1=*s* = −*α*, the contribution from *R* vanishes and all that is left is
\(P = \left.{1 \over s+\beta}\right|_{s=-\alpha} = {1 \over \beta - \alpha}.\)

Similarly, the residue *R* is given by
\(R = \left.{1 \over s + \alpha}\right|_{s=-\beta} = {1 \over \alpha - \beta}.\)

Note that
\(R = {-1 \over \beta - \alpha} = - P\)
and so the substitution of *R* and *P* into the expanded expression for *H*(*s*) gives
\(H(s)  = \left(\frac{1}{\beta - \alpha} \right) \cdot \left( { 1 \over s + \alpha } - { 1  \over s + \beta }  \right).\)

Finally, using the linearity property and the known transform for exponential decay (see *Item* #*3* in the *Table of Laplace Transforms*, above), we can take the inverse Laplace transform of *H*(*s*) to obtain
\(h(t) = \mathcal{L}^{-1}\{H(s)\} = \frac{1}{\beta - \alpha}\left(e^{-\alpha t} - e^{-\beta t}\right),\)
which is the impulse response of the system.

; Convolution :
The same result can be achieved using the convolution property as if the system is a series of filters with transfer functions 1/(*s* + *α*) and 1/(*s* + *β*). That is, the inverse of
\(H(s) = \frac{1}{(s + \alpha)(s + \beta)} = \frac{1}{s+\alpha} \cdot \frac{1}{s + \beta}\)
is
\(\mathcal{L}^{-1}\! \left\{ \frac{1}{s + \alpha} \right\} * \mathcal{L}^{-1}\! \left\{\frac{1}{s + \beta} \right\} = e^{-\alpha t} * e^{-\beta t} = \int_0^t e^{-\alpha x}e^{-\beta (t - x)}\, dx = \frac{e^{-\alpha t}-e^{-\beta t}}{\beta - \alpha}.\)

### Phase delay
{| class="wikitable"
|-
 ! scope="col" | Time function
 ! scope="col" | Laplace transform
|-
 | \(\sin{(\omega t + \varphi)}\)
 | \(\frac{s\sin(\varphi) + \omega \cos(\varphi)}{s^2 + \omega^2}\)
|-
 | \(\cos{(\omega t + \varphi)}\)
 | \(\frac{s\cos(\varphi) - \omega \sin(\varphi)}{s^2 + \omega^2}.\)
|}

Starting with the Laplace transform,
\(X(s) = \frac{s\sin(\varphi) + \omega \cos(\varphi)}{s^2 + \omega^2}\)
we find the inverse by first rearranging terms in the fraction:
\(\begin{align}
  X(s) &= \frac{s \sin(\varphi)}{s^2 + \omega^2} + \frac{\omega \cos(\varphi)}{s^2 + \omega^2} \\
       &= \sin(\varphi) \left(\frac{s}{s^2 + \omega^2} \right) + \cos(\varphi) \left(\frac{\omega}{s^2 + \omega^2} \right).
\end{align}\)

We are now able to take the inverse Laplace transform of our terms:
\(\begin{align}
  x(t) &= \sin(\varphi) \mathcal{L}^{-1}\left\{\frac{s}{s^2 + \omega^2} \right\} + \cos(\varphi) \mathcal{L}^{-1}\left\{\frac{\omega}{s^2 + \omega^2} \right\} \\
       &= \sin(\varphi)\cos(\omega t) + \cos(\varphi)\sin(\omega t).
\end{align}\)

This is just the sine of the sum of the arguments, yielding:
\(x(t) = \sin (\omega t + \varphi).\)

We can apply similar logic to find that
\(\mathcal{L}^{-1} \left\{ \frac{s\cos\varphi - \omega \sin\varphi}{s^2 + \omega^2} \right\} = \cos{(\omega t + \varphi)}.\)

### Statistical mechanics
In statistical mechanics, the Laplace transform of the density of states \(g(E)\) defines the partition function. That is, the canonical partition function \(Z(\beta)\) is given by
\(Z(\beta) = \int_0^\infty e^{-\beta E}g(E)\,dE\)
and the inverse is given by
\(g(E) = \frac{1}{2\pi i} \int_{\beta_0-i\infty}^{\beta_0+i\infty} e^{\beta E}Z(\beta) \, d\beta\)

### Spatial (not time) structure from astronomical spectrum
The wide and general applicability of the Laplace transform and its inverse is illustrated by an application in astronomy which provides some information on the *spatial distribution* of matter of an astronomical source of radiofrequency thermal radiation too distant to resolve as more than a point, given its flux density spectrum, rather than relating the *time* domain with the spectrum (frequency domain).

Assuming certain properties of the object, e.g. spherical shape and constant temperature, calculations based on carrying out an inverse Laplace transformation on the spectrum of the object can produce the only possible model of the distribution of matter in it (density as a function of distance from the center) consistent with the spectrum. When independent information on the structure of an object is available, the inverse Laplace transform method has been found to be in good agreement.

### Birth and death processes
Consider a random walk, with steps \(\{+1,-1\}\) occurring with probabilities . Suppose also that the time step is a Poisson process, with parameter .  Then the probability of the walk being at the lattice point \(n\) at time \(t\) is

\[
P_n(t) = \int_0^t\lambda e^{-\lambda(t-s)}(pP_{n-1}(s) + qP_{n+1}(s))\,ds\quad (+e^{-\lambda t}\quad\text{when}\ n=0).
\]

This leads to a system of integral equations (or equivalently a system of differential equations).  However, because it is a system of convolution equations, the Laplace transform converts it into a system of linear equations for

\[
\pi_n(s) = \mathcal L(P_n)(s),
\]

namely:

\[
\pi_n(s) = \frac{\lambda}{\lambda+s}(p\pi_{n-1}(s) + q\pi_{n+1}(s))\quad (+\frac1{\lambda + s}\quad \text{when}\ n=0)
\]

which may now be solved by standard methods.

### Tauberian theory
The Laplace transform of the measure \(\mu\) on \([0,\infty)\) is given by

\[
\mathcal L\mu(s) = \int_0^\infty e^{-st}d\mu(t).
\]

It is intuitively clear that, for small , the exponentially decaying integrand will become more sensitive to the concentration of the measure \(\mu\) on larger subsets of the domain.  To make this more precise, introduce the distribution function:

\[
M(t) = \mu([0,t)).
\]

Formally, we expect a limit of the following kind:

\[
\lim_{s\to 0^+}\mathcal L\mu(s) = \lim_{t\to\infty} M(t).
\]

Tauberian theorems are theorems relating the asymptotics of the Laplace transform, as , to those of the distribution of \(\mu\) as .  They are thus of importance in asymptotic formulae of probability and statistics, where often the spectral side has asymptotics that are simpler to infer.

Two Tauberian theorems of note are the Hardy–Littlewood Tauberian theorem and  Wiener's Tauberian theorem.  The Wiener theorem generalizes the Ikehara Tauberian theorem, which is the following statement:

Let  be a non-negative, monotonic nondecreasing function of , defined for . Suppose that

\[
f(s)=\int_0^\infty A(x) e^{-xs}\,dx
\]

converges for  to the function  and that, for some non-negative number ,

\[
f(s) - \frac{c}{s-1}
\]

has an extension as a continuous function for .
Then the limit as  goes to infinity of {{tmath|e^{-x} A(x)}} is equal to .

This statement can be applied in particular to the logarithmic derivative of Riemann zeta function, and thus provides an extremely short way to prove the prime number theorem.

## See also


* Analog signal processing
* Bernstein's theorem on monotone functions
* Continuous-repayment mortgage
* Dirichlet integral
* Differential equation
* Generating function
* Hamburger moment problem
* Hardy–Littlewood Tauberian theorem
* Laplace–Carson transform
* Moment-generating function
* Nonlocal operator
* Partial fraction decomposition
* Post's inversion formula
* Signal-flow graph
* Transfer function
* Z-transform


## Notes


## References
### Modern
*
*
*
*
*
*
*

### Historical

*
*
* , Chapters 3–5
*
*
*

## Further reading
*
*
*
*
*
*
*
*
*  |ref=none }} – see Chapter VI. The Laplace transform
*
*
*

## External links


*
* Online Computation of the transform or inverse transform, wims.unice.fr
* Tables of Integral Transforms at EqWorld: The World of Mathematical Equations.
*
* Good explanations of the initial and final value theorems
* Laplace Transforms at MathPages
* Computational Knowledge Engine allows to easily calculate Laplace Transforms and its inverse Transform.
* Laplace Calculator to calculate Laplace Transforms online easily.
* Code to visualize Laplace Transforms and many example videos.

