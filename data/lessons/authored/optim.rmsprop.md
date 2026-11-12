# RMSProp

RMSProp divides the gradient by a decaying root-mean-square of recent gradients. It fixes Adagrad's shrinking step on long runs.

## The denominator

A decayed average of squared gradients is updated at every step. Its square root, plus a small constant, divides the gradient, so coordinates with consistently large gradients receive smaller moves.

## The decay

A decay near one gives a long memory and a smooth denominator. A smaller decay reacts faster to recent changes and is more sensitive to spikes in the gradient.

## Comparison with Adagrad

Adagrad accumulates every past square, so its denominator grows without bound and the step decays to zero. RMSProp forgets old squares, so the denominator settles and the step stays usable.

To check your understanding, average squared gradients, choose a long memory and forget old squares.

RMSProp keeps the per-coordinate scaling of Adagrad and adds forgetting, which is what makes it practical for long training runs.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
