# The wave equation in one dimension

The wave equation describes a vibrating string whose solutions travel without changing shape.

## The travelling wave form

With unit speed the equation is \(u_{tt}=u_{xx}\), whose general solution is \(f(x-t)+g(x+t)\).

One profile moves right and the other left, both at speed 1. Each solves the equation on its own.

Take \(f(x)=x^{2}\). Then \(f(x-t)=(x-t)^{2}\) has second time derivative 2 and second space derivative 2, so it solves the equation.

## D'Alembert's formula

The solution with initial position \(p\) and velocity \(v\) is \(\frac{p(x-t)+p(x+t)}{2}+\frac{1}{2}\int_{x-t}^{x+t}v(s)\,ds\).

The first part splits the initial shape into two halves. The second integrates the velocity. With zero velocity and \(p(x)=\sin x\) the formula gives \(\sin x\cos t\).

At \(x=\pi/6\) and \(t=\pi/3\) the factors are 0.5 and 0.5, so \(u=0.25\).

## A worked initial profile

A string released from rest keeps its shape and oscillates in place, since the two halves travel and then recombine.

For the sine profile the value at time zero is \(\sin x\), as required. The peak falls from 1 to 0.5 at \(t=\pi/3\).

A plucked string behaves differently, because its corner splits into two travelling corners.

The two descriptions agree everywhere. The choice between them is one of convenience.
