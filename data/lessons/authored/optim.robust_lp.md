# Robust linear programs

A robust linear program protects a constraint against every coefficient inside an uncertainty set. A norm ball turns that infinite family into one cone constraint.

## Uncertainty as a norm ball

The coefficient vector is allowed to move by at most a radius in some norm. The constraint must hold at the worst point of that ball, not at its centre.

## The worst case

For a ball in the dual norm the worst shift of the left side is the radius times the norm of the variable. The worst case therefore depends on the variable only through its norm.

## The robust counterpart

Moving the worst case to the left of the bound gives one second-order cone constraint. The radius sets its right-hand side and the norm sets its tail.

To check your understanding, compute the norm of the variable, subtract the nominal value from the bound and evaluate the worst case.

The robust counterpart is conservative on purpose. It holds for every coefficient in the ball, so it can be tighter than the nominal model at the centre.
