# The heat kernel on the real line

The heat kernel is the temperature profile produced by a unit of heat released at a single point.

## The heat equation

The heat equation is \(u_{t}=u_{xx}\) for a rod with unit diffusivity.

The rate of change of temperature at a point equals the curvature of the profile there. A hot spot has negative curvature, so it cools while its neighbours warm.

The profile \(\sin x\) decays as \(e^{-t}\sin x\), so at \(t=\ln 2\) its amplitude has halved.

## The Gaussian kernel

The kernel is \(K(t,x)=\frac{1}{\sqrt{4\pi t}}e^{-x^{2}/(4t)}\) for positive time.

Its integral over the line is 1 for every time, so the total heat is conserved. Its width grows like the square root of time.

At \(t=1\) and \(x=0\) the value is \(1/\sqrt{4\pi}\), about 0.2821. At \(x=2\) the factor \(e^{-1}\) appears and the value is about 0.1037.

## The convolution formula

The solution with initial data \(f\) is \(u(t,x)=\int K(t,x-y)f(y)\,dy\).

The kernel slides across the data and averages it with a Gaussian weight. For \(f(y)=y^{2}\) the result is \(x^{2}+2t\).

At \(t=1\) and \(x=2\) that gives 6. The derivatives \(u_{t}=2\) and \(u_{xx}=2\) agree.

The kernel is also a limit, since the Gaussian narrows to a point as time falls to zero.
