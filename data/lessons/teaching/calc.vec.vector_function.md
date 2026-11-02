# Vector-Valued Functions

A vector-valued function assigns a vector to each value of a parameter. It describes a curve in the plane or in space.

## Parametrised curves

A vector-valued function has the form \(\mathbf r(t)=\langle x(t),y(t),z(t)\rangle\). Each component is an ordinary function of the parameter \(t\). The endpoint of the vector traces a curve as \(t\) varies.

For \(\mathbf r(t)=\langle 3\cos t,3\sin t,4t\rangle\), the x and y components trace a circle of radius 3 while the z component rises steadily. The result is a helix.

## Derivatives and velocity

Differentiate a vector-valued function component by component: \(\mathbf r'(t)=\langle x'(t),y'(t),z'(t)\rangle\). The derivative is the velocity vector, tangent to the curve at the point.

For \(\mathbf r(t)=\langle 3\cos t,3\sin t,4t\rangle\), the velocity is \(\langle -3\sin t,3\cos t,4\rangle\). At \(t=0\) the velocity is \(\langle 0,3,4\rangle\), which points upward and in the positive y direction.

## Speed and arc length

Speed is the magnitude of the velocity: \(\|\mathbf r'(t)\|\). For the helix above, speed is \(\sqrt{9\sin^2t+9\cos^2t+16}=\sqrt{25}=5\). The speed is constant even though the velocity direction changes.

Arc length from \(t=a\) to \(t=b\) is \(\int_a^b\|\mathbf r'(t)\|\,dt\). For constant speed 5 over an interval of length 2, the arc length is 10.
