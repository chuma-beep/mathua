# The gamma function and its integral

The gamma function extends the factorial to non-integer arguments through an integral that converges for positive real part.

## The integral definition

The definition is the integral from 0 to infinity of \(t^{z-1}e^{-t}\). The exponential makes the upper end converge. The power makes the lower end converge when the real part of \(z\) is positive.

At \(z=4\) the integrand is \(t^{3}e^{-t}\), whose integral is 6. That matches the factorial of 3.

## The functional equation

Integrating by parts gives \(\Gamma(z+1)=z\Gamma(z)\). The boundary terms vanish at both ends. What remains is the recursion.

With \(\Gamma(1)=1\) the equation gives the factorial at every positive integer, so \(\Gamma(5)=24\). It also continues the function to the left of the origin, where poles appear at the nonpositive integers.

## Values at half-integers

The value at one half is \(\sqrt{\pi}\), about 1.7725. The equation then gives \(\Gamma(3/2)=\sqrt{\pi}/2\), about 0.8862. It also gives \(\Gamma(5/2)=3\sqrt{\pi}/4\), about 1.3293.

Each step multiplies by the previous half-integer. The pattern is a product of odd numbers over a power of 2, times \(\sqrt{\pi}\).

The gamma function is the continuous version of the factorial. Its half-integer values appear in the volumes of spheres and in the beta integral.
