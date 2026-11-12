# Stability regions

The stability region is the set of step parameters for which a method damps the test mode. Its shape explains why some methods suit stiff decay and others suit oscillatory problems.

## Definition

The region is the set of complex numbers where the stability function has modulus at most one. It is drawn in the complex plane and contains the origin for every consistent method.

## Explicit methods

An explicit method has a polynomial stability function, so its region is bounded. The fourth-order Runge-Kutta region reaches about minus two and eight tenths along the negative real axis and does not cover the imaginary axis beyond a small interval.

## Reading the boundary

Where the boundary crosses the imaginary axis determines the stability of an undamped oscillator. A method whose region includes a segment of that axis can integrate a conservative system without artificial damping.

To check your understanding, write the modulus condition, identify the polynomial and find the imaginary crossing.

The region is a picture of a method's competence, and the two axes answer different questions: the negative real axis for decay, the imaginary axis for oscillation.
