> Content sourced from [Algebrica](https://algebrica.org/simple-harmonic-motion/) — CC BY-NC 4.0

## What is simple harmonic motion

A simple harmonic motion is a straight-line motion obtained by projecting the uniform circular motion of a body onto a fixed diameter of the circle. In this case, a material point repeatedly moves along the diameter, oscillating back and forth, and returning to the same position at regular [intervals](<../intervals/>) of time equal to the period of the motion.

![](/diagrams/algebrica/harmonic-motion-1.png)

This projection creates a smooth and continuous oscillation, where the displacement from the center varies [sinusoidally](<../sine-and-cosine>) with time. The position as a function of time for simple harmonic motion is given by the equation:

\\[x(t) = A \sin(\omega t) \\]

  * \\( A \\) is the amplitude, the maximum displacement from the equilibrium position (the function \\( \sin(\omega t) \\) oscillates between the values \\( +1 \\) and \\( -1 \\)).
  * \\( \omega \\) is the angular frequency (in radians per second).
  * \\( t \\) is the time variable.


Since the function is periodic with period \\( T \\), meaning it repeats itself from \\( t \\) to \\( t+T \\), the argument of the trigonometric function must change by \\( 2\pi \\).  
Therefore, we have:

\\[\omega(t+T) - \omega t = 2\pi \\]

From this, we derive the fundamental relation: \\[\omega = \frac{2\pi}{T} \\] where \\( \nu \\) represents the frequency of the harmonic motion, and is related to the period by:

\\[\nu = \frac{1}{T} \\]


In simple harmonic motion where the center of oscillation does not coincide with the origin but is located at \\( x_0 \\), the general equation of motion is:

\\[x - x_0 = A \sin(\omega t + \varphi) \\]

  * \\( x(t) \\) is the position of the material point at time ( t ),
  * \\( x_0 \\) is the center around which the motion oscillates,
  * \\( A \\) is the amplitude,
  * \\( \omega \\) is the angular frequency,
  * \\( \varphi \\) is the initial phase.


## Example 1

For example, calculate the equation of a simple harmonic motion with period \\( T = 3 \, \text{s} \\) and amplitude \\( A = 0.5 \, \text{m} \\).

We find the angular frequency.

\\[\omega = \frac{2\pi}{T} = \frac{2\pi}{3} \, \text{rad/s} \\]

Thus, the motion equation is:

\\[x(t) = 0.5 \sin\left( \frac{2\pi}{3} t \right) \\]

## Velocity

To find the instantaneous [scalar velocity](<../velocity>) \\( v(t) \\), we [differentiate](<../derivatives>) \\( x(t) \\) with respect to time:

\\[v(t) = \frac{dx}{dt} \\]

Applying the derivative: \\[\frac{d}{dt} \left( A \sin(\omega t) \right) = A \omega \cos(\omega t) \\]

Thus, the instantaneous velocity of the material point is: \\[v(t) = A \omega \cos(\omega t) \\]

This expression shows that the velocity varies over time following a [cosine function](<../cosine-function>), reaching its maximum when the particle passes through the equilibrium position.

  * The velocity in simple harmonic motion is not constant: it changes over time following the trend of a cosine function.
  * When the body passes through the equilibrium position (that is, \\(x = 0\\), the center of motion), the velocity is maximum.
  * When the body reaches the extremes (that is, at the points \\(x = +A\\) or \\(x = -A\\)), the velocity is zero.
  * The velocity, being a function of the cosine, is phase-advanced with respect to the displacement \\(x\\) in simple harmonic motion.


## Example 2

For example, calculate the instantaneous velocity of a simple harmonic motion with period \\( T = 3 \, \text{s} \\) and amplitude \\( A = 0.5 \, \text{m} \\).


We find the angular frequency.

\\[\omega = \frac{2\pi}{T} = \frac{2\pi}{3} \, \text{rad/s} \\]

The position function is:

\\[x(t) = 0.5 \sin\left( \frac{2\pi}{3} t \right) \\]

Differentiating with respect to time to find the velocity:

\\[v(t) = \frac{dx}{dt} = 0.5 \times \frac{2\pi}{3} \cos\left( \frac{2\pi}{3} t \right) \\]

Simplifying we obtain:

\\[v(t) = \frac{\pi}{3} \cos\left( \frac{2\pi}{3} t \right) \\]

## Acceleration

To find the instantaneous [acceleration](<../acceleration>) \\( a(t) \\), we [differentiate](<../derivatives>) the velocity function \\( v(t) \\) with respect to time:

\\[a(t) = \frac{dv}{dt} \\]

Applying the derivative:

\\[\frac{d}{dt} \left( A \omega \cos(\omega t) \right) = -A \omega^2 \sin(\omega t) \\]

Thus, the instantaneous acceleration of the material point is:

\\[a(t) = -A \omega^2 \sin(\omega t) \\]

This expression shows that the acceleration varies over time following a [sine function](<../sine-function>), reaching its maximum magnitude when the particle is at the maximum displacement from the equilibrium position.

  * The acceleration in simple harmonic motion is not constant: it changes over time following the trend of a sine function.
  * When the body passes through the equilibrium position (that is, \\(x = 0\\), the center of motion), the acceleration is zero.
  * When the body reaches the extremes (that is, at the points \\(x = +A\\) or \\(x = -A\\)), the acceleration reaches its maximum magnitude.


The acceleration, being a function of the sine, is in phase opposition with respect to the displacement \\(x\\). In simple harmonic motion the acceleration has only a normal component \\( a_n \\) and it is centripetal, always directed toward the center of oscillation. This means that as the particle moves, the acceleration acts continuously to pull it back toward the equilibrium position.

> In uniformly [accelerated rectilinear motion](<../acceleration>), the acceleration [vector](<../vectors/>) has two components: a normal component \\( a_n \\) and a tangential component \\( a_t \\). Since the trajectory is still a straight line, the normal component \\( a_n = 0 \\). However, the tangential component \\( a_t \\) is constant and nonzero, representing a steady change in the velocity’s magnitude along the direction of motion.

## Example 3

For example, calculate the acceleration in a simple harmonic motion with period \\( T = 3 , \text{s} \\) and amplitude \\( A = 0.5 \, \text{m} \\).


Let’s review the steps followed in the previous examples. First, we find the angular frequency.

\\[\omega = \frac{2\pi}{T} = \frac{2\pi}{3} \, \text{rad/s} \\]

Starting from the motion equation:

\\[x(t) = 0.5 \sin\left( \frac{2\pi}{3} t \right) \\]

The velocity is:

\\[v(t) = 0.5 \times \frac{2\pi}{3} \cos\left( \frac{2\pi}{3} t \right) = \frac{\pi}{3} \cos\left( \frac{2\pi}{3} t \right) \\]

Differentiating again, we find the acceleration:

\\[a(t) = \frac{dv}{dt} = -\left( \frac{\pi}{3} \times \frac{2\pi}{3} \right) \sin\left( \frac{2\pi}{3} t \right) \\]

Simplifying:

\\[a(t) = -\frac{2\pi^2}{9} \sin\left( \frac{2\pi}{3} t \right) \\]
