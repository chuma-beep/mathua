> Content sourced from [Algebrica](https://algebrica.org/velocity/) — CC BY-NC 4.0

## Introduction

Kinematics is the study of the motion of objects, describing their position, velocity, and [acceleration](<../acceleration>) over time. Before diving into a detailed explanation of these phenomena, it is essential to introduce some key concepts.

  * A material point is an idealized object whose size is considered negligible relative to the distances involved in its motion.
  * The trajectory is the path traced by a material point as it moves through space.
  * A motion is said to be rectilinear if its trajectory lies along a straight line.


If a material point moves along a straight-line path at constant velocity, meaning that the distances traveled are proportional to the time intervals taken, the motion is called **uniform rectilinear motion.**

## Velocity

Let’s consider two points, \\( x_1 \\) and \\( x_2 \\), representing the position \\( P \\) of a point at two successive moments in time, \\( t_1 \\) and \\( t_2 \\), respectively.

![](https://algebrica.org/wp-content/uploads/resources/images/velocity-1-1.png)

We can express the following relationship:

\\[\frac{x_2 - x_1}{t_2 - t_1} = v \\]

This relationship shows that the distance traveled is proportional to the elapsed time, and the ratio remains constant, equal to the magnitude \\( v \\). The **instantaneous scalar speed** is defined as the limit of the ratio as the time [interval](<../intervals/>) approaches zero:

\\[v_s = \lim_{\Delta t \to 0} \frac{\Delta x}{\Delta t} = \lim_{\Delta t \to 0} \frac{x(t + \Delta t) - x(t)}{\Delta t} \\]

where \\( v \\) represents the instantaneous speed, \\( \Delta x \\) is the displacement, and \\( \Delta t \\) is the time interval. Therefore, the instantaneous scalar speed is given by the [derivative](<../derivatives>) of \\(x = x(t)\\) with respect to time.


Let us now imagine that the position \\( P \\) of the point at time \\( t \\) is identified by the displacement [vector](<../vectors/>) \\( \vec{r} \\). The displacement from \\( P \\) to \\( P’ \\) occurs over a time interval \\( \Delta t \\) and is represented by the vector \\( \Delta \vec{r} \\). We have:

\\[\lim_{\Delta t \to 0} \frac{\Delta \mathbf{r}}{\Delta t} = \frac{d\mathbf{r}}{dt} = \mathbf{v} \\]

![](https://algebrica.org/wp-content/uploads/resources/images/velocity-2.png)

Thus, we can define the **velocity vector** as:

\\[\mathbf{v} = \frac{dx}{dt} = \mathbf{i} \frac{dx(t)}{dt} \\]

where \\(\mathbf{i}\\) represents a directed and oriented vector. The velocity vector is tangent to the trajectory at each point, oriented according to the direction of motion, and has a magnitude equal to the scalar speed.

  * In uniform rectilinear motion, the velocity vector remains constant.
  * In uniform rectilinear motion, the trajectory’s position-time equation is a [straight line](<../lines>), meaning the position varies linearly with time.


Velocity is measured in units of length multiplied by time raised to the power of \\(-1\\), and its standard unit is meters per second \\((\text{ms}^{-1})\\).

## Example

Imagine a car traveling along a straight road at a constant speed of \\( v = 20\ \mathrm{m/s} \\). Since the velocity is constant, the distance traveled by the car is directly proportional to the elapsed time. The position \\( x(t) \\) of the car at any time \\( t \\) can be expressed as:

\\[x(t) = x_0 + v t \\]

where:

  * \\( x_0 \\) is the initial position (at \\( t = 0 \\)).
  * \\( v \\) is the constant velocity.
  * \\( t \\) is the time elapsed.


This means that for every second that passes, the car moves exactly 20 meters forward, without speeding up or slowing down. Let’s summarize the data in a table:

\\[\begin{array}{c|c} \text{Time } (\mathrm{s}) & \text{Position } (\mathrm{m}) \\\ \hline 0 & 0 \\\ 1 & 20 \\\ 2 & 40 \\\ 3 & 60 \\\ 4 & 80 \\\ \vdots & \vdots \end{array} \\]

## Glossary

  * Kinematics: the study of the motion of objects, describing their position, velocity, and acceleration over time.

  * Material point: an idealized object whose size is considered negligible relative to the distances involved in its motion.

  * Trajectory: the path traced by a material point as it moves through space.

  * Rectilinear motion: motion whose trajectory lies along a straight line.

  * Uniform rectilinear motion: motion along a straight line with constant velocity, where distances traveled are proportional to time intervals.

  * Scalar speed: the limit of the ratio of displacement to time interval as the time interval approaches zero; the magnitude of the velocity vector.

  * Velocity vector: the rate of change of the displacement vector with respect to time; a vector tangent to the trajectory, oriented in the direction of motion, with a magnitude equal to the scalar speed.


## What is velocity

  * Velocity describes how fast and in what direction an object moves.
  * Scalar velocity refers to the absolute value of velocity, representing only the speed of the object without considering the direction.
  * Vector velocity is the rate of change of position with respect to time, expressed as a vector tangent to the trajectory and oriented in the direction of motion.


Kinematics

Kinematics studies motion through trajectories, velocities, and accelerations.

1.2k

[Uniformly Accelerated Motion: Acceleration](https://algebrica.org/acceleration/)

1.1k

[Simple Harmonic Motion](https://algebrica.org/simple-harmonic-motion/)
