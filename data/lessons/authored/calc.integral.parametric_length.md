# Arc length of parametric curves

## The speed of a parametrisation

A parametrisation \(x(t), y(t)\) moves along a curve with speed \(\sqrt{x'(t)^2 + y'(t)^2}\). For \(x = t^2\) and \(y = t^3\) the derivatives at \(t = 1\) are 2 and 3, so the speed there is \(\sqrt{13} \approx 3.606\). Speed is a rate per unit of the parameter.

## The length integral

The length is the integral of speed: \(L = \int_a^b \sqrt{x'^2 + y'^2}\,dt\). A circle of radius 3 has \(x = 3\cos t\) and \(y = 3\sin t\), with speed 3 and length \(6\pi \approx 18.85\) over one turn. The result matches the circumference formula.

## A straight segment

A straight segment is the check case. With \(x = 2t\) and \(y = 3t + 1\) on the interval from 0 to 1 the speed is constant at \(\sqrt{13}\), so the length is 3.606. The endpoints are \((0,1)\) and \((2,4)\). Their distance is \(\sqrt{4+9} = \sqrt{13}\), so both methods agree.

A complete answer names the relation, substitutes the data and states the result in the language of the problem. When a second method is available, compare the two results before moving on.
