# The Winding Number

## Winding around a point

The winding number of a closed curve around a point counts how many times the curve travels counterclockwise around that point. It is an integer, positive for counterclockwise motion, negative for clockwise motion and zero when the point lies outside the curve. The winding number is constant on each connected component of the complement of the curve.

For a curve parametrized by z of t, the winding number around a is the change in the argument of z of t minus a, divided by two pi. This definition is analytic rather than topological, but it agrees with the intuitive count. The curve must be closed and must not pass through the point.

Winding numbers are the reason contour integrals detect holes. The integral of one over z minus a around a closed curve equals two pi i times the winding number. Cauchy's theorem therefore holds for curves that wind zero times around every singularity, even if the curve is complicated.

## Computing the winding number

For a circle traversed once counterclockwise, the winding number around its centre is one. Around a point outside, it is zero. Reversing the orientation changes the sign. These two facts cover most textbook examples and provide a quick check before computing a full argument change.

For a curve that winds several times, the winding number is the net number of turns. A figure-eight curve can wind once around one lobe and once in the opposite direction around the other, giving zero for a point between the lobes. The winding number is additive when curves are concatenated, which makes it a homotopy invariant.

To compute it from a parametrization, track the argument continuously and count crossings of the negative real axis. Each counterclockwise crossing adds one, and each clockwise crossing subtracts one. The total is the winding number. This method works even when the curve is not a circle, as long as the point is not on the curve.

## Winding and argument change

The argument principle is the statement that the winding number of the image curve equals the zero count minus the pole count. This connects the topological winding number to the algebraic data of a meromorphic function. It is one of the clearest examples of topology influencing analysis.

The winding number also classifies maps from the circle to the punctured plane up to homotopy. Two such maps are homotopic exactly when they have the same winding number. This is the first nontrivial homotopy group computation and the prototype for degree theory in higher dimensions.

In applications, winding numbers detect whether a curve encloses a point and how many times. They are used in contour integration, in the stability criterion for feedback systems and in the definition of the index of a vector field. The integer is robust: small perturbations of the curve do not change it unless the curve crosses the point.

