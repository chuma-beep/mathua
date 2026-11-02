# Forced oscillations

Forced oscillations is part of the differential equations curriculum. This lesson develops the idea in three sections, from the defining picture through the main computation to the interpretation.

## Driven Oscillators

A forced oscillator adds an external driving term to the spring-mass equation. The homogeneous part describes the natural motion of the system, while the forcing term describes what is being done to it from outside. For a sinusoidal force the equation has the form m y double prime + c y prime + k y = F0 cos(omega t). The constants m, c and k are mass, damping and stiffness, and F0 and omega describe the force.

The general solution is a natural part plus a particular part. With positive damping the natural part decays, so the long-run behaviour is governed by the particular solution. That steady-state response has the same frequency as the driving force but usually a different amplitude and phase. The transient and steady-state labels refer to time behaviour, not to different equations.

The natural frequency is sqrt(k/m) when damping is absent. It depends only on the system, while the driving frequency is chosen by the external force. The relationship between those two frequencies determines whether the response is small, large or resonant.

## Resonance

Resonance occurs when the driving frequency approaches a natural frequency of the undamped system. In the undamped case the particular solution then contains a factor of t, so its amplitude grows without bound. This is why a repeated root in the characteristic equation changes the form of the trial solution: the usual cosine guess is already a solution of the homogeneous equation and cannot work again.

With damping the response remains finite, but it can still be very large near the resonant frequency. The damping coefficient limits the peak amplitude and shifts it slightly below the undamped natural frequency. A lightly damped system has a sharp peak, while heavy damping produces a broad and low response.

Resonance is not only a mathematical curiosity. It explains why soldiers break step on a bridge, why a washing machine shakes at a particular spin speed and why tuned circuits respond strongly to one frequency. Engineers either avoid the resonant range or add damping so that the response stays within safe limits.

## The Particular Solution

For a sinusoidal force and no damping, the trial particular solution is A cos(omega t) plus B sin(omega t). Substituting into the equation and matching coefficients gives A and B. At resonance the trial must be multiplied by t, producing terms such as t sin(omega t). The factor t is the signature of a repeated forcing frequency.

With damping, the same trial form works, but the algebra is longer because the first derivative contributes sine and cosine terms. The amplitude can be written as F0 divided by the square root of (k - m omega^2)^2 + (c omega)^2. That denominator measures how far the driving frequency is from resonance and how much damping is present.

After finding the particular solution, add the homogeneous solution and apply the initial conditions to the combined expression. A common mistake is to apply initial conditions to the particular solution alone. The constants in the homogeneous part are what make the full solution satisfy the starting displacement and velocity.
