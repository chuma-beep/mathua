> Content sourced from [Algebrica](https://algebrica.org/acceleration/) — CC BY-NC 4.0

## Introduction

Kinematics is the study of the motion of objects, describing their position, [velocity](<../velocity>), and acceleration over time. Before diving into a detailed explanation of these phenomena, it is essential to introduce some key concepts.

  * A material point is an idealized object whose size is considered negligible relative to the distances involved in its motion.
  * The trajectory is the path traced by a material point as it moves through space.
  * A motion is said to be rectilinear if its trajectory lies along a straight line.


If a material point moves along a [straight-line](<../lines>) path under constant acceleration, meaning that the rate of change of velocity remains uniform over time, the motion is called **uniformly accelerated rectilinear motion**.

## Acceleration

Let us consider a particle moving along a straight-line trajectory, where the position as a function of time is not described by a [linear equation](<../linear-equations/>). Let \\( P_1 \\) and \\( P_2 \\) denote two positions of the material point along the \\( x \\)-axis at times \\( t_1 \\) and \\( t_2 \\), respectively. We denote by \\( \mathbf{v}_1 \\) and \\( \mathbf{v}_2 \\) the corresponding velocity [vectors](<../vectors/>), with \\( \mathbf{v}_1 \neq \mathbf{v}_2 \\). The **vector acceleration** is defined as the following limit:

\\[\mathbf{a} = \lim_{\Delta t \to 0} \frac{\Delta \mathbf{v}}{\Delta t} = \frac{d\mathbf{v}}{dt} \\]

We have seen, by analyzing the [velocity](<../velocity>), that:

\\[\lim_{\Delta t \to 0} \frac{\Delta \mathbf{r}}{\Delta t} = \frac{d\mathbf{r}}{dt} = \mathbf{v} \\]

Thus, we have:

\\[\mathbf{a} = \frac{d}{dt}\left( \frac{d\mathbf{r}}{dt} \right) = \frac{d^2 \mathbf{r}}{dt^2} \\]


Starting from the general expression of acceleration it is possible to introduce the concept of **tangential acceleration** As a point \\( P \\) travels along a given path, the acceleration vector \\( \mathbf{a} \\) can be broken down into two components:

  * One tangential to the trajectory.
  * One normal to the trajectory (also called centripetal acceleration that points toward the center of the curvature of the path).


The tangential acceleration, denoted by \\( \mathbf{a}_t \\), corresponds to the variation of the speed over time. It is defined as:

\\[a_t = \frac{dv}{dt} = \mathbf{i} \, a_t \\]

where \\( v \\) represents the magnitude of the velocity vector \\( \mathbf{v} \\) and \\(\mathbf{i}\\) represents a directed and oriented vector.

![The acceleration vector consists of two parts: a tangential component and a normal component.](https://algebrica.org/wp-content/uploads/resources/images/acceleration-1.png)

  * If the magnitude of the velocity changes, there is tangential acceleration \\((a_t \neq 0)\\).
  * If the magnitude of the velocity remains constant, the tangential acceleration is zero \\((a_t = 0)\\).


Uniformly accelerated motion is a type of motion in which the tangential acceleration \\( a_t \\) is constant at every point and equal to the average acceleration over any time [interval](<../intervals/>). We have:

\\[\frac{v - v_0}{t-t_0} = a_t\\]

Starting from this formula, solving for \\( v \\) and assuming \\( t_0 = 0 \\), we obtain:

\\[v = v_0 + a_t t \\]

In this way, derived the expression for velocity based on the definition of acceleration. Starting from the expression of velocity as a function of time we can derive the equation of motion by [integrating](<../integrals>) with respect to time:

\\[y = \int_0^t v(t) \, dt = \int_0^t (v_0 + a_t t) \, dt \\]

Evaluating the integral, we obtain:

\\[y = v_0 t + \frac{1}{2} a_t t^2 \\]

where \\( y \\) represents the displacement of the material point along the trajectory as a function of time.

Kinematics

Kinematics studies motion through trajectories, velocities, and accelerations.

1.7k

[Uniform Linear Motion: Velocity](https://algebrica.org/velocity/)

1.1k

[Simple Harmonic Motion](https://algebrica.org/simple-harmonic-motion/)
