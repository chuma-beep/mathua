# Weight decay as a constraint

Weight decay can be read as a penalty or as a projection. The two readings agree in the simplest case and diverge once the step is adaptive.

## The penalty view

Adding a squared norm to the objective produces a gradient term proportional to the parameters. The strength of the shrink then depends on how the optimiser scales that term.

## The projection view

Shrinking the parameters by a fixed factor at each step is a soft projection towards the origin. It is independent of the gradient and therefore of the adaptive denominator.

## Choosing the strength

The decay rate sets the time scale on which a parameter without gradient signal returns to zero. It should be slow compared with learning and fast compared with the run length.

To check your understanding, add a squared norm, shrink by a factor and set a time scale.

Separating the two views makes it clear why decoupled decay is preferred with adaptive methods and why the same coefficient behaves differently under plain gradient descent.
