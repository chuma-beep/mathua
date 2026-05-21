> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Laplace_transform) — CC BY-SA 4.0

# Laplace transforms

In mathematics, the **Laplace transform**, named after Pierre-Simon Laplace (), is an integral transform that converts a function of a real variable (usually , in the *time domain*) to a function of a complex variable \(s\) (in the complex-valued frequency domain, also known as ***s*-domain** or ***s*-plane**). The functions are often denoted using a lowercase symbol for the time-domain function and the corresponding uppercase symbol for the frequency-domain function, e.g. \(x(t)\) and. The transform is useful for converting differentiation and integration in the time domain into the algebraic operations multiplication and division in the Laplace domain (analogous to how logarithms are useful for simplifying multiplication and division into addition and subtraction). This gives the transform many applications in science and engineering, mostly as a tool for solving linear differential equations and dynamical systems by replacing ordinary differential equations and integral equations with algebraic polynomial equations, and by replacing convolution with multiplication.

For example, through the Laplace transform, the equation of the simple harmonic oscillator (Hooke's law) \(x''(t) + kx(t)=0\) is converted into the algebraic equation \(s^2X(s) - s x(0) - x'(0) +kX(s) = 0,\) which incorporates the initial conditions \(x(0)\) and , and can be solved for the unknown function. Once solved, the inverse Laplace transform can be used to transform it to the original domain. This is often aided by referencing tables such as that given below.

The Laplace transform is defined (for suitable functions ) by the integral

\[
\mathcal{L}\{f\}(s) = \int_0^\infty f(t) e^{-st} \, dt,
\]

where is a complex number.

The Laplace transform is related to many other transforms. It is essentially the same as the Mellin transform and is closely related to the Fourier transform. Unlike for the Fourier transform, the Laplace transform of a function is often an analytic function, meaning that it can be expressed as a power series that converges locally, the coefficients of which represent the moments of the original function. Moreover, the techniques of complex analysis, especially contour integrals, can be used for simplifying calculations.

## History

The Laplace transform is named after mathematician and astronomer Pierre-Simon, Marquis de Laplace, who used a similar transform in his work on probability theory. Laplace wrote extensively about the use of generating functions (1814), and the integral form of the Laplace transform evolved naturally as a result.

Laplace's use of generating functions was similar to what is now known as the z-transform, and he gave little attention to the continuous variable case which was discussed by Niels Henrik Abel.

From 1744, Leonhard Euler investigated integrals of the form
\(z = \int X(x) e^{ax}\, dx \quad\text{ and }\quad z = \int X(x) x^A \, dx\)
as solutions of differential equations, introducing in particular the gamma function. Joseph-Louis Lagrange was an admirer of Euler and, in his work on integrating probability density functions, investigated expressions of the form
\(\int X(x) e^{- a x } a^x\, dx,\)
which resembles a Laplace transform.

These types of integrals seem first to have attracted Laplace's attention in 1782, where he was following in the spirit of Euler in using the integrals themselves as solutions of equations. However, in 1785, Laplace took the critical step forward when, rather than simply looking for a solution in the form of an integral, he started to apply the transforms in the sense that was later to become popular. He used an integral of the form
\(\int x^s \varphi (x)\, dx,\)
akin to a Mellin transform, to transform the whole of a difference equation, in order to look for solutions of the transformed equation. He then went on to apply the Laplace transform in the same way and started to derive some of its properties, beginning to appreciate its potential power.

Laplace also recognised that Joseph Fourier's method of Fourier series for solving the diffusion equation could only apply to a limited region of space, because those solutions were periodic. In 1809, Laplace applied his transform to find solutions that diffused indefinitely in space. In 1821, Cauchy developed an operational calculus for the Laplace transform that could be used to study linear differential equations in much the same way the transform is now used in basic engineering. This method was popularized, and perhaps rediscovered, by Oliver Heaviside around the turn of the century.

Bernhard Riemann used the Laplace transform in his 1859 paper *On the number of primes less than a given magnitude*, in which he also developed the inversion theorem. Riemann used the Laplace transform to develop the functional equation of the Riemann zeta function, and his method is still used to relate the modular transformation law of the Jacobi theta function, which is readily proved via Poisson summation, to the functional equation.

Hjalmar Mellin was among the first to study the Laplace transform rigorously in the Karl Weierstrass school of analysis, and apply it to the study of differential equations and special functions, at the turn of the 20th century. At around the same time, Heaviside was busy with his operational calculus. Thomas Joannes Stieltjes considered a generalization of the Laplace transform connected to his work on moments. Other contributors in this time period included Mathias Lerch, Oliver Heaviside, and Thomas Bromwich.

In 1929, Vannevar Bush and Norbert Wiener published *Operational Circuit Analysis* as a text for engineering analysis of electrical circuits, applying both Fourier transforms and operational calculus, and in which they included one of the first predecessors of the modern table of Laplace transforms.
In 1934, Raymond Paley and Norbert Wiener published the important work *Fourier transforms in the complex domain*, about what is now called the Laplace transform (see below). Also during the 30s, the Laplace transform was instrumental in Godfrey Harold Hardy and John Edensor Littlewood's study of tauberian theorems, and this application was later expounded on by , who developed other aspects of the theory such as a new method for inversion. Edward Charles Titchmarsh wrote the influential *Introduction to the theory of the Fourier integral* (1937).

The current widespread use of the transform (mainly in engineering) came about during and soon after World War II, replacing the earlier Heaviside operational calculus. The advantages of the Laplace transform had been emphasized by Gustav Doetsch.

## Formal definition

The Laplace transform of a function *f*(*t*), defined for all real numbers *t* ≥ 0, is the function *F*(*s*) defined by
 \(\frac{s\cos(\varphi) - \omega \sin(\varphi)}{s^2 + \omega^2}.\)
